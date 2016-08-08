package yuyin

var API_ResponseErrEnum map[int]API_responseErr

func init()  {
	API_ResponseErrEnum = map[int]API_responseErr{
		0:{Err_code:0, Meaning:"成功"},
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
	Corpus_no string   `json:"corpus_no,omitempty"`//这个参数在官方的文档上我没有发现
	Err_no    int      `json:"err_no"`//错误码
	Err_msg   string   `json:"err_msg"`//错误码描述
	Sn        string   `json:"sn"`//语音数据唯一标识，系统内部产生，用于 debug
	Result    []string `json:"result"`//识别结果数组，提供1-5 个候选结果，string 类型为识别的字符串， utf-8 编码
}