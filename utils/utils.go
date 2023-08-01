package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/store"
)

func RandStr(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)

	return str[:len]
}

func ZeroValue(v interface{}) bool {
	value := reflect.ValueOf(v)
	if value.Interface() == reflect.Zero(value.Type()).Interface() {
		return true
	} else {
		return false
	}
}

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}

func SplitRedisPort(ports string) (split []string) {
	split = strings.Split(ports, ",")
	if len(split) == 0 {
		panic("redis 端口解析错误：检查yaml配置：redis-port-> " + ports)
	}
	return split
}

// ObtainClientId 获取当前的租户id
func ObtainClientId() (clientId int64, err error) {
	value, ok := store.Get(fmt.Sprintf("%p", &store.PoInterKey) + _const.ClientID)
	if !ok {
		return 0, errors.New("租户不确定")
	}
	//租户id参数类型转换(string -> int64)
	cId, err1 := strconv.ParseInt(value.(string), 10, 32)
	if err1 != nil {
		return 0, errors.New("租户参数转换失败")
	}
	return cId, nil
}
