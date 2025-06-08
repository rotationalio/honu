package metadata

import (
	"go.rtnl.ai/honu/pkg/store/lani"
)

//===========================================================================
// Schema Version
//===========================================================================

type SchemaVersion struct {
	Name  string `json:"name" msg:"name"`
	Major uint32 `json:"major" msg:"major"`
	Minor uint32 `json:"minor" msg:"minor"`
	Patch uint32 `json:"patch" msg:"patch"`
}

var _ lani.Encodable = &SchemaVersion{}
var _ lani.Decodable = &SchemaVersion{}

// The static size of a zero valued SchemaVersion object; see TestSchemaVersionSize for details.
const schemaVersionStaticSize = 25

func (o *SchemaVersion) Size() (s int) {
	s = schemaVersionStaticSize
	s += len([]byte(o.Name))
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
