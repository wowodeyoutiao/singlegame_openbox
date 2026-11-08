package api

import (
	"context"
	"encoding/json"
	"math"
	"nanoserver/db"
	"nanoserver/global"
	"nanoserver/internal/game"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

// calculateBeginTime 计算活动开始时间
// StartTime 规则：
// - 0: 从角色注册时间起算
// - >= 100: 作为时间戳使用
// - < 100: 作为开服第几天使用
func calculateBeginTime(startTime int32, registerAt int64, serverID int32) int64 {
	if startTime == 0 {
		return registerAt
	}
	if startTime >= 100 {
		return int64(startTime)
	}
	// startTime < 100，按开服第几天计算
	if serverID <= 0 {
		// 如果没有服务器ID，回退到注册时间
		return registerAt
	}
	// 获取服务器开服时间
	serverInfo, err := global.ServerListMgrInstance.Get(int(serverID))
	if err != nil {
		// 如果获取失败，回退到注册时间
		return registerAt
	}
	// 计算开服第 startTime 天的时间戳（开服时间 + (startTime - 1) * 24小时）
	// OpenBy 是毫秒时间戳，需要转换为秒
	openTime := time.Unix(serverInfo.OpenBy/1000, 0)
	openDayMidnight := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 0, 0, 0, 0, openTime.Location())
	beginTime := openDayMidnight.AddDate(0, 0, int(startTime)-1)
	return beginTime.Unix()
}

func MakeActivityService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/activity/listall", nex.Handler(activityListAll)).Methods("POST")                          //所有活动情况2.0
	router.Handle("/v1/activity/listcontinuity", nex.Handler(activityListContinuity)).Methods("POST")            //连续登陆2.0
	router.Handle("/v1/activity/receiveday", nex.Handler(activityReceiveDay)).Methods("PUT")                     //连续登陆的签到2.0
	router.Handle("/v1/activity/listcumulative", nex.Handler(activityListCumulative)).Methods("POST")            //累计登录
	router.Handle("/v1/activity/receivecumulative", nex.Handler(activityReceiveCumulative)).Methods("PUT")       //累计登录一键领取
	router.Handle("/v1/activity/listcharge", nex.Handler(activityListCharge)).Methods("POST")                    //首充列表2.0
	router.Handle("/v1/activity/receivecharge", nex.Handler(activityReceiveCharge)).Methods("PUT")               //首充领取2.0
	router.Handle("/v1/activity/listgrade", nex.Handler(activityListGrade)).Methods("POST")                      //等级列表
	router.Handle("/v1/activity/receivegrade", nex.Handler(activityReceiveGrade)).Methods("PUT")                 //等级礼包获取
	router.Handle("/v1/activity/listturntable", nex.Handler(activityListTurntable)).Methods("POST")              //大转盘抽取情况
	router.Handle("/v1/activity/receiveturntable", nex.Handler(activityReceiveTurntable)).Methods("PUT")         //大转盘抽奖
	router.Handle("/v1/activity/listtreasurebox", nex.Handler(activityListTreasurebox)).Methods("POST")          //百宝箱抽取情况
	router.Handle("/v1/activity/receivetreasurebox", nex.Handler(activityReceiveTreasureBox)).Methods("PUT")     //百宝箱抽奖
	router.Handle("/v1/activity/listtrial", nex.Handler(activityListTrial)).Methods("POST")                      //累计登录
	router.Handle("/v1/activity/receivetrial", nex.Handler(activityReceiveTrial)).Methods("PUT")                 //累计登录一键领取
	router.Handle("/v1/activity/groupup", nex.Handler(activityListGroupup)).Methods("POST")                      //成长基金2.0
	router.Handle("/v1/activity/receivegroupup", nex.Handler(activityReceiveGroupup)).Methods("PUT")             //成长基金一键领取2.0
	router.Handle("/v1/activity/listrechargegrade", nex.Handler(activityListRechargeGrade)).Methods("POST")      //累充礼包列表
	router.Handle("/v1/activity/receiverechargegrade", nex.Handler(activityReceiveRechargeGrade)).Methods("PUT") //累充礼包获取
	router.Handle("/v1/activity/listgradebytype", nex.Handler(activityListGradeByType)).Methods("POST")          //通用礼包列表（通过type区分）
	router.Handle("/v1/activity/receivegradebytype", nex.Handler(activityReceiveGradeByType)).Methods("PUT")     //通用礼包领取（通过type区分）
	return router
}

func activityListAll(ctx context.Context) (*protocol.ActivityAllResponse, error) {
	//[{2024-05-09 2} {2024-05-10 1}]
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	list, err := ActivityListAll(request)
	if err != nil {
		return &protocol.ActivityAllResponse{}, err
	}
	return &protocol.ActivityAllResponse{Data: list}, nil
}

func ActivityListAll(r *protocol.ActivityOperateRequest) ([]protocol.ActivityAll, error) {
	//比较等级
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, err
	}
	activityCfg := global.ActivityCfg
	result := make([]protocol.ActivityAll, 0)
	//活动内容
	activity := make(map[int32]map[int32]int)
	for _, activityn := range activityCfg {
		if len(activity[activityn.Acttag]) == 0 {
			activity[activityn.Acttag] = map[int32]int{}
		}
		if activity[activityn.Acttag][activityn.Type] != 1 {
			now := time.Now().Unix()
			beginTime := calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)

			// 对于 StartTime < 100 的活动，表示按服务器开服天数开启，需要检查当前服务器开服天数是否已达到开启条件
			if activityn.StartTime > 0 && activityn.StartTime < 100 {
				// 检查当前服务器开服天数是否已达到开启条件
				if res.ServerId > 0 {
					currentOpenDay := global.ServerListMgrInstance.GetServerOpenDay(res.ServerId)
					if int32(currentOpenDay) < activityn.StartTime {
						// 还没到开启天数，跳过该活动
						continue
					}
				}
			}

			if activityn.Type == constant.ActivityTypeContinuity ||
				activityn.Type == constant.ActivityTypeCumulative ||
				activityn.Type == constant.ActivityTypeTrunTable ||
				activityn.Type == constant.ActivityTypeBounty ||
				activityn.Type == constant.ActivityTypeSingleCharge ||
				activityn.Type == constant.ActivityTypeCumulativeCharge ||
				activityn.Type == constant.ActivityTypePartnerGift ||
				activityn.Type == constant.ActivityTypePetGift ||
				activityn.Type == constant.ActivityTypeDailyGift ||
				activityn.Type == constant.ActivityTypeElementGift ||
				activityn.Type == constant.ActivityTypeSevenDaysFun ||
				activityn.Type == constant.ActivityTypeAssist ||
				activityn.Type == constant.ActivityTypeSevenDaysFunGift ||
				activityn.Type == constant.ActivityTypeStarWishSummon ||
				activityn.Type == constant.ActivityTypeCrazyRefine ||
				activityn.Type == constant.ActivityTypeLottery ||
				activityn.Type == constant.ActivityTypeRefine ||
				activityn.Type == constant.ActivityTypeInstance ||
				activityn.Type == constant.ActivityTypePinkEntry ||
				activityn.Type == constant.ActivityTypeRedPartnerStar ||
				activityn.Type == constant.ActivityTypeEquipmentStar ||
				activityn.Type == constant.ActivityTypePetSynthesis ||
				activityn.Type == constant.ActivityTypeSoulStoneLevel ||
				activityn.Type == constant.ActivityTypeFiveElementLevel {
				// 将时间戳转换为time.Time类型
				originalTime := time.Unix(beginTime, 0)
				// 获取第二天凌晨的时间
				dayss := 0
				// if activityn.Type == constant.ActivityTypeContinuity {
				// 	dayss = 1
				// }
				nextDayMidnight := originalTime.AddDate(0, 0, int(activityn.ResetTime)-dayss) // 添加N天
				nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
				// 转换回时间戳
				nextDayMidnightTimestamp := nextDayMidnight.Unix()
				if now < nextDayMidnightTimestamp {
					item := protocol.ActivityAll{
						Id:        activityn.Acttag,
						ActType:   activityn.Type,
						StartTime: beginTime,
						EndTime:   nextDayMidnightTimestamp,
						ResetTime: activityn.ResetTime,
					}
					result = append(result, item)
				}
				activity[activityn.Acttag][activityn.Type] = 1
			}
			if activityn.Type == constant.ActivityTypeCharge {
				if activityn.Content != 1 {
					continue
				}
				goodsId := make([]int64, 0)
				goodsId = append(goodsId, int64(activityn.PreLimit))
				StoType := int32(goodsId[0] / 1000)
				//查询用户购买情况
				resultb, _ := db.StoreList(
					r.Uid,
					StoType,
					int32(365*50),
					goodsId,
					false,
				)

				if beginTime == int64(0) && len(resultb) == 0 {
					//t, err := time.Parse("2006-01-02", resultd[0].CreatedAt)
					item := protocol.ActivityAll{
						Id:        activityn.Acttag,
						ActType:   activityn.Type,
						StartTime: beginTime,
						EndTime:   -1,
						ResetTime: activityn.ResetTime,
					}
					result = append(result, item)
				} else {
					if len(resultb) > 0 {
						beginTime = resultb[0].CreateTime
						// 将时间戳转换为time.Time类型
						originalTime := time.Unix(beginTime, 0)
						// 获取第二天凌晨的时间
						nextDayMidnight := originalTime.AddDate(0, 0, int(activityn.ResetTime)) // 添加N天
						nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
						// 转换回时间戳
						nextDayMidnightTimestamp := nextDayMidnight.Unix()
						if now < nextDayMidnightTimestamp {
							item := protocol.ActivityAll{
								Id:        activityn.Acttag,
								ActType:   activityn.Type,
								StartTime: beginTime,
								EndTime:   nextDayMidnightTimestamp,
								ResetTime: activityn.ResetTime,
							}
							result = append(result, item)
						}
					}
				}
				//if err != nil {
				//	return nil, errutil.ErrDBOperation
				//}

				activity[activityn.Acttag][activityn.Type] = 1
			}

		}

	}
	return result, nil
}
func activityListContinuity(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityListResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityListResponse{}, errutil.ErrInvalidParameter
	}
	list, days, err := ActivityListDay(request)
	if err != nil {
		return &protocol.ActivityListResponse{}, err
	}
	return &protocol.ActivityListResponse{Data: list, Day: days}, nil
}
func ActivityListDay(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListDay, int64, error) {

	actType := constant.ActivityTypeContinuity
	activityCfg := global.ActivityCfg
	result := make([]protocol.ActivityDetailDay, 0)
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				result = append(result, protocol.ActivityDetailDay{
					UUID:    int64(activityn.UUID),
					Content: activityn.Content,
				})

			}
		}
	}
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, 0, err
	}
	// 计算开始时间
	var beginTime int64 = 0
	var UserBeginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId && activityn.Type == actType {
			beginTime = calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)
			break
		}
	}
	UserBeginTime = beginTime
	endTime := time.Now()
	// 计算时间差
	createdTime := time.Unix(UserBeginTime, 0)
	createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
	duration := endTime.Sub(createdDayMidnight)
	// 计算相差的天数，即目前第几天
	days := int64(duration.Hours()/24) + 1

	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		r.ActivityId,
	)
	roleReceiveInfo := make(map[int32]bool)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = true
	}
	resultEnd := make([]protocol.ActivityListDay, 0)
	for _, activityn := range result {
		activitynDay := activityn.Content

		//判断时间是否超过当前日期
		//状态 1-待领取，2-已领取，3-不可领取，4-需要补签
		if int64(activitynDay) > days {
			resultEnd = append(resultEnd, protocol.ActivityListDay{
				Day:    activityn.Content,
				Status: 3,
			})
		} else {
			//是否领取
			if roleReceiveInfo[activityn.Content] {
				resultEnd = append(resultEnd, protocol.ActivityListDay{
					Day:    activityn.Content,
					Status: 2,
				})
			} else {
				//判断是否登录
				if activitynDay == res.CreateDay {
					resultEnd = append(resultEnd, protocol.ActivityListDay{
						Day:    activityn.Content,
						Status: 1,
					})

				} else {
					resultEnd = append(resultEnd, protocol.ActivityListDay{
						Day:    activityn.Content,
						Status: 4,
					})
				}
			}
		}

	}
	return resultEnd, days + 1, nil
}

