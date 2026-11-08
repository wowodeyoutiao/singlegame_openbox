package global

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"nanoserver/pkg/astar"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/randx"
	_ "nanoserver/pkg/redis_factory"
	"nanoserver/protocol/protobuf/pb"

	"nanoserver/protocol"
)

type ItemUseConfigType int8

const (
	ItemUseConfigTypeFixed       ItemUseConfigType = iota + 1
	ItemUseConfigTypeRandom                        // 随机宝箱
	ItemUseConfigTypeBuff                          // Buff
	ItemUseConfigTypeSelect                        // 选择宝箱
	ItemUseConfigTypeDirectLogic                   // 直接转换为对应的业务逻辑
)

type (
	PointF struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	PointI struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	// 全局表 config.xml
	XMLConfig struct {
		BaseValue       BaseValue             `xml:"BaseValue"`
		VIP             []Key                 `xml:"VIP>key"`
		Revive          []Revive              `xml:"Revive>key"`
		PVP             PVP                   `xml:"PVP"`
		PVPSeason       PVPSeason             `xml:"PVPseason"`
		CreatAward      CreatAward            `xml:"Creataward"`
		Rifinerandom    RefinedCostConfig     `xml:"Rifinerandom"`
		Rifinedirect    RefinedCostConfig     `xml:"Rifinedirect"`
		PVPReward       PVPReward             `xml:"PVPreward"`
		Starstone       Starstone             `xml:"starstone"`
		Switch          Switch                `xml:"switch"`
		Bounty          Bounty                `xml:"bounty"`
		GVG             GVG                   `xml:"GVG"`
		Fate            Fate                  `xml:"fate"`
		WorldBoss       WorldBoss             `xml:"worldBoss"`
		Diamond         Diamond               `xml:"Diamond"`
		DailyMaps       DailyMaps             `xml:"dailymap"`
		Initials        []Initial             `xml:"initial>key"`
		Equipenhances   []Equipenhance        `xml:"Equipenhance>quality"`
		AutoDecompose   AutoDecompose         `xml:"autodecompose>key"`
		Redbag          []RedBag              `xml:"redbag>key"`
		EquipQualitys   []EquipQuality        `xml:"Equipquality>key"`
		Partnerquality  []IDValue             `xml:"Partnerquality>key"`
		Petquality      []IDValue             `xml:"Petquality>key"`
		PetLevel        []IDValue             `xml:"Petlevel>key"`
		Petbag          []LevelValue          `xml:"Petbag>key"`
		PetAdvance      PetAdvance            `xml:"PetAdvance"`
		Copyset         []XMLIdNum            `xml:"Copyset>key"`
		Guild           Guild                 `xml:"Guild"`
		Openslot        Openslot              `xml:"openslot"`
		Partner         Partner               `xml:"partener"`
		FightEnd        FightEnd              `xml:"fightend>set"`
		GuardIsland     GuardIsland           `xml:"GuardIsland"`
		Element         Element               `xml:"Element"`
		War             War                   `xml:"War"`
		Firstgacha      IdNumAttr             `xml:"firstgacha"`
		Gacha           Gacha                 `xml:"gacha"`
		LimitRanking    LimitRanking          `xml:"LimitRanking"`
		GearRefining    GearRefining          `xml:"GearRefining"`
		MainBattle      MainBattle            `xml:"mainBattle"`
		FateCoinTrigger FateCoinTrigger       `xml:"Fatecointrigger"`
		EpicWeapon      PetAdvance            `xml:"EpicWeapon"`
		PushingWall     PushingWall           `xml:"PushingWall"`
		PartnerTalent   PartnerTalent         `xml:"Talent"`
		SameDayDelivery SameDayDeliveryConfig `xml:"SameDayDelivery"`
	}
	BaseValue struct {
		Global Global `xml:"global"`
		// Roleval          Roleval                `xml:"roleval"`
		Crit             CritConfig             `xml:"crit"`
		Power            PowerConfig            `xml:"power"`
		Energy           EnergyConfig           `xml:"energy"`
		Rename           RenameConfig           `xml:"rename"`
		Mailfreshtime    MailfreshtimeConfig    `xml:"mailfreshtime"`
		Equipstrenghthen EquipstrenghthenConfig `xml:"equipstrenghthen"`
		Chatconsume      XMLIdNum               `xml:"chatconsume"`
		Equipstar        Equipstar              `xml:"equipstar"`
		Rolelevel        Rolelevel              `xml:"rolelevel"`
		BagNum           BagNum                 `xml:"bagnum"`
		Bonus            Bonus                  `xml:"bonus"`
		Mailmax          Mailmax                `xml:"mailmax"`
		Autobattle       Autobattle             `xml:"autobattle"`
		Guidemap         Guidemap               `xml:"guidemap"`
		Autodecompose    Autodecompose          `xml:"autodecompose"`
		Onlinebattle     Onlinebattle           `xml:"onlinebattle"`
		PetbagMax        PetbagMax              `xml:"petbag"`
		Petbuy           Petbuy                 `xml:"petbuy"`
		Petfree          Petfree                `xml:"petfree"`
		Boxlevelup       Boxlevelup             `xml:"boxlevelup"`
		Expstore         Expstore               `xml:"expstore"`
		Expcopy          Expcopy                `xml:"expcopy"`
		Serverlimit      Serverlimit            `xml:"serverlimit"`
		PVPinitscore     PVPinitscore           `xml:"PVPinitscore"`
		Maxuselimit      Maxuselimit            `xml:"Maxuselimit"`
		InitialIcon      InitialIcon            `xml:"InitialIcon"`
		ValueRefine      ValueRefine            `xml:"ValueRefine"`
		FateCoin         FateCoin               `xml:"Fatecoin"`
		Extrovalue       Extrovalue             `xml:"extrovalue"`
		Offlinebox       Offlinebox             `xml:"Offlinebox"`
		Epicweapon       Epicweapon             `xml:"Epicweapon"`
		BattleSet        BattleSet              `xml:"battleset"`
		Official         Official               `xml:"Official"`
	}
	// <ValueRefine RandomID="1" RandomNum="100000" AssignID="5" AssignNum="10"  /><!--用来配置随机洗炼所需道具ID&数量以及定向洗炼所需道具ID&数量-->

	Basetalent struct {
		MaxNum int `xml:"maxnum,attr"`
	}
	PartnerTalent struct {
		Basetalent  Basetalent        `xml:"basetalent"`
		Cost        []XMLQualityIdNum `xml:"quality>stage"`
		CommonPoint []XMLQualityIdNum `xml:"commonpoint>stage"`
	}

	Epicweapon struct {
		MaxNum int `xml:"maxnum,attr"`
	}
	Official struct {
		ItemID       int32 `xml:"itemID,attr"`
		AttributeID  int32 `xml:"attributeid,attr"`
		AttributeNum int64 `xml:"attributenum,attr"`
	}
	Maxuselimit struct {
		Num int64 `xml:"num,attr"`
	}
	InitialIcon struct {
		Avatar      int32 `xml:"avatar,attr"`
		Avatarframe int32 `xml:"avatarframe,attr"`
	}
	ValueRefine struct {
		RandomID  int32 `xml:"RandomID,attr"`
		RandomNum int64 `xml:"RandomNum,attr"`
		AssignID  int32 `xml:"AssignID,attr"`
		AssignNum int64 `xml:"AssignNum,attr"`
	}

	Serverlimit struct {
		MaxNum int64 `xml:"maxnum,attr"`
	}
	PVPinitscore struct {
		Num int64 `xml:"num,attr"`
	}
	FightEnd struct {
		WeekSet   int    `xml:"weekset,attr"`
		StartTime string `xml:"starttime,attr"`
		EndTime   string `xml:"endtime,attr"`
	}

	Expcopy struct {
		Dailytime int32 `xml:"dailytime,attr"`
		Lastday   int32 `xml:"lastday,attr"`
	}

	Partner struct {
		Forge     Forge      `xml:"forge"`
		Decompose Decompose  `xml:"decompose>stage"`
		ExpItem   []XMLIdNum `xml:"exp>range"`
	}
	Forge struct {
		Free   Free              `xml:"free"`
		Charge []XMLQualityIdNum `xml:"charge"`
	}
	Free struct {
		Rate float64 `xml:"rate,attr"`
	}

	Decompose struct {
		Quality   int     `xml:"quality,attr"`
		Starlevel int     `xml:"starlevel,attr"`
		Rate      float64 `xml:"rate,attr"`
		ID        int     `xml:"id,attr"`
	}

	SameDayDeliveryConfig struct {
		Rank   SameDayDeliveryRank   `xml:"rank"`
		Limits []SameDayDeliveryItem `xml:"xianzhi>list"`
	}

	SameDayDeliveryRank struct {
		Num int64 `xml:"num,attr"`
	}

	SameDayDeliveryItem struct {
		ID  int32 `xml:"id,attr"`
		Num int64 `xml:"num,attr"`
	}

	// 用来配置修改角色名称所需要的道具和数量
	RenameConfig struct {
		Needitem int `xml:"needitem,attr"`
		Neednum  int `xml:"neednum,attr"`
	}
	// 来配置邮件的有效时间，单位为天
	MailfreshtimeConfig struct {
		Time int `xml:"time,attr"`
	}
	// 用来配置装备强化最高等级与角色等级之间的倍率
	EquipstrenghthenConfig struct {
		Rate int32 `xml:"rate,attr"`
	}
	// 用来配置基础攻击力和攻击一次消耗金币的系数
	Global struct {
		ConsumeVar float64 `xml:"consumeVar,attr"`
		BaseAttack int     `xml:"baseAttack,attr"`
		Rolenum    int     `xml:"rolenum,attr"`
	}
	// 用来配置角色的默认攻击速度和移动速度
	Roleval struct {
		AttackSpeed int64 `xml:"attackSpeed,attr"`
		MoveSpeed   int64 `xml:"moveSpeed,attr"`
	}
	Initial struct {
		Job         int32 `xml:"job,attr"`
		AttackSpeed int32 `xml:"attackspeed,attr"`
		MoveSpeed   int32 `xml:"movespeed,attr"`
	}

	RedBag struct {
		Id    int `xml:"ID,attr"`
		Value int `xml:"value,attr"`
	}
	// CritConfig 定义暴击率的配置结构体
	CritConfig struct {
		K              float64 `xml:"k,attr"`
		PreValue       int     `xml:"prevalue,attr"`
		CriticalDamage float64 `xml:"critdamage,attr"`
	}
	// PowerConfig 定义战斗力计算的配置结构体
	PowerConfig struct {
		DefValue         float64 `xml:"defvalue,attr"`
		LifeValue        float64 `xml:"lifevalue,attr"`
		AttackValue      float64 `xml:"attackvalue,attr"`
		Hitvalue         float64 `xml:"hitvalue,attr"`
		Dodgevalue       float64 `xml:"dodgevalue,attr"`
		Critvalue        float64 `xml:"critvalue,attr"`
		Toughness        float64 `xml:"toughness,attr"`
		Stunvalue        float64 `xml:"stunvalue,attr"`
		StunResistance   float64 `xml:"stun_resistance,attr"`
		Poisonvalue      float64 `xml:"poisonvalue,attr"`
		PoisonResistance float64 `xml:"poison_resistance,attr"`
		Reduce_N         float64 `xml:"reduce_N,attr"`
		Reduce_M         float64 `xml:"reduce_M,attr"`
	}
	// 体力
	EnergyConfig struct {
		Fullvalue          int `xml:"fullvalue,attr"`
		Maxvalue           int `xml:"maxvalue,attr"`
		Timeinterval       int `xml:"timeinterval,attr"`
		BuytimePeace       int `xml:"buytime_peace,attr"`
		ConsumeenergyPeace int `xml:"consumeenergy_peace,attr"`
		BuytimeWar         int `xml:"buytime_war,attr"`
		ConsumeenergyWar   int `xml:"consumeenergy_war,attr"`
	}

	//抗晕值的分母参数、毒抗值的分母参数
	Extrovalue struct {
		StunvalueEx   int `xml:"stunvalue,attr"`
		PoisonvalueEx int `xml:"poisonvalue,attr"`
		StunRound     int `xml:"stunround,attr"`
		PoisonRound   int `xml:"poisonround,attr"`
	}

	Key struct {
		VipLevel  int `xml:"vipLevel,attr"`
		NeedMoney int `xml:"needMoney,attr"`
	}
	LevelValue struct {
		Level int `xml:"Level,attr"`
		Value int `xml:"value,attr"`
	}
	// 原地复活消耗
	Revive struct {
		Time        int `xml:"time,attr"`
		NeedMoneyID int `xml:"needMoneyID,attr"`
		Num         int `xml:"num,attr"`
	}

	// PVP 定义1对1 PK玩法活动配置
	PVP struct {
		Match      Match      `xml:"match"`
		DailySet   DailySet   `xml:"dailyset"`
		TimeSet    TimeSet    `xml:"timeset"`
		LevelSet   LevelSet   `xml:"levelset"`
		Settlement Settlement `xml:"settlement"`
		Lamp       Lamp       `xml:"lamp"`
		Score      PVPScore   `xml:"score"`
	}

	PVPScore struct {
		Victory int32 `xml:"victory,attr"`
		Lose    int32 `xml:"lose,attr"`
	}

	GVGScore struct {
		Items []XMLIdWinLose `xml:"item"` // 将Items定义为Item类型的切片
	}

	Lamp struct {
		LampWinTimes int `xml:"time,attr"`
	}

	GVG struct {
		Match        GVGMatch     `xml:"match"`
		DailySet     DailySet     `xml:"dailyset"`
		TimeSet      TimeSet      `xml:"timeset"`
		LevelSet     LevelSet     `xml:"levelset"`
		Settlement   Settlement   `xml:"settlement"`
		BuffGap      BuffGap      `xml:"buffgap"`
		BuffGenerate BuffGenerate `xml:"buffgenerate"`
		Buff         Buff         `xml:"buff"`
		// Reward       []GVGReward  `xml:"reward>key"` // 将Reward定义为切片
		Relive    Relive    `xml:"relive"`
		GVGScore  GVGScore  `xml:"pvpscore"`
		GVGReward PVPReward `xml:"reward"`
		Maxround  Maxround  `xml:"Maxround"`
	}
	Maxround struct {
		Num int `xml:"num,attr"`
	}
	Fate struct {
		Basevalue Basevalue `xml:"basevalue"`
		Rate      Rate      `xml:"rate"`
	}
	PetAdvance struct {
		Basevalue Basevalue `xml:"basevalue"`
		Rate      Rate      `xml:"rate"`
	}

	Basevalue struct {
		Item []XMLIdNum `xml:"item"`
	}
	Guild struct {
		GuildMatch  GuildMatch  `xml:"match"`
		WishValue   WishValue   `xml:"wishvalue"`
		CreatCost   CreatCost   `xml:"creatcost"`
		InterLimit  InterLimit  `xml:"interlimit"`
		GoldDonate  GoldDonate  `xml:"golddonate"`
		MoneyDonate MoneyDonate `xml:"moneydonate"`
		DonateLimit DonateLimit `xml:"donatelimit"`
		GuildCopy   GuildCopy   `xml:"guildcopy"`
		GuildBoss   GuildBoss   `xml:"guildboss"`
	}
	GuildMatch struct {
		GuildNum     int `xml:"guildnum,attr"`
		ApplicantNum int `xml:"applicantnum,attr"`
	}
	WishValue struct {
		RefreshTime int `xml:"refreshtime,attr"`
		MaxValue    int `xml:"maxvalue,attr"`
		GetitemNum  int `xml:"getitemnum,attr"`
		GetSKillExp int `xml:"getskillexp,attr"`
	}
	CreatCost struct {
		ItemId  int `xml:"itemid,attr"`
		NeedNum int `xml:"neednum,attr"`
	}
	InterLimit struct {
		Time int `xml:"time,attr"`
	}
	GoldDonate struct {
		NeedItemId  int `xml:"needitemid,attr"`
		NeedItemNum int `xml:"needitemnum,attr"`
		GetItemNum  int `xml:"getitemnum,attr"`
		GetGuildExp int `xml:"getguildexp,attr"`
	}
	MoneyDonate struct {
		NeedItemId  int `xml:"needitemid,attr"`
		NeedItemNum int `xml:"needitemnum,attr"`
		GetItemNum  int `xml:"getitemnum,attr"`
		GetGuildExp int `xml:"getguildexp,attr"`
	}
	DonateLimit struct {
		DailyTime int `xml:"dailytime,attr"`
	}
	GuildCopy struct {
		Starttime string `xml:"starttime,attr"`
		Endtime   string `xml:"endtime,attr"`
		Lasttime  int    `xml:"lasttime,attr"`
	}
	GuildBoss struct {
		MaxNum int32 `xml:"maxnum,attr"`
	}

	Rate struct {
		Stage []Stage `xml:"stage"`
	}
	Stage struct {
		Rate int `xml:"rate,attr"`
		Pro  int `xml:"Pro,attr"`
	}
	BuffGap struct {
		RefreshTime int `xml:"refreshtime,attr"`
	}

	BuffGenerate struct {
		Keys []struct {
			X int `xml:"x,attr"`
			Y int `xml:"y,attr"`
		} `xml:"key"`
	}

	Buff struct {
		Keys []struct {
			ID int32 `xml:"id,attr"`
		} `xml:"key"`
	}

	GVGReward struct {
		Round int        `xml:"round,attr"`
		Items []XMLIdNum `xml:"item"` // 将Items定义为Item类型的切片
	}
	GVGMatch struct {
		MultiplePVP int `xml:"multiplePVP,attr"`
	}

	// PVPReward 定义了1vs1奖励结算的配置
	PVPReward struct {
		Win struct {
			Keys []RewardKey `xml:"key"`
		} `xml:"win"`
		Lose struct {
			Keys []RewardKey `xml:"key"`
		} `xml:"lose"`
	}

	// Key 定义了奖励的键值对
	RewardKey struct {
		ID   int `xml:"id,attr"`
		Num  int `xml:"num,attr"`
		Mark int `xml:"mark,attr"`
	}

	Match struct {
		SinglePVP     int `xml:"singlePVP,attr"`
		SecurityTimes int `xml:"securitytime,attr"`
	}

	DailySet struct {
		DailyTime int `xml:"dailytime,attr"`
		// DailyBuy      int `xml:"dailybuy,attr"`
		// BuyUseMoneyID int `xml:"buyusemoneyid,attr"`
		// NeedMoneyNum  int `xml:"needmoneynum,attr"`
	}

	TimeSet struct {
		MatchTime int `xml:"matchtime,attr"`
		WaitTime  int `xml:"waittime,attr"`
		Duration  int `xml:"duration,attr"`
	}

	LevelSet struct {
		LevelLimit int `xml:"levellimit,attr"` // 用来配置玩法中匹配玩家之间的最大等级差
	}

	Settlement struct {
		Settlement string `xml:"time,attr"`       // 结算时间
		Dailybegin string `xml:"dailybegin,attr"` // 每天开始时间
		Dailyend   string `xml:"dailyend,attr"`   // 每天结束时间
	}

	// PVPSeason 定义1对1 PK玩法赛季配置
	PVPSeason struct {
		Season []Season `xml:"season"`
	}

	// Season 定义赛季的开始和结束时间
	Season struct {
		ID        int    `xml:"id,attr"`
		StartTime string `xml:"starttime,attr"`
		EndTime   string `xml:"endtime,attr"`
	}
	//
	CreatAward struct {
		Warrior AwardItem `xml:"Warrior"`
		Mage    AwardItem `xml:"mage"`
	}
	AwardItem struct {
		Item []AwardIdNum `xml:"item"`
	}
	AwardIdNum struct {
		ID  int32 `xml:"id,attr"`
		Num int32 `xml:"num,attr"`
	}

	RefineCost struct {
		ID      int32 `xml:"ID,attr"`
		ItemID  int32 `xml:"itemid,attr"`
		ItemNum int64 `xml:"itemnum,attr"`
	}
	RefinedCostConfig struct {
		Cost []RefineCost `xml:"quality"`
	}
	XMLQualityIdNum struct {
		Quality int32 `xml:"quality,attr"`
		ID      int32 `xml:"id,attr"`
		Num     int64 `xml:"num,attr"`
	}
	XMLIdNum struct {
		ID  int32 `xml:"id,attr"`
		Num int64 `xml:"num,attr"`
	}
	XMLIdNumMaxtime struct {
		ID      int32 `xml:"id,attr"`
		Num     int64 `xml:"num,attr"`
		MaxTime int32 `xml:"maxtime,attr"`
	}
	XMLIdValue struct {
		ID    int32 `xml:"id,attr"`
		Value int64 `xml:"value,attr"`
	}
	XMLKeyValue struct {
		Key   string `xml:"key,attr"`
		Value int64  `xml:"value,attr"`
	}
	XMLIdWinLose struct {
		Id   int32 `xml:"id,attr"`
		Win  int32 `xml:"win,attr"`
		Lose int32 `xml:"lose,attr"`
	}

	XMLTime struct {
		ID   int32 `xml:"id,attr"`
		Time int32 `xml:"time,attr"`
	}

	Equipstar struct {
		Maxlevel int32 `xml:"maxlevel,attr"`
	}
	Rolelevel struct {
		Maxlevel int32 `xml:"maxlevel,attr"`
	}
	Openslot struct {
		IDLevel []XMLIDLevel `xml:"stage"`
	}
	XMLIDLevel struct {
		ID        int32 `xml:"id,attr"`
		Needlevel int32 `xml:"needlevel,attr"`
		Openday   int32 `xml:"openday,attr"`
		Minbattle int32 `xml:"minBattel,attr"`
		ServerDay int32 `xml:"severday,attr"`
	}
	Bonus struct {
		Rate float32 `xml:"rate,attr"`
	}
	Mailmax struct {
		Num int32 `xml:"num,attr"`
	}
	Autobattle struct {
		Viplevel int32 `xml:"viplevel,attr"`
	}

	Guidemap struct {
		LastTime int32 `xml:"time,attr"`
	}

	Autodecompose struct {
		Viplevel  int32 `xml:"viplevel,attr"`
		Rolelevel int32 `xml:"rolelevel,attr"`
	}
	Onlinebattle struct {
		Basetime    int32 `xml:"basetime,attr"`
		Triggertime int32 `xml:"triggertime,attr"`
		Limittime   int32 `xml:"limittime,attr"`
	}
	PetbagMax struct {
		Totlenum int32 `xml:"totlenum,attr"`
	}
	Petbuy struct {
		ItemId    int32 `xml:"itemid,attr"`
		NeedMoney int32 `xml:"needmoney,attr"`
	}
	Petfree struct {
		Recovertime int32 `xml:"recovertime,attr"`
		Maxnum      int32 `xml:"maxnum,attr"`
	}
	Boxlevelup struct {
		Itemid     int32 `xml:"itemid,attr"`
		Accelerate int32 `xml:"accelerate,attr"`
	}
	Expstore struct {
		Rate int64 `xml:"rate,attr"`
	}
	Starstone struct {
		Starlevel []Starlevel `xml:"starlevel"`
	}
	Starlevel struct {
		ID    int32 `xml:"ID,attr"`
		Level int32 `xml:"level,attr"`
		Slot  int32 `xml:"slot,attr"`
	}
	Module struct {
		ID        int32  `xml:"ID,attr"`
		Name      string `xml:"name,attr"`
		Openlevel int32  `xml:"openlevel,attr"`
		OpenDay   int32  `xml:"openday,attr"`
		MinBattel int32  `xml:"minBattel,attr"`
		Severday  int32  `xml:"severday,attr"`
	}
	Switch struct {
		Module []Module `xml:"module"`
	}
	BagNum struct {
		TotleNum int `xml:"totlenum,attr"`
		EquipNum int `xml:"equipnum,attr"`
		ItemNum  int `xml:"itemnum,attr"`
	}
	Bounty struct {
		Grade []Grade `xml:"grade"`
	}
	Grade struct {
		ID     int32 `xml:"ID,attr"`
		Weight int32 `xml:"weight,attr"`
	}

	Relive struct {
		Alivetime   int `xml:"alivetime,attr"`
		Protecttime int `xml:"protecttime,attr"`
	}

	// WorldBoss 定义了世界boss的配置结构
	WorldBoss struct {
		Blood         Blood             `xml:"blood"`
		EnterTime     EnterTime         `xml:"entertime"`
		Set           Set               `xml:"set"`
		Challengetime BossChallengetime `xml:"challengetime"`
		Accommodation Accommodation     `xml:"accommodation"`
		TurnIds       []int32           `xml:"opentime>turn"`
	}
	Diamond struct {
		Openlevel Openlevel `xml:"openlevel"`
		Master    Master    `xml:"master"`
	}
	Openlevel struct {
		Stage []OpenlevelStage `xml:"stage"`
	}
	OpenlevelStage struct {
		Num   int32 `xml:"num,attr"`
		Level int32 `xml:"level,attr"`
	}
	Master struct {
		Stage []XMLIdNum `xml:"stage"`
	}
	Set struct {
		Weekset   string `xml:"weekset,attr"`
		StartTime string `xml:"starttime,attr"`
		EndTime   string `xml:"endtime,attr"`
	}

	// Blood 定义了血条数量的结构
	Blood struct {
		Num int `xml:"num,attr"`
	}

	// EnterTime 定义了每天挑战次数的结构
	EnterTime struct {
		Num int `xml:"num,attr"`
	}

	BossChallengetime struct {
		Duration int `xml:"value,attr"`
	}

	Accommodation struct {
		Num int `xml:"num,attr"`
	}

	// 定义基础结构体，用于包含所有副本类型
	DailyMaps struct {
		GoldMap       GoldMap       `xml:"goldmap"`
		ExpMap        ExpMap        `xml:"expmap"`
		StrengthenMap StrengthenMap `xml:"strengthenstonemap"`
		SkillPageMap  SkillPageMap  `xml:"skillpagemap"`
		StarStoneMap  StarStoneMap  `xml:"starstonemap"`
		SoulStoneMap  SoulStoneMap  `xml:"soulstonemap"`
	}

	// 各个副本的配置结构体
	GoldMap struct {
		ChallengeTime ChallengeTime `xml:"challengetime"`
		BaseValue     XMLIdNum      `xml:"basevalue"`
		Score         Score         `xml:"score"`
		QuickAccess   QuickAccess   `xml:"quickaccess"`
		DoubleAccess  XMLIdNum      `xml:"doubleaccess"`
		Accommodation Accommodation `xml:"accommodation"`
		BuffGenerate  BuffGenerate  `xml:"buffgenerate"`
		Buff          Buff          `xml:"buff"`
		BuffGap       BuffGap       `xml:"buffgap"`
	}

	ExpMap struct {
		ChallengeTime ChallengeTime `xml:"challengetime"`
		BaseValue     XMLIdNum      `xml:"basevalue"`
		Score         Score         `xml:"score"`
		QuickAccess   QuickAccess   `xml:"quickaccess"`
		DoubleAccess  XMLIdNum      `xml:"doubleaccess"`
		Accommodation Accommodation `xml:"accommodation"`
		BuffGenerate  BuffGenerate  `xml:"buffgenerate"`
		Buff          Buff          `xml:"buff"`
		BuffGap       BuffGap       `xml:"buffgap"`
	}

	StrengthenMap struct {
		ChallengeTime ChallengeTime `xml:"challengetime"`
		BaseValue     XMLIdNum      `xml:"basevalue"`
		TotalCycle    TotalCycle    `xml:"totlecycle"` // 注意：XML中的标签是totlecycle，可能是个拼写错误
		QuickAccess   QuickAccess   `xml:"quickaccess"`
		DoubleAccess  XMLIdNum      `xml:"doubleaccess"`
		Accommodation Accommodation `xml:"accommodation"`
	}

	SkillPageMap struct {
		ChallengeTime ChallengeTime `xml:"challengetime"`
		BaseValue     XMLIdNum      `xml:"basevalue"`
		TotalCycle    TotalCycle    `xml:"totlecycle"` // 注意：XML中的标签是totlecycle，可能是个拼写错误
		QuickAccess   QuickAccess   `xml:"quickaccess"`
		DoubleAccess  XMLIdNum      `xml:"doubleaccess"`
		Accommodation Accommodation `xml:"accommodation"`
	}

	StarStoneMap struct {
		ChallengeTime ChallengeTime     `xml:"challengetime"`
		QuickAccess   QuickAccessStages `xml:"quickaccess"`
		DoubleAccess  XMLIdNum          `xml:"doubleaccess"`
	}

	SoulStoneMap struct {
		ChallengeTime ChallengeTime     `xml:"challengetime"`
		QuickAccess   QuickAccessStages `xml:"quickaccess"`
		DoubleAccess  XMLIdNum          `xml:"doubleaccess"`
	}

	// 辅助结构体
	ChallengeTime struct {
		Num      int `xml:"num,attr"`
		Duration int `xml:"time,attr"`
		VipLimit int `xml:"viplimit,attr"`
	}

	Score struct {
		Stages []XMLKeyValue `xml:"stage"`
	}

	QuickAccess struct {
		Score int `xml:"score,attr"`
	}

	QuickAccessStages struct {
		Stages []XMLIdNum `xml:"stage"`
	}

	TotalCycle struct {
		Num int `xml:"num,attr"`
	}

	Equipenhance struct {
		ID    int32   `xml:"ID,attr"`
		Value float32 `xml:"value,attr"`
	}
	// AutoDecompose 代表自动分解的配置结构
	AutoDecompose struct {
		Quality          int `xml:"quality,attr"`
		Level            int `xml:"level,attr"`
		QualitySwitch    int `xml:"qualityswitch,attr"`
		LevelSwitch      int `xml:"levelswitch,attr"`
		KeepBetterSwitch int `xml:"keepbetterswitch,attr"`
	}
	EquipQuality struct {
		ID     int32   `xml:"ID,attr"`
		Value  float64 `xml:"value,attr"`
		Prefix string  `xml:"prefix,attr"`
	}
	IDValue struct {
		ID    int32   `xml:"ID,attr"`
		Value float64 `xml:"value,attr"`
	}
	Element struct {
		Battle      ElementBattle     `xml:"battle"`
		MainElement MainElement       `xml:"mainelement"`
		Daily       ElementDaily      `xml:"daily"`
		Fresh       ElementFresh      `xml:"fresh"`
		Rank        ElementRank       `xml:"rank"`
		EnterLimit  ElementEnterLimit `xml:"enterlimit"`
		MaxNum      ElementMaxNum     `xml:"maxnum"`
	}

	ElementBattle struct {
		Constant       int32 `xml:"constant,attr"`
		RestraintRatio int32 `xml:"restraintratio,attr"`
	}

	MainElement struct {
		AttributeID int32 `xml:"AttributeID,attr"`
	}

	ElementDaily struct {
		InitTime int32 `xml:"inittime,attr"`
		NeedItem int32 `xml:"needitem,attr"`
	}

	ElementFresh struct {
		Day int32 `xml:"day,attr"`
	}

	ElementRank struct {
		MaxNum int32 `xml:"maxnum,attr"`
	}

	ElementEnterLimit struct {
		StartTime string `xml:"starttime,attr"`
		EndTime   string `xml:"endtime,attr"`
		LastTime  int32  `xml:"lasttime,attr"`
	}

	ElementMaxNum struct {
		CopyNum int32 `xml:"copynum,attr"`
	}

	// 战场相关的配置内容
	War struct {
		EnterLimit EnterLimitConfig `xml:"enterlimit"`
		OpenTime   []OpenTimeConfig `xml:"openday>copy"`
		//OpenDay     []int            `xml:"openday>day"`
		Protect     ProtectConfig  `xml:"protect"`
		Revive      []ReviveConfig `xml:"revive>order"`
		Autorevive  Autorevive     `xml:"autorevive"`
		Notice      Notice         `xml:"notice"`
		PK          PK             `xml:"PK"`
		Luckticket  Luckticket     `xml:"luckticket"`
		Titlestar   TitleStar      `xml:"titlestar"`
		Monstergrow Monstergrow    `xml:"monstergrow"`
	}

	Monstergrow struct {
		MonsterAttr []XMLIdNumMaxtime `xml:"monster"`
	}
	TitleStar struct {
		Maxnum int `xml:"maxnum,attr"`
	}
	Autorevive struct {
		Time int `xml:"time,attr"`
	}

	EnterLimitConfig struct {
		MaxNum    int `xml:"maxnum,attr"`
		CopyNum   int `xml:"copynum,attr"`
		Bornrange int `xml:"bornrange,attr"`
	}

	OpenTimeConfig struct {
		RoomId                 int32  `xml:"room,attr"`
		WeekDay                int32  `xml:"date,attr"`
		Mark                   int    `xml:"mark,attr"`
		StartTime              string `xml:"starttime,attr"`
		EndTime                string `xml:"endtime,attr"`
		LastTime               int    `xml:"lasttime,attr"`
		StartHour, StartMinute int
		EndHour, EndMinute     int
	}

	ProtectConfig struct {
		SafyTime int `xml:"safytime,attr"`
	}

	ReviveConfig struct {
		ID       int `xml:"id,attr"`
		NeedItem int `xml:"needitem,attr"`
		Num      int `xml:"num,attr"`
	}
	Notice struct {
		Time []XMLIdNum `xml:"time"`
	}

	PK struct {
		Killperson int `xml:"killperson,attr"`
		Daylimit   int `xml:"daylimit,attr"`
	}
	Luckticket struct {
		Item []XMLIdNum `xml:"item"`
	}

	IdNumAttr struct {
		ID  int32 `xml:"id,attr"`
		Num int32 `xml:"num,attr"`
	}
	Gacha struct {
		GachaLuckly  GachaLuckly  `xml:"gachaLuckly"`
		Guidneed     Guidneed     `xml:"guidneed"`
		IntMustOut   IntMustOut   `xml:"intMustOut"`
		RefreshGacha RefreshGacha `xml:"refreshGacha"`
		Output       Output       `xml:"output"`
	}
	Guidneed struct {
		ID int32 `xml:"id,attr"`
	}
	GachaLuckly struct {
		MaxNum int32 `xml:"maxnum,attr"`
	}
	IntMustOut struct {
		BagId int32 `xml:"bagid,attr"`
		Id    int32 `xml:"id,attr"`
		Num   int32 `xml:"num,attr"`
	}
	Output struct {
		Select []RefreshGacha `xml:"select"`
	}
	RefreshGacha struct {
		Id      int32 `xml:"id,attr"`
		UseType int32 `xml:"useType,attr"`
		Num     int32 `xml:"num,attr"`
		Luck    int32 `xml:"luck,attr"`
	}
	LimitRanking struct {
		Key []XMLIdValue `xml:"key"`
	}
	GearRefiningKey struct {
		Slot  int32 `xml:"slot,attr"`
		Limit int32 `xml:"limit,attr"`
		Type  int32 `xml:"type,attr"`
	}
	GearRefining struct {
		Keys []GearRefiningKey `xml:"key"`
	}
	FateCoin struct {
		Rate     int32 `xml:"rate,attr"`
		AwardId  int32 `xml:"awardid,attr"`
		AwardNum int32 `xml:"awardnum,attr"`
	}
	FateCoinTrigger struct {
		Key []XMLTime `xml:"key"`
	}
	PushingWall struct {
		DailyLimit struct {
			Max int32 `xml:"max,attr"`
		} `xml:"dailylimit"`
		FightTime struct {
			Max int32 `xml:"max,attr"`
		} `xml:"fightTime"`
	}
	MainBattle struct {
		Keys []struct {
			Flag    int32 `xml:"Flag,attr"`    // 怪物类型标记：0-小怪，1-精英，2-BOSS
			MinDrop int32 `xml:"MinDrop,attr"` // 最小掉落数量
			MaxDrop int32 `xml:"MaxDrop,attr"` // 最大掉落数量
		} `xml:"key"`
	}

	Offlinebox struct {
		Interval int32 `xml:"interval,attr"`
		Maxnum   int32 `xml:"maxnum,attr"`
		Equipnum int32 `xml:"equipnum,attr"`
	}

	BattleSet struct {
		Constvalue int32 `xml:"constvalue,attr"`
		SkipCD     int32 `xml:"skipCD,attr"`
	}

	//==================================================================
	//房间配置
	RoomConfig struct {
		ID            int32             // 房间ID
		Name          string            // 房间名称
		GameType      protocol.GameType // 游戏类型
		MinEntryLimit int               // 最低进入限制
		MaxEntryLimit int               // 最高进入限制
		TotalTables   int               // 房间总的桌数，为0时则为动态
		PlayerNum     int32             // 人数
		MapId         int32             // 地图ID
		Roomtype      int32             // 房间类型 room.Roomtype
		Needenerge    int32             // 进入副本消耗的体力
		Addtime       int32             // 消耗体力增加时间（秒）
	}

	IdNumP struct {
		ID    int32  `json:"id"`
		Num   int32  `json:"num"`
		Point PointI `json:"point"`
	}

	Cell struct {
		L PointI
		V string
	}
	// 地图阻挡点xml
	MapInfo struct {
		MapId  int    `xml:"ID,attr"`
		Height int    `xml:"height,attr"`
		Wide   int    `xml:"width,attr"`
		Values string `xml:"values"`
	}
	// 地图配置
	MapConfig struct {
		ID         int32
		MapID      int32
		High       int
		Wide       int
		MonsterIds []IdNumP
		MaxMonster int32
		AstarMgr   astar.Map
		StrCells   []string
		Cells      map[PointI]bool
		BornPoints []PointI
	}

	Interval struct {
		Min int32 `json:"min"`
		Max int32 `json:"max"`
	}
	DropConfig struct {
		ID    []int32  `json:"id"`
		Num   Interval `json:"num"`
		Ratio int32    `json:"ratio"`
	}
	// 怪物配置
	MonsterConfig struct {
		ID              int32
		Name            string
		Level           int32
		BassProperty    []KeyValue
		Drop            []DropConfig // 掉落
		Relive          int32        // 重生时间
		Delay           int32        // 延迟刷新时间
		MType           int32        // 0普通怪，1小boss，2，世界boss, 3采集怪
		Isdisplaybar    bool         // 是否显示血条
		Trigger         int32        // 触发事件
		Experience      int32        // 经验
		Isdisplaylevel  bool
		ViewRange       int // 巡逻范围（也是刷怪范围)
		AttackRange     int // 主动攻击范围（0,不主动攻击)
		DistributeRange int
		Ai              int
		Skills          []int // 技能列表
		AttackDistance  int
		Chase           int32          // 追逐距离
		Quality         pb.ItemQuality // 品质
	}

	Duration struct {
		Type     int32 `json:"type"` // 1：代表从道具获得按照天数进行回收的道具。2：代表按照自然天数进行回收的道具。
		Time     int32 `json:"time"`
		StarTime int64 `json:"starttime"`
		EndTime  int64 `json:"endtime"`
	}
	Stats struct {
		Key   pb.Attribute
		Value int32
	}

	ItemConfig struct {
		ID       int32
		Name     string
		MainType int32 // 类型
		// 0：代表数值类道具
		// 1：代表材料类道具
		// 2：代表装备类道具
		// 3：代表宝石类道具
		// 4：代表消耗类道具
		SonType         int32
		Maxnum          int64 // 最大堆叠数量
		Desc            string
		Icon            string
		Quality         int32    // 品质
		Priority        int32    // 优先级
		Timeliness      Duration // 道具时效性
		Iconname        string
		Equiptype       int   // 装备部位
		DiamondPosition []int // 宝石镶嵌位置
		Levellimit      int
		Occupylimit     int
		Genderlimit     int
		Stats           []Stats
		Decompose       []KeyValue
		RandomVariable  int32
		NatureMark      int32 // 自然属性标识
		Itemlevel       int32 // 物品等级
		UserNow         int32 // 是否立即使用
	}

	UUIdNum struct {
		UID int64 `json:"uid"`
		Num int64 `json:"num"`
	}
	IdNum struct {
		ID  int32 `json:"id"`
		Num int32 `json:"num"`
	}
	IdNumQuality struct {
		ID      int32 `json:"id"`
		Num     int32 `json:"num"`
		Quality int32 `json:"quality"`
	}
	IdKeyNum struct {
		ID  int32 `json:"id"`
		Num int32 `json:"num"`
		Key int32 `json:"key"`
	}
	IdNumJobWeight struct {
		ID     int32 `json:"id"`
		Num    int32 `json:"num"`
		Job    int32 `json:"job"`
		Weight int32 `json:"rateval"`
	}
	IdNumJobGender struct {
		ID     int32 `json:"id"`
		Num    int32 `json:"num"`
		Job    int32 `json:"job"`
		Gender int32 `json:"gender"`
	}

	ItemUseConfig struct {
		ID         int32  // 使用 json 标签来指定 JSON 字段名
		Name       string // 同样，这里指定 JSON 字段名
		Type       int32  // 使用类型
		GiveGoods  []IdNumJobWeight
		Buffeffect []IdKeyNum
		UseLimit   int64 // 使用限制 0无限制
	}

	RoleConfig struct {
		Level        int32
		Job          int32
		NeedExp      int64 // 所需经验
		BassProperty []KeyValue
		// WizardBassProperty  []KeyValue
	}

	SkillEffect struct {
		Duration           int32 `json:"duration"`        // 持续时间
		AttackSpeed        int32 `json:"attackspeed"`     // 攻击速度
		MovementSpeed      int32 `json:"movementspeed"`   // 移动速度
		Attack             int32 `json:"attack"`          // 攻击
		AttackPercentage   int32 `json:"attackpercent"`   // 攻击百分比
		Defense            int32 `json:"defense"`         // 防御
		DefensePercentage  int32 `json:"defensepercent"`  // 防御百分比
		Health             int32 `json:"health"`          // 生命
		HealthPercentage   int32 `json:"healthpercent"`   // 生命百分比
		CriticalStrike     int32 `json:"criticalstrike"`  // 暴击
		CriticalPercentage int32 `json:"criticalpercent"` // 暴击百分比
		FixedDamage        int32 `json:"fixeddamage"`     // 固定伤害
		DamagePercentage   int32 `json:"damagepercent"`   // 基础攻击百分比伤害
	}

	SkillUpperlevel struct {
		DamagePer   int32 `json:"DamagePer"`   // 伤害百分比提升
		DurationCon int32 `json:"DurationCon"` // 持续时间固定值增加(毫秒)
		Distance    int32 `json:"Distance"`    // 技能范围增加
		Control     int32 `json:"Control"`     // 减控百分比
		Trigger     int32 `json:"Trigger"`     // 触发几率
		Cooldown    int32 `json:"Cooldown"`    // 冷却时间减少
		SlayPer     int32 `json:"SlayPer"`     // 斩杀血量比例提升
	}

	SpecialEffect struct {
		BackDistence int32 `json:"BackDistence"` // 击退
		FlyDistence  int32 `json:"FlyDistence"`  // 瞬移
		MulAttact    int32 `json:"MulAttact"`    // 多段攻击
		DamTransfer  int32 `json:"DamTransfer"`  // 伤害转移
		Slay         int32 `json:"Slay"`         // 斩杀
		HPReducePer  int32 `json:"HPReducePer"`  // 百分比减少最大生命值
		Callout      int32 `json:"Callout"`      // 召唤
		Realdam      int32 `json:"Realdam"`      // 无视防御
	}

	SkillConfig struct {
		UID               int32      // 技能的唯一标识符
		SkillID           int32      // 技能的内部ID
		SType             int32      // 技能类型
		Level             int32      // 技能等级
		SkillProfession   int32      // 职业
		ActiveorPassive   int32      // 是否主动技能
		Name              string     // 技能名称
		RoleLevelLimit    int32      // 角色等级限制
		PromoteConsume    []IdNum    // 升级消耗
		Effect            []KeyValue // 技能效果描述
		Description       string     // 技能详细描述
		CoolDown          int32      // 冷却时间
		DelayTime         int32      // 延迟生效时间 毫秒
		MaxDistance       int32      // 技能最大距离
		Radius            int32      // 技能半径
		TriggerBuff       []int32    // 触发buffids
		CanupperLevel     int32      // 可进阶等级
		UpperlevelConsume []IdNum
		UpperlevelSetting SkillUpperlevel
		SpecialEffect     SpecialEffect

		Powerfactor   int32      // 战力影响值
		Skillbuff     []KeyValue // 技能buff效果
		Object        int32      // 作用对象
		Needaim       int32      // 是否需要目标
		Effecttrigger int32      // buff效果触发几率
		Bufftime      int32      // Buff持续时间(毫秒)
		Skilltrigger  int32      // 技能效果触发几率
		Effectnum     int32      // 作用范围类型
		Buffround     int32      // buff回合数
		RoundMark     int32      // 回合制技能标记
		RoundCD       int32      // 回合制技能释放顺序
	}

	SignConfig struct {
		Id       int32
		Round    int32
		Days     int32
		VIPLimit string
		// Contentday Award
		Content []Award
	}

	Award struct {
		Id  int32 `json:"id"`
		Num int32 `json:"num"`
	}
	StoreConfig struct {
		Id            int32
		StoreType     int32
		GoodsId       int64
		PreLimit      int32 // 前置条件
		Limit         Limit
		Props         []Props
		Entrance      string
		Method        string
		Consumechange []KeyValue
	}
	Limit struct {
		NeedMoney int32 `json:"needmoney"` // 是否需要人民币
		LimitType int64 `json:"limittype"` // 0：代表终生限制；1代表每日限制；2代表每周限制；3代表每月限制
		LimitNum  int64 `json:"limitnum"`  // 购买次数限制
		Num       int64 `json:"num"`       // 总数量
		Times     int32 `json:"times"`     // 购买多少次免费
	}
	Props struct {
		Id  int32 `json:"itemID"`
		Num int32 `json:"Num"`
	}
	KeyValue struct {
		Key   int32 `json:"key"`
		Value int32 `json:"value"`
	}

	KeyValues struct {
		Key    int32 `json:"key"`
		Values int32 `json:"values"`
	}
	AttrWithJob struct {
		Key   int32 `json:"key"`
		Value int32 `json:"value"`
		Job   int32 `json:"job"`
	}
	IdValues struct {
		Id     int32 `json:"id"`
		Values int32 `json:"values"`
	}
	KeyValueKind struct {
		Kind  int32 `json:"kind"`
		Key   int32 `json:"key"`
		Value int32 `json:"value"`
	}
	KeyValueEx struct {
		Id     int32      `json:"id"`
		Values []KeyValue `json:"values"`
	}
	KeyValueEx2 struct {
		Key    int32   `json:"key"`
		Values []IdNum `json:"values"`
	}
	KeyValueGap struct {
		Id     int32          `json:"id"`
		Values []KeyValueKind `json:"values"`
	}
	FightEndInfo struct {
		Uid           int64 `json:"uid"`
		BattleId      int32 `json:"battleid"`
		ChallengeTime int32 `json:"challengetime"`
	}
	ElementInfo struct {
		BattleId      int32 `json:"battleid"`
		ChallengeTime int32 `json:"challengetime"`
	}
	ActivityConfig struct {
		UUID                int64      // 编号
		Type                int32      // 活动类型
		Content             int32      // 内容编号
		Acttag              int32      // 活动标记
		CommonItem          []KeyValue // 普通道具
		PayItem             []KeyValue // 付费道具
		PreLimit            int32      // 前置条件
		CompensationConsume []KeyValue // 补领消耗
		ActivityTime        string     // 活动时间
		RuleDes             string     // 活动规则
		Name                string     // 活动名称
		ResetTime           int32      // 重置时间
		StartTime           int32      // 开启时间
	}

	BuffConfig struct {
		ID          int32      `json:"ID"`          // Buff的ID
		Type        int32      `json:"Type"`        // Buff的类型
		Mark        int32      `json:"Mark"`        // Buff标识
		Describe    string     `json:"Describe"`    // 说明
		Value       []KeyValue `json:"Value"`       // 数值配置，JSON格式
		Rate        int32      `json:"Rate"`        // 触发几率
		LastTime    int32      `json:"LastTime"`    // 持续时间(毫秒)
		CDTime      int32      `json:"CDTime"`      // 冷却时间(毫秒)
		EffectTime  int32      `json:"EffectTime"`  // 生效间隔
		Icon        string     `json:"Icon"`        // Buff图标
		BuffEffect  string     `json:"BuffEffect"`  // Buff特效
		CanPickUp   int32      `json:"CanPickUp"`   // 是否可拾取
		EffectRange int32      `json:"EffectRange"` // 作用范围
		DamIncrease int32      `json:"DamIncrease"` // 伤害增量
		StackNum    int32      `json:"StackNum"`    // 叠加层数
		Pretype     int32      `json:"Pretype"`     // 大类
	}

	StartEnd struct {
		Start string `json:"start"` // 唯一标识
		End   string `json:"end"`   // 唯一标识
	}

	TaskConfig struct {
		UUID            int32      `json:"uuid"`            // 唯一标识
		Class           int32      `json:"class"`           // 任务分类
		Type            int32      `json:"type"`            // 类型定义
		Typeson         int32      `json:"typeson"`         // 子分类
		Limit           int32      `json:"limit"`           // 达成条件
		Reward          []KeyValue `json:"reward"`          // 一次性奖励，使用json.RawMessage以支持多种类型
		NextAchievement int32      `json:"nextAchievement"` // 配置完成该成就后关联的下一个成就任务的 UUID（后续任务的任务ID）。
		Level           []int32    `json:"level"`           // 接取该任务的等级段数组（双闭区间），每对值表示 [最小等级, 最大等级]。
		Job             int32      `json:"job"`             // 职业限制
		VipLevel        int32      `json:"vipLevel"`        // vip任务
		ActivID         int32      `json:"activID"`         // 活动ID
		SevenDay        int32      `json:"sevenDay"`        // 七日任务
		LinkPanel       string     `json:"linkPanel"`       // 跳转界面
		TaskName        string     `json:"taskName"`        // 成就名称
		TaskDes         string     `json:"taskDes"`         // 条件描述
		TaskIcon        string     `json:"taskIcon"`        // 成就图标
		ResetMark       bool       `json:"ResetMark"`       // 指示继承任务在接取时是否重置计数：1/true 表示重置，0/false 或未设置表示不重置。
		LimitTime       StartEnd   `json:"LimitTime"`
		ActivityReward  KeyValue   `json:"ActivityReward"` // 活跃度奖励
		FirstRound      bool       `json:"FirstRound"`     // task.TaskMark 表示该任务是否为有后续任务的首个任务：为1/true 表示是首个并有后续，0/false 或未填表示需完成前置任务后才能接取。
	}

	PrivilegeConfig struct {
		Pid           int32      // 特权ID
		Sid           int32      // 商品ID
		Pname         string     // 特权名称
		Dailyaward    []Award    // 每日奖励
		Pcontent      []KeyValue // 特权
		Pdisplay      string     // 特权展示
		Duration      int32      // 持续时间
		Energystate   int32      // 活力满满增加的时间
		SkipAd        int32      // 跳过广告
		StaminaUp     int32      // 体力增加
		QuicklyIsiand int32      // 扫荡海岛
		QuicklyDaily  int32      // 扫荡日常
		AutoBattle    int32      // 自动战斗
	}
	RankingAwardConfig struct {
		UUID  int32      // uuid
		Type  int32      // 类型
		Rank  int32      // 名次
		Award []KeyValue // 每日奖励
	}
	// 随机属性表
	RandomCreat struct {
		RandomID     int32   // 随机属性ID
		SetID        []int32 // 装备部位
		Quality      int32
		CatalogueNum []KeyValue // 条数配置
		KindNum      []KeyValue // 种类配置
		DropSet      []KeyValue // 掉落分段配置
		RefiningSet  []KeyValue // 洗炼分段配置
	}
	RandomDistrib struct {
		UUID     int32  `json:"uuid"`     // 唯一标识
		Kind     int32  `json:"kind"`     // 属性种类
		Name     string `json:"name"`     // 属性名称
		GroupID  int32  `json:"group_id"` // 等级分组
		ValueID  int32  `json:"value_id"` // 分段ID
		Color    int32  `json:"color"`    // 品质颜色
		Interval map[int32]Interval
	}
	VIPInfo struct {
		VIPLevel    int32      `json:"vip_level"`   // VIP等级
		VIPName     string     `json:"vip_name"`    // VIP名称
		NeedExp     int32      `json:"need_exp"`    // 充值额度
		Partnerslot int32      `json:"partnerslot"` // 伙伴槽位
		Energyadd   int32      `json:"energyadd"`   // 体力增加
		Pcontent    []KeyValue `json:"pcontent"`    // VIP特权内容
		Pdisplay    string     `json:"pdisplay"`    // 特权展示
		Vdisplay    int32      `json:"vdisplay"`    // 显示价值
		DailyTime   []KeyValue `json:"dailytime"`   // 每日副本次数增加
	}

	RobotConfig struct {
		UUID        int32      `json:"uuid"`        // 使用字符串类型来表示UUID
		Job         int32      `json:"job"`         // 职业
		Level       int32      `json:"level"`       // 等级
		Power       int32      `json:"power"`       // 战力
		Feature     []KeyValue `json:"feature"`     // 属性配置
		Skill       []int      `json:"skill"`       // 技能配置
		Partners    []IdKeyNum `json:"partners"`    // 伙伴配置
		Name        string     `json:"name"`        // 名称
		RankMark    int32      `json:"rankMark"`    // 排行标识
		OfficialEXP int64      `json:"OfficialEXP"` // 官职经验
	}

	EnhanceConfig struct {
		UUID            int32       `json:"uuid"`            // uuid
		Job             int32       `json:"job"`             // 职业
		Strengthenlevel int32       `json:"strengthenlevel"` // 强化等级
		Consume         []IdNum     `json:"consume"`         // 道具消耗
		Attribute       []Attribute `json:"attribute"`       // 属性加成
	}
	StarsConfig struct {
		UUID        int32       `json:"uuid"`         // uuid
		Job         int32       `json:"job"`          // 职业
		StarLevel   int32       `json:"star_level"`   // 升星等级
		DecLevel    int32       `json:"dec_level"`    // 失败降级
		FailLevel   int32       `json:"fail_level"`   // 失败后等级
		SuccessRate int32       `json:"success_rate"` // 成功概率
		MaterialID  []IdNum     `json:"material_id"`  // 所需材料
		Attribute   []Attribute `json:"attribute"`    // 升星属性
	}
	Attribute struct {
		Slot      int32 `json:"slot"`      // 格位
		Attribute int32 `json:"attribute"` // 属性id
		Value     int32 `json:"value"`     // 值
	}

	Offlineraward struct {
		IdNum
		Gap int `json:"gap"`
		Max int `json:"max"`
	}
	// 无尽试炼
	EndlesstrialsConfig struct {
		LevelLimit    int32           `json:"Levellimit"`    // 等级限制
		Layers        int32           `json:"Layers"`        // 闯关层数
		Monster       []IdNumP        `json:"Monster"`       // 刷怪配置
		TimeLimit     int32           `json:"Timelimit"`     // 挑战时间
		Settlement    []IdNum         `json:"Settlement"`    // 结算奖励
		OfflineReward []Offlineraward `json:"Offlinereward"` // 离线收益效率
		Preview       string          `json:"Preview"`       // 产出预览
	}

	EquipMasterConfig struct {
		UUID      int32      `json:"UUID"`
		Category  int32      `json:"Category"`
		Limit     int32      `json:"Limit"`
		Attribute []KeyValue `json:"Attribute"`
	}

	PrizeWheelConfig struct {
		Uid        int32   `json:"Uid"`
		Round      int32   `json:"Round"`    // 奖励轮次
		Consume    int32   `json:"Consume"`  // 抽取消耗
		Rate       float64 `json:"rate"`     // 奖励倍率
		Priority   int32   `json:"Priority"` // 权值
		LevelLimit int32   `json:"LevelLimit"`
	}
	TreasureBoxConfig struct {
		Uid         int32                   `json:"Uid"`
		LevelRange  []int                   `json:"LevelRange"`
		Consume     []TreasureBoxConsume    `json:"Consume"`     // 抽取消耗
		CommonPool  []TreasureBoxCommonPool `json:"CommonPool"`  // 普通奖池配置
		SpecialPool []TreasureBoxCommonPool `json:"SpecialPool"` // 必中奖池配置
		Node        []int                   `json:"Node"`        // 必中次数配置
		MaxTime     int32                   `json:"MaxTime"`     // 最大抽取次数
	}
	TreasureBoxConsume struct {
		Type int32 `json:"type"`
		Id   int32 `json:"id"`
		Num  int32 `json:"num"`
	}
	TreasureBoxCommonPool struct {
		Position int32 `json:"position"`
		Id       int32 `json:"id"`
		Num      int32 `json:"num"`
		Priority int32 `json:"priority"` // 权重
	}
	RandNameConfig struct {
		ID         int32  // 编号，使用int32类型
		Firstname  string // 名，使用string类型
		Familyname string // 姓，使用string类型
		Lastname   string // 字，使用string类型
	}
	FateConfig struct {
		FateLevel     int32         `json:"FateLevel"`     // 命格等级
		Multiplicity  string        `json:"Multiplicity"`  // 重数显示
		Level         string        `json:"Level"`         // 级数显示
		FateExp       int32         `json:"FateExp"`       // 升级经验
		UpValue       []FateUpValue `json:"UpValue"`       // 升级属性
		AdvancedNeed  []IdNum       `json:"AdvancedNeed"`  // 进阶消耗
		AdvancedValue []KeyValue    `json:"AdvancedValue"` // 进阶属性
	}
	ComposeConfig struct {
		ComposeID   int32            `json:"ComposeID"`   // 合成编号
		Raw         []int32          `json:"Raw"`         // 源材料
		Produced    []ProducedConfig `json:"Produced"`    // 结果材料
		NeedGold    IdNum            `json:"NeedGold"`    // 消耗金币
		NeedNum     int32            `json:"NeedNum"`     // 消耗数量
		SuccessRate int32            `json:"SuccessRate"` // 成功几率
		Return      []KeyValueEx2    `json:"Return"`      // 失败返还
	}
	ProducedConfig struct {
		ID    []int32 `json:"id"`
		Ratio int32   `json:"ratio"`
	}
	FateUpValue struct {
		Job   int32 `json:"job"`
		Key   int32 `json:"key"`
		Value int32 `json:"value"`
	}

	// 掉落控制
	ItemDropControlConfig struct {
		Uid          int32
		ItemIDs      []int32  // 道具ID
		ControlTime  int32    // 控制间隔 秒
		ControlLevel Interval // 可掉落等级段，由于可能包含多种数据类型，使用 RawMessage
		Maxnum       int32    // 上限数量
		CandropMaps  []int32  // 可掉落地图，由于可能包含多种数据类型，使用 RawMessage
		Scope        int32    // 1 表示个人, 2 表示全服掉落
	}

	// 等级奖励
	LevelGiftConfig struct {
		Uid   int32
		Level int32
		Gift  []IdNumJobGender
	}

	// 伙伴配置
	PartnerConfig struct {
		Cid            int32      `json:"ID"`
		Name           string     `json:"Name"`
		Level          int32      `json:"Level"`
		Quality        int32      `json:"Quality"`
		Star           int32      `json:"Star"`
		Job            int32      `json:"job"`
		AttackRange    int32      `json:"AttackRange"`
		Skills         []KeyValue `json:"Skills"`
		AttackDistance int        `json:"AttackDistance"`
		AI             int32      `json:"AI"`
		Icon           string     `json:"Icon"`
		IconName       string     `json:"Iconname"`
		ItemID         int32      `json:"ItemID"`
		StarConsume    []KeyValueEx
		ChipNum        []KeyValue
		StarAttribute  []KeyValueEx
		BornValue      []KeyValue `json:"BornValue"`
		Extropower     int32
		TalentList     []int32
		TalentPoints   []KeyValue
	}

	//伙伴天赋
	PartnerTalentConfig struct {
		UID     int32      `json:"UID"`
		Name    string     `json:"Name"`
		Slot    int32      `json:"Slot"`
		Level   int32      `json:"Level"`
		Value   []KeyValue `json:"Value"`
		Consume KeyValue   `json:"Consume"`
		Limit   KeyValue   `json:"Limit"`
		Path    string     `json:"Path"`
		Icon    string     `json:"Icon"`
	}

	MaintriggerConfig struct {
		ID              int32           // 编号，使用int32类型
		Level           []LevelInterval `json:"Level"`
		MonsterIds      []MonsterIds    `json:"MonsterIds"`
		BackgroundmapID int32           `json:"BackgroundmapID"`
	}
	LevelInterval struct {
		MinLevel int32 `json:"minlevel"`
		MaxLevel int32 `json:"maxlevel"`
	}
	MonsterIds struct {
		Type int32 `json:"type"`
		Id   int32 `json:"id"`
		Num  int32 `json:"num"`
		Rate int32 `json:"rate"`
	}
	OpenBoxConfig struct {
		ID            int32           // 编号，使用int32类型
		Level         []LevelInterval `json:"Level"`
		Equiprandom   []IdRateJob     `json:"Equiprandom"`
		Qualityrandom []IdRate        `json:"Qualityrandom"`
	}

	IdRate struct {
		Id   int32 `json:"id"`
		Rate int32 `json:"rate"`
	}
	IdRateJob struct {
		Id   int32 `json:"id"`
		Rate int32 `json:"rate"`
		Job  int32 `json:"job"`
	}
	IdNumRate struct {
		Id     int32 `json:"id"`
		ItemId int32 `json:"itemid"`
		Num    int32 `json:"num"`
		Rate   int32 `json:"rate"`
	}
	BoxLevelConfig struct {
		ID            int32    // 编号，使用int32类型
		Level         int32    `json:"Level"`
		Rolelevel     int32    `json:"Rolelevel"`
		BattleID      int32    `json:"BattleID"`
		Needtime      int32    `json:"Needtime"`
		Consume       []IdNum  `json:"Consume"`
		Maxtime       int32    `json:"Maxtime"`
		Qualityrandom []IdRate `json:"Qualityrandom"`
	}
	SummonConfig struct {
		Uid         int32              `json:"Uid"`
		DrawModule  int32              `json:"DrawModule"`
		Drawtype    int32              `json:"Drawtype"`
		LevelRange  []int              `json:"LevelRange"`
		Consume     []SummonConsume    `json:"Consume"`     // 抽取消耗
		CommonPool  []SummonCommonPool `json:"CommonPool"`  // 普通奖池配置
		SpecialPool []SummonCommonPool `json:"SpecialPool"` // 必中奖池配置
		Node        []int              `json:"Node"`        // 必中次数配置
		MaxTime     int32              `json:"MaxTime"`     // 最大抽取次数
	}
	SummonConsume struct {
		Type int32 `json:"type"`
		Id   int32 `json:"id"`
		Num  int32 `json:"num"`
	}
	SummonCommonPool struct {
		Kind     int32 `json:"kind"`
		Id       int32 `json:"id"`
		Num      int32 `json:"num"`
		Priority int32 `json:"priority"` // 权重
	}

	PetConfig struct {
		ID            int32              `json:"ID"`            // id
		Name          string             `json:"Name"`          // 名字
		PetKind       KeyValue           `json:"PetKind"`       // 宠物类型
		Quality       int32              `json:"Quality"`       // 品质
		Mark          int32              `json:"Mark"`          // 角标
		StarLimit     []IdValues         `json:"StarLimit"`     // 升星消耗
		StarConsume   []IdValues         `json:"StarConsume"`   // 升星消耗
		StarAttribute []KeyValueEx       `json:"StarAttribute"` // 升星属性
		Skills        []map[string]int32 `json:"CommonPool"`    // 普通奖池配置
		Icon          string             `json:"Icon"`          // 图标
		Iconname      string             `json:"Iconname"`      // 图标名字
		ItemID        int32              `json:"ItemID"`        // 模板id
		OfflineAward  []KeyValuesOffline `json:"OfflineAward"`  // 离线收益效率
		Talent        []Talent           `json:"Talent"`        // 宠物特殊天赋效果
	}
	Talent struct {
		ID    int32   `json:"id"`
		Rate  int32   `json:"rate"` // 几率
		Value []IdNum `json:"value"`
	}
	KeyValuesOffline struct {
		Key    int32           `json:"key"`
		Values []Offlineraward `json:"values"`
	}
	PetStarConfig struct {
		Star          int32 `json:"Star"`          // 星级
		Exp           int32 `json:"exp"`           // 升阶所需经验
		StarAttribute int32 `json:"starAttribute"` // 升星属性
		StarSKill     int32 `json:"starSKill"`     // 升星属性
	}
	PetIdValues struct {
		Id     int32      `json:"id"`
		Values []KeyValue `json:"values"`
	}

	// 战役配置
	BattleConfig struct {
		UID          int32         `json:"UID"`          // 唯一ID
		CopyType     int32         `json:"Copytype"`     // 副本类型
		LevelLimit   int32         `json:"LevelLimit"`   // 等级限制
		PowerLimit   int32         `json:"Powerlimit"`   // 战力限制
		Chapter      int32         `json:"Chapter"`      // 章节
		Layers       int32         `json:"Layers"`       // 闯关层数
		Monster      []IdNumP      `json:"Monster"`      // 刷怪配置（JSON格式）
		TimeLimit    int32         `json:"TimeLimit"`    // 挑战时间限制（单位：秒）
		Settlement   []IdNum       `json:"Settlement"`   // 结算奖励（JSON格式）
		Randomaward  []IdNumRate   `json:"Randomaward"`  // 随机奖励（JSON格式） [{"id":14,"num":12,"rate":1000},{"id":2,"num":20000,"rate":1000},{"id":1,"num":20000,"rate":1000}]
		Preview      []int32       `json:"Preview"`      // 产出预览（JSON格式）
		ChapterAward []IdNum       `json:"ChapterAward"` // 章节奖励（JSON格式）
		StarAward    []KeyValueEx2 `json:"StarAward"`    // 星级奖励（JSON格式）
		MapMark      int32         `json:"MapMark"`      // 章节末标记
		BattleEffect []IdNumRate   `json:"Battleeffect"` // 战斗效果（JSON格式）[{"id":1,"rate":300},{"id":2,"rate":500},{"id":3,"rate":800},{"id":4,"rate":800},{"id":5,"rate":800}]
		RecommendPOW int32         `json:"recommendPOW"` // 推荐战力
	}
	// GuildConfig 表示公会等级的配置。
	GuildConfig struct {
		GuildLevel   int32      `json:"guild_level"`   // 公会等级
		NeedEXP      int32      `json:"need_exp"`      // 升级所需经验
		AttributeAdd []KeyValue `json:"attribute_add"` // 增加的属性
		PrivilegeDes string     `json:"privilege_des"` // 特权描述
		WarehouseAdd int32      `json:"warehouse_add"` // 仓库增加的格子数
		MemberNum    int32      `json:"member_num"`    // 最大成员数
	}
	// GuildStoreConfig 表示公会商店物品的配置。
	GuildStoreConfig struct {
		ID            int32   `json:"id"`             // 唯一编号
		Goods         []IdNum `json:"goods"`          // 出售商品
		Price         []IdNum `json:"price"`          // 出售价格
		GuildLimit    int32   `json:"guild_limit"`    // 行会每日上限
		PersonalLimit int32   `json:"personal_limit"` // 个人每日限购
		UnlockLevel   int32   `json:"unlock_level"`   // 行会解锁等级
	}
	// Skill represents the structure for skill configuration.
	GuildSkillConfig struct {
		UUID       int32      `json:"uuid"`        // Unique identifier
		SkillID    int32      `json:"skill_id"`    // Skill ID
		SkillName  string     `json:"skill_name"`  // Skill name
		SkillLevel int32      `json:"skill_level"` // Skill level
		EXP        int32      `json:"exp"`         // Experience required for upgrade
		ValueAdd   []KeyValue `json:"value_add"`   // Additional attributes in JSON format
	}

	// 跑马灯
	LampConfig struct {
		ID       int32  `json:"ID"`
		Broad    string `json:"Broad"`
		Type     string `json:"Type"`
		Panel    string `json:"Panel"`
		Receptor int32  `json:"Receptor"`
	}

	// 技能释放时间延迟
	ActionsetConfig struct {
		SkillID  int32 `json:"SkillID"`
		Lasttime int32 `json:"Lasttime"`
		Job      int32 `json:"Job"`
	}

	// 活动控制表
	ActivityControllerConfig struct {
		UUID          int64 `json:"UUID"`          // 唯一标识
		PayId         int64 `json:"PayId"`         // 付费ID
		StartTimeType int32 `json:"StartTimeType"` // 开始时间类型
		StartTime     int64 `json:"StartTime"`     // 活动开始时间
		EndTimeType   int32 `json:"EndTimeType"`   // 结束时间类型
		EndTime       int64 `json:"EndTime"`       // 活动结束时间
	}

	ActivityQuestConfig struct {
		UUID                 int32      `json:"UUID"`                 // 唯一标识
		ActivityControllerID int32      `json:"ActivityControllerID"` // 活动控制ID
		SubGroupID           int32      `json:"SubGroupID"`           // 子组ID
		Day                  int32      `json:"Day"`                  // 第几天的
		Type                 int32      `json:"Type"`                 // 类型，对应 task 表
		Limit                int32      `json:"Limit"`                // 限制条件
		Reward               []KeyValue `json:"Reward"`               // 奖励内容
		Activity             int32      `json:"Activity"`             // 活跃度
	}
	GuardIsland struct {
		Hangup []GuardIslandRecord `xml:"record>rate"` // 挂机倍率
		Quick  []Quick             `xml:"Quick>range"`
	}
	GuardIslandRecord struct {
		ID         int32 `xml:"id,attr"`         // 唯一标识
		Type       int32 `xml:"type,attr"`       // 类型：1 free, 2 ad, 3 item, 4 charge
		Multiplier int32 `xml:"multiplier,attr"` // 挂机倍率
	}
	Quick struct {
		Minlevel  int32   `xml:"minlevel,attr"`  // 等级限制
		Maxlevel  int32   `xml:"maxlevel,attr"`  // 等级限制
		Battle    float32 `xml:"battle,attr"`    // 挂机效率 杀一只怪的时间(秒)
		Expbattle float32 `xml:"Expbattle,attr"` // 挂机效率 杀一只怪的经验(秒)
	}

	BonusConfig struct {
		ID    int32 `json:"ID"`    // 唯一标识
		Level int32 `json:"Level"` // 等级
		Gift  []IdNumJobGender
	}

	// 羁绊
	PartnerBondConfig struct {
		ID         int32      // 唯一ID
		BondID     int32      // 羁绊ID
		BondName   string     // 羁绊名称
		BondLevel  int32      // 羁绊等级
		PartnerIDs []int32    // 羁绊伙伴ID列表（存为逗号分隔字符串，如 "1001,1002"）
		Star       int32      // 星级
		Attribute  []KeyValue // 属性
	}

	// 主界面战斗
	MainBattleConfig struct {
		ID          int32     // 唯一ID
		Name        string    // 名称
		Flag        int32     // 标识
		Cost        int32     //消耗
		FirstReward []IdNum   // 首次奖励
		Drop        []IdNum   // 普通奖励
		Monster     [][]int32 // 刷怪
	}
	// 元素赋能配置
	ElementConfig struct {
		UID          int32      `json:"UID"`          // 唯一编号
		Part         int32      `json:"Part"`         // 元素部位
		Partname     string     `json:"Partname"`     // 部位名称
		Level        int32      `json:"Level"`        // 赋能等级
		Rate         int32      `json:"Rate"`         // 赋能成功率
		Levelvalue   []KeyValue `json:"Levelvalue"`   // 赋能属性
		Levelconsume []IdNum    `json:"Levelconsume"` // 赋能消耗
		Breaklevel   int32      `json:"Breaklevel"`   // 突破阶段
		Breakconsume []IdNum    `json:"Breakconsume"` // 突破消耗
		Breakvalue   []KeyValue `json:"Breakvalue"`   // 突破属性
	}
	ElementTrialsConfig struct {
		UID            int32        `json:"UID"`             // 唯一ID
		LevelLimit     int32        `json:"Levellimit"`      // 等级限制
		Layers         int32        `json:"Layers"`          // 闯关层数
		Monster        []IdPosition `json:"Monster"`         // 刷怪配置
		TimeLimit      int32        `json:"Timelimit"`       // 挑战回合
		Settlement     []IdNum      `json:"Settlement"`      // 结算奖励
		FirstPassAward []IdNum      `json:"Firstpassraward"` // 首通奖励
		Preview        []int32      `json:"Preview"`         // 产出预览
	}
	IdPosition struct {
		ID       int32 `json:"id"`
		Position int32 `json:"position"`
	}

	// 称号
	TitleConfig struct {
		UID         int32        `json:"uid"`         // 唯一编号
		Name        string       `json:"name"`        // 称号名称
		Slot        int32        `json:"slot"`        // 部位
		Goods       []KeyValueEx `json:"goods"`       // 激活属性（JSON）
		Limit       []KeyValue   `json:"limit"`       // 佩戴属性（JSON）
		Activity    []KeyValueEx `json:"activity"`    // 激活消耗（JSON）
		StarRate    []KeyValue   `json:"starRate"`    // 升星几率
		StarConsume []KeyValueEx `json:"starConsume"` // 升星消耗（JSON）
		Period      int32        `json:"period"`      // 时效(天) 0：永久
		Path        string       `json:"path"`        // 图标路径
		Icon        string       `json:"icon"`        // 称号图标
	}

	//称号共鸣
	TitleActivationConfig struct {
		ID        int32      `json:"id"`        // 唯一ID
		Slot      int32      `json:"slot"`      // 部位
		Limit     string     `json:"limit"`     // 共鸣条件
		StarLimit int32      `json:"starLimit"` // 共鸣星级
		Attribute []KeyValue `json:"attribute"` // 共鸣属性（JSON）
	}

	ActivationCodeConfig struct {
		ChannelType int32     `json:"channel_type"` // 渠道类型
		CodeType    int32     `json:"code_type"`    // 激活码种类
		Start       time.Time `json:"start"`        // 开始时间
		End         time.Time `json:"end"`          // 结束时间
		Reward      []IdNum   `json:"reward"`       // 道具内容
	}
	LotteryConfig struct {
		ID          int32  `json:"ID"`          // 唯一ID
		BagId       int32  `json:"bagId"`       // 背包ID
		Name        string `json:"name"`        // 名字
		Studioname  int32  `json:"studioname"`  // 图标
		IfBag       int32  `json:"ifBag"`       // 是否包
		ItemId      int32  `json:"itemId"`      // 道具ID
		ItemName    string `json:"itemName"`    // 道具名
		ItemNum     int32  `json:"itemNum"`     // 道具数量
		Probability int32  `json:"probability"` // 道具权重
		Type        int32  `json:"type"`        // 计算类型
		TimesOut    int32  `json:"timesOut"`    // 必出次数
		CountNum    int32  `json:"countNum"`    // 获得上限
		Show        int32  `json:"show"`        // 是否显示
	}
	KillMonControlConfig struct {
		Uid        int32
		MonsterIDs []int32 // 怪物ID
		Maxnum     int32   // 上限数量
	}
	SwitchConfig struct {
		OpenLevel int32 `json:"open_level"` // 开启等级
		OpenDay   int32 `json:"open_day"`   // 开服天数
		MinBattel int32 `json:"minBattel"`  // 首页副本光卡id
		ServerDay int32 `json:"serverDay"`  // 是否用服务器时间，如果是0，就是玩家注册时间
	}
	AvatarConfig struct {
		UID      int32         `json:"uid"`      // 唯一编号
		Name     string        `json:"name"`     // 名称
		ItemType int32         `json:"itemType"` // 类别
		Value    []AttrWithJob `json:"value"`    // 佩戴属性 (JSON格式)
		Period   int32         `json:"period"`   // 时效
	}
	// RandomCreatV2Config 定义了随机属性V2配置结构体，用于装备随机属性生成 pk: EquipSlot, MaxLevel
	RandomCreatV2Config struct {
		EquipSlot                     int32      `json:"equip_slot"`                     // 装备槽位
		MaxLevel                      int32      `json:"max_level"`                      // 等级下限
		AttributeTypeSet              []KeyValue `json:"attribute_type_set"`             // 种类配置， key: 种类ID, value:权重
		RefiningSegmentRandomSet      []KeyValue `json:"refing_segment_random_set"`      // 洗炼分段配置，key: 分段ID, value:权重
		RefiningSegmentDirectionalSet []KeyValue `json:"refing_segment_directional_set"` // 定向分段配置，key: 分段ID, value:权重
	}
	// AttributeSegment 定义了属性分段结构体，用于描述装备属性的不同分段信息  pk: AttributeType, MaxLevel，RefingSegment
	AttributeSegmentV2Config struct {
		AttributeType int32  // 属性
		MaxLevel      int32  // 等级下限
		RefingSegment int32  // 分段ID
		ValueMin      int64  // 数值分段下限
		ValueMax      int64  // 数值分段上限
		QualityColor  uint32 // 品质颜色
	}

	FateCoinConfig struct {
		ID         int32             `json:"ID"`        // 唯一编号
		EffectType pb.FateCoinEffect `json:"Type"`      // 类别
		BuffID     int32             `json:"BuffID"`    // buffID
		BuffTime   int32             `json:"Bufftime"`  // Buff持续时间(毫秒)
		SkillBuff  []KeyValue        `json:"Skillbuff"` // 技能buff效果
		EffectNum  int32             `json:"effectNum"` // 作用数量
	}

	ArtifactStrengthConfig struct {
		UID        int32   `json:"UID"`        // 唯一编号
		WeaponID   int32   `json:"weaponId"`   // 器魂编号
		WeaponName string  `json:"weaponName"` // 器魂名称
		OpenLimit  int32   `json:"openLimit"`  // 器魂等级（解锁门槛）
		LevelLimit int32   `json:"levelLimit"` // 升级条件（等级下限）
		LevelExp   int64   `json:"levelExp"`   // 升级经验
		LevelValue []IdNum `json:"levelValue"` // 升级属性（JSON 串）
	}
	ArtifactAdvanceConfig struct {
		UID          int32   `json:"uid"`          // 唯一编号
		WeaponID     int32   `json:"weaponId"`     // 神兵编号
		Rank         int32   `json:"rank"`         // 神兵阶数
		WeaponName   string  `json:"weaponName"`   // 神兵名称
		SkillName    string  `json:"skillName"`    // 特技名称
		OpenLimit    int32   `json:"openLimit"`    // 解锁等级
		AdvanceLimit []IdNum `json:"advanceLimit"` // 进阶条件
		Consume      []IdNum `json:"consume"`      // 进阶消耗（Json 串）
		AdvanceValue []IdNum `json:"advanceValue"` // 进阶属性（Json 串）
		SkillID      int32   `json:"skillId"`      // 技能ID
	}
	PushingWallMonsterConfig struct {
		ID    int32      `json:"id"`
		Level int32      `json:"level"`
		Item  []KeyValue `json:"item"`
	}
	OfficialConfig struct {
		UID     int32  `json:"uid"`     // 唯一编号
		Name    string `json:"name"`    // 官职名称
		EXP     int64  `json:"exp"`     // 升级所需经验
		Quality int32  `json:"quality"` // 官职品质
	}
)

