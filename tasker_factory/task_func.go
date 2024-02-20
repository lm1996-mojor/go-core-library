package tasker_factory

import (
	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/consul"
	"github.com/lm1996-mojor/go-core-library/log"
)

func ObtainSpecifyingConfigServicesFromTheRegistrationCenter() {
	log.Info("获取指定服务列表...")
	consul.ServiceLib = make([]consul.ServiceLibrary, 0)
	serviceMap, err := consul.FindServiceList("")
	if err != nil {
		log.Error("Failed to get model list：" + err.Error())
	}
	serviceList := make([]consul.ServiceLibrary, 0)
	for _, service := range serviceMap {
		serviceList = append(serviceList, consul.ServiceLibrary{
			ServiceName:     service.Service,
			ServiceId:       service.ID,
			ServiceMetadata: service.Meta,
			Host:            service.Address,
			Port:            service.Port,
			Proto:           "http",
			Weight:          1, // TODO： 后面要动态更改当前的权重情况
		})

	}
	if len(config.Sysconfig.Consul.Service.DesignatedServices) <= 0 {
		consul.ServiceLib = serviceList
	} else {
		for _, list := range serviceList {
			for _, designatedService := range config.Sysconfig.Consul.Service.DesignatedServices {
				if designatedService.ServiceName == list.ServiceName || config.Sysconfig.Detection.AuthService == list.ServiceName || config.Sysconfig.Detection.TokenService == list.ServiceName {
					consul.ServiceLib = append(consul.ServiceLib, list)
				}
			}
		}
	}

	log.Info("获取指定服务列表完成")
}
