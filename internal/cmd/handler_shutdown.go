// Package cmd
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cmd

import (
	"context"
	"os"
	"sync"

	"friberg/internal/consts"
	"friberg/utility/event.go"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
)

var (
	serverCloseSignal = make(chan struct{}, 1)
	serverWg          = sync.WaitGroup{}
	once              sync.Once
)

// signalHandlerForOverall 关闭信号处理
func signalHandlerForOverall(sig os.Signal) {
	serverCloseEvent(gctx.GetInitCtx())
	serverCloseSignal <- struct{}{}
}

// signalListen 信号监听
func signalListen(ctx context.Context, handler ...gproc.SigHandler) {
	go func(ctx context.Context) {
		gproc.AddSigHandlerShutdown(handler...)
		gproc.Listen()
	}(ctx)
}

// serverCloseEvent 关闭事件
// 区别于服务收到退出信号后的处理，只会执行一次
func serverCloseEvent(ctx context.Context) {
	once.Do(func() {
		event.Event().Call(consts.EventServerClose, ctx)
	})
}
