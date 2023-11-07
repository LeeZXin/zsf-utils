package hashmap

import "errors"

type Map[K comparable, V any] interface {
	Get(K) (V, bool)
	Put(K, V)
	Contains(K) bool
	Remove(...K)
	AllKeys() []K
	Range(func(K, V) bool)
	Size() int
	Clear()
	GetOrDefault(K, V) V
	GetOrPut(K, V) (V, bool)
	GetOrPutWithLoader(K, func() (V, error)) (V, bool, error)
	ToMap() map[K]V
}

type ConcurrentMap[K comparable, V any] interface {
	Map[K, V]
	RangeWithRLock(func(K, V) bool)
}

type HashMap[K comparable, V any] struct {
	m map[K]V
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		m: make(map[K]V, 8),
	}
}

func NewHashMapWithMap[K comparable, V any](m map[K]V) *HashMap[K, V] {
	ret := &HashMap[K, V]{
		m: make(map[K]V, len(m)),
	}
	for k, v := range m {
		ret.Put(k, v)
	}
	return ret
}

func (h *HashMap[K, V]) Get(k K) (V, bool) {
	v, b := h.m[k]
	return v, b
}

func (h *HashMap[K, V]) Put(k K, v V) {
	h.m[k] = v
}

func (h *HashMap[K, V]) Contains(k K) bool {
	_, b := h.Get(k)
	return b
}
func (h *HashMap[K, V]) Remove(ks ...K) {
	for _, k := range ks {
		delete(h.m, k)
	}
}
func (h *HashMap[K, V]) AllKeys() []K {
	ret := make([]K, 0, len(h.m))
	for k := range h.m {
		ret = append(ret, k)
	}
	return ret
}

func (h *HashMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	for k, v := range h.m {
		if !fn(k, v) {
			return
		}
	}
}

func (h *HashMap[K, V]) Size() int {
	return len(h.m)
}

func (h *HashMap[K, V]) Clear() {
	for k := range h.m {
		delete(h.m, k)
	}
}

func (h *HashMap[K, V]) GetOrDefault(k K, v V) V {
	ret, b := h.Get(k)
	if b {
		return ret
	}
	return v
}

func (h *HashMap[K, V]) GetOrPut(k K, v V) (V, bool) {
	ret, b := h.Get(k)
	if b {
		return ret, true
	}
	h.Put(k, v)
	return v, false
}

func (h *HashMap[K, V]) GetOrPutWithLoader(k K, fn func() (V, error)) (V, bool, error) {
	v, b := h.Get(k)
	if b {
		return v, true, nil
	}
	var err error
	if fn == nil {
		err = errors.New("nil loader")
	} else {
		v, err = fn()
		if err == nil {
			h.Put(k, v)
		}
	}
	return v, false, err
}

func (h *HashMap[K, V]) ToMap() map[K]V {
	ret := make(map[K]V, h.Size())
	h.Range(func(k K, v V) bool {
		ret[k] = v
		return true
	})
	return ret
}
