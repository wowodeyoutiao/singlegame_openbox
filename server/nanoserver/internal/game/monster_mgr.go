package game

import (
	"context"
	"strconv"
	"sync"
	"time"

	"nanoserver/dispatcher"
	"nanoserver/global"
	"nanoserver/pkg/astar"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/goredisx"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/redis/go-redis/v9"
)

const DefaultAttackSpeed = 500
const DefaultMoveSpeed = 250

const (
	MTypeNormal     int32 = iota // 0：代表不是boss
	MTypeCommonBoss              // 1：代表普通boss
	MTypeWorldBoss               // 2：代表世界boss
	MTypeGathering               // 3：代表采集怪物
	MTypeWarBoss                 // 4：代表战场boss
)

type MonsterInfo struct {
	lock sync.Mutex
	*MapObject
	level           int32 // 等级
	createTime      int64 // 创建时间(毫秒时间戳) time.Now().UnixNano()/1e6
	reliveTime      int64 // 重生时间
	cfg             *global.MonsterConfig
	roamPath        []*astar.APoint
	searchTime      time.Time
	roamTime        time.Time
	miningPlayer    *Player   // 采集怪物的采集者
	miningStartTime time.Time // 采集怪物的采集开始时间
	ownershipStart  time.Time
	ownershipType   pb.OwnershipType // 怪物归属类型
	ownershipUid    int64            // 怪物归属ID
	ownershipName   string           // 归属名字
}

