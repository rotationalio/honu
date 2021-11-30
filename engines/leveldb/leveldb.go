package leveldb

import (
	"bytes"
	"errors"
	"sync"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/iterator"
	opts "github.com/rotationalio/honu/options"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Open a leveldb engine as the backend to the Honu database.
func Open(path string, conf config.ReplicaConfig) (_ *LevelDBEngine, err error) {
	engine := &LevelDBEngine{}
	if engine.ldb, err = leveldb.OpenFile(path, nil); err != nil {
		return nil, err
	}
	return engine, nil
}

// LevelDBEngine implements Engine and Store
type LevelDBEngine struct {
	sync.RWMutex
	ldb *leveldb.DB
}

// Transaction implements Transaction
type Transaction struct {
	db *LevelDBEngine
	ro bool
}

// Returns the name of the engine type.
func (db *LevelDBEngine) Engine() string {
	return "leveldb"
}

// Close the database and flush all remaining writes to disk. LevelDB requires a close
// for graceful shutdown to ensure there is no data loss.
func (db *LevelDBEngine) Close() error {
	return db.ldb.Close()
}

// Begin a multi-operation transaction (used primarily for Put and Delete)
func (db *LevelDBEngine) Begin(readonly bool) (engine.Transaction, error) {
	if readonly {
		db.RLock()
	} else {
		db.Lock()
	}
	return &Transaction{db: db, ro: readonly}, nil
}

// Finish a multi-operation transaction
func (tx *Transaction) Finish() error {
	if tx.ro {
		tx.db.RUnlock()
	} else {
		tx.db.Unlock()
	}
	return nil
}

// Get the latest version of the object stored by the key. This is the Transaction Get
// method which can be used in either readonly or write modes. This is the preferred
// mechanism to access the underlying engine.
func (tx *Transaction) Get(key []byte, options *opts.Options) (value []byte, err error) {
	return tx.db.get(key, options)
}

// Get the latest version of the object stored by the key. This is the Store Get method
// which can be used directly without a transaction. It is a unary operation that will
// read lock the database.
func (db *LevelDBEngine) Get(key []byte, options *opts.Options) (value []byte, err error) {
	db.RLock()
	defer db.RUnlock()
	return db.get(key, options)
}

// Thread-unsafe get that is called both by the Transaction and the Store.
func (db *LevelDBEngine) get(key []byte, options *opts.Options) (value []byte, err error) {
	//Create a default to prevent SegFaults when accessing options.
	if options == nil {
		if options, err = opts.New(); err != nil {
			return nil, err
		}
	}
	// Namespaces in leveldb are provided not by buckets but by namespace:: prefixed keys
	if options.Namespace != "" {
		key = prepend(options.Namespace, key)
	}

	if value, err = db.ldb.Get(key, options.LevelDBRead); err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return value, engine.ErrNotFound
	}
	return value, err
}

// Put a new value to the specified key. This is the Transaction Put method which is
// used by Honu to ensure consistency across version updates and can only be used in a
// write transaction.
func (tx *Transaction) Put(key, value []byte, options *opts.Options) (err error) {
	if tx.ro {
		return engine.ErrReadOnlyTx
	}
	return tx.db.put(key, value, options)
}

// Put a new value to the specified key and update the version. This is the Store Put
// method which can be used directly without a transaction. It is a unary operation that
// will lock the database to synchronize it with respect to other transactions.
func (db *LevelDBEngine) Put(key, value []byte, options *opts.Options) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.put(key, value, options)
}

// Thread-unsafe put that is called both by the Transaction and the Store.
func (db *LevelDBEngine) put(key, value []byte, options *opts.Options) (err error) {
	//Create a default to prevent SegFaults when accessing options.
	if options == nil {
		if options, err = opts.New(); err != nil {
			return err
		}
	}
	// Namespaces in leveldb are provided not by buckets but by namespace:: prefixed keys
	if options.Namespace != "" {
		key = prepend(options.Namespace, key)
	}
	return db.ldb.Put(key, value, options.LevelDBWrite)
}

// Delete the object represented by the key, removing it from the database entirely.
// This is the Transaction Delete method which is used by Honu to clean up and vacuum
// the database or to reset an object back to the first version. Note that normal Honu
// deletes Put a tombstone rather than directly deleting data from the database.
func (tx *Transaction) Delete(key []byte, options *opts.Options) (err error) {
	if tx.ro {
		return engine.ErrReadOnlyTx
	}
	return tx.db.delete(key, options)
}

// Delete the object represented by the key, removing it from the database entirely.
// This is the Store Delete method which can be used directly without a transaction. It
// is a unary operation that will lock the database to synchronize it with respect to
// other transactions. Note that normal Honu deletes Put a tombstone rather than
// directly deleting data from the database.
func (db *LevelDBEngine) Delete(key []byte, options *opts.Options) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.delete(key, options)
}

// Thread-unsafe delete that is called both by the Transaction and the Store.
func (db *LevelDBEngine) delete(key []byte, options *opts.Options) (err error) {
	//Create a default to prevent SegFaults when accessing options.
	if options == nil {
		if options, err = opts.New(); err != nil {
			return err
		}
	}
	// Namespaces in leveldb are provided not by buckets but by namespace:: prefixed keys
	if options.Namespace != "" {
		key = prepend(options.Namespace, key)
	}
	return db.ldb.Delete(key, options.LevelDBWrite)
}

func (db *LevelDBEngine) Iter(prefix []byte, options *opts.Options) (i iterator.Iterator, err error) {
	//Create a default to prevent SegFaults when accessing options.
	if options == nil {
		if options, err = opts.New(); err != nil {
			return nil, err
		}
	}
	// Namespaces in leveldb are provided not by buckets but by namespace:: prefixed keys
	if options.Namespace != "" {
		prefix = prepend(options.Namespace, prefix)
	}
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	return NewLevelDBIterator(db.ldb.NewIterator(slice, options.LevelDBRead)), nil
}

var nssep = []byte("::")

// prepend the namespace to the key
func prepend(namespace string, key []byte) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(namespace),
			key,
		}, nssep,
	)
}
