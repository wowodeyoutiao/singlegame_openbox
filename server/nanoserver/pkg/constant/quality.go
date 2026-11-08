package constant

// Quality 伙伴, 宠物颜色枚举
type Quality int32

const (
	QualityWhite  Quality = iota + 1 // 1：白色
	QualityGreen                     // 2：绿色
	QualityBlue                      // 3：蓝色
	QualityPurple                    // 4：紫色
	QualityPink                      // 5：粉色
	QualityOrange                    // 6：橙色
	QualityRed                       // 7：红色
	QualityYellow                    // 8: 黄色
)

func (c Quality) String() string {
	switch c {
	case QualityWhite:
		return "白色"
	case QualityGreen:
		return "绿色"
	case QualityBlue:
		return "蓝色"
	case QualityPurple:
		return "紫色"
	case QualityPink:
		return "粉色"
	case QualityOrange:
		return "橙色"
	case QualityRed:
		return "红色"
	case QualityYellow:
		return "黄色"
	default:
		return "未知"
	}
}
