package metadata_test

import (
	"encoding/binary"
	"testing"

	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestCollection(t *testing.T) {
	// Compute the static size of the Collection struct
	var staticSize int
	staticSize += 16                        // ID (ULID) is fixed length.
	staticSize += binary.MaxVarintLen64     // Length of Name
	staticSize += 1                         // Version not nil bool
	staticSize += 16 + 16 + 1               // Owner, Group, and Permissions (ULIDs and uint8)
	staticSize += 2 * binary.MaxVarintLen64 // Length of ACL and WriteRegions lists
	staticSize += 4                         // Publisher, Schema, Encryption, and Compression not nil bool
	staticSize += 1                         // Flags
	staticSize += binary.MaxVarintLen64     // Length of Indexes list
	staticSize += 2 * binary.MaxVarintLen64 // Created, and Modified (time.Time)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Collection",
		Fixture:     "collection.json",
		StaticSize:  staticSize,
		FixtureSize: 643,
		New:         func() TestObject { return &metadata.Collection{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}
