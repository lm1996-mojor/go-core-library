package white_list

import (
	"regexp"

	"github.com/lm1996-mojor/go-core-library/databases"
	clog "github.com/lm1996-mojor/go-core-library/log"
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
	clog.Info("获取token白名单....")
	defaultWhiteList := make([]Url, 0)
	var tokenWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_white_list = ?", 1).Where("status = ?", 1).Where("menu_type = ? or menu_type = ?", 3, 4).
		Select("req_url").Find(&tokenWhiteList)
	for _, url := range tokenWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: "T"})
	}
	clog.Info("获取权限白名单....")
	var authWhiteList []string
	databases.GetDbByName("platform_management").Table("permissions_menu").
		Where("is_enable_auth = ?", 1).Where("status = ?", 1).Where("menu_type = ? or menu_type = ?", 3, 4).
		Select("req_url").Find(&authWhiteList)
	for _, url := range authWhiteList {
		defaultWhiteList = append(defaultWhiteList, Url{ReqUrl: url, CheckType: "A"})
	}
	AppendList(defaultWhiteList)
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
