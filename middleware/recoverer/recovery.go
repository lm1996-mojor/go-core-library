package recoverer

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"
)

// Recover 统一错误处理中心
func Recover(ctx iris.Context) {
	defer func() {
		err := recover()
		//databases.DisposeCustomizedTx(ctx, err)
		//databases.DisposeMasterDbTx(ctx, err)
		//databases.DisposeClientTx(ctx, err)
		databases.TransactionHandler(err)
		if err != nil {
			log.Error("服务器错误：" + fmt.Sprint(err))
			ctx.JSON(rest.FailCustom(500, fmt.Sprint(err), rest.ERROR))
		}
		store.DelCurrent(fmt.Sprint(&ctx))
	}()
	ctx.Next()
}
