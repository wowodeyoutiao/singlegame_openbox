package game

import (
	"math"
	"math/rand"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/maputil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"time"
)

// 迅影追击 SkillTypeSwiftStrike
func (o *MapObject) swiftStrikeEnd(skillId int32, direction pb.Direction) int32 {

	var maxDistance int32 = 0
	if o.dead {
		return maxDistance
	}
	skill := o.getSkill(skillId)
	if skill == nil {
		return maxDistance
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return maxDistance
	}

	var endCell *Cell
	for i := 0; i < int(skillCfg.Radius); i++ {
		point := maputil.Move(&o.cell.location, direction, int(i))
		if cell := o.desk.mapInfo.getCell(point); cell != nil {
			if !cell.block {
				maxDistance = int32(i + 1)
				endCell = cell
				continue
			}
		}
		break
	}
	if endCell == nil {
		return maxDistance
	}

	if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
		o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
	}
	return maxDistance
}

// 拒敌千里 SkillTypeRepelFoe 受击者后退
func (o *MapObject) repelFoeEnd(skillId int32, attacker IMapObject) {
	if o.dead {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}
	var backDir pb.Direction
	var endCell *Cell
	backDir = maputil.DirectionFromPointI(&attacker.GetCell().location, &o.cell.location)
	for i := 1; i <= int(skillCfg.SpecialEffect.BackDistence); i++ {
		point := maputil.Move(&o.cell.location, backDir, int(i))
		if cell := o.desk.mapInfo.getCell(point); cell != nil {
			if !cell.block {
				endCell = cell
				continue
			}
		}
		break
	}
	if endCell == nil {
		return
	}

	//p.logger.Debugf("拒敌千里 起始坐标 %v 终点坐标 %v  距离 %d", p.cell.location, endCell.location, maputil.Distance(p.location, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}))
	if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
		o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
	}

	// 采集中断
	if o.GetObAction() == ObjectActionMining {
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			p.miningBreak()
		}
	}
}

// 替身宝宝 SkillTypeDecoy
func (o *MapObject) decoyEnd(skillId int32, direction pb.Direction) {
	if o.dead {
		return
	}
	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}
	location := o.cell.location
	// 人先飞走
	var endCell *Cell
	for i := 0; i < int(skillCfg.SpecialEffect.FlyDistence); i++ {
		point := maputil.Move(&o.cell.location, direction, int(i))
		if cell := o.desk.mapInfo.getCell(point); cell != nil {
			if !cell.block {
				endCell = cell
				continue
			}
		}
		break
	}

	if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
		o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
	}

	// 生成替身
	for _, tri := range skillCfg.TriggerBuff {
		bufCfg, err := global.GetBuffCfg(tri)
		if err != nil {
			continue
		}

		// 筛除特殊buff
		if bufCfg.Pretype == constant.BuffDecoy {
			stats := &Stats{Attrs: map[pb.Attribute]int32{}}
			for _, v := range skillCfg.Skillbuff {
				stats.setStats(pb.Attribute(v.Key), v.Value)
			}

			// 技能进阶效果
			if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skill.AdvancedLevel); err == nil {
				if skillUpperlevel.DamagePer > 0 {
					stats.setStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE, skillUpperlevel.DamagePer)
				}
			}

			// 替身的最大生命值
			decoyHp := float64(o.GetMaxHP()) * float64(stats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
			stats.setStats(pb.Attribute_HEALTH, int32(decoyHp))
			stats.changeStats(pb.Attribute_MAX_HEALTH, int32(decoyHp))
			mc := float64(o.GetMC()) * float64(stats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE)) / BaseRate
			mc = math.Max(mc, 1)
			stats.setStats(pb.Attribute_MAGIC_ATTACK, int32(mc))

			if buffOb := CreateBuffObject(bufCfg.ID, ObjectTypeDecoy, o, skillCfg.Bufftime, skillId, stats); buffOb != nil {
				o.desk.addObject(buffOb, location)
			}
		}
	}
}

