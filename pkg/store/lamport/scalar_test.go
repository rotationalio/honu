package lamport_test

import (
	"encoding/json"
	"math/rand/v2"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	. "go.rtnl.ai/honu/pkg/store/lamport"
)

func TestScalar(t *testing.T) {
	zero := &Scalar{}
	one := &Scalar{1, 1}

	t.Run("IsZero", func(t *testing.T) {
		require.True(t, zero.IsZero(), "empty scalar is not zero")
		require.False(t, one.IsZero(), "1.1 should not be zero")
		require.False(t, (&Scalar{1, 0}).IsZero(), "1.0 should not be zero")
		require.False(t, (&Scalar{0, 1}).IsZero(), "0.1 should not be zero")
	})

	t.Run("Serialize", func(t *testing.T) {
		current := &Scalar{}
		for i := 0; i < 128; i++ {
			data, err := current.MarshalBinary()
			require.NoError(t, err, "could not marshal %s", current)

			cmpr := &Scalar{}
			err = cmpr.UnmarshalBinary(data)
			require.NoError(t, err, "could not unmarshal %d bytes", len(data))

			require.Equal(t, current, cmpr, "unmarshaled scalar does not match marshaled one")

			current = randNextScalar(current)
		}
	})

	t.Run("Text", func(t *testing.T) {
		current := &Scalar{}
		for i := 0; i < 128; i++ {
			data, err := current.MarshalText()
			require.NoError(t, err, "could not marshal %s", current)

			cmpr := &Scalar{}
			err = cmpr.UnmarshalText(data)
			require.NoError(t, err, "could not unmarshal %d bytes", len(data))

			require.Equal(t, current, cmpr, "unmarshaled scalar does not match marshaled one")

			current = randNextScalar(current)
		}
	})

	t.Run("BadText", func(t *testing.T) {
		testCases := []string{
			"123",
			"a.b",
			"a.123",
			"1.abc",
			"",
			"1.1.1",
			"1.",
			".1",
		}

		for i, tc := range testCases {
			err := (&Scalar{}).UnmarshalText([]byte(tc))
			require.Error(t, err, "expected errror on test case %d", i)
		}
	})

	t.Run("Binary", func(t *testing.T) {
		vers := &Scalar{42, 198}
		data, err := vers.MarshalBinary()
		require.NoError(t, err, "could not marshal scalar as a binary value")
		require.Equal(t, []byte{0x2a, 0xc6, 0x1}, data)
	})

	t.Run("JSON", func(t *testing.T) {
		s := &Scalar{}
		data := []byte(`"8.16"`)

		require.NoError(t, json.Unmarshal(data, s), "could not unmarshal s")
		require.Equal(t, &Scalar{8, 16}, s, "incorrect unmarshal")

		cmpt, err := json.Marshal(s)
		require.NoError(t, err, "could not marshal s")
		require.Equal(t, data, cmpt, "unexpected marshaled data")
	})
}

func TestCompare(t *testing.T) {
	testCases := []struct {
		a, b     *Scalar
		expected int
	}{
		{
			nil, nil, 0,
		},
		{
			&Scalar{}, &Scalar{}, 0,
		},
		{
			nil, &Scalar{}, 0,
		},
		{
			&Scalar{}, nil, 0,
		},
		{
			&Scalar{1, 1}, &Scalar{}, 1,
		},
		{
			&Scalar{}, &Scalar{1, 1}, -1,
		},
		{
			&Scalar{1, 2}, &Scalar{1, 1}, 1,
		},
		{
			&Scalar{1, 1}, &Scalar{1, 2}, -1,
		},
		{
			&Scalar{1, 1}, &Scalar{1, 1}, 0,
		},
		{
			&Scalar{2, 1}, &Scalar{1, 1}, 1,
		},
		{
			&Scalar{1, 1}, &Scalar{2, 1}, -1,
		},
		{
			&Scalar{2, 2}, &Scalar{2, 2}, 0,
		},
		{
			&Scalar{4, 4}, &Scalar{4, 9}, -1,
		},
		{
			&Scalar{1, 10}, &Scalar{10, 10}, -1,
		},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, Compare(tc.a, tc.b), "test case %d failed", i)

		// Handle semantic comparisons
		switch tc.expected {
		case 0:
			require.True(t, tc.a.Equals(tc.b), "test case %d equality failed", i)
			require.True(t, tc.b.Equals(tc.a), "test case %d equality failed", i)

			require.False(t, tc.a.Before(tc.b), "test case %d before failed", i)
			require.False(t, tc.b.Before(tc.a), "test case %d before failed", i)

			require.False(t, tc.a.After(tc.b), "test case %d after failed", i)
			require.False(t, tc.b.After(tc.a), "test case %d after failed", i)
		case 1:
			require.False(t, tc.a.Equals(tc.b), "test case %d equality failed", i)
			require.False(t, tc.b.Equals(tc.a), "test case %d equality failed", i)

			require.False(t, tc.a.Before(tc.b), "test case %d before failed", i)
			require.True(t, tc.b.Before(tc.a), "test case %d before failed", i)

			require.True(t, tc.a.After(tc.b), "test case %d after failed", i)
			require.False(t, tc.b.After(tc.a), "test case %d after failed", i)
		case -1:
			require.False(t, tc.a.Equals(tc.b), "test case %d equality failed", i)
			require.False(t, tc.b.Equals(tc.a), "test case %d equality failed", i)

			require.True(t, tc.a.Before(tc.b), "test case %d before failed", i)
			require.False(t, tc.b.Before(tc.a), "test case %d before failed", i)

			require.False(t, tc.a.After(tc.b), "test case %d after failed", i)
			require.True(t, tc.b.After(tc.a), "test case %d after failed", i)
		}

	}

}

func TestSort(t *testing.T) {
	// Create a list of random scalars
	scalars := make(Scalars, 128)
	for i := range scalars {
		scalars[i] = randScalar()
	}

	// Sort the scalars
	sort.Sort(scalars)

	// Ensure that the scalars are sorted
	for i := 1; i < len(scalars); i++ {
		require.True(t, Compare(scalars[i-1], scalars[i]) <= 0, "scalars[%d] is not before or equal to scalars[%d]", i-1, i)
	}
}

func randScalar() *Scalar {
	return &Scalar{
		PID: uint32(rand.Int32N(24)),
		VID: uint64(rand.Int64N(48)),
	}
}

func randNextScalar(prev *Scalar) *Scalar {
	s := &Scalar{}
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
