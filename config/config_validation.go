package config

func validation() {
	dbValidation()
}

func dbValidation() {
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
}
