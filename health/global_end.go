package health

import (
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/consul"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"
)

func ServiceEndGlobal() {
	if config.Sysconfig.Consul.Addr != "" && config.Sysconfig.Consul.Addr != "null" && len(config.Sysconfig.Consul.Addr) > 0 {
		value, ok := store.Get(_const.ConsulEndId)
		if ok {
			log.Error("获取本地缓存数据失败：consulId")
		}
		err := consul.ServiceDeregister(value.(string))
		if err != nil {
			log.Error("consul服务注销失败：" + err.Error())
		}
	}
}
