package aoi9

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子类型
type Grid struct {
	// 格子ID
	GID int
	// 左界
	MinX int
	// 右界
	MaxX int
	// 上界
	MinY int
	// 下界
	MaxY int
	// 当前玩家集合
	playerIDs map[int]bool
	// 保护锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化格子的方法
func NewGrid(id, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       id,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: map[int]bool{},
	}
}

// Add 格子添加玩家/物品
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true

}

// Delete 格子删除玩家/物品
func (g *Grid) Delete(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 获取所有玩家ID
func (g *Grid) GetPlayerIDs() (ids []int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	for k, _ := range g.playerIDs {
		ids = append(ids, k)
	}
	return
}

// 打印格子信息(调试)
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d，minX:%d，"+
		"maxX:%d，minY:%d，"+
		"maxY:%d，playerIDs:%v\n", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)

}
