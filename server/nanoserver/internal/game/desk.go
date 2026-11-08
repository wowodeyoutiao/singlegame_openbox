package game

import (
	"fmt"
	"math"
	"nanoserver/pkg/event"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/aoi9"
	"nanoserver/pkg/astar"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/randx"
	"nanoserver/pkg/utils"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano"
	"github.com/lonng/nano/scheduler"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	LoopTick = 5
)

type IGameLogic interface {
	OnCreate()
	OnStart() error
	OnGameOver()
	OnEnter(p *Player)
	OnExit(p *Player)
}

type DelayAction struct {
	ObjType   ObjectType
	ObjUid    int64
	SkillId   int32
	Target    IMapObject
	Direction pb.Direction
	Location  global.PointI
}

type AttackRequest struct {
	Uid int64
	Req *pb.AttackRequest
}

type MoveRequest struct {
	Uid      int64
	Location *pb.Point
}

type MiningRequest struct {
	Uid       int64
	TargetUid int64
}

type SoloTower struct {
	Layers       int32
	SoloTowerCfg *global.EndlesstrialsConfig
}

type MonRelive struct {
	CfgID         int32
	Location      global.PointI
	OwnershipType pb.OwnershipType // 怪物归属类型
	OwnershipUid  int64            // 怪物归属ID
	OwnershipName string           // 归属名字
	Lifespan      int64            // 怪物存活时间(毫秒)
	linkedSkill   int32            // 怪物关联技能
}

type PlayerRelive struct {
	uid        int64
	reliveType int32 //0:表示复活点复活, 1:表示原地复活
}

type ReliveData struct {
	reliveType pb.ReliveType
	uid        int64
	objectType ObjectType
	reviveHp   int32
}

type DWorldBoss struct {
	// hpChangeChannel    chan int32
	nextBuffTime   time.Time
	curRoundDamage int64 // 本局伤害
}

type DDailyMap struct {
	round           int32
	scoring         global.XMLKeyValue
	totalDamage     int32
	reward          global.IdNum
	monBuffStatus   *Stats
	buffRefreshTime time.Time
	mingCount       int32 // 挖矿次数
}

type BattleMap struct {
	GuildId           int64
	BattleId          int32
	BattleCfg         *global.BattleConfig
	GuildReward       *pb.GuildCopyRewardPush
	gr_lock           sync.Mutex
	GuildRewardMember map[int64]int64
	RobotId           int64
}

type MonsterRefresh struct {
	CfgId       int32
	RefreshTime time.Time
	Location    global.PointI
}

type Kills struct {
	killNum  int32
	lastKill time.Time
}

type PveWarMap struct {
	ConsecutiveKills map[int64]*Kills
	refreshList      []*MonsterRefresh
	monsterDeadTimes map[int32]int32
}

type Leave struct {
	Uid        int64
	Disconnect bool
}

type Desk struct {
	lock      sync.Mutex
	gameLogic IGameLogic

	cfgID     int32 // 配表ID
	deskID    int64 // 唯一ID
	mod       int32 // 房间类型
	mapId     int32 // 地图ID
	state     constant.DeskStatus
	createdAt int64 // 创建时间
	players   []*Player
	robots    []*RobotInfo
	monsters  map[int64]*MonsterInfo
	objects   []*MapObject
	beginTime time.Time
	serverID  int32 // 房间归属服务器
	zoneID    int32 // 房间归属战区
	isInitMon bool

	group *nano.Group // 组播通道
	die   chan struct{}

	chJoin chan int64

	chLeave chan *Leave

	chMove chan *MoveRequest // 移动

	chAttack chan *AttackRequest

	chMining chan *MiningRequest

	chDelayAttack chan *DelayAction

	chMonsterRelive chan *MonRelive

	chPlayerRelive chan *PlayerRelive

	chSkipBattle chan int64

	eventMgr *event.EventManager

	astarMgr *astar.Map       // 寻路
	aoiMgr   *aoi9.AOIManager // 视野
	mapCfg   *global.MapConfig
	roomCfg  *global.RoomConfig
	mapInfo  *MapInfo

	// pvp房间
	matchUid int64
	*MatchInfo
	// 无尽试炼房间
	*SoloTower
	*DWorldBoss
	*DDailyMap
	*BattleMap
	*PveWarMap
	*FateCoin

	logger *log.Entry
}

func NewBaseDesk(deskID int64, roomID int32) *Desk {
	d := &Desk{
		cfgID:     roomID,
		deskID:    deskID,
		mod:       protocol.DeskModNomal,
		state:     constant.DeskStatusCreate,
		players:   []*Player{},
		monsters:  map[int64]*MonsterInfo{},
		SoloTower: &SoloTower{},

		group: nano.NewGroup(uuid.New()),

		die:        make(chan struct{}),
		logger:     log.WithField(fieldDesk, deskID),
		DWorldBoss: &DWorldBoss{},
		DDailyMap:  &DDailyMap{},
		BattleMap:  &BattleMap{GuildReward: &pb.GuildCopyRewardPush{}, GuildRewardMember: map[int64]int64{}},

		PveWarMap: &PveWarMap{ConsecutiveKills: map[int64]*Kills{}, refreshList: []*MonsterRefresh{}, monsterDeadTimes: map[int32]int32{}},
		FateCoin:  NewFateCoin(),
		eventMgr:  event.NewEventManager(),
	}

	// 注册事件处理器
	d.registerEventHandlers()
	return d
}

func (d *Desk) getTargets(effectNum int32, p *Player) []IMapObject {
	var targets []IMapObject
	switch effectNum {
	case constant.SkillEffectnumSlef:
		targets = append(targets, p)
	case constant.SkillEffectnumAlliance:
		targets = p.Friends()
	case constant.SkillEffectnumEnemyAll:
		targets = d.getMonsters()
	case constant.SkillEffectnumPartner:
		for _, v := range p.partnerSlots {
			if v == nil {
				continue
			}
			if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
				targets = append(targets, partner)
			}
		}
	}
	return targets
}

// 广播消息
func (d *Desk) Broadcast(route string, packet interface{}) {
	if d.group == nil {
		return
	}

	if err := d.group.Broadcast(route, packet); err != nil {
		logger.Debugln(stack())
		logger.Errorf("Session push message error, Error=%s", err.Error())
	}
}

// 房间容量
func (d *Desk) totalPlayerCount() int {
	return int(roomManager.getRoom(d.cfgID).opt.PlayerNum)
}

// 玩家数量
func (d *Desk) realPlayerCount() int {
	var count int
	for _, v := range d.players {
		if !v.isRobot() {
			count++
		}
	}
	return count
}

// 同步信息
func (d *Desk) syncDeskStatus() {
	packet := pb.PlayerEnterDeskPush{Data: []*pb.PlayerInfo{}}

	for _, p := range d.players {
		if p.curState != PlayerStatusPlaying || p.session == nil {
			d.logger.Warningf("玩家状态异常,玩家id=%d, 玩家状态=%d", p.uid, p.curState)
			continue
		}

		guildname := ""
		if guild := guildManager.GetGuild(p.guildId); guild != nil {
			guildname = guild.Guild.Name
		}
		packet.Data = append(packet.Data,
			&pb.PlayerInfo{
				Uid:         p.Uid(),
				Nickname:    p.name,
				HeadUrl:     "",
				MaxHp:       p.stat(pb.Attribute_MAX_HEALTH),
				CurHp:       p.currentHP,
				Job:         int32(p.job),
				Gender:      int32(p.sex),
				Level:       p.level,
				Location:    &pb.Point{X: p.location.X, Y: p.location.Y},
				Direction:   p.direction,
				BuffList:    p.getPBBuffList(),
				PartnerList: p.getScenePartner(),
				GuildId:     p.guildId,
				GuildName:   guildname,
				TitleId:     p.equipTitleId,
				OfficialExp: p.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
			})

	}

	for _, r := range d.robots {
		packet.Data = append(packet.Data,
			&pb.PlayerInfo{
				Uid:         r.Uid(),
				Nickname:    r.name,
				HeadUrl:     "",
				MaxHp:       r.stat(pb.Attribute_MAX_HEALTH),
				CurHp:       r.currentHP,
				Job:         int32(r.job),
				Gender:      int32(r.sex),
				Level:       r.level,
				Location:    r.location,
				Direction:   r.direction,
				BuffList:    r.getPBBuffList(),
				PartnerList: r.getScenePartner(),
				GuildId:     0,
				GuildName:   "",
			})
	}

	d.Broadcast(protocol.RoutePlayerEnterDeskPush, &packet)
}

func (d *Desk) setStatus(s constant.DeskStatus) {
	atomic.StoreInt32((*int32)(&d.state), int32(s))
}

func (d *Desk) status() constant.DeskStatus {
	return constant.DeskStatus(atomic.LoadInt32((*int32)(&d.state)))
}

func (d *Desk) playerWithId(uid int64) (*Player, error) {
	for _, p := range d.players {
		if p.uid == uid {
			return p, nil
		}
	}

	return nil, errutil.ErrPlayerNotFound
}

func (d *Desk) objectWithId(uid int64) (*MapObject, error) {
	for _, o := range d.objects {
		if o.uid == uid {
			return o, nil
		}
	}

	return nil, errutil.ErrPlayerNotFound
}

func (d *Desk) robotWithId(uid int64) (*RobotInfo, error) {
	for _, r := range d.robots {
		if r.uid == uid {
			return r, nil
		}
	}

	return nil, errutil.ErrPlayerNotFound
}

func (d *Desk) isDestroy() bool {
	return d.status() == constant.DeskStatusDestory
}

// 摧毁桌子
func (d *Desk) destroy() {
	if d.status() == constant.DeskStatusDestory {
		d.logger.Info("桌子已经解散")
		return
	}

	close(d.die)

	// 标记销毁
	d.setStatus(constant.DeskStatusDestory)

	d.logger.Infof("销毁房间,deskid=%d", d.deskID)
	for i := range d.players {
		p := d.players[i]
		d.gameLogic.OnExit(p)

		d.players[i] = nil

		p.reset()
	}

	// 释放资源
	d.group.Close()

	//事件关闭
	d.eventMgr.Close()

	// 删除桌子
	scheduler.PushTask(func() {
		if room := roomManager.getRoom(d.cfgID); room != nil {
			room.setDesk(d.deskID, nil)
		}
	})
}

// 重置
func (d *Desk) reset() {
	d.matchUid = 0

	d.setStatus(constant.DeskStatusIdle)

	d.players = []*Player{}
	d.monsters = map[int64]*MonsterInfo{}
	d.objects = []*MapObject{}
	d.robots = []*RobotInfo{}

	d.chMove = make(chan *MoveRequest, chLen)
	d.chAttack = make(chan *AttackRequest, chLen)
	d.chMining = make(chan *MiningRequest, chLen)

	d.chDelayAttack = make(chan *DelayAction, chLen)
	d.chMonsterRelive = make(chan *MonRelive, chLen)
	d.chPlayerRelive = make(chan *PlayerRelive, chLen)

	d.chJoin = make(chan int64, chLen)
	d.chLeave = make(chan *Leave, chLen)
	d.chSkipBattle = make(chan int64, chLen)
	d.GuildRewardMember = make(map[int64]int64, 0)
	d.refreshList = []*MonsterRefresh{}
	d.monsterDeadTimes = map[int32]int32{}
	d.ConsecutiveKills = map[int64]*Kills{}
	d.mapInfo.clear()
	//事件重置
	d.eventMgr.Reset()
	d.FateCoin.reset()
	d.isInitMon = false
	d.group.LeaveAll()
}

