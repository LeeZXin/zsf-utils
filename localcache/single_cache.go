package localcache

import (
	"context"
	"github.com/LeeZXin/zsf-utils/concurrentutil"
	"sync"
	"time"
)

// SingleCache 单个数据缓存
// 带过期时间

type SingleCacheEntry[T any] struct {
	expireDuration time.Duration
	expireTime     *concurrentutil.Value[time.Time]
	data           *concurrentutil.Value[T]
	mu             sync.Mutex
	supplier       Supplier[T]
}

func NewSingleCacheEntry[T any](supplier Supplier[T], duration time.Duration) (*SingleCacheEntry[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	return &SingleCacheEntry[T]{
		expireDuration: duration,
		expireTime:     concurrentutil.NewValue[time.Time](),
		data:           concurrentutil.NewValue[T](),
		mu:             sync.Mutex{},
		supplier:       supplier,
	}, nil
}

func (e *SingleCacheEntry[T]) LoadData(ctx context.Context) (T, error) {
	var (
		result T
		err    error
	)
	etime := e.expireTime.Load()
	// 首次加载
	if etime.IsZero() {
		e.mu.Lock()
		defer e.mu.Unlock()
		if etime = e.expireTime.Load(); !etime.IsZero() {
			return e.data.Load(), nil
		}
		result, err = e.supplier(ctx)
		if err != nil {
			return result, err
		}
		e.data.Store(result)
		e.expireTime.Store(time.Now().Add(e.expireDuration))
		return result, nil
	}
	if e.expireDuration > 0 {
		if etime.Before(time.Now()) {
			// 过期
			if e.mu.TryLock() {
				defer e.mu.Unlock()
				result, err = e.supplier(ctx)
				if err == nil {
					e.data.Store(result)
					e.expireTime.Store(time.Now().Add(e.expireDuration))
				}
				return result, nil
			}
		}
	}
	return e.data.Load(), nil
}

func (e *SingleCacheEntry[T]) Clear() {
	e.expireTime.Store(time.Time{})
}
