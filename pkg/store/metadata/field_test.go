package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestField(t *testing.T) {
	// Compute the size of a Field struct
	var staticSize int
	staticSize += binary.MaxVarintLen64 // Length of Name
	staticSize += 1                     // Type (uint8) is fixed length.
	staticSize += 16                    // Collection (ULID) is fixed length.

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Field",
		Fixture:     "field.json",
		StaticSize:  staticSize,
		FixtureSize: 36,
		New:         func() TestObject { return &metadata.Field{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}

func TestFieldType(t *testing.T) {
	testCase := &TestEnumCase{
		Name: "IndexType",
		Values: []TestEnum{
			metadata.StringField,
			metadata.BlobField,
			metadata.ULIDField,
			metadata.UUIDField,
			metadata.IntField,
			metadata.UIntField,
			metadata.FloatField,
			metadata.TimeField,
			metadata.VectorField,
		},
		Strings: []string{
			"STRING",
			"BLOB",
			"ULID",
			"UUID",
			"INT",
			"UINT",
			"FLOAT",
			"TIME",
			"VECTOR",
		},
		Unknowns: "",
		ICase:    true,
		ISpace:   true,
		Parse:    func(s string) (TestEnum, error) { return metadata.ParseFieldType(s) },
		New:      func(i uint8) Serializable { val := metadata.FieldType(i); return &val },
	}

	t.Run("String", testCase.TestString)
	t.Run("StringBounds", testCase.TestStringBounds)
	t.Run("Parse", testCase.TestParse)
	t.Run("JSON", testCase.TestJSON)
}
