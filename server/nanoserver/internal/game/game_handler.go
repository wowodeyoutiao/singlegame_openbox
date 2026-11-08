package game

import (
	"context"
	"nanoserver/dispatcher"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/goredisx"
	"nanoserver/pkg/utils"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func PVEWarMapBossDieCountHandler(eventData *dispatcher.PVEWarMapBossDieCountEvent) {
	key := "pve_war_map:boss_die:" + strconv.FormatInt(int64(eventData.ZoneID), 10) + ":" + strconv.FormatInt(int64(eventData.MonsterID), 10)
	if err := goredisx.Instance.Incr(context.Background(), key).Err(); err != nil {
		logger.Errorf("PVEWarMapBossDieCountHandler Incr error: %v", err)
		return
	}
}

// ActivatePurchasedSlotHandler 激活已购买的槽位处理函数
func ActivatePurchasedSlotHandler(eventData *dispatcher.ActivatePurchasedSlot) {
	player, err := defaultManager.GetPlayer(eventData.UID)
	if err != nil {
		logger.Errorf("ActivatePurchasedSlotHandler GetPlayer error: %v", err)
		return
	}

	player.lock.Lock()
	defer player.lock.Unlock()

	_, _, refinedSlot := player.refine.parseKey(player.refine.GetTempPurchasedSlot())
	if refinedSlot != pb.RefinedSlot(eventData.PurchasedSlot) {
		logger.Warnf("ActivatePurchasedSlotHandler refinedSlot not match: %v != %v", refinedSlot, eventData.PurchasedSlot)
		return
	}

	if err := player.refine.ActivatePurchasedSlot(player); err != nil {
		logger.Errorf("ActivatePurchasedSlotHandler ActivatePurchasedSlot error: %v", err)
		return
	}

	player.dirty(true)

	player.pushMsg(protocol.RouteCommonNotificationPush, &pb.CommonNotificationPush{
		ModuleID: int32(constant.ModelRefined),
	})
}

// UserShowRefreshHandler 处理用户展示信息刷新事件
func UserShowRefreshHandler(eventData *dispatcher.UserShowRefreshEvent) {
	// 	20user:2610324357658001408
	// -1103user:2610324357658001408
	// serverUser := strconv.FormatInt(int64(eventData.ServerID), 10) + "user:" + strconv.FormatInt(eventData.UserID, 10)
	// zoneUser := strconv.FormatInt(int64(eventData.ZoneID), 10) + "user:" + strconv.FormatInt(eventData.UserID, 10)

	// if goredisx.Instance.HMSet(context.TODO(), serverUser, map[string]interface{}{
	// 	"name":      eventData.UserShowRefresh.Name,
	// 	"avatar":    eventData.UserShowRefresh.Avatar,
	// 	"frame":     eventData.UserShowRefresh.Frame,
	// 	"sex":       eventData.UserShowRefresh.Sex,
	// 	"job":       eventData.UserShowRefresh.Job,
	// 	"robot_cid": 0,
	// 	"server_id": eventData.UserShowRefresh.ServerID,
	// }).Err() != nil {
	// 	logger.Errorf("UserShowRefreshHandler HMSet error: %v", eventData)
	// }

	// if goredisx.Instance.HMSet(context.TODO(), zoneUser, map[string]interface{}{
	// 	"name":      eventData.UserShowRefresh.Name,
	// 	"avatar":    eventData.UserShowRefresh.Avatar,
	// 	"frame":     eventData.UserShowRefresh.Frame,
	// 	"sex":       eventData.UserShowRefresh.Sex,
	// 	"job":       eventData.UserShowRefresh.Job,
	// 	"robot_cid": 0,
	// 	"server_id": eventData.UserShowRefresh.ServerID,
	// }).Err() != nil {
	// 	logger.Errorf("UserShowRefreshHandler HMSet error: %v", eventData)
	// }
	err := provide.SaveUserInfo(&protocol.RankUserInfoInput{
		Uid:      eventData.UserID,
		Name:     eventData.UserShowRefresh.Name,
		Sex:      int32(eventData.UserShowRefresh.Sex),
		Job:      int32(eventData.UserShowRefresh.Job),
		RobotCId: 0,
		ServerId: int32(eventData.UserShowRefresh.ServerID),
		AvatarV2: int32(eventData.UserShowRefresh.Avatar),
		Frame:    int32(eventData.UserShowRefresh.Frame),
	})
	if err != nil {
		logger.Errorf("UserShowRefreshHandler SaveUserInfo error: %v", err)
	}
}

func PayRefreshPushHandler(eventData *pb.PayRefreshPush) {
	p, err := defaultManager.GetPlayer(int64(eventData.UserId))
	if err != nil {
		logger.Errorf("获取玩家失败: %v", err)
		return
	}

	logger.Infof("PayRefreshPushHandler: %v", eventData)

	p.pushMsg(protocol.RoutePayRefreshPush, &pb.PayRefreshPush{
		RequestId: eventData.RequestId,
		UserId:    eventData.UserId,
		StoreId:   eventData.StoreId,
	})
}

// AdClickHandler 广告点击处理函数
func AdClickHandler(eventData *dispatcher.AdClickEvent) {
	logger.Debugf("AdClickHandler event: %+v", eventData)

	result, err := utils.ParseAdQuery(eventData.ClickID)
	if err != nil {
		logger.Warnf("AdClickHandler ParseAdQuery %v", err)
		return
	}

	// 去数据库获取 Redis 记录
	switch result.Source {
	case utils.AdByteDance:
		callback := ""
		if result.ByteDance.Clickid != "" {
			callback = result.ByteDance.Clickid
		} else {
			callback = result.ByteDance.ClueToken
		}
		if callback != "" {
			// https://event-manager.oceanengine.com/docs/8650/app_api_docs
			conversionEventRequest := &protocol.ConversionEventRequest{
				EventType: "active_register",
				Context: protocol.Context{
					Ad: protocol.Ad{
						Callback: callback,
					},
				},
				Timestamp: time.Now().UnixMilli(),
			}
			conversionEventResponse := &protocol.ConversionEventResponse{}
			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(conversionEventRequest).
				SetResult(conversionEventResponse).
				Post(constant.ByteDanceConversionURL)
			if err != nil {
				logger.Errorf("AdClickHandler: HTTP request to ByteDance failed: %v", err)
			}
			if !resp.IsSuccess() {
				logger.Errorf("AdClickHandler: HTTP request to ByteDance failed with status code %d: %s", resp.StatusCode(), resp.String())
			}
			if conversionEventResponse.Code != 0 {
				logger.Errorf("AdClickHandler: ByteDance conversion event failed with code %d: %s", conversionEventResponse.Code, conversionEventResponse.Message)
			}
		}
		// if err := thinkingdatax.TrackerIns.SetOnce(strconv.FormatInt(eventData.Uid, 10), eventData.OpenID, tracking.M{
		// 	thinkingdatax.AdSource.String(): result.Source.String(),
		// }); err != nil {
		// 	logger.Warnf("ThinkingData UserSetOnce failed: accountID=%s, error=%v", eventData.OpenID, err)
		// }
	case utils.AdGDT:
		// if err := thinkingdatax.TrackerIns.SetOnce(strconv.FormatInt(eventData.Uid, 10), eventData.OpenID, tracking.M{
		// 	thinkingdatax.AdSource.String(): result.Source.String(),
		// }); err != nil {
		// 	logger.Warnf("ThinkingData UserSetOnce failed: accountID=%s, error=%v", eventData.OpenID, err)
		// }
	default:
		logger.Warnf("AdClickHandler")
	}
}
