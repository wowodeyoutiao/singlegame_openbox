package game

import (
	"fmt"
	"nanoserver/global"
	"nanoserver/pkg/utils"
	"sync"

	"github.com/lonng/nano"
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
)

var playerMgr = NewPlayerManager()

// PlayerManager 是游戏玩家管理器，负责管理所有玩家和分组
type PlayerManager struct {
	slicePool sync.Pool
	playerMap *utils.SafeMap[int64, *Player] // 线程安全的玩家管理映射，键为玩家UID，值为玩家对象
	groupMap  map[int32]*nano.Group          // 服务器分组管理，键为服务器ID，值为nano.Group
}

// NewPlayerManager 创建一个新的玩家管理器
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		playerMap: utils.NewSafeMap[int64, *Player](), // 初始化线程安全的映射
		groupMap:  make(map[int32]*nano.Group, 0),     // 初始化分组映射
	}
}

// Init 初始化玩家管理器，创建服务器分组
func (m *PlayerManager) Init() {
	m.slicePool.New = func() interface{} {
		return make([]*Player, 0, 200)
	}
	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
		// 为每个服务器ID创建一个分组，用于广播消息
		m.groupMap[i] = nano.NewGroup(fmt.Sprintf("%d_SYSTEM_MESSAGE_BROADCAST", i))
	}
}

// GetPlayer 根据玩家UID获取玩家对象
func (m *PlayerManager) GetPlayer(uid int64) (*Player, bool) {
	return m.playerMap.Get(uid)
}

// PutPlayer 添加或更新玩家对象到玩家管理映射中
func (m *PlayerManager) PutPlayer(player *Player) error {
	if _, ok := m.playerMap.Get(player.uid); ok {
		log.Warnf("玩家已经存在，正在覆盖玩家， UID=%d", player.uid)
	}
	m.playerMap.Set(player.uid, player)
	return nil
}

// RemovePlayer 从玩家管理映射中删除玩家对象
func (m *PlayerManager) RemovePlayer(uid int64) error {
	if _, ok := m.playerMap.Get(uid); ok {
		m.playerMap.Delete(uid)
		log.Infof("玩家: %d 从在线列表中删除, 剩余：%d", uid, m.playerMap.Len())
	}
	return nil
}

// IteratePlayers 遍历玩家
func (m *PlayerManager) IteratePlayers(f func(player *Player) bool) {
	if m.playerMap.Len() == 0 {
		return
	}
	//避免死锁
	tempSlice := m.slicePool.Get().([]*Player)
	tempSlice = tempSlice[:0]
	m.playerMap.Range(func(key int64, val *Player) bool {
		tempSlice = append(tempSlice, val)
		return true
	})
	for _, player := range tempSlice {
		if !f(player) {
			break
		}
	}
	m.slicePool.Put(tempSlice)
}

// Len 返回当前在线玩家的数量
func (m *PlayerManager) Len() int {
	return m.playerMap.Len()
}

// AddSessionToGroup 将玩家的会话添加到对应的分组中
func (m *PlayerManager) AddSessionToGroup(session *session.Session) {
	player, ok := m.playerMap.Get(session.UID())
	if !ok {
		log.Errorf("添加session找不绑定玩家,  UID=%d", session.UID())
		return
	}
	group, ok := m.groupMap[player.serverID]
	if !ok {
		return
	}
	err := group.Add(session)
	if err != nil {
		return
	}
}

// RemoveSessionFromGroup 将玩家的会话从对应的分组中移除
func (m *PlayerManager) RemoveSessionFromGroup(session *session.Session) {
	player, ok := m.playerMap.Get(session.UID())
	if !ok {
		log.Errorf("移除session找不绑定玩家, UID=%d", session.UID())
		return
	}
	group, ok := m.groupMap[player.serverID]
	if !ok {
		return
	}
	err := group.Leave(session)
	if err != nil {
		return
	}
}

// BroadcastToAllGroups 向所有服务器分组广播消息
func (m *PlayerManager) BroadcastToAllGroups(route string, msg interface{}) {
	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
		m.BroadcastToGroupByServerID(i, route, msg)
	}
}

// BroadcastToGroupByServerID 向指定服务器分组广播消息
func (m *PlayerManager) BroadcastToGroupByServerID(sid int32, route string, msg interface{}) error {
	group, _ := m.groupMap[sid]
	if group == nil {
		return nil
	}
	err := group.Broadcast(route, msg)
	if err != nil {
		return err
	}
	return nil
}
