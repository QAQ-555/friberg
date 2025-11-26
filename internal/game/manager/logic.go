package manager

import (
	"context"
	"friberg/internal/consts"
	"friberg/internal/game/player"
	"friberg/internal/game/room"
	isocket "friberg/internal/library/websocket"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

//
// ========== 公共工具函数 ==========
//

// 判断玩家是否合法存在
func validatePlayer(uuid string) error {
	if !PM.IsPlayerExist(uuid) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player not exist")
	}
	return nil
}

// 构建房间列表输出格式
func buildRoomList() []room.OutRoom {
	RM.mu.RLock()
	defer RM.mu.RUnlock()

	roomList := make([]room.OutRoom, 0, len(RM.rooms))
	for _, r := range RM.rooms {
		owner := PM.get(r.Owner)
		roomList = append(roomList, room.OutRoom{
			Roomid:    r.Roomid,
			Owner:     owner.UserName,
			NumPlayer: len(r.PlayerList),
		})
	}
	return roomList
}

// 推送消息给某个玩家
func pushToPlayer(uuid string, msg isocket.Message) error {
	ws := PM.GetWs(uuid)
	if ws == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player ws not exist")
	}
	return ws.WriteJSON(msg)
}

//
// ========== 业务逻辑部分 ==========
//

// 创建房间
func CreateRoomBusiness(r *room.Room) error {
	g.Log().Infof(context.TODO(), "player %s create room:%s", r.Owner, r.Roomid)

	if err := validatePlayer(r.Owner); err != nil {
		return err
	}

	if IRPM.IsExist(r.Owner) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player already in room")
	}

	RM.AddRoom(r)
	BroadcastRoomList()
	return nil
}

// 删除房间
func DeleteRoomBusiness(roomID string) error {
	g.Log().Infof(context.TODO(), "delete room:%s", roomID)

	if !RM.IsRoomExist(roomID) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "room not exist")
	}

	players := RM.GetPlayerList(roomID)
	for _, uuid := range players {
		IRPM.Leave(uuid)
		pushToPlayer(uuid, isocket.Message{MsgType: consts.MsgType_RoomExit})
	}

	RM.RemoveRoom(roomID)
	BroadcastRoomList()
	return nil
}

// 玩家加入房间
func AddPlayerToRoomBusiness(roomID, uuid string) error {
	g.Log().Infof(context.TODO(), "add player %s to room %s", uuid, roomID)

	if err := validatePlayer(uuid); err != nil {
		return err
	}

	if IRPM.IsExist(uuid) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player already in room")
	}

	if err := RM.AddPlayerToRoom(roomID, uuid); err != nil {
		return err
	}

	IRPM.AddPl2Room(uuid, roomID)
	BroadcastPlayerList(roomID)
	return nil
}

// 玩家退出房间（主动）
func PlayerExitRoom(uuid string) error {
	g.Log().Infof(context.TODO(), "player %s exit room", uuid)

	if err := validatePlayer(uuid); err != nil {
		return err
	}

	if !IRPM.IsExist(uuid) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "player not in room")
	}

	roomID := IRPM.GetRoomId(uuid)
	RM.RemovePlayerFromRoom(roomID, uuid)
	IRPM.Leave(uuid)

	BroadcastPlayerList(roomID)
	return nil
}

// 玩家断线
func RemovePlayerBusiness(uuid string) error {
	g.Log().Infof(context.TODO(), "player %s disconnect", uuid)

	if err := validatePlayer(uuid); err != nil {
		return err
	}

	if IRPM.IsExist(uuid) {
		roomID := IRPM.GetRoomId(uuid)

		RM.RemovePlayerFromRoom(roomID, uuid)

		// 房间没人了就删房
		if len(RM.GetPlayerList(roomID)) == 0 {
			DeleteRoomBusiness(roomID)
		} else {
			RM.SetNewOwner(roomID, RM.GetPlayerList(roomID)[0])
			BroadcastPlayerList(roomID)
		}
		IRPM.Leave(uuid)
	}

	PM.Delete(uuid)
	return nil
}

//
// ========== 广播部分 ==========
//

// 广播房间列表
func BroadcastRoomList() {
	roomList := buildRoomList()
	msg := isocket.Message{
		MsgType: consts.MsgType_RoomList,
		Data:    room.OutRoomlist{RoomList: roomList},
	}

	for _, p := range PM.players {
		pushToPlayer(p.Uuid, msg)
	}
}

// 广播房间内玩家列表
func BroadcastPlayerList(roomID string) {
	players := RM.GetPlayerList(roomID)
	if len(players) == 0 {
		return
	}

	names := make([]string, 0, len(players))
	for _, uuid := range players {
		names = append(names, PM.get(uuid).UserName)
	}

	msg := isocket.Message{
		MsgType: consts.MsgType_RoomPlayerList,
		Data:    player.OutClient{UserName: names},
	}

	for _, uuid := range players {
		if err := pushToPlayer(uuid, msg); err != nil {
			g.Log().Error(context.TODO(), "write json error", err)
			PM.EndCtx(uuid)
		}
	}
}

// 给某个玩家发送房间列表
func SendRoomListToPlayer(uuid string) error {
	if err := validatePlayer(uuid); err != nil {
		return err
	}

	msg := isocket.Message{
		MsgType: consts.MsgType_RoomList,
		Data:    room.OutRoomlist{RoomList: buildRoomList()},
	}

	return pushToPlayer(uuid, msg)
}
