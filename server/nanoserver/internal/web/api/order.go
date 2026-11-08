package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nanoserver/global"
	"nanoserver/internal/game"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

func MakeOrderService() http.Handler {
	router := mux.NewRouter()
	// router.Handle("/v1/order/console/", nex.Handler(orderList)).Methods("POST") // 订单列表
	// router.Handle("/v1/order/", nex.Handler(createOrder)).Methods("POST")       // 创建订单（未使用）
	// router.Handle("/v1/order/create/wechat", nex.Handler(createOrderWeChat)).Methods("POST")    // 创建微信订单（未使用）
	// router.Handle("/v1/order/notify/wechat", nex.Handler(wechatCallback)).Methods("POST")       // 微信订单回调
	// router.Handle("/v1/order/notify/alipay", nex.Handler(alipayCallback)).Methods("POST")       // 支付宝订单回调
	// router.Handle("/v1/order/notify/byrty0uxi", nex.Handler(byrty0uxiCallback)).Methods("POST") // byrty0uxi订单回调

	// 微信 Android
	// router.HandleFunc("/v1/order/wx/android/callbacks/paymentns", handleWechatMessagePushVerifyURL).Methods("GET")     // 微信订单回调 GET 请求验签（模拟推送）
	// router.Handle("/v1/order/wx/android/callbacks/paymentns", nex.Handler(handleMessagePushPlainText)).Methods("POST") // 微信订单回调
	// 微信 IOS
	// router.HandleFunc("/v1/order/wx/ios/callbacks/customer/message", processCustomerMessageSignature).Methods("GET")  // WX SDK 的客服息息回调（模拟推送）
	// router.Handle("/v1/order/wx/ios/callbacks/customer/message", nex.Handler(processCustomerMessage)).Methods("POST") // WX SDK 的客服息息回调请求
	// router.HandleFunc("/v1/order/wx/ios/html", generateWechatRechargeHTML).Methods("GET")                             // 微信公众号充值页面（兼容客户端URL格式）
	// router.Handle("/v1/order/wx/ios/callbacks/payments", nex.Handler(processIOSPaymentCallback)).Methods("POST")      // WX SDK 的支付结果通知

	return router
}

// func CreateOrder(r *protocol.CreateOrderRequest) (interface{}, error) {
// 	storeCfg, err := global.GetStoreCfg(int32(r.ProductId))
// 	if err != nil {
// 		return nil, err
// 	}

// 	order := &model.Order{
// 		OrderId:      strings.Replace(uuid.New(), "-", "", -1),
// 		AppId:        r.AppID,
// 		Uid:          r.Uid,
// 		ChannelId:    r.ChannelID,
// 		PayPlatform:  r.Platform,
// 		Extra:        r.Extra,
// 		ProductName:  r.ProductionName,
// 		ProductCount: r.ProductCount,
// 		CreatedAt:    time.Now().Unix(),
// 		Status:       db.OrderStatusCreated,
// 		Remote:       r.Device.Remote,
// 		Ip:           r.Device.Remote,
// 		Imei:         r.Device.IMEI,
// 		Model:        r.Device.Model,
// 		Os:           r.Device.OS,
// 		ProductId:    r.ProductId,
// 		Money:        int(storeCfg.Limit.Num),
// 	}

// 	resp, err := provider.Byrtypay.CreateOrderResponse(order)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		return nil, err
// 	}

