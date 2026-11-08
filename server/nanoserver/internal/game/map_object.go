package game

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/event"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
)

type ObjectType int

const (
	ObjectTypePlayer  ObjectType = 1
	ObjectTypePartner ObjectType = 2
	ObjectTypeMonster ObjectType = 3
	ObjectTypeRobot   ObjectType = 4
	ObjectTypeBuff    ObjectType = 5
	ObjectTypeDecoy   ObjectType = 6
)

type ObjectAction int

const (
	ObjectActionStand  ObjectAction = 1
	ObjectActionMove   ObjectAction = 2
	ObjectActionAttack ObjectAction = 3
	ObjectActionDead   ObjectAction = 4
	ObjectActionMining ObjectAction = 5
)

type IMapObject interface {
	PushMsg(route string, v interface{}) error

	ObjectType() ObjectType // 实例类型
	Location() *pb.Point    // 位置
	GetStats() *Stats       // 全部属性
	Uid() int64
	Name() string
	GetDC() int32    // 物攻
	GetMC() int32    // 魔攻
	GetAR() int32    // 物防
	GetMR() int32    // 魔防
	GetMaxHP() int32 // 最大血量
	GetCurHP() int32 // 当前血量
	Dodged()         // 闪避
	Critical()       // 暴击
	Reflect()        // 反伤
	AttackId(uid int64)
	Dead() bool
	SetLocation(l *pb.Point)
	SetCell(cell *Cell)
	GetCell() *Cell
	GetDirection() pb.Direction
	SettDirection(dir pb.Direction)
	SetDesk(desk *Desk)
	Desk() *Desk
	HaveBuff(pretype constant.BuffType) bool
	CfgID() int32
	Owner() IMapObject
	GetBuff(pretype constant.BuffType) *BuffInfo
	GetObAction() ObjectAction
	SetObAction(state ObjectAction)
	Level() int32
	Job() constant.MowJob
	IsWarrior() bool
	IsRobot() bool
	ServerId() int32
	GuildId() int64
	IsSameGuild(guildId int64, obj IMapObject) bool

	AttackAttributes() pb.Attribute // 攻击属性 怪物:根据配表确定物攻/魔攻， 玩家:根据职业/技能确定物攻/魔攻
	Attacked(attacker IMapObject, skillId int32, skillAdvanceLevel int32, power int32, bReflect bool)
	Revive(skill *PlayerSkill, skillAdvanceLevel, reviveHp int32)
	CleanDebuff()
	AddBuff(id int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32)
	AddBuff2(buffId int64, id int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32)
	Target() IMapObject
	ClearTarget()
	Friends() []IMapObject
	SkillBuffTrigger(skill *PlayerSkill)
	ChangeHp(hp int32, obj IMapObject) int32
	RefreshStats()
	//CanAddBuff(buffId int32) bool
	Die(attacker IMapObject)
	RefreshStatsAsync()
	DelayAction(target IMapObject, skillId int32, direction pb.Direction, location global.PointI)
	GetSceneOfficialExp() int64
}

type HpDisplay struct {
	// 血量变化相关
	currentHP           int32
	displayCrit         bool // 暴击
	displayMiss         bool // 闪避
	displayReflect      bool // 反伤
	displayHP           int32
	attackId            int64
	displaySkillId      int32
	displaySuckingBlood bool
	displaymulAttact    int32 // 多段
}

// 称号
type TitleInfo struct {
	cfgId          int32
	star           int32
	activationTime time.Time
}

// 地图实例基类
type MapObject struct {
	IMapObject
	*HpDisplay
	objectType ObjectType
	cfgId      int32 // 配置id
	uid        int64
	//account_name  string // 账号名
	name          string // 名字
	location      *pb.Point
	direction     pb.Direction
	cell          *Cell
	attackTime    time.Time     // 攻击时间
	moveTime      time.Time     // 移动时间
	dead          bool          // 是否死亡
	deadTime      time.Time     // 副本内死亡时间
	reviveTime    time.Time     // 副本内复活时间
	bornLocation  global.PointI // 树怪位置(用于刷怪，复活)
	spawnLocation global.PointI // 出生位置(用于追逐距离)
	reviveTimes   int           // 复活次数
	createTime    time.Time     // 创建时间
	delTime       time.Time     // 消失时间

	target IMapObject // 目标
	//block  sync.Mutex
	//buffs              map[int32]*BuffInfo // buff
	buffMgr            *BuffManager
	slock              sync.Mutex
	stats              *Stats     // 数据统计
	desk               *Desk      // 当前桌
	owner              IMapObject // 主人
	linkedSkill        int32
	obAction           ObjectAction    // 当前动作
	level              int32           // 等级
	skills             []*PlayerSkill  // 技能
	job                constant.MowJob // 职业
	serverID           int32           // 实例所在服务器
	bRobot             bool            // 是否机器人
	robotType          RobotType       // 机器人类型
	titles             *pb.DBTitle     // 称号列表
	equipTitleId       int32           // 已装备的称号
	activateTitleBonds []int32         // 激活称号共鸣
}

func NewMapObject(t ObjectType, uid int64) *MapObject {
	mapObject := &MapObject{
		objectType:    t,
		uid:           uid,
		location:      &pb.Point{},
		attackTime:    time.Now(),
		moveTime:      time.Now(),
		dead:          false,
		direction:     pb.Direction_UP,
		bornLocation:  global.PointI{},
		spawnLocation: global.PointI{},
		HpDisplay:     &HpDisplay{},
		//buffs:       map[int32]*BuffInfo{},
		buffMgr:     new(BuffManager),
		cell:        &Cell{},
		reviveTimes: 0,
		stats:       &Stats{Attrs: make(map[pb.Attribute]int32)},
		createTime:  time.Now(),
		linkedSkill: 0,
		obAction:    ObjectActionStand,
		skills:      []*PlayerSkill{},
		titles:      &pb.DBTitle{},
	}
	mapObject.buffMgr.init(mapObject)
	return mapObject
}

func CreateBuffObject(cfgId int32, obType ObjectType, owner IMapObject, addMilliSec int32, linkedSkill int32, exStats *Stats) *MapObject {
	if buffCfg, err := global.GetBuffCfg(cfgId); err == nil {
		uid := snowflake.GenID()
		mob := NewMapObject(obType, uid)
		mob.cfgId = buffCfg.ID
		mob.owner = owner
		mob.linkedSkill = linkedSkill

		if exStats != nil {
			for k, v := range exStats.Attrs {
				mob.stats.setStats(k, v)
				mob.stats.setStats(pb.Attribute_HEALTH, mob.stat(pb.Attribute_MAX_HEALTH))
			}
			mob.currentHP = exStats.getStats(pb.Attribute_MAX_HEALTH)
		}

		mob.delTime = time.Now().Add(time.Duration(addMilliSec) * time.Millisecond)

		return mob
	}
	return nil
}

func (o *MapObject) CanAttack() bool {
	if o.dead {
		return false
	}

	if o.desk == nil {
		return false
	}

	if time.Now().Before(o.attackTime) {
		return false
	}

	if time.Now().Before(o.moveTime) {
		return false
	}
	if o.getBuffByType(constant.BuffFreeze) != nil {
		return false
	}

	return true
}

func (o *MapObject) CanAttackTarget(target IMapObject) bool {
	if target == nil {
		return false
	}

	if o.desk == nil || target.Desk() == nil {
		return false
	}

	if target.Dead() {
		// logger.Debugf(" 无法攻击  目标已死 uid=%d ", o.uid)
		return false
	}

	if o.uid == target.Uid() {
		return false
	}

	if o.dead && o.objectType != ObjectTypeDecoy {
		return false
	}

	if target.ObjectType() == ObjectTypeBuff {
		return false
	}

	if o.Owner() != nil && o.Owner().Uid() == target.Uid() {
		return false
	}

	if o.getBuffByType(constant.BuffFreeze) != nil {
		// logger.Debugf(" 无法攻击 身中眩晕 uid=%d ", o.uid)
		return false
	}

	if target.ObjectType() == ObjectTypeMonster {

		if m, err := o.desk.monsterWithUid(target.Uid()); err == nil {
			if m.cfg.MType == MTypeGathering {
				return false
			}
		}
		return o.objectType != ObjectTypeMonster
	}

	if (target.ObjectType() == ObjectTypeDecoy || target.ObjectType() == ObjectTypePartner) && target.Owner() != nil {
		if target.Owner().Uid() == o.uid {
			return false
		}
		if o.Owner() != nil && target.Owner().Uid() == o.owner.Uid() {
			return false
		}
		if o.Owner() != nil && target.Owner() != nil && target.Owner().Owner() != nil && target.Owner().Owner().Uid() == o.Owner().Uid() {
			return false
		}
	}

	if o.objectType == ObjectTypeDecoy && o.owner != nil {
		if o.owner.Uid() == target.Uid() {
			return false
		}
	}

	if target.HaveBuff(constant.BuffGodMode) {
		return false
	}

	if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if o.GuildId() != 0 && target.GuildId() == o.GuildId() {
			return false
		}
	}

	return true
}

func (o *MapObject) CanHelpTarget(target IMapObject) bool {
	if target == nil {
		return false
	}
	if o.desk == nil || target.Desk() == nil {
		return false
	}

	if o.dead {
		return false
	}

	if o.objectType == ObjectTypeDecoy || o.objectType == ObjectTypeMonster || target.ObjectType() == ObjectTypeMonster {
		return false
	}

	if target.HaveBuff(constant.BuffGodMode) {
		return false
	}

	bFriend := false
	for _, v := range o.Friends() {
		if v.Uid() == target.Uid() {
			bFriend = true
		}
	}

	return bFriend
}

//func (o *MapObject) CanAddBuff(buffId int32) bool {
//cfg, err := global.GetBuffCfg(buffid)
//
//if err != nil {
//	logger.Warningf("获取buff配置失败 %d:%s buffid=%d", o.uid, o.name, buffid)
//	return false
//}
//
//if o.dead {
//	return false
//}
//
//if cfg.Mark == constant.BuffTypeChainSkill && o.desk == nil {
//	return false
//}
//
//if cfg.Type == constant.BuffTypeDebuff {
//	if o.HaveBuff(constant.BuffGodMode) {
//		return false
//	}
//
//	if o.HaveBuff(constant.BuffImmune) {
//		return false
//	}
//	// 免疫debuff
//	if immuneBuff := o.getBuffByType(constant.BuffNoNegativeEffect); immuneBuff != nil {
//		if global.RandIn(int(BaseRate)) < int(immuneBuff.Rate) {
//			logger.Debugf("免疫debuff  %d:%s  buffid=%d", o.uid, o.name, buffid)
//			return false
//		}
//	}
//}
//
//return true
//}

func (o *MapObject) SkillEffectType(skillCfg *global.SkillConfig, target IMapObject) bool {
	if skillCfg.Object == constant.SkillEffectTypeMonster && target.ObjectType() != ObjectTypeMonster {
		return false
	}
	if skillCfg.Object == constant.SkillEffectTypePlayer && target.ObjectType() == ObjectTypeMonster {
		return false
	}
	return true
}

