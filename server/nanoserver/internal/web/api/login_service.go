package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/internal/game"
	"nanoserver/internal/web/provide"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/jwtv5x"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/timex"
	"nanoserver/protocol"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gopkg.in/chanxuehong/wechat.v2/open/oauth2"
)

var (
	host     string                // 服务器地址
	port     int                   // 服务器端口
	config   protocol.ClientConfig // 远程配置
	messages []string              // 广播消息
	logger   = log.WithFields(log.Fields{"component": "http", "service": "login"})

	// 游客登陆
	enableGuest   = false
	guestChannels = []string{}

	enableDebug = false

	// 注册人数
	lock          sync.Mutex
	registerCount map[int32]int64 = make(map[int32]int64)
)

const defaultCoin = 10000

func AddMessage(message string) {
	messages = append(messages, message)
}

func MakeLoginService() http.Handler {
	host = viper.GetString("game-server.host")
	port = viper.GetInt("game-server.port")

	// 更新相关配置
	config.Version = viper.GetString("update.version")
	config.Android = viper.GetString("update.android")
	config.IOS = viper.GetString("update.ios")
	config.Heartbeat = viper.GetInt("core.heartbeat")

	// 分享相关配置
	config.Title = viper.GetString("share.title")
	config.Desc = viper.GetString("share.desc")

	// 客服相关配置
	config.Daili1 = viper.GetString("contact.daili1")
	config.Daili2 = viper.GetString("contact.daili2")
	config.Kefu1 = viper.GetString("contact.kefu1")

	// 游客相关配置
	enableGuest = viper.GetBool("login.guest")
	guestChannels = viper.GetStringSlice("login.lists")
	logger.Infof("是否开启游客登陆: %t, 渠道列表: %v", enableGuest, guestChannels)

	// 语音相关配置
	config.AppId = viper.GetString("voice.appid")
	config.AppKey = viper.GetString("voice.appkey")
	// 后台地址的URl
	config.AdminUrl = viper.GetString("admin.adminUrl")
	if config.Heartbeat < 5 {
		config.Heartbeat = 5
	}

	messages = viper.GetStringSlice("broadcast.message")
	logger.Debugf("version infomation: %+v", config)
	logger.Debugf("广播消息: %v", messages)

	fu := viper.GetBool("update.force")
	logger.Infof("是否强制更新: %t", fu)
	config.ForceUpdate = fu

	if register, err := db.QueryUserCountByServerId(); err == nil {
		for _, v := range register {
			registerCount[int32(v.ServerId)] = v.Num
		}
	}

	router := mux.NewRouter()
	router.Handle("/v1/user/login/query", limitHandler(nex.Handler(queryHandler), "login_query")).Methods("POST")      // 三方登录
	router.Handle("/v1/user/login/3rd", limitHandler(nex.Handler(thirdUserLoginHandler), "login_3rd")).Methods("POST") // 三方登录
	router.Handle("/v1/user/login/guest", nex.Handler(guestLoginHandler)).Methods("POST")                              // 三方登录
	router.Handle("/v1/user/club", nex.Handler(clubListHandler)).Methods("GET")                                        // 获取俱乐部列表
	router.Handle("/v1/user/login/create_char", nex.Handler(createCharacterHandler)).Methods("PUT")                    // 创建角色
	router.Handle("/v1/user/login/select_char", nex.Handler(selectCharacterHandler)).Methods("POST")                   // 选择角色
	router.Handle("/v1/user/login/del_char", nex.Handler(deleteCharacterHandler)).Methods("DELETE")                    // 角色删除
	router.Handle("/v1/user/login/realname", nex.Handler(realNameHandler)).Methods("POST")                             // 实名认证
	router.Handle("/v1/user/login/sms", nex.Handler(smsHandler)).Methods("GET")                                        // 短信

	return router
}

func ip(addr string) string {
	addr = strings.TrimSpace(addr)
	deflt := "127.0.0.1"
	if addr == "" {
		return deflt
	}

	if parts := strings.Split(addr, ":"); len(parts) > 0 {
		return parts[0]
	}

	return deflt
}

func ipport(addr string) string {
	addr = strings.TrimSpace(addr)
	deflt := "12307"
	if addr == "" {
		return deflt
	}

	if parts := strings.Split(addr, ":"); len(parts) > 0 {
		return parts[1]
	}

	return deflt
}

// 过滤emoji表情, 设置为*
func filterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		} else {
			new_content += "*"
		}
	}
	return new_content
}

func checkSession(uid int64) {
	// fixed: 之前已有session未断开
	// 检查是否该玩家是否有未断开的网络连接, 把之前的号顶掉
	// game.Kick(uid)
}

func clubs(uid int64) []protocol.ClubItem {
	list, err := db.ClubList(uid)
	if err != nil {
		return []protocol.ClubItem{}
	}

	ret := make([]protocol.ClubItem, len(list))
	for i := range list {
		ret[i] = protocol.ClubItem{
			Id:        list[i].ClubId,
			Name:      list[i].Name,
			Desc:      list[i].Desc,
			Member:    list[i].Member,
			MaxMember: list[i].MaxMember,
		}
	}
	return ret
}

func clubListHandler(form *nex.Form) (*protocol.ClubListResponse, error) {
	uid := form.Int64OrDefault("uid", -1)
	if uid < 0 {
		return nil, errors.New("服务器内部错误")
	}

	log.Println("获取俱乐部列表")

	list := clubs(uid)
	return &protocol.ClubListResponse{Data: list}, nil
}

