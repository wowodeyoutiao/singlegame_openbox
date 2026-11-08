package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/internal/game"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	"github.com/spf13/viper"
)

func MakeMailService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/mail/add", nex.Handler(mailAdd)).Methods("Put") // 添加
	// router.Handle("/v1/mail/addall", nex.Handler(mailAddAll)).Methods("Put")                                           // 添加
	router.Handle("/v1/mail/list", nex.Handler(mailList)).Methods("POST")                                              // 列表
	router.Handle("/v1/mail/detail", nex.Handler(mailDetail)).Methods("POST")                                          // 详情
	router.Handle("/v1/mail/del", nex.Handler(mailDel)).Methods("Delete")                                              // 删除
	router.Handle("/v1/mail/delall", nex.Handler(mailDelAll)).Methods("Delete")                                        // 删除全部
	router.Handle("/v1/mail/extract", limitHandler(nex.Handler(mailExtract), "mail_extract")).Methods("Post")          // 提取
	router.Handle("/v1/mail/extractall", limitHandler(nex.Handler(mailExtractAll), "mail_extractall")).Methods("Post") // 提取全部
	// router.Handle("/v1/mail/MailTask", limitHandler(nex.Handler(mailTask))).Methods("Post")         //登录触发发邮件

	return router
}

func MailList(r *protocol.MailListRequest) ([]protocol.MailInfo, int, error) {
	if r == nil {
		return nil, 0, errutil.ErrIllegalParameter
	}
	//id, err := strconv.ParseInt(r.Uid, 10, 0)
	//if err != nil {
	//	return nil, 0, err
	//}

	list, total, err := db.MailList(
		r.Uid,
		r.Offset,
		r.Count,
	)
	if err != nil {
		return nil, 0, err
	}

	result := make([]protocol.MailInfo, len(list))
	for i, Mail := range list {
		if Mail.SenderUid == r.Uid {
			result[i] = protocol.MailInfo{
				MailId:        Mail.MailId,
				CreatedAt:     Mail.CreatedAt,
				ExpiryAt:      Mail.ExpiryAt,
				Title:         Mail.Title,
				Desc:          Mail.Desc,
				IsAttachments: Mail.IsAttachments,
				IsExtract:     0,
				IsRead:        1,
			}
		} else {
			result[i] = protocol.MailInfo{
				MailId:        Mail.MailId,
				CreatedAt:     Mail.CreatedAt,
				ExpiryAt:      Mail.ExpiryAt,
				Title:         Mail.Title,
				Desc:          Mail.Desc,
				IsAttachments: Mail.IsAttachments,
				IsExtract:     Mail.IsExtract,
				IsRead:        Mail.IsRead,
			}
		}
	}

	if err := game.SetMailCounter(r.Uid, constant.MailCounterOPTypeSet, len(result)); err != nil {
		logger.Errorf("SetMailCounter error: %v", err)
	}

	return result, total, nil
}

func MailDetail(r *protocol.MailListRequest) ([]protocol.MailInfo, error) {
	if r == nil {
		return nil, errutil.ErrIllegalParameter
	}
	//id, err := strconv.ParseInt(r.Uid, 10, 0)
	//if err != nil {
	//	return nil, 0, err
	//}

	list, err := db.MailDetail(
		r.Uid,
		r.MailId,
	)
	if err != nil {
		return nil, err
	}

	result := make([]protocol.MailInfo, len(list))
	for i, Mail := range list {
		if Mail.SenderUid == r.Uid {
			result[i] = protocol.MailInfo{
				MailId:        Mail.MailId,
				CreatedAt:     Mail.CreatedAt,
				ExpiryAt:      Mail.ExpiryAt,
				Title:         Mail.Title,
				Desc:          Mail.Desc,
				IsAttachments: Mail.IsAttachments,
				IsExtract:     0,
				IsRead:        1,
				Attachments:   make([]protocol.Attachment, len(Mail.Attachments)),
			}
		} else {
			result[i] = protocol.MailInfo{
				MailId:        Mail.MailId,
				CreatedAt:     Mail.CreatedAt,
				ExpiryAt:      Mail.ExpiryAt,
				Title:         Mail.Title,
				Desc:          Mail.Desc,
				IsAttachments: Mail.IsAttachments,
				IsExtract:     Mail.IsExtract,
				IsRead:        Mail.IsRead,
				Attachments:   make([]protocol.Attachment, len(Mail.Attachments)),
			}
		}

		// 复制 Attachments 字段
		for j, attachment := range Mail.Attachments {
			// 假设 Attachment 结构体中有可以复制的字段
			result[i].Attachments[j] = protocol.Attachment{
				// ... 根据字段一一复制值
				MailId:      attachment.MailId,
				UserPropsId: attachment.UserPropsId,
				Num:         attachment.Num,
			}
		}
	}
	return result, nil
}