// 能否移动
func (o *MapObject) CanMove() bool {
	if o.dead {
		return false
	}

	if o.desk == nil {
		return false
	}

	if time.Now().Before(o.moveTime) {
		return false
	}
	// if time.Now().Before(o.attackTime) {
	// 	return false
	// }

	if o.getBuffByType(constant.BuffFreeze) != nil {
		return false
	}

	if o.getBuffByType(constant.BuffStun) != nil {
		return false
	}
	return true
}

func (o *MapObject) Location() *pb.Point {
	return o.location
}

// 攻击范围
func (o *MapObject) inAttackRange(distance int) bool {
	return maputil.InRange(o.location, o.target.Location(), distance)
}

func (o *MapObject) stat(k pb.Attribute) int32 {
	return o.stats.getStats(k)
}

func (o *MapObject) getBuffByType(pretype constant.BuffType) *BuffInfo {
	//defer func() {
	//	if err := recover(); err != nil {
	//		logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
	//	}
	//}()
	//if o.block.TryLock() {
	//	defer o.block.Unlock()
	//	for _, v := range o.buffs {
	//		if v.Pretype == int32(pretype) {
	//			if time.Now().Before(v.OverTime) {
	//				return v
	//			}
	//			//logger.Debugf(" [%d:%s] buff存在但是已过期,类型 %d BuffID %d 当前时间 %d  过期时间 %d CD过期时间 %d", o.uid, o.name, pretype, v.Id, time.Now().Unix(), v.OverTime.Unix(), v.CDTime.Unix())
	//		}
	//	}
	//} else {
	//	//todo processBuff 锁block ，在processBuff 中造成伤害，会调用 die ， die 中会调用 getBuffByType ，造成死锁
	//	logger.Warningf("getBuffByType 锁未获取到 uid=%d 查找buff类型 %d", o.Uid(), pretype)
	//}

	//return nil
	return o.buffMgr.getBuffByType(pretype)
}

// client buff list
func (o *MapObject) getPBBuffList() []*pb.BuffInfo {
	//defer func() {
	//	if err := recover(); err != nil {
	//		logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
	//	}
	//}()
	//o.block.Lock()
	//defer o.block.Unlock()
	//
	//var buffs []*pb.BuffInfo
	//for _, v := range o.buffs {
	//	one := &pb.BuffInfo{
	//		BuffId:        v.Id,
	//		RemainingTime: int32(v.OverTime.UnixMilli() - v.StartTime.UnixMilli()),
	//		BuffStats:     []*pb.Stats{},
	//	}
	//
	//	for bk, bv := range v.BuffStats.Attrs {
	//		one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: bk, Num: bv})
	//	}
	//
	//	buffs = append(buffs, one)
	//}
	//
	//return buffs
	return o.buffMgr.getPBBuffList()
}

// 血量变化
func (o *MapObject) changeHp(hp int32, mapobj IMapObject) int32 {
	var realHp float64
	if hp > 0 {
		if o.currentHP == o.GetMaxHP() {
			return o.currentHP
		}
		realHp = float64(hp) * (1 + float64(o.stat(pb.Attribute_LIFE_RECOVERY_INCREASE_PERCENTAGE))/BaseRate)
		realHp = float64(realHp) * (1 - float64(o.stat(pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE))/BaseRate)
		if realHp <= 0 {
			return o.currentHP
		}
	} else {
		realHp = float64(hp)
	}

	// 同步世界boss血量
	if o.ObjectType() == ObjectTypeMonster {
		if cfg := global.GetMonsterCfg(o.cfgId); cfg != nil && cfg.MType == MTypeWorldBoss {
			wb := worldBossMgr.getBoss(o.serverID)
			if wb != nil && mapobj != nil {
				// logger.Infof("changeHp 世界boss 当前血量变化 uid=%d  name=%s  change=%f, stack:\n%s", o.Uid(), o.Name(), realHp, debug.Stack())
				wb.changeHp(int32(realHp), mapobj)
			}
		}
	}

	o.currentHP += int32(realHp)
	if o.currentHP <= 0 {
		o.currentHP = 0
		o.processHp()
		o.die(mapobj)
		return 0
	}
	if o.currentHP > o.GetMaxHP() {
		o.currentHP = o.GetMaxHP()
	}

	o.stats.changeStats(pb.Attribute_HEALTH, o.currentHP)
	return o.currentHP
}

func (o *MapObject) processHp() {
	if o.dead {
		return
	}
	if o.displayCrit || o.displayMiss || o.displayReflect || o.displayHP != o.currentHP {
		change := o.currentHP - o.displayHP

		o.displayHP = o.currentHP

		// 判断是否是玩家并且目标是世界boss
		isPlayer := o.objectType == ObjectTypePlayer
		isWorldBoss := false
		if o.target != nil && o.target.ObjectType() == ObjectTypeMonster {
			if cfg := global.GetMonsterCfg(o.target.CfgID()); cfg != nil && cfg.MType == MTypeWorldBoss {
				isWorldBoss = true
			}
		}

		if isPlayer && isWorldBoss {
			o.pushMsg(protocol.RouteObjectHpChangePush, &pb.ObjectHpChangePush{
				TargetId:   o.uid,
				AttackerId: o.attackId,
				ChangeHp:   change,
				SkillId:    o.displaySkillId,
				Critical:   o.displayCrit,
				Miss:       o.displayMiss,
				Reflect:    o.displayReflect,
				MaxHp:      o.GetMaxHP(),
				MulAttact:  o.displaymulAttact,
			})
		} else if o.desk != nil {
			o.desk.Broadcast(protocol.RouteObjectHpChangePush, &pb.ObjectHpChangePush{
				TargetId:   o.uid,
				AttackerId: o.attackId,
				ChangeHp:   change,
				SkillId:    o.displaySkillId,
				Critical:   o.displayCrit,
				Miss:       o.displayMiss,
				Reflect:    o.displayReflect,
				MaxHp:      o.GetMaxHP(),
				MulAttact:  o.displaymulAttact,
			})
		}

		o.displaySkillId = 0
		o.displayCrit = false
		o.displayMiss = false
		o.displayReflect = false
		o.attackId = 0
		o.displaymulAttact = 0
		o.displaySuckingBlood = false
	}
}

func (o *MapObject) processBuff() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	o.buffMgr.ProcessBuff()
	//o.block.Lock()
	//defer o.block.Unlock()
	//
	//if o.dead {
	//	return
	//}
	//now := time.Now()
	//var delbufIds []int32
	//var bRefreshStats bool
	//for k, v := range o.buffs {
	//	// CD 和 持续时间都到， 删除buff
	//	if now.After(v.OverTime) && now.After(v.CDTime) && v.Mark != constant.BuffMarkPassive {
	//		delbufIds = append(delbufIds, k)
	//		continue
	//	}
	//	// 持续时间到 , 清除效果保留至CD结束
	//	if now.After(v.OverTime) && v.Mark != constant.BuffMarkPassive {
	//		if o.buffs[k].Status == int32(constant.BuffStatusStart) {
	//			o.buffs[k].BuffStats.clear()
	//			bRefreshStats = true
	//			o.buffs[k].Status = int32(constant.BuffStatusStop)
	//		}
	//		continue
	//	}
	//
	//	// 嘲讽状态 释放者退出 死亡 超过6米 均接触嘲讽状态
	//	if v.BuffType == constant.BuffTaunt {
	//		if v.Owner == nil || v.Owner.Dead() || maputil.Distance(o.location, v.Owner.Location()) > 6 {
	//			delbufIds = append(delbufIds, k)
	//			continue
	//		}
	//	}
	//
	//	if v.TickFrequency == 0 {
	//		continue
	//	}
	//
	//	if v.Rate != 0 && global.RandIn(int(BaseRate)) >= int(v.Rate) {
	//		continue
	//	}
	//
	//	// 间隔未到
	//	if !now.After(v.TickTime) {
	//		continue
	//	}
	//
	//	var damage, recover float64
	//	switch constant.BuffType(v.Pretype) {
	//	case constant.BuffFixedHeal, constant.BuffPercentageHeal:
	//		for i, e := range v.BuffStats.Attrs {
	//			if i == pb.Attribute_HEALTH {
	//				recover += float64(e)
	//			}
	//			if i == pb.Attribute_HEALTH_PERCENTAGE {
	//				recover += float64(o.GetMaxHP()) * float64(e) / BaseRate
	//			}
	//		}
	//	case constant.BuffFixedDamage, constant.BuffPercentageDamage, constant.BuffFixedDamage2:
	//		for i, e := range v.BuffStats.Attrs {
	//			if i == pb.Attribute_HEALTH {
	//				damage += float64(e)
	//			}
	//			if i == pb.Attribute_HEALTH_PERCENTAGE {
	//				damage += float64(o.GetMaxHP()) * float64(e) / BaseRate
	//			}
	//
	//			//出血，中毒，灭世，三个技能，输出类型的最高伤害不超过当前攻击力3倍；掉血类型的最高掉血量不超过当前释放者血量上限的3倍
	//			if damage > float64(v.Owner.GetMaxHP()*3) {
	//				if v.Owner != nil {
	//					damage = float64(v.Owner.GetMaxHP() * 3)
	//				} else {
	//					damage = 1
	//				}
	//			}
	//		}
	//	case constant.BuffFixedHealOnece:
	//		for i, e := range v.BuffStats.Attrs {
	//			if i == pb.Attribute_HEALTH {
	//				recover += float64(e)
	//			}
	//			if i == pb.Attribute_HEALTH_PERCENTAGE {
	//				recover += float64(o.GetMaxHP()) * float64(e) / BaseRate
	//			}
	//		}
	//		delbufIds = append(delbufIds, k)
	//	case constant.BuffAcidRain:
	//		damageRain := float64(o.GetMC()) * float64(v.BuffStats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE)) / BaseRate
	//		damageRain += damageRain * (1 + float64(v.BuffStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate)
	//		o.magicAcidRain(damageRain, v.SkillID)
	//
	//	case constant.BuffCurHpPercentageDamage:
	//		damage += float64(o.currentHP) * float64(v.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
	//		if damage > float64(v.Owner.GetMaxHP()*3) {
	//			if v.Owner != nil {
	//				damage = float64(v.Owner.GetMaxHP() * 3)
	//			} else {
	//				damage = 1
	//			}
	//		}
	//	}
	//
	//	// 挫折意志 减伤效果只针对sype = 4 或 5
	//	var damage_percentage_decrease, flat_damage_decrease int32
	//	if v.SkillID > 0 {
	//		if skillCfg, err := global.GetSkillCfg(v.SkillID); err == nil {
	//			if skillCfg.SType == constant.SkillTypeCareerSpecial {
	//				damage_percentage_decrease += o.stat(pb.Attribute_CONTROL_DAMAGE_REDUCTION_PERCENTAGE)
	//				flat_damage_decrease += o.stat(pb.Attribute_CONTROL_DAMAGE_REDUCTION_FIXED_VALUE)
	//			}
	//		}
	//	}
	//	damage -= damage * float64(damage_percentage_decrease) / BaseRate
	//	damage -= float64(flat_damage_decrease)
	//
	//	if damage > 0 && o.currentHP > 1 {
	//		damage = math.Min(damage, float64(o.currentHP-1))
	//		if o.ChangeHp(-int32(damage), v.Owner) == 0 {
	//			return
	//		}
	//	}
	//	if recover > 0 {
	//		o.attackId = v.Owner.Uid()
	//		if o.ChangeHp(int32(recover), v.Owner) == 0 {
	//			return
	//		}
	//	}
	//	o.buffs[k].TickTime = time.Now().Add(time.Duration(v.TickFrequency) * time.Millisecond)
	//	//logger.Debugf("  Buff生效 buffid = %d , uid:name %d:%s 伤害 %f  治疗 %f  到期时间 %v ", v.Id, o.uid, o.name, damage, recover, v.OverTime)
	//}
	//
	//if len(delbufIds) > 0 {
	//	for _, v := range delbufIds {
	//		o.delBuff(int32(v))
	//	}
	//	bRefreshStats = true
	//}
	//
	//if bRefreshStats {
	//	o.RefreshStatsAsync()
	//}
}

