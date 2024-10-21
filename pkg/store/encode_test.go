package store_test

import (
	"testing"

	. "github.com/rotationalio/honu/pkg/store"
	"github.com/stretchr/testify/require"
)

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

}
