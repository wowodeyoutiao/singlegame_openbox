package game

import (
	"context"
	"errors"
	"fmt"
	"math"
	"nanoserver/pkg/event"
	"nanoserver/pkg/jwtv5x"
	"nanoserver/pkg/timex"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"slices"
	"strconv"
	"strings"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/randx"
	"nanoserver/pkg/sensitive"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano/scheduler"
	"github.com/pborman/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
)

const (
	kickResetBacklog = 8
)

// Prometheus metrics
var (
	onlinePlayersGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "game_online_players_total",
		Help: "The total number of online players",
	})
)

var defaultManager = NewManager()

type (
	Manager struct {
		//lock sync.Mutex
		component.Base
		//group       *nano.Group                // 广播channel
		//players     map[int64]*Player          // 所有玩家
		chKick      chan int64                 // 退出队列
		chReset     chan int64                 // 重置队列
		chRecharge  chan protocol.RechargeInfo // 充值信息
		chAddSource chan protocol.AddSource
		chGmCmd     chan protocol.GmCmd
		history     map[int32]map[int32]map[int64][]*pb.ChatHistory // 历史记录
	}
)

func NewManager() *Manager {
	return &Manager{
		//group:       nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
		//players:     map[int64]*Player{},
		chKick:      make(chan int64, kickResetBacklog),
		chReset:     make(chan int64, kickResetBacklog),
		chRecharge:  make(chan protocol.RechargeInfo, 32),
		chAddSource: make(chan protocol.AddSource, 1024),
		chGmCmd:     make(chan protocol.GmCmd, 32),
		history:     map[int32]map[int32]map[int64][]*pb.ChatHistory{},
	}
}

// 添加到聊天记录
func (m *Manager) addChatHistory(serverID int32, ptype int32, uid int64, msg string, name string, guildId int64) {
	msgg := &pb.ChatHistory{
		Uid:  uid,
		Name: name,
		Time: time.Now().Unix(),
		Msg:  msg,
	}

	serverID32 := int32(serverID)

	if _, exists := defaultManager.history[serverID32]; !exists {
		defaultManager.history[serverID32] = make(map[int32]map[int64][]*pb.ChatHistory)
	}

	if _, exists := defaultManager.history[serverID32][ptype]; !exists {
		defaultManager.history[serverID32][ptype] = make(map[int64][]*pb.ChatHistory)
	}

	if _, exists := defaultManager.history[serverID32][ptype][guildId]; !exists {
		defaultManager.history[serverID32][ptype][guildId] = make([]*pb.ChatHistory, 0)
	}

	defaultManager.history[serverID32][ptype][guildId] = append(defaultManager.history[serverID32][ptype][guildId], msgg)

	if len(defaultManager.history[serverID32][ptype][guildId]) > 30 {
		defaultManager.history[serverID32][ptype][guildId] = defaultManager.history[serverID32][ptype][guildId][1:]
	}
}

// getChatHistory 安全地获取聊天历史记录
func (m *Manager) getChatHistory(serverID int32, ptype int32, guildId int64) []*pb.ChatHistory {
	if ptype == int32(pb.ChatType_ZONE) {
		serverID = -global.ServerListMgrInstance.GetZoneId(serverID)
	}
	// 逐层检查map是否存在
	if serverMap, exists := defaultManager.history[serverID]; exists {
		if ptypeMap, exists := serverMap[ptype]; exists {
			if history, exists := ptypeMap[guildId]; exists {
				return history
			}
		}
	}
	// 如果不存在，返回空slice
	return []*pb.ChatHistory{}
}

// 查找聊天记录
func (m *Manager) ChatHistory(s *session.Session, req *pb.ChatHistoryRequest) error {
	res := &pb.ChatHistoryResponse{Code: 0}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if req.Channel == 0 {
		// 获取所有频道的历史记录
		for i := pb.ChatType_SYSTEM; i <= pb.ChatType_ZONE; i++ {
			var guildId int64
			if i == pb.ChatType_WORLD || i == pb.ChatType_ZONE {
				guildId = 0 // 世界/战区频道使用guildId=0
			} else {
				guildId = p.guildId
			}

			history := m.getChatHistory(p.serverID, int32(i), guildId)
			res.ChannelHistories = append(res.ChannelHistories, &pb.ChatHistoryResponse_ChatChannelHistory{
				Channel: int32(i),
				History: history,
			})
		}
	} else {
		// 获取指定频道的历史记录
		var guildId int64
		if req.Channel == int32(pb.ChatType_WORLD) || req.Channel == int32(pb.ChatType_ZONE) {
			guildId = 0 // 世界/战区频道使用guildId=0
		} else {
			guildId = p.guildId
		}

		history := m.getChatHistory(p.serverID, req.Channel, guildId)
		res.ChannelHistories = append(res.ChannelHistories, &pb.ChatHistoryResponse_ChatChannelHistory{
			Channel: req.Channel,
			History: history,
		})
	}

	return s.Response(res)
}

func (m *Manager) Init() {
	playerMgr.Init()
}

func (m *Manager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		//m.group.Leave(s)
		playerMgr.RemoveSessionFromGroup(s)

		// 更新在线玩家数量指标
		onlinePlayersGauge.Dec()
	})

	// 处理踢出玩家和重置玩家消息(来自http)
	scheduler.NewTimer(time.Second, func() {
	ctrl:
		for {
			select {
			case uid := <-m.chKick:
				p, ok := defaultManager.player(uid)
				if !ok {
					logger.Errorf("玩家不存在 userId: %d", uid)
					continue
				}
				if p.session == nil {
					logger.Errorf("玩家不在线 userId: %d", uid)
					continue
				}
				p.session.Close()
				logger.Infof("踢出玩家, UID=%d", uid)
			case uid := <-m.chReset:
				p, ok := defaultManager.player(uid)
				if !ok {
					return
				}
				if p.session != nil {
					logger.Errorf("玩家正在游戏中，不能重置 userId: %d", uid)
					return
				}
				p.desk = nil
				logger.Infof("重置玩家, UID=%d", uid)

				// case ri := <-m.chRecharge:
				// 	player, ok := m.player(ri.Uid)
				// 	if ok {
				// 		player.recharge4vip(ri.Coin)
				// 	} else {
				// 		async.Run(func() {
				// 			u, err := db.QueryUser(ri.Uid)
				// 			if err != nil {
				// 				logger.Errorf("玩家chRecharge同步错误,userId: %d Error=%v", ri.Uid, err)
				// 				return
				// 			}
				// 			u.Recharge += ri.Coin

				// 			if err := db.UpdateUser(u); err != nil {
				// 				logger.Errorf("保存角色信息错误, Error=%v", err)
				// 				return
				// 			}
				// 		})
				// 	}

			case ads := <-m.chAddSource:
				player, ok := m.player(ads.Uid)
				if ok {
					items := []global.IdNum{}
					for _, source := range ads.Source {
						if source.Num == 0 {
							continue
						}
						if source.Num > 0 {
							items = append(items, global.IdNum{ID: source.Id, Num: int32(source.Num)})
							if ads.ShopType > constant.ShopTypeUnknown && ads.ShopType < constant.ShopTypeAll {
								if shopType, ok := constant.ShopMap[ads.ShopType]; ok {
									player.tasks.dispatch(shopType, 0, 1)
								}
							}
						} else {
							player.takeSource(source.Id, int64(math.Abs(float64(source.Num))), ads.Method)
						}
					}
					if err := player.checkAddSourceBatch(items); err == nil {
						player.addSourceBatch(items, ads.Method)
						//player.asyncSaveDb()
					} else {
						attachment := []protocol.Attachment{}
						for _, item := range items {
							attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
						}
						player.SendMail(constant.ItemAddGame, attachment)
						logger.Errorf("Web Interface 发放道具失败,包裹空间不足,邮件发送 userId:%d  道具数量 :%d", ads.Uid, len(items))
					}
				} else {
					logger.Errorf("Web Interface 发放道具失败,玩家未找到userId:%d  SourceNum : %d", ads.Uid, len(ads.Source))
				}

			case cmd := <-m.chGmCmd:
				m.GmCmd(cmd)

			default:
				break ctrl
			}
		}
	})

	// 每分钟保存一次数据
	scheduler.NewTimer(time.Minute, func() {
		//for _, v := range m.players {
		//	if v.curState != PlayerStatusIdle {
		//		continue
		//	}
		//	if !v.bDirty {
		//		continue
		//	}
		//	v.asyncSaveDb()
		//}

		playerMgr.IteratePlayers(func(val *Player) bool {
			if val.curState != PlayerStatusIdle {
				return true
			}
			if !val.bDirty {
				return true
			}
			val.asyncSaveDb()
			return true
		})
	})
	scheduler.NewTimer(time.Second, func() {
		kickIds := []int64{}
		//for _, v := range m.players {
		//	// if v.curState >= PlayerStatusDisconnect && v.curState < PlayerStatusRobot {
		//	// 	//离线10分钟加入离线机器人
		//	// 	if time.Now().After(v.disconnectTime.Add(10 * time.Minute)) {
		//	// 		kickIds = append(kickIds, val.uid)
		//	// 	}
		//	// }
		//	if v.curState == PlayerStatusRobot {
		//		// 离线5天清理内存
		//		if time.Now().After(v.disconnectTime.Add(5 * 24 * 60 * time.Minute)) {
		//			kickIds = append(kickIds, v.uid)
		//			continue
		//		}
		//	}
		//	v.update()
		//}
		playerMgr.IteratePlayers(func(val *Player) bool {
			if val.curState == PlayerStatusRobot {
				// 离线5天清理内存
				if time.Now().After(val.disconnectTime.Add(5 * 24 * 60 * time.Minute)) {
					kickIds = append(kickIds, val.uid)
					return true
				}
			}
			val.update()
			val.eventMgr.ProcessEvents()
			return true
		})

		// 离线超过x分钟清理内存
		for _, v := range kickIds {
			m.offline(v)
		}
	})
}

func (m *Manager) TimerFor0Oclock() {
	logger.Println("触发每日0点任务")
	//for _, v := range m.players {
	//	if v.curState >= PlayerStatusDisconnect {
	//		continue
	//	}
	//	v.newDay()
	//}
	playerMgr.IteratePlayers(func(val *Player) bool {
		if val.curState >= PlayerStatusDisconnect {
			return true
		}
		val.newDay()
		return true
	})
	// 数据库清理
	// db.DeleteUserItem()
}

// 每日凌晨清楚当日在线时长，累加在线时长
func (m *Manager) DayOnlineTime() {
	logger.Println("触发每日凌晨任务")
}

func (m *Manager) player(uid int64) (*Player, bool) {
	//m.lock.Lock()
	//defer m.lock.Unlock()
	//
	//p, ok := m.players[uid]
	//return p, ok
	return playerMgr.GetPlayer(uid)
}

func (m *Manager) setPlayer(uid int64, p *Player) {
	//if _, ok := m.players[uid]; ok {
	//	log.Warnf("玩家已经存在，正在覆盖玩家， UID=%d", uid)
	//}
	//m.players[uid] = p
	playerMgr.PutPlayer(p)
}

func (m *Manager) sessionCount() int {
	//return len(m.players)
	return playerMgr.Len()
}

func (m *Manager) offline(uid int64) {
	//if _, ok := m.player(uid); ok {
	//	delete(m.players, uid)
	//	log.Infof("玩家: %d 从在线列表中删除, 剩余：%d", uid, len(m.players))
	//}
	playerMgr.RemovePlayer(uid)
}

// (弃用)
func (m *Manager) reconnect(p *Player) {
	logger.Debugf("断线重连 Uid=%d , 断线前房间ID=%d", p.uid, p.roomId)
	if p.desk != nil {
		p.desk.playerKick(p, constant.GameOverType)
		return
	}
	if room, ok := roomManager.rooms[p.roomId]; ok {
		room.ReJoin(p)
	}
}

// web interface
// 检查是否能送道具
func (m *Manager) CheckSourceAdd(source *protocol.AddSource) error {
	p, ok := m.player(source.Uid)

	if !ok {
		return errutil.ErrPlayerNotFound
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	items := []global.IdNum{}
	for _, v := range source.Source {
		items = append(items, global.IdNum{ID: v.Id, Num: int32(v.Num)})
	}
	return p.checkAddSourceBatch(items)
}

// 检查道具数量
func (m *Manager) CheckSourceNum(source *protocol.AddSource) error {
	p, ok := m.player(source.Uid)

	if !ok {
		return errutil.ErrPlayerNotFound
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, v := range source.Source {
		if v.Num <= 0 {
			logger.Warningf("Web CheckSourceNum 参数错误 ID:%d Num:%d", v.Id, v.Num)
			return errutil.ErrIllegalParameter
		}
		if !p.checkTakeSource(v.Id, v.Num) {
			return errutil.ErrItemNotEnough
		}
	}

	return nil
}

// Web添加资源接口
func (m *Manager) AddSourceBatch(source protocol.AddSource) {
	if len(source.Source) > 0 {
		m.chAddSource <- source
	}
}

// 增加特权
func (m *Manager) PrivilegeAdd(shopType int32, uid int64, privilegeId int32, expirationTime time.Time) error {
	p, ok := m.player(uid)

	if !ok {
		return errutil.ErrPlayerNotFound
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	return p.addPrivilege(shopType, privilegeId, expirationTime)
}

// Client Message Route
func (m *Manager) Login(s *session.Session, req *pb.LoginToGameServerRequest) error {
	uid := req.Uid
	s.Bind(uid)

	logger.Debugf("玩家: %d 登录: %+v", uid, req)

	res := &pb.LoginToGameServerResponse{
		Uid:        s.UID(),
		HeadUrl:    req.HeadUrl,
		Name:       req.Name,
		Timestamp:  time.Now().Unix(),
		GameStatus: 0,
	}
	if !global.LoginOpen {
		res.Code = -1
		res.Error = errutil.ServerIsBusy
		return s.Response(res)
	}

	// 校验token
	// tokenFactory := token.CreateUserFactory()
	// if check := tokenFactory.IsEffective(req.Token); !check {
	// 	logger.Debugf("Token 无效,校验失败: Uid %d 登录: %s", uid, req.Token)
	// 	res.Code = -2
	// 	res.Error = errutil.ErrInvalidToken.Error()
	// 	return s.Response(res)
	// }
	// if userInfo, err := tokenFactory.ParseToken(req.Token); err != nil {
	// 	logger.Debugf("Token 解析失败: Uid %d 登录: %s", uid, req.Token)
	// 	res.Code = -3
	// 	res.Error = errutil.ErrInvalidToken.Error()
	// 	return s.Response(res)
	// } else if userInfo.UserId != req.Uid {
	// 	logger.Debugf("Token UID不一致,校验失败: Uid %d 登录: %s", uid, req.Token)
	// 	res.Code = -4
	// 	res.Error = errutil.ErrInvalidToken.Error()
	// 	return s.Response(res)
	// }

	claims, err := jwtv5x.Instance.VerifyToken(context.TODO(), req.Token)
	if err != nil {
		logger.Debugf("Token 解析失败: Uid %d 登录: %s", uid, req.Token)
		res.Code = -3
		res.Error = errutil.ErrInvalidToken.Error()
		return s.Response(res)
	}
	if claims.UserId != strconv.FormatInt(req.Uid, 10) {
		logger.Debugf("Token UID不一致,校验失败: Uid %d 登录: %s", uid, req.Token)
		res.Code = -4
		res.Error = errutil.ErrInvalidToken.Error()
		return s.Response(res)
	}

	if p, ok := m.player(uid); !ok {

		logger.Debugf("玩家: %d 登录: 缓存无数据，加载数据  ", uid)
		var ip string
		if parts := strings.Split(s.RemoteAddr().String(), ":"); len(parts) > 0 {
			ip = parts[0]
		}

		p = newPlayer(s, uid, req.Name, req.HeadUrl, ip)
		p.bindSession(s)

		if !p.syncFromDB() {
			res.Code = -1
			res.Error = errutil.ErrDBOperation.Error()
			logger.Errorf("玩家: %d 登录失败,数据加载失败", uid)
			return s.Response(res)
		}

		// 检查封号状态（通过 Register 表检查，封号是针对账号的）
		register, err := db.QueryRegisterByUid(uid)
		if err == nil && register != nil {
			if register.BanUntil > 0 && time.Now().Unix() < register.BanUntil {
				res.Code = -1
				remainingSeconds := register.BanUntil - time.Now().Unix()
				remainingMinutes := remainingSeconds / 60
				if remainingMinutes < 1 {
					remainingMinutes = 1
				}
				banUntilTime := time.Unix(register.BanUntil, 0).Format("2006-01-02 15:04:05")
				res.Error = fmt.Sprintf("账号已被封禁，剩余时间：%d分钟，解封时间：%s", remainingMinutes, banUntilTime)
				logger.Warnf("玩家: %d 登录失败,账号被封禁,解封时间: %s", uid, banUntilTime)
				return s.Response(res)
			}
		}

		s.Response(res)
		// 隔天登录
		// 查询他是第几天
		p.syncUserDay()
		if p.lastLoginAt.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
			p.newDay()
		}

		m.setPlayer(uid, p)

		p.online()

		p.reset()

	} else {
		// 重置之前的session
		if prevSession := p.session; prevSession != nil && prevSession != s {
			logger.Warningf("玩家:%d 已经在线,顶号", uid)
			// 告知这个用户，他被顶号了
			packet := &pb.CompulsoryWithdrawalPush{
				Status: pb.CompulsoryWithdrawalPush_WithdrawalStatus(pb.CompulsoryWithdrawalPush_Pushoff),
			}
			p.pushMsg(protocol.RouteCompulsoryWithdrawalPush, packet)
			// 移除广播频道
			//m.group.Leave(prevSession)
			playerMgr.RemoveSessionFromGroup(prevSession)

			// 如果之前房间存在，则退出来
			if p, err := playerWithSession(prevSession); err == nil && p != nil && p.desk != nil && p.desk.group != nil {
				//p.desk.playerExit(p, true)
				p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: false}
				p.desk = nil
				p.setPlayerStatus(PlayerStatusDisconnect)
			}
			time.AfterFunc(time.Millisecond*10, func() {
				prevSession.Clear()
				prevSession.Close()
			})
		} else {
			// 断线重连情况，用户之前可能已经离线，需要增加计数
		}

		// 绑定新session
		p.bindSession(s)
		// 查询他是第几天
		p.syncUserDay()
		// 隔天登录
		if p.lastLoginAt.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
			p.newDay()
		}
		p.OfflineAdd()
		p.online()
		// 断线重连
		if p.roomId != 0 && p.desk != nil && p.curState == PlayerStatusDisconnect && time.Now().Before(p.disconnectTime) {
			res.GameStatus = int32(PlayerStatusPlaying)
			s.Response(res)
			defer func() {
				p.desk.playerRejoin(p)
			}()
			p.logger.Debugf("玩家: %d : %s 断线重连: 房间号 %d ", uid, p.name, p.roomId)
		} else {
			if p.desk != nil {
				p.desk.playerKick(p, constant.GameOverTypeTime)

				p.logger.Debugf("玩家: %d : %s 断线重连: 房间号 %d 超过短线重连时间", uid, p.name, p.roomId)
			}
			p.reset()
			s.Response(res)
			p.logger.Debugf("玩家: %d : %s 登录成功，读取缓存内数据，进入大厅", uid, p.name)
		}
	}

	// 添加到广播频道
	//m.group.Add(s)
	playerMgr.AddSessionToGroup(s)

	onlinePlayersGauge.Inc()

	p, err := m.GetPlayer(uid)
	if err != nil {
		log.Errorf("thinkingdatax get player failed: %v", err)
	}
	if err := thinkingdatax.TrackerIns.Set(strconv.FormatInt(uid, 10), p.accountName, tracking.M{
		"login_time_v2": time.Now().Format(timex.FormatDateTimeMillis),
	}); err != nil {
		log.Warnf("thinkingdatax set login_time_v2 failed: %v", err)
	}

	return nil
}

func (m *Manager) CheckOrder(s *session.Session, msg *protocol.CheckOrderReqeust) error {
	log.Infof("%+v", msg)

	return s.Response(&protocol.CheckOrderResponse{
		FangKa: 20,
	})
}

