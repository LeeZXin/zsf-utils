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
	handler                     func(context.Context)
	leaser                      lease.Leaser
	subCancelFn                 atomic.Value
	waitDuration, renewDuration time.Duration
}

func RunMainLoopTask(handler func(context.Context), leaser lease.Leaser, waitDuration, renewDuration time.Duration) (Stopper, error) {
	if handler == nil || leaser == nil || waitDuration <= 0 || renewDuration <= 0 {
		return nil, errors.New("invalid args")
	}
	task := &mainLoopTask{
		handler:       handler,
		leaser:        leaser,
		waitDuration:  waitDuration,
		renewDuration: renewDuration,
	}
	return task.Start(), nil
}

func (t *mainLoopTask) Start() Stopper {
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
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
	})
}

func (t *mainLoopTask) do() {
	ctx, cancelFn := context.WithCancel(context.Background())
	t.subCancelFn.Store(cancelFn)
	defer cancelFn()
	// 尝试加锁
	releaser, renewer, b, err := t.leaser.TryGrant()
	if err == nil && b {
		defer releaser.Release()
		t.handler(ctx)
		// 续期
		go func() {
			time.Sleep(t.renewDuration)
			if ctx.Err() != nil {
				return
			}
			renewRet, err := renewer.Renew(ctx)
			if err != nil || !renewRet {
				return
			}
		}()
	}
}
