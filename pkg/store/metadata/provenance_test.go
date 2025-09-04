package metadata_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestPublisherSize(t *testing.T) {
	// Compute the static size of the Encryption struct
	var staticSize int
	staticSize += 16 + 16               // PublisherID and ClientID (both ULID)
	staticSize += binary.MaxVarintLen64 // Length of IPAddress
	staticSize += binary.MaxVarintLen64 // Length of UserAgent

	t.Run("StaticSize", func(t *testing.T) {
		publisher := &metadata.Publisher{}
		require.Equal(t, staticSize, publisher.Size(), "expected zero valued publisher to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var publisher metadata.Publisher
		loadFixture(t, "publisher.json", &publisher)
		require.Equal(t, 77, publisher.Size(), "expected publisher to have a size of 197 bytes as computed from fixture")
	})
}

func TestPublisherSerialization(t *testing.T) {
	var obj *metadata.Publisher
	loadFixture(t, "publisher.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not encode publisher")

	cmp := &metadata.Publisher{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not decode publisher")
	require.Equal(t, obj, cmp, "deserialized publisher does not match original")
}
