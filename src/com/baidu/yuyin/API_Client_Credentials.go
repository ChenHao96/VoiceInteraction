package yuyin

/*
	客户凭证授权
	获取调用API的token
 */
const Credentials_Url = "https://openapi.baidu.com/oauth/2.0/token"

var Credentials_ResponseErrEnum map[string]Credentials_ResponseErr

func init() {
	Credentials_ResponseErrEnum = map[string]Credentials_ResponseErr{
		"invalid_request":{Error:"invalid_request", Error_description:"invalid refresh token", Description:"请求缺少某个必需参数，包含一个不支持的参数或参数值，或者格式不正确。"},
		"invalid_client":{Error:"invalid_client", Error_description:"unknown client id", Description:"client_id或client_secret参数无效"},
		"invalid_grant":{Error:"invalid_grant", Error_description:"The provided authorization grant is revoked", Description:"提供的Access Grant是无效的、过期的或已撤销的，例如，Authorization Code无效(一个授权码只能使用一次)、Refresh Token无效、redirect_uri与获取Authorization Code时提供的不一致、Devie Code无效(一个设备授权码只能使用一次)等。"},
		"unauthorized_client":{Error:"unauthorized_client", Error_description:"The client is not authorized to use this authorization grant type", Description:"应用没有被授权，无法使用所指定的grant_type。"},
		"unsupported_grant_type":{Error:"unsupported_grant_type", Error_description:"The authorization grant type is not supported", Description:"“grant_type”百度OAuth2.0服务不支持该参数。"},
		"invalid_scope":{Error:"invalid_scope", Error_description:"The requested scope is exceeds the scope granted by the resource owner", Description:"请求的“scope”参数是无效的、未知的、格式不正确的、或所请求的权限范围超过了数据拥有者所授予的权限范围。"},
		"expired_token":{Error:"expired_token", Error_description:"refresh token has been used", Description:"提供的Refresh Token已过期"},
		"redirect_uri_mismatch":{Error:"redirect_uri_mismatch", Error_description:"Invalid redirect uri", Description:"“redirect_uri”所在的根域与开发者注册应用时所填写的根域名不匹配。"},
		"unsupported_response_type":{Error:"unsupported_response_type", Error_description:"The response type is not supported", Description:"“response_type”参数值不为百度OAuth2.0服务所支持，或者应用已经主动禁用了对应的授权模式"},
		"slow_down":{Error:"slow_down", Error_description:"The device is polling too frequently", Description:"Device Flow中，设备通过Device Code换取Access Token的接口过于频繁，两次尝试的间隔应大于5秒。"},
		"authorization_pending":{Error:"authorization_pending", Error_description:"User has not yet completed the authorization", Description:"Device Flow中，用户还没有对Device Code完成授权操作。"},
		"authorization_declined":{Error:"authorization_declined", Error_description:"User has declined the authorization", Description:"Device Flow中，用户拒绝了对Device Code的授权操作。"},
		"invalid_referer":{Error:"invalid_referer", Error_description:"Invalid Referer", Description:"Implicit Grant模式中，浏览器请求的Referer与根域名绑定不匹配"},
	}
}

type Credentials_Request struct {
	Grant_type    string //必填参数 固定为“client_credentials”；
	Client_id     string //必填参数 应用的API Key
	Client_secret string //必须参数 应用的Secret Key;
	/*
		非必须参数。
		以空格分隔的权限列表，采用本方式获取Access Token时只能申请跟用户数据无关的数据访问权限。
		关于权限的具体信息请参考
		http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	*/
	Scope         string
}

type Credentials_Response struct {
	Access_token   string `json:"access_token"`   //要获取的Access Token
	Expires_in     int    `json:"expires_in"`     //Access Token的有效期,以秒为单位
	Refresh_token  string `json:"refresh_token"`  //用于刷新Access Token 的 Refresh Token,所有应用都会返回该参数;（10年的有效期）
	Session_key    string `json:"session_key"`    //基于http调用Open API时所需要的Session Key,其有效期与Access Token一致;
	Session_secret string `json:"session_secret"` //基于http调用Open API时计算参数签名用的签名密钥.
	/*
	    Access Token最终的访问范围，
	    即用户实际授予的权限列表（用户在授权页面时，有可能会取消掉某些请求的权限），
	    关于权限的具体信息参考
	    http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	*/
	Scope          string `json:"scope"`
}

type Credentials_ResponseErr struct {
	Error             string
	Error_description string
	Description       string
}