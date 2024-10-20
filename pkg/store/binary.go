package store

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"github.com/oklog/ulid"
)

type Writer struct {
	bytes.Buffer
}

var _ io.Writer = &Writer{}

func (w *Writer) WriteBytes(p []byte) (n int, err error) {
	// TODO: reduce the number of allocations by directly writing to a bytes array.
	buf := make([]byte, binary.MaxVarintLen64, binary.MaxVarintLen64+len(p))
	i := binary.PutUvarint(buf, uint64(len(p)))
	buf = buf[:i+len(p)]
	copy(buf[i:], p)
	return w.Write(buf)
}

func (w *Writer) WriteString(s string) (n int, err error) {
	return w.WriteBytes([]byte(s))
}

func (w *Writer) WriteStrings(s []string) (n int, err error) {
	var m int
	if m, err = w.WriteUint64(uint64(len(s))); err != nil {
		return m, err
	}
	n += m

	for _, v := range s {
		if m, err = w.WriteString(v); err != nil {
			return n + m, err
		}
		n += m
	}
	return
}

func (w *Writer) WriteBool(b bool) (n int, err error) {
	if b {
		return w.WriteUint8(1)
	}
	return w.WriteUint8(0)
}

func (w *Writer) WriteUint8(n uint8) (int, error) {
	return w.Write([]byte{n})
}

func (w *Writer) WriteUint32(n uint32) (int, error) {
	buf := make([]byte, binary.MaxVarintLen32)
	i := binary.PutUvarint(buf, uint64(n))
	return w.Write(buf[:i])
}

func (w *Writer) WriteUint64(n uint64) (int, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	i := binary.PutUvarint(buf, n)
	return w.Write(buf[:i])
}

func (w *Writer) WriteInt64(n int64) (int, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	i := binary.PutVarint(buf, n)
	return w.Write(buf[:i])
}

func (w *Writer) WriteULID(u ulid.ULID) (int, error) {
	return w.Write(u[:])
}

func (w *Writer) WriteTime(t time.Time) (int, error) {
	return w.WriteInt64(t.UnixNano())
}

type Reader struct {
	buf []byte
	i   int64
}

var _ io.Reader = &Reader{}

func NewReader(b []byte) *Reader {
	return &Reader{buf: b}
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	n = copy(b, r.buf[r.i:])
	r.i += int64(n)
	return
}

func (r *Reader) ReadBytes() (out []byte, err error) {
	// Read the length of the bytes data
	var rl int64
	if rl, err = r.readLength(); err != nil {
		return nil, err
	}

	if rl == 0 {
		return nil, nil
	}

	// Find stop index
	si := r.i + rl + 1

	// TODO: check the buffer before reading
	out = make([]byte, rl)
	copy(out, r.buf[r.i:si])

	// Advance the read index
	r.i += rl
	return out, nil
}

func (r *Reader) ReadString() (_ string, err error) {
	var out []byte
	if out, err = r.ReadBytes(); err != nil {
		return "", err
	}
	return string(out), nil
}

func (r *Reader) ReadStrings() (out []string, err error) {
	var nstrs uint64
	if nstrs, err = r.ReadUint64(); err != nil {
		return nil, err
	}

	if nstrs == 0 {
		return nil, nil
	}

	out = make([]string, nstrs)
	for i := uint64(0); i < nstrs; i++ {
		if out[i], err = r.ReadString(); err != nil {
			return nil, err
		}
	}

	return
}

func (r *Reader) ReadBool() (_ bool, err error) {
	var v uint8
	if v, err = r.ReadUint8(); err != nil {
		return false, err
	}

	if v == 1 {
		return true, nil
	}

	if v == 0 {
		return false, nil
	}

	return false, ErrParseBoolean
}

func (r *Reader) ReadUint8() (uint8, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	n := uint8(r.buf[r.i])
	r.i += 1
	return n, nil
}

func (r *Reader) ReadUint32() (uint32, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	rl, ln := binary.Uvarint(r.buf[r.i : r.i+binary.MaxVarintLen32])
	if ln <= 0 {
		return 0, ErrNoLength
	}

	r.i += int64(ln)
	return uint32(rl), nil
}

func (r *Reader) ReadUint64() (uint64, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	rl, ln := binary.Uvarint(r.buf[r.i : r.i+binary.MaxVarintLen64])
	if ln <= 0 {
		return 0, ErrNoLength
	}

	r.i += int64(ln)
	return rl, nil
}

func (r *Reader) ReadInt64() (int64, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	rl, ln := binary.Varint(r.buf[r.i : r.i+binary.MaxVarintLen64])
	if ln <= 0 {
		return 0, ErrNoLength
	}

	r.i += int64(ln)
	return rl, nil
}

func (r *Reader) ReadULID() (ulid.ULID, error) {
	if r.i >= int64(len(r.buf)) {
		return ulid.ULID{}, io.EOF
	}

	out := ulid.ULID(r.buf[r.i : r.i+16])
	r.i += 16
	return out, nil
}

func (r *Reader) ReadTime() (_ time.Time, err error) {
	var ts int64
	if ts, err = r.ReadInt64(); err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, ts), nil
}

func (r *Reader) readLength() (int64, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	rl, ln := binary.Uvarint(r.buf[r.i : r.i+binary.MaxVarintLen64])
	if ln <= 0 {
		return -1, ErrNoLength
	}

	// Advance r.i and return the read length
	r.i += int64(ln)
	return int64(rl), nil
}
