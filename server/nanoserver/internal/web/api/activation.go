package api

import (
	"context"
	"nanoserver/db"
	"nanoserver/global"
	"nanoserver/internal/game"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakeActivationService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/activation/check", nex.Handler(checkAddV2)).Methods("Put") //兑换验证码

	return router
}
func checkAdd(ctx context.Context, r *http.Request, Activation *protocol.ActivationListRequest) (*protocol.ActivationOperateResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request := &protocol.MailAddRequest{}

	//request.Title = resulta.Title
	//request.Desc = resulta.Desc
	//request.Attachment = resulta.Attachment

	//校验激活码ok不ok
	channelCode := Activation.ActivationCode

	//查询激活码是否有效
	resulta, err := db.ActivationDetail(channelCode)
	if channelCode == "" {
		return &protocol.ActivationOperateResponse{Data: "激活码不存在"}, errutil.ErrActivation
	}
	if len(resulta) == 0 {
		return &protocol.ActivationOperateResponse{Code: -1, Data: "激活码不存在"}, nil
	}

	resultb, err := db.MailTaskCheck(Uid, "Activation-"+channelCode)
	if len(resultb) > 0 {
		return &protocol.ActivationOperateResponse{Code: -1, Data: "已激活过"}, nil
	}
	if resulta[0].UseType == 2 {
		resultc, err3 := db.ActivationUse(channelCode)
		if !resultc {
			return &protocol.ActivationOperateResponse{Code: -1, Data: "服务器故障"}, err3
		}
	}
	request.Uid = -1
	request.Title = "系统邮件"
	request.Desc = "激活码兑换成功"
	request.RecipientUid = Uid
	request.ChannelCode = "Activation-" + channelCode
	request.Attachment = make([]protocol.Attachment, len(resulta[0].Attachments))
	for i, v2 := range resulta[0].Attachments {
		request.Attachment[i] = protocol.Attachment{
			UserPropsId: v2.PropsId,
			Num:         v2.Num,
		}
	}
	result, err := Mailadd(request)
	if err != nil || !result {
		return nil, err
	}
	return &protocol.ActivationOperateResponse{Data: "成功"}, nil

}

func checkAddV2(ctx context.Context, r *http.Request, req *protocol.ActivationListRequest) (*protocol.ActivationOperateResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	// 判断 Code 是否存在
	activationCodeConfig, ok := global.ActivationCodeCfg[strings.TrimSpace(req.ActivationCode)]
	if !ok {
		logger.Warnf("Activation code %s not found", req.ActivationCode)
		return &protocol.ActivationOperateResponse{Code: -1, Data: "激活码不存在"}, nil
	}

	// 判断 Code 是否在有效期
	now := time.Now()
	if activationCodeConfig.Start.After(now) || activationCodeConfig.End.Before(now) {
		logger.Warnf("Activation code %s is not valid at this time", req.ActivationCode)
		return &protocol.ActivationOperateResponse{Code: -1, Data: "兑换码已经过期"}, nil
	}

	// 获取用户信息
	userInfo, err := game.UserInfo(Uid)
	if err != nil {
		logger.Errorf("CreateOrderV2: 获取用户信息失败: %v", err)
		return &protocol.ActivationOperateResponse{Code: -1, Data: "获取用户信息失败"}, err
	}

	// 检查Code是否已被使用
	verifyActiveCodeRequest := &protocol.VerifyActiveCodeRequest{
		Uid:      uint64(Uid),
		ServerId: int32(userInfo.ServerId),
		Code:     req.ActivationCode,
		CodeType: int32(activationCodeConfig.CodeType),
		Channel:  int32(activationCodeConfig.ChannelType),
	}
	logger.Debugf("CreateOrderV2: VerifyActiveCodeRequest: %+v", verifyActiveCodeRequest)
	verifyActiveCodeResponse, err := NewPaymentClient().VerifyActiveCode(context.Background(), verifyActiveCodeRequest)
	if err != nil {
		logger.Errorf("CreateOrderV2: VerifyActiveCode failed: %v", err)
		return &protocol.ActivationOperateResponse{Code: -1, Data: errutil.ErrServerInternal.Error()}, err
	}
	if verifyActiveCodeResponse.Code != 0 {
		logger.Warnf("Activation code verification failed: %s, code: %s, uid: %d", verifyActiveCodeResponse.Message, req.ActivationCode, Uid)
		return &protocol.ActivationOperateResponse{Code: -1, Data: verifyActiveCodeResponse.Message}, nil
	}

	// 发放邮件
	rewardAttachments := make([]protocol.Attachment, len(activationCodeConfig.Reward))
	for i, v := range activationCodeConfig.Reward {
		rewardAttachments[i] = protocol.Attachment{
			UserPropsId: int64(v.ID),
			Num:         int64(v.Num),
		}
	}
	mailAddRequest := &protocol.MailAddRequest{
		Uid:          -1,
		Title:        "系统邮件",
		Desc:         "激活码兑换成功",
		RecipientUid: Uid,
		ChannelCode:  "Activation-" + req.ActivationCode,
		Attachment:   rewardAttachments,
	}
	result, err := Mailadd(mailAddRequest)
	if err != nil || !result {
		logger.Errorf("MailAdd: 发送邮件失败: %v, result: %v", err, result)
		return &protocol.ActivationOperateResponse{Code: -1, Data: errutil.ErrServerInternal.Error()}, err
	}
	upCDKShushuDate(Uid, req.ActivationCode, rewardAttachments)

	return &protocol.ActivationOperateResponse{Code: 0, Data: "兑换成功"}, nil
}
func upCDKShushuDate(Uid int64, channel_code string, sourceInfo []protocol.Attachment) {
	shushuDate := make(map[string]interface{}, 0)
	//在合并点其他的数据
	shushuDate["CDkey_content"] = channel_code
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
	shushuDate["name"] = "CDkey"
	game.Upshushu(Uid, shushuDate, true, false)
}
