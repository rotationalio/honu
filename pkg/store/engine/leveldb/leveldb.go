package leveldb

import (
	"errors"
	"fmt"

	"github.com/rotationalio/honu/pkg/config"
	"github.com/rotationalio/honu/pkg/store/engine"
	"github.com/rotationalio/honu/pkg/store/iterator"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/object"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func Open(conf config.StoreConfig) (engine *Engine, err error) {
	engine = &Engine{}
	opts := Options(conf)

	if engine.ldb, err = leveldb.OpenFile(conf.DataPath, opts); err != nil {
		return nil, fmt.Errorf("could not open leveldb: %w", err)
	}
	return engine, nil
}

type Engine struct {
	ldb *leveldb.DB
}

var (
	_ engine.Engine   = (*Engine)(nil)
	_ engine.Store    = (*Engine)(nil)
	_ engine.Iterator = (*Engine)(nil)
)

// Returns the underlying database object for direct access (e.g. for batches or backups).
func (e *Engine) DB() *leveldb.DB {
	return e.ldb
}

//===========================================================================
// Implement engine.Engine Interface
//===========================================================================

// Returns the name of the engine type
func (e *Engine) Engine() string {
	return "leveldb"
}

// Close the database and flush all remaining writes to disk. LevelDB requires a close
// during graceful shutdown to ensure there is no data loss or corruption.
func (e *Engine) Close() error {
	return e.ldb.Close()
}

//===========================================================================
// Implement engine.Store Interface
//===========================================================================

func (e *Engine) Has(key key.Key) (exists bool, err error) {
	if exists, err = e.ldb.Has(key, nil); err != nil {
		return exists, Wrap(err)
	}
	return exists, nil
}

func (e *Engine) Get(key key.Key) (_ object.Object, err error) {
	var val []byte
	if val, err = e.ldb.Get(key, nil); err != nil {
		return nil, Wrap(err)
	}
	return object.Object(val), nil
}

func (e *Engine) Put(key key.Key, obj object.Object) (err error) {
	if err = e.ldb.Put(key, obj, nil); err != nil {
		return Wrap(err)
	}
	return nil
}

func (e *Engine) Delete(key key.Key) (err error) {
	if err = e.ldb.Delete(key, nil); err != nil {
		return Wrap(err)
	}
	return nil
}

//===========================================================================
// Implement engine.Iterator Interface
//===========================================================================

func (e *Engine) Iter(prefix []byte) (_ iterator.Iterator, err error) {
	return NewIterator(e.ldb.NewIterator(util.BytesPrefix(prefix), nil)), nil
}

func (e *Engine) Range(start, limit []byte) (_ iterator.Iterator, err error) {
	return NewIterator(e.ldb.NewIterator(&util.Range{Start: start, Limit: limit}, nil)), nil
}

//===========================================================================
// Error Handling
//===========================================================================

func Wrap(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, leveldb.ErrNotFound):
		return engine.ErrNotFound
	case errors.Is(err, leveldb.ErrReadOnly):
		return engine.ErrReadOnlyDB
	case errors.Is(err, leveldb.ErrIterReleased):
		return iterator.ErrIterReleased
	case errors.Is(err, leveldb.ErrClosed):
		return engine.ErrClosed
	default:
		return fmt.Errorf("could not complete engine operation: %w", err)
	}
}
