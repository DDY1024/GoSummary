package sync

// sync.Cond 条件变量，提供 watch 机制
// 1. broadcast 所有等待的 goroutine
// 2. signal 某一个等待的 goroutine
//
// 注意事项
// 1. sync.Cond 初始化时需要传入一个 Locker 接口的实现，这个 Locker 用于保护条件变量
// 2. Broadcast 和 Signal 不需要加锁调用，但是调用 Wait 时需要加锁。
// 3. Wait 执行中有解锁重加锁的过程，在这个期间对临界区是没有保护的。
//    3.1 调用Wait 前 Lock --> Wait 内部解锁  --> 不受保护期 --> Wait 内部加锁
// 4. for { Wait() } 通过 for 循环来检查条件是否满足，因为随时都可以触发通知
//

// a. broadcast
/*
	var m sync.Mutex
	c := sync.NewCond(&m)
	ready := make(chan struct{})
	isReady := false
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			m.Lock()
			time.Sleep(rand...)
			ready <- struct{}{}
			for !isReady {
				c.Wait()
			}
			m.Unlock()
		}()
	}

	for i := 0; i < 10; i++ {
		<-ready
	}

	isReady = true
	c.Broadcast()
*/
