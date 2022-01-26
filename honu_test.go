package honu_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/config"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

var pairs = [][]string{
	{"aa", "first"},
	{"ab", "second"},
	{"ba", "third"},
	{"bb", "fourth"},
	{"bc", "fifth"},
	{"ca", "sixth"},
	{"cb", "seventh"},
}

// Returns a constant list of namespace strings.
// TODO: Share with engines/leveldb/leveldb_test.go
var testNamespaces = []string{
	"",
	"basic",
	"namespace with spaces",
	"namespace::with::colons",
}

func setupHonuDB(t testing.TB) (db *honu.DB, tmpDir string) {
	// Create a new leveldb database in a temporary directory
	tmpDir, err := ioutil.TempDir("", "honuldb-*")
	require.NoError(t, err)

	// Open a Honu leveldb database with default configuration
	uri := fmt.Sprintf("leveldb:///%s", tmpDir)
	db, err = honu.Open(uri, config.WithReplica(config.ReplicaConfig{PID: 8, Region: "us-southwest-16", Name: "testing"}))
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

	for _, namespace := range testNamespaces {
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

		// Attempt to directly update the object in the database with a later version
		obj.Data = []byte("directly updated")
		obj.Owner = "me"
		obj.Version.Parent = nil
		obj.Version.Version = 42
		obj.Version.Pid = 93
		obj.Version.Region = "here"
		obj.Version.Tombstone = false
		_, err = db.Update(obj)
		require.NoError(t, err)

		// TODO: figure out what to do with this testcase.
		// Iter currently grabs the namespace by splitting
		// on :: and grabbing the first string, so it only
		// grabs "namespace".
		if namespace == "namespace::with::colons" {
			continue
		}

		// Put a range of data into the database
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

func TestUpdate(t *testing.T) {
	// Create a test database to attempt to update
	db, tmpDir := setupHonuDB(t)

	// Cleanup when we're done with the test
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	// Create a random object in the database to start update tests on
	key := randomData(32)
	namespace := "yeti"
	root, err := db.Put(key, randomData(96), options.WithNamespace(namespace))
	require.NoError(t, err, "could not put random data")

	// Generate new1 - a linear object from root as though it were from a different replica
	new1 := &object.Object{
		Key:       key,
		Namespace: namespace,
		Version: &object.Version{
			Pid:     113,
			Version: 2,
			Parent:  root.Version,
		},
		Region: "the-void",
		Owner:  root.Owner,
		Data:   randomData(112),
	}

	// Should be able to update with no namespace option
	update, err := db.Update(new1)
	require.NoError(t, err, "could not update db with new1")
	require.Equal(t, honu.UpdateLinear, update, "expected new1 update to be linear")
	requireObjectEqual(t, db, new1, key, namespace)

	// Should not be be able to update with the same version twice, since it is now no
	// longer later than previous version (it is the equal version on disk).
	update, err = db.Update(new1)
	require.EqualError(t, err, "cannot update object, it is not a later version then the current object")
	require.Equal(t, honu.UpdateNoChange, update)

	// Should be able to force the update to apply the same object back to disk.
	update, err = db.Update(new1, options.WithForce())
	require.NoError(t, err, "could not force update with new1")
	require.Equal(t, honu.UpdateForced, update)

	// Generate new2 - an object stomping new1 as though it were from a different replica
	new2 := &object.Object{
		Key:       key,
		Namespace: namespace,
		Version: &object.Version{
			Pid:     42,
			Version: 2,
			Parent:  root.Version,
		},
		Region: "the-other-void",
		Owner:  root.Owner,
		Data:   randomData(112),
	}

	// Update with the wrong namespace should error
	update, err = db.Update(new2, options.WithNamespace("this is not the right namespace for sure"))
	require.EqualError(t, err, "options namespace does not match object namespace")
	require.Equal(t, honu.UpdateNoChange, update)
	requireObjectEqual(t, db, new1, key, namespace)

	// Update with the wrong namespace but with force should not error and create a new object
	// NOTE: this is kind of a wild force since now the object has the wrong namespace metadata.
	update, err = db.Update(new2, options.WithNamespace("trashcan"), options.WithForce())
	require.NoError(t, err)
	require.Equal(t, honu.UpdateForced, update)
	requireObjectEqual(t, db, new1, key, namespace)
	requireObjectEqual(t, db, new2, key, "trashcan")

	// Update with same namespace option should not error.
	update, err = db.Update(new2, options.WithNamespace(namespace))
	require.NoError(t, err, "could not update new2")
	require.Equal(t, honu.UpdateStomp, update)
	requireObjectEqual(t, db, new2, key, namespace)

	// Generate new3 - an object skipping new2 as though it were from the same replica
	new3 := &object.Object{
		Key:       key,
		Namespace: namespace,
		Version: &object.Version{
			Pid:     42,
			Version: 12,
			Parent:  root.Version,
		},
		Region: "the-other-void",
		Owner:  root.Owner,
		Data:   randomData(112),
	}

	// Ensure UpdateSkip is returned
	update, err = db.Update(new3)
	require.NoError(t, err, "could not update new3")
	require.Equal(t, honu.UpdateSkip, update)
	requireObjectEqual(t, db, new3, key, namespace)

	// Update with an earlier version should error
	update, err = db.Update(new1)
	require.EqualError(t, err, "cannot update object, it is not a later version then the current object")
	require.Equal(t, honu.UpdateNoChange, update)
	requireObjectEqual(t, db, new3, key, namespace)

	// Should be able to force the update to apply the earlier object back to disk
	update, err = db.Update(new1, options.WithForce())
	require.NoError(t, err, "could not force update with new1")
	require.Equal(t, honu.UpdateForced, update)
	requireObjectEqual(t, db, new1, key, namespace)

	// Update an object that does not exist should not error.
	stranger := &object.Object{
		Key:       randomData(18),
		Namespace: "default",
		Version: &object.Version{
			Pid:     1,
			Version: 1,
			Parent:  nil,
		},
		Region: "the-void",
		Owner:  "me",
		Data:   randomData(8),
	}

	update, err = db.Update(stranger)
	require.NoError(t, err)
	require.Equal(t, honu.UpdateLinear, update)
}

// Helper assertion function to check to make sure an object matches what is in the database
func requireObjectEqual(t *testing.T, db *honu.DB, expected *object.Object, key []byte, namespace string) {
	actual, err := db.Object(key, options.WithNamespace(namespace))
	require.NoError(t, err, "could not fetch expected object from the database")

	// NOTE: we cannot do a require.Equal(t, expected, actual) because the test will hang
	// it's not clear if there is a recursive loop with version comparisons or some other
	// deep equality is causing the problem. Instead we directly compare the data.
	require.True(t, bytes.Equal(expected.Key, actual.Key), "key is not equal")
	require.Equal(t, expected.Namespace, actual.Namespace, "namespace not equal")
	require.Equal(t, expected.Region, actual.Region, "region not equal")
	require.Equal(t, expected.Owner, actual.Owner, "owner not equal")
	require.True(t, expected.Version.Equal(actual.Version), "versions not equal")
	require.Equal(t, expected.Version.Region, actual.Version.Region, "version region not equal")
	require.Equal(t, expected.Version.Tombstone, actual.Version.Tombstone, "version tombstone not the same")
	if expected.Version.Parent != nil {
		require.True(t, expected.Version.Parent.Equal(actual.Version.Parent), "parents not equal")
		require.Equal(t, expected.Version.Parent.Region, actual.Version.Parent.Region, "parent regions not equal")
		require.Equal(t, expected.Version.Parent.Tombstone, actual.Version.Parent.Tombstone, "parent tombstone not the same")
	} else {
		require.Nil(t, actual.Version.Parent, "expected parent is nil")
	}
	require.True(t, bytes.Equal(expected.Data, actual.Data), "value is not equal")
}

// Helper function to generate random data
func randomData(len int) []byte {
	data := make([]byte, len)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	return data
}
