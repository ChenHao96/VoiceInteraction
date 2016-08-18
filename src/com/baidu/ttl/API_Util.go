package ttl

import (
	"com/baidu/public"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"errors"
	"net/url"
)

//REST API Url
const API_URL = "http://tsn.baidu.com/text2audio"

type API_Request struct {
	Tex  string `json:"tex"`           //必填；合成文本，使用UTF-8编码，请注意文本长度必须小于1024
	Lan  string `json:"lan"`           //必填；语言选择，填写zh
	Tok  string `json:"tok"`           //必填；开放平台获取到的开发者access_token
	Ctp  int    `json:"ctp"`           //必填；客户端类型选择,web端填写1
	Cuid string `json:"cuid"`          //必填；用户唯一标识，用来区分用户，填写机器 MAC 地址或 IMEI 码，长度为60以内
	Spd  int    `json:"spd,omitempty"` //选填；语速，取值0-9，默认5
	Pit  int    `json:"pit,omitempty"` //选填；语调，取值0-9，默认5
	Vol  int    `json:"vol,omitempty"` //选填；音量，取值0-9，默认5
	Per  int    `json:"per,omitempty"` //选填；发音人选择，取值0-1；默认0女声 1男声
}

type API_Response struct {
	Err_no  int        `json:"err_no"`
	Err_msg string     `json:"err_msg"`
	Sn      string     `json:"sn"`
	Idx     int        `json:"idx"`
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

func (this API_Util) Text2AudioFile(filePath, text string) (err error) {

	body, err := this.Text2AudioBytes(text)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filePath, body, 0666)
	return
}

func (this API_Util) Text2AudioBytes(text string) (data []byte, err error) {

	param := url.Values{}
	param.Set("ctp", "1")
	param.Set("lan", "zh")
	param.Set("tex", text)
	param.Set("cuid", this.Cuid)
	param.Set("tok", this.Credentials.Access_token)

	response, err := http.PostForm(API_URL, param)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	contentType := response.Header.Get("Content-type")
	if "audio/mp3" == contentType {
		data = body
	} else {
		var errMsg API_Response
		err = json.Unmarshal(body, &errMsg);
		if err == nil {
			err = errors.New(errMsg.Err_msg)
		}
	}

	return
}