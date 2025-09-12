package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestIndex(t *testing.T) {
	// Compute the static size of the Index struct
	var staticSize int
	staticSize += 16                    // ID (ULID) is fixed length.
	staticSize += binary.MaxVarintLen64 // Length of Name
	staticSize += 1                     // Type (uint8) is fixed length.
	staticSize += 2                     // Field and Ref not nil

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Index",
		Fixture:     "index.json",
		StaticSize:  staticSize,
		FixtureSize: 104,
		New:         func() TestObject { return &metadata.Index{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}

func TestIndexType(t *testing.T) {
	testCase := &TestEnumCase{
		Name: "IndexType",
		Values: []TestEnum{
			metadata.IndexTypeUnknown,
			metadata.UNIQUE,
			metadata.INDEX,
			metadata.FOREIGN_KEY,
			metadata.VECTOR,
			metadata.SEARCH,
			metadata.COLUMN,
			metadata.BLOOM,
		},
		Strings: []string{
			"UNKNOWN",
			"UNIQUE",
			"INDEX",
			"FOREIGN_KEY",
			"VECTOR",
			"SEARCH",
			"COLUMN",
			"BLOOM",
		},
		Unknowns: "UNKNOWN",
		ICase:    true,
		ISpace:   true,
		Parse:    func(s string) (TestEnum, error) { return metadata.ParseIndexType(s) },
		New:      func(i uint8) Serializable { val := metadata.IndexType(i); return &val },
	}

	t.Run("String", testCase.TestString)
	t.Run("StringBounds", testCase.TestStringBounds)
	t.Run("Parse", testCase.TestParse)
	t.Run("JSON", testCase.TestJSON)
}
