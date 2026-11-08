package constant

import "nanoserver/protocol/protobuf/pb"

type DeskStatus int32

const (
	// 创建桌子
	DeskStatusCreate DeskStatus = iota
	// 刷怪
	DeskStatusReady
	// 游戏
	DeskStatusPlaying
	// 空闲
	DeskStatusIdle
	// 已销毁
	DeskStatusDestory
)

var stringify = [...]string{
	DeskStatusCreate:  "创建",
	DeskStatusReady:   "准备好",
	DeskStatusPlaying: "游戏中",
	DeskStatusIdle:    "空闲",
	DeskStatusDestory: "已销毁",
}

func (s DeskStatus) String() string {
	return stringify[s]
}

const (
	// token相关
	JwtTokenOK            int    = 200100                      // token有效
	JwtTokenInvalid       int    = -400100                     // 无效的token
	JwtTokenExpired       int    = -400101                     // 过期的token
	JwtTokenFormatErrCode int    = -400102                     // 提交的 token 格式错误
	JwtTokenFormatErrMsg  string = "提交的 token 格式错误"            // 提交的 token 格式错误
	JwtTokenMustValid     string = "token为必填项,请在请求header部分提交!" // 提交的 token 格式错误
)

// 道具时效性类型
const (
	ItemTimelinessShiChang = 1 // 固定时长
	ItemTimelinessQuJian   = 2 // 时间范围
)

// 职业
type MowJob int32

const (
	Warrior          MowJob = iota + 1 // 战士
	Wizard                             // 法师
	Ritualist                          // 祭祀
	PartnerWarrior   = 11              // 伙伴战士
	PartnerWizard    = 12              // 伙伴法师
	PartnerRitualist = 13              // 伙伴祭祀
	Pet              = 21              // 宠物
)

// 对应职业名称
var JobName = map[MowJob]string{
	Warrior:          "战士",
	Wizard:           "法师",
	Ritualist:        "祭祀",
	PartnerWarrior:   "伙伴战士",
	PartnerWizard:    "伙伴法师",
	PartnerRitualist: "伙伴祭祀",
	Pet:              "宠物",
}

// 对应task名称
var TaskName = map[int32]string{
	1:  "日常任务",
	2:  "周常任务",
	3:  "成就任务",
	4:  "每日福利",
	5:  "日常活跃度奖励",
	6:  "周常活跃度奖励",
	7:  "主线任务",
	8:  "VIP任务",
	9:  "开服七天乐",
	10: "开服七天乐所有礼包",
}

// 地图类型枚举
const (
	// PVEPeacefulMap 代表PVE玩法中的和平模式地图
	PVEPeacefulMap int32 = iota + 1

	// PVEWarMap 代表PVE玩法中的战争模式地图
	PVEWarMap

	// PVPOneOnOneMap 代表PVP玩法中的1对1模式地图
	PVPOneOnOneMap

	// PVPMultiOnMultiMap 代表PVP玩法中的多对多模式地图
	PVPMultiOnMultiMap

	// SoloTowerMap 代表个人爬塔类玩法地图
	SoloTowerMap

	// 世界boss
	WorldBossMap

	// ExperienceFBMap 代表经验副本 boss
	ExperienceFBMap

	// GoldFBMap 代表金币副本 boss
	GoldFBMap

	// EnhanceStoneFBMap 代表强化石副本 波次
	EnhanceStoneFBMap

	// BookPageFBMap 代表书页副本 波次
	BookPageFBMap

	// StarUpStoneFBMap 代表升星石副本 采集
	StarUpStoneFBMap

	// SoulStoneFBMap 代表魂石副本 采集
	SoulStoneFBMap

	// GuidanceFBMap 代表新手引导
	GuidanceFBMap

	// 战役（主线）副本
	BattleMap
	// 首领副本
	BossMap
	// 试炼塔
	TrialTowerMap
	// 遗迹
	RelicMap
	// 工会
	GuildMap
	// 一战到底
	FightEndMap
	//回合制
	RoundMap = 10000
)

// BuffID 枚举定义了游戏中的各种buff类型。
type BuffID int

// 定义枚举常量，使用iota自动递增。
const (
	BUFF_UNSPECIFIED                BuffID = iota // 未指定的buff，可以作为默认值
	BLEED_5SEC                                    // 每5秒损失1%的血量，持续30秒
	ARMOR_BREAK                                   // 物防和魔防降低10%，持续10秒
	SLOW                                          // 移动速度降低50%，持续5秒
	BLEED_1SEC                                    // 每1秒损失1%的血量，持续10秒
	ARMOR_BREAK_ONLY                              // 物防降低10%，持续10秒
	POISON                                        // 攻击力降低30%，持续5秒
	ATTACK_PERCENT_BOOST                          // 攻击力得到百分比提升
	HEAL                                          // 持续回复一定比例的生命值
	DEFENSE_DOUBLE_BOOST                          // 物理防御和魔法防御得到百分比提升
	ARMOR_BOOST                                   // 物理防御得到提升
	MAGIC_ARMOR_BOOST                             // 魔法防御得到提升
	ATTACK_SPEED_MOVE_SPEED_BOOST                 // 攻击速度和移动速度得到提升
	ATTACK_DEFENSE_BOOST                          // 攻击力和双防都将会得到提升
	STAR_UPGRADE_CHANCE_BOOST                     // 升星的成功几率得到提升
	FREEZE                                        // 冰冻效果
	INJURY                                        // 伤血效果
	ATTACK_PERCENT_BOOST_2          BuffID = 20   // 此状态下，攻击力得到百分比提升（第二个同名buff）
	DEFENSE_DOUBLE_BOOST_2          BuffID = 21   // 物理防御和魔法防御得到百分比提升（第二个同名buff）
	HEAL_PERCENT                    BuffID = 22   // 一次性回复一定百分比的生命
	ATTACK_SPEED_MOVE_SPEED_BOOST_2 BuffID = 23   // 攻速移速都会提升（第二个同名buff）
	GOD_MODE                        BuffID = 24   // 无敌状态，免疫所有伤害和控制技能
)

