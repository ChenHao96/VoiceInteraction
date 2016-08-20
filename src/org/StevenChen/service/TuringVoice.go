package service

import "com/tuling123"

const APT_KEY string = "d975f8141aa550cea27b7f48dd50c48d"

var tuLingAPI = tuling123.NewAPI_Request(APT_KEY)

func DoSomething(result []string) string {
	/**
		do something...
		添加你自己的回答
	 */
	return tuLingAPI.Talk(result[0]);
}