package tasker_factory

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
)

func GetTaskId(key string) cron.EntryID {
	return TaskMap[key].TaskId
}

func GetTaskBody(key string) *cron.Cron {
	return TaskMap[key].TaskBody
}

func InitTask(opts ...cron.Option) Task {
	var task Task
	if len(opts) > 0 {
		task.TaskBody = cron.New(opts...)
	} else {
		task.TaskBody = cron.New()
	}
	return task
}

func BatchedRunTasker() {
	for _, t := range TaskMap {
		if t.TaskStatus {
			continue
		}
		t.TaskStatus = true
		t.TaskBody.Run()
	}
}

func StopTask(key string) {
	t, ok := TaskMap[key]
	if ok {
		if t.TaskStatus {
			t.TaskStatus = false
			t.TaskBody.Stop()
			return
		} else {
			log.WarnF("当前任务已经停止，无需重复操作key:[%s:id:(%d)]", key, t.TaskId)
			return
		}
	}
	log.Warn("没有找到对应的定时任务，停止任务失败key:" + key)
}

func BatchedStopTask(keys []string) {
	for _, key := range keys {
		StopTask(key)
	}
}

func RegexpStopTask(r *regexp.Regexp) {
	for key, _ := range TaskMap {
		result := r.MatchString(key)
		if result {
			StopTask(key)
		}
	}
}

func StopAndRemoveTask(key string) error {
	StopTask(key)
	return RemoveTask(key)
}

func RemoveTask(key string) error {
	t := TaskMap[key]
	taskId := t.TaskId
	if t.TaskStatus {
		return errors.New("当前任务(" + key + ")id为：(" + cast.ToString(taskId) + ")正在执行中，请先停止任务。再进行移除")
	}
	t.TaskBody.Remove(taskId)
	delete(TaskMap, key)
	log.Info("指定任务已移除：(" + key + ")，任务id为：[" + cast.ToString(taskId) + "]")
	return nil
}

func BatchRemoveTask(keys []string) error {
	errMsg := "这些任务$正在执行中，请先停止任务。再进行移除"
	replErrStr := ""
	errFlag := false
	for i, key := range keys {
		err := RemoveTask(key)
		if err != nil {
			errFlag = true
			t := TaskMap[key]
			if i == len(keys)-1 {
				replErrStr += fmt.Sprintf("[taskKey:[%s,id:(%d)]]", key, t.TaskId)
			} else {
				replErrStr += fmt.Sprintf("[taskKey:[%s,id:(%d)]],", key, t.TaskId)
			}
		}
	}
	if errFlag {
		strings.ReplaceAll(errMsg, "$", replErrStr)
		return errors.New(errMsg)
	}
	return nil
}

func RegexpRemoveTask(r *regexp.Regexp) error {
	for key, _ := range TaskMap {
		result := r.MatchString(key)
		if result {
			if err := RemoveTask(key); err != nil {
				return err
			}
		}
	}
	return nil
}
