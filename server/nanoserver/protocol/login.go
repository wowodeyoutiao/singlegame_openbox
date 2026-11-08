package protocol

type ThirdUserLoginRequest struct {
	Platform    string `json:"platform"`     //三方平台/渠道  wechat , mobile , unity
	AppID       string `json:"appId"`        //用户来自于哪一个应用
	ChannelID   string `json:"channelId"`    //用户来自于哪一个渠道
	Device      Device `json:"device"`       //设备信息
	Name        string `json:"name"`         //微信平台名
	OpenID      string `json:"openid"`       //微信平台openid/手机号
	AccessToken string `json:"access_token"` //微信AccessToken/验证码
	ServerId    int32  `json:"serverId"`     //服务器id
}

type LoginInfo struct {
	// 三方登录字段
	Platform     string `json:"platform"`      //三方平台
	ThirdAccount string `json:"third_account"` //三方平台唯一ID
	ThirdName    string `json:"account"`       //三方平台账号名

	Token      string `json:"token"`       //用户Token
	ExpireTime int64  `json:"expire_time"` //Token过期时间

	AccountID int64 `json:"acId"` //用户的uuid,即user表的pk

	GameServerIP   string `json:"ip"` //游戏服的ip&port
	GameServerPort int    `json:"port"`
}

type UserLoginResponse struct {
	Code int32     `json:"code"` //状态码
	Data LoginInfo `json:"data"`
}

type LoginRequest struct {
	AppID     string `json:"appId"`     //用户来自于哪一个应用
	ChannelID string `json:"channelId"` //用户来自于哪一个渠道
	IMEI      string `json:"imei"`
	Device    Device `json:"device"`
	ServerId  int32  `json:"serverId"` //服务器id
	ClickID   string `json:"clickid"`  //字节的点击ID
}

type ClientConfig struct {
	Version     string `json:"version"`
	Android     string `json:"android"`
	IOS         string `json:"ios"`
	Heartbeat   int    `json:"heartbeat"`
	ForceUpdate bool   `json:"forceUpdate"`

	Title string `json:"title"` // 分享标题
	Desc  string `json:"desc"`  // 分享描述

	Daili1 string `json:"daili1"`
	Daili2 string `json:"daili2"`
	Kefu1  string `json:"kefu1"`

	AppId    string `json:"appId"`
	AppKey   string `json:"appKey"`
	AdminUrl string `json:"adminUrl"`
}

type CharacterInfo struct {
	Name       string `json:"name"`
	Uid        int64  `json:"uid"`
	HeadUrl    string `json:"headUrl"`
	Coin       int64  `json:"coin"`      //金币
	Diamanta   int64  `json:"diamanta"`  //钻石
	Coupon     int64  `json:"coupon"`    //点券
	Soul       int64  `json:"soul"`      //战魂
	Energy     int64  `json:"energy"`    //体力
	Job        int    `json:"job"`       //职业
	Sex        int    `json:"sex"`       //[0]未知 [1]男 [2]女
	Exp        int64  `json:"exp"`       //经验
	Level      int    `json:"level"`     //等级
	Vip        int    `json:"vip"`       //vip等级
	Recharge   int    `json:"recharge"`  //vip经验
	FateLevel  int    `json:"fateLevel"` //vip经验
	FateExp    int    `json:"fateExp"`   //vip经验
	RegisterAt int64  `json:"registerAt"`

	IP       string `json:"ip"`
	Port     int    `json:"port"`
	PlayerIP string `json:"playerIp"`

	Debug int `json:"debug"`
}
type LoginResponse struct {
	Code          int              `json:"code"`
	CharacterList []*CharacterInfo `json:"characterList"`
	Config        ClientConfig     `json:"config"`
	Messages      []string         `json:"messages"`
	Token         string           `json:"token"`
	IsRealName    bool             `json:"realname"` //是否已经实名认证
}

// 创建角色
type CreateCharacterRequest struct {
	Account string `json:"account"` //账号
	Sex     int    `json:"sex"`     //[0]未知 [1]男 [2]女
	Name    string `json:"name"`    //名字 N1--N2 个字符
	Job     int    `json:"job"`     //职业
	ClickID string `json:"clickid"` //字节的点击ID
}

type CreateCharacterResponse struct {
	Code          int            `json:"code"`
	Date          string         `json:"date"`
	CharacterList *CharacterInfo `json:"characterList"`
}

type SelectCharacterRequest struct {
	Uid     int64 `json:"uid"`
	BeginAt int64 `json:"begin_at"`
	EndAt   int64 `json:"end_at"`
}
type SelectCharacterResponse struct {
	Code    int    `json:"code"` //状态码
	Message string `json:"message"`
	Token   string `json:"token"`
}
type LoginrateResponse struct {
	Code    int    `json:"code"`    //状态码
	Message string `json:"message"` //备注
	Data    bool   `json:"data"`    //是否成功
}
type LoginNumResponse struct {
	Code int            `json:"code"` //状态码
	Data []LoginRoleNum `json:"data"` //是否成功
}

// 实名认证
type RealNameRequest struct {
	Name   string `json:"username"` //姓名
	IDCard string `json:"idcard"`   //身份证
}
type RealNameResponse struct {
	Code    int    `json:"code"`    //状态码
	Message string `json:"message"` //备注
}

type SMSRequest struct {
	Mobile string `json:"mobile"` //手机号
}
type SMSResponse struct {
	Code int `json:"code"` //0成功
}
