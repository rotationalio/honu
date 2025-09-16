package region

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.rtnl.ai/honu/pkg/store/lani"
)

// Region enumerates the clouds and regions that are available to Ensign in order to
// ensure region identification and serialiation is as small a data type as possible.
// Region codes are generally broken into parts: the first digit represents the cloud,
// e.g. a region code that starts with 1 is Linode. The second series of three digits
// represents the country, e.g. USA is 840 in the ISO 3166 standard. The three digits
// represents the zone of the datacenter, and is usually cloud specific.
//
// NOTE: this guide to the enumeration representation is generally about making the
// definition easier to see and parse; but the exact information of the region should
// be looked up using the RegionInfo struct.
type Region uint32

type Regions []Region

// Geographic metadata for compliance and region-awareness.
type RegionInfo struct {
	ID       Region    `json:"id" msg:"id"`
	Name     string    `json:"name" msg:"name"`
	Country  string    `json:"country" msg:"country"`
	Zone     string    `json:"zone,omitempty" msg:"zone,omitempty"`
	Cloud    string    `json:"cloud,omitempty" msg:"cloud,omitempty"`
	Cluster  string    `json:"cluster,omitempty" msg:"cluster,omitempty"`
	Created  time.Time `json:"created,omitempty" msg:"created,omitempty"`
	Modified time.Time `json:"modified,omitempty" msg:"modified,omitempty"`
}

// Returns a list of all available regions that HonuDB is aware of.
func List() Regions {
	regions := make(Regions, 0, len(regionNames))
	for r := range regionNames {
		if r == 0 {
			continue
		}
		regions = append(regions, Region(r))
	}
	return regions
}

// Parse a region from a string or integer representation.
func Parse(r any) (Region, error) {
	switch v := r.(type) {
	case string:
		v = strings.Replace(strings.ToUpper(strings.TrimSpace(v)), "-", "_", -1)
		if id, ok := regionValues[v]; ok {
			return Region(id), nil
		}
		return UNKNOWN, fmt.Errorf("unknown region name: %q", v)
	case uint32:
		return Region(v), nil
	case Region:
		return v, nil
	default:
		return UNKNOWN, fmt.Errorf("cannot parse %T to Region", r)
	}
}

func (r Region) String() string {
	if name, ok := regionNames[uint32(r)]; ok {
		return name
	}
	return regionNames[0]
}

//===========================================================================
// Serialization
//===========================================================================

// MarshalJSON defaults to marshaling the region as its string representation.
func (r Region) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON supports unmarshaling a region from either its string or numeric representation.
func (r *Region) UnmarshalJSON(data []byte) (err error) {
	var v any
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	if *r, err = Parse(v); err != nil {
		return err
	}
	return nil
}

//===========================================================================
// Encoding and Decoding
//===========================================================================

var _ lani.Encodable = Region(0)
var _ lani.Decodable = (*Region)(nil)

func (r Region) Size() int {
	return binary.MaxVarintLen32
}

func (r Region) Encode(e *lani.Encoder) (n int, err error) {
	return e.EncodeUint32(uint32(r))
}

func (r *Region) Decode(d *lani.Decoder) (err error) {
	var v uint32
	if v, err = d.DecodeUint32(); err != nil {
		return err
	}

	*r = Region(v)
	return nil
}

var _ lani.Encodable = Regions(nil)
var _ lani.Decodable = (*Regions)(nil)

func (r Regions) Size() (s int) {
	s += binary.MaxVarintLen64
	s += len(r) * binary.MaxVarintLen32
	return s
}

func (r Regions) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeUint64(uint64(len(r))); err != nil {
		return n + m, err
	}
	n += m

	for _, region := range r {
		if m, err = e.EncodeUint32(uint32(region)); err != nil {
			return n + m, err
		}
		n += m
	}

	return n, nil
}

func (r *Regions) Decode(d *lani.Decoder) (err error) {
	var length uint64
	if length, err = d.DecodeUint64(); err != nil {
		return err
	}

	*r = make(Regions, length)
	for i := range *r {
		var v uint32
		if v, err = d.DecodeUint32(); err != nil {
			return err
		}
		(*r)[i] = Region(v)
	}
	return nil
}
