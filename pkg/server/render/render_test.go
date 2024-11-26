package render_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rotationalio/honu/pkg/api/v1"
	"github.com/rotationalio/honu/pkg/mime"
	"github.com/rotationalio/honu/pkg/server/render"
	"github.com/stretchr/testify/require"
)

func TestText(t *testing.T) {
	w := httptest.NewRecorder()
	err := render.Text(http.StatusOK, w, "hello world")
	require.NoError(t, err, "could not render text")

	rep := w.Result()
	require.Equal(t, http.StatusOK, rep.StatusCode, "unexpected status code")
	require.Equal(t, mime.TEXT.ContentType(), rep.Header.Get(render.ContentType), "unexpected content type")

	body, _ := io.ReadAll(rep.Body)
	require.Equal(t, "hello world\n", string(body))
}

func TestTextf(t *testing.T) {
	w := httptest.NewRecorder()
	err := render.Textf(http.StatusAccepted, w, "hello %s", "bobby")
	require.NoError(t, err, "could not render text")

	rep := w.Result()
	require.Equal(t, http.StatusAccepted, rep.StatusCode, "unexpected status code")
	require.Equal(t, mime.TEXT.ContentType(), rep.Header.Get(render.ContentType), "unexpected content type")

	body, _ := io.ReadAll(rep.Body)
	require.Equal(t, "hello bobby", string(body))
}

func TestMsgPack(t *testing.T) {
	w := httptest.NewRecorder()
	out := &api.StatusReply{Status: "foo"}
	err := render.MsgPack(http.StatusCreated, w, out)
	require.NoError(t, err, "could not render msgpack")

	rep := w.Result()
	require.Equal(t, http.StatusCreated, rep.StatusCode, "unexpected status code")
	require.Equal(t, mime.MSGPACK.ContentType(), rep.Header.Get(render.ContentType), "unexpected content type")

	body, _ := io.ReadAll(rep.Body)
	expected, _ := out.MarshalMsg(nil)
	require.Equal(t, expected, body)
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	out := &api.StatusReply{Status: "foo"}
	err := render.JSON(http.StatusEarlyHints, w, out)
	require.NoError(t, err, "could not render json")

	rep := w.Result()
	require.Equal(t, http.StatusEarlyHints, rep.StatusCode, "unexpected status code")
	require.Equal(t, mime.JSON.ContentType(), rep.Header.Get(render.ContentType), "unexpected content type")

	body, _ := io.ReadAll(rep.Body)
	expected, _ := json.Marshal(out)
	expected = append(expected, '\n')
	require.Equal(t, expected, body)
}
