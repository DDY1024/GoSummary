package goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// 一个简单的 goroutine local storage 实现
// 具体参考: https://chai2010.cn/advanced-go-programming-book/ch3-asm/ch3-08-goroutine-id.html
// 开源实现: https://github.com/modern-go/gls

// 应用场景: goroutine 内部的局部存储
// 性能优化: 分段锁优化
type GLS struct {
	m map[int64]map[interface{}]interface{} // goroutine_id --> kv map
	sync.Mutex
}

func NewGLS() *GLS {
	return &GLS{
		m: make(map[int64]map[interface{}]interface{}),
	}
}

func GetGoroutineID() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		// false: 只输出当前 goroutine 堆栈信息并返回写入字节数
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}
	return int64(id)
}

func (gls *GLS) GetMap() map[interface{}]interface{} {
	gls.Lock()
	defer gls.Unlock()

	goid := GetGoroutineID()
	if m := gls.m[goid]; m != nil {
		return m
	}

	m := make(map[interface{}]interface{})
	gls.m[goid] = m
	return m
}

func (gls *GLS) Get(key interface{}) interface{} {
	return gls.GetMap()[key]
}

func (gls *GLS) Put(key interface{}, v interface{}) {
	gls.GetMap()[key] = v
}

func (gls *GLS) Delete(key interface{}) {
	delete(gls.GetMap(), key)
}

func (gls *GLS) Clean() {
	gls.Lock()
	defer gls.Unlock()

	delete(gls.m, GetGoroutineID())
}

func main() {
	var wg sync.WaitGroup
	gls := NewGLS()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			defer gls.Clean()

			defer func() {
				fmt.Printf("%d: number = %d\n", idx, gls.Get("number"))
			}()
			gls.Put("number", idx+100)
		}(i)
	}
	wg.Wait()
}
