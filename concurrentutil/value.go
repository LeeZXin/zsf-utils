package concurrentutil

import (
	"sync"
)

type Value[T any] struct {
	v T
	m sync.RWMutex
}

func NewValue[T any]() *Value[T] {
	return &Value[T]{
		m: sync.RWMutex{},
	}
}

func NewWithValue[T any](v T) *Value[T] {
	return &Value[T]{
		v: v,
		m: sync.RWMutex{},
	}
}

func (v *Value[T]) Store(t T) {
	v.m.Lock()
	defer v.m.Unlock()
	v.v = t
}

func (v *Value[T]) Load() T {
	v.m.RLock()
	defer v.m.RUnlock()
	return v.v
}
