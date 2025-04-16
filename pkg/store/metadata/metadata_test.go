package metadata_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lani"
	. "go.rtnl.ai/honu/pkg/store/metadata"
)

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
