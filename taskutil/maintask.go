package taskutil

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/lease"
	"sync/atomic"
	"time"
)

/*
mainLoopTask 集群节点中，只有单个节点执行，其他节点会等待执行
*/
type mainLoopTask struct {
	handler                      func(context.Context)
	leaser                       lease.Leaser
	subCancelFn, releaser        atomic.Value
	waitDuration, renewDuration  time.Duration
	grantCallback, renewCallback func(error, bool)
	releaseCallback              func()
}

type MainLoopTaskOpts struct {
	Handler         func(context.Context)
	Leaser          lease.Leaser
	WaitDuration    time.Duration
	RenewDuration   time.Duration
	GrantCallback   func(error, bool)
	RenewCallback   func(error, bool)
	ReleaseCallback func()
}

func (m *MainLoopTaskOpts) IsValid() bool {
	if m.Handler == nil || m.Leaser == nil || m.WaitDuration <= 0 || m.RenewDuration <= 0 {
		return false
	}
	return true
}

func RunMainLoopTask(opts MainLoopTaskOpts) (StopFunc, error) {
	if !opts.IsValid() {
		return nil, errors.New("invalid args")
	}
	task := &mainLoopTask{
		handler:         opts.Handler,
		leaser:          opts.Leaser,
		waitDuration:    opts.WaitDuration,
		renewDuration:   opts.RenewDuration,
		grantCallback:   opts.GrantCallback,
		renewCallback:   opts.RenewCallback,
		releaseCallback: opts.ReleaseCallback,
	}
	return task.Start(), nil
}

func (t *mainLoopTask) Start() StopFunc {
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		for ctx.Err() == nil {
			t.do()
			time.Sleep(t.waitDuration)
		}
	}()
	return NewContextStopper(func() {
		cancelFn()
		fn := t.subCancelFn.Load()
		if fn != nil {
			fn.(context.CancelFunc)()
		}
		r := t.releaser.Load()
		if r != nil {
			r.(lease.Releaser).Release()
			if t.releaseCallback != nil {
				t.releaseCallback()
			}
		}
	})
}

func (t *mainLoopTask) do() {
	ctx, cancelFn := context.WithCancel(context.Background())
	t.subCancelFn.Store(cancelFn)
	defer cancelFn()
	// 尝试加锁
	releaser, renewer, b, err := t.leaser.TryGrant()
	if t.grantCallback != nil {
		t.grantCallback(err, b)
	}
	if err == nil && b {
		t.releaser.Store(releaser)
		defer func() {
			releaser.Release()
			if t.releaseCallback != nil {
				t.releaseCallback()
			}
		}()
		// 续期
		go func() {
			for {
				time.Sleep(t.renewDuration)
				if ctx.Err() != nil {
					return
				}
				renewRet, err := renewer.Renew(ctx)
				if t.renewCallback != nil {
					t.renewCallback(err, renewRet)
				}
				if err != nil || !renewRet {
					cancelFn()
					return
				}
			}
		}()
		// 执行循环函数
		t.handler(ctx)
	}
}
