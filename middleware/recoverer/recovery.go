package recoverer

import (
	"fmt"

	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/store"

	"github.com/kataras/iris/v12"
)

// Recover 统一错误处理中心
func Recover(ctx iris.Context) {
	defer func() {
		err := recover()
		//databases.DisposeCustomizedTx(err)
		//databases.DisposeMasterDbTx(err)
		//databases.DisposeClientTx(err)
		if err != nil {
			log.Error("服务器错误：" + fmt.Sprint(err))
			ctx.JSON(rest.FailCustom(500, fmt.Sprint(err), rest.ERROR))
		}
		store.Clean()
	}()
	ctx.Next()
}
