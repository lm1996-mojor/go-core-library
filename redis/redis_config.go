package redis

import (
	"strings"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/utils"
	"golang.org/x/net/context"

	"github.com/kataras/iris/v12"
	"github.com/redis/go-redis/v9"
)

var redisDb redis.UniversalClient

// Init 初始化redis
func Init(app *iris.Application) {
	if config.Sysconfig.Redis.Host != "" || len(config.Sysconfig.Redis.Host) != 0 {
		if strings.Contains(config.Sysconfig.Redis.Ports, ",") {
			log.Info("redis集群初始化")
			split := utils.SplitRedisPort(config.Sysconfig.Redis.Ports)
			redisConnectionInfos := make([]string, 0, len(split))
			for _, port := range split {
				redisConnectionInfos = append(redisConnectionInfos, config.Sysconfig.Redis.Host+":"+port)
			}
			RedisInit(redisConnectionInfos, context.Background())
		} else {
			log.Info("自定义redis初始化")
			RedisInit([]string{config.Sysconfig.Redis.Host + ":" + config.Sysconfig.Redis.Ports}, context.Background())
		}
	} else {
		log.Info("没有redis相关配置，无需进行redis初始化")
	}

}

func RedisInit(addrs []string, ctx context.Context) {
	redisDb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addrs,
	})
	if err := redisDb.Ping(ctx).Err(); err != nil {
		log.Error("redis初始化错误")
		panic(err)
	}
}

const runLevel = -4

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})

}
