package engine

import (
	"github.com/rotationalio/honu/iterator"
	opts "github.com/rotationalio/honu/options"
)

// Engines are the disk storage mechanism that Honu wraps. Users may chose different
// engines for a variety of reasons, including variable performance benefits, different
// features, or even implement heterogeneous Honu networks composed of different engines.
type Engine interface {
	// Engine returns the engine name and is used for debugging and logging.
	Engine() string

	// Close the engine so that it no longer can be accessed.
	Close() error

	// Begin a transaction to issue multiple commands (this is an internal transaction
	// for Honu-specific version management, not an external interface).
	Begin(readonly bool) (tx Transaction, err error)
}

// Store is a simple key/value interface that allows for Get, Put, and Delete. Nearly
// all engines should support the Store interface.
type Store interface {
	Has(key []byte, options *opts.Options) (exists bool, err error)
	Get(key []byte, options *opts.Options) (value []byte, err error)
	Put(key, value []byte, options *opts.Options) error
	Delete(key []byte, options *opts.Options) error
}

// Iterator engines allow queries that scan a range of consecutive keys.
type Iterator interface {
	Iter(prefix []byte, options *opts.Options) (i iterator.Iterator, err error)
}

type Transaction interface {
	Store
	Finish() error
}
