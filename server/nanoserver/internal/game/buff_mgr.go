package game

import (
	"math"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/utils"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"sync"
	"time"
)

type IBuffMgr interface {
}

type BuffManager struct {
	object   *MapObject
	buffs    *utils.SafeMap[int32, *BuffInfo]
	sPool    sync.Pool
	tickTime time.Time
}

func (m *BuffManager) init(object *MapObject) {
	m.object = object
	m.buffs = utils.NewSafeMap[int32, *BuffInfo]()
	m.sPool.New = func() interface{} {
		return make([]*BuffInfo, 0)
	}
	m.tick()
}

func (m *BuffManager) IterateBuff(f func(buff *BuffInfo) bool) {
	if m.buffs.Len() == 0 {
		return
	}
	tempSlice := m.sPool.Get().([]*BuffInfo)
	tempSlice = tempSlice[:0]
	m.buffs.Range(func(key int32, val *BuffInfo) bool {
		tempSlice = append(tempSlice, val)
		return true
	})
	for _, b := range tempSlice {
		if !f(b) {
			break
		}
	}
	m.sPool.Put(tempSlice)
}

func (m *BuffManager) rangeBuff(f func(buff *BuffInfo) bool) {
	if m.buffs.Len() == 0 {
		return
	}
	m.buffs.Range(func(key int32, val *BuffInfo) bool {
		if !f(val) {
			return false
		}
		return true
	})
}

