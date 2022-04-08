package lock

import (
	"runtime"
	"sync"
)

const (
	cacheLineSize = 64
)

var (
	shardsLen int
)

type RWLock []RWMutexShard

type RWMutexShard struct {
	_ [cacheLineSize]byte
	sync.RWMutex
}

func init() {
	// 分段锁个数，一般为 P 的个数
	shardsLen = runtime.GOMAXPROCS(0)
}

func New() RWLock {
	return RWLock(make([]RWMutexShard, shardsLen))
}

func (this RWLock) Lock() {
	for shard := range this {
		this[shard].Lock()
	}
}

func (this RWLock) Unlock() {
	for shard := range this {
		this[shard].Unlock()
	}
}

func getPid() int {
	return 0xff
}

func (this RWLock) RLocker() sync.Locker {
	tid := getPid()
	return this[tid%shardsLen].RWMutex.RLocker()
}
