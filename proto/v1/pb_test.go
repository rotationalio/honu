package pb_test

import (
	"testing"

	. "github.com/rotationalio/honu/proto/v1"
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
	v1 := &Version{Pid: 8, Version: 42}
	require.True(t, v1.IsLater(&Version{Pid: 8, Version: 40}))
	require.True(t, v1.IsLater(&Version{Pid: 9, Version: 42}))
	require.False(t, v1.IsLater(&Version{Pid: 7, Version: 42}))
	require.False(t, v1.IsLater(&Version{Pid: 8, Version: 44}))
}