// 道具增加/消耗日志
const (
	ItemAddWeb                = "Web服务端增加"
	ItemAddGame               = "游戏服务端增加"
	ItemAddPvp                = "PVP增加"
	ItemAddTask               = "任务"
	ItemAddDrop               = "掉落"
	ItemAddGM                 = "GM增加"
	ItemAddDistroy            = "销毁"
	ItemAddDisintegrate       = "分解"
	ItemAddSoloTower          = "爬塔"
	ItemAddCompose            = "合成材料"
	ItemAddGemstone           = "宝石卸下"
	ItemAddDailymap           = "每日必做"
	ItemAddLevelGift          = "等级奖励"
	ItemAddAFKMonster         = "挂机奖励怪物掉落"
	ItemAddAFKBox             = "挂机奖励道具"
	ItemAddOpenBox            = "挂机奖励箱子"
	ItemAddPartner            = "伙伴"
	ItemAddPetAuth            = "宠物自动给宠物蛋"
	ItemAddPetGold            = "宠物技能提供金币"
	ItemAddPetExp             = "宠物技能提供经验"
	ItemAddPetAFK             = "宠物离线"
	ItemAddGuildWish          = "自动给祈愿点"
	ItemAddGuildMap           = "公会副本"
	ItemAddGuildWishExp       = "祈愿提供经验"
	ItemAddGuildShopBuy       = "行会商店购买"
	ItemAddGuildDonate        = "行会捐献"
	ItemAddDropHaiDao         = "守卫海岛掉落"
	ItemAddDropStarUpStoneFB  = "升星石掉落"
	ItemAddDropSoulStoneFBMap = "魂石掉落"
	ItemAddPetFull            = "宠物满了"
	ItemAddContinuity         = "连续登录"
	ItemAddCumulative         = "累计登录"
	ItemAddFirstCharge        = "首冲"
	ItemAddActivityLevel      = "活动等级礼包"
	ItemAddTurntable          = "大转盘抽奖"
	ItemAddTreasureBox        = "百宝箱抽奖"
	ItemAddTrial              = "试炼基金"
	ItemAddGroupUp            = "成长基金"
	ItemAddMail               = "邮件"
	ItemAddMainBattle         = "主界面战斗"
	ItemAddElement            = "元素"
	ItemAddElementSweep       = "元素扫荡"
	ItemAddWar                = "群雄战场"
	ItemAddTitle              = "称号"
	ItemAddOfflineBox         = "离线恢复宝箱"
	ItemAddParTalent          = "伙伴天赋"

	ItemDelWeb            = "Web服务端消耗"
	ItemDelGame           = "游戏服务端消耗"
	ItemDelEnhance        = "强化消耗"
	ItemDelStar           = "升星消耗"
	ItemDelChangeRoleName = "改名消耗"
	ItemDelRefine         = "洗炼"
	ItemDelChat           = "聊天"
	ItemDelDistroy        = "销毁"
	ItemDelExpired        = "过期"
	ItemDelDisintegrate   = "分解"
	ItemDelSkillUpgrade   = "升级技能"
	ItemDelCompose        = "合成材料"
	ItemDelGemstone       = "宝石穿戴"
	ItemDelHomeOpenBox    = "宝石穿戴"
	ItemDelPetLevel       = "宠物升级"
	ItemDelPetStar        = "宠物升星"
	ItemDelGuildCreate    = "行会创建"
	ItemDelGuildDonate    = "行会捐献"
	ItemDelGuildShopBuy   = "行会商店购买"
	ItemDelVipExp         = "vip经验"
	ItemDelGuildWishExp   = "祈愿提供经验"
	ItemDelGuardIsland    = "挂机守卫海岛"
	ItemDelElementUpgrade = "元素升级"
	ItemDelTitle          = "称号"
	ItemDelArtifact       = "神兵"
	ItemDelParTalent      = "伙伴天赋"

	ItemFateUp           = "命格升级"
	ItemFateBreak        = "命格升阶"
	ItemOrder            = "订单"
	ItemPrivilege        = "特权卡"
	ItemPrizeDrawPartner = "抽奖伙伴"
	ItemPrizeDrawPet     = "抽奖宠物"
	ItemSign             = "签到"
)

// 道具ID ,拓展必须对应item.xlsx ID 字段
const (
	ItemTypeCoin                   int32 = 1     // 金币	游戏内的通用货币
	ItemTypeExp                    int32 = 2     // 经验
	ItemTypeGuildContribution      int32 = 3     // 工会贡献
	ItemTypeSoul                   int32 = 4     // 战魂	游戏内获得的数值货币，通过PVP模块可获取
	ItemTypeCoupon                 int32 = 5     // 点券	游戏内的充值货币，对应点券
	ItemTypeEnergy                 int32 = 6     // 体力	游戏内获得的数值货币，对应能量
	ItemTypeVipExp                 int32 = 7     // vip经验	游戏内获得的数值货币
	ItemTypeTrialVoucher           int32 = 8     // 试炼凭证	游戏内获得的数值货币
	ItemTypePVPTimes               int32 = 10    // PVP次数
	ItemTypeTreasureNum            int32 = 14    // 开箱次数
	ItemTypeEnhanceStoneChallenge  int32 = 15    // 强化石副本挑战卷轴
	ItemTypeStarUpStoneChallenge   int32 = 16    // 升星石副本挑战卷轴
	ItemTypeSoulStoneChallenge     int32 = 17    // 魂石副本挑战卷轴
	ItemTypeBookPageChallenge      int32 = 18    // 书页副本挑战卷轴
	ItemTypeGoldChallenge          int32 = 19    // 金币副本挑战卷轴
	ItemTypeExpChallenge           int32 = 20    // 经验副本挑战卷轴
	ItemTypeLoudspeaker            int32 = 29    // 喇叭
	ItemTypeAfkCard                int32 = 56    // 15分钟挂机卡
	ItemTypeStarWishCoin           int32 = 58    // 星愿币
	ItemTypeGVGTimes               int32 = 71    // GVG次数
	ItemTypeElementChallengeTimes  int32 = 72    // 属性强化挑战次数
	ItemTypeBossChallenge          int32 = 73    // Boss挑战卷轴
	ItemTypeRelicChallenge         int32 = 74    // 遗迹挑战卷轴
	ItemTypePetNormal              int32 = 80    // 普通宠物蛋
	ItemTypePetExp                 int32 = 83    // 宠物经验
	ItemTypeBossChallengeTimes     int32 = 90    // Boss挑战次数
	ItemTypeRelicChallengeTimes    int32 = 91    // 遗迹挑战次数
	ItemTypeGuildWish              int32 = 93    // 工会祈愿点
	ItemTypeAccelerateScroll       int32 = 99    // 加速卷轴
	ItemTypeExchangeCoupon         int32 = 25072 // 副本卷轴兑换券
	ItemTypeLuckyBagOne            int32 = 51004 // 福袋●壹
	ItemTypeLuckyBagTwo            int32 = 51005 // 福袋●贰
	ItemTypeLuckyBagThree          int32 = 51006 // 福袋●叁
	ItemTypePushWallChallengeTimes int32 = 51014 // 推墙挑战次数
	ItemTypePartnerTalentPoint     int32 = 52001 // 伙伴天赋点
	ItemTypeOfficialExp            int32 = 52002 // 官职经验
)

type EquipmentType int

// 定义装备类型的常量，使用iota自动递增。
const (
	Weapon        EquipmentType = iota + 1 // 1: 代表武器
	Clothes                                // 2: 代表衣服
	Helmet                                 // 3: 代表头盔
	Necklace                               // 4: 代表项链
	LeftRing                               // 5: 代表左戒指
	RightRing                              // 6: 代表右戒指
	LeftBracelet                           // 7: 代表左手镯
	RightBracelet                          // 8: 代表右手镯
	Belt                                   // 9: 代表腰带
	Boots                                  // 10: 代表靴子
	MaxEquipSlot
)

