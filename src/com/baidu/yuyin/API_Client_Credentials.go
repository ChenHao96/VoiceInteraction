package yuyin

/*
	客户凭证授权
	获取调用API的token
 */
const Credentials_Url = "https://openapi.baidu.com/oauth/2.0/token"

const (
	INVALID_REQUEST = &Credentials_ResponseErr{error:"invalid_request",error_description:"invalid refresh token",description:"请求缺少某个必需参数，包含一个不支持的参数或参数值，或者格式不正确。"}
	INVALID_CLIENT = &Credentials_ResponseErr{error:"invalid_client",error_description:"unknown client id",description:"client_id”、“client_secret”参数无效"}
	INVALID_GRANT = &Credentials_ResponseErr{error:"invalid_grant",error_description:"The provided authorization grant is revoked",description:"提供的Access Grant是无效的、过期的或已撤销的，例如，Authorization Code无效(一个授权码只能使用一次)、Refresh Token无效、redirect_uri与获取Authorization Code时提供的不一致、Devie Code无效(一个设备授权码只能使用一次)等。"}
	UNAUTHORIZED_CLIENT = &Credentials_ResponseErr{error:"unauthorized_client",error_description:"The client is not authorized to use this authorization grant type",description:"应用没有被授权，无法使用所指定的grant_type。"}
	UNSUPPORTED_GRANT_TYPE = &Credentials_ResponseErr{error:"unsupported_grant_type",error_description:"The authorization grant type is not supported",description:"“grant_type”百度OAuth2.0服务不支持该参数。"}
	INVALID_SCOPE = &Credentials_ResponseErr{error:"invalid_scope",error_description:"The requested scope is exceeds the scope granted by the resource owner",description:"请求的“scope”参数是无效的、未知的、格式不正确的、或所请求的权限范围超过了数据拥有者所授予的权限范围。"}
	EXPIRED_TOKEN = &Credentials_ResponseErr{error:"expired_token",error_description:"refresh token has been used",description:"提供的Refresh Token已过期"}
	REDIRECT_URI_MISMATCH = &Credentials_ResponseErr{error:"redirect_uri_mismatch",error_description:"Invalid redirect uri",description:"“redirect_uri”所在的根域与开发者注册应用时所填写的根域名不匹配。"}
	UNSUPPORTED_RESPONSE_TYPE = &Credentials_ResponseErr{error:"unsupported_response_type",error_description:"The response type is not supported",description:"“response_type”参数值不为百度OAuth2.0服务所支持，或者应用已经主动禁用了对应的授权模式"}
	SLOW_DOWN = &Credentials_ResponseErr{error:"slow_down",error_description:"The device is polling too frequently",description:"Device Flow中，设备通过Device Code换取Access Token的接口过于频繁，两次尝试的间隔应大于5秒。"}
	AUTHORIZATION_PENDING = &Credentials_ResponseErr{error:"authorization_pending",error_description:"User has not yet completed the authorization",description:"Device Flow中，用户还没有对Device Code完成授权操作。"}
	AUTHORIZATION_DECLINED = &Credentials_ResponseErr{error:"authorization_declined",error_description:"User has declined the authorization",description:"Device Flow中，用户拒绝了对Device Code的授权操作。"}
	INVALID_REFERER = &Credentials_ResponseErr{error:"invalid_referer",error_description:"Invalid Referer",description:"Implicit Grant模式中，浏览器请求的Referer与根域名绑定不匹配"}
)

type Credentials_Request struct {
	grant_type    string "client_credentials" //必填参数 固定为“client_credentials”；
	client_id     string                      //必填参数 应用的API Key
	client_secret string                      //必须参数 应用的Secret Key;
	/*
	  非必须参数。
	  以空格分隔的权限列表，采用本方式获取Access Token时只能申请跟用户数据无关的数据访问权限。
	  关于权限的具体信息请参考
	  http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	*/
	scope         string
}

type Credentials_Response struct {
	access_token   string	//要获取的Access Token
	expires_in     int	//Access Token的有效期,以秒为单位
	refresh_token  string	//用于刷新Access Token 的 Refresh Token,所有应用都会返回该参数;（10年的有效期）
	session_key    string	//基于http调用Open API时所需要的Session Key,其有效期与Access Token一致;
	session_secret string	//基于http调用Open API时计算参数签名用的签名密钥.
	/*
	   Access Token最终的访问范围，
	   即用户实际授予的权限列表（用户在授权页面时，有可能会取消掉某些请求的权限），
	   关于权限的具体信息参考
	   http://developer.baidu.com/wiki/index.php?title=docs/oauth/list
	 */
	scope          string
}

type Credentials_ResponseErr struct {
	error string
	error_description string
	description string
}
