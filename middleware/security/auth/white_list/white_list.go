package white_list

import (
	"regexp"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/databases"
	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/tasker_factory"
	"github.com/rs/zerolog/log"
)

type Url struct {
	ReqUrl    string //接口路径
	CheckType string //接口检查类型(T:免token A：免鉴权)
}

var authWhiteListRegex = make([]*regexp.Regexp, 0)

var tokenWhiteListRegex = make([]*regexp.Regexp, 0)

func Init() {
	clog.Info("初始化路由白名单")
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
	if len(defaultWhiteList) > 0 {
		AppendList(defaultWhiteList)
		TimedExecution()
	} else {
		clog.Info("没有检测要求,无需初始化")
	}
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
	defaultWhiteList := make([]Url, 0)
	defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: "/consul/ser/health", CheckType: "T"})
	clog.Info("获取token白名单....")
	var tokenWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_white_list = ?", 1).Where("req_url != '' or req_url is not null").Where("status = ?", 1).Select("req_url").Find(&tokenWhiteList)
	for _, url := range tokenWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: "T"})
	}
	return defaultWhiteList
}

func authWhiteListInit() []Url {
	defaultWhiteList := make([]Url, 0)
	clog.Info("获取权限白名单....")
	var authWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_enable_auth = ?", 2).Where("req_url != '' or req_url is not null").Where("status = ?", 1).Where("menu_type = ? or menu_type = ?", 3, 4).
		Select("req_url").Find(&authWhiteList)
	for _, url := range authWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: "A"})
	}
	return defaultWhiteList
}

func TimedExecution() {
	spec := "@every 11s"
	err := tasker_factory.AddTask("DelayRefreshList", "发现服务定时任务", spec, DelayRefreshList)
	if err != nil {
		panic("添加延迟刷新白名单列表定时任务添加失败" + err.Error())
	}
}

func DelayRefreshList() {
	clog.Info("白名单刷新中...")
	authWhiteListRegex = make([]*regexp.Regexp, 0)
	tokenWhiteListRegex = make([]*regexp.Regexp, 0)
	items := InitSystemList()
	if len(items) > 0 {
		var reg *regexp.Regexp
		var err error
		for _, item := range items {
			clog.Infof("新增的白名单：", item)
			reg, err = regexp.Compile(item.ReqUrl)
			if err == nil {
				if item.CheckType == "T" {
					tokenWhiteListRegex = append(tokenWhiteListRegex, reg)
				} else if item.CheckType == "A" {
					authWhiteListRegex = append(authWhiteListRegex, reg)
				} else {
					log.Error().Msg("接口检查类型不符合规范，仅支持T（免token）/A（免鉴权）")
				}
			} else {
				log.Error().Msg("invalid regex in white list: " + item.ReqUrl)
			}
		}
	}
	clog.Info("白名单刷新完成")
}

// AppendList append to URL white list
func AppendList(items []Url) {
	var reg *regexp.Regexp
	var err error
	for _, item := range items {
		clog.Infof("新增的白名单：", item)
		reg, err = regexp.Compile(item.ReqUrl)
		if err == nil {
			if item.CheckType == "T" {
				tokenWhiteListRegex = append(tokenWhiteListRegex, reg)
			} else if item.CheckType == "A" {
				authWhiteListRegex = append(authWhiteListRegex, reg)
			} else {
				panic("接口检查类型不符合规范，仅支持T（免token）/A（免鉴权）")
			}
		} else {
			log.Error().Msg("invalid regex in white list: " + item.ReqUrl)
		}
	}
}

func InList(path string, checkType int) bool {
	list := make([]*regexp.Regexp, 0)
	msgStr := ""
	if checkType == 1 {
		msgStr = "token"
		list = tokenWhiteListRegex
	} else {
		msgStr = "权限"
		list = authWhiteListRegex
	}
	//判断请求路径是否在白名单中
	for _, g := range list {
		if g.Match([]byte(path)) {
			clog.Info(path + "：" + msgStr + "白名单匹配结果：成功")
			return true
		}
	}
	clog.Info(path + "：" + msgStr + "白名单匹配结果：失败")
	return false
}
