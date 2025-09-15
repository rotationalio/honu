package store

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
	berrors "go.etcd.io/bbolt/errors"
	"go.rtnl.ai/honu/pkg/errors"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
	"go.rtnl.ai/ulid"
)

type Tx struct {
	tx          *bbolt.Tx
	opts        *TxOptions
	closed      bool
	rollbackErr error
	commitErr   error

	// Collections metadata bucket and names index.
	cmbkt   *bbolt.Bucket
	cmnames *bbolt.Bucket

	// Cache of opened buckets for collections.
	collections map[ulid.ULID]*Collection
}

type TxOptions struct {
	ReadOnly    bool
	ClosedError bool
}

// Commits the transaction if it is writeable. If the transaction is read-only, then
// this is a no-op and returns nil (unlike bolt which will return an error). Commit
// can also be called multiple times safely without an error being returned.
func (t *Tx) Commit() error {
	if t.writeable() {
		t.commitErr = t.tx.Commit()
		t.closed = true

		if !t.opts.ClosedError && errors.Is(t.commitErr, berrors.ErrTxClosed) {
			t.commitErr = nil
		}

		// Clear the cached collections after a commit.
		t.collections = nil
		t.opts = nil
	}
	return t.commitErr
}

// Rollback the transaction if it is still open. If the transaction has already been
// committed or rolled back, then this is a no-op and returns nil (unlike bolt which
// will return an error). Rollback can also be called multiple times safely without
// an error being returned.
func (t *Tx) Rollback() error {
	if !t.closed {
		t.rollbackErr = t.tx.Rollback()
		t.closed = true

		if !t.opts.ClosedError && errors.Is(t.rollbackErr, berrors.ErrTxClosed) {
			t.rollbackErr = nil
		}

		// Clear the cached collections after a rollback.
		t.collections = nil
		t.opts = nil
	}
	return t.rollbackErr
}

// Collection retrieves a bucket by name if a string is supplied by looking the name up
// in the collections name index or by ID if a ulid.ULID. If the collection does not
// exist an error is returned.
func (t *Tx) Collection(identifier any) (c *Collection, err error) {
	if t.closed {
		return nil, errors.ErrTxClosed
	}

	// Initialize the collections management system if it hasn't already been.
	if err = t.initialize(); err != nil {
		return nil, err
	}

	var (
		collectionID   ulid.ULID
		collectionName string
	)

	// Returns an error if either the ID or name is zero valued, the name is not a
	// valid collection name, or the ID is a system collection ID.
	if collectionID, collectionName, err = collectionIdentifier(identifier); err != nil {
		return nil, err
	}

	// If a name was specified, look up the ID in the name index.
	if collectionID.IsZero() {
		if v := t.cmnames.Get([]byte(collectionName)); v != nil {
			copy(collectionID[:], v)
		} else {
			return nil, errors.ErrNoCollection
		}
	}

	// Check the cache of opened collections first.
	if c, ok := t.collections[collectionID]; ok {
		return c, nil
	}

	// Initialize the collection and cache it
	c = &Collection{
		bkt: t.tx.Bucket(collectionID[:]),
	}

	if c.bkt == nil {
		log.Error().Str("collectionID", collectionID.String()).Msg("collection bucket does not exist")
		return nil, errors.ErrRepairCollection
	}

	// Get the latest metadata for the collection.
	cursor := t.cmbkt.Cursor()
	key, meta := cursor.Seek(collectionID[:])
	if key != nil && bytes.HasPrefix(key, collectionID[:]) {
		c.Collection = metadata.Collection{}
		if err = object.UnmarshalSystem(object.Object(meta), &c.Collection); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal collection metadata")
			return nil, errors.ErrRepairCollection
		}
	} else {
		log.Error().Str("collectionID", collectionID.String()).Msg("collection metadata does not exist")
		return nil, errors.ErrRepairCollection
	}

	t.collections[collectionID] = c
	return c, nil

}

// Has returns true if the object with the specified ID has any version (including
// tombstones) stored in the specified collection. See Exists() for checking if the
// latest version of the object is not a tombstone.
func (t *Tx) Has(collection ulid.ULID, id ulid.ULID) (exists bool, err error) {
	if t.closed {
		return false, errors.ErrTxClosed
	}
	panic("not implemented yet")
}

// Exists returns true if the object with the specified ID exists in the specified
// collection and the latest version is not a tombstone.
func (t *Tx) Exists(collection ulid.ULID, id ulid.ULID) (exists bool, err error) {
	if t.closed {
		return false, errors.ErrTxClosed
	}
	panic("not implemented yet")
}

// writeable returns true if the transaction is not read-only and has not been closed.
func (t *Tx) writeable() bool {
	return !t.opts.ReadOnly && !t.closed
}

// initialize the collection management system on the transaction if it hasn't already
// been initialized (including the collections bucket, names index, and cache)
func (t *Tx) initialize() error {
	if t.collections == nil {
		t.collections = make(map[ulid.ULID]*Collection)
	}

	if t.cmbkt == nil {
		t.cmbkt = t.tx.Bucket(SystemCollections[:])
		if t.cmbkt == nil {
			log.Error().Msg("system collections bucket does not exist")
			return errors.ErrNotInitialized
		}
	}

	if t.cmnames == nil {
		t.cmnames = t.cmbkt.Bucket(SystemCollectionNames[:])
		if t.cmnames == nil {
			log.Error().Msg("system collection names index bucket does not exist")
			return errors.ErrNotInitialized
		}
	}

	return nil
}
