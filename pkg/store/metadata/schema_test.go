package metadata_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestSchemaVersionSize(t *testing.T) {
	var staticSize int
	staticSize += binary.MaxVarintLen64     // Length of name string
	staticSize += 3 * binary.MaxVarintLen32 // Major, Minor, Patch (all uint32)

	t.Run("StaticSize", func(t *testing.T) {
		schemaVersion := &metadata.SchemaVersion{}
		require.Equal(t, staticSize, schemaVersion.Size(), "expected zero valued SchemaVersion to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var schemaVersion metadata.SchemaVersion
		loadFixture(t, "schema.json", &schemaVersion)
		require.Equal(t, 32, schemaVersion.Size(), "expected SchemaVersion to have a size of 32 bytes as computed from fixture")
	})
}

func TestSchemaVersionSerialization(t *testing.T) {
	var obj *metadata.SchemaVersion
	loadFixture(t, "schema.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal schema version")

	cmp := &metadata.SchemaVersion{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal schema version")
	require.Equal(t, obj, cmp, "deserialized schema version does not match original")
}
