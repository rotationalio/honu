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
	"github.com/rotationalio/honu/engines/leveldb"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

// a test set of key/value pairs used to evaluate iteration
// note because :: is the namespace separator in leveldb, we want to ensure that keys
// with colons are correctly iterated on.
var pairs = [][]string{
	{"aa", "first"},
	{"ab", "second"},
	{"b::a", "third"},
	{"b::b", "fourth"},
	{"b::c", "fifth"},
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
	tmpDir, err := ioutil.TempDir("", "honudb-*")
	require.NoError(t, err)

	// Open a Honu leveldb database with default configuration
	uri := fmt.Sprintf("leveldb:///%s", tmpDir)
	db, err = honu.Open(uri, config.WithReplica(config.ReplicaConfig{PID: 8, Region: "us-southwest-16", Name: "testing"}))
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

func TestLevelDBInteractions(t *testing.T) {
	db, _ := setupHonuDB(t)

	totalKeys := 0
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
		totalKeys++

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

		// Put a range of data into the database
		for _, pair := range pairs {
			key := []byte(pair[0])
			value := []byte(pair[1])
			_, err := db.Put(key, value, options.WithNamespace(namespace))
			require.NoError(t, err)
			totalKeys++
		}

		// Iterate over a prefix in the database
		iter, err := db.Iter([]byte("b"), options.WithNamespace(namespace))
		require.NoError(t, err)
		collected := 0
		for iter.Next() {
			key := iter.Key()
			require.Equal(t, string(key), pairs[collected+2][0])

			value := iter.Value()
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

	// Test iteration over all the namespaces
	_, ok := db.Engine().(*leveldb.LevelDBEngine)
	require.True(t, ok, "the engine type returned should be a leveldb.DB")
	requireDatabaseLen(t, db, totalKeys)
}

func TestUpdate(t *testing.T) {
	// Create a test database to attempt to update
	db, _ := setupHonuDB(t)

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

func TestTombstones(t *testing.T) {
	// Create a test database
	db, _ := setupHonuDB(t)

	// Assert that there is nothing in the namespace as an initial check
	requireNamespaceLen(t, db, "graveyard", 0)
	requireDatabaseLen(t, db, 0)

	// Create a list of keys with integer values
	keys := make([][]byte, 0, 20)
	for i := 0; i < 20; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		keys = append(keys, key)
	}

	// Add data to the database
	for _, key := range keys {
		db.Put(key, randomData(256), options.WithNamespace("graveyard"))
	}
	requireNamespaceLen(t, db, "graveyard", 20)
	requireDatabaseLen(t, db, 20)

	// Delete all even keys
	for i, key := range keys {
		if i%2 == 0 {
			db.Delete(key, options.WithNamespace("graveyard"))
		}
	}

	// Ensure that the iterator returns 10 items but that there are still 20 objects
	// including tombstones still stored in the database.
	requireNamespaceLen(t, db, "graveyard", 10)
	requireGraveyardLen(t, db, "graveyard", 20)
	requireDatabaseLen(t, db, 20)

	// Sanity check, attempt to get Get all keys and verify tombstones
	for i, key := range keys {
		if i%2 == 0 {
			// This is a tombstone
			val, err := db.Get(key, options.WithNamespace("graveyard"))
			require.EqualError(t, err, "not found", "tombstone did not return a not found error")
			require.Nil(t, val, "tombstone returned a non nil value")

			obj, err := db.Object(key, options.WithNamespace("graveyard"))
			require.NoError(t, err, "tombstone did not return an object")
			require.True(t, obj.Tombstone())
		} else {
			// Not a tombstone
			val, err := db.Get(key, options.WithNamespace("graveyard"))
			require.NoError(t, err, "a live object returned error on get")
			require.Len(t, val, 256)

			obj, err := db.Object(key, options.WithNamespace("graveyard"))
			require.NoError(t, err, "live object did not return an object")
			require.False(t, obj.Tombstone())
		}
	}

	// "Resurrect" every 4th tombstone and give it a new value
	for i, key := range keys {
		if i%4 == 0 {
			db.Put(key, randomData(192), options.WithNamespace("graveyard"))
		}
	}

	// Ensure that the iterator returns 15 items but that there are still 20 objects
	// including tombstones still stored in the database.
	requireNamespaceLen(t, db, "graveyard", 15)
	requireGraveyardLen(t, db, "graveyard", 20)
	requireDatabaseLen(t, db, 20)

	// Sanity check, attempt to get Get all keys and verify tombstones and undead keys
	for i, key := range keys {
		if i%2 == 0 {
			if i%4 == 0 {
				// This is an undead version
				val, err := db.Get(key, options.WithNamespace("graveyard"))
				require.NoError(t, err, "undead object returned error on get")
				require.Len(t, val, 192)

				obj, err := db.Object(key, options.WithNamespace("graveyard"))
				require.NoError(t, err, "undead object did not return an object")
				require.False(t, obj.Tombstone())
			} else {
				// This is a tombstone
				val, err := db.Get(key, options.WithNamespace("graveyard"))
				require.EqualError(t, err, "not found", "tombstone did not return a not found error")
				require.Nil(t, val, "tombstone returned a non nil value")

				obj, err := db.Object(key, options.WithNamespace("graveyard"))
				require.NoError(t, err, "tombstone did not return an object")
				require.True(t, obj.Tombstone())
			}
		} else {
			// Not a tombstone
			val, err := db.Get(key, options.WithNamespace("graveyard"))
			require.NoError(t, err, "a live object returned error on get")
			require.Len(t, val, 256)

			obj, err := db.Object(key, options.WithNamespace("graveyard"))
			require.NoError(t, err, "live object did not return an object")
			require.False(t, obj.Tombstone())
		}
	}

	// Test Seek, Next, and Prev with and without Tombstones
	iter, err := db.Iter(nil, options.WithNamespace("graveyard"))
	require.NoError(t, err, "could not create honu iterator")

	itert, err := db.Iter(nil, options.WithNamespace("graveyard"), options.WithTombstones())
	require.NoError(t, err, "could not create honu tombstone iterator")

	// Seek to a non-tombstone key
	require.True(t, iter.Seek(keys[9]), "could not seek to a non-tombstone key")
	require.True(t, itert.Seek(keys[9]), "could not seek to a non-tombstone key with tombstone iterator")
	require.True(t, bytes.Equal(iter.Key(), keys[9]), "unexpected key at iter cursor")
	require.True(t, bytes.Equal(itert.Key(), keys[9]), "unexpected key at iter cursor with tombstone iterator")

	// Seek to a tombstone key (move to 15 and 14 respectively)
	require.True(t, iter.Seek(keys[14]), "could not seek to a tombstone key")
	require.True(t, itert.Seek(keys[14]), "could not seek to a tombstone key with tombstone iterator")
	require.True(t, bytes.Equal(iter.Key(), keys[15]), "unexpected key at iter cursor")
	require.True(t, bytes.Equal(itert.Key(), keys[14]), "unexpected key at iter cursor with tombstone iterator")

	// Prev should move us to keys[13] for both two iterators
	require.True(t, iter.Prev(), "could not prev to a non-tombstone key")
	require.True(t, itert.Prev(), "could not prev to a non-tombstone key with tombstone iterator")
	require.True(t, bytes.Equal(iter.Key(), keys[13]), "unexpected key at iter cursor")
	require.True(t, bytes.Equal(itert.Key(), keys[13]), "unexpected key at iter cursor with tombstone iterator")

	// Next should move us back to 15 and 14 respectively
	require.True(t, iter.Next(), "could not next to a non-tombstone key")
	require.True(t, itert.Next(), "could not next to a tombstone key with tombstone iterator")
	require.True(t, bytes.Equal(iter.Key(), keys[15]), "unexpected key at iter cursor")
	require.True(t, bytes.Equal(itert.Key(), keys[14]), "unexpected key at iter cursor with tombstone iterator")
}

func TestTombstonesMultipleNamespaces(t *testing.T) {
	// Create a test database
	db, _ := setupHonuDB(t)
	namespaces := []string{"graveyard", "cemetery", "catacombs"}

	// Assert that there is nothing in the namespaces as an initial check
	for _, ns := range namespaces {
		requireNamespaceLen(t, db, ns, 0)
	}
	requireDatabaseLen(t, db, 0)

	// Create a list of keys with integer values
	keys := make([][]byte, 0, 100)
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		keys = append(keys, key)
	}

	// Add data to the database
	for _, key := range keys {
		for _, ns := range namespaces {
			db.Put(key, randomData(256), options.WithNamespace(ns))
		}
	}

	for _, ns := range namespaces {
		requireNamespaceLen(t, db, ns, 100)
	}
	requireDatabaseLen(t, db, 300)

	// Delete all even keys
	for i, key := range keys {
		if i%2 == 0 {
			for _, ns := range namespaces {
				db.Delete(key, options.WithNamespace(ns))
			}
		}
	}

	// Ensure that the iterator returns 50 items but that there are still 100 objects
	// including tombstones still stored in the database. Also ensure that the entire
	// database still contains 300 objects.
	for _, ns := range namespaces {
		requireNamespaceLen(t, db, ns, 50)
		requireGraveyardLen(t, db, ns, 100)
	}
	requireDatabaseLen(t, db, 300)

	// "Resurrect" every 4th tombstone and give it a new value
	for i, key := range keys {
		if i%4 == 0 {
			for _, ns := range namespaces {
				db.Put(key, randomData(192), options.WithNamespace(ns))
			}
		}
	}

	// Ensure that the iterator returns 75 items but that there are still 100 objects
	// including tombstones still stored in the database. Also ensure that the entire
	// database still contains 300 objects.
	for _, ns := range namespaces {
		requireNamespaceLen(t, db, ns, 75)
		requireGraveyardLen(t, db, ns, 100)
	}
	requireDatabaseLen(t, db, 300)
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

func requireNamespaceLen(t *testing.T, db *honu.DB, namespace string, expected int) {
	iter, err := db.Iter(nil, options.WithNamespace(namespace))
	require.NoError(t, err)

	actual := 0
	for iter.Next() {
		actual++
	}

	require.NoError(t, iter.Error())
	iter.Release()
	require.Equal(t, expected, actual)
}

func requireGraveyardLen(t *testing.T, db *honu.DB, namespace string, expected int) {
	iter, err := db.Iter(nil, options.WithNamespace(namespace), options.WithTombstones())
	require.NoError(t, err)

	actual := 0
	for iter.Next() {
		actual++
	}

	require.NoError(t, iter.Error())
	iter.Release()
	require.Equal(t, expected, actual)
}

func requireDatabaseLen(t *testing.T, db *honu.DB, expected int) {
	engine, ok := db.Engine().(*leveldb.LevelDBEngine)
	require.True(t, ok, "database len requires a leveldb engine")
	ldb := engine.DB()

	actual := 0
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		actual++
	}

	require.NoError(t, iter.Error(), "could not iterate using leveldb directly")
	iter.Release()

	require.Equal(t, expected, actual, "database key count does not match")
}

// Helper function to generate random data
func randomData(len int) []byte {
	data := make([]byte, len)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	return data
}