func (b *BuffManager) ProcessBuff() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	p := b.object
	//o.block.Lock()
	//defer o.block.Unlock()
	if p.dead {
		return
	}
	//大于50毫秒处理一次
	if !b.tick() {
		return
	}
	now := time.Now()
	var delBufIds []int32
	var bRefreshStats bool
	b.IterateBuff(func(buff *BuffInfo) bool {
		cfgId := buff.Id
		// CD 和 持续时间都到， 删除buff
		if now.After(buff.OverTime) && now.After(buff.CDTime) && buff.Mark != constant.BuffMarkPassive {
			delBufIds = append(delBufIds, cfgId)
			return true
		}
		// 持续时间到 , 清除效果保留至CD结束
		if now.After(buff.OverTime) && buff.Mark != constant.BuffMarkPassive {
			if buff.Status == int32(constant.BuffStatusStart) {
				buff.BuffStats.clear()
				bRefreshStats = true
				buff.Status = int32(constant.BuffStatusStop)
			}
			return true
		}
		// 嘲讽状态 释放者退出 死亡 超过6米 均接触嘲讽状态
		if buff.BuffType == constant.BuffTaunt {
			if buff.Owner == nil || buff.Owner.Dead() || maputil.Distance(p.location, buff.Owner.Location()) > 6 {
				delBufIds = append(delBufIds, cfgId)
				return true
			}
		}
		if buff.TickFrequency == 0 {
			return true
		}
		if buff.Rate != 0 && global.RandIn(int(BaseRate)) >= int(buff.Rate) {
			return true
		}
		// 间隔未到
		if !now.After(buff.TickTime) {
			return true
		}

		var damage, recoverHp float64
		switch constant.BuffType(buff.Pretype) {
		case constant.BuffFixedHeal, constant.BuffPercentageHeal:
			buff.BuffStats.Range(func(i pb.Attribute, e int32) bool {
				if i == pb.Attribute_HEALTH {
					recoverHp += float64(e)
				}
				if i == pb.Attribute_HEALTH_PERCENTAGE {
					recoverHp += float64(p.GetMaxHP()) * float64(e) / BaseRate
				}
				return true
			})
			//for i, e := range buff.BuffStats.Attrs {
			//	if i == pb.Attribute_HEALTH {
			//		recover += float64(e)
			//	}
			//	if i == pb.Attribute_HEALTH_PERCENTAGE {
			//		recover += float64(p.GetMaxHP()) * float64(e) / BaseRate
			//	}
			//}
		case constant.BuffFixedDamage, constant.BuffPercentageDamage, constant.BuffFixedDamage2:
			buff.BuffStats.Range(func(i pb.Attribute, e int32) bool {
				if i == pb.Attribute_HEALTH {
					damage += float64(e)
				}
				if i == pb.Attribute_HEALTH_PERCENTAGE {
					damage += float64(p.GetMaxHP()) * float64(e) / BaseRate
				}

				//出血，中毒，灭世，三个技能，输出类型的最高伤害不超过当前攻击力3倍；掉血类型的最高掉血量不超过当前释放者血量上限的3倍
				if damage > float64(buff.Owner.GetMaxHP()*3) {
					if buff.Owner != nil {
						damage = float64(buff.Owner.GetMaxHP() * 3)
					} else {
						damage = 1
					}
				}
				return true
			})
			//for i, e := range buff.BuffStats.Attrs {
			//	if i == pb.Attribute_HEALTH {
			//		damage += float64(e)
			//	}
			//	if i == pb.Attribute_HEALTH_PERCENTAGE {
			//		damage += float64(p.GetMaxHP()) * float64(e) / BaseRate
			//	}
			//
			//	//出血，中毒，灭世，三个技能，输出类型的最高伤害不超过当前攻击力3倍；掉血类型的最高掉血量不超过当前释放者血量上限的3倍
			//	if damage > float64(buff.Owner.GetMaxHP()*3) {
			//		if buff.Owner != nil {
			//			damage = float64(buff.Owner.GetMaxHP() * 3)
			//		} else {
			//			damage = 1
			//		}
			//	}
			//}
		case constant.BuffFixedHealOnece:
			//for i, e := range buff.BuffStats.Attrs {
			//	if i == pb.Attribute_HEALTH {
			//		recover += float64(e)
			//	}
			//	if i == pb.Attribute_HEALTH_PERCENTAGE {
			//		recover += float64(p.GetMaxHP()) * float64(e) / BaseRate
			//	}
			//}
			buff.BuffStats.Range(func(i pb.Attribute, e int32) bool {
				if i == pb.Attribute_HEALTH {
					recoverHp += float64(e)
				}
				if i == pb.Attribute_HEALTH_PERCENTAGE {
					recoverHp += float64(p.GetMaxHP()) * float64(e) / BaseRate
				}
				return true
			})
			delBufIds = append(delBufIds, cfgId)
		case constant.BuffAcidRain:
			damageRain := float64(p.GetMC()) * float64(buff.BuffStats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE)) / BaseRate
			damageRain += damageRain * (1 + float64(buff.BuffStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate)
			p.magicAcidRain(damageRain, buff.SkillID)

		case constant.BuffCurHpPercentageDamage:
			damage += float64(p.currentHP) * float64(buff.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
			if damage > float64(buff.Owner.GetMaxHP()*3) {
				if buff.Owner != nil {
					damage = float64(buff.Owner.GetMaxHP() * 3)
				} else {
					damage = 1
				}
			}
		}
		// 挫折意志 减伤效果只针对sype = 4 或 5
		var damage_percentage_decrease, flat_damage_decrease int32
		if buff.SkillID > 0 {
			if skillCfg, err := global.GetSkillCfg(buff.SkillID); err == nil {
				if skillCfg.SType == constant.SkillTypeCareerSpecial {
					damage_percentage_decrease += p.stat(pb.Attribute_CONTROL_DAMAGE_REDUCTION_PERCENTAGE)
					flat_damage_decrease += p.stat(pb.Attribute_CONTROL_DAMAGE_REDUCTION_FIXED_VALUE)
				}
			}
		}
		damage -= damage * float64(damage_percentage_decrease) / BaseRate
		damage -= float64(flat_damage_decrease)

		if damage > 0 && p.currentHP > 1 {
			damage = math.Min(damage, float64(p.currentHP-1))
			if p.ChangeHp(-int32(damage), buff.Owner) == 0 {
				return false
			}
		}
		if recoverHp > 0 {
			p.attackId = buff.Owner.Uid()
			if p.ChangeHp(int32(recoverHp), buff.Owner) == 0 {
				return false
			}
		}
		buff.TickTime = time.Now().Add(time.Duration(buff.TickFrequency) * time.Millisecond)
		//logger.Debugf("  Buff生效 buffid = %d , uid:name %d:%s 伤害 %f  治疗 %f  到期时间 %v ", v.Id, o.uid, o.name, damage, recover, v.OverTime)
		return true
	})

	if !p.dead {
		if len(delBufIds) > 0 {
			for _, v := range delBufIds {
				b.delBuff(v)
			}
			bRefreshStats = true
		}

		if bRefreshStats {
			p.RefreshStatsAsync()
		}
	}

}

