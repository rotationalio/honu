package versions

import (
	"errors"
	"fmt"

	"github.com/rotationalio/honu/config"
	pb "github.com/rotationalio/honu/object/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// New creates a new manager for handling lamport scalar versions.
func New(conf config.ReplicaConfig) (v *Manager, err error) {
	// Make sure we don't create a VersionManager that is unable to do its job.
	if conf.PID == 0 || conf.Region == "" {
		return nil, errors.New("improperly configured: version manager requires PID and Region")
	}

	v = &Manager{PID: conf.PID, Region: conf.Region}

	// Compute the owner name
	if conf.Name != "" {
		// The common name
		v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Name)
	} else {
		// The owner name is just the pid:region
		v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Region)
	}

	return v, nil
}

// Manager is configured with information associated with the local Replica in
// order to correctly implement Lamport clocks for sequential, conflict-free replicated
// versions.
type Manager struct {
	PID    uint64
	Owner  string
	Region string
}

// Update the version of an object in place.
// If the object was previously a tombstone (e.g. deleted) then it is "undeleted".
func (v Manager) Update(meta *pb.Object) error {
	if meta == nil {
		return errors.New("cannot update version on empty object")
	}

	// Update the parent to the current version of the object.
	if meta.Version != nil && !meta.Version.IsZero() {
		meta.Version.Parent = meta.Version.Clone()
	} else {
		// This is the first version of the object. Also set provenance on the object.
		meta.Version = &pb.Version{}
		meta.Region = v.Region
		meta.Owner = v.Owner
	}

	// Update the version to the new version of the local version manager,
	// Undelete the version if it was a Tombstone before
	v.updateVersion(meta, false)
	return nil
}

// Delete creates a tombstone version of the object in place.
func (v Manager) Delete(meta *pb.Object) error {
	if meta == nil {
		return errors.New("cannot create tombstone version on empty object")
	}

	if meta.Version == nil || meta.Version.IsZero() {
		// This is the first version of the object, it cannot be deleted.
		return errors.New("cannot delete version that doesn't exist yet")
	}

	// Cannot delete a tombstone
	if meta.Version.Tombstone {
		return errors.New("cannot delete an already deleted object")
	}
	meta.Version.Parent = meta.Version.Clone()

	//Update Pid, Version, Region and Tombstone for the version.
	v.updateVersion(meta, true)
	return nil
}

// Assigns the attributes of the passed versionManager to the object.
func (v Manager) updateVersion(meta *pb.Object, delete_version bool) {
	meta.Version.Pid = v.PID
	meta.Version.Version++
	meta.Version.Region = v.Region
	meta.Version.Tombstone = delete_version
	meta.Version.Modified = timestamppb.Now()
}
