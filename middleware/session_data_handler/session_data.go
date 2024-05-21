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
	reqPath := ctx.Path()
	ctx.Values().Set("pass_label", "N")
	if white_list.InList(reqPath, ctx.Request().Method, 1) || strings.Contains(reqPath, "platform_management") {
		ctx.Values().Set("pass_label", "Y")
		ctx.Next()
		return
	}
	currentReqHeader := ctx.GetHeader("Content-Type")
	if strings.Contains(currentReqHeader, "multipart/form-data") {
		sessionParam := ctx.FormValue(_const.HttpSessionParam)
		storeParamHandler(ctx, []byte(sessionParam))
		ctx.Request().Form.Del(_const.HttpSessionParam)
	} else {
		all, _ := io.ReadAll(ctx.Request().Body)
		param := storeParamHandler(ctx, all)
		marshal, _ := json.Marshal(param[_const.OriginalReqParam])
		ctx.Request().Body = io.NopCloser(bytes.NewReader(marshal))
	}
	ctx.Next()
}

func storeParamHandler(ctx iris.Context, sessionParam []byte) map[string]interface{} {
	param := make(map[string]interface{})
	json.Unmarshal(sessionParam, &param)
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientID, param[_const.ClientID].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientCode, param[_const.ClientCode].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserId, param[_const.UserId].(string))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserCode, param[_const.UserCode].(string))
	//将解析后的token中的用户信息存入local store
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.JwtData, param[_const.JwtData].(map[string]interface{}))
	store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.TokenOriginal, param[_const.TokenOriginal].(string))
	delete(param, _const.ClientID)
	delete(param, _const.ClientCode)
	delete(param, _const.UserId)
	delete(param, _const.UserCode)
	delete(param, _const.JwtData)
	delete(param, _const.TokenOriginal)
	return param
}
