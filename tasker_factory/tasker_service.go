package tasker_factory

import (
	"context"
	"strings"
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/redis"
)

var mutex sync.Mutex

//var TaskerMap = map[string]map[string]*Tasker{}
//
//// 任务主函数
//func taskerMain(tasker *Tasker, taskerFunc func()) {
//	time.Sleep(time.Second * 1)
//	defer tasker.stop()
//	for {
//		select {
//		case <-tasker.t.C: //实现每隔1s执行任务
//			taskerFunc()
//		case <-tasker.stopChan:
//			log.Info("任务已停止")
//			return
//		}
//	}
//}
//
//// CreateTasker 创建任务
////
//// Description: 创建一个定时任务，根据执行时间
//func CreateTasker(execTime time.Duration, clientCode string, taskerName string, ownerObjName string) (taskerId string) {
//	log.Info("任务信息初始化")
//	if taskerName == "" {
//		taskerName = "taskerCli_" + clientCode
//	}
//	if clientCode == "" {
//		log.Error("任务id不能为空")
//		panic("服务器错误")
//	}
//	mutex.Lock()
//	taskerId = utils.GenerateCodeByUUID(18)
//	mutex.Unlock()
//	tasker := NewTasker(execTime, taskerId, taskerName, ownerObjName)
//	TaskerMap[clientCode][taskerId] = tasker
//	return taskerId
//}
//
//// ExecTasker 执行任务
//func ExecTasker(clientCode string, taskerId string, taskerFunc func()) {
//	tasker := TaskerMap[clientCode][taskerId]
//	log.Info("任务执行")
//	mutex.Lock()
//	tasker.TaskerStatus = 2
//	go taskerMain(tasker, taskerFunc)
//	mutex.Unlock()
//}
//
//// StopTasker 停止指定任务
//func StopTasker(clientCode string, taskerId string) {
//	GetTasker(clientCode, taskerId).stop()
//}
//
//// GetTasker 获取指定任务
//func GetTasker(clientCode string, taskerId string) *Tasker {
//	tasker, ok := TaskerMap[clientCode][taskerId]
//	if !ok {
//		log.Error("任务不存在，请检查任务id")
//		panic("服务器错误")
//	}
//	return tasker
//}
//
//// GetTaskerList 获取指定租户任务列表
//func GetTaskerList(clientCode string) (taskerList []*Tasker) {
//	for _, tasker := range TaskerMap[clientCode] {
//		taskerList = append(taskerList, tasker)
//	}
//	return taskerList
//}
//
//// RestTaskerNextExecTime 重置指定任务下次执行时间
//func RestTaskerNextExecTime(clientCode string, taskerId string, execTime time.Duration) {
//	GetTasker(clientCode, taskerId).reset(execTime)
//}
//
//// DeleteTasker 删除指定的任务
////
//// 如果租户code不为空且所属业务也不为空，则代表要删除指定租户下的指定任务，反之则删除该租户下的所有任务
//func DeleteTasker(clientCode string, taskerId string) {
//	if clientCode != "" && taskerId != "" {
//		delete(TaskerMap[clientCode], taskerId)
//	} else {
//		delete(TaskerMap, clientCode)
//	}
//}
//
//// DeleteStoppedTask 删除指定租户全部已停止的任务
//func DeleteStoppedTask(clientCode string) {
//	for cCode, tasker := range TaskerMap[clientCode] {
//		if tasker.TaskerStatus == 4 {
//			DeleteTasker(cCode, "")
//		}
//	}
//}

func GetSubscriptionMessagesFromCache() {
	for {
		//【Subscribe】订阅频道
		sub := redis.RedisPSubscribe(context.Background(), "client_db_add_?*")
		if sub != nil {

			dbDnsMap := make(map[string]string)
			// 订阅者实时接收频道中的消息
			select {
			case msg := <-sub.Channel():
				log.Info("发现新的数据源订阅，处理订阅信息：" + msg.Channel)
				split := strings.Split(msg.Channel, "add")
				dbDnsMap[split[1]] = msg.Payload
			}
			// 遍历数据
			log.Info("装载数据源")
			for dbKey, dns := range dbDnsMap {
				// 判断新连接是否已经在缓存中
				if _, ok := databases.GetDbMap()[dbKey]; ok {
					continue
				}
				// 连接数据库
				db, err := databases.ConnectDB(dns)
				if err != nil {
					log.Error("连接数据库失败:" + dbKey + "，连接为【" + dns + "】")
				}
				mutex.Lock()
				databases.SetDbMap(dbKey, db)
				mutex.Unlock()
			}
		}
		//time.Sleep(30 * time.Second)
	}
}

const runLevel = -1

// 初始化数据库操作
func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

// Init 初始化数据库信息实现方法
func Init(app *iris.Application) {
	if config.Sysconfig.DataBases.EnableDbDynamicAddition {
		go GetSubscriptionMessagesFromCache()
	}
}
