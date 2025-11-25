package room

import (
	"context"

	"friberg/api/room/room"
	"friberg/internal/game/manager"
)

func (c *ControllerRoom) RoomJoin(ctx context.Context, req *room.RoomJoinReq) (res *room.RoomJoinRes, err error) {
	err = manager.RM.AddPlayer(req.RoomId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return
}
