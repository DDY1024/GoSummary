package main

import (
	"fmt"
	"reflect"
	"sync"
)

// 观察者模式：广播/多播机制、发布/订阅模式
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

// 一种事件总线的实现方式
type Bus interface {
	Subscribe(topic string, handler interface{}) error
	Publish(topic string, args ...interface{})
}

// AsyncEventBus 异步事件总线
type AsyncEventBus struct {
	handlers map[string][]reflect.Value
	lock     sync.Mutex
}

// NewAsyncEventBus new
func NewAsyncEventBus() *AsyncEventBus {
	return &AsyncEventBus{
		handlers: map[string][]reflect.Value{},
		lock:     sync.Mutex{},
	}
}

// Subscribe 订阅
func (bus *AsyncEventBus) Subscribe(topic string, f interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	v := reflect.ValueOf(f)
	if v.Type().Kind() != reflect.Func { // reflect.Func
		return fmt.Errorf("handler is not a function")
	}

	handler, ok := bus.handlers[topic]
	if !ok {
		handler = []reflect.Value{}
	}

	handler = append(handler, v)
	bus.handlers[topic] = handler
	return nil
}

// Publish 发布
func (bus *AsyncEventBus) Publish(topic string, args ...interface{}) {
	handlers, ok := bus.handlers[topic]
	if !ok {
		fmt.Println("not found handlers in topic:", topic)
		return
	}

	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	for i := range handlers {
		go handlers[i].Call(params) // 异步执行，并不会等待返回结果 go func
	}
}
