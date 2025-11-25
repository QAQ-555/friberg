package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// LoggerMiddleware 使用 g.Log 打印请求信息和 TraceID
func LoggerMiddleware(r *ghttp.Request) {
	ctx := r.Context()
	start := time.Now()

	// 打印请求方法和路径
	g.Log().Infof(ctx, "%s %s", r.Method, r.RequestURI)

	// 继续执行请求
	r.Middleware.Next()

	// 打印耗时
	duration := time.Since(start)
	g.Log().Infof(ctx, "%s %s completed in %v", r.Method, r.RequestURI, duration)
}
