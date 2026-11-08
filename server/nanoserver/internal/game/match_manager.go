package game

import (
	"fmt"
	"math"
	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
)

type (
	MatchPlayer struct {
		*Player
		startTime time.Time //匹配开始时间
		ready     bool
		win       bool
		bornIndex int
		killNum   int
	}
	MatchMonster struct {
		monsters map[int64]*MonsterInfo
		win      bool
		killNum  int
	}
	MatchInfo struct {
		*Desk
		matchId          int64
		roomId           int64
		mplayers         []*MatchPlayer
		startTime        time.Time //比赛开始时间
		mtype            pb.MatchType
		buffRefreshTime  time.Time
		boardRefreshTime time.Time
		dirty            bool
	}
	MatchManager struct {
		component.Base
		lock        sync.Mutex
		sessionOpen bool

		matchingOpenDuel                                    bool //单挑开启
		matchingOpenMelee                                   bool //混战开启
		curSessionNum                                       int
		dailyBeginDuel, dailyEndDuel, dailSettlementDuel    time.Time
		dailyBeginMelee, dailyEndMelee, dailSettlementMelee time.Time

		matchingListDuel       []*MatchPlayer       //匹配列表单挑
		matchingListMelee      []*MatchPlayer       //匹配列表混战
		fightingList           map[int64]*MatchInfo //比赛列表
		updateMatchingListTime time.Time
		matchResult            chan int64
	}
)

var matchManager = NewMatchManager()

func NewMatchManager() *MatchManager {
	mgr := &MatchManager{
		matchingListDuel: []*MatchPlayer{},
		fightingList:     make(map[int64]*MatchInfo),
		matchResult:      make(chan int64, chLen)}

	return mgr
}

func (mm *MatchManager) Init() {}

func (mm *MatchManager) AfterInit() {
	for _, v := range global.GlobalCfg.PVPSeason.Season {
		if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
			if time.Now().Before(start) {
				continue
			}
		} else {
			logger.Error("PVP赛季时间配置错误")
			continue
		}

		if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
			if time.Now().After(end) {
				continue
			}
		}
		mm.curSessionNum = v.ID
		mm.sessionOpen = true
		mm.updateMatchingListTime = time.Now()
		break
	}
	var err error
	if mm.dailyBeginDuel, err = time.Parse(time.TimeOnly, global.GlobalCfg.PVP.Settlement.Dailybegin); err != nil {
		logger.Errorf("PVP每日开启时间配置错误 %s", err.Error())
	}
	if mm.dailyEndDuel, err = time.Parse(time.TimeOnly, global.GlobalCfg.PVP.Settlement.Dailyend); err != nil {
		logger.Errorf("PVP每日结束时间配置错误 %s", err.Error())
	}
	if mm.dailSettlementDuel, err = time.Parse(time.TimeOnly, global.GlobalCfg.PVP.Settlement.Settlement); err != nil {
		logger.Errorf("PVP每日结算时间配置错误 %s", err.Error())
	}

	if mm.dailyBeginMelee, err = time.Parse(time.TimeOnly, global.GlobalCfg.GVG.Settlement.Dailybegin); err != nil {
		logger.Errorf("PVP每日开启时间配置错误 %s", err.Error())
	}
	if mm.dailyEndMelee, err = time.Parse(time.TimeOnly, global.GlobalCfg.GVG.Settlement.Dailyend); err != nil {
		logger.Errorf("PVP每日结束时间配置错误 %s", err.Error())
	}
	if mm.dailSettlementMelee, err = time.Parse(time.TimeOnly, global.GlobalCfg.GVG.Settlement.Settlement); err != nil {
		logger.Errorf("PVP每日结算时间配置错误 %s", err.Error())
	}

	mm.matchingOn()

	// 匹配功能已关闭
	// scheduler.NewTimer(time.Millisecond*30, func() {
	// 	mm.checkMatchingDuel()
	// 	mm.checkMatchingMelee()
	// })

	// scheduler.NewTimer(time.Minute*5, func() {
	// 	robotMgr.updateRank()
	// })

	scheduler.NewTimer(time.Second*1, func() {
		mm.checkMatchingTime()
	})

	scheduler.NewTimer(time.Millisecond*10, func() {
	ctrl:
		for {
			select {
			case mid := <-mm.matchResult:
				match := mm.matchInfoById(mid)
				if match == nil {
					continue
				}
				if match.mtype == pb.MatchType_SINGLE_COMBAT {
					mm.matchOverDuel(match)
				} else if match.mtype == pb.MatchType_FREE_FOR_ALL {
					mm.matchOverMelee(match)
				}

			default:
				break ctrl
			}
		}
	})
}

func (mm *MatchManager) checkMatchingTime() {
	for _, v := range global.GlobalCfg.PVPSeason.Season {
		if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
			if time.Now().Before(start) {
				continue
			}
		} else {
			logger.Error("PVP赛季时间配置错误")
			continue
		}

		if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
			if time.Now().After(end) {
				if mm.sessionOpen {
					mm.sessionOpen = false
				}
				continue
			}
		}
		mm.curSessionNum = v.ID
		mm.sessionOpen = true
		break
	}
	mm.matchingOn()
}

func (mm *MatchManager) matchingOn() {
	if mm.sessionOpen {
		t := time.Now()
		startTime := time.Date(t.Year(), t.Month(), t.Day(), mm.dailyBeginDuel.Hour(), mm.dailyBeginDuel.Minute(), 0, 0, t.Location())
		endTime := time.Date(t.Year(), t.Month(), t.Day(), mm.dailyEndDuel.Hour(), mm.dailyEndDuel.Minute(), 0, 0, t.Location())
		if t.After(startTime) && t.Before(endTime) {
			mm.matchingOpenDuel = true
		} else {
			if !mm.matchingOpenDuel {
				mm.matchingOff()
			}
		}

		startTimeMelee := time.Date(t.Year(), t.Month(), t.Day(), mm.dailyBeginMelee.Hour(), mm.dailyBeginMelee.Minute(), 0, 0, t.Location())
		endTimeMelee := time.Date(t.Year(), t.Month(), t.Day(), mm.dailyEndMelee.Hour(), mm.dailyEndMelee.Minute(), 0, 0, t.Location())
		if t.After(startTimeMelee) && t.Before(endTimeMelee) {
			mm.matchingOpenMelee = true
		} else {
			if !mm.matchingOpenMelee {
				mm.matchingOpenMelee = false
			}
		}
	}
}

func (mm *MatchManager) matchingOff() {
	mm.matchingOpenDuel = false
}

func (mm *MatchManager) matchInfoById(id int64) *MatchInfo {
	info, ok := mm.fightingList[id]
	if ok {
		return info
	}
	return nil
}

func (mm *MatchManager) robotCountDuel(level int) int {
	num := 0
	for _, mp := range mm.matchingListDuel {
		if mp.Player.isRobot() {
			if maputil.ABS(mp.Player.level-int32(level)) < 30 {
				num++
			}
		}
	}
	return num
}

func (mm *MatchManager) robotCountMelee(level int) int {
	num := 0
	for _, mp := range mm.matchingListMelee {
		if mp.Player.isRobot() {
			if maputil.ABS(mp.Player.level-int32(level)) < 30 {
				num++
			}
		}
	}
	return num
}

// 匹配1v1
func (mm *MatchManager) checkMatchingDuel() {
	// mm.lock.Lock()
	// defer mm.lock.Unlock()

	// if len(mm.matchingListDuel) == 0 {
	// 	return
	// }
	// for _, mp := range mm.matchingListDuel {
	// 	if mp.curState == PlayerStatusOffLine || mp.session == nil || time.Now().After(mp.startTime.Add(time.Duration(global.GlobalCfg.PVP.TimeSet.MatchTime*2)*time.Second)) {
	// 		mm.removeMatchingPlayerDuel(mp)
	// 		packet := &pb.MatchingResultPush{Status: 2}
	// 		mp.pushMsg(protocol.RouteMatchingResultPush, packet)
	// 		break
	// 	}
	// 	if mp.Player.isRobot() {
	// 		continue
	// 	}
	// 	//15 秒未匹配到 加入机器人
	// 	maxMatchingTime := global.GlobalCfg.PVP.TimeSet.MatchTime
	// 	if time.Now().After(mp.startTime.Add(time.Duration(maxMatchingTime)*time.Second)) || (mp.losingStreakCount >= int32(global.GlobalCfg.PVP.Match.SecurityTimes) && time.Now().After(mp.startTime.Add(2*time.Second))) {
	// 		if mm.robotCountDuel(int(mp.Player.level)) <= 0 {
	// 			if r := robotMgr.getFreeRobot(mp.level); r != nil {
	// 				mm.matchingListDuel = append(mm.matchingListDuel, &MatchPlayer{Player: r.Player, startTime: time.Now()})
	// 				logger.Debugf("匹配超过15秒 加入机器人 %s   等级 %d", r.name, r.level)
	// 			}
	// 		}
	// 	}
	// }

	// if len(mm.matchingListDuel) < global.GlobalCfg.PVP.Match.SinglePVP+1 {
	// 	return
	// }

	// sort.Slice(mm.matchingListDuel, func(i, j int) bool {
	// 	return mm.matchingListDuel[i].level < mm.matchingListDuel[j].level
	// })

	// levellimit := int32(global.GlobalCfg.PVP.LevelSet.LevelLimit)
	// var player1, player2 *MatchPlayer
	// minDiff := int32(9999) // 最大的int32值
	// for i := 0; i < len(mm.matchingListDuel); i++ {
	// 	if mm.matchingListDuel[i].session == nil {
	// 		continue
	// 	}

	// 	if i < len(mm.matchingListDuel)-1 {
	// 		if mm.matchingListDuel[i].Player.isRobot() && mm.matchingListDuel[i+1].Player.isRobot() {
	// 			continue
	// 		}
	// 		diff := maputil.ABS(mm.matchingListDuel[i].level - mm.matchingListDuel[i+1].level)
	// 		if diff < minDiff {
	// 			minDiff = diff
	// 			player1 = mm.matchingListDuel[i]
	// 			player2 = mm.matchingListDuel[i+1]
	// 		}
	// 	}
	// }

	// if player1 == nil || player2 == nil {
	// 	return
	// }
	// if maputil.ABS(player1.level-player2.level) >= 30 {
	// 	return
	// }

	// if maputil.ABS(player1.level-player2.level) > levellimit && time.Now().After(player1.startTime.Add(3*time.Second)) {
	// 	return
	// }
	// if maputil.ABS(player1.level-player2.level) > levellimit*2 && time.Now().After(player1.startTime.Add(6*time.Second)) {
	// 	return
	// }

	// if maputil.ABS(player1.level-player2.level) < 30 {
	// 	logger.Debugf("Matching players: %s and %s\n", player1.name, player2.name)

	// 	matchplayerinfo2 := &pb.PlayerInfo{
	// 		Uid:      player2.uid,
	// 		Nickname: player2.name,
	// 		HeadUrl:  player2.head,
	// 		MaxHp:    player2.stat(pb.Attribute_MAX_HEALTH),
	// 		CurHp:    player2.currentHP,
	// 		Job:      int32(player2.job),
	// 		Gender:   int32(player2.sex),
	// 		Level:    player2.level,
	// 	}
	// 	matchplayerinfo1 := &pb.PlayerInfo{
	// 		Uid:      player1.uid,
	// 		Nickname: player1.name,
	// 		HeadUrl:  player1.head,
	// 		MaxHp:    player1.stat(pb.Attribute_MAX_HEALTH),
	// 		CurHp:    player1.currentHP,
	// 		Job:      int32(player1.job),
	// 		Gender:   int32(player1.sex),
	// 		Level:    player1.level,
	// 	}
	// 	packet := &pb.MatchingResultPush{
	// 		Status:  0,
	// 		Players: []*pb.PlayerInfo{}}
	// 	packet.Players = append(packet.Players, matchplayerinfo1)
	// 	packet.Players = append(packet.Players, matchplayerinfo2)

	// 	player1.pushMsg(protocol.RouteMatchingResultPush, packet)
	// 	player2.pushMsg(protocol.RouteMatchingResultPush, packet)

	// 	player1.bornIndex = 0
	// 	player2.bornIndex = 1

	// 	matchid := snowflake.GenID()
	// 	matchInfo := &MatchInfo{matchId: matchid, mplayers: []*MatchPlayer{}, startTime: time.Now()}
	// 	matchInfo.mplayers = append(matchInfo.mplayers, player1)
	// 	matchInfo.mplayers = append(matchInfo.mplayers, player2)
	// 	matchInfo.mtype = pb.MatchType_SINGLE_COMBAT

	// 	mm.fightingList[matchid] = matchInfo
	// 	// 删除已经匹配成功的玩家
	// 	mm.removeMatchingPlayerDuel(player1)
	// 	mm.removeMatchingPlayerDuel(player2)

	// 	if room := roomManager.getRoomByType(constant.PVPOneOnOneMap, 1); room != nil {
	// 		for _, d := range room.desks {
	// 			if d.status() == constant.DeskStatusDestory {
	// 				continue
	// 			}
	// 			if len(d.players) > 0 {
	// 				continue
	// 			}
	// 			if d.mapInfo == nil {
	// 				d.setStatus(constant.DeskStatusDestory)
	// 				continue
	// 			}
	// 			mm.fightingList[matchid].roomId = d.deskID
	// 			mm.fightingList[matchid].Desk = d
	// 			d.matchUid = matchid
	// 			d.MatchInfo = mm.fightingList[matchid]

	// 			if err := d.playerJoin(player1.Player, false, 0); err == nil {
	// 				packet := &pb.MatchRoomJoinPush{RoomCfgId: d.cfgID, RoomUid: d.deskID, Location: player1.location, Direction: player1.direction}
	// 				player1.pushMsg(protocol.RouteMatchRoomJoinPush, packet)
	// 			}
	// 			if err := d.playerJoin(player2.session, false, 0); err == nil {
	// 				packet := &pb.MatchRoomJoinPush{RoomCfgId: d.cfgID, RoomUid: d.deskID, Location: player2.location, Direction: player2.direction}
	// 				player2.pushMsg(protocol.RouteMatchRoomJoinPush, packet)
	// 			}

	// 			d.syncDeskStatus()
	// 			logger.Debugf("匹配成功  玩家1: ID=%d  玩家2: ID=%d  队列剩余人数=%d", player1.uid, player2.uid, len(mm.matchingListDuel))
	// 			break
	// 		}
	// 	}
	// }
}

// 匹配混战
func (mm *MatchManager) checkMatchingMelee() {
	// mm.lock.Lock()
	// defer mm.lock.Unlock()

	// if len(mm.matchingListMelee) == 0 {
	// 	return
	// }
	// for _, mp := range mm.matchingListMelee {
	// 	if mp.curState == PlayerStatusOffLine || mp.session == nil || time.Now().After(mp.startTime.Add(time.Duration(global.GlobalCfg.GVG.TimeSet.MatchTime*2)*time.Second)) {
	// 		mm.removeMatchingPlayerMelee(mp)
	// 		packet := &pb.MatchingResultPush{Status: 2}
	// 		mp.pushMsg(protocol.RouteMatchingResultPush, packet)
	// 		break
	// 	}
	// 	if mp.Player.isRobot() {
	// 		continue
	// 	}

	// }

	// sort.Slice(mm.matchingListMelee, func(i, j int) bool {
	// 	return mm.matchingListMelee[i].level < mm.matchingListMelee[j].level
	// })
	// var result []*MatchPlayer

	// for i := 0; i < len(mm.matchingListMelee); i++ {
	// 	count := 1
	// 	sumLevel := mm.matchingListMelee[i].level
	// 	var matchedPlayers []*MatchPlayer
	// 	matchedPlayers = append(matchedPlayers, mm.matchingListMelee[i])

	// 	for j := i + 1; j < len(mm.matchingListMelee); j++ {
	// 		if count == 6 {
	// 			break
	// 		}

	// 		diff := mm.matchingListMelee[i].level - mm.matchingListMelee[j].level
	// 		if diff >= -30 && diff <= 30 {
	// 			count++
	// 			sumLevel += mm.matchingListMelee[j].level
	// 			matchedPlayers = append(matchedPlayers, mm.matchingListMelee[j])
	// 		}
	// 	}

	// 	if count == global.GlobalCfg.GVG.Match.MultiplePVP {
	// 		result = append(result, matchedPlayers...)
	// 		break
	// 	} else {
	// 		result = append(result, matchedPlayers...)
	// 		break
	// 	}
	// }
	// //15 秒未匹配到 加入机器人
	// //maxMatchingTime := global.GlobalCfg.GVG.TimeSet.MatchTime
	// if len(result) > 0 && len(result) < global.GlobalCfg.GVG.Match.MultiplePVP {
	// 	fmp := result[0]
	// 	for _, fp := range result {
	// 		if !fp.isRobot() {
	// 			fmp = fp
	// 			break
	// 		}
	// 	}

	// 	if len(result) < 3 && time.Now().After(fmp.startTime.Add(5*time.Second)) {
	// 		if r := robotMgr.getFreeRobot(fmp.level); r != nil {
	// 			mm.matchingListMelee = append(mm.matchingListMelee, &MatchPlayer{Player: r.Player, startTime: time.Now()})
	// 			logger.Debugf("匹配超过5秒 加入机器人 %s   等级 %d", r.name, r.level)
	// 			return
	// 		}
	// 	}
	// 	ct := global.RandIn(5) + 6
	// 	if len(result) < 5 && time.Now().After(fmp.startTime.Add(time.Duration(ct)*time.Second)) {
	// 		if r := robotMgr.getFreeRobot(fmp.level); r != nil {
	// 			mm.matchingListMelee = append(mm.matchingListMelee, &MatchPlayer{Player: r.Player, startTime: time.Now()})
	// 			logger.Debugf("匹配超过 %d 秒 加入机器人 %s   等级 %d", ct, r.name, r.level)
	// 			return
	// 		}
	// 	}
	// 	ct = global.RandIn(5) + 11
	// 	if len(result) < 6 && time.Now().After(fmp.startTime.Add(time.Duration(ct)*time.Second)) {
	// 		if r := robotMgr.getFreeRobot(fmp.level); r != nil {
	// 			mm.matchingListMelee = append(mm.matchingListMelee, &MatchPlayer{Player: r.Player, startTime: time.Now()})
	// 			logger.Debugf("匹配超过 %d 秒 加入机器人 %s   等级 %d", ct, r.name, r.level)
	// 			return
	// 		}
	// 	}
	// }
	// var bUpdate bool
	// if time.Now().After(mm.updateMatchingListTime) {
	// 	bUpdate = true
	// 	mm.updateMatchingListTime = time.Now().Add(time.Second)
	// }
	// if bUpdate && len(result) < 6 {
	// 	packet := &pb.MatchingResultPush{
	// 		Players: []*pb.PlayerInfo{},
	// 		Status:  1}

	// 	for _, p := range result {
	// 		packet.Players = append(packet.Players, &pb.PlayerInfo{
	// 			Uid:      p.uid,
	// 			Nickname: p.name,
	// 			HeadUrl:  p.head,
	// 			MaxHp:    p.stat(pb.Attribute_MAX_HEALTH),
	// 			CurHp:    p.currentHP,
	// 			Job:      int32(p.job),
	// 			Gender:   int32(p.sex),
	// 			Level:    p.level,
	// 		})
	// 	}
	// 	for _, p := range result {
	// 		p.pushMsg(protocol.RouteMatchingResultPush, packet)
	// 	}
	// 	bUpdate = false
	// }
	// if len(result) == 6 {
	// 	packet := &pb.MatchingResultPush{
	// 		Status:  0,
	// 		Players: []*pb.PlayerInfo{}}

	// 	for _, p := range result {
	// 		packet.Players = append(packet.Players, &pb.PlayerInfo{
	// 			Uid:      p.uid,
	// 			Nickname: p.name,
	// 			HeadUrl:  p.head,
	// 			MaxHp:    p.stat(pb.Attribute_MAX_HEALTH),
	// 			CurHp:    p.currentHP,
	// 			Job:      int32(p.job),
	// 			Gender:   int32(p.sex),
	// 			Level:    p.level,
	// 		})
	// 	}
	// 	matchid := snowflake.GenID()
	// 	matchInfo := &MatchInfo{matchId: matchid, mplayers: []*MatchPlayer{}, startTime: time.Now()}
	// 	matchInfo.mtype = pb.MatchType_FREE_FOR_ALL
	// 	matchInfo.buffRefreshTime = time.Now().Add(time.Duration(global.GlobalCfg.GVG.BuffGap.RefreshTime) * time.Second)
	// 	matchInfo.boardRefreshTime = time.Now()

	// 	for i, p := range result {
	// 		p.pushMsg(protocol.RouteMatchingResultPush, packet)
	// 		p.bornIndex = i
	// 		matchInfo.mplayers = append(matchInfo.mplayers, p)
	// 		mm.removeMatchingPlayerMelee(p)
	// 	}

	// 	mm.fightingList[matchid] = matchInfo

	// 	if room := roomManager.getRoomByType(constant.PVPMultiOnMultiMap, 1); room != nil {
	// 		for _, d := range room.desks {
	// 			if len(d.players) > 0 {
	// 				continue
	// 			}
	// 			if d.mapInfo == nil {
	// 				d.setStatus(constant.DeskStatusDestory)
	// 				logger.Printf("房间异常 无地图信息  roomID = %d", d.cfgID)
	// 				continue
	// 			}
	// 			mm.fightingList[matchid].roomId = d.deskID
	// 			mm.fightingList[matchid].Desk = d
	// 			d.matchUid = matchid
	// 			d.MatchInfo = mm.fightingList[matchid]

	// 			for _, p := range result {
	// 				if err := d.playerJoin(p.session, false, 0); err == nil {
	// 					packet := &pb.MatchRoomJoinPush{RoomCfgId: d.cfgID, RoomUid: d.deskID, Location: p.location, Direction: p.direction}
	// 					p.pushMsg(protocol.RouteMatchRoomJoinPush, packet)
	// 				} else {
	// 					//有一个人加入失败 ,房间结束
	// 					mm.matchResult <- d.matchUid
	// 					logger.Debugf("有玩家加入房间失败 ,房间结束 Uid=%d , name= %s", p.uid, p.name)
	// 					return
	// 				}
	// 			}

	// 			d.syncDeskStatus()
	// 			logger.Debugf("混战模式 匹配成功 ，队列剩余人数=%d", len(mm.matchingListDuel))
	// 			break
	// 		}
	// 	}
	// }
}

