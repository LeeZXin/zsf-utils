package localcache

import (
	"context"
	"errors"
)

var (
	NilSupplierErr    = errors.New("nil supplier")
	IllegalMaxSizeErr = errors.New("maxSize should greater than 0")
)

type Supplier[T any] func(context.Context) (T, error)

type SupplierWithKey[T any] func(context.Context, string) (T, error)

type ExpireCache[T any] interface {
	// LoadData 获取数据
	LoadData(context.Context, string) (T, error)
	// RemoveKey 删除key
	RemoveKey(string)
	// AllKeys 获取所有的key
	AllKeys() []string
	// Clear 清除
	Clear()
	// ContainsKey 包含某个key
	ContainsKey(string) bool
}
