package tasker_factory

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var TaskMap sync.Map
var taskKeys = make([]string, 0) // 用于存储map的key用于顺序遍历

// Task
/** 定时任务对象
 *  使用的cron/v3包
 *  相关资料：https://www.cnblogs.com/beatle-go/p/17472105.html
 */
type Task struct {
	TaskBody   *cron.Cron    // 任务主体对象
	TaskId     cron.EntryID  // 任务启动后的id
	TaskDesc   string        // 任务描述
	TaskStatus bool          // 当前任务状态
	TaskJob    func()        // 定时任务执行的方法
	TaskOpt    []cron.Option // 自定义的负载
}
