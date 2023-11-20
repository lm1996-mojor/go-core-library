package config

import (
	_const "github.com/lm1996-mojor/go-core-library/const"
	"github.com/lm1996-mojor/go-core-library/global"
	"github.com/lm1996-mojor/go-core-library/log"

	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

// 全局配置结构体-私有化
type sysconfig struct {
	//服务信息结构体
	App struct {
		Cluster  string
		Name     string
		Host     string
		Port     string
		TimeZone string
		Language string
	}

	//数据库内部结构体
	DataBases struct {
		ClientEnable bool       // 多数据源配置 默认为关闭状态
		MasterDbName string     // 主数据库名称
		PDns         string     // 平台数据库连接地址
		DbInfoList   []dataBase // 多个自定义数据库源
	}

	//Redis配置
	Redis struct {
		Host  string
		Ports string
	}
	//系统环境参数配置
	SystemEnv struct {
		Env            string //系统环境（dev/test/prod）
		Id             int64  //主键id
		Domain         string //主机域名
		IsUse          int8   //是否正在使用（1、正在使用 2、未使用）
		EnvType        int8   //环境类型（1、开发环境 2、测试环境 3、生产线环境 ）
		ConnectHostIp  string //连接主机ip
		BelongUserName string //主机所属人
	}
	// 检测配置
	Detection struct {
		Token          bool // 是否开启token检测
		Authentication bool // 是否开启鉴权检测
	}
}

// 数据库内部结构体
type dataBase struct {
	Host   string // 数据库访问地址
	Port   string // 数据库访问端口
	DbName string // 项目连接的数据库
	DbUser string // 数据库用户
	DbPass string // 数据密码
}

// 启动级别 （值越小，优先级越高）
const runLevel = -11

// Sysconfig 将私有的全局配置结构体公开化
var Sysconfig = sysconfig{}

func Init(app *iris.Application) {
	// viper 参考链接：https://www.cnblogs.com/randysun/p/15889494.html
	//导入配置文件
	//确定文件类型 支持从JSON、TOML、YAML、HCL、INI和Java properties文件中读取配置数据。
	viper.SetConfigType("yaml")
	//从固定位置读取配置文件
	viper.SetConfigFile("./" + _const.CONFIG)
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic("read configuration file failed: " + err.Error())
	}
	err = viper.Unmarshal(&Sysconfig)
	if err != nil {
		panic("load configuration failed: " + err.Error())
	}
	// 配置有效性检验
	validation()
	log.Infof("初始化配置: \n %v", Sysconfig)
}

// 初始化方法
func init() {
	// 全局初始化注册，根据runLevel变量的值来进行选择谁先初始化
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}
