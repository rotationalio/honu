package key_test

import (
	"bytes"
	crand "crypto/rand"
	"math"
	"math/rand/v2"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/ulid"
)

func TestNewKey(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 1, PID: 2}

	k := key.New(oid, vers)
	require.NotNil(t, k)
	require.Equal(t, 29, len(k))
	require.Equal(t, oid, k.ObjectID())
	require.Equal(t, *vers, k.Version())
}

func TestKeyLexicographic(t *testing.T) {
	// For the same collection and object ID the keys should be lexicographically sorted
	// by the version to ensure that we can read the latest version either by choosing
	// the first item in a list of sorted keys or the last item.
	oid := ulid.Make()

	t.Run("Static", func(t *testing.T) {
		versions := lamport.Scalars{
			{PID: 1, VID: 1},
			{PID: 1, VID: 2},
			{PID: 1, VID: 3},
			{PID: 2, VID: 3},
			{PID: 3, VID: 3},
			{PID: 1, VID: 4},
			{PID: 4, VID: 4},
			{PID: 4, VID: 9},
			{PID: 1, VID: 10},
			{PID: 2, VID: 10},
			{PID: 2, VID: 10},
			{PID: 4, VID: 11},
			{PID: 13, VID: 11},
		}

		keys := make(key.Keys, len(versions))
		for i, v := range versions {
			keys[i] = key.New(oid, v)
		}

		// Ensure the keys are sorted both by monotonically increasing version and by
		// lexicographic byte order.
		for i := 1; i < len(keys); i++ {
			versa, versb := keys[i-1].Version(), keys[i].Version()
			require.True(t, lamport.Compare(&versa, &versb) <= 0, "%s is not before %s version", versa.String(), versb.String())
			require.True(t, bytes.Compare(keys[i-1][:], keys[i][:]) <= 0, "keys[%d] (%s) is not less than or equal to keys[%d] (%s)", i-1, versa.String(), i, versb.String())
		}
	})

	t.Run("Random", func(t *testing.T) {
		keys := make(key.Keys, 512)
		vers := &lamport.Scalar{VID: 1, PID: 1}

		// Create a list of keys with monotonically increasing versions.
		for i := 0; i < len(keys); i++ {
			keys[i] = key.New(oid, vers)
			vers = randNextScalar(vers)
		}

		// Ensure the keys are sorted both by monotonically increasing version and by
		// lexicographic byte order.
		for i := 1; i < len(keys); i++ {
			versa, versb := keys[i-1].Version(), keys[i].Version()
			require.True(t, lamport.Compare(&versa, &versb) <= 0, "%s is not before %s version", versa.String(), versb.String())
			require.True(t, bytes.Compare(keys[i-1][:], keys[i][:]) <= 0, "keys[%d] (%s) is not less than or equal to keys[%d] (%s)", i-1, versa.String(), i, versb.String())
		}
	})

	t.Run("RandomSort", func(t *testing.T) {
		keys := make(key.Keys, 512)
		vers := &lamport.Scalar{VID: 1, PID: 1}

		// Create a list of keys with monotonically increasing versions.
		for i := 0; i < len(keys); i++ {
			keys[i] = key.New(oid, vers)
			vers = randNextScalar(vers)
		}

		sort.Sort(keys)

		// Ensure the keys are sorted both by monotonically increasing version and by
		// lexicographic byte order.
		for i := 1; i < len(keys); i++ {
			versa, versb := keys[i-1].Version(), keys[i].Version()
			require.True(t, lamport.Compare(&versa, &versb) <= 0, "%s is not before %s version", versa.String(), versb.String())
			require.True(t, bytes.Compare(keys[i-1][:], keys[i][:]) <= 0, "keys[%d] (%s) is not less than or equal to keys[%d] (%s)", i-1, versa.String(), i, versb.String())
		}
	})
}

func TestKeyCheck(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 1, PID: 2}

	t.Run("Valid", func(t *testing.T) {
		k := key.New(oid, vers)
		require.NoError(t, k.Check())
	})

	t.Run("BadSize", func(t *testing.T) {
		badKey := key.Key(make([]byte, 42))
		require.ErrorIs(t, badKey.Check(), key.ErrBadSize)
	})

	t.Run("BadVersion", func(t *testing.T) {
		badKey := key.New(oid, vers)
		badKey[0] = 0x28
		require.ErrorIs(t, badKey.Check(), key.ErrBadVersion)
	})
}

