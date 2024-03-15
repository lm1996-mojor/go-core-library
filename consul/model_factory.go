package consul

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
	localConfig "github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
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

var CheckFlag bool
var failCheckCount int
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

func TimedExecution() {
	spec := "@every 11s"
	var err error
	if localConfig.Sysconfig.Consul.EnableObtainService {
		if localConfig.Sysconfig.Consul.Service.Spec != "" && localConfig.Sysconfig.Consul.Service.Spec != "null" && len(localConfig.Sysconfig.Consul.Service.Spec) > 0 {
			spec = localConfig.Sysconfig.Consul.Service.Spec
		}
		err = tasker_factory.AddTask("ObtainSpecifyingConfigServicesFromTheRegistrationCenter", "发现服务定时任务", spec, ObtainSpecifyingConfigServicesFromTheRegistrationCenter)
		if err != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
	}
	err = tasker_factory.AddTask("LocalCheckServiceHealth", "本地服务状态检查定时任务", spec, LocalCheckServiceHealth)
	if err != nil {
		panic("添加本地服务状态检查定时任务添加失败" + err.Error())
	}
	err = tasker_factory.AddTask("RestartCheckFlag", "重置检查目标定时任务", spec, RestartCheckFlag)
	if err != nil {
		panic("添加本地服务状态检查定时任务添加失败" + err.Error())
	}
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

func LocalCheckServiceHealth() {
	log.Info("本地检查服务健康情况...")
	if failCheckCount >= 3 {
		// 重新注册服务
		log.Info("当前服务已经无法连接到服务中心3次，已停止重新注册服务相关的定时任务")
		err := tasker_factory.StopAndRemoveTask("LocalCheckServiceHealth")
		if err != nil {
			log.Error(err.Error())
		}
		err = tasker_factory.StopAndRemoveTask("RestartCheckFlag")
		if err != nil {
			log.Error(err.Error())
		}
		consulServiceId := Register()
		store.Set(_const.ConsulEndId, consulServiceId)
		spec := "@every 11s" // 需要根据配置文件中的checkInterval参数进行动态修改: 增加1秒
		err = tasker_factory.AddTask("LocalCheckServiceHealth", "本地服务状态检查定时任务", spec, LocalCheckServiceHealth)
		if err != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
		err = tasker_factory.AddTask("RestartCheckFlag", "重置检查目标定时任务", spec, RestartCheckFlag)
		if err != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
		tasker_factory.BatchedRunTasker()
		return
	}
	if !CheckFlag {
		failCheckCount++
	}
	failCheckCount = 0
	log.Info("检查完成...")
}

func RestartCheckFlag() {
	CheckFlag = false
}