// 瞬移 SkillTypeInstantaneousMove
func (o *MapObject) instantaneousMoveEnd(skillId int32) {
	if o.dead {
		return
	}
	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	var endCell *Cell
	if global.RandIn(int(BaseRate)) <= int(skillCfg.Skilltrigger) {
		endCell = o.desk.mapInfo.getAvailableCell(&o.cell.location, skillCfg.Radius)
	}
	if endCell == nil {
		endCell = o.cell
	}

	if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
		o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
	}
}

// 剑开天门（特技）SkillTypeSwordHeavenBreaker
func (o *MapObject) swordHeavenBreakerEnd(skillId int32) bool {
	if o.dead {
		return false
	}
	if skill := o.getSkill(skillId); skill == nil {
		return false
	}
	return true
}

// 魔法酸雨 SkillTypeMagicAcidRain
func (o *MapObject) magicAcidRain(damage float64, skillID int32) {
	if o.dead {
		return
	}

	skill := o.getSkill(skillID)
	if skill == nil {
		return
	}

	skillCfg, err := global.GetSkillCfg(int32(skillID))
	if err != nil {
		return
	}
	attackCells := o.desk.mapInfo.getSquareCells(&o.cell.location, skillCfg.Radius)
	var targets []IMapObject
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.CanAttackTarget(obj) && o.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == o.uid {
			continue
		}

		if t.ObjectType() < ObjectTypePlayer && t.ObjectType() > ObjectTypeDecoy {
			continue
		}

		if t.ObjectType() == ObjectTypePlayer {
			if o.desk.roomCfg.Roomtype <= constant.PVEPeacefulMap || o.desk.roomCfg.Roomtype >= constant.SoloTowerMap {
				continue
			}
		}

		// 和平模式不能攻击除怪以外目标
		if o.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
			if t.ObjectType() != ObjectTypeMonster {
				continue
			}
		}

		// 安全区内不能攻击玩家
		if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (o.desk.mapInfo.safeArea(o.location) || o.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		// 命中-闪避
		hit := o.stats.getStats(pb.Attribute_HIT)
		if tstats := t.GetStats(); tstats != nil {
			dodge := t.GetStats().getStats(pb.Attribute_DODGE)
			if hit < dodge && global.RandIn(int(dodge)) > int(hit) {
				t.Dodged()
				continue
			}
		}

		if t.ObjectType() == ObjectTypeMonster {
			damage += float64(o.stat(pb.Attribute_DMT_INCREASE_FIXED_MONSTER))
			damage += damage * float64(o.stat(pb.Attribute_PVE_DAMAGE_PERCENTAGE)) / BaseRate
		} else if t.ObjectType() == ObjectTypePlayer {
			damage += float64(o.stat(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))
			damage += damage * float64(o.stat(pb.Attribute_PVP_DAMAGE_PERCENTAGE)) / BaseRate
		}

		damage += float64(o.stat(pb.Attribute_DMT_FINAL_INCREASE_FIXED))
		damage += damage * float64(o.stat(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE)) / BaseRate

		//攻击-防御
		if tstats := t.GetStats(); tstats != nil {
			damage -= float64(tstats.getStats(maputil.DefenseAttributes(o.AttackAttributes())))
		}

		if damage <= 0 {
			damage = float64(rand.Intn(2)) + 1
		}

		//p.logger.Debugf("P 魔法酸雨 targetID: %d  目标位置: %v skillID: %d  Power:%f", t.Uid(), t.Location(), skillID, damage)
		t.Attacked(o, skillID, skill.AdvancedLevel, int32(damage), true)
	}
}

// 回光返照
func (o *MapObject) lastGasp(target IMapObject, skillID int32) bool {
	if o.dead {
		return false
	}
	if !o.CanHelpTarget(target) {
		return false
	}

	var skill = o.getSkill(skillID)
	if skill == nil {
		return false
	}
	if time.Now().Before(skill.Cooldown) {
		return false
	}
	skillCfg, err := global.GetSkillCfg(int32(skillID))
	if err != nil {
		return false
	}

	if skillCfg.CoolDown > 0 {
		skill.Cooldown = time.Now().Add(time.Duration(float64(skillCfg.CoolDown)*0.99) * time.Millisecond)
		//o.pushMsg(protocol.RouteSkillCooldownPush, &pb.SkillCooldownPush{CfgId: skillID, Delay: skillCfg.CoolDown})
	}

	//延迟生效
	pdesk := o.desk
	pt := target
	ilocation := &o.cell.location
	time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
		if pdesk != nil {
			pdesk.addDelayAction(o.objectType, o.uid, skillID, pt, o.direction, ilocation)
		}
	})
	return true
}

