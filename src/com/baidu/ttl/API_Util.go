package ttl

import (
	"com/baidu/public"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"errors"
	"net/url"
)

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

func (this API_Util) Text2AudioFile(filePath, fileNam, text string) (err error) {

	param := url.Values{}
	param.Set("ctp","1")
	param.Set("lan","zh")
	param.Set("tex",text)
	param.Set("cuid",this.Cuid)
	param.Set("tok",this.Credentials.Access_token)

	response, err := http.PostForm(API_URL,param)
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
		err = ioutil.WriteFile(filePath+fileNam, body, 0666)
	} else {
		var errMsg API_Response
		err = json.Unmarshal(body, &errMsg);
		if err == nil{
			err = errors.New(errMsg.Err_msg)
		}
	}

	return
}

func (this API_Util) Text2AudioBytes(text string) (data []byte, err error) {

	param := url.Values{}
	param.Set("ctp","1")
	param.Set("lan","zh")
	param.Set("tex",text)
	param.Set("cuid",this.Cuid)
	param.Set("tok",this.Credentials.Access_token)

	response, err := http.PostForm(API_URL,param)
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
		if err == nil{
			err = errors.New(errMsg.Err_msg)
		}
	}

	return
}