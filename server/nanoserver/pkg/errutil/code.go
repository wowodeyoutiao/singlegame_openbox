package errutil

const (
	codeBase = 1000
)

const (
	Unknown = codeBase + iota
	yxBadRoute
	yxNotFound
	yxWrongType
	yxUserNameExists
	yxAuthFailed
	yxIllegalParameter
	yxNotRSAPublicKey
	yxNotRSAPrivateKey
	yxIllegalLoginType
	yxWrongPassword
	yxUserNameNotFound
	yxInitFailed
	yxServerInternal
	yxDbOperation
	yxCacheOperation
	yxPermissionDenied
	yxNotImplemented

	yxUserNotFound
	yxTokenNotFound
	yxInvalidToken
	yxIllegalName
	yxTokenMismatchUser
	yxInvalidPlatform
	yxUnsupportSignType
	yxOrderNotFound
	yxInvalidParameter
	yxRequestFailed
	yxIllegalOrderType
	yxCoinNotEnough
	yxFrequencyLimited
	yxVerifyFailed
	yxSignFailed
	yxTradeExisted
	yxProviderNotFound
	yxThirdAccountNotFound
	yxWrongThirdLoginType
	yxDirNotExists
	yxPayTestDisable
	yxPropertyNotFound
	yxUuidNotFound
	yxProductionNotFound
	yxRequestPrePayIDFailed
	yxDeskNotFound
	yxAlreadyInRoom
)

var errs = map[error]int{
	ErrBadRoute:              yxBadRoute,
	ErrNotFound:              yxNotFound,
	ErrWrongType:             yxWrongType,
	ErrUserNameExists:        yxUserNameExists,
	ErrIllegalParameter:      yxIllegalParameter,
	ErrNotRSAPublicKey:       yxNotRSAPublicKey,
	ErrNotRSAPrivateKey:      yxNotRSAPrivateKey,
	ErrAuthFailed:            yxAuthFailed,
	ErrIllegalLoginType:      yxIllegalLoginType,
	ErrWrongPassword:         yxWrongPassword,
	ErrUserNameNotFound:      yxUserNameNotFound,
	ErrInitFailed:            yxInitFailed,
	ErrServerInternal:        yxServerInternal,
	ErrDBOperation:           yxDbOperation,
	ErrCacheOperation:        yxCacheOperation,
	ErrPermissionDenied:      yxPermissionDenied,
	ErrNotImplemented:        yxNotImplemented,
	ErrUserNotFound:          yxUserNotFound,
	ErrTokenNotFound:         yxTokenNotFound,
	ErrInvalidToken:          yxInvalidToken,
	ErrIllegalName:           yxIllegalName,
	ErrTokenMismatchUser:     yxTokenMismatchUser,
	ErrInvalidPayPlatform:    yxInvalidPlatform,
	ErrUnsupportSignType:     yxUnsupportSignType,
	ErrOrderNotFound:         yxOrderNotFound,
	ErrInvalidParameter:      yxInvalidParameter,
	ErrRequestFailed:         yxRequestFailed,
	ErrIllegalOrderType:      yxIllegalOrderType,
	ErrCoinNotEnough:         yxCoinNotEnough,
	ErrFrequencyLimited:      yxFrequencyLimited,
	ErrVerifyFailed:          yxVerifyFailed,
	ErrSignFailed:            yxSignFailed,
	ErrTradeExisted:          yxTradeExisted,
	ErrProviderNotFound:      yxProviderNotFound,
	ErrThirdAccountNotFound:  yxThirdAccountNotFound,
	ErrWrongThirdLoginType:   yxWrongThirdLoginType,
	ErrDirNotExists:          yxDirNotExists,
	ErrPayTestDisable:        yxPayTestDisable,
	ErrPropertyNotFound:      yxPropertyNotFound,
	ErrUUIDNotFound:          yxUuidNotFound,
	ErrProductionNotFound:    yxProductionNotFound,
	ErrRequestPrePayIDFailed: yxRequestPrePayIDFailed,
	ErrDeskNotFound:          yxDeskNotFound,
	ErrAlreadyInRoom:         yxAlreadyInRoom,
}

const (
	ServerIsBusy        = "服务器繁忙"
	MsgKickTimeOver     = "副本倒计时结束，已自动退出！"
	MsgKickDieTimeOver  = "死亡复活倒计时结束，已自动退出！"
	MsgSafeArea         = "安全区内不能攻击！"
	EventIsPending      = "活动未开始"
	LimitHasBeenReached = "今日次数已用完"
	MsgBossKilled       = "玩家%s成功击杀了boss %s，获得稀有掉落！"
	MsgSoloTower        = "玩家%s成功通过了无尽试炼%d层！"
	VIPNotEnough        = "VIP等级不足"
	EnergyMax           = "体力值已达上限"
	IllegalName         = "不合法身份证姓名"
	IllegalIDCard       = "不合法身份证号码"
	IllegalAge          = "根据相关规定，您当前未满18周岁，无法继续进行游戏"
	NoTimes             = "今日次数已用完"
	IllegalInBattle     = "战斗中，无法进行该操作"
	IllegalSkillType    = "技能类型不合法"
	EquipOnUse          = "物品正在使用中"
	ModClosed           = "功能已关闭"
	TargetDead          = "目标已死亡"
	TargetNotFound      = "目标未找到"
	TargetInBattle      = "当前玩家正在战斗中"
	EnergyNotEnough     = "没有体力了，请休息一下再战斗吧！"
	CanotMining         = "无法采集"
	RoomPlayerMax       = "房间人数已满"
	PrivilegeNotEnough  = "特权不足"
	WarKilled           = "你已被玩家<color=#ff5855>%s</color>击杀！"
	KillMonControl      = "已达今日可采集上限"
	ErrModuleNotOpen    = "功能未开放"
)