func (b *BuffManager) CanAddBuff(buffId int32) bool {
	o := b.object
	cfg, err := global.GetBuffCfg(buffId)

	if err != nil {
		logger.Warningf("获取buff配置失败 %d:%s buffid=%d", o.uid, o.name, buffId)
		return false
	}

	if o.dead {
		return false
	}

	if cfg.Mark == constant.BuffTypeChainSkill && o.desk == nil {
		return false
	}

	if cfg.Type == constant.BuffTypeDebuff {
		if o.HaveBuff(constant.BuffGodMode) {
			return false
		}

		if o.HaveBuff(constant.BuffImmune) {
			return false
		}
		// 免疫debuff
		if immuneBuff := o.getBuffByType(constant.BuffNoNegativeEffect); immuneBuff != nil {
			if global.RandIn(int(BaseRate)) < int(immuneBuff.Rate) {
				logger.Debugf("免疫debuff  %d:%s  buffid=%d", o.uid, o.name, buffId)
				return false
			}
		}
	}

	return true
}

func (b *BuffManager) addBuff(buffIID int64, cfgId int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	//o.block.Lock()
	//defer o.block.Unlock()
	cfg, err := global.GetBuffCfg(cfgId)
	if err != nil {
		return
	}
	now := time.Now()
	buff, ok := b.buffs.Get(cfgId)
	if ok {
		if now.Before(buff.CDTime) {
			return
		}
	}

	//if buff, ok := o.buffs[cfgId]; ok {
	//	if now.Before(buff.CDTime) {
	//		return
	//	}
	//}

	// 同类型buff 效果强的替换掉效果弱的,弱的不生效
	bRecover := false
	recoverId := 0
	var restoreHP float64

	b.IterateBuff(func(buff *BuffInfo) bool {
		if buff.Pretype != cfg.Pretype {
			return true
		}
		for oldK, oldV := range buff.BuffStats.copyStats() {
			for _, cfgV := range cfg.Value {
				if cfgV.Key == int32(oldK) {
					if oldV < cfgV.Value {
						bRecover = true
						recoverId = int(buff.Id)
						break
					} else {
						recoverId = -1
						return false
					}
				}
			}
			if exStats != nil {
				//exStats.Attrs 这里属性都是配置
				for exK, exV := range exStats.Attrs {
					if exK == oldK {
						if oldV < exV {
							bRecover = true
							recoverId = int(buff.Id)
							break
						} else {
							recoverId = -1
							return false
						}
					}
				}
			}
		}
		return false
	})

	if recoverId == -1 {
		return
	}

	p := b.object
	if bRecover && recoverId != 0 {
		b.delBuff(int32(recoverId))
	}

	// 免疫时间 >= buff时间
	if cfg.LastTime+addSec <= 0 {
		return
	}

	buff = &BuffInfo{
		Id:            cfgId,
		BuffType:      cfg.Type,
		Mark:          cfg.Mark,
		StartTime:     now,
		TickTime:      now.Add(time.Duration(cfg.EffectTime) * time.Millisecond),
		OverTime:      now.Add(time.Duration(cfg.LastTime+addSec+100) * time.Millisecond),
		CDTime:        now.Add(time.Duration(cfg.CDTime) * time.Millisecond),
		Pretype:       cfg.Pretype,
		Owner:         owner,
		Status:        int32(constant.BuffStatusStart),
		BuffStats:     Stats{Attrs: make(map[pb.Attribute]int32)},
		TickFrequency: cfg.EffectTime,
		Rate:          cfg.Rate + addRate,
		SkillID:       skillID,
		BuffId:        buffIID,
	}

	switch constant.BuffType(buff.Pretype) {
	case constant.BuffDecoy:
		return
	case constant.BuffPercentageDefenseDecrease,
		constant.BuffPercentageAttackDecrease,
		constant.BuffPercentageDefenseHealthDecrease,
		constant.BuffPercentageCritChanceDecrease,
		constant.BuffPercentageAccuracyDecrease,
		constant.BuffPercentageAttackSpeedDecrease,
		constant.BuffPercentageMovementSpeedDecrease:
		for _, v := range cfg.Value {
			if v.Key == int32(pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE) {
				buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
			} else {
				buff.BuffStats.setStats(pb.Attribute(v.Key), -v.Value)
			}
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				if k == pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE {
					buff.BuffStats.setStats(k, v)
				} else {
					buff.BuffStats.setStats(k, -v)
				}
			}
		}

	case constant.BuffShield:
		shieldHp := float64(p.GetMaxHP()) * float64(exStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
		shieldHp = math.Max(shieldHp, 1)
		buff.BuffStats.setStats(pb.Attribute_HEALTH, int32(shieldHp))

	case constant.BuffTaunt:
		p.target = owner // 嘲讽强制改变攻击目标
	case constant.BuffCommonDebuff: // 通用减益
		for _, v := range cfg.Value {
			buff.BuffStats.setStats(pb.Attribute(v.Key), -v.Value)
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.BuffStats.setStats(k, -v)
			}
		}

	case constant.BuffCommonBuff: // 通用增益
		for _, v := range cfg.Value {
			buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.BuffStats.setStats(k, v)
			}
		}
	case constant.BuffPercentageMaxHealthIncrease2:
		if maxHPP, ok := exStats.Attrs[pb.Attribute_HEALTH_MAX_PERCENTAGE]; ok {
			restoreHP = float64(p.GetMaxHP()) * float64(exStats.getStats(pb.Attribute_HEALTH_MAX_PERCENTAGE)) / BaseRate
			p.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(p.GetMaxHP())*float64(maxHPP)/BaseRate))
			buff.Shield = int32(restoreHP)
		}
		for k, v := range exStats.Attrs {
			buff.BuffStats.setStats(k, v)
		}

		logger.Debugf("[%d:%s]添加 强脉buff 剩余护盾 %d  最大生命值 %d 当前生命值 %d", p.uid, p.name, buff.Shield, p.GetMaxHP(), p.currentHP)

	default:
		for _, v := range cfg.Value {
			buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
		}
		//logger.Debugf("[%s-%v] 增加buff: %+v", o.name, o.uid, cfg.Value)
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.BuffStats.setStats(k, v)
			}
		}
	}

	b.buffs.Set(cfgId, buff)

	//o.buffs[cfgId] = buff

	packet := &pb.AddBuffPush{
		TargetId: p.Uid(),
		Buff: &pb.BuffInfo{
			BuffId:        cfgId,
			RemainingTime: cfg.LastTime + addSec,
			BuffStats:     []*pb.Stats{},
			OwnerUid:      owner.Uid(),
		},
	}
	buff.BuffStats.Range(func(bk pb.Attribute, bv int32) bool {
		packet.Buff.BuffStats = append(packet.Buff.BuffStats, &pb.Stats{Attr: bk, Num: bv})
		return true
	})
	//for bk, bv := range buff.BuffStats.Attrs {
	//	packet.Buff.BuffStats = append(packet.Buff.BuffStats, &pb.Stats{Attr: bk, Num: bv})
	//}
	if p.desk != nil {
		p.desk.Broadcast(protocol.RouteAddBuffPush, packet)
	} else {
		p.pushMsg(protocol.RouteAddBuffPush, packet)
	}

	if restoreHP > 0 {
		p.ChangeHp(int32(restoreHP), owner)
	}
	if buff.BuffId != 0 {
		logger.Debugf("++++[%s]++++ 增加buff,  BuffId=%d BuffCfgId=%d skillId= %d  addSec= %d , stats = %v  maxhp = %d  currHp = %d 结束时间 %v",
			p.name, buff.BuffId, buff.Id, skillID, addSec, buff.BuffStats.Attrs, p.GetMaxHP(), p.currentHP, buff.OverTime.Unix())
	}

}

