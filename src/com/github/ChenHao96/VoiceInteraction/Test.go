package main

import (
	"com/baidu/yuyin"
	"os"
	"fmt"
	"encoding/json"
	"com/baidu/ttl"
)

const (
	API_Key = "6NzYjugkwzipGIMbLRPKjCaQ"
	Secret_Key = "63950f36a6d2026ad8d9a7afbbb66895"
)

var resourcePath = os.Getenv("GOPATH") + "/resources/"

func main() {
	//voice()
	text()
}

func text() {

	util, err := ttl.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err.Error())
	}

	err = util.Text2AudioFile(resourcePath, "test.pm3", "你好吗")
	if err != nil {
		panic(err.Error())
	}
}

func voice() {

	filePath := resourcePath + "test.pcm"
	util, err := yuyin.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err.Error())
	}

	result, err := util.SendFileRequest(filePath, "pcm", 8000)
	if err != nil {
		panic(err.Error())
	}
	value, err := json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("SendFileRequest:", string(value))

	result, err = util.SendBytesRequest(filePath, "pcm", 8000)
	if err != nil {
		panic(err.Error())
	}
	value, err = json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("SendBytesRequest:", string(value))
}