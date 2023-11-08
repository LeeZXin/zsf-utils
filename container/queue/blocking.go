package queue

import (
	"errors"
	"github.com/LeeZXin/zsf-utils/threadutil"
	"sync"
	"time"
)

var (
	TimeoutError    = errors.New("timeout")
	ChanClosedError = errors.New("chan closed")
)

type ConcurrentBlockingQueue[T any] struct {
	c chan T
	o sync.Once
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	if maxSize <= 0 {
		maxSize = 1
	}
	return &ConcurrentBlockingQueue[T]{
		c: make(chan T, maxSize),
	}
}

func (q *ConcurrentBlockingQueue[T]) Put(t T) {
	_ = threadutil.RunSafe(func() {
		q.c <- t
	})
}

func (q *ConcurrentBlockingQueue[T]) Offer(t T, waitDuration time.Duration) (err error) {
	err2 := threadutil.RunSafe(func() {
		select {
		case <-time.After(waitDuration):
			err = TimeoutError
		case q.c <- t:
		}
	})
	if err2 != nil {
		err = ChanClosedError
	}
	return
}

func (q *ConcurrentBlockingQueue[T]) Take() (T, error) {
	select {
	case t, ok := <-q.c:
		if ok {
			return t, nil
		}
		return t, ChanClosedError
	}
}

func (q *ConcurrentBlockingQueue[T]) Poll(duration time.Duration) (T, error) {
	select {
	case t, ok := <-q.c:
		if ok {
			return t, nil
		}
		return t, ChanClosedError
	case <-time.After(duration):
		var t T
		return t, TimeoutError
	}
}

func (q *ConcurrentBlockingQueue[T]) Close() {
	q.o.Do(func() {
		close(q.c)
	})
}
