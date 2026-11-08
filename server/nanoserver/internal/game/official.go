package game

import (
	"nanoserver/global"
	"nanoserver/pkg/constant"
)

// CheckOfficialSuppression 判断玩家 A 对玩家 B 的官职压制关系
func CheckOfficialSuppression(aExp, bExp int64, cfgs []*global.OfficialConfig) constant.OfficialSuppressionState {
	getOfficialLevel := func(exp int64, cfgs []*global.OfficialConfig) int32 {
		if len(cfgs) == 0 || exp < cfgs[0].EXP {
			return 0 // 无官职
		}
		low, high := 0, len(cfgs)-1
		level := int32(0)
		for low <= high {
			mid := low + (high-low)/2
			if exp >= cfgs[mid].EXP {
				level = cfgs[mid].UID
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
		return level
	}

	aLevel := getOfficialLevel(aExp, cfgs)
	bLevel := getOfficialLevel(bExp, cfgs)

	switch {
	case aLevel > bLevel:
		return constant.SuppressionApplied // A 压制 B
	case aLevel < bLevel:
		return constant.SuppressionReversed // A 被 B 压制
	default:
		return constant.SuppressionNone // 官职相同，无压制
	}
}
