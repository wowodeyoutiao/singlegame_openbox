package model

import (
	"nanoserver/protocol/protobuf/pb"
)

type AdminRecharge struct {
	Id           int64
	AgentId      int64  `xorm:"not null BIGINT(20) default"`
	AgentName    string `xorm:"not null VARCHAR(32) default"`
	AgentAccount string `xorm:"not null VARCHAR(32) default"`
	AdminId      string `xorm:"not null VARCHAR(32) default"`
	AdminName    string `xorm:"not null VARCHAR(32) default"`
	AdminAccount string `xorm:"not null VARCHAR(32) default"`
	Extra        string `xorm:"not null VARCHAR(255) default"`
	CreateAt     int64  `xorm:"not null BIGINT(20) default"`
	CardCount    int64  `xorm:"not null BIGINT(20) default"`
}

type Agent struct {
	Id             int64
	Name           string `xorm:"not null VARCHAR(32) default"`
	Account        string `xorm:"not null VARCHAR(32) default"`
	Password       string `xorm:"not null VARCHAR(64) default"`
	Phone          string `xorm:"not null VARCHAR(11) default"`
	Wechat         string `xorm:"not null VARCHAR(32) default"`
	Salt           string `xorm:"not null VARCHAR(32) default"`
	Role           int    `xorm:"not null TINYINT(4) default"`
	Status         int    `xorm:"not null TINYINT(4) default"`
	Extra          string `xorm:"not null VARCHAR(255) default"`
	CreateAt       int64  `xorm:"not null BIGINT(20) default"`
	DeleteAt       int64  `xorm:"not null BIGINT(20) default"`
	DeleteAccount  string `xorm:"not null VARCHAR(32) default"`
	CreateAccount  string `xorm:"not null VARCHAR(32) default"`
	ConfirmAccount string `xorm:"not null VARCHAR(32) default"`
	CardCount      int64  `xorm:"not null BIGINT(20) default"`
	Level          int    `xorm:"not null INT(20) default"`
	Discount       int    `xorm:"not null INT(20) default"`
}

type CardConsume struct {
	Id        int64
	UserId    int64  `xorm:"not null index BIGINT(20) default"`
	CardCount int    `xorm:"not null TINYINT(4) default"`
	DeskId    int64  `xorm:"not null BIGINT(20) default"`
	ClubId    int64  `xorm:"not null index BIGINT(20) default"`
	DeskNo    string `xorm:"not null VARCHAR(32) default"`
	ConsumeAt int64  `xorm:"not null BIGINT(20) default"`
	Extra     string `xorm:"not null VARCHAR(255) default"`
}

type Desk struct {
	Id           int64
	Creator      int64  `xorm:"not null index BIGINT(20) default"`
	ClubId       int64  `xorm:"not null index BIGINT(20) default"`
	Round        int    `xorm:"not null INT(11) default 8"`
	Mode         int    `xorm:"not null INT(11) default 3"`
	DeskNo       string `xorm:"not null index VARCHAR(6) default"`
	Player0      int64  `xorm:"not null index BIGINT(20) default 0"`
	Player1      int64  `xorm:"not null index BIGINT(20) default 0"`
	Player2      int64  `xorm:"not null index BIGINT(20) default 0"`
	Player3      int64  `xorm:"not null index BIGINT(20) default 0"`
	PlayerName0  string `xorm:"not null VARCHAR(255) default"`
	PlayerName1  string `xorm:"not null VARCHAR(255) default"`
	PlayerName2  string `xorm:"not null VARCHAR(255) default"`
	PlayerName3  string `xorm:"not null VARCHAR(255) default"`
	ScoreChange0 int    `xorm:"not null INT(255) default 0"`
	ScoreChange1 int    `xorm:"not null INT(255) default 0"`
	ScoreChange2 int    `xorm:"not null INT(255) default 0"`
	ScoreChange3 int    `xorm:"not null INT(255) default 0"`
	CreatedAt    int64  `xorm:"not null index BIGINT(255) default 0"`
	DismissAt    int64  `xorm:"not null BIGINT(255) default 0"`
	Extras       string `xorm:"not null TEXT default"`
}

type History struct {
	Id           int64
	DeskId       int64  `xorm:"not null index BIGINT(20) default 0"`
	Mode         int    `xorm:"not null index INT(255) default 3"`
	BeginAt      int64  `xorm:"not null BIGINT(255) default 0 "`
	EndAt        int64  `xorm:"not null BIGINT(255) default 0 "`
	PlayerName0  string `xorm:"not null VARCHAR(255) default"`
	PlayerName1  string `xorm:"not null VARCHAR(255) default"`
	PlayerName2  string `xorm:"not null VARCHAR(255) default"`
	PlayerName3  string `xorm:"not null VARCHAR(255) default"`
	ScoreChange0 int    `xorm:"not null INT(255) default 0"`
	ScoreChange1 int    `xorm:"not null INT(255) default 0"`
	ScoreChange2 int    `xorm:"not null INT(255) default 0"`
	ScoreChange3 int    `xorm:"not null INT(255) default 0"`
	Snapshot     string `xorm:"not null TEXT default"`
}

type Login struct {
	Id        int64
	Uid       int64  `xorm:"not null index BIGINT(20) default"`
	Remote    string `xorm:"not null VARCHAR(40) default"`
	Ip        string `xorm:"not null VARCHAR(40) default"`
	Model     string `xorm:"not null VARCHAR(64) default"`
	Imei      string `xorm:"not null VARCHAR(32) default"`
	Os        string `xorm:"not null VARCHAR(64) default"`
	AppId     string `xorm:"not null VARCHAR(64) default"`
	ChannelId string `xorm:"not null VARCHAR(32) default"`
	LoginAt   int64  `xorm:"not null BIGINT(11) default"`
	LogoutAt  int64  `xorm:"not null BIGINT(11) default"`
}
type LoginRole struct {
	Id       int64
	Uid      int64  `xorm:"not null index BIGINT(20) default"`
	Remote   string `xorm:"not null VARCHAR(40) default"`
	Ip       string `xorm:"not null VARCHAR(40) default"`
	Model    string `xorm:"not null VARCHAR(64) default"`
	Imei     string `xorm:"not null VARCHAR(32) default"`
	Os       string `xorm:"not null VARCHAR(64) default"`
	LoginAt  int64  `xorm:"not null BIGINT(11) default"`
	LogoutAt int64  `xorm:"not null BIGINT(11) default"`
}

