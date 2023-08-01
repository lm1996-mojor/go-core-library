package repo

import (
	"errors"
	"fmt"
	"strconv"

	_const "github.com/lm1996-mojor/go-core-library/const"
	dbLib "github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/store"

	"gorm.io/gorm"
)

//---------------------------- 自定义数据源工具方法部分 ---------------------------------

// ObtainCustomDbByDbName 根据自定义的数据源名称获取自定义数据源对象
func ObtainCustomDbByDbName(dbName string) (db *gorm.DB) {
	return dbLib.GetCustomizedDbByName(dbName)
}

// ObtainCustomTxDbByDbName 根据自定义的数据源名称获取带事务的自定义数据源对象
func ObtainCustomTxDbByDbName(dbName string) (tx *gorm.DB) {
	return dbLib.GetCustomDbTxByDbName(dbName)
}

// ObtainMasterDb 获取常规主数据源
func ObtainMasterDb() (db *gorm.DB) {
	return dbLib.GetMasterDb()
}

// ObtainMasterDbTx 获取带事务的数据源
func ObtainMasterDbTx() (tx *gorm.DB) {
	return dbLib.GetMasterDbTx()
}

// ObtainClientDb 获取常规动态租户数据源
func ObtainClientDb() (db *gorm.DB) {
	clientId, err := obtainClientId()
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetClientDb(fmt.Sprintf("%d", clientId))
}

// ObtainClientDbTx 获取带事务的动态租户数据源
func ObtainClientDbTx() (db *gorm.DB) {
	clientId, err := obtainClientId()
	if err != nil {
		log.Error("租户id获取失败，请检查token情况，和本地缓存情况" + err.Error())
		panic("服务器错误")
	}
	return dbLib.GetClientDbTX(fmt.Sprintf("%d", clientId))
}

// ObtainClientId 获取当前的租户id
func obtainClientId() (clientId int64, err error) {
	value, ok := store.Get(fmt.Sprintf("%p", &store.PoInterKey) + _const.ClientID)
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

/**
-----------------  老版架构数据库工具分界线---------------------
repo_util所需方法
*/

type OrderType int

const (
	Desc OrderType = iota
	Asc
)

type SortBo struct {
	FiledName string
	Order     OrderType
}

const DeletedAtFilterSql = "deleted_at IS NULL"

func (s *SortBo) ConvertSortString() string {
	return fmt.Sprintf("%s %s", s.FiledName, s.Order.ConvertString())
}

func (o OrderType) ConvertString() string {
	switch o {
	case Desc:
		return "desc"
	case Asc:
		return "asc"
	}
	return "desc"
}
