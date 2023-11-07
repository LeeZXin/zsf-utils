package hashmap

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

type LinkedHashMap[K comparable, V any] struct {
	m Map[K, V]
	e Map[K, *list.Element]
	l *list.List

	limitSize   int
	accessOrder bool

	whenRemoveEldestKey atomic.Value
}

func NewLinkedHashMap[K comparable, V any](accessOrder bool) *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		m: NewHashMap[K, V](),
		e: NewHashMap[K, *list.Element](),
		l: list.New(),

		accessOrder: accessOrder,
	}
}

func NewLinkedHashMapWithLimitSize[K comparable, V any](accessOrder bool, limitSize int) *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{
		m: NewHashMap[K, V](),
		e: NewHashMap[K, *list.Element](),
		l: list.New(),

		limitSize:   limitSize,
		accessOrder: accessOrder,
	}
}

func (s *LinkedHashMap[K, V]) WhenRemoveEldestKey(fn func(K, V)) {
	s.whenRemoveEldestKey.Store(fn)
}

func (s *LinkedHashMap[K, V]) Get(k K) (V, bool) {
	v, b := s.m.Get(k)
	if b {
		if s.accessOrder {
			e, _ := s.e.Get(k)
			s.l.Remove(e)
			e = s.l.PushBack(k)
			s.e.Put(k, e)
		}
		return v, true
	}
	return v, false
}

func (s *LinkedHashMap[K, V]) Put(k K, v V) {
	e, b := s.e.Get(k)
	if b {
		if !s.accessOrder {
			s.l.Remove(e)
			e = s.l.PushBack(k)
			s.e.Put(k, e)
		}
	} else {
		e = s.l.PushBack(k)
		s.e.Put(k, e)
	}
	s.m.Put(k, v)
	if s.limitSize > 0 && s.limitSize < s.Size() {
		n := s.l.Front()
		nk := n.Value.(K)
		nv, _ := s.m.Get(nk)
		s.l.Remove(n)
		s.e.Remove(nk)
		s.m.Remove(nk)
		f := s.whenRemoveEldestKey.Load()
		if f != nil {
			f.(func(K, V))(nk, nv)
		}
	}
}

func (s *LinkedHashMap[K, V]) Contains(k K) bool {
	_, b := s.m.Get(k)
	return b
}

func (s *LinkedHashMap[K, V]) Remove(ks ...K) {
	for _, k := range ks {
		if s.Contains(k) {
			e, _ := s.e.Get(k)
			s.e.Remove(k)
			s.l.Remove(e)
			s.m.Remove(k)
		}
	}
}

func (s *LinkedHashMap[K, V]) AllKeys() []K {
	ret := make([]K, 0, s.m.Size())
	n := s.l.Front()
	for n != nil {
		ret = append(ret, n.Value.(K))
		n = n.Next()
	}
	return ret
}

func (s *LinkedHashMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	keys := s.AllKeys()
	for _, k := range keys {
		v, _ := s.m.Get(k)
		if !fn(k, v) {
			return
		}
	}
}

func (s *LinkedHashMap[K, V]) Size() int {
	return s.m.Size()
}

func (s *LinkedHashMap[K, V]) Clear() {
	s.l.Init()
	s.e.Clear()
	s.m.Clear()
}

func (s *LinkedHashMap[K, V]) GetOrDefault(k K, v V) V {
	ret, b := s.Get(k)
	if b {
		return ret
	}
	return v
}

func (s *LinkedHashMap[K, V]) GetOrPut(k K, v V) (V, bool) {
	ret, b := s.Get(k)
	if b {
		return ret, true
	}
	s.Put(k, v)
	return v, false
}

func (s *LinkedHashMap[K, V]) GetOrPutWithLoader(k K, fn func() (V, error)) (V, bool, error) {
	v, b := s.Get(k)
	if b {
		return v, true, nil
	}
	var err error
	if fn == nil {
		err = errors.New("nil fn")
	} else {
		v, err = fn()
		if err == nil {
			s.Put(k, v)
		}
	}
	return v, false, err
}

func (s *LinkedHashMap[K, V]) ToMap() map[K]V {
	return s.m.ToMap()
}

type ConcurrentLinkedHashMap[K comparable, V any] struct {
	m  *LinkedHashMap[K, V]
	mu sync.RWMutex

	accessOrder bool
}

func NewConcurrentLinkedHashMap[K comparable, V any](accessOrder bool) *ConcurrentLinkedHashMap[K, V] {
	return &ConcurrentLinkedHashMap[K, V]{
		m:  NewLinkedHashMap[K, V](accessOrder),
		mu: sync.RWMutex{},

		accessOrder: accessOrder,
	}
}

func NewConcurrentLinkedHashMapWithLimitSize[K comparable, V any](accessOrder bool, limitSize int) *ConcurrentLinkedHashMap[K, V] {
	return &ConcurrentLinkedHashMap[K, V]{
		m:  NewLinkedHashMapWithLimitSize[K, V](accessOrder, limitSize),
		mu: sync.RWMutex{},

		accessOrder: accessOrder,
	}
}

func (m *ConcurrentLinkedHashMap[K, V]) WhenRemoveEldestKey(fn func(K, V)) {
	m.m.WhenRemoveEldestKey(fn)
}

func (m *ConcurrentLinkedHashMap[K, V]) Get(k K) (V, bool) {
	if m.accessOrder {
		m.mu.Lock()
		defer m.mu.Unlock()
	} else {
		m.mu.RLock()
		defer m.mu.RUnlock()
	}
	return m.m.Get(k)
}

func (m *ConcurrentLinkedHashMap[K, V]) Put(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Put(k, v)
}

func (m *ConcurrentLinkedHashMap[K, V]) Contains(k K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Contains(k)
}

func (m *ConcurrentLinkedHashMap[K, V]) Remove(ks ...K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Remove(ks...)
}

func (m *ConcurrentLinkedHashMap[K, V]) AllKeys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.AllKeys()
}

func (m *ConcurrentLinkedHashMap[K, V]) Range(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Range(fn)
}

func (m *ConcurrentLinkedHashMap[K, V]) RangeWithRLock(fn func(K, V) bool) {
	if fn == nil {
		return
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.m.Range(fn)
}

func (m *ConcurrentLinkedHashMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.Size()
}

func (m *ConcurrentLinkedHashMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m.Clear()
}

func (m *ConcurrentLinkedHashMap[K, V]) GetOrDefault(k K, v V) V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.GetOrDefault(k, v)
}

func (m *ConcurrentLinkedHashMap[K, V]) GetOrPut(k K, v V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.GetOrPut(k, v)
}

func (m *ConcurrentLinkedHashMap[K, V]) GetOrPutWithLoader(k K, fn func() (V, error)) (V, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.m.GetOrPutWithLoader(k, fn)
}

func (m *ConcurrentLinkedHashMap[K, V]) ToMap() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m.ToMap()
}