// 包裹信息
func (m *Manager) GetBag(s *session.Session, _ []byte) error {
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	packet := &pb.GetBagResponse{ItemList: &pb.Bag{ItemList: []*pb.ItemInfo{}}}

	// 背包检测道具过期
	for _, item := range p.bag.GetItemList() {
		itemCfg, err := global.GetItemCfg(item.CfgId)
		if err != nil {
			p.logger.Errorf("道具配置不存在， userId: %d, itemid: %d", p.Uid(), item.CfgId)
			continue
		}
		if itemCfg.Timeliness.Type == constant.ItemTimelinessShiChang {
			expired := time.Unix(item.CreateTime, 0).Add(time.Duration(itemCfg.Timeliness.Time) * time.Hour * 24)
			if expired.Before(time.Now()) {
				p.takeSourceByUid(item.Uid, 1, constant.ItemDelExpired)
				break
			}
		} else if itemCfg.Timeliness.Type == constant.ItemTimelinessQuJian {
			if time.Now().Unix() < itemCfg.Timeliness.StarTime || time.Now().Unix() > itemCfg.Timeliness.EndTime {
				p.takeSourceByUid(item.Uid, 1, constant.ItemDelExpired)
				break
			}
		}

		// 过滤穿戴装备
		if p.hasEquipment(item.Uid) {
			continue
		}

		// 过滤宝箱候选区装备
		if p.hasHomeOpenBox(item.Uid) {
			continue
		}
		packet.ItemList.ItemList = append(packet.ItemList.ItemList, item)
	}

	dailyItemUsed := make([]*pb.DailyItemUsed, 0, len(p.DailyItemUsed))
	for itemCfg, itemUsageCount := range p.DailyItemUsed {
		dailyItemUsed = append(dailyItemUsed, &pb.DailyItemUsed{
			CfgId: itemCfg,
			Num:   itemUsageCount,
		})
	}
	packet.DailyItemUsedList = dailyItemUsed

	return s.Response(packet)
}

// 使用道具
// func (m *Manager) ItemUse(s *session.Session, req *pb.ItemUseRequest) error {
// 	res := &pb.ItemUseResponse{ItemList: []*pb.ItemInfo{}}

// 	p, err := playerWithSession(s)
// 	if err != nil {
// 		res.Code = -1
// 		res.Error = err.Error()
// 		return s.Response(res)
// 	}

// 	itemInfo := p.getItemByUid(req.ItemUid)
// 	if itemInfo == nil {
// 		res.Code = -2
// 		res.Error = errutil.ErrItemNotEnough.Error()
// 		return s.Response(res)
// 	}

// 	itemCfg, err := global.GetItemCfg(itemInfo.CfgId)
// 	if err != nil {
// 		res.Code = -2
// 		res.Error = errutil.ErrItemNotEnough.Error()
// 		return s.Response(res)
// 	}

// 	if p.level < int32(itemCfg.Levellimit) {
// 		res.Code = -3
// 		res.Error = errutil.ErrLevel.Error()
// 		return s.Response(res)
// 	}

// 	if !p.checkTakeSourceByUid(req.ItemUid, int64(req.Num)) {
// 		res.Code = -4
// 		res.Error = errutil.ErrItemNotEnough.Error()
// 		return s.Response(res)
// 	}

// 	if itemUseCfg, ok := global.ItemUseCfg[itemInfo.CfgId]; ok {
// 		var itembagcheck, equipbagcheck bool
// 		for _, v := range itemUseCfg.GiveGoods {
// 			if getItemCfg, err := global.GetItemCfg(v.ID); err == nil {
// 				if getItemCfg.SonType == constant.ItemTypeEquipment {
// 					equipbagcheck = true
// 				}

// 				if getItemCfg.SonType == constant.ItemTypeMaterial ||
// 					getItemCfg.SonType == constant.ItemTypeGem ||
// 					getItemCfg.SonType == constant.ItemTypeConsumable {
// 					itembagcheck = true
// 				}
// 			}
// 		}

// 		if equipbagcheck && p.bagLeftCount(pb.ItemType_EQUIPMENT) == 0 {
// 			res.Code = -6
// 			res.Error = errutil.ErrBagFull.Error()
// 			return s.Response(res)
// 		}

// 		if itembagcheck && p.bagLeftCount(pb.ItemType_IN_BAG_DISPOSABLE) == 0 {
// 			res.Code = -6
// 			res.Error = errutil.ErrBagFull.Error()
// 			return s.Response(res)
// 		}
// 	}

// 	reward := []global.IdNum{}
// 	buffs := []global.IdKeyNum{}
// 	useNum := req.Num // 实际使用数量
// 	if itemUseCfg, ok := global.ItemUseCfg[itemInfo.CfgId]; ok {
// 		// 检查使用道具的上限
// 		for _, v := range itemUseCfg.GiveGoods {
// 			if v.ID == constant.ItemTypeEnergy {
// 				if p.energy >= int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) {
// 					res.Code = -6
// 					res.Error = errutil.EnergyMax
// 					return s.Response(res)
// 				} else {
// 					needEnergy := int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) - p.energy
// 					needNum := int32(math.Ceil(float64(needEnergy) / float64(v.Num)))
// 					if useNum > needNum {
// 						useNum = needNum
// 					}
// 					break
// 				}
// 			}
// 		}

// 		buffs = itemUseCfg.Buffeffect
// 		// 如身上有高级buff，则不添加低级buff
// 		if len(buffs) > 0 {
// 			for _, ib := range buffs {
// 				for _, v := range p.buffs {
// 					if ib.ID == v.Id {
// 						for _, attr := range v.BuffStats.Attrs {
// 							if ib.Key == int32(pb.Attribute_EQUIP_STAR_PERCENTAGE) && attr >= ib.Num {
// 								logger.Warnf("身上有高级buff，不能添加低级buff, buffId: %d, buffKey: %d, buffNum: %d, uid: %d", ib.ID, ib.Key, ib.Num, p.Uid())
// 								res.Error = errutil.ErrBuffAlreadyExist.Error()
// 								return s.Response(res)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 		jobreward := []global.IdNumJobWeight{}
// 		jobrewardweight := []global.IdNumJobWeight{}
// 		if len(itemUseCfg.GiveGoods) > 0 {
// 			for _, gg := range itemUseCfg.GiveGoods {
// 				if p.job == constant.MowJob(gg.Job) || gg.Job == 0 {
// 					jobreward = append(jobreward, gg)
// 				}
// 			}
// 		}
// 		var weight int32
// 		for _, jr := range jobreward {
// 			weight += jr.Weight
// 			jobrewardweight = append(jobrewardweight, global.IdNumJobWeight{ID: jr.ID, Num: jr.Num, Job: jr.Job, Weight: int32(weight)})
// 		}
// 		for i := 0; i < int(useNum); i++ {
// 			tmpReward := []global.IdNumJobWeight{}
// 			if weight == 0 {
// 				tmpReward = jobrewardweight
// 			} else {
// 				w := global.RandIn(int(weight))
// 				for _, jrw := range jobrewardweight {
// 					if w < int(jrw.Weight) {
// 						tmpReward = append(tmpReward, global.IdNumJobWeight{ID: jrw.ID, Num: jrw.Num, Job: jrw.Job, Weight: weight})
// 						break
// 					}
// 				}
// 			}

// 			for _, tr := range tmpReward {
// 				bFind := false
// 				for k, r := range reward {
// 					if r.ID == tr.ID {
// 						if rewardCfg, err := global.GetItemCfg(r.ID); err == nil {
// 							if rewardCfg.Maxnum > 1 || rewardCfg.SonType == 0 {
// 								reward[k].Num += tr.Num
// 								bFind = true
// 							}
// 						}
// 						break
// 					}
// 				}
// 				if !bFind {
// 					reward = append(reward, global.IdNum{ID: tr.ID, Num: tr.Num})
// 				}
// 			}

// 		}
// 	}
// 	res.ItemList = append(res.ItemList, &pb.ItemInfo{Uid: req.ItemUid, Num: int64(useNum), CfgId: itemInfo.CfgId, CreateTime: itemInfo.CreateTime})

// 	// 升星石
// 	num, err2 := p.IsStarsLevel(itemInfo.CfgId, useNum)
// 	if err2 != nil && num == 0 {
// 		res.Code = -6
// 		res.Error = err2.Error()
// 		return s.Response(res)
// 	}
// 	// 剩余的
// 	useNum = num
// 	// 命格
// 	_, num, err = p.IsFate(itemInfo.CfgId, useNum)
// 	if err != nil && num == 0 {
// 		res.Code = -7
// 		res.Error = err.Error()
// 		return s.Response(res)
// 	}
// 	// 剩余的
// 	useNum = num
// 	// 删除道具
// 	p.takeSourceByUid(req.ItemUid, int64(useNum), constant.ItemDelGame)

// 	// 发放奖励
// 	if len(reward) > 0 {
// 		if err = p.checkAddSourceBatch(reward); err == nil {
// 			p.addSourceBatch(reward, constant.ItemAddDisintegrate)
// 		} else {
// 			attachment := []protocol.Attachment{}
// 			for _, item := range reward {
// 				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
// 			}
// 			p.SendMail(constant.ItemAddDisintegrate, attachment)
// 		}
// 		// for _, v := range reward {
// 		// 	p.addSource(v.ID, int64(v.Num), constant.ItemAddGame)
// 		// 	res.ItemList = append(res.ItemList, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
// 		// }
// 	}

// 	// 道具buff
// 	for _, ib := range buffs {
// 		bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
// 		bufStats.setStats(pb.Attribute(ib.Key), ib.Num)
// 		p.AddBuff(ib.ID, p, 0, bufStats, 0, 0)
// 	}

// 	return s.Response(res)
// }

func (m *Manager) itemUse(ctx context.Context, player *Player, req *pb.ItemUseRequest) (resp *pb.ItemUseResponse, err error) {
	var (
		buffs   []global.IdKeyNum
		rewards []global.IdNum
	)

	itemInfo := player.getItemByUid(req.ItemUid)
	if itemInfo == nil {
		return nil, errutil.NewBizError(-2, fmt.Sprint("Item not found, uid: ", player.Uid(), ", itemUid: ", req.ItemUid), errutil.ErrItemNotEnough.Error())
	}

	itemCfg, err := global.GetItemCfg(itemInfo.CfgId)
	if err != nil {
		return nil, errutil.NewBizError(-2, fmt.Sprintf("Item config not found, uid: %d, itemUid: %d, cfgId: %d", player.Uid(), req.ItemUid, itemInfo.CfgId), errutil.ErrItemNotEnough.Error())
	}

	if player.level < int32(itemCfg.Levellimit) {
		return nil, errutil.NewBizError(-3, fmt.Sprintf("Player level too low, uid: %d, itemUid: %d, required level: %d, current level: %d", player.Uid(), req.ItemUid, itemCfg.Levellimit, player.level), errutil.ErrLevel.Error())
	}

	//单次使用道具数量限制
	if req.Num > int32(global.GlobalCfg.BaseValue.Maxuselimit.Num) {
		req.Num = int32(global.GlobalCfg.BaseValue.Maxuselimit.Num)
	}

	if !player.checkTakeSourceByUid(req.ItemUid, int64(req.Num)) {
		return nil, errutil.NewBizError(-4, fmt.Sprintf("Not enough items, uid: %d, itemUid: %d, required num: %d", player.Uid(), req.ItemUid, req.Num), errutil.ErrItemNotEnough.Error())
	}

	itemUseCfg, ok := global.ItemUseCfg[itemInfo.CfgId]
	if !ok {
		return nil, errutil.NewBizError(-5, fmt.Sprintf("Item use config not found, uid: %d, itemUid: %d, cfgId: %d", player.Uid(), req.ItemUid, itemInfo.CfgId), errutil.ErrServerInternal.Error())
	}

	// 检查背包和体力使用上限
	useNum := req.Num
	var isItemBagChecked, isEquipmentBagChecked bool
	for _, v := range itemUseCfg.GiveGoods {
		if getItemCfg, err := global.GetItemCfg(v.ID); err == nil {
			if getItemCfg.SonType == constant.ItemTypeEquipment {
				isEquipmentBagChecked = true
			}
			if getItemCfg.SonType == constant.ItemTypeMaterial ||
				getItemCfg.SonType == constant.ItemTypeGem ||
				getItemCfg.SonType == constant.ItemTypeConsumable {
				isItemBagChecked = true
			}
		}

		if v.ID == constant.ItemTypeEnergy {
			if player.energy >= int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) {
				return nil, errutil.NewBizError(-6, fmt.Sprintf("Energy is already at maximum, uid: %d", player.Uid()), errutil.EnergyMax)
			} else {
				needEnergy := int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) - player.energy
				needNum := int32(math.Ceil(float64(needEnergy) / float64(v.Num)))
				if useNum > needNum {
					useNum = needNum
				}
				break
			}
		}
	}
	if isEquipmentBagChecked && player.bagLeftCount(pb.ItemType_EQUIPMENT) == 0 {
		return nil, errutil.NewBizError(-6, fmt.Sprintf("Equipment bag is full, uid: %d", player.Uid()), errutil.ErrBagFull.Error())
	}
	if isItemBagChecked && player.bagLeftCount(pb.ItemType_IN_BAG_DISPOSABLE) == 0 {
		return nil, errutil.NewBizError(-6, fmt.Sprintf("Item bag is full, uid: %d", player.Uid()), errutil.ErrBagFull.Error())
	}

	// 检查使用道具的上限
	dailyItemUsageCount := player.DailyItemUsed[itemInfo.CfgId]
	if itemUseCfg.UseLimit > 0 {
		if (dailyItemUsageCount + int64(useNum)) > itemUseCfg.UseLimit {
			return nil, errutil.NewBizError(-8, fmt.Sprintf("Item daily use limit exceeded, uid: %d, itemUid: %d, cfgId: %d, dailyLimit: %d, alreadyUsed: %d, tryingToUse: %d", player.Uid(), req.ItemUid, itemInfo.CfgId, itemUseCfg.UseLimit, dailyItemUsageCount, useNum), errutil.ErrItemUseLimit.Error())
		}
		player.DailyItemUsed[itemInfo.CfgId] = dailyItemUsageCount + int64(useNum)
	}
	if itemInfo.CfgId == constant.ItemTypeLuckyBagOne || itemInfo.CfgId == constant.ItemTypeLuckyBagTwo || itemInfo.CfgId == constant.ItemTypeLuckyBagThree {
		player.tasks.dispatch(pb.TargetType_LUCKY_BAG_OPENED_TOTAL, 0, useNum)
	}

	switch global.ItemUseConfigType(itemUseCfg.Type) {
	case global.ItemUseConfigTypeFixed:
		if len(itemUseCfg.GiveGoods) > 0 {
			for _, giveGood := range itemUseCfg.GiveGoods {
				// if player.job == constant.MowJob(giveGood.Job) || giveGood.Job == 0 {
				rewards = append(rewards, global.IdNum{ID: giveGood.ID, Num: giveGood.Num * req.Num})
				// }
			}
		}
	case global.ItemUseConfigTypeRandom:
		picker := randx.New(itemUseCfg.GiveGoods, func(i global.IdNumJobWeight) int { return int(i.Weight) })
		for i := 0; i < int(req.Num); i++ {
			item, err := picker.Pick()
			if err != nil {
				return nil, errutil.NewBizError(-5, fmt.Sprintf("Failed to pick random item, uid: %d, itemUid: %d, cfgId: %d, error: %v", player.Uid(), req.ItemUid, itemInfo.CfgId, err), errutil.ErrServerInternal.Error())
			}
			rewards = append(rewards, global.IdNum{ID: item.ID, Num: item.Num})
		}
	case global.ItemUseConfigTypeBuff:
		buffs = itemUseCfg.Buffeffect
		// 如身上有高级buff，则不添加低级buff
		if len(buffs) > 0 {
			var ib *global.IdKeyNum
			player.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
				for _, b := range buffs {
					if b.ID != buff.Id {
						continue
					}
					buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
						if b.Key == int32(pb.Attribute_EQUIP_STAR_PERCENTAGE) && value >= b.Num {
							ib = &b
							return false
						}
						return true
					})
					if ib != nil {
						return false
					}
				}
				return true
			})
			if ib != nil {
				return nil, errutil.NewBizError(-7, fmt.Sprintf("Higher level buff already exists, cannot add lower level buff, buffId: %d, buffKey: %d, buffNum: %d, uid: %d", ib.ID, ib.Key, ib.Num, player.Uid()), errutil.ErrBuffAlreadyExist.Error())
			}
			//for _, v := range player.buffs {
			//	if ib.ID == v.Id {
			//		for _, attr := range v.BuffStats.Attrs {
			//			if ib.Key == int32(pb.Attribute_EQUIP_STAR_PERCENTAGE) && attr >= ib.Num {
			//				return nil, errutil.NewBizError(-7, fmt.Sprintf("Higher level buff already exists, cannot add lower level buff, buffId: %d, buffKey: %d, buffNum: %d, uid: %d", ib.ID, ib.Key, ib.Num, player.Uid()), errutil.ErrBuffAlreadyExist.Error())
			//			}
			//		}
			//	}
			//}
		}
		for _, ib := range buffs {
			bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
			bufStats.setStats(pb.Attribute(ib.Key), ib.Num)
			player.AddBuff(ib.ID, player, 0, bufStats, 0, 0)
		}
	case global.ItemUseConfigTypeSelect:
		var found bool
		// 自选的只能开一个
		if req.Num > 1 {
			return nil, errutil.NewBizError(-5, fmt.Sprintf("Cannot use select item with quantity greater than 1, uid: %d, itemUid: %d, cfgId: %d, selectedItemId: %d", player.Uid(), req.ItemUid, itemInfo.CfgId, req.SelectItemID), errutil.ErrIllegalParameter.Error())
		}
		for _, v := range itemUseCfg.GiveGoods {
			if v.ID == req.SelectItemID {
				found = true
				rewards = append(rewards, global.IdNum{ID: v.ID, Num: v.Num})
				break
			}
		}
		if !found {
			return nil, errutil.NewBizError(-5, fmt.Sprintf("Selected item not in options, uid: %d, itemUid: %d, cfgId: %d, selectedItemId: %d", player.Uid(), req.ItemUid, itemInfo.CfgId, req.SelectItemID), errutil.ErrIllegalParameter.Error())
		}
	case global.ItemUseConfigTypeDirectLogic:
		// NOTE 现在只有 15分钟挂机卡 业务逻辑直接和次数关联，这块代码暂时没有
		// switch itemInfo.CfgId {
		// case constant.ItemTypeAfkCard: // 15分钟挂机卡
		// }
	default:
		return nil, errutil.NewBizError(-5, fmt.Sprintf("Unsupported item use config type: %d, uid: %d, itemUid: %d", itemUseCfg.Type, player.Uid(), req.ItemUid), errutil.ErrServerInternal.Error())
	}

	// 发放奖励
	if len(rewards) > 0 {
		if err = player.checkAddSourceBatch(rewards); err == nil {
			player.addSourceBatch(rewards, constant.ItemAddDisintegrate)
		} else {
			attachment := []protocol.Attachment{}
			for _, item := range rewards {
				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
			}
			player.SendMail(constant.ItemAddDisintegrate, attachment)
		}
	}

	// 删除道具
	player.takeSourceByUid(req.ItemUid, int64(useNum), constant.ItemDelGame)

	return &pb.ItemUseResponse{
		ItemList: []*pb.ItemInfo{{Uid: req.ItemUid, Num: int64(useNum), CfgId: itemInfo.CfgId, CreateTime: itemInfo.CreateTime}},
	}, nil
}

