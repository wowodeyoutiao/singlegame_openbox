package game

import (
	"context"
	"nanoserver/global"
	"nanoserver/pkg/aoi9"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/goredisx"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const chLen int = 1024

// World interface
type MowWorld struct {
	gameType protocol.GameType
}

func NewMowWorld() *MowWorld {
	return &MowWorld{gameType: protocol.GameTypeMow}
}

func (w *MowWorld) NewDesk(roomNo int64, roomID int32) *Desk {
	d := NewBaseDesk(roomNo, roomID)
	d.gameLogic = &Mow{desk: d}

	roomCfg := global.GetRoomCfg(roomID)
	if roomCfg == nil {
		logger.Errorf("创建房间失败,roomcfg 不存在 RoomId: %d", roomID)
		return nil
	}
	mapCfg := global.GetMapCfg(roomCfg.MapId)
	if mapCfg == nil {
		logger.Errorf("创建房间失败,mapCfg 不存在 MapId: %d", roomCfg.MapId)
		return nil
	}

	d.astarMgr = &mapCfg.AstarMgr

	gh := mapCfg.High / 10
	gw := mapCfg.Wide / 10
	d.aoiMgr = aoi9.NewAOIManager(0, mapCfg.Wide, gh, 0, mapCfg.High, gw)
	d.mapId = mapCfg.ID
	d.mapCfg = mapCfg
	d.roomCfg = roomCfg
	d.mapInfo = InitMapInfo(mapCfg)

	d.gameLogic.OnCreate()

	return d
}

func (w *MowWorld) GetName() string {
	return "Mow"
}

// ****************************************************************************************************

type Mow struct {
	desk *Desk

	additionList map[int32]interface{} //房间自定义参数
}

// IGameLogic interface Begin
func (mow *Mow) OnCreate() {
	mow.desk.chMove = make(chan *MoveRequest, chLen)

	mow.desk.chAttack = make(chan *AttackRequest, chLen)

	mow.desk.chMining = make(chan *MiningRequest, chLen)

	mow.desk.chDelayAttack = make(chan *DelayAction, chLen)
	mow.desk.chMonsterRelive = make(chan *MonRelive, chLen)
	mow.desk.chPlayerRelive = make(chan *PlayerRelive, chLen)
	mow.desk.chJoin = make(chan int64, chLen)
	mow.desk.chLeave = make(chan *Leave, chLen)
	mow.desk.chSkipBattle = make(chan int64, chLen)
}

func (mow *Mow) OnEnter(p *Player) {

	mow.desk.serverID = p.serverID
	p.setPlayerStatus(PlayerStatusPlaying)
	p.stats.changeStats(pb.Attribute_HEALTH, p.GetMaxHP())
	p.currentHP = p.GetMaxHP()
	p.displayHP = p.currentHP
	p.reviveTimes = 0
	if p.isRobot() {
		robotMgr.onEnterDesk(p.serverID, p.Uid())
	}
}

func (mow *Mow) OnExit(p *Player) {

	if p.isRobot() {
		robotMgr.onExitDesk(p.serverID, p.Uid())
	} else {
		p.setPlayerStatus(PlayerStatusIdle)
		p.asyncSaveDb()
	}
	p.stats.setStats(pb.Attribute_HEALTH, p.GetMaxHP())
	p.currentHP = p.GetMaxHP()
	p.displayHP = p.currentHP
	p.miningBreak()
	//删除地图数据
	mow.desk.mapInfo.removeObject(p, p.GetCell().location)

	for _, v := range p.partnerSlots {
		if v == nil {
			continue
		}
		if v.PartnerUid == 0 {
			continue
		}
		//删除地图数据
		mow.desk.delObject(v.PartnerUid)

		if par := p.getPartnerInfoBySlotId(v.SlotId); par != nil {
			par.stats.setStats(pb.Attribute_HEALTH, par.GetMaxHP())
			par.currentHP = par.GetMaxHP()
			par.displayHP = par.currentHP
			par.target = nil
			par.desk = nil
		}
	}

	delete(mow.desk.ConsecutiveKills, p.uid)
}

func (mow *Mow) OnStart() error {

	roomCfg := global.GetRoomCfg(mow.desk.cfgID)
	if roomCfg == nil {
		return errutil.ErrConfigNotFound
	}

	switch mow.desk.roomCfg.Roomtype {
	case constant.SoloTowerMap:
		mow.desk.logger.Debugln("玩家进入无尽试炼，开始刷怪 ")
		if stcfg := global.GetEndlesstrialsCfg(mow.desk.SoloTower.Layers); stcfg != nil {
			for _, v := range stcfg.Monster {
				monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)
				//mow.desk.logger.Debugf("爬塔地【%d】层 ，开始刷怪 %+v ", mow.desk.SoloTowerCfg.Layers, v)
				for _, mon := range monlist {
					if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
						if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {
							mon.setDesk(mow.desk)
							mow.desk.addMonster(mon)
						} else {
							mow.desk.logger.Debugf("爬塔地【%d】层 ，刷怪失败 ,未找到出生点 ", mow.desk.SoloTowerCfg.Layers)
						}
					}
				}
			}
		} else {
			return errutil.ErrConfigNotFound
		}
		mow.desk.logger.Debugln("无尽试炼刷怪结束")

	case constant.PVEPeacefulMap, constant.GuidanceFBMap:
		mow.desk.logger.Debugln("第一个玩家进入PVE房间,开始刷怪 ")
		if len(mow.desk.mapCfg.MonsterIds) == 0 {
			return errutil.ErrConfigNotFound
		}
		for _, v := range mow.desk.mapCfg.MonsterIds {

			monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)

			for _, mon := range monlist {

				if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
					if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {
						mon.setDesk(mow.desk)
						mow.desk.addMonster(mon)
					}
				}
			}
		}
		mow.desk.logger.Debugln("PVE刷怪结束")

	case constant.PVEWarMap:
		mow.desk.logger.Debugln("第一个玩家进入PVE房间,开始刷怪 ")
		if len(mow.desk.mapCfg.MonsterIds) == 0 {
			return errutil.ErrConfigNotFound
		}
		for _, v := range mow.desk.mapCfg.MonsterIds {

			monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)

			for _, mon := range monlist {
				if mon.cfg.MType != MTypeGathering && mon.cfg.Delay > 0 {
					refreshTime := time.Now().Add(time.Duration(mon.cfg.Delay) * time.Millisecond)

					mow.desk.refreshList = append(mow.desk.refreshList, &MonsterRefresh{
						CfgId:       v.ID,
						RefreshTime: refreshTime,
						Location:    v.Point})

					//mow.desk.logger.Debugf("增加延迟刷新boss  mcfgid %d delay %v ", v.ID, refreshTime)

					continue
				}

				if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
					if mow.desk.mapInfo.safeArea(global.ToPointF(&cell.location)) {
						logger.Warningf("刷怪失败, 位置在安全区, mcfgid %d  坐标 %v    Cell.location %v", mon.cfg.ID, mon.location, cell.location)
						continue
					}
					if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {

						mon.setDesk(mow.desk)
						mow.desk.addMonster(mon)
					}
				}

				// 给怪物加上永久增益buff，从redis获取
				for _, entity := range global.GlobalCfg.War.Monstergrow.MonsterAttr {
					if mon.cfg.ID == entity.ID {
						key := "pve_war_map:boss_die:" + strconv.FormatInt(int64(mow.desk.zoneID), 10) + ":" + strconv.FormatInt(int64(entity.ID), 10)
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
						bufStats.Attrs[pb.Attribute_HEALTH_PERCENTAGE] = times * int32(v.Num)
						bufStats.Attrs[pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE] = times * int32(v.Num)
						bufStats.Attrs[pb.Attribute_MAGIC_DEFENSE_PERCENTAGE] = times * int32(v.Num)
						mon.AddBuff(101, mon, 1000*60*60*24, bufStats, 0, 0)
					}
				}
			}
		}
		mow.desk.logger.Debugln("PVE刷怪结束")

	case constant.WorldBossMap:
		if len(mow.desk.mapCfg.MonsterIds) == 0 {
			return errutil.ErrConfigNotFound
		}

		if curBoss := worldBossMgr.getBoss(mow.desk.serverID); curBoss != nil {
			if boss := monsterMgr.CreateMonsterByLocation(curBoss.cfgId, &curBoss.bornLocation); boss != nil {
				boss.stats.changeStats(pb.Attribute_HEALTH, curBoss.stat(pb.Attribute_HEALTH))
				boss.currentHP = boss.stat(pb.Attribute_HEALTH)
				boss.displayHP = boss.currentHP
				boss.serverID = mow.desk.serverID
				if cell := mow.desk.mapInfo.getAvailableCell(&curBoss.bornLocation, int32(boss.cfg.DistributeRange)); cell != nil {
					if boss.cell = mow.desk.mapInfo.addObject(boss, &cell.location); boss.cell != nil {
						boss.SettDirection(pb.Direction_DOWN)
						boss.setDesk(mow.desk)
						mow.desk.addMonster(boss)

						mow.desk.nextBuffTime = time.Now().Add(time.Duration(global.GlobalCfg.WorldBoss.Accommodation.Num) * time.Second)
					}
				}
			}
		} else {
			return errutil.ErrConfigNotFound
		}

		mow.desk.logger.Debugln("世界boss刷怪结束")

	case constant.ExperienceFBMap, constant.GoldFBMap, constant.EnhanceStoneFBMap, constant.BookPageFBMap, constant.StarUpStoneFBMap, constant.SoulStoneFBMap:
		mow.desk.DDailyMap.mingCount = 0
		if len(mow.desk.mapCfg.MonsterIds) == 0 {
			return errutil.ErrConfigNotFound
		}
		for _, v := range mow.desk.mapCfg.MonsterIds {
			monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)

			for _, mon := range monlist {
				if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
					if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {
						mon.setDesk(mow.desk)
						mow.desk.addMonster(mon)
					}
				}
			}
		}
		mow.desk.logger.Debugln("每日必做副本刷怪结束")

	case constant.BattleMap, constant.BossMap, constant.TrialTowerMap, constant.RelicMap:
		mow.desk.logger.Debugln("玩家进入副本，开始刷怪 ")
		if btcfg := global.GetBattleMapCfg(mow.desk.BattleId); btcfg != nil {
			for _, v := range btcfg.Monster {
				monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)
				//mow.desk.logger.Debugf("爬塔地【%d】层 ，开始刷怪 %+v ", mow.desk.SoloTowerCfg.Layers, v)
				for _, mon := range monlist {
					if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
						if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {
							mon.setDesk(mow.desk)
							mow.desk.addMonster(mon)
						} else {
							mow.desk.logger.Debugf("战役副本【%d】 ，刷怪失败 ,未找到出生点 ", mow.desk.SoloTowerCfg.Layers)
						}
					}
				}
			}
		} else {
			return errutil.ErrConfigNotFound
		}
	case constant.GuildMap:
		mow.desk.logger.Debugln("玩家进入行会副本，开始刷怪 ")
		beginId := constant.GuildMap*100000 + 1
		mow.desk.BattleId = beginId
		if btcfg := global.GetBattleMapCfg(beginId); btcfg != nil {
			for _, v := range btcfg.Monster {
				monlist := monsterMgr.CreateMonsterBatch(v.ID, v.Num, &v.Point)
				for _, mon := range monlist {
					if cell := mow.desk.mapInfo.getAvailableCell(&mon.bornLocation, int32(mon.cfg.DistributeRange)); cell != nil {
						if mon.cell = mow.desk.mapInfo.addObject(mon, &cell.location); mon.cell != nil {
							mon.setDesk(mow.desk)
							mow.desk.addMonster(mon)
						} else {
							mow.desk.logger.Debugf("行会副本【%d】 ，刷怪失败 ,未找到出生点 ", mow.desk.SoloTowerCfg.Layers)
						}
					}
				}
			}

			for _, v := range btcfg.Settlement {
				mow.desk.GuildReward.NexReward = append(mow.desk.GuildReward.NexReward, &pb.KeyValue{K: v.ID, V: v.Num})
			}

			mow.desk.logger.Debugln("行会副本刷怪结束")

		} else {
			return errutil.ErrConfigNotFound
		}

	case constant.FightEndMap:
		if btcfg := global.GetBattleMapCfg(mow.desk.BattleId); btcfg != nil {
			for _, v := range btcfg.Monster {
				if robot := robotMgr.getFreeRobotByCfgId(v.ID); robot != nil {
					robot.refreshStats()
					robot.dead = false
					robot.currentHP = robot.GetMaxHP()
					robot.displayHP = robot.currentHP

					logger.Debugf("一战到底对手 %d:%s Uid: %d  ", robot.cfg.UUID, robot.cfg.Name, robot.uid)
					robot.target = nil
					for _, v := range robot.partners {
						v.target = nil
					}
					mow.desk.robotJoin(robot, mow.desk.BattleId)
				}
			}

			mow.desk.logger.Debugln("一战到底副本刷怪结束")
		} else {
			return errutil.ErrConfigNotFound
		}
	}
	return nil
}

func (mow *Mow) OnGameOver() {
	mow.desk.reset()
}

// IGameLogic Interface End

func (mow *Mow) GetAddition(name int32) interface{} {
	unit, ok := mow.additionList[name]
	if ok {
		return unit
	}
	return nil
}

func (mow *Mow) SetAddition(name int32, v interface{}) {
	mow.additionList[name] = v
}
