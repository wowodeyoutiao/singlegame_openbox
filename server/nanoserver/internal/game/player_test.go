package game

import (
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/protocol/protobuf/pb"
	"testing"
)

func TestCalculateBattlePower(t *testing.T) {
	// 创建属性数据
	stats := &Stats{
		Attrs: make(map[pb.Attribute]int32),
	}

	// 设置基础属性值
	// stats.setStats(pb.Attribute_PHYSICAL_ATTACK, 6000) // 物理攻击
	stats.setStats(pb.Attribute_MAGIC_ATTACK, 2672)     // 魔法攻击 (总主攻击值=6000)
	stats.setStats(pb.Attribute_MAX_HEALTH, 31396)      // 最大生命值
	stats.setStats(pb.Attribute_PHYSICAL_DEFENSE, 1799) // 物理防御
	stats.setStats(pb.Attribute_MAGIC_DEFENSE, 1947)    // 魔法防御
	stats.setStats(pb.Attribute_HIT, 163)               // 命中值
	stats.setStats(pb.Attribute_DODGE, 1)               // 闪避值
	stats.setStats(pb.Attribute_CRITICAL_STRIKE, 45)    // 暴击值
	stats.setStats(pb.Attribute_TOUGHNESS, 1)           // 韧性值
	stats.setStats(pb.Attribute_STUN, 1)                // 眩晕效果值
	stats.setStats(pb.Attribute_STUN_RESISTANCE, 1)     // 眩晕抗性值
	stats.setStats(pb.Attribute_POISON, 1)              // 中毒效果值
	stats.setStats(pb.Attribute_POISON_RESISTANCE, 1)   // 中毒抗性值

	// 设置百分比属性
	stats.setStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE, 1200) // 暴击伤害百分比 (0.12 * 10000)

	// 设置固定增减伤属性
	stats.setStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER, 10) // 对怪增伤固定值
	stats.setStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER, 10)  // 对人增伤固定值
	stats.setStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER, 0)    // 对怪减伤固定值
	stats.setStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER, 0)     // 对人减伤固定值
	stats.setStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED, 0)      // 最终减伤固定值
	stats.setStats(pb.Attribute_DMT_FINAL_INCREASE_FIXED, 0)    // 最终增伤固定值

	// 设置百分比增减伤属性
	stats.setStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER, 0) // 对怪减伤百分比
	stats.setStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER, 0)  // 对人减伤百分比
	stats.setStats(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE, 0)   // 最终减伤百分比
	stats.setStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE, 0)         // 对怪增伤百分比
	stats.setStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE, 0)         // 对人增伤百分比
	stats.setStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, 0) // 最终增伤百分比

	// 设置吸血和反伤属性（都为0）
	stats.setStats(pb.Attribute_SUCKING_BLOOD_PERCENTAGE, 0)  // 吸血比例
	stats.setStats(pb.Attribute_REFLECT_DAMAGE_PERCENTAGE, 0) // 反伤比例

	// 创建战力配置
	powerCfg := global.PowerConfig{
		DefValue:         100.0, // 防御系数
		LifeValue:        1.0,   // 有效生命系数
		AttackValue:      1.0,   // 有效输出系数
		Hitvalue:         10.0,  // 命中系数
		Dodgevalue:       10.0,  // 闪避系数
		Critvalue:        30.0,  // 暴击系数
		Toughness:        10.0,  // 韧性系数
		Stunvalue:        10.0,
		StunResistance:   5.0,
		Poisonvalue:      10.0,
		PoisonResistance: 5.0,
	}

	// 创建技能配置映射
	// 技能id  技能等级  技能阶数
	// 2819			19			3
	// 2913			13			2
	// 3620			20			4
	skillCfg := map[int32]*global.SkillConfig{
		2819: {Powerfactor: 95},
		2913: {Powerfactor: 66},
		3620: {Powerfactor: 108},
	}

	// 创建技能列表
	skills := []*PlayerSkill{
		{CfgId: 2819, AdvancedLevel: 3},
		{CfgId: 2913, AdvancedLevel: 2},
		{CfgId: 3620, AdvancedLevel: 4},
	}

	// 调用计算函数
	result := CalculateBattlePower(stats, powerCfg, skillCfg, skills, constant.Wizard)

	// 验证结果
	expected := int64(92326)
	if result != expected {
		t.Errorf("CalculateBattlePower() = %d, expected %d", result, expected)
	}

	t.Logf("计算战力: %d, 预期战力: %d", result, expected)
}
