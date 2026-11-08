package protocol

import "time"

type CreatePrivilegeRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type PrivilegeInfo struct {
}
type PrivilegeOperateRequest struct {
	Uid         int64 //Token
	PrivilegeId int32 `json:"privilege_id"`
	PID         int32 `json:"pid"`
	Buy         int64 `json:"vip"`
}
type PrivilegeOperateResponse struct {
	Code int  `json:"code"` //状态码
	Data bool `json:"data"` //渠道列表
}
type PrivilegeAdd struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type PrivilegeListRequest struct {
	Uid int64 //Token
}
type PrivilegeAllResponse struct {
	Code int            `json:"code"` //状态码
	Data []PrivilegeAll `json:"data"` //详情
}
type PrivilegeAll struct {
	Id        int32 `json:"id"`         //活动id
	ActType   int32 `json:"act_type"`   //活动类型
	StartTime int64 `json:"start_time"` //开始时间
	EndTime   int64 `json:"end_time"`   //开始时间
	ResetTime int32 `json:"reset_time"` //重置时间
}

// PrivilegeListResponse 订单列表
type PrivilegeListResponse struct {
	Code int                  `json:"code"` //状态码
	Data []PrivilegeDetailDay `json:"data"` //详情
}
type PrivilegeReceiveDayResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //
}

//	type PrivilegeListCharge struct {
//		Day    int32 `json:"day"`    //第几天
//		Status int   `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）
//	}
type PrivilegeDetailDay struct {
	Pid     int32 `json:"pid"`      //内容编号
	Buy     int   `json:"buy"`      //1-未购买，2-已购买
	Get     int   `json:"get"`      //状态 1-待领取，2-已领取，
	EndTime int64 `json:"end_time"` //到期时间
}
type PrivilegeDetailTCP struct {
	Pid     int32     `json:"pid"`      //内容编号
	EndTime time.Time `json:"end_time"` //到期时间
}

// PrivilegeListResponse 订单列表
type PrivilegeChargeResponse struct {
	Code   int    `json:"code"`   //状态码
	Day    int64  `json:"day"`    //当前天数
	Data   string `json:"data"`   //返回信息
	Status int    `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）
}
