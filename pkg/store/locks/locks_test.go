package locks_test

import (
	"crypto/rand"
	random "math/rand/v2"
	"sync"
	"testing"

	"github.com/rotationalio/honu/pkg/store/locks"
	"github.com/stretchr/testify/require"
)

func TestLocking(t *testing.T) {
	mu := locks.New(0)
	k1 := []byte{28, 38, 142, 19, 36, 13, 138, 224, 208, 120, 174, 202, 245, 249, 3, 170, 219, 30, 199, 222, 214, 181, 20, 190, 160, 228, 153, 169, 50, 227, 233, 158, 213, 202, 71, 43, 228, 172, 50, 245, 129, 145, 45, 241, 52}
	k2 := []byte{130, 255, 26, 113, 234, 106, 222, 190, 243, 205, 105, 236, 252, 191, 170, 22, 103, 99, 205, 208, 252, 137, 32, 197, 180, 40, 46, 90, 120, 224, 109, 215, 124, 190, 67, 208, 114, 208, 64, 13, 144, 88, 188, 112, 144}
	k3 := []byte{43, 190, 204, 198, 110, 31, 120, 141, 165, 0, 172, 63, 39, 97, 37, 6, 189, 186, 23, 208, 124, 212, 163, 232, 153, 81, 105, 133, 147, 247, 192, 66, 253, 243, 163, 225, 35, 22, 230, 42, 55, 4, 83, 74, 240}

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
}

func TestContention(t *testing.T) {
	mu := locks.New(128)
	wg := sync.WaitGroup{}

	wg.Add(1024)
	for i := 0; i < 1024; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 512; j++ {
				k := RandomKey()
				mu.Lock(k)
				mu.Unlock(k)
			}
		}()
	}

	wg.Add(2048)
	for i := 0; i < 2048; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1024; j++ {
				k := RandomKey()
				mu.RLock(k)
				mu.RUnlock(k)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkLocks(b *testing.B) {
	mu := sync.Mutex{}
	keys := make([][]byte, 128)
	for i := 0; i < 1024; i++ {
		keys[i] = RandomKey()
	}

	runContention := func() {
		wg := sync.WaitGroup{}
		wg.Add(1024)
		for i := 0; i < 1024; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 512; j++ {
					k := keys[random.IntN(1024)]
					mu.Lock()
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		mu.Unlock()
	}
}

func RandomKey() []byte {
	k := make([]byte, 45)
	rand.Read(k)
	return k
}
