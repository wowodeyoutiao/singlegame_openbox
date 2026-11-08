package api

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"nanoserver/db"
	"nanoserver/dispatcher"
	"nanoserver/global"
	"nanoserver/internal/game"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/timex"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type PayUseCase struct{}

func NewPayUseCase() *PayUseCase {
	return &PayUseCase{}
}

func (p *PayUseCase) PayNotify(req *protocol.PayNotifyRequest) (*protocol.PayNotifyResponse, error) {
	var (
		totalCharacterSpend = req.TotalCharacterSpend / 100 // 转换为元
	)

	createTime := time.Now().Unix()

	// 获取配置
	storeCfg, ok := global.StoreCfg[int32(req.StoreId)]
	if !ok {
		return nil, errors.Errorf("store configuration not found for ID: %d", req.StoreId)
	}
	multiplier := int64(1)
	//首充双倍，特别关照
	if storeCfg.StoreType == constant.ShopTypeRecharge {
		list, err := StoreList(&protocol.StoreListRequest{
			GoodsUuid: int32(req.StoreId),
			Num:       req.Num,
			Uid:       int64(req.Uid),
			StoreType: storeCfg.StoreType,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get store list for store ID: %d", req.StoreId)
		}

		if len(list) > 0 && list[0].Buy == 0 {
			multiplier++
		}
	}

	//购买多份，给到的数量就x 倍
	if req.Num > 1 {
		multiplier = int64(req.Num) * multiplier
	} else {
		req.Num = 1
	}
	sourceAdd := make([]protocol.SourceInfo, 0, len(storeCfg.Props))
	for _, item := range storeCfg.Props {
		if item.Id == constant.ItemTypeLuckyBagOne || item.Id == constant.ItemTypeLuckyBagTwo || item.Id == constant.ItemTypeLuckyBagThree {
			game.TaskDispatch(int64(req.Uid), pb.TargetType_LUCKY_BAG_OBTAINED_TOTAL, 0, 1)
		}
		sourceAdd = append(sourceAdd, protocol.SourceInfo{
			Id:       item.Id,
			Num:      int64(item.Num) * multiplier,
			ShopType: storeCfg.StoreType,
			Method:   storeCfg.Method,
		})
	}

	// 获取用户信息
	userInfo, err := game.UserInfo(int64(req.Uid))
	if err != nil {
		return nil, err
	}

	// 添加道具
	if storeCfg.StoreType != constant.ShopTypeSubscription {
		err := CheckSourceAdd(sourceAdd, int64(req.Uid), true, storeCfg.Method+"-"+storeCfg.Entrance)
		if err != nil {
			if err == errutil.ErrPlayerNotFound {
				return nil, err
			}
			if err == errutil.ErrBagFull {
				MailId := strconv.FormatInt(snowflake.GenID(), 10)
				result3 := make([]protocol.Attachment, len(sourceAdd))
				for i, Attachment := range sourceAdd {
					result3[i] = protocol.Attachment{
						MailId:      MailId,
						UserPropsId: int64(Attachment.Id),
						Num:         Attachment.Num,
					}
				}
				mailAddRequest := &protocol.MailAddRequest{
					Uid:          int64(req.Uid),
					RecipientUid: -1,
					Title:        "系统邮件",
					Desc:         "包裹已满，请及时清理",
					Attachment:   result3,
					ChannelCode:  "0",
				}
				mailAddResult, err := Mailadd(mailAddRequest)
				if !mailAddResult || err != nil {
					return nil, err
				}
			}
		}
	}

	bizHandler, exist := PayCallbackHandlesBefore[storeCfg.StoreType]
	if exist {
		if err := bizHandler(storeCfg, userInfo); err != nil {
			return nil, errors.Wrapf(err, "failed to handle payment callback for store type: %d", storeCfg.StoreType)
		}
	}

	log.Infof("Payment notification processed successfully for user: %d, store type: %d", req.Uid, storeCfg.StoreType)

	jsonData, err := json.Marshal(sourceAdd)
	if err != nil {
		return nil, err
	}
	jsonDataStr := string(jsonData)
	storeAdditionResult, err := db.StoreAdd(
		int64(req.Uid),
		storeCfg.Id,
		jsonDataStr,
		req.Num,
		createTime,
	)
	if err != nil {
		logger.Errorf("db.StoreAdd failed for user: %d, store ID: %d, err: %v", req.Uid, storeCfg.Id, err)
		return nil, errors.Wrap(err, "db.StoreAdd failed")
	}
	if !storeAdditionResult {
		return nil, nil
	}
	bizHandlerAfter, exist := PayCallbackHandlesAfter[storeCfg.StoreType]
	if exist {
		if err := bizHandlerAfter(storeCfg, userInfo); err != nil {
			return nil, errors.Wrapf(err, "failed to handle payment callback for store type: %d", storeCfg.StoreType)
		}
	}

	//数数上报充值
	getItemInfo := make([]map[string]interface{}, 0)
	uid := strconv.FormatInt(int64(req.Uid), 10)
	for _, v := range sourceAdd {
		if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
			getItemInfo = append(getItemInfo, map[string]interface{}{
				"item_id":   ItemCfg.ID,
				"item_name": ItemCfg.Name,
				"item_num":  math.Abs(float64(v.Num)),
			})
		}
	}
	if storeCfg.Limit.NeedMoney == int32(constant.NeedMoney) {
		if err := thinkingdatax.TrackerIns.E(uid, userInfo.AccountName, "order_finish", tracking.M{
			"server_id":             strconv.Itoa(int(userInfo.ServerId)), // 服务器ID，格式为字符串
			"currency_type":         "CNY",                                // 币种，文本
			"is_first_pay":          req.IsFirstRecharge,                  // 是否首次充值，布尔
			"order_id":              req.OutTradeNo,                       // 订单id，文本
			"pay_amount":            storeCfg.Limit.Num,                   // 支付金额，数值
			"payment_name":          storeCfg.Method,                      // 购买项目名，文本
			"payment_entrance":      storeCfg.Entrance,                    // 付费触发位置，文本
			"get_item_info":         getItemInfo,                          // 获得物品详情，对象组
			"total_character_spend": totalCharacterSpend,                  // 本次充值后角色的总消费金额，数值, 单位是元
			"is_success":            true,                                 // 是否成功，布尔
		}); err != nil {
			log.Warnf("thinkingdatax order_finish event error: %v", err)
		}
		// upStoreOrderFinishShushu(storeCfg.Entrance, storeCfg.Method, req.OutTradeNo, req.TransactionId, int(storeCfg.Limit.Num), int64(req.Uid), sourceAdd, req.IsFirstRecharge)
		if req.IsFirstRecharge {
			if err := thinkingdatax.TrackerIns.SetOnce(uid, userInfo.AccountName, tracking.M{
				"first_pay_time": time.Now().Format(timex.FormatDateTimeMillis),
			}); err != nil {
				log.Warnf("thinkingdatax set first_pay_time error: %v", err)
			}
		}
		if err := thinkingdatax.TrackerIns.Set(uid, userInfo.AccountName, tracking.M{
			"total_character_spend": totalCharacterSpend,
		}); err != nil {
			log.Warnf("thinkingdatax set first_pay_time error: %v", err)
		}
		game.AddUpRecord(int64(req.Uid), constant.CumulativeChargeAmount, int64(storeCfg.Limit.Num))
	}

	return &protocol.PayNotifyResponse{
		Code:    0,
		Message: "Payment notification processed successfully",
	}, nil
}

