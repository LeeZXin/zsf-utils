package selector

import (
	"sync"
)

// WeightedRoundRobinSelector 加权平滑路由选择器
type WeightedRoundRobinSelector[T any] struct {
	nodes       []Node[T]
	selectMutex sync.Mutex
	current     int
	gcd         int
	max         int
}

func (s *WeightedRoundRobinSelector[T]) Select(...string) (Node[T], error) {
	s.selectMutex.Lock()
	defer s.selectMutex.Unlock()
	for {
		s.current = (s.current + 1) % len(s.nodes)
		if s.current == 0 {
			s.max -= s.gcd
			if s.max <= 0 {
				s.max = s.maxWeight()
			}
		}
		if s.nodes[s.current].Weight >= s.max {
			return s.nodes[s.current], nil
		}
	}
}

func (s *WeightedRoundRobinSelector[T]) maxWeight() int {
	m := 0
	for _, server := range s.nodes {
		if server.Weight > m {
			m = server.Weight
		}
	}
	return m
}

func (s *WeightedRoundRobinSelector[T]) init() {
	nodes := s.nodes
	weights := make([]int, len(nodes))
	for i, node := range nodes {
		if node.Weight <= 0 {
			weights[i] = 1
		} else {
			weights[i] = node.Weight
		}
	}
	s.gcd = gcd(weights)
	s.max = max(weights)
}

func (s *WeightedRoundRobinSelector[T]) GetNodes() []Node[T] {
	return s.nodes
}

func NewWeightedRoundRobinSelector[T any](nodes []Node[T]) Selector[T] {
	if nodes == nil || len(nodes) == 0 {
		return &errorSelector[T]{
			Err: EmptyNodesErr,
		}
	}
	if len(nodes) == 1 {
		return newSingleNodeSelector(nodes[0])
	}
	w := &WeightedRoundRobinSelector[T]{nodes: nodes}
	w.init()
	return w
}

func gcd(numbers []int) int {
	result := numbers[0]
	for _, number := range numbers[1:] {
		result = gcdTwoNumbers(result, number)
	}
	return result
}

func gcdTwoNumbers(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func max(numbers []int) int {
	m := numbers[0]
	for _, number := range numbers[1:] {
		if number > m {
			m = number
		}
	}
	return m
}
