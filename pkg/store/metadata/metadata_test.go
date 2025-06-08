package metadata_test

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	. "go.rtnl.ai/honu/pkg/store/metadata"
)

func TestMetadataSize(t *testing.T) {
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

	t.Run("StaticSize", func(t *testing.T) {
		metadata := &Metadata{}
		require.Equal(t, staticSize, metadata.Size(), "expected zero valued metadata to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var metadata Metadata
		loadFixture(t, "metadata.json", &metadata)
		require.Equal(t, 573, metadata.Size(), "expected metadata to have a size of 573 bytes as computed from fixture")
	})
}

func TestMetadataSerialization(t *testing.T) {
	var obj *Metadata
	loadFixture(t, "metadata.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal metdata")

	cmp := &Metadata{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal metdata")
	require.Equal(t, obj, cmp, "deserialized metdata does not match original")
}

func TestMetadataKey(t *testing.T) {
	var obj *Metadata
	loadFixture(t, "metadata.json", &obj)

	key := obj.Key()
	require.Equal(t, obj.CollectionID, key.CollectionID())
	require.Equal(t, obj.ObjectID, key.ObjectID())
	require.Equal(t, obj.Version.Scalar, key.Version())
}

func loadFixture(t *testing.T, name string, v interface{}) {
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	require.NoError(t, err, "could not open %s", path)
	defer f.Close()

	err = json.NewDecoder(f).Decode(v)
	require.NoError(t, err, "could not decode %s", path)
}