// 销毁道具
func (m *Manager) ItemDestroy(s *session.Session, req *pb.ItemDestroyRequest) error {
	res := &pb.ItemDestroyResponse{Item: &pb.ItemInfo{Uid: req.ItemUid, Num: int64(req.Num)}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if p.isOnEquip(req.ItemUid) {
		res.Code = -3
		res.Error = errutil.EquipOnUse
		return s.Response(res)
	}

	if !p.checkTakeSourceByUid(req.ItemUid, int64(req.Num)) {
		res.Code = -2
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}

	var getitem []global.KeyValue

	itemInfo := p.getItemByUid(req.ItemUid)
	if itemInfo == nil {
		res.Code = -3
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	itemCfg, err := global.GetItemCfg(itemInfo.CfgId)
	if err != nil {
		res.Code = -4
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}

	for _, des := range itemCfg.Decompose {
		bFind := false
		for _, gi := range getitem {
			if gi.Key == des.Key {
				gi.Value += des.Value
				bFind = true
				break
			}
		}
		if !bFind {
			getitem = append(getitem, des)
		}
	}

	p.takeSourceByUid(req.ItemUid, int64(req.Num), constant.ItemDelDistroy)
	for _, v := range getitem {
		p.addSource(v.Key, int64(v.Value*req.Num), constant.ItemAddDistroy)
	}
	return s.Response(res)
}

// 分解道具
func (m *Manager) ItemDisassembly(s *session.Session, req *pb.ItemDisassemblyRequest) error {
	res := &pb.ItemDisassemblyResponse{}
	res.ItemUids = req.ItemUids

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	var getitem []global.KeyValue
	for _, v := range req.ItemUids {
		if !p.checkTakeSourceByUid(v, 1) {
			p.boxResIntoBag(v)
			continue
			// res.Code = -2
			// res.Error = errutil.ErrItemNotEnough.Error()
			// return s.Response(res)
		}

		if p.isOnEquip(v) {
			p.boxResIntoBag(v)
			res.Code = -3
			res.Error = errutil.EquipOnUse
			return s.Response(res)
		}

		itemInfo := p.getItemByUid(v)
		if itemInfo == nil {
			continue
		}
		itemCfg, err := global.GetItemCfg(itemInfo.CfgId)
		if err != nil {
			continue
		}

		for _, des := range itemCfg.Decompose {
			// bFind := false
			// for i, gi := range getitem {
			// 	if gi.Key == des.Key {
			// 		if itemCfg.Quality > 0 && itemCfg.Quality <= 8 {
			// 			enhance := global.GlobalCfg.Equipenhances[itemCfg.Quality-1].Value
			// 			getitem[i].Value += int32(float32(des.Value) * enhance)
			// 			bFind = true
			// 			break
			// 		}
			// 	}
			// }
			// if !bFind {
			// 	getitem = append(getitem, des)
			// }
			if itemCfg.Quality > 0 && itemCfg.Quality <= 8 {
				if itemInfo.Equip != nil && itemInfo.Equip.Quality > 0 &&
					int(itemInfo.Equip.Quality-1) < len(global.GlobalCfg.Equipenhances) {
					enhance := global.GlobalCfg.Equipenhances[itemInfo.Equip.Quality-1].Value
					getitem = append(getitem, global.KeyValue{Key: des.Key, Value: int32(float32(des.Value) * enhance)})
				}
			}
		}
		// 如果分解的是首页的装备，需要把首页的装备删掉
		p.boxResIntoBag(v)
	}

	delSource := []global.UUIdNum{}
	for _, v := range req.ItemUids {
		delSource = append(delSource, global.UUIdNum{UID: v, Num: 1})
	}
	p.takeSourceByUidBatch(delSource, constant.ItemDelDisintegrate)

	addSource := []global.IdNum{}
	for _, v := range getitem {
		addSource = append(addSource, global.IdNum{ID: v.Key, Num: v.Value})
	}
	p.addSourceBatch(addSource, constant.ItemAddDisintegrate)

	return s.Response(res)
}

// 自动分解
func (m *Manager) AutoDisassembly(s *session.Session, req *pb.AutoDisassemblyRequest) error {
	p, err := playerWithSession(s)
	if err != nil {
		return nil
	}

	res := &pb.AutoDisassemblyResponse{}
	if p.vip < global.GlobalCfg.BaseValue.Autodecompose.Viplevel && p.level < global.GlobalCfg.BaseValue.Autodecompose.Rolelevel {
		res.Code = -1
		res.Error = errutil.ErrLevel.Error()
		return s.Response(res)
	}

	if req.Info.Level <= 0 || req.Info.Quality <= pb.ItemQuality_INVALID_Quality || req.Info.Quality > pb.ItemQuality_RED {
		res.Code = -1
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}

	p.ADLevel = req.Info.Level
	p.ADQuality = int32(req.Info.Quality)
	p.ADOpen = req.Info.Open
	p.ADHighCombat = req.Info.HighCombat
	p.ADLevelSwitch = req.Info.LevelSwitch
	p.ADQualitySwitch = req.Info.QualitySwitch

	res.Open = p.ADOpen
	res.Level = p.ADLevel
	res.Quality = pb.ItemQuality(p.ADQuality)
	res.LevelSwitch = p.ADLevelSwitch
	res.QualitySwitch = p.ADQualitySwitch
	res.HighCombat = p.ADHighCombat

	return s.Response(res)
}

// 穿脱装备
func (m *Manager) ChangeEquip(s *session.Session, req *pb.ChangeEquipRequest) error {
	res := &pb.ChangeEquipResponse{OpType: req.OpType, Slot: req.Slot, PartnerSlotId: req.PartnerSlotId, Item: &pb.ItemInfo{}, OldItem: &pb.ItemInfo{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	if req.Slot <= pb.EquipSlot_INVALID_EquipSlot || req.Slot >= pb.EquipSlot_MaxEquipSlot {
		res.Code = -2
		res.Error = errutil.ErrItemNotEnough.Error()
		logger.Debugf("玩家: %d 穿脱装备伙伴失败,道具不存在了, uid=%d", p.Uid(), req.ItemUid)
		return s.Response(res)
	}
	ItemInfo := p.getItemByUid(req.ItemUid)
	if ItemInfo == nil {
		logger.Debugf("玩家: %d 穿脱装备失败,道具不存在, uid=%d", p.Uid(), req.ItemUid)
		res.Code = -3
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	res.Item = ItemInfo
	if req.PartnerSlotId == 0 {
		// 前端要求如果身上有装备，那么老装备就是身上装备，如果身上没装备，老装备就是自身
		if p.equips[req.Slot].ItemId != 0 {
			res.OldItem = p.getItemByUid(p.equips[req.Slot].ItemId)
		} else {
			res.OldItem = res.Item
		}
	} else {
		// 前端要求如果身上有装备，那么老装备就是身上装备，如果身上没装备，老装备就是自身
		partnerSlotInfo := p.getPartnerSlotInfo(req.PartnerSlotId)
		if partnerSlotInfo == nil {
			res.Code = -4
			res.Error = errutil.ErrItemNotEnough.Error()
			logger.Debugf("玩家: %d 穿脱装备伙伴失败,道具不存在, uid=%d", p.Uid(), req.ItemUid)
			return s.Response(res)
		}
		if len(partnerSlotInfo.Equips) == 0 {
			res.Code = -5
			res.Error = errutil.ErrItemNotEnough.Error()
			logger.Debugf("玩家: %d 穿脱装备伙伴失败,道具不存在了, uid=%d", p.Uid(), req.ItemUid)
			return s.Response(res)
		}
		equipSlot := partnerSlotInfo.Equips[req.Slot]
		if equipSlot.ItemId != 0 {
			endInfo := p.getItemByUid(equipSlot.ItemId)
			if endInfo == nil {
				res.Code = -6
				res.Error = errutil.ErrItemNotEnough.Error()
				logger.Debugf("玩家: %d 穿脱装备伙伴失败了,道具不存在, uid=%d", p.Uid(), req.ItemUid)
				return s.Response(res)
			}
			res.OldItem = endInfo
		} else {
			res.OldItem = res.Item
		}
	}
	if err = p.changeEquip(req.OpType, req.Slot, req.ItemUid, req.PartnerSlotId); err != nil {
		res.Code = -7
		res.Error = err.Error()
		return s.Response(res)
	}
	if req.PartnerSlotId == 0 {
		p.SyncUserEquipsInfo()
	} else {
		p.SyncUserPartnerSlot()
	}

	p.refreshStats()
	return s.Response(res)
}

// 完成任务
func (m *Manager) TaskComplate(s *session.Session, req *pb.TaskComplateRequest) error {
	res := &pb.TaskComplateResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if err = p.tasks.complate(req.TaskId); err != nil {
		res.Code = -2
		res.Error = err.Error()
	}

	return s.Response(res)
}

// 洗炼
func (m *Manager) EquipRefine(s *session.Session, req *pb.EquipRefineRequest) error {
	res := &pb.EquipRefineResponse{NewStats: &pb.Stats{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if err = p.equipRefine(req.ItemUid, req.Attr, req.Random); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}

	res.NewStats = &p.lastEquipRefine.newStat

	return s.Response(res)
}

func (m *Manager) EquipRefineReplace(s *session.Session, req *pb.EquipRefineReplaceRequest) error {
	res := &pb.EquipRefineReplaceResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		return err
	}

	if item, err := p.equipRefineReplace(req.Replace); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	} else {
		res.Item = item
	}

	return s.Response(res)
}

func (m *Manager) chat(p *Player, req *pb.ChatRequest) (resp *pb.ChatResponse, err error) {
	//禁言不能说话
	if p.muteUntil > 0 && time.Now().Unix() < p.muteUntil {
		// 计算剩余禁言时间（分钟）
		remainingSeconds := p.muteUntil - time.Now().Unix()
		remainingMinutes := remainingSeconds / 60
		if remainingMinutes < 1 {
			remainingMinutes = 1 // 至少显示1分钟
		}
		errorMsg := fmt.Sprintf("您当前处于禁言状态，距离解封还有%d分钟", remainingMinutes)
		return &pb.ChatResponse{Code: -1, Error: errorMsg}, nil
	}
	// gm指令
	appEnv := strings.ToLower(viper.GetString("core.env"))
	if len(req.Context) > 0 && req.Context[0] == '@' && appEnv != constant.EnvProd && appEnv != constant.EnvReview && appEnv != constant.EnvStag {
		cmd := strings.Split(req.Context, " ")
		var id int
		var num int
		var err error
		switch cmd[0] {
		case "@add":
			if len(cmd) == 3 {
				if id, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if num, err = strconv.Atoi(cmd[2]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -2, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				p.addSource(int32(id), int64(num), constant.ItemAddGM)
			} else if len(cmd) == 4 {
				num := 999
				begin := 1
				end := 1
				if begin, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if end, err = strconv.Atoi(cmd[2]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if num, err = strconv.Atoi(cmd[3]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				var items []global.IdNum

				for i := begin; i <= end; i++ {
					items = append(items, global.IdNum{ID: int32(i), Num: int32(num)})
				}

				p.addSourceBatch(items, constant.ItemAddGM)
			}
		case "@clear":
			newbag := &pb.Bag{ItemList: []*pb.ItemInfo{}}
			packet := pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
			for _, v := range p.bag.ItemList {
				if p.hasEquipment(v.Uid) {
					newbag.ItemList = append(newbag.ItemList, v)
					continue
				}
				packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: v.CfgId, Uid: v.Uid, Num: v.Num})
			}
			p.pushMsg(protocol.RouteTakeSourcePush, &packet)
			p.bag = newbag
		case "@god":
			godLevel := int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)
			if len(cmd) == 2 {
				if godLevel, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -3, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
			}

			if godLevel > int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel) {
				godLevel = int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)
			}
			if godLevel < int(p.level) {
				logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
				return &pb.ChatResponse{Code: -4, Error: errutil.ErrIllegalParameter.Error()}, nil
			}

			roleCfg, err := global.GetRoleCfg(int32(p.job), int32(godLevel))
			if err != nil {
				logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
				return &pb.ChatResponse{Code: -5, Error: errutil.ErrIllegalParameter.Error()}, nil
			}

			// p.coin = 200000
			// p.soul = 0
			// p.coupon = 0
			p.tokens[constant.ItemTypeCoin] = 200000
			p.tokens[constant.ItemTypeSoul] = 10
			p.tokens[constant.ItemTypeCoupon] = 10
			p.tokens[constant.ItemTypeTrialVoucher] = 120

			p.level = int32(godLevel - 1)
			p.exp = roleCfg.NeedExp

			p.vip = 3
			p.energy = 500
			p.levelUp()
			p.artifactUnlock()

			guildManager.GuildAddExp(p.guildId, 1081000)

			for k := range p.dailyMaps {
				p.dailyMaps[k] = 10
			}

			p.PVPDuel.dailyBuyNum = 10
			p.PVPMelee.dailyBuyNum = 10
			p.dailyWorldBossTimes = 0
			p.totalWorldBossDamage = 0

			p.afkDailyTime = 0
			p.mainBattleId = 1000

			for _, v := range p.partnerSlots {
				if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
					partner.level = int32(godLevel)
					partner.SyncUserPartner()
				}
			}

			items := []global.IdNum{}
			for _, v := range global.ItemCfg {
				if v.Itemlevel == int32(math.Min(15, 15*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel))) {
					items = append(items, global.IdNum{ID: v.ID, Num: 15})
				}
				if v.MainType != 5 {
					continue
				}
				if v.Occupylimit != 0 && v.Occupylimit != int(p.job) {
					continue
				}
				if v.Levellimit != int(p.level/10*10) {
					continue
				}

				items = append(items, global.IdNum{ID: v.ID, Num: 1})
			}
			p.addSourceBatch(items, constant.ItemAddGM)

			for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
				p.genstoneChangeAll(i, 0)
			}
			for i := 1; i < len(p.equips); i++ {
				p.equips[i].Stars = 30

				p.equips[i].Enhance = p.level / 2

				for _, parslot := range p.partnerSlots {
					parslot.Equips[i].Stars = p.equips[i].Stars
					parslot.Equips[i].Enhance = p.equips[i].Enhance
				}

				for _, v := range p.bag.ItemList {
					if equipCfg, err := global.GetItemCfg(v.CfgId); err == nil {
						if equipCfg.Equiptype == i {
							p.equips[i].ItemId = v.Uid
						}
					}
				}
			}

			p.skillsSlot = [pb.SkillsSlot_MaxSkillsSlot]int32{}
			for i := 0; i < len(p.skills); i++ {
				if scfg, err := global.GetSkillCfg(p.skills[i].CfgId); err == nil {
					if scfg.Level == 0 || scfg.ActiveorPassive == PassiveSkill || scfg.SType == constant.SkillTypeArtifact {
						continue
					}

					p.skills[i].AdvancedLevel = int32(math.Min(6, 6*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)))
					p.skills[i].CfgId = scfg.SkillID*100 + int32(math.Min(30, 30*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)))

				}
			}
		case "@addpartner":
			if len(cmd) == 2 {
				if id, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}

				if err := p.addPartner(int32(id)); err == nil {
					logger.Debugf("@添加伙伴成功 cmd=%s partnerId=%d ", req.Context, id)
				}
			} else {
				for _, v := range global.PartnerCfg {
					if err := p.addPartner(v.Cid); err == nil {
						logger.Debugf("@添加伙伴成功 cmd=%s  partnerId=%d ", req.Context, v.Cid)
					}
				}
				for _, v := range p.partners {
					v.level = 80
					v.star = 20

					v.learnSkill()

					v.SyncUserPartner()
				}
			}
		case "@addpet":
			if len(cmd) == 3 {
				if id, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if num, err = strconv.Atoi(cmd[2]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -2, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if _, err := p.addPet(int32(id), 0, 0, int32(num), 0, "gm"); err == nil {
					logger.Debugf("@添加宠物成功 cmd=%s petId=%d ", req.Context, id)
				}
			}
		case "@setbat":
			if len(cmd) == 2 {
				if num, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				// if num, err = strconv.Atoi(cmd[2]); err != nil {
				// 	logger.Debugf("@参数错误 cmd=%s ", req.Context)
				// 	res.Code = -2
				// 	break
				// }
				id := num / 100000
				cfg := global.GetBattleMapCfg(int32(num))
				if cfg == nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -2, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if id == int(constant.BattleMap) {
					p.InstanceData.BattleMap.BattleId = int32(num)
				} else if id == int(constant.BossMap) {
					p.InstanceData.BossMap.BattleId = int32(num)
				} else if id == int(constant.TrialTowerMap) {
					p.InstanceData.TrialTowerMap.BattleId = int32(num)
				} else if id == int(constant.RelicMap) {
					for _, v := range global.BattleCfg {
						if v.CopyType != int32(id) {
							continue
						}
						p.InstanceData.RelicMap.Evaluation = append(p.InstanceData.RelicMap.Evaluation, global.KeyValue{Key: v.UID, Value: 3})
					}
					p.InstanceData.RelicMap.BattleId = int32(num)
				}
			}
		case "@addtitle":
			num := 999
			if len(cmd) == 2 {
				if num, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
			}
			var items []global.IdNum

			for i := 30001; i <= 30963; i++ {
				items = append(items, global.IdNum{ID: int32(i), Num: int32(num)})
			}

			p.addSourceBatch(items, constant.ItemAddGM)

		case "@power":
			num := 9999
			id := p.AttackAttributes()
			if len(cmd) == 2 {
				if arg1, err := strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				} else {
					num = arg1
				}
			}
			if len(cmd) == 3 {
				if arg1, err := strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				} else {
					id = pb.Attribute(arg1)
				}

				if arg2, err := strconv.Atoi(cmd[2]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				} else {
					num = arg2
				}
			}
			if p.equips[pb.EquipSlot_Weapon].ItemId == 0 {
				return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
			}

			if equip := p.getItemByUid(p.equips[pb.EquipSlot_Weapon].ItemId); equip != nil {
				new := true
				for _, v := range equip.Equip.Attributes {
					if v.Attr == id {
						v.Num = int32(num)
						new = false
						break
					}
				}
				if new {
					equip.Equip.Attributes = append(equip.Equip.Attributes, &pb.Stats{Attr: id, Num: int32(num)})
				}
			}

			for _, v := range p.partnerSlots {
				if equip := p.getItemByUid(v.Equips[pb.EquipSlot_Weapon].ItemId); equip != nil {
					new := true
					for _, v := range equip.Equip.Attributes {
						if v.Attr == id {
							v.Num = int32(num)
							new = false
							break
						}
					}
					if new {
						equip.Equip.Attributes = append(equip.Equip.Attributes, &pb.Stats{Attr: id, Num: int32(num)})
					}
				}
			}

		case "@addmon":
			id := 0
			num := 1
			if len(cmd) == 2 {
				if id, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if cfg := global.GetMonsterCfg(int32(id)); cfg == nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
			}
			if len(cmd) == 3 {
				if num, err = strconv.Atoi(cmd[2]); err != nil || num > 100 {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
			}
			if p.desk != nil {
				for i := 0; i < num; i++ {
					p.desk.addMonsterRelive(int32(id), p.cell.location, 0, 0, 0, 0, "")
				}
			}

		case "@main":
			mainId := 0
			maxMain := len(global.MainBattleCfg)
			if len(cmd) == 2 {
				if mainId, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -3, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if mainId < 1 {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -3, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
				if mainId >= maxMain {
					mainId = maxMain - 500
				}
				p.mainBattleId = int32(mainId)
			}
		case "@fc":
			if len(cmd) != 2 {
				logger.Errorf("gm指令错误，参数错误[格式:@fc 1] cmd=%s ", req.Context)
				return &pb.ChatResponse{Code: -3, Error: errutil.ErrIllegalParameter.Error()}, nil
			}
			cfgId, err := strconv.Atoi(cmd[1])
			if err != nil {
				logger.Errorf("gm指令错误，参数错误[格式:@fc 1] cmd=%s ", req.Context)
				return &pb.ChatResponse{Code: -3, Error: errutil.ErrIllegalParameter.Error()}, nil
			}
			desk := p.desk
			if desk == nil {
				logger.Errorf("该Gm命令只能在场景里使用 cmd=%s ", req.Context)
				return &pb.ChatResponse{Code: -3, Error: errutil.ErrNotFound.Error()}, nil
			}
			desk.eventMgr.Publish(event.EventTypeFateCoinEffect, &event.EventData{UId: p.uid, CfgId: int32(cfgId)})

		case "@artifact":
			num := 9
			if len(cmd) == 2 {
				if num, err = strconv.Atoi(cmd[1]); err != nil {
					logger.Errorf("gm指令错误，参数错误 cmd=%s ", req.Context)
					return &pb.ChatResponse{Code: -1, Error: errutil.ErrIllegalParameter.Error()}, nil
				}
			}

			addartifacts := make(map[int32]*pb.ArtifactInfo)
			delartifacts := []int32{}
			p.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {

				// p.pushMsg(protocol.RouteArtifactUpdatePush, &pb.ArtifactUpdatePush{
				// 	Artifact: &pb.ArtifactInfo{CfgId: value.CfgId},
				// 	Type:     pb.ArtifactUpdatePush_DELETE,
				// })

				oldId := value.CfgId
				value.CfgId = value.CfgId/10*10 + int32(num)
				for _, v := range value.Soul {
					v.Level = 10 * int32(num)
				}

				delartifacts = append(delartifacts, oldId)
				addartifacts[value.CfgId] = value

				// p.pushMsg(protocol.RouteArtifactUpdatePush, &pb.ArtifactUpdatePush{
				// 	Artifact: value,
				// 	Type:     pb.ArtifactUpdatePush_NEW,
				// })

				return true
			})

			for _, v := range addartifacts {
				p.artifacts.Set(v.CfgId, v)

			}

			for _, v := range delartifacts {
				p.artifacts.Delete(v)
			}

			newskills := []int32{} // 技能
			for _, v := range p.skills {
				if scfg, err := global.GetSkillCfg(v.CfgId); err == nil {
					if scfg.SType == constant.SkillTypeArtifact {
						continue
					}
					newskills = append(newskills, v.CfgId)
				}
			}
			p.skills = []*PlayerSkill{}
			p.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
				if cfg, ok := global.ArtifactAdvanceCfg[key]; ok {
					p.addSkill(cfg.SkillID)
				}
				return true
			})

			for _, v := range newskills {
				p.addSkill(v)
			}

			p.syncArtifact()
		}

		p.syncUserBase()
		p.SyncUserAdvance()
		p.syncUserSkillSlot()
		p.refreshStats()
		p.SyncUserPartnerSlot()

		return &pb.ChatResponse{
			Code: 0,
		}, nil
	}
	if !p.isModuleOpen(constant.ModuleChatFunc) {
		return &pb.ChatResponse{Code: -21, Error: errutil.ErrModuleNotOpen}, nil
	}
	// 聊天
	if req.Type == pb.ChatType_WORLD {
		if p.getItemCountInBag(constant.ItemTypeLoudspeaker) >= 1 {
			p.takeSource(constant.ItemTypeLoudspeaker, 1)
		} else {
			if !p.checkTakeSource(global.GlobalCfg.BaseValue.Chatconsume.ID, global.GlobalCfg.BaseValue.Chatconsume.Num) {
				logger.Errorf("玩家: %d 世界聊天道具不足, uid=%d", p.Uid(), global.GlobalCfg.BaseValue.Chatconsume.ID)
				return &pb.ChatResponse{Code: -1, Error: errutil.ErrItemNotEnough.Error()}, nil
			}
			p.takeSource(global.GlobalCfg.BaseValue.Chatconsume.ID, global.GlobalCfg.BaseValue.Chatconsume.Num, constant.ItemDelChat)
		}
		packet := &pb.ChatPush{Type: req.Type, Context: sensitive.ReplaceSensitiveWords(req.Context), Uid: p.uid, Name: p.name}
		if p.serverID != 0 {
			//for _, v := range m.players {
			//	if v.serverID == p.serverID {
			//		if err := v.pushMsg(protocol.RouteChatPush, packet); err != nil {
			//			logger.Errorf("玩家: %d 世界聊天失败,推送失败, uid=%d, err=%s", p.Uid(), p.uid, err.Error())
			//		}
			//		// 添加到聊天记录
			//		// if v.serverID == p.serverID {
			//		// 	v.pushMsg(protocol.RouteChatPush, packet)
			//	}
			//}
			if err := playerMgr.BroadcastToGroupByServerID(p.serverID, protocol.RouteChatPush, packet); err != nil {
				logger.Errorf("玩家: %d 世界聊天失败,推送失败, uid=%d, err=%s", p.Uid(), p.uid, err.Error())
			}
			defaultManager.addChatHistory(p.serverID, int32(pb.ChatType_WORLD), p.uid, sensitive.ReplaceSensitiveWords(req.Context), p.name, 0)
		}
	}

	// 如果是行会聊天
	if req.Type == pb.ChatType_GUILD {
		if p.guildId == 0 {
			logger.Errorf("玩家: %d 行会聊天失败,没有行会, uid=%d", p.Uid(), p.uid)
			return &pb.ChatResponse{Code: -1, Error: errutil.ErrGuildNotFound.Error()}, nil
		}

		packet := &pb.ChatPush{Type: req.Type, Context: sensitive.ReplaceSensitiveWords(req.Context), Uid: p.uid, Name: p.name}
		if guildManager.guilds[p.guildId] != nil {
			for _, v := range guildManager.guilds[p.guildId].Members {
				// 在线的才给推送
				// v.MemberId
				// 判断是否在线
				if pe, ok := m.player(v.MemberId); ok {
					if pe.curState <= PlayerStatusDisconnect {
						if err := pe.pushMsg(protocol.RouteChatPush, packet); err != nil {
							logger.Errorf("玩家: %d 行会聊天失败,推送失败, uid=%d, err=%s", p.Uid(), p.uid, err.Error())
						}
					}
				}
			}
			defaultManager.addChatHistory(p.serverID, int32(pb.ChatType_GUILD), p.uid, sensitive.ReplaceSensitiveWords(req.Context), p.name, p.guildId)
		}
	}

	if req.Type == pb.ChatType_ZONE {
		name := fmt.Sprintf("[%d服]%s", p.serverID%1000, p.name)
		var zoneID int32
		packet := &pb.ChatPush{Type: req.Type, Context: sensitive.ReplaceSensitiveWords(req.Context), Uid: p.uid, Name: name}
		if serverInfo, err := global.ServerListMgrInstance.Get(int(p.serverID)); err == nil {
			//for _, p := range m.players {
			//	for _, v := range serverlist {
			//		if v == int(p.serverID) {
			//			p.pushMsg(protocol.RouteChatPush, packet)
			//			break
			//		}
			//	}
			//}
			zoneID = serverInfo.ZoneId
			serverlist := global.ServerListMgrInstance.GetZoneServerIds(zoneID)
			for _, v := range serverlist {
				playerMgr.BroadcastToGroupByServerID(int32(v), protocol.RouteChatPush, packet)
			}
		}
		defaultManager.addChatHistory(-zoneID, int32(pb.ChatType_ZONE), p.uid, sensitive.ReplaceSensitiveWords(req.Context), name, 0)
	}

	return &pb.ChatResponse{
		Code: 0,
	}, nil
}

