package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/ulid"
)

func TestIsSystemCollection(t *testing.T) {
	collections := []ulid.ULID{
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	for _, collectionID := range collections {
		collection := &store.Collection{Collection: metadata.Collection{ID: collectionID}}
		require.True(t, collection.IsSystem(), "expected %s to be a system collection", collection)
	}

	// Test a non-system collection
	nonSystemCollection := &store.Collection{Collection: metadata.Collection{ID: ulid.MustNew(ulid.Now(), nil)}}
	require.False(t, nonSystemCollection.IsSystem(), "expected %s to not be a system collection", nonSystemCollection)
}
