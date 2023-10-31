package poolutil

import (
	"github.com/LeeZXin/zsf-utils/threadutil"
	"sync"
	"time"
)

type waitFuture[T any] struct {
	c chan Object[T]
	d time.Duration
	o sync.Once
}

func newWaitFuture[T any](d time.Duration) *waitFuture[T] {
	return &waitFuture[T]{
		c: make(chan Object[T], 1),
		d: d,
		o: sync.Once{},
	}
}

func (f *waitFuture[T]) get() (Object[T], error) {
	if f.d > 0 {
		t := time.NewTimer(f.d)
		defer func() {
			t.Stop()
			f.close()
		}()
		for {
			select {
			case <-t.C:
				return nil, TimeoutErr
			case ret, ok := <-f.c:
				if !ok {
					return nil, CancelErr
				}
				return ret, nil
			}
		}
	} else {
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
}

func (f *waitFuture[T]) notify(o Object[T]) {
	if o == nil {
		return
	}
	threadutil.RunSafe(func() {
		f.c <- o
	}, f.close)
}

func (f *waitFuture[T]) close() {
	f.o.Do(func() {
		close(f.c)
	})
}
