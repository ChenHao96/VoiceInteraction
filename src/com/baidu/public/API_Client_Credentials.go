package public

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
	客户凭证授权
	获取调用API的token
*/
const Credentials_Url = "https://openapi.baidu.com/oauth/2.0/token"

var Credentials_ResponseErrEnum = map[string][2]string{
	"invalid_request":           {"invalid refresh token", "请求缺少某个必需参数，包含一个不支持的参数或参数值，或者格式不正确。"},
	"invalid_client":            {"unknown client id", "client_id或client_secret参数无效"},
	"invalid_grant":             {"The provided authorization grant is revoked", "提供的Access Grant是无效的、过期的或已撤销的，例如，Authorization Code无效(一个授权码只能使用一次)、Refresh Token无效、redirect_uri与获取Authorization Code时提供的不一致、Devie Code无效(一个设备授权码只能使用一次)等。"},
	"unauthorized_client":       {"The client is not authorized to use this authorization grant type", "应用没有被授权，无法使用所指定的grant_type。"},
	"unsupported_grant_type":    {"The authorization grant type is not supported", "“grant_type”百度OAuth2.0服务不支持该参数。"},
	"invalid_scope":             {"The requested scope is exceeds the scope granted by the resource owner", "请求的“scope”参数是无效的、未知的、格式不正确的、或所请求的权限范围超过了数据拥有者所授予的权限范围。"},
	"expired_token":             {"refresh token has been used", "提供的Refresh Token已过期"},
	"redirect_uri_mismatch":     {"Invalid redirect uri", "“redirect_uri”所在的根域与开发者注册应用时所填写的根域名不匹配。"},
	"unsupported_response_type": {"The response type is not supported", "“response_type”参数值不为百度OAuth2.0服务所支持，或者应用已经主动禁用了对应的授权模式"},
	"slow_down":                 {"The device is polling too frequently", "Device Flow中，设备通过Device Code换取Access Token的接口过于频繁，两次尝试的间隔应大于5秒。"},
	"authorization_pending":     {"User has not yet completed the authorization", "Device Flow中，用户还没有对Device Code完成授权操作。"},
	"authorization_declined":    {"User has declined the authorization", "Device Flow中，用户拒绝了对Device Code的授权操作。"},
	"invalid_referer":           {"Invalid Referer", "Implicit Grant模式中，浏览器请求的Referer与根域名绑定不匹配"},
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

type Credentials_Response_Err struct {
	Error             string `json:"error"`
	Error_description string `json:"error_description"`
}

func GetCredentials(request Credentials_Request) Credentials_Response {

	postValue := url.Values{}
	postValue.Set("scope", request.Scope)
	postValue.Set("client_id", request.Client_id)
	postValue.Set("grant_type", "client_credentials")
	postValue.Set("client_secret", request.Client_secret)

	postResponse, err := http.PostForm(Credentials_Url, postValue)
	if err != nil {
		panic(err.Error())
	}
	defer postResponse.Body.Close()

	body, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		panic(err.Error())
	}

	var result Credentials_Response_Err
	if nil == json.Unmarshal(body, &result) {
		if description ,ok := Credentials_ResponseErrEnum[result.Error];ok{
			panic(description[1])
		}
	}

	var response Credentials_Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err.Error())
	}

	return response
}