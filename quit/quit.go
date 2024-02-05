package quit

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// 监听程序kill事件, 并执行注销函数
// 用于关闭资源等， 如httpServer，数据库等

type ShutdownHook func()

var (
	hookList     = make([]ShutdownHook, 0)
	priorityList = make([]ShutdownHook, 0)
	mu           = sync.Mutex{}
)

// AddShutdownHook 添加关闭钩子 isPriority 放入优先队列
func AddShutdownHook(hook ShutdownHook, isPriority ...bool) {
	if hook != nil {
		mu.Lock()
		defer mu.Unlock()
		if len(isPriority) > 0 && isPriority[0] {
			priorityList = append(priorityList, hook)
		} else {
			hookList = append(hookList, hook)
		}
	}
}

// Wait 注册signal事件，并无限等待
func Wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	mu.Lock()
	defer mu.Unlock()
	// priorityList 优先关闭priorityList
	for _, fn := range priorityList {
		fn()
	}
	for _, fn := range hookList {
		fn()
	}
}
