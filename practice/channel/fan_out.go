package main

import "reflect"

// fan-out 两种模式
// 1. 广播模式：全部
// 2. 单播模式：随机选择一个

func FanOutByGoroutine(ch <-chan interface{}, out []chan interface{}, async bool) {
	defer func() {
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	for v := range ch {
		v := v
		for i := 0; i < len(out); i++ { // 广播模式
			i := i
			if async {
				go func() {
					out[i] <- v
				}()
			} else {
				out[i] <- v
			}
		}
	}
}

func FanOutByReflect(ch <-chan interface{}, out []chan interface{}) {
	defer func() {
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	cases := make([]reflect.SelectCase, len(out))
	for i := range cases {
		cases[i].Dir = reflect.SelectSend
		cases[i].Chan = reflect.ValueOf(out[i])
	}

	for v := range ch {
		v := v

		// random
		for i := range cases {
			cases[i].Send = reflect.ValueOf(v)
		}
		reflect.Select(cases)
	}
}

// round-robin
func FanOutByRoundRobin(ch <-chan interface{}, out []chan interface{}) {
	defer func() {
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	idx, n := -1, len(out)
	for v := range ch {
		idx = (idx + 1) % n
		out[idx] <- v
	}
}
