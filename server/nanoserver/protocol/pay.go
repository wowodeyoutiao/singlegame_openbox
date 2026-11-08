package protocol

// 米大师支付原串数据
type SignData struct {
	Mode         string `protobuf:"bytes,1,opt,name=mode,proto3" json:"mode,omitempty"`                 // 支付的类型
	OfferId      string `protobuf:"bytes,2,opt,name=offerId,proto3" json:"offerId,omitempty"`           // 在米大师侧申请的应用id
	BuyQuantity  int32  `protobuf:"varint,3,opt,name=buyQuantity,proto3" json:"buyQuantity,omitempty"`  // 购买数量
	Env          int32  `protobuf:"varint,4,opt,name=env,proto3" json:"env,omitempty"`                  // 环境配置
	CurrencyType string `protobuf:"bytes,5,opt,name=currencyType,proto3" json:"currencyType,omitempty"` // 币种
	Platform     string `protobuf:"bytes,6,opt,name=platform,proto3" json:"platform,omitempty"`         // 平台
	ZoneId       string `protobuf:"bytes,7,opt,name=zoneId,proto3" json:"zoneId,omitempty"`             // 分区ID
	ProductId    string `protobuf:"bytes,8,opt,name=productId,proto3" json:"productId,omitempty"`       // 道具ID
	GoodsPrice   int32  `protobuf:"varint,9,opt,name=goodsPrice,proto3" json:"goodsPrice,omitempty"`    // 道具单价（分）
	OutTradeNo   string `protobuf:"bytes,10,opt,name=outTradeNo,proto3" json:"outTradeNo,omitempty"`    // 业务订单号
	Attach       string `protobuf:"bytes,11,opt,name=attach,proto3" json:"attach,omitempty"`            // 透传数据
}

type MidasPaymentGameItemRequest struct {
	SignData *SignData `protobuf:"bytes,1,opt,name=sign_data,json=signData,proto3" json:"sign_data,omitempty" binding:"required"` // 米大师支付原串数据
	// string pay_sig = 2;          // 支付签名
	// string signature = 3;        // 用户态签名
	OpenId    string `protobuf:"bytes,4,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty" binding:"required"`        // 用户的 openid
	UserId    uint64 `protobuf:"varint,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" binding:"required"`       // 用户的 user_id
	ServerId  uint64 `protobuf:"varint,6,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty" binding:"required"` // 服务器 ID
	RequestId int64  `protobuf:"varint,7,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`                 // 请求 ID
}

type MidasPaymentGameItemResponse struct {
	SignData  string `protobuf:"bytes,1,opt,name=signData,proto3" json:"signData,omitempty"`   // 支付原串数据
	PaySig    string `protobuf:"bytes,2,opt,name=paySig,proto3" json:"paySig,omitempty"`       // 支付签名
	Signature string `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"` // 用户态签名
}

type PrepayRequest struct {
	Appid         string               `protobuf:"bytes,1,opt,name=appid,proto3" json:"appid,omitempty"`                                           // 公众号ID
	Mchid         string               `protobuf:"bytes,2,opt,name=mchid,proto3" json:"mchid,omitempty"`                                           // 直连商户号
	Description   string               `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`                               // 商品描述
	OutTradeNo    string               `protobuf:"bytes,4,opt,name=out_trade_no,json=outTradeNo,proto3" json:"out_trade_no,omitempty"`             // 商户订单号
	TimeExpire    string               `protobuf:"bytes,5,opt,name=time_expire,json=timeExpire,proto3" json:"time_expire,omitempty"`               // 订单失效时间，格式为rfc3339格式
	Attach        string               `protobuf:"bytes,6,opt,name=attach,proto3" json:"attach,omitempty"`                                         // 附加数据
	NotifyUrl     string               `protobuf:"bytes,7,opt,name=notify_url,json=notifyUrl,proto3" json:"notify_url,omitempty"`                  // 有效性：1. HTTPS；2. 不允许携带查询串。
	GoodsTag      string               `protobuf:"bytes,8,opt,name=goods_tag,json=goodsTag,proto3" json:"goods_tag,omitempty"`                     // 商品标记，代金券或立减优惠功能的参数。
	LimitPay      []string             `protobuf:"bytes,9,rep,name=limit_pay,json=limitPay,proto3" json:"limit_pay,omitempty"`                     // 指定支付方式
	SupportFapiao bool                 `protobuf:"varint,10,opt,name=support_fapiao,json=supportFapiao,proto3" json:"support_fapiao,omitempty"`    // 传入true时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效。
	Amount        *PrepayRequestAmount `protobuf:"bytes,11,opt,name=amount,proto3" json:"amount,omitempty"`                                        // 订单金额信息
	Payer         *PrepayRequestPayer  `protobuf:"bytes,12,opt,name=payer,proto3" json:"payer,omitempty"`                                          // 支付者信息
	Detail        *PrepayRequestDetail `protobuf:"bytes,13,opt,name=detail,proto3" json:"detail,omitempty"`                                        // 优惠功能详情
	SceneInfo     *SceneInfo           `protobuf:"bytes,14,opt,name=scene_info,json=sceneInfo,proto3" json:"scene_info,omitempty"`                 // 支付场景描述
	SettleInfo    *SettleInfo          `protobuf:"bytes,15,opt,name=settle_info,json=settleInfo,proto3" json:"settle_info,omitempty"`              // 结算信息
	UserId        uint64               `protobuf:"varint,16,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                         // 用户ID，非必填，若传入则用于记录用户信息
	ServerId      uint64               `protobuf:"varint,6,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty" binding:"required"` // 服务器 ID
	RequestId     int64                `protobuf:"varint,17,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`                // 请求 ID
}