func (o *MapObject) addBuff(buffIID int64, cfgId int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	o.buffMgr.addBuff(buffIID, cfgId, owner, addSec, exStats, addRate, skillID)
	//o.block.Lock()
	//defer o.block.Unlock()
	//
	//cfg, err := global.GetBuffCfg(cfgId)
	//if err != nil {
	//	return
	//}
	//now := time.Now()
	//
	//if buff, ok := o.buffs[cfgId]; ok {
	//	if now.Before(buff.CDTime) {
	//		return
	//	}
	//}
	//
	//// 同类型buff 效果强的替换掉效果弱的,弱的不生效
	//bRcover := false
	//recoverId := 0
	//
	//var restoreHP float64
	//for i, old_buff := range o.buffs {
	//	if old_buff.Pretype == cfg.Pretype {
	//		for old_i, old_v := range old_buff.BuffStats.Attrs {
	//			for _, cfg_v := range cfg.Value {
	//				if cfg_v.Key == int32(old_i) {
	//					if old_v < cfg_v.Value {
	//						bRcover = true
	//						recoverId = int(i)
	//						break
	//					} else {
	//						return
	//					}
	//				}
	//			}
	//			if exStats != nil {
	//				for ex_k, ex_v := range exStats.Attrs {
	//					if ex_k == old_i {
	//						if old_v < ex_v {
	//							bRcover = true
	//							recoverId = int(i)
	//							break
	//						} else {
	//							return
	//						}
	//					}
	//				}
	//			}
	//		}
	//		break
	//	}
	//}
	//if bRcover && recoverId != 0 {
	//	o.delBuff(int32(recoverId))
	//}
	//
	//// 免疫时间 >= buff时间
	//if cfg.LastTime+addSec <= 0 {
	//	return
	//}
	//
	//buff := &BuffInfo{
	//	Id:            cfgId,
	//	BuffType:      cfg.Type,
	//	Mark:          cfg.Mark,
	//	StartTime:     now,
	//	TickTime:      now.Add(time.Duration(cfg.EffectTime) * time.Millisecond),
	//	OverTime:      now.Add(time.Duration(cfg.LastTime+addSec+100) * time.Millisecond),
	//	CDTime:        now.Add(time.Duration(cfg.CDTime) * time.Millisecond),
	//	Pretype:       cfg.Pretype,
	//	Owner:         owner,
	//	Status:        int32(constant.BuffStatusStart),
	//	BuffStats:     Stats{Attrs: make(map[pb.Attribute]int32)},
	//	TickFrequency: cfg.EffectTime,
	//	Rate:          cfg.Rate + addRate,
	//	SkillID:       skillID,
	//	BuffId:        buffIID,
	//}
	//
	//switch constant.BuffType(buff.Pretype) {
	//case constant.BuffDecoy:
	//	return
	//case constant.BuffPercentageDefenseDecrease,
	//	constant.BuffPercentageAttackDecrease,
	//	constant.BuffPercentageDefenseHealthDecrease,
	//	constant.BuffPercentageCritChanceDecrease,
	//	constant.BuffPercentageAccuracyDecrease,
	//	constant.BuffPercentageAttackSpeedDecrease,
	//	constant.BuffPercentageMovementSpeedDecrease:
	//	for _, v := range cfg.Value {
	//		if v.Key == int32(pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE) {
	//			buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
	//		} else {
	//			buff.BuffStats.setStats(pb.Attribute(v.Key), -v.Value)
	//		}
	//	}
	//	if exStats != nil {
	//		for k, v := range exStats.Attrs {
	//			if k == pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE {
	//				buff.BuffStats.setStats(k, v)
	//			} else {
	//				buff.BuffStats.setStats(k, -v)
	//			}
	//		}
	//	}
	//
	//case constant.BuffShield:
	//	shieldHp := float64(o.GetMaxHP()) * float64(exStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
	//	shieldHp = math.Max(shieldHp, 1)
	//	buff.BuffStats.setStats(pb.Attribute_HEALTH, int32(shieldHp))
	//
	//case constant.BuffTaunt:
	//	o.target = owner // 嘲讽强制改变攻击目标
	//case constant.BuffCommonDebuff: // 通用减益
	//	for _, v := range cfg.Value {
	//		buff.BuffStats.setStats(pb.Attribute(v.Key), -v.Value)
	//	}
	//	if exStats != nil {
	//		for k, v := range exStats.Attrs {
	//			buff.BuffStats.setStats(k, -v)
	//		}
	//	}
	//
	//case constant.BuffCommonBuff: // 通用增益
	//	for _, v := range cfg.Value {
	//		buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
	//	}
	//	if exStats != nil {
	//		for k, v := range exStats.Attrs {
	//			buff.BuffStats.setStats(k, v)
	//		}
	//	}
	//case constant.BuffPercentageMaxHealthIncrease2:
	//	if maxhpp, ok := exStats.Attrs[pb.Attribute_HEALTH_MAX_PERCENTAGE]; ok {
	//		restoreHP = float64(o.GetMaxHP()) * float64(exStats.getStats(pb.Attribute_HEALTH_MAX_PERCENTAGE)) / BaseRate
	//		o.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(o.GetMaxHP())*float64(maxhpp)/BaseRate))
	//		buff.Shield = int32(restoreHP)
	//	}
	//	for k, v := range exStats.Attrs {
	//		buff.BuffStats.setStats(k, v)
	//	}
	//
	//	logger.Debugf("[%d:%s]添加 强脉buff 剩余护盾 %d  最大生命值 %d 当前生命值 %d", o.uid, o.name, buff.Shield, o.GetMaxHP(), o.currentHP)
	//
	//default:
	//	for _, v := range cfg.Value {
	//		buff.BuffStats.setStats(pb.Attribute(v.Key), v.Value)
	//	}
	//	//logger.Debugf("[%s-%v] 增加buff: %+v", o.name, o.uid, cfg.Value)
	//	if exStats != nil {
	//		for k, v := range exStats.Attrs {
	//			buff.BuffStats.setStats(k, v)
	//		}
	//	}
	//}
	//
	//o.buffs[cfgId] = buff
	//
	//packet := &pb.AddBuffPush{
	//	TargetId: o.Uid(),
	//	Buff: &pb.BuffInfo{
	//		BuffId:        cfgId,
	//		RemainingTime: cfg.LastTime + addSec,
	//		BuffStats:     []*pb.Stats{},
	//		OwnerUid:      owner.Uid(),
	//	},
	//}
	//for bk, bv := range buff.BuffStats.Attrs {
	//	packet.Buff.BuffStats = append(packet.Buff.BuffStats, &pb.Stats{Attr: bk, Num: bv})
	//}
	//
	//if o.desk != nil {
	//	o.desk.Broadcast(protocol.RouteAddBuffPush, packet)
	//} else {
	//	o.pushMsg(protocol.RouteAddBuffPush, packet)
	//}
	//
	//if restoreHP > 0 {
	//	o.ChangeHp(int32(restoreHP), owner)
	//}
	//if buff.BuffId != 0 {
	//	logger.Debugf("++++[%s]++++ 增加buff,  BuffId=%d BuffCfgId=%d skillId= %d  addSec= %d , stats = %v  maxhp = %d  currHp = %d 结束时间 %v",
	//		o.name, buff.BuffId, buff.Id, skillID, addSec, buff.BuffStats.Attrs, o.GetMaxHP(), o.currentHP, buff.OverTime.Unix())
	//}
	//// logger.Debugf("++++[%s]++++ 增加buff,  BuffId=%d  skillId= %d  addSec= %d , stats = %v  maxhp = %d  currHp = %d 结束时间 %v",
	//// 	o.name, buff.Id, skillID, addSec, buff.BuffStats.Attrs, o.GetMaxHP(), o.currentHP, buff.OverTime.Unix())
}

func (o *MapObject) delBuff(id int32) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	o.buffMgr.delBuff(id)
	//// 装备升星几率的buf如果没有过期则不清理
	//var buffId int64
	//if buff, ok := o.buffs[id]; ok {
	//	switch buff.Pretype {
	//	case constant.BuffPercentEquipStarIncrease:
	//		if time.Now().Before(buff.OverTime) {
	//			return
	//		}
	//	case constant.BuffPercentageMaxHealthIncrease2:
	//		if buff.Shield > 0 {
	//			o.ChangeHp(-buff.Shield, o)
	//
	//			logger.Debugf("[%d:%s]移除 强脉buff 剩余护盾 %d  最大生命值 %d  当前生命值 %d", o.uid, o.name, buff.Shield, o.GetMaxHP(), o.currentHP)
	//
	//			buff.Shield = 0
	//		}
	//	}
	//	// if buff.Pretype == constant.BuffPercentEquipStarIncrease {
	//	// 	if time.Now().Before(buff.OverTime) {
	//	// 		return
	//	// 	}
	//	// }
	//	buffId = buff.BuffId
	//}
	//
	//delete(o.buffs, id)
	//
	//o.RefreshStatsAsync()
	//
	//packet := &pb.RemoveBuffPush{BuffId: id, TargetId: o.Uid()}
	//if o.desk != nil {
	//	o.desk.Broadcast(protocol.RouteRemoveBuffPush, packet)
	//} else {
	//	o.pushMsg(protocol.RouteRemoveBuffPush, packet)
	//}
	//
	//if buffId != 0 {
	//	logger.Debugf("----[%s]---- 删除buff ,BuffId=%d BuffCfgId=%d ", o.name, buffId, id)
	//}
	////logger.Debugf("----[%s]---- 删除buff ,BuffId=%d  ", o.name, id)
}

