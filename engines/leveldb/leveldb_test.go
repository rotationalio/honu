package leveldb_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/engines/leveldb"
	"github.com/stretchr/testify/require"
)

func setupLeveldbEngine() (_ *leveldb.LevelDBEngine, path string, err error) {
	tempDir, err := ioutil.TempDir("", "leveldb-*")
	ldbPath := fmt.Sprintf("leveldb:///%s", tempDir)
	if err != nil {
		return nil, "", err
	}
	conf := config.ReplicaConfig{}
	ldb, err := leveldb.Open(ldbPath, conf)
	if err != nil {
		return nil, "", err

	}
	return ldb, ldbPath, nil
}

func TestLeveldbEngine(t *testing.T) {
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)

	//Ensure the db was created.
	require.DirExists(t, ldbPath)

	//Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer require.NoError(t, ldbEngine.Close())

	require.Equal(t, "leveldb", ldbEngine.Engine())

	//Perform a Put without options/namespaces.
	key := []byte("put_without")
	putValue := []byte("options")
	err = ldbEngine.Put(key, putValue, nil)
	require.NoError(t, err)

	//Read the previously Put value without options/namespaces.
	Getvalue, err := ldbEngine.Get(key, nil)
	require.NoError(t, err)
	require.Equal(t, putValue, Getvalue)
}

func TestLeveldbTransactions(t *testing.T) {
	ldbEngine, ldbPath, err := setupLeveldbEngine()
	require.NoError(t, err)

	//Teardown after finishing the test
	defer os.RemoveAll(ldbPath)
	defer ldbEngine.Close()

	//Start a transaction with readonly set to false.
	tx, err := ldbEngine.Begin(false)
	require.NoError(t, err)

	//
	//
	//

	//Complete the transaction.
	require.NoError(t, tx.Finish())

	//Create a transaction with readonly set to true.
	tx, err = ldbEngine.Begin(true)
	require.NoError(t, err)

	//Ensure Put fails with readonly set.
	err = tx.Put([]byte("bad"), []byte("write"), nil)
	require.Equal(t, err, engine.ErrReadOnlyTx)

	//Complete the transaction.
	require.NoError(t, tx.Finish())
}
