package db

import (
	"time"

	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

// 插入
func ActivityAddDay(uid int64, uuid int64, buy int, actId int32) (bool, error) {
	atype := uuid / 1000
	act := actId
	day := uuid - int64(atype*1000) - int64(actId*100)

	u := &model.Activity{
		Uid:       uid,
		Uuid:      uuid,
		Tag:       int32(act),
		Type:      int(atype),
		CreatedAt: time.Now().Unix(),
		Day:       int32(day),
		Buy:       buy,
		Status:    StatusNormal,
	}
	_, err := database.Insert(u)
	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}

func ActivityAddDayUpdate(uid int64, uuid int64, atype int, buy int) (bool, error) {
	u := &model.Activity{
		UpdatedAt: time.Now().Unix(),

		Buy: buy,
	}

	if _, err := database.Where("uid=? and uuid=?", uid, uuid).Update(u); err != nil {
		return false, err
	}

	return true, nil
}

// 用户领取情况
func ActivityListDay(uid int64, atype int32, activity int32) ([]model.Activity, error) {
	Activity := &model.Activity{}

	result := make([]model.Activity, 0)
	err := database.Where("uid = ? and type =? and  tag=? and status =?", uid, atype, activity, StatusNormal).
		Desc("created_at").Find(&result, Activity)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 用户最新累计签到情况
func ActivityListAccumulate(uid int64, userRound int32, round int32) ([]model.Activity, error) {
	Activity := &model.Activity{}

	result := make([]model.Activity, 0)
	err := database.Where("Uid = ? and type =? and user_round =? and round = ? ", uid, 2, userRound, round).
		Desc("created_at").Find(&result, Activity)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 用户最新大转盘情况
func ActivityListTurntable(uid int64, tag int32) ([]model.Turntable, error) {
	Activity := &model.Turntable{}

	result := make([]model.Turntable, 0)
	err := database.Where("Uid = ? and Tag =?  ", uid, tag).
		Desc("created_at").Find(&result, Activity)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 大转盘插入
func ActivityAddTurntable(uid int64, uuid int64, coupon int32) (bool, error) {
	day := uuid % 100
	// atype := uuid / 1000
	act := (uuid % 1000) / 100
	u := &model.Turntable{
		Uid:       uid,
		Tag:       int32(act),
		Content:   int32(day),
		Coupon:    coupon,
		CreatedAt: time.Now().Unix(),
		Status:    StatusNormal,
	}
	_, err := database.Insert(u)
	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}

// 用户最新百宝箱情况
func ActivityListTreasureBox(uid int64) (float64, error) {
	now := time.Now()

	// 创建一个新的时间，年份、月份、日与当前时间相同，其余为0
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 转换为时间戳（秒）
	nowTime := midnight.Unix()
	// result := make([]model.TreasureBox, 0)
	// result, err := database.Query("SELECT SUM(num) AS nums,Uid FROM treasure_box WHERE Uid = ? and  created_at > ?  ", uid, nowTime)
	result := &model.TreasureBox{}
	Activity, err := database.Where("Uid = ? and created_at>?  ", uid, nowTime).Sum(*result, "num")
	if err != nil {
		return 0, errutil.ErrDBOperation
	}
	return Activity, nil
}

// 百宝箱盘插入
func ActivityAddTreasureBox(uid int64, sourceInfo string, num int32) error {
	// atype := uuid / 1000

	u := &model.TreasureBox{
		Uid:        uid,
		SourceInfo: sourceInfo,
		CreatedAt:  time.Now().Unix(),
		Status:     StatusNormal,
		Num:        num,
	}
	_, err := database.Insert(u)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// ActivityRecordUpsert is a placeholder function for upserting activity records.
func ActivityRecordUpsert(ar *model.ActivityRecord) (int64, error) {
	// 获取主键
	id := database.IDOf(ar)
	if len(id) == 0 {
		// 没有主键，直接插入
		return database.Insert(ar)
	}

	// 先尝试更新
	affected, err := database.ID(id[0]).Update(ar)
	if err != nil {
		return 0, err
	}

	// 如果没有更新任何行，则插入
	if affected == 0 {
		return database.Insert(ar)
	}

	return affected, nil
}

// ActivityRecordGet retrieves the activity record from the database.
func ActivityRecordGet(PID int64) (*model.ActivityRecord, error) {
	var ar model.ActivityRecord
	has, err := database.Where("uid = ?", PID).Get(&ar)
	if err != nil {
		return nil, errutil.ErrDBOperation
	}
	if !has {
		return nil, errutil.ErrNotFound
	}

	return &ar, nil
}
