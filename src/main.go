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
	"regexp"
)

const (
	API_Key = "ZC2NNfFUkg8rxgmVkfBC6ycX"
	Secret_Key = "9a98e53b2ef7339bf03793f0b53fc7e4"

	playCommand = "play"
	recordCommand = "arecord"
)

func checkFile(resourcePath, playFile, recordFile string) (string, string) {

	message := "程序自动生成的音频文件在本地磁盘出现同名文件!!!"

	checkFile:
	_, err := os.Stat(playFile);
	_, err2 := os.Stat(recordFile);

	if (nil == err || os.IsExist(err)) || (nil == err2 || os.IsExist(err2)) {
		fmt.Println(message)
		fmt.Println("是否更改生成文件名称?直接回车忽略,或输入文件名(非法字符将被剔除)回车.")
		fmt.Print("文件名:")

		line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		if len(line) > 0 {
			fileName := regexp.MustCompile("\\W").ReplaceAllString(string(line), "")

			if len([]byte(fileName)) > 0 {
				playFile = resourcePath + "/" + fileName + ".mp3"
				recordFile = resourcePath + "/" + fileName + ".wav"
				fmt.Println("播放文件:", playFile, "\n 录音文件:", recordFile)
				message = "修改后的文件名称，文件依旧存在!!!"

				goto checkFile
			} else {
				fmt.Println("输入的字符均为非法字符，系统将强制覆盖同名文件!!!")
			}
		} else {
			playFile = resourcePath + "/test.mp3"
			recordFile = resourcePath + "/test.wav"
			fmt.Println("播放文件:", playFile, "\n录音文件:", recordFile)
		}
	}

	return playFile, recordFile
}

func funClose(playFile, recordFile string) {

	fmt.Println("正在退出.")

	_, err := os.Stat(playFile)
	if nil == err || os.IsExist(err) {
		if err := os.Remove(playFile); nil != err {
			fmt.Println(err)
		} else {
			fmt.Print(".")
		}
	}

	_, err = os.Stat(recordFile)
	if nil == err || os.IsExist(err) {
		if err := os.Remove(recordFile); nil != err {
			fmt.Println(err)
		} else {
			fmt.Print(".")
		}
	}

	syscall.Exit(0)
}

func recordSound(voiceUtil yuyin.API_Util, recordFile string) []string {

	/*
		打开linux下的录音软件
		-r 采样率
		-t 文件格式
		-c 声道
		-f 位深
		-d 录音时间(单位；秒)
		arecord的操作可以查看：http://blog.chinaunix.net/uid-29616823-id-4761787.html
	*/
	record := exec.Command(recordCommand, "-r", "16000", "-t", "wav", "-c", "1", "-f", "S16_LE", recordFile)
	go func() {
		fmt.Println("->回车结束录音")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("开始聆听...")
		if err := record.Run(); nil != err {
			fmt.Println(err)
		}
	}()

	/**********************************************************/
	/*********************此处还需要改进************************/
	bufio.NewReader(os.Stdin).ReadByte()
	/**********************************************************/
	time.Sleep(200 * time.Millisecond)
	record.Process.Kill()
	fmt.Println("结束聆听。")

	fmt.Println("正在识别...")
	result := voiceUtil.SendFileRequest(recordFile, "wav", 16000)

	return result.Result
}

func playAnswer(ttsUtil tts.API_Util, playFile, answer string) {

	/*
		linux下安装下面的库
		sudo apt-get install sox libsox-fmt-all
	*/
	ttsUtil.Text2AudioFile(playFile, answer)
	play := exec.Command(playCommand, playFile)
	if err := play.Run(); nil != err {
		fmt.Println(err)
	}
}

/*
该测试在linux下运行
*/
func main() {

	resourcePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	playFile := resourcePath + "/test.mp3"
	recordFile := resourcePath + "/test.wav"

	playFile, recordFile = checkFile(resourcePath, playFile, recordFile)

	ttsUtil := tts.NewAPI_Util(API_Key, Secret_Key)
	voiceUtil := yuyin.NewAPI_Util(API_Key, Secret_Key)

	for {

		result := recordSound(voiceUtil, recordFile)

		if len(result) <= 0 {
			fmt.Println("未识别到语音内容，识别无结果。")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		fmt.Println("语音识别结果:", result)

		beExit := false
		var answer string
		if strings.Contains(result[0], "闭嘴") {
			answer = "虽然你这么说我会很伤心，但我还是要走的，再见。"
			beExit = true
		} else {
			answer = service.DoSomething(result)
		}
		fmt.Println("响应结果(回答):", answer)

		playAnswer(ttsUtil, playFile, answer)

		if beExit {
			funClose(playFile, recordFile)
		}
	}
}