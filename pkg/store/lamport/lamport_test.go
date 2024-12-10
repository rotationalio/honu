package lamport_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rotationalio/honu/pkg/store/lamport"
	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	clock := lamport.New(1)
	require.Equal(t, lamport.Scalar{1, 1}, clock.Next(), "expected next timestamp to be 1")

	clock.Update(lamport.Scalar{})
	require.Equal(t, lamport.Scalar{1, 2}, clock.Next(), "expected next timestamp to be 2")

	clock.Update(lamport.Scalar{2, 2})
	require.Equal(t, lamport.Scalar{1, 3}, clock.Next(), "expected next timestamp to be 3")

	clock.Update(lamport.Scalar{2, 6})
	require.Equal(t, lamport.Scalar{1, 7}, clock.Next(), "expected next timestamp to be 7")
}

func TestClockConcurrency(t *testing.T) {
	// Test concurrent clock operations by running a large number of threads with
	// indpendent read and write clocks and ensure that the next version number is
	// always after the previous version number when generated.
	var (
		wg, rg   sync.WaitGroup
		failures int64
	)

	// Create broadcast channel to simulate synchronization
	C := make([]chan *lamport.Scalar, 0, 3)

	for i := 1; i < 4; i++ {
		bc := make(chan *lamport.Scalar, 1)

		wg.Add(1)
		go func(wg *sync.WaitGroup, c chan *lamport.Scalar) {
			defer wg.Done()
			defer close(c)

			var sg sync.WaitGroup
			clock := lamport.New(uint32(i))

			// Kick off the receiver routine
			rg.Add(1)
			go func(rg *sync.WaitGroup, c <-chan *lamport.Scalar) {
				defer rg.Done()
				for v := range c {
					clock.Update(*v)
				}
			}(&rg, c)

			// Kick off several updater routines
			for j := 0; j < 3; j++ {
				sg.Add(1)
				go func(sg *sync.WaitGroup) {
					defer sg.Done()
					prev := &lamport.Scalar{}
					for k := 0; k < 16; k++ {
						time.Sleep(time.Millisecond * time.Duration(rand.Int63n(32)+2))

						now := clock.Next()
						if !now.After(prev) {
							atomic.AddInt64(&failures, 1)
							return
						}

						// Broadcast the update
						prev = &now
						for _, c := range C {
							c <- prev
						}
					}
				}(&sg)
			}

			sg.Wait()
		}(&wg, bc)
	}

	wg.Wait()
	for _, c := range C {
		close(c)
	}
	rg.Wait()

	require.Zero(t, failures, "expected zero failures to have occurred")
}
