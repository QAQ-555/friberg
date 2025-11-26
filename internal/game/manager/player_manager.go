package manager

import (
	"friberg/internal/game/player"
	"sync"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	mu      sync.RWMutex
	players map[string]*player.WsClient
}

type InRoomPlayerManager struct {
	mu            sync.RWMutex
	inRoomPlayers map[string]*player.WsClient
}

var PM = &PlayerManager{
	players: make(map[string]*player.WsClient),
}

var IRPM = &InRoomPlayerManager{
	inRoomPlayers: make(map[string]*player.WsClient),
}

// ----------------- 基础接口 -----------------
func (pm *PlayerManager) Add(player *player.WsClient) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.players[player.Uuid] = player
}

func (pm *PlayerManager) Delete(uuid string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.players, uuid)
}

func (pm *PlayerManager) EndCtx(uuid string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.players[uuid].Cancel()
}

func (pm *PlayerManager) get(uuid string) *player.WsClient {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.players[uuid]
}

func (pm *PlayerManager) IsPlayerExist(uuid string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, ok := pm.players[uuid]
	return ok
}

func (pm *PlayerManager) Remove(uuid string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.players, uuid)
}
func (pm *PlayerManager) GetWs(uuid string) *websocket.Conn {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	pl := pm.players[uuid]
	if pl == nil {
		return nil
	}
	return pl.Ws
}
func (irpm *InRoomPlayerManager) AddPl2Room(uuid string, roomID string) {
	irpm.mu.Lock()
	defer irpm.mu.Unlock()
	pl := PM.get(uuid)
	if pl == nil {
		return
	}
	pl.RoomId = roomID
	irpm.inRoomPlayers[uuid] = pl
}

func (irpm *InRoomPlayerManager) Leave(uuid string) {
	irpm.mu.Lock()
	defer irpm.mu.Unlock()
	irpm.inRoomPlayers[uuid].RoomId = ""
	delete(irpm.inRoomPlayers, uuid)
}

func (irpm *InRoomPlayerManager) IsExist(uuid string) bool {
	irpm.mu.RLock()
	defer irpm.mu.RUnlock()
	_, ok := irpm.inRoomPlayers[uuid]
	return ok
}

func (irpm *InRoomPlayerManager) GetRoomId(uuid string) string {
	irpm.mu.RLock()
	defer irpm.mu.RUnlock()
	return irpm.inRoomPlayers[uuid].RoomId
}