type Online struct {
	Id        int64
	Time      int64 `xorm:"not null BIGINT(20) default"`
	UserCount int   `xorm:"not null INT(20) default"`
	DeskCount int   `xorm:"not null INT(11) default"`
}

// type Order struct {
// 	Id             int64
// 	OrderId        string `xorm:"not null unique VARCHAR(32)   comment('游戏内订单号')"`
// 	Type           int    `xorm:"not null TINYINT(1) default 0"`
// 	AppId          string `xorm:"not null index VARCHAR(32) default"`
// 	ChannelId      string `xorm:"not null index VARCHAR(32) default"`
// 	PayPlatform    string `xorm:"not null VARCHAR(32) default"`
// 	ChannelOrderId string `xorm:"not null VARCHAR(255) default   comment('平台订单号')"`
// 	Currency       string `xorm:"not null VARCHAR(255) default"`
// 	Extra          string `xorm:"not null VARCHAR(1024) default"`
// 	Money          int    `xorm:"not null INT(11) default   comment('价格')"`
// 	RealMoney      int    `xorm:"not null INT(11) default   comment('实际付款')"`
// 	Uid            int64  `xorm:"not null index BIGINT(20) default   comment('用户唯一ID')"`
// 	RoleId         string `xorm:"not null VARCHAR(255) default"`
// 	RoleName       string `xorm:"not null VARCHAR(255) default"`
// 	ServerId       string `xorm:"not null VARCHAR(255) default"`
// 	ServerName     string `xorm:"not null VARCHAR(255) default"`
// 	CreatedAt      int64  `xorm:"not null BIGINT(11) default"`
// 	ProductId      int    `xorm:"not null INT(11) default   comment('商品配置ID store表')"`
// 	ProductCount   int    `xorm:"not null INT(10) default"`
// 	ProductName    string `xorm:"not null VARCHAR(255) default"`
// 	ProductExtra   string `xorm:"not null VARCHAR(255) default"`
// 	NotifyUrl      string `xorm:"not null VARCHAR(2048) default"`
// 	Status         int    `xorm:"not null TINYINT(2) default 1"`
// 	Remote         string `xorm:"not null VARCHAR(40) default"`
// 	Ip             string `xorm:"not null VARCHAR(40) default"`
// 	Imei           string `xorm:"not null VARCHAR(64) default"`
// 	Os             string `xorm:"not null VARCHAR(20) default"`
// 	Model          string `xorm:"not null VARCHAR(20) default"`
// 	WeChatPrePayID string `xorm:"'wechat_prepay_id' VARCHAR(64) default ''"`
// 	OpenID         string `xorm:"'openid' VARCHAR(64) not null index(openid_nonce_mul)"`
// 	Nonce          string `xorm:"'nonce' VARCHAR(64) default '' index(openid_nonce_mul)"`
// }

// type Recharge struct {
// 	Id           int64
// 	AgentId      string `xorm:"not null VARCHAR(32) default"`
// 	AgentName    string `xorm:"not null VARCHAR(32) default"`
// 	AgentAccount string `xorm:"not null VARCHAR(32) default"`
// 	PlayerId     int64  `xorm:"not null BIGINT(20) default"`
// 	Extra        string `xorm:"not null VARCHAR(255) default"`
// 	CreateAt     int64  `xorm:"not null BIGINT(20) default"`
// 	CardCount    int64  `xorm:"not null BIGINT(20) default"`
// }

type Register struct {
	Id           int64  `xorm:"pk autoincr"`
	Uid          int64  `xorm:"not null index BIGINT(20) default"`
	Remote       string `xorm:"not null VARCHAR(40) default"`
	Ip           string `xorm:"not null VARCHAR(40) default"`
	Imei         string `xorm:"not null unique(imei_appid_serverid) VARCHAR(128) default"`
	Os           string `xorm:"not null VARCHAR(50) default"`
	Model        string `xorm:"not null VARCHAR(20) default"`
	AppId        string `xorm:"not null unique(imei_appid_serverid) VARCHAR(32) default"`
	ChannelId    string `xorm:"not null index VARCHAR(32) default"`
	RegisterAt   int64  `xorm:"not null index BIGINT(20) default"`
	RegisterType int    `xorm:"not null index TINYINT(8) default"`
	ServerId     int32  `xorm:"not null unique(imei_appid_serverid) INT(11) default 1 comment('服务器ID')"`
	BanUntil     int64  `xorm:"not null index BIGINT(20) default 0 comment('封号到期时间')"`
	AdSource     string `xorm:"not null VARCHAR(64) default '' comment('广告来源')"`
}

type ThirdAccount struct {
	Id           int64
	ThirdAccount string `xorm:"not null index VARCHAR(128) default"`
	Uid          int64  `xorm:"not null BIGINT(20) default"`
	Platform     string `xorm:"not null index VARCHAR(32) default"`
	ThirdName    string `xorm:"not null VARCHAR(64) default"`
	HeadUrl      string `xorm:"not null VARCHAR(512) default"`
	Sex          int    `xorm:"not null TINYINT(4)  default 0"`
	RealName     int    `xorm:"not null TINYINT(4)  default 0"`
	Pi           string `xorm:"not null VARCHAR(128) default"`
	Name         string `xorm:"not null VARCHAR(64) default"`
	IDCard       string `xorm:"not null VARCHAR(64) default"`
	ServerId     int32  `xorm:"not null INT(11) default 1 comment('服务器ID')"`
}

// type Trade struct {
// 	Id            int64
// 	OrderId       string `xorm:"not null unique VARCHAR(32) default"`
// 	PayOrderId    string `xorm:"not null VARCHAR(255) default"`
// 	PayPlatform   string `xorm:"not null VARCHAR(32) default"`
// 	PayAt         int64  `xorm:"not null BIGINT(11) default"`
// 	PayCreateAt   int64  `xorm:"not null BIGINT(11) default"`
// 	ComsumerId    string `xorm:"not null VARCHAR(128) default"`
// 	MerchantId    string `xorm:"not null VARCHAR(128) default"`
// 	ComsumerEmail string `xorm:"not null VARCHAR(64) default"`
// 	Raw           string `xorm:"not null VARCHAR(2048) default"`
// }

