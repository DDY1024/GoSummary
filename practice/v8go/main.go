package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"rogchap.com/v8go"
	v8 "rogchap.com/v8go"
)

// 直接运行
func directCall() {
	ctx := v8.NewContext()                                // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript("const add = (a, b) => a + b", "x1.js") // executes a script on the global context
	ctx.RunScript("const result = add(3, 4)", "x2.js")    // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "x3.js")            // return a value in JavaScript back to Go
	fmt.Println("addition result:", val)
}

// v8 vm: context、vm、v8go
func vmRun() {
	iso := v8.NewIsolate() // creates a new JavaScript VM
	// ctx1 := v8.NewContext(iso) // new context within the VM
	// v1, _ := ctx1.RunScript("const multiply = (a, b) => a * b", "math.js")
	// v2, _ := ctx1.RunScript("result = multiply(1, 2)", "xx.js")

	// ctx 具有串联上下文的能力
	ctx2 := v8.NewContext(iso) // another context on the same VM
	if _, err := ctx2.RunScript("multiply(3, 4)", "main.js"); err != nil {
		// this will error as multiply is not defined in this context
		fmt.Println(err)
	}
}

func funcTplRun() {
	iso := v8.NewIsolate() // create a new VM

	// a template that represents a JS function
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%d\n", len(info.Args()))
		fmt.Printf("%v\n", info.Args()) // when the JS function is called this Go callback will execute
		return nil                      // you can return a value back to the JS caller if required
	})

	global := v8.NewObjectTemplate(iso) // a template that represents a JS Object (create a global object in a Context)
	global.Set("print", printfn)        // sets the "print" property of the Object to our function

	ctx := v8.NewContext(iso, global)                    // new Context with the global Object set to our object template
	fmt.Println(ctx.RunScript("print('foo')", "xxx.js")) // will execute the Go callback with a single argunent 'foo'
}

func globalObjRun() {
	ctx := v8.NewContext() // new context with a default VM

	obj := ctx.Global()          // get the global object from the context
	obj.Set("version", "v1.0.0") // set the property "version" on the object

	val, _ := ctx.RunScript("version", "version.js") // global object will have the property set within the JS VM
	fmt.Printf("version: %s\n", val)

	if obj.Has("version") { // check if a property exists on the object
		obj.Delete("version") // remove the property from the object
	}
}

// iso、global、ctx
// iso --> ctx : 一对多
// ctx --> global: 一对一 NewObjectTemplate
// global --> func: 一对多 NewFunctionTemplate

// 终止 js 脚本执行
func terminateRun() {
	data, _ := os.ReadFile("hello.js")
	// fmt.Println(string(data))

	iso := v8.NewIsolate()
	// iso.TerminateExecution(): 终止 js 脚本的执行
	ctx := v8.NewContext(iso)

	_, err := ctx.RunScript(string(data), "hello.js")
	if err != nil {
		fmt.Println(err)
	}
}

// // js 脚本预编译
// func main() {
// 	source := "const multiply = (a, b) => a * b"
// 	iso1 := v8.NewIsolate()                                                         // creates a new JavaScript VM
// 	ctx1 := v8.NewContext(iso1)                                                     // new context within the VM
// 	script1, _ := iso1.CompileUnboundScript(source, "math.js", v8.CompileOptions{}) // compile script to get cached data
// 	val, _ := script1.Run(ctx1)

// 	cachedData := script1.CreateCodeCache()

// 	iso2 := v8.NewIsolate()     // create a new JavaScript VM
// 	ctx2 := v8.NewContext(iso2) // new context within the VM

// 	script2, _ := iso2.CompileUnboundScript(source, "math.js", v8.CompileOptions{CachedData: cachedData}) // compile script in new isolate with cached data
// 	val, _ = script2.Run(ctx2)
// }

// func main() {

// 	testError()
// 	vm := v8.NewIsolate()

// 	f := v8.NewFunctionTemplate(vm, func(info *v8.FunctionCallbackInfo) *v8.Value {
// 		panic("wo ca, shen me qing kuang!")
// 	})

