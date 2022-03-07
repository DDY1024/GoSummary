package lru

import "container/list"

// leetcode 中关于 lru cache 的实现
// https://leetcode-cn.com/problems/lru-cache/
// 生产环境中 lru cache 的实现参考：https://github.com/hashicorp/golang-lru

type Item struct {
	key int
	val int
}

type LRUCache struct {
	visit *list.List
	index map[int]*list.Element
	cap   int
	size  int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		visit: list.New(),
		index: make(map[int]*list.Element, 0),
		cap:   capacity,
		size:  0,
	}
}

func (self *LRUCache) Get(key int) int {
	e, ok := self.index[key]
	if !ok {
		return -1
	}
	self.visit.Remove(e)
	e = self.visit.PushFront(e.Value)
	self.index[key] = e
	return e.Value.(*Item).val
}

func (self *LRUCache) Put(key int, value int) {
	if e, ok := self.index[key]; ok {
		e.Value.(*Item).val = value
		self.visit.Remove(e)
		e = self.visit.PushFront(e.Value)
		self.index[key] = e
		return
	}

	item := &Item{
		key: key,
		val: value,
	}

	e := self.visit.PushFront(item)
	self.index[key] = e
	self.size++

	if self.size > self.cap {
		self.size--
		item = self.visit.Back().Value.(*Item)
		delete(self.index, item.key)
		self.visit.Remove(self.visit.Back())
	}
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
