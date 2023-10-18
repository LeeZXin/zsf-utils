package localcache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type AutoCleanCache[T any] struct {
	c              *cache.Cache
	mu             sync.Mutex
	expireDuration time.Duration
	supplier       SupplierWithKey[T]
}

func NewAutoCleanCache[T any](supplier SupplierWithKey[T], expireDuration time.Duration, cleanupDuration time.Duration) (*AutoCleanCache[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	return &AutoCleanCache[T]{
		c:              cache.New(expireDuration, cleanupDuration),
		mu:             sync.Mutex{},
		expireDuration: expireDuration,
		supplier:       supplier,
	}, nil
}

// LoadData 获取数据
func (c *AutoCleanCache[T]) LoadData(ctx context.Context, key string) (T, error) {
	r, b := c.c.Get(key)
	if b {
		return r.(T), nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	r, b = c.c.Get(key)
	if b {
		return r.(T), nil
	}
	ret, err := c.supplier(ctx, key)
	if err != nil {
		var t T
		return t, nil
	}
	c.c.Set(key, ret, c.expireDuration)
	return ret, nil
}

// RemoveKey 删除key
func (c *AutoCleanCache[T]) RemoveKey(key string) {
	c.c.Delete(key)
}

// AllKeys 获取所有的key
func (c *AutoCleanCache[T]) AllKeys() []string {
	items := c.c.Items()
	ret := make([]string, 0, len(items))
	for key := range items {
		ret = append(ret, key)
	}
	return ret
}

// Clear 清除
func (c *AutoCleanCache[T]) Clear() {
	c.c.Flush()
}

// ContainsKey 包含某个key
func (c *AutoCleanCache[T]) ContainsKey(key string) bool {
	_, b := c.c.Get(key)
	return b
}