func newMonster(cfgId int32) *MonsterInfo {
	now := time.Now().Unix()
	uid := snowflake.GenID()

	cfg := global.GetMonsterCfg(cfgId)
	if cfg == nil {
		return nil
	}

	m := &MonsterInfo{
		level:      cfg.Level,
		createTime: now,
		reliveTime: int64(cfg.Relive),
		MapObject:  NewMapObject(ObjectTypeMonster, uid),
		cfg:        cfg,
		roamPath:   []*astar.APoint{},
	}

	for _, v := range cfg.BassProperty {
		m.stats.setStats(pb.Attribute(v.Key), v.Value)
	}
	for _, v := range cfg.Skills {
		m.skills = append(m.skills, &PlayerSkill{CfgId: int32(v), Cooldown: time.Now()})
	}
	m.stats.setStats(pb.Attribute(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE), int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	m.cfgId = cfgId
	m.name = cfg.Name
	m.currentHP = m.stats.getStats(pb.Attribute_MAX_HEALTH)
	m.displayHP = m.currentHP
	m.stats.setStats(pb.Attribute_HEALTH, m.stats.getStats(pb.Attribute_MAX_HEALTH))
	m.roamTime = time.Now().Add(time.Second * time.Duration(global.RandIn(90)+1))
	return m
}

// 死亡
func (m *MonsterInfo) die(attacker IMapObject) {
	if m.desk == nil {
		logger.Errorf("异常 怪物死亡,房间nil , MonsterUId : %d , MonsterCfgID: %d ", m.uid, m.cfgId)
		return
	}
	desk := m.desk
	cfgId := m.cfgId

	var player *Player
	// 计算掉落
	if attacker != nil && m.cfg.MType != MTypeWorldBoss {
		var uid int64

		if attacker.ObjectType() == ObjectTypePlayer {
			uid = attacker.Uid()
		} else if attacker.ObjectType() == ObjectTypePartner && attacker.Owner() != nil {
			uid = attacker.Owner().Uid()
		}
		if player, err := m.desk.playerWithId(uid); err == nil {
			m.drop(player)
			player.addKillMonControl(m.cfgId, 1)
		}
	}

	packet := pb.ObjectDiePush{ObjectId: m.uid, ReviveTip: false}
	desk.Broadcast(protocol.RouteObjectDiePush, &packet)

	//logger.Debugf(" ==死亡== 怪物 , MonsterCfgID:%d uid %d : %s", m.cfgId, m.uid, m.name)

	// 怪物重生
	if desk.roomCfg.Roomtype == constant.PVEWarMap {
		if m.reliveTime > 0 {
			if m.cfg.MType == MTypeWarBoss {
				m.desk.refreshList = append(m.desk.refreshList, &MonsterRefresh{
					CfgId:       m.cfgId,
					RefreshTime: time.Now().Add(time.Duration(m.reliveTime) * time.Millisecond),
					Location:    m.bornLocation})
				m.desk.refreshMonster()
			} else {
				mdesk := desk
				reliveTime := m.reliveTime
				time.AfterFunc(time.Millisecond*time.Duration(reliveTime), func() {
					if mdesk != nil {
						mdesk.addMonsterRelive(cfgId, m.bornLocation, 0, 0, 0, 0, "")
					}
				})
			}
		}

		for _, v := range global.GlobalCfg.War.Monstergrow.MonsterAttr {
			if v.ID == m.cfg.ID {
				dispatcher.MustGetInstance().Publish(dispatcher.TopicPVEWarMapBossDieCount.String(), &dispatcher.PVEWarMapBossDieCountEvent{
					ZoneID:    desk.zoneID,
					MonsterID: m.cfgId,
				})
			}
		}
	}

	if desk.roomCfg.Roomtype == constant.PVEPeacefulMap ||
		desk.roomCfg.Roomtype == constant.StarUpStoneFBMap || desk.roomCfg.Roomtype == constant.SoulStoneFBMap {

		mdesk := desk
		reliveTime := m.reliveTime
		if reliveTime > 0 {
			time.AfterFunc(time.Millisecond*time.Duration(reliveTime), func() {
				if mdesk != nil {
					mdesk.addMonsterRelive(cfgId, m.bornLocation, 0, 0, 0, 0, "")
				}
			})
		}
	}

	if player != nil {
		player.tasks.dispatch(pb.TargetType_KILL_MONSTER_REQUIREMENT, 0, 1)
		if m.cfg.MType == MTypeCommonBoss {
			player.tasks.dispatch(pb.TargetType_TOTAL_KILLS_OF_BOSS_MONSTERS, 0, 1)
		} else {
			player.tasks.dispatch(pb.TargetType_TOTAL_KILLS_OF_NORMAL_MONSTERS, 0, 1)
		}

		if desk.roomCfg.Roomtype == constant.PVEWarMap {
			player.tasks.dispatch(pb.TargetType_BATTLEFIELD_TOTAL_KILLS, 0, 1)
		}
	}

	if desk.roomCfg.Roomtype == constant.GuildMap {
		for _, nr := range desk.GuildReward.NexReward {
			bFinde := false
			for _, cr := range desk.GuildReward.CurReward {
				if nr.K == cr.K {
					cr.V += nr.V
					bFinde = true
					break
				}
			}
			if !bFinde {
				desk.GuildReward.CurReward = append(desk.GuildReward.CurReward, nr)
			}
		}

		desk.GuildReward.NexReward = []*pb.KeyValue{}
		if bcfg := global.GetBattleMapCfg(desk.BattleId + 1); bcfg != nil {
			for _, v := range bcfg.Settlement {
				desk.GuildReward.NexReward = append(desk.GuildReward.NexReward, &pb.KeyValue{K: v.ID, V: v.Num})
			}
		}
		desk.Broadcast(protocol.RouteGuildCopyRewardPush, desk.GuildReward)
	}
	// 战场boss死亡计数
	if desk.roomCfg.Roomtype == constant.PVEWarMap {
		for _, v := range global.GlobalCfg.War.Monstergrow.MonsterAttr {
			if v.ID == m.cfgId {
				desk.PveWarMap.monsterDeadTimes[m.cfgId]++
				break
			}
		}
	}
	if desk.roomCfg.Roomtype == constant.BattleMap {
		desk.FateCoin.putDieMonster(m)
	}
	// 重置并归还
	desk.delMonster(m)
	m.reset()
	monsterMgr.GiveBackMonster(m)
}

// 重生
func Relive(desk *Desk, mr *MonRelive) {
	if desk == nil {
		return
	}

	m, err := monsterMgr.TakeMonster(mr.CfgID)
	if err == nil {
		m.bornLocation = mr.Location

		if cell := desk.mapInfo.getAvailableCell(&m.bornLocation, int32(m.cfg.DistributeRange)); cell != nil {
			m.uid = snowflake.GenID()
			m.ownershipType = mr.OwnershipType
			m.ownershipName = mr.OwnershipName
			m.ownershipUid = mr.OwnershipUid
			m.ownershipStart = time.Now()
			m.linkedSkill = mr.linkedSkill
			lifespan := mr.Lifespan
			if lifespan <= 0 {
				lifespan = 24 * 60 * 60 * 1000
			}
			m.delTime = time.Now().Add(time.Duration(lifespan) * time.Millisecond)

			if m.cell = desk.mapInfo.addObject(m, &cell.location); m.cell != nil {
				m.setDesk(desk)
				desk.addMonster(m)
			}

			// 战场怪物死亡成长机制.增加buff 13
			if desk.roomCfg.Roomtype == constant.PVEWarMap {
				//desk.logger.Debugf("战场怪物重生  MonsterCfgID: %d  位置%v  归属类型 :%d  归属名字:%s  归属uid:%d", m.cfgId, cell.location, m.ownershipType, m.ownershipName, m.ownershipUid)

				// if deadTimes, ok := desk.PveWarMap.monsterDeadTimes[m.cfgId]; ok && deadTimes > 0 {
				// 	for _, v := range global.GlobalCfg.War.Monstergrow.MonsterAttr {
				// 		if v.ID == m.cfgId {
				// 			bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
				// 			bufStats.Attrs[pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE] = deadTimes * int32(v.Num)
				// 			bufStats.Attrs[pb.Attribute_MAGIC_ATTACK_PERCENTAGE] = deadTimes * int32(v.Num)
				// 			bufStats.Attrs[pb.Attribute_HEALTH_PERCENTAGE] = deadTimes * int32(v.Num)
				// 			bufStats.Attrs[pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE] = deadTimes * int32(v.Num)
				// 			bufStats.Attrs[pb.Attribute_MAGIC_DEFENSE_PERCENTAGE] = deadTimes * int32(v.Num)
				// 			m.AddBuff(13, m, 1000*60*60*24, bufStats, 0, 0)

				// 			desk.logger.Debugf("怪物死亡成长机制.增加buff 13  怪物id:%d 死亡次数:%d  bufStats:%v", m.cfgId, deadTimes, bufStats)
				// 			break
				// 		}
				// 	}
				// }

				// 给怪物加上永久增益buff，从redis获取
				for _, entity := range global.GlobalCfg.War.Monstergrow.MonsterAttr {
					if mr.CfgID == entity.ID {
						key := "pve_war_map:boss_die:" + strconv.FormatInt(int64(desk.zoneID), 10) + ":" + strconv.FormatInt(int64(entity.ID), 10)
						bossDieCount, err := goredisx.Instance.Get(context.TODO(), key).Int()
						if err != nil {
							if err != redis.Nil {
								logger.Errorf("PVEWarMapBossDieCountHandler Get error: %v", err)
							}
						}
						var times int32
						if int32(bossDieCount) >= entity.MaxTime {
							times = entity.MaxTime
						} else {
							times = int32(bossDieCount)
						}
						bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
						bufStats.Attrs[pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE] = times * int32(entity.Num)
						bufStats.Attrs[pb.Attribute_MAGIC_ATTACK_PERCENTAGE] = times * int32(entity.Num)
						bufStats.Attrs[pb.Attribute_HEALTH_PERCENTAGE] = times * int32(entity.Num)
						bufStats.Attrs[pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE] = times * int32(entity.Num)
						bufStats.Attrs[pb.Attribute_MAGIC_DEFENSE_PERCENTAGE] = times * int32(entity.Num)
						m.AddBuff(101, m, 1000*60*60*24, bufStats, 0, 0)
					}
				}
			}

			packet := pb.MonsterInPush{MonsterList: []*pb.MonsterInfo{}}

			packet.MonsterList = append(packet.MonsterList,
				&pb.MonsterInfo{
					CfgId:     m.cfgId,
					Uid:       m.uid,
					Name:      m.name,
					MaxHp:     m.GetMaxHP(),
					CurHp:     m.currentHP,
					Location:  &pb.Point{X: m.location.X, Y: m.location.Y},
					Direction: m.direction,
					Ownership: &pb.Ownership{Type: m.ownershipType, Name: m.ownershipName, Uid: m.ownershipUid},
				})

			desk.Broadcast(protocol.RouteMonsterInPush, &packet)
		}
	}
}

// 掉落
func (m *MonsterInfo) drop(attacker *Player) {
	if attacker == nil || attacker.isRobot() {
		return
	}
	monCfg := global.GetMonsterCfg(m.cfgId)
	if monCfg == nil {
		logger.Debugf("获取怪物掉落配置失败 cfgId :%d", m.cfgId)
		return
	}

	packet := &pb.ItemDropPush{Items: []*pb.ItemInfo{}}
	packet.Location = &pb.Point{X: m.location.X, Y: m.location.Y}
	packet.ObjectId = m.uid
	reason := constant.ItemAddDrop

	// 根据房间类型设置掉落原因
	switch m.desk.roomCfg.Roomtype {
	case constant.PVEPeacefulMap:
		reason = constant.ItemAddDropHaiDao
	case constant.StarUpStoneFBMap:
		reason = constant.ItemAddDropStarUpStoneFB
	case constant.SoulStoneFBMap:
		reason = constant.ItemAddDropSoulStoneFBMap
	}

	// 是否需要处理掉落
	shouldProcessDrop := m.desk.roomCfg.Roomtype == constant.PVEWarMap ||
		m.desk.roomCfg.Roomtype == constant.WorldBossMap ||
		m.desk.roomCfg.Roomtype == constant.StarUpStoneFBMap ||
		m.desk.roomCfg.Roomtype == constant.SoulStoneFBMap ||
		m.desk.roomCfg.Roomtype == constant.GuidanceFBMap ||
		m.desk.roomCfg.Roomtype == constant.PVEPeacefulMap

	var multiplier int32 = 0

	if shouldProcessDrop {
		items := []global.IdNum{}
		miningMonsters := []global.IdNum{}

		for _, v := range monCfg.Drop {
			if len(v.ID) == 0 {
				continue
			}

			// 计算掉落倍率
			if m.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
				// 守卫海岛挂机倍率
				if guardIslandHangup, exist := global.GuardIslandCfg[attacker.AfkMultiplierID]; exist {
					// 当倍率为-1时，代表前端需要弹窗, 这个时候倍率设置为0
					if guardIslandHangup.Multiplier == -1 {
						multiplier = 0
					} else {
						multiplier = guardIslandHangup.Multiplier
					}
				} else {
					logger.Errorf("守卫海岛挂机倍率配置不存在, UserId:%d, MonsterCfgID: %d , MonsterUid: %d", attacker.Uid(), m.cfgId, m.uid)
				}
			}

			r := global.RandIn(int(BaseRate))
			ratio := float64(v.Ratio) * (1 + float64(attacker.stat(pb.Attribute_ITEM_DROP_RATE_PERCENTAGE)+multiplier)/BaseRate)

			if r < int(ratio) {
				if m.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
					logger.Debugf("守卫海岛怪物掉落, UserId:%d, MonsterCfgID: %d , MonsterUid: %d ,ratio: %v", attacker.Uid(), m.cfgId, m.uid, ratio)
				}

				index := global.RandIn(len(v.ID))
				if index >= len(v.ID) {
					index = 0
				}
				id := v.ID[index]
				num := v.Num.Min + int32(global.RandIn(int(v.Num.Max)-int(v.Num.Min)))

				// if !itemRecoderMgr.check(attacker.uid, attacker.level, id, num, m.desk.roomCfg.Roomtype) {
				// 	logger.Infof("掉落控制 drop UserId:%d, MonsterCfgID: %d ,  ItemID:%d ItemNum:%d", attacker.Uid(), m.cfgId, id, num)
				// 	continue
				// }
				if m.desk != nil && m.desk.roomCfg != nil {
					if !DropControlInstance.TryDrop(id, int64(num), m.desk.roomCfg.MapId, attacker) {
						continue
					}
				}

				if itemCfg, err := global.GetItemCfg(id); err == nil {
					// 掉落可采集的宝箱
					if itemCfg.MainType == constant.ItemMainTypeMiningBox {
						miningMonsters = append(miningMonsters, global.IdNum{ID: id, Num: num})
					} else {
						items = append(items, global.IdNum{ID: id, Num: num})
						//packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: id, Num: int64(num)})
					}
				}

			}
		}

		// 处理物品添加或邮件发送
		if err := attacker.checkAddSourceBatch(items); err == nil {
			if m.cfg.Trigger == constant.MonsterTrigerTypeEquipFuDai {
				for _, v := range items {
					if newItem := attacker.addEquipWithQuality(v.ID, m.cfg.Quality, reason); newItem != nil {
						packet.Items = append(packet.Items, newItem)
					}
				}
			} else {
				newItems := attacker.addSourceBatch(items, reason)
				packet.Items = append(packet.Items, newItems...)
			}
		} else {
			attachment := []protocol.Attachment{}
			for _, item := range items {
				if item.ID == constant.ItemTypeCoin {
					attacker.addSource(item.ID, int64(item.Num), reason)
				} else {
					attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
				}
			}
			attacker.SendMail(reason, attachment)
		}

		for _, v := range miningMonsters {
			ownershipUid := attacker.uid
			ownershipType := pb.OwnershipType_OWNERSHIP_PLAYER
			ownershipName := attacker.name

			if attacker.guildId > 0 {
				ownershipUid = attacker.guildId
				ownershipType = pb.OwnershipType_OWNERSHIP_GUILD
				if guild := guildManager.GetGuild(attacker.guildId); guild != nil {
					ownershipName = guild.Guild.Name
				}
			}

			m.desk.addMonsterRelive(v.ID, m.cell.location, 0, 0, ownershipUid, ownershipType, ownershipName)
		}
	}

	// 经验处理
	if monCfg.Experience > 0 {
		num := int64(float64(monCfg.Experience) * (1 + float64(attacker.stat(pb.Attribute_EXP_INCREASE_PERCENTAGE)+multiplier)/BaseRate))
		attacker.addSource(constant.ItemTypeExp, num, reason)
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: constant.ItemTypeExp, Num: int64(num)})
	}

	if len(packet.Items) > 0 {
		attacker.pushMsg(protocol.RouteItemDropPush, packet)
	}

}

