package thinkingdatax

// UserAttribute 表示用户属性的枚举类型，底层为字符串，用于标识用户画像或行为数据中的字段。
type UserAttribute string

const (
	Channel               UserAttribute = "channel"                  // 来源渠道：第三方归因平台导入，默认官方微信
	System                UserAttribute = "system"                   // 操作系统：ios/android
	CurrentLevel          UserAttribute = "current_level"            // 当前等级：玩家等级变动时上报
	CurrentVip            UserAttribute = "current_vip"              // 当前vip等级：玩家VIP变动时上报
	CurrentIngot          UserAttribute = "current_ingot"            // 当前元宝：在用户登出前上报
	CurrentIngotCost      UserAttribute = "current_ingot_cost"       // 当前累计消耗元宝：在用户登出前上报
	CurrentCoin           UserAttribute = "current_coin"             // 当前金币：在用户登出前上报
	CurrentCoinCost       UserAttribute = "current_coin_cost"        // 当前累计消耗金币：在用户登出前上报
	FirstPayTime          UserAttribute = "first_pay_time"           // 首次付费时间：首次付费时上报
	LastLoginTime         UserAttribute = "last_login_time"          // 最后登出时间：在用户登出前上报
	CurrentMainId         UserAttribute = "current_main_id"          // 当前主线关卡进度：每次通关新的主线关卡时记录
	CurrentBoss           UserAttribute = "current_boss"             // 当前boss关卡进度：每次通关新的boss时记录
	CurrentTrial          UserAttribute = "current_trial"            // 最高试炼塔进度：每次刷新试炼塔最高层数时记录
	CurrentRuins          UserAttribute = "current_ruins"            // 当前遗迹关卡进度：每次通关新的遗迹关卡时记录
	RegisterTime          UserAttribute = "register_time"            // 用户新增时间：账号新增时设置
	TotalPayAmount        UserAttribute = "total_pay_amount"         // 累计付费金额：每次充值时上报
	FirstRegistrationTime UserAttribute = "first_registration_time"  // 用户注册时间（时间）：在用户注册时上报
	AdTotalToday          UserAttribute = "ad_total_today"           // 今日总看广告次数：在用户登出前上报
	CurrentPower          UserAttribute = "current_power"            // 当前战力：在用户登出前上报
	CurrentGuildId        UserAttribute = "current_guild_id"         // 当前联盟id：每次联盟变化时上报
	CurrentDay            UserAttribute = "current_day"              // 当前累计游戏天数：当日首次登录时上报
	CurrentTime           UserAttribute = "current_time"             // 当前累计游戏时长：在用户登出前上报
	CurrentBoxLevel       UserAttribute = "current_box_level"        // 当前宝箱等级：宝箱升级时上报
	CurrentBoxes          UserAttribute = "current_boxes"            // 当前累计开宝箱数：在用户登出前上报
	CurrentBoxesLeft      UserAttribute = "current_boxes_left"       // 当前剩余宝箱数量：在用户登出前上报
	Ip                    UserAttribute = "ip"                       // ip属地：账号新增时设置
	AccountId             UserAttribute = "account_id"               // 游戏账号：账号新增时设置
	DeviceId              UserAttribute = "device_id"                // 设备id：账号新增时设置
	RoleId                UserAttribute = "role_id"                  // 角色id：创建角色时上报
	RoleJob               UserAttribute = "role_job"                 // 角色职业：创建角色时上报
	RoleNum               UserAttribute = "role_num"                 // 第几个角色?：创建角色时上报
	CurrentBattleSoul     UserAttribute = "current_battle_soul"      // 当前战魂数量：在用户登出前上报
	CurrentBattleSoulCost UserAttribute = "current_battle_soul_cost" // 当前累计消耗战魂：在用户登出前上报
	CurrentTrialCoin      UserAttribute = "current_trial_coin"       // 当前试炼币数量：在用户登出前上报
	CurrentTrialCoinCost  UserAttribute = "current_trial_coin_cost"  // 当前累计消耗试炼币：在用户登出前上报
	CurrentGuildCoin      UserAttribute = "current_guild_coin"       // 当前行会币数量：在用户登出前上报
	CurrentGuildCoinCost  UserAttribute = "current_guild_coin_cost"  // 当前累计消耗行会币：在用户登出前上报
	AllPartnerLevel       UserAttribute = "all_partner_level"        // 当前伙伴总等级：在用户登出前上报
	IsMonthCard           UserAttribute = "is_month_card"            // 是否月卡：月卡状态变更时上报
	IsLifetimeCard        UserAttribute = "is_lifetime_card"         // 是否终身卡：终身卡状态购买后上报
	CurrentPetZun         UserAttribute = "current_pet_zun"          // 当前宠物是否为尊：在用户登出前上报
	ServerId              UserAttribute = "server_id"                // 服务器id：创建角色时上报
	ImproveResonance      UserAttribute = "improve_resonance"        // 强化大师等级：共鸣等级升级时上报
	StarResonance         UserAttribute = "star_resonance"           // 升星大师等级：共鸣等级升级时上报
	StoneResonance        UserAttribute = "stone_resonance"          // 魂石大师等级：共鸣等级升级时上报
	Nickname              UserAttribute = "nickname"                 // 角色昵称：创建角色和修改时上报
	NicknameOne           UserAttribute = "nickname_one"             // 第一个角色昵称：创建第一个角色时上报
	NicknameTwo           UserAttribute = "nickname_two"             // 第二个角色昵称：创建第二个角色时上报
	RoleIdOne             UserAttribute = "role_id_one"              // 第一个角色访客id：创建第一个角色时上报
	RoleIdTwo             UserAttribute = "role_id_two"              // 第二个角色访客id：创建第二个角色时上报
	RoleJobOne            UserAttribute = "role_job_one"             // 第一个角色职业：创建第一个角色时上报
	RoleJobTwo            UserAttribute = "role_job_two"             // 第二个角色职业：创建第二个角色时上报
	PhoneNumber           UserAttribute = "phone_number"             // 手机号：授权手机号时上报
	DeviceModel           UserAttribute = "device_model"             // 手机型号：授权时上报
	CampaignId            UserAttribute = "campaign_id"              // 广告id：账号新增时设置
	AdTotal               UserAttribute = "ad_total"                 // 总看广告次数：在用户登出前上报
	IsPc                  UserAttribute = "is_pc"                    // 是否电脑端：在用户登出前上报
	TmpPvp1Rank           UserAttribute = "tmp_pvp1_rank"            // 当前即时 PVP排名：在用户登出前上报
	TmpPvp2Rank           UserAttribute = "tmp_pvp2_rank"            // 当前回合 PVP排名：在用户登出前上报
	TotalCharacterSpend   UserAttribute = "total_character_spend"    // 本次充值后角色的总消费金额（单位：分）：用户支付后上报
	AdSource              UserAttribute = "ad_source"                // 广告进入的用户标识：角色创建时上报（腾讯：gdt，字节：ocean）
	LoginTimeV2           UserAttribute = "login_time_v2"            // 最近登录时间：用户登录时上报
)

