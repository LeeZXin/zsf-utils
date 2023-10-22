package completable

import (
	"errors"
)

type iNotify interface {
	notify(base iBase)
}

type relayBase struct {
	f      iNotify
	parent iBase
}

func (r *relayBase) fire() {
	r.f.notify(r.parent)
}

func (r *relayBase) checkCompletedAndAppendToStack(iBase) bool {
	return false
}

func (r *relayBase) getResult() (any, error) {
	return nil, nil
}

type allOfFuture struct {
	*anyOfFuture
}

func (f *allOfFuture) fire() {
	defer f.postComplete()
	if len(f.waits) == 0 {
		return
	}
	count := 0
	for {
		select {
		case b, ok := <-f.baseChan:
			if !ok {
				return
			}
			_, err := b.getResult()
			if err != nil {
				f.setResultAndErr(nil, err)
				return
			}
			count += 1
			if count == len(f.waits) {
				return
			}
		case <-f.doneChan:
			return
		}
	}
}

func ThenAllOf(bases ...iBase) IFuture[any] {
	return thenAllOf(false, bases...)
}

func ThenAllOfAsync(bases ...iBase) IFuture[any] {
	return thenAllOf(true, bases...)
}

func thenAllOf(isAsync bool, bases ...iBase) IFuture[any] {
	if len(bases) == 0 {
		return newKnownResultFutureWithErr[any](errors.New("nil futures"))
	}
	f := &allOfFuture{
		anyOfFuture: newAnyOfFuture(bases...),
	}
	for _, i := range bases {
		if !i.checkCompletedAndAppendToStack(&relayBase{
			f:      f,
			parent: i,
		}) {
			_, err := i.getResult()
			if err != nil {
				return newKnownResultFutureWithErr[any](err)
			}
		}
	}
	if isAsync {
		go f.fire()
	} else {
		f.fire()
	}
	return f
}
