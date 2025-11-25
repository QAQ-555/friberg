package room

import (
	"context"
	"friberg/api/room/room"
	"friberg/internal/game/manager"
	iroom "friberg/internal/game/room"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerRoom) RoomCreate(ctx context.Context, req *room.RoomCreateReq) (res *room.RoomCreateRes, err error) {
	if !manager.PM.IsPlayerExist(req.Uuid) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "player not exist")
	}
	newRoom := iroom.NewRoom(req.Uuid)
	g.Log().Info(ctx, "create room", newRoom.Roomid)
	err = manager.RM.CreateRoom(newRoom)
	if err != nil {
		return nil, err
	}
	err = manager.RM.AddPlayer(newRoom.Roomid, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &room.RoomCreateRes{
		RoomId: newRoom.Roomid,
	}, nil
}
