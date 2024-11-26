package render

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rotationalio/honu/pkg/api/v1"
	"github.com/rotationalio/honu/pkg/mime"
	"github.com/tinylib/msgp/msgp"
)

const (
	defaultContentType = mime.JSON
)

func Negotiate(r *http.Request) Renderer {
	n := &negotiate{}
	accepted := Accepted(r)

	if len(accepted) == 0 {
		n.contentType = defaultContentType
		return n
	}

	// Find the first accept header that we can parse.
	for _, accept := range accepted {
		var err error
		if n.contentType, err = mime.Parse(accept); err == nil {
			if n.contentType == mime.ANY {
				n.contentType = defaultContentType
			}
			return n
		}
	}

	// Return UNKNOWN to cause a 406 error to be returned.
	return n
}

type negotiate struct {
	contentType mime.MIME
}

func (n *negotiate) Render(code int, w http.ResponseWriter, obj any) error {
	switch n.contentType {
	case mime.MSGPACK, mime.OCTET_STREAM:
		if msg, ok := obj.(msgp.Encodable); ok {
			return MsgPack(code, w, msg)
		}
		return fmt.Errorf("object of type %T cannot be serialized as msgpack", obj)
	case mime.JSON, mime.ANY:
		return JSON(code, w, obj)
	case mime.UNKNOWN:
		// Return a 406 Not Acceptable HTTP code to the user
		return JSON(http.StatusNotAcceptable, w, api.NotAcceptable)
	default:
		return fmt.Errorf("unhandled mimetype %s", n.contentType)
	}
}

func Accepted(r *http.Request) []string {
	if headers, ok := r.Header[Accept]; ok {
		out := make([]string, 0)
		for _, header := range headers {
			out = append(out, ParseAccept(header)...)
		}
		return out
	}
	return nil
}

func ParseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		// Remove any extensions from the mediatype
		if i := strings.IndexByte(part, ';'); i > 0 {
			part = part[:i]
		}
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
