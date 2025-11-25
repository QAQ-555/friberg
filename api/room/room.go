// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package room

import (
	"context"

	"friberg/api/room/room"
	"friberg/api/room/socket"
)

type IRoomRoom interface {
	RoomCreate(ctx context.Context, req *room.RoomCreateReq) (res *room.RoomCreateRes, err error)
	RoomExit(ctx context.Context, req *room.RoomExitReq) (res *room.RoomExitRes, err error)
	RoomGameBegin(ctx context.Context, req *room.RoomGameBeginReq) (res *room.RoomGameBeginRes, err error)
	RoomJoin(ctx context.Context, req *room.RoomJoinReq) (res *room.RoomJoinRes, err error)
}

type IRoomSocket interface {
	SocketIo(ctx context.Context, req *socket.SocketIoReq) (res *socket.SocketIoRes, err error)
}
