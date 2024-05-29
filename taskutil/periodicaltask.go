package taskutil

import (
	"context"
	"errors"
	"time"
)

func RunPeriodicalTask(delay, interval time.Duration, handler func(context.Context)) (StopFunc, error) {
	if interval <= 0 || handler == nil {
		return nil, errors.New("invalid args")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		for {
			if ctx.Err() != nil {
				return
			}
			handler(ctx)
			time.Sleep(interval)
		}
	}()
	return NewContextStopper(cancelFunc), nil
}
