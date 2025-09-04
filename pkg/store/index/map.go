package index

type Map[K ordered, V any] map[K]V

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) Has(k K) bool {
	_, ok := m[k]
	return ok
}

func (m Map[K, V]) Get(k K) (V, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Map[K, V]) Put(k K, v V) (V, bool) {
	prev, ok := m[k]
	m[k] = v
	return prev, ok
}

func (m Map[K, V]) Del(k K) (V, bool) {
	v, ok := m[k]
	if ok {
		delete(m, k)
	}
	return v, ok
}
