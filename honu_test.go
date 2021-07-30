package honu_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/config"
	"github.com/stretchr/testify/require"
)

func TestLevelDBInteractions(t *testing.T) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "honuldb-*")
	require.NoError(t, err)

	// Cleanup after the test is over
	defer os.RemoveAll(tmpDir)

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
	require.NoError(t, err)

	// Close the database when we're done with the test
	defer db.Close()

	// Put a version to the database
	err = db.Put([]byte("foo"), []byte("this is the value of foo"))
	require.NoError(t, err)

	// Get the version of foo from the database
	value, err := db.Get([]byte("foo"))
	require.NoError(t, err)
	require.Equal(t, []byte("this is the value of foo"), value)

	// Get the meta data from foo
	obj, err := db.Object([]byte("foo"))
	require.NoError(t, err)
	require.Equal(t, uint64(1), obj.Version.Version)

	// Delete the version from the database
	err = db.Delete([]byte("foo"))
	require.NoError(t, err)

	// Should not be able to get the deleted version
	value, err = db.Get([]byte("foo"))
	require.Error(t, err)
	require.Empty(t, value)

	// Get the tombstone from the database
	obj, err = db.Object([]byte("foo"))
	require.NoError(t, err)
	require.Equal(t, uint64(2), obj.Version.Version)
	require.True(t, obj.Tombstone())
	require.Empty(t, obj.Data)

	// Be able to "undelete" a tombstone
	err = db.Put([]byte("foo"), []byte("this is the undead foo"))
	require.NoError(t, err)

	value, err = db.Get([]byte("foo"))
	require.NoError(t, err)
	require.Equal(t, []byte("this is the undead foo"), value)

	// Get the tombstone from the database
	obj, err = db.Object([]byte("foo"))
	require.NoError(t, err)
	require.Equal(t, uint64(3), obj.Version.Version)
	require.False(t, obj.Tombstone())

	// Put a range of data into the database
	require.NoError(t, db.Put([]byte("aa"), []byte("123456")))
	require.NoError(t, db.Put([]byte("ab"), []byte("7890123")))
	require.NoError(t, db.Put([]byte("ba"), []byte("4567890")))
	require.NoError(t, db.Put([]byte("bb"), []byte("1234567")))
	require.NoError(t, db.Put([]byte("bc"), []byte("9012345")))
	require.NoError(t, db.Put([]byte("ca"), []byte("67890123")))
	require.NoError(t, db.Put([]byte("cb"), []byte("4567890123")))

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
