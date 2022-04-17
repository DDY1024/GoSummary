package go_design_pattern

import "fmt"

// 抽象工厂接口定义
type Parser interface {
	Parse(data []byte)
}

// 具体实现
type jsonParser struct{}

func (jp jsonParser) Parse(data []byte) {
	fmt.Println("json parser")
}

type yamlParser struct{}

func (yp yamlParser) Parse(data []byte) {
	fmt.Println("yaml parser")
}

// 提供一个 new 方法来创建指定的实现
func NewParser(option string) Parser {
	switch option {
	case "json":
		return jsonParser{}
	case "yaml":
		return yamlParser{}
	}
	return nil
}
