package game

import (
	"math"
	"sort"
	"sync"
	"time"

	"nanoserver/global"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/astar"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
)

type RobotStatus int16

const (
	StatusFree RobotStatus = iota
	StatusReady
	StatusPlaying
	StatusPlayerOnline
	StatusDelete
)

type RobotType int16

const (
	RobotTypeRobot RobotType = iota

	RobotTypePlayer
)

// 1：代表即时赛排行榜
// 2：代表回合赛排行榜
// 3：代表一战到底专用的机器人

type RobotMark int16

const (
	RobotMarkNone RobotMark = iota
	RobotMark1v1
	RobotMarkTurnBased
	RobotMarkFightEnd
)

type RobotInfo struct {
	*Player
	robotId    int32
	status     RobotStatus
	searchTime time.Time
	cfg        *global.RobotConfig
	rank       int64 // 机器人排名
	zoneId     int32
}

type RobotMgr struct {
	lock         sync.Mutex
	robots       map[int64]*RobotInfo
	rank1v1      map[int32]map[int64]*RobotInfo // 1v1排行榜 单服榜6
	rankTurnBase map[int32]map[int64]*RobotInfo // 回合制排行榜 战区榜
}

var robotMgr *RobotMgr

func init() {
	robotMgr = &RobotMgr{robots: make(map[int64]*RobotInfo), rank1v1: make(map[int32]map[int64]*RobotInfo), rankTurnBase: make(map[int32]map[int64]*RobotInfo)}

	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
		robotMgr.rank1v1[i] = make(map[int64]*RobotInfo)
	}
	// todo 回合制榜单按战区存储
	robotMgr.rankTurnBase[0] = make(map[int64]*RobotInfo)
}

func newRobotPlayer(uid int64, name string, sex, job, level, combat int32) *Player {
	r := &Player{
		sex:          int(sex),
		combat:       combat,
		energy:       int64(global.GlobalCfg.BaseValue.Energy.Fullvalue),
		logger:       log.WithField(fieldPlayer, uid),
		bag:          &pb.Bag{ItemList: []*pb.ItemInfo{}},
		curState:     PlayerStatusIdle,
		MapObject:    NewMapObject(ObjectTypeRobot, uid),
		PVPInfo:      &PVPInfo{PVPDuel: &PVPDuel{}, PVPMelee: &PVPMelee{}},
		partnerSlots: [MaxPartnerNum]*PartnerSlot{},
		partners:     map[int64]*Partner{},
		TurnBased:    &TurnBased{},
		InstanceData: &InstanceData{},
		ElementInfo:  &ElementInfo{},
		refine:       NewRefineV2(),
		avatar:       NewAvatar(),
		tokens:       make(map[int32]int64),
	}
	r.tasks = NewTask(r)
	r.name = name
	r.level = level
	r.job = constant.MowJob(job)
	r.PVPDuel.pvpScore = DefaultPvpScore
	r.PVPMelee.gvgScore = DefaultPvpScore
	r.uid = uid
	r.InstanceData.ChallengeRecord.Records = map[int32]int32{}
	r.bRobot = true
	r.ElementInfo.elementMain = int32(pb.ElementSlot_FireEmpower)

	r.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(job).AttackSpeed)
	r.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(job).MoveSpeed)
	r.Artifact = NewArtifact(r)
	r.PetMgr = NewPetMgr(r)
	s := session.New(nil)
	s.Bind(uid)
	r.bindSession(s)
	return r
}

func (rm *RobotMgr) deepCopyRobotPlayer(p *Player) *Player {
	r := newPlayer(nil, p.uid, p.name, "", "")

	r.vip = p.vip
	r.sex = int(p.sex)
	r.level = p.level
	r.job = constant.MowJob(p.job)
	r.bag = p.bag
	r.equips = p.equips
	r.skillsSlot = p.skillsSlot

	r.totalCombat = p.totalCombat
	r.bRobot = true
	r.curState = PlayerStatusIdle
	r.objectType = ObjectTypeRobot
	r.pvpScore = p.pvpScore
	r.gvgScore = p.gvgScore
	r.serverID = p.serverID
	r.robotType = RobotTypePlayer
	r.guildId = p.guildId
	r.elementMain = p.elementMain
	r.guildName = p.guildName
	r.refine = p.refine
	r.avatar = p.avatar
	r.Artifact = p.Artifact

	r.skills = []*PlayerSkill{}
	for _, v := range p.skills {
		r.skills = append(r.skills, &PlayerSkill{CfgId: v.CfgId, AdvancedLevel: v.AdvancedLevel, Cooldown: time.Now(), Distance: v.Distance})
	}

	for k, v := range p.tokens {
		r.tokens[k] = v
	}

	p.stats.Range(func(k pb.Attribute, v int32) bool {
		r.stats.setStats(k, v)
		return true
	})
	//for k, v := range p.stats.copyStats() {
	//	r.stats.setStats(k, v)
	//}

	// 复制伙伴
	r.partners = make(map[int64]*Partner)
	r.partnerSlots = [MaxPartnerNum]*PartnerSlot{}
	for i := 0; i < int(MaxPartnerNum); i++ {
		sid_ := i + 1
		r.partnerSlots[i] = &PartnerSlot{SlotId: int32(sid_)}
		for j := 0; j < int(pb.EquipSlot_MaxEquipSlot); j++ {
			r.partnerSlots[i].Equips[j] = Slot{}
		}
	}

	for _, v := range p.partners {
		basePar := p.getPartnerInfo(v.uid)
		if basePar == nil {
			continue
		}
		copyPar := basePar.deepCopy()
		copyPar.uid = basePar.uid
		copyPar.owner = r

		r.partners[copyPar.uid] = copyPar
	}

	for k, v := range p.partnerSlots {
		if v.PartnerUid == 0 {
			continue
		}

		r.partnerSlots[k].PartnerUid = v.PartnerUid

		for ei := 0; ei < len(v.Equips); ei++ {
			if v.Equips[ei].ItemId == 0 {
				continue
			}
			r.partnerSlots[k].Equips[ei].ItemId = v.Equips[ei].ItemId
			r.partnerSlots[k].Equips[ei].Stars = v.Equips[ei].Stars
			r.partnerSlots[k].Equips[ei].Enhance = v.Equips[ei].Enhance
			r.partnerSlots[k].Equips[ei].Quality = v.Equips[ei].Quality
			r.partnerSlots[k].Equips[ei].Gemstone = v.Equips[ei].Gemstone
		}
	}
	pPet := p.PetMgr.GetBattlePet()
	rPet := r.PetMgr.GetBattlePet()
	if pPet.ItemId > 0 {
		rPet.PetUid = pPet.PetUid
		rPet.PetKind = pPet.PetKind
		rPet.PetCfgId = pPet.PetCfgId
		rPet.Level = pPet.Level
		rPet.Star = pPet.Star
		rPet.Enhance = pPet.Enhance
		rPet.Quality = pPet.Quality
		rPet.Exp = pPet.Exp
		rPet.ExpStar = pPet.ExpStar
		rPet.ItemId = pPet.ItemId
	}
	// 复制宠物
	//if len(pPetMgr.pet) > 0 && pPetMgr.pet[0].ItemId > 0 && len(rPetMgr.pet) > 0 {
	//	rPetMgr.pet[0].PetUid = pPetMgr.pet[0].PetUid
	//	rPetMgr.pet[0].PetKind = pPetMgr.pet[0].PetKind
	//	rPetMgr.pet[0].PetCfgId = pPetMgr.pet[0].PetCfgId
	//	rPetMgr.pet[0].Level = pPetMgr.pet[0].Level
	//	rPetMgr.pet[0].Star = pPetMgr.pet[0].Star
	//	rPetMgr.pet[0].Enhance = pPetMgr.pet[0].Enhance
	//	rPetMgr.pet[0].Quality = pPetMgr.pet[0].Quality
	//	rPetMgr.pet[0].Exp = pPetMgr.pet[0].Exp
	//	rPetMgr.pet[0].ExpStar = pPetMgr.pet[0].ExpStar
	//	rPetMgr.pet[0].ItemId = pPetMgr.pet[0].ItemId
	//}

	for i := 0; i < int(TurnBasedMax); i++ {
		r.turnBasedSlot[i] = p.turnBasedSlot[i]
	}

	for _, v := range p.activateBonds {
		r.activateBonds = append(r.activateBonds, v)
	}

	s := session.New(nil)
	s.Bind(p.uid)
	r.bindSession(s)
	return r
}

