package honu

import (
	"errors"
	"fmt"
	"net"

	"github.com/rotationalio/honu/config"
	pb "github.com/rotationalio/honu/proto/v1"
)

// NewVersionManager creates a new manager for handling lamport scalar versions.
func NewVersionManager(conf config.ReplicaConfig) (v *VersionManager, err error) {
	// Make sure we don't create a VersionManager that is unable to do its job.
	if conf.PID == 0 || conf.Region == "" {
		return nil, errors.New("improperly configured: version manager requires PID and Region")
	}

	v = &VersionManager{PID: conf.PID, Region: conf.Region}

	// Compute the owner name
	if conf.Name != "" {
		// The common name
		v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Name)
	} else {
		// Check to see if there is a domain name in the bindaddr
		var host string
		if host, _, err = net.SplitHostPort(conf.BindAddr); err == nil {
			v.Owner = fmt.Sprintf("%d:%s", conf.PID, host)
		} else {
			// The owner name is just the pid:region in the last case
			v.Owner = fmt.Sprintf("%d:%s", conf.PID, conf.Region)
		}
	}

	return v, nil
}

type VersionManager struct {
	PID    uint64
	Owner  string
	Region string
}

// Update the version of an object in place.
// If the object was previously a tombstone (e.g. deleted) then it is "undeleted".
func (v VersionManager) Update(meta *pb.Object) error {
	if meta == nil {
		return errors.New("cannot update version on empty object")
	}

	// Update the parent to the current version of the object.
	if meta.Version != nil && !meta.Version.IsZero() {
		meta.Version.Parent = &pb.Version{
			Pid:     meta.Version.Pid,
			Version: meta.Version.Version,
			Region:  meta.Version.Region,
		}
	} else {
		// This is the first version of the object. Also set provenance on the object.
		meta.Version = &pb.Version{}
		meta.Region = v.Region
		meta.Owner = v.Owner
	}

	// Update the version to the new version of the local version manager
	meta.Version.Pid = v.PID
	meta.Version.Version++
	meta.Version.Region = v.Region

	// Undelete the version if it was a Tombstone before
	meta.Version.Tombstone = false
	return nil
}

// Delete creates a tombstone version of the object in place.
func (v VersionManager) Delete(meta *pb.Object) error {
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

	meta.Version.Parent = &pb.Version{
		Pid:     meta.Version.Pid,
		Version: meta.Version.Version,
		Region:  meta.Version.Region,
	}

	meta.Version.Pid = v.PID
	meta.Version.Version++
	meta.Version.Region = v.Region
	meta.Version.Tombstone = true
	return nil
}
