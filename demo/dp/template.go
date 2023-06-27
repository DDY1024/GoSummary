package main

import (
	"fmt"
)

// 1. 通过 interface 定义一组标准操作流程（一组方法）
// 2. 通过定义不同的接口实现来达到不同场景下的特定执行效果
// 		往往会定义一个 base 实现，然后针对特殊场景改写 base 对应的实现方法

type ISMS interface {
	send(content string, phone int) error
}

type sms struct {
	ISMS
}

func (s *sms) Valid(content string) error {
	if len(content) > 63 {
		return fmt.Errorf("content is too long")
	}
	return nil
}

func (s *sms) Send(content string, phone int) error {
	if err := s.Valid(content); err != nil {
		return err
	}

	return s.send(content, phone)
}

type TelecomSms struct {
}

func NewTelecomSms() *TelecomSms {
	tel := &TelecomSms{}
	return tel
}

func (tel *TelecomSms) send(content string, phone int) error {
	fmt.Println("send by telecom success", content, phone)
	return nil
}

func main() {
	tsms := NewTelecomSms()
	sms := &sms{
		ISMS: tsms,
	}
	sms.Send("wxy", 10086)
}
