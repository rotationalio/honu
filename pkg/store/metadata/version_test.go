package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestVersion(t *testing.T) {
	var staticSize int
	staticSize += 1                     // Parent not nil bool
	staticSize += 1                     // Tombstone bool
	staticSize += binary.MaxVarintLen64 // Timestamp int64

	// Must also add the scalar size here because it is not nilable
	t.Logf("Version static size without scalar is %d", staticSize)
	expectedSize := staticSize + binary.MaxVarintLen64 + binary.MaxVarintLen32

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Version",
		Fixture:     "version.json",
		StaticSize:  expectedSize,
		FixtureSize: 42,
		New:         func() TestObject { return &metadata.Version{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)

}