func (mm *MatchManager) removeMatchingPlayerDuel(player *MatchPlayer) {
	for i, p := range mm.matchingListDuel {
		if p == player {
			mm.matchingListDuel = append(mm.matchingListDuel[:i], mm.matchingListDuel[i+1:]...)
			break
		}
	}
}

func (mm *MatchManager) removeMatchingPlayerMelee(player *MatchPlayer) {
	for i, p := range mm.matchingListMelee {
		if p == player {
			mm.matchingListMelee = append(mm.matchingListMelee[:i], mm.matchingListMelee[i+1:]...)
			break
		}
	}
}

// 游戏结束(1v1)
func (mm *MatchManager) matchOverDuel(m *MatchInfo) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if _, ok := mm.fightingList[m.matchId]; !ok {
		return
	}
	var winid int64

	for _, p := range m.mplayers {
		if p.win {
			winid = p.uid
			break
		}
	}
	duration := time.Now().Unix() - m.startTime.Unix()
	var p1, p2 *MatchPlayer
	if winid == 0 {
		if len(m.mplayers) == 2 {
			p1 = m.mplayers[0]
			p2 = m.mplayers[1]

			if p1.currentHP > p2.currentHP {
				p1.win = true
			} else if p1.currentHP < p2.currentHP {
				p2.win = true
			} else {
				if p1.combat < p2.combat {
					p1.win = true
				} else {
					p2.win = true
				}
			}
			// 掉线的输
			if p1.curState == PlayerStatusOffLine && p2.curState == PlayerStatusOffLine {
				p1.win = false
				p1.win = false
			} else if p1.curState <= PlayerStatusPlaying && p2.curState >= PlayerStatusDisconnect {
				p1.win = true
				p2.win = false
			} else if p1.curState >= PlayerStatusDisconnect && p2.curState <= PlayerStatusPlaying {
				p1.win = false
				p2.win = true
			}
		}
	}
	is_first_try := false

	packet := &pb.GameResultMatchPush{Items: []*pb.ItemInfo{}}
	for _, p := range m.mplayers {
		var bRobot bool
		if p.isRobot() {
			// 机器人不算分，离线玩家算分
			if p.robotType == RobotTypeRobot {
				if p.win {
					winid = p.uid
				}
				continue
			} else if p.robotType == RobotTypePlayer {
				bRobot = true
			}
		} else {
			if p.pvpScore == 0 {
				is_first_try = true
			}
		}

		if p.win {
			var items []global.IdNum

			for _, v := range global.GlobalCfg.PVPReward.Win.Keys {
				num := v.Num
				if v.Mark == 1 {
					num = v.Num * int(p.level)
				}
				items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})

				packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
			}
			if !bRobot {
				if err := p.checkAddSourceBatch(items); err == nil {
					p.addSourceBatch(items, constant.ItemAddPvp, WithAddSourceUnNotify(AddSourceTypeUnNotify))
					attachment := []protocol.Attachment{}
					for _, v := range items {
						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
					}
					mm.upPvpOverShushuDate(m.mplayers[0].Player, m.mplayers[1].Player, attachment, true, is_first_try, duration, "即时赛")
				} else {
					attachment := []protocol.Attachment{}
					for _, v := range items {
						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
					}
					p.SendMail(constant.ItemAddPvp, attachment)
					mm.upPvpOverShushuDate(m.mplayers[0].Player, m.mplayers[1].Player, attachment, true, is_first_try, duration, "即时赛")

				}

				if p.pvpScore == 0 {
					p.pvpScore = DefaultPvpScore
				}

				p.winningStreakCount++
				p.losingStreakCount = 0
				p.pvpWinStreakTimes++

				// 跑马灯
				if p.pvpWinStreakTimes >= global.GlobalCfg.PVP.Lamp.LampWinTimes {
					p.pvpWinStreakTimes = 0
					context := []string{p.name}
					defaultManager.Lamp(int32(constant.LampUnstoppableInBrawl), p.serverID, context)
				}

				winid = p.uid
			}
		} else {
			var items []global.IdNum

			for _, v := range global.GlobalCfg.PVPReward.Lose.Keys {
				num := v.Num
				if v.Mark == 1 {
					num = v.Num * int(p.level)
				}
				items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})
				packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
			}

			if !bRobot {
				if err := p.checkAddSourceBatch(items); err == nil {
					p.addSourceBatch(items, constant.ItemAddPvp, WithAddSourceUnNotify(AddSourceTypeUnNotify))
					attachment := []protocol.Attachment{}
					for _, v := range items {
						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
					}
					mm.upPvpOverShushuDate(m.mplayers[0].Player, m.mplayers[1].Player, attachment, false, is_first_try, duration, "即时赛")
				} else {
					attachment := []protocol.Attachment{}
					for _, v := range items {
						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
					}
					p.SendMail(constant.ItemAddPvp, attachment)
					mm.upPvpOverShushuDate(m.mplayers[0].Player, m.mplayers[1].Player, attachment, false, is_first_try, duration, "即时赛")
				}

				//p.pushMsg(protocol.RouteGameResultMatchPush, losePacket)

				p.winningStreakCount = 0
				p.losingStreakCount++
				p.pvpWinStreakTimes = 0
			}
		}
	}

	//算分
	if len(m.mplayers) == 2 {
		blue := m.mplayers[0]
		red := m.mplayers[1]
		changeScore := int64(0)
		pvpScore := mm.calcScore(blue.pvpScore, red.pvpScore)
		if pvpScore != nil {
			if blue.win {
				blue.pvpScore += pvpScore.Win
				red.pvpScore -= pvpScore.Win
				changeScore = int64(pvpScore.Win)
				packet.Score = int32(changeScore)
				packet.Win = true
			} else {
				blue.pvpScore -= pvpScore.Lose
				red.pvpScore += pvpScore.Lose
				changeScore = int64(pvpScore.Lose)
				packet.Score = int32(-changeScore)
				packet.Win = false
			}
		}

		blue.pushMsg(protocol.RouteGameResultMatchPush, packet)

		blue.rankChangeTCP(constant.RankTypeSoul, int64(blue.pvpScore))
		red.rankChangeTCP(constant.RankTypeSoul, int64(red.pvpScore))

		//保存PVP记录
		mm.savePvpRecord(PvpRecordTypeJiShi, blue, red, int64(changeScore))
	}

	var serverId int32
	for _, p := range m.mplayers {
		if !p.isRobot() {
			p.addChallengeRecord(constant.PVPOneOnOneMap)
			serverId = p.serverID
		}
	}

	m.setStatus(constant.DeskStatusIdle)

	delete(mm.fightingList, m.matchId)

	// 战斗测试代码
	//combatManagerr.FightEnd(m.matchUid, winid)

	m.matchUid = 0
	m.MatchInfo = nil
	// m.mplayers[0].Player.tasks.dispatch(pb.TargetType_TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE, 0, 1)

	robotMgr.updateRank(serverId, 0, constant.RankTypeSoul)
	logger.Debugf("1v1比赛结束  deskid=%d  winer=%d", m.deskID, winid)
}
func (mm *MatchManager) upPvpOverShushuDate(p *Player, red *Player, sourceInfo []protocol.Attachment, isSuccess bool, is_first_try bool, duration int64, pvp_mode string) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 主角信息
	//查看战斗信息
	partnerInfo := make([]map[string]interface{}, 0)
	shushuDate["name"] = "realtime_pvp_end"
	shushuDate["pvp_mode"] = pvp_mode

	// 1v1对手信息
	shushuDate["opponent_name"] = red.name
	shushuDate["opponent_id"] = fmt.Sprintf("%d", red.uid)
	shushuDate["opponent_power"] = red.totalCombat
	if red.robotType == RobotTypePlayer {
		shushuDate["is_robot"] = true
	} else {
		shushuDate["is_robot"] = false
	}
	//查询是否是首次挑战
	getChallengeRecord := p.getChallengeRecord(constant.PVPOneOnOneMap)
	if getChallengeRecord == 0 {
		shushuDate["is_first_try"] = true
	} else {
		shushuDate["is_first_try"] = false
	}
	// shushuDate["opponent_rank"] = red.rank
	rankInfo1 := &protocol.RankOperateRequest{
		Uid:  red.uid,
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
	shushuDate["role_rank"] = 0
	rankInfo2 := &protocol.RankOperateRequest{
		Uid:  p.uid,
		Type: 3,
		All:  3,
	}
	res2, err := provide.RankList(rankInfo2)
	if err == nil {
		if res2.RankOwn.Num > 0 {
			shushuDate["role_rank"] = res2.RankOwn.Num
		}
	}
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
	shushuDate["is_success"] = isSuccess
	shushuDate["duration"] = duration

	// 是否首次尝试
	shushuDate["is_first_try"] = is_first_try
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
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			//查询itemname
			if ItemCfg, err := global.GetItemCfg(int32(v.UserPropsId)); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["get_item_info"] = get_item_info
	shushumain.ReportData(shushuDate)

}

// 游戏结束(混战)
func (mm *MatchManager) matchOverMelee(m *MatchInfo) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if _, ok := mm.fightingList[m.matchId]; !ok {
		return
	}

	var totalKill int
	for _, mp := range m.mplayers {
		totalKill += mp.killNum
	}

	sort.Slice(m.mplayers, func(i, j int) bool {
		return m.mplayers[i].killNum > m.mplayers[j].killNum
	})

	var reward []global.XMLIdNum
	// for _, r := range global.GlobalCfg.GVG.Reward {
	// 	if r.Round == mm.curSessionNum%4 {
	// 		reward = r.Items
	// 	}
	// }

	packet := &pb.GameResultMatchMeleePush{Result: []*pb.MeleeResult{}}
	for i := range m.mplayers {
		OneResult := &pb.MeleeResult{
			Player: &pb.PlayerInfo{
				Uid:         m.mplayers[i].uid,
				Nickname:    m.mplayers[i].name,
				OfficialExp: m.mplayers[i].tokens[global.GlobalCfg.BaseValue.Official.ItemID],
			},
			Items: []*pb.ItemInfo{},
		}
		OneResult.KillNum = int32(m.mplayers[i].killNum)
		packet.Result = append(packet.Result, OneResult)
		logger.Debugf("混战比赛结果  uid:%d  name:%s  kill:%d   totalKill:%d ", m.mplayers[i].uid, m.mplayers[i].name, m.mplayers[i].killNum, totalKill)
		if m.mplayers[i].killNum == 0 || totalKill == 0 {
			m.mplayers[i].pvpWinStreakTimes = 0
			continue
		}

		// 跑马灯
		if i == 0 {
			m.mplayers[i].pvpWinStreakTimes++
			m.mplayers[i].pvpScore++

			m.mplayers[i].rankChangeTCP(constant.RankTypeSoul, int64(m.mplayers[i].pvpScore))
			if m.mplayers[i].pvpWinStreakTimes >= global.GlobalCfg.PVP.Lamp.LampWinTimes {
				m.mplayers[i].pvpWinStreakTimes = 0
				context := []string{m.mplayers[i].name}
				defaultManager.Lamp(int32(constant.LampUnstoppableInBrawl), m.mplayers[i].serverID, context)
			}
		}

		var items []global.IdNum
		for _, r := range reward {
			num := float32(r.Num) * float32(m.mplayers[i].killNum) / float32(totalKill)
			items = append(items, global.IdNum{ID: r.ID, Num: int32(num)})
			OneResult.Items = append(OneResult.Items, &pb.ItemInfo{CfgId: r.ID, Num: int64(num)})
		}
		if len(items) == 0 || m.mplayers[i].isRobot() {
			continue
		}

		if err := m.mplayers[i].checkAddSourceBatch(items); err == nil {
			m.mplayers[i].addSourceBatch(items, constant.ItemAddPvp, WithAddSourceUnNotify(AddSourceTypeUnNotify))
		} else {
			attachment := []protocol.Attachment{}
			for _, v := range items {
				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
			}
			m.mplayers[i].SendMail(constant.ItemAddPvp, attachment)
		}
	}

	if len(m.mplayers) > 0 && m.mplayers[0].killNum > 0 {
		m.mplayers[0].winningStreakCount++
		// 跑马灯
	}

	for _, p := range m.mplayers {
		p.pushMsg(protocol.RouteGameResultMatchMeleePush, packet)
	}

	for _, p := range m.mplayers {
		m.playerKick(p.Player, constant.GameOverType)
	}

	m.setStatus(constant.DeskStatusIdle)

	delete(mm.fightingList, m.matchId)

	m.matchUid = 0
	m.MatchInfo = nil

	logger.Debugf("混战比赛结束  deskid=%d  ", m.deskID)
}

func (mm *MatchManager) getBornIndex(id int64, uid int64) int {
	if match := mm.matchInfoById(id); match != nil {
		for _, v := range match.mplayers {
			if v.uid == uid {
				return v.bornIndex
			}
		}
	}
	return 0
}

func (mm *MatchManager) leave(matchid int64, uid int64) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if match := mm.matchInfoById(matchid); match != nil {
		for _, v := range match.mplayers {
			if v.uid == uid {
				v.win = false
				logger.Debugf("1v1中途退出判定失败 uid=%d name=%s  matchid=%d", uid, v.name, matchid)
				continue
			}
			if match.mtype == pb.MatchType_SINGLE_COMBAT {
				v.win = true
				mm.matchResult <- matchid
			}
		}
	}
}

// 保存pvp记录
func (mm *MatchManager) savePvpRecord(rtype PvpRecordType, blue *MatchPlayer, red *MatchPlayer, score int64) {

	now := time.Now().Unix()

	//玩家1记录
	recoder1 := model.PvpRecord{RecordType: int8(rtype)}
	recoder1.Uid = snowflake.GenID()
	recoder1.SelfUid = blue.uid
	recoder1.Status = 1
	recoder1.OpponentUid = red.uid
	if blue.win {
		recoder1.Result = 1
		recoder1.Score = score
	} else {
		recoder1.Result = 0
		recoder1.Score = -score
	}
	recoder1.MatchTime = now
	recoder1.Name = red.name
	recoder1.Power = int64(red.totalCombat)
	if red.robotType == RobotTypePlayer {
		if redPlayer, ok := defaultManager.player(red.uid); ok {
			recoder1.Power = int64(redPlayer.totalCombat)
		}
	}
	recoder1.Gender = int8(red.sex)
	recoder1.Job = int8(red.job)
	avatar := int32(0)
	frame := int32(0)
	if red.avatar != nil {
		avatar = red.avatar.GetCurrentAvatar()
		frame = red.avatar.GetCurrentFrame()
	}
	recoder1.AvatarV2 = avatar
	recoder1.Frame = frame
	recoder1.Level = red.level

	if rtype == PvpRecordTypeHuiHe {
		recoder1.SelfScore = int64(blue.gvgScore)
		recoder1.OpponentScore = int64(red.gvgScore)
	} else {
		recoder1.SelfScore = int64(blue.pvpScore)
		recoder1.OpponentScore = int64(red.pvpScore)
	}

	blue.pvpRecorder = append(blue.pvpRecorder, &recoder1)
	go func() {
		db.InsertPvpRecord(&recoder1)
	}()

	//删除超过3天的pvp记录
	del_ids1 := []int64{}
	restRecoder1 := []*model.PvpRecord{}
	for _, v := range blue.pvpRecorder {
		if now-v.MatchTime > 3*24*3600 {
			del_ids1 = append(del_ids1, v.ID)
		} else {
			restRecoder1 = append(restRecoder1, v)
		}
	}
	blue.pvpRecorder = restRecoder1

	for _, v := range del_ids1 {
		go db.DeletePvpRecord(v)
	}

	//玩家2记录
	if red.robotType == RobotTypePlayer {
		redPlayer := defaultManager.ObserveUserBaseInfo(red.serverID, red.uid)
		recoder2 := model.PvpRecord{RecordType: int8(rtype)}
		recoder2.Uid = snowflake.GenID()
		recoder2.SelfUid = redPlayer.uid
		recoder2.Status = 1
		recoder2.OpponentUid = blue.uid
		if red.win {
			recoder2.Result = 1
			recoder2.Score = score
		} else {
			recoder2.Result = 0
			recoder2.Score = -score
		}
		recoder2.MatchTime = now
		recoder2.Name = blue.name
		recoder2.Power = int64(blue.totalCombat)
		recoder2.Gender = int8(blue.sex)
		recoder2.Job = int8(blue.job)
		recoder2.Level = blue.level

		if rtype == PvpRecordTypeHuiHe {
			recoder2.SelfScore = int64(red.gvgScore)
			recoder2.OpponentScore = int64(blue.gvgScore)
		} else {
			recoder2.SelfScore = int64(red.pvpScore)
			recoder2.OpponentScore = int64(blue.pvpScore)
		}

		redPlayer.pvpRecorder = append(redPlayer.pvpRecorder, &recoder2)

		go func() {
			db.InsertPvpRecord(&recoder2)
		}()

		//删除超过3天的pvp记录
		del_ids2 := []int64{}
		restRecoder2 := []*model.PvpRecord{}
		for _, v := range redPlayer.pvpRecorder {
			if now-v.MatchTime > 3*24*3600 {
				del_ids2 = append(del_ids2, v.ID)
			} else {
				restRecoder2 = append(restRecoder2, v)
			}
		}
		redPlayer.pvpRecorder = restRecoder2

		for _, v := range del_ids2 {
			go db.DeletePvpRecord(v)
		}
	}

}

func (mm *MatchManager) Matching(s *session.Session, req *pb.MatchingRequest) error {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	if p.desk != nil {
		return errutil.ErrDeskNotFound
	}

	res := &pb.MatchingResponse{Cancel: req.Cancel, MatchType: req.MatchType}

	res.Code = -1
	res.Error = errutil.ModClosed
	return s.Response(res)

	// if p.level < global.SwitchCfg[ModuleLuanDouZhiWang] {
	// 	res.Code = -9
	// 	res.Error = errutil.ErrLevel.Error()
	// 	return s.Response(res)
	// }
	switch req.MatchType {
	case pb.MatchType_SINGLE_COMBAT:
		if !req.Cancel {
			room := roomManager.getRoomByType(int32(req.MatchType), p.level)
			if room == nil {
				res.Code = -1
				return s.Response(res)
			}
			if !mm.matchingOpenDuel {
				res.Code = -2
				res.Error = errutil.EventIsPending
				return s.Response(res)
			}

			if p.PVPInfo.PVPDuel.dailyMatchTimes >= int32(global.GlobalCfg.PVP.DailySet.DailyTime+int(p.PVPInfo.PVPDuel.dailyBuyNum)) {
				res.Code = -3
				res.Error = errutil.LimitHasBeenReached
				return s.Response(res)
			}

			if len(mm.matchingListDuel) >= 100 {
				res.Code = -4
				res.Error = errutil.ServerIsBusy
				return s.Response(res)
			}

			for _, mp := range mm.matchingListDuel {
				if mp.uid == p.uid {
					res.Code = -5
					return s.Response(res)
				}
			}

			matchp := &MatchPlayer{Player: p, startTime: time.Now(), ready: false, win: false}
			mm.matchingListDuel = append(mm.matchingListDuel, matchp)
		} else {
			for i, mp := range mm.matchingListDuel {
				if mp.uid == p.uid {
					mm.matchingListDuel = append(mm.matchingListDuel[:i], mm.matchingListDuel[i+1:]...)
				}
			}
		}
	case pb.MatchType_FREE_FOR_ALL:
		if !req.Cancel {
			room := roomManager.getRoomByType(int32(req.MatchType), p.level)
			if room == nil {
				res.Code = -1
				return s.Response(res)
			}
			if !mm.matchingOpenMelee {
				res.Code = -2
				res.Error = errutil.EventIsPending
				return s.Response(res)
			}

			if p.PVPInfo.PVPMelee.dailyMatchTimes >= int32(global.GlobalCfg.GVG.DailySet.DailyTime+int(p.PVPInfo.PVPMelee.dailyBuyNum)) {
				res.Code = -3
				res.Error = errutil.LimitHasBeenReached
				return s.Response(res)
			}

			if len(mm.matchingListMelee) >= 100 {
				res.Code = -4
				res.Error = errutil.ServerIsBusy
				return s.Response(res)
			}

			for _, mp := range mm.matchingListMelee {
				if mp.uid == p.uid {
					res.Code = -5
					return s.Response(res)
				}
			}

			matchp := &MatchPlayer{Player: p, startTime: time.Now(), ready: false, win: false}
			mm.matchingListMelee = append(mm.matchingListMelee, matchp)
			logger.Debugf("混战匹配  uid:%d level:%d name:%s  队列人数:%d", p.uid, p.level, p.name, len(mm.matchingListMelee))
		} else {
			for i, mp := range mm.matchingListMelee {
				if mp.uid == p.uid {
					mm.matchingListMelee = append(mm.matchingListMelee[:i], mm.matchingListMelee[i+1:]...)
				}
			}
		}
	}

	return s.Response(res)
}

