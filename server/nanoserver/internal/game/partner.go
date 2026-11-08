package game

import (
	"math"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/astar"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"sort"
	"time"
)

// const PartnerWarriorNormalSkillID = 2000000
// const PartnerWizardNormalSkillID = 2000100
// const PartnerRitualistNormalSkillID = 2000200

type Talent struct {
	points    int64          //当前点数
	activated []global.IdNum //已激活天赋
}

func (t *Talent) talentUID(slot, level int32) int32 {
	return slot*1000 + level
}

type Partner struct {
	*MapObject
	*Talent
	star   int32 //星级
	cfg    *global.PartnerConfig
	combat int32

	searchTime time.Time
	roamPath   []*astar.APoint
	attrRatio  float64
}

func CreatePartner(cid int32) *Partner {
	uid := snowflake.GenID()
	partner := &Partner{
		MapObject: NewMapObject(ObjectTypePartner, uid),
		star:      0,
		Talent:    &Talent{points: 0, activated: make([]global.IdNum, 0)},
	}
	partner.cfg = global.GetPartnerCfg(cid)
	partner.cfgId = cid
	partner.level = 1
	partner.name = partner.cfg.Name
	partner.job = constant.MowJob(partner.cfg.Job)
	partner.attrRatio = BaseRate

	partner.initTalent()

	return partner
}

func (par *Partner) setLevel(level int32) bool {

	player := par.owner.(*Player)
	if player == nil {
		return false
	}
	if level > global.GlobalCfg.BaseValue.Rolelevel.Maxlevel {
		return false
	}

	pcfg := global.GetPartnerCfg(par.cfgId)
	if pcfg == nil {
		return false
	}

	if _, err := global.GetRoleCfg(int32(pcfg.Job), par.level); err != nil {
		player.logger.Errorf("升级失败 userId: %d,  err: %s", par.uid, err.Error())
		return false
	}

	par.level = level

	oldtotalCombat := par.combat
	player.refreshStats()
	newtotalCombat := par.combat
	par.upPartnerLevelShushuDate(player, par.level, newtotalCombat-oldtotalCombat)
	par.currentHP = par.GetMaxHP()
	par.stats.setStats(pb.Attribute_HEALTH, par.currentHP)
	player.pushMsg(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: par.level, Uid: par.uid})
	player.tasks.dispatch(pb.TargetType_UPGRADE_ANY_PARTNER_LEVEL, par.level, 1)
	return true
}

func (par *Partner) levelUp() bool {

	player := par.owner.(*Player)
	if player == nil {
		return false

	}
	if par.level >= global.GlobalCfg.BaseValue.Rolelevel.Maxlevel {
		return false
	}

	pcfg := global.GetPartnerCfg(par.cfgId)
	if pcfg == nil {
		return false
	}

	cfg, err := global.GetRoleCfg(int32(pcfg.Job), par.level)
	if err != nil {
		player.logger.Errorf("升级失败 userId: %d,  err: %s", par.uid, err.Error())
		return false
	}

	// var bUpgrade bool
	// for player.exp >= cfg.NeedExp {
	// 	player.exp -= cfg.NeedExp
	// 	if player.exp <= 0 {
	// 		player.exp = 0
	// 	}

	// 	cfg, err = global.GetRoleCfg(int32(pcfg.Job), par.level+1)
	// 	if err != nil {
	// 		player.logger.Errorf("升级失败 userId: %d,  err: %s", par.uid, err.Error())
	// 		return false
	// 	}
	// 	par.level++
	// 	bUpgrade = true
	// }
	var bUpgrade bool
	if player.exp >= cfg.NeedExp {
		player.exp -= cfg.NeedExp
		if player.exp <= 0 {
			player.exp = 0
		}
		par.level++
		bUpgrade = true
	}
	oldtotalCombat := par.combat
	if bUpgrade {
		//par.logger.Debugf("升级 userId: %d,  level: %d  vip: %d  exp: %d", par.uid, par.level, par.vip, par.exp)
		player.refreshStats()
		newtotalCombat := par.combat
		par.upPartnerLevelShushuDate(player, par.level, newtotalCombat-oldtotalCombat)
		par.currentHP = par.GetMaxHP()
		par.stats.setStats(pb.Attribute_HEALTH, par.currentHP)
		player.pushMsg(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: par.level, Uid: par.uid})
		player.tasks.dispatch(pb.TargetType_UPGRADE_ANY_PARTNER_LEVEL, par.level, 1)
	}
	return true
}
func (par *Partner) upPartnerLevelShushuDate(player *Player, level int32, combat int32) {
	//创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	//upBaseInfo的数据合并进去
	baseInfo := player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//在合并点其他的数据
	shushuDate["name"] = "partner_level_upgrade"
	shushuDate["level_before"] = level - 1
	shushuDate["level_after"] = level
	shushuDate["partner_job"] = constant.JobName[par.job]
	shushuDate["partner_power"] = int32(par.combat)
	shushuDate["partner_quality"] = par.cfg.Quality
	shushuDate["partner_star"] = par.star
	shushuDate["partner_power_up"] = combat
	shushuDate["partner_name"] = par.name
	shushumain.ReportData(shushuDate)

}

// 技能学习
func (par *Partner) learnSkill() {
	for _, v := range par.cfg.Skills {
		if par.star >= v.Key {
			if skill := par.getSkill(v.Value); skill != nil {
				continue
			}
			if skill := par.getSkillByType(v.Value); skill != nil {
				continue
			}
			_, err := global.GetSkillCfg(v.Value)
			if err == nil {
				par.skills = append(par.skills, &PlayerSkill{CfgId: v.Value})
			}
		}
	}
}

