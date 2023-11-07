package localcache

import (
	"context"
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"time"
)

// lru缓存 带过期时间
// 双向链表 + map

type LRUCache[T any] struct {
	cache          *hashmap.ConcurrentLinkedHashMap[string, *SingleCacheEntry[T]]
	supplier       SupplierWithKey[T]
	expireDuration time.Duration
}

func NewLRUCache[T any](supplier SupplierWithKey[T], duration time.Duration, maxSize int) (*LRUCache[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	return &LRUCache[T]{
		cache: hashmap.NewConcurrentLinkedHashMapWithLimitSize[string, *SingleCacheEntry[T]](true, maxSize),

		supplier: supplier,

		expireDuration: duration,
	}, nil
}

func (c *LRUCache[T]) LoadData(ctx context.Context, key string) (T, error) {
	entry, _, _ := c.cache.GetOrPutWithLoader(key, func() (*SingleCacheEntry[T], error) {
		return NewSingleCacheEntry(func(ctx context.Context) (T, error) {
			return c.supplier(ctx, key)
		}, c.expireDuration)
	})
	return entry.LoadData(ctx)
}

func (c *LRUCache[T]) AllKeys() []string {
	return c.cache.AllKeys()
}

func (c *LRUCache[T]) RemoveKey(key string) {
	c.cache.Remove(key)
}

func (c *LRUCache[T]) ContainsKey(key string) bool {
	return c.cache.Contains(key)
}

func (c *LRUCache[T]) Clear() {
	c.cache.Clear()
}
