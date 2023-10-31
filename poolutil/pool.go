package poolutil

import (
	"errors"
	"sync"
	"time"
)

var (
	CancelErr        = errors.New("wait for object canceled")
	TimeoutErr       = errors.New("wait timeout")
	PoolExhaustedErr = errors.New("pool exhausted")
	PoolClosedErr    = errors.New("pool closed")
)

type PoolConfig[T any] struct {
	MinIdle   int
	MaxIdle   int
	MaxActive int
	Factory   ObjectFactory[T]
}

type Pool[T any] interface {
	BorrowObject() (Object[T], error)
	BorrowObjectUntil(time.Duration) (Object[T], error)
	ReturnObject(Object[T])
	Close()
	IsClosed() bool
	GetMaxIdle() int
	GetMinIdle() int
	GetIdleNum() int
}

type Object[T any] interface {
	IsActive() bool
	GetObject() T
	Close()
}

type ObjectFactory[T any] interface {
	CreateObject() (Object[T], error)
}

type GenericPool[T any] struct {
	minIdle   int
	maxIdle   int
	maxActive int
	activeNum int
	factory   ObjectFactory[T]
	pool      []Object[T]
	waitList  []*waitFuture[T]
	objectMu  sync.Mutex
	closed    bool
	closeOnce sync.Once
}

func NewGenericPool[T any](config PoolConfig[T]) (*GenericPool[T], error) {
	if config.MinIdle < 0 {
		config.MinIdle = 0
	}
	if config.MaxIdle < 0 {
		config.MaxIdle = 0
	}
	if config.MaxActive < 0 {
		config.MaxActive = 0
	}
	if config.MaxIdle < config.MinIdle {
		return nil, errors.New("maxIdle should be greater than minIdle")
	}
	if config.MaxActive < config.MaxIdle {
		return nil, errors.New("maxIdle should be greater than minIdle")
	}
	if config.Factory == nil {
		return nil, errors.New("nil object factory")
	}
	pool := make([]Object[T], 0, config.MaxIdle)
	for i := 0; i < config.MinIdle; i++ {
		object, err := config.Factory.CreateObject()
		if err == nil {
			pool = append(pool, object)
		}
	}
	return &GenericPool[T]{
		minIdle:   config.MinIdle,
		maxIdle:   config.MaxIdle,
		maxActive: config.MaxActive,
		factory:   config.Factory,
		pool:      pool,
	}, nil
}

func (p *GenericPool[T]) borrowObject(shouldWait bool, d time.Duration) (Object[T], error) {
	p.objectMu.Lock()
	if p.closed {
		p.objectMu.Unlock()
		return nil, PoolClosedErr
	}
	for len(p.pool) > 0 {
		object := p.pool[0]
		p.pool = p.pool[1:]
		if object.IsActive() {
			p.activeNum += 1
			p.objectMu.Unlock()
			return object, nil
		} else {
			object.Close()
		}
	}
	if p.maxActive <= 0 || p.maxActive > p.activeNum {
		object, err := p.factory.CreateObject()
		if err == nil {
			p.activeNum += 1
		}
		p.objectMu.Unlock()
		return object, err
	}
	if !shouldWait {
		p.objectMu.Unlock()
		return nil, PoolExhaustedErr
	}
	future := newWaitFuture[T](d)
	p.waitList = append(p.waitList, future)
	p.objectMu.Unlock()
	return future.get()
}

func (p *GenericPool[T]) BorrowObject() (Object[T], error) {
	return p.borrowObject(false, 0)
}

func (p *GenericPool[T]) BorrowObjectUntil(d time.Duration) (Object[T], error) {
	return p.borrowObject(true, d)
}

func (p *GenericPool[T]) ReturnObject(t Object[T]) {
	if t == nil {
		return
	}
	p.objectMu.Lock()
	defer p.objectMu.Unlock()
	if p.closed {
		t.Close()
		return
	}
	if t.IsActive() {
		if len(p.waitList) > 0 {
			f := p.waitList[0]
			p.waitList = p.waitList[1:]
			f.notify(t)
			return
		}
		if p.activeNum > 0 {
			p.activeNum -= 1
		}
		if len(p.pool) < p.maxIdle {
			p.pool = append(p.pool, t)
		} else {
			t.Close()
		}
	} else {
		if p.activeNum > 0 {
			p.activeNum -= 1
		}
		t.Close()
	}
}

func (p *GenericPool[T]) Close() {
	p.closeOnce.Do(func() {
		p.objectMu.Lock()
		p.closed = true
		pool := p.pool[:]
		waitList := p.waitList[:]
		p.pool = nil
		p.waitList = nil
		p.objectMu.Unlock()
		for _, o := range pool {
			o.Close()
		}
		for _, w := range waitList {
			w.close()
		}
	})
}

func (p *GenericPool[T]) IsClosed() bool {
	p.objectMu.Lock()
	defer p.objectMu.Unlock()
	return p.closed
}

func (p *GenericPool[T]) GetMaxIdle() int {
	return p.maxIdle
}

func (p *GenericPool[T]) GetMinIdle() int {
	return p.minIdle
}

func (p *GenericPool[T]) GetIdleNum() int {
	p.objectMu.Lock()
	defer p.objectMu.Unlock()
	return len(p.pool)
}
