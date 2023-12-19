package main

// 支持的操作
// 1. AddXXX
// 2. LoadXXX
// 3. StoreXXX
// 4. CAS 操作
// 5. Swap

// Mutex vs Atomic
// 1. mutex 更多用于保护一段临界区代码，atomic 更多用于保护一个变量
// 2. mutex 由操作系统调度实现，atomic 由底层硬件指令支持（lock-free）

// 自旋锁：for + CAS

func CAS() {
	// atomic.CompareAndSwapInt64()
	// atomic.CompareAndSwapPointer()
	// atomic.CompareAndSwapUintptr()
}

// atomic.Value 任意数据类型的原子读写操作