// 掉落
// func (m *MonsterInfo) drop(attacker *Player) {
// 	if attacker == nil {
// 		return
// 	}
// 	monCfg := global.GetMonsterCfg(m.cfgId)
// 	if monCfg == nil {
// 		logger.Debugf("获取怪物掉落配置失败 cfgId :%d", m.cfgId)
// 		return
// 	}

// 	packet := &pb.ItemDropPush{Items: []*pb.ItemInfo{}}
// 	packet.Location = &pb.Point{X: m.location.X, Y: m.location.Y}
// 	packet.ObjectId = m.uid
// 	reason := constant.ItemAddDrop

// 	if m.desk.roomCfg.Roomtype == constant.PVEPeacefulMap || m.desk.roomCfg.Roomtype == constant.PVEWarMap ||
// 		m.desk.roomCfg.Roomtype == constant.WorldBossMap || m.desk.roomCfg.Roomtype == constant.StarUpStoneFBMap ||
// 		m.desk.roomCfg.Roomtype == constant.SoulStoneFBMap || m.desk.roomCfg.Roomtype == constant.GuidanceFBMap {

// 		items := []global.IdNum{}

// 		for _, v := range monCfg.Drop {
// 			if len(v.ID) == 0 {
// 				continue
// 			}
// 			r := global.RandIn(int(BaseRate))
// 			ratio := float64(v.Ratio) * (1 + float64(attacker.stat(pb.Attribute_ITEM_DROP_RATE_PERCENTAGE))/BaseRate)
// 			if r < int(ratio) {
// 				//logger.Infof("怪物掉落, UserId:%d, MonsterCfgID: %d , MonsterUid: %d ,Drop: %v", attacker.Uid(), m.cfgId, m.uid, v)
// 				index := global.RandIn(len(v.ID))
// 				if index >= len(v.ID) {
// 					index = 0
// 				}
// 				id := v.ID[index]
// 				num := v.Num.Min + int32(global.RandIn(int(v.Num.Max)-int(v.Num.Min)))

