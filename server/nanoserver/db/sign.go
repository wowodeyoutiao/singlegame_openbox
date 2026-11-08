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

// 签到
func SignAddDay(uid int64, vipLevel int32, userRound int32, round int32, day int32, SourceInfo string) (bool, error) {
	u := &model.Sign{
		Uid:        uid,
		VipLevel:   vipLevel,
		Round:      round,
		UserRound:  userRound,
		Type:       1,
		CreatedAt:  time.Now().Unix(),
		Day:        day,
		Status:     StatusNormal,
		SourceInfo: SourceInfo,
	}
	_, err := database.Insert(u)

	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}

// 用户最新日签到情况
func SignListDay(uid int64) ([]model.Sign, error) {
	Sign := &model.Sign{}

	result := make([]model.Sign, 0)
	err := database.Where("uid = ? and type =? and status=? ", uid, 1, StatusNormal).
		Desc("created_at").Limit(1).Find(&result, Sign)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 用户最新累计签到情况
func SignListAccumulate(uid int64, userRound int32, round int32) ([]model.Sign, error) {
	Sign := &model.Sign{}

	result := make([]model.Sign, 0)
	err := database.Where("Uid = ? and type =? and user_round =? and round = ? and status=?", uid, 2, userRound, round, StatusNormal).
		Desc("created_at").Find(&result, Sign)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

func SignExtract(uid int64, vipLevel int32, userRound int32, round int32, day int32, SourceInfo string) (bool, error) {
	u := &model.Sign{
		Uid:        uid,
		VipLevel:   vipLevel,
		Round:      round,
		UserRound:  userRound,
		Type:       2,
		CreatedAt:  time.Now().Unix(),
		Day:        day,
		Status:     StatusNormal,
		SourceInfo: SourceInfo,
	}
	_, err := database.Insert(u)

	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}