// 准备
func (mm *MatchManager) Ready(s *session.Session, req *pb.MatchingReadyNotify) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	if p.desk == nil {
		return nil
	}

	matchInfo := mm.matchInfoById(p.desk.matchUid)
	if matchInfo == nil {
		return nil
	}

	for _, v := range matchInfo.mplayers {
		if v.isRobot() {
			v.ready = true
			continue
		}
		if v.uid == p.uid {
			v.ready = true
		}
	}

	readyNum := 0
	for _, v := range matchInfo.mplayers {
		if v.ready && !v.isRobot() {
			readyNum++
		}
	}
	if readyNum == p.desk.realPlayerCount() {
		p.desk.Broadcast(protocol.RouteMatchingStartPush, &pb.MatchingStartPush{})

		for i := range p.desk.players {
			actionTime := time.Now().Add(3100 * time.Millisecond)
			p.desk.players[i].moveTime = actionTime
			p.desk.players[i].attackTime = actionTime
			//logger.Debugf("主角 %d:%s  攻击时间 %v  移动时间 %v ", p.desk.players[i].uid, p.desk.players[i].name, p.desk.players[i].attackTime, p.desk.players[i].moveTime)
			for _, v2 := range p.desk.players[i].partnerSlots {
				if v2.PartnerUid == 0 {
					continue
				}
				if par := p.desk.players[i].getPartnerInfoBySlotId(v2.SlotId); par != nil {
					par.moveTime = actionTime
					par.attackTime = actionTime

					//logger.Debugf("伙伴 %d:%s  攻击时间 %v  移动时间 %v ", par.uid, par.name, par.attackTime, par.moveTime)
				}
			}
			if !p.desk.players[i].isRobot() {
				if matchInfo.mtype == pb.MatchType_SINGLE_COMBAT {
					p.desk.players[i].PVPInfo.PVPDuel.dailyMatchTimes++
					p.desk.players[i].tasks.dispatch(pb.TargetType_TOTAL_CHALLENGES_IN_ARENA, 0, 1)
				} else if matchInfo.mtype == pb.MatchType_FREE_FOR_ALL {
					p.desk.players[i].PVPInfo.PVPMelee.dailyMatchTimes++
					p.desk.players[i].tasks.dispatch(pb.TargetType_TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE, 0, 1)
				}

			}
		}
		for i := range p.desk.robots {
			actionTime := time.Now().Add(3100 * time.Millisecond)
			p.desk.robots[i].moveTime = actionTime
			p.desk.robots[i].attackTime = actionTime
			//logger.Debugf("主角 %d:%s  攻击时间 %v  移动时间 %v ", p.desk.players[i].uid, p.desk.players[i].name, p.desk.players[i].attackTime, p.desk.players[i].moveTime)
			for _, v2 := range p.desk.robots[i].partnerSlots {
				if v2.PartnerUid == 0 {
					continue
				}
				if par := p.desk.robots[i].getPartnerInfoBySlotId(v2.SlotId); par != nil {
					par.moveTime = actionTime
					par.attackTime = actionTime

					//logger.Debugf("伙伴 %d:%s  攻击时间 %v  移动时间 %v ", par.uid, par.name, par.attackTime, par.moveTime)
				}
			}
		}
		p.desk.isInitMon = true
		p.desk.start()
	}

	return nil
}

func (mm *MatchManager) PVPInfo(s *session.Session, req *pb.PVPInfoRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	res := &pb.PVPInfoResponse{MatchType: req.MatchType}
	if req.MatchType == pb.MatchType_SINGLE_COMBAT {
		res.DailyMatchTimes = p.PVPDuel.dailyMatchTimes
		res.DailyBuyNum = p.PVPDuel.dailyBuyNum
	} else if req.MatchType == pb.MatchType_FREE_FOR_ALL {
		res.DailyMatchTimes = p.PVPMelee.dailyMatchTimes
		res.DailyBuyNum = p.PVPMelee.dailyBuyNum
	}
	return s.Response(res)
}

// 1v1测试
func (mm *MatchManager) combatMatch(red *Player, blue *Player) int64 {
	// mm.lock.Lock()
	// defer mm.lock.Unlock()

	// player1 := &MatchPlayer{Player: red}
	// player2 := &MatchPlayer{Player: blue}

	// player1.bornIndex = 0
	// player2.bornIndex = 1

	// matchid := snowflake.GenID()
	// matchInfo := &MatchInfo{matchId: matchid, mplayers: []*MatchPlayer{}, startTime: time.Now()}
	// matchInfo.mplayers = append(matchInfo.mplayers, player1)
	// matchInfo.mplayers = append(matchInfo.mplayers, player2)
	// matchInfo.mtype = pb.MatchType_SINGLE_COMBAT

	// mm.fightingList[matchid] = matchInfo

	// if room := roomManager.getRoomByType(constant.PVPOneOnOneMap, 1); room != nil {
	// 	for _, d := range room.desks {
	// 		if d.status() == constant.DeskStatusDestory {
	// 			continue
	// 		}
	// 		if len(d.players) > 0 {
	// 			continue
	// 		}
	// 		if d.mapInfo == nil {
	// 			d.setStatus(constant.DeskStatusDestory)
	// 			continue
	// 		}
	// 		mm.fightingList[matchid].roomId = d.deskID
	// 		mm.fightingList[matchid].Desk = d
	// 		d.matchUid = matchid
	// 		d.MatchInfo = mm.fightingList[matchid]

	// 		if err := d.playerJoin(player1.Player, false, 0); err != nil {
	// 			logger.Debugf(" combatMatch 加入房间失败 err=%s", err.Error())
	// 		}
	// 		if err := d.playerJoin(player2.session, false, 0); err != nil {
	// 			logger.Debugf(" combatMatch 加入房间失败 err=%s", err.Error())
	// 		}

	// 		d.checkStart()

	// 		return matchid
	// 	}
	// }
	return 0
}

// 前5后5十个对手
func (mm *MatchManager) PvpMatching(s *session.Session, req *pb.PvpMatchingRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	res := &pb.PvpMatchingResponse{Code: 0, MatchType: req.MatchType, Players: []*pb.PvpMatchingResponse_PVPPlayer{}}

	modeleOpen := constant.ModuleID(0)
	var result []*RobotInfo
	if req.MatchType == 1 {
		modeleOpen = constant.ModulePK1v1
		result = robotMgr.getTenClosestPlayers(constant.RankTypeSoul, p.serverID, p.uid)
	} else if req.MatchType == 2 {
		modeleOpen = constant.ModulePKMulti
		result = robotMgr.getTenClosestPlayers(constant.RankTypeTurnBased, p.serverID, p.uid)
	}

	if !p.isModuleOpen(modeleOpen) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}
	if len(result) == 0 {

		logger.Debugf("[%d:%s] 匹配对手人数异常  数量 %d, 模式 %v  ", p.uid, p.name, len(result), req.MatchType)
		res.Code = -1
		res.Error = errutil.ServerIsBusy
		return s.Response(res)
	}

	for _, v := range result {
		if v.robotType == RobotTypePlayer {
			v.refreshStats()
			if realPlayer, ok := defaultManager.player(v.uid); ok {
				logger.Debugf("匹配到排行榜上玩家离线数据[%d:%s]  战力 %d , 真实战力 %d", realPlayer.uid, realPlayer.name, v.totalCombat, realPlayer.totalCombat)
				v.totalCombat = realPlayer.totalCombat
			}
		}

		if v.uid == p.uid {
			logger.Warningf("[%d] 匹配对手时发现自己 对战类型 %d ", p.uid, req.MatchType)
			continue
		}

		score := v.pvpScore
		if req.MatchType == 2 {
			score = v.gvgScore
		}
		AvatarV2 := int32(0)
		Frame := int32(0)
		if v.avatar != nil {
			AvatarV2 = v.avatar.GetCurrentAvatar()
			Frame = v.avatar.GetCurrentFrame()
		}
		//获取用户头像和头像框
		res.Players = append(res.Players, &pb.PvpMatchingResponse_PVPPlayer{
			Uid:      v.uid,
			Name:     v.name,
			Level:    v.level,
			Combat:   v.totalCombat,
			Score:    score,
			Rank:     int32(v.rank),
			Job:      int32(v.job),
			Gender:   int32(v.sex),
			Avatar:   "avatar",
			AvatarV2: AvatarV2,
			Frame:    Frame,
		})
	}
	mm.upPvpChallengeListDate(p, req.MatchType)
	logger.Debugf("[%d] 匹配对手人数 %d ", p.uid, len(res.Players))

	return s.Response(res)
}

//PVPChallengeRequest

func (mm *MatchManager) upPvpChallengeListDate(p *Player, mode int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//查询是否是首次挑战
	getUpRecord := p.getUpRecord("challenge_list" + fmt.Sprintf("%d", mode))
	is_first_time := true
	if getUpRecord > 0 {
		is_first_time = false
	}
	shushuDate["is_first_time"] = is_first_time
	shushuDate["name"] = "challenge_list"
	switch mode {
	case 1:
		shushuDate["mode"] = "单人挑战"
	case 2:
		shushuDate["mode"] = "回合制"
	default:
		shushuDate["mode"] = "单人挑战"
	}
	p.addUpRecord("challenge_list" + fmt.Sprintf("%d", mode))
	shushumain.ReportData(shushuDate)
}
func (mm *MatchManager) PVPChallenge(s *session.Session, req *pb.PVPChallengeRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	logger.Debugf("[%d:%s] 即时赛挑战  对手 %d ", p.uid, p.name, req.Uid)
	res := &pb.PVPChallengeResponse{}

	if p.desk != nil {
		res.Code = -1
		return s.Response(res)
	}

	if p.uid == req.Uid {
		res.Code = -4
		res.Error = errutil.ErrIllegalParameter.Error()
		logger.Debugf("[%d:%s] 不能挑战自己 ", p.uid, p.name)
		return s.Response(res)
	}

	if p.PVPInfo.PVPDuel.dailyMatchTimes >= int32(global.GlobalCfg.PVP.DailySet.DailyTime+int(p.PVPInfo.PVPDuel.dailyBuyNum)) {
		res.Code = -5
		res.Error = errutil.LimitHasBeenReached
		return s.Response(res)
	}

	robot := robotMgr.getRankRobotInfoByUid(constant.RankTypeSoul, p.serverID, req.Uid)
	if robot == nil {
		res.Code = -2
		logger.Debugf("[%d:%s] 挑战对手排行榜数据不存在, 对手UID %d ", p.uid, p.name, req.Uid)
		return s.Response(res)
	}

	if robot.status == StatusPlaying {
		res.Code = -3
		res.Error = errutil.TargetInBattle

		logger.Debugf("[%d:%s] 挑战对手正在战斗, 对手 %d : %s ", p.uid, p.name, robot.uid, robot.name)
		return s.Response(res)
	}

	blue := &MatchPlayer{Player: p, startTime: time.Now(), ready: false, win: false}
	red := &MatchPlayer{Player: robot.Player, startTime: time.Now(), ready: false, win: false}
	blue.refreshStats(WithRefreshStatsPushMessage(false))
	robot.refreshStats()

	blue.bornIndex = 0
	red.bornIndex = 1

	matchid := snowflake.GenID()
	matchInfo := &MatchInfo{matchId: matchid, mplayers: []*MatchPlayer{}, startTime: time.Now()}
	matchInfo.mplayers = append(matchInfo.mplayers, blue)
	matchInfo.mplayers = append(matchInfo.mplayers, red)
	matchInfo.mtype = pb.MatchType_SINGLE_COMBAT

	mm.fightingList[matchid] = matchInfo

	if room := roomManager.getRoomByType(constant.PVPOneOnOneMap, 1); room != nil {
		for _, d := range room.desks {
			if d.status() == constant.DeskStatusDestory {
				continue
			}
			if len(d.players) > 0 || len(d.robots) > 0 {
				continue
			}
			if d.mapInfo == nil {
				d.setStatus(constant.DeskStatusDestory)
				continue
			}
			d.reset()
			mm.fightingList[matchid].roomId = d.deskID
			mm.fightingList[matchid].Desk = d
			d.matchUid = matchid
			d.MatchInfo = mm.fightingList[matchid]

			if err := d.playerJoin(blue.Player, false, 0); err == nil {
				d.join(p)
				packet := &pb.MatchRoomJoinPush{RoomCfgId: d.cfgID, RoomUid: d.deskID, Location: blue.location, Direction: blue.direction}
				blue.pushMsg(protocol.RouteMatchRoomJoinPush, packet)

				logger.Debugf("1v1自己属性:  Uid %d  Name:%s 等级:%d 战斗力 %d 血量 %d  物理攻击 %d  魔法攻击 %d  物防 %d  魔防 %d  移动速度 %d  攻击速度 %d  命中 %d  闪避 %d   ",
					blue.uid, blue.name, blue.level, blue.totalCombat, blue.GetMaxHP(), blue.GetDC(), blue.GetMC(), blue.GetAR(), blue.GetMR(), blue.getMoveSpeed(), blue.getAttackSpeed(), blue.stat(pb.Attribute_HIT), blue.stat(pb.Attribute_DODGE))
			}
			if err := d.robotJoin(robot, 0); err == nil {
				packet := &pb.MatchRoomJoinPush{RoomCfgId: d.cfgID, RoomUid: d.deskID, Location: red.location, Direction: red.direction}
				red.pushMsg(protocol.RouteMatchRoomJoinPush, packet)

				logger.Debugf("1v1对手属性:  Uid %d  Name : %s 等级:%d 战斗力 %d 血量 %d  物理攻击 %d  魔法攻击 %d  物防 %d  魔防 %d  移动速度 %d  攻击速度 %d  命中 %d 闪避 %d   ",
					red.uid, red.name, red.level, red.totalCombat, red.GetMaxHP(), red.GetDC(), red.GetMC(), red.GetAR(), red.GetMR(), red.getMoveSpeed(), red.getAttackSpeed(), red.stat(pb.Attribute_HIT), blue.stat(pb.Attribute_DODGE))
			}

			d.syncDeskStatus()
			p.addUpRecord(constant.InstanceNum) // 增加统计
			mm.upPvpBeginShushuDate(p, red.Player, "即时赛")
			break
		}
	}

	return s.Response(res)
}