func upStoreOrderFinishShushu(buy_name string, buy_entrance string, order_id string, transaction_id string, pay_amount int, uid int64, sourceInfo []protocol.SourceInfo, is_first_pay bool) {
	data := make(map[string]interface{}, 0)
	data["name"] = "order_finish"
	data["order_id"] = order_id
	data["transaction_id"] = transaction_id
	data["pay_amount"] = pay_amount
	data["payment_entrance"] = buy_name
	data["payment_name"] = buy_entrance
	data["currency_type"] = "CNY"
	data["is_success"] = true
	data["is_first_pay"] = is_first_pay
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		// 查询itemname
		if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
			get_item_info = append(get_item_info, map[string]interface{}{
				"item_id":   ItemCfg.ID,
				"item_name": ItemCfg.Name,
				"item_num":  math.Abs(float64(v.Num)),
			})
		}
	}
	data["get_item_info"] = get_item_info

	game.Upshushu(uid, data, true, false)
}
func (p *PayUseCase) PayUserInfo(req *protocol.PayUserInfoRequest) (*protocol.PayUserInfoResponse, error) {
	if req.Uid != 0 {
		player, err := game.UserInfo(req.Uid)
		if err != nil {
			return nil, err
		}
		return &protocol.PayUserInfoResponse{
			Uid:      req.Uid,
			Name:     player.Name,
			ServerId: int64(player.ServerId),
			Level:    int(player.Level),
			VipLevel: int(player.Vip),
		}, nil
	}

	if req.OpenID != "" {
		player, err := db.QueryPlayersByOpenID(req.OpenID)
		if err != nil && !errors.Is(err, errutil.ErrUserNotFound) {
			return nil, err
		} else if errors.Is(err, errutil.ErrUserNotFound) {
			return nil, errors.Wrap(errutil.ErrUserNotFound, fmt.Sprintf("OpenID: %s", req.OpenID))
		}

		return &protocol.PayUserInfoResponse{
			Uid:      player.Uid,
			Name:     player.Name,
			ServerId: int64(player.ServerId),
			Level:    int(player.Level),
			VipLevel: int(player.Vip),
		}, nil
	}

	return nil, errors.New("either uid or openid must be provided")
}

