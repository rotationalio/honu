package lani_test

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

func TestMarshal(t *testing.T) {
	mock := &Mock{[]byte("hello world"), nil}
	data, err := Marshal(mock)
	require.NoError(t, err, "could not marshall mock encodable")
	require.GreaterOrEqual(t, mock.Size(), len(data), "data was not the expected size")
}

func TestMarshalError(t *testing.T) {
	mock := &Mock{nil, errors.New("whoopsie")}
	data, err := Marshal(mock)
	require.EqualError(t, err, "whoopsie")
	require.Nil(t, data)
}

func TestEncoder(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		enc := &Encoder{}
		require.Zero(t, enc.Len(), "expected length to be zero")
		require.Zero(t, enc.Cap(), "expected capacity to be zero")
	})

	t.Run("Grow", func(t *testing.T) {
		t.Run("Small", func(t *testing.T) {
			enc := &Encoder{}
			enc.Grow(32)

			require.Zero(t, enc.Len(), "expected length to be zero")
			require.Equal(t, 64, enc.Cap(), "expected capacity to be the small buffer size")
		})

		t.Run("SingleAlloc", func(t *testing.T) {
			enc := &Encoder{}
			enc.Grow(512)

			require.Zero(t, enc.Len(), "expected length to be zero")
			require.Equal(t, 512, enc.Cap(), "expected capacity to be the amount grown by")

			enc.Grow(512)
			require.Zero(t, enc.Len(), "expected length to be zero")
			require.Equal(t, 512, enc.Cap(), "expected capacity to be unchanged")
		})

		t.Run("TooLarge", func(t *testing.T) {
			// Causes a panic in the grow method
			growBy := int(^uint(0)>>1) - 256
			require.Panics(t, func() {
				enc := &Encoder{}
				enc.Grow(512)
				enc.Grow(growBy)
			}, "expected a panic when capacity exceeded max integer")
		})

		t.Run("GrowSlicePanic", func(t *testing.T) {
			// This differs from the TooLarge panic test in that it does not initially
			// grow the buffer, ensuring a call is made to growSlice to cause the panic
			growBy := int(^uint(0)>>1) - 256
			require.Panics(t, func() {
				enc := &Encoder{}
				enc.Grow(growBy)
			}, "expected a panic when capacity exceeded max integer")
		})

		t.Run("Negative", func(t *testing.T) {
			require.Panics(t, func() {
				enc := &Encoder{}
				enc.Grow(-42)
			}, "expected a panic when trying to grow by a negative number")
		})

		t.Run("CauseReset", func(t *testing.T) {
			enc := &Encoder{}
			enc.Encode([]byte("this is a larger string than the next string"))
			oldcap := enc.Cap()

			enc.Bytes()
			enc.Encode([]byte("this is a smaller string"))
			require.Equal(t, oldcap, enc.Cap())
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("Single", func(t *testing.T) {
			data := []byte("I'm a little teapot short and stout, here is my handle, here is my spout. When I'm feeling steamed I jump and shout; tip me over and pour me out!")

			enc := &Encoder{}
			n, err := enc.Encode(data)
			require.NoError(t, err, "could not write bytes to the encoder")
			require.Equal(t, n, enc.Len(), "expected the amount written to be the new length")
			require.Equal(t, 2+len(data), n, "expected the frame to be 2+len(data)")

			cmp := enc.Bytes()
			require.Equal(t, data, cmp[2:], "expected data to be written after length")
		})

		t.Run("Multiple", func(t *testing.T) {
			data := [][]byte{
				[]byte("I'm a little teapot short and stout"),
				[]byte("here is my handle, here is my spout"),
				[]byte("when I'm feeling steamed I jump and shout"),
				[]byte("pour me over and tip me out!"),
			}

			N := len(data)
			enc := &Encoder{}

			for _, row := range data {
				N += len(row)

				n, err := enc.Encode(row)
				require.NoError(t, err, "could not write row")
				require.Equal(t, 1+len(row), n, "expected 1+len(row) data to be written; did the text get longer?")
			}

			require.Equal(t, N, enc.Len(), "expected len(data) frames + the length of the data")
		})

		t.Run("GrowBeforeMultiple", func(t *testing.T) {
			data := [][]byte{
				[]byte("I'm a little teapot short and stout"),
				[]byte("here is my handle, here is my spout"),
				[]byte("when I'm feeling steamed I jump and shout"),
				[]byte("pour me over and tip me out!"),
			}

			N := len(data)
			for _, row := range data {
				N += len(row)
			}

			enc := &Encoder{}
			enc.Grow(N)

			for _, row := range data {
				n, err := enc.Encode(row)
				require.NoError(t, err, "could not write row")
				require.Equal(t, 1+len(row), n, "expected 1+len(row) data to be written; did the text get longer?")
			}

			require.Equal(t, N, enc.Len(), "expected len(data) frames + the length of the data")
		})
	})

	t.Run("WriteReadReset", func(t *testing.T) {
		data := [][][]byte{
			{
				[]byte(""),
			},
			{
				[]byte("hello world"),
			},
			{
				[]byte("hello world"),
				[]byte("this is the computer"),
			},
			{
				[]byte("hello world"),
				[]byte("this is the computer"),
				[]byte("running tests"),
			},
		}

		enc := &Encoder{}
		for i, rows := range data {
			N := 0
			for j, row := range rows {
				n, err := enc.Encode(row)
				require.NoError(t, err, "could not encode dataset %d, row %d", i, j)
				require.Equal(t, 1+len(row), n, "unexpected number of bytes written")

				N += 1 + len(row)
				require.Equal(t, N, enc.Len(), "unexpected encoder length after write")
			}

			cmp := enc.Bytes()
			require.Len(t, cmp, N)
			require.Equal(t, 0, enc.Len())

			require.Nil(t, enc.Bytes(), "the second bytes read should be nil and cause a reset")
		}

		// After all of this the capacity should be less than 64
		require.LessOrEqual(t, enc.Cap(), 64, "capacity larger than expected")
	})

	t.Run("StringVsBytes", func(t *testing.T) {
		e := &Encoder{}
		s := "the brown bear jumps through the woods"

		_, err := e.EncodeString(s)
		require.NoError(t, err, "could not encode string")
		s1 := e.Bytes()

		_, err = e.Encode([]byte(s))
		require.NoError(t, err, "could not encode bytes")
		s2 := e.Bytes()

		require.Equal(t, s1, s2, "encoding string and bytes did not match")
	})

	t.Run("EncodeStringSlice", func(t *testing.T) {
		tests := []struct {
			n    int
			strs []string
		}{
			{1, nil},
			{1, []string{}},
			{13, []string{"hello world"}},
			{25, []string{"hello world", "this is bob"}},
			{26, []string{"hello world", "", "this is bob"}},
			{4, []string{"", "", ""}},
		}

		for i, tc := range tests {
			enc := &Encoder{}
			n, err := enc.EncodeStringSlice(tc.strs)
			require.NoError(t, err, "could not encode slice on test case %d", i)
			require.Equal(t, tc.n, n, "mismatch in expected amount of data written")
			require.Equal(t, tc.n, enc.Len(), "expected length to be the same as written")
		}
	})

	t.Run("OneByteData", func(t *testing.T) {
		// Use only one encoder that should only have one byte in it
		e := &Encoder{}

		t.Run("Byte", func(t *testing.T) {
			n, err := e.EncodeByte(0x42)
			require.NoError(t, err, "could not encode byte")
			require.Equal(t, 1, n, "expected only one byte written")
			require.Equal(t, []byte{0x42}, e.Bytes())
		})

		t.Run("Uint8", func(t *testing.T) {
			n, err := e.EncodeUint8(42)
			require.NoError(t, err, "could not encode uint8")
			require.Equal(t, 1, n, "expected only one byte written")
			require.Equal(t, []byte{0x2a}, e.Bytes())
		})

		t.Run("True", func(t *testing.T) {
			n, err := e.EncodeBool(true)
			require.NoError(t, err, "could not encode true")
			require.Equal(t, 1, n, "expected only one byte written")
			require.Equal(t, []byte{0x01}, e.Bytes())
		})

		t.Run("False", func(t *testing.T) {
			n, err := e.EncodeBool(false)
			require.NoError(t, err, "could not encode false")
			require.Equal(t, 1, n, "expected only one byte written")
			require.Equal(t, []byte{0x00}, e.Bytes())
		})

	})

	t.Run("Integers", func(t *testing.T) {
		t.Run("Uint32", func(t *testing.T) {
			tests := []struct {
				i    uint32
				size int
			}{
				{127, 1},
				{16383, 2},
				{2097151, 3},
				{268435455, 4},
				{4294967295, 5},
			}

			enc := &Encoder{}
			for i, tc := range tests {
				n, err := enc.EncodeUint32(tc.i)
				require.NoError(t, err, "could not encode uint32 in test case %d", i)
				require.Equal(t, tc.size, n, "wrong number of bytes written in test case %d", i)
				require.Equal(t, n, enc.Len(), "unexpected length of byte array in test case %d", i)

				// Clear the encoder
				enc.Reset()
			}
		})

		t.Run("Uint64", func(t *testing.T) {
			tests := []struct {
				i    uint64
				size int
			}{
				{127, 1},
				{16383, 2},
				{2097151, 3},
				{268435455, 4},
				{34359738367, 5},
				{4398046511103, 6},
				{562949953421311, 7},
				{72057594037927928, 8},
				{9223372036854775807, 9},
				{18446744073709551615, 10},
			}

			enc := &Encoder{}
			for i, tc := range tests {
				n, err := enc.EncodeUint64(tc.i)
				require.NoError(t, err, "could not encode uint64 in test case %d", i)
				require.Equal(t, tc.size, n, "wrong number of bytes written in test case %d", i)
				require.Equal(t, n, enc.Len(), "unexpected length of byte array in test case %d", i)

				// Clear the encoder
				enc.Reset()
			}
		})

		t.Run("Int64", func(t *testing.T) {
			tests := []struct {
				i    int64
				size int
			}{
				{-9223372036854775807, 10},
				{-72057594037927928, 9},
				{-562949953421311, 8},
				{-4398046511103, 7},
				{-34359738367, 6},
				{-268435455, 5},
				{-2097151, 4},
				{-16383, 3},
				{-127, 2},
				{-32, 1},
				{0, 1},
				{32, 1},
				{127, 2},
				{16383, 3},
				{2097151, 4},
				{268435455, 5},
				{34359738367, 6},
				{4398046511103, 7},
				{562949953421311, 8},
				{72057594037927928, 9},
				{9223372036854775807, 10},
			}

			enc := &Encoder{}
			for i, tc := range tests {
				n, err := enc.EncodeInt64(tc.i)
				require.NoError(t, err, "could not encode int64 in test case %d", i)
				require.Equal(t, tc.size, n, "wrong number of bytes written in test case %d", i)
				require.Equal(t, n, enc.Len(), "unexpected length of byte array in test case %d", i)

				// Clear the encoder
				enc.Reset()
			}
		})
	})

	t.Run("Fixed", func(t *testing.T) {
		enc := &Encoder{}
		data := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}
		n, err := enc.EncodeFixed(data)
		require.NoError(t, err, "could not encode fixed")
		require.Equal(t, len(data), n, "expected only the data to be written")
	})

	t.Run("ULID", func(t *testing.T) {
		enc := &Encoder{}
		data := ulid.MustNew(ulid.Now(), rand.Reader)
		n, err := enc.EncodeULID(data)
		require.NoError(t, err, "could not encode ulid")
		require.Equal(t, len(data), n, "expected only the data to be written")
	})

	t.Run("Time", func(t *testing.T) {
		enc := &Encoder{}
		data := time.Now()
		n, err := enc.EncodeTime(data)
		require.NoError(t, err, "could not encode timestamp")
		require.LessOrEqual(t, n, binary.MaxVarintLen64, "unexpected number of bytes written")
	})

	t.Run("ZeroTime", func(t *testing.T) {
		enc := &Encoder{}
		n, err := enc.EncodeTime(time.Time{})
		require.NoError(t, err, "could not encode timestamp")
		require.Equal(t, n, 1, "unexpected number of bytes written")
		require.Equal(t, []byte{0x00}, enc.Bytes(), "zero valud time not written as zero")
	})

	t.Run("Struct", func(t *testing.T) {
		t.Run("Nil", func(t *testing.T) {
			var obj *Mock
			enc := &Encoder{}
			n, err := enc.EncodeStruct(obj)
			require.NoError(t, err, "could not encode nil struct")
			require.Equal(t, 1, n, "expected 1 byte encoded for nil structs")
		})
	})

}