// 加入房间
func (d *Desk) playerJoin(p *Player, reJoin bool, BattleId int32) error {
	if p.session == nil {
		return errutil.ErrNotFound
	}

	if p.desk != nil {
		p.logger.Debugf("加入房间失败,玩家已经在房间中")
		return errutil.ErrAlreadyInRoom
	}

	if roomCfg := global.GetRoomCfg(d.cfgID); roomCfg != nil {

		switch roomCfg.Roomtype {
		case constant.SoloTowerMap:
			// 爬塔房间设置层数
			d.Layers = p.SoloTowerInfo.Layers + 1
			d.SoloTowerCfg = global.GetEndlesstrialsCfg(p.SoloTowerInfo.Layers + 1)
			d.beginTime = time.Now()
		case constant.PVEPeacefulMap, constant.PVEWarMap:
			if !reJoin {
				p.tasks.dispatch(pb.TargetType_PVE_INSTANCE_ENTRY_REQUIREMENT, 0, 1)

				d.beginTime = time.Now()
			}
		case constant.WorldBossMap:
			d.beginTime = time.Now()
			p.dailyWorldBossTimes++
			d.DWorldBoss.curRoundDamage = 0
			d.serverID = p.serverID
		case constant.ExperienceFBMap, constant.GoldFBMap, constant.EnhanceStoneFBMap, constant.BookPageFBMap, constant.StarUpStoneFBMap, constant.SoulStoneFBMap:
			// isfirstchallge := p.dailyMapFirstChallenge(roomCfg.Roomtype)

			p.subDailyMapTime(roomCfg.Roomtype)
			d.beginTime = time.Now()
			d.scoring = global.XMLKeyValue{}
			d.reward = global.IdNum{}
			d.DDailyMap.buffRefreshTime = time.Now().Add(time.Duration(global.GlobalCfg.DailyMaps.ExpMap.BuffGap.RefreshTime) * time.Second)

			if roomCfg.Roomtype != constant.StarUpStoneFBMap && roomCfg.Roomtype != constant.SoulStoneFBMap {
				// 成长机制 攻防血全属性＝配置基础攻防血属性×[ INT（玩家等级÷30）＋1]
				d.monBuffStatus = &Stats{Attrs: make(map[pb.Attribute]int32)}
				accommodation := global.GetDailyMapAccommodation(roomCfg.Roomtype)
				if accommodation > 0 {
					v := (p.level/int32(accommodation)*int32(BaseRate) + 1)
					d.monBuffStatus.setStats(pb.Attribute_HEALTH_PERCENTAGE, v)
					d.monBuffStatus.setStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE, v)
					d.monBuffStatus.setStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE, v)
					d.monBuffStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE, v)
					d.monBuffStatus.setStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE, v)
				}
			}

		case constant.BattleMap, constant.BossMap, constant.TrialTowerMap, constant.RelicMap:
			d.BattleId = BattleId
			d.BattleCfg = global.GetBattleMapCfg(BattleId)
			d.beginTime = time.Now()

		case constant.GuildMap:
			if d.GuildId == 0 {
				d.beginTime = time.Now()
				startTime, _, err := global.GetGuildCopyTime()
				if err == nil {
					d.beginTime = startTime
				}

				d.GuildId = p.guildId
			}
		case constant.FightEndMap:
			if p.ObjectType() == ObjectTypePlayer {
				d.BattleId = BattleId
				d.BattleCfg = global.GetBattleMapCfg(BattleId)
				d.beginTime = time.Now()

			}
		}

	}

	return nil
}

// Loop 内加入房间
func (d *Desk) join(p *Player) error {

	var cell *Cell
	var dir pb.Direction
	switch d.roomCfg.Roomtype {
	case constant.FightEndMap, constant.PVPOneOnOneMap:
		if len(d.mapCfg.BornPoints) < 2 {
			logger.Warningf("进入地图失败, 未找到出生点 uid=%d  name=%s mapId=%d", p.uid, p.name, d.mapId)
			return errutil.ErrNotFound
		}
		var bornLocation *global.PointI
		if !p.isRobot() {
			bornLocation = &d.mapCfg.BornPoints[0]
			dir = pb.Direction_UP
		} else {
			bornLocation = &d.mapCfg.BornPoints[1]
			dir = pb.Direction_DOWN
		}
		if bornLocation == nil {
			logger.Warningf("进入地图失败, 未找到出生点 uid=%d  name=%s mapId=%d  出生点 %v", p.uid, p.name, d.mapId, d.mapCfg.BornPoints)
			return errutil.ErrNotFound
		}
		cell = d.mapInfo.getAvailableCell(bornLocation, 2)
		if cell == nil {
			logger.Warningf("进入地图失败, 未找到出生点 uid=%d  name=%s mapId=%d 出生点 %v", p.uid, p.name, d.mapId, d.mapCfg.BornPoints)
			return errutil.ErrNotFound
		}

		d.beginTime = time.Now()

	default:
		bornLocation := &d.mapCfg.BornPoints[global.RandIn(len(d.mapCfg.BornPoints))]
		cell = d.mapInfo.getAvailableCell(bornLocation, 2)
		if cell == nil {
			logger.Warningf("进入地图失败, 未找到出生点 uid=%d  name=%s mapId=%d", p.uid, p.name, d.mapId)
			return errutil.ErrNotFound
		}
		dir = pb.Direction_UP
	}

	p.dead = false
	p.bornLocation = cell.location
	p.direction = dir
	p.cell = d.mapInfo.addObject(p, &p.bornLocation)
	p.spawnLocation = p.cell.location
	p.attackTime = time.Now()
	p.moveTime = time.Now()
	p.setDesk(d)

	d.gameLogic.OnEnter(p)

	if d.mapInfo.safeArea(p.location) {
		p.logger.Debugf("玩家【%s】进入地图【%d】安全区。坐标:[%v]", p.name, d.mapCfg.ID, p.location)
	} else {
		p.logger.Warningf("玩家【%s】进入地图【%d】非安全区。坐标:[%v]", p.name, d.mapCfg.ID, p.location)
	}
	// 机器人不加入
	if !p.isRobot() {
		d.group.Add(p.session)
		p.defaultBuff()
		d.players = append(d.players, p)
	}

	// 全部伙伴进入房间
	for _, v := range p.partnerSlots {
		if v != nil && v.PartnerUid > 0 {

			partner := p.getPartnerInfo(v.PartnerUid)
			if partner == nil {
				continue
			}
			partner.dead = false
			partner.target = nil

			if bornLocation := partner.standLocation(int(v.SlotId)); bornLocation != nil {
				if d.addObject(partner.MapObject, bornLocation.location) {
					partner.direction = p.direction
				}
			} else {
				if d.addObject(partner.MapObject, p.bornLocation) {
					partner.direction = p.direction
				}
			}
			logger.Debugf("伙伴进入房间  uid %d  name %s location %v", partner.uid, partner.name, partner.location)
		}
	}
	return nil
}

func (d *Desk) robotJoin(r *RobotInfo, BattleId int32) error {

	r.defaultBuff()
	d.join(r.Player)
	d.robots = append(d.robots, r)
	d.logger.Debugf("机器人【%d : %s】加入房间【%d】", r.uid, r.name, d.deskID)
	return nil
}

func (d *Desk) playerRejoin(p *Player) error {
	if p, err := d.playerWithId(p.uid); err != nil && p != nil {
		p.logger.Warnf("[%s] 重连房间失败， 房间已解散", p.name)
		return nil
	}
	// 机器人不加入
	if !p.isRobot() {
		d.group.Add(p.session)
	}
	p.attackTime = time.Now()
	p.moveTime = time.Now()
	leftTime := p.expirationTime.Unix() - time.Now().Unix()
	packet := &pb.ReJoinDeskPush{
		TableInfo: &pb.TableInfo{
			RoomCfgId:   d.cfgID,
			RoomUid:     d.deskID,
			MonsterList: []*pb.MonsterInfo{},
			BattleId:    d.BattleId,
		},
		Location:  p.location,
		Direction: p.direction,
		LeftTime:  leftTime,
	}

	// 发送怪物列表
	for _, m := range d.monsters {
		packet.TableInfo.MonsterList = append(packet.TableInfo.MonsterList, &pb.MonsterInfo{
			CfgId:     m.cfgId,
			Uid:       m.uid,
			Name:      m.name,
			MaxHp:     m.GetMaxHP(),
			CurHp:     m.displayHP,
			Location:  m.location,
			Direction: m.direction,
		})
	}

	p.pushMsg(protocol.RouteReJoinDeskPush, packet)
	p.syncBuff()

	p.setPlayerStatus(PlayerStatusPlaying)
	d.syncDeskStatus()
	return nil
}

func (d *Desk) playerExit(p *Player, disconnect bool) {
	if p == nil {
		return
	}

	uid := p.uid
	if p.session != nil {
		d.group.Leave(p.session)
	} else {
		d.logger.Warningf("玩家【%d:%s】退出房间，session为空", uid, p.name)
	}

	var exitPlayer *Player
	restPlayers := []*Player{}
	for _, dp := range d.players {
		if dp.Uid() != uid {
			restPlayers = append(restPlayers, dp)
		} else {
			exitPlayer = dp
		}
	}

	if exitPlayer == nil {
		d.logger.Debugf("退出房间异常 玩家uid=%d   disconnect=%v", uid, disconnect)
		return
	}

	d.logger.Debugf("退出房间 玩家uid=%d name %s  matchid=%d  disconnect=%v", uid, p.name, d.matchUid, disconnect)

	switch d.roomCfg.Roomtype {
	case constant.ExperienceFBMap, constant.GoldFBMap, constant.EnhanceStoneFBMap, constant.BookPageFBMap:
		exitPlayer.desk.dailymapGameOver(constant.GameOverTypeTime)

	case constant.PVPOneOnOneMap:
		matchManager.leave(d.matchId, exitPlayer.uid)

	case constant.PVEPeacefulMap:
		delAfkTime := time.Now().Unix() - d.beginTime.Unix()
		exitPlayer.delAfkLeftTime(delAfkTime)
		// 守卫海岛上报
		d.upGuardingEndShushuDate(exitPlayer, constant.PVEPeacefulMap, int64(exitPlayer.afkLeftTime()), delAfkTime)
		exitPlayer.addChallengeRecord(d.roomCfg.Roomtype)
		if delAfkTime > 5*60 {
			// 挑战时间大于5分钟才算完成任务
			//p.tasks.dispatch(pb.TargetType_GUARD_ISLAND_INSTANCE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_GUARD_ISLAND_INSTANCE, Id: 0, Num: 1})
		}

	case constant.StarUpStoneFBMap, constant.SoulStoneFBMap:
		CompletionTime := time.Now().Unix() - d.beginTime.Unix()
		d.upBattleEndShushuDate(exitPlayer, nil, d.BattleId, nil, 3, d.round, false, d.roomCfg.Roomtype, CompletionTime, d.DDailyMap.mingCount, 0, 0)
		exitPlayer.addChallengeRecord(d.roomCfg.Roomtype)
		if d.roomCfg.Roomtype == constant.StarUpStoneFBMap {
			//p.tasks.dispatch(pb.TargetType_STAR_UP_STONE_INSTANCE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_STAR_UP_STONE_INSTANCE, Id: 0, Num: 1})
		}
		if d.roomCfg.Roomtype == constant.SoulStoneFBMap {
			//p.tasks.dispatch(pb.TargetType_SOUL_STONE_INSTANCE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_SOUL_STONE_INSTANCE, Id: 0, Num: 1})
		}
		p.addUpRecord(constant.InstanceNum)
	case constant.WorldBossMap:
		exitPlayer.rankChangeTCP(constant.RankTypeWorldBoss, d.DWorldBoss.curRoundDamage)

	}

	d.players = restPlayers

	d.gameLogic.OnExit(exitPlayer)

	exitPlayer.reset()

	if disconnect {
		exitPlayer.removeSession()
	}

	res := &pb.ExitPush{
		Uid: uid,
	}
	d.Broadcast(protocol.RouteExitPush, res)
}

