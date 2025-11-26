package manager

import (
	"context"
	"friberg/internal/game/player"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

type PlayerManager struct {
	mu      sync.RWMutex
	players map[string]*player.WsClient
}
type InRoomPlayerManager struct {
	inRoomPlayers map[string]*player.WsClient
	mu            sync.RWMutex
}

var PM = &PlayerManager{
	players: make(map[string]*player.WsClient),
}

var IRPM = &InRoomPlayerManager{
	inRoomPlayers: make(map[string]*player.WsClient),
}

func (pm *PlayerManager) Add(player *player.WsClient) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.players[player.Uuid] = player
}
func (pm *PlayerManager) Get(uuid string) *player.WsClient {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.players[uuid]
}
func (pm *PlayerManager) IsPlayerExist(uuid string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	_, ok := pm.players[uuid]
	return ok
}

func (pm *PlayerManager) Remove(uuid string) {
	g.Log().Info(context.TODO(), "Remove player", uuid)
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pl := pm.players[uuid]
	if pl.RoomId != "" {
		delete(IRPM.inRoomPlayers, pl.RoomId)
		RM.RemovePlayer(pl.RoomId, uuid)
	}
	delete(pm.players, uuid)
}
