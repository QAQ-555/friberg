package manager

import (
	"friberg/internal/game/player"
	"sync"
)

type PlayerManager struct {
	mu            sync.Mutex
	players       map[string]*player.WsClient
	inRoomPlayers map[string]*player.WsClient
}

var PM = &PlayerManager{
	players:       make(map[string]*player.WsClient),
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