// 踢出房间 code 1副本时间到，2死亡复活超时
func (d *Desk) playerKick(p *Player, code int32) {
	restPlayers := []*Player{}
	for _, dp := range d.players {
		if dp.Uid() != p.Uid() {
			restPlayers = append(restPlayers, dp)
		}
	}
	d.players = restPlayers

	// 扣除每日挂机时间
	if d.roomCfg.Roomtype == constant.PVEPeacefulMap {
		delAfkTime := time.Now().Unix() - d.beginTime.Unix()
		// 守卫海岛上报
		// p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_GUARD_ISLAND_INSTANCE, Id: 0, Num: 1})
		p.tasks.dispatch(pb.TargetType_GUARD_ISLAND_INSTANCE, 0, 1)
		p.delAfkLeftTime(delAfkTime)
		d.upGuardingEndShushuDate(p, constant.PVEPeacefulMap, int64(p.afkLeftTime()), delAfkTime)
		p.addChallengeRecord(d.roomCfg.Roomtype)
	}

	packet := &pb.ExitPush{
		Uid:  p.Uid(),
		Code: code,
	}
	d.Broadcast(protocol.RouteExitPush, packet)

	d.gameLogic.OnExit(p)

	if p.session != nil {
		d.group.Leave(p.session)
		p.reset()
	}
}

func (d *Desk) robotKick(r *RobotInfo, code int32) {
	restRobots := []*RobotInfo{}
	for _, dr := range d.robots {
		if dr.Uid() != r.Uid() {
			restRobots = append(restRobots, dr)
		}
	}
	d.robots = restRobots

	packet := &pb.ExitPush{
		Uid:  r.Uid(),
		Code: code,
	}
	d.Broadcast(protocol.RouteExitPush, packet)

	d.gameLogic.OnExit(r.Player)
}

func (d *Desk) upGuardingEndShushuDate(p *Player, roomType int32, leftTime int64, duration int64) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "guarding_island_end"
	// 主角信息
	if p.level > 110 {
		shushuDate["main_id"] = "高级"
		shushuDate["main_name"] = "高级"
	} else {
		shushuDate["main_id"] = "初级"
		shushuDate["main_name"] = "初级"
	}
	// 查询是否是首次挑战
	getChallengeRecord := p.getChallengeRecord(roomType)
	if getChallengeRecord == 0 {
		shushuDate["is_first_try"] = true
	} else {
		shushuDate["is_first_try"] = false
	}
	// 第一把挑战时间与挂机时间
	if p.afkDailyTime == duration {
		shushuDate["first_time_daily"] = true
	} else {
		shushuDate["first_time_daily"] = false
	}
	shushuDate["remaining_time"] = leftTime
	shushuDate["duration"] = duration
	if leftTime <= 0 {
		shushuDate["is_finished"] = true
	} else {
		shushuDate["is_finished"] = false
	}
	shushumain.ReportData(shushuDate)
}

// 刷怪
// func (d *Desk) checkStart() error {
// 	if d.status() == constant.DeskStatusCreate || d.status() == constant.DeskStatusIdle {
// 		if err := d.gameLogic.OnStart(); err == nil {
// 			//d.setStatus(constant.DeskStatusReady)
// 			return nil
// 		} else {
// 			d.logger.Errorf("刷怪失败 房间类型 %v err =%v", d.roomCfg.Roomtype, err)
// 			return err
// 		}
// 	} else {
// 		d.logger.Errorf("刷怪失败,房间状态错误, 房间类型 %v", d.roomCfg.Roomtype)
// 	}
// 	return errutil.ErrServerInternal
// }

func (d *Desk) initMonster() error {
	if d.status() == constant.DeskStatusReady {
		if err := d.gameLogic.OnStart(); err == nil {
			d.isInitMon = true
			return nil
		} else {
			d.logger.Errorf("刷怪失败 房间类型 %v err =%v", d.roomCfg.Roomtype, err)
			return err
		}
	} else {
		d.logger.Errorf("刷怪失败,房间状态错误, 房间类型 %v", d.roomCfg.Roomtype)
	}
	return errutil.ErrServerInternal
}

// 开始游戏
func (d *Desk) start() {
	if d.status() != constant.DeskStatusPlaying && d.status() != constant.DeskStatusReady {
		go d.loop()
	}
}

