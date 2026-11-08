package db

import (
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
)

// 查询记录
func QueryPvpRecords(uid int64) ([]*model.PvpRecord, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}

	var result []*model.PvpRecord

	err := database.Where("self_uid=? and status=1", uid).Find(&result)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// InsertPvpRecord 插入记录
func InsertPvpRecord(record *model.PvpRecord) error {
	_, err := database.Insert(record)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// DeletePvpRecord 删除记录
func DeletePvpRecord(uid int64) error {
	_, err := database.Where("uid=?", uid).Delete(model.PvpRecord{})
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// UpdatePvpRecord 更新记录
func UpdatePvpRecord(record *model.PvpRecord) error {
	_, err := database.Where("uid=?", record.Uid).Update(record)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// UpdatePvpRecordStatus 更新记录状态 并且更新user表的分数
func UpdateAllPvpRecord(status int8) error {
	_, err1 := database.Exec("UPDATE `pvp_record` SET `status` = ? ", status)
	if err1 != nil {
		return err1
	}

	_, err2 := database.Exec("UPDATE `user` SET `pvp_score`=1000 ,`gvg_score`=1000 ")
	if err2 != nil {
		return err2
	}
	return nil
}

// /////////////战场///////////////
// 插入记录 War
func InsertWar(war *model.War) error {
	_, err := database.Insert(war)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}

// 查询记录 War
func QueryWar(uid int64) (*model.War, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	u := &model.War{
		Uid: uid,
	}
	has, err := database.Get(u)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return u, nil
}

// 更新记录 War
func UpdateWar(war *model.War) error {
	_, err := database.Where("uid=?", war.Uid).AllCols().Update(war)
	if err != nil {
		logger.Error(err)
		return errutil.ErrDBOperation
	}
	return nil
}
