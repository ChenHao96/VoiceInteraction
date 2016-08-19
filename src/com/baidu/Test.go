package main

import (
	"com/baidu/yuyin"
	"os"
	"fmt"
	"com/baidu/tts"
	"os/exec"
)

const (
	API_Key = "6NzYjugkwzipGIMbLRPKjCaQ"
	Secret_Key = "63950f36a6d2026ad8d9a7afbbb66895"
)

var resourcePath = os.Getenv("GOPATH") + "/resources/"

func main() {

	voice()

	text()

	c := exec.Command(resourcePath+"gplay.exe", resourcePath + "test.mp3")
	if err := c.Run(); "exit status 1" != fmt.Sprintf("%s", err) {
		fmt.Println("Error:", err)
	}
}

func text() {

	util := tts.NewAPI_Util(API_Key, Secret_Key)

	util.Text2AudioFile(resourcePath + "test.mp3", "百度语音提供技术支持，")

	fmt.Println("Text2AudioFile：百度语音提供技术支持，")
}

func voice() {

	filePath := resourcePath + "test.pcm"
	util := yuyin.NewAPI_Util(API_Key, Secret_Key)

	result := util.SendFileRequest(filePath, "pcm", 8000)
	fmt.Println("SendFileRequest:", result.Result)

	result = util.SendBytesRequest(filePath, "pcm", 8000)
	fmt.Println("SendBytesRequest:", result.Result)
}
