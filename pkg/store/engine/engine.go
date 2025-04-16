package engine

import (
	"go.rtnl.ai/honu/pkg/store/iterator"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/object"
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
// all engines should support the Store interface. Note that transactions and options
// are a higher level construct than engine stores and are provided by the Honu store.
type Store interface {
	Has(key.Key) (exists bool, err error)
	Get(key.Key) (object.Object, error)
	Put(key.Key, object.Object) error
	Delete(key.Key) error
}

// Iterator engines allow queries that scan a range of consecutive keys.
type Iterator interface {
	Iter(prefix []byte) (i iterator.Iterator, err error)
	Range(start, limit []byte) (i iterator.Iterator, err error)
}
