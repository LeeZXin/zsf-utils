package completable

import "errors"

func Call[T any](callable Callable[T]) IFuture[T] {
	return call(callable, false)
}

func CallAsync[T any](callable Callable[T]) IFuture[T] {
	return call(callable, true)
}

func call[T any](callable Callable[T], isAsync bool) IFuture[T] {
	if callable == nil {
		return newKnownResultFutureWithErr[T](errors.New("nil futures"))
	}
	f := newCallFuture[T](callable, isAsync)
	f.fire()
	return f
}
