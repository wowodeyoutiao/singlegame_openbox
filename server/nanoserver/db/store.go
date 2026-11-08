package db

import (
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"time"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

type Store struct {
	GoodsUuiD  int64 `xorm:"goods_uuid"`
	Num        int64 `xorm:"num"`
	CreateTime int64 `xorm:"create_time"`
}

// 用户购买情况
func StoreList(uid int64, storeType int32, dayType int32, goodsIds []int64, normal bool) ([]Store, error) {

	store := []Store{}
	//err := database.Where("uid = ? and storetype =?  and created_at> ? and goods_id in(?)",
	//	uid, storeType, dayType*3600*24, goodsIds).
	//	Desc("created_at").Limit(1).Find(&result, Sign)
	now := time.Now()
	var timea int64
	if storeType == constant.ShopTypeSubscription || storeType == constant.ShopTypeAFK {
		normal = true
	}
	if normal {
		tomorrow := now.AddDate(0, 0, 1)                                                                          // 增加一天
		midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()) // 设置为0点0分0秒
		timea = midnight.Unix() - int64((dayType)*3600*24)

	} else {
		switch dayType {
		case 1:
			tomorrow := now.AddDate(0, 0, 1)
			midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
			timea = midnight.Unix() - int64(dayType*3600*24)
		case 7: // 按周统计，跨周重置
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7 // 周日
			}
			monday := time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, now.Location())
			timea = monday.Unix()
		case 30: // 按月统计，跨月重置
			firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			timea = firstDayOfMonth.Unix()
		case 90: // 按季度统计，跨季度重置
			month := int(now.Month())
			var quarterStartMonth int
			switch {
			case month >= 1 && month <= 3:
				quarterStartMonth = 1
			case month >= 4 && month <= 6:
				quarterStartMonth = 4
			case month >= 7 && month <= 9:
				quarterStartMonth = 7
			default:
				quarterStartMonth = 10
			}
			firstDayOfQuarter := time.Date(now.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, now.Location())
			timea = firstDayOfQuarter.Unix()
		case 365: // 按年统计，跨年重置
			firstDay := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			timea = firstDay.Unix()
		default: // 不限时间
			timea = 0
		}

	}

	var sqlStoreType []int
	if storeType != constant.ShopTypeOfferReward {
		if storeType == 0 {
			for i := 1; i < int(constant.ShopTypeAll); i++ {
				if i != int(constant.ShopTypeOfferReward) {
					sqlStoreType = append(sqlStoreType, i)
				}
			}
		} else {
			sqlStoreType = append(sqlStoreType, int(storeType))
		}
		storeOther := []Store{}
		err := database.Select("goods_uuid,sum(buy_num) as  num,min(created_at) create_time").Where("status = ?   and created_at>? ",
			StatusNormal, timea).In("uid", uid).In("store_type", sqlStoreType).In("goods_uuid", goodsIds).GroupBy("goods_uuid").
			OrderBy("create_time").Find(&storeOther)
		if err != nil {
			logger.Error(err)
			return nil, errutil.ErrDBOperation
		}
		store = append(store, storeOther...)

	}
	if storeType == constant.ShopTypeOfferReward || storeType == 0 {
		OfferReward := []Store{}
		var uids []int64
		userInfo := &model.User{Uid: uid}
		has, err := database.Get(userInfo)
		if err != nil {
			return nil, err
		}
		if !has {
			return nil, err
		}
		users := []model.User{}
		err = database.Select("uid").Where("account_name = ?  ",
			userInfo.AccountName).Find(&users)
		if err != nil {
			logger.Error(err)
			return nil, errutil.ErrDBOperation
		}
		// 查找所有的用户
		for _, resb := range users {
			uids = append(uids, resb.Uid)
		}
		err = database.Select("goods_uuid,sum(buy_num) as  num,min(created_at) create_time").Where("status = ?   and created_at>? ",
			StatusNormal, timea).In("uid", uids).In("store_type", constant.ShopTypeOfferReward).In("goods_uuid", goodsIds).GroupBy("goods_uuid").
			OrderBy("create_time").Find(&OfferReward)
		if err != nil {
			logger.Error(err)
			return nil, errutil.ErrDBOperation
		}
		store = append(store, OfferReward...)
	}

	return store, nil
}

// 特殊情况，game需要查询这个商城购买了几次
func StoreTotal(uid int64, storeType int32) ([]Store, error) {
	storeOther := []Store{}
	sqlStoreType := []int{int(storeType)}
	err := database.Select("goods_uuid,sum(buy_num) as  num,min(created_at) create_time").Where("status = ?   ",
		StatusNormal).In("uid", uid).In("store_type", sqlStoreType).GroupBy("goods_uuid").
		OrderBy("create_time").Find(&storeOther)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return storeOther, nil
}
func StoreAdd(uid int64, goodsUuid int32, SourceInfo string, buyNum int32, createTime int64) (bool, error) {

	u := &model.Store{
		Uid:        uid,
		StoreType:  goodsUuid / 1000,
		GoodsId:    int64(goodsUuid % 1000),
		GoodsUuid:  goodsUuid,
		SourceInfo: SourceInfo,
		BuyNum:     buyNum,
		CreatedAt:  createTime,
		Status:     StatusNormal,
	}
	_, err := database.Insert(u)

	if err != nil {
		logger.Error(err)
		return false, errutil.ErrDBOperation
	}
	return true, nil
}

func ExtractAccountInfo(uid int64) ([]model.UserExtractAccount, error) {
	result := make([]model.UserExtractAccount, 0)
	err := database.Select("type,detail,name,created_at").Where("status = ? and uid = ?  ",
		StatusNormal, uid).OrderBy("created_at").Find(&result)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}
func ExtractAccountAdd(uid int64, detail string, etype int32, name string) error {

	u := &model.UserExtractAccount{
		Uid:       uid,
		Type:      etype,
		Detail:    detail,
		CreatedAt: time.Now().Unix(),
		Status:    StatusNormal,
		Name:      name,
	}
	_, err := database.Insert(u)

	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}
func ExtractAccountUpdate(uid int64, detail string, etype int32, name string) error {
	u := &model.UserExtractAccount{
		Type:      etype,
		Detail:    detail,
		Name:      name,
		UpdatedAt: time.Now().Unix(),
	}
	_, err := database.Where("uid =?  ", uid).Update(u)

	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// 用户购买情况
func Storetotal(storeType int32) ([]Store, error) {

	store := []Store{}
	//err := database.Where("uid = ? and storetype =?  and created_at> ? and goods_id in(?)",
	//	uid, storeType, dayType*3600*24, goodsIds).
	//	Desc("created_at").Limit(1).Find(&result, Sign)
	now := time.Now()                                                                                         // 获取当前时间
	tomorrow := now.AddDate(0, 0, -1)                                                                         // 增加一天
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()) // 设置为0点0分0秒
	timea := midnight.Unix()
	var sqlStoreType []int
	if storeType == 0 {
		store := global.StoreCfg
		while := make(map[int32]bool)
		for _, v := range store {
			if v.Limit.NeedMoney == constant.NeedMoney && !while[v.StoreType] {
				sqlStoreType = append(sqlStoreType, int(v.StoreType))
				while[v.StoreType] = true
			}
		}
	} else {
		sqlStoreType = append(sqlStoreType, int(storeType))
	}
	err := database.Select("sum(props_id*props_num) as  num,min(created_at) create_time").Where("status = ? and created_at>? ",
		StatusNormal, timea).In("store_type", sqlStoreType).
		OrderBy("created_at").Find(&store)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return store, nil
}
