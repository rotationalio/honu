package store

import (
	"bytes"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
	"go.rtnl.ai/honu/pkg/config"
	"go.rtnl.ai/honu/pkg/errors"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
	"go.rtnl.ai/ulid"
)

// System collections, indexes, and users are used for internal database management and
// are directly identified by a unique ID that is less than any ULID that would be
// generated since they are before March 20, 1984 (and all ULIDs should be generated
// after September, 2025). These IDs are hardcoded to ensure that they are always
// available and do not change between database instances.
var (
	SystemPrefix          = [5]byte{0x00, 0x68, 0x6f, 0x6e, 0x75}
	SystemHonuAgent       = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x68, 0x6f, 0x6e, 0x75, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x00})
	SystemCollections     = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e})
	SystemReplicas        = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x01, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6c, 0x69, 0x73, 0x74})
	SystemAccessControl   = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x02, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67})
	SystemCollectionNames = ulid.ULID([16]byte{0x00, 0x68, 0x6f, 0x6e, 0x75, 0x00, 0x63, 0x6f, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x69, 0x64, 0x78})
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
	conf config.StoreConfig
	pid  lamport.PID
	db   *bbolt.DB
}

// Open a new Store with the provided configuration. Only one Store can be opened for a
// single configuration (database path), which locks the files on disk to prevent other
// processes from opening the same database. If the database already exists, it will be
// opened and the existing data will be used. If the database does not exist, it will
// be created. The Store can then be used by multiple goroutines, to sequence access
// to disk.
func Open(conf config.Config) (s *Store, err error) {
	s = &Store{
		conf: conf.Store,
		pid:  lamport.PID(conf.PID),
	}

	// TODO: better open with options.
	if s.db, err = bbolt.Open(conf.Store.DataPath, 0600, nil); err != nil {
		return nil, err
	}

	// Ensure the database is initialized and ready for use.
	if err = s.initialize(); err != nil {
		s.db.Close()
		return nil, err
	}

	// Check that the database in in a ready state.
	if err = s.check(); err != nil {
		s.db.Close()
		return nil, err
	}

	return s, nil
}

// Close the store and release all resources associated with it.
func (s *Store) Close() error {
	err := s.db.Close()
	s.db = nil
	return err
}

//===========================================================================
// Transactions Handling
//===========================================================================

// Begin starts a new transaction. Multiple read-only transactions can be used
// concurrently but only one write transaction can be used at a time. Starting multiple
// write transactions will cause the calls to block and be serialized until the current
// write transaction finishes.
//
// Transactions must be either committed or rolled back when they are no longer needed
// to release the associated resources. If a transaction is not committed or rolled
// back, pages in the database will not be freed and other transactions may be remain
// deadlocked.
func (s *Store) Begin(opts *TxOptions) (tx *Tx, err error) {
	if s.db == nil {
		return nil, errors.ErrClosed
	}

	if opts == nil {
		opts = &TxOptions{}
	}

	tx = &Tx{
		opts: opts,
	}

	if tx.tx, err = s.db.Begin(!opts.ReadOnly); err != nil {
		return nil, err
	}

	return tx, nil
}

//===========================================================================
// Collection Management
//===========================================================================

// Returns a list of the collections that the store is maintaining excluding any system
// collections that are used for internal database management. Only the metadata is
// returned for each collection, so any collection operations must be performed after
// opening the collection by its ID or name.
func (s *Store) Collections() (collections []metadata.Collection, err error) {
	var tx *bbolt.Tx
	if tx, err = s.db.Begin(false); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var bucket *bbolt.Bucket
	if bucket = tx.Bucket(SystemCollections[:]); bucket == nil {
		return nil, errors.New("collections: system collections bucket does not exist")
	}

	bucket.ForEach(func(key, data []byte) error {
		// Skip system collections.
		if bytes.HasPrefix(key, SystemPrefix[:]) {
			return nil
		}

		var c metadata.Collection
		if err := lani.Unmarshal(data, &c); err != nil {
			return err
		}
		collections = append(collections, c)
		return nil
	})

	return collections, nil
}

// Creates a new collection with the given name and associates it with a unique ID.
// The name is case-insensitive and should be unique within the store. It must contain
// only no spaces or punctuation and cannot start with a number. The name also must not
// be a ULID string, which is reserved for collection IDs.
func (s *Store) New(info *metadata.Collection) (err error) {
	// Validate the info to ensure a collection can be created.
	if err = info.Validate(); err != nil {
		return err
	}

	// A collection must not have an ID set.
	if !info.ID.IsZero() {
		return ErrCreateID
	}

	// TODO: lock the collection key
	// TODO: check name in index to ensure uniqueness.
	// TODO: normalize name to ensure it is valid and does not contain any illegal characters.
	// TODO: save the collection index to disk.
	return nil
}

