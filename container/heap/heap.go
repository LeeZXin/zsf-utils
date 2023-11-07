package heap

import (
	"container/heap"
	"sync"
)

type Object[T any] interface {
	GetObject() T
	GetPriority() int64
}

type objectWrapper[T any] struct {
	Object[T]
	index int
}

type container[T any] struct {
	o []*objectWrapper[T]

	positive bool
}

func (c *container[T]) Len() int {
	return len(c.o)
}

func (c *container[T]) Less(i, j int) bool {
	if c.positive {
		return c.o[i].GetPriority() < c.o[j].GetPriority()
	}
	return c.o[i].GetPriority() > c.o[j].GetPriority()
}

func (c *container[T]) Swap(i, j int) {
	if i < 0 || j < 0 {
		return
	}
	c.o[i], c.o[j] = c.o[j], c.o[i]
	c.o[i].index = i
	c.o[j].index = j
}

func (c *container[T]) Push(x any) {
	n := c.Len()
	item := &objectWrapper[T]{
		Object: x.(Object[T]),
		index:  n,
	}
	c.o = append(c.o, item)
}

func (c *container[T]) Pop() any {
	if c.Len() == 0 {
		return nil
	}
	old := c.o
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	c.o = old[0 : n-1]
	return item.Object
}

type Heap[T any] struct {
	c *container[T]
}

func NewHeap[T any](positive bool) *Heap[T] {
	c := &container[T]{
		o:        make([]*objectWrapper[T], 0),
		positive: positive,
	}
	heap.Init(c)
	return &Heap[T]{
		c: c,
	}
}

func (h *Heap[T]) Push(t Object[T]) {
	if t == nil {
		return
	}
	heap.Push(h.c, t)
}

func (h *Heap[T]) Pop() (Object[T], bool) {
	pop := heap.Pop(h.c)
	if pop == nil {
		return nil, false
	}
	return pop.(Object[T]), true
}

func (h *Heap[T]) Peek() (Object[T], bool) {
	o := h.c.o
	if len(o) == 0 {
		return nil, false
	}
	return o[0].Object, true
}

type ConcurrentHeap[T any] struct {
	h  *Heap[T]
	mu sync.Mutex
}

func NewConcurrentHeap[T any](positive bool) *ConcurrentHeap[T] {
	return &ConcurrentHeap[T]{
		h: NewHeap[T](positive),
	}
}

func (h *ConcurrentHeap[T]) Push(t Object[T]) {
	if t == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.h.Push(t)
}

func (h *ConcurrentHeap[T]) Pop() (Object[T], bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.h.Pop()
}

func (h *ConcurrentHeap[T]) Peek() (Object[T], bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.h.Peek()
}
