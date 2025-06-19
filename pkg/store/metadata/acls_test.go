package metadata_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
)

func TestAccessControlSize(t *testing.T) {
	staticSize := 17 // 16 bytes for ClientID and 1 byte for permissions (uint8)

	t.Run("StaticSize", func(t *testing.T) {
		acl := &metadata.AccessControl{}
		require.Equal(t, staticSize, acl.Size(), "expected zero valued compression to have a static size of %d bytes", staticSize)
	})

	t.Run("VariableSize", func(t *testing.T) {
		var acl metadata.AccessControl
		loadFixture(t, "acl.json", &acl)
		require.Equal(t, 17, acl.Size(), "expected compression to have a size of 11 bytes as computed from fixture")
	})
}

func TestAccessControlSerialization(t *testing.T) {
	var obj *metadata.AccessControl
	loadFixture(t, "acl.json", &obj)

	data, err := lani.Marshal(obj)
	require.NoError(t, err, "could not marshal access control")

	cmp := &metadata.AccessControl{}
	err = lani.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal access control")
	require.Equal(t, obj, cmp, "deserialized access control does not match original")
}
