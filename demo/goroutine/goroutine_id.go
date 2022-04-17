package goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// 利用 goroutine 堆栈信息获取 goroutine_id
// 具体参考: https://chai2010.cn/advanced-go-programming-book/ch3-asm/ch3-08-goroutine-id.html

/*
goroutine 1 [running]:
main.main()
    /path/to/main.g
*/
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