func (rm *RobotMgr) loadRobot() {
	if robotMgr == nil {
		robotMgr = &RobotMgr{robots: make(map[int64]*RobotInfo), rank1v1: make(map[int32]map[int64]*RobotInfo), rankTurnBase: make(map[int32]map[int64]*RobotInfo)}

		for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
			robotMgr.rank1v1[i] = make(map[int64]*RobotInfo)
		}
		// todo 回合制榜单按战区存储
		robotMgr.rankTurnBase[0] = make(map[int64]*RobotInfo)
	}

	sort.Slice(global.RobotCfg, func(i, j int) bool {
		return global.RobotCfg[i].Level < global.RobotCfg[j].Level
	})

	for _, v := range global.RobotCfg {
		for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {

			if v.RankMark == int32(RobotMark1v1) {
				rm.createRobotByServerAndNum(i, 0, 1, v)
			} else if v.RankMark == int32(RobotMarkFightEnd) {
				rm.createRobotByServerAndNum(i, 0, 50, v)
			}
		}

		zonelist := global.ServerListMgrInstance.GetLocalServerZoneIds()
		for _, i := range zonelist {
			if v.RankMark == int32(RobotMark1v1) {
				rm.createRobotByServerAndNum(0, i, 1, v)
			} else if v.RankMark == int32(RobotMarkFightEnd) {
				rm.createRobotByServerAndNum(0, i, 50, v)
			} else if v.RankMark == int32(RobotMarkTurnBased) {
				rm.createRobotByServerAndNum(0, i, 1, v)
			}
		}
	}
	logger.Debugf("机器人登录完毕, 共【%d】", len(rm.robots))

	rm.updateAllServerRank()
}

func (rm *RobotMgr) createRobotByServerAndNum(serverId, zoneId, num int32, v *global.RobotConfig) {
	for i := 0; i < int(num); i++ {

		if r := rm.createRobot(v); r != nil {
			r.serverID = serverId
			r.zoneId = zoneId

			rm.robots[r.Uid()] = r
		}
	}
}

// 更新排行榜
func (rm *RobotMgr) rankChangeTCP(r *RobotInfo, k int, v int64) {
	provide.RankChangeTCP(&protocol.RankUserInfoInput{
		Uid:          r.uid,
		Type:         int32(k),
		Avatar:       r.head,
		Name:         r.name,
		Value:        v,
		Sex:          int32(r.sex),
		Job:          int32(r.job),
		Mysql:        true,
		Compare:      k <= constant.RankTypeTrial,
		ServerId:     r.serverID,
		UserServerId: r.serverID,
		RobotCId:     r.robotId,
		ZoneId:       r.zoneId,
	})

	if r.cfgId == 0 || r.robotId == 0 {
		logger.Errorf("rm.rankChangeTCP机器人配置id为0  id: %d  name %s  cfgID %d robotid %d", r.uid, r.name, r.cfgId, r.robotId)
	}
}

// 获取对应等级的机器人，level=0 表示不限制等级
func (rm *RobotMgr) getFreeRobot(level int32) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, v := range rm.robots {
		if v.status != StatusFree {
			continue
		}

		if v.serverID != 0 {
			continue
		}

		if v.cfg.RankMark != int32(RobotMark1v1) {
			continue
		}

		if level == 0 {
			v.status = StatusReady
			return v
		}

		if maputil.ABS(level-v.level) <= 5 {
			v.status = StatusReady
			return v
		}
	}

	for _, v := range global.RobotCfg {
		if maputil.ABS(level-v.Level) <= 5 {
			if r := rm.createRobot(v); r != nil {
				r.status = StatusReady
				return r
			}
		}
	}

	return nil
}

func (rm *RobotMgr) createRobot(v *global.RobotConfig) *RobotInfo {
	info := &RobotInfo{Player: &Player{}}

	uid := snowflake.GenID()
	rName := v.Name
	if rName == "" {
		rName = global.GenRandName()
	}
	info.Player = newRobotPlayer(uid, rName, 1, v.Job, v.Level, v.Power)

	for _, attr := range v.Feature {
		info.stats.setStats(pb.Attribute(attr.Key), attr.Value)
	}
	info.stats.setStats(pb.Attribute_HEALTH, info.stat(pb.Attribute_MAX_HEALTH))

	for _, skill := range v.Skill {
		info.addSkill(int32(skill))
	}

	for i := 0; i < int(MaxPartnerNum); i++ {
		sid_ := i + 1
		info.partnerSlots[i] = &PartnerSlot{SlotId: int32(sid_)}
		for j := 0; j < int(pb.EquipSlot_MaxEquipSlot); j++ {
			info.partnerSlots[i].Equips[j] = Slot{}
		}
	}

	// 伙伴装备到槽位
	slot := int32(1)
	for _, par := range v.Partners {
		if err := info.addPartner(par.Key); err == nil {
			if partner := info.getPartnerInfoByCfgId(par.Key); partner != nil {
				partner.attrRatio = float64(par.Num)
				partner.skills = append(partner.skills, &PlayerSkill{CfgId: par.ID})
				info.partnerCombat(slot, partner.uid)
				slot++
			}
		}
	}

	// 回合制伙伴站位
	tbslot := int32(0)
	info.turnBasedSlot[tbslot] = info.uid
	for _, v := range info.partnerSlots {
		tbslot++
		info.turnBasedSlot[tbslot] = v.PartnerUid
	}

	info.robotId = v.UUID
	info.cfg = v
	info.cfgId = info.cfg.UUID
	info.currentHP = info.stat(pb.Attribute_HEALTH)
	info.robotType = RobotTypeRobot
	info.tokens[constant.ItemTypeOfficialExp] = v.OfficialEXP

	return info
}

func (rm *RobotMgr) getRobotByCfgId(serverId, zoneId int32, cid int32) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, v := range rm.robots {
		if v.robotId == cid && v.serverID == serverId && v.zoneId == zoneId {
			return v
		}
	}
	return nil
}

func (rm *RobotMgr) getUniqueRobotByCfgId(serverID, zoneId, rankType int32, cid int32) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, v := range rm.robots {
		if rankType == constant.RankTypeSoul && v.cfg.RankMark != int32(RobotMark1v1) {
			continue
		}
		if rankType == constant.RankTypeTurnBased && v.cfg.RankMark != int32(RobotMarkTurnBased) {
			continue
		}
		if v.robotId == cid && v.serverID == serverID && v.zoneId == zoneId {

			if _, ok := rm.rank1v1[serverID][v.uid]; ok {
				continue
			}
			if _, ok := rm.rankTurnBase[zoneId][v.uid]; ok {
				continue
			}
			return v
		}
	}
	return nil
}

func (rm *RobotMgr) getFreeRobotByCfgId(cid int32) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	for _, v := range rm.robots {
		if v.robotId == cid && v.status == StatusFree {
			return v
		}
	}

	return nil
}

func (rm *RobotMgr) onGameOver(robotID int64) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	r, ok := rm.robots[robotID]
	if ok {

		logger.Debugf("[%s] -- [机器人类型 : %d] 游戏结束  改变机器人游戏状态 [%d] ", r.name, r.robotType, r.status)
		r.status = StatusFree

		if r.desk != nil {
			r.desk = nil
		}
		r.currentHP = r.stat(pb.Attribute_MAX_HEALTH)
		r.dead = false
		r.stats.setStats(pb.Attribute_HEALTH, r.stat(pb.Attribute_MAX_HEALTH))
	}
}

