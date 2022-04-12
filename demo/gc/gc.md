**参考**：https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap-26c2462549a2/

1. Heap Size 
- some useful(live memory)
- some garbage

2. Live Memory
-  refers to all the allocations that are currently `being referenced by the running application`; not garbage

3. GC cost
- `gc mark` >> `gc sweep`：mark 贡献了 gc 大部分耗时，现代操作系统中 gc sweep 是一个很快速的工作

4. GC mark
- 遍历所有当前应用所引用的对象，因此其跟 live memory 是呈正比关系。内存中的 `extra garbage` 对 gc 耗时是不会有太大影响的。alloct_object 影响 gc 频率，in_use_object 才是真正影响 gc 耗时的。
- `garbage --> less gc freq --> less gc mark --> less cpu spent`
  
5. GC pacer
- `pacer`：determine when to trigger the next GC cycle --> target heap size 
- `mark termination`：during the mark termination phase of the current GC cycle --> 2 * live_set
- `GOGC` --> adjust gc pacer

6. GC ballast
```go
func main() {

	// Create a large heap allocation of 10 GiB
    // 这 10GB 内存常驻在 heap 中，从而提升 GC 触发的条件，降低 GC 触发的频率
    // 其最终本质是提高 GC 门槛，降低 GC 频率来达到我们想要的效果
    // []byte 对 GC Mark 影响是最低的
    // 采用这种方式相比较于调整 GOGC 来说，其对内存的控制是比较精确的
	ballast := make([]byte, 10<<30)

	// Application execution continues
	// ...

    runtime.KeepAlive(ballast)
}
```

7. GC assists
- `辅助标记`：gcAssistAlloc
- It should now be clear that any goroutine that does work that involves allocating will incur the GCAssist penalty during a GC cycle.

```go
// runtime.makeslice --> runtime.mallocgc
someObject :=  make([]int, 5)
```

8. Summary
- We noticed our applications were doing a lot of GC work
- We deployed a `memory ballast`
- It reduced GC cycles by allowing the heap to grow larger（提高 gc trigger 门槛）
- API latency improved since the `Go GC delayed our work less with assists`
- The ballast allocation is mostly free because `it resides in virtual memory`
- Ballasts are easier to reason about than configuring a GOGC value
- Start with a small ballast and increase with testing（测试取值）