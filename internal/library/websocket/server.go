package isocket

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

var (
	wsUpGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// In production, you should implement proper origin checking
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// Error callback function.
		},
	}
)

func HandleWsRequest(r *ghttp.Request, ch chan *websocket.Conn) {
	ctx := r.Context()
	ws, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		ch <- nil
		r.Response.Write(err.Error())
		return
	}
	defer HandleWsDisconnect(ctx, ws)

	ch <- ws
	g.Log().Info(ctx, "HandleWsRequest success")
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		g.Log().Infof(ctx, "Received message[%d]: %s", messageType, p)
	}
}

func HandleWsDisconnect(ctx context.Context, ws *websocket.Conn) error {
	//TODO: 处理断开连接
	return ws.Close()
}

func SendRoomList(ctx context.Context, ws *websocket.Conn) error {
	//TODO: 发送房间列表
	return nil
}