func (rm *RobotMgr) onExitDesk(serverID int32, robotID int64) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if r, ok := rm.rank1v1[serverID][robotID]; ok {
		r.status = StatusFree
		r.desk = nil
		return
	}
	if r, ok := rm.robots[robotID]; ok {
		r.status = StatusFree
		r.desk = nil
		return
	}
}

func (rm *RobotMgr) onEnterDesk(serverID int32, robotID int64) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if r, ok := rm.rank1v1[serverID][robotID]; ok {
		r.status = StatusPlaying
		return
	}
	if r, ok := rm.robots[robotID]; ok {
		r.status = StatusPlaying
		return
	}
}

func (rm *RobotMgr) OnGameOpMsg(uid int64, route string, msg interface{}) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	r, ok := rm.robots[uid]
	if !ok {
		logger.Warningf("R 不在线 ID : [%d]", uid)
		return
	}

	if r.desk == nil {
		logger.Warningf("R 不在桌子上 ID : [%d]", uid)
		return
	}

	switch route {
	// case protocol.RouteObjectHpChangePush:
	// 	logger.Debugf("R 血量变化 %v", msg.(*pb.ObjectHpChangePush))
	case protocol.RouteMatchRoomJoinPush:
		data := msg.(*pb.MatchRoomJoinPush)
		logger.Debugf("R 加入比赛房间 roomid=%d  rUid =%d  matchID=%d , %v", r.desk.deskID, r.uid, r.desk.matchUid, data)

	default:

	}
}

// 获取排行榜机器人
func (rm *RobotMgr) getRankRobotInfo(serverID int32, uid int64) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if list, ok := rm.rank1v1[serverID]; ok {
		r, ok := list[uid]
		if ok {
			return r
		}
	}

	for _, v := range rm.rankTurnBase {
		if r, ok := v[uid]; ok {
			return r
		}
	}

	if r, ok := rm.robots[uid]; ok {
		return r
	}

	return nil
}

// 获取排行榜机器人
func (rm *RobotMgr) getRankRobotInfoByRank(rankType int32, serverID int32, rank int64) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	switch rankType {
	case constant.RankTypeSoul:
		if list, ok := rm.rank1v1[serverID]; ok {
			for _, v := range list {
				if v.rank == rank {
					return v
				}
			}
			logger.Warningf("按排名获取机器人失败, 榜上无此机排名 RankType= %d serverID = %d , rank = %d", rankType, serverID, rank)
		} else {
			logger.Errorf("按排名获取机器人失败,服务器榜单异常 RankType= %d serverID = %d , rank = %d", rankType, serverID, rank)
		}

	case constant.RankTypeTurnBased:

		if serverInfo, err := global.ServerListMgrInstance.Get(int(serverID)); err == nil {
			logger.Errorf("获取回合制机器人, 战区ID [%d], RankType= %d serverID = %d , rank = %d ", serverInfo.ZoneId, rankType, serverID, rank)
			if list, ok := rm.rankTurnBase[serverInfo.ZoneId]; ok {
				for _, v := range list {
					if v.rank == rank {
						return v
					}
				}
			}
		} else {
			logger.Errorf("获取回合制机器人失败, 战区ID获取失败, RankType= %d serverID = %d , rank = %d , err=%v", rankType, serverID, rank, err)
		}

	}

	return nil
}

// 获取排行榜机器人
func (rm *RobotMgr) getRankRobotInfoByUid(rankType int32, serverID int32, uid int64) *RobotInfo {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	switch rankType {
	case constant.RankTypeSoul:
		if list, ok := rm.rank1v1[serverID]; ok {
			if r, ok := list[uid]; ok {
				return r
			}
		}
	case constant.RankTypeTurnBased:
		for _, v := range rm.rankTurnBase {
			if r, ok := v[uid]; ok {
				return r
			}
		}
		//logger.Errorf("获取回合制机器人失败, RankType= %d serverID = %d , uid = %d", rankType, serverID, uid)
	}

	return nil
}

// 获取排行榜十个名次相近的对手
func (rm *RobotMgr) getTenClosestPlayers(rankType int32, serverID int32, uid int64) []*RobotInfo {
	var players []*RobotInfo
	self := rm.getRankRobotInfoByUid(rankType, serverID, uid)
	if self == nil {
		logger.Debugf("自己不在排行榜上 [%d]  服务器ID [%d]", uid, serverID)
		for i := 100; i > 90; i-- {
			if rp := rm.getRankRobotInfoByRank(rankType, serverID, int64(i)); rp != nil {
				if rp.uid == uid {
					continue
				}
				exist := false
				for _, v := range players {
					if v.uid == rp.uid {
						exist = true
						break
					}
				}
				if !exist {
					players = append(players, rp)
				}
			}
		}
	} else {
		logger.Debugf("自己在排行榜上 [%d] 名次 %d 服务器ID [%d]", uid, self.rank, serverID)
		// 获取前面的玩家
		for i := self.rank - 1; i > 0 && len(players) < 5; i-- {
			if rp := rm.getRankRobotInfoByRank(rankType, serverID, int64(i)); rp != nil {
				if rp.curState != PlayerStatus(StatusFree) {
					continue
				}
				exist := false
				for _, v := range players {
					if v.uid == rp.uid {
						exist = true
						break
					}
				}
				if !exist {
					players = append(players, rp)
				}
			}
		}
		// 获取后面的玩家
		for i := self.rank + 1; i <= int64(len(rm.rank1v1[serverID])) && len(players) < 10; i++ {
			if rp := rm.getRankRobotInfoByRank(rankType, serverID, i); rp != nil {
				if rp.curState != PlayerStatus(StatusFree) {
					continue
				}
				exist := false
				for _, v := range players {
					if v.uid == rp.uid {
						exist = true
						break
					}
				}
				if !exist {
					players = append(players, rp)
				}
			}
		}
	}

	//如果还是不满10人
	for next := 100; len(players) < 10 && next > 1; next-- {
		if rp := rm.getRankRobotInfoByRank(rankType, serverID, int64(next)); rp != nil {
			exist := false
			for _, v := range players {
				if v.uid == rp.uid {
					exist = true
					break
				}
			}
			if !exist {
				players = append(players, rp)
			}
		}
	}

	return players
}

func (rm *RobotMgr) updateAllServerRank() {
	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
		robotMgr.rank1v1[i] = make(map[int64]*RobotInfo)
	}

	zonelist := global.ServerListMgrInstance.GetLocalServerZoneIds()
	for _, i := range zonelist {
		robotMgr.rankTurnBase[i] = make(map[int64]*RobotInfo)
	}

	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {
		rm.updateRank(int32(i), 0, constant.RankTypeSoul)

		logger.Debugf("新增即时赛排行榜机器人, 服务器ID [%d]", i)
	}

	for _, i := range zonelist {
		rm.updateRank(0, i, constant.RankTypeTurnBased)
		logger.Debugf("新增回合制排行榜机器人, 战区ID [%d]", i)
	}

	for i := global.ServerIDRangeStart; i <= global.ServerIDRangeEnd; i++ {

		time.AfterFunc(time.Millisecond*time.Duration(global.RandIn(1000)+500), func() {
			rm.updateRank(int32(i), 0, constant.RankTypeSoul)
		})
		logger.Debugf("更新即时赛排行榜机器人, 服务器ID [%d]", i)
	}

	for _, i := range zonelist {
		time.AfterFunc(time.Millisecond*time.Duration(global.RandIn(1000)+500), func() {
			rm.updateRank(0, i, constant.RankTypeTurnBased)
		})
		logger.Debugf("更新回合制排行榜机器人, 战区ID [%d]", i)
	}
}

