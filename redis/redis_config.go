package redis

import (
	"strings"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/utils"

	"github.com/go-redis/redis"
	"github.com/kataras/iris/v12"
)

var redisClusterDb *redis.ClusterClient
var redisSingleDb *redis.Client

// Init 初始化redis
func Init(app *iris.Application) {
	if config.Sysconfig.Redis.Host != "" || len(config.Sysconfig.Redis.Host) != 0 {
		if strings.Contains(config.Sysconfig.Redis.Ports, ",") {
			log.Info("redis集群初始化")
			ClusterInit()
		} else {
			if config.Sysconfig.Redis.Host != "127.0.0.1" || config.Sysconfig.Redis.Host != "localhost" {
				log.Info("自定义redis初始化")
				CustomSingerInit()
			} else {
				log.Info("默认redis初始化(本地redis)")
				SingerInit()
			}
		}
	} else {
		log.Info("没有redis相关配置，无需进行redis初始化")
	}

}

func ClusterInit() {
	split := utils.SplitRedisPort(config.Sysconfig.Redis.Ports)
	redisConnectionInfos := make([]string, 0, len(split))
	for _, port := range split {
		redisConnectionInfos = append(redisConnectionInfos, config.Sysconfig.Redis.Host+":"+port)
	}
	redisClusterDb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: redisConnectionInfos,
	})
	_, err := redisClusterDb.Ping().Result()
	if err != nil {
		log.Error("redis集群初始化错误")
		panic(err)
	}
}

func SingerInit() {
	redisSingleDb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	_, err := redisSingleDb.Ping().Result()
	if err != nil {
		log.Error("redis初始化错误")
		panic(err)
	}
}

func CustomSingerInit() {
	redisSingleDb = redis.NewClient(&redis.Options{
		Addr: config.Sysconfig.Redis.Host + ":" + config.Sysconfig.Redis.Ports,
		DB:   0,
	})
	_, err := redisSingleDb.Ping().Result()
	if err != nil {
		log.Error("redis初始化错误")
		panic(err)
	}
}

const runLevel = -4

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})

}
