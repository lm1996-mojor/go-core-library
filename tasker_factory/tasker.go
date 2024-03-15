package tasker_factory

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var TaskMap sync.Map
var taskKeys = make([]string, 0) // 用于存储map的key用于顺序遍历

type Task struct {
	TaskBody   *cron.Cron   // 任务主体对象
	TaskId     cron.EntryID // 任务启动后的id
	TaskDesc   string       // 任务描述
	TaskStatus bool         // 当前任务状态
}
