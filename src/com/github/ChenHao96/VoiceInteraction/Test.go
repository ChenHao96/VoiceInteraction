package main

import (
	"fmt"
	"os"
	"syscall"
	"com/baidu/yuyin"
)

const (
	API_Key = "************************"
	Secret_Key = "********************************"
)

func main() {

	filePath := os.Getenv("GOPATH") + "/resources/test.pcm"

	token, err := yuyin.GetToken(API_Key, Secret_Key);
	if err != nil {
		fmt.Println("token获取异常:", err)
		syscall.Exit(-1)
	}

	cuid, err := yuyin.GetCUID()
	if err != nil {
		fmt.Println("cuid获取异常:", err)
		syscall.Exit(-1)
	}

	result, err := yuyin.SendBytesRequest(filePath, token, cuid)
	if err != nil {
		fmt.Println("byteRequest获取异常:", err)
		syscall.Exit(-1)
	}
	fmt.Println("byteRequest:",result)

	result, err = yuyin.SendFileRequest(filePath, token, cuid)
	if err != nil {
		fmt.Println("fileRequest获取异常:", err)
		syscall.Exit(-1)
	}
	fmt.Println("fileRequest:",result)
}