// 计算角色属性
func (par *Partner) refreshStats(packet *pb.SyncUserStatsPush, owner *Player) {
	par.slock.Lock()
	defer par.slock.Unlock()

	now := time.Now()
	par.stats.clear()
	// 基础属性
	roleCfg, err := global.GetRoleCfg(int32(par.cfg.Job), par.level)
	if err != nil {
		owner.logger.Errorf("获取角色等级配置失败 level %d", par.level)
		return
	}

	quality_rate := global.GetpartnerQualityRate(par.cfg.Quality)
	for _, rc := range roleCfg.BassProperty {
		par.stats.setStats(pb.Attribute(rc.Key), int32(math.Round(float64(rc.Value)*quality_rate*par.attrRatio/BaseRate)))
	}

	// 星级
	for _, v := range par.cfg.StarAttribute {
		if v.Id == par.star {
			for _, v2 := range v.Values {
				par.stats.setStats(pb.Attribute(v2.Key), v2.Value)
			}
			break
		}
	}

	for _, v := range par.cfg.BornValue {
		par.stats.setStats(pb.Attribute(v.Key), v.Value)
	}

	par.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(par.cfg.Job)).AttackSpeed)
	par.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(par.cfg.Job)).MoveSpeed)

	// 魂石增加生命值基数=角色基础属性+命格
	gemHealthAdd := par.stat(pb.Attribute_MAX_HEALTH)

	// 装备 强化大师
	equipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	gemStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	var maxEnhance, maxStars int32
	maxEnhance = math.MaxInt32
	maxStars = math.MaxInt32
	var starCfg, enhanceCfg *global.EquipMasterConfig

	slot_equip := [pb.EquipSlot_MaxEquipSlot]Slot{}
	slot_id := int32(0)
	slot_equip_title := int32(0)

	allGems := []int32{}
	for _, slot := range owner.partnerSlots {
		if slot == nil {
			continue
		}
		if slot.PartnerUid == par.uid {
			slot_equip = slot.Equips

			slot_id = slot.SlotId
			slot_equip_title = slot.TitleId
			break
		}
	}
	for slot, equip := range slot_equip {
		if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
			continue
		}
		gemstoneNature := make([]int32, 0)
		for GemstoneSlot, Gemstone := range equip.Gemstone {
			if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
				continue
			}
			genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
			if err != nil {
				continue
			}

			if genstoneCfg.NatureMark > 0 {
				gemstoneNature = append(gemstoneNature, genstoneCfg.NatureMark)
			}
			//对应宝石等级的数目加1
			allGems = append(allGems, genstoneCfg.Itemlevel)

			for _, attr := range genstoneCfg.Stats {
				if attr.Key == pb.Attribute_HEALTH_PERCENTAGE {
					gemStatus.setStats(attr.Key, attr.Value)
				} else {
					par.stats.setStats(attr.Key, attr.Value)
				}
			}
		}

		// 魂石自然属性
		if len(gemstoneNature) == int(pb.GemstoneSlot_MaxGemstoneSlot) && gemstoneNature[0] == gemstoneNature[1] && gemstoneNature[1] == gemstoneNature[2] && gemstoneNature[2] == gemstoneNature[3] {
			gemNatureCfg, ok := global.GetEquipMaster(constant.EquipMasterCategoryGemStoneNature*1000 + gemstoneNature[0])
			if ok != nil {
				continue
			}
			for _, v := range gemNatureCfg.Attribute {
				par.stats.setStats(pb.Attribute(v.Key), v.Value)
			}
		}

		if equip.Enhance < maxEnhance {
			maxEnhance = equip.Enhance
		}
		if equip.Stars < maxStars {
			maxStars = equip.Stars
		}
		if equip.ItemId == 0 {
			continue
		}
		equipInfo := owner.getItemByUid(equip.ItemId)
		if equipInfo == nil {
			continue
		}
		equipCfg, err := global.GetItemCfg(equipInfo.CfgId)
		if err != nil {
			continue
		}

		// 基础值攻\防: （装备基础值+槽位强化基础值)*升星百分比 , 基础属性: 装备基础值*品质百分比
		for _, attr := range equipCfg.Stats {
			if attr.Key == pb.Attribute_PHYSICAL_ATTACK || attr.Key == pb.Attribute_MAGIC_ATTACK ||
				attr.Key == pb.Attribute_PHYSICAL_DEFENSE || attr.Key == pb.Attribute_MAGIC_DEFENSE ||
				attr.Key == pb.Attribute_MAX_HEALTH {

				equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))))
			} else {
				par.stats.setStats(attr.Key, attr.Value)
			}
		}
		// 极品属性直接计算
		if equipInfo.Equip != nil {
			for _, attr := range equipInfo.Equip.Attributes {
				par.stats.setStats(attr.Attr, attr.Num)
			}
		}

		star_uid := int32(par.cfg.Job%10*100) + equip.Stars
		if equipStarCfg, err := global.GetStarsCfg(star_uid); err == nil {
			for _, star_attr := range equipStarCfg.Attribute {
				if star_attr.Slot == int32(slot) {
					equipStatus.setStats(pb.Attribute(star_attr.Attribute), star_attr.Value)
				}
			}
		}

		enhanc_uid := int32(par.cfg.Job%10*1000) + equip.Enhance
		if equipEnhancCfg, err := global.GetEnhanceCfg(enhanc_uid); err == nil {
			for _, star_attr := range equipEnhancCfg.Attribute {
				if star_attr.Slot == int32(slot) {
					if equipCfg.Quality > 0 && equipCfg.Quality <= 7 {
						equipStatus.setStats(pb.Attribute(star_attr.Attribute), int32(float32(star_attr.Value)*global.GlobalCfg.Equipenhances[equipCfg.Quality-1].Value))
					}
				}
			}
		}
	}
	equipStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAX_HEALTH, int32(float64(equipStatus.getStats(pb.Attribute_MAX_HEALTH))*(float64(equipStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	for ek, ev := range equipStatus.Attrs {
		if ek == pb.Attribute_PHYSICAL_ATTACK || ek == pb.Attribute_MAGIC_ATTACK ||
			ek == pb.Attribute_PHYSICAL_DEFENSE || ek == pb.Attribute_MAGIC_DEFENSE ||
			ek == pb.Attribute_MAX_HEALTH {
			par.stats.setStats(ek, ev)
		}
	}

	// 魂石增加生命百分比
	par.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(gemHealthAdd)*(float64(gemStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryEnhance && maxEnhance >= equipmaster.Limit {
			enhanceCfg = equipmaster
		}
		if equipmaster.Category == constant.EquipMasterCategoryStars && maxStars >= equipmaster.Limit {
			starCfg = equipmaster
		}
	}

	//魂石大师
	sort.Slice(allGems, func(i, j int) bool {
		return allGems[i] > allGems[j]
	})
	for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
		if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
			equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
			if err != nil {
				continue
			}
			if allGems[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
				for _, v := range equipMasterCfg.Attribute {
					par.stats.setStats(pb.Attribute(v.Key), v.Value)
				}
				break
			}
		}
	}

	//强化大师
	if enhanceCfg != nil {
		for _, v := range enhanceCfg.Attribute {
			par.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	//升星大师
	if starCfg != nil {
		for _, v := range starCfg.Attribute {
			par.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	//特权卡
	privilegeStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, v := range owner.privilege {
		cfg, ok := global.PrivilegeCfg[int32(v.Pid)]
		if !ok {
			continue
		}
		for _, att := range cfg.Pcontent {
			av := privilegeStatus.getStats(pb.Attribute(att.Key))
			if att.Value > av {
				privilegeStatus.changeStats(pb.Attribute(att.Key), att.Value)
			}
		}
	}
	for k, v := range privilegeStatus.Attrs {
		par.stats.setStats(k, v)
	}

	// VIP
	vipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, gv := range global.VIPCfg {
		if gv.VIPLevel <= owner.vip {
			for _, att := range gv.Pcontent {
				av := vipStatus.getStats(pb.Attribute(att.Key))
				if att.Value > av {
					vipStatus.changeStats(pb.Attribute(att.Key), att.Value)
				}
			}
		}
	}
	for k, v := range vipStatus.Attrs {
		par.stats.setStats(k, v)
	}

	par.stats.setStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE, int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	//被动技能
	for _, v := range par.skills {
		if skillCfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			if skillCfg.ActiveorPassive == PassiveSkill {
				for _, e := range skillCfg.Effect {
					par.stats.setStats(pb.Attribute(e.Key), e.Value)
				}
			}
		}
	}

	//称号共鸣
	for _, v := range par.activateTitleBonds {
		if cfg := global.GetTitleActivationCfg(v); cfg != nil {
			for _, a := range cfg.Attribute {
				par.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}
	//五行石计算
	if owner.ElementInfo != nil {
		par.stats.setStats(pb.Attribute(pb.Attribute_Main_EMPOWER), owner.ElementInfo.elementMain)
		for _, v := range owner.ElementInfo.element {
			if v.Eid == 0 {
				continue
			}
			if cfg := global.GetElementCfg(v.Eid); cfg != nil {
				for _, a := range cfg.Levelvalue {
					par.stats.setStats(pb.Attribute(a.Key), a.Value)
				}
				if v.Break != cfg.Breaklevel {
					if v.Eid-1 == 0 {
						continue
					}
					if cfg := global.GetElementCfg(v.Eid - 1); cfg != nil {
						for _, a := range cfg.Breakvalue {
							if a.Value > 0 {
								par.stats.setStats(pb.Attribute(a.Key), a.Value)
							}
						}
					}
				} else {
					if v.Eid == 0 {
						continue
					}
					if cfg := global.GetElementCfg(v.Eid); cfg != nil {
						for _, a := range cfg.Breakvalue {
							if a.Value > 0 {
								par.stats.setStats(pb.Attribute(a.Key), a.Value)
							}
						}
					}
				}
			}
		}
	}

	// 称号升星属性
	for _, v := range owner.titles.Titles {
		if cfg := global.GetTitleCfg(v.CfgId); cfg != nil {
			if cfg.Slot-1 == slot_id {
				for _, g := range cfg.Goods {
					if v.Star == g.Id {
						for _, a := range g.Values {
							par.stats.setStats(pb.Attribute(a.Key), a.Value)
						}
						break
					}
				}
			}
		}
	}
	//称号装备属性
	if slot_equip_title > 0 {
		if cfg := global.GetTitleCfg(slot_equip_title); cfg != nil {
			for _, a := range cfg.Limit {
				par.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}

	if owner.refine != nil {
		owner.refine.Stats(nil, par, pb.CharacterSlot(slot_id+1))
	}

	// 神兵
	owner.artifactStatus(par.MapObject)

	// 天赋
	for _, talent := range par.activated {
		if cfg, ok := global.PartnerTalentCfg[par.talentUID(talent.ID, talent.Num)]; ok {
			for _, attr := range cfg.Value {
				par.stats.setStats(pb.Attribute(attr.Key), attr.Value)
			}
		}
	}

	// 参与计算战斗力的属性写在上面!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	combatStats := Stats{Attrs: par.stats.copyStats()}

	par.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
		buffCfg, err := global.GetBuffCfg(buff.Id)
		if err != nil {
			return true
		}

		if buffCfg.EffectTime != 0 {
			return true
		}

		if now.After(buff.OverTime) {
			return true
		}

		// 魔法酸雨的增伤只针对技能
		if buffCfg.Pretype == constant.SkillTypeMagicAcidRain {
			return true
		}

		buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
			par.stats.setStats(key, value)
			return true
		})
		//for i, e := range buff.BuffStats.Attrs {
		//	par.stats.setStats(i, e)
		//}
		return true
	})
	//for _, buff := range par.buffs {
	//	buffCfg, err := global.GetBuffCfg(buff.Id)
	//	if err != nil {
	//		continue
	//	}
	//
	//	if buffCfg.EffectTime != 0 {
	//		continue
	//	}
	//
	//	if now.After(buff.OverTime) {
	//		continue
	//	}
	//
	//	// 魔法酸雨的增伤只针对技能
	//	if buffCfg.Pretype == constant.SkillTypeMagicAcidRain {
	//		continue
	//	}
	//
	//	for i, e := range buff.BuffStats.Attrs {
	//		par.stats.setStats(i, e)
	//	}
	//}

	// <0  减益 ; >0 增益
	par.stats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(par.GetDC())*(float64(par.stats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(par.GetMC())*(float64(par.stats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(par.GetAR())*(float64(par.stats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(par.GetMR())*(float64(par.stats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	par.stats.setStats(pb.Attribute_DODGE, int32(float64(par.stats.getStats(pb.Attribute_DODGE))*(float64(par.stats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_HIT, int32(float64(par.stats.getStats(pb.Attribute_HIT))*(float64(par.stats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(par.stats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(par.stats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_TOUGHNESS, int32(float64(par.stats.getStats(pb.Attribute_TOUGHNESS))*(float64(par.stats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	par.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(par.GetMaxHP())*float64(par.stats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))
	par.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(par.GetMaxHP())*float64(par.stats.getStats(pb.Attribute_HEALTH_MAX_PERCENTAGE))/BaseRate))

	// 攻速移速增加   基础*(1-加成)
	par.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(par.stats.getStats(pb.Attribute_ATTACK_SPEED))*(1-float64(par.stats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate)))
	par.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(par.stats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(par.stats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	par.stats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(par.GetDC()))))
	par.stats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(par.GetMC()))))
	par.stats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(par.GetAR()))))
	par.stats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(par.GetMR()))))

	par.stats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_DODGE)))))
	par.stats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_HIT)))))
	par.stats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	par.stats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_TOUGHNESS)))))
	par.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_ATTACK_SPEED)))))
	par.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(par.stats.getStats(pb.Attribute_MOVE_SPEED)))))
	par.stats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(par.GetMaxHP()))))
	if owner.curState == PlayerStatusPlaying {
		if par.currentHP > par.GetMaxHP() {
			par.currentHP = par.GetMaxHP()
		}
		par.stats.changeStats(pb.Attribute_HEALTH, par.currentHP)
		combatStats.changeStats(pb.Attribute_HEALTH, par.currentHP)
	} else {
		par.stats.changeStats(pb.Attribute_HEALTH, par.GetMaxHP())
		par.currentHP = par.GetMaxHP()
		par.displayHP = par.currentHP
		combatStats.changeStats(pb.Attribute_HEALTH, par.GetMaxHP())
	}

	combat := int32(par.calcCombatRealV2(&combatStats, true))
	par.combat = combat

	//par.logger.Debugf(" [%s] 真实战力 %d , 虚假战力  %d", par.name, combat, par.calcCombat())

	//tmpStats2 := par.stats.copyStats()

	ps := &pb.SyncUserStatsPush_PartnerStats{Attributes: []*pb.Stats{}, PartnerUid: par.uid, Combat: combat}
	par.stats.Range(func(k pb.Attribute, v int32) bool {
		if v <= 0 {
			return true
		}
		ps.Attributes = append(ps.Attributes, &pb.Stats{Attr: k, Num: v})
		return true
	})
	//for k, v := range tmpStats2 {
	//	if v <= 0 {
	//		continue
	//	}
	//	ps.Attributes = append(ps.Attributes, &pb.Stats{Attr: k, Num: v})
	//}
	if packet != nil {
		packet.Partners = append(packet.Partners, ps)
	}

	//par.logger.Debugf("属性计算 当前生命:%d  最大生命:%d  物攻: %d 魔攻: %d", par.currentHP, par.GetMaxHP(), par.GetDC(), par.GetMC())
}

