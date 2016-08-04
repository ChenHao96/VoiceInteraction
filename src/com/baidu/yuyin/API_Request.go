package yuyin

//REST API Url
const API_URL = "http://vop.baidu.com/server_api"

type API_Request struct {
	format   string //必填	语音压缩的格式，请填写上述格式之一，不区分大小写
	rate     int    //必填	采样率，支持 8000 或者 16000
	channel  int    //必填	声道数，仅支持单声道，请填写 1
	cuid     string //必填	用户唯一标识，用来区分用户，填写机器 MAC 地址或 IMEI 码，长度为60以内
	token    string //必填	开放平台获取到的开发者 access_token
	ptc      string //选填	协议号，下行识别结果选择，默认 nbest 结果
	lan      string //选填	语种选择，中文=zh、粤语=ct、英文=en，不区分大小写，默认中文
	url      string //选填	语音下载地址
	callback string //选填	识别结果回调地址
	speech   string //选填	真实的语音数据 ，需要进行base64 编码
	len      int    //选填	原始语音长度，单位字节
}