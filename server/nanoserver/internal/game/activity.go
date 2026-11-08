package game

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"nanoserver/db"
// 	"nanoserver/db/model"
// 	"nanoserver/dispatcher"
// 	"nanoserver/global"
// 	"nanoserver/pkg/constant"
// 	"nanoserver/pkg/errutil"
// 	"nanoserver/protocol"
// 	"nanoserver/protocol/protobuf/pb"

// 	"google.golang.org/protobuf/proto"
// )

// type ActivityService interface {
// 	// 领取奖励
// 	ReceiveReward(pid int64, player *Player, questIDs []int64) error
// 	// ToPB 将活动管理器转换为协议缓冲格式
// 	ToPB() proto.Message
// }

// // ActivityMgr 活动管理器
// type ActivityMgr struct {
// 	CumulativeRecharge   *CumulativeRecharge                                                                     // 累计充值活动管理器
// 	OpenServerCarnival   *OpenServerCarnival                                                                     // 开服狂欢
// 	ActivityLoginHandler map[int32]func(ctx context.Context, activityType int32, ar *model.ActivityRecord) error // 活动登录处理器
// }

// // NewActivityManager 创建一个新的活动管理器实例
// func NewActivityManager() *ActivityMgr {
// 	ar := &ActivityMgr{
// 		CumulativeRecharge: &CumulativeRecharge{
// 			Quests: make(map[int64]*pb.Quest),
// 		},
// 		OpenServerCarnival: &OpenServerCarnival{
// 			Quests: make(map[int64]*pb.Quest),
// 		},
// 	}

// 	ar.ActivityLoginHandler = map[int32]func(ctx context.Context, activityType int32, ar *model.ActivityRecord) error{
// 		constant.ActivityTypeSevenDaysFun: ar.SevenDaysFunLoginHandler,
// 	}

// 	return ar
// }

// // Save 保存活动状态到数据库
// func (am *ActivityMgr) Save() (err error) {
// 	cr, err1 := proto.Marshal(am.CumulativeRecharge.ToPB())
// 	osc, err2 := proto.Marshal(am.OpenServerCarnival.ToPB())
// 	if err = errors.Join(err1, err2); err != nil {
// 		return err
// 	}

// 	ar := &model.ActivityRecord{
// 		CumulativeRecharge: cr,
// 		OpenServerCarnival: osc,
// 	}
// 	_, err = db.ActivityRecordUpsert(ar)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Load 从数据库加载活动状态
// func (am *ActivityMgr) Load(PID int64) error {
// 	ar, err := db.ActivityRecordGet(PID)
// 	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
// 		return err
// 	} else if errors.Is(err, errutil.ErrNotFound) {
// 		return nil // 如果没有找到记录，返回 nil
// 	}

// 	// 根据活动 id 获取时间段
// 	for _, actCtrCfg := range global.ActivityControllerCfg {
// 		// 判断是否开启
// 		if !am.IsOpen(actCtrCfg) {
// 			logger.Warn("活动未开启", "PID", PID, "ActivityID", actCtrCfg.UUID)
// 			continue
// 		}

// 		// 加载并推送
// 		handler, ok := am.ActivityLoginHandler[int32(actCtrCfg.UUID)]
// 		if !ok {
// 			return fmt.Errorf("未找到活动类型 %d 的登录处理器", actCtrCfg.UUID)
// 		}
// 		if err := handler(context.Background(), int32(actCtrCfg.UUID), ar); err != nil {
// 			return fmt.Errorf("活动类型 %d 的登录处理器执行失败: %w", actCtrCfg.UUID, err)
// 		}
// 	}

// 	// cr := &pb.ActivityCumulativeRechargePush{}
// 	// ao := &pb.ActivityOpenServerCarnivalPush{}