// 				if !itemRecoderMgr.check(attacker.uid, attacker.level, id, num, m.desk.roomCfg.Roomtype) {
// 					logger.Infof("掉落控制 drop UserId:%d, MonsterCfgID: %d ,  ItemID:%d ItemNum:%d", attacker.Uid(), m.cfgId, id, num)
// 					continue
// 				}

// 				items = append(items, global.IdNum{ID: id, Num: num})

// 				packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: id, Num: int64(num)})
// 			}
// 		}

// 		if m.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
// 			reason = constant.ItemAddDropHaiDao
// 		} else if m.desk.roomCfg.Roomtype == constant.StarUpStoneFBMap {
// 			reason = constant.ItemAddDropStarUpStoneFB
// 		} else if m.desk.roomCfg.Roomtype == constant.SoulStoneFBMap {
// 			reason = constant.ItemAddDropSoulStoneFBMap
// 		}

// 		if err := attacker.checkAddSourceBatch(items); err == nil {
// 			attacker.addSourceBatch(items, reason)
// 		} else {
// 			attachment := []protocol.Attachment{}
// 			for _, item := range items {
// 				if item.ID == constant.ItemTypeCoin {
// 					attacker.addSource(item.ID, int64(item.Num), reason)
// 				} else {
// 					attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
// 				}
// 			}
// 			attacker.SendMail(reason, attachment)
// 		}
// 	}

// 	if monCfg.Experience > 0 {
// 		// 经验加成
// 		num := int64(float64(monCfg.Experience) * (1 + float64(attacker.stat(pb.Attribute_EXP_INCREASE_PERCENTAGE))/BaseRate))
// 		attacker.addSource(constant.ItemTypeExp, num, reason)
// 		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: constant.ItemTypeExp, Num: int64(num)})
// 	}

// 	if len(packet.Items) > 0 {
// 		attacker.pushMsg(protocol.RouteItemDropPush, packet)
// 	}
// }