func thirdUserLoginHandler(r *http.Request, data *protocol.ThirdUserLoginRequest) (*protocol.LoginResponse, error) {
	logger.Infof("第三方登录: %+v", data)

	if data == nil {
		return nil, errutil.ErrInvalidParameter
	}

	serverId := data.ServerId
	if serverId == 0 || data.OpenID == "" || data.AppID == "" {
		logger.Errorf("登录信息错误!!! LoginData: %v", data)
		return nil, errutil.ErrInvalidParameter
	}

	resp := &protocol.LoginResponse{CharacterList: []*protocol.CharacterInfo{}}
	var userInfo *oauth2.UserInfo
	var err error
	if data.Platform == constant.PlatformWechat {
		userInfo, err = oauth2.GetUserInfo(data.AccessToken, data.OpenID, "zh_CN", nil)
		if err != nil {
			resp.Code = -1
			return resp, nil
		}
	} else if data.Platform == constant.PlatformMobile {
		if err = CheckCode(data.OpenID, data.AccessToken); err != nil {
			resp.Code = -2
			return resp, nil
		}
	}

	var accountUid int64
	thirdUser, err := db.QueryThirdAccount(data.OpenID, data.Platform, data.ServerId)
	if err != nil && err != errutil.ErrThirdAccountNotFound {
		resp.Code = -3
		return resp, nil
	}
	var Nickname, HeadImageURL string
	var Sex int
	if data.Platform == constant.PlatformUnity || data.Platform == constant.PlatformMobile {
		Nickname = filterEmoji(thirdUser.ThirdName)
		HeadImageURL = thirdUser.HeadUrl
		Sex = thirdUser.Sex
	} else if userInfo != nil {
		Nickname = filterEmoji(userInfo.Nickname)
		HeadImageURL = userInfo.HeadImageURL
		Sex = userInfo.Sex
	}

	var characterList []model.User

	// 用户存在
	if err == nil {
		characterList, err = db.QueryUserListByAccountName(data.OpenID, serverId)
		if err != nil {
			logger.Errorf("QueryUserListByAccount faild %s  err: %s", data.OpenID, err.Error())
			resp.Code = -4
			return resp, nil
		}
		// 更新昵称
		thirdUser.ThirdName = Nickname
		thirdUser.HeadUrl = HeadImageURL
		thirdUser.Sex = Sex
		db.UpdateThirdAccount(thirdUser)

		accountUid = thirdUser.Uid

	} else {
		accountUid = snowflake.GenID()
		thirdUser = &model.ThirdAccount{
			ThirdAccount: userInfo.OpenId,
			ThirdName:    Nickname,
			Platform:     data.Platform,
			HeadUrl:      HeadImageURL,
			Sex:          Sex,
			Uid:          accountUid,
			ServerId:     serverId,
		}

		if err := db.InsertThirdAccount(thirdUser); err != nil {
			resp.Code = -5
			return resp, nil
		}
	}

	// create account token
	// token.DelTokenCacheFromRedis(accountUid)
	// tokenFactory := token.CreateUserFactory()
	// tokenStr, err := tokenFactory.GenerateToken(accountUid, data.OpenID, data.Platform, serverId, 3600*24*7)
	// if err != nil {
	// 	resp.Code = -6
	// 	return resp, nil
	// }
	// tokenFactory.RecordLoginToken(tokenStr, r.RemoteAddr)
	tokenStr, err := jwtv5x.Instance.GenerateToken(context.TODO(), strconv.FormatInt(accountUid, 10), data.Platform, data.OpenID, serverId, "")
	if err != nil {
		resp.Code = -6
		return resp, nil
	}

	resp.Config = config
	resp.Messages = messages
	resp.Token = tokenStr
	resp.IsRealName = thirdUser.RealName == 1

	for _, c := range characterList {
		char := &protocol.CharacterInfo{
			Name:      c.Name,
			Uid:       c.Uid,
			HeadUrl:   HeadImageURL,
			Sex:       c.Sex,
			IP:        host,
			Port:      port,
			Coin:      c.Coin,
			PlayerIP:  ip(r.RemoteAddr),
			Exp:       c.Exp,
			Diamanta:  c.Diamanta,
			Soul:      c.Soul,
			Coupon:    c.Coupon,
			Energy:    c.Energy,
			Job:       c.Job,
			Vip:       int(c.Vip),
			Level:     int(c.Level),
			Recharge:  int(c.Recharge),
			Debug:     0, // u.Debug,
			FateLevel: int(c.FateLevel),
			FateExp:   int(c.FateExp),
		}
		resp.CharacterList = append(resp.CharacterList, char)
	}
	// addr := viper.GetString("webserver.addr")
	// http := ""
	// enableSSL := viper.GetBool("webserver.enable_ssl")
	// if enableSSL {
	// 	http = "https://"
	// } else {
	// 	http = "http://"
	// }
	UpLoginInfo(userInfo.OpenId, host, serverId)
	// 插入登陆记录
	device := protocol.Device{
		IP:     ip(r.RemoteAddr),
		Remote: r.RemoteAddr,
		IMEI:   data.OpenID,
	}
	db.InsertLoginLog(accountUid, device, data.AppID, data.Platform)

	logger.Infof("第三方登录 account: %s platform:%s token :%s", data.OpenID, data.Platform, tokenStr)

	return resp, nil
}

