package localcache

import (
	"context"
	"github.com/LeeZXin/zsf-utils/collections/hashmap"
	"hash/crc32"
	"time"
)

// 分段segment 可过期map
// 默认分64个segment

const (
	segmentSize = 64
)

type segment[T any] struct {
	expireDuration time.Duration
	cache          hashmap.Map[string, *SingleCacheEntry[T]]
	supplier       SupplierWithKey[T]
}

func newSegment[T any](supplier SupplierWithKey[T], expireDuration time.Duration) *segment[T] {
	return &segment[T]{
		cache:          hashmap.NewConcurrentHashMap[string, *SingleCacheEntry[T]](),
		supplier:       supplier,
		expireDuration: expireDuration,
	}
}

func (e *segment[T]) getEntry(key string) *SingleCacheEntry[T] {
	ret, _, _ := e.cache.GetOrPutWithLoader(key, func() (*SingleCacheEntry[T], error) {
		return NewSingleCacheEntry(func(ctx context.Context) (T, error) {
			return e.supplier(ctx, key)
		}, e.expireDuration)
	})
	return ret
}

func (e *segment[T]) allKeys() []string {
	return e.cache.AllKeys()
}

func (e *segment[T]) clear() {
	e.cache.Clear()
}

func (e *segment[T]) removeKey(key string) {
	e.cache.Remove(key)
}

func (e *segment[T]) containsKey(key string) bool {
	return e.cache.Contains(key)
}

type LocalCache[T any] struct {
	supplier SupplierWithKey[T]
	segments []*segment[T]
}

func NewLocalCache[T any](supplier SupplierWithKey[T], duration time.Duration) (*LocalCache[T], error) {
	if supplier == nil {
		return nil, NilSupplierErr
	}
	segments := make([]*segment[T], 0, segmentSize)
	for i := 0; i < segmentSize; i++ {
		segments = append(segments, newSegment(supplier, duration))
	}
	return &LocalCache[T]{
		segments: segments,
		supplier: supplier,
	}, nil
}

func (e *LocalCache[T]) LoadData(ctx context.Context, key string) (T, error) {
	return e.getSegment(key).getEntry(key).LoadData(ctx)
}

func (e *LocalCache[T]) getSegment(key string) *segment[T] {
	// mod 64
	index := hash(key) & 0x3f
	return e.segments[index]
}

func (e *LocalCache[T]) RemoveKey(key string) {
	e.getSegment(key).removeKey(key)
}

func (e *LocalCache[T]) AllKeys() []string {
	ret := make([]string, 0)
	for _, seg := range e.segments {
		ret = append(ret, seg.allKeys()...)
	}
	return ret
}

func (e *LocalCache[T]) Clear() {
	for _, seg := range e.segments {
		seg.clear()
	}
}

func (e *LocalCache[T]) ContainsKey(key string) bool {
	return e.getSegment(key).containsKey(key)
}

func hash(key string) int {
	ret := crc32.ChecksumIEEE([]byte(key))
	return int(ret)
}
