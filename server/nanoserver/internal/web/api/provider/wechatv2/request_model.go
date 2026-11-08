package wechatv2

// import (
// 	"time"

// 	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
// )

// // AccessTokenResponse 微信 access_token 接口响应结构
// type AccessTokenResponse struct {
// 	AccessToken string `json:"access_token"`
// 	ExpiresIn   int    `json:"expires_in"`
// 	ErrCode     int    `json:"errcode"`
// 	ErrMsg      string `json:"errmsg"`
// }

// // Code2SessionResponse 微信小程序登录凭证校验响应结构
// type Code2SessionResponse struct {
// 	OpenID     string `json:"openid"`      // 用户唯一标识
// 	SessionKey string `json:"session_key"` // 会话密钥
// 	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符（可选）
// 	ErrCode    int    `json:"errcode"`     // 错误码
// 	ErrMsg     string `json:"errmsg"`      // 错误信息
// }

// // CheckSessionKeyResponse 微信小程序校验登录态 session_key 响应结构
// type CheckSessionKeyResponse struct {
// 	ErrCode int    `json:"errcode"` // 错误码，0表示成功，87009表示签名错误
// 	ErrMsg  string `json:"errmsg"`  // 错误信息
// }

// // ResetUserSessionKeyResponse 微信小程序重置登录态 session_key 响应结构
// type ResetUserSessionKeyResponse struct {
// 	ErrCode    int    `json:"errcode"`     // 错误码
// 	ErrMsg     string `json:"errmsg"`      // 错误信息
// 	OpenID     string `json:"openid"`      // 用户唯一标识符
// 	SessionKey string `json:"session_key"` // 重置后的用户登录态
// }

// type RequestMidasPaymentGameItem struct {
// 	SignData  string `json:"signData"`  // 签名数据JSON字符串
// 	PaySig    string `json:"paySig"`    // 支付签名
// 	Signature string `json:"signature"` // 用户态签名
// }

// // PrepayRequest 预支付参数结构体
// type PrepayRequest struct {
// 	Description string    `json:"description"`  // 商品描述 [必填]
// 	OutTradeNo  string    `json:"out_trade_no"` // 商户订单号 [必填]
// 	TimeExpire  *time.Time `json:"time_expire"`  // 订单失效时间 [非必填]
// 	Attach      string    `json:"attach"`       // 附加数据 [非必填]
// 	// NotifyUrl     string            `json:"notify_url"`     // 通知地址 [必填]
// 	GoodsTag      string            `json:"goods_tag"`      // 商品标记 [非必填]
// 	LimitPay      []string          `json:"limit_pay"`      // 限制支付方式 [非必填]
// 	SupportFapiao bool              `json:"support_fapiao"` // 是否支持发票 [非必填]
// 	Amount        *jsapi.Amount     `json:"amount"`         // 金额信息 [必填]
// 	Payer         *jsapi.Payer      `json:"payer"`          // 支付者信息 [必填]
// 	Detail        *jsapi.Detail     `json:"detail"`         // 商品详情 [非必填]
// 	SceneInfo     *jsapi.SceneInfo  `json:"scene_info"`     // 场景信息 [非必填]
// 	SettleInfo    *jsapi.SettleInfo `json:"settle_info"`    // 结算信息 [非必填]
// }

// // PrepayResponse 预支付响应结构体
// type PrepayResponse struct {
// 	PrepayID string `json:"prepay_id"` // 预支付交易会话标识 [必填] string
// }

// // PrepayErrorResponse 预支付响应结构体
// type PrepayErrorResponse struct {
// 	Code    string             `json:"code"`    // 错误码 [必填] string
// 	Message string             `json:"message"` // 错误信息 [必填] string
// 	Detail  *PrepayErrorDetail `json:"detail"`  // 错误详情 [可选] struct
// }

// // 错误详情结构体
// type PrepayErrorDetail struct {
// 	Field    string `json:"field"`    // 出错字段路径 [必填] string
// 	Value    string `json:"value"`    // 错误值 [必填] string
// 	Issue    string `json:"issue"`    // 错误原因 [必填] string
// 	Location string `json:"location"` // 错误位置 [必填] string
// }

// // SendCustomMessageRequest 发送客服消息请求结构体
// type SendCustomMessageRequest struct {
// 	ToUser          string                  `json:"touser"`                    // 用户的OpenID
// 	MsgType         string                  `json:"msgtype"`                   // 消息类型
// 	Text            *TextMessage            `json:"text,omitempty"`            // 文本消息，msgtype="text"时必填
// 	Image           *ImageMessage           `json:"image,omitempty"`           // 图片消息，msgtype="image"时必填
// 	Link            *LinkMessage            `json:"link,omitempty"`            // 图文链接，msgtype="link"时必填
// 	MiniProgramPage *MiniProgramPageMessage `json:"miniprogrampage,omitempty"` // 小程序卡片，msgtype="miniprogrampage"时必填
// }

// // TextMessage 文本消息结构体
// type TextMessage struct {
// 	Content string `json:"content"` // 文本消息内容
// }

// // ImageMessage 图片消息结构体
// type ImageMessage struct {
// 	MediaID string `json:"media_id"` // 发送的图片的媒体ID
// }

// // LinkMessage 图文链接消息结构体
// type LinkMessage struct {
// 	Title       string `json:"title"`       // 消息标题
// 	Description string `json:"description"` // 图文链接消息描述
// 	URL         string `json:"url"`         // 图文链接消息被点击后跳转的链接
// 	ThumbURL    string `json:"thumb_url"`   // 图文链接消息的图片链接
// }

// // MiniProgramPageMessage 小程序卡片消息结构体
// type MiniProgramPageMessage struct {
// 	Title        string `json:"title"`          // 消息标题
// 	PagePath     string `json:"pagepath"`       // 小程序的页面路径
// 	ThumbMediaID string `json:"thumb_media_id"` // 小程序消息卡片的封面media_id
// }

// // SendCustomMessageResponse 发送客服消息返回结构体
// type SendCustomMessageResponse struct {
// 	ErrCode int    `json:"errcode"` // 错误码
// 	ErrMsg  string `json:"errmsg"`  // 错误信息
// }
