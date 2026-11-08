package game

import (
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"sort"
	"time"
)

type CombatPlayer struct {
	cp  *RobotInfo
	win int
}

type CombatManager struct {
	red        *CombatPlayer
	blue       *CombatPlayer
	totalRound int
	curRound   int
	matchingID int64
	fighting   bool
}

var combatManagerr = NewCombatManager()

func NewCombatManager() *CombatManager {
	mgr := &CombatManager{}

	return mgr
}

func (cm *CombatManager) Init(v *protocol.CombatSetting) error {

	if cm.fighting {
		return errutil.ErrRequestFailed
	}

	if v.Red == nil || v.Blue == nil || v.TotalRound == 0 || v.Red.Job == 0 ||
		v.Blue.Job == 0 || v.Red.Level == 0 || v.Blue.Level == 0 {
		return errutil.ErrInvalidParameter
	}

	cs := protocol.CombatSetting{TotalRound: 1}
	cs.Red = v.Red
	cs.Blue = v.Blue
	cs.TotalRound = v.TotalRound

	cm.Combat(&cs)

	logger.Infof(" 模拟战斗开始 参数 :%d  red %v  blue %v", v.TotalRound, v.Red, v.Blue)
	return nil
}

func (cm *CombatManager) Init2() error {

	if cm.fighting {
		return errutil.ErrRequestFailed
	}

	cs := protocol.CombatSetting{TotalRound: 10}
	cs.Red = &protocol.PlayerSetting{
		Job:                1,
		Level:              50,
		EquipLevel:         30,
		Enhance:            50,
		Stars:              10,
		SkillLevel:         1,
		SkillAdvancedLevel: 1,
		FateLevel:          20,
		GemLevel:           0,
	}
	cs.Blue = &protocol.PlayerSetting{
		Job:                2,
		Level:              50,
		EquipLevel:         30,
		Enhance:            50,
		Stars:              10,
		SkillLevel:         1,
		SkillAdvancedLevel: 1,
		FateLevel:          20,
		GemLevel:           0,
	}

	cm.Combat(&cs)
	return nil
}

func (cm *CombatManager) Combat(s *protocol.CombatSetting) bool {

	items := []global.IdNum{}
	red := &CombatPlayer{cp: &RobotInfo{}}
	red.cp.Player = newRobotPlayer(snowflake.GenID(), "red", 1, int32(s.Red.Job), int32(s.Red.Level), 0)

	robotMgr.robots[red.cp.uid] = red.cp

	for _, v := range global.ItemCfg {

		// if  v.MainType != 6 && v.Itemlevel == int32(s.Red.GemLevel) {
		// 	items = append(items, global.IdNum{ID: v.ID, Num: 12})
		// }

		if v.MainType == 5 && v.Levellimit == s.Red.EquipLevel && v.Occupylimit == s.Red.Job {
			items = append(items, global.IdNum{ID: v.ID, Num: 1})
		}
	}

	red.cp.addSourceBatch(items, constant.ItemAddGM)

	// 装备
	for i := 0; i < len(red.cp.equips); i++ {
		red.cp.equips[i].Enhance = int32(s.Red.Enhance)
		red.cp.equips[i].Stars = int32(s.Red.Stars)

		for _, v := range red.cp.bag.ItemList {
			if equipCfg, err := global.GetItemCfg(v.CfgId); err == nil {
				if equipCfg.MainType != 5 {
					continue
				}
				if equipCfg.Equiptype == i {
					red.cp.equips[i].ItemId = v.Uid
					break
				}
			}
		}
	}

	// 宝石
	// for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
	// 	red.cp.genstoneChangeAll(i)
	// }

	for _, v := range global.SkillCfg {
		if v.UID > 50000 {
			continue
		}
		if (v.Level == int32(s.Red.SkillLevel) || v.Level == 0) && v.SkillProfession == int32(s.Red.Job) && v.RoleLevelLimit < 999 {
			red.cp.skills = append(red.cp.skills, &PlayerSkill{CfgId: v.UID, AdvancedLevel: int32(s.Red.SkillAdvancedLevel), Distance: v.MaxDistance})
		}
	}
	sort.Slice(red.cp.skills, func(i, j int) bool {
		return red.cp.skills[i].Distance > red.cp.skills[j].Distance
	})

	red.cp.fateLevel = int32(s.Red.FateLevel)

	red.cp.refreshStats()

	items = []global.IdNum{}
	blue := &CombatPlayer{cp: &RobotInfo{}}
	blue.cp.Player = newRobotPlayer(snowflake.GenID(), "blue", 2, int32(s.Blue.Job), int32(s.Blue.Level), 0)
	robotMgr.robots[blue.cp.uid] = blue.cp

	for _, v := range global.ItemCfg {

		// if v.MainType != 6 && v.Itemlevel == int32(s.Blue.GemLevel) {
		// 	items = append(items, global.IdNum{ID: v.ID, Num: 12})
		// }

		if v.MainType == 5 && v.Levellimit == s.Blue.EquipLevel && v.Occupylimit == s.Blue.Job {
			items = append(items, global.IdNum{ID: v.ID, Num: 1})
		}
	}

	blue.cp.addSourceBatch(items, constant.ItemAddGM)

	// 装备
	for i := 0; i < len(blue.cp.equips); i++ {
		blue.cp.equips[i].Enhance = int32(s.Blue.Enhance)
		blue.cp.equips[i].Stars = int32(s.Blue.Stars)

		for _, v := range blue.cp.bag.ItemList {
			if equipCfg, err := global.GetItemCfg(v.CfgId); err == nil {
				if equipCfg.MainType != 5 {
					continue
				}
				if equipCfg.Equiptype == i {
					blue.cp.equips[i].ItemId = v.Uid
					break
				}
			}
		}
	}

	// 宝石
	// for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
	// 	blue.cp.genstoneChangeAll(i)
	// }

	for _, v := range global.SkillCfg {
		if v.UID > 50000 {
			continue
		}

		if (v.Level == int32(s.Blue.SkillLevel) || v.Level == 0) && v.SkillProfession == int32(s.Blue.Job) && v.RoleLevelLimit < 999 {
			blue.cp.skills = append(blue.cp.skills, &PlayerSkill{CfgId: v.UID, AdvancedLevel: int32(s.Blue.SkillAdvancedLevel), Distance: v.MaxDistance})
		}
	}

	sort.Slice(blue.cp.skills, func(i, j int) bool {
		return blue.cp.skills[i].Distance > blue.cp.skills[j].Distance
	})

	blue.cp.fateLevel = int32(s.Blue.FateLevel)
	blue.cp.refreshStats()

	cm.red = red
	cm.blue = blue
	cm.totalRound = s.TotalRound
	cm.curRound = 1
	cm.fighting = true

	cm.FightStart()
	return true
}

