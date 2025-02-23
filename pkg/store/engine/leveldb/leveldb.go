package leveldb

import (
	"errors"
	"fmt"

	"github.com/rotationalio/honu/pkg/config"
	"github.com/rotationalio/honu/pkg/store/engine"
	"github.com/rotationalio/honu/pkg/store/iterator"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/object"
	"github.com/rotationalio/honu/pkg/store/opts"
	"github.com/syndtr/goleveldb/leveldb"
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
	return "honu.LevelDBEngine"
}

// Close the database and flush all remaining writes to disk. LevelDB requires a close
// during graceful shutdown to ensure there is no data loss or corruption.
func (e *Engine) Close() error {
	return e.ldb.Close()
}

//===========================================================================
// Implement engine.Store Interface
//===========================================================================

func (e *Engine) Has(key key.Key, ro *opts.ReadOptions) (exists bool, err error) {
	// TODO: tombstone handling
	if exists, err = e.ldb.Has(key[:], nil); err != nil {
		return exists, Wrap(err)
	}
	return exists, nil
}

func (e *Engine) Get(key key.Key, ro *opts.ReadOptions) (_ object.Object, err error) {
	// TODO: tombstone handling
	var val []byte
	if val, err = e.ldb.Get(key[:], nil); err != nil {
		return nil, Wrap(err)
	}
	return object.Object(val), nil
}

func (e *Engine) Put(key key.Key, obj object.Object, wo *opts.WriteOptions) (err error) {
	// TODO: do we need a transaction here?
	if wo.GetNoOverwrite() || wo.GetCheckUpdate() {
		var exists bool
		if exists, err = e.ldb.Has(key, nil); err != nil {
			return Wrap(err)
		}

		if exists && wo.GetNoOverwrite() {
			return engine.ErrAlreadyExists
		}

		if !exists && wo.GetCheckUpdate() {
			return engine.ErrNotFound
		}
	}

	if err = e.ldb.Put(key[:], obj[:], nil); err != nil {
		return Wrap(err)
	}
	return nil
}

func (e *Engine) Delete(key key.Key, wo *opts.WriteOptions) (err error) {
	// TODO: do we need a transaction here?
	if wo.CheckDelete {
		var exists bool
		if exists, err = e.ldb.Has(key, nil); err != nil {
			return Wrap(err)
		}

		if !exists {
			return engine.ErrNotFound
		}
	}

	if err = e.ldb.Delete(key[:], nil); err != nil {
		return Wrap(err)
	}
	return nil
}

//===========================================================================
// Implement engine.Iterator Interface
//===========================================================================

func (e *Engine) Iter(prefix []byte, ro *opts.ReadOptions) (_ iterator.Iterator, err error) {
	return nil, nil
}

//===========================================================================
// Error Handling
//===========================================================================

func Wrap(err error) error {
	switch {
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