func MailEextract(r *protocol.MailListRequest) (bool, error) {
	if r == nil {
		return false, errutil.ErrIllegalParameter
	}

	result, err, title := db.MailEextractDetail(
		r.Uid,
		r.MailId,
	)

	if err != nil {
		return false, err
	}
	////发放签到的奖励
	err2 := CheckSourceAdd(result, r.Uid, true, constant.ItemAddMail)
	if err2 != nil {
		return false, err2
	}
	// 改变状态
	var mailIds []string
	mailIds = append(mailIds, r.MailId)
	resultb, err3 := db.MailEextract(
		r.Uid,
		mailIds,
	)
	if !resultb {
		return false, err3
	}
	upMailGetShushu(r.Uid, result, r.MailId, title)
	return true, nil
}

func MailEextractAll(r *protocol.MailListRequest) (bool, error) {
	if r == nil {
		return false, errutil.ErrIllegalParameter
	}
	// 查询所有有附件的且没有被提取的
	resulta, err := db.MailEextractNeed(
		r.Uid,
	)
	if err != nil {
		return false, err
	}

	var mailIds []string
	for _, item := range resulta {
		mailIds = append(mailIds, item.MailId)
	}
	result, err := db.MailEextractAll(r.Uid, mailIds)
	if err != nil {
		logger.Errorf("MailExtract error: %v", err)
		return false, err
	}
	////发放签到的奖励
	err = CheckSourceAdd(result, r.Uid, true, constant.ItemAddMail)
	if err != nil {
		return false, err
	}
	// 改变状态
	_, err = db.MailEextract(
		r.Uid,
		mailIds,
	)
	if err != nil {
		return false, err
	}
	upMailGetShushu(r.Uid, result, r.MailId, "一键提取")
	return true, nil
}

func upMailGetShushu(uid int64, sourceInfo []protocol.SourceInfo, mail_id string, mail_type string) {
	data := make(map[string]interface{}, 0)
	data["name"] = "mail_open"
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	data["get_item_info"] = get_item_info

	data["mail_type"] = mail_type
	data["mail_id"] = mail_id
	game.Upshushu(uid, data, true, false)
}

func MailDel(r *protocol.MailListRequest) (int, error) {
	if r == nil {
		return 2, errutil.ErrIllegalParameter
	}
	result, err := db.MailDel(
		r.Uid,
		r.MailId,
	)

	if err := game.SetMailCounter(r.Uid, constant.MailCounterOPTypeSub, 1); err != nil {
		logger.Errorf("SetMailCounter error: %v", err)
	}

	return result, err
}

func MailDelAll(r *protocol.MailListRequest) (bool, error) {
	if r == nil {
		return false, errutil.ErrIllegalParameter
	}
	result, count, err := db.MailDelAll(
		r.Uid,
	)
	if err != nil || !result {
		return false, err
	}

	if err := game.SetMailCounter(r.Uid, constant.MailCounterOPTypeSub, int(count)); err != nil {
		logger.Errorf("SetMailCounter error: %v", err)
	}

	return true, nil
}

