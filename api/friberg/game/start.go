package game

import "github.com/gogf/gf/v2/frame/g"

type StartBaseReq struct {
	g.Meta `path:"/friberg/begin/base" tags:"game" method:"get" summary:"start a base single game"`
}

type StartBaseRes struct {
	g.Meta `mime:"application/cookie"`
}

type StartMultiReq struct {
	g.Meta `path:"/friberg/begin/multi" tags:"game" method:"get" summary:"start a multi game"`
}

type StartMultiRes struct {
	g.Meta `mime:"application/cookie"`
}
