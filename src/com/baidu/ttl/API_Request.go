package ttl

//REST API Url
const API_URL = "http://tsn.baidu.com/text2audio"

type API_Request struct {
	Tex  string `json:"tex"`           //必填；合成文本，使用UTF-8编码，请注意文本长度必须小于1024
	Lan  string `json:"lan"`           //必填；语言选择，填写zh
	Tok  string `json:"tok"`           //必填；开放平台获取到的开发者access_token
	Ctp  int    `json:"ctp"`           //必填；客户端类型选择,web端填写1
	Cuid string `json:"cuid"`          //必填；用户唯一标识，用来区分用户，填写机器 MAC 地址或 IMEI 码，长度为60以内
	Spd  int    `json:"spd,omitempty"` //选填；语速，取值0-9，默认5
	Pit  int    `json:"pit,omitempty"` //选填；语调，取值0-9，默认5
	Vol  int    `json:"vol,omitempty"` //选填；音量，取值0-9，默认5
	Per  int    `json:"per,omitempty"` //选填；发音人选择，取值0-1；默认0女声 1男声
}