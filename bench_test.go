package honu_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	pb "github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

// Global variables to prevent compiler optimizations
var (
	gKey   []byte
	gValue []byte
	gErr   error
	gObj   *pb.Object
)

func setupLevelDB(t testing.TB) (*leveldb.DB, string) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "leveldb-*")
	require.NoError(t, err)

	// Open a leveldb database directly without honu wrapper
	db, err := leveldb.OpenFile(tmpDir, nil)
	if err != nil && tmpDir != "" {
		os.RemoveAll(tmpDir)
	}
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(tmpDir)
	})

	return db, tmpDir
}

func BenchmarkHonuGet(b *testing.B) {
	db, _ := setupHonuDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	_, err = db.Put(key, value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gValue, gErr = db.Get(key)
	}

	require.NoError(b, gErr)
	require.Equal(b, value, gValue)
}

func BenchmarkLevelDBGet(b *testing.B) {
	db, _ := setupLevelDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	require.NoError(b, db.Put(key, value, nil))

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gValue, gErr = db.Get(key, nil)
	}

	require.NoError(b, gErr)
	require.Equal(b, value, gValue)
}

func BenchmarkHonuPut(b *testing.B) {
	db, _ := setupHonuDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, gErr = db.Put(key, value)
	}

	require.NoError(b, gErr)
}

func BenchmarkLevelDBPut(b *testing.B) {
	db, _ := setupLevelDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gErr = db.Put(key, value, nil)
	}

	require.NoError(b, gErr)
}

func BenchmarkHonuDelete(b *testing.B) {
	db, _ := setupHonuDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		_, err = db.Put(key, value)
		require.NoError(b, err)
		b.StartTimer()
		_, gErr = db.Delete(key)
	}

	require.NoError(b, gErr)
}

func BenchmarkLevelDBDelete(b *testing.B) {
	db, _ := setupLevelDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		require.NoError(b, db.Put(key, value, nil))
		b.StartTimer()
		gErr = db.Delete(key, nil)
	}

	require.NoError(b, gErr)
}

func BenchmarkHonuIter(b *testing.B) {
	db, _ := setupHonuDB(b)

	// Create a key and value
	for _, key := range []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"} {
		value := make([]byte, 4096)
		_, err := rand.Read(value)
		require.NoError(b, err)

		_, err = db.Put([]byte(key), value)
		require.NoError(b, err)
	}

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter, err := db.Iter(nil)
		require.NoError(b, err)
		for iter.Next() {
			gKey = iter.Key()
			gValue = iter.Value()
		}

		gErr = iter.Error()
		iter.Release()
	}

	require.NoError(b, gErr)
	require.Len(b, gKey, 2)
	require.Len(b, gValue, 4096)
}

func BenchmarkLevelDBIter(b *testing.B) {
	db, _ := setupLevelDB(b)

	// Create a key and value
	for _, key := range []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"} {
		value := make([]byte, 4096)
		_, err := rand.Read(value)
		require.NoError(b, err)

		require.NoError(b, db.Put([]byte(key), value, nil))
	}

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			gKey = iter.Key()
			gValue = iter.Value()
		}

		gErr = iter.Error()
		iter.Release()
	}

	require.NoError(b, gErr)
	require.Len(b, gKey, 2)
	require.Len(b, gValue, 4096)
}

func BenchmarkHonuObject(b *testing.B) {
	db, _ := setupHonuDB(b)

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err := rand.Read(value)
	require.NoError(b, err)

	_, err = db.Put(key, value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gObj, gErr = db.Object(key)
	}

	require.NoError(b, gErr)
	require.NotEmpty(b, gObj)
}
