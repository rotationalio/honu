package iterator_test

import (
	"crypto/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
	"go.rtnl.ai/honu/pkg/store/iterator"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
	"go.rtnl.ai/ulid"
)

var (
	testBucket = []byte("test-bucket")
	expected   []key.Key
)

func TestCursor(t *testing.T) {
	// Generate a temporary BoltDB database and bucket, populated with test data.
	db := setupDatabase(t)

	t.Run("List", func(t *testing.T) {
		tx, err := db.Begin(false)
		require.NoError(t, err, "failed to begin read-only transaction")
		t.Cleanup(func() { tx.Rollback() })

		bkt := tx.Bucket(testBucket)
		require.NotNil(t, bkt, "test bucket does not exist")

		iter := iterator.New(bkt.Cursor())

		idx := -1
		for iter.Next() {
			idx++
			key := iter.Key()
			require.NotNil(t, key, "expected non-nil key from iterator")
			require.Equal(t, expected[idx], key, "unexpected key from iterator at index %d", idx)

			obj := iter.Object()
			require.NotNil(t, obj, "expected non-nil object from iterator")

			key, err := obj.Key()
			require.NoError(t, err, "failed to get object key")
			require.Equal(t, expected[idx], key, "unexpected object key")
		}
		require.NoError(t, iter.Error(), "unexpected iterator error")
		require.Equal(t, len(expected), idx+1, "unexpected number of keys from iterator")

		iter.Release()

		require.False(t, iter.Next(), "expected Next() to return false after Release()")
	})

	t.Run("Reverse", func(t *testing.T) {
		tx, err := db.Begin(false)
		require.NoError(t, err, "failed to begin read-only transaction")
		t.Cleanup(func() { tx.Rollback() })

		bkt := tx.Bucket(testBucket)
		require.NotNil(t, bkt, "test bucket does not exist")

		iter := iterator.New(bkt.Cursor())
		idx := len(expected)

		for iter.Prev() {
			idx--
			key := iter.Key()
			require.NotNil(t, key, "expected non-nil key from iterator")
			require.Equal(t, expected[idx], key, "unexpected key from iterator at index %d", idx)
		}

		require.NoError(t, iter.Error(), "unexpected iterator error")
		require.Equal(t, 0, idx, "unexpected number of keys from iterator")

		iter.Release()
		require.False(t, iter.Prev(), "expected Prev() to return false after Release()")
	})

	t.Run("Seek", func(t *testing.T) {
		tx, err := db.Begin(false)
		require.NoError(t, err, "failed to begin read-only transaction")
		t.Cleanup(func() { tx.Rollback() })

		bkt := tx.Bucket(testBucket)
		require.NotNil(t, bkt, "test bucket does not exist")

		iter := iterator.New(bkt.Cursor())

		// Seek to the 17th key.
		require.True(t, iter.Seek(expected[16]), "expected Seek() to find existing key")
		require.Equal(t, expected[16], iter.Key(), "unexpected key from iterator after Seek()")

		require.True(t, iter.Next(), "expected Next() to return true after Seek()")
		require.Equal(t, expected[17], iter.Key(), "unexpected key from iterator after Next()")

		// Seek to the 42nd key.
		require.True(t, iter.Seek(expected[41]), "expected Seek() to find existing key")
		require.Equal(t, expected[41], iter.Key(), "unexpected key from iterator after Seek()")

		require.True(t, iter.Prev(), "expected Prev() to return true after Seek()")
		require.Equal(t, expected[40], iter.Key(), "unexpected key from iterator after Prev()")

		// Seek to a non-existent key (between 64 and 65).
		nonExistent := make(key.Key, len(expected[0]))
		copy(nonExistent, expected[63])
		nonExistent[len(nonExistent)-1]++
		require.True(t, iter.Seek(nonExistent), "expected Seek() to find next key after non-existent key")
		require.Equal(t, expected[64], iter.Key(), "unexpected key from iterator after Seek() to non-existent key")

		// Seek to a key beyond the end of the collection.
		beyondEnd := make(key.Key, len(expected[0]))
		copy(beyondEnd, expected[len(expected)-1])
		beyondEnd[len(beyondEnd)-1]++
		require.False(t, iter.Seek(beyondEnd), "expected Seek() to return false for key beyond end of collection")
		require.Nil(t, iter.Key(), "expected nil key from iterator after Seek() beyond end of collection")

		require.NoError(t, iter.Error(), "unexpected iterator error")

		iter.Release()
		require.False(t, iter.Seek(expected[0]), "expected Seek() to return false after Release()")
	})

	t.Run("First", func(t *testing.T) {
		tx, err := db.Begin(false)
		require.NoError(t, err, "failed to begin read-only transaction")
		t.Cleanup(func() { tx.Rollback() })

		bkt := tx.Bucket(testBucket)
		require.NotNil(t, bkt, "test bucket does not exist")

		iter := iterator.New(bkt.Cursor())

		require.True(t, iter.First(), "expected First() to return true")
		require.Equal(t, expected[0], iter.Key(), "unexpected key from iterator after First()")

		require.True(t, iter.Next(), "expected Next() to return true after First()")
		require.Equal(t, expected[1], iter.Key(), "unexpected key from iterator after Next()")

		require.True(t, iter.Prev(), "expected Prev() to return true after Next()")
		require.Equal(t, expected[0], iter.Key(), "unexpected key from iterator after Prev()")

		require.False(t, iter.Prev(), "expected Prev() to return false at beginning of collection")
		require.Nil(t, iter.Key(), "expected nil key from iterator after Prev() at beginning of collection")

		iter.Release()
		require.False(t, iter.First(), "expected First() to return false after Release()")
	})

	t.Run("Last", func(t *testing.T) {
		tx, err := db.Begin(false)
		require.NoError(t, err, "failed to begin read-only transaction")
		t.Cleanup(func() { tx.Rollback() })

		bkt := tx.Bucket(testBucket)
		require.NotNil(t, bkt, "test bucket does not exist")

		iter := iterator.New(bkt.Cursor())

		require.True(t, iter.Last(), "expected Last() to return true")
		require.Equal(t, expected[len(expected)-1], iter.Key(), "unexpected key from iterator after Last()")

		require.True(t, iter.Prev(), "expected Prev() to return true after Last()")
		require.Equal(t, expected[len(expected)-2], iter.Key(), "unexpected key from iterator after Prev()")

		require.True(t, iter.Next(), "expected Next() to return true after Prev()")
		require.Equal(t, expected[len(expected)-1], iter.Key(), "unexpected key from iterator after Next()")

		require.False(t, iter.Next(), "expected Next() to return false at end of collection")
		require.Nil(t, iter.Key(), "expected nil key from iterator after Next() at end of collection")

		iter.Release()
		require.False(t, iter.Last(), "expected Last() to return false after Release()")
	})

}

