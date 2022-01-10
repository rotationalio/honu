package honu_test

import (
	"testing"

	. "github.com/rotationalio/honu"
	"github.com/rotationalio/honu/config"
	pb "github.com/rotationalio/honu/object"
	"github.com/stretchr/testify/require"
)

func TestVersionManager(t *testing.T) {
	conf := config.ReplicaConfig{}

	// Check required settings
	_, err := NewVersionManager(conf)
	require.Error(t, err)

	conf.PID = 8
	_, err = NewVersionManager(conf)
	require.Error(t, err)

	conf.Region = "us-east-2c"
	vers1, err := NewVersionManager(conf)
	require.NoError(t, err)
	require.Equal(t, "8:us-east-2c", vers1.Owner)

	conf.Name = "mitchell"
	vers1, err = NewVersionManager(conf)
	require.NoError(t, err)
	require.Equal(t, "8:mitchell", vers1.Owner)

	// Check update system
	require.Error(t, vers1.Update(nil))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		PID:     8,
	// 		Version: 1,
	// 		Region:  "us-east-2c",
	// 		Parent:  nil,
	// 	},
	// }

	obj := &pb.Object{Key: []byte("foo"), Namespace: "awesome"}
	require.NoError(t, vers1.Update(obj))
	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers1.PID, obj.Version.Pid)
	require.Equal(t, uint64(1), obj.Version.Version)
	require.Equal(t, vers1.Region, obj.Version.Region)
	require.Empty(t, obj.Version.Parent)
	require.False(t, obj.Tombstone())

	// Create a new remote versioner
	conf.PID = 13
	conf.Region = "europe-west-3"
	conf.Name = "jacques"
	vers2, err := NewVersionManager(conf)
	require.NoError(t, err)

	// Update the previous version
	require.NoError(t, vers2.Update(obj))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		Pid:     13,
	// 		Version: 2,
	// 		Region:  "europe-west-3",
	// 		Parent: &Version{
	// 			Pid:     8,
	// 			Version: 1,
	// 			Region:  "us-east-2c",
	// 		},
	// 	},
	// }

	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers2.PID, obj.Version.Pid)
	require.Equal(t, uint64(2), obj.Version.Version)
	require.Equal(t, vers2.Region, obj.Version.Region)
	require.NotEmpty(t, obj.Version.Parent)
	require.False(t, obj.Version.Parent.IsZero())
	require.Equal(t, vers1.PID, obj.Version.Parent.Pid)
	require.Equal(t, uint64(1), obj.Version.Parent.Version)
	require.Equal(t, vers1.Region, obj.Version.Parent.Region)
	require.False(t, obj.Tombstone())

	// Test Delete - creating a tombstone
	require.NoError(t, vers1.Delete(obj))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		Pid:     8,
	// 		Version: 3,
	// 		Region:  "us-east-2c",
	// 		Parent: &Version{
	// 			Pid:     13,
	// 			Version: 2,
	// 			Region:  "europe-west-3",
	// 		},
	//	    Tombstone: true,
	// 	},
	// }

	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers1.PID, obj.Version.Pid)
	require.Equal(t, uint64(3), obj.Version.Version)
	require.Equal(t, vers1.Region, obj.Version.Region)
	require.NotEmpty(t, obj.Version.Parent)
	require.False(t, obj.Version.Parent.IsZero())
	require.Equal(t, vers2.PID, obj.Version.Parent.Pid)
	require.Equal(t, uint64(2), obj.Version.Parent.Version)
	require.Equal(t, vers2.Region, obj.Version.Parent.Region)
	require.True(t, obj.Tombstone())

	// Cannot delete a deleted object
	require.Error(t, vers1.Delete(obj))

	// Cannot delete a nil object
	require.Error(t, vers1.Delete(nil))

	// Cannot delete an empty object
	require.Error(t, vers1.Delete(&pb.Object{}))
	require.Error(t, vers1.Delete(&pb.Object{Version: &pb.Version{}}))

	// Test Undelete the object
	require.NoError(t, vers1.Update(obj))

	// Expected object definition:
	// &Object{
	// 	Key:       "foo",
	// 	Namespace: "awesome",
	// 	Region:    "us-east-2c",
	// 	Owner:     "8:mitchell",
	// 	Version: &Version{
	// 		Pid:     8,
	// 		Version: 4,
	// 		Region:  "us-east-2c",
	// 		Parent: &Version{
	// 			Pid:     9,
	// 			Version: 3,
	// 			Region:  "us-east-2c",
	// 		},
	//	    Tombstone: false,
	// 	},
	// }

	require.Equal(t, vers1.Region, obj.Region)
	require.Equal(t, vers1.Owner, obj.Owner)
	require.Equal(t, vers1.PID, obj.Version.Pid)
	require.Equal(t, uint64(4), obj.Version.Version)
	require.Equal(t, vers1.Region, obj.Version.Region)
	require.NotEmpty(t, obj.Version.Parent)
	require.False(t, obj.Version.Parent.IsZero())
	require.Equal(t, vers1.PID, obj.Version.Parent.Pid)
	require.Equal(t, uint64(3), obj.Version.Parent.Version)
	require.Equal(t, vers1.Region, obj.Version.Parent.Region)
	require.False(t, obj.Tombstone())
}
