package web

import (
	"github.com/kataras/iris/v12"
)

// URLParamPage returns parameter page and size
func URLParamPage(ctx iris.Context) (page int, size int, err error) {
	page, err = ctx.URLParamInt("page")
	if err != nil {
		return
	}
	size, err = ctx.URLParamInt("size")
	if err != nil {
		return
	}
	return
}
