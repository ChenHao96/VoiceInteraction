package secret

import (
	"com/tuling123/interfaces"
	"encoding/json"
	"time"
	"strconv"
	"crypto/md5"
	"net/http"
	"bytes"
	"io/ioutil"
	"com/tuling123"
	"org/StevenChen/util"
	"encoding/base64"
	"fmt"
)

type API_Request struct {
	Key    string `json:"key"`  //必须	32位	1ca8089********736b8ce41591426(32位)		注册之后在机器人接入页面获得（参见本文档第2部分）
	Info   string `json:"info"` //必须	1-30位	打招呼“你好”;查天气“北京今天天气”;			请求内容，编码方式为UTF-8
	secret string
}

type request struct {
	Key       string        `json:"key"`
	Timestamp string        `json:"timestamp"`
	Data      string        `json:"data"`
}

func NewAPI_Request(key, secret string) interfaces.API_Util {

	var request API_Request
	request.Key = key
	request.secret = secret

	return request
}

func (this API_Request) getData(data []byte) []byte {

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)+"000"

	keyParam := this.secret + timestamp + this.Key
	key := md5.Sum([]byte(keyParam))

	//这里有个毛病   数组无法直接赋值个切片
	key2 := make([]byte, 0, 16)
	for _, value := range key {
		key2 = append(key2, value)
	}

	/*bug 一堆堆的  golang的AES加密我不知为什么得到的数据与java的不一样*/
	key2 = []byte(fmt.Sprintf("%x",key2))
	data, err := util.AesEncrypt(data, key2)
	if nil != err {
		panic(err.Error())
	}

	dataP := base64.StdEncoding.EncodeToString(data)

	requestParam := &request{Key:this.Key, Timestamp:timestamp, Data:dataP}
	param, err := json.Marshal(requestParam)
	if nil != err {
		panic(err.Error())
	}

	return param
}

func (this API_Request) Talk(worlds string) string {

	this.Info = worlds
	data, err := json.Marshal(this)
	if nil != err {
		panic(err.Error())
	}

	data = this.getData(data)

	response, err := http.Post(tuling123.API_URL, "application/json; charset=utf-8", bytes.NewReader(data));
	defer response.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		panic(err.Error())
	}

	var first tuling123.API_Response
	json.Unmarshal(body, &first);

	if 100000 != first.Code {
		if content, ok := tuling123.API_Response_Err[first.Code]; ok {
			panic(content)
		} else {
			panic("Unknow,该信息暂时无法识别")
		}
	}

	return first.Text
}