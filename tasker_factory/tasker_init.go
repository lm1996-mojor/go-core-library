package tasker_factory

import (
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
)

const runLevel = 99

// 初始化数据库操作
func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

func Init(app *iris.Application) {
	log.Info("定时任务初始化....")
	BatchedRunTasker()
	log.Info("定时任务初始化完成")
}
