package honu_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/config"
	pb "github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

// Returns a constant list of namespace strings.
// TODO: Share with engines/leveldb/leveldb_test.go
func getNamespaces() []string {
	return []string{
		"",
		"basic",
		"namespace with spaces",
		"namespace::with::colons",
	}
}

func setupHonuDB(t testing.TB) (db *honu.DB, tmpDir string) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "honuldb-*")
	require.NoError(t, err)

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

	db, err = honu.Open(uri, conf)
	require.NoError(t, err)
	if err != nil && tmpDir != "" {
		fmt.Println(tmpDir)
		os.RemoveAll(tmpDir)
	}
	return db, tmpDir
}

func TestLevelDBInteractions(t *testing.T) {
	db, tmpDir := setupHonuDB(t)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	for _, namespace := range getNamespaces() {
		// Use a constant key to ensure namespaces
		// are working correctly.
		key := []byte("foo")
		//append a constant to namespace as the value
		//because when the empty namespace is returned
		//as a key it is unmarsheled as []byte(nil)
		//instead of []byte{}
		expectedValue := []byte(namespace + "this is the value of foo")

		// Put a version to the database
		obj, err := db.Put(key, expectedValue, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.False(t, obj.Tombstone())

		// Get the version of foo from the database
		value, err := db.Get(key, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.Equal(t, expectedValue, value)

		// Get the meta data from foo
		obj, err = db.Object(key, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.Equal(t, uint64(1), obj.Version.Version)
		require.False(t, obj.Tombstone())

		// Delete the version from the database and ensure you
		// are not able to get the deleted version
		_, err = db.Delete(key, options.WithNamespace(namespace))
		require.NoError(t, err)

		value, err = db.Get(key, options.WithNamespace(namespace))
		require.Error(t, err)
		require.Empty(t, value)

		// Get the tombstone from the database
		obj, err = db.Object(key, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.Equal(t, uint64(2), obj.Version.Version)
		require.True(t, obj.Tombstone())
		require.Empty(t, obj.Data)

		// Be able to "undelete" a tombstone
		undeadValue := []byte("this is the undead foo")
		obj, err = db.Put(key, undeadValue, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.False(t, obj.Tombstone())

		// Get the metadata from the database (should no longer be a tombstone)
		obj, err = db.Object(key, options.WithNamespace(namespace))
		require.NoError(t, err)
		require.Equal(t, uint64(3), obj.Version.Version)
		require.False(t, obj.Tombstone())

		// TODO: figure out what to do with this testcase.
		// Iter currently grabs the namespace by splitting
		// on :: and grabbing the first string, so it only
		// grabs "namespace".
		if namespace == "namespace::with::colons" {
			continue
		}
		// Put a range of data into the database
		pairs := [][]string{
			{"aa", "first"},
			{"ab", "second"},
			{"ba", "third"},
			{"bb", "fourth"},
			{"bc", "fifth"},
			{"ca", "sixth"},
			{"cb", "seventh"},
		}
		for _, pair := range pairs {
			key := []byte(pair[0])
			value := []byte(pair[1])
			_, err := db.Put(key, value, options.WithNamespace(namespace))
			require.NoError(t, err)
		}

		// Iterate over a prefix in the database
		iter, err := db.Iter([]byte("b"), options.WithNamespace(namespace))
		require.NoError(t, err)
		collected := 0
		for iter.Next() {
			key := iter.Key()
			require.Equal(t, string(key), pairs[collected+2][0])

			value := iter.Value()
			fmt.Println(value)
			require.Equal(t, string(value), string(pairs[collected+2][1]))

			obj, err := iter.Object()
			require.NoError(t, err)
			require.Equal(t, uint64(1), obj.Version.Version)

			collected++
		}

		require.Equal(t, 3, collected)
		require.NoError(t, iter.Error())
		iter.Release()
	}
}

// Global variables to prevent compiler optimizations
var (
	gKey   []byte
	gValue []byte
	gErr   error
	gObj   *pb.Object
)
