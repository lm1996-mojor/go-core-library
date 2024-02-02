package http_session

import (
	"github.com/kataras/iris/v12"
)

// GetCurrentHttpSessionUniqueKey 根据ConstValue获取key
func GetCurrentHttpSessionUniqueKey(ctx iris.Context) string {
	sessionId := ctx.GetID().(string)
	return sessionId
}