// 核心逻辑
func (d *Desk) loop() {
	defer func() {
		if err := recover(); err != nil {
			d.logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	d.setStatus(constant.DeskStatusReady)
	if d.isInitMon {
		d.setStatus(constant.DeskStatusPlaying)
	}
	d.logger.Debugln("******开始游戏循环******")

	for !d.isGameOver() {

		now := time.Now()

		select {
		case uid, ok := <-d.chJoin:
			if !ok {
				continue
			}
			if !d.isInitMon {
				if err := d.initMonster(); err != nil {
					d.gameOver(constant.GameOverType)
					return
				}
			}

			if np, ok := defaultManager.player(uid); ok {
				d.join(np)
			} else {
				d.logger.Warningf("玩家加入房间失败 房间CfgID %d  玩家Uid %d ", d.roomCfg.ID, uid)
				continue
			}

			if p, err := d.playerWithId(uid); err == nil {

				packet := pb.MonsterInPush{MonsterList: []*pb.MonsterInfo{}}
				for _, m := range d.monsters {
					if m.dead {
						continue
					}

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
				}

				p.pushMsg(protocol.RouteMonsterInPush, &packet)
				if d.roomCfg.Roomtype >= constant.ExperienceFBMap && d.roomCfg.Roomtype <= constant.SoulStoneFBMap {
					// 同步评级
					packet2 := &pb.DailyMapScorePush{
						MType:      pb.MapType(d.roomCfg.Roomtype),
						MonsterNum: int32(len(d.monsters)),
						Round:      d.round}
					p.pushMsg(protocol.RouteDailyMapScorePush, packet2)
				}
				d.syncDeskStatus()
				d.refreshMonster()

				if count, num := len(d.players), d.totalPlayerCount(); count < 1 {
					d.logger.Debugf("当前房间玩家数量不足，不能开始游戏，当前玩家=%d, 房间容量=%d", count, num)
					continue
				}
				d.setStatus(constant.DeskStatusPlaying)
			}
		default:
		}

		select {
		case leave, ok := <-d.chLeave:
			if !ok {
				continue
			}

			if p, err := d.playerWithId(leave.Uid); err == nil {
				d.playerExit(p, leave.Disconnect)
				continue
			}
		default:
		}

		select {
		case attack, ok := <-d.chAttack:
			if !ok {
				continue
			}
			var targets []IMapObject
			if p, err := d.playerWithId(attack.Uid); err == nil {
				for _, tar := range attack.Req.TargetUids {
					obj := d.objectByUid(tar)
					if obj == nil {
						break
					}
					targets = append(targets, obj)
				}
				p.attack(targets, attack.Req.SkillId, attack.Req.Location, attack.Req.Direction)
			}
		default:

		}

		select {
		case move, ok := <-d.chMove:
			if !ok {
				continue
			}
			if p, err := d.playerWithId(move.Uid); err == nil {
				p.move(move.Location)
			}

		default:
		}

		select {
		case da, ok := <-d.chDelayAttack:
			if !ok {
				break
			}

			if attacker := d.objectByUid(da.ObjUid); attacker != nil {
				attacker.DelayAction(da.Target, da.SkillId, da.Direction, da.Location)
			}

		default:
		}

		select {
		case mid, ok := <-d.chMonsterRelive:
			if !ok {
				break
			}
			Relive(d, mid)

		default:
		}

		select {
		case pr, ok := <-d.chPlayerRelive:
			if !ok {
				break
			}
			if p, err := d.playerWithId(pr.uid); err == nil {
				p.revive(pr.reliveType)
			}

		default:
		}

		select {
		case mining, ok := <-d.chMining:
			if !ok {
				continue
			}
			if p, err := d.playerWithId(mining.Uid); err == nil {
				p.miningStart(mining.TargetUid)
			}

		default:

		}

		select {
		case uid, ok := <-d.chSkipBattle:
			if !ok {
				continue
			}
			d.skipBattle(uid)

		default:

		}

		// 处理事件队列
		d.eventMgr.ProcessEvents()

		// 玩家
		for _, p := range d.players {
			if p.curState != PlayerStatusPlaying || p.session == nil {
				d.playerKick(p, constant.GameOverType)
				d.logger.Warningf("玩家状态异常,踢出房间,玩家id=%d, 玩家状态=%d", p.uid, p.curState)
				break
			}
			// 超踢出房间
			if d.roomCfg.Roomtype == constant.PVEPeacefulMap || d.roomCfg.Roomtype == constant.PVEWarMap {
				if now.After(p.expirationTime) || now.Unix() > p.expirationTime.Unix() {
					d.playerKick(p, constant.GameOverTypeTime)
					break
				}
			}

			if p.dead {
				if d.roomCfg.Roomtype == constant.PVPMultiOnMultiMap {
					if now.After(p.reviveTime) {
						if !p.revive(0) {
							p.reviveTime = now.Add(time.Second)
						}
					}
				} else if d.roomCfg.Roomtype == constant.PVEPeacefulMap || d.roomCfg.Roomtype == constant.PVEWarMap {
					if now.After(p.reviveTime) {
						if !p.revive(0) {
							p.reviveTime = now.Add(time.Second)
						}
					}
				} else if d.roomCfg.Roomtype == constant.WorldBossMap {
					if now.After(p.reviveTime) {
						d.worldBossGameOver(constant.GameOverTypeDeadFaild)
						return
					}
				} else if d.roomCfg.Roomtype == constant.EnhanceStoneFBMap || d.roomCfg.Roomtype == constant.BookPageFBMap ||
					d.roomCfg.Roomtype == constant.ExperienceFBMap || d.roomCfg.Roomtype == constant.GoldFBMap {

					if now.After(p.deadTime.Add(time.Duration(global.GlobalCfg.War.Autorevive.Time) * time.Second)) {
						d.dailymapGameOver(constant.GameOverTypePDead)
						return
					}
				} else if d.roomCfg.Roomtype == constant.SoloTowerMap {
					d.soloTowerGameOver(false)
					return
				} else if d.roomCfg.Roomtype == constant.BattleMap || d.roomCfg.Roomtype == constant.BossMap ||
					d.roomCfg.Roomtype == constant.TrialTowerMap || d.roomCfg.Roomtype == constant.RelicMap {
					if now.After(p.deadTime.Add(1500 * time.Millisecond)) {
						d.battleGameOver(constant.GameOverTypeDeadFaild)
						return
					}
				} else if d.roomCfg.Roomtype == constant.GuildMap {
					if now.After(p.reviveTime) {
						if !p.revive(0) {
							p.reviveTime = now.Add(time.Second)
						}
					}
				} else if d.roomCfg.Roomtype == constant.FightEndMap {
					if now.After(p.deadTime.Add(1500 * time.Millisecond)) {
						d.fightEndGameove(constant.GameOverTypeDeadFaild)
						return
					}
				}
			}

			if time.Now().After(p.moveTime.Add(time.Second*1)) && p.GetObAction() == ObjectActionMove {
				p.SetObAction(ObjectActionStand)
			}

			p.processBuff()
			p.processHp()

			if p.curState == PlayerStatusPlaying {
				for _, v := range p.partnerSlots {
					if v != nil && v.PartnerUid > 0 {
						if par := p.getPartnerInfo(v.PartnerUid); par != nil {
							par.processAI()
						}
					}
				}
				if p.GetObAction() == ObjectActionMining {
					p.checkMining()
				}
			}
		}

		for _, r := range d.robots {
			r.processAI()

			for _, v := range r.partnerSlots {
				if v != nil && v.PartnerUid > 0 {
					if par := r.getPartnerInfo(v.PartnerUid); par != nil {
						par.processAI()
					}
				}
			}

			if r.dead {
				if d.roomCfg.Roomtype == constant.PVEPeacefulMap || d.roomCfg.Roomtype == constant.PVEWarMap {
					if now.After(r.reviveTime) {
						if !r.revive(0) {
							r.reviveTime = now.Add(time.Second)
						}
					}
				} else if d.roomCfg.Roomtype == constant.FightEndMap {
					if now.After(r.deadTime.Add(1500 * time.Millisecond)) {
						d.fightEndGameove(constant.GameOverTypeDeadFaild)
						break
					}
				}
			}
		}
		// 怪物
		for _, m := range d.monsters {
			if m == nil {
				continue
			}
			m.processAI()
		}

		// 其他实体
		for _, v := range d.objects {
			if v == nil {
				continue
			}
			v.processAI()
		}

		switch d.roomCfg.Roomtype {
		case constant.SoloTowerMap:
			if now.After(d.beginTime.Add(time.Duration(d.SoloTowerCfg.TimeLimit) * time.Second)) {
				d.soloTowerGameOver(false)
				return
			}
			if len(d.monsters) == 0 && now.After(d.beginTime.Add(3*time.Second)) {
				d.soloTowerGameOver(true)
				return
			}
		case constant.PVPMultiOnMultiMap:
			if d.matchUid != 0 && d.MatchInfo != nil {
				if now.After(d.MatchInfo.startTime.Add(time.Duration(global.GlobalCfg.GVG.TimeSet.Duration) * time.Second)) {
					matchManager.matchResult <- d.matchUid
					break
				}
				if len(d.players) == 1 {
					matchManager.matchResult <- d.matchUid
				}
				if now.After(d.MatchInfo.buffRefreshTime) {
					buff_id := global.GlobalCfg.GVG.Buff.Keys[global.RandIn(len(global.GlobalCfg.GVG.Buff.Keys))]
					location := global.GlobalCfg.GVG.BuffGenerate.Keys[global.RandIn(len(global.GlobalCfg.GVG.BuffGenerate.Keys))]
					d.createObjectBuff(buff_id.ID, global.PointI{X: location.X, Y: location.Y}, int32(global.GlobalCfg.GVG.BuffGap.RefreshTime))
					d.DDailyMap.buffRefreshTime = time.Now().Add(time.Duration(global.GlobalCfg.GVG.BuffGap.RefreshTime) * time.Second)
				}
				if now.After(d.MatchInfo.boardRefreshTime) {
					d.MatchInfo.boardRefreshTime = time.Now().Add(time.Second)
					if d.MatchInfo.dirty {
						killpacket := &pb.PVPMeleeKillNumPush{List: []*pb.KillInfo{}}
						for _, mp := range d.MatchInfo.mplayers {
							killpacket.List = append(killpacket.List, &pb.KillInfo{Uid: mp.uid, Name: mp.name, Num: int32(mp.killNum)})
						}
						d.Broadcast(protocol.RoutePVPMeleeKillPush, killpacket)
						d.MatchInfo.dirty = false
						d.logger.Debugf("同步GVG击杀榜单")
					}
				}
			}
		case constant.PVPOneOnOneMap:
			if d.matchUid != 0 && d.MatchInfo != nil {
				if now.Before(d.beginTime.Add(3 * time.Second)) {
					break
				}
				if now.After(d.beginTime.Add(time.Duration(global.GlobalCfg.PVP.TimeSet.Duration) * time.Second)) {
					matchManager.matchResult <- d.matchUid
					break
				}
				var aliveCount int
				for _, p := range d.players {
					if !p.dead {
						aliveCount++
					}
					if p.dead && now.Before(p.deadTime.Add(1000*time.Millisecond)) {
						aliveCount++
					}
				}
				for _, r := range d.robots {
					if !r.dead {
						aliveCount++
					}
					if r.dead && now.Before(r.deadTime.Add(1000*time.Millisecond)) {
						aliveCount++
					}
				}
				if aliveCount == 1 {
					matchManager.matchResult <- d.matchUid
				}
				if len(d.players)+len(d.robots) <= 1 {
					matchManager.matchResult <- d.matchUid
				}
			}
		case constant.WorldBossMap:
			if now.After(d.beginTime.Add(time.Duration(global.GlobalCfg.WorldBoss.Challengetime.Duration) * time.Second)) {
				d.worldBossGameOver(constant.GameOverTypeTime)
				return
			}

			// select {
			// case updateHealth, ok := <-d.hpChangeChannel:
			// 	if ok {
			// 		for _, m := range d.monsters {
			// 			change := m.currentHP - updateHealth.HP
			// 			logger.Infof("loop 世界boss 当前血量变化 uid=%d  name=%s  m.currentHP=%d  updateHealth.HP=%d  updateHealth.HP=%d", m.uid, m.name, m.currentHP, updateHealth.HP, change)
			// 			if change > 0 {
			// 				logger.Infof("loop 世界boss 当前血量变化 uid=%d  name=%s  change=%d", m.uid, m.name, change)
			// 				m.changeHp(-change, updateHealth.Attacker)
			// 			}
			// 		}
			// 	}
			// default:
			// }

			// 成长机制 当前攻击值＝配置基础攻击值×[ INT（当前挑战时间÷10）＋1]
			if now.After(d.nextBuffTime) {
				ct := now.Unix() - d.beginTime.Unix()
				buffStat := &Stats{Attrs: make(map[pb.Attribute]int32)}
				buffStat.setStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE, int32((ct/int64(global.GlobalCfg.WorldBoss.Accommodation.Num)*int64(BaseRate) + 1)))
				buffStat.setStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE, int32((ct/int64(global.GlobalCfg.WorldBoss.Accommodation.Num)*int64(BaseRate) + 1)))
				for _, v := range d.monsters {
					v.AddBuff(7, v, int32(global.GlobalCfg.WorldBoss.Accommodation.Num*1000), buffStat, 0, 0)

					d.logger.Debugf(" 世界boss 增加buff , 物攻 %d   魔攻 %d", v.stat(pb.Attribute_PHYSICAL_ATTACK), v.stat(pb.Attribute_MAGIC_ATTACK))
				}

				d.nextBuffTime = now.Add(time.Duration(global.GlobalCfg.WorldBoss.Accommodation.Num) * time.Second)
			}

		case constant.ExperienceFBMap, constant.GoldFBMap:
			if now.After(d.beginTime.Add(time.Duration(global.GetDailyMapTime(d.roomCfg.Roomtype).Duration) * time.Second)) {
				d.dailymapGameOver(constant.GameOverTypeTime)
				return
			}
			if len(d.monsters) == 0 && now.After(d.beginTime.Add(3*time.Second)) {
				d.dailymapGameOver(constant.GameOverTypeMDead)
				return
			}

			var boss *MonsterInfo
			for _, v := range d.monsters {
				boss = v
				break
			}

			oldStage := d.scoring
			if boss != nil {
				d.totalDamage = boss.stat(pb.Attribute_MAX_HEALTH) - boss.currentHP
				bossHpPercent := float64(boss.stat(pb.Attribute_MAX_HEALTH)-boss.currentHP) / float64(boss.stat(pb.Attribute_MAX_HEALTH)) * BaseRate

				score := global.GetDailyMapScore(d.roomCfg.Roomtype).Stages
				for i := 0; i < len(score); i++ {
					if bossHpPercent >= float64(score[i].Value) {
						d.scoring = score[i]
					}
				}
			}

			// 同步评级
			if oldStage.Key != d.scoring.Key {
				packet := &pb.DailyMapScorePush{
					MType:       pb.MapType(d.roomCfg.Roomtype),
					TotalDamage: int64(d.totalDamage),
					Scoring:     d.scoring.Key,
					MonsterNum:  int32(len(d.monsters)),
					Round:       0,
				}

				d.Broadcast(protocol.RouteDailyMapScorePush, packet)
			}

			// 刷buff
			if now.After(d.DDailyMap.buffRefreshTime) {
				if d.roomCfg.Roomtype == constant.ExperienceFBMap {
					buff_id := global.GlobalCfg.DailyMaps.ExpMap.Buff.Keys[global.RandIn(len(global.GlobalCfg.DailyMaps.ExpMap.Buff.Keys))]
					location := global.GlobalCfg.DailyMaps.ExpMap.BuffGenerate.Keys[global.RandIn(len(global.GlobalCfg.DailyMaps.ExpMap.BuffGenerate.Keys))]
					d.createObjectBuff(buff_id.ID, global.PointI{X: location.X, Y: location.Y}, int32(global.GlobalCfg.DailyMaps.ExpMap.BuffGap.RefreshTime))
					d.DDailyMap.buffRefreshTime = time.Now().Add(time.Duration(global.GlobalCfg.DailyMaps.ExpMap.BuffGap.RefreshTime) * time.Second)
				} else if d.roomCfg.Roomtype == constant.GoldFBMap {
					buff_id := global.GlobalCfg.DailyMaps.GoldMap.Buff.Keys[global.RandIn(len(global.GlobalCfg.DailyMaps.GoldMap.Buff.Keys))]
					location := global.GlobalCfg.DailyMaps.GoldMap.BuffGenerate.Keys[global.RandIn(len(global.GlobalCfg.DailyMaps.GoldMap.BuffGenerate.Keys))]
					d.createObjectBuff(buff_id.ID, global.PointI{X: location.X, Y: location.Y}, int32(global.GlobalCfg.DailyMaps.GoldMap.BuffGap.RefreshTime))
					d.DDailyMap.buffRefreshTime = time.Now().Add(time.Duration(global.GlobalCfg.DailyMaps.GoldMap.BuffGap.RefreshTime) * time.Second)
				}
			}

		case constant.EnhanceStoneFBMap, constant.BookPageFBMap:
			if now.After(d.beginTime.Add(time.Duration(global.GetDailyMapTime(d.roomCfg.Roomtype).Duration) * time.Second)) {
				d.dailymapGameOver(constant.GameOverTypeTime)
				return
			}
			var totlecycle int32
			if d.roomCfg.Roomtype == constant.EnhanceStoneFBMap {
				totlecycle = int32(global.GlobalCfg.DailyMaps.StrengthenMap.TotalCycle.Num)
			} else if d.roomCfg.Roomtype == constant.BookPageFBMap {
				totlecycle = int32(global.GlobalCfg.DailyMaps.SkillPageMap.TotalCycle.Num)
			}
			if len(d.monsters) == 0 && now.After(d.beginTime.Add(3*time.Second)) {
				if d.round >= totlecycle {
					d.dailymapGameOver(constant.GameOverTypeMDead)
					return
				} else {
					d.round++
					// 刷怪
					var mNum int32
					for _, v := range d.mapCfg.MonsterIds {
						mNum += v.Num
						for i := 0; i < int(v.Num); i++ {
							d.addMonsterRelive(v.ID, v.Point, 0, 0, 0, 0, "")
						}
					}
					// 同步评级
					packet := &pb.DailyMapScorePush{
						MType:      pb.MapType(d.roomCfg.Roomtype),
						MonsterNum: mNum,
						Round:      d.round,
					}
					d.Broadcast(protocol.RouteDailyMapScorePush, packet)
				}
			}

		case constant.StarUpStoneFBMap, constant.SoulStoneFBMap, constant.PVEPeacefulMap:
			if d.roomCfg.Roomtype == constant.StarUpStoneFBMap || d.roomCfg.Roomtype == constant.SoulStoneFBMap {
				if now.After(d.beginTime.Add(time.Duration(global.GetDailyMapTime(d.roomCfg.Roomtype).Duration) * time.Second)) {
					d.dailymapGameOver(constant.GameOverTypeTime)
					return
				}
			}

		case constant.PVEWarMap:
			if !global.InWarTime(d.roomCfg.ID) {
				d.gameOver(constant.GameOverTypeTime)
				return
			}
			for i := 0; i < len(d.refreshList); i++ {
				if now.After(d.refreshList[i].RefreshTime) {
					mr := &MonRelive{CfgID: d.refreshList[i].CfgId, Location: d.refreshList[i].Location}
					Relive(d, mr)
					d.refreshList = append(d.refreshList[:i], d.refreshList[i+1:]...)

					d.logger.Debugf("战场刷怪， 怪物ID %d  剩余队列: %d", mr.CfgID, len(d.refreshList))
					break
				}
			}

		case constant.BattleMap, constant.BossMap, constant.RelicMap, constant.TrialTowerMap:
			if constant.BattleMap == d.roomCfg.Roomtype {
				//命运金币
				d.FateCoin.trigger(now, d)
			}

			if now.After(d.beginTime.Add(time.Duration(d.BattleCfg.TimeLimit) * time.Second)) {
				d.battleGameOver(constant.GameOverTypeTime)
				return
			}
			if len(d.monsters) == 0 && now.After(d.beginTime.Add(3*time.Second)) && d.isInitMon {
				d.battleGameOver(constant.GameOverTypeMDead)
				return
			}
		case constant.GuildMap:
			if now.After(d.beginTime.Add(time.Duration(global.GlobalCfg.Guild.GuildCopy.Lasttime) * time.Second * 60)) {
				d.guildGameOver(constant.GameOverTypeTime)
				return
			}

			if len(d.monsters) == 0 && now.After(d.beginTime.Add(3*time.Second)) {
				if d.BattleId < constant.GuildMap*100000+global.GlobalCfg.Guild.GuildBoss.MaxNum {
					d.BattleId++
					if g := guildManager.GetGuild(d.GuildId); g != nil {
						g.SetGuidBattleId(d.BattleId)
					}
					// 刷怪
					if guildMapCfg := global.GetBattleMapCfg(d.BattleId); guildMapCfg != nil {
						for _, v := range guildMapCfg.Monster {
							for i := 0; i < int(v.Num); i++ {
								d.addMonsterRelive(v.ID, v.Point, 0, 0, 0, 0, "")
								d.logger.Debugf("行会副本刷怪 ,房间ID %d 怪物ID %d  ", d.deskID, v.ID)
							}
						}
					}
				}
			}

		case constant.FightEndMap:
			if now.After(d.beginTime.Add(time.Duration(d.BattleCfg.TimeLimit) * time.Second)) {
				d.battleGameOver(constant.GameOverTypeTime)
				return
			}
		}

		time.Sleep(time.Millisecond * LoopTick)

	}

	d.gameOver(constant.GameOverType)
	d.logger.Debugln("******停止游戏循环******")
}

