package protocol

type CreateSignRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type SignInfo struct {
}
type SignOperateRequest struct {
	Uid         int64 //
	AccountName string
}
type SignOperateResponse struct {
	Code int  `json:"code"` //状态码
	Data bool `json:"data"` //渠道列表
}
type SignAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type SignListRequest struct {
	Uid int64 //Token
}

// SignListResponse 订单列表
type SignListResponse struct {
	Code int    `json:"code"` //状态码
	Data Awsard `json:"data"` //
}
type Awsard struct {
	DayAward        DayAward          `json:"day_award"`        //每日签到
	AccumulateAward []AccumulateAward `json:"accumulate_award"` //累计签到提取情况
}
type DayAward struct {
	Round int32 `json:"round"` //用户当前轮数
	Day   int32 `json:"day"`   //用户签到第几天
	Sign  bool  `json:"sign"`  //
}

type AccumulateAward struct {
	Round int32 `json:"round"` //提取的第几轮
	Day   int32 `json:"day"`   //提取的第几天
}
type AccumulateAwardRequest struct {
	Uid         int64  `json:"uid"`          //Token
	AccountName string `json:"account_name"` //
	Round       int32  `json:"round"`        //提取的第几轮
	Day         int32  `json:"day"`          //提取的第几天
}
type UserProps struct {
	Props []Props `json:"props"` //商品信息
}
type Props struct {
	PropsId int64 `json:"props_id"` //商品
	Num     int64 `json:"num"`
}
