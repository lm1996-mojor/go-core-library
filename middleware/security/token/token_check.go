package token

import (
	"encoding/json"
	"fmt"

	_const "github.com/lm1996-mojor/go-core-library/const"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/proxy"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"

	"github.com/kataras/iris/v12"
)

func CheckIdentity(ctx iris.Context) {
	//获取token
	author := ctx.GetHeader(_const.TokenName)

	//获取请求路径
	reqPath := ctx.Path()
	clog.Info("请求路径: " + reqPath)
	//clog.Info("token_util：" + author)
	/**
	token不存在
	判断token是否存在，如果不存在，
	则判断访问的接口是否存在白名单中
		如果存在则继续路由访问
		如果不存在则返回错误信息
	*/
	//ctx.Values().Set("pass_label", "N")
	//if requestfilter.InList(reqPath, 1) {
	//	ctx.Values().Set("pass_label", "Y")
	//	ctx.Next()
	//	return
	//} else if reflect.DeepEqual(author, "") || len(author) < 7 {
	//	if requestfilter.InList(reqPath, 2) {
	//		ctx.Next()
	//		return
	//	}
	//	ctx.JSON(rest.FailCustom(403, "尚未登录", rest.ERROR))
	//	return
	//}
	//
	///**
	//token存在
	//请求统一登录系统进行token验证并解析
	//	如果验证通过则进行后续操作
	//	如果验证不通过则返回错误信息
	//*/
	//去除token 头部信息
	token := author[7:]
	// 获取解析后的token信息
	respBody, err := proxy.GetParseToken(token, "http://192.168.31.113:9901/platform_inlet/sso/parse/token")
	if err != nil {
		panic("CheckIdentity->token解析: 响应出错")
		return
	}
	// 解析响应体中的数据
	userClaims, parseRespErr := parseResponseBody(respBody)
	if parseRespErr != nil {
		ctx.JSON(rest.FailCustom(500, parseRespErr.Error(), rest.ERROR))
		return
	}
	if int64(userClaims["code"].(float64)) == 200 {
		//判断自定义的token类型是否正确
		tokenClaims := userClaims["data"].(map[string]interface{})["parse_token"].(map[string]interface{})
		if t, ok := tokenClaims["token_type"].(string); ok && t != _const.TokenType { //不是access token
			clog.Info("令牌类型认证无效: " + err.Error())
			ctx.JSON(rest.FailCustom(403, "登录信息无效，请重新登录", rest.ERROR))
			return
		}
		if t, ok := tokenClaims["token_single"].(string); ok && t != _const.TokenSignature { //不是access token
			clog.Info("令牌签名认证无效: " + err.Error())
			ctx.JSON(rest.FailCustom(403, "登录信息无效，请重新登录", rest.ERROR))
			return
		}
		// 以下所有数据都会在单次回话完成后进行清空
		//将从token中获取到的租户id存入tls中，用于动态数据源
		store.Set(fmt.Sprintf("%p", &store.PoInterKey)+_const.ClientID, tokenClaims[_const.ClientID].(string))
		store.Set(fmt.Sprintf("%p", &store.PoInterKey)+_const.ClientCode, tokenClaims[_const.ClientCode].(string))
		store.Set(fmt.Sprintf("%p", &store.PoInterKey)+_const.UserId, tokenClaims[_const.UserId].(string))
		store.Set(fmt.Sprintf("%p", &store.PoInterKey)+_const.UserCode, tokenClaims[_const.UserCode].(string))
		//将解析后的token中的用户信息存入local store
		store.Set(fmt.Sprintf("%p", &store.PoInterKey)+_const.JwtData, tokenClaims["client_code"].(map[string]interface{}))
		ctx.Next()
	} else {
		ctx.JSON(rest.FailCustom(int(userClaims["code"].(float64)), userClaims["msg"].(string), rest.ERROR))
		return
	}
}

func parseResponseBody(respBody []byte) (map[string]interface{}, error) {
	var userClaims map[string]interface{}
	//使用json解析响应体中的数据，并存入输出结构体中
	err := json.Unmarshal(respBody, &userClaims)
	if err != nil {
		clog.Errorf("security_check.go 解析json到结构体出错 ", err)
		return nil, err
	}
	return userClaims, nil
}
