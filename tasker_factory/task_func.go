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

// AddTask 往库中新增任务
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

// InitTask 初始化任务对象
func InitTask(opts ...cron.Option) Task {
	var task Task
	task.TaskBody = cron.New(opts...)
	return task
}

// BatchedRunTasker 批量运行库中的所有任务
func BatchedRunTasker() {
	for _, key := range taskKeys {
		value, _ := TaskMap.Load(key)
		task := value.(Task)
		if task.TaskStatus {
			continue
		}
		task.TaskStatus = true
		go task.TaskBody.Run()
	}
}

// BatchedRunSpecifiedTask 批量运行指定任务
func BatchedRunSpecifiedTask(keys []string) {
	for _, key := range keys {
		value, _ := TaskMap.Load(key)
		task := value.(Task)
		if task.TaskStatus {
			continue
		}
		task.TaskStatus = true
		go task.TaskBody.Run()
	}
}

// RunTask 启动单个指定任务
func RunTask(key string) {
	value, _ := TaskMap.Load(key)
	task := value.(Task)
	if task.TaskStatus {
		return
	}
	task.TaskStatus = true
	go task.TaskBody.Run()
}

// StopTask 停止单个指定任务
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

// BatchedStopTask 批量停止指定任务
func BatchedStopTask(keys []string) {
	for _, key := range keys {
		StopTask(key)
	}
}

// RegexpStopTask 正则匹配停止任务
func RegexpStopTask(r *regexp.Regexp) {
	for _, key := range taskKeys {
		result := r.MatchString(key)
		if result {
			StopTask(key)
		}
	}
}

// StopAndRemoveTask 停止并删除单个任务
func StopAndRemoveTask(key string) error {
	StopTask(key)
	return RemoveTask(key)
}

// BatchedStopAndRemoveTask 批量停止并删除指定任务
func BatchedStopAndRemoveTask(keys []string) error {
	BatchedStopTask(keys)
	return BatchRemoveTask(keys)
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

// BatchRemoveTask 批量删除指定任务
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

// RegexpRemoveTask 正泽匹配删除任务
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

// RenewTaskNextTime 更新指定任务的下次执行时间
func RenewTaskNextTime(key string, nextSpec string) error {
	value, _ := TaskMap.Load(key)
	task := value.(Task)
	job := task.TaskJob
	opts := task.TaskOpt
	desc := task.TaskDesc
	err := StopAndRemoveTask(key)
	if err != nil {
		return err
	}
	err = AddTask(key, desc, nextSpec, job, opts...)
	if err != nil {
		return err
	}
	return nil
}