func (d *Desk) isGameOver() bool {
	if d.status() != constant.DeskStatusPlaying && d.status() != constant.DeskStatusReady {
		return true
	}

	//真人全部退出，结束房间
	if d.realPlayerCount() == 0 && d.status() == constant.DeskStatusPlaying {
		if d.roomCfg.Roomtype != constant.PVEWarMap {
			d.setStatus(constant.DeskStatusIdle)
			return true
		}
	}
	return false
}

func (d *Desk) gameOver(code int32) {
	for _, p := range d.players {
		d.playerKick(p, code)
	}

	for _, r := range d.robots {
		d.robotKick(r, code)
		robotMgr.onGameOver(r.Uid())
	}

	for _, m := range d.monsters {
		d.delMonster(m)
	}

	for _, ob := range d.objects {
		d.delObject(ob.Uid())
	}

	d.gameLogic.OnGameOver()
}

func (d *Desk) soloTowerGameOver(win bool) {
	d.setStatus(constant.DeskStatusIdle)
	if len(d.players) == 1 {
		p := d.players[0]
		packet := &pb.GameResultEndlesstrialsPush{Win: win, Layer: d.SoloTowerCfg.Layers, Items: []*pb.ItemInfo{}}
		if win {
			if p.SoloTowerInfo.Layers == 0 {
				p.SoloTowerInfo.RewardTime = time.Now()
			} else {
				// 通关后自动领取累计离线奖励
				if err := p.soloTowerOfflineReward(); err != nil {
					p.logger.Errorf("通关后自动发放离线奖励失败 err=%s", err.Error())
				}
			}

			bSendMail := false
			for _, i := range d.SoloTowerCfg.Settlement {
				if !p.checkAddSource(i.ID, int64(i.Num)) {
					bSendMail = true
				}
				packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: i.ID, Num: int64(i.Num)})
			}
			if !bSendMail {
				// for _, i := range d.SoloTowerCfg.Settlement {
				// 	p.addSource(i.ID, int64(i.Num), constant.ItemAddSoloTower)
				// }
				sourceInfo := []global.IdNum{}
				for _, v := range d.SoloTowerCfg.Settlement {
					sourceInfo = append(sourceInfo, global.IdNum{ID: v.ID, Num: v.Num})
				}
				p.addSourceBatch(sourceInfo, constant.ItemAddSoloTower)
			} else {
				attachment := []protocol.Attachment{}
				for _, i := range d.SoloTowerCfg.Settlement {
					attachment = append(attachment, protocol.Attachment{UserPropsId: int64(i.ID), Num: int64(i.Num)})
				}
				p.SendMail(constant.ItemAddSoloTower, attachment)
			}

			// 跑马灯
			if p.SoloTowerInfo.Layers > 0 && p.SoloTowerInfo.Layers%10 == 0 {
				context := []string{p.name, strconv.Itoa(int(p.SoloTowerInfo.Layers))}
				defaultManager.Lamp(int32(constant.LampTrialTowerLevelUp), p.serverID, context)
			}
			p.SoloTowerInfo.Layers++
			//p.tasks.dispatch(pb.TargetType_TOTAL_ATTEMPTS_IN_ENDLESS_TRIALS, 0, p.desk.SoloTower.Layers)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_ATTEMPTS_IN_ENDLESS_TRIALS, Id: 0, Num: p.desk.SoloTower.Layers})
			p.rankChangeTCP(constant.RankTypeTrial, int64(p.SoloTowerInfo.Layers))
		}
		p.pushMsg(protocol.RouteGameResultEndlesstrialsPush, packet)
		d.logger.Debugf("爬塔完成  userid=%d , win =%t, layer=%d ", p.uid, win, d.SoloTowerCfg.Layers)
	}
	d.SoloTower.Layers = 0
	// d.SoloTower.SoloTowerCfg = nil
	d.gameOver(constant.GameOverType)
}

func (d *Desk) worldBossGameOver(code int32) {
	d.setStatus(constant.DeskStatusIdle)
	if len(d.players) == 1 {
		if p := d.players[0]; p != nil {
			d.gameOver(code)
			p.rankChangeTCP(constant.RankTypeWorldBoss, d.DWorldBoss.curRoundDamage)
			d.logger.Debugf("世界boss战斗结束  userid=%d , totalWorldBossDamage =%d ", p.uid, p.totalWorldBossDamage)
			// 战斗时长
			CompletionTime := time.Now().Unix() - d.beginTime.Unix()
			current_damage := d.curRoundDamage
			all_damage := p.totalWorldBossDamage

			d.upBattleEndShushuDate(p, nil, d.BattleId, nil, code, d.round, false, constant.WorldBossMap, CompletionTime, 0, current_damage, all_damage)
			p.addChallengeRecord(d.roomCfg.Roomtype)
		}
	}
}

func (d *Desk) dailymapGameOver(code int32) {
	d.setStatus(constant.DeskStatusIdle)
	if len(d.players) == 0 {
		return
	}
	CompletionTime := time.Now().Unix() - d.beginTime.Unix()
	if p := d.players[0]; p != nil {
		var baseValue global.XMLIdNum
		if d.roomCfg.Roomtype == constant.ExperienceFBMap {
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_EXPERIENCE_FB_CLEARED, Id: 0, Num: 1})
			//p.tasks.dispatch(pb.TargetType_TOTAL_EXPERIENCE_FB_CLEARED, 0, 1)
			baseValue = global.GlobalCfg.DailyMaps.ExpMap.BaseValue
		} else if d.roomCfg.Roomtype == constant.GoldFBMap {
			//p.tasks.dispatch(pb.TargetType_TOTAL_GOLD_FB_CLEARED, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_GOLD_FB_CLEARED, Id: 0, Num: 1})
			baseValue = global.GlobalCfg.DailyMaps.GoldMap.BaseValue
		} else if d.roomCfg.Roomtype == constant.EnhanceStoneFBMap {
			//p.tasks.dispatch(pb.TargetType_ENHANCE_STONE_INSTANCE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_ENHANCE_STONE_INSTANCE, Id: 0, Num: 1})
			baseValue = global.GlobalCfg.DailyMaps.StrengthenMap.BaseValue
		} else if d.roomCfg.Roomtype == constant.BookPageFBMap {
			//p.tasks.dispatch(pb.TargetType_BOOK_PAGE_INSTANCE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_BOOK_PAGE_INSTANCE, Id: 0, Num: 1})
			baseValue = global.GlobalCfg.DailyMaps.SkillPageMap.BaseValue
		} else {
			if d.roomCfg.Roomtype == constant.SoulStoneFBMap {
				//p.tasks.dispatch(pb.TargetType_SOUL_STONE_INSTANCE, 0, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_SOUL_STONE_INSTANCE, Id: 0, Num: 1})

			} else if d.roomCfg.Roomtype == constant.StarUpStoneFBMap {
				//p.tasks.dispatch(pb.TargetType_STAR_UP_STONE_INSTANCE, 0, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_STAR_UP_STONE_INSTANCE, Id: 0, Num: 1})
			}
			//p.tasks.dispatch(pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE, 0, 1)
			p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE, Id: 0, Num: 1})
			d.gameOver(code)
			p.addChallengeRecord(d.roomCfg.Roomtype)
			return
		}

		var num float64
		if d.roomCfg.Roomtype == constant.ExperienceFBMap || d.roomCfg.Roomtype == constant.GoldFBMap {
			num = float64(baseValue.Num) * (1 + float64(d.scoring.Value)/BaseRate) * float64(p.level)
		}
		if d.roomCfg.Roomtype == constant.EnhanceStoneFBMap || d.roomCfg.Roomtype == constant.BookPageFBMap {
			num = float64(baseValue.Num) * float64(d.round) * float64(p.level)
		}
		num = math.Max(num, 1)

		d.reward = global.IdNum{ID: baseValue.ID, Num: int32(num)}
		p.DailyMapReward.reward = append(p.DailyMapReward.reward, d.reward)
		p.DailyMapReward.roomType = d.roomCfg.Roomtype
		if p.checkAddSource(baseValue.ID, int64(num)) {
			p.addSource(baseValue.ID, int64(num), constant.ItemAddDailymap)
		} else {
			attachment := []protocol.Attachment{}
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(baseValue.ID), Num: int64(num)})
			p.SendMail(constant.ItemAddDailymap, attachment)
		}

		packet := &pb.GameResultBattlePush{Items: []*pb.ItemInfo{}, RoomId: d.roomCfg.ID, OverType: pb.GameOverType(code)}
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: baseValue.ID, Num: int64(num)})
		p.pushMsg(protocol.GameResultBattleMapPush, packet)
		//p.tasks.dispatch(pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE, 0, 1)
		p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE, Id: 0, Num: 1})

		packetItems := []global.IdNum{{ID: baseValue.ID, Num: int32(num)}}
		// d.logger.Debugf("每日必做副本结算 Type: %d  奖励 %+v  评分 %+v  总伤害 %d  波次 %d", d.roomCfg.Roomtype, d.reward, d.scoring, d.totalDamage, d.round)
		p.addUpRecord(constant.InstanceNum) // 增加统计
		d.upBattleEndShushuDate(p, nil, d.BattleId, packetItems, code, d.round, false, 0, CompletionTime, 0, 0, 0)
		p.addChallengeRecord(d.roomCfg.Roomtype)
	}
	d.gameOver(code)
}

