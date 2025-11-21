package v1

import (
	"friberg/internal/model/iostruct"

	"github.com/gogf/gf/v2/frame/g"
)

type FribergHelloReq struct {
	g.Meta `path:"/friberg/hello" tags:"Hello" method:"get" summary:"You first hello api"`
}
type FribergHelloRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
type FuzzReq struct {
	g.Meta `path:"/fuzz" tags:"select" method:"get" summary:"search fuzz on name"`
	Name   string `json:"name" dc:"name to search"`
}

type FuzzRes struct {
	g.Meta   `mime:"text/html" example:"string"`
	GameInfo []iostruct.GameInfo `json:"game_info" dc:"names found"`
}

type GuessReq struct {
	g.Meta `path:"/guess" tags:"guess" method:"get" summary:"guess the game"`
	Id     string `json:"id" dc:"game id"`
}

type GuessRes struct {
	g.Meta  `mime:"application/json"` // 返回 JSON
	Success bool                      `json:"success" dc:"guess success"`
	Result  iostruct.Game             `json:"result"` // 只有一个 game 对象
}
