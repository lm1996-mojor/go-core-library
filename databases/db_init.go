package databases

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/config"
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/global"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/redis"
	localCipher "github.com/lm1996-mojor/go-core-library/utils/cipher"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var mutex sync.Mutex                  // 锁对象
var dbMap = make(map[string]*gorm.DB) // 租户数据库map

// ClientDb  租户数据库信息
type clientDb struct {
	ClientId int64  `json:"-,omitempty"`       // 租户id
	DbHost   string `json:"dbHost,omitempty"`  // 租户专属数据库连接地址
	DbPort   string `json:"dbPort,omitempty"`  // 租户专属数据库连接端口
	DbName   string `json:"dbName,omitempty"`  // 租户专属数据库名称
	DbUser   string `json:"dbUser,omitempty"`  // 租户专属数据库账户
	DbPass   string `json:"dbPass,omitempty"`  // 租户专属数据库密码
	DbType   string `json:"dbType,omitempty"`  // 租户专属数据库类型（mysql/Oracle/PostgreSQL/DB2/SQL Server、MariaDB）
	EnvType  int8   `json:"envType,omitempty"` // 数据库环境类型（1 线上 2 开发  3 测试 4 体验）
}

// GormLogger 自定义Gorm日志结构体
type GormLogger struct{}

// NewGormLogger 获取自定义Gorm日志对象
func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

// Printf 自定义格式打印Gorm日志
func (*GormLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

var newLogger = logger.New(
	NewGormLogger(), // io writer
	logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Info, // Log level
		Colorful:      false,       // Disable color
	},
)

const runLevel = -3

// 初始化数据库操作
func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

// Init 初始化数据库信息实现方法
func Init(app *iris.Application) {
	if config.Sysconfig.DataBases.EnableDbDynamicAddition {
		go GetSubscriptionMessagesFromCache()
	}
	if config.Sysconfig.DataBases.ClientEnable {
		initClientDB()
	}
	initCustomizedDB() //初始化自定义数据库信息
}

// InitCustomizedDB 初始化自定义数据库信息
func initCustomizedDB() {
	//获取自定义的数据库信息：从配置文件中获取，即config/Sysconfig结构体中获取
	dbInfoList := config.Sysconfig.DataBases.DbInfoList
	for _, database := range dbInfoList {
		//创建临时数据库连接变量
		dsn := database.DbUser + ":" + database.DbPass + "@tcp(" + database.Host + ":" +
			database.Port + ")/" + database.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		//打开连接
		clog.Info("自定义数据库连接：" + dsn)
		db, err := connectDB(dsn)
		if err != nil {
			panic("自定义数据库连接错误: " + err.Error())
		}
		//定制key,将打开的连接存入到map中
		databaseName := database.DbName
		mutex.Lock()
		dbMap[databaseName] = db
		mutex.Unlock()
	}
}

// InitClientDB 初始化租户数据库信息
func initClientDB() {
	platformDbConnectAddress := "root:123.com@tcp(192.168.0.62:62232)/platform_management?charset=utf8&parseTime=True&loc=Local"
	if config.Sysconfig.SystemEnv.Env == "prod" {
		if config.Sysconfig.DataBases.PDns != "" {
			// 链接初次解密
			decodeString, err := base64.StdEncoding.DecodeString(config.Sysconfig.DataBases.PDns)
			if err != nil {
				panic("初次解密错误")
			}
			// 201dd1f39f184638 = MD5(link_cipher)加密后的16位
			pDns := localCipher.AesDecryptECB(decodeString, []byte(_const.DbLinkEncryptKey))
			if len(pDns) <= 0 {
				panic("数据库链接解密错误")
			}
			platformDbConnectAddress = string(pDns)
		} else {
			panic("检测到目前配置为生产线环境，请配置平台数据库连接地址")
		}
	}
	platformDb, err1 := gorm.Open(mysql.Open(platformDbConnectAddress), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // love表将是love，不再是loves，即可成功取消表明被加s
	})
	if err1 != nil {
		clog.Error("初始化平台数据库连接错误")
		panic(err1)
	}
	dbInfoList := make([]clientDb, 0)
	platformDb.Table("client_db").Where("env_type = ?", 1).Scan(&dbInfoList)
	if len(dbInfoList) <= 0 {
		clog.Info("没有租户，暂不进行租户初始化")
		return
	}
	//获取自定义的数据库信息：从配置文件中获取，即config/Sysconfig结构体中获取
	for _, database := range dbInfoList {
		//创建临时数据库连接变量
		dsn := database.DbUser + ":" + database.DbPass + "@tcp(" + database.DbHost + ":" +
			database.DbPort + ")/" + database.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		//打开连接
		clog.Info("租户数据库连接：" + dsn)
		db, err := connectDB(dsn)
		if err != nil {
			panic("租户数据库连接错误: " + err.Error())
		}
		//定制key,将打开的连接存入到map中
		mutex.Lock()
		dbMap[fmt.Sprintf("%d", database.ClientId)] = db
		mutex.Unlock()
	}
}

// 打开数据库连接
func connectDB(dsn string) (db *gorm.DB, err error) {
	//通过传输进来的dsn信息，使用mysql.open方法打开数据的连接，并配置gorm.config结构体相关的信息
	// NamingStrategy ：取消默认表名

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // love表将是love，不再是loves，即可成功取消表明被加s
		Logger:         newLogger,                                  //指定自定义的gorm日志结构体
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err1 := db.DB() // 通过连接池创建出单例的数据库连接对象,并存入池中
	if err1 != nil {
		panic(err)
	}

	//配置连接对象的连接信息
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	return
}

func GetSubscriptionMessagesFromCache() {
	for {
		ctx := context.Background()
		//【Subscribe】订阅频道
		sub := redis.RedisPSubscribe(ctx, "client_db_add@*")
		if sub != nil {
			dbDnsMap := make(map[string]string)
			// 订阅者实时接收频道中的消息
			for msg := range sub.Channel() {
				// 打印频道号和消息内容
				//fmt.Printf("接收到来自频道%s的消息: %s\n",
				//	 msg.Channel, msg.Payload)
				split := strings.Split(msg.Channel, "@")
				dbDnsMap[split[1]] = msg.Payload
			}
			// 遍历数据
			for dbKey, dns := range dbDnsMap {
				// 判断新连接是否已经在缓存中
				if _, ok := dbMap[dbKey]; ok {
					continue
				}
				// 连接数据库
				db, err := connectDB(dns)
				if err != nil {
					clog.Error("连接数据库失败:" + dbKey + "，连接为【" + dns + "】")
				}
				mutex.Lock()
				dbMap[dbKey] = db
				mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Minute)
	}

}
