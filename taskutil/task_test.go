package taskutil

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRunPeriodicalTask(t *testing.T) {
	fmt.Println(time.Now(), "start")
	stopFunc, _ := RunPeriodicalTask(time.Second, 3*time.Second, func(ctx context.Context) {
		fmt.Println(time.Now(), "run")
	})
	time.Sleep(10 * time.Second)
	stopFunc()
	fmt.Println(time.Now(), "stop")
}

func TestRunChunkTask(t *testing.T) {
	executeFunc, _, stopFunc, _ := RunChunkTask[string](3, func(chunks []Chunk[string]) {
		fmt.Println(time.Now(), chunks)
	}, 10*time.Second)
	defer stopFunc()
	for i := 0; i < 10; i++ {
		executeFunc(strconv.Itoa(i), 1)
	}
	time.Sleep(20 * time.Second)
}