const (
	BuffTypeBuff   = 1 // 增益
	BuffTypeDebuff = 2 // 减益
)

const (
	BuffMarkPassive    int32 = 1    // 被动
	BuffTypeChainSkill int32 = 1000 // 关联技能数值
)

// BuffType 表示游戏中的增益或减益效果的枚举类型。
type BuffType int

// 定义BuffType的枚举值
const (
	// 持续性百分比回血
	BuffPercentageHeal BuffType = iota + 1
	// 持续性固定值回血
	BuffFixedHeal
	// 百分比加双防
	BuffPercentageDefenseBoost
	// 百分比加攻击
	BuffPercentageAttackBoost
	// 持续性百分比扣血
	BuffPercentageDamage
	// 持续性固定值扣血
	BuffFixedDamage
	// 百分比减防御
	BuffPercentageDefenseDecrease
	// 百分比增加攻速
	BuffPercentageAttackSpeedIncrease
	// 百分比增加移速
	BuffPercentageMovementSpeedIncrease
	// 百分比降低攻速
	BuffPercentageAttackSpeedDecrease
	// 百分比降低移速
	BuffPercentageMovementSpeedDecrease
	// 百分比减攻击
	BuffPercentageAttackDecrease
	// 百分比增加攻击力和双防
	BuffPercentageAttackBothDefenseIncrease
	// 百分比增加攻速移速
	BuffPercentageAttackMoveSpeedIncrease
	// 一次性百分比回血
	BuffFixedHealOnece
	// 无敌模式
	BuffGodMode

	// 百分比增加装备星级
	BuffPercentEquipStarIncrease = 50
	// 冰冻
	BuffFreeze = 51
	// 以角色的攻击力作为伤血计算基准值，走正常的战斗公式计算伤血效果
	BuffFixedDamage2 = 52
	// 嘲讽
	BuffTaunt = 53
	// 百分比最终减伤
	BuffDmtFinalReducePercentage = 54
	// 受控时间降低
	BuffTauntTimeReduce = 55
	// 反伤
	BuffReflectDamage = 56
	// 降低双防降低回血速度
	BuffPercentageDefenseHealthDecrease = 57
	// 代表替身宝宝技能
	BuffDecoy = 58
	// 代表护盾效果
	BuffShield = 59
	// 代表普攻带有的击退效果
	BuffKnockback = 60
	// 代表法师魔法酸雨技能
	BuffAcidRain = 61
	// 代表法师吸血盾
	BuffLifeLeechShield = 62
	// 挫折意志
	BuffFrustratedWill = 63
	// 代表技能群体净化术的效果
	BuffMassPurification = 64
	// 代表技能对怪伤害的百分比提升效果
	BuffIncreaseFixedMonster = 65
	// 代表宠物天赋中增加武器百分比攻击力的效果
	BuffPetWeaponAttackPercentageIncrease = 66
	// 百分比降低目标的暴击值
	BuffPercentageCritChanceDecrease = 67
	// 代表不受到任何负面效果的影响
	BuffNoNegativeEffect = 68
	// 百分比降低命中值
	BuffPercentageAccuracyDecrease = 69
	// 百分比提升暴击伤害
	BuffPercentageCritDamageIncrease = 70
	// 百分比提升命中值
	BuffPercentageAccuracyIncrease = 71
	// 百分比提升生命上限
	BuffPercentageMaxHealthIncrease = 72
	// 百分比提升PK伤害
	BuffPercentagePKDamageIncrease = 73
	// 74：代表攻击者百分比降低攻击、防御和暴击率
	BuffPercentageAttackDefenseCritDecrease = 74
	// 通用减益效果
	BuffCommonDebuff = 75
	// 通用增益效果
	BuffCommonBuff = 76
	// 77：代表受到伤害不会超出一定比例
	BuffDamageLimit = 77
	// 78：代表一场战斗中有一次复活的机会
	BuffReviveOnce = 78
	// 79：代表免疫效果
	BuffImmune = 79
	// 80 百分比提升生命上限
	BuffPercentageMaxHealthIncrease2 = 80
	// 龙卷风
	BuffTornado = 81
	// 按照剩余生命值百分比扣血
	BuffCurHpPercentageDamage = 82
	//禁锢
	BuffStun = 83
)

type BuffStatus int32

const (
	BuffStatusStart BuffType = iota + 1 // 开始
	BuffStatusStop                      // 结束
)

const (
	SkillEffectTypeAll     int32 = iota + 1 // 全部
	SkillEffectTypePlayer                   // 对人生效
	SkillEffectTypeMonster                  // 对怪生效
)

// 技能产生buff影响范围类型
const (
	// 自身
	SkillEffectnumSlef int32 = iota
	// 友方全体
	SkillEffectnumAlliance
	// 己方血量最低
	SkillEffectnumAllianceLowHp
	// 敌方一列(2个)
	SkillEffectnumEnemyColumn
	// 敌方一排(3个)
	SkillEffectnumEnemyRow
	// 敌方范围
	SkillEffectnumEnemyRange
	// 敌方单体
	SkillEffectnumEnemySingle
	// 敌方全体
	SkillEffectnumEnemyAll
	// 己方第一个死亡的角色
	SkillEffectnumAllianceFirstDie
	// 友方全体和行会队友
	SkillEffectnumAllianceAndGuild
	// 己方伙伴
	SkillEffectnumPartner = 10
)