// 清理游戏内产生的buff
func (o *MapObject) resetBuff() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	o.buffMgr.resetBuff()
	//if o.block.TryLock() {
	//	defer o.block.Unlock()
	//	whileBuff := make([]int32, 0)
	//	if player, ok := defaultManager.player(o.uid); ok {
	//		if len(player.pet) != 0 {
	//			if player.pet[0].PetCfgId != 0 {
	//				// 查询宠物buff
	//				petInfo := global.GetPetCfg(player.pet[0].PetCfgId)
	//				for _, v := range petInfo.Talent {
	//					whileBuff = append(whileBuff, v.ID)
	//				}
	//			}
	//		}
	//	}
	//	var delbuffs []int32
	//	for k := range o.buffs {
	//		if o.buffs[k].Pretype == constant.BuffPercentEquipStarIncrease {
	//			continue
	//		}
	//
	//		if o.buffs[k].Mark == constant.BuffMarkPassive {
	//			continue
	//		}
	//		// 在白名单的buff也要跳过
	//		in := false
	//		for _, v := range whileBuff {
	//			if v == k {
	//				in = true
	//				break
	//			}
	//		}
	//		if in {
	//			continue
	//		}
	//		delbuffs = append(delbuffs, k)
	//	}
	//
	//	for i := 0; i < len(delbuffs); i++ {
	//		o.delBuff(delbuffs[i])
	//	}
	//}
	//
	//for i := range o.skills {
	//	o.skills[i].Cooldown = time.Now()
	//}
}

func (o *MapObject) pushMsg(route string, msg interface{}) {
	if o.objectType == ObjectTypePlayer {
		if player, ok := defaultManager.player(o.uid); ok {
			player.PushMsg(route, msg)
		}
	} else if o.objectType == ObjectTypePartner {
		if o.owner == nil {
			return
		}
		if player, ok := defaultManager.player(o.owner.Uid()); ok {
			player.PushMsg(route, msg)
		}
	}
}

func (o *MapObject) RefreshStatsAsync() {
	if o.desk != nil {
		o.desk.eventMgr.Publish(event.EventTypeRefreshStats, o.uid)
	} else {
		o.RefreshStats()
	}
}

// 房间中刷新属性
func (o *MapObject) RefreshStats() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	if o.desk == nil {
		if p, ok := playerMgr.GetPlayer(o.uid); ok {
			p.refreshStats()
		}
		return
	}
	// 主角和伙伴属性独立计算
	uid := o.uid
	if o.objectType == ObjectTypePartner {
		if o.owner == nil {
			logger.Debugf("查找伙伴的主人失败 uid:name %d:%s", o.uid, o.name)
			return
		}
		uid = o.owner.Uid()
		if player, err := o.desk.playerWithId(uid); err == nil {
			player.refreshStats()
		}
		if robot, err := o.desk.robotWithId(uid); err == nil {
			robot.refreshStats()
		}
		return
	} else if o.objectType == ObjectTypePlayer {
		if player, err := o.desk.playerWithId(uid); err == nil {
			player.refreshStats()
		}
		return
	} else if o.objectType == ObjectTypeRobot {
		if r, err := o.desk.robotWithId(uid); err == nil {
			r.refreshStats()
		} else {
			logger.Debugf("机器人查找失败 uid %d name %s", o.uid, o.name)
		}
		return
	}

	// 怪物
	now := time.Now()
	o.stats.clear()
	baseStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	if o.objectType == ObjectTypeMonster {
		cfg := global.GetMonsterCfg(o.cfgId)
		if cfg == nil {
			logger.Debugf("怪物查找配置失败 cfgId %d", o.cfgId)
			return
		}
		// 基础属性
		for _, v := range cfg.BassProperty {
			baseStatus.setStats(pb.Attribute(v.Key), v.Value)
		}
		baseStatus.setStats(pb.Attribute(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE), int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	} else if o.objectType == ObjectTypeRobot && o.desk != nil {
		if r, err := o.desk.robotWithId(o.uid); err == nil {
			r.refreshStats()
			return
		} else {
			logger.Debugf("机器人查找失败 uid %d name %s", o.uid, o.name)
			return
		}
	}

	for bk, bv := range baseStatus.Attrs {
		o.stats.setStats(bk, bv)
	}

	o.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
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
		// 部分buff中的属性不叠加
		if buffCfg.Pretype == constant.BuffCurHpPercentageDamage ||
			buffCfg.Pretype == int32(constant.BuffPercentageDamage) ||
			buffCfg.Pretype == int32(constant.BuffFixedHealOnece) {
			return true
		}

		buff.BuffStats.Range(func(i pb.Attribute, e int32) bool {
			o.stats.setStats(i, e)
			return true
		})
		return true
	})

	//for _, buff := range o.buffs {
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
	//		o.stats.setStats(i, e)
	//	}
	//}

	// < 0  攻速移速 增益，其他属性减益 ; >0 攻速移速 减益，其他属性增益
	o.stats.changeStats(pb.Attribute_HEALTH, o.currentHP)

	o.stats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(o.GetDC())*(float64(o.stats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(o.GetMC())*(float64(o.stats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(o.GetAR())*(float64(o.stats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(o.GetMR())*(float64(o.stats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	o.stats.setStats(pb.Attribute_DODGE, int32(float64(o.stats.getStats(pb.Attribute_DODGE))*(float64(o.stats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_HIT, int32(float64(o.stats.getStats(pb.Attribute_HIT))*(float64(o.stats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(o.stats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(o.stats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_TOUGHNESS, int32(float64(o.stats.getStats(pb.Attribute_TOUGHNESS))*(float64(o.stats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	o.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(o.GetMaxHP())*(float64(o.stats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	// 攻速移速增加  基础*(1-加成)
	o.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(o.stats.getStats(pb.Attribute_ATTACK_SPEED))*(1-(float64(o.stats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate))))
	o.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(o.stats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(o.stats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	o.stats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(o.GetDC()))))
	o.stats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(o.GetMC()))))
	o.stats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(o.GetAR()))))
	o.stats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(o.GetMR()))))

	o.stats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_DODGE)))))
	o.stats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_HIT)))))
	o.stats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	o.stats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_TOUGHNESS)))))
	o.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_ATTACK_SPEED)))))
	o.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(o.stats.getStats(pb.Attribute_MOVE_SPEED)))))
	o.stats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(o.GetMaxHP()))))
}

