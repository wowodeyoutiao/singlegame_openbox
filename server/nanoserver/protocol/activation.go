package protocol

type CreateActivationRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type ActivationInfo struct {
}
type ActivationOperateRequest struct {
	Uid         int64 //
	AccountName string
}
type ActivationOperateResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //返回信息
}
type ActivationAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type ActivationListRequest struct {
	Uid            int64  `json:"uid"`             //用户id
	ActivationCode string `json:"activation_code"` //商城类别
}

// ActivationListResponse 订单列表
type ActivationListResponse struct {
	Code int              `json:"code"` //状态码
	Data []ActivationList `json:"data"` //渠道列表
}

type ActivationList struct {
	Buy        int64      `json:"buy"`        //已购买数量
	GoodsId    int64      `json:"goods_id"`   //商品id
	Num        int64      `json:"num"`        //总数量
	NeedMoney  int32      `json:"need_money"` //前置条件
	PreLimit   int32      `json:"pre_limit"`  //前置条件
	SourceInfo SourceInfo `json:"props"`      //商品信息
	//LimitType int64 `json:"limit_type"` //0：代表终生限制；1代表每日限制；2代表每周限制；3代表每月限制
	//LimitNum  int64 `json:"limit_num"`  //购买次数限制
}
