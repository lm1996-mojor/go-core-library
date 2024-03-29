package consul

import (
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/rest"
)

type Controller struct {
	Ctx iris.Context
}

func NewController() *Controller {
	return &Controller{}
}

// GetSerHealth
/** 服务健康检查*/
func (c *Controller) GetSerHealth() rest.Result {
	return rest.SuccessResult(nil)
}
