package protocol

import (
	"time"
)

type CreateOrderRequest struct {
	AppID          string // 来自哪个应用的订单
	ChannelID      string // 来自哪个渠道的订单
	Platform       string // 支付平台
	ProductId      int    // 商品唯一ID
	ProductionName string
	ProductCount   int    // 数量
	Extra          string // 描述信息
	Device         Device // 设备信息
	Uid            int64  // 用户Io
	Amount         int    // 金额，单位分
	OpenID         string // 微信OpenID
	ServerID       uint64 // 服务器ID
	Attach         string
	RequestID      int64 // 请求ID
}

type CreateOrderByAdminRequest struct {
	AppID     string `json:"appid"`      // 来自哪个应用的订单
	ChannelID string `json:"channel_id"` // 来自哪个渠道的订单
	Extra     string `json:"extra"`      // 描述信息
	Operator  string `json:"operator"`   // 管理员账号
	Money     int    `json:"money"`      // 金额
	Uid       int64  `json:"uid"`        // 用户ID
	Device    Device `json:"device"`     // 设备信息
}

type BalanceListRequest struct {
	Uids []string `json:"uids"` // uid列表
}

type BalanceListResponse struct {
	Code int               `json:"code"` // 状态码
	Data map[string]string `json:"data"` // 渠道列表
}

// OrderByAdminListRequest 由管理员创建的订单列表
type OrderByAdminListRequest struct {
	Offset    int    `json:"offset"`
	Count     int    `json:"count"`
	Start     int64  `json:"start"` // 时间起点
	End       int64  `json:"end"`   // 时间终点
	Uid       int64  `json:"uid"`   // 用户id
	OrderId   string `json:"order_id"`
	AppID     string `json:"appid"`      // 来自哪个应用的订单
	ChannelID string `json:"channel_id"` // 来自哪个渠道的订单
}

// PayOrderListRequest 由管理员创建的订单列表
type PayOrderListRequest struct {
	Offset    int    `json:"offset"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`       // 查询类型: 1-购买代币 2-消费代币
	Start     int64  `json:"start"`      // 时间起点
	End       int64  `json:"end"`        // 时间终点
	Uid       int64  `json:"uid"`        // 用户id
	OrderID   string `json:"order_id"`   // 订单号
	AppID     string `json:"appid"`      // 来自哪个应用的订单
	ChannelID string `json:"channel_id"` // 来自哪个渠道的订单
}

// OrderListRequest 订单列表
type OrderListRequest struct {
	Offset    int    `json:"offset"`
	Count     int    `json:"count"`
	Status    uint8  `json:"status"`
	Start     int64  `json:"start"` // 时间起点
	End       int64  `json:"end"`   // 时间终点
	PayBy     string `json:"pay_by"`
	Uid       string `json:"uid"`        // 用户id
	OrderID   string `json:"order_id"`   // 订单号
	AppID     string `json:"appid"`      // 来自哪个应用的订单
	ChannelID string `json:"channel_id"` // 来自哪个渠道的订单
}

// TradeListRequest 交易列表
type TradeListRequest struct {
	Offset    int    `json:"offset"`
	Count     int    `json:"count"`
	Start     int64  `json:"start"`      // 时间起点
	End       int64  `json:"end"`        // 时间终点
	OrderID   string `json:"order_id"`   // 订单号
	AppID     string `json:"appid"`      // 来自哪个应用的订单
	ChannelID string `json:"channel_id"` // 来自哪个渠道的订单
}

type PayOrderListResponse struct {
	Code  int                 `json:"code"`  // 状态码
	Data  []SnakePayOrderInfo `json:"data"`  // 渠道列表
	Total int                 `json:"total"` // 总数量
}

type OrderListResponse struct {
	Code  int         `json:"code"`  // 状态码
	Data  []OrderInfo `json:"data"`  // 渠道列表
	Total int         `json:"total"` // 总数量
}

type TradeListResponse struct {
	Code  int         `json:"code"`  // 状态码
	Data  []TradeInfo `json:"data"`  // 渠道列表
	Total int         `json:"total"` // 总数量
}

type ObtainBalanceReqeust struct {
	Token string `json:"token"`
}

type ObtainBalanceResponse struct {
	Code int   `json:"code"`
	Data int64 `json:"data"`
}

type CreateOrderSnakeResponse struct {
	Result      string `json:"result"`
	PayPlatform string `json:"pay_platform"`
}

