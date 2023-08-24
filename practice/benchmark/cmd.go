package main

/*
#### 基准测试
- `go test -bench=. -run=none`
    - `-benchtime=3s`：指定性能测试时间
    - `-benchmem`：内存统计，包括每次操作内存分配次数，每次操作分配内存字节数
    - `-cpu=1,2,4`：指定基准测试时 `GOMAXPROCS`数量
    - `-count=10`：基础测试执行多次，针对操作耗时非常小（微秒、纳秒）的场景
- 更好的工具 `benchstat`

#### 其它命令
- `go test -v -run={func_name}`：指定运行函数
- `go test -v -coverprofile=c.out`：生成代码测试覆盖率结果
- `go tool cover -html=c.out -o=tag.html`：显示覆盖率测试结果
*/
