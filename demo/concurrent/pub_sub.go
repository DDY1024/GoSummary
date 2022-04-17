package concurrent

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 发布-订阅并发模型 --> 多播模型
// 1. 发布者会向每个订阅者 channel 塞入消息
// 2. 订阅者 channel 根据消息类型进行过滤

type (
	subscriber  chan interface{}
	filterTopic func(v interface{}) bool
)

type Publisher struct {
	sync.RWMutex
	bufSize     int                        // 订阅队列大小
	pubTimeOut  time.Duration              // 发布超时时间
	subscribers map[subscriber]filterTopic // 订阅者对象
}

func NewPublisher(size int, t time.Duration) *Publisher {
	return &Publisher{
		bufSize:     size,
		pubTimeOut:  t,
		subscribers: make(map[subscriber]filterTopic),
	}
}

// 清空订阅者
func (pub *Publisher) Close() {
	pub.Lock()
	defer pub.Unlock()
	for sub := range pub.subscribers {
		delete(pub.subscribers, sub)
		close(sub)
	}
}

// 摘除某个订阅者
func (pub *Publisher) Evict(sub chan interface{}) {
	pub.Lock()
	defer pub.Unlock()

	delete(pub.subscribers, sub)
	close(sub)
}

func (pub *Publisher) Publish(v interface{}) {
	pub.Lock()
	defer pub.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range pub.subscribers {
		wg.Add(1)
		go pub.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

func (pub *Publisher) sendTopic(sub subscriber, topic filterTopic, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(pub.pubTimeOut): // 存在问题：发送一个消息，新建一个 timer
	}
	return
}

func (pub *Publisher) Subscribe(topic filterTopic) chan interface{} {
	pub.Lock()
	defer pub.Unlock()
	sub := make(chan interface{}, pub.bufSize)
	pub.subscribers[sub] = topic
	return sub
}

func main() {
	pub := NewPublisher(10, time.Duration(100*time.Millisecond))
	all := pub.Subscribe(nil)
	cc := pub.Subscribe(func(v interface{}) bool {
		if vv, ok := v.(string); ok {
			return strings.Contains(vv, "go")
		}
		return false
	})

	pub.Publish("wang xi yang")
	pub.Publish("wang xi yang is a go coder")
	go func() {
		for v := range all {
			fmt.Println("All Subscriber: ", v)
		}
	}()

	go func() {
		for v := range cc {
			fmt.Println("Go Subscriber: ", v)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Service End")
}
