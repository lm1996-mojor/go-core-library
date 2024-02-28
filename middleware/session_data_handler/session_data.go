package session_data_handler

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/kataras/iris/v12"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/middleware/security/auth/white_list"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/utils"
)

func SessionDataInit(ctx iris.Context) {
	utils.PrintCallerInfo(ctx)
	reqPath := ctx.Request().RequestURI
	ctx.Values().Set("pass_label", "N")
	if white_list.InList(reqPath, 1) || strings.Contains(reqPath, "platform_management") || strings.Contains(reqPath, "platform_inlet") {
		ctx.Values().Set("pass_label", "Y")
		ctx.Next()
		return
	}
	param := make(map[string]interface{})
	all, _ := io.ReadAll(ctx.Request().Body)
	if len(all) > 0 {
		json.Unmarshal(all, &param)
	}
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientID, param[_const.ClientID].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientCode, param[_const.ClientCode].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserId, param[_const.UserId].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserCode, param[_const.UserCode].(string))
	//将解析后的token中的用户信息存入local store
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.JwtData, param[_const.JwtData].(map[string]interface{}))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.TokenOriginal, param[_const.TokenOriginal].(string))
	marshal, _ := json.Marshal(param)
	ctx.Request().Body = io.NopCloser(bytes.NewReader(marshal))
	ctx.Next()
}
