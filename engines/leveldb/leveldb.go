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

func (db *LevelDBEngine) Engine() string {
	return "leveldb"
}

func (db *LevelDBEngine) Close() error {
	return db.ldb.Close()
}

// Get the latest version of the object stored by the key.
func (db *LevelDBEngine) Get(key []byte, options ...opts.SetOptions) (value []byte, err error) {
	var cfg *opts.Options
	cfg.LeveldbRead = nil
	for _, option := range options {
		if err = option(cfg); err != nil {
			return nil, err
		}
	}
	if value, err = db.ldb.Get(key, cfg.LeveldbRead); err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return value, engine.ErrNotFound
	}
	return value, err
}

// Put a new value to the specified key and update the version.
func (db *LevelDBEngine) Put(key, value []byte, options ...opts.SetOptions) (err error) {
	var cfg *opts.Options
	cfg.LeveldbWrite = nil
	for _, option := range options {
		if err = option(cfg); err != nil {
			return err
		}
	}
	return db.ldb.Put(key, value, cfg.LeveldbWrite)
}

// Delete the object represented by the key, creating a tombstone object.
func (db *LevelDBEngine) Delete(key []byte, options ...opts.SetOptions) (err error) {
	var cfg *opts.Options
	cfg.LeveldbWrite = nil
	for _, option := range options {
		if err = option(cfg); err != nil {
			return err
		}
	}
	return db.ldb.Delete(key, cfg.LeveldbWrite)
}

func (db *LevelDBEngine) Iter(prefix []byte) (i iterator.Iterator, err error) {
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	return NewLevelDBIterator(db.ldb.NewIterator(slice, nil)), nil
}
