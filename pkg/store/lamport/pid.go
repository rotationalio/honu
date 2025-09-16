package lamport

import (
	"sync"
)

// The PID type makes it easy to create scalar versions from a previous version.
type PID uint32

func (p PID) Next(v *Scalar) Scalar {
	if v == nil {
		return Scalar{PID: uint32(p), VID: 1}
	}
	return Scalar{PID: uint32(p), VID: v.VID + 1}
}

// A Global PID can be set by the application to ensure that all modules have access
// to the same PID for generating versions. This is thread-safe to read from and should
// only be written once at the start of a process. If not set, will panic when used.
var (
	pidMu     sync.RWMutex
	processID *uint32
)

func SetProcessID(pid uint32) {
	pidMu.Lock()
	defer pidMu.Unlock()

	processID = &pid
}

func ProcessID() PID {
	pidMu.RLock()
	defer pidMu.RUnlock()

	if processID == nil {
		panic("process ID not set")
	}
	return PID(*processID)
}

func Next(v *Scalar) Scalar {
	pidMu.RLock()
	defer pidMu.RUnlock()

	if processID == nil {
		panic("process ID not set")
	}

	if v == nil {
		return Scalar{PID: uint32(*processID), VID: 1}
	}
	return Scalar{PID: uint32(*processID), VID: v.VID + 1}
}
