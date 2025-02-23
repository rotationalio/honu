package engine

import (
	"github.com/rotationalio/honu/pkg/store/iterator"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/object"
	"github.com/rotationalio/honu/pkg/store/opts"
)

// Engines are the disk storage mechanism that Honu wraps. Users may chose different
// engines for a variety of reasons, including variable performance benefits, different
// features, or even implement heterogeneous Honu networks composed of different engines.
type Engine interface {
	// Engine returns the engine name and is used for debugging and logging.
	Engine() string

	// Close the engine so that it no longer can be accessed.
	Close() error
}

// Store is a simple key/value interface that allows for Get, Put, and Delete. Nearly
// all engines should support the Store interface.
type Store interface {
	Has(key.Key, *opts.ReadOptions) (exists bool, err error)
	Get(key.Key, *opts.ReadOptions) (object.Object, error)
	Put(key.Key, object.Object, *opts.WriteOptions) error
	Delete(key.Key, *opts.WriteOptions) error
}

// Iterator engines allow queries that scan a range of consecutive keys.
type Iterator interface {
	Iter(prefix []byte, ro *opts.ReadOptions) (i iterator.Iterator, err error)
}
