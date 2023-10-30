package sortutil

import "sort"

type Comparable[T any] interface {
	CompareTo(T) bool
}

func SliceStable[T Comparable[T]](c []T) {
	if len(c) > 0 {
		sort.SliceStable(c, func(i, j int) bool {
			return c[i].CompareTo(c[j])
		})
	}
}

func SliceStableReverse[T Comparable[T]](c []T) {
	if len(c) > 0 {
		sort.SliceStable(c, func(i, j int) bool {
			return !c[i].CompareTo(c[j])
		})
	}
}
