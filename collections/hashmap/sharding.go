package hashmap

import (
	"hash/crc32"
	"sync"
)

type segment[T any] struct {
	sync.RWMutex
	data map[string]T
}

func newSegment[T any]() *segment[T] {
	return &segment[T]{
		data: make(map[string]T, 8),
	}
}

func (s *segment[T]) get(k string) (T, bool) {
	s.RLock()
	defer s.RUnlock()
	r, b := s.data[k]
	return r, b
}

func (s *segment[T]) put(k string, v T) {
	s.Lock()
	defer s.Unlock()
	s.data[k] = v
}

func (s *segment[T]) del(k string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, k)
}

func (s *segment[T]) getOrPut(k string, o T) (T, bool) {
	v, b := s.get(k)
	if b {
		return v, false
	}
	s.Lock()
	defer s.Unlock()
	v, b = s.data[k]
	if b {
		return v, false
	}
	s.data[k] = o
	return o, true
}

func (s *segment[T]) allKeys() []string {
	s.RLock()
	defer s.RUnlock()
	ret := make([]string, 0, len(s.data))
	for k := range s.data {
		ret = append(ret, k)
	}
	return ret
}

func (s *segment[T]) Range(fn func(string, T)) {
	for _, k := range s.allKeys() {
		v, b := s.get(k)
		if b {
			fn(k, v)
		}
	}
}

type ShardingMap[T any] struct {
	segments []*segment[T]
}

func NewShardingMap[T any](segmentSize int) *ShardingMap[T] {
	if segmentSize <= 0 {
		segmentSize = 1
	}
	segments := make([]*segment[T], 0, segmentSize)
	for i := 0; i < segmentSize; i++ {
		segments = append(segments, newSegment[T]())
	}
	return &ShardingMap[T]{
		segments: segments,
	}
}

func (m *ShardingMap[T]) getSegment(k string) *segment[T] {
	hashRet := int(crc32.ChecksumIEEE([]byte(k)))
	return m.segments[hashRet%len(m.segments)]
}

func (m *ShardingMap[T]) Get(k string) (T, bool) {
	return m.getSegment(k).get(k)
}

func (m *ShardingMap[T]) Put(k string, v T) {
	m.getSegment(k).put(k, v)
}

func (m *ShardingMap[T]) Del(k string) {
	m.getSegment(k).del(k)
}

func (m *ShardingMap[T]) GetOrPut(k string, v T) (T, bool) {
	return m.getSegment(k).getOrPut(k, v)
}

func (m *ShardingMap[T]) AllKeys() []string {
	ret := make([]string, 0)
	for _, seg := range m.segments {
		ret = append(ret, seg.allKeys()...)
	}
	return ret
}

func (m *ShardingMap[T]) Range(fn func(string, T)) {
	for _, seg := range m.segments {
		seg.Range(fn)
	}
}
