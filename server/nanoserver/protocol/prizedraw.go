package protocol

type PrizedrawOperateRequest struct {
	Uid        int64 //Token
	ActivityId int32 `json:"activity_id"`
	UUID       int64 `json:"uuid"`
	Buy        int64 `json:"vip"`
	Check      int32 `json:"check"`
}

type PrizedrawRequest struct {
	Uid    int64 //Token
	Module int32 `json:"module"` // 1-抽   2-10抽
	Type   int32 `json:"type"`   // 1-抽   2-10抽
	PType  int32 `json:"ptype"`  //1青铜 2白银 3黄金
}

// 伙伴抽奖
type PrizedrawResponse struct {
	Code  int            `json:"code"`  //状态码
	Data  string         `json:"data"`  //返回信息
	Award []PartnerAward `json:"award"` //
}

// 伙伴奖池
type PartnerAward struct {
	Kind int32 `json:"kind"` // 类型 1-碎片 2-整卡
	Id   int32 `json:"id"`   // 物品id
	Num  int32 `json:"num"`  // 数量
}
type PrizedrawListResponse struct {
	Code    int               `json:"code"` //状态码
	Data    string            `json:"data"` //返回信息
	CanPlay []PrizedrawDetail `json:"canplay"`
}
type PrizedrawDetail struct {
	Type    int32 `json:"type"`     // 类型 1-青铜白银黄金
	CanPlay int32 `json:"can_play"` // 剩余可抽奖次数
	NeedNum int32 `json:"need_num"` // 还剩多少抽
}
