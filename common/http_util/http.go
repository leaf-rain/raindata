package http_util

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type RequestBody struct {
	Cmd  int64       `json:"cmd"`
	Data interface{} `json:"data,omitempty"`
}

type ResponseBody struct {
	Code int64       `json:"code"`
	Cmd  int64       `json:"cmd,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func PostJSONReceiveJSON(client *http.Client, url string, send, receive interface{}) (int, error) {
	code, body, err := PostJson(client, url, send)
	if err != nil {
		return code, err
	}
	err = json.Unmarshal(body, receive)
	if err != nil {
		return code, err
	}
	return code, nil
}

func PostJson(client *http.Client, url string, data interface{}) (int, []byte, error) {
	js, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	//创建一个新的post请求
	var request *http.Request
	request, err = http.NewRequest("POST", url, strings.NewReader(string(js)))
	if err != nil {
		return 0, nil, err
	}
	//请求头设置
	request.Header.Add("Content-Type", "application/json") //json请求
	//发送请求到服务端
	var resp *http.Response
	resp, err = client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}
