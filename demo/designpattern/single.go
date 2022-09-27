package go_design_pattern

import "sync"

// 单例模式存在两种实现方式 "饿汉式" 和 "懒汉式"
// 1. "饿汉式" --> init 直接初始化创建对象
// 2. "懒汉式" --> 第一次访问时初始化创建

// "饿汉式" 实现: init 函数中直接初始化创建单个全局变量实例
type Singleton struct{}

var singleton *Singleton

func init() {
	singleton = &Singleton{}
}

func GetInstance() *Singleton {
	return singleton
}

// "懒汉式"：利用 go 中 sync.Once 特性，实现懒创建，且保证该实例只会被创建一次
// "懒汉式" 不适应场景: 对第一次创建操作比较耗时的场景，懒创建会导致这次处理耗时非常大，在一些场景下是不能接受的
var (
	lazySingleton *Singleton
	once          sync.Once
)

func GetLazySingleton() *Singleton {
	once.Do(func() {
		lazySingleton = &Singleton{}
	})
	return lazySingleton
}
