package hashset

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
	m map[T]struct{}
}

func NewHashSet[T comparable](arr []T) *HashSet[T] {
	ret := &HashSet[T]{
		m: make(map[T]struct{}, len(arr)),
	}
	for _, t := range arr {
		ret.Add(t)
	}
	return ret
}

func (s *HashSet[T]) Add(ts ...T) {
	for _, t := range ts {
		s.m[t] = struct{}{}
	}
}

func (s *HashSet[T]) Remove(ts ...T) {
	for _, t := range ts {
		delete(s.m, t)
	}
}

func (s *HashSet[T]) Contains(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *HashSet[T]) AllKeys() []T {
	ret := make([]T, 0, len(s.m))
	for k := range s.m {
		ret = append(ret, k)
	}
	return ret
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
	if fn == nil {
		return
	}
	keys := s.AllKeys()
	for _, key := range keys {
		if !fn(key) {
			return
		}
	}
}

func (s *HashSet[T]) Size() int {
	return len(s.m)
}

func (s *HashSet[T]) Clear() {
	for k := range s.m {
		delete(s.m, k)
	}
}