// 强化
func (m *Manager) Enhance(s *session.Session, req *pb.EnhanceRequest) error {
	res := &pb.EnhanceResponse{Code: 0, Slot: req.Slot, PartnerSlotId: req.PartnerSlotId}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleEquipStrengthen) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	if req.Slot <= pb.EquipSlot_INVALID_EquipSlot || req.Slot >= pb.EquipSlot_MaxEquipSlot {
		res.Code = -2
		res.Error = errutil.ErrEquipSlot.Error()
		return s.Response(res)
	}
	_, level, err := p.enhance(req.Slot, req.PartnerSlotId)
	if err != nil {
		res.Code = -3
		res.Error = err.Error()
		res.NewLevel = level
		res.Slot = req.Slot
		return s.Response(res)
	}

	res.NewLevel = level
	res.Slot = req.Slot
	return s.Response(res)
}

// 升星
func (m *Manager) Stars(s *session.Session, req *pb.StarsRequest) error {
	res := &pb.StarsResponse{Code: 0, Slot: req.Slot, PartnerSlotId: req.PartnerSlotId}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	if !p.isModuleOpen(constant.ModuleEquipStarUp) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	if req.Slot <= pb.EquipSlot_INVALID_EquipSlot || req.Slot >= pb.EquipSlot_MaxEquipSlot {
		res.Code = -2
		res.Error = errutil.ErrEquipSlot.Error()
		return s.Response(res)
	}
	check, level, err := p.stars(req.Slot, req.PartnerSlotId)
	if err != nil {
		res.Code = -3
		res.Error = err.Error()
		res.NewLevel = level
		res.Slot = req.Slot
		return s.Response(res)
	}
	if check {
		res.Error = errutil.StarsSuccess.Error()
	} else {
		res.Code = 1
		res.Error = errutil.StarsError.Error()
	}
	res.NewLevel = level
	res.Slot = req.Slot
	return s.Response(res)
}

// 使用升星道具
func (m *Manager) StarsLevel(s *session.Session, req *pb.ItemUseRequest) error {
	res := &pb.StarsLevelResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.checkTakeSourceByUid(req.ItemUid, int64(req.Num)) {
		res.Code = -3
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	Starlevel := global.GlobalCfg.Starstone.Starlevel
	if len(Starlevel) == 0 {
		res.Code = -4
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	checkId := int32(-1)
	slot := int32(0)
	for _, v := range Starlevel {
		if int64(v.ID) == req.ItemUid {
			checkId = v.Level
			slot = v.Slot
		}
	}
	if checkId == -1 {
		res.Code = -4
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	// selectId := int64(-1)
	if slot != 0 {
		if p.equips[slot].Stars < checkId {
			p.equips[slot].Stars = checkId
		} else {
			res.Code = -5
			res.Error = errutil.StarsNoUseError.Error()
			return s.Response(res)
		}
	} else {
		selectId := make([]int, 0)
		for i, v := range p.equips {
			// 需要随机给他一个槽位升满
			if v.Stars < checkId {
				selectId = append(selectId, i)
			}
		}
		if len(selectId) == 0 {
			res.Code = -5
			res.Error = errutil.StarsNoUseError.Error()
			return s.Response(res)
		}
		weight := global.RandIn(len(selectId))
		p.equips[weight].Stars = checkId
	}

	// 删除道具
	p.takeSourceByUid(req.ItemUid, int64(req.Num), constant.ItemDelGame)

	return s.Response(res)
}

// 穿脱技能
func (m *Manager) SkillChange(s *session.Session, req *pb.SkillChangeRequest) error {
	res := &pb.SkillChangeResponse{OpType: req.OpType, Slot: req.Slot, SkillId: req.SkillId}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if err = p.skillChange(req.OpType, req.Slot, req.SkillId); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

// 升级技能
func (m *Manager) SkillUpgrade(s *session.Session, req *pb.SkillUpgradeRequest) error {
	res := &pb.SkillUpgradeResponse{Code: 0, PartnerUid: req.PartnerUid}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleSkillUpgrade) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	newSkillId, err := p.skillsUpgrade(req.SkillId, req.PartnerUid)
	if err != nil {
		res.Code = -3
		res.Error = err.Error()
		res.NewSkillId = newSkillId
		return s.Response(res)
	}
	res.OldSkillId = req.SkillId
	res.NewSkillId = newSkillId
	return s.Response(res)
}

// 进阶技能
func (m *Manager) SkillAdvanced(s *session.Session, req *pb.SkillAdvancedRequest) error {
	res := &pb.SkillAdvancedResponse{Code: 0, PartnerUid: req.PartnerUid}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleSkillAdvance) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	newSkillId, err := p.skillsAdvanced(req.SkillId, req.PartnerUid)
	if err != nil {
		res.Code = -3
		res.Error = err.Error()
		res.NewLevel = newSkillId
		return s.Response(res)
	}
	res.OldSkillId = req.SkillId
	res.NewLevel = newSkillId
	return s.Response(res)
}

// 个人信息+战区id
func (m *Manager) GetPlayerInfo(uid int64) (*protocol.UserInfo, error) {
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return nil, errutil.ErrPlayerNotFound
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	userInfo := &protocol.UserInfo{
		UID:                    p.uid,
		Name:                   p.name,
		Level:                  p.level,
		Vip:                    p.vip,
		BagNumEquip:            p.bagEquipLeftCount(),
		BagNumItem:             p.bagItemLeftCount(),
		FateLevel:              p.fateLevel,
		FateExp:                p.fateExp,
		Status:                 int(p.curState),
		TrialTowerLevel:        p.InstanceData.TrialTowerMap.BattleId % 1000,
		TotalCombat:            p.totalCombat,
		GuildId:                p.guildId,
		ServerId:               p.serverID,
		AccountName:            p.accountName,
		CreateDay:              p.createDay,
		RegisterAt:             p.registerAt,
		ZoneId:                 p.getZoneId(),
		Guidance:               p.Guidance,
		LotteryNum:             p.getUpRecord(constant.LotteryNum),
		RefineNum:              p.getUpRecord(constant.RefineNum),
		InstanceNum:            p.getUpRecord(constant.InstanceNum), // 活动次数
		CumulativeChargeAmount: p.getUpRecord(constant.CumulativeChargeAmount),
		PinkEntryNum:           p.calcPinkEntryNum(),
		RedPartnerStarNum:      p.calcRedPartnerStarNum(),
		EquipmentStarNum:       p.calcEquipmentStarNum(),
		PetSynthesisNum:        p.calcPetSynthesisNum(),
		SoulStoneLevelNum:      p.calcSoulStoneLevelNum(),
		FiveElementLevelNum:    p.calcFiveElementLevelNum(),
		MuteUntil:              p.muteUntil,
		IsOnline:               int(p.curState),
	}
	return userInfo, nil
}

// 个人信息
func (m *Manager) SetGuidance(uid int64, setp int32) error {
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return errutil.ErrPlayerNotFound
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Guidance = []int32{setp}
	return nil
}

// 个人信息
func (m *Manager) SetPlayerInfo(uid int64, fateLevel int32, fateExp int64) error {
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return errutil.ErrPlayerNotFound
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.SetPlayerInfo(fateLevel, fateExp)
	return nil
}

// 设置命格属性
func (m *Manager) userFate(uid int64, cfgId int32, num int32) ([]int, int32, error) {
	// 手动使用
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return nil, num, errutil.ErrPlayerNotFound
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	res, resnum, err := p.IsFate(cfgId, num)
	if err != nil {
		return nil, resnum, err
	}
	return res, resnum, nil
}

func (m *Manager) Notification(uid int64, Type pb.Notification) {
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return
	}
	packet := &pb.SystemNotificationPush{
		Type: Type,
	}
	p.pushMsg(protocol.RouteSystemNotificationPush, packet)
}

// 获取宠物空闲位置数量
func (m *Manager) GetPetEmptySlot(uid int64) (int, error) {
	p, ok := m.player(uid)
	if !ok || p.session == nil {
		return 0, errutil.ErrPlayerNotFound
	}
	emptySlot := p.PetMgr.getEmptyPetSlotNum()
	if emptySlot == 0 {
		return 0, errutil.ErrPetFull
	}
	return int(emptySlot), nil
}

// 合成
func (m *Manager) Compose(s *session.Session, req *pb.ComposeRequest) error {
	res := &pb.ComposeResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	newComposeId, err := p.compose(req.Info)
	log.Debug("err", err)
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
		res.NewConpose = newComposeId
		return s.Response(res)
	}
	res.NewConpose = newComposeId
	return s.Response(res)
}

// 穿脱宝石
func (m *Manager) GemstoneChange(s *session.Session, req *pb.GemstoneChangeRequest) error {
	res := &pb.GemstoneChangeResponse{OpType: req.OpType, Slot: req.Slot, PartnerSlotId: req.PartnerSlotId}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleEquipSoul) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	if err = p.genstoneChange(req.OpType, req.Slot, req.GemstoneSlot, req.ItemUid, req.PartnerSlotId); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}
	res.GemstoneSlot = req.GemstoneSlot
	res.ItemUid = req.ItemUid

	p.tasks.dispatch(pb.TargetType_MAIN_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST, 0, p.GetMainRoleAllSlotMinSoulStoneLevel())
	p.tasks.dispatch(pb.TargetType_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST, 0, p.GetAllSlotMinSoulStoneLevel())

	return s.Response(res)
}

// 一键穿脱宝石
func (m *Manager) GemstoneChangeAll(s *session.Session, req *pb.GemstoneChangeAllRequest) error {
	res := &pb.GemstoneChangeAllResponse{Code: 0, Error: "", Slot: req.Slot, PartnerSlotId: req.PartnerSlotId}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.GemStoneSlot = nil
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleEquipSoul) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	newdate, err := p.genstoneChangeAll(req.Slot, req.PartnerSlotId)
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
		res.GemStoneSlot = nil
		return s.Response(res)
	}
	res.GemStoneSlot = newdate

	p.tasks.dispatch(pb.TargetType_MAIN_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST, 0, p.GetMainRoleAllSlotMinSoulStoneLevel())
	p.tasks.dispatch(pb.TargetType_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST, 0, p.GetAllSlotMinSoulStoneLevel())

	return s.Response(res)
}

// 跑马灯
func (m *Manager) Lamp(id int32, serverId int32, context []string) {
	packet := &pb.LampPush{
		Uid:     id,
		Context: context,
	}
	if cfg, ok := global.LampCfg[id]; ok {
		if cfg.Receptor == 1 {
			if serverId == 0 {
				//m.group.Broadcast(protocol.RouteLampPush, packet)
				playerMgr.BroadcastToAllGroups(protocol.RouteLampPush, packet)
			} else {
				//for _, v := range m.players {
				//	if v.serverID == serverId {
				//		v.pushMsg(protocol.RouteLampPush, packet)
				//	}
				//}
				playerMgr.BroadcastToGroupByServerID(serverId, protocol.RouteLampPush, packet)
			}
		} else if cfg.Receptor == 3 { // 战区跑马灯
			if serverInfo, err := global.ServerListMgrInstance.Get(int(serverId)); err == nil {
				//for _, p := range m.players {
				//	for _, v := range serverlist {
				//		if v == int(p.serverID) {
				//			p.pushMsg(protocol.RouteLampPush, packet)
				//			break
				//		}
				//	}
				//}
				serverlist := global.ServerListMgrInstance.GetZoneServerIds(serverInfo.ZoneId)
				for _, v := range serverlist {
					playerMgr.BroadcastToGroupByServerID(int32(v), protocol.RouteLampPush, packet)
				}
			}
		}
	}
}

// 跑马灯
func (m *Manager) SystemLamp(rtype int32, rvalue int64, context string) {
	packet := &pb.ChatPush{Type: 18, Context: context}
	switch rtype {
	case 1:
	case 2:
		//战区
	case 3:
		//服务器
		//for _, v := range m.players {
		//	if v.serverID == int32(rvalue) {
		//		v.pushMsg(protocol.RouteLampPush, packet)
		//	}
		//}
		playerMgr.BroadcastToGroupByServerID(int32(rvalue), protocol.RouteLampPush, packet)
	case 4:
		//个人
		//for k, v := range m.players {
		//	if k == int64(rvalue) {
		//		v.pushMsg(protocol.RouteLampPush, packet)
		//	}
		//}
		if player, ok := playerMgr.GetPlayer(rvalue); ok {
			player.pushMsg(protocol.RouteLampPush, packet)
		}
	default:
	}

}
func (m *Manager) SendLamp(lamp *pb.LampPush) {
	packet := &pb.LampPush{
		Uid:     lamp.Uid,
		Context: lamp.Context,
	}
	if cfg, ok := global.LampCfg[lamp.Uid]; ok {
		if cfg.Receptor == 1 {
			//m.group.Broadcast(protocol.RouteLampPush, packet)
			playerMgr.BroadcastToAllGroups(protocol.RouteLampPush, packet)
		}
	}
}

func (m *Manager) GmCmd(cmd protocol.GmCmd) {
	switch cmd.Cmd {

	case "shutdown":
		global.LoginOpen = false

		for _, r := range roomManager.rooms {
			for _, d := range r.desks {
				d.setStatus(constant.DeskStatusDestory)
			}
		}

		time.Sleep(3 * time.Second)

		g := new(errgroup.Group)
		g.SetLimit(20)

		playerMgr.IteratePlayers(func(val *Player) bool {
			logger.Infof("服务器关闭, 保存玩家数据 %d : %s", val.uid, val.name)

			g.Go(func() error {
				val.syncSaveDb()
				return nil
			})
			return true
		})

		if err := g.Wait(); err == nil {
			logger.Info("全部玩家保存完毕.")
		}

		logger.Info("结束进程")

	case "login":
		if len(cmd.Param) == 1 {
			if cmd.Param[0] == "on" {
				global.LoginOpen = true
			} else if cmd.Param[0] == "off" {
				global.LoginOpen = false
			}

			logger.Debugf("GM设置服务器登录状态 %t ", global.LoginOpen)
		}

	case "mute":
		// 禁言指令: mute <uid> <duration>
		// duration: 禁言时长（秒），0表示解除禁言
		if len(cmd.Param) < 2 {
			logger.Errorf("GM指令错误: mute 需要2个参数 <uid> <duration>")
			break
		}
		uid, err := strconv.ParseInt(cmd.Param[0], 10, 64)
		if err != nil {
			logger.Errorf("GM指令错误: mute uid参数错误 %s", cmd.Param[0])
			break
		}
		duration, err := strconv.ParseInt(cmd.Param[1], 10, 64)
		if err != nil {
			logger.Errorf("GM指令错误: mute duration参数错误 %s", cmd.Param[1])
			break
		}

		var muteUntil int64
		if duration > 0 {
			muteUntil = time.Now().Unix() + duration
		} else {
			muteUntil = 0 // 0表示解除禁言
		}

		if err := m.SetMuteUntil(uid, muteUntil); err != nil {
			logger.Errorf("GM指令执行失败: mute uid=%d, duration=%d, error=%v", uid, duration, err)
		} else {
			logger.Infof("GM指令执行成功: mute uid=%d, muteUntil=%d", uid, muteUntil)
		}

	case "ban":
		// 封号指令: ban <uid> <duration>
		// duration: 封号时长（秒），0表示解除封号
		if len(cmd.Param) < 2 {
			logger.Errorf("GM指令错误: ban 需要2个参数 <uid> <duration>")
			break
		}
		uid, err := strconv.ParseInt(cmd.Param[0], 10, 64)
		if err != nil {
			logger.Errorf("GM指令错误: ban uid参数错误 %s", cmd.Param[0])
			break
		}
		duration, err := strconv.ParseInt(cmd.Param[1], 10, 64)
		if err != nil {
			logger.Errorf("GM指令错误: ban duration参数错误 %s", cmd.Param[1])
			break
		}

		var banUntil int64
		if duration > 0 {
			banUntil = time.Now().Unix() + duration
		} else {
			banUntil = 0 // 0表示解除封号
		}

		if err := m.SetBanUntil(uid, banUntil); err != nil {
			logger.Errorf("GM指令执行失败: ban uid=%d, duration=%d, error=%v", uid, duration, err)
		} else {
			logger.Infof("GM指令执行成功: ban uid=%d, banUntil=%d", uid, banUntil)
		}

	case "delete":
		// 删除/恢复玩家指令: delete <uid> <action>
		// action: 1=删除（软删除），2=恢复
		if len(cmd.Param) < 2 {
			logger.Errorf("GM指令错误: delete 需要2个参数 <uid> <action>")
			break
		}
		uid, err := strconv.ParseInt(cmd.Param[0], 10, 64)
		if err != nil {
			logger.Errorf("GM指令错误: delete uid参数错误 %s", cmd.Param[0])
			break
		}
		action, err := strconv.ParseInt(cmd.Param[1], 10, 32)
		if err != nil {
			logger.Errorf("GM指令错误: delete action参数错误 %s", cmd.Param[1])
			break
		}

		if err := m.DeleteOrRestorePlayer(uid, int32(action)); err != nil {
			logger.Errorf("GM指令执行失败: delete uid=%d, action=%d, error=%v", uid, action, err)
		} else {
			actionStr := "删除"
			if action == 2 {
				actionStr = "恢复"
			}
			logger.Infof("GM指令执行成功: %s uid=%d", actionStr, uid)
		}
	}
}

