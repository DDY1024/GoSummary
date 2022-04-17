package runtime

import (
	"fmt"
	"runtime"
	"time"
)

// 容器化部署 go 应用一些监控指标
func main() {
	var m0, m1 runtime.MemStats

	// 1. 时间点一
	runtime.ReadMemStats(&m0)

	lastCgoCall := runtime.NumCgoCall()
	time.Sleep(10 * time.Second)

	// 2. 时间点二
	runtime.ReadMemStats(&m1)

	// runtime.alloc.objects
	fmt.Println(m1.Mallocs - m0.Mallocs)

	// runtime.alloc.bytes
	// TotalAlloc is cumulative bytes allocated for heap objects.
	fmt.Println(m1.TotalAlloc - m0.TotalAlloc)

	// runtime.gc.count
	fmt.Println(m1.NumGC - m0.NumGC)

	// runtime.goroutine.num
	fmt.Println(runtime.NumGoroutine()) // 实时统计值

	// runtime.cgocall.count
	fmt.Println(runtime.NumCgoCall() - lastCgoCall)

	// runtime.gc.pause
	fmt.Println(GCPauseNs(m1, m0))

	// cgroup 三个关键的监控指标
	// cgroup.cpustat.nr_throttled
	// cgroup.cpustat.nr_periods
	// cgroup.cputime
	//
	//
	// 1. /proc/self/cgroup
	// 2. /sys/fs/cgroup/cpu,cpuacct/{pod_name}
	// 3. cpu.stat、cpuacct.stat、cpuacct.usage
	//
	// cur.nr_throttled - last.nr_throttled
	// cur.nr_periods - last.nr_periods
	// cur.cpu_time - last.cpu_time
}

func GCPauseNs(new runtime.MemStats, old runtime.MemStats) uint64 {
	if new.NumGC <= old.NumGC {
		return new.PauseNs[(new.NumGC+255)%256]
	}

	n := new.NumGC - old.NumGC
	if n > 256 {
		n = 256
	}

	// The most recent pause is at PauseNs[(NumGC+255)%256]
	// NumGC - 0 + 255
	// NumGC - 1 + 255
	// ....
	// NumGC - k + 255
	var maxPauseNs uint64
	for i := uint32(0); i < n; i++ {
		if pauseNs := new.PauseNs[(new.NumGC-i+255)%256]; pauseNs > maxPauseNs {
			maxPauseNs = pauseNs
		}
	}

	return maxPauseNs
}
