package main

import (
	"com/baidu/yuyin"
	"os"
	"fmt"
	"com/baidu/tts"
	"os/exec"
)

const (
	API_Key = "ZC2NNfFUkg8rxgmVkfBC6ycX"
	Secret_Key = "9a98e53b2ef7339bf03793f0b53fc7e4"
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
