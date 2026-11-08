package game

import (
	"math"
	"nanoserver/global"
	"nanoserver/pkg/maputil"
	"nanoserver/protocol/protobuf/pb"
	"sort"
)

type Cell struct {
	location global.PointI
	block    bool
	objects  []IMapObject
}

func (c *Cell) addObject(o IMapObject) {
	if c.objects == nil {
		c.objects = []IMapObject{o}
	} else {
		c.objects = append(c.objects, o)
	}
}
func (c *Cell) removeObject(o IMapObject) bool {
	bDel := false
	var newObjects []IMapObject
	for _, v := range c.objects {
		if v.Uid() == o.Uid() {
			bDel = true
			continue
		}
		newObjects = append(newObjects, v)
	}
	c.objects = newObjects
	return bDel
}

// 清空格子
func (c *Cell) clear() {
	c.objects = nil
}

type MapInfo struct {
	high  int
	wide  int
	cfg   *global.MapConfig
	cells map[global.PointI]*Cell
}

func InitMapInfo(cfg *global.MapConfig) *MapInfo {

	m := &MapInfo{cells: make(map[global.PointI]*Cell)}
	m.cfg = cfg
	m.high = cfg.High
	m.wide = cfg.Wide
	for i, v := range cfg.Cells {
		m.cells[i] = &Cell{location: i, block: !v}
	}
	return m
}

func (m *MapInfo) getCell(point *global.PointI) *Cell {
	cell, ok := m.cells[*point]
	if ok {
		return cell
	}
	return nil
}

func (m *MapInfo) addObject(o IMapObject, location *global.PointI) *Cell {
	dstCell := m.getCell(location)
	if dstCell == nil {
		return nil
	}

	if dstCell.block {
		return nil
	}

	dstCell.addObject(o)
	o.SetCell(dstCell)
	o.SetLocation(&pb.Point{X: float64(dstCell.location.X), Y: float64(dstCell.location.Y)})

	//logger.Debugf("加入地图  MapId=%d , ObjUid:%d , 坐标: %v", m.cfg.ID, o.Uid(), dstCell.location)

	return dstCell
}

func (m *MapInfo) removeObject(o IMapObject, location global.PointI) {

	dstCell := m.getCell(&location)
	if dstCell == nil {
		logger.Warningf("删除物体失败,目标格子nil 地图: %d , Uid %d ", m.cfg.ID, o.Uid())
		return
	}

	if ok := dstCell.removeObject(o); !ok {
		for k, c := range m.cells {
			for _, obj := range c.objects {
				if obj.Uid() == o.Uid() {
					if ok = m.cells[k].removeObject(o); ok {
						return
					}
				}
			}
		}
		if !ok {
			logger.Warningf("删除物体失败 地图: %d , Uid %d  Name %s", m.cfg.ID, o.Uid(), o.Name())
		}
	}
}

// 清空地图
func (m *MapInfo) clear() {
	for _, c := range m.cells {
		c.clear()
	}
}

func (m *MapInfo) locationChange(o IMapObject, dst *pb.Point) bool {
	dstCell := m.getCell(&global.PointI{X: int(dst.X), Y: int(dst.Y)})
	if dstCell == nil {
		return false
	}

	if dstCell.block {
		return false
	}

	if o.GetCell() != nil {
		if !o.GetCell().removeObject(o) {
			logger.Warningf("删除实例失败,整个地图查找删除 地图: %d , Uid %d  ", m.cfg.ID, o.Uid())
			m.removeObject(o, o.GetCell().location)
		}
	}

	dstCell.addObject(o)

	//怪坐标使用整数, 人物坐标使用浮点数
	o.SetLocation(dst)
	if dir := maputil.DirectionFromPointI(&o.GetCell().location, &dstCell.location); dir != pb.Direction_NONE {
		o.SettDirection(dir)
	}

	o.SetCell(dstCell)
	return true
}

func (m *MapInfo) getAvailableCell(location *global.PointI, radius int32) *Cell {

	cells := m.getSquareCells(location, radius)

	// for _, v := range cells {
	// 	if len(v.objects) == 0 {
	// 		return v
	// 	}
	// }
	var cells2 []*Cell
	for _, v := range cells {
		if len(v.objects) == 0 {
			cells2 = append(cells2, v)
		}
	}
	if len(cells2) > 0 {
		randid := global.RandIn(len(cells2))
		return cells2[randid]
	}
	return m.getAvailableCellPvp(location, radius)
}

// pvp 伙伴刷新点可以重叠
func (m *MapInfo) getAvailableCellPvp(location *global.PointI, radius int32) *Cell {
	cells := m.getSquareCells(location, radius)
	if len(cells) > 0 {
		randid := global.RandIn(len(cells))
		return cells[randid]
	}
	return nil
}

