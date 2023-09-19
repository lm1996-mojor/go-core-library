package tasker_factory

import (
	"time"
)

// Tasker 定时任务对象
type Tasker struct {
	t            *time.Ticker  // 任务主要执行对象
	stopChan     chan struct{} // 停止通道
	TaskerStatus int8          // 任务状态（1 就绪 2 执行中 3 阻塞中 4 已停止）
	TaskerId     string        // 任务id
	TaskerName   string        // 任务名称
	OwnerObjName string        // 任务所属对象名称
	ExecTime     string        // 执行间隔时间（年-月-周-日&次）
}

func NewTasker(d time.Duration, taskerId string, taskerName string, ownerObjName string) *Tasker {
	return &Tasker{
		t:            time.NewTicker(d),
		stopChan:     make(chan struct{}),
		TaskerStatus: 1,
		TaskerId:     taskerId,
		TaskerName:   taskerName,
		OwnerObjName: ownerObjName,
	}
}

// Stop 停止任务
func (t *Tasker) stop() {
	t.TaskerStatus = 4
	t.stopChan <- struct{}{}
}

// Reset 重置任务：在指定时间点执行下次任务
func (t *Tasker) reset(execTime time.Duration) {
	t.t.Reset(execTime)
}
