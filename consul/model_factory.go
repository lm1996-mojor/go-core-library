package consul

import (
	"fmt"

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
