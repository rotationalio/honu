package metadata_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestVersionSize(t *testing.T) {
	var staticSize int
	staticSize += 1                     // Parent not nil bool
	staticSize += 1                     // Tombstone bool
	staticSize += binary.MaxVarintLen64 // Timestamp int64

	t.Run("StaticSize", func(t *testing.T) {
		// Must also add the scalar size here because it is not nilable
		expectedSize := staticSize + binary.MaxVarintLen64 + binary.MaxVarintLen32

		version := &metadata.Version{}
		require.Equal(t, expectedSize, version.Size(), "expected zero valued Version to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var version metadata.Version
		loadFixture(t, "version.json", &version)
		require.Equal(t, 42, version.Size(), "expected Version to have a size of 42 bytes as computed from fixture")
	})
}

func TestVersionSerialization(t *testing.T) {
	var obj *metadata.Version
	loadFixture(t, "version.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal version")

	cmp := &metadata.Version{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal version")
	require.Equal(t, obj, cmp, "deserialized version does not match original")
}
