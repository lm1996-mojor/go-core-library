package consul

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
	localConfig "github.com/lm1996-mojor/go-core-library/config"
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
