package metadata

import (
	"encoding/binary"
	"time"

	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/lani"
)

//===========================================================================
// Metadata Version
//===========================================================================

type Version struct {
	Scalar    lamport.Scalar  `json:"scalar" msg:"scalar"`
	Region    string          `json:"region" msg:"region"`
	Parent    *lamport.Scalar `json:"parent,omitempty" msg:"parent,omitempty"`
	Tombstone bool            `json:"tombstone,omitempty" msg:"tombstone,omitempty"`
	Created   time.Time       `json:"created" msg:"created"`
}

var _ lani.Encodable = &Version{}
var _ lani.Decodable = &Version{}

func (o *Version) Size() (s int) {
	s += o.Scalar.Size() // Scalar uint32 + uint64
	s += 1               // Add 1 for the parent nil bool

	if o.Parent != nil {
		s += o.Parent.Size()
	}

	s += 1                     // Tombstone bool
	s += binary.MaxVarintLen64 // Timestamp int64

	return
}

func (o *Version) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = o.Scalar.Encode(e); err != nil {
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
	if err = o.Scalar.Decode(d); err != nil {
		return err
	}

	if o.Region, err = d.DecodeString(); err != nil {
		return err
	}

	var isNil bool
	o.Parent = &lamport.Scalar{}

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
