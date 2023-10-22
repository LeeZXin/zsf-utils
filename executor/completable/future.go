package completable

import (
	"errors"
	"sync"
	"time"
)

var (
	timeoutError = errors.New("task timeout")
)

type Callable[T any] func() (T, error)
type ApplyFunc[T, K any] func(T) (K, error)
type CombineFunc[T, K, L any] func(T, K) (L, error)

type iBase interface {
	fire()
	getResult() (any, error)
	checkCompletedAndAppendToStack(i iBase) bool
}

type IFuture[T any] interface {
	iBase
	Get() (T, error)
	GetWithTimeout(timeout time.Duration) (T, error)
}

// futureResult promise err
type futureResult[T any] struct {
	Result T
	Err    error
}

type callFuture[T any] struct {
	sync.Mutex
	result   *futureResult[T]
	done     chan struct{}
	doneOnce sync.Once
	callable Callable[T]

	isAsync bool
	stack   []iBase
}

func newCallFuture[T any](callable Callable[T], isAsync bool) *callFuture[T] {
	return &callFuture[T]{
		Mutex:    sync.Mutex{},
		done:     make(chan struct{}),
		doneOnce: sync.Once{},
		callable: callable,
		isAsync:  isAsync,
		stack:    make([]iBase, 0),
	}
}

func (b *callFuture[T]) postComplete() {
	b.Lock()
	result := b.result
	stack := b.stack[:]
	b.Unlock()
	if result == nil {
		return
	}
	for _, i := range stack {
		i.fire()
	}
}

func (b *callFuture[T]) fire() {
	if b.isAsync {
		go func() {
			res, err := b.callable()
			b.setResult(&futureResult[T]{
				Result: res,
				Err:    err,
			})
			b.completed()
		}()
	} else {
		res, err := b.callable()
		b.setResult(&futureResult[T]{
			Result: res,
			Err:    err,
		})
		b.completed()
	}
}

func (b *callFuture[T]) getResult() (any, error) {
	return b.Get()
}

func (b *callFuture[T]) setResult(result *futureResult[T]) bool {
	b.Lock()
	defer b.Unlock()
	if b.result != nil {
		return false
	}
	b.result = result
	return true
}

func (b *callFuture[T]) completed() {
	b.doneOnce.Do(func() {
		close(b.done)
		b.postComplete()
	})
}

func (b *callFuture[T]) getFutureResult() *futureResult[T] {
	b.Lock()
	defer b.Unlock()
	return b.result
}

func (b *callFuture[T]) checkCompletedAndAppendToStack(i iBase) bool {
	if i == nil {
		return false
	}
	b.Lock()
	defer b.Unlock()
	if b.result != nil {
		return false
	}
	b.stack = append(b.stack, i)
	return true
}

func (b *callFuture[T]) IsCompleted() bool {
	b.Lock()
	defer b.Unlock()
	return b.result != nil
}

// Get 阻塞获取结果 无限期等待
func (b *callFuture[T]) Get() (T, error) {
	return b.GetWithTimeout(0)
}

// GetWithTimeout 带超时返回结果 超时返回timeoutErr
func (b *callFuture[T]) GetWithTimeout(timeout time.Duration) (T, error) {
	val := b.getFutureResult()
	if val == nil {
		if timeout > 0 {
			timer := time.NewTimer(timeout)
			defer timer.Stop()
			select {
			case <-b.done:
				break
			case <-timer.C:
				var t T
				return t, timeoutError
			}
		} else {
			select {
			case <-b.done:
				break
			}
		}
		val = b.getFutureResult()
	}
	return val.Result, val.Err
}

type knownResultFuture[T any] struct {
	Result T
	Err    error
}

func newKnownResultFuture[T any](result T) IFuture[T] {
	return &knownResultFuture[T]{
		Result: result,
	}
}

func newKnownResultFutureWithErr[T any](err error) IFuture[T] {
	return &knownResultFuture[T]{
		Err: err,
	}
}

func (f *knownResultFuture[T]) fire() {}

func (f *knownResultFuture[T]) getResult() (any, error) {
	return f.Result, f.Err
}

func (f *knownResultFuture[T]) Get() (T, error) {
	return f.Result, f.Err
}

func (f *knownResultFuture[T]) GetWithTimeout(time.Duration) (T, error) {
	return f.Result, f.Err
}

func (f *knownResultFuture[T]) checkCompletedAndAppendToStack(i iBase) bool {
	return false
}
