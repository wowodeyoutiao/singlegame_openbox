package game

const (
	kCurPlayer = "player"
)

const (
	fieldDesk    = "desk"
	fieldPlayer  = "player"
	fieldMonster = "monster"
)

const BaseRate float64 = 10000.0 // 定义万分比的基数为 10000.0

type PlayerStatus int32 //玩家状态

const (
	PlayerStatusIdle    PlayerStatus = iota //空闲
	PlayerStatusPlaying                     //游戏中
	PlayerStatusDisconnect
	PlayerStatusOffLine
	PlayerStatusRobot //离线机器人
)

type PvpRecordType int32

const (
	PvpRecordTypeJiShi PvpRecordType = iota + 1 //计时赛
	PvpRecordTypeHuiHe                          //回合制
)

const (
	ActiveSkill  = iota + 1 // 主动技能，值为 1
	PassiveSkill            // 被动技能，值为 2
)

// 获得道具类型
const (
	AddSourceTypeNew      = iota + 1 // 新获得
	AddSourceTypeUnEquip             // 装备卸下
	AddSourceTypeAutoDis             // 自动分解
	AddSourceTypeUnNotify            // 获得不飘字
)

const MaxPartnerNum int32 = 4     // 伙伴数量
const TurnBasedMax int32 = 6      // 回合制上阵数量
const ElementPartnerMax int32 = 6 // 回合制上阵数量
const DefaultPvpScore = 1000
const LotteryPoolMax int = 4 // 抽奖卡池人数
