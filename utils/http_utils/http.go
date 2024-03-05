package http_utils

import (
	"bytes"
	"encoding/json"
	"io"

	_const "github.com/lm1996-mojor/go-core-library/const"
)

func AddBodyParam(srcBody io.Reader, addParam map[string]interface{}) (newBody io.Reader) {
	var srcBodyMap interface{}
	//var srcBodyMap interface{}
	b, _ := io.ReadAll(srcBody)
	if len(b) > 0 {
		err := json.Unmarshal(b, &srcBodyMap)
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
