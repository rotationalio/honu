package metadata

import (
	"encoding/binary"
	"time"

	"github.com/rotationalio/honu/pkg/store/lani"
)

//===========================================================================
// Metadata Version
//===========================================================================

type Version struct {
	PID       uint64    `json:"pid" msg:"pid"`
	Version   uint64    `json:"version" msg:"version"`
	Region    string    `json:"region" msg:"region"`
	Parent    *Version  `json:"parent,omitempty" msg:"parent,omitempty"`
	Tombstone bool      `json:"tombstone,omitempty" msg:"tombstone,omitempty"`
	Created   time.Time `json:"created" msg:"created"`
}

var _ lani.Encodable = &Version{}
var _ lani.Decodable = &Version{}

func (o *Version) Size() (s int) {
	s += 2 * binary.MaxVarintLen64
	s += len([]byte(o.Region)) + binary.MaxVarintLen64

	if o.Parent != nil {
		s += o.Parent.Size() + 1 // Add 1 for the not nil bool
	} else {
		s += 1 // Add 1 for the nil bool
	}

	s += 1                     // Tombstone bool
	s += binary.MaxVarintLen64 // Timestamp int64

	return
}

func (o *Version) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint64(o.PID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint64(o.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.Region); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Parent); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeBool(o.Tombstone); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Created); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Version) Decode(d *lani.Decoder) (err error) {
	if o.PID, err = d.DecodeUint64(); err != nil {
		return err
	}

	if o.Version, err = d.DecodeUint64(); err != nil {
		return err
	}

	if o.Region, err = d.DecodeString(); err != nil {
		return err
	}

	var isNil bool
	o.Parent = &Version{}

	if isNil, err = d.DecodeStruct(o.Parent); err != nil {
		return err
	}

	if isNil {
		o.Parent = nil
	}

	if o.Tombstone, err = d.DecodeBool(); err != nil {
		return err
	}

	if o.Created, err = d.DecodeTime(); err != nil {
		return err
	}

	return nil
}
