package translate

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

const (
	API_URL = "http://translate.google.cn/translate_a/single?client=t&sl="
	constParam = "&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&ssel=0&tsel=0&kc=0&tk=335813.205684&q="
)

//http://translate.google.cn/translate_a/single?client=t&sl=en&tl=zh-CN&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&srcrom=0&ssel=0&tsel=0&kc=1&tk=371077.234804&q=target
//http://translate.google.cn/translate_a/single?client=t&sl=en&tl=zh-CN&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&ssel=3&tsel=3&kc=1&tk=152456.290617&q=who
//http://translate.google.cn/translate_a/single?client=t&sl=zh-CN&tl=en&hl=zh-CN&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&source=btn&ssel=3&tsel=3&kc=0&tk=335813.205684&q=%E7%9B%AE%E6%A0%87

func TestHello(text,targetLang,originLang string) error {

	response, err := http.Get(API_URL+originLang+"&tl="+targetLang+constParam+text)
	defer response.Body.Close()
	if nil != err {
		return err;
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err;
	}

	fmt.Println(string(body))
	return nil
}