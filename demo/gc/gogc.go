package gc

// 优化手段: 动态调整 GOGC 参数来进行 GC 优化（自适应降频）
// 原文参考: https://xargin.com/dynamic-gogc/

// 1. GOGC 计算公式
// GOGC= (hard_target - live_dataset)/live_dataset *100
// hard_target: 容器内存大小，保险起见我们使用 70%
// live_dataset: 存活对象占用空间大小，运行时可获取
///////// m := &runtime.MemStats{}
///////// runtime.ReadMemStats(m)
///////// 使用 m.HeapSys 来近似表示 live_dataset
//
// 2. 运行时设置 GOGC
// debug.SetGCPercent(GOGC)
//
// 3. runtime.SetFinalizer 妙用
// 这个方法的作用是为一个对象设置一个关联函数，当这个对象在某次 GC 中变成不可达状态时，触发关联函数，并且取消关联。由于关联函数的存在，此对象的状态又变成了可达，本次 GC 不会回收这个对象。
// 但在下次 GC 时，此对象又会变成不可达状态，然后被回收。(该对象会被延长一个 GC 周期)
// 因此，存在这样一种方案，我们可以人为 new 一个对象，然后通过 runtime.SetFinalizer 为该对象关联一个操作函数
// 这个操作函数我们可以动态调整 GOGC 的值，这样不断 runtime.SetFinalizer 相当于该对象一直没被 GC 回收，且每次 GC
// mark 时都会执行该操作。具体代码如下所示:
/*
type finalizer struct {
	C   chan time.Time
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

func finalizerHandler(f *finalizerRef) {
	select {
	case f.parent.C <- time.Now():
	default:
	}
	runtime.SetFinalizer(f, finalizerHandler)
}

func newTicker() *finalizer {
	f := &finalizer{
		C: make(chan time.Time, 1),
	}
	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil  // 取消 ref 指针应用，finalizerRef 在 GC mark 时会被判定为不可达对象，进一步触发 finalizerHandler 关联函数
	return f
}

// 定时调整 GOGC 参数
newTicker()
go func() {
		for t := range newTicker().C {  // 一次 GC 触发一次调用
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)
			total := m.HeapSys
			live := m.HeapSys - m.HeapIdle
			gogc := 0

			// 动态调整计算
			if mhard > total && live > 0 {
				gogc = int((mhard - live) * 100 / live)
			}

			// 设置一个 GOGC 上限，动态调整不允许超过该上限
			if gogc > upperLimit {
				gogc = upperLimit
			}

			// 运行时动态调整 GOGC 参数
			debug.SetGCPercent(gogc)
		}
}()
*/
