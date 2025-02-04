// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package metaloop_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/common/errs2"
	"storj.io/common/memory"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
	"storj.io/storj/private/testplanet"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/metainfo/metabase"
	"storj.io/storj/satellite/metainfo/metaloop"
	"storj.io/uplink/private/multipart"
)

// TestLoop does the following
// * upload 5 remote files with 1 segment
// * (TODO) upload 3 remote files with 2 segments
// * upload 2 inline files
// * connect two observers to the metainfo loop
// * run the metainfo loop
// * expect that each observer has seen:
//    - 5 remote files
//    - 5 remote segments
//    - 2 inline files/segments
//    - 7 unique path items
func TestLoop(t *testing.T) {
	// TODO: figure out how to configure testplanet so we can upload 2*segmentSize to get two segments
	segmentSize := 8 * memory.KiB

	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 4,
		UplinkCount:      1,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Metainfo.Loop.CoalesceDuration = 1 * time.Second
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		ul := planet.Uplinks[0]
		satellite := planet.Satellites[0]
		metaLoop := satellite.Metainfo.Loop

		// upload 5 remote files with 1 segment
		for i := 0; i < 5; i++ {
			testData := testrand.Bytes(segmentSize)
			path := "/some/remote/path/" + strconv.Itoa(i)
			err := ul.Upload(ctx, satellite, "bucket", path, testData)
			require.NoError(t, err)
		}

		// (TODO) upload 3 remote files with 2 segments
		// for i := 0; i < 3; i++ {
		// 	testData := testrand.Bytes(2 * segmentSize)
		// 	path := "/some/other/remote/path/" + strconv.Itoa(i)
		// 	err := ul.Upload(ctx, satellite, "bucket", path, testData)
		// 	require.NoError(t, err)
		// }

		// upload 2 inline files
		for i := 0; i < 2; i++ {
			testData := testrand.Bytes(segmentSize / 8)
			path := "/some/inline/path/" + strconv.Itoa(i)
			err := ul.Upload(ctx, satellite, "bucket", path, testData)
			require.NoError(t, err)
		}

		// create 2 observers
		obs1 := newTestObserver(nil)
		obs2 := newTestObserver(nil)

		var group errgroup.Group
		group.Go(func() error {
			return metaLoop.Join(ctx, obs1)
		})
		group.Go(func() error {
			return metaLoop.Join(ctx, obs2)
		})

		err := group.Wait()
		require.NoError(t, err)

		projectID := ul.Projects[0].ID
		for _, obs := range []*testObserver{obs1, obs2} {
			assert.EqualValues(t, 7, obs.objectCount)
			assert.EqualValues(t, 5, obs.remoteSegCount)
			assert.EqualValues(t, 2, obs.inlineSegCount)
			assert.EqualValues(t, 7, len(obs.uniquePaths))
			for _, path := range obs.uniquePaths {
				assert.EqualValues(t, path.BucketName, "bucket")
				assert.EqualValues(t, path.ProjectID, projectID)
			}
			// TODO we need better calulation
			assert.NotZero(t, obs.totalMetadataSize)
		}
	})
}

func TestLoop_AllData(t *testing.T) {
	segmentSize := 8 * memory.KiB
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 4,
		UplinkCount:      3,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Metainfo.Loop.CoalesceDuration = 1 * time.Second
				config.Metainfo.Loop.ListLimit = 2
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		bucketNames := strings.Split("abc", "")

		data := testrand.Bytes(segmentSize)
		for _, up := range planet.Uplinks {
			for _, bucketName := range bucketNames {
				err := up.Upload(ctx, planet.Satellites[0], "zzz"+bucketName, "1", data)
				require.NoError(t, err)
			}
		}

		metaLoop := planet.Satellites[0].Metainfo.Loop

		obs := newTestObserver(nil)
		err := metaLoop.Join(ctx, obs)
		require.NoError(t, err)

		gotItems := len(obs.uniquePaths)
		require.Equal(t, len(bucketNames)*len(planet.Uplinks), gotItems)
	})
}

