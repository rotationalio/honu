package locks_test

import (
	"crypto/rand"
	random "math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/locks"
)

func TestLocking(t *testing.T) {
	mu := locks.New(0)
	k1 := []byte{28, 38, 142, 19, 36, 13, 138, 224, 208, 120, 174, 202, 245, 249, 3, 170, 219, 30, 199, 222, 214, 181, 20, 190, 160, 228, 153, 169, 50, 227, 233, 158, 213, 202, 71, 43, 228, 172, 50, 245, 129, 145, 45, 241, 52}
	k2 := []byte{130, 255, 26, 113, 234, 106, 222, 190, 243, 205, 105, 236, 252, 191, 170, 22, 103, 99, 205, 208, 252, 137, 32, 197, 180, 40, 46, 90, 120, 224, 109, 215, 124, 190, 67, 208, 114, 208, 64, 13, 144, 88, 188, 112, 144}
	k3 := []byte{43, 190, 204, 198, 110, 31, 120, 141, 165, 0, 172, 63, 39, 97, 37, 6, 189, 186, 23, 208, 124, 212, 163, 232, 153, 81, 105, 133, 147, 247, 192, 66, 253, 243, 163, 225, 35, 22, 230, 42, 55, 4, 83, 74, 240}

	t.Run("NoKey", func(t *testing.T) {
		mu.Lock()
		mu.Unlock()

		mu.RLock()
		mu.RLock()
		mu.RLock()

		mu.RUnlock()
		mu.RUnlock()
		mu.RUnlock()
	})

	t.Run("Single", func(t *testing.T) {
		// Make sure the three keys have different indexes or this test will fail.
		require.NotEqual(t, mu.Index(k1), mu.Index(k2))
		require.NotEqual(t, mu.Index(k1), mu.Index(k3))
		require.NotEqual(t, mu.Index(k2), mu.Index(k3))

		mu.Lock(k1)
		mu.Lock(k2)
		mu.Lock(k3)

		mu.Unlock(k2)
		mu.RLock(k2)
		mu.RLock(k2)
		mu.RLock(k2)

		mu.RUnlock(k2)
		mu.RUnlock(k2)
		mu.RUnlock(k2)

		mu.Unlock(k1)
		mu.Unlock(k3)
	})

	t.Run("Multiple", func(t *testing.T) {
		mu.Lock(k1, k2, k3)
		mu.Unlock(k1, k2, k3)

		mu.RLock(k1, k2, k3)
		mu.RLock(k1, k2, k3)
		mu.RLock(k1, k2, k3)

		mu.RUnlock(k1, k2, k3)
		mu.RUnlock(k1, k2, k3)
		mu.RUnlock(k1, k2, k3)
	})
}

