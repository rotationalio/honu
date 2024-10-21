package store

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/oklog/ulid"
)

const (
	maxInt          = int(^uint(0) >> 1)
	smallBufferSize = 64
)

// Marshal an encodable object into a byte slice for storage or serialization.
func Marshal(v Encodable) (_ []byte, err error) {
	encoder := &Encoder{}
	encoder.Grow(v.Size())
	if _, err = v.Encode(encoder); err != nil {
		return nil, err
	}
	return encoder.Bytes(), nil
}

// All objects in this package must be encodable.
type Encodable interface {
	Size() int
	Encode(*Encoder) (int, error)
}

// Encoder is similar to a bytes.Buffer in that it maintains an internal buffer that
// it keeps at the largest capacity its seen for repeated encodings. To use the
// Encoder, sequentially encode values or objects to the buffer then read the encoded
// data using the Bytes() method.
type Encoder struct {
	buf []byte
	off int
}

//===========================================================================
// Get Encoded Data
//===========================================================================

// Bytes returns the data encoded so far and advances the offset so that the same
// buffer can be reused for additional encoding in the future.
func (e *Encoder) Bytes() []byte {
	if e.empty() {
		// If empty, reset to recover space
		e.Reset()
		return nil
	}

	// Return the byte slice currently available and increment the offset as
	// the last encoded object has been read off of the encoder.
	out := e.buf[e.off:]
	e.off += len(out)
	return out
}

//===========================================================================
// Encoder Write Methods
//===========================================================================

// Grow the encoder's write capacity, if necessary, to guarantee space for
// another n bytes. Grow should be called with an Encodable.Size() to ensure
// the entire encodable can be written without additional allocations.
func (e *Encoder) Grow(n int) {
	if n < 0 {
		panic(ErrGrowNegative)
	}
	m := e.grow(n)
	e.buf = e.buf[:m]
}

// Encode a frame to the underlying byte slice consisting of the length of the
// byte array and the specified byte array to be written.
func (e *Encoder) Encode(p []byte) (int, error) {
	// Ensure we have capacity to write the frame (length + bytes)
	size := len(p) + binary.MaxVarintLen64
	m, ok := e.tryGrowByReslice(size)
	if !ok {
		m = e.grow(size)
	}

	// Write the length portion of the frame
	i := binary.PutUvarint(e.buf[m:], uint64(len(p)))

	// Copy in the data and truncate any extra capacity
	n := copy(e.buf[m+i:], p)
	e.buf = e.buf[:m+n+i]
	return n + i, nil
}

// Encode a string frame with length to the underlying byte slice (see Encode).
// Shortcut for Encode([]byte(s))
func (e *Encoder) EncodeString(s string) (int, error) {
	return e.Encode([]byte(s))
}

// Encode a string slice as a multi-frame where the first component is the length
// of the slice and the following components are individual string frames.
func (e *Encoder) EncodeStringSlice(strs []string) (b int, err error) {
	// Compute the size needed to encode the string array
	size := (len(strs) + 1) * binary.MaxVarintLen64
	for _, s := range strs {
		size += len([]byte(s))
	}

	m, ok := e.tryGrowByReslice(size)
	if !ok {
		m = e.grow(size)
	}

	// Write the number of elements in the array
	i := binary.PutUvarint(e.buf[m:], uint64(len(strs)))
	b += i

	// Reset the write index and the buffer to include the array length
	m += i

	// Encode the strings in the string slice
	for _, s := range strs {
		// Encode the string frame length
		sb := []byte(s)
		j := binary.PutUvarint(e.buf[m:], uint64(len(sb)))

		// Copy in the string and truncate extra
		n := copy(e.buf[m+j:], sb)
		m += n + j
		b += n + j
	}

	e.buf = e.buf[:m]
	return
}

// Encode a single byte by writing only that byte to the underlying buffer without
// an array. We expect the decoder to understand it is only reading a single byte.
func (e *Encoder) EncodeByte(c byte) (int, error) {
	m, ok := e.tryGrowByReslice(1)
	if !ok {
		m = e.grow(1)
	}
	e.buf[m] = c
	return 1, nil
}

// Encode a bool as a single byte (1 for true, 0 for false).
func (e *Encoder) EncodeBool(b bool) (int, error) {
	if b {
		return e.EncodeByte(0x1)
	}
	return e.EncodeByte(0x0)
}

// Encode a uint8 to the underlying buffer without a frame.
// Primarily used for short enums and bools.
// Shortcut for EncodeByte(byte(c))
func (e *Encoder) EncodeUint8(c uint8) (int, error) {
	return e.EncodeByte(byte(c))
}

// Encode a Uint32 to the underlying buffer.
func (e *Encoder) EncodeUint32(n uint32) (int, error) {
	m, ok := e.tryGrowByReslice(binary.MaxVarintLen32)
	if !ok {
		m = e.grow(binary.MaxVarintLen32)
	}

	i := binary.PutUvarint(e.buf[m:], uint64(n))
	e.buf = e.buf[:m+i]
	return i, nil
}