type RoleKey struct {
	Job   int32
	Level int32
}

var (
	BasePath           string
	ServerIDRangeStart int32
	ServerIDRangeEnd   int32
	Rng                *rand.Rand
	LoginOpen          bool

	GlobalCfg                      XMLConfig
	RoomCfg                        []*RoomConfig
	MonsterCfg                     map[int32]*MonsterConfig
	MapCfg                         []*MapConfig
	ItemCfg                        map[int32]*ItemConfig
	ItemStoneCfg                   map[int32]map[int32]int32
	RoleCfg                        map[RoleKey]*RoleConfig
	SkillCfg                       map[int32]*SkillConfig
	SignCfg                        map[string]map[int32]*SignConfig
	ItemUseCfg                     map[int32]*ItemUseConfig
	StoreCfg                       map[int32]*StoreConfig
	ActivityCfg                    map[int32]*ActivityConfig
	ActivityCfgByType              map[int32][]*ActivityConfig // 对活动表二次解析后的数据，避免每次循环整个表
	BuffCfg                        map[int32]*BuffConfig
	TaskCfg                        map[int32]*TaskConfig
	PrivilegeCfg                   map[int32]*PrivilegeConfig
	RankingAwardCfg                map[int32]*RankingAwardConfig
	RandomCreatCfg                 []*RandomCreat
	RandomDistribCfg               []*RandomDistrib
	VIPCfg                         []*VIPInfo
	RobotCfg                       []*RobotConfig
	EnhanceCfg                     map[int32]*EnhanceConfig
	StarsCfg                       map[int32]*StarsConfig
	EndlesstrialsCfg               map[int32]*EndlesstrialsConfig
	EquipMasterCfg                 []*EquipMasterConfig
	EquipMasterCfgMap              map[int32]*EquipMasterConfig
	SwitchCfg                      map[int32]*SwitchConfig
	PrizeWheelCfg                  map[int32]*PrizeWheelConfig
	TreasureBoxCfg                 map[int32]*TreasureBoxConfig
	RandNameCfg                    []*RandNameConfig
	FateCfg                        map[int32]*FateConfig
	ComposeCfg                     map[int32]*ComposeConfig
	ItemDropControlCfg             []*ItemDropControlConfig
	ItemDropControlCfgMap          map[string]*ItemDropControlConfig // itemID-mapID -> ItemDropControlConfig
	LevelGiftCfg                   map[int32]*LevelGiftConfig
	PartnerCfg                     map[int32]*PartnerConfig
	PartnerTalentCfg               map[int32]*PartnerTalentConfig
	MaintriggerCfg                 map[int32]*MaintriggerConfig
	OpenBoxCfg                     map[int32]*OpenBoxConfig
	BoxLevelCfg                    map[int32]*BoxLevelConfig
	SummonCfg                      map[int32]*SummonConfig
	PetCfg                         map[int32]*PetConfig
	PetStarCfg                     map[int32]*PetStarConfig
	BattleCfg                      map[int32]*BattleConfig
	GuildCfg                       map[int32]*GuildConfig
	GuildStoreCfg                  map[int32]*GuildStoreConfig
	GuildSkillCfg                  map[int32]*GuildSkillConfig
	LampCfg                        map[int32]*LampConfig
	ActionsetCfg                   map[int32]*ActionsetConfig
	ActivityControllerCfg          map[int32]*ActivityControllerConfig
	ActivityQuestCfg               map[int32]*ActivityQuestConfig
	ActivityQuestCfgByControllerID map[int32][]*ActivityQuestConfig // 对活动任务表二次解析后的数据，避免每次循环整个表
	GuardIslandCfg                 map[int64]*GuardIslandRecord
	BonusCfg                       map[int32]*BonusConfig
	PartnerBondCfg                 []*PartnerBondConfig
	MainBattleCfg                  map[int32]*MainBattleConfig
	ElementCfg                     map[int32]*ElementConfig
	ElementTrialsCfg               map[int32]*ElementTrialsConfig
	TitleCfg                       map[int32]*TitleConfig
	TitleActivationCfg             map[int32]*TitleActivationConfig
	ActivationCodeCfg              map[string]*ActivationCodeConfig // ID -> *ActivationCodeConfig
	KillMonControlCfg              []*KillMonControlConfig
	LotteryCfg                     map[int32]*LotteryConfig
	LotteryCfgByType               map[int32][]*LotteryConfig
	AvatarCfg                      map[int32]*AvatarConfig
	RandomCreatCfgV2               map[int32][]*RandomCreatV2Config     // EquipSlot -> []*RandomCreatV2Config
	AttributeSegmentCfgV2          map[string]*AttributeSegmentV2Config // AttributeType-MaxLevel -> RefingSegment -> *AttributeSegmentV2Config
	AttributeSegmentCfgV2MaxLevel  []int32                              // MaxLevel 列表
	FateCoinCfg                    map[int32]*FateCoinConfig
	ArtifactStrengthCfg            map[int32]*ArtifactStrengthConfig
	ArtifactAdvanceCfg             map[int32]*ArtifactAdvanceConfig
	PushingWallMonsterCfg          map[int32]*PushingWallMonsterConfig
	OfficialCfg                    map[int32]*OfficialConfig
	OfficalListCfg                 []*OfficialConfig
)