func (cm *CombatManager) FightStart() {

	cm.matchingID = matchManager.combatMatch(cm.red.cp.Player, cm.blue.cp.Player)
	logger.Debugf(" ooo第【%d】轮 战斗开始 ,[血量-移速-物攻-魔攻-物防-魔防] red=[%d-%d-%d-%d-%d-%d]  blue=[%d-%d-%d-%d-%d-%d]", cm.curRound,
		cm.red.cp.currentHP, cm.red.cp.stat(pb.Attribute_MOVE_SPEED), cm.red.cp.stat(pb.Attribute_PHYSICAL_ATTACK), cm.red.cp.stat(pb.Attribute_MAGIC_ATTACK), cm.red.cp.stat(pb.Attribute_PHYSICAL_DEFENSE), cm.red.cp.stat(pb.Attribute_MAGIC_DEFENSE),
		cm.blue.cp.currentHP, cm.blue.cp.stat(pb.Attribute_MOVE_SPEED), cm.blue.cp.stat(pb.Attribute_PHYSICAL_ATTACK), cm.blue.cp.stat(pb.Attribute_MAGIC_ATTACK), cm.blue.cp.stat(pb.Attribute_PHYSICAL_DEFENSE), cm.blue.cp.stat(pb.Attribute_MAGIC_DEFENSE))
}

func (cm *CombatManager) FightEnd(matchingID int64, winId int64) {

	if cm.matchingID == matchingID {
		if winId == cm.red.cp.uid {
			cm.red.win++
		} else {
			cm.blue.win++
		}

		logger.Debugf("xxx第【%d】轮 战斗结束", cm.curRound)

		if cm.curRound < cm.totalRound {
			time.AfterFunc(time.Second, func() {
				cm.FightStart()
			})
			cm.curRound++
		} else {
			//测试结束
			cm.fighting = false
			logger.Debugf(" 测试完成 共 【%d】轮  红胜 %d  蓝胜 %d", cm.curRound, cm.red.win, cm.blue.win)
		}

	}
}

// 查询
func (cm *CombatManager) CombatResult() *protocol.CombatResult {
	if cm.red == nil || cm.blue == nil {
		return &protocol.CombatResult{}
	}
	r := &protocol.CombatResult{CurRound: cm.curRound, RedWin: cm.red.win, BlueWin: cm.blue.win, Finish: !cm.fighting}
	return r
}