func (d *Desk) battleGameOver(code int32) {
	d.setStatus(constant.DeskStatusIdle)
	if len(d.players) >= 1 {
		p := d.players[0]
		packet := &pb.GameResultBattlePush{BattleId: d.BattleId, OverType: pb.GameOverType(code), Items: []*pb.ItemInfo{}, RoomId: d.roomCfg.ID}

		// 战斗时长
		packet.CompletionTime = time.Now().Unix() - d.beginTime.Unix()
		packetItems := []global.IdNum{}
		is_first_try := false
		if code == constant.GameOverTypeMDead {
			// bSendMail := false
			// for _, i := range d.BattleCfg.Settlement {
			// 	if !p.checkAddSource(i.ID, int64(i.Num)) {
			// 		bSendMail = true
			// 	}
			// 	packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: i.ID, Num: int64(i.Num)})
			// }
			// if !bSendMail {
			// 	// for _, i := range d.BattleCfg.Settlement {
			// 	// 	packetItems = append(packetItems, global.IdNum{ID: i.ID, Num: i.Num})
			// 	// 	p.addSource(i.ID, int64(i.Num), constant.ItemAddSoloTower)
			// 	// }
			// 	sourceInfo := []global.IdNum{}
			// 	for _, v := range d.BattleCfg.Settlement {
			// 		sourceInfo = append(sourceInfo, global.IdNum{ID: v.ID, Num: v.Num})
			// 	}
			// 	p.addSourceBatch(sourceInfo, constant.ItemAddSoloTower)
			// } else {
			// 	attachment := []protocol.Attachment{}
			// 	for _, i := range d.BattleCfg.Settlement {
			// 		packetItems = append(packetItems, global.IdNum{ID: i.ID, Num: i.Num})
			// 		attachment = append(attachment, protocol.Attachment{UserPropsId: int64(i.ID), Num: int64(i.Num)})
			// 	}
			// 	p.SendMail(constant.ItemAddSoloTower, attachment)
			// }
			p.distributeRewards(p, d.BattleCfg.Settlement, constant.ItemAddSoloTower, false)
			for _, i := range d.BattleCfg.Settlement {
				packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: i.ID, Num: int64(i.Num)})
				packetItems = append(packetItems, global.IdNum{ID: i.ID, Num: i.Num})
			}

			switch d.BattleCfg.CopyType {
			case constant.BattleMap:
				is_first_try = true
				d.upMainShushuDate(p, d.BattleId)
				p.BattleMap.BattleId = d.BattleId
				//p.tasks.dispatch(pb.TargetType_BATTLE_CHAPTERS_CLEARED, d.BattleId, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_BATTLE_CHAPTERS_CLEARED, Id: d.BattleId, Num: 1})
				p.rankChangeTCP(constant.RankTypeCampaignStage, int64(p.BattleMap.BattleId))
				// 跑马灯 TODO 这个是对应 battle 表的副本类型 14 没有很好的方法能动态获取每个章节的副本数
				if p.BattleMap.BattleId > 0 && p.BattleMap.BattleId%40 == 0 {
					context := []string{p.name, strconv.Itoa(int(p.BattleMap.BattleId % 100000 / 1000))}
					defaultManager.Lamp(int32(constant.LampCampaignChapterCompleted), p.serverID, context)
				}
				// 战役副本的翻牌奖励
				picker := randx.New(d.BattleCfg.Randomaward, func(i global.IdNumRate) int {
					return int(i.Rate)
				})
				filpReward, err := picker.Pick()
				if err != nil {
					d.logger.Errorf("副本随机奖励失败 err=%s", err.Error())
				}
				rewards := make([]global.IdNum, 0)
				if filpReward.Id != 0 {
					rewards = append(rewards, global.IdNum{ID: filpReward.Id, Num: filpReward.Num})
					packet.FlipCardReward = &pb.ItemInfo{CfgId: filpReward.Id, Num: int64(filpReward.Num)}
					//p.distributeRewards(p, []global.IdNum{{ID: filpReward.Id, Num: filpReward.Num}}, constant.ModuleCampaign.String(), true)
				}
				//命运金币
				rewardCfg := global.GlobalCfg.BaseValue.FateCoin
				if d.FateCoin.isReward() && rewardCfg.AwardId != 0 {
					rewards = append(rewards, global.IdNum{ID: rewardCfg.AwardId, Num: rewardCfg.AwardNum})
					packet.FateCoinReward = &pb.ItemInfo{CfgId: rewardCfg.AwardId, Num: int64(rewardCfg.AwardNum)}
				}
				if len(rewards) != 0 {
					p.distributeRewards(p, rewards, constant.ModuleCampaign.String(), true)
				}

			case constant.BossMap:
				if p.BossMap.BattleId < d.BattleId {
					d.upBossShushuDate(p, d.BattleId)
					is_first_try = true
					p.BossMap.BattleId = d.BattleId
				}
				p.BossMap.ChallengeTime--
				p.addUpRecord(constant.InstanceNum) // 增加统计
				//p.tasks.dispatch(pb.TargetType_BOSS_LEVELS_CLEARED, d.BattleId, 1)
				//p.tasks.dispatch(pb.TargetType_TOTAL_BOSS_CHALLENGE, 0, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_BOSS_LEVELS_CLEARED, Id: d.BattleId, Num: 1})
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_BOSS_CHALLENGE, Id: 0, Num: 1})

				// 跑马灯
				if p.BossMap.BattleId > 0 && p.BossMap.BattleId%5 == 0 {
					context := []string{p.name, strconv.Itoa(int(p.BossMap.BattleId % 1000))}
					defaultManager.Lamp(int32(constant.LampBossChapterCompleted), p.serverID, context)
				}

			case constant.TrialTowerMap:
				is_first_try = true
				d.upTrialShushuDate(p, d.BattleId)
				p.TrialTowerMap.BattleId = d.BattleId
				//p.tasks.dispatch(pb.TargetType_TRIAL_TOWER_FLOORS_CLEARED, d.BattleId, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TRIAL_TOWER_FLOORS_CLEARED, Id: d.BattleId, Num: 1})

				// 跑马灯
				if p.TrialTowerMap.BattleId > 0 && p.TrialTowerMap.BattleId%10 == 0 {
					context := []string{p.name, strconv.Itoa(int(p.TrialTowerMap.BattleId % 1000))}
					defaultManager.Lamp(int32(constant.LampTrialTowerLevelUp), p.serverID, context)
				}
			case constant.RelicMap:
				//p.tasks.dispatch(pb.TargetType_RELIC_CHAPTERS_CLEARED, d.BattleId, 1)
				//p.tasks.dispatch(pb.TargetType_TOTAL_RELIC_CHALLENGE, 0, 1)
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_RELIC_CHAPTERS_CLEARED, Id: d.BattleId, Num: 1})
				p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_TOTAL_RELIC_CHALLENGE, Id: 0, Num: 1})

				if p.RelicMap.BattleId < d.BattleId {
					d.upRuinsShushuDate(p, d.BattleId)
					is_first_try = true

					p.RelicMap.BattleId = d.BattleId
				}

				// 成功通关都扣次数
				p.RelicMap.ChallengeTime--
				p.addUpRecord(constant.InstanceNum) // 增加统计

				// 副本评价
				chapterId := utils.GetChapterID(d.BattleId/1000) % 100 // 章节id
				chapterStars := 0                                      // 某章节的所有关卡获取的星数
				timestar := 1
				allalive := 1
				allhealth := 1
				totalHp := int32(0) // 剩余总血量
				curHp := int32(0)   // 最大总血量
				for _, v := range p.partnerSlots {
					if v.PartnerUid == 0 {
						continue
					}

					if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
						if partner.dead {
							allalive = 0
						}

						curHp += partner.currentHP
						totalHp += partner.GetMaxHP()
						partner.dead = false
					}
				}
				// 总血量不能小于50%
				curHp += p.currentHP
				totalHp += p.GetMaxHP()
				if curHp < totalHp/2 {
					allhealth = 0
				}

				allstar := timestar + allalive + allhealth

				bFind := false
				add := false
				for i := 0; i < len(p.RelicMap.Evaluation); i++ {
					if p.RelicMap.Evaluation[i].Key == d.BattleId {
						if int32(allstar) > p.RelicMap.Evaluation[i].Value {
							p.RelicMap.Evaluation[i].Value = int32(allstar)
							add = true
						}
						bFind = true
					}

					if utils.GetChapterID(p.RelicMap.Evaluation[i].Key) == chapterId {
						chapterStars += int(p.RelicMap.Evaluation[i].Value)
					}
				}
				if !bFind {
					add = true
					p.RelicMap.Evaluation = append(p.RelicMap.Evaluation, global.KeyValue{Key: d.BattleId, Value: int32(allstar)})
					chapterStars += allstar
				}
				if add {
					//p.tasks.dispatch(pb.TargetType_RELIC_TOTAL_STAR, 0, int32(allstar))
					p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_RELIC_TOTAL_STAR, Id: 0, Num: int32(allstar)})
				}

				// battleID: 1701005
				p.tasks.dispatch(pb.TargetType_RELIC_CHAPTERS_ALL_3STARS, chapterId, int32(chapterStars))

			}
		}

		p.pushMsg(protocol.GameResultBattleMapPush, packet)
		CompletionTime := time.Now().Unix() - d.beginTime.Unix()
		d.upBattleEndShushuDate(p, nil, d.BattleId, packetItems, code, d.round, is_first_try, d.BattleCfg.CopyType, CompletionTime, 0, 0, 0)
		p.addChallengeRecord(d.roomCfg.Roomtype)
	}

	d.gameOver(constant.GameOverType)
}

