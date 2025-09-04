package store

import (
	"errors"
	"time"

	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/store/engine"
	"go.rtnl.ai/honu/pkg/store/engine/leveldb"
	"go.rtnl.ai/honu/pkg/store/index"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/locks"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
	"go.rtnl.ai/ulid"
)

// System collections are used for internal database management and are not named but
// are directly identified by a unique ID that is less than any ULID that would be
// generated since they are before
var (
	SystemPrefix        = [5]byte{0x00, 0x68, 0x6f, 0x6e, 0x75}
	SystemHonuAgent     = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x68, 0x6f, 0x6e, 0x75, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x00})
	SystemCollections   = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e})
	SystemReplicas      = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x01, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6c, 0x69, 0x73, 0x74})
	SystemAccessControl = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x02, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67})
)

// Collection or store related errors.
var (
	ErrCreateID = errors.New("create collection: cannot specify ID")
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

	collectionNames index.Index[string, ulid.ULID] // Index of collection names to their IDs
	collections     map[ulid.ULID]*Collection      // In memory collection objects to reduce allocations
}

// Open a new Store with the provided configuration. Only one Store can be opened for a
// single configuration (database path), which locks the files on disk to prevent other
// processes from opening the same database. If the database already exists, it will be
// opened and the existing data will be used. If the database does not exist, it will
// be created. The Store can then be used by multiple goroutines, to sequence access
// to disk.
func Open(conf config.Config) (s *Store, err error) {
	s = &Store{
		pid:             lamport.PID(conf.PID),
		mu:              locks.New(conf.Store.Concurrency),
		collections:     make(map[ulid.ULID]*Collection),
		collectionNames: index.Map[string, ulid.ULID](make(map[string]ulid.ULID)),
	}

	if s.db, err = leveldb.Open(conf.Store); err != nil {
		return nil, err
	}

	// Ensure the database is initialized and ready for use.
	if err = s.initialize(); err != nil {
		s.db.Close()
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

	// Clean up in-memory items and indexes.
	s.collections = nil
	s.collectionNames = nil

	return s.db.Close()
}

//===========================================================================
// Collection Management
//===========================================================================

// Returns a list of the collections that the store is maintaining excluding any system
// collections that are used for internal database management. Only the metadata is
// returned for each collection, so any collection operations must be performed after
// opening the collection by its ID or name.
func (s *Store) Collections() []metadata.Collection {
	// TODO: should we check the collections on disk?
	collections := make([]metadata.Collection, 0, len(s.collections))
	for _, c := range s.collections {
		collections = append(collections, c.Collection)
	}
	return collections
}

// Creates a new collection with the given name and associates it with a unique ID.
// The name is case-insensitive and should be unique within the store. It must contain
// only no spaces or punctuation and cannot start with a number. The name also must not
// be a ULID string, which is reserved for collection IDs.
func (s *Store) New(info *metadata.Collection) (_ *Collection, err error) {
	// Validate the info to ensure a collection can be created.
	if err = info.Validate(); err != nil {
		return nil, err
	}

	// A collection must not have an ID set.
	if !info.ID.IsZero() {
		return nil, ErrCreateID
	}

	// TODO: lock the collection key
	// TODO: check name in index to ensure uniqueness.
	// TODO: normalize name to ensure it is valid and does not contain any illegal characters.
	// TODO: save the collection index to disk.
	collection := &Collection{
		Collection: *info,

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

// Modifies the metadata of an existing collection; the collection should either have
// an ID or a name for reference and must already exist in the store.
func (s *Store) Modify(info *metadata.Collection) error {
	// Validate the info to ensure a collection can be modified.
	if err := info.Validate(); err != nil {
		return err
	}

	return nil
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

// Returns the underlying engine that the store is using for persistence. This is
// primarily used for testing and debugging purposes, and should be used with caution.
func (s *Store) Engine() engine.Engine {
	return s.db
}

//===========================================================================
// Initialization
//===========================================================================

func (s *Store) initialize() (err error) {
	// Ensures that the system collections are created and initialized in the store and
	// that the system user is created with the correct permissions.
	now := time.Now()
	defaultCollections := []*metadata.Collection{
		{
			ID:   SystemCollections,
			Name: "honu_collections",
			Version: &metadata.Version{
				Scalar:  lamport.Scalar{PID: 0, VID: 1}, // This is the default first version for system collections.
				Created: now,
			},
			Owner:    SystemHonuAgent,
			Group:    SystemHonuAgent,
			Created:  now,
			Modified: now,
		},
		{
			ID:   SystemReplicas,
			Name: "honu_replicas",
			Version: &metadata.Version{
				Scalar:  lamport.Scalar{PID: 0, VID: 1}, // This is the default first version for system collections.
				Created: now,
			},
			Owner:    SystemHonuAgent,
			Group:    SystemHonuAgent,
			Created:  now,
			Modified: now,
		},
		{
			ID:   SystemAccessControl,
			Name: "honu_access_control",
			Version: &metadata.Version{
				Scalar:  lamport.Scalar{PID: 0, VID: 1}, // This is the default first version for system collections.
				Created: now,
			},
			Owner:    SystemHonuAgent,
			Group:    SystemHonuAgent,
			Created:  now,
			Modified: now,
		},
	}

	for _, collection := range defaultCollections {
		collectionKey := key.New(SystemCollections, collection.ID, &collection.Version.Scalar)

		var exists bool
		if exists, err = s.db.Has(collectionKey); err != nil {
			return err
		}

		if !exists {
			// Encode the collection metadata to store it in the database.
			var data object.Object
			if data, err = object.MarshalSystem(collection); err != nil {
				return err
			}

			if err = s.db.Put(collectionKey, data); err != nil {
				return err
			}
		}
	}

	return nil
}
