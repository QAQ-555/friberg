package friberg

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"friberg/api/friberg/v1"
)

func (c *ControllerV1) FribergHello(ctx context.Context, req *v1.FribergHelloReq) (res *v1.FribergHelloRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
