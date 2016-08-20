package tuling123

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
	"org/StevenChen/util"
)

const API_URL string = "http://www.tuling123.com/openapi/api"

type API_Request struct {
	Key    string `json:"key"`           //必须	32位	1ca8089********736b8ce41591426(32位)		注册之后在机器人接入页面获得（参见本文档第2部分）
	Info   string `json:"info"`          //必须	1-30位	打招呼“你好”;查天气“北京今天天气”;			请求内容，编码方式为UTF-8
	UserId string `json:"userid"`        //必须	1-32位	abc123（支持0-9，a-z,A-Z组合，不能包含特殊字符）	开发者给自己的用户分配的唯一标志（对应自己的每一个用户）
	Loc    string `json:"loc,omitempty"` //非必须	1-30位	北京市中关村					位置信息，请求跟地理位置相关的内容时使用，编码方式UTF-8
}

var API_Response_Err = map[int]string{
	40001:"参数key错误",
	40002:"请求内容info为空",
	40004:"当天请求次数已使用完",
	40007:"数据格式异常",
}

type API_Response struct {
	Code int`json:"code"`
	Text string `json:"text"`
}

func NewAPI_Request(key string) API_Request {

	var request API_Request
	request.Key = key
	request.UserId = strings.Join(strings.Split(util.GetCUID(), ":"), "")

	return request
}

func (this API_Request) Talk(worlds string) string {

	this.Info = worlds
	data, err := json.Marshal(this)
	if nil != err {
		panic(err.Error())
	}

	response, err := http.Post(API_URL, "application/json; charset=utf-8", bytes.NewReader(data));
	defer response.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		panic(err.Error())
	}

	var first API_Response
	json.Unmarshal(body, &first);

	if 100000 != first.Code {
		if content, ok := API_Response_Err[first.Code]; ok {
			panic(content)
		} else {
			panic("Unknow,该信息暂时无法识别")
		}
	}

	return first.Text
}