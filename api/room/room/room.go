package room

import "github.com/gogf/gf/v2/frame/g"

// 房间创建
type RoomCreateReq struct {
	g.Meta `path:"/friberg/room/create" tags:"room" method:"post" summary:"创建房间"`
	Uuid   string `json:"uuid" dc:"用户唯一标识"`
}
type RoomCreateRes struct {
	g.Meta `mime:"application/json" example:"string"`
	RoomId string `json:"roomId" dc:"房间唯一标识"`
}

// 房间退出
type RoomExitReq struct {
	g.Meta `path:"/friberg/room/exit" tags:"room" method:"post" summary:"退出房间"`
	RoomId string `json:"roomId" dc:"房间唯一标识"`
	Uuid   string `json:"uuid" dc:"用户唯一标识"`
}
type RoomExitRes struct{}

// 房间开始游戏
type RoomGameBeginReq struct {
	g.Meta `path:"/friberg/room/begin" tags:"room" method:"post" summary:"开始游戏"`
	RoomId string `json:"roomId" dc:"房间唯一标识"`
	Uuid   string `json:"uuid" dc:"用户唯一标识"`
}
type RoomGameBeginRes struct{}

// 加入房间
type RoomJoinReq struct {
	g.Meta `path:"/friberg/room/join" tags:"room" method:"post" summary:"加入房间"`
	RoomId string `json:"roomId" dc:"房间唯一标识"`
	Uuid   string `json:"uuid" dc:"用户唯一标识"`
}

type RoomJoinRes struct{}
