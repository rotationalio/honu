package mime_test

import (
	"mime"
	"testing"

	. "github.com/rotationalio/honu/pkg/mime"
	"github.com/stretchr/testify/require"
)

func TestStringParse(t *testing.T) {
	testCases := []MIME{
		ANY,
		OCTET_STREAM,
		TEXT,
		JSON,
		MSGPACK,
	}

	for i, tc := range testCases {
		actual, err := Parse(tc.String())
		require.NoError(t, err, "test case %d errored", i)
		require.Equal(t, tc, actual, "no match on test case %d", i)
	}
}

func TestContentType(t *testing.T) {
	t.Run("WithCharset", func(t *testing.T) {
		testCases := []MIME{
			TEXT,
			JSON,
		}

		for i, tc := range testCases {
			_, params, err := mime.ParseMediaType(tc.ContentType())
			require.NoError(t, err, "test case %d errored", i)
			require.Contains(t, params, "charset", "no charset in test case %d", i)
			require.Equal(t, params["charset"], "utf-8", "incorrect charset in test case %d", i)
		}
	})

	t.Run("WithoutCharset", func(t *testing.T) {
		testCases := []MIME{
			ANY,
			OCTET_STREAM,
			MSGPACK,
		}

		for i, tc := range testCases {
			mt, params, err := mime.ParseMediaType(tc.ContentType())
			require.NoError(t, err, "test case %d errored", i)
			require.NotContains(t, params, "charset", "charset in test case %d", i)
			require.Equal(t, mt, tc.String(), "no match on test case %d", i)
		}
	})
}

func TestError(t *testing.T) {
	t.Run("ParseMediaType", func(t *testing.T) {
		m, err := Parse("foo/bar; param=")
		require.ErrorIs(t, err, mime.ErrInvalidMediaParameter, "expected an error when parsing a bad mediatype")
		require.Equal(t, UNKNOWN, m, "expected unknown returned on error")
	})

	t.Run("UnknownMediaType", func(t *testing.T) {
		m, err := Parse("foo")
		require.Error(t, err, "expected an error when parsing a bad mediatype")
		require.Equal(t, UNKNOWN, m, "expected unknown returned on error")
	})

	t.Run("StringPanic", func(t *testing.T) {
		m := MIME(65499)
		require.Panics(t, func() { m.String() })
	})
}
