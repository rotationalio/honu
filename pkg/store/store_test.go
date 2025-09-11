package store_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/logger"
	"go.rtnl.ai/honu/pkg/store"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Test Suite with underlying database
//===========================================================================

type honuTestSuite struct {
	suite.Suite
	conf  config.Config
	store *store.Store
}

func TestWriteableStore(t *testing.T) {
	tests := &honuTestSuite{
		conf: config.Config{
			PID:          1,
			Maintenance:  false,
			LogLevel:     logger.LevelDecoder(zerolog.ErrorLevel),
			ConsoleLog:   false,
			BindAddr:     ":11111",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
			Store: config.StoreConfig{
				ReadOnly:    false,
				DataPath:    filepath.Join(t.TempDir(), "honu-test.db"),
				Concurrency: 128,
			},
		},
	}

	var err error
	tests.store, err = store.Open(tests.conf)
	require.NoError(t, err, "failed to open store, could not start tests")

	suite.Run(t, tests)
}

func TestReadonlyStore(t *testing.T) {
	// Before implementing read-only stores, we need to create a database to read from.
	// If the database does not exist, then it cannot create it in read-only mode.
	t.Skip("skipping read-only tests until we have fixtures")

	tests := &honuTestSuite{
		conf: config.Config{
			PID:          1,
			Maintenance:  false,
			LogLevel:     logger.LevelDecoder(zerolog.ErrorLevel),
			ConsoleLog:   false,
			BindAddr:     ":11111",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
			Store: config.StoreConfig{
				ReadOnly:    true,
				DataPath:    filepath.Join(t.TempDir(), "honu-test.db"),
				Concurrency: 128,
			},
		},
	}
	fmt.Println(tests.conf.Store.DataPath)

	var err error
	tests.store, err = store.Open(tests.conf)
	require.NoError(t, err, "failed to open store, could not start tests")

	suite.Run(t, tests)
}

//===========================================================================
// Writeable and Read-Only Store Tests
//===========================================================================

func (s *honuTestSuite) TestInitialized() {
	// In read-only mode, the collections should not be overwritten.
	require := s.Require()

	db := s.store.DB()
	require.NotNil(db, "store db should not be nil")

	collections := []ulid.ULID{
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	for _, collection := range collections {
		exists, err := s.store.Has(collection)
		require.NoError(err, "failed to check if collection exists")
		require.True(exists, "collection %s should exist", collection)
	}
}

//===========================================================================
// Unit tests without underlying database
//===========================================================================

func TestSystemCollections(t *testing.T) {
	// System collections must be less than any ULID that would be generated.
	// They also must be sorted in a specific order for efficient access.
	collections := []ulid.ULID{
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	require.True(t, sort.SliceIsSorted(collections, func(i, j int) bool {
		return bytes.Compare(collections[i][:], collections[j][:]) < 0
	}))

	// All collection ULID timestamps must be less than 1984-04-07T09:21:07-06:00
	epoch, _ := time.Parse(time.RFC3339, "1984-04-07T09:21:07-06:00")
	for _, collection := range collections {
		ts := collection.Timestamp()
		require.True(t, ts.Before(epoch), "collection %s timestamp %s is not before epoch %s", collection, ts, epoch)
	}
}

func Example_system_names() {
	convert := func(b []byte) string {
		s := make([]byte, 0, len(b))
		for i, c := range b {
			if c < 0x20 || c > 0x7E {
				if i == 0 {
					continue
				}
				s = append(s, ' ')
			} else {
				s = append(s, c)
			}
		}

		return string(s)
	}

	systemIDs := []ulid.ULID{
		store.SystemHonuAgent,
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
		store.SystemCollectionNames,
	}

	for _, sysid := range systemIDs {
		ts := sysid.Timestamp()
		name := convert(sysid[:])

		fmt.Printf("%s (%s)\n", name, ts.UTC().Format(time.RFC3339))
	}

	// Output:
	// honu honuagent  (1984-03-19T12:08:28Z)
	// honu collection (1984-03-19T12:08:28Z)
	// honu accesslist (1984-03-19T12:08:28Z)
	// honu networking (1984-03-19T12:08:28Z)
	// honu colnameidx (1984-03-19T12:08:28Z)
}

func TestInitializedEmpty(t *testing.T) {
	// Test Setup
	conf := config.Config{
		PID: uint32(8),
		Store: config.StoreConfig{
			DataPath:    filepath.Join(t.TempDir(), "honu-test.db"),
			ReadOnly:    false,
			Concurrency: 16,
		},
	}

	version := lamport.Scalar{PID: 0, VID: 1}
	collections := []ulid.ULID{
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	// Ensure the store is intialized when the database is empty.
	bdb, err := bbolt.Open(conf.Store.DataPath, 0600, nil)
	require.NoError(t, err, "could not open bbolt for testing")

	// Helper method to check if a key exists in the database.
	hasKey := func(bdb *bbolt.DB, key []byte) (exists bool, err error) {
		err = bdb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket(store.SystemCollections[:])
			if b == nil {
				return nil
			}

			if v := b.Get(key); v != nil {
				exists = true
			}
			return nil
		})
		return exists, err
	}

	// Ensure the collections do not exist.
	for _, collection := range collections {
		cKey := key.New(collection, &version)
		exists, err := hasKey(bdb, cKey)
		require.NoError(t, err, "failed to check if collection exists")
		require.False(t, exists, "collection %s should exist", collection)
	}

	// Open the store to initialize the database.
	require.NoError(t, bdb.Close(), "could not close bbolt")

	db, err := store.Open(conf)
	require.NoError(t, err, "could not open store")

	// Ensure the collections now exist.
	for _, collection := range collections {
		cKey := key.New(collection, &version)
		exists, err := hasKey(db.DB(), cKey)
		require.NoError(t, err, "failed to check if collection exists")
		require.True(t, exists, "collection %s should exist", collection)
	}
}

//===========================================================================
// Fixtures Management
//===========================================================================