func ConvertRangeToSlice(rangeStr string) ([]int, error) {
	// 去除字符串两端的方括号
	rangeStr = strings.Trim(rangeStr, "[]")

	// 将字符串分割成多个数字
	parts := strings.Split(rangeStr, ",")

	// 创建切片用于存储转换后的整数
	slice := make([]int, 0, len(parts))

	// 遍历分割后的字符串，转换每个部分为整数并添加到切片中
	for _, part := range parts {
		// 去除可能存在的空格
		trimmedPart := strings.TrimSpace(part)

		// 将字符串转换为整数
		num, err := strconv.Atoi(trimmedPart)
		if err != nil {
			return nil, fmt.Errorf("invalid number format in range: %s", trimmedPart)
		}
		slice = append(slice, num)
	}

	return slice, nil
}

func RandIn(max int) int {
	if max <= 0 {
		return 0
	}
	num, err := randx.RandUint32(uint32(max))
	if err != nil {
		return max / 2
	}
	return int(num)
}

func ToPointI(l *pb.Point) *PointI {
	return &PointI{X: int(l.X), Y: int(l.Y)}
}

func ToPointF(l *PointI) *pb.Point {
	return &pb.Point{X: float64(l.X), Y: float64(l.Y)}
}

func GetRoomCfg(roomID int32) *RoomConfig {
	for _, v := range RoomCfg {
		if v.ID == roomID {
			return v
		}
	}
	return nil
}

