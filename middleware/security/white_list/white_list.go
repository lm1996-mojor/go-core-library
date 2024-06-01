package white_list

import (
	"strings"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/databases"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
)

type Url struct {
	ReqUrl    string // 接口路径
	Method    string // 请求方式
	CheckType int    //接口检查类型(1:免token 2：免鉴权)
}

var tokenWhiteListMap = make(map[string]string)

var authWhiteListMap = make(map[string]string)

func Init() {
	clog.Info("初始化路由白名单")
	defaultWhiteList := make([]Url, 0)
	defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: "/consul/ser/health", Method: "get", CheckType: 1})
	list := InitSystemList()
	if len(list) > 0 {
		defaultWhiteList = append(defaultWhiteList, list...)
	}
	if len(defaultWhiteList) > 0 {
		AppendList(defaultWhiteList)
	} else {
		clog.Info("没有检测要求,无需初始化")
	}
	timedExecution()
	clog.Info("初始化完成")
}

func InitSystemList() []Url {
	defaultWhiteList := make([]Url, 0)
	list := tokenWhiteListInit()
	if len(list) > 0 {
		defaultWhiteList = append(defaultWhiteList, list...)
	}
	list = authWhiteListInit()
	if len(list) > 0 {
		defaultWhiteList = append(defaultWhiteList, list...)
	}
	return defaultWhiteList
}

func tokenWhiteListInit() (tokenWhiteList []Url) {
	clog.Info("获取token白名单....")
	db := databases.GetDbByName("platform_management", "mysql").Table("permissions_menu").
		Where("is_white_list = ?", 1).Where("req_url != '' or req_url is not null").Where("status = ?", 1)
	if config.Sysconfig.App.GlobalReqPathPrefix != "" && len(config.Sysconfig.App.GlobalReqPathPrefix) > 0 && config.Sysconfig.App.GlobalReqPathPrefix != "null" {
		db = db.Where("req_url like ?", config.Sysconfig.App.GlobalReqPathPrefix+"%")
	}
	db.Select("req_url,method").Find(&tokenWhiteList)
	for i := 0; i < len(tokenWhiteList); i++ {
		tokenWhiteList[i].CheckType = 1
	}
	return tokenWhiteList
}

func authWhiteListInit() (authWhiteList []Url) {
	clog.Info("获取权限白名单....")
	db := databases.GetDbByName("platform_management", "mysql").Table("permissions_menu").
		Where("is_auth_white_list = ?", 1).Where("req_url != '' or req_url is not null").Where("status = ?", 1).Where("menu_type = ? or menu_type = ?", 3, 4)
	if config.Sysconfig.App.GlobalReqPathPrefix != "" && len(config.Sysconfig.App.GlobalReqPathPrefix) > 0 && config.Sysconfig.App.GlobalReqPathPrefix != "null" {
		db = db.Where("req_url like ?", config.Sysconfig.App.GlobalReqPathPrefix+"%")
	}
	db.Select("req_url,method").Find(&authWhiteList)
	for i := 0; i < len(authWhiteList); i++ {
		authWhiteList[i].CheckType = 2
	}
	return authWhiteList
}

func timedExecution() {
	spec := "@every 11s"
	err := tasker_factory.AddTask("DelayRefreshList", "白名单刷新定时任务", spec, DelayRefreshList)
	if err != nil {
		panic("添加延迟刷新白名单列表定时任务添加失败" + err.Error())
	}
}

func DelayRefreshList() {
	clog.Info("白名单刷新中...")
	items := InitSystemList()
	for _, item := range items {
		if item.CheckType == 1 {
			tokenWhiteListMap[item.ReqUrl] = item.Method
		} else {
			authWhiteListMap[item.ReqUrl] = item.Method
		}
	}
	clog.Info("白名单刷新完成")
}

// AppendList append to URL white list
func AppendList(items []Url) {
	if len(items) <= 0 {
		clog.WarnF("新增白名单，传入参数为0,items：", len(items))
		return
	}
	for _, item := range items {
		if item.Method == "" {
			panic("白名单[" + item.ReqUrl + "]请求方式不能为空")
		}
		if item.CheckType < 1 || item.CheckType > 2 {
			panic("接口检查类型不符合规范，仅支持T（免token）/A（免鉴权）")
		}
		if item.CheckType == 1 {
			if _, ok := tokenWhiteListMap[item.ReqUrl]; !ok {
				tokenWhiteListMap[item.ReqUrl] = item.Method
			}
			continue
		} else {
			if _, ok := authWhiteListMap[item.ReqUrl]; !ok {
				authWhiteListMap[item.ReqUrl] = item.Method
			}
			continue
		}
	}
}

func InList(path string, method string, checkType int) bool {
	msgStr := ""
	if checkType == 1 {
		msgStr = "token"
		return match(path, tokenWhiteListMap, method, msgStr)
	} else {
		msgStr = "权限"
		return match(path, authWhiteListMap, method, msgStr)
	}
}

func match(reqPath string, srcReqPathSlice map[string]string, method string, msgStr string) bool {
	// user/{id}
	for srcReqPath, srcMethod := range srcReqPathSlice {
		if strings.ToUpper(srcMethod) == strings.ToUpper(method) {
			if strings.Contains(srcReqPath, "{") {
				index := strings.IndexByte(srcReqPath, '{')
				if len(reqPath) > index {
					if srcReqPath[0:index] == reqPath[0:index] {
						if strings.Count(srcReqPath[strings.Index(srcReqPath, "{")-1:], "/") == strings.Count(reqPath[strings.Index(srcReqPath, "{")-1:], "/") {
							if srcReqPath[0:strings.Index(srcReqPath, "{")] == reqPath[0:len(srcReqPath[0:strings.Index(srcReqPath, "{")])] {
								clog.Info(reqPath + "：" + msgStr + "路径参数匹配成功")
								clog.Info(reqPath + "：" + msgStr + "白名单匹配结果：在名单中")
								return true
							}
						}
					}
				}
			}
			if srcReqPath == reqPath {
				clog.Info(reqPath + "：" + msgStr + "普通路径匹配成功")
				clog.Info(reqPath + "：" + msgStr + "白名单匹配结果：在名单中")
				return true
			}
		}
	}
	clog.Info(reqPath + "：" + msgStr + "白名单匹配结果：不在名单中")
	return false
}
