package atomicutil

import "sync/atomic"

type Value[T any] struct {
	v atomic.Value
}

func NewValue[T any]() *Value[T] {
	v := atomic.Value{}
	return &Value[T]{
		v: v,
	}
}

func (v *Value[T]) Store(t T) {
	v.v.Store(t)
}

func (v *Value[T]) Load() (T, bool) {
	var t T
	val := v.v.Load()
	if val == nil {
		return t, false
	}
	return val.(T), true
}
