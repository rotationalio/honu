package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestSchemaVersion(t *testing.T) {
	var staticSize int
	staticSize += binary.MaxVarintLen64     // Length of name string
	staticSize += 3 * binary.MaxVarintLen32 // Major, Minor, Patch (all uint32)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "SchemaVersion",
		Fixture:     "schema.json",
		StaticSize:  staticSize,
		FixtureSize: 32,
		New:         func() TestObject { return &metadata.SchemaVersion{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}
