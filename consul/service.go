package consul

import (
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	libConfig "github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
	"github.com/lm1996-mojor/go-core-library/utils/sys_environment"
	"github.com/spf13/cast"
)

func Register() string {
	log.Info("服务注册中...")
	protocol := "http://"
	host := ""
	ipAddrList := sys_environment.GetInternalIP()
	networkSegment := strings.Split(libConfig.Sysconfig.Consul.Addr, ".")[2]
	host = ""
	if len(ipAddrList) > 1 {
		for _, ip := range ipAddrList {
			if networkSegment == "127.0.0.1" || networkSegment == "localhost" {
				host = ip
				break
			}
			if strings.Split(ip, ".")[2] == networkSegment {
				host = ip
				break
			}
		}
	} else {
		host = ipAddrList[0]
	}

	if host == "" {
		panic("您当前计算机所处的网段，consul无法连接，当前consul的网络ip为：" + libConfig.Sysconfig.Consul.Addr + "，请将您当前的计算机所处网络，与consul同步")
	}
	// 解决服务器不在同一个网段的问题
	//if libConfig.Sysconfig.SystemEnv.Env != "prod" {
	//	protocol = "http" + protocol
	//	host = strings.Split(ipAddrList[0], "/")[0]
	//	//host = strings.ReplaceAll(strings.Split(ipAddrList[0], "/")[0], ".", "_")
	//} else {
	//	protocol = "https" + protocol
	//	if strings.Contains(sys_environment.GetExternal(), "/") {
	//		host = strings.Split(sys_environment.GetExternal(), "/")[0]
	//	} else {
	//		host = sys_environment.GetExternal()
	//	}
	//	//host = strings.ReplaceAll(strings.Split(sys_environment.GetExternal(), "/")[0], ".", "_")
	//}
	serviceCheck := &api.AgentServiceCheck{
		HTTP:                           protocol + host + ":" + libConfig.Sysconfig.App.Port + "/consul/ser/health",
		Timeout:                        libConfig.Sysconfig.Consul.Check.CheckTimeout,
		Interval:                       libConfig.Sysconfig.Consul.Check.CheckInterval,
		DeregisterCriticalServiceAfter: libConfig.Sysconfig.Consul.Check.InvalidServiceLogoutTime,
	}
	meta := make(map[string]string)
	for i := 0; i < len(ipAddrList); i++ {
		ipHost := ""
		if strings.Contains(ipAddrList[i], "/") {
			ipHost = strings.Split(ipAddrList[i], "/")[0]
		} else {
			ipHost = ipAddrList[i]
		}
		if i == len(ipAddrList)-1 {
			meta["intranet"] = meta["intranet"] + protocol + ipHost + ":" + libConfig.Sysconfig.App.Port
		} else {
			meta["intranet"] = meta["intranet"] + protocol + ipHost + ":" + libConfig.Sysconfig.App.Port + ","
		}
	}
	meta["public_network"] = protocol + sys_environment.GetExternal() + ":" + libConfig.Sysconfig.App.Port
	meta["protocol_host"] = protocol + host
	uuId, _ := uuid.GenerateUUID()
	registration := &api.AgentServiceRegistration{
		Address: host,
		ID:      libConfig.Sysconfig.App.Name + "_" + strings.ReplaceAll(host, ".", "_") + "_" + libConfig.Sysconfig.App.Port + "_" + strings.Split(uuId, "-")[0],
		Name:    libConfig.Sysconfig.App.Name,
		Port:    cast.ToInt(libConfig.Sysconfig.App.Port),
		Tags:    []string{libConfig.Sysconfig.App.Name},
		Check:   serviceCheck,
		Meta:    meta,
	}
	err := GetClient().Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	log.Info("服务已注册...")
	return registration.ID
}

func RefreshRegister() {
	log.Info("重新注册中...")
	consulServiceId := Register()
	store.Set(_const.ConsulEndId, consulServiceId)
	log.Info("重新注册完成")
}

func FindServiceList(searchConditionValue string) (map[string]*api.AgentService, error) {
	return GetClient().Agent().ServicesWithFilter(searchConditionValue)
}

func ServiceDeregister(serviceId string) error {
	return GetClient().Agent().ServiceDeregister(serviceId)
}

func FindSpecifyingServiceList(serviceName string) (serviceList []ServiceLibrary) {
	if strings.Contains(serviceName, "/") {
		serviceName = strings.ReplaceAll(serviceName, "/", "")
	}
	for _, service := range ServiceLib {
		if strings.Contains(service.ServiceName, serviceName) {
			serviceList = append(serviceList, service)
		}
	}
	if len(serviceName) <= 0 {
		log.Error("没有找到对应服务器")
		panic("服务器错错误")
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

func TimedExecution() {
	spec := "@every 11s"
	var err error
	if libConfig.Sysconfig.Consul.EnableObtainService {
		if libConfig.Sysconfig.Consul.Service.Spec != "" && libConfig.Sysconfig.Consul.Service.Spec != "null" && len(libConfig.Sysconfig.Consul.Service.Spec) > 0 {
			spec = libConfig.Sysconfig.Consul.Service.Spec
		}
		err = tasker_factory.AddTask("ObtainSpecifyingConfigServicesFromTheRegistrationCenter", "发现服务定时任务", spec, ObtainSpecifyingConfigServicesFromTheRegistrationCenter)
		if err != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
	}
	if libConfig.Sysconfig.Consul.EnableServRegister && libConfig.Sysconfig.Consul.Addr != "" && libConfig.Sysconfig.Consul.Port > 0 {
		err = tasker_factory.AddTask("ServerStatusCheck", "服务状态定时检查", spec, ServerStatusCheck)
		if err != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
	}
}

func ServerStatusCheck() {
	log.Info("服务状态检查中...")
	value, ok := store.Get(_const.ConsulEndId)
	if !ok {
		RefreshRegister()
	}
	filter, err := GetClient().Agent().ChecksWithFilter("ServiceID == " + value.(string))
	if err != nil {
		log.Error("Failed to get check list：" + err.Error())
		tasker_factory.StopTask("ServerStatusCheck")
		return
	}
	if _, flag := filter["service:"+value.(string)]; !flag {
		RefreshRegister()
	}

	log.Info("服务状态检查完成")
}

func ObtainSpecifyingConfigServicesFromTheRegistrationCenter() {
	log.Info("获取指定服务列表...")

	serviceMap, err := FindServiceList("")
	if err != nil {
		log.Error("Failed to get model list：" + err.Error())
	}
	serviceList := make([]ServiceLibrary, 0)
	for _, service := range serviceMap {
		serviceList = append(serviceList, ServiceLibrary{
			ServiceName:     service.Service,
			ServiceId:       service.ID,
			ServiceMetadata: service.Meta,
			Host:            service.Address,
			Port:            service.Port,
			Proto:           "http",
			Weight:          1, // TODO： 后面要动态更改当前的权重情况
		})
	}

	if len(libConfig.Sysconfig.Consul.Service.DesignatedServices) <= 0 {
		ServiceLib = serviceList
	} else {
		services := make([]ServiceLibrary, 0)
		designatedServiceMap := make(map[string]string)
		for _, designatedService := range libConfig.Sysconfig.Consul.Service.DesignatedServices {
			designatedServiceMap[designatedService.ServiceName] = designatedService.ServiceName
		}
		for _, list := range serviceList {
			_, ok := designatedServiceMap[list.ServiceName]
			if ok {
				services = append(services, list)
			}
		}
		ServiceLib = services
	}
	log.Info("获取指定服务列表完成")
}
