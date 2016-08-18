package main

import (
	"com/baidu/yuyin"
	"os"
	"fmt"
	"com/baidu/ttl"
	"os/exec"
)

const (
	//API_Key = "************************"
	//Secret_Key = "********************************"
	API_Key = "6NzYjugkwzipGIMbLRPKjCaQ"
	Secret_Key = "63950f36a6d2026ad8d9a7afbbb66895"
)

var resourcePath = os.Getenv("GOPATH") + "/resources/"

func main() {

	voice()

	text()

	//window下将gplay加到system32下就能听到声音
	c := exec.Command("gplay", resourcePath + "test.mp3")
	if err := c.Run(); "exit status 1" != fmt.Sprintf("%s", err) {
		fmt.Println("Error:", err)
	}
}

func text() {

	util, err := ttl.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err.Error())
	}

	err = util.Text2AudioFile(resourcePath + "test.mp3", "百度语音提供技术支持，")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("百度语音提供技术支持，")
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
	fmt.Println("SendFileRequest:", result.Result)

	result, err = util.SendBytesRequest(filePath, "pcm", 8000)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("SendBytesRequest:", result.Result)
}