// 掉落
func (m *MonsterInfo) worldboosDrop(attacker *Player) {
	if attacker == nil {
		return
	}
	monCfg := global.GetMonsterCfg(m.cfgId)
	if monCfg == nil {
		logger.Debugf("获取怪物掉落配置失败 cfgId :%d", m.cfgId)
		return
	}

	packet := &pb.ItemDropPush{Items: []*pb.ItemInfo{}}
	packet.Location = &pb.Point{X: m.location.X, Y: m.location.Y}
	packet.ObjectId = m.uid

	items := []global.IdNum{}

	for _, v := range monCfg.Drop {
		if len(v.ID) == 0 {
			continue
		}
		r := global.RandIn(int(BaseRate))
		ratio := float64(v.Ratio) * (1 + float64(attacker.stat(pb.Attribute_ITEM_DROP_RATE_PERCENTAGE))/BaseRate)
		if r < int(ratio) {
			logger.Infof("世界boss怪物掉落, UserId:%d, MonsterCfgID: %d , MonsterUid: %d ,ratio: %v", attacker.Uid(), m.cfgId, m.uid, ratio)
			index := global.RandIn(len(v.ID))
			if index >= len(v.ID) {
				index = 0
			}
			id := v.ID[index]
			num := v.Num.Min + int32(global.RandIn(int(v.Num.Max)-int(v.Num.Min)))

			// if !itemRecoderMgr.check(attacker.uid, attacker.level, id, num, constant.WorldBossMap) {
			// 	logger.Infof("掉落控制 worldboosDrop UserId:%d, MonsterCfgID: %d ,  ItemID:%d ItemNum:%d", attacker.Uid(), m.cfgId, id, num)
			// 	continue
			// }

			if !DropControlInstance.TryDrop(id, int64(num), constant.WorldBossMap, attacker) {
				logger.Debugf("掉落控制过滤 worldboosDrop, UserId:%d, MonsterCfgID: %d, ItemID:%d ItemNum:%d", attacker.Uid(), m.cfgId, id, num)
				continue
			}

			items = append(items, global.IdNum{ID: id, Num: num})

			packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: id, Num: int64(num)})
		}
	}

	if err := attacker.checkAddSourceBatch(items); err == nil {
		attacker.addSourceBatch(items, constant.ItemAddDrop)
	} else {
		attachment := []protocol.Attachment{}
		for _, item := range items {
			if item.ID == constant.ItemTypeCoin {
				attacker.addSource(item.ID, int64(item.Num), constant.ItemAddDrop)
			} else {
				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
			}
		}
		attacker.SendMail(constant.ItemAddDrop, attachment)
	}

	if monCfg.Experience > 0 {
		// 经验加成
		num := int64(float64(monCfg.Experience) * (1 + float64(attacker.stat(pb.Attribute_EXP_INCREASE_PERCENTAGE))/BaseRate))
		attacker.addSource(constant.ItemTypeExp, num, constant.ItemAddDrop)
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: constant.ItemTypeExp, Num: int64(num)})
	}

	if len(packet.Items) > 0 {
		attacker.pushMsg(protocol.RouteItemDropPush, packet)
	}
}

// 重置
func (m *MonsterInfo) reset() {
	m.uid = 0
	m.currentHP = m.GetMaxHP()
	m.displayHP = m.GetMaxHP()
	m.location.X = 0
	m.location.Y = 0
	m.createTime = 0
	m.moveTime = time.Now()
	m.attackTime = time.Now()
	m.searchTime = time.Now()
	m.dead = false
	m.desk = nil
	m.target = nil
	m.roamPath = []*astar.APoint{}
	//m.buffs = make(map[int32]*BuffInfo)
	m.ownershipName = ""
	m.ownershipUid = 0
	m.ownershipType = pb.OwnershipType_OWNERSHIP_NONE
	m.miningPlayer = nil
	m.spawnLocation = global.PointI{}
}

func (m *MonsterInfo) setDesk(d *Desk) {
	m.desk = d
}

// 巡逻
func (m *MonsterInfo) processRoam() {
	if time.Now().Before(m.roamTime) {
		return
	}

	if !m.CanMove() {
		return
	}
	if m.desk == nil {
		return
	}

	if m.target != nil {
		return
	}

	if m.cfg.ViewRange == 0 {
		return
	}

	if m.cfg.MType == MTypeGathering {
		return
	}

	m.moveTo(maputil.RandomNeighbor(m.cell.location.X, m.cell.location.Y))

	m.roamTime = time.Now().Add(time.Second * time.Duration(global.RandIn(90)+1))
}

