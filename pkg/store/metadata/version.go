package metadata

import (
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

var _ lani.Encodable = (*Version)(nil)
var _ lani.Decodable = (*Version)(nil)

// The static size of a zero valued Version object; see TestVersionSize for details.
const versionStaticSize = 12

func (o *Version) Size() (s int) {
	s = versionStaticSize

	// Scalar is technically static, but we add its size here in case there are changes
	// to the Scalar struct in the future (e.g. a promotion of integers).
	s += o.Scalar.Size()

	if o.Parent != nil {
		s += o.Parent.Size()
	}
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
