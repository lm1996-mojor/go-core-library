package consul

import (
	"sort"
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

var ServiceLib []ServiceLibrary

type ServiceLibrary struct {
	ServiceName     string
	ServiceId       string
	ServiceMetadata map[string]string
	Host            string
	Port            int
	Proto           string
	Weight          int
}

func FindSpecifyingServiceList(serviceName string) (serviceList []ServiceLibrary) {
	if strings.Contains(serviceName, "/") {
		serviceName = strings.ReplaceAll(serviceName, "/", "")
	}
	for _, service := range ServiceLib {
		if service.ServiceName == serviceName {
			serviceList = append(serviceList, service)
		}
	}
	// 做排序操作：降序（将权重最高的服务放在前面）
	sort.Slice(serviceList, func(i, j int) bool {
		return serviceList[i].Weight > serviceList[j].Weight
	})
	return serviceList
}

func ObtainHighestWeightInServiceList(serviceName string) ServiceLibrary {
	return FindSpecifyingServiceList(serviceName)[0]
}

func Init(app *iris.Application) {
	if libConfig.Sysconfig.Consul.Addr != "" && libConfig.Sysconfig.Consul.Addr != "null" && len(libConfig.Sysconfig.Consul.Addr) > 0 {
		host := ""
		ipAddrList := sys_environment.GetInternalIP()
		if libConfig.Sysconfig.SystemEnv.Env != "prod" {
			host = strings.ReplaceAll(strings.Split(ipAddrList[0], "/")[0], ".", "_")
		} else {
			host = strings.ReplaceAll(strings.Split(sys_environment.GetExternal(), "/")[0], ".", "_")
		}
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
		mvc.New(app.Party("/consul")).Handle(NewController())
		log.Info("初始化服务治理-服务健康检查接口")
	} else {
		log.Info("无服务治理要求...")
	}
}

const runLevel = 10

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel, ControllerReqPath: "", IsAddContextPath: false})
}
