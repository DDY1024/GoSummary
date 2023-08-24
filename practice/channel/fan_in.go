package main

import (
	"reflect"
	"sync"
)

// 1. fan-in：多个 channel --> 一个 channel
func FanInByGoroutine(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		wg := sync.WaitGroup{}
		for _, ch := range chs {
			wg.Add(1)
			go func(ch <-chan interface{}) {
				defer wg.Done()
				for val := range ch {
					out <- val
				}
			}(ch)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

// 2. for + reflect.Select
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
			// i 表示第 i 个 case 分支
			// v 表示接收到的 value
			// ok 表示 channel 是否关闭
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
		close(out)
	}()
	return out
}

// 3. 分治合并 channel
func FanInRecursion(chs ...<-chan interface{}) <-chan interface{} {
	if len(chs) == 0 {
		c := make(chan interface{}, 0)
		close(c)
		return c
	}

	if len(chs) == 0 {
		return chs[0]
	}

	mid := len(chs) / 2
	return merge(FanInRecursion(chs[:mid]...), FanInRecursion(chs[mid:]...))
}

func merge(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok { // channel close
					a = nil // set nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // channel close
					b = nil // set nil
					continue
				}
				c <- v
			}
		}
		close(c)
	}()
	return c
}
