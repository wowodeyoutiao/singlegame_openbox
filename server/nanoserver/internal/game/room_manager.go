package game

import (
	"nanoserver/db"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/async"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/timex"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"nanoserver/protocol/protobuf/pb"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
)

type RoomManager struct {
	component.Base
	//key:roomCfgID
	rooms map[int32]*DeskManager
}

var roomManager = NewRoomManager()

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}

func (roommgr *RoomManager) AfterInit() {
	roomManager.rooms = make(map[int32]*DeskManager, len(global.RoomCfg))
	for i := range global.RoomCfg {
		room := NewDeskManager(global.RoomCfg[i])
		roomManager.rooms[global.RoomCfg[i].ID] = room
	}

	//初始化游戏房间
	for _, v := range roomManager.rooms {
		for i := 0; i < v.opt.TotalTables; i++ {
			nd := worldMgr.GetWorld(v.opt.GameType).NewDesk(snowflake.GenID(), int32(v.ID()))
			if nd != nil {
				v.desks[nd.deskID] = nd
				nd.createdAt = time.Now().Unix()
			}
		}
	}

	// 战场房间分给不同战区
	zonelist := global.ServerListMgrInstance.GetLocalServerZoneIds()
	markRoomids := []int32{}
	for _, v := range roomManager.rooms {
		if v.opt.Roomtype == constant.PVEWarMap {
			for _, mrv := range markRoomids {
				if v.ID() == mrv {
					break
				}
			}
			markRoomids = append(markRoomids, v.ID())

			for zi := 0; zi < len(zonelist); zi++ {
				for dk, dv := range v.desks {
					if dv.zoneID == 0 {
						v.desks[dk].zoneID = zonelist[zi]
						break
					}
				}
			}
		}
	}

	// 延迟加载机器人，等数据库初始化完毕
	time.AfterFunc(time.Second, func() { robotMgr.loadRobot() })

	session.Lifetime.OnClosed(func(s *session.Session) {
		// Fixed: 玩家WIFI切换到4G网络不断开, 重连时，将UID设置为illegalSessionUid
		if s.UID() > 0 {
			if err := roommgr.onPlayerDisconnect(s); err != nil {
				logger.Errorf("玩家退出: userId=%d, Error=%s", s.UID(), err.Error())
			}
		}
	})

	// 每5分钟清空一次已摧毁的房间信息
	scheduler.NewTimer(300*time.Second, func() {
		for _, v := range roommgr.rooms {
			destroyDesk := map[int64]*Desk{}
			for no, d := range v.desks {
				// 清除创建超过24小时的房间
				if d.status() == constant.DeskStatusDestory {
					destroyDesk[no] = d
				}
			}
			for _, d := range destroyDesk {
				d.destroy()
			}

			//v.dumpDeskInfo()

			// 统计结果异步写入数据库
			sCount := defaultManager.sessionCount()
			dCount := len(v.desks)
			async.Run(func() {
				db.InsertOnline(sCount, dCount)
			})
		}
	})

	// 每个房间最多进入3个机器人
	scheduler.NewTimer(time.Second*10, func() {
		for _, v := range roommgr.rooms {
			if v.opt.Roomtype != constant.PVEWarMap {
				continue
			}
			// for _, d := range v.desks {
			// 	if len(d.players) >= 1 && len(d.robots) < 1 {
			// 		r := robotMgr.getFreeRobot(0)
			// 		if r != nil {
			// 			r.refreshStats()
			// 			d.robotJoin(r, 0)
			// 			d.syncDeskStatus()
			// 			break
			// 		}
			// 	}
			// }
		}
	})

	// itemRecoderMgr.load()
}

// 断线
func (roommgr *RoomManager) onPlayerDisconnect(s *session.Session) error {
	p, err := playerWithSession(s)
	if err != nil {
		logger.Warningf("玩家网络断开!!! session 已释放 uid: %d ", s.UID())
		return err
	}
	p.disconnect()
	p.InstanceData.PushingWall.TmpUnixTime = 0

	roommgr.UpOfflineShushuDate(p)

	if err := thinkingdatax.TrackerIns.Set(strconv.FormatInt(p.Uid(), 10), p.accountName, tracking.M{
		"account_id": p.accountName,
	}); err != nil {
		logger.Warnf("ThinkingData UserSet failed: accountID=%s, error=%v", p.accountName, err)
	}

	return nil
}

