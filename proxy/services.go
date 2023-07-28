package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	clog "mojor/go-core-library/log"
	"mojor/go-core-library/rest"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func Get(url string) (respBody []byte, err error) {
	// 进行请求，并获取响应流
	resp, err1 := http.Get(url)
	if err1 != nil {
		return nil, err1
	}
	//定制defer方法，在方法结束前固定执行。用于关闭响应流
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return checkResp(resp, nil)
}

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

	return checkResp(resp, req)
}

// GetClientDBList 获取全部租户数据库连接信息
func GetClientDBList(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	return
}

//func Post(url string) (respBody []byte, err error) {
//	// 进行请求，并获取响应流
//	resp, err := http.Post(url)
//	if err != nil {
//		return nil, err
//	}
//	//定制defer方法，在方法结束前固定执行。用于关闭响应流
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//
//		}
//	}(resp.Body)
//
//	return checkResp(resp, nil)
//}

// PostRequestWithParamsUsingJson 带参的post请求
// RequestType: POST
// RequestParam: map[string]interface{}
// ParamFormat: JSON
// Response: map[string]interface{}
// ERRORS: "业务报错/解析json到结构体出错/流关闭出现问题"
func PostRequestWithParamsUsingJson(params map[string]interface{}, url string) (map[string]interface{}, error) {
	//将请求数据转为二进制数组
	respData, _ := json.Marshal(params)
	resp, err1 := http.Post(url, "application/json", bytes.NewReader(respData))
	if err1 != nil {
		clog.Error(err1.Error())
		return nil, err1
	}
	//解析响应体数据为二进制数组([]byte)
	respBody, err := checkResp(resp, nil)
	if err != nil {
		return nil, err
	}
	//闭包关流
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)
	//将二进制数据解析为map[string]interface{}类型，以便后面取用
	body, err2 := ParseResponseBody(respBody)
	if err2 != nil {
		return nil, err
	}
	return body, nil
}

// DeleteRequest 带参的Delete请求
// RequestType: DELETE
// RequestParam: map[string]interface{} 请求参数
// ParamFormat: JSON 参数类型
// Response: map[string]interface{} 返回数据类型
// ERRORS: "业务报错/解析json到结构体出错/流关闭出现问题"
func DeleteRequest(params map[string]interface{}, url string) (map[string]interface{}, error) {
	//将请求数据转为二进制数组
	respData, _ := json.Marshal(params)
	req, _ := http.NewRequest("DELETE", url, bytes.NewReader(respData))
	req.Header.Add("Content-Type", "application/json")
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		return nil, err2
	}
	//解析响应体数据为二进制数组([]byte)
	respBody, err := checkResp(resp, nil)
	if err != nil {
		return nil, err
	}
	//闭包关流
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)
	//将二进制数据解析为map[string]interface{}类型，以便后面取用
	body, err1 := ParseResponseBody(respBody)
	if err1 != nil {
		return nil, err
	}
	return body, nil
}

// PutRequest 带参的Delete请求
// RequestType: DELETE
// RequestParam: map[string]interface{} 请求参数
// ParamFormat: JSON 参数类型
// Response: map[string]interface{} 返回数据类型
// ERRORS: "业务报错/解析json到结构体出错/流关闭出现问题"
func PutRequest(params map[string]interface{}, url string) (map[string]interface{}, error) {
	//将请求数据转为二进制数组
	respData, _ := json.Marshal(params)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(respData))
	req.Header.Add("Content-Type", "application/json")
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		return nil, err2
	}
	//解析响应体数据为二进制数组([]byte)
	respBody, err := checkResp(resp, nil)
	if err != nil {
		return nil, err
	}
	//闭包关流
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)
	//将二进制数据解析为map[string]interface{}类型，以便后面取用
	body, err1 := ParseResponseBody(respBody)
	if err1 != nil {
		return nil, err
	}
	return body, nil
}

// CustomRequestHeadAndMethod 自定义请求头和请求方式
// RequestType: Post
// RequestParam: bodyParams 请求体参数  headParams 请求头参数 url 请求地址
// ParamFormat: JSON 参数类型
// Response: map[string]interface{} 返回数据类型
// ERRORS: "业务报错/解析json到结构体出错/流关闭出现问题"
func CustomRequestHeadAndMethod(bodyParams map[string]interface{}, headParams map[string]string, url string, method string) (map[string]interface{}, error) {
	//将请求数据转为二进制数组
	method = strings.ToUpper(method)
	respData, _ := json.Marshal(bodyParams)
	req, _ := http.NewRequest(method, url, bytes.NewReader(respData))
	setFlag := false
	for key, value := range headParams {
		if key == "Content-Type" {
			setFlag = true
			req.Header.Add(key, value)
		}
		req.Header.Add(key, value)
	}
	//如果没有设置请求格式则设置默认请求格式
	if !setFlag {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		return nil, err2
	}
	//解析响应体数据为二进制数组([]byte)
	respBody, err := checkResp(resp, nil)
	if err != nil {
		return nil, err
	}
	//闭包关流
	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {
			panic("流关闭出现问题")
		}
	}(resp.Body)
	//将二进制数据解析为map[string]interface{}类型，以便后面取用
	body, err1 := ParseResponseBody(respBody)
	if err1 != nil {
		return nil, err
	}
	return body, nil
}

func checkResp(resp *http.Response, req *http.Request) ([]byte, error) {
	b, e := io.ReadAll(resp.Body)
	if e == nil && resp.StatusCode != 200 {
		var resp rest.Result
		err := json.Unmarshal(b, &resp)
		if err == nil && resp.Code != 200 && resp.Code != 0 {
			e = errors.New(resp.Msg)
			return nil, e
		}
	}
	return b, e
}

func parseTokenUtil(method string, url string, token string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)
	req.Header.Set("Content-type", "application/json")
	return
}

func ParseResponseBody(respBody []byte) (map[string]interface{}, error) {
	var resultMap map[string]interface{}
	//使用json解析响应体中的数据，并存入输出结构体中
	err := json.Unmarshal(respBody, &resultMap)
	if err != nil {
		clog.Errorf("services.go 解析json到结构体出错 ", err)
		return nil, errors.New("proxy->解析json到结构体出错:109")
	}
	return resultMap, nil
}
