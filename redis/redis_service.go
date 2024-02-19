package redis

import (
	"context"
	"strings"
	"sync"

	"github.com/lm1996-mojor/go-core-library/databases"
	"github.com/lm1996-mojor/go-core-library/log"
)

var mutex sync.Mutex

func GetSubscriptionMessagesFromCache() {
	for {
		//【Subscribe】订阅频道
		sub := RedisPSubscribe(context.Background(), "client_db_add_?*")
		if sub != nil {
			dbDnsMap := make(map[string]string)
			// 订阅者实时接收频道中的消息
			select {
			case msg := <-sub.Channel():
				log.Info("发现新的数据源订阅，处理订阅信息：" + msg.Channel)
				clientId := strings.ReplaceAll(msg.Channel, "client_db_add_", "")
				dbDnsMap[clientId] = msg.Payload
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
	}
}