// 	obj := v8.NewObjectTemplate(vm)
// 	obj.Set("fuck", f)

// 	data, _ := os.ReadFile("hello.js")
// 	ctx := v8.NewContext(vm, obj)
// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				fmt.Println("Panic:", err)
// 			}
// 		}()
// 		if _, err := ctx.RunScript(string(data), "hello.js"); err != nil {
// 			fmt.Println("XXXX:", err)
// 		}
// 	}()

// 	time.Sleep(2 * time.Second)
// 	fmt.Println("success")

//		// var err error
//		// var err error
//	}

var (
	globalIso *v8.Isolate
	globalCtx *v8.Context
)

var script1 = `
function gcd(a, b) {
    if (b === 0) {
        return a
    } else {
        return gcd(b, a % b)
    }
}
`

func init() {
	globalIso = v8.NewIsolate()
	globalCtx = v8.NewContext(globalIso)

	// pre-compile
	ins, err := globalIso.CompileUnboundScript(script1, "gcd.js", v8.CompileOptions{})
	if err != nil {
		fmt.Println("XXXXX:", err)
	}
	if _, err := ins.Run(globalCtx); err != nil {
		fmt.Println("YYYYY:", err)
	}

	// obj1 := v8.NewObjectTemplate(globalIso)
	// f := v8.NewFunctionTemplate(globalIso, func(info *v8.FunctionCallbackInfo) *v8.Value {
	// 	a := info.Args()[0].Integer()
	// 	b := info.Args()[1].Integer()
	// 	c := gcd(a, b)
	// 	r, _ := v8.NewValue(globalIso, c)
	// 	return r
	// })
	// obj1.Set("gcd", f)

	// obj2 := v8.NewObjectTemplate(globalIso)
	// obj2.Set("math", obj1)
	// globalCtx = v8.NewContext(globalIso, obj2)
	// globalCtx.Global().Set("xxx", )
}

// func gcd(a, b int64) int64 {
// 	if b == 0 {
// 		return a
// 	}
// 	return gcd(b, a%b)
// }

// func jsonParse(value *v8go.Value, v interface{}) error {
// 	data, err := value.MarshalJSON()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("XXXX:", string(data))

// 	// v8go.NewValue()

// 	return json.Unmarshal(data, &v)
// }

func JsValue(ctx *v8go.Context, value interface{}) (*v8go.Value, error) {
	if value == nil {
		return v8go.Null(ctx.Isolate()), nil
	}

	switch v := value.(type) {
	case string, int32, uint32, bool, *big.Int, float64, []byte: // basic type
		return v8go.NewValue(ctx.Isolate(), v)

	case int64:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case uint64:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case int:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case int8:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case int16:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case uint:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case uint8:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case uint16:
		return v8go.NewValue(ctx.Isolate(), int32(v))

	case float32:
		return v8go.NewValue(ctx.Isolate(), float64(v))

	default:
		return jsValueParse(ctx, v)
	}
}

func jsValueParse(ctx *v8go.Context, value interface{}) (*v8go.Value, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	jsValue, err := v8go.JSONParse(ctx, string(data))
	if err != nil {
		return nil, err
	}
	return jsValue, nil
}