func (rm *RobotMgr) updateRank(serverId int32, zoneId int32, rankType int64) {
	// 加载排行榜
	req := protocol.RankOperateRequest{Type: rankType, All: 1, ServerId: serverId, ZoneId: zoneId}
	if rankList, err := provide.RankList(&req); err == nil {
		if len(rankList.RankAll) > 0 {
			for _, v := range rankList.RankAll {
				if v.RobotCId > 0 {
					if r := rm.getRankRobotInfoByUid(int32(rankType), serverId, v.Uid); r != nil {
						r.rank = v.Num
						r.name = v.Name
						r.pvpScore = int32(v.Value)
						r.robotType = RobotTypeRobot
						r.uid = v.Uid
						r.job = constant.MowJob(v.Job)
						r.sex = int(v.Sex)
						r.cfgId = v.RobotCId
						r.serverID = serverId
						r.zoneId = zoneId

						s := session.New(nil)
						s.Bind(r.uid)
						r.bindSession(s)

						rm.setRankRobot(serverId, zoneId, rankType, r)
						r.refreshStats()

					} else if r := rm.getUniqueRobotByCfgId(serverId, zoneId, int32(rankType), v.RobotCId); r != nil {
						r.rank = v.Num
						r.name = v.Name
						r.pvpScore = int32(v.Value)
						r.robotType = RobotTypeRobot
						r.uid = v.Uid
						r.job = constant.MowJob(v.Job)
						r.sex = int(v.Sex)
						r.cfgId = v.RobotCId
						r.serverID = serverId
						r.zoneId = zoneId

						tbslot := int32(0)
						r.turnBasedSlot[tbslot] = r.uid
						for _, v := range r.partnerSlots {
							tbslot++
							r.turnBasedSlot[tbslot] = v.PartnerUid
						}

						s := session.New(nil)
						s.Bind(r.uid)
						r.bindSession(s)

						rm.setRankRobot(serverId, zoneId, rankType, r)
						r.refreshStats()
						// logger.Debugf(" 服务器ID: %d 排行榜机器人数据 UID : %d  Name : %s 积分r:%d  Rank : %d", serverId, r.Uid(), r.name, r.pvpScore, r.rank)
					}
				} else {
					if player := defaultManager.ObserveUserBaseInfo(serverId, v.Uid); player != nil {

						if r := rm.getRankRobotInfoByUid(int32(rankType), serverId, v.Uid); r != nil {
							if r.desk != nil {
								r.rank = v.Num
								logger.Debugf("正在战斗的机器人不刷新属性  %d %s ", r.uid, r.name)
								continue
							}
							r.Player = rm.deepCopyRobotPlayer(player)
							r.defaultBuff()
							r.serverID = serverId
							r.zoneId = zoneId
							if rankType == constant.RankTypeSoul {
								r.pvpScore = int32(v.Value)
							} else if rankType == constant.RankTypeTurnBased {
								r.gvgScore = int32(v.Value)
							}
							player.setScore(v.Value, int(rankType))
							r.rank = v.Num
							continue
						}

						dp := rm.deepCopyRobotPlayer(player)
						r := &RobotInfo{Player: dp}
						r.rank = v.Num
						r.defaultBuff()

						if rankType == constant.RankTypeSoul {
							r.pvpScore = int32(v.Value)
						} else if rankType == constant.RankTypeTurnBased {
							r.gvgScore = int32(v.Value)
						}
						player.setScore(v.Value, int(rankType))

						r.robotType = RobotTypePlayer
						r.job = constant.MowJob(v.Job)
						r.sex = int(v.Sex)

						r.serverID = serverId
						r.zoneId = zoneId

						rm.setRankRobot(serverId, zoneId, rankType, r)

						logger.Debugf(" 服务器ID: %d 排行榜真实玩家数据 UID : %d  Name : %s 积分r:%d  Rank : %d", serverId, r.Uid(), r.name, r.pvpScore, r.rank)
					}
				}
			}
		} else {
			rankRobot := []*RobotInfo{}

			if rankType == constant.RankTypeSoul {
				for _, v := range rm.robots {
					if v.cfg.RankMark == int32(RobotMark1v1) && v.serverID == serverId {
						rankRobot = append(rankRobot, v)
						if len(rankRobot) >= 100 {
							break
						}
					}
				}
			} else if rankType == constant.RankTypeTurnBased {
				for _, v := range rm.robots {
					if v.cfg.RankMark == int32(RobotMarkTurnBased) && v.zoneId == zoneId {
						rankRobot = append(rankRobot, v)
						if len(rankRobot) >= 100 {
							break
						}
					}
				}
			}
			sort.Slice(rankRobot, func(i, j int) bool {
				return rankRobot[i].cfg.Level > rankRobot[j].cfg.Level
			})

			for i := 0; i < len(rankRobot); i++ {
				rankRobot[i].pvpScore = DefaultPvpScore
				rankRobot[i].gvgScore = DefaultPvpScore

				rm.rankChangeTCP(rankRobot[i], int(rankType), int64(rankRobot[i].pvpScore))

				// logger.Debugf("服务器ID:%d 战区ID:%d 首次加载排行榜， 加入机器人 名次 %d  名字 %s  UID: %d  CfgID: %d  积分 : %d  性别 %d  职业 %d",
				// 	serverId, zoneId, i, rankRobot[i].name, rankRobot[i].uid, rankRobot[i].robotId, rankRobot[i].pvpScore, rankRobot[i].sex, rankRobot[i].job)
				time.Sleep(time.Millisecond)
			}

		}
	}

}

func (rm *RobotMgr) setRankRobot(serverId int32, zoneId int32, rankType int64, r *RobotInfo) {
	switch rankType {
	case constant.RankTypeSoul:
		{
			if _, ok := rm.rank1v1[serverId]; ok {
				rm.rank1v1[serverId][r.uid] = r
			} else {
				rm.rank1v1[serverId] = make(map[int64]*RobotInfo)
				rm.rank1v1[serverId][r.uid] = r
			}

		}
	case constant.RankTypeTurnBased:
		{
			if _, ok := rm.rankTurnBase[zoneId]; ok {
				rm.rankTurnBase[zoneId][r.uid] = r
			} else {
				rm.rankTurnBase[zoneId] = make(map[int64]*RobotInfo)
				rm.rankTurnBase[zoneId][r.uid] = r
			}
		}
	}
}

// 更新排行榜玩家数据
func (rm *RobotMgr) updateRankPlayer(player *Player) {
	// 更新即时赛
	if r := rm.getRankRobotInfoByUid(constant.RankTypeSoul, player.serverID, player.uid); r != nil {
		if r.desk == nil && r.status == StatusFree && r.robotType == RobotTypePlayer {
			r.Player = rm.deepCopyRobotPlayer(player)
			r.defaultBuff()
			r.totalCombat = player.totalCombat
		}
	}

	if r := rm.getRankRobotInfoByUid(constant.RankTypeTurnBased, 0, player.uid); r != nil {
		if r.desk == nil && r.status == StatusFree && r.robotType == RobotTypePlayer {
			r.Player = rm.deepCopyRobotPlayer(player)
			r.defaultBuff()
			r.totalCombat = player.totalCombat
		}
	}
}

/*
---------------------------------------
   RobotMgr End
---------------------------------------
*/

// 获取攻击距离
func (r *RobotInfo) getAttackDistance() int {
	d := 2
	switch r.job {
	case constant.Warrior, constant.PartnerWarrior:
		d = 2

	case constant.PartnerWizard, constant.PartnerRitualist:
		d = 5

	case constant.Wizard, constant.Ritualist:
		d = 7
	}

	return d
}

func (r *RobotInfo) processAI() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	r.processSearch()
	r.processTarget()
	r.processBuff()
	r.processHp()
}