func Mailadd(r *protocol.MailAddRequest) (bool, error) {
	if r == nil {
		return false, errutil.ErrIllegalParameter
	}
	mailId, result, err := db.MailAdd(
		r.Uid,
		r.RecipientUid,
		r.Title,
		r.Desc,
		r.Attachment,
		r.ChannelCode,
	)
	if err != nil || !result {
		return false, err
	}
	Notification(r.RecipientUid)
	upMailpShushuDate(mailId, r.RecipientUid, r.Attachment)
	return true, nil
}
func upMailpShushuDate(mail_id string, Uid int64, sourceInfo []protocol.Attachment) {
	shushuDate := make(map[string]interface{}, 0)
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(int32(v.UserPropsId)); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["get_item_info"] = get_item_info
	shushuDate["mail_id"] = mail_id
	shushuDate["name"] = "mail_receive"
	game.Upshushu(Uid, shushuDate, true, false)
}

func mailList(ctx context.Context, r *http.Request, form *nex.Form) (*protocol.MailListResponse, error) {
	request := &protocol.MailListRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	//utoken := r.Header.Get("Authorization")
	//tokenFactory := token.CreateUserFactory()
	//userInfo, err := tokenFactory.ParseToken(utoken)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println("GetUserToken", userInfo)

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	// request.Uid = strconv.FormatInt(Uid, 10)
	request.Offset = form.IntOrDefault("offset", 0)
	request.Count = form.IntOrDefault("count", -1)

	//_, err2 := mailTask(Uid)
	//if err2 != nil {
	//}
	list, total, err := MailList(request)
	if err != nil {
		return nil, err
	}
	return &protocol.MailListResponse{Data: list, Total: total}, nil
}

func mailDetail(ctx context.Context, r *http.Request, date *protocol.MailOperateRequest) (*protocol.MailDetailResponse, error) {
	request := &protocol.MailListRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	if date.MailId == "" {
		return &protocol.MailDetailResponse{Data: protocol.MailInfo{}}, errutil.ErrInvalidParameter
	}
	request.MailId = date.MailId

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	detail, err := MailDetail(request)
	if err != nil {
		return nil, err
	}
	return &protocol.MailDetailResponse{Data: detail[0]}, nil
}

