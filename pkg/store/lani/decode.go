package lani

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

// Decoder is similar to a bytes.Reader, allowing a sequential decoding of byte frames,
// so that every repeated decoding advances the internal index to the next frame.
type Decoder struct {
	buf []byte
	i   int
}

// Create a new decoder to decode the given byte array.
func NewDecoder(data []byte) *Decoder {
	return &Decoder{buf: data, i: 0}
}

//===========================================================================
// Decoder Read Methods
//===========================================================================

// Decode a frame from the underlying byte slice by first reading the length of
// bytes from the frame then returning a byte array of that length. Returns an
// error if the underlying data does not contain a frame representation.
func (d *Decoder) Decode() (out []byte, err error) {
	var rl int
	if rl, err = d.readLength(); err != nil {
		return nil, err
	}

	// If there is no following data return a nil slice.
	if rl == 0 {
		return nil, nil
	}

	// Find the stop index:
	j := d.i + rl

	// Ensure we've got enough data for the read
	if j > len(d.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	// Read the data
	out = make([]byte, rl)
	copy(out, d.buf[d.i:j])

	// Advance the read index
	d.i += rl
	return out, nil
}

// Decode a string frame from the underlying byte slice as a frame (see Decode).
func (d *Decoder) DecodeString() (_ string, err error) {
	var out []byte
	if out, err = d.Decode(); err != nil {
		return "", err
	}
	return string(out), nil
}

// Decode a string slice from the underlying buffer where the string slice is a
// multi-frame structure. The first element is the number of elements in the
// array, and then each array after that is a string frame.
func (d *Decoder) DecodeStringSlice() (out []string, err error) {
	// Find the number of elements in the string slice.
	var nstrs int
	if nstrs, err = d.readLength(); err != nil {
		return nil, err
	}

	// If the array is empty return nil
	if nstrs == 0 {
		return nil, nil
	}

	// Begin decoding the string slice
	out = make([]string, nstrs)
	for i := 0; i < nstrs; i++ {
		if out[i], err = d.DecodeString(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// Decode a single byte from the buffer
func (d *Decoder) DecodeByte() (byte, error) {
	if d.i >= len(d.buf) {
		return 0x0, io.EOF
	}

	c := d.buf[d.i]
	d.i += 1
	return c, nil
}

// Decode a bool from the next byte in the buffer
func (d *Decoder) DecodeBool() (_ bool, err error) {
	var c byte
	if c, err = d.DecodeByte(); err != nil {
		return false, err
	}

	switch c {
	case 0x0:
		return false, nil
	case 0x1:
		return true, nil
	default:
		return false, ErrParseBoolean
	}
}

// Decode a uint8 from the next byte in the buffer
func (d *Decoder) DecodeUint8() (uint8, error) {
	return d.DecodeByte()
}

// Decode a uint32 from the underlying buffer
func (d *Decoder) DecodeUint32() (uint32, error) {
	if d.i >= len(d.buf) {
		return 0, io.EOF
	}

	// Find the stop index without going over the length of the buffer
	j := d.i + binary.MaxVarintLen32
	if j > len(d.buf) {
		j = len(d.buf)
	}

	// Decode the Uvarint
	v, k := binary.Uvarint(d.buf[d.i:j])
	if k <= 0 {
		return 0, ErrParseVarInt
	}

	d.i += k
	return uint32(v), nil
}

// Decode a uint64 from the underlying buffer
func (d *Decoder) DecodeUint64() (uint64, error) {
	if d.i >= len(d.buf) {
		return 0, io.EOF
	}

	// Find the stop index without going over the length of the buffer
	j := d.i + binary.MaxVarintLen64
	if j > len(d.buf) {
		j = len(d.buf)
	}

	// Decode the Uvarint
	v, k := binary.Uvarint(d.buf[d.i:j])
	if k <= 0 {
		return 0, ErrParseVarInt
	}

	d.i += k
	return v, nil
}

// Decode an int64 from the underlying buffer
func (d *Decoder) DecodeInt64() (int64, error) {
	if d.i >= len(d.buf) {
		return 0, io.EOF
	}

	// Find the stop index without going over the length of the buffer
	j := d.i + binary.MaxVarintLen64
	if j > len(d.buf) {
		j = len(d.buf)
	}

	// Decode the Uvarint
	v, k := binary.Varint(d.buf[d.i:j])
	if k <= 0 {
		return 0, ErrParseVarInt
	}

	d.i += k
	return v, nil
}

// Decode n fixed bytes from the buffer without using a frame
func (d *Decoder) DecodeFixed(n int) ([]byte, error) {
	if d.i >= len(d.buf) {
		return nil, io.EOF
	}

	j := d.i + n
	if j > len(d.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	out := d.buf[d.i:j]
	d.i += n
	return out, nil
}

// Decode a ULID as 16 fixed bytes from the underlying buffer.
func (d *Decoder) DecodeULID() (ulid.ULID, error) {
	if d.i >= len(d.buf) {
		return ulid.ULID{}, io.EOF
	}

	if d.i+16 > len(d.buf) {
		return ulid.ULID{}, io.ErrUnexpectedEOF
	}

	out := ulid.ULID(d.buf[d.i : d.i+16])
	d.i += 16
	return out, nil
}

// Decode a timestamp from the underlying buffer.
func (d *Decoder) DecodeTime() (_ time.Time, err error) {
	var ts int64
	if ts, err = d.DecodeInt64(); err != nil {
		return time.Time{}, err
	}

	// Handle zero-valued time
	if ts == 0 {
		return time.Time{}, nil
	}

	return time.Unix(0, ts), nil
}

// Decode a struct (must implement the Decodable interface). This function performs
// a nil check and returns if the underlying data is nil or not; otherwise it decodes
// the struct into s.
func (d *Decoder) DecodeStruct(s Decodable) (isNil bool, err error) {
	var notNil bool
	if notNil, err = d.DecodeBool(); err != nil {
		return false, err
	}

	if notNil {
		if err = s.Decode(d); err != nil {
			return !notNil, err
		}
	}

	return !notNil, nil
}

//===========================================================================
// Internal Decoder Methods
//===========================================================================

func (d *Decoder) readLength() (int, error) {
	// Check to make sure we haven't reached the read limit.
	if d.i >= len(d.buf) {
		return 0, io.EOF
	}

	// Find the stop index without going over the length of the buffer
	j := d.i + binary.MaxVarintLen64
	if j > len(d.buf) {
		j = len(d.buf)
	}

	// get the read length and the index that we read to.
	rl, k := binary.Uvarint(d.buf[d.i:j])
	if k <= 0 {
		return -1, ErrNoLength
	}

	// Advance the read index and return the read length
	d.i += k
	return int(rl), nil
}
