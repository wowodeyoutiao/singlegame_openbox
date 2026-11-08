package game

import "time"

type IBuff interface {
}

type BuffInfo struct {
	Id            int32      // buffID
	BuffType      int32      // 类型 1代表为增益效果 2代表为减益效果
	Mark          int32      // 1 被动，1000 技能引用的效果
	BuffStats     Stats      // 属性
	StartTime     time.Time  // 创建时间
	OverTime      time.Time  // 结束时间
	TickTime      time.Time  // 下次生效时间
	CDTime        time.Time  // 冷却时间
	TickFrequency int32      // 生效间隔
	Pretype       int32      // 大类
	Owner         IMapObject // 释放着
	Status        int32      // 状态
	Rate          int32      // 触发概率
	SkillID       int32      // 触发buff的原技能ID
	Shield        int32      // buff产生的护盾/生命上限增加
	BuffId        int64      // buff唯一Id
}
