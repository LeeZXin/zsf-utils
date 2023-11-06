package hashmap

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

func (i *ImmutableMap[K, V]) Contains(k K) bool {
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

func (i *ImmutableMap[K, V]) Remove(_ ...K) {}

func (i *ImmutableMap[K, V]) Put(_ K, _ V) {}

func (i *ImmutableMap[K, V]) GetOrPut(k K, v V) (V, bool) {
	return i.GetOrDefault(k, v), false
}

func (i *ImmutableMap[K, V]) GetOrPutWithLoader(k K, fn func() (V, error)) (V, bool, error) {
	ret, b := i.Get(k)
	if b {
		return ret, true, nil
	}
	v, err := fn()
	return v, false, err
}

func (i *ImmutableMap[K, V]) Clear() {}

func (i *ImmutableMap[K, V]) ToMap() map[K]V {
	ret := make(map[K]V, i.Size())
	i.Range(func(k K, v V) bool {
		ret[k] = v
		return true
	})
	return ret
}
