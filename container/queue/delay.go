package queue

import (
	"github.com/LeeZXin/zsf-utils/container/heap"
	"sync"
	"time"
)

type Delayed[T any] interface {
	GetObject() T
	GetDelayedDuration() time.Duration
}

type objectWrapper[T any] struct {
	d Delayed[T]
	t time.Time
}

func (w *objectWrapper[T]) GetObject() T {
	return w.d.GetObject()
}

func (w *objectWrapper[T]) GetPriority() int64 {
	return int64(w.d.GetDelayedDuration())
}

func (w *objectWrapper[T]) GetRemainTime() time.Duration {
	return w.t.Add(w.d.GetDelayedDuration()).Sub(time.Now())
}

type DelayedQueue[T any] struct {
	h      *heap.Heap[T]
	mu     sync.Mutex
	notify chan struct{}
}

func NewDelayedQueue[T any]() *DelayedQueue[T] {
	return &DelayedQueue[T]{
		h:      heap.NewHeap[T](true),
		mu:     sync.Mutex{},
		notify: make(chan struct{}, 1),
	}
}

func (h *DelayedQueue[T]) Push(d Delayed[T]) {
	if d == nil {
		return
	}
	h.mu.Lock()
	h.h.Push(&objectWrapper[T]{
		d: d,
		t: time.Now(),
	})
	p, _ := h.peek()
	if p == d {
		h.mu.Unlock()
		select {
		case h.notify <- struct{}{}:
		default:
		}
	} else {
		h.mu.Unlock()
	}
}

func (h *DelayedQueue[T]) peek() (Delayed[T], bool) {
	p, b := h.h.Peek()
	if b {
		return p.(*objectWrapper[T]).d, true
	}
	return nil, false
}

func (h *DelayedQueue[T]) Take() Delayed[T] {
	for {
		h.mu.Lock()
		p, b := h.h.Peek()
		if b {
			w := p.(*objectWrapper[T])
			duration := w.GetRemainTime()
			if duration <= 0 {
				pop, _ := h.h.Pop()
				h.mu.Unlock()
				return pop.(*objectWrapper[T]).d
			} else {
				h.mu.Unlock()
				select {
				case <-time.After(duration):
					continue
				case <-h.notify:
					continue
				}
			}
		} else {
			h.mu.Unlock()
			select {
			case <-h.notify:
				continue
			}
		}
	}
}
