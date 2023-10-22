package completable

import "errors"

func ThenApply[T, K any](f IFuture[T], fn ApplyFunc[T, K]) IFuture[K] {
	return thenApply(f, fn, false)
}

func ThenApplyAsync[T, K any](f IFuture[T], fn ApplyFunc[T, K]) IFuture[K] {
	return thenApply(f, fn, true)
}

func thenApply[T, K any](f IFuture[T], fn ApplyFunc[T, K], isAsync bool) IFuture[K] {
	if f == nil {
		return newKnownResultFutureWithErr[K](errors.New("nil futures"))
	}
	if fn == nil {
		return newKnownResultFutureWithErr[K](errors.New("nil apply fn"))
	}
	cf := newCallFuture[K](func() (K, error) {
		fret, err := f.Get()
		if err != nil {
			var k K
			return k, err
		}
		return fn(fret)
	}, isAsync)
	if f.checkCompletedAndAppendToStack(cf) {
		cf.fire()
		return cf
	}
	fret, err := f.Get()
	if err != nil {
		return newKnownResultFutureWithErr[K](err)
	}
	return call(func() (K, error) {
		return fn(fret)
	}, isAsync)
}
