package db

import (
	"errors"
	"strconv"
	"time"

	"nanoserver/db/model"
	"nanoserver/global"

	//"nanoserver/internal/web"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

func MailDetail(uid int64, MailId string) ([]model.Mail, error) {
	// 判断是否是自己的
	check, err := MailCheck(uid, MailId)
	if err != nil {
		return nil, errutil.ErrIllegalParameter
	}
	if len(check) == 0 {
		return nil, errutil.ErrIllegalParameter
	}
	// 如果有附件，查询附件内容
	if check[0].IsAttachments == 1 {
		// Attachment := &model.Attachment{}
		// result := make([]model.Attachment, 0)
		// err := database.Where(" mail_id = ?", MailId).Find(&result, Attachment)

		Attachment := []model.Attachment{}
		// result := make([]model.Attachment, 0)
		err := database.Where(" mail_id = ?", MailId).Find(&Attachment)
		if err != nil {
			logger.Error(err)
			return nil, errutil.ErrDBOperation
		}
		check[0].Attachments = Attachment
	}
	// 如果没有读过改成已读
	if check[0].IsRead == 0 && check[0].RecipientUid == uid {
		u := &model.Mail{IsRead: 1}
		if _, err := database.Where("mail_id=?", MailId).Update(u); err != nil {
			return nil, err
		}
	}
	check[0].IsRead = 1
	return check, nil
}

func MailTaskDetail() ([]model.MailTask, error) {
	// 判断是否是自己的
	Mail := &model.MailTask{}

	result := make([]model.MailTask, 0)

	// 判断是不是有官方邮件
	err := database.Where(" expiry_at > ? and created_at < ?", time.Now().Unix(), time.Now().Unix()).Desc("id").Find(&result, Mail)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	// 如果有附件，查询附件内容
	for i, item := range result {
		if result[i].IsAttachments == 1 {
			// Attachment := &model.Attachment{}
			// result := make([]model.Attachment, 0)
			// err := database.Where(" mail_id = ?", MailId).Find(&result, Attachment)

			AttachmentTask := []model.AttachmentTask{}
			// result := make([]model.Attachment, 0)
			err := database.Where(" mail_task_id = ?", item.MailTaskId).Find(&AttachmentTask)
			if err != nil {
				logger.Error(err)
				return nil, errutil.ErrDBOperation
			}
			result[i].Attachments = AttachmentTask
		}
	}
	return result, nil
}

func InsertMail(Mail *model.Mail) error {
	if Mail == nil {
		return errutil.ErrInvalidParameter
	}
	_, err := database.Insert(Mail)
	if err != nil {
		return errutil.ErrDBOperation
	}
	return nil
}

func MailCheck(uid int64, MailId string) ([]model.Mail, error) {
	Mail := &model.Mail{}

	result := make([]model.Mail, 0)
	err := database.Where("((sender_uid = ? and is_delete_sender=0) or (recipient_uid = ? and is_delete_recipient=0)) and mail_id = ?", uid, uid, MailId).
		Desc("id").Find(&result, Mail)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 判断是否已经获取过
func MailTaskCheck(uid int64, ChannelCode string) ([]model.Mail, error) {
	Mail := &model.Mail{}

	result := make([]model.Mail, 0)
	err := database.Where("recipient_uid =? and channel_code = ?", uid, ChannelCode).
		Desc("id").Find(&result, Mail)
	if err != nil {
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 用户邮件list
func MailList(uid int64, offset, count int) ([]model.Mail, int, error) {
	Mail := &model.Mail{}
	//total, err := database.Where("(sender_uid = ? and is_delete_sender=0) or (recipient_uid = ? and is_delete_recipient=0)", uid, uid).Count(Mail)
	//if err != nil {
	//	logger.Error(err)
	//	return nil, 0, errutil.ErrDBOperation
	//}
	currentTime := time.Now()
	// 获取当前时间戳（秒）
	timestamp := currentTime.Unix()
	result := make([]model.Mail, 0)
	total, err := database.Where("(sender_uid = ? and is_delete_sender=0) or (recipient_uid = ? and is_delete_recipient=0) and expiry_at>?", uid, uid, timestamp).
		Desc("created_at").Limit(count, offset).FindAndCount(&result, Mail)
	if err != nil {
		logger.Error(err)
		return nil, 0, errutil.ErrDBOperation
	}
	return result, int(total), nil
}

func MailEextractDetail(uid int64, MailId string) ([]protocol.SourceInfo, error, string) {
	// 判断是否是自己的
	check, err := MailCheck(uid, MailId)
	if err != nil {
		return nil, errutil.ErrChooseRole, ""
	}
	if len(check) == 0 {
		return nil, errutil.ErrChooseRole, ""
	}

	if check[0].IsExtract == 1 {
		return nil, errutil.ErrGetAward, ""
	}
	if check[0].RecipientUid != uid {
		return nil, errutil.ErrChooseRole, ""
	}
	// 需要配合到账接口
	// 发放签到的奖励
	Attachment := []model.Attachment{}
	// result := make([]model.Attachment, 0)
	err3 := database.Where(" mail_id = ?", MailId).Find(&Attachment)
	result := make([]protocol.SourceInfo, len(Attachment))
	for i, info := range Attachment {
		result[i].Id = int32(info.UserPropsId)
		result[i].Num = info.Num
	}
	if err3 != nil {
		// logger.Error(err3)
		return nil, errutil.ErrDBOperation, ""
	}
	return result, nil, check[0].Title
	//_, err2 := api.Send(result, uid)
	//if err2 != nil {
	//	return false, err2
	//}
}

func MailEextract(uid int64, MailId []string) (bool, error) {
	u := &model.Mail{
		IsRead:    1,
		IsExtract: 1,
	}
	currentTime := time.Now()
	// 获取当前时间戳（秒）
	timestamp := currentTime.Unix()
	_, errr := database.Where("recipient_uid=?  and expiry_at>?  ", uid, timestamp).In("mail_id", MailId).Update(u)
	if errr != nil {
		return false, errr
	}
	return true, errr
}

func MailEextractNeed(uid int64) ([]model.Mail, error) {
	Mail := &model.Mail{}
	//total, err := database.Where("(sender_uid = ? and is_delete_sender=0) or (recipient_uid = ? and is_delete_recipient=0)", uid, uid).Count(Mail)
	//if err != nil {
	//	logger.Error(err)
	//	return nil, 0, errutil.ErrDBOperation
	//}
	result := make([]model.Mail, 0)
	currentTime := time.Now()
	// 获取当前时间戳（秒）
	timestamp := currentTime.Unix()
	errr := database.Where("is_extract =0 and is_attachments =1 and recipient_uid=? and  expiry_at>?", uid, timestamp).Find(&result, Mail)
	if errr != nil {
		return nil, errr
	}

	return result, nil
}

func MailEextractAll(uid int64, MailId []string) ([]protocol.SourceInfo, error) {
	Attachment := []model.Attachment{}
	// result := make([]model.Attachment, 0)
	err3 := database.In("mail_id", MailId).Find(&Attachment)
	if err3 != nil {
		// logger.Error(err3)
		return nil, errutil.ErrDBOperation
	}
	result := make([]protocol.SourceInfo, len(Attachment))
	for i, info := range Attachment {
		result[i].Id = int32(info.UserPropsId)
		result[i].Num = info.Num
	}
	return result, nil
}

func MailDel(uid int64, MailId string) (int, error) {
	// 判断是否是自己的
	check, err := MailCheck(uid, MailId)
	if err != nil {
		return StatusDeleted, errutil.ErrIllegalParameter
	}
	if len(check) == 0 {
		return StatusDeleted, errutil.ErrIllegalParameter
	}
	sess := database.NewSession()

	if check[0].SenderUid == uid {
		u := &model.Mail{IsDeleteSender: 1}
		_, err1 := database.Where("mail_id =? and sender_uid=? ", MailId, uid).Update(u)
		if err1 != nil {
			sess.Rollback()
			return StatusDeleted, err1
		}
	}
	if check[0].RecipientUid == uid {
		if check[0].IsRead == 0 {
			return StatusFreezed, errors.New("文件未读")
		}
		// 判断是不是有附件没提取
		if check[0].IsAttachments == 1 && check[0].IsExtract == 0 {
			sess.Rollback()
			return StatusFreezed, errors.New("有附件没提取")
		}
		u := &model.Mail{IsDeleteRecipient: 1}
		_, err2 := database.Where("mail_id =?  and  is_read =1 ", MailId).Update(u)
		if err2 != nil {
			sess.Rollback()
			return StatusDeleted, err2
		}
	}
	sess.Commit()
	return StatusNormal, err
}

// func MailDelAll(uid int64) (bool, error) {
// 	sess := database.NewSession()

// 	// 开始事务
// 	sess.Begin()

// 	u := &model.Mail{IsDeleteSender: 1}
// 	_, err1 := sess.Where("sender_uid  =?   ", uid).Update(u)
// 	if err1 != nil {
// 		// 处理错误
// 		sess.Rollback()
// 	}
// 	m := &model.Mail{IsDeleteRecipient: 1}
// 	_, err2 := sess.Where("((is_attachments =1 and is_extract =1) or is_attachments=0 )and  is_read =1 and recipient_uid =?", uid).Update(m)
// 	if err2 != nil {
// 		// 处理错误
// 		sess.Rollback()
// 	}
// 	sess.Commit()
// 	return true, nil
// }

func MailDelAll(uid int64) (bool, int64, error) {
	sess := database.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return false, 0, err
	}

	u := &model.Mail{IsDeleteSender: 1}
	senderCount, err1 := sess.Where("sender_uid = ?", uid).Update(u)
	if err1 != nil {
		sess.Rollback()
		return false, 0, err1
	}

	m := &model.Mail{IsDeleteRecipient: 1}
	recipientCount, err2 := sess.Where("((is_attachments = 1 and is_extract = 1) or is_attachments = 0) and is_read = 1 and recipient_uid = ?", uid).Update(m)
	if err2 != nil {
		sess.Rollback()
		return false, 0, err2
	}

	if err := sess.Commit(); err != nil {
		return false, 0, err
	}

	totalDeleted := senderCount + recipientCount
	return true, totalDeleted, nil
}

func MailCount(RecipientUid int64) error {
	Mail := []model.Mail{}
	// result := make([]model.Attachment, 0)
	currentTime := time.Now()
	// 获取当前时间戳（秒）
	timestamp := currentTime.Unix()
	err := database.Where(" recipient_uid = ? and is_delete_recipient =? and  expiry_at>?", RecipientUid, 0, timestamp).Find(&Mail)
	if err != nil {
		// logger.Error(err3)
		return errutil.ErrDBOperation
	}

	if len(Mail) >= int(global.GlobalCfg.BaseValue.Mailmax.Num) {
		return errutil.ErrMailMax
	}
	return nil
}

func MailAdd(uid int64, RecipientUid int64, title string, desc string, attachments []protocol.Attachment, ChannelCode string) (string, bool, error) {
	mailCountRes := MailCount(RecipientUid)
	if mailCountRes != nil {
		return "", false, mailCountRes
	}
	sess := database.NewSession()
	sess.Begin()
	MailId := strconv.FormatInt(snowflake.GenID(), 10)
	is_attachments := 0
	result := make([]protocol.Attachment, len(attachments))
	if len(attachments) > 0 {
		for i, Attachment := range attachments {
			result[i] = protocol.Attachment{
				MailId:      MailId,
				UserPropsId: Attachment.UserPropsId,
				Num:         Attachment.Num,
			}
		}
		_, err := sess.Insert(result)
		if err != nil {
			sess.Rollback()
			return "", false, err
		}
		is_attachments = 1
	}
	IsOfficial := 0
	if uid == -1 {
		IsOfficial = 1
	}
	u := &model.Mail{
		MailId:        MailId,
		SenderUid:     uid,
		RecipientUid:  RecipientUid,
		IsOfficial:    IsOfficial,
		Title:         title,
		Desc:          desc,
		CreatedAt:     time.Now().Unix(),
		ExpiryAt:      time.Now().AddDate(0, 0, 7).Unix(),
		IsAttachments: is_attachments,
		ChannelCode:   ChannelCode,
	}
	_, err := sess.Insert(u)
	if err != nil {
		sess.Rollback()
		return "", false, err
	}
	// sess.Commit()
	sess.Commit()
	return MailId, true, err
}
