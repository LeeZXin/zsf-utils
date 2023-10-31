package hashset

import "sync"

type ConcurrentHashSet[T comparable] struct {
	s  Set[T]
	mu sync.RWMutex
}

func NewConcurrentHashSet[T comparable](t []T) *ConcurrentHashSet[T] {
	return &ConcurrentHashSet[T]{
		s:  NewHashSet(t),
		mu: sync.RWMutex{},
	}
}

func NewConcurrentHashSetWithSet[T comparable](s Set[T]) *ConcurrentHashSet[T] {
	if s == nil {
		s = NewHashSet[T](nil)
	}
	return &ConcurrentHashSet[T]{
		s:  s,
		mu: sync.RWMutex{},
	}
}

func (c *ConcurrentHashSet[T]) Add(t T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Add(t)
}

func (c *ConcurrentHashSet[T]) Delete(t ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Delete(t...)
}

func (c *ConcurrentHashSet[T]) Contains(t T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Contains(t)
}

func (c *ConcurrentHashSet[T]) Copy() Set[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Copy()
}

func (c *ConcurrentHashSet[T]) AllKeys() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.AllKeys()
}

func (c *ConcurrentHashSet[T]) Intersect(h Set[T]) Set[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Intersect(h)
}

func (c *ConcurrentHashSet[T]) Range(fn func(T) bool) {
	if fn == nil {
		return
	}
	keys := c.AllKeys()
	for _, key := range keys {
		if !fn(key) {
			return
		}
	}
}
