package redis

import (
	"context"
	"strings"
	"sync"

	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
)

var mutex sync.Mutex

func renewDb() {
	for {
		//【Subscribe】订阅频道
		sub := RedisPSubscribe(context.Background(), "client_db_?*")
		if sub != nil {
			//dbDnsMap := make(map[string]string)
			// 订阅者实时接收频道中的消息
			select {
			case msg := <-sub.Channel():
				log.Info("发现新的数据源订阅，处理订阅信息：" + msg.Channel)

				clientOperate := strings.ReplaceAll(msg.Channel, _const.ClientDbRedisPSubscribe, "")
				split := strings.Split(clientOperate, "_")
				switch split[0] {
				case _const.DbOperateAdd:
					// 遍历数据
					log.Info("装载数据源")
					// 判断新连接是否已经在缓存中
					if _, ok := databases.GetDbMap()[split[1]]; ok {
						break
					}
					// 连接数据库
					db, err := databases.ConnectDB(msg.Payload)
					if err != nil {
						log.Error("连接数据库失败:" + split[1] + "，连接为【" + msg.Payload + "】")
					}
					mutex.Lock()
					databases.SetDbMap(split[1], db)
					mutex.Unlock()
				case _const.DbOperateUpdate:
					// 连接数据库
					db, err := databases.ConnectDB(msg.Payload)
					if err != nil {
						log.Error("连接数据库失败:" + split[1] + "，连接为【" + msg.Payload + "】")
					}
					mutex.Lock()
					databases.SetDbMap(split[1], db)
					mutex.Unlock()
				case _const.DbOperateDel:
					mutex.Lock()
					databases.DelDb(split[1])
					mutex.Unlock()
				}
			}
		}
	}
}