// 技能类型
const (
	//===战士技能===
	SkillTypeNormalAttackWarrior int32 = iota + 1 // 战士普通攻击
	SkillTypeHalfMoon                             // 月弧斩
	SkillTypeFourSides                            // 四方斩
	SkillTypeBreakArmy                            // 破军斩
	SkillTypeIronBlood                            // 铁血狂潮
	SkillTypeWarGod                               // 战神荣耀
	SkillTypeSwiftStrike                          // 迅影追击 朝着当前的方向突进6米，并对路径上的敌对目标造成（110%攻击力+50点攻击力）的伤害，移动过程中为无敌状态。
	SkillTypeDuel                                 // 一决高下 对周围3米范围内的敌对目标发出嘲讽，使得目标有3%的几率处于嘲讽状态。（不能转换攻击目标和移动，普攻和技能均可使用）。
	SkillTypeDiscouragedWill                      // 挫折意志 释放后降低10%控制类技能的控制时间，并使得所受控制类技能的伤害降低9%，状态持续3秒。
	SkillTypeCounterStorm                         // 反击风暴 10%的几率将自身所受伤害的2%返还给攻击者，状态持续3秒
	SkillTypeFrenzyTouch                          // 狂暴之触
	SkillTypePlagueBlade                          // 瘟疫之刃
	SkillTypeStunSmash                            // 眩晕猛击

	//===法师技能===
	SkillTypeNormalAttackWizard // 法师普通攻击
	SkillTypeMagicMissile       // 魔法飞弹
	SkillTypeStellarBurst       // 星爆术
	SkillTypeAuroraPulse        // 极光脉冲
	SkillTypeTimeControl        // 时间操控
	SkillTypeMagicAmplification // 魔法增幅
	SkillTypeRepelFoe           // 拒敌千里 将身边2*2范围内的目标向外推出4米，并对目标造成（110%攻击力+40点攻击力）的伤害。
	SkillTypeDecoy              // 替身宝宝 在当前位置释放一个自身幻影（幻影无法移动攻击只能原地站立自身血量为本身血量的10%），本身则向当前朝向瞬移6米。
	SkillTypeInstantaneousMove  // 瞬息移动
	SkillTypeBloodthirstyCurse  // 嗜血之咒
	SkillTypeMagicBarrier       // 魔法结界
	SkillTypeInvincible         // 无懈可击
	SkillTypeManaExpel          // 魔力驱逐

	//===祭祀技能===
	SkillTypeNormalAttackRitualist // 祭祀普通攻击
	SkillTypeSoulTalisman          // 灵魂符咒
	SkillTypeEarthquakeWave        // 震地波
	SkillTypeNetherSpell           // 幽冥术
	SkillTypeDivineBattleArmor     // 神圣战甲术
	SkillTypeVenomGrimoire         // 毒蛊术
	SkillTypeToxicSand             // 毒砂术
	SkillTypeDivineStrength        // 神力精通
	SkillTypeArcaneMastery         // 魔力精通
	SkillTypeMassHeal              // 群体治愈术
	SkillTypeMassPurification      // 群体净化术
	SkillTypeLastGasp              // 回光返照
	SkillTypeBlessingOfCourage     // 勇气祝福

	// == boss 设计的技能 ==
	// 护盾
	SkillTypeShield = 51
	// 践踏
	SkillTypeStomp = 52
	// 吸取
	SkillTypeAbsorb = 53
	// 金身
	SkillTypeGoldenBody = 54
	// 重击
	SkillTypeHeavyBlow = 55
	// 分身
	SkillTypeDecoyBoss = 56
	// 斩杀
	SkillTypeExecute = 57
	// 嘲讽
	SkillTypeTaunt = 58
	// 怒气
	SkillTypeRage = 59
	// 卸甲
	SkillTypeDisarm = 60
	// 毁伤
	SkillTypeWreck = 61
	// 恐惧
	SkillTypeFear = 62
	// 中毒
	SkillTypePoison = 63
	// 致盲
	SkillTypeBlind = 64
	// 灭世
	SkillTypeApocalypse = 65
	// 出血
	SkillTypeBleeding = 66
	// 破兵
	SkillTypeBreakWeapon = 67
	// 回血
	SkillTypeHeal = 68
	// 群防
	SkillTypeGroupDefense = 69
	// 重生
	SkillTypeReborn = 70
	// 群体输出
	SkillTypeMassAttack = 71
	// 群体净化
	SkillTypeMassPurificationBoss = 72
	// 群体爆发
	SkillTypeMassBurst = 73
	// 回生术
	SkillTypeRebirth = 74
	// 专注
	SkillTypeConcentration = 75
	// 免疫
	SkillTypeImmune = 76
	// 强身
	SkillTypeBodyguard = 77
	// 重甲
	SkillTypeHeavyArmor = 78
	// 重伤
	SkillTypeHeavyWound = 79
	// 爆发
	SkillTypeBurstBoss = 80
	// 封锁
	SkillTypeSeal = 81
	// 战意
	SkillTypeBattleWill = 82
	// 护甲
	SkillTypeArmor = 83
	// 反击大事
	SkillTypeCounterAttack = 84
	// 自愈
	SkillTypeSelfHeal = 85
	// 强脉
	SkillTypeStrongPulse = 86

	// ===特殊技能===
	SkillTypeFreezeTime         = 500 // 千里冰封
	SkillTypeSwordHeavenBreaker = 501 // 剑开天门
	SkillTypeTimeSpaceLock      = 600 // 时空凝结
	SkillTypeMagicAcidRain      = 601 // 魔法酸雨

	// ===神器技能===
	SkillTypeStarShatter      = 3001 // 碎星击
	SkillTypeNinjaClone       = 3002 // 忍者分身
	SkillTypeBrambleCage      = 3003 // 荆棘牢笼
	SkillTypeGoddessDescent   = 3004 // 女神降临
	SkillTypeWaterDragonSlash = 3005 // 水龙斩
	SkillTypeWheelOfTime      = 3006 // 时空之轮
	SkillTypeVoidLava         = 3007 // 虚空熔岩
	SkillTypeVoidCollapse     = 3008 // 虚空坍缩
	SkillTypeSwordRain        = 3009 // 剑雨

	//===被动技能===
	SkillTypeAgile           = 27000 // 【灵活】
	SkillTypeResilient       = 28000 // 【坚韧】
	SkillTypeAbundant        = 29000 // 【厚泽】
	SkillTypePrecise         = 30000 // 【精准】
	SkillTypeBurst           = 31000 // 【爆发】
	SkillTypeValiant         = 32000 // 【勇猛】
	SkillTypeIronWallI       = 36000 // 【铁壁Ⅰ】
	SkillTypeSpiritBarrierI  = 37000 // 【灵障Ⅰ】
	SkillTypeIronWallII      = 38000 // 【铁壁Ⅱ】
	SkillTypeSpiritBarrierII = 39000 // 【灵障Ⅱ】
	SkillTypeSlaughterI      = 40000 // 【屠戮Ⅰ】
	SkillTypeEnemyBreakerI   = 41000 // 【破敌Ⅰ】
	SkillTypeSlaughterII     = 42000 // 【屠戮Ⅱ】
	SkillTypeEnemyBreakerII  = 43000 // 【破敌Ⅱ】
	SkillTypeLifeSourceI     = 44000 // 【命源Ⅰ】
	SkillTypeHealingSpiritI  = 45000 // 【愈灵Ⅰ】
	SkillTypeLifeSourceII    = 46000 // 【命源Ⅱ】
	SkillTypeHealingSpiritII = 47000 // 【愈灵Ⅱ】

	//===怪物技能===
	SkillTypeNormalAttackMonster = 9999  // 怪物普攻
	WorldBossSkill1              = 10000 // 世界boss技能
	ThorSkill0                   = 10001 // 雷神普攻
	ThorSkill1                   = 10002 // 雷神 1 技能
	ThorSkill2                   = 10003 // 雷神 2 技能
	OctopusSkill0                = 10004 // 章鱼普攻
	OctopusSkill1                = 10005 // 章鱼 1 技能
	OctopusSkill2                = 10006 // 章鱼 2 技能
	OctopusSkill4                = 10007 // 章鱼爆炸
	QiongqiSkill0                = 10008 // 穷奇普攻
	QiongqiSkill1                = 10009 // 穷奇 1 技能
	QiongqiSkill2                = 10010 // 穷奇 2 技能
	HeshangSkill0                = 10011 // 和尚普攻
	HeshangSkill1                = 10012 // 和尚 1 技能
	HeshangSkill2                = 10013 // 和尚 2 技能

	PartnerWarriorNormalSkillID   = 20000
	PartnerWizardNormalSkillID    = 20001
	PartnerRitualistNormalSkillID = 20002
)

