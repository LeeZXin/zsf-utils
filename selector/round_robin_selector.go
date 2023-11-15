package selector

import (
	"github.com/LeeZXin/zsf-utils/randutil"
	"sync/atomic"
)

// RoundRobinSelector 轮询路由选择器
type RoundRobinSelector[T any] struct {
	nodes []Node[T]
	index uint64
}

func (s *RoundRobinSelector[T]) Select(...string) (Node[T], error) {
	index := atomic.AddUint64(&s.index, 1)
	return s.nodes[index%uint64(len(s.nodes))], nil
}

func (s *RoundRobinSelector[T]) GetNodes() []Node[T] {
	return s.nodes
}

func (s *RoundRobinSelector[T]) Size() int {
	return len(s.nodes)
}

func NewRoundRobinSelector[T any](nodes []Node[T]) Selector[T] {
	if len(nodes) == 0 {
		return &errorSelector[T]{
			Err: EmptyNodesErr,
		}
	}
	if len(nodes) == 1 {
		return newSingleNodeSelector(nodes[0])
	}
	r := &RoundRobinSelector[T]{nodes: nodes}
	r.index = uint64(randutil.Intn(len(r.nodes)))
	return r
}
