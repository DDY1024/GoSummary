#### 0. goroutine 内部
- 在一个 goroutine 内部，程序的执行顺序和其代码指定的顺序是相同的，即使编译器或 CPU 重排了指令
- Go 只保证 goroutine 内部重排对读写的顺序没有影响

#### 1. init 
- 程序初始化是在单一 goroutine 内执行的；如果 pkg A 导入了 pkg B，那么 B 的初始化一定在 A 之前执行
- **特殊情况**：main 函数一定在导入 pkg 的 init 函数之后执行
- 同一 pkg 下的多个文件，会按照文件名排序一次进行初始化
- pkg 级别的变量在同一个文件中是按照声明顺序逐个初始化的，除非初始化它的时候依赖其它的变量

#### 2. goroutine
- 启动 goroutine 的 go 语句执行，一定 happens before 此 goroutine 内代码的执行
  - 如果 go 语句传入的参数是一个函数执行的结果，那么该函数一定优先于 goroutine 内部代码执行
  
#### 3. channel 
- 往 channel 发送数据，一定 happens before 从该 channel 接受相应数据的**动作完成之前**；即第 n 个 send 一定 happens before 第 n 个 receive 的**完成**
- close 一个 Channel 的调用，肯定 happens before 从关闭的 Channel 中读取出一个零值
- **对于 unbuffered 的 Channel**，也就是容量是 0 的 Channel，从此 Channel 中**读取数据**的调用一定 happens before 往此 Channel **发送数据的调用完成**
- 如果 Channel 的容量是 m（m>0），那么，第 **n 个 receive** 一定 happens before 第 **n+m 个 send 的完成**

#### 4. Mutex/RWMutex
- 第 **n 次的 m.Unlock** 一定 happens before 第 **n+1 m.Lock 方法的返回**
- 对于读写锁 RWMutex m，如果它的第 n 个 m.Lock 方法的调用已返回，那么它的第 n 个 m.Unlock 的方法调用一定 happens before 任何一个 m.RLock 方法调用的返回，只要这些 m.RLock 方法调用 happens after 第 n 次 m.Lock 的调用的返回。这就可以保证，只有释放了持有的写锁，那些等待的读请求才能请求到读锁
- 对于读写锁 RWMutex m，如果它的第 n 个 m.RLock 方法的调用已返回，那么它的第 k （k<=n）个成功的 m.RUnlock 方法的返回一定 happens before 任意的 m.Lock 方法调用，只要这些 m.Lock 方法调用 happens after 第 n 次 m.RLock

#### 5. WaitGroup 
- wg.Add 或 wg.Done 一定 happens before wg.Wait 方法的返回
- wait 方法等到计数值归零之后才返回

#### 6. sync.Once
- 对于 once.Do(f) 调用，f 函数的单次调用一定 happens-before 任何 once.Do(f) 调用的返回
  - 函数 f 一定会在 Do 方法返回之前执行

#### 7. atomic
- 可以保证使用 atomic 的 Load/Store 的变量之间的顺序性，**但是过于复杂，现阶段不建议使用 atomic 保证顺序性**