func UpLoginInfo(code string, ip string, serverId int32) {
	// 上报给中台告知最近登录情况
	// reqBody := BaseRequest{
	// 	Code:   code,
	// 	Ipaddr: ip,
	// }
	// reqJson, _ := json.Marshal(&reqBody)
	// urlReq, _ := http.NewRequest(http.MethodPost, config.AdminUrl+"/api/v1/servernameuser", bytes.NewReader(reqJson))
	// parma := urlReq.URL.Query()
	// urlReq.URL.RawQuery = parma.Encode()
	// HttpDo(urlReq)
	// var resp2 respParma
	// err := json.Unmarshal(ress, &resp2)
	// if err != nil {
	// 	log.Error("解析失败", err)
	// }
	// log.Debug("vvvvvvvvvvvvvv", resp2, config.AdminUrl+"/api/v1/servernameuser")

	provide.RequestWeb(http.MethodPost, viper.GetString("admin.adminUrl")+"/api/v1/servernameuser", BaseRequest{
		Code:   code,
		Ipaddr: strconv.Itoa(int(serverId)),
	})
}

type BaseRequest struct {
	Code   string `form:"code"`
	Ipaddr string `form:"ipaddr"`
}

func HttpDo(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 根据具体返回的值填写
type respParma struct {
	RequestId string `json:"requestId"`
	Code      int32  `json:"code"`
	Msg       string `json:"msg"`
	Data      string `json:"data"`
}

func guestLoginHandler(r *http.Request, data *protocol.LoginRequest) (resp *protocol.LoginResponse, err error) {
	// 1. 早期验证和初始化
	data.Device.IMEI = data.IMEI
	data.Device.IP = ip(r.RemoteAddr)
	logger.Infof("游客登录IEMEI: %+#v", data)

	// 基础参数验证
	if data.ServerId == 0 || data.IMEI == "" || data.AppID == "" {
		logger.Errorf("登录信息错误!!! LoginData: %v", data)
		return &protocol.LoginResponse{
			Code:     -1,
			Messages: []string{"登录信息错误, 请重新登录"},
		}, nil
	}

	// 服务器ID验证
	if data.ServerId < global.ServerIDRangeStart || data.ServerId > global.ServerIDRangeEnd {
		logger.Errorf("登录信息错误, 无效的服务器ID LoginData: %v", data)
		return &protocol.LoginResponse{
			Code:     -1,
			Messages: []string{"无效的服务器ID, 请重新选择其他区服"},
		}, nil
	}

	// 业务逻辑处理
	resp, err = NewLoginUseCase().GuestLogin(context.Background(), data)
	if err != nil {
		var bizErr *errutil.BizError
		if errors.As(err, &bizErr) {
			logger.Warnf("GuestLogin failed for IMEI %s: %v", data.Device.IMEI, bizErr.Reason)
			return &protocol.LoginResponse{
				Code:     -1,
				Messages: []string{bizErr.Message},
			}, nil
		} else {
			logger.Errorf("GuestLogin failed for IMEI %s: %+v", data.Device.IMEI, err)
			return &protocol.LoginResponse{
				Code:     -1,
				Messages: []string{"服务器错误, 请稍后再试"},
			}, nil
		}
	}

	return
}

// func guestLoginHandler(r *http.Request, data *protocol.LoginRequest) (resp *protocol.LoginResponse, err error) {
// 	data.Device.IMEI = data.IMEI
// 	IP := ip(r.RemoteAddr)
// 	logger.Infof("游客登录IEMEI: %v", data)

// 	var characterList []model.User
// 	var accountUid int64
// 	serverId := data.ServerId
// 	if serverId == 0 || data.IMEI == "" || data.AppID == "" {
// 		logger.Errorf("登录信息错误!!! LoginData: %v", data)
// 		return &protocol.LoginResponse{
// 			Code:     -1,
// 			Messages: []string{"登录信息错误, 请重新登录"},
// 		}, nil
// 	}

// 	if serverId < global.ServerIDRangeStart || serverId > global.ServerIDRangeEnd {
// 		logger.Errorf("登录信息错误, 无效的服务器ID LoginData: %v", data)
// 		return &protocol.LoginResponse{
// 			Code:     -1,
// 			Messages: []string{"无效的服务器ID, 请重新选择其他区服"},
// 		}, nil
// 	}

// 	account, err := db.QueryRegister(data.AppID, data.Device.IMEI, serverId)
// 	if err != nil {
// 		if !errors.Is(err, errutil.ErrUserNotFound) {
// 			logger.Errorf("QueryRegister faild %s  err: %s", data.Device.IMEI, err.Error())
// 			return &protocol.LoginResponse{
// 				Code:     -1,
// 				Messages: []string{"服务器错误, 请稍后再试"},
// 			}, nil
// 		}
// 		if strings.ToLower(viper.GetString("core.env")) == constant.EnvProd {
// 			maxServerID := global.ServerListMgrInstance.GetMaxServerID()
// 			if maxServerID != int(data.ServerId) {
// 				logger.Errorf("游客注册失败, 服务器ID不匹配, 当前最大服务器ID: %d, 请求的服务器ID: %d", maxServerID, data.ServerId)
// 				return &protocol.LoginResponse{
// 					Code:     -1,
// 					Messages: []string{"服务器关闭注册, 请重新选择其他区服"},
// 				}, nil
// 			}
// 			// 如果不是开服的同一天, 则不允许登录
// 			// serverInfo, err := global.ServerListMgrInstance.Get(int(data.ServerId))
// 			// if err != nil {
// 			// 	logger.Errorf("获取服务器信息失败: %v, err: %v", data.ServerId, err)
// 			// 	return &protocol.LoginResponse{
// 			// 		Code:     -1,
// 			// 		Messages: []string{"服务器错误, 请稍后再试"},
// 			// 	}, nil
// 			// }
// 			// t, err := timex.ParseTimeWithFormat(serverInfo.OpenBy, timex.FormatUnixMilli)
// 			// if err != nil {
// 			// 	logger.Errorf("服务器错误: %v, err: %v", serverInfo.OpenBy, err)
// 			// 	return &protocol.LoginResponse{
// 			// 		Code:     -1,
// 			// 		Messages: []string{"服务器错误, 请稍后再试"},
// 			// 	}, nil
// 			// }
// 			// if !timex.IsSameDay(t, time.Now()) {
// 			// 	logger.Errorf("服务器 %d 不是当天开的服, 开服时间: %v", data.ServerId, t)
// 			// 	return &protocol.LoginResponse{
// 			// 		Code:     -1,
// 			// 		Messages: []string{"服务器关闭注册, 请重新选择其他区服"},
// 			// 	}, nil
// 			// }
// 		}
// 		accountUid = snowflake.GenID()
// 		db.RegisterUserLog(accountUid, data.Device, data.AppID, data.ChannelID, protocol.RegTypeThird, serverId) // 注册记录
// 		// upUserRegisterShushu(data.Device.IMEI, accountUid, data.AppID, IP)
// 	} else {
// 		accountUid = account.Uid
// 		characterList, err = db.QueryUserListByAccountName(data.Device.IMEI, serverId)
// 		if err != nil {
// 			logger.Errorf("QueryUserListByImei faild %s  err: %s", data.Device.IMEI, err.Error())
// 			return &protocol.LoginResponse{
// 				Code:     -1,
// 				Messages: []string{"服务器错误, 请稍后再试"},
// 			}, nil
// 		}
// 	}

// 	// create account token
// 	token.DelTokenCacheFromRedis(accountUid)
// 	tokenFactory := token.CreateUserFactory()
// 	tokenStr, err := tokenFactory.GenerateToken(accountUid, data.Device.IMEI, "guest", serverId, 3600*24*7)
// 	if err != nil {
// 		logger.Errorf("GenerateToken faild %s  err: %v", data.Device.IMEI, err)
// 		return &protocol.LoginResponse{
// 			Code:     -1,
// 			Messages: []string{"服务器错误, 请稍后再试"},
// 		}, nil
// 	}
// 	tokenFactory.RecordLoginToken(tokenStr, r.RemoteAddr)

// 	logger.Infof("游客登录IEMEI account: %s token :%s", data.Device.IMEI, tokenStr)

// 	resp = &protocol.LoginResponse{
// 		CharacterList: []*protocol.CharacterInfo{},
// 		Config:        config,
// 		Messages:      messages,
// 		Token:         tokenStr,
// 		IsRealName:    true,
// 	}

// 	for _, c := range characterList {
// 		char := &protocol.CharacterInfo{
// 			Name:      c.Name,
// 			Uid:       c.Uid,
// 			HeadUrl:   "",
// 			Sex:       c.Sex,
// 			IP:        host,
// 			Port:      port,
// 			Coin:      c.Coin,
// 			PlayerIP:  ip(r.RemoteAddr),
// 			Exp:       c.Exp,
// 			Diamanta:  c.Diamanta,
// 			Soul:      c.Soul,
// 			Coupon:    c.Coupon,
// 			Energy:    c.Energy,
// 			Job:       c.Job,
// 			Vip:       int(c.Vip),
// 			Level:     int(c.Level),
// 			Recharge:  int(c.Recharge),
// 			Debug:     0, // u.Debug,
// 			FateLevel: int(c.FateLevel),
// 			FateExp:   int(c.FateExp),
// 		}
// 		resp.CharacterList = append(resp.CharacterList, char)
// 	}
// 	addr := viper.GetString("webserver.addr")
// 	// 获取请求协议是不是http还是https
// 	http := ""
// 	enableSSL := viper.GetBool("webserver.enable_ssl")
// 	if enableSSL {
// 		http = "https://"
// 	} else {
// 		http = "http://"
// 	}
// 	UpLoginInfo(data.IMEI, http+host+":"+ipport(addr), serverId)
// 	upUserShushu(accountUid)

// 	// 插入登陆记录
// 	device := protocol.Device{
// 		IP:     IP,
// 		Remote: r.RemoteAddr,
// 	}
// 	db.InsertLoginLog(accountUid, device, data.AppID, data.ChannelID)

// 	return resp, nil
// }

func upUserRegisterShushu(accountId string, uid int64, device_id string, ip string) {
	data := make(map[string]interface{}, 0)
	data["name"] = "register"
	data["register_type"] = "wechat"
	data["tmp_role_id"] = uid
	data["tmp_account"] = accountId
	game.Upshushu(uid, data, false, false)

	userSet := make(map[string]interface{}, 0)
	userSet["ip"] = ip
	userSet["account_id"] = uid
	userSet["device_id"] = device_id
	userSet["register_time"] = time.Now().Unix()
	userSet["tmp_role_id"] = uid
	userSet["tmp_account"] = accountId
	game.UpshushuSetOnce(uid, userSet, false)
}

func upUserShushu(uid int64, os string) {
	data := make(map[string]interface{}, 0)
	data["name"] = "login"
	data["login_type"] = "wechat"
	data["login_time"] = time.Now().Unix()
	data["login_time_v2"] = time.Now().Format(timex.FormatDateTimeMillis)
	data["login_type"] = os
	data["tmp_role_id"] = uid
	data["tmp_account"] = fmt.Sprintf("%d", uid)
	game.Upshushu(uid, data, false, false)
}

// 查询是否使用游客登陆
type (
	queryRequest struct {
		AppId     string `json:"appId"`
		ChannelId string `json:"channelId"`
	}

	queryResponse struct {
		Code  int  `json:"code"`
		Guest bool `json:"guest"`
	}
)

var (
	forbidGuest  = &queryResponse{Guest: false}
	accepetGuest = &queryResponse{Guest: true}
)

func queryHandler(query *queryRequest) (*queryResponse, error) {
	logger.Infof("%v", query)
	if !enableGuest {
		return forbidGuest, nil
	}

	for _, s := range guestChannels {
		if query.ChannelId == s {
			return accepetGuest, nil
		}
	}

	return forbidGuest, nil
}

func createCharacterHandler(ctx context.Context, r *http.Request, data *protocol.CreateCharacterRequest) (resp *protocol.CreateCharacterResponse, err error) {
	userData, ok := ctx.Value("userDate").(my_jwt.CustomClaims)
	if !ok {
		logger.Warnf("CreateCharacter: userData not found in context")
		return &protocol.CreateCharacterResponse{
			Code: -1,
			Date: "请求错误, 请重新登录",
		}, nil
	}
	accountName := userData.AccountName
	serverId := userData.ServerId

	log.Infof("创建角色 account: %s, serverId: %d, data: %+#v", accountName, serverId, data)

	if !CanCreateCharacter(serverId) {
		logger.Warnf("当前区服 %d 已满, 不能创建角色", serverId)
		resp = &protocol.CreateCharacterResponse{Code: -2, Date: "当前区服已满,请选择其他区服"}
		return resp, nil
	}

	if err := NewPaymentClient().CheckContentSecurity(accountName, data.Name, 1); err != nil {
		logger.Warnf("内容安全检查失败: %v, accountName: %s, roleName: %s", err, accountName, data.Name)
		return &protocol.CreateCharacterResponse{Code: -1, Date: "内容含有敏感信息"}, nil
	}

	if len(data.Name) < 6 || len(data.Name) > 21 {
		logger.Warnf("角色名长度不合法: %s, accountName: %v", data.Name, accountName)
		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrRoleNameLen.Error()}
		return resp, nil
	}

	if data.Job > int(constant.Ritualist) || data.Job < int(constant.Warrior) {
		logger.Warnf("角色职业不合法: %d, accountName: %s", data.Job, accountName)
		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrInvalidParameter.Error()}
		return resp, nil
	}

	if data.Sex != 1 && data.Sex != 2 {
		logger.Warnf("角色性别不合法: %d, accountName: %s", data.Sex, accountName)
		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrIllegalParameter.Error()}
		return resp, nil
	}

	resp, err = NewLoginUseCase().CreateCharacter(ctx, data)
	if err != nil {
		var bizErr *errutil.BizError
		if errors.As(err, &bizErr) {
			logger.Warnf("CreateCharacter failed for account %s: %v", accountName, bizErr)
			return &protocol.CreateCharacterResponse{
				Code: -1,
				Date: err.Error(),
			}, nil
		} else {
			logger.Errorf("CreateCharacter failed for account %s: %+#v", accountName, err)
			return &protocol.CreateCharacterResponse{
				Code: -1,
				Date: "服务器错误, 请稍后再试",
			}, nil
		}
	}

	return resp, nil
}

