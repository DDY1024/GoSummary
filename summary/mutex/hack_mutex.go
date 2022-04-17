package Mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 鸟窝文章: https://time.geekbang.org/column/article/296793

const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

type Mutex struct{ sync.Mutex }

/*
type Mutex struct {
	state int32
	sema  uint32
}
*/

// 获取竞争锁的 goroutine 数目
func (m *Mutex) GoroutineNum() int {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift // 得到等待者的数值
	v = v + (v & mutexLocked) // 再加上锁持有者的数量 0 或 1
	return int(v)
}

// 锁是否已经被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// 是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}
