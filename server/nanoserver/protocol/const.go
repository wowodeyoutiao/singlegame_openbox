package protocol

// 桌子类型普通/好友房
const (
	DeskModNomal  = 1
	DeskModFriend = 2
)

const (
	ExitTypeExitDeskUI           = -1
	ExitTypeDissolve             = 6
	ExitTypeSelfRequest          = 0
	ExitTypeClassicCoinNotEnough = 1
	ExitTypeDailyMatchEnd        = 2
	ExitTypeNotReadyForStart     = 3
	ExitTypeChangeDesk           = 4
	ExitTypeRepeatLogin          = 5
)

type GameType int

const (
	GameTypeMow GameType = iota + 100
)
