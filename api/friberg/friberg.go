// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package friberg

import (
	"context"

	"friberg/api/friberg/v1"
)

type IFribergV1 interface {
	FribergHello(ctx context.Context, req *v1.FribergHelloReq) (res *v1.FribergHelloRes, err error)
}
