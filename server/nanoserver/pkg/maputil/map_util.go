package maputil

import (
	"fmt"
	"math"
	"nanoserver/global"
	"nanoserver/pkg/astar"
	"nanoserver/protocol/protobuf/pb"
)

// 8 方向单位向量（dx,dy），索引与常量对应
var DirVec = [8]global.PointI{
	{X: 0, Y: 1},   // UP
	{X: 1, Y: 1},   // UP_RIGHT
	{X: 1, Y: 0},   // RIGHT
	{X: 1, Y: -1},  // DOWN_RIGHT
	{X: 0, Y: -1},  // DOWN
	{X: -1, Y: -1}, // DOWN_LEFT
	{X: -1, Y: 0},  // LEFT
	{X: -1, Y: 1},  // UP_LEFT
}

// Distance 计算两个 APoint 点之间的距离
func Distance(a *pb.Point, b *pb.Point) int {
	// 直接计算差值的最大绝对值，避免重复计算
	if a == nil || b == nil {
		return math.MaxInt32
	}
	maxAbsDiff := math.Max(math.Abs(float64(a.X-b.X)), math.Abs(float64(a.Y-b.Y)))
	// 将 float64 的绝对值最大差转换为 int
	return int(maxAbsDiff)
}

// InRange 检查点a是否在点b周围指定的整数距离i内
func InRange(a *pb.Point, b *pb.Point, i int) bool {
	return inRange(int(a.X), int(a.Y), int(b.X), int(b.Y), i)
}

// inRange 是内部辅助函数，接受坐标的整数形式
func inRange(x1, y1, x2, y2, i int) bool {
	return math.Abs(float64(x1-x2)) <= float64(i) && math.Abs(float64(y1-y2)) <= float64(i)
}

// DirectionFromPoint 根据起点和终点返回相应的方向
func DirectionFromPointF(source, dest *pb.Point) pb.Direction {
	if source.X < dest.X {
		if source.Y < dest.Y {
			return pb.Direction_UP_RIGHT
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN_RIGHT
		} else {
			return pb.Direction_RIGHT
		}
	} else if source.X > dest.X {
		if source.Y < dest.Y {
			return pb.Direction_UP_LEFT
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN_LEFT
		} else {
			return pb.Direction_LEFT
		}
	} else {
		if source.Y < dest.Y {
			return pb.Direction_UP
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN
		} else {
			return pb.Direction_NONE // 如果起点和终点相同，则表示无方向
		}
	}
}

func DirectionFromPointI(source, dest *global.PointI) pb.Direction {
	if source.X < dest.X {
		if source.Y < dest.Y {
			return pb.Direction_UP_RIGHT
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN_RIGHT
		} else {
			return pb.Direction_RIGHT
		}
	} else if source.X > dest.X {
		if source.Y < dest.Y {
			return pb.Direction_UP_LEFT
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN_LEFT
		} else {
			return pb.Direction_LEFT
		}
	} else {
		if source.Y < dest.Y {
			return pb.Direction_UP
		} else if source.Y > dest.Y {
			return pb.Direction_DOWN
		} else {
			return pb.Direction_NONE // 如果起点和终点相同，则表示无方向
		}
	}
}

// 顶点在左上角
func Move1(location *global.PointI, direction pb.Direction, distance int) *global.PointI {
	switch direction {
	case pb.Direction_UP:
		location.Y -= distance
	case pb.Direction_UP_RIGHT:
		location.X += distance
		location.Y -= distance
	case pb.Direction_RIGHT:
		location.X += distance
	case pb.Direction_DOWN_RIGHT:
		location.X += distance
		location.Y += distance
	case pb.Direction_DOWN:
		location.Y += distance
	case pb.Direction_DOWN_LEFT:
		location.X -= distance
		location.Y += distance
	case pb.Direction_LEFT:
		location.X -= distance
	case pb.Direction_UP_LEFT:
		location.X -= distance
		location.Y -= distance
	case pb.Direction_NONE:
		// 无方向，不进行移动
	default:
		fmt.Printf("Unknown direction: %d\n", direction)
	}
	return location
}

