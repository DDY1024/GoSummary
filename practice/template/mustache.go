package main

import (
	"fmt"

	"github.com/aymerick/raymond"
)

var tpl = `<div class="entry">
<h1>{{title}}</h1>
<div class="body">
  {{body}}
</div>
<h1>{{person.name}}</h1>
</div>
`

var tpl2 = `wwwww`

func main() {
	ctx := map[string]interface{}{
		"title":  "My New Post",
		"body":   "This is my first post!",
		"person": map[string]interface{}{"name": "wxy"},
	}

	result, err := raymond.Render(tpl2, ctx)
	if err != nil {
		panic("Please report a bug :)")
	}

	fmt.Println(result)

	r, err := raymond.ParseFile("a.tpl")
	if err != nil {
		panic(err)
	}

	fmt.Println(r.Exec(ctx))

	r2, err := raymond.ParseFile("b.docx")
	if err != nil {
		panic(err)
	}

	fmt.Println(r2.Exec(ctx))
}

// template.xxx.render

// template.xxx.render
