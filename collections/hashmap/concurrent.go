package hashmap

import (
	"sync"
)

type ConcurrentHashMap[K comparable, V any] struct {
	m  *HashMap[K, V]
	mu sync.RWMutex
}

func NewConcurrentHashMap[K comparable, V any]() *ConcurrentHashMap[K, V] {
	return &ConcurrentHashMap[K, V]{
		m:  NewHashMap[K, V](),
		mu: sync.RWMutex{},
	}
}

func (m *ConcurrentHashMap[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Get(k)
}

func (m *ConcurrentHashMap[K, V]) Contains(k K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Contains(k)
}

func (m *ConcurrentHashMap[K, V]) Remove(ks ...K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Remove(ks...)
}

func (m *ConcurrentHashMap[K, V]) AllKeys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.AllKeys()
}

func (m *ConcurrentHashMap[K, V]) Put(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Put(k, v)
}

func (m *ConcurrentHashMap[K, V]) GetOrPut(k K, v V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.GetOrPut(k, v)
}

func (m *ConcurrentHashMap[K, V]) GetOrPutWithLoader(k K, fn func() (V, error)) (V, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.GetOrPutWithLoader(k, fn)
}

func (m *ConcurrentHashMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	keys := m.AllKeys()
	for _, key := range keys {
		val, b := m.Get(key)
		if b {
			if !fn(key, val) {
				return
			}
		}
	}
}

func (m *ConcurrentHashMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Size()
}

func (m *ConcurrentHashMap[K, V]) GetOrDefault(k K, v V) V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.GetOrDefault(k, v)
}

func (m *ConcurrentHashMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Clear()
}

func (m *ConcurrentHashMap[K, V]) ToMap() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.ToMap()
}
