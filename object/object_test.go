package object_test

import (
	"testing"

	. "github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/require"
)

func TestTombstone(t *testing.T) {
	obj := &Object{
		Key:       []byte("foo"),
		Namespace: "general",
		Version: &Version{
			Pid:       8,
			Version:   1,
			Region:    "us-southwest-14",
			Tombstone: false,
		},
	}

	require.False(t, obj.Tombstone())
	obj.Version.Tombstone = true
	require.True(t, obj.Tombstone())
}

func TestVersionIsLater(t *testing.T) {
	// Create a version with non-zero values to test against
	v1 := &Version{Pid: 8, Version: 42}

	// Tests against version zero
	require.True(t, v1.IsLater(&VersionZero), "version should be later than zero value")
	require.True(t, v1.IsLater(nil), "version should treat nil as zero value")
	require.False(t, VersionZero.IsLater(v1), "zero value should not be later than version")

	require.False(t, VersionZero.IsLater(&VersionZero), "zero value should not be later than itself")
	require.False(t, VersionZero.IsLater(nil), "zero value should not be later than nil")

	// Test against other versions
	require.True(t, v1.IsLater(&Version{Pid: 8, Version: 40}), "version should be later than version with lesser scalar but same pid")
	require.True(t, v1.IsLater(&Version{Pid: 9, Version: 42}), "concurrent version should be later than lower precedence pid")
	require.False(t, v1.IsLater(&Version{Pid: 7, Version: 42}), "concurrent version should not be later than higher precedence pid")
	require.False(t, v1.IsLater(&Version{Pid: 8, Version: 44}), "version should not be later than version with greater scaler but same pid")
	require.False(t, v1.IsLater(&Version{Pid: 8, Version: 42}), "version should not be later than an equal version")
}

func TestVersionEqual(t *testing.T) {
	// Test if a version is equal to an identical version
	v1 := &Version{Pid: 8, Version: 42}
	v2 := &Version{Pid: 8, Version: 42}
	require.True(t, v1.Equal(v2), "two versions with same PID and scalar should be equal")
	require.True(t, v2.Equal(v1), "version equality should be reciprocal")
	require.True(t, v1.Equal(v1), "a version object should be equal to itself")

	// Test no equality
	require.False(t, v1.Equal(&VersionZero), "a version should not be equal to zero value")
	require.False(t, v1.Equal(nil), "a version should not be equal to nil")
	require.False(t, VersionZero.Equal(v2), "zero value should not be equal to a version")
	require.False(t, v1.Equal(&Version{Pid: 9, Version: 42}), "versions with different PIDs should not be equal")
	require.False(t, v1.Equal(&Version{Pid: 9, Version: 41}), "versions with different PIDs and scalars should not be equal")
	require.False(t, v1.Equal(&Version{Pid: 8, Version: 43}), "versions with a greater scalar should not be equal")
	require.False(t, v1.Equal(&Version{Pid: 8, Version: 41}), "versions with a lesser scalar should not be equal")

	// Test zero-versioned equality
	require.True(t, VersionZero.Equal(nil), "zero value should equal nil")
	require.True(t, VersionZero.Equal(&VersionZero), "zero value should equal zero value")
}

func TestVersionConcurrent(t *testing.T) {
	// Test two concurrent versions
	v1 := &Version{Pid: 8, Version: 42}
	v2 := &Version{Pid: 12, Version: 42}

	require.True(t, v1.Concurrent(v2), "versions should be concurrent with same scalar and different PID")
	require.True(t, v2.Concurrent(v1), "concurrent should be reciprocal")

	// Equal, later, and earlier versions should not be concurrent even if they have the same PID
	require.False(t, v1.Concurrent(&Version{Pid: 8, Version: 42}), "equal versions should not be concurrent")
	require.False(t, v1.Concurrent(&Version{Pid: 12, Version: 48}), "later versions should not be concurrent")
	require.False(t, v1.Concurrent(&Version{Pid: 12, Version: 21}), "earlier versions should not be concurrent")
	require.False(t, v1.Concurrent(&VersionZero), "zero version should not be concurrent")
	require.False(t, v1.Concurrent(nil), "nil should be treated as zero version")

	// Test zero-versioned concurrency
	require.False(t, VersionZero.Concurrent(&VersionZero), "zero version should not be concurrent with itself")
	require.False(t, VersionZero.Concurrent(nil), "zero version should not be concurrent with nil")
	require.False(t, VersionZero.Concurrent(v1), "zero version should not be concurrent with another version")
}

