package service

import (
	"com/tuling123/normal"
)

//const (
//	SECRET string = "4145a1cb5f92901b"
//	API_KEY string = "d975f8141aa550cea27b7f48dd50c48d"
//)
//var tuLingAPI = secret.NewAPI_Request(API_KEY, SECRET)

var tuLingAPI = normal.NewAPI_Request("89b298ba45ec4a479dd9f20076d82b81")

func DoSomething(result []string) string {
	/**
		do something...
		添加你自己的回答
	 */
	return tuLingAPI.Talk(result[0]);
}