// 伙伴
func (m *Manager) PartnerList(s *session.Session, req *pb.PartnerListRequest) error {
	res := &pb.PartnerListResponse{List: []*pb.PartnerInfo{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	for _, v := range p.partners {
		partner := &pb.PartnerInfo{
			Uid:   v.uid,
			CfgId: v.cfgId,
			Level: v.level,
			Star:  v.star,
		}
		partner.Skills = []*pb.SkillInfo{}
		for _, v := range v.skills {
			partner.Skills = append(partner.Skills, &pb.SkillInfo{CfgId: v.CfgId, Cooldown: 0, Advanced: v.AdvancedLevel})
		}

		partner.Points = v.points
		for _, v2 := range v.activated {
			partner.Talents = append(partner.Talents, &pb.KeyValue{K: v2.ID, V: v2.Num})
		}
		res.List = append(res.List, partner)
	}
	return s.Response(res)
}

func (m *Manager) PartnerSlotInfo(s *session.Session, req *pb.PartnerSlotInfoRequest) error {
	res := &pb.PartnerSlotInfoResponse{List: []*pb.PartnerSlotInfo{}}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	for _, v := range p.partnerSlots {
		res.List = append(res.List, &pb.PartnerSlotInfo{
			SlotId:     v.SlotId,
			PartnerUid: v.PartnerUid,
			Equips:     p.FillEquipInfo(v.Equips),
		})
	}

	return s.Response(res)
}

// 首页宝箱设置
func (m *Manager) HomeBoxSetting(s *session.Session, req *pb.HomeBoxSettingRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_HomeBoxSetting)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	err = p.homeBoxSetting(req.Level, req.Job)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

// 首页开宝箱
func (m *Manager) HomeOpenBox(s *session.Session, req *pb.HomeOpenBoxRequest) error {
	res := &pb.HomeOpenBoxResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	// 判断之前开的宝箱有没有数据
	if len(p.box.NowOpenRes) > 0 || req.Num <= 0 {
		res.Code = -2
		res.Error = errutil.ErrOpenBoxEnough.Error()
		return s.Response(res)
	}
	if !p.checkTakeSource(constant.ItemTypeTreasureNum, int64(req.Num)) {
		res.Code = -3
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	if req.Level > p.level {
		res.Code = -4
		res.Error = errutil.ErrMaxLevel.Error()
		return s.Response(res)
	}
	openbox, err := p.openboxV2(req.Level, req.Num, req.Job)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	p.takeSource(constant.ItemTypeTreasureNum, int64(req.Num), constant.ItemDelHomeOpenBox)
	res.ItemInfo = openbox
	return s.Response(res)
}

// 首页宝箱升级
func (m *Manager) HomeUpgradeBox(s *session.Session, req *pb.HomeUpgradeBoxRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_HomeUpgradeBox)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	err = p.upgradeBox()
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

// 加快首页宝箱升级
func (m *Manager) HomeUpgradeBoxQuick(s *session.Session, req *pb.HomeUpgradeBoxQuickRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_HomeUpgradeBoxQuick)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	if req.Num <= 0 {
		res.Code = -2
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}
	err = p.UpgradeBoxQuick(req.Num)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

// 装备入背包
func (m *Manager) BoxResIntoBag(s *session.Session, req *pb.BoxResIntoBagRequest) error {
	res := &pb.BoxResIntoBagResponse{}
	res.Uid = req.Uid
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	err = p.boxResIntoBag(req.Uid)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

func (m *Manager) UpgradeLevel(s *session.Session, req *pb.UpgradeLevelRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_UpgradeLevel)}

	// p, err := playerWithSession(s)
	// if err != nil {
	// 	res.Code = -1
	// 	res.Error = err.Error()
	// 	return s.Response(res)
	// }
	// if req.Type == pb.ObjectType_PLAYER {
	// 	p.levelUp()
	// } else {
	// 	partner := p.getPartnerInfo(req.Uid)
	// 	if partner == nil {
	// 		res.Code = -2
	// 		return s.Response(res)
	// 	}
	// 	if partner.level >= p.level {
	// 		res.Code = -3
	// 		res.Error = errutil.ErrMaxLevel.Error()
	// 		return s.Response(res)
	// 	}
	// 	if !partner.levelUp() {
	// 		res.Code = -3
	// 		return s.Response(res)
	// 	}
	// }
	// p.syncUserBase()

	return s.Response(res)
}

// /伙伴
func (m *Manager) PartnerCombat(s *session.Session, req *pb.PartnerCombatRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerCombat)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	if p.desk != nil {
		res.Code = -2
		res.Error = errutil.IllegalInBattle
		return s.Response(res)
	}

	if req.SlotId < 1 || req.SlotId > 4 {
		res.Code = -3
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}
	needLevel := int32(0)
	needOpenDay := int32(0)
	serverDay := int32(0)
	for _, v := range global.GlobalCfg.Openslot.IDLevel {
		if v.ID == req.SlotId {
			serverDay = v.ServerDay
			needLevel = v.Needlevel
			needOpenDay = v.Openday
			break
		}
	}
	if p.level < needLevel {
		res.Code = -4
		res.Error = errutil.ErrLevel.Error()
		return s.Response(res)
	}

	if serverDay > 0 {
		if global.ServerListMgrInstance.GetServerOpenDay(p.serverID) < needOpenDay {
			res.Code = -5
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
	} else {
		if p.createDay < needOpenDay {
			res.Code = -5
			res.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(res)
		}
	}

	if err := p.partnerCombat(req.SlotId, req.PartnerUid); err != nil {
		res.Code = -5
		res.Error = err.Error()
	}
	return s.Response(res)
}

func (m *Manager) PartnerRest(s *session.Session, req *pb.PartnerRestRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerRest)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if p.desk != nil {
		res.Code = -2
		res.Error = errutil.IllegalInBattle
		return s.Response(res)
	}
	if err := p.partnerRest(req.SlotId); err != nil {
		res.Code = -3
		res.Error = err.Error()
	}
	return s.Response(res)
}

func (m *Manager) PartnerReplace(s *session.Session, req *pb.PartnerReplaceRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerReplace)}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if req.SlotId < 1 || req.SlotId > 4 {
		res.Code = -2
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}

	needLevel := int32(0)
	for _, v := range global.GlobalCfg.Openslot.IDLevel {
		if v.ID == req.SlotId {
			needLevel = v.Needlevel
			break
		}
	}

	if p.level < needLevel {
		res.Code = -3
		res.Error = errutil.ErrLevel.Error()
		return s.Response(res)
	}

	// 下阵伙伴
	if p.desk != nil {
		res.Code = -4
		res.Error = errutil.IllegalInBattle
		return s.Response(res)
	}

	slotInfo := p.getPartnerSlotInfo(req.SlotId)
	if slotInfo == nil {
		return errutil.ErrPartnerInvalid
	}
	partner := p.getPartnerInfoBySlotId(req.SlotId)
	if partner == nil {
		return errutil.ErrPartnerInvalid
	}

	equipnum := 0
	for _, v := range slotInfo.Equips {
		if v.ItemId > 0 {
			equipnum++
		}
	}

	if equipnum > p.bagEquipLeftCount() {
		res.Code = -5
		res.Error = errutil.ErrBagFull.Error()
		return s.Response(res)
	}

	partnerSlot := p.getPartnerSlotInfo(req.SlotId)
	if partnerSlot == nil {
		res.Code = -6
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}

	// p.partnerSlots[slotId].PartnerUid = 0
	for k, v := range partnerSlot.Equips {
		if v.ItemId <= 0 {
			continue
		}
		p.changeEquip(2, pb.EquipSlot(k), v.ItemId, partnerSlot.SlotId)
	}

	// 自动下阵伙伴技能
	bskillslotchange := false
	for _, v := range partner.skills {
		for slot := 0; slot < len(p.skillsSlot); slot++ {
			if p.skillsSlot[slot] == v.CfgId {
				p.skillsSlot[slot] = 0
				bskillslotchange = true
			}
		}
	}
	if bskillslotchange {
		p.syncUserSkillSlot()
	}

	partnerSlot.PartnerUid = 0
	p.SyncUserPartnerSlot()

	// 上阵伙伴
	if err := p.partnerCombat(req.SlotId, req.PartnerUid); err != nil {
		res.Code = -7
		res.Error = err.Error()
	}
	return s.Response(res)
}

func (m *Manager) PartnerStar(s *session.Session, req *pb.PartnerStarRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerStar)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	partner := p.getPartnerInfo(req.PartnerUid)
	if partner == nil {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}

	if err = partner.starUp(); err != nil {
		res.Code = -3
		res.Error = err.Error()
		return s.Response(res)
	}
	if err = p.syncSavePartnerDb(partner.uid, model.Update); err != nil {
		logger.Errorf("伙伴更新数据失败 paruid :%d err:%s", partner.uid, err.Error())
	}

	allPartnerMinStar := p.GetMinStarByQuality(constant.Quality(partner.cfg.Quality))
	switch constant.Quality(partner.cfg.Quality) {
	case constant.QualityGreen:
		p.tasks.dispatch(pb.TargetType_ALL_GREEN_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	case constant.QualityBlue:
		p.tasks.dispatch(pb.TargetType_ALL_BLUE_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	case constant.QualityPurple:
		p.tasks.dispatch(pb.TargetType_ALL_PURPLE_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	// case constant.PartnerQualityPink:
	// 	p.tasks.dispatch(pb.TargetType_ALL_PINK_PARTNERS_REACH_STAR, 0, 0)
	case constant.QualityOrange:
		p.tasks.dispatch(pb.TargetType_ALL_ORANGE_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	case constant.QualityRed:
		p.tasks.dispatch(pb.TargetType_ALL_RED_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	case constant.QualityYellow:
		p.tasks.dispatch(pb.TargetType_ALL_YELLOW_PARTNERS_REACH_STAR, 0, allPartnerMinStar)
	default:
		logger.Warnf("PartnerStar 未知伙伴品质 star up quality:%d", partner.cfg.Quality)
	}

	return s.Response(res)
}

func (m *Manager) PartnerCompose(s *session.Session, req *pb.PartnerComposeRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerCompose)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if err = p.partnerCompose(int32(req.PartnerCfgId)); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}

	return s.Response(res)
}

// PartnerReforge
func (m *Manager) PartnerReforge(s *session.Session, req *pb.PartnerReforgeRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerReforge)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if err = p.partnerReforge(req.PartnerUid, req.IsFree); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}

	return s.Response(res)
}

// PartnerDisassemble
func (m *Manager) PartnerDisassemble(s *session.Session, req *pb.PartnerDisassembleRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerDisassemble)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if err = p.partnerDisassemble(req.PartnerUid); err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}

	return s.Response(res)
}

// PartnerTalentActivate
func (m *Manager) PartnerTalentActivate(s *session.Session, req *pb.PartnerTalentActivateRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerTalentActivate)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	par := p.getPartnerInfo(req.PartnerUid)
	if par == nil {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}

	if err := par.talentActivate(req.TalentSlot); err != nil {
		res.Code = -3
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

// PartnerTalentReset
func (m *Manager) PartnerTalentReset(s *session.Session, req *pb.PartnerTalentResetRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PartnerTalentReset)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}
	par := p.getPartnerInfo(req.PartnerUid)
	if par == nil {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}

	if err := par.talentReset(); err != nil {
		res.Code = -3
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

func (m *Manager) EquipBest(s *session.Session, req *pb.EquipBestRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_EquipBest)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}

	if req.Type == pb.ObjectType_PARTNER && req.PartnerSlotId == 0 {
		res.Code = -2
		res.Error = errutil.ErrPartnerInvalid.Error()
		return s.Response(res)
	}
	p.equipBest(req.Type, req.PartnerSlotId)
	return s.Response(res)
}

// 批量添加伙伴
func (m *Manager) addPartnerBatch(AddSource *protocol.AddSource) error {
	p, ok := m.player(AddSource.Uid)
	if !ok {
		return errutil.ErrPlayerNotFound
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, v := range AddSource.Source {
		if err := p.addPartner(v.Id); err != nil {
			return errutil.ErrPartnerInvalid
		}
	}
	return nil
}

// 批量添加宠物
func (m *Manager) addPetBatch(AddSource *protocol.AddSource, NeedDel int32, reason string) error {
	p, ok := m.player(AddSource.Uid)
	if !ok {
		return errutil.ErrPlayerNotFound
	}
	for _, v := range AddSource.Source {
		p.tasks.dispatch(pb.TargetType_TOTAL_PET_EGGS_OPENED, 0, 1)
		if _, err := p.addPet(v.Id, 0, 0, 0, NeedDel, reason); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) PetStar(s *session.Session, req *pb.PetStarRequest) error {
	res := &pb.PetStarResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
	}
	if req.Num <= 0 {
		res.Code = -2
		res.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(res)
	}
	resnum, err := p.petStar(req.SlotId, req.CfgId, req.Num)
	if err != nil {
		res.Code = -3
		res.Error = err.Error()
		return s.Response(res)
	}
	res.ResPer = resnum
	return s.Response(res)
}

func (m *Manager) PetLevel(s *session.Session, req *pb.PetLevelRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_PetLevel)}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	pet := p.getPetInfo(req.SlotId)
	if pet == nil {
		res.Code = -2
		res.Error = errutil.ErrPetInvalid.Error()
		return s.Response(res)
	}

	if err = p.petLevel(req.SlotId); err != nil {
		res.Code = -3
		res.Error = err.Error()
		return s.Response(res)
	}
	return s.Response(res)
}

func (m *Manager) PetMove(s *session.Session, req *pb.PetMoveRequest) error {
	res := &pb.PetMoveResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModulePet) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	newSlot, newPet, newPetCfgId, err := p.petMove(req.SlotId, req.TargetSlotId)
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
	}
	res.PetcfgId = newPet
	res.SlotId = newSlot
	res.NewPetCfgId = newPetCfgId
	return s.Response(res)
}

func (m *Manager) InstanceInfo(s *session.Session, req *pb.InstanceInfoRequest) error {
	res := &pb.InstanceInfoResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		return s.Response(res)
	}

	p.fillInstanceInfo(res)

	logger.Debugf("InstanceInfoResponse.InstanceInfoResponse_AFK.AfkMultiplierID: %d, uid: %d", res.GetAfk().GetAfkMultiplierID(), p.Uid())

	return s.Response(res)
}

func (m *Manager) SendMail(uid int64, desc string, source string, attachment []protocol.Attachment) {
	mail := protocol.MailAddRequest{
		Uid:          -1,
		RecipientUid: uid,
		Title:        "系统邮件",
		Desc:         desc,
		ChannelCode:  source,
		Attachment:   attachment,
	}
	if err := provide.Mailadd(&mail); err != nil {
		if p, ok := m.player(uid); ok {
			packet := &pb.SystemNotificationPush{Type: pb.Notification_MAIL_FULL}
			p.pushMsg(protocol.RouteSystemNotificationPush, packet)
		}
	}
}

func (m *Manager) ObservePlayerInfo(s *session.Session, req *pb.ObservePlayerInfoRequest) error {
	res := &pb.ObservePlayerInfoResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		return s.Response(res)
	}
	// 获取观察玩家信息
	observePlayer := m.ObserveUserBaseInfo(p.serverID, req.Uid)
	if observePlayer == nil {
		return s.Response(res)
	}

	// Fill response with observe player info
	res.ObservePlayerInfo = &pb.ObservePlayerInfo{
		TotalCombat:   observePlayer.totalCombat, // Using combat as total combat since it's the main combat value
		Combat:        observePlayer.combat,
		VipLevel:      int32(observePlayer.vip),
		Sex:           int32(observePlayer.sex),
		Level:         int32(observePlayer.level),
		Job:           int32(observePlayer.job),
		Nickname:      observePlayer.name,
		Attributes:    []*pb.Stats{},
		PartnersStats: make([]*pb.ObservePlayerInfo_PartnerStats, 0),
		Partners:      make([]*pb.PartnerSlotInfo, 0),
		//Pets: &pb.PetSlotInfo{
		//	CfgId:   observePlayer.pet[0].PetCfgId,
		//	Level:   int32(observePlayer.pet[0].Level),
		//	Star:    int32(observePlayer.pet[0].Star),
		//	Enhance: int32(observePlayer.pet[0].Enhance),
		//	Quality: int32(observePlayer.pet[0].Quality),
		//	Exp:     int32(observePlayer.pet[0].Exp),
		//	ExpStar: int32(observePlayer.pet[0].ExpStar),
		//},
		Pets:      observePlayer.PetMgr.GetBattlePetInfo(),
		SlotInfo:  make([]*pb.SlotInfo, 0),
		GuildName: observePlayer.guildName,
		Avatar:    observePlayer.avatar.currentAvatarID,
		Frame:     observePlayer.avatar.currentFrameID,
	}

	observePlayer.stats.Range(func(k pb.Attribute, v int32) bool {
		if v <= 0 {
			return true
		}
		res.ObservePlayerInfo.Attributes = append(res.ObservePlayerInfo.Attributes, &pb.Stats{Attr: k, Num: v})
		return true
	})
	// Add partners stats and info
	for _, partner := range observePlayer.partners {
		// Add partner stats
		partnerStats := &pb.ObservePlayerInfo_PartnerStats{
			Combat:     partner.combat,
			Attributes: []*pb.Stats{},
			PartnerUid: partner.uid,
		}
		partner.stats.Range(func(k pb.Attribute, v int32) bool {
			if v <= 0 {
				return true
			}
			partnerStats.Attributes = append(partnerStats.Attributes, &pb.Stats{Attr: k, Num: v})
			return true
		})
		res.ObservePlayerInfo.PartnersStats = append(res.ObservePlayerInfo.PartnersStats, partnerStats)

		res.ObservePlayerInfo.PartnerInfo = append(res.ObservePlayerInfo.PartnerInfo, &pb.PartnerInfo{
			Uid:   partner.uid,
			CfgId: partner.cfgId,
			Level: partner.level,
			Star:  partner.star,
		})
		// Add partner slot info

	}
	for _, slot := range observePlayer.partnerSlots {
		partnerSlotInfo := &pb.PartnerSlotInfo{
			SlotId:     slot.SlotId,
			PartnerUid: slot.PartnerUid,
			Equips:     observePlayer.FillEquipInfo(slot.Equips),
		}
		res.ObservePlayerInfo.Partners = append(res.ObservePlayerInfo.Partners, partnerSlotInfo)
	}

	// Add  slot info

	for _, v := range observePlayer.equips {
		Gemstone := make([]int32, 0)
		for _, v2 := range v.Gemstone {
			Gemstone = append(Gemstone, int32(v2))
		}
		slot := &pb.SlotInfo{Stars: v.Stars, Enhance: v.Enhance, Item: observePlayer.getItemByUid(v.ItemId), Gemstone: Gemstone}
		res.ObservePlayerInfo.SlotInfo = append(res.ObservePlayerInfo.SlotInfo, slot)
	}

	characterRefineData, _ := observePlayer.refine.ToPb()
	res.ObservePlayerInfo.CharacterRefines = characterRefineData

	// 官职
	res.ObservePlayerInfo.OfficialExp = p.tokens[global.GlobalCfg.BaseValue.Official.ItemID]

	return s.Response(res)
}

// 获取某个人的信息
func (m *Manager) ObserveUserBaseInfo(serverId int32, uid int64) *Player {
	// 获取用户信息
	// 检查玩家是否在线 defaultManager.players(uid)

	if p, ok := defaultManager.player(uid); ok {
		return p
	}

	if r := robotMgr.getRankRobotInfo(serverId, uid); r != nil {
		if r.robotType == RobotTypePlayer {
			roleCfg, err := global.GetRoleCfg(int32(r.job), r.level)
			if err != nil {
				logger.Errorf("获取角色等级配置失败 ")
			}

			for _, rc := range roleCfg.BassProperty {
				r.stats.setStats(pb.Attribute(rc.Key), rc.Value)
			}

			r.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(r.job)).AttackSpeed)
			r.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(r.job)).MoveSpeed)
		}
		return r.Player
	}
	p := newPlayer(nil, uid, "", "", "")
	if !p.syncFromDB() {
		return nil
	}
	p.curState = PlayerStatusOffLine
	p.checkPartnerBonds()
	p.refreshStats(WithRefreshStatsPushMessage(false))
	defaultManager.setPlayer(uid, p)
	// 获取用户信息
	return p
}

