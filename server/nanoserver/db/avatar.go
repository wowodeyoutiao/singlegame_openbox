package db

import (
	"fmt"

	"nanoserver/db/model"
)

// GetAvatarByUid 根据用户ID获取头像信息
func GetAvatarByUid(uid int64) (*model.Avatar, error) {
	avatar := &model.Avatar{UID: uid}

	has, err := database.Get(avatar)
	if err != nil {
		return nil, fmt.Errorf("get avatar by uid failed: %w", err)
	}

	if !has {
		return nil, nil
	}

	return avatar, nil
}

// InsertAvatar 插入新的头像信息
func InsertAvatar(avatar *model.Avatar) error {
	_, err := database.Insert(avatar)
	return err
}

// UpdateAvatar 更新头像信息
func UpdateAvatar(avatar *model.Avatar) error {
	_, err := database.Where("uid = ?", avatar.UID).
		Cols("owned_avatars", "owned_frames", "current_avatar", "current_frame").
		Update(avatar)
	return err
}

// UpsertAvatar 保存或更新用户头像信息（Upsert操作）
// func UpsertAvatar(avatar *model.Avatar) error {
// 	session := database.NewSession()
// 	defer session.Close()

// 	if err := session.Begin(); err != nil {
// 		return fmt.Errorf("begin transaction failed: %w", err)
// 	}

// 	// 先检查记录是否存在
// 	exists, err := session.Exist(&model.Avatar{Id: avatar.Id})
// 	if err != nil {
// 		session.Rollback()
// 		return fmt.Errorf("check avatar existence failed: %w", err)
// 	}

// 	if exists {
// 		// 记录存在，进行更新
// 		_, err := session.Where("id = ?", avatar.Id).
// 			Cols("owned_avatars", "owned_frames", "current_avatar", "current_frame").
// 			Update(avatar)
// 		if err != nil {
// 			session.Rollback()
// 			return fmt.Errorf("update avatar failed: %w", err)
// 		}
// 	} else {
// 		// 记录不存在，进行插入
// 		if _, err := session.Insert(avatar); err != nil {
// 			session.Rollback()
// 			return fmt.Errorf("insert avatar failed: %w", err)
// 		}
// 	}

// 	return session.Commit()
// }
