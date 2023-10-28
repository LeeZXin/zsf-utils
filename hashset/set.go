package hashset

type Set[T comparable] interface {
	Add(T)
	Delete(...T)
	Contains(T) bool
	Copy() Set[T]
	AllKeys() []T
	Intersect(h Set[T]) Set[T]
	Range(func(T) bool)
}

type HashSet[T comparable] struct {
	m map[T]struct{}
}

func NewHashSet[T comparable](arr []T) Set[T] {
	ret := &HashSet[T]{
		m: make(map[T]struct{}, len(arr)),
	}
	for _, t := range arr {
		ret.Add(t)
	}
	return ret
}

func (s *HashSet[T]) Add(key T) {
	s.m[key] = struct{}{}
}

func (s *HashSet[T]) Delete(keys ...T) {
	for _, key := range keys {
		delete(s.m, key)
	}
}

func (s *HashSet[T]) Contains(key T) bool {
	_, ok := s.m[key]
	return ok
}

func (s *HashSet[T]) Copy() Set[T] {
	return NewHashSet(s.AllKeys())
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
