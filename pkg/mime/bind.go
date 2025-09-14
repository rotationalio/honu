package mime

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tinylib/msgp/msgp"
	"go.rtnl.ai/honu/pkg/api/v1"
)

const (
	// Maximum size in bytes for non-document payload: 10MiB
	MaxPayloadSize = 1.049e+7

	// Content Type header key
	ContentType = "Content-Type"
)

// Bind reads the request body into the destination structure, using the Content-Type
// header to determine how to decode the body. If the Content-Type is not set, a 415
// error is returned. The request body is limited to MaxPayloadSize to prevent abuse.
func Bind(w http.ResponseWriter, r *http.Request, dst any) (err error) {
	// Determine the content type of the request sent by the client.
	var mime MIME
	contentType := r.Header.Get(ContentType)
	if mime, err = Parse(contentType); err != nil {
		return &api.StatusError{
			StatusCode: http.StatusUnsupportedMediaType,
			Reply:      api.Error(err),
		}
	}

	// Create a maximum bytes reader to prevent abuse.
	r.Body = http.MaxBytesReader(w, r.Body, MaxPayloadSize)
	defer r.Body.Close()

	switch mime {
	case JSON:
		return bindJSON(r.Body, dst)
	case MSGPACK:
		return bindMsgPack(r.Body, dst)
	default:
		return &api.StatusError{
			StatusCode: http.StatusUnsupportedMediaType,
			Reply:      api.Error("unsupported content type for request body"),
		}
	}
}

// Bind the JSON payload from the request to the destination structure. This method
// requires that the Content-Type header is set to application/json. If the header
// is not set, or if the content type is not application/json, an error is returned.
func BindJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	return bind(JSON, bindJSON, w, r, dst)
}

// Bind the MessagePack payload from the request to the destination structure. This method
// requires that the Content-Type header is set to application/msgpack. If the header
// is not set, or if the content type is not application/msgpack, an error is returned.
func BindMsgPack(w http.ResponseWriter, r *http.Request, dst any) error {
	return bind(MSGPACK, bindMsgPack, w, r, dst)
}

// Universal bind function to make the type-specific binders easier to implement.
func bind(contentType MIME, binder binder, w http.ResponseWriter, r *http.Request, dst any) (err error) {
	// Determine the content type of the request sent by the client.
	var mime MIME
	if mime, err = Parse(r.Header.Get(ContentType)); err != nil {
		return &api.StatusError{
			StatusCode: http.StatusUnsupportedMediaType,
			Reply:      api.Error(err),
		}
	}

	if mime != contentType {
		return &api.StatusError{
			StatusCode: http.StatusUnsupportedMediaType,
			Reply:      api.Error(fmt.Errorf("content type %s required for this endpoint", contentType)),
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxPayloadSize)
	defer r.Body.Close()

	return binder(r.Body, dst)
}

type binder func(io.Reader, any) error

func bindJSON(body io.Reader, dst any) error {
	// Create the JSON decoder
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	// Decode the JSON data and handle any errors
	if err := decoder.Decode(&dst); err != nil {
		var (
			syntaxError *json.SyntaxError
			typeError   *json.UnmarshalTypeError
			maxBytes    *http.MaxBytesError
		)

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error(msg)}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error("request body contains badly-formed JSON")}

		case errors.As(err, &typeError):
			msg := fmt.Sprintf("request body contains an invalid value for field %q at position %d", typeError.Field, typeError.Offset)
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error(msg)}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error(msg)}

		case errors.Is(err, io.EOF):
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error("no data in request body")}

		case errors.As(err, &maxBytes):
			return &api.StatusError{StatusCode: http.StatusRequestEntityTooLarge, Reply: api.Error("maximum size limit exceeded")}

		default:
			return err
		}
	}

	// Ensure the request body only contains a single JSON object
	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error("request body must contain a single JSON object")}
	}
	return nil
}

func bindMsgPack(body io.Reader, dst any) error {
	decodable, ok := dst.(msgp.Decodable)
	if !ok {
		return fmt.Errorf("destination does not implement msgp.Decodable")
	}

	// Create a new MessagePack decoder
	decoder := msgp.NewReader(body)
	if err := decodable.DecodeMsg(decoder); err != nil {
		var maxBytes *http.MaxBytesError

		switch {
		case errors.Is(err, io.EOF):
			return &api.StatusError{StatusCode: http.StatusBadRequest, Reply: api.Error("no data in request body")}

		case errors.As(err, &maxBytes):
			return &api.StatusError{StatusCode: http.StatusRequestEntityTooLarge, Reply: api.Error("maximum size limit exceeded")}

		default:
			return &api.StatusError{
				StatusCode: http.StatusBadRequest,
				Reply:      api.Error(err),
			}
		}
	}

	return nil
}