// 寻找
func (r *RobotInfo) processSearch() {
	if r.desk == nil {
		return
	}

	if time.Now().Before(r.searchTime) {
		return
	}
	r.searchTime = time.Now().Add(time.Second)

	if !r.CanAttack() {
		return
	}

	// 被嘲讽不能转换目标
	if r.getBuffByType(constant.BuffTaunt) != nil && r.target != nil {
		return
	}

	if r.Target() != nil && maputil.Distance(r.Location(), r.Target().Location()) <= 3 {
		// logger.Debugf("R 目标不变 自己ID: %d , 自身坐标:%v, 目标ID: %d  目标坐标: %v 目标类型: %v", r.uid, r.location, r.Target().Uid(), r.target.Location(), r.target.ObjectType())
		return
	}

	distance := math.MaxInt
	var closest IMapObject
	if r.desk.roomCfg.Roomtype == constant.PVEPeacefulMap || r.desk.roomCfg.Roomtype == constant.PVEWarMap {
		for _, v := range r.desk.monsters {
			if v.dead {
				continue
			}
			if v.uid == r.uid {
				continue
			}
			tmpDistance := maputil.Distance(r.location, v.Location())
			if tmpDistance < distance {
				distance = tmpDistance
				closest = v
			}
		}
	} else {
		for _, v := range r.desk.objects {
			if v.objectType == ObjectTypeDecoy || v.objectType == ObjectTypePartner {
				tmpDistance := maputil.Distance(r.location, v.Location())
				if tmpDistance > r.getAttackDistance()+20 {
					continue
				}
				if !r.CanAttackTarget(v) {
					continue
				}

				if tmpDistance < distance {
					distance = tmpDistance
					closest = v
				}
			}
		}

		for _, v := range r.desk.players {
			if v.dead {
				continue
			}
			if v.uid == r.uid {
				continue
			}
			if v.isRobot() {
				continue
			}
			tmpDistance := maputil.Distance(r.location, v.Location())
			if tmpDistance < distance {
				distance = tmpDistance
				closest = v
			}
		}

		for _, v := range r.desk.monsters {
			if v.dead {
				continue
			}
			if v.uid == r.uid {
				continue
			}
			tmpDistance := maputil.Distance(r.location, v.Location())
			if tmpDistance < distance {
				distance = tmpDistance
				closest = v
			}
		}

	}

	if closest != nil {
		r.target = closest

		//logger.Debugf("R 锁定目标 mUid: %d   自身坐标:%v, 目标ID: %d 目标名字 %s, 目标坐标: %v 目标类型: %v", r.uid, r.location, r.target.Uid(), r.target.Name(), r.target.Location(), r.target.ObjectType())
	}
}

// 目标
func (r *RobotInfo) processTarget() {
	if r.desk == nil {
		return
	}

	if r.target == nil {
		return
	}

	if r.dead || r.target.Dead() {
		r.target = nil
		return
	}

	//战场安全区内不攻击
	if r.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if r.desk.mapInfo.safeArea(r.target.Location()) {
			return
		}
	}

	if !r.CanAttackTarget(r.target) {
		return
	}

	var castSkill *PlayerSkill
	for _, v := range r.skills {
		skillCfg, err := global.GetSkillCfg(int32(v.CfgId))
		if err != nil {
			return
		}

		if skillCfg.SType == constant.SkillTypeNormal || skillCfg.SType == constant.SkillTypeArtifact {
			continue
		}

		if skillCfg.ActiveorPassive == PassiveSkill {
			continue
		}

		if skillCfg.CoolDown > 0 && time.Now().Before(v.Cooldown) {
			continue
		}

		if skillCfg.Effectnum < 3 {
			castSkill = v
			break
		} else {
			if r.inAttackRange(int(skillCfg.MaxDistance)) {
				castSkill = v
				break
			}
		}
	}

	distance := r.getAttackDistance()

	if !r.inAttackRange(distance) {
		if !r.CanMove() {
			return
		}
		var notargetSkill *PlayerSkill
		for _, v := range r.skills {
			skillCfg, err := global.GetSkillCfg(int32(v.CfgId))
			if err != nil {
				return
			}

			if skillCfg.SType == constant.SkillTypeNormal || skillCfg.SType == constant.SkillTypeArtifact {
				continue
			}

			if skillCfg.CoolDown > 0 && time.Now().Before(v.Cooldown) {
				continue
			}
			if skillCfg.SkillID == constant.SkillTypeDecoy || skillCfg.SkillID == constant.SkillTypeDecoyBoss || skillCfg.SkillID == constant.SkillTypeInstantaneousMove {
				continue
			}

			if skillCfg.SkillID == constant.SkillTypeSwiftStrike {
				notargetSkill = v
				break
			}
			// 30% 几率给自己上buff
			if skillCfg.Needaim == 0 && global.RandIn(100) < 30 {
				notargetSkill = v
				break
			}
		}
		if notargetSkill != nil {
			castSkill = notargetSkill
		} else {
			if path := maputil.GetPath(r.location, r.target.Location(), &r.desk.mapCfg.AstarMgr); len(path) > 1 {
				r.moveTo(path[1])
			} else {
				logger.Debugf("R 寻路失败 起点=%v  终点=%v", r.location, r.target.Location())
			}
			return
		}

	} else {
		if (r.job == constant.Wizard || r.job == constant.Ritualist) && maputil.Distance(r.location, r.target.Location()) < 2 && r.getBuffByType(constant.BuffKnockback) == nil {
			if r.CanMove() && r.getBuffByType(constant.BuffTaunt) == nil {
				var nearSkill *PlayerSkill
				for _, v := range r.skills {
					skillCfg, err := global.GetSkillCfg(int32(v.CfgId))
					if err != nil {
						return
					}
					if skillCfg.SType == constant.SkillTypeArtifact {
						continue
					}

					if skillCfg.CoolDown > 0 && time.Now().Before(v.Cooldown) {
						continue
					}

					if skillCfg.SkillID == constant.SkillTypeDecoy || skillCfg.SkillID == constant.SkillTypeDecoyBoss || skillCfg.SkillID == constant.SkillTypeInstantaneousMove {
						nearSkill = v
						break
					}

					if skillCfg.SkillID == constant.SkillTypeManaExpel || skillCfg.SkillID == constant.SkillTypeMagicBarrier {
						nearSkill = v
						break
					}
				}

				if nearSkill != nil {
					castSkill = nearSkill
				} else {
					var cell *Cell
					var point *global.PointI
					dir := maputil.DirectionFromPointI(&r.cell.location, &r.target.GetCell().location)
					point = maputil.Move(&r.cell.location, maputil.ShiftDirection(dir, 4), 1)
					if cell = r.desk.mapInfo.getCell(point); cell != nil && !cell.block {
						r.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
						r.moveTime = r.moveTime.Add(time.Duration(r.stat(pb.Attribute_MOVE_SPEED)*3) * time.Millisecond)
						return
					}
					point = maputil.Move(&r.cell.location, maputil.ShiftDirection(dir, 2), 1)
					if cell = r.desk.mapInfo.getCell(point); cell != nil && !cell.block {
						r.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
						r.moveTime = r.moveTime.Add(time.Duration(r.stat(pb.Attribute_MOVE_SPEED)*2) * time.Millisecond)
						return
					}
					point = maputil.Move(&r.cell.location, maputil.ShiftDirection(dir, -2), 1)
					if cell = r.desk.mapInfo.getCell(point); cell != nil && !cell.block {
						r.moveTo(&astar.APoint{X: cell.location.X, Y: cell.location.Y})
						r.moveTime = r.moveTime.Add(time.Duration(r.stat(pb.Attribute_MOVE_SPEED)*2) * time.Millisecond)
						return
					}
				}
			}
		}
		if len(r.cell.objects) > 1 {
			cell := r.desk.mapInfo.getAvailableCell(global.ToPointI(r.target.Location()), int32(r.getAttackDistance()))
			if cell != nil {
				if path := maputil.GetPath(r.location, &pb.Point{X: float64(cell.location.X), Y: float64(cell.location.Y)}, &r.desk.mapCfg.AstarMgr); len(path) > 1 {
					r.moveTo(path[1])
					return
				}
			}
		}
	}
	r.attack(castSkill)
}

