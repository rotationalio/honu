package mime

import (
	"fmt"
	"mime"
)

const (
	UNKNOWN MIME = iota
	ANY
	OCTET_STREAM
	TEXT
	JSON
	MSGPACK
	HONU
)

const (
	ctTEXT = "text/plain; charset=utf-8"
	ctJSON = "application/json; charset=utf-8"
)

const (
	mtANY          = "*/*"
	mtTEXT         = "text/plain"
	mtOCTET_STREAM = "application/octet-stream"
	mtMSGPACK      = "application/msgpack"
	mtJSON         = "application/json"
	mtHonu         = "application/honu+lani"
)

var mimeValues = map[MIME]string{
	ANY:          mtANY,
	OCTET_STREAM: mtOCTET_STREAM,
	TEXT:         mtTEXT,
	JSON:         mtJSON,
	MSGPACK:      mtMSGPACK,
	HONU:         mtHonu,
}

var mimeNames = map[string]MIME{
	mtANY:          ANY,
	mtOCTET_STREAM: OCTET_STREAM,
	mtTEXT:         TEXT,
	mtJSON:         JSON,
	mtMSGPACK:      MSGPACK,
	mtHonu:         HONU,
}

type MIME uint16

func Parse(s string) (_ MIME, err error) {
	var mediatype string
	if mediatype, _, err = mime.ParseMediaType(s); err != nil {
		return UNKNOWN, err
	}

	if m, ok := mimeNames[mediatype]; ok {
		return m, nil
	}

	return UNKNOWN, fmt.Errorf("unknown mediatype %q", mediatype)
}

func (m MIME) String() string {
	if m == UNKNOWN || m == OCTET_STREAM {
		return mtOCTET_STREAM
	}

	if s, ok := mimeValues[m]; ok {
		return s
	}

	panic(fmt.Errorf("unknown mimetype enum %d", m))
}

func (m MIME) ContentType() string {
	switch m {
	case TEXT:
		return ctTEXT
	case JSON:
		return ctJSON
	default:
		return m.String()
	}
}
