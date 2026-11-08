package game

import (
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/event"
	"nanoserver/pkg/randx"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol/protobuf/pb"
)

// registerEventHandlers 注册事件处理器
func (d *Desk) registerEventHandlers() {
	d.eventMgr.Subscribe(event.EventTypeFateCoinEffect, d.onFateCoinEffect)
	d.eventMgr.Subscribe(event.EventTypeRefreshStats, d.onRefreshStats)
}

// onFateCoinEffect 命运金币效果
func (d *Desk) onFateCoinEffect(data interface{}) {
	eventData, ok := data.(*event.EventData)
	if !ok {
		return
	}
	cfg, ok := global.FateCoinCfg[eventData.CfgId]
	if !ok {
		return
	}
	p, err := d.playerWithId(eventData.UId)
	if err != nil {
		return
	}
	switch cfg.EffectType {
	case pb.FateCoinEffect_FC_EFFECT_NONE:
		return

	case pb.FateCoinEffect_FC_EFFECT_BUFF:
		d.handleFateCoinBuff(cfg, p)

	case pb.FateCoinEffect_FC_EFFECT_RANDOM_DEATH:
		d.handleFateCoinRandomDeath(cfg, p)

	case pb.FateCoinEffect_FC_EFFECT_REVIVE:
		d.handleFateCoinRevive(cfg)

	case pb.FateCoinEffect_FC_EFFECT_REVIVE_RECOVER_HP:
		d.handleFateCoinReviveRecoverHP(cfg, p)
	}
}

// handleFateCoinBuff 命运金币buff
func (d *Desk) handleFateCoinBuff(cfg *global.FateCoinConfig, p *Player) {
	_, err := global.GetBuffCfg(cfg.BuffID)
	if err != nil {
		return
	}
	bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, v := range cfg.SkillBuff {
		bufStats.setStats(pb.Attribute(v.Key), v.Value)
	}
	targets := d.getTargets(cfg.EffectNum, p)
	for _, v := range targets {
		if v.Dead() {
			continue
		}
		if cfg.EffectNum != constant.SkillEffectnumEnemyAll {
			buffId := snowflake.GenID()
			v.AddBuff2(buffId, cfg.BuffID, p, cfg.BuffTime, bufStats, 0, 0)
		} else {
			v.AddBuff(cfg.BuffID, p, cfg.BuffTime, bufStats, 0, 0)
		}
	}
}

// handleFateCoinRevive 命运金币复活
func (d *Desk) handleFateCoinRevive(cfg *global.FateCoinConfig) {
	for _, v := range d.dieMonsterMap {
		monRelive := &MonRelive{CfgID: v.cfgId, Location: global.PointI{X: v.bornLocation.X, Y: v.bornLocation.Y}}
		Relive(d, monRelive)
	}
	d.dieMonsterMap = map[int64]*MonsterInfo{}
}

// handleFateCoinReviveRecoverHP 命运金币复活和回血
func (d *Desk) handleFateCoinReviveRecoverHP(cfg *global.FateCoinConfig, p *Player) {
	for _, target := range d.getTargets(cfg.EffectNum, p) {
		if target.Dead() {
			var hp float64
			for _, skillBuff := range cfg.SkillBuff {
				if skillBuff.Key == int32(pb.Attribute_HEALTH_PERCENTAGE) {
					hp += float64(target.GetMaxHP()) * float64(skillBuff.Value) / BaseRate
				}
			}
			d.reliveNormal(&ReliveData{pb.ReliveType_RT_RELIVE_NORMAL, target.Uid(), target.ObjectType(), int32(hp)})
		}
	}
}

// handleFateCoinRandomDeath 命运金币随机死亡
func (d *Desk) handleFateCoinRandomDeath(cfg *global.FateCoinConfig, p *Player) {
	aliveTargets := d.getFateCoinAliveTargets(cfg.EffectNum, p)
	if len(aliveTargets) == 0 {
		return
	}
	index, _ := randx.RandUint32(uint32(len(aliveTargets)))
	aliveTargets[index].Die(nil)
}

// getFateCoinAliveTargets 命运金币存活目标
func (d *Desk) getFateCoinAliveTargets(effectNum int32, p *Player) []IMapObject {
	var aliveTargets []IMapObject
	for _, target := range d.getTargets(effectNum, p) {
		if !target.Dead() {
			aliveTargets = append(aliveTargets, target)
		}
	}
	return aliveTargets
}

// 房间内刷新属性
func (d *Desk) onRefreshStats(data interface{}) {
	uid, ok := data.(int64)
	if !ok {
		return
	}
	if obj := d.objectByUid(uid); obj != nil {
		obj.RefreshStats()
	}
}
