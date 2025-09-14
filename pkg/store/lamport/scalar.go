package lamport

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"go.rtnl.ai/honu/pkg/store/lani"
)

// A lamport scalar indicates a timestamp or version in a vector clock that provides a
// "happens before" relationship between two scalars and is used to create distributed
// version numbers for eventually consistent systems and uses a "latest writer wins"
// policy to determine how to solve conflicts.
//
// Each scalar is 12 bytes and is composed of an 4 byte PID which should be unique to
// every process in the system, and an 8 byte monotonically increasing VID that
// represents the next latest version. In the case of a scalar with the same VID, the
// scalar with the smaller PID happens before the scalar with the larger PID (e.g.
// newer processes, processes with bigger PIDs, win ties).
type Scalar struct {
	PID uint32
	VID uint64
}

var (
	zero                            = &Scalar{0, 0}
	scre                            = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	_    lani.Encodable             = (*Scalar)(nil)
	_    lani.Decodable             = (*Scalar)(nil)
	_    encoding.BinaryMarshaler   = (*Scalar)(nil)
	_    encoding.BinaryUnmarshaler = (*Scalar)(nil)
	_    encoding.TextMarshaler     = (*Scalar)(nil)
	_    encoding.TextUnmarshaler   = (*Scalar)(nil)
)

const scalarSize = binary.MaxVarintLen32 + binary.MaxVarintLen64

var (
	ErrInvalidFormat = errors.New("could not parse text representation of a scalar")
)

// Compare returns an integer comparing two scalars using a happens before relationship.
// The result will be 0 if a == b, -1 if a < b (e.g. a happens before b), and
// +1 if a > b (e.g. b happens before a). A nil argument is equivalent to a zero scalar.
func Compare(a, b *Scalar) int {
	if a == nil && b == nil {
		return 0
	}

	if a == nil {
		a = zero
	}

	if b == nil {
		b = zero
	}

	if a.VID == b.VID {
		switch {
		case a.PID < b.PID:
			return -1
		case a.PID > b.PID:
			return 1
		default:
			return 0
		}
	}

	if a.VID > b.VID {
		return 1
	}
	return -1
}

// Returns true if the scalar is the zero-valued scalar (0.0)
func (s *Scalar) IsZero() bool {
	return s.PID == 0 && s.VID == 0
}

// Returns true if the scalar is equal to the input scalar.
func (s *Scalar) Equals(o *Scalar) bool {
	return Compare(s, o) == 0
}

// Returns true if the scalar is less than the input scalar (e.g. this scalar happens
// before the input scalar).
func (s *Scalar) Before(o *Scalar) bool {
	return Compare(s, o) < 0
}

// Returns true if the scalar is grater than the input scalar (e.g. the input scalar
// hapens before this scalar).
func (s *Scalar) After(o *Scalar) bool {
	return Compare(s, o) > 0
}

func (s *Scalar) Size() int {
	return scalarSize
}

func (s *Scalar) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint32(s.PID); err != nil {
		return n, err
	}
	n += m

	if m, err = e.EncodeUint64(s.VID); err != nil {
		return n, err
	}
	n += m

	return n, nil
}

func (s *Scalar) Decode(d *lani.Decoder) (err error) {
	if s.PID, err = d.DecodeUint32(); err != nil {
		return err
	}

	if s.VID, err = d.DecodeUint64(); err != nil {
		return err
	}

	return nil
}

func (s *Scalar) MarshalBinary() ([]byte, error) {
	e := &lani.Encoder{}
	e.Grow(s.Size())

	if _, err := s.Encode(e); err != nil {
		return nil, err
	}

	return e.Bytes(), nil
}

func (s *Scalar) UnmarshalBinary(data []byte) (err error) {
	d := lani.NewDecoder(data)
	return s.Decode(d)
}

func (s *Scalar) MarshalText() (_ []byte, err error) {
	return []byte(s.String()), nil
}

func (s *Scalar) UnmarshalText(text []byte) (err error) {
	if !scre.Match(text) {
		return ErrInvalidFormat
	}

	parts := bytes.Split(text, []byte{'.'})

	var pid uint64
	if pid, err = strconv.ParseUint(string(parts[0]), 10, 32); err != nil {
		panic("pid is not parseable even though regular expression matched.")
	}
	s.PID = uint32(pid)

	if s.VID, err = strconv.ParseUint(string(parts[1]), 10, 32); err != nil {
		panic("pid is not parseable even though regular expression matched.")
	}

	return nil
}

// Returns a scalar version representation in the form PID.VID using decimal notation.
func (s *Scalar) String() string {
	return fmt.Sprintf("%d.%d", s.PID, s.VID)
}

//===========================================================================
// Sort Interface
//===========================================================================

type Scalars []*Scalar

func (s Scalars) Len() int           { return len(s) }
func (s Scalars) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s Scalars) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