func activityReceiveDay(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.UUID = r.UUID
	request.Ad = r.Ad
	if request.UUID == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	Cumulative, err := ActivityCumulative(request)
	if !Cumulative {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// 签到，给到日期
func ActivityCumulative(r *protocol.ActivityOperateRequest) (bool, string) {
	actType := constant.ActivityTypeContinuity
	activityCfg := global.ActivityCfg
	result := protocol.ActivityDetailDay{}
	beginTime := int64(0)
	UserBeginTime := int64(0)
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.UUID == r.UUID {
			result = protocol.ActivityDetailDay{
				UUID:                activityn.UUID,
				Content:             activityn.Content,
				CompensationConsume: make([]protocol.SourceInfo, len(activityn.CompensationConsume)),
				CommonItem:          make([]protocol.SourceInfo, len(activityn.CommonItem)),
			}
			for i, attachment := range activityn.CompensationConsume {
				// 假设 Attachment 结构体中有可以复制的字段
				result.CompensationConsume[i] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment.Key,
					Num: int64(attachment.Value),
				}
			}

			for i, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				result.CommonItem[i] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}
			// 复制 Attachments 字段

		}
	}
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err.Error()
	}
	// 计算开始时间
	for _, activityn := range activityCfg {
		if activityn.UUID == r.UUID {
			beginTime = calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)
			break
		}
	}
	UserBeginTime = beginTime
	endTime := time.Now()
	// 计算时间差
	createdTime := time.Unix(UserBeginTime, 0)
	createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
	duration := endTime.Sub(createdDayMidnight)
	// 计算相差的天数，即目前第几天
	days := int64(duration.Hours()/24) + 1
	ActivityId := int32(r.UUID / 1000)
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		ActivityId,
	)
	roleReceiveInfo := make(map[int32]bool)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = true
	}

	//需要领取的日期
	needDay := r.UUID % 100
	//状态 1-待领取，2-已领取，3-需要补签，4-不可领取

	//这个是领取情况，如果已经领取了，不能再领
	if roleReceiveInfo[int32(needDay)] {
		return false, errutil.ErrGetAward.Error()
	}
	//判断时间是否超过当前日期
	if int64(needDay) > days {
		return false, errutil.ErrChooseRound.Error()
	}

	//没有签到，需要补签
	//减去消耗品
	resultc := make([]protocol.SourceInfo, len(result.CompensationConsume))
	if r.Ad > 0 {
		for i, info := range result.CompensationConsume {
			resultc[i].Id = info.Id
			resultc[i].Num = -info.Num
		}
		err2 := CheckSourceNum(resultc, r.Uid, true, constant.ItemAddContinuity)
		if err2 != nil {
			return false, err2.Error()
		}
	}

	//校验背包情况
	result2 := make([]protocol.SourceInfo, len(result.CommonItem))

	for i, info := range result.CommonItem {
		result2[i].Id = info.Id
		result2[i].Num = info.Num
	}
	err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddContinuity)
	if err3 != nil {
		return false, err3.Error()
	}

	_, err = db.ActivityAddDay(
		r.Uid,
		r.UUID,
		1,
		ActivityId,
	)
	if err != nil {
		return false, ""
	}
	//领取奖励
	return true, ""
}
func activityListCumulative(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityListResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	buys := int64(1)
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	request2.StoreType = 6
	request2.GoodsUuid = 6*1000 + r.ActivityId
	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	lists, err2 := StoreList(request2)
	if err2 != nil || lists == nil {
		return &protocol.ActivityListResponse{Code: -1}, nil
	}
	if lists[0].Buy > 0 {
		buys = 2
	}

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityListResponse{}, errutil.ErrInvalidParameter
	}
	list, days, err := ActivityListCumulative(request)
	if err != nil {
		return &protocol.ActivityListResponse{}, err
	}
	return &protocol.ActivityListResponse{Data: list, Day: days, Buy: buys}, nil
}

func ActivityListCumulative(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListDay, int64, error) {
	actType := constant.ActivityTypeCumulative
	activityCfg := global.ActivityCfg
	result := make([]protocol.ActivityDetailDay, 0)
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				result = append(result, protocol.ActivityDetailDay{
					UUID:    int64(activityn.UUID),
					Content: activityn.Content,
				})

			}
		}
	}
	// 计算开始时间
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, 0, err
	}
	var beginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId && activityn.Type == actType {
			beginTime = calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)
			break
		}
	}

	resulta, err := db.RoleLoginNum(
		r.Uid,
		beginTime,
		time.Now().Unix(),
	)
	if err != nil {
		return nil, 0, errutil.ErrDBOperation
	}

	days := int64(len(resulta))
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		r.ActivityId,
	)
	resultEnd := make([]protocol.ActivityListDay, 0)
	//领取的都是2
	for i, item := range resultd {
		resultEnd = append(resultEnd, protocol.ActivityListDay{
			Day:    int32(i) + 1,
			Status: 2,
			Buy:    item.Buy,
		})
	}

	return resultEnd, days, nil
}
func activityReceiveCumulative(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Buy = 1
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	request2.StoreType = constant.ShopTypeSignin30Days
	request2.GoodsUuid = constant.ShopTypeSignin30Days*1000 + r.ActivityId
	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	list, err2 := StoreList(request2)
	if err2 != nil || list == nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err2.Error()}, nil
	}
	//这里需要判断是否购买vip，没有购买也会是空
	if len(list) > 0 {
		//购买了，后面需要判断vip礼包
		if list[0].Buy > 0 {
			request.Buy = 2
		}
	}
	Cumulative, err := ActivityReceiveCumulative(request)
	if !Cumulative {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// 一键签到
func ActivityReceiveCumulative(r *protocol.ActivityOperateRequest) (bool, string) {
	actType := constant.ActivityTypeCumulative

	activityCfg := global.ActivityCfg
	result := []protocol.ActivityDetailDay{}
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				item := protocol.ActivityDetailDay{
					UUID:                activityn.UUID,
					Content:             activityn.Content,
					CompensationConsume: make([]protocol.SourceInfo, len(activityn.CompensationConsume)),
					CommonItem:          make([]protocol.SourceInfo, len(activityn.CommonItem)),
					PayItem:             make([]protocol.SourceInfo, len(activityn.PayItem)),
				}
				for j, attachment := range activityn.CompensationConsume {
					// 假设 Attachment 结构体中有可以复制的字段
					item.CompensationConsume[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment.Key,
						Num: int64(attachment.Value),
					}
				}

				for j, attachment2 := range activityn.CommonItem {
					// 假设 Attachment 结构体中有可以复制的字段
					item.CommonItem[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment2.Key,
						Num: int64(attachment2.Value),
					}
				}

				for j, attachment3 := range activityn.PayItem {
					// 假设 Attachment 结构体中有可以复制的字段
					item.PayItem[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment3.Key,
						Num: int64(attachment3.Value),
					}
				}
				result = append(result, item)
				// 复制 Attachments 字段
			}
		}
	}
	// 计算开始时间
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err.Error()
	}
	var beginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId && activityn.Type == actType {
			beginTime = calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)
			break
		}
	}

	resulta, _ := db.RoleLoginNum(
		r.Uid,
		beginTime,
		time.Now().Unix(),
	)
	days := int64(len(resulta))

	ActivityId := r.ActivityId
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		ActivityId,
	)
	roleReceiveInfo := make(map[int32]int)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = item.Buy
	}

	//状态 1-待领取，2-已领取，3-不可领取，4-需要补签
	for _, item := range result {
		//需要领取的日期
		needDay := item.UUID % 100

		//判断时间是否超过当前日期
		if int64(needDay) > days {
			continue
		}
		if roleReceiveInfo[int32(needDay)] > 0 {
			//领取了且当初也是vip，那就跳过
			if roleReceiveInfo[int32(needDay)] > 1 {
				continue
			}
			//领取了且当初不是vip，补发vip的
			if r.Buy > 1 {
				//校验背包情况
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddContinuity)
				if err3 != nil {
					continue
				}

				_, err := db.ActivityAddDayUpdate(
					r.Uid,
					item.UUID,
					2,
					int(r.Buy),
				)
				if err != nil {
					continue
				}
			}
		} else {
			//没领取，那么发，但是要校验是不是vip
			if r.Buy > 1 {
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddCumulative)
				if err3 != nil {
					continue
				}
			}
			result3 := make([]protocol.SourceInfo, len(item.CommonItem))
			for i, info := range item.CommonItem {
				result3[i].Id = info.Id
				result3[i].Num = info.Num
			}
			err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddCumulative)
			if err3 != nil {
				continue
			}
			_, err := db.ActivityAddDay(
				r.Uid,
				item.UUID,
				int(r.Buy),
				r.ActivityId,
			)
			if err != nil {
				continue
			}
		}

	}

	//领取奖励
	return true, ""
}