func (par *Partner) calcCombatRealV2(stats *Stats, calculatePanelStats bool) int64 {
	if calculatePanelStats {
		return CalculateBattlePower(par.stats4Combat(stats), global.GlobalCfg.BaseValue.Power, global.SkillCfg, par.skills, par.job)
	} else {
		return CalculateBattlePower(stats, global.GlobalCfg.BaseValue.Power, global.SkillCfg, par.skills, par.job)
	}

	//return CalculateBattlePower(stats, global.GlobalCfg.BaseValue.Power, global.SkillCfg, par.skills, par.job)
}

// 战力（不计算buff）
func (par *Partner) calcCombatReal(stats *Stats) int32 {
	combatStats := par.stats4Combat(stats)
	if combatStats == nil {
		return 0
	}

	var skillCombat float64
	for _, v := range par.skills {
		if cfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			skillCombat += float64(cfg.Powerfactor) * (1 + float64(v.AdvancedLevel))
		}
	}

	// 减伤阈值75%
	dmt_reduce := combatStats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER) + combatStats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER) + combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)
	if dmt_reduce > 7500 {
		dmt_reduce = 7500
	}

	// 有效生命
	effective_hp := (float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))+global.GlobalCfg.BaseValue.Power.DefValue*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))+float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))/2.0)/(1-float64(dmt_reduce)/BaseRate) +
		(float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))/2.0 +
		float64(combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED))
	// logger.Debugf("  第1部分 %f   第2部分 %f    第3部分 %f   ===有效生命 %f ", (float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))+global.GlobalCfg.BaseValue.Power.DefValue*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))+float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))/2.0)/(1-float64(dmt_reduce)/BaseRate),
	// 	(float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))/2.0,
	// 	float64(combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED)), effective_hp)
	// 有效输出
	effective_power := float64(combatStats.getStats(par.AttackAttributes()))*
		(1+float64(combatStats.getStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate)*(1+skillCombat/100) +
		(float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)))/2.0 +
		float64(combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_FIXED))

	// logger.Debugf(" 面板攻击 %f  (1+增伤百分比+爆伤) %f  技能加成 %f   固定值 %f   ===有效输出 %f ", float64(combatStats.getStats(par.AttackAttributes())),
	// 	(1 + float64(combatStats.getStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate),
	// 	(1 + skillCombat/100), (float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)))/2.0, effective_power)

	combat := math.Sqrt(effective_hp*global.GlobalCfg.BaseValue.Power.LifeValue)*math.Sqrt(effective_power*global.GlobalCfg.BaseValue.Power.AttackValue) +
		float64(combatStats.getStats(pb.Attribute_HIT))*global.GlobalCfg.BaseValue.Power.Hitvalue +
		float64(combatStats.getStats(pb.Attribute_DODGE))*global.GlobalCfg.BaseValue.Power.Dodgevalue +
		float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE))*global.GlobalCfg.BaseValue.Power.Critvalue +
		float64(combatStats.getStats(pb.Attribute_TOUGHNESS))*global.GlobalCfg.BaseValue.Power.Toughness

	combat += float64(par.cfg.Extropower)
	//logger.Debugf(" 终最战力 %f", combat)
	par.combat = int32(math.Round(combat))

	return par.combat
}

