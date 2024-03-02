package consul

import (
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	libConfig "github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/log"
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
	for _, ip := range ipAddrList {
		if strings.Split(ip, ".")[2] == networkSegment {
			host = strings.Split(ip, "/")[0]
		}
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

func FindServiceList(searchConditionValue string) (map[string]*api.AgentService, error) {
	return GetClient().Agent().ServicesWithFilter(searchConditionValue)
}

func ServiceDeregister(serviceId string) error {
	return GetClient().Agent().ServiceDeregister(serviceId)
}