const (
	SkillType               int32 = iota //
	SkillTypeNormal                      // 普攻
	SkillTypeCareer                      // 技能
	SkillTypeCareerSpecial               // 特殊技能技能
	SkillTypeCareerSpecial4              // 特殊技能技能4
	SkillTypeCareerSpecial5              // 特殊技能技能5
	SkillTypeArtifact                    // 神兵技能
)

const (
	SkillSlotType int32 = iota //
)

// 排行榜
const (
	RankTypeLevel            = 1  // 等级榜
	RankTypeFight            = 2  // 战力榜
	RankTypeSoul             = 3  // PVP胜场数量
	RankTypeTrial            = 4  // 试炼榜
	RankTypePvp              = 6  // 每日9点的计算(当天战魂)
	RankTypeSeason           = 7  // 赛季排行榜
	RankTypeWorldBoss        = 20 // 世界boss
	RankTypeWorldBossKf      = 21 // 世界boss的跨服榜
	RankTypeGightEnd         = 22 // 一战到底总榜
	RankTypeGightEndlayers   = 23 // 一战到底每一层榜单
	RankTypeTurnBased        = 24 // 回合制排行榜
	RankTypeSoulWeekly       = 25 // pvp胜场周榜
	RankTypeTurnBasedWeekly  = 26 // 回合制周榜
	RankTypeFightWeekly      = 27 // 战力周榜
	RankTypeElement          = 28 // 五行层数榜
	RankTypeCampaignStage    = 29 // 代表限时榜单中的战役通关榜单
	RankTypeMainBattle       = 30 // 代表限时榜单中的主线关卡通关榜单
	RankTypePinkEntry        = 31 // 粉色词条榜单
	RankTypeRedPartnerStar   = 32 // 红色伙伴总星级榜单
	RankTypeEquipmentStar    = 33 // 装备总升星榜单
	RankTypePetSynthesis     = 34 // 宠物合成榜单
	RankTypeSoulStoneLevel   = 35 // 魂石总等级榜单
	RankTypeFiveElementLevel = 36 // 五行赋能总等级榜单
	RankTypeSimpleFormation  = 37 // 简兵列阵排行榜
)

// 排行榜对应的中文名称
var RankTypeName = map[int]string{
	RankTypeLevel:            "等级日榜",
	RankTypeFight:            "战力日榜",
	RankTypeSoul:             "即时赛日榜",
	RankTypeTrial:            "试炼日榜",
	RankTypePvp:              "战魂日榜",
	RankTypeSeason:           "赛季排行榜",
	RankTypeWorldBoss:        "世界boss日榜",
	RankTypeWorldBossKf:      "世界boss跨服榜",
	RankTypeGightEndlayers:   "一战到底每一层榜单",
	RankTypeGightEnd:         "一战到底总榜",
	RankTypeTurnBased:        "回合制日榜",
	RankTypeSoulWeekly:       "即时赛周榜",
	RankTypeTurnBasedWeekly:  "回合制周榜",
	RankTypeFightWeekly:      "战力周榜",
	RankTypeElement:          "五行层数榜",
	RankTypeCampaignStage:    "战役通关榜",
	RankTypeMainBattle:       "主线关卡通关榜",
	RankTypePinkEntry:        "粉色词条榜",
	RankTypeRedPartnerStar:   "红色伙伴总星级榜",
	RankTypeEquipmentStar:    "装备总升星榜",
	RankTypePetSynthesis:     "宠物合成榜",
	RankTypeSoulStoneLevel:   "魂石总等级榜",
	RankTypeFiveElementLevel: "五行赋能总等级榜",
	RankTypeSimpleFormation:  "简兵列阵榜",
}

// 不发跑马灯的榜
var RankTypeServerNoLamp = map[int64]bool{
	RankTypeGightEndlayers: true, // 一战到底每一层榜单
	RankTypeCampaignStage:  true, // 战役通关榜
	RankTypeMainBattle:     true, // 主线关卡通关榜
}

// 需要分服务器的排行榜
var RankTypeServer = map[int64]bool{
	RankTypeFight:            true, // 战力
	RankTypeSoul:             true, // 即时赛榜
	RankTypeCampaignStage:    true, // 战役通关榜
	RankTypeMainBattle:       true, // 主线关卡通关榜
	RankTypePinkEntry:        true, // 粉色词条榜
	RankTypeRedPartnerStar:   true, // 红色伙伴总星级榜
	RankTypeEquipmentStar:    true, // 装备总升星榜
	RankTypePetSynthesis:     true, // 宠物合成榜
	RankTypeSoulStoneLevel:   true, // 魂石总等级榜
	RankTypeFiveElementLevel: true, // 五行赋能总等级榜
	RankTypeSimpleFormation:  true, // 简兵列阵榜
}

// 发奖励需要判断开服几天的榜
var RankTypeServerOpenDay = map[int64]bool{
	RankTypeCampaignStage: true, // 战役通关榜
	RankTypeMainBattle:    true, // 主线关卡通关榜
}

// 需要查询服务器名字的排行榜
var RankTypeNameServer = map[int64]bool{
	RankTypeWorldBoss: true, //
}

// 需要比较数据的排行榜
var RankTypeServerCompare = map[int64]bool{
	RankTypeLevel:          true, //
	RankTypeFight:          true, //
	RankTypeTrial:          true, //
	RankTypeWorldBoss:      true,
	RankTypeGightEndlayers: true,
}

// 日榜对应的周榜Id
var RankTypeWeekly = map[int64]int64{
	RankTypeFight:     RankTypeFightWeekly,
	RankTypeSoul:      RankTypeSoulWeekly,
	RankTypeTurnBased: RankTypeTurnBasedWeekly,
}

const (
	NeedMoney = 100 // 定义商城人民币支付的Id
	NeedAd    = 101 // 定义商城广告支付的Id
)

