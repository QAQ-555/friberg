package v1

import "github.com/gogf/gf/v2/frame/g"

type FribergHelloReq struct {
	g.Meta `path:"/friberg/hello" tags:"Hello" method:"get" summary:"You first hello api"`
}
type FribergHelloRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