func (b *BuffManager) delBuff(id int32) {
	// 装备升星几率的buf如果没有过期则不清理
	p := b.object
	var buffId int64
	if buff, ok := b.buffs.Get(id); ok {
		switch buff.Pretype {
		case constant.BuffPercentEquipStarIncrease:
			if time.Now().Before(buff.OverTime) {
				return
			}
		case constant.BuffPercentageMaxHealthIncrease2:
			if buff.Shield > 0 {
				p.ChangeHp(-buff.Shield, p)

				logger.Debugf("[%d:%s]移除 强脉buff 剩余护盾 %d  最大生命值 %d  当前生命值 %d", p.uid, p.name, buff.Shield, p.GetMaxHP(), p.currentHP)

				buff.Shield = 0
			}
		}
		// if buff.Pretype == constant.BuffPercentEquipStarIncrease {
		// 	if time.Now().Before(buff.OverTime) {
		// 		return
		// 	}
		// }
		buffId = buff.BuffId
	}

	b.buffs.Delete(id)

	p.RefreshStatsAsync()

	packet := &pb.RemoveBuffPush{BuffId: id, TargetId: p.Uid()}
	if p.desk != nil {
		p.desk.Broadcast(protocol.RouteRemoveBuffPush, packet)
	} else {
		p.pushMsg(protocol.RouteRemoveBuffPush, packet)
	}

	if buffId != 0 {
		logger.Debugf("----[%s]---- 删除buff ,BuffId=%d BuffCfgId=%d ", p.name, buffId, id)
	} else {
		//logger.Debugf("----[%s]---- 删除buff ,BuffId=%d  ", p.name, id)
	}
}

