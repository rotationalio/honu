package iterator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/errors"
	. "go.rtnl.ai/honu/pkg/store/iterator"
)

func TestEmptyIterator(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		// Check that the empty iterator returns expected values
		iter := Empty(nil)
		require.False(t, iter.Next())
		require.False(t, iter.Prev())
		require.False(t, iter.Seek([]byte("foo")))
		require.Nil(t, iter.Key())
		require.NoError(t, iter.Error())
		require.Nil(t, iter.Object())

		// After calling release the empty iterator should still have no error
		iter.Release()
		require.NoError(t, iter.Error())

		// However if next is called after release, then the iterator should error
		require.False(t, iter.Next())
		require.EqualError(t, iter.Error(), errors.ErrIterReleased.Error())
	})

	t.Run("Error", func(t *testing.T) {
		// Check that the empty iterator can be initialized with an error
		iter := Empty(errors.New("something bad happened"))
		require.EqualError(t, iter.Error(), "something bad happened")

		// Ensure that calling any of the iterator methods do not change the error
		require.False(t, iter.Next())
		require.False(t, iter.Prev())
		require.False(t, iter.Seek([]byte("foo")))
		require.Nil(t, iter.Key())
		require.Nil(t, iter.Object())

		require.EqualError(t, iter.Error(), "something bad happened")

		// Ensure calling Release doesn't change the error
		iter.Release()
		require.EqualError(t, iter.Error(), "something bad happened")

		// Ensure calling Next after Release doesn't change the error
		iter.Next()
		require.EqualError(t, iter.Error(), "something bad happened")
	})
}