type User struct {
	Id                  int64
	Uid                 int64  `xorm:"not null index BIGINT(20) default 0 comment('角色UID')"`
	Algo                string `xorm:"not null VARCHAR(16) default"`
	Hash                string `xorm:"not null VARCHAR(64) default"`
	Salt                string `xorm:"not null VARCHAR(64) default"`
	Role                int    `xorm:"not null TINYINT(3) default 1"`
	Status              int    `xorm:"not null TINYINT(3) default 1"`
	IsOnline            int    `xorm:"not null TINYINT(1) default 1"`
	LastLoginAt         int64  `xorm:"not null index BIGINT(11) default"`
	PrivKey             string `xorm:"not null VARCHAR(512) default"`
	PubKey              string `xorm:"not null VARCHAR(128) default"`
	Coin                int64  `xorm:"not null BIGINT(20) default 0 comment('金币')"`
	Diamanta            int64  `xorm:"not null BIGINT(20) default 0 comment('钻石/血玉')"`
	Soul                int64  `xorm:"not null BIGINT(20) default 0 comment('战魂')"`
	Coupon              int64  `xorm:"not null BIGINT(20) default 0 comment('点券')"`
	Energy              int64  `xorm:"not null BIGINT(20) default 0 comment('能量')"`
	Recharge            int64  `xorm:"not null BIGINT(20) default 0 comment('充值总数')"`
	Vip                 int32  `xorm:"not null TINYINT(3) default 0 comment('VIP等级')"`
	RegisterAt          int64  `xorm:"not null index BIGINT(20) default 0"`
	FirstRechargeAt     int64  `xorm:"not null index BIGINT(20) default 0"`
	Debug               int    `xorm:"not null index TINYINT(1) default 0"`
	IsRobot             int    `xorm:"not null TINYINT(1) default 0 comment('是否机器人')"`
	Bag                 []byte `xorm:"blob null comment('包裹')"`
	Exp                 int64  `xorm:"not null BIGINT(20) default 0 comment('经验')"`
	Level               int32  `xorm:"not null INT(11) default 0 comment('等级')"`
	Job                 int    `xorm:"not null TINYINT(3) default 0 comment('职业')"`
	Name                string `xorm:"not null unique VARCHAR(30) default"`
	Sex                 int    `xorm:"not null TINYINT(3) default 0 comment('性别')"`
	Skill               string `xorm:"not null json default"`
	AccountName         string `xorm:"not null VARCHAR(128) default"`
	EnergyRecoveryTime  int64  `xorm:"not null index BIGINT(20) default"`
	Equips              string `xorm:"not null json default"`
	SkillSlot           string `xorm:"not null json default"`
	Task                []byte `xorm:"blob null comment('任务')"`
	DailySoul           int    `xorm:"not null BIGINT(20) default 0 comment('每天pvp战魂')"`
	SessionSoul         int    `xorm:"not null BIGINT(20) default 0 comment('赛季pvp战魂')"`
	PvpScore            int    `xorm:"not null BIGINT(20) default 0 comment('pvp积分')"`
	GvgScore            int    `xorm:"not null BIGINT(20) default 0 comment('gvg积分')"`
	ServerId            int32  `xorm:"not null INT(11) default 1 comment('服务器ID')"`
	Buff                []byte `xorm:"blob null comment('BUFF')"`
	SoloTowerLayers     int    `xorm:"not null INT(11) default 0 comment('爬塔层数')"`
	SoloTowerRewardTime int64  `xorm:"not null index BIGINT(20) default 0 comment('领取奖励时间')"`
	TrialVoucher        int64  `xorm:"not null BIGINT(20) default 0 comment('试炼凭证')"`
	MatchTimesDuel      int    `xorm:"not null index TINYINT(11) default 0 comment('每日比赛次数')"`
	BuyMatchNumDuel     int    `xorm:"not null index TINYINT(11) default 0 comment('每日购买比赛次数')"`
	MatchTimesMelee     int    `xorm:"not null index TINYINT(11) default 0 comment('每日比赛次数')"`
	BuyMatchNumMelee    int    `xorm:"not null index TINYINT(11) default 0 comment('每日购买比赛次数')"`
	FateExp             int64  `xorm:"not null BIGINT(20) default 0 comment('命格经验')"`
	FateLevel           int32  `xorm:"not null INT(11) default '1' comment('命格等级')"`
	WorldBossTimes      int    `xorm:"not null index INT(11) default 0 comment('世界boss每日挑战次数')"`
	WorldBossDamage     int64  `xorm:"not null index BIGINT(20) default 0 comment('世界boss每日挑战总伤害')"`
	DailyMap            string `xorm:"not null json default"`
	Disassembly         string `xorm:"not null json default"`
	Tokens              string `xorm:"not null json default"`
	TokenDeletes        string `xorm:"not null json default"`
	FirstUse            string `xorm:"not null json default"`
	DayOnlineTime       int64  `xorm:"not null index BIGINT(20) default 0 comment('当日在线时长')"`
	TotalOnlineTime     int64  `xorm:"not null index BIGINT(20) default 0 comment('总在线时长')"`
	Nowday              int32  `xorm:"not null index INT(11) default 1 comment('活跃天数')"`
	TotalCombat         int32  `xorm:"not null index BIGINT(20) default 0 comment('总战斗力')"`
	AfkExpirationAt     int64  `xorm:"not null index BIGINT(20) default 0 comment('自动挂机权限到期时间')"`
	AfkDailyTime        int64  `xorm:"not null index BIGINT(20) default 0 comment('当日挂机时长 秒')"`
	AfkMultiplierID     int64  `xorm:"not null index BIGINT(20) default -1 comment('挂机倍率ID')"`
	VitalityTime        int64  `xorm:"not null index BIGINT(20) default 0 comment('活力时间')"`
	OfflineTime         int64  `xorm:"not null index BIGINT(20) default 0 comment('离线时间')"`
	TurnBasedSlot       string `xorm:"not null json default"`
	Guidance            string `xorm:"not null json default"`
	FirstUp             string `xorm:"not null json default"`
	DailyItemUsed       string `xorm:"not null json default"`
	AdTotalToday        int32  `xorm:"not null index BIGINT(20) default 0 comment('今日广告总次数')"`
	ModuleOpen          string `xorm:"not null json default "`
	ModuleNew           []byte `xorm:"blob null comment('功能New图标激活')"`
	MuteUntil           int64  `xorm:"not null index BIGINT(20) default 0 comment('禁言到期时间')"`
}

