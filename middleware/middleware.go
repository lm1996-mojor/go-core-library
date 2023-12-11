package middleware

import (
	"errors"
	"sort"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"
	clog "github.com/lm1996-mojor/go-core-library/log"
	cors "github.com/lm1996-mojor/go-core-library/middleware/cors_handler"
	"github.com/lm1996-mojor/go-core-library/middleware/recoverer"
	"github.com/lm1996-mojor/go-core-library/middleware/security/auth"
	"github.com/lm1996-mojor/go-core-library/middleware/security/auth/white_list"
	"github.com/lm1996-mojor/go-core-library/middleware/security/token"
	"github.com/lm1996-mojor/go-core-library/middleware/session_data_handler"

	"github.com/kataras/iris/v12/context"
)

const runLevel = 9

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

// Init 中间件初始化
func Init(app *iris.Application) {
	// 初始化白名单
	white_list.Init()
	// 注册中间件
	RegisterMiddleWare(app)
}

func RegisterMiddleWare(app *iris.Application) {
	app.Configure(iris.WithOptimizations)
	// 配置会话id
	app.UseGlobal(requestid.New(requestid.DefaultGenerator))
	// 根据情况选定使用哪个中间件（token 和 会话数据初始化）
	if !config.Sysconfig.Detection.Token {
		tempSlice := make([]MiddleWare, 0)
		for _, ware := range globalMiddleWares {
			if ware.HandlerEnDesc != "token_check" {
				tempSlice = append(tempSlice, ware)
			}
		}
		globalMiddleWares = tempSlice
	} else {
		tempSlice := make([]MiddleWare, 0)
		for _, ware := range globalMiddleWares {
			if ware.HandlerEnDesc != "session_data_init" {
				tempSlice = append(tempSlice, ware)
			}
		}
		globalMiddleWares = tempSlice
	}
	// 关闭权限检测
	if !config.Sysconfig.Detection.Authentication {
		tempSlice := make([]MiddleWare, 0)
		for _, ware := range globalMiddleWares {
			if ware.HandlerEnDesc != "auth_check" {
				tempSlice = append(tempSlice, ware)
			}
		}
		globalMiddleWares = tempSlice
	}

	// 配置跨域处理
	cors.InitCors(app)
	clog.Info("中间件中心注册中间件中.....")

	if len(globalMiddleWares) > 0 {
		// 按照等级排序: 升序
		sort.Slice(globalMiddleWares, func(i, j int) bool {
			return globalMiddleWares[i].MiddleWareLevel < globalMiddleWares[j].MiddleWareLevel
		})
		// 注册中间件
		for i := 0; i < len(globalMiddleWares); i++ {
			clog.Info(globalMiddleWares[i].HandlerCnDesc + "注册中...")
			app.UseGlobal(globalMiddleWares[i].Handler)
		}
	} else {
		clog.Info("没有全局中间件，无需处理")
	}
	if len(singleMiddleWares) > 0 {
		// 按照等级排序: 升序
		sort.Slice(singleMiddleWares, func(i, j int) bool {
			return singleMiddleWares[i].MiddleWareLevel < singleMiddleWares[j].MiddleWareLevel
		})
		for i := 0; i < len(singleMiddleWares); i++ {
			clog.Info(singleMiddleWares[i].HandlerCnDesc + "注册中...")
			app.UseGlobal(singleMiddleWares[i].Handler)
		}
	} else {
		clog.Info("没有个体中间件，无需处理")
	}
	clog.Info("中间件处理中心已加载完成.....")
}

// MiddleWare 中间件结构体
type MiddleWare struct {
	Handler         context.Handler // 中间件处理器
	HandlerCnDesc   string          // 中间件处理器描述
	HandlerEnDesc   string          // 中间件处理器英文描述
	HandlerServer   string          // 中间件所属服务，用于解决所属服务在使用公共库时。不会重复注册中间件。
	MiddleWareLevel int32           // 中间件等级(影响中间件运行顺序,数值越大，等级越小)
}

// 全局化web中间件，先于其他中间件执行
var globalMiddleWares = []MiddleWare{
	{recoverer.Recover, "统一错误处理(固定插件)", "err_recover", "global", 1},
	{session_data_handler.SessionDataInit, "会话数据初始化(固定插件)", "session_data_init", "global", 2},
	{token.CheckIdentity, "token检查(固定插件)", "token_check", "global", 99},
	{auth.Verify, "权限检查(固定插件)", "auth_check", "global", 100},
}

// web中间件，比global中间件晚运行
var singleMiddleWares = []MiddleWare{}

// AppendSingleMiddleWares 新增全局路由的中间件，比global中间件晚运行
func AppendSingleMiddleWares(item []MiddleWare) {
	for _, ware := range singleMiddleWares {
		for _, customWare := range item {
			if &ware.Handler == &customWare.Handler {
				clog.Error("请勿重复注册中间件,重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware repeat： please check the middleware"))
			}
			if ware.HandlerEnDesc == customWare.HandlerEnDesc {
				clog.Error("英文描述为检测项，请勿重复，重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware duplicate English description： please check the middleware"))
			}
			if ware.MiddleWareLevel == customWare.MiddleWareLevel {
				clog.Error("中间件为顺序运行，等级请勿重复，重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware duplicate level： please check the middleware"))
			}
		}
	}
	singleMiddleWares = append(singleMiddleWares, item...)
}

// AppendGlobalMiddleWares 新增全局化中间件，先于其他中间件执行
func AppendGlobalMiddleWares(item []MiddleWare) {
	for _, ware := range singleMiddleWares {
		for _, customWare := range item {
			if &ware.Handler == &customWare.Handler {
				clog.Error("请勿重复注册中间件,重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware repeat： please check the middleware"))
			}
			if ware.HandlerEnDesc == customWare.HandlerEnDesc {
				clog.Error("英文描述为检测项，请勿重复，重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware duplicate English description： please check the middleware"))
			}
			if ware.MiddleWareLevel == customWare.MiddleWareLevel {
				clog.Error("中间件为顺序运行，等级请勿重复，重复项:[" + ware.HandlerCnDesc + "] 与 [" + customWare.HandlerCnDesc + "]")
				panic(errors.New("middleware duplicate level： please check the middleware"))
			}
		}
	}
	globalMiddleWares = append(globalMiddleWares, item...)
}

func closeMiddleWares(srcData []MiddleWare, closeKey string) (returnData []MiddleWare) {
	index := -1
	for i := 0; i < len(srcData); i++ {
		if srcData[i].HandlerEnDesc == closeKey {
			index = i
		}
	}

}