// 寻找
func (m *MonsterInfo) processSearch_bak() {
	// 嘲讽状态不在攻击其他目标
	if m.getBuffByType(constant.BuffTaunt) != nil && m.target != nil {
		return
	}

	if time.Now().Before(m.searchTime) {
		return
	}
	m.searchTime = time.Now().Add(time.Second)

	if !m.CanAttack() {
		return
	}

	if m.cfg.AttackRange == 0 {
		return
	}

	var distance int
	var closest IMapObject
	// objs := m.desk.aoiMgr.GetPidsByPos(float32(m.location.X), float32(m.location.Y))
	// for _, v := range objs {
	// 	obj := m.desk.objectByUid(int64(v))
	// 	if obj.ObjectType() == ObjectTypeMonster {
	// 		continue
	// 	}
	// 	tmpDistance := maputil.Distance(m.location, obj.Location())
	// 	if tmpDistance > m.cfg.AttackRange {
	// 		continue
	// 	}
	// 	if closest == nil || tmpDistance < distance {
	// 		closest = obj
	// 	}
	// }
	for _, v := range m.desk.objects {
		if v.objectType == ObjectTypeDecoy || v.objectType == ObjectTypePartner {
			if v.dead {
				continue
			}

			tmpDistance := maputil.Distance(m.location, v.Location())
			if tmpDistance > m.cfg.AttackRange {
				continue
			}

			// 优先以战士职业位目标
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

	if closest == nil || !closest.IsWarrior() {
		for _, v := range m.desk.players {
			if v.dead {
				continue
			}
			tmpDistance := maputil.Distance(m.location, v.Location())
			if tmpDistance > m.cfg.AttackRange {
				continue
			}
			if m.desk.roomCfg.Roomtype == constant.PVEWarMap {
				if m.desk.mapInfo.safeArea(v.Location()) {
					continue
				}
			}
			if distance == 0 || tmpDistance < distance || v.IsWarrior() {
				distance = tmpDistance
				closest = v
			}
		}
	}

	if closest != nil {
		m.target = closest
		m.roamPath = []*astar.APoint{}

		// logger.Debugf("M 锁定目标 mUid: %d , 自身坐标:%v, 目标ID: %d  目标坐标: %v  目标类型:%v", m.uid, m.location, m.target.Uid(), m.target.Location(), m.target.ObjectType())
	}
}

// 寻找目标：优先伙伴战士职业，其次战士主角，最后其他
func (m *MonsterInfo) processSearch() {
	// 检查嘲讽状态
	if m.getBuffByType(constant.BuffTaunt) != nil && m.target != nil {
		return
	}

	// 检查搜索冷却时间
	if time.Now().Before(m.searchTime) {
		return
	}
	m.searchTime = time.Now().Add(time.Second)

	// 检查基本攻击条件
	if !m.CanAttack() || m.cfg.AttackRange == 0 {
		return
	}

	if m.target != nil && m.target.Job() == constant.PartnerWarrior {
		return
	}

	// 定义目标优先级
	const (
		PriorityPartnerWarrior = iota + 1
		PriorityPlayerWarrior
		PriorityOther
	)

	var (
		bestTarget   IMapObject
		bestPriority int
		bestDistance int
	)

	// 获取攻击范围内的所有格子
	attackCells := m.desk.mapInfo.getSquareCells(&m.cell.location, int32(m.cfg.AttackRange))

	// 遍历所有可能的目标
	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}

		for _, obj := range cell.objects {
			if obj.Dead() {
				continue
			}
			// 跳过无效对象类型
			switch obj.ObjectType() {
			case ObjectTypeMonster, ObjectTypeDecoy, ObjectTypeBuff:
				continue
			}

			// 计算距离
			distance := maputil.Distance(m.location, obj.Location())

			// 确定目标优先级
			priority := PriorityOther
			if obj.ObjectType() == ObjectTypePartner && obj.Job() == constant.PartnerWarrior {
				priority = PriorityPartnerWarrior
			} else if obj.ObjectType() == ObjectTypePlayer && obj.Job() == constant.Warrior {
				priority = PriorityPlayerWarrior
			}

			// 比较目标优先级和距离
			if bestTarget == nil ||
				priority < bestPriority ||
				(priority == bestPriority && distance < bestDistance) {
				bestTarget = obj
				bestPriority = priority
				bestDistance = distance
			}
		}
	}

	// 更新目标
	if bestTarget != nil {
		if m.target != nil && m.target.Uid() != bestTarget.Uid() {
			logger.Debugf("%d:%s 切换目标 原目标: %d %s 新目标: %d %s",
				m.uid, m.name, m.target.Uid(), m.target.Name(), bestTarget.Uid(), bestTarget.Name())
		}
		m.target = bestTarget
		m.roamPath = []*astar.APoint{}
	}
}

// 目标
func (m *MonsterInfo) processTarget() {
	if m.desk == nil {
		return
	}

	if m.target == nil {
		return
	}

	if m.target.Dead() {
		m.target = nil
		return
	}

	if !m.CanAttackTarget(m.target) {
		return
	}
	if m.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if m.desk.mapInfo.safeArea(m.target.Location()) {
			m.target = nil
			return
		}
	}

	if !m.inAttackRange(m.cfg.AttackDistance) {
		if !m.CanMove() {
			return
		}
		if maputil.Distance(m.location, &pb.Point{X: float64(m.spawnLocation.X), Y: float64(m.spawnLocation.Y)}) > int(m.cfg.Chase) {
			if m.desk.roomCfg.Roomtype == constant.PVEWarMap && (m.cfg.MType == MTypeCommonBoss || m.cfg.MType == MTypeWarBoss) {
				m.target = nil
				if cell := m.desk.mapInfo.getAvailableCell(&m.spawnLocation, 2); cell != nil {
					if m.desk.mapInfo.locationChange(m, &pb.Point{X: float64(cell.location.X), Y: float64(cell.location.Y)}) {
						m.desk.Broadcast(protocol.RouteTeleportPush, &pb.TeleportPush{Uid: m.uid, Location: m.location, Direction: m.direction})
						logger.Debugf("回到原点 起点=%v  终点=%v", m.location, cell.location)
					} else {
						logger.Debugf("回到原点寻路失败 起点=%v  终点=%v", m.location, m.target.Location())
					}
				}
			} else if m.cfg.MType == MTypeNormal {
				m.target = nil
				cell := m.desk.mapInfo.getAvailableCell(&m.spawnLocation, 2)

				if cell != nil {
					if path := maputil.GetPath(m.location, &pb.Point{X: float64(cell.location.X), Y: float64(cell.location.Y)}, &m.desk.mapCfg.AstarMgr); len(path) > 1 {
						m.roamPath = path[1:]
						// logger.Debugf("回到原点 起点=%v  终点=%v", m.location, cell.location)
					} else {
						logger.Debugf("回到原点寻路失败 起点=%v  终点=%v", m.location, m.target.Location())
					}
				}
			}

		} else {
			if path := maputil.GetPath(m.location, m.target.Location(), &m.desk.mapCfg.AstarMgr); len(path) > 1 {
				m.moveTo(path[1])
			} else {
				logger.Debugf("寻路失败 起点=%v  终点=%v", m.location, m.target.Location())
			}
		}
	} else {
		if len(m.cell.objects) > 1 {
			cell := m.desk.mapInfo.getAvailableCell(global.ToPointI(m.target.Location()), int32(m.cfg.AttackDistance))
			if cell != nil {
				if path := maputil.GetPath(m.location, &pb.Point{X: float64(cell.location.X), Y: float64(cell.location.Y)}, &m.desk.mapCfg.AstarMgr); len(path) > 1 {
					m.moveTo(path[1])
					return
				}
			} else {
				m.attack()
			}
		}

		m.attack()
	}
}

