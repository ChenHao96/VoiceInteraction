package tuling123

const API_URL string = "http://www.tuling123.com/openapi/api"

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