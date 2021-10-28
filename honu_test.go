package honu_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/config"
	pb "github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func setupHonuDB() (*honu.DB, string, error) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "honuldb-*")
	if err != nil {
		return nil, "", err
	}

	// Open a Honu leveldb database with default configuration
	uri := fmt.Sprintf("leveldb:///%s", tmpDir)
	conf := config.ReplicaConfig{
		Enabled:        true,
		BindAddr:       ":443",
		PID:            8,
		Region:         "us-southwest-16",
		Name:           "testing",
		GossipInterval: 1 * time.Minute,
		GossipSigma:    15 * time.Second,
	}

	db, err := honu.Open(uri, conf)
	return db, tmpDir, err
}

func setupLevelDB() (*leveldb.DB, string, error) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "honuldb-*")
	if err != nil {
		return nil, "", err
	}

	// Open a leveldb database directly without honu wrapper
	db, err := leveldb.OpenFile(tmpDir, nil)
	return db, tmpDir, err
}

func TestLevelDBInteractions(t *testing.T) {
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(t, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Put a version to the database
	err = db.Put([]byte("foo"), []byte("this is the value of foo"), nil)
	require.NoError(t, err)

	// Get the version of foo from the database
	value, err := db.Get([]byte("foo"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("this is the value of foo"), value)

	// Get the meta data from foo
	obj, err := db.Object([]byte("foo"), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(1), obj.Version.Version)

	// Delete the version from the database
	err = db.Delete([]byte("foo"), nil)
	require.NoError(t, err)

	// Should not be able to get the deleted version
	value, err = db.Get([]byte("foo"), nil)
	require.Error(t, err)
	require.Empty(t, value)

	// Get the tombstone from the database
	obj, err = db.Object([]byte("foo"), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(2), obj.Version.Version)
	require.True(t, obj.Tombstone())
	require.Empty(t, obj.Data)

	// Be able to "undelete" a tombstone
	err = db.Put([]byte("foo"), []byte("this is the undead foo"), nil)
	require.NoError(t, err)

	value, err = db.Get([]byte("foo"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("this is the undead foo"), value)

	// Get the tombstone from the database
	obj, err = db.Object([]byte("foo"), nil)
	require.NoError(t, err)
	require.Equal(t, uint64(3), obj.Version.Version)
	require.False(t, obj.Tombstone())

	// Put a range of data into the database
	require.NoError(t, db.Put([]byte("aa"), []byte("123456"), nil))
	require.NoError(t, db.Put([]byte("ab"), []byte("7890123"), nil))
	require.NoError(t, db.Put([]byte("ba"), []byte("4567890"), nil))
	require.NoError(t, db.Put([]byte("bb"), []byte("1234567"), nil))
	require.NoError(t, db.Put([]byte("bc"), []byte("9012345"), nil))
	require.NoError(t, db.Put([]byte("ca"), []byte("67890123"), nil))
	require.NoError(t, db.Put([]byte("cb"), []byte("4567890123"), nil))

	// Iterate over a prefix in the database
	iter, err := db.Iter([]byte("b"))
	require.NoError(t, err)
	collected := 0
	for iter.Next() {
		collected++

		key := iter.Key()
		require.Len(t, key, 2)

		value := iter.Value()
		require.Len(t, value, 7)

		obj, err := iter.Object()
		require.NoError(t, err)
		require.Equal(t, uint64(1), obj.Version.Version)
	}

	require.Equal(t, 3, collected)
	require.NoError(t, iter.Error())
	iter.Release()
}

// Global variables to prevent compiler optimizations
var (
	gKey   []byte
	gValue []byte
	gErr   error
	gObj   *pb.Object
)

func BenchmarkHonuGet(b *testing.B) {
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
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

func BenchmarkLevelDBGet(b *testing.B) {
	db, tmpDir, err := setupLevelDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
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
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gErr = db.Put(key, value, nil)
	}

	require.NoError(b, gErr)
}

func BenchmarkLevelDBPut(b *testing.B) {
	db, tmpDir, err := setupLevelDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
	require.NoError(b, err)

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gErr = db.Put(key, value, nil)
	}

	require.NoError(b, gErr)
}

func BenchmarkHonuDelete(b *testing.B) {
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
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

func BenchmarkLevelDBDelete(b *testing.B) {
	db, tmpDir, err := setupLevelDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
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
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	for _, key := range []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"} {
		value := make([]byte, 4096)
		_, err = rand.Read(value)
		require.NoError(b, err)

		require.NoError(b, db.Put([]byte(key), value, nil))
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
	db, tmpDir, err := setupLevelDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	for _, key := range []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"} {
		value := make([]byte, 4096)
		_, err = rand.Read(value)
		require.NoError(b, err)

		require.NoError(b, db.Put([]byte(key), value, nil))
	}

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := db.NewIterator(nil, nil)
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

func BenchmarkHonuObject(b *testing.B) {
	db, tmpDir, err := setupHonuDB()
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	require.NoError(b, err)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a key and value
	key := []byte("foo")
	value := make([]byte, 4096)
	_, err = rand.Read(value)
	require.NoError(b, err)

	require.NoError(b, db.Put(key, value, nil))

	// Reset the timer to focus only on the get call
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gObj, gErr = db.Object(key, nil)
	}

	require.NoError(b, gErr)
	require.NotEmpty(b, gObj)
}