// String 实现 fmt.Stringer 接口，返回属性名的字符串形式
func (ua UserAttribute) String() string {
	return string(ua)
}

// Event 表示埋点事件的枚举类型，底层为字符串，用于标识用户行为或系统事件。
type Event string

const (
	TaMpLaunch              Event = "ta_mp_launch"              // 小程序初始化：自动采集事件
	TaMpShow                Event = "ta_mp_show"                // 小程序启动：在小程序被启动、或从后台调回前台时触发
	TaMpHide                Event = "ta_mp_hide"                // 小程序隐藏：被调入后台时触发，并记录本次使用时长
	TaMpShare               Event = "ta_mp_share"               // 小程序分享：自动采集事件
	NewDevice               Event = "new_device"                // 设备激活：设备首次启动时通过首次事件上报
	Register                Event = "register"                  // 用户注册：用户注册账号后上报
	Login                   Event = "login"                     // 用户登录：用户登录账号后上报
	Logout                  Event = "logout"                    // 用户登出：用户登出账号后上报
	CreateRole              Event = "create_role"               // 创建角色：用户创建角色后上报
	RoleLogin               Event = "role_login"                // 登录角色：用户登录角色后上报
	RoleLogout              Event = "role_logout"               // 登出角色：用户登出角色后上报
	GuideCompleted          Event = "guide_completed"           // 新手引导：用户完成阶段新手引导后上报
	Online                  Event = "online"                    // 在线数据：通过服务端上报，以任意唯一ID生成虚拟用户进行在线人数上报（间隔如1分钟、5分钟等）
	VipUpgrad               Event = "vip_upgrad"                // vip升级：vip升级时上报
	TaskComplete            Event = "task_complete"             // 任务完成：完成任务后上报
	RoleLevelUpgrade        Event = "role_level_upgrade"        // 主角升级：升级后上报
	PartnerLevelUpgrade     Event = "partner_level_upgrade"     // 伙伴升级：升级后上报
	PartnerStarUpgrade      Event = "partner_star_upgrade"      // 伙伴升星：升星后上报
	InitOrder               Event = "init_order"                // 发起订单：用户发起充值订单
	OrderFinish             Event = "order_finish"              // 付费完成：用户完成充值后上报
	OpenBox                 Event = "open_box"                  // 宝箱获取：开箱时上报
	BoxLevelUp              Event = "box_level_up"              // 宝箱升级：宝箱升级完成时上报
	PetCatch                Event = "pet_catch"                 // 宠物获得：获得宠物时上报
	PetGobble               Event = "pet_gobble"                // 宠物升级：升级时上报
	PetSynthesis            Event = "pet_synthesis"             // 宠物合成：合成时上报
	PetOvershoot            Event = "pet_overshoot"             // 宠物突破：使用突破石时上报
	PetEquip                Event = "pet_equip"                 // 宠物上阵：上阵变更时上报
	EquipGet                Event = "equip_get"                 // 装备获取（宝箱来源以外）：从非宝箱来源获取装备时上报
	EquipPositionStrengthen Event = "equip_position_strengthen" // 装备位强化：强化装备位时上报（如+1、+2、+3）
	EquipStarUp             Event = "equip_star_up"             // 装备升星：装备升星时上报
	EquipInlay              Event = "equip_inlay"               // 魂石镶嵌：魂石镶嵌时上报
	SkillLevelUp            Event = "skill_level_up"            // 升级技能：升级技能时上报
	SkillUpgrade            Event = "skill_upgrade"             // 技能升阶：技能升阶时上报
	ChallengeBattleStart    Event = "challenge_battle_start"    // 挑战玩法挑战开始：挑战开始时上报（主线、五行塔等合并）
	ChallengeBattleEnd      Event = "challenge_battle_end"      // 挑战战役挑战结束：挑战结束后上报（主线、五行塔等合并）
	RealtimePvpStart        Event = "realtime_pvp_start"        // 开始pvp挑战：开始PVP时上报
	RealtimePvpEnd          Event = "realtime_pvp_end"          // pvp挑战结束：PVP结束后上报
	ShopBuyPackage          Event = "shop_buy_package"          // 商店购买礼包：商店购买礼包后上报（非人民币购买）
	ShopBuyItem             Event = "shop_buy_item"             // 商店购买单品：商店购买单品后上报（非人民币购买）
	SummonParnter           Event = "summon_parnter"            // 召唤伙伴：召唤伙伴时上报
	OtherItemAcquire        Event = "other_item_acquire"        // 其他物品获得：物品产出时上报（开宝箱获得的装备不上报）
	OtherItemConsume        Event = "other_item_consume"        // 其他物品消耗：物品消耗时上报（开宝箱消耗的箱子不上报）
	CreateGuild             Event = "create_guild"              // 创建行会：创建行会时上报
	JoinGuild               Event = "join_guild"                // 加入行会：加入行会时上报
	GuildSkill              Event = "guild_skill"               // 行会技能：行会技能祈愿时上报
	GuildDonate             Event = "guild_donate"              // 行会捐献：行会捐献时上报
	GuardingIslandStart     Event = "guarding_island_start"     // 守卫海岛开始：守卫海岛开始时上报
	GuardingIslandEnd       Event = "guarding_island_end"       // 守卫海岛结束：守卫海岛结束或中断时上报
	CheckIn                 Event = "check_in"                  // 签到：完成签到时上报（现 ptype）
	WorldBossStart          Event = "world_boss_start"          // 开始世界boss：开始世界boss时上报
	WorldBossEnd            Event = "world_boss_end"            // 完成世界boss：完成世界boss时上报
	YizhandaodiStart        Event = "yizhandaodi_start"         // 一站到底开始：开始一站到底时上报
	YizhandaodiEnd          Event = "yizhandaodi_end"           // 一站到底结束：一站到底打完一局时上报
	MailReceive             Event = "mail_receive"              // 收到邮件：收到邮件时上报
	MailOpen                Event = "mail_open"                 // 领取邮件：领取邮件时上报
	ChallengeList           Event = "challenge_list"            // 打开挑战列表：打开PVP挑战列表时上报
	Resonance               Event = "resonance"                 // 共鸣等级：共鸣等级提升时上报
	AdClick                 Event = "ad_click"                  // 点击广告：点击广告时上报
	AdShow                  Event = "ad_show"                   // 广告曝光：广告曝光时（变现广告）
	AdClose                 Event = "ad_close"                  // 广告关闭：广告被关闭时（变现广告）
	ClickEvent              Event = "click_event"               // 点击事件：玩家游戏过程中，具体点击了哪些事件
	WuxingGet               Event = "wuxing_get"                // 五行赋能：五行赋能时上报（需确认消耗道具及位置命名）
	InlaySynthesis          Event = "inlay_synthesis"           // 魂石合成：魂石合成时上报
	ItemSynthesis           Event = "item_synthesis"            // 道具合成：道具合成时上报
	CraftSynthesis          Event = "craft_synthesis"           // 材料合成：看后续是否拓展，目前游戏内为空
	BondActivate            Event = "bond_activate"             // 羁绊激活：激活羁绊时上报
	QuitGuild               Event = "quit_guild"                // 退出行会：退出行会时上报
	ChatMessage             Event = "message"                   // 聊天：聊天频道发言时上报
	CDkey                   Event = "CDkey"                     // 激活码：激活码兑换时上报
	VipCard                 Event = "vip_card"                  // 特权卡：领取特权卡奖励时上报
	WuxingTowerStart        Event = "wuxing_tower_start"        // 五行塔挑战开始：挑战五行塔时上报
	WuxingTowerEnd          Event = "wuxing_tower_end"          // 五行塔挑战结束：挑战结束后上报
	MainPve                 Event = "main_pve"                  // 主线挑战结束：主线挑战结束后上报
	EquipRefining           Event = "equip_refining"            // 装备洗练：装备洗练时上报
	QunxiongWar             Event = "qunxiong_war"              // 群雄战场：挑战结束后上报
	Combo                   Event = "combo"                     // 组合技：发生连击行为时上报（含系统默认滚动）
	ComboClickNumberOne     Event = "combo_click_number_one"    // Combo首个按钮点击：玩家点击首个按钮触发combo
	ComboAppear             Event = "combo_appear"              // Combo出现：Combo效果出现时上报
	ComboFirstFail          Event = "combo_first_fail"          // Combo首次失败：Combo首次未成功触发时上报
)

// String 实现 fmt.Stringer 接口
func (e Event) String() string {
	return string(e)
}
