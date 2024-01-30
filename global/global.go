package global

import (
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/health"
	"github.com/lm1996-mojor/go-core-library/log"
)

type Initiator struct {
	Action  func(app *iris.Application)
	Level   int
	EndFlag bool
	id      string
}

var initiators []Initiator

func RegisterInit(initiator Initiator) {
	if initiators == nil {
		initiators = make([]Initiator, 0, 5)
	}
	initiator.id = uuid.New().String()
	if initiator.EndFlag {
		for _, info := range initiators {
			if info.id != initiator.id {
				if info.EndFlag && initiator.EndFlag {
					panic("系统仅限于全局存在一个末尾业务注册")
				}
			}
		}
	}
	initiators = append(initiators, initiator)
}

var mu sync.Mutex

func runInitiator(first, last int, app *iris.Application) {
	app.Configure(iris.WithConfiguration(iris.Configuration{
		TimeFormat: "2006-01-02 15:04:05",
	}))
	wg := sync.WaitGroup{}
	wg.Add(last - first + 1)
	for j := first; j < last+1; j++ {
		go func(j int) {
			if initiators[j].Level >= 10 {
				mu.Lock()
			}
			initiators[j].Action(app)
			if initiators[j].Level >= 10 {
				mu.Unlock()
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
}

func RunApp(app *iris.Application, level int, appConfigs ...iris.Configurator) {
	sort.Slice(initiators, func(i, j int) bool {
		return initiators[i].Level < initiators[j].Level
	})
	for i := 0; i < len(initiators); i++ {
		mu.Lock()
		initiators[i].Action(app)
		mu.Unlock()
	}
	defer func() {
		health.ServiceEndGlobal()
	}()
	var err error
	if len(appConfigs) > 0 {
		err = app.Run(iris.Addr(":"+config.Sysconfig.App.Port), appConfigs...)
	} else {
		err = app.Run(iris.Addr(":" + config.Sysconfig.App.Port))
	}
	if err != nil {
		log.Error("服务停止：" + err.Error())
		panic(err)
	}
	//first := 0
	//last := 0
	//for i, initiator := range initiators {
	//	if initiator.Level > level {
	//		runInitiator(first, last, app)
	//		break
	//	}
	//	if initiator.Level == initiators[first].Level {
	//		last = i
	//		if i < len(initiators)-1 {
	//			continue
	//		} else {
	//			runInitiator(first, last, app)
	//			break
	//		}
	//	}
	//
	//	runInitiator(first, last, app)
	//
	//	first = i
	//	last = i
	//
	//	if i == len(initiators)-1 {
	//		runInitiator(first, last, app)
	//	}
	//}
}
