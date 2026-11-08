package protocol

type WindowLtvListRequest struct {
	PageIndex int  `json:"page_index"`
	PageSize  int  `json:"page_size"`
	Export    bool `json:"export"`
}
type WindowLtvResponse struct {
	Code int           `json:"code"` //状态码
	Data WindowLtvList `json:"data"` //备注信息
}

type WindowLtvList struct {
	Count     int               `json:"count"`      // 总数
	PageIndex int               `json:"page_index"` // 次日
	PageSize  int               `json:"page_size"`  //
	Detail    []WindowLtvDetail `json:"detail"`     // 总数
}

// Ltv
type WindowLtvDetail struct {
	Date   string `json:"date"`   // 日期
	Today  int    `json:"today"`  // 当天
	Two    int    `json:"two"`    // 次日
	Three  int    `json:"three"`  // 3日
	Seven  int    `json:"seven"`  // 7日
	Thirty int    `json:"thirty"` // 7日
}
type WindowShushuRequest struct {
	Name string `json:"name"`
}
type WindowShushuRegisterRequest struct {
	Ip       string `json:"ip"`
	Name     string `json:"name"`
	Register int64  `json:"register"`
	Channel  string `json:"channel"`
}
type WindowShushuHumpHideRequest struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}
type WindowShushuCompletedRequest struct {
	Name       string `json:"name"`
	Guide_step int    `json:"guide_step"`
	Guide_name string `json:"guide_name"`
}
type WindowShusOpenBoxRequest struct {
	Name              string        `json:"name"`
	BoxItemInfo       []BoxItemInfo `json:"box_item_info"`
	IsAuto            bool          `json:"is_auto"`
	PowerUp           int           `json:"power_up"`
	OpenNum           int           `json:"open_num"`
	RolePowerStop     bool          `json:"role_power_stop"`
	PartnerPowerStop  bool          `json:"partner_power_stop"`
	BoxNumBefore      int           `json:"box_num_before"`
	BoxNumAfter       int           `json:"box_num_after"`
	BoxJobLimited     string        `json:"box_job_limited"`
	BoxLevelLimited   int           `json:"box_level_limited"`
	RetainProperty    string        `json:"retain_property"`
	RetainPropertyNum string        `json:"retain_property_num"`
}
type WindowShushuAdShowRequest struct {
	Name       string `json:"name"`
	AdID       string `json:"ad_id"`
	AdType     string `json:"ad_type"`
	AdScene    string `json:"ad_scene"`
	AdPosition string `json:"ad_position"`
}

type WindowShushuAdCloseRequest struct {
	Name         string        `json:"name"`
	AdID         string        `json:"ad_id"`
	AdType       string        `json:"ad_type"`
	AdScene      string        `json:"ad_scene"`
	AdTime       int32         `json:"ad_time"`
	AdDuration   int32         `json:"ad_duration"`
	ViewComplete bool          `json:"view_complete"`
	RewardType   string        `json:"reward_type"`
	AdPosition   string        `json:"ad_position"`
	GetItemInfo  []GetItemInfo `json:"get_item_info"`
}
type WindowShushuClickEventRequest struct {
	Name       string `json:"name"`
	Page       string `json:"page"`
	Component  string `json:"component"`
	Position   string `json:"position"`
	EventClick string `json:"event_click"`
}
type GetItemInfo struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	ItemNum  int    `json:"item_num"`
}
type BoxItemInfo struct {
	ItemID       string `json:"item_id"`
	ItemLevel    int    `json:"item_level"`
	ItemQuality  int    `json:"item_quality"`
	ItemName     string `json:"item_name"`
	Count        int    `json:"count"`
	ItemProperty string `json:"item_property"` //item_property
}
type WindowShushuResponse struct {
	Code int `json:"code"` //状态码
}

type QueryUserInfoResponse struct {
	Code int       `json:"code"` //状态码
	Data *UserInfo `json:"data"` //用户信息
}
type SendLampRequest struct {
	Title  string `json:"title"`  // 灯标题
	Stype  int32  `json:"stype"`  // 灯类型
	Svalue int64  `json:"svalue"` // 灯对象
}

type SendMailRequest struct {
	Name       string       `json:"name"`       // 发送者名称
	Title      string       `json:"title"`      // 邮件标题
	Desc       string       `json:"desc"`       // 邮件内容
	Attachment []Attachment `json:"attachment"` // 附件
	Stype      int32        `json:"stype"`      // 邮件类型
	SValue     int64        `json:"svalue"`     // 邮件对象
	Key        string       `json:"key"`
}

type SendMuteRequest struct {
	Uid      int64  `json:"uid"`      // 玩家UID
	Duration int64  `json:"duration"` // 禁言时长（秒），0表示解除禁言
	Key      string `json:"key"`      // 加密验证key
}

type BanUserRequest struct {
	Uid      int64  `json:"uid"`      // 玩家UID
	Duration int64  `json:"duration"` // 封号时长（秒），0表示解除封号
	Key      string `json:"key"`      // 加密验证key
}

type DeleteUserRequest struct {
	Uid    int64  `json:"uid"`    // 玩家UID
	Action int32  `json:"action"` // 操作类型：1=删除，2=恢复
	Key    string `json:"key"`    // 加密验证key
}

type QueryUserInfoRequest struct {
	Uid int64  `json:"uid"` // 玩家UID
	Key string `json:"key"` // 加密验证key
}