type PrepayRequestAmount struct {
	Total    int64  `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`      // 订单总金额，单位为分
	Currency string `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"` // CNY：人民币，境内商户号仅支持人民币。
}

type PrepayRequestPayer struct {
	Openid string `protobuf:"bytes,1,opt,name=openid,proto3" json:"openid,omitempty"` // 用户在商户appid下的唯一标识。
}

type PrepayRequestDetail struct {
	CostPrice int64 `protobuf:"varint,1,opt,name=cost_price,json=costPrice,proto3" json:"cost_price,omitempty"` // 1.商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。
	// 2.当订单原价与支付金额不相等，则不享受优惠。 3.该字段主要用于防止同一张小票分多次支付，以享受多次优惠的情况，正常支付订单不必上传此参数。
	InvoiceId   string         `protobuf:"bytes,2,opt,name=invoice_id,json=invoiceId,proto3" json:"invoice_id,omitempty"`       // 商家小票ID。
	GoodsDetail []*GoodsDetail `protobuf:"bytes,3,rep,name=goods_detail,json=goodsDetail,proto3" json:"goods_detail,omitempty"` // 商品详情列表
}

type GoodsDetail struct {
	MerchantGoodsId  string `protobuf:"bytes,1,opt,name=merchant_goods_id,json=merchantGoodsId,proto3" json:"merchant_goods_id,omitempty"`    // 商户侧商品编号
	WechatpayGoodsId string `protobuf:"bytes,2,opt,name=wechatpay_goods_id,json=wechatpayGoodsId,proto3" json:"wechatpay_goods_id,omitempty"` // 微信支付定义的统一商品编号（没有可不传）
	GoodsName        string `protobuf:"bytes,3,opt,name=goods_name,json=goodsName,proto3" json:"goods_name,omitempty"`                        // 商品名称
	Quantity         int64  `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`                                          // 商品数量
	UnitPrice        int64  `protobuf:"varint,5,opt,name=unit_price,json=unitPrice,proto3" json:"unit_price,omitempty"`                       // 商品单价，单位为分
}

type SceneInfo struct {
	PayerClientIp string     `protobuf:"bytes,1,opt,name=payer_client_ip,json=payerClientIp,proto3" json:"payer_client_ip,omitempty"` // 用户终端IP
	DeviceId      string     `protobuf:"bytes,2,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`                  // 商户端设备号
	StoreInfo     *StoreInfo `protobuf:"bytes,3,opt,name=store_info,json=storeInfo,proto3" json:"store_info,omitempty"`               // 商户门店信息
}

type SettleInfo struct {
	ProfitSharing bool `protobuf:"varint,1,opt,name=profit_sharing,json=profitSharing,proto3" json:"profit_sharing,omitempty"` // 是否指定分账
}

type PrepayResponse struct {
	PrepayId string `protobuf:"bytes,1,opt,name=prepay_id,json=prepayId,proto3" json:"prepay_id,omitempty"` // 预支付交易会话标识
}

