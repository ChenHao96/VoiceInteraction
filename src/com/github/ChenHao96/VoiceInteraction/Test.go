package main

import (
	"com/baidu/yuyin"
	"os"
	"fmt"
	"encoding/json"
)

const (
	API_Key = "************************"
	Secret_Key = "********************************"
)

func main() {

	filePath := os.Getenv("GOPATH") + "/resources/test.pcm"
	util, err := yuyin.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err)
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