// 创建角色
// func createCharacterHandler(ctx context.Context, r *http.Request, data *protocol.CreateCharacterRequest) (*protocol.CreateCharacterResponse, error) {
// 	logger.Infof("创建角色 : %+v", data)

// 	accountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName
// 	serverId := ctx.Value("userDate").(my_jwt.CustomClaims).ServerId

// 	resp := &protocol.CreateCharacterResponse{}
// 	if !CanCreateCharacter(serverId) {
// 		logger.Error("服务器创角数量已满")
// 		resp = &protocol.CreateCharacterResponse{Code: -2, Date: "当前区服已满,请选择其他区服"}
// 		return resp, nil
// 	}
// 	characterList, err := db.QueryUserListByAccountName(accountName, serverId)
// 	if err != nil {
// 		logger.Errorf("QueryUserListByImei faild %s  err: %s", accountName, err.Error())
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: err.Error()}
// 		return resp, nil
// 	}
// 	Rolenum := global.GlobalCfg.BaseValue.Global.Rolenum
// 	if len(characterList) >= Rolenum {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: "最多创建" + strconv.Itoa(Rolenum) + "个角色"}
// 		return resp, nil
// 	}

// 	if db.IsUserExistsByName(data.Name) {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrAccountExists.Error()}
// 		return resp, nil
// 	}

