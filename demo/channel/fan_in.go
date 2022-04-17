package channel

import (
	"reflect"
	"sync"
)

// fan-in 模式：多个 channel 合并为一个 channel

// 1. 方式一: goroutine 汇总
func FanInByGoroutine(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chs))

		for _, c := range chs {
			go func(c <-chan interface{}) {
				defer wg.Done()
				for v := range c {
					out <- v
				}
			}(c)
		}

		wg.Wait()
		close(out) // 不再有数据传输，直接关闭 channel
	}()

	return out
}

// 2. for 循环 + reflect.Select
func FanInByReflect(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 接收的 channel 已经关闭，则剔除掉
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}

		close(out)
	}()

	return out
}

// 3. 递归处理方式(二分、归并)
func FanInRecursion(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chs[0]
	case 2:
		return mergeTwo(chs[0], chs[1])
	default:
		m := len(chs) / 2
		return mergeTwo(
			FanInRecursion(chs[:m]...),
			FanInRecursion(chs[m:]...))
	}
}

func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok { // channel 关闭
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // channel 关闭
					b = nil
					continue
				}
				c <- v
			}
		}
		close(c)
	}()
	return c
}
