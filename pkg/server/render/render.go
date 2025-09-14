package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tinylib/msgp/msgp"
	"go.rtnl.ai/honu/pkg/mime"
)

type Renderer interface {
	Render(code int, w http.ResponseWriter, obj any) error
}

// Header keys for http requests and responses
const (
	Accept      = "Accept"
	ContentType = "Content-Type"
)

// Content type values
var (
	plainContentType   = mime.TEXT.ContentType()
	msgpackContentType = mime.MSGPACK.ContentType()
	jsonContentType    = mime.JSON.ContentType()
)

func Text(code int, w http.ResponseWriter, text string) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintln(w, text)
	return nil
}

func Textf(code int, w http.ResponseWriter, text string, a ...any) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintf(w, text, a...)
	return nil
}

// MsgPack marshals the given interface object to the binary MessagePack representation
// and writes it to the response with the correct ContentType
func MsgPack(code int, w http.ResponseWriter, obj msgp.Encodable) error {
	w.Header().Set(ContentType, msgpackContentType)
	w.WriteHeader(code)

	// Create a buffered encodable stream to prevent additional allocations.
	wm := msgp.NewWriter(w)
	if err := obj.EncodeMsg(wm); err != nil {
		return err
	}

	// Ensure the message is flushed to the underlying writer
	return wm.Flush()
}

// JSON marshals the given interface object and writes it with the correct ContentType.
func JSON(code int, w http.ResponseWriter, obj any) error {
	w.Header().Set(ContentType, jsonContentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}
	return nil
}
