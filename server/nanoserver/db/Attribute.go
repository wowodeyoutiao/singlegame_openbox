package db

import (
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
)

// QueryUser get the user by id
func QueryUserAttribute(uid int64) (*model.UserAttribute, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	uAttribute := &model.UserAttribute{
		Uid: uid,
	}
	has, err := database.Get(uAttribute)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	if !has {
		return nil, errutil.ErrUserNotFound
	}
	return uAttribute, nil

}

// UpdateUser update user's info
func UpdateUserAttribute(p *model.UserAttribute) error {
	if p == nil {
		return nil
	}
	_, err := database.Where("uid=? ", p.Uid).Update(p)
	return err
}

// InsertUser insert a new user
func InsertUserAttribute(u *model.UserAttribute) error {
	if u == nil {
		return nil
	}
	_, err := database.Insert(u)
	return err
}