// 顶点在左下角
// <param name="location">原始坐标</param>
// <param name="direction">方向</param>
// <param name="distance">距离</param>
// <returns>偏移后的坐标</returns>
func Move(location *global.PointI, direction pb.Direction, distance int) *global.PointI {
	shiftLocation := *location
	switch direction {
	case pb.Direction_UP:
		shiftLocation.Y += distance
	case pb.Direction_UP_RIGHT:
		shiftLocation.X += distance
		shiftLocation.Y += distance
	case pb.Direction_RIGHT:
		shiftLocation.X += distance
	case pb.Direction_DOWN_RIGHT:
		shiftLocation.X += distance
		shiftLocation.Y -= distance
	case pb.Direction_DOWN:
		shiftLocation.Y -= distance
	case pb.Direction_DOWN_LEFT:
		shiftLocation.X -= distance
		shiftLocation.Y -= distance
	case pb.Direction_LEFT:
		shiftLocation.X -= distance
	case pb.Direction_UP_LEFT:
		shiftLocation.X -= distance
		shiftLocation.Y += distance
	case pb.Direction_NONE:
		// 无方向，不进行移动
	default:
		fmt.Printf("Unknown direction: %d\n", direction)
	}
	return &shiftLocation
}

// ShiftDirection 根据给定的方向和整数i，返回偏移后的方向
func ShiftDirection(dir pb.Direction, i int) pb.Direction {
	// 由于MowDirection是int的别名，可以直接进行整数运算
	// 加8确保在环状方向上正确地环绕，% 8确保结果仍然在枚举值范围内
	return pb.Direction((int(dir) + i + 8) % 8)
}

// 根据攻击属性获取对应的防御属性
func DefenseAttributes(attr pb.Attribute) pb.Attribute {
	if attr == pb.Attribute_PHYSICAL_ATTACK {
		return pb.Attribute_PHYSICAL_DEFENSE
	} else if attr == pb.Attribute_MAGIC_ATTACK {
		return pb.Attribute_MAGIC_DEFENSE
	}
	return pb.Attribute_INVALID_Attribute
}

func GetPathAstar(startP *pb.Point, endP *pb.Point, aMap *astar.Map) []*astar.APoint {

	astarRoad := []*astar.APoint{}
	p := astar.NewSearchRoad(int(startP.X), int(startP.Y), int(endP.X), int(endP.Y), aMap)
	if b, road := p.FindoutRoad(); b {
		for i, j := 0, len(road)-1; i < j; i, j = i+1, j-1 {
			road[i], road[j] = road[j], road[i]
		}
		astarRoad = append(astarRoad, road...)
	}
	return astarRoad
}

func ABS(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// 随机一个相邻点
func RandomNeighbor(x, y int) *astar.APoint {
	// 8 方向偏移
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	i := global.RandIn(8)
	return &astar.APoint{X: x + dx[i], Y: y + dy[i]}
}

// 计算移动方向（-1, 0, 1）
func direction(a, b int) int {
	if a < b {
		return 1
	}
	if a > b {
		return -1
	}
	return 0
}

// 计算两点之间的最短路径（支持斜向移动）
func GetPath(startP *pb.Point, endP *pb.Point, _ *astar.Map) []*astar.APoint {
	path := []*astar.APoint{}
	// 起点和终点相同
	if startP == endP {
		return path
	}

	current := astar.APoint{X: int(startP.X), Y: int(startP.Y)}
	end := astar.APoint{X: int(endP.X), Y: int(endP.Y)}
	path = append(path, &current)
	for current != end {
		// 计算移动方向
		dx := direction(current.X, end.X)
		dy := direction(current.Y, end.Y)

		// 移动到下一个点
		current.X += dx
		current.Y += dy
		path = append(path, &astar.APoint{X: current.X, Y: current.Y})
	}

	return path
}
