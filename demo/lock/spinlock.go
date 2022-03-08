package lock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// 自旋锁实现，参考自: https://github.com/tidwall/spinlock/blob/master/locker.go
//
// https://zh.wikipedia.org/wiki/%E8%87%AA%E6%97%8B%E9%94%81
// https://www.cnblogs.com/cxuanBlog/p/11679883.html
//
// 由于线程在这一过程中保持执行，因此是一种忙等待。一旦获取了自旋锁，线程会一直保持该锁，直至显式释放自旋锁。
// 自旋锁避免了进程上下文的调度开销，因此对于线程只会阻塞很短时间的场合是有效的。
// 因此操作系统的实现在很多地方往往用自旋锁。
//
// *****自旋锁优缺点*****
// 自旋锁能尽可能的减少线程的阻塞，对于锁的竞争不激烈，且占用锁时间非常短的代码块来说性能能大幅度的提升，
// 因为自旋的消耗会小于线程阻塞挂起再唤醒的操作的消耗，这些操作会导致线程发生两次上下文切换
//
// 如果锁的竞争激烈，或者持有锁的线程需要长时间占用锁执行同步块，这时候就不适合使用自旋锁了，
// 因为自旋锁在获取锁前一直都是占用 CPU 做无用功，占着 XX 不 XX，同时有大量线程在竞争一个锁，
// 会导致获取锁的时间很长，线程自旋的消耗大于线程阻塞挂起操作的消耗，其它需要 cpu 的线程又不能获取到 cpu，
// 造成 cpu 的浪费。所以这种情况下我们要关闭自旋锁。
//
// Locker is a spinlock implementation.
//
// A Locker must not be copied after first use.
type Locker struct {
	// noCopy 应用可以参考这篇文章: https://jiajunhuang.com/articles/2018_11_12-golang_nocopy.md.html
	// noCopy
	_    sync.Mutex // for copy protection compiler warning
	lock uintptr
}

// Lock locks l.
// If the lock is already in use, the calling goroutine
// blocks until the locker is available.
func (l *Locker) Lock() {
loop: // 不断循环等待
	if !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched() // 调度 CPU
		goto loop
	}
}

// Unlock unlocks l.
func (l *Locker) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}