func (mm *MatchManager) upPvpBeginShushuDate(p *Player, red *Player, pvp_mode string) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//查询是否是首次挑战
	getChallengeRecord := p.getChallengeRecord(constant.PVPOneOnOneMap)
	if getChallengeRecord == 0 {
		shushuDate["is_first_try"] = true
	} else {
		shushuDate["is_first_try"] = false
	}
	// 主角信息
	//查看战斗信息
	partnerInfo := make([]map[string]interface{}, 0)
	shushuDate["name"] = "realtime_pvp_start"
	shushuDate["pvp_mode"] = pvp_mode

	// 1v1对手信息
	shushuDate["opponent_name"] = red.name
	shushuDate["opponent_id"] = fmt.Sprintf("%d", red.uid)
	shushuDate["opponent_power"] = red.totalCombat
	if red.robotType == RobotTypePlayer {
		shushuDate["is_robot"] = true
	} else {
		shushuDate["is_robot"] = false
	}
	rankInfo1 := &protocol.RankOperateRequest{
		Uid:  red.uid,
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
	shushuDate["role_rank"] = 0
	rankInfo2 := &protocol.RankOperateRequest{
		Uid:  p.uid,
		Type: 3,
		All:  3,
	}
	res2, err := provide.RankList(rankInfo2)
	if err == nil {
		if res2.RankOwn.Num > 0 {
			shushuDate["role_rank"] = res2.RankOwn.Num
		}
	}
	// shushuDate["opponent_rank"] = red.rank

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

	// 是否首次尝试
	// shushuDate["is_first_try"] = p.isFirstTry

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

func (mm *MatchManager) PvpRecoder(s *session.Session, req *pb.PvpRecorderRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	res := &pb.PvpRecorderResponse{Code: 0, Type: req.Type, Recorders: []*pb.Information{}}

	for _, v := range p.pvpRecorder {
		if v.RecordType != int8(req.Type) {
			continue
		}
		if v.Status == 2 {
			continue
		}

		one := &pb.Information{
			Uid:         v.Uid,
			SelfUid:     v.SelfUid,
			MatchTime:   v.MatchTime,
			Result:      int32(v.Result),
			OpponentUid: v.OpponentUid,
			Name:        v.Name,
			Power:       v.Power,
			Gender:      int32(v.Gender),
			Job:         int32(v.Job),
			Score:       v.Score,
			Level:       v.Level,
			AvatarV2:    v.AvatarV2,
			Frame:       v.Frame,
		}

		res.Recorders = append(res.Recorders, one)
	}

	return s.Response(res)
}

func (mm *MatchManager) TurnBasedSlot(s *session.Session, req *pb.TurnBasedSlotRequest) error {
	res := &pb.TurnBasedSlotResponse{Code: 0, Uids: []int64{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		return s.Response(res)
	}

	bFind := false
	for _, v := range req.Uids {
		if v == 0 {
			continue
		}
		if v == p.uid {
			bFind = true
		} else {
			if par := p.getPartnerInfo(v); par == nil {
				res.Code = -1
				res.Error = errutil.ErrPartnerInvalid.Error()
				return s.Response(res)
			}
		}
	}
	if !bFind {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}
	p.turnBasedSlot = [TurnBasedMax]int64{}

	copy(p.turnBasedSlot[:], req.Uids)

	for _, v := range p.turnBasedSlot {
		res.Uids = append(res.Uids, v)
	}

	if r := robotMgr.getRankRobotInfoByUid(constant.RankTypeTurnBased, 0, p.uid); r != nil {
		r.turnBasedSlot = [TurnBasedMax]int64{}
		copy(r.turnBasedSlot[:], req.Uids)
	}

	return s.Response(res)
}

func (mm *MatchManager) ElementPartnerSlot(s *session.Session, req *pb.ElementPartnerSlotRequest) error {
	res := &pb.ElementPartnerSlotResponse{Code: 0, Uids: []int64{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		return s.Response(res)
	}

	bFind := false
	for _, v := range req.Uids {
		if v == 0 {
			continue
		}
		if v == p.uid {
			bFind = true
		} else {
			if par := p.getPartnerInfo(v); par == nil {
				res.Code = -1
				res.Error = errutil.ErrPartnerInvalid.Error()
				return s.Response(res)
			}
		}
	}
	if !bFind {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}
	p.ElementInfo.elementInwork = [ElementPartnerMax]int64{}

	copy(p.ElementInfo.elementInwork[:], req.Uids)

	for _, v := range p.ElementInfo.elementInwork {
		res.Uids = append(res.Uids, v)
	}

	// if r := robotMgr.getRankRobotInfoByUid(constant.RankTypeElement, 0, p.uid); r != nil {
	// 	r.elementPartnerSlot = [ElementPartnerMax]int64{}
	// 	copy(r.elementPartnerSlot[:], req.Uids)
	// }

	return s.Response(res)
}

// ========================================================================
// 								回合制
// ========================================================================

func (mm *MatchManager) TurnBasedChallenge(s *session.Session, req *pb.TurnBasedChallengeRequest) error {
	res := &pb.TurnBasedChallengeResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModulePKMulti) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	if p.desk != nil {
		res.Code = -2
		return s.Response(res)
	}

	if p.uid == req.Uid {
		res.Code = -4
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}

	// 等级限制 次数限制
	// if p.level < global.SwitchCfg[ModuleLuanDouZhiWang] {
	// 	res.Code = -9
	// 	res.Error = errutil.ErrLevel.Error()
	// 	return s.Response(res)
	// }

	if p.PVPInfo.PVPMelee.dailyMatchTimes >= int32(global.GlobalCfg.GVG.DailySet.DailyTime+int(p.PVPInfo.PVPMelee.dailyBuyNum)) {
		res.Code = -3
		res.Error = errutil.LimitHasBeenReached
		return s.Response(res)
	}

	p.PVPInfo.PVPMelee.dailyMatchTimes++
	p.tasks.dispatch(pb.TargetType_TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE, 0, 1)

	robot := robotMgr.getRankRobotInfoByUid(constant.RankTypeTurnBased, p.serverID, req.Uid)
	if robot == nil {
		res.Code = -3
		return s.Response(res)
	}

	blue := &MatchPlayer{Player: p, startTime: time.Now(), ready: false, win: false}
	red := &MatchPlayer{Player: robot.Player, startTime: time.Now(), ready: false, win: false}
	blue.refreshStats(WithRefreshStatsPushMessage(false))
	robot.refreshStats()

	res.Info = mm.turnBasedStart(blue, red)

	//奖励
	items := []global.IdNum{}
	scoreCfg := mm.calcScore(blue.gvgScore, red.gvgScore)
	if scoreCfg == nil {
		res.Code = -4
		res.Error = errutil.ErrConfigNotFound.Error()
		return s.Response(res)
	}

	logger.Debugf(" 回合制结算 是否胜利 %v 结算分段 %v 己方 %d:%s 分数 %d   对手 %d:%s 分数 %d", blue.win, scoreCfg, blue.uid, blue.name, blue.gvgScore, red.uid, red.name, red.gvgScore)

	changeScore := int64(0)
	if res.Info.Result.Win {
		blue.win = true
		res.Info.Result.Score = scoreCfg.Win
		blue.gvgScore += scoreCfg.Win
		red.gvgScore -= scoreCfg.Win
		if red.gvgScore < 0 {
			red.gvgScore = 0
		}
		changeScore = int64(scoreCfg.Win)

		res.Info.Result.Items = []*pb.ItemInfo{}
		for _, v := range global.GlobalCfg.PVPReward.Win.Keys {
			num := v.Num
			if v.Mark == 1 {
				num = v.Num * int(p.level)
			}
			items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})

			res.Info.Result.Items = append(res.Info.Result.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
		}

	} else {
		res.Info.Result.Score = -scoreCfg.Lose
		blue.gvgScore -= scoreCfg.Lose
		red.win = true
		if blue.gvgScore < 0 {
			blue.gvgScore = 0
		}
		red.gvgScore += scoreCfg.Lose

		changeScore = int64(scoreCfg.Lose)

		res.Info.Result.Items = []*pb.ItemInfo{}
		for _, v := range global.GlobalCfg.PVPReward.Lose.Keys {
			num := v.Num
			if v.Mark == 1 {
				num = v.Num * int(p.level)
			}
			items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})

			res.Info.Result.Items = append(res.Info.Result.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
		}
	}

	blue.rankChangeTCP(constant.RankTypeTurnBased, int64(blue.gvgScore))
	red.rankChangeTCP(constant.RankTypeTurnBased, int64(red.gvgScore))
	p.addUpRecord(constant.InstanceNum) // 增加统计
	if err := blue.checkAddSourceBatch(items); err == nil {
		blue.addSourceBatch(items, constant.ItemAddPvp, WithAddSourceUnNotify(AddSourceTypeUnNotify))
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range items {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		blue.SendMail(constant.ItemAddPvp, attachment)
	}

	//添加对战记录
	mm.savePvpRecord(PvpRecordTypeHuiHe, blue, red, changeScore)
	if serverInfo, err := global.ServerListMgrInstance.Get(int(p.serverID)); err == nil {
		robotMgr.updateRank(0, serverInfo.ZoneId, constant.RankTypeTurnBased)
	}

	// //查询是否是首次挑战
	getChallengeRecord := blue.Player.getChallengeRecord(constant.PVPMultiOnMultiMap)
	is_first_try := false
	if getChallengeRecord == 0 {
		is_first_try = true
	} else {
		is_first_try = false
	}
	attachment2 := []protocol.Attachment{}
	for _, v := range items {
		attachment2 = append(attachment2, protocol.Attachment{
			UserPropsId: int64(v.ID),
			Num:         int64(v.Num),
		})
	}
	blue.Player.addChallengeRecord(constant.PVPMultiOnMultiMap)
	mm.upPvpOverShushuDate(blue.Player, red.Player, attachment2, res.Info.Result.Win, is_first_try, int64(res.Info.MatchTime), "回合赛")
	return s.Response(res)
}

// 回合制战斗开始
func (mm *MatchManager) turnBasedStart(blue *MatchPlayer, red *MatchPlayer) *pb.TurnBasedBattle {

	packet := &pb.TurnBasedBattle{
		Result:    &pb.TurnBasedBattle_TBResult{Items: []*pb.ItemInfo{}},
		MatchTime: time.Now().Unix()}

	//构造对战数据
	packet.Blue = &pb.TurnBasedBattle_TBPlayer{
		Uid:   blue.uid,
		Name:  blue.name,
		Level: int32(blue.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	//所有角色初始状态
	packet.OriginBlue = &pb.TurnBasedBattle_TBPlayer{
		Uid:   blue.uid,
		Name:  blue.name,
		Level: int32(blue.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	for pos, v := range blue.turnBasedSlot {
		if v == blue.uid {
			obj := &pb.TurnBasedBattle_TBObject{
				Type:        pb.ObjectType_PLAYER,
				Uid:         blue.uid,
				CfgID:       0,
				Level:       blue.level,
				Name:        blue.name,
				MaxHp:       blue.GetMaxHP(),
				CurHp:       blue.GetMaxHP(),
				State:       pb.TurnBasedBattle_TBObject_NORMAL,
				Buffs:       []*pb.TurnBasedBattle_TBBuff{},
				Attr:        []*pb.Stats{},
				Position:    int32(pos),
				Camp:        pb.TurnBasedBattle_BLUE_CAMP,
				Job:         int32(blue.job),
				Gender:      int32(blue.sex),
				ElementMain: blue.elementMain,
				GuildName:   blue.guildName,
				OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
			}

			blue.GetStats().Range(func(k pb.Attribute, v int32) bool {
				if k == pb.Attribute_HEALTH_PERCENTAGE ||
					k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
					k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
					k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
					k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
					return true
				}
				obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				return true
			})

			//for k, v := range blue.GetStats().Attrs {
			//	if k == pb.Attribute_HEALTH_PERCENTAGE ||
			//		k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
			//		k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
			//		k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
			//		k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
			//		continue
			//	}
			//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
			//}

			blue.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
				attr := []*pb.Stats{}
				//for k, v := range buff.BuffStats.Attrs {
				//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//}
				buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
					attr = append(attr, &pb.Stats{Attr: key, Num: value})
					return true
				})
				obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: blue.uid, Attr: attr})
				return true
			})

			//for _, v := range blue.buffs {
			//	attr := []*pb.Stats{}
			//	for k, v := range v.BuffStats.Attrs {
			//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
			//	}
			//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: blue.uid, Attr: attr})
			//}

			//obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(obj.Job)), CD: 0, AdvancedLevel: 0})
			for _, v := range blue.skills {
				if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.SType != constant.SkillTypeArtifact {
					obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
				}
			}
			sort.Slice(obj.Skills, func(i, j int) bool {
				return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
			})

			packet.Blue.Objs = append(packet.Blue.Objs, obj)
		} else {
			if par := blue.getPartnerInfo(v); par != nil {
				par.refreshStats(nil, blue.Player)
				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_PARTNER,
					Uid:         par.uid,
					CfgID:       par.cfgId,
					Name:        par.name,
					Level:       par.level,
					MaxHp:       par.GetMaxHP(),
					CurHp:       par.GetMaxHP(),
					State:       pb.TurnBasedBattle_TBObject_NORMAL,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    int32(pos),
					Camp:        pb.TurnBasedBattle_BLUE_CAMP,
					Job:         int32(par.job),
					Gender:      int32(blue.sex),
					ElementMain: blue.elementMain,
					GuildName:   blue.guildName,
					OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}

				par.GetStats().Range(func(k pb.Attribute, v int32) bool {
					if k == pb.Attribute_HEALTH_PERCENTAGE ||
						k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
						k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
						k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
						k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
						return true
					}
					obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
					return true
				})
				//
				//for k, v := range par.GetStats().Attrs {
				//	if k == pb.Attribute_HEALTH_PERCENTAGE ||
				//		k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
				//		k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
				//		k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
				//		k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
				//		continue
				//	}
				//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				//}
				//for _, v := range par.buffs {
				//	attr := []*pb.Stats{}
				//
				//	for k, v := range v.BuffStats.Attrs {
				//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//	}
				//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: par.uid, Attr: attr})
				//}

				par.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
					attr := []*pb.Stats{}
					buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
						attr = append(attr, &pb.Stats{Attr: key, Num: value})
						return true
					})
					//for k, v := range buff.BuffStats.Attrs {
					//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
					//}
					obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: par.uid, Attr: attr})

					return true
				})

				obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(par.job)), CD: 0, AdvancedLevel: 0})
				for _, v := range par.skills {
					if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.RoundMark == 1 && cfg.SType != constant.SkillTypeArtifact {
						obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
					}
				}
				sort.Slice(obj.Skills, func(i, j int) bool {
					return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
				})

				packet.Blue.Objs = append(packet.Blue.Objs, obj)
			} else {
				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_INVALID_ObjectType,
					Uid:         0,
					CfgID:       0,
					Level:       0,
					MaxHp:       0,
					CurHp:       0,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    int32(pos),
					Camp:        pb.TurnBasedBattle_BLUE_CAMP,
					Job:         0,
					Gender:      int32(blue.sex),
					ElementMain: blue.elementMain,
					GuildName:   blue.guildName,
					OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}
				packet.Blue.Objs = append(packet.Blue.Objs, obj)
			}
		}

	}

	packet.Red = &pb.TurnBasedBattle_TBPlayer{
		Uid:   red.uid,
		Name:  red.name,
		Level: int32(red.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	//所有角色初始状态
	packet.OriginRed = &pb.TurnBasedBattle_TBPlayer{
		Uid:   red.uid,
		Name:  red.name,
		Level: int32(red.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	for pos, v := range red.turnBasedSlot {
		if v == red.uid {
			obj := &pb.TurnBasedBattle_TBObject{
				Type:        pb.ObjectType_PLAYER,
				Uid:         red.uid,
				CfgID:       0,
				Level:       red.level,
				Name:        red.name,
				MaxHp:       red.GetMaxHP(),
				CurHp:       red.GetMaxHP(),
				State:       pb.TurnBasedBattle_TBObject_NORMAL,
				Buffs:       []*pb.TurnBasedBattle_TBBuff{},
				Attr:        []*pb.Stats{},
				Position:    int32(pos),
				Camp:        pb.TurnBasedBattle_RED_CAMP,
				Job:         int32(red.job),
				Gender:      int32(red.sex),
				ElementMain: red.elementMain,
				GuildName:   red.guildName,
				OfficialExp: red.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
			}

			red.GetStats().Range(func(k pb.Attribute, v int32) bool {
				if k == pb.Attribute_HEALTH_PERCENTAGE ||
					k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
					k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
					k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
					k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
					return true
				}
				obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				return true
			})
			//for k, v := range red.GetStats().Attrs {
			//	if k == pb.Attribute_HEALTH_PERCENTAGE ||
			//		k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
			//		k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
			//		k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
			//		k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
			//		continue
			//	}
			//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
			//}

			red.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
				attr := []*pb.Stats{}
				buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
					attr = append(attr, &pb.Stats{Attr: key, Num: value})
					return true
				})
				//for k, v := range buff.BuffStats.Attrs {
				//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//}
				obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: blue.uid, Attr: attr})
				return true
			})
			//for _, v := range red.buffs {
			//	attr := []*pb.Stats{}
			//	for k, v := range v.BuffStats.Attrs {
			//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
			//	}
			//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: blue.uid, Attr: attr})
			//}
			//obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(obj.Job)), CD: 0, AdvancedLevel: 0})
			for _, v := range red.skills {
				if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.SType != constant.SkillTypeArtifact {
					obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
				}
			}
			sort.Slice(obj.Skills, func(i, j int) bool {
				return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
			})
			packet.Red.Objs = append(packet.Red.Objs, obj)
		} else {
			if par := red.getPartnerInfo(v); par != nil {
				par.refreshStats(nil, red.Player)

				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_PARTNER,
					Uid:         par.uid,
					CfgID:       par.cfgId,
					Name:        par.name,
					Level:       par.level,
					MaxHp:       par.GetMaxHP(),
					CurHp:       par.GetMaxHP(),
					State:       pb.TurnBasedBattle_TBObject_NORMAL,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    int32(pos),
					Camp:        pb.TurnBasedBattle_RED_CAMP,
					Job:         int32(par.job),
					Gender:      int32(red.sex),
					ElementMain: red.elementMain,
					GuildName:   red.guildName,
					OfficialExp: red.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}

				par.GetStats().Range(func(k pb.Attribute, v int32) bool {
					if k == pb.Attribute_HEALTH_PERCENTAGE ||
						k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
						k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
						k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
						k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
						return true
					}
					obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
					return true
				})
				//for k, v := range par.GetStats().Attrs {
				//	if k == pb.Attribute_HEALTH_PERCENTAGE ||
				//		k == pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE || k == pb.Attribute_MAGIC_DEFENSE_PERCENTAGE ||
				//		k == pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE || k == pb.Attribute_MAGIC_ATTACK_PERCENTAGE ||
				//		k == pb.Attribute_DODGE_PERCENTAGE || k == pb.Attribute_HIT_PERCENTAGE ||
				//		k == pb.Attribute_CRITICAL_STRIKE_PERCENTAGE || k == pb.Attribute_TOUGHNESS_PERCENTAGE {
				//		continue
				//	}
				//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				//}

				par.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
					attr := []*pb.Stats{}
					buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
						attr = append(attr, &pb.Stats{Attr: key, Num: value})
						return true
					})
					//for k, v := range buff.BuffStats.Attrs {
					//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
					//}
					obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: red.uid, Attr: attr})
					return true
				})
				//for _, v := range par.buffs {
				//	attr := []*pb.Stats{}
				//
				//	for k, v := range v.BuffStats.Attrs {
				//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//	}
				//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: red.uid, Attr: attr})
				//}

				obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(par.cfg.Job)), CD: 0, AdvancedLevel: 0})
				for _, v := range par.skills {
					if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.SType != constant.SkillTypeArtifact {
						obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
					}
				}
				sort.Slice(obj.Skills, func(i, j int) bool {
					return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
				})
				packet.Red.Objs = append(packet.Red.Objs, obj)
			} else {
				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_INVALID_ObjectType,
					Uid:         0,
					CfgID:       0,
					Level:       0,
					MaxHp:       0,
					CurHp:       0,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    0,
					Camp:        pb.TurnBasedBattle_RED_CAMP,
					Job:         0,
					Gender:      int32(blue.sex),
					ElementMain: red.elementMain,
					GuildName:   red.guildName,
					OfficialExp: red.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}
				packet.Red.Objs = append(packet.Red.Objs, obj)
			}
		}
	}

	for _, bobj := range packet.Blue.Objs {
		packet.OriginBlue.Objs = append(packet.OriginBlue.Objs, mm.copyObject(bobj))
	}
	for _, robj := range packet.Red.Objs {
		packet.OriginRed.Objs = append(packet.OriginRed.Objs, mm.copyObject(robj))
	}

	logger.Debugf(" 回合制对战开始， 双方数据 蓝上阵数量 %d  |||  红上阵数量  %d ", len(packet.Blue.Objs), len(packet.Red.Objs))

	logger.Debugf(" 蓝详细数据 uid %d name %s  level %d  ", packet.Red.Uid, packet.Blue.Name, packet.Blue.Level)

	for k, v := range packet.Blue.Objs {
		logger.Debugf("【蓝】 第[%d]号  uid %d Name%s  cfgId %d  level %d  maxHp %d  curHp %d  camp %d  gender %d Job %d ##### Attr %v",
			k, v.Uid, v.Name, v.CfgID, v.Level, v.MaxHp, v.CurHp, v.Camp, v.Gender, v.Job, v.Attr)
	}

	logger.Debugf(" 红详细数据 uid %d name %s  level %d   ", packet.Red.Uid, packet.Red.Name, packet.Red.Level)
	for k, v := range packet.Red.Objs {
		logger.Debugf("【红】 第[%d]号  uid %d name %s cfgId %d  level %d  maxHp %d  curHp %d  camp %d  gender %d job %d ##### Attr %v",
			k, v.Uid, v.Name, v.CfgID, v.Level, v.MaxHp, v.CurHp, v.Camp, v.Gender, v.Job, v.Attr)
	}

	actionNum := int32(0)                         //出手次数
	roundId := int32(1)                           //回合数
	var acctionCamp pb.TurnBasedBattle_ActionCamp //行动方阵营
	attackOrder := mm.makeAttackOrder(packet)     //本回合出手顺序

	// 回合开始
	for roundId <= int32(global.GlobalCfg.GVG.Maxround.Num) {
		if mm.getBlueAliveNum(packet) == 0 || mm.getRedAliveNum(packet) == 0 {
			logger.Debugf(" 全部角色死亡 ，游戏结束。")
			break
		}
		actionNum++

		if len(attackOrder) == 0 {
			roundId++
			attackOrder = mm.makeAttackOrder(packet)
			if len(attackOrder) == 0 || roundId > int32(global.GlobalCfg.GVG.Maxround.Num) {
				break
			}

			logger.Debugf(" 第 [%d] 回合 出手顺序 %v ", roundId, attackOrder)
			for k, av := range attackOrder {
				attobj := mm.getNextObjByUid(packet, av)

				logger.Debugf(" 第 [%d] 个人  ID:Name %d:%s 阵营 %v  剩余血量 %d", k, attobj.Uid, attobj.Name, attobj.Camp, attobj.CurHp)
			}
		}

		curObj := mm.getNextObjByUid(packet, attackOrder[0])
		attackOrder = attackOrder[1:]
		if curObj == nil {
			break
		}

		if curObj.State == pb.TurnBasedBattle_TBObject_DEAD {
			continue
		}

		acctionCamp = curObj.Camp

		targetObj := mm.getTargetObj(packet, acctionCamp)
		if targetObj == nil {
			break
		}

		//每次轮到出手技能cd-1
		for i := 0; i < len(curObj.Skills); i++ {
			if curObj.Skills[i].CD > 0 {
				curObj.Skills[i].CD--
			}
		}

		logger.Debugf(" 第[%d]回合 , [%d]次出手, 行动方 [ %d:%s ]  首个目标 [ %d:%s ]", roundId, actionNum, curObj.Uid, curObj.Name, targetObj.Uid, targetObj.Name)

		castSkillId, skillAdvancedLevel := mm.getCastSkill(curObj, packet)

		if castSkillId != 0 {
			castSkill := mm.getSkill(curObj, castSkillId)
			if castSkill == nil {
				castSkill = &pb.TurnBasedBattle_TBSkill{Id: castSkillId, CD: 0, AdvancedLevel: 0}
				castSkillId = castSkill.Id
			}
			castSkill.CD += int32(len(curObj.Skills))

			logger.Debugf(" [%d:%s] 本次行动选中技能 %d 进阶 %d ", curObj.Uid, curObj.Name, castSkillId, skillAdvancedLevel)

			// if castSkillId == 5201 {
			// 	logger.Debugf("释放技能 ID %d   释放技能角色 %d : %s", castSkillId, curObj.Uid, curObj.Name)
			// }
		} else {
			logger.Debugf("角色被控制 [%d:%s]  本轮不行动 ", curObj.Uid, curObj.Name)
		}

		oneRound := &pb.TurnBasedBattle_Round{
			RoundId:    roundId,
			ActionNum:  actionNum,
			Uid:        curObj.Uid,
			SkillId:    castSkillId,
			ActionCamp: acctionCamp,
			Effects:    []*pb.TurnBasedBattle_Round_Effect{},
			BlueStats:  []*pb.TurnBasedBattle_TBObject{},
			RedStats:   []*pb.TurnBasedBattle_TBObject{}}

		packet.Rounds = append(packet.Rounds, oneRound)

		rate := global.RandIn(100)

		//影响范围
		effectObjects := mm.getEffectList(packet, castSkillId, oneRound.ActionCamp, curObj, targetObj)

		for _, obj := range effectObjects {
			if obj == nil {
				continue
			}
			if obj.Uid == 0 || obj.Camp == pb.TurnBasedBattle_INVALID_CAMP {
				logger.Warningf("目标无法攻击 Uid==0 %t : INVALID_CAMP %t : DEAD %t", obj.Uid == 0, obj.Camp == pb.TurnBasedBattle_INVALID_CAMP, obj.State == pb.TurnBasedBattle_TBObject_DEAD)
				continue
			}

			skillCfg, err := global.GetSkillCfg(castSkillId)
			if err != nil {
				continue
			}

			// 不同阵营才能造成伤害
			var damage, finalDamage float64
			if curObj.Camp != obj.Camp && (len(skillCfg.Effect) > 0 || castSkillId == global.GetNormalSkillId(constant.MowJob(curObj.Job))) {

				if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}

				// 无敌免疫伤害
				if mm.getBuff(obj, constant.BuffType(constant.BuffGodMode)) != nil {
					continue
				}

				// 命中-闪避
				attack_hit := mm.getAttr(curObj, pb.Attribute_HIT)
				def_dodge := mm.getAttr(obj, pb.Attribute_DODGE)

				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}}

				if attack_hit < def_dodge && global.RandIn(int(def_dodge)) > int(attack_hit) {
					effect.Miss = true
					oneRound.Effects = append(oneRound.Effects, effect)
					continue
				}

				var basePower int32
				var power, skillPower float64
				var attack_attr pb.Attribute
				// 怪物攻击取攻击力最高的属性
				basePower, attack_attr = mm.getMainAttackAttr(curObj)

				power += float64(basePower)

				//技能增伤
				for _, v := range skillCfg.Effect {
					switch pb.Attribute(v.Key) {
					case pb.Attribute_PHYSICAL_ATTACK, pb.Attribute_MAGIC_ATTACK:
						skillPower += float64(v.Value)
					case pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE, pb.Attribute_MAGIC_ATTACK_PERCENTAGE, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE:
						skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
					case pb.Attribute_DMT_FINAL_INCREASE_FIXED:
						skillPower += float64(v.Value)
					}
				}

				// 技能进阶效果
				if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skillAdvancedLevel); err == nil {
					skillPower += skillPower * float64(skillUpperlevel.DamagePer) / BaseRate
				}

				power += skillPower

				// var damage_percentage_increase int32

				// // 技能进阶效果
				// if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, castSkill.AdvancedLevel); err == nil {
				// 	damage_percentage_increase += skillUpperlevel.DamagePer
				// }

				//获取对应的防御属性
				defense := mm.getAttr(obj, pb.Attribute_PHYSICAL_DEFENSE)
				if attack_attr == pb.Attribute_MAGIC_ATTACK {
					defense = mm.getAttr(obj, pb.Attribute_MAGIC_DEFENSE)
				}

				damage = float64(power) - float64(defense)

				critical := mm.getAttr(curObj, pb.Attribute_CRITICAL_STRIKE)
				tenacity := mm.getAttr(obj, pb.Attribute_TOUGHNESS)
				criticalRate := (math.Max(0, float64(critical-tenacity))*global.GlobalCfg.BaseValue.Crit.K + float64(global.GlobalCfg.BaseValue.Crit.PreValue))

				if rate < int(criticalRate) {
					criticalDamageRate := float64(mm.getAttr(curObj, pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)) / BaseRate
					damage = damage * (1 + criticalDamageRate)
					effect.Crit = true
					// oneRound.Effects = append(oneRound.Effects, effect)
				}
				//元素属性的最终增伤
				//元素的强化 强化素质减去对手的抗性 %，属性克制
				mainBase := global.GlobalCfg.Element.MainElement.AttributeID
				main := mm.getAttr(curObj, pb.Attribute_Main_EMPOWER)
				if main > 0 {
					if constant.ElementSlot_attribute[pb.ElementSlot(main)] > 0 {
						main = int32(constant.ElementSlot_attribute[pb.ElementSlot(main)])
					}
				} else {
					main = mainBase
				}
				tarmainValue := global.GlobalCfg.Element.MainElement.AttributeID
				tarValue := mm.getAttr(targetObj, pb.Attribute(mainBase))
				for j := pb.Attribute_GOLD_EMPOWER; j <= pb.Attribute_EARTH_EMPOWER; j++ {

					if tarValue < mm.getAttr(targetObj, j) {
						tarmainValue = int32(j)
						tarValue = mm.getAttr(targetObj, j)
					}
					//最终增伤
					//改版，每个强化都会给他作对比
					// logger.Infof("部位%s 法强%d,对方抗性%d,差距%d", j, (mm.getAttr(curObj, j)), mm.getAttr(targetObj, constant.Attribute_emphasis[j]),
					// 	(mm.getAttr(curObj, j)-mm.getAttr(targetObj, constant.Attribute_emphasis[j]))*100/global.GlobalCfg.Element.Battle.Constant)
					if mm.getAttr(curObj, j)-mm.getAttr(targetObj, constant.Attribute_emphasis[j]) > 0 {
						mm.setAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, (mm.getAttr(curObj, j)-mm.getAttr(targetObj, constant.Attribute_emphasis[j]))*100/global.GlobalCfg.Element.Battle.Constant)
					}
				}

				//查看是否克制
				if constant.Attribute_know_attribute[pb.Attribute(main)] == pb.Attribute(tarmainValue) {
					mm.setAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, global.GlobalCfg.Element.Battle.RestraintRatio*100)
					// logger.Debug("有克制")
				}
				// 增伤
				var final_damage_percentage_increase float64
				final_damage_percentage_increase += float64(mm.getAttr(curObj, pb.Attribute_PVP_DAMAGE_PERCENTAGE))
				final_damage_percentage_increase += float64(mm.getAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))

				//官职压制
				if curObj.OfficialExp >= 0 && obj.OfficialExp >= 0 && CheckOfficialSuppression(curObj.OfficialExp, obj.OfficialExp, global.OfficalListCfg) == constant.SuppressionApplied {
					final_damage_percentage_increase += float64(global.GlobalCfg.BaseValue.Official.AttributeNum)
					//logger.Debugf("官职压制 攻击者: %d:%s  官职经验 %d   受击者 %d:%s 官职经验 %d ", curObj.Uid, curObj.Name, curObj.OfficialExp, obj.Uid, obj.Name, obj.OfficialExp)
				}

				damage += damage * final_damage_percentage_increase / BaseRate
				damage += float64(mm.getAttr(curObj, pb.Attribute_DMT_INCREASE_FIXED_PLAYER))
				damage += float64(mm.getAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_FIXED))

				// 减伤
				damage -= float64(mm.getAttr(obj, pb.Attribute_DMT_FINAL_REDUCE_FIXED))
				damage -= float64(mm.getAttr(obj, pb.Attribute_DMT_REDUCE_FIXED_PLAYER))
				damage -= damage * float64(mm.getAttr(obj, pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER)) / BaseRate
				damage -= damage * float64(mm.getAttr(obj, pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)) / BaseRate
				damage = math.Max(damage, float64(global.RandIn(2)+1))

				//护盾buff
				if shield := mm.getBuff(obj, constant.BuffShield); shield != nil {
					if shield.Shield > int64(damage) {
						shield.Shield -= int64(damage)
						damage = 0
					} else {
						damage -= float64(shield.Shield)
						mm.delBuff(obj, shield.Id, mm.getOriginalData(packet, obj.Uid))
					}
				}
				//如果有神圣天赋，豁免超过生命值百分比的伤害
				if buff := mm.getBuff(obj, constant.BuffDamageLimit); buff != nil {
					for _, v := range buff.Attr {
						if v.Attr == pb.Attribute_HEALTH_PERCENTAGE {
							if v.Num > 0 {
								if damage > float64(obj.MaxHp)*float64(v.Num)/BaseRate {
									damage = float64(obj.MaxHp) * float64(v.Num) / BaseRate
								}
							}
						}
					}
				}

				finalDamage = damage

				if finalDamage <= 0 {
					// finalDamage = float64(global.RandIn(2) + 1)
					finalDamage = 0
				}

				if buff := mm.getBuff(obj, constant.BuffImmune); buff != nil {
					finalDamage = 0
				}
				if buff := mm.getBuff(obj, constant.BuffPercentageMaxHealthIncrease2); buff != nil && buff.Shield > 0 {
					buff.Shield -= int64(finalDamage)
					if buff.Shield < 0 {
						buff.Shield = 0
					}
					logger.Debugf("%d:%s  强脉buff 生命值扣除 %f 剩余生命值上限 %d", obj.Uid, obj.Name, finalDamage, buff.Shield)
				}

				//吸血
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_SUCKING_BLOOD_PERCENTAGE) > 0 && curObj.MaxHp > curObj.CurHp {
					recover := finalDamage * float64(mm.getAttr(curObj, pb.Attribute_SUCKING_BLOOD_PERCENTAGE)) / BaseRate

					recover = math.Min(recover, float64(curObj.MaxHp-curObj.CurHp))
					recover = math.Max(recover, 1)
					logger.Debugf("吸血属性 %d:%s 造成伤害 %f 吸血量[%f] 生命值 最大:当前[ %d : %d ]", curObj.Uid, curObj.Name, damage, recover, curObj.MaxHp, curObj.CurHp)
					if recover > 0 {
						mm.changeHp(curObj, int32(recover))
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Recover: int32(recover), Dead: false, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
					}
				}

				if finalDamage > 0 {
					effect.Damage = int32(finalDamage)
					if mm.changeHp(obj, -int32(finalDamage)) == 0 {
						obj.State = pb.TurnBasedBattle_TBObject_DEAD
						effect.Dead = true
						oneRound.Effects = append(oneRound.Effects, effect)
						logger.Debugf("%d:%s 正常死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
						//宠物复活
						if buff := mm.getBuff(obj, constant.BuffReviveOnce); buff != nil {
							for _, v := range buff.Attr {
								if v.Attr == pb.Attribute_HEALTH_PERCENTAGE && v.Num > 0 {
									RecoverHp := float64(obj.GetMaxHp()) * float64(v.Num) / BaseRate
									effect2 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: int32(RecoverHp), Revive: true}
									obj.State = pb.TurnBasedBattle_TBObject_NORMAL
									effect2.Dead = false
									oneRound.Effects = append(oneRound.Effects, effect2)
									mm.delBuff(obj, buff.Id, mm.getOriginalData(packet, obj.Uid))
									mm.changeHp(obj, int32(RecoverHp))
									logger.Debugf(" 复活 ，最大生命值 %d,恢复生命值 %s", int32(obj.GetMaxHp()), strconv.FormatFloat(RecoverHp, 'f', 2, 64))
								}
							}
						}
						if mm.getBlueAliveNum(packet) == 0 || mm.getRedAliveNum(packet) == 0 {
							logger.Debugf(" 全部角色死亡 ，游戏结束.")
							break
						}
						continue
					}
				}

				oneRound.Effects = append(oneRound.Effects, effect)
				//反伤
				if buff := mm.getBuff(obj, constant.BuffReflectDamage); buff != nil {
					if global.RandIn(int(BaseRate)) < int(buff.Rate) {
						for _, v := range buff.Attr {
							if v.Attr == pb.Attribute_REFLECT_DAMAGE_PERCENTAGE && v.Num > 0 {
								reflectDamage := damage * float64(v.Num) / BaseRate

								//加入一条回合内攻击记录
								effect2 := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: int32(reflectDamage), Reflect: true}

								if mm.changeHp(curObj, -int32(reflectDamage)) == 0 {
									curObj.State = pb.TurnBasedBattle_TBObject_DEAD
									effect2.Dead = true
								}

								oneRound.Effects = append(oneRound.Effects, effect2)
							}
						}
					}
				}
				//胆怯
				if buff := mm.getBuff(obj, constant.BuffPercentageAttackDefenseCritDecrease); buff != nil {
					if global.RandIn(int(BaseRate)) < int(buff.Rate) {
						//给攻击者加buff
						var bufStats2 *Stats
						bufStats2 = &Stats{Attrs: make(map[pb.Attribute]int32)}
						for _, v := range buff.Attr {
							bufStats2.setStats(pb.Attribute(v.Attr), v.Num)
						}
						effect3 := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
						buff2 := mm.addBuff(curObj, 100, curObj.Uid, 100, 10000, 0, bufStats2)
						effect3.Addbuffs = append(effect3.Addbuffs, buff2.Id)
						oneRound.Effects = append(oneRound.Effects, effect3)
					}
				}

				// 眩晕
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_STUN) > 0 {
					stum := mm.getAttr(curObj, pb.Attribute_STUN)
					stum_resistance := mm.getAttr(obj, pb.Attribute_STUN_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.StunvalueEx)
					if stum >= stum_resistance || global.RandIn(int(stum_resistance)) < int(stum) {
						if cfg, err := global.GetBuffCfg(1002); err == nil {
							effect4 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
							buff2 := mm.addBuff(obj, cfg.ID, curObj.Uid, int32(global.GlobalCfg.BaseValue.Extrovalue.StunRound), 10000, 0, nil)
							effect4.Addbuffs = append(effect4.Addbuffs, buff2.Id)
							oneRound.Effects = append(oneRound.Effects, effect4)
							logger.Debugf("触发眩晕属性 %d:%s 添加眩晕buff [%d]  攻击者 %d:%s", obj.Uid, obj.Name, buff2.Id, curObj.Uid, curObj.Name)
						}
					}
				}

				// 中毒
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_POISON) > 0 {
					poison := mm.getAttr(curObj, pb.Attribute_POISON)
					poison_resistance := mm.getAttr(obj, pb.Attribute_POISON_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.PoisonvalueEx)
					if poison >= poison_resistance || global.RandIn(int(poison_resistance)) < int(poison) {
						if cfg, err := global.GetBuffCfg(1001); err == nil {
							effect5 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
							buff2 := mm.addBuff(obj, cfg.ID, curObj.Uid, int32(global.GlobalCfg.BaseValue.Extrovalue.PoisonRound), 10000, 0, nil)
							effect5.Addbuffs = append(effect5.Addbuffs, buff2.Id)
							oneRound.Effects = append(oneRound.Effects, effect5)
							logger.Debugf("触发中毒属性 %d:%s 添加中毒buff [%d]  攻击者 %d:%s", obj.Uid, obj.Name, buff2.Id, curObj.Uid, curObj.Name)
						}
					}
				}

			}

			switch skillCfg.SkillID {
			case constant.SkillTypeReborn, constant.SkillTypeLastGasp:
				if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
					obj.State = pb.TurnBasedBattle_TBObject_NORMAL
				}
				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: true}
				oneRound.Effects = append(oneRound.Effects, effect)
			case constant.SkillTypeExecute:
				if float64(obj.CurHp)/float64(obj.MaxHp) <= float64(skillCfg.SpecialEffect.Slay)/BaseRate {

					effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: obj.CurHp, Dead: true, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
					oneRound.Effects = append(oneRound.Effects, effect)

					obj.CurHp = 0
					obj.State = pb.TurnBasedBattle_TBObject_DEAD
					logger.Debugf("%d:%s 斩杀死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
				}
			case constant.SkillTypeMassPurificationBoss:
				delIds := []int32{}
				for _, v := range obj.Buffs {
					if cfg, err := global.GetBuffCfg(v.Id); err == nil {
						if cfg.Type == constant.BuffTypeDebuff {
							delIds = append(delIds, v.Id)
						}
					}
				}
				for _, v := range delIds {
					mm.delBuff(obj, v, mm.getOriginalData(packet, obj.Uid))
				}
			case constant.SkillTypeAbsorb:
				if finalDamage > 0 {
					recover := float64(finalDamage) * float64(skillCfg.SpecialEffect.DamTransfer) / BaseRate
					recover = math.Min(recover, float64(curObj.MaxHp-curObj.CurHp))
					if recover > 0 {
						mm.changeHp(curObj, int32(recover))
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Recover: int32(recover), Dead: false, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
					}
				}
			//星爆术，
			case constant.SkillTypeStellarBurst:
				//概率斩杀
				if global.RandIn(int(BaseRate)) < int(skillCfg.Skilltrigger) {
					//斩杀
					if float64(obj.CurHp)/float64(obj.MaxHp) <= float64(skillCfg.SpecialEffect.Slay)/BaseRate {
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: obj.CurHp, Dead: true, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
						obj.CurHp = 0
						obj.State = pb.TurnBasedBattle_TBObject_DEAD
						logger.Debugf("%d:%s 斩杀死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
					}
				}

			default:
			}

			//技能buff触发
			for _, tri := range skillCfg.TriggerBuff {
				bufCfg, err := global.GetBuffCfg(tri)
				if err != nil {
					continue
				}
				if bufCfg.Type == constant.BuffTypeDebuff && curObj.Camp == obj.Camp {
					continue
				}
				if bufCfg.Type == constant.BuffTypeBuff && curObj.Camp != obj.Camp {
					continue
				}

				// 筛除特殊buff
				if bufCfg.Pretype == constant.BuffDecoy {
					continue
				}
				if bufCfg.Type == constant.BuffTypeDebuff {
					// 免疫debuff
					if immuneBuff := mm.getBuff(obj, constant.BuffNoNegativeEffect); immuneBuff != nil {
						// 有几率免疫
						if global.RandIn(int(BaseRate)) < int(immuneBuff.Rate) {
							continue
						}
					}
				}
				var addSec, addRate, addDam, subCtl, shield int32
				addRate += skillCfg.Effecttrigger
				addSec += skillCfg.Buffround

				// 技能进阶效果 (回合制进阶效果结束)
				// if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, castSkill.AdvancedLevel); err == nil {
				// 	subCtl += skillUpperlevel.Control
				// 	addDam += skillUpperlevel.DamagePer

				// 	addSec += skillUpperlevel.DurationCon
				// 	addRate += skillUpperlevel.Trigger
				// 	addDistance += skillUpperlevel.Distance
				// }

				// 技能关联的buff 并且 有触发几率 才计算概率
				if bufCfg.Mark == constant.BuffTypeChainSkill && skillCfg.Skilltrigger > 0 && global.RandIn(int(BaseRate)) >= int(skillCfg.Skilltrigger+addRate) {
					continue
				}

				if skillCfg.SkillID == constant.SkillTypeMassPurificationBoss && (constant.MowJob(obj.Job) == constant.Ritualist || constant.MowJob(obj.Job) == constant.PartnerRitualist ||
					constant.MowJob(obj.Job) == constant.Warrior || constant.MowJob(obj.Job) == constant.PartnerWarrior) {
					continue
				}

				var bufStats *Stats
				if bufCfg.Mark == constant.BuffTypeChainSkill {
					bufStats = &Stats{Attrs: make(map[pb.Attribute]int32)}
					for _, v := range skillCfg.Skillbuff {
						if bufCfg.Pretype == constant.BuffShield && v.Key == int32(pb.Attribute_HEALTH_PERCENTAGE) {
							shield = v.Value
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}

					if addDam > 0 {
						bufStats.setStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, addDam)
					}
					if subCtl > 0 {
						bufStats.setStats(pb.Attribute_CONTROL_TIME_REDUCTION_PERCENTAGE, subCtl)
					}
				}

				buff := mm.addBuff(obj, bufCfg.ID, curObj.Uid, addSec, addRate, int64(shield), bufStats)

				//logger.Debugf(" 增加buff buffId %d  buff属性 %v", buff.Id, buff.Attr)
				//logger.Debugf("00000000000 buff 前 属性  %v ", obj.Attr)
				mm.refreshStats(obj, packet)
				// 部分技能增加的buff要立即生效
				if skillCfg.SkillID == constant.SkillTypeReborn {
					mm.handleBuff(oneRound, obj, mm.getOriginalData(packet, obj.Uid), true)
				}
				if skillCfg.SkillID == constant.SkillTypeStrongPulse {
					mm.handleBuff(oneRound, obj, mm.getOriginalData(packet, obj.Uid), false)
				}
				//logger.Debugf("11111111111 buff 后 属性  %v ", obj.Attr)
				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid}
				effect.Addbuffs = append(effect.Addbuffs, buff.Id)
				oneRound.Effects = append(oneRound.Effects, effect)
			}

			// 复活技能只能用一次
			if skillCfg.SkillID == constant.SkillTypeReborn {
				mm.delSkill(curObj, castSkillId)
			}

			logger.Debugf(" 一次攻击结束, RoundId %d 攻击者 %d:%s  技能 %d 本次攻击产生的所有效果 %v ", oneRound.RoundId, oneRound.Uid, curObj.Name, oneRound.SkillId, oneRound.Effects)
		}

		// 处理buff
		mm.handleBuff(oneRound, curObj, mm.getOriginalData(packet, curObj.Uid), true)
		mm.refreshStats(curObj, packet)

		for _, bobj := range packet.Blue.Objs {
			oneRound.BlueStats = append(oneRound.BlueStats, mm.copyObject(bobj))
		}
		for _, robj := range packet.Red.Objs {
			oneRound.RedStats = append(oneRound.RedStats, mm.copyObject(robj))
		}
	}

	mm.calcResult(packet)

	logger.Debugf(" 游戏结束  是否胜利 %v", packet.Result.Win)
	mm.upPvpBeginShushuDate(blue.Player, red.Player, "回合制")
	blue.tasks.dispatch(pb.TargetType_TURN_BASED_PK_TIMES, 0, 1)

	return packet
}

