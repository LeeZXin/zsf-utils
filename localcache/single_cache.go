package localcache

import (
	"context"
	"github.com/LeeZXin/zsf-utils/atomicutil"
	"sync"
	"time"
)

// SingleCache 单个数据缓存
// 带过期时间

type SingleCacheEntry[T any] struct {
	expireDuration time.Duration
	expireTime     *atomicutil.Value[time.Time]
	data           *atomicutil.Value[T]
	mu             sync.Mutex
	supplier       Supplier[T]
}

func NewSingleCacheEntry[T any](supplier Supplier[T], duration time.Duration) (*SingleCacheEntry[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	return &SingleCacheEntry[T]{
		expireDuration: duration,
		expireTime:     atomicutil.NewValue[time.Time](),
		data:           atomicutil.NewValue[T](),
		mu:             sync.Mutex{},
		supplier:       supplier,
	}, nil
}

func (e *SingleCacheEntry[T]) LoadData(ctx context.Context) (T, error) {
	var (
		result T
		err    error
	)
	etime, b := e.expireTime.Load()
	// 首次加载
	if !b {
		e.mu.Lock()
		defer e.mu.Unlock()
		if _, b = e.expireTime.Load(); b {
			ret, _ := e.data.Load()
			return ret, nil
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
	ret, _ := e.data.Load()
	return ret, nil
}
