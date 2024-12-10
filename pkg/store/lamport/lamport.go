package lamport

import "sync"

// New returns a new clock with the specified PID. The returned clock is thread-safe
// and uses a mutex to guard against updates from multiple threads.
func New(pid uint32) Clock {
	return &clock{pid: pid}
}

// A Lamport Clock keeps track of the current process instance and monotonically
// increasing sequence timestamp. It is used to track the conflict-free distributed
// version scalar inside of a replicated system.
type Clock interface {
	// Return the next timestamp using the internal process ID of the clock.
	Next() Scalar

	// Update the clock with a timestamp scalar; if the scalar happens before the
	// current timestamp in the clock it is ignored.
	Update(Scalar)
}

var _ Clock = &clock{}

type clock struct {
	sync.Mutex
	pid     uint32
	current Scalar
}

func (c *clock) Next() Scalar {
	c.Lock()
	defer c.Unlock()

	c.current = Scalar{
		PID: c.pid,
		VID: c.current.VID + 1,
	}

	// Ensure a copy is returned so that the clock cannot be modified externally
	return c.current
}

func (c *clock) Update(now Scalar) {
	c.Lock()
	defer c.Unlock()
	if now.After(&c.current) {
		c.current = now
	}
}
