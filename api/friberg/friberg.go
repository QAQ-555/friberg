// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package friberg

import (
	"context"

	"friberg/api/friberg/game"
	"friberg/api/friberg/v1"
)

type IFribergGame interface {
	StartBase(ctx context.Context, req *game.StartBaseReq) (res *game.StartBaseRes, err error)
}

type IFribergV1 interface {
	FribergHello(ctx context.Context, req *v1.FribergHelloReq) (res *v1.FribergHelloRes, err error)
	Fuzz(ctx context.Context, req *v1.FuzzReq) (res *v1.FuzzRes, err error)
	Guess(ctx context.Context, req *v1.GuessReq) (res *v1.GuessRes, err error)
}
