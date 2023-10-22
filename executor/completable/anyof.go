package completable

import (
	"errors"
	"sync"
	"time"
)

type anyOfFuture struct {
	sync.Mutex
	result   any
	err      error
	done     bool
	doneChan chan struct{}
	baseChan chan iBase
	doneOnce sync.Once
	stack    []iBase
	waits    []iBase
}

func newAnyOfFuture(bases ...iBase) *anyOfFuture {
	return &anyOfFuture{
		Mutex:    sync.Mutex{},
		doneChan: make(chan struct{}),
		baseChan: make(chan iBase, len(bases)),
		doneOnce: sync.Once{},
		stack:    make([]iBase, 0),
		waits:    bases,
	}
}

func (f *anyOfFuture) checkCompletedAndAppendToStack(i iBase) bool {
	if i == nil {
		return false
	}
	f.Lock()
	defer f.Unlock()
	if f.done {
		return false
	}
	f.stack = append(f.stack, i)
	return true
}

func (f *anyOfFuture) notify(base iBase) {
	if base == nil {
		return
	}
	f.Lock()
	defer f.Unlock()
	if !f.done {
		f.baseChan <- base
	}
}

func (f *anyOfFuture) Get() (any, error) {
	return f.GetWithTimeout(0)
}

func (f *anyOfFuture) getResult() (any, error) {
	return f.Get()
}

func (f *anyOfFuture) getResultAndErr() (any, error) {
	f.Lock()
	defer f.Unlock()
	return f.result, f.err
}

func (f *anyOfFuture) setResultAndErr(result any, err error) {
	f.Lock()
	defer f.Unlock()
	f.result = result
	f.err = err
	f.done = true
}

func (f *anyOfFuture) GetWithTimeout(timeout time.Duration) (any, error) {
	if f.isDone() {
		return f.getResultAndErr()
	}
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case <-timer.C:
			return nil, timeoutError
		case <-f.doneChan:
			return f.getResultAndErr()
		}
	} else {
		select {
		case <-f.doneChan:
			return f.getResultAndErr()
		}
	}
}

func (f *anyOfFuture) fire() {
	defer f.postComplete()
	if len(f.waits) == 0 {
		return
	}
	for {
		select {
		case b, ok := <-f.baseChan:
			if !ok {
				return
			}
			result, err := b.getResult()
			f.setResultAndErr(result, err)
			return
		case <-f.doneChan:
			return
		}
	}
}

func (f *anyOfFuture) isDone() bool {
	f.Lock()
	defer f.Unlock()
	return f.done
}

func (f *anyOfFuture) postComplete() {
	f.doneOnce.Do(func() {
		close(f.doneChan)
		close(f.baseChan)
		f.Lock()
		stack := f.stack[:]
		f.Unlock()
		for _, i := range stack {
			i.fire()
		}
	})
}

func ThenAnyOf(bases ...iBase) IFuture[any] {
	return thenAnyOf(false, bases...)
}

func ThenAnyOfAsync(bases ...iBase) IFuture[any] {
	return thenAnyOf(true, bases...)
}

func thenAnyOf(isAsync bool, bases ...iBase) IFuture[any] {
	if len(bases) == 0 {
		return newKnownResultFutureWithErr[any](errors.New("nil futures"))
	}
	f := newAnyOfFuture(bases...)
	for _, i := range bases {
		if !i.checkCompletedAndAppendToStack(&relayBase{
			f:      f,
			parent: i,
		}) {
			f.notify(i)
		}
	}
	if isAsync {
		go f.fire()
	} else {
		f.fire()
	}
	return f
}