func (d *Desk) upMainShushuDate(p *Player, battleId int32) {
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_main_id"] = battleId
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (d *Desk) upBossShushuDate(p *Player, battleId int32) {
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_boss"] = battleId
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (d *Desk) upTrialShushuDate(p *Player, battleId int32) {
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_trial"] = battleId
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (d *Desk) upRuinsShushuDate(p *Player, battleId int32) {
	// 创建一interface{}类型的map
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_ruins"] = battleId
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (d *Desk) guildGameOver(code int32) {
	d.setStatus(constant.DeskStatusIdle)

	for k := range d.GuildRewardMember {
		if p, err := d.playerWithId(k); err == nil {

			attachment := []protocol.Attachment{}
			packetItems := []global.IdNum{}
			for _, item := range d.GuildReward.CurReward {
				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.K), Num: int64(item.V)})
				packetItems = append(packetItems, global.IdNum{ID: item.K, Num: item.V})
			}
			p.SendMail(constant.ItemAddGame, attachment, constant.EmailContentGuidMap)
			CompletionTime := time.Now().Unix() - d.beginTime.Unix()
			d.upBattleEndShushuDate(p, nil, d.BattleId, packetItems, code, d.round, false, constant.GuildMap, CompletionTime, 0, 0, 0)
			p.addChallengeRecord(d.roomCfg.Roomtype)
		}
	}
	d.GuildId = 0
	d.GuildReward = &pb.GuildCopyRewardPush{}

	d.logger.Debugf("公会副本结算 Type: %d  奖励 %+v  伤害 %+v  ", d.roomCfg.Roomtype, d.GuildReward, d.GuildRewardMember)
	d.gameOver(code)
}

func (d *Desk) fightEndGameove(code int32) {
	packet := &pb.GameResultBattlePush{BattleId: d.BattleId, OverType: pb.GameOverType(code), Items: []*pb.ItemInfo{}, RoomId: d.roomCfg.ID}
	// 战斗时长
	packet.CompletionTime = time.Now().Unix() - d.beginTime.Unix()

	var p *Player
	var p2 *Player
	if code != constant.GameOverTypeTime {
		for _, v := range d.players {
			if v.ObjectType() == ObjectTypeRobot {
				p2 = v
			}
			if v.ObjectType() == ObjectTypePlayer {
				p = v
				if !v.dead {
					code = constant.GameOverTypeMDead
					packet.OverType = pb.GameOverType(code)
				}
			}
		}
	}
	packetItems := make([]global.IdNum, 0)
	is_first_try := false
	if p != nil {

		if code == constant.GameOverTypeMDead {
			if p.FightEndMap.BattleId < d.BattleId {
				is_first_try = true
				p.FightEndMap.BattleId = d.BattleId
				bSendMail := false
				for _, i := range d.BattleCfg.Settlement {
					if !p.checkAddSource(i.ID, int64(i.Num)) {
						bSendMail = true
					}
					packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: i.ID, Num: int64(i.Num)})
				}
				if !bSendMail {
					// for _, i := range d.BattleCfg.Settlement {
					// 	packetItems = append(packetItems, global.IdNum{ID: i.ID, Num: i.Num})
					// 	p.addSource(i.ID, int64(i.Num), constant.ItemAddSoloTower)
					// }
					sourceInfo := []global.IdNum{}
					for _, v := range d.BattleCfg.Settlement {
						sourceInfo = append(sourceInfo, global.IdNum{ID: v.ID, Num: v.Num})
					}
					p.addSourceBatch(sourceInfo, constant.ItemAddSoloTower)
				} else {
					attachment := []protocol.Attachment{}
					for _, i := range d.BattleCfg.Settlement {
						packetItems = append(packetItems, global.IdNum{ID: i.ID, Num: i.Num})
						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(i.ID), Num: int64(i.Num)})
					}
					p.SendMail(constant.ItemAddSoloTower, attachment)
				}
			}

			for i := 0; i < int(len(p.FightEndMap.Evaluation)); i++ {
				if p.FightEndMap.Evaluation[i].BattleId == d.BattleId {
					if p.FightEndMap.Evaluation[i].ChallengeTime == 0 {
						p.FightEndMap.Evaluation[i].ChallengeTime = int32(packet.CompletionTime)
						// 推给排行榜
						time := int32(0)
						maxLevel := 0
						for _, v := range p.FightEndMap.Evaluation {
							time = time + v.ChallengeTime
							if v.ChallengeTime > 0 && int(v.BattleId)%100000 > maxLevel {
								maxLevel = int(v.BattleId) % 100000
							}
						}
						point := int64(maxLevel)*1000000 + (1000000 - int64(time))
						p.rankChangeTCP(constant.RankTypeGightEnd, int64(point))
						p.rankChangeTCP(constant.RankTypeGightEndlayers, int64(d.BattleId))
					} else if p.FightEndMap.Evaluation[i].ChallengeTime > int32(packet.CompletionTime) {
						p.FightEndMap.Evaluation[i].ChallengeTime = int32(packet.CompletionTime)
						// 推给排行榜
						// 推给排行榜
						time := int32(0)
						maxLevel := 0
						for _, v := range p.FightEndMap.Evaluation {
							time = time + v.ChallengeTime
							if v.ChallengeTime > 0 && int(v.BattleId)%100000 > maxLevel {
								maxLevel = int(v.BattleId) % 100000
							}
						}
						point := int64(maxLevel)*1000000 + (1000000 - int64(time))
						p.rankChangeTCP(constant.RankTypeGightEnd, int64(point))
						p.rankChangeTCP(constant.RankTypeGightEndlayers, int64(d.BattleId))
					}
					break
				}
			}
		}
		//p.tasks.dispatch(pb.TargetType_ENDLESS_TRIALS_STAGE, d.BattleId, 1)
		p.dispatchTask(&event.TaskEventData{Key: pb.TargetType_ENDLESS_TRIALS_STAGE, Id: d.BattleId, Num: 1})
		d.upBattleEndShushuDate(p, p2, d.BattleId, packetItems, code, d.round, is_first_try, constant.FightEndMap, packet.CompletionTime, 0, 0, 0)
		p.addChallengeRecord(d.roomCfg.Roomtype)
		p.pushMsg(protocol.GameResultBattleMapPush, packet)
	}

	d.gameOver(code)
}

func (d *Desk) upBattleEndShushuDate(p *Player, p2 *Player, BattleId int32, packetItems []global.IdNum,
	code int32, round int32, is_first_try bool, roomType int32, duration int64, mingCount int32,
	current_damage int64, all_damage int64) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 查询是否是首次挑战
	getChallengeRecord := p.getChallengeRecord(roomType)
	if getChallengeRecord == 0 {
		shushuDate["is_first_try"] = true
	} else {
		shushuDate["is_first_try"] = false
	}

	// 主角信息
	// 查看战斗信息
	main_type := roomType
	main_name := ""
	cfg := global.GetBattleMapCfg(BattleId)
	if cfg != nil {
		main_name = fmt.Sprintf("%d-%d", cfg.Chapter, cfg.Layers)
	}
	partnerInfo := make([]map[string]interface{}, 0)
	switch roomType {
	case constant.GuildMap:
		shushuDate["name"] = "guild_boss_end"
		if p.guildId > 0 {
			shushuDate["guild_id"] = p.guildId
			// 查询工会信息
			res, err := guildManager.validateGuild(p.guildId)
			if err == nil {
				// 查询boss名字
				monsterId := cfg.Monster[0].ID
				if monsterId != 0 {
					if cfg2 := global.GetMonsterCfg(monsterId); cfg2 != nil {
						shushuDate["boss_name"] = cfg2.Name
					}
				}
				// 查询工会信息
				shushuDate["guild_level"] = res.Guild.Level
				shushuDate["guild_name"] = res.Guild.Name
			}
		}
	case constant.WorldBossMap:
		shushuDate["name"] = "world_boss_end"
		shushuDate["current_damage"] = current_damage
		shushuDate["all_damage"] = all_damage
		bossMsg := worldBossMgr.getBoss(p.serverID)
		if bossMsg != nil {
			shushuDate["main_name"] = bossMsg.name
		}
	case constant.FightEndMap:
		shushuDate["name"] = "yizhandaodi_end"
		shushuDate["role_stage"] = main_name

		// shushuDate["opponent_name"] = robotInfo.name
		// shushuDate["opponent_id"] = robotInfo.uid
		// shushuDate["opponent_power"] = robotInfo.combat
		// shushuDate["opponent_rank"] = robotInfo.rank
		// shushuDate["is_robot"] = true
	default:
		shushuDate["name"] = "challenge_battle_end"
	}
	if p2 != nil {
		shushuDate["opponent_name"] = p2.name
		shushuDate["opponent_id"] = fmt.Sprintf("%d", p2.uid)
		shushuDate["opponent_power"] = p2.combat
		if p2.robotType == RobotTypePlayer {
			shushuDate["is_robot"] = true
		} else {
			shushuDate["is_robot"] = false
		}
		rankInfo1 := &protocol.RankOperateRequest{
			Uid:  p2.uid,
			Type: 3,
			All:  3,
		}
		res1, err := provide.RankList(rankInfo1)
		shushuDate["opponent_rank"] = 0
		if err == nil {
			if res1.RankOwn.Num > 0 {
				shushuDate["opponent_rank"] = res1.RankOwn.Num
			}
		}

	}
	shushuDate["main_name"] = main_name
	shushuDate["main_id"] = BattleId
	shushuDate["main_type"] = main_type
	shushuDate["duration"] = duration
	shushuDate["open_boxes"] = mingCount

	// 伙伴信息
	for _, v := range p.partnerSlots {
		if v != nil && v.PartnerUid > 0 {
			partner := p.getPartnerInfo(v.PartnerUid)
			if partner != nil {
				partnerInfo = append(partnerInfo, map[string]interface{}{
					"partner_id":      partner.cfgId,
					"partner_name":    partner.name,
					"partner_quality": partner.cfg.Quality,
					"partner_level":   partner.level,
					"partner_job":     int32(partner.job),
					"partner_star":    partner.star,
					"partner_power":   partner.combat,
				})
			}
		}
	}
	shushuDate["partner_info"] = partnerInfo

	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range packetItems {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(v.ID); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["get_item_info"] = get_item_info
	shushuDate["wave_num"] = round
	shushuDate["open_boxes"] = 0
	// 是否首次尝试
	// shushuDate["is_first_try"] = p.isFirstTry
	if pb.GameOverType(code) == pb.GameOverType_GameOverTypeMDead || pb.GameOverType(code) == pb.GameOverType_GameOverTypeCopySweep {
		shushuDate["is_success"] = true
	} else {
		shushuDate["is_success"] = false
	}
	// 技能信息
	skillInfo := make([]map[string]interface{}, 0)
	for _, v := range p.skillsSlot {
		if v != 0 && v > 0 {
			if scfg, err := global.GetSkillCfg(v); err == nil {
				skill := map[string]interface{}{
					"skill_id":    v,
					"skill_name":  scfg.Name,
					"skill_level": scfg.Level,
				}
				if getSkill := p.getSkill(v); getSkill != nil {
					skill["skill_grade"] = getSkill.AdvancedLevel
				}
				skillInfo = append(skillInfo, skill)
			}
		}
	}
	shushuDate["skill_info"] = skillInfo
	shushumain.ReportData(shushuDate)
}

