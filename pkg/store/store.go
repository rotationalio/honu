package store

import (
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/store/engine"
	"go.rtnl.ai/honu/pkg/store/engine/leveldb"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/locks"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/ulid"
)

// System collections are used for internal database management and are not named but
// are directly identified by a unique ID that is less than any ULID that would be
// generated since they are before
var (
	SystemPrefix        = [5]byte{0x00, 0x68, 0x6f, 0x6e, 0x75}
	SystemCollections   = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e})
	SystemReplicas      = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x01, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6c, 0x69, 0x73, 0x74})
	SystemAccessControl = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x02, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67})
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
	pid lamport.PID
	db  engine.Engine
	mu  locks.Keys

	collections map[ulid.ULID]*Collection
}

// Open a new Store with the provided configuration. Only one Store can be opened for a
// single configuration (database path), which locks the files on disk to prevent other
// processes from opening the same database. If the database already exists, it will be
// opened and the existing data will be used. If the database does not exist, it will
// be created. The Store can then be used by multiple goroutines, to sequence access
// to disk.
func Open(conf config.Config) (s *Store, err error) {
	s = &Store{
		pid:         lamport.PID(conf.PID),
		mu:          locks.New(conf.Store.Concurrency),
		collections: make(map[ulid.ULID]*Collection),
	}

	if s.db, err = leveldb.Open(conf.Store); err != nil {
		return nil, err
	}

	// TODO: should we load the collections from disk here?
	return s, nil
}

// Close the store and release all resources associated with it. This will close the
// underlying database engine and release any locks held by the store.
func (s *Store) Close() error {
	// Acquire all locks to wait for any ongoing operations to complete and to ensure
	// that no new operations can start while closing the store.
	s.mu.LockAll()
	defer s.mu.UnlockAll()
	return s.db.Close()
}

//===========================================================================
// Collection Management
//===========================================================================

// Returns all of the collections that the store is maintaining including any system
// collections that are used for internal database management.
func (s *Store) Collections() []*Collection {
	// TODO: should we check the collections on disk?
	collections := make([]*Collection, 0, len(s.collections))
	for _, c := range s.collections {
		collections = append(collections, c)
	}
	return collections
}

// Creates a new collection with the given name and associates it with a unique ID.
// The name is case-insensitive and should be unique within the store. It must contain
// only no spaces or punctuation and cannot start with a number. The name also must not
// be a ULID string, which is reserved for collection IDs.
func (s *Store) New(name string) (*Collection, error) {
	// TODO: lock the collection key
	// TODO: check name in index to ensure uniqueness.
	// TODO: normalize name to ensure it is valid and does not contain any illegal characters.
	// TODO: save the collection index to disk.
	collection := &Collection{
		Collection: metadata.Collection{
			ID:   ulid.Make(),
			Name: name,
		},

		pid: s.pid,
		db:  s.db,
		mu:  s.mu,
	}

	s.collections[collection.ID] = collection
	return collection, nil
}

// Opens an existing collection either by its ID or by its name. If the collection does
// not exist, it will return an error. The collection is ready for access when returned.
func (s *Store) Open(identifier any) (*Collection, error) {
	return nil, nil
}

// Drop a collection, removing it from the store and deleting all of its contained
// objects. The collection can be recreated later with the same name, but all of the
// previous objects and their versions will be deleted.
func (s *Store) Drop(identifier any) error {
	return nil
}

// Truncate a collection, removing all of its contained objects and versions, but
// keeping the collection and its indexes intact. This is fundamenally different to the
// Collection.Empty method which adds a new tombstone version to remove objects from
// the collection but keeps the version history intact.
//
// NOTE: truncate adds a truncated record to every object in the collection, thereby
// replicating the truncation operation to all replicas. Objects that are added after
// the truncation operation will start from version 1, even if that truncation happens
// concurrently with the creation of the new object.
func (s *Store) Truncate(identifier any) error {
	return nil
}
