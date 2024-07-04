/*
Package iterator provides an interface and implementations to traverse over the contents
of an embedded database while maintaining and reading replicated object metadata.
*/
package iterator

import (
	"errors"

	pb "github.com/rotationalio/honu/object/v1"
)

// Standard iterator errors that may be returned for error type checking.
var (
	ErrIterReleased = errors.New("iterator has been released")
)

// Iterator retrieves multiple results from the underlying database, allowing users to
// loop over the results one at a time in a memory-safe fashion. The iterator may wrap
// a leveldb iterator or a sqlite rows context, fetching one row at a time in a Next
// loop. The Iterator also provides access to the versioned metadata for low-level
// interactions with the replicated data types.
// TODO: Implement IteratorSeeker interface from leveldb
type Iterator interface {
	// Next moves the iterator to the next key/value pair or row.
	// It returns false if the iterator has been exhausted.
	Next() bool

	// Prev moves the iterator to the previous key/value pair or row.
	// It returns false if the iterator has been exhausted.
	Prev() bool

	// Error returns any accumulated error. Exhausting all rows or key/value pairs is
	// not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, the object key of a row, or
	// nil if done. The caller should not modify the contents of the returned slice, and
	// its contents may change as the iterator progresses across the database.
	Key() []byte

	// Value returns the data of the current key/value pair, the object data of a row,
	// or nil if done. The caller should not modify the contents of the returned slice,
	// and its contents may change as the iterator progresses across the database.
	Value() []byte

	// Object returns the replicated object metadata and version information without
	// data. This method can be used to read meta-information and is also used for
	// replication. The object's Data property needs to be populated with Value() after
	// the object has been loaded from disk.
	Object() (*pb.Object, error)

	// When called, Release will close and release any resources associated with the
	// iterator. Release can be called multiple times without error but after it has
	// been called, no Iterator methods will return data.
	Release()

	// Seek moves the iterator to the first key/value pair whose key is greater than or
	// equal to the given key. It returns whether such pair exists.
	Seek(key []byte) bool

	// Namespace returns the current namespace the iterator is operating on.
	Namespace() string
}
