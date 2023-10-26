package selector

// singleNodeSelector 单节点选择器 当节点只有一个时
type singleNodeSelector[T any] struct {
	node Node[T]
}

func newSingleNodeSelector[T any](node Node[T]) Selector[T] {
	return &singleNodeSelector[T]{
		node: node,
	}
}

func (s *singleNodeSelector[T]) Select(...string) (Node[T], error) {
	return s.node, nil
}

func (s *singleNodeSelector[T]) GetNodes() []Node[T] {
	return []Node[T]{s.node}
}
