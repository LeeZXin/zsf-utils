package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"os"
	"runtime"
	"time"
)

var (
	sys = getSysInfo()
)

type StatInfo struct {
	MemInfo    MemInfo
	NetInfo    net.IOCountersStat
	CpuPercent float64
	Time       time.Time
}

type MemInfo struct {
	MemAll  uint64
	MemFree uint64
	MemUsed uint64

	MemUsedPercent float64
}

func GetCurrentMemInfo() MemInfo {
	unit := uint64(1024 * 1024)
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemInfo{}
	}
	info := MemInfo{}
	info.MemAll = v.Total
	info.MemFree = v.Free
	info.MemUsed = info.MemAll - info.MemFree
	info.MemUsedPercent = v.UsedPercent
	info.MemAll /= unit
	info.MemUsed /= unit
	info.MemFree /= unit
	return info
}

func GetCurrentNetworkInfo() net.IOCountersStat {
	ret, err := net.IOCounters(true)
	if err != nil {
		return net.IOCountersStat{}
	}
	return ret[0]
}

func GetCurrentCpuPercent() float64 {
	cc, err := cpu.Percent(0, false)
	if err != nil {
		return 0
	}
	return cc[0]
}

type SysInfo struct {
	Platform string
	Os       string
	CpuNum   int
	Arch     string
	HostName string
	WorkDir  string
	BootTime time.Time
	Pid      int

	PhysicalCpuCnt int
	LogicalCpuCnt  int
}

func GetCurrentStatInfo() StatInfo {
	return StatInfo{
		MemInfo:    GetCurrentMemInfo(),
		NetInfo:    GetCurrentNetworkInfo(),
		CpuPercent: GetCurrentCpuPercent(),
		Time:       time.Now(),
	}
}

func getSysInfo() SysInfo {
	info, err := host.Info()
	if err != nil {
		return SysInfo{}
	}
	ret := SysInfo{
		Platform: fmt.Sprintf("%s(%s) %s", info.Platform, info.PlatformFamily, info.PlatformVersion),
		HostName: info.Hostname,
		Os:       info.OS,
		CpuNum:   runtime.NumCPU(),
		Arch:     info.KernelArch,
		BootTime: time.Now(),
		Pid:      os.Getpid(),
	}
	dir, err := os.Getwd()
	if err == nil {
		ret.WorkDir = dir
	}
	name, err := os.Hostname()
	if err == nil {
		ret.HostName = name
	}
	cnt, err := cpu.Counts(false)
	if err == nil {
		ret.PhysicalCpuCnt = cnt
	}
	cnt, err = cpu.Counts(true)
	if err == nil {
		ret.LogicalCpuCnt = cnt
	}
	return ret
}

func GetSysInfo() SysInfo {
	return sys
}
