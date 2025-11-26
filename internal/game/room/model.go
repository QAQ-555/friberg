package room

import (
	"friberg/internal/game/player"

	"github.com/gogf/gf/v2/util/guid"
)

type Room struct {
	Roomid     string
	Owner      string
	PlayerList []*player.WsClient
}

type OutRoomlist struct {
	RoomList []OutRoom `json:"room_list"`
}

type OutRoom struct {
	Roomid    string `json:"roomid"`
	Owner     string `json:"owner"`
	NumPlayer int    `json:"num_player"`
}

func NewRoom(uuid string, args ...string) *Room {
	str := ""
	for _, arg := range args {
		str += arg
	}
	return &Room{
		Roomid:     guid.S([]byte(str)),
		Owner:      uuid,
		PlayerList: make([]*player.WsClient, 0),
	}
}
