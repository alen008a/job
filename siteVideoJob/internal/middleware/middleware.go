package middleware

import (
	"siteVideoJob/internal/context"
	"siteVideoJob/mdata"
	"siteVideoJob/xxl"
	"strconv"
	"time"
)

// CustomMiddleware 自定义中间件
func CustomMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(ctx *context.Context, param *xxl.RunReq) string {
		startTime := time.Now()
		ctx.Infof("[middleware] Start at: %v", startTime)
		ctx.Console("<<<< 执行脚本: %s >>>> <br>", param.ExecutorHandler)
		ctx.Console(
			"<<<< 脚本参数 >>>> <br> %s <br>------------------------------------------------------------------------------------------------<br>",
			param.ExecutorParams)
		var executorParams xxl.ExecutorParams
		mdata.Cjson.UnmarshalFromString(param.ExecutorParams, &executorParams)
		ctx.SiteId = strconv.Itoa(executorParams.SiteId)
		res := tf(ctx, param)
		since := time.Since(startTime)
		ctx.Infof("[middleware] End... execution time: %v", since)
		return res
	}
}
