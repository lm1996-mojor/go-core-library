package token

import (
	"encoding/json"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/consul"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/middleware/security/auth/white_list"
	"github.com/lm1996-mojor/go-core-library/proxy"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/utils"
)

func CheckIdentity(ctx iris.Context) {
	//获取请求路径
	utils.PrintCallerInfo(ctx)
	reqPath := ctx.Path()
	ctx.Values().Set("pass_label", "N")
	if white_list.InList(reqPath, 1) {
		ctx.Values().Set("pass_label", "Y")
		ctx.Next()
		return
	}
	//获取token(根据请求头信息获取不同形式的token：http/websocket)
	author := ""
	if strings.Contains(ctx.Request().Proto, "HTTP") {
		author = ctx.GetHeader(_const.TokenName)
	} else {
		author = ctx.GetHeader(_const.WebSocketTokenStoreHttpRequestHeaderName)
	}
	if author == "" || author == "null" || len(author) <= 0 {
		ctx.JSON(rest.FailCustom(401, "尚未登录,请登录后再进行操作", rest.ERROR))
		return
	}

	//去除token 头部信息
	var token string
	if strings.Contains(author, "Bearer ") {
		token = author[7:]
	} else {
		token = author
	}
	tokenService := consul.ObtainHighestWeightInServiceList(config.Sysconfig.Detection.TokenService)
	url := tokenService.Proto + "://" + tokenService.Host + config.Sysconfig.Detection.TokenCheckServiceApiUrl
	// 获取解析后的token信息
	respBody, err := proxy.GetParseToken(token, url)
	if err != nil {
		clog.Error("token解析: 响应出错" + err.Error())
		panic("服务器错误")
		return
	}
	// 解析响应体中的数据
	result, parseRespErr := parseResponseBody(respBody)
	if parseRespErr != nil {
		clog.Error(parseRespErr.Error())
		panic("服务器错误")
		return
	}

	if result.Code == 200 {
		//判断自定义的token类型是否正确
		tokenClaims := result.Data.(map[string]interface{})["parse_token"].(map[string]interface{})
		if t, ok := tokenClaims["token_type"].(string); ok && t != _const.TokenType { //不是access token
			clog.Warn("令牌类型认证无效: " + err.Error())
			clog.Warn("无效令牌：" + token)
			ctx.JSON(rest.FailCustom(401, "登录信息无效，请重新登录", rest.ERROR))
			return
		}
		if t, ok := tokenClaims["token_single"].(string); ok && t != _const.TokenSignature {
			//不是access token
			clog.Warn("令牌签名认证无效: " + err.Error())
			clog.Warn("无效令牌：" + token)
			ctx.JSON(rest.FailCustom(401, "登录信息无效，请重新登录", rest.ERROR))
			return
		}
		// 以下所有数据都会在单次回话完成后进行清空
		// 用于判断是否为超级管理员，主要用在鉴权时是否需要走权限系统
		if config.Sysconfig.Detection.Authentication {
			ctx.Values().Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+"isSuperAdmin", tokenClaims["isSuperAdmin"])
		}
		//将从token中获取到的租户id存入tls中，用于动态数据源
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientID, tokenClaims[_const.ClientID].(string))
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientCode, tokenClaims[_const.ClientCode].(string))
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserId, tokenClaims[_const.UserId].(string))
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.UserCode, tokenClaims[_const.UserCode].(string))
		//将解析后的token中的用户信息存入local store
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.JwtData, tokenClaims)
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.TokenOriginal, token)
		ctx.Next()
	} else {
		ctx.JSON(result)
		return
	}
}

func parseResponseBody(respBody []byte) (rest.Result, error) {
	var result rest.Result
	//使用json解析响应体中的数据，并存入输出结构体中
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		clog.Errorf("token检查解析json到结构体出错 ", err)
		return rest.Result{}, err
	}
	return result, nil
}
