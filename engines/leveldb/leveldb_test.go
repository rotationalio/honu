package leveldb_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/engines/leveldb"
	"github.com/rotationalio/honu/iterator"
	pb "github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// Returns a constant list of namespace strings.
func getNamespaces() []string {
	return []string{
		"",
		"basic",
		"namespace with spaces",
		"namespace::with::colons",
	}
}

// Returns a constant pair of Key/Value string pairs.
func getIterPairs() [][]string {
	return [][]string{
		{"aa", "first"},
		{"ab", "second"},
		{"ba", "third"},
		{"bb", "fourth"},
		{"bc", "fifth"},
		{"ca", "sixth"},
		{"cb", "seventh"},
	}
}

// Returns a LevelDBEngine and the path were it was created.
func setupLeveldbEngine(t *testing.T) (_ *leveldb.LevelDBEngine, path string) {
	tempDir, err := ioutil.TempDir("", "leveldb-*")
	ldbPath := fmt.Sprintf("leveldb:///%s", tempDir)
	require.NoError(t, err)
	conf := config.ReplicaConfig{}
	engine, err := leveldb.Open(ldbPath, conf)
	require.NoError(t, err)
	return engine, ldbPath
}

// Creates an options.Options struct with namespace set and returns
// a pointer to it.
func getOpts(namespace string, t *testing.T) *options.Options {
	opts, err := options.New(options.WithNamespace(namespace))
	require.NoError(t, err)
	return opts
}

// Wraps engine.Store.Put with testing checks.
func wrappedPut(ldbStore engine.Store, opts *options.Options, key []byte, value []byte, t *testing.T) {
	err := ldbStore.Put(key, value, opts)
	require.NoError(t, err)
}

// Wraps engine.Store.Get with testing checks.
func wrappedGet(ldbStore engine.Store, opts *options.Options, key []byte, expectedValue []byte, t *testing.T) {
	getValue, err := ldbStore.Get(key, opts)
	require.NoError(t, err)
	require.Equal(t, getValue, expectedValue)
}

// Wraps engine.Store.Delete with testing checks.
func wrappedDelete(ldbStore engine.Store, opts *options.Options, key []byte, t *testing.T) {
	err := ldbStore.Delete(key, opts)
	require.NoError(t, err)

	value, err := ldbStore.Get(key, opts)
	require.Equal(t, err, engine.ErrNotFound)
	require.Empty(t, value)
}

func TestLeveldbEngine(t *testing.T) {
	// Setup a levelDB Engine.
	ldbEngine, ldbPath := setupLeveldbEngine(t)
	require.Equal(t, "leveldb", ldbEngine.Engine())

	// Ensure the db was created.
	require.DirExists(t, ldbPath)

	// Teardown after finishing the test.
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	// Use a constant key to ensure namespaces
	// are working correctly.
	key := []byte("foo")

	// Check Put, Get and Delete with a nil namespace.
	value := []byte("nil")
	wrappedPut(ldbEngine, nil, key, value, t)
	wrappedGet(ldbEngine, nil, key, value, t)
	wrappedDelete(ldbEngine, nil, key, t)

	// Iterate through a list of namespaces and ensure
	// Put, Get and Delete are working.
	for _, namespace := range getNamespaces() {
		opts := getOpts(namespace, t)
		value := []byte(namespace)
		wrappedPut(ldbEngine, opts, key, value, t)
		wrappedGet(ldbEngine, opts, key, value, t)
		wrappedDelete(ldbEngine, opts, key, t)
	}
}

func TestLeveldbTransactions(t *testing.T) {
	ldbEngine, ldbPath := setupLeveldbEngine(t)

	// Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	// Use a constant key to ensure namespaces
	// are working correctly.
	key := []byte("tx")

	// Iterate through a list of namespaces and ensure
	// Put, Get and Delete are working with transactions.
	for _, namespace := range getNamespaces() {
		// Start a transaction with readonly set to false.
		tx, err := ldbEngine.Begin(false)
		require.NoError(t, err)

		opts := getOpts(namespace, t)
		value := []byte(namespace)
		wrappedPut(tx, opts, key, value, t)
		wrappedGet(tx, opts, key, value, t)
		wrappedDelete(tx, opts, key, t)

		// Complete the transaction.
		require.NoError(t, tx.Finish())
	}

	// Begin a transaction with readonly set to true.
	tx, err := ldbEngine.Begin(true)
	require.NoError(t, err)

	// Ensure Put fails with readonly set.
	err = tx.Put([]byte("bad"), []byte("write"), nil)
	require.Equal(t, err, engine.ErrReadOnlyTx)
	err = tx.Delete([]byte("bad"), nil)
	require.Equal(t, err, engine.ErrReadOnlyTx)

	// Complete the transaction.
	require.NoError(t, tx.Finish())
}

func TestLevelDBIter(t *testing.T) {
	ldbEngine, ldbPath := setupLeveldbEngine(t)

	// Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	for _, namespace := range getNamespaces() {
		// TODO: figure out what to do with this testcase.
		// Iter currently grabs the namespace by splitting
		// on :: and grabbing the first string, so it only
		// grabs "namespace".
		if namespace == "namespace::with::colons" {
			continue
		}
		// Add data to the database to iterate over.
		opts := getOpts(namespace, t)
		pairs := getIterPairs()
		addIterPairsToDB(ldbEngine, opts, pairs, t)

		// Try to iterate over all keys
		prefix := []byte("")
		iter, err := ldbEngine.Iter(prefix, opts)
		require.NoError(t, err)
		collected := iterate(iter, prefix, pairs, t)
		require.Equal(t, len(pairs), collected)

		// Try to iterate over all keys that start with "b"
		prefix = []byte("b")
		iter, err = ldbEngine.Iter(prefix, opts)
		require.NoError(t, err)
		bPairs := [][]string{
			pairs[2],
			pairs[3],
			pairs[4],
		}
		collected = iterate(iter, prefix, bPairs, t)
		require.Equal(t, 3, collected)
	}
}

// Helper function for TestLevelDBIter that iterates over the keys starting
// with 'prefix' using iter, checking that the values are correct with 'pairs'
// and returns the number of values iterated over.
func iterate(iter iterator.Iterator, prefix []byte, pairs [][]string, t *testing.T) int {
	collected := 0
	for iter.Next() {
		key := string(iter.Key())
		expectedKey := pairs[collected][0]
		require.Equal(t, expectedKey, key)

		value := string(iter.Value())
		expectedValue := pairs[collected][1]
		require.Equal(t, expectedValue, value)

		collected++
	}
	return collected
}

// Helper function for TestLevelDBIter that puts a range of data into the database.
func addIterPairsToDB(ldbStore engine.Store, opts *options.Options, pairs [][]string, t *testing.T) {
	for _, pair := range pairs {
		key := []byte(pair[0])
		value := []byte(pair[1])

		obj := &pb.Object{
			Key:       key,
			Namespace: opts.Namespace,
			Data:      value,
		}

		data, err := proto.Marshal(obj)
		require.NoError(t, err)
		wrappedPut(ldbStore, opts, key, data, t)
	}
}