// 空实现
func HandleShopCallbackNoop(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	return nil
}

var PayCallbackHandlesBefore = map[int32]func(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error{
	// constant.ShopTypeRecharge:              HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeCoupon:                HandleShopCallbackNoop,
	// constant.ShopTypeGold:                  HandleShopCallbackNoop,
	// constant.ShopTypeSoul:                  HandleShopCallbackNoop,
	// constant.ShopTypeTrial:                 HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeSignin30Days:          HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeActivityRecommendGift: HandleShopCallbackNoop,
	constant.ShopTypeFirstCharge:  HandleShopTypeFirstCharge,  // 有充值
	constant.ShopTypeSubscription: HandleShopTypeSubscription, // 有充值
	// constant.ShopTypeBoost:                 HandleShopCallbackNoop,
	// constant.ShopTypeVIPGift:               HandleShopCallbackNoop,
	// constant.ShopTypeOfferReward:           HandleShopCallbackNoop,
	// constant.ShopTypePVP:                   HandleShopCallbackNoop,
	// constant.ShopTypeRelic:                 HandleShopCallbackNoop,
	constant.ShopTypeTrialReward: HandleShopTypeTrialReward, // 有充值
	// constant.ShopTypeGuild:                 HandleShopCallbackNoop,
	// constant.ShopTypePartner:               HandleShopCallbackNoop, // 有充值
	// constant.ShopTypePet:                   HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeAFK:                   HandleShopCallbackNoop,
	// constant.ShopTypeDailyGift:             HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeSoulStoneGift:         HandleShopCallbackNoop, // 有充值
	// constant.ShopTypeGrowthFund:            HandleShopCallbackNoop,
	// constant.ShopTypeLimitedTimeGift:       HandleShopCallbackNoop,
	// constant.ShopTypeElementGift:           HandleShopCallbackNoop,     // 有充值
	constant.ShopTypeGuardIsland: HandleShopTypeGuardIsland, // 有充值
	// constant.ShopTypeAd:                    HandleShopCallbackNoop,
	// constant.ShopTypeElementChallenge:      HandleShopCallbackNoop,
	// constant.ShopTypeExchange: HandleShopTypeExchange,
	// constant.ShopTypeEnergy:                HandleShopCallbackNoop,
	constant.ShopTypeRefine: HandleShopRefine, // 洗炼
}
var PayCallbackHandlesAfter = map[int32]func(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error{
	constant.ShopTypeSevenDaysFun: HandleShopTypeSevenDaysFun, // 有充值
}

// func HandleShopTypeExchange(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
// 	game.TaskDispatch(userInfo.UID, pb.TargetType_TICKET_EXCHANGE_TIMES, 0, 1)
// 	return nil
// }

func HandleShopTypeTrialReward(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	strSend := []string{userInfo.Name}
	game.Lamp(int32(constant.LampTrialFundPurchased), userInfo.ServerId, strSend)
	return nil
}

func HandleShopTypeSevenDaysFun(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	game.TaskDispatch(userInfo.UID, pb.TargetType_SEVENDAY_TASK_GIFT_PURCHASED, 0, 1)
	return nil
}

func HandleShopTypeGuardIsland(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	var (
		guardIslandRecrodID int64
		found               bool
	)

	for _, v := range global.GuardIslandCfg {
		if v.ID == storeCfg.Id {
			guardIslandRecrodID = int64(v.ID)
			found = true
			break
		}
	}
	if !found {
		return errors.Errorf("guard island record not found for uuid: %d", storeCfg.Id)
	}

	if err := game.GuardIslandMultiplier(userInfo.UID, guardIslandRecrodID); err != nil {
		return errors.Wrapf(err, "failed to set guard island multiplier for user: %d", userInfo.UID)
	}

	return nil
}

func HandleShopRefine(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	dispatcher.MustGetInstance().Publish(dispatcher.TopicActivatePurchasedSlot.String(), &dispatcher.ActivatePurchasedSlot{
		UID:           userInfo.UID,
		PurchasedSlot: int32(storeCfg.GoodsId),
	})
	return nil
}

func HandleShopTypeSubscription(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	sourceAdd := make([]protocol.SourceInfo, 0, len(storeCfg.Props))
	for _, item := range storeCfg.Props {
		sourceAdd = append(sourceAdd, protocol.SourceInfo{
			Id:       item.Id,
			Num:      int64(item.Num),
			ShopType: storeCfg.StoreType,
			Method:   storeCfg.Method,
		})
	}

	mailId := strconv.FormatInt(snowflake.GenID(), 10)
	result3 := make([]protocol.Attachment, len(sourceAdd))
	for i, Attachment := range sourceAdd {
		result3[i] = protocol.Attachment{
			MailId:      mailId,
			UserPropsId: int64(Attachment.Id),
			Num:         Attachment.Num,
		}
	}
	mailAddRequest := &protocol.MailAddRequest{
		Uid:          -1,
		RecipientUid: userInfo.UID,
		Title:        "系统邮件：" + storeCfg.Entrance + "奖励",
		Desc:         storeCfg.Entrance + "购买立返奖励",
		Attachment:   result3,
		ChannelCode:  "0",
	}
	mailAddResult, err := Mailadd(mailAddRequest)
	if !mailAddResult {
		return errors.Wrap(err, "failed to add mail for subscription reward")
	}

	var pName string
	for _, privilegen := range global.PrivilegeCfg {
		if privilegen.Sid == storeCfg.Id {
			pName = privilegen.Pname
			originalTime := time.Now()
			nextDayMidnight := originalTime.AddDate(0, 0, int(privilegen.Duration)) // 添加N天
			nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
			err := game.PrivilegeAdd(storeCfg.StoreType, userInfo.UID, privilegen.Pid, nextDayMidnight)
			if err != nil {
				return err
			}
		}
	}

	strSend := []string{userInfo.Name, pName}
	game.Lamp(int32(constant.LampPrivilegeCardPurchased), userInfo.ServerId, strSend)

	return nil
}

func HandleShopTypeFirstCharge(storeCfg *global.StoreConfig, userInfo *protocol.UserInfo) error {
	strSend := []string{userInfo.Name}
	game.Lamp(int32(constant.LampFirstRecharge), userInfo.ServerId, strSend)
	return nil
}
