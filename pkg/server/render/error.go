package render

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.rtnl.ai/honu/pkg/api/v1"
	"go.rtnl.ai/honu/pkg/mime"
)

// Error replies to the request with the specified error message. If the error is a
// StatusError then the status code will be used, otherwise
// http.StatusInternalServerError will be used. It does not otherwise end the request;
// the caller should ensure no further writes are done to w.
//
// While http.Error ensures that the response is plain text, this function will attempt
// to negotiate the content type of the response based on the request's Accept header.
// If there is not acceptable content type, then it will fall back to plain text.
//
// Error deletes the Content-Length header, and sets the Content-Type. It also sets
// the X-Content-Type-Options header to "nosniff". This configures the header properly
// for the error message in case the caller had set it up expecting a successful output.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del(ContentLength)

	// Set the X-Content-Type-Options header to nosniff to prevent the browser
	// from MIME-sniffing a response away from the declared Content-Type.
	h.Set(XContentTypeOptions, "nosniff")

	// Determine the status code to use from the error type.
	statusError := &api.StatusError{}
	if !errors.As(err, &statusError) {
		statusError.Code = http.StatusInternalServerError
		statusError.Reply = api.Error(err)
	}

	// Negotiate the content type to use for the error response.
	// However, do not use the negotiated renderer directly, to ensure that our error
	// is written and not an not acceptable or parsing error.
	renderer := Negotiate(r).(*negotiate)
	switch renderer.contentType {
	case mime.MSGPACK, mime.OCTET_STREAM:
		err = MsgPack(statusError.Code, w, &statusError.Reply)
	case mime.JSON, mime.ANY:
		err = JSON(statusError.Code, w, &statusError.Reply)
	default:
		Text(statusError.Code, w, statusError.Reply.Error)
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("could not write error response")
		Text(statusError.Code, w, statusError.Reply.Error)
	}
}
