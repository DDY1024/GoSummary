package main

import (
	"fmt"
	"sync"
)

// scanline: 读入完整的一行数据

// func main() {
// 	buf := bufio.NewScanner(os.Stdin)
// 	for {
// 		if !buf.Scan() {
// 			break
// 		}
// 		line := buf.Text()
// 		fmt.Println(line)
// 	}
// }

// byte scan
func Scanline() string {
	var (
		c   byte
		bs  []byte
		err error
	)
	for {
		_, err = fmt.Scanf("%c", &c)
		if err != nil || c == '\n' {
			// fmt.Println(err, c)
			break
		}
		bs = append(bs, c)
	}
	return string(bs)
}

func main() {
	s := Scanline()
	fmt.Println(s)
	var m sync.Map
	m.Load()
	m.Store()
}