// 	if len(data.Name) < 6 || len(data.Name) > 21 {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrRoleNameLen.Error()}
// 		return resp, nil
// 	}

// 	if err := CheckContentSecurity(accountName, data.Name, 1); err != nil {
// 		logger.Errorf("createCharacterHandler 安全检查失败: %v", err)
// 		return &protocol.CreateCharacterResponse{
// 			Code: -1,
// 			Date: "内容含有敏感信息",
// 		}, nil
// 	}

// 	// if sensitive.ContainsSensitiveWord(data.Name) || !sensitive.IsValidName(data.Name) {
// 	// 	resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrSensitiveWord.Error()}
// 	// 	return resp, nil
// 	// }

// 	if data.Job > int(constant.Ritualist) || data.Job < int(constant.Warrior) {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrInvalidParameter.Error()}
// 		return resp, nil
// 	}

// 	if data.Sex != 1 && data.Sex != 2 {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: errutil.ErrIllegalParameter.Error()}
// 		return resp, nil
// 	}

// 	c := &model.User{}
// 	c.Uid = snowflake.GenID()
// 	c.Sex = data.Sex
// 	c.Name = data.Name
// 	c.Job = data.Job
// 	c.Level = 1
// 	c.RegisterAt = time.Now().Unix()
// 	c.AccountName = accountName
// 	c.Status = db.StatusNormal
// 	c.Energy = int64(global.GlobalCfg.BaseValue.Energy.Fullvalue)
// 	c.ServerId = serverId
// 	c.FateLevel = 1
// 	c.FateExp = 0

