package leveldb

import (
	"errors"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/iterator"
	opts "github.com/rotationalio/honu/options"
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
	return d.ldb.Close()
}

// Get the latest version of the object stored by the key.
func (d *LevelDBEngine) Get(key []byte, options string) (value []byte, err error) {
	opts := opts.LeveldbOptions{}
	readOptions, err := opts.Read(options)
	if err != nil {
		return nil, err
	}
	if value, err = d.ldb.Get(key, readOptions); err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return value, engine.ErrNotFound
	}
	return value, err
}

// Put a new value to the specified key and update the version.
func (d *LevelDBEngine) Put(key, value []byte, options string) (err error) {
	opts := opts.LeveldbOptions{}
	writeOptions, err := opts.Write(options)
	if err != nil {
		return err
	}
	return d.ldb.Put(key, value, writeOptions)
}

// Delete the object represented by the key, creating a tombstone object.
func (d *LevelDBEngine) Delete(key []byte, options string) (err error) {
	opts := opts.LeveldbOptions{}
	writeOptions, err := opts.Write(options)
	if err != nil {
		return err
	}
	return d.ldb.Delete(key, writeOptions)
}

func (d *LevelDBEngine) Iter(prefix []byte) (i iterator.Iterator, err error) {
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	return NewLevelDBIterator(d.ldb.NewIterator(slice, nil)), nil
}
