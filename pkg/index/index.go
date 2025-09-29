package index

import (
	"slices"

	"go.rtnl.ai/ulid"
)

// Index arrays are uses to perform reverse lookups of object IDs; e.g. mapping an
// object to all the objects that reference it. The critical property of these arrays is
// that they are always sorted in ascending order and contain no duplicates. This allows
// for efficient merging of arrays and searching using binary search.
type Array []ulid.ULID

func (a Array) Insert(id ulid.ULID) Array {
	if len(a) == 0 {
		return append(a, id)
	} else {
		if cmp := a[len(a)-1].Compare(id); cmp == 0 {
			// Already present, do not add a duplicate.
			return a
		} else if cmp < 0 {
			// Fast path: we can append to the end.
			return append(a, id)
		}
	}

	// Use binary search to find the insertion point.
	i, found := slices.BinarySearchFunc(a, id, func(a, b ulid.ULID) int {
		return a.Compare(b)
	})
	if found {
		// Already present, do not add a duplicate.
		return a
	}

	// Insert at the found index.
	a = append(a, ulid.Zero)
	copy(a[i+1:], a[i:])
	a[i] = id
	return a
}

func (a Array) Contains(id ulid.ULID) bool {
	_, found := slices.BinarySearchFunc(a, id, func(a, b ulid.ULID) int {
		return a.Compare(b)
	})
	return found
}

func (a Array) Remove(id ulid.ULID) Array {
	i, found := slices.BinarySearchFunc(a, id, func(a, b ulid.ULID) int {
		return a.Compare(b)
	})
	if !found {
		return a
	}

	// Remove the element at index i.
	copy(a[i:], a[i+1:])
	a = a[:len(a)-1]
	return a
}
