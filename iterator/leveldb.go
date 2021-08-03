package iterator

import (
	pb "github.com/rotationalio/honu/object"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"google.golang.org/protobuf/proto"
)

// NewLevelDBIterator creates a new iterator that wraps a leveldb Iterator with object
// management access and Honu-specific serialization.
func NewLevelDBIterator(iter iterator.Iterator) Iterator {
	return &ldbIterator{ldb: iter}
}

// Wraps the underlying leveldb iterator to provide object management access.
type ldbIterator struct {
	ldb iterator.Iterator
}

// Type check for the ldbIterator
var _ Iterator = &ldbIterator{}

func (i *ldbIterator) Next() bool   { return i.ldb.Next() }
func (i *ldbIterator) Prev() bool   { return i.ldb.Prev() }
func (i *ldbIterator) Key() []byte  { return i.ldb.Key() }
func (i *ldbIterator) Error() error { return i.ldb.Error() }
func (i *ldbIterator) Release()     { i.ldb.Release() }

func (i *ldbIterator) Value() []byte {
	obj, _ := i.Object()
	if obj != nil {
		return obj.Data
	}
	return nil
}

func (i *ldbIterator) Object() (obj *pb.Object, err error) {
	obj = new(pb.Object)
	if err = proto.Unmarshal(i.ldb.Value(), obj); err != nil {
		return nil, err
	}
	return obj, nil
}
