package db

import (
	"fmt"

	"nanoserver/db/model"

	//"nanoserver/internal/web"
	"nanoserver/pkg/errutil"

	"github.com/go-xorm/xorm"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

func QueryUserBox(uid int64) (*model.UserBox, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	ubox := &model.UserBox{
		Uid: uid,
	}
	has, err := database.Get(ubox)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return ubox, nil
}

func UpdateUserBox(ubox *model.UserBox) error {
	_, err := database.Where("uid=?", ubox.Uid).Update(ubox)
	return err
}

func InsertUserBox(ubox *model.UserBox) error {
	_, err := database.Insert(ubox)
	return err
}

func InsertUserBoxV2(session *xorm.Session, ubox *model.UserBox) (affected int64, err error) {
	if ubox == nil {
		return 0, fmt.Errorf("InsertUserBoxV2: user box is nil")
	}
	return session.Insert(ubox)
}
