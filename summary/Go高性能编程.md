#### 内容整理自文章: https://mp.weixin.qq.com/s/wUrRQYz9rueAC8A4ugBYaQ
#### 1. 反射虽好，切莫贪杯
##### 1.1 strconv vs fmt
```go
// Bad
// 反射实现
for i := 0; i < b.N; i++ {
 s := fmt.Sprint(rand.Int())
}

BenchmarkFmtSprint-4    143 ns/op    2 allocs/op

// Good
for i := 0; i < b.N; i++ {
 //     strconv.FormatInt
 s := strconv.Itoa(rand.Int())
}

BenchmarkStrconv-4    64.2 ns/op    1 allocs/op
```

##### 1.2 思考反射使用的必要性
```go
// 本段代码实际上就是想利用反射达到一定的泛型化能力，过滤掉 list 中的指定元素
func DeleteSliceElms(i interface{}, elms ...interface{}) interface{} {
 // 构建 map set。
 m := make(map[interface{}]struct{}, len(elms))
 for _, v := range elms {
  m[v] = struct{}{}
 }
 // 创建新切片，过滤掉指定元素。
 v := reflect.ValueOf(i) 
 t := reflect.MakeSlice(reflect.TypeOf(i), 0, v.Len())
 for i := 0; i < v.Len(); i++ {
  if _, ok := m[v.Index(i).Interface()]; !ok {
   t = reflect.Append(t, v.Index(i))
  }
 }
 return t.Interface()
}

// 指定具体类型
// DeleteU64liceElms 从 []uint64 过滤指定元素。注意：不修改原切片。
func DeleteU64liceElms(i []uint64, elms ...uint64) []uint64 {
 // 构建 map set。
 m := make(map[uint64]struct{}, len(elms))
 for _, v := range elms {
  m[v] = struct{}{}
 }
 // 创建新切片，过滤掉指定元素。
 t := make([]uint64, 0, len(i))
 for _, v := range i {
  if _, ok := m[v]; !ok {
   t = append(t, v)
  }
 }
 return t
}

// benchmark 显示有好几倍的性能差距
go test -bench=. -benchmem main/reflect 
goos: darwin
goarch: amd64
pkg: main/reflect
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkDeleteSliceElms-12              1226868               978.2 ns/op           296 B/op         16 allocs/op
BenchmarkDeleteU64liceElms-12            8249469               145.3 ns/op            80 B/op          1 allocs/op
PASS
ok      main/reflect    3.809s
```

##### 1.3 慎用 binary.Read 和 binary.Write
- encoding/binary 包实现了数字和字节序列之间的简单转换以及 varints 的编码和解码
- encoding/binary 包实现，引入了反射，性能较低


#### 2. 避免重复的 string <--> []byte
```go
// Bad
for i := 0; i < b.N; i++ {
 w.Write([]byte("Hello world"))
}

BenchmarkBad-4   50000000   22.2 ns/op

// Good
data := []byte("Hello world") 
for i := 0; i < b.N; i++ {
 w.Write(data)
}

BenchmarkGood-4  500000000   3.25 ns/op
```

#### 3. 指定容器容量
```go
// 1. map
// 注意: map cap 并不保证完全的抢占式分配，而是用于估计所需的 hashmap bucket 数量
make(map[T1]T2, hint)

// 2. slice
make([]T, length, capacity)

```

#### 4. 字符串拼接
- `+`：行内字符串拼接
- `fmt.Sprintf`
- `strings.Builder`
- `bytes.Buffer`
- `[]byte`

#### 5. struct{}
- `set`：`map[T]struct{}`
- `channel`: `make(chan struct{})`
- `type T struct{}`


#### 6. 内存对齐
##### 为什么要内存对齐？
- CPU 以字长为单位访问内存
- 字节对齐可以减少 CPU 访存次数，提高 CPU 吞吐量

##### Go 内存对齐规则
- unsafe.Alignof 用于获取变量的对齐系数
- 对于任意类型的变量 x ，unsafe.Alignof(x) 至少为 1
- 对于结构体类型的变量 x，计算 x 每一个字段 f 的 unsafe.Alignof(x.f)，unsafe.Alignof(x) 等于其中的最大值
- 对于数组类型的变量 x，unsafe.Alignof(x) 等于构成数组的元素类型的对齐系数

##### struct 合理布局
```go
// 8 字节
type demo1 struct {
 a int8
 b int16
 c int32
}

// 12 字节
type demo2 struct {
 a int8
 c int32
 b int16
}

最优布局策略，field 从小到大进行排序
```

##### 空 struct 与 空数组对内存对齐的影响
- 空 struct{} 和 没有任何元素的 array 占据的内存空间大小为 0
- `正常情况下，空 struct{} 和 空 array 作为其它 struct 字段时，一般不需要内存对齐`
- `当空 struct{} 或 空 array 作为结构体最后一个字段时，需要进行内存对齐`
```go
// 4 byte
type demo3 struct {
 a struct{}
 b int32
}

// 8 byte
type demo4 struct {
 b int32
 a struct{}  // 作为最后一个字段，需要内存对齐
}

func main() {
 fmt.Println(unsafe.Sizeof(demo3{})) // 4
 fmt.Println(unsafe.Sizeof(demo4{})) // 8
}
```

