package databases

import (
	"context"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/middleware/http_session"
	"github.com/lm1996-mojor/go-core-library/store"
	"gorm.io/gorm"
)

// GetDbByName 根据key获取数据库操作对象
func GetDbByName(key string) (db *gorm.DB) {
	if key == "" {
		key = config.Sysconfig.DataBases.MasterDbName
	}
	return dbMap[key].WithContext(context.Background())
}

// ---------- 自定义数据源处理代码块 ----------------

func GetCustomDbTxByDbName(ctx iris.Context, name string) (tx *gorm.DB) {
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.CustomTx)
	if ok {
		tx = value.(*gorm.DB)
	} else {
		tx = GetDbByName(name).Begin()
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.CustomTx, tx)
	}
	return
}

// DisposeCustomizedTx commit the transaction if err is nil otherwise rollback
//func DisposeCustomizedTx(ctx iris.Context, err interface{}) {
//	transaction(ctx, _const.CustomTx, err)
//}

// ---------- 主数据源处理代码块 ----------------

func GetMasterDbTx(ctx iris.Context) (tx *gorm.DB) {
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.MasterTx)
	if ok {
		tx = value.(*gorm.DB)
	} else {
		tx = GetDbByName("").Begin()
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.MasterTx, tx)
	}
	return
}

//func DisposeMasterDbTx(ctx iris.Context, err interface{}) {
//	transaction(ctx, _const.MasterTx, err)
//}

// ---------- 租户数据源处理代码块 ----------------

func GetClientDbTX(ctx iris.Context, clientId string) (tx *gorm.DB) {
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.ClientTx)
	if ok {
		tx = value.(*gorm.DB)
	} else {
		tx = GetDbByName(clientId).Begin()
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientTx, tx)
	}
	return
}

// ---------- 获取指定的数据源，且每个数据源占有独立的本地存储空间 ----------------

func GetDbTxObjByDbName(ctx iris.Context, name string) (tx *gorm.DB) {
	value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + _const.ClientTx + name)
	if ok {
		tx = value.(*gorm.DB)
	} else {
		tx = GetDbByName(name).Begin()
		store.Set(http_session.GetCurrentHttpSessionUniqueKey(ctx)+_const.ClientTx+name, tx)
	}
	return
}

// DisposeClientTx commit the transaction if err is nil otherwise rollback
//func DisposeClientTx(ctx iris.Context, err interface{}) {
//	transaction(ctx, _const.ClientTx, err)
//}

func TransactionHandler(ctx iris.Context, err interface{}) {
	values := store.GetValueByCondition(fmt.Sprint(&ctx))
	values.Range(func(key, value any) bool {
		if strings.Contains(key.(string), "_db_") {
			tx := value.(*gorm.DB)
			if err == nil {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}
		return true
	})
	//txObjKey := ""
	//switch dbType {
	//case _const.ClientTx:
	//	txObjKey = _const.ClientTx
	//case _const.MasterTx:
	//	txObjKey = _const.MasterTx
	//default:
	//	txObjKey = _const.CustomTx
	//}
	//// 获取到单次会话获取过的数据库操作对象
	//value, ok := store.Get(http_session.GetCurrentHttpSessionUniqueKey(ctx) + txObjKey)
	//if ok {
	//	tx := value.(*gorm.DB)
	//	if err == nil {
	//		tx.Commit()
	//	} else {
	//		tx.Rollback()
	//	}
	//	store.Del(http_session.GetCurrentHttpSessionUniqueKey(ctx) + txObjKey)
	//}
}
