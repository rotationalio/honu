package store

import (
	"bytes"

	"go.etcd.io/bbolt"
	"go.rtnl.ai/honu/pkg/errors"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/ulid"
)

// Collections are subsets of the Store that allow access to related objects. Each
// object in a collection is prefixed by the collection ID, ensuring that the objects
// are grouped together and can be accessed efficiently.
type Collection struct {
	metadata.Collection
	bkt *bbolt.Bucket `json:"-" msg:"-"`
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

// Has returns true if the object with the specified ID has any version (including
// tombstones) stored in the collection. See Exists() for checking if the latest version
// of the object is not a tombstone.
func (c *Collection) Has(id ulid.ULID) (exists bool, err error) {
	panic("not implemented yet")
}

// Exists returns true if the object with the specified ID exists in the collection
// and the latest version is not a tombstone.
func (c *Collection) Exists(id ulid.ULID) (exists bool, err error) {
	panic("not implemented yet")
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

//===========================================================================
// Collection Helper Methods
//===========================================================================

// Returns either an ULID or a name from the specified identifier, returning an error
// if the identifier is not valid (e.g. zero valued or not a collection name).
// NOTE: this method will not return a system collection ID or name.
func collectionIdentifier(identifier any) (id ulid.ULID, name string, err error) {
	switch v := identifier.(type) {
	case string:
		name = v
	case ulid.ULID:
		id = v
	default:
		return ulid.Zero, "", errors.ErrCollectionIdentifier
	}

	if id.IsZero() {
		if name == "" {
			return ulid.Zero, "", errors.ErrCollectionIdentifier
		}

		if err = metadata.ValidateName(name); err != nil {
			return ulid.Zero, "", err
		}
	} else {
		if bytes.HasPrefix(id[:], SystemPrefix[:]) {
			return ulid.Zero, "", errors.ErrCollectionIdentifier
		}
	}

	return id, name, nil
}
