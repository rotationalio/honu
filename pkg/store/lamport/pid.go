package lamport

// The PID type makes it easy to create scalar versions from a previous version.
type PID uint32

func (p PID) Next(v *Scalar) Scalar {
	if v == nil {
		return Scalar{PID: uint32(p), VID: 1}
	}
	return Scalar{PID: uint32(p), VID: v.VID + 1}
}
