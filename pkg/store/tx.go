package store

import (
	"errors"

	"go.etcd.io/bbolt"
	berrors "go.etcd.io/bbolt/errors"
	"go.rtnl.ai/ulid"
)

type Tx struct {
	tx          *bbolt.Tx
	opts        *TxOptions
	closed      bool
	rollbackErr error
	commitErr   error
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
	}
	return t.rollbackErr
}

// Collection retrieves a bucket by name if a string is supplied by looking the name up
// in the collections name index or by ID if a ulid.ULID. If the collection does not
// exist an error is returned.
func (t *Tx) Collection(identifier any) (c *Collection, err error) {
	var bucket *bbolt.Bucket
	switch v := identifier.(type) {
	case string:
		panic("lookup by name not implemented yet")
	case ulid.ULID:
		bucket = t.tx.Bucket(v[:])
	default:
		// TODO: return a better error
		return nil, errors.New("invalid collection identifier")
	}

	if bucket == nil {
		// TODO: return a better error
		return nil, errors.New("collection not found")
	}

	c = &Collection{bck: bucket}
	return c, nil

}

// Has returns true if the object with the specified ID has any version (including
// tombstones) stored in the specified collection. See Exists() for checking if the
// latest version of the object is not a tombstone.
func (t *Tx) Has(collection ulid.ULID, id ulid.ULID) (exists bool, err error) {
	panic("not implemented yet")
}

// Exists returns true if the object with the specified ID exists in the specified
// collection and the latest version is not a tombstone.
func (t *Tx) Exists(collection ulid.ULID, id ulid.ULID) (exists bool, err error) {
	panic("not implemented yet")
}

// writeable returns true if the transaction is not read-only and has not been closed.
func (t *Tx) writeable() bool {
	return !t.opts.ReadOnly && !t.closed
}