// 	if err := db.InsertOrder(order); err != nil {
// 		logger.Error(err.Error())
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func createOrder(ctx context.Context, r *http.Request) (interface{}, error) {
// 	uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

// 	if err := r.ParseForm(); err != nil {
// 		return nil, err
// 	}
// 	count, err := strconv.Atoi(strings.TrimSpace(r.Form.Get("count")))
// 	if err != nil {
// 		return nil, err
// 	}
// 	productId, err := strconv.Atoi(strings.TrimSpace(r.Form.Get("productId")))
// 	if err != nil {
// 		return nil, err
// 	}

// 	appId := strings.TrimSpace(r.Form.Get("appId"))
// 	channelId := strings.TrimSpace(r.Form.Get("channelId"))
// 	platform := strings.TrimSpace(r.Form.Get("platform"))
// 	extra := strings.TrimSpace(r.Form.Get("extra"))
// 	name := strings.TrimSpace(r.Form.Get("name"))
// 	if appId == "" || channelId == "" || count == 0 || name == "" {
// 		return nil, errutil.ErrIllegalParameter
// 	}

// 	request := &protocol.CreateOrderRequest{
// 		AppID:          appId,
// 		ChannelID:      channelId,
// 		Platform:       platform,
// 		ProductId:      productId,
// 		ProductionName: name,
// 		ProductCount:   count,
// 		Extra:          extra,
// 		Uid:            uid,
// 		Device:         protocol.Device{Remote: r.RemoteAddr},
// 	}

// 	return CreateOrder(request)
// }

func CreateOrderV2(r *protocol.CreateOrderRequest) (resp *protocol.StoreOperateResponse, err error) {
	storeCfg, err := global.GetStoreCfg(int32(r.ProductId))
	if err != nil {
		return nil, err
	}

	orderID := strings.Replace(uuid.New(), "-", "", -1)
	money := int(storeCfg.Limit.Num)
	uid := r.Uid
	extra := r.Extra
	// order := &model.Order{
	// 	OrderId:      orderID,
	// 	AppId:        r.AppID,
	// 	Uid:          r.Uid,
	// 	ChannelId:    r.ChannelID,
	// 	PayPlatform:  r.Platform,
	// 	Extra:        r.Extra,
	// 	ProductName:  r.ProductionName,
	// 	ProductCount: r.ProductCount,
	// 	CreatedAt:    time.Now().Unix(),
	// 	Status:       db.OrderStatusCreated,
	// 	Remote:       r.Device.Remote,
	// 	Ip:           r.Device.Remote,
	// 	Imei:         r.Device.IMEI,
	// 	Model:        r.Device.Model,
	// 	Os:           r.Device.OS,
	// 	ProductId:    r.ProductId,
	// 	Money:        int(storeCfg.Limit.Num),
	// 	OpenID:       r.OpenID,
	// }

	// 微信Android支付
	switch r.Platform {
	// case constant.PayPlatformTypeWeChatAndroid.String():
	// 	logger.Infof("CreateOrderV2: 创建微信Android订单: %+v", r)

	// 	userInfo, err := game.UserInfo(order.Uid)
	// 	if err != nil {
	// 		logger.Error("CreateOrderV2: 虚拟支付获取用户信息失败: ", err)
	// 		return nil, errutil.ErrIllegalParameter
	// 	}

	// 	signData := &protocol.SignData{}
	// 	if err = json.Unmarshal([]byte(order.Extra), signData); err != nil {
	// 		logger.Error("CreateOrderV2: 虚拟支付解析订单额外信息失败: ", err)
	// 		return nil, errutil.ErrIllegalParameter
	// 	}

	// 	attach := &protocol.Attach{
	// 		StoreId: int64(storeCfg.Id),
	// 	}
	// 	attachStr, err := json.Marshal(attach)
	// 	if err != nil {
	// 		logger.Error("CreateOrderV2: 虚拟支付解析订单额外信息失败: ", err)
	// 		return nil, errutil.ErrIllegalParameter
	// 	}

	// 	midasPaymentGameItemRequest := &protocol.MidasPaymentGameItemRequest{
	// 		OpenId: userInfo.AccountName,
	// 		UserId: uint64(r.Uid),
	// 		SignData: &protocol.SignData{
	// 			Mode:         signData.Mode,
	// 			OfferId:      signData.OfferId,
	// 			BuyQuantity:  signData.BuyQuantity,
	// 			Env:          signData.Env,
	// 			CurrencyType: signData.CurrencyType,
	// 			Platform:     signData.Platform,
	// 			ZoneId:       signData.ZoneId,
	// 			ProductId:    fmt.Sprintf("%d", storeCfg.Limit.Num*100), // ProductId 统一改成金额单位为分
	// 			GoodsPrice:   signData.GoodsPrice,
	// 			OutTradeNo:   orderID,
	// 			Attach:       string(attachStr),
	// 		},
	// 		ServerId:  r.ServerID, // 服务器 ID
	// 		RequestId: r.RequestID,
	// 	}

	// 	midasPaymentGameItemResponse, err := NewPaymentClient().MidasPaymentGameItem(context.Background(), midasPaymentGameItemRequest)
	// 	if err != nil {
	// 		logger.Error("CreateOrderV2: 虚拟支付生成微信支付请求失败: ", err)
	// 		return nil, errutil.ErrCreatePreOrderFailed
	// 	}
	// 	resp = &protocol.StoreOperateResponse{
	// 		Code:      0,
	// 		SignData:  midasPaymentGameItemResponse.SignData,
	// 		PaySig:    midasPaymentGameItemResponse.PaySig,
	// 		Signature: midasPaymentGameItemResponse.Signature,
	// 		Platform:  r.Platform,
	// 	}

	case constant.PayPlatformTypeWeChatIOS.String():
		logger.Infof("CreateOrderV2: 创建微信IOS订单: %+v", r)
		attach := &protocol.Attach{
			StoreId: int64(storeCfg.Id),
		}
		attachStr, err := json.Marshal(attach)
		if err != nil {
			logger.Error("CreateOrderV2: 虚拟支付解析订单额外信息失败: ", err)
			return nil, errutil.ErrIllegalParameter
		}

		// PrepayRequest 信息是不全的需要在pay补充
		prepayRequest := &protocol.PrepayRequest{
			Description: r.ProductionName,
			OutTradeNo:  orderID,
			TimeExpire:  time.Now().Add(30 * time.Minute).Format(time.RFC3339), // 订单失效时间，格式为rfc3339格式
			Amount: &protocol.PrepayRequestAmount{
				Total:    int64(storeCfg.Limit.Num * 100), // 金额单位为
				Currency: "CNY",                           // 人民币
			},
			Attach: string(attachStr), // 透传信息
			Payer: &protocol.PrepayRequestPayer{
				Openid: r.OpenID, // 用户的OpenID
			},
			Detail: &protocol.PrepayRequestDetail{
				CostPrice: int64(money * 100), // 商品原价，单位为分
				GoodsDetail: []*protocol.GoodsDetail{
					{
						MerchantGoodsId: fmt.Sprintf("%d", storeCfg.Id),
						Quantity:        1,
						UnitPrice:       int64(storeCfg.Limit.Num * 100),
					},
				},
			},
			UserId:    uint64(r.Uid), // 用户ID
			ServerId:  r.ServerID,    // 服务器ID
			RequestId: r.RequestID,   // 请求ID
		}
		prepayResponse, err := NewPaymentClient().JSAPIPrepay(context.Background(), prepayRequest)
		if err != nil {
			logger.Error("CreateOrderV2: 商户支付生成微信预支付订单失败: ", err)
			return nil, errutil.ErrCreatePreOrderFailed
		}
		logger.Debugf("CreateOrderV2: 微信预支付订单生成成功: %+v", prepayResponse)
		resp = &protocol.StoreOperateResponse{
			Code:             0,
			PrePayInternalID: prepayRequest.OutTradeNo,
			Platform:         r.Platform,
		}
	default:
		// logger.Errorf("CreateOrderV2: 不支持的支付平台: %s", r.Platform)
		// return nil, errutil.ErrIllegalParameter
		logger.Infof("CreateOrderV2: 创建微信Android订单: %+v, platform: %v", r, r.Platform)

		userInfo, err := game.UserInfo(uid)
		if err != nil {
			logger.Error("CreateOrderV2: 虚拟支付获取用户信息失败: ", err)
			return nil, errutil.ErrIllegalParameter
		}

		signData := &protocol.SignData{}
		if err = json.Unmarshal([]byte(extra), signData); err != nil {
			logger.Error("CreateOrderV2: 虚拟支付解析订单额外信息失败: ", err)
			return nil, errutil.ErrIllegalParameter
		}
		signData.GoodsPrice = int32(storeCfg.Limit.Num * 100) // 商品价格，单位为分

		attach := &protocol.Attach{
			StoreId: int64(storeCfg.Id),
		}
		attachStr, err := json.Marshal(attach)
		if err != nil {
			logger.Error("CreateOrderV2: 虚拟支付解析订单额外信息失败: ", err)
			return nil, errutil.ErrIllegalParameter
		}

		midasPaymentGameItemRequest := &protocol.MidasPaymentGameItemRequest{
			OpenId: userInfo.AccountName,
			UserId: uint64(r.Uid),
			SignData: &protocol.SignData{
				Mode:         signData.Mode,
				OfferId:      signData.OfferId,
				BuyQuantity:  signData.BuyQuantity,
				Env:          signData.Env,
				CurrencyType: signData.CurrencyType,
				Platform:     signData.Platform,
				ZoneId:       signData.ZoneId,
				ProductId:    fmt.Sprintf("%d", storeCfg.Limit.Num*100), // ProductId 统一改成金额单位为分
				GoodsPrice:   signData.GoodsPrice,
				OutTradeNo:   orderID,
				Attach:       string(attachStr),
			},
			ServerId:  r.ServerID, // 服务器 ID
			RequestId: r.RequestID,
		}

		midasPaymentGameItemResponse, err := NewPaymentClient().MidasPaymentGameItem(context.Background(), midasPaymentGameItemRequest)
		if err != nil {
			logger.Error("CreateOrderV2: 虚拟支付生成微信支付请求失败: ", err)
			return nil, errutil.ErrCreatePreOrderFailed
		}
		resp = &protocol.StoreOperateResponse{
			Code:             0,
			SignData:         midasPaymentGameItemResponse.SignData,
			PaySig:           midasPaymentGameItemResponse.PaySig,
			Signature:        midasPaymentGameItemResponse.Signature,
			PrePayInternalID: orderID,
			Platform:         r.Platform,
		}
	}

	// if err := db.InsertOrder(order); err != nil {
	// 	logger.Error(err.Error())
	// 	return nil, errutil.ErrDBOperation
	// }

	return resp, nil
}

// func TradeList(r *protocol.TradeListRequest) ([]protocol.TradeInfo, int, error) {
// 	if r == nil {
// 		return nil, 0, errutil.ErrIllegalParameter
// 	}

// 	list, total, err := db.TradeList(
// 		r.AppID,
// 		r.ChannelID,
// 		r.OrderID,
// 		r.Start,
// 		r.End,
// 		r.Offset,
// 		r.Count)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	result := make([]protocol.TradeInfo, len(list))
// 	for i, order := range list {
// 		result[i] = protocol.TradeInfo{
// 			OrderId:        order.OrderId,
// 			Uid:            strconv.FormatInt(order.Uid, 10),
// 			Money:          order.Money,
// 			RealMoney:      order.RealMoney,
// 			ProductName:    order.ProductName,
// 			ProductCount:   order.ProductCount,
// 			ServerName:     order.ServerName,
// 			RoleId:         order.RoleId,
// 			PayBy:          order.PayPlatform,
// 			AppId:          order.AppId,
// 			ChannelId:      order.ChannelId,
// 			PayPlatformUid: order.ComsumerId,
// 			PayAt:          order.PayAt,
// 			Currency:       order.Currency,
// 		}
// 	}
// 	return result, total, nil
// }

// func OrderList(r *protocol.OrderListRequest) ([]protocol.OrderInfo, int, error) {
// 	if r == nil {
// 		return nil, 0, errutil.ErrIllegalParameter
// 	}

// 	id, err := strconv.ParseInt(r.Uid, 10, 0)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	list, total, err := db.OrderList(
// 		id,
// 		r.AppID,
// 		r.ChannelID,
// 		r.OrderID,
// 		r.PayBy,
// 		r.Start,
// 		r.End,
// 		int(r.Status),
// 		r.Offset,
// 		r.Count)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	result := make([]protocol.OrderInfo, len(list))
// 	for i, order := range list {
// 		result[i] = protocol.OrderInfo{
// 			OrderId:      order.OrderId,
// 			Uid:          strconv.FormatInt(order.Uid, 10),
// 			Money:        order.Money,
// 			RealMoney:    order.RealMoney,
// 			ProductName:  order.ProductName,
// 			ProductCount: order.ProductCount,
// 			ServerName:   order.ServerName,
// 			RoleID:       order.RoleId,
// 			PayBy:        order.PayPlatform,
// 			AppId:        order.AppId,
// 			Imei:         order.Imei,
// 			Status:       order.Status,
// 			Extra:        order.Extra,
// 			CreatedAt:    order.CreatedAt,
// 		}
// 	}
// 	return result, total, nil
// }

// 订单列表
// func orderList(r *http.Request, form *nex.Form) (*protocol.OrderListResponse, error) {
// 	if !whitelist.VerifyIP(r.RemoteAddr) {
// 		return nil, errutil.ErrPermissionDenied
// 	}

// 	request := &protocol.OrderListRequest{}
// 	err := r.ParseForm()
// 	if err != nil {
// 		return nil, err
// 	}

// 	request.Offset = form.IntOrDefault("offset", 0)
// 	request.Count = form.IntOrDefault("count", -1)
// 	request.PayBy = strings.ToLower(form.Get("pay_by"))
// 	request.Status = uint8(form.IntOrDefault(form.Get("status"), 0))
// 	request.AppID = form.Get("appid")
// 	request.ChannelID = form.Get("channel_id")
// 	request.Start = form.Int64OrDefault(form.Get("start"), -1)
// 	request.End = form.Int64OrDefault(form.Get("end"), -1)
// 	request.Uid = form.Get("uid")
// 	request.OrderID = form.Get("order_id")

// 	list, total, err := OrderList(request)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &protocol.OrderListResponse{Data: list, Total: total}, nil
// }

// func WechatNotify(r *protocol.WechatOrderCallbackRequest) (resp interface{}, err error) {
// 	var trade *model.Trade
// 	var order *model.Order
// 	if trade, resp, err = provider.Wechat.Notify(r); err != nil {
// 		logger.Error(err.Error())
// 		return nil, err
// 	}

// 	if order, err = db.QueryOrder(trade.OrderId); err != nil {
// 		logger.Error(err.Error())
// 		return nil, err
// 	}

// 	if err := db.InsertTrade(trade); err != nil {
// 		// 如果是重复通知,直接忽略之
// 		if err == errutil.ErrTradeExisted {
// 			return resp, nil
// 		}

// 		logger.Error(err.Error())
// 		return nil, err
// 	}

// 	if err := db.UserAddCoin(order.Uid, int64(10)); err != nil {
// 		logger.Error(err.Error())
// 		return nil, err
// 	}
// 	return resp, nil
// }

// func wechatCallback(_ context.Context, r *http.Request) (interface{}, error) {
// 	request := &protocol.WechatOrderCallbackRequest{}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := xml.Unmarshal(data, request); err != nil {
// 		return nil, err
// 	}

// 	return WechatNotify(request)
// }

// 虚拟支付2.0道具直购 https://developers.weixin.qq.com/minigame/dev/guide/open-ability/virtual-payment/goods.html
// func handleMessagePushPlainText(_ context.Context, r *http.Request) (interface{}, error) {
// 	// 验证消息
// 	query := r.URL.Query()
// 	signature := query.Get("signature")
// 	timestamp := query.Get("timestamp")
// 	nonce := query.Get("nonce")

// 	token := viper.GetString("wechatv2.token")
// 	if token == "" {
// 		logger.Errorf("微信Android, 验证失败: token 未配置")
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "Token not configured"}, nil
// 	}