// 移动
func (r *RobotInfo) moveTo(point *astar.APoint) {
	if r.desk == nil {
		return
	}

	if !r.CanMove() {
		return
	}
	// from := r.location
	to := &pb.Point{X: float64(point.X), Y: float64(point.Y)}
	if !r.desk.mapInfo.locationChange(r, to) {
		return
	}

	r.moveTime = time.Now().Add(time.Duration(r.getMoveSpeed()+100) * time.Millisecond)

	packet := &pb.MovePush{Uid: r.uid, Location: r.location, Direction: r.direction}
	r.desk.Broadcast(protocol.RouteMovePush, packet)
	r.SetObAction(ObjectActionMove)

	// logger.Debugf("%s 移动 rUid: %d , 起点:%v , 目标: %v  与目标距离 %d  移动速度 %d", r.name, r.uid, from, to, maputil.Distance(r.location, r.target.Location()), r.getMoveSpeed())
}

// 攻击
func (r *RobotInfo) attack(castSkill *PlayerSkill) {
	if r.desk == nil {
		logger.Errorf("机器人攻击失败，不在房间中 %d:%s", r.uid, r.name)
		return
	}

	if !r.CanAttack() {
		return
	}

	if !r.thresholdDetect() {
		return
	}

	var castSkillId, skillIndex int32
	if castSkill != nil {
		castSkillId = castSkill.CfgId
	} else {
		castSkillId = global.GetNormalSkillId(r.job)
	}
	skillCfg, err := global.GetSkillCfg(int32(castSkillId))
	if err != nil {
		return
	}
	for i := range r.skills {
		if r.skills[i].CfgId == castSkillId {
			skillIndex = int32(i)
		}
	}

	Cooldown := int32(0)
	Radius := skillCfg.Radius
	if castSkill != nil {
		if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, castSkill.AdvancedLevel); err == nil {
			Cooldown = skillUpperlevel.Cooldown
			Radius += skillUpperlevel.Distance
		}
	}

	attackCells := []*Cell{}
	targets := []IMapObject{}
	targets = append(targets, r.target)

	//优先设置攻击间隔
	r.setAttackTime(global.GetActionsetCfg(skillCfg.SkillID))

	// 技能影响到的全部格子
	switch skillCfg.SkillID {

	case constant.SkillTypeHalfMoon: // 扇形
		attackCells = r.desk.mapInfo.getSectorCells(&r.cell.location, r.direction, Radius)

	case constant.SkillTypeFourSides: // 四方向
		for i := 0; i < 4; i++ {
			direction := maputil.ShiftDirection(r.direction, i*2)
			attackCells = append(attackCells, r.desk.mapInfo.getLineCells(&r.cell.location, direction, Radius)...)
		}

	case constant.SkillTypeFreezeTime, constant.SkillTypeTimeSpaceLock, constant.SkillTypeDuel, constant.SkillTypeRepelFoe,
		constant.SkillTypeMassHeal, constant.SkillTypeDivineStrength, constant.SkillTypeArcaneMastery, constant.SkillTypeBlessingOfCourage,
		constant.SkillTypeNetherSpell, constant.SkillTypeDivineBattleArmor, constant.SkillTypeTaunt, constant.SkillTypeStomp,
		constant.SkillTypeBleeding, constant.SkillTypeWreck, constant.SkillTypeApocalypse: // 正方形
		attackCells = r.desk.mapInfo.getSquareCells(&r.cell.location, Radius)

	case constant.SkillTypeBreakArmy, constant.SkillTypeAuroraPulse, constant.SkillTypeEarthquakeWave: // 长方形
		attackCells = r.desk.mapInfo.getLineCells(&r.cell.location, r.direction, Radius)
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
				if r.CanAttackTarget(obj) && r.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == r.uid {
			continue
		}

		// 安全区内不能攻击玩家
		if r.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (r.desk.mapInfo.safeArea(r.location) || r.desk.mapInfo.safeArea(t.Location())) {
				continue
			}
		}

		// 技能距离
		if skillCfg.SType != constant.SkillTypeNormal && skillCfg.MaxDistance > 0 && maputil.Distance(r.location, t.Location()) > int(skillCfg.MaxDistance) {
			continue
		}
		//伤害延迟生效
		pardesk := r.desk
		part := t
		ilocation := &r.cell.location
		time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
			if pardesk != nil {
				pardesk.addDelayAction(ObjectTypeRobot, r.uid, castSkillId, part, r.direction, ilocation)
			}
		})
	}

	r.direction = maputil.DirectionFromPointF(r.location, r.target.Location())

	packet := &pb.AttackPush{
		AttackerId: r.uid,
		TargetUids: []int64{r.target.Uid()},
		SkillId:    int32(castSkillId),
		Location:   r.target.Location(),
		Direction:  r.direction,
	}
	r.desk.Broadcast(protocol.RouteAttackPush, packet)

	// 技能CD
	if skillCfg.CoolDown > 0 {
		CD := skillCfg.CoolDown - Cooldown
		if CD < 0 {
			CD = 0
		}
		r.skills[skillIndex].Cooldown = time.Now().Add(time.Duration(CD) * time.Millisecond)
	}

	r.SkillBuffTrigger(castSkill)

	r.SetObAction(ObjectActionAttack)
	//logger.Debugf("%s 攻击 attack rUid: %d , 目标: %d  技能 %d  ", r.name, r.uid, r.Target().Uid(), skillId)
}

