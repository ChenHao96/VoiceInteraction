package util

import (
	"fmt"
	"runtime"
	"net"
)

/*
	获取一个本地的MAC地址作为API的 cuid
*/
func GetCUID() string {

	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}

	var result string
	switch runtime.GOOS {
	case "windows":
		result = fmt.Sprintf("%s", interfaces[0].HardwareAddr)
	case "linux":
		result = fmt.Sprintf("%s", interfaces[1].HardwareAddr)
	default:
		result = "01:02:03:04:05:06"
	}

	return result
}
