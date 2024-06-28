package executor

import "errors"

//拒绝策略
//默认实现两种
// AbortStrategy 直接丢弃
// CallerRunsStrategy 由主协程执行
// StillQueuedStrategy 继续排队

var (
	AbortStrategy RejectStrategy = func(runnable Runnable, executor *Executor) error {
		return errors.New("task rejected by executor")
	}

	CallerRunsStrategy RejectStrategy = func(runnable Runnable, executor *Executor) error {
		runnable.Run()
		return nil
	}

	StillQueuedStrategy RejectStrategy = func(runnable Runnable, executor *Executor) error {
		executor.queue <- runnable
		return nil
	}
)

type RejectStrategy func(runnable Runnable, executor *Executor) error
