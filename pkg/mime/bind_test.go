package mime_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"
	"go.rtnl.ai/honu/pkg/api/v1"
	"go.rtnl.ai/honu/pkg/mime"
)

func TestBind(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.json"))
		r.Header.Set("Content-Type", "application/json")

		actual := &api.StatusReply{}
		err := mime.Bind(w, r, actual)
		require.NoError(t, err, "should bind JSON without error")
		require.Equal(t, expected, actual, "bound object should match expected")
	})

	t.Run("MsgPack", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.msgpack"))
		r.Header.Set("Content-Type", "application/msgpack")

		actual := &api.StatusReply{}
		err := mime.Bind(w, r, actual)
		require.NoError(t, err, "should bind MsgPack without error")
		require.Equal(t, expected, actual, "bound object should match expected")
	})

	t.Run("Text", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Set("Content-Type", "text/plain")

		err := mime.Bind(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[415] unsupported content type for request body")
	})

	t.Run("Empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Del("Content-Type")

		err := mime.Bind(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[415] mime: no media type")
	})

	t.Run("Unknown", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Set("Content-Type", "application/foo")

		err := mime.Bind(w, r, &api.StatusReply{})
		require.EqualError(t, err, `[415] unknown mediatype "application/foo"`)
	})
}

func TestBindJSON(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.json"))
		r.Header.Set("Content-Type", "application/json")

		actual := &api.StatusReply{}
		err := mime.BindJSON(w, r, actual)
		require.NoError(t, err, "should bind JSON without error")
		require.Equal(t, expected, actual, "bound object should match expected")
	})

	t.Run("WrongType", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.msgpack"))
		r.Header.Set("Content-Type", "application/msgpack")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, `[415] content type application/json required for this endpoint`)
	})

	t.Run("NoType", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.json"))
		r.Header.Set("Content-Type", "")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.EqualError(t, err, `[415] mime: no media type`)
	})

	t.Run("BadJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString("this is not valid json"))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] request body contains badly-formed JSON (at position 2)")
	})

	t.Run("Truncated", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"status": "ok"`))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] request body contains badly-formed JSON")
	})

	t.Run("JSONType", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"status": 1.23}`))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] request body contains an invalid value for field \"status\" at position 15")
	})

	t.Run("UnknownField", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"status": "ok", "foo": "bar"}`))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] request body contains unknown field \"foo\"")
	})

	t.Run("Empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] no data in request body")
	})

	t.Run("TooLarge", func(t *testing.T) {
		large := make([]byte, mime.MaxPayloadSize+100)
		rand.Read(large)

		status := &api.StatusReply{
			Status: base64.StdEncoding.EncodeToString(large),
		}
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		require.NoError(t, encoder.Encode(status), "should encode large status without error")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", &buf)
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[413] maximum size limit exceeded")
	})

	t.Run("Multiple", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"status": "ok"}{"status": "again"}`))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindJSON(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] request body must contain a single JSON object")
	})
}

func TestBindMsgPack(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.msgpack"))
		r.Header.Set("Content-Type", "application/msgpack")

		actual := &api.StatusReply{}
		err := mime.BindMsgPack(w, r, actual)
		require.NoError(t, err, "should bind MsgPack without error")
		require.Equal(t, expected, actual, "bound object should match expected")
	})

	t.Run("WrongType", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.json"))
		r.Header.Set("Content-Type", "application/json")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.EqualError(t, err, `[415] content type application/msgpack required for this endpoint`)
	})

	t.Run("NoType", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", readFixture(t, "data.msgpack"))
		r.Header.Set("Content-Type", "")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.EqualError(t, err, `[415] mime: no media type`)
	})

	t.Run("Decodable", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Set("Content-Type", "application/msgpack")

		err := mime.BindMsgPack(w, r, &struct{}{})
		require.EqualError(t, err, "destination does not implement msgp.Decodable")
	})

	t.Run("Empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", nil)
		r.Header.Set("Content-Type", "application/msgpack")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[400] no data in request body")
	})

	t.Run("TooLarge", func(t *testing.T) {
		large := make([]byte, mime.MaxPayloadSize+100)
		rand.Read(large)

		status := &api.StatusReply{
			Status: base64.StdEncoding.EncodeToString(large),
		}
		var buf bytes.Buffer
		encoder := msgp.NewWriter(&buf)
		require.NoError(t, status.EncodeMsg(encoder), "should encode large status without error")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", &buf)
		r.Header.Set("Content-Type", "application/msgpack")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.EqualError(t, err, "[413] maximum size limit exceeded")
	})

	t.Run("BadData", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString("this is not valid msgpack"))
		r.Header.Set("Content-Type", "application/msgpack")

		err := mime.BindMsgPack(w, r, &api.StatusReply{})
		require.Error(t, err, "should return error for invalid msgpack data")
	})
}

// TestMain is the entry point for testing in this package. It generates the test data
// if it doesn't already exist. To regenerate testdata, simply delete all of the files
// from the testdata directory.
func TestMain(m *testing.M) {
	fixtures := map[string]func(path string) error{
		"data.json":    generateJSON,
		"data.msgpack": generateMsgPack,
	}

	for name, gen := range fixtures {
		path := filepath.Join("testdata", name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := gen(path); err != nil {
				fmt.Fprintf(os.Stderr, "could not generate fixture %s: %v\n", name, err)
				os.Exit(7)
			}
		}
	}

	// Run all the tests in the package
	os.Exit(m.Run())
}

var expected = &api.StatusReply{
	Status:  "testing",
	Uptime:  (time.Duration(12345) * time.Millisecond).String(),
	Version: "v1.2.3-test.1+abcdefg",
}

func readFixture(t *testing.T, path string) io.Reader {
	t.Helper()
	path = filepath.Join("testdata", path)
	f, err := os.Open(path)
	require.NoError(t, err, "could not open fixture %s", path)
	t.Cleanup(func() { f.Close() })
	return f
}

func generateJSON(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err = encoder.Encode(expected); err != nil {
		return err
	}
	return nil
}

func generateMsgPack(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	encoder := msgp.NewWriter(f)
	if err = expected.EncodeMsg(encoder); err != nil {
		return err
	}

	if err = encoder.Flush(); err != nil {
		return err
	}
	return nil
}