func (par *Partner) stats4Combat(combatStats *Stats) *Stats {

	// <0  减益 ; >0 增益
	combatStats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	combatStats.setStats(pb.Attribute_DODGE, int32(float64(combatStats.getStats(pb.Attribute_DODGE))*(float64(combatStats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_HIT, int32(float64(combatStats.getStats(pb.Attribute_HIT))*(float64(combatStats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_TOUGHNESS, int32(float64(combatStats.getStats(pb.Attribute_TOUGHNESS))*(float64(combatStats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))*float64(combatStats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))

	// 攻速移速增加  基础*(1-加成)
	combatStats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED))*(1-float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate)))
	combatStats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(combatStats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(combatStats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	combatStats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK)))))
	combatStats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK)))))
	combatStats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE)))))
	combatStats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))))

	combatStats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_DODGE)))))
	combatStats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_HIT)))))
	combatStats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	combatStats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_TOUGHNESS)))))
	combatStats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED)))))
	combatStats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MOVE_SPEED)))))
	combatStats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAX_HEALTH)))))

	return combatStats
}

// 单个道具战力
func (par *Partner) calcCombatEquip(equipInfo *pb.ItemInfo) int32 {
	if equipInfo == nil {
		return 0
	}
	equipCfg, err := global.GetItemCfg(equipInfo.CfgId)
	if err != nil {
		return 0
	}
	totalStatus := &Stats{Attrs: make(map[pb.Attribute]int32)}
	equipStatus := &Stats{Attrs: make(map[pb.Attribute]int32)}
	// 基础值攻\防: （装备基础值+槽位强化基础值)*升星百分比
	for _, attr := range equipCfg.Stats {
		if attr.Key == pb.Attribute_PHYSICAL_ATTACK || attr.Key == pb.Attribute_MAGIC_ATTACK ||
			attr.Key == pb.Attribute_PHYSICAL_DEFENSE || attr.Key == pb.Attribute_MAGIC_DEFENSE ||
			attr.Key == pb.Attribute_MAX_HEALTH {
			equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))))
		} else {
			totalStatus.setStats(attr.Key, attr.Value)
		}
	}
	// 极品属性直接计算
	if equipInfo.Equip != nil {
		for _, attr := range equipInfo.Equip.Attributes {
			totalStatus.setStats(attr.Attr, attr.Num)
		}
	}

	equipStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, par.GetDC()+int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_ATTACK, par.GetMC()+int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, par.GetAR()+int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_DEFENSE, par.GetMR()+int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAX_HEALTH, par.GetMaxHP()+int32(float64(equipStatus.getStats(pb.Attribute_MAX_HEALTH))*(float64(equipStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	// equipStatus.setStats(pb.Attribute_DODGE, par.GetDodge()+int32(float64(equipStatus.getStats(pb.Attribute_DODGE))*(float64(equipStatus.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_HIT, par.GetHit()+int32(float64(equipStatus.getStats(pb.Attribute_HIT))*(float64(equipStatus.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_CRITICAL_STRIKE, par.GetCrit()+int32(float64(equipStatus.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(equipStatus.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_TOUGHNESS, par.GetTough()+int32(float64(equipStatus.getStats(pb.Attribute_TOUGHNESS))*(float64(equipStatus.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))

	for ek, ev := range equipStatus.Attrs {
		if ek == pb.Attribute_PHYSICAL_ATTACK || ek == pb.Attribute_MAGIC_ATTACK ||
			ek == pb.Attribute_PHYSICAL_DEFENSE || ek == pb.Attribute_MAGIC_DEFENSE ||
			ek == pb.Attribute_MAX_HEALTH {
			totalStatus.setStats(ek, ev)
		}
	}

	//logger.Debugf("--------装备属性 %v", totalStatus)
	return int32(par.calcCombatRealV2(totalStatus, true))
}

func (par *Partner) starUp() error {
	if par.cfg == nil {
		return errutil.ErrPartnerInvalid
	}

	var consume *global.KeyValueEx
	for _, v := range par.cfg.StarConsume {
		if v.Id == par.star+1 {
			consume = &v
			break
		}
	}

	if consume == nil {
		return errutil.ErrPartnerInvalid
	}

	owner := par.owner.(*Player)
	if owner == nil {
		return errutil.ErrPartnerInvalid
	}

	for _, v := range consume.Values {
		if !owner.checkTakeSource(v.Key, int64(v.Value)) {
			return errutil.ErrItemNotEnough
		}
	}

	for _, v := range consume.Values {
		owner.takeSource(v.Key, int64(v.Value))
	}

	par.star++

	//解锁技能
	par.learnSkill()
	par.addTalentPoint(par.star)

	par.SyncUserPartner()
	oldcombat := par.combat
	owner.checkPartnerBonds()
	owner.refreshStats()
	newcombat := par.combat
	owner.rankChangeTCP(constant.RankTypeRedPartnerStar, owner.calcRedPartnerStarNum())

	par.upPartnerStarShushuDate(owner, par.star, newcombat-oldcombat)
	return nil
}
func (par *Partner) upPartnerStarShushuDate(player *Player, star int32, combat int32) {
	//创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	//upBaseInfo的数据合并进去
	baseInfo := player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//在合并点其他的数据
	shushuDate["name"] = "partner_star_upgrade"
	shushuDate["star_before"] = star - 1
	shushuDate["star_after"] = star
	shushuDate["partner_job"] = constant.JobName[par.job]
	shushuDate["partner_power"] = int32(par.combat)
	shushuDate["partner_quality"] = par.cfg.Quality
	shushuDate["partner_power_up"] = combat
	shushuDate["partner_level"] = par.level
	shushuDate["partner_star"] = par.star
	shushuDate["partner_name"] = par.name
	shushumain.ReportData(shushuDate)
}
func (par *Partner) SyncUserPartner() {

	packet := &pb.PartnerUpdatePush{Partner: &pb.PartnerInfo{}}
	packet.Partner.Uid = par.uid
	packet.Partner.CfgId = par.cfgId
	packet.Partner.Level = par.level
	packet.Partner.Star = par.star
	packet.Partner.Skills = make([]*pb.SkillInfo, 0)
	for _, v := range par.skills {
		packet.Partner.Skills = append(packet.Partner.Skills, &pb.SkillInfo{CfgId: v.CfgId, Advanced: v.AdvancedLevel})
	}
	packet.Partner.Points = par.points
	for _, v := range par.activated {
		packet.Partner.Talents = append(packet.Partner.Talents, &pb.KeyValue{K: v.ID, V: v.Num})
	}

	par.owner.PushMsg(protocol.RoutePartnerUpdatePush, packet)
}

func (par *Partner) processAI() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	if par.desk == nil {
		return
	}
	if par.dead {
		owner_, _ := par.desk.playerWithId(par.owner.Uid())
		if owner_ == nil {
			return
		}
		if skill, caster := owner_.getSlotSkillByType(constant.SkillTypeReborn); skill != nil {
			if time.Now().After(skill.Cooldown) && !caster.dead {
				caster.lastGasp(par, skill.CfgId)
			}
		}
		return
	}

	par.processHp()
	par.processRoam()
	par.processSearch()
	par.processTarget()
	par.processBuff()
}

// 巡逻
func (par *Partner) processRoam() {
	if par.desk == nil {
		return
	}

	if par.owner == nil {
		return
	}

	owner_, _ := par.desk.playerWithId(par.owner.Uid())
	if owner_ == nil {
		if r, err := par.desk.robotWithId(par.owner.Uid()); err == nil {
			owner_ = r.Player
		}
	}

	if owner_ == nil {
		return
	}

	if !par.CanMove() {
		return
	}

	if par.desk.roomCfg.Roomtype == constant.PVPOneOnOneMap {
		if maputil.Distance(par.location, owner_.Location()) <= 10 && par.target != nil {
			return
		}
	} else {
		if maputil.Distance(par.location, owner_.Location()) <= 4 && par.target != nil {
			return
		}
	}

	self_slot := 0
	flow_slot := -1
	var flow_point *pb.Point
	for i := 1; i <= 4; i++ {
		if spar := owner_.getPartnerInfoBySlotId(int32(i)); spar != nil && spar.uid == par.uid {
			self_slot = i
			break
		}
	}

	for i := 4; i >= 1; i-- {
		if fpar := owner_.getPartnerInfoBySlotId(int32(i)); fpar != nil && fpar.uid != par.uid && !fpar.dead {
			if self_slot > i {
				flow_point = fpar.Location()
				flow_slot = i
				break
			}
		}
	}
	if flow_point == nil {
		flow_point = owner_.Location()
	}
	distance := 2
	if par.desk.roomCfg.Roomtype == constant.PVPOneOnOneMap {
		distance = 8
	}

	if owner_.GetObAction() == ObjectActionMove && maputil.Distance(par.location, flow_point) > distance {
		//logger.Debugf("[走走走]  我是【%d】号 , 跟随[%d] 号  终点【%v】, 距离 【%d】", self_slot, flow_slot, flow_point, maputil.Distance(par.location, flow_point))
		par.target = nil
		if path := maputil.GetPath(par.location, flow_point, &par.desk.mapCfg.AstarMgr); len(path) > 1 {
			par.moveTo(path[1])
		} else {
			logger.Debugf("Par 回到原点寻路失败 起点=%v  终点=%v  跟随[%d]", par.location, flow_point, flow_slot)
		}
	}

	if owner_.GetObAction() == ObjectActionStand {
		owner_.target = nil
		if stand_cell := par.standLocation(self_slot); stand_cell != nil {

			if maputil.Distance(par.location, global.ToPointF(&stand_cell.location)) > 1 {

				//logger.Debugf("[停停停]   我是【%d】号,坐标[%v]  距离主人 【%d】 我要移动到 [%v]", self_slot, par.location, maputil.Distance(par.location, flow_point), stand_cell.location)
				if path := maputil.GetPath(par.location, global.ToPointF(&stand_cell.location), &par.desk.mapCfg.AstarMgr); len(path) > 1 {
					par.moveTo(path[1])
				} else if len(path) == 1 {
					par.moveTo((*astar.APoint)(&stand_cell.location))
				} else {
					logger.Debugf("Par 回到原点寻路失败 起点=%v  终点=%v", par.location, flow_point)
				}
			}
		} else if maputil.Distance(par.location, owner_.Location()) > distance {
			cell := par.desk.mapInfo.getAvailableCell(&owner_.GetCell().location, 2)
			if cell != nil {
				if path := maputil.GetPath(par.location, global.ToPointF(&cell.location), &par.desk.mapCfg.AstarMgr); len(path) > 1 {
					par.moveTo(path[1])
				} else {
					logger.Debugf("Par 回到原点寻路失败 起点=%v  终点=%v", par.location, flow_point)
				}
			}
		}

	}
}

func (par *Partner) standLocation(par_slot int) *Cell {
	if par.owner == nil {
		return nil
	}

	var cell *Cell
	var point *global.PointI
	dir := par.owner.GetDirection()

	desk := par.owner.Desk()
	if desk == nil {
		return nil
	}
	switch par_slot {
	case 1:
		point = maputil.Move(&par.owner.GetCell().location, maputil.ShiftDirection(dir, 1), 3)
		if cell = desk.mapInfo.getCell(point); cell != nil && !cell.block {
			return cell
		}
	case 2:
		point = maputil.Move(&par.owner.GetCell().location, maputil.ShiftDirection(dir, -1), 3)
		if cell = desk.mapInfo.getCell(point); cell != nil && !cell.block {
			return cell
		}
	case 3:
		point = maputil.Move(&par.owner.GetCell().location, maputil.ShiftDirection(dir, 3), 2)
		if cell = desk.mapInfo.getCell(point); cell != nil && !cell.block {
			return cell
		}
	case 4:
		point = maputil.Move(&par.owner.GetCell().location, maputil.ShiftDirection(dir, -3), 2)
		if cell = desk.mapInfo.getCell(point); cell != nil && !cell.block {
			return cell
		}
	default:
		return nil
	}
	return nil
}

// 寻找
func (par *Partner) processSearch() {
	// 嘲讽状态不在攻击其他目标
	if par.getBuffByType(constant.BuffTaunt) != nil && par.target != nil {
		return
	}

	if time.Now().Before(par.searchTime) {
		return
	}
	par.searchTime = time.Now().Add(time.Millisecond * 300)

	if !par.CanAttack() {
		return
	}

	var distance int
	var closest IMapObject

	if par.owner.Target() != nil && par.owner.Target().Desk() != nil && par.desk.roomCfg.Roomtype != constant.PVPOneOnOneMap {
		par.target = par.owner.Target()
	} else {

		if par.desk.roomCfg.Roomtype == constant.PVEPeacefulMap || par.desk.roomCfg.Roomtype == constant.BattleMap || par.desk.roomCfg.Roomtype == constant.BossMap ||
			par.desk.roomCfg.Roomtype == constant.TrialTowerMap || par.desk.roomCfg.Roomtype == constant.RelicMap || par.desk.roomCfg.Roomtype == constant.GuildMap ||
			par.desk.roomCfg.Roomtype == constant.PVEWarMap {

			for _, v := range par.desk.monsters {
				if v.dead {
					continue
				}
				if v.uid == par.uid {
					continue
				}
				if v.cfg != nil && v.cfg.MType == MTypeGathering {
					continue
				}
				tmpDistance := maputil.Distance(par.owner.Location(), v.Location())
				if tmpDistance > par.cfg.AttackDistance {
					continue
				}

				if distance == 0 || tmpDistance < distance {
					distance = tmpDistance
					closest = v
					if distance <= 2 {
						break
					}
				}
			}
		} else {
			for _, v := range par.desk.players {
				if v.dead {
					continue
				}
				if v.uid == par.uid {
					continue
				}
				if par.owner == v {
					continue
				}

				tmpDistance := maputil.Distance(par.location, v.Location())
				if distance == 0 || tmpDistance < distance {
					distance = tmpDistance
					closest = v
				}
			}
			for _, v := range par.desk.robots {
				if v.dead {
					continue
				}
				if v.uid == par.uid {
					continue
				}
				if par.owner != nil && par.owner.Uid() == v.uid {
					continue
				}

				tmpDistance := maputil.Distance(par.location, v.Location())
				if distance == 0 || tmpDistance < distance {
					distance = tmpDistance
					closest = v
				}
			}
			for _, v := range par.desk.objects {
				if v.objectType == ObjectTypeDecoy || v.objectType == ObjectTypePartner {
					tmpDistance := maputil.Distance(par.location, v.Location())
					//pvp 房间伙伴扩大索敌范围
					if tmpDistance > par.cfg.AttackDistance+20 {
						continue
					}
					if !par.CanAttackTarget(v) {
						continue
					}

					//优先以战士职业位目标
					if closest != nil {
						if closest.IsWarrior() {
							if distance-tmpDistance < 3 {
								continue
							}
						} else {
							if v.IsWarrior() && distance-tmpDistance < 2 {
								distance = tmpDistance
								if closest.Job() == constant.PartnerWarrior {
									continue
								}
								closest = v
								continue
							}
						}
					}

					if distance == 0 || tmpDistance < distance {
						distance = tmpDistance
						closest = v
					}
				}
			}
		}

	}
	if closest != nil {
		par.target = closest
		par.roamPath = []*astar.APoint{}

		//logger.Debugf("Par 锁定目标 mUid: %d , 自身坐标:%v, 目标ID: %d  目标坐标: %v  目标类型:%v  距离 %d",
		//par.uid, par.location, par.target.Uid(), par.target.Location(), par.target.ObjectType(), maputil.Distance(par.location, closest.Location()))
	}

}

// 目标
func (par *Partner) processTarget() {

	if par.desk == nil {
		return
	}

	if !par.CanAttack() {
		return
	}

	if par.target == nil {
		return
	}

	if par.target.Dead() || par.target.Uid() == 0 || par.target.Desk() == nil {
		par.target = nil
		return
	}

	if par.target.Desk() != par.desk {
		par.target = nil
		return
	}

	//战场安全区内不攻击
	if par.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if par.desk.mapInfo.safeArea(par.location) {
			return
		}
	}

	if par.owner == nil {
		logger.Warningf(" %d %s Parther owner is nil ", par.uid, par.cfg.Name)
		return
	}

	tmpdistance := 8
	if par.desk.roomCfg.Roomtype == constant.PVPOneOnOneMap {
		tmpdistance = 16
	}
	// 目标如果离主人过远 ，则不攻击
	if maputil.Distance(par.owner.Location(), par.target.Location()) > tmpdistance {
		par.target = nil
		return
	}

	if !par.CanAttackTarget(par.target) {
		return
	}

	if par.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if par.desk.mapInfo.safeArea(par.target.Location()) {
			par.target = nil
			return
		}
	}

	if !par.inAttackRange(int(par.cfg.AttackDistance)) {
		if !par.CanMove() {
			return
		}

		if path := maputil.GetPath(par.location, par.target.Location(), &par.desk.mapCfg.AstarMgr); len(path) > 1 {
			par.moveTo(path[1])
		} else {
			logger.Debugf("Par 寻路失败 起点=%v  终点=%v  目标坐标 %v", par.location, par.target.Location(), par.target.Location())
		}

	} else {
		if len(par.cell.objects) > 1 {
			cell := par.desk.mapInfo.getAvailableCell(global.ToPointI(par.target.Location()), int32(par.cfg.AttackDistance))
			if cell != nil {
				if path := maputil.GetPath(par.location, &pb.Point{X: float64(cell.location.X), Y: float64(cell.location.Y)}, &par.desk.mapCfg.AstarMgr); len(path) > 1 {
					par.moveTo(path[1])
					return
				}
			}
		}

		if !par.IsWarrior() && maputil.Distance(par.location, par.target.Location()) < 2 && par.getBuffByType(constant.BuffKnockback) == nil {
			if par.CanMove() && par.getBuffByType(constant.BuffTaunt) == nil {
				var cell *Cell
				var point *global.PointI
				dir := maputil.DirectionFromPointI(&par.cell.location, &par.target.GetCell().location)
				point = maputil.Move(&par.cell.location, maputil.ShiftDirection(dir, 4), 1)
				if cell = par.desk.mapInfo.getCell(point); cell != nil && !cell.block {
					par.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
					par.moveTime = par.moveTime.Add(time.Duration(par.stat(pb.Attribute_MOVE_SPEED)*4) * time.Millisecond)
					return
				}
				point = maputil.Move(&par.cell.location, maputil.ShiftDirection(dir, 2), 1)
				if cell = par.desk.mapInfo.getCell(point); cell != nil && !cell.block {
					par.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
					par.moveTime = par.moveTime.Add(time.Duration(par.stat(pb.Attribute_MOVE_SPEED)*4) * time.Millisecond)
					return
				}
				point = maputil.Move(&par.cell.location, maputil.ShiftDirection(dir, -2), 1)
				if cell = par.desk.mapInfo.getCell(point); cell != nil && !cell.block {
					par.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
					par.moveTime = par.moveTime.Add(time.Duration(par.stat(pb.Attribute_MOVE_SPEED)*4) * time.Millisecond)
					return
				}
			}
		}

		par.attack()
	}
}

