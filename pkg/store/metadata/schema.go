package metadata

import (
	"encoding/binary"

	"github.com/rotationalio/honu/pkg/store/lani"
)

//===========================================================================
// Schema Version
//===========================================================================

type SchemaVersion struct {
	Name  string
	Major uint32
	Minor uint32
	Patch uint32
}

var _ lani.Encodable = &SchemaVersion{}
var _ lani.Decodable = &SchemaVersion{}

func (o *SchemaVersion) Size() (s int) {
	s += len([]byte(o.Name)) + binary.MaxVarintLen64
	s += 3 * binary.MaxVarintLen32
	return
}

func (o *SchemaVersion) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeString(o.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Major); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Minor); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint32(o.Patch); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *SchemaVersion) Decode(d *lani.Decoder) (err error) {
	if o.Name, err = d.DecodeString(); err != nil {
		return err
	}

	if o.Major, err = d.DecodeUint32(); err != nil {
		return err
	}

	if o.Minor, err = d.DecodeUint32(); err != nil {
		return err
	}

	if o.Patch, err = d.DecodeUint32(); err != nil {
		return err
	}

	return nil
}
