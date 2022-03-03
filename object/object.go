package object

import "time"

var VersionZero = Version{}

// Tombstone returns true if the version of the object is a Tombstone (for a deleted object)
func (o *Object) Tombstone() bool {
	if o.Version == nil {
		panic("object improperly initialized without version")
	}
	return o.Version.Tombstone
}

// IsZero determines if the version is zero valued (e.g. the PID and Version are zero).
// Note that zero-valuation does not check parent or region.
func (v *Version) IsZero() bool {
	return v.Pid == 0 && v.Version == 0
}

// IsLater returns true if the specified version is later than the other version. It
// returns false if the other version is later or equal to the specified version. If
// other is nil, it is considered equivalent to the zero-valued version.
//
// Versions are Lamport Scalars composed of a monotonically increasing scalar component
// along with a tiebreaking component called the PID (the process ID of a replica). If
// the scalar is greater than the other scaler, then the Version is later. If the
// scalars are equal, then the version with the lower PID has higher precedence.
// Versions are conflict free so long as the processes that increment the scalar have a
// unique PID. E.g. if two processes concurrently increment the scalar to the same value
// then unique PIDs guarantee that one of those processes will have precedence. Lower
// PIDs are used for the higher precedence because of the PIDs are assigned in
// increasing order, then the lower PIDs are older replicas.
func (v *Version) IsLater(other *Version) bool {
	// If other is nil, then we assume it represents the zero-valued version.
	if other == nil {
		other = &VersionZero
	}

	// Version is monotonically increasing, if it's greater than the other, then this
	// version is later than the other.
	if v.Version > other.Version {
		return true
	}

	// If the versions are equal, then the version with the lower PID has higher precedence
	if v.Version == other.Version && v.Pid < other.Pid {
		return true
	}

	// Either v.Version < other.Version, other.Pid > v.Pid, or the versions are equal.
	return false
}

// Equal returns true if and only if the Version scalars and PIDs are equal. Nil is
// considered to be the zero valued Version.
func (v *Version) Equal(other *Version) bool {
	// If other is nil, then we assume it represents the zero-valued version.
	if other == nil {
		other = &VersionZero
	}
	return v.Version == other.Version && v.Pid == other.Pid
}

// Concurrent returns true if and only if the Version scalars are equal but the PIDs are
// not equal. Nil is considered to be the zero valued Version.
func (v *Version) Concurrent(other *Version) bool {
	// If other is nil, then we assume it represents the zero-valued version.
	if other == nil {
		other = &VersionZero
	}
	return v.Version == other.Version && v.Pid != other.Pid
}

// LinearFrom returns true if and only if the the parent of this version is equal to the
// other version. E.g. is this version created from the other version (a child of).
// LinearFrom implies that this version is later than the other version so long as the
// child is always later than the parent version.
//
// NOTE: this method cannot detect a linear chain through multiple ancestors to the
// current version - it is a direct relationship only. A complete version history is
// required to compute a longer linear chain and branch points.
func (v *Version) LinearFrom(other *Version) bool {
	// Handle the case of the root version (no parent)
	if v.Parent == nil {
		return other == nil || other.IsZero()
	}
	return v.Parent.Equal(other)
}

// Stomps returns true if and only if this version is both concurrent and later than the
// other version; e.g. this version is concurrent and would have precedence.
func (v *Version) Stomps(other *Version) bool {
	return v.Concurrent(other) && v.IsLater(other)
}

// Skips returns true if and only if this version is later, is not concurrent (e.g. is
// not a stomp) and is not linear from the other version, e.g. this version is at least
// one version farther in the version history than the other version. Skips does not
// imply a stomp, though a stomp is possible in the version chain between the skip.
// Skips really mean that the  replica does not have enough information to determine if
// a stomp has occurred or if we've just moved forward in a linear chain.
func (v *Version) Skips(other *Version) bool {
	return v.IsLater(other) && !v.Concurrent(other) && !v.LinearFrom(other)
}

//Copies the child's attributes before updating to the parent.
func (v *Version) Clone() *Version {
	parent := &Version{
		Pid:       v.Pid,
		Version:   v.Version,
		Region:    v.Region,
		Tombstone: v.Tombstone,
		Created:   v.Created,
	}
	return parent
}

// Returns the Version timestamp
func (v *Version) Timestamp() time.Time {
	return time.Unix(v.Created, 0)
}