func (m *Manager) RankLogin(serverID int32, uid int64) {
	// 获取观察玩家信息
	observePlayer := m.ObserveUserBaseInfo(serverID, uid)
	if observePlayer == nil {
		return
	}
	rank := provide.IsFirst(uid, 0, observePlayer.serverID)
	if len(rank) > 0 {
		for _, v := range rank {
			if v == constant.RankTypeGightEndlayers {
				continue
			}
			//有时间的榜单，超过时间直接不录入
			send := true
			if constant.RankTypeServerOpenDay[int64(v)] {
				openDay := global.ServerListMgrInstance.GetServerOpenDay(int32(observePlayer.serverID))
				//天数
				for _, v2 := range global.GlobalCfg.LimitRanking.Key {
					//实际天数与配置天数一样，需要发奖励
					if v2.Value < int64(openDay) && int64(v2.ID) == int64(v) {
						send = false
						break
					}
				}
			}
			if !send {
				continue
			}
			strSend := []string{constant.RankTypeName[v], observePlayer.name}
			// 发送消息给所有在线用户
			defaultManager.Lamp(int32(constant.LampTopRankOnline), observePlayer.serverID, strSend)
		}
	}
}

func (m *Manager) GetPlayer(uid int64) (*Player, error) {
	if p, ok := defaultManager.player(uid); ok {
		return p, nil
	}
	return nil, errutil.ErrPlayerNotFound
}

// SetMuteUntil 设置玩家禁言到期时间（支持在线玩家）
// muteUntil: 禁言到期时间（Unix时间戳），0表示解除禁言
func (m *Manager) SetMuteUntil(uid int64, muteUntil int64) error {
	// 如果玩家在线，直接更新内存数据即可，玩家保存时会自动同步到数据库
	if p, ok := m.player(uid); ok {
		p.muteUntil = muteUntil
		// 异步保存到数据库（确保数据一致性）
		return nil
	}

	// 玩家不在线，直接操作数据库
	if err := db.UpdateUserMuteUntil(uid, muteUntil); err != nil {
		return err
	}

	return nil
}

// SetBanUntil 设置账号封号到期时间（封号是针对账号的，不是针对角色的）
// banUntil: 封号到期时间（Unix时间戳），0表示解除封号
func (m *Manager) SetBanUntil(uid int64, banUntil int64) error {
	// 更新 Register 表中的封号状态
	if err := db.UpdateRegisterBanUntil(uid, banUntil); err != nil {
		return err
	}

	// 如果玩家在线且被封号，立即踢下线
	if banUntil > 0 && time.Now().Unix() < banUntil {
		if p, ok := m.player(uid); ok {
			// 告知用户被封号
			packet := &pb.CompulsoryWithdrawalPush{
				Status: pb.CompulsoryWithdrawalPush_WithdrawalStatus(pb.CompulsoryWithdrawalPush_ForceOffline),
			}
			p.pushMsg(protocol.RouteCompulsoryWithdrawalPush, packet)
			// 移除广播频道
			playerMgr.RemoveSessionFromGroup(p.session)

			// 如果玩家在房间，则退出房间
			if p.desk != nil && p.desk.group != nil {
				p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: false}
				p.desk = nil
				p.setPlayerStatus(PlayerStatusDisconnect)
			}

			// 延迟关闭session，确保消息能够发送
			time.AfterFunc(time.Millisecond*10, func() {
				p.session.Clear()
				p.session.Close()
			})
		}
	}

	return nil
}

// DeleteOrRestorePlayer 删除/恢复玩家（软删除，支持恢复）
// action: 1=删除（软删除），2=恢复
func (m *Manager) DeleteOrRestorePlayer(uid int64, action int32) error {
	// 如果玩家在线，先踢下线
	if p, ok := m.player(uid); ok {
		if p.session != nil {
			// 告知用户被踢下线
			packet := &pb.CompulsoryWithdrawalPush{
				Status: pb.CompulsoryWithdrawalPush_WithdrawalStatus(pb.CompulsoryWithdrawalPush_ForceOffline),
			}
			p.pushMsg(protocol.RouteCompulsoryWithdrawalPush, packet)
			// 移除广播频道
			playerMgr.RemoveSessionFromGroup(p.session)

			// 如果玩家在房间，则退出房间
			if p.desk != nil && p.desk.group != nil {
				p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: false}
				p.desk = nil
				p.setPlayerStatus(PlayerStatusDisconnect)
			}

			// 延迟关闭session，确保消息能够发送
			time.AfterFunc(time.Millisecond*10, func() {
				p.session.Clear()
				p.session.Close()
			})
		}
	}

	// 更新数据库中的玩家状态（软删除或恢复）
	if err := db.DeleteOrRestoreUser(uid, action); err != nil {
		return fmt.Errorf("delete or restore user failed: %w", err)
	}

	return nil
}

// DeletePlayer 删除玩家（先踢下线，再删除数据库）- 保留兼容性
func (m *Manager) DeletePlayer(uid int64) error {
	return m.DeleteOrRestorePlayer(uid, 1) // 默认删除
}

// func (m *Manager) ActivityRewardReceive(s *session.Session, req *pb.ActivityRewardReceiveRequest) error {
// 	resp := &pb.ActivityRewardReceiveResponse{}
// 	p, err := playerWithSession(s)
// 	if err != nil {
// 		resp.Code = -1
// 		resp.Error = err.Error()
// 		return s.Response(resp)
// 	}

// 	p.activityMgr.RewardReceive(s.UID(), p, req.GetRewards())

// 	return s.Response(resp)
// }

// 新手引导
func (m *Manager) GuidanceUpdate(s *session.Session, req *pb.GuidanceUpdateRequest) error {
	resp := &pb.GuidanceUpdateResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if len(p.Guidance) > 0 && len(req.Guidance) > 0 {
		if p.Guidance[0] > req.Guidance[0] {
			resp.Guidance = p.Guidance
			return s.Response(resp)
		}
		// m.guideCompletedShushuDate(p, p.Guidance[0], req.Guidance[0])
	}

	p.Guidance = req.Guidance
	//兼容一下老的引导数据是数组模式
	resp.Guidance = append(resp.Guidance, p.Guidance...)
	return s.Response(resp)
}
func (m *Manager) guideCompletedShushuDate(p *Player, old_guide_step int32, new_guide_step int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//在合并点其他的数据
	shushuDate["guide_old_step"] = old_guide_step
	shushuDate["guide_new_step"] = new_guide_step
	shushuDate["guide_name"] = new_guide_step
	shushuDate["name"] = "guide_completed_admin"
	shushumain.ReportData(shushuDate)
}

func (m *Manager) UpShushuDate(uid int64, date map[string]interface{}, have bool, first bool) {
	shushuDate := make(map[string]interface{}, 0)
	if have {
		p, ok := defaultManager.player(uid)
		if ok {
			p.lock.Lock()
			defer p.lock.Unlock()

			baseInfo := p.ToMapForShushu()
			for k, v := range baseInfo {
				shushuDate[k] = v
			}

		}
		if first {
			getUpRecord := p.getUpRecord(date["name"].(string))
			is_first_time := true
			if getUpRecord > 0 {
				is_first_time = false
			}
			shushuDate["is_first_time"] = is_first_time
			p.addUpRecord(date["name"].(string))
		}
	}
	for k, v := range date {
		shushuDate[k] = v
	}

	shushumain.ReportData(shushuDate)
}
func (m *Manager) AddUpRecord(uid int64, name string, num int64) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return
	}
	p.addUpRecordByAmount(name, num)
}
func (m *Manager) UpShushuDateSet(uid int64, date map[string]interface{}, have bool) {
	shushuDate := make(map[string]interface{}, 0)
	for k, v := range date {
		shushuDate[k] = v
	}
	shushumain.ReportUserSetData(shushuDate)
}

func (m *Manager) UpShushuDateSetOnce(uid int64, date map[string]interface{}, have bool) {
	shushuDate := make(map[string]interface{}, 0)
	for k, v := range date {
		shushuDate[k] = v
	}
	shushumain.ReportUserSetOnceData(shushuDate)
}

func (m *Manager) UpShushuDateNo(date map[string]interface{}) {
	shushuDate := make(map[string]interface{}, 0)

	for k, v := range date {
		shushuDate[k] = v
	}
	shushumain.ReportData(shushuDate)
}

func (m *Manager) UpShushuTencent(date map[string]interface{}) {
	shushuDate := make(map[string]interface{}, 0)
	for k, v := range date {
		shushuDate[k] = v
	}
	shushumain.ReportTencentData(shushuDate)
}

// 获取背包内道具数量
func (m *Manager) GetItemCount(uid int64, itemUid int32) (int64, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return 0, errutil.ErrPlayerNotFound
	}
	petNormalCount := p.getItemCountInBag(itemUid)
	if petNormalCount < 0 {
		return 0, errutil.ErrItemNotFound
	}
	return petNormalCount, nil
}

// MailCounterSet 设置MailCounter
func (m *Manager) MailCounterSet(uid int64, op constant.MailCounterOPType, num int) error {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return errutil.ErrPlayerNotFound
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.SetMailCounter(op, num)
}

// 同步任务
func (m *Manager) TaskDispatch(uid int64, key pb.TargetType, Typeson int32, limit int32) error {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return errutil.ErrPlayerNotFound
	}
	p.eventMgr.Publish(event.EventTypeTaskProgress, &event.TaskEventData{
		Key: key,
		Id:  Typeson,
		Num: limit,
	})

	return nil
}

// GuardIslandMultiplier 设置守护岛倍数
func (m *Manager) GuardIslandMultiplierSvc(uid, guardIslandRecordID int64) error {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return errutil.ErrPlayerNotFound
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	return p.GuardIslandMultiplier(guardIslandRecordID)
}

// ChangeRoleNameREST 修改角色名
func (m *Manager) ChangeRoleNameREST(uid int64, name string) (string, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return "", errutil.ErrPlayerNotFound
	}
	return p.changeRoleName(name)
}

func (m *Manager) GuildCreateREST(uid int64, req *pb.GuildCreateRequest) (*pb.GuildCreateResponse, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return nil, errutil.ErrPlayerNotFound
	}

	return guildManager.guildCreate(p, req)
}

func (m *Manager) ChatREST(uid int64, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return nil, errutil.ErrPlayerNotFound
	}

	return defaultManager.chat(p, req)
}

// 重置pvp积分
func (m *Manager) resetPvpScore() error {

	logger.Info("每周重置赛季 pvp 分数")

	//for _, v := range defaultManager.players {
	//	v.resetPvpScore()
	//}
	playerMgr.IteratePlayers(func(val *Player) bool {
		val.resetPvpScore()
		return true
	})

	if err := db.UpdateAllPvpRecord(2); err != nil {
		logger.Errorf("重置pvp 分数失败 error: %v", err)
		return err
	}
	return nil
}

func (m *Manager) AddAdCount(uid int64) (int32, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return 0, errutil.ErrPlayerNotFound
	}
	p.AddAdCount()

	return p.GetAdCount(), nil
}

func (m *Manager) ElementInfo(s *session.Session, req *pb.ElementInfoRequest) error {
	resp := &pb.ElementInfoResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	res := p.elementInfo()
	resp = res
	return s.Response(resp)
}

func (m *Manager) ElementUpgrade(s *session.Session, req *pb.ElementUpgradeRequest) error {
	p, err := playerWithSession(s)
	resp := &pb.ElementUpgradeResponse{SlotId: req.SlotId}
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	success, oldLevel, newLevel, main, err := p.elementUpgrade(pb.ElementSlot(req.SlotId))
	resp.OldLevel = oldLevel
	resp.NewLevel = newLevel
	resp.ElementMain = main
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if success {
		resp.Success = 1
	} else {
		resp.Success = 0
	}
	p.refreshStats()
	return s.Response(resp)
}

func (m *Manager) ElementBreakthrough(s *session.Session, req *pb.ElementBreakthroughRequest) error {
	p, err := playerWithSession(s)
	resp := &pb.ElementBreakthroughResponse{SlotId: req.SlotId}
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	oldLevel, newLevel, err := p.ElementBreakthrough(pb.ElementSlot(req.SlotId))
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	resp.Code = 0
	resp.OldLevel = oldLevel
	resp.NewLevel = newLevel
	p.refreshStats()
	return s.Response(resp)
}
func (m *Manager) ElementSweep(s *session.Session, req *pb.ElementSweepRequest) error {
	p, err := playerWithSession(s)
	resp := &pb.CommonResponse{RequestCmdId: uint32(pb.CmdId_C2S_ElementSweep)}
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	err = p.elementSweep()
	if err != nil {
		resp.Code = -2
		resp.Error = err.Error()
		return s.Response(resp)
	}
	resp.Code = 0
	return s.Response(resp)
}
func (m *Manager) ElementSweepForward(s *session.Session, req *pb.ElementSweepForWardRequest) error {
	p, err := playerWithSession(s)
	resp := &pb.ElementSweepForWardResponse{Num: 0}
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if req.Num <= 0 {
		resp.Code = -3
		resp.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(resp)
	}
	if p.ElementMap.ChallengeTime < req.Num {
		resp.Code = -4
		resp.Error = errutil.ErrEnergyFull.Error()
		return s.Response(resp)
	}
	err = p.elementSweepForward(req.Num)
	if err != nil {
		resp.Code = -2
		resp.Error = err.Error()
		return s.Response(resp)
	}
	resp.Code = 0
	resp.Num = req.Num
	return s.Response(resp)
}

// 称号列表 TitleListRequest
func (m *Manager) TitleList(s *session.Session, req *pb.TitleListRequest) error {
	resp := &pb.TitleListResponse{Titles: []*pb.TitleInfo{}}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	for _, v := range p.titles.Titles {
		resp.Titles = append(resp.Titles, v)
	}
	return s.Response(resp)
}

// 称号合成 C2S_TitleCompose
func (m *Manager) TitleCompose(s *session.Session, req *pb.TitleComposeRequest) error {
	resp := &pb.TitleComposeResponse{CfgId: req.CfgId}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if err := p.titleCompose(req.CfgId); err != nil {
		resp.Code = -2
		resp.Error = err.Error()
	}
	return s.Response(resp)
}

// 称号升星 C2S_TitleStart
func (m *Manager) TitleStart(s *session.Session, req *pb.TitleStartRequest) error {
	resp := &pb.TitleStartResponse{CfgId: req.CfgId}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	star, err := p.titleStart(req.CfgId, req.LuckItemCfgId)
	if err != nil {
		resp.Code = -2
		resp.Error = err.Error()
	}
	resp.Star = star
	return s.Response(resp)
}

// 装备称号 C2S_TitleEquip
func (m *Manager) TitleEquip(s *session.Session, req *pb.TitleEquipRequest) error {
	resp := &pb.TitleEquipResponse{CfgId: req.CfgId}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	if !p.isModuleOpen(constant.ModuleTitle) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}
	if req.Slot < 0 || req.Slot > 4 {
		resp.Code = -2
		resp.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(resp)
	}
	err = p.titleEquip(req.CfgId, req.Slot)
	if err != nil {
		resp.Code = -3
		resp.Error = err.Error()
	}

	return s.Response(resp)
}

// 获取抽奖列表
func (m *Manager) LotteryInfo(s *session.Session, req *pb.LotteryInfoRequest) error {
	resp := &pb.LotteryInfoResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	// if !p.isModuleOpen(constant.ModuleLottery) {
	// 	resp.Code = -21
	// 	resp.Error = errutil.ErrModuleNotOpen
	// 	return s.Response(resp)
	// }
	//对比时间戳
	// 检查是否需要刷新
	currentTime := time.Now().Unix()
	shouldRefresh := false

	// 计算今天凌晨12点的时间戳
	todayMidnight := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix()

	// 计算从凌晨12点开始到现在应该经过的半小时数
	halfHoursSinceMidnight := (currentTime - todayMidnight) / (30 * 60)
	// 计算最近一次应该刷新的时间点
	latestRefreshTime := todayMidnight + halfHoursSinceMidnight*(30*60)

	// 如果上次刷新时间早于最近一次应该刷新的时间点，则需要刷新
	if p.LotteryInfo.RefreshTime < latestRefreshTime {
		shouldRefresh = true
		p.LotteryInfo.RefreshTime = latestRefreshTime
		// 这里可以添加刷新奖池的逻辑
		// 例如：重新生成PrizePool内容
	}
	prizePool := make([]*pb.PrizePool, LotteryPoolMax)
	HavePrizePool := make([]*pb.PrizePool, 0)
	if shouldRefresh {
		//初始化所有的数据
		p.LotteryInfo.PrizePool = make([]int32, 0)
		p.LotteryInfo.HavePrizePool = make([]int32, 0)
		//查询是不是第一次玩这个活动，如果是第一次，且配置了
		getUpRecord := p.getUpRecord("lottery")
		is_first_time := true
		if getUpRecord > 0 {
			is_first_time = false
		}
		if is_first_time {
			pool := make([]*global.IdNumRate, 0)
			if global.GlobalCfg.Gacha.Guidneed.ID != 0 {
				// 刷新奖池
				pool, err = m.LotteryRefreshDetail(LotteryPoolMax - 1)
				if err != nil {
					resp.Code = -1
					resp.Error = err.Error()
					return s.Response(resp)
				}
				lotteryInfo := global.GetLotteryCfgByItem(global.GlobalCfg.Gacha.Guidneed.ID)
				if lotteryInfo == nil {
					resp.Code = -2
					resp.Error = errutil.ErrIllegalParameter.Error()
					return s.Response(resp)
				}
				pool = append(pool, &global.IdNumRate{
					Id:     lotteryInfo.ID,
					ItemId: lotteryInfo.ItemId,
					Num:    lotteryInfo.ItemNum,
					Rate:   lotteryInfo.Probability,
				})
			} else {
				// 刷新奖池
				pool, err = m.LotteryRefreshDetail(LotteryPoolMax)
				if err != nil {
					resp.Code = -1
					resp.Error = err.Error()
					return s.Response(resp)
				}
			}
			// Convert []*global.IdNumRate to []*pb.PrizePool
			for i, v := range pool {
				prizePool[i] = &pb.PrizePool{
					Id:     v.Id,
					ItemId: v.ItemId,
					Num:    v.Num,
					Rate:   v.Rate,
				}
				p.LotteryInfo.PrizePool = append(p.LotteryInfo.PrizePool, v.Id)
			}

		} else {
			// 刷新奖池
			pool, err := m.LotteryRefreshDetail(LotteryPoolMax)
			if err != nil {
				resp.Code = -1
				resp.Error = err.Error()
				return s.Response(resp)
			}
			// Convert []*global.IdNumRate to []*pb.PrizePool
			for i, v := range pool {
				prizePool[i] = &pb.PrizePool{
					Id:     v.Id,
					ItemId: v.ItemId,
					Num:    v.Num,
					Rate:   v.Rate,
				}
				p.LotteryInfo.PrizePool = append(p.LotteryInfo.PrizePool, v.Id)
			}
		}

	} else {
		for i, v := range p.LotteryInfo.PrizePool {
			lottery := global.GetLotteryCfg(v)
			if lottery == nil {
				continue
			}
			prizePool[i] = &pb.PrizePool{
				Id:     lottery.ID,
				ItemId: lottery.ItemId,
				Num:    lottery.ItemNum,
				Rate:   lottery.Probability,
			}
		}
	}
	for _, v := range p.LotteryInfo.HavePrizePool {
		lottery := global.GetLotteryCfg(v)
		if lottery == nil {
			continue
		}
		HavePrizePool = append(HavePrizePool, &pb.PrizePool{
			Id:     v,
			ItemId: lottery.ItemId,
			Num:    lottery.ItemNum,
			Rate:   lottery.Probability,
		})
	}
	resp.LotteryInfos = &pb.LotteryPoolInfo{
		RefreshTime:   p.LotteryInfo.RefreshTime, //刷新时间
		PrizePool:     prizePool,                 //上阵槽位
		HavePrizePool: HavePrizePool,             //已有槽位
		Progress:      p.LotteryInfo.Progress,    //进度
		NeedItemId:    p.LotteryInfo.NeedItemId,  //需要道具id
		FreeNum:       p.LotteryInfo.FreeNum,
	}

	return s.Response(resp)
}

