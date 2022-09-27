package main

import (
	"fmt"
)

// 观察者模式: 广播、多播机制
// 发布-订阅机制
// https://mp.weixin.qq.com/s/4NqjkXVqFPamEc_QsyRipA

type ISubject interface {
	Register(observer IObsever) // 注册添加一个观察者
	Remove(observer IObsever)   // 移除一个观察者
	Notify(observer IObsever)   // 消息通知
}

type IObsever interface { // 定义观察者的接口类型
	Update(msg string)
}

// 订阅对象
type Subject struct {
	observers []IObsever
}

func (sub *Subject) Register(observer IObsever) {
	sub.observers = append(sub.observers, observer)
}

func (sub *Subject) Remove(observer IObsever) {
	for i, ob := range sub.observers {
		if ob == observer {
			sub.observers = append(sub.observers[:i], sub.observers[i+1:]...)
		}
	}
}

func (sub *Subject) Notify(msg string) {
	for _, o := range sub.observers {
		o.Update(msg)
	}
}

type Obsever1 struct{}

func (Obsever1) Update(msg string) {
	fmt.Printf("Obsever1: %s", msg)
}

type Obsever2 struct{}

func (Obsever2) Update(msg string) {
	fmt.Printf("Obsever2: %s", msg)
}
