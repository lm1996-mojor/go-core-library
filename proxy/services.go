package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest"
)

const (
	GET = "GET"
)

func GetParseToken(token string, url string) (respBody []byte, err error) {
	client := &http.Client{Timeout: 60 * time.Second}
	req, err1 := parseTokenUtil(GET, url, token)
	if err1 != nil {
		return nil, err1
	}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)

	return checkResp(resp)
}

// RequestAction 请求执行
//
// reqMdl 请求内容结构体
//
// reqMod 请求模式（sync[同步]/其他值默认为异步）
func RequestAction(reqMdl *RemoteReqMdl, reqMod string) (respParam map[string]interface{}, err error) {
	if reqMdl == nil {
		return nil, errors.New("请求对象为空")
	}
	if reqMod == "" {
		reqMod = "sync"
	}
	if reqMod == "sync" {
		// 同步请求
		remoteRequestHandler(reqMdl, &respParam, &err)
	} else {
		// 异步请求
		go remoteRequestHandler(reqMdl, &respParam, &err)
	}
	return respParam, err
}

// 远程请求处理器
func remoteRequestHandler(reqMdl *RemoteReqMdl, respParam *map[string]interface{}, err *error) {
	reqMdl.Method = strings.ToUpper(reqMdl.Method)
	//将请求数据转为二进制数组
	var reader *bytes.Reader
	if len(reqMdl.ReqParam) > 0 {
		reader = bytes.NewReader(reqMdl.ReqParam)
	} else {
		reader = nil
	}

	req, _ := http.NewRequest(reqMdl.Method, reqMdl.Url, reader)
	// 判断是否需要携带token
	if reqMdl.CarryToken {
		builderTokenHead(req, reqMdl.Token)
	}
	// 判断是否需要构建请求头数据
	if len(reqMdl.HeadParams) > 0 {
		builderHttpRequestHeader(req, reqMdl.HeadParams)
	}
	// 开始请求
	resp, err2 := http.DefaultClient.Do(req)
	err = &err2
	//解析响应体数据为二进制数组([]byte)
	respBody, err3 := checkResp(resp)
	err = &err3
	//闭包关流
	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)
	//将二进制数据解析为map[string]interface{}类型，以便后面取用
	body, err4 := ParseResponseBody(respBody)
	respParam = &body
	err = &err4
}

func checkResp(resp *http.Response) ([]byte, error) {
	b, e := io.ReadAll(resp.Body)
	if e == nil && resp.StatusCode != 200 {
		var result rest.Result
		err := json.Unmarshal(b, &result)
		if err == nil && result.Code != 200 && result.Code != 0 {
			e = errors.New(result.Msg)
			return nil, e
		}
	}
	return b, e
}

func builderHttpRequestHeader(req *http.Request, headParams map[string]string) {
	for key, value := range headParams {
		req.Header.Add(key, value)
	}
}

func builderTokenHead(req *http.Request, token string) {
	req.Header.Add("Authorization", token)
	req.Header.Set("Content-type", "application/json")
}

func parseTokenUtil(method string, url string, token string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	builderTokenHead(req, token)
	return
}

func ParseResponseBody(respBody []byte) (map[string]interface{}, error) {
	var resultMap map[string]interface{}
	//使用json解析响应体中的数据，并存入输出结构体中
	err := json.Unmarshal(respBody, &resultMap)
	if err != nil {
		clog.Errorf("解析json到结构体出错 ", err)
		return nil, errors.New("proxy->解析json到结构体出错:109")
	}
	return resultMap, nil
}