#### 7. 减少变量逃逸
##### 7.1. 逃逸一般发生的情况
- `变量较大`: 不同 Go 版本的大小限制可能不一样，一般 < 64KB，局部变量将不会逃逸到堆上
- 变量大小不确定
- 变量类型不确定
- 返回指针
- 返回引用
- 闭包

##### 7.2 小的拷贝好过引用
- 减少变量逃逸
- `对于一些短小的对象，栈上复制的成本远小于在堆上分配和回收操作`

```go
const capacity = 1024

// 数组拷贝
func arrayFibonacci() [capacity]int {
 var d [capacity]int
 for i := 0; i < len(d); i++ {
  if i <= 1 {
   d[i] = 1
   continue
  }
  d[i] = d[i-1] + d[i-2]
 }
 return d
}

// 切片引用
func sliceFibonacci() []int {
 d := make([]int, capacity)
 for i := 0; i < len(d); i++ {
  if i <= 1 {
   d[i] = 1
   continue
  }
  d[i] = d[i-1] + d[i-2]
 }
 return d
}

// -gcflags="-l": 禁止内联优化
// -gcflags="-m": 逃逸分析
go test -bench=. -benchmem -gcflags="-l" main/copy
```

##### 7.3 返回值 vs 返回指针
- `返回值`：拷贝整个对象
- `返回指针`：避免拷贝，但会导致内存分配逃逸到堆中，增加垃圾回收的负担。
- 一般，修改原对象值或较大 struct，采用指针处理；对于小 struct ，直接返回值会更好一些

##### 7.4 返回值使用确定类型
- 如果变量类型不确定，那么将会逃逸到堆上；所以，`函数返回值如果能确定的类型，就不要使用 interface{}`

#### 8. sync.Pool 复用对象
##### 特性
- sync.Pool 中对象可能被无通知的回收（GC）
- sync.Pool 可伸缩，并发安全，其容量受限于内存大小
- 利用 sync.Pool 进行对象复用的场景还是蛮多的

```go
// bytes.Buffer 复用
var bufferPool = sync.Pool{
 New: func() interface{} {
  return &bytes.Buffer{}
 },
}

var data = make([]byte, 10000)

func BenchmarkBufferWithPool(b *testing.B) {
 for n := 0; n < b.N; n++ {
  buf := bufferPool.Get().(*bytes.Buffer)
  buf.Write(data)
  buf.Reset()
  bufferPool.Put(buf)
 }
}

func BenchmarkBuffer(b *testing.B) {
 for n := 0; n < b.N; n++ {
  var buf bytes.Buffer
  buf.Write(data)
 }
}
```

#### 9. 锁
- 无锁化
  - `lock-free 数据结构`: 利用硬件支持的原子操作可以实现无锁的数据结构，原子操作可以在 lock-free 的情况下保证并发安全（CAS）
    - `for + CAS`：for 循环 + 配合 CAS 操作
  - 串行无锁
    - 多个 loader 返回多个独立的结果，然后最终进行汇总；避免 loader 边返回结果，边对一个 map 进行并发操作
- 分片 lock
  - bigcache、freecache；需要设置合理的分片数，不宜过多
- 优先共享锁而不是互斥锁
  - `sync.RWMutex > sync.Mutex`

```go
type LockFreeList struct {
	Head atomic.Value  // atomic.Value
}

// for 循环 + cas 操作
func (l *LockFreeList) Push(v interface{}) {
	for {
		head := l.Head.Load()
		headNode, _ := head.(*Node)
		n := &Node{
			Value: v,
			Next:  headNode,
		}
        // cas 操作能够保证并发安全
		if l.Head.CompareAndSwap(head, n) {
			break
		}
	}
}

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
```

#### 10. goroutine 数量控制
##### 10.1. 协程数过多
- 内存开销: 一个 goroutine 初始栈空间大小为 4KB
- 调度开销: `runtime.Gosched()`，协程切换消耗 cpu 时间在 ns 级，线程切换消耗 cpu 时间在 us 级
- GC 压力: go 内存回收

##### 10.2 限制协程数量
- 缓冲 channel：`make(chan struct{}, 3)`
- 携程池化
  - https://github.com/Jeffail/tunny
  - **ants**

#### 11. sync.Once 
- 单例模式
- init vs sync.Once
  - init 在 pkg 首次加载时被执行
  - sync.Once 属于懒处理，在一些场景下 sync.Once 会引入问题，慎用；
    - 如果能够确定 sync.Once 懒处理没啥问题，在一些初始化的场景下比较有用


#### 12. sync.Cond 不建议使用
- 后续再进行研究