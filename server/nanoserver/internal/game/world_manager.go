package game

import (
	"fmt"
	"nanoserver/global"
	"nanoserver/protocol"
)

type IWorld interface {
	GetName() string
	NewDesk(roomNo int64, roomID int32) *Desk
}

// 游戏世界管理器
type WorldMgr struct {
	IWorld
	worlds map[protocol.GameType]IWorld
}

var worldMgr = NewWorldManager()

func NewWorldManager() *WorldMgr {
	return &WorldMgr{}
}

// 根据配置,初始化游戏世界管理器
func (w *WorldMgr) InitWorld() {

	w.worlds = make(map[protocol.GameType]IWorld)

	for _, v := range global.RoomCfg {
		switch v.GameType {
		case protocol.GameTypeMow:
			_, ok := w.worlds[v.GameType]
			if !ok {
				w.worlds[v.GameType] = NewMowWorld()
			}
		default:
			panic(fmt.Sprintf("invalid game type [%d]", v.GameType))
		}
	}
}

func (w *WorldMgr) GetWorld(t protocol.GameType) IWorld {
	wd, ok := w.worlds[t]
	if ok {
		return wd
	}
	return nil
}
