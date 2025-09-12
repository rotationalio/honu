package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

type Field struct {
	Name       string    `json:"name" msg:"name"`
	Type       FieldType `json:"type" msg:"type"`
	Collection ulid.ULID `json:"collection" msg:"collection"`
}

type FieldType uint8

const (
	StringField FieldType = iota
	BlobField
	ULIDField
	UUIDField
	IntField
	UIntField
	FloatField
	TimeField
	VectorField
)

var _ lani.Encodable = (*Field)(nil)
var _ lani.Decodable = (*Field)(nil)

// The static size of a zero valued Field object; see TestFieldSize for details.
const fieldStaticSize = 27

func (o *Field) Size() int {
	return fieldStaticSize + len(o.Name)
}

func (o *Field) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeString(o.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(uint8(o.Type)); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeULID(o.Collection); err != nil {
		return n + m, err
	}
	n += m

	return
}

func (o *Field) Decode(d *lani.Decoder) (err error) {
	if o.Name, err = d.DecodeString(); err != nil {
		return err
	}

	var t uint8
	if t, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.Type = FieldType(t)

	if o.Collection, err = d.DecodeULID(); err != nil {
		return err
	}

	return nil
}

var fieldTypeNames = [9]string{
	"STRING", "BLOB", "ULID", "UUID", "INT",
	"UINT", "FLOAT", "TIME", "VECTOR",
}

func ParseFieldType(s string) (FieldType, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	for i, name := range fieldTypeNames {
		if s == name {
			return FieldType(i), nil
		}
	}
	return FieldType(0), fmt.Errorf("unknown field type: %q", s)
}

func (t FieldType) String() string {
	if int(t) < len(fieldTypeNames) {
		return fieldTypeNames[t]
	}
	return ""
}

func (t FieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *FieldType) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *t, err = ParseFieldType(alg); err != nil {
		return err
	}
	return nil
}

func (t FieldType) Value() uint8 {
	return uint8(t)
}
