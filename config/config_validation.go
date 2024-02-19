package config

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func validation() {
	dbValidation()
	consulConfigValidate()
}

func dbValidation() {
	log.Info().Msg("数据源配置检查...")
	if Sysconfig.DataBases.MasterDbName == "" || len(Sysconfig.DataBases.MasterDbName) <= 0 {
		panic("Master 数据源不能为空，请检查yaml中的 masterDbName")
	}
	flag := false
	for _, dbInfo := range Sysconfig.DataBases.DbInfoList {
		if dbInfo.DbName == Sysconfig.DataBases.MasterDbName {
			flag = true
		}
	}
	if !flag {
		panic("Master 数据源，在多数据源列表(DbInfoList)中不存在，请检查yaml中的 masterDbName")
	}
	log.Info().Msg("数据源配置检查完成")
}

func consulConfigValidate() {
	log.Info().Msg("服务治理中心配置检查...")
	serviceListStr := ""
	for _, service := range Sysconfig.Consul.Service.DesignatedServices {
		serviceListStr += service.ServiceName + ","
	}
	if Sysconfig.Detection.Token || Sysconfig.Detection.Authentication {
		if !Sysconfig.Consul.EnableObtainService {
			panic("检测到开启了token或者权限检测，但是没有开启服务定时检索，请在yml文件中添加[consul.enableObtainService:true]")
		}
	}
	if Sysconfig.Detection.Token {
		if !strings.Contains(serviceListStr, Sysconfig.Detection.TokenService) {
			panic("检测到服务中开启了token检查，但是服务发现中没有对应的服务检索关键词，请在yml文件中的[consul.service.designatedServices]添加" + Sysconfig.Detection.TokenService)
		}
	}
	if Sysconfig.Detection.Authentication {
		if !strings.Contains(serviceListStr, Sysconfig.Detection.AuthService) {
			panic("检测到服务中开启了权限检查，但是服务发现中没有对应的服务检索关键词，请在yml文件中的[consul.service.designatedServices]添加" + Sysconfig.Detection.AuthService)
		}
	}
	log.Info().Msg("服务治理中心数据源配置检查完成")
}
