package lani_test

import (
	"crypto/rand"
	"errors"
	"io"
	"math"
	"testing"
	"time"

	. "github.com/rotationalio/honu/pkg/store/lani"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/ulid"
)

func TestUnmarshal(t *testing.T) {
	mock := &Mock{[]byte("hello world"), nil}
	data, err := Marshal(mock)
	require.NoError(t, err, "could not marshall mock encodable")

	t.Run("Happy", func(t *testing.T) {
		cmp := &Mock{}
		err := Unmarshal(data, cmp)
		require.NoError(t, err, "could not unmarshall mock encodable")
		require.Equal(t, mock.Data, cmp.Data)
	})

	t.Run("Error", func(t *testing.T) {
		cmp := &Mock{nil, errors.New("whoopsie")}
		err := Unmarshal(data, cmp)
		require.EqualError(t, err, "whoopsie")
		require.Nil(t, cmp.Data)
	})
}

func TestDecoder(t *testing.T) {
	// Decoder tests essentially just ensure that the decoder can correctly deserialize
	// the output of an encoder for the specified type. There are a couple of
	// specialized tests and error checking here and there, but the tests attempt to
	// keep things as simple as possible for future maintainability.
	enc := &Encoder{}

	t.Run("Decode", func(t *testing.T) {
		t.Run("Nil", func(t *testing.T) {
			dec := NewDecoder(nil)
			out, err := dec.Decode()
			require.ErrorIs(t, err, io.EOF, "expected EOF when decoding nil")
			require.Nil(t, out)
		})

		t.Run("InvalidLength", func(t *testing.T) {
			dec := NewDecoder([]byte{0xff, 0xff})
			out, err := dec.Decode()
			require.ErrorIs(t, err, ErrNoLength, "expected incorrect length error")
			require.Nil(t, out)
		})

		t.Run("UnexpectedEOF", func(t *testing.T) {
			dec := NewDecoder([]byte{0xff, 0x12, 0x23, 0x42, 0xf2, 0x21})
			out, err := dec.Decode()
			require.ErrorIs(t, err, io.ErrUnexpectedEOF, "expected EOF before read complete")
			require.Nil(t, out)
		})

		t.Run("NoData", func(t *testing.T) {
			dec := NewDecoder([]byte{0x00})
			out, err := dec.Decode()
			require.NoError(t, err, "expected no error for a nil bytes frame")
			require.Nil(t, out)
		})

		t.Run("Bytes", func(t *testing.T) {
			tests := []struct {
				in  []byte
				err error
			}{
				{nil, nil},
				{[]byte{0x42}, nil},
				{[]byte("hello world"), nil},
			}

			for i, tc := range tests {
				_, err := enc.Encode(tc.in)
				require.NoError(t, err, "could not encode input in test case %d", i)

				out, err := NewDecoder(enc.Bytes()).Decode()
				CompareOrErr(t, i, tc.in, out, err, tc.err)
			}

		})
	})

	t.Run("String", func(t *testing.T) {
		tests := []struct {
			in  string
			err error
		}{
			{"", nil},
			{"hello world", nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeString(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeString()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("StringSlice", func(t *testing.T) {
		tests := []struct {
			in  []string
			err error
		}{
			{nil, nil},
			{[]string{""}, nil},
			{[]string{"hello world"}, nil},
			{[]string{"", "", "", "", "", ""}, nil},
			{[]string{"", "hello", "", "", "world", ""}, nil},
			{[]string{"apples", "oranges", "pineapples", "pumpkins"}, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeStringSlice(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeStringSlice()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Byte", func(t *testing.T) {
		tests := []struct {
			in  byte
			err error
		}{
			{0x0, nil},
			{0xff, nil},
			{0x12, nil},
			{0xf0, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeByte(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeByte()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Bool", func(t *testing.T) {
		tests := []struct {
			in  bool
			err error
		}{
			{false, nil},
			{true, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeBool(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeBool()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("BadBool", func(t *testing.T) {
		_, err := enc.EncodeByte(0xf2)
		require.NoError(t, err, "could not encode input")

		out, err := NewDecoder(enc.Bytes()).DecodeBool()
		require.ErrorIs(t, err, ErrParseBoolean)
		require.False(t, out)
	})

	t.Run("Uint8", func(t *testing.T) {
		tests := []struct {
			in  uint8
			err error
		}{
			{0x0, nil},
			{0xff, nil},
			{0x12, nil},
			{0xf0, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeUint8(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeUint8()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Uint32", func(t *testing.T) {
		tests := []struct {
			in  uint32
			err error
		}{
			{0x0, nil},
			{0xff, nil},
			{0xffff, nil},
			{0xffffff, nil},
			{0xffffffff, nil},
			{math.MaxUint32, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeUint32(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeUint32()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Uint64", func(t *testing.T) {
		tests := []struct {
			in  uint64
			err error
		}{
			{0x0, nil},
			{0xff, nil},
			{0xffff, nil},
			{0xffffff, nil},
			{0xffffffff, nil},
			{0xffffffffff, nil},
			{0xffffffffffff, nil},
			{0xffffffffffffff, nil},
			{0xffffffffffffffff, nil},
			{math.MaxUint64, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeUint64(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeUint64()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Int64", func(t *testing.T) {
		tests := []struct {
			in  int64
			err error
		}{
			{0x0, nil},
			{0xff, nil},
			{0xffff, nil},
			{0xffffff, nil},
			{0xffffffff, nil},
			{0xffffffffff, nil},
			{0xffffffffffff, nil},
			{0xffffffffffffff, nil},
			{-0xff, nil},
			{-0xffff, nil},
			{-0xffffff, nil},
			{-0xffffffff, nil},
			{-0xffffffffff, nil},
			{-0xffffffffffff, nil},
			{-0xffffffffffffff, nil},
			{math.MaxInt64, nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeInt64(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeInt64()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Fixed", func(t *testing.T) {
		tests := []struct {
			in  []byte
			err error
		}{
			{nil, io.EOF},
			{[]byte{0x42}, nil},
			{[]byte("hello world"), nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeFixed(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeFixed(len(tc.in))
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("ULID", func(t *testing.T) {
		tests := []struct {
			in  ulid.ULID
			err error
		}{
			{ulid.ULID{}, nil},
			{ulid.MustParse("01JB24RGJ9FH1X04F9S5A26QAT"), nil},
			{ulid.MustNew(ulid.Now(), rand.Reader), nil},
		}

		for i, tc := range tests {
			_, err := enc.EncodeULID(tc.in)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeULID()
			CompareOrErr(t, i, tc.in, out, err, tc.err)
		}

	})

	t.Run("Time", func(t *testing.T) {
		tests := []time.Time{
			{},
			time.Date(2024, 04, 07, 12, 32, 41, 821, time.UTC),
			time.Now(),
		}

		for i, tc := range tests {
			_, err := enc.EncodeTime(tc)
			require.NoError(t, err, "could not encode input in test case %d", i)

			out, err := NewDecoder(enc.Bytes()).DecodeTime()
			require.NoError(t, err, "could not decode timestamp")
			require.True(t, tc.Equal(out), "timestamps are not equal in testcase %d", i)
		}
	})

	t.Run("Struct", func(t *testing.T) {
		t.Run("Nil", func(t *testing.T) {
			var obj *Mock

			_, err := enc.EncodeStruct(obj)
			require.NoError(t, err, "could not encode nil struct input")

			var cmp *Mock
			isNil, err := NewDecoder(enc.Bytes()).DecodeStruct(cmp)
			require.NoError(t, err, "could not decode nil struct")
			require.True(t, isNil, "expected isNil to be true")
		})

	})

}

func CompareOrErr(t *testing.T, i int, in, out any, err, target error) {
	if target == nil {
		require.NoError(t, err, "could not decode input in test case %d", i)
		require.Equal(t, in, out, "original input did not match decoded output in test case %d", i)
	} else {
		require.ErrorIs(t, err, target, "expected error on test case %d", i)
		require.Nil(t, out, "expected nil on an expected decode error for test case %d", i)
	}
}