// 移动
func (par *Partner) moveTo(point *astar.APoint) bool {
	if par.desk == nil {
		return false
	}

	if !par.CanMove() {
		return false
	}

	//from := par.location
	to := &pb.Point{X: float64(point.X), Y: float64(point.Y)}
	if !par.desk.mapInfo.locationChange(par, to) {
		return false
	}

	par.moveTime = time.Now().Add(time.Duration(par.getMoveSpeed()) * time.Millisecond)

	packet := &pb.MovePush{Uid: par.uid, Location: par.location, Direction: par.direction}
	par.desk.Broadcast(protocol.RouteMovePush, packet)

	// if par.target != nil {
	// 	logger.Debugf("Par 向目标移动 mUid: %d  起点 %v, =====目标位置: %v   目标ID %d", par.uid, from, par.target.Location(), par.target.Uid())
	// } else {
	// 	logger.Debugf("Par 没有目标移动 mUid: %d  起点 %v, 目标位置: %v", par.uid, from, to)
	// }

	return true
}

// 攻击
func (par *Partner) attack() {
	if par.desk == nil {
		logger.Errorf("par 攻击失败，伙伴不在房间内 %d:%s ", par.uid, par.name)
		return
	}
	if !par.CanAttack() {
		return
	}

	if !par.thresholdDetect() {
		return
	}

	// 技能CD ,
	var castSkill *PlayerSkill
	var castSkillId int32
	normalSkillId := global.GetNormalSkillId(constant.MowJob(par.cfg.Job))
	if par.Owner() != nil && par.Owner().IsRobot() {
		for _, v := range par.skills {
			if v.CfgId == 0 {
				continue
			}
			if cfg, err := global.GetSkillCfg(castSkillId); err == nil {
				if cfg.ActiveorPassive == PassiveSkill {
					continue
				}
				if cfg.SType == constant.SkillTypeNormal {
					continue
				}
			}

			if v.Cooldown.Before(time.Now()) {
				castSkill = v
				break
			}
		}
	}
	if castSkill == nil {
		castSkillId = normalSkillId
	} else {
		castSkillId = castSkill.CfgId
	}

	skillCfg, _ := global.GetSkillCfg(castSkillId)
	if skillCfg == nil {
		return
	}

	//优先设置攻击间隔
	par.setAttackTime(global.GetActionsetCfg(skillCfg.SkillID))

	Radius := skillCfg.Radius
	Cooldown := int32(0)
	//技能CD
	if castSkillId != normalSkillId && skillCfg.CoolDown > 0 {

		// 技能进阶效果
		if castSkill != nil {
			if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, castSkill.AdvancedLevel); err == nil {
				Cooldown = skillUpperlevel.Cooldown
				Radius += skillUpperlevel.Distance
			}
		}
		CD := skillCfg.CoolDown - Cooldown
		if CD < 0 {
			CD = 0
		}
		castSkill.Cooldown = time.Now().Add(time.Duration(CD) * time.Millisecond)
	}

	attackCells := []*Cell{}
	targets := []IMapObject{}
	targets = append(targets, par.target)

	// 技能影响到的全部格子
	switch skillCfg.SkillID {

	case constant.SkillTypeHalfMoon: // 扇形
		attackCells = par.desk.mapInfo.getSectorCells(&par.cell.location, par.direction, Radius)

	case constant.SkillTypeFourSides: // 四方向
		for i := 0; i < 4; i++ {
			direction := maputil.ShiftDirection(par.direction, i*2)
			attackCells = append(attackCells, par.desk.mapInfo.getLineCells(&par.cell.location, direction, Radius)...)
		}

	case constant.SkillTypeFreezeTime, constant.SkillTypeTimeSpaceLock, constant.SkillTypeDuel, constant.SkillTypeRepelFoe,
		constant.SkillTypeMassHeal, constant.SkillTypeDivineStrength, constant.SkillTypeArcaneMastery, constant.SkillTypeBlessingOfCourage,
		constant.SkillTypeNetherSpell, constant.SkillTypeDivineBattleArmor, constant.SkillTypeTaunt, constant.SkillTypeStomp,
		constant.SkillTypeBleeding, constant.SkillTypeWreck, constant.SkillTypeApocalypse: // 正方形
		attackCells = par.desk.mapInfo.getSquareCells(&par.cell.location, Radius)

	case constant.SkillTypeBreakArmy, constant.SkillTypeAuroraPulse, constant.SkillTypeEarthquakeWave: // 长方形
		attackCells = par.desk.mapInfo.getLineCells(&par.cell.location, par.direction, Radius)
	}

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
				if par.CanAttackTarget(obj) && par.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == par.uid {
			continue
		}

		// 安全区内不能攻击玩家
		if par.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (par.desk.mapInfo.safeArea(par.location) || par.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		// 技能距离
		if skillCfg.SType != constant.SkillTypeNormal && skillCfg.MaxDistance > 0 && maputil.Distance(par.location, t.Location()) > int(skillCfg.MaxDistance) {
			continue
		}
		//伤害延迟生效
		pardesk := par.desk
		part := t
		ilocation := &par.cell.location
		time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
			if pardesk != nil {
				pardesk.addDelayAction(ObjectTypePartner, par.uid, castSkillId, part, par.direction, ilocation)
			}
		})
	}

	par.direction = maputil.DirectionFromPointF(par.location, par.target.Location())

	packet := &pb.AttackPush{AttackerId: par.uid,
		TargetUids: []int64{par.target.Uid()},
		SkillId:    int32(castSkillId),
		Location:   par.target.Location(),
		Direction:  par.direction}
	par.desk.Broadcast(protocol.RouteAttackPush, packet)

	par.SkillBuffTrigger(castSkill)
	//logger.Debugf("伙伴攻击 ParUid:%d Name:%s, 技能ID:%d 目标uid:%d,  目标坐标 %v 目标类型 %v ", par.uid, par.name, castSkillId, par.target.Uid(), par.target.Location(), par.target.ObjectType())
}

