package redis

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"
	"golang.org/x/net/context"

	"github.com/kataras/iris/v12"
	"github.com/redis/go-redis/v9"
)

var redisDb redis.UniversalClient

// Init 初始化redis
func Init(app *iris.Application) {
	if config.Sysconfig.Redis.ConnInfo != "" || len(config.Sysconfig.Redis.ConnInfo) != 0 {
		if strings.Contains(config.Sysconfig.Redis.ConnInfo, ",") {
			log.Info("redis集群初始化")
			redisConnections := strings.Split(config.Sysconfig.Redis.ConnInfo, ",")
			for _, s := range redisConnections {
				redisConn := strings.Split(s, ":")
				if len(redisConn) > 2 {
					panic("redis连接有问题请检查")
				}
				if matched, err := regexp.MatchString("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}", redisConn[0]); !matched || err != nil {
					panic("redis连接有问题请检查:连接ip有问题" + redisConn[0])
				}
				if toInt, err := strconv.Atoi(redisConn[1]); !(toInt > 0) || err != nil {
					panic("redis连接有问题请检查:端口有问题" + redisConn[1])
				}
			}
			ConnectionRedis(redisConnections, context.Background())
		} else {
			log.Info("自定义redis初始化")
			ConnectionRedis([]string{config.Sysconfig.Redis.ConnInfo}, context.Background())
		}
		if config.Sysconfig.DataBases.EnableDbDynamicManage {
			go GetSubscriptionMessagesFromCache()
		}
	} else {
		if config.Sysconfig.DataBases.EnableDbDynamicManage {
			panic("检测到系统中需要同步数据源，但没有相应的redis配置，请在项目根目录的yml文件中，添加redis配置")
		}
		log.Info("没有redis相连接需求")
	}
}

func ConnectionRedis(addrs []string, ctx context.Context) {
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
