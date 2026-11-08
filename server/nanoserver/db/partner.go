package db

import (
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
)

// QueryUser get the user by id
func QueryPartner(uid int64) ([]model.Partner, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}

	result := make([]model.Partner, 0)

	err := database.Where(" uid=? ", uid).Desc("id").Find(&result)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil

}

// UpdateUser update user's info
func UpdatePartner(u *model.Partner) error {
	if u == nil {
		return nil
	}

	_, err := database.Where("uid=? and partner_uid=? ", u.Uid, u.PartnerUid).Cols("level, star, skill,talnet_activated,talnet_points").Update(u)
	return err
}

// InsertUser insert a new user
func InsertPartner(u *model.Partner) error {
	if u == nil {
		return nil
	}
	_, err := database.Insert(u)
	return err
}

// DeleteUser delete the user
func DeletePartner(u *model.Partner) error {

	if u == nil {
		return nil
	}
	_, err := database.Where("uid=? and partner_uid=? ", u.Uid, u.PartnerUid).Delete(u)
	return err
}

// PARTNER SLOT OPERATION
func QueryPartnerSlot(uid int64) ([]model.PartnerSlot, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}

	result := make([]model.PartnerSlot, 0)

	err := database.Where(" uid=? ", uid).Find(&result)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

func UpdatePartnerSlot(u *model.PartnerSlot) error {
	if u == nil {
		return nil
	}

	_, err := database.Where("uid=? and slot_id=?", u.Uid, u.SlotId).Cols("partner_uid, equips").Update(u)
	return err
}

func InsertPartnerSlot(u *model.PartnerSlot) error {
	if u == nil {
		return nil
	}
	_, err := database.Insert(u)
	return err
}