func activityListCharge(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityChargeResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityChargeResponse{}, errutil.ErrInvalidParameter
	}
	chargeList, days, err := ActivityListCharge(request)
	if err != nil {
		if err == errutil.ErrActivityExpire {
			return &protocol.ActivityChargeResponse{Code: -1, Data: []protocol.ActivityListDay{}}, nil
		}
		return &protocol.ActivityChargeResponse{}, err
	}
	// 为了兼容性，保留 Status 字段（取最后一天的状态，如果列表为空则为3）
	status := 3
	if len(chargeList) > 0 {
		status = int(chargeList[len(chargeList)-1].Status)
	}
	return &protocol.ActivityChargeResponse{Data: chargeList, Day: days, Status: status}, nil
}

// 领取情况
func ActivityListCharge(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListDay, int64, error) {
	actType := constant.ActivityTypeCharge
	days := int64(1)
	activityCfg := global.ActivityCfg
	ResetTime := int32(0)
	//活动内容
	goodsId := make([]int64, 0)
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				goodsId = append(goodsId, int64(activityn.PreLimit))
				ResetTime = activityn.ResetTime
			}
		}
	}
	StoType := int32(goodsId[0] / 1000)

	//查询用户购买情况
	resultb, _ := db.StoreList(
		r.Uid,
		StoType,
		int32(365*50),
		goodsId,
		false,
	)
	// 计算开始时间
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, 0, err
	}
	var beginTime int64 = 0
	var UserBeginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId && activityn.Type == actType {
			beginTime = calculateBeginTime(activityn.StartTime, res.RegisterAt, res.ServerId)
			break
		}
	}
	if beginTime > 0 {
		UserBeginTime = beginTime
	} else {
		if len(resultb) > 0 {
			timeObj := time.Unix(resultb[0].CreateTime, 0)
			midnight := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 0, 0, 0, 0, timeObj.Location())
			UserBeginTime = midnight.Unix()
		} else {
			// 没有购买记录，返回空列表
			return []protocol.ActivityListDay{}, 1, nil
		}
	}
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		r.ActivityId,
	)

	endTime := time.Now()
	duration := endTime.Sub(time.Unix(int64(UserBeginTime), 0))
	if UserBeginTime > 0 {
		days = int64(duration.Hours()/24) + 1
	}
	if UserBeginTime+int64(ResetTime*3600*24) < endTime.Unix() {
		return nil, 0, errutil.ErrActivityExpire
	}

	// 构建所有购买记录的天数映射（用于判断是否可双倍领取）
	resultbInfo := make(map[int32]bool)
	for _, item := range resultb {
		// 计算时间差
		duration1 := time.Unix(int64(item.CreateTime), 0).Sub(time.Unix(int64(UserBeginTime), 0))
		// 计算相差的天数，即目前第几天
		days1 := int64(math.Abs(duration1.Hours()/24)) + 1
		resultbInfo[int32(days1)] = true
	}

	// 构建所有领取记录的天数映射
	resultdInfo := make(map[int32]int)
	for _, item := range resultd {
		resultdInfo[item.Day] = item.Buy
	}

	// 生成所有天的状态列表
	resultList := make([]protocol.ActivityListDay, 0)
	firstDayUnlocked := false
	if claimedBuy, ok := resultdInfo[1]; ok && claimedBuy > 0 {
		firstDayUnlocked = true
	} else if resultbInfo[1] {
		firstDayUnlocked = true
	}

	for day := int32(1); day <= int32(days); day++ {
		var status int
		buy := 1 // 默认普通领取

		if day == 1 {
			if claimedBuy, ok := resultdInfo[day]; ok {
				firstDayUnlocked = true
				status = 2
				if claimedBuy > 0 {
					buy = claimedBuy
				}
			} else if resultbInfo[day] {
				firstDayUnlocked = true
				status = 2 // 首日购买后视为已领取
				buy = 2    // 首日购买属于双倍领取
			} else {
				status = 3 // 首日无免费领取
			}
		} else {
			if !firstDayUnlocked {
				status = 3 // 首日未解锁，后续不可领取
			} else {
				if claimedBuy, ok := resultdInfo[day]; ok {
					status = 2 // 已领取
					if claimedBuy > 0 {
						buy = claimedBuy
					}
				} else if resultbInfo[day] {
					// 当天付费，视为已领取
					status = 2
					buy = 2
				} else {
					// 未领取，也未付费，维持免费可领取状态
					status = 1
				}
			}
		}
		resultList = append(resultList, protocol.ActivityListDay{
			Day:    day,
			Status: status,
			Buy:    buy,
		})
	}

	return resultList, days, nil
}

func activityReceiveCharge(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {

	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}

	chargeList, days, err := ActivityListCharge(request)
	if err != nil {
		return &protocol.ActivityReceiveDayResponse{}, err
	}
	//状态 1-待领取，2-已领取，3-不可领取（没充值）
	// 获取当天的状态（列表最后一项）
	status := 3
	if len(chargeList) > 0 {
		status = int(chargeList[len(chargeList)-1].Status)
	}
	if status != 1 {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: "状态不对"}, nil
	}
	// 设置当天领取的类型（1-普通，2-双倍）
	request.Buy = int64(chargeList[len(chargeList)-1].Buy)
	request.UUID = int64(5000) + int64(100*r.ActivityId) + days

	Cumulative, err2 := ActivityReceiveCharge(request)
	if !Cumulative {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err2}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// 首充领取
func ActivityReceiveCharge(r *protocol.ActivityOperateRequest) (bool, string) {
	activityCfg := global.ActivityCfg
	//result := make(map[int32][]protocol.SourceInfo)
	result := protocol.ActivityDetailDay{}
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.UUID == r.UUID {

			item := protocol.ActivityDetailDay{
				UUID:                activityn.UUID,
				Content:             activityn.Content,
				CompensationConsume: make([]protocol.SourceInfo, len(activityn.CompensationConsume)),
				CommonItem:          make([]protocol.SourceInfo, len(activityn.CommonItem)),
				PayItem:             make([]protocol.SourceInfo, len(activityn.PayItem)),
			}
			for j, attachment := range activityn.CompensationConsume {
				// 假设 Attachment 结构体中有可以复制的字段
				item.CompensationConsume[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment.Key,
					Num: int64(attachment.Value),
				}
			}

			for j, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.CommonItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}

			for j, attachment3 := range activityn.PayItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.PayItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment3.Key,
					Num: int64(attachment3.Value),
				}
			}
			result = item
			// 复制 Attachments 字段
		}
	}

	result3 := make([]protocol.SourceInfo, len(result.CommonItem))
	for i, info := range result.CommonItem {
		result3[i].Id = info.Id
		result3[i].Num = info.Num
	}
	err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddFirstCharge)
	if err3 != nil {
		return false, err3.Error()
	}
	_, err := db.ActivityAddDay(
		r.Uid,
		result.UUID,
		int(r.Buy),
		r.ActivityId,
	)
	if err != nil {
		return false, err.Error()
	}

	//领取奖励
	return true, ""
}

