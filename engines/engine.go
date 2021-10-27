package engine

import (
	"github.com/rotationalio/honu/iterator"
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
	Get(key []byte) (value []byte, err error)
	Put(key, value []byte) error
	Delete(key []byte) error
}

// Iterator engines allow queries that scan a range of consecutive keys.
type Iterator interface {
	Iter(prefix []byte) (i iterator.Iterator, err error)
}
