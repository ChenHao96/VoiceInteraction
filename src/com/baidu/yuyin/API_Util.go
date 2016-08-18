package yuyin

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"
	"bytes"
	"strconv"
	"errors"
	"com/baidu/public"
	"io"
)

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

var API_ResponseErrEnum map[int]API_responseErr

func init() {
	API_ResponseErrEnum = map[int]API_responseErr{
		3300:{Err_code:3300, Meaning:"输入参数不正确"},
		3301:{Err_code:3301, Meaning:"识别错误"},
		3302:{Err_code:3302, Meaning:"验证失败"},
		3303:{Err_code:3303, Meaning:"语音服务器后端问题"},
		3304:{Err_code:3304, Meaning:"请求 GPS 过大，超过限额"},
		3305:{Err_code:3305, Meaning:"产品线当前日请求数超过限额"},
	}
}

type API_responseErr struct {
	Err_code int    //错误码
	Meaning  string //含义
}

type API_Response struct {
	Corpus_no string   `json:"corpus_no,omitempty"` //这个参数在官方的文档上我没有发现
	Err_no    int      `json:"err_no"`              //错误码
	Err_msg   string   `json:"err_msg"`             //错误码描述
	Sn        string   `json:"sn"`                  //语音数据唯一标识，系统内部产生，用于 debug
	Result    []string `json:"result"`              //识别结果数组，提供1-5 个候选结果，string 类型为识别的字符串， utf-8 编码
}

type API_Util struct {
	Credentials public.Credentials_Response
	Cuid        string
}

func NewAPI_Util(api_key, secret_key string) (util API_Util, err error) {

	cuid, err := public.GetCUID()
	if err != nil {
		return
	}

	res, err := public.GetCredentials(public.Credentials_Request{
		Client_id:api_key, Client_secret:secret_key})

	util.Cuid = cuid
	util.Credentials = res

	return
}

func getResult(url, contentType string, data io.Reader) (result API_Response, err error) {

	response, err := http.Post(url, contentType, data);
	defer response.Body.Close()
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	var first = make(map[string]string)
	err = json.Unmarshal(body, &first);

	if value, ok := first["err_code"]; ok {
		code, _ := strconv.Atoi(value)
		errMean, ok := API_ResponseErrEnum[code]
		if ok {
			err = errors.New(errMean.Meaning)
			return
		}
	}

	err = json.Unmarshal(body, &result);
	return
}

func printBase64Binary(val []byte) string {

	base64Table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encoding := base64.NewEncoding(base64Table);

	return encoding.EncodeToString(val)
}

/*
 不太推荐使用效率很低
*/
func (this API_Util) SendBytesRequest(filePath, format string, rate int) (API_Response, error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return API_Response{}, err
	}

	soundStr := printBase64Binary(data)

	param := &API_Request{Speech:soundStr, Cuid:this.Cuid,
		Token:this.Credentials.Refresh_token}
	param.Rate = rate
	param.Channel = 1
	param.Len = len(data)
	param.Format = format

	postValue, err := json.Marshal(param)
	if err != nil {
		return API_Response{}, err
	}

	return getResult(API_URL, "application/json; charset=utf-8", bytes.NewReader(postValue))
}

func (this API_Util) SendFileRequest(filePath, format string, rate int) (API_Response, error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return API_Response{}, err
	}

	url := API_URL + "?cuid=" + this.Cuid + "&token=" + this.Credentials.Refresh_token
	contentType := "audio/" + format + "; rate=" + strconv.Itoa(rate)

	return getResult(url, contentType, bytes.NewReader(data))
}