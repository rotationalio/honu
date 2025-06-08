package metadata

import (
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Access Control List
//===========================================================================

type AccessControl struct {
	ClientID    ulid.ULID `json:"client_id" msg:"client_id"`
	Permissions uint8     `json:"permissions" msg:"permissions"`
}

var _ lani.Encodable = &AccessControl{}
var _ lani.Decodable = &AccessControl{}

const accessControlStaticSize = 17 // 16 bytes for ULID + 1 byte for permissions

func (o *AccessControl) Size() int {
	return accessControlStaticSize
}

func (o *AccessControl) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(o.ClientID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Permissions); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *AccessControl) Decode(d *lani.Decoder) (err error) {
	if o.ClientID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Permissions, err = d.DecodeUint8(); err != nil {
		return err
	}

	return nil
}
