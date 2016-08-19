package main

import (
	"bufio"
	"com/baidu/tts"
	"com/baidu/yuyin"
	"fmt"
	"os"
	"os/exec"
	"org/StevenChen/service"
)

const (
	API_Key = "************************"
	Secret_Key = "********************************"
)

/*
该测试在linux下运行
*/
func main() {

	resourcePath := os.Getenv("GOPATH") + "/resources/"

	ttsUtil := tts.NewAPI_Util(API_Key, Secret_Key)
	voiceUtil := yuyin.API_Util{Credentials:ttsUtil.Credentials,Cuid:ttsUtil.Cuid}

	playFile := resourcePath + "test.mp3"
	recordFile := resourcePath + "test.wav"

	/*
		打开linux下的录音软件
		-r 采样率
		-t 文件格式
		-c 声道
		-f 位深
		-d 录音时间(单位；秒)
		arecord的操作可以查看：http://blog.chinaunix.net/uid-29616823-id-4761787.html
	*/
	record := exec.Command("arecord", "-r", "16000", "-t", "wav", "-c", "1", "-f", "S16_LE", recordFile)
	record.Stderr = os.Stderr
	record.Stdout = os.Stdout

	/*
		linux 下安装下面的库
		sudo apt-get install sox libsox-fmt-all
	*/
	play := exec.Command("play", playFile)
	play.Stderr = os.Stderr
	play.Stdout = os.Stdout

	for {
		go func() {
			fmt.Println("开始聆听...")
			if err := record.Run(); nil != err {
				fmt.Println(err)
			}
		}()

		/**********************************************************/
		/*********************此处还需要改进************************/
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
		/**********************************************************/
		record.Process.Kill()
		fmt.Println("结束聆听。")

		fmt.Println("正在识别...")
		result := voiceUtil.SendFileRequest(recordFile, "wav", 16000)

		if len(result.Result) <= 0 {
			fmt.Println("未识别到语音内容，识别无结果。")
			continue
		}
		fmt.Println("语音识别结果:", result.Result)

		answer := service.DoSomething(result.Result)
		fmt.Println("响应结果(回答):", answer)

		ttsUtil.Text2AudioFile(playFile, answer)
		if err := play.Run(); nil != err {
			fmt.Println(err)
		}
	}
}
