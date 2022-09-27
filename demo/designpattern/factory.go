package go_design_pattern

import "fmt"

// 三类工厂模式
// 1. 简单工厂
// 2. 工厂方法
// 3. 抽象工厂

// 1. 简单工厂实现
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

// 2. 工厂方法
// 相当于创建相应类型 parser 的方法下沉到子类中实现，而不是一个大而全的 NewXXX 方法

type ParserCreator interface {
	Create() Parser
}

type jsonParserCreator struct{}

func (jpc jsonParserCreator) Create() Parser {
	return jsonParser{}
}

type yamlParserCreator struct{}

func (ypc yamlParserCreator) Create() Parser {
	return yamlParser{}
}

// 3. 抽象工厂
type AbstractFactory interface { // 创建多种不同类型的产品
	CreateTelevision() ITelevision
	CreateAirConditioner() IAirConditioner
}

type ITelevision interface {
	Watch()
}

type IAirConditioner interface {
	SetTemperature(int)
}

type HuaWeiFactory struct{}

func (hf *HuaWeiFactory) CreateTelevision() ITelevision {
	return &HuaWeiTV{}
}
func (hf *HuaWeiFactory) CreateAirConditioner() IAirConditioner {
	return &HuaWeiAirConditioner{}
}

type HuaWeiTV struct{}

func (ht *HuaWeiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type HuaWeiAirConditioner struct{}

func (ha *HuaWeiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("HuaWei AirConditioner set temperature to %d ℃\n", temp)
}

type MiFactory struct{}

func (mf *MiFactory) CreateTelevision() ITelevision {
	return &MiTV{}
}
func (mf *MiFactory) CreateAirConditioner() IAirConditioner {
	return &MiAirConditioner{}
}

type MiTV struct{}

func (mt *MiTV) Watch() {
	fmt.Println("Watch HuaWei TV")
}

type MiAirConditioner struct{}

func (ma *MiAirConditioner) SetTemperature(temp int) {
	fmt.Printf("Mi AirConditioner set temperature to %d ℃\n", temp)
}

func main() {
	var factory AbstractFactory
	var tv ITelevision
	var air IAirConditioner

	// 相应具体实例的创建，下沉到子类中去实现
	factory = &HuaWeiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(25)

	factory = &MiFactory{}
	tv = factory.CreateTelevision()
	air = factory.CreateAirConditioner()
	tv.Watch()
	air.SetTemperature(26)
}
