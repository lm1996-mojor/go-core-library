package consul

import (
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
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
	log.Info("服务健康检查")
	CheckFlag = true
	spec := "@every 1s"
	err := tasker_factory.AddTask("RestartCheckFlag", "重置检查目标定时任务", spec, RestartCheckFlag)
	if err != nil {
		panic("添加本地服务状态检查定时任务添加失败" + err.Error())
	}
	return rest.SuccessResult(nil)
}