// 定义商城类型的枚举
const (
	ShopTypeUnknown               int32 = iota
	ShopTypeRecharge                    // 充值商城
	ShopTypeCoupon                      // 点券商城
	ShopTypeGold                        // 金币商城
	ShopTypeSoul                        // 战魂商城
	ShopTypeTrial                       // 试炼商城
	ShopTypeSignin30Days                // 30日累计签到的人民币购买
	ShopTypeActivityRecommendGift       // 活动中心的推荐礼包
	ShopTypeFirstCharge                 // 首冲三日活动
	ShopTypeSubscription                // 周、月、季卡
	ShopTypeBoost                       // 购买体力
	ShopTypeVIPGift                     // VIP礼包
	ShopTypeOfferReward                 // 悬赏
	ShopTypePVP                         // pvp购买次数
	ShopTypeRelic                       // 玩法的挑战次数购买
	ShopTypeTrialReward                 // 试炼基金
	ShopTypeGuild                       // 工会商店
	ShopTypePartner                     // 伙伴商店
	ShopTypePet                         // 宠物商店
	ShopTypeAFK                         // 挂机奖励
	ShopTypeDailyGift                   // 每日礼包
	ShopTypeSoulStoneGift               // 魂石礼包（宝石礼包）
	ShopTypeGrowthFund                  // 成长基金
	ShopTypeLimitedTimeGift             // 限时礼包
	ShopTypeElementGift                 // 24: 元素礼包
	ShopTypeSevenDaysFun                // 25: 七日活动专用
	ShopTypeGuardIsland                 // 26: 守卫海岛专用倍率
	ShopTypeAd                          // 27: 广告商店
	ShopTypeElementChallenge            // 28: 元素挑战次数
	ShopTypeEnergy                      // 29: 体力礼包
	ShopTypeExchange                    // 30兑换商店
	ShopTypeBless                       // 31福袋礼包
	ShopTypeOnHook                      // 32挂机商店
	ShopTypeAssistGift                  // 33代表助战礼包
	ShopTypeRefine                      // 34洗炼
	ShopTypeRefinePkg                   // 35洗炼礼包
	ShopTypeStarWishSummon              // 36: 星愿召唤礼包
	ShopTypeCrazyRefine                 // 37: 疯狂洗炼礼包
	ShopTypeRefineFund                  // 38: 洗练基金
	ShopTypeLotteryFund                 // 39: 抽奖基金
	ShopTypeInstanceFund                // 40: 副本基金
	ShopTypePinkEntry                   // 41: 粉色词条活动
	ShopTypeRedPartnerStar              // 42: 红色伙伴总星级活动
	ShopTypeEquipmentStar               // 43: 装备总升星活动
	ShopTypePetSynthesis                // 44: 宠物合成活动
	ShopTypeSoulStoneLevel              // 45: 魂石总等级活动
	ShopTypeFiveElementLevel            // 46: 五行赋能总等级活动

	ShopTypeAll
)

// 定义商城类型的枚举
const (
	ShopGoodsIdSpecial = 12001
)

// 定义活动类型的枚举
const (
	ActivityTypeUnknown          int32 = iota // 0: 未知
	ActivityTypeContinuity                    // 1: 连续登录
	ActivityTypeCumulative                    // 2: 累计登录
	ActivityTypeGrade                         // 3: 等级奖励
	ActivityTypeRecommend                     // 4: 推荐礼包
	ActivityTypeCharge                        // 5: 首充多日礼包
	ActivityTypeTrunTable                     // 6: 点券大转盘
	ActivityTypeBounty                        // 7: 限时悬赏活动
	ActivityTypeSingleCharge                  // 8: 单笔充值返利
	ActivityTypeCumulativeCharge              // 9: 累计充值返利
	ActivityTypeTrial                         // 10: 试炼塔玩法中的试炼基金
	ActivityTypePartnerGift                   // 11: 活动中心的伙伴礼包
	ActivityTypePetGift                       // 12: 活动中心的宠物礼包
	ActivityTypeDailyGift                     // 13: 每日礼包
	ActivityTypeSoulStoneGift                 // 14: 魂石礼包（宝石礼包）
	ActivityTypeGrowthFund                    // 15: 成长基金
	ActivityTypeLimitedTimeGift               // 16: 限时礼包
	ActivityTypeElementGift                   // 17: 元素礼包
	ActivityTypeSevenDaysFun                  // 18: 开服七天乐
	ActivityTypeSevenDaysFunGift              // 19: 开服七天乐礼包
	ActivityTypeAdGift                        // 20: 广告礼包
	ActivityTypeEnergyGift                    // 21: 体力礼包
	ActivityTypeLuckyBag                      // 22: 福袋礼包
	ActivityTypeTimeGift                      // 23: 时间礼包
	ActivityTypeAssist                        // 24: 助力活动
	ActivityTypeRefineGift                    // 25: 洗炼礼包
	ActivityTypeStarWishSummon                // 26: 星愿召唤活动
	ActivityTypeCrazyRefine                   // 27: 疯狂洗炼活动
	ActivityTypeRefine                        // 28: 洗练活动
	ActivityTypeLottery                       // 29: 抽奖活动
	ActivityTypeInstance                      // 30: 副本活动
	ActivityTypePinkEntry                     // 31: 粉色词条活动
	ActivityTypeRedPartnerStar                // 32: 红色伙伴总星级活动
	ActivityTypeEquipmentStar                 // 33: 装备总升星活动
	ActivityTypePetSynthesis                  // 34: 宠物合成活动
	ActivityTypeSoulStoneLevel                // 35: 魂石总等级活动
	ActivityTypeFiveElementLevel              // 36: 五行赋能总等级活动

)

var ShopMap map[int32]pb.TargetType

func init() {
	ShopMap = make(map[int32]pb.TargetType)
	ShopMap[ShopTypeCoupon] = pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_COUPONS
	ShopMap[ShopTypeGold] = pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_GOLD
	ShopMap[ShopTypeSoul] = pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_SOULS
	ShopMap[ShopTypeTrial] = pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_TRIAL_TOKENS
	ShopMap[ShopTypeBoost] = pb.TargetType_TOTAL_BOOST_PURCHASES
	ShopMap[ShopTypeExchange] = pb.TargetType_TICKET_EXCHANGE_TIMES
}

var ItemQuality = [...]string{
	"",
	"#ffffff",
	"#6fb469",
	"#6d99d0",
	"#a36ef6",
	"#e67be2",
	"#d29435",
	"#d55756",
}

// 定义
const (
	EquipMasterCategoryUnknown        int32 = iota
	EquipMasterCategoryEnhance              // 强化
	EquipMasterCategoryStars                // 升星
	EquipMasterCategoryGemStone             // 魂石大师
	EquipMasterCategoryGemStoneNature       // 自然加成
)

// 对应的中文名字
var EquipMasterCategoryName = map[int32]string{
	EquipMasterCategoryEnhance:        "强化",
	EquipMasterCategoryStars:          "升星",
	EquipMasterCategoryGemStone:       "魂石大师",
	EquipMasterCategoryGemStoneNature: "自然加成",
}

const (
	GameOverType          int32 = iota
	GameOverTypeTime            // 时间到
	GameOverTypePDead           // 死亡复活超时
	GameOverTypeMDead           // 怪死光
	GameOverTypeDeadFaild       // 死亡失败
)

// LampType 表示玩家事件的类型。
type LampType int