func TestLinearFrom(t *testing.T) {
	// Root version
	v1 := &Version{Pid: 8, Version: 1, Parent: nil}

	// Linear from v1 but with different PIDs
	v2 := &Version{Pid: 8, Version: 2, Parent: &Version{Pid: 8, Version: 1}}
	v3 := &Version{Pid: 9, Version: 2, Parent: &Version{Pid: 8, Version: 1}}

	// Linear from v3
	v4 := &Version{Pid: 9, Version: 3, Parent: &Version{Pid: 9, Version: 2}}

	require.True(t, v1.LinearFrom(&VersionZero), "root version should be linear from version zero")
	require.True(t, v1.LinearFrom(nil), "should treat nil as the zero version")
	require.False(t, v1.LinearFrom(v2), "parent should not be linear from child")

	require.True(t, v2.LinearFrom(v1), "version should be linear from its parent, same PID")
	require.True(t, v3.LinearFrom(v1), "version should be linear from its parent, different PID")
	require.True(t, v4.LinearFrom(v3), "a child should be linear from its parent")

	require.False(t, v2.LinearFrom(v3), "concurrent should not be linear from, other version")
	require.False(t, v3.LinearFrom(v2), "stomp should not be linear from other version")

	require.False(t, v4.LinearFrom(v2), "a child should not be linear from stomped version")
	require.False(t, v4.LinearFrom(v1), "a skip should not be linear from earlier version")
	require.False(t, v1.LinearFrom(v4), "version should not be linear from later version")
}

func TestStomps(t *testing.T) {
	// Root version
	v1 := &Version{Pid: 8, Version: 1, Parent: nil}

	require.False(t, v1.Stomps(&VersionZero), "root version should not stomp version zero")
	require.False(t, v1.Stomps(nil), "should treat nil as the zero version")
	require.False(t, v1.Stomps(&Version{Pid: 8, Version: 1}), "version should not stomp equal version")

	// Linear from v1 but with different PIDs (v2 stomps v3)
	v2 := &Version{Pid: 8, Version: 2, Parent: &Version{Pid: 8, Version: 1}}
	v3 := &Version{Pid: 9, Version: 2, Parent: &Version{Pid: 8, Version: 1}}

	require.True(t, v2.Stomps(v3), "stomps should be detected for concurrent versions")
	require.False(t, v3.Stomps(v2), "only one concurrent version should stomp the other")
	require.False(t, v2.Stomps(v2), "version should not stomp itself")

	// Linear from v3 (should not stomp v2 even though it's later)
	v4 := &Version{Pid: 9, Version: 3, Parent: &Version{Pid: 9, Version: 2}}
	require.False(t, v4.Stomps(v2), "later version should not stomp earlier version if not concurent")
	require.False(t, v4.Stomps(v3), "version should not stomp its parent")
}

func TestSkips(t *testing.T) {
	// Root version
	v1 := &Version{Pid: 8, Version: 1, Parent: nil}

	require.False(t, v1.Skips(&VersionZero), "root version should not skip zero value")
	require.False(t, v1.Skips(nil), "should treat nil as the zero version")
	require.False(t, v1.Skips(&Version{Pid: 8, Version: 1}), "version should not skip equal version")

	// Linear from v1 but with different PIDs (v2 stomps v3)
	v2 := &Version{Pid: 8, Version: 2, Parent: &Version{Pid: 8, Version: 1}}
	v3 := &Version{Pid: 9, Version: 2, Parent: &Version{Pid: 8, Version: 1}}

	require.False(t, v2.Skips(v1), "version should not skip its parent")
	require.False(t, v2.Skips(v3), "stomp should not be a skip")
	require.False(t, v3.Skips(v2), "concurrent version should not be a skip")

	// Linear from v3 (skips v3 from v1)
	v4 := &Version{Pid: 9, Version: 3, Parent: &Version{Pid: 9, Version: 2}}
	require.True(t, v4.Skips(v1), "a later version should skip an earlier version if parent is not equal")
	require.True(t, v4.Skips(v2), "a branch should be considered a skip if the parent was stomped")
	require.False(t, v4.Skips(v3), "a child should not skip a non-root parent")
}

func TestVersionHistory(t *testing.T) {
	// Create a version history without reusing pointers
	v71 := &Version{Pid: 7, Version: 1, Parent: nil}
	v72 := &Version{Pid: 7, Version: 2, Parent: &Version{Pid: 7, Version: 1}}
	v73 := &Version{Pid: 7, Version: 3, Parent: &Version{Pid: 7, Version: 2}}
	v83 := &Version{Pid: 8, Version: 3, Parent: &Version{Pid: 7, Version: 2}}
	v74 := &Version{Pid: 7, Version: 4, Parent: &Version{Pid: 7, Version: 3}}
	v84 := &Version{Pid: 8, Version: 4, Parent: &Version{Pid: 8, Version: 3}}
	v85 := &Version{Pid: 8, Version: 5, Parent: &Version{Pid: 8, Version: 4}}
	v95 := &Version{Pid: 9, Version: 5, Parent: &Version{Pid: 7, Version: 4}}
	v76 := &Version{Pid: 7, Version: 6, Parent: &Version{Pid: 9, Version: 5}}
	v97 := &Version{Pid: 9, Version: 7, Parent: &Version{Pid: 7, Version: 6}}

	// Test the 71 through 97 linear history
	history := []*Version{v71, v72, v73, v74, v95, v76, v97}
	for i := 1; i < len(history); i++ {
		parent := history[i-1]
		current := history[i]
		require.True(t, current.LinearFrom(parent), "%s should be linear from %s", current, parent)
	}

	// Test the 71 through 85 linear history
	history = []*Version{v71, v72, v83, v84, v85}
	for i := 1; i < len(history); i++ {
		parent := history[i-1]
		current := history[i]
		require.True(t, current.LinearFrom(parent), "%s should be linear from %s", current, parent)
	}
}
