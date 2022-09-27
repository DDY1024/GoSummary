package main

import (
	"fmt"
	"strings"
)

func reorderSpaces(s string) (ans string) {
	words := strings.Fields(s) // 切分单词，一个或多个空格
	space := strings.Count(s, " ")

	lw := len(words) - 1

	if lw == 0 {
		return words[0] + strings.Repeat(" ", space) // 重复
	}

	return strings.Join(words, strings.Repeat(" ", space/lw)) + strings.Repeat(" ", space%lw)
}

func main() {
	fmt.Println("hello, world!")
	fmt.Println(strings.Repeat("wxy", 10))
}
