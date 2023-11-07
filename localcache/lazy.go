package localcache

import (
	"context"
	"sync"
	"sync/atomic"
)

type LazyLoader[T any] struct {
	t        atomic.Value
	mu       sync.Mutex
	supplier Supplier[T]
}

func NewLazyLoader[T any](supplier Supplier[T]) (*LazyLoader[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	return &LazyLoader[T]{
		t:        atomic.Value{},
		mu:       sync.Mutex{},
		supplier: supplier,
	}, nil
}

func (l *LazyLoader[T]) Load(ctx context.Context) (T, error) {
	v := l.t.Load()
	if v != nil {
		return v.(T), nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	v = l.t.Load()
	if v != nil {
		return v.(T), nil
	}
	t, err := l.supplier(ctx)
	if err != nil {
		return t, err
	}
	l.t.Store(t)
	return t, nil
}
