package hashmap

import (
	"hash/crc32"
)

type ShardingMap[T any] struct {
	segments []*ConcurrentHashMap[string, T]
}

func NewShardingMap[T any](segmentSize int) *ShardingMap[T] {
	if segmentSize <= 0 {
		segmentSize = 1
	}
	segments := make([]*ConcurrentHashMap[string, T], 0, segmentSize)
	for i := 0; i < segmentSize; i++ {
		segments = append(segments, NewConcurrentHashMap[string, T]())
	}
	return &ShardingMap[T]{
		segments: segments,
	}
}

func (m *ShardingMap[T]) getSegment(k string) *ConcurrentHashMap[string, T] {
	hashRet := int(crc32.ChecksumIEEE([]byte(k)))
	return m.segments[hashRet%len(m.segments)]
}

func (m *ShardingMap[T]) Get(k string) (T, bool) {
	return m.getSegment(k).Get(k)
}

func (m *ShardingMap[T]) Put(k string, v T) {
	m.getSegment(k).Put(k, v)
}

func (m *ShardingMap[T]) Remove(k string) {
	m.getSegment(k).Remove(k)
}

func (m *ShardingMap[T]) GetOrPut(k string, v T) (T, bool) {
	return m.getSegment(k).GetOrPut(k, v)
}

func (m *ShardingMap[T]) AllKeys() []string {
	ret := make([]string, 0)
	for _, seg := range m.segments {
		ret = append(ret, seg.AllKeys()...)
	}
	return ret
}

func (m *ShardingMap[T]) Range(fn func(string, T)) {
	for _, seg := range m.segments {
		seg.Range(fn)
	}
}
