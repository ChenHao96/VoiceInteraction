package yuyin

//REST API Url
const API_URL = "http://vop.baidu.com/server_api"

type API_Request struct {
	Format   string `json:"format"`             //必填	语音压缩的格式，请填写上述格式之一，不区分大小写
	Rate     int    `json:"rate"`               //必填	采样率，支持 8000 或者 16000
	Channel  int    `json:"channel"`            //必填	声道数，仅支持单声道，请填写 1
	Cuid     string `json:"cuid"`               //必填	用户唯一标识，用来区分用户，填写机器 MAC 地址或 IMEI 码，长度为60以内
	Token    string `json:"token"`              //必填	开放平台获取到的开发者 access_token
	Ptc      string `json:"ptc,omitempty"`      //选填	协议号，下行识别结果选择，默认 nbest 结果
	Lan      string `json:"lan,omitempty"`      //选填	语种选择，中文=zh、粤语=ct、英文=en，不区分大小写，默认中文
	Url      string `json:"url,omitempty"`      //选填	语音下载地址
	Callback string `json:"callback,omitempty"` //选填	识别结果回调地址
	Speech   string `json:"speech,omitempty"`   //选填	真实的语音数据 ，需要进行base64 编码
	Len      int    `json:"len,omitempty"`      //选填	原始语音长度，单位字节
}