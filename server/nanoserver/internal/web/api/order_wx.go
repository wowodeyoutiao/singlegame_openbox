package api

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"time"

// 	"nanoserver/db"
// 	"nanoserver/internal/web/api/provider"
// 	constant "nanoserver/pkg/constant"
// 	"nanoserver/pkg/errutil"
// 	"nanoserver/pkg/utils"
// 	"nanoserver/protocol"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/go-xorm/xorm"
// 	"github.com/spf13/viper"
// )

// const (
// 	customerServiceBaseURL  = "https://webchat-sh.clink.cn/chat.html"
// 	customerServiceAccessID = "94a993b1-9fad-4b9e-a732-8074e1cc5578"
// )

// // MessageHandler 定义消息处理器函数类型
// type MessageHandler func(interface{}) (interface{}, error)

// // messageHandlers 消息处理器映射
// var messageHandlers = map[string]MessageHandler{
// 	"*protocol.CustomTextMessage":            handleTextMessage,
// 	"*protocol.CustomImageMessage":           handleImageMessage,
// 	"*protocol.CustomMiniProgramPageMessage": handleMiniProgramMessage,
// 	"*protocol.CustomUserEnterEvent":         handleUserEnterEvent,
// }

// handleMessage 根据消息类型分发处理
// func handleMessage(msg interface{}) (interface{}, error) {
// 	msgType := fmt.Sprintf("%T", msg)

// 	logger.Debugf("处理消息类型: %s, 内容: %+v", msgType, msg)

// 	// 这里只需要处理 EnterEvent 事件，其他的消息类型可以忽略
// 	if handler, exists := messageHandlers[msgType]; exists {
// 		return handler(msg)
// 	}

// 	return nil, fmt.Errorf("unsupported message type: %s", msgType)
// }

// // handleTextMessage 处理文本消息
// func handleTextMessage(msg interface{}) (interface{}, error) {
// 	textMsgReq := msg.(*protocol.CustomTextMessage)

// 	logger.Debugf("收到文本消息: %+v", textMsgReq)

// 	customerServiceURL, err := createCustomerServiceURL(textMsgReq.FromUserName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	textMsgResp := &protocol.TextMessage{
// 		ToUser:  textMsgReq.FromUserName,
// 		MsgType: "text",
// 		Text: &protocol.TextContent{
// 			Content: fmt.Sprintf(`亲爱的玩家，若您在游戏过程中遇到任何问题，欢迎<a href="%s" target="_blank">点击此处联系客服</a>，我们将竭诚为您服务。`, customerServiceURL),
// 		},
// 	}

// 	if err := SendCustomMessage(context.Background(), textMsgResp); err != nil {
// 		logger.Errorf("发送文本消息失败: %v, textMsgResp: %v", err, textMsgResp)
// 		return nil, fmt.Errorf("failed to send custom message: %w, textMsgResp: %v", err, textMsgResp)
// 	}

// 	return nil, nil
// }

type PlayerInfo struct {
	Uid                 int64  `json:"uid"`                   // 用户ID
	Name                string `json:"name"`                  // 用户名
	ServerId            int32  `json:"server_id"`             // 服务器ID
	TotalOnlineTime     int64  `json:"total_online_time"`     // 总在线时长
	Level               int32  `json:"level"`                 // 玩家等级
	VipLevel            int32  `json:"vip_level"`             // VIP等级
	TotalRechargeAmount int    `json:"total_recharge_amount"` // 总充值金额
}

// createCustomerServiceURL 创建客服链接
// func createCustomerServiceURL(openID string) (string, error) {
// 	//  TODO 后续这个操作需要加 Redis 缓存

// 	// 查询用户信息
// 	player, err := db.QueryPlayersByOpenID(openID)
// 	if err != nil && !errors.Is(err, errutil.ErrUserNotFound) {
// 		return "", fmt.Errorf("failed to query user by openID: %w", err)
// 	} else if errors.Is(err, errutil.ErrUserNotFound) {
// 		return fmt.Sprint(customerServiceBaseURL, "?accessId=", customerServiceAccessID, "&language=zh_CN"), nil
// 	}

// 	// 根据用户id 获取他的所有订单
// 	var totalRechargeAmount int
// 	orders, err := db.QueryOrdersByUserID(player.Uid)
// 	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
// 		return "", fmt.Errorf("failed to query orders by user ID: %w", err)
// 	}