// ========================================================================
// 								五行副本
// ========================================================================

func (mm *MatchManager) ElementChallenge(s *session.Session, req *pb.ElementChallengeRequest) error {
	res := &pb.ElementChallengeResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleElementTrial) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	if p.desk != nil {
		res.Code = -2
		return s.Response(res)
	}

	// 等级限制 次数限制
	// if p.level < global.SwitchCfg[ModuleElementSystem] {
	// 	res.Code = -9
	// 	res.Error = errutil.ErrLevel.Error()
	// 	return s.Response(res)
	// }

	//查询是否需要扫荡
	if p.ElementMap.NeedSweep != 0 {
		res.Code = -10
		res.Error = errutil.ErrNeedSweep.Error()
		return s.Response(res)
	}
	nextBattle := p.ElementMap.BattleId + 1

	//获取怪物信息
	ElementTrials := global.GetElementTrialsCfg(nextBattle)
	if ElementTrials == nil {
		res.Code = -3
		res.Error = errutil.LimitHasBeenReached
		return s.Response(res)
	}
	//还有个等级限制
	// if p.InstanceData.ElementMap.C < nexBattle {
	//还有个等级限制
	if p.level < ElementTrials.LevelLimit {
		res.Code = -9
		res.Error = errutil.ErrLevel.Error()
		return s.Response(res)
	}
	if p.ElementMap.ChallengeTime <= 0 {
		res.Code = -3
		res.Error = errutil.NoTimes
		return s.Response(res)
	}

	// p.tasks.dispatch(pb.TargetType_TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE, 0, 1)

	monsters := make(map[int64]*MonsterInfo) // <- 加这一行

	for _, v := range ElementTrials.Monster {
		monster := global.GetMonsterCfg(v.ID)
		if monster == nil {
			res.Code = -3
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}

		m, err := monsterMgr.TakeMonster(v.ID)
		if err != nil {
			res.Code = -3
			res.Error = errutil.LimitHasBeenReached
			return s.Response(res)
		}
		//查询这个怪兽的属性
		monsters[int64(v.Position)] = m
	}
	blue := &MatchPlayer{Player: p, win: false}
	red := &MatchMonster{monsters: monsters, win: false}
	blue.refreshStats(WithRefreshStatsPushMessage(false))
	// red.refreshStats(WithRefreshStatsPushMessage(false))

	res.Info = mm.elementStart(blue, red)

	//奖励
	items := []global.IdNum{}
	if res.Info.Result.Win {
		blue.win = true
		res.Info.Result.Items = []*pb.ItemInfo{}
		for _, v := range ElementTrials.Settlement {
			num := v.Num
			items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})
			res.Info.Result.Items = append(res.Info.Result.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
		}
		for _, v := range ElementTrials.FirstPassAward {
			num := v.Num
			items = append(items, global.IdNum{ID: int32(v.ID), Num: int32(num)})
			res.Info.Result.Items = append(res.Info.Result.Items, &pb.ItemInfo{CfgId: int32(v.ID), Num: int64(num)})
		}
		p.ElementMap.BattleId = nextBattle
		blue.rankChangeTCP(constant.RankTypeElement, int64(p.ElementMap.BattleId)*100000000000+int64(p.totalCombat))
		//减去挑战次数
		p.ElementMap.ChallengeTime--
		p.addUpRecord(constant.InstanceNum) // 增加统计
	} else {
		res.Info.Result.Items = []*pb.ItemInfo{}
	}
	if err := blue.checkAddSourceBatch(items); err == nil {
		blue.addSourceBatch(items, constant.ItemAddElement, WithAddSourceUnNotify(AddSourceTypeUnNotify))
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range items {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		blue.SendMail(constant.ItemAddElement, attachment)
	}
	mm.upElementOverShushuDate(blue.Player, items, res.Info.Result.Win, len(res.Info.Rounds), nextBattle)
	p.tasks.dispatch(pb.TargetType_FIVE_ELEMENTS_INSTANCE, 0, 1)
	if res.Info != nil && res.Info.Result != nil && res.Info.Result.Win {
		p.tasks.dispatch(pb.TargetType_FIVE_ELEMENTS_INSTANCE_CLEARED_FLOOR, 0, p.ElementMap.BattleId)
	}
	return s.Response(res)
}

