package yuyin

import (
	"bytes"
	"com/baidu/public"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

//REST API Url
const API_URL = "http://vop.baidu.com/server_api"

var API_ResponseErrEnum = map[int]string{
	3300: "输入参数不正确",
	3301: "识别错误",
	3302: "验证失败",
	3303: "语音服务器后端问题",
	3304: "请求 GPS 过大，超过限额",
	3305: "产品线当前日请求数超过限额",
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
	api_key     string
	secret_key  string
}

func NewAPI_Util(api_key, secret_key string) API_Util {

	cuid := public.GetCUID()

	res := public.GetCredentials(public.Credentials_Request{
		Client_id: api_key, Client_secret: secret_key})

	var util API_Util
	util.Cuid = cuid
	util.Credentials = res
	util.api_key = api_key
	util.secret_key = secret_key

	return util
}

func (this *API_Util) getResult(url, contentType string, data io.Reader) API_Response {

	response, err := http.Post(url, contentType, data)
	defer response.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	var first = make(map[string]string)
	json.Unmarshal(body, &first)

	if value, ok := first["err_code"]; ok {
		code, _ := strconv.Atoi(value)
		if 3302 == code {
			*this = NewAPI_Util(this.api_key, this.secret_key)
			return this.getResult(url, contentType, data)
		} else if errMean, ok := API_ResponseErrEnum[code]; ok {
			panic(errMean)
		}
	}

	var result API_Response
	err = json.Unmarshal(body, &result)
	if nil != err {
		panic(err.Error())
	}

	return result
}

func (this API_Util) SendFileRequest(filePath, format string, rate int) API_Response {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	url := API_URL + "?cuid=" + this.Cuid + "&token=" + this.Credentials.Refresh_token
	contentType := "audio/" + format + "; rate=" + strconv.Itoa(rate)

	return this.getResult(url, contentType, bytes.NewReader(data))
}
