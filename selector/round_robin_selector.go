package selector

import (
	"math/rand"
	"sync/atomic"
)

// RoundRobinSelector 轮询路由选择器
type RoundRobinSelector[T any] struct {
	Nodes []Node[T]
	index uint64
}

func (s *RoundRobinSelector[T]) Select(...string) (Node[T], error) {
	index := atomic.AddUint64(&s.index, 1)
	return s.Nodes[index%uint64(len(s.Nodes))], nil
}

func NewRoundRobinSelector[T any](nodes []Node[T]) (Selector[T], error) {
	if nodes == nil || len(nodes) == 0 {
		return nil, EmptyNodesErr
	} else if len(nodes) == 1 {
		return &SingleNodeSelector[T]{Node: nodes[0]}, nil
	}
	r := &RoundRobinSelector[T]{Nodes: nodes}
	r.index = uint64(rand.Intn(len(r.Nodes)))
	return r, nil
}
