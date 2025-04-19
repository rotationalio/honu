/*
The locks package implements a key-based lock mechanism that uses a crc32 hash to
distribute keys across a fixed number of locks. This allows for concurrent access
to different keys without contention, but fixes the number of locks (and therefore the
amount of available concurrency) to ensure that memory usage is bounded.

Locks help us implement simple transactions for the storage engine. At the start of a
transaction before any other operation is performed, the transaction must acquire
its read and write locks to all keys it wants to access in the transaction. The locks
are acquired in a consistent order to prevent deadlocks.

Reference: https://medium.com/@ksandeeptech07/handling-deadlocks-in-golang-gracefully-1f661c341a1d
*/
package locks

import (
	"hash/crc32"
	"sort"
	"sync"
)

const DefaultCount = 1024

type Keys interface {
	Lock(...[]byte)
	Unlock(...[]byte)
	RLock(...[]byte)
	RUnlock(...[]byte)
	Index([]byte) int
	Indices(...[]byte) []int
}

// KeyLocks are used to prevent concurrent writes to the same key and to allow multiple
// concurrent reads using a sync.RWMutex. The keys are distributed across a fixed number
// of mutexes to preven unbounded memory growth; so it is possible that two different
// keys will share the same lock. Collection keys (keys without object ids or versions),
// object prefix keys, and specific version keys are all locked with the same data
// structure.
//
// In a storage transaction, the transaction must acquire its read and write locks to
// all keys it wants to access in the transaction. The locks are acquired and released
// in a consistent order to prevent deadlocks. Multi-key locks must be acquired and
// released at the same time without multiple calls to Lock or RLock.
type KeyLock struct {
	count uint32
	locks []sync.RWMutex
	table *crc32.Table
}

var _ Keys = (*KeyLock)(nil)

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

// Acquire a write lock for the given keys. Note that the same exact keyset should be
// unlocked all at once using the Unlock function and multiple calls to Lock before
// unlocking should be avoided.
//
// The locks associated with the keys are identified by the index function, which is
// deduplicated and sorted so that the highest index lock is always acquired first. As
// long as the order of acquisition is consistent between different goroutines, you
// won't get deadlocks, even for arbitrary subsets of the locks.
func (k *KeyLock) Lock(keys ...[]byte) {
	switch len(keys) {
	case 0:
		// TODO: should a no key lock allow locking the entire keyspace?
		return
	case 1:
		k.locks[crc32.Checksum(keys[0], k.table)%k.count].Lock()
	default:
		indices := k.Indices(keys...)
		for _, index := range indices {
			k.locks[index].Lock()
		}
	}
}

// Unlock the write lock for the specified keys. Note that the same exact keyset should
// be unlocked all at once using the Unlock function before any calls to Lock or RLock.
//
// The locks associated with the keys are identified by the index function, which is
// deduplicated and sorted so that the highest index lock is always acquired first. As
// long as the order of acquisition is consistent between different goroutines, you
// won't get deadlocks, even for arbitrary subsets of the locks.
func (k *KeyLock) Unlock(keys ...[]byte) {
	switch len(keys) {
	case 0:
		return
	case 1:
		k.locks[crc32.Checksum(keys[0], k.table)%k.count].Unlock()
	default:
		indices := k.Indices(keys...)
		for _, index := range indices {
			k.locks[index].Unlock()
		}
	}
}

// Acquire a read lock for the given keys. Multiple read locks can be acquired
// concurrently, but a write lock cannot be acquired while any read locks are held and
// no read locks can be acquired while a write lock is held. The same exact keyset
// should be used with RUnlock to release the read locks and multiple calls to RLock or
// Lock before unlocking should be avoided.
//
// The locks associated with the keys are identified by the index function, which is
// deduplicated and sorted so that the highest index lock is always acquired first. As
// long as the order of acquisition is consistent between different goroutines, you
// won't get deadlocks, even for arbitrary subsets of the locks.
func (k *KeyLock) RLock(keys ...[]byte) {
	switch len(keys) {
	case 0:
		// TODO: should a no key lock allow read locking the entire keyspace?
		return
	case 1:
		k.locks[crc32.Checksum(keys[0], k.table)%k.count].RLock()
	default:
		indices := k.Indices(keys...)
		for _, index := range indices {
			k.locks[index].RLock()
		}
	}
}

// Release a read lock for the specified keys. Note that the same exact keyset should
// be unlocked all at once using the RUnlock function before any calls to Lock or RLock.
//
// The locks associated with the keys are identified by the index function, which is
// deduplicated and sorted so that the highest index lock is always acquired first. As
// long as the order of acquisition is consistent between different goroutines, you
// won't get deadlocks, even for arbitrary subsets of the locks.
func (k *KeyLock) RUnlock(keys ...[]byte) {
	switch len(keys) {
	case 0:
		return
	case 1:
		k.locks[crc32.Checksum(keys[0], k.table)%k.count].RUnlock()
	default:
		indices := k.Indices(keys...)
		for _, index := range indices {
			k.locks[index].RUnlock()
		}
	}
}

// Return the index of the lock for the given key (useful for debugging).
func (k *KeyLock) Index(key []byte) int {
	return int(crc32.Checksum(key, k.table) % k.count)
}

// Return a sorted, deduplicated slice of indices for the specified keys.
func (k *KeyLock) Indices(keys ...[]byte) (out []int) {
	// Find all indices associated with the keys
	out = make([]int, len(keys))
	for i, key := range keys {
		out[i] = int(crc32.Checksum(key, k.table) % k.count)
	}

	// Return the indices if there are no duplicates
	if len(out) < 2 {
		return out
	}

	// Deduplicate the indices to ensure that within a keyspace, there is no locking
	// deadlock (e.g. if we ask for locks for two keys with the same index in the same
	// goroutine we would deadlock if we didn't deduplicate).
	// Slice sort based deduplication is faster than map based deduplication.
	sort.Slice(out, func(i, j int) bool { return out[i] > out[j] })
	var e = 1
	for i := 1; i < len(out); i++ {
		if out[i] == out[i-1] {
			continue
		}
		out[e] = out[i]
		e++
	}
	return out[:e]
}
