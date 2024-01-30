package consul

import (
	"strings"

	"github.com/hashicorp/consul/api"
	libConfig "github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/utils/sys_environment"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

func Register() string {
	log.Info("服务注册中...")
	protocol := "://"
	host := ""
	ipAddrList := sys_environment.GetInternalIP()
	if libConfig.Sysconfig.SystemEnv.Env != "prod" {
		protocol = "http" + protocol
		host = ipAddrList[0]
	} else {
		protocol = "https" + protocol
		host = sys_environment.GetExternal()
	}
	serviceCheck := &api.AgentServiceCheck{
		HTTP:                           protocol + host + ":" + libConfig.Sysconfig.App.Port + "/consul/ser/health",
		Timeout:                        libConfig.Sysconfig.Consul.Check.CheckTimeout,
		Interval:                       libConfig.Sysconfig.Consul.Check.CheckInterval,
		DeregisterCriticalServiceAfter: libConfig.Sysconfig.Consul.Check.InvalidServiceLogoutTime,
	}
	meta := make(map[string]string)
	for i := 0; i < len(ipAddrList); i++ {
		if i == len(ipAddrList)-1 {
			meta["intranet"] = meta["intranet"] + ipAddrList[i] + ":" + libConfig.Sysconfig.App.Port
		} else {
			meta["intranet"] = meta["intranet"] + ipAddrList[i] + ":" + libConfig.Sysconfig.App.Port + ","
		}
	}
	meta["public_network"] = sys_environment.GetExternal() + ":" + libConfig.Sysconfig.App.Port
	registration := &api.AgentServiceRegistration{
		Address: host,
		ID:      libConfig.Sysconfig.App.Name + "-" + host + ":" + libConfig.Sysconfig.App.Port + "-" + strings.Split(uuid.NewV3(uuid.NewV4(), "abc").String(), "-")[0],
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