func (o *MapObject) lastGaspEnd(target IMapObject, skillID int32) {
	if !target.Dead() {
		return
	}

	var skillAdvanceLevel int32
	if skill := o.getSkill(skillID); skill != nil {
		skillAdvanceLevel = skill.AdvancedLevel
		//技能CD
		if skillCfg, err := global.GetSkillCfg(skillID); err == nil {
			if skillCfg.CoolDown > 0 {
				skill.Cooldown = time.Now().Add(time.Duration(float64(skillCfg.CoolDown)*0.99) * time.Millisecond)
				o.pushMsg(protocol.RouteSkillCooldownPush, &pb.SkillCooldownPush{CfgId: skillID, Delay: skillCfg.CoolDown})
			}
			var hp float64
			for _, v := range skillCfg.Skillbuff {
				if v.Key == int32(pb.Attribute_HEALTH_PERCENTAGE) {
					hp += float64(target.GetMaxHP()) * float64(v.Value) / BaseRate
				}
			}

			target.Revive(skill, skillAdvanceLevel, int32(hp))

		}
		logger.Debugf(" %d:%s 伙伴复活 目标[%s] 复活后血量 %d 总血量 %d", o.uid, o.name, target.Name(), target.GetMaxHP(), target.GetMaxHP())
	}
}

// 群体净化
func (o *MapObject) MassPurificationEnd(skillID int32) {
	if o.dead {
		return
	}
	if skill := o.getSkill(skillID); skill != nil {

		for _, friend := range o.Friends() {
			friend.CleanDebuff()
		}
		o.SkillBuffTrigger(skill)
	}
}

// 幻象爆炸
func (o *MapObject) DecoyExplosionEnd() {
	var targets []IMapObject
	var targetsUids []int64
	Radius := 2
	if skillCfg, err := global.GetSkillCfg(o.linkedSkill); err == nil {
		Radius = int(skillCfg.Radius)
	}
	attackCells := o.desk.mapInfo.getSquareCells(&o.cell.location, int32(Radius))
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.CanAttackTarget(obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == o.uid {
			continue
		}

		if t.ObjectType() < ObjectTypePlayer && t.ObjectType() > ObjectTypeMonster {
			continue
		}

		// 非PVP地图不能攻击玩家
		if t.ObjectType() == ObjectTypePlayer || t.ObjectType() == ObjectTypeDecoy {
			if o.desk.roomCfg.Roomtype <= constant.PVEPeacefulMap || o.desk.roomCfg.Roomtype >= constant.SoloTowerMap {
				logger.Debugf("目标无法攻击,房间类型%+v userId: %d , targetId:%d, targetType:%+v", o.desk.roomCfg.Roomtype, o.uid, t.Uid(), t.ObjectType())
				continue
			}
		}

		// 和平模式不能攻击除怪以外目标
		if o.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
			if t.ObjectType() != ObjectTypeMonster {
				continue
			}
		}

		// 安全区内不能攻击玩家
		if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if (t.ObjectType() != ObjectTypeMonster || t.ObjectType() == ObjectTypeDecoy) && (o.desk.mapInfo.safeArea(o.location) || o.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		targetsUids = append(targetsUids, t.Uid())

		packet := &pb.AttackPush{
			AttackerId: o.uid,
			TargetUids: targetsUids,
			SkillId:    o.linkedSkill,
			Location:   o.location,
			Direction:  o.direction}
		o.desk.Broadcast(protocol.RouteAttackPush, packet)

		//伤害延迟生效
		odesk := o.desk
		ot := t
		ilocation := &o.cell.location
		time.AfterFunc(time.Duration(500)*time.Millisecond, func() {
			if odesk != nil {
				odesk.addDelayAction(ObjectTypeDecoy, o.uid, o.linkedSkill, ot, o.direction, ilocation)
			}
		})
	}

	o.delTime = time.Now().Add(600 * time.Millisecond) // 替身延迟消失
}