type CreateOrderWechatReponse struct {
	AppID     string `json:"appid"`
	PartnerId string `json:"partnerid"`
	OrderId   string `json:"orderid"`
	PrePayID  string `json:"prepayid"`
	NonceStr  string `json:"noncestr"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Extra     string `json:"extData"`
}

type UnifyOrderCallbackRequest struct {
	PayPlatform string
	RawRequest  interface{}
}

type WechatOrderCallbackRequest struct {
	ReturnMsg        string `xml:"return_msg,omitempty"`
	DeviceInfo       string `xml:"device_info,omitempty"`
	ErrCode          string `xml:"err_code,omitempty"`
	ErrCodeDes       string `xml:"err_code_des,omitempty"`
	Attach           string `xml:"attach,omitempty"`
	CashFeeType      string `xml:"cash_fee_type,omitempty"`
	CouponFee        int    `xml:"coupon_fee,omitempty"`
	CouponCount      int    `xml:"coupon_count,omitempty"`
	CouponIDDollarN  string `xml:"coupon_id_$n,omitempty"`
	CouponFeeDollarN string `xml:"coupon_fee_$n,omitempty"`

	ReturnCode    string `xml:"return_code"`
	Appid         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	Nonce         string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	Openid        string `xml:"openid"`
	IsSubscribe   string `xml:"is_subscribe"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      int    `xml:"total_fee"`
	FeeType       string `xml:"fee_type"`
	CashFee       int    `xml:"cash_fee"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	TimeEnd       string `xml:"time_end"`

	Raw string
}

type WechatOrderCallbackResponse struct {
	ReturnCode string `xml:"return_code,cdata"`
	ReturnMsg  string `xml:"return_msg,cdata"`
}

type RechargeRequest struct {
	Count int64 `json:"count"`
	Uid   int64 `json:"uid"`
}

type FundTransParams struct {
	Name         string `json:"name"`         // 商户name，外部获得（可选）
	AliName      string `json:"aliName"`      // 用户姓名
	Identity     string `json:"identity"`     // 支付宝账号
	IdentityType string `json:"identityType"` // 类型（可选）
	OrderTitle   string `json:"orderTitle"`   // 标题
	OutBizNo     string `json:"outBizNo"`     // 订单号
	Remark       string `json:"remark"`       // 备注
	TransAmount  string `json:"transAmount"`  // 金额，元
}

// 微信回调主结构体
type WxPayNotifyRequest struct {
	ToUserName   string    `json:"ToUserName"`   // 小游戏原始ID
	FromUserName string    `json:"FromUserName"` // 固定为微信官方OpenID
	CreateTime   int64     `json:"CreateTime"`   // 消息发送时间戳
	MsgType      string    `json:"MsgType"`      // 消息类型(event)
	Event        string    `json:"Event"`        // 事件类型(minigame_game_pay_goods_deliver_notify)
	MiniGame     *MiniGame `json:"MiniGame"`     // 道具直购参数
}

// MiniGame 嵌套结构体
type MiniGame struct {
	Payload     string `json:"Payload"`     // JSON字符串格式的负载数据
	PayEventSig string `json:"PayEventSig"` // 事件签名 [[1]]
	IsMock      bool   `json:"IsMock"`      // 是否为测试推送
}

// Payload 结构体
type Payload struct {
	OpenId     string     `json:"OpenId"`     // 接收道具的玩家OpenID
	Env        int        `json:"Env"`        // 环境配置(0:正式 1:沙箱)
	OutTradeNo string     `json:"OutTradeNo"` // 订单号
	GoodsInfo  *GoodsInfo `json:"GoodsInfo"`  // 发货道具信息
	Attach     string     `json:"Attach"`     // 透传数据
}

// GoodsInfo 结构体
type GoodsInfo struct {
	ProductId   string `json:"ProductId"`   // 游戏道具ID标识
	Quantity    int    `json:"Quantity"`    // 购买道具数量
	ZoneId      string `json:"ZoneId"`      // 分区
	OrigPrice   int    `json:"OrigPrice"`   // 原始价格(单位：分)
	ActualPrice int    `json:"ActualPrice"` // 实际支付价格(单位：分)
	Attach      string `json:"Attach"`      // 透传数据
	OrderSource int    `json:"OrderSource"` // 订单来源(1:游戏内 2:商城 3:商城测试)
}

// 微信回调响应结构体
type WxPayAndroidNotifyResponse struct {
	ErrCode int    `json:"ErrCode"` // 状态码：0 表示成功，其他值表示失败
	ErrMsg  string `json:"ErrMsg"`  // 错误信息（非0时返回）
}

type WxPayIOSNotifyResponse struct {
	Code    string `json:"code"`    // 状态码：0 表示成功，其他值表示失败
	Message string `json:"message"` // 错误信息（非0时返回）
}

// IOS 支付结果通知请求体
type MerchantPayNotifyRequest struct {
	ID           string          `json:"id"`            // 通知唯一ID [必填]
	CreateTime   string          `json:"create_time"`   // 通知创建时间 [必填] time.Time
	EventType    string          `json:"event_type"`    // 事件类型（固定TRANSACTION.SUCCESS）[必填]
	ResourceType string          `json:"resource_type"` // 资源类型（固定encrypt-resource）[必填]
	Summary      string          `json:"summary"`       // 通知摘要 [必填]
	Resource     ResourceWrapper `json:"resource"`      // 加密资源数据 [必填]
}

// 资源包装结构体（用于解密）
type ResourceWrapper struct {
	Algorithm      string `json:"algorithm"`       // 加密算法（AEAD_AES_256_GCM）[必填]
	Ciphertext     string `json:"ciphertext"`      // Base64密文 [必填]
	AssociatedData string `json:"associated_data"` // 附加数据 [可选]
	Nonce          string `json:"nonce"`           // 随机串 [必填]
}

type MerchantNotifyResponse struct {
	Code    string `xml:"code"`    // SUCCESS/FAIL
	Message string `xml:"message"` // OK/错误描述
}

// BaseCustomMessage 定义了所有客服消息/事件共有的基础字段。
type BaseCustomMessage struct {
	ToUserName   string `json:"ToUserName" xml:"ToUserName"`           // 小程序的原始 ID
	FromUserName string `json:"FromUserName" xml:"FromUserName"`       // 发送者的 OpenID
	CreateTime   int64  `json:"CreateTime" xml:"CreateTime"`           // 消息创建时间（Unix 时间戳）
	MsgType      string `json:"MsgType" xml:"MsgType"`                 // 消息类型，如 text, image, event 等
	MsgId        int64  `json:"MsgId,omitempty" xml:"MsgId,omitempty"` // 消息唯一标识 ID（部分消息无此字段）
}

// CustomTextMessage 表示用户发送的文本消息。
type CustomTextMessage struct {
	BaseCustomMessage
	Content string `json:"Content" xml:"Content"` // 文本消息内容
}

// CustomImageMessage 表示用户发送的图片消息。
type CustomImageMessage struct {
	BaseCustomMessage
	PicUrl  string `json:"PicUrl" xml:"PicUrl"`   // 图片链接地址（由微信服务器生成）
	MediaId string `json:"MediaId" xml:"MediaId"` // 图片媒体资源 ID，可用于获取临时素材
}

// CustomMiniProgramPageMessage 表示用户发送的小程序卡片消息。
type CustomMiniProgramPageMessage struct {
	BaseCustomMessage
	Title        string `json:"Title" xml:"Title"`               // 卡片标题
	AppId        string `json:"AppId" xml:"AppId"`               // 被分享小程序的 AppID
	PagePath     string `json:"PagePath" xml:"PagePath"`         // 小程序页面路径
	ThumbUrl     string `json:"ThumbUrl" xml:"ThumbUrl"`         // 缩略图 CDN 地址
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId"` // 缩略图素材 ID
}

