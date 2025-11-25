package friberg

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"friberg/api/friberg/game"
)

func (c *ControllerGame) StartMulti(ctx context.Context, req *game.StartMultiReq) (res *game.StartMultiRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
