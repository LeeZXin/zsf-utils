package listutil

import (
	"encoding/json"
	"errors"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/randutil"
	"math"
)

type ComparableList[T comparable] []T

func (l *ComparableList[T]) Contains(l2 T) bool {
	for _, item := range *l {
		if item == l2 {
			return true
		}
	}
	return false
}

func (l *ComparableList[T]) FromDB(content []byte) error {
	if l == nil {
		*l = make([]T, 0)
	}
	return json.Unmarshal(content, l)
}

func (l *ComparableList[T]) ToDB() ([]byte, error) {
	return json.Marshal(l)
}

type List[T any] []T

func (l *List[T]) FromDB(content []byte) error {
	if l == nil {
		*l = make([]T, 0)
	}
	return json.Unmarshal(content, l)
}

func (l *List[T]) ToDB() ([]byte, error) {
	return json.Marshal(l)
}

func Contains[T any](arr []T, fn func(T) (bool, error)) (bool, error) {
	_, b, err := FindFirst(arr, fn)
	return b, err
}

func ContainsNe[T any](arr []T, fn func(T) bool) bool {
	_, b := FindFirstNe(arr, fn)
	return b
}

func FindFirst[T any](arr []T, fn func(T) (bool, error)) (T, bool, error) {
	if fn == nil {
		var t T
		return t, false, errors.New("nil fn")
	}
	for _, t := range arr {
		b, err := fn(t)
		if err != nil {
			return t, false, err
		}
		if b {
			return t, true, nil
		}
	}
	var t T
	return t, false, nil
}

func FindFirstNe[T any](arr []T, fn func(T) bool) (T, bool) {
	if fn == nil {
		var t T
		return t, false
	}
	for _, t := range arr {
		if fn(t) {
			return t, true
		}
	}
	var t T
	return t, false
}

func Filter[T any](data []T, fn func(T) (bool, error)) ([]T, error) {
	if fn == nil {
		return nil, errors.New("nil filter")
	}
	ret := make([]T, 0)
	for _, d := range data {
		b, err := fn(d)
		if err != nil {
			return nil, err
		}
		if b {
			ret = append(ret, d)
		}
	}
	return ret, nil
}

func FilterNe[T any](data []T, fn func(T) bool) []T {
	if fn == nil {
		return nil
	}
	ret := make([]T, 0)
	for _, d := range data {
		if fn(d) {
			ret = append(ret, d)
		}
	}
	return ret
}

func Map[T, K any](data []T, mapper func(T) (K, error)) ([]K, error) {
	if mapper == nil {
		return nil, errors.New("nil mapper")
	}
	ret := make([]K, 0, len(data))
	for _, d := range data {
		k, err := mapper(d)
		if err != nil {
			return nil, err
		}
		ret = append(ret, k)
	}
	return ret, nil
}

// MapNe Map with no error
func MapNe[T, K any](data []T, mapper func(T) K) []K {
	if mapper == nil {
		return nil
	}
	ret := make([]K, 0, len(data))
	for _, d := range data {
		ret = append(ret, mapper(d))
	}
	return ret
}

func MapWithIndex[T, K any](data []T, mapper func(T, int) (K, error)) ([]K, error) {
	if mapper == nil {
		return nil, errors.New("nil mapper")
	}
	ret := make([]K, 0, len(data))
	for i, d := range data {
		k, err := mapper(d, i)
		if err != nil {
			return nil, err
		}
		ret = append(ret, k)
	}
	return ret, nil
}

// MapNeWithIndex MapWithIndex with no error
func MapNeWithIndex[T, K any](data []T, mapper func(T, int) K) []K {
	if mapper == nil {
		return nil
	}
	ret := make([]K, 0, len(data))
	for i, d := range data {
		ret = append(ret, mapper(d, i))
	}
	return ret
}

func Distinct[T comparable](data ...T) []T {
	return hashset.NewHashSet(data...).AllKeys()
}

func Partition[T any](data []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}
	if data == nil {
		return [][]T{}
	}
	ret := make([][]T, 0, int(math.Ceil(float64(len(data)/size))))
	start := 0
	for start+size < len(data) {
		ret = append(ret, data[start:start+size])
		start += size
	}
	if start < len(data) {
		ret = append(ret, data[start:])
	}
	return ret
}

func CollectToMap[T any, N comparable, K any](data []T, nameFn func(T) (N, error), valFn func(T) (K, error)) (map[N]K, error) {
	if nameFn == nil {
		return nil, errors.New("nil nameFn")
	}
	if valFn == nil {
		return nil, errors.New("nil valFn")
	}
	ret := make(map[N]K, len(data))
	for _, d := range data {
		name, err := nameFn(d)
		if err != nil {
			return nil, err
		}
		val, err := valFn(d)
		if err != nil {
			return nil, err
		}
		ret[name] = val
	}
	return ret, nil
}

func Copy[T any](data []T) []T {
	if data == nil {
		return nil
	}
	ret := make([]T, 0, len(data))
	for _, d := range data {
		ret = append(ret, d)
	}
	return ret
}

func Shuffle[T any](data []T) []T {
	if data == nil {
		return nil
	}
	data = Copy(data)
	ShuffleSelf(data)
	return data
}

func ShuffleSelf[T any](data []T) {
	if data == nil {
		return
	}
	length := len(data)
	tail := length - 1
	for i := 0; i < tail; i++ {
		j := i + randutil.Intn(length-i)
		data[i], data[j] = data[j], data[i]
	}
}