// 	for _, v := range orders {
// 		totalRechargeAmount += v.RealMoney
// 	}

// 	playerInfo := &PlayerInfo{
// 		Uid:                 player.Uid,
// 		Name:                player.Name,
// 		ServerId:            player.ServerId,
// 		TotalOnlineTime:     player.TotalOnlineTime,
// 		Level:               player.Level,
// 		VipLevel:            player.Vip,
// 		TotalRechargeAmount: totalRechargeAmount,
// 	}

// 	// 构建客服链接
// 	customerServiceURL, err := buildCustomerServiceURL(playerInfo)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to build customer service URL: %w", err)
// 	}

// 	logger.Debugf("构建的客服链接: %s", customerServiceURL)

// 	return customerServiceURL, nil
// }

// // handleImageMessage 处理图片消息
// func handleImageMessage(msg interface{}) (interface{}, error) {
// 	// imageMsg := msg.(*protocol.CustomImageMessage)
// 	// 处理图片消息逻辑
// 	return nil, nil
// }

// handleMiniProgramMessage 处理小程序卡片消息
// func handleMiniProgramMessage(msg interface{}) (interface{}, error) {
// 	// miniMsg := msg.(*protocol.CustomMiniProgramPageMessage)
// 	// 处理小程序卡片信息
// 	return nil, nil
// }

// // handleUserEnterEvent 处理用户进入事件
// func handleUserEnterEvent(msg interface{}) (interface{}, error) {
// 	userEventReq := msg.(*protocol.CustomUserEnterEvent)

// 	customerServiceURL, err := createCustomerServiceURL(userEventReq.FromUserName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if userEventReq.SessionFrom == "customer service" {
// 		// 如果是客服会话，则不处理
// 		logger.Debugf("收到客服会话事件: %s, 来自用户: %s", userEventReq.Event, userEventReq.FromUserName)
// 		userEventResp := &protocol.TextMessage{
// 			ToUser:  userEventReq.FromUserName,
// 			MsgType: "text",
// 			Text: &protocol.TextContent{
// 				Content: fmt.Sprintf(`亲爱的玩家，若您在游戏过程中遇到任何问题，欢迎<a href="%s" target="_blank">点击此处联系客服</a>，我们将竭诚为您服务。`, customerServiceURL),
// 			},
// 		}

// 		if err := SendCustomMessage(context.Background(), userEventResp); err != nil {
// 			logger.Errorf("发送客服会话消息失败: %v, msg: %v", err, userEventResp)
// 			return nil, fmt.Errorf("failed to send custom message: %w, msg: %v", err, userEventResp)
// 		}

// 		return nil, nil
// 	}

// 	nonceStr, err := utils.GenerateNonceStr(32)
// 	if err != nil {
// 		logger.Error("微信IOS支付生成随机字符串失败: ", err)
// 		return nil, fmt.Errorf("failed to generate nonce string: %w", err)
// 	}

// 	_, err = db.QueryAndUpdateOrderNonce(userEventReq.SessionFrom, nonceStr)
// 	if err != nil {
// 		logger.Error("微信IOS支付查询或更新订单失败: ", err)
// 		return &protocol.WxPayAndroidNotifyResponse{ErrCode: -1, ErrMsg: "查询或更新订单失败"}, nil
// 	}

// 	// 发送给wechatsdk
// 	linkContent := &protocol.LinkContent{
// 		Title:       "点击充值",
// 		Description: "这岛有剧本",
// 		URL:         fmt.Sprintf("%s/v1/order/wx/ios/html?gzhOpenId=%s&nonce=%s", "https://pay.shbinqiwangluo.com", userEventReq.FromUserName, nonceStr),
// 		ThumbURL:    "",
// 	}
// 	linkMessageResp := &protocol.LinkMessage{
// 		ToUser:  userEventReq.FromUserName,
// 		MsgType: "link",
// 		Link:    linkContent,
// 	}

// 	if err := SendCustomMessage(context.Background(), linkMessageResp); err != nil {
// 		logger.Errorf("发送链接消息失败: %v, msg: %v", err, linkMessageResp)
// 		return nil, fmt.Errorf("failed to send custom message: %w, msg: %v", err, linkMessageResp)
// 	}

// 	return nil, nil
// }

// func processCustomerMessage(ctx context.Context, r *http.Request) (interface{}, error) {
// 	// 处理微信验证请求
// 	query := r.URL.Query()
// 	signature := query.Get("signature")
// 	timestamp := query.Get("timestamp")
// 	nonce := query.Get("nonce")

