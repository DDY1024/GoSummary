package main

import "fmt"

// 适配器模式: 顾明思义就是做适配工作，将不兼容接口通过包装转化成兼容接口
// 使用目的: 对外统一接口形式，内部做各种各样的适配工作

// 通过 ICreateServer 接口形式统一云主机创建接口
type ICreateServer interface {
	CreateServer(cpu, mem float64) error
}

type AWSClient struct{}

func (c *AWSClient) RunInstance(cpu, mem float64) error {
	fmt.Printf("aws client run success, cpu： %f, mem: %f", cpu, mem)
	return nil
}

// AwsClientAdapter 适配: CreateServer 内部
type AwsClientAdapter struct {
	Client AWSClient
}

func (a *AwsClientAdapter) CreateServer(cpu, mem float64) error {
	a.Client.RunInstance(cpu, mem)
	return nil
}

type AliyunClient struct{}

func (c *AliyunClient) CreateServer(cpu, mem int) error {
	fmt.Printf("aws client run success, cpu： %d, mem: %d", cpu, mem)
	return nil
}

type AliyunClientAdapter struct {
	Client AliyunClient
}

// CreateServer 启动实例： CreateServer(int, int) --> CreateServer(float, float)
func (a *AliyunClientAdapter) CreateServer(cpu, mem float64) error {
	a.Client.CreateServer(int(cpu), int(mem))
	return nil
}
