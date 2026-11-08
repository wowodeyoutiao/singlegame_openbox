package db

import (
	"fmt"

	"nanoserver/db/model"
	"nanoserver/pkg/errutil"

	"github.com/go-xorm/xorm"
)

// QueryUser get the user by id
func QueryUserPet(uid int64) (*model.UserPet, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	upet := &model.UserPet{
		Uid: uid,
	}
	has, err := database.Get(upet)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return upet, nil
}

// UpdateUser update user's info
func UpdateUserPet(p *model.UserPet) error {
	if p == nil {
		return nil
	}
	_, err := database.Where("uid=? ", p.Uid).Update(p)
	return err
}

// InsertUser insert a new user
func InsertUserPet(u *model.UserPet) error {
	if u == nil {
		return nil
	}
	_, err := database.Insert(u)
	return err
}

func InsertUserPetV2(session *xorm.Session, u *model.UserPet) (affected int64, err error) {
	if u == nil {
		return 0, fmt.Errorf("InsertUserPetV2: user pet is nil")
	}
	return session.Insert(u)
}