// 	token := viper.GetString("wechatv2.token")
// 	if token == "" {
// 		logger.Error("微信验证失败: token 未配置")
// 		return &protocol.WxPayAndroidNotifyResponse{
// 			ErrCode: -1,
// 			ErrMsg:  "微信验证失败: token 未配置",
// 		}, nil
// 	}

// 	if signature == "" || timestamp == "" || nonce == "" {
// 		logger.Error("微信验证失败: 缺少必要的查询参数")
// 		return &protocol.WxPayAndroidNotifyResponse{
// 			ErrCode: -1,
// 			ErrMsg:  "微信验证失败: 缺少必要的查询参数",
// 		}, nil
// 	}

// 	if !provider.WeChatV2.HandleWechatMessagePushVerifyURL(token, timestamp, nonce, signature) {
// 		logger.Error("微信验证失败: 签名验证不通过")
// 		return &protocol.WxPayAndroidNotifyResponse{
// 			ErrCode: -1,
// 			ErrMsg:  "微信验证失败: 签名验证不通过",
// 		}, nil
// 	}

// 	// 读取请求体
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		logger.Error("读取请求体失败: ", err)
// 		return nil, fmt.Errorf("failed to read request body: %w", err)
// 	}
// 	defer r.Body.Close()

// 	var base protocol.BaseCustomMessage
// 	if err := json.Unmarshal(body, &base); err != nil || base.MsgType == "" {
// 		return nil, fmt.Errorf("failed to parse JSON request body: %w", err)
// 	}

// 	var msg interface{}
// 	var msgDesc string

// 	switch base.MsgType {
// 	case "text":
// 		msg = &protocol.CustomTextMessage{}
// 		msgDesc = "text message"
// 	case "image":
// 		msg = &protocol.CustomImageMessage{}
// 		msgDesc = "image message"
// 	case "event":
// 		msg = &protocol.CustomUserEnterEvent{}
// 		msgDesc = "event message"
// 	default:
// 		return nil, fmt.Errorf("unknown message type: %s", base.MsgType)
// 	}

// 	if err := json.Unmarshal(body, msg); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal %s: %w", msgDesc, err)
// 	}

// 	// 处理
// 	_, err = handleMessage(msg)
// 	if err != nil {
// 		logger.Errorf("处理客服消息失败: %s, 错误: %v", msgDesc, err)
// 	}

// 	return nil, nil
// }

const paymentTemplate = `<!DOCTYPE html>
<html>
  <head />
  <body>
    <script>
      function onBridgeReady() {
        WeixinJSBridge.invoke(
          "getBrandWCPayRequest",
          {
            appId: "{{.AppId}}",
            nonceStr: "{{.NonceStr}}",
            package: "{{.Package}}",
            paySign: "{{.PaySign}}",
            signType: "{{.SignType}}",
            timeStamp: "{{.TimeStamp}}",
          },
          function (res) {
            WeixinJSBridge.call("closeWindow");
          }
        );
      }
      if (typeof WeixinJSBridge == "undefined") {
        if (document.addEventListener) {
          document.addEventListener(
            "WeixinJSBridgeReady",
            onBridgeReady,
            false
          );
        } else if (document.attachEvent) {
          document.attachEvent("WeixinJSBridgeReady", onBridgeReady);
          document.attachEvent("onWeixinJSBridgeReady", onBridgeReady);
        }
      } else {
        onBridgeReady();
      }
    </script>
  </body>
</html>`

// generateWechatRechargeHTML 处理微信公众号充值页面请求（兼容客户端URL格式）
// func generateWechatRechargeHTML(w http.ResponseWriter, r *http.Request) {
// 	// 设置响应头
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// 获取查询参数
// 	query := r.URL.Query()
// 	gzhOpenId := query.Get("gzhOpenId")
// 	// timeStamp := query.Get("timestamp")
// 	nonce := query.Get("nonce")
// 	if gzhOpenId == "" || nonce == "" {
// 		logger.Errorf("缺少必要的查询参数: gzhOpenId=%s, timestamp=%s, nonce=%s", gzhOpenId, nonce)
// 		http.Error(w, "缺少必要的查询参数", http.StatusBadRequest)
// 		return
// 	}

// 	logger.Debugf("微信公众号充值请求: gzhOpenId=%s, nonce=%s", gzhOpenId, nonce)