// 斩杀
func (o *MapObject) ExecuteEnd(target IMapObject, skillID int32) {
	if target.Dead() {
		return
	}
	target.Attacked(o, skillID, 0, 0, false)
}

// 章鱼技能1
func (o *MapObject) octopusSkill1End(target IMapObject, skillID int32) {
	if target.Dead() {
		return
	}
	if o.desk == nil {
		return
	}
	var Radius int32
	skillCfg, err := global.GetSkillCfg(skillID)
	if err != nil {
		return
	}
	Radius = skillCfg.Radius
	callNum := 10
	callId := skillCfg.SpecialEffect.Callout

	monCfg := global.GetMonsterCfg(callId)
	if monCfg == nil {
		return
	}

	for i := 0; i < callNum; i++ {
		if cell := o.desk.mapInfo.getAvailableCell(&o.cell.location, Radius); cell != nil {
			o.desk.addMonsterRelive(callId, cell.location, 26000, skillID, 0, 0, "")
		}
	}
	logger.Debugf("章鱼技能1 [%d:%s] 技能ID %d  召唤怪物ID %d  数量 %d  ", o.uid, o.name, skillID, callId, callNum)
}

// 章鱼爆炸
func (o *MapObject) octopusSkill3End(skillId int32) {
	var targets []IMapObject
	Radius := 2

	if skillCfg, err := global.GetSkillCfg(skillId); err == nil {
		Radius = int(skillCfg.Radius)

	}
	attackCells := o.desk.mapInfo.getSquareCells(&o.cell.location, int32(Radius))
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.CanAttackTarget(obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == o.uid {
			continue
		}

		if t.ObjectType() < ObjectTypePlayer && t.ObjectType() > ObjectTypeMonster {
			continue
		}

		// 非PVP地图不能攻击玩家
		if t.ObjectType() == ObjectTypePlayer || t.ObjectType() == ObjectTypeDecoy {
			if o.desk.roomCfg.Roomtype <= constant.PVEPeacefulMap || o.desk.roomCfg.Roomtype >= constant.SoloTowerMap {
				logger.Debugf("目标无法攻击,房间类型%+v userId: %d , targetId:%d, targetType:%+v", o.desk.roomCfg.Roomtype, o.uid, t.Uid(), t.ObjectType())
				continue
			}
		}

		// 和平模式不能攻击除怪以外目标
		if o.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
			if t.ObjectType() != ObjectTypeMonster {
				continue
			}
		}

		// 安全区内不能攻击玩家
		if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if (t.ObjectType() != ObjectTypeMonster || t.ObjectType() == ObjectTypeDecoy) && (o.desk.mapInfo.safeArea(o.location) || o.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		o.delayAttack(t, skillId, 1)
	}

	o.linkedSkill = 0
	logger.Debugf("章鱼爆炸 [%d:%s] 技能ID %d  位置 %v 半径 %d 影响目标数量 %d", o.uid, o.name, skillId, o.cell.location, Radius, len(targets))
	o.die(nil)

}

// 穷奇1技能 龙卷风
func (o *MapObject) qiongqiSkill1End(skillId int32) {
	if o.dead {
		return
	}
	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}
	location := o.cell.location

	// 生成 5 个龙卷风
	for _, tri := range skillCfg.TriggerBuff {
		bufCfg, err := global.GetBuffCfg(tri)
		if err != nil {
			continue
		}

		if bufCfg.Pretype == constant.BuffTornado {
			for i := 0; i < 5; i++ {
				stats := &Stats{Attrs: map[pb.Attribute]int32{}}
				for _, v := range skillCfg.Skillbuff {
					stats.setStats(pb.Attribute(v.Key), v.Value)
				}

				// 技能进阶效果
				if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skill.AdvancedLevel); err == nil {
					if skillUpperlevel.DamagePer > 0 {
						stats.setStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE, skillUpperlevel.DamagePer)
					}
				}
				if cell := o.desk.mapInfo.getAvailableCell(&location, skillCfg.Radius); cell != nil {
					if buffOb := CreateBuffObject(bufCfg.ID, ObjectTypeBuff, o, skillCfg.Bufftime, skillId, stats); buffOb != nil {
						o.desk.addObject(buffOb, cell.location)
					}
				}
			}
		}
	}
}

