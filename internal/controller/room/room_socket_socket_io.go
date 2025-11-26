package room

import (
	"context"
	"fmt"

	"friberg/api/room/socket"
	"friberg/internal/consts"
	"friberg/internal/game/manager"
	"friberg/internal/game/player"
	isocket "friberg/internal/library/websocket"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

func (c *ControllerSocket) SocketIo(ctx context.Context, req *socket.SocketIoReq) (res *socket.SocketIoRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	ch := make(chan *websocket.Conn, 1)
	var (
		address = r.RemoteAddr
		header  = fmt.Sprintf("%v", r.Header)
	)
	client := player.NewWsClient(nil, ctx, req.UserName, address, header)
	go isocket.HandleWsRequest(r, ch, func() {
		manager.RemovePlayerBusiness(client.Uuid)
	})
	ws := <-ch // 等待连接初始化完成
	if ws == nil {
		return nil, gerror.New("websocket upgrade failed")
	}
	client.Ws = ws
	manager.PM.Add(client)
	msg := isocket.Message{
		MsgType: consts.MsgType_Init,
		Data:    client.Uuid,
	}
	ws.WriteJSON(msg)
	return &socket.SocketIoRes{
		Uuid: client.Uuid,
	}, nil
}
