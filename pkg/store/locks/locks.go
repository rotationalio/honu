/*
The locks package implements a key-based lock mechanism that uses a crc32 hash to
distribute keys across a fixed number of locks. This allows for concurrent access
to different keys without contention, but fixes the number of locks (and therefore the
amount of available concurrency) to ensure that memory usage is bounded.
*/
package locks

import (
	"hash/crc32"
	"sync"
)

const DefaultCount = 1024

type Keys interface {
	Lock([]byte)
	Unlock([]byte)
	RLock([]byte)
	RUnlock([]byte)
	Index([]byte) int
}

// KeyLocks are used to prevent concurrent writes to the same key and to allow multiple
// concurrent reads using a sync.RWMutex. The keys are distributed across a fixed number
// to preven unbounded memory growth; so it is possible that two different keys will
// share the same lock. Collection keys (keys without object ids or versions), object
// prefix keys, and specific version keys are all locked with the same data structure.
type KeyLock struct {
	count uint32
	locks []sync.RWMutex
	table *crc32.Table
}

// Create a new KeyLock with the given number of locks. The greater nlocks is, the
// greater concurrency there is across the entire key space at the cost of more memory.
// We recommend allocating at least 1024 locks to ensure suitable performance.
func New(nlocks uint32) *KeyLock {
	if nlocks == 0 {
		nlocks = DefaultCount
	}

	return &KeyLock{
		count: nlocks,
		locks: make([]sync.RWMutex, nlocks),
		table: crc32.MakeTable(crc32.Koopman),
	}
}

// Acquire a write lock for the given key.
func (k *KeyLock) Lock(key []byte) {
	k.locks[crc32.Checksum(key, k.table)%k.count].Lock()
}

// Unlock the write lock for the specified key.
func (k *KeyLock) Unlock(key []byte) {
	k.locks[crc32.Checksum(key, k.table)%k.count].Unlock()
}

// Acquire a read lock for the given key. Multiple read locks can be acquired
// concurrently, but a write lock cannot be acquired while any read locks are held.
func (k *KeyLock) RLock(key []byte) {
	k.locks[crc32.Checksum(key, k.table)%k.count].RLock()
}

// Release a read lock for the specified key.
func (k *KeyLock) RUnlock(key []byte) {
	k.locks[crc32.Checksum(key, k.table)%k.count].RUnlock()
}

// Return the index of the lock for the given key (useful for debugging).
func (k *KeyLock) Index(key []byte) int {
	return int(crc32.Checksum(key, k.table) % k.count)
}
