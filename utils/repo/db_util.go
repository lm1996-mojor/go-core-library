package repo

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
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
	return dbLib.GetDbByName(dbName)
}

// ObtainCustomTxDbByDbName 根据自定义的数据源名称获取带事务的自定义数据源对象
func ObtainCustomTxDbByDbName(ctx iris.Context, dbName string) (tx *gorm.DB) {
	return dbLib.GetCustomDbTxByDbName(ctx, dbName)
}

// ObtainMasterDb 获取常规主数据源
func ObtainMasterDb() (db *gorm.DB) {
	return dbLib.GetDbByName("")
}

// ObtainMasterDbTx 获取带事务的数据源
func ObtainMasterDbTx(ctx iris.Context) (tx *gorm.DB) {
	return dbLib.GetMasterDbTx(ctx)
}

// ObtainClientDb 获取常规动态租户数据源
func ObtainClientDb(ctx iris.Context) (db *gorm.DB) {
	clientId, err := ObtainClientId(ctx)
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetDbByName(fmt.Sprintf("%d", clientId))
}

// ObtainClientDbTx 获取带事务的动态租户数据源
func ObtainClientDbTx(ctx iris.Context) (db *gorm.DB) {
	clientId, err := ObtainClientId(ctx)
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetClientDbTX(ctx, fmt.Sprintf("%d", clientId))
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

// ObtainDbObjByDbName 根据名称获取独立存储空间的db对象
func ObtainDbObjByDbName(dbName string) (db *gorm.DB) {
	return dbLib.GetDbByName(dbName)
}

// ObtainDbTxObjByDbName 根据名称获取独立存储空间且带事务的db对象,该方法仅限用于一次请求需要操作多个数据源的场景
func ObtainDbTxObjByDbName(ctx iris.Context, dbName string) (tx *gorm.DB) {
	return dbLib.GetDbTxObjByDbName(ctx, dbName)
}