func TestLoop_ObjectNoSegments(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 4,
		UplinkCount:      1,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Metainfo.Loop.CoalesceDuration = 1 * time.Second
				config.Metainfo.Loop.ListLimit = 2
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		err := planet.Uplinks[0].CreateBucket(ctx, planet.Satellites[0], "abcd")
		require.NoError(t, err)

		project, err := planet.Uplinks[0].OpenProject(ctx, planet.Satellites[0])
		require.NoError(t, err)
		defer ctx.Check(project.Close)

		expectedNumberOfObjects := 5
		for i := 0; i < expectedNumberOfObjects; i++ {
			info, err := multipart.NewMultipartUpload(ctx, project, "abcd", "t"+strconv.Itoa(i), nil)
			require.NoError(t, err)

			_, err = multipart.CompleteMultipartUpload(ctx, project, "abcd", "t"+strconv.Itoa(i), info.StreamID, nil)
			require.NoError(t, err)
		}

		metaLoop := planet.Satellites[0].Metainfo.Loop

		obs := newTestObserver(nil)
		err = metaLoop.Join(ctx, obs)
		require.NoError(t, err)

		require.Equal(t, expectedNumberOfObjects, obs.objectCount)
		require.Zero(t, obs.inlineSegCount)
		require.Zero(t, obs.remoteSegCount)

		// add object with single segment
		data := testrand.Bytes(8 * memory.KiB)
		err = planet.Uplinks[0].Upload(ctx, planet.Satellites[0], "dcba", "1", data)
		require.NoError(t, err)

		obs = newTestObserver(nil)
		err = metaLoop.Join(ctx, obs)
		require.NoError(t, err)

		require.Equal(t, expectedNumberOfObjects+1, obs.objectCount)
		require.Zero(t, obs.inlineSegCount)
		require.Equal(t, 1, obs.remoteSegCount)
	})
}

// TestLoopObserverCancel does the following:
// * upload 3 remote segments
// * hook three observers up to metainfo loop
// * let observer 1 run normally
// * let observer 2 return an error from one of its handlers
// * let observer 3's context be canceled
// * expect observer 1 to see all segments
// * expect observers 2 and 3 to finish with errors.
func TestLoopObserverCancel(t *testing.T) {
	segmentSize := 8 * memory.KiB

	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 4,
		UplinkCount:      1,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Metainfo.Loop.CoalesceDuration = 1 * time.Second
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		ul := planet.Uplinks[0]
		satellite := planet.Satellites[0]
		metaLoop := satellite.Metainfo.Loop

		// upload 3 remote files with 1 segment
		for i := 0; i < 3; i++ {
			testData := testrand.Bytes(segmentSize)
			path := "/some/remote/path/" + strconv.Itoa(i)
			err := ul.Upload(ctx, satellite, "bucket", path, testData)
			require.NoError(t, err)
		}

		// create 1 "good" observer
		obs1 := newTestObserver(nil)
		obs1x := newTestObserver(nil)

		// create observer that will return an error from RemoteSegment
		obs2 := newTestObserver(func(ctx context.Context) error {
			return errors.New("test error")
		})

		// create observer that will cancel its own context from RemoteSegment
		obs3Ctx, cancel := context.WithCancel(ctx)
		var once int64
		obs3 := newTestObserver(func(ctx context.Context) error {
			if atomic.AddInt64(&once, 1) == 1 {
				cancel()
				<-obs3Ctx.Done() // ensure we wait for cancellation to propagate
			} else {
				panic("multiple calls to observer after loop cancel")
			}
			return nil
		})

		var group errgroup.Group
		group.Go(func() error {
			return metaLoop.Join(ctx, obs1, obs1x)
		})
		group.Go(func() error {
			err := metaLoop.Join(ctx, obs2)
			if err == nil {
				return errors.New("got no error")
			}
			if !strings.Contains(err.Error(), "test error") {
				return errors.New("expected to find error")
			}
			return nil
		})
		group.Go(func() error {
			err := metaLoop.Join(obs3Ctx, obs3)
			if !errs2.IsCanceled(err) {
				return errors.New("expected canceled")
			}
			return nil
		})

		err := group.Wait()
		require.NoError(t, err)

		// expect that obs1 saw all three segments, but obs2 and obs3 only saw the first one
		assert.EqualValues(t, 3, obs1.remoteSegCount)
		assert.EqualValues(t, 3, obs1x.remoteSegCount)
		assert.EqualValues(t, 1, obs2.remoteSegCount)
		assert.EqualValues(t, 1, obs3.remoteSegCount)
	})
}

