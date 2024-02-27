package consul

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
	localConfig "github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/log"
)

func GetClient() *api.Client {
	config := api.DefaultConfig()
	config.Address = localConfig.Sysconfig.Consul.Addr + ":" + fmt.Sprint(localConfig.Sysconfig.Consul.Port)
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client
}

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
		if strings.Contains(service.ServiceName, serviceName) {
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

func ObtainSpecifyingConfigServicesFromTheRegistrationCenter() {
	log.Info("获取指定服务列表...")
	ServiceLib = make([]ServiceLibrary, 0)
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
	if len(localConfig.Sysconfig.Consul.Service.DesignatedServices) <= 0 {
		ServiceLib = serviceList
	} else {
		for _, list := range serviceList {
			for _, designatedService := range localConfig.Sysconfig.Consul.Service.DesignatedServices {
				if designatedService.ServiceName == list.ServiceName || localConfig.Sysconfig.Detection.AuthService == list.ServiceName || localConfig.Sysconfig.Detection.TokenService == list.ServiceName {
					ServiceLib = append(ServiceLib, list)
				}
			}
		}
	}

	log.Info("获取指定服务列表完成")
}
