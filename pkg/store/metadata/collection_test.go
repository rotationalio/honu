package metadata_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestCollectionSize(t *testing.T) {
	// Compute the static size of the Collection struct
	var staticSize int
	staticSize += 16                        // ID (ULID) is fixed length.
	staticSize += binary.MaxVarintLen64     // Length of Name
	staticSize += 1                         // Version not nil bool
	staticSize += 16 + 16 + 1               // Owner, Group, and Permissions (ULIDs and uint8)
	staticSize += 2 * binary.MaxVarintLen64 // Length of ACL and WriteRegions lists
	staticSize += 4                         // Publisher, Schema, Encryption, and Compression not nil bool
	staticSize += 1                         // Flags
	staticSize += 2 * binary.MaxVarintLen64 // Created, and Modified (time.Time)

	t.Run("StaticSize", func(t *testing.T) {
		collection := &metadata.Collection{}
		require.Equal(t, staticSize, collection.Size(), "expected zero valued Collection to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var collection metadata.Collection
		loadFixture(t, "collection.json", &collection)
		require.Equal(t, 414, collection.Size(), "expected Collection to have a size of 414 bytes as computed from fixture")
	})
}

func TestCollectionSerialization(t *testing.T) {
	var obj *metadata.Collection
	loadFixture(t, "collection.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal Collection")

	cmp := &metadata.Collection{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal Collection")
	require.Equal(t, obj, cmp, "deserialized Collection does not match original")
}
