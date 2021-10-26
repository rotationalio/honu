package leveldb

import (
	"errors"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/iterator"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func Open(path string, conf config.ReplicaConfig) (_ *LevelDBEngine, err error) {
	engine := &LevelDBEngine{}
	if engine.ldb, err = leveldb.OpenFile(path, nil); err != nil {
		return nil, err
	}
	return engine, nil
}

type LevelDBEngine struct {
	ldb *leveldb.DB
}

func (d *LevelDBEngine) Engine() string {
	return "leveldb"
}

func (d *LevelDBEngine) Close() error {
	return nil
}

// Get the latest version of the object stored by the key.
// TODO: provide read options to the underlying database.
func (d *LevelDBEngine) Get(key []byte) (value []byte, err error) {
	if value, err = d.ldb.Get(key, nil); err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return value, engine.ErrNotFound
	}
	return value, err
}

// Put a new value to the specified key and update the version.
// TODO: provide write options to the underlying database.
func (d *LevelDBEngine) Put(key, value []byte) (err error) {
	return d.ldb.Put(key, value, nil)
}

// Delete the object represented by the key, creating a tombstone object.
// TODO: provide write options to the underlying database.
func (d *LevelDBEngine) Delete(key []byte) (err error) {
	return d.ldb.Delete(key, nil)
}

func (d *LevelDBEngine) Iter(prefix []byte) (i iterator.Iterator, err error) {
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	return NewLevelDBIterator(d.ldb.NewIterator(slice, nil)), nil
}
