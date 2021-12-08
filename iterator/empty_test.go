package iterator_test

import (
	"errors"
	"testing"

	. "github.com/rotationalio/honu/iterator"
	"github.com/stretchr/testify/require"
)

func TestEmptyIterator(t *testing.T) {
	// Check that the empty iterator returns expected values
	iter := NewEmptyIterator(nil)
	require.False(t, iter.Next())
	require.False(t, iter.Prev())
	require.False(t, iter.Seek([]byte("foo")))
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	obj, err := iter.Object()
	require.NoError(t, err)
	require.Nil(t, obj)

	// After calling release the empty iterator should still have no error
	iter.Release()
	require.NoError(t, iter.Error())

	// However if next is called after release, then the iterator should error
	require.False(t, iter.Next())
	require.EqualError(t, iter.Error(), ErrIterReleased.Error())

	// Check that the empty iterator can be initialized with an error
	iter = NewEmptyIterator(errors.New("something bad happened"))
	require.EqualError(t, iter.Error(), "something bad happened")

	// Ensure that calling any of the iterator methods do not change the error
	require.False(t, iter.Next())
	require.False(t, iter.Prev())
	require.False(t, iter.Seek([]byte("foo")))
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())

	obj, err = iter.Object()
	require.NoError(t, err)
	require.Nil(t, obj)

	require.EqualError(t, iter.Error(), "something bad happened")

	// Ensure calling Release doesn't change the error
	iter.Release()
	require.EqualError(t, iter.Error(), "something bad happened")

	// Ensure calling Next after Release doesn't change the error
	iter.Next()
	require.EqualError(t, iter.Error(), "something bad happened")
}
