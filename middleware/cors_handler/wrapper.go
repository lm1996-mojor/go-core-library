package cors_handler

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

// InitCors init wrapper for cors_handler
func InitCors(app *iris.Application) {
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		ctx := app.ContextPool.Acquire(w, r)
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Overwrite, Destination, Content-Type, Depth, User-Agent, Translate, Range, Content-Range, Timeout, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control, Location, Lock-Token, If, Authorization, watermark, angle")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Authorization, Content-Disposition")
		ctx.Header("Access-Control-Max-Age", "3600")
		if r.Method == "OPTIONS" {
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
			ctx.StatusCode(204)
			return
		}
		router(w, r)
		return
	})
}
