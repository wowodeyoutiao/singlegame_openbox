package wechatv2

// // MiniGameGoodsDeliverNotifyRequest 微信小程序道具发货通知请求结构
// // https://docs.qq.com/doc/DR1hhWlpnQXJXWHRh?_t=1749823652587&nlc=1
// type MiniGameGoodsDeliverNotifyRequest struct {
// 	ToUserName   string    `json:"ToUserName"`   // 小游戏原始ID
// 	FromUserName string    `json:"FromUserName"` // 该事件消息的openid，道具发货场景固定为微信官方的openid
// 	CreateTime   int64     `json:"CreateTime"`   // 消息发送时间
// 	MsgType      string    `json:"MsgType"`      // 消息类型，道具发货场景固定为：event
// 	Event        string    `json:"Event"`        // 事件类型
// 	MiniGame     *MiniGame `json:"MiniGame"`     // 道具直购发货参数
// }

// // MiniGame 道具直购发货参数
// type MiniGame struct {
// 	Payload     string `json:"Payload"`     // 携带的具体内容，格式为json
// 	PayEventSig string `json:"PayEventSig"` // 支付事件签名
// 	IsMock      bool   `json:"IsMock"`      // True: 模拟测试推送  False：真实推送
// }

// // GoodsDeliverPayload 道具发货载荷数据
// type GoodsDeliverPayload struct {
// 	OpenId     string     `json:"OpenId"`     // 接收道具的玩家openid
// 	Env        int        `json:"Env"`        // 环境配置 0：现网环境 1：沙箱环境
// 	OutTradeNo string     `json:"OutTradeNo"` // 订单号
// 	GoodsInfo  *GoodsInfo `json:"GoodsInfo"`  // 发货道具
// }

// // GoodsInfo 发货道具信息
// type GoodsInfo struct {
// 	ProductId   string `json:"ProductId"`   // 游戏道具id标识
// 	Quantity    int    `json:"Quantity"`    // 购买道具数量
// 	ZoneId      string `json:"ZoneId"`      // 分区
// 	OrigPrice   int    `json:"OrigPrice"`   // 物品原始价格（单位：分）
// 	ActualPrice int    `json:"ActualPrice"` // 物品实际支付价格（单位：分）
// 	Attach      string `json:"Attach"`      // 透传数据
// 	OrderSource int    `json:"OrderSource"` // 1 游戏内 2 商城下单 3 商城测试下单
// }

// // MiniGameGoodsDeliverNotifyResponse 微信小程序道具发货通知响应结构
// type MiniGameGoodsDeliverNotifyResponse struct {
// 	ErrCode int    `json:"ErrCode"`          // 发送状态。0：成功，其他：失败
// 	ErrMsg  string `json:"ErrMsg,omitempty"` // 错误原因，用于调试。在errcode非0的情况下可以返回
// }

// // HandleMessagePushPlainTextRequest Android 支付回调，只验证签名
// type HandleMessagePushPlainTextRequest struct {
// 	// URL 参数
// 	Signature string `form:"signature"` // 微信加密签名
// 	Timestamp string `form:"timestamp"` // 时间戳
// 	Nonce     string `form:"nonce"`     // 随机数
// 	Token     string // 微信开发者设置的 token

// 	// 消息体字段
// 	MiniGameRequest *MiniGameGoodsDeliverNotifyRequest `json:"MiniGameGoodsDeliverNotifyRequest"` // 小游戏道具发货通知请求
// }

// type HandleMessagePushPlainTextResponse struct {
// 	ErrCode int    `json:"ErrCode"`          // 发送状态。0：成功，其他：失败
// 	ErrMsg  string `json:"ErrMsg,omitempty"` // Success
// }

// // HandleMessagePushCipherTextRequest Android 支付回调，验证签名和解密消息体
// type HandleMessagePushCipherTextRequest struct {
// 	// URL 参数
// 	Signature    string `form:"signature"`     // 微信加密签名
// 	Timestamp    string `form:"timestamp"`     // 时间戳
// 	Nonce        string `form:"nonce"`         // 随机数
// 	OpenID       string `form:"openid"`        // 消息发送者的 OpenID
// 	EncryptType  string `form:"encrypt_type"`  // 加密类型，固定为 aes
// 	MsgSignature string `form:"msg_signature"` // 消息签名
// 	Token        string // 微信开发者设置的 token

// 	// 消息体字段
// 	ToUserName string `json:"ToUserName"` // 小程序的原始ID
// 	Encrypt    string `json:"Encrypt"`    // 加密的消息体
// }

// type HandleMessagePushCipherTextResponse struct {
// 	ErrCode int    `json:"ErrCode"`          // 发送状态。0：成功，其他：失败
// 	ErrMsg  string `json:"ErrMsg,omitempty"` // Success
// }