// 	// 从数据库获取订单
// 	order, err := db.QueryOrderByOpenIDAndNonce(gzhOpenId, nonce)
// 	if err != nil {
// 		logger.Errorf("微信公众号充值查询订单失败: %v", err)
// 		http.Error(w, "查询订单失败", http.StatusBadRequest)
// 	}

// 	timeStamp := strconv.FormatInt(order.CreatedAt, 10)
// 	packageStr := fmt.Sprintf("prepay_id=%s", order.WeChatPrePayID)
// 	signatureBase64, err := provider.WeChatV2.JSAPISign(timeStamp, nonce, packageStr)
// 	if err != nil {
// 		logger.Errorf("微信公众号充值订单签名失败: %v", err)
// 		http.Error(w, "签名失败", http.StatusBadRequest)
// 	}

// 	// 构建模板数据 - 使用真实的微信支付参数
// 	templateData := map[string]string{
// 		"AppId":     viper.GetString("wechatv2.appid"), // 微信支付应用ID
// 		"TimeStamp": timeStamp,                         // 当前时间戳
// 		"NonceStr":  nonce,                             // 生成随机字符串
// 		"Package":   packageStr,                        // 使用订单ID作为预支付ID
// 		"PaySign":   signatureBase64,                   // 使用生成的支付签名
// 		"SignType":  "RSA",                             // 微信支付v3使用HMAC-SHA256
// 	}

// 	tmpl, err := template.New("payment").Parse(paymentTemplate)
// 	if err != nil {
// 		logger.Error("解析模板失败: ", err)
// 		http.Error(w, "解析模板失败", http.StatusInternalServerError)
// 		return
// 	}

// 	// 渲染模板
// 	var buf bytes.Buffer
// 	err = tmpl.Execute(&buf, templateData)
// 	if err != nil {
// 		logger.Error("渲染模板失败: ", err)
// 		http.Error(w, "渲染模板失败", http.StatusInternalServerError)
// 		return
// 	}

// 	// 返回渲染后的HTML
// 	_, _ = w.Write(buf.Bytes())
// }

// func processIOSPaymentCallback(ctx context.Context, r *http.Request) (interface{}, error) {
// 	// 验签
// 	data, err := provider.WeChatV2.ParseNotifyRequest(r)
// 	if err != nil {
// 		logger.Errorf("微信IOS支付通知, 析失败: %v", err)
// 		return &protocol.WxPayIOSNotifyResponse{Code: "FAIL", Message: "解析通知内容失败"}, nil
// 	}

// 	logger.Debugf("收到微信支付通知: %+v", data)

// 	// 处理通知数据
// 	if data.EventType != "TRANSACTION.SUCCESS" {
// 		logger.Warnf("收到非交易成功的通知: %s", data.EventType)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "非交易成功的通知",
// 		}, nil
// 	}

// 	wechatPayCallback := &protocol.WechatPayCallback{}
// 	if err := json.Unmarshal([]byte(data.Resource.Plaintext), &wechatPayCallback); err != nil {
// 		logger.Error("微信IOS支付通知, 内容解析失败: ", err)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "解析通知内容失败",
// 		}, nil
// 	}

// 	logger.Debugf("微信IOS支付通知, 通知内容: %+v", wechatPayCallback)

// 	// 处理数据库操作
// 	order, err := db.QueryOrder(wechatPayCallback.OutTradeNo)
// 	if err != nil || order == nil {
// 		logger.Errorf("微信IOS支付通知, 查询订单失败或订单不存在: %s, err: %v", wechatPayCallback.OutTradeNo, err)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "订单不存在或查询失败",
// 		}, nil
// 	}

// 	// 检查订单状态
// 	if order.Status != db.OrderStatusCreated {
// 		if order.Status == db.OrderStatusPayed {
// 			logger.Infof("微信IOS支付通知, 订单已支付: %s, 当前状态: %s", order.OrderId, order.Status)
// 			return &protocol.WxPayIOSNotifyResponse{
// 				Code:    "SUCCESS",
// 				Message: "订单已支付",
// 			}, nil
// 		}
// 		logger.Errorf("微信IOS支付通知, 订单状态错误: %s, 当前状态: %s", order.OrderId, order.Status)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "订单状态错误",
// 		}, nil
// 	}

