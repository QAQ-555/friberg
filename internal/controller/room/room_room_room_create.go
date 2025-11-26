package room

import (
	"context"
	"fmt"
	"friberg/api/room/room"
	"friberg/internal/game/manager"
	iroom "friberg/internal/game/room"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerRoom) RoomCreate(ctx context.Context, req *room.RoomCreateReq) (res *room.RoomCreateRes, err error) {
	request := g.RequestFromCtx(ctx)
	var (
		address = request.RemoteAddr
		header  = fmt.Sprintf("%v", request.Header)
	)
	newRoom := iroom.NewRoom(req.Uuid, address, header)
	g.Log().Info(ctx, "create room", newRoom.Roomid)
	err = manager.CreateRoomBusiness(newRoom)
	if err != nil {
		return nil, err
	}
	err = manager.AddPlayerToRoomBusiness(newRoom.Roomid, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &room.RoomCreateRes{
		RoomId: newRoom.Roomid,
	}, nil
}
