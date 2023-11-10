package taskutil

import (
	"context"
	"errors"
	"github.com/LeeZXin/zsf-utils/threadutil"
	"sync"
	"time"
)

type PeriodicalTask struct {
	interval time.Duration
	fn       func()

	ctx      context.Context
	cancelFn context.CancelFunc

	startOnce sync.Once
	stopOnce  sync.Once
}

func NewPeriodicalTask(interval time.Duration, fn func()) (*PeriodicalTask, error) {
	if interval == 0 || fn == nil {
		return nil, errors.New("invalid task arguments")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &PeriodicalTask{
		interval:  interval,
		fn:        fn,
		ctx:       ctx,
		cancelFn:  cancelFunc,
		startOnce: sync.Once{},
		stopOnce:  sync.Once{},
	}, nil
}

func (t *PeriodicalTask) Start() {
	t.startOnce.Do(func() {
		go func() {
			for {
				if t.ctx.Err() != nil {
					return
				}
				t.Execute()
				time.Sleep(t.interval)
			}
		}()
	})
}

func (t *PeriodicalTask) Stop() {
	t.stopOnce.Do(func() {
		t.cancelFn()
	})
}

func (t *PeriodicalTask) Execute() {
	_ = threadutil.RunSafe(t.fn)
}
