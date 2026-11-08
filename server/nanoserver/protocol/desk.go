package protocol

import (
	"nanoserver/pkg/constant"
)

type EnterDeskInfo struct {
	DeskPos  int    `json:"deskPos"`
	Uid      int64  `json:"acId"`
	Nickname string `json:"nickname"`
	IsReady  bool   `json:"isReady"`
	Sex      int    `json:"sex"`
	IsExit   bool   `json:"isExit"`
	HeadUrl  string `json:"headURL"`
	Score    int    `json:"score"`
	IP       string `json:"ip"`
	Offline  bool   `json:"offline"`
}

type ExitResponse struct {
	AccountId int64 `json:"acid"`
	IsExit    bool  `json:"isexit"`
	ExitType  int   `json:"exitType"`
	DeskPos   int   `json:"deskPos"`
}

type PlayerEnterDesk struct {
	Data []EnterDeskInfo `json:"data"`
}

type ExitRequest struct {
	IsDestroy bool `json:"isDestroy"`
}

type DeskBasicInfo struct {
	DeskID string `json:"deskId"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Mode   int    `json:"mode"`
}

type ScoreInfo struct {
	Uid   int64 `json:"acId"`
	Score int   `json:"score"`
}

type GameEndScoreChange struct {
	Uid    int64 `json:"acId"`
	Score  int   `json:"score"`
	Remain int   `json:"remain"`
}

type RoundOverStats struct {
	Title       string               `json:"title"`
	Round       string               `json:"round"`
	ScoreChange []GameEndScoreChange `json:"scoreChange"`
}

type DeskPlayerData struct {
	Uid int64 `json:"acId"`
}

type SyncDesk struct {
	Status    constant.DeskStatus `json:"status"` //1,2,3,4,5
	Players   []DeskPlayerData    `json:"players"`
	ScoreInfo []ScoreInfo         `json:"scoreInfo"`
}

type ReConnect struct {
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
	HeadUrl string `json:"headUrl"`
	Sex     int    `json:"sex"`
}

type Desk struct {
	Id           int64  `json:"id"`
	Creator      int64  `json:"creator"`
	Round        int    `json:"round"`
	DeskNo       string `json:"desk_no"`
	Mode         int    `json:"mode"`
	CreatedAt    int64  `json:"created_at"`
	CreatedAtStr string `json:"created_at_str"`
	DismissAt    int64  `json:"dismiss_at"`
	Extras       string `json:"extras"`
}

type DestroyDeskResponse struct {
	RoundStats       *RoundOverStats `json:"roundStats"`
	Title            string          `json:"title"`
	IsNormalFinished bool            `json:"isNormalFinished"`
}

type DeskListResponse struct {
	Code  int    `json:"code"`
	Total int64  `json:"total"` //总数量
	Data  []Desk `json:"data"`
}

type DeskByIDResponse struct {
	Code int   `json:"code"`
	Data *Desk `json:"data"`
}
