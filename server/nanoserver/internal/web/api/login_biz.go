package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/dispatcher"
	"nanoserver/global"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/jwtv5x"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/timex"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"nanoserver/pkg/utils"
	"nanoserver/protocol"

	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type LoginUseCase struct{}

func NewLoginUseCase() *LoginUseCase {
	return &LoginUseCase{}
}

func (uc *LoginUseCase) GuestLogin(ctx context.Context, req *protocol.LoginRequest) (resp *protocol.LoginResponse, err error) {
	var (
		accountUID    int64
		characterList []model.User
		adSource      string
	)

	result, err := utils.ParseAdQuery(req.ClickID)
	if err != nil {
		logger.Warnf("AdClickHandler ParseAdQuery %v", err)
	} else {
		adSource = result.Source.String()
	}

	// 注册用户
	if err = db.WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		account, err := db.QueryRegisterV2(session, req.AppID, req.Device.IMEI, req.ServerId)
		if err != nil && !errors.Is(err, errutil.ErrOrderNotFound) {
			if strings.ToLower(viper.GetString("core.env")) == constant.EnvProd {
				maxServerID := global.ServerListMgrInstance.GetMaxProdServerID()
				if maxServerID != int(req.ServerId) {
					return errutil.NewBizError(-1, fmt.Sprintf("服务器关闭注册, IMEI: %s, 当前服务器ID: %d, 最大服务器ID: %d", req.Device.IMEI, req.ServerId, maxServerID), "服务器关闭注册, 请重新选择其他区服")
				}
			}
			accountUID = snowflake.GenID()
			if err := db.RegisterUserLogV2(session, accountUID, req.Device, req.AppID, req.ChannelID, protocol.RegTypeThird, req.ServerId, adSource); err != nil {
				return errors.Wrapf(err, "register failed imei: %v", req.Device.IMEI)
			}
		} else {
			// 检查账号是否被封号（封号是针对账号的，存储在 Register 表中）
			if account.BanUntil > 0 && time.Now().Unix() < account.BanUntil {
				remainingSeconds := account.BanUntil - time.Now().Unix()
				remainingMinutes := remainingSeconds / 60
				if remainingMinutes < 1 {
					remainingMinutes = 1
				}
				banUntilTime := time.Unix(account.BanUntil, 0).Format("2006-01-02 15:04:05")
				return errutil.NewBizError(-1,
					fmt.Sprintf("账号已被封禁, IMEI: %s, 剩余时间: %d分钟, 解封时间: %s", req.Device.IMEI, remainingMinutes, banUntilTime),
					fmt.Sprintf("账号已被封禁，距离解封还有%d分钟", remainingMinutes))
			}
			accountUID = account.Uid
			adSource = account.AdSource
		}
		return nil
	}); err != nil {
		return
	}

	characterList, err = db.QueryUserListByAccountName(req.Device.IMEI, req.ServerId)
	if err != nil {
		return nil, err
	}

	// create account token
	// token.DelTokenCacheFromRedis(accountUID)
	// tokenFactory := token.CreateUserFactory()
	// tokenStr, err := tokenFactory.GenerateToken(accountUID, req.Device.IMEI, "guest", req.ServerId, 3600*24*7)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "failed to generate token for account %s", req.Device.IMEI)
	// }
	// if !tokenFactory.RecordLoginToken(tokenStr, req.Device.IP) {
	// 	return nil, errors.Wrapf(errutil.ErrDBOperation, "failed to record login token for account %s", req.Device.IMEI)
	// }

	tokenStr, err := jwtv5x.Instance.GenerateToken(context.TODO(), strconv.FormatInt(accountUID, 10), "guest", req.Device.IMEI, req.ServerId, adSource)
	if err != nil {
		logger.Errorf("游客登录 token创建失败")
		return nil, errutil.ErrDBOperation
	}

	resp = &protocol.LoginResponse{
		CharacterList: []*protocol.CharacterInfo{},
		Config:        config,
		Messages:      messages,
		Token:         tokenStr,
		IsRealName:    true,
	}

	for _, c := range characterList {
		char := &protocol.CharacterInfo{
			Name:       c.Name,
			Uid:        c.Uid,
			HeadUrl:    "",
			Sex:        c.Sex,
			IP:         host,
			Port:       port,
			Coin:       c.Coin,
			PlayerIP:   req.Device.IP,
			Exp:        c.Exp,
			Diamanta:   c.Diamanta,
			Soul:       c.Soul,
			Coupon:     c.Coupon,
			Energy:     c.Energy,
			Job:        c.Job,
			Vip:        int(c.Vip),
			Level:      int(c.Level),
			Recharge:   int(c.Recharge),
			Debug:      0,
			FateLevel:  int(c.FateLevel),
			FateExp:    int(c.FateExp),
			RegisterAt: c.RegisterAt,
		}
		resp.CharacterList = append(resp.CharacterList, char)
	}

	addr := viper.GetString("webserver.addr")
	// 获取请求协议是不是http还是https
	http := ""
	enableSSL := viper.GetBool("webserver.enable_ssl")
	if enableSSL {
		http = "https://"
	} else {
		http = "http://"
	}
	UpLoginInfo(req.IMEI, http+host+":"+ipport(addr), req.ServerId)
	upUserShushu(accountUID, req.Device.OS)

	// 插入登陆记录
	device := protocol.Device{
		IP:     host,
		Remote: req.Device.IP,
	}
	db.InsertLoginLog(accountUID, device, req.AppID, req.ChannelID)

	logger.Infof("游客登录成功, account: %s, uid: %d, serverId: %d, token: %s", req.Device.IMEI, accountUID, req.ServerId, tokenStr)

	// 异步的 HTTP 回调字节
	// if req.ClickID != "" {
	// 	async.Run(func() {
	// 		logger.Debugf("GuestLogin: Sending conversion event to ByteDance for ClickID: %s", req.ClickID)

	// 		launchOptionsQuery := &protocol.LaunchOptionsSyncQuery{}
	// 		if err := json.Unmarshal([]byte(req.ClickID), launchOptionsQuery); err != nil {
	// 			logger.Warnf("GuestLogin: Failed to parse LaunchOptionsSyncQuery from ClickID: %v", err)
	// 		} else {
	// 			// 获取广告callback
	// 			var callback string
	// 			if launchOptionsQuery.Clickid != "" {
	// 				callback = launchOptionsQuery.Clickid
	// 			} else {
	// 				callback = launchOptionsQuery.ClueToken
	// 			}

	// 			if callback != "" {
	// 				// https://event-manager.oceanengine.com/docs/8650/app_api_docs
	// 				conversionEventRequest := &protocol.ConversionEventRequest{
	// 					EventType: "active_register",
	// 					Context: protocol.Context{
	// 						Ad: protocol.Ad{
	// 							Callback: callback,
	// 						},
	// 					},
	// 					Timestamp: time.Now().UnixMilli(),
	// 				}
	// 				conversionEventResponse := &protocol.ConversionEventResponse{}
	// 				resp, err := resty.New().R().
	// 					SetHeader("Content-Type", "application/json").
	// 					SetBody(conversionEventRequest).
	// 					SetResult(conversionEventResponse).
	// 					Post(constant.ByteDanceConversionURL)
	// 				if err != nil {
	// 					logger.Errorf("GuestLogin: HTTP request to ByteDance failed: %v", err)
	// 				}
	// 				if !resp.IsSuccess() {
	// 					logger.Errorf("GuestLogin: HTTP request to ByteDance failed with status code %d: %s", resp.StatusCode(), resp.String())
	// 				}
	// 				if conversionEventResponse.Code != 0 {
	// 					logger.Errorf("GuestLogin: ByteDance conversion event failed with code %d: %s", conversionEventResponse.Code, conversionEventResponse.Message)
	// 				}
	// 			}
	// 		}
	// 	})
	// }

	if err := thinkingdatax.TrackerIns.Set("", req.Device.IMEI, tracking.M{
		"account_id": req.Device.IMEI,
	}); err != nil {
		logger.Warnf("ThinkingData UserSet failed: accountID=%s, error=%v", req.Device.IMEI, err)
	}

	return resp, nil
}

