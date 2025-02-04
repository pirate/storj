// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package metabase_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
	"storj.io/storj/satellite/metainfo/metabase"
)

func TestNodeAliases(t *testing.T) {
	All(t, func(ctx *testcontext.Context, t *testing.T, db *metabase.DB) {
		t.Run("Zero", func(t *testing.T) {
			defer DeleteAll{}.Check(ctx, t, db)

			nodes := []storj.NodeID{
				testrand.NodeID(),
				{},
			}
			EnsureNodeAliases{
				Opts: metabase.EnsureNodeAliases{
					Nodes: nodes,
				},
				ErrClass: &metabase.Error,
				ErrText:  "tried to add alias to zero node",
			}.Check(ctx, t, db)
		})

		t.Run("Empty", func(t *testing.T) {
			defer DeleteAll{}.Check(ctx, t, db)

			aliasesAfter := ListNodeAliases{}.Check(ctx, t, db)
			require.Len(t, aliasesAfter, 0)
		})

		t.Run("Valid", func(t *testing.T) {
			defer DeleteAll{}.Check(ctx, t, db)

			nodes := []storj.NodeID{
				testrand.NodeID(),
				testrand.NodeID(),
				testrand.NodeID(),
			}
			nodes = append(nodes, nodes...) // add duplicates to our slice

			EnsureNodeAliases{
				Opts: metabase.EnsureNodeAliases{
					Nodes: nodes,
				},
			}.Check(ctx, t, db)

			EnsureNodeAliases{
				Opts: metabase.EnsureNodeAliases{
					Nodes: nodes,
				},
			}.Check(ctx, t, db)

			aliases := ListNodeAliases{}.Check(ctx, t, db)
			require.Len(t, aliases, 3)

			for _, entry := range aliases {
				require.True(t, nodesContains(nodes, entry.ID))
				require.LessOrEqual(t, int(entry.Alias), 3)
			}

			EnsureNodeAliases{
				Opts: metabase.EnsureNodeAliases{
					Nodes: []storj.NodeID{testrand.NodeID()},
				},
			}.Check(ctx, t, db)

			aliasesAfter := ListNodeAliases{}.Check(ctx, t, db)
			require.Len(t, aliasesAfter, 4)
		})

		t.Run("Concurrent", func(t *testing.T) {
			defer DeleteAll{}.Check(ctx, t, db)

			nodes := make([]storj.NodeID, 128)
			for i := range nodes {
				nodes[i] = testrand.NodeID()
			}

			var group errgroup.Group
			for k := range nodes {
				node := nodes[k]
				group.Go(func() error {
					return db.EnsureNodeAliases(ctx, metabase.EnsureNodeAliases{
						Nodes: []storj.NodeID{node},
					})
				})
			}
			require.NoError(t, group.Wait())

			aliases := ListNodeAliases{}.Check(ctx, t, db)
			seen := map[metabase.NodeAlias]bool{}
			require.Len(t, aliases, len(nodes))
			for _, entry := range aliases {
				require.True(t, nodesContains(nodes, entry.ID))
				require.LessOrEqual(t, int(entry.Alias), len(nodes))

				require.False(t, seen[entry.Alias])
				seen[entry.Alias] = true
			}
		})

		t.Run("Stress Concurrent", func(t *testing.T) {
			defer DeleteAll{}.Check(ctx, t, db)

			nodes := make([]storj.NodeID, 128)
			for i := range nodes {
				nodes[i] = testrand.NodeID()
			}
			group, gctx := errgroup.WithContext(ctx)
			for k := 0; k < 16; k++ {
				group.Go(func() error {
					loc := nodes
					for len(loc) > 0 {
						k := testrand.Intn(10)
						if k > len(loc) {
							k = len(loc)
						}
						var batch []storj.NodeID
						batch, loc = loc[:k], loc[k:]
						err := db.EnsureNodeAliases(gctx,
							metabase.EnsureNodeAliases{Nodes: batch},
						)
						if err != nil {
							panic(err)
						}

						if gctx.Err() != nil {
							break
						}
					}
					return nil
				})
			}
			require.NoError(t, group.Wait())
		})
	})
}

func nodesContains(nodes []storj.NodeID, v storj.NodeID) bool {
	for _, n := range nodes {
		if n == v {
			return true
		}
	}
	return false
}