func (par *Partner) getSkill(skillId int32) *PlayerSkill {
	for _, v := range par.skills {
		if v.CfgId == skillId {
			return v
		}
	}
	return nil
}

// 获取技能根据类型
func (par *Partner) getSkillByType(skillType int32) *PlayerSkill {
	for _, v := range par.skills {

		if v.CfgId/100 == skillType/100 {
			return v
		}
	}
	return nil
}

func (par *Partner) revive() {

	if par.desk == nil || par.Owner() == nil {
		logger.Errorf(" 伙伴复活失败 desk:[%t] 或者 owner:[%t] 为空 伙伴ID %d", par.desk == nil, par.Owner() == nil, par.uid)
		return
	}
	cell := par.desk.mapInfo.getAvailableCell(&par.Owner().GetCell().location, 3)
	if cell == nil {
		logger.Errorf("[%s] 伙伴复活失败,未找到复活点  伙伴ID %d", par.Owner().Name(), par.uid)
		return
	}
	par.desk.mapInfo.locationChange(par, global.ToPointF(&cell.location))

	par.dead = false
	par.reviveTimes++
	par.currentHP = par.stat(pb.Attribute_MAX_HEALTH)
	par.displayHP = par.currentHP
	par.stats.setStats(pb.Attribute_HEALTH, par.currentHP)
	par.resetBuff()
	par.desk.Broadcast(protocol.RouteRevivePush, &pb.RevivePush{
		Uid:       par.uid,
		Location:  par.location,
		Direction: par.direction,
		CurHp:     par.currentHP})

	par.AddBuff(int32(constant.GOD_MODE), par, int32(global.GlobalCfg.GVG.Relive.Protecttime*1000), nil, 0, 0)

	//logger.Debugf("[%s]伙伴复活！！！   当前血量 %d -- %d --- %d -- [%v]", par.name, par.currentHP, par.stat(pb.Attribute_HEALTH), par.GetMaxHP(), par.Location())
}

