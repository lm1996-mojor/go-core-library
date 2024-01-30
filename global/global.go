package global

import (
	"sort"
	"sync"

	"github.com/kataras/iris/v12"
)

type Initiator struct {
	Action func(app *iris.Application)
	Level  int
}

var initiators []Initiator

func RegisterInit(initiator Initiator) {
	if initiators == nil {
		initiators = make([]Initiator, 0, 5)
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

func RunApp(app *iris.Application, level int) {
	sort.Slice(initiators, func(i, j int) bool {
		return initiators[i].Level < initiators[j].Level
	})

	first := 0
	last := 0
	for i, initiator := range initiators {
		if initiator.Level > level {
			runInitiator(first, last, app)
			break
		}
		if initiator.Level == initiators[first].Level {
			last = i
			if i < len(initiators)-1 {
				continue
			} else {
				runInitiator(first, last, app)
				break
			}
		}

		runInitiator(first, last, app)

		first = i
		last = i

		if i == len(initiators)-1 {
			runInitiator(first, last, app)
		}
	}
}
