package protocol

type CreateActivityRequest struct {
	RegistrationId string //来自哪个渠道的订单
	Uid            int64  //Token
	Offset         int    `json:"offset"`
	Count          int    `json:"count"`
}
type ActivityInfo struct {
}
type ActivityOperateRequest struct {
	Uid        int64 //Token
	ActivityId int32 `json:"activity_id"`
	UUID       int64 `json:"uuid"`
	Buy        int64 `json:"vip"`
	Check      int32 `json:"check"`
	Ad         int   `json:"ad"`
	Atype      int32 `json:"atype"`
}
type ActivityOperateResponse struct {
	Code int  `json:"code"` //状态码
	Data bool `json:"data"` //渠道列表
}
type ActivityAddRequest struct {
	PropsId int64 `json:"props_id"` //渠道列表
	Num     int64 `json:"num"`
}
type ActivityListRequest struct {
	Uid int64 //Token
}
type ActivityAllResponse struct {
	Code int           `json:"code"` //状态码
	Data []ActivityAll `json:"data"` //详情
}
type ActivityAll struct {
	Id        int32 `json:"id"`         //活动id
	ActType   int32 `json:"act_type"`   //活动类型
	StartTime int64 `json:"start_time"` //开始时间
	EndTime   int64 `json:"end_time"`   //开始时间
	ResetTime int32 `json:"reset_time"` //重置时间
}

// ActivityListResponse 订单列表
type ActivityListResponse struct {
	Code int               `json:"code"` //状态码
	Day  int64             `json:"day"`  //当前天数
	Buy  int64             `json:"buy"`  //
	Data []ActivityListDay `json:"data"` //渠道列表
	Num  int64             `json:"num"`  //数量
}
type ActivityReceiveDayResponse struct {
	Code int    `json:"code"` //状态码
	Data string `json:"data"` //
}
type ActivityReceiveTurntableResponse struct {
	Code int    `json:"code"` //状态码
	Puid int32  `json:"puid"` //中奖编号
	Data string `json:"data"` //
}

type ActivityListDay struct {
	Day    int32 `json:"day"`    //第几天
	Buy    int   `json:"buy"`    //付费状态 1-待领取，2-已领取，3-不可领取
	Status int   `json:"status"` //状态 1-待领取，2-已领取，3-不可领取，4-需要补签
}

//	type ActivityListCharge struct {
//		Day    int32 `json:"day"`    //第几天
//		Status int   `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）
//	}
type ActivityDetailDay struct {
	UUID                int64        `json:"uuid"`                 //内容编号
	Content             int32        `json:"content"`              //内容编号
	CommonItem          []SourceInfo `json:"common_item"`          //普通道具
	PayItem             []SourceInfo `json:"pay_item"`             //普通道具
	PreLimit            int32        `json:"pre_limit"`            //前置条件
	CompensationConsume []SourceInfo `json:"compensation_consume"` // 补领消耗
	Status              int          `json:"status"`               //状态 1-待领取，2-已领取，3-不可领取，4-需要补签
}

// ActivityListResponse 订单列表
type ActivityChargeResponse struct {
	Code   int               `json:"code"`   //状态码
	Day    int64             `json:"day"`    //当前天数
	Data   []ActivityListDay `json:"data"`   //所有首充情况列表
	Status int               `json:"status"` //状态 1-待领取，2-已领取，3-不可领取（没充值）（保留兼容性）
}

// ActivityGradeResponse 等级奖励
type ActivityGradeResponse struct {
	Code              int                 `json:"code"`                //状态码
	Data              string              `json:"data"`                //返回信息
	ActivityListGrade []ActivityListGrade `json:"activity_list_grade"` //等级列表详情
	Num               int64               `json:"num"`                 //数量
}
type ActivityListGrade struct {
	Uuid   int32 `json:"uuid"`   //内容编号
	Grade  int32 `json:"grade"`  //多少等级
	Status int   `json:"status"` //状态 1-待领取，2-已领取，3-不可领取
}

// 大转盘抽奖
type ActivityTruntableResponse struct {
	Code    int   `json:"code"`    //状态码
	Content int32 `json:"content"` //当前准备抽的内容编号 -1就是已经没有可以抽的了
}

// 百宝箱抽奖
type ActivityTreasureBoxRequest struct {
	Uid  int64 //Token
	Type int32 `json:"type"` // 1-抽   2-10抽
}

// 百宝箱抽奖
type ActivityTreasureBoxResponse struct {
	Code     int     `json:"code"`     //状态码
	Data     string  `json:"data"`     //返回信息
	Position []int32 `json:"position"` //
}
type ActivityListTreasureBoxResponse struct {
	Code    int    `json:"code"`     //状态码
	Data    string `json:"data"`     //返回信息
	CanPlay int32  `json:"can_play"` // 剩余可抽奖次数
	NeedNum int32  `json:"need_num"` // 还剩多少抽

}
