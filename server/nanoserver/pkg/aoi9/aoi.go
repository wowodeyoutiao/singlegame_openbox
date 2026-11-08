package aoi9

import "fmt"

// AOIManager AOI区域管理模块
type AOIManager struct {
	// 左
	MinX int
	// 右
	MaxX int
	// X方向格子的数量
	CntsX int
	// 上
	MinY int
	// 下
	MaxY int
	// Y方向格子的数量
	CntsY int
	// 当前区域中有哪些格子Id
	Grids map[int]*Grid
}

// NewAOIManager 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoi := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}
	// 给aoi初始化区域中所有的格子进行编号和初始化
	gh := aoi.gridHeight()
	gl := aoi.gridLength()
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			/*
				这里是关键
			*/
			// 根据x，y编号，计算格子ID：idy*cntsX+x
			gid := y*cntsX + x
			// 初始化gid
			aoi.Grids[gid] = NewGrid(gid, aoi.MinX+x*gh,
				aoi.MinX+(x+1)*gh,
				aoi.MinY+y*gl,
				aoi.MinY+(y+1)*gl)
		}
	}
	return aoi
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridHeight() int {
	return (m.MaxX - m.MinX) / m.CntsX

}

// 得到每个格子在y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子的信息
func (m *AOIManager) String() string {
	// 打印aoi信息
	s := fmt.Sprintf("AOIManager:\n"+
		"MinX:%d,MaxX:%d,CntsX:%d\n"+
		"MinY:%d,MaxX:%d,CntsX:%d\n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	// 打印格子的信息
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundGridsByGid 根据格子GID得到周边就宫格的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := m.Grids[gID]; !ok {
		return nil
	}
	// 初始化返回值数组
	grids = append(grids, m.Grids[gID])
	// 判断gID左边是否有格子、右边是否有格子
	indexX := gID % m.CntsX
	// 需要通过gID得到当前格子X轴的编号 idx:= id % cnx
	// 判断idx编号坐标右边是否还有格子
	if indexX > 0 {
		grids = append(grids, m.Grids[gID-1])
	}
	// 判断idx编号坐标左边是否还有格子
	if indexX < m.CntsX-1 {
		grids = append(grids, m.Grids[gID+1])
	}
	// 遍历一个slice
	for _, grid := range grids {
		if grid.GID/m.CntsY > 0 {
			grids = append(grids, m.Grids[grid.GID-5])
		}
		if grid.GID/m.CntsY < m.CntsY-1 {
			grids = append(grids, m.Grids[grid.GID+5])
		}
	}
	return
}

// GetPidsByPos 通过横纵坐标得到周边9宫格内全部的PlayersIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的GID格子id
	gID := m.GetGidByPos(x, y)
	// 通过GID得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gID)
	// 将九宫格的信息里的全部的Player的id累加到playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		//fmt.Printf("========> grid ID:%d ,pid:%v <===========", grid.GID, grid.GetPlayerIDs())
	}
	return
}

// GetGidByPos 通过x、y横纵轴坐标得到当前的GID格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridLength()
	idy := (int(y) - m.MinY) / m.gridLength()
	return idy*m.CntsX + idx
}

// AddPidToGrid 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].Add(pID)
}

// RemovePidFromGrid 移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].Delete(pID)
}

// GetPidsByGid 通过GID获得全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetPlayerIDs()
	return
}

// AddToGridByPos 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.Add(pID)
}

// RemoveFromGridByPos 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.Delete(pID)
}