// CustomUserEnterEvent 表示用户点击客服会话按钮进入会话时触发的事件。
type CustomUserEnterEvent struct {
	BaseCustomMessage
	Event       string `json:"Event" xml:"Event"`             // 事件类型，固定为 user_enter_tempsession
	SessionFrom string `json:"SessionFrom" xml:"SessionFrom"` // 会话来源标识，由前端传入
}

// LinkMessage 定义发送 link 消息的整体结构
type LinkMessage struct {
	ToUser  string       `json:"touser"`
	MsgType string       `json:"msgtype"` // 固定为 "link"
	Link    *LinkContent `json:"link"`
}

// LinkContent 定义 link 消息的具体内容
type LinkContent struct {
	Title       string `json:"title"`       // 消息标题
	Description string `json:"description"` // 消息描述
	URL         string `json:"url"`         // 消息链接
	ThumbURL    string `json:"thumb_url"`   // 缩略图 URL
}

// TextMessage 定义发送 text 消息的整体结构
type TextMessage struct {
	ToUser  string       `json:"touser"`
	MsgType string       `json:"msgtype"` // 固定为 "text"
	Text    *TextContent `json:"text"`
}

// TextContent 定义 text 消息的具体内容
type TextContent struct {
	Content string `json:"content"` // 消息内容
}

