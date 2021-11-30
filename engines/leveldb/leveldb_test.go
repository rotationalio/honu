package leveldb_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/engines/leveldb"
	"github.com/rotationalio/honu/options"
	"github.com/stretchr/testify/require"
)

func getNamespaces() []string {
	return []string{
		"",
		"basic",
		"namespace with spaces",
		"namespace::with::colons",
	}
}

func setupLeveldbEngine() (_ *leveldb.LevelDBEngine, path string, err error) {
	tempDir, err := ioutil.TempDir("", "leveldb-*")
	ldbPath := fmt.Sprintf("leveldb:///%s", tempDir)
	if err != nil {
		return nil, "", err
	}
	conf := config.ReplicaConfig{}
	engine, err := leveldb.Open(ldbPath, conf)
	if err != nil {
		return nil, "", err

	}
	return engine, ldbPath, nil
}

func getOpts(namespace string, t *testing.T) *options.Options {
	opts, err := options.New(options.WithNamespace(namespace))
	require.NoError(t, err)
	return opts
}

func wrappedPut(ldbStore engine.Store, opts *options.Options, key []byte, value []byte, t *testing.T) {
	err := ldbStore.Put(key, value, opts)
	require.NoError(t, err)
}

func wrappedGet(ldbStore engine.Store, opts *options.Options, key []byte, expectedValue []byte, t *testing.T) {
	getValue, err := ldbStore.Get(key, opts)
	require.NoError(t, err)
	require.Equal(t, getValue, expectedValue)
}

func wrappedDelete(ldbStore engine.Store, opts *options.Options, key []byte, t *testing.T) {
	err := ldbStore.Delete(key, opts)
	require.NoError(t, err)

	value, err := ldbStore.Get(key, opts)
	require.Equal(t, err, engine.ErrNotFound)
	require.Empty(t, value)
}

func TestLeveldbEngine(t *testing.T) {
	//Setup a levelDB Engine
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)
	require.Equal(t, "leveldb", ldbEngine.Engine())

	//Ensure the db was created.
	require.DirExists(t, ldbPath)

	//Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	key := []byte("foo")

	for _, namespace := range getNamespaces() {
		opts := getOpts(namespace, t)
		value := []byte(namespace)
		wrappedPut(ldbEngine, opts, key, value, t)
		wrappedGet(ldbEngine, opts, key, value, t)
		wrappedDelete(ldbEngine, opts, key, t)
	}

	value := []byte("nil")
	wrappedPut(ldbEngine, nil, key, value, t)
	wrappedGet(ldbEngine, nil, key, value, t)
	wrappedDelete(ldbEngine, nil, key, t)
}

func TestLeveldbEngineWithNullOptions(t *testing.T) {
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)

	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	key := []byte("nilcheck")
	value := []byte("abc")
	wrappedPut(ldbEngine, nil, key, value, t)

	opts := getOpts("default", t)
	wrappedGet(ldbEngine, opts, key, value, t)
}

func TestLeveldbTransactions(t *testing.T) {
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)

	//Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	key := []byte("tx")

	//Start a transaction with readonly set to false.
	for _, namespace := range getNamespaces() {
		tx, err := ldbEngine.Begin(false)
		require.NoError(t, err)

		opts := getOpts(namespace, t)
		value := []byte(namespace)
		wrappedPut(tx, opts, key, value, t)
		wrappedGet(tx, opts, key, value, t)
		wrappedDelete(tx, opts, key, t)

		//Complete the transaction.
		require.NoError(t, tx.Finish())
	}

	//Create a transaction with readonly set to true.
	tx, err := ldbEngine.Begin(true)
	require.NoError(t, err)

	//Ensure Put fails with readonly set.
	err = tx.Put([]byte("bad"), []byte("write"), nil)
	require.Equal(t, err, engine.ErrReadOnlyTx)
	err = tx.Delete([]byte("bad"), nil)
	require.Equal(t, err, engine.ErrReadOnlyTx)

	//Complete the transaction.
	require.NoError(t, tx.Finish())
}

func TestLevelDBIter(t *testing.T) {
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)

	//Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	// Put a range of data into the database

	for _, pair := range [][]string{
		{"aa", "aa"},
		{"ab", "ab"},
		{"ba", "ba"},
		{"bb", "bb"},
		{"bc", "bc"},
		{"ca", "ca"},
		{"cb", "cb"},
	} {
		key := []byte(pair[0])
		value := []byte(pair[1])
		wrappedPut(ldbEngine, nil, key, value, t)
	}

	prefix := []byte("")
	iter, err := ldbEngine.Iter(prefix, nil)
	require.NoError(t, err)

	collected := 0
	for iter.Next() {
		key := iter.Key()
		fmt.Print(key)

		value := iter.Value()
		fmt.Println(value)

		collected++
	}
	fmt.Println(collected)
}
