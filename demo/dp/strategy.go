package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type StorageStrategy interface {
	Save(name string, data []byte) error
}

// 策略模型: map 维护所有注册的具体实现，根据需要的具体类型选择对应的执行策略
// 典型模式: map[type]strategy
var strategys = map[string]StorageStrategy{
	"file":         &fileStorage{},
	"encrypt_file": &encryptFileStorage{},
}

func NewStorageStrategy(t string) (StorageStrategy, error) {
	s, ok := strategys[t]
	if !ok {
		return nil, fmt.Errorf("not found StorageStrategy: %s", t)
	}
	return s, nil
}

type fileStorage struct{}

func (s *fileStorage) Save(name string, data []byte) error {
	return ioutil.WriteFile(name, data, os.ModeAppend)
}

type encryptFileStorage struct{}

func (s *encryptFileStorage) Save(name string, data []byte) error {
	data, err := encrypt(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(name, data, os.ModeAppend)
}

func encrypt(data []byte) ([]byte, error) {
	return data, nil
}
