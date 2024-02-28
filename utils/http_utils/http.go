package http_utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func AddBodyParam(srcBody io.Reader, addParam map[string]interface{}) (newBody io.Reader) {
	var srcBodyMap map[string]interface{}
	b, _ := io.ReadAll(srcBody)
	err := json.Unmarshal(b, &srcBodyMap)
	if err != nil {
		panic(err)
	}
	for key, value := range addParam {
		srcBodyMap[key] = value
	}
	marshal, err1 := json.Marshal(srcBodyMap)
	if err1 != nil {
		panic(err1)
	}
	newBody = bytes.NewReader(marshal)
	return newBody
}
