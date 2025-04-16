package iterator

import (
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/object"
)

// Empty creates an empty iterator that returns nothing. The err parameter
// can be nil, but if not nil the given err will be returned by the Error method.
func Empty(err error) Iterator {
	return &emptyIterator{err: err}
}

type emptyIterator struct {
	released bool
	err      error
}

var _ Iterator = &emptyIterator{}

func (i *emptyIterator) rErr() {
	if i.err == nil && i.released {
		i.err = ErrIterReleased
	}
}

func (i *emptyIterator) Next() bool           { i.rErr(); return false }
func (i *emptyIterator) Prev() bool           { i.rErr(); return false }
func (i *emptyIterator) First() bool          { i.rErr(); return false }
func (i *emptyIterator) Last() bool           { i.rErr(); return false }
func (i *emptyIterator) Seek(key []byte) bool { i.rErr(); return false }
func (*emptyIterator) Key() key.Key           { return nil }
func (*emptyIterator) Object() object.Object  { return nil }
func (i *emptyIterator) Error() error         { return i.err }
func (i *emptyIterator) Release()             { i.released = true }
