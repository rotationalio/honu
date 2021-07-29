package iterator

import pb "github.com/rotationalio/honu/proto/v1"

// Wraps the underlying leveldb iterator to provide object management access.
type ldbIterator struct{}

// Type check for the ldbIterator
var _ Iterator = &ldbIterator{}

func (i *ldbIterator) Next() bool                { return false }
func (i *ldbIterator) Prev() bool                { return false }
func (*ldbIterator) Key() []byte                 { return nil }
func (*ldbIterator) Value() []byte               { return nil }
func (*ldbIterator) Object() (*pb.Object, error) { return nil, nil }
func (i *ldbIterator) Error() error              { return nil }
func (i *ldbIterator) Release()                  {}
