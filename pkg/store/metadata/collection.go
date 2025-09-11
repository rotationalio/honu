package metadata

import (
	"encoding/binary"
	"time"

	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

// Collections are subsets of the Store that allow access to related objects. Each
// object in a collection is prefixed by the collection ID, ensuring that the objects
// are grouped together and can be accessed efficiently.
type Collection struct {
	ID           ulid.ULID        `json:"id" msg:"id"`
	Name         string           `json:"name" msg:"name"`
	Version      *Version         `json:"version" msg:"version"`
	Owner        ulid.ULID        `json:"owner" msg:"owner"`
	Group        ulid.ULID        `json:"group" msg:"group"`
	Permissions  uint8            `json:"permissions" msg:"permissions"`
	ACL          []*AccessControl `json:"acl,omitempty" msg:"acl,omitempty"`
	WriteRegions []string         `json:"write_regions,omitempty" msg:"write_regions,omitempty"`
	Publisher    *Publisher       `json:"publisher,omitempty" msg:"publisher,omitempty"`
	Schema       *SchemaVersion   `json:"schema,omitempty" msg:"schema,omitempty"`
	Encryption   *Encryption      `json:"encryption,omitempty" msg:"encryption,omitempty"`
	Compression  *Compression     `json:"compression,omitempty" msg:"compression,omitempty"`
	Flags        uint8            `json:"flags" msg:"flags"`
	Created      time.Time        `json:"created" msg:"created"`
	Modified     time.Time        `json:"modified" msg:"modified"`
}

var _ lani.Encodable = (*Collection)(nil)
var _ lani.Decodable = (*Collection)(nil)

// Returns the name of the collection if it is set, otherwise returns the ULID.
func (c *Collection) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.ID.String()
}

func (c *Collection) Validate() (err error) {
	if err = ValidateName(c.Name); err != nil {
		return err
	}

	return nil
}

// The static size of a zero valued Collection object; see TestCollectionSize for details.
const collectionStaticSize = 105

func (c *Collection) Size() (s int) {
	s = collectionStaticSize

	// Name length
	s += len([]byte(c.Name))

	// Version size
	if c.Version != nil {
		s += c.Version.Size()
	}

	// ACL List
	s += len(c.ACL) * binary.MaxVarintLen64
	for _, ac := range c.ACL {
		s += 1
		if ac != nil {
			s += ac.Size()
		}
	}

	// Write Regions List
	s += len(c.WriteRegions) * binary.MaxVarintLen64
	for _, wr := range c.WriteRegions {
		s += len([]byte(wr))
	}

	// Publisher size
	if c.Publisher != nil {
		s += c.Publisher.Size()
	}

	// Schema size
	if c.Schema != nil {
		s += c.Schema.Size()
	}

	// Encryption size
	if c.Encryption != nil {
		s += c.Encryption.Size()
	}

	// Compression size
	if c.Compression != nil {
		s += c.Compression.Size()
	}

	return
}

func (c *Collection) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(c.ID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(c.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(c.Version); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(c.Owner); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(c.Group); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(c.Permissions); err != nil {
		return n + m, err
	}
	n += m

	// Encode ACL length
	if m, err = e.EncodeUint64(uint64(len(c.ACL))); err != nil {
		return n + m, err
	}
	n += m

	// Encode each ACL entry
	for _, ac := range c.ACL {
		if m, err = e.EncodeStruct(ac); err != nil {
			return n + m, err
		}
		n += m
	}

	if m, err = e.EncodeStringSlice(c.WriteRegions); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(c.Publisher); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(c.Schema); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(c.Encryption); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeStruct(c.Compression); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(c.Flags); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(c.Created); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeTime(c.Modified); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (c *Collection) Decode(d *lani.Decoder) (err error) {
	// Setup nested structs
	c.Version = &Version{}
	c.Publisher = &Publisher{}
	c.Schema = &SchemaVersion{}
	c.Encryption = &Encryption{}
	c.Compression = &Compression{}

	if c.ID, err = d.DecodeULID(); err != nil {
		return err
	}

	if c.Name, err = d.DecodeString(); err != nil {
		return err
	}

	var isNil bool
	if isNil, err := d.DecodeStruct(c.Version); err != nil {
		return err
	} else if isNil {
		c.Version = nil
	}

	if c.Owner, err = d.DecodeULID(); err != nil {
		return err
	}

	if c.Group, err = d.DecodeULID(); err != nil {
		return err
	}

	if c.Permissions, err = d.DecodeUint8(); err != nil {
		return err
	}

	// Read the number of ACLs stored
	var nACLs uint64
	if nACLs, err = d.DecodeUint64(); err != nil {
		return err
	}

	// Decode all ACLs
	if nACLs > 0 {
		c.ACL = make([]*AccessControl, nACLs)
		for i := uint64(0); i < nACLs; i++ {
			c.ACL[i] = &AccessControl{}
			if isNil, err = d.DecodeStruct(c.ACL[i]); err != nil {
				return err
			}
			if isNil {
				c.ACL[i] = nil
			}
		}
	}

	if c.WriteRegions, err = d.DecodeStringSlice(); err != nil {
		return err
	}

	if isNil, err = d.DecodeStruct(c.Publisher); err != nil {
		return err
	} else if isNil {
		c.Publisher = nil
	}

	if isNil, err = d.DecodeStruct(c.Schema); err != nil {
		return err
	} else if isNil {
		c.Schema = nil
	}

	if isNil, err = d.DecodeStruct(c.Encryption); err != nil {
		return err
	} else if isNil {
		c.Encryption = nil
	}

	if isNil, err = d.DecodeStruct(c.Compression); err != nil {
		return err
	} else if isNil {
		c.Compression = nil
	}

	if c.Flags, err = d.DecodeUint8(); err != nil {
		return err
	}

	if c.Created, err = d.DecodeTime(); err != nil {
		return err
	}

	if c.Modified, err = d.DecodeTime(); err != nil {
		return err
	}

	return nil
}