// 回合制战斗开始
func (mm *MatchManager) elementStart(blue *MatchPlayer, red *MatchMonster) *pb.TurnBasedBattle {
	if len(red.monsters) < 1 {
		return nil
	}
	//第一个怪物为主怪
	redMain := &MonsterInfo{}
	for _, v := range red.monsters {
		redMain = v
		break
	}
	packet := &pb.TurnBasedBattle{
		Result:    &pb.TurnBasedBattle_TBResult{Items: []*pb.ItemInfo{}},
		MatchTime: time.Now().Unix()}

	//构造对战数据
	packet.Blue = &pb.TurnBasedBattle_TBPlayer{
		Uid:   blue.uid,
		Name:  blue.name,
		Level: int32(blue.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	//所有角色初始状态
	packet.OriginBlue = &pb.TurnBasedBattle_TBPlayer{
		Uid:   blue.uid,
		Name:  blue.name,
		Level: int32(blue.level),
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	for pos, v := range blue.ElementInfo.elementInwork {
		if v == blue.uid {
			obj := &pb.TurnBasedBattle_TBObject{
				Type:        pb.ObjectType_PLAYER,
				Uid:         blue.uid,
				CfgID:       0,
				Level:       blue.level,
				Name:        blue.name,
				MaxHp:       blue.GetMaxHP(),
				CurHp:       blue.GetMaxHP(),
				State:       pb.TurnBasedBattle_TBObject_NORMAL,
				Buffs:       []*pb.TurnBasedBattle_TBBuff{},
				Attr:        []*pb.Stats{},
				Position:    int32(pos),
				Camp:        pb.TurnBasedBattle_BLUE_CAMP,
				Job:         int32(blue.job),
				Gender:      int32(blue.sex),
				ElementMain: blue.elementMain,
				GuildName:   blue.guildName,
				OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
			}

			blue.GetStats().Range(func(k pb.Attribute, v int32) bool {
				obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				return true
			})

			//for k, v := range blue.GetStats().Attrs {
			//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
			//}

			blue.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
				attr := []*pb.Stats{}
				//for k, v := range buff.BuffStats.Attrs {
				//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//}
				buff.BuffStats.Range(func(k pb.Attribute, v int32) bool {
					attr = append(attr, &pb.Stats{Attr: k, Num: v})
					return true
				})

				obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: blue.uid, Attr: attr})
				return true
			})
			//for _, v := range blue.buffs {
			//	attr := []*pb.Stats{}
			//	for k, v := range v.BuffStats.Attrs {
			//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
			//	}
			//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: blue.uid, Attr: attr})
			//}

			//obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(obj.Job)), CD: 0, AdvancedLevel: 0})
			for _, v := range blue.skills {
				if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.SType != constant.SkillTypeArtifact {
					obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
				}
			}
			sort.Slice(obj.Skills, func(i, j int) bool {
				return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
			})

			packet.Blue.Objs = append(packet.Blue.Objs, obj)
		} else {
			if par := blue.getPartnerInfo(v); par != nil {
				par.refreshStats(nil, blue.Player)
				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_PARTNER,
					Uid:         par.uid,
					CfgID:       par.cfgId,
					Name:        par.name,
					Level:       par.level,
					MaxHp:       par.GetMaxHP(),
					CurHp:       par.GetMaxHP(),
					State:       pb.TurnBasedBattle_TBObject_NORMAL,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    int32(pos),
					Camp:        pb.TurnBasedBattle_BLUE_CAMP,
					Job:         int32(par.job),
					Gender:      int32(blue.sex),
					ElementMain: blue.elementMain,
					GuildName:   blue.guildName,
					OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}

				par.GetStats().Range(func(k pb.Attribute, v int32) bool {
					obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
					return true
				})
				//for k, v := range par.GetStats().Attrs {
				//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
				//}

				par.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
					attr := []*pb.Stats{}
					//for k, v := range buff.BuffStats.Attrs {
					//	attr = append(attr, &pb.Stats{Attr: k, Num: v})
					//}
					buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
						attr = append(attr, &pb.Stats{Attr: key, Num: value})
						return true
					})

					obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: par.uid, Attr: attr})
					return true
				})
				//for _, v := range par.buffs {
				//	attr := []*pb.Stats{}
				//
				//	for k, v := range v.BuffStats.Attrs {
				//		attr = append(attr, &pb.Stats{Attr: k, Num: v})
				//	}
				//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: par.uid, Attr: attr})
				//}

				obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(par.job)), CD: 0, AdvancedLevel: 0, RoundCD: 0})
				for _, v := range par.skills {
					if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.RoundMark == 1 && cfg.SType != constant.SkillTypeArtifact {
						obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: v.AdvancedLevel, RoundCD: cfg.RoundCD})
					}
				}
				sort.Slice(obj.Skills, func(i, j int) bool {
					return obj.Skills[i].RoundCD < obj.Skills[j].RoundCD
				})
				packet.Blue.Objs = append(packet.Blue.Objs, obj)
			} else {
				obj := &pb.TurnBasedBattle_TBObject{
					Type:        pb.ObjectType_INVALID_ObjectType,
					Uid:         0,
					CfgID:       0,
					Level:       0,
					MaxHp:       0,
					CurHp:       0,
					Buffs:       []*pb.TurnBasedBattle_TBBuff{},
					Attr:        []*pb.Stats{},
					Position:    int32(pos),
					Camp:        pb.TurnBasedBattle_BLUE_CAMP,
					Job:         0,
					Gender:      int32(blue.sex),
					ElementMain: blue.elementMain,
					GuildName:   blue.guildName,
					OfficialExp: blue.tokens[global.GlobalCfg.BaseValue.Official.ItemID],
				}
				packet.Blue.Objs = append(packet.Blue.Objs, obj)
			}
		}

	}

	packet.Red = &pb.TurnBasedBattle_TBPlayer{
		Name:  redMain.name,
		Level: redMain.level,
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	//所有角色初始状态
	packet.OriginRed = &pb.TurnBasedBattle_TBPlayer{
		Uid:   redMain.uid,
		Name:  redMain.name,
		Level: redMain.level,
		Objs:  []*pb.TurnBasedBattle_TBObject{},
	}

	//记录站位
	poss := []int32{}
	for i := int32(1); i <= ElementPartnerMax; i++ {
		poss = append(poss, i)
	}
	for pos, v := range red.monsters {
		obj := &pb.TurnBasedBattle_TBObject{
			Type:     pb.ObjectType_MONSTER,
			Uid:      v.uid,
			CfgID:    v.cfgId,
			Level:    v.level,
			Name:     v.name,
			MaxHp:    v.GetMaxHP(),
			CurHp:    v.GetMaxHP(),
			State:    pb.TurnBasedBattle_TBObject_NORMAL,
			Buffs:    []*pb.TurnBasedBattle_TBBuff{},
			Attr:     []*pb.Stats{},
			Position: int32(pos - 1),
			Camp:     pb.TurnBasedBattle_RED_CAMP,
			Job:      int32(v.job),
		}
		v.GetStats().Range(func(key pb.Attribute, value int32) bool {
			obj.Attr = append(obj.Attr, &pb.Stats{Attr: key, Num: value})
			return true
		})
		//for k, v := range v.GetStats().Attrs {
		//	obj.Attr = append(obj.Attr, &pb.Stats{Attr: k, Num: v})
		//}
		//for _, v := range v.buffs {
		//	obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 15, Rate: v.Rate, OwnerUid: blue.uid, Attr: []*pb.Stats{}}) //&pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 1, Rate: v.Rate, OwnerUid: v.uid, Attr: []*pb.Stats{}})
		//}
		v.buffMgr.rangeBuff(func(buff *BuffInfo) bool {
			obj.Buffs = append(obj.Buffs, &pb.TurnBasedBattle_TBBuff{Id: buff.Id, Round: 15, Rate: buff.Rate, OwnerUid: blue.uid, Attr: []*pb.Stats{}}) //&pb.TurnBasedBattle_TBBuff{Id: v.Id, Round: 1, Rate: v.Rate, OwnerUid: v.uid, Attr: []*pb.Stats{}})
			return true
		})

		// obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: global.GetNormalSkillId(constant.MowJob(v.job)), CD: 0, AdvancedLevel: 0})
		for _, v := range v.skills {
			if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.ActiveorPassive == ActiveSkill && cfg.SType != constant.SkillTypeArtifact {
				obj.Skills = append(obj.Skills, &pb.TurnBasedBattle_TBSkill{Id: v.CfgId, CD: 0, AdvancedLevel: 0})
			}
		}
		packet.Red.Objs = append(packet.Red.Objs, obj)
		//删除poss
		for i, v := range poss {
			if int64(v) == pos {
				poss = append(poss[:i], poss[i+1:]...)
			}
		}
	}
	//pos定位减去1
	for _, v := range poss {
		obj := &pb.TurnBasedBattle_TBObject{
			Type:        pb.ObjectType_INVALID_ObjectType,
			Uid:         0,
			CfgID:       0,
			Level:       0,
			MaxHp:       0,
			CurHp:       0,
			Buffs:       []*pb.TurnBasedBattle_TBBuff{},
			Attr:        []*pb.Stats{},
			Position:    v - 1,
			Camp:        pb.TurnBasedBattle_RED_CAMP,
			Job:         0,
			ElementMain: 4,
			GuildName:   "",
		}
		packet.Red.Objs = append(packet.Red.Objs, obj)
	}
	//给 packet.Red.Objs按position排序
	sort.Slice(packet.Red.Objs, func(i, j int) bool {
		return packet.Red.Objs[i].Position < packet.Red.Objs[j].Position
	})

	for _, bobj := range packet.Blue.Objs {
		packet.OriginBlue.Objs = append(packet.OriginBlue.Objs, mm.copyObject(bobj))
	}
	for _, robj := range packet.Red.Objs {
		packet.OriginRed.Objs = append(packet.OriginRed.Objs, mm.copyObject(robj))
	}

	logger.Debugf(" 回合制对战开始， 双方数据 蓝上阵数量 %d  |||  红上阵数量  %d ", len(packet.Blue.Objs), len(packet.Red.Objs))

	logger.Debugf(" 蓝详细数据 uid %d name %s  level %d  ", packet.Red.Uid, packet.Blue.Name, packet.Blue.Level)

	for k, v := range packet.Blue.Objs {
		logger.Debugf("【蓝】 第[%d]号  uid %d  cfgId %d  level %d  maxHp %d  curHp %d  camp %d  gender %d Job %d  SKill %v POS %d",
			k, v.Uid, v.CfgID, v.Level, v.MaxHp, v.CurHp, v.Camp, v.Gender, v.Job, v.Skills, v.Position)
	}

	logger.Debugf(" 红详细数据 uid %d name %s  level %d   ", packet.Red.Uid, packet.Red.Name, packet.Red.Level)
	for k, v := range packet.Red.Objs {
		logger.Debugf("【红】 第[%d]号  uid %d  cfgId %d  level %d  maxHp %d  curHp %d  camp %d  gender %d job %d SKill %v POS %d",
			k, v.Uid, v.CfgID, v.Level, v.MaxHp, v.CurHp, v.Camp, v.Gender, v.Job, v.Skills, v.Position)
	}

	actionNum := int32(0)                         //出手次数
	roundId := int32(1)                           //回合数
	var acctionCamp pb.TurnBasedBattle_ActionCamp //行动方阵营
	attackOrder := mm.makeAttackOrder(packet)     //本回合出手顺序

	// 回合开始
	for roundId <= int32(global.GlobalCfg.GVG.Maxround.Num) {
		if mm.getBlueAliveNum(packet) == 0 || mm.getRedAliveNum(packet) == 0 {
			logger.Debugf(" 全部角色死亡 ，游戏结束。")
			break
		}
		actionNum++

		if len(attackOrder) == 0 {
			roundId++
			attackOrder = mm.makeAttackOrder(packet)
			if len(attackOrder) == 0 || roundId > int32(global.GlobalCfg.GVG.Maxround.Num) {
				break
			}

			logger.Debugf(" 第 [%d] 回合 出手顺序 %v ", roundId, attackOrder)
			for k, av := range attackOrder {
				attobj := mm.getNextObjByUid(packet, av)

				logger.Debugf(" 第 [%d] 个人  ID:Name %d:%s 阵营 %v  剩余血量 %d", k, attobj.Uid, attobj.Name, attobj.Camp, attobj.CurHp)
			}
		}

		curObj := mm.getNextObjByUid(packet, attackOrder[0])
		attackOrder = attackOrder[1:]
		if curObj == nil {
			break
		}

		if curObj.State == pb.TurnBasedBattle_TBObject_DEAD {
			continue
		}

		acctionCamp = curObj.Camp

		targetObj := mm.getTargetObj(packet, acctionCamp)
		if targetObj == nil {
			break
		}

		//每次轮到出手技能cd-1
		for i := 0; i < len(curObj.Skills); i++ {
			if curObj.Skills[i].CD > 0 {
				curObj.Skills[i].CD--
			}
		}

		castSkillId, skillAdvancedLevel := mm.getCastSkill(curObj, packet)
		if castSkillId != 0 {
			castSkill := mm.getSkill(curObj, castSkillId)
			if castSkill == nil {
				castSkill = &pb.TurnBasedBattle_TBSkill{Id: castSkillId, CD: 0, AdvancedLevel: 0}
				castSkillId = castSkill.Id
			}
			castSkill.CD += int32(len(curObj.Skills))

			// if castSkillId == 6501 {
			// 	logger.Debugf("释放技能 ID %d   释放技能角色 %d : %s", castSkillId, curObj.Uid, curObj.Name)
			// }
		}

		oneRound := &pb.TurnBasedBattle_Round{
			RoundId:    roundId,
			ActionNum:  actionNum,
			Uid:        curObj.Uid,
			SkillId:    castSkillId,
			ActionCamp: acctionCamp,
			Effects:    []*pb.TurnBasedBattle_Round_Effect{},
			BlueStats:  []*pb.TurnBasedBattle_TBObject{},
			RedStats:   []*pb.TurnBasedBattle_TBObject{}}

		packet.Rounds = append(packet.Rounds, oneRound)

		rate := global.RandIn(100)

		//影响范围
		effectObjects := mm.getEffectList(packet, castSkillId, oneRound.ActionCamp, curObj, targetObj)

		logger.Debugf(" 第[%d]回合 , [%d]次出手, 行动方 [ %d:%s ] 释放技能ID [%d] 影响目标数量 [%d]  首个目标 [ %d:%s ] 目标位置[ %d]", roundId, actionNum, curObj.Uid, curObj.Name, castSkillId, len(effectObjects), targetObj.Uid, targetObj.Name, targetObj.Position)

		for _, obj := range effectObjects {
			if obj == nil {
				continue
			}
			if obj.Uid == 0 || obj.Camp == pb.TurnBasedBattle_INVALID_CAMP {
				logger.Warningf("目标无法攻击 Uid==0 %t : INVALID_CAMP %t : DEAD %t", obj.Uid == 0, obj.Camp == pb.TurnBasedBattle_INVALID_CAMP, obj.State == pb.TurnBasedBattle_TBObject_DEAD)
				continue
			}

			skillCfg, err := global.GetSkillCfg(castSkillId)
			if err != nil {
				continue
			}

			// 不同阵营才能造成伤害
			var damage, finalDamage float64
			if curObj.Camp != obj.Camp && (len(skillCfg.Effect) > 0 || castSkillId == global.GetNormalSkillId(constant.MowJob(curObj.Job))) {

				if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}

				// 无敌免疫伤害
				if mm.getBuff(obj, constant.BuffType(constant.BuffGodMode)) != nil {
					continue
				}

				// 命中-闪避
				attack_hit := mm.getAttr(curObj, pb.Attribute_HIT)
				def_dodge := mm.getAttr(obj, pb.Attribute_DODGE)

				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}}

				if attack_hit < def_dodge && global.RandIn(int(def_dodge)) > int(attack_hit) {
					effect.Miss = true
					oneRound.Effects = append(oneRound.Effects, effect)
					continue
				}

				var basePower int32
				var power, skillPower float64
				var attack_attr pb.Attribute
				// 怪物攻击取攻击力最高的属性
				basePower, attack_attr = mm.getMainAttackAttr(curObj)
				power += float64(basePower)

				//技能增伤
				for _, v := range skillCfg.Effect {
					switch pb.Attribute(v.Key) {
					case pb.Attribute_PHYSICAL_ATTACK, pb.Attribute_MAGIC_ATTACK:
						skillPower += float64(v.Value)
					case pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE, pb.Attribute_MAGIC_ATTACK_PERCENTAGE, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE:
						skillPower += float64(basePower) * (float64(v.Value) / BaseRate)
					case pb.Attribute_DMT_FINAL_INCREASE_FIXED:
						skillPower += float64(v.Value)
					}
				}
				// 技能进阶效果
				if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skillAdvancedLevel); err == nil {
					skillPower += skillPower * float64(skillUpperlevel.DamagePer) / BaseRate
				}

				power += skillPower

				//获取对应的防御属性
				defense := mm.getAttr(obj, pb.Attribute_PHYSICAL_DEFENSE)
				if attack_attr == pb.Attribute_MAGIC_ATTACK {
					defense = mm.getAttr(obj, pb.Attribute_MAGIC_DEFENSE)
				}

				damage = float64(power) - float64(defense)

				critical := mm.getAttr(curObj, pb.Attribute_CRITICAL_STRIKE)
				tenacity := mm.getAttr(obj, pb.Attribute_TOUGHNESS)
				criticalRate := (math.Max(0, float64(critical-tenacity))*global.GlobalCfg.BaseValue.Crit.K + float64(global.GlobalCfg.BaseValue.Crit.PreValue))

				if rate < int(criticalRate) {
					criticalDamageRate := float64(mm.getAttr(curObj, pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)) / BaseRate
					damage = damage * (1 + criticalDamageRate)
					effect.Crit = true
					// oneRound.Effects = append(oneRound.Effects, effect)
				}
				//元素属性的最终增伤
				//元素的强化 强化素质减去对手的抗性 %，属性克制
				mainBase := global.GlobalCfg.Element.MainElement.AttributeID
				main := mm.getAttr(curObj, pb.Attribute_Main_EMPOWER)
				if main > 0 {
					if constant.ElementSlot_attribute[pb.ElementSlot(main)] > 0 {
						main = int32(constant.ElementSlot_attribute[pb.ElementSlot(main)])
					}
				} else {
					main = mainBase
				}
				tarmainValue := global.GlobalCfg.Element.MainElement.AttributeID
				tarValue := mm.getAttr(obj, pb.Attribute(mainBase))
				for j := pb.Attribute_GOLD_EMPOWER; j <= pb.Attribute_EARTH_EMPOWER; j++ {

					if tarValue < mm.getAttr(obj, j) {
						tarmainValue = int32(j)
						tarValue = mm.getAttr(obj, j)
					}
					//最终增伤
					//改版，每个强化都会给他作对比
					// logger.Infof("部位%s 法强%d,对方抗性%d,差距%d", j, (mm.getAttr(curObj, j)), mm.getAttr(obj, constant.Attribute_emphasis[j]),
					// 	(mm.getAttr(curObj, j)-mm.getAttr(obj, constant.Attribute_emphasis[j]))*100/global.GlobalCfg.Element.Battle.Constant)
					if mm.getAttr(curObj, j)-mm.getAttr(obj, constant.Attribute_emphasis[j]) > 0 {
						mm.setAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, (mm.getAttr(curObj, j)-mm.getAttr(obj, constant.Attribute_emphasis[j]))*100/global.GlobalCfg.Element.Battle.Constant)
					}
				}

				//查看是否克制
				if constant.Attribute_know_attribute[pb.Attribute(main)] == pb.Attribute(tarmainValue) {
					mm.setAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, global.GlobalCfg.Element.Battle.RestraintRatio*100)
					// logger.Debug("有克制")
				}
				// 增伤
				var final_damage_percentage_increase float64
				final_damage_percentage_increase += float64(mm.getAttr(curObj, pb.Attribute_PVE_DAMAGE_PERCENTAGE))
				final_damage_percentage_increase += float64(mm.getAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))
				damage += damage * final_damage_percentage_increase / BaseRate
				damage += float64(mm.getAttr(curObj, pb.Attribute_DMT_INCREASE_FIXED_MONSTER))
				damage += float64(mm.getAttr(curObj, pb.Attribute_DMT_FINAL_INCREASE_FIXED))

				// 减伤
				damage -= float64(mm.getAttr(obj, pb.Attribute_DMT_FINAL_REDUCE_FIXED))
				damage -= float64(mm.getAttr(obj, pb.Attribute_DMT_REDUCE_FIXED_MONSTER))
				damage -= damage * float64(mm.getAttr(obj, pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER)) / BaseRate
				damage -= damage * float64(mm.getAttr(obj, pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)) / BaseRate
				damage = math.Max(damage, float64(global.RandIn(2)+1))

				//护盾buff
				if shield := mm.getBuff(obj, constant.BuffShield); shield != nil {
					if shield.Shield > int64(damage) {
						shield.Shield -= int64(damage)
						damage = 0
					} else {
						damage -= float64(shield.Shield)
						mm.delBuff(obj, shield.Id, mm.getOriginalData(packet, obj.Uid))
					}
				}
				//如果有神圣天赋，豁免超过生命值百分比的伤害
				if buff := mm.getBuff(obj, constant.BuffDamageLimit); buff != nil {
					for _, v := range buff.Attr {
						if v.Attr == pb.Attribute_HEALTH_PERCENTAGE {
							if v.Num > 0 {
								if damage > float64(obj.MaxHp)*float64(v.Num)/BaseRate {
									damage = float64(obj.MaxHp) * float64(v.Num) / BaseRate
								}
							}
						}
					}
				}
				finalDamage = damage
				if finalDamage <= 0 {
					// finalDamage = float64(global.RandIn(2) + 1)
					finalDamage = 0
				}

				if buff := mm.getBuff(obj, constant.BuffImmune); buff != nil {
					finalDamage = 0
				}
				if buff := mm.getBuff(obj, constant.BuffPercentageMaxHealthIncrease2); buff != nil && buff.Shield > 0 {
					buff.Shield -= int64(finalDamage)
					if buff.Shield < 0 {
						buff.Shield = 0
					}
					logger.Debugf("%d:%s  强脉buff 生命值扣除 %f 剩余生命值上限 %d", obj.Uid, obj.Name, finalDamage, buff.Shield)
				}

				//吸血
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_SUCKING_BLOOD_PERCENTAGE) > 0 && curObj.MaxHp > curObj.CurHp {
					recover := finalDamage * float64(mm.getAttr(curObj, pb.Attribute_SUCKING_BLOOD_PERCENTAGE)) / BaseRate

					recover = math.Min(recover, float64(curObj.MaxHp-curObj.CurHp))
					recover = math.Max(recover, 1)
					logger.Debugf("吸血属性 %d:%s 造成伤害 %f 吸血量[%f] 生命值 最大:当前[ %d : %d ]", curObj.Uid, curObj.Name, damage, recover, curObj.MaxHp, curObj.CurHp)
					if recover > 0 {
						mm.changeHp(curObj, int32(recover))
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Recover: int32(recover), Dead: false, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
					}
				}

				if finalDamage > 0 {
					effect.Damage = int32(finalDamage)
					if mm.changeHp(obj, -int32(finalDamage)) == 0 {
						obj.State = pb.TurnBasedBattle_TBObject_DEAD
						effect.Dead = true
						oneRound.Effects = append(oneRound.Effects, effect)
						logger.Debugf("%d:%s 正常死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
						//宠物复活
						if buff := mm.getBuff(obj, constant.BuffReviveOnce); buff != nil {
							for _, v := range buff.Attr {
								if v.Attr == pb.Attribute_HEALTH_PERCENTAGE && v.Num > 0 {
									RecoverHp := float64(obj.GetMaxHp()) * float64(v.Num) / BaseRate
									effect2 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: int32(RecoverHp), Revive: true}
									obj.State = pb.TurnBasedBattle_TBObject_NORMAL
									effect2.Dead = false
									oneRound.Effects = append(oneRound.Effects, effect2)
									mm.delBuff(obj, buff.Id, mm.getOriginalData(packet, obj.Uid))
									mm.changeHp(obj, int32(RecoverHp))
									logger.Debugf(" 复活 ，最大生命值 %d,恢复生命值 %s", int32(obj.GetMaxHp()), strconv.FormatFloat(RecoverHp, 'f', 2, 64))
								}
							}
						}
						if mm.getBlueAliveNum(packet) == 0 || mm.getRedAliveNum(packet) == 0 {
							logger.Debugf(" 全部角色死亡 ，游戏结束.")
							break
						}
						continue
					}
				}
				oneRound.Effects = append(oneRound.Effects, effect)

				//反伤
				if buff := mm.getBuff(obj, constant.BuffReflectDamage); buff != nil {
					if global.RandIn(int(BaseRate)) < int(buff.Rate) {
						for _, v := range buff.Attr {
							if v.Attr == pb.Attribute_REFLECT_DAMAGE_PERCENTAGE && v.Num > 0 {
								reflectDamage := damage * float64(v.Num) / BaseRate

								//加入一条回合内攻击记录
								effect2 := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: int32(reflectDamage), Reflect: true}

								if mm.changeHp(curObj, -int32(reflectDamage)) == 0 {
									curObj.State = pb.TurnBasedBattle_TBObject_DEAD
									effect2.Dead = true
								}

								oneRound.Effects = append(oneRound.Effects, effect2)
							}
						}
					}
				}
				//胆怯
				if buff := mm.getBuff(obj, constant.BuffPercentageAttackDefenseCritDecrease); buff != nil {
					if global.RandIn(int(BaseRate)) < int(buff.Rate) {
						//给攻击者加buff
						var bufStats2 *Stats
						bufStats2 = &Stats{Attrs: make(map[pb.Attribute]int32)}
						for _, v := range buff.Attr {
							bufStats2.setStats(pb.Attribute(v.Attr), v.Num)
						}
						effect3 := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
						buff2 := mm.addBuff(curObj, 100, curObj.Uid, 100, 10000, 0, bufStats2)
						effect3.Addbuffs = append(effect3.Addbuffs, buff2.Id)
						oneRound.Effects = append(oneRound.Effects, effect3)
					}
				}

				// 眩晕
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_STUN) > 0 {
					stum := mm.getAttr(curObj, pb.Attribute_STUN)
					stum_resistance := mm.getAttr(obj, pb.Attribute_STUN_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.StunvalueEx)
					if stum >= stum_resistance || global.RandIn(int(stum_resistance)) < int(stum) {
						if cfg, err := global.GetBuffCfg(1002); err == nil {
							effect4 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
							buff2 := mm.addBuff(obj, cfg.ID, curObj.Uid, int32(global.GlobalCfg.BaseValue.Extrovalue.StunRound), 10000, 0, nil)
							effect4.Addbuffs = append(effect4.Addbuffs, buff2.Id)
							oneRound.Effects = append(oneRound.Effects, effect4)
							logger.Debugf("触发眩晕属性 %d:%s 添加眩晕buff [%d]  攻击者 %d:%s", obj.Uid, obj.Name, buff2.Id, curObj.Uid, curObj.Name)
						}
					}
				}

				// 中毒
				if finalDamage > 0 && mm.getAttr(curObj, pb.Attribute_POISON) > 0 {
					poison := mm.getAttr(curObj, pb.Attribute_POISON)
					poison_resistance := mm.getAttr(obj, pb.Attribute_POISON_RESISTANCE) + int32(global.GlobalCfg.BaseValue.Extrovalue.PoisonvalueEx)
					if poison >= poison_resistance || global.RandIn(int(poison_resistance)) < int(poison) {
						if cfg, err := global.GetBuffCfg(1001); err == nil {
							effect5 := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Recover: 0, Revive: false}
							buff2 := mm.addBuff(obj, cfg.ID, curObj.Uid, int32(global.GlobalCfg.BaseValue.Extrovalue.PoisonRound), 10000, 0, nil)
							effect5.Addbuffs = append(effect5.Addbuffs, buff2.Id)
							oneRound.Effects = append(oneRound.Effects, effect5)
							logger.Debugf("触发中毒属性 %d:%s 添加中毒buff [%d]  攻击者 %d:%s", obj.Uid, obj.Name, buff2.Id, curObj.Uid, curObj.Name)
						}
					}
				}

			}

			switch skillCfg.SkillID {
			case constant.SkillTypeReborn, constant.SkillTypeLastGasp:
				if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
					obj.State = pb.TurnBasedBattle_TBObject_NORMAL
				}
				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: 0, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: true}
				oneRound.Effects = append(oneRound.Effects, effect)
			case constant.SkillTypeExecute:
				if float64(obj.CurHp)/float64(obj.MaxHp) <= float64(skillCfg.SpecialEffect.Slay)/BaseRate {

					effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: obj.CurHp, Dead: true, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
					oneRound.Effects = append(oneRound.Effects, effect)

					obj.CurHp = 0
					obj.State = pb.TurnBasedBattle_TBObject_DEAD
					logger.Debugf("%d:%s 斩杀死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
				}
			case constant.SkillTypeMassPurificationBoss:
				delIds := []int32{}
				for _, v := range obj.Buffs {
					if cfg, err := global.GetBuffCfg(v.Id); err == nil {
						if cfg.Type == constant.BuffTypeDebuff {
							delIds = append(delIds, v.Id)
						}
					}
				}
				for _, v := range delIds {
					mm.delBuff(obj, v, mm.getOriginalData(packet, obj.Uid))
				}
			case constant.SkillTypeAbsorb:
				if finalDamage > 0 {
					recover := float64(finalDamage) * float64(skillCfg.SpecialEffect.DamTransfer) / BaseRate
					recover = math.Min(recover, float64(curObj.MaxHp-curObj.CurHp))
					if recover > 0 {
						mm.changeHp(curObj, int32(recover))
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: curObj.Uid, Damage: 0, Recover: int32(recover), Dead: false, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
					}
				}
			//星爆术
			case constant.SkillTypeStellarBurst:
				//概率斩杀
				if global.RandIn(int(BaseRate)) < int(skillCfg.Skilltrigger) {
					//斩杀
					if float64(obj.CurHp)/float64(obj.MaxHp) <= float64(skillCfg.SpecialEffect.Slay)/BaseRate {
						effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: obj.CurHp, Dead: true, Miss: false, Crit: false, Reflect: false, Addbuffs: []int32{}, Revive: false}
						oneRound.Effects = append(oneRound.Effects, effect)
						obj.CurHp = 0
						obj.State = pb.TurnBasedBattle_TBObject_DEAD
						logger.Debugf("%d:%s 斩杀死亡!!!!!!!!!!!!", obj.Uid, obj.Name)
					}
				}

			default:
			}

			//技能buff触发
			for _, tri := range skillCfg.TriggerBuff {
				bufCfg, err := global.GetBuffCfg(tri)
				if err != nil {
					continue
				}
				if bufCfg.Type == constant.BuffTypeDebuff && curObj.Camp == obj.Camp {
					continue
				}
				if bufCfg.Type == constant.BuffTypeBuff && curObj.Camp != obj.Camp {
					continue
				}

				// 筛除特殊buff
				if bufCfg.Pretype == constant.BuffDecoy {
					continue
				}
				if bufCfg.Type == constant.BuffTypeDebuff {
					// 免疫debuff
					if immuneBuff := mm.getBuff(obj, constant.BuffNoNegativeEffect); immuneBuff != nil {
						// 有几率免疫
						if global.RandIn(int(BaseRate)) < int(immuneBuff.Rate) {
							continue
						}
					}
				}
				var addSec, addRate, addDam, subCtl, shield int32
				addRate += skillCfg.Effecttrigger
				addSec += skillCfg.Buffround

				// 技能进阶效果 (回合制进阶效果结束)
				// if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, castSkill.AdvancedLevel); err == nil {
				// 	subCtl += skillUpperlevel.Control
				// 	addDam += skillUpperlevel.DamagePer

				// 	addSec += skillUpperlevel.DurationCon
				// 	addRate += skillUpperlevel.Trigger
				// 	addDistance += skillUpperlevel.Distance
				// }

				// 技能关联的buff 并且 有触发几率 才计算概率
				if bufCfg.Mark == constant.BuffTypeChainSkill && skillCfg.Skilltrigger > 0 && global.RandIn(int(BaseRate)) >= int(skillCfg.Skilltrigger+addRate) {
					continue
				}

				if skillCfg.SkillID == constant.SkillTypeMassPurificationBoss && (constant.MowJob(obj.Job) == constant.Ritualist || constant.MowJob(obj.Job) == constant.PartnerRitualist ||
					constant.MowJob(obj.Job) == constant.Warrior || constant.MowJob(obj.Job) == constant.PartnerWarrior) {
					continue
				}

				var bufStats *Stats
				if bufCfg.Mark == constant.BuffTypeChainSkill {
					bufStats = &Stats{Attrs: make(map[pb.Attribute]int32)}
					for _, v := range skillCfg.Skillbuff {
						if bufCfg.Pretype == constant.BuffShield && v.Key == int32(pb.Attribute_HEALTH_PERCENTAGE) {
							shield = v.Value
							continue
						}
						bufStats.setStats(pb.Attribute(v.Key), v.Value)
					}

					if addDam > 0 {
						bufStats.setStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE, addDam)
					}
					if subCtl > 0 {
						bufStats.setStats(pb.Attribute_CONTROL_TIME_REDUCTION_PERCENTAGE, subCtl)
					}
				}

				buff := mm.addBuff(obj, bufCfg.ID, curObj.Uid, addSec, addRate, int64(shield), bufStats)

				//logger.Debugf(" 增加buff buffId %d  buff属性 %v", buff.Id, buff.Attr)
				//logger.Debugf("00000000000 buff 前 属性  %v ", obj.Attr)
				mm.refreshStats(obj, packet)
				// 部分技能增加的buff要立即生效
				if skillCfg.SkillID == constant.SkillTypeReborn {
					mm.handleBuff(oneRound, obj, mm.getOriginalData(packet, obj.Uid), true)
				}
				if skillCfg.SkillID == constant.SkillTypeStrongPulse {
					mm.handleBuff(oneRound, obj, mm.getOriginalData(packet, obj.Uid), false)
				}
				//logger.Debugf("11111111111 buff 后 属性  %v ", obj.Attr)
				effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid}
				effect.Addbuffs = append(effect.Addbuffs, buff.Id)
				oneRound.Effects = append(oneRound.Effects, effect)
			}

			// 复活技能只能用一次
			if skillCfg.SkillID == constant.SkillTypeReborn {
				mm.delSkill(curObj, castSkillId)
			}

			logger.Debugf(" 一次攻击结束, RoundId %d 攻击UID %d  技能 %d 本次攻击产生的所有效果  %v", oneRound.RoundId, oneRound.Uid, oneRound.SkillId, oneRound.Effects)
		}

		// 处理buff
		mm.handleBuff(oneRound, curObj, mm.getOriginalData(packet, curObj.Uid), true)
		mm.refreshStats(curObj, packet)

		for _, bobj := range packet.Blue.Objs {
			oneRound.BlueStats = append(oneRound.BlueStats, mm.copyObject(bobj))
		}
		for _, robj := range packet.Red.Objs {
			oneRound.RedStats = append(oneRound.RedStats, mm.copyObject(robj))
		}
	}

	mm.calcResultElement(packet)
	logger.Debugf(" 游戏结束  是否胜利 %v", packet.Result.Win)
	// mm.upPvpBeginShushuDate(blue.Player, red.Player, "回合制")

	return packet
}
func (mm *MatchManager) upElementOverShushuDate(p *Player, sourceInfo []global.IdNum, isSuccess bool, duration int, battleId int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 主角信息
	//查看战斗信息
	partnerInfo := make([]map[string]interface{}, 0)
	shushuDate["name"] = "wuxing_tower"
	shushuDate["main_id"] = battleId

	//查询是否是首次挑战
	getChallengeRecord := p.getChallengeRecord(constant.RoundMap)
	if getChallengeRecord == 0 {
		shushuDate["is_first_try"] = true
	} else {
		shushuDate["is_first_try"] = false
	}
	// shushuDate["opponent_rank"] = red.rank
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
	shushuDate["is_success"] = isSuccess
	shushuDate["duration"] = duration

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
	shushumain.ReportData(shushuDate)

}