const (
	// 1. 排行榜榜首变化
	LampTopRankChanged LampType = iota + 1 // 玩家暂列XX排行榜榜首，可喜可贺！

	// 2. 乱斗之王玩法中连续获得10场胜利
	LampUnstoppableInBrawl // 玩家在乱斗之王玩法中已经暴走，无人可挡！

	// 3. 试炼塔玩法中通过10的整数倍层数
	LampTrialTowerLevelUp // 玩家通过了试炼塔第XX层，太厉害了！

	// 4. 战役玩法中通过每一章的最后一关
	LampCampaignChapterCompleted // 玩家通过了战役第X章，一往无前！

	// 5. 首领玩法中通过5的整数倍关数
	LampBossChapterCompleted // 玩家通过了首领第X关，收获了大量奖励！

	// 6. 遗迹玩法中领取每一章的30星宝箱奖励
	LampRelicChapterReward // 玩家领取了遗迹星级宝箱获取珍稀奖励，令人羡慕！

	// 7. 购买试炼塔的试炼基金
	LampTrialFundPurchased // 玩家购买了试炼基金，一次性领取了海量大奖！

	// 8. 伙伴召唤中抽取到橙色或红色品质的整卡英雄
	LampHeroSummon // 玩家运气爆棚，召唤出了（伙伴名称）！

	// 9. 宠物蛋抽取得到橙色品质以上的尊品宠物
	LampPetSummon // 玩家手气奇佳，召唤出了（宠物名称）！

	// 10. 宠物蛋合成得到橙色品质以上的尊品宠物
	LampPetComposeSummon // 玩家手气奇佳，合成出了（宠物名称）！

	// 11. 玩家成功进行了首充并获得了首充奖励
	LampFirstRecharge // 玩家获得了三日首充大礼，赢在起点！

	// 12. 玩家成功购买了特权卡
	LampPrivilegeCardPurchased // 玩家购买了X卡，获得了超值豪礼！

	// 13. 世界boss开启挑战
	LampWorldBossChallenge // 世界boss挑战开启，海量奖励等你来拿！

	// 14. 装备槽位升星到指定星级
	LampEquipmentStarUpgrade // 玩家成功打造X星槽位，傲视群雄！

	// 15. 合成出了8级以上且带有自然属性的宝石
	LampGemSynthesis // 玩家幸运的合成出了X级X魂石（X），可喜可贺！

	// 16. 行会会长上线
	LampGuildMasterLogin // 会长XXX已上线！

	// 17. 榜单排名第一的玩家上线
	LampTopRankOnline // XX榜第一名的玩家（角色名称）上线了，大家快来围观！
	// 18. common
	LampCommon
	// 19. 连斩
	LampConsecutiveKills // 玩家{0}在战场中已经达成了{1}，谁能与之一战？
)

// 登录方式
const (
	PlatformGuest  = "guest"
	PlatformUnity  = "unity"
	PlatformWechat = "wechat"
	PlatformMobile = "mobile"
)

// 首页怪物
const (
	MaintriggerType int32 = iota //
	MaintriggerNormal
	MaintriggerMonster
)

// 抽奖类型
const (
	LotteryType int32 = iota //
	LotteryTypeBronze
	LotteryTypeSilver
	LotteryTypeGold
	LotteryTypePlatinum
	LotteryMax
)

// 中奖类型
const (
	LotteryAwardType int32 = iota //
	LotteryAwardFragment
	LotteryAwardWhole
	LotteryAwardMax
)

// 中奖模块
const (
	LotteryModuleType int32 = iota //
	LotteryModulePartner
	LotteryModulePet
	LotteryModuleMax
)

// 邮件内容
const (
	EmailContentSystem  = "系统邮件"
	EmailContentBagFull = "包裹已满"
	EmailContentGuidMap = "行会首领击杀奖励"
)

// Docx 枚举定义了游戏中不同类型的文档。用于活动
type Docx int

const (
	DocxActivity Docx = iota + 1 // 活动
)

type PayPlatformType string

const (
	PayPlatformTypeUnknown       PayPlatformType = "unknown"
	PayPlatformTypeWeChatAndroid PayPlatformType = "wc_android"
	PayPlatformTypeWeChatIOS     PayPlatformType = "wc_ios"
)

func (p PayPlatformType) String() string {
	return string(p)
}

// 修改邮件的计数器，操作为设置，增加，减少
type MailCounterOPType int

const (
	MailCounterOPTypeSet MailCounterOPType = iota + 1 // 设置
	MailCounterOPTypeAdd                              // 增加
	MailCounterOPTypeSub                              // 减少
)

const (
	EnvDev    = "dev"    // 开发环境
	EnvTest   = "test"   // 内网测试环境
	EnvProd   = "prod"   // 生产环境
	EnvStag   = "stag"   // 外网测试(体验服)
	EnvReview = "review" // 提审服
)

// GuardIsLandType 守卫海岛的类型*(挂机副本)
type GuardIsLandType int8

const (
	GuardIsLandTypeFree     GuardIsLandType = iota + 1 // 免费挂机
	GuardIsLandTypeAd                                  // 广告挂机
	GuardIsLandTypeItem                                // 道具挂机
	GuardIsLandTypeRecharge                            // 充值挂机
)

const (
	PaymentURL             = "https://payment.shbinqigame.com"
	ByteDanceConversionURL = "https://analytics.oceanengine.com/api/v2/conversion"
)

// 元素强度对应的抗性
var (
	Attribute_emphasis = map[pb.Attribute]pb.Attribute{
		pb.Attribute_GOLD_EMPOWER:  pb.Attribute_GOLD_RESISTANCE,
		pb.Attribute_WOOD_EMPOWER:  pb.Attribute_WOOD_RESISTANCE,
		pb.Attribute_WATER_EMPOWER: pb.Attribute_WATER_RESISTANCE,
		pb.Attribute_FIRE_EMPOWER:  pb.Attribute_FIRE_RESISTANCE,
		pb.Attribute_EARTH_EMPOWER: pb.Attribute_EARTH_RESISTANCE}
)

// 槽位对应id
var (
	ElementSlot_attribute = map[pb.ElementSlot]pb.Attribute{
		pb.ElementSlot_GoldEmpower:  pb.Attribute_GOLD_EMPOWER,
		pb.ElementSlot_WoodEmpower:  pb.Attribute_WOOD_EMPOWER,
		pb.ElementSlot_WaterEmpower: pb.Attribute_WATER_EMPOWER,
		pb.ElementSlot_FireEmpower:  pb.Attribute_FIRE_EMPOWER,
		pb.ElementSlot_EarthEmpower: pb.Attribute_EARTH_EMPOWER}
)

// 元素对应克制属性
var (
	Attribute_know_attribute = map[pb.Attribute]pb.Attribute{
		pb.Attribute_GOLD_EMPOWER:  pb.Attribute_WOOD_EMPOWER,
		pb.Attribute_WOOD_EMPOWER:  pb.Attribute_EARTH_EMPOWER,
		pb.Attribute_EARTH_EMPOWER: pb.Attribute_WATER_EMPOWER,
		pb.Attribute_WATER_EMPOWER: pb.Attribute_FIRE_EMPOWER,
		pb.Attribute_FIRE_EMPOWER:  pb.Attribute_GOLD_EMPOWER}
)

const (
	//0：代表没有触发内容
	MonsterTrigerTypeNone int32 = iota
	//1：代表触发了随机特定品质装备的怪物，此类怪物需要将掉落装备乘上怪物的品质系数
	MonsterTrigerTypeEquipFuDai
)

// 新手引导特殊的类型
const (
	GuidanceTypeBoxOne  = 8
	GuidanceTypeBoxTwo  = 10
	GuidanceTypePartner = 18
)

