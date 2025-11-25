package room

import (
	"context"

	"friberg/api/room/room"
	"friberg/internal/game/manager"
)

func (c *ControllerRoom) RoomExit(ctx context.Context, req *room.RoomExitReq) (res *room.RoomExitRes, err error) {
	err = manager.RM.RemovePlayer(req.RoomId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return
}
