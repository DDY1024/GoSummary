package channel

import (
	"reflect"
	"sync"
)

// or-channel 工作模式: 同时监听多个信号，任何一个信号返回，便执行后续逻辑，同时忽略其它信号（first-return-process）

// 1. goroutine
func OrByGoroutine(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		var once sync.Once
		for _, ch := range chs {
			go func(ch <-chan interface{}) {
				select {
				case <-out: // 第一个 close(out) 之后，后续的 channel 不再监听
				case <-ch:
					once.Do(func() { close(out) }) // 接收到任意信号以后，仅 close(out) 一次，进行通知
				}
			}(ch)
		}
	}()

	return out
}

// 2. reflect.Select
func OrByReflect(chs ...<-chan interface{}) <-chan interface{} {

	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases) // reflect.Select 监听 **不定数目** 的任一 case 的返回，然后执行后续操作
	}()

	return orDone
}

// 3. 递归合并 channel
func OrByRecursion(chs ...<-chan interface{}) <-chan interface{} {

	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	// 简单的只做信号通知
	orDone := make(chan interface{})

	go func() {
		defer close(orDone)
		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		default:
			// 二分递归进行处理
			m := len(chs) / 2
			select {
			case <-OrByRecursion(chs[:m]...): // 左半部分 channel 合并
			case <-OrByRecursion(chs[m:]...): // 右半部分 channel 合并
			}
		}
	}()

	return orDone
}