// 龙卷风攻击
func (o *MapObject) tornadoAttackEnd() {

	skillCfg, err := global.GetSkillCfg(o.linkedSkill)
	if err != nil {
		return
	}

	attackCells := []*Cell{}
	targets := []IMapObject{}
	attackCells = o.desk.mapInfo.getSquareCells(&o.cell.location, 5)

	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.CanAttackTarget(obj) && o.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == o.uid {
			continue
		}

		// 安全区内不能攻击玩家
		if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (o.desk.mapInfo.safeArea(o.location) || o.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}
		damage := float64(t.GetMaxHP()) * float64(o.stat(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
		if damage <= 0 {
			continue
		}
		t.Attacked(o, o.linkedSkill, 0, int32(damage), false)

		logger.Debugf("龙卷风攻击 [%d:%s] 技能ID %d  位置 %v 半径 %d  目标 [%d:%s] 位置 %v 伤害 %f", o.uid, o.name, o.linkedSkill, o.cell.location, 5, t.Uid(), t.Name(), t.Location(), damage)
	}

}

// 和尚1技能 凝视，扇形120度范围内，光柱攻击宽度3-5米(具体看效果)，被光柱照射，每秒造成200%攻击+8000点最终伤害，CD25秒，持续5秒，在扇形范围内来回照射
// 根据特效播放时间计算光柱移动速度，特效完整时间2.6秒
func (o *MapObject) heshangSkill1Begin(skillId int32) {

	if !o.CanAttack() {
		return
	}
	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	location := o.cell.location
	direction := o.direction
	odesk := o.desk

	o.attackTime = time.Now().Add(1500 * time.Millisecond)
	o.moveTime = time.Now().Add(1500 * time.Millisecond)

	for i := 1; i <= 3; i++ {

		switch i {
		case 1:
			direction = maputil.ShiftDirection(o.direction, -1)
		case 2:
			direction = o.direction
		case 3:
			direction = maputil.ShiftDirection(o.direction, 1)
		}

		idir := direction
		ilocation := &location
		time.AfterFunc(time.Duration(skillCfg.DelayTime*int32(i))*time.Millisecond, func() {
			if odesk != nil {
				odesk.addDelayAction(o.objectType, o.uid, skillId, nil, idir, ilocation)
			}
		})
		logger.Debugf("%d:%s Begin和尚凝视技能第[%d]次攻击 , 位置 %v 方向 %d", o.uid, o.name, i, ilocation, idir)
	}
}

func (o *MapObject) heshangSkill1End(skillId int32, direction pb.Direction, location global.PointI) {

	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	var d1, d2 pb.Direction
	d1 = direction
	d2 = maputil.ShiftDirection(d1, 1)
	attackCells := o.desk.mapInfo.getAngleCells(&location, d1, d2, skillCfg.Radius)
	logger.Debugf("%d:%s End和尚凝视技能 , 位置 %v 左方向 %d 右方向 %d 半径 %d 攻击坐标数量 %d", o.uid, o.name, location, d1, d2, skillCfg.Radius, len(attackCells))
	o.attackLocation(attackCells, skillCfg)
}

// 被动技能效果
func (o *MapObject) passiveSkillTriger(attacker IMapObject, skillId int32, skillAdvanceLevel int32, power float64) {
	skillCfg, err := global.GetSkillCfg(skillId)
	if err != nil {
		return
	}
	switch skillCfg.SkillID {
	case constant.QiongqiSkill2:
		rate := float64(o.GetMaxHP()-o.currentHP) / float64(o.GetMaxHP()) * 100
		stats := &Stats{Attrs: make(map[pb.Attribute]int32)}
		for _, v := range skillCfg.Skillbuff {
			stats.setStats(pb.Attribute(v.Key), v.Value*int32(rate))
		}
		o.AddBuff(101, o, skillCfg.Bufftime, stats, skillCfg.Effecttrigger, skillId)

	case constant.HeshangSkill2:
		for _, b := range skillCfg.TriggerBuff {
			stats := &Stats{Attrs: make(map[pb.Attribute]int32)}
			for _, v := range skillCfg.Skillbuff {
				stats.setStats(pb.Attribute(v.Key), v.Value)
			}
			attacker.AddBuff(b, o, skillCfg.Bufftime, stats, skillCfg.Effecttrigger, skillId)
		}

	}
}

// SkillTypeStarShatter 受击者后退
func (o *MapObject) starShatterEnd(skillId int32, attacker IMapObject) {
	if o.dead {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}
	var backDir pb.Direction
	var endCell *Cell
	backDir = maputil.DirectionFromPointI(&attacker.GetCell().location, &o.cell.location)
	for i := 1; i <= int(skillCfg.SpecialEffect.BackDistence); i++ {
		point := maputil.Move(&o.cell.location, backDir, int(i))
		if cell := o.desk.mapInfo.getCell(point); cell != nil {
			if !cell.block {
				endCell = cell
				continue
			}
		}
		break
	}
	if endCell == nil {
		return
	}

	logger.Debugf("碎星击 起始坐标 %v 终点坐标 %v  距离 %d", o.cell.location, endCell.location, maputil.Distance(o.location, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}))
	if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
		o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
	}

	// 采集中断
	if o.GetObAction() == ObjectActionMining {
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			p.miningBreak()
		}
	}
}