// 登出上报
func (roommgr *RoomManager) UpOfflineShushuDate(p *Player) {
	//创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	//upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "role_logout"
	shushuDate["logout_time"] = time.Now().Unix()
	shushuDate["#duration"] = p.wellenOnlineTime

	zun := false
	//TmpPetID := int64(0)
	//if len(p.PetMgr.pet) > 0 {
	//	TmpPetID = int64(p.PetMgr.pet[0].PetCfgId)
	//}
	tmpPetID := p.PetMgr.GetBattlePet().PetCfgId
	// 查询宠物信息
	if petCfg := global.GetPetCfg(tmpPetID); petCfg != nil {
		if petCfg.Mark == 1 {
			zun = true
		}
	}
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_ingot"] = int32(p.tokens[constant.ItemTypeCoupon])
	shushuUserSetDate["current_ingot_cost"] = int32(p.tokenDeletes[constant.ItemTypeCoupon])
	shushuUserSetDate["current_coin"] = int32(p.tokens[constant.ItemTypeCoin])
	shushuUserSetDate["current_coin_cost"] = int32(p.tokenDeletes[constant.ItemTypeCoin])
	shushuUserSetDate["current_battle_soul"] = int32(p.tokens[constant.ItemTypeSoul])
	shushuUserSetDate["current_battle_soul_cost"] = int32(p.tokenDeletes[constant.ItemTypeSoul])
	shushuUserSetDate["current_trial_coin"] = int32(p.tokens[constant.ItemTypeTrialVoucher])
	shushuUserSetDate["current_trial_coin_cost"] = int32(p.tokenDeletes[constant.ItemTypeTrialVoucher])
	shushuUserSetDate["current_boxes_left"] = int32(p.tokens[constant.ItemTypeTreasureNum])
	shushuUserSetDate["current_boxes"] = int32(p.tokenDeletes[constant.ItemTypeTreasureNum])
	shushuUserSetDate["current_guild_coin"] = int32(p.tokens[constant.ItemTypeGuildContribution])
	shushuUserSetDate["current_guild_coin_cost"] = int32(p.tokenDeletes[constant.ItemTypeGuildContribution])
	shushuUserSetDate["last_login_time"] = p.lastLoginAt.Unix()
	shushuUserSetDate["last_login_time_v2"] = time.Unix(p.lastLoginAt.Unix(), 0).Format(timex.FormatDateTimeMillis)
	shushuUserSetDate["current_pet_zun"] = zun

	TmpPartnerLevel := int32(0)
	for _, v := range p.partnerSlots {
		if v.PartnerUid > 0 {
			if par := p.getPartnerInfo(v.PartnerUid); par != nil {
				TmpPartnerLevel += par.Level()
			}
		}
	}
	TmpMonthCard := false
	TmpLifetimeCard := false
	for _, v := range p.privilege {
		cfg, ok := global.PrivilegeCfg[int32(v.Pid)]
		if !ok {
			continue
		}
		if cfg.Duration < 31 {
			TmpMonthCard = true
		} else {
			TmpLifetimeCard = true
		}
	}
	ad_total := p.getUpRecord("ad_click")
	shushuUserSetDate["all_partner_level"] = TmpPartnerLevel
	shushuUserSetDate["is_month_card"] = TmpMonthCard
	shushuUserSetDate["is_lifetime_card"] = TmpLifetimeCard
	shushuUserSetDate["current_time"] = p.totalOnlineTime
	shushuUserSetDate["current_power"] = p.totalCombat
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushuUserSetDate["ad_total"] = ad_total
	shushuUserSetDate["ad_total_today"] = atomic.LoadInt32(&p.adTotalToday)
	shushumain.ReportUserSetData(shushuUserSetDate)
	shushumain.ReportData(shushuDate)
}

// 根据游戏房间号获取游戏房间
func (roommgr *RoomManager) getRoom(roomId int32) *DeskManager {
	room, ok := roomManager.rooms[roomId]
	if !ok {
		return nil
	}
	return room
}

func (roommgr *RoomManager) getRoomByType(roomType int32, level int32) *DeskManager {

	for _, v := range roomManager.rooms {
		if v.opt.Roomtype == roomType && v.opt.MinEntryLimit <= int(level) && v.opt.MaxEntryLimit >= int(level) {
			return v
		}
	}
	return nil
}

