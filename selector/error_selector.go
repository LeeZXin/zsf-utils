package selector

// errorSelector 配置错误选择器
type errorSelector[T any] struct {
	Err error
}

func newErrorSelector[T any](err error) Selector[T] {
	return &errorSelector[T]{
		Err: err,
	}
}

func (e *errorSelector[T]) Select(...string) (node Node[T], err error) {
	err = e.Err
	return
}

func (e *errorSelector[T]) GetNodes() []Node[T] {
	return nil
}

func (e *errorSelector[T]) Size() int {
	return 0
}