func setupDatabase(t *testing.T) *bbolt.DB {
	// Create a bbolt database in a temporary file with 128 objects in it.
	tempDir := t.TempDir()
	db, err := bbolt.Open(filepath.Join(tempDir, "cursor_test.db"), 0644, nil)
	if err != nil {
		t.Fatalf("failed to create temporary bbolt database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := db.Update(populateBucket); err != nil {
		t.Fatalf("failed to populate temporary bbolt database: %v", err)
	}
	return db
}

func populateBucket(tx *bbolt.Tx) error {
	bkt, err := tx.CreateBucketIfNotExists(testBucket)
	if err != nil {
		return err
	}

	expected = make([]key.Key, 128)
	vers := &lamport.Scalar{VID: 1, PID: 8}
	entropy := ulid.Monotonic(rand.Reader, 0)

	meta := &metadata.Metadata{
		CollectionID: ulid.MustNew(ulid.Now(), entropy),
		Version: &metadata.Version{
			Scalar: *vers,
		},
		Schema: &metadata.SchemaVersion{
			Name:  "RandomData",
			Major: 1, Minor: 0, Patch: 0,
		},
	}

	for i := 0; i < 128; i++ {
		meta.ObjectID = ulid.MustNew(ulid.Now(), entropy)
		meta.Version.Created = time.Now()
		meta.Created = meta.Version.Created
		meta.Modified = meta.Version.Created

		data := make([]byte, 128)
		rand.Read(data)
		value, err := object.Marshal(meta, data)
		if err != nil {
			return err
		}

		k, _ := value.Key()
		expected[i] = k

		if err := bkt.Put(k, value); err != nil {
			return err
		}
	}

	return nil
}
