package hashset

import "sync"

type ConcurrentHashSet[T comparable] struct {
	s  Set[T]
	mu sync.RWMutex
}

func NewConcurrentHashSet[T comparable]() *ConcurrentHashSet[T] {
	return &ConcurrentHashSet[T]{
		s:  NewHashSet[T](),
		mu: sync.RWMutex{},
	}
}

func (c *ConcurrentHashSet[T]) Add(ts ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range ts {
		c.s.Add(t)
	}
}

func (c *ConcurrentHashSet[T]) Remove(t ...T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Remove(t...)
}

func (c *ConcurrentHashSet[T]) Contains(t T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Contains(t)
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

func (c *ConcurrentHashSet[T]) Range(fn func(T)) {
	if fn == nil {
		return
	}
	keys := c.AllKeys()
	for _, key := range keys {
		fn(key)
	}
}

func (c *ConcurrentHashSet[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.s.Size()
}

func (c *ConcurrentHashSet[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Clear()
}
