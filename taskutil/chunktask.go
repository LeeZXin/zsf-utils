package taskutil

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Chunk[T any] struct {
	Size int
	Data T
}

type FlushFunc[T any] func([]Chunk[T])

type ChunkTaskExecuteFunc[T any] func(data T, dataSize int)

type ChunkTaskFlushFunc func()

type chunkTask[T any] struct {
	mu          sync.Mutex
	triggerSize int
	dataSize    int
	chunkList   []Chunk[T]
	flushFunc   FlushFunc[T]
}

func RunChunkTask[T any](triggerSize int, flushFunc FlushFunc[T], flushInterval time.Duration) (ChunkTaskExecuteFunc[T], ChunkTaskFlushFunc, StopFunc, error) {
	if flushFunc == nil || flushInterval <= 0 || triggerSize <= 0 {
		return nil, nil, nil, errors.New("invalid args")
	}
	task := &chunkTask[T]{
		triggerSize: triggerSize,
		chunkList:   make([]Chunk[T], 0, triggerSize),
		flushFunc:   flushFunc,
	}
	ptStop, _ := RunPeriodicalTask(flushInterval, flushInterval, func(context.Context) {
		task.Flush()
	})
	return task.Execute, task.Flush, NewContextStopper(func() {
		task.Flush()
		ptStop()
	}), nil
}

func (t *chunkTask[T]) Execute(data T, dataSize int) {
	t.addAndFlush(data, dataSize)
}

func (t *chunkTask[T]) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flush()
}

func (t *chunkTask[T]) flush() {
	if len(t.chunkList) > 0 {
		t.flushFunc(t.chunkList)
		t.chunkList = make([]Chunk[T], 0, t.triggerSize)
		t.dataSize = 0
	}
}

func (t *chunkTask[T]) addAndFlush(data T, size int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.chunkList = append(t.chunkList, Chunk[T]{
		Size: size,
		Data: data,
	})
	t.dataSize += size
	if t.dataSize >= t.triggerSize {
		t.flush()
	}
}
