package Mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 鸟窝文章：https://time.geekbang.org/column/article/296793
const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	// 1. 锁初始状态: 直接抢锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 2. 检测当前锁的状态
	//    a. 锁被占有
	//    b. 饥饿状态: 存在 goroutine 长时间没有占用到锁
	//    c. 唤醒状态: 存在 goroutine 被唤醒抢锁
	// 上述这几种情况，我们不参与锁的竞争
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// 3. 尝试抢锁
	new := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, new)
}
