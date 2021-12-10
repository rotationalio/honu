package leveldb

import (
	"bytes"

	honuiter "github.com/rotationalio/honu/iterator"
	pb "github.com/rotationalio/honu/object"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"google.golang.org/protobuf/proto"
)

// NewLevelDBIterator creates a new iterator that wraps a leveldb Iterator with object
// management access and Honu-specific serialization.
func NewLevelDBIterator(iter iterator.Iterator, namespace string) honuiter.Iterator {
	return &ldbIterator{ldb: iter, namespace: namespace}
}

// Wraps the underlying leveldb iterator to provide object management access.
type ldbIterator struct {
	ldb       iterator.Iterator
	namespace string
}

// Type check for the ldbIterator
var _ honuiter.Iterator = &ldbIterator{}

func (i *ldbIterator) Next() bool   { return i.ldb.Next() }
func (i *ldbIterator) Prev() bool   { return i.ldb.Prev() }
func (i *ldbIterator) Error() error { return i.ldb.Error() }
func (i *ldbIterator) Release()     { i.ldb.Release() }

func (i *ldbIterator) Seek(key []byte) bool {
	// We need to prefix the seek with the correct namespace
	key = prepend(i.namespace, key)
	return i.ldb.Seek(key)
}

func (i *ldbIterator) Key() []byte {
	// Fetch the key then split the namespace from the key
	key := i.ldb.Key()
	parts := bytes.SplitN(key, nssep, 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return key
}

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

func (i *ldbIterator) Namespace() string {
	return i.namespace
}
