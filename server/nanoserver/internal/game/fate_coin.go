package game

import (
	"nanoserver/global"
	"nanoserver/pkg/event"
	"nanoserver/pkg/randx"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"time"
)

// 命运金币
type FateCoin struct {
	lastTime      time.Time
	dieMonsterMap map[int64]*MonsterInfo
	cfgId         int32
	count         int32
	triggerMap    map[int32]int32
	triggerTime   int32
}

func NewFateCoin() *FateCoin {
	fateCoin := &FateCoin{dieMonsterMap: map[int64]*MonsterInfo{}, lastTime: time.Time{}, triggerMap: make(map[int32]int32)}
	var id int32
	for _, v := range global.GlobalCfg.FateCoinTrigger.Key {
		fateCoin.triggerMap[v.ID] = v.Time
		if v.ID > id {
			fateCoin.triggerTime = v.Time
		}
	}
	return fateCoin
}

func (d *FateCoin) getTriggerTime() int32 {
	cdTime := d.triggerTime
	triggerCount := d.count + 1
	if cTime, ok := d.triggerMap[triggerCount]; ok {
		cdTime = cTime
	}
	return cdTime
}

// trigger 定时触发
func (d *FateCoin) trigger(now time.Time, desk *Desk) {
	if d.triggerTime == 0 {
		return
	}
	if d.lastTime.IsZero() {
		d.lastTime = desk.beginTime.Add(time.Duration(d.getTriggerTime()) * time.Second)
	}
	if now.After(d.lastTime) {
		rate := int32(global.RandIn(100))
		if rate > global.GlobalCfg.BaseValue.FateCoin.Rate {
			return
		}
		d.registerEffect(desk)
		d.lastTime = d.lastTime.Add(time.Duration(d.getTriggerTime()) * time.Second)
	}
}

// registerEffect
func (d *FateCoin) registerEffect(desk *Desk) {
	if len(desk.BattleCfg.BattleEffect) == 0 {
		return
	}
	picker := randx.New(desk.BattleCfg.BattleEffect, func(i global.IdNumRate) int {
		return int(i.Rate)
	})
	effect, err := picker.Pick()
	if err != nil {
		desk.logger.Errorf("命运金币触发失败 err=%s", err.Error())
		return
	}
	if effect.Id == 0 {
		return
	}
	d.cfgId = effect.Id
	d.count++
	if len(desk.players) == 0 {
		return
	}
	p := desk.players[0]
	packet := &pb.FateCoinPush{CfgId: effect.Id}
	p.pushMsg(protocol.RouteFateCoinPush, packet)
	desk.logger.Debugf("命运金币触发，事件CfgId=%d, 玩家=%d", effect.Id, p.uid)
}

func (d *FateCoin) onEffect(p *Player) (pb.ErrorCode, int32) {
	cfgId := d.cfgId
	if cfgId > 0 {
		desk := p.desk
		if desk != nil {
			desk.eventMgr.Publish(event.EventTypeFateCoinEffect, &event.EventData{UId: p.uid, CfgId: cfgId})
		}
		d.cfgId = 0
	}
	p.logger.Debugf("命运金币领取事件CfgId=%d, 玩家=%d", cfgId, p.uid)
	return pb.ErrorCode_ERROR_CODE_SUCCESS, cfgId
}

func (d *FateCoin) isReward() bool {
	return d.count > 0
}

func (d *FateCoin) putDieMonster(monster *MonsterInfo) {
	d.dieMonsterMap[monster.uid] = monster
}

// 重置
func (d *FateCoin) reset() {
	d.lastTime = time.Time{}
	d.dieMonsterMap = map[int64]*MonsterInfo{}
	d.cfgId = 0
	d.count = 0
}
