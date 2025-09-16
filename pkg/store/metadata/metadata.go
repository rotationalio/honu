package metadata

import (
	"encoding/binary"
	"time"

	"go.rtnl.ai/honu/pkg/region"
	"go.rtnl.ai/honu/pkg/store/keys"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

//===========================================================================
// Metadata Storage
//===========================================================================

type Metadata struct {
	ObjectID     ulid.ULID        `json:"oid" msg:"oid"`
	CollectionID ulid.ULID        `json:"collection" msg:"collection"`
	Version      *Version         `json:"version" msg:"version"`
	Schema       *SchemaVersion   `json:"schema,omitempty" msg:"schema,omitempty"`
	MIME         string           `json:"mime" msg:"mime"`
	Owner        ulid.ULID        `json:"owner" msg:"owner"`
	Group        ulid.ULID        `json:"group" msg:"group"`
	Permissions  uint8            `json:"permissions" msg:"permissions"`
	ACL          []*AccessControl `json:"acl,omitempty" msg:"acl,omitempty"`
	WriteRegions region.Regions   `json:"write_regions,omitempty" msg:"write_regions,omitempty"`
	Publisher    *Publisher       `json:"publisher,omitempty" msg:"publisher,omitempty"`
	Encryption   *Encryption      `json:"encryption,omitempty" msg:"encryption,omitempty"`
	Compression  *Compression     `json:"compression,omitempty" msg:"compression,omitempty"`
	Flags        uint8            `json:"flags" msg:"flags"`
	Created      time.Time        `json:"created" msg:"created"`
	Modified     time.Time        `json:"modified" msg:"modified"`
	key          keys.Key         `json:"-" msg:"-"`
}

//===========================================================================
// Metadata Helper Methods
//===========================================================================

func (m *Metadata) IsTombstone() bool {
	return m.Version != nil && m.Version.Tombstone
}

//===========================================================================
// Metadata Serialization
//===========================================================================

var _ lani.Encodable = (*Metadata)(nil)
var _ lani.Decodable = (*Metadata)(nil)

func (o *Metadata) Key() keys.Key {
	if o.key == nil {
		o.key = keys.New(o.ObjectID, &o.Version.Scalar)
	}
	return o.key
}

// The static size of a zero valued Metadata object; see TestMetadataSize for details.
const metadataStaticSize = 121

func (o *Metadata) Size() (s int) {
	s = metadataStaticSize

	// Version size
	if o.Version != nil {
		s += o.Version.Size()
	}

	// SchemaVersion size
	if o.Schema != nil {
		s += o.Schema.Size()
	}

	// Length of MIME string
	s += len([]byte(o.MIME))

	// ACL List
	s += len(o.ACL) * binary.MaxVarintLen64
	for _, ac := range o.ACL {
		s += 1
		if ac != nil {
			s += ac.Size()
		}
	}

	// Write Regions List
	s += binary.MaxVarintLen64
	s += len(o.WriteRegions) * binary.MaxVarintLen32

	// Publisher size
	if o.Publisher != nil {
		s += o.Publisher.Size()
	}

	// Encryption size
	if o.Encryption != nil {
		s += o.Encryption.Size()
	}

	// Compression size
	if o.Compression != nil {
		s += o.Compression.Size()
	}

	return
}

func (o *Metadata) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(o.ObjectID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.CollectionID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Schema); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.MIME); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.Owner); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.Group); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Permissions); err != nil {
		return n + m, err
	}
	n += m

	// Encode ACL length
	if m, err = e.EncodeUint64(uint64(len(o.ACL))); err != nil {
		return n + m, err
	}
	n += m

	// Encode each ACL
	for _, ac := range o.ACL {
		if m, err = e.EncodeStruct(ac); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = o.WriteRegions.Encode(e); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Publisher); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Encryption); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(o.Compression); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(o.Flags); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Created); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(o.Modified); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Metadata) Decode(d *lani.Decoder) (err error) {
	// Setup nested structs
	o.Version = &Version{}
	o.Schema = &SchemaVersion{}
	o.Publisher = &Publisher{}
	o.Encryption = &Encryption{}
	o.Compression = &Compression{}

	if o.ObjectID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.CollectionID, err = d.DecodeULID(); err != nil {
		return err
	}

	var isNil bool
	if isNil, err = d.DecodeStruct(o.Version); err != nil {
		return err
	} else if isNil {
		o.Version = nil
	}

	if isNil, err = d.DecodeStruct(o.Schema); err != nil {
		return err
	} else if isNil {
		o.Schema = nil
	}

	if o.MIME, err = d.DecodeString(); err != nil {
		return err
	}

	if o.Owner, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Group, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Permissions, err = d.DecodeUint8(); err != nil {
		return err
	}

	// Read the number of ACLs stored
	var nACLs uint64
	if nACLs, err = d.DecodeUint64(); err != nil {
		return err
	}

	// Decode all ACLs
	if nACLs > 0 {
		o.ACL = make([]*AccessControl, nACLs)
		for i := uint64(0); i < nACLs; i++ {
			o.ACL[i] = &AccessControl{}
			if isNil, err = d.DecodeStruct(o.ACL[i]); err != nil {
				return err
			}
			if isNil {
				o.ACL[i] = nil
			}
		}
	}

	if err = o.WriteRegions.Decode(d); err != nil {
		return err
	}

	if isNil, err = d.DecodeStruct(o.Publisher); err != nil {
		return err
	} else if isNil {
		o.Publisher = nil
	}

	if isNil, err = d.DecodeStruct(o.Encryption); err != nil {
		return err
	} else if isNil {
		o.Encryption = nil
	}

	if isNil, err = d.DecodeStruct(o.Compression); err != nil {
		return err
	} else if isNil {
		o.Compression = nil
	}

	if o.Flags, err = d.DecodeUint8(); err != nil {
		return err
	}

	if o.Created, err = d.DecodeTime(); err != nil {
		return err
	}

	if o.Modified, err = d.DecodeTime(); err != nil {
		return err
	}

	return nil
}
