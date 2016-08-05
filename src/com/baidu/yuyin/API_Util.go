package yuyin

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"net"
	"encoding/base64"
	"net/url"
	"fmt"
	"bytes"
)

/*
	获取一个本地的MAC地址作为API的 cuid
 */
func GetCUID() (cuId string, err error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	cuId = string(interfaces[0].HardwareAddr)
	return
}

func GetToken(API_Key, Secret_Key string) (token string, err error) {

	//不知道这里有没有更好的方案，传递结构体就能作为参数
	postValue := url.Values{};
	postValue.Set("client_id", API_Key)
	postValue.Set("client_secret", Secret_Key)
	postValue.Set("grant_type", "client_credentials")

	response, err := http.PostForm(Credentials_Url, postValue)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(body), &result);
	if err != nil {
		return
	}

	token = fmt.Sprintf("%s", result["access_token"])
	return
}

func SendBytesRequest(filePath, token, cuid string) (result string, err error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	soundStr := printBase64Binary(data)

	param := &API_Request{Speech:soundStr, Token:token, Cuid:cuid}
	param.Format = "pcm"
	param.Rate = 8000
	param.Channel = 1
	param.Len = len(data)

	postValue, err := json.Marshal(param)
	if err != nil {
		return
	}

	response, err := http.Post(API_URL,
		"application/json; charset=utf-8",
		bytes.NewReader(postValue));

	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	result = string(body)
	return
}

func SendFileRequest(filePath, token, cuid string) (result string, err error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	param := bytes.NewReader(data)
	url := API_URL+"?cuid="+cuid+"&token="+token
	response, err := http.Post(url,"audio/pcm; rate=8000", param);
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	result = string(body)
	return

}

func printBase64Binary(val []byte) string {

	base64Table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encoding := base64.NewEncoding(base64Table);

	return encoding.EncodeToString(val)
}