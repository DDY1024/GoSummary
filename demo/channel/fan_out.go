package channel

import "reflect"

// fan-out 分为两种模式
// 1. 广播模式：输入 channel 中的数据，广播传输给所有下游 channel
// 2. 单播模式：输入 channel 中的数据，随机选择一个下游 channel 进行传输

// 1. goroutine 方式
func FanOutByGoroutine(ch <-chan interface{}, out []chan interface{}, async bool) {
	defer func() {
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	for v := range ch {
		v := v
		for i := 0; i < len(out); i++ {
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

// 2. reflect
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

		// 随机选择一个下游 channel 进行传输
		for i := range cases {
			cases[i].Send = reflect.ValueOf(v)
		}
		reflect.Select(cases)

		// reflect.Select 每次只会选择一个 channel 进行 send 操作
		// 下述操作相当于模拟了广播数据这样的场景
		// for i := range cases {
		// 	cases[i].Chan = reflect.ValueOf(out[i])
		// 	cases[i].Send = reflect.ValueOf(v)
		// }
		// for _ = range cases { // for each channel
		// 	chosen, _, _ := reflect.Select(cases)
		// 	cases[chosen].Chan = reflect.ValueOf(nil)  // 选择完一个 channel 后关闭，后续从剩下 channel 选取
		// }
	}
}

// 3. 轮询策略
func FanOutByRoundRobin(ch <-chan interface{}, out []chan interface{}) {
	defer func() {
		for i := 0; i < len(out); i++ {
			close(out[i])
		}
	}()

	idx, n := 0, len(out)
	for v := range ch {
		v := v
		out[idx] <- v
		idx = (idx + 1) % n
	}
}