// 计算角色属性
func (r *RobotInfo) refreshStats() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	if !r.slock.TryLock() {
		return
	}
	defer r.slock.Unlock()

	now := time.Now()
	r.stats.clear()
	// 基础属性

	if r.robotType == RobotTypePlayer {
		roleCfg, err := global.GetRoleCfg(int32(r.job), r.level)
		if err != nil {
			logger.Errorf("获取角色等级配置失败 level %d", r.level)
			return
		}

		for _, rc := range roleCfg.BassProperty {
			r.stats.setStats(pb.Attribute(rc.Key), rc.Value)
		}

		r.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(r.job)).AttackSpeed)
		r.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(r.job)).MoveSpeed)
	} else {
		if r.cfg != nil {
			for _, attr := range r.cfg.Feature {
				r.stats.setStats(pb.Attribute(attr.Key), attr.Value)
			}

			r.stats.changeStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(r.job)).AttackSpeed)
			r.stats.changeStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(r.job)).MoveSpeed)

			// 刷新伙伴属性
			for _, v := range r.partnerSlots {
				if v == nil {
					continue
				}
				partner := r.getPartnerInfo(v.PartnerUid)
				if partner == nil {
					continue
				}
				partner.refreshStats(nil, r.Player)
			}
		}
	}

	r.stats.setStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE, int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	// 魂石增加生命值基数=角色基础属性+命格
	gemHealthAdd := r.stat(pb.Attribute_MAX_HEALTH)

	// 装备 强化大师
	equipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	gemStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	petStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	newpetStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// guildStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// guildSkillStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// elementBaseStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// elementAdvanceStatus := Stats{Attrs: make(map[pb.Attribute]int32)}

	var maxEnhance, maxStars int32
	maxEnhance = math.MaxInt32
	maxStars = math.MaxInt32
	var starCfg, enhanceCfg *global.EquipMasterConfig

	// 宠物
	// 宠物因为有些属性加给装备，所以提到前面来了
	//if len(r.pet) > 0 {
	//	if r.pet[0].PetCfgId != 0 {
	//		petCfg := global.GetPetCfg(r.pet[0].PetCfgId)
	//		if petCfg != nil {
	//
	//			for _, e2 := range petCfg.StarAttribute {
	//				if e2.Id == r.pet[0].Star {
	//					for _, e3 := range e2.Values {
	//						petStatus.setStats(pb.Attribute(e3.Key), e3.Value)
	//					}
	//					break
	//				}
	//			}
	//
	//			roleCfg, _ := global.GetRoleCfg(constant.Pet, r.pet[0].Level)
	//			if roleCfg != nil {
	//				rate := global.GetpetQualityRate(r.pet[0].Quality)
	//				if rate == 0 || petCfg.PetKind.Key == 2 {
	//					rate = 1
	//				}
	//				for _, e := range roleCfg.BassProperty {
	//					petStatus.setStats(pb.Attribute(e.Key), int32(math.Round(float64(e.Value)*rate)))
	//				}
	//			}
	//		}
	//	}
	//}
	r.PetMgr.setStats(&petStatus)

	// 计算属性
	newpetStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(petStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(1+(float64(petStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(petStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(1+(float64(petStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(petStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(1+float64(petStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	newpetStatus.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(petStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(1+(float64(petStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_MAX_HEALTH, int32(float64(petStatus.getStats(pb.Attribute_MAX_HEALTH))*(1+(float64(petStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))))
	for pk, pv := range newpetStatus.Attrs {
		if pk == pb.Attribute_PHYSICAL_ATTACK || pk == pb.Attribute_MAGIC_ATTACK ||
			pk == pb.Attribute_PHYSICAL_DEFENSE || pk == pb.Attribute_MAGIC_DEFENSE ||
			pk == pb.Attribute_MAX_HEALTH {
			r.stats.setStats(pk, pv)
		}
	}

	allGems := []int32{}
	for slot, equip := range r.equips {
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
			// 对应宝石等级的数目加1
			allGems = append(allGems, genstoneCfg.Itemlevel)

			for _, attr := range genstoneCfg.Stats {
				if attr.Key == pb.Attribute_HEALTH_PERCENTAGE {
					gemStatus.setStats(attr.Key, attr.Value)
				} else {
					r.stats.setStats(attr.Key, attr.Value)
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
				r.stats.setStats(pb.Attribute(v.Key), v.Value)
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
		equipInfo := r.getItemByUid(equip.ItemId)
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
				if equipInfo.Equip != nil {
					// 如果是武器，会额外算上武器基础值
					if pb.EquipSlot(equipCfg.Equiptype) == pb.EquipSlot_Weapon {
						per := float64(1)
						// 获取某个buff的属性
						if buff := r.getBuffByType(constant.BuffPetWeaponAttackPercentageIncrease); buff != nil {
							per = float64(buff.BuffStats.getStats(pb.Attribute_WEAPON_ATTACK_PERCENTAGE)) / BaseRate
						}
						equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))*per))
					} else {
						equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))))
					}
				}
			} else {
				r.stats.setStats(attr.Key, attr.Value)
			}
		}
		// 极品属性直接计算
		if equipInfo.Equip != nil {
			for _, attr := range equipInfo.Equip.Attributes {
				r.stats.setStats(attr.Attr, attr.Num)
			}
		}

		star_uid := int32(r.job%10*100) + equip.Stars
		if equipStarCfg, err := global.GetStarsCfg(star_uid); err == nil {
			for _, star_attr := range equipStarCfg.Attribute {
				if star_attr.Slot == int32(slot) {
					equipStatus.setStats(pb.Attribute(star_attr.Attribute), star_attr.Value)
				}
			}
		}

		enhanc_uid := int32(r.job%10*1000) + equip.Enhance
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
			r.stats.setStats(ek, ev)
		}
	}

	// 魂石增加生命百分比
	r.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(gemHealthAdd)*(float64(gemStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryEnhance && maxEnhance >= equipmaster.Limit {
			enhanceCfg = equipmaster
		}
		if equipmaster.Category == constant.EquipMasterCategoryStars && maxStars >= equipmaster.Limit {
			starCfg = equipmaster
		}
	}

	// 魂石大师
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
					r.stats.setStats(pb.Attribute(v.Key), v.Value)
				}
				break
			}
		}
	}

	// 强化大师
	if enhanceCfg != nil {
		for _, v := range enhanceCfg.Attribute {
			r.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	// 升星大师
	if starCfg != nil {
		for _, v := range starCfg.Attribute {
			r.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	// 特权卡
	privilegeStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, v := range r.privilege {
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
		r.stats.setStats(k, v)
	}

	// VIP
	vipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, gv := range global.VIPCfg {
		if gv.VIPLevel <= r.vip {
			for _, att := range gv.Pcontent {
				av := vipStatus.getStats(pb.Attribute(att.Key))
				if att.Value > av {
					vipStatus.changeStats(pb.Attribute(att.Key), att.Value)
				}
			}
		}
	}
	for k, v := range vipStatus.Attrs {
		r.stats.setStats(k, v)
	}

	r.stats.setStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE, int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	// 工会基础属性
	if r.guildId > 0 {
		// 工会技能
		// guildManager := NewGuildManager()
		guild, err := guildManager.validateGuild(r.guildId)
		if err == nil {
			guildCfg := global.GetGuildCfg(int32(guild.Guild.Level))
			if guildCfg != nil {
				for _, e := range guildCfg.AttributeAdd {
					r.stats.setStats(pb.Attribute(e.Key), e.Value)
				}
			}
			// 查询工会目前技能情况
			for _, e := range guild.Skills {
				skillCfg := global.GetGuildSkillCfg(e.UUid)
				if skillCfg != nil {
					for _, e2 := range skillCfg.ValueAdd {
						r.stats.setStats(pb.Attribute(e2.Key), e2.Value)
					}
				}
			}
		}

	}
	// 被动技能
	for _, v := range r.skills {
		if skillCfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			if skillCfg.ActiveorPassive == PassiveSkill {
				for _, e := range skillCfg.Effect {
					r.stats.setStats(pb.Attribute(e.Key), e.Value)
				}
			}
		}
	}
	// buff
	petBuffs := r.getPetBuffs()
	//if r.block.TryLock() {
	//logger.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxx", r.buffs)
	r.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
		buffCfg, err := global.GetBuffCfg(buff.Id)
		if err == nil {
			// buff白名单，有的buff不计算战斗力，
			// 判断是不是在白名单里面
			// 获取宠物buff
			if len(petBuffs) > 0 {
				for _, petBuff := range petBuffs {
					// 获取buffid对应的大类
					buffCfgDetail, err := global.GetBuffCfg(petBuff.Id)
					if err == nil {
						//logger.Debug("yyyyyyyyyyyyyyyyyyyyy", buffCfg.ID, petBuff.Id)
						if petBuff.Id == buffCfg.ID && (constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageAttackBoost ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageDefenseBoost ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageMaxHealthIncrease ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffReflectDamage) {
							for i, e := range petBuff.Stats.Attrs {
								//logger.Infof("dddddddddddddddddddddddddddd%d------------%d", i, e)
								r.stats.setStats(i, e)
							}
						}
					}

				}
			}
		}
		return true
	})

	//for _, buff := range r.buffs {
	//	buffCfg, err := global.GetBuffCfg(buff.Id)
	//	if err == nil {
	//		// buff白名单，有的buff不计算战斗力，
	//		// 判断是不是在白名单里面
	//		// 获取宠物buff
	//		if len(petBuffs) > 0 {
	//			for _, petBuff := range petBuffs {
	//				// 获取buffid对应的大类
	//				buffCfgDetail, err := global.GetBuffCfg(petBuff.Id)
	//				if err == nil {
	//					//logger.Debug("yyyyyyyyyyyyyyyyyyyyy", buffCfg.ID, petBuff.Id)
	//					if petBuff.Id == buffCfg.ID && (constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageAttackBoost ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageDefenseBoost ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageMaxHealthIncrease ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffReflectDamage) {
	//						for i, e := range petBuff.Stats.Attrs {
	//							//logger.Infof("dddddddddddddddddddddddddddd%d------------%d", i, e)
	//							r.stats.setStats(i, e)
	//						}
	//					}
	//				}
	//
	//			}
	//		}
	//	}
	//
	//}
	//r.block.Unlock()
	//}

	// 羁绊
	for _, v := range r.activateBonds {
		if cfg := global.GetPartnerBondCfg(v); cfg != nil {
			for _, a := range cfg.Attribute {
				r.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}
	//五行石计算
	if r.ElementInfo != nil {
		r.stats.setStats(pb.Attribute(pb.Attribute_Main_EMPOWER), r.ElementInfo.elementMain)
		for _, v := range r.ElementInfo.element {
			if v.Eid == 0 {
				continue
			}
			if cfg := global.GetElementCfg(v.Eid); cfg != nil {
				for _, a := range cfg.Levelvalue {
					r.stats.setStats(pb.Attribute(a.Key), a.Value)
				}
				if v.Break != cfg.Breaklevel {
					if v.Eid-1 == 0 {
						continue
					}
					if cfg := global.GetElementCfg(v.Eid - 1); cfg != nil {
						for _, a := range cfg.Breakvalue {
							if a.Value > 0 {
								r.stats.setStats(pb.Attribute(a.Key), a.Value)
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
								r.stats.setStats(pb.Attribute(a.Key), a.Value)
							}
						}
					}
				}
			}
		}
	}

	//称号共鸣
	for _, v := range r.activateTitleBonds {
		if cfg := global.GetTitleActivationCfg(v); cfg != nil {
			for _, a := range cfg.Attribute {
				r.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}
	// 称号升星属性
	for _, v := range r.titles.Titles {
		if cfg := global.GetTitleCfg(v.CfgId); cfg != nil {
			if cfg.Slot == 1 {
				for _, g := range cfg.Goods {
					if v.Star == g.Id {
						for _, a := range g.Values {
							r.stats.setStats(pb.Attribute(a.Key), a.Value)
						}
						break
					}
				}
			}
		}
	}
	//称号装备属性
	if r.equipTitleId > 0 {
		if cfg := global.GetTitleCfg(r.equipTitleId); cfg != nil {
			for _, a := range cfg.Limit {
				r.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}

	// 头像和头像框属性
	r.avatar.Stats(r.Player)

	// 槽位精炼属性
	r.refine.Stats(r.Player, nil, pb.CharacterSlot_CHARACTER_SLOT_1)

	// 神兵属性
	r.artifactStatus(r.MapObject)

	// 参与计算战斗力的属性写在上面!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	// combatStats := Stats{Attrs: r.stats.copyStats()}
	// buff
	//if r.block.TryLock() {
	r.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
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
			r.stats.setStats(key, value)
			return true
		})
		//for i, e := range buff.BuffStats.Attrs {
		//	r.stats.setStats(i, e)
		//}
		return true
	})

	//for _, buff := range r.buffs {
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
	//		r.stats.setStats(i, e)
	//	}
	//}
	//	r.block.Unlock()
	//}
	// <0  减益 ; >0 增益
	r.stats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(r.GetDC())*(float64(r.stats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(r.GetMC())*(float64(r.stats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(r.GetAR())*(float64(r.stats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(r.GetMR())*(float64(r.stats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	r.stats.setStats(pb.Attribute_DODGE, int32(float64(r.stats.getStats(pb.Attribute_DODGE))*(float64(r.stats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_HIT, int32(float64(r.stats.getStats(pb.Attribute_HIT))*(float64(r.stats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(r.stats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(r.stats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_TOUGHNESS, int32(float64(r.stats.getStats(pb.Attribute_TOUGHNESS))*(float64(r.stats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	r.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(r.GetMaxHP())*float64(r.stats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))
	r.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(r.GetMaxHP())*float64(r.stats.getStats(pb.Attribute_HEALTH_MAX_PERCENTAGE))/BaseRate))

	// 攻速移速增加  基础*(1-加成)
	r.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(r.stats.getStats(pb.Attribute_ATTACK_SPEED))*(1-float64(r.stats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate)))
	r.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(r.stats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(r.stats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	r.stats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(r.GetDC()))))
	r.stats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(r.GetMC()))))
	r.stats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(r.GetAR()))))
	r.stats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(r.GetMR()))))

	r.stats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_DODGE)))))
	r.stats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_HIT)))))
	r.stats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	r.stats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_TOUGHNESS)))))
	r.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_ATTACK_SPEED)))))
	r.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(r.stats.getStats(pb.Attribute_MOVE_SPEED)))))
	r.stats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(r.GetMaxHP()))))

	if r.curState == PlayerStatusPlaying {
		if r.currentHP > r.GetMaxHP() {
			r.currentHP = r.GetMaxHP()
		}
		r.stats.changeStats(pb.Attribute_HEALTH, r.currentHP)
		// combatStats.changeStats(pb.Attribute_HEALTH, r.currentHP)
	} else {
		r.stats.changeStats(pb.Attribute_HEALTH, r.GetMaxHP())
		r.currentHP = r.GetMaxHP()
		r.displayHP = r.currentHP
		// combatStats.changeStats(pb.Attribute_HEALTH, r.GetMaxHP())
	}

	combat := int32(r.calcCombatRealV2(&Stats{Attrs: r.stats.copyStats()}, false))

	//r.logger.Debugf(" 机器人【%t】  [%s] 真实战力  %d  属性 %v ", r.IsRobot(), r.name, combat, &combatStats)

	packet := &pb.SyncUserStatsPush{Attributes: []*pb.Stats{}, Partners: []*pb.SyncUserStatsPush_PartnerStats{}}

	packet.Combat = combat
	// 刷新伙伴属性
	for _, v := range r.partnerSlots {
		if v == nil {
			continue
		}
		partner := r.getPartnerInfo(v.PartnerUid)
		if partner == nil {
			continue
		}
		partner.refreshStats(packet, r.Player)
	}

	packet.TotalCombat += packet.Combat
	for _, v := range packet.Partners {
		packet.TotalCombat += v.Combat
	}

	r.totalCombat = packet.TotalCombat

	// r.logger.Debugf("属性计算 当前生命:%d  最大生命:%d  物攻: %d 魔攻: %d", r.currentHP, r.GetMaxHP(), r.GetDC(), r.GetMC())
}

