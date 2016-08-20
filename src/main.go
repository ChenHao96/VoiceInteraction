package main

import (
	"bufio"
	"com/baidu/tts"
	"com/baidu/yuyin"
	"fmt"
	"os"
	"os/exec"
	"org/StevenChen/service"
	"time"
	"strings"
	"syscall"
	"path/filepath"
)

const (
	API_Key = "ZC2NNfFUkg8rxgmVkfBC6ycX"
	Secret_Key = "9a98e53b2ef7339bf03793f0b53fc7e4"
)

/*
该测试在linux下运行
*/
func main() {

	resourcePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	ttsUtil := tts.NewAPI_Util(API_Key, Secret_Key)
	voiceUtil := yuyin.API_Util{Credentials:ttsUtil.Credentials, Cuid:ttsUtil.Cuid}

	playFile := resourcePath + "/test.mp3"
	recordFile := resourcePath + "/test.wav"

	for {
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
		go func() {
			fmt.Println("开始聆听...")
			time.Sleep(200 * time.Millisecond)
			if err := record.Run(); nil != err {
				fmt.Println(err)
			}
		}()

		/**********************************************************/
		/*********************此处还需要改进************************/
		bufio.NewReader(os.Stdin).ReadByte()
		/**********************************************************/
		record.Process.Kill()
		fmt.Println("结束聆听。")

		fmt.Println("正在识别...")
		result := voiceUtil.SendFileRequest(recordFile, "wav", 16000)

		if len(result.Result) <= 0 {
			fmt.Println("未识别到语音内容，识别无结果。")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		fmt.Println("语音识别结果:", result.Result)

		beExit := false
		answer := service.DoSomething(result.Result)
		if strings.Contains(answer, "闭嘴") {
			answer = "虽然你这么说我会很伤心，但我还是要走的，再见。"
		}
		fmt.Println("响应结果(回答):", answer)

		/*
			linux下安装下面的库
			sudo apt-get install sox libsox-fmt-all
		*/
		ttsUtil.Text2AudioFile(playFile, answer)
		play := exec.Command("play", playFile)
		if err := play.Run(); nil != err {
			fmt.Println(err)
		}

		if beExit {
			fmt.Println("正在退出...")
			syscall.Exit(0)
		}
	}
}