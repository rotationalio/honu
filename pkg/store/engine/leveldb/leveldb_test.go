package leveldb_test

import (
	crand "crypto/rand"
	"math/rand/v2"
	"testing"

	"github.com/rotationalio/honu/pkg/config"
	"github.com/rotationalio/honu/pkg/store/engine"
	. "github.com/rotationalio/honu/pkg/store/engine/leveldb"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/lamport"
	"github.com/rotationalio/honu/pkg/store/metadata"
	"github.com/rotationalio/honu/pkg/store/object"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Basic Tests
//===========================================================================

func TestEngine(t *testing.T) {
	dir := t.TempDir()
	engine, err := Open(config.StoreConfig{DataPath: dir})
	require.NoError(t, err)
	require.NotNil(t, engine)
	require.Equal(t, "leveldb", engine.Engine())
	require.NotNil(t, engine.DB())
	require.NoError(t, engine.Close())
}

func TestStore(t *testing.T) {
	db := OpenLevelDB(t, false)
	cid, oid := ulid.Make(), ulid.Make()
	k, obj := RandomObject(t, cid, oid, nil)

	// The database should be empty
	exists, err := db.Has(k)
	require.NoError(t, err)
	require.False(t, exists)

	// Get should return a not found error
	_, err = db.Get(k)
	require.ErrorIs(t, err, engine.ErrNotFound)

	// Should be able to Put a value
	err = db.Put(k, obj)
	require.NoError(t, err)

	// The database should now have the value
	exists, err = db.Has(k)
	require.NoError(t, err)
	require.True(t, exists)

	// Get should return the value
	got, err := db.Get(k)
	require.NoError(t, err)
	require.Equal(t, obj, got)

	// Should be able to Delete the value
	err = db.Delete(k)
	require.NoError(t, err)

	// The database should be empty again
	exists, err = db.Has(k)
	require.NoError(t, err)
	require.False(t, exists)

	// Get should return a not found error
	_, err = db.Get(k)
	require.ErrorIs(t, err, engine.ErrNotFound)
}

func TestReadOnly(t *testing.T) {
	// Create a LevelDB and populate it with some data
	dir := t.TempDir()
	db, err := Open(config.StoreConfig{DataPath: dir})
	require.NoError(t, err)

	cid, oid := ulid.Make(), ulid.Make()
	k, obj := RandomObject(t, cid, oid, nil)
	Populate(t, db, 10, cid, oid)
	require.NoError(t, db.Put(k, obj))
	require.NoError(t, db.Close())

	// Open the database in read-only mode
	db, err = Open(config.StoreConfig{DataPath: dir, ReadOnly: true})
	require.NoError(t, err)

	// Should be able to check if the key exists
	exists, err := db.Has(k)
	require.NoError(t, err)
	require.True(t, exists)

	// Get should return the value
	got, err := db.Get(k)
	require.NoError(t, err)
	require.Equal(t, obj, got)

	// Should not be able to Put a value
	err = db.Put(k, obj)
	require.ErrorIs(t, err, engine.ErrReadOnlyDB)

	// Should not be able to Delete a value
	err = db.Delete(k)
	require.ErrorIs(t, err, engine.ErrReadOnlyDB)
}

func TestClosed(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(config.StoreConfig{DataPath: dir})
	require.NoError(t, err)
	require.NoError(t, db.Close())

	cid, oid := ulid.Make(), ulid.Make()
	vers := RandomVersion(nil)
	k := key.New(cid, oid, vers)

	_, err = db.Has(k)
	require.ErrorIs(t, err, engine.ErrClosed)

	_, err = db.Get(k)
	require.ErrorIs(t, err, engine.ErrClosed)

	err = db.Put(k, nil)
	require.ErrorIs(t, err, engine.ErrClosed)

	err = db.Delete(k)
	require.ErrorIs(t, err, engine.ErrClosed)
}

//===========================================================================
// Test Helpers
//===========================================================================

// Open a new LevelDB engine for testing and make sure that the database is closed at
// the end of the tests and that the data is cleaned up from the temp directory.
func OpenLevelDB(t *testing.T, readonly bool) *Engine {
	dir := t.TempDir()
	engine, err := Open(config.StoreConfig{DataPath: dir, ReadOnly: readonly})
	require.NoError(t, err, "could not open leveldb at %s", dir)

	t.Cleanup(func() {
		require.NoError(t, engine.Close(), "could not close leveldb")
	})
	return engine
}

// Populate n random records with the specified collection ID and object ID.
func Populate(t *testing.T, engine *Engine, n int, cid, oid ulid.ULID) {
	var vers *lamport.Scalar
	for i := 0; i < n; i++ {
		k, obj := RandomObject(t, cid, oid, vers)
		require.NoError(t, engine.Put(k, obj))

		next := k.Version()
		require.True(t, vers.Before(&next))
		vers = &next
	}
}

func RandomObject(t *testing.T, cid, oid ulid.ULID, prev *lamport.Scalar) (key.Key, object.Object) {
	vers := RandomVersion(prev)
	meta := RandomMetadata(cid, oid, vers)

	data := make([]byte, rand.IntN(1024)+64)
	_, err := crand.Read(data)
	require.NoError(t, err)

	obj, err := object.Marshal(meta, data)
	require.NoError(t, err)

	return meta.Key(), obj
}

func RandomMetadata(cid, oid ulid.ULID, prev *lamport.Scalar) *metadata.Metadata {
	meta := &metadata.Metadata{
		ObjectID:     oid,
		CollectionID: cid,
		Version: &metadata.Version{
			Scalar: *RandomVersion(prev),
			Parent: prev,
		},
		MIME: "application/octet-stream",
	}

	return meta
}

func RandomVersion(prev *lamport.Scalar) *lamport.Scalar {
	s := &lamport.Scalar{
		PID: uint32(rand.Int32N(24)),
		VID: 1,
	}

	if prev != nil {
		s.VID = uint64(rand.Int64N(5)) + prev.VID

		if !prev.Before(s) {
			s.PID = prev.PID + 1
			if !prev.Before(s) {
				panic("failed to generate next version")
			}
		}
	}

	return s
}