// 	err = db.Insert(c)
// 	if err != nil {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: err.Error()}
// 		return resp, nil
// 	}
// 	ub := &model.UserBox{}
// 	ub.Uid = c.Uid
// 	ub.NowLevel = 1
// 	ub.NowOpenRes = ""
// 	ub.NextBoxRecoveryTime = -1
// 	err = db.InsertUserBox(ub)
// 	if err != nil {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: err.Error()}
// 		return resp, nil
// 	}
// 	up := &model.UserPet{}
// 	up.Uid = c.Uid
// 	up.PetInfo = ""
// 	up.PetEverHave = ""
// 	err = db.InsertUserPet(up)
// 	if err != nil {
// 		resp = &protocol.CreateCharacterResponse{Code: -1, Date: err.Error()}
// 		return resp, nil
// 	}
// 	request := protocol.CharacterInfo{
// 		Uid:       c.Uid,
// 		HeadUrl:   "",
// 		Sex:       data.Sex,
// 		Name:      data.Name,
// 		Job:       data.Job,
// 		IP:        host,
// 		Port:      port,
// 		Level:     int(c.Level),
// 		FateLevel: int(c.FateLevel),
// 		FateExp:   int(c.FateExp),
// 		Exp:       c.Exp,
// 	}
// 	if _, err = CreateRoleSendProps(c.Uid, data.Job); err != nil {
// 		logger.Errorln("创建角色赠送装备失败", err.Error())
// 	}
// 	resp = &protocol.CreateCharacterResponse{
// 		CharacterList: &request,
// 	}
// 	ChangeRegisterCount(serverId)
// 	upUserCreateShushu(data.Job, len(characterList), c.Uid, accountName, c.Uid, serverId, data.Name)
// 	return resp, nil
// }

func upUserCreateShushu(job int, role_number int, uid int64, account_name string, distinct_id int64, serverId int32, nickname string) {
	data := make(map[string]interface{}, 0)
	data["name"] = "create_role"
	data["job"] = job
	data["tmp_account"] = account_name
	data["tmp_role_id"] = distinct_id
	data["role_number"] = role_number
	data["create_time"] = time.Now().Unix()
	game.Upshushu(uid, data, false, false)

	userSet := make(map[string]interface{}, 0)

	distinctId := strconv.FormatInt(distinct_id, 10)
	if role_number == 0 {
		userSet["tmp_account"] = account_name
		userSet["tmp_role_id"] = distinct_id
		userSet["nickname_one"] = nickname
		userSet["role_id_one"] = distinctId
		userSet["role_job_one"] = job
		userSet["server_id"] = serverId
		game.UpshushuSet(uid, userSet, false)
	}
	if role_number == 1 {
		userSet["tmp_account"] = account_name
		userSet["tmp_role_id"] = distinct_id
		userSet["nickname_two"] = nickname
		userSet["role_id_two"] = distinctId
		userSet["role_job_two"] = job
		userSet["server_id"] = serverId
		game.UpshushuSet(uid, userSet, false)
	}
	// uptencent := make(map[string]interface{}, 0)
	// uptencent["claim_type"] = "CREATE_ROLE"
	// game.Uptencent(uptencent)
}

