package main

import (
	"os/exec"
	"os"
	"fmt"
	"bufio"
	"com/baidu/yuyin"
	"com/baidu/ttl"
)

const (
	API_Key = "************************"
	Secret_Key = "********************************"
)

func main() {
	fileName := "test.wav"

	cmd := exec.Command("arecord", "-r", "16000", "-t", "wav", "-c", "1", "-f", "S16_LE", fileName)
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

	util, err := yuyin.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err.Error())
	}

	result, err := util.SendFileRequest(fileName, "wav", 16000)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("SendFileRequest:", result.Result[0])

	ttlUtil, err := ttl.NewAPI_Util(API_Key, Secret_Key)
	if err != nil {
		panic(err.Error())
	}

	ttlUtil.Text2AudioFile(fileName,result.Result[0])

	cmd = exec.Command("aplay", fileName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); nil != err {
		fmt.Println(err)
	}
}
