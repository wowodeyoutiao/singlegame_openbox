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
func PrivilegeAddDay(uid int64, pid int32) (bool, error) {
	u := &model.Privilege{
		Uid:       uid,
		Type:      int(pid),
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

// 用户最新签到情况
func PrivilegeListDay(uid int64) ([]model.Privilege, error) {
	Privilege := &model.Privilege{}
	now := time.Now()

	// 设置时间为当天凌晨
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 获取凌晨的时间戳（Unix时间戳）
	timestamp := midnight.Unix()
	result := make([]model.Privilege, 0)
	err := database.Where("uid = ? and type in(1,2,3) and status=?  and created_at>? ", uid, StatusNormal, timestamp).
		Desc("created_at").Find(&result, Privilege)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}
