package session_data_handler

import (
	"github.com/kataras/iris/v12"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/store"
)

func SessionDataInit(ctx iris.Context) {
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientID, "0")
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientCode, "0")
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserId, "0")
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserCode, "0")
	ctx.Next()
}
