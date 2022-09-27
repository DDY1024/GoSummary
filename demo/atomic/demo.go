package main

// 支持的几大操作
// 1. 增减操作，AddXXX
// 2. 载入操作，LoadXXX
// 3. 存储操作，StoreXXX
// 4. 比较并交换，CAS 操作，常用
// 5. 交换操作（很少用）

// mutex vs atomic
// 1. mutex 更多的是用来保护一段临界区代码，atomic 更多的用户保护一个变量
// 2. mutex 由操作系统调度器实现；而 atomic 操作由底层硬件指令直接支持（lock-free）

// CAS 类操作
// 1. 常见使用方式: for 循环 + cas 操作

func casDemo() {
	// atomic.LoadInt32()
	// atomic.CompareAndSwapInt64()  // 数据库操作中的乐观锁
	// atomic.CompareAndSwapPointer()
}

// atomic.Value 针对任意数据类型的读写操作的原子性

func valueDemo() {
	// var data atomic.Value
	// data.Load()
	// data.Store()
	// data.CompareAndSwap()
	// data.Swap()
}
