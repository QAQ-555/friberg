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

func NewRoom(uuid string) *Room {
	return &Room{
		Roomid:     guid.S(),
		Owner:      uuid,
		PlayerList: make([]*player.WsClient, 0),
	}
}
