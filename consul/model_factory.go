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
	task := tasker_factory.InitTask()
	spec := "@every 11s"
	if localConfig.Sysconfig.Consul.EnableObtainService {
		if localConfig.Sysconfig.Consul.Service.Spec != "" && localConfig.Sysconfig.Consul.Service.Spec != "null" && len(localConfig.Sysconfig.Consul.Service.Spec) > 0 {
			spec = localConfig.Sysconfig.Consul.Service.Spec
		}

		OSCSFTRCTaskId, err := task.TaskBody.AddFunc(spec, ObtainSpecifyingConfigServicesFromTheRegistrationCenter)
		if err != nil {
			panic("添加发现服务定时任务添加失败" + err.Error())
		}
		task.TaskId = OSCSFTRCTaskId
		task.TaskDesc = "发现服务定时任务"
		tasker_factory.TaskMap["ObtainSpecifyingConfigServicesFromTheRegistrationCenter"] = task

	}
	task = tasker_factory.InitTask()
	LCSHTaskId, err1 := task.TaskBody.AddFunc(spec, LocalCheckServiceHealth)
	if err1 != nil {
		panic("添加本地服务状态检查定时任务添加失败" + err1.Error())
	}
	task.TaskId = LCSHTaskId
	task.TaskDesc = "本地服务状态检查定时任务"
	tasker_factory.TaskMap["LocalCheckServiceHealth"] = task

	task = tasker_factory.InitTask()
	RCFTaskId, err2 := task.TaskBody.AddFunc(spec, RestartCheckFlag)
	if err2 != nil {
		panic("添加重置检查目标定时任务添加失败" + err2.Error())
	}
	task.TaskId = RCFTaskId
	task.TaskDesc = "重置检查目标定时任务"
	tasker_factory.TaskMap["RestartCheckFlag"] = task
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
	var FailCheckCount int
	log.Info("本地检查服务健康情况...")
	if FailCheckCount >= 3 {
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
		task := tasker_factory.InitTask()
		LCSHTaskId, err1 := task.TaskBody.AddFunc(spec, LocalCheckServiceHealth)
		if err1 != nil {
			panic("添加本地服务状态检查定时任务添加失败" + err.Error())
		}
		task.TaskId = LCSHTaskId
		task.TaskDesc = "本地服务状态检查定时任务"
		tasker_factory.TaskMap["LocalCheckServiceHealth"] = task

		task = tasker_factory.InitTask()
		RCFTaskId, err2 := task.TaskBody.AddFunc(spec, RestartCheckFlag)
		if err2 != nil {
			panic("添加重置检查目标定时任务添加失败" + err.Error())
		}
		task.TaskId = RCFTaskId
		task.TaskDesc = "重置检查目标定时任务"
		tasker_factory.TaskMap["RestartCheckFlag"] = task
		tasker_factory.BatchedRunTasker()
		return
	}
	if !CheckFlag {
		FailCheckCount++
	}
	FailCheckCount = 0
	log.Info("检查完成...")
}

func RestartCheckFlag() {
	CheckFlag = false
}
