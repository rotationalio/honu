package object_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
)

func TestSystemSerialize(t *testing.T) {
	tests := []struct {
		path string
		src  lani.Encodable
		dst  lani.Decodable
	}{
		{
			"collection.json",
			&metadata.Collection{},
			&metadata.Collection{},
		},
	}

	for _, tc := range tests {
		loadFixture(t, tc.path, tc.src)
		data, err := object.MarshalSystem(tc.src)
		require.NoError(t, err, "could not marshal %T", tc.src)
		require.NotZero(t, data, "marshaled %T should not be empty", tc.src)

		err = object.UnmarshalSystem(data, tc.dst)
		require.NoError(t, err, "could not unmarshal %T", tc.dst)
		require.Equal(t, tc.src, tc.dst, "deserialized %T does not match original", tc.src)
	}
}
