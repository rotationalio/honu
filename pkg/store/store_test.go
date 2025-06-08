package store_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store"
	"go.rtnl.ai/ulid"
)

func TestSystemCollections(t *testing.T) {
	// System collections must be less than any ULID that would be generated.
	// They also must be sorted in a specific order for efficient access.
	collections := []ulid.ULID{
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	require.True(t, sort.SliceIsSorted(collections, func(i, j int) bool {
		return bytes.Compare(collections[i][:], collections[j][:]) < 0
	}))

	// All collection ULID timestamps must be less than 1984-04-07T09:21:07-06:00
	epoch, _ := time.Parse(time.RFC3339, "1984-04-07T09:21:07-06:00")
	for _, collection := range collections {
		ts := collection.Timestamp()
		require.True(t, ts.Before(epoch), "collection %s timestamp %s is not before epoch %s", collection, ts, epoch)
	}
}

//===========================================================================
// Helper Functions
//===========================================================================

func loadFixture(t *testing.T, name string, v interface{}) {
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	require.NoError(t, err, "could not open %s", path)
	defer f.Close()

	err = json.NewDecoder(f).Decode(v)
	require.NoError(t, err, "could not decode %s", path)
}
