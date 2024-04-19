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
	if len(Sysconfig.DataBases.DbInfoList) > 0 {
		for i, dbinfo := range Sysconfig.DataBases.DbInfoList {
			for j, dbInfo2 := range Sysconfig.DataBases.DbInfoList {
				if dbinfo.DbName == dbInfo2.DbName && dbinfo.DbType == dbInfo2.DbType && i != j {
					panic("自定义数据源列表中，不能存在两个相同的数据源(名称和类型相同):[" + dbinfo.DbName + "]")
				}
				if dbinfo.DbName == dbInfo2.DbName && dbinfo.DbType != dbInfo2.DbType && i != j {
					log.Warn().Msg("请注意！有两个名称相同，类型不同的数据源。请注意使用:[" + dbinfo.DbName + "]")
				}
			}
		}
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
	} else {
		log.Info().Msg("无数据源要求")
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
	log.Info().Msg("服务治理中心配置检查完成")
}
