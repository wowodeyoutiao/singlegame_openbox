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

// 用户最新情况
func PrizedrawList(uid int64, ptype int32, module int32, date int) (float64, error) {
	// now := time.Now()
	// // 创建一个新的时间，年份、月份、日与当前时间相同，其余为0
	// midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// // 转换为时间戳（秒）
	// nowTime := midnight.Unix()
	nowTime := int64(0)
	if date > 0 {
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// 转换为时间戳（秒）
		nowTime = midnight.Unix()
	}
	// result := make([]model.Partner, 0)
	// result, err := database.Query("SELECT SUM(num) AS nums,Uid FROM treasure_box WHERE Uid = ? and  created_at > ?  ", uid, nowTime)
	result := &model.Summon{}
	Prizedraw, err := database.Where("Uid = ? and created_at>? and ptype=?  and module=?", uid, nowTime, ptype, module).Sum(*result, "num")
	if err != nil {
		return 0, errutil.ErrDBOperation
	}
	return Prizedraw, nil
}

// 抽奖插入插入
// PrizedrawAdd 添加抽奖记录
// uid: 用户ID
// sourceInfo: 来源信息
// num: 数量
// ptype: 奖品类型
// module: 模块类型
// 返回: 错误信息
func PrizedrawAdd(uid int64, sourceInfo string, num int32, ptype int32, module int32) error {
	// atype := uuid / 1000

	u := &model.Summon{
		Uid:        uid,
		Module:     module,
		Ptype:      ptype,
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
