package leveldb

import (
	"bytes"

	honuiter "github.com/rotationalio/honu/iterator"
	pb "github.com/rotationalio/honu/object"
	opts "github.com/rotationalio/honu/options"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"google.golang.org/protobuf/proto"
)

// NewLevelDBIterator creates a new iterator that wraps a leveldb Iterator with object
// management access and Honu-specific serialization.
func NewLevelDBIterator(iter iterator.Iterator, options *opts.Options) honuiter.Iterator {
	return &ldbIterator{ldb: iter, options: options}
}

// Wraps the underlying leveldb iterator to provide object management access.
type ldbIterator struct {
	ldb     iterator.Iterator
	options *opts.Options
}

// Type check for the ldbIterator
var _ honuiter.Iterator = &ldbIterator{}

func (i *ldbIterator) Error() error { return i.ldb.Error() }
func (i *ldbIterator) Release()     { i.ldb.Release() }

func (i *ldbIterator) Next() bool {
	if ok := i.ldb.Next(); !ok {
		return false
	}

	// If we aren't including Tombstones, we need to check if the next version is a
	// tombstone before we know if we have a next value or not.
	if !i.options.Tombstones {
		if obj, err := i.Object(); err != nil || obj.Tombstone() {
			return i.Next()
		}
	}
	return true
}

func (i *ldbIterator) Prev() bool {
	if ok := i.ldb.Prev(); !ok {
		return false
	}

	// If we aren't including Tombstones, we need to check if the next version is a
	// tombstone before we know if we have a next value or not.
	if !i.options.Tombstones {
		if obj, err := i.Object(); err != nil || obj.Tombstone() {
			return i.Prev()
		}
	}
	return true
}

func (i *ldbIterator) Seek(key []byte) bool {
	// We need to prefix the seek with the correct namespace
	if i.options.Namespace != "" {
		key = prepend(i.options.Namespace, key)
	}

	if ok := i.ldb.Seek(key); !ok {
		return false
	}

	// If we aren't including Tombstones, we need to check if the check if the current
	// version is a tombstone, and if not, continue to the next non-tombstone object
	if !i.options.Tombstones {
		if obj, err := i.Object(); err != nil || obj.Tombstone() {
			return i.Next()
		}
	}
	return true
}

func (i *ldbIterator) Key() []byte {
	// Fetch the key then split the namespace from the key
	// Note that because the namespace itself might have colons in it, we
	// strip off the namespace prefix then remove any preceding colons.
	key := i.ldb.Key()
	if i.options.Namespace != "" {
		prefix := prepend(i.options.Namespace, nil)
		return bytes.TrimPrefix(key, prefix)
	}
	return key
}

func (i *ldbIterator) Value() []byte {
	obj, err := i.Object()
	if err != nil {
		// NOTE: if err is not nil, it's up to the caller to get the error from Object
		return nil
	}
	return obj.Data
}

func (i *ldbIterator) Object() (obj *pb.Object, err error) {
	obj = new(pb.Object)
	if err = proto.Unmarshal(i.ldb.Value(), obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (i *ldbIterator) Namespace() string {
	return i.options.Namespace
}