func (b *BuffManager) getBuffByType(pretype constant.BuffType) *BuffInfo {
	var buffInfo *BuffInfo
	b.IterateBuff(func(buff *BuffInfo) bool {
		if buff.Pretype == int32(pretype) {
			if time.Now().Before(buff.OverTime) {
				buffInfo = buff
				return false
			}
			//logger.Debugf(" [%d:%s] buff存在但是已过期,类型 %d BuffID %d 当前时间 %d  过期时间 %d CD过期时间 %d", o.uid, o.name, pretype, v.Id, time.Now().Unix(), v.OverTime.Unix(), v.CDTime.Unix())
		}
		return true
	})
	return buffInfo
}

// client buff list
func (b *BuffManager) getPBBuffList() []*pb.BuffInfo {
	now := time.Now()
	var buffs []*pb.BuffInfo
	b.rangeBuff(func(buff *BuffInfo) bool {
		one := &pb.BuffInfo{
			BuffId:        buff.Id,
			RemainingTime: int32(buff.OverTime.UnixMilli() - now.UnixMilli()),
			BuffStats:     []*pb.Stats{},
		}

		buff.BuffStats.Range(func(bk pb.Attribute, bv int32) bool {
			one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: bk, Num: bv})
			return true
		})

		//for bk, bv := range buff.BuffStats.Attrs {
		//	one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: bk, Num: bv})
		//}

		buffs = append(buffs, one)
		return true
	})
	return buffs
}

// 清理游戏内产生的buff
func (b *BuffManager) resetBuff() {
	p := b.object
	whileBuff := make([]int32, 0)
	if player, ok := defaultManager.player(p.uid); ok {
		//if len(player.pet) != 0 {
		//	if player.pet[0].PetCfgId != 0 {
		//		// 查询宠物buff
		//		petInfo := global.GetPetCfg(player.pet[0].PetCfgId)
		//		for _, v := range petInfo.Talent {
		//			whileBuff = append(whileBuff, v.ID)
		//		}
		//	}
		//}
		player.PetMgr.getTalent(&whileBuff)
	}
	var delBuffs []int32
	b.rangeBuff(func(buff *BuffInfo) bool {
		if buff.Pretype == constant.BuffPercentEquipStarIncrease {
			return true
		}

		if buff.Mark == constant.BuffMarkPassive {
			return true
		}
		// 在白名单的buff也要跳过
		in := false
		for _, v := range whileBuff {
			if v == buff.Id {
				in = true
				break
			}
		}
		if in {
			return true
		}
		delBuffs = append(delBuffs, buff.Id)
		return true
	})
	for i := 0; i < len(delBuffs); i++ {
		b.delBuff(delBuffs[i])
	}
	for i := range p.skills {
		p.skills[i].Cooldown = time.Now()
	}
}

func (b *BuffManager) CleanDebuff() {
	var delBuffs []int32
	b.rangeBuff(func(buff *BuffInfo) bool {
		if buff.BuffType == constant.BuffTypeDebuff {
			delBuffs = append(delBuffs, buff.Id)
		}
		return true
	})
	for i := 0; i < len(delBuffs); i++ {
		b.delBuff(delBuffs[i])
	}
}

func (b *BuffManager) getBuffIds() []int32 {
	var ids []int32
	b.rangeBuff(func(buff *BuffInfo) bool {
		ids = append(ids, buff.Id)
		return true
	})
	return ids
}

func (m *BuffManager) tick() bool {
	now := time.Now()
	if now.After(m.tickTime) {
		m.tickTime = now.Add(time.Duration(100) * time.Millisecond)
		return true
	}
	return false
}
