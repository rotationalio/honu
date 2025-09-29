package index_test

import (
	"crypto/rand"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/index"
	"go.rtnl.ai/ulid"
)

func TestInsort(t *testing.T) {
	t.Run("Random", func(t *testing.T) {
		// Create an array of random ULIDs
		arr := make(index.Array, 0, 128)
		for i := 0; i < 128; i++ {
			id := make([]byte, 16)
			rand.Read(id)

			arr = arr.Insert(ulid.ULID(id))
		}

		require.Len(t, arr, 128)
		IsSorted(t, arr)
		NoDuplicates(t, arr)
	})

	t.Run("Ordered", func(t *testing.T) {
		entropy := ulid.Monotonic(rand.Reader, 0)
		arr := make(index.Array, 0, 128)

		for i := 0; i < 128; i++ {
			id := ulid.MustNew(ulid.Now(), entropy)
			arr = arr.Insert(id)
		}

		require.Len(t, arr, 128)
		IsSorted(t, arr)
		NoDuplicates(t, arr)
	})

	t.Run("Reverse", func(t *testing.T) {
		entropy := ulid.Monotonic(rand.Reader, 0)
		ids := make([]ulid.ULID, 0, 128)

		for i := 0; i < 128; i++ {
			id := ulid.MustNew(ulid.Now(), entropy)
			ids = append(ids, id)
		}

		arr := make(index.Array, 0, 128)
		for i := len(ids) - 1; i >= 0; i-- {
			arr = arr.Insert(ids[i])
		}

		require.Len(t, arr, 128)
		IsSorted(t, arr)
		NoDuplicates(t, arr)
	})

	t.Run("Duplicates", func(t *testing.T) {
		data := make([]ulid.ULID, 0, 384)
		for i := 0; i < 128; i++ {
			id := make([]byte, 16)
			rand.Read(id)
			data = append(data, ulid.ULID(id))
			data = append(data, ulid.ULID(id))
			data = append(data, ulid.ULID(id))
		}

		r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })

		arr := make(index.Array, 0, 128)
		for _, id := range data {
			arr = arr.Insert(id)
		}

		require.Len(t, arr, 128)
		IsSorted(t, arr)
		NoDuplicates(t, arr)
	})

	t.Run("DoubleAppend", func(t *testing.T) {
		id := ulid.MustNew(ulid.Now(), rand.Reader)
		arr := make(index.Array, 0, 16)
		arr = arr.Insert(id)
		arr = arr.Insert(id)
		IsSorted(t, arr)
		NoDuplicates(t, arr)
	})
}

