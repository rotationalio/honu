package store

import (
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/store/engine"
	"go.rtnl.ai/honu/pkg/store/engine/leveldb"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/locks"
	"go.rtnl.ai/ulid"
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

// Collections are subsets of the Store that allow access to related objects. Each
// object in a collection is prefixed by the collection ID, ensuring that the objects
// are grouped together and can be accessed efficiently.
type Collection struct {
	pid lamport.PID
	db  engine.Engine
	mu  locks.Keys

	ID   ulid.ULID
	Name string
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
		pid: s.pid,
		db:  s.db,
		mu:  s.mu,

		ID:   ulid.Make(),
		Name: name,
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

//===========================================================================
// Object Management
//===========================================================================

// List all of the objects in the collection, returning an iterator that will allow the
// caller to either simply iterate over the keys or to actually retreive the objects in
// a memory-efficient manner.
func (c *Collection) List() error {
	return nil
}

// List all of the objects in the collection that match the specified query. An iterator
// is returned that will allow the caller to either simply iterate over the keys
// or to actually retrieve the objects in a memory-efficient manner.
func (c *Collection) Query() error {
	return nil
}

// Empty the collection by adding a tombstone version to all of the objects in the
// collection. These objects cannot be accessed directly any longer but their version
// history is preserved. This is different from the Collection.Truncate method which
// removes all objects and their versions from the collection.
func (c *Collection) Empty() error {
	return nil
}

// Create a new object in the collection with the given key and value. The key must be
// unique within the collection, and if it already exists, it will return an error.
// Note that because of the replicated nature of Honu, we can't guarantee that the
// object was uniquely created, however it will guarantee that the object version
// history starts from the current object and branches can be detected later.
func (c *Collection) Create() error {
	return nil
}

// Retrieve the latest version of the object with the given key from the collection. If
// the object is a tombstone record or if the key is not in the store, then a not found
// error will be returned. If a version is specified that version will be retrieved,
// even if it is a tombstone record; version does not exist is returned instead of
// not found in this case.
func (c *Collection) Retrieve() error {
	return nil
}

// Returns an iterator of all versions of the object; iterating from the most recent
// version to the oldest. Tombstone versions are included by the iterator.
func (c *Collection) Versions() error {
	return nil
}

// Create a new version record of the object for the given key. If the object does not
// already exist, it will return an error. Because of the replicated nature of Honu,
// we can't guarantee that the object doesn't exist somewhere else in the cluster but
// this will prevent updates locally until that created version is replicated.
func (c *Collection) Update() error {
	return nil
}

// Merge performs an upsert operation on the object, creating a new version of the key
// if it does not exist, or updating the existing version if it does. Merge provides
// simpler semantics than Create or Update as the caller does not need to worry about
// whether the object exists on the cluster or not, and in single replica queries its
// better to use Merge.
func (c *Collection) Merge() error {
	return nil
}

// Delete an object from the collection by adding a tombstone version; the object will
// not be returned in list queries or retrieval but the version history of the object
// will be preserved.
func (c *Collection) Delete() error {
	return nil
}

// Destroy the object and all of its versions from the collection. This method adds a
// truncated record to the object, which is replicated to all replicas. Any object that
// gets created with the same key in the future will start from version 1, even if
// the truncation happens concurrently with the creation of the new object.
func (c *Collection) Destroy() error {
	return nil
}