type Uuid struct {
	Id        int64
	UidInUse  int64  `xorm:"not null index BIGINT(20) default 0"`
	UidOrigin int64  `xorm:"not null BIGINT(20) default 0"`
	Appid     string `xorm:"not null index VARCHAR(32) default"`
	Uuid      string `xorm:"not null VARCHAR(64) default"`
}

type Club struct {
	Id        int64
	Balance   int64  `xorm:"not null BIGINT(20) default 0"`
	ClubId    int64  `xorm:"not null index BIGINT(20) default 0"`
	AgentId   int64  `xorm:"not null index BIGINT(20) default 0"`
	Name      string `xorm:"not null VARCHAR(128) default"`
	Desc      string `xorm:"not null VARCHAR(512) default"`
	Member    int    `xorm:"not null INT(11) default"`
	MaxMember int    `xorm:"not null INT(11) default 500"`
	CreatedAt int64  `xorm:"not null BIGINT(20) default"`
}

type UserClub struct {
	Id        int64
	Uid       int64 `xorm:"not null index BIGINT(20) default"`
	ClubId    int64 `xorm:"not null index BIGINT(20) default"`
	CreatedAt int64 `xorm:"not null BIGINT(20) default"`
	Status    int   `xorm:"not null TINYINT(3) default 1"`
}

// mail 代表邮件结构
type Mail struct {
	Id                int64
	MailId            string       `xorm:"not null VARCHAR(128) default"`
	SenderUid         int64        `xorm:"not null index BIGINT(20) default 0 comment('发送方id')"`
	RecipientUid      int64        `xorm:"not null index BIGINT(20) default 0 comment('接收方id')"`
	Title             string       `xorm:"not null VARCHAR(128) default '-1' comment('标题')"`
	Desc              string       `xorm:"not null VARCHAR(512) default '-1' comment('内容')"`
	CreatedAt         int64        `xorm:"not null BIGINT(20) default"`
	ExpiryAt          int64        `xorm:"not null BIGINT(20) default"`
	IsOfficial        int          `xorm:"not null tinyint(1) default 0 comment('是否是官方')"`
	IsRead            int          `xorm:"not null tinyint(1) default 0 comment('是否已读')"`
	IsExtract         int          `xorm:"not null tinyint(1) default 0 comment('是否提取')"`
	IsAttachments     int          `xorm:"not null tinyint(1) default 0 comment('有没有附件')"`
	IsDeleteSender    int          `xorm:"not null tinyint(1) default 0 comment('发送者方是否删除')"`
	IsDeleteRecipient int          `xorm:"not null tinyint(1) default 0  comment('接收方是否删除')"`
	ChannelCode       string       `xorm:"not null  VARCHAR(128) default '-1' comment('渠道来源')"`
	Attachments       []Attachment `xorm:"-"`
}

// Attachment 代表邮件附件结构
type Attachment struct {
	Id          int
	MailId      string `xorm:"not null index varchar(50) default"`
	UserPropsId int64  `xorm:"not null index BIGINT(20) default"`
	Num         int64  `xorm:"not null index BIGINT(20) default"`
}

// 群发邮件的任务表
type MailTask struct {
	Id            int64
	MailTaskId    int64            `xorm:"not null BIGINT(20) default"`
	Title         string           `xorm:"not null VARCHAR(128) default '-1' comment('标题')"`
	Desc          string           `xorm:"not null VARCHAR(512) default '-1' comment('内容')"`
	IsAttachments int              `xorm:"not null tinyint(1) default 0 comment('有没有附件')"`
	CreatedAt     int64            `xorm:"not null BIGINT(20) default"`
	ExpiryAt      int64            `xorm:"not null BIGINT(20) default"`
	Attachments   []AttachmentTask `xorm:"-"`
}

// Attachment 代表邮件附件结构
type AttachmentTask struct {
	Id         int
	MailTaskId int64 `xorm:"not null index BIGINT(20) default"`
	PropsId    int64 `xorm:"not null index BIGINT(20) default"`
	Num        int64 `xorm:"not null index BIGINT(20) default"`
}

// 激活码
type Activation struct {
	Id             int64
	ActivationCode string                 `xorm:"not null varchar(255) default '-1' comment('激活码号') "`
	UseType        int                    `xorm:"not null tinyint(1) default 1 comment('1-无限使用,2-1次使用') "`
	CreatedAt      int64                  `xorm:"not null BIGINT(20) default"`
	ExpiryAt       int64                  `xorm:"not null BIGINT(20) default"`
	Attachments    []AttachmentActivation `xorm:"-"`
	Status         int                    `xorm:"not null TINYINT(3) default 1"`
}

