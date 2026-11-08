package db

import (
	"fmt"

	"nanoserver/db/model"

	//"nanoserver/internal/web"
	"nanoserver/pkg/errutil"

	"github.com/go-xorm/xorm"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

func QueryLottery(uid int64) (*model.Lottery, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	ulottery := &model.Lottery{
		Uid: uid,
	}
	has, err := database.Get(ulottery)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Info(err)
		return nil, err
	}
	return ulottery, nil
}

// UpdateLottery 更新抽奖信息的函数
// 参数:
//
//	ulottery: 指向model.Lottery结构体的指针，包含要更新的抽奖信息
//
// 返回值:
//
//	error: 操作过程中遇到的错误，如果更新成功则返回nil
func UpdateUserLottery(ulottery *model.Lottery) error {
	// 使用database的Where方法筛选出指定uid的记录，然后更新为ulottery中的内容
	// 返回受影响的行数和可能的错误
	_, err := database.Where("uid=?", ulottery.Uid).Cols("refresh_time, prize_pool, have_prize_pool, progress, need_item_id, free_num").Update(ulottery)
	// 将操作结果返回给调用者
	return err
}

func InsertLottery(ulottery *model.Lottery) error {
	_, err := database.Insert(ulottery)
	return err
}

func InsertLotteryV2(session *xorm.Session, ulottery *model.Lottery) (affected int64, err error) {
	if ulottery == nil {
		return 0, fmt.Errorf("InsertLotteryV2: user lottery is nil")
	}
	return session.Insert(ulottery)
}