// 死亡
func (o *MapObject) die(object IMapObject) {
	if o.desk == nil {
		logger.Errorf("%d 死亡异常 , 桌子已经被清理 uid %d", o.objectType, o.uid)
		return
	}

	if o.dead {
		return
	}

	if object != nil {
		object.ClearTarget()
		if object.ObjectType() == ObjectTypePartner {
			object.Owner().ClearTarget()
		}
	}

	// 替身死亡AOE
	if o.linkedSkill > 0 {
		cfg, err := global.GetSkillCfg(o.linkedSkill)
		if err == nil {
			switch cfg.SkillID {
			case constant.SkillTypeDecoy:
				o.DecoyExplosionEnd()
			case constant.OctopusSkill1:
				o.octopusSkill3End(o.linkedSkill)
			}
		}
		return
	}

	if o.objectType == ObjectTypePlayer {
		if p, err := o.desk.playerWithId(o.uid); err == nil && p != nil {
			p.die(object)

			//优先使用宠物复活,宠物只能复活主角
			if buff := p.getBuffByType(constant.BuffReviveOnce); buff != nil {
				hp := float64(p.GetMaxHP()) * float64(buff.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
				time.AfterFunc(500*time.Millisecond, func() {
					p.Revive(nil, 0, int32(hp))
				})
				p.delBuff(buff.Id)
				logger.Debugf("宠物复活 目标[%d:%s] 复活血量%d  总血量 %d", o.uid, o.name, int32(hp), p.GetMaxHP())
			} else if skill, caster := p.getSlotSkillByType(constant.SkillTypeReborn); skill != nil {
				if caster != nil && !caster.dead && caster.uid != o.uid {
					caster.lastGasp(o, skill.CfgId)
				}
			}
		}
		return
	} else if o.objectType == ObjectTypePartner {
		if owner, err := o.desk.playerWithId(o.Owner().Uid()); err == nil {
			if skill, caster := owner.getSlotSkillByType(constant.SkillTypeReborn); skill != nil {
				if caster != nil && !caster.dead && caster.uid != o.uid {
					caster.lastGasp(o, skill.CfgId)
				}
			}
		}
	} else if o.objectType == ObjectTypeMonster {
		if m, ok := o.desk.monsters[o.uid]; ok {
			m.die(object)
		}
		return
	} else if o.objectType == ObjectTypeRobot {
		if r, err := o.desk.robotWithId(o.uid); err == nil && r != nil {
			r.die(object)

			//优先使用宠物复活,宠物只能复活主角
			if buff := r.getBuffByType(constant.BuffReviveOnce); buff != nil {
				hp := float64(r.GetMaxHP()) * float64(buff.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
				time.AfterFunc(500*time.Millisecond, func() {
					r.Revive(nil, 0, int32(hp))
				})
				r.delBuff(buff.Id)
				logger.Debugf("宠物复活 目标[%d:%s] 复活血量%d  总血量 %d", o.uid, o.name, int32(hp), r.GetMaxHP())
			} else if skill, caster := r.getSlotSkillByType(constant.SkillTypeReborn); skill != nil {
				caster.lastGasp(o, skill.CfgId)
			}
		} else {
			logger.Errorf("机器人死亡异常 serverID %d uid %d", o.serverID, o.uid)
		}
		return
	}

	o.dead = true
	o.deadTime = time.Now()
	o.target = nil

	packet := pb.ObjectDiePush{ObjectId: o.uid, ReviveTip: false}
	o.desk.Broadcast(protocol.RouteObjectDiePush, &packet)

	//logger.Debugf(" ==死亡== OBType [%d],  uid %d : %s", o.objectType, o.uid, o.name)
}

func (o *MapObject) reset() {
	o.desk = nil
	o.target = nil
	o.dead = false

	o.stats.setStats(pb.Attribute_HEALTH, o.stats.getStats(pb.Attribute_MAX_HEALTH))

	o.currentHP = o.GetMaxHP()
	o.displayHP = o.currentHP
	o.reviveTimes = 0

	o.resetBuff()
}

func (o *MapObject) delayAction(target IMapObject, skillId int32, direction pb.Direction, location global.PointI) {
	cfg, err := global.GetSkillCfg(skillId)
	if err != nil {
		logger.Errorf(" 获取技能配置失败 uid: %d, skillId: %d", o.uid, skillId)
		return
	}

	switch cfg.SkillID {
	case constant.SkillTypeDecoy, constant.SkillTypeDecoyBoss:
		o.decoyEnd(skillId, direction)
	case constant.SkillTypeInstantaneousMove:
		o.instantaneousMoveEnd(skillId)
	case constant.SkillTypeSwordHeavenBreaker:
		if o.swordHeavenBreakerEnd(skillId) {
			attackCells := o.desk.mapInfo.getSquareCells(&location, cfg.Radius)
			o.attackLocation(attackCells, cfg)
		}
	case constant.SkillTypeSwiftStrike:
		maxDistance := o.swiftStrikeEnd(skillId, direction)
		if maxDistance > 0 {
			attackCells := o.desk.mapInfo.getLineCells(&o.cell.location, direction, maxDistance)
			o.attackLocation(attackCells, cfg)
		}

	case constant.SkillTypeReborn: //回光返照
		o.lastGaspEnd(target, skillId)
	case constant.SkillTypeMassPurification,
		constant.SkillTypeMassPurificationBoss: //群体净化术
		o.MassPurificationEnd(skillId)
	case constant.OctopusSkill1:
		o.octopusSkill1End(target, skillId)
	case constant.OctopusSkill4:
		o.octopusSkill3End(skillId)
	case constant.QiongqiSkill1:
		o.qiongqiSkill1End(skillId)
	case constant.HeshangSkill1:
		o.heshangSkill1End(skillId, direction, location)
	case constant.SkillTypeBrambleCage: //荆棘牢笼
		o.brambleCageEnd(skillId, location)
	case constant.SkillTypeGoddessDescent:
		o.goddessDescentEnd(skillId)
	case constant.SkillTypeSwordRain:
		o.swordRainEnd(skillId, location)

	default:
		if target == nil {
			return
		}
		o.delayAttack(target, skillId, 1)
	}
}

func (o *MapObject) delayAttack(target IMapObject, skillId int32, mulAttact int32) {
	if !o.CanAttackTarget(target) {
		return
	}

	skillAdvancedLevel := int32(0)
	skill := o.getSkill(skillId)
	if skill != nil {
		skillAdvancedLevel = skill.AdvancedLevel
	}
	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		return
	}

	// 命中-闪避
	hit := o.stats.getStats(pb.Attribute_HIT)
	if tstats := target.GetStats(); tstats != nil {
		dodge := target.GetStats().getStats(pb.Attribute_DODGE)
		if hit < dodge && global.RandIn(int(dodge)) > int(hit) {
			target.AttackId(o.uid)
			target.Dodged()
			return
		}
	}

	// 技能效果
	var power, basePower, skillPower, skill_final_flat_damage_increase, skill_final_damage_percentage_increase float64

	basePower = float64(o.stats.getStats(o.AttackAttributes()))

	power += basePower
	for _, v := range skillCfg.Effect {
		switch pb.Attribute(v.Key) {
		case pb.Attribute_PHYSICAL_ATTACK, pb.Attribute_MAGIC_ATTACK:
			skillPower += float64(v.Value)
		case pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE:
			if o.AttackAttributes() == pb.Attribute_PHYSICAL_ATTACK {
				skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
			}
		case pb.Attribute_MAGIC_ATTACK_PERCENTAGE:
			if o.AttackAttributes() == pb.Attribute_MAGIC_ATTACK {
				skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
			}

		case pb.Attribute_PVP_DAMAGE_PERCENTAGE:
			if target.ObjectType() == ObjectTypePlayer {
				skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
			}
		case pb.Attribute_PVE_DAMAGE_PERCENTAGE:
			if target.ObjectType() == ObjectTypeMonster {
				skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
			}
		case pb.Attribute_DMT_FINAL_INCREASE_FIXED:
			skill_final_flat_damage_increase += float64(v.Value)

		case pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE:
			skill_final_damage_percentage_increase += float64(v.Value)
		}
	}

	power += skillPower

	var damage_percentage_increase, flat_damage_increase int32 // 技能进阶中的增伤百分比，增伤固定值
	var ignore_defense_percentage float64                      //无视防御百分比

	// 技能进阶效果
	if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skillAdvancedLevel); err == nil {
		damage_percentage_increase += skillUpperlevel.DamagePer
	}

	//技能特殊效果
	if skillCfg.SpecialEffect.MulAttact > 0 && mulAttact > 1 {
		damage_percentage_increase += skillCfg.SpecialEffect.MulAttact
	}
	ignore_defense_percentage += float64(skillCfg.SpecialEffect.Realdam) / BaseRate
	// 攻击-防御
	if tstats := target.GetStats(); tstats != nil {
		ignore_defense_value := float64(tstats.getStats(maputil.DefenseAttributes(o.AttackAttributes()))) * (1 - ignore_defense_percentage)
		ignore_defense_value = math.Max(0, ignore_defense_value)
		power -= ignore_defense_value
		//power -= float64(tstats.getStats(maputil.DefenseAttributes(o.AttackAttributes())))
	}

	power += power * float64(damage_percentage_increase) / BaseRate
	power += float64(flat_damage_increase)
	power = math.Max(power, 0)

	// 最终增伤
	var final_damage_percentage_increase, final_flat_damage_increase int32
	final_flat_damage_increase += int32(skill_final_flat_damage_increase)
	final_damage_percentage_increase += int32(skill_final_damage_percentage_increase)

	if target.ObjectType() == ObjectTypeMonster {
		final_flat_damage_increase += o.stat(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)
		final_damage_percentage_increase += o.stat(pb.Attribute_PVE_DAMAGE_PERCENTAGE)
	} else {
		final_flat_damage_increase += o.stat(pb.Attribute_DMT_INCREASE_FIXED_PLAYER)
		final_damage_percentage_increase += o.stat(pb.Attribute_PVP_DAMAGE_PERCENTAGE)
	}
	final_damage_percentage_increase += o.stat(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE)
	final_flat_damage_increase += o.stat(pb.Attribute_DMT_FINAL_INCREASE_FIXED)

	//官职压制
	if o.GetSceneOfficialExp() >= 0 && target.GetSceneOfficialExp() >= 0 && CheckOfficialSuppression(o.GetSceneOfficialExp(), target.GetSceneOfficialExp(), global.OfficalListCfg) == constant.SuppressionApplied {
		final_damage_percentage_increase += int32(global.GlobalCfg.BaseValue.Official.AttributeNum)
		//logger.Debugf("官职压制 攻击者: %d:%s  官职经验 %d   受击者 %d:%s 官职经验 %d ", o.uid, o.name, o.GetSceneOfficialExp(), target.Uid(), target.Name(), target.GetSceneOfficialExp())
	}

	power += power * float64(final_damage_percentage_increase) / BaseRate
	power += float64(final_flat_damage_increase)

	if power <= 0 {
		power = float64(rand.Intn(2)) + 1
	}
	if len(skillCfg.Effect) == 0 && skillCfg.SType != constant.SkillTypeNormal {
		power = 0
	}

	//logger.Debugf(" [%d:%s]攻击 targetID: %d  目标位置: %v skillID: %d  Power:%f", o.uid, o.name, target.Uid(), target.Location(), skillId, power)

	target.Attacked(o, skillId, skillAdvancedLevel, int32(power), true)
}

func (o *MapObject) attackLocation(cells []*Cell, skillCfg *global.SkillConfig) {
	var targets []IMapObject
	for _, cell := range cells {
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

		// 技能距离
		if skillCfg.MaxDistance > 0 && maputil.Distance(o.location, t.Location()) > int(skillCfg.MaxDistance) {
			continue
		}

		o.delayAttack(t, skillCfg.UID, 1)
	}
}

func (o *MapObject) getMoveSpeed() int32 {
	if o.stats.getStats(pb.Attribute_MOVE_SPEED) <= 0 {
		o.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(o.job)).MoveSpeed)
		//return DefaultMoveSpeed
	}
	return o.stats.getStats(pb.Attribute_MOVE_SPEED)
}

func (o *MapObject) getAttackSpeed() int32 {
	if o.stats.getStats(pb.Attribute_ATTACK_SPEED) <= 0 {
		o.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(o.job)).AttackSpeed)
		//return DefaultAttackSpeed
	}
	return o.stats.getStats(pb.Attribute_ATTACK_SPEED)
}

func (o *MapObject) getSkill(skillId int32) *PlayerSkill {
	for _, v := range o.skills {
		if v.CfgId == skillId {
			return v
		}
	}
	return nil
}

func (o *MapObject) setAttackTime(delay int32) {
	delay = o.getAttackSpeed() - 1
	if delay <= 0 {
		delay = DefaultAttackSpeed
	}
	o.attackTime = time.Now().Add(time.Duration(delay) * time.Millisecond)
}

// 验证攻速移速
func (o *MapObject) thresholdDetect() bool {
	if o.getMoveSpeed() <= 1 || o.getAttackSpeed() <= 1 {
		o.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(o.job)).AttackSpeed)
		o.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(o.job)).MoveSpeed)
		logger.Errorf("攻击失败，攻击速度异常 %d:%s  攻速 %d 移速 %d", o.uid, o.name, o.stat(pb.Attribute_ATTACK_SPEED), o.stat(pb.Attribute_MOVE_SPEED))
		return false
	}
	return true
}

// Base Interface
func (o *MapObject) HaveBuff(pretype constant.BuffType) bool     { return o.getBuffByType(pretype) != nil }
func (o *MapObject) GetBuff(pretype constant.BuffType) *BuffInfo { return o.getBuffByType(pretype) }
func (o *MapObject) CfgID() int32                                { return o.cfgId }
func (o *MapObject) GetCell() *Cell                              { return o.cell }
func (o *MapObject) SetCell(cell *Cell)                          { o.cell = cell }
func (o *MapObject) SettDirection(dir pb.Direction)              { o.direction = dir }
func (o *MapObject) GetDirection() pb.Direction                  { return o.direction }
func (o *MapObject) SetLocation(l *pb.Point)                     { o.location = l }
func (o *MapObject) Reflect()                                    { o.displayReflect = true }
func (o *MapObject) SuckingBlood()                               { o.displaySuckingBlood = true }
func (o *MapObject) DisplaymulAttact(v int32)                    { o.displaymulAttact = v }
func (o *MapObject) AttackId(uid int64)                          { o.attackId = uid }
func (o *MapObject) Dead() bool                                  { return o.dead }
func (o *MapObject) SetDesk(desk *Desk)                          { o.desk = desk }
func (o *MapObject) Desk() *Desk                                 { return o.desk }
func (o *MapObject) Uid() int64                                  { return o.uid }
func (o *MapObject) Name() string                                { return o.name }
func (o *MapObject) ObjectType() ObjectType                      { return o.objectType }
func (o *MapObject) GetStats() *Stats                            { return o.stats }
func (o *MapObject) GetDC() int32                                { return o.stats.getStats(pb.Attribute_PHYSICAL_ATTACK) }
func (o *MapObject) GetMC() int32                                { return o.stats.getStats(pb.Attribute_MAGIC_ATTACK) }
func (o *MapObject) GetAR() int32                                { return o.stats.getStats(pb.Attribute_PHYSICAL_DEFENSE) }
func (o *MapObject) GetMR() int32                                { return o.stats.getStats(pb.Attribute_MAGIC_DEFENSE) }
func (o *MapObject) GetMaxHP() int32                             { return o.stats.getStats(pb.Attribute_MAX_HEALTH) }
func (o *MapObject) GetCurHP() int32                             { return o.stats.getStats(pb.Attribute_HEALTH) }
func (o *MapObject) Dodged()                                     { o.displayMiss = true }
func (o *MapObject) Critical()                                   { o.displayCrit = true }
func (o *MapObject) Owner() IMapObject                           { return o.owner }
func (o *MapObject) GetObAction() ObjectAction                   { return o.obAction }
func (o *MapObject) SetObAction(state ObjectAction)              { o.obAction = state }
func (o *MapObject) Level() int32                                { return o.level }
func (o *MapObject) Job() constant.MowJob                        { return o.job }
func (o *MapObject) ServerId() int32                             { return o.serverID }
func (o *MapObject) GetDodge() int32                             { return o.stats.getStats(pb.Attribute_DODGE) }
func (o *MapObject) GetHit() int32                               { return o.stats.getStats(pb.Attribute_HIT) }
func (o *MapObject) GetCrit() int32                              { return o.stats.getStats(pb.Attribute_CRITICAL_STRIKE) }
func (o *MapObject) GetTough() int32                             { return o.stats.getStats(pb.Attribute_TOUGHNESS) }
func (o *MapObject) Target() IMapObject                          { return o.target }
func (o *MapObject) ClearTarget()                                { o.target = nil }
func (o *MapObject) DelayAction(target IMapObject, skillId int32, direction pb.Direction, location global.PointI) {
	o.delayAction(target, skillId, direction, location)
}

