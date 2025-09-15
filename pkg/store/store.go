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
	conf   config.StoreConfig
	region string
	pid    lamport.PID
	db     *bbolt.DB
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

	// ReadOnly checks
	if s.conf.ReadOnly && !opts.ReadOnly {
		return nil, errors.ErrReadOnlyDB
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
// TODO: check permissions and ACLs to ensure the user is allowed to read collections.
func (s *Store) Collections() (collections []*metadata.Collection, err error) {
	var tx *bbolt.Tx
	if tx, err = s.db.Begin(false); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var bucket *bbolt.Bucket
	if bucket = tx.Bucket(SystemCollections[:]); bucket == nil {
		return nil, errors.ErrNotInitialized
	}

	bucket.ForEach(func(key, data []byte) error {
		// Skip system collections.
		if bytes.HasPrefix(key, SystemPrefix[:]) {
			return nil
		}

		var c *metadata.Collection
		if err := lani.Unmarshal(data, c); err != nil {
			return err
		}
		collections = append(collections, c)
		return nil
	})

	return collections, nil
}

// Creates a new collection with the given name and associates it with a unique ID.
// The name is case-insensitive and should be unique within the store. It must contain
// only no spaces or punctuation and cannot start with a number.
//
// This method will set the collection version, ID, creation, and modification time;
// any data already set in the collection info will be overridden without error.
//
// Any indexes defined on the collection will be created when the collection is created.
// TODO: check permissions and ACLs to ensure the user is allowed to create the collection.
func (s *Store) New(info *metadata.Collection) (err error) {
	// Readonly checks
	if s.conf.ReadOnly {
		return errors.ErrReadOnlyDB
	}

	// Validate the info to ensure a collection can be created.
	if err = info.Validate(); err != nil {
		return err
	}

	// A collection must not have an ID set.
	if !info.ID.IsZero() {
		return errors.ErrCreateID
	}

	// Update the collection info to set the ID and creation time.
	info.ID = ulid.MakeSecure()
	info.Version = &metadata.Version{
		Scalar:    s.pid.Next(nil),
		Region:    s.region,
		Parent:    nil,
		Tombstone: false,
		Created:   time.Now(),
	}
	info.Created = info.Version.Created
	info.Modified = info.Version.Created

	var tx *bbolt.Tx
	if tx, err = s.db.Begin(true); err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the collections bucket to create the new collection info in.
	// No nil checks are required because if not initialized correctly a nil panic will occur.
	collections := tx.Bucket(SystemCollections[:])
	nameIndex := collections.Bucket(SystemCollectionNames[:])

	// Ensure the collection name is unique by checking the name index.
	// NOTE: the collection ID should be unique by ULID generation.
	if v := nameIndex.Get([]byte(info.Name)); v != nil {
		return errors.ErrCollectionExists
	}

	// Add name to the name index.
	if err = nameIndex.Put([]byte(info.Name), info.ID[:]); err != nil {
		return fmt.Errorf("could not index collection name %s: %w", info.Name, err)
	}

	// Encode the collection metadata to store in the database.
	var data object.Object
	if data, err = object.MarshalSystem(info); err != nil {
		return fmt.Errorf("could not marshal collection metadata %s: %w", info.Name, err)
	}

	// Store the collection metadata in the collections bucket.
	key := key.New(info.ID, &info.Version.Scalar)
	if err = collections.Put(key, data); err != nil {
		return fmt.Errorf("could not store collection metadata %s: %w", info.Name, err)
	}

	// Create the bucket for the collection itself to hold its objects.
	var bucket *bbolt.Bucket
	if bucket, err = tx.CreateBucketIfNotExists(info.ID[:]); err != nil {
		return fmt.Errorf("could not create collection bucket %s: %w", info.Name, err)
	}

	// If indexes are defined on the collection, ensure they are created.
	for _, idx := range info.Indexes {
		if _, err = bucket.CreateBucket(idx.ID[:]); err != nil {
			return fmt.Errorf("could not create index %s in %s: %w", idx.Name, info.Name, err)
		}
	}

	return tx.Commit()
}

// Has returns true if the collection with the specified ID or name exists in the store.
// TODO: check permissions and ACLs to ensure the user is allowed to read the collection.
func (s *Store) Has(identifier any) (exists bool, err error) {
	var (
		collectionID   ulid.ULID
		collectionName string
	)

	// Returns an error if either the ID or name is zero valued, the name is not a
	// valid collection name, or the ID is a system collection ID.
	if collectionID, collectionName, err = collectionIdentifier(identifier); err != nil {
		return false, err
	}

	var tx *bbolt.Tx
	if tx, err = s.db.Begin(false); err != nil {
		return false, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	collections := tx.Bucket(SystemCollections[:])

	// If a name was provided, look up the ID first.
	if collectionID.IsZero() {
		nameIndex := collections.Bucket(SystemCollectionNames[:])
		if v := nameIndex.Get([]byte(collectionName)); v != nil {
			copy(collectionID[:], v)
		} else {
			return false, nil
		}
	}

	// collectionID is the prefix, we need to look to see if any objects start with
	// that prefix, because that will indicate that there is at least one version of
	// the collection stored in the database.
	cursor := collections.Cursor()
	key, _ := cursor.Seek(collectionID[:])
	exists = key != nil && bytes.HasPrefix(key, collectionID[:])

	return exists, err
}

// Collection returns the latest metadata version fo the specified collection using its
// identifier (e.g. either the collection ID or name). If the collection does not exist,
// an ErrNoCollection error is returned.
// TODO: check permissions and ACLs to ensure the user is allowed to read the collection.
func (s *Store) Collection(identifier any) (info *metadata.Collection, err error) {
	return info, nil
}

// Modifies the metadata of an existing collection; the collection should either have
// an ID or a name for reference and must already exist in the store.
func (s *Store) Modify(info *metadata.Collection) error {
	return errors.ErrNotImplemented
}

// Drop a collection, removing it from the store and deleting all of its contained
// objects. The collection can be recreated later with the same name, but all of the
// previous objects and their versions will be deleted.
//
// If a collection does not exist an ErrNoCollection error is returned.
//
// This operation is handled by deleting the collection bucket and any indexes from
// BoltDB then adding a tombstone record to the collection to replicate the deletion
// to all replicas and deletes all prior collection metadata versions to free up space.
//
// Use this operation to free up space in the database, but note that data deletion,
// history, and metadata will be lost. If you want to keep the version history of
// the objects in the collection, use the Collection.Empty method instead which adds
// a tombstone version to each object in the collection but keeps the version history
// intact.
//
// Edge case: there is a case where a collection is dropped concurrently with a
// collection modification and the modification wins the last writer wins policy. In
// this case, the collection should eventually be recreated on a subsequent replication.
// To avoid this, ensure that the drop operation happens in quorum or on the highest
// priority replica.
//
// TODO: check permissions and ACLs to ensure the user is allowed to drop the collection.
func (s *Store) Drop(identifier any) (err error) {
	var (
		collectionID   ulid.ULID
		collectionName string
	)

	// Returns an error if either the ID or name is zero valued, the name is not a
	// valid collection name, or the ID is a system collection ID.
	if collectionID, collectionName, err = collectionIdentifier(identifier); err != nil {
		return err
	}

	var tx *bbolt.Tx
	if tx, err = s.db.Begin(true); err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	collections := tx.Bucket(SystemCollections[:])

	// If a name was provided, look up the ID first.
	if collectionID.IsZero() {
		nameIndex := collections.Bucket(SystemCollectionNames[:])
		if v := nameIndex.Get([]byte(collectionName)); v != nil {
			copy(collectionID[:], v)
		} else {
			return errors.ErrNoCollection
		}
	}

	// Fetch the collection meta to get its current version.
	var meta metadata.Collection
	cursor := collections.Cursor()
	if key, data := cursor.Seek(collectionID[:]); key == nil || !bytes.HasPrefix(key, collectionID[:]) {
		return errors.ErrNoCollection
	} else {
		if err = lani.Unmarshal(data, &meta); err != nil {
			return fmt.Errorf("could not unmarshal collection meta: %w", err)
		}
	}

	// Remove the name from the name index.
	nameIndex := collections.Bucket(SystemCollectionNames[:])
	if err = nameIndex.Delete([]byte(meta.Name)); err != nil {
		return fmt.Errorf("could not remove collection name from index: %w", err)
	}

	// Delete all collection versions.
	for key, _ := cursor.Seek(collectionID[:]); key != nil && bytes.HasPrefix(key, collectionID[:]); key, _ = cursor.Next() {
		if err = collections.Delete(key); err != nil {
			return fmt.Errorf("could not delete collection version: %w", err)
		}
	}

	// Create the tombstone version for the collection to replicate the deletion.
	// Modify the current collection inline to prevent an allocation.
	meta.Tombstone(s.pid, s.region)

	var tdata object.Object
	if tdata, err = object.MarshalSystem(&meta); err != nil {
		return fmt.Errorf("could not marshal tombstone collection meta: %w", err)
	}

	tkey := key.New(meta.ID, &meta.Version.Scalar)
	if err = collections.Put(tkey, tdata); err != nil {
		return fmt.Errorf("could not store tombstone collection meta: %w", err)
	}

	// Delete the collection bucket to remove all of its objects and indexes.
	// NOTE: DeleteBucket removes the bucket and all nested buckets (including indexes)
	// and marks the pages as free.
	if err = tx.DeleteBucket(collectionID[:]); err != nil {
		return fmt.Errorf("could not delete collection bucket: %w", err)
	}

	return tx.Commit()
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
	return errors.ErrNotImplemented
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
