package cmd

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gsession"

	"friberg/internal/controller/friberg"
	"friberg/internal/controller/hello"
	"friberg/internal/controller/room"
	imongo "friberg/internal/library/mongo"
	"friberg/internal/logic/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			imongo.Register()
			s.SetSessionMaxAge(time.Minute)
			s.SetSessionStorage(gsession.NewStorageRedis(g.Redis()))
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(middleware.MiddlewareCORS)
				group.Middleware(middleware.LoggerMiddleware)
				group.Bind(
					hello.NewV1(),
					friberg.NewV1(),
					friberg.NewGame(),
					room.NewSocket(),
					room.NewRoom(),
				)
				group.ALL("/get", func(r *ghttp.Request) {
					sessionData, err := r.Session.Data()
					if err != nil {
						r.Response.Write(err.Error())
						return
					}
					r.Response.Write(sessionData)
				})
			})
			signalListen(ctx, signalHandlerForOverall)
			s.Run()
			return nil
		},
	}
)
