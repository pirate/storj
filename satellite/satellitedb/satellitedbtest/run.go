// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedbtest

// This package should be referenced only in test files!

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"storj.io/common/testcontext"
	"storj.io/storj/private/dbutil"
	"storj.io/storj/private/dbutil/pgtest"
	"storj.io/storj/private/dbutil/pgutil"
	"storj.io/storj/private/dbutil/tempdb"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/metainfo"
	"storj.io/storj/satellite/satellitedb"
)

// SatelliteDatabases maybe name can be better.
type SatelliteDatabases struct {
	Name       string
	MasterDB   Database
	PointerDB  Database
	MetabaseDB Database
}

// Database describes a test database.
type Database struct {
	Name    string
	URL     string
	Message string
}

type ignoreSkip struct{}

func (ignoreSkip) Skip(...interface{}) {}

// Databases returns default databases.
func Databases() []SatelliteDatabases {
	cockroachConnStr := pgtest.PickCockroach(ignoreSkip{})
	postgresConnStr := pgtest.PickPostgres(ignoreSkip{})
	return []SatelliteDatabases{
		{
			Name:       "Postgres",
			MasterDB:   Database{"Postgres", postgresConnStr, "Postgres flag missing, example: -postgres-test-db=" + pgtest.DefaultPostgres + " or use STORJ_TEST_POSTGRES environment variable."},
			PointerDB:  Database{"Postgres", postgresConnStr, ""},
			MetabaseDB: Database{"Postgres", postgresConnStr, ""},
		},
		{
			Name:       "Cockroach",
			MasterDB:   Database{"Cockroach", cockroachConnStr, "Cockroach flag missing, example: -cockroach-test-db=" + pgtest.DefaultCockroach + " or use STORJ_TEST_COCKROACH environment variable."},
			PointerDB:  Database{"Cockroach", cockroachConnStr, ""},
			MetabaseDB: Database{"Cockroach", cockroachConnStr, ""},
		},
	}
}

// SchemaSuffix returns a suffix for schemas.
func SchemaSuffix() string {
	return pgutil.CreateRandomTestingSchemaName(6)
}

// SchemaName returns a properly formatted schema string.
func SchemaName(testname, category string, index int, schemaSuffix string) string {
	// postgres has a maximum schema length of 64
	// we need additional 6 bytes for the random suffix
	//    and 4 bytes for the satellite index "/S0/""

	indexStr := strconv.Itoa(index)

	var maxTestNameLen = 64 - len(category) - len(indexStr) - len(schemaSuffix) - 2
	if len(testname) > maxTestNameLen {
		testname = testname[:maxTestNameLen]
	}

	if schemaSuffix == "" {
		return strings.ToLower(testname + "/" + category + indexStr)
	}

	return strings.ToLower(testname + "/" + schemaSuffix + "/" + category + indexStr)
}

// tempMasterDB is a satellite.DB-implementing type that cleans up after itself when closed.
type tempMasterDB struct {
	satellite.DB
	tempDB *dbutil.TempDatabase
}

// Close closes a tempMasterDB and cleans it up afterward.
func (db *tempMasterDB) Close() error {
	return errs.Combine(db.DB.Close(), db.tempDB.Close())
}

// CreateMasterDB creates a new satellite database for testing.
func CreateMasterDB(ctx context.Context, log *zap.Logger, name string, category string, index int, dbInfo Database) (db satellite.DB, err error) {
	if dbInfo.URL == "" {
		return nil, fmt.Errorf("Database %s connection string not provided. %s", dbInfo.Name, dbInfo.Message)
	}

	schemaSuffix := SchemaSuffix()
	log.Debug("creating", zap.String("suffix", schemaSuffix))
	schema := SchemaName(name, category, index, schemaSuffix)

	tempDB, err := tempdb.OpenUnique(ctx, dbInfo.URL, schema)
	if err != nil {
		return nil, err
	}

	return CreateMasterDBOnTopOf(ctx, log, tempDB)
}

// CreateMasterDBOnTopOf creates a new satellite database on top of an already existing
// temporary database.
func CreateMasterDBOnTopOf(ctx context.Context, log *zap.Logger, tempDB *dbutil.TempDatabase) (db satellite.DB, err error) {
	masterDB, err := satellitedb.Open(ctx, log.Named("db"), tempDB.ConnStr, satellitedb.Options{ApplicationName: "satellite-satellitdb-test"})
	return &tempMasterDB{DB: masterDB, tempDB: tempDB}, err
}

// tempPointerDB is a satellite.DB-implementing type that cleans up after itself when closed.
type tempPointerDB struct {
	metainfo.PointerDB
	tempDB *dbutil.TempDatabase
}

// Close closes a tempPointerDB and cleans it up afterward.
func (db *tempPointerDB) Close() error {
	return errs.Combine(db.PointerDB.Close(), db.tempDB.Close())
}

