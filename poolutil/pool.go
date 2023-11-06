package poolutil

import (
	"errors"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"sync"
	"time"
)

var (
	CancelErr        = errors.New("wait for object canceled")
	PoolExhaustedErr = errors.New("pool exhausted")
	PoolClosedErr    = errors.New("pool closed")
)

type PoolConfig[T any] struct {
	MinIdle            int
	MaxIdle            int
	MaxActive          int
	MaxWait            int
	IdleDuration       time.Duration
	CleanupDuration    time.Duration
	Factory            ObjectFactory[T]
	BlockWhenExhausted bool
}

type Pool[T any] interface {
	BorrowObject() (Object[T], error)
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

type objectWrapper[T any] struct {
	Object[T]
	t time.Time
}

func newObjectWrapper[T any](object Object[T]) *objectWrapper[T] {
	return &objectWrapper[T]{
		Object: object,
		t:      time.Now(),
	}
}

func (w *objectWrapper[T]) IsNotExpired(idleDuration time.Duration) bool {
	if idleDuration <= 0 {
		return true
	}
	return w.t.Add(idleDuration).After(time.Now())
}

type ObjectFactory[T any] interface {
	CreateObject() (Object[T], error)
}

type GenericPool[T any] struct {
	config      PoolConfig[T]
	activeNum   int
	pool        []*objectWrapper[T]
	waitList    []*waitFuture[T]
	objectMu    sync.Mutex
	closed      bool
	closeOnce   sync.Once
	cleanupTask *taskutil.PeriodicalTask
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
	if config.CleanupDuration <= 0 {
		config.CleanupDuration = time.Minute
	}
	pool := make([]*objectWrapper[T], 0, config.MaxIdle)
	for i := 0; i < config.MinIdle; i++ {
		object, err := config.Factory.CreateObject()
		if err == nil {
			pool = append(pool, newObjectWrapper(object))
		}
	}
	p := &GenericPool[T]{
		config: config,
		pool:   pool,
	}
	p.cleanupTask, _ = taskutil.NewPeriodicalTask(config.CleanupDuration, p.cleanup)
	p.cleanupTask.Start()
	return p, nil
}

func (p *GenericPool[T]) cleanup() {
	if p.objectMu.TryLock() {
		shouldClose := make([]Object[T], 0)
		newPool := make([]*objectWrapper[T], 0, p.config.MaxIdle)
		for i := range p.pool {
			o := p.pool[i]
			if o.IsNotExpired(p.config.IdleDuration) && o.IsActive() {
				newPool = append(newPool, o)
			} else {
				shouldClose = append(shouldClose, o)
			}
		}
		if len(newPool) != len(p.pool) {
			p.pool = newPool
		}
		p.objectMu.Unlock()
		for _, o := range shouldClose {
			o.Close()
		}
	}
}

func (p *GenericPool[T]) BorrowObject() (Object[T], error) {
	shouldClose := make([]Object[T], 0)
	defer func() {
		for _, o := range shouldClose {
			o.Close()
		}
	}()
	p.objectMu.Lock()
	if p.closed {
		p.objectMu.Unlock()
		return nil, PoolClosedErr
	}
	for len(p.pool) > 0 {
		obj := p.pool[0]
		p.pool = p.pool[1:]
		if obj.IsNotExpired(p.config.IdleDuration) && obj.IsActive() {
			p.activeNum += 1
			p.objectMu.Unlock()
			return obj.Object, nil
		} else {
			shouldClose = append(shouldClose, obj)
		}
	}
	if p.config.MaxActive <= 0 || p.config.MaxActive > p.activeNum {
		object, err := p.config.Factory.CreateObject()
		if err == nil {
			p.activeNum += 1
		}
		p.objectMu.Unlock()
		return object, err
	}
	if p.config.BlockWhenExhausted {
		if p.config.MaxWait < 0 || p.config.MaxWait > len(p.waitList) {
			future := newWaitFuture[T]()
			p.waitList = append(p.waitList, future)
			p.objectMu.Unlock()
			return future.get()
		}
	}
	p.objectMu.Unlock()
	return nil, PoolExhaustedErr
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
		for len(p.waitList) > 0 {
			f := p.waitList[0]
			p.waitList = p.waitList[1:]
			if err := f.notify(t); err == nil {
				return
			}
		}
		if p.activeNum > 0 {
			p.activeNum -= 1
		}
		if len(p.pool) < p.config.MaxIdle {
			p.pool = append(p.pool, newObjectWrapper(t))
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
		p.cleanupTask.Stop()
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
	return p.config.MaxIdle
}

func (p *GenericPool[T]) GetMinIdle() int {
	return p.config.MinIdle
}

func (p *GenericPool[T]) GetIdleNum() int {
	p.objectMu.Lock()
	defer p.objectMu.Unlock()
	return len(p.pool)
}
