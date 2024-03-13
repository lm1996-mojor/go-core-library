package global

import (
	"fmt"
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

func Init(app *iris.Application, level int) {
	app.Configure(iris.WithConfiguration(iris.Configuration{
		TimeFormat: "2006-01-02 15:04:05",
	}))
	sort.Slice(initiators, func(i, j int) bool {
		return initiators[i].Level < initiators[j].Level
	})
	for i := 0; i < len(initiators); i++ {
		mu.Lock()
		initiators[i].Action(app)
		mu.Unlock()
	}
	tuAn()
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

func tuAn() {
	tuan := "....................阿弥陀佛....................\n" +
		"                    _oo0oo_                              \n" +
		"                   o8888888o                             \n" +
		"                   88“ 卍 ”88                             \n" +
		"                   (| -_- |                             \n" +
		"                   0\\  =  /0                            \n" +
		"                 __/ ‘---***’ \\__                          \n" +
		"              .' \\|        |/ '.                        \n" +
		"             / \\\\|||   :   |||// \\                    \n" +
		"            / _|||||  -卍-  |||||_\\                     \n" +
		"            |   |\\\\\\   -   /// | |                    \n" +
		"            | \\_|  ‘’\\---/     |  |                    \n" +
		"           \\ . -\\_    '-'    _/-. /                    \n" +
		"         __ ' . .'       /--.--\\  '. .'__               \n" +
		"       “”   ‘<   '.__\\_<|>_/__.'  “”                    \n" +
		"    |  | :  '- \\'.;' \\ _ /'  ;.'/ - ' : |   |          \n" +
		"     \\  \\   ’-\\_  __  \\ / __ _/ .-'   /  /           \n" +
		"  ===== ’-.     \\_ _ \\____/ __. -'  __. -' =====       \n" +
		"                    '=---='                              \n" +
		"...............佛祖保佑，永无BUG...............\n" +
		"------------------------------------------------------------------------\n"
	fmt.Println(tuan)
	//"年复一年春光度，度得他人做老板；老板扣我薄酒钱，没有酒钱怎过年。\n" +
	//	"春光逝去皱纹起，作起程序也委靡；来到水源把水灌，打死不做程序员。\n" +
	//	"别人笑我忒疯癫，我笑他人命太贱；状元三百六十行，偏偏来做程序员。\n" +
	//	"但愿老死电脑间，不愿鞠躬老板前；奔驰宝马贵者趣，公交自行程序员。\n" +
	//	"别人笑我忒疯癫，我笑自己命太贱；不见满街漂亮妹，哪个归得程序员。")
}