// 功能模块枚举
type ModuleID int32

const (
	ModuleBoxPlay         ModuleID = iota + 1 // 1. 开箱玩法（宝箱开启+宝箱升级+宝箱设置）
	ModuleTask                                // 2. 任务模块（主线+日常+周常+成就+福利）
	ModuleDailySign                           // 3. 每日签到
	ModuleEventCenter                         // 4. 活动中心（七日连登+30日累登+等级奖励+折扣礼包）
	ModuleFirstRecharge                       // 5. 首充功能
	ModulePrivilegeCard                       // 6. 特权卡功能
	ModuleVIP                                 // 7. VIP功能
	ModuleSystem                              // 8. 系统功能（设置+邮件+激活码）
	ModuleMainRole                            // 9. 角色模块-主角（装备+技能）
	ModuleBag                                 // 10. 包裹模块
	ModuleShop                                // 11. 商城模块（充值商城+元宝商城+金币商城）
	ModuleChat                                // 12. 聊天模块
	ModulePartner                             // 13. 角色模块-伙伴（召唤+技能）
	ModuleExpPool                             // 14. 经验池功能（主角+伙伴升级）
	ModuleCampaign                            // 15. 玩法模块-战役
	ModulePartnerStarUp                       // 16. 伙伴升星
	ModuleBoss                                // 17. 玩法模块-首领
	ModuleTrialShop                           // 18. 商城模块（试炼商城）
	ModuleTrialTower                          // 19. 玩法模块-试炼塔
	ModulePet                                 // 20. 战宠模块（抽取+合成+突破）
	ModuleRelic                               // 21. 玩法模块-遗迹
	ModuleSoulShop                            // 22. 商城模块（战魂商城）
	ModulePK1v1                               // 23. 玩法模块-限时-1V1PK
	ModulePKMulti                             // 24. 玩法模块-限时-多V多PK
	ModuleEnhanceFB                           // 25. 玩法模块-副本-强化石副本
	ModuleEquipStrengthen                     // 26. 打造模块-装备强化
	ModuleStarsFB                             // 27. 玩法模块-副本-升星石副本
	ModuleEquipStarUp                         // 28. 打造模块-装备升星
	ModuleBookFB                              // 29. 玩法模块-副本-书页副本
	ModuleSkillUpgrade                        // 30. 打造模块-技能升级
	ModuleSoulFB                              // 31. 玩法模块-副本-魂石副本
	ModuleEquipSoul                           // 32. 打造模块-魂石
	ModuleGuild                               // 33. 家族模块（家族升级+家族技能+家族副本+家族商店）
	ModuleSkillAdvance                        // 34. 打造模块-技能进阶
	ModuleGuardIsland                         // 35. 玩法模块-守卫海岛
	ModuleWorldBoss                           // 36. 世界boss
	ModuleOneBattle                           // 37. 一战到底
	ModuleElementTrial                        // 38. 五行试炼
	ModuleBattlefield                         // 39. 群雄战场
	ModuleTitle                               // 40. 角色模块-称号
	ModuleChatFunc                            // 41. 聊天功能
	ModuleDoubleSpeed                         // 42. 双倍速功能
	ModuleAutoBattle                          // 43. 主界面自动战斗
	ModuleLottery                             // 44. 抽奖模块
	ModuleAvatar                              // 45. 头像-头像框
	ModelRefined                              // 46. 槽位洗炼
	ModuleExpFB                               // 47. 经验副本
	ModuleCoinFB                              // 48. 金币副本
	ModulePushingWall                         // 49. 这怪有点强
	ModuleArtifact                            // 50. 神器
)

// ModuleID 对应的中文名称
var ModuleName = map[ModuleID]string{
	ModuleBoxPlay:         "开箱玩法",
	ModuleTask:            "任务模块",
	ModuleDailySign:       "每日签到",
	ModuleEventCenter:     "活动中心",
	ModuleFirstRecharge:   "首充功能",
	ModulePrivilegeCard:   "特权卡功能",
	ModuleVIP:             "VIP功能",
	ModuleSystem:          "系统功能",
	ModuleMainRole:        "角色模块-主角",
	ModuleBag:             "包裹模块",
	ModuleShop:            "商城模块",
	ModuleChat:            "聊天模块",
	ModulePartner:         "角色模块-伙伴",
	ModuleExpPool:         "经验池功能",
	ModuleCampaign:        "玩法模块-战役",
	ModulePartnerStarUp:   "伙伴升星",
	ModuleBoss:            "玩法模块-首领",
	ModuleTrialShop:       "商城模块-试炼商城",
	ModuleTrialTower:      "玩法模块-试炼塔",
	ModulePet:             "战宠模块",
	ModuleRelic:           "玩法模块-遗迹",
	ModuleSoulShop:        "商城模块-战魂商城",
	ModulePK1v1:           "玩法模块-1V1PK",
	ModulePKMulti:         "玩法模块-多V多PK",
	ModuleEnhanceFB:       "玩法模块-强化石副本",
	ModuleEquipStrengthen: "打造模块-装备强化",
	ModuleStarsFB:         "玩法模块-升星石副本",
	ModuleEquipStarUp:     "打造模块-装备升星",
	ModuleBookFB:          "玩法模块-书页副本",
	ModuleSkillUpgrade:    "打造模块-技能升级",
	ModuleSoulFB:          "玩法模块-魂石副本",
	ModuleEquipSoul:       "打造模块-魂石",
	ModuleGuild:           "家族模块",
	ModuleSkillAdvance:    "打造模块-技能进阶",
	ModuleGuardIsland:     "玩法模块-守卫海岛",
	ModuleWorldBoss:       "世界boss",
	ModuleOneBattle:       "一战到底",
	ModuleElementTrial:    "五行试炼",
	ModuleBattlefield:     "群雄战场",
	ModuleTitle:           "角色模块-称号",
	ModuleChatFunc:        "聊天功能",
	ModuleDoubleSpeed:     "双倍速功能",
	ModuleAutoBattle:      "主界面自动战斗",
	ModuleLottery:         "抽奖模块",
	ModuleAvatar:          "头像-头像框",
	ModelRefined:          "槽位洗炼",
	ModuleExpFB:           "经验副本",
	ModuleCoinFB:          "金币副本",
	ModulePushingWall:     "这怪有点强",
}

// String 返回模块ID对应的中文名称
func (m ModuleID) String() string {
	if name, ok := ModuleName[m]; ok {
		return name
	}
	return "未知模块"
}

// "lottery_num" 可能需要定义为常量
const LotteryNum = "lottery_num"
const RefineNum = "refine_num"
const InstanceNum = "instance_num"
const CumulativeChargeAmount = "cumulative_charge_amount"
const PinkEntryNum = "pink_entry_num"
const RedPartnerStarNum = "red_partner_star_num"
const EquipmentStarNum = "equipment_star_num"
const PetSynthesisNum = "pet_synthesis_num"
const SoulStoneLevelNum = "soul_stone_level_num"
const FiveElementLevelNum = "five_element_level_num"