// 	// if ar != nil && ar.CumulativeRecharge != nil {
// 	// 	if err = proto.Unmarshal(ar.CumulativeRecharge, cr); err != nil {
// 	// 		return fmt.Errorf("failed to unmarshal CumulativeRecharge: %w", err)
// 	// 	}
// 	// }

// 	// if ar != nil && ar.OpenServerCarnival != nil {
// 	// 	if err = proto.Unmarshal(ar.OpenServerCarnival, ao); err != nil {
// 	// 		return fmt.Errorf("failed to unmarshal OpenServerCarnival: %w", err)
// 	// 	}
// 	// }

// 	// am.CumulativeRecharge.TotalAmount = float64(cr.TotalRecharge)
// 	// am.CumulativeRecharge.Quests = make(map[int64]*pb.Quest, len(cr.Quests))
// 	// for _, quest := range cr.Quests {
// 	// 	am.CumulativeRecharge.Quests[quest.Id] = &pb.Quest{
// 	// 		Id:      quest.Id,
// 	// 		Receive: quest.Receive,
// 	// 	}
// 	// }

// 	// am.OpenServerCarnival.Quests = make(map[int64]*pb.Quest, len(ao.Quests))
// 	// for _, quest := range ao.Quests {
// 	// 	am.OpenServerCarnival.Quests[quest.Id] = &pb.Quest{
// 	// 		Id:      quest.Id,
// 	// 		Receive: quest.Receive,
// 	// 	}
// 	// }

// 	logger.Info("ActivityMgr loaded successfully", "PID", PID,
// 		"CumulativeRecharge", am.CumulativeRecharge,
// 		"OpenServerCarnival", am.OpenServerCarnival.Quests)

// 	return err
// }

// func (am *ActivityMgr) IsOpen(acc *global.ActivityControllerConfig) bool {
// 	if acc == nil {
// 		return false
// 	}

// 	// if

// 	return false
// }

// func (am *ActivityMgr) RewardReceive(PID int64, player *Player, req *pb.RewardReceive) error {
// 	var (
// 		instance     ActivityService
// 		activityType int32
// 	)

// 	// 工厂获取对象
// 	if req.GetDesign() == int32(constant.DocxActivity) && len(req.GetQuestId()) > 0 {
// 		activityInfo, ok := global.ActivityCfg[int32(req.GetQuestId()[0])]
// 		if !ok {
// 			return errutil.ErrNotFound
// 		}
// 		switch activityInfo.Type {
// 		case constant.ActivityTypeCumulativeCharge:
// 			instance = am.CumulativeRecharge
// 		}
// 	} else {
// 		return errutil.ErrIllegalParameter
// 	}

// 	// 领取奖励
// 	if err := instance.ReceiveReward(PID, player, req.GetQuestId()); err != nil {
// 		return err
// 	}

// 	// 推送活动数据
// 	dispatcher.GetInstance().Publish(dispatcher.TopicActivityPush.String(), &dispatcher.ActivityPushEvent{
// 		Type: dispatcher.ActivityType(activityType),
// 		PID:  PID,
// 		Data: instance.ToPB(),
// 	})

// 	return nil
// }

// // =============== Cumulative Recharge Activity Management ===============
// // CumulativeRecharge 累计充值管理器

// var _ ActivityService = (*CumulativeRecharge)(nil)

// type CumulativeRecharge struct {
// 	TotalAmount float64             // 累计充值总金额
// 	Quests      map[int64]*pb.Quest // 累计充值任务列表
// }

// // Update 更新累计充值任务状态
// func (cr *CumulativeRecharge) Update(PID int64, amount string) error {
// 	cumulativeRechargeQuests, ok := global.ActivityCfgByType[constant.ActivityTypeCumulativeCharge]
// 	if !ok {
// 		logger.Error("No cumulative recharge quests found")
// 		return nil
// 	}

