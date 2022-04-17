package timer

import (
	"fmt"
	"time"
)

// 1. bad case: for + select 模式中，创建过多无效的 timer
func runOne(messages <-chan string) {
	for {
		select {
		// 1. time.After 每次会 new 一个 timer，造成 timer 泄漏
		case <-time.After(time.Minute):
			return
		// 2. 每发送一条 message，便会 new 一个 timer
		case msg := <-messages:
			fmt.Println(msg)
		}
	}
}

// 2. good case: for + select 模式中，timer 尽可能复用
func runTwo(messages <-chan string) {
	timer := time.NewTimer(time.Minute)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return
		case msg := <-messages:
			fmt.Println(msg)
			// 1. 如果 timer 已经过期，则 stop 返回 false，timer.C 仍然存在一条过期通知消息，需要取出
			// 2. 如果 timer 没有过期，则 stop 返回 true，继续执行后续操作
			if !timer.Stop() {
				<-timer.C
			}
		}
		// 每处理完一条消息后重置下 timer
		timer.Reset(time.Minute)
	}
}