func GetMonsterCfg(cfgId int32) *MonsterConfig {
	m, ok := MonsterCfg[cfgId]
	if ok {
		return m
	}
	return nil
}

func GetMapCfg(id int32) *MapConfig {
	for _, v := range MapCfg {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func GetMapCfgByMapID(mapid int32) []int32 {
	maps := []int32{}
	for _, v := range MapCfg {
		if v.MapID == mapid {
			maps = append(maps, v.ID)
		}
	}
	return maps
}

func GetItemCfg(id int32) (*ItemConfig, error) {
	item, ok := ItemCfg[id]
	if ok {
		return item, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetItemStoneCfg(id int32) (map[int32]int32, error) {
	item, ok := ItemStoneCfg[id]
	if ok {
		return item, nil
	}
	return nil, errutil.ErrConfigNotFound
}

//	for _, ItemInfo := range  global.ItemCfg{
//		for _, stoneSlotInfo := range ItemInfo.DiamondPosition{
//			if
//		}
//	}
func GetItemUseCfg(id int32) (*ItemUseConfig, error) {
	itemuse, ok := ItemUseCfg[id]
	if ok {
		return itemuse, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetRoleCfg(job, level int32) (*RoleConfig, error) {
	key := RoleKey{Job: job, Level: level}
	role, ok := RoleCfg[key]
	if ok {
		return role, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetSkillCfg(uid int32) (*SkillConfig, error) {
	skill, ok := SkillCfg[uid]
	if ok {
		return skill, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetUpperlevelSettingCfg(uid int32, advanced int32) (*SkillUpperlevel, error) {
	id := uid / 100
	for _, v := range SkillCfg {
		if v.SkillID == id && advanced == v.CanupperLevel {
			return &v.UpperlevelSetting, nil
		}
	}
	return nil, errutil.ErrConfigNotFound
}

func GetBuffCfg(id int32) (*BuffConfig, error) {
	buff, ok := BuffCfg[id]
	if ok {
		return buff, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetTaskCfg(id int32) (*TaskConfig, error) {
	task, ok := TaskCfg[id]
	if ok {
		return task, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetRandomDistribCfg(UUID int32, part int32) *RandomDistrib {
	for _, v := range RandomDistribCfg {
		if v.UUID == UUID && v.ValueID == part {
			return v
		}
	}
	return nil
}

func GetRandomDistribCfgByKind(kind int32, part int32) *RandomDistrib {
	for _, v := range RandomDistribCfg {
		if v.Kind == kind && v.ValueID == part {
			return v
		}
	}
	return nil
}

func GetEnhanceCfg(id int32) (*EnhanceConfig, error) {
	enhance, ok := EnhanceCfg[id]
	if ok {
		return enhance, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetStarsCfg(id int32) (*StarsConfig, error) {
	stars, ok := StarsCfg[id]
	if ok {
		return stars, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetActivityCfg(id int32) (*ActivityConfig, error) {
	activity, ok := ActivityCfg[id]
	if ok {
		return activity, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetEndlesstrialsCfg(id int32) *EndlesstrialsConfig {
	endlesstrials, ok := EndlesstrialsCfg[id]

	if ok {
		return endlesstrials
	}
	return nil
}

func GetPeizeWheelCfg(id int32) (*PrizeWheelConfig, error) {
	prizeWheel, ok := PrizeWheelCfg[id]

	if ok {
		return prizeWheel, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetTreasureBox(level int32) *TreasureBoxConfig {
	for _, v := range TreasureBoxCfg {
		if level >= int32(v.LevelRange[0]) && level <= int32(v.LevelRange[1]) {
			return v
		}
	}
	return nil
}

func GetSummon(dmType int32, ptype int32, level int32) *SummonConfig {
	for _, v := range SummonCfg {
		if dmType == v.DrawModule && ptype == v.Drawtype && level >= int32(v.LevelRange[0]) && level <= int32(v.LevelRange[1]) {
			return v
		}
	}
	return nil
}

func GetStoreCfg(id int32) (*StoreConfig, error) {
	StoreInfo, ok := StoreCfg[id]
	if ok {
		return StoreInfo, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GenRandName() string {
	firstNameIndex := RandIn(len(RandNameCfg))
	familyNameIndex := RandIn(len(RandNameCfg))
	lastNameIndex := RandIn(len(RandNameCfg))

	return RandNameCfg[firstNameIndex].Firstname + RandNameCfg[familyNameIndex].Familyname + RandNameCfg[lastNameIndex].Lastname
}

func GetFateCfg(level int32) (*FateConfig, error) {
	fate, ok := FateCfg[level]
	if ok {
		return fate, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetFateMaxCfg() (*FateConfig, error) {
	maxItem := FateCfg[1]
	for i := 2; i <= len(FateCfg); i++ { // 从第二个元素开始遍历
		if FateCfg[int32(i)].FateLevel > maxItem.FateLevel {
			maxItem = FateCfg[int32(i)] // 发现更大的元素，更新maxItem
		}
	}
	if maxItem.FateLevel > 0 {
		return maxItem, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetComposeCfg(composeId int32) (*ComposeConfig, error) {
	compose, ok := ComposeCfg[composeId]
	if ok {
		return compose, nil
	}
	return nil, errutil.ErrConfigNotFound
}

func GetEquipMaster(uid int32) (*EquipMasterConfig, error) {
	equipMaster, ok := EquipMasterCfgMap[uid]
	if ok {
		return equipMaster, nil
	}
	return nil, errutil.ErrConfigNotFound
}

//	func GetISCompose(composeId []int32) (*ComposeConfig, error) {
//		for
//		compose, ok := ComposeCfg[composeId]
//		if ok {
//			return compose, nil
//		}
//		return nil, errutil.ErrConfigNotFound
//	}
func GetVipCfg(vipLevel int32) *VIPInfo {
	for _, v := range VIPCfg {
		if v.VIPLevel == vipLevel {
			return v
		}
	}
	return nil
}

func GetBoxLevelCfg(boxLevel int32) *BoxLevelConfig {
	for _, v := range BoxLevelCfg {
		if v.Level == boxLevel {
			return v
		}
	}
	return nil
}

func GetPetCfg(id int32) *PetConfig {
	pet, ok := PetCfg[id]
	if ok {
		return pet
	}
	return nil
}

func GetDailyMapTime(roomType int32) *ChallengeTime {
	switch roomType {
	case constant.ExperienceFBMap:
		return &GlobalCfg.DailyMaps.ExpMap.ChallengeTime
	case constant.GoldFBMap:
		return &GlobalCfg.DailyMaps.GoldMap.ChallengeTime
	case constant.EnhanceStoneFBMap:
		return &GlobalCfg.DailyMaps.StrengthenMap.ChallengeTime
	case constant.BookPageFBMap:
		return &GlobalCfg.DailyMaps.SkillPageMap.ChallengeTime
	case constant.StarUpStoneFBMap:
		return &GlobalCfg.DailyMaps.StarStoneMap.ChallengeTime
	case constant.SoulStoneFBMap:
		return &GlobalCfg.DailyMaps.SoulStoneMap.ChallengeTime
	}
	return nil
}

func GetDailyMapScore(roomType int32) *Score {
	switch roomType {
	case constant.ExperienceFBMap:
		return &GlobalCfg.DailyMaps.ExpMap.Score
	case constant.GoldFBMap:
		return &GlobalCfg.DailyMaps.GoldMap.Score
	}
	return nil
}

func GetDailyMapDoubleaccess(roomType int32) *XMLIdNum {
	switch roomType {
	case constant.ExperienceFBMap:
		return &GlobalCfg.DailyMaps.ExpMap.DoubleAccess
	case constant.GoldFBMap:
		return &GlobalCfg.DailyMaps.GoldMap.DoubleAccess
	case constant.EnhanceStoneFBMap:
		return &GlobalCfg.DailyMaps.StrengthenMap.DoubleAccess
	case constant.BookPageFBMap:
		return &GlobalCfg.DailyMaps.SkillPageMap.DoubleAccess
	case constant.StarUpStoneFBMap:
		return &GlobalCfg.DailyMaps.StarStoneMap.DoubleAccess
	case constant.SoulStoneFBMap:
		return &GlobalCfg.DailyMaps.SoulStoneMap.DoubleAccess
	}
	return nil
}

func GetDailyMapAccommodation(roomType int32) int {
	switch roomType {
	case constant.ExperienceFBMap:
		return GlobalCfg.DailyMaps.ExpMap.Accommodation.Num
	case constant.GoldFBMap:
		return GlobalCfg.DailyMaps.GoldMap.Accommodation.Num
	case constant.EnhanceStoneFBMap:
		return GlobalCfg.DailyMaps.StrengthenMap.Accommodation.Num
	case constant.BookPageFBMap:
		return GlobalCfg.DailyMaps.SkillPageMap.Accommodation.Num
	}
	return 0
}

// IsToken 判断是否是货币
func IsToken(id int32) bool {
	switch id {
	case constant.ItemTypeCoin,
		constant.ItemTypeSoul,
		constant.ItemTypeCoupon,
		constant.ItemTypeTrialVoucher,
		constant.ItemTypeTreasureNum,
		constant.ItemTypePetExp,
		constant.ItemTypeAccelerateScroll,
		constant.ItemTypeGuildContribution,
		constant.ItemTypeVipExp,
		constant.ItemTypeGuildWish,
		constant.ItemTypeStarWishCoin,
		constant.ItemTypePartnerTalentPoint:
		return true
	default:
		return false
	}
}

func GetPartnerCfg(id int32) *PartnerConfig {
	partner, ok := PartnerCfg[id]
	if ok {
		return partner
	}
	return nil
}

func GetPartnerCfgByItemId(id int32) *PartnerConfig {
	for _, partner := range PartnerCfg {
		if partner.ItemID == id {
			return partner
		}
	}
	return nil
}

func GetEquipqualityCoefficient(q pb.ItemQuality) float64 {
	for _, v := range GlobalCfg.EquipQualitys {
		if v.ID == int32(q) {
			return v.Value
		}
	}
	return 1
}

func GetRoleInitial(job int32) *Initial {
	for _, v := range GlobalCfg.Initials {
		if v.Job == job {
			return &v
		}
	}
	return &Initial{Job: job, AttackSpeed: 500, MoveSpeed: 250}
}

func GetpartnerQualityRate(quality int32) float64 {
	for _, v := range GlobalCfg.Partnerquality {
		if v.ID == quality {
			return v.Value
		}
	}
	return 1
}

func GetpetQualityRate(quality int32) float64 {
	for _, v := range GlobalCfg.Petquality {
		if v.ID == quality {
			return v.Value
		}
	}
	return 1
}

func GetpetLevelRate(quality int32) float64 {
	for _, v := range GlobalCfg.PetLevel {
		if v.ID == quality {
			return v.Value
		}
	}
	return 1
}

func GetBattleMapCfg(id int32) *BattleConfig {
	battle, ok := BattleCfg[id]
	if ok {
		return battle
	}
	return nil
}

// 当前应该挑战的关卡ID
func GetNextBattleId(CopyType int32, curBattleId int32) int32 {
	var battleId int32

	var copyType_ int32
	var chapter_ int32
	var layers_ int32

	copyType_ = CopyType
	switch CopyType {
	case constant.BattleMap:
		if curBattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			copyType_ = curBattleId / 100000
			chapter_ = curBattleId % 100000 / 1000
			layers_ = curBattleId % 1000

			if layers_ == 20 {
				chapter_++
				layers_ = 1

			} else {
				layers_++
			}
			battleId = copyType_*100000 + chapter_*1000 + layers_
		}
	case constant.BossMap:
		if curBattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			battleId = curBattleId + 1
		}
	case constant.TrialTowerMap:
		if curBattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			battleId = curBattleId + 1
		}
	case constant.RelicMap:
		if curBattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			copyType_ = curBattleId / 100000
			chapter_ = curBattleId % 100000 / 1000
			layers_ = curBattleId % 1000

			if layers_ == 10 {
				chapter_++
				layers_ = 1
			} else {
				layers_++
			}
			battleId = copyType_*100000 + chapter_*1000 + layers_
		}
	case constant.GuildMap:
		battleId = copyType_*100000 + curBattleId + 1

	default:
		return 0
	}
	return battleId
}

func GetCopySetCfg(id int32) int32 {
	for _, v := range GlobalCfg.Copyset {
		if v.ID == id {
			return int32(v.Num)
		}
	}
	return 0
}

// 品质 等级上限 部位
func GetRandomCreatCfg(quality, level, setid int32) *RandomCreat {
	for _, cfg := range RandomCreatCfg {
		if cfg.Quality != quality {
			continue
		}
		for _, set := range cfg.SetID {
			if set == setid && cfg.RandomID >= level {
				return cfg
			}
		}
	}
	return nil
}

// GetGuildCfg 获取指定等级的公会配置。
func GetGuildCfg(level int32) *GuildConfig {
	cfg, ok := GuildCfg[level]
	if ok {
		return cfg
	}
	return nil
}

// GetGuildStoreCfg 获取指定ID的公会商店配置。
func GetGuildStoreCfg(id int32) *GuildStoreConfig {
	cfg, ok := GuildStoreCfg[id]
	if ok {
		return cfg
	}
	return nil
}

// GetGuildSkillCfg 获取指定ID的公会技能配置。
func GetGuildSkillCfg(id int32) *GuildSkillConfig {
	cfg, ok := GuildSkillCfg[id]
	if ok {
		return cfg
	}
	return nil
}

// 查询特权卡
func GetPrivilegeCfg(id int32) *PrivilegeConfig {
	cfg, ok := PrivilegeCfg[id]
	if ok {
		return cfg
	}
	return nil
}

func GetRankAwardCfg(id int32) *RankingAwardConfig {
	cfg, ok := RankingAwardCfg[id]
	if ok {
		return cfg
	}
	return nil
}

// 经验转换位经验道具
func Exp2ExpItems(exp int64) []IdNum {
	var items []IdNum
	var left int64 = exp
	for i := len(GlobalCfg.Partner.ExpItem) - 1; i >= 0; i-- {
		num := left / GlobalCfg.Partner.ExpItem[i].Num
		left = left % GlobalCfg.Partner.ExpItem[i].Num
		if num > 0 {
			items = append(items, IdNum{ID: GlobalCfg.Partner.ExpItem[i].ID, Num: int32(num)})
		}
		if left <= 0 {
			break
		}
	}
	return items
}

// 获取技能释放时间延迟
func GetActionsetCfg(skillId int32) int32 {
	if cfg, ok := ActionsetCfg[skillId]; ok {
		return cfg.Lasttime
	}

	return 500
}

// 行会副本开启关闭时间
func GetGuildCopyTime() (time.Time, time.Time, error) {
	var err error
	var starttime, endtime time.Time
	if starttime, err = time.Parse(time.TimeOnly, GlobalCfg.Guild.GuildCopy.Starttime); err != nil {
		return time.Now(), time.Now(), errutil.ErrIllegalParameter
	}

	if endtime, err = time.Parse(time.TimeOnly, GlobalCfg.Guild.GuildCopy.Endtime); err != nil {
		return time.Now(), time.Now(), errutil.ErrIllegalParameter
	}
	t := time.Now()

	startTime_ := time.Date(t.Year(), t.Month(), t.Day(), starttime.Hour(), starttime.Minute(), 0, 0, t.Location())
	endTime_ := time.Date(t.Year(), t.Month(), t.Day(), endtime.Hour(), endtime.Minute(), 0, 0, t.Location())
	return startTime_, endTime_, err
}

// 一战到底开启关闭时间
func GetFightEndTime() (time.Time, time.Time, error) {
	var err error
	var starttime, endtime time.Time
	if starttime, err = time.Parse(time.TimeOnly, GlobalCfg.FightEnd.StartTime); err != nil {
		return time.Now(), time.Now(), errutil.ErrIllegalParameter
	}

	if endtime, err = time.Parse(time.TimeOnly, GlobalCfg.FightEnd.EndTime); err != nil {
		return time.Now(), time.Now(), errutil.ErrIllegalParameter
	}
	t := time.Now()

	startTime_ := time.Date(t.Year(), t.Month(), t.Day(), starttime.Hour(), starttime.Minute(), 0, 0, t.Location())
	endTime_ := time.Date(t.Year(), t.Month(), t.Day(), endtime.Hour(), endtime.Minute(), 0, 0, t.Location())
	return startTime_, endTime_, err
}

// 获取普通攻击技能ID
func GetNormalSkillId(job constant.MowJob) int32 {
	skillId := int32(0)
	switch job {
	case constant.Warrior:
		skillId = constant.SkillTypeNormalAttackWarrior
	case constant.Wizard:
		skillId = constant.SkillTypeNormalAttackWizard
	case constant.Ritualist:
		skillId = constant.SkillTypeNormalAttackRitualist
	case constant.PartnerWarrior:
		skillId = constant.PartnerWarriorNormalSkillID
	case constant.PartnerWizard:
		skillId = constant.PartnerWizardNormalSkillID
	case constant.PartnerRitualist:
		skillId = constant.PartnerRitualistNormalSkillID
	}
	return skillId * 100
}

func GetBonusCfg(day int32) []IdNumJobGender {
	if cfg, ok := BonusCfg[day]; ok {
		return cfg.Gift
	}
	return nil
}

// 伙伴羁绊
func GetPartnerBondCfg(id int32) *PartnerBondConfig {
	for _, v := range PartnerBondCfg {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func GetMainBattleCfg(id int32) *MainBattleConfig {
	if cfg, ok := MainBattleCfg[id]; ok {
		return cfg
	}
	return nil
}
func GetElementCfg(id int32) *ElementConfig {
	if cfg, ok := ElementCfg[id]; ok {
		return cfg
	}
	return nil
}
func GetElementTrialsCfg(id int32) *ElementTrialsConfig {
	if cfg, ok := ElementTrialsCfg[id]; ok {
		return cfg
	}
	return nil
}
func GetElementTrialsSwttlementCfg(Battle int32) map[int32]int32 {
	settleMent := make(map[int32]int32)
	for _, v := range ElementTrialsCfg {
		if v.UID <= Battle {
			for _, v := range v.Settlement {
				settleMent[v.ID] += v.Num
			}
		}
	}
	return settleMent
}

// 跨服战场开启关闭时间
// func GetWarTime() (time.Time, time.Time, error) {
// 	for _, v := range GlobalCfg.War.OpenTime {
// 		var err error
// 		var starttime, endtime time.Time
// 		if starttime, err = time.Parse(time.TimeOnly, v.StartTime); err != nil {
// 			continue
// 		}

// 		if endtime, err = time.Parse(time.TimeOnly, v.EndTime); err != nil {
// 			continue
// 		}
// 		t := time.Now()

// 		startTime_ := time.Date(t.Year(), t.Month(), t.Day(), starttime.Hour(), starttime.Minute(), 0, 0, t.Location())
// 		endTime_ := time.Date(t.Year(), t.Month(), t.Day(), endtime.Hour(), endtime.Minute(), 0, 0, t.Location())

// 		if time.Now().After(startTime_) && time.Now().Before(endTime_) {
// 			return startTime_, endTime_, err
// 		}
// 	}

// 	return time.Now(), time.Now(), errutil.ErrConfigNotFound
// }

// func InWarTime() bool {
// 	if GlobalCfg.War.OpenTime == nil {
// 		return false
// 	}

// 	now := time.Now()
// 	currentHour := now.Hour()
// 	currentMinute := now.Minute()
// 	currentSecond := now.Second()

// 	// 判断是否在开放日
// 	var weekOpen bool
// 	for _, v := range GlobalCfg.War.OpenDay {
// 		if v == int(now.Weekday()) {
// 			weekOpen = true
// 			break
// 		}
// 	}
// 	if !weekOpen {
// 		return false
// 	}

// 	// 判断是否在开放时间
// 	currentTotalMinutes := currentHour*60 + currentMinute
// 	currentTotalSeconds := currentTotalMinutes*60 + currentSecond

// 	for _, tr := range GlobalCfg.War.OpenTime {
// 		startTotalMinutes := tr.StartHour*60 + tr.StartMinute
// 		endTotalMinutes := tr.EndHour*60 + tr.EndMinute

// 		startTotalSeconds := startTotalMinutes * 60
// 		endTotalSeconds := endTotalMinutes * 60

// 		if currentTotalSeconds > startTotalSeconds && currentTotalSeconds < endTotalSeconds {
// 			return true
// 		}
// 	}
// 	return false
// }

// 跨服战场开启关闭时间
func GetWarTime(roomID int32) (time.Time, time.Time, error) {
	for _, v := range GlobalCfg.War.OpenTime {
		if v.RoomId != roomID {
			continue
		}
		var err error
		var starttime, endtime time.Time
		if starttime, err = time.Parse(time.TimeOnly, v.StartTime); err != nil {
			continue
		}

		if endtime, err = time.Parse(time.TimeOnly, v.EndTime); err != nil {
			continue
		}
		t := time.Now()

		startTime_ := time.Date(t.Year(), t.Month(), t.Day(), starttime.Hour(), starttime.Minute(), 0, 0, t.Location())
		endTime_ := time.Date(t.Year(), t.Month(), t.Day(), endtime.Hour(), endtime.Minute(), 0, 0, t.Location())

		if time.Now().After(startTime_) && time.Now().Before(endTime_) {
			return startTime_, endTime_, err
		}
	}

	return time.Now(), time.Now(), errutil.ErrConfigNotFound
}

func InWarTime(roomID int32) bool {
	if GlobalCfg.War.OpenTime == nil {
		return false
	}

	now := time.Now()
	currentHour := now.Hour()
	currentMinute := now.Minute()
	currentSecond := now.Second()

	// 判断是否在开放时间
	currentTotalMinutes := currentHour*60 + currentMinute
	currentTotalSeconds := currentTotalMinutes*60 + currentSecond

	for _, tr := range GlobalCfg.War.OpenTime {
		if tr.RoomId != roomID {
			continue
		}
		if tr.Mark == 0 {
			return false
		}

		if tr.WeekDay != int32(now.Weekday()) {
			return false
		}

		startTotalMinutes := tr.StartHour*60 + tr.StartMinute
		endTotalMinutes := tr.EndHour*60 + tr.EndMinute

		startTotalSeconds := startTotalMinutes * 60
		endTotalSeconds := endTotalMinutes * 60

		if currentTotalSeconds > startTotalSeconds && currentTotalSeconds < endTotalSeconds {
			return true
		}
	}
	return false
}

// 获取称号配置
func GetTitleCfg(id int32) *TitleConfig {
	if cfg, ok := TitleCfg[id]; ok {
		return cfg
	}
	return nil
}

// TitleActivationCfg
func GetTitleActivationCfg(id int32) *TitleActivationConfig {
	if cfg, ok := TitleActivationCfg[id]; ok {
		return cfg
	}
	return nil
}
func GetLotteryCfg(id int32) *LotteryConfig {
	if cfg, ok := LotteryCfg[id]; ok {
		return cfg
	}
	return nil
}
func GetLotteryCfgByType(id int32) []*LotteryConfig {
	if cfg, ok := LotteryCfgByType[id]; ok {
		return cfg
	}
	return nil
}
func GetLotteryCfgByItem(ItemId int32) *LotteryConfig {
	for _, v := range LotteryCfg {
		if v.ItemId == ItemId {
			return v
		}
	}
	return nil
}

// GetRandomCreatCfgV2 获取随机属性V2配置
func GetRandomCreatCfgV2(equipSlot, level int32) (*RandomCreatV2Config, error) {
	if cfgs, ok := RandomCreatCfgV2[equipSlot]; ok {
		low, high := 0, len(cfgs)-1
		ans := -1
		for low <= high {
			mid := low + (high-low)/2
			if cfgs[mid].MaxLevel >= level {
				ans = mid
				high = mid - 1
			} else {
				low = mid + 1
			}
		}

		if ans != -1 {
			return cfgs[ans], nil
		}
	}
	return nil, errutil.ErrConfigNotFound
}

// GetAttributeSegmentCfgV2 获取属性分段配置
func GetAttributeSegmentCfgV2(attributeType, level, refingSeqment int32) (*AttributeSegmentV2Config, error) {
	var maxLevel int32
	log, high := 0, len(AttributeSegmentCfgV2MaxLevel)-1
	ans := -1
	for log <= high {
		mid := log + (high-log)/2
		if AttributeSegmentCfgV2MaxLevel[mid] >= level {
			ans = mid
			high = mid - 1
		} else {
			log = mid + 1
		}
	}
	if ans == -1 {
		return nil, errutil.ErrConfigNotFound
	}

	maxLevel = AttributeSegmentCfgV2MaxLevel[ans]

	key := fmt.Sprintf("%d_%d_%d", attributeType, maxLevel, refingSeqment)
	if cfg, ok := AttributeSegmentCfgV2[key]; ok {
		return cfg, nil
	}

	return nil, errutil.ErrConfigNotFound
}
