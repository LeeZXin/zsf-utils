package maputil

type ImmutableMap[K comparable, V any] struct {
	m map[K]V
}

func NewImmutableMap[K comparable, V any](m map[K]V) *ImmutableMap[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return &ImmutableMap[K, V]{
		m: m,
	}
}

func (i *ImmutableMap[K, V]) Get(k K) (V, bool) {
	var (
		v V
		b bool
	)
	if i.m != nil {
		v, b = i.m[k]
	}
	return v, b
}

func (i *ImmutableMap[K, V]) HasKey(k K) bool {
	_, b := i.Get(k)
	return b
}

func (i *ImmutableMap[K, V]) AllKeys() []K {
	ks := make([]K, 0, len(i.m))
	for k := range i.m {
		ks = append(ks, k)
	}
	return ks
}

func (i *ImmutableMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	for k, v := range i.m {
		if !fn(k, v) {
			return
		}
	}
}

func (i *ImmutableMap[K, V]) Size() int {
	return len(i.m)
}

func (i *ImmutableMap[K, V]) GetOrDefault(k K, v V) V {
	ret, b := i.Get(k)
	if b {
		return ret
	}
	return v
}
