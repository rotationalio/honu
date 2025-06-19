package index

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string | ~[16]byte
}

// An Index provides high-performance in-memory lookups and inserts for queries and
// other database operations. There are multiple types of indexes available that all
// share this common interface.
type Index[K ordered, V any] interface {
	Len() int
	Has(K) bool
	Get(K) (V, bool)
	Put(K, V) (V, bool)
	Del(K) (V, bool)
}
