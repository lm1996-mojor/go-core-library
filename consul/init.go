package consul

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	libConfig "github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/utils/sys_environment"
)

func Init(app *iris.Application) {
	if libConfig.Sysconfig.Consul.Addr != "" && libConfig.Sysconfig.Consul.Addr != "null" && len(libConfig.Sysconfig.Consul.Addr) > 0 {
		host := ""
		//ipAddrList := make([]string, 0)
		ipAddrList := ""
		if libConfig.Sysconfig.SystemEnv.Env != "prod" {
			ipAddrList = sys_environment.GetIp()
		} else {
			//ipAddrList = append(ipAddrList, sys_environment.GetExternal())
			ipAddrList = sys_environment.GetExternal()
		}
		log.Infof("%s", ipAddrList)
		//if strings.Contains(ipAddrList[0], "/") {
		//	ipAddrList[0] = strings.Split(ipAddrList[0], "/")[0]
		//}
		//host = strings.ReplaceAll(ipAddrList[0], ".", "_")
		if strings.Contains(ipAddrList, "/") {
			ipAddrList = strings.Split(ipAddrList, "/")[0]
		}
		host = strings.ReplaceAll(ipAddrList, ".", "_")
		searchConditionValue := "ID contains " + libConfig.Sysconfig.App.Name + "_" + host + "_" + libConfig.Sysconfig.App.Port
		consulServiceInfoList, err := FindServiceList(searchConditionValue)
		if err != nil {
			panic("查询服务列表失败：" + err.Error())
		}
		if len(consulServiceInfoList) > 0 {
			for serviceId, _ := range consulServiceInfoList {
				if err1 := ServiceDeregister(serviceId); err1 != nil {
					panic("服务注销失败：" + err1.Error())
				}
			}
		}
		consulServiceId := Register()
		store.Set(_const.ConsulEndId, consulServiceId)
		ObtainSpecifyingConfigServicesFromTheRegistrationCenter()
		mvc.New(app.Party("/consul")).Handle(NewController())
		log.Info("初始化服务治理-服务健康检查接口")
	} else {
		log.Info("无服务治理要求...")
	}
}

const runLevel = 10

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}