// 等级奖励列表
func activityListGrade(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityGradeResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityGradeResponse{}, errutil.ErrInvalidParameter
	}
	res, err := ActivityListGrade(request)
	if err != nil {
		if err == errutil.ErrActivityExpire {
			return &protocol.ActivityGradeResponse{Code: -1, Data: errutil.ErrActivityExpire.Error()}, nil
		}
		return &protocol.ActivityGradeResponse{}, err
	}
	return &protocol.ActivityGradeResponse{Data: "成功", ActivityListGrade: res}, nil
}
func ActivityListGrade(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListGrade, error) {
	actType := constant.ActivityTypeGrade
	activityCfg := global.ActivityCfg

	//活动内容
	goodsId := make([]int64, 0)
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				if r.UUID > 0 {
					if r.UUID == activityn.UUID {
						goodsId = append(goodsId, activityn.UUID)
					}
				} else {
					goodsId = append(goodsId, activityn.UUID)
				}
			}
		}
	}

	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	getInfo := make(map[int64]bool)

	if len(resultb) > 0 {
		for _, item := range resultb {
			getInfo[item.GoodsUuiD] = true
		}
	}

	ActivityListGradeAll := make([]protocol.ActivityListGrade, len(goodsId))
	for i, item := range goodsId {
		GradeInfo, err := global.GetActivityCfg(int32(item))
		if err != nil {
			return nil, err
		}
		GetStatus := 1
		if getInfo[item] {
			GetStatus = 2
		}
		//比较等级
		res, err := game.UserInfo(r.Uid)
		if err != nil {
			return nil, err
		}
		if GradeInfo.PreLimit > res.Level {
			GetStatus = 3
		}
		ActivityListGradeAll[i] = protocol.ActivityListGrade{
			Grade:  GradeInfo.PreLimit,
			Status: GetStatus,
		}
	}
	return ActivityListGradeAll, nil
}

// 领取等级奖励
func activityReceiveGrade(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.UUID = r.UUID
	if request.UUID == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	_, err := ActivityReceiveGade(request)
	if err != nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}
func ActivityReceiveGade(r *protocol.ActivityOperateRequest) (bool, error) {

	activityCfg := global.ActivityCfg
	result := protocol.ActivityDetailDay{}
	//活动内容
	for _, activityn := range activityCfg {
		if activityn.UUID == r.UUID {
			item := protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				PreLimit:   activityn.PreLimit,
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
			}
			for j, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.CommonItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}
			result = item
		}
	}
	//需要领取物品的等级
	//比较等级
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err
	}
	if res.Level < result.PreLimit {
		return false, errutil.ErrLevel
	}
	goodsId := make([]int64, 0)
	goodsId = append(goodsId, result.UUID)
	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	if len(resultb) > 0 {
		return false, errutil.ErrGetAward
	}
	result3 := make([]protocol.SourceInfo, len(result.CommonItem))
	for i, info := range result.CommonItem {
		result3[i].Id = info.Id
		result3[i].Num = info.Num
	}
	//领取奖励
	err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddActivityLevel)
	if err3 != nil {
		return false, err3
	}
	_, err = db.GradeAdd(
		r.Uid,
		int32(result.UUID),
		result3,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func activityListTurntable(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityTruntableResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityTruntableResponse{}, errutil.ErrInvalidParameter
	}
	content, err := ActivityListTurntable(request)
	if err != nil {
		return &protocol.ActivityTruntableResponse{Code: -1, Content: -1}, nil
	}
	return &protocol.ActivityTruntableResponse{Content: content}, nil
}

func ActivityListTurntable(r *protocol.ActivityOperateRequest) (int32, error) {

	actType := constant.ActivityTypeTrunTable
	activityCfg := global.ActivityCfg
	//result := make(map[int32][]protocol.SourceInfo)
	result := make([]protocol.ActivityDetailDay, 0)
	ResetTime := int64(0)
	//活动内容
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				ResetTime = int64(activityn.ResetTime)
				result = append(result, protocol.ActivityDetailDay{
					UUID:     int64(activityn.UUID),
					Content:  activityn.Content,
					PreLimit: activityn.PreLimit,
				})
			}
		}
	}
	// 计算开始时间
	userInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return -1, err
	}
	var beginTime int64 = 0
	var UserBeginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId && activityn.Type == actType {
			beginTime = calculateBeginTime(activityn.StartTime, userInfo.RegisterAt, userInfo.ServerId)
			break
		}
	}

	resulta, err := db.RoleLoginNum(
		r.Uid,
		beginTime,
		time.Now().Unix(),
	)
	if err != nil {
		return -1, errutil.ErrDBOperation
	}
	roleLoginInfo := make(map[string]string)
	for _, item := range resulta {
		roleLoginInfo[item.Date] = item.Num
	}
	if beginTime > 0 {
		UserBeginTime = beginTime
	} else {
		if len(resulta) > 0 {
			t, err := time.Parse("2006-01-02", resulta[0].Date)
			if err != nil {
				return -1, err
			}
			UserBeginTime = t.Unix()
		} else {
			UserBeginTime = userInfo.RegisterAt
		}
	}
	endTime := time.Now()
	// 计算时间差
	createdTime := time.Unix(UserBeginTime, 0)
	createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
	duration := endTime.Sub(createdDayMidnight)
	// 计算相差的天数
	days := int64(duration.Hours() / 24)
	if days > ResetTime-1 {
		return -1, errutil.ErrActivityExpire
	}
	//查询用户目领取情况
	resultd, _ := db.ActivityListTurntable(
		r.Uid,
		r.ActivityId,
	)
	content := int32(0)
	for _, item := range resultd {
		if item.Content > content {
			content = item.Content
		}
	}

	uuid := actType*1000 + r.ActivityId*100 + content + 1
	//uuid对应信息
	ActRes, err := global.GetActivityCfg(uuid)
	if err != nil {
		return -1, nil
	}
	//个人等级比较
	if userInfo.Level < ActRes.PreLimit && r.Check == 1 {
		return -1, nil
	} else {
		return uuid, nil
	}

}

func activityReceiveTurntable(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveTurntableResponse, error) {

	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Check = 1
	Cumulative, err := ActivityReceiveTruntable(request)
	if err != "" {
		return &protocol.ActivityReceiveTurntableResponse{Code: -1, Data: err}, nil
	}
	return &protocol.ActivityReceiveTurntableResponse{Puid: Cumulative}, nil
}

// 抽奖，给到抽到什么
func ActivityReceiveTruntable(r *protocol.ActivityOperateRequest) (int32, string) {
	ress, err := ActivityListTurntable(r)
	if err != nil {
		return -1, err.Error()
	}
	activityCfg := global.ActivityCfg
	result := protocol.ActivityDetailDay{}
	//活动内容
	ResetTime := int64(0)
	for _, activityn := range activityCfg {

		if activityn.UUID == int64(ress) {
			ResetTime = int64(activityn.ResetTime)
			result = protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				PayItem:    make([]protocol.SourceInfo, len(activityn.PayItem)),
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
			}
			for i, attachment := range activityn.PayItem {
				// 假设 Attachment 结构体中有可以复制的字段
				result.PayItem[i] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment.Key,
					Num: int64(attachment.Value),
				}
			}

			for i, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				result.CommonItem[i] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}
			// 复制 Attachments 字段

		}
	}
	// 计算开始时间
	userInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return -1, err.Error()
	}
	var beginTime int64 = 0
	var UserBeginTime int64 = 0
	for _, activityn := range activityCfg {
		if activityn.UUID == int64(ress) {
			beginTime = calculateBeginTime(activityn.StartTime, userInfo.RegisterAt, userInfo.ServerId)
			break
		}
	}
	resulta, err := db.RoleLoginNum(
		r.Uid,
		beginTime,
		time.Now().Unix(),
	)
	if err != nil {
		return -1, errutil.ErrDBOperation.Error()
	}
	// roleLoginInfo := make(map[string]string)
	// for _, item := range resulta {
	// 	roleLoginInfo[item.Date] = item.Num
	// }
	if beginTime > 0 {
		UserBeginTime = beginTime
	} else {
		if len(resulta) == 0 {
			return -1, err.Error()
		}
		t, err := time.Parse("2006-01-02", resulta[0].Date)
		if err != nil {
			return -1, err.Error()
		}
		UserBeginTime = t.Unix()
	}
	endTime := time.Now()
	// 计算时间差
	duration := endTime.Sub(time.Unix(int64(UserBeginTime), 0))
	// 计算相差的天数
	days := int64(duration.Hours() / 24)
	if days > ResetTime-1 {
		return -1, errutil.ErrActivityExpire.Error()
	}
	//需要领取的日期
	needContent := ress % 100

	//校验背包情况
	result2 := make([]protocol.SourceInfo, len(result.CommonItem))

	for i, info := range result.CommonItem {
		result2[i].Id = info.Id
		result2[i].Num = info.Num
	}
	err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddTurntable)
	if err3 != nil {
		return -1, err3.Error()
	}
	//减去消耗品
	resultc := make([]protocol.SourceInfo, len(result.PayItem))
	sendId := int32(2)
	for i, info := range result.PayItem {
		sendId = info.Id
		resultc[i].Id = info.Id
		resultc[i].Num = -info.Num
	}
	err = CheckSourceNum(resultc, r.Uid, true, constant.ItemAddTurntable)
	if err != nil {
		return -1, err.Error()
	}
	truntableRes, puid := TruntableDetail(resultc, int32(needContent))
	if truntableRes == -1 {
		return -1, ""
	}
	resulte := make([]protocol.SourceInfo, 1)
	resulte[0].Id = sendId
	resulte[0].Num = int64(math.Abs(truntableRes))
	err = CheckSourceNum(resulte, r.Uid, true, constant.ItemAddTurntable)
	if err != nil {
		return -1, err.Error()
	}
	_, err = db.ActivityAddTurntable(
		r.Uid,
		int64(ress),
		int32(truntableRes),
	)
	if err != nil {
		return -1, ""
	}
	//个人信息
	// userInfo, err := game.UserInfo(r.Uid)
	// if err != nil {
	// 	return -1, err.Error()
	// }
	// strSend := []string{userInfo.Name, strconv.Itoa(int(math.Abs(res)))}
	// game.Lamp(int32(constant.LampLuckyPoints), strSend)
	//领取奖励
	return puid, ""
}

