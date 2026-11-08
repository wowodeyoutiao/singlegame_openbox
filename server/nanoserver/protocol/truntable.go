package protocol

type CreateTurntableRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type TurntableInfo struct {
}
type TurntableOperateRequest struct {
	Uid         int64 //
	AccountName string
}
type TurntableOperateResponse struct {
	Code int  `json:"code"` //状态码
	Data bool `json:"data"` //渠道列表
}
type TurntableAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type TurntableListRequest struct {
	Uid int64 //Token
}

// TurntableListResponse 订单列表
type TurntableListResponse struct {
	Code int    `json:"code"` //状态码
	Data Awsard `json:"data"` //
}
