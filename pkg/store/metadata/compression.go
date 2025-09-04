package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.rtnl.ai/honu/pkg/store/lani"
)

//===========================================================================
// Compression
//===========================================================================

type CompressionAlgorithm uint8

const (
	None CompressionAlgorithm = iota
	GZIP
	COMPRESS
	DEFLATE
	BROTLI
)

type Compression struct {
	Algorithm CompressionAlgorithm `json:"algorithm" msg:"algorithm"`
	Level     int64                `json:"level,omitempty" msg:"level,omitempty"`
}

var _ lani.Encodable = &Compression{}
var _ lani.Decodable = &Compression{}

// The static size of a zero valued Compression object; see TestCompressionSize for details.
const compressionStaticSize = 11

func (o *Compression) Size() int {
	return compressionStaticSize
}

func (o *Compression) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint8(uint8(o.Algorithm)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeInt64(o.Level); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Compression) Decode(d *lani.Decoder) (err error) {
	var a uint8
	if a, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.Algorithm = CompressionAlgorithm(a)

	if o.Level, err = d.DecodeInt64(); err != nil {
		return err
	}

	return nil
}

func ParseCompressionAlgorithm(s string) (CompressionAlgorithm, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "NONE":
		return None, nil
	case "GZIP":
		return GZIP, nil
	case "COMPRESS":
		return COMPRESS, nil
	case "DEFLATE":
		return DEFLATE, nil
	case "BROTLI":
		return BROTLI, nil
	default:
		return 0, fmt.Errorf("%q is not a valid compression algorithm", s)
	}
}

func (o CompressionAlgorithm) String() string {
	switch o {
	case None:
		return "NONE"
	case GZIP:
		return "GZIP"
	case COMPRESS:
		return "COMPRESS"
	case DEFLATE:
		return "DEFLATE"
	case BROTLI:
		return "BROTLI"
	default:
		return "UNKNOWN"
	}
}

func (o *CompressionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *CompressionAlgorithm) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *o, err = ParseCompressionAlgorithm(alg); err != nil {
		return err
	}
	return nil
}
