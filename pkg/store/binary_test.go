package store_test

import (
	"testing"

	"github.com/rotationalio/honu/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestReadWriteBytes(t *testing.T) {
	t.Run("EncodeDecode", func(t *testing.T) {
		data := []byte("I'm a little teapot, short and stout")

		w := &store.Writer{}
		n, err := w.WriteBytes(data)
		require.NoError(t, err, "could not write bytes")
		require.Equal(t, len(data)+1, n, "the length component should be 1 byte long + the length of the data")

		r := store.NewReader(w.Bytes())
		out, err := r.ReadBytes()
		require.NoError(t, err, "could not read bytes")
		require.Equal(t, data, out, "did not read the expected encoded data")
	})

	t.Run("MultipleBytes", func(t *testing.T) {
		data := [][]byte{
			[]byte("I'm a little teapot, short and stout"),
			[]byte("here is my handle, here is my spout"),
			[]byte("when I'm feeling steamed up, hear me shout."),
			[]byte("Tip me over and pour me out!"),
		}

		w := &store.Writer{}
		for _, row := range data {
			_, err := w.WriteBytes(row)
			require.NoError(t, err, "could not write data to buffer")
		}

		r := store.NewReader(w.Bytes())
		for _, row := range data {
			out, err := r.ReadBytes()
			require.NoError(t, err, "could not read data from buffer")
			require.Equal(t, row, out, "read bytes did not match original bytes")
		}
	})

	t.Run("SingleByte", func(t *testing.T) {
		data := []byte{112}

		w := &store.Writer{}
		n, err := w.WriteBytes(data)
		require.NoError(t, err, "could not write byte")
		require.Equal(t, 2, n, "expected two bytes written")

		r := store.NewReader(w.Bytes())
		out, err := r.ReadBytes()
		require.NoError(t, err, "could not read byte from buffer")
		require.Equal(t, data, out, "read byte did not match original byte")
	})
}
