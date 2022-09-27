1. slice、array 越界判断处理

2. 指针非空判断（nil 判断）

3. 对于外部输入数据做合法性校验
    - 对于外部输入 slice 数据做长度校验，避免非法输入过长的数据

4. runtime.SetFinalizer 使用问题
    - 当一个对象被 GC 选中到回收内存之前，runtime.SetFinalizer 都不会执行，即使程序正常结束或者发生错误
    - 由指针构成的**循环引用**虽然能被 GC 正确处理，但是由于无法确定 Finalizer 依赖顺序，因此无法调用 SetFinalizer，无法取消对象关系，从进一步导致无法确定对象变成不可达状态，进一步导致内存泄漏
```go
// bad
func foo() {
	var a, b Data
	a.o = &b
	b.o = &a
	// 指针循环引⽤，SetFinalizer()⽆法正常调⽤
	runtime.SetFinalizer(&a, func(d *Data) {
		fmt.Printf("a %p final.\n", d)
	})
	runtime.SetFinalizer(&b, func(d *Data) {
		fmt.Printf("b %p final.\n", d)
	})
}
func main() {
	for { // 注意: for 循环会不断地创建 a 和 b 对象，且 a 和 b 无法被 GC 回收，从而导致内存泄漏
		foo()
		time.Sleep(time.Millisecond)
	}
}
``` 

5. 避免对同一 channel 执行两次 close 操作
    - close 两次同一个 channel 会引发 panic

6. 避免 goroutine 泄漏
    - 确保每个 goroutine 均存在退出条件（如果创建多个同类 goroutine，且不存在退出条件，会导致 goroutine 泄漏）

7. unsafe 操作慎用

8. slice 作为函数入参应该注意的问题
    - slice 共同引用同一底层数组，函数内的修改会影响到外部变量
    - slice 函数内部 append 操作可能会触发扩容操作，导致底层 data 数组发生改变，需要注意

9. 针对外部输入变量需要做严格的校验
    - 外部输入文件路径
    - 外部输入命令参数 

10. 敏感数据处理
    - 对称加密算法、非对称加密算法
    - 加密算法选择
        - 推荐：crypto/rsa、crypto/aes 等
        - 不推荐：crypto/des、crypto/md5、crypto/sha1、crypto/rc4 等
11. 正则表达式
    - regexp 库进行正则表达式匹配；regexp 保证了线性时间性能和优雅失败；对解析器、编译器、执行引擎均进行了内存限制
    - regexp 不支持以下正则表达式特性
        - 回溯引用
    - regexp.MatchString()

12. 输入参数校验
    - validator 组件针对入参做校验
