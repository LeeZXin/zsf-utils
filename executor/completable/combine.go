package completable

import "errors"

func ThenCombine[T, K, L any](t IFuture[T], k IFuture[K], fn CombineFunc[T, K, L]) IFuture[L] {
	return thenCombine(t, k, fn, false)
}

func ThenCombineAsync[T, K, L any](t IFuture[T], k IFuture[K], fn CombineFunc[T, K, L]) IFuture[L] {
	return thenCombine(t, k, fn, true)
}

func thenCombine[T, K, L any](t IFuture[T], k IFuture[K], fn CombineFunc[T, K, L], isAsync bool) IFuture[L] {
	if t == nil || k == nil {
		return newKnownResultFutureWithErr[L](errors.New("nil futures"))
	}
	if fn == nil {
		return newKnownResultFutureWithErr[L](errors.New("nil fn"))
	}
	all := ThenAllOf(t, k)
	return thenApply(all, func(any) (L, error) {
		tret, _ := t.Get()
		kret, _ := k.Get()
		return fn(tret, kret)
	}, isAsync)
}