// 荆棘牢笼
func (o *MapObject) brambleCageEnd(skillId int32, location global.PointI) {

	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	attackCells := o.desk.mapInfo.getSquareCells(&location, skillCfg.Radius)
	o.attackLocation(attackCells, skillCfg)
}

// 女神降临
func (o *MapObject) goddessDescentEnd(skillId int32) {

	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}
	target_friends := []IMapObject{}
	attackCells := o.desk.mapInfo.getSquareCells(&o.cell.location, skillCfg.Radius)
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range target_friends {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.CanHelpTarget(obj) {
					target_friends = append(target_friends, obj)
					continue
				}
				if o.GuildId() != 0 && o.GuildId() == obj.GuildId() {
					target_friends = append(target_friends, obj)
				}
			}
		}
	}

	for _, t := range target_friends {

		if t.ObjectType() < ObjectTypePlayer || t.ObjectType() > ObjectTypePartner {
			continue
		}

		// 先复活再加buff
		if t.Dead() {
			t.Revive(skill, 0, 1)
		}

		for _, v := range skillCfg.TriggerBuff {
			cfg, err := global.GetBuffCfg(v)
			if err != nil {
				continue
			}
			bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
			for _, s := range skillCfg.Skillbuff {
				if constant.BuffFixedHealOnece == constant.BuffType(cfg.Pretype) &&
					(pb.Attribute(s.Key) == pb.Attribute_HEALTH_PERCENTAGE || pb.Attribute(s.Key) == pb.Attribute_HEALTH) {
					bufStats.setStats(pb.Attribute(s.Key), s.Value)
				}
				if constant.BuffPercentageDefenseBoost == constant.BuffType(cfg.Pretype) &&
					(pb.Attribute(s.Key) == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || pb.Attribute(s.Key) == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE) {
					bufStats.setStats(pb.Attribute(s.Key), s.Value)
				}
			}
			t.AddBuff(v, o, skillCfg.Bufftime, bufStats, skillCfg.Effecttrigger, skillId)
		}
	}

}

// SwordRain
func (o *MapObject) swordRainEnd(skillId int32, location global.PointI) {
	skill := o.getSkill(skillId)
	if skill == nil {
		return
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	attackCells := o.desk.mapInfo.getSquareCells(&location, skillCfg.Radius)

	var targets []IMapObject
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
					if (obj.ObjectType() != ObjectTypeMonster || obj.ObjectType() == ObjectTypeDecoy) && (o.desk.mapInfo.safeArea(o.location) || o.desk.mapInfo.safeArea(obj.Location())) {
						continue
					}
				}
				if o.CanAttackTarget(obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		o.delayAttack(t, skillId, 1)
	}

	// 剑雨有概率额外攻击一次
	if global.RandIn(int(BaseRate)) < int(skillCfg.Skilltrigger) {
		logger.Debugf("%d:%s 剑雨触发额外攻击 ", o.uid, o.name)
		for _, t := range targets {
			o.delayAttack(t, skillId, 2)
		}
	}
}