// 主动刷新
func (m *Manager) LotteryRefresh(s *session.Session, req *pb.LotteryRefreshResponse) error {
	resp := &pb.LotteryRefreshResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if !p.isModuleOpen(constant.ModuleLottery) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}
	//减去消耗品
	if !p.checkTakeSource(global.GlobalCfg.Gacha.RefreshGacha.UseType, int64(global.GlobalCfg.Gacha.RefreshGacha.Num)) {
		resp.Code = ResponseCodeInvalid
		resp.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(resp)
	}
	p.takeSource(global.GlobalCfg.Gacha.RefreshGacha.UseType, int64(global.GlobalCfg.Gacha.RefreshGacha.Num), constant.ModuleName[constant.ModuleLottery])
	pool := make([]*global.IdNumRate, 0)
	if p.LotteryInfo.Progress >= global.GlobalCfg.Gacha.GachaLuckly.MaxNum {
		pool, err = m.LotteryRefreshDetail(LotteryPoolMax - 1)
		if err != nil {
			resp.Code = -1
			resp.Error = err.Error()
			return s.Response(resp)
		}
		NeedItemId := global.GlobalCfg.Gacha.IntMustOut.BagId
		if p.LotteryInfo.NeedItemId != 0 {
			NeedItemId = p.LotteryInfo.NeedItemId
		}
		lottery := global.GetLotteryCfgByItem(NeedItemId)
		if lottery == nil {
			resp.Code = -2
			resp.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(resp)
		}
		if lottery.IfBag == 1 {
			if itemId, ok := LotteryPool(NeedItemId, 16); ok {
				pool = append(pool, itemId)
			}
		} else {
			pool = append(pool, &global.IdNumRate{
				Id:     lottery.ID,
				ItemId: lottery.ItemId,
				Num:    lottery.ItemNum,
				Rate:   lottery.Probability,
			})
		}
		//重置
		p.LotteryInfo.Progress = 0
		p.LotteryInfo.NeedItemId = 0
	} else {
		pool, err = m.LotteryRefreshDetail(LotteryPoolMax)
		if err != nil {
			resp.Code = -1
			resp.Error = err.Error()
			return s.Response(resp)
		}
	}

	prizePool := make([]*pb.PrizePool, len(pool))
	p.LotteryInfo.PrizePool = []int32{}
	for i, v := range pool {
		prizePool[i] = &pb.PrizePool{
			Id:     v.Id,
			ItemId: v.ItemId,
			Num:    v.Num,
			Rate:   v.Rate,
		}
		p.LotteryInfo.PrizePool = append(p.LotteryInfo.PrizePool, v.Id)
	}
	p.LotteryInfo.Progress += global.GlobalCfg.Gacha.RefreshGacha.Luck
	if p.LotteryInfo.Progress > global.GlobalCfg.Gacha.GachaLuckly.MaxNum {
		p.LotteryInfo.Progress = global.GlobalCfg.Gacha.GachaLuckly.MaxNum
	}
	p.LotteryInfo.HavePrizePool = []int32{}
	resp.LotteryInfos = &pb.LotteryPoolInfo{
		RefreshTime:   p.LotteryInfo.RefreshTime, //刷新时间
		PrizePool:     prizePool,                 //上阵槽位
		HavePrizePool: []*pb.PrizePool{},         //已有槽位
		Progress:      p.LotteryInfo.Progress,    //进度
		NeedItemId:    p.LotteryInfo.NeedItemId,  //需要道具id
		FreeNum:       p.LotteryInfo.FreeNum,
	}
	p.addUpRecord("lottery")
	return s.Response(resp)
}

// 刷新卡池
func (m *Manager) LotteryRefreshDetail(needNum int) ([]*global.IdNumRate, error) {
	Items := make([]*global.IdNumRate, 0, needNum)
	for i := 0; i < needNum; i++ {
		if itemId, ok := LotteryPool(100000, 16); ok {
			Items = append(Items, itemId)
		}
	}
	return Items, nil
}

// 指定某一个道具
func (m *Manager) LotterySelectByItem(s *session.Session, req *pb.LotterySelectRequest) error {
	resp := &pb.LotterySelectResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	NeedItemId := int32(0)
	if req.ItemId != 0 {
		lotteryInfo := global.GetLotteryCfgByItem(req.ItemId)
		if lotteryInfo == nil {
			resp.Code = -2
			resp.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(resp)
		}
		//减去消耗品
		if !p.checkTakeSource(global.GlobalCfg.Gacha.IntMustOut.Id, int64(global.GlobalCfg.Gacha.IntMustOut.Num)) {
			resp.Code = ResponseCodeInvalid
			resp.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(resp)
		}
		p.takeSource(global.GlobalCfg.Gacha.IntMustOut.Id, int64(global.GlobalCfg.Gacha.IntMustOut.Num), constant.ModuleName[constant.ModuleLottery])
		NeedItemId = req.ItemId
	}
	//查询是不是合规
	p.LotteryInfo.NeedItemId = NeedItemId
	resp.ItemId = req.ItemId
	return s.Response(resp)
}

// drawFinalLotteryItem 在指定奖池中进行一次抽取；若抽中为道具组，则进入其子奖池继续抽，直到抽到具体道具或达到最大层级。
func LotteryPool(poolId int32, maxDepth int) (*global.IdNumRate, bool) {
	for depth := 0; depth < maxDepth; depth++ {
		poolCfg := global.GetLotteryCfgByType(poolId)
		if len(poolCfg) == 0 {
			return nil, false
		}
		totalPriority := 0
		for _, it := range poolCfg {
			totalPriority += int(it.Probability)
		}
		if totalPriority <= 0 {
			return nil, false
		}
		r := global.RandIn(totalPriority) + 1
		acc := 0
		for _, it := range poolCfg {
			acc += int(it.Probability)
			if acc >= r {
				// IfBag==1 表示是道具组；BagId 为子奖池ID
				if it.IfBag == 1 {
					poolId = it.ItemId
					break
				}
				// 命中具体道具，返回具体道具ID
				return &global.IdNumRate{Id: it.ID, ItemId: it.ItemId, Num: it.ItemNum, Rate: it.Probability}, true
			}
		}
	}
	return nil, false
}
func drawFinalLotteryItem(p *Player) (*global.IdNumRate, bool) {
	//p.LotteryInfo.PrizePool卡池要去除掉已经抽到的
	// 过滤出未抽中的奖池项（多重集差集：PrizePool 次数 - HavePrizePool 次数）
	prizeCount := make(map[int32]int)
	for _, id := range p.LotteryInfo.PrizePool {
		prizeCount[id]++
	}
	haveCount := make(map[int32]int)
	for _, id := range p.LotteryInfo.HavePrizePool {
		haveCount[id]++
	}
	DellPool := make([]*pb.PrizePool, 0)
	for id, cnt := range prizeCount {
		remain := cnt - haveCount[id]
		for i := 0; i < remain; i++ {
			lottery := global.GetLotteryCfg(id)
			if lottery == nil {
				continue
			}
			DellPool = append(DellPool, &pb.PrizePool{
				Id:     lottery.ID,
				ItemId: lottery.ItemId,
				Num:    lottery.ItemNum,
				Rate:   lottery.Probability,
			})
		}
	}

	totalPriority := 0
	for _, it := range DellPool {
		totalPriority += int(it.Rate)
	}
	if totalPriority <= 0 {
		return nil, false
	}
	r := global.RandIn(totalPriority) + 1
	acc := 0
	for _, it := range DellPool {
		acc += int(it.Rate)
		if acc >= r {
			// 命中具体道具，返回具体道具ID
			return &global.IdNumRate{Id: it.Id, ItemId: it.ItemId, Num: it.Num, Rate: it.Rate}, true
		}
	}
	return nil, false
}

// 抽取奖励
func (m *Manager) Lottery(s *session.Session, req *pb.LotteryRequest) error {
	resp := &pb.LotteryResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	if !p.isModuleOpen(constant.ModuleLottery) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}
	//查询所有的抽奖情况
	Chou := global.GlobalCfg.Gacha.Output.Select
	needItemId := int32(0)
	needItemNum := int64(0)
	needItemLuck := int32(0)
	if len(p.LotteryInfo.HavePrizePool) >= 4 {
		resp.Code = ResponseCodeInvalid
		resp.Error = "抽奖次数已用完"
		return s.Response(resp)
	}
	for _, v := range Chou {
		if p.LotteryInfo.FreeNum > 0 {
			needItemId = 0
			p.LotteryInfo.FreeNum--
			break
		} else {
			if int(v.Id) == len(p.LotteryInfo.HavePrizePool)+1 {
				needItemId = v.UseType
				needItemNum = int64(v.Num)
				needItemLuck = v.Luck
				break
			}
		}

	}
	if needItemId != 0 {
		//减去消耗品
		if !p.checkTakeSource(needItemId, needItemNum) {
			resp.Code = ResponseCodeInvalid
			resp.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(resp)
		}
		p.takeSource(needItemId, needItemNum, constant.ModuleName[constant.ModuleLottery])
	}
	//抽奖
	getUpRecord := p.getUpRecord("lottery")
	is_first_time := true
	if getUpRecord > 0 {
		is_first_time = false
	}
	p.addUpRecord("lottery")
	seedId := int32(0)
	seedItem := int32(0)
	seedNum := int64(0)
	if is_first_time && global.GlobalCfg.Gacha.Guidneed.ID != 0 {
		lotteryInfo := global.GetLotteryCfgByItem(global.GlobalCfg.Gacha.Guidneed.ID)
		if lotteryInfo == nil {
			resp.Code = -2
			resp.Error = errutil.ErrIllegalParameter.Error()
			return s.Response(resp)
		}
		seedId = lotteryInfo.ID
		seedItem = lotteryInfo.ItemId
		seedNum = int64(lotteryInfo.ItemNum)
	} else {
		itemInfo, ok := drawFinalLotteryItem(p)
		if !ok {
			resp.Code = ResponseCodeInvalid
			return s.Response(resp)
		}
		seedId = itemInfo.Id
		seedItem = itemInfo.ItemId
		seedNum = int64(itemInfo.Num)
	}
	if seedItem == 0 || seedNum == 0 {
		resp.Code = ResponseCodeInvalid
		resp.Error = "抽奖道具不存在"
		return s.Response(resp)
	}

	if p.checkAddSource(seedItem, seedNum) {
		p.addSource(seedItem, seedNum, constant.ModuleName[constant.ModuleLottery])
	} else {
		attachment := []protocol.Attachment{}
		attachment = append(attachment, protocol.Attachment{UserPropsId: int64(seedItem), Num: int64(seedNum)})
		p.SendMail(constant.ModuleName[constant.ModuleLottery], attachment)
	}
	partnerInfo := global.GetPartnerCfgByItemId(seedItem)
	//红色以上会发走马灯
	if partnerInfo != nil && partnerInfo.Quality >= 7 {
		strSend := []string{p.name, strconv.Itoa(int(partnerInfo.ItemID))}
		m.Lamp(int32(constant.LampHeroSummon), p.serverID, strSend)
	}
	//如果都抽完了，
	if len(p.LotteryInfo.HavePrizePool) == 3 {
		p.LotteryInfo.HavePrizePool = []int32{}
		//需要自动刷新
		p.LotteryInfo.PrizePool = make([]int32, 0)
		// 刷新奖池
		pool, err := m.LotteryRefreshDetail(4)
		if err != nil {
			resp.Code = -1
			resp.Error = err.Error()
			return s.Response(resp)
		}
		// Convert []*global.IdNumRate to []*pb.PrizePool
		for _, v := range pool {
			p.LotteryInfo.PrizePool = append(p.LotteryInfo.PrizePool, v.Id)
		}
	} else {
		p.LotteryInfo.HavePrizePool = append(p.LotteryInfo.HavePrizePool, seedId)
	}
	resp.Id = seedId
	resp.ItemId = seedItem
	resp.ItemNum = int32(seedNum)
	p.LotteryInfo.Progress += needItemLuck
	if p.LotteryInfo.Progress > global.GlobalCfg.Gacha.GachaLuckly.MaxNum {
		p.LotteryInfo.Progress = global.GlobalCfg.Gacha.GachaLuckly.MaxNum
	}
	p.tasks.dispatch(pb.TargetType_PARTNER_RECEIVE, 0, 1)
	//发放奖励
	p.addUpRecord(constant.LotteryNum)
	return s.Response(resp)
}

// UseAvatar 替换用头像或头像框
func (m *Manager) UseAvatar(s *session.Session, req *pb.UseAvatarRequest) error {
	resp := &pb.UseAvatarResponse{
		Code: 0,
		Id:   req.GetId(),
		Type: req.GetType(),
	}

	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	if !p.isModuleOpen(constant.ModuleAvatar) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}

	avatarCfg, exist := global.AvatarCfg[req.GetId()]
	if !exist {
		resp.Code = -1
		resp.Error = errutil.ErrItemNotFound.Error()
		return s.Response(resp)
	}

	if pb.AvatarType(avatarCfg.ItemType) != req.GetType() {
		resp.Code = -1
		resp.Error = errutil.ErrItemNotFound.Error()
		return s.Response(resp)
	}

	switch req.GetType() {
	case pb.AvatarType_AVATAR:
		if err := p.avatar.SetCurrentAvatar(req.GetId()); err != nil {
			resp.Code = -1
			resp.Error = err.Error()
			logger.Warnf("UseAvatar err %v, pid=%d", err, p.uid)
			return s.Response(resp)
		}
	case pb.AvatarType_FRAME:
		if err := p.avatar.SetCurrentFrame(req.GetId()); err != nil {
			resp.Code = -1
			resp.Error = err.Error()
			logger.Warnf("UseAvatar err %v, pid=%d", err, p.uid)
			return s.Response(resp)
		}
	default:
		resp.Code = -1
		resp.Error = errutil.ErrIllegalParameter.Error()
		logger.Warnf("UseAvatar err %v, pid=%d", errutil.ErrIllegalParameter, p.Uid())
		return s.Response(resp)
	}

	if err := s.Response(resp); err != nil {
		logger.Errorf("UseAvatar s.Response error: %v", err)
	}

	// 设置排行榜信息
	p.userShowRefresh()

	// 不是穿戴上加属性，是获取就加。刷新战力
	// p.refreshStats()

	return nil
}

// GuardIslandMultiplier 挂机副本首次进入选择倍率
func (m *Manager) GuardIslandMultiplier(s *session.Session, req *pb.GuardIslandMultiplierRequest) error {
	res := &pb.GuardIslandMultiplierResponse{
		Code:  0,
		Error: "",
		ID:    req.GetID(),
	}

	logger.Debugf("GuardIslandMultiplierService request.ID: %+v, session: %d", req.GetID(), s.ID())

	// 获取玩家对象
	player, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		logger.Errorf("GuardIslandMultiplier playerWithSession error: %v", err)
		return s.Response(res)
	}

	// 执行业务逻辑
	err = player.GuardIslandMultiplier(req.ID)
	if err != nil {
		var bizErr *errutil.BizError
		if errors.As(err, &bizErr) {
			res.Code = bizErr.Code
			res.Error = bizErr.Error()
			logger.Warnf("GuardIslandMultiplier business error, pid: %d, err: %+v", player.Uid(), err)
		} else {
			res.Code = -1
			res.Error = err.Error()
			logger.Errorf("GuardIslandMultiplier error, pid: %d, err: %+v", player.Uid(), err)
		}
		return s.Response(res)
	}

	// 返回响应
	return s.Response(res)
}

// 聊天
// func (m *Manager) Chat(s *session.Session, req *pb.ChatRequest) error {
// 	res := &pb.ChatResponse{}
// 	p, err := playerWithSession(s)
// 	if err != nil {
// 		return err
// 	}

// 	if strings.ToLower(viper.GetString("core.env")) == constant.EnvProd {
// 		if err := provider.WeChatV2.CheckMessageSecurity(p.accountName, req.Context); err != nil {
// 			res.Code = -1
// 			res.Error = err.Error()
// 			return s.Response(res)
// 		}
// 	}

// 	// gm指令
// 	if len(req.Context) > 0 && req.Context[0] == '@' && strings.ToLower(viper.GetString("core.env")) != constant.EnvProd {
// 		cmd := strings.Split(req.Context, " ")
// 		var id int
// 		var num int
// 		var err error
// 		switch cmd[0] {
// 		case "@add":
// 			if len(cmd) == 3 {
// 				if id, err = strconv.Atoi(cmd[1]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -1
// 					break
// 				}
// 				if num, err = strconv.Atoi(cmd[2]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -2
// 					break
// 				}
// 				p.addSource(int32(id), int64(num), constant.ItemAddGM)
// 			}
// 		case "@clear":
// 			newbag := &pb.Bag{ItemList: []*pb.ItemInfo{}}
// 			packet := pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
// 			for _, v := range p.bag.ItemList {
// 				if p.hasEquipment(v.Uid) {
// 					newbag.ItemList = append(newbag.ItemList, v)
// 					continue
// 				}
// 				packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: v.CfgId, Uid: v.Uid, Num: v.Num})
// 			}
// 			p.pushMsg(protocol.RouteTakeSourcePush, &packet)
// 			p.bag = newbag
// 		case "@god":
// 			godLevel := int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)
// 			if len(cmd) == 2 {
// 				if godLevel, err = strconv.Atoi(cmd[1]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -3
// 					break
// 				}
// 			}

// 			if godLevel > int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel) {
// 				godLevel = int(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)
// 			}
// 			if godLevel < int(p.level) {
// 				logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 				res.Code = -4
// 				break
// 			}

// 			roleCfg, err := global.GetRoleCfg(int32(p.job), int32(godLevel))
// 			if err != nil {
// 				logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 				res.Code = -5
// 				break
// 			}

// 			// p.coin = 200000
// 			// p.soul = 0
// 			// p.coupon = 0
// 			p.tokens[constant.ItemTypeCoin] = 200000
// 			p.tokens[constant.ItemTypeSoul] = 10
// 			p.tokens[constant.ItemTypeCoupon] = 10
// 			p.tokens[constant.ItemTypeTrialVoucher] = 120

// 			p.level = int32(godLevel - 1)
// 			p.exp = roleCfg.NeedExp
// 			if p.level >= global.SwitchCfg[ModuleDestinyFunction] {
// 				p.fateLevel = int32(godLevel)
// 			}

// 			// p.vip = 3
// 			p.energy = 120
// 			p.levelUp()

// 			for _, v := range p.partnerSlots {
// 				if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
// 					partner.level = int32(godLevel)
// 					partner.SyncUserPartner()
// 				}
// 			}

// 			items := []global.IdNum{}
// 			for _, v := range global.ItemCfg {
// 				if v.Itemlevel == int32(math.Min(15, 15*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel))) {
// 					items = append(items, global.IdNum{ID: v.ID, Num: 15})
// 				}
// 				if v.MainType != 5 {
// 					continue
// 				}
// 				if v.Occupylimit != 0 && v.Occupylimit != int(p.job) {
// 					continue
// 				}
// 				if v.Levellimit != int(p.level/10*10) {
// 					continue
// 				}

// 				items = append(items, global.IdNum{ID: v.ID, Num: 1})
// 			}
// 			p.addSourceBatch(items, constant.ItemAddGM)

// 			// if p.level >= global.SwitchCfg[ModuleEquipmentSystemSoulFB] {
// 			// 	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
// 			// 		p.genstoneChangeAll(i, 0)
// 			// 	}
// 			// }

// 			for i := 1; i < len(p.equips); i++ {
// 				if p.level >= global.SwitchCfg[ModuleZhuangBeiShengXing] {
// 					p.equips[i].Stars = int32(math.Min(30, 30*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel-int32(global.SwitchCfg[ModuleZhuangBeiQiangHua]))))
// 				}
// 				if p.level >= global.SwitchCfg[ModuleZhuangBeiQiangHua] {
// 					p.equips[i].Enhance = p.level / 2 // int32(math.Min(300, 300*float64(int32(godLevel)-global.SwitchCfg[ModuleZhuangBeiQiangHua])/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel-int32(global.SwitchCfg[ModuleZhuangBeiQiangHua]))))
// 				}

// 				for _, parslot := range p.partnerSlots {
// 					parslot.Equips[i].Stars = p.equips[i].Stars
// 					parslot.Equips[i].Enhance = p.equips[i].Enhance
// 				}

// 				for _, v := range p.bag.ItemList {
// 					if equipCfg, err := global.GetItemCfg(v.CfgId); err == nil {
// 						if equipCfg.Equiptype == i {
// 							p.equips[i].ItemId = v.Uid
// 						}
// 					}
// 				}
// 			}

// 			p.skillsSlot = [pb.SkillsSlot_MaxSkillsSlot]int32{}
// 			for i := 0; i < len(p.skills); i++ {
// 				if scfg, err := global.GetSkillCfg(p.skills[i].CfgId); err == nil {
// 					if scfg.Level == 0 || scfg.ActiveorPassive == PassiveSkill {
// 						continue
// 					}
// 					if p.level >= global.SwitchCfg[ModuleJiNengShengJi] {
// 						p.skills[i].CfgId = scfg.SkillID*100 + int32(math.Min(30, 30*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)))
// 					}
// 					if p.level >= global.SwitchCfg[ModuleJiNengJinJie] {
// 						p.skills[i].AdvancedLevel = int32(math.Min(6, 6*float64(godLevel)/float64(global.GlobalCfg.BaseValue.Rolelevel.Maxlevel)))
// 					}
// 					if i < int(pb.SkillsSlot_MaxSkillsSlot) && i != 0 {
// 						p.skillsSlot[i] = p.skills[i].CfgId
// 					}

// 				}
// 			}
// 		case "@addpartner":
// 			if len(cmd) == 2 {
// 				if id, err = strconv.Atoi(cmd[1]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -1
// 					break
// 				}

