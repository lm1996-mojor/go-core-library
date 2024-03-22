package white_list

import (
	"strings"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/databases"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
)

type Url struct {
	ReqUrl    string //接口路径
	CheckType int    //接口检查类型(1:免token 2：免鉴权)
}

var tokenWhiteListMap = make(map[string]string)

var authWhiteListMap = make(map[string]string)

func Init() {
	clog.Info("初始化路由白名单")
	defaultWhiteList := make([]Url, 0)
	defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: "/consul/ser/health", CheckType: 1})
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
	if config.Sysconfig.Detection.Token {
		list := tokenWhiteListInit()
		if len(list) > 0 {
			defaultWhiteList = append(defaultWhiteList, list...)
		}
	}
	if config.Sysconfig.Detection.Authentication {
		list := authWhiteListInit()
		if len(list) > 0 {
			defaultWhiteList = append(defaultWhiteList, list...)
		}
	}
	return defaultWhiteList
}

func tokenWhiteListInit() []Url {
	clog.Info("获取token白名单....")
	defaultWhiteList := make([]Url, 0)
	var tokenWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_white_list = ?", 1).Where("req_url != '' or req_url is not null").Where("status = ?", 1).
		Where("req_url like ?", config.Sysconfig.App.GlobalReqPathPrefix+"%").
		Select("req_url").Find(&tokenWhiteList)
	for _, url := range tokenWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: 1})
	}
	return defaultWhiteList
}

func authWhiteListInit() []Url {
	clog.Info("获取权限白名单....")
	defaultWhiteList := make([]Url, 0)
	var authWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_enable_auth = ?", 2).Where("req_url != '' or req_url is not null").Where("status = ?", 1).Where("menu_type = ? or menu_type = ?", 3, 4).
		Where("req_url like ?", config.Sysconfig.App.GlobalReqPathPrefix+"%").
		Select("req_url").Find(&authWhiteList)
	for _, url := range authWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: 2})
	}
	return defaultWhiteList
}

func timedExecution() {
	spec := "@every 11s"
	err := tasker_factory.AddTask("DelayRefreshList", "发现服务定时任务", spec, DelayRefreshList)
	if err != nil {
		panic("添加延迟刷新白名单列表定时任务添加失败" + err.Error())
	}
}

func DelayRefreshList() {
	clog.Info("白名单刷新中...")
	items := InitSystemList()
	for _, item := range items {
		if item.CheckType == 1 {
			tokenWhiteListMap[item.ReqUrl] = "TOKEN"
		} else {
			authWhiteListMap[item.ReqUrl] = "AUTHENTICATION"
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
		if item.CheckType < 1 || item.CheckType > 2 {
			panic("接口检查类型不符合规范，仅支持T（免token）/A（免鉴权）")
		}
		if item.CheckType == 1 {
			if _, ok := tokenWhiteListMap[item.ReqUrl]; !ok {
				tokenWhiteListMap[item.ReqUrl] = "TOKEN"
			}
			continue
		} else {
			if _, ok := authWhiteListMap[item.ReqUrl]; !ok {
				authWhiteListMap[item.ReqUrl] = "AUTHENTICATION"
			}
			continue
		}
	}
}

func InList(path string, checkType int) bool {
	msgStr := ""
	if checkType == 1 {
		msgStr = "token"
		return match(path, tokenWhiteListMap, msgStr)
	} else {
		msgStr = "权限"
		return match(path, authWhiteListMap, msgStr)
	}
}

func match(reqPath string, srcReqPathSlice map[string]string, msgStr string) bool {
	for key, _ := range srcReqPathSlice {
		if strings.Contains(key, reqPath) {
			clog.Info(reqPath + "：" + msgStr + "白名单匹配结果：成功")
			return true
		}
	}
	clog.Info(reqPath + "：" + msgStr + "白名单匹配结果：失败")
	return false
}