// Encode a Uint64 to the underlying buffer.
func (e *Encoder) EncodeUint64(n uint64) (int, error) {
	m, ok := e.tryGrowByReslice(binary.MaxVarintLen64)
	if !ok {
		m = e.grow(binary.MaxVarintLen64)
	}

	i := binary.PutUvarint(e.buf[m:], n)
	e.buf = e.buf[:m+i]
	return i, nil
}

func (e *Encoder) EncodeInt64(n int64) (int, error) {
	m, ok := e.tryGrowByReslice(binary.MaxVarintLen64)
	if !ok {
		m = e.grow(binary.MaxVarintLen64)
	}

	i := binary.PutVarint(e.buf[m:], n)
	e.buf = e.buf[:m+i]
	return i, nil
}

// Encode a fixed number of bytes without a frame to the underlying buffer. We
// expect the decoder to understand exactly how many bytes it needs to read.
func (e *Encoder) EncodeFixed(p []byte) (int, error) {
	m, ok := e.tryGrowByReslice(len(p))
	if !ok {
		m = e.grow(len(p))
	}
	return copy(e.buf[m:], p), nil
}

// Encode a ULID which is always 16 bytes to the frame.
// Shortcut for EncodeFixed(u[:])
func (e *Encoder) EncodeULID(u ulid.ULID) (int, error) {
	return e.EncodeFixed(u[:])
}

// Encode a timestamp as a unix epoch int64 with nanoseconds resolution.
// Shortcut for EncodeInt64(t.UnixNano())
func (e *Encoder) EncodeTime(t time.Time) (int, error) {
	return e.EncodeInt64(t.UnixNano())
}

// Encode a struct (must implement the Encodable interface) as a bool indicating
// if the struct is nil or not with the raw data of the struct following.
func (e *Encoder) EncodeStruct(s Encodable) (n int, err error) {
	var m int
	isNil := isNilEncodable(s)
	if m, err = e.EncodeBool(!isNil); err != nil {
		return n + m, err
	}
	n += m

	if !isNil {
		if m, err = s.Encode(e); err != nil {
			return n + m, err
		}
		n += m
	}

	return
}

func isNilEncodable(s Encodable) bool {
	switch t := s.(type) {
	case *Object:
		return t == nil
	case *Version:
		return t == nil
	case *SchemaVersion:
		return t == nil
	case *AccessControl:
		return t == nil
	case *Publisher:
		return t == nil
	case *Encryption:
		return t == nil
	case *Compression:
		return t == nil
	default:
		panic(fmt.Errorf("unknown type %T", t))
	}
}

//===========================================================================
// Encoder Information Methods
//===========================================================================

// Len returns the number of bytes of the unread portion of the encoder.
func (e *Encoder) Len() int {
	return len(e.buf) - e.off
}

// Cap returns the capacity of the encoder's underlying byte slice.
func (e *Encoder) Cap() int {
	return cap(e.buf)
}

// Resets the buffer to be empty, but retains the underlying storage.
func (e *Encoder) Reset() {
	e.buf = e.buf[:0]
	e.off = 0
}

// empty reports whether the unread portion of the encoder is empty.
func (e *Encoder) empty() bool { return len(e.buf) <= e.off }

//===========================================================================
// Growing the Encoder Buffer
//===========================================================================

// Grow the buffer to guarantee space for n more bytes.
// If the buffer can't grow it will panic with ErrTooLarge.
func (e *Encoder) grow(n int) int {
	m := e.Len()

	// If empty, reset to recover space
	if m == 0 && e.off != 0 {
		e.Reset()
	}

	// Try to grow by means of a reslice.
	if i, ok := e.tryGrowByReslice(n); ok {
		return i
	}

	// Allocate new underlying storage with a minimum buffer size.
	if e.buf == nil && n <= smallBufferSize {
		e.buf = make([]byte, n, smallBufferSize)
		return 0
	}

	// Find the space or allocate new space if necessary
	c := cap(e.buf)
	if n <= c/2-m {
		// Try to slide data down instead of allocating a new slice
		copy(e.buf, e.buf[e.off:])
	} else if c > maxInt-c-n {
		// If we can't grow the slice then we need to panic
		panic(ErrTooLarge)
	} else {
		// Allocate just enough for the n bytes we need
		e.buf = growSlice(e.buf[e.off:], e.off+n)
	}

	// Restore offset and length
	e.off = 0
	e.buf = e.buf[:m+n]
	return m
}

// An inlineable version of grow for the fast-case where the internal
// buffer only needs to be resliced.
func (e *Encoder) tryGrowByReslice(n int) (int, bool) {
	if l := len(e.buf); n <= cap(e.buf)-l {
		e.buf = e.buf[:l+n]
		return l, true
	}
	return 0, false
}

// growSlice grows b by n, preserving the original content of b but not doubling the
// capacity as the buffer in the standard library does.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()

	buf := append([]byte(nil), make([]byte, len(b)+n)...)
	copy(buf, b)
	return buf[:len(b)]
}
