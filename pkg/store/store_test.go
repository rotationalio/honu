package store_test

import (
	"bytes"
	"fmt"
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

func Example_system_names() {
	convert := func(b []byte) string {
		s := make([]byte, 0, len(b))
		for i, c := range b {
			if c < 0x20 || c > 0x7E {
				if i == 0 {
					continue
				}
				s = append(s, ' ')
			} else {
				s = append(s, c)
			}
		}

		return string(s)
	}

	systemIDs := []ulid.ULID{
		store.SystemHonuAgent,
		store.SystemCollections,
		store.SystemReplicas,
		store.SystemAccessControl,
	}

	for _, sysid := range systemIDs {
		ts := sysid.Timestamp()
		name := convert(sysid[:])

		fmt.Printf("%s (%s)\n", name, ts.UTC().Format(time.RFC3339))
	}

	// Output:
	// honu honuagent  (1984-03-19T12:08:28Z)
	// honu collection (1984-03-19T12:08:28Z)
	// honu accesslist (1984-03-19T12:08:28Z)
	// honu networking (1984-03-19T12:08:28Z)
}
