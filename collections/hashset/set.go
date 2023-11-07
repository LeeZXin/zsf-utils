package hashset

import "github.com/LeeZXin/zsf-utils/collections/hashmap"

type Set[T comparable] interface {
	Add(...T)
	Remove(...T)
	Contains(T) bool
	AllKeys() []T
	Intersect(h Set[T]) Set[T]
	Range(func(T) bool)
	Size() int
	Clear()
}

type HashSet[T comparable] struct {
	m *hashmap.HashMap[T, struct{}]
}

func NewHashSet[T comparable](ts []T) *HashSet[T] {
	ret := &HashSet[T]{
		m: hashmap.NewHashMap[T, struct{}](),
	}
	ret.Add(ts...)
	return ret
}

func (s *HashSet[T]) Add(ts ...T) {
	for _, t := range ts {
		s.m.Put(t, struct{}{})
	}
}

func (s *HashSet[T]) Remove(ts ...T) {
	s.m.Remove(ts...)
}

func (s *HashSet[T]) Contains(t T) bool {
	return s.m.Contains(t)
}

func (s *HashSet[T]) AllKeys() []T {
	return s.m.AllKeys()
}

func (s *HashSet[T]) Intersect(h Set[T]) Set[T] {
	ret := NewHashSet[T](nil)
	h.Range(func(t T) bool {
		if s.Contains(t) {
			ret.Add(t)
		}
		return true
	})
	return ret
}

func (s *HashSet[T]) Range(fn func(T) bool) {
	s.m.Range(func(t T, s struct{}) bool {
		return fn(t)
	})
}

func (s *HashSet[T]) Size() int {
	return s.m.Size()
}

func (s *HashSet[T]) Clear() {
	s.m.Clear()
}
