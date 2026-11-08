package db

import (
	"nanoserver/db/model"
	"time"

	//"nanoserver/internal/web"
	"nanoserver/pkg/errutil"
)

const (
// noLimitFlag  = -1 //如果count == -1则表示返回所有数据
// noTimeFilter = -1 //如果start/end == -1则表示无时间筛选
)

func ActivationDetail(code string) ([]model.Activation, error) {
	//判断

	result := make([]model.Activation, 0)

	err := database.Where(" activation_code = ? and expiry_at>? and status =?", code, time.Now().Unix(), StatusNormal).
		Desc("id").Find(&result)
	if err != nil {
		return result, nil
	}
	//如果有附件，查询附件内容
	if len(result) > 0 {
		AttachmentTask := []model.AttachmentActivation{}
		err2 := database.Where(" mail_activation_id = ? ", result[0].ActivationCode).Find(&AttachmentTask)
		if err2 != nil {
			logger.Error(err2)
			return nil, errutil.ErrDBOperation
		}
		result[0].Attachments = AttachmentTask
	}
	return result, nil
}
func ActivationUse(code string) (bool, error) {
	//判断

	u := &model.Activation{Status: StatusBound}
	_, err2 := database.Where("activation_code =?  ", code).Update(u)
	if err2 != nil {
		return false, err2
	}
	return true, nil
}
