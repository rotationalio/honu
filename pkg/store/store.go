package store

import (
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/store/engine"
	"go.rtnl.ai/honu/pkg/store/engine/leveldb"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/locks"
)

// Store implements local database functionality for interaction with objects and their
// metadata on disk. All external accessors of the store (with the possible exception
// of database backups) should use the Store to ensure proper isolation, consistency,
// durability, and atomicity of operations.
//
// The Store is thread-safe and can be used safely from multiple goroutines. Go
// routines should provide a cancelable context to ensure database operations do not
// proceed after cancellation.
//
// The Store maintains the versioning of object accesses, so all writes must be
// serialized through the store. Additionally the Store maintains all of the indexes
// associated with the database, and maintains all constraints such as uniqueness.
type Store struct {
	lamport.Clock
	pid lamport.PID
	db  engine.Engine
	mu  locks.Keys
}

func Open(conf config.Config) (s *Store, err error) {
	s = &Store{
		Clock: lamport.New(conf.PID),
		pid:   lamport.PID(conf.PID),
		mu:    locks.New(conf.Store.Concurrency),
	}

	if s.db, err = leveldb.Open(conf.Store); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

//===========================================================================
// Collection Management
//===========================================================================

//===========================================================================
// Object Management
//===========================================================================

func (s *Store) List() error {
	return nil
}

func (s *Store) Query() error {
	return nil
}

func (s *Store) Create() error {
	return nil
}

func (s *Store) Retrieve() error {
	return nil
}

func (s *Store) Update() error {
	return nil
}

func (s *Store) Delete() error {
	return nil
}
