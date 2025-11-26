package hello

import (
	"context"

	v1 "friberg/api/hello/v1"
	"friberg/internal/game/manager"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	manager.RM.ShowList()
	return
}
