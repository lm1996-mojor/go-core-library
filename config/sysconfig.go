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
		Name                string // 当前服务名称
		Host                string // 当前服务的本地地址
		Port                string // 当前服务的访问端口
		TimeZone            string // 当前服务的时区
		Language            string // 当前服务的语言
		GlobalReqPathPrefix string // 当前服务的全局请求的地址前缀
	}
	Consul struct {
		Addr                string   // 服务治理中心的地址
		Port                int      // 服务治理中心的端口
		EnableObtainService bool     // 开启获取服务工具
		Check               struct { // 服务检查对象
			CheckTimeout             string // 服务检查请求超时时间
			CheckInterval            string // 服务检查间隔时间
			InvalidServiceLogoutTime string // 无效的服务注销时间
		}
		Service struct {
			Spec               string     // 定时获取服务时间命令
			DesignatedServices []struct { // 当前服务需要获取的服务列表
				ServiceName string // 需要获取的服务的名称
				/** 获取的服务负载均衡模式
				 * rr：顺序轮询，轮流分配到后端服务器；
				 * wrr：权重轮询，根据后端服务器负载情况来分配；
				 * lc：最小连接，分配已建立连接最少的服务器上；
				 * wlc：权重最小连接，根据后端服务器处理能力来分配。
				 * ssi: 请求会话id，根据前端发送的会话中的id来绑定某个固定服务进行指定服务请求（一般用于文件传输等操作）
				 */
				LoadBalanceMode string
			}
		}
	}
	//数据库内部结构体
	DataBases struct {
		ClientEnable          bool       // 租户数据源配置 默认为关闭状态
		MasterDbName          string     // 主数据库名称
		PDns                  string     // 平台数据库连接地址
		EnableDbDynamicManage bool       // 开始数据源动态管理(默认为关闭状态)
		DbInfoList            []struct { // 多个自定义数据库源
			Host   string // 数据库访问地址
			Port   string // 数据库访问端口
			DbName string // 项目连接的数据库
			DbUser string // 数据库用户
			DbPass string // 数据密码
		}
	}

	//Redis配置
	Redis struct {
		ConnInfo string // 示例：指定ip:指定端口号 / :6379(默认为本地ip)
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
		Token                   bool   // 是否开启token检测
		TokenService            string // token检查使用的服务
		TokenCheckServiceApiUrl string // token检查服务接口地址
		Authentication          bool   // 是否开启鉴权检测
		AuthService             string // 权限检查使用的服务
		AuthCheckServiceApiUrl  string // 权限检查使用的服务地址
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
	//viper.SetConfigType("yaml")
	////从固定位置读取配置文件
	//viper.SetConfigFile("./" + _const.CONFIG)
	////读取配置文件
	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic("read configuration file failed: " + err.Error())
	//}
	//err = viper.Unmarshal(&Sysconfig)
	//if err != nil {
	//	panic("load configuration failed: " + err.Error())
	//}
	ReadConfigFile("./"+_const.CONFIG, "yaml", &Sysconfig)
	// 配置有效性检验
	validation()
	log.Infof("初始化配置: \n %v", Sysconfig)
}

func ReadConfigFile(filePath string, fileType string, rawObj interface{}) interface{} {
	// viper 参考链接：https://www.cnblogs.com/randysun/p/15889494.html
	//导入配置文件
	//确定文件类型 支持从JSON、TOML、YAML、HCL、INI和Java properties文件中读取配置数据。
	viper.SetConfigType(fileType)
	//从固定位置读取配置文件
	viper.SetConfigFile(filePath)
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic("read configuration file failed: " + err.Error())
	}
	err = viper.Unmarshal(&rawObj)
	if err != nil {
		panic("load configuration failed: " + err.Error())
	}
	return rawObj
}

// 初始化方法
func init() {
	// 全局初始化注册，根据runLevel变量的值来进行选择谁先初始化
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}