// 				if err := p.addPartner(int32(id)); err == nil {
// 					logger.Debugf("@添加伙伴成功 cmd=%s partnerId=%d ", req.Context, id)
// 				}
// 			} else {
// 				for _, v := range global.PartnerCfg {
// 					if err := p.addPartner(v.Cid); err == nil {
// 						logger.Debugf("@添加伙伴成功 cmd=%s  partnerId=%d ", req.Context, v.Cid)
// 					}
// 				}
// 				for _, v := range p.partners {
// 					v.level = 80
// 					v.star = 20

// 					v.learnSkill()

// 					v.SyncUserPartner()
// 				}
// 			}
// 		case "@addpet":
// 			if len(cmd) == 3 {
// 				if id, err = strconv.Atoi(cmd[1]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -1
// 					break
// 				}
// 				if num, err = strconv.Atoi(cmd[2]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -2
// 					break
// 				}
// 				if _, err := p.addPet(int32(id), 0, 0, int32(num), 0, "gm"); err == nil {
// 					logger.Debugf("@添加宠物成功 cmd=%s petId=%d ", req.Context, id)
// 				}
// 			}
// 		case "@setbat":
// 			if len(cmd) == 2 {
// 				if num, err = strconv.Atoi(cmd[1]); err != nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -1
// 				}
// 				// if num, err = strconv.Atoi(cmd[2]); err != nil {
// 				// 	logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 				// 	res.Code = -2
// 				// 	break
// 				// }
// 				id := num / 100000
// 				cfg := global.GetBattleMapCfg(int32(num))
// 				if cfg == nil {
// 					logger.Debugf("@参数错误 cmd=%s ", req.Context)
// 					res.Code = -3
// 					break
// 				}
// 				if id == int(constant.BattleMap) {
// 					p.InstanceData.BattleMap.BattleId = int32(num)
// 				} else if id == int(constant.BossMap) {
// 					p.InstanceData.BossMap.BattleId = int32(num)
// 				} else if id == int(constant.TrialTowerMap) {
// 					p.InstanceData.TrialTowerMap.BattleId = int32(num)
// 				} else if id == int(constant.RelicMap) {
// 					for _, v := range global.BattleCfg {
// 						if v.CopyType != int32(id) {
// 							continue
// 						}
// 						p.InstanceData.RelicMap.Evaluation = append(p.InstanceData.RelicMap.Evaluation, global.KeyValue{Key: v.UID, Value: 3})
// 					}
// 					p.InstanceData.RelicMap.BattleId = int32(num)
// 				}
// 			}
// 		}

// 		p.syncUserBase()
// 		p.SyncUserAdvance()
// 		p.syncUserSkillSlot()
// 		p.refreshStats()
// 		p.SyncUserPartnerSlot()

// 		return s.Response(res)
// 	}

// 	// 聊天
// 	if req.Type == pb.ChatType_WORLD {
// 		if p.getItemCountInBag(constant.ItemTypeLoudspeaker) >= 1 {
// 			p.takeSource(constant.ItemTypeLoudspeaker, 1)
// 		} else {
// 			if !p.checkTakeSource(global.GlobalCfg.BaseValue.Chatconsume.ID, global.GlobalCfg.BaseValue.Chatconsume.Num) {
// 				res.Code = -1
// 				res.Error = errutil.ErrItemNotEnough.Error()
// 				return s.Response(res)
// 			}
// 			p.takeSource(global.GlobalCfg.BaseValue.Chatconsume.ID, global.GlobalCfg.BaseValue.Chatconsume.Num, constant.ItemDelChat)
// 		}
// 		packet := &pb.ChatPush{Type: req.Type, Context: sensitive.ReplaceSensitiveWords(req.Context), Uid: p.uid, Name: p.name}
// 		if p.serverID != 0 {
// 			for _, v := range m.players {
// 				if v.serverID == p.serverID {
// 					v.pushMsg(protocol.RouteChatPush, packet)
// 					// 添加到聊天记录
// 					// if v.serverID == p.serverID {
// 					// 	v.pushMsg(protocol.RouteChatPush, packet)
// 				}
// 			}
// 			defaultManager.addChatHistory(p.serverID, int32(pb.ChatType_WORLD), p.uid, req.Context, p.name, 0)

// 		}
// 	}

// 	// 如果是行会聊天
// 	if req.Type == pb.ChatType_GUILD {
// 		if p.guildId == 0 {
// 			res.Code = -1
// 			res.Error = errutil.ErrGuildNotFound.Error()
// 			return s.Response(res)
// 		}

// 		packet := &pb.ChatPush{Type: req.Type, Context: sensitive.ReplaceSensitiveWords(req.Context), Uid: p.uid, Name: p.name}
// 		if guildManager.guilds[p.guildId] != nil {
// 			for _, v := range guildManager.guilds[p.guildId].Members {
// 				// 在线的才给推送
// 				// v.MemberId
// 				// 判断是否在线
// 				if pe, ok := m.players[v.MemberId]; ok {
// 					if pe.curState <= PlayerStatusDisconnect {
// 						pe.pushMsg(protocol.RouteChatPush, packet)
// 					}
// 				}
// 			}
// 			defaultManager.addChatHistory(p.serverID, int32(pb.ChatType_GUILD), p.uid, req.Context, p.name, p.guildId)
// 		}

// 	}

// 	return s.Response(res)
// }

// Chat 聊天
func (m *Manager) Chat(s *session.Session, req *pb.ChatRequest) error {
	res := &pb.ChatResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	resp, err := m.chat(p, req)
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
		return s.Response(res)
	}

	return s.Response(resp)
}

// 更新角色姓名
func (m *Manager) ChangeRoleName(s *session.Session, req *pb.ChangeRoleNameRequest) error {
	res := &pb.ChangeRoleNameResponse{}

	p, err := playerWithSession(s)
	if err != nil {
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	name, err := p.changeRoleName(req.RoleName)
	if err != nil {
		res.Code = -2
		res.Error = err.Error()
		res.RoleName = name
		return s.Response(res)
	}
	res.RoleName = name
	return s.Response(res)
}

// 使用道具
func (m *Manager) ItemUse(s *session.Session, req *pb.ItemUseRequest) error {
	player, err := playerWithSession(s)
	if err != nil {
		res := &pb.ItemUseResponse{}
		res.Code = -1
		res.Error = err.Error()
		return s.Response(res)
	}

	resp, err := m.itemUse(context.Background(), player, req)
	if err != nil {
		resp = &pb.ItemUseResponse{}
		var bizErr *errutil.BizError
		if errors.As(err, &bizErr) {
			resp.Code = bizErr.Code
			resp.Error = bizErr.Error()
			logger.Warnf("ItemUse business error, pid: %d, err: %+v", player.Uid(), bizErr.Reason)
		} else {
			resp.Code = -1
			resp.Error = errutil.ErrServerInternal.Error()
			logger.Errorf("ItemUse error, pid: %d, err: %+v", player.Uid(), err)
		}
		return s.Response(resp)
	}

	return s.Response(resp)
}

// RefineInfo 获取洗练信息
func (m *Manager) RefineInfo(s *session.Session, req *pb.RefineInfoRequest) error {
	resp := &pb.RefineInfoResponse{}
	player, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	characterRefineDataList, resultPurchasedSlots := player.refine.ToPb()
	resp.CharacterRefines = characterRefineDataList
	resp.PurchasedRefineSlots = resultPurchasedSlots

	return s.Response(resp)
}

// Refine 洗练
func (m *Manager) Refine(s *session.Session, req *pb.RefineRequest) error {
	resp := &pb.RefineResponse{
		Type:          req.Type,
		CharacterSlot: req.CharacterSlot,
		EquipSlot:     req.EquipSlot,
		RefinedSlot:   req.RefinedSlot,
	}

	player, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	if !player.isModuleOpen(constant.ModelRefined) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}

	stats, err := player.refine.Refine(req, player)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	resp.NewAttributes = stats

	return s.Response(resp)
}

// RefineSelect 洗练属性替换
func (m *Manager) RefineSelect(s *session.Session, req *pb.RefineSelectRequest) error {
	var (
		accountID           string               // 角色ID
		refine_type         string               // 洗练类型（随机或定向）
		power_before_refine int64                // 洗练前战力
		power_after_refine  int64                // 洗练后战力
		before_refining     []thinkingdatax.Attr // 洗练前的槽位属性
		after_refining      []thinkingdatax.Attr // 洗练后的槽位属性
	)

	resp := &pb.RefineSelectResponse{
		CharacterSlot: req.CharacterSlot,
		EquipSlot:     req.EquipSlot,
		RefinedSlot:   req.RefinedSlot,
	}

	player, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	if !player.isModuleOpen(constant.ModelRefined) {
		resp.Code = -21
		resp.Error = errutil.ErrModuleNotOpen
		return s.Response(resp)
	}

	beferStats, _ := player.refine.get(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
	if beferStats != nil {
		before_refining = append(before_refining, thinkingdatax.Attr{
			Attr:  int64(beferStats.GetAttr()),
			Value: int64(beferStats.GetNum()),
		})
	} else {
		before_refining = append(before_refining, thinkingdatax.Attr{
			Attr:  0,
			Value: 0,
		})
	}

	resp, err = player.refine.RefineSelect(req)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	afterStats, _ := player.refine.get(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
	if afterStats != nil {
		after_refining = append(after_refining, thinkingdatax.Attr{
			Attr:  int64(afterStats.GetAttr()),
			Value: int64(afterStats.GetNum()),
		})
	} else {
		after_refining = append(after_refining, thinkingdatax.Attr{
			Attr:  0,
			Value: 0,
		})
	}

	s.Response(resp)

	power_before_refine = int64(player.totalCombat)

	// 如果在活动内，推送到RankTypePinkEntry
	player.rankChangeTCP(constant.RankTypePinkEntry, player.calcPinkEntryNum())
	player.refreshStats()

	power_after_refine = int64(player.totalCombat)

	accountID = strconv.FormatInt(player.Uid(), 10)

	slot3_unlocked, slot4_unlocked := player.refine.GetPurchasedSlotCount()

	switch player.refine.lastRefineType {
	case pb.RefineType_REFINE_TYPE_RANDOM:
		refine_type = "随机"
	case pb.RefineType_REFINE_TYPE_DIRECTIONAL:
		refine_type = "定向"
	default:
		refine_type = "unknown"
	}

	if !proto.Equal(beferStats, afterStats) {
		thinkingdatax.TrackerIns.E(accountID, player.accountName, "equip_refining", player.BuildUserReportMap(tracking.M{
			"refine_type":         refine_type,            // 洗练类型（随机或定向）
			"character_slot":      req.GetCharacterSlot(), // 角色槽位 1-5
			"equip_slot":          req.GetEquipSlot(),     // 装备槽位 1-12
			"refine_slot":         req.GetRefinedSlot(),   // 洗练槽位 1-4
			"power_before_refine": power_before_refine,    // 洗练前战力
			"power_after_refine":  power_after_refine,     // 洗练后战力
			"before_refining":     before_refining,        // 洗练前的槽位属性 [{attr, value}]
			"after_refining":      after_refining,         // 洗练后的槽位属性 [{attr, value}]
			"slot3_unlocked":      slot3_unlocked,         // 当前槽位3解锁的数量
			"slot4_unlocked":      slot4_unlocked,         // 当前槽位4解锁的数量
		}))
	}

	return nil
}
func (m *Manager) SetTempPurchasedSlot(uid int64, characterSlot, equipSlot, refinedSlot int32) error {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return err
	}

	if err := p.refine.validateIndex(pb.CharacterSlot(characterSlot), pb.EquipSlot(equipSlot), pb.RefinedSlot(refinedSlot)); err != nil {
		return err
	}

	p.refine.SetTempPurchasedSlot(p.refine.calculateKey(pb.CharacterSlot(characterSlot), pb.EquipSlot(equipSlot), pb.RefinedSlot(refinedSlot)))

	return nil
}

// 模块打开录入
func (m *Manager) ModuleOpenRecord(s *session.Session, req *pb.ModuleOpenRecordRequest) error {
	resp := &pb.ModuleOpenRecordResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	// 校验是否是合规
	_, ok := global.SwitchCfg[req.ModuleId]
	if !ok {
		resp.Code = -2
		resp.Error = errutil.ErrIllegalParameter.Error()
		return s.Response(resp)
	}
	//去重
	if slices.Contains(p.moduleOpen, req.ModuleId) {
		resp.Code = 1
		return s.Response(resp)
	}
	p.moduleOpen = append(p.moduleOpen, req.ModuleId)
	resp.ModuleIds = p.moduleOpen
	return s.Response(resp)
}

func (m *Manager) RefineIsPurchased(uid int64, characterSlot, equipSlot, refinedSlot int32) (bool, error) {
	p, err := m.GetPlayer(uid)
	if err != nil {
		return false, err
	}

	if err := p.refine.validateIndex(pb.CharacterSlot(characterSlot), pb.EquipSlot(equipSlot), pb.RefinedSlot(refinedSlot)); err != nil {
		return false, err
	}

	return p.refine.IsPurchased(pb.CharacterSlot(characterSlot), pb.EquipSlot(equipSlot), pb.RefinedSlot(refinedSlot)), nil
}
func (m *Manager) MainBattlePrepare(s *session.Session, req *pb.MainBattlePrepareRequest) error {
	resp := &pb.MainBattlePrepareResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	preview, err := p.mainBattlePrepare(req.Job, req.MainBattleId)
	for _, v := range preview {
		resp.Items = append(resp.Items, &pb.MainBattlePrepareResponse_IDQuality{
			ID:      v.ID,
			Num:     v.Num,
			Quality: v.Quality,
		})
	}
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	resp.Code = 1
	return s.Response(resp)
}

func (m *Manager) ModuleNew(s *session.Session, req *pb.ModuleNewRequest) error {
	resp := &pb.ModuleNewResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.ErrorCode = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND
		return s.Response(resp)
	}
	info := req.GetInfo()
	if info != nil {
		_, ok := global.SwitchCfg[info.GetCfgId()]
		if !ok {
			resp.ErrorCode = pb.ErrorCode_ERROR_CODE_NOT_CONFIG
			return s.Response(resp)
		}
		p.ModuleNew.Store(info.GetCfgId(), info)
	}
	resp.ErrorCode = pb.ErrorCode_ERROR_CODE_SUCCESS
	resp.Info = info
	return s.Response(resp)
}

// 神兵强化
func (m *Manager) ArtifactStrength(s *session.Session, req *pb.ArtifactStrengthRequest) error {
	resp := &pb.ArtifactStrengthResponse{}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND
		resp.Error = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND.String()
		return s.Response(resp)
	}

	if exp, err := p.artifactStrength(req.CfgId, req.SoulCfgId); err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	} else {
		resp.AddExp = exp
	}
	return s.Response(resp)
}

// 神兵进阶
func (m *Manager) ArtifactAdvance(s *session.Session, req *pb.ArtifactAdvanceRequest) error {
	resp := &pb.CommonResponse{RequestCmdId: uint32(pb.CmdId_C2S_ArtifactAdvance)}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = int32(pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND)
		resp.Error = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND.String()
		return s.Response(resp)
	}

	if err := p.artifactAdvance(req.CfgId); err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	return s.Response(resp)
}

// 神兵装/卸
func (m *Manager) ArtifactEquip(s *session.Session, req *pb.ArtifactEquipRequest) error {
	resp := &pb.CommonResponse{RequestCmdId: uint32(pb.CmdId_C2S_ArtifactEquip)}
	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = int32(pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND)
		resp.Error = pb.ErrorCode_ERROR_CODE_PLAYER_NOT_FOUND.String()
		return s.Response(resp)
	}

	if err := p.artifactEquip(req.CfgId, req.IsEquip); err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}
	return s.Response(resp)
}

// PushingWall 这怪有点儿墙
func (m *Manager) PushingWall(s *session.Session, req *pb.PushingWallRequest) error {
	resp := &pb.PushingWallResponse{
		WallId: req.GetWallId(),
	}

	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		return s.Response(resp)
	}

	if !p.isModuleOpen(constant.ModulePushingWall) {
		resp.Code = pb.ErrorCode_ERROR_CODE_FUNCTION_NOT_OPEN
		resp.Error = errutil.ErrModuleNotOpen
		logger.Warnf("PushingWall module not open, pid=%d", p.uid)
		return s.Response(resp)
	}

	if req.GetUUID() == "" {
		resp.Code = pb.ErrorCode_ERROR_CODE_ILLEGAL_OPERATION
		resp.Error = "Illegal Operation"
		logger.Warnf("PushingWall UUID is empty, pid=%d", p.uid)
		return s.Response(resp)
	}

	pushingWallMonsterCfg, ok := global.PushingWallMonsterCfg[req.WallId]
	if !ok {
		resp.Code = pb.ErrorCode_ERROR_CODE_CONFIG_NOT_FOUND
		resp.Error = errutil.ErrItemNotFound.Error()
		logger.Warnf("PushingWall config not found, pid=%d, wallId=%d", p.uid, req.WallId)
		return s.Response(resp)
	}

	resp.PushingWallCount = p.PushingWall.UsedCount
	resp.TicketPurchaseCount = p.PushingWall.TicketPurchaseCount

	items := make([]global.IdNum, 0, len(pushingWallMonsterCfg.Item))
	for _, v := range pushingWallMonsterCfg.Item {
		items = append(items, global.IdNum{ID: v.Key, Num: v.Value})
	}

	if p.InstanceData.PushingWall.SessionUUID != req.GetUUID() {
		resp.Code = pb.ErrorCode_ERROR_CODE_ILLEGAL_OPERATION
		resp.Error = "Illegal Operation"
		logger.Warnf("PushingWall UUID mismatch, pid=%d, serverUUID=%s, clientUUID=%s", p.uid, p.InstanceData.PushingWall.SessionUUID, req.GetUUID())
		return s.Response(resp)
	}

	p.distributeRewards(p, items, constant.ModulePushingWall.String(), false)

	p.InstanceData.PushingWall.SessionUUID = ""
	p.InstanceData.PushingWall.TmpUnixTime = 0

	if err := s.Response(resp); err != nil {
		logger.Errorf("PushingWall s.Response error: %v, pid=%d", err, p.uid)
	}

	packet := &pb.GameResultBattlePush{
		OverType: pb.GameOverType_GameOverTypeMDead,
		BattleId: -1,
	}
	for _, v := range items {
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	}

	if err = p.pushMsg(protocol.GameResultBattleMapPush, packet); err != nil {
		logger.Errorf("PushingWall p.pushMsg error: %v, pid=%d", err, p.uid)
	}

	return nil
}

func (m *Manager) Shutdown() {
	log.Info("Manager shutdown completed")
}

// PushingWallInfo 获取这怪有点儿墙信息
func (m *Manager) PushingWallInfo(s *session.Session, req *pb.PushingWallInfoRequest) error {

	resp := &pb.PushingWallInfoResponse{
		WallId: req.GetWallId(),
		OPType: req.GetOPType(),
	}

	p, err := playerWithSession(s)
	if err != nil {
		resp.Code = -1
		resp.Error = err.Error()
		logger.Errorf("PushingWallInfo playerWithSession error: %v, pid=%d", err, p.uid)
		return s.Response(resp)
	}

	if !p.isModuleOpen(constant.ModulePushingWall) {
		resp.Code = pb.ErrorCode_ERROR_CODE_FUNCTION_NOT_OPEN
		resp.Error = errutil.ErrModuleNotOpen
		logger.Warnf("PushingWallInfo module not open, pid=%d", p.uid)
		return s.Response(resp)
	}

	if p.PushingWall.UsedCount >= (global.GlobalCfg.PushingWall.DailyLimit.Max + p.PushingWall.TicketPurchaseCount) {
		resp.Code = pb.ErrorCode_ERROR_CODE_CHALLENGE_NO_TIMES
		resp.Error = errutil.ErrPushingWallNotEnough.Error()
		logger.Warnf("PushingWallInfo not enough times, pid=%d, usedCount=%d, purchaseCount=%d", p.uid, p.PushingWall.UsedCount, p.PushingWall.TicketPurchaseCount)
		return s.Response(resp)
	}

	now := time.Now()
	// 1.进入游戏 2.刷新游戏
	if req.GetOPType() == 1 {
		if p.PushingWall.TmpUnixTime <= now.Unix() {
			sessionUUID := uuid.New()
			resp.UUID = sessionUUID
			p.PushingWall.SessionUUID = sessionUUID
			p.PushingWall.UsedCount++
			// NOTE: 设置3秒的缓冲时间，防止客户端瞬间多次请求进入, 后续会接入新得战斗系统
			p.PushingWall.TmpUnixTime = now.Add(3 * time.Second).Unix()
			//上报完成任务
			p.tasks.dispatch(pb.TargetType_PUSHING_WALL_PLAY_TIMES, 0, 1)
		} else {
			resp.UUID = p.PushingWall.SessionUUID
		}
		p.dirty(true)
	}

	resp.PushingWallCount = p.PushingWall.UsedCount
	resp.TicketPurchaseCount = p.PushingWall.TicketPurchaseCount

	logger.Debugf("PushingWallInfo enter game, pid=%d, req=%+v, resp=%+v", p.uid, req, resp)
	return s.Response(resp)
}
