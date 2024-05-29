package taskutil

import (
	"context"
	"sync"
)

type StopFunc func()

type contextStopper struct {
	stopOnce sync.Once
	cancelFn context.CancelFunc
}

func (s *contextStopper) stop() {
	s.stopOnce.Do(func() {
		if s.cancelFn != nil {
			s.cancelFn()
		}
	})
}

func NewContextStopper(cancelFunc context.CancelFunc) StopFunc {
	stopper := &contextStopper{
		cancelFn: cancelFunc,
	}
	return stopper.stop
}