// 深度拷贝伙伴信息
func (par *Partner) deepCopy() *Partner {
	newPar := CreatePartner(par.cfgId)
	newPar.level = par.level
	newPar.star = par.star
	newPar.skills = []*PlayerSkill{}

	for _, v := range par.skills {
		newPar.skills = append(newPar.skills, &PlayerSkill{CfgId: v.CfgId, AdvancedLevel: v.AdvancedLevel, Cooldown: time.Now(), Distance: v.Distance})
	}

	par.stats.Range(func(k pb.Attribute, v int32) bool {
		newPar.stats.setStats(k, v)
		return true
	})
	//for k, v := range par.stats.copyStats() {
	//	newPar.stats.setStats(k, v)
	//}
	newPar.currentHP = par.GetMaxHP()
	newPar.displayHP = newPar.currentHP

	newPar.Talent = par.Talent

	return newPar
}

// 初始化天赋
func (par *Partner) initTalent() {

	if par.Talent != nil && len(par.activated) > 0 {
		return
	}
	par.Talent = &Talent{activated: make([]global.IdNum, 0)}

	for _, id := range par.cfg.TalentList {
		par.activated = append(par.activated, global.IdNum{ID: id, Num: 1})
	}
	for _, tp := range par.cfg.TalentPoints {
		if par.star >= tp.Key {
			par.points += int64(tp.Value)
		}
	}

	logger.Debugf("%d:%s 已有伙伴天赋初始化, 伙伴星级 %d ,天赋点数 %d", par.uid, par.name, par.star, par.points)
}