// 加入
func (roommgr *RoomManager) Join(s *session.Session, data *pb.JoinDeskRequest) error {

	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.JoinDeskResponse{}
	if p.desk != nil {
		logger.Warnf("进入地图失败, 已经有房间 uid=%d  name=%s 进入房间=%d", p.uid, p.name, data.RoomId)
		res.Code = -1
		res.Error = errutil.ErrAlreadyInRoom.Error()
		return s.Response(res)
	}
	room, ok := roomManager.rooms[data.RoomId]
	if !ok {
		logger.Debugf("进入地图失败, 房间类型不存在 uid=%d  name=%s 进入房间=%d", p.uid, p.name, data.RoomId)
		res.Code = -2
		res.Error = errutil.ErrDeskNotFound.Error()
		return s.Response(res)
	}

	if p.level < int32(room.opt.MinEntryLimit) || p.level > int32(room.opt.MaxEntryLimit) {
		res.Code = -3
		res.Error = errutil.ErrLevel.Error()
		return s.Response(res)
	}

	switch room.opt.Roomtype {
	case constant.SoloTowerMap:

		if endlesstrialsCfg := global.GetEndlesstrialsCfg(p.SoloTowerInfo.Layers + 1); endlesstrialsCfg != nil {
			if endlesstrialsCfg.LevelLimit > p.level {
				res.Code = -4
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		} else {
			res.Code = -5
			res.Error = errutil.ErrConfigNotFound.Error()
			return s.Response(res)
		}
	case constant.PVEPeacefulMap:

		if !p.isModuleOpen(constant.ModuleGuardIsland) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}

		if p.afkLeftTime() <= 0 {
			res.Code = -9
			res.Error = errutil.NoTimes
			return s.Response(res)
		}
	case constant.WorldBossMap:
		if !p.isModuleOpen(constant.ModuleWorldBoss) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}

		if p.dailyWorldBossTimes >= int32(global.GlobalCfg.WorldBoss.EnterTime.Num) {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !worldBossMgr.checkWorldBossOpen() {
			res.Code = -11
			res.Error = errutil.EventIsPending
			return s.Response(res)
		}
		if worldBossMgr.checkWorldBossDead(p.serverID) {
			res.Code = -12
			res.Error = errutil.ErrActivityExpire.Error()
			return s.Response(res)
		}
	case constant.StarUpStoneFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleStarsFB) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}
	case constant.ExperienceFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleExpFB) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}
	case constant.BookPageFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleBookFB) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}
	case constant.SoulStoneFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleEquipSoul) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}

	case constant.EnhanceStoneFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleEnhanceFB) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}
	case constant.GoldFBMap:
		if p.dailyMaps[room.opt.Roomtype] <= 0 {
			res.Code = -10
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		if !p.isModuleOpen(constant.ModuleCoinFB) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}

	case constant.GuidanceFBMap:
		if p.bCreate {
			res.Code = -11
			res.Error = errutil.ErrActivityExpire.Error()
			return s.Response(res)
		}
	case constant.BattleMap, constant.TrialTowerMap:
		if data.BattleId != p.getNextBattleId(room.opt.Roomtype) {
			res.Code = -12
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
		if cfg := global.GetBattleMapCfg(data.BattleId); cfg != nil {
			if cfg.LevelLimit > p.level || cfg.PowerLimit > p.totalCombat {
				res.Code = -13
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		}

	case constant.BossMap:
		if data.BattleId > p.getNextBattleId(room.opt.Roomtype) {
			res.Code = -14
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
		if cfg := global.GetBattleMapCfg(data.BattleId); cfg != nil {
			if cfg.LevelLimit > p.level || cfg.PowerLimit > p.totalCombat {
				res.Code = -15
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		}
		if p.BossMap.ChallengeTime <= 0 {
			res.Code = -16
			res.Error = errutil.NoTimes
			return s.Response(res)

		}
	case constant.RelicMap:
		if data.BattleId > p.getNextBattleId(room.opt.Roomtype) {
			res.Code = -17
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
		if cfg := global.GetBattleMapCfg(data.BattleId); cfg != nil {
			if cfg.LevelLimit > p.level || cfg.PowerLimit > p.totalCombat {
				res.Code = -18
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		}
		if p.RelicMap.ChallengeTime <= 0 {
			res.Code = -19
			res.Error = errutil.NoTimes
			return s.Response(res)
		}
	case constant.FightEndMap:
		if !p.isModuleOpen(constant.ModuleOneBattle) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}
		//检查副本开启时间
		if data.BattleId > p.getNextBattleId(room.opt.Roomtype) {
			res.Code = -17
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
		if cfg := global.GetBattleMapCfg(data.BattleId); cfg != nil {
			if cfg.LevelLimit > p.level || cfg.PowerLimit > p.totalCombat {
				res.Code = -18
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		}
		t := time.Now()
		startTime_, endTime_, err := global.GetFightEndTime()
		if err != nil {
			logger.Errorf("一战到底每日开启时间配置错误 %s", err.Error())
			res.Code = -19
			res.Error = errutil.EventIsPending
			return s.Response(res)
		}

		if t.Before(startTime_) || t.After(endTime_) {
			res.Code = -19
			res.Error = errutil.EventIsPending
			return s.Response(res)
		}
	case constant.GuildMap:
		t := time.Now()
		startTime_, endTime_, err := global.GetGuildCopyTime()
		if err != nil {
			logger.Errorf("行会BOSS每日开启时间配置错误 %s", err.Error())
			res.Code = -19
			res.Error = errutil.EventIsPending
			return s.Response(res)
		}

		if t.Before(startTime_) || t.After(endTime_) {
			res.Code = -19
			res.Error = errutil.EventIsPending
			return s.Response(res)
		}

		if cfg := global.GetBattleMapCfg(data.BattleId); cfg != nil {
			if cfg.LevelLimit > p.level || cfg.PowerLimit > p.totalCombat {
				res.Code = -18
				res.Error = errutil.ErrLevel.Error()
				return s.Response(res)
			}
		}
	case constant.PVEWarMap:
		if !p.isModuleOpen(constant.ModuleBattlefield) {
			res.Code = -21
			res.Error = errutil.ErrModuleNotOpen
			return s.Response(res)
		}

		if !global.InWarTime(data.RoomId) {
			res.Code = -21
			res.Error = errutil.EventIsPending

			for _, v := range global.GlobalCfg.War.OpenTime {
				if v.RoomId == data.RoomId {
					logger.Debugf("战场不在开启时间内 时间配置 %v ", v)
				}
			}
			return s.Response(res)
		}

		for _, v := range room.desks {
			if v.zoneID == p.getZoneId() {
				if len(v.players) >= global.GlobalCfg.War.EnterLimit.MaxNum {
					res.Code = -22
					res.Error = errutil.RoomPlayerMax
					return s.Response(res)
				}
			}
		}

	default:
		res.Code = -100
		res.Error = errutil.ErrDeskNotFound.Error()
		return s.Response(res)
	}

	if err := room.Join(p, data); err != nil {
		res.Code = -6
		res.Error = err.Error()
		return s.Response(res)
	}
	return nil
}

// 退出
func (roommgr *RoomManager) Exit(s *session.Session, _ []byte) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	d := p.desk
	if d == nil || d.isDestroy() {
		p.logger.Debugln("退出房间失败 玩家不在房间内")
		return s.Response(&pb.ExitResponse{})
	}

	p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: false}

	select {
	case p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: false}:
	default:
		p.logger.Warningln(" 移动失败 chMove 通道已满")
	}

	s.Response(&pb.ExitResponse{})

	return nil
}

// 移动
func (roommgr *RoomManager) Move(s *session.Session, data *pb.MoveRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	//p.logger.Debugf("[ %s ] 移动请求！！！！！！！！！！！！！！！ Uid: %d ", p.name, p.Uid())
	res := &pb.MoveResponse{Location: data.Location}

	if p.desk == nil || p.desk.isDestroy() {
		p.logger.Debugf("XXXXXX移动失败,不在房间内 Uid: %d ", p.Uid())
		res.Code = -1
		//res.Error = errutil.ErrDeskNotFound.Error()
		return s.Response(res)
	}

	// if !p.CanMove() {
	// 	p.logger.Debugf("XXXXXXXXX移动失败,移动频率过快 Uid: %d  下次移动间隔 %d", p.Uid(), p.moveTime.UnixMilli()-time.Now().UnixMilli())
	// 	res.Code = -2
	// 	res.Error = errutil.ErrDeskNotFound.Error()
	// 	return s.Response(res)
	// }

	//检查越界
	if data.Location.X < 0 || data.Location.X > float64(p.desk.mapCfg.High) ||
		data.Location.Y < 0 || data.Location.Y > float64(p.desk.mapCfg.Wide) {

		p.logger.Debugf("XXXXXXXXXXXXX移动失败,地图边缘 Uid: %d ", p.Uid())
		res.Code = -3
		res.Error = errutil.ErrInvalidPoint.Error()
		return s.Response(res)
	}

	select {
	case p.desk.chMove <- &MoveRequest{Uid: p.uid, Location: &pb.Point{X: data.Location.X, Y: data.Location.Y}}:
	default:
		p.logger.Warningln(" 移动失败 chMove 通道已满")
		res.Code = -10
		return s.Response(res)
	}

	return s.Response(res)
}

// 攻击
func (roommgr *RoomManager) Attack(s *session.Session, data *pb.AttackRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	res := &pb.AttackResponse{
		TargetUids: data.TargetUids,
		SkillId:    data.SkillId,
		Location:   data.Location,
		Direction:  data.Direction}

	//检查是否在游戏中
	if p.desk == nil || p.desk.isDestroy() || p.dead {
		p.logger.Debugf("攻击失败 玩家不在房间内 Uid: %d", p.Uid())
		res.Code = -1
		return s.Response(res)
	}

	// skill, cast_uid := p.getSkillex(skillId)
	// if skill == nil {
	// 	p.logger.Debugf("攻击失败 使用了未学习的技能  skillID: %d ", skillId)
	// 	res.Code = -2
	// 	return s.Response(res)
	// }

	skillCfg, err := global.GetSkillCfg(data.SkillId)
	if err != nil {
		p.logger.Debugf("攻击失败 无效的技能ID  skillID: %d ", data.SkillId)
		res.Code = -3
		res.Error = errutil.ErrConfigNotFound.Error()
		return s.Response(res)
	}

	// if cast_uid != p.uid {
	// 	if partner := p.getPartnerInfo(cast_uid); partner != nil {
	// 		if partner.dead {
	// 			res.Code = -4
	// 			return s.Response(res)
	// 		}
	// 		/*
	// 			if skillCfg.SType != 3 && time.Now().Before(partner.attackTime) {
	// 				p.logger.Debugf("伙伴 %s 攻击频率过快 当前攻速 %d , 离下次攻击时间 %d 毫秒", partner.name, partner.stat(pb.Attribute_ATTACK_SPEED), partner.attackTime.UnixMilli()-time.Now().UnixMilli())

	// 				//err_msg := fmt.Sprintf("debug 伙伴[%s]攻击频率过快, 距离下次攻击时间[%d]毫秒, 技能[%d]", partner.name, partner.attackTime.UnixMilli()-time.Now().UnixMilli(), skill.CfgId)
	// 				res.Error = errutil.ErrCannotAttack.Error()
	// 				res.Code = -5
	// 				return s.Response(res)
	// 			}
	// 		*/
	// 	} else {
	// 		res.Code = -6
	// 		return s.Response(res)
	// 	}
	// } else {
	// 	if skillCfg.SType != 3 && time.Now().Before(p.attackTime) {

	// 		p.logger.Debugf("主角 %s 攻击频率过快 当前攻速 %d , 离下次攻击时间 %d 毫秒 下次攻击 %d ", p.name, p.stat(pb.Attribute_ATTACK_SPEED), p.attackTime.UnixMilli()-time.Now().UnixMilli(), p.attackTime.UnixMilli())

	// 		//err_msg := fmt.Sprintf("debug 主角[%s]攻击频率过快, 距离下次攻击时间[%d]毫秒, 技能[%d] -- %d", p.name, p.attackTime.UnixMilli()-time.Now().UnixMilli(), skill.CfgId, p.attackTime.UnixMilli())
	// 		//res.Error = errutil.ErrCannotAttack.Error()
	// 		res.Code = -7
	// 		return s.Response(res)
	// 	}
	// }

	// 技能CD
	// if time.Now().Before(skill.Cooldown) {
	// 	p.logger.Debugf("攻击失败 技能冷却时间未到  skillID: %d name %s  剩余 %d 毫秒", skillId, skillCfg.Name, skill.Cooldown.UnixMilli()-time.Now().UnixMilli())
	// 	res.Error = errutil.ErrSkillCoolDown.Error()
	// 	res.Code = -8
	// 	return s.Response(res)
	// }

	if skillCfg.SkillID == constant.SkillTypeLastGasp || skillCfg.SkillID == constant.SkillTypeReborn {
		p.logger.Debugf(" %s 技能释放失败， 不能主动释放 %s ", p.name, skillCfg.Name)

		res.Error = errutil.IllegalSkillType
		res.Code = -9
		return s.Response(res)
	}

	select {
	case p.desk.chAttack <- &AttackRequest{Uid: p.uid, Req: data}:
	default:
		p.logger.Warningln(" 攻击失败 chAttack 通道已满")
		res.Error = errutil.ErrCannotAttack.Error()
		res.Code = -10
		return s.Response(res)
	}

	return s.Response(res)
}

// 续费
func (roommgr *RoomManager) InstanceTimeAdd(s *session.Session, data *pb.InstanceTimeAddRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.InstanceTimeAddResponse{Code: -1}

	if data.AutoBuy && data.EnergyNum != 0 {
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}
	//检查是否在游戏中
	if p.desk == nil || p.desk.isDestroy() {
		p.logger.Debugf("增加时间失败,不在房间内 Uid: %d", p.Uid())
		res.Error = errutil.ErrDeskNotFound.Error()
		return s.Response(res)
	}

	if data.EnergyNum != 0 {
		var needEnergy, addTime int
		if p.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
			needEnergy = global.GlobalCfg.BaseValue.Energy.ConsumeenergyPeace
			addTime = global.GlobalCfg.BaseValue.Energy.BuytimePeace

		} else if p.desk.roomCfg.Roomtype == constant.PVEWarMap {
			needEnergy = global.GlobalCfg.BaseValue.Energy.ConsumeenergyWar
			addTime = global.GlobalCfg.BaseValue.Energy.BuytimeWar
		} else {
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}

		if p.getItemCountInBag(constant.ItemTypeEnergy) < int64(needEnergy) {
			p.logger.Debugf("增加副本时间失败,体力不足 Uid: %d", p.Uid())
			res.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(res)
		}
		p.takeSource(constant.ItemTypeEnergy, int64(needEnergy))
		p.expirationTime = p.expirationTime.Add(time.Duration(addTime) * time.Second)

	} else {
		p.autoBuyInstanceTime = data.AutoBuy
	}

	res.Code = 0
	res.InstanceTimeLeft = int32(p.expirationTime.Unix() - time.Now().Unix())
	res.AutoBuy = p.autoBuyInstanceTime
	return s.Response(res)
}

// 复活
func (roommgr *RoomManager) Revive(s *session.Session, data *pb.ReviveRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.ReviveResponse{Location: p.location, Direction: p.direction}

	//检查是否在游戏中
	if p.desk == nil || p.desk.isDestroy() {
		p.logger.Debugf("复活失败 不在房间内 Uid: %d", p.Uid())
		res.Code = -1
		res.Error = errutil.ErrDeskNotFound.Error()
		return s.Response(res)
	}

	if !p.dead {
		res.Code = -2
		return s.Response(res)
	}

	if p.desk.roomCfg.Roomtype == constant.PVPOneOnOneMap || p.desk.roomCfg.Roomtype == constant.PVPMultiOnMultiMap ||
		p.desk.roomCfg.Roomtype == constant.SoloTowerMap || p.desk.roomCfg.Roomtype == constant.WorldBossMap {
		res.Error = errutil.ErrWrongType.Error()
		res.Code = -3
		return s.Response(res)
	}

	costId := int32(0)
	costNum := int64(0)
	if p.desk.roomCfg.Roomtype == constant.PVEWarMap {
		if data.ReviveType == 1 {
			var bFind bool
			for _, v := range global.GlobalCfg.War.Revive {
				if v.ID-1 == p.reviveTimes {
					costId = int32(global.GlobalCfg.War.Revive[p.reviveTimes].NeedItem)
					costNum = int64(global.GlobalCfg.War.Revive[p.reviveTimes].Num)
					bFind = true
					break
				}
			}
			if !bFind {
				costId = int32(global.GlobalCfg.War.Revive[len(global.GlobalCfg.War.Revive)-1].NeedItem)
				costNum = int64(global.GlobalCfg.War.Revive[len(global.GlobalCfg.War.Revive)-1].Num)
			}
		} else {
			p.reviveTimes = 0
		}

	} else if p.desk.roomCfg.Roomtype == constant.EnhanceStoneFBMap || p.desk.roomCfg.Roomtype == constant.BookPageFBMap {
		if p.reviveTimes >= len(global.GlobalCfg.Revive) {
			res.Code = -4
			return s.Response(res)
		}
		costId = int32(global.GlobalCfg.Revive[p.reviveTimes].NeedMoneyID)
		costNum = int64(global.GlobalCfg.Revive[p.reviveTimes].Num)
	}

	// 检查道具
	if costId > 0 && costNum > 0 {
		if !p.checkTakeSource(costId, costNum) {
			res.Code = -5
			res.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(res)
		}
		p.takeSource(costId, costNum)
	}

	select {
	case p.desk.chPlayerRelive <- &PlayerRelive{p.uid, data.ReviveType}:
	default:
		p.logger.Warningln("复活失败 chPlayerRelive 通道已满")
	}

	res.CurHp = p.currentHP
	return s.Response(res)
}

// 爬塔
func (roommgr *RoomManager) EndlesstrialsReward(s *session.Session, req *pb.EndlesstrialsRewardRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	res := &pb.EndlesstrialsRewardResponse{Reward: []*pb.ItemInfo{}}
	res.Layers = p.Layers

	cfg := global.GetEndlesstrialsCfg(p.SoloTowerInfo.Layers)
	if cfg == nil {
		res.Code = -1
		res.Error = errutil.ErrConfigNotFound.Error()
		return s.Response(res)
	}

	offsetMinite := (time.Now().Unix() - p.SoloTowerInfo.RewardTime.Unix())
	if req.Type == 1 {
		for _, v := range cfg.OfflineReward {
			multiple := offsetMinite / int64(v.Gap)
			num := int64(v.Num) * multiple
			if num > int64(v.Max) {
				num = int64(v.Max)
			}
			res.Reward = append(res.Reward, &pb.ItemInfo{CfgId: v.ID, Num: num})
		}
		if p.SoloTowerInfo.Layers > 0 && time.Now().After(p.SoloTowerInfo.RewardTime.Add(30*time.Minute)) {
			res.Ready = true
		}

	} else if req.Type == 2 {
		if p.SoloTowerInfo.Layers <= 0 || time.Now().Before(p.SoloTowerInfo.RewardTime.Add(30*time.Minute)) {
			res.Code = -2
			return s.Response(res)
		}

		if err := p.soloTowerOfflineReward(); err != nil {
			res.Code = -3
			res.Error = err.Error()
			return s.Response(res)
		}
		res.Ready = false
	}
	return s.Response(res)
}

// 双倍
func (roommgr *RoomManager) DoubleReward(s *session.Session, req *pb.DoubleRewardRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	c, e := p.dailyMapDoubleReward(req.Double)
	res := &pb.DoubleRewardResponse{Code: c, Error: e}
	return s.Response(res)
}

// 扫荡
func (roommgr *RoomManager) CopySweep(s *session.Session, req *pb.CopySweepRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	if cfg := global.GetRoomCfg(req.RoomId); cfg != nil {
		if cfg.Roomtype == constant.ExperienceFBMap || cfg.Roomtype == constant.GoldFBMap || cfg.Roomtype == constant.EnhanceStoneFBMap ||
			cfg.Roomtype == constant.BookPageFBMap || cfg.Roomtype == constant.StarUpStoneFBMap || cfg.Roomtype == constant.SoulStoneFBMap {
			c, e := p.dailyMapCopySweep(req.RoomId)
			if cfg.Roomtype == constant.ExperienceFBMap {
				p.tasks.dispatch(pb.TargetType_TOTAL_EXPERIENCE_FB_CLEARED, 0, 1)
			}
			if cfg.Roomtype == constant.GoldFBMap {
				p.tasks.dispatch(pb.TargetType_TOTAL_GOLD_FB_CLEARED, 0, 1)
			}
			if cfg.Roomtype == constant.SoulStoneFBMap {
				p.tasks.dispatch(pb.TargetType_SOUL_STONE_INSTANCE, 0, 1)
			}
			if cfg.Roomtype == constant.StarUpStoneFBMap {
				p.tasks.dispatch(pb.TargetType_STAR_UP_STONE_INSTANCE, 0, 1)
			}
			if cfg.Roomtype == constant.EnhanceStoneFBMap {
				p.tasks.dispatch(pb.TargetType_ENHANCE_STONE_INSTANCE, 0, 1)
			}
			if cfg.Roomtype == constant.BookPageFBMap {
				p.tasks.dispatch(pb.TargetType_BOOK_PAGE_INSTANCE, 0, 1)
			}
			if req.SweepTimes <= 0 {
				req.SweepTimes = 1
			}
			p.addUpRecordByAmount(constant.InstanceNum, int64(req.SweepTimes)) // 增加统计
			res := &pb.CopySweepResponse{Code: c, Error: e}
			return s.Response(res)
		} else if cfg.Roomtype == constant.BossMap || cfg.Roomtype == constant.RelicMap {
			if req.SweepTimes <= 0 {
				res := &pb.CopySweepResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}
				return s.Response(res)
			}
			c, e := p.battleMapCopySweep(req.RoomId, req.BattleId, req.SweepTimes)
			if cfg.Roomtype == constant.BossMap {
				p.tasks.dispatch(pb.TargetType_TOTAL_BOSS_CHALLENGE, 0, 1)
			}
			if cfg.Roomtype == constant.RelicMap {
				p.tasks.dispatch(pb.TargetType_TOTAL_RELIC_CHALLENGE, 0, 1)
			}
			// for i := 0; i < int(req.SweepTimes); i++ {
			p.addUpRecordByAmount(constant.InstanceNum, int64(req.SweepTimes)) // 增加统计
			// }
			res := &pb.CopySweepResponse{Code: c, Error: e}
			return s.Response(res)
		} else if cfg.Roomtype == constant.PVEPeacefulMap {
			c, e := p.guardIslandSweep(req.RoomId)
			res := &pb.CopySweepResponse{Code: c, Error: e}
			p.tasks.dispatch(pb.TargetType_GUARD_ISLAND_INSTANCE, 0, 1)
			if req.SweepTimes <= 0 {
				req.SweepTimes = 1
			}
			p.addUpRecordByAmount(constant.InstanceNum, int64(req.SweepTimes)) // 增加统计
			return s.Response(res)
		}
	}
	return errutil.ErrConfigNotFound
}

