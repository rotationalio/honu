package leveldb_test

import (
	"math/rand/v2"
	"testing"

	"github.com/rotationalio/honu/pkg/store/iterator"
	"github.com/rotationalio/honu/pkg/store/key"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/ulid"
)

func TestIterator(t *testing.T) {
	// Create a database and populate it with some data
	db := OpenLevelDB(t, false)

	// Testing Data Structures
	cid := ulid.Make()
	oids := make(map[ulid.ULID]int)
	nkeys := 0

	// Create three different objects with multiple versions each in the same collection
	var target ulid.ULID
	for i := 0; i < 3; i++ {
		count := rand.IntN(20) + 4
		nkeys += count

		oid := ulid.Make()
		oids[oid] = count

		Populate(t, db, count, cid, oid)

		if i == 1 {
			target = oid
		}
	}

	t.Run("All", func(t *testing.T) {
		actual := 0
		iter, err := db.Iter(nil)
		require.NoError(t, err)

		for iter.Next() {
			actual++

			k := iter.Key()
			require.Equal(t, cid, k.CollectionID())

			obj := iter.Object()
			require.NotNil(t, obj)
		}

		iter.Release()
		require.NoError(t, iter.Error())
		require.Equal(t, nkeys, actual)
	})

	t.Run("Prefix", func(t *testing.T) {
		k := key.New(cid, target, nil)
		iter, err := db.Iter(k.ObjectPrefix())
		require.NoError(t, err)

		actual := 0
		for iter.Next() {
			actual++
			require.Equal(t, target, iter.Key().ObjectID())
		}

		iter.Release()
		require.NoError(t, iter.Error())
		require.Equal(t, oids[target], actual)
	})

	t.Run("Range", func(t *testing.T) {
		k := key.New(cid, target, nil)
		iter, err := db.Range(k.ObjectPrefix(), k.ObjectLimit())
		require.NoError(t, err)

		actual := 0
		for iter.Next() {
			actual++
			require.Equal(t, target, iter.Key().ObjectID())
		}

		iter.Release()
		require.NoError(t, iter.Error())
		require.Equal(t, oids[target], actual)
	})

	t.Run("Error", func(t *testing.T) {
		iter, err := db.Iter(nil)
		require.NoError(t, err)

		iter.Release()
		require.False(t, iter.Next())
		require.ErrorIs(t, iter.Error(), iterator.ErrIterReleased)
	})
}
