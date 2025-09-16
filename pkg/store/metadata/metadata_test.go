package metadata_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	. "go.rtnl.ai/honu/pkg/store/metadata"
)

func TestMetadata(t *testing.T) {
	// Compute the static size of the Metadata struct
	var staticSize int
	staticSize += 16 + 16                   // ObjectID and CollectionID (ULIDs) are fixed length.
	staticSize += 2                         // Version and SchemaVersion not nil bool
	staticSize += binary.MaxVarintLen64     // Length of MIME
	staticSize += 16 + 16 + 1               // Owner, Group, and Permissions (ULIDs and uint8)
	staticSize += 2 * binary.MaxVarintLen64 // Length of ACL and WriteRegions lists
	staticSize += 3                         // Publisher, Encryption, and Compression not nil bool
	staticSize += 1                         // Flags
	staticSize += 2 * binary.MaxVarintLen64 // Created, and Modified (time.Time)

	// Create a test generic case and execute the tests
	testCase := &TestCase{
		Name:        "Metadata",
		Fixture:     "metadata.json",
		StaticSize:  staticSize,
		FixtureSize: 557,
		New:         func() TestObject { return &Metadata{} },
	}

	t.Run("StaticSize", testCase.TestStaticSize)
	t.Run("VariableSize", testCase.TestVariableSize)
	t.Run("Serialization", testCase.TestSerialization)
}

func TestMetadataKey(t *testing.T) {
	var obj *Metadata
	loadFixture(t, "metadata.json", &obj)

	key := obj.Key()
	require.Equal(t, obj.ObjectID, key.ObjectID())
	require.Equal(t, obj.Version.Scalar, key.Version())
}
