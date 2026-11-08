package db

import (
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
	"time"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

// 用户最新日签到情况
func TurntableListDay(uid int64) ([]model.Turntable, error) {
	Turntable := &model.Turntable{}

	result := make([]model.Turntable, 0)
	err := database.Where("uid = ? and type =? and status=? ", uid, 1, StatusNormal).
		Desc("created_at").Limit(1).Find(&result, Turntable)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 用户最新累计签到情况
func TurntableList(uid int64, content int32, tag int32) ([]model.Turntable, error) {
	Turntable := &model.Turntable{}

	result := make([]model.Turntable, 0)
	err := database.Where("Uid = ? and tag =? and  content = ? and status=?", uid, tag, content, StatusNormal).
		Desc("created_at").Find(&result, Turntable)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

func TurntableAdd(uid int64, content int32, tag int32, coupon int32) (bool, error) {
	u := &model.Turntable{
		Uid:       uid,
		Tag:       tag,
		Content:   content,
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