func TestContains(t *testing.T) {
	arr := index.Array{
		ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727"),
		ulid.MustParse("01K6APMNR3Z8G114TRYR8224D4"),
		ulid.MustParse("01K6APP9VQ5DXV98BQ8ECZ6D78"),
		ulid.MustParse("01K6APPEGF6XHJ4S8JVMGPC704"),
		ulid.MustParse("01K6APPKG3JRW6X1GZEZJADSV2"),
		ulid.MustParse("01K6APPRPDXVVJF76DGSTY02KN"),
		ulid.MustParse("01K6APPYJB3TDQ1DYQR43K70R3"),
		ulid.MustParse("01K6APQ4T9BHNGBFXHVEDF0R1X"),
		ulid.MustParse("01K6APQBEEK6XBP3VDJVE7VW77"),
		ulid.MustParse("01K6APQGYV9CGCBKVTXD1GD5NG"),
		ulid.MustParse("01K6APQPF23472MPVGMW4WK5TP"),
		ulid.MustParse("01K6APR3QYFNGK5323E5JSN3V0"),
		ulid.MustParse("01K6APR969SDEX3FFSRATCXGVM"),
		ulid.MustParse("01K6APREDFY7CG1P4HRCBW5706"),
		ulid.MustParse("01K6APRK84PRH5ESY053TA83NB"),
		ulid.MustParse("01K6APRRJVF2FCJNMWCRNRGFQ8"),
		ulid.MustParse("01K6APRXZE47845QY3RT5CV74G"),
		ulid.MustParse("01K6APS2V1FGKS3RMXF3BH4P65"),
		ulid.MustParse("01K6APS7YM1KRJMMQ5EMYCV0ES"),
		ulid.MustParse("01K6APSCXKCQ2FS8F9GDJRF1EZ"),
	}

	IsSorted(t, arr)
	require.True(t, arr.Contains(ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727")))
	require.True(t, arr.Contains(ulid.MustParse("01K6APSCXKCQ2FS8F9GDJRF1EZ")))
	require.True(t, arr.Contains(ulid.MustParse("01K6APQPF23472MPVGMW4WK5TP")))
	require.True(t, arr.Contains(ulid.MustParse("01K6APPRPDXVVJF76DGSTY02KN")))
	require.False(t, arr.Contains(ulid.MustParse("01K6APMNR3Z8G114TRYR8224D5")))
}

func TestRemove(t *testing.T) {
	arr := index.Array{
		ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727"),
		ulid.MustParse("01K6APMNR3Z8G114TRYR8224D4"),
		ulid.MustParse("01K6APP9VQ5DXV98BQ8ECZ6D78"),
		ulid.MustParse("01K6APPEGF6XHJ4S8JVMGPC704"),
		ulid.MustParse("01K6APPKG3JRW6X1GZEZJADSV2"),
		ulid.MustParse("01K6APPRPDXVVJF76DGSTY02KN"),
		ulid.MustParse("01K6APPYJB3TDQ1DYQR43K70R3"),
		ulid.MustParse("01K6APQ4T9BHNGBFXHVEDF0R1X"),
		ulid.MustParse("01K6APQBEEK6XBP3VDJVE7VW77"),
		ulid.MustParse("01K6APQGYV9CGCBKVTXD1GD5NG"),
		ulid.MustParse("01K6APQPF23472MPVGMW4WK5TP"),
		ulid.MustParse("01K6APR3QYFNGK5323E5JSN3V0"),
		ulid.MustParse("01K6APR969SDEX3FFSRATCXGVM"),
		ulid.MustParse("01K6APREDFY7CG1P4HRCBW5706"),
		ulid.MustParse("01K6APRK84PRH5ESY053TA83NB"),
		ulid.MustParse("01K6APRRJVF2FCJNMWCRNRGFQ8"),
		ulid.MustParse("01K6APRXZE47845QY3RT5CV74G"),
		ulid.MustParse("01K6APS2V1FGKS3RMXF3BH4P65"),
		ulid.MustParse("01K6APS7YM1KRJMMQ5EMYCV0ES"),
		ulid.MustParse("01K6APSCXKCQ2FS8F9GDJRF1EZ"),
	}

	IsSorted(t, arr)
	require.Len(t, arr, 20)

	// Remove an element that is present.
	require.True(t, arr.Contains(ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727")))
	arr = arr.Remove(ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727"))
	require.Len(t, arr, 19)
	require.False(t, arr.Contains(ulid.MustParse("01K6APM9QD21HZ3N7Y74E65727")))

	// Remove an ID that is not present, length should not change.
	arr = arr.Remove(ulid.MustParse("01K6APM9QD21HZ3N7Y74E65728"))
	require.Len(t, arr, 19)

}

func IsSorted(t *testing.T, a index.Array) {
	t.Helper()
	for i := 1; i < len(a); i++ {
		if a[i-1].Compare(a[i]) >= 0 {
			t.Errorf("array is not sorted at index %d: %s >= %s", i, a[i-1], a[i])
		}
	}
}

func NoDuplicates(t *testing.T, a index.Array) {
	t.Helper()
	seen := make(map[ulid.ULID]struct{}, len(a))
	for i, id := range a {
		if _, ok := seen[id]; ok {
			t.Errorf("duplicate id found at index %d: %s", i, id)
		}
		seen[id] = struct{}{}
	}

	// Really we shouldn't make it this far if there are duplicates, but just in case...
	if len(seen) != len(a) {
		t.Errorf("duplicate ids found: seen %d, array length %d", len(seen), len(a))
	}
}
