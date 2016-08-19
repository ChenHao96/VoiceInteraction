package service

import "com/tuling123"

const APT_KEY string = "********************************"

func DoSomething(result []string) string {
	/**
		do something...
		添加你自己的回答
	 */
	return tuling123.NewAPI_Request(APT_KEY).Talk(result[0]);
}