package manager

import (
	"context"
	"friberg/internal/consts"
	"friberg/internal/game/player"
	"friberg/internal/game/room"
	isocket "friberg/internal/library/websocket"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type RoomManager struct {
	mu    sync.Mutex
	rooms map[string]*room.Room
}

var RM = &RoomManager{
	rooms: make(map[string]*room.Room),
}

// 创建房间
func (rm *RoomManager) CreateRoom(room *room.Room) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	// 检查创建者是否已在其他房间
	if _, ok := PM.inRoomPlayers[room.Owner]; ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player already in room")
	}
	//添加房间对象
	rm.rooms[room.Roomid] = room
	//广播房间列表
	rm.BroadcastRoomList()
	return nil
}

// 删除房间
func (rm *RoomManager) DeleteRoom(roomID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	//检查房间是否存在
	if _, ok := rm.rooms[roomID]; !ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}
	for _, p := range rm.rooms[roomID].PlayerList {
		//更新玩家房间ID
		PM.players[p.Uuid].RoomId = ""
		delete(PM.inRoomPlayers, p.Uuid)
		msg := isocket.Message{
			MsgType: consts.MsgType_RoomExit,
			Data:    nil,
		}
		p.Ws.WriteJSON(msg)
	}
	delete(rm.rooms, roomID)
	//广播房间列表
	rm.BroadcastRoomList()
	return nil
}

// 添加玩家到房间
func (rm *RoomManager) AddPlayer(roomID, uuid string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	//检查房间是否存在
	if _, ok := rm.rooms[roomID]; !ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}
	//检查玩家是否已在其他房间
	if _, ok := PM.inRoomPlayers[uuid]; ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player already in room")
	}
	//添加玩家到在房间内玩家列表
	PM.inRoomPlayers[uuid] = PM.players[uuid]
	//添加玩家到房间内玩家列表
	r := rm.rooms[roomID]
	r.PlayerList = append(r.PlayerList, PM.players[uuid])
	//广播房间内玩家列表
	rm.BroadcastPlayer(r)
	//更新玩家房间ID
	PM.players[uuid].RoomId = roomID
	return nil
}

// 从房间移除玩家
func (rm *RoomManager) RemovePlayer(roomID, uuid string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	//检查房间是否存在
	if _, ok := rm.rooms[roomID]; !ok {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}
	r := rm.rooms[roomID]
	//检查将玩家从房间内玩家列表移除
	for i, p := range r.PlayerList {
		if p.Uuid == uuid {
			r.PlayerList = append(r.PlayerList[:i], r.PlayerList[i+1:]...)
			break
		}
	}
	//广播房间内玩家列表
	rm.BroadcastPlayer(r)
	//更新玩家房间ID
	PM.players[uuid].RoomId = ""
	delete(PM.inRoomPlayers, uuid)
	return nil
}

// BroadcastRoom 向所有在线玩家广播房间列表信息
// 该函数会遍历所有在线玩家，并将当前的房间列表数据发送给他们
// 使用互斥锁确保线程安全，防止并发访问rooms数据
func (rm *RoomManager) BroadcastRoomList() {
	roomList := room.OutRoomlist{}
	for _, r := range rm.rooms {
		owner := PM.players[r.Owner]
		g.Log().Info(context.TODO(), r)
		outRoom := room.OutRoom{
			Roomid:    r.Roomid,
			Owner:     owner.UserName,
			NumPlayer: len(r.PlayerList),
		}
		roomList.RoomList = append(roomList.RoomList, outRoom)
	}
	for _, c := range PM.players {
		c.Ws.WriteJSON(isocket.Message{
			MsgType: consts.MsgType_RoomList,
			Data:    roomList,
		})
	}
}

func (rm *RoomManager) BroadcastPlayer(r *room.Room) {
	if len(r.PlayerList) == 0 {
		return
	}

	// 构建 UUID 列表
	userNames := make([]string, 0, len(r.PlayerList))
	for _, c := range r.PlayerList {
		userNames = append(userNames, c.UserName)
	}

	msg := isocket.Message{
		MsgType: consts.MsgType_RoomPlayerList,
		Data: player.OutClient{
			UserName: userNames,
		},
	}

	// 发送消息
	for _, c := range r.PlayerList {
		err := c.Ws.WriteJSON(msg)
		if err != nil {
			g.Log().Error(c.Ctx, "write json error", err)
			c.Cancel()
		}
	}
}
