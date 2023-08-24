package main

import (
	"fmt"
	"mime"
	"path/filepath"
)

func main() {
	fmt.Println(mime.TypeByExtension(""))
	fmt.Println(mime.TypeByExtension(".png"))
	fmt.Println(mime.TypeByExtension(".txt"))

	fmt.Println(filepath.Ext("aa.txt"))
	fmt.Println(filepath.Ext("a/b.png"))
}
