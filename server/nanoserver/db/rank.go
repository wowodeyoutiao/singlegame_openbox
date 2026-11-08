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

type Rank struct {
	GoodsUuiD  int64 `xorm:"goods_uuid"`
	Num        int64 `xorm:"num"`
	CreateTime int64 `xorm:"created_at"`
}

func RankList() ([]model.RankList, error) {

	RankModel := &model.RankList{}
	RankInfo := []model.RankList{}
	//now := time.Now()
	// 获得当前时间的年份、月份、日
	//year, month, day := now.Date()
	//// 创建一个时间，代表当天的开始（00:00:00）
	//midnight := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	err := database.Find(&RankInfo, RankModel)

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return RankInfo, nil
}
func RankUpdate(rtype int32, List []byte) (bool, error) {
	RankModel := &model.RankList{}
	RankInfo := []model.RankList{}
	// 获取当前时间
	now := time.Now()
	//// 获得当前时间的年份、月份、日
	//year, month, day := now.Date()
	//// 创建一个时间，代表当天的开始（00:00:00）
	//midnight := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	//result := make([]model.Attachment, 0)
	err := database.Where(" type = ? ", rtype).Find(&RankInfo, RankModel)
	if err != nil {
		return false, err
	}
	if len(RankInfo) == 0 {

		// 返回当前时间戳
		r := &model.RankList{
			Type:      rtype,
			List:      List,
			CreatedAt: now.Unix(),
			Status:    StatusNormal,
		}
		_, err1 := database.Insert(r)
		if err1 != nil {
			return false, err1
		}
	} else {
		r := &model.RankList{
			List:      List,
			UpdatedAt: now.Unix(),
			Status:    StatusNormal,
		}
		_, err1 := database.Where("type = ?  ", rtype).Update(r)
		if err1 != nil {
			return false, err1
		}
	}
	return true, nil
}
func RankUpdateSend(rtype int32) (bool, error) {
	// 获取当前时间
	now := time.Now()
	//// 获得当前时间的年份、月份、日
	//year, month, day := now.Date()
	//// 创建一个时间，代表当天的开始（00:00:00）
	//midnight := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	//result := make([]model.Attachment, 0)

	r := &model.RankList{
		SendAt:     now.Unix(),
		SendStatus: StatusNormal,
	}
	_, err1 := database.Where("type = ?  ", rtype).Update(r)
	if err1 != nil {
		return false, err1
	}
	return true, nil
}