func (o *MapObject) IsWarrior() bool {
	return o.job == constant.Warrior || o.job == constant.PartnerWarrior
}

func (o *MapObject) IsRobot() bool {
	return o.bRobot
}

func (o *MapObject) AttackAttributes() pb.Attribute {
	if o.job == constant.Warrior || o.job == constant.PartnerWarrior {
		return pb.Attribute_PHYSICAL_ATTACK
	} else if o.job == constant.Wizard || o.job == constant.PartnerWizard || o.job == constant.Ritualist || o.job == constant.PartnerRitualist {
		return pb.Attribute_MAGIC_ATTACK
	} else {
		if o.GetDC() >= o.GetMC() {
			return pb.Attribute_PHYSICAL_ATTACK
		} else {
			return pb.Attribute_MAGIC_ATTACK
		}
	}
}

// 战斗场景内实例的官职经验
func (o *MapObject) GetSceneOfficialExp() int64 {
	officialExp := int64(-1)
	if o.desk == nil {
		return officialExp
	}

	switch o.objectType {
	case ObjectTypePlayer:
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			officialExp = p.getItemCountInBag(constant.ItemTypeOfficialExp)
		}
	case ObjectTypeRobot:
		if r, err := o.desk.robotWithId(o.uid); err == nil {
			officialExp = r.getItemCountInBag(constant.ItemTypeOfficialExp)
		}
	case ObjectTypePartner:
		if o.owner != nil {
			if p, err := o.desk.playerWithId(o.owner.Uid()); err == nil {
				officialExp = p.getItemCountInBag(constant.ItemTypeOfficialExp)
			} else if r, err := o.desk.robotWithId(o.owner.Uid()); err == nil {
				officialExp = r.getItemCountInBag(constant.ItemTypeOfficialExp)
			}
		}
	}
	return officialExp
}

