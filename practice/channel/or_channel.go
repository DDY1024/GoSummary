package main

import (
	"reflect"
	"sync"
)

// or-channel: 同时监听多个 channel，任意一个 channel 准备好，接收数据进行处理并忽略其它 channel

func OrByGoroutine(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		var once sync.Once
		for _, ch := range chs {
			go func(ch <-chan interface{}) {
				select {
				case <-out:
				case <-ch:
					once.Do(func() { close(out) })
				}
			}(ch)
		}
	}()

	return out
}

// 2. reflect.Select
func OrByReflect(chs ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	if len(chs) == 1 {
		return chs[0]
	}

	go func() {
		defer close(done)
		var cases []reflect.SelectCase
		for _, ch := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}
		reflect.Select(cases) // reflect.Select 任一满足，直接返回
	}()
	return done
}

// 3. 分治合并 channel
func OrByRecursion(chs ...<-chan interface{}) <-chan interface{} {
	if len(chs) == 1 {
		return chs[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		mid := len(chs) / 2
		select {
		case <-OrByRecursion(chs[:mid]...):
		case <-OrByRecursion(chs[mid:]...):
		}
	}()
	return orDone
}