// 大转盘抽奖
func TruntableDetail(sourceInfo []protocol.SourceInfo, round int32) (float64, int32) {
	prizeWheelCfg := global.PrizeWheelCfg

	allPriority := 0
	beginPriority := 0

	//活动内容
	for _, prizeWheelitem := range prizeWheelCfg {
		if prizeWheelitem.Round == round {
			allPriority = allPriority + int(prizeWheelitem.Priority)
		}
	}
	rand := global.RandIn(allPriority) + 1
	for _, prizeWheelitem := range prizeWheelCfg {
		if prizeWheelitem.Round == round {
			if beginPriority+int(prizeWheelitem.Priority) >= rand {
				return float64(sourceInfo[0].Num) * prizeWheelitem.Rate, prizeWheelitem.Uid
			}
			beginPriority = beginPriority + int(prizeWheelitem.Priority)
		}
	}
	return -1, 0
}

// 判断活动是否已经结束了
func checkTime(uid int64, beginTime int64) (map[string]string, int64, int64, string) {
	resulta, _ := db.RoleLoginNum(
		uid,
		beginTime,
		time.Now().Unix(),
	)
	roleLoginInfo := make(map[string]string)
	for _, item := range resulta {
		roleLoginInfo[item.Date] = item.Num
	}
	UserBeginTime := int64(0)

	if beginTime > 0 {
		UserBeginTime = beginTime
	} else {
		t, err := time.Parse("2006-01-02", resulta[0].Date)
		if err != nil {
			return roleLoginInfo, -1, -1, "转换出错"
		}
		UserBeginTime = t.Unix()
	}
	endTime := time.Now()
	// 计算时间差
	createdTime := time.Unix(UserBeginTime, 0)
	createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
	duration := endTime.Sub(createdDayMidnight)
	// 计算相差的天数，即目前第几天
	days := int64(duration.Hours()/24) + 1
	return roleLoginInfo, UserBeginTime, days, ""
}
func activityReceiveTreasureBox(ctx context.Context, r *protocol.ActivityTreasureBoxRequest) (*protocol.ActivityTreasureBoxResponse, error) {
	request := &protocol.ActivityTreasureBoxRequest{}
	//[{2024-05-09 2} {202405-10 1}]
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request.Uid = Uid

	if r.Type != 1 && r.Type != 2 {
		return &protocol.ActivityTreasureBoxResponse{Code: -1}, nil
	}
	request.Type = r.Type
	_, position, err := ActivityReceiveTreasureBox(request)
	if err != nil {
		return &protocol.ActivityTreasureBoxResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.ActivityTreasureBoxResponse{Position: position}, nil
}

// 百宝箱抽奖
func ActivityReceiveTreasureBox(r *protocol.ActivityTreasureBoxRequest) ([]protocol.SourceInfo, []int32, error) {
	//个人等级
	uerInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, nil, err
	}

	//活动内容
	res := global.GetTreasureBox(uerInfo.Level)
	if res == nil {
		return nil, nil, errutil.ErrConfigNotFound
	}

	// //需要领取的日期
	// needContent := ress % 100
	needNum := -1
	switch r.Type {
	case 1:
		needNum = 1
	case 2:
		needNum = 10
	}
	if needNum == -1 {
		return nil, nil, errutil.ErrConfigNotFound
	}
	//减去消耗品
	resultc := make([]protocol.SourceInfo, 1)
	for _, info := range res.Consume {
		if r.Type == info.Type {
			resultc[0].Id = info.Id
			resultc[0].Num = -int64(info.Num)
		}

	}
	//校验背包情况
	err = CheckSourceNum(resultc, r.Uid, true, constant.ItemAddTurntable)
	if err != nil {
		return nil, nil, err
	}

	haveNum, err := db.ActivityListTreasureBox(r.Uid)
	if err != nil {
		return nil, nil, err
	}
	randInfo, position, err := TreasureBoxDetail(uerInfo, needNum, int(haveNum))
	if err != nil {
		return nil, nil, err
	}

	result2 := make([]protocol.SourceInfo, len(randInfo))

	for i, info := range randInfo {
		result2[i].Id = info.Id
		result2[i].Num = info.Num
	}

	err = CheckSourceAdd(result2, r.Uid, true, constant.ItemAddTurntable)
	if err != nil {
		if err == errutil.ErrPlayerNotFound {
			return nil, nil, err
		}
		if err == errutil.ErrBagFull {

			MailId := strconv.FormatInt(snowflake.GenID(), 10)
			result3 := make([]protocol.Attachment, len(result2))
			for i, Attachment := range result2 {
				result3[i] = protocol.Attachment{
					MailId:      MailId,
					UserPropsId: int64(Attachment.Id),
					Num:         Attachment.Num,
				}
			}
			request3 := &protocol.MailAddRequest{
				Uid:          r.Uid,
				RecipientUid: -1,
				Title:        "系统邮件",
				Desc:         "包裹已满，请及时清理",
				Attachment:   result3,
				ChannelCode:  "0",
			}
			resultc, err := Mailadd(request3)
			if !resultc || err != nil {
				return nil, nil, err
			}
		}

	}
	jsonDate, err := json.Marshal(result2)
	if err != nil {
		return nil, nil, err
	}

	jsonDateString := string(jsonDate)
	err = db.ActivityAddTreasureBox(r.Uid, jsonDateString, int32(needNum))
	if err != nil {
		return nil, nil, err
	}
	// return int32(resulte[0].Num), ""
	return result2, position, nil
}

// 百宝箱抽奖
func TreasureBoxDetail(userInfo *protocol.UserInfo, needNum int, haveNum int) ([]protocol.SourceInfo, []int32, error) {
	position := make([]int32, 0)

	treasureBoxlCfg := global.GetTreasureBox(userInfo.Level)
	if needNum+haveNum > int(treasureBoxlCfg.MaxTime) {
		return nil, position, errutil.ErrNotCanBuy
	}
	CommonPriority := 0
	SpecialPriority := 0
	resultc := make([]protocol.SourceInfo, 0)

	//普通奖池
	for _, treasureBoxcitem := range treasureBoxlCfg.CommonPool {
		CommonPriority = CommonPriority + int(treasureBoxcitem.Priority)
	}
	//必中奖池
	for _, treasureBoxsitem := range treasureBoxlCfg.SpecialPool {
		SpecialPriority = SpecialPriority + int(treasureBoxsitem.Priority)
	}
	for i := 1; i <= needNum; i++ {
		fa := make(map[int]bool)
		faBegin := 0
		for _, item := range treasureBoxlCfg.Node {
			faBegin = faBegin + item
			fa[faBegin] = true
		}
		if fa[i+int(haveNum)] {
			beginPriority := 0
			rand := global.RandIn(SpecialPriority) + 1
			for _, prizeWheelitem := range treasureBoxlCfg.SpecialPool {
				if beginPriority+int(prizeWheelitem.Priority) >= rand {
					resultc = append(resultc, protocol.SourceInfo{
						Id:  prizeWheelitem.Id,
						Num: int64(prizeWheelitem.Num),
					})
					// itemInfo, err := global.GetItemCfg(prizeWheelitem.Id)
					// if err != nil {
					// 	return nil, nil, nil
					// }
					// if itemInfo.Quality >= 5 {
					// 	strSend := []string{userInfo.Name, strconv.Itoa(int(itemInfo.ID))}
					// 	game.Lamp(int32(constant.LampObtainedTreasure), strSend)

					// }
					position = append(position, prizeWheelitem.Position)
					break
				}
				beginPriority = beginPriority + int(prizeWheelitem.Priority)
			}
		} else {
			beginPriority := 0
			rand := global.RandIn(CommonPriority) + 1
			for _, prizeWheelitem := range treasureBoxlCfg.CommonPool {
				if beginPriority+int(prizeWheelitem.Priority) >= rand {
					resultc = append(resultc, protocol.SourceInfo{
						Id:  prizeWheelitem.Id,
						Num: int64(prizeWheelitem.Num),
					})
					// itemInfo, err := global.GetItemCfg(prizeWheelitem.Id)
					// if err != nil {
					// 	return nil, nil, nil
					// }
					// if itemInfo.Quality >= 5 {
					// 	strSend := []string{userInfo.Name, strconv.Itoa(int(itemInfo.ID))}
					// 	game.Lamp(int32(constant.LampObtainedTreasure), strSend)
					// }
					position = append(position, prizeWheelitem.Position)
					break
				}
				beginPriority = beginPriority + int(prizeWheelitem.Priority)
			}
		}

	}
	return resultc, position, nil
}

