package db

import (
	"github.com/lonng/nano/serialize/protobuf"
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"time"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

type Grade struct {
	GoodsUuiD  int64 `xorm:"goods_uuid"`
	Num        int64 `xorm:"num"`
	CreateTime int64 `xorm:"create_time"`
}

// 用户购买情况
func GradeList(uid int64, dayType int32, goodsIds []int64) ([]Grade, error) {

	grade := []Grade{}
	//err := database.Where("uid = ? and gradetype =?  and created_at> ? and goods_id in(?)",
	//	uid, gradeType, dayType*3600*24, goodsIds).
	//	Desc("created_at").Limit(1).Find(&result, Sign)
	timea := time.Now().Unix() - int64(dayType*3600*24)
	err := database.Select("goods_uuid,count(*) as  num,min(created_at) create_time").Where("status = ? and uid = ?  and created_at>? ",
		StatusNormal, uid, timea).In("goods_uuid", goodsIds).GroupBy("goods_uuid").
		OrderBy("create_time").Find(&grade)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return grade, nil
}
func GradeAdd(uid int64, goodsUuid int32, SourceInfo []protocol.SourceInfo) (bool, error) {
	data, err := protobuf.NewSerializer().Marshal(SourceInfo)
	dataDel := ""
	if err == nil {
		dataDel = string(data)
	}

	u := &model.Grade{
		Uid:       uid,
		GoodsId:   int64(goodsUuid % 1000),
		GoodsUuid: goodsUuid,
		PropsInfo: dataDel,
		CreatedAt: time.Now().Unix(),
		Status:    StatusNormal,
	}
	_, err = database.Insert(u)
	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}
