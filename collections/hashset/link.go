package hashset

import (
	"container/list"
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"sync"
)

type LinkedHashSet[T comparable] struct {
	s *HashSet[T]
	e *hashmap.HashMap[T, *list.Element]
	l *list.List
}

func NewLinkedHashSet[T comparable]() *LinkedHashSet[T] {
	return &LinkedHashSet[T]{
		s: NewHashSet[T](nil),
		e: hashmap.NewHashMap[T, *list.Element](),
		l: list.New(),
	}
}

func (s *LinkedHashSet[T]) Add(ts ...T) {
	for _, t := range ts {
		if !s.Contains(t) {
			e := s.l.PushBack(t)
			s.e.Put(t, e)
			s.s.Add(t)
		}
	}
}

func (s *LinkedHashSet[T]) Remove(ts ...T) {
	for _, t := range ts {
		if s.Contains(t) {
			e, _ := s.e.Get(t)
			s.s.Remove(t)
			s.e.Remove(t)
			s.l.Remove(e)
		}
	}
}

func (s *LinkedHashSet[T]) Contains(t T) bool {
	return s.s.Contains(t)
}

func (s *LinkedHashSet[T]) AllKeys() []T {
	ret := make([]T, 0, s.s.Size())
	n := s.l.Front()
	for n != nil {
		ret = append(ret, n.Value.(T))
		n = n.Next()
	}
	return ret
}

func (s *LinkedHashSet[T]) Intersect(h Set[T]) Set[T] {
	ret := NewHashSet[T](nil)
	h.Range(func(t T) {
		if s.Contains(t) {
			ret.Add(t)
		}
	})
	return ret
}

func (s *LinkedHashSet[T]) Range(fn func(T)) {
	if fn == nil {
		return
	}
	keys := s.AllKeys()
	for _, key := range keys {
		fn(key)
	}
}

func (s *LinkedHashSet[T]) Size() int {
	return s.s.Size()
}

func (s *LinkedHashSet[T]) Clear() {
	s.s.Clear()
	s.e.Clear()
	s.l.Init()
}

type ConcurrentLinkedHashSet[T comparable] struct {
	s  *LinkedHashSet[T]
	mu sync.RWMutex
}

func NewConcurrentLinkedHashSet[T comparable]() *ConcurrentLinkedHashSet[T] {
	return &ConcurrentLinkedHashSet[T]{
		s:  NewLinkedHashSet[T](),
		mu: sync.RWMutex{},
	}
}

func (c *ConcurrentLinkedHashSet[T]) Add(ts ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range ts {
		c.s.Add(t)
	}
}

func (c *ConcurrentLinkedHashSet[T]) Remove(t ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Remove(t...)
}

func (c *ConcurrentLinkedHashSet[T]) Contains(t T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Contains(t)
}

func (c *ConcurrentLinkedHashSet[T]) AllKeys() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.AllKeys()
}

func (c *ConcurrentLinkedHashSet[T]) Intersect(h Set[T]) Set[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Intersect(h)
}

func (c *ConcurrentLinkedHashSet[T]) Range(fn func(T)) {
	if fn == nil {
		return
	}
	keys := c.AllKeys()
	for _, key := range keys {
		fn(key)
	}
}

func (c *ConcurrentLinkedHashSet[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Size()
}

func (c *ConcurrentLinkedHashSet[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Clear()
}
