package iterator

import (
	pb "github.com/rotationalio/honu/object/v1"
	"github.com/rotationalio/honu/options"
)

// NewEmptyIterator creates an empty iterator that returns nothing. The err parameter
// can be nil, but if not nil the given err will be returned by the Error method.
func NewEmptyIterator(err error, namespace string) Iterator {
	if namespace == "" {
		namespace = options.NamespaceDefault
	}
	return &emptyIterator{err: err, namespace: namespace}
}

type emptyIterator struct {
	released  bool
	err       error
	namespace string
}

var _ Iterator = &emptyIterator{}

func (i *emptyIterator) rErr() {
	if i.err == nil && i.released {
		i.err = ErrIterReleased
	}
}

func (i *emptyIterator) Next() bool                { i.rErr(); return false }
func (i *emptyIterator) Prev() bool                { i.rErr(); return false }
func (i *emptyIterator) Seek(key []byte) bool      { i.rErr(); return false }
func (*emptyIterator) Key() []byte                 { return nil }
func (*emptyIterator) Value() []byte               { return nil }
func (*emptyIterator) Object() (*pb.Object, error) { return nil, nil }
func (i *emptyIterator) Error() error              { return i.err }
func (i *emptyIterator) Release()                  { i.released = true }
func (i *emptyIterator) Namespace() string         { return i.namespace }