// Has returns true if the collection with the specified ID or name exists in the store.
func (s *Store) Has(identifier any) (exists bool, err error) {
	var collectionID ulid.ULID
	switch v := identifier.(type) {
	case string:
		panic("not implemented yet")
	case ulid.ULID:
		collectionID = v
	default:
		// TODO: return a better error
		return false, errors.New("invalid collection identifier")
	}

	// TODO: handle versions
	key := key.New(collectionID, &lamport.Scalar{PID: 0, VID: 1})
	err = s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(SystemCollections[:])
		if b == nil {
			return nil
		}

		if v := b.Get(key); v != nil {
			exists = true
		}
		return nil
	})
	return exists, err
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
func (s *Store) DB() *bbolt.DB {
	return s.db
}

//===========================================================================
// Initialization
//===========================================================================

func defaultCollections() []*metadata.Collection {
	now := time.Now()
	return []*metadata.Collection{
		{
			ID:   SystemCollections,
			Name: "honu_collections",
			Version: &metadata.Version{
				Scalar:  lamport.Scalar{PID: 0, VID: 1}, // This is the default first version for system collections.
				Created: now,
			},
			Owner: SystemHonuAgent,
			Group: SystemHonuAgent,
			Indexes: []*metadata.Index{
				{
					ID:   SystemCollectionNames,
					Name: "honu_collections_unique_name",
					Type: metadata.UNIQUE,
					Field: &metadata.Field{
						Name:       "name",
						Type:       metadata.StringField,
						Collection: SystemCollections,
					},
					Ref: &metadata.Field{
						Name:       "id",
						Type:       metadata.ULIDField,
						Collection: SystemCollections,
					},
				},
			},
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
}

func (s *Store) initialize() (err error) {
	if s.conf.ReadOnly {
		// If the store is read-only, do not attempt to initialize it.
		return nil
	}

	// Ensures that the system collections are created and initialized in the store and
	// that the system user is created with the correct permissions.
	defaultCollections := defaultCollections()

	var tx *bbolt.Tx
	if tx, err = s.db.Begin(true); err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the collections bucket to create the systems collection info in.
	var collectionsBucket *bbolt.Bucket
	if collectionsBucket, err = tx.CreateBucketIfNotExists(SystemCollections[:]); err != nil {
		return fmt.Errorf("could not create system collections bucket: %w", err)
	}

	for _, collection := range defaultCollections {
		// System collections are not updated except between versions of honu, so they
		// can be fetched directly with the hardcoded ID.
		collectionKey := key.New(collection.ID, &collection.Version.Scalar)

		var exists bool
		if collectionMeta := collectionsBucket.Get(collectionKey); collectionMeta != nil {
			exists = true
		}

		if !exists {
			// Encode the collection metadata to store it in the database.
			var data object.Object
			if data, err = object.MarshalSystem(collection); err != nil {
				return fmt.Errorf("could not marshal collection metadata %s: %w", collection.Name, err)
			}

			if err = collectionsBucket.Put(collectionKey, data); err != nil {
				return fmt.Errorf("could not store collection metadata %s: %w", collection.Name, err)
			}
		}

		// Create the bucket for the collection itself to hold its objects.
		var cbckt *bbolt.Bucket
		if cbckt, err = tx.CreateBucketIfNotExists(collection.ID[:]); err != nil {
			return fmt.Errorf("could not create collection bucket %s: %w", collection.Name, err)
		}

		// If indexes are defined on the collection, ensure they are created.
		for _, idx := range collection.Indexes {
			if _, err = cbckt.CreateBucketIfNotExists(idx.ID[:]); err != nil {
				return fmt.Errorf("could not create index %s in %s: %w", idx.Name, collection.Name, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit initialize transaction: %w", err)
	}
	return nil
}

func (s *Store) check() (err error) {
	// Ensure that the system collections and indexes exist in the database.
	defaultCollections := defaultCollections()

	var tx *bbolt.Tx
	if tx, err = s.db.Begin(false); err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the collections bucket to check the systems collection info in.
	var collectionsBucket *bbolt.Bucket
	if collectionsBucket = tx.Bucket(SystemCollections[:]); collectionsBucket == nil {
		return fmt.Errorf("missing bucket %s for system collections", SystemCollections)
	}

	for _, collection := range defaultCollections {
		// System collections are not updated except between versions of honu, so they
		// can be fetched directly with the hardcoded ID.
		collectionKey := key.New(collection.ID, &collection.Version.Scalar)

		if meta := collectionsBucket.Get(collectionKey); meta == nil {
			err = errors.Join(err, fmt.Errorf("missing metadata for collection %s (%s)", collection.Name, collection.ID))
			continue
		}

		// Ensure the bucket for the collection itself exists to hold its objects.
		var cbckt *bbolt.Bucket
		if cbckt = tx.Bucket(collection.ID[:]); cbckt == nil {
			err = errors.Join(err, fmt.Errorf("missing bucket for collection %s (%s)", collection.Name, collection.ID))
			continue
		}

		// If indexes are defined on the collection, ensure they are created.
		for _, idx := range collection.Indexes {
			if ibckt := cbckt.Bucket(idx.ID[:]); ibckt == nil {
				err = errors.Join(err, fmt.Errorf("missing index bucket %s in %s", idx.Name, collection.Name))
			}
		}
	}

	return err
}
