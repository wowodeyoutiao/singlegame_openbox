package astar

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type OpenList []*_AstarPoint

func (ol OpenList) Len() int           { return len(ol) }
func (ol OpenList) Less(i, j int) bool { return ol[i].fVal < ol[j].fVal }
func (ol OpenList) Swap(i, j int)      { ol[i], ol[j] = ol[j], ol[i] }

func (ol *OpenList) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*ol = append(*ol, x.(*_AstarPoint))
}

func (ol *OpenList) Pop() interface{} {
	old := *ol
	n := len(old)
	x := old[n-1]
	*ol = old[0 : n-1]
	return x
}

type _Point struct {
	x    int
	y    int
	view string
}

//========================================================================================

// 保存地图的基本信息
type Map struct {
	points [][]_Point
	blocks map[string]*_Point
	maxX   int
	maxY   int
}

func NewMap(charMap []string) (m Map) {
	m.points = make([][]_Point, len(charMap))
	m.blocks = make(map[string]*_Point, len(charMap)*2)
	for x, row := range charMap {
		cols := strings.Split(row, " ")
		m.points[x] = make([]_Point, len(cols))
		for y, view := range cols {
			m.points[x][y] = _Point{x, y, view}
			if view == "X" {
				m.blocks[pointAsKey(x, y)] = &m.points[x][y]
			}
		} // end of cols
	} // end of row

	m.maxX = len(m.points)
	m.maxY = len(m.points[0])

	return m
}

func NewMapEx(high, wide int) (m Map) {
	m.points = make([][]_Point, high*wide)
	m.blocks = make(map[string]*_Point, high*wide*2)

	for ix := 0; ix < high; ix++ {
		m.points[ix] = make([]_Point, wide)
		for iy := 0; iy < wide; iy++ {
			m.points[ix][iy] = _Point{ix, iy, "."}
		}
	}
	m.maxX = high
	m.maxY = wide
	return m
}

func NewMapExx(high, wide int, charMap []string) (m Map) {
	m.points = make([][]_Point, len(charMap))
	m.blocks = make(map[string]*_Point, len(charMap)*2)

	for x := 0; x < wide; x++ {
		m.points[x] = make([]_Point, high)
		for y := 0; y < high; y++ {

			view := charMap[x*high+y]
			m.points[x][y] = _Point{x, y, view}
			if view == "1" {
				m.blocks[pointAsKey(x, y)] = &m.points[x][y]
			}
		}
	}
	m.maxX = wide
	m.maxY = high

	return m
}

func (m *Map) getAdjacentPoint(curPoint *_Point) (adjacents []*_Point) {
	if x, y := curPoint.x, curPoint.y-1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y-1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y+1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x, curPoint.y+1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y+1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y-1; x >= 0 && x < m.maxX && y >= 0 && y < m.maxY {
		adjacents = append(adjacents, &m.points[x][y])
	}
	return adjacents
}

func (m *Map) PrintMap(path *SearchRoad) {
	fmt.Println("map's border:", m.maxX, m.maxY)
	for x := 0; x < m.maxX; x++ {
		for y := 0; y < m.maxY; y++ {
			if path != nil {
				if x == path.start.x && y == path.start.y {
					fmt.Print("S")
					goto NEXT
				}
				if x == path.end.x && y == path.end.y {
					fmt.Print("E")
					goto NEXT
				}
				for i := 0; i < len(path.TheRoad); i++ {
					if path.TheRoad[i].x == x && path.TheRoad[i].y == y {
						fmt.Print("*")
						goto NEXT
					}
				}
			}
			fmt.Print(m.points[x][y].view)
		NEXT:
		}
		fmt.Println()
	}
}

func pointAsKey(x, y int) (key string) {
	key = strconv.Itoa(x) + "," + strconv.Itoa(y)
	return key
}

func (m *Map) GetVisionArea(center APoint) []*APoint {
	visionArea := []*APoint{}
	for x := center.X - 1; x <= center.X+1; x++ {
		for y := center.Y - 1; y <= center.Y+1; y++ {
			if x >= 0 && x < m.maxX && y >= 0 && y < m.maxY { // 假设九宫格从0到8
				if x != center.X || y != center.Y { // 排除中心点本身
					visionArea = append(visionArea, &APoint{X: x, Y: y})
				}
			}
		}
	}
	return visionArea
}

//========================================================================================

type _AstarPoint struct {
	_Point
	father *_AstarPoint
	gVal   int
	hVal   int
	fVal   int
}

func NewAstarPoint(p *_Point, father *_AstarPoint, end *_AstarPoint) (ap *_AstarPoint) {
	ap = &_AstarPoint{*p, father, 0, 0, 0}
	if end != nil {
		ap.calcFVal(end)
	}
	return ap
}

