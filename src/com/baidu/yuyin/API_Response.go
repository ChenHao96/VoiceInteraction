package yuyin

const (
	PARAM_ERR = &API_ResponseErr{Err_code:3300, Meaning:"输入参数不正确"}
	RECOGNITION_ERR = &API_ResponseErr{Err_code:3301, Meaning:"识别错误"}
	VALIDATOR_ERR = &API_ResponseErr{Err_code:3302, Meaning:"验证失败"}
	SERVER_ERR = &API_ResponseErr{Err_code:3303, Meaning:"语音服务器后端问题"}
	TOO_BIG_REQUEST_ERR = &API_ResponseErr{Err_code:3304, Meaning:"请求 GPS 过大，超过限额"}
	REQUEST_EXCEED_ERR = &API_ResponseErr{Err_code:3305, Meaning:"产品线当前日请求数超过限额"}
)

type API_ResponseErr struct {
	Err_code int    //错误码
	Meaning  string //含义
}

type API_Response struct {
	err_no  int      //必填	错误码
	err_msg string   //必填	错误码描述
	sn      string   //必填	语音数据唯一标识，系统内部产生，用于 debug
	result  []string //选填	识别结果数组，提供1-5 个候选结果，string 类型为识别的字符串， utf-8 编码
}