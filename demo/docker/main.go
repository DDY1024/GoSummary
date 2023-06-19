package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	// tpl := fasttemplate.New("{{name}} is da hao ren!", "{{", "}}")
	// fmt.Println(tpl.ExecuteString(map[string]interface{}{
	// 	"name": "wangxiyang",
	// }))

	// 	var str string
	// POINT:
	// 	for {
	// 		_, err := fmt.Scan(&str)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			os.Exit(1)
	// 		}

	// 		switch str {
	// 		case "end":
	// 			break POINT
	// 		default:
	// 			fmt.Println("Input:", str)
	// 		}
	// 	}
	// 	fmt.Println("End")

	http.HandleFunc("/ping", ping)
	fmt.Println("http start")
	if err := http.ListenAndServe(":18080", nil); err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, strings.NewReader("Hello!"))
}
