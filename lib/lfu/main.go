package lfu

// 实现参考: http://dhruvbird.com/lfu.pdf
// 本 lfu cache 实现主要用于热 key 防御
// 热 key 检测策略: hit >= total * threshold
//
// 如何实现 O(1) 复杂度的驱逐策略？
// 思路: 维护一个递增 freq list，其中指定频率情况下用 map 维护 freq = x 的所有元素
//
//
// 关于 runtime.SetFinalizer 在 cache 中的应用可以参考: https://zhuanlan.zhihu.com/p/76504936
//
//
// LFU 热 key 防御使用建议
// 1. 对于 cache size 设置为 1000 即可有效的拦截热度为 1% 的热 key
// 2. 为了保证 cache 对热 key 感知的灵敏度，"清理周期"不宜设置的过长，1s就足够了
// 3. key 的过期时间会直接影响数据的准确性（需要结合具体服务设定合适的值）
// 4. 如果想只拦截热度达到某个阈值的热 key，可以在 Get 操作中设置 key 阈值，降低数据不一致对服务的影响
//
// 实验结论证明: lfu 较 lru 在热 key 命中率上更优
//
//
// 进一步优化：
// lfu 实现中对于 kv、freq_list 操作会进行加锁，为了降低锁的粒度，同时创建多个 lfu cache，进行 hash 取余选取
// 大锁拆小锁的方式

import (
	"container/list"
	"errors"
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	*cache
}

type cache struct {
	values map[interface{}]*cacheEntry
	freqs  *list.List // freq 链表，单调递增，freq = 1, 2, 3, 4, ..., N

	// If len > maxLen, cache will automatically evict down to MaxLen
	len    int
	maxLen int // 驱逐阈值

	cleanInterval time.Duration
	janitor       *janitor

	lock  *sync.Mutex
	name  string
	total int // 统计总的请求次数
}

type cacheEntry struct {
	key        interface{}
	value      interface{}
	expiresAt  *time.Time
	lastAccess time.Time
	freqNode   *list.Element
	hit        int
}

type listEntry struct {
	entries map[*cacheEntry]byte // 标记
	freq    int
}

// maxLen: cache max size
// ci: delete keys that have expired or not accessed for long time
func newCache(maxLen int, ci time.Duration, name string) *cache {
	c := new(cache)
	c.values = make(map[interface{}]*cacheEntry)
	c.freqs = list.New()
	c.lock = new(sync.Mutex)
	c.maxLen = maxLen
	c.name = name
	c.cleanInterval = ci
	// c.enableMetrics = len(name) != 0
	return c
}

// threshold: 热 key 判定标准
func (c *cache) Get(key interface{}, threshold float64) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.total++
	if e, ok := c.values[key]; ok {
		if e.expiresAt != nil {
			if e.expiresAt.After(time.Now()) {
				return c.value(e, threshold)
			}
			// 当前 key 过期，直接删除
			c.del(e)
			return nil, nil
		}
		return c.value(e, threshold)
	}
	return nil, nil
}

func (c *cache) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next() // 下一个节点
	}

	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		// create a new list entry
		li := new(listEntry)
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry]byte)
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	e.lastAccess = time.Now()
	e.hit++
	nextPlace.Value.(*listEntry).entries[e] = 1
	if currentPlace != nil {
		c.remEntry(currentPlace, e)
	}
}

func (c *cache) remEntry(place *list.Element, entry *cacheEntry) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}

func (c *cache) value(e *cacheEntry, threshold float64) (interface{}, error) {
	c.increment(e)
	// 热 key 诊断: hit >= total * threshold
	if e.hit <= 1 || float64(e.hit) < float64(c.total)*threshold {
		return nil, errors.New("Not hit threshold")
	}
	return e.value, nil
}

func (c *cache) Set(key interface{}, value interface{}, expiresAt *time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		// value already exists for key.  overwrite
		e.value = value
		e.expiresAt = expiresAt
		c.increment(e)
	} else {
		// value doesn't exist.  insert
		e := new(cacheEntry)
		e.key = key
		e.value = value
		e.expiresAt = expiresAt
		c.values[key] = e
		c.increment(e)
		c.len++
		// 元素数大于两倍 maxLen，执行驱逐策略
		if c.len > c.maxLen*2 {
			c.evict(c.len - c.maxLen)
		}
	}
}

func (c *cache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if entry, ok := c.values[key]; ok {
		c.del(entry)
	}
}

func (c *cache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.len
}

func (c *cache) Evict(count int) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.evict(count)
}

// 驱逐策略: 按照 freq 从小到大依次驱逐
func (c *cache) evict(count int) int {
	// No lock here so it can be called
	// from within the lock (during Set)
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry := range place.Value.(*listEntry).entries {
				if i < count {
					c.del(entry)
					evicted++
					i++
				}
			}
		}
	}
	return evicted
}

func (c *cache) del(entry *cacheEntry) {
	delete(c.values, entry.key)
	c.remEntry(entry.freqNode, entry)
	c.len--
}

func (c *cache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	purgeTime := time.Now().Add(-c.cleanInterval)
	for _, e := range c.values {
		e.hit = 0 // reset hit
		if e.expiresAt != nil && e.expiresAt.Before(time.Now()) ||
			e.lastAccess.Before(purgeTime) {
			c.del(e)
		}
	}

	c.total = 0 // reset total
}

func (c *cache) MGet(keys []interface{}, threshold float64) (map[interface{}]interface{}, map[interface{}]error) {
	r := make(map[interface{}]interface{})
	e := make(map[interface{}]error)
	for _, k := range keys {
		r[k], e[k] = c.Get(k, threshold)
	}

	return r, e
}

func (c *cache) MSet(kvs map[interface{}]interface{}, expiresAt *time.Time) {
	for k, v := range kvs {
		c.Set(k, v, expiresAt)
	}
}

// 定时清理策略
type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.Purge()
		case <-j.stop: // 停止
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

// ci: lfu local cache 清除过期 key 的周期
func New(maxLen int, ci time.Duration, name string) *Cache {
	c := newCache(maxLen, ci, name)
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		// 在 cache 对象没有再被引用时，停止后台定时清理的 goroutine
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}