// 挖矿
func (roommgr *RoomManager) Mining(s *session.Session, req *pb.MiningRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.MiningResponse{}
	if p.desk == nil || p.desk.isDestroy() || p.dead {
		p.logger.Debugf("采集失败,不在房间内 Uid: %d ", p.uid)
		res.Code = -1
		return s.Response(res)
	}

	if time.Now().Before(p.attackTime) {
		p.logger.Debugf("采集失败,频率过快 Uid: %d ", p.uid)

		res.Code = -2
		return s.Response(res)
	}

	if p.bagLeftCount(pb.ItemType_EQUIPMENT) <= 0 {
		p.logger.Debugf("采集失败,包裹已满 Uid: %d ", p.uid)

		res.Code = -3
		res.Error = errutil.ErrBagFull.Error()
		return s.Response(res)
	}

	select {
	case p.desk.chMining <- &MiningRequest{Uid: p.uid, TargetUid: req.Uid}:
	default:
		logger.Warningln("挖矿失败 chMining 已满")
	}

	return s.Response(res)
}

// 每日
func (roommgr *RoomManager) DailyMapInfo(s *session.Session, req *pb.DailyMapInfoRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	packet := &pb.DailyMapInfoResponse{Info: []*pb.ChallengeInfo{}}

	for i := constant.ExperienceFBMap; i <= constant.SoulStoneFBMap; i++ {
		packet.Info = append(packet.Info, &pb.ChallengeInfo{MType: pb.MapType(i), Times: p.dailyMapTimes(i)})
	}
	return s.Response(packet)
}

