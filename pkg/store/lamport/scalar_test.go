package lamport_test

import (
	"math/rand/v2"
	"testing"

	. "github.com/rotationalio/honu/pkg/store/lamport"
	"github.com/stretchr/testify/require"
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

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "1.1", one.String())
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

func randNextScalar(prev *Scalar) *Scalar {
	s := &Scalar{}
	s.PID = uint32(rand.Int32N(24))
	s.VID = uint64(rand.Int64N(32)) + prev.VID
	return s
}