// 增加天赋点
func (par *Partner) addTalentPoint(star int32) {
	for _, v := range par.cfg.TalentPoints {
		if v.Key == star {
			par.points += int64(v.Value)
			break
		}
	}
}

// 天赋加点
func (par *Partner) talentActivate(talentId int32) error {
	for K, v := range par.activated {
		if v.ID == talentId {
			if _, ok := global.PartnerTalentCfg[par.talentUID(v.ID, v.Num+1)]; !ok {
				return errutil.ErrConfigNotFound
			}
			cfg, ok := global.PartnerTalentCfg[par.talentUID(v.ID, v.Num)]
			if !ok {
				return errutil.ErrConfigNotFound
			}
			owner := par.owner.(*Player)
			if owner == nil {
				return errutil.ErrPartnerInvalid
			}
			if par.points+owner.getItemCountInBag(cfg.Consume.Key) < int64(cfg.Consume.Value) {
				return errutil.ErrItemNotEnough
			}
			// 伙伴自身天赋点 ， 已使用的通用天赋点 ， 总投入的天赋点， 本次升级需要的通用天赋点
			var parPoints, comPoints, totalPoints, needComPoints int64

			for _, v := range par.activated {
				for i := 1; i < int(v.Num); i++ {
					cfg, ok := global.PartnerTalentCfg[par.talentUID(v.ID, int32(i))]
					if !ok {
						continue
					}
					totalPoints += int64(cfg.Consume.Value)
				}
			}
			for _, tp := range par.cfg.TalentPoints {
				if par.star >= tp.Key {
					parPoints += int64(tp.Value)
				}
			}
			// 此伙伴已使用的通用天赋点 = 总投入的天赋点 - （伙伴自身总天赋点-剩余未使用的天赋点)
			comPoints = totalPoints - (parPoints - par.points)

			pointLimit := global.XMLQualityIdNum{}
			for _, v := range global.GlobalCfg.PartnerTalent.Cost {
				if v.Quality == par.cfg.Quality {
					pointLimit = v
					break
				}
			}

			if comPoints >= pointLimit.Num {
				return errutil.ErrCommonTalentPointLimit
			}
			//优先使用自身点数，然后再消耗通用点数，如果通用点数使用达到上限 ：通用天赋点数使用已达上限！
			if par.points < int64(cfg.Consume.Value) {
				needComPoints = int64(cfg.Consume.Value - int32(par.points))
				if needComPoints > pointLimit.Num-comPoints {
					return errutil.ErrCommonTalentPointLimit
				}
			}

			//前置天赋
			prePoints := 0
			for _, pp := range par.activated {
				if pp.ID == cfg.Limit.Key {
					prePoints = int(pp.Num)
					break
				}
			}
			if prePoints < int(cfg.Limit.Value) {
				return errutil.ErrPreTalentNotEnough
			}

			// 扣除天赋点
			needPoints := int64(cfg.Consume.Value)
			if par.points >= needPoints {
				par.points -= needPoints
			} else {
				needPoints -= par.points
				par.points = 0
				owner.takeSource(cfg.Consume.Key, needPoints)
			}

			par.activated[K].Num++
			owner.logger.Debugf("%d:%s 升级天赋， 天赋ID %d 等级 %d 剩余通用天赋点 %d 伙伴自身天赋点 %d",
				par.uid, par.cfg.Name, talentId, par.activated[K].Num, owner.getItemCountInBag(constant.ItemTypePartnerTalentPoint), par.points)

			par.SyncUserPartner()
			owner.syncUserBase()
			owner.RefreshStats()
			break
		}
	}
	return nil
}

// 天赋重置
func (par *Partner) talentReset() error {

	owner := par.owner.(*Player)
	if owner == nil {
		return errutil.ErrPartnerInvalid
	}

	cost := global.XMLQualityIdNum{}
	for _, v := range global.GlobalCfg.PartnerTalent.Cost {
		if v.Quality == par.cfg.Quality {
			cost = v
			break
		}
	}
	if cost.ID == 0 || cost.Num == 0 {
		return errutil.ErrConfigNotFound
	}

	if !owner.checkTakeSource(cost.ID, int64(cost.Num)) {
		return errutil.ErrItemNotEnough
	}

	// 重置天赋点
	var parPoints, comPoints, totalPoints int64

	for _, v := range par.activated {
		for i := 1; i < int(v.Num); i++ {
			cfg, ok := global.PartnerTalentCfg[par.talentUID(v.ID, int32(i))]
			if !ok {
				continue
			}
			totalPoints += int64(cfg.Consume.Value)
		}
	}
	for _, tp := range par.cfg.TalentPoints {
		if par.star >= tp.Key {
			parPoints += int64(tp.Value)
		}
	}
	comPoints = totalPoints - (parPoints - par.points)

	if totalPoints == 0 {
		return errutil.ErrPartnerInvalid
	}

	owner.takeSource(cost.ID, int64(cost.Num), constant.ItemDelParTalent)

	par.points = parPoints
	for K := range par.activated {
		par.activated[K].Num = 1
	}

	owner.logger.Debugf("%d:%s 重置天赋 总投入天赋点 %d 返还通用天赋点 %d . 返回后剩余通用天赋点 %d 伙伴自身天赋点 %d",
		par.uid, par.cfg.Name, totalPoints, comPoints, owner.getItemCountInBag(constant.ItemTypePartnerTalentPoint), par.points)

	owner.addSource(constant.ItemTypePartnerTalentPoint, comPoints, constant.ItemAddParTalent)

	par.SyncUserPartner()
	owner.syncUserBase()
	owner.RefreshStats()
	return nil
}
