package tasker_factory

import (
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
)

const runLevel = -1

// 初始化数据库操作
func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

func Init(app *iris.Application) {
	log.Info("定时任务初始化....")
	go jobRun()
	log.Info("定时任务初始化完成")
}

func jobRun() {
	if config.Sysconfig.Consul.EnableObtainService {
		spec := "*/3 * * * *"
		if config.Sysconfig.Consul.Service.Spec != "" && config.Sysconfig.Consul.Service.Spec != "null" && len(config.Sysconfig.Consul.Service.Spec) > 0 {
			spec = config.Sysconfig.Consul.Service.Spec
		}
		c := GetCornTasker()
		_, err := c.AddFunc(spec, ObtainSpecifyingConfigServicesFromTheRegistrationCenter)
		if err != nil {
			panic("定时任务添加失败" + err.Error())
		}
		c.Run()
	}
}
