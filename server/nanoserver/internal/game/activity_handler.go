package game

// import (
// 	"nanoserver/dispatcher"
// 	"nanoserver/protocol"
// )

// func ActivityCumulativeRechargeHandler(eventData *dispatcher.RechargeSuccessEvent) {
// 	p, err := defaultManager.GetPlayer(eventData.PID)
// 	if err != nil {
// 		logger.Errorf("获取玩家失败: %v", err)
// 		return
// 	}

// 	if err = p.activityMgr.CumulativeRecharge.Update(eventData.PID, eventData.Amount); err != nil {
// 		logger.Errorf("更新累计充值任务失败: %v", err)
// 		return
// 	}
// }

// func ActivityPushHandler(eventData *dispatcher.ActivityPushEvent) {
// 	p, err := defaultManager.GetPlayer(eventData.PID)
// 	if err != nil {
// 		logger.Errorf("获取玩家失败: %v", err)
// 		return
// 	}

// 	switch eventData.Type {
// 	case dispatcher.ActivityCumulativeRecharge:
// 		if err := p.pushMsg(protocol.RouteActivityCumulativeRechargePush, eventData.Data); err != nil {
// 			logger.Errorf("推送累计充值活动失败: %v", err)
// 			return
// 		}
// 	default:
// 		logger.Warnf("未知活动类型: %v", eventData.Type)
// 		return
// 	}
// }

// var map[uint32]func()
