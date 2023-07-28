package rest

import (
	"encoding/json"
)

// Result HTTP request result
type Result struct {
	Code    int         `json:"code"`    //响应状态码
	Msg     string      `json:"msg"`     //响应信息
	Data    interface{} `json:"data"`    //响应数据
	MsgType string      `json:"msgType"` //响应类型
}

const (
	INFO    = "info"
	Success = "success"
	WARNING = "warning"
	ERROR   = "error"
	LOADING = "loading"
)

// NewResult create new result
func NewResult(data interface{}, c int, msg string) Result {
	return Result{Data: data, Code: c, Msg: msg}
}

// NewQueryPage create new page result
func NewQueryPage(items interface{}, page, size int, total int64) Result {
	pageable := NewPageable(page, size, total)
	data := map[string]interface{}{
		"content":  items,
		"pageable": pageable,
	}
	return NewResult(data, 200, "ok")
}

// SuccessResult 常用成功响应
// @@Param:data 成功需要返回的数据，可以为nil
// @Description: 固定code：200，message：请求成功
func SuccessResult(data map[string]interface{}) Result {
	return Result{
		Code: GetHttpStatus(200).Value,
		Msg:  GetHttpStatus(200).ReasonPhrase,
		Data: data,
		//MsgType: Success,
	}
}

// SuccessCustom 自定义消息成功响应
// @@Param: data 成功需要返回的数据，可以为nil
// @@Param: msg 自定义需要返回的消息
// @Description: 固定code：200
func SuccessCustom(msg string, data map[string]interface{}, msgType string) Result {
	return Result{
		Code: GetHttpStatus(200).Value,
		Msg:  msg,
		Data: data,
		//MsgType: msgType,
	}
}

// FailResult 常用失败响应
// @@Param:data 成功需要返回的数据，可以为nil
// @Description: 固定code：500，message：服务器错误
func FailResult() Result {
	return Result{
		Code:    GetHttpStatus(500).Value,
		Msg:     GetHttpStatus(500).ReasonPhrase,
		MsgType: ERROR,
	}
}

// BadResult 常用失败响应
// @@Param:data 成功需要返回的数据，可以为nil
// @Description: 固定code：500，message：自定义
func BadResult(msg string) Result {
	return Result{
		Code:    GetHttpStatus(500).Value,
		Msg:     msg,
		MsgType: ERROR,
	}
}

// FailCustom 自定义失败响应
// @@Param:data 成功需要返回的数据，可以为nil
// @Description: 内容自定义
func FailCustom(code int, msg string, msgType string) Result {
	return Result{
		Code:    code,
		Msg:     msg,
		MsgType: msgType,
	}
}

// FailCustomBinaryResponse 自定义失败响应(二进制返回)
// @@Param:data 成功需要返回的数据，可以为nil
// @Description: 内容自定义
func FailCustomBinaryResponse(code int, msg string) []byte {
	marshal, _ := json.Marshal(Result{
		Code:    code,
		Msg:     msg,
		MsgType: ERROR,
	})
	return marshal
}
