package http_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"strings"

	"github.com/kataras/iris/v12"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/log"
)

func AddBodyParam(ctx iris.Context, srcBody io.Reader, addParam map[string]interface{}) iris.Context {
	currentReqHeader := ctx.GetHeader("Content-Type")
	if strings.Contains(currentReqHeader, "multipart/form-data") {
		body, contentType := fileReqHandler(ctx, addParam)
		ctx.Request().Body = io.NopCloser(body)
		ctx.Header("Content-Type", contentType)
		return ctx
	} else {
		ctx.Request().Body = io.NopCloser(usualReqHandler(srcBody, addParam))
		return ctx
	}
	//var srcBodyMap interface{}
	//srcData, _ := io.ReadAll(srcBody)
	//if len(srcData) > 0 {
	//	err := json.Unmarshal(srcData, &srcBodyMap)
	//	if err != nil {
	//		panic(err)
	//	}
	//	addParam[_const.OriginalReqParam] = srcBodyMap
	//}
	//marshal, err1 := json.Marshal(addParam)
	//if err1 != nil {
	//	panic(err1)
	//}
	//newBody = bytes.NewReader(marshal)
	//return newBody
}
func usualReqHandler(srcBody io.Reader, addParam map[string]interface{}) (newBody io.Reader) {
	var srcBodyMap interface{}
	srcData, _ := io.ReadAll(srcBody)
	if len(srcData) > 0 {
		err := json.Unmarshal(srcData, &srcBodyMap)
		if err != nil {
			panic(err)
		}
		addParam[_const.OriginalReqParam] = srcBodyMap
	}
	marshal, err1 := json.Marshal(addParam)
	if err1 != nil {
		panic(err1)
	}
	newBody = bytes.NewReader(marshal)
	return newBody
}

func fileReqHandler(ctx iris.Context, addParam map[string]interface{}) (newBody io.Reader, contentType string) {
	//var srcBodyMap interface{}
	_, fileHeader, _ := ctx.FormFile(_const.FileRequestKey)
	//fileHeaders := ctx.Request().MultipartForm.File[fileKey]
	// 创建一个新的 buffer 来构建新的 multipart/form-data 请求体
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// 写入原始的文件（如果需要的话）
	if fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			log.Error("文件读取异常:" + err.Error())
			panic("服务器错误")
		}
		defer file.Close()

		part, err1 := w.CreateFormFile(_const.FileRequestKey, fileHeader.Filename)
		if err1 != nil {
			log.Error("文件读取异常:" + err.Error())
			panic("服务器错误")
		}
		_, err = io.Copy(part, file)
		if err != nil {
			log.Error("文件读取异常:" + err.Error())
			panic("服务器错误")
		}
	}
	for k, v := range ctx.FormValues() {
		if len(v) > 1 {
			values := "" // 1,1,1,1,1,1,1
			for i, s := range v {
				// [ 1, 2,2,3,4,5]
				if i >= len(v)-1 {
					values += s
				} else {
					values += s + ","
				}
			}
			w.WriteField(k, values)
		} else {
			w.WriteField(k, v[0])
		}
	}
	newParamJSON, _ := json.Marshal(addParam)
	w.WriteField(_const.HttpSessionParam, string(newParamJSON))
	w.Close()
	return &b, w.FormDataContentType()
}
