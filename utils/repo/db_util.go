package repo

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	dbLib "github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/store"

	"gorm.io/gorm"
)

//---------------------------- 自定义数据源工具方法部分 ---------------------------------

// ObtainCustomDbByDbName 根据自定义的数据源名称获取自定义数据源对象
func ObtainCustomDbByDbName(dbName string) (db *gorm.DB) {
	return dbLib.GetDbByName(dbName, "mysql")
}

// ObtainCustomTxDbByDbName 根据自定义的数据源名称获取带事务的自定义数据源对象
func ObtainCustomTxDbByDbName(ctx iris.Context, dbName string) (tx *gorm.DB) {
	return dbLib.GetCustomDbTxByDbName(ctx, dbName, "mysql")
}

// ObtainMasterDb 获取常规主数据源
func ObtainMasterDb() (db *gorm.DB) {
	return dbLib.GetDbByName("", "mysql")
}

// ObtainMasterDbTx 获取带事务的数据源
func ObtainMasterDbTx(ctx iris.Context) (tx *gorm.DB) {
	return dbLib.GetMasterDbTx(ctx, "mysql")
}

// ObtainClientDb 获取常规动态租户数据源
func ObtainClientDb(ctx iris.Context) (db *gorm.DB) {
	clientId, err := ObtainClientId(ctx)
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetDbByName(fmt.Sprintf("%d", clientId), "mysql")
}

// ObtainClientDbTx 获取带事务的动态租户数据源
func ObtainClientDbTx(ctx iris.Context) (db *gorm.DB) {
	clientId, err := ObtainClientId(ctx)
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetClientDbTX(ctx, fmt.Sprintf("%d", clientId), "mysql")
}

// ObtainDb 获取数据源
//
// @Param ctx http会话对象
//
// @Param txFlag 是否获取带事务的数据源标识（true 是  false 否）
func ObtainDb(ctx iris.Context, txFlag bool) *gorm.DB {
	clientId, err := ObtainClientId(ctx)
	if config.Sysconfig.App.Name == "platform_management" {
		if txFlag {
			return dbLib.GetCustomDbTxByDbName(ctx, "platform_management", "mysql")
		} else {
			return dbLib.GetDbByName("platform_management", "mysql")
		}
	}
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		//panic("服务器错误")
	}
	// 判断是否需要进入租户库
	if clientId <= 0 {
		// 使用默认数据库（如果当前操作的是业务数据库，则都会报错）
		if txFlag {
			return dbLib.GetCustomDbTxByDbName(ctx, "platform_management", "mysql")
		} else {
			return dbLib.GetDbByName("platform_management", "mysql")
		}
	} else {
		if config.Sysconfig.SystemEnv.Env == "prod" && config.Sysconfig.DataBases.ClientEnable {
			clientIdStr := fmt.Sprintf("%d", clientId)
			if txFlag {
				return dbLib.GetClientDbTX(ctx, clientIdStr, "mysql")
			} else {
				return dbLib.GetDbByName(clientIdStr, "mysql")
			}
		} else {
			if txFlag {
				return dbLib.GetMasterDbTx(ctx, "mysql")
			} else {
				return dbLib.GetDbByName("", "mysql")
			}
		}
	}
}

// ObtainClientId 获取当前的租户id
func ObtainClientId(ctx iris.Context) (clientId int64, err error) {
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.ClientID)
	if !ok {
		return 0, errors.New("租户不确定")
	}
	//租户id参数类型转换(string -> int64)
	cId, err1 := strconv.ParseInt(value.(string), 10, 32)
	if err1 != nil {
		return 0, errors.New("租户参数转换失败")
	}
	return cId, nil
}

// ObtainDbByDbType 根据db类型获取数据源
func ObtainDbByDbType(ctx iris.Context, txFlag bool, dbType string) *gorm.DB {
	switch dbType {
	case "mysql":
		return ObtainDb(ctx, txFlag)
	case "clickhouse":
		clientId, err := ObtainClientId(ctx)
		if err != nil {
			log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		}
		if config.Sysconfig.SystemEnv.Env == "prod" && config.Sysconfig.DataBases.ClientEnable {
			clientIdStr := fmt.Sprintf("%d", clientId)
			if txFlag {
				return dbLib.GetClientDbTX(ctx, clientIdStr, "clickhouse")
			} else {
				return dbLib.GetDbByName(clientIdStr, "clickhouse")
			}
		} else {
			if txFlag {
				return ObtainCustomTxDbByDbNameAndDbType(ctx, "human_resource_management", "clickhouse")
			} else {
				return ObtainCustomDbByDbNameWithDbType("human_resource_management", "clickhouse")
			}
		}
	default:
		panic("无法识别的数据库类型[" + dbType + "]")
	}
}

// ObtainCustomDbByDbNameWithDbType 根据数据类型和自定义的数据源名称获取自定义数据源对象
func ObtainCustomDbByDbNameWithDbType(dbName string, dbType string) (db *gorm.DB) {
	return dbLib.GetDbByName(dbName, dbType)
}

// ObtainCustomTxDbByDbNameAndDbType 根据数据类型和自定义的数据源名称获取带事务的自定义数据源对象
func ObtainCustomTxDbByDbNameAndDbType(ctx iris.Context, dbName string, dbType string) (tx *gorm.DB) {
	return dbLib.GetCustomDbTxByDbName(ctx, dbName, dbType)
}
