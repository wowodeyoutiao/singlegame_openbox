package protocol

import (
	"nanoserver/pkg/constant"
)

type ReJoinDeskRequest struct {
	DeskNo string `json:"deskId"`
}

type ReJoinDeskResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ReEnterDeskRequest struct {
	DeskNo string `json:"deskId"`
}

type TableInfo struct {
	DeskNo    string              `json:"deskId"`
	CreatedAt int64               `json:"createdAt"`
	Creator   int64               `json:"creator"`
	Title     string              `json:"title"`
	Desc      string              `json:"desc"`
	Status    constant.DeskStatus `json:"status"`
	Round     uint32              `json:"round"`
	Mode      int32               `json:"mode"`
}

type JoinDeskResponse struct {
	Code      int       `json:"code"`
	Error     string    `json:"error"`
	TableInfo TableInfo `json:"tableInfo"`
}

type CheckOrderReqeust struct {
	OrderID string `json:"orderid"`
}

type CheckOrderResponse struct {
	Code   int    `json:"code"`
	Error  string `json:"error"`
	FangKa int    `json:"fangka"`
}

type DissolveStatusItem struct {
	DeskPos int    `json:"deskPos"`
	Status  string `json:"status"`
}

type DissolveResponse struct {
	DissolveUid    int64                `json:"dissolveUid"`
	DissolveStatus []DissolveStatusItem `json:"dissolveStatus"`
	RestTime       int32                `json:"restTime"`
}

type DissolveStatusResponse struct {
	DissolveStatus []DissolveStatusItem `json:"dissolveStatus"`
	RestTime       int32                `json:"restTime"`
}

type DissolveResult struct {
	DeskPos int `json:"deskPos"`
}

type PlayerOfflineStatus struct {
	Uid     int64 `json:"uid"`
	Offline bool  `json:"offline"`
}

// 属性设置
type PlayerSetting struct {
	Job                int `json:"job"`           // 职业
	Level              int `json:"level"`         // 等级
	EquipLevel         int `json:"equip_level"`   // 装备等级
	Enhance            int `json:"enhance"`       // 装备的强化等级
	Stars              int `json:"stars"`         // 装备的星级
	SkillLevel         int `json:"skill_level"`   // 技能等级
	SkillAdvancedLevel int `json:"skill_advance"` // 技能进阶
	FateLevel          int `json:"fate_level"`    // 命格等级
	FateUpgrade        int `json:"fate_upgrade"`  // 命格等阶
	GemLevel           int `json:"gem_level"`     // 宝石等级
}

// 对战设置
type CombatSetting struct {
	TotalRound int            `json:"round"` // 战斗次数
	Red        *PlayerSetting `json:"red"`   // 红方
	Blue       *PlayerSetting `json:"blue"`  // 蓝方
}

// 进度查询
type CombatResult struct {
	CurRound int  `json:"cur_round"` // 当前战斗次数
	RedWin   int  `json:"red_win"`   // 红方胜
	BlueWin  int  `json:"blue_win"`  // 蓝方胜
	Finish   bool `json:"finish"`    // 是否结束
}
