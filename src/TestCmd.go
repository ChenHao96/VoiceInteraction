package main

import (
	"os/exec"
	"os"
	"fmt"
	"bufio"
	"strconv"
	"com/baidu/yuyin"
	"com/baidu/tts"
)

const (
	API_Key = "************************"
	Secret_Key = "********************************"
)

var resourcePath = os.Getenv("GOPATH") + "/resources/"

/*
该测试在linux下运行
 */
func main() {

	suffix := "wav"
	fileName := "test." + suffix
	rate := 16000

	/*
		打开linux下的录音软件
		-r 采样率
		-t 文件格式
		-c 声道
		-f 位深
		-d 录音时间(单位；秒)
		arecord的操作可以查看：http://blog.chinaunix.net/uid-29616823-id-4761787.html
	 */
	cmd := exec.Command("arecord", "-r", strconv.Itoa(rate), "-t", suffix, "-c", "1", "-f", "S16_LE", resourcePath + fileName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	go func() {
		fmt.Println("开始聆听...")
		if err := cmd.Run(); nil != err {
			fmt.Println(err)
		}
	}()

	readLine := bufio.NewReader(os.Stdin)
	for {
		if key, err := readLine.ReadByte(); nil != err {
			fmt.Println(err)
		} else {
			if byte('e') == key || byte('q') == key {
				break
			}
		}
	}
	fmt.Println("结束聆听。")
	cmd.Process.Kill()

	util := yuyin.NewAPI_Util(API_Key, Secret_Key)
	result := util.SendFileRequest(resourcePath + fileName, suffix, rate)

	resultText := result.Result[0]
	fmt.Println("SendFileRequest:", resultText)

	ttsUtil := tts.NewAPI_Util(API_Key, Secret_Key)
	ttsUtil.Text2AudioFile(resourcePath + "test.mp3", resultText)

	/*
		linux 下安装下面的库
		sudo apt-get install sox libsox-fmt-all
	 */
	cmd = exec.Command("play", resourcePath + "test.mp3")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); nil != err {
		fmt.Println(err)
	}
}
