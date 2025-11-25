package player

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	Uuid     string             `json:"uuid"` // 客户端唯一标识
	UserName string             // 用户名
	Ws       *websocket.Conn    // ws 链接
	HttpCtx  context.Context    // http上下文
	Ctx      context.Context    // 客户端的Ctx，客户端退出时会通知到Ctx.Done()，可用于取消当前客户端的异步任务
	Cancel   context.CancelFunc // 取消Ctx方法
	RoomId   string             // 所在房间
	Once     sync.Once          // 安全锁
	Mutex    sync.Mutex         // 用户锁
}

type OutClient struct {
	UserName []string `json:"userNames"` // 用户名
}

func NewWsClient(ws *websocket.Conn, ctx context.Context, userName string) *WsClient {
	Newctx, cancel := context.WithCancel(gctx.New())
	return &WsClient{
		Uuid:     guid.S(),
		Ws:       ws,
		HttpCtx:  ctx,
		Ctx:      Newctx,
		Cancel:   cancel,
		UserName: userName,
	}
}