func TestContention(t *testing.T) {
	var (
		writers    = 64
		writes     = 256
		writesleep = 32
		readers    = 256
		reads      = 512
		readsleep  = 16
		keyspace   = 512
		locksize   = 128
	)

	t.Run("NoKey", func(t *testing.T) {
		mu := locks.New(uint32(locksize))
		wg := sync.WaitGroup{}

		wg.Add(writers)
		for range writers {
			go func() {
				defer wg.Done()
				for range writes {
					mu.Lock()
					time.Sleep(time.Duration(random.IntN(writesleep)) * time.Microsecond)
					mu.Unlock()
				}
			}()
		}

		wg.Add(readers)
		for range readers {
			go func() {
				defer wg.Done()
				for range reads {
					mu.RLock()
					time.Sleep(time.Duration(random.IntN(readsleep)) * time.Microsecond)
					mu.RUnlock()
				}
			}()
		}

		wg.Wait()
	})

	t.Run("Single", func(t *testing.T) {
		mu := locks.New(uint32(locksize))
		wg := sync.WaitGroup{}

		keys := make([][]byte, keyspace)
		for i := range keyspace {
			keys[i] = RandomKey()
		}

		wg.Add(writers)
		for range writers {
			go func() {
				defer wg.Done()
				k := keys[random.IntN(len(keys))]

				for range writes {
					mu.Lock(k)
					time.Sleep(time.Duration(random.IntN(writesleep)) * time.Microsecond)
					mu.Unlock(k)
				}
			}()
		}

		wg.Add(readers)
		for range readers {
			go func() {
				defer wg.Done()
				k := keys[random.IntN(len(keys))]

				for range reads {
					mu.RLock(k)
					time.Sleep(time.Duration(random.IntN(readsleep)) * time.Microsecond)
					mu.RUnlock(k)
				}
			}()
		}

		wg.Wait()
	})

	t.Run("Multiple", func(t *testing.T) {
		mu := locks.New(uint32(locksize))
		wg := sync.WaitGroup{}

		keys := make([][]byte, keyspace)
		for i := range keyspace {
			keys[i] = RandomKey()
		}

		wg.Add(writers)
		for range writers {
			go func() {
				defer wg.Done()

				n := random.IntN(32)
				ks := make([][]byte, n)
				for j := range n {
					ks[j] = keys[random.IntN(len(keys))]
				}

				for range writes {
					mu.Lock(ks...)
					time.Sleep(time.Duration(random.IntN(writesleep)) * time.Microsecond)
					mu.Unlock(ks...)
				}
			}()
		}

		wg.Add(readers)
		for range readers {
			go func() {
				defer wg.Done()

				n := random.IntN(32)
				ks := make([][]byte, n)
				for j := range n {
					ks[j] = keys[random.IntN(len(keys))]
				}

				for range reads {
					mu.RLock(ks...)
					time.Sleep(time.Duration(random.IntN(readsleep)) * time.Microsecond)
					mu.RUnlock(ks...)
				}
			}()
		}

		wg.Wait()
	})
}

func TestIndices(t *testing.T) {
	t.Run("NoDuplicates", func(t *testing.T) {
		mu := locks.New(128)
		keys := [][]byte{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}

		indices := mu.Indices(keys...)
		require.Equal(t, []int{81, 42, 20}, indices)
	})

	t.Run("Duplicates", func(t *testing.T) {
		mu := locks.New(8)
		keys := make([][]byte, 128)
		for i := 0; i < 128; i++ {
			keys[i] = RandomKey()
		}

		indices := mu.Indices(keys...)
		require.Equal(t, []int{7, 6, 5, 4, 3, 2, 1, 0}, indices)
	})

}

func BenchmarkLocks(b *testing.B) {
	keys := make([][]byte, 128)
	for i := 0; i < 128; i++ {
		keys[i] = RandomKey()
	}

	runContention := func(f func([]byte)) {
		wg := sync.WaitGroup{}
		wg.Add(1024)
		for i := 0; i < 1024; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 512; j++ {
					k := keys[random.IntN(128)]
					f(k)
				}
			}()
		}
		wg.Wait()
	}

	b.Run("Sync", func(b *testing.B) {
		var mu sync.Mutex
		for i := 0; i < b.N; i++ {
			runContention(func(k []byte) {
				mu.Lock()
				_ = len(k)
				mu.Unlock()
			})
		}
	})

	b.Run("KeyLock128", func(b *testing.B) {
		mu := locks.New(128)
		for i := 0; i < b.N; i++ {
			runContention(func(k []byte) {
				mu.Lock(k)
				_ = len(k)
				mu.Unlock(k)
			})
		}
	})

	b.Run("KeyLock1024", func(b *testing.B) {
		mu := locks.New(1024)
		for i := 0; i < b.N; i++ {
			runContention(func(k []byte) {
				mu.Lock(k)
				_ = len(k)
				mu.Unlock(k)
			})
		}
	})

}

func RandomKey() []byte {
	k := make([]byte, 45)
	rand.Read(k)
	return k
}