// 挂机
func (roommgr *RoomManager) AutoPlay(s *session.Session, req *pb.AutoPlayRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.AutoPlayResponse{}
	if p.vip < global.GlobalCfg.BaseValue.Autobattle.Viplevel {
		res.Code = -1
		res.Error = errutil.VIPNotEnough
	}
	return s.Response(res)
}

func (roommgr *RoomManager) BattleReward(s *session.Session, req *pb.BattleRewardRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_BattleReward)}
	if err = p.getBattleReward(req.MType, req.ChapterId); err != nil {
		res.Code = -1
		res.Error = err.Error()
	}
	return s.Response(res)
}

// 主界面战斗
func (roommgr *RoomManager) MainBattle(s *session.Session, req *pb.MainBattleRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	p.logger.Debugf("主界面战斗,收到通关请求 %s MainBattleId: %d  已通关关卡ID %d", p.name, req.MainBattleId, p.mainBattleId)

	res := &pb.MainBattleResponse{MainBattleId: req.MainBattleId, Items: []*pb.ItemInfo{}}
	// 判断之前开的宝箱有没有数据
	if len(p.box.NowOpenRes) > 0 {
		res.Code = -2
		res.Error = errutil.ErrOpenBoxEnough.Error()
		return s.Response(res)
	}
	cfg := global.GetMainBattleCfg(req.MainBattleId)

	if cfg == nil {
		res.Code = -1
		res.Error = errutil.ErrConfigNotFound.Error()
		return s.Response(res)
	}

	if req.MainBattleId > p.mainBattleId+1 {
		res.Code = -2
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}

	if p.energy < int64(cfg.Cost) {
		res.Code = -3
		res.Error = errutil.EnergyNotEnough
		return s.Response(res)
	}

	if p.bagItemLeftCount() < 1 {
		res.Code = -3
		res.Error = errutil.ErrBagFull.Error()
		return s.Response(res)
	}

	Reward := []global.IdNum{}
	if req.MainBattleId > p.mainBattleId {
		for _, v := range cfg.FirstReward {
			Reward = append(Reward, v)
		}
	}

	for _, v := range cfg.Drop {
		Reward = append(Reward, v)
	}

	if err := p.checkAddSourceBatch(Reward); err != nil {
		res.Code = -4
		res.Error = errutil.ErrBagFull.Error()
		return s.Response(res)
	}
	GeTNum, err := defaultManager.GetItemCount(p.uid, constant.ItemTypeEnergy)
	if err != nil {
	}
	p.takeSource(constant.ItemTypeEnergy, int64(cfg.Cost), constant.ItemAddMainBattle)
	p.addSourceBatch(Reward, constant.ItemAddMainBattle)

	if req.MainBattleId > p.mainBattleId {
		p.mainBattleId = req.MainBattleId
		//完成任务
		p.tasks.dispatch(pb.TargetType_MAIN_BATTLE, int32(req.MainBattleId), 1)
		p.rankChangeTCP(constant.RankTypeMainBattle, int64(p.mainBattleId))
		p.logger.Debugf("主界面战斗,记录新通关ID %s MainBattleId: %d  已通关关卡ID %d", p.name, req.MainBattleId, p.mainBattleId)
		num := int32(0)
		//主线更新的时候，看有没有解锁伙伴槽
		for _, v := range global.GlobalCfg.Openslot.IDLevel {
			if v.ServerDay > 0 {
				if p.level >= v.Needlevel && global.ServerListMgrInstance.GetServerOpenDay(p.serverID) >= v.Openday && p.mainBattleId >= v.Minbattle {
					num++
				}
			} else {
				if p.level >= v.Needlevel && p.createDay >= v.Openday && p.mainBattleId >= v.Minbattle {
					num++
				}
			}

		}
		p.tasks.dispatch(pb.TargetType_UNLOCK_PARTNER_SLOT, 0, num)
	}

	p.dirty(true)

	for _, v := range Reward {
		res.Items = append(res.Items, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	}
	//还有其他奖励，开箱奖励
	preview, err := p.openbox()
	if err != nil {
		res.Code = -5
		res.Error = err.Error()
		return s.Response(res)
	}
	if preview != nil {
		res.Items = append(res.Items, preview...)
	}
	roommgr.upMainPveShushuDate(p, Reward, cfg.Cost, int32(GeTNum), p.mainBattleId, cfg.Name)
	return s.Response(res)
}
func (roommgr *RoomManager) upMainPveShushuDate(p *Player, sourceInfo []global.IdNum, stamina_num int32, stamina_before int32, MainBattleId int32, mainName string) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 主角信息
	shushuDate["name"] = "main_pve"
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			//查询itemname
			if ItemCfg, err := global.GetItemCfg(int32(v.ID)); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["get_item_info"] = get_item_info
	shushuDate["stamina_num"] = stamina_num
	shushuDate["stamina_before"] = stamina_before
	shushuDate["stamina_after"] = stamina_before - stamina_num
	shushuDate["main_id"] = MainBattleId
	shushumain.ReportData(shushuDate)

	shushuUserSetDate := make(map[string]interface{}, 0)

	shushuUserSetDate["current_mainbattle_name"] = mainName
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)

}