// 被攻击
func (o *MapObject) Attacked(attacker IMapObject, skillId int32, skillAdvanceLevel int32, power int32, bReflect bool) {
	if o.dead || o.desk == nil || o.uid == 0 {
		logger.Warningf(" ===异常===  [%v] 被攻击失败 dead:%t  muid: %d", o.objectType, o.dead, o.uid)
		return
	}

	skillCfg, err := global.GetSkillCfg(skillId)
	if err != nil {
		return
	}

	// 斩杀先判断
	if skillCfg.SkillID == constant.SkillTypeExecute {
		canSlay := true
		if o.ObjectType() == ObjectTypeMonster {
			if m, err := o.desk.monsterWithUid(o.uid); err == nil {
				if cfg := global.GetMonsterCfg(m.cfgId); cfg != nil && (cfg.MType == MTypeWorldBoss || cfg.MType == MTypeWarBoss) {
					canSlay = false
				}
			}
		}

		if canSlay && float64(o.currentHP)/float64(o.GetMaxHP()) <= float64(skillCfg.SpecialEffect.Slay)/BaseRate {
			o.ChangeHp(-o.currentHP, attacker)
			return
		}
	}

	if power > 0 {

		//元素增伤 世界boss，一战到底，战场
		if o.desk.roomCfg.Roomtype == constant.WorldBossMap || o.desk.roomCfg.Roomtype == constant.FightEndMap || o.desk.roomCfg.Roomtype == constant.PVEWarMap {
			var final_damage_element_percentage_increase int32
			mainBase := global.GlobalCfg.Element.MainElement.AttributeID
			//攻击者的最高元素属性
			main := attacker.GetStats().getStats(pb.Attribute_Main_EMPOWER)
			if main > 0 && main <= 5 {
				if constant.ElementSlot_attribute[pb.ElementSlot(main)] > 0 {
					main = int32(constant.ElementSlot_attribute[pb.ElementSlot(main)])
				}
			} else {
				main = mainBase
			}
			//被攻击者的最高元素属性索引
			tarmainValue := global.GlobalCfg.Element.MainElement.AttributeID
			tarValue := o.stats.getStats(pb.Attribute(mainBase))
			for j := pb.Attribute_GOLD_EMPOWER; j <= pb.Attribute_EARTH_EMPOWER; j++ {

				if tarValue < o.stats.getStats(pb.Attribute(j)) {
					tarmainValue = int32(j)
					tarValue = o.stats.getStats(pb.Attribute(j))
				}
				if attacker.GetStats().getStats(pb.Attribute(j)) > o.stats.getStats(constant.Attribute_emphasis[j]) {
					final_damage_element_percentage_increase = (attacker.GetStats().getStats(pb.Attribute(j)) - o.stats.getStats(constant.Attribute_emphasis[j])) * 100 / global.GlobalCfg.Element.Battle.Constant
				}
			}

			//元素克制
			if constant.Attribute_know_attribute[pb.Attribute(main)] == pb.Attribute(tarmainValue) {
				final_damage_element_percentage_increase += global.GlobalCfg.Element.Battle.RestraintRatio * 100

				//logger.Debugf("[%d:%s]被攻击， 属性克制最终元素增伤  %d", o.uid, o.name, final_damage_element_percentage_increase)
			}
			power += int32(float64(power) * float64(final_damage_element_percentage_increase) / BaseRate)
		}

		//暴击
		critical := attacker.GetStats().getStats(pb.Attribute_CRITICAL_STRIKE)
		tenacity := o.stats.getStats(pb.Attribute_TOUGHNESS)
		criticalRate := (math.Max(0, float64(critical-tenacity))*global.GlobalCfg.BaseValue.Crit.K + float64(global.GlobalCfg.BaseValue.Crit.PreValue))

		var damage, finalDamage float64
		damage = float64(power)
		if global.RandIn(100) < int(criticalRate) {
			criticalDamageRate := float64(attacker.GetStats().getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)) / BaseRate
			damage = float64(power) * (1 + criticalDamageRate)
			o.Critical()
		}

		// 减伤,包含buff中的减伤
		damage -= float64(o.stat(pb.Attribute_DMT_FINAL_REDUCE_FIXED))
		if attacker.ObjectType() == ObjectTypeMonster {
			damage -= float64(o.stat(pb.Attribute_DMT_REDUCE_FIXED_MONSTER))
			damage -= damage * float64(o.stat(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER)) / BaseRate
		} else {
			damage -= float64(o.stat(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))
			damage -= damage * float64(o.stat(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER)) / BaseRate
		}
		damage -= damage * float64(o.stat(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)) / BaseRate

		damage = math.Max(damage, float64(global.RandIn(2)+1))

		// 护盾
		if buff := o.getBuffByType(constant.BuffShield); buff != nil {
			if buff.BuffStats.getStats(pb.Attribute_HEALTH) > int32(damage) {
				buff.BuffStats.setStats(pb.Attribute_HEALTH, -int32(damage))
				damage = 0
			} else {
				damage = damage - float64(buff.BuffStats.getStats(pb.Attribute_HEALTH))
				o.delBuff(buff.Id)
			}
		}
		// 强脉
		if buff := o.getBuffByType(constant.BuffPercentageMaxHealthIncrease2); buff != nil {
			if buff.Shield > 0 {
				buff.Shield -= int32(damage)
				//logger.Debugf("[%d:%s] 强脉buff 护盾扣除 %f  剩余护盾 %d", o.uid, o.name, damage, buff.Shield)
			}
		}
		// 如果有神圣天赋，豁免超过生命值百分比的伤害
		if buff := o.getBuffByType(constant.BuffDamageLimit); buff != nil {
			if damage > float64(o.GetMaxHP())*float64(buff.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate {
				damage = float64(o.GetMaxHP()) * float64(buff.BuffStats.getStats(pb.Attribute_HEALTH_PERCENTAGE)) / BaseRate
			}
		}

		finalDamage = damage

		if buff := o.getBuffByType(constant.BuffImmune); buff != nil {
			finalDamage = 0
		}

		// 吸血&技能中的伤害转移
		suckblood := finalDamage * float64(attacker.GetStats().getStats(pb.Attribute_SUCKING_BLOOD_PERCENTAGE)) / float64(BaseRate)
		suckblood += finalDamage * float64(skillCfg.SpecialEffect.DamTransfer) / float64(BaseRate)
		if suckblood > 0 {
			attacker.ChangeHp(int32(suckblood), attacker)
		}

		o.attackId = attacker.Uid()
		o.displaySkillId = skillId

		// logger.Debugf("[%v:%s]被攻击。 攻击者Uid: %d-%s ,  技能ID: %d  伤害: %f  暴击：%t  闪避：%t",
		// 	o.uid, o.name, attacker.Uid(), attacker.Name(), skillId, damage, o.displayCrit, o.displayMiss)

		// if o.desk.roomCfg.Roomtype == constant.GuildMap && attacker.ObjectType() == ObjectTypePlayer {
		// 	o.desk.recordGuildBossDamage(attacker.Uid(), int64(damage))
		// }

		if finalDamage > 0 && o.ChangeHp(-int32(finalDamage), attacker) == 0 {
			return
		}

		// 被动怪物被攻击后选中攻击者，优先选择战士职业
		if o.objectType == ObjectTypeMonster && o.CanAttackTarget(attacker) {
			if m, err := o.desk.monsterWithUid(o.uid); err == nil && m.cfg != nil {
				if m.cfg.AttackRange == 0 {
					if m.target == nil {
						o.target = attacker
					} else if attacker.IsWarrior() {
						o.target = attacker
					}
				}
			}
		}

		// 反伤
		if bReflect && o.stat(pb.Attribute_REFLECT_DAMAGE_PERCENTAGE) > 0 {
			if buff := o.getBuffByType(constant.BuffReflectDamage); buff != nil {
				if global.RandIn(int(BaseRate)) < int(buff.Rate) {
					reflectDamage := damage * float64(o.stat(pb.Attribute_REFLECT_DAMAGE_PERCENTAGE)) / BaseRate
					reflectDamage = math.Max(reflectDamage, 1)
					attacker.ChangeHp(-int32(reflectDamage), o)
					attacker.Reflect()
					logger.Debugf(" 反伤  ID= %d:%s 伤害: %f  攻击者 %d:%s  ", o.uid, o.name, reflectDamage, attacker.Uid(), attacker.Name())
				}
			}
		}
		// 胆怯
		if buff := o.getBuffByType(constant.BuffPercentageAttackDefenseCritDecrease); buff != nil {
			attacker.AddBuff(100, o, 10000000, &Stats{Attrs: buff.BuffStats.copyStats()}, 0, 0)
		}

		// 伤害转化
		if o.stat(pb.Attribute_DAMAGE_CONVERT_PERCENTAGE) > 0 {
			suckingBlood := damage * float64(o.stat(pb.Attribute_DAMAGE_CONVERT_PERCENTAGE)/int32(BaseRate))
			suckingBlood = math.Max(suckingBlood, 1)
			o.ChangeHp(int32(suckingBlood), attacker)
		}

		// 眩晕
		if finalDamage > 0 && attacker.GetStats().getStats(pb.Attribute_STUN) > 0 {
			stum := attacker.GetStats().getStats(pb.Attribute_STUN)
			stum_resistance := o.stat(pb.Attribute_STUN_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.StunvalueEx)
			if stum >= stum_resistance || global.RandIn(int(stum_resistance)) < int(stum) {
				if cfg, err := global.GetBuffCfg(1002); err == nil {
					o.AddBuff(cfg.ID, attacker, 0, nil, 0, 0)
				}
			}
		}

		// 中毒
		if finalDamage > 0 && attacker.GetStats().getStats(pb.Attribute_POISON) > 0 {
			poison := attacker.GetStats().getStats(pb.Attribute_POISON)
			poison_resistance := o.stat(pb.Attribute_POISON_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.PoisonvalueEx)
			if poison >= poison_resistance || global.RandIn(int(poison_resistance)) < int(poison) {
				if cfg, err := global.GetBuffCfg(1001); err == nil {
					o.AddBuff(cfg.ID, attacker, 0, nil, 0, 0)
				}
			}
		}

		//被动技能触发
		for _, v := range o.skills {
			if time.Now().Before(v.Cooldown) {
				continue
			}

			if cfg, err := global.GetSkillCfg(v.CfgId); err == nil {
				if cfg.ActiveorPassive == PassiveSkill {
					o.passiveSkillTriger(attacker, v.CfgId, v.AdvancedLevel, finalDamage)
					v.Cooldown = time.Now().Add(time.Duration(cfg.CoolDown) * time.Millisecond)
				}
			}
		}
	}

	// 技能关联buff触发
	if !o.dead {
		for _, v := range skillCfg.TriggerBuff {
			bufCfg, err := global.GetBuffCfg(v)
			if err != nil {
				continue
			}
			if bufCfg.Type == constant.BuffTypeBuff {
				continue
			}

			// 免疫debuff
			// if immuneBuff := o.getBuffByType(constant.BuffNoNegativeEffect); immuneBuff != nil {
			// 	// 有几率免疫
			// 	if global.RandIn(int(BaseRate)) < int(immuneBuff.Rate) {
			// 		continue
			// 	}
			// }
			// if immuneBuff := o.getBuffByType(constant.BuffImmune); immuneBuff != nil {
			// 	continue
			// }

			var addSec, addRate, addDam, addDistance, subCtl int32
			// 基础buff持续时间 与触发概率
			addSec += skillCfg.Bufftime
			addRate += skillCfg.Effecttrigger
			// 技能进阶效果
			if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skillAdvanceLevel); err == nil {
				subCtl += skillUpperlevel.Control
				addDam += skillUpperlevel.DamagePer

				addSec += skillUpperlevel.DurationCon
				addRate += skillUpperlevel.Trigger
				addDistance += skillUpperlevel.Distance
			}

			if skillCfg.Skilltrigger > 0 && global.RandIn(int(BaseRate)) >= int(skillCfg.Skilltrigger+addRate) {
				continue
			}

			var bufStats *Stats
			if bufCfg.Mark == constant.BuffTypeChainSkill {
				bufStats = &Stats{Attrs: make(map[pb.Attribute]int32)}

				if bufCfg.Pretype == constant.BuffFixedDamage2 {
					for _, v := range skillCfg.Skillbuff {
						bleeding := float64(attacker.GetStats().getStats(attacker.AttackAttributes())) * float64(v.Value) / BaseRate
						bleeding = math.Max(1, bleeding-float64(o.stat(maputil.DefenseAttributes(attacker.AttackAttributes()))))
						bufStats.setStats(pb.Attribute_HEALTH, int32(bleeding))
					}
				} else if bufCfg.Pretype == int32(constant.BuffPercentageAttackDecrease) { // 减双攻
					for _, v := range skillCfg.Skillbuff {
						if v.Key != int32(pb.Attribute_PHYSICAL_ATTACK) && v.Key != int32(pb.Attribute_HEALTH_PERCENTAGE) &&
							v.Key != int32(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE) && v.Key != int32(pb.Attribute_MAGIC_ATTACK_PERCENTAGE) {
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}
				} else if bufCfg.Pretype == int32(constant.BuffPercentageDefenseHealthDecrease) { // 减双防 减回血
					for _, v := range skillCfg.Skillbuff {
						if v.Key != int32(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE) && v.Key != int32(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE) && v.Key != int32(pb.Attribute_LIFE_RECOVERY_DECREASE_PERCENTAGE) {
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}
				} else if bufCfg.Pretype == int32(constant.BuffPercentageAttackSpeedDecrease) { // 减攻击速度
					for _, v := range skillCfg.Skillbuff {
						if v.Key != int32(pb.Attribute_ATTACK_SPEED_PERCENTAGE) {
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}
				} else if bufCfg.Pretype == int32(constant.BuffPercentageAccuracyDecrease) { // 减命中
					for _, v := range skillCfg.Skillbuff {
						if v.Key != int32(pb.Attribute_HIT_PERCENTAGE) {
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}
				} else {
					for _, v := range skillCfg.Skillbuff {
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}
				}

				// 控制时间减去角色免控时间
				if bufCfg.Pretype == constant.BuffFreeze || bufCfg.Pretype == constant.BuffTaunt {
					if myBuff := o.getBuffByType(constant.BuffTauntTimeReduce); myBuff != nil {
						redeceSec := float64(addSec+bufCfg.LastTime) * float64(o.stats.getStats(pb.Attribute_REDUCE_CONTROL_TIME_PERCENTAGE)) / BaseRate
						addSec -= int32(redeceSec)
					}
					if skillCfg.SType == constant.SkillTypeCareerSpecial4 || skillCfg.SType == constant.SkillTypeCareerSpecial5 {
						redeceSec := float64(o.stat(pb.Attribute_CONTROL_TIME_REDUCTION_FIXED_VALUE))
						redeceSec += float64(addSec) * float64(o.stat(pb.Attribute_CONTROL_TIME_REDUCTION_PERCENTAGE)) / BaseRate
						addSec -= int32(redeceSec)
					}
				}
			}
			if addSec > 0 {
				o.AddBuff(bufCfg.ID, attacker, addSec, bufStats, addRate, skillId)
			}
		}
	}

	// 技能特殊效果触发
	switch skillCfg.SkillID {
	case constant.SkillTypeRepelFoe:
		o.repelFoeEnd(skillId, attacker)

	case constant.SkillTypeNormalAttackWizard:
		if attacker.HaveBuff(constant.BuffKnockback) {
			if o.HaveBuff(constant.BuffFreeze) {
				break
			}

			if buff := attacker.GetBuff(constant.BuffKnockback); buff != nil {
				if buffSkillCfg, err := global.GetSkillCfg(buff.SkillID); err == nil {
					if global.RandIn(int(BaseRate)) >= int(buff.Rate) {
						return
					}

					var backDir pb.Direction
					var endCell *Cell
					backDir = maputil.DirectionFromPointI(&attacker.GetCell().location, &o.cell.location)
					for i := 1; i <= int(buffSkillCfg.SpecialEffect.BackDistence); i++ {
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
					if o.desk.mapInfo.locationChange(o, &pb.Point{X: float64(endCell.location.X), Y: float64(endCell.location.Y)}) {
						o.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: o.uid, Location: o.location, Direction: o.direction})
					}
				}
			}
		}

	case constant.SkillTypeStellarBurst:
		// 星爆术增加几率斩杀,只对玩家生效
		if o.objectType != ObjectTypeMonster && global.RandIn(int(BaseRate)) < int(skillCfg.Skilltrigger) {

			curhpPre := float64(o.currentHP) / float64(o.GetMaxHP())
			slayPer := float64(skillCfg.SpecialEffect.Slay) / BaseRate
			if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skillAdvanceLevel); err == nil {
				slayPer += float64(skillUpperlevel.SlayPer) / BaseRate
			}
			if curhpPre < slayPer {
				if o.desk.roomCfg.Roomtype == constant.FightEndMap {

					//星爆术在一战到底，最高伤害不超过当前攻击力3倍
					dmg := float64(o.currentHP)
					if dmg > float64(attacker.GetStats().getStats(attacker.AttackAttributes())*3) {
						dmg = float64(attacker.GetStats().getStats(attacker.AttackAttributes()) * 3)
					}
					o.ChangeHp(int32(-dmg), attacker)

				} else {
					//logger.Debugf("[%d:%s]被星爆术斩杀,当前血量 %d  总血量 %d  [%f:%f]", o.uid, o.name, o.currentHP, o.GetMaxHP(), curhpPre, slayPer)
					o.ChangeHp(int32(-o.currentHP), attacker)
				}
			}
		}
	case constant.ThorSkill1:
		if skillCfg.SpecialEffect.HPReducePer > 0 {
			hpReducePer := float64(o.GetMaxHP()) * float64(skillCfg.SpecialEffect.HPReducePer) / float64(BaseRate)
			o.ChangeHp(int32(-hpReducePer), attacker)
			//logger.Debugf("[%d:%s]被攻击,雷神1技能,当前血量 %d  总血量 %d  [%f]", o.uid, o.name, o.currentHP, o.GetMaxHP(), hpReducePer)
		}

	case constant.SkillTypeStarShatter:
		o.starShatterEnd(skillId, attacker)
	}
}

// 死亡
func (o *MapObject) Die(attacker IMapObject) {
	if o.dead {
		return
	}
	o.die(attacker)
}

// 复活
func (o *MapObject) Revive(skill *PlayerSkill, skillAdvanceLevel, reviveHp int32) {
	if o.desk == nil {
		logger.Errorf("[%s]: %d 玩家复活！！！ desk为空", o.name, o.uid)
		return
	}

	o.dead = false
	//o.reviveTimes++
	o.currentHP = reviveHp
	o.displayHP = o.currentHP
	o.stats.setStats(pb.Attribute_HEALTH, o.currentHP)
	o.resetBuff()

	o.SkillBuffTrigger(skill)

	o.desk.Broadcast(protocol.RouteRevivePush, &pb.RevivePush{
		Uid:       o.uid,
		Location:  o.location,
		Direction: o.direction,
		CurHp:     o.currentHP,
	})

	o.RefreshStats()

	// logger.Debugf("玩家复活！！！   当前血量 %d -- %d --- %d", o.currentHP, o.stat(pb.Attribute_HEALTH), o.GetMaxHP())
}

func (o *MapObject) CleanDebuff() {
	//if o.block.TryLock() {
	//	defer o.block.Unlock()
	//
	//	var delbuffs []int32
	//	for k := range o.buffs {
	//		if o.buffs[k].BuffType == constant.BuffTypeDebuff {
	//			delbuffs = append(delbuffs, k)
	//		}
	//	}
	//
	//	for i := 0; i < len(delbuffs); i++ {
	//		o.delBuff(delbuffs[i])
	//	}
	//}
	o.buffMgr.CleanDebuff()
}

// TODO 注意玩家在场景里,须在场景AddBuff
func (o *MapObject) AddBuff(id int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32) {
	o.AddBuff2(0, id, owner, addSec, exStats, addRate, skillID)
}

func (o *MapObject) AddBuff2(buffID int64, id int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32) {
	if !o.buffMgr.CanAddBuff(id) {
		return
	}

	o.addBuff(buffID, id, owner, addSec, exStats, addRate, skillID)

	//logger.Debugf(" [%s] 增加buff 之前,  BuffId=%d  生命 %d  最大生命 %d , 攻击 %d : %d  防御 %d:%d 暴击 %d", o.name, id, o.currentHP, o.GetMaxHP(), o.GetDC(), o.GetMC(), o.GetAR(), o.GetMR(), o.GetStats().getStats(pb.Attribute_CRITICAL_STRIKE))

	o.RefreshStats()

	//logger.Debugf(" [%s] 增加buff 之后,  BuffId=%d  生命 %d  最大生命 %d , 攻击 %d : %d  防御 %d:%d 暴击 %d", o.name, id, o.currentHP, o.GetMaxHP(), o.GetDC(), o.GetMC(), o.GetAR(), o.GetMR(), o.GetStats().getStats(pb.Attribute_CRITICAL_STRIKE))
}

func (o *MapObject) AddBuffNoRefresh(id int32, owner IMapObject, addSec int32, exStats *Stats, addRate int32, skillID int32) {
	o.addBuff(0, id, owner, addSec, exStats, addRate, skillID)
}

func (o *MapObject) Friends() []IMapObject {
	var friends []IMapObject
	if o.desk == nil {
		return friends
	}

	if o.objectType == ObjectTypePlayer {
		friends = append(friends, o)
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			for _, slot := range p.partnerSlots {
				if slot.PartnerUid == 0 {
					continue
				}
				if partner, err := o.desk.objectWithId(slot.PartnerUid); err == nil {
					friends = append(friends, partner)
				}
			}
		}
	} else if o.objectType == ObjectTypePartner && o.owner != nil {
		if p, err := o.desk.playerWithId(o.owner.Uid()); err == nil {
			for _, slot := range p.partnerSlots {
				if slot.PartnerUid == 0 {
					continue
				}
				if partner := p.getPartnerInfoBySlotId(slot.SlotId); partner != nil {
					friends = append(friends, partner)
				}
			}
			friends = append(friends, p)
		}
		if r, err := o.desk.robotWithId(o.owner.Uid()); err == nil {
			for _, slot := range r.partnerSlots {
				if slot.PartnerUid == 0 {
					continue
				}
				if partner := r.getPartnerInfoBySlotId(slot.SlotId); partner != nil {
					friends = append(friends, partner)
				}
			}
			friends = append(friends, r)
		}
	} else if o.objectType == ObjectTypeRobot {
		friends = append(friends, o)
		if r, err := o.desk.robotWithId(o.uid); err == nil {
			for _, slot := range r.partnerSlots {
				if slot.PartnerUid == 0 {
					continue
				}
				if partner, err := o.desk.objectWithId(slot.PartnerUid); err == nil {
					friends = append(friends, partner)
				}
			}
		}
	}
	return friends
}

