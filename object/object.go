package object

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
// returns false if the other version is later or equal to the specified version.
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

	// Either v.Version < other.Version or other.Pid > v.Pid
	return false
}

//Copies the child's attributes before updating to the parent.
func (v *Version) Clone() *Version {
	parent := &Version{
		Pid:       v.Pid,
		Version:   v.Version,
		Region:    v.Region,
		Tombstone: v.Tombstone,
	}
	return parent
}