// 删除怪物
func (d *Desk) delMonster(m *MonsterInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.monsters[m.uid]

	if ok {
		d.mapInfo.removeObject(m, m.cell.location)
		delete(d.monsters, m.uid)
	}
}

// 添加怪物
func (d *Desk) addMonster(m *MonsterInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, ok := d.monsters[m.uid]

	if !ok {
		if d.monBuffStatus != nil {
			if cfg := global.GetDailyMapTime(d.roomCfg.Roomtype); cfg != nil {
				duration := cfg.Duration * 1000
				if duration > 0 {
					m.AddBuff(101, m, int32(duration), d.monBuffStatus, 0, 0)
				}
			}
		}
		if d.roomCfg.Roomtype != constant.WorldBossMap {
			m.currentHP = m.GetMaxHP()
			m.stats.changeStats(pb.Attribute_HEALTH, m.GetMaxHP())
		}

		m.displayHP = m.currentHP

		for i, v := range m.skills {
			if scfg, err := global.GetSkillCfg(v.CfgId); err == nil {
				if scfg.SkillID == constant.SkillTypeNormalAttackMonster {
					continue
				}
				m.skills[i].Cooldown = time.Now().Add(time.Duration(scfg.CoolDown) * time.Millisecond)
			}
		}

		d.monsters[m.uid] = m

		m.setDesk(d)
		m.moveTime = time.Now().Add(time.Second * 2)

		m.spawnLocation = m.cell.location
	}
}

// 刷新怪物
func (d *Desk) refreshMonster() {

	packet := &pb.MonsterRefreshPush{}
	packet.List = make([]*pb.MonsterRefreshPush_Info, 0)

	for _, v := range d.refreshList {
		packet.List = append(packet.List, &pb.MonsterRefreshPush_Info{
			CfgId:       v.CfgId,
			RefreshTime: int32(v.RefreshTime.Unix()),
		})
	}

	d.Broadcast(protocol.RouteMonsterRefreshPush, packet)
}

func (d *Desk) objectByUid(uid int64) IMapObject {
	d.lock.Lock()
	defer d.lock.Unlock()

	p, err := d.playerWithId(uid)
	if err == nil {
		return p
	}
	m, ok := d.monsters[uid]
	if ok {
		return m
	}
	for _, v := range d.objects {
		if v.Uid() == uid {
			return v
		}
	}
	for _, v := range d.robots {
		if v.Uid() == uid {
			return v
		}
	}
	return nil
}

func (d *Desk) monsterWithUid(uid int64) (*MonsterInfo, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, v := range d.monsters {
		if v.uid == uid {
			return v, nil
		}
	}
	return nil, errutil.ErrNotFound
}

func (d *Desk) getMonsters() []IMapObject {
	d.lock.Lock()
	defer d.lock.Unlock()

	var targets []IMapObject
	for _, v := range d.monsters {
		targets = append(targets, v)
	}
	return targets
}

// 地图上随机创建buff
func (d *Desk) createObjectBuff(id int32, location global.PointI, last_time int32) {
	if buffOb := CreateBuffObject(id, ObjectTypeBuff, nil, last_time*1000, 0, nil); buffOb != nil {
		d.addObject(buffOb, location)
	}
}

func (d *Desk) addObject(ob *MapObject, l global.PointI) bool {

	//d.logger.Debugf("开始生成实例 类型%+v UID: %d 名字 %s 坐标 %+v  ID: %d 剩余实例数量 %d", ob.objectType, ob.uid, ob.Name(), l, ob.cfgId, len(d.objects))

	cell := d.mapInfo.getAvailableCellPvp(&l, 1)
	if cell != nil {
		d.mapInfo.addObject(ob, &cell.location)
		ob.SetDesk(d)
		d.objects = append(d.objects, ob)
		ob.bornLocation = cell.location
		ob.spawnLocation = cell.location

		var ownerid_ int64
		if ob.owner != nil {
			ownerid_ = ob.owner.Uid()
		}

		packet := &pb.GameObjectAddPush{
			Info: &pb.ObjectInfo{
				Type:       pb.ObjectType(ob.ObjectType()),
				Uid:        ob.Uid(),
				CfgID:      ob.cfgId,
				Location:   ob.Location(),
				Direction:  ob.GetDirection(),
				Level:      ob.Level(),
				MaxHp:      ob.GetMaxHP(),
				CurHp:      ob.currentHP,
				OwnerUid:   ownerid_,
				BuffList:   ob.getPBBuffList(),
				Attributes: nil,
				Name:       ob.Name(),
			},
		}

		d.Broadcast(protocol.RouteGameObjectAddPush, packet)

		d.logger.Debugf("生成实例 类型%+v UID: %d 名字 %s 坐标 %+v  ID: %d 剩余实例数量 %d", ob.objectType, ob.uid, ob.Name(), ob.location, ob.cfgId, len(d.objects))
		return true
	}
	d.logger.Debugf("生成实例失败 类型%+v UID: %d 名字 %s 坐标 %+v  ID: %d 剩余实例数量 %d", ob.objectType, ob.uid, ob.Name(), l, ob.cfgId, len(d.objects))
	return false
}

func (d *Desk) delObject(uid int64) bool {
	for i, mob := range d.objects {
		if mob.uid == uid {
			if mob.objectType == ObjectTypeDecoy && !mob.dead {
				// if mob.owner.ObjectType() == ObjectTypePlayer {
				// 	owner := mob.owner.(*Player)
				// 	owner.decoyBoom(mob)
				// }
				mob.die(nil)
				return true
			}
			d.mapInfo.removeObject(mob, mob.cell.location)

			d.objects, _ = utils.RemoveElementInPlace(d.objects, i)
			ownerUid_ := int64(0)
			if mob.owner != nil {
				ownerUid_ = mob.owner.Uid()
			}
			packet := &pb.GameObjectDelPush{
				Info: &pb.ObjectInfo{
					Type:      pb.ObjectType(mob.objectType),
					Uid:       uid,
					CfgID:     mob.cfgId,
					Level:     mob.Level(),
					Name:      mob.Name(),
					Location:  mob.location,
					Direction: mob.direction,
					MaxHp:     mob.GetMaxHP(),
					CurHp:     mob.currentHP,
					OwnerUid:  ownerUid_,
				},
			}
			d.Broadcast(protocol.RouteGameObjectDelPush, packet)

			// d.logger.Debugf("删除实例 类型 %+v 坐标 %v  ID: %d  剩余实例数量 %d", mob.objectType, mob.location, mob.cfgId, len(d.objects))
			return true
		}
	}
	return false
}

// 记录行会boos 伤害
func (d *Desk) recordGuildBossDamage(uid int64, damage int64) {
	if d.gr_lock.TryLock() {
		defer d.gr_lock.Unlock()

		d.GuildRewardMember[uid] += damage
	}
}

// 连续击杀记录
func (d *Desk) recordKill(uid int64) {
	if d.roomCfg.Roomtype == constant.PVEWarMap {
		if d.gr_lock.TryLock() {
			defer d.gr_lock.Unlock()
			if _, ok := d.ConsecutiveKills[uid]; ok {
				if time.Now().After(d.ConsecutiveKills[uid].lastKill.Add(time.Second * 60)) {
					d.ConsecutiveKills[uid].killNum = 1
					d.logger.Debugf("战场连斩,超时重置计数  玩家 %d 连斩次数 %d", uid, d.ConsecutiveKills[uid].killNum)
				} else {
					d.ConsecutiveKills[uid].killNum++

					d.logger.Debugf("战场连斩  玩家 %d 连斩次数 %d", uid, d.ConsecutiveKills[uid].killNum)
				}
			} else {
				d.ConsecutiveKills[uid] = &Kills{killNum: 1, lastKill: time.Now()}
			}
			d.ConsecutiveKills[uid].lastKill = time.Now()

			for i := len(global.GlobalCfg.War.Notice.Time) - 1; i >= 0; i-- {
				if d.ConsecutiveKills[uid].killNum == int32(global.GlobalCfg.War.Notice.Time[i].Num) {
					//跑马灯
					if p, err := d.playerWithId(uid); err == nil {
						strSend := []string{p.name, strconv.Itoa(int(d.ConsecutiveKills[uid].killNum))}
						defaultManager.Lamp(int32(constant.LampConsecutiveKills), p.serverID, strSend)
						d.logger.Debugf("战场连斩 跑马灯 玩家 %d 连斩次数 %d", uid, d.ConsecutiveKills[uid].killNum)

						p.pushMsg(protocol.RouteConsecutiveKillsPush, &pb.ConsecutiveKillsPush{Uid: p.uid, Count: d.ConsecutiveKills[uid].killNum})
					}
					break
				}
			}
		}
	}
}

func (d *Desk) addDelayAction(objType ObjectType, objUid int64, skillId int32, target IMapObject, direction pb.Direction, location *global.PointI) {
	if d.status() != constant.DeskStatusPlaying {
		d.logger.Warningf("攻击失败 desk状态不对 deskStatus: %d", d.status())
		return
	}

	select {
	case d.chDelayAttack <- &DelayAction{
		ObjType:   objType,
		ObjUid:    objUid,
		SkillId:   skillId,
		Target:    target,
		Direction: direction,
		Location:  *location,
	}:
	default:
		d.logger.Warningf(" chDelayAttack 通道已满")
	}
}

func (d *Desk) addMonsterRelive(cfgID int32, location global.PointI, lifespan int64, linkedSkill int32, ownershipUid int64, ownershipType pb.OwnershipType, ownershipName string) {
	if d.status() != constant.DeskStatusPlaying {
		d.logger.Warningf("复活怪物失败 desk状态不对 deskStatus: %d", d.status())
		return
	}

	select {
	case d.chMonsterRelive <- &MonRelive{
		CfgID:         cfgID,
		Location:      location,
		Lifespan:      lifespan,
		linkedSkill:   linkedSkill,
		OwnershipType: ownershipType,
		OwnershipUid:  ownershipUid,
		OwnershipName: ownershipName,
	}:
	default:
		d.logger.Warningf(" chMonsterRelive 通道已满")
	}
}

func (d *Desk) skipBattle(uid int64) {

	if time.Now().Before(d.beginTime.Add(time.Duration(global.GlobalCfg.BaseValue.BattleSet.SkipCD-1) * time.Second)) {
		return
	}
	if p, err := d.playerWithId(uid); err == nil {

		gameOverType := constant.GameOverTypeDeadFaild

		// 每一关增加推荐战力B配置，假定当前角色总战力为A，则点击跳过时做如下判断：
		// 如果A小于等于B，则点击跳过必然失败
		// 如果A大于B，则点击跳过成功的几率为：[（A×参数）/B]-1，所得到的值即为成功几率
		if p.combat > d.BattleCfg.RecommendPOW {
			rate := (p.combat*global.GlobalCfg.BaseValue.BattleSet.Constvalue/d.BattleCfg.RecommendPOW - 1) * 100

			if global.RandIn(100) < int(rate) {
				gameOverType = constant.GameOverTypeMDead
			}
		}

		d.logger.Debugf("%s 跳过战役,战力:%d  推荐战力:%d  战斗结果:%v", p.name, p.combat, d.BattleCfg.RecommendPOW, gameOverType)

		d.battleGameOver(gameOverType)
	}
}
