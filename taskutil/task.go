package taskutil

import (
	"context"
	"sync"
)

type Stopper interface {
	Stop()
}

type contextStopper struct {
	stopOnce sync.Once
	cancelFn context.CancelFunc
}

func (s *contextStopper) Stop() {
	s.stopOnce.Do(func() {
		if s.cancelFn != nil {
			s.cancelFn()
		}
	})
}

func NewContextStopper(cancelFunc context.CancelFunc) Stopper {
	return &contextStopper{
		cancelFn: cancelFunc,
	}
}