// CreatePointerDB creates a new satellite pointer database for testing.
func CreatePointerDB(ctx context.Context, log *zap.Logger, name string, category string, index int, dbInfo Database) (db metainfo.PointerDB, err error) {
	if dbInfo.URL == "" {
		return nil, fmt.Errorf("Database %s connection string not provided. %s", dbInfo.Name, dbInfo.Message)
	}

	schemaSuffix := SchemaSuffix()
	log.Debug("creating", zap.String("suffix", schemaSuffix))

	schema := SchemaName(name, category, index, schemaSuffix)

	tempDB, err := tempdb.OpenUnique(ctx, dbInfo.URL, schema)
	if err != nil {
		return nil, err
	}

	return CreatePointerDBOnTopOf(ctx, log, tempDB)
}

// CreatePointerDBOnTopOf creates a new satellite database on top of an already existing
// temporary database.
func CreatePointerDBOnTopOf(ctx context.Context, log *zap.Logger, tempDB *dbutil.TempDatabase) (db metainfo.PointerDB, err error) {
	pointerDB, err := metainfo.OpenStore(ctx, log.Named("pointerdb"), tempDB.ConnStr, "satellite-satellitdb-test")
	if err != nil {
		return nil, err
	}
	err = pointerDB.MigrateToLatest(ctx)
	return &tempPointerDB{PointerDB: pointerDB, tempDB: tempDB}, err
}

// tempMetabaseDB is a metabase.DB-implementing type that cleans up after itself when closed.
type tempMetabaseDB struct {
	metainfo.MetabaseDB
	tempDB *dbutil.TempDatabase
}

// Close closes a tempPointerDB and cleans it up afterward.
func (db *tempMetabaseDB) Close() error {
	return errs.Combine(db.MetabaseDB.Close(), db.tempDB.Close())
}

// CreateMetabaseDB creates a new satellite metabase for testing.
func CreateMetabaseDB(ctx context.Context, log *zap.Logger, name string, category string, index int, dbInfo Database) (db metainfo.MetabaseDB, err error) {
	if dbInfo.URL == "" {
		return nil, fmt.Errorf("Database %s connection string not provided. %s", dbInfo.Name, dbInfo.Message)
	}

	schemaSuffix := SchemaSuffix()
	log.Debug("creating", zap.String("suffix", schemaSuffix))

	schema := SchemaName(name, category, index, schemaSuffix)

	tempDB, err := tempdb.OpenUnique(ctx, dbInfo.URL, schema)
	if err != nil {
		return nil, err
	}

	return CreateMetabaseDBOnTopOf(ctx, log, tempDB)
}

// CreateMetabaseDBOnTopOf creates a new metabase on top of an already existing
// temporary database.
func CreateMetabaseDBOnTopOf(ctx context.Context, log *zap.Logger, tempDB *dbutil.TempDatabase) (db metainfo.MetabaseDB, err error) {
	metabaseDB, err := metainfo.OpenMetabase(ctx, log.Named("metabase"), tempDB.ConnStr)
	if err != nil {
		return nil, err
	}
	return &tempMetabaseDB{MetabaseDB: metabaseDB, tempDB: tempDB}, err
}

// Run method will iterate over all supported databases. Will establish
// connection and will create tables for each DB.
func Run(t *testing.T, test func(ctx *testcontext.Context, t *testing.T, db satellite.DB)) {
	for _, dbInfo := range Databases() {
		dbInfo := dbInfo
		if strings.EqualFold(dbInfo.MasterDB.URL, "omit") {
			continue
		}

		t.Run(dbInfo.Name, func(t *testing.T) {
			t.Parallel()

			ctx := testcontext.New(t)
			defer ctx.Cleanup()

			if dbInfo.MasterDB.URL == "" {
				t.Skipf("Database %s connection string not provided. %s", dbInfo.MasterDB.Name, dbInfo.MasterDB.Message)
			}

			db, err := CreateMasterDB(ctx, zaptest.NewLogger(t), t.Name(), "T", 0, dbInfo.MasterDB)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err := db.Close()
				if err != nil {
					t.Fatal(err)
				}
			}()

			err = db.TestingMigrateToLatest(ctx)
			if err != nil {
				t.Fatal(err)
			}

			test(ctx, t, db)
		})
	}
}

// Bench method will iterate over all supported databases. Will establish
// connection and will create tables for each DB.
func Bench(b *testing.B, bench func(b *testing.B, db satellite.DB)) {
	for _, dbInfo := range Databases() {
		dbInfo := dbInfo
		if strings.EqualFold(dbInfo.MasterDB.URL, "omit") {
			continue
		}

		b.Run(dbInfo.Name, func(b *testing.B) {
			if dbInfo.MasterDB.URL == "" {
				b.Skipf("Database %s connection string not provided. %s", dbInfo.MasterDB.Name, dbInfo.MasterDB.Message)
			}

			ctx := testcontext.New(b)
			defer ctx.Cleanup()

			db, err := CreateMasterDB(ctx, zap.NewNop(), b.Name(), "X", 0, dbInfo.MasterDB)
			if err != nil {
				b.Fatal(err)
			}
			defer func() {
				err := db.Close()
				if err != nil {
					b.Fatal(err)
				}
			}()

			err = db.MigrateToLatest(ctx)
			if err != nil {
				b.Fatal(err)
			}

			// TODO: pass the ctx down
			bench(b, db)
		})
	}
}
