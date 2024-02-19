package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/consul"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/middleware/security/auth/white_list"
	"github.com/lm1996-mojor/go-core-library/proxy"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"
)

func Verify(ctx iris.Context) {
	// 判断是否需要鉴权(鉴权必须要有token)
	if ctx.Values().Get("pass_label").(string) == "Y" || ctx.Values().Get(http_session.GetCurrentHttpSessionUniqueKey(ctx)+"isSuperAdmin").(bool) {
		ctx.Next()
		return
	}
	reqUrl := ctx.Request().URL.Path
	if white_list.InList(reqUrl, 2) {
		log.Info("当前接口无需鉴权")
		ctx.Next()
		return
	}
	// 权限系统-鉴权路径
	authService := consul.ObtainHighestWeightInServiceList(config.Sysconfig.Detection.AuthService)
	url := authService.Proto + "://" + authService.Host + config.Sysconfig.Detection.AuthCheckServiceApiUrl
	actionUrl := url + "?reqUrl=" + reqUrl
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.TokenOriginal)
	if !ok {
		ctx.JSON(rest.FailCustom(401, "暂未登录，请重新登录", rest.ERROR))
		return
	}
	// 获取远程请求对象
	remoteReqMdl := proxy.ParametricConstructionOfRemoteReqMdl(nil, nil, actionUrl, "GET", true, value.(string))
	// 进行远程请求
	json, respErr := proxy.RequestAction(&remoteReqMdl, "asynchronous")
	if respErr != nil {
		log.Error("远程服务响应错误")
		panic("服务器错误")
	}
	// 响应code如果不为200 则打印对应消息
	if json["code"].(float64) != 200 {
		log.Error(json["msg"].(string))
		ctx.JSON(rest.FailCustom(int(json["code"].(float64)), json["msg"].(string), rest.ERROR))
		return
	}
	ctx.Next()
}
