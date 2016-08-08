package yuyin

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"net"
	"encoding/base64"
	"bytes"
	"fmt"
	"strconv"
	"net/url"
	"errors"
	"time"
)

type API_Util struct {
	Credentials Credentials_Response
	Cuid        string
}

/*
	获取一个本地的MAC地址作为API的 cuid
 */
func getCUID() (cuId string, err error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	cuId = fmt.Sprintf("%s", interfaces[0].HardwareAddr)
	return
}

func getCredentials(request Credentials_Request) (response Credentials_Response, err error) {

	postValue := url.Values{};
	postValue.Set("scope", request.Scope)
	postValue.Set("client_id", request.Client_id)
	postValue.Set("grant_type", "client_credentials")
	postValue.Set("client_secret", request.Client_secret)

	postResponse, err := http.PostForm(Credentials_Url, postValue)
	if err != nil {
		return
	}
	defer postResponse.Body.Close()

	body, err := ioutil.ReadAll(postResponse.Body);
	if err != nil {
		return
	}

	var result = make(map[string]string)
	err = json.Unmarshal(body, &result);

	if err == nil {
		if value, ok := result["error"]; ok {
			err = errors.New(Credentials_ResponseErrEnum[value].Description)
			return
		}
	}

	err = json.Unmarshal(body, &response);

	return
}

func NewAPI_Util(api_key, secret_key string) (util API_Util, err error) {

	cuid, err := getCUID()
	if err != nil {
		return
	}

	res, err := getCredentials(Credentials_Request{
		Client_id:api_key, Client_secret:secret_key})

	util.Cuid = cuid
	util.Credentials = res

	return
}

func printBase64Binary(val []byte) string {

	base64Table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encoding := base64.NewEncoding(base64Table);

	return encoding.EncodeToString(val)
}

func (this API_Util) SendBytesRequest(filePath, format string, rate int) (result API_Response, err error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	soundStr := printBase64Binary(data)

	param := &API_Request{Speech:soundStr, Cuid:this.Cuid,
		Token:this.Credentials.Refresh_token}
	param.Rate = rate
	param.Channel = 1
	param.Len = len(data)
	param.Format = format

	postValue, err := json.Marshal(param)
	if err != nil {
		return
	}

	begin := time.Now()
	response, err := http.Post(API_URL,
		"application/json; charset=utf-8",
		bytes.NewReader(postValue));

	if err != nil {
		return
	}
	defer response.Body.Close()
	end := time.Now()

	fmt.Println("SendBytesRequest用时:",end.Sub(begin))

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	var first = make(map[string]string)
	err = json.Unmarshal(body, &result);

	if err == nil {
		if value, ok := first["err_code"]; ok {
			code, _ := strconv.Atoi(value)
			err = errors.New(API_ResponseErrEnum[code].Meaning)
			return
		}
	}

	err = json.Unmarshal(body, &result);
	return
}

func (this API_Util) SendFileRequest(filePath, format string, rate int) (result API_Response, err error) {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	param := bytes.NewReader(data)
	url := API_URL + "?cuid=" + this.Cuid + "&token=" + this.Credentials.Refresh_token
	contentType := "audio/" + format + "; rate=" + strconv.Itoa(rate)

	begin := time.Now()
	response, err := http.Post(url, contentType, param);
	if err != nil {
		return
	}
	defer response.Body.Close()
	end := time.Now()

	fmt.Println("SendFileRequest用时:",end.Sub(begin))

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return
	}

	var first = make(map[string]string)
	err = json.Unmarshal(body, &result);

	if err == nil {
		if value, ok := first["err_code"]; ok {
			code, _ := strconv.Atoi(value)
			err = errors.New(API_ResponseErrEnum[code].Meaning)
			return
		}
	}

	err = json.Unmarshal(body, &result);
	return
}