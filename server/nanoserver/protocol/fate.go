package protocol

type CreateFateRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type FateInfo struct {
}
type FateOperateRequest struct {
	Uid int64 `json:"uid"` //渠道列表
}
type FateOperateResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //渠道列表
}
type FateUpRequest struct {
	Uid      int64 `json:"uid"`       //渠道列表
	ConfigId int32 `json:"config_id"` //渠道列表
	Num      int32 `json:"num"`       //数量
}
type FateUpResponse struct {
	Code int    `json:"code"` //状态码
	Per  []int  `json:"per"`  //倍率
	Num  int32  `json:"num"`  //实际生效数目
	Data string `json:"data"` //
}
type FateAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type FateListRequest struct {
	Uid int64 //Token
}

// FateListResponse 订单列表
type FateListResponse struct {
	Code      int   `json:"code"`       //状态码
	FateLevel int32 `json:"fate_level"` //
	Exp       int32 `json:"exp"`        //
}