func activityListTreasurebox(ctx context.Context) (*protocol.ActivityListTreasureBoxResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	residue, residueBig, err := ActivityListTreasurebox(request)
	if err != nil {
		return &protocol.ActivityListTreasureBoxResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.ActivityListTreasureBoxResponse{CanPlay: residue, NeedNum: residueBig}, nil
}

// 百宝箱抽奖
func ActivityListTreasurebox(r *protocol.ActivityOperateRequest) (int32, int32, error) {
	//个人等级
	uerInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return 0, 0, err
	}

	haveNum, err := db.ActivityListTreasureBox(r.Uid)
	if err != nil {
		return 0, 0, err
	}
	treasureBoxlCfg := global.GetTreasureBox(uerInfo.Level)
	residue := float64(treasureBoxlCfg.MaxTime) - haveNum

	residueBig := float64(0)
	faBegin := float64(haveNum)
	for _, item := range treasureBoxlCfg.Node {
		faBegin = faBegin - float64(item)
		if faBegin >= 0 {
			continue
		} else {
			residueBig = math.Abs(faBegin)
			break
		}
	}
	return int32(residue), int32(residueBig), nil
}
func activityListTrial(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityListResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	buys := int64(1)
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	request2.StoreType = constant.ShopTypeTrialReward
	request2.GoodsUuid = constant.ShopTypeTrialReward*1000 + r.ActivityId
	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	lists, err2 := StoreList(request2)
	if err2 != nil || lists == nil {
		return &protocol.ActivityListResponse{Code: -1}, nil
	}
	if len(lists) == 0 {
		return &protocol.ActivityListResponse{Code: -1}, nil
	}
	if lists[0].Buy > 0 {
		buys = 2
	}

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityListResponse{}, errutil.ErrInvalidParameter
	}
	list, err := ActivityListtrial(request)
	if err != nil {
		return &protocol.ActivityListResponse{}, err
	}
	res, err := game.UserInfo(Uid)
	if err != nil {
		return &protocol.ActivityListResponse{}, err
	}
	days := int64(res.TrialTowerLevel)
	return &protocol.ActivityListResponse{Data: list, Buy: buys, Day: days}, nil
}

func ActivityListtrial(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListDay, error) {
	actType := constant.ActivityTypeTrial
	//活动内容

	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		r.ActivityId,
	)
	resultEnd := make([]protocol.ActivityListDay, 0)
	//领取的都是2
	for _, item := range resultd {
		resultEnd = append(resultEnd, protocol.ActivityListDay{
			Day:    item.Day,
			Status: 2,
			Buy:    item.Buy,
		})
	}

	return resultEnd, nil
}

func activityReceiveTrial(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Buy = 1
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	request2.StoreType = constant.ShopTypeTrialReward
	request2.GoodsUuid = constant.ShopTypeTrialReward*1000 + r.ActivityId
	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	list, err2 := StoreList(request2)
	if err2 != nil || list == nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err2.Error()}, nil
	}
	// 检查是否有购买记录，没有购买记录也能往下走，
	if len(list) > 0 {
		//有购买记录会加一个标志
		if list[0].Buy > 0 {
			request.Buy = 2
		}
	}
	Cumulative, err := ActivityReceiveTrial(request)
	if !Cumulative {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// 一键领取
func ActivityReceiveTrial(r *protocol.ActivityOperateRequest) (bool, string) {
	actType := constant.ActivityTypeTrial
	//查询用户登录情况
	activityCfg := global.ActivityCfg
	result := []protocol.ActivityDetailDay{}
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {

				item := protocol.ActivityDetailDay{
					UUID:                activityn.UUID,
					Content:             activityn.Content,
					CompensationConsume: make([]protocol.SourceInfo, len(activityn.CompensationConsume)),
					CommonItem:          make([]protocol.SourceInfo, len(activityn.CommonItem)),
					PayItem:             make([]protocol.SourceInfo, len(activityn.PayItem)),
				}
				for j, attachment := range activityn.CompensationConsume {
					// 假设 Attachment 结构体中有可以复制的字段
					item.CompensationConsume[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment.Key,
						Num: int64(attachment.Value),
					}
				}

				for j, attachment2 := range activityn.CommonItem {
					// 假设 Attachment 结构体中有可以复制的字段
					item.CommonItem[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment2.Key,
						Num: int64(attachment2.Value),
					}
				}

				for j, attachment3 := range activityn.PayItem {
					// 假设 Attachment 结构体中有可以复制的字段
					item.PayItem[j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  attachment3.Key,
						Num: int64(attachment3.Value),
					}
				}
				result = append(result, item)
				// 复制 Attachments 字段
			}
		}
	}

	//
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err.Error()
	}
	days := int64(res.TrialTowerLevel)
	ActivityId := r.ActivityId
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		ActivityId,
	)
	roleReceiveInfo := make(map[int32]int)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = item.Buy
	}

	//状态 1-待领取，2-已领取，3-不可领取，4-需要补签
	for _, item := range result {
		//需要领取的日期
		needDay := item.UUID - int64(r.ActivityId*100) - int64(actType*1000)

		//判断时间是否超过当前日期if int64(needDay) > int64(days) {
		if int64(needDay) > days {
			continue
		}
		if roleReceiveInfo[int32(needDay)] > 0 {
			//领取了且当初也是vip，那就跳过
			if roleReceiveInfo[int32(needDay)] > 1 {
				continue
			}
			//领取了且当初不是vip，补发vip的
			if r.Buy > 1 {
				//校验背包情况
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddCumulative)
				if err3 != nil {
					continue
				}

				_, err := db.ActivityAddDayUpdate(
					r.Uid,
					item.UUID,
					2,
					int(r.Buy),
				)
				if err != nil {
					continue
				}
			}
		} else {
			//没领取，那么发，但是要校验是不是vip
			if r.Buy > 1 {
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddCumulative)
				if err3 != nil {
					continue
				}
			}
			result3 := make([]protocol.SourceInfo, len(item.CommonItem))
			for i, info := range item.CommonItem {
				result3[i].Id = info.Id
				result3[i].Num = info.Num
			}
			err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddCumulative)
			if err3 != nil {
				continue
			}
			_, err := db.ActivityAddDay(
				r.Uid,
				item.UUID,
				int(r.Buy),
				r.ActivityId,
			)
			if err != nil {
				continue
			}
		}

	}

	//领取奖励
	return true, ""
}
func activityListGroupup(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityListResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	buys := int64(1)
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	switch r.Atype {
	case constant.ActivityTypeGrowthFund:
		request2.StoreType = constant.ShopTypeGrowthFund
		request2.GoodsUuid = constant.ShopTypeGrowthFund*1000 + r.ActivityId
	case constant.ActivityTypeLottery:
		request2.StoreType = constant.ShopTypeLotteryFund
		request2.GoodsUuid = constant.ShopTypeLotteryFund*1000 + r.ActivityId
	case constant.ActivityTypeRefine:
		request2.StoreType = constant.ShopTypeRefineFund
		request2.GoodsUuid = constant.ShopTypeRefineFund*1000 + r.ActivityId
	case constant.ActivityTypeInstance:
		request2.StoreType = constant.ShopTypeInstanceFund
		request2.GoodsUuid = constant.ShopTypeInstanceFund*1000 + r.ActivityId
	default:
		r.Atype = constant.ActivityTypeGrowthFund
		request2.StoreType = constant.ShopTypeGrowthFund
		request2.GoodsUuid = constant.ShopTypeGrowthFund*1000 + r.ActivityId
		// return &protocol.ActivityListResponse{Code: -1}, nil
	}

	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	lists, err2 := StoreList(request2)
	if err2 != nil || lists == nil {
		return &protocol.ActivityListResponse{Code: -1}, nil
	}
	if len(lists) > 0 {
		if lists[0].Buy > 0 {
			buys = 2
		}
	}

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Atype = r.Atype
	if request.ActivityId == 0 {
		return &protocol.ActivityListResponse{}, errutil.ErrInvalidParameter
	}
	//比较等级
	ress, err := game.UserInfo(Uid)
	if err != nil {
		return nil, err
	}
	num := int64(0)
	switch r.Atype {
	case constant.ActivityTypeGrowthFund:
		num = int64(ress.Level)
	case constant.ActivityTypeLottery:
		num = ress.LotteryNum
	case constant.ActivityTypeRefine:
		num = ress.RefineNum
	case constant.ActivityTypeInstance:
		num = ress.InstanceNum
	default:
	}
	// //活动时间过了就不给他玩了
	// for _, activityn := range global.ActivityCfg {
	// 	if activityn.Acttag == r.ActivityId && activityn.Type == r.Atype {
	// 		if activityn.Type == actType {
	// 		}
	// 	}
	// }

	list, err := ActivityListGroupup(request)
	if err != nil {
		return &protocol.ActivityListResponse{}, err
	}
	return &protocol.ActivityListResponse{Data: list, Buy: buys, Num: num}, nil
}

