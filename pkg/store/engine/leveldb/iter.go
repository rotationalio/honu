package leveldb

import (
	"go.rtnl.ai/honu/pkg/store/iterator"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/object"

	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
)

func NewIterator(iter ldbiter.Iterator) iterator.Iterator {
	return &ldbIterator{Iterator: iter}
}

type ldbIterator struct {
	ldbiter.Iterator
}

// Iterator accessor methods.
func (i *ldbIterator) Key() key.Key          { return key.Key(i.Iterator.Key()) }
func (i *ldbIterator) Object() object.Object { return object.Object(i.Iterator.Value()) }
func (i *ldbIterator) Error() error          { return Wrap(i.Iterator.Error()) }
