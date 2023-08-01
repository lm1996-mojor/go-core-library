package recoverer

import (
	"github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"

	"github.com/kataras/iris/v12"
)

// Recover middleware to recover the transaction
// 统一事务处理,事务自动提交
func Recover(ctx iris.Context) {
	defer func() {
		err := recover()
		databases.DisposeCustomizedTx(err)
		databases.DisposeMasterDbTx(err)
		databases.DisposeClientTx(err)
		if err != nil {
			log.Error("服务器错误：" + err.(error).Error())
			ctx.JSON(rest.FailCustom(500, err.(error).Error(), rest.ERROR))
		}
		store.Clean()
	}()
	ctx.Next()
}