// 	// 解析商品UUID
// 	goodsUUID, err := strconv.Atoi(wechatPayCallback.OutTradeNo)
// 	if err != nil {
// 		logger.Errorf("微信IOS支付通知, 解析商品UUID失败: %s, err: %v", wechatPayCallback.OutTradeNo, err)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "解析商品UUID失败",
// 		}, nil
// 	}

// 	// 更新订单状态为已支付
// 	order.Status = db.OrderStatusPayed
// 	order.RealMoney = int(wechatPayCallback.Amount.PayerTotal) / 100

// 	if err := db.UpdateOrder(order); err != nil {
// 		logger.Errorf("微信IOS支付通知, 更新订单状态失败: %s, err: %v", order.OrderId, err)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "更新订单状态失败",
// 		}, nil
// 	}

// 	storeBuyReq := &protocol.StoreListRequest{
// 		Num:       1,
// 		Notify:    true,
// 		Uid:       order.Uid,
// 		GoodsUuid: int32(goodsUUID),
// 	}

// 	_, err = StoreBuy(storeBuyReq)
// 	if err != nil {
// 		logger.Errorf("微信IOS支付通知, 购买商品失败: %s, err: %v", order.OrderId, err)
// 		return &protocol.WxPayIOSNotifyResponse{
// 			Code:    "FAIL",
// 			Message: "购买商品失败",
// 		}, nil
// 	}

// 	return &protocol.WxPayIOSNotifyResponse{
// 		Code:    "SUCCESS",
// 		Message: "支付成功",
// 	}, nil
// }

// buildCustomerServiceURL 构建客服链接，支持传入用户信息和自定义字段
// https://develop.clink.cn/develop/integration/channel.html#_%E7%BD%91%E9%A1%B5%E6%B8%A0%E9%81%93
// func buildCustomerServiceURL(player *PlayerInfo) (string, error) {
// 	logger.Debugf("构建客服链接, Id: %d, Name: %s, ServerId: %d, TotalOnlineTime: %d, Level: %d, VipLevel: %d",
// 		player.Uid, player.Name, player.ServerId, player.TotalOnlineTime, player.Level, player.VipLevel)

// 	baseURL, err := url.Parse(customerServiceBaseURL)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse base URL: %w", err)
// 	}
// 	params := url.Values{}

// 	params.Add("accessId", customerServiceAccessID)
// 	params.Add("language", "zh_CN")

// 	params.Add("externalId", strconv.FormatInt(player.Uid, 10))
// 	params.Add("visitorName", player.Name)

// 	// 添加玩家信息
// 	visitorExtraInfo := map[string]interface{}{
// 		"uid":                   player.Uid,
// 		"name":                  player.Name,
// 		"server_id":             player.ServerId,
// 		"total_online_time":     player.TotalOnlineTime,
// 		"level":                 player.Level,
// 		"vip_level":             player.VipLevel,
// 		"total_recharge_amount": player.TotalRechargeAmount,
// 	}

// 	visitorExtraInfoJSON, err := json.Marshal(visitorExtraInfo)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal visitor extra info: %w", err)
// 	}
// 	params.Add("visitorExtraInfo", url.QueryEscape(string(visitorExtraInfoJSON)))

// 	baseURL.RawQuery = params.Encode()

// 	return baseURL.String(), nil
// }

// Deprecated: In pay service already have this function
// func SendCustomMessage(ctx context.Context, msg interface{}) error {
// 	if msg == nil {
// 		return fmt.Errorf("message cannot be nil")
// 	}
// 	paymentAddr := fmt.Sprintf("%s/v1/wx/message/custom/send", constant.PaymentURL)
// 	httpClient := resty.New().SetTimeout(10 * time.Second).SetRetryCount(3).SetRetryWaitTime(2 * time.Second).AddRetryCondition(
// 		func(r *resty.Response, err error) bool {
// 			return err != nil || (r.StatusCode() >= 500 && r.StatusCode() <= 599)
// 		},
// 	)

// 	msgSendCustomMessageResp, err := httpClient.R().
// 		SetContext(ctx).
// 		SetHeader("Content-Type", "application/json; charset=utf-8").
// 		SetBody(msg).
// 		Post(paymentAddr)
// 	if err != nil {
// 		return fmt.Errorf("failed to send custom message: %v", err)
// 	}

// 	if msgSendCustomMessageResp.StatusCode() != http.StatusOK {
// 		return fmt.Errorf("failed to send custom message, status code: %d", msgSendCustomMessageResp.StatusCode())
// 	}

// 	return nil
// }