func (m *MonsterInfo) processMining() {
	if m.desk == nil || m.dead {
		return
	}
	if m.ownershipType == pb.OwnershipType_OWNERSHIP_NONE {
		return
	}

	if time.Now().After(m.ownershipStart.Add(time.Duration(global.GlobalCfg.War.Protect.SafyTime) * time.Second)) {
		m.ownershipType = pb.OwnershipType_OWNERSHIP_NONE
		m.ownershipName = ""
		m.ownershipUid = 0

		m.desk.Broadcast(protocol.RouteMonsterUpdatePush, &pb.MonsterUpdatePush{
			Monster: m.fillPBMonster(),
		})

		m.desk.logger.Debugf(" 采集怪归属时间到， 变无归属  muid:name :%d:%s  m.ownershipType : %d", m.uid, m.name, m.ownershipType)
	}
}

// 召唤
func (m *MonsterInfo) processSummon() {
	//召唤怪物必须有召唤技能ID
	if m.linkedSkill == 0 {
		return
	}
	if time.Now().After(m.delTime) {
		//logger.Debugf("召唤怪死亡， 召唤技能失效  muid:name :%d:%s  m.linkedSkill : %d", m.uid, m.name, m.linkedSkill)

		m.linkedSkill = 0
		m.Die(nil)
		return
	}
}

func (m *MonsterInfo) fillPBMonster() *pb.MonsterInfo {
	mm := pb.MonsterInfo{
		CfgId:     m.cfgId,
		Uid:       m.uid,
		Location:  m.location,
		Direction: m.direction,
		CurHp:     m.currentHP,
		MaxHp:     m.GetMaxHP(),
		Ownership: &pb.Ownership{},
	}
	mm.Ownership.Type = m.ownershipType
	mm.Ownership.Name = m.ownershipName
	mm.Ownership.Uid = m.ownershipUid
	return &mm
}

// 移动
func (m *MonsterInfo) moveTo(point *astar.APoint) bool {
	if point == nil {
		return false
	}

	if m.desk == nil {
		return false
	}

	if time.Now().Before(m.attackTime) {
		return false
	}

	if !m.CanMove() {
		return false
	}

	if m.cfg.MType == MTypeGathering {
		return false
	}

	to := &pb.Point{X: float64(point.X), Y: float64(point.Y)}
	if !m.desk.mapInfo.locationChange(m, to) {
		return false
	}

	m.moveTime = time.Now().Add(time.Duration(m.getMoveSpeed()) * time.Millisecond)

	packet := &pb.MovePush{Uid: m.uid, Location: m.location, Direction: m.direction}
	// m.boardcastMsg(protocol.RouteMovePush, packet)

	m.desk.Broadcast(protocol.RouteMovePush, packet)

	//logger.Debugf("M 移动 mUid: %d:%s   目标位置: %v 移动速度 %d", m.uid, m.name, to, m.getMoveSpeed())
	return true
}

// 攻击
func (m *MonsterInfo) attack() {
	if !m.CanAttack() {
		return
	}

	if len(m.cfg.Skills) == 0 {
		return
	}

	if !m.thresholdDetect() {
		return
	}

	skillIndex := 0
	castSkillId := m.GetNormalSkillID()

	for i, v := range m.skills {
		if v.CfgId == m.GetNormalSkillID() {
			continue
		}
		sc, err := global.GetSkillCfg(v.CfgId)
		if err != nil {
			continue
		}
		if sc.ActiveorPassive == PassiveSkill {
			continue
		}

		if time.Now().After(m.skills[i].Cooldown) {
			castSkillId = v.CfgId
			skillIndex = i
			break
		}
	}

	if castSkillId == 0 {
		return
	}

	skillCfg, err := global.GetSkillCfg(int32(castSkillId))
	if err != nil {
		logger.Errorf("M 攻击 获取技能配置失败 mUid: %d, skillId: %d", m.uid, castSkillId)
		return
	}

	attackCells := []*Cell{}
	targets := []IMapObject{}
	targets = append(targets, m.target)
	switch skillCfg.SkillID {
	case constant.WorldBossSkill1, constant.ThorSkill1:
		attackCells = m.desk.mapInfo.getSquareCells(&m.cell.location, skillCfg.Radius)
	case constant.HeshangSkill1:
		attackCells = []*Cell{}
		targets = []IMapObject{}
		m.heshangSkill1Begin(castSkillId)
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
				if m.CanAttackTarget(obj) && m.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == m.uid {
			continue
		}

		if t.ObjectType() == ObjectTypeBuff || t.ObjectType() == ObjectTypeMonster {
			continue
		}

		// 安全区内不能攻击玩家
		if m.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (m.desk.mapInfo.safeArea(m.location) || m.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		// 技能距离
		if skillCfg.SType != constant.SkillTypeNormal && skillCfg.MaxDistance > 0 && maputil.Distance(m.location, t.Location()) > int(skillCfg.MaxDistance) {
			continue
		}

		// 伤害延迟生效
		mdesk := m.desk
		mt := t
		ilocation := &m.cell.location
		time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
			if mdesk != nil {
				mdesk.addDelayAction(ObjectTypeMonster, m.uid, castSkillId, mt, m.direction, ilocation)
			}
		})
		//logger.Debugf("MM 攻击 attack mUid:Name  %d:%s , 技能 %d 目标: %d:%s, 坐标 %v  目标坐标 %v ", m.uid, m.name, castSkillId, t.Uid(), t.Name(), m.location, t.Location())
	}

	m.setAttackTime(global.GetActionsetCfg(skillCfg.SkillID))

	m.direction = maputil.DirectionFromPointF(m.location, m.target.Location())

	packet := &pb.AttackPush{
		AttackerId: m.uid,
		TargetUids: []int64{m.target.Uid()},
		SkillId:    castSkillId,
		Location:   m.target.Location(),
		Direction:  m.direction,
	}
	m.desk.Broadcast(protocol.RouteAttackPush, packet)

	m.SkillBuffTrigger(m.skills[skillIndex])

	// 技能CD
	if castSkillId != m.GetNormalSkillID() && skillCfg.CoolDown > 0 {
		m.skills[skillIndex].Cooldown = time.Now().Add(time.Duration(skillCfg.CoolDown) * time.Millisecond)
	}
}