func ActivityListGroupup(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListDay, error) {
	actType := r.Atype
	activityCfg := global.ActivityCfg
	result := make([]protocol.ActivityDetailDay, 0)
	// UserBeginTime := int64(0)
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				if activityn.StartTime > 0 {
					// beginTime = int64(activityn.StartTime)
				}
				result = append(result, protocol.ActivityDetailDay{
					UUID:    int64(activityn.UUID),
					Content: activityn.Content,
				})

			}
		}
	}

	// resulta, err := db.RoleLoginNum(
	// 	r.Uid,
	// 	beginTime,
	// 	time.Now().Unix(),
	// )
	// if err != nil {
	// 	return nil, 0, errutil.ErrDBOperation
	// }

	// days := int64(len(resulta))
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		r.ActivityId,
	)
	resultEnd := make([]protocol.ActivityListDay, 0)
	//领取的都是2
	for _, item := range resultd {
		resultEnd = append(resultEnd, protocol.ActivityListDay{
			Day:    item.Day,
			Status: 2,
			Buy:    item.Buy,
		})
	}

	return resultEnd, nil
}
func activityReceiveGroupup(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Buy = 1
	request.Atype = r.Atype
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	request2 := &protocol.StoreListRequest{}
	request2.Uid = Uid
	switch r.Atype {
	case constant.ActivityTypeGrowthFund:
		request2.StoreType = constant.ShopTypeGrowthFund
		request2.GoodsUuid = constant.ShopTypeGrowthFund*1000 + r.ActivityId
	case constant.ActivityTypeLottery:
		request2.StoreType = constant.ShopTypeLotteryFund
		request2.GoodsUuid = constant.ShopTypeLotteryFund*1000 + r.ActivityId
	case constant.ActivityTypeRefine:
		request2.StoreType = constant.ShopTypeRefineFund
		request2.GoodsUuid = constant.ShopTypeRefineFund*1000 + r.ActivityId
	case constant.ActivityTypeInstance:
		request2.StoreType = constant.ShopTypeInstanceFund
		request2.GoodsUuid = constant.ShopTypeInstanceFund*1000 + r.ActivityId
	default:
		r.Atype = constant.ActivityTypeGrowthFund
		request2.StoreType = constant.ShopTypeGrowthFund
		request2.GoodsUuid = constant.ShopTypeGrowthFund*1000 + r.ActivityId
		// return &protocol.ActivityReceiveDayResponse{Code: -1}, nil
	}

	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	list, err2 := StoreList(request2)
	if err2 != nil || list == nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err2.Error()}, nil
	}
	// 检查是否有购买记录，没有购买记录也能往下走，
	if len(list) > 0 {
		//有购买记录会加一个标志
		if list[0].Buy > 0 {
			request.Buy = 2
		}
	}
	Cumulative, err := ActivityReceiveGroupup(request)
	if !Cumulative {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// 一键签到
func ActivityReceiveGroupup(r *protocol.ActivityOperateRequest) (bool, string) {
	//比较等级
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err.Error()
	}
	actType := r.Atype

	activityCfg := global.ActivityCfg
	result := []protocol.ActivityDetailDay{}
	//活动内容
	for _, activityn := range activityCfg {

		if activityn.Acttag == r.ActivityId && activityn.Type == actType {

			if activityn.ResetTime > 0 {
				if res.CreateDay > activityn.ResetTime {
					return false, "活动时间已过"
				}
				// beginTime = int64(activityn.StartTime)
			}
			item := protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
				PayItem:    make([]protocol.SourceInfo, len(activityn.PayItem)),
				PreLimit:   activityn.PreLimit,
			}

			for j, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.CommonItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}

			for j, attachment3 := range activityn.PayItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.PayItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制值
					Id:  attachment3.Key,
					Num: int64(attachment3.Value),
				}
			}
			result = append(result, item)
			// 复制 Attachments 字段
		}
	}

	ActivityId := r.ActivityId
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		r.Uid,
		actType,
		ActivityId,
	)
	roleReceiveInfo := make(map[int32]int)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = item.Buy
	}

	//状态 1-待领取，2-已领取，3-不可领取，4-需要补签
	for _, item := range result {
		//需要领取的日期
		needDay := item.UUID % 100

		//判断时间是否超过当前日期
		// if int64(needDay) > days {
		// 	continue
		// }
		//不同活动，判断条件不同
		switch actType {
		case constant.ActivityTypeGrowthFund:
			if res.Level < item.PreLimit {
				continue
			}
		case constant.ActivityTypeLottery:
			// 处理抽奖活动
			if int32(res.LotteryNum) < item.PreLimit {
				continue
			}
		case constant.ActivityTypeRefine:
			// 处理洗练
			if int32(res.RefineNum) < item.PreLimit {
				continue
			}
		case constant.ActivityTypeInstance:
			// 处理副本活动
			if int32(res.InstanceNum) < item.PreLimit {
				continue
			}
		}

		if roleReceiveInfo[int32(needDay)] > 0 {
			//领取了且当初也是vip，那就跳过
			if roleReceiveInfo[int32(needDay)] > 1 {
				continue
			}
			//领取了且当初不是vip，补发vip的
			if r.Buy > 1 {
				//校验背包情况
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddGroupUp)
				if err3 != nil {
					continue
				}

				_, err := db.ActivityAddDayUpdate(
					r.Uid,
					item.UUID,
					2,
					int(r.Buy),
				)
				if err != nil {
					continue
				}
			}
		} else {
			//没领取，那么发，但是要校验是不是vip
			if r.Buy > 1 {
				result2 := make([]protocol.SourceInfo, len(item.PayItem))

				for i, info := range item.PayItem {
					result2[i].Id = info.Id
					result2[i].Num = info.Num
				}
				err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemAddGroupUp)
				if err3 != nil {
					continue
				}
			}
			result3 := make([]protocol.SourceInfo, len(item.CommonItem))
			for i, info := range item.CommonItem {
				result3[i].Id = info.Id
				result3[i].Num = info.Num
			}
			err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddGroupUp)
			if err3 != nil {
				continue
			}
			_, err := db.ActivityAddDay(
				r.Uid,
				item.UUID,
				int(r.Buy),
				r.ActivityId,
			)
			if err != nil {
				continue
			}
		}

	}

	//领取奖励
	return true, ""
}

// 累充礼包列表
func activityListRechargeGrade(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityGradeResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	if request.ActivityId == 0 {
		return &protocol.ActivityGradeResponse{}, errutil.ErrInvalidParameter
	}
	res, err := ActivityListRechargeGrade(request)
	if err != nil {
		if err == errutil.ErrActivityExpire {
			return &protocol.ActivityGradeResponse{Code: -1, Data: errutil.ErrActivityExpire.Error()}, nil
		}
		return &protocol.ActivityGradeResponse{}, err
	}
	//比较累计充值金�?
	ress, err := game.UserInfo(Uid)
	if err != nil {
		return nil, err
	}
	return &protocol.ActivityGradeResponse{Data: "成功", ActivityListGrade: res, Num: ress.CumulativeChargeAmount}, nil
}

// ActivityListRechargeGrade 累充礼包列表
func ActivityListRechargeGrade(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListGrade, error) {
	actType := constant.ActivityTypeCumulativeCharge
	activityCfg := global.ActivityCfg

	//活动内容
	goodsId := make([]int64, 0)
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				if r.UUID > 0 {
					if r.UUID == activityn.UUID {
						goodsId = append(goodsId, activityn.UUID)
					}
				} else {
					goodsId = append(goodsId, activityn.UUID)
				}
			}
		}
	}

	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	getInfo := make(map[int64]bool)

	if len(resultb) > 0 {
		for _, item := range resultb {
			getInfo[item.GoodsUuiD] = true
		}
	}

	ActivityListRechargeAll := make([]protocol.ActivityListGrade, len(goodsId))
	for i, item := range goodsId {
		RechargeInfo, err := global.GetActivityCfg(int32(item))
		if err != nil {
			return nil, err
		}
		GetStatus := 1
		if getInfo[item] {
			GetStatus = 2
		}
		//比较累计充值金�?
		res, err := game.UserInfo(r.Uid)
		if err != nil {
			return nil, err
		}
		// PreLimit 在这里表示需要的累计充值金额（单位：元），CumulativeChargeAmount 单位是分
		if int64(RechargeInfo.PreLimit) > res.CumulativeChargeAmount {
			GetStatus = 3
		}
		ActivityListRechargeAll[i] = protocol.ActivityListGrade{
			Uuid:   int32(item),
			Grade:  RechargeInfo.PreLimit, // 这里 Grade 字段实际表示需要的充值金额（元）
			Status: GetStatus,
		}
	}
	return ActivityListRechargeAll, nil
}

