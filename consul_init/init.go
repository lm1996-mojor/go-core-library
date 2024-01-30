package consul_init

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	libConfig "github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/consul_utils"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/utils/sys_environment"
)

func Init(app *iris.Application) {
	mvc.New(app.Party("/consul_utils")).Handle(consul_utils.NewController())
	log.Info("初始化服务治理-服务健康检查接口")
}

const runLevel = 10

func init() {
	if libConfig.Sysconfig.Consul.Addr != "" && libConfig.Sysconfig.Consul.Addr != "null" && len(libConfig.Sysconfig.Consul.Addr) > 0 {
		host := ""
		ipAddrList := sys_environment.GetInternalIP()
		if libConfig.Sysconfig.SystemEnv.Env != "prod" {
			host = ipAddrList[0]
		} else {
			host = sys_environment.GetExternal()
		}
		searchConditionValue := "ID contains " + libConfig.Sysconfig.App.Name + "-" + host + ":" + libConfig.Sysconfig.App.Port

		consulServiceInfoList, err := consul_utils.FindServiceList(searchConditionValue)
		if err != nil {
			panic("查询服务列表失败：" + err.Error())
		}
		if len(consulServiceInfoList) > 0 {
			for serviceId, _ := range consulServiceInfoList {
				if err1 := consul_utils.ServiceDeregister(serviceId); err1 != nil {
					panic("服务注销失败：" + err1.Error())
				}
			}
		}
		consulServiceId := consul_utils.Register()
		store.Set(_const.ConsulEndId, consulServiceId)
		global.RegisterInit(global.Initiator{Action: Init, Level: runLevel, EndFlag: false})
	}
	log.Info("无服务治理要求...")
}