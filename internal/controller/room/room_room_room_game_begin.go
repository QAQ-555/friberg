package room

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"friberg/api/room/room"
)

func (c *ControllerRoom) RoomGameBegin(ctx context.Context, req *room.RoomGameBeginReq) (res *room.RoomGameBeginRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
