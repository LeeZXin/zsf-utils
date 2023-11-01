package poolutil

import (
	"github.com/LeeZXin/zsf-utils/threadutil"
	"sync"
)

type waitFuture[T any] struct {
	c chan Object[T]
	o sync.Once
}

func newWaitFuture[T any]() *waitFuture[T] {
	return &waitFuture[T]{
		c: make(chan Object[T], 1),
		o: sync.Once{},
	}
}

func (f *waitFuture[T]) get() (Object[T], error) {
	defer f.close()
	for {
		select {
		case ret, ok := <-f.c:
			if !ok {
				return nil, CancelErr
			}
			return ret, nil
		}
	}
}

func (f *waitFuture[T]) notify(o Object[T]) error {
	return threadutil.RunSafe(func() {
		f.c <- o
	}, f.close)
}

func (f *waitFuture[T]) close() {
	f.o.Do(func() {
		close(f.c)
	})
}
