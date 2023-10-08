package white_list

import (
	"regexp"

	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/rs/zerolog/log"
)

type Url struct {
	Url       string //接口路径
	CheckType string //接口检查类型(T:免token A：免鉴权)
}

var defaultWhiteList = []Url{
	{"/add/db", "T"},
	{"/replace/db", "T"},
	{"/drop/db", "T"},
	{"/clear/db", "T"},
}
var authList = make([]*regexp.Regexp, 0)

var whiteListRegex = make([]*regexp.Regexp, 0)

func Init() {
	clog.Info("初始化路由白名单")
	AppendList(defaultWhiteList)
}

// AppendList append to URL white list
func AppendList(items []Url) {
	var reg *regexp.Regexp
	var err error
	for _, item := range items {
		clog.Infof("新增的白名单：", item)
		reg, err = regexp.Compile(item.Url)
		if err == nil {
			if item.CheckType == "T" {
				whiteListRegex = append(whiteListRegex, reg)
			} else if item.CheckType == "A" {
				authList = append(authList, reg)
			} else {
				panic("接口检查类型不符合规范，仅支持T（免token）/A（免鉴权）")
			}
		} else {
			log.Error().Msg("invalid regex in white list: " + item.Url)
		}
	}
}

func InList(path string, checkType int) bool {
	list := make([]*regexp.Regexp, 0)
	if checkType == 1 {
		list = whiteListRegex
	} else {
		list = authList
	}
	//判断请求路径是否在白名单中
	for _, g := range list {
		if g.Match([]byte(path)) {
			clog.Info(path + " 白名单路由匹配结果：成功")
			return true
		}
	}
	clog.Info(path + " 白名单路由匹配结果：失败")
	return false
}