// Attachment 代表邮件附件结构
type AttachmentActivation struct {
	Id               int
	MailActivationId string `xorm:"not null index varchar(50) default"`
	PropsId          int64  `xorm:"not null index BIGINT(20) default"`
	Num              int64  `xorm:"not null index BIGINT(20) default"`
}
type Sign struct {
	Id         int64
	Uid        int64  `xorm:"not null index BIGINT(20) default"`
	VipLevel   int32  `xorm:"not null index BIGINT(20) default 0 comment('打卡时的vip等级')"`
	Round      int32  `xorm:"not null index BIGINT(20) default 1  comment('用户目前所处的循环')" `
	Day        int32  `xorm:"not null index BIGINT(20) default 1 comment('当前的天数')"`
	UserRound  int32  `xorm:"not null index BIGINT(20) default 0 comment('已经循环几轮')"`
	Type       int    `xorm:"not null TINYINT(3) default 1"`
	CreatedAt  int64  `xorm:"not null BIGINT(20) default"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
	SourceInfo string `xorm:"not null json default"`
}
type Store struct {
	Id         int64
	Uid        int64  `xorm:"not null index BIGINT(20) default"`
	StoreType  int32  `xorm:"not null index BIGINT(20) default 0 comment('商城类型')"` //
	GoodsId    int64  `xorm:"not null index BIGINT(20) default 0 comment('物品编号')"` //
	GoodsUuid  int32  `xorm:"not null index BIGINT(20) default 0 comment('物品唯一')"` //
	SourceInfo string `xorm:"not null json default"`
	BuyNum     int32  `xorm:"not null index int(3) default 1 comment('购买次数')"` //
	CreatedAt  int64  `xorm:"not null BIGINT(20) default"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
}
type Grade struct {
	Id        int64
	Uid       int64  `xorm:"not null index BIGINT(20) default"`
	GoodsId   int64  `xorm:"not null index BIGINT(20) default 0 comment('物品编号')"`
	GoodsUuid int32  `xorm:"not null index BIGINT(20) default 0 comment('物品唯一')"`
	PropsInfo string `xorm:" json default "`
	CreatedAt int64  `xorm:"not null BIGINT(20) default"`
	Status    int    `xorm:"not null TINYINT(3) default 1"`
}
type Activity struct {
	Id         int64
	Uid        int64  `xorm:"not null index BIGINT(20) default"`
	Uuid       int64  `xorm:"not null index BIGINT(20) default"`
	Tag        int32  `xorm:"not null index BIGINT(20) default 1  comment('活动标记')" `
	Day        int32  `xorm:"not null index BIGINT(20) default 1 comment('当前的天数')"`
	Type       int    `xorm:"not null TINYINT(3) default 1 comment('1-7天,2-30天')" `
	Buy        int    `xorm:"not null TINYINT(3) default 1 comment('是否氪金')" `
	CreatedAt  int64  `xorm:"not null BIGINT(20) default 0"`
	UpdatedAt  int64  `xorm:"not null BIGINT(20) default 0"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
	SourceInfo string `xorm:"not null json default"`
}
type ActivityGet struct {
	Id        int64
	Uid       int64 `xorm:"not null index BIGINT(20) default"`
	Tag       int32 `xorm:"not null index BIGINT(20) default 1  comment('活动标记')" `
	Day       int32 `xorm:"not null index BIGINT(20) default 1 comment('当前的天数')"`
	Type      int   `xorm:"not null TINYINT(3) default 1 comment('1-7天,2-30天')" `
	Vip       int   `xorm:"not null TINYINT(3) default 1 comment('是否氪金')" `
	CreatedAt int64 `xorm:"not null BIGINT(20) default"`
	Status    int   `xorm:"not null TINYINT(3) default 1"`
}
type Privilege struct {
	Id        int64
	Uid       int64 `xorm:"not null index BIGINT(20) default"`
	Type      int   `xorm:"not null TINYINT(3) default 1 comment('1-周,2-月,3-季')"`
	CreatedAt int64 `xorm:"not null BIGINT(20) default"`
	Status    int   `xorm:"not null TINYINT(3) default 1"`
}

type RankList struct {
	Id         int64
	Type       int32  `xorm:"not null int(8) default 1 comment('1:等级榜2:战力榜3:战魂榜4:炼榜6:点的计算7:年度第一个赛季8:年度第二个赛季9:年度第三个赛季10:年度第四赛季')"`
	List       []byte `xorm:"blob null comment('排行榜内容')"`
	CreatedAt  int64  `xorm:"not null BIGINT(20) default"`
	UpdatedAt  int64  `xorm:"not null BIGINT(20) default"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
	SendStatus int    `xorm:"not null TINYINT(3) default"`
	SendAt     int64  `xorm:"not null BIGINT(20) default"`
}
type Turntable struct {
	Id        int64
	Uid       int64 `xorm:"not null index BIGINT(20) default"`
	Tag       int32 `xorm:"not null index BIGINT(20) default 1  comment('活动标记')" `
	Content   int32 `xorm:"not null int(8) default 1 comment('内容')"`
	Coupon    int32 `xorm:"not null BIGINT(20) default 0 comment('中奖数量')"`
	CreatedAt int64 `xorm:"not null BIGINT(20) default"`
	UpdatedAt int64 `xorm:"not null BIGINT(20) default"`
	Status    int   `xorm:"not null TINYINT(3) default 1"`
}
type TreasureBox struct {
	Id         int64
	Uid        int64  `xorm:"not null index BIGINT(20) default"`
	Num        int32  `xorm:"not null index TINYINT(20) default   comment('抽奖次数')"`
	SourceInfo string `xorm:"not null json default"`
	CreatedAt  int64  `xorm:"not null BIGINT(20) default"`
	UpdatedAt  int64  `xorm:"not null BIGINT(20) default"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
}
type UserExtractAccount struct {
	Id        int64
	Uid       int64  `xorm:"not null index BIGINT(20) default"`
	Type      int32  `xorm:"not null index TINYINT(20) default   comment('类型 1-微信,2-支付宝')"`
	Name      string `xorm:"not null index varchar(64) default   comment('姓名')"`
	Detail    string `xorm:"not null index varchar(64) default   comment('详情')"`
	CreatedAt int64  `xorm:"not null BIGINT(20) default"`
	UpdatedAt int64  `xorm:"not null BIGINT(20) default"`
	Status    int    `xorm:"not null TINYINT(3) default 1"`
}

// type ItemControlRefresh struct {
// 	Id          int64
// 	Uid         int32  `xorm:"not null  INT(11) default 0"`
// 	ItemCfgIds  string `xorm:"not null json default"`
// 	RecoderTime int64  `xorm:"not null index BIGINT(20) default"`
// }

// type ItemControlLast struct {
// 	Id        int64
// 	Uid       int64 `xorm:"not null  BIGINT(20) default 0"`
// 	ItemCfgId int32 `xorm:"not null INT(11) default 0 comment('道具ID')"`
// 	ItemNum   int32 `xorm:"not null INT(11) default 0 comment('道具数量')"`
// 	IsServer  int32 `xorm:"not null INT(11) default 0 comment('0server,1person')"`
// }

type ItemRecoders struct {
	Id          int64
	Uid         int64  `xorm:"not null  BIGINT(20) default 0"`
	ItemCfgId   int32  `xorm:"not null INT(11) default 0 comment('道具ID')"`
	ItemNum     int32  `xorm:"not null INT(11) default 0 comment('道具数量')"`
	RecoderTime int64  `xorm:"not null BIGINT(20) default 0"`
	Reason      string `xorm:"not null VARCHAR(30) default"`
	MapId       int32  `xorm:"not null INT(11) default 0 comment('地图ID')"`
	OPType      int32  `xorm:"not null INT(11) default 0 comment('类型 1获取 2失去')"`
}

// 伙伴
type Partner struct {
	Id              int64
	Uid             int64  `xorm:"not null  BIGINT(20) default 0"`
	PartnerUid      int64  `xorm:"not null  BIGINT(20) default 0"`
	CfgId           int32  `xorm:"not null INT(11) default 0 comment('伙伴ID')"`
	Level           int32  `xorm:"not null INT(11) default 0 comment('伙伴等级')"`
	Star            int32  `xorm:"not null INT(11) default 0 comment('伙伴星级')"`
	Skill           string `xorm:"not null json default"`
	TalnetPoints    int64  `xorm:"not null  BIGINT(20) default 0"`
	TalnetActivated string `xorm:"not null json default"`
}

// 伙伴槽
type PartnerSlot struct {
	Id           int64
	Uid          int64  `xorm:"not null  BIGINT(20) default 0"`
	SlotId       int32  `xorm:"not null INT(11) default 0 comment('伙伴槽ID')"`
	PartnerUid   int64  `xorm:"not null  BIGINT(20) default 0"`
	Equips       string `xorm:"not null json default"`
	EquipTitleid int32  `xorm:"not null index INT(11) default   comment('装备的称号')"`
}

// 宝箱
type UserBox struct {
	Id                  int64
	Uid                 int64  `xorm:"not null  BIGINT(20) default 0"`
	NowLevel            int32  `xorm:"not null INT(11) default 1 comment('当前宝箱等级')"`
	NextBoxRecoveryTime int64  `xorm:"not null BIGINT(20) default -1 comment('多少秒升级')"`
	NowOpenRes          string `xorm:"not null json default"`
	NowOpenLevel        int32  `xorm:"not null INT(11) default 0 comment('当前宝箱开出来的等级,0代表自身等级')"`
	Job                 int32  `xorm:"not null INT(11) default 0 comment('职业,0代随机')"`
	BeginAt             int64  `xorm:"not null BIGINT(20) default"`
	IsAccelerate        int32  `xorm:"not null int(3) default"`
	AccelerateNum       int32  `xorm:"not null int(8) default"`
	CostGold            int32  `xorm:"not null int(8) default"`
	TimeCost            int32  `xorm:"not null int(8) default"`
	OfflineTime         int64  `xorm:"not null BIGINT(20) default"`
}
type Summon struct {
	Id         int64
	Uid        int64  `xorm:"not null index BIGINT(20) default"`
	Module     int32  `xorm:"not null index TINYINT(4) default   comment('1-伙伴,2-宠物')"`
	Ptype      int32  `xorm:"not null index TINYINT(4) default   comment('类型 青铜白银黄金')"`
	Num        int32  `xorm:"not null index int(11) default   comment('抽奖次数')"`
	SourceInfo string `xorm:"not null json default"`
	CreatedAt  int64  `xorm:"not null BIGINT(20) default"`
	UpdatedAt  int64  `xorm:"not null BIGINT(20) default"`
	Status     int    `xorm:"not null TINYINT(3) default 1"`
}

// 宠物
type UserPet struct {
	Id                  int64
	Uid                 int64  `xorm:"not null  BIGINT(20) default 0"`
	PetRecoveryTime     int64  `xorm:"not null index BIGINT(20) default"`
	PetGold             int32  `xorm:"not null index int(11) default"`
	PetRecoveryTimeGold int64  `xorm:"not null index BIGINT(20) default"`
	PetExp              int32  `xorm:"not null index int(11) default"`
	PetRecoveryTimeExp  int64  `xorm:"not null index BIGINT(20) default"`
	PetInfo             string `xorm:"not null json default"` // 宠物信息
	PetEverHave         string `xorm:"not null json default"` // 之前拥有过的宠物
}

// 各种副本数据 都放到这个结构体中
type InstanceData struct {
	Id                int64
	Uid               int64 `xorm:"not null index BIGINT(20) default"`
	BattleCurStage    int32 `xorm:"not null index INT(11) default   comment('当前战役阶段')"`
	BattleRewardStage int32 `xorm:"not null index INT(11) default   comment('战役奖励阶段')"`

	BossCurStage      int32 `xorm:"not null index INT(11) default   comment('当前BOSS阶段')"`
	BossChallengeTime int32 `xorm:"not null index INT(11) default   comment('BOSS挑战次数')"`

	TrialTowerCurStage int32 `xorm:"not null index INT(11) default   comment('当前爬塔阶段')"`

	RelicCurStage        int32  `xorm:"not null index INT(11) default   comment('当前遗迹阶段')"`
	RelicChallengeTime   int32  `xorm:"not null index INT(11) default   comment('遗迹挑战次数')"`
	RelicEvaluation      string `xorm:"not null json default"`
	RelicChapterReward   string `xorm:"not null json default"`
	ElementCurStage      int32  `xorm:"not null index INT(11) default   comment('当前五行挑战阶段')"`
	ElementChallengeTime int32  `xorm:"not null index INT(11) default   comment('五行挑战次数')"`
	ElementNeedSweep     int32  `xorm:"not null index INT(11) default   comment('扫荡状态 0-不能扫荡 1-需要扫荡  2-已经扫荡')"`

	FightEndCurStage    int32  `xorm:"not null index INT(11) default   comment('当前一战到底阶段')"`
	FightEndEvaluation  string `xorm:"not null json default"`
	FightEndRefreshTime int64  `xorm:"not null index BIGINT(20) default"`

	ChallengeRecord string `xorm:"not null json default"`
	MainBattleId    int32  `xorm:"not null index INT(11) default   comment('主界面战斗关卡')"`

	KillMonControl string `xorm:"not null json default"`

	PushingWallCount               int32  `xorm:"'pushing_wall_count' not null INT(11) default 0 comment('今天推墙次数')"`
	PushingWallTicketPurchaseCount int32  `xorm:"'pushing_wall_ticket_purchase_count' not null INT(11) default 0 comment('今天购买推墙次数')"`
}

// 工会表
type Guild struct {
	GuildID          int64  `xorm:"'guild_id' comment('工会ID,雪花ID')"`                   // 主键自增
	Name             string `xorm:"'name' varchar(50)  comment('工会名称,唯一')"`            // 唯一名称
	Description      string `xorm:"'description' text comment('工会描述')"`                // 描述
	LeaderID         int64  `xorm:"'leader_id' index comment('工会会长ID,外键关联用户表')"`       // 外键关联用户表
	CreatedBy        int64  `xorm:"'created_by' comment('创建人ID,外键关联用户表')"`             // 创建人ID
	CreatedAt        int64  `xorm:"'created_at' not null comment('工会创建时间')"`           // 创建时间
	IconURL          string `xorm:"'icon_url' varchar(255) comment('工会图标URL')"`        // 图标URL
	MaxMembers       int32  `xorm:"'max_members' default 0 comment('工会最大成员数')"`        // 最大成员数
	AllowApplication int8   `xorm:"'allow_application' default 1 comment('是否允许申请加入')"` // 是否允许申请
	RequireApproval  int8   `xorm:"'require_approval' default 1 comment('是否需要审批加入')"`  // 是否需要审批
	Level            int    `xorm:"'level' default 1 comment('工会等级')"`                 // 工会等级
	Points           int    `xorm:"'points' default 0 comment('工会积分')"`                // 工会积分
	Status           int8   `xorm:"'status' default 1 comment('工会状态:1正常 2解散')"`        // 状态
	GuildExp         int32  `xorm:"'guild_exp' default 0 comment('工会经验')"`             // 工会经验
	Skill            string `xorm:"not null json default"`
	ServerId         int32  `xorm:"not null INT(11) default 1 comment('服务器ID')"` // 服务器ID
}

// 工会成员表（用户只能属于一个工会）
type GuildMember struct {
	MemberID          int64 `xorm:"'member_id' pk autoincr comment('工会成员ID,主键自增')"`   // 主键自增
	GuildID           int64 `xorm:"'guild_id' index comment('工会ID,外键关联工会表')"`         // 外键关联工会表
	UserID            int64 `xorm:"'user_id'  comment('用户ID,外键关联用户表')"`               // 外键关联用户表（唯一）
	JoinedAt          int64 `xorm:"'joined_at' not null comment('加入工会时间')"`           // 加入时间
	LeftAt            int64 `xorm:"'left_at' comment('退出工会时间')"`                      // 退出时间
	Role              int8  `xorm:"'role' default 3 comment('成员职位:1会长 2管理员 3成员')"`    // 职位
	Status            int8  `xorm:"'status' default 1 comment('成员状态:1正常 2退出 3被踢')"`   // 状态
	Points            int32 `xorm:"'points' default 0 comment('成员积分')"`               // 积分
	GuildContribution int32 `xorm:"'Guild_contribution' default 0 comment('成员贡献')"`   // 工会积分
	GuildWish         int32 `xorm:"'guild_wish' default 0 comment('心愿点')"`            // 心愿点
	RecoveryTimeExp   int64 `xorm:"'recovery_time_exp' default 0 comment('祈愿点恢复时间')"` // 祈愿点恢复时间
}

// 工会申请记录表
type GuildApplication struct {
	ApplicationID int64  `xorm:"'application_id' pk autoincr comment('工会申请记录ID,主键自增')"`                   // 主键自增
	GuildID       int64  `xorm:"'guild_id' index(guild_status) comment('工会ID,外键关联工会表')"`                  // 联合索引
	ApplicantID   int64  `xorm:"'applicant_id'  comment('申请用户ID,外键关联用户表')"`                               // 联合唯一索引
	AppliedAt     int64  `xorm:"'applied_at' not null comment('申请时间')"`                                   // 申请时间
	Status        int8   `xorm:"'status' default 1 index(guild_status) comment('申请状态:1待处理 2通过 3拒绝 4撤销')"` // 状态
	HandledBy     int64  `xorm:"'handled_by' null comment('处理人ID,可为空')"`                                  // 处理人（可为空）
	HandledAt     int64  `xorm:"'handled_at' comment('处理时间')"`                                            // 处理时间
	Reason        string `xorm:"'reason' text comment('申请理由')"`                                           // 申请理由
}

// 工会操作日志表
type GuildLog struct {
	LogID        int64  `xorm:"'log_id' pk autoincr comment('工会日志ID,主键自增')"`        // 主键自增
	GuildID      int64  `xorm:"'guild_id' index comment('工会ID,外键关联工会表')"`           // 外键关联工会表
	ActionType   string `xorm:"'action_type' varchar(50) not null comment('操作类型')"` // 操作类型
	OperatorID   int64  `xorm:"'operator_id' comment('操作用户ID')"`                    // 操作者用户ID
	TargetUserID int64  `xorm:"'target_user_id' null comment('目标用户ID,可为空')"`        // 目标用户ID（可为空）
	ActionTime   int64  `xorm:"'action_time' not null comment('操作时间')"`             // 操作时间
	Details      string `xorm:"'details' text comment('操作详情')"`                     // 操作详情
}

// 工会捐献表
type GuildDonate struct {
	ID         int64 `xorm:"'id' pk autoincr comment('工会捐献ID,主键自增')"`         // 主键自增
	GuildID    int64 `xorm:"'guild_id' index comment('工会ID,外键关联工会表')"`        // 外键关联工会表
	DonatorID  int64 `xorm:"'donator_id' comment('捐献者ID,外键关联用户表')"`           // 捐献者ID
	DonateType int32 `xorm:"'donate_type' default 1 comment('捐献类型:1金币 2钻石')"` // 捐献类型
	DonateNum  int32 `xorm:"'donate_num' default 0 comment('捐献数量')"`          // 捐献数量
	DonateExp  int64 `xorm:"'donate_exp' default 0 comment('累加经验')"`          // 累加经验
	DonateTime int64 `xorm:"'donate_time' not null comment('捐献时间')"`          // 捐献时间
}

// 工会商城购买记录
type GuildStoreBuy struct {
	ID      int64 `xorm:"'id' pk autoincr comment('工会商城购买记录ID,主键自增')"` // 主键自增
	UserID  int64 `xorm:"'user_id' comment('用户ID,外键关联用户表')"`           // 外键关联用户表
	GuildID int64 `xorm:"'guild_id' index comment('工会ID,外键关联工会表')"`    // 外键关联工会表
	GoodsID int64 `xorm:"'goods_id' comment('购买商品ID')"`                // 购买商品ID
	BuyNum  int32 `xorm:"'buy_num' default 0 comment('购买数量')"`         // 购买数量
	BuyTime int64 `xorm:"'buy_time' not null comment('购买时间')"`         // 购买时间
	Status  int8  `xorm:"'status' default 1 comment('状态:1正常 2取消')"`    // 状态
}

// pvp记录
type PvpRecord struct {
	ID            int64
	Uid           int64  `xorm:"not null index BIGINT(20) default 0 comment('唯一ID')"`
	RecordType    int8   `xorm:"not null TINYINT(4) default 0 comment('记录类型')"`
	CreateAt      int64  `xorm:"not null BIGINT(20) default"`
	DeleteAt      int64  `xorm:"not null BIGINT(20) default"`
	SelfUid       int64  `xorm:"not null index BIGINT(20) default 0 comment('用户ID')"`
	MatchTime     int64  `xorm:"not null BIGINT(20) default 0 comment('比赛时间')"`
	Result        int8   `xorm:"not null TINYINT(4) default 0 comment('比赛结果')"`
	OpponentUid   int64  `xorm:"not null index BIGINT(20) default 0 comment('对手用户ID')"`
	Name          string `xorm:"not null VARCHAR(50) default '' comment('名字')"`
	Power         int64  `xorm:"not null BIGINT(20) default 0 comment('战力')"`
	Gender        int8   `xorm:"not null TINYINT(4) default 0 comment('性别')"`
	Job           int8   `xorm:"not null TINYINT(4) default 0 comment('职业')"`
	Score         int64  `xorm:"not null BIGINT(20) default 0 comment('积分')"`
	Level         int32  `xorm:"not null INT(11) default 0 comment('等级')"`
	Status        int8   `xorm:"'status' default 1 comment('状态:1正常 2删除')"` // 状态
	SelfScore     int64  `xorm:"not null BIGINT(20) default 0 comment('己方积分')"`
	OpponentScore int64  `xorm:"not null BIGINT(20) default 0 comment('对手积分')"`
	AvatarV2      int32  `xorm:"not null INT(11) default 0 comment('头像V2')"`
	Frame         int32  `xorm:"not null INT(11) default 0 comment('头像框')"`
}

// 活动
type ActivityRecord struct {
	ID                 int64  `xorm:"'id' pk autoincr comment('主键自增')"`
	Uid                int64  `xorm:"not null index BIGINT(20) default 0 comment('角色UID')"`
	CumulativeRecharge []byte `xorm:"blob null comment('累计充值')"`
	OpenServerCarnival []byte `xorm:"blob null comment('开服狂欢')"`
}
type UserAttribute struct {
	Id            int64  `xorm:"'id' pk autoincr comment('主键自增')"`
	Uid           int64  `xorm:"not null index BIGINT(20) default 0 comment('角色UID')"`
	Elements      []byte `xorm:"blob null comment('元素槽位')"`
	ElementInWork []byte `xorm:"blob null comment('元素上阵')"`
	ElementMain   int32  `xorm:"not null index INT(11) default 0 comment('元素主元素')"`
}

// 战场
type War struct {
	Id           int64
	Uid          int64  `xorm:"not null  BIGINT(20) default 0"`
	Title        []byte `xorm:"blob null comment('称号')"`
	Soul         int64  `xorm:"not null BIGINT(20) default 0 comment('单日已获得战魂')"`
	EquipTitleid int32  `xorm:"not null index INT(11) default   comment('装备的称号')"`
}

// ItemDrop 代表物品掉落记录
type ItemDrop struct {
	Type            int32 `xorm:"not null unique(type_instance_id) INT(11) default 0 comment('掉落类型')"`    // ItemDropTypePlayer 1, ItemDropTypeServer 2, ItemDropTypeBattleZone 3
	InstanceId      int64 `xorm:"not null unique(type_instance_id) BIGINT(20) default 0 comment('实例ID')"` // 分别对应用户ID，服务器ID，战区ID
	ItemId          int32 `xorm:"not null unique(type_instance_id) BIGINT(20) default 0 comment('物品ID')"` // 物品ID
	MapId           int64 `xorm:"not null unique(type_instance_id) BIGINT(20) default 0 comment('地图ID')"`
	BucketStartTime int64 `xorm:"not null BIGINT(20) default 0 comment('当前桶的开始时间')"`
	CurrentCount    int32 `xorm:"not null INT(11) default 0 comment('当前桶的计数')"`
}

// 抽奖
type Lottery struct {
	Id            int64
	Uid           int64  `xorm:"not null  BIGINT(20) default 0"`
	RefreshTime   int64  `xorm:"not null BIGINT(20) default 0 comment('系统上次刷新时间')"`
	PrizePool     string `xorm:"text null comment('奖励池')"`
	HavePrizePool string `xorm:"text null comment('已有奖励池')"`
	Progress      int32  `xorm:"not null INT(11) default 0 comment('进度')"`
	NeedItemId    int32  `xorm:"not null INT(11) default 0 comment('指定道具Id')"`
	FreeNum       int32  `xorm:"not null INT(11) default 0 comment('免费次数')"`
}

// Avatar 头像头像框
type Avatar struct {
	UID           int64           `xorm:"'uid' pk comment('主键用户ID')"` // 主键用户ID
	OwnedAvatars  map[int32]int64 `xorm:"'owned_avatars' json not null comment('拥有的头像ID列表,JSON格式')"`
	OwnedFrames   map[int32]int64 `xorm:"'owned_frames' json not null comment('拥有的头像框ID列表,JSON格式')"`
	CurrentAvatar int32           `xorm:"'current_avatar' not null default 0 comment('当前使用的头像ID')"` // 当前使用的头像ID
	CurrentFrame  int32           `xorm:"'current_frame' not null default 0 comment('当前使用的头像框ID')"` // 当前使用的头像框ID
}

// Refine 洗练
type Refine struct {
	UID            int64               `xorm:"'uid' pk comment('主键用户ID')"`
	PurchasedSlots []int32             `xorm:"'purchased_slots' json not null comment('已购买槽位列表,JSON格式')"` // 已购买槽位列表
	Attrs          map[int32]*pb.Stats `xorm:"'attrs' json not null comment('洗练属性列表,JSON格式')"`            // 洗练属性列表
}

// 神兵
type Artifact struct {
	UID int64 `xorm:"'uid' pk comment('主键用户ID')"`

	Artifacts map[int32]*pb.ArtifactInfo `xorm:"'aritfacts' json not null comment('神兵信息,JSON格式')"`
}