func CreateRoleSendProps(uid int64, Type int) (bool, error) {
	channelCode := "CreateRole"
	request := protocol.MailAddRequest{}
	request.Uid = -1
	request.Title = "系统邮件"
	request.Desc = "登录送装备"
	request.RecipientUid = uid
	request.ChannelCode = "Activation-" + channelCode
	resulta := global.GlobalCfg.CreatAward

	switch constant.MowJob(Type) {
	case constant.Warrior:
		request.Attachment = make([]protocol.Attachment, len(resulta.Warrior.Item))
		for i, v2 := range resulta.Warrior.Item {
			request.Attachment[i] = protocol.Attachment{
				UserPropsId: int64(v2.ID),
				Num:         int64(v2.Num),
			}
		}
	case constant.Wizard:
		request.Attachment = make([]protocol.Attachment, len(resulta.Mage.Item))
		for i, v2 := range resulta.Mage.Item {
			request.Attachment[i] = protocol.Attachment{
				UserPropsId: int64(v2.ID),
				Num:         int64(v2.Num),
			}
		}
	default:
		request.Attachment = make([]protocol.Attachment, len(resulta.Warrior.Item))
		for i, v2 := range resulta.Warrior.Item {
			request.Attachment[i] = protocol.Attachment{
				UserPropsId: int64(v2.ID),
				Num:         int64(v2.Num),
			}
		}
	}
	if len(request.Attachment) > 0 {
		_, err := Mailadd(&request)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// 选择角色
func selectCharacterHandler(ctx context.Context, r *http.Request, data *protocol.SelectCharacterRequest) (*protocol.SelectCharacterResponse, error) {
	logger.Infof("选择登录角色 : %+v", data)

	accountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName
	platform := ctx.Value("userDate").(my_jwt.CustomClaims).Platform
	serverId := ctx.Value("userDate").(my_jwt.CustomClaims).ServerId
	adSource := ctx.Value("userDate").(my_jwt.CustomClaims).Ad
	// accountUid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	// 校验是不是自己
	check, err := db.UserCheck(data.Uid, accountName)
	if err != nil {
		logger.Errorf("UserCheck 选择登录角色失败 : %s", err.Error())
		return nil, err
	}
	if len(check) == 0 {
		// logger.Errorf("QueryUserListByImei faild %s  err: %s", accountName, err.Error())
		logger.Errorf(" 选择登录角色失败 len(check): %d", len(check))
		return nil, errutil.ErrChooseRole
		// resp := &protocol.SelectCharacterResponse{Code: -1, Message: "请选择正确账号", Token: ""}
	}
	// create user token
	// token.DelTokenCacheFromRedis(data.Uid)
	// tokenFactory := token.CreateUserFactory()
	// tokenStr, err := tokenFactory.GenerateToken(data.Uid, accountName, platform, serverId, 3600*24*7)
	// if err != nil {
	// 	logger.Errorf(" 选择登录角色失败 token创建失败")
	// 	return nil, errutil.ErrDBOperation
	// }
	// tokenFactory.RecordLoginToken(tokenStr, r.RemoteAddr)
	tokenStr, err := jwtv5x.Instance.GenerateToken(context.TODO(), strconv.FormatInt(data.Uid, 10), platform, accountName, serverId, adSource)
	if err != nil {
		logger.Errorf(" 选择登录角色失败 token创建失败")
		return nil, errutil.ErrDBOperation
	}

	// 插入登陆记录
	device := protocol.Device{
		IP:     ip(r.RemoteAddr),
		Remote: r.RemoteAddr,
	}
	db.InsertRoleLoginLog(check[0].Uid, device)
	// 走马灯
	// game.GuildLogin(check[0].Uid)
	// game.RankLogin(check[0].Uid)
	// 上报shushu
	// shushu.Upload()
	resp := &protocol.SelectCharacterResponse{Code: 0, Token: tokenStr}
	logger.Infof(" 选择登录角色uid: %d token :%s", check[0].Uid, tokenStr)
	return resp, nil
}

func loginUpdate() {
	// name := "login"
	// accountId := "1"
	// distinctId := "1"
	// properties := make(map[string]interface{})
	// properties = map[string]interface{}{

	// }
	// shushu.UploadDate(te, accountId, distinctId, name, properties)
}

// 删除角色角色
func deleteCharacterHandler(ctx context.Context, r *http.Request, data *protocol.SelectCharacterRequest) (*protocol.LoginrateResponse, error) {
	accountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName
	// 校验是不是自己
	check, err := db.UserCheck(data.Uid, accountName)
	if err != nil {
		return nil, err
	}
	if len(check) == 0 {
		// logger.Errorf("QueryUserListByImei faild %s  err: %s", accountName, err.Error())
		return nil, errutil.ErrChooseRole
		// resp := &protocol.SelectCharacterResponse{Code: -1, Message: "请选择正确账号", Token: ""}
	}

	delres, err2 := db.DelUserByAccountName(accountName, data.Uid)
	if err2 != nil || !delres {
		return nil, errutil.ErrDBOperation
	}

	// 查询角色基础信息
	characterList, err := db.QueryUserListByAccountName(accountName, check[0].ServerId)
	if err != nil {
		return nil, errutil.ErrDBOperation
	}
	upUserDeleteShushu(check[0].Job, len(characterList), data.Uid, accountName, data.Uid)
	return &protocol.LoginrateResponse{Data: true}, err2
}

func upUserDeleteShushu(job int, role_number int, uid int64, account_name string, distinct_id int64) {
	data := make(map[string]interface{}, 0)
	data["name"] = "delete_role"
	data["job"] = job
	data["role_number"] = role_number
	data["account_id"] = account_name
	data["distinct_id"] = distinct_id
	game.Upshushu(uid, data, false, false)
}

// 用户签到天数
func RoleLoginNum(ctx context.Context, data *protocol.SelectCharacterRequest) (*protocol.LoginNumResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	// 校验是不是自己
	check, err := db.RoleLoginNum(Uid, data.BeginAt, data.EndAt)
	if err != nil {
		return nil, errutil.ErrDBOperation
	}
	return &protocol.LoginNumResponse{Data: check}, nil
}

// 实名认证
type RealNameReq struct {
	Username string `json:"username"`
	Idcard   string `json:"idcard"`
}
type RealNameRes struct {
	Status int    `json:"status"`
	Pi     string `json:"pi"`
}

func realNameHandler(ctx context.Context, r *http.Request, data *protocol.RealNameRequest) (*protocol.RealNameResponse, error) {
	accountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName
	serverId := ctx.Value("userDate").(my_jwt.CustomClaims).ServerId
	res := &protocol.RealNameResponse{Code: -1, Message: errutil.ErrRealName.Error()}

	thirdUser, err := db.QueryThirdAccount(accountName, "wechat", serverId)
	if (err != nil && err != errutil.ErrThirdAccountNotFound) || thirdUser == nil {
		logger.Error(err)
		return res, nil
	}

	if thirdUser.RealName == 1 {
		return res, nil
	}

	if !VerifyIDCard(data.IDCard) {
		res.Code = -2
		res.Message = errutil.IllegalIDCard
		return res, nil
	}
	if GetAgeFromIDCard(data.IDCard) < 18 {
		res.Code = -3
		res.Message = errutil.IllegalAge
		return res, nil
	}

	if !VerifyRealname(data.Name) {
		res.Code = -4
		res.Message = errutil.IllegalName
		return res, nil
	}

	reqdata, err := json.Marshal(data)
	if err != nil {
		return res, nil
	}

	response, err := http.Post(viper.GetString("byrty0uxi.real_name_url"), "application/json", strings.NewReader(string(reqdata)))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return res, nil
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return res, nil
	}
	realNameRes := &RealNameRes{}
	if err := json.Unmarshal(body, &realNameRes); err != nil {
		logger.Errorf("unify order response unmarsharl failed: %s", err.Error())
		res.Message = err.Error()
		return res, nil
	}
	if realNameRes.Status != 0 {
		logger.Errorf("实名验证失败 Status= %d", realNameRes.Status)
		res.Code = realNameRes.Status
		res.Message = errutil.IllegalName
		return res, nil
	}

	thirdUser.Pi = realNameRes.Pi
	thirdUser.RealName = 1
	thirdUser.Name = data.Name
	thirdUser.IDCard = data.IDCard

	if err := db.UpdateThirdAccount(thirdUser); err != nil {
		return res, nil
	}

	res.Code = 0
	return res, nil
}

// VerifyRealname ...
func VerifyRealname(realname string) bool {
	regular := "^[\u4E00-\u9FA5\uf900-\ufa2d·s]{2,20}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(realname)
}

// VerifyIDCard ...
func VerifyIDCard(idcard string) bool {
	regular := "(^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$)|(^[1-9]\\d{5}\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{2}[0-9Xx]$)"
	// regular := "^[1-9]\\d{5}(18|19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"
	reg := regexp.MustCompile(regular)
	if !reg.MatchString(idcard) {
		return false
	}

	if len(idcard) == 18 {
		factor := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		parity := []int{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}
		val := []byte(idcard)

		sum := 0
		for i := 0; i < 17; i++ {
			// fmt.Println("pos:", i, "val:", int(val[i])-48)
			sum += (int(val[i]) - 48) * factor[i]
		}
		code := parity[sum%11]
		// 字符X 不需要-48
		incode := int(val[17])
		if incode != 88 {
			incode -= 48
		}
		if code != incode {
			return false
		}
	}

	checkProv := func(val string) bool {
		ival, _ := strconv.Atoi(val)
		provs := map[int]string{11: "北京", 12: "天津", 13: "河北", 14: "山西", 15: "内蒙古", 21: "辽宁", 22: "吉林", 23: "黑龙江 ", 31: "上海", 32: "江苏", 33: "浙江", 34: "安徽", 35: "福建", 36: "江西", 37: "山东", 41: "河南", 42: "湖北 ", 43: "湖南", 44: "广东", 45: "广西", 46: "海南", 50: "重庆", 51: "四川", 52: "贵州", 53: "云南", 54: "西藏 ", 61: "陕西", 62: "甘肃", 63: "青海", 64: "宁夏", 65: "新疆", 71: "台湾", 81: "香港", 82: "澳门"}
		if _, exist := provs[ival]; exist {
			return true
		}
		return false
	}
	if !checkProv(idcard[:2]) {
		return false
	}

	checkDate := func(val string) bool {
		if len(val) == 8 {
			_, err := time.ParseInLocation("20060102", val, time.Now().Location())
			if err != nil {
				return false
			}
		} else {
			_, err := time.ParseInLocation("060102", val, time.Now().Location())
			if err != nil {
				return false
			}
		}
		return true
	}
	if len(idcard) == 18 {
		if !checkDate(idcard[6:14]) {
			return false
		}
	} else {
		if !checkDate(idcard[6:12]) {
			return false
		}
	}

	return true
}

// GetAgeFromIDCard ...
func GetAgeFromIDCard(idcard string) int {
	if len(idcard) < 14 {
		return 0
	}

	dayStr := idcard[6:14]
	birthDay, err := time.ParseInLocation("20060102", dayStr, time.Now().Location())
	if err != nil {
		return 0
	}

	now := time.Now()

	ny, nm, nd := now.Date()
	by, bm, bd := birthDay.Date()

	age := ny - by
	if nm < bm ||
		(nm == bm && nd < bd) {
		age = age - 1
	}

	return age
}

type HSSms struct {
	Code string
	Time time.Time
}

// 保存验证码
var Sms_map = make(map[string]HSSms)

func smsHandler(ctx context.Context, r *http.Request, data *protocol.SMSRequest) (*protocol.SMSResponse, error) {
	// 基本参数配置
	apiUrl := "http://v.juhe.cn/vercodesms/send"
	apiKey := "2fe9565f7c1ea714e8711cb2074ae5b1"
	// 接口请求入参配置
	code := global.RandIn(8999) + 1000
	tplValue := fmt.Sprintf("#code#=%d", code)
	requestParams := url.Values{}
	requestParams.Set("key", apiKey)
	requestParams.Set("mobile", data.Mobile)
	requestParams.Set("tplId", "1")
	requestParams.Set("tplValue", tplValue)

	// 发起接口网络请求
	resp, err := http.Get(apiUrl + "?" + requestParams.Encode())
	if err != nil {
		fmt.Println("网络请求异常:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var responseResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseResult)
	if err != nil {
		fmt.Println("解析响应结果异常:", err)
		return nil, err
	}
	fmt.Println(responseResult)
	res := &protocol.SMSResponse{Code: 0}
	Sms_map[data.Mobile] = HSSms{Code: strconv.Itoa(code), Time: time.Now().Add(time.Minute * 30)}
	return res, nil
}

// 校验验证码
func CheckCode(mobile string, code string) error {
	sms, ok := Sms_map[mobile]
	now := time.Now()
	if ok && sms.Code == code {
		if now.Before(sms.Time) {
			delete(Sms_map, mobile)
			return nil
		} else {
			delete(Sms_map, mobile)
			return errutil.SMSExpire
		}
	}
	return errutil.SMSExpire
}

// 是否能注册新角色
func CanCreateCharacter(serverId int32) bool {
	lock.Lock()
	defer lock.Unlock()
	if num, ok := registerCount[serverId]; ok {
		logger.Debugf("当前人数 %d   最大人数 %d", num, global.GlobalCfg.BaseValue.Serverlimit.MaxNum)
		return num < global.GlobalCfg.BaseValue.Serverlimit.MaxNum
	}
	return true
}

// 修改人数
func ChangeRegisterCount(serverId int32) {
	lock.Lock()
	defer lock.Unlock()
	registerCount[serverId] += 1
}