// 	if !provider.WeChatV2.HandleWechatMessagePushVerifyURL(token, timestamp, nonce, signature) {
// 		logger.Errorf("微信Android, 验证失败: 无效的签名")
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "Signature verification failed"}, nil
// 	}

// 	// 读取请求体
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		logger.Errorf("微信Android, 读取请求体失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "读取请求体失败"}, nil
// 	}

// 	// 解析请求JSON
// 	var request protocol.WxPayNotifyRequest
// 	if err := json.Unmarshal(body, &request); err != nil {
// 		logger.Errorf("微信Android, 解析JSON失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "解析JSON失败"}, nil
// 	}

// 	logger.Debugf("微信Android支付回调请求: %+v", request)

// 	// 验证消息类型和事件类型
// 	if request.MsgType != "event" {
// 		logger.Errorf("微信Android支付回调, 消息类型错误: %s", request.MsgType)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "消息类型错误"}, nil
// 	}

// 	if request.Event != "minigame_game_pay_goods_deliver_notify" {
// 		logger.Errorf("微信Android支付回调, 事件类型错误: %s", request.Event)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "事件类型错误"}, nil
// 	}

// 	// 解析Payload - 使用动态解析来处理可能缺失的字段
// 	var payload protocol.Payload
// 	if err := json.Unmarshal([]byte(request.MiniGame.Payload), &payload); err != nil {
// 		logger.Errorf("微信Android支付回调, 解析Payload失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "解析Payload失败"}, nil
// 	}

// 	// 验证签名 PayEventSig
// 	if !CalcPayEventSig(request.MiniGame.PayEventSig, request.Event, request.MiniGame.Payload, provider.AppKey) {
// 		logger.Errorf("微信Android支付回调, 签名验证失败: %s", request.MiniGame.PayEventSig)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "签名验证失败"}, nil
// 	}

// 	// 查询订单
// 	order, err := db.QueryOrder(payload.OutTradeNo)
// 	if err != nil || order == nil {
// 		logger.Errorf("微信Android支付回调, 查询订单失败或订单不存在: %s, err: %v", payload.OutTradeNo, err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "订单不存在"}, nil
// 	}

// 	// 检查订单状态
// 	if order.Status != db.OrderStatusCreated {
// 		if order.Status == db.OrderStatusNotified {
// 			logger.Infof("微信Android支付回调, 订单已通知: %s", order.OrderId)
// 			return &protocol.WxPayAndroidNotifyResponse{ErrCode: 0, ErrMsg: "Success"}, nil
// 		}
// 		logger.Errorf("微信Android支付回调, 订单状态错误: %s, 当前状态: %v", order.OrderId, order.Status)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "订单状态错误"}, nil
// 	}

// 	// 验证金额（转换为分）
// 	if !request.MiniGame.IsMock { // 非测试环境才验证金额
// 		expectedAmount := order.Money * 100 // 转换为分
// 		if int(payload.GoodsInfo.ActualPrice) != expectedAmount {
// 			logger.Errorf("微信Android支付回调, 金额验证失败: 预期金额 %d, 实际金额 %d", expectedAmount, payload.GoodsInfo.ActualPrice)
// 			return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "金额验证失败"}, nil
// 		}
// 	}

// 	// 解析商品UUID
// 	goodsUUID, err := strconv.Atoi(payload.GoodsInfo.ProductId)
// 	if err != nil {
// 		logger.Errorf("微信Android支付回调, 解析商品UUID失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "解析商品UUID失败"}, nil
// 	}

// 	// 更新订单状态为已支付
// 	order.Status = db.OrderStatusPayed
// 	order.RealMoney = int(payload.GoodsInfo.ActualPrice) / 100 // 转换回元

// 	if err := db.UpdateOrder(order); err != nil {
// 		logger.Errorf("微信Android支付回调, 更新订单失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "更新订单失败"}, nil
// 	}

// 	storeBuyReq := &protocol.StoreListRequest{
// 		Num:       1,
// 		Notify:    true,
// 		Uid:       order.Uid,
// 		GoodsUuid: int32(goodsUUID),
// 	}
// 	_, err = StoreBuy(storeBuyReq)
// 	if err != nil {
// 		logger.Errorf("微信Android支付回调, 购买商品失败: %v", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "购买商品失败"}, nil
// 	}

// 	logger.Debugf("微信Android支付回调, 订单处理成功: %s", order.OrderId)
// 	return &protocol.WxPayAndroidNotifyResponse{ErrCode: 0, ErrMsg: "Success"}, nil
// }

// 计算 pay_event_sig
// func CalcPayEventSig(payEventSig, event, payload, appKey string) bool {
// 	signMsg := event + "&" + payload

// 	mac := hmac.New(sha256.New, []byte(appKey))
// 	mac.Write([]byte(signMsg))

// 	return hex.EncodeToString(mac.Sum(nil)) == payEventSig
// }

// processCustomerMessageSignature验证微信服务器签名
// func processCustomerMessageSignature(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	signature := query.Get("signature")
// 	timestamp := query.Get("timestamp")
// 	nonce := query.Get("nonce")
// 	echostr := query.Get("echostr")

// 	token := viper.GetString("wechatv2.token")
// 	if token == "" {
// 		logger.Error("微信验证失败: token 未配置")
// 		http.Error(w, "Token not configured", http.StatusInternalServerError)
// 		return
// 	}

// 	if signature == "" || timestamp == "" || nonce == "" {
// 		http.Error(w, "Missing required parameters", http.StatusBadRequest)
// 		return
// 	}

// 	if !provider.WeChatV2.HandleWechatMessagePushVerifyURL(token, timestamp, nonce, signature) {
// 		http.Error(w, "Invalid signature", http.StatusUnauthorized)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	_, err := w.Write([]byte(echostr))
// 	if err != nil {
// 		logger.Error("Failed to write response: ", err)
// 		return
// 	}
// }

// handleWechatMessagePushVerifyURL 验证微信服务器签名 - 原生HTTP Handler
// func handleWechatMessagePushVerifyURL(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
// 	signature := query.Get("signature")
// 	timestamp := query.Get("timestamp")
// 	nonce := query.Get("nonce")
// 	echostr := query.Get("echostr")

// 	token := viper.GetString("wechatv2.token")
// 	if token == "" {
// 		logger.Error("微信验证失败: token 未配置")
// 		http.Error(w, "Token not configured", http.StatusInternalServerError)
// 		return
// 	}

// 	if signature == "" || timestamp == "" || nonce == "" {
// 		http.Error(w, "Missing required parameters", http.StatusBadRequest)
// 		return
// 	}

// 	if !provider.WeChatV2.HandleWechatMessagePushVerifyURL(token, timestamp, nonce, signature) {
// 		http.Error(w, "Invalid signature", http.StatusUnauthorized)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	_, err := w.Write([]byte(echostr))
// 	if err != nil {
// 		logger.Error("Failed to write response: ", err)
// 		return
// 	}
// }

// func alipayCallback(_ context.Context, r *http.Request) (interface{}, error) {
// 	if err := provider.Alipay.VerifySign(r); err != nil {
// 		return nil, err
// 	}

// 	// 读取支付宝POST请求的数据
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 获取支付宝的参数
// 	params := strings.Split(string(body), "&")
// 	amap := make(map[string]string)
// 	for _, param := range params {
// 		pair := strings.Split(param, "=")
// 		if len(pair) == 2 {
// 			amap[pair[0]] = pair[1]
// 		}
// 	}
// 	// 获取订单信息
// 	outTradeNo := amap["out_trade_no"]
// 	tradeNo := amap["trade_no"]
// 	amount := amap["receipt_amount"]
// 	if amap["trade_status"] == "TRADE_FINISHED" || amap["trade_status"] == "TRADE_SUCCESS" {
// 		// 支付成功，更新订单状态
// 		order, err := db.QueryOrder(outTradeNo)
// 		if err != nil || order == nil {
// 			return nil, errutil.ErrOrderNotFound
// 		}
// 		if order.Status != db.OrderStatusCreated {
// 			return nil, errutil.ErrIllegalOrderType
// 		}

// 		productCfg, ok := global.StoreCfg[int32(order.ProductId)]
// 		if ok {
// 			iAmount, _ := strconv.Atoi(amount)
// 			if productCfg.Limit.Num != int64(iAmount) {
// 				return nil, errutil.ErrVerifyFailed
// 			}
// 		}

// 		order.Status = db.OrderStatusPayed
// 		order.ChannelOrderId = tradeNo

// 		if err = db.UpdateOrder(order); err == nil {
// 			sources := []protocol.SourceInfo{}
// 			for _, v := range productCfg.Props {
// 				sources = append(sources, protocol.SourceInfo{Id: v.Id, Num: int64(v.Num), ShopType: constant.ShopTypeRecharge})
// 			}

// 			if err := CheckSourceAdd(sources, order.Uid, true, constant.ItemOrder); err == nil {
// 				order.Status = db.OrderStatusNotified
// 				go db.UpdateOrder(order)
// 			} else {
// 				return nil, err
// 			}
// 		}

// 		log.Println(outTradeNo, "success")
// 	} else {
// 		// 支付失败或者未支付
// 		order, err := db.QueryOrder(outTradeNo)
// 		if err != nil || order == nil {
// 			return nil, errutil.ErrOrderNotFound
// 		}
// 		order.Status = db.OrderStatusNotified
// 		if err = db.UpdateOrder(order); err == nil {
// 			return nil, err
// 		}

// 		log.Println(outTradeNo, "fail")
// 	}
// 	return nil, nil
// }

// Byrty0uxiTradeInfo 定义了支付宝交易信息的结构体
type Byrty0uxiTradeInfo struct {
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号
	TradeNo       string `json:"trade_no"`       // 支付宝订单号
	AppID         string `json:"app_id"`         // 支付宝APPID
	BuyerID       string `json:"buyer_id"`       // 用户支付宝ID
	BuyerLogonID  string `json:"buyer_logon_id"` // 用户支付宝账号
	TotalAmount   string `json:"total_amount"`   // 充值金额
	ReceiptAmount string `json:"receipt_amount"` // 付款金额
	GmtPayment    string `json:"gmt_payment"`    // 充值时间
	NotifyTime    string `json:"notify_time"`    // 通知时间
	Sign          string `json:"sign"`           // 签名，和下单算法一致
}

// byrty0uxi 支付回调
// func byrty0uxiCallback(_ context.Context, r *http.Request) (interface{}, error) {
// 	// 解析表单数据
// 	if err := r.ParseForm(); err != nil {
// 		log.Printf("Error parsing form: %v", err)
// 		return nil, err
// 	}

// 	// 从表单中获取参数
// 	tradeInfo := Byrty0uxiTradeInfo{
// 		OutTradeNo:    r.FormValue("out_trade_no"),
// 		TradeNo:       r.FormValue("trade_no"),
// 		AppID:         r.FormValue("app_id"),
// 		BuyerID:       r.FormValue("buyer_id"),
// 		BuyerLogonID:  r.FormValue("buyer_logon_id"),
// 		TotalAmount:   r.FormValue("total_amount"),
// 		ReceiptAmount: r.FormValue("receipt_amount"),
// 		GmtPayment:    r.FormValue("gmt_payment"),
// 		NotifyTime:    r.FormValue("notify_time"),
// 		Sign:          r.FormValue("sign"),
// 	}

// 	logger.Debugf("支付回调: %v", tradeInfo)

// 	m := make(map[string]interface{})
// 	m["out_trade_no"] = tradeInfo.OutTradeNo
// 	m["trade_no"] = tradeInfo.TradeNo
// 	m["app_id"] = tradeInfo.AppID
// 	m["buyer_id"] = tradeInfo.BuyerID
// 	m["buyer_logon_id"] = tradeInfo.BuyerLogonID
// 	m["total_amount"] = tradeInfo.TotalAmount
// 	m["receipt_amount"] = tradeInfo.ReceiptAmount
// 	m["gmt_payment"] = tradeInfo.GmtPayment
// 	m["notify_time"] = tradeInfo.NotifyTime

// 	if !provider.Byrtypay.CheckSign(m, tradeInfo.Sign) {
// 		logger.Debugf("支付回调，失败 签名错误 %s", tradeInfo.Sign)
// 		return "sign failed ", nil
// 	}

// 	// var PID int64
// 	if tradeInfo.OutTradeNo != "" && tradeInfo.TotalAmount != "" {
// 		// 支付成功，更新订单状态
// 		order, err := db.QueryOrder(tradeInfo.OutTradeNo)
// 		if err != nil || order == nil {
// 			logger.Debugf("支付回调，失败 更新订单失败 %s", tradeInfo.OutTradeNo)
// 			return nil, errutil.ErrOrderNotFound
// 		}
// 		if order.Status != db.OrderStatusCreated {
// 			logger.Debugf("支付回调，订单状态错误 %s %d", tradeInfo.OutTradeNo, order.Status)
// 			return nil, errutil.ErrIllegalOrderType
// 		}

// 		if viper.GetBool("core.debug") {
// 			if tradeInfo.TotalAmount != "0.01" {
// 				logger.Debugf("支付回调，金额校验失败失败1 %s", tradeInfo.TotalAmount)
// 				return nil, errutil.ErrVerifyFailed
// 			}
// 		} else {
// 			iAmount, _ := strconv.Atoi(tradeInfo.TotalAmount)
// 			// 验证金额
// 			if order.Money != iAmount {
// 				logger.Debugf("支付回调，金额校验失败失败2 %s", tradeInfo.TotalAmount)
// 				return nil, errutil.ErrVerifyFailed
// 			}
// 			order.RealMoney = iAmount
// 		}

// 		order.Status = db.OrderStatusPayed
// 		order.ChannelOrderId = tradeInfo.TradeNo

// 		if err = db.UpdateOrder(order); err == nil {
// 			// TODO:发放奖励并更新订单状态
// 			request := &protocol.StoreListRequest{}

// 			request.Num = 1
// 			request.Notify = true
// 			request.Uid = order.Uid
// 			// PID = order.Uid
// 			request.GoodsUuid = int32(order.ProductId)
// 			StoreBuy(request)

// 			order.Status = db.OrderStatusNotified
// 			if err := db.UpdateOrder(order); err != nil {
// 				logger.Debugf("支付回调，更新订单状态失败 %s", tradeInfo.OutTradeNo)
// 			}
// 		}

// 	} else {
// 		// 支付失败或者未支付
// 		order, err := db.QueryOrder(tradeInfo.TradeNo)
// 		if err != nil || order == nil {
// 			logger.Debugf("支付回调失败,订单不存在 %s", tradeInfo.TradeNo)
// 			return nil, errutil.ErrOrderNotFound
// 		}
// 		order.Status = db.OrderStatusNotified
// 		if err = db.UpdateOrder(order); err == nil {
// 			logger.Debugf("支付回调失败,更新订单状态失败 %s", tradeInfo.TradeNo)
// 			return nil, err
// 		}
// 	}

// 	// dispatcher.MustGetInstance().Publish(dispatcher.TopicRechargeSuccess.String(), &dispatcher.RechargeSuccessEvent{
// 	// 	PID:    PID,
// 	// 	Amount: tradeInfo.TotalAmount,
// 	// })

// 	logger.Debugf("支付回调,订单已完成 %s", tradeInfo.TradeNo)
// 	return "sucess", nil
// }