func (uc *LoginUseCase) CreateCharacter(ctx context.Context, req *protocol.CreateCharacterRequest) (resp *protocol.CreateCharacterResponse, err error) {
	userData, ok := ctx.Value("userDate").(my_jwt.CustomClaims)
	if !ok {
		logger.Warnf("CreateCharacter: userData not found in context")
		return &protocol.CreateCharacterResponse{
			Code: -1,
			Date: "请求错误, 请重新登录",
		}, nil
	}
	adSource := ctx.Value("userDate").(my_jwt.CustomClaims).Ad

	var (
		user          *model.User
		characterList []model.User
		accountName   = userData.AccountName
		serverId      = userData.ServerId
		// accountID     = userData.UserId
	)

	if err = db.WithTx(ctx, func(ctx context.Context, session *xorm.Session) error {
		var err error
		characterList, err = db.QueryUserListByAccountNameV2(session, req.Name, serverId)
		if err != nil {
			return errors.Wrapf(err, "QueryUserListByAccountNameV2 failed for accountName %s on server %d", accountName, serverId)
		}
		userRoleNumber := global.GlobalCfg.BaseValue.Global.Rolenum
		if len(characterList) >= userRoleNumber {
			return errutil.NewBizError(-1,
				fmt.Sprintf("character limit exceeded, accountName: %s, serverID: %d, current: %d, limit: %d", accountName, serverId, len(characterList), userRoleNumber),
				fmt.Sprintf("最多创建%d个角色", userRoleNumber),
			)
		}

		exist, err := db.IsUserExistsByNameV2(session, req.Name)
		if err != nil {
			return errors.Wrapf(err, "IsUserExistsByNameV2 failed for name %s", req.Name)
		}
		if exist {
			return errutil.NewBizError(-1, fmt.Sprintf("user already exists: %s", req.Name), errutil.ErrAccountExists.Error())
		}

		// 创建角色
		user = &model.User{
			Uid:         snowflake.GenID(),
			Sex:         req.Sex,
			Name:        req.Name,
			Job:         req.Job,
			Level:       1,
			RegisterAt:  time.Now().Unix(),
			AccountName: accountName,
			Status:      db.StatusNormal,
			Energy:      int64(global.GlobalCfg.BaseValue.Energy.Fullvalue),
			ServerId:    serverId,
			FateLevel:   1,
			FateExp:     0,
		}
		if _, err = db.InsertV2(session, user); err != nil {
			return errors.Wrapf(err, "InsertV2 failed for name %s, account: %s, accountName: %s, serverID: %d", req.Name, req.Account, accountName, serverId)
		}

		// 创建宠物
		pet := &model.UserPet{
			Uid:         user.Uid,
			PetInfo:     "",
			PetEverHave: "",
		}
		if _, err = db.InsertV2(session, pet); err != nil {
			return errors.Wrapf(err, "InsertUserPetV2 failed for user %d", user.Uid)
		}

		// 创建宝箱
		box := &model.UserBox{
			Uid:                 user.Uid,
			NowLevel:            1,
			NowOpenRes:          "",
			NextBoxRecoveryTime: -1,
		}
		if _, err = db.InsertV2(session, box); err != nil {
			return errors.Wrapf(err, "InsertUserBoxV2 failed for user %d", user.Uid)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	characterInfo := &protocol.CharacterInfo{
		Uid:        user.Uid,
		HeadUrl:    "",
		Sex:        req.Sex,
		Name:       req.Name,
		Job:        req.Job,
		IP:         host,
		Port:       port,
		Level:      int(user.Level),
		FateLevel:  int(user.FateLevel),
		FateExp:    int(user.FateExp),
		Exp:        user.Exp,
		RegisterAt: user.RegisterAt,
	}
	if _, err = CreateRoleSendProps(user.Uid, user.Job); err != nil {
		return nil, errors.Wrapf(err, "CreateRoleSendProps failed for user %d", user.Uid)
	}

	resp = &protocol.CreateCharacterResponse{
		Code:          0,
		Date:          "",
		CharacterList: characterInfo,
	}

	ChangeRegisterCount(serverId)
	upUserCreateShushu(req.Job, len(characterList), user.Uid, accountName, user.Uid, serverId, req.Name)

	// 回调字节
	clickId := req.ClickID
	if clickId != "" && clickId != "{}" {
		dispatcher.MustGetInstance().Publish(dispatcher.TopicAdClick.String(), &dispatcher.AdClickEvent{
			OpenID:  accountName,
			ClickID: clickId,
		})
	}

	if adSource != "" {
		if err := thinkingdatax.TrackerIns.SetOnce(strconv.FormatInt(user.Uid, 10), accountName, tracking.M{
			thinkingdatax.AdSource.String(): adSource,
		}); err != nil {
			logger.Warnf("ThinkingData UserSetOnce failed: accountID=%s, error=%v", accountName, err)
		}
	}

	if err := thinkingdatax.TrackerIns.Set(strconv.FormatInt(user.Uid, 10), accountName, tracking.M{
		thinkingdatax.AccountId.String(): accountName,
		thinkingdatax.Nickname.String():  req.Name,
	}); err != nil {
		logger.Warnf("ThinkingData UserSet failed: accountID=%s, error=%v", accountName, err)
	}

	if err := thinkingdatax.TrackerIns.SetOnce(strconv.FormatInt(user.Uid, 10), accountName, tracking.M{
		thinkingdatax.FirstRegistrationTime.String(): time.Now().Format(timex.FormatDateTimeMillis),
	}); err != nil {
		logger.Warnf("ThinkingData UserSetOnce failed: accountID=%s, error=%v", accountName, err)
	}

	// 创角同步到Payment Service
	// dispatcher.MustGetInstance().Publish(dispatcher.TopicCreateCharacter.String(), &dispatcher.CreateCharacter{
	// 	AccountId:     accountID,
	// 	CharacterId:   user.Uid,
	// 	CharacterName: req.Name,
	// 	ServerId:      int64(serverId),
	// 	// ServerName:    global.ServerListMgrInstance.GetServerNameByID(serverId),
	// 	Level:       1,
	// 	Avatar:      int32(global.GlobalCfg.BaseValue.InitialIcon.Avatar),
	// 	AvatarFrame: int32(global.GlobalCfg.BaseValue.InitialIcon.Avatarframe),
	// })

	return
}
