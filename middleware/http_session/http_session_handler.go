package http_session

import (
	"fmt"

	"github.com/kataras/iris/v12"
	_const "github.com/lm1996-mojor/go-core-library/const"
)

func SetCurrentHttpSessionUniqueKey(ctx iris.Context) {
	// 设置单次会话时获取数据源时使用的key
	ctx.Values().Set(_const.CurrentHttpSessionUniqueKey, fmt.Sprintf("%p", &ctx))
	ctx.Next()
}

// GetCurrentHttpSessionUniqueKey 根据ConstValue获取key
func GetCurrentHttpSessionUniqueKey(ctx iris.Context) string {
	return ctx.Values().Get(_const.CurrentHttpSessionUniqueKey).(string)
}
