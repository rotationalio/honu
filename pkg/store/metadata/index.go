package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/ulid"
)

// Indexes are nested btree buckets inside of collections that allow for efficient
// lookups of objects based on specific attributes and aid in querying and retrieval.
// Index metadata defines how the index is structured and what it contains.
type Index struct {
	ID    ulid.ULID `json:"id" msg:"id"`
	Name  string    `json:"name" msg:"name"`
	Type  IndexType `json:"type" msg:"type"`
	Field Field     `json:"field" msg:"field"`
	Ref   Field     `json:"ref" msg:"ref"`
}

type IndexType uint8

const (
	IndexTypeUnknown IndexType = iota
	UNIQUE                     // Enforces uniqueness of indexed values
	INDEX                      // Non-unique index for faster lookups
	FOREIGN_KEY                // Ensures referential integrity between collections
	VECTOR                     // Vector index for semantic search
	SEARCH                     // Full-text search index
	COLUMN                     // Stores data in a columnar format for aggregations
	BLOOM                      // Probabilistic data structure for membership queries
)

var indexTypeNames = [8]string{
	"UNKNOWN", "UNIQUE", "INDEX", "FOREIGN_KEY",
	"VECTOR", "SEARCH", "COLUMN", "BLOOM",
}

var _ lani.Encodable = (*Index)(nil)
var _ lani.Decodable = (*Index)(nil)

// The static size of a zero valued Index object; see TestIndexSize for details.
const indexStaticSize = 11

func (o *Index) Size() int {
	return indexStaticSize + len(o.Name)
}

func (o *Index) Encode(e *lani.Encoder) (n int, err error) {
	var m int
	if m, err = e.EncodeULID(o.ID); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeString(o.Name); err != nil {
		return n + m, err
	}
	n += m

	if m, err = e.EncodeUint8(uint8(o.Type)); err != nil {
		return n + m, err
	}
	n += m

	return n, nil
}

func (o *Index) Decode(d *lani.Decoder) (err error) {
	if o.ID, err = d.DecodeULID(); err != nil {
		return err
	}

	if o.Name, err = d.DecodeString(); err != nil {
		return err
	}

	var t uint8
	if t, err = d.DecodeUint8(); err != nil {
		return err
	}
	o.Type = IndexType(t)

	return nil
}

func ParseIndexType(s string) (IndexType, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	for i, name := range indexTypeNames {
		if s == name {
			return IndexType(i), nil
		}
	}
	return IndexType(0), fmt.Errorf("unknown index type: %q", s)
}

func (t IndexType) String() string {
	if int(t) < len(indexTypeNames) {
		return indexTypeNames[t]
	}
	return indexTypeNames[0]
}

func (t *IndexType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *IndexType) UnmarshalJSON(data []byte) (err error) {
	var alg string
	if err := json.Unmarshal(data, &alg); err != nil {
		return err
	}
	if *t, err = ParseIndexType(alg); err != nil {
		return err
	}
	return nil
}
