package protocol

type CreateRemainRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type RemainInfo struct {
}
type RemainOperateRequest struct {
	Uid int64 `json:"uid"` //渠道列表
}
type RemainOperateResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //渠道列表
}
type RemainUpRequest struct {
	Uid      int64 `json:"uid"`       //渠道列表
	ConfigId int32 `json:"config_id"` //渠道列表
	Num      int32 `json:"num"`       //数量
}
type RemainUpResponse struct {
	Code int    `json:"code"` //状态码
	Per  []int  `json:"per"`  //倍率
	Num  int32  `json:"num"`  //实际生效数目
	Data string `json:"data"` //
}
type RemainAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type RemainListRequest struct {
	Uid int64 //Token
}

// RemainListResponse 订单列表
type RemainListResponse struct {
	Code        int   `json:"code"`         //状态码
	RemainLevel int32 `json:"remain_level"` //
	Exp         int32 `json:"exp"`          //
}