func (m *MonsterInfo) processAI() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	m.processRoam()
	m.processSearch()
	m.processTarget()
	m.processHp()
	m.processBuff()
	m.processMining()
	m.processSummon()
}

///////////////////////////////////////////////////////////////////////////////////

type MonsterMgr struct {
	sync.Mutex
	increment int64

	monsterPool map[int32][]*MonsterInfo // 怪物回收池 key:怪物配置ID
}

var monsterMgr *MonsterMgr = NewMonsterMgr()

func NewMonsterMgr() *MonsterMgr {
	return &MonsterMgr{increment: 1, monsterPool: map[int32][]*MonsterInfo{}}
}

// 产怪
func (mgr *MonsterMgr) CreateMonsterByLocation(cfgId int32, l *global.PointI) *MonsterInfo {
	mgr.increment++
	cfg := global.GetMonsterCfg(cfgId)
	if cfg == nil {
		return nil
	}

	m := newMonster(cfgId)

	// rx := rand.Intn(cfg.ViewRange)
	// ry := rand.Intn(cfg.ViewRange)
	// m.location = &pb.Point{X: float64(l.X) + float64(rx), Y: float64(l.Y) + float64(ry)}

	m.bornLocation = *l
	return m
}

func (mgr *MonsterMgr) CreateMonster(cfgId int32) *MonsterInfo {
	mgr.increment++

	m := newMonster(cfgId)
	return m
}

// 批量产怪
func (mgr *MonsterMgr) CreateMonsterBatch(cfgId, num int32, l *global.PointI) []*MonsterInfo {
	mgr.Lock()
	defer mgr.Unlock()

	monlist := []*MonsterInfo{}
	for i := 0; i < int(num); i++ {
		m := mgr.CreateMonsterByLocation(cfgId, l)

		monlist = append(monlist, m)
	}

	return monlist
}

// 归还到怪物池
func (mgr *MonsterMgr) GiveBackMonster(mon *MonsterInfo) {
	mgr.Lock()
	defer mgr.Unlock()

	_, ok := mgr.monsterPool[mon.cfgId]
	if ok {
		// 每个类型的怪最多缓存1000只
		if len(mgr.monsterPool[mon.cfgId]) >= 1000 {
			logger.Debugf("归还到怪物池失败  Id : %d , 此类怪物剩余数量：%d", mon.cfgId, len(mgr.monsterPool[mon.cfgId]))
			return
		}
		mgr.monsterPool[mon.cfgId] = append(mgr.monsterPool[mon.cfgId], mon)
	} else {
		mgr.monsterPool[mon.cfgId] = []*MonsterInfo{}
		mgr.monsterPool[mon.cfgId] = append(mgr.monsterPool[mon.cfgId], mon)
	}

	//logger.Debugf("归还怪物池成功  CfgId : %d Name : %s , 此类怪物剩余数量：%d", mon.cfgId, mon.cfg.Name, len(mgr.monsterPool[mon.cfgId]))
}

// 获取空闲怪物
func (mgr *MonsterMgr) TakeMonster(cfgId int32) (*MonsterInfo, error) {
	mgr.Lock()
	defer mgr.Unlock()

	var newMon *MonsterInfo
	_, ok := mgr.monsterPool[cfgId]
	if ok {
		if len(mgr.monsterPool[cfgId]) > 0 {
			newMon = mgr.monsterPool[cfgId][0]
			mgr.monsterPool[cfgId] = mgr.monsterPool[cfgId][1:]

			//logger.Debugf("取出怪物池成功  CfgId : %d Name : %s , 此类怪物剩余数量：%d", cfgId, newMon.cfg.Name, len(mgr.monsterPool[cfgId]))
		}
	}
	if newMon == nil {
		logger.Debugf("怪物池已空,创建新的怪物实例  CfgId : %d  ", cfgId)
		newMon = mgr.CreateMonster(cfgId)
	}
	if newMon != nil {
		return newMon, nil
	} else {
		return nil, errutil.ErrCreateMonter
	}
}

func (mgr *MonsterMgr) MonsterDorp(cfgId int32, num int32) []global.IdNum {
	items := []global.IdNum{}

	monCfg := global.GetMonsterCfg(cfgId)
	if monCfg == nil {
		logger.Debugf("获取怪物掉落配置失败 cfgId :%d", cfgId)
		return items
	}
	for i := int32(0); i < num; i++ {
		for _, v := range monCfg.Drop {
			if len(v.ID) == 0 {
				continue
			}
			r := global.RandIn(int(BaseRate))
			ratio := float64(v.Ratio)
			if r < int(ratio) {

				index := global.RandIn(len(v.ID))
				if index >= len(v.ID) {
					index = 0
				}
				id := v.ID[index]
				num := v.Num.Min + int32(global.RandIn(int(v.Num.Max)-int(v.Num.Min)))

				items = append(items, global.IdNum{ID: id, Num: num})
			}
		}
	}
	return items
}
