package manager

import (
	"context"
	"friberg/internal/game/room"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*room.Room
}

var RM = &RoomManager{
	rooms: make(map[string]*room.Room),
}

// ----------------- 基础接口 -----------------
func (rm *RoomManager) AddRoom(r *room.Room) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.rooms[r.Roomid] = r
}

func (rm *RoomManager) RemoveRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.rooms, roomID)
}

func (rm *RoomManager) GetRoom(roomID string) (*room.Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.rooms[roomID]
	return r, ok
}

func (rm *RoomManager) IsRoomExist(roomID string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	_, ok := rm.rooms[roomID]
	return ok
}
func (rm *RoomManager) SetNewOwner(roomID, uuid string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.rooms[roomID]
	if !ok {
		return false
	}
	r.Owner = uuid
	return true
}
func (rm *RoomManager) GetPlayerList(roomID string) []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	r, ok := rm.rooms[roomID]
	if !ok {
		return nil
	}
	var playerList []string
	for _, p := range r.PlayerList {
		playerList = append(playerList, p.Uuid)
	}
	return playerList
}

func (rm *RoomManager) AddPlayerToRoom(roomID string, uuid string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	r, ok := rm.rooms[roomID]
	if !ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}

	pl := PM.get(uuid)
	if pl == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player not exist")
	}
	r.PlayerList = append(r.PlayerList, pl)
	return nil
}

func (rm *RoomManager) RemovePlayerFromRoom(roomID, uuid string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	r, ok := rm.rooms[roomID]
	if !ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}
	for i, p := range r.PlayerList {
		if p.Uuid == uuid {
			r.PlayerList = append(r.PlayerList[:i], r.PlayerList[i+1:]...)
			break
		}
	}
	return nil
}

func (rm *RoomManager) ShowList() {
	rm.mu.RLock()
	PM.mu.RLock()
	IRPM.mu.RLock()
	defer PM.mu.RUnlock()
	defer IRPM.mu.RUnlock()
	defer rm.mu.RUnlock()
	g.Log().Info(context.TODO(), "room list:")
	for _, r := range rm.rooms {
		g.Log().Info(context.TODO(), "[room]", r.Roomid, r.Owner, len(r.PlayerList))
	}
	g.Log().Info(context.TODO(), "inRoomPlayers:")
	for _, p := range IRPM.inRoomPlayers {
		g.Log().Info(context.TODO(), "[inRoomPlayers]", p.UserName, p.RoomId, p.Uuid)
	}
	g.Log().Info(context.TODO(), "players:")
	for _, p := range PM.players {
		g.Log().Info(context.TODO(), "[players]", p.UserName, p.RoomId, p.Uuid)
	}
}