type JSAPIPaymentHTMLRequest struct {
	OpenId    string `protobuf:"bytes,1,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty" form:"open_id" binding:"required"` // @gotags: form:"open_id" binding:"required"
	Nonce     string `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty" form:"nonce" binding:"required"`                   // @gotags: form:"nonce" binding:"required"
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty" form:"timestamp" binding:"required"`      // @gotags: form:"timestamp" binding:"required"
}

type JSAPIPaymentHTMLResponse struct {
	Html []byte `protobuf:"bytes,1,opt,name=html,proto3" json:"html,omitempty"` // 支付HTML内容
}

type MsgSecCheckRequest struct {
	Openid    string `protobuf:"bytes,2,opt,name=openid,proto3" json:"openid,omitempty" binding:"required"`
	Scene     int32  `protobuf:"varint,3,opt,name=scene,proto3" json:"scene,omitempty" binding:"required"`
	Version   int32  `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty" binding:"required"`
	Content   string `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty" binding:"required"`
	Title     string `protobuf:"bytes,6,opt,name=title,proto3" json:"title,omitempty"`         // 可选：文本标题，需使用UTF-8编码
	Nickname  string `protobuf:"bytes,7,opt,name=nickname,proto3" json:"nickname,omitempty"`   // 可选：用户昵称，需使用UTF-8编码
	Signature string `protobuf:"bytes,8,opt,name=signature,proto3" json:"signature,omitempty"` // 可选：个性签名，仅在资料类场景有效(scene=1)
}

type MsgSecCheckResponse struct {
	Errcode int32          `protobuf:"varint,1,opt,name=errcode,proto3" json:"errcode,omitempty"`               // 错误码
	Errmsg  string         `protobuf:"bytes,2,opt,name=errmsg,proto3" json:"errmsg,omitempty"`                  // 错误信息
	TraceId string         `protobuf:"bytes,3,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"` // 唯一请求标识，标记单次请求
	Result  *ResultSummary `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`                  // 综合结果
	Detail  []*DetailItem  `protobuf:"bytes,5,rep,name=detail,proto3" json:"detail,omitempty"`                  // 详细检测结果
}

type ResultSummary struct {
	Suggest string `protobuf:"bytes,1,opt,name=suggest,proto3" json:"suggest,omitempty"` // 建议，有risky、pass、review三种值
	Label   int32  `protobuf:"varint,2,opt,name=label,proto3" json:"label,omitempty"`    // 命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003
}

type DetailItem struct {
	Strategy string `protobuf:"bytes,1,opt,name=strategy,proto3" json:"strategy,omitempty"` // 策略类型
	Errcode  int32  `protobuf:"varint,2,opt,name=errcode,proto3" json:"errcode,omitempty"`  // 错误码，仅当该值为0时，该项结果有效
	Suggest  string `protobuf:"bytes,3,opt,name=suggest,proto3" json:"suggest,omitempty"`   // 建议，有risky、pass、review三种值
	Label    int32  `protobuf:"varint,4,opt,name=label,proto3" json:"label,omitempty"`      // 命中标签枚举值，100 正常；10001 广告；20001 时政；20002 色情；20003
	// 辱骂；20006 违法犯罪；20008 欺诈；20012 低俗；20013 版权；21000 其他
	Keyword string  `protobuf:"bytes,5,opt,name=keyword,proto3" json:"keyword,omitempty"` // 命中的自定义关键词
	Prob    float32 `protobuf:"fixed32,6,opt,name=prob,proto3" json:"prob,omitempty"`     // 0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
}

type VerifyActiveCodeRequest struct {
	Uid      uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty" binding:"required"`
	ServerId int32  `protobuf:"varint,2,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty" binding:"required"`
	Code     string `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty" binding:"required"`
	CodeType int32  `protobuf:"varint,4,opt,name=code_type,json=codeType,proto3" json:"code_type,omitempty" binding:"required"`
	Channel  int32  `protobuf:"varint,5,opt,name=channel,proto3" json:"channel,omitempty" binding:"required"`
}

type VerifyActiveCodeResponse struct {
	Code    int32  `protobuf:"varint,1,opt,name=code,proto3,oneof" json:"code,omitempty"`      // 错误码，0表示成功
	Message string `protobuf:"bytes,2,opt,name=message,proto3,oneof" json:"message,omitempty"` // 错误信息
}
