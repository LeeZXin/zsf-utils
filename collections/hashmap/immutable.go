package hashmap

type ImmutableMap[K comparable, V any] struct {
	m *HashMap[K, V]
}

func NewImmutableMap[K comparable, V any](m map[K]V) *ImmutableMap[K, V] {
	return &ImmutableMap[K, V]{
		m: NewHashMapWithMap(m),
	}
}

func (i *ImmutableMap[K, V]) Get(k K) (V, bool) {
	return i.m.Get(k)
}

func (i *ImmutableMap[K, V]) Contains(k K) bool {
	return i.m.Contains(k)
}

func (i *ImmutableMap[K, V]) AllKeys() []K {
	return i.m.AllKeys()
}

func (i *ImmutableMap[K, V]) Range(fn func(K, V)) {
	i.m.Range(fn)
}

func (i *ImmutableMap[K, V]) Size() int {
	return i.m.Size()
}

func (i *ImmutableMap[K, V]) GetOrDefault(k K, v V) V {
	return i.m.GetOrDefault(k, v)
}

func (i *ImmutableMap[K, V]) ToMap() map[K]V {
	return i.m.ToMap()
}