// 微信支付回调数据结构体
type WechatPayCallback struct {
	ID              string      `json:"id"`                         // 微信支付订单号
	SpAppID         string      `json:"sp_appid"`                   // 机构APPID
	SpMchID         string      `json:"sp_mchid"`                   // 机构商户号
	SubMchID        string      `json:"sub_mchid"`                  // 子商户号
	OutTradeNo      string      `json:"out_trade_no"`               // 商户订单号
	Payer           Payer       `json:"payer"`                      // 支付者信息
	Amount          Amount      `json:"amount"`                     // 订单金额信息
	TradeType       string      `json:"trade_type"`                 // 交易类型
	TradeState      string      `json:"trade_state"`                // 交易状态
	TradeStateDesc  string      `json:"trade_state_desc"`           // 交易状态描述
	BankType        string      `json:"bank_type"`                  // 付款银行
	Attach          string      `json:"attach"`                     // 商户附加数据
	SuccessTime     time.Time   `json:"success_time"`               // 支付完成时间
	PromotionDetail []Promotion `json:"promotion_detail,omitempty"` // 优惠详情（可选）
}

// 支付者信息
type Payer struct {
	SpOpenID string `json:"sp_openid"` // 用户在机构appid下的唯一标识
}

// 订单金额信息
type Amount struct {
	Total         int64         `json:"total"`                    // 订单总金额（单位：分）
	Currency      string        `json:"currency,omitempty"`       // 货币类型
	PayerTotal    int64         `json:"payer_total"`              // 用户实际支付金额（单位：分）
	PayerCurrency string        `json:"payer_currency,omitempty"` // 支付货币类型
	ExchangeRate  *ExchangeRate `json:"exchange_rate,omitempty"`  // 汇率信息（可选）
}

// 汇率信息
type ExchangeRate struct {
	Type string `json:"type"` // 汇率类型
	Rate int64  `json:"rate"` // 汇率值（例如：8000000 表示 8.000000）
}

// 优惠信息
type Promotion struct {
	PromotionID               string  `json:"promotion_id"`                          // 券ID
	Name                      string  `json:"name,omitempty"`                        // 优惠名称
	Scope                     string  `json:"scope,omitempty"`                       // 优惠范围
	Type                      string  `json:"type,omitempty"`                        // 优惠类型
	Amount                    int64   `json:"amount"`                                // 优惠券面额
	Currency                  string  `json:"currency,omitempty"`                    // 货币类型
	ActivityID                string  `json:"activity_id,omitempty"`                 // 活动ID
	WechatpayContributeAmount int64   `json:"wechatpay_contribute_amount,omitempty"` // 微信出资金额
	MerchantContributeAmount  int64   `json:"merchant_contribute_amount,omitempty"`  // 商户出资金额
	OtherContributeAmount     int64   `json:"other_contribute_amount,omitempty"`     // 其他出资金额
	GoodsDetail               []Goods `json:"goods_detail,omitempty"`                // 单品列表（可选）
}

// 单品信息
type Goods struct {
	GoodsID     string `json:"goods_id"`               // 商品编码
	GoodsRemark string `json:"goods_remark,omitempty"` // 商品备注
	Quantity    int64  `json:"quantity"`               // 商品数量
	Price       int64  `json:"price"`                  // 商品价格（单位：分）
}

type PayNotifyRequest struct {
	Uid                 uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`                                                              // 用户ID
	StoreId             int64  `protobuf:"varint,2,opt,name=store_id,json=storeId,proto3" json:"store_id,omitempty"`                                       // 商店ID
	Num                 int32  `protobuf:"varint,3,opt,name=num,proto3" json:"num,omitempty"`                                                              // 购买数量, pay 不会传
	RequestID           int64  `protobuf:"varint,4,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`                                 // 请求ID
	IsFirstRecharge     bool   `protobuf:"varint,4,opt,name=is_first_recharge,json=isFirstRecharge,proto3" json:"is_first_recharge,omitempty"`             // 是否首冲
	TransactionId       string `protobuf:"bytes,5,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`                      // 微信支付订单号
	OutTradeNo          string `protobuf:"bytes,6,opt,name=out_trade_no,json=outTradeNo,proto3" json:"out_trade_no,omitempty"`                             // 商户订单号
	TotalCharacterSpend int64  `protobuf:"varint,7,opt,name=total_character_spend,json=totalCharacterSpend,proto3" json:"total_character_spend,omitempty"` // 本次充值后用户总消费金额
}

type PayNotifyResponse struct {
	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`      // 返回状态码，如 "SUCCESS"
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // 错误信息
}

type PayUserInfoRequest struct {
	Uid    int64  `json:"uid"`     // 	用户ID
	OpenID string `json:"open_id"` // 微信OpenID
}

type PayUserInfoResponse struct {
	Uid      int64  `json:"uid"`       // 用户ID
	Name     string `json:"name"`      // 用户名
	ServerId int64  `json:"server_id"` // 服务器ID
	Level    int    `json:"level"`     // 用户等级
	VipLevel int    `json:"vip_level"` // VIP等级
}
