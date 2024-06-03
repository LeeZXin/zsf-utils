package taskutil

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf-utils/lease"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"testing"
	"time"
	"xorm.io/xorm"
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

func TestRunMainLoopTask(t *testing.T) {
	engine, _ := xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/hhhh?charset=utf8")
	leaser, _ := lease.NewDbLease("timer-lock", idutil.RandomUuid(), "ztimer_lock", engine, 10*time.Second)
	RunMainLoopTask(MainLoopTaskOpts{
		Handler: func(ctx context.Context) {
			for ctx.Err() == nil {
				fmt.Println(time.Now(), "fufufufu")
				time.Sleep(5 * time.Second)
			}
		},
		Leaser:        leaser,
		WaitDuration:  10 * time.Second,
		RenewDuration: 5 * time.Second,
		GrantCallback: func(err error, b bool) {
			fmt.Println(time.Now(), "grant", err, b)
		},
		RenewCallback: func(err error, b bool) {
			fmt.Println(time.Now(), "renew", err, b)
		},
		ReleaseCallback: func() {
			fmt.Println(time.Now(), "release")
		},
	})
	select {}
}
