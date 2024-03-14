package tasker_factory

import (
	"github.com/robfig/cron/v3"
)

var TaskMap = make(map[string]Task)

type Task struct {
	TaskBody   *cron.Cron   // 任务主体对象
	TaskId     cron.EntryID // 任务启动后的id
	TaskDesc   string       // 任务描述
	TaskStatus bool         // 当前任务状态
}