func TestObjectID(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 80, PID: 122}

	t.Run("Ok", func(t *testing.T) {
		k := key.New(oid, vers)
		require.Equal(t, oid, k.ObjectID())
	})

	t.Run("Panics", func(t *testing.T) {
		badKey := key.Key(make([]byte, 42))
		require.Panics(t, func() {
			badKey.ObjectID()
		})
	})
}

func TestVersion(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 5, PID: 1}

	t.Run("Ok", func(t *testing.T) {
		k := key.New(oid, vers)
		require.Equal(t, *vers, k.Version())
	})

	t.Run("Panics", func(t *testing.T) {
		badKey := key.Key(make([]byte, 42))
		require.Panics(t, func() {
			badKey.Version()
		})
	})
}

func TestObjectPrefix(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 1, PID: 2}

	t.Run("Ok", func(t *testing.T) {
		k := key.New(oid, vers)
		prefix := k.ObjectPrefix()
		require.Len(t, prefix, 17)
		require.Equal(t, uint8(0x01), prefix[0])
		require.Equal(t, oid, ulid.ULID(prefix[1:17]))
	})

	t.Run("Panics", func(t *testing.T) {
		badKey := key.Key(make([]byte, 42))
		require.Panics(t, func() {
			badKey.ObjectPrefix()
		})
	})
}

func TestObjectLimit(t *testing.T) {
	oid := ulid.Make()
	vers := &lamport.Scalar{VID: 1, PID: 2}

	t.Run("Ok", func(t *testing.T) {
		k := key.New(oid, vers)
		limit := k.ObjectLimit()
		require.Len(t, limit, 17)
		require.Equal(t, uint8(0x01), limit[0])
		require.True(t, bytes.Equal(oid[:15], limit[1:16]))
		require.Equal(t, oid[15]+1, limit[16])
	})

	t.Run("Panics", func(t *testing.T) {
		badKey := key.Key(make([]byte, 42))
		require.Panics(t, func() {
			badKey.ObjectLimit()
		})
	})
}

func TestObjectRange(t *testing.T) {
	oid := ulid.Make()

	k := key.New(oid, nil)
	prefix, limit := k.ObjectPrefix(), k.ObjectLimit()

	// The prefix should come before the key
	require.True(t, bytes.Compare(prefix, k[:]) == -1)

	// The limit should be greater than the maximum possible version
	maxver := key.New(oid, &lamport.Scalar{VID: math.MaxUint64, PID: math.MaxUint32})
	require.True(t, bytes.Compare(limit, maxver[:]) == 1)

	// The limit should be less than the zero valued version of the next object
	next := make([]byte, 16)
	copy(next[:15], oid[:15])
	next[15] = oid[15] + 1
	nver := key.New(ulid.ULID(next), nil)
	require.True(t, bytes.Compare(limit, nver[:]) == -1)
}

func TestHasVersion(t *testing.T) {
	oid := ulid.Make()

	t.Run("Nil", func(t *testing.T) {
		k1 := key.New(oid, nil)
		require.False(t, k1.HasVersion())

		k2 := key.New(oid, &lamport.Scalar{})
		require.False(t, k2.HasVersion())
	})

	t.Run("NonNil", func(t *testing.T) {
		cases := [][]byte{
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0x0, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		}

		k1 := key.New(oid, nil)
		k2 := key.New(ulid.Null, nil)

		for _, vers := range cases {
			copy(k1[17:29], vers)
			copy(k2[17:29], vers)
			require.True(t, k1.HasVersion())
			require.True(t, k2.HasVersion())
		}
	})

	t.Run("Panics", func(t *testing.T) {
		badKey := key.Key(make([]byte, 26))
		require.Panics(t, func() {
			badKey.HasVersion()
		})
	})
}

func TestSort(t *testing.T) {
	keys := make(key.Keys, 128)
	for i := 0; i < 128; i++ {
		keys[i] = randKey()
	}

	sort.Sort(keys)

	for i := 1; i < len(keys); i++ {
		require.True(t, bytes.Compare(keys[i-1][:], keys[i][:]) <= 0, "keys[%d] is not less than or equal to keys[%d]", i-1, i)
	}
}

func randKey() key.Key {
	b := make([]byte, 29)
	crand.Read(b)
	return key.Key(b)
}

func randNextScalar(prev *lamport.Scalar) *lamport.Scalar {
	s := &lamport.Scalar{}
	s.PID = uint32(rand.Int32N(24))
	s.VID = uint64(rand.Int64N(5)) + prev.VID

	if !prev.Before(s) {
		s.PID = prev.PID + 1
		if !prev.Before(s) {
			panic("failed to generate next scalar")
		}
	}
	return s
}