// TestLoopCancel does the following:
// * upload 3 remote segments
// * hook two observers up to metainfo loop
// * cancel loop context partway through
// * expect both observers to exit with an error and see fewer than 3 remote segments
// * expect that a new observer attempting to join at this point receives a loop closed error.
func TestLoopCancel(t *testing.T) {
	segmentSize := 8 * memory.KiB

	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 4,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		ul := planet.Uplinks[0]
		satellite := planet.Satellites[0]

		// upload 3 remote files with 1 segment
		for i := 0; i < 3; i++ {
			testData := testrand.Bytes(segmentSize)
			path := "/some/remote/path/" + strconv.Itoa(i)
			err := ul.Upload(ctx, satellite, "bucket", path, testData)
			require.NoError(t, err)
		}

		// create a new metainfo loop
		metaLoop := metaloop.New(metaloop.Config{
			CoalesceDuration: 1 * time.Second,
			ListLimit:        10000,
		}, satellite.Metainfo.Metabase)

		// create a cancelable context to pass into metaLoop.Run
		loopCtx, cancel := context.WithCancel(ctx)

		// create 1 normal observer
		obs1 := newTestObserver(nil)

		var once int64
		// create another normal observer that will wait before returning during RemoteSegment so we can sync with context cancelation
		obs2 := newTestObserver(func(ctx context.Context) error {
			// cancel context during call to obs2.RemoteSegment inside loop
			if atomic.AddInt64(&once, 1) == 1 {
				cancel()
				<-ctx.Done() // ensure we wait for cancellation to propagate
			} else {
				panic("multiple calls to observer after loop cancel")
			}
			return nil
		})

		var group errgroup.Group

		// start loop with cancelable context
		group.Go(func() error {
			err := metaLoop.Run(loopCtx)
			if !errs2.IsCanceled(err) {
				return errors.New("expected context canceled")
			}
			return nil
		})
		group.Go(func() error {
			err := metaLoop.Join(ctx, obs1)
			if !errs2.IsCanceled(err) {
				return errors.New("expected context canceled")
			}
			return nil
		})
		group.Go(func() error {
			err := metaLoop.Join(ctx, obs2)
			if !errs2.IsCanceled(err) {
				return errors.New("expected context canceled")
			}
			return nil
		})

		err := group.Wait()
		require.NoError(t, err)

		err = metaLoop.Close()
		require.NoError(t, err)

		obs3 := newTestObserver(nil)
		err = metaLoop.Join(ctx, obs3)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loop closed")

		// expect that obs1 and obs2 each saw fewer than three remote segments
		assert.True(t, obs1.remoteSegCount < 3)
		assert.True(t, obs2.remoteSegCount < 3)
	})
}

type testObserver struct {
	objectCount       int
	remoteSegCount    int
	inlineSegCount    int
	totalMetadataSize int
	uniquePaths       map[string]metabase.SegmentLocation
	onSegment         func(context.Context) error // if set, run this during RemoteSegment()
}

func newTestObserver(onSegment func(context.Context) error) *testObserver {
	return &testObserver{
		objectCount:       0,
		remoteSegCount:    0,
		inlineSegCount:    0,
		totalMetadataSize: 0,
		uniquePaths:       make(map[string]metabase.SegmentLocation),
		onSegment:         onSegment,
	}
}

func (obs *testObserver) RemoteSegment(ctx context.Context, segment *metaloop.Segment) error {
	obs.remoteSegCount++

	key := segment.Location.Encode()
	if _, ok := obs.uniquePaths[string(key)]; ok {
		// TODO: collect the errors and check in test
		panic("Expected unique path in observer.RemoteSegment")
	}
	obs.uniquePaths[string(key)] = segment.Location

	if obs.onSegment != nil {
		return obs.onSegment(ctx)
	}

	return nil
}

func (obs *testObserver) Object(ctx context.Context, object *metaloop.Object) error {
	obs.objectCount++
	obs.totalMetadataSize += object.EncryptedMetadataSize
	return nil
}

func (obs *testObserver) InlineSegment(ctx context.Context, segment *metaloop.Segment) error {
	obs.inlineSegCount++
	key := segment.Location.Encode()
	if _, ok := obs.uniquePaths[string(key)]; ok {
		// TODO: collect the errors and check in test
		panic("Expected unique path in observer.InlineSegment")
	}
	obs.uniquePaths[string(key)] = segment.Location
	return nil
}