// 技能buff触发
func (o *MapObject) SkillBuffTrigger(skill *PlayerSkill) {
	if skill == nil {
		return
	}

	skillCfg, err := global.GetSkillCfg(skill.CfgId)
	if err == nil {

		//女神降临技能有2个buff goddessDescentEnd()中处理
		if skillCfg.SkillID == constant.SkillTypeGoddessDescent {
			return
		}

		for _, tri := range skillCfg.TriggerBuff {
			bufCfg, err := global.GetBuffCfg(tri)
			if err != nil {
				continue
			}
			if bufCfg.Type == constant.BuffTypeDebuff {
				continue
			}

			// 筛除特殊buff
			if bufCfg.Pretype == constant.BuffDecoy {
				continue
			}

			var addSec, addRate, addDam, addDistance, subCtl int32
			addRate += skillCfg.Effecttrigger
			addSec += skillCfg.Bufftime
			// 技能进阶效果
			if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skill.AdvancedLevel); err == nil {
				subCtl += skillUpperlevel.Control
				addDam += skillUpperlevel.DamagePer

				addSec += skillUpperlevel.DurationCon
				addRate += skillUpperlevel.Trigger
				addDistance += skillUpperlevel.Distance
			}

			// 技能关联的buff 并且 有触发几率 才计算概率
			if bufCfg.Mark == constant.BuffTypeChainSkill && skillCfg.Skilltrigger > 0 && global.RandIn(int(BaseRate)) >= int(skillCfg.Skilltrigger+addRate) {
				continue
			}

			var bufStats *Stats
			if bufCfg.Mark == constant.BuffTypeChainSkill {
				bufStats = &Stats{Attrs: make(map[pb.Attribute]int32)}
				for _, v := range skillCfg.Skillbuff {
					bufStats.setStats(pb.Attribute(v.Key), v.Value)
				}

				if addDam > 0 {
					bufStats.setStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, addDam)
				}
				if subCtl > 0 {
					bufStats.setStats(pb.Attribute_CONTROL_TIME_REDUCTION_PERCENTAGE, subCtl)
				}
			}

			// 影响范围
			var target_friend []IMapObject
			if skillCfg.Effectnum == constant.SkillEffectnumSlef {
				target_friend = append(target_friend, o)
			} else if skillCfg.Effectnum == constant.SkillEffectnumAlliance {
				target_friend = o.Friends()
			} else if skillCfg.Effectnum == constant.SkillEffectnumAllianceLowHp {

				low_hp := float64(o.currentHP) / float64(o.GetMaxHP())
				var target IMapObject = o
				for _, v := range o.Friends() {
					if v.Dead() {
						continue
					}

					tmp_low := float64(v.GetStats().getStats(pb.Attribute_HEALTH)) / float64(v.GetMaxHP())
					if tmp_low < low_hp {
						low_hp = tmp_low
						target = v
					}
				}
				target_friend = append(target_friend, target)
				logger.Debugf("%s 释放技能 [%d] 选择血量最低目标 %d:%s , 血量百分比 %.2f", o.name, skill.CfgId, target.Uid(), target.Name(), low_hp)
			} else if skillCfg.Effectnum == constant.SkillEffectnumAllianceAndGuild {
				target_friend = o.Friends()
				attackCells := o.desk.mapInfo.getSquareCells(&o.cell.location, skillCfg.Radius)
				for _, cell := range attackCells {
					if cell == nil || cell.block {
						continue
					}
					for _, obj := range cell.objects {
						bFind := false
						for _, have := range target_friend {
							if have.Uid() == obj.Uid() {
								bFind = true
								break
							}
						}
						if !bFind {
							if o.desk.roomCfg.Roomtype == constant.PVEWarMap {
								if o.GuildId() != 0 && o.GuildId() == obj.GuildId() {
									target_friend = append(target_friend, obj)
								}
							}
						}
					}
				}

			}

			for _, v := range target_friend {
				//只对法师生效
				if skillCfg.SkillID == constant.SkillTypeMassPurificationBoss && (v.Job() == constant.Ritualist || v.Job() == constant.PartnerRitualist ||
					v.Job() == constant.Warrior || v.Job() == constant.PartnerWarrior) {
					continue
				}
				if v.Dead() {
					continue
				}
				v.AddBuff(bufCfg.ID, o, addSec, bufStats, addRate, skill.CfgId)
			}
		}
	}
}

func (o *MapObject) ChangeHp(hp int32, obj IMapObject) int32 {
	if obj != nil {
		o.AttackId(obj.Uid())
	}

	return o.changeHp(hp, obj)
}

func (o *MapObject) IsSameGuild(guildId int64, obj IMapObject) bool {
	if o.desk == nil {
		return false
	}

	if guildId == 0 {
		return false
	}

	if o.objectType == ObjectTypePlayer {
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			return p.guildId == guildId
		}
	} else if o.objectType == ObjectTypePartner {
		if p, err := o.desk.playerWithId(o.owner.Uid()); err == nil {
			return p.guildId == guildId
		}
	}
	return false
}

func (o *MapObject) GuildId() int64 {
	if o.desk == nil {
		return 0
	}

	if o.objectType == ObjectTypePlayer {
		if p, err := o.desk.playerWithId(o.uid); err == nil {
			return p.guildId
		}
	} else if o.objectType == ObjectTypePartner {
		if p, ok := defaultManager.player(o.owner.Uid()); ok {
			return p.guildId
		}
	}
	return 0
}

func (o *MapObject) GetNormalSkillID() int32 {
	for _, v := range o.skills {
		if cfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			if cfg.SType == constant.SkillTypeNormal {
				return v.CfgId
			}
		}
	}
	return 0
}

// 技能生成的实例有自己的AI
func (o *MapObject) processAI() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	if o.desk == nil {
		return
	}
	if o.objectType != ObjectTypeBuff {
		return
	}
	if o.dead {
		o.desk.delObject(o.uid)
		return
	}

	if time.Now().After(o.delTime) {
		o.desk.delObject(o.uid)
		return
	}

	o.processRoam()
	o.processSearch()
	o.processTarget()
}

func (o *MapObject) processRoam() {

}
func (o *MapObject) processSearch() {

}
func (o *MapObject) processTarget() {
	if !o.CanAttack() {
		return
	}

	if o.objectType != ObjectTypeBuff {
		return
	}

	buffCfg, err := global.GetBuffCfg(o.cfgId)
	if err != nil {
		return
	}

	o.setAttackTime(buffCfg.EffectTime)

	switch buffCfg.Pretype {
	case constant.BuffTornado:
		o.tornadoAttackEnd()
	}
}