// 计算结果
func (mm *MatchManager) calcResult(packet *pb.TurnBasedBattle) {
	//人死光的失败
	if mm.getRedAliveNum(packet) == 0 {
		packet.Result.Win = true
	}
	//存活人数多的获胜
	blueAliveNum := mm.getBlueAliveNum(packet)
	redAliveNum := mm.getRedAliveNum(packet)
	if blueAliveNum > 0 && redAliveNum == 0 {
		packet.Result.Win = true
	} else if blueAliveNum == 0 && redAliveNum > 0 {
		packet.Result.Win = false
	} else if blueAliveNum > redAliveNum {
		packet.Result.Win = true
	} else if blueAliveNum < redAliveNum {
		packet.Result.Win = false
	} else {
		//剩余血量高的获胜
		logger.Debugf(" 蓝色剩余血量 %d   红色剩余血量 %d", mm.getBlueCampLeftHp(packet), mm.getRedCampLeftHp(packet))
		if mm.getBlueCampLeftHp(packet) >= mm.getRedCampLeftHp(packet) {
			packet.Result.Win = true
		} else {
			packet.Result.Win = false
		}
	}
}

// 计算结果五行
func (mm *MatchManager) calcResultElement(packet *pb.TurnBasedBattle) {
	//人死光的失败
	if mm.getRedAliveNum(packet) == 0 {
		packet.Result.Win = true
	} else {
		packet.Result.Win = false

	}

}
func (mm *MatchManager) getNextObjByUid(packet *pb.TurnBasedBattle, uid int64) *pb.TurnBasedBattle_TBObject {
	for _, b := range packet.Blue.Objs {
		if b.Uid == uid {
			return b
		}
	}
	for _, r := range packet.Red.Objs {
		if r.Uid == uid {
			return r
		}
	}
	return nil
}

// 目标 交替攻击
func (mm *MatchManager) getTargetObj(packet *pb.TurnBasedBattle, acctionCamp pb.TurnBasedBattle_ActionCamp) *pb.TurnBasedBattle_TBObject {
	if acctionCamp == pb.TurnBasedBattle_BLUE_CAMP {
		for _, r := range packet.Red.Objs {
			if r.Uid == 0 || r.State == pb.TurnBasedBattle_TBObject_DEAD {
				continue
			}
			return r
		}
	} else {
		for _, b := range packet.Blue.Objs {
			if b.Uid == 0 || b.State == pb.TurnBasedBattle_TBObject_DEAD {
				continue
			}
			return b
		}
	}
	return nil
}

// 下一个要释放的技能
func (mm *MatchManager) getCastSkill(obj *pb.TurnBasedBattle_TBObject, packet *pb.TurnBasedBattle) (int32, int32) {

	// 有眩晕buff不能攻击
	if mm.getBuff(obj, constant.BuffType(constant.BuffFreeze)) != nil {
		return 0, 0
	}

	// logger.Debugf(" %d:%s 全部技能 ", obj.Uid, obj.Name)
	// for k, v := range obj.Skills {
	// 	logger.Debugf(" %d:%s [%d]技能   ID %d  CD %d ADV %d", obj.Uid, obj.Name, k, v.Id, v.CD, v.AdvancedLevel)
	// }

	for _, v := range obj.Skills {

		if v.CD == 0 {
			if cfg, err := global.GetSkillCfg(v.Id); err == nil {
				// if cfg.SkillID == global.GetNormalSkillId(constant.MowJob(obj.Job)) {
				// 	continue
				// }
				if cfg.SkillID == constant.SkillTypeReborn || cfg.SkillID == constant.SkillTypeLastGasp {
					if deadman := mm.getFirstDead(packet, obj.Camp); deadman == nil {
						continue
					}
				}
			}
			// logger.Debugf(" %d:%s 本次行动选中技能 %d  %d ", obj.Uid, obj.Name, v.Id, v.AdvancedLevel)
			return v.Id, v.AdvancedLevel
		}
	}
	//logger.Debugf(" %d:%s 本次行动普通攻击 ", obj.Uid, obj.Name)
	return global.GetNormalSkillId(constant.MowJob(obj.Job)), 0
}

// 获取技能
func (mm *MatchManager) getSkill(obj *pb.TurnBasedBattle_TBObject, skillId int32) *pb.TurnBasedBattle_TBSkill {
	for _, v := range obj.Skills {
		if v.Id == skillId {
			return v
		}
	}
	return nil
}

// 获取buff
func (mm *MatchManager) getBuff(obj *pb.TurnBasedBattle_TBObject, pretype constant.BuffType) *pb.TurnBasedBattle_TBBuff {

	for _, v := range obj.Buffs {
		if cfg, err := global.GetBuffCfg(v.Id); err == nil {
			if cfg.Pretype == int32(pretype) {
				return v
			}
		}
	}
	return nil
}

// 删除buff
func (mm *MatchManager) delBuff(obj *pb.TurnBasedBattle_TBObject, buffId int32, originalObj *pb.TurnBasedBattle_TBObject) {

	tmpbuff := []*pb.TurnBasedBattle_TBBuff{}
	var delbuff *pb.TurnBasedBattle_TBBuff
	for _, v := range obj.Buffs {
		if v.Id != buffId {
			tmpbuff = append(tmpbuff, v)
		} else {
			delbuff = v
		}
	}
	obj.Buffs = tmpbuff

	if delbuff != nil {
		cfg, err := global.GetBuffCfg(buffId)
		if err != nil {
			return
		}
		var damage int32
		switch constant.BuffType(cfg.Pretype) {
		case constant.BuffPercentageMaxHealthIncrease:
			obj.MaxHp = originalObj.MaxHp
			if obj.CurHp > obj.MaxHp {
				damage = obj.CurHp - obj.MaxHp
			}
		case constant.BuffPercentageMaxHealthIncrease2:
			damage = int32(delbuff.Shield)
			logger.Debugf("[%d:%s]移除 强脉buff 剩余护盾 %d  最大生命值 %d 当前生命值 %d", obj.Uid, obj.Name, delbuff.Shield, obj.MaxHp, obj.CurHp)
		}

		if damage > 0 {
			if mm.changeHp(obj, -int32(damage)) == 0 {
				obj.State = pb.TurnBasedBattle_TBObject_DEAD
			}
		}

		logger.Debugf("[%d:%s]移除buff BuffId %d", obj.Uid, obj.Name, delbuff.Id)
	}

}

// 添加buff
func (mm *MatchManager) addBuff(obj *pb.TurnBasedBattle_TBObject, id int32, ownerId int64, addRound int32, addRate int32, shield int64, exStats *Stats) *pb.TurnBasedBattle_TBBuff {
	for _, v := range obj.Buffs {
		if v.Id == id {
			v.Round = addRound
			v.OwnerUid = ownerId
			if shield > v.Shield {
				v.Shield = shield
			}

			logger.Debugf("[%d:%s] 刷新 buffId %d  剩余回合数 %d  剩余护盾 %d  ", obj.Uid, obj.Name, v.Id, v.Round, v.Shield)
			return v
		}
	}

	buff := &pb.TurnBasedBattle_TBBuff{Id: id, Round: addRound, OwnerUid: ownerId, Rate: addRate, Shield: shield, Attr: []*pb.Stats{}}
	cfg, err := global.GetBuffCfg(id)
	if err != nil {
		return nil
	}
	var restoreHP float64
	switch constant.BuffType(cfg.Pretype) {
	case constant.BuffPercentageDefenseDecrease,
		constant.BuffPercentageAttackDecrease,
		constant.BuffPercentageDefenseHealthDecrease,
		constant.BuffPercentageCritChanceDecrease,
		constant.BuffPercentageAccuracyDecrease:
		for _, v := range cfg.Value {
			buff.Attr = append(buff.Attr, &pb.Stats{Attr: pb.Attribute(v.Key), Num: -v.Value})
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.Attr = append(buff.Attr, &pb.Stats{Attr: k, Num: -v})
			}
		}

	case constant.BuffCommonDebuff: // 通用减益
		for _, v := range cfg.Value {
			buff.Attr = append(buff.Attr, &pb.Stats{Attr: pb.Attribute(v.Key), Num: -v.Value})
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.Attr = append(buff.Attr, &pb.Stats{Attr: k, Num: -v})
			}
		}

	case constant.BuffCommonBuff: // 通用增益
		for _, v := range cfg.Value {
			buff.Attr = append(buff.Attr, &pb.Stats{Attr: pb.Attribute(v.Key), Num: v.Value})
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.Attr = append(buff.Attr, &pb.Stats{Attr: k, Num: v})
			}
		}
	case constant.BuffPercentageMaxHealthIncrease2:
		if maxhpp, ok := exStats.Attrs[pb.Attribute_HEALTH_MAX_PERCENTAGE]; ok {
			restoreHP = float64(obj.MaxHp) * float64(maxhpp) / BaseRate
			obj.MaxHp += int32(restoreHP)
			buff.Shield = int64(restoreHP)
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.Attr = append(buff.Attr, &pb.Stats{Attr: k, Num: v})
			}
		}
		if mm.changeHp(obj, int32(restoreHP)) == 0 {
			obj.State = pb.TurnBasedBattle_TBObject_DEAD
		}
		logger.Debugf("[%d:%s]添加 强脉buff 剩余护盾 %d  最大生命值 %d 当前生命值 %d", obj.Uid, obj.Name, buff.Shield, obj.MaxHp, obj.CurHp)

	default:
		for _, v := range cfg.Value {
			buff.Attr = append(buff.Attr, &pb.Stats{Attr: pb.Attribute(v.Key), Num: v.Value})
		}
		if exStats != nil {
			for k, v := range exStats.Attrs {
				buff.Attr = append(buff.Attr, &pb.Stats{Attr: k, Num: v})
			}
		}
	}

	obj.Buffs = append(obj.Buffs, buff)

	return buff

}

