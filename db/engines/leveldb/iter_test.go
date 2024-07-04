package leveldb_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	pb "github.com/rotationalio/honu/object/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// This test verifies LevelDB functionality to ensure that our iterator functionality
// matches the expected iterator API for leveldb.
func TestLevelDBFunctionality(t *testing.T) {
	// Setup a levelDB Engine and create a bunch of keys.
	engine, _ := setupLevelDBEngine(t)
	ldb := engine.DB()

	keys := make([][]byte, 0, 100)
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		keys = append(keys, key)
		require.NoError(t, ldb.Put(key, randomData(192), nil), "could not put fixture data")
	}

	// Create an iterator to test Seek/Prev/Next functionality
	iter := ldb.NewIterator(nil, nil)

	// When next has not been called, what does Prev do?
	require.False(t, iter.Prev(), "if Prev is called before Next, we expect it to return false")

	// If we seek to the first key, the value of the key should be the first key
	// If we call Prev, we should get false, because the Seek was to the first key but
	// what is the value of the cursor, do we have to call Next again to get back to
	// the first key?
	require.True(t, iter.Seek(keys[0]), "we should be able to seek to the first key")
	require.True(t, bytes.Equal(iter.Key(), keys[0]), "the cursor is now at the first key")
	require.False(t, iter.Prev(), "if we're at the first key, prev should be false")
	require.Nil(t, iter.Key(), "call Prev moved us behind the first key - so it should now be nil")
	require.True(t, iter.Next(), "we should now be able to move back to the first key")
	require.True(t, bytes.Equal(iter.Key(), keys[0]), "the cursor is now at the first key")
}

// Test Honu seek behavior matches LevelDB  API
func TestHonuSeek(t *testing.T) {
	// Setup a levelDB Engine and create a bunch of keys.
	db, _ := setupLevelDBEngine(t)

	keys := make([][]byte, 0, 100)
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		keys = append(keys, key)
		require.NoError(t, db.Put(key, randomObject(key, 192), nil), "could not put fixture data")
	}

	// Create an iterator to test Seek/Prev/Next functionality
	iter, err := db.Iter(nil, nil)
	require.NoError(t, err, "could not create honu leveldb iterator")

	// When next has not been called, what does Prev do?
	require.False(t, iter.Prev(), "if Prev is called before Next, we expect it to return false")

	// If we seek to the first key, the value of the key should be the first key
	// If we call Prev, we should get false, because the Seek was to the first key but
	// what is the value of the cursor, do we have to call Next again to get back to
	// the first key?
	require.True(t, iter.Seek(keys[0]), "we should be able to seek to the first key")
	require.True(t, bytes.Equal(iter.Key(), keys[0]), "the cursor is now at the first key")
	require.False(t, iter.Prev(), "if we're at the first key, prev should be false")
	require.Nil(t, iter.Key(), "call Prev moved us behind the first key - so it should now be nil")
	require.True(t, iter.Next(), "we should now be able to move back to the first key")
	require.True(t, bytes.Equal(iter.Key(), keys[0]), "the cursor is now at the first key")
}

// Helper function to generate random data
func randomData(len int) []byte {
	data := make([]byte, len)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	return data
}

// Helper function to generate a random object data
func randomObject(key []byte, len int) []byte {
	obj := &pb.Object{
		Key:       key,
		Namespace: "default",
		Version: &pb.Version{
			Pid:       1,
			Version:   1,
			Region:    "testing",
			Parent:    nil,
			Tombstone: false,
		},
		Region: "testing",
		Owner:  "testing",
		Data:   randomData(len),
	}

	// Marshal the object
	data, err := proto.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return data
}