func (r *RobotInfo) defaultBuff() {
	// 宠物加buff
	// 宠物buff数组，每个元素包含ID和对应的属性

	// 初始化buff
	//for _, v := range r.buffs {
	//	r.delBuff(v.Id)
	//}
	ids := r.buffMgr.getBuffIds()
	for _, v := range ids {
		r.delBuff(v)
	}
	// 获取宠物天赋Buff列表
	petBuffs := r.getPetBuffs()
	for k, v := range global.BuffCfg {
		if v.Mark == constant.BuffMarkPassive {
			r.AddBuffNoRefresh(k, r, 24*60*60*1000, nil, 0, 0)
		}
		// 如果在这个数组里面，就加buff
		if len(petBuffs) > 0 {
			// 如果这个buff在这个数组里面，就加buff
			for _, vv := range petBuffs {
				if vv.Id == k {
					r.AddBuffNoRefresh(k, r, 24*60*60*1000, vv.Stats, vv.Rate, 0)
				}
			}
		}

	}
	r.refreshStats()
}

// 增加伙伴
func (r *RobotInfo) addPartner(cfgId int32) error {
	if cfgId == 0 {
		return nil
	}

	// 重复获得伙伴转换成碎片
	for _, v := range r.partners {
		if v.cfg.Cid == cfgId {
			return nil
		}
	}

	par := CreatePartner(cfgId)
	par.owner = r
	r.partners[par.Uid()] = par

	// 解锁技能
	par.learnSkill()

	r.checkPartnerBonds()
	r.refreshStats()

	return nil
}

// 伙伴上阵
func (r *RobotInfo) partnerCombat(slotId int32, partnerUid int64) error {
	partnerSlot := r.getPartnerSlotInfo(slotId)
	if partnerSlot == nil {
		return errutil.ErrPartnerInvalid
	}

	partner := r.getPartnerInfo(partnerUid)
	if partner == nil {
		return errutil.ErrPartnerInvalid
	}

	if partnerSlot.PartnerUid > 0 {
		return errutil.ErrPartnerInvalid
	}

	for _, v := range r.partnerSlots {
		if v.PartnerUid == partnerUid {
			return errutil.ErrPartnerInvalid
		}

		info := r.getPartnerInfo(v.PartnerUid)
		if info != nil {
			if info.cfg.Cid == partner.cfg.Cid {
				return errutil.ErrPartnerInvalid
			}
		}
	}

	partnerSlot.PartnerUid = partnerUid

	r.refreshStats()
	return nil
}