func main() {
	// a, _ := JsValue(globalCtx, 1)
	// b, _ := JsValue(globalCtx, 2)
	a, _ := v8go.NewValue(globalIso, int32(1)) // 注意: 此处输入 1 默认类型为 int，不行；必须显示指定 int32
	b, _ := v8go.NewValue(globalIso, int32(2))
	fmt.Println(globalCtx.Global().MethodCall("gcd", a, b))
	// c, _ := v8go.NewValue(globalCtx.Isolate(), 3)
	// args := []v8.Valuer{a, b}
	// v, _ := globalCtx.Global().Get("gcd")
	// fmt.Println(globalCtx.RunScript(script1, "xx.js"))
	// fmt.Println(v.IsFunction()) // true
	// f, _ := v.AsFunction()
	// f.Call(v8.Valuer(v8.Undefined(globalIso)), a, b)

	// globalCtx.Global().MethodCall("gcd", a, b)
	// a, _ := v8.NewValue(globalIso, 1)
	// b, _ := v8.NewValue(globalIso, 2)
	// args := []v8.Valuer{a, b}

	// fmt.Println(globalCtx.Global().MethodCall("gcd", args...))

	// val, err := v8go.NewValue(globalIso, []byte{1, 2, 3}) // []byte{1, 2, 3}
	// fmt.Println(val, err)
	// fmt.Println(val.IsUint8Array())
	// fmt.Println(val.Uint8Array())
	// // fmt.Println(val.IsArray())
	// // // data, _ := val.MarshalJSON()
	// // // fmt.Println(data)

	// var arr map[string]interface{}
	// fmt.Println(jsonParse(val, &arr))
	// fmt.Println(arr)

	// data, _ := json.Marshal([]int{1, 2, 3})
	// fmt.Println("XXX:", string(data))
	// r, err := v8go.JSONParse(globalCtx, string(data))
	// fmt.Println(r, err)
	// fmt.Println(r.IsArray())

	// // obj := globalCtx.Global()
	// vv, _ := globalCtx.Global().Get("math")
	// if vv.IsObject() {
	// 	// fmt.Println(val.AsObject())
	// 	obj, _ := vv.AsObject()
	// 	fmt.Println(obj.Get("gcd"))
	// }

	// source := `math.gcd(1,2);`
	// vvv, err := globalCtx.RunScript(source, "xx.js")
	// fmt.Println("XXXXX:", vvv, err)

	// source2 := `throw "wc";`
	// zzz, err := globalCtx.RunScript(source2, "yy.js")
	// fmt.Println("YYYYY:", zzz)
	// fmt.Println(err)

	// v, _ := globalCtx.Global().Get("Error")
	// fmt.Println("WC:", v)
	// fmt.Println(globalCtx.Global().Get("math"))

	// `try catch`

	// source := "a=10;"
	// wg := sync.WaitGroup{}
	// for i := 0; i < 100; i++ {
	// 	go func(i int) {
	// 		defer wg.Done()
	// 		val, err := globalCtx.RunScript(source, "test.json")
	// 		if err != nil {
	// 			fmt.Println("XXX:", i, err)
	// 			return
	// 		}
	// 		fmt.Println(val)
	// 		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	// 	}(i)
	// }
	// wg.Wait()
}

// func main() {
// 	// source := "const multiply = (a, b) => a * b"
// 	iso := v8.NewIsolate()

// 	mul := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
// 		fmt.Println("XXXX:", len(info.Args()))
// 		if len(info.Args()) == 2 {
// 			r, _ := v8.NewValue(iso, info.Args()[0].Integer()*info.Args()[1].Integer())
// 			// info.Args()[0] *
// 			return r
// 		}

// 		val, err := v8.NewValue(iso, 0)
// 		if err != nil {
// 			return v8.Undefined(iso)
// 		}
// 		return val
// 		// fmt.Printf("%d\n", len(info.Args()))
// 		// fmt.Printf("%v\n", info.Args()) // when the JS function is called this Go callback will execute
// 		// return nil                      // you can return a value back to the JS caller if required
// 	})
// 	obj := v8.NewObjectTemplate(iso)
// 	obj.Set("mul", mul)
// 	// obj.Set                                            // creates a new JavaScript VM

// 	source := "mul(2, 3)"
// 	ctx1 := v8.NewContext(iso, obj)
// 	start := time.Now()                                                            // new context within the VM
// 	script1, _ := iso.CompileUnboundScript(source, "math.js", v8.CompileOptions{}) // compile script to get cached data
// 	fmt.Println("XXXX:", time.Since(start).Milliseconds())

// 	time.Sleep(2 * time.Second)
// 	val, _ := script1.Run(ctx1)
// 	fmt.Println(val)
// }