// 	for _, v := range cumulativeRechargeQuests {
// 		questID := v.UUID
// 		if _, exists := cr.Quests[questID]; !exists {
// 			if cr.TotalAmount >= float64(v.PreLimit) {
// 				cr.Quests[questID] = &pb.Quest{
// 					Id:      questID,
// 					Receive: false,
// 				}
// 			}
// 		}
// 	}

// 	// 触发则推送活动
// 	dispatcher.GetInstance().Publish(dispatcher.TopicActivityPush.String(), &dispatcher.ActivityPushEvent{
// 		Type: dispatcher.ActivityCumulativeRecharge,
// 		PID:  PID,
// 		Data: cr.ToPB(),
// 	})

// 	return nil
// }

// // ToPB 将活动管理器转换为协议缓冲格式
// func (cr *CumulativeRecharge) ToPB() proto.Message {
// 	quests := make([]*pb.Quest, 0, len(cr.Quests))
// 	for _, quest := range cr.Quests {
// 		quests = append(quests, &pb.Quest{
// 			Id:      quest.Id,
// 			Receive: quest.Receive,
// 		})
// 	}
// 	return &pb.ActivityCumulativeRechargePush{
// 		TotalRecharge: int64(cr.TotalAmount),
// 		Quests:        quests,
// 	}
// }

// // 领取奖励
// func (cr *CumulativeRecharge) ReceiveReward(pid int64, player *Player, questIDs []int64) error {
// 	ads := protocol.AddSource{
// 		Uid:    pid,
// 		Source: make([]*protocol.SourceInfo, 0, len(questIDs)),
// 	}

// 	// 判断是否可以领取
// 	for _, v := range questIDs {
// 		quest, exists := cr.Quests[v]
// 		if !exists {
// 			return fmt.Errorf("任务未触发ID %d: %w", v, errutil.ErrNotFound)
// 		}

// 		if quest.Receive {
// 			return errutil.ErrAlreadyReceived
// 		}

// 		cumulativeRechargeQuests, ok := global.ActivityCfg[int32(v)]
// 		if !ok {
// 			return fmt.Errorf("未找到任务ID %d 的活动配置: %w", v, errutil.ErrNotFound)
// 		}

// 		for _, v := range cumulativeRechargeQuests.CommonItem {
// 			ads.Source = append(ads.Source, &protocol.SourceInfo{
// 				Id:  v.Key,
// 				Num: int64(v.Value),
// 			})
// 		}

// 		if cr.TotalAmount < float64(cumulativeRechargeQuests.PreLimit) {
// 			return fmt.Errorf("累计充值金额不足, 无法领取任务ID %d 的奖励: %w", v, errutil.ErrNotEligibleForReward)
// 		}
// 	}

// 	// 领取奖励
// 	SourceAdd(ads)

// 	// 修改领取状态
// 	for _, v := range questIDs {
// 		quest := cr.Quests[v] // 上面检查过了，这里忽略
// 		quest.Receive = true
// 	}

// 	return errutil.ErrNotEligibleForReward
// }

// // =============== Open Server Carnival Activity Management ===============
// // OpenServerCarnival 开服狂欢管理器
// type OpenServerCarnival struct {
// 	// Add fields and methods for OpenServerCarnival management
// 	Quests map[int64]*pb.Quest // 累计充值任务列表
// }

// // Update 更新开服狂欢任务状态
// func (osc *OpenServerCarnival) Update(PID int64, amount string) error {
// 	return nil
// }

// // ToPB 将开服狂欢管理器转换为协议缓冲格式
// func (osc *OpenServerCarnival) ToPB() *pb.ActivityOpenServerCarnivalPush {
// 	quests := make([]*pb.Quest, 0, len(osc.Quests))
// 	for _, quest := range osc.Quests {
// 		quests = append(quests, &pb.Quest{
// 			Id:      quest.Id,
// 			Receive: quest.Receive,
// 		})
// 	}
// 	return &pb.ActivityOpenServerCarnivalPush{
// 		Quests: quests,
// 	}
// }
