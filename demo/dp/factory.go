package go_design_pattern

import "fmt"

// 常见模式
// 1. 简单工厂
// 2. 工厂方法
// 3. 抽象工厂
// 4. 依赖注入

// 接口定义
type Parser interface {
	Parse(data []byte)
}

type jsonParser struct{}

func (jp jsonParser) Parse(data []byte) {
	fmt.Println("json parser")
}

type yamlParser struct{}

func (yp yamlParser) Parse(data []byte) {
	fmt.Println("yaml parser")
}

// 工厂创建方法，返回 Parser 接口类型
func NewParser(tp string) Parser {
	switch tp {
	case "json":
		return jsonParser{}
	case "yaml":
		return yamlParser{}
	}
	return nil
}

// 2. 工厂方法模式
// IRuleConfigParser IRuleConfigParser
type IRuleConfigParser interface {
	Parse(data []byte)
}

// jsonRuleConfigParser jsonRuleConfigParser
type jsonRuleConfigParser struct {
}

// Parse Parse
func (J jsonRuleConfigParser) Parse(data []byte) {
	panic("implement me")
}

// yamlRuleConfigParser yamlRuleConfigParser
type yamlRuleConfigParser struct {
}

// Parse Parse
func (Y yamlRuleConfigParser) Parse(data []byte) {
	panic("implement me")
}

// IRuleConfigParserFactory 工厂方法接口
type IRuleConfigParserFactory interface {
	CreateParser() IRuleConfigParser
}

// yamlRuleConfigParserFactory yamlRuleConfigParser 的工厂类
type yamlRuleConfigParserFactory struct {
}

// CreateParser CreateParser
func (y yamlRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return yamlRuleConfigParser{}
}

// jsonRuleConfigParserFactory jsonRuleConfigParser 的工厂类
type jsonRuleConfigParserFactory struct {
}

// CreateParser CreateParser
func (j jsonRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return jsonRuleConfigParser{}
}

// NewIRuleConfigParserFactory 用一个简单工厂封装工厂方法
func NewIRuleConfigParserFactory(t string) IRuleConfigParserFactory {
	switch t {
	case "json":
		return jsonRuleConfigParserFactory{}
	case "yaml":
		return yamlRuleConfigParserFactory{}
	}
	return nil
}

// 3. 抽象工厂
// ISystemConfigParser ISystemConfigParser
type ISystemConfigParser interface {
	ParseSystem(data []byte)
}

// jsonSystemConfigParser jsonSystemConfigParser
type jsonSystemConfigParser struct{}

// Parse Parse
func (j jsonSystemConfigParser) ParseSystem(data []byte) {
	panic("implement me")
}

// IConfigParserFactory 工厂方法接口
type IConfigParserFactory interface { // 提供创建各种 parser 方法
	CreateRuleParser() IRuleConfigParser
	CreateSystemParser() ISystemConfigParser
}

type jsonConfigParserFactory struct{}

func (j jsonConfigParserFactory) CreateRuleParser() IRuleConfigParser {
	return jsonRuleConfigParser{}
}

func (j jsonConfigParserFactory) CreateSystemParser() ISystemConfigParser {
	return jsonSystemConfigParser{}
}

// 4. 对于复杂初始化依赖逻辑，推荐使用依赖注入框架
// 	  a. 使用反射实现的: https://github.com/uber-go/dig
//    b. 使用 generate 实现的: https://github.com/google/wire