func mailExtract(ctx context.Context, date *protocol.MailOperateRequest) (*protocol.MailOperateResponse, error) {
	request := &protocol.MailListRequest{}
	if date.MailId == "" {
		return &protocol.MailOperateResponse{Data: "缺少参数"}, errutil.ErrInvalidParameter
	}
	request.MailId = date.MailId

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	result, err := MailEextract(request)
	if err == errutil.ErrGetAward {
		return &protocol.MailOperateResponse{Code: -1, Data: errutil.ErrGetAward.Error()}, nil
	}
	if err != nil || !result {
		return &protocol.MailOperateResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.MailOperateResponse{Data: "成功"}, nil
}

func mailExtractAll(ctx context.Context) (*protocol.MailOperateResponse, error) {
	request := &protocol.MailListRequest{}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	result, err := MailEextractAll(request)
	if err != nil || !result {
		return &protocol.MailOperateResponse{Data: err.Error(), Code: -1}, nil
	}
	return &protocol.MailOperateResponse{Data: "成功"}, nil
}

func mailDel(ctx context.Context, date *protocol.MailOperateRequest) (*protocol.MailOperateResponse, error) {
	request := &protocol.MailListRequest{}

	if date.MailId == "" {
		return &protocol.MailOperateResponse{Data: "缺少参数"}, errutil.ErrInvalidParameter
	}
	request.MailId = date.MailId

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	result, err := MailDel(request)

	if result == db.StatusDeleted {
		return &protocol.MailOperateResponse{Data: err.Error()}, err
	}
	if result == db.StatusFreezed {
		return &protocol.MailOperateResponse{Code: 0, Data: err.Error()}, nil
	}
	return &protocol.MailOperateResponse{Data: "成功"}, nil
}

func mailDelAll(ctx context.Context) (*protocol.MailOperateResponse, error) {
	request := &protocol.MailListRequest{}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	result, err := MailDelAll(request)
	if err != nil || !result {
		return &protocol.MailOperateResponse{}, err
	}
	return &protocol.MailOperateResponse{}, err
}

func mailAdd(ctx context.Context, r *http.Request, data *protocol.MailAddRequest) (*protocol.MailOperateResponse, error) {
	if strings.ToLower(viper.GetString("core.env")) == constant.EnvProd {
		return &protocol.MailOperateResponse{Data: "正式不能发"}, errors.New("正式不能发")
	}
	request := &protocol.MailAddRequest{}
	err := r.ParseForm()
	Uid := int64(0)
	if err != nil {
		return &protocol.MailOperateResponse{}, err
	}
	if ctx.Value("userDate") == nil {
	} else {
		Uid = ctx.Value("userDate").(my_jwt.CustomClaims).UserId
		if Uid == 0 {
			Uid = -1
		}
	}

	request.Uid = Uid
	request.RecipientUid = data.RecipientUid
	request.Title = data.Title
	request.Desc = data.Desc
	request.Attachment = data.Attachment
	if request.RecipientUid == Uid {
		return &protocol.MailOperateResponse{}, errors.New("不能发给自己")
	}
	if request.RecipientUid == 0 || request.Title == "" || request.Desc == "" {
		return &protocol.MailOperateResponse{}, errutil.ErrInvalidParameter
	}
	// todo 校验是不是好友
	// 插入
	result, err := Mailadd(request)
	if err != nil || !result {
		return &protocol.MailOperateResponse{}, err
	}
	return &protocol.MailOperateResponse{Data: "成功"}, err
}
func mailAddAll(ctx context.Context, r *http.Request, data *protocol.MailAddRequest) (*protocol.MailOperateResponse, error) {
	request := &protocol.MailAddRequest{}
	err := r.ParseForm()
	Uid := int64(0)
	if err != nil {
		return &protocol.MailOperateResponse{}, err
	}
	if ctx.Value("userDate") == nil {
	} else {
		Uid = ctx.Value("userDate").(my_jwt.CustomClaims).UserId
		if Uid == 0 {
			Uid = -1
		}
	}

	request.Uid = Uid
	request.RecipientUid = data.RecipientUid
	request.Title = data.Title
	request.Desc = data.Desc
	request.Attachment = data.Attachment
	if request.RecipientUid == Uid {
		return &protocol.MailOperateResponse{}, errors.New("不能发给自己")
	}
	if request.RecipientUid == 0 || request.Title == "" || request.Desc == "" {
		return &protocol.MailOperateResponse{}, errutil.ErrInvalidParameter
	}
	// todo 校验是不是好友

	result, err := Mailadd(request)
	if err != nil || !result {
		return &protocol.MailOperateResponse{}, err
	}
	return &protocol.MailOperateResponse{Data: "成功"}, err
}
func MailTask(Uid int64) (bool, error) {
	request := &protocol.MailAddRequest{}
	resulta, err := Mailtask(request)

	// request.Title = resulta.Title
	// request.Desc = resulta.Desc
	// request.Attachment = resulta.Attachment

	// 插入
	for _, v := range resulta {
		channelCode := strconv.Itoa(int(v.MailTaskId))

		resultb, _ := db.MailTaskCheck(Uid, "MailTask-"+channelCode)
		if len(resultb) > 0 {
			continue
		}
		request.Uid = -1
		request.RecipientUid = Uid
		request.Title = v.Title
		request.Desc = v.Desc

		request.ChannelCode = "MailTask-" + channelCode
		request.Attachment = make([]protocol.Attachment, len(v.Attachments))
		for i, v2 := range v.Attachments {
			request.Attachment[i] = protocol.Attachment{
				UserPropsId: v2.PropsId,
				Num:         v2.Num,
			}
		}
		result, err2 := Mailadd(request)
		if err2 != nil || !result {
			continue
		}
	}
	return true, err
}

func Mailtask(r *protocol.MailAddRequest) ([]model.MailTask, error) {
	if r == nil {
		return nil, errutil.ErrIllegalParameter
	}
	result, err := db.MailTaskDetail()
	if err != nil || len(result) == 0 {
		return nil, err
	}
	return result, nil
}
