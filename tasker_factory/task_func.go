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

func AddTask(taskStoreKey, taskDesc, spec string, cmd func(), opts ...cron.Option) error {
	task := InitTask(opts...)
	taskId, err := task.TaskBody.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	task.TaskId = taskId
	task.TaskDesc = taskDesc
	TaskMap.Store(taskStoreKey, task)
	taskKeys = append(taskKeys, taskStoreKey)
	return nil
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
	for _, key := range taskKeys {
		value, _ := TaskMap.Load(key)
		task := value.(Task)
		task.TaskStatus = true
		task.TaskBody.Run()
	}
}

func StopTask(key string) {
	v, ok := TaskMap.Load(key)
	if ok {
		t := v.(Task)
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
	for _, key := range taskKeys {
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
	v, ok := TaskMap.Load(key)
	if !ok {
		return errors.New("没有找到指定的任务: " + key)
	}
	t := v.(Task)
	taskId := t.TaskId
	if t.TaskStatus {
		return errors.New("当前任务(" + key + ")id为：(" + cast.ToString(taskId) + ")正在执行中，请先停止任务。再进行移除")
	}
	t.TaskBody.Remove(taskId)
	TaskMap.Delete(key)
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
			v, ok := TaskMap.Load(key)
			if !ok {
				return errors.New("没有找到指定的任务: " + key)
			}
			t := v.(Task)
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
	for _, key := range taskKeys {
		result := r.MatchString(key)
		if result {
			if err := RemoveTask(key); err != nil {
				return err
			}
		}
	}
	return nil
}
