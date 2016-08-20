package service

import "com/tuling123"

const (
	API_KEY string = "89b298ba45ec4a479dd9f20076d82b81"
)

var tuLingAPI = tuling123.NewAPI_Request(API_KEY)

func DoSomething(result []string) string {
	/**
		do something...
		添加你自己的回答
	 */
	return tuLingAPI.Talk(result[0]);
}