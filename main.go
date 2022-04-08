package main

import (
	"fmt"
	"sync/atomic"
)

// LockFreeList 无锁单向链表
type LockFreeList struct {
	Head atomic.Value
}

// Push 有锁
func (l *LockFreeList) Push(v interface{}) {
	for {
		head := l.Head.Load()
		headNode, _ := head.(*Node)
		n := &Node{
			Value: v,
			Next:  headNode,
		}
		if l.Head.CompareAndSwap(head, n) {
			break
		}
	}
}

// String 有锁链表的字符串形式输出
func (l LockFreeList) String() string {
	s := ""
	cur := l.Head.Load().(*Node)
	for {
		if cur == nil {
			break
		}
		if s != "" {
			s += ","
		}
		s += fmt.Sprintf("%v", cur.Value)
		cur = cur.Next
	}
	return s
}