// 直线: min max 直线宽度 ，distance 直线长度
func (m *MapInfo) getLineCells(location *global.PointI, direction pb.Direction, distance int32) []*Cell {
	var cells []*Cell

	for i := 1; i <= int(distance); i++ {
		offsetPoint := maputil.Move(location, direction, i)
		cell := m.getCell(offsetPoint)

		if cell == nil {
			continue
		}
		cells = append(cells, cell)

		switch direction {
		case pb.Direction_UP, pb.Direction_RIGHT, pb.Direction_DOWN, pb.Direction_LEFT:
			cell = m.getCell(maputil.Move(offsetPoint, maputil.ShiftDirection(direction, -2), 1))
			cells = append(cells, cell)
			cell = m.getCell(maputil.Move(offsetPoint, maputil.ShiftDirection(direction, 2), 1))
			cells = append(cells, cell)
			// 处理上下左右方向的逻辑
		case pb.Direction_UP_RIGHT, pb.Direction_DOWN_RIGHT, pb.Direction_DOWN_LEFT, pb.Direction_UP_LEFT:
			// 处理对角线方向的逻辑
			cell = m.getCell(maputil.Move(offsetPoint, maputil.ShiftDirection(direction, -1), 1))
			cells = append(cells, cell)
			cell = m.getCell(maputil.Move(offsetPoint, maputil.ShiftDirection(direction, 1), 1))
			cells = append(cells, cell)
		}
	}

	// var strPoints string
	// for _, v := range cells {
	// 	strPoints += fmt.Sprintf("--%v", v.location)
	// }
	// logger.Printf("直线 : 起点=%v  长度=%d  方向=%d 选取区域=%s", location, distance, direction, strPoints)

	return cells
}

// 矩形: radius 半径
func (m *MapInfo) getSquareCells(location *global.PointI, radius int32) []*Cell {
	var cells []*Cell

	if location == nil || m == nil {
		logger.Panicln("getSquareCells faild !!!!!!!!!!")
	}
	for d := 0; d <= int(radius); d++ {
		for y := location.Y - d; y <= location.Y+d; y++ {
			if y < 0 || y >= m.high {
				continue
			}

			step := d * 2
			if int(math.Abs(float64(y-location.Y))) == d {
				step = 1
			}

			for x := location.X - d; x <= location.X+d; x += step {
				if x < 0 || x >= m.wide {
					continue
				}

				point := global.PointI{X: x, Y: y}
				cell := m.cells[point]

				if cell == nil {
					continue
				}

				if cell.block {
					continue
				}

				cells = append(cells, cell)
			}
		}
	}

	// var strPoints string
	// for _, v := range cells {
	// 	strPoints += fmt.Sprintf("--%v", v.location)
	// }
	// logger.Printf("矩形 : 起点=%v  半径=%d   选取区域=%s", location, radius, strPoints)

	return cells
}

// 扇形 location 起点 ，cell: radius半径  direction 朝向
func (m *MapInfo) getSectorCells(location *global.PointI, direction pb.Direction, radius int32) []*Cell {
	var sectorCells []*Cell
	var directions []pb.Direction

	// 根据初始方向计算扇形区域内需要考虑的方向
	for i := 1; i <= int(radius); i++ {
		directions = append(directions, maputil.ShiftDirection(direction, i))
		directions = append(directions, maputil.ShiftDirection(direction, -i))
	}

	// 对每个方向进行遍历，获取该方向上的所有坐标
	for _, dir := range directions {
		for dist := 1; dist <= int(radius); dist++ {
			offsetPoint := maputil.Move(location, dir, dist)

			// 检查坐标点是否在地图范围内
			if offsetPoint.X < 0 || offsetPoint.X >= m.wide || offsetPoint.Y < 0 || offsetPoint.Y >= m.high {
				continue
			}

			// 添加坐标点对应的单元格到结果中
			if cell := m.getCell(offsetPoint); cell != nil && !cell.block {

				sectorCells = append(sectorCells, cell)
			}
		}
	}

	// var strPoints string
	// for _, v := range sectorCells {
	// 	strPoints += fmt.Sprintf("--%v", v.location)
	// }
	// logger.Printf("扇形 : 起点=%v  半径=%d  方向=%d 选取区域=%s", location, radius, direction, strPoints)

	return sectorCells
}

// 三角形 获取相邻两个方向上夹角内所有的坐标, 从左到右
func (m *MapInfo) getAngleCells(location *global.PointI, d1, d2 pb.Direction, distance int32) []*Cell {
	var sectorCells []*Cell
	if d1 == d2 || distance <= 0 {
		return nil
	}
	// 保证 d1 逆时针到 d2 且相邻
	if (d1+1)%8 != d2 && (d2+1)%8 != d1 {
		return nil
	}
	if (d1+1)%8 != d2 { // 让 d1 永远是逆时针侧
		d1, d2 = d2, d1
	}

	set := make(map[global.PointI]struct{})
	v1, v2 := maputil.DirVec[d1], maputil.DirVec[d2]

	// 三角形扫法：层数 1..r
	for i := 1; i <= int(distance); i++ {
		// 从 v1 到 v2 线性插值，步长 1/i
		for t := 0; t <= i; t++ {
			x := location.X + v1.X*t + v2.X*(i-t)
			y := location.Y + v1.Y*t + v2.Y*(i-t)
			if x < 0 || x >= m.wide || y < 0 || y >= m.high {
				continue
			}

			set[global.PointI{X: x, Y: y}] = struct{}{}
		}
	}

	out := make([]global.PointI, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].X != out[j].X {
			return out[i].X < out[j].X
		}
		return out[i].Y < out[j].Y
	})

	for _, offsetPoint := range out {
		if cell := m.getCell(&offsetPoint); cell != nil && !cell.block {
			sectorCells = append(sectorCells, cell)
		}
	}
	return sectorCells
}

func (m *MapInfo) safeArea(location *pb.Point) bool {

	for _, v := range m.cfg.BornPoints {
		if maputil.InRange(location, &pb.Point{X: float64(v.X), Y: float64(v.Y)}, global.GlobalCfg.War.EnterLimit.Bornrange) {
			return true
		}
	}
	return false
}
