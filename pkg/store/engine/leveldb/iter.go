package leveldb

import (
	"github.com/rotationalio/honu/pkg/store/iterator"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/rotationalio/honu/pkg/store/object"
	"github.com/rotationalio/honu/pkg/store/opts"

	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
)

func NewIterator(iter ldbiter.Iterator, ro *opts.ReadOptions) iterator.Iterator {
	return &ldbIterator{ldb: iter, ro: ro}
}

type ldbIterator struct {
	ldb ldbiter.Iterator
	ro  *opts.ReadOptions
}

// Iterator management and error handling.
func (i *ldbIterator) Error() error { return i.ldb.Error() }
func (i *ldbIterator) Release()     { i.ldb.Release() }

// Iterator accessor methods.
func (i *ldbIterator) Key() key.Key          { return key.Key(i.ldb.Key()) }
func (i *ldbIterator) Object() object.Object { return object.Object(i.ldb.Value()) }

//===========================================================================
// Implement iterator.Seeker interface
//===========================================================================

func (i *ldbIterator) Seek(key []byte) bool {
	// Attempt to seek to the specified key.
	if ok := i.ldb.Seek(key); !ok {
		return false
	}

	// If we aren't including tombstones, skip over any tombstones.
	if !i.ro.Tombstones {
		if i.Object().Tombstone() {
			return i.Next()
		}
	}
	return true
}

func (i *ldbIterator) Next() bool {
	if ok := i.ldb.Next(); !ok {
		return false
	}

	// Keep skipping over tombstones until we find a non-tombstone object by
	// recursively calling Next()
	if !i.ro.Tombstones {
		if i.Object().Tombstone() {
			return i.Next()
		}
	}

	return true
}

func (i *ldbIterator) Prev() bool {
	if ok := i.ldb.Prev(); !ok {
		return false
	}

	// Keep skipping over tombstones until we find a non-tombstone object by
	// recursively calling Prev()
	if !i.ro.Tombstones {
		if i.Object().Tombstone() {
			return i.Prev()
		}
	}

	return true
}

func (i *ldbIterator) First() bool {
	if ok := i.ldb.First(); !ok {
		return false
	}

	// If we aren't including tombstones, skip over any tombstones by calling Next
	if !i.ro.Tombstones {
		if i.Object().Tombstone() {
			return false
		}
	}

	return true
}

func (i *ldbIterator) Last() bool {
	if ok := i.ldb.Last(); !ok {
		return false
	}

	// If we aren't including tombstones, skip over any tombstones by calling Prev
	if !i.ro.Tombstones {
		if i.Object().Tombstone() {
			return i.Prev()
		}
	}

	return true
}