// 领取累充礼包奖励
func activityReceiveRechargeGrade(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.UUID = r.UUID
	if request.UUID == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	_, err := ActivityReceiveRechargeGrade(request)
	if err != nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// ActivityReceiveRechargeGrade 领取累充礼包
func ActivityReceiveRechargeGrade(r *protocol.ActivityOperateRequest) (bool, error) {

	activityCfg := global.ActivityCfg
	result := protocol.ActivityDetailDay{}
	resetTime := int32(0)
	//活动内容
	for _, activityn := range activityCfg {
		if activityn.UUID == r.UUID {
			item := protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				PreLimit:   activityn.PreLimit,
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
			}
			for j, attachment2 := range activityn.CommonItem {
				// 假设 Attachment 结构体中有可以复制的字段
				item.CommonItem[j] = protocol.SourceInfo{
					// ... 根据字段一一复制�?
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}
			result = item
			resetTime = activityn.ResetTime
		}
	}
	//需要领取物品的累计充值金�? //比较累计充值金�?
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err
	}
	if resetTime > 0 && res.CreateDay > resetTime {
		return false, errutil.ErrActivityExpire
	}
	// PreLimit 表示需要的累计充值金额（单位：元），CumulativeChargeAmount 单位是分
	if int64(result.PreLimit) > res.CumulativeChargeAmount {
		return false, errutil.ErrNotEligibleForReward
	}
	goodsId := make([]int64, 0)
	goodsId = append(goodsId, result.UUID)
	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	if len(resultb) > 0 {
		return false, errutil.ErrGetAward
	}
	result3 := make([]protocol.SourceInfo, len(result.CommonItem))
	for i, info := range result.CommonItem {
		result3[i].Id = info.Id
		result3[i].Num = info.Num
	}
	//领取奖励
	err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddActivityLevel)
	if err3 != nil {
		return false, err3
	}
	_, err = db.GradeAdd(
		r.Uid,
		int32(result.UUID),
		result3,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 通用礼包列表（通过type区分）
func activityListGradeByType(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityGradeResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.Atype = r.Atype
	if request.ActivityId == 0 {
		return &protocol.ActivityGradeResponse{}, errutil.ErrInvalidParameter
	}
	if request.Atype == 0 {
		return &protocol.ActivityGradeResponse{}, errutil.ErrInvalidParameter
	}
	res, err := ActivityListGradeByType(request)
	if err != nil {
		if err == errutil.ErrActivityExpire {
			return &protocol.ActivityGradeResponse{Code: -1, Data: errutil.ErrActivityExpire.Error()}, nil
		}
		return &protocol.ActivityGradeResponse{}, err
	}
	// 获取用户信息用于返回相关数值
	ress, err := game.UserInfo(Uid)
	if err != nil {
		return nil, err
	}
	// 根据不同的活动类型返回不同的数值
	var num int64 = 0
	switch request.Atype {
	case constant.ActivityTypePinkEntry:
		// 粉色词条活动，返回粉色词条数量
		num = ress.PinkEntryNum
	case constant.ActivityTypeRedPartnerStar:
		// 红色伙伴总星级活动，返回红色伙伴总星级
		num = ress.RedPartnerStarNum
	case constant.ActivityTypeEquipmentStar:
		// 装备总升星活动，返回装备总升星数
		num = ress.EquipmentStarNum
	case constant.ActivityTypePetSynthesis:
		// 宠物合成活动，返回宠物合成数量
		num = ress.PetSynthesisNum
	case constant.ActivityTypeSoulStoneLevel:
		// 魂石总等级活动，返回魂石总等级
		num = ress.SoulStoneLevelNum
	case constant.ActivityTypeFiveElementLevel:
		// 五行赋能总等级活动，返回五行赋能总等级
		num = ress.FiveElementLevelNum
	}
	return &protocol.ActivityGradeResponse{Data: "成功", ActivityListGrade: res, Num: num}, nil
}

// ActivityListGradeByType 通用礼包列表（通过type区分）
func ActivityListGradeByType(r *protocol.ActivityOperateRequest) ([]protocol.ActivityListGrade, error) {
	actType := r.Atype
	activityCfg := global.ActivityCfg

	//活动内容
	goodsId := make([]int64, 0)
	for _, activityn := range activityCfg {
		if activityn.Acttag == r.ActivityId {
			if activityn.Type == actType {
				if r.UUID > 0 {
					if r.UUID == activityn.UUID {
						goodsId = append(goodsId, activityn.UUID)
					}
				} else {
					goodsId = append(goodsId, activityn.UUID)
				}
			}
		}
	}

	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	getInfo := make(map[int64]bool)

	if len(resultb) > 0 {
		for _, item := range resultb {
			getInfo[item.GoodsUuiD] = true
		}
	}

	// 获取用户信息用于比较数值
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, err
	}

	ActivityListAll := make([]protocol.ActivityListGrade, len(goodsId))
	for i, item := range goodsId {
		ActivityInfo, err := global.GetActivityCfg(int32(item))
		if err != nil {
			return nil, err
		}
		GetStatus := 1
		if getInfo[item] {
			GetStatus = 2
		}
		// 根据不同的活动类型比较不同的数值
		var currentValue int64 = 0
		switch actType {
		case constant.ActivityTypePinkEntry:
			// 粉色词条活动，比较粉色词条数量
			currentValue = res.PinkEntryNum
		case constant.ActivityTypeRedPartnerStar:
			// 红色伙伴总星级活动，比较红色伙伴总星级
			currentValue = res.RedPartnerStarNum
		case constant.ActivityTypeEquipmentStar:
			// 装备总升星活动，比较装备总升星数
			currentValue = res.EquipmentStarNum
		case constant.ActivityTypePetSynthesis:
			// 宠物合成活动，比较宠物合成数量
			currentValue = res.PetSynthesisNum
		case constant.ActivityTypeSoulStoneLevel:
			// 魂石总等级活动，比较魂石总等级
			currentValue = res.SoulStoneLevelNum
		case constant.ActivityTypeFiveElementLevel:
			// 五行赋能总等级活动，比较五行赋能总等级
			currentValue = res.FiveElementLevelNum
		}
		// PreLimit 表示需要的数值
		if int64(ActivityInfo.PreLimit) > currentValue {
			GetStatus = 3
		}
		ActivityListAll[i] = protocol.ActivityListGrade{
			Uuid:   int32(item),
			Grade:  ActivityInfo.PreLimit,
			Status: GetStatus,
		}
	}
	return ActivityListAll, nil
}

// 通用礼包领取（通过type区分）
func activityReceiveGradeByType(ctx context.Context, r *protocol.ActivityOperateRequest) (*protocol.ActivityReceiveDayResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.ActivityOperateRequest{}
	request.Uid = Uid
	request.ActivityId = r.ActivityId
	request.UUID = r.UUID
	request.Atype = r.Atype
	if request.UUID == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	if request.ActivityId == 0 {
		return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
	}
	if request.Atype == 0 {
		if cfg, ok := global.ActivityCfg[int32(request.UUID)]; ok {
			request.Atype = cfg.Type
		} else {
			return &protocol.ActivityReceiveDayResponse{}, errutil.ErrInvalidParameter
		}
	}
	_, err := ActivityReceiveGradeByType(request)
	if err != nil {
		return &protocol.ActivityReceiveDayResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.ActivityReceiveDayResponse{Data: "成功"}, nil
}

// ActivityReceiveGradeByType 通用礼包领取（通过type区分）
func ActivityReceiveGradeByType(r *protocol.ActivityOperateRequest) (bool, error) {
	actType := r.Atype
	activityCfg := global.ActivityCfg
	result := protocol.ActivityDetailDay{}
	startTime := int32(0)
	resetTime := int32(0)
	//活动内容
	for _, activityn := range activityCfg {
		if activityn.UUID == r.UUID && activityn.Type == actType {
			item := protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				PreLimit:   activityn.PreLimit,
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
			}
			for j, attachment2 := range activityn.CommonItem {
				item.CommonItem[j] = protocol.SourceInfo{
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}
			result = item
			startTime = activityn.StartTime
			resetTime = activityn.ResetTime
			break
		}
	}
	if result.UUID == 0 {
		return false, errutil.ErrConfigNotFound
	}
	//比较前置条件
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err
	}
	serverOpenDay := global.ServerListMgrInstance.GetServerOpenDay(res.ServerId)
	if startTime > 0 && int32(serverOpenDay) < startTime && int32(serverOpenDay) > resetTime+startTime {
		return false, errutil.ErrActivityExpire
	}
	// 根据不同的活动类型比较不同的数值
	var currentValue int64 = 0
	switch actType {
	case constant.ActivityTypePinkEntry:
		// 粉色词条活动，比较粉色词条数量
		currentValue = res.PinkEntryNum
	case constant.ActivityTypeRedPartnerStar:
		// 红色伙伴总星级活动，比较红色伙伴总星级
		currentValue = res.RedPartnerStarNum
	case constant.ActivityTypeEquipmentStar:
		// 装备总升星活动，比较装备总升星数
		currentValue = res.EquipmentStarNum
	case constant.ActivityTypePetSynthesis:
		// 宠物合成活动，比较宠物合成数量
		currentValue = res.PetSynthesisNum
	case constant.ActivityTypeSoulStoneLevel:
		// 魂石总等级活动，比较魂石总等级
		currentValue = res.SoulStoneLevelNum
	case constant.ActivityTypeFiveElementLevel:
		// 五行赋能总等级活动，比较五行赋能总等级
		currentValue = res.FiveElementLevelNum
	}
	// PreLimit 表示需要的数值
	if int64(result.PreLimit) > currentValue {
		return false, errutil.ErrNotEligibleForReward
	}
	goodsId := make([]int64, 0)
	goodsId = append(goodsId, result.UUID)
	//查询用户领取情况
	resultb, _ := db.GradeList(
		r.Uid,
		int32(365*50),
		goodsId,
	)
	if len(resultb) > 0 {
		return false, errutil.ErrGetAward
	}
	result3 := make([]protocol.SourceInfo, len(result.CommonItem))
	for i, info := range result.CommonItem {
		result3[i].Id = info.Id
		result3[i].Num = info.Num
	}
	//领取奖励
	err3 := CheckSourceAdd(result3, r.Uid, true, constant.ItemAddActivityLevel)
	if err3 != nil {
		return false, err3
	}
	_, err = db.GradeAdd(
		r.Uid,
		int32(result.UUID),
		result3,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
