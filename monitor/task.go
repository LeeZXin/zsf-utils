package monitor

import (
	"github.com/LeeZXin/zsf-utils/taskutil"
	"sync"
	"time"
)

type StatInfoFetchTask struct {
	mu      sync.RWMutex
	data    []StatInfo
	MaxSize int
	task    *taskutil.PeriodicalTask
}

func NewStatInfoFetchTask(tick time.Duration, maxSize int) *StatInfoFetchTask {
	if tick <= 0 {
		tick = time.Second
	}
	if maxSize <= 0 {
		maxSize = 60
	}
	ret := &StatInfoFetchTask{
		mu:      sync.RWMutex{},
		data:    make([]StatInfo, 0, maxSize),
		MaxSize: maxSize,
	}
	task, _ := taskutil.NewPeriodicalTask(tick, ret.fetch)
	task.Start()
	ret.task = task
	return ret
}

func (t *StatInfoFetchTask) fetch() {
	info := GetCurrentStatInfo()
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.data) == t.MaxSize {
		t.data = t.data[1:]
		t.data = append(t.data, info)
	} else {
		t.data = append(t.data, info)
	}
}

func (t *StatInfoFetchTask) GetInfo() []StatInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ret := make([]StatInfo, 0, len(t.data))
	for _, datum := range t.data {
		ret = append(ret, StatInfo{
			MemInfo:    datum.MemInfo,
			NetInfo:    datum.NetInfo,
			CpuPercent: datum.CpuPercent,
		})
	}
	return ret
}

func (t *StatInfoFetchTask) Stop() {
	if t.task != nil {
		t.task.Stop()
	}
}
