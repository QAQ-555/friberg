package socket

import "github.com/gogf/gf/v2/frame/g"

type SocketIoReq struct {
	g.Meta   `path:"/friberg/socket" tags:"socket.io" method:"ALL" summary:"socket.io"`
	UserName string `json:"userName" dc:"userName"  v:"required#userName is required"`
}

type SocketIoRes struct {
	g.Meta `mime:"text/html" example:"string"`
	Uuid   string `json:"uuid" dc:"uuid"`
}
