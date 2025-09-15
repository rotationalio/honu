package iterator

import (
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/object"
)

// Iterator retrieves multiple results from the underlying database, allowing users to
// loop over the results one at a time in a memory-safe fashion.
type Iterator interface {
	Seeker
	Releaser

	// Key returns the key of the current key/value pair, the object key of a row, or
	// nil if done. The caller should not modify the contents of the returned slice, and
	// its contents may change as the iterator progresses across the database.
	Key() key.Key

	// Object returns the replicated object metadata and version information without
	// data. This method can be used to read meta-information and is also used for
	// replication. The object's Data property needs to be populated with Value() after
	// the object has been loaded from disk.
	Object() object.Object

	// Error returns any accumulated errors. Exhausting all rows or key/value pairs is
	// not considered to be an error.
	Error() error
}

type Seeker interface {
	// Seek moves the iterator to the first key/value pair whose key is greater than or
	// equal to the given key. It returns whether such pair exists.
	Seek(key []byte) bool

	// Next moves the iterator to the next key/value pair or row.
	// It returns false if the iterator has been exhausted.
	Next() bool

	// Prev moves the iterator to the previous key/value pair or row.
	// It returns false if the iterator has been exhausted.
	Prev() bool

	// First moves the iterator to the first key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	First() bool

	// Last moves the iterator to the last key/value pair. If the iterator
	// only contains one key/value pair then First and Last would moves
	// to the same key/value pair.
	// It returns whether such pair exist.
	Last() bool
}

type Releaser interface {
	// When called, Release will close and release any resources associated with the
	// iterator. Release can be called multiple times without error but after it has
	// been called, no Iterator methods will return data.
	Release()
}