func (asp *_AstarPoint) calcGVal() int {
	if asp.father != nil {
		deltaX := math.Abs(float64(asp.father.x - asp.x))
		deltaY := math.Abs(float64(asp.father.y - asp.y))
		if deltaX == 1 && deltaY == 0 {
			asp.gVal = asp.father.gVal + 10
		} else if deltaX == 0 && deltaY == 1 {
			asp.gVal = asp.father.gVal + 10
		} else if deltaX == 1 && deltaY == 1 {
			asp.gVal = asp.father.gVal + 14
		} else {
			panic("father point is invalid!")
		}
	}
	return asp.gVal
}

func (asp *_AstarPoint) calcHVal(end *_AstarPoint) int {
	asp.hVal = int(math.Abs(float64(end.x-asp.x)) + math.Abs(float64(end.y-asp.y)))
	return asp.hVal
}

func (asp *_AstarPoint) calcFVal(end *_AstarPoint) int {
	asp.fVal = asp.calcGVal() + asp.calcHVal(end)
	return asp.fVal
}

//========================================================================================

type SearchRoad struct {
	theMap  *Map
	start   _AstarPoint
	end     _AstarPoint
	closeLi map[string]*_AstarPoint
	openLi  OpenList
	openSet map[string]*_AstarPoint
	TheRoad []*_AstarPoint
}

type APoint struct {
	X int
	Y int
}

func NewSearchRoad(startx, starty, endx, endy int, m *Map) *SearchRoad {
	sr := &SearchRoad{}
	sr.theMap = m
	sr.start = *NewAstarPoint(&_Point{startx, starty, "S"}, nil, nil)
	sr.end = *NewAstarPoint(&_Point{endx, endy, "E"}, nil, nil)
	sr.TheRoad = make([]*_AstarPoint, 0)
	sr.openSet = make(map[string]*_AstarPoint, m.maxX+m.maxY)
	sr.closeLi = make(map[string]*_AstarPoint, m.maxX+m.maxY)

	heap.Init(&sr.openLi)
	heap.Push(&sr.openLi, &sr.start) // 首先把起点加入开放列表
	sr.openSet[pointAsKey(sr.start.x, sr.start.y)] = &sr.start
	// 将障碍点放入关闭列表
	for k, v := range m.blocks {
		sr.closeLi[k] = NewAstarPoint(v, nil, nil)
	}

	return sr
}

func (sr *SearchRoad) FindoutRoad() (b bool, or []*APoint) {
	for len(sr.openLi) > 0 {
		// 将节点从开放列表移到关闭列表当中。
		x := heap.Pop(&sr.openLi)
		curPoint := x.(*_AstarPoint)
		delete(sr.openSet, pointAsKey(curPoint.x, curPoint.y))
		sr.closeLi[pointAsKey(curPoint.x, curPoint.y)] = curPoint

		//fmt.Println("curPoint :", curPoint.x, curPoint.y)
		adjacs := sr.theMap.getAdjacentPoint(&curPoint._Point)
		for _, p := range adjacs {
			//fmt.Println("\t adjact :", p.x, p.y)
			theAP := NewAstarPoint(p, curPoint, &sr.end)
			if pointAsKey(theAP.x, theAP.y) == pointAsKey(sr.end.x, sr.end.y) {
				// 找出路径了, 标记路径
				for theAP.father != nil {
					sr.TheRoad = append(sr.TheRoad, theAP)
					theAP.view = "*"
					theAP = theAP.father
					or = append(or, &APoint{X: theAP.x, Y: theAP.y})
				}
				return true, or
			}

			_, ok := sr.closeLi[pointAsKey(p.x, p.y)]
			if ok {
				continue
			}

			existAP, ok := sr.openSet[pointAsKey(p.x, p.y)]
			if !ok {
				heap.Push(&sr.openLi, theAP)
				sr.openSet[pointAsKey(theAP.x, theAP.y)] = theAP
			} else {
				oldGVal, oldFather := existAP.gVal, existAP.father
				existAP.father = curPoint
				existAP.calcGVal()
				// 如果新的节点的G值还不如老的节点就恢复老的节点
				if existAP.gVal > oldGVal {
					// restore father
					existAP.father = oldFather
					existAP.gVal = oldGVal
				}
			}

		}
	}

	return false, nil
}

//========================================================================================

func main() {
	presetMap := []string{
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		"X . X X X X X X X X X X X X X X X X X X X X X X X X X",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		"X X X X X X X X X X X X . X X X X X X X X X X X . X X",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
	}
	m := NewMap(presetMap)
	m.PrintMap(nil)

	searchRoad := NewSearchRoad(0, 0, 14, 0, &m)
	if b, _ := searchRoad.FindoutRoad(); b {
		fmt.Println("找到了， 你看！")
		m.PrintMap(searchRoad)
	} else {
		fmt.Println("找不到路径！")
	}
}
