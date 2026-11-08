package db

import (
	"nanoserver/db/model"
)

// GetRefineByUid 根据用户ID获取洗练信息
func GetRefineByUid(uid int64) (*model.Refine, error) {
	refine := &model.Refine{UID: uid}

	has, err := database.Get(refine)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}
	return refine, nil
}

// InsertRefine 插入新的洗练信息
func InsertRefine(refine *model.Refine) error {
	_, err := database.Insert(refine)
	return err
}

// UpdateRefine 更新洗练信息
func UpdateRefine(refine *model.Refine) error {
	_, err := database.Where("uid = ?", refine.UID).
		Cols("purchased_slots", "attrs").
		Update(refine)
	return err
}