// 获取所有战场房间
func (roommgr *RoomManager) GetRoomList(RoomType int32, zoneId int32) []*Desk {
	var desks []*Desk
	for _, v := range roommgr.rooms {
		if v.opt.Roomtype == RoomType {
			for _, d := range v.desks {
				if d.zoneID == zoneId {
					desks = append(desks, d)
				}
			}
		}
	}
	return desks
}

func (roommgr *RoomManager) WarRoomList(s *session.Session, req *pb.WarRoomListRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.WarRoomListResponse{Rooms: []*pb.WarRoomListResponse_RoomInfo{}}

	desks := roommgr.GetRoomList(constant.PVEWarMap, p.getZoneId())
	for _, v := range desks {
		res.Rooms = append(res.Rooms, &pb.WarRoomListResponse_RoomInfo{
			RoomCfgId:    v.cfgID,
			CurPlayerNum: int32(len(v.players) + len(v.robots)),
			MaxPlayerNum: int32(global.GlobalCfg.War.EnterLimit.MaxNum),
			DeskUid:      v.deskID})
	}

	return s.Response(res)
}

// 命运金币
func (roommgr *RoomManager) FateCoin(s *session.Session, data *pb.FateCoinRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	res := &pb.FateCoinResponse{}
	desk := p.desk
	if desk == nil || desk.isDestroy() {
		res.Code = pb.ErrorCode_ERROR_CODE_SCENE_NO_EXIST
		return s.Response(res)
	}
	res.Code = pb.ErrorCode_ERROR_CODE_SUCCESS
	if constant.BattleMap == desk.roomCfg.Roomtype &&
		desk.FateCoin != nil {
		res.Code, res.CfgId = desk.FateCoin.onEffect(p)
	}
	return s.Response(res)
}

// 战役跳过战斗
func (roommgr *RoomManager) SkipBattle(s *session.Session, _ *pb.SkipBattleRequest) error {
	resp := &pb.CommonResponse{RequestCmdId: uint32(pb.CmdId_C2S_SkipBattle)}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = int32(pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND)
		resp.Error = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND.String()
		return s.Response(resp)
	}

	if p.desk == nil {
		resp.Code = int32(pb.ErrorCode_ERROR_CODE_SCENE_NO_EXIST)
		resp.Error = pb.ErrorCode_ERROR_CODE_SCENE_NO_EXIST.String()
		return s.Response(resp)
	}
	p.desk.chSkipBattle <- p.uid

	return s.Response(resp)
}