// 获取影响列表
func (mm *MatchManager) getEffectList(packet *pb.TurnBasedBattle, castSkillId int32, actionCamp pb.TurnBasedBattle_ActionCamp,
	curObj *pb.TurnBasedBattle_TBObject, targetObj *pb.TurnBasedBattle_TBObject) []*pb.TurnBasedBattle_TBObject {

	effectObjects := []*pb.TurnBasedBattle_TBObject{}

	// 有眩晕buff不能攻击
	if mm.getBuff(curObj, constant.BuffType(constant.BuffFreeze)) != nil {
		return effectObjects
	}

	skillCfg, err := global.GetSkillCfg(castSkillId)
	if err != nil {
		return effectObjects
	}

	switch skillCfg.Effectnum {
	case constant.SkillEffectnumSlef:
		if curObj != nil {
			effectObjects = append(effectObjects, curObj)
		}

	case constant.SkillEffectnumAlliance:
		if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
			for _, v := range packet.Blue.Objs {
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}
				if v.Uid == 0 {
					continue
				}
				effectObjects = append(effectObjects, v)
			}
		} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
			for _, v := range packet.Red.Objs {
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}
				if v.Uid == 0 {
					continue
				}
				effectObjects = append(effectObjects, v)
			}
		}
	case constant.SkillEffectnumAllianceLowHp:
		lowHpObj := curObj
		low_hp := float64(curObj.CurHp) / float64(curObj.MaxHp)
		if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
			for _, v := range packet.Blue.Objs {
				if v.Uid == 0 || v.CurHp == 0 {
					continue
				}
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}

				tmp_low := float64(v.CurHp) / float64(v.MaxHp)
				if tmp_low < low_hp {
					low_hp = tmp_low
					lowHpObj = v
				}
			}
		} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
			for _, v := range packet.Red.Objs {
				if v.Uid == 0 || v.CurHp == 0 {
					continue
				}
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}

				tmp_low := float64(v.CurHp) / float64(v.MaxHp)
				if tmp_low < low_hp {
					low_hp = tmp_low
					lowHpObj = v
				}
			}
		}
		if lowHpObj.Uid != 0 {
			effectObjects = append(effectObjects, lowHpObj)
		}
	case constant.SkillEffectnumEnemyColumn:
		if targetObj != nil {
			effectObjects = append(effectObjects, targetObj)

			pos := targetObj.Position + 3
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}
		}
	case constant.SkillEffectnumEnemyRow:
		if targetObj != nil {
			effectObjects = append(effectObjects, targetObj)

			pos := targetObj.Position + 1
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}

			pos = targetObj.Position - 1
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}
		}
	case constant.SkillEffectnumEnemyRange:
		if targetObj != nil {
			effectObjects = append(effectObjects, targetObj)

			pos := targetObj.Position + 1
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}

			pos = targetObj.Position - 1
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}
			pos = targetObj.Position + 3
			if pos >= 0 && pos < 6 {
				if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
					if packet.Red.Objs[pos].Uid > 0 && packet.Red.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Red.Objs[pos])
					}
				} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
					if packet.Blue.Objs[pos].Uid > 0 && packet.Blue.Objs[pos].State != pb.TurnBasedBattle_TBObject_DEAD {
						effectObjects = append(effectObjects, packet.Blue.Objs[pos])
					}
				}
			}
		}
	case constant.SkillEffectnumEnemySingle:
		if targetObj != nil {
			effectObjects = append(effectObjects, targetObj)
		}
	case constant.SkillEffectnumEnemyAll:
		if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
			for _, v := range packet.Red.Objs {
				if v.Uid > 0 && v.State != pb.TurnBasedBattle_TBObject_DEAD {
					effectObjects = append(effectObjects, v)
				}
			}
		} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
			for _, v := range packet.Blue.Objs {
				if v.Uid > 0 && v.State != pb.TurnBasedBattle_TBObject_DEAD {
					effectObjects = append(effectObjects, v)
				}
			}
		}
	case constant.SkillEffectnumAllianceFirstDie:
		if obj := mm.getFirstDead(packet, curObj.Camp); obj != nil {
			effectObjects = append(effectObjects, obj)
		}
	case constant.SkillEffectnumAllianceAndGuild:
		if actionCamp == pb.TurnBasedBattle_BLUE_CAMP {
			for _, v := range packet.Blue.Objs {
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}
				if v.Uid == 0 {
					continue
				}
				effectObjects = append(effectObjects, v)
			}
		} else if actionCamp == pb.TurnBasedBattle_RED_CAMP {
			for _, v := range packet.Red.Objs {
				if v.State == pb.TurnBasedBattle_TBObject_DEAD {
					continue
				}
				if v.Uid == 0 {
					continue
				}
				effectObjects = append(effectObjects, v)
			}
		}
	}

	return effectObjects
}

// 处理buff
func (mm *MatchManager) handleBuff(round *pb.TurnBasedBattle_Round, obj *pb.TurnBasedBattle_TBObject, originalObj *pb.TurnBasedBattle_TBObject, delRound bool) {
	if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
		return
	}
	for _, v := range obj.Buffs {
		cfg, err := global.GetBuffCfg(v.Id)
		if err != nil {
			return
		}
		// if cfg.EffectTime == 0 {
		// 	continue
		// }
		var damage, recover float64
		switch constant.BuffType(cfg.Pretype) {
		case constant.BuffFixedHeal, constant.BuffPercentageHeal:
			for _, e := range v.Attr {
				if e.Attr == pb.Attribute_HEALTH {
					recover += float64(e.Num)
				}
				if e.Attr == pb.Attribute_HEALTH_PERCENTAGE {
					recover += float64(obj.MaxHp) * float64(e.Num) / BaseRate
				}
			}
		case constant.BuffFixedDamage, constant.BuffPercentageDamage, constant.BuffFixedDamage2:
			for _, e := range v.Attr {
				if e.Attr == pb.Attribute_HEALTH {
					damage += float64(e.Num)
				}
				if e.Attr == pb.Attribute_HEALTH_PERCENTAGE {
					damage += float64(obj.MaxHp) * float64(e.Num) / BaseRate
				}
			}
		case constant.BuffFixedHealOnece:
			for _, e := range v.Attr {
				if e.Attr == pb.Attribute_HEALTH {
					recover += float64(e.Num)
				}
				if e.Attr == pb.Attribute_HEALTH_PERCENTAGE {
					recover += float64(obj.MaxHp) * float64(e.Num) / BaseRate
				}
			}
		case constant.BuffPercentageMaxHealthIncrease2:
			for _, e := range v.Attr {
				if e.Attr == pb.Attribute_HEALTH {
					recover += float64(e.Num)
				}
				if e.Attr == pb.Attribute_HEALTH_PERCENTAGE {
					recover += float64(obj.MaxHp) * float64(e.Num) / BaseRate
				}
			}

		}

		effect := &pb.TurnBasedBattle_Round_Effect{TargetUid: obj.Uid, Damage: int32(damage), Recover: int32(recover)}
		if recover > 0 {
			mm.changeHp(obj, int32(recover))
		}

		if damage > 0 {
			if mm.changeHp(obj, -int32(damage)) == 0 {
				obj.State = pb.TurnBasedBattle_TBObject_DEAD
				effect.Dead = true
			}
		}
		if (effect.TargetUid != 0) && (effect.Damage != 0 || effect.Recover != 0 || effect.Dead) {
			round.Effects = append(round.Effects, effect)
		}
		if delRound {
			v.Round--
		}

		logger.Debugf("buff 生效一次 buffid %d  剩余回合数 %d  effectUid %d : %s  当前血量 %d  回血 %d  伤害 %d  死亡 %v", v.Id, v.Round, obj.Uid, obj.Name, obj.CurHp, effect.Recover, effect.Damage, effect.Dead)

		if v.Round <= 0 {
			logger.Debugf("buff 消失 buffid %d  effectUid %d : %s", v.Id, obj.Uid, obj.Name)
			mm.delBuff(obj, v.Id, originalObj)
		}

	}
}

// blue剩余活人
func (mm *MatchManager) getBlueAliveNum(packet *pb.TurnBasedBattle) int32 {
	var aliveNum int32

	for _, v := range packet.Blue.Objs {
		if v.Uid != 0 && v.State != pb.TurnBasedBattle_TBObject_DEAD {
			aliveNum++
		}
	}
	return aliveNum
}

// red剩余活人
func (mm *MatchManager) getRedAliveNum(packet *pb.TurnBasedBattle) int32 {
	var aliveNum int32
	for _, v := range packet.Red.Objs {
		if v.Uid != 0 && v.State != pb.TurnBasedBattle_TBObject_DEAD {
			aliveNum++
		}
	}
	return aliveNum
}

// 获取第一个死人
func (mm *MatchManager) getFirstDead(packet *pb.TurnBasedBattle, camp pb.TurnBasedBattle_ActionCamp) *pb.TurnBasedBattle_TBObject {
	if camp == pb.TurnBasedBattle_BLUE_CAMP {
		for _, v := range packet.Blue.Objs {
			if v.Uid == 0 {
				continue
			}
			if v.State == pb.TurnBasedBattle_TBObject_DEAD {
				return v
			}
		}
	} else if camp == pb.TurnBasedBattle_RED_CAMP {
		for _, v := range packet.Red.Objs {
			if v.Uid == 0 {
				continue
			}
			if v.State == pb.TurnBasedBattle_TBObject_DEAD {
				return v
			}
		}
	}
	return nil
}

// blue剩余血量
func (mm *MatchManager) getBlueCampLeftHp(packet *pb.TurnBasedBattle) int32 {
	var hp int32
	for _, v := range packet.Blue.Objs {
		if v.State != pb.TurnBasedBattle_TBObject_DEAD {
			hp += v.CurHp
		}
	}
	return hp
}

// red剩余血量
func (mm *MatchManager) getRedCampLeftHp(packet *pb.TurnBasedBattle) int32 {
	var hp int32
	for _, v := range packet.Red.Objs {
		if v.State != pb.TurnBasedBattle_TBObject_DEAD {
			hp += v.CurHp
		}
	}
	return hp
}

// 计算输赢分数
func (mm *MatchManager) calcScore(blueScore, redScore int32) *global.XMLIdWinLose {
	diff := redScore - blueScore

	itemLen := len(global.GlobalCfg.GVG.GVGScore.Items)

	if itemLen == 0 {
		logger.Errorf("pvp 积分配置错误")
		return nil
	}
	if diff >= 0 {
		for i := 0; i < itemLen; i++ {
			if diff >= global.GlobalCfg.GVG.GVGScore.Items[i].Id {
				return &global.GlobalCfg.GVG.GVGScore.Items[i]
			}
		}
	} else {
		for i := itemLen - 1; i >= 0; i-- {
			if diff <= global.GlobalCfg.GVG.GVGScore.Items[i].Id {
				return &global.GlobalCfg.GVG.GVGScore.Items[i]
			}
		}
	}

	return &global.GlobalCfg.GVG.GVGScore.Items[itemLen-1]
}

// 合并两个数组，去掉0元素，将较长的数组放入新数组的偶数或奇数下标位置，将较短的数组依次放入中间的空位
func (mm *MatchManager) makeAttackOrderOld(packet *pb.TurnBasedBattle) []int64 {

	var A, B []int64
	for _, v := range packet.Blue.Objs {
		if v.State == pb.TurnBasedBattle_TBObject_DEAD {
			A = append(A, 0)
			continue
		}
		A = append(A, v.Uid)
	}
	for _, v := range packet.Red.Objs {
		if v.State == pb.TurnBasedBattle_TBObject_DEAD {
			B = append(B, 0)
			continue
		}
		B = append(B, v.Uid)
	}

	// 去掉数组中的 0 元素
	nonZeroA := make([]int64, 0)
	nonZeroB := make([]int64, 0)
	for _, v := range A {
		if v != 0 {
			nonZeroA = append(nonZeroA, v)
		}
	}
	for _, v := range B {
		if v != 0 {
			nonZeroB = append(nonZeroB, v)
		}
	}

	// 确定哪个数组更长
	var longer, shorter []int64
	isALonger := false
	if len(nonZeroA) >= len(nonZeroB) {
		longer = nonZeroA
		shorter = nonZeroB
		isALonger = true
	} else {
		longer = nonZeroB
		shorter = nonZeroA
	}

	// 初始化结果数组
	result := make([]int64, len(longer)*2)

	// 将较长的数组放入新数组的偶数或奇数下标位置
	for i, v := range longer {
		if isALonger {
			result[i*2] = v
		} else {
			result[i*2+1] = v
		}
	}

	// 将较短的数组依次放入中间的空位
	shortIndex := 0
	for i := 0; i < len(result); i++ {
		if result[i] == 0 {
			result[i] = shorter[shortIndex]
			shortIndex = (shortIndex + 1) % len(shorter)
		}
	}

	return result
}
func (mm *MatchManager) makeAttackOrder(packet *pb.TurnBasedBattle) []int64 {
	// 提取并过滤非0元素的蓝队和红队UID
	nonZeroA := make([]int64, 0)
	for _, v := range packet.Blue.Objs {
		if v.State != pb.TurnBasedBattle_TBObject_DEAD && v.Uid != 0 {
			nonZeroA = append(nonZeroA, v.Uid)
		}
	}
	nonZeroB := make([]int64, 0)
	for _, v := range packet.Red.Objs {
		if v.State != pb.TurnBasedBattle_TBObject_DEAD && v.Uid != 0 {
			nonZeroB = append(nonZeroB, v.Uid)
		}
	}

	// 交替合并两个数组（蓝队先开始）
	result := make([]int64, 0, len(nonZeroA)+len(nonZeroB))
	i, j := 0, 0
	for i < len(nonZeroA) && j < len(nonZeroB) {
		result = append(result, nonZeroA[i])
		i++
		result = append(result, nonZeroB[j])
		j++
	}
	// 添加剩余元素
	result = append(result, nonZeroA[i:]...)
	result = append(result, nonZeroB[j:]...)

	return result
}

// 拷贝一个对象
func (mm *MatchManager) copyObject(obj *pb.TurnBasedBattle_TBObject) *pb.TurnBasedBattle_TBObject {
	newObj := &pb.TurnBasedBattle_TBObject{
		Uid:         obj.Uid,
		Type:        obj.Type,
		CfgID:       obj.CfgID,
		Level:       obj.Level,
		Name:        obj.Name,
		MaxHp:       obj.MaxHp,
		CurHp:       obj.CurHp,
		State:       obj.State,
		Job:         obj.Job,
		Position:    obj.Position,
		Gender:      obj.Gender,
		Camp:        obj.Camp,
		Attr:        obj.Attr,
		Skills:      obj.Skills,
		Buffs:       []*pb.TurnBasedBattle_TBBuff{},
		ElementMain: obj.ElementMain,
		GuildName:   obj.GuildName,
	}
	for _, v := range obj.Buffs {
		buff := &pb.TurnBasedBattle_TBBuff{
			Id:       v.Id,
			Round:    v.Round,
			Rate:     v.Rate,
			OwnerUid: v.OwnerUid,
			Shield:   v.Shield,
		}
		newObj.Buffs = append(newObj.Buffs, buff)
	}
	return newObj
}

// 删除技能
func (mm *MatchManager) delSkill(obj *pb.TurnBasedBattle_TBObject, skillId int32) {

	var tmpSkills []*pb.TurnBasedBattle_TBSkill

	for _, v := range obj.Skills {
		if v.Id == skillId {
			continue
		}

		tmpSkills = append(tmpSkills, v)
	}
	obj.Skills = tmpSkills
}

// 获取原始数据
func (mm *MatchManager) getOriginalData(packet *pb.TurnBasedBattle, uid int64) *pb.TurnBasedBattle_TBObject {
	for _, v := range packet.OriginBlue.Objs {
		if v.Uid == uid {
			return v
		}
	}
	for _, v := range packet.OriginRed.Objs {
		if v.Uid == uid {
			return v
		}
	}
	return nil
}

func (mm *MatchManager) getMainAttackAttr(obj *pb.TurnBasedBattle_TBObject) (int32, pb.Attribute) {
	if constant.MowJob(obj.Job) == constant.Warrior || obj.Job == constant.PartnerWarrior {
		return mm.getAttr(obj, pb.Attribute_PHYSICAL_ATTACK), pb.Attribute_PHYSICAL_ATTACK
	} else if constant.MowJob(obj.Job) == constant.Wizard || obj.Job == constant.PartnerWizard || constant.MowJob(obj.Job) == constant.Ritualist || obj.Job == constant.PartnerRitualist {
		return mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK), pb.Attribute_MAGIC_ATTACK
	} else {

		power := mm.getAttr(obj, pb.Attribute_PHYSICAL_ATTACK)
		attack_attr := pb.Attribute_PHYSICAL_ATTACK
		if mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK) > power {
			power = mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK)
			attack_attr = pb.Attribute_MAGIC_ATTACK
		}
		return power, attack_attr
	}
}

// 获取属性
func (mm *MatchManager) getAttr(obj *pb.TurnBasedBattle_TBObject, attr pb.Attribute) int32 {
	for _, v := range obj.Attr {
		if v.Attr == attr {
			return v.Num
		}
	}
	return 0
}

// 设置属性
func (mm *MatchManager) setAttr(obj *pb.TurnBasedBattle_TBObject, attr pb.Attribute, num int32) {
	for _, v := range obj.Attr {
		if v.Attr == attr {
			v.Num += num
			if v.Num < 0 {
				v.Num = 0
			}
			return
		}
	}
	obj.Attr = append(obj.Attr, &pb.Stats{Attr: attr, Num: num})
}

// 属性重新计算
func (mm *MatchManager) refreshStats(obj *pb.TurnBasedBattle_TBObject, packet *pb.TurnBasedBattle) {
	originObj := mm.getOriginalData(packet, obj.Uid)
	if originObj == nil {
		return
	}

	obj.Attr = []*pb.Stats{}
	for _, v := range originObj.Attr {
		if v.Attr == pb.Attribute_HEALTH_PERCENTAGE {
			continue
		}
		mm.setAttr(obj, v.Attr, v.Num)
	}

	for _, v := range obj.Buffs {
		//为了不要重复计算，宠物的buff需要过滤掉
		if cfg, err := global.GetBuffCfg(v.Id); err == nil {
			if cfg.Pretype == int32(constant.BuffPercentageHeal) ||
				cfg.Pretype == int32(constant.BuffPercentageAttackDefenseCritDecrease) ||
				cfg.Pretype == int32(constant.BuffFixedHeal) ||
				cfg.Pretype == int32(constant.BuffDamageLimit) ||
				cfg.Pretype == int32(constant.BuffFixedHealOnece) ||
				cfg.Pretype == int32(constant.BuffReviveOnce) ||
				cfg.Pretype == int32(constant.BuffPercentageMaxHealthIncrease) {
				continue
			}

			for _, v2 := range v.Attr {
				mm.setAttr(obj, v2.Attr, v2.Num)
			}
		}
	}

	//重新计算所有属性
	mm.setAttr(obj, pb.Attribute_PHYSICAL_ATTACK, int32(float64(mm.getAttr(obj, pb.Attribute_PHYSICAL_ATTACK))*float64(mm.getAttr(obj, pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_MAGIC_ATTACK, int32(float64(mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK))*float64(mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_PHYSICAL_DEFENSE, int32(float64(mm.getAttr(obj, pb.Attribute_PHYSICAL_DEFENSE))*float64(mm.getAttr(obj, pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_MAGIC_DEFENSE, int32(float64(mm.getAttr(obj, pb.Attribute_MAGIC_DEFENSE))*float64(mm.getAttr(obj, pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_DODGE, int32(float64(mm.getAttr(obj, pb.Attribute_DODGE))*float64(mm.getAttr(obj, pb.Attribute_DODGE_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_HIT, int32(float64(mm.getAttr(obj, pb.Attribute_HIT))*float64(mm.getAttr(obj, pb.Attribute_HIT_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_CRITICAL_STRIKE, int32(float64(mm.getAttr(obj, pb.Attribute_CRITICAL_STRIKE))*float64(mm.getAttr(obj, pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_TOUGHNESS, int32(float64(mm.getAttr(obj, pb.Attribute_TOUGHNESS))*float64(mm.getAttr(obj, pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_MAX_HEALTH, int32(float64(mm.getAttr(obj, pb.Attribute_MAX_HEALTH))*float64(mm.getAttr(obj, pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))
	mm.setAttr(obj, pb.Attribute_MAX_HEALTH, int32(float64(mm.getAttr(obj, pb.Attribute_MAX_HEALTH))*float64(mm.getAttr(obj, pb.Attribute_HEALTH_MAX_PERCENTAGE))/BaseRate))

	obj.MaxHp = mm.getAttr(obj, pb.Attribute_MAX_HEALTH)

	// logger.Debugf("[%d:%s] ##### 刷新属性后  生命 %d  最大生命 %d  物攻 %d  法攻 %d  物防 %d  法防 %d ",
	// 	obj.Uid, obj.Name, obj.CurHp, obj.MaxHp, mm.getAttr(obj, pb.Attribute_PHYSICAL_ATTACK), mm.getAttr(obj, pb.Attribute_MAGIC_ATTACK),
	// 	mm.getAttr(obj, pb.Attribute_PHYSICAL_DEFENSE), mm.getAttr(obj, pb.Attribute_MAGIC_DEFENSE))
}

// 血量改变
func (mm *MatchManager) changeHp(obj *pb.TurnBasedBattle_TBObject, num int32) int32 {
	if obj.State == pb.TurnBasedBattle_TBObject_DEAD {
		return 0
	}
	obj.CurHp += num
	if obj.CurHp < 0 {
		obj.CurHp = 0
	}
	if obj.CurHp > obj.MaxHp {
		obj.CurHp = obj.MaxHp
	}

	return obj.CurHp
}
