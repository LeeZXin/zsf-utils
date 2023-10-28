package maputil

import (
	"errors"
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func NewConcurrentMap[K comparable, V any](m map[K]V) *ConcurrentMap[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return &ConcurrentMap[K, V]{
		m:  m,
		mu: sync.RWMutex{},
	}
}

func (m *ConcurrentMap[K, V]) Load(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, b := m.m[k]
	return v, b
}

func (m *ConcurrentMap[K, V]) HasKey(k K) bool {
	_, b := m.Load(k)
	return b
}

func (m *ConcurrentMap[K, V]) RemoveKey(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, k)
}

func (m *ConcurrentMap[K, V]) AllKeys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ret := make([]K, 0, len(m.m))
	for k := range m.m {
		ret = append(ret, k)
	}
	return ret
}

func (m *ConcurrentMap[K, V]) Store(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[k] = v
}

func (m *ConcurrentMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v2, b := m.m[k]
	if b {
		return v2, true
	}
	m.m[k] = v
	return v, false
}

func (m *ConcurrentMap[K, V]) LoadOrStoreWithLoader(k K, fn func() (V, error)) (V, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, b := m.m[k]
	if b {
		return v, nil
	}
	var err error
	if fn == nil {
		err = errors.New("nil loader")
	} else {
		v, err = fn()
		if err == nil {
			m.m[k] = v
		}
	}
	return v, err
}

func (m *ConcurrentMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	keys := m.AllKeys()
	for _, key := range keys {
		val, b := m.Load(key)
		if b {
			if !fn(key, val) {
				return
			}
		}
	}
}

func (m *ConcurrentMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

func (m *ConcurrentMap[K, V]) GetOrDefault(k K, v V) V {
	ret, b := m.Load(k)
	if b {
		return ret
	}
	return v
}
