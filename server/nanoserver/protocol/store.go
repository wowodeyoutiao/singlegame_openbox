package protocol

import "nanoserver/protocol/protobuf/pb"

type CreateStoreRequest struct {
	RegistrationId string // 来自哪个渠道的订单
	Uid            int64  // Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type (
	StoreInfo           struct{}
	StoreOperateRequest struct {
		Uid         int64 //
		AccountName string
	}
)

type StoreOperateResponse struct {
	Code    int    `json:"code"`    // 状态码
	Error   string `json:"error"`   // 状态码
	Data    bool   `json:"data"`    // 渠道列表
	OrderSn string `json:"ordersn"` // 订单号
	// Rmb     bool   `json:"rmb"`     //是否是rmb
	Platform         string `json:"platform"`           // 支付平台标识
	SignData         string `json:"signData"`           // 微信Android支付：签名数据JSON字符串
	PaySig           string `json:"paySig"`             // 微信Android支付：支付签名
	Signature        string `json:"signature"`          // 微信Android支付：用户态签名
	PrePayInternalID string `json:"prepay_internal_id"` // 微信IOS支付：预支付内部交易会话标识
}

type StoreAddRequest struct {
	PropsId int64 `json:"props_id"` // 渠道列表
	Num     int64 `json:"num"`
}

type StoreListRequest struct {
	Uid        int64          `json:"uid"`                   // 用户id
	StoreType  int32          `json:"store_type"`            // 商城类别
	GoodsId    int64          `json:"goods_id"`              // 商品id
	GoodsUuid  int32          `json:"goods_uuid"`            // 商品id
	Detail     int32          `json:"detail"`                // 备注信息
	Num        int32          `json:"num"`                   // 购买数量
	Notify     bool           `json:"notify"`                // 是否是回调订单 默认为false
	Ip         string         `json:"ip"`                    // IP
	Extra      string         `json:"extra"`                 // 额外信息
	Platform   string         `json:"platform"`              // 支付平台标识
	RequestID  int64          `json:"request_id"`            // 请求ID
	RefineSlot *pb.RefineSlot `json:"refine_slot,omitempty"` // 只用于洗练槽位
}

// StoreListResponse 订单列表
type StoreListResponse struct {
	Code int               `json:"code"` // 状态码
	Data []StoreListLittle `json:"data"` // 渠道列表
}

type StoreList struct {
	Buy        int64        `json:"buy"`        // 已购买数量
	GoodsId    int64        `json:"goods_id"`   // 商品id
	GoodsUuid  int64        `json:"goods_uuid"` // 商品id
	Num        int64        `json:"num"`        // 花费的数量
	LimitNum   int64        `json:"limit_num"`  // 总数量
	NeedMoney  int32        `json:"need_money"` // 需要的类型
	Times      int32        `json:"times"`      // 购买多少次免费
	PreLimit   int32        `json:"pre_limit"`  // 前置条件
	SourceInfo []SourceInfo `json:"props"`      // 商品信息
	CreateAt   int32        `json:"create_at"`
	Date       int32        `json:"date"`
	BuyAllNum  int64        `json:"buy_all_num"` // 购买总次数
	Entrance   string       `json:"entrance"`    // 入口
	Method     string       `json:"method"`      // 分类
	// LimitType int64 `json:"limit_type"` //0：代表终生限制；1代表每日限制；2代表每周限制；3代表每月限制
	// LimitNum  int64 `json:"limit_num"`  //购买次数限制
}

type StoreListLittle struct {
	Buy       int64 `json:"buy"`         // 已购买数量
	GoodsUuid int64 `json:"goods_uuid"`  // 商品id
	Date      int32 `json:"date"`        // 备注信息
	OtherNum  int32 `json:"other_num"`   // 其他购买数量
	BuyAllNum int64 `json:"buy_all_num"` // 购买总数量
	// LimitType int64 `json:"limit_type"` //0：代表终生限制；1代表每日限制；2代表每周限制；3代表每月限制
	// LimitNum  int64 `json:"limit_num"`  //购买次数限制
}

type StoreExtractAccountRequest struct {
	Uid int64 `json:"uid"` // 用户id
}

type StoreExtractAccountResponse struct {
	Code   int    `json:"code"`   // 状态码
	Uid    int64  `json:"uid"`    // 用户id
	Type   int32  `json:"type"`   // 1-微信，2-支付宝
	Name   string `json:"name"`   // 姓名
	Detail string `json:"detail"` // 详情
	Data   string `json:"data"`   // 备注信息
}

type StoreExtractAccountAddRequest struct {
	Uid    int64  `json:"uid"`    // 用户id
	Type   int32  `json:"type"`   // 1-微信，2-支付宝
	Name   string `json:"name"`   // 姓名
	Detail string `json:"detail"` // 详情
}

type Attach struct {
	StoreId int64 `protobuf:"varint,1,opt,name=store_id,json=storeId,proto3" json:"store_id,omitempty"`
}
