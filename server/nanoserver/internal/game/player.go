package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/dispatcher"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/internal/web/provide"
	"nanoserver/protocol"

	"nanoserver/pkg/async"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/event"
	"nanoserver/pkg/maputil"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"nanoserver/pkg/utils"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano/serialize/protobuf"
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
)

type (
	PlayerSkill struct {
		CfgId         int32     `json:"cfgid"`
		AdvancedLevel int32     `json:"advanced_level"`
		Cooldown      time.Time `json:"-"`
		Distance      int32     `json:"-"`
	}

	PVPInfo struct {
		*PVPDuel
		*PVPMelee
		dailySoul   int // 每天pvp累计获得的战魂
		sessionSoul int // 每赛季pvp累计获得的战魂

		pvpWinStreakTimes int // 所有模式连胜次数

		pvpRecorder []*model.PvpRecord // PVP记录
	}

	PVPDuel struct {
		pvpScore           int32 // pvp分数
		dailyMatchTimes    int32
		dailyBuyNum        int32
		losingStreakCount  int32 // 1v1连胜
		winningStreakCount int32 // 1v1连败
	}

	PVPMelee struct {
		gvgScore        int32 // gvg分数
		dailyMatchTimes int32
		dailyBuyNum     int32
	}

	ElementSlot struct {
		Eid    int32 `json:"eid"`
		SlotId int32 `json:"id"`    // 槽位Id
		Level  int32 `json:"level"` // 元素的等级
		Break  int32 `json:"break"` //突破次数
	}
	TurnBased struct {
		turnBasedSlot [TurnBasedMax]int64
	}

	Slot struct {
		ItemId   int64                                            `json:"id"`       // 装备的唯一标识符
		Stars    int32                                            `json:"stars"`    // 装备的星级
		Enhance  int32                                            `json:"enhance"`  // 装备的强化等级
		Gemstone [pb.GemstoneSlot_MaxGemstoneSlot]pb.GemstoneSlot `json:"gemstone"` // 装备的宝石
		Quality  int32                                            `json:"quality"`  // 装备的品质
	}

	// 最后一次洗练的属性
	LastEquipRefine struct {
		itemId  int64
		attr    pb.Attribute
		random  bool
		newStat pb.Stats
	}

	// 无尽试炼
	SoloTowerInfo struct {
		Layers     int32
		RewardTime time.Time
	}

	// 世界boss
	WorldBossInfo struct {
		dailyWorldBossTimes  int32 // 当天已经挑战次数
		totalWorldBossDamage int64 // 每天总伤害
	}
	// 宝石的阈值
	StoneThreshold struct {
		Level int32 // 宝石等级
		Num   int32 // 数目
	}

	// DailyMap map[int32]int32 // key 副本类型 value 挑战次数
	// DailyMap struct {
	// 	RoomType      int32 `json:"room_type"`
	// 	ChallengeTime int32 `json:"challenge_time"`
	// }

	DailyMapReward struct {
		roomType int32
		reward   []global.IdNum
	}

	AutoDisassembly struct {
		ADOpen          bool  `json:"open"`
		ADLevelSwitch   bool  `json:"level_switch"`
		ADLevel         int32 `json:"level"`
		ADQualitySwitch bool  `json:"quality_switch"`
		ADQuality       int32 `json:"quality"`
		ADHighCombat    bool  `json:"high_combat"`
	}

	PartnerSlot struct {
		PartnerUid int64
		SlotId     int32                           // 伙伴槽位ID
		Equips     [pb.EquipSlot_MaxEquipSlot]Slot // 穿戴装备列表
		TitleId    int32                           // 称号ID
	}
	BoxInfo struct {
		NowLevel            int32
		NextBoxRecoveryTime int64
		BoxRecoveryTime     time.Time // 最后一次宝箱刷新时间
		BeginAt             int64
		IsAccelerate        int32
		AccelerateNum       int32
		CostGold            int32
		TimeCost            int32
		NowOpenRes          []pb.ItemInfo
		Preview             []global.IdNumQuality // 预览装备列表
		OfflineTime         int64
		// NowOpenLevel        int32
		// Job                 int32
	}
	PetSlotInfo struct {
		PetUid   int64 `json:"petUid"`   // 宠物ID
		PetKind  int32 `json:"petKind"`  // 宠物配置ID
		PetCfgId int32 `json:"petcfgid"` // 宠物配置ID
		Level    int32 `json:"level"`    // 宠物的等级
		Star     int32 `json:"star"`     // 宠物的星级
		Enhance  int32 `json:"enhance"`  // 宠物的强化等级
		Quality  int32 `json:"quality"`  // 装备的品质
		Exp      int32 `json:"exp"`      // 宠物的经验
		ExpStar  int32 `json:"expStar"`  // 宠物的经验
		ItemId   int32 `json:"itemId"`   // 宠物item的id
	}

	InstanceData struct {
		// 战役副本
		BattleMap struct {
			BattleId    int32
			RewardStage int32
		}
		// 世界boss
		BossMap struct {
			BattleId      int32
			ChallengeTime int32
		}
		// 试炼塔
		TrialTowerMap struct {
			BattleId int32
		}
		// 遗迹
		RelicMap struct {
			BattleId      int32
			ChallengeTime int32
			Evaluation    []global.KeyValue `json:"evaluation"`
			ChapterReward []global.KeyValue `json:"chapter_reward"` // key 章节id value 星数
		}
		// 一战到底
		FightEndMap struct {
			BattleId   int32
			Evaluation []global.FightEndInfo `json:"evaluation"` // key 对手UID value 通关时间
			RefreshAt  time.Time
		}
		// 全部类型副本挑战记录
		ChallengeRecord struct {
			Records map[int32]int32 `json:"records"` // key 副本类型 value 挑战次数
		}
		ElementMap struct {
			BattleId      int32
			ChallengeTime int32 // 挑战次数
			NeedSweep     int32 // 扫荡状态 0-不能扫荡 1-需要扫荡  2-已经扫荡
			NeedBreak     int32 // 突破状态 0-不能突破 1-需要突破 2-已经突破
		}
		MonKillRecord struct {
			MonKillRecords map[int32]int32 `json:"mon_kill_records"` // key killmoncontrol表uid value 数量
		}
		PushingWall struct {
			UsedCount           int32  // 已经使用的次数
			TicketPurchaseCount int32  // 门票购买次数
			SessionUUID         string // 挑战唯一ID
			TmpUnixTime         int64  // 临时挑战时间戳
		}
	}
)

type Player struct {
	lock sync.Mutex

	session *session.Session
	logger  *log.Entry // 日志

	//------------------------------------------------------------------
	//战斗中使用的属性
	*MapObject
	curState     PlayerStatus
	miningObj    *MonsterInfo                       // 采集对象
	partnerSlots [MaxPartnerNum]*PartnerSlot        // 伙伴槽位
	partners     map[int64]*Partner                 // 伙伴列表
	skillsSlot   [pb.SkillsSlot_MaxSkillsSlot]int32 // 穿戴技能列表

	//仅刷新属性使用
	privilege     []protocol.PrivilegeDetailTCP // 特权卡
	activateBonds []int32                       // 激活羁绊
	avatar        *Avatar
	refine        *RefineV2
	*Artifact
	bag     *pb.Bag                         // 包裹
	equips  [pb.EquipSlot_MaxEquipSlot]Slot // 穿戴装备列表
	*PetMgr                                 // 宠物

	//加入房间 或者 结算 使用
	dailyMaps map[int32]int32 // key 副本类型 value 挑战次数
	*DailyMapReward
	*WorldBossInfo
	*PVPInfo
	*InstanceData
	*ElementInfo
	//------------------------------------------------------------------

	*SoloTowerInfo
	*TurnBased
	*LotteryInfo

	// 玩家数据
	head         string           // 头像地址
	ip           string           // ip地址
	sex          int              // 性别
	tokens       map[int32]int64  // 代币
	tokenDeletes map[int32]int64  // 代币消耗
	firstUse     map[string]int64 // 代币消耗
	exp          int64            // 经验池
	energy       int64            // 体力	进入副本消耗
	vip          int32            // vip等级
	vipExp       int64            // vip经验
	recharge     int64            // 充值总数
	registerAt   int64            // 注册时间

	accountName string
	lastLoginAt time.Time
	combat      int32 // 战斗力
	fateLevel   int32 // 命格等级
	fateExp     int64 // 命格经验

	tasks *TaskMgr // 任务

	autoBuyInstanceTime bool      // 副本时间自动续费
	expirationTime      time.Time // 房间内剩余时间
	roomId              int32     // 最后房间ID

	energyRecoveryTime     time.Time // 最后一次体力恢复时间
	nextEnergyRecoveryTime int       // 下一次恢复体力时间（秒)
	disconnectTime         time.Time // 断线时间

	// petNormal           int32     //普通宠物次数
	//petRecoveryTime     time.Time // 最后一次宠物蛋恢复时间
	//nextPetRecoveryTime int       // 下一次恢复宠物蛋时间（秒)
	//
	//petGold                 int32     // 普通宠物送的金币
	//petRecoveryTimeGold     time.Time // 最后一次宠物送的金币恢复时间
	//nextPetRecoveryTimeGold int       // 下一次恢复宠物送的金币时间（秒）
	//
	//petExp                 int32     // 普通宠物送的经验
	//petRecoveryTimeExp     time.Time // 最后一次宠物送的经验恢复时间
	//nextPetRecoveryTimeExp int       // 下一次恢复宠物送的经验时间（秒)

	dayOnlineTime    int64 // 当日在线时长
	totalOnlineTime  int64 // 累计在线时长
	vitalityTime     int64 // 活力满满时间
	wellenOnlineTime int64 // 本次在线时长

	nowDay    int32 // 在线第几天
	createDay int32 // 出生第几天
	// vitalityRecoveryTime int   //下一次给奖励时间（秒)

	lastEquipRefine LastEquipRefine
	bCreate         bool // 新角色第一次登录gameserver
	bDirty          bool

	*AutoDisassembly // 自动分解

	box BoxInfo // 宝箱信息
	//pet          [pb.PetsSlot_MaxPetsSlot]PetSlotInfo // 宠物列表
	//petEverHave  []int32                              // 曾经 已拥有的宠物列表
	// partnerEverHave []int32                              //曾经已拥有的伙伴列表
	totalCombat          int32     // 总战斗力
	guildId              int64     // 工会ID
	guildName            string    // 工会名称
	offlineTime          int64     // 离线时间
	guildContribution    int64     // 工会贡献
	guildWish            int64     // 工会祈愿
	wishRecoveryTime     time.Time // 最后一次工会祈愿恢复时间
	nextWishRecoveryTime int       // 下一次恢复工会祈愿时间（秒）)
	afkExpirationAt      time.Time // 自动挂机特权卡过期时间
	afkDailyTime         int64     // 每日自动挂机时
	AfkMultiplierID      int64     // 自动挂机倍率的id

	// activityMgr *ActivityMgr
	Guidance    []int32                       // 引导任务
	mailCounter int                           // 用户当前邮件数量
	firstUp     *utils.SafeMap[string, int64] // 上报次数

	mainBattleId int32 // 主界面战斗关卡ID

	warDailySoul int64 // 群雄战场每天获得的魂

	DailyItemUsed map[int32]int64 // 每日使用的物品 key 物品cfgid value 使用次数
	adTotalToday  int32           //今日看广告次数
	taskProgress  int             // 任务进度
	moduleOpen    []int32         // 已开启的模块
	ModuleNew     sync.Map        // 功能New图标激活 map[int32]*pb.ModuleInfo

	eventMgr  *event.EventManager
	muteUntil int64 // 禁言到期时间（Unix时间戳，0表示未禁言）

}

func newPlayer(s *session.Session, uid int64, name, head, ip string) *Player {
	p := &Player{
		session: s,
		// uid:     uid,
		head: head,
		ip:   ip,

		logger: log.WithField(fieldPlayer, uid),

		curState:      PlayerStatusIdle,
		bag:           &pb.Bag{ItemList: []*pb.ItemInfo{}},
		equips:        [pb.EquipSlot_MaxEquipSlot]Slot{},
		PVPInfo:       &PVPInfo{PVPDuel: &PVPDuel{}, PVPMelee: &PVPMelee{}},
		SoloTowerInfo: &SoloTowerInfo{},

		MapObject: NewMapObject(ObjectTypePlayer, uid),

		expirationTime:      time.Now(),
		autoBuyInstanceTime: false,
		privilege:           []protocol.PrivilegeDetailTCP{},
		WorldBossInfo:       &WorldBossInfo{},
		dailyMaps:           make(map[int32]int32, 8),
		DailyMapReward:      &DailyMapReward{},
		AutoDisassembly:     &AutoDisassembly{},
		partnerSlots:        [MaxPartnerNum]*PartnerSlot{},
		partners:            map[int64]*Partner{},
		tokens:              map[int32]int64{},
		tokenDeletes:        map[int32]int64{},
		InstanceData:        &InstanceData{},
		// activityMgr:         NewActivityManager(),
		TurnBased: &TurnBased{},
		//pet:           [pb.PetsSlot_MaxPetsSlot]PetSlotInfo{},
		firstUp:       utils.NewSafeMap[string, int64](),
		activateBonds: []int32{},
		ElementInfo:   &ElementInfo{},
		DailyItemUsed: make(map[int32]int64, 0),
		adTotalToday:  0,
		LotteryInfo:   &LotteryInfo{},
		avatar:        NewAvatar(),
		refine:        NewRefineV2(),
		eventMgr:      event.NewEventManager(),
		muteUntil:     0, // 禁言到期时间，0表示未禁言
	}
	p.registerEventHandlers()
	for i := 0; i < int(MaxPartnerNum); i++ {
		sid_ := i + 1
		p.partnerSlots[i] = &PartnerSlot{SlotId: int32(sid_)}
		for j := 0; j < int(pb.EquipSlot_MaxEquipSlot); j++ {
			p.partnerSlots[i].Equips[j] = Slot{}
		}
	}

	p.InstanceData.FightEndMap.Evaluation = []global.FightEndInfo{}
	beginId := int32(1900001)
	maxstep := 0
	for _, v := range global.BattleCfg {
		if v.CopyType == 19 {
			maxstep++
		}
	}
	for i := 0; i < maxstep; i++ {
		p.InstanceData.FightEndMap.Evaluation = append(p.InstanceData.FightEndMap.Evaluation, global.FightEndInfo{Uid: 0, BattleId: beginId, ChallengeTime: 0})
		beginId++
	}

	p.tasks = NewTask(p)
	p.name = name
	p.uid = uid
	p.InstanceData.RelicMap.Evaluation = []global.KeyValue{}
	p.InstanceData.ChallengeRecord.Records = make(map[int32]int32)
	p.InstanceData.MonKillRecord.MonKillRecords = make(map[int32]int32)
	p.pvpRecorder = []*model.PvpRecord{}
	p.Guidance = []int32{}
	p.bRobot = false
	p.ElementInfo.element = [pb.ElementSlot_MaxElementpSlot]ElementSlot{}
	p.ElementInfo.elementInwork = [ElementPartnerMax]int64{}
	p.ElementInfo.elementMain = int32(pb.ElementSlot_FireEmpower)
	p.LotteryInfo.HavePrizePool = []int32{}
	p.LotteryInfo.PrizePool = []int32{}
	p.LotteryInfo.RefreshTime = 0
	p.LotteryInfo.NeedItemId = 0
	p.LotteryInfo.Progress = 0
	p.LotteryInfo.FreeNum = 1
	p.Artifact = NewArtifact(p)
	p.PetMgr = NewPetMgr(p)
	return p
}

func (p *Player) PushMsg(route string, v interface{}) error {
	return p.pushMsg(route, v)
}

func (p *Player) worldBossDamage(damage int64) {
	p.totalWorldBossDamage += damage
	if p.desk != nil {
		p.desk.DWorldBoss.curRoundDamage += damage
	}
	// p.logger.Debugf(" 本次伤害 %d 世界BOSS伤害累计 %d", damage, p.totalWorldBossDamage)
}

// IMapObject interface end

// RefreshStatsOption 定义刷新属性的选项函数
type RefreshStatsOption func(*refreshStatsConfig)

// refreshStatsConfig 刷新属性的配置
type refreshStatsConfig struct {
	shouldPushMessage bool // 是否推送消息
}

// WithRefreshStatsPushMessage 设置是否推送消息
func WithRefreshStatsPushMessage(push bool) RefreshStatsOption {
	return func(cfg *refreshStatsConfig) {
		cfg.shouldPushMessage = push
	}
}

// newRefreshStatsConfig 创建默认的刷新属性配置
func newRefreshStatsConfig() *refreshStatsConfig {
	return &refreshStatsConfig{
		shouldPushMessage: true, // 默认推送消息
	}
}

// 计算角色属性
func (p *Player) refreshStats(opts ...RefreshStatsOption) {
	defer func() {
		if err := recover(); err != nil {
			p.logger.Errorf("Panic Error=%v  Stack=%v", err, stack())
		}
	}()
	if !p.slock.TryLock() {
		return
	}
	defer p.slock.Unlock()

	cfg := newRefreshStatsConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if p.isRobot() {
		p.logger.Errorf("错误调用，机器人刷新属性使用 r.refreshStats %d %s", p.uid, p.name)
		return
	}

	now := time.Now()
	p.stats.clear()
	// 基础属性
	roleCfg, err := global.GetRoleCfg(int32(p.job), p.level)
	if err != nil {
		p.logger.Errorf("获取角色等级配置失败 level %d", p.level)
		return
	}

	for _, rc := range roleCfg.BassProperty {
		p.stats.setStats(pb.Attribute(rc.Key), rc.Value)
	}

	p.stats.setStats(pb.Attribute_ATTACK_SPEED, global.GetRoleInitial(int32(p.job)).AttackSpeed)
	p.stats.setStats(pb.Attribute_MOVE_SPEED, global.GetRoleInitial(int32(p.job)).MoveSpeed)

	// 命格
	for _, gv := range global.FateCfg {
		if gv.FateLevel == p.fateLevel {
			for _, att := range gv.UpValue {
				if att.Job == int32(p.job) {
					p.stats.setStats(pb.Attribute(att.Key), att.Value)
				}
			}
			break
		}
	}

	// 魂石增加生命值基数=角色基础属性+命格
	gemHealthAdd := p.stat(pb.Attribute_MAX_HEALTH)

	// 装备 强化大师
	equipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	gemStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	petStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	newpetStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// guildStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// guildSkillStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// elementBaseStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	// elementAdvanceStatus := Stats{Attrs: make(map[pb.Attribute]int32)}

	var maxEnhance, maxStars int32
	maxEnhance = math.MaxInt32
	maxStars = math.MaxInt32
	var starCfg, enhanceCfg *global.EquipMasterConfig

	// 宠物
	// 宠物因为有些属性加给装备，所以提到前面来了
	//if len(p.pet) > 0 {
	//	if p.pet[0].PetCfgId != 0 {
	//		petCfg := global.GetPetCfg(p.pet[0].PetCfgId)
	//		if petCfg != nil {
	//
	//			for _, e2 := range petCfg.StarAttribute {
	//				if e2.Id == p.pet[0].Star {
	//					for _, e3 := range e2.Values {
	//						petStatus.setStats(pb.Attribute(e3.Key), e3.Value)
	//					}
	//					break
	//				}
	//			}
	//
	//			roleCfg, _ := global.GetRoleCfg(constant.Pet, p.pet[0].Level)
	//			if roleCfg != nil {
	//				rate := global.GetpetQualityRate(p.pet[0].Quality)
	//				if rate == 0 || petCfg.PetKind.Key == 2 {
	//					rate = 1
	//				}
	//				for _, e := range roleCfg.BassProperty {
	//					petStatus.setStats(pb.Attribute(e.Key), int32(math.Round(float64(e.Value)*rate)))
	//				}
	//			}
	//		}
	//	}
	//}
	p.PetMgr.setStats(&petStatus)

	// 计算属性
	newpetStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(petStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(1+(float64(petStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(petStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(1+(float64(petStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(petStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(1+float64(petStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	newpetStatus.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(petStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(1+(float64(petStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate))))
	newpetStatus.setStats(pb.Attribute_MAX_HEALTH, int32(float64(petStatus.getStats(pb.Attribute_MAX_HEALTH))*(1+(float64(petStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))))
	for pk, pv := range newpetStatus.Attrs {
		if pk == pb.Attribute_PHYSICAL_ATTACK || pk == pb.Attribute_MAGIC_ATTACK ||
			pk == pb.Attribute_PHYSICAL_DEFENSE || pk == pb.Attribute_MAGIC_DEFENSE ||
			pk == pb.Attribute_MAX_HEALTH {
			p.stats.setStats(pk, pv)
		}
	}

	allGems := []int32{}
	for slot, equip := range p.equips {
		if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
			continue
		}
		gemstoneNature := make([]int32, 0)
		for GemstoneSlot, Gemstone := range equip.Gemstone {
			if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
				continue
			}
			genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
			if err != nil {
				continue
			}

			if genstoneCfg.NatureMark > 0 {
				gemstoneNature = append(gemstoneNature, genstoneCfg.NatureMark)
			}
			// 对应宝石等级的数目加1
			allGems = append(allGems, genstoneCfg.Itemlevel)

			for _, attr := range genstoneCfg.Stats {
				if attr.Key == pb.Attribute_HEALTH_PERCENTAGE {
					gemStatus.setStats(attr.Key, attr.Value)
				} else {
					p.stats.setStats(attr.Key, attr.Value)
				}
			}
		}

		// 魂石自然属性
		if len(gemstoneNature) == int(pb.GemstoneSlot_MaxGemstoneSlot) && gemstoneNature[0] == gemstoneNature[1] && gemstoneNature[1] == gemstoneNature[2] && gemstoneNature[2] == gemstoneNature[3] {
			gemNatureCfg, ok := global.GetEquipMaster(constant.EquipMasterCategoryGemStoneNature*1000 + gemstoneNature[0])
			if ok != nil {
				continue
			}
			for _, v := range gemNatureCfg.Attribute {
				p.stats.setStats(pb.Attribute(v.Key), v.Value)
			}
		}

		if equip.Enhance < maxEnhance {
			maxEnhance = equip.Enhance
		}
		if equip.Stars < maxStars {
			maxStars = equip.Stars
		}
		if equip.ItemId == 0 {
			continue
		}
		equipInfo := p.getItemByUid(equip.ItemId)
		if equipInfo == nil {
			continue
		}
		equipCfg, err := global.GetItemCfg(equipInfo.CfgId)
		if err != nil {
			continue
		}

		// 基础值攻\防: （装备基础值+槽位强化基础值)*升星百分比 , 基础属性: 装备基础值*品质百分比
		for _, attr := range equipCfg.Stats {
			if attr.Key == pb.Attribute_PHYSICAL_ATTACK || attr.Key == pb.Attribute_MAGIC_ATTACK ||
				attr.Key == pb.Attribute_PHYSICAL_DEFENSE || attr.Key == pb.Attribute_MAGIC_DEFENSE ||
				attr.Key == pb.Attribute_MAX_HEALTH {
				if equipInfo.Equip != nil {
					// 如果是武器，会额外算上武器基础值
					if pb.EquipSlot(equipCfg.Equiptype) == pb.EquipSlot_Weapon {
						per := float64(1)
						// 获取某个buff的属性
						if buff := p.getBuffByType(constant.BuffPetWeaponAttackPercentageIncrease); buff != nil {
							per = float64(buff.BuffStats.getStats(pb.Attribute_WEAPON_ATTACK_PERCENTAGE)) / BaseRate
						}
						equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))*per))
					} else {
						equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))))
					}
				}
			} else {
				p.stats.setStats(attr.Key, attr.Value)
			}
		}
		// 极品属性直接计算
		if equipInfo.Equip != nil {
			for _, attr := range equipInfo.Equip.Attributes {
				p.stats.setStats(attr.Attr, attr.Num)
			}
		}

		star_uid := int32(p.job%10*100) + equip.Stars
		if equipStarCfg, err := global.GetStarsCfg(star_uid); err == nil {
			for _, star_attr := range equipStarCfg.Attribute {
				if star_attr.Slot == int32(slot) {
					equipStatus.setStats(pb.Attribute(star_attr.Attribute), star_attr.Value)
				}
			}
		}

		enhanc_uid := int32(p.job%10*1000) + equip.Enhance
		if equipEnhancCfg, err := global.GetEnhanceCfg(enhanc_uid); err == nil {
			for _, star_attr := range equipEnhancCfg.Attribute {
				if star_attr.Slot == int32(slot) {
					if equipCfg.Quality > 0 && equipCfg.Quality <= 7 {
						equipStatus.setStats(pb.Attribute(star_attr.Attribute), int32(float32(star_attr.Value)*global.GlobalCfg.Equipenhances[equipCfg.Quality-1].Value))
					}
				}
			}
		}
	}

	equipStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAX_HEALTH, int32(float64(equipStatus.getStats(pb.Attribute_MAX_HEALTH))*(float64(equipStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	for ek, ev := range equipStatus.Attrs {
		if ek == pb.Attribute_PHYSICAL_ATTACK || ek == pb.Attribute_MAGIC_ATTACK ||
			ek == pb.Attribute_PHYSICAL_DEFENSE || ek == pb.Attribute_MAGIC_DEFENSE ||
			ek == pb.Attribute_MAX_HEALTH {
			p.stats.setStats(ek, ev)
		}
	}

	// 魂石增加生命百分比
	p.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(gemHealthAdd)*(float64(gemStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))

	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryEnhance && maxEnhance >= equipmaster.Limit {
			enhanceCfg = equipmaster
		}
		if equipmaster.Category == constant.EquipMasterCategoryStars && maxStars >= equipmaster.Limit {
			starCfg = equipmaster
		}
	}

	// 魂石大师
	sort.Slice(allGems, func(i, j int) bool {
		return allGems[i] > allGems[j]
	})
	for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
		if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
			equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
			if err != nil {
				continue
			}
			if allGems[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
				for _, v := range equipMasterCfg.Attribute {
					p.stats.setStats(pb.Attribute(v.Key), v.Value)
				}
				break
			}
		}
	}

	// 强化大师
	if enhanceCfg != nil {
		for _, v := range enhanceCfg.Attribute {
			p.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	// 升星大师
	if starCfg != nil {
		for _, v := range starCfg.Attribute {
			p.stats.setStats(pb.Attribute(v.Key), v.Value)
		}
	}

	// 特权卡
	privilegeStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, v := range p.privilege {
		cfg, ok := global.PrivilegeCfg[int32(v.Pid)]
		if !ok {
			continue
		}
		for _, att := range cfg.Pcontent {
			av := privilegeStatus.getStats(pb.Attribute(att.Key))
			if att.Value > av {
				privilegeStatus.changeStats(pb.Attribute(att.Key), att.Value)
			}
		}
	}
	for k, v := range privilegeStatus.Attrs {
		p.stats.setStats(k, v)
	}

	// VIP
	vipStatus := Stats{Attrs: make(map[pb.Attribute]int32)}
	for _, gv := range global.VIPCfg {
		if gv.VIPLevel <= p.vip {
			for _, att := range gv.Pcontent {
				av := vipStatus.getStats(pb.Attribute(att.Key))
				if att.Value > av {
					vipStatus.changeStats(pb.Attribute(att.Key), att.Value)
				}
			}
		}
	}
	for k, v := range vipStatus.Attrs {
		p.stats.setStats(k, v)
	}

	p.stats.setStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE, int32(global.GlobalCfg.BaseValue.Crit.CriticalDamage))

	// 工会基础属性
	if p.guildId > 0 {
		// 工会技能
		// guildManager := NewGuildManager()
		guild, err := guildManager.validateGuild(p.guildId)
		if err == nil {
			guildCfg := global.GetGuildCfg(int32(guild.Guild.Level))
			if guildCfg != nil {
				for _, e := range guildCfg.AttributeAdd {
					p.stats.setStats(pb.Attribute(e.Key), e.Value)
				}
			}
			// 查询工会目前技能情况
			for _, e := range guild.Skills {
				skillCfg := global.GetGuildSkillCfg(e.UUid)
				if skillCfg != nil {
					for _, e2 := range skillCfg.ValueAdd {
						p.stats.setStats(pb.Attribute(e2.Key), e2.Value)
					}
				}
			}
		}

	}
	// 被动技能
	for _, v := range p.skills {
		if skillCfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			if skillCfg.ActiveorPassive == PassiveSkill {
				for _, e := range skillCfg.Effect {
					p.stats.setStats(pb.Attribute(e.Key), e.Value)
				}
			}
		}
	}
	// buff
	petBuffs := p.getPetBuffs()
	//if p.block.TryLock() {
	//logger.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxx", p.buffs)
	p.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
		buffCfg, err := global.GetBuffCfg(buff.Id)
		if err == nil {
			// buff白名单，有的buff不计算战斗力，
			// 判断是不是在白名单里面
			// 获取宠物buff
			if len(petBuffs) > 0 {
				for _, petBuff := range petBuffs {
					// 获取buffid对应的大类
					buffCfgDetail, err := global.GetBuffCfg(petBuff.Id)
					if err == nil {
						//logger.Debug("yyyyyyyyyyyyyyyyyyyyy", buffCfg.ID, petBuff.Id)
						if petBuff.Id == buffCfg.ID && (constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageAttackBoost ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageDefenseBoost ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageMaxHealthIncrease ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffIncreaseFixedMonster ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffCommonBuff ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPetWeaponAttackPercentageIncrease ||
							constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentagePKDamageIncrease) {
							for i, e := range petBuff.Stats.Attrs {
								//logger.Infof("dddddddddddddddddddddddddddd%d------------%d", i, e)
								p.stats.setStats(i, e)
							}
						}
					}

				}
			}
		}
		return true
	})

	//for _, buff := range p.buffs {
	//	buffCfg, err := global.GetBuffCfg(buff.Id)
	//	if err == nil {
	//		// buff白名单，有的buff不计算战斗力，
	//		// 判断是不是在白名单里面
	//		// 获取宠物buff
	//		if len(petBuffs) > 0 {
	//			for _, petBuff := range petBuffs {
	//				// 获取buffid对应的大类
	//				buffCfgDetail, err := global.GetBuffCfg(petBuff.Id)
	//				if err == nil {
	//					//logger.Debug("yyyyyyyyyyyyyyyyyyyyy", buffCfg.ID, petBuff.Id)
	//					if petBuff.Id == buffCfg.ID && (constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageAttackBoost ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageDefenseBoost ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentageMaxHealthIncrease ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffIncreaseFixedMonster ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffCommonBuff ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPetWeaponAttackPercentageIncrease ||
	//						constant.BuffType(buffCfgDetail.Pretype) == constant.BuffPercentagePKDamageIncrease) {
	//						for i, e := range petBuff.Stats.Attrs {
	//							//logger.Infof("dddddddddddddddddddddddddddd%d------------%d", i, e)
	//							p.stats.setStats(i, e)
	//						}
	//					}
	//				}
	//
	//			}
	//		}
	//	}
	//
	//}
	//	p.block.Unlock()
	//}

	// 羁绊
	for _, v := range p.activateBonds {
		if cfg := global.GetPartnerBondCfg(v); cfg != nil {
			for _, a := range cfg.Attribute {
				p.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}
	//五行石计算
	if p.ElementInfo != nil {
		p.stats.setStats(pb.Attribute(pb.Attribute_Main_EMPOWER), p.ElementInfo.elementMain)
		for _, v := range p.ElementInfo.element {
			if v.Eid == 0 {
				continue
			}
			if cfg := global.GetElementCfg(v.Eid); cfg != nil {
				for _, a := range cfg.Levelvalue {
					p.stats.setStats(pb.Attribute(a.Key), a.Value)
				}
				if v.Break != cfg.Breaklevel {
					if v.Eid-1 == 0 {
						continue
					}
					if cfg := global.GetElementCfg(v.Eid - 1); cfg != nil {
						for _, a := range cfg.Breakvalue {
							if a.Value > 0 {
								p.stats.setStats(pb.Attribute(a.Key), a.Value)
							}
						}
					}
				} else {
					if v.Eid == 0 {
						continue
					}
					if cfg := global.GetElementCfg(v.Eid); cfg != nil {
						for _, a := range cfg.Breakvalue {
							if a.Value > 0 {
								p.stats.setStats(pb.Attribute(a.Key), a.Value)
							}
						}
					}
				}
			}
		}
	}

	//称号共鸣
	for _, v := range p.activateTitleBonds {
		if cfg := global.GetTitleActivationCfg(v); cfg != nil {
			for _, a := range cfg.Attribute {
				p.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}
	// 称号升星属性
	for _, v := range p.titles.Titles {
		if cfg := global.GetTitleCfg(v.CfgId); cfg != nil {
			if cfg.Slot == 1 {
				for _, g := range cfg.Goods {
					if v.Star == g.Id {
						for _, a := range g.Values {
							p.stats.setStats(pb.Attribute(a.Key), a.Value)
						}
						break
					}
				}
			}
		}
	}
	//称号装备属性
	if p.equipTitleId > 0 {
		if cfg := global.GetTitleCfg(p.equipTitleId); cfg != nil {
			for _, a := range cfg.Limit {
				p.stats.setStats(pb.Attribute(a.Key), a.Value)
			}
		}
	}

	// 头像和头像框属性
	p.avatar.Stats(p)

	// 槽位精炼属性
	p.refine.Stats(p, nil, pb.CharacterSlot_CHARACTER_SLOT_1)

	// 神兵属性
	p.artifactStatus(p.MapObject)

	// 参与计算战斗力的属性写在上面!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	combatStats := Stats{Attrs: p.stats.copyStats()}
	// buff
	//if p.block.TryLock() {
	p.buffMgr.IterateBuff(func(buff *BuffInfo) bool {
		buffCfg, err := global.GetBuffCfg(buff.Id)
		if err != nil {
			return true
		}

		if buffCfg.EffectTime != 0 {
			return true
		}

		if now.After(buff.OverTime) {
			return true
		}

		// 魔法酸雨的增伤只针对技能
		if buffCfg.Pretype == constant.SkillTypeMagicAcidRain {
			return true
		}

		have := false
		if len(petBuffs) > 0 {
			for _, petBuff := range petBuffs {
				if petBuff.Id == buffCfg.ID {
					have = true
					break
				}
			}
		}
		if have {
			return true
		}
		//for i, e := range buff.BuffStats.Attrs {
		//	p.stats.setStats(i, e)
		//}
		buff.BuffStats.Range(func(key pb.Attribute, value int32) bool {
			p.stats.setStats(key, value)
			return true
		})
		return true
	})

	//for _, buff := range p.buffs {
	//	buffCfg, err := global.GetBuffCfg(buff.Id)
	//	if err != nil {
	//		continue
	//	}
	//
	//	if buffCfg.EffectTime != 0 {
	//		continue
	//	}
	//
	//	if now.After(buff.OverTime) {
	//		continue
	//	}
	//
	//	// 魔法酸雨的增伤只针对技能
	//	if buffCfg.Pretype == constant.SkillTypeMagicAcidRain {
	//		continue
	//	}
	//
	//	have := false
	//	if len(petBuffs) > 0 {
	//		for _, petBuff := range petBuffs {
	//			if petBuff.Id == buffCfg.ID {
	//				have = true
	//				break
	//			}
	//		}
	//	}
	//	if have {
	//		continue
	//	}
	//	for i, e := range buff.BuffStats.Attrs {
	//		p.stats.setStats(i, e)
	//	}
	//}
	//	p.block.Unlock()
	//}
	// <0  减益 ; >0 增益
	p.stats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(p.GetDC())*(float64(p.stats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(p.GetMC())*(float64(p.stats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(p.GetAR())*(float64(p.stats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(p.GetMR())*(float64(p.stats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	p.stats.setStats(pb.Attribute_DODGE, int32(float64(p.stats.getStats(pb.Attribute_DODGE))*(float64(p.stats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_HIT, int32(float64(p.stats.getStats(pb.Attribute_HIT))*(float64(p.stats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(p.stats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(p.stats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_TOUGHNESS, int32(float64(p.stats.getStats(pb.Attribute_TOUGHNESS))*(float64(p.stats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	p.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(p.GetMaxHP())*float64(p.stats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))
	p.stats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(p.GetMaxHP())*float64(p.stats.getStats(pb.Attribute_HEALTH_MAX_PERCENTAGE))/BaseRate))

	// 攻速移速增加 基础*(1-加成)
	p.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(p.stats.getStats(pb.Attribute_ATTACK_SPEED))*(1-float64(p.stats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate)))
	p.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(p.stats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(p.stats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	p.stats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(p.GetDC()))))
	p.stats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(p.GetMC()))))
	p.stats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(p.GetAR()))))
	p.stats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(p.GetMR()))))

	p.stats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_DODGE)))))
	p.stats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_HIT)))))
	p.stats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	p.stats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_TOUGHNESS)))))
	p.stats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_ATTACK_SPEED)))))
	p.stats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(p.stats.getStats(pb.Attribute_MOVE_SPEED)))))
	p.stats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(p.GetMaxHP()))))

	if p.curState == PlayerStatusPlaying {
		if p.currentHP > p.GetMaxHP() {
			p.currentHP = p.GetMaxHP()
		}
		p.stats.changeStats(pb.Attribute_HEALTH, p.currentHP)
		combatStats.changeStats(pb.Attribute_HEALTH, p.currentHP)
	} else {
		p.stats.changeStats(pb.Attribute_HEALTH, p.GetMaxHP())
		p.currentHP = p.GetMaxHP()
		p.displayHP = p.currentHP
		combatStats.changeStats(pb.Attribute_HEALTH, p.GetMaxHP())
	}

	//p.logger.Debugf(" 机器人【%t】  [%s] 真实战力  %d  属性 %v ", p.IsRobot(), p.name, combat, &combatStats)

	//tmpStats2 := p.stats.copyStats()

	combat := int32(p.calcCombatRealV2(&combatStats, true))
	p.combat = combat

	packet := &pb.SyncUserStatsPush{Attributes: []*pb.Stats{}, Partners: []*pb.SyncUserStatsPush_PartnerStats{}}
	p.stats.Range(func(k pb.Attribute, v int32) bool {
		if v <= 0 {
			return true
		}
		packet.Attributes = append(packet.Attributes, &pb.Stats{Attr: k, Num: v})
		return true
	})

	//for k, v := range tmpStats2 {
	//	if v <= 0 {
	//		continue
	//	}
	//	packet.Attributes = append(packet.Attributes, &pb.Stats{Attr: k, Num: v})
	//}
	packet.Combat = combat
	// 刷新伙伴属性
	for _, v := range p.partnerSlots {
		if v == nil {
			continue
		}
		partner := p.getPartnerInfo(v.PartnerUid)
		if partner == nil {
			continue
		}
		partner.refreshStats(packet, p)
	}

	packet.TotalCombat += packet.Combat
	for _, v := range packet.Partners {
		packet.TotalCombat += v.Combat
	}
	oldTotalCombat := p.totalCombat

	p.totalCombat = packet.TotalCombat
	if p.totalCombat > oldTotalCombat {
		p.tasks.dispatch(pb.TargetType_ROLE_POWER_REQUIREMENT, 0, p.totalCombat)
		if !p.bRobot {
			p.rankChangeTCP(constant.RankTypeFight, int64(p.totalCombat))
			//因为五行副本也有战斗力，所以也要刷新一下
			if p.ElementMap.BattleId > 0 && p.ElementMap.NeedBreak != 1 {
				p.rankChangeTCP(constant.RankTypeElement, int64(p.ElementMap.BattleId)*100000000000+int64(p.totalCombat))
			}
		}
		if p.desk == nil && !p.bRobot {
			robotMgr.updateRankPlayer(p)
		}
	}
	if cfg.shouldPushMessage {
		p.pushMsg(protocol.RouteSyncUserStatsPush, packet)
	}
	// p.logger.Debugf("属性计算 当前生命:%d  最大生命:%d  物攻: %d 魔攻: %d", p.currentHP, p.GetMaxHP(), p.GetDC(), p.GetMC())
}

// 异步从数据库同步
func (p *Player) syncFromDB() bool {
	ch := make(chan bool)
	async.Run(func() {
		u, err := db.QueryUser(p.Uid())
		if err != nil {
			p.logger.Errorf("玩家基础数据同步错误,userId: %d Error=%v", p.Uid(), err)
			ch <- false
			return
		}
		// ru, err := db.QueryRegisterByUid(p.Uid())
		// if err != nil {
		// 	p.logger.Errorf("玩家基础数据同步错误Register,userId: %d Error=%v", p.Uid(), err)
		// 	ch <- false
		// 	return
		// }
		p.lock.Lock()
		defer p.lock.Unlock()
		p.totalOnlineTime = u.TotalOnlineTime
		p.dayOnlineTime = u.DayOnlineTime
		p.accountName = u.AccountName
		p.name = u.Name
		p.energy = u.Energy
		p.exp = u.Exp
		p.level = u.Level
		p.vip = u.Vip
		p.recharge = u.Recharge
		p.job = constant.MowJob(u.Job)
		p.accountName = u.AccountName
		p.energyRecoveryTime = time.Unix(u.EnergyRecoveryTime, 0)
		p.lastLoginAt = time.Unix(u.LastLoginAt, 0)
		if u.LastLoginAt == 0 {
			p.bCreate = true
		}
		p.nowDay = u.Nowday
		p.registerAt = u.RegisterAt

		p.dailySoul = u.DailySoul
		p.sessionSoul = u.SessionSoul
		p.pvpScore = int32(u.PvpScore)
		p.gvgScore = int32(u.GvgScore)
		p.PVPDuel.dailyBuyNum = int32(u.BuyMatchNumDuel)
		p.PVPDuel.dailyMatchTimes = int32(u.MatchTimesDuel)
		p.PVPMelee.dailyBuyNum = int32(u.BuyMatchNumMelee)
		p.PVPMelee.dailyMatchTimes = int32(u.MatchTimesMelee)

		p.SoloTowerInfo.Layers = int32(u.SoloTowerLayers)
		p.SoloTowerInfo.RewardTime = time.Unix(u.SoloTowerRewardTime, 0)
		p.sex = u.Sex
		p.fateExp = u.FateExp
		p.fateLevel = u.FateLevel
		p.vitalityTime = u.VitalityTime
		p.dailyWorldBossTimes = int32(u.WorldBossTimes)
		p.totalWorldBossDamage = u.WorldBossDamage
		p.serverID = u.ServerId

		// 新的同步总战斗力
		p.totalCombat = u.TotalCombat
		p.afkExpirationAt = time.Unix(u.AfkExpirationAt, 0)
		p.afkDailyTime = u.AfkDailyTime
		p.AfkMultiplierID = u.AfkMultiplierID
		p.offlineTime = u.OfflineTime
		// p.AdTotalToday = u.AdTotalToday
		atomic.StoreInt32(&p.adTotalToday, u.AdTotalToday)
		p.muteUntil = u.MuteUntil
		// p.channelId = u.ChannelId
		// 体力恢复
		// p.nextEnergyRecoveryTime = 0
		// if p.energy < int64(global.GlobalCfg.BaseValue.Energy.Fullvalue) {
		// 	// 离线恢复总体力
		// 	recoverNum := (time.Now().Unix() - u.EnergyRecoveryTime) / int64(global.GlobalCfg.BaseValue.Energy.Timeinterval*60)
		// 	// 不足恢复一次体力的时间
		// 	leftTime := (time.Now().Unix() - u.EnergyRecoveryTime) % int64(global.GlobalCfg.BaseValue.Energy.Timeinterval*60)
		// 	p.energy += recoverNum
		// 	if p.energy > int64(global.GlobalCfg.BaseValue.Energy.Fullvalue) {
		// 		p.energy = int64(global.GlobalCfg.BaseValue.Energy.Fullvalue)
		// 		p.energyRecoveryTime = time.Now()
		// 	} else {
		// 		p.nextEnergyRecoveryTime = global.GlobalCfg.BaseValue.Energy.Timeinterval*60 - int(leftTime)
		// 		p.energyRecoveryTime = time.Now().Add(-time.Duration(global.GlobalCfg.BaseValue.Energy.Timeinterval) * time.Minute)
		// 	}
		// }

		if err = protobuf.NewSerializer().Unmarshal(u.Bag, p.bag); err != nil {
			p.logger.Errorf("包裹数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}

		if err = json.Unmarshal([]byte(u.Skill), &p.skills); u.Skill != "" && err != nil {
			p.logger.Errorf("技能数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}

		if err = json.Unmarshal([]byte(u.Equips), &p.equips); u.Equips != "" && err != nil {
			p.logger.Errorf("装备数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}
		if err = json.Unmarshal([]byte(u.SkillSlot), &p.skillsSlot); u.SkillSlot != "" && err != nil {
			p.logger.Errorf("技能框数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}
		dbtask := &pb.DBTask{}
		if err = protobuf.NewSerializer().Unmarshal(u.Task, dbtask); err == nil {
			p.tasks.fromPB(dbtask)
			for _, v := range p.tasks.tasks {
				if v.Class == pb.TaskClass_MAINLINE {
					if v.Status == pb.TaskStatus_REWARDED {
						if p.taskProgress < v.ID {
							p.taskProgress = v.ID
						}
						continue
					}
					//查主线任务对应的信息
					cfg, err := global.GetTaskCfg(int32(v.ID))
					if err != nil {
						continue
					}
					if int32(v.Target.TargetType) != cfg.Type {
						v.Target.TargetType = pb.TargetType(cfg.Type)
						v.Target.TargetId = int(cfg.Typeson)
						v.Target.TargetNum = int(cfg.Limit)
						v.Target.Amount = 0
						v.Target.Completed = false

						p.logger.Debugf("任务配置修改,任务目标更新 :%v ", v.Target)
					}
				}
			}
		} else {
			p.logger.Errorf("包裹数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}

		// if err = json.Unmarshal([]byte(u.DailyMap), &p.dailyMaps); u.DailyMap != "" && err != nil {
		// 	p.logger.Errorf("每日必做数据反序列化失败 err: %s", err.Error())
		// 	ch <- false
		// 	return
		// }
		if u.DailyMap != "" {
			if err = json.Unmarshal([]byte(u.DailyMap), &p.dailyMaps); err != nil {
				p.logger.Infof("每日必做数据反序列化失败 err: %s, 已初始化为默认值", err.Error())
				p.dailyMapReset()
			}
		}

		if u.DailyItemUsed != "" {
			if err = json.Unmarshal([]byte(u.DailyItemUsed), &p.DailyItemUsed); err != nil {
				p.logger.Infof("每日道具使用限制反序列化失败 err: %s", err.Error())
			}
		}

		if err = json.Unmarshal([]byte(u.TurnBasedSlot), &p.turnBasedSlot); u.TurnBasedSlot != "" && err != nil {
			p.logger.Errorf("回合制数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}

		if err = json.Unmarshal([]byte(u.Guidance), &p.Guidance); u.Guidance != "" && err != nil {
			p.logger.Errorf("引导数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}
		if u.ModuleOpen != "" {
			if err = json.Unmarshal([]byte(u.ModuleOpen), &p.moduleOpen); err != nil {
				p.logger.Infof("模块打开数据反序列化失败 err: %s", err.Error())
			}
		}

		// moduleNew
		p.fromDbModuleNew(u.ModuleNew)

		// buff
		p.fromDbBuff(u.Buff)

		// 加载特权卡
		if privilege, err := provide.PrivilegeListTCP(p.uid); err == nil {
			for _, v := range privilege {
				p.addPrivilege(constant.ShopTypeSubscription, v.Pid, v.EndTime)
			}
		}

		if err = json.Unmarshal([]byte(u.Disassembly), &p.AutoDisassembly); u.Disassembly != "" && err != nil {
			p.logger.Errorf("自动分解数据反序列化失败 err: %s", err.Error())
			ch <- false
			return
		}

		// 代币
		if err = json.Unmarshal([]byte(u.Tokens), &p.tokens); u.Tokens != "" && err != nil {
			p.logger.Errorf("代币数据反序列化 err: %s", err.Error())
			ch <- false
			return
		}
		if u.FirstUp != "" {
			// 兼容老数据：先反序列化到普通 map，再转换到 SafeMap
			var tempMap map[string]int64
			if err = json.Unmarshal([]byte(u.FirstUp), &tempMap); err != nil {
				p.logger.Errorf("是否第一次数据反序列化 err: %s", err.Error())
				ch <- false
				return
			}
			// 将老数据迁移到 SafeMap
			for k, v := range tempMap {
				p.firstUp.Set(k, v)
			}
		}
		// 代币消耗
		if err = json.Unmarshal([]byte(u.TokenDeletes), &p.tokenDeletes); u.TokenDeletes != "" && err != nil {
			p.logger.Errorf("代币消耗数据反序列化 err: %s", err.Error())
			ch <- false
			return
		}
		if err = json.Unmarshal([]byte(u.FirstUse), &p.firstUse); u.FirstUse != "" && err != nil {
			p.logger.Errorf("代币消耗数据反序列化 err: %s", err.Error())
			ch <- false
			return
		}
		if err = p.syncPartnerFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.syncBoxFromDB(); err != nil {
			p.logger.Errorf("宝箱同步失败 err: %s", err.Error())
			ch <- false
			return
		}
		if err = p.syncPetFromDB(); err != nil {
			p.logger.Errorf("宠物信息同步失败 err: %s", err.Error())
			ch <- false
			return
		}
		// 工会信息
		if err = p.syncGuildFromDB(); err != nil {
			p.logger.Errorf("工会信息同步失败 err: %s", err.Error())
			ch <- false
			return
		}
		if err = p.syncInstanceDataFromDB(); err != nil {
			ch <- false
			return
		}

		if err = p.syncPvpRecorderFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.syncAttributeFromDB(); err != nil {
			p.logger.Errorf("五行副本同步失败 err: %s", err.Error())
			ch <- false
			return
		}
		if err = p.syncWarFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.syncLotteryFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.syncAvatarFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.syncRefineFromDB(); err != nil {
			ch <- false
			return
		}
		if err = p.loadArtifactDB(); err != nil {
			ch <- false
			return
		}
		ch <- true
	})
	return <-ch
}

// 定义包级缓存和时间戳
// var (
// 	userBoxCache      = sync.Map{} // key: uid(int64), value: *model.UserBox
// 	userBoxCacheTimes = sync.Map{} // key: uid(int64), value: int64(Unix秒)
// )

func (p *Player) syncBoxFromDB() error {
	// var ubox *model.UserBox
	// var err error
	uid := p.Uid()
	now := time.Now().Unix()
	// cacheDuration := int64(10) // 缓存10秒

	// // 检查缓存时间
	// if val, ok := userBoxCache.Load(uid); ok {
	// 	if ts, ok2 := userBoxCacheTimes.Load(uid); ok2 {
	// 		if now-int64(ts.(int64)) < cacheDuration {
	// 			ubox = val.(*model.UserBox)
	// 		}
	// 	}
	// }
	// if ubox == nil {
	// 	ubox, err = db.QueryUserBox(uid)
	// 	if err == nil && ubox != nil {
	// 		userBoxCache.Store(uid, ubox)
	// 		userBoxCacheTimes.Store(uid, now)
	// 	}
	// }
	ubox, err := db.QueryUserBox(uid)
	if err != nil {
		p.logger.Errorf("玩家宝箱同步错误,userId: %d Error=%v", uid, err)
		return err
	}
	timeCost := ubox.TimeCost
	// 当前时间减最后在线时间
	if p.offlineTime > 0 {
		if ubox.NextBoxRecoveryTime != -1 {
			delta := now - p.offlineTime
			if delta > ubox.NextBoxRecoveryTime {
				ubox.NextBoxRecoveryTime = 1
				timeCost += int32(ubox.NextBoxRecoveryTime)
			} else {
				ubox.NextBoxRecoveryTime = ubox.NextBoxRecoveryTime - delta
				timeCost += int32(ubox.NextBoxRecoveryTime)
			}
		}
	}

	p.box = BoxInfo{
		NowLevel:            ubox.NowLevel,
		NextBoxRecoveryTime: ubox.NextBoxRecoveryTime,
		NowOpenRes:          []pb.ItemInfo{},
		BoxRecoveryTime:     time.Unix(ubox.BeginAt, 0),
		BeginAt:             ubox.BeginAt,
		IsAccelerate:        ubox.IsAccelerate,
		AccelerateNum:       ubox.AccelerateNum,
		CostGold:            ubox.CostGold,
		TimeCost:            timeCost,
		OfflineTime:         ubox.OfflineTime,
		// NowOpenLevel:        ubox.NowOpenLevel,
		// Job:                 ubox.Job,
	}
	if ubox.NowOpenRes != "" && ubox.NowOpenRes != "0" {
		if err = json.Unmarshal([]byte(ubox.NowOpenRes), &p.box.NowOpenRes); err != nil {
			p.logger.Errorf("首页装备反序列化失败 err: %s", err.Error())
			return err
		}
	}
	return nil
}

func (p *Player) syncPartnerFromDB() error {
	list, err := db.QueryPartner(p.Uid())
	if err != nil {
		p.logger.Errorf("加载伙伴错误,userId: %d Error=%v", p.Uid(), err)

		return err
	}

	for _, v := range list {
		partner := &Partner{
			MapObject: NewMapObject(ObjectTypePartner, v.PartnerUid),
		}

		partner.cfgId = v.CfgId
		partner.star = v.Star
		partner.level = v.Level

		if err = json.Unmarshal([]byte(v.Skill), &partner.skills); v.Skill != "" && err != nil {
			p.logger.Errorf("技能数据反序列化 err: %s", err.Error())

			return err
		}
		partner.owner = p
		partner.cfg = global.GetPartnerCfg(partner.cfgId)
		partner.name = partner.cfg.Name
		partner.job = constant.MowJob(partner.cfg.Job)
		partner.attrRatio = BaseRate
		partner.Talent = &Talent{activated: make([]global.IdNum, 0)}
		partner.points = v.TalnetPoints
		if err = json.Unmarshal([]byte(v.TalnetActivated), &partner.activated); v.TalnetActivated != "" && err != nil {
			p.logger.Errorf("伙伴同步错误 天赋反序列失败,userId: %d Error=%v", p.Uid(), err)
		}

		partner.initTalent()
		p.partners[v.PartnerUid] = partner
	}

	slots, err := db.QueryPartnerSlot(p.Uid())
	if err != nil {
		p.logger.Errorf("伙伴同步错误,userId: %d Error=%v", p.Uid(), err)
		return err
	}

	for i, v := range slots {
		if i < int(MaxPartnerNum) {
			pslot := &PartnerSlot{SlotId: v.SlotId, PartnerUid: v.PartnerUid, TitleId: v.EquipTitleid}
			// p.partnerSlots[v.SlotId] = &PartnerSlot{sid: v.SlotId, puid: v.PartnerUid}
			if err = json.Unmarshal([]byte(v.Equips), &pslot.Equips); v.Equips != "" && err != nil {
				p.logger.Errorf("伙伴装备数据反序列化失败 err: %s", err.Error())
				return err
			}
			p.partnerSlots[i] = pslot
		} else {
			p.logger.Errorf("伙伴槽位数量异常 %d ", len(slots))
		}
	}

	return nil
}

// 加载war
func (p *Player) syncWarFromDB() error {
	u, err := db.QueryWar(p.uid)
	if err != nil {
		p.logger.Errorf("初始化群雄战争DB,userId: %d ", p.uid)
		go db.InsertWar(&model.War{Uid: p.uid})
	} else {
		p.equipTitleId = u.EquipTitleid
		p.warDailySoul = u.Soul
		title := &pb.DBTitle{}
		if err = protobuf.NewSerializer().Unmarshal([]byte(u.Title), title); err == nil {
			p.titles = title
		} else {
			p.logger.Errorf("包裹数据反序列化失败 err: %s", err.Error())
			return err
		}
	}

	return nil
}

// 加载抽奖
func (p *Player) syncLotteryFromDB() error {
	l, err := db.QueryLottery(p.uid)
	if err != nil {
		p.logger.Infof("初始化抽奖DB,userId: %d ", p.uid)
		if err == errutil.ErrUserNotFound {
			go db.InsertLottery(&model.Lottery{Uid: p.uid, FreeNum: 0, RefreshTime: 0, Progress: 0, NeedItemId: 0, HavePrizePool: "[]", PrizePool: "[]"})
		} else {
			return err
		}
	} else {
		p.LotteryInfo.FreeNum = l.FreeNum
		p.LotteryInfo.RefreshTime = int64(l.RefreshTime)
		p.LotteryInfo.Progress = l.Progress
		p.LotteryInfo.NeedItemId = l.NeedItemId
		havePrizePool := []int32{}
		if err = json.Unmarshal([]byte(l.HavePrizePool), &havePrizePool); err == nil {
			p.LotteryInfo.HavePrizePool = havePrizePool
		} else {
			p.logger.Errorf("抽奖反序列化失败 err: %s", err.Error())
			return err
		}
		prizePool := []int32{}
		if err = json.Unmarshal([]byte(l.PrizePool), &prizePool); err == nil {
			p.LotteryInfo.PrizePool = prizePool
		} else {
			p.logger.Errorf("抽奖反序列化失败 err: %s", err.Error())
			return err
		}
	}
	return nil
}

// 加载头像头像框
func (p *Player) syncAvatarFromDB() error {
	avatardb, err := db.GetAvatarByUid(p.uid)
	if err != nil {
		p.logger.Errorf("加载头像错误,userId: %d Error=%v", p.Uid(), err)
		return err
	}

	if avatardb == nil {
		p.logger.Infof("初始化头像DB,userId: %d ", p.uid)
		async.Run(func() {
			p.avatar.AddInitialAvatarAndFrame()
			if err = db.InsertAvatar(p.avatar.ToDB(p)); err != nil {
				p.logger.Errorf("初始化头像DB失败,userId: %d Error=%v", p.Uid(), err)
			}
		})

		return nil
	}

	p.avatar.LoadFromDB(p, avatardb)
	p.avatar.AddInitialAvatarAndFrame()

	return nil
}

// syncRefineFromDB 加载洗练
func (p *Player) syncRefineFromDB() error {
	r, err := db.GetRefineByUid(p.uid)
	if err != nil {
		p.logger.Errorf("加载洗练错误,userId: %d Error=%v", p.Uid(), err)
	}
	if r == nil {
		p.logger.Infof("初始化洗练DB,userId: %d ", p.uid)
		async.Run(func() {
			// db.InsertRefine(&model.Refine{Uid: p.uid, PurchasedSlots: []byte("{}"), Attrs: []byte("{}")})
			purchasedSlots := make([]int32, 0, p.refine.purchasedSlots.Len())
			p.refine.purchasedSlots.Range(func(key int32, _ struct{}) bool {
				purchasedSlots = append(purchasedSlots, key)
				return true
			})

			attrs := make(map[int32]*pb.Stats, p.refine.attrs.Len())
			p.refine.attrs.Range(func(key int32, value *pb.Stats) bool {
				attrs[key] = value
				return true
			})

			db.InsertRefine(&model.Refine{
				UID:            p.uid,
				PurchasedSlots: purchasedSlots,
				Attrs:          attrs,
			})
		})
		return nil
	}
	return p.refine.FromDB(r)
}

// // 定义包级缓存和时间戳
// var (
// 	userPetCache      = sync.Map{} // key: uid(int64), value: *model.UserPet
// 	userPetCacheTimes = sync.Map{} // key: uid(int64), value: int64(Unix秒)
// )

func (p *Player) syncPetFromDB() error {
	return p.PetMgr.syncPetFromDB()
	//	// var upet *model.UserPet
	//	// var err error
	//	uid := p.Uid()
	//	// now := time.Now().Unix()
	//	// cacheDuration := int64(10) // 缓存10秒
	//
	//	// // 检查缓存时间
	//	// if val, ok := userPetCache.Load(uid); ok {
	//	// 	if ts, ok2 := userPetCacheTimes.Load(uid); ok2 {
	//	// 		if now-int64(ts.(int64)) < cacheDuration {
	//	// 			upet = val.(*model.UserPet)
	//	// 		}
	//	// 	}
	//	// }
	//	// if upet == nil {
	//	// 	upet, err = db.QueryUserPet(uid)
	//	// 	if err == nil && upet != nil {
	//	// 		userPetCache.Store(uid, upet)
	//	// 		userPetCacheTimes.Store(uid, now)
	//	// 	}
	//	// }
	//	upet, err := db.QueryUserPet(uid)
	//	if err != nil {
	//		p.logger.Errorf("玩家宠物同步错误,userId: %d Error=%v", uid, err)
	//		return err
	//	}
	//	if upet == nil {
	//		return nil
	//	}
	//	if upet.PetInfo != "" && upet.PetInfo != "0" {
	//		if err = json.Unmarshal([]byte(upet.PetInfo), &p.pet); err != nil {
	//			p.logger.Errorf("宠物列表反序列化失败 err: %s", err.Error())
	//			return err
	//		}
	//	}
	//	p.petEverHave = []int32{}
	//	if upet.PetEverHave != "" && upet.PetEverHave != "0" {
	//		if err = json.Unmarshal([]byte(upet.PetEverHave), &p.petEverHave); err != nil {
	//			p.logger.Errorf("宠物曾经拥有反序列化失败 err: %s", err.Error())
	//			return err
	//		}
	//	}
	//	// ...existing code...
	//	// 宠物恢复
	//	p.nextPetRecoveryTime = 0
	//	petNormalCount := p.getItemCountInBag(constant.ItemTypePetNormal)
	//	if petNormalCount >= int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) {
	//		p.petRecoveryTime = time.Now()
	//	} else {
	//		if upet.PetRecoveryTime == 0 {
	//			// 当前时间的前一秒
	//			p.petRecoveryTime = time.Now().Add(-time.Second)
	//			p.nextPetRecoveryTime = int(global.GlobalCfg.BaseValue.Petfree.Recovertime * 60)
	//		}
	//		if upet.PetRecoveryTime > 0 {
	//			// 当前时间减去最后更新时间，就是需要恢复的次数
	//			recoverNum := (time.Now().Unix() - upet.PetRecoveryTime) / int64(global.GlobalCfg.BaseValue.Petfree.Recovertime*60)
	//			leftTime := (time.Now().Unix() - upet.PetRecoveryTime) % int64(global.GlobalCfg.BaseValue.Petfree.Recovertime*60)
	//			// 次数大于最大的次数，取最大次数
	//			if recoverNum+petNormalCount > int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) {
	//				recoverNum = int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) - petNormalCount
	//			}
	//			if recoverNum < 0 {
	//				recoverNum = 0
	//			}
	//			// 下次恢复时间
	//			p.nextPetRecoveryTime = int(global.GlobalCfg.BaseValue.Petfree.Recovertime*60) - int(leftTime)
	//			// 次数大于0才进行恢复
	//			if recoverNum > 0 {
	//				go p.addSource(constant.ItemTypePetNormal, int64(recoverNum), constant.ItemAddPetAFK)
	//				// 最后一次宠物蛋恢复时间
	//				p.petRecoveryTime = time.Now().Add(-time.Duration(leftTime) * time.Second)
	//			} else {
	//				p.petRecoveryTime = time.Unix(upet.PetRecoveryTime, 0)
	//			}
	//		}
	//	}
	//	// ...existing code...
	//	if p.pet[0].PetCfgId != 0 {
	//		petInfo := global.GetPetCfg(p.pet[0].PetCfgId)
	//		if petInfo != nil && len(petInfo.OfflineAward) > 0 {
	//			petGoldNum := int32(0)
	//			petGoldGap := int32(0)
	//			petGoldMax := int32(0)
	//			petExpNum := int32(0)
	//			petExpGap := int32(0)
	//			petExpMax := int32(0)
	//			for _, v := range petInfo.OfflineAward {
	//				if v.Key == p.pet[0].Star {
	//					for _, v2 := range v.Values {
	//						if v2.ID == constant.ItemTypeCoin {
	//							petGoldNum = v2.Num
	//							petGoldGap = int32(v2.Gap)
	//							petGoldMax = int32(v2.Max)
	//						}
	//						if v2.ID == constant.ItemTypeExp {
	//							petExpNum = v2.Num
	//							petExpGap = int32(v2.Gap)
	//							petExpMax = int32(v2.Max)
	//						}
	//					}
	//				}
	//			}
	//			if petGoldNum > 0 {
	//				if p.petGold < petGoldMax {
	//					if upet.PetRecoveryTimeGold > 0 {
	//						p.nextPetRecoveryTimeGold = 0
	//						recoverNumGold := (time.Now().Unix() - upet.PetRecoveryTimeGold) / int64(petGoldGap) * int64(petGoldNum)
	//						leftTime := (time.Now().Unix() - upet.PetRecoveryTimeGold) % int64(petGoldGap)
	//						if p.petGold+int32(recoverNumGold) > petGoldMax {
	//							p.petGold = petGoldMax
	//							p.petRecoveryTimeGold = time.Now()
	//							// upet.PetRecoveryTimeGold = time.Now().Unix()
	//							if petGoldMax-p.petGold > 0 {
	//								go p.addSource(constant.ItemTypeCoin, int64(petGoldMax-p.petGold), constant.ItemAddPetGold)
	//							}
	//						} else {
	//							p.petGold += int32(recoverNumGold)
	//							p.nextPetRecoveryTimeGold = int(petGoldGap) - int(leftTime)
	//							p.petRecoveryTimeGold = time.Now().Add(-time.Duration(petGoldGap))
	//							// upet.PetRecoveryTimeGold = time.Now().Unix() - int64(petGoldGap)
	//							if recoverNumGold > 0 {
	//								go p.addSource(constant.ItemTypeCoin, int64(recoverNumGold), constant.ItemAddPetGold)
	//							}
	//						}
	//					}
	//				}
	//			}
	//			if petExpNum > 0 {
	//				if p.petExp < petExpMax {
	//					if upet.PetRecoveryTimeExp > 0 {
	//						recoverNumExp := (time.Now().Unix() - upet.PetRecoveryTimeExp) / int64(petExpGap) * int64(petExpNum)
	//						leftTime := (time.Now().Unix() - upet.PetRecoveryTimeExp) % int64(petExpGap)
	//						if p.petExp+int32(recoverNumExp) > petExpMax {
	//							p.petExp = petExpMax
	//							p.petRecoveryTimeExp = time.Now()
	//							// upet.PetRecoveryTimeExp = time.Now().Unix()
	//							if petExpMax-p.petExp > 0 {
	//								go p.addSource(constant.ItemTypeExp, int64(petExpMax-p.petExp), constant.ItemAddPetExp)
	//							}
	//						} else {
	//							p.petExp += int32(recoverNumExp)
	//							p.nextPetRecoveryTimeExp = int(petExpGap) - int(leftTime)
	//							p.petRecoveryTimeExp = time.Now().Add(-time.Duration(petExpGap))
	//							// upet.PetRecoveryTimeExp = time.Now().Unix() - int64(petExpGap)
	//							if recoverNumExp > 0 {
	//								go p.addSource(constant.ItemTypeExp, int64(recoverNumExp), constant.ItemAddPetExp)
	//							}
	//						}
	//					}
	//				}
	//			}
	//
	//		}
	//	}
	//	return nil
}

func (p *Player) syncAttributeFromDB() error {
	// 查询玩家属性
	attribute, err := db.QueryUserAttribute(p.Uid())
	if err != nil {

		if err == errutil.ErrUserNotFound {
			//插入数据，初始化
			// 将 ElementSlot 数组转换为 pb.ElementEmpower 数组

			elements, err := protobuf.NewSerializer().Marshal(p.ElementInfoElementToPb())
			if err != nil {
				p.logger.Errorf("序列化元素数据失败 err: %s", err.Error())
				return err
			}

			uids, err := protobuf.NewSerializer().Marshal(p.ElementInfoElementInworkToPb())
			if err != nil {
				p.logger.Errorf("序列化元素数据失败 err: %s", err.Error())
				return err
			}
			attribute = &model.UserAttribute{
				Uid:           p.Uid(),
				Elements:      elements,
				ElementInWork: uids,
				ElementMain:   int32(pb.ElementSlot_FireEmpower),
			}
			if err = db.InsertUserAttribute(attribute); err != nil {
				p.logger.Errorf("插入玩家属性失败 err: %s", err.Error())
				return err
			}
		} else {
			return err
		}

	} else {
		empower := pb.ElementEmpowerList{}
		if err = protobuf.NewSerializer().Unmarshal([]byte(attribute.Elements), &empower); err != nil {
			p.logger.Errorf("元素数据反序列化 err: %s", err.Error())
			return err
		}

		// 将 pb.ElementEmpower 数组转换为 ElementSlot 数组
		for k, v := range empower.Elements {
			// if v.SlotId >= 1 && v.SlotId < int32(pb.ElementSlot_MaxElementpSlot) {
			// 根据 Level 推断 Eid
			var eid int32
			if v.Level > 0 {
				// 如果等级大于0，说明已经升级过，需要计算 Eid
				// 根据 elementUpgrade 函数的逻辑，Eid 的计算公式是 int32(slotId)*1000 + level
				eid = int32(k)*1000 + v.Level
			} else {
				eid = 0
			}

			p.ElementInfo.element[k] = ElementSlot{
				Eid:    eid,
				SlotId: int32(k),
				Level:  v.Level,
				Break:  v.BreakLevel,
			}
			// }
		}

		uid := pb.ElementUids{}
		if err = protobuf.NewSerializer().Unmarshal([]byte(attribute.ElementInWork), &uid); err != nil {
			p.logger.Errorf("元素数据反序列化 err: %s", err.Error())
			return err
		}
		// p.elementInwork = make([]int32, len(uid.Uids))
		for k, v := range uid.Uids {
			if v > 0 {
				p.ElementInfo.elementInwork[k] = v
			}
		}
		p.ElementInfo.elementMain = attribute.ElementMain
	}
	return nil
}

// 同步工会信息
func (p *Player) syncGuildFromDB() error {
	uid := p.Uid()

	guildMember, err := db.QueryGuildMember(uid, 1)
	if err != nil {
		p.logger.Errorf("工会信息同步失败 err: %s", err.Error())
		return err
	}
	if guildMember == nil {
		p.guildId = 0
		p.guildName = ""
	} else {
		p.guildId = guildMember.GuildID
		p.guildContribution = int64(guildMember.GuildContribution)
		p.guildWish = int64(guildMember.GuildWish)
		//查询工会信息
		guild, err2 := guildManager.validateGuild(guildMember.GuildID)
		if err2 == nil {
			p.guildName = guild.Guild.Name
		} else {
			p.guildName = ""
		}
	}
	if p.guildId > 0 {
		p.nextWishRecoveryTime = 0
		selectTakeSource := p.getItemCountInBag(constant.ItemTypeGuildWish)
		if selectTakeSource > int64(global.GlobalCfg.Guild.WishValue.MaxValue) {
			p.wishRecoveryTime = time.Now()
		} else {
			if guildMember.RecoveryTimeExp == 0 {
				p.wishRecoveryTime = time.Now().Add(-time.Second)
				p.nextWishRecoveryTime = int(global.GlobalCfg.Guild.WishValue.RefreshTime * 60)
			}
			// Guild wish recovery
			if guildMember.RecoveryTimeExp > 0 {
				// 离线恢复总祈愿点
				recoverNum := (time.Now().Unix() - guildMember.RecoveryTimeExp) / int64(global.GlobalCfg.Guild.WishValue.RefreshTime*60)
				// 不足恢复一次祈愿点的时间
				leftTime := (time.Now().Unix() - guildMember.RecoveryTimeExp) % int64(global.GlobalCfg.Guild.WishValue.RefreshTime*60)
				// 祈愿点恢复
				if recoverNum+selectTakeSource > int64(global.GlobalCfg.Guild.WishValue.MaxValue) {
					recoverNum = int64(global.GlobalCfg.Guild.WishValue.MaxValue) - selectTakeSource
				}
				p.nextWishRecoveryTime = int(global.GlobalCfg.Guild.WishValue.RefreshTime*60) - int(leftTime)
				if recoverNum > 0 {
					go p.addSource(constant.ItemTypeGuildWish, int64(recoverNum), constant.ItemAddGuildWish)
					p.wishRecoveryTime = time.Now().Add(-time.Duration(leftTime) * time.Second)
				} else {
					p.wishRecoveryTime = time.Unix(guildMember.RecoveryTimeExp, 0)
				}

			}
		}
	}
	return nil
}

func (p *Player) syncInstanceDataFromDB() error {
	u, err := db.QueryInstanceData(p.Uid())
	if err == nil {
		p.BattleMap.BattleId = u.BattleCurStage
		p.BattleMap.RewardStage = u.BattleRewardStage
		p.BossMap.BattleId = u.BossCurStage
		p.BossMap.ChallengeTime = u.BossChallengeTime
		p.TrialTowerMap.BattleId = u.TrialTowerCurStage

		p.RelicMap.BattleId = u.RelicCurStage
		p.RelicMap.ChallengeTime = u.RelicChallengeTime

		p.FightEndMap.BattleId = u.FightEndCurStage
		p.FightEndMap.RefreshAt = time.Unix(u.FightEndRefreshTime, 0)
		p.mainBattleId = u.MainBattleId

		p.ElementMap.BattleId = u.ElementCurStage
		p.ElementMap.ChallengeTime = u.ElementChallengeTime
		p.ElementMap.NeedSweep = u.ElementNeedSweep
		p.PushingWall.UsedCount = u.PushingWallCount
		p.PushingWall.TicketPurchaseCount = u.PushingWallTicketPurchaseCount

		if err = json.Unmarshal([]byte(u.RelicEvaluation), &p.RelicMap.Evaluation); u.RelicEvaluation != "" && err != nil {
			p.logger.Errorf("RelicMap副本评价数据反序列化失败 err: %s", err.Error())

			return errutil.ErrDBOperation
		}
		if err = json.Unmarshal([]byte(u.RelicChapterReward), &p.RelicMap.ChapterReward); u.RelicChapterReward != "" && err != nil {
			p.logger.Errorf("RelicMap副本奖励数据反序列化失败 err: %s", err.Error())

			return errutil.ErrDBOperation
		}

		var fightEndInfo []global.FightEndInfo
		if err = json.Unmarshal([]byte(u.FightEndEvaluation), &fightEndInfo); u.FightEndEvaluation != "" && err != nil {
			p.logger.Errorf("BattleToTheEndMap副本数据反序列化失败 err: %s", err.Error())
			return errutil.ErrDBOperation
		} else {
			for _, v := range fightEndInfo {
				for i := 0; i < len(p.FightEndMap.Evaluation); i++ {
					if p.FightEndMap.Evaluation[i].BattleId == v.BattleId {
						p.FightEndMap.Evaluation[i].ChallengeTime = v.ChallengeTime
						break
					}
				}
			}
		}

		if err = json.Unmarshal([]byte(u.ChallengeRecord), &p.ChallengeRecord.Records); u.ChallengeRecord != "" && err != nil {
			p.logger.Errorf("ChallengeRecord副本数据反序列化失败 err: %s", err.Error())

			return errutil.ErrDBOperation
		}
		if err = json.Unmarshal([]byte(u.KillMonControl), &p.MonKillRecord.MonKillRecords); u.KillMonControl != "" && err != nil {
			p.logger.Errorf("KillMonControl副本数据反序列化失败 err: %s", err.Error())

			return errutil.ErrDBOperation
		}
	}

	return nil
}

// 同步pvprecoder
func (p *Player) syncPvpRecorderFromDB() error {
	u, err := db.QueryPvpRecords(p.Uid())
	if err == nil {
		p.pvpRecorder = u
	}
	return nil
}

// 同步角色基础信息
func (p *Player) syncUserBase() {
	vitality := int64(0)
	// 查询vip对应的活力分数
	if p.vip > 0 {
		if vipCfg := global.GetVipCfg(p.vip); vipCfg != nil {
			vitality += int64(vipCfg.Energyadd) * 60
		}
	}
	// 查询特权卡对应的活力分数
	if len(p.privilege) > 0 {
		for _, v := range p.privilege {
			// 查询特权卡对应的活力分数
			if privilegeCfg := global.GetPrivilegeCfg(v.Pid); privilegeCfg != nil {
				vitality += int64(privilegeCfg.Energystate) * 60
			}
		}
	}

	packet := &pb.SyncUserBasePush{
		Uid:                 p.Uid(),
		Coin:                p.tokens[constant.ItemTypeCoin],
		Diamanta:            p.tokens[constant.ItemTypeGuildContribution],
		Coupon:              p.tokens[constant.ItemTypeCoupon],
		Soul:                p.tokens[constant.ItemTypeSoul],
		TrialVoucher:        p.tokens[constant.ItemTypeTrialVoucher],
		TreasuerBox:         p.tokens[constant.ItemTypeTreasureNum],
		PetExp:              p.tokens[constant.ItemTypePetExp],
		BoxQuick:            p.tokens[constant.ItemTypeAccelerateScroll],
		GuildContribution:   p.tokens[constant.ItemTypeGuildContribution],
		GuildWish:           p.tokens[constant.ItemTypeGuildWish],
		VipExp:              p.tokens[constant.ItemTypeVipExp],
		StarWishCoin:        p.tokens[constant.ItemTypeStarWishCoin],
		Energy:              p.energy,
		Exp:                 p.exp,
		Level:               p.level,
		Vip:                 p.vip,
		EnergyRecoveryTime:  int32(p.nextEnergyRecoveryTime),
		Combat:              p.combat,
		FateLevel:           p.fateLevel,
		FateExp:             p.fateExp,
		HomeAFK:             int64(global.GlobalCfg.BaseValue.Onlinebattle.Basetime) + vitality - p.vitalityTime,
		GuildId:             p.guildId,
		WishRecoveryTime:    int32(time.Now().Unix() + int64(p.nextWishRecoveryTime)),
		MainBattleId:        p.mainBattleId,
		TitleId:             p.equipTitleId,
		ModuleIds:           p.moduleOpen,
		CreateDay:           p.createDay,
		PartnerTalentPoints: p.tokens[constant.ItemTypePartnerTalentPoint],
	}
	p.pushMsg(protocol.RouteSyncUserBasePush, packet)
}

func (p *Player) SyncUserAdvance() {
	// 其他信息
	packet := &pb.SyncUserAdvanceInfoPush{
		Skills: []*pb.SkillInfo{}, AutoDisassembly: &pb.AutoDisassemblyInfo{},
		TurnBasedUids: []int64{}, Guidance: []int32{}, Box: &pb.Userbox{},
	}
	for _, v := range p.skills {
		if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.SType == constant.SkillTypeArtifact {
			continue
		}
		packet.Skills = append(packet.Skills, &pb.SkillInfo{CfgId: v.CfgId, Cooldown: 0, Advanced: v.AdvancedLevel})
	}
	for _, v := range p.equips {
		Gemstone := make([]int32, 0)
		for _, v2 := range v.Gemstone {
			Gemstone = append(Gemstone, int32(v2))
		}
		slot := &pb.SlotInfo{Stars: v.Stars, Enhance: v.Enhance, Item: p.getItemByUid(v.ItemId), Gemstone: Gemstone}
		packet.Equips = append(packet.Equips, slot)
	}
	// 技能
	for _, v := range p.skillsSlot {
		packet.SkillsSlot = append(packet.SkillsSlot, v)
	}
	packet.Tasks = p.tasks.toPB().GetTasks()

	packet.Guidance = append(packet.Guidance, p.Guidance...)

	packet.AutoDisassembly.Open = p.ADOpen
	packet.AutoDisassembly.Level = p.ADLevel
	packet.AutoDisassembly.Quality = pb.ItemQuality(p.ADQuality)
	packet.AutoDisassembly.HighCombat = p.ADHighCombat
	packet.AutoDisassembly.QualitySwitch = p.ADQualitySwitch
	packet.AutoDisassembly.LevelSwitch = p.ADLevelSwitch
	// 回合制阵容
	for _, v := range p.turnBasedSlot {
		packet.TurnBasedUids = append(packet.TurnBasedUids, v)
	}

	p.pushMsg(protocol.RouteSyncUserAdvanceInfoPush, packet)
}

func (p *Player) syncUserSkillSlot() {
	packet := &pb.SyncUserSkillSlotPush{}
	for _, v := range p.skillsSlot {
		packet.SkillsSlot = append(packet.SkillsSlot, v)
	}
	p.pushMsg(protocol.SyncUserSkillSlotPush, packet)
}

func (p *Player) syncParterListPush() {
	packet := &pb.SyncPartnerListPush{}
	for _, v := range p.partners {
		partner := &pb.PartnerInfo{
			Uid:   v.uid,
			CfgId: v.cfgId,
			Level: v.level,
			Star:  v.star,
		}
		partner.Skills = []*pb.SkillInfo{}
		for _, v := range v.skills {
			partner.Skills = append(partner.Skills, &pb.SkillInfo{CfgId: v.CfgId, Cooldown: 0, Advanced: v.AdvancedLevel})
		}
		partner.Points = v.points
		for _, v2 := range v.activated {
			partner.Talents = append(partner.Talents, &pb.KeyValue{K: v2.ID, V: v2.Num})
		}

		packet.List = append(packet.List, partner)
	}
	p.pushMsg(protocol.RouteSyncPartnerListPush, packet)
}

// 同步宝箱信息
func (p *Player) SyncUserBox() {
	packet := &pb.SyncUserBoxInfoPush{Box: &pb.Userbox{}}
	packet.Box.Nowlevel = p.box.NowLevel
	// packet.Box.NowOpenLevel = p.box.NowOpenLevel
	// packet.Box.Job = p.box.Job
	// 箱子，这里本来是给他还剩多少，但是因为前端的问题，需要给他时间戳
	if p.box.NextBoxRecoveryTime <= 0 {
		packet.Box.NextBoxRecoveryTime = 0
	} else {
		packet.Box.NextBoxRecoveryTime = time.Now().Unix() + p.box.NextBoxRecoveryTime
	}
	packet.Box.NowOpenRes = make([]*pb.ItemInfo, 0)
	for i := range p.box.NowOpenRes {
		v := &p.box.NowOpenRes[i]
		packet.Box.NowOpenRes = append(packet.Box.NowOpenRes, v)
	}
	p.pushMsg(protocol.RouteSyncUserBoxInfoPush, packet)
}

// 同步宠物信息
func (p *Player) SyncUserPet(needDel int32) {
	//packet := &pb.SyncUserPetInfoPush{Pet: []*pb.PetSlotInfo{}, PetRecoveryTime: 0, PetEverHave: []int32{}}
	//packet.Pet = []*pb.PetSlotInfo{}
	//// 赋值
	//for i := range p.pet {
	//	v := p.pet[i]
	//	packet.Pet = append(packet.Pet, &pb.PetSlotInfo{
	//		Uid:     v.PetUid,
	//		CfgId:   v.PetCfgId,
	//		Star:    v.Star,
	//		Level:   v.Level,
	//		Exp:     v.Exp,
	//		ExpStar: v.ExpStar,
	//	})
	//}
	//checkNum := int64(0)
	//// 当前时间+下一次恢复时间
	//checkNum = int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) + int64(needDel)
	//if p.checkTakeSource(constant.ItemTypePetNormal, checkNum) {
	//	packet.PetRecoveryTime = 0
	//} else {
	//	// 这里是给前端的时间，就是当前时间加上需要等待的时间
	//	packet.PetRecoveryTime = int32(time.Now().Unix() + int64(p.nextPetRecoveryTime))
	//}
	//logger.Infof("SyncUserPetInfoPush PetRecoveryTime:%d, nextPetRecoveryTime:%d", packet.PetRecoveryTime, p.nextPetRecoveryTime)
	//packet.PetEverHave = p.petEverHave
	//p.pushMsg(protocol.RouteSyncUserPetInfoPush, packet)
	p.PetMgr.SyncUserPet(needDel)
}

// 同步头像和头像框
func (p *Player) SyncUserAvatar(refresh bool) {
	var rankResetFlag bool
	// 判断头像和头像框是否过期，过期则删除
	currentTime := time.Now().Unix()

	delAvatarIDs := []int32{}
	p.avatar.ownedAvatars.Range(func(avatarID int32, expiry int64) bool {
		if expiry != 0 && expiry <= currentTime {
			delAvatarIDs = append(delAvatarIDs, avatarID)
			if p.avatar.currentAvatarID == avatarID {
				p.avatar.currentAvatarID = global.GlobalCfg.BaseValue.InitialIcon.Avatar // 重置当前头像
				rankResetFlag = true
			}
		}
		return true
	})
	for _, avatarID := range delAvatarIDs {
		p.avatar.ownedAvatars.Delete(avatarID)
	}

	delFrameIDs := []int32{}
	p.avatar.ownedFrames.Range(func(frameID int32, expiry int64) bool {
		if expiry != 0 && expiry <= currentTime {
			delFrameIDs = append(delFrameIDs, frameID)
			if p.avatar.currentFrameID == frameID {
				p.avatar.currentFrameID = global.GlobalCfg.BaseValue.InitialIcon.Avatarframe // 重置当前头像框
				rankResetFlag = true
			}
		}
		return true
	})
	for _, frameID := range delFrameIDs {
		p.avatar.ownedFrames.Delete(frameID)
	}

	if rankResetFlag {
		// 设置排行榜信息
		p.dirty(true)
		p.userShowRefresh()
	}

	if refresh {
		p.refreshStats()
	}

	p.pushMsg(protocol.RouteSyncUserAvatarPush, p.avatar.ToProto())
}

func (p *Player) SyncUserEquipsInfo() {
	packet := &pb.SyncUserEquipsInfoPush{Equips: []*pb.SlotInfo{}}
	for _, v := range p.equips {
		var Gemstone []int32
		for _, v2 := range v.Gemstone {
			Gemstone = append(Gemstone, int32(v2))
		}
		slot := &pb.SlotInfo{Stars: v.Stars, Enhance: v.Enhance, Item: p.getItemByUid(v.ItemId), Gemstone: Gemstone}
		packet.Equips = append(packet.Equips, slot)
	}
	p.pushMsg(protocol.RouteSyncUserEquipsInfoPush, packet)
}

// 伙伴槽位
func (p *Player) SyncUserPartnerSlot() {
	packet := &pb.PartnerSlotUpdatePush{List: []*pb.PartnerSlotInfo{}}
	for _, v := range p.partnerSlots {
		packet.List = append(packet.List, &pb.PartnerSlotInfo{
			SlotId:     v.SlotId,
			PartnerUid: v.PartnerUid,
			Equips:     p.FillEquipInfo(v.Equips),
			TitleId:    v.TitleId,
		})
	}
	p.pushMsg(protocol.RoutePartnerSlotUpdatePush, packet)
}

func (p *Player) FillEquipInfo(equips [pb.EquipSlot_MaxEquipSlot]Slot) []*pb.SlotInfo {
	eslot := make([]*pb.SlotInfo, 0)
	for _, v := range equips {
		Gemstone := make([]int32, 0)
		for _, v2 := range v.Gemstone {
			Gemstone = append(Gemstone, int32(v2))
		}
		slot := &pb.SlotInfo{Stars: v.Stars, Enhance: v.Enhance, Item: p.getItemByUid(v.ItemId), Gemstone: Gemstone}
		eslot = append(eslot, slot)
	}
	return eslot
}

func (p *Player) online() {
	if p.bCreate {
		p.tasks.init()
		p.AutoDisassembly.ADOpen = false
		p.AutoDisassembly.ADLevelSwitch = global.GlobalCfg.AutoDecompose.LevelSwitch == 1
		p.AutoDisassembly.ADLevel = int32(global.GlobalCfg.AutoDecompose.Level)
		p.AutoDisassembly.ADQualitySwitch = global.GlobalCfg.AutoDecompose.QualitySwitch == 1
		p.AutoDisassembly.ADQuality = int32(global.GlobalCfg.AutoDecompose.Quality)
		p.AutoDisassembly.ADHighCombat = global.GlobalCfg.AutoDecompose.KeepBetterSwitch == 1

		time.AfterFunc(time.Second, func() {
			p.levelGift()
		})
		for _, v := range p.partnerSlots {
			slot := model.PartnerSlot{
				PartnerUid:   0,
				Uid:          p.uid,
				SlotId:       v.SlotId,
				EquipTitleid: 0,
			}
			data, err := json.Marshal(v.Equips)
			if err == nil {
				slot.Equips = string(data)
			}
			if err := db.InsertPartnerSlot(&slot); err != nil {
				return
			}
		}

		instance_data := model.InstanceData{
			Uid:               p.uid,
			BattleCurStage:    0,
			BattleRewardStage: 0,
		}
		if err := db.AddInstanceData(&instance_data); err != nil {
			p.logger.Errorf("玩家副本数据同步错误,userId: %d Error=%v", p.uid, err)
		}
		p.PVPDuel.pvpScore = DefaultPvpScore
		p.PVPMelee.gvgScore = DefaultPvpScore

	} else {
		//会长登录发走马灯
		guildManager.GuildLoginPush(p.uid)

		//排行榜榜首登陆
		defaultManager.RankLogin(p.serverID, p.uid)

	}
	p.curState = PlayerStatusIdle

	p.offlineRecovery()
	p.syncUserBase()
	p.skillAutoLearn()
	p.SyncUserAdvance()
	p.SyncUserPartnerSlot()
	p.PetMgr.online()

	p.tasks.check()
	p.syncUserSkillSlot()
	p.syncParterListPush()

	//p.refreshStats()
	p.currentHP = p.stat(pb.Attribute_MAX_HEALTH)
	p.defaultBuff()
	p.syncBuff()

	db.Online(p.uid)

	p.wellenOnlineTime = 0

	p.dirty(false)
	p.SyncUserPet(0)
	p.SyncUserBox()
	p.bCreate = false
	p.bRobot = false
	p.lastLoginAt = time.Now()

	// robotMgr.updateOfflinePlayer(p.uid)
	// 登录上报数据
	p.upLoginShushuDate()
	// if err := p.activityMgr.Load(p.uid); err != nil {
	// 	p.logger.Errorf("加载活动信息失败 err: %s", err.Error())
	// }

	p.checkPartnerBonds()
	p.SyncUserAvatar(false)
	p.PushModuleNew()
	p.initArtifact()
	p.syncArtifact()
	p.logger.Infof("玩家上线,uid: %d , name: %s, serverId: %d ", p.uid, p.name, p.serverID)
}

func (p *Player) disconnect() {
	p.lastLoginAt = time.Now()
	p.disconnectTime = time.Now()

	if d := p.desk; d != nil {
		p.desk.chLeave <- &Leave{Uid: p.uid, Disconnect: true}
		p.setPlayerStatus(PlayerStatusDisconnect)
	} else {
		p.setPlayerStatus(PlayerStatusOffLine)
		// 移除session
		p.removeSession()
	}

	p.equipRefineReplace(true)
	p.asyncSaveDb()
	go db.Offline(p.uid)

	logger.Debugf("玩家网络断开!!! uid: %d name: %s ", p.uid, p.name)
}
func (p *Player) syncUserDay() {
	// 当前时间减去注册时间，算出当天是第几天
	if p.registerAt > 0 {
		t1 := time.Unix(p.registerAt, 0)
		t2 := time.Now()
		y1, m1, d1 := t1.Date()
		y2, m2, d2 := t2.Date()
		start := time.Date(y1, m1, d1, 0, 0, 0, 0, t1.Location())
		end := time.Date(y2, m2, d2, 0, 0, 0, 0, t2.Location())
		p.createDay = int32(end.Sub(start).Hours()/24) + 1
	}
}

func (p *Player) upLoginShushuDate() {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "role_login"
	shushumain.ReportData(shushuDate)
	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["login_time"] = time.Now().Unix()
	shushuUserSetDate["current_day"] = p.nowDay
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (p *Player) dirty(v bool) {
	p.bDirty = v
}

// 异步保存数据库
func (p *Player) asyncSaveDb() {
	if p.curState == PlayerStatusPlaying {
		return
	}

	async.Run(func() {
		p.syncSaveDb()
	})

	p.dirty(false)
}

// 保存数据库
func (p *Player) syncSaveDb() {
	p.lock.Lock()
	defer p.lock.Unlock()

	u := &model.User{}

	u.Uid = p.Uid()

	u.Energy = p.energy
	u.Recharge = p.recharge
	u.Vip = p.vip
	// u.VipExp = p.vipExp
	u.Level = p.level
	u.Exp = p.exp
	u.Name = p.name
	u.Job = int(p.job)
	u.Sex = p.sex
	u.AccountName = p.accountName
	u.EnergyRecoveryTime = p.energyRecoveryTime.Unix()
	u.DailySoul = int(p.dailySoul)
	u.SessionSoul = int(p.sessionSoul)
	u.OfflineTime = p.offlineTime
	u.PvpScore = int(p.pvpScore)
	u.GvgScore = int(p.gvgScore)
	u.Nowday = p.nowDay
	u.MatchTimesDuel = int(p.PVPDuel.dailyMatchTimes)
	u.BuyMatchNumDuel = int(p.PVPDuel.dailyBuyNum)

	u.MatchTimesMelee = int(p.PVPMelee.dailyMatchTimes)
	u.BuyMatchNumMelee = int(p.PVPMelee.dailyBuyNum)

	u.SoloTowerLayers = int(p.SoloTowerInfo.Layers)
	u.SoloTowerRewardTime = p.SoloTowerInfo.RewardTime.Unix()

	u.FateExp = p.fateExp
	u.FateLevel = p.fateLevel

	u.WorldBossTimes = int(p.dailyWorldBossTimes)
	u.WorldBossDamage = p.totalWorldBossDamage
	u.DayOnlineTime = int64(p.dayOnlineTime)
	u.TotalOnlineTime = int64(p.totalOnlineTime)
	u.VitalityTime = p.vitalityTime
	u.TotalCombat = p.totalCombat
	u.AfkMultiplierID = p.AfkMultiplierID
	u.AfkExpirationAt = p.afkExpirationAt.Unix()
	u.AfkDailyTime = p.afkDailyTime
	// u.AdTotalToday = p.AdTotalToday
	u.AdTotalToday = atomic.LoadInt32(&p.adTotalToday)
	u.MuteUntil = p.muteUntil
	data, err := protobuf.NewSerializer().Marshal(p.bag)
	if err == nil {
		u.Bag = data
		if len(u.Bag) == 0 {
			u.Bag = nil
		}
	}
	uBox := &model.UserBox{}
	uBox.Uid = p.Uid()
	uBox.NowLevel = p.box.NowLevel
	uBox.NextBoxRecoveryTime = p.box.NextBoxRecoveryTime
	uBox.BeginAt = p.box.BoxRecoveryTime.Unix()
	uBox.IsAccelerate = p.box.IsAccelerate
	uBox.AccelerateNum = p.box.AccelerateNum
	uBox.CostGold = p.box.CostGold
	uBox.TimeCost = p.box.TimeCost
	uBox.OfflineTime = p.box.OfflineTime
	// uBox.NowOpenLevel = p.box.NowOpenLevel
	// uBox.Job = p.box.Job
	data, err = json.Marshal(p.box.NowOpenRes)
	if err == nil {
		uBox.NowOpenRes = string(data)
	}
	//uPet := &model.UserPet{}
	//uPet.Uid = p.Uid()
	//if !p.petRecoveryTime.IsZero() {
	//	uPet.PetRecoveryTime = p.petRecoveryTime.Unix()
	//}
	//if !p.petRecoveryTimeGold.IsZero() {
	//	uPet.PetRecoveryTimeGold = p.petRecoveryTimeGold.Unix()
	//}
	//if !p.petRecoveryTimeExp.IsZero() {
	//	uPet.PetRecoveryTimeExp = p.petRecoveryTimeExp.Unix()
	//}
	//uPet.PetExp = p.petExp
	//uPet.PetGold = p.petGold
	//data, err = json.Marshal(p.pet)
	//if err == nil {
	//	uPet.PetInfo = string(data)
	//}
	//data, err = json.Marshal(p.petEverHave)
	//if err == nil {
	//	uPet.PetEverHave = string(data)
	//}
	data, err = json.Marshal(p.skills)
	if err == nil {
		u.Skill = string(data)
	}

	data, err = json.Marshal(p.equips)
	if err == nil {
		u.Equips = string(data)
	}
	data, err = json.Marshal(p.skillsSlot)
	if err == nil {
		u.SkillSlot = string(data)
	}

	data, err = protobuf.NewSerializer().Marshal(p.tasks.toPB())
	if err == nil {
		u.Task = data
		if len(u.Task) == 0 {
			u.Task = nil
		}
	}

	u.Buff = p.toDbBuff()

	data, err = json.Marshal(p.dailyMaps)
	if err == nil {
		u.DailyMap = string(data)
	}

	data, err = json.Marshal(p.DailyItemUsed)
	if err == nil {
		u.DailyItemUsed = string(data)
	}

	data, err = json.Marshal(p.AutoDisassembly)
	if err == nil {
		u.Disassembly = string(data)
	}
	// 代币
	data, err = json.Marshal(p.tokens)
	if err == nil {
		u.Tokens = string(data)
	}
	// 是否第一次：从 SafeMap 转换回普通 map 再序列化
	tempMap := make(map[string]int64)
	p.firstUp.Range(func(key string, value int64) bool {
		tempMap[key] = value
		return true
	})
	data, err = json.Marshal(tempMap)
	if err == nil {
		u.FirstUp = string(data)
	}
	// 代币消耗
	data, err = json.Marshal(p.tokenDeletes)
	if err == nil {
		u.TokenDeletes = string(data)
	}
	data, err = json.Marshal(p.firstUse)
	if err == nil {
		u.FirstUse = string(data)
	}
	data, err = json.Marshal(p.moduleOpen)
	if err == nil {
		u.ModuleOpen = string(data)
	}

	u.ModuleNew = p.toDbModuleNew()

	for i := 0; i < int(MaxPartnerNum); i++ {
		p.syncSavePartnerSlotDb(p.partnerSlots[i].PartnerUid, p.partnerSlots[i].SlotId, p.partnerSlots[i].TitleId)
	}
	for _, v := range p.partners {
		p.syncSavePartnerDb(v.Uid(), model.Update)
	}

	data, err = json.Marshal(p.turnBasedSlot)
	if err == nil {
		u.TurnBasedSlot = string(data)
	}

	data, err = json.Marshal(p.Guidance)
	if err == nil {
		u.Guidance = string(data)
	}

	if err := db.UpdateUser(u); err != nil {
		p.logger.Errorf("保存角色信息错误, Error=%v", err)
		return
	}
	if err := db.UpdateUserBox(uBox); err != nil {
		p.logger.Errorf("保存用户开箱信息错误, Error=%v", err)
		return
	}
	//if err := db.UpdateUserPet(uPet); err != nil {
	//	p.logger.Errorf("保存用户宠物信息错误, Error=%v", err)
	//	return
	//}
	if p.PetMgr.syncSaveDb() != nil {
		return
	}
	if p.guildId > 0 {
		uGuild := &model.GuildMember{}
		uGuild.GuildID = p.guildId
		uGuild.UserID = p.uid
		uGuild.RecoveryTimeExp = p.wishRecoveryTime.Unix()
		uGuild.GuildWish = int32(p.guildWish)
		uGuild.GuildContribution = int32(p.guildContribution)
		if err := db.UpdateUserGuild(uGuild); err != nil {
			p.logger.Errorf("保存用户工会信息错误, Error=%v", err)
			return
		}
	}
	if err := p.syncSaveInstanceDataDb(); err != nil {
		p.logger.Errorf("保存玩家实例信息错误, Error=%v", err)
		return
	}

	uAttribute := &model.UserAttribute{
		Uid:         p.uid,
		ElementMain: p.ElementInfo.elementMain,
	}
	data, err = protobuf.NewSerializer().Marshal(p.ElementInfoElementToPb())
	if err == nil {
		uAttribute.Elements = data
	}
	data, err = protobuf.NewSerializer().Marshal(p.ElementInfoElementInworkToPb())
	if err == nil {
		uAttribute.ElementInWork = data
	}

	if err := db.UpdateUserAttribute(uAttribute); err != nil {
		p.logger.Errorf("保存玩家元素信息错误, Error=%v", err)
		return
	}

	if err := p.syncSaveWarDataDb(); err != nil {
		p.logger.Errorf("保存玩家群雄战场信息错误, Error=%v", err)
		return
	}
	ULotteryInfo := &model.Lottery{
		Uid:         p.uid,
		RefreshTime: p.LotteryInfo.RefreshTime,
		Progress:    p.LotteryInfo.Progress,
		NeedItemId:  p.LotteryInfo.NeedItemId,
		FreeNum:     p.LotteryInfo.FreeNum,
	}
	if data, err := json.Marshal(p.LotteryInfo.PrizePool); err == nil {
		ULotteryInfo.PrizePool = string(data)
	}
	if data, err := json.Marshal(p.LotteryInfo.HavePrizePool); err == nil {
		ULotteryInfo.HavePrizePool = string(data)
	}
	//保存抽奖
	if err := db.UpdateUserLottery(ULotteryInfo); err != nil {
		p.logger.Errorf("保存玩家抽奖信息错误, Error=%v", err)
		return
	}

	p.avatar.AddInitialAvatarAndFrame()
	if err := db.UpdateAvatar(p.avatar.ToDB(p)); err != nil {
		p.logger.Errorf("syncSaveDb 保存玩家头像信息错误, Error=%v", err)
	}

	if err := db.UpdateRefine(p.refine.ToDB(p)); err != nil {
		p.logger.Errorf("syncSaveDb 保存玩家精炼信息错误, Error=%v", err)
	}

	if err := p.saveArtifactDB(); err != nil {
		p.logger.Errorf("syncSaveDb 保存玩家神兵信息错误, Error=%v", err)
	}

	p.logger.Debugf("保存数据库, name=%s  level=%d", p.name, p.level)
}

// 伙伴相关
func (p *Player) syncSavePartnerSlotDb(puid int64, sid int32, titleid int32) error {
	slot := model.PartnerSlot{
		Uid:          p.uid,
		PartnerUid:   puid,
		SlotId:       sid,
		EquipTitleid: titleid,
	}
	data, err := json.Marshal(p.getPartnerSlotInfo(sid).Equips)
	if err == nil {
		slot.Equips = string(data)
	}

	if err := db.UpdatePartnerSlot(&slot); err != nil {
		return err
	}
	return nil
}

func (p *Player) syncSavePartnerDb(puid int64, optype model.DBOpt) error {

	if p.isRobot() {
		return nil
	}

	async.Run(func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		info, ok := p.partners[puid]
		if !ok {
			return
		}
		switch optype {
		case model.Add:
			actiTelnet, _ := json.Marshal(info.activated)
			partners := model.Partner{
				Uid:             p.uid,
				PartnerUid:      puid,
				CfgId:           info.cfgId,
				Level:           info.level,
				Star:            info.star,
				TalnetPoints:    info.points,
				TalnetActivated: string(actiTelnet),
			}

			data, err := json.Marshal(info.skills)
			if err == nil {
				partners.Skill = string(data)
			}

			async.Run(func() {
				if err := db.InsertPartner(&partners); err != nil {
					p.logger.Errorf("DB新增玩家伙伴信息失败:%v", err)
				}
			})

		case model.Del:
			partners := model.Partner{
				Uid:        p.uid,
				PartnerUid: puid,
			}

			async.Run(func() {
				if err := db.DeletePartner(&partners); err != nil {
					p.logger.Errorf("DB删除玩家伙伴信息失败:%v", err)
					return
				}
			})

		case model.Update:
			actiTelnet, _ := json.Marshal(info.activated)
			partners := model.Partner{
				Uid:             p.uid,
				PartnerUid:      puid,
				Level:           info.level,
				Star:            info.star,
				TalnetPoints:    info.points,
				TalnetActivated: string(actiTelnet),
			}
			data, err := json.Marshal(info.skills)
			if err == nil {
				partners.Skill = string(data)
			}

			async.Run(func() {
				if err := db.UpdatePartner(&partners); err != nil {
					p.logger.Errorf("DB更新玩家伙伴信息失败:%v", err)
				}
			})
		}
	})

	return nil
}

func (p *Player) syncSaveInstanceDataDb() error {
	instanceData := model.InstanceData{
		Uid:               p.uid,
		BattleCurStage:    p.BattleMap.BattleId,
		BattleRewardStage: p.BattleMap.RewardStage,

		BossCurStage:      p.BossMap.BattleId,
		BossChallengeTime: p.BossMap.ChallengeTime,

		TrialTowerCurStage: p.TrialTowerMap.BattleId,

		RelicCurStage:      p.RelicMap.BattleId,
		RelicChallengeTime: p.RelicMap.ChallengeTime,

		FightEndCurStage:     p.FightEndMap.BattleId,
		FightEndRefreshTime:  p.FightEndMap.RefreshAt.Unix(),
		ElementCurStage:      p.ElementMap.BattleId,
		ElementChallengeTime: p.ElementMap.ChallengeTime,
		ElementNeedSweep:     p.ElementMap.NeedSweep,

		MainBattleId:                   p.mainBattleId,
		PushingWallCount:               p.PushingWall.UsedCount,
		PushingWallTicketPurchaseCount: p.PushingWall.TicketPurchaseCount,
	}
	data, err := json.Marshal(p.RelicMap.Evaluation)
	if err == nil {
		instanceData.RelicEvaluation = string(data)
	}
	data, err = json.Marshal(p.RelicMap.ChapterReward)
	if err == nil {
		instanceData.RelicChapterReward = string(data)
	}

	data, err = json.Marshal(p.FightEndMap.Evaluation)
	if err == nil {
		instanceData.FightEndEvaluation = string(data)
	}

	data, err = json.Marshal(p.ChallengeRecord.Records)
	if err == nil {
		instanceData.ChallengeRecord = string(data)
	}
	data, err = json.Marshal(p.MonKillRecord.MonKillRecords)
	if err == nil {
		instanceData.KillMonControl = string(data)
	}

	if err := db.UpdateInstanceData(&instanceData); err != nil {
		return err
	}
	return nil
}

// 保存war 数据
func (p *Player) syncSaveWarDataDb() error {
	warData := model.War{Uid: p.uid, Soul: p.warDailySoul, EquipTitleid: p.equipTitleId}

	data, err := protobuf.NewSerializer().Marshal(p.titles)
	if err == nil {
		warData.Title = data
	}

	if err := db.UpdateWar(&warData); err != nil {
		return err
	}
	return nil
}

// 获取道具数量 By CfgID
func (p *Player) getItemCountInBag(id int32) int64 {
	var count int64
	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		return 0
	}

	if itemCfg.MainType == constant.NotInBag {
		count, ok = p.tokens[id]
		if !ok {
			count = 0
		}
	} else {
		for _, item := range p.bag.ItemList {
			if id == item.CfgId {
				count += item.Num
			}
		}
	}

	return count
}

// 包裹剩余格子
func (p *Player) bagLeftCount(itemType pb.ItemType) int {
	if itemType == pb.ItemType_NOT_IN_BAG {
		if len(p.bag.ItemList)-p.getEquipNum() >= global.GlobalCfg.BaseValue.BagNum.TotleNum {
			return 0
		}
		return global.GlobalCfg.BaseValue.BagNum.TotleNum - len(p.bag.ItemList) - p.getEquipNum()
	}

	num := 0
	if itemType == pb.ItemType_EQUIPMENT {
		for _, item := range p.bag.GetItemList() {
			if item.Type == pb.ItemType_EQUIPMENT {
				num++
			}
		}

		if num-p.getEquipNum() >= global.GlobalCfg.BaseValue.BagNum.EquipNum {
			return 0
		}
		return global.GlobalCfg.BaseValue.BagNum.EquipNum - num + p.getEquipNum()
	} else {
		for _, item := range p.bag.GetItemList() {
			if item.Type != pb.ItemType_EQUIPMENT {
				num++
			}
		}

		if num >= global.GlobalCfg.BaseValue.BagNum.ItemNum {
			return 0
		}
		return global.GlobalCfg.BaseValue.BagNum.ItemNum - num
	}
}

func (p *Player) bagEquipLeftCount() int {
	num := 0
	for _, item := range p.bag.GetItemList() {
		if item.Type == pb.ItemType_EQUIPMENT {
			num++
		}
	}

	if num-p.getEquipNum() >= global.GlobalCfg.BaseValue.BagNum.EquipNum {
		return 0
	}
	return global.GlobalCfg.BaseValue.BagNum.EquipNum - num + p.getEquipNum()
}

func (p *Player) bagItemLeftCount() int {
	num := 0
	for _, item := range p.bag.GetItemList() {
		if item.Type != pb.ItemType_EQUIPMENT {
			num++
		}
	}

	if num >= global.GlobalCfg.BaseValue.BagNum.ItemNum {
		return 0
	}
	return global.GlobalCfg.BaseValue.BagNum.ItemNum - num
}

// 批量增加道具检查包裹空间
func (p *Player) checkAddSourceBatch(items []global.IdNum) error {
	var itemcount, equipcount, petCount int
	for _, s := range items {
		if item, err := global.GetItemCfg(s.ID); err == nil {
			if item.MainType == constant.ItemPet {
				petCount += int(s.Num)
			}

			if item.MainType != constant.NotInBag {
				if item.MainType == constant.Equipment {
					equipcount += int(s.Num)
				} else {
					num := s.Num/int32(item.Maxnum) + 1
					itemcount += int(num)
				}
			}
		} else {
			return err
		}
	}
	// p.logger.Debugf("checkAddSourceBatch 添加道具数量: %d , 装备包裹剩余格子: %d 道具包裹剩余格子: %d", len(items), p.bagEquipLeftCount(), p.bagItemLeftCount())

	targetPetNum := p.PetMgr.getEmptyPetSlotNum()
	if targetPetNum < petCount {
		return errutil.ErrPetFull
	}

	if equipcount > p.bagEquipLeftCount() || itemcount > p.bagItemLeftCount() {
		packet := &pb.SystemNotificationPush{Type: pb.Notification_BAG_FULL}
		p.pushMsg(protocol.RouteSystemNotificationPush, packet)
		return errutil.ErrBagFull
	}

	return nil
}

// 检查包裹空间
func (p *Player) checkAddSource(id int32, num int64) bool {
	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		return false
	}

	if itemCfg.MainType == constant.NotInBag {
		return true
	}

	if itemCfg.Maxnum > 0 {
		needBag := num / itemCfg.Maxnum
		if num%itemCfg.Maxnum != 0 {
			needBag++
		}

		if p.bagLeftCount(pb.ItemType(itemCfg.MainType)) < int(needBag) {
			packet := &pb.SystemNotificationPush{Type: pb.Notification_BAG_FULL}
			p.pushMsg(protocol.RouteSystemNotificationPush, packet)
			return false
		}
	}

	return true
}

// 道具数量检查 By CfgID
func (p *Player) checkTakeSource(id int32, num int64) bool {
	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		return false
	}
	num = int64(math.Abs(float64(num)))

	if itemCfg.MainType == constant.NotInBag {
		count, ok := p.tokens[id]
		if !ok {
			count = 0
		}
		if count < num {
			return false
		}

	} else {
		if p.getItemCountInBag(id) < num {
			return false
		}
	}

	return true
}

// 检查道具数量 By CfgID
func (p *Player) checkTakeSourceByUid(uid int64, num int64) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if num <= 0 {
		return false
	}

	for _, item := range p.bag.GetItemList() {
		if item.Uid == uid {
			if item.Num >= num {
				return true
			}
		}
	}
	return false
}

func (p *Player) addToken(id int32, num int64, reason string) {
	if global.IsToken(id) {
		if p.tokens == nil {
			p.logger.Errorf("增加货币失败   %s  tokens nil", p.name)
			return
		}
		if _, ok := p.tokens[id]; !ok {
			p.tokens[id] = num
		} else {
			p.tokens[id] += num
		}
		if _, ok := p.tokenDeletes[id]; !ok {
			if num < 0 {
				p.tokenDeletes[id] = int64(math.Abs(float64(num)))
			}
		} else {
			if num < 0 {
				p.tokenDeletes[id] += int64(math.Abs(float64(num)))
			}
		}
	}
	switch id {
	case constant.ItemTypeCoin:
		if num < 0 {
			p.tasks.dispatch(pb.TargetType_TOTAL_GOLD_COINS_CONSUMED_REQUIREMENT, 0, int32(math.Abs(float64(num))))
		}
	case constant.ItemTypeSoul:
		if reason == constant.ItemAddPvp {
			p.PVPInfo.dailySoul += int(num)
			p.PVPInfo.sessionSoul += int(num)

			if reason == constant.ItemAddPvp {
				// 更新排行榜
				p.rankChangeTCP(constant.RankTypePvp, int64(p.PVPInfo.dailySoul))
				p.rankChangeTCP(constant.RankTypeSeason, int64(p.PVPInfo.sessionSoul))
			}
			p.tasks.dispatch(pb.TargetType_TOTAL_ACQUIRED_SOULS, 0, int32(num))
		} else if reason == constant.ItemAddWar {
			p.warDailySoul += num
		}

	case constant.ItemTypeEnergy:
		p.energy += num
		if p.energy < 0 {
			p.energy = 0
		} else if p.energy > int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) {
			p.energy = int64(global.GlobalCfg.BaseValue.Energy.Maxvalue)
		}

	case constant.ItemTypeExp:
		p.exp += num
		cfg, err := global.GetRoleCfg(int32(p.job), p.level)
		if err == nil {
			if p.exp > cfg.NeedExp*global.GlobalCfg.BaseValue.Expstore.Rate {
				p.exp = cfg.NeedExp * global.GlobalCfg.BaseValue.Expstore.Rate
			}
		}
		if p.exp < 0 {
			p.exp = 0
		}
		// if p.getCombatPartnerNum() == 1 { // 只有主角上阵直接升级
		// 	p.levelUp()
		// }
		p.levelUp()

	case constant.ItemTypeVipExp:
		p.vipExp += int64(num)
		if num > 0 {
			p.vipLevelUp()
		}

	case constant.ItemTypeTrialVoucher:
		if num > 0 {
			p.tasks.dispatch(pb.TargetType_TOTAL_ACQUIRED_TRIAL_TOKENS, 0, int32(num))
		}

	case constant.ItemTypePVPTimes:
		p.PVPInfo.PVPDuel.dailyBuyNum += int32(num)

	case constant.ItemTypeGVGTimes:
		p.PVPInfo.PVPMelee.dailyBuyNum += int32(num)

	case constant.ItemTypeBossChallengeTimes:
		p.InstanceData.BossMap.ChallengeTime += int32(num)

	case constant.ItemTypeRelicChallengeTimes:
		p.InstanceData.RelicMap.ChallengeTime += int32(num)

	case constant.ItemTypeGuildContribution:
		p.guildContribution += num
		if p.guildContribution < 0 {
			p.guildContribution = 0
		}
	case constant.ItemTypeGuildWish:
		p.guildWish += num
		if p.guildWish < 0 {
			p.guildWish = 0
		}
	case constant.ItemTypeElementChallengeTimes:
		p.ElementMap.ChallengeTime += int32(num)
	case constant.ItemTypeBossChallenge:
		if p.InstanceData != nil {
			p.InstanceData.BossMap.ChallengeTime += int32(num)
		}
	case constant.ItemTypeRelicChallenge:
		if p.InstanceData != nil {
			p.InstanceData.RelicMap.ChallengeTime += int32(num)
		}
	case constant.ItemTypeEnhanceStoneChallenge:
		p.dailyMaps[constant.EnhanceStoneFBMap] += int32(num)
	case constant.ItemTypeStarUpStoneChallenge:
		p.dailyMaps[constant.StarUpStoneFBMap] += int32(num)
	case constant.ItemTypeSoulStoneChallenge:
		p.dailyMaps[constant.SoulStoneFBMap] += int32(num)
	case constant.ItemTypeBookPageChallenge:
		p.dailyMaps[constant.BookPageFBMap] += int32(num)
	case constant.ItemTypeGoldChallenge:
		p.dailyMaps[constant.GoldFBMap] += int32(num)
	case constant.ItemTypeExpChallenge:
		p.dailyMaps[constant.ExperienceFBMap] += int32(num)
	case constant.ItemTypeStarWishCoin:
		p.dailyMaps[constant.ItemTypeStarWishCoin] += int32(num)
	case constant.ItemTypePushWallChallengeTimes:
		p.InstanceData.PushingWall.TicketPurchaseCount += int32(num)
	}
}

// 增加资源 By CfgID
func (p *Player) addSource(id int32, num int64, reason string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var avatarFlag bool

	if !p.checkAddSource(id, num) {
		p.logger.Warnf("添加道具失败,包裹已满 userId %d  itemId: %d  itemNum: %d", p.Uid(), id, num)
		return
	}

	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		p.logger.Warnf("添加道具失败,未找到配置 userId %d  itemId: %d  itemNum: %d", p.Uid(), id, num)
		return
	}

	packet := pb.AddSourcePush{Type: AddSourceTypeNew, Item: []*pb.ItemInfo{}}

	if itemCfg.MainType == constant.NotInBag {
		p.addToken(id, num, reason)

		packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Num: num})
		p.syncUserBase()
	} else if itemCfg.MainType == constant.ItemPartner {
		for _, par_cfg := range global.PartnerCfg {
			if par_cfg.ItemID == id {
				for i := 0; i < int(num); i++ {
					p.addPartner(par_cfg.Cid)
				}
				break
			}
		}
	} else if itemCfg.MainType == constant.ItemPet {
		for _, pet_cfg := range global.PetCfg {
			if pet_cfg.ItemID == id {
				for i := 0; i < int(num); i++ {
					p.addPet(pet_cfg.ID, 0, 0, 0, 1, reason)
				}
				break
			}
		}
		packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Num: num})
	} else if itemCfg.MainType == constant.ItemMainTypeAvatar {
		p.avatar.AddAvatar(itemCfg.ID)
		avatarFlag = true
	} else if itemCfg.MainType == constant.ItemMainTypeAvatarFrame {
		p.avatar.AddFrame(itemCfg.ID)
		avatarFlag = true
	} else {
		var leftNum int64
		leftNum = num
		for _, item := range p.bag.GetItemList() {
			if item.CfgId == id {
				if item.Num < itemCfg.Maxnum {
					if num <= itemCfg.Maxnum-item.Num {
						item.Num += num
						leftNum = 0
						packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: num})
						break
					} else {
						addNum := itemCfg.Maxnum - item.Num
						leftNum = leftNum - addNum
						item.Num = itemCfg.Maxnum
						packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: addNum})
					}
				}
			}
		}
		for leftNum > 0 {
			// 道具获取时间，如果Timeliness=1，取获取时间当天0点
			getTime := time.Now().Unix()
			if itemCfg.Timeliness.Type == constant.ItemTimelinessShiChang {
				currentTime := time.Now()
				year, month, day := currentTime.Date()
				newTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
				getTime = newTime.Unix()
			}
			if leftNum > itemCfg.Maxnum {
				leftNum -= itemCfg.Maxnum
				itemUid := snowflake.GenID()
				newItem := &pb.ItemInfo{CfgId: id, Uid: itemUid, Num: itemCfg.Maxnum, CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType)}

				// 随机品质&随机属性
				p.randomEquipQuality(newItem)
				//p.upgradeRandomAttr(newItem, itemCfg)
				packet.Item = append(packet.Item, newItem)
				p.bag.ItemList = append(p.bag.ItemList, newItem)
			} else {
				itemUid := snowflake.GenID()
				newItem := &pb.ItemInfo{CfgId: id, Uid: itemUid, Num: leftNum, CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType)}

				// 随机品质&随机属性
				p.randomEquipQuality(newItem)
				//p.upgradeRandomAttr(newItem, itemCfg)
				packet.Item = append(packet.Item, newItem)
				p.bag.ItemList = append(p.bag.ItemList, newItem)
				leftNum = 0
			}
		}
	}
	if reason != constant.ItemAddAFKBox && reason != constant.ItemAddAFKMonster && reason != constant.ItemAddDisintegrate &&
		reason != constant.ItemAddDrop && reason != constant.ItemAddPetExp && reason != constant.ItemAddPetAFK && reason != constant.ItemAddPetGold {
		sourceInfo := make([]global.IdNum, 0)
		sourceInfo = append(sourceInfo, global.IdNum{ID: id, Num: int32(num)})
		p.upAddSourceShushuDate(sourceInfo, reason)
	}
	p.pushMsg(protocol.RouteAddSourcePush, &packet)

	p.dirty(true)

	if avatarFlag {
		p.SyncUserAvatar(true)
	}

	// p.logger.Infof("addSource 增加道具  ItemCfgID:%d , Num:%d ,Reason:%s", id, num, reason)
}

// 添加指定品质装备
func (p *Player) addEquipWithQuality(id int32, quality pb.ItemQuality, reason string) *pb.ItemInfo {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.checkAddSource(id, 1) {
		p.logger.Warnf("添加道具失败,包裹已满 userId %d  itemId: %d ", p.Uid(), id)
		return nil
	}

	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		p.logger.Warnf("添加道具失败,未找到配置 userId %d  itemId: %d  ", p.Uid(), id)
		return nil
	}

	if itemCfg.MainType != constant.Equipment {
		return nil
	}

	packet := pb.AddSourcePush{Type: AddSourceTypeNew, Item: []*pb.ItemInfo{}}
	itemUid := snowflake.GenID()
	getTime := time.Now().Unix()
	newItem := &pb.ItemInfo{CfgId: id, Uid: itemUid, Num: 1, CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType), Equip: &pb.EquipInfo{Quality: quality}}
	//p.upgradeRandomAttr(newItem, itemCfg)
	packet.Item = append(packet.Item, newItem)
	p.bag.ItemList = append(p.bag.ItemList, newItem)

	p.pushMsg(protocol.RouteAddSourcePush, &packet)

	p.dirty(true)

	if reason != constant.ItemAddAFKBox && reason != constant.ItemAddAFKMonster && reason != constant.ItemAddDisintegrate &&
		reason != constant.ItemAddDrop && reason != constant.ItemAddPetExp && reason != constant.ItemAddPetAFK && reason != constant.ItemAddPetGold {
		sourceInfo := make([]global.IdNum, 0)
		sourceInfo = append(sourceInfo, global.IdNum{ID: id, Num: 1})
		p.upAddSourceShushuDate(sourceInfo, reason)
	}
	return newItem
}

func (p *Player) addEquipBatchWithQuality(items []global.IdNumQuality, reason string, opts ...AddSourceBatchOption) []*pb.ItemInfo {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(items) == 0 {
		return nil
	}

	if p.bagEquipLeftCount() < len(items) {
		p.logger.Warnf("添加道具失败,包裹已满 %v", items)
		return nil
	}

	sourceInfo := make([]global.IdNum, 0)
	cfg := newAddSourceBatchConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	packet := pb.AddSourcePush{Type: int32(cfg.useNotifyType), Item: []*pb.ItemInfo{}}

	for _, v := range items {
		itemCfg, ok := global.ItemCfg[v.ID]
		if !ok {
			p.logger.Warnf("添加道具失败,未找到配置  item: %v  ", v)
			continue
		}
		if itemCfg.MainType != constant.Equipment {
			p.logger.Warnf("添加道具失败,不是装备  item: %v  ", v)
			continue
		}
		if v.Num != 1 {
			p.logger.Warnf("添加道具失败,数量错误 item: %v  ", v)
			continue
		}
		itemUid := snowflake.GenID()
		getTime := time.Now().Unix()
		newItem := &pb.ItemInfo{CfgId: v.ID, Uid: itemUid, Num: int64(v.Num), CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType), Equip: &pb.EquipInfo{Quality: pb.ItemQuality(v.Quality)}}

		packet.Item = append(packet.Item, newItem)
		p.bag.ItemList = append(p.bag.ItemList, newItem)

		if reason != constant.ItemAddAFKBox && reason != constant.ItemAddAFKMonster && reason != constant.ItemAddDisintegrate &&
			reason != constant.ItemAddDrop && reason != constant.ItemAddPetExp && reason != constant.ItemAddPetAFK && reason != constant.ItemAddPetGold {
			sourceInfo = append(sourceInfo, global.IdNum{ID: v.ID, Num: 1})
		}
	}

	p.pushMsg(protocol.RouteAddSourcePush, &packet)

	if len(sourceInfo) > 0 {
		p.upAddSourceShushuDate(sourceInfo, reason)
	}

	p.dirty(true)

	return packet.Item
}

// 减少资源 By CfgId
func (p *Player) takeSource(id int32, num int64, args ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	packet := pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
	sourceInfo := make([]global.IdNum, 0)
	itemCfg, ok := global.ItemCfg[id]
	if !ok {
		return
	}
	if itemCfg.MainType == constant.NotInBag {
		p.addToken(id, -num, "")
		sourceInfo = append(sourceInfo, global.IdNum{ID: id, Num: int32(num)})
		packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Num: num})
		p.syncUserBase()
	} else {
		var leftNum int64
		leftNum = num
		var rmUids []int64
		for _, item := range p.bag.GetItemList() {
			if item.CfgId == id && leftNum > 0 {
				sourceInfo = append(sourceInfo, global.IdNum{ID: item.CfgId, Num: int32(leftNum)})
				if item.Num >= leftNum {
					item.Num -= leftNum
					packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: leftNum})
					if item.Num == 0 {
						rmUids = append(rmUids, item.Uid)
					}
					break
				} else {
					packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: item.Num})
					rmUids = append(rmUids, item.Uid)
					leftNum = leftNum - item.Num
					if leftNum <= 0 {
						break
					}
				}
			}
		}
		if len(rmUids) > 0 {
			newBag := pb.Bag{}
			for _, item := range p.bag.ItemList {
				bFind := false
				for _, v := range rmUids {
					if v == item.Uid {
						bFind = true
						break
					}
				}
				if !bFind {
					newBag.ItemList = append(newBag.ItemList, item)
				}
				bFind = false
			}
			p.bag = &newBag
		}

		p.tasks.dispatch(pb.TargetType_TOTAL_ITEM_CONSUMPTION, itemCfg.ID, int32(num))
	}
	method := ""
	if len(args) > 0 {
		method = args[0]
	}
	if method != constant.ItemDelHomeOpenBox && method != constant.ItemDelGemstone && method != constant.ItemDelDisintegrate {
		p.upTakeSourceShushuDate(sourceInfo, method)
	}

	p.pushMsg(protocol.RouteTakeSourcePush, &packet)
	p.dirty(true)
	// p.logger.Infof("删除道具 ItemCfgID:%d , Num:%d ,Reason:%v", id, num, args)
}

// 减少资源 By Uid
func (p *Player) takeSourceByUid(uid int64, num int64, args ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	packet := pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
	sourceInfo := make([]global.IdNum, 0)
	defer p.pushMsg(protocol.RouteTakeSourcePush, &packet)

	for i, item := range p.bag.GetItemList() {
		if item.Uid == uid {
			item.Num -= num
			packet.Item = append(packet.Item, &pb.ItemInfo{Uid: uid, Num: num, CfgId: item.CfgId})
			sourceInfo = append(sourceInfo, global.IdNum{ID: item.CfgId, Num: int32(num)})
			if item.Num == 0 {
				p.bag.ItemList = append(p.bag.ItemList[:i], p.bag.ItemList[i+1:]...)
			}

			p.dirty(true)
			p.logger.Infof("takeSourceByUid 删除道具 ItemCfgID:%d , Num:%d ,Reason:%v", item.CfgId, num, args)
			break
		}
	}
	method := ""
	if len(args) > 0 {
		method = args[0]
	}
	if method != constant.ItemDelHomeOpenBox && method != constant.ItemDelGemstone && method != constant.ItemDelDisintegrate {
		p.upTakeSourceShushuDate(sourceInfo, method)
	}
}

// 批量删除资源
func (p *Player) takeSourceByUidBatch(items []global.UUIdNum, args string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(items) == 0 {
		return
	}
	packet := pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
	defer p.pushMsg(protocol.RouteTakeSourcePush, &packet)
	sourceInfo := make([]global.IdNum, 0)
	for _, delitem := range items {
		for i, item := range p.bag.GetItemList() {
			if item.Uid == delitem.UID {
				sourceInfo = append(sourceInfo, global.IdNum{ID: item.CfgId, Num: int32(delitem.Num)})
				item.Num -= int64(delitem.Num)
				packet.Item = append(packet.Item, &pb.ItemInfo{Uid: delitem.UID, Num: int64(delitem.Num), CfgId: item.CfgId})

				if item.Num <= 0 {
					p.bag.ItemList = append(p.bag.ItemList[:i], p.bag.ItemList[i+1:]...)
				}
				break
			}
		}
	}
	if args != constant.ItemDelHomeOpenBox && args != constant.ItemDelGemstone && args != constant.ItemDelDisintegrate {
		p.upTakeSourceShushuDate(sourceInfo, args)
	}
	p.dirty(true)
	// p.logger.Infof("takeSourceByUidBatch 批量删除道具 Items:%+v ,Reason:%v", items, args)
}

func (p *Player) upTakeSourceShushuDate(sourceInfo []global.IdNum, entrance string) {
	if p.bCreate {
		return
	}
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "other_item_consume"
	cost_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(v.ID); err == nil {
				cost_item_info = append(cost_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["cost_item_info"] = cost_item_info
	if idx := strings.Index(entrance, "-"); idx != -1 {
		parts := strings.SplitN(entrance, "-", 2)
		shushuDate["entrance"] = parts[0]
		shushuDate["method"] = parts[1]
	} else {
		shushuDate["entrance"] = entrance
	}
	shushumain.ReportData(shushuDate)
}

func mergeItems(items []global.IdNum) []global.IdNum {
	m := make(map[int32]int32) // key: ID, value: 累计 Num
	for _, v := range items {
		m[v.ID] += v.Num
	}

	out := make([]global.IdNum, 0, len(m))
	for id, num := range m {
		out = append(out, global.IdNum{ID: id, Num: num})
	}
	return out
}

// AddSourceBatchOption 定义批量添加资源的选项函数
type AddSourceBatchOption func(*addSourceBatchConfig)

// addSourceBatchConfig 批量添加资源的配置
type addSourceBatchConfig struct {
	useNotifyType int // 通知类型 AddSourceTypeUnNotify
}

// WithAddSourceUnNotify 设置是否使用 AddSourceTypeUnNotify
func WithAddSourceUnNotify(useNotifyType int) AddSourceBatchOption {
	return func(cfg *addSourceBatchConfig) {
		cfg.useNotifyType = useNotifyType
	}
}

// newAddSourceBatchConfig 创建默认的批量添加资源配置
func newAddSourceBatchConfig() *addSourceBatchConfig {
	return &addSourceBatchConfig{
		useNotifyType: AddSourceTypeNew, // 默认使用 AddSourceTypeNew
	}
}

// 批量增加资源 By CfgID
func (p *Player) addSourceBatch(items []global.IdNum, reason string, opts ...AddSourceBatchOption) []*pb.ItemInfo {
	p.lock.Lock()
	defer p.lock.Unlock()

	var avatarFlag bool

	if err := p.checkAddSourceBatch(items); err != nil {
		p.logger.Warnf("批量添加道具失败,包裹已满 userId %d  itemIs: %v  ", p.Uid(), items)
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	// 应用选项配置
	cfg := newAddSourceBatchConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	//合并同类相
	items = mergeItems(items)

	packet := pb.AddSourcePush{Type: int32(cfg.useNotifyType), Item: []*pb.ItemInfo{}}

	sourceInfo := make([]global.IdNum, 0)
	for _, additem := range items {
		id := additem.ID
		num := int64(additem.Num)
		itemCfg, ok := global.ItemCfg[id]
		if !ok {
			p.logger.Warnf("添加道具失败,未找到配置 userId %d  itemId: %d  itemNum: %d", p.Uid(), id, num)
			return nil
		}

		if itemCfg.MainType == constant.NotInBag {
			p.addToken(id, num, reason)
			sourceInfo = append(sourceInfo, additem)

			packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Num: num})
			p.syncUserBase()
		} else if itemCfg.MainType == constant.ItemPartner {
			for _, par_cfg := range global.PartnerCfg {
				if par_cfg.ItemID == id {
					for i := 0; i < int(num); i++ {
						p.addPartner(par_cfg.Cid)
					}
					break
				}
			}
		} else if itemCfg.MainType == constant.ItemPet {
			for _, pet_cfg := range global.PetCfg {
				if pet_cfg.ItemID == id {
					for i := 0; i < int(num); i++ {
						p.addPet(pet_cfg.ID, 0, 0, 0, 1, reason)
					}
					break
				}
			}
			packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Num: num})
		} else if itemCfg.MainType == constant.ItemMainTypeAvatar {
			p.avatar.AddAvatar(itemCfg.ID)
			avatarFlag = true
		} else if itemCfg.MainType == constant.ItemMainTypeAvatarFrame {
			p.avatar.AddFrame(itemCfg.ID)
			avatarFlag = true
		} else {
			sourceInfo = append(sourceInfo, additem)
			var leftNum int64
			leftNum = num
			for _, item := range p.bag.GetItemList() {
				if item.CfgId == id {
					if item.Num < itemCfg.Maxnum {
						if leftNum <= itemCfg.Maxnum-item.Num {
							item.Num += leftNum
							leftNum = 0
							packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: num})
							break
						} else {
							addNum := itemCfg.Maxnum - item.Num
							leftNum = leftNum - addNum
							item.Num = itemCfg.Maxnum
							packet.Item = append(packet.Item, &pb.ItemInfo{CfgId: id, Uid: item.Uid, Num: addNum})
						}
					}
				}
			}
			for leftNum > 0 {
				// 道具获取时间，如果Timeliness=1，取获取时间当天0点
				getTime := time.Now().Unix()
				if itemCfg.Timeliness.Type == constant.ItemTimelinessShiChang {
					currentTime := time.Now()
					year, month, day := currentTime.Date()
					newTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
					getTime = newTime.Unix()
				}
				if leftNum > itemCfg.Maxnum {
					leftNum -= itemCfg.Maxnum
					itemUid := snowflake.GenID()
					newItem := &pb.ItemInfo{CfgId: id, Uid: itemUid, Num: itemCfg.Maxnum, CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType)}

					// 随机品质&随机属性
					p.randomEquipQuality(newItem)
					//p.upgradeRandomAttr(newItem, itemCfg)
					packet.Item = append(packet.Item, newItem)
					p.bag.ItemList = append(p.bag.ItemList, newItem)
				} else {
					itemUid := snowflake.GenID()
					newItem := &pb.ItemInfo{CfgId: id, Uid: itemUid, Num: leftNum, CreateTime: getTime, Type: pb.ItemType(itemCfg.MainType)}

					// 随机品质&随机属性
					p.randomEquipQuality(newItem)
					//p.upgradeRandomAttr(newItem, itemCfg)

					packet.Item = append(packet.Item,
						&pb.ItemInfo{
							CfgId:      id,
							Uid:        itemUid,
							Num:        leftNum,
							CreateTime: getTime,
							Type:       pb.ItemType(itemCfg.MainType),
							Equip:      newItem.Equip})

					// 自动分解
					if reason == constant.ItemAddDrop && itemCfg.MainType == constant.Equipment && p.ADOpen {
						bDecompose := false
						if p.ADLevelSwitch && itemCfg.Levellimit <= int(p.ADLevel) {
							bDecompose = true
						}
						if p.ADQualitySwitch && itemCfg.Quality <= p.ADQuality {
							bDecompose = true
						}
						if p.ADHighCombat && p.calcCombatEquip(newItem) > p.calcCombatEquip(p.getItemByUid(p.equips[itemCfg.Equiptype].ItemId)) {
							bDecompose = false
						}
						if bDecompose {
							// p.logger.Debugf("自动分解 穿戴装备战力 %d  新装备战力 %d ", p.calcCombatEquip(p.getItemByUid(p.equips[itemCfg.Equiptype].ItemId)), p.calcCombatEquip(newItem))
							p.ItemDisassembly(newItem)
							packet.Type = AddSourceTypeAutoDis
						}
					} else {
						p.bag.ItemList = append(p.bag.ItemList, newItem)
					}

					leftNum = 0
				}
			}
		}
	}
	if reason != constant.ItemAddOpenBox {
		if reason != constant.ItemAddAFKBox && reason != constant.ItemAddAFKMonster && reason != constant.ItemAddDisintegrate &&
			reason != constant.ItemAddDrop && reason != constant.ItemAddPetExp && reason != constant.ItemAddPetAFK && reason != constant.ItemAddPetGold {
			p.upAddSourceShushuDate(sourceInfo, reason)
		}
		p.pushMsg(protocol.RouteAddSourcePush, &packet)
	}
	p.dirty(true)

	if avatarFlag {
		p.SyncUserAvatar(true)
	}

	return packet.Item
	// p.logger.Infof("addSourceBatch 批量增加道具  Items:%+v ,Reason:%s", items, reason)
}

// 异步添加/删除道具
func (p *Player) addSourceEvent(source *protocol.AddSource) {
	p.eventMgr.Publish(event.EventTypeAddResource, source)
}

func (p *Player) upAddSourceShushuDate(sourceInfo []global.IdNum, entrance string) {
	if p.bCreate {
		return
	}
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "other_item_get"
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(v.ID); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	shushuDate["get_item_info"] = get_item_info
	if idx := strings.Index(entrance, "-"); idx != -1 {
		parts := strings.SplitN(entrance, "-", 2)
		shushuDate["entrance"] = parts[0]
		shushuDate["method"] = parts[1]
	} else {
		shushuDate["entrance"] = entrance
		shushuDate["method"] = entrance
	}
	shushumain.ReportData(shushuDate)
}

// 获取包裹道具信息 By Uid
func (p *Player) getItemByUid(uid int64) *pb.ItemInfo {
	if uid == 0 {
		return nil
	}

	for _, item := range p.bag.GetItemList() {
		if item.Uid == uid {
			return item
		}
	}
	return nil
}

// 获取包裹道具信息 By configId
func (p *Player) getItemByCfgId(cfgId int32) *pb.ItemInfo {
	if cfgId == 0 {
		return nil
	}

	for _, item := range p.bag.GetItemList() {
		if item.CfgId == cfgId {
			return item
		}
	}
	return nil
}

func (p *Player) newDay() {
	// 新的一周刷新试练塔层数
	// 判断是否新的一周（以周一为一周的开始）
	now := time.Now()
	_, week := now.ISOWeek()
	_, lastLoginWeek := p.lastLoginAt.ISOWeek()
	if week != lastLoginWeek {
		// // 新的一周，重置试炼塔层数
		// p.TrialTowerMap.BattleId = 0
		// if p.ElementMap.BattleId > 0 {
		// 	p.ElementMap.NeedSweep = 1
		// }

		if !p.bCreate {
			p.tasks.resetNewweek()
		}
	}
	// 查询他是第几天
	p.nowDay++
	p.syncUserDay()
	// p.newDayAward()
	if !p.bCreate {
		p.tasks.resetNewday()
	}

	p.PVPInfo.PVPDuel.dailyBuyNum = 0
	p.PVPInfo.PVPDuel.dailyMatchTimes = 0
	p.PVPInfo.dailySoul = 0
	if !matchManager.sessionOpen {
		p.PVPInfo.sessionSoul = 0
	}
	p.PVPInfo.PVPMelee.dailyBuyNum = 0
	p.PVPInfo.PVPMelee.dailyMatchTimes = 0

	p.dailyWorldBossTimes = 0
	p.totalWorldBossDamage = 0
	p.dayOnlineTime = 0
	p.wellenOnlineTime = 0
	p.vitalityTime = 0
	p.dailyMapReset()
	p.DailyItemUsed = make(map[int32]int64, 0)
	//p.petGold = 0
	//p.petExp = 0
	p.PetMgr.newDay()
	p.InstanceData.BossMap.ChallengeTime = global.GetCopySetCfg(constant.BossMap)
	p.InstanceData.RelicMap.ChallengeTime = global.GetCopySetCfg(constant.RelicMap)
	p.ElementMap.ChallengeTime = global.GlobalCfg.Element.Daily.InitTime
	p.InstanceData.MonKillRecord.MonKillRecords = make(map[int32]int32)
	p.InstanceData.PushingWall.UsedCount = 0
	p.InstanceData.PushingWall.TicketPurchaseCount = 0

	// 跨天提出挂机副本
	if p.desk != nil && p.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
		p.desk.playerKick(p, constant.GameOverTypeTime)
	}

	p.afkDailyTime = 0
	p.AfkMultiplierID = -1 // 每天重置为-1,代表前端需要弹窗
	p.warDailySoul = 0
	// p.AdTotalToday = 0
	atomic.StoreInt32(&p.adTotalToday, 0)
	p.LotteryInfo.FreeNum = 1
	//有些活动，没有领取奖励，需要发给他
	p.checkAndGiveDailyActivityRewards()
	p.SyncUserAvatar(false)
	//累充奖励也要过期给他
	p.checkAndGiveCumulativeChargeRewards()
	//通用任务活动（31-37）过期奖励也要给他
	p.checkAndGiveGradeActivityRewards()
}

// 有些活动，没有领取奖励，需要发给他 checkAndGiveDailyActivityRewards
func (p *Player) checkAndGiveDailyActivityRewards() {
	whileActivity := []int32{
		constant.ActivityTypeLottery,
		constant.ActivityTypeRefine,
		constant.ActivityTypeInstance,
	}

	// 遍历活动配置
	for _, activityType := range whileActivity {
		// 获取活动配置
		actOne := activityType*1000 + 101
		actCfg, err := global.GetActivityCfg(actOne)
		if err != nil {
			continue
		}
		key := "send" + strconv.Itoa(int(activityType))
		if p.createDay > actCfg.ResetTime && p.getUpRecord(key) < 1 {
			// 如果活动已经重置且奖励未领取，则发放奖励
			p.giveActivityReward(activityType, actCfg.Acttag)

			p.addUpRecord(key)
		}
	}
}

// 查询
func (p *Player) giveActivityReward(actType int32, ActivityId int32) {
	//查询用户目领取情况
	resultd, _ := db.ActivityListDay(
		p.uid,
		actType,
		ActivityId,
	)
	roleReceiveInfo := make(map[int32]int)
	for _, item := range resultd {
		roleReceiveInfo[item.Day] = item.Buy
	}
	// 检查是否为VIP（这里需要根据实际情况判断）
	request2 := &protocol.StoreListRequest{}
	request2.Uid = p.uid
	switch actType {
	case constant.ActivityTypeGrowthFund:
		request2.StoreType = constant.ShopTypeGrowthFund
		request2.GoodsUuid = constant.ShopTypeGrowthFund*1000 + ActivityId
	case constant.ActivityTypeLottery:
		request2.StoreType = constant.ShopTypeLotteryFund
		request2.GoodsUuid = constant.ShopTypeLotteryFund*1000 + ActivityId
	case constant.ActivityTypeRefine:
		request2.StoreType = constant.ShopTypeRefineFund
		request2.GoodsUuid = constant.ShopTypeRefineFund*1000 + ActivityId
	case constant.ActivityTypeInstance:
		request2.StoreType = constant.ShopTypeInstanceFund
		request2.GoodsUuid = constant.ShopTypeInstanceFund*1000 + ActivityId
	default:
	}
	delstoren := []int64{
		int64(request2.GoodsUuid),
	}
	//request.StoreType = r.StoreType
	//request.GoodsId = r.GoodsId
	day := 365 * 50
	// 查询用户购买情况
	resultb, _ := db.StoreList(
		p.uid,
		request2.StoreType,
		int32(day),
		delstoren,
		false,
	)
	isVip := false
	if len(resultb) > 0 {
		isVip = true
	}
	// 获取活动配置
	activityCfg := global.ActivityCfg
	result := []protocol.ActivityDetailDay{}
	sendAward := make([]protocol.Attachment, 0)
	// 活动内容
	actvityName := ""
	for _, activityn := range activityCfg {
		if activityn.Acttag == ActivityId && activityn.Type == actType {
			actvityName = activityn.Name
			item := protocol.ActivityDetailDay{
				UUID:       activityn.UUID,
				Content:    activityn.Content,
				CommonItem: make([]protocol.SourceInfo, len(activityn.CommonItem)),
				PayItem:    make([]protocol.SourceInfo, len(activityn.PayItem)),
				PreLimit:   activityn.PreLimit,
			}

			for j, attachment2 := range activityn.CommonItem {
				item.CommonItem[j] = protocol.SourceInfo{
					Id:  attachment2.Key,
					Num: int64(attachment2.Value),
				}
			}

			for j, attachment3 := range activityn.PayItem {
				item.PayItem[j] = protocol.SourceInfo{
					Id:  attachment3.Key,
					Num: int64(attachment3.Value),
				}
			}
			result = append(result, item)
		}
	}

	//状态 1-待领取，2-已领取，3-不可领取，4-需要补签
	for _, item := range result {
		//需要领取的日期
		needDay := item.UUID % 100

		//不同活动，判断条件不同
		switch actType {
		case constant.ActivityTypeGrowthFund:
			if p.Level() < item.PreLimit {
				continue
			}
		case constant.ActivityTypeLottery:
			// 处理抽奖活动
			if int32(p.getUpRecord(constant.LotteryNum)) < item.PreLimit {
				continue
			}
		case constant.ActivityTypeRefine:
			// 处理洗练
			if int32(p.getUpRecord(constant.RefineNum)) < item.PreLimit {
				continue
			}
		case constant.ActivityTypeInstance:
			// 处理副本活动
			if int32(p.getUpRecord(constant.InstanceNum)) < item.PreLimit {
				continue
			}
		}

		if roleReceiveInfo[int32(needDay)] > 0 {
			//领取了且当初也是vip，那就跳过
			if roleReceiveInfo[int32(needDay)] > 1 {
				continue
			}
			//领取了且当初不是vip，补发vip的
			if isVip {
				for _, info := range item.PayItem {
					sendAward = append(sendAward, protocol.Attachment{
						MailId:      "",
						UserPropsId: int64(info.Id),
						Num:         int64(info.Num),
					})
				}

			}
		} else {
			//没领取，那么发，但是要校验是不是vip
			if isVip {
				for _, info := range item.PayItem {
					sendAward = append(sendAward, protocol.Attachment{
						MailId:      "",
						UserPropsId: int64(info.Id),
						Num:         int64(info.Num),
					})
				}

			}
			for _, info := range item.CommonItem {
				sendAward = append(sendAward, protocol.Attachment{
					MailId:      "",
					UserPropsId: int64(info.Id),
					Num:         int64(info.Num),
				})
			}

		}
	}
	if sendAward == nil || len(sendAward) == 0 {
		return
	}
	mailId, _, err := db.MailAdd(
		-1,
		p.uid,
		"系统邮件："+actvityName+"未领取",
		actvityName,
		sendAward,
		constant.ItemAddTask,
	)
	if err != nil {
		logger.Error(err)
	}
	upMailShushuDate(mailId, sendAward)
}
func (p *Player) checkAndGiveCumulativeChargeRewards() {
	actType := constant.ActivityTypeCumulativeCharge
	activityCfg := global.ActivityCfg

	// 按活动标记聚合累充配置
	type cumulativeGroup struct {
		resetTime int32
		configs   []*global.ActivityConfig
	}
	groups := make(map[int32]*cumulativeGroup)

	for _, cfg := range activityCfg {
		if cfg.Type != actType {
			continue
		}

		group, ok := groups[cfg.Acttag]
		if !ok {
			group = &cumulativeGroup{resetTime: cfg.ResetTime}
			groups[cfg.Acttag] = group
		}

		if group.resetTime == 0 && cfg.ResetTime > 0 {
			group.resetTime = cfg.ResetTime
		}

		group.configs = append(group.configs, cfg)
	}

	if len(groups) == 0 {
		return
	}

	totalCharge := p.getUpRecord(constant.CumulativeChargeAmount)
	const gradeHistoryDays = int32(365 * 50)

	for actTag, group := range groups {
		if len(group.configs) == 0 {
			continue
		}

		if group.resetTime <= 0 {
			// 永久活动，无需补发
			continue
		}

		if p.createDay <= group.resetTime {
			// 活动尚未结束
			continue
		}

		progressKey := fmt.Sprintf("cumulative_reward_sent_%d", actTag)
		if p.getUpRecord(progressKey) > 0 {
			continue
		}

		goodsIDs := make([]int64, 0, len(group.configs))
		for _, cfg := range group.configs {
			goodsIDs = append(goodsIDs, cfg.UUID)
		}

		if len(goodsIDs) == 0 {
			p.addUpRecord(progressKey)
			continue
		}

		claimedList, err := db.GradeList(p.uid, gradeHistoryDays, goodsIDs)
		if err != nil {
			logger.Errorf("checkAndGiveCumulativeChargeRewards: query grade list failed, uid=%d, actTag=%d, err=%v", p.uid, actTag, err)
			continue
		}

		claimed := make(map[int64]struct{}, len(claimedList))
		for _, item := range claimedList {
			claimed[item.GoodsUuiD] = struct{}{}
		}

		rewardMap := make(map[int64]int64)
		for _, cfg := range group.configs {
			if _, ok := claimed[cfg.UUID]; ok {
				continue
			}

			if int64(cfg.PreLimit) > totalCharge {
				continue
			}

			for _, item := range cfg.CommonItem {
				rewardMap[int64(item.Key)] += int64(item.Value)
			}
		}

		if len(rewardMap) == 0 {
			p.addUpRecord(progressKey)
			continue
		}

		itemIDs := make([]int64, 0, len(rewardMap))
		for id := range rewardMap {
			itemIDs = append(itemIDs, id)
		}
		sort.Slice(itemIDs, func(i, j int) bool { return itemIDs[i] < itemIDs[j] })

		attachments := make([]protocol.Attachment, 0, len(itemIDs))
		for _, id := range itemIDs {
			attachments = append(attachments, protocol.Attachment{
				UserPropsId: id,
				Num:         rewardMap[id],
			})
		}

		mailTitle := "系统邮件：累充返利补发"
		mailDesc := "累充返利活动已结束，未领取的奖励已通过附件补发，请及时查收。"
		mailID, _, err := db.MailAdd(-1, p.uid, mailTitle, mailDesc, attachments, constant.ItemAddActivityLevel)
		if err != nil {
			logger.Errorf("checkAndGiveCumulativeChargeRewards: send mail failed, uid=%d, actTag=%d, err=%v", p.uid, actTag, err)
			continue
		}

		upMailShushuDate(mailID, attachments)
		p.addUpRecord(progressKey)
	}
}

// checkAndGiveGradeActivityRewards 通用任务活动（31-37）过期奖励补发
func (p *Player) checkAndGiveGradeActivityRewards() {
	// 新的通用任务活动类型
	activityTypes := []int32{
		constant.ActivityTypePinkEntry,
		constant.ActivityTypeRedPartnerStar,
		constant.ActivityTypeEquipmentStar,
		constant.ActivityTypePetSynthesis,
		constant.ActivityTypeSoulStoneLevel,
		constant.ActivityTypeFiveElementLevel,
	}

	now := time.Now().Unix()
	const gradeHistoryDays = int32(365 * 50)

	// 第一步：收集所有需要检查的活动配置和 goodsIDs
	type activityGroupInfo struct {
		actType      int32
		actTag       int32
		groupConfigs []*global.ActivityConfig
		currentValue int64
		progressKey  string
	}

	var allGroups []activityGroupInfo
	allGoodsIDs := make([]int64, 0)

	for _, actType := range activityTypes {
		configs, ok := global.ActivityCfgByType[actType]
		if !ok || len(configs) == 0 {
			continue
		}

		// 获取玩家当前值
		var currentValue int64 = 0
		switch actType {
		case constant.ActivityTypePinkEntry:
			currentValue = p.calcPinkEntryNum()
		case constant.ActivityTypeRedPartnerStar:
			currentValue = p.calcRedPartnerStarNum()
		case constant.ActivityTypeEquipmentStar:
			currentValue = p.calcEquipmentStarNum()
		case constant.ActivityTypePetSynthesis:
			currentValue = p.calcPetSynthesisNum()
		case constant.ActivityTypeSoulStoneLevel:
			currentValue = p.calcSoulStoneLevelNum()
		case constant.ActivityTypeFiveElementLevel:
			currentValue = p.calcFiveElementLevelNum()
		}

		// 按 Acttag 分组处理
		groups := make(map[int32][]*global.ActivityConfig)
		for _, cfg := range configs {
			groups[cfg.Acttag] = append(groups[cfg.Acttag], cfg)
		}

		for actTag, groupConfigs := range groups {
			if len(groupConfigs) == 0 {
				continue
			}

			// 检查是否已经处理过
			progressKey := fmt.Sprintf("grade_activity_reward_sent_%d_%d", actType, actTag)
			if p.getUpRecord(progressKey) > 0 {
				continue
			}

			// 检查活动是否过期
			// 使用第一个配置的 ResetTime 作为参考
			firstCfg := groupConfigs[0]
			if firstCfg.ResetTime <= 0 {
				// 永久活动，无需补发
				continue
			}

			beginTime := p.calculateActivityBeginTime(firstCfg.StartTime)
			if beginTime == 0 {
				continue
			}

			// 计算活动结束时间
			originalTime := time.Unix(beginTime, 0)
			endTime := time.Date(originalTime.Year(), originalTime.Month(), originalTime.Day(), 0, 0, 0, 0, originalTime.Location()).
				AddDate(0, 0, int(firstCfg.ResetTime))

			// 如果活动还未结束，跳过
			if now < endTime.Unix() {
				continue
			}

			// 收集所有 UUID
			for _, cfg := range groupConfigs {
				allGoodsIDs = append(allGoodsIDs, cfg.UUID)
			}

			allGroups = append(allGroups, activityGroupInfo{
				actType:      actType,
				actTag:       actTag,
				groupConfigs: groupConfigs,
				currentValue: currentValue,
				progressKey:  progressKey,
			})
		}
	}

	if len(allGroups) == 0 {
		return
	}

	// 第二步：一次性查询所有已领取的奖励
	claimed := make(map[int64]struct{})
	if len(allGoodsIDs) > 0 {
		claimedList, err := db.GradeList(p.uid, gradeHistoryDays, allGoodsIDs)
		if err != nil {
			logger.Errorf("checkAndGiveGradeActivityRewards: query grade list failed, uid=%d, err=%v", p.uid, err)
			return
		}
		for _, item := range claimedList {
			claimed[item.GoodsUuiD] = struct{}{}
		}
	}

	// 第三步：处理每个活动组
	for _, groupInfo := range allGroups {
		// 收集未领取且满足条件的奖励
		rewardMap := make(map[int64]int64)
		for _, cfg := range groupInfo.groupConfigs {
			if _, ok := claimed[cfg.UUID]; ok {
				continue
			}

			// 检查是否满足前置条件
			if int64(cfg.PreLimit) > groupInfo.currentValue {
				continue
			}

			// 合并奖励
			for _, item := range cfg.CommonItem {
				rewardMap[int64(item.Key)] += int64(item.Value)
			}
		}

		if len(rewardMap) == 0 {
			p.addUpRecord(groupInfo.progressKey)
			continue
		}

		// 转换为邮件附件格式
		itemIDs := make([]int64, 0, len(rewardMap))
		for id := range rewardMap {
			itemIDs = append(itemIDs, id)
		}
		sort.Slice(itemIDs, func(i, j int) bool { return itemIDs[i] < itemIDs[j] })

		attachments := make([]protocol.Attachment, 0, len(itemIDs))
		for _, id := range itemIDs {
			attachments = append(attachments, protocol.Attachment{
				UserPropsId: id,
				Num:         rewardMap[id],
			})
		}

		// 发送邮件
		// 通过活动类型找到对应的排行榜类型，然后获取名称
		rankType := 0
		for rt, at := range rankTypeActivityMap {
			if at == groupInfo.actType {
				rankType = rt
				break
			}
		}
		actTypeName := "活动"
		if rankType > 0 {
			if name, ok := constant.RankTypeName[rankType]; ok {
				actTypeName = name
			}
		}
		mailTitle := fmt.Sprintf("系统邮件：%s活动奖励补发", actTypeName)
		mailDesc := fmt.Sprintf("%s活动已结束，未领取的奖励已通过附件补发，请及时查收。", actTypeName)
		mailID, _, err := db.MailAdd(-1, p.uid, mailTitle, mailDesc, attachments, constant.ItemAddActivityLevel)
		if err != nil {
			logger.Errorf("checkAndGiveGradeActivityRewards: send mail failed, uid=%d, actType=%d, actTag=%d, err=%v", p.uid, groupInfo.actType, groupInfo.actTag, err)
			continue
		}

		upMailShushuDate(mailID, attachments)
		p.addUpRecord(groupInfo.progressKey)
	}
}

func (p *Player) dailyMapReset() {
	// 重置每日副本挑战次数
	p.dailyMaps = make(map[int32]int32, 8)

	// 计算VIP额外次数
	vTimes := make(map[int32]int32, 8)
	if vipCfg := global.GetVipCfg(p.vip); vipCfg != nil {
		for _, v := range vipCfg.DailyTime {
			vTimes[v.Key] = v.Value
		}
	}

	// 定义每日副本类型
	maps := []int32{
		constant.ExperienceFBMap,
		constant.GoldFBMap,
		constant.EnhanceStoneFBMap,
		constant.BookPageFBMap,
		constant.StarUpStoneFBMap,
		constant.SoulStoneFBMap,
	}

	// 重置每个副本的挑战次数
	for _, m := range maps {
		baseTimes := int32(0)
		if cfg := global.GetDailyMapTime(m); cfg != nil {
			baseTimes = int32(cfg.Num)
		}
		p.dailyMaps[m] = baseTimes + vTimes[m]
	}
}

func (p *Player) newDayAward() {
	// 新的一天发奖励
	Cfg := global.GetBonusCfg(p.createDay)
	if Cfg == nil {
		return
	}
	result := make([]protocol.Attachment, 0)
	for _, v := range Cfg {
		result = append(result, protocol.Attachment{
			UserPropsId: int64(v.ID),
			Num:         int64(v.Num),
		})
	}
	mailData := protocol.MailAdd{
		Uid:          -1,
		RecipientUid: p.uid,
		Desc:         "欢迎来到《这岛有剧本》，附件是今日道具福利，请查收。本次测试开服前6天都会有道具福利哦，请及时上线领取",
		Title:        "系统邮件:6天狂欢礼包:第" + strconv.Itoa(int(p.createDay)) + "天",
		ChannelCode:  "6天狂欢",
		Attachments:  result,
	}
	shushumain.ReportMailData(mailData)
}

func (p *Player) update() {
	p.lock.Lock()
	defer p.lock.Unlock()
	nowTime := time.Now().Unix()
	// 不在线，不回复
	if p.curState > PlayerStatusDisconnect {
		return
	}
	// 体力恢复检测
	if p.nextEnergyRecoveryTime > 0 {
		p.nextEnergyRecoveryTime--
	}
	// 体力恢复速度加成,换算成秒
	timeinterval := int(float64(global.GlobalCfg.BaseValue.Energy.Timeinterval*60) * (1 - float64(p.stat(pb.Attribute_ENERGY_RECOVERY_SPEED_PERCENTAGE))/BaseRate))
	if p.energy < p.maxEnergy() {
		if time.Now().After(p.energyRecoveryTime.Add(time.Duration(timeinterval) * time.Second)) {
			p.energy++
			if p.energy > int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) {
				p.energy = int64(global.GlobalCfg.BaseValue.Energy.Maxvalue)
				p.nextEnergyRecoveryTime = 0
			} else {
				p.nextEnergyRecoveryTime = timeinterval
			}
			p.energyRecoveryTime = time.Now()
			p.syncUserBase()
		}
	}

	if p.guildId > 0 {
		// 祈愿点恢复检测
		if p.nextWishRecoveryTime > 0 {
			p.nextWishRecoveryTime--
		}
		// 祈愿点恢复速度加成,换算成秒
		wishTimeInterval := int(global.GlobalCfg.Guild.WishValue.RefreshTime * 60)
		if p.checkTakeSource(constant.ItemTypeGuildWish, int64(global.GlobalCfg.Guild.WishValue.MaxValue)) {
			p.nextWishRecoveryTime = wishTimeInterval
			p.wishRecoveryTime = time.Now()
		} else {
			if p.wishRecoveryTime.IsZero() {
				p.nextWishRecoveryTime = wishTimeInterval
				p.wishRecoveryTime = time.Now()
				p.syncUserBase()
			} else {
				if time.Now().After(p.wishRecoveryTime.Add(time.Duration(wishTimeInterval) * time.Second)) {
					go p.addSource(constant.ItemTypeGuildWish, 1, constant.ItemAddGuildWish)
					p.nextWishRecoveryTime = wishTimeInterval
					p.wishRecoveryTime = time.Now()
					p.syncUserBase()
				}
			}
		}

	}

	p.PetMgr.update()
	//// 宠物恢复检测
	//if p.nextPetRecoveryTime > 0 {
	//	p.nextPetRecoveryTime--
	//}
	//
	//// 宠物恢复速度加成,换算成秒//体力恢复速度加成,换算成秒
	//pettimeinterval := int(global.GlobalCfg.BaseValue.Petfree.Recovertime * 60)
	//// 查询背包里面的宠物蛋数目,如果大于设置的数目就不回复了
	//if p.checkTakeSource(constant.ItemTypePetNormal, int64(global.GlobalCfg.BaseValue.Petfree.Maxnum)) {
	//	p.nextPetRecoveryTime = pettimeinterval
	//	p.petRecoveryTime = time.Now()
	//} else {
	//	if time.Now().After(p.petRecoveryTime.Add(time.Duration(pettimeinterval) * time.Second)) {
	//		go p.addSource(constant.ItemTypePetNormal, 1, constant.ItemAddPetAuth)
	//		p.nextPetRecoveryTime = pettimeinterval
	//		p.petRecoveryTime = time.Now()
	//		p.SyncUserPet(-1)
	//	}
	//}
	//
	//if len(p.pet) > 0 {
	//	if p.pet[0].PetCfgId != 0 {
	//		petInfo := global.GetPetCfg(p.pet[0].PetCfgId)
	//		if petInfo != nil && len(petInfo.OfflineAward) > 0 {
	//			petGoldNum := int32(0)
	//			petGoldGap := int32(0)
	//			petGoldMax := int32(0)
	//			petExpNum := int32(0)
	//			petExpGap := int32(0)
	//			petExpMax := int32(0)
	//			for _, v := range petInfo.OfflineAward {
	//				if v.Key == p.pet[0].Star {
	//					for _, v2 := range v.Values {
	//						if v2.ID == constant.ItemTypeCoin {
	//							petGoldNum = v2.Num
	//							petGoldGap = int32(v2.Gap)
	//							petGoldMax = int32(v2.Max)
	//						}
	//						if v2.ID == constant.ItemTypeExp {
	//							petExpNum = v2.Num
	//							petExpGap = int32(v2.Gap)
	//							petExpMax = int32(v2.Max)
	//						}
	//					}
	//				}
	//			}
	//			// log.Debug("", p.petGold, "-----------", p.petRecoveryTimeGold, p.petExp, "-----------", p.petRecoveryTimeExp)
	//			// log.Debug("", p.pet[0].Star, p.petExp, "-----------", petGoldNum, "---", petGoldGap, "---", petGoldMax, "---", petExpNum, "---", petExpGap, "---", petExpMax)
	//			if petGoldNum > 0 {
	//				// 宠物金币恢复检测
	//				if p.nextPetRecoveryTimeGold > 0 {
	//					p.nextPetRecoveryTimeGold--
	//				}
	//				// 宠物金币速度加成,换算成秒
	//				timeintervalgold := int(float64(petGoldGap))
	//				if p.petGold < petGoldMax {
	//					if time.Now().After(p.petRecoveryTimeGold.Add(time.Duration(timeintervalgold) * time.Second)) {
	//						if p.petGold+petGoldNum > petGoldMax {
	//							p.petGold = petGoldMax
	//							p.nextPetRecoveryTimeGold = 0
	//							if petGoldMax-p.petGold > 0 {
	//								go p.addSource(constant.ItemTypeCoin, int64(petGoldMax-p.petGold), constant.ItemAddPetGold)
	//							}
	//						} else {
	//							p.petGold += petGoldNum
	//							p.nextPetRecoveryTimeGold = timeintervalgold
	//							if petGoldNum > 0 {
	//								go p.addSource(constant.ItemTypeCoin, int64(petGoldNum), constant.ItemAddPetGold)
	//							}
	//						}
	//						p.petRecoveryTimeGold = time.Now()
	//						p.syncUserBase()
	//					}
	//				}
	//			}
	//
	//			if petExpNum > 0 {
	//				// 宠物经验恢复检测
	//				if p.nextPetRecoveryTimeExp > 0 {
	//					p.nextPetRecoveryTimeExp--
	//				}
	//				// 宠物经验速度加成,换算成秒
	//				timeintervalexp := int(float64(petExpGap))
	//				if p.petExp < petExpMax {
	//					if time.Now().After(p.petRecoveryTimeExp.Add(time.Duration(timeintervalexp) * time.Second)) {
	//						if p.petExp+petExpNum > petExpMax {
	//							p.petExp = petExpMax
	//							p.nextPetRecoveryTimeExp = 0
	//							if petExpMax-p.petExp > 0 {
	//								go p.addSource(constant.ItemTypeExp, int64(petExpMax-p.petExp), constant.ItemAddPetExp)
	//							}
	//						} else {
	//							p.petExp += petExpNum
	//							p.nextPetRecoveryTimeExp = timeintervalexp
	//							if petExpNum > 0 {
	//								go p.addSource(constant.ItemTypeExp, int64(petExpNum), constant.ItemAddPetExp)
	//							}
	//						}
	//						p.petRecoveryTimeExp = time.Now()
	//						p.syncUserBase()
	//					}
	//				}
	//			}
	//
	//		}
	//
	//	}
	//}

	p.offlineTime = nowTime

	// 在线时长累加
	p.dayOnlineTime++
	p.totalOnlineTime++
	p.wellenOnlineTime++

	//2025年7月25日17:27:53 活力满满因为需求变更，暂时隐藏
	// 超过阈值就没有奖励了
	// 查询vip对应的活力分数
	// vitality := int64(0)
	// if p.vip > 0 {
	// 	if vipCfg := global.GetVipCfg(p.vip); vipCfg != nil {
	// 		vitality += int64(vipCfg.Energyadd) * 60
	// 	}
	// }
	// // 查询特权卡对应的活力分数
	// if len(p.privilege) > 0 {
	// 	for _, v := range p.privilege {
	// 		// 查询特权卡对应的活力分数
	// 		if privilegeCfg := global.GetPrivilegeCfg(v.Pid); privilegeCfg != nil {
	// 			vitality += int64(privilegeCfg.Energystate) * 60
	// 		}
	// 	}
	// }
	// // go p.openbox(3)
	// // todo vip

	// if p.vitalityTime < int64(global.GlobalCfg.BaseValue.Onlinebattle.Basetime)+vitality {
	// 	p.vitalityTime++
	// 	// 每隔一段时间发奖励
	// 	if p.vitalityTime%int64(global.GlobalCfg.BaseValue.Onlinebattle.Triggertime) == 0 {
	// 		// 随机
	// 		randomType, id, randomReward := p.randomReward(p.level)
	// 		packet := &pb.HomeAutoEncounterPush{
	// 			Type:     randomType,
	// 			Id:       id,
	// 			ItemInfo: make([]*pb.ItemInfo, 0),
	// 		}

	// 		if len(randomReward) > 0 {
	// 			var reason string
	// 			if randomType == constant.MaintriggerMonster {
	// 				reason = constant.ItemAddAFKMonster
	// 			}
	// 			if randomType == constant.MaintriggerNormal {
	// 				reason = constant.ItemAddAFKBox
	// 			}
	// 			for _, v := range randomReward {
	// 				packet.ItemInfo = append(packet.ItemInfo, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	// 			}

	// 			// 延迟添加
	// 			tmpReward := make([]global.IdNum, 0)
	// 			for _, v := range randomReward {
	// 				getItemCfg, ok := global.ItemCfg[v.ID]
	// 				if !ok {
	// 					continue
	// 				}
	// 				// 如果是使用道具
	// 				if getItemCfg.UserNow > 0 {
	// 					res, err := p.itemAutoUse(v.ID, v.Num)
	// 					if err != nil {
	// 						tmpReward = append(tmpReward, global.IdNum{ID: v.ID, Num: int32(v.Num)})
	// 						// 使用失败
	// 						continue
	// 					}
	// 					for _, v2 := range res {
	// 						tmpReward = append(tmpReward, global.IdNum{ID: v2.ID, Num: int32(v2.Num)})
	// 					}
	// 				} else {
	// 					tmpReward = append(tmpReward, global.IdNum{ID: v.ID, Num: int32(v.Num)})
	// 				}
	// 			}

	// 			var delayTime time.Duration
	// 			delayTime = time.Millisecond * 500
	// 			if reason == constant.ItemAddAFKMonster {
	// 				delayTime = time.Duration(math.Floor(float64(5)/float64(p.getCombatPartnerNum()))*1.2) * time.Second
	// 			}
	// 			// 告知前端的奖励
	// 			err := p.checkAddSourceBatch(tmpReward)
	// 			if err == nil {
	// 				p.pushMsg(protocol.RouteHomeAutoEncounterPush, packet)
	// 			}

	// 			time.AfterFunc(time.Duration(delayTime), func() {
	// 				// 检查背包有没有满
	// 				if err == nil {
	// 					p.addSourceBatch(tmpReward, reason)
	// 				} else {
	// 					attachment := []protocol.Attachment{}
	// 					for _, v := range tmpReward {
	// 						attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
	// 					}
	// 					p.SendMail(constant.ItemAddDailymap, attachment)
	// 				}
	// 			})
	// 		}
	// 	}
	// }
	// 时间任务检测
	//p.tasks.dispatch(pb.TargetType_TIME_RANGE, 0, 0)
	// 特权卡
	newPrivilege := []protocol.PrivilegeDetailTCP{}
	for _, v := range p.privilege {
		if time.Now().After(v.EndTime) {
			continue
		}
		newPrivilege = append(newPrivilege, v)
	}
	if len(p.privilege) != len(newPrivilege) {
		p.privilege = newPrivilege
	}

	// 游戏内玩家,buff由房间goroutine处理
	if p.desk == nil {
		p.processBuff()
		p.processHp()
	}
}

func (p *Player) itemAutoUse(CfgId int32, useNum int32) ([]global.IdNum, error) {
	reward := []global.IdNum{}
	if itemUseCfg, ok := global.ItemUseCfg[CfgId]; ok {
		// 检查使用道具的上限
		for _, v := range itemUseCfg.GiveGoods {
			if v.ID == constant.ItemTypeEnergy {
				if p.energy >= int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) {
					return nil, errutil.ErrEnergyFull
				} else {
					needEnergy := int64(global.GlobalCfg.BaseValue.Energy.Maxvalue) - p.energy
					needNum := int32(math.Ceil(float64(needEnergy) / float64(v.Num)))
					if useNum > needNum {
						useNum = needNum
					}
					break
				}
			}
		}

		jobreward := []global.IdNumJobWeight{}
		jobrewardweight := []global.IdNumJobWeight{}
		if len(itemUseCfg.GiveGoods) > 0 {
			for _, gg := range itemUseCfg.GiveGoods {
				if p.job == constant.MowJob(gg.Job) || gg.Job == 0 {
					jobreward = append(jobreward, gg)
				}
			}
		}
		var weight int32
		for _, jr := range jobreward {
			weight += jr.Weight
			jobrewardweight = append(jobrewardweight, global.IdNumJobWeight{ID: jr.ID, Num: jr.Num, Job: jr.Job, Weight: int32(weight)})
		}
		for i := 0; i < int(useNum); i++ {
			tmpReward := []global.IdNumJobWeight{}
			if weight == 0 {
				tmpReward = jobrewardweight
			} else {
				w := global.RandIn(int(weight))
				for _, jrw := range jobrewardweight {
					if w < int(jrw.Weight) {
						tmpReward = append(tmpReward, global.IdNumJobWeight{ID: jrw.ID, Num: jrw.Num, Job: jrw.Job, Weight: weight})
						break
					}
				}
			}

			for _, tr := range tmpReward {
				bFind := false
				for k, r := range reward {
					if r.ID == tr.ID {
						if rewardCfg, err := global.GetItemCfg(r.ID); err == nil {
							if rewardCfg.Maxnum > 1 {
								reward[k].Num += tr.Num
								bFind = true
							}
						}
						break
					}
				}
				if !bFind {
					reward = append(reward, global.IdNum{ID: tr.ID, Num: tr.Num})
				}
			}

		}
	}
	return reward, nil
}

// 给到随机的奖励
func (p *Player) randomReward(level int32) (returnType int32, id int32, rewards []global.IdNum) {
	createCfg := global.MaintriggerCfg
	if len(createCfg) == 0 {
		return constant.MaintriggerType, 0, []global.IdNum{}
	}
	for _, v := range createCfg {
		for _, v1 := range v.Level {
			if level >= v1.MinLevel && level <= v1.MaxLevel {
				var rate int32
				for _, v2 := range v.MonsterIds {
					rate += v2.Rate
				}
				weight := global.RandIn(int(rate))
				DemoRandom := global.MonsterIds{}
				beginPriority := 0
				for _, v3 := range v.MonsterIds {
					if beginPriority+int(v3.Rate) >= weight {
						DemoRandom = v3
						break
					}
					beginPriority += int(v3.Rate)
				}
				// 不是怪物就直接给他
				if DemoRandom.Type == constant.MaintriggerNormal {
					return constant.MaintriggerNormal, DemoRandom.Id, []global.IdNum{{ID: DemoRandom.Id, Num: DemoRandom.Num}}
				}
				if DemoRandom.Type == constant.MaintriggerMonster {
					// 是怪物就给他怪物的奖励
					// log.Debug("怪物", DemoRandom.Id, DemoRandom.Num)
					dropInfo := monsterMgr.MonsterDorp(DemoRandom.Id, DemoRandom.Num)
					// log.Debug("怪物", dropInfo)
					return constant.MaintriggerMonster, DemoRandom.Id, dropInfo
				}
				return constant.MaintriggerType, 0, []global.IdNum{}
			} else {
				continue
			}
		}
	}
	return constant.MaintriggerType, 0, []global.IdNum{}
}

func (p *Player) upBoxShushuDate(boxinfo BoxInfo) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "box_level_up"
	shushuDate["box_level_before"] = boxinfo.NowLevel - 1
	shushuDate["box_level_after"] = boxinfo.NowLevel
	if boxinfo.IsAccelerate > 0 {
		shushuDate["is_accelerate"] = true
	} else {
		shushuDate["is_accelerate"] = false
	}
	shushuDate["accelerate_num"] = boxinfo.AccelerateNum
	shushuDate["cost_gold"] = boxinfo.CostGold
	shushuDate["time_cost"] = boxinfo.TimeCost
	shushumain.ReportData(shushuDate)

	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_box_level"] = boxinfo.NowLevel
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

// func (p *Player) levelUp() bool {
// 	if p.level >= global.GlobalCfg.BaseValue.Rolelevel.Maxlevel {
// 		return false
// 	}

// 	cfg, err := global.GetRoleCfg(int32(p.job), p.level)
// 	if err != nil {
// 		p.logger.Errorf("升级失败 userId: %d,  err: %s", p.uid, err.Error())
// 		return false
// 	}

// 	var bUpgrade bool
// 	if p.getCombatPartnerNum() == 1 {
// 		for p.exp >= cfg.NeedExp {
// 			p.exp -= cfg.NeedExp
// 			if p.exp <= 0 {
// 				p.exp = 0
// 			}

// 			cfg, err = global.GetRoleCfg(int32(p.job), p.level+1)
// 			if err != nil {
// 				return false
// 			}
// 			p.level++
// 			p.levelGift()
// 			bUpgrade = true
// 		}
// 	} else {
// 		if p.exp < cfg.NeedExp {
// 			return false
// 		}
// 		p.exp -= cfg.NeedExp
// 		if p.exp <= 0 {
// 			p.exp = 0
// 		}
// 		p.level++
// 		p.levelGift()
// 		bUpgrade = true
// 	}
// 	oldtotalCombat := p.totalCombat
// 	if bUpgrade {

// 		// p.logger.Debugf("升级 userId: %d,  level: %d  vip: %d  exp: %d", p.uid, p.level, p.vip, p.exp)
// 		p.refreshStats()
// 		newtotalCombat := p.totalCombat
// 		p.upMyLevelShushuDate(p.level, newtotalCombat-oldtotalCombat)
// 		p.currentHP = p.GetMaxHP()
// 		p.stats.setStats(pb.Attribute_HEALTH, p.currentHP)
// 		p.skillAutoLearn()

// 		p.pushMsg(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: p.level, Uid: p.uid})
// 		if p.desk != nil {
// 			p.desk.Broadcast(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: p.level, Uid: p.uid})
// 		}
// 		// 更新排行与任务
// 		p.rankChangeTCP(constant.RankTypeLevel, int64(p.level))
// 		p.tasks.dispatch(pb.TargetType_ROLE_LEVEL_REQUIREMENT, 0, p.level)

// 		num := int32(0)
// 		for _, v := range global.GlobalCfg.Openslot.IDLevel {
// 			if p.level >= v.Needlevel {
// 				num++
// 			}
// 		}
// 		p.tasks.dispatch(pb.TargetType_UNLOCK_PARTNER_SLOT, 0, num)

// 		p.tasks.check()
// 	}
// 	return true
// }

func (p *Player) levelUp() bool {
	if p.level >= global.GlobalCfg.BaseValue.Rolelevel.Maxlevel {
		return false
	}

	cfg, err := global.GetRoleCfg(int32(p.job), p.level)
	if err != nil {
		p.logger.Errorf("升级失败 userId: %d,  err: %s", p.uid, err.Error())
		return false
	}

	var bUpgrade bool
	for p.exp >= cfg.NeedExp {
		p.exp -= cfg.NeedExp
		if p.exp <= 0 {
			p.exp = 0
		}

		cfg, err = global.GetRoleCfg(int32(p.job), p.level+1)
		if err != nil {
			return false
		}
		p.level++
		p.levelGift()
		bUpgrade = true
	}
	oldtotalCombat := p.totalCombat
	if bUpgrade {

		for _, par := range p.partners {
			par.setLevel(p.level)
		}

		// p.logger.Debugf("升级 userId: %d,  level: %d  vip: %d  exp: %d", p.uid, p.level, p.vip, p.exp)
		p.refreshStats()
		newtotalCombat := p.totalCombat
		p.upMyLevelShushuDate(p.level, newtotalCombat-oldtotalCombat)
		p.currentHP = p.GetMaxHP()
		p.stats.setStats(pb.Attribute_HEALTH, p.currentHP)
		p.skillAutoLearn()

		p.pushMsg(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: p.level, Uid: p.uid})
		if p.desk != nil {
			p.desk.Broadcast(protocol.RouteLevelUpPush, &pb.LevelUpPush{NewLevel: p.level, Uid: p.uid})
		}
		// 更新排行与任务
		p.rankChangeTCP(constant.RankTypeLevel, int64(p.level))
		p.tasks.dispatch(pb.TargetType_ROLE_LEVEL_REQUIREMENT, 0, p.level)

		//废弃，现在是按照主线关卡来
		// num := int32(0)
		// for _, v := range global.GlobalCfg.Openslot.IDLevel {
		// 	if p.level >= v.Needlevel {
		// 		num++
		// 	}
		// }
		// p.tasks.dispatch(pb.TargetType_UNLOCK_PARTNER_SLOT, 0, num)

		p.tasks.check()
		p.artifactUnlock()
	}
	return true
}

func (p *Player) upMyLevelShushuDate(level int32, combat int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	shushuUserSetDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
		// shushuUserSetDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "role_level_upgrade"
	shushuDate["level_before"] = level - 1
	shushuDate["level_after"] = level
	shushuDate["role_job"] = constant.JobName[p.job]
	shushuDate["role_power"] = int32(p.combat)
	shushuDate["role_power_up"] = combat
	shushumain.ReportData(shushuDate)

	shushuUserSetDate["current_level"] = level
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

// vip升级
func (p *Player) vipLevelUp() {
	cfg := global.GetVipCfg(p.vip + 1)
	if cfg == nil {
		p.logger.Errorf("vip升级失败 userId: %d", p.uid)
		return
	}
	if p.vip >= int32(len(global.VIPCfg)) {
		return
	}
	var bUpgrade bool
	delExp := int64(0)
	vip := p.getItemCountInBag(constant.ItemTypeVipExp)
	for vip >= int64(cfg.NeedExp) {
		delExp += int64(cfg.NeedExp)
		vip -= int64(cfg.NeedExp)
		bUpgrade = true
		p.vip++
		// Check if VIP level has reached the maximum
		if p.vip >= int32(len(global.VIPCfg)) {
			break
		}
		// Update the required experience for the next VIP level
		cfg = global.GetVipCfg(p.vip + 1)
		if cfg == nil {
			break
		}
	}
	if bUpgrade {
		p.RefreshStatsAsync()
		// 上报
		p.upVipShushuDate(p.vip)
		p.tasks.dispatch(pb.TargetType_VIP_LEVEL, 0, p.vip)
		p.tasks.check()

		p.addToken(constant.ItemTypeVipExp, -delExp, constant.ItemDelVipExp)
		p.syncUserBase()
	}
}

func (p *Player) upVipShushuDate(vip int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	shushuUserSetDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
		shushuUserSetDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "vip_upgrad"
	shushuDate["vip_before"] = vip - 1
	shushuDate["vip_after"] = vip
	shushumain.ReportData(shushuDate)
	shushuUserSetDate["current_vip"] = vip
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (p *Player) setDesk(d *Desk) {
	if d == nil {
		p.logger.Error("桌号为空")
		return
	}

	p.desk = d
	p.roomId = d.cfgID

	p.logger = log.WithFields(log.Fields{fieldDesk: p.desk.deskID, fieldPlayer: p.Uid()})
}

func (p *Player) bindSession(s *session.Session) {
	p.session = s
	p.session.Set(kCurPlayer, p)
}

func (p *Player) removeSession() {
	if p.session != nil {
		p.session.Remove(kCurPlayer)
		p.session = nil
	}
}

func (p *Player) setPlayerStatus(s PlayerStatus) {
	p.curState = s
}

func (p *Player) isRobot() bool {
	return p.bRobot
}

func (p *Player) pushMsg(route string, v interface{}) error {
	if p.isRobot() {
		return nil
	}

	if p.curState == PlayerStatusOffLine || p.curState == PlayerStatusDisconnect {
		return nil
	}

	if p.session == nil {
		p.logger.Warningf("pushMsg Failed ErrPlayerNotFound ,route:%s ", route)
		return errutil.ErrPlayerNotFound
	}

	if err := p.session.Push(route, v); err != nil {
		p.logger.Errorf("pushMsg Error:%s ,route:%s msg:%v", err.Error(), route, v)
		return err
	}
	return nil
}

func (p *Player) reset() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.desk = nil
	p.target = nil
	p.expirationTime = time.Now()
	p.dead = false
	p.roomId = 0
	p.setPlayerStatus(PlayerStatusIdle)
	p.stats.setStats(pb.Attribute_HEALTH, p.stats.getStats(pb.Attribute_MAX_HEALTH))

	p.currentHP = p.GetMaxHP()
	p.displayHP = p.currentHP
	p.reviveTimes = 0
	p.autoBuyInstanceTime = false

	p.resetBuff()

	for _, v := range p.partnerSlots {
		if v == nil {
			continue
		}
		if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
			partner.reset()
		}
	}

	p.refreshStats(WithRefreshStatsPushMessage(false))

	p.logger = log.WithField(fieldPlayer, p.uid)
}

// 学习技能
func (p *Player) skillAutoLearn() {
	for k, v := range global.SkillCfg {
		if v.SkillProfession != int32(p.job) {
			continue
		}

		if s := p.getSkill(k); s != nil {
			continue
		}

		bLearn := false
		for _, skill := range p.skills {
			if skill.CfgId/100 == v.SkillID {
				bLearn = true
				break
			}
		}
		if bLearn {
			continue
		}

		if v.RoleLevelLimit <= p.level && v.Level <= 1 {
			p.skills = append(p.skills, &PlayerSkill{CfgId: v.UID, Cooldown: time.Now(), AdvancedLevel: 0})
			if v.Level > 0 && v.ActiveorPassive != PassiveSkill {
				for i, v2 := range p.skillsSlot { // 主动技能学习后自动装备到技能槽位
					if v2 == 0 {
						// 如果是特殊技能，需要放到第一个槽位
						if v.SType == constant.SkillTypeCareerSpecial {
							p.skillsSlot[constant.SkillSlotType] = v.UID
						} else {
							if i == int(constant.SkillSlotType) {
								continue
							}
							p.skillsSlot[i] = v.UID
						}
						break
					}
				}
				// 技能
				skillslot := make([]int32, 0)
				for _, v1 := range p.skillsSlot {
					skillslot = append(skillslot, v1)
				}
				p.pushMsg(protocol.RouteNewSkill, &pb.NewSkillPush{Skill: &pb.SkillInfo{CfgId: v.UID, Cooldown: 0, Advanced: 0}, SkillsSlot: skillslot})
				p.tasks.dispatch(pb.TargetType_SKILL_UPGRADING_REQUIREMENT, v.UID, v.Level)

				p.refreshStats()
			}
		}
	}
}

func (p *Player) addSkill(skillId int32) {
	if s := p.getSkill(skillId); s != nil {
		return
	}

	cfg, err := global.GetSkillCfg(skillId)
	if err != nil {
		return
	}

	if cfg.SType != constant.SkillTypeArtifact && cfg.SkillProfession != int32(p.job) {
		return
	}

	p.skills = append(p.skills, &PlayerSkill{CfgId: cfg.UID, Cooldown: time.Now(), Distance: cfg.MaxDistance})
}

func (p *Player) getSkill(cfgID int32) *PlayerSkill {
	for _, s := range p.skills {
		if s.CfgId == cfgID {
			return s
		}
	}

	return nil
}

func (p *Player) getSkillex(cfgID int32) (*PlayerSkill, int64) {
	for _, s := range p.skills {
		if s.CfgId == cfgID {
			return s, p.uid
		}
	}
	for _, v := range p.partnerSlots {
		if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
			if skill := partner.getSkill(cfgID); skill != nil {
				return skill, partner.uid
			}
		}
	}
	return nil, 0
}

func (p *Player) die(object IMapObject) {
	p.dead = true
	p.deadTime = time.Now()
	p.target = nil

	p.miningBreak()

	p.reviveTime = time.Now().Add(1500 * time.Millisecond)

	packet := pb.ObjectDiePush{ObjectId: p.uid,
		ReviveTimes: int32(p.reviveTimes),
		CanRevive:   true,
		Countdown:   int32(global.GlobalCfg.War.Autorevive.Time),
		ReviveTip:   false}

	if p.desk != nil {
		if p.desk.roomCfg.Roomtype == constant.EnhanceStoneFBMap || p.desk.roomCfg.Roomtype == constant.BookPageFBMap {
			p.reviveTime = time.Now().Add(time.Duration(global.GlobalCfg.War.Autorevive.Time) * time.Second)
			if p.reviveTimes >= len(global.GlobalCfg.Revive) {
				p.reviveTime = time.Now().Add(time.Duration(global.GlobalCfg.War.Autorevive.Time) * time.Second)
				packet.CanRevive = false
			}
		} else if p.desk.roomCfg.Roomtype == constant.PVEWarMap {
			p.reviveTime = time.Now().Add(time.Duration(global.GlobalCfg.War.Autorevive.Time) * time.Second)
		}
	}

	p.resetBuff()
	p.SetObAction(ObjectActionDead)

	// 没有复活技能，全都得死
	allDie := true
	if skill, obj := p.getSlotSkillByType(constant.SkillTypeReborn); skill != nil && obj != nil {
		if !obj.dead && time.Now().After(skill.Cooldown) {
			allDie = false
		}
	}
	if buff := p.getBuffByType(constant.BuffReviveOnce); buff != nil {
		allDie = false
	}

	if allDie {
		for _, v := range p.partnerSlots {
			if v == nil {
				continue
			}
			if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
				partner.die(object)
			}
		}
	}

	packet.ReviveTip = allDie
	if p.desk != nil {
		p.desk.Broadcast(protocol.RouteObjectDiePush, &packet)

		uid := object.Uid()
		if object.ObjectType() == ObjectTypePartner {
			if object.Owner() != nil {
				uid = object.Owner().Uid()
			}
		}

		if killer, err := p.desk.playerWithId(uid); err == nil && killer != nil && !p.IsRobot() {
			killer.pveWarMapAddSoul(int64(global.GlobalCfg.War.PK.Killperson))
			p.desk.recordKill(uid)

			context := fmt.Sprintf(errutil.WarKilled, killer.Name())
			packet := &pb.ChatPush{Type: pb.ChatType_SYSTEM, Context: context, Uid: p.uid, Name: "系统"}
			p.pushMsg(protocol.RouteChatPush, packet)
		}

		// else {
		// 	logger.Errorf("[%d:%s]被杀， 杀手不是人[%d: %s], 桌子ID %d ", p.uid, p.name, object.Uid(), object.Name(), p.desk.deskID)
		// }
	} else {
		p.logger.Error("玩家死亡 , 桌子已经被清理")
	}
	// p.logger.Debugf("玩家 [%s] 死亡！！！", p.name)
}

// 参数说明：reviveType: 0 表示复活点复活, 1 表示原地复活
func (p *Player) revive(reviveType int32) bool {
	bornLocation := &p.desk.mapCfg.BornPoints[global.RandIn(len(p.desk.mapCfg.BornPoints))]
	if reviveType == 1 {
		bornLocation = &p.cell.location
		p.reviveTimes++
	} else {
		p.reviveTimes = 0
	}
	cell := p.desk.mapInfo.getAvailableCell(bornLocation, 2)
	if cell == nil {
		logger.Warningf("复活失败, 未找到出生点 uid=%d  name=%s mapId=%d", p.uid, p.name, p.desk.mapId)
		return false
	}
	p.desk.mapInfo.locationChange(p, global.ToPointF(&cell.location))
	p.dead = false
	p.currentHP = p.stat(pb.Attribute_MAX_HEALTH)
	p.displayHP = p.currentHP
	p.stats.setStats(pb.Attribute_HEALTH, p.currentHP)
	p.resetBuff()
	p.desk.Broadcast(protocol.RouteRevivePush, &pb.RevivePush{
		Uid:       p.uid,
		Location:  p.location,
		Direction: p.direction,
		CurHp:     p.currentHP,
	})

	for _, v := range p.partnerSlots {
		if v.PartnerUid == 0 {
			continue
		}
		if partner := p.getPartnerInfoBySlotId(v.SlotId); partner != nil {
			partner.revive()
		}
	}

	for _, v := range p.skillsSlot {
		if v == 0 {
			continue
		}
		if skill, _ := p.getSkillex(v); skill != nil {
			if scfg, _ := global.GetSkillCfg(skill.CfgId); scfg != nil && scfg.CoolDown > 10000 {
				p.pushMsg(protocol.RouteSkillCooldownPush, &pb.SkillCooldownPush{CfgId: skill.CfgId, Delay: 0})
			}
		}
	}

	p.AddBuff(int32(constant.GOD_MODE), p, int32(global.GlobalCfg.GVG.Relive.Protecttime*1000), nil, 0, 0)
	// p.logger.Debugf("***********[%s]玩家复活！！！%t   当前血量 %d -- %d --- %d --[%v]", p.name, p.isRobot(), p.currentHP, p.stat(pb.Attribute_HEALTH), p.GetMaxHP(), p.location)

	return true
}

func (p *Player) attack(targets []IMapObject, skillId int32, location *pb.Point, direction pb.Direction) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.desk == nil {
		return
	}
	skill, cast_uid := p.getSkillex(skillId)
	if skill == nil {
		p.logger.Debugf("攻击失败 使用了未学习的技能 userId: %d skillID: %d ", p.uid, skillId)
		return
	}

	castDir := p.direction
	castObjType := ObjectTypePlayer
	if cast_uid != p.uid {
		castObjType = ObjectTypePartner
		if partner := p.getPartnerInfo(cast_uid); partner != nil {
			castDir = partner.direction
		}
	}
	caster := p.MapObject
	if cast_uid != p.uid && p.getPartnerInfo(cast_uid) != nil {
		caster = p.getPartnerInfo(cast_uid).MapObject
	}

	skillCfg, err := global.GetSkillCfg(int32(skillId))
	if err != nil {
		p.logger.Debugf("攻击失败 技能配置读取失败 userId: %d skillID: %d ", p.uid, skillId)
		return
	}

	// 技能CD
	if skillCfg.CoolDown > 0 && time.Now().Before(skill.Cooldown) {
		p.logger.Debugf("攻击失败 技能冷却时间未到 userId: %d  skillID: %d ", p.uid, skillId)
		return
	}

	if time.Now().Before(caster.attackTime) && skillCfg.SType != constant.SkillTypeArtifact {
		p.logger.Debugf("攻击失败 攻击频率过快 userId: %d  skillID: %d ", p.uid, skillId)
		return
	}

	if skillCfg.SType == constant.SkillTypeArtifact && !p.isArtifactSkillValid(skillCfg.UID) {
		p.logger.Debugf("攻击失败 使用了未装备的神兵技能 userId: %d  skillID: %d ", p.uid, skillId)
		return
	}

	roomCfg := global.GetRoomCfg(p.desk.cfgID)
	if roomCfg == nil {
		p.logger.Debugf("攻击失败 房间配置不存在 userId: %d  desk.cfgID: %d ", p.uid, p.desk.cfgID)
		return
	}

	p.miningBreak()

	// 不需要目标的技能清空前端发送的目标
	if skillCfg.Needaim == 0 {
		targets = nil
	} else {
		if len(targets) == 0 {
			p.logger.Debugf("[%s]攻击失败 此技能需要目标 userId: %d  skillID: %d ", p.name, p.uid, skillId)
			return
		} else {
			// 被嘲讽的状态下不能攻击嘲讽者以外的目标
			if debuff := p.getBuffByType(constant.BuffTaunt); debuff != nil {
				var bFind bool
				for _, v := range targets {
					if v.Uid() == debuff.Owner.Uid() {
						bFind = true
					}
				}
				if !bFind {
					p.logger.Debugf(" 被嘲讽的状态下不能攻击嘲讽者以外的目标 ， 嘲讽者ID %d", debuff.Owner.Uid())
					return
				}
			}

			// 技能距离
			for _, t := range targets {
				if skillCfg.MaxDistance > 0 && maputil.Distance(caster.location, t.Location()) > int(skillCfg.MaxDistance) {
					p.logger.Debugf("攻击失败 超出攻击距离 userId: %d, 位置:%v , 目标坐标:%v  skillID: %d ", p.uid, p.location, t.Location(), skillId)
					return
				}
			}

		}
	}

	var targetsUids []int64
	attackCells := []*Cell{}
	var attackLocation *global.PointI
	if targets == nil {
		targets = []IMapObject{}
	}
	if location != nil {
		attackLocation = &global.PointI{X: int(location.X), Y: int(location.Y)}
	}
	if direction >= pb.Direction_UP && direction <= pb.Direction_UP_LEFT {
		p.direction = direction
	}

	Radius := skillCfg.Radius
	Cooldown := int32(0)
	// 技能进阶效果
	if skillUpperlevel, err := global.GetUpperlevelSettingCfg(skillCfg.UID, skill.AdvancedLevel); err == nil {
		Radius += skillUpperlevel.Distance
		Cooldown = skillUpperlevel.Cooldown
	}
	// 技能影响到的全部格子
	switch skillCfg.SkillID {

	case constant.SkillTypeHalfMoon: // 扇形
		attackCells = p.desk.mapInfo.getSectorCells(&caster.cell.location, caster.direction, Radius)

	case constant.SkillTypeFourSides: // 四方向
		for i := 0; i < 4; i++ {
			direction := maputil.ShiftDirection(caster.direction, i*2)
			attackCells = append(attackCells, p.desk.mapInfo.getLineCells(&caster.cell.location, direction, Radius)...)
		}

	case constant.SkillTypeMagicMissile, constant.SkillTypeStellarBurst, constant.SkillTypeSoulTalisman,
		constant.SkillTypeVoidLava, constant.SkillTypeVoidCollapse: // 正方形
		if attackLocation != nil {
			attackCells = p.desk.mapInfo.getSquareCells(attackLocation, Radius)
		}

	case constant.SkillTypeFreezeTime, constant.SkillTypeTimeSpaceLock, constant.SkillTypeDuel, constant.SkillTypeRepelFoe,
		constant.SkillTypeMassHeal, constant.SkillTypeDivineStrength, constant.SkillTypeArcaneMastery, constant.SkillTypeBlessingOfCourage,
		constant.SkillTypeNetherSpell, constant.SkillTypeDivineBattleArmor, constant.SkillTypeTaunt, constant.SkillTypeStomp,
		constant.SkillTypeBleeding, constant.SkillTypeWreck, constant.SkillTypeApocalypse: // 正方形
		attackCells = p.desk.mapInfo.getSquareCells(&caster.cell.location, Radius)

	case constant.SkillTypeBreakArmy, constant.SkillTypeAuroraPulse, constant.SkillTypeEarthquakeWave: // 长方形
		attackCells = p.desk.mapInfo.getLineCells(&caster.cell.location, caster.direction, Radius)

		// 独立处理有个特殊效果的技能
	case constant.SkillTypeSwiftStrike, constant.SkillTypeDecoy, constant.SkillTypeDecoyBoss, constant.SkillTypeInstantaneousMove,
		constant.SkillTypeMassPurification, constant.SkillTypeMassPurificationBoss,
		constant.SkillTypeGoddessDescent, constant.SkillTypeBrambleCage, constant.SkillTypeSwordRain:
		targets = nil
		pdesk := p.desk
		time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
			if pdesk != nil {
				pdesk.addDelayAction(castObjType, cast_uid, skillId, nil, castDir, attackLocation)
			}
		})
	}

	for _, cell := range attackCells {
		if cell == nil || cell.block {
			continue
		}
		for _, obj := range cell.objects {
			bFind := false
			for _, have := range targets {
				if have.Uid() == obj.Uid() {
					bFind = true
					break
				}
			}
			if !bFind {
				if p.CanAttackTarget(obj) && p.SkillEffectType(skillCfg, obj) {
					targets = append(targets, obj)
				}
			}
		}
	}

	for _, t := range targets {
		if t.Uid() == p.uid {
			continue
		}

		if t.ObjectType() < ObjectTypePlayer && t.ObjectType() > ObjectTypeDecoy {
			continue
		}

		if t.ObjectType() == ObjectTypePlayer {
			if roomCfg.Roomtype <= constant.PVEPeacefulMap || roomCfg.Roomtype >= constant.SoloTowerMap {
				p.logger.Debugf("目标无法攻击,房间类型 %+v userId: %d , targetId:%d ", roomCfg.Roomtype, p.uid, t.Uid())
				continue
			}
		}
		// 和平模式不能攻击除怪以外目标
		if p.desk.roomCfg.Roomtype == constant.PVEPeacefulMap {
			if t.ObjectType() != ObjectTypeMonster {
				continue
			}
		}
		// 安全区内不能攻击玩家
		if p.desk.roomCfg.Roomtype == constant.PVEWarMap {
			if t.ObjectType() != ObjectTypeMonster && (p.desk.mapInfo.safeArea(p.location) || p.desk.mapInfo.safeArea(t.Location())) {
				p.logger.Debugf("安全去不能攻击玩家  位置:%v , 目标坐标:%v  skillID: %d ", p.location, t.Location(), skillId)
				continue
			}
		}

		caster.target = t

		caster.direction = maputil.DirectionFromPointF(p.location, t.Location())
		castDir = caster.direction

		targetsUids = append(targetsUids, t.Uid())

		// 伤害延迟生效
		pdesk := p.desk
		pt := t
		ilocation := &p.cell.location
		time.AfterFunc(time.Duration(skillCfg.DelayTime)*time.Millisecond, func() {
			if pdesk != nil {
				pdesk.addDelayAction(castObjType, cast_uid, skillId, pt, castDir, ilocation)
			}
		})
	}

	// 技能buff触发
	caster.SkillBuffTrigger(skill)

	// 技能CD
	if skillCfg.CoolDown > 0 {
		// 减 600 防止客户端技能CD未结束，服务器校验不通过
		CD := skillCfg.CoolDown - Cooldown
		if CD < 600 {
			CD = 600
		}
		cooldown_ := time.Now().Add(time.Duration(CD-600) * time.Millisecond)
		skill.Cooldown = cooldown_
		p.pushMsg(protocol.RouteSkillCooldownPush, &pb.SkillCooldownPush{CfgId: skillId, Delay: CD})

		if skillCfg.SType == constant.SkillTypeArtifact {
			for _, v := range p.skills {
				if cfg, err := global.GetSkillCfg(v.CfgId); err == nil && cfg.SType == constant.SkillTypeArtifact {
					v.Cooldown = cooldown_
					p.pushMsg(protocol.RouteSkillCooldownPush, &pb.SkillCooldownPush{CfgId: v.CfgId, Delay: CD})
				}
			}
		}
	}

	if cast_uid == p.uid {
		p.setAttackTime(global.GetActionsetCfg(skillCfg.SkillID))
	}

	packet := &pb.AttackPush{
		AttackerId: cast_uid,
		TargetUids: targetsUids,
		SkillId:    skillId,
		Location:   location,
		Direction:  castDir,
	}
	p.desk.Broadcast(protocol.RouteAttackPush, packet)

	// p.logger.Debugf("====[%s]====释放技能   玩家位置=%v skillID %d skillName: %s    朝向=%d", caster.Name(), p.location, skillId, skillCfg.Name, p.direction)

	p.SetObAction(ObjectActionAttack)
}

// ====================================================================
func (p *Player) move(location *pb.Point) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.desk == nil || !p.desk.mapCfg.Cells[*global.ToPointI(location)] {
		p.logger.Debugf("P 移动失败,不能到达   坐标: %v", location)
		return
	}

	if !p.CanMove() {
		p.logger.Debugf("P 移动失败,移动频率过快!/ 已死亡:[%t]   距离下次移动时间 %d", p.dead, p.moveTime.UnixMilli()-time.Now().UnixMilli())
		return
	}

	if maputil.Distance(p.location, location) > 2 {
		p.logger.Debugf("P 移动失败,距离过远 起点 %v  终点 %v", p.location, location)
		return
	}

	p.miningBreak()

	if buff := p.getBuffByType(constant.BuffTaunt); buff != nil {
		ownerDir := maputil.DirectionFromPointF(p.location, buff.Owner.Location())
		goDir := maputil.DirectionFromPointF(p.location, location)
		if ownerDir != goDir {
			p.logger.Warningf("P 移动方向与嘲讽者方向不一致,被嘲讽状态只能朝嘲讽者移动 ownerDir %v  goDir %v", ownerDir, goDir)
			// return
		}
	}

	if p.desk.mapInfo.locationChange(p, location) {
		packet := &pb.MovePush{Uid: p.uid, Location: p.location, Direction: p.direction}
		p.desk.Broadcast(protocol.RouteMovePush, packet)

		//地面可拾取buff
		cells := p.desk.mapInfo.getSquareCells(&p.cell.location, 1)
		for _, v := range cells {
			for _, obj := range v.objects {
				if obj.ObjectType() == ObjectTypeBuff {
					if cfg, err := global.GetBuffCfg(obj.CfgID()); err == nil {
						if cfg.CanPickUp == 1 {
							if p.desk.delObject(obj.Uid()) {
								p.AddBuff(obj.CfgID(), p, 0, nil, 0, 0)
								break
							}
						}
					}
				}
			}
		}
	} else {
		p.pushMsg(protocol.RouteMovePush, &pb.MovePush{Uid: p.Uid(), Location: p.location, Direction: p.direction})
	}
	p.moveTime = time.Now().Add(time.Duration(p.stat(pb.Attribute_MOVE_SPEED)-100) * time.Millisecond)
	p.SetObAction(ObjectActionMove)
}

func (p *Player) toDbBuff() []byte {
	dbbuff := &pb.DBBuffList{Buffs: []*pb.DBBuffInfo{}}

	p.buffMgr.IterateBuff(func(buf *BuffInfo) bool {
		if buf.Mark == constant.BuffMarkPassive {
			return true
		}

		switch constant.BuffType(buf.Pretype) {
		case constant.BuffPercentEquipStarIncrease:
			one := &pb.DBBuffInfo{
				Id:            buf.Id,
				BuffType:      buf.BuffType,
				StartTime:     buf.StartTime.Unix(),
				OverTime:      buf.OverTime.Unix(),
				TickTime:      buf.TickTime.Unix(),
				CDTime:        buf.CDTime.Unix(),
				TickFrequency: buf.TickFrequency,
				Pretype:       buf.Pretype,
				OwnerId:       buf.Owner.Uid(),
				Status:        buf.Status,
				BuffStats:     []*pb.Stats{},
			}
			buf.BuffStats.Range(func(k pb.Attribute, v int32) bool {
				one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: k, Num: v})
				return true
			})
			//for k, v := range buf.BuffStats.Attrs {
			//	one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: k, Num: v})
			//}
			dbbuff.Buffs = append(dbbuff.Buffs, one)
		}

		return true
	})
	//for _, buf := range p.buffs {
	//	if buf.Mark == constant.BuffMarkPassive {
	//		continue
	//	}
	//
	//	switch constant.BuffType(buf.Pretype) {
	//	case constant.BuffPercentEquipStarIncrease:
	//		one := &pb.DBBuffInfo{
	//			Id:            buf.Id,
	//			BuffType:      buf.BuffType,
	//			StartTime:     buf.StartTime.Unix(),
	//			OverTime:      buf.OverTime.Unix(),
	//			TickTime:      buf.TickTime.Unix(),
	//			CDTime:        buf.CDTime.Unix(),
	//			TickFrequency: buf.TickFrequency,
	//			Pretype:       buf.Pretype,
	//			OwnerId:       buf.Owner.Uid(),
	//			Status:        buf.Status,
	//			BuffStats:     []*pb.Stats{},
	//		}
	//
	//		for k, v := range buf.BuffStats.Attrs {
	//			one.BuffStats = append(one.BuffStats, &pb.Stats{Attr: k, Num: v})
	//		}
	//
	//		dbbuff.Buffs = append(dbbuff.Buffs, one)
	//	}
	//}
	if data, err := protobuf.NewSerializer().Marshal(dbbuff); err == nil {
		return data
	}
	return nil
}

func (p *Player) fromDbBuff(buffData []byte) {
	data := &pb.DBBuffList{}
	if err := protobuf.NewSerializer().Unmarshal(buffData, data); err == nil {
		for _, dbbuff := range data.Buffs {
			if dbbuff.OverTime < time.Now().Unix() {
				continue
			}

			buff := &BuffInfo{
				Id:            dbbuff.Id,
				BuffType:      dbbuff.BuffType,
				Mark:          dbbuff.Mark,
				StartTime:     time.Unix(dbbuff.StartTime, 0),
				OverTime:      time.Unix(dbbuff.OverTime, 0),
				TickTime:      time.Unix(dbbuff.TickTime, 0),
				CDTime:        time.Unix(dbbuff.CDTime, 0),
				TickFrequency: dbbuff.TickFrequency,
				Pretype:       dbbuff.Pretype,
				BuffStats:     Stats{Attrs: make(map[pb.Attribute]int32)},
			}

			if owner, ok := defaultManager.player(dbbuff.OwnerId); ok {
				buff.Owner = owner
			} else {
				buff.Owner = p
			}
			for _, v := range dbbuff.BuffStats {
				buff.BuffStats.setStats(v.Attr, v.Num)
			}
			//p.buffs[buff.Id] = buff
			p.buffMgr.buffs.Set(buff.Id, buff)
		}
	}
}

func (p *Player) HoneB(buffData []byte) {
	data := &pb.DBBuffList{}
	if err := protobuf.NewSerializer().Unmarshal(buffData, data); err == nil {
		for _, dbbuff := range data.Buffs {
			if dbbuff.OverTime < time.Now().Unix() {
				continue
			}

			buff := &BuffInfo{
				Id:            dbbuff.Id,
				BuffType:      dbbuff.BuffType,
				Mark:          dbbuff.Mark,
				StartTime:     time.Unix(dbbuff.StartTime, 0),
				OverTime:      time.Unix(dbbuff.OverTime, 0),
				TickTime:      time.Unix(dbbuff.TickTime, 0),
				CDTime:        time.Unix(dbbuff.CDTime, 0),
				TickFrequency: dbbuff.TickFrequency,
				Pretype:       dbbuff.Pretype,
				BuffStats:     Stats{Attrs: make(map[pb.Attribute]int32)},
			}

			if owner, ok := defaultManager.player(dbbuff.OwnerId); ok {
				buff.Owner = owner
			} else {
				buff.Owner = p
			}
			for _, v := range dbbuff.BuffStats {
				buff.BuffStats.setStats(v.Attr, v.Num)
			}
			//p.buffs[buff.Id] = buff
			p.buffMgr.buffs.Set(buff.Id, buff)
		}
	}
}

func (p *Player) defaultBuff() {
	// 宠物加buff
	// 宠物buff数组，每个元素包含ID和对应的属性

	// 初始化buff
	//for _, v := range p.buffs {
	//	p.delBuff(v.Id)
	//}

	ids := p.buffMgr.getBuffIds()
	for _, id := range ids {
		p.delBuff(id)
	}
	// 获取宠物天赋Buff列表
	petBuffs := p.getPetBuffs()
	for k, v := range global.BuffCfg {
		if v.Mark == constant.BuffMarkPassive {
			p.AddBuffNoRefresh(k, p, 24*60*60*1000, nil, 0, 0)
		}
		// 如果在这个数组里面，就加buff
		if len(petBuffs) > 0 {
			// 如果这个buff在这个数组里面，就加buff
			for _, vv := range petBuffs {
				if vv.Id == k {
					p.AddBuffNoRefresh(k, p, 24*60*60*1000, vv.Stats, vv.Rate, 0)
				}
			}
		}

	}
	p.refreshStats()
}

type PetBuff struct {
	Id    int32
	Rate  int32
	Stats *Stats
}

// 获取宠物天赋Buff列表
func (p *Player) getPetBuffs() []*PetBuff {
	//petBuffs := []*PetBuff{}
	//if len(p.pet) > 0 {
	//	if p.pet[0].PetCfgId != 0 {
	//		if petcfg := global.GetPetCfg(p.pet[0].PetCfgId); petcfg != nil {
	//			for _, v := range petcfg.Talent {
	//				bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
	//				if len(v.Value) > 0 {
	//					for _, vv := range v.Value {
	//						bufStats.setStats(pb.Attribute(vv.ID), vv.Num)
	//					}
	//				}
	//				petBuffs = append(petBuffs, &PetBuff{Id: v.ID, Stats: bufStats, Rate: v.Rate})
	//			}
	//		}
	//	}
	//}
	//return petBuffs
	return p.PetMgr.getPetBuffs()
}

func (p *Player) syncBuff() {
	packet := &pb.BuffListPush{BuffList: p.getPBBuffList()}
	p.pushMsg(protocol.RouteBuffListPush, packet)
}

// 穿脱装备opType 1穿 ，2脱  partnerSlotId 伙伴槽位 0 为主角
func (p *Player) changeEquip(opType int32, slot pb.EquipSlot, itemUid int64, partnerSlotId int32) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if slot <= pb.EquipSlot_INVALID_EquipSlot || slot >= pb.EquipSlot_MaxEquipSlot {
		return errutil.ErrEquipSlot
	}

	if partnerSlotId < 0 || partnerSlotId > 5 {
		return errutil.ErrEquipSlot
	}

	var equipSlot *Slot
	var job_ int
	var level_ int
	if partnerSlotId == 0 {
		equipSlot = &p.equips[slot]
		job_ = int(p.job)
		level_ = int(p.level)
	} else {
		equipSlot = &p.getPartnerSlotInfo(partnerSlotId).Equips[slot]

		if par := p.getPartnerInfoBySlotId(partnerSlotId); par != nil {
			job_ = int(par.cfg.Job) % 10
			level_ = int(par.level)
		} else {
			return errutil.ErrPartnerInvalid
		}
	}

	if opType == 1 {
		// item是最新的装备的信息
		if item := p.getItemByUid(itemUid); item != nil {
			itemCfg, err := global.GetItemCfg(item.CfgId)
			if err != nil {
				return err
			}
			if itemCfg.Equiptype != int(slot) {
				return errutil.ErrEquiptype
			}
			if itemCfg.Levellimit > level_ {
				return errutil.ErrLevel
			}
			if itemCfg.Occupylimit != 0 && itemCfg.Occupylimit != job_ {
				return errutil.ErrJob
			}
			// if itemCfg.Genderlimit != 0 && itemCfg.Genderlimit != p.sex {
			// 	return errutil.ErrGender
			// }
			// 检查是否已经在身上了
			if p.isOnEquip(item.Uid) {
				return errors.New(errutil.EquipOnUse)
			}
			// 同步包裹
			// 如果当前装备槽位有装备
			if equipSlot.ItemId != 0 {
				unequip := p.getItemByUid(equipSlot.ItemId)
				if unequip != nil {
					packet := &pb.AddSourcePush{Type: AddSourceTypeUnEquip, Item: []*pb.ItemInfo{}}
					packet.Item = append(packet.Item, unequip)
					p.pushMsg(protocol.RouteAddSourcePush, packet)
					// 在首页开宝箱穿装备时，需要把之前的装备替换进之前的宝箱
					// 遍历当前打开的宝箱
					for i := range p.box.NowOpenRes {
						// 在候选区找到这个装备
						if p.box.NowOpenRes[i].Uid == itemUid {
							// 替换这个值
							p.box.NowOpenRes[i].CfgId = unequip.CfgId
							p.box.NowOpenRes[i].Uid = equipSlot.ItemId
							p.box.NowOpenRes[i].Num = 1
							p.box.NowOpenRes[i].Equip = &pb.EquipInfo{Quality: pb.ItemQuality(equipSlot.Quality)}
							break

						}
					}
				}

			} else {
				for i := range p.box.NowOpenRes {
					if p.box.NowOpenRes[i].Uid == itemUid {
						// 删除这个值
						p.box.NowOpenRes = append(p.box.NowOpenRes[:i], p.box.NowOpenRes[i+1:]...)
						break
					}
				}
			}
			packet2 := &pb.TakeSourcePush{Item: []*pb.ItemInfo{}}
			packet2.Item = append(packet2.Item, item)
			p.pushMsg(protocol.RouteTakeSourcePush, packet2)

			equipSlot.ItemId = item.Uid
			equipSlot.Quality = int32(item.Equip.Quality)
		}
		if partnerSlotId == 0 {
			p.tasks.dispatch(pb.TargetType_EQUIP_ITEMS_WORN_TOTAL, 0, int32(len(p.equips)))
		}
	} else if opType == 2 {

		if p.bagEquipLeftCount() == 0 {
			return errutil.ErrBagFull
		}
		item := p.getItemByUid(equipSlot.ItemId)
		if item != nil {
			packet := &pb.AddSourcePush{Type: AddSourceTypeUnEquip, Item: []*pb.ItemInfo{}}
			packet.Item = append(packet.Item, item)
			p.pushMsg(protocol.RouteAddSourcePush, packet)
			equipSlot.ItemId = 0
		}
	}

	return nil
}

// equipItem 穿装备
func (p *Player) equipItem(slot pb.EquipSlot, itemUid int64, partnerSlotId int32) error {
	if slot <= pb.EquipSlot_INVALID_EquipSlot || slot >= pb.EquipSlot_MaxEquipSlot {
		return errutil.ErrEquipSlot
	}

	if partnerSlotId < 0 || partnerSlotId > 5 {
		return errutil.ErrEquipSlot
	}

	var equipSlot *Slot
	var job_ int
	var level_ int
	if partnerSlotId == 0 {
		equipSlot = &p.equips[slot]
		job_ = int(p.job)
		level_ = int(p.level)
	} else {
		equipSlot = &p.getPartnerSlotInfo(partnerSlotId).Equips[slot]
		if par := p.getPartnerInfoBySlotId(partnerSlotId); par != nil {
			job_ = int(par.cfg.Job) % 10
			level_ = int(par.level)
		} else {
			return errutil.ErrPartnerInvalid
		}
	}

	item := p.getItemByUid(itemUid)
	if item == nil {
		return errutil.ErrItemNotFound
	}

	itemCfg, err := global.GetItemCfg(item.CfgId)
	if err != nil {
		return err
	}

	// 检查装备类型
	if itemCfg.Equiptype != int(slot) {
		return errutil.ErrEquiptype
	}

	// 检查等级限制
	if itemCfg.Levellimit > level_ {
		return errutil.ErrLevel
	}

	// 检查职业限制
	if itemCfg.Occupylimit != 0 && itemCfg.Occupylimit != job_ {
		return errutil.ErrJob
	}

	// 检查是否已经在身上
	if p.isOnEquip(item.Uid) {
		return errors.New(errutil.EquipOnUse)
	}

	// 如果当前装备槽位有装备，先卸下
	if equipSlot.ItemId != 0 {
		unequip := p.getItemByUid(equipSlot.ItemId)
		if unequip != nil {
			// 在首页开宝箱穿装备时，需要把之前的装备替换进之前的宝箱
			for i := range p.box.NowOpenRes {
				if p.box.NowOpenRes[i].Uid == itemUid {
					p.box.NowOpenRes[i].CfgId = unequip.CfgId
					p.box.NowOpenRes[i].Uid = equipSlot.ItemId
					p.box.NowOpenRes[i].Num = 1
					p.box.NowOpenRes[i].Equip = &pb.EquipInfo{Quality: pb.ItemQuality(equipSlot.Quality)}
					break
				}
			}
		}
	} else {
		// 从宝箱候选区移除该装备
		for i := range p.box.NowOpenRes {
			if p.box.NowOpenRes[i].Uid == itemUid {
				p.box.NowOpenRes = append(p.box.NowOpenRes[:i], p.box.NowOpenRes[i+1:]...)
				break
			}
		}
	}

	// 装备到槽位
	equipSlot.ItemId = item.Uid
	equipSlot.Quality = int32(item.Equip.Quality)

	return nil
}

// unequipItem 脱装备
func (p *Player) unequipItem(slot pb.EquipSlot, partnerSlotId int32) error {
	if slot <= pb.EquipSlot_INVALID_EquipSlot || slot >= pb.EquipSlot_MaxEquipSlot {
		return errutil.ErrEquipSlot
	}

	if partnerSlotId < 0 || partnerSlotId > 5 {
		return errutil.ErrEquipSlot
	}

	var equipSlot *Slot
	if partnerSlotId == 0 {
		equipSlot = &p.equips[slot]
	} else {
		equipSlot = &p.getPartnerSlotInfo(partnerSlotId).Equips[slot]
	}

	// 检查包裹空间
	if p.bagEquipLeftCount() == 0 {
		return errutil.ErrBagFull
	}

	// 如果槽位有装备则卸下
	if equipSlot.ItemId != 0 {
		equipSlot.ItemId = 0
	}

	return nil
}

// 已装备道具数量
func (p *Player) getEquipNum() int {
	num := 0
	for i, v := range p.equips {
		if pb.EquipSlot(i) <= pb.EquipSlot_INVALID_EquipSlot || pb.EquipSlot(i) >= pb.EquipSlot_MaxEquipSlot {
			continue
		}
		if v.ItemId != 0 {
			num++
		}
	}
	for _, v := range p.partnerSlots {
		for _, vv := range v.Equips {
			if vv.ItemId != 0 {
				num++
			}
		}
	}
	return num
}

func (p *Player) hasEquipment(itemUid int64) bool {
	for _, v := range p.equips {
		if v.ItemId == itemUid {
			return true
		}
	}
	for _, v := range p.partnerSlots {
		for _, vv := range v.Equips {
			if vv.ItemId == itemUid {
				return true
			}
		}
	}
	return false
}

func (p *Player) hasHomeOpenBox(itemUid int64) bool {
	for i := range p.box.NowOpenRes {
		if p.box.NowOpenRes[i].Uid == itemUid {
			return true
		}
	}
	return false
}

func (p *Player) stats4Combat(combatStats *Stats) *Stats {
	// <0  减益 ; >0 增益
	combatStats.setStats(pb.Attribute_PHYSICAL_ATTACK, int32(float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAGIC_ATTACK, int32(float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_PHYSICAL_DEFENSE, int32(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAGIC_DEFENSE, int32(float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))

	combatStats.setStats(pb.Attribute_DODGE, int32(float64(combatStats.getStats(pb.Attribute_DODGE))*(float64(combatStats.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_HIT, int32(float64(combatStats.getStats(pb.Attribute_HIT))*(float64(combatStats.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_CRITICAL_STRIKE, int32(float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_TOUGHNESS, int32(float64(combatStats.getStats(pb.Attribute_TOUGHNESS))*(float64(combatStats.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))
	combatStats.setStats(pb.Attribute_MAX_HEALTH, int32(float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))*float64(combatStats.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate))

	// 攻速移速增加  基础*(1-加成)
	combatStats.changeStats(pb.Attribute_ATTACK_SPEED, int32(float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED))*(1-float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED_PERCENTAGE))/BaseRate)))
	combatStats.changeStats(pb.Attribute_MOVE_SPEED, int32(float64(combatStats.getStats(pb.Attribute_MOVE_SPEED))*(1-(float64(combatStats.getStats(pb.Attribute_MOVE_SPEED_PERCENTAGE))/BaseRate))))

	// 基础属性 >= 0
	combatStats.changeStats(pb.Attribute_PHYSICAL_ATTACK, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_PHYSICAL_ATTACK)))))
	combatStats.changeStats(pb.Attribute_MAGIC_ATTACK, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAGIC_ATTACK)))))
	combatStats.changeStats(pb.Attribute_PHYSICAL_DEFENSE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE)))))
	combatStats.changeStats(pb.Attribute_MAGIC_DEFENSE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))))

	combatStats.changeStats(pb.Attribute_DODGE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_DODGE)))))
	combatStats.changeStats(pb.Attribute_HIT, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_HIT)))))
	combatStats.changeStats(pb.Attribute_CRITICAL_STRIKE, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE)))))
	combatStats.changeStats(pb.Attribute_TOUGHNESS, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_TOUGHNESS)))))
	combatStats.changeStats(pb.Attribute_ATTACK_SPEED, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_ATTACK_SPEED)))))
	combatStats.changeStats(pb.Attribute_MOVE_SPEED, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MOVE_SPEED)))))
	combatStats.changeStats(pb.Attribute_MAX_HEALTH, int32(math.Max(0, float64(combatStats.getStats(pb.Attribute_MAX_HEALTH)))))

	return combatStats
}

// 战力（不计算buff）
func (p *Player) calcCombatReal(stats *Stats) int32 {
	combatStats := p.stats4Combat(stats)
	if combatStats == nil {
		return 0
	}
	// oldcombat := p.combat
	var skillCombat float64
	for _, v := range p.skills {
		if cfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			skillCombat += float64(cfg.Powerfactor) * (1 + float64(v.AdvancedLevel))
		}
	}

	// 减伤阈值75%
	dmt_reduce := combatStats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER) + combatStats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER) + combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)
	if dmt_reduce > 7500 {
		dmt_reduce = 7500
	}

	// 有效生命
	effective_hp := (float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))+global.GlobalCfg.BaseValue.Power.DefValue*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))+float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))/2.0)/(1-float64(dmt_reduce)/BaseRate) +
		(float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))/2.0 +
		float64(combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED))
	// p.logger.Debugf("  第1部分 %f   第2部分 %f    第3部分 %f   ===有效生命 %f ", (float64(combatStats.getStats(pb.Attribute_MAX_HEALTH))+global.GlobalCfg.BaseValue.Power.DefValue*(float64(combatStats.getStats(pb.Attribute_PHYSICAL_DEFENSE))+float64(combatStats.getStats(pb.Attribute_MAGIC_DEFENSE)))/2.0)/(1-float64(dmt_reduce)/BaseRate),
	// 	(float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))/2.0,
	// 	float64(combatStats.getStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED)), effective_hp)
	// 有效输出
	effective_power := float64(combatStats.getStats(p.AttackAttributes()))*
		(1+float64(combatStats.getStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate)*(1+skillCombat/100) +
		(float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)))/2.0 +
		float64(combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_FIXED))

	// p.logger.Debugf(" 面板攻击 %f  (1+增伤百分比+爆伤) %f  技能加成 %f   固定值 %f   ===有效输出 %f ", float64(combatStats.getStats(p.AttackAttributes())),
	// 	(1 + float64(combatStats.getStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)+combatStats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate),
	// 	(1 + skillCombat/100), (float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))+float64(combatStats.getStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)))/2.0, effective_power)

	combat := math.Sqrt(effective_hp*global.GlobalCfg.BaseValue.Power.LifeValue)*math.Sqrt(effective_power*global.GlobalCfg.BaseValue.Power.AttackValue) +
		float64(combatStats.getStats(pb.Attribute_HIT))*global.GlobalCfg.BaseValue.Power.Hitvalue +
		float64(combatStats.getStats(pb.Attribute_DODGE))*global.GlobalCfg.BaseValue.Power.Dodgevalue +
		float64(combatStats.getStats(pb.Attribute_CRITICAL_STRIKE))*global.GlobalCfg.BaseValue.Power.Critvalue +
		float64(combatStats.getStats(pb.Attribute_TOUGHNESS))*global.GlobalCfg.BaseValue.Power.Toughness

	// p.logger.Debugf(" 最终战力 %f", combat)
	p.combat = int32(math.Round(combat))

	return p.combat
}

func (p *Player) calcCombatRealV2(stats *Stats, calculatePanelStats bool) int64 {
	if calculatePanelStats {
		return CalculateBattlePower(p.stats4Combat(stats), global.GlobalCfg.BaseValue.Power, global.SkillCfg, p.skills, p.job)
	} else {
		return CalculateBattlePower(stats, global.GlobalCfg.BaseValue.Power, global.SkillCfg, p.skills, p.job)
	}
}

// CalculateBattlePower 计算战力, 传入的 stats 是面属性
func CalculateBattlePower(stats *Stats, powerCfg global.PowerConfig, skillCfg map[int32]*global.SkillConfig, skills []*PlayerSkill, job constant.MowJob) int64 {
	var (
		// 一级属性
		attackValue           = float64(0)                                                           // 面板主攻击值 1	2
		maxHealth             = math.Max(0, float64(stats.getStats(pb.Attribute_MAX_HEALTH)))        // 面板生命值 12
		physicalDefense       = math.Max(0, float64(stats.getStats(pb.Attribute_PHYSICAL_DEFENSE)))  // 面板物防值	4
		magicDefense          = math.Max(0, float64(stats.getStats(pb.Attribute_MAGIC_DEFENSE)))     // 面板魔防值	5
		hitValue              = math.Max(0, float64(stats.getStats(pb.Attribute_HIT)))               // 命中值	7
		dodgeValue            = math.Max(0, float64(stats.getStats(pb.Attribute_DODGE)))             // 闪避值	6
		criticalValue         = math.Max(0, float64(stats.getStats(pb.Attribute_CRITICAL_STRIKE)))   // 暴击值	8
		toughnessValue        = math.Max(0, float64(stats.getStats(pb.Attribute_TOUGHNESS)))         // 韧性值	9
		stunValue             = math.Max(0, float64(stats.getStats(pb.Attribute_STUN)))              // 眩晕值	53
		stunResistanceValue   = math.Max(0, float64(stats.getStats(pb.Attribute_STUN_RESISTANCE)))   // 眩晕抗性值	54
		poisonValue           = math.Max(0, float64(stats.getStats(pb.Attribute_POISON)))            // 中毒值	55
		poisonResistanceValue = math.Max(0, float64(stats.getStats(pb.Attribute_POISON_RESISTANCE))) // 中毒抗性值	56

		// 二级属性
		// lifeStealTriggerPct     = float64(0)                                                                       // 吸血触发百分比
		lifeStealPct            = math.Max(0, float64(stats.getStats(pb.Attribute_SUCKING_BLOOD_PERCENTAGE)))      // 吸血比例	40
		finalDamageReductionPct = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)))   // 最终减伤百分比	35
		PvEDamageReductionPct   = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER))) // 对怪减伤百分比	33
		PvPDamageReductionPct   = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER)))  // 对人减伤百分比	32
		criticalDamagePct       = math.Max(0, float64(stats.getStats(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)))    // 暴击伤害增加百分比	24
		reflectDamageTriggerPct = float64(0)                                                                       // 反伤触发百分比 暂无投放
		reflectDamagePct        = math.Max(0, float64(stats.getStats(pb.Attribute_REFLECT_DAMAGE_PERCENTAGE)))     // 反伤比例	42
		PvEDamageIncreasePct    = math.Max(0, float64(stats.getStats(pb.Attribute_PVE_DAMAGE_PERCENTAGE)))         // 对怪增伤百分比	26
		PvPDamageIncreasePct    = math.Max(0, float64(stats.getStats(pb.Attribute_PVP_DAMAGE_PERCENTAGE)))         // 对人增伤百分比	25
		finalDamageIncreasePct  = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))) // 最终增伤百分比	39

		// 三级属性
		PvEDamageIncrease    = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_INCREASE_FIXED_MONSTER))) // 对怪增伤固定值	37
		PvEDamageReduction   = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))   // 对怪减伤固定值	31
		PvPDamageReduction   = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_REDUCE_FIXED_PLAYER)))    // 对人减伤固定值	30
		PvPDamageIncrease    = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_INCREASE_FIXED_PLAYER)))  // 对人增伤固定值	36
		finalDamageReduction = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_FINAL_REDUCE_FIXED)))     // 最终减伤固定值	34
		finalDamageIncrease  = math.Max(0, float64(stats.getStats(pb.Attribute_DMT_FINAL_INCREASE_FIXED)))   // 最终增伤固定值	38

		deffenseCoefficient         = powerCfg.DefValue         // 防御系数
		effectiveHPCoefficient      = powerCfg.LifeValue        // 有效生命系数
		effectivePowerCoefficient   = powerCfg.AttackValue      // 有效输出系数
		hitCoefficient              = powerCfg.Hitvalue         // 命中系数
		dodgeCoefficient            = powerCfg.Dodgevalue       // 闪避系数
		criticalCoefficient         = powerCfg.Critvalue        // 暴击系数
		toughnessCoefficient        = powerCfg.Toughness        // 韧性系数
		stunvalueCoefficient        = powerCfg.Stunvalue        // 眩晕系数
		stunResistanceCoefficient   = powerCfg.StunResistance   // 眩晕抗性系数
		poisonvalueCoefficient      = powerCfg.Poisonvalue      // 中毒系数
		poisonResistanceCoefficient = powerCfg.PoisonResistance // 中毒抗性系数
		reduceN                     = powerCfg.Reduce_N         // 减益系数
		reduceM                     = powerCfg.Reduce_M         // 减益系数
	)

	// 获取攻击值
	switch job {
	case constant.Warrior, constant.PartnerWarrior:
		attackValue = float64(stats.getStats(pb.Attribute_PHYSICAL_ATTACK))
	case constant.Wizard, constant.PartnerWizard:
		attackValue = float64(stats.getStats(pb.Attribute_MAGIC_ATTACK))
	case constant.Ritualist, constant.PartnerRitualist:
		attackValue = float64(stats.getStats(pb.Attribute_MAGIC_ATTACK))
	}

	// 技能加成效果=∑_1^n▒{(〖技能战力加成系数〗_n )×(〖1＋当前技能等阶〗_n )}
	skillBonus := float64(0)
	for _, v := range skills {
		cfg, ok := skillCfg[v.CfgId]
		if !ok {
			continue
		}
		skillBonus += float64(cfg.Powerfactor) * float64((1 + v.AdvancedLevel))
	}

	// 有效生命值计算：(面板生命值×【1+吸血触发百分比×（1+吸血比例）】+防御系数×(（面板物防值+面板魔防值）)⁄2)/(（1-最终减伤百分比-对怪减伤百分比/2-对人减伤百分比/2）)+(（对怪减伤固定值+对人减伤固定值）)/2+最终减伤固定值
	// lifeWithBlood := maxHealth * (1 + (lifeStealTriggerPct/BaseRate)*(1+lifeStealPct/BaseRate))
	lifeWithBlood := maxHealth * (1 + lifeStealPct/BaseRate)
	defenseValue := deffenseCoefficient * ((physicalDefense + magicDefense) / 2.0)
	// 计算总减伤百分比，并确保不超过阈值
	totalReductionPct := finalDamageReductionPct + PvEDamageReductionPct/2 + PvPDamageReductionPct/2
	// 减伤系数，确保不为0或负数
	reductionCoeff := 1 - totalReductionPct/BaseRate
	if reductionCoeff <= 0 {
		reductionCoeff = 0.01 // 最小减伤系数，避免除零或负数
	}

	effectiveHP := (lifeWithBlood+defenseValue)/reductionCoeff + ((PvEDamageReduction+PvPDamageReduction)/2.0)*reduceN + finalDamageReduction*reduceM

	// 有效输出计算：面板主攻击值×【1+增伤百分比+对人增伤百分比/2+对怪增伤百分比/2+暴击伤害增加百分比+反伤触发百分比×（1+反伤比例）】×(1+技能加成效果/100)+(（对怪增伤固定值+对人增伤固定值）)/2+最终增伤固定值
	damagePercentage := (finalDamageIncreasePct+PvEDamageIncreasePct/2+PvPDamageIncreasePct/2+criticalDamagePct)/BaseRate + (reflectDamageTriggerPct/BaseRate)*(1+reflectDamagePct/BaseRate)
	skillMultiplier := float64(1 + skillBonus/100)
	fixedDamageIncrease := (PvEDamageIncrease+PvPDamageIncrease)/2 + finalDamageIncrease

	effectivePower := attackValue*(1+damagePercentage)*skillMultiplier + fixedDamageIncrease

	// 战力计算：√(有效生命×有效生命系数)×√(有效输出×有效输出系数)+命中值×命中系数+闪避值×闪避系数+暴击值×暴击系数+韧性值×韧性系数+击晕值×击晕系数+抗击晕值×抗击晕系数+中毒值×中毒系数+毒抗值×毒抗系数
	combat := math.Sqrt(effectiveHP*effectiveHPCoefficient)*
		math.Sqrt(effectivePower*effectivePowerCoefficient) +
		hitValue*hitCoefficient +
		dodgeValue*dodgeCoefficient +
		criticalValue*criticalCoefficient +
		toughnessValue*toughnessCoefficient +
		stunValue*stunvalueCoefficient +
		stunResistanceValue*stunResistanceCoefficient +
		poisonValue*poisonvalueCoefficient +
		poisonResistanceValue*poisonResistanceCoefficient

	// fmt.Printf("战力计算细节: skillBonus=%f, skillMultiplier=%.4f, attackValue=%.2f, effectivePower=%.4f, effectiveHP=%.2f, combat=%.2f\n",skillBonus, skillMultiplier, attackValue, effectivePower, effectiveHP, combat)

	return int64(math.Round(combat))
}

// 战力
func (p *Player) calcCombat() int32 {
	var skillCombat float64
	for _, v := range p.skills {
		if cfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			skillCombat += float64(cfg.Powerfactor) * (1 + float64(cfg.CanupperLevel))
		}
	}

	// 减伤阈值75%
	dmt_reduce := p.stat(pb.Attribute_DMT_REDUCE_PERCENTAGE_PLAYER) + p.stat(pb.Attribute_DMT_REDUCE_PERCENTAGE_MONSTER) + p.stat(pb.Attribute_DMT_FINAL_REDUCE_PERCENTAGE)
	if dmt_reduce > 7500 {
		dmt_reduce = 7500
	}

	// 有效生命
	effective_hp := (float64(p.GetMaxHP())+global.GlobalCfg.BaseValue.Power.DefValue*(float64(p.GetAR())+float64(p.GetMR()))/2.0)/(1-float64(dmt_reduce)/BaseRate) +
		(float64(p.stat(pb.Attribute_DMT_REDUCE_FIXED_PLAYER))+float64(p.stat(pb.Attribute_DMT_REDUCE_FIXED_MONSTER)))/2.0 +
		float64(p.stat(pb.Attribute_DMT_FINAL_REDUCE_FIXED))

	// 有效输出
	effective_power := float64(p.stat(p.AttackAttributes()))*
		(1+float64(p.stat(pb.Attribute_PVP_DAMAGE_PERCENTAGE)+p.stat(pb.Attribute_PVE_DAMAGE_PERCENTAGE)+p.stat(pb.Attribute_CRITICAL_DAMAGE_PERCENTAGE)+p.stat(pb.Attribute_DMT_FINAL_INCREASE_PERCENTAGE))/BaseRate)*(1+skillCombat/100) +
		(float64(p.stat(pb.Attribute_DMT_INCREASE_FIXED_PLAYER))+float64(p.stat(pb.Attribute_DMT_INCREASE_FIXED_MONSTER)))/2.0 +
		float64(p.stat(pb.Attribute_DMT_FINAL_INCREASE_FIXED))

	combat := math.Sqrt(effective_hp*global.GlobalCfg.BaseValue.Power.LifeValue)*math.Sqrt(effective_power*global.GlobalCfg.BaseValue.Power.AttackValue) +
		float64(p.stat(pb.Attribute_HIT))*global.GlobalCfg.BaseValue.Power.Hitvalue +
		float64(p.stat(pb.Attribute_DODGE))*global.GlobalCfg.BaseValue.Power.Dodgevalue +
		float64(p.stat(pb.Attribute_CRITICAL_STRIKE))*global.GlobalCfg.BaseValue.Power.Critvalue +
		float64(p.stat(pb.Attribute_TOUGHNESS))*global.GlobalCfg.BaseValue.Power.Toughness

	return int32(math.Round(combat))
}

// 单个道具战力
func (p *Player) calcCombatEquip(equipInfo *pb.ItemInfo) int32 {
	if equipInfo == nil {
		return 0
	}
	equipCfg, err := global.GetItemCfg(equipInfo.CfgId)
	if err != nil {
		return 0
	}
	totalStatus := &Stats{Attrs: make(map[pb.Attribute]int32)}
	equipStatus := &Stats{Attrs: make(map[pb.Attribute]int32)}
	// 基础值攻\防: （装备基础值+槽位强化基础值)*升星百分比
	for _, attr := range equipCfg.Stats {
		if attr.Key == pb.Attribute_PHYSICAL_ATTACK || attr.Key == pb.Attribute_MAGIC_ATTACK ||
			attr.Key == pb.Attribute_PHYSICAL_DEFENSE || attr.Key == pb.Attribute_MAGIC_DEFENSE ||
			attr.Key == pb.Attribute_MAX_HEALTH {
			equipStatus.setStats(attr.Key, int32(math.Round(float64(attr.Value)*global.GetEquipqualityCoefficient(equipInfo.Equip.Quality))))
		} else {
			totalStatus.setStats(attr.Key, attr.Value)
		}
	}
	// 极品属性直接计算
	if equipInfo.Equip != nil {
		for _, attr := range equipInfo.Equip.Attributes {
			totalStatus.setStats(attr.Attr, attr.Num)
		}
	}

	equipStatus.setStats(pb.Attribute_PHYSICAL_ATTACK, p.GetDC()+int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_ATTACK, p.GetMC()+int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_ATTACK_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_PHYSICAL_DEFENSE, p.GetAR()+int32(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_PHYSICAL_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAGIC_DEFENSE, p.GetMR()+int32(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE))*(float64(equipStatus.getStats(pb.Attribute_MAGIC_DEFENSE_PERCENTAGE))/BaseRate)))
	equipStatus.setStats(pb.Attribute_MAX_HEALTH, p.GetMaxHP()+int32(float64(equipStatus.getStats(pb.Attribute_MAX_HEALTH))*(float64(equipStatus.getStats(pb.Attribute_HEALTH_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_DODGE, p.GetDodge()+int32(float64(equipStatus.getStats(pb.Attribute_DODGE))*(float64(equipStatus.getStats(pb.Attribute_DODGE_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_HIT, p.GetHit()+int32(float64(equipStatus.getStats(pb.Attribute_HIT))*(float64(equipStatus.getStats(pb.Attribute_HIT_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_CRITICAL_STRIKE, p.GetCrit()+int32(float64(equipStatus.getStats(pb.Attribute_CRITICAL_STRIKE))*(float64(equipStatus.getStats(pb.Attribute_CRITICAL_STRIKE_PERCENTAGE))/BaseRate)))
	// equipStatus.setStats(pb.Attribute_TOUGHNESS, p.GetTough()+int32(float64(equipStatus.getStats(pb.Attribute_TOUGHNESS))*(float64(equipStatus.getStats(pb.Attribute_TOUGHNESS_PERCENTAGE))/BaseRate)))

	for ek, ev := range equipStatus.Attrs {
		if ek == pb.Attribute_PHYSICAL_ATTACK || ek == pb.Attribute_MAGIC_ATTACK ||
			ek == pb.Attribute_PHYSICAL_DEFENSE || ek == pb.Attribute_MAGIC_DEFENSE ||
			ek == pb.Attribute_MAX_HEALTH {
			totalStatus.setStats(ek, ev)
		}
	}

	return int32(p.calcCombatRealV2(totalStatus, false))
}

func (p *Player) addPrivilege(shopType int32, privilegeId int32, expirationTime time.Time) error {
	if shopType == constant.ShopTypeSubscription {
		for _, v := range p.privilege {
			if v.Pid == privilegeId {
				return errutil.ErrPrivilegeExist
			}
		}
		p.privilege = append(p.privilege, protocol.PrivilegeDetailTCP{Pid: privilegeId, EndTime: expirationTime})
	} else if shopType == constant.ShopTypeAFK {
		// 当前时间加30天后的0点
		day := global.GlobalCfg.BaseValue.Expcopy.Lastday
		p.afkExpirationAt = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, int(day))
		p.afkDailyTime = 0
		p.syncSaveDb()
	}

	p.refreshStats()
	p.syncUserBase()
	p.logger.Debugf("添加特权卡成功 privilegeId=%d expirationTime=%v", privilegeId, expirationTime)

	return nil
}

// 更新排行榜

func (p *Player) rankChangeTCP(k int, v int64) {
	// if p.isRobot() {
	// 	if p.robotType == RobotTypeRobot {
	// 		return
	// 	}
	// 	if p.robotType == RobotTypePlayer && k != constant.RankTypeSoul {
	// 		return
	// 	}
	// }
	if p.isRobot() && p.robotType == RobotTypeRobot {
		if k != constant.RankTypeSoul && k != constant.RankTypeTurnBased {
			return
		}
	}

	if !p.isRankActivityOpen(k, v) {
		return
	}
	//有时间的榜单，超过时间直接不录入
	if constant.RankTypeServerOpenDay[int64(k)] {
		openDay := global.ServerListMgrInstance.GetServerOpenDay(int32(p.serverID))
		//天数
		for _, v := range global.GlobalCfg.LimitRanking.Key {
			//实际天数与配置天数一样，需要发奖励
			if v.Value < int64(openDay) && int64(v.ID) == int64(k) {
				return
			}
		}
	}
	// 不是第一名的时候发消息
	cansend := 0
	servername := ""
	rankMember := p.uid
	point := v
	RankServerId := int32(0)
	//一站到底榜单的特殊性
	if k == constant.RankTypeGightEndlayers {
		for k1, v1 := range p.FightEndMap.Evaluation {
			if v1.BattleId == int32(v) {
				point = int64(1000000 - p.FightEndMap.Evaluation[k1].ChallengeTime)
			}
		}
		rankMember = int64(v % 100000)
	}

	if k == constant.RankTypeCampaignStage || k == constant.RankTypeMainBattle {
		point = point<<32 | math.MaxInt32 - time.Now().Unix()
	}

	//需要区分服务器的榜单
	if RankTypeServer, ok := constant.RankTypeServer[int64(k)]; ok && RankTypeServer {
		//一站到底的每层榜单不需要发送
		if k != constant.RankTypeGightEndlayers {
			rank := provide.IsFirst(rankMember, k, p.serverID)
			if len(rank) == 0 {
				cansend = 1
			}
		}
		// serverId = p.serverID
	}
	//需要查询服务器名字的榜单
	if RankTypeNameServer, ok := constant.RankTypeNameServer[int64(k)]; ok && RankTypeNameServer {
		allServers := global.ServerListMgrInstance.GetAll()
		if serverInfo, exists := allServers[int(p.serverID)]; exists {
			servername = serverInfo.Name
			// 这里与一站到底冲突了
			RankServerId = p.serverID * 100
		}
	}
	AvatarV2 := int32(0)
	Frame := int32(0)
	if p.avatar != nil {
		AvatarV2 = int32(p.avatar.GetCurrentAvatar())
		Frame = int32(p.avatar.GetCurrentFrame())
	}
	provide.RankChangeTCP(&protocol.RankUserInfoInput{
		Uid:          rankMember,
		Type:         int32(k),
		AvatarV2:     AvatarV2,
		Frame:        Frame,
		Name:         p.name,
		Value:        point,
		Sex:          int32(p.sex),
		Job:          int32(p.job),
		Mysql:        true,
		Compare:      constant.RankTypeServerCompare[int64(k)],
		ServerId:     p.serverID,
		RankServerId: RankServerId,
		UserServerId: p.serverID,
		ServerName:   servername,
		RobotCId:     p.cfgId,
		ZoneId:       p.getZoneId(),
	})

	if p.bRobot && p.robotType == RobotTypeRobot && p.cfgId == 0 {
		logger.Errorf("p.rankChangeTCP机器人配置id为0  id: %d  name %s  cfgID %d ", p.uid, p.name, p.cfgId)
	}

	//符合条件的发送走马灯
	if RankTypeServer, ok := constant.RankTypeServer[int64(k)]; ok && RankTypeServer {
		if cansend == 1 {
			rank := provide.IsFirst(rankMember, k, p.serverID)
			if len(rank) > 0 {
				strSend := []string{p.name, constant.RankTypeName[k]}
				// 发送消息给所有在线用户
				defaultManager.Lamp(int32(constant.LampTopRankChanged), p.serverID, strSend)
			}
		}
	}
}

var rankTypeActivityMap = map[int]int32{
	constant.RankTypePinkEntry:        constant.ActivityTypePinkEntry,
	constant.RankTypeRedPartnerStar:   constant.ActivityTypeRedPartnerStar,
	constant.RankTypeEquipmentStar:    constant.ActivityTypeEquipmentStar,
	constant.RankTypePetSynthesis:     constant.ActivityTypePetSynthesis,
	constant.RankTypeSoulStoneLevel:   constant.ActivityTypeSoulStoneLevel,
	constant.RankTypeFiveElementLevel: constant.ActivityTypeFiveElementLevel,
}

func (p *Player) isRankActivityOpen(rankType int, value int64) bool {
	// 如果不是活动排行榜，直接允许
	activityType, ok := rankTypeActivityMap[rankType]
	if !ok {
		return true
	}

	// 检查活动是否在有效期内
	configs, ok := global.ActivityCfgByType[activityType]
	if !ok || len(configs) == 0 {
		return true
	}
	now := time.Now().Unix()
	activityOpen := false
	for _, cfg := range configs {
		beginTime := p.calculateActivityBeginTime(cfg.StartTime)
		if beginTime == 0 {
			continue
		}
		if cfg.ResetTime <= 0 {
			if now >= beginTime {
				activityOpen = true
				break
			}
			continue
		}
		originalTime := time.Unix(beginTime, 0)
		endTime := time.Date(originalTime.Year(), originalTime.Month(), originalTime.Day(), 0, 0, 0, 0, originalTime.Location()).
			AddDate(0, 0, int(cfg.ResetTime))
		if now >= beginTime && now < endTime.Unix() {
			activityOpen = true
			break
		}
	}
	if !activityOpen {
		return false
	}

	// 活动开放后，检查数量限制
	// 优先使用该排行榜的特定限制
	limitNum := int64(0)
	for _, limit := range global.GlobalCfg.SameDayDelivery.Limits {
		if int(limit.ID) == rankType {
			limitNum = limit.Num
			break
		}
	}
	// // 如果没有找到特定限制，使用默认限制
	// if limitNum == 0 {
	// 	limitNum = global.GlobalCfg.SameDayDelivery.Rank.Num
	// }
	// 检查是否满足数量限制
	if limitNum > 0 && value < limitNum {
		return false
	}

	return true
}

func (p *Player) calculateActivityBeginTime(startTime int32) int64 {
	if startTime == 0 {
		return p.registerAt
	}
	if startTime >= 100 {
		return int64(startTime)
	}
	if p.serverID <= 0 {
		return p.registerAt
	}
	serverInfo, err := global.ServerListMgrInstance.Get(int(p.serverID))
	if err != nil {
		return p.registerAt
	}
	openTime := time.Unix(serverInfo.OpenBy/1000, 0)
	openDayMidnight := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 0, 0, 0, 0, openTime.Location())
	return openDayMidnight.AddDate(0, 0, int(startTime)-1).Unix()
}

// 更新用户显示信息
func (p *Player) userShowRefresh() {
	dispatcher.MustGetInstance().Publish(dispatcher.TopicUserShowRefresh.String(), &dispatcher.UserShowRefreshEvent{
		UserID:   p.Uid(),
		ServerID: p.serverID,
		ZoneID:   p.getZoneId(),
		UserShowRefresh: &dispatcher.UserShowRefresh{
			Name:     p.name,
			Avatar:   int(p.avatar.GetCurrentAvatar()),
			Frame:    int(p.avatar.GetCurrentFrame()),
			Sex:      p.sex,
			Job:      int(p.job),
			RobotCID: 0,
			ServerID: int(p.serverID),
		},
	})
}

// 生成随机属性
func (p *Player) upgradeRandomAttr(item *pb.ItemInfo, itemCfg *global.ItemConfig) error {

	if item.Type != pb.ItemType_EQUIPMENT || item.Equip == nil {
		return errutil.ErrConfigNotFound
	}

	createCfg := global.GetRandomCreatCfg(int32(item.Equip.Quality), int32(itemCfg.Levellimit), int32(itemCfg.Equiptype))
	if createCfg == nil {
		return errutil.ErrConfigNotFound
	}

	var attrNum, distribUid, attrPart, weight int
	weight = global.RandIn(10000)
	for i := 0; i < len(createCfg.CatalogueNum); i++ {
		if weight < int(createCfg.CatalogueNum[i].Value) {
			attrNum = int(createCfg.CatalogueNum[i].Key)
			break
		}
	}

	if attrNum == 0 {
		return nil
	} else if attrNum > len(createCfg.KindNum) {
		attrNum = len(createCfg.KindNum)
	}

	attrList := []int{}
	tmpKindWeight := []global.KeyValue{}
	for _, v := range createCfg.KindNum {
		tmpKindWeight = append(tmpKindWeight, global.KeyValue{Key: v.Key, Value: v.Value})
	}

	fcontains := func(slice []int, val int) bool {
		for _, v := range slice {
			if v == val {
				return true
			}
		}
		return false
	}

	for i := 0; i < attrNum; i++ {
		if len(tmpKindWeight) == 0 {
			break
		}

		sort.Slice(tmpKindWeight, func(i, j int) bool {
			return tmpKindWeight[i].Value > tmpKindWeight[j].Value
		})
		weight = global.RandIn(int(tmpKindWeight[0].Value))
		sort.Slice(tmpKindWeight, func(i, j int) bool {
			return tmpKindWeight[i].Value < tmpKindWeight[j].Value
		})

		for wi, wv := range tmpKindWeight {
			if weight < int(wv.Value) && !fcontains(attrList, int(wv.Key)) {
				distribUid = int(wv.Key)
				attrList = append(attrList, distribUid)
				// 去重
				copy(tmpKindWeight[wi:], tmpKindWeight[wi+1:])
				tmpKindWeight = tmpKindWeight[:len(tmpKindWeight)-1]
				break
			}
		}

		// 随机分段
		weight = global.RandIn(10000)
		for _, v := range createCfg.DropSet {
			if weight < int(v.Value) {
				attrPart = int(v.Key)
				break
			}
		}

		if distribUid == 0 || attrPart == 0 {
			continue
		}
		distribCfg := global.GetRandomDistribCfg(int32(distribUid), int32(attrPart))
		if distribCfg == nil {
			continue
		}

		disCfg, ok := distribCfg.Interval[int32(item.Equip.Quality)]
		if !ok {
			continue
		}
		var statValue int32
		if disCfg.Max <= disCfg.Min {
			statValue = disCfg.Max
		} else {
			statValue = int32(global.RandIn(int(disCfg.Max)-int(disCfg.Min))+1) + disCfg.Min
		}
		if item.Equip == nil {
			item.Equip = &pb.EquipInfo{Attributes: []*pb.Stats{}}
		}
		item.Equip.Attributes = append(item.Equip.Attributes, &pb.Stats{Attr: pb.Attribute(distribCfg.Kind), Num: statValue})
		distribUid = 0
		attrPart = 0

		p.logger.Debugf("生成随机属性i=%d itemid=%d UUID=%d  key=%d  value=%d   attrlist=%v", i, item.CfgId, distribCfg.UUID, distribCfg.Kind, statValue, attrList)
	}

	return nil
}

func (p *Player) recharge4vip(cny int64) {
	old := p.vip
	p.recharge += cny
	for _, v := range global.VIPCfg {
		if p.recharge >= int64(v.NeedExp) && p.vip < v.VIPLevel {
			p.vip = v.VIPLevel
		}
	}
	p.syncUserBase()
	if old != p.vip {
		p.refreshStats()
	}
}

func (p *Player) changeRoleName(name string) (string, error) {
	if db.IsUserExistsByName(name) {
		return p.name, errutil.ErrAccountExists
	}

	if len(name) > 21 || len(name) < 6 {
		return p.name, errutil.ErrRoleNameLen
	}

	// if sensitive.ContainsSensitiveWord(name) || !sensitive.IsValidName(name) {
	// 	return p.name, errutil.ErrSensitiveWord
	// }

	Rename := global.GlobalCfg.BaseValue.Rename
	// 校验背包情况
	res := p.checkTakeSource(int32(Rename.Needitem), int64(Rename.Neednum))
	if !res {
		return p.name, errutil.ErrItemNotEnough
	}
	// 减少道具
	p.takeSource(int32(Rename.Needitem), int64(Rename.Neednum), constant.ItemDelChangeRoleName)
	p.name = name
	u := &model.User{}
	u.Id = p.uid
	u.Name = name
	if err2 := db.UpdateUser(u); err2 != nil {
		return p.name, err2
	}

	if err := thinkingdatax.TrackerIns.Set(strconv.FormatInt(p.Uid(), 10), p.accountName, tracking.M{
		thinkingdatax.Nickname.String(): name,
	}); err != nil {
		logger.Warnf("ThinkingData UserSet failed: accountID=%s, error=%v", p.accountName, err)
	}

	return name, nil
}

// 强化
func (p *Player) enhance(id pb.EquipSlot, partnerslotid int32) (bool, int32, error) {
	// 查询当前槽位的信息
	var equipsInfo *Slot
	if partnerslotid == 0 {
		equipsInfo = &p.equips[pb.EquipSlot(id)]

		// 最高等级不让操作
		if equipsInfo.Enhance >= (p.level * global.GlobalCfg.BaseValue.Equipstrenghthen.Rate) {
			return false, -1, errutil.ErrMaxLevel
		}
	} else {
		equipsInfo = &p.getPartnerSlotInfo(partnerslotid).Equips[pb.EquipSlot(id)]

		if partner := p.getPartnerInfoBySlotId(partnerslotid); partner != nil {
			if equipsInfo.Enhance >= (partner.level * global.GlobalCfg.BaseValue.Equipstrenghthen.Rate) {
				return false, -1, errutil.ErrMaxLevel
			}
		} else {
			return false, -1, errutil.ErrPartnerSlotNil
		}
	}
	var maxEnhance, maxEnhance2, enhanceLevel, enhanceLevel2 int32
	maxEnhance = math.MaxInt32
	maxEnhance2 = math.MaxInt32
	for i, equip := range p.equips {
		if equip.Enhance < maxEnhance {
			if i == 0 {
				continue
			}

			maxEnhance = equip.Enhance
		}
	}
	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryEnhance && maxEnhance >= equipmaster.Limit {
			enhanceLevel = equipmaster.Limit
		}
	}

	// 槽位的强化，拼接uuid
	uuid := int32(p.job%10*1000) + equipsInfo.Enhance
	// 目前升星所需要的信息
	enhanceNeed, err := global.GetEnhanceCfg(uuid)
	if err != nil {
		return false, -1, err
	}
	// 校验背包情况
	// 目前强化所需要的材料
	for _, v := range enhanceNeed.Consume {
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return false, equipsInfo.Enhance, errutil.ErrItemNotEnough
		}
	}

	Cnum := int32(0)
	Pnum := int32(0)
	// 减少道具
	for _, v := range enhanceNeed.Consume {
		if v.ID == int32(constant.ItemTypeCoin) {
			Cnum = v.Num
		} else {
			Pnum = v.Num
		}
		p.takeSource(int32(v.ID), int64(v.Num), constant.ItemDelEnhance)
	}
	// 升级
	oldLevel := equipsInfo.Enhance
	equipsInfo.Enhance++
	newLevel := equipsInfo.Enhance
	// 任务上报：全身装备强化等级需求
	minEnhance := int32(math.MaxInt32)
	for i, equip := range p.equips {
		if i == 0 {
			continue
		}
		if equip.Enhance < minEnhance {
			minEnhance = equip.Enhance
		}
	}
	if minEnhance > 0 {
		p.tasks.dispatch(pb.TargetType_EQUIPMENT_STRENGTHENING_LEVEL_REQUIREMENT, 0, minEnhance)
	}

	p.tasks.dispatch(pb.TargetType_ANY_EQUIP_STRENGTHENED_TO_LEVEL, newLevel, 1)

	p.tasks.dispatch(pb.TargetType_TOTAL_STRENGTHEN_TIMES, 0, 1)

	if partnerslotid == 0 {
		p.tasks.dispatch(pb.TargetType_MAIN_ALL_EQUIP_SLOTS_STRENGTHENED_TO_LEVEL, 0, p.GetMainRoleAllSlotMinEnhanceLevel())
	}

	p.tasks.dispatch(pb.TargetType_ALL_EQUIP_SLOTS_STRENGTHENED_TO_LEVEL, 0, p.GetAllSlotMinEnhanceLevel())

	oldcombat := p.totalCombat

	p.refreshStats()

	for i, equip := range p.equips {
		if i == 0 {
			continue
		}
		if equip.Enhance < maxEnhance2 {
			maxEnhance2 = equip.Enhance
		}
	}
	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryEnhance && maxEnhance2 >= equipmaster.Limit {
			enhanceLevel2 = equipmaster.Limit
		}
	}
	if enhanceLevel2 > enhanceLevel {
		p.upSuitShushuDate(constant.EquipMasterCategoryEnhance, enhanceLevel2, partnerslotid)
	}
	newcombat := p.totalCombat
	p.upEnhanceShushuDate(partnerslotid, int32(id), Cnum, Pnum, newcombat-oldcombat, oldLevel, newLevel)

	return true, equipsInfo.Enhance, nil
}

func (p *Player) upSuitShushuDate(resonance_type int32, resonance_level int32, lineup_slot int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "resonance"
	// 新增字段
	if constant.EquipMasterCategoryName[resonance_type] == "" {
		return
	}
	shushuDate["resonance_type"] = constant.EquipMasterCategoryName[resonance_type] // 可根据实际情况填写伙伴数量
	shushuDate["resonance_level"] = resonance_level
	shushuDate["lineup_slot"] = lineup_slot
	// 可根据实际情况填写装备部位
	shushumain.ReportData(shushuDate)
	shushuUserSetDate := make(map[string]interface{}, 0)

	switch resonance_type {
	case constant.EquipMasterCategoryEnhance:
		shushuUserSetDate["improve_resonance"] = resonance_level
	case constant.EquipMasterCategoryStars:
		shushuUserSetDate["star_resonance"] = resonance_level
	case constant.EquipMasterCategoryGemStone:
		shushuUserSetDate["stone_resonance"] = resonance_level
	}
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (p *Player) upEnhanceShushuDate(partnerslotid int32, equip_part int32, Cnum int32, Pnum int32, combat int32, oldLevel int32, newLevel int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "equip_position_strengthen"
	// 新增字段
	shushuDate["partner_num"] = partnerslotid // 可根据实际情况填写伙伴数量
	shushuDate["equip_part"] = equip_part     // 可根据实际情况填写装备部位
	shushuDate["cost_item_num"] = Pnum        // 可根据实际消耗道具数量填写
	shushuDate["cost_gold"] = Cnum            // 可根据实际消耗金币数量填写
	// shushuDate["equip_power"] = 0              // 可根据实际装备战力填写
	shushuDate["equip_power_up"] = combat      // 可根据实际装备战力提升填写
	shushuDate["strengthen_after"] = newLevel  // 可根据实际强化等级填写
	shushuDate["strengthen_before"] = oldLevel // 可根据实际强化等级填写
	shushumain.ReportData(shushuDate)
}

// 升星
func (p *Player) stars(id pb.EquipSlot, partnerslotid int32) (bool, int32, error) {
	// 查询当前槽位的信息
	var equipsInfo *Slot
	if partnerslotid == 0 {
		equipsInfo = &p.equips[pb.EquipSlot(id)]
	} else {
		equipsInfo = &p.getPartnerSlotInfo(partnerslotid).Equips[pb.EquipSlot(id)]
	}
	if equipsInfo == nil {
		return false, -1, errutil.ErrEquipSlot
	}
	// 最高等级不让操作
	if equipsInfo.Stars >= global.GlobalCfg.BaseValue.Equipstar.Maxlevel {
		return false, -1, errutil.ErrMaxLevel
	}
	// 槽位的强化，拼接uuid
	uuid := int32(p.job*100) + equipsInfo.Stars
	// 目前升星所需要的信息
	StarsNeed, err := global.GetStarsCfg(uuid)
	if err != nil {
		return false, -1, err
	}
	// 校验背包情况
	// 目前升星所需要的材料
	for _, v := range StarsNeed.MaterialID {
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return false, -1, errutil.ErrItemNotEnough
		}
	}
	var maxStars, maxStars2, starLevel, starLevel2 int32
	maxStars = math.MaxInt32
	maxStars2 = math.MaxInt32
	for i, equip := range p.equips {
		if equip.Stars < maxStars {
			if i == 0 {
				continue
			}
			maxStars = equip.Stars
		}
	}
	for _, equipmaster := range global.EquipMasterCfg {
		if equipmaster.Category == constant.EquipMasterCategoryStars && maxStars >= equipmaster.Limit {
			starLevel = equipmaster.Limit
		}
	}
	Cnum := int32(0)
	Pnum := int32(0)
	oldstars := equipsInfo.Stars
	success := false
	// 减少道具
	for _, v := range StarsNeed.MaterialID {
		if v.ID == int32(constant.ItemTypeCoin) {
			Cnum = v.Num
		} else {
			Pnum = v.Num
		}
		p.takeSource(int32(v.ID), int64(v.Num), constant.ItemDelStar)
	}
	// 随机分段
	weight := global.RandIn(int(BaseRate))
	// 查询装备信息
	var equipCfg *global.ItemConfig
	equip_quality := int32(0)
	if equipsInfo.ItemId > 0 {
		equipCfg2 := p.getItemByUid(equipsInfo.ItemId)
		if equipCfg2 == nil {
			return false, equipsInfo.Stars, err
		}
		equip_quality = int32(equipCfg2.Equip.Quality)
		equipCfg, err = global.GetItemCfg(int32(equipCfg2.CfgId))
		if err != nil {
			return false, equipsInfo.Stars, err
		}
	}

	haveUp := false
	if p.stat(pb.Attribute_EQUIP_STAR_PERCENTAGE) > 0 {
		haveUp = true
	}
	if float64(weight) < float64(StarsNeed.SuccessRate)*(1+float64(p.stat(pb.Attribute_EQUIP_STAR_PERCENTAGE))/BaseRate) {
		// 升级
		equipsInfo.Stars++
		if equipsInfo.Stars >= 10 && equipsInfo.Stars%5 == 0 {
			context := []string{p.name, strconv.Itoa(int(equipsInfo.Stars))}
			defaultManager.Lamp(int32(constant.LampEquipmentStarUpgrade), p.serverID, context)
		}
		success = true
		newstars := equipsInfo.Stars
		oldcombat := p.totalCombat
		p.refreshStats()
		newcombat := p.totalCombat
		for i, equip := range p.equips {
			if equip.Stars < maxStars2 {
				if i == 0 {
					continue
				}
				maxStars2 = equip.Stars
			}
		}
		for _, equipmaster := range global.EquipMasterCfg {
			if equipmaster.Category == constant.EquipMasterCategoryStars && maxStars2 >= equipmaster.Limit {
				starLevel2 = equipmaster.Limit
			}
		}
		if starLevel2 > starLevel {
			p.upSuitShushuDate(constant.EquipMasterCategoryStars, starLevel2, partnerslotid)
		}
		p.rankChangeTCP(constant.RankTypeEquipmentStar, p.calcEquipmentStarNum())
		p.tasks.dispatch(pb.TargetType_ANY_EQUIP_STAR_UP_TIMES, 0, 1)
		p.tasks.dispatch(pb.TargetType_EQUIP_STAR_TO_X, newstars, 1)
		if partnerslotid == 0 {
			p.tasks.dispatch(pb.TargetType_MAIN_ALL_EQUIP_SLOTS_STAR_TO_LEVEL, 0, p.GetMainRoleAllSlotMinStar())
		}
		p.tasks.dispatch(pb.TargetType_ALL_EQUIP_SLOTS_STAR_TO_LEVEL, 0, p.GetAllSlotMinStar())

		// 任务上报：全身装备升星等级需求
		minStar := int32(math.MaxInt32)
		for i, equip := range p.equips {
			if i == 0 {
				continue
			}
			if equip.Stars < minStar {
				minStar = equip.Stars
			}
		}
		if minStar > 0 {
			p.tasks.dispatch(pb.TargetType_EQUIPMENT_UPGRADING_STAR_LEVEL_REQUIREMENT, 0, minStar)
		}

		p.upStarsShushuDate(equipCfg, partnerslotid, int32(id), Cnum, Pnum, newcombat-oldcombat, oldstars, newstars, success, weight, haveUp, p.stat(pb.Attribute_EQUIP_STAR_PERCENTAGE), equip_quality)
		return true, equipsInfo.Stars, nil
	} else {
		equipsInfo.Stars = StarsNeed.FailLevel
	}
	newstars := equipsInfo.Stars
	oldcombat := p.totalCombat
	p.refreshStats()
	newcombat := p.totalCombat
	p.upStarsShushuDate(equipCfg, partnerslotid, int32(id), Cnum, Pnum, newcombat-oldcombat, oldstars, newstars, success, weight, haveUp, p.stat(pb.Attribute_EQUIP_STAR_PERCENTAGE), equip_quality)

	return false, equipsInfo.Stars, nil
}

func (p *Player) upStarsShushuDate(equipsDetail *global.ItemConfig, partnerslotid int32, equip_part int32, Cnum int32, Pnum int32, combat int32, oldLevel int32, newLevel int32, success bool, weight int, haveUp bool, per int32, equip_quality int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "equip_star_up"
	// 新增字段
	if equipsDetail == nil {
		shushuDate["equip_id"] = 0      // 可根据实际装备ID填写
		shushuDate["equip_quality"] = 0 // 可根据实际装备品质填写
		shushuDate["equip_num"] = 0     // 可根据实际装备数量填写
		shushuDate["equip_name"] = ""   // 可根据实际装备名称填写
		shushuDate["equip_job"] = 0     // 可根据实际装备职业填写
	} else {
		shushuDate["equip_id"] = equipsDetail.ID           // 可根据实际装备ID填写
		shushuDate["equip_num"] = 1                        // 可根据实际装备数量填写
		shushuDate["equip_quality"] = equip_quality        // 可根据实际装备品质填写
		shushuDate["equip_name"] = equipsDetail.Name       // 可根据实际装备名称填写
		shushuDate["equip_job"] = equipsDetail.Occupylimit // 可根据实际装备职业填写
	}

	shushuDate["equip_part"] = equip_part               // 可根据实际装备部位填写
	shushuDate["equip_level"] = equipsDetail.Levellimit // 可根据实际装备等级填写
	shushuDate["partner_num"] = partnerslotid           // 可根据实际情况填写伙伴数量
	shushuDate["cost_item_num"] = Pnum                  // 可根据实际消耗道具数量填写
	shushuDate["cost_gold"] = Cnum                      // 可根据实际消耗金币数量填写
	shushuDate["star_after"] = newLevel                 // 可根据实际升星后星级填写
	shushuDate["star_before"] = oldLevel                // 可根据实际升星前星级填写
	shushuDate["is_success"] = success                  // 可根据实际升星是否成功填写
	shushuDate["success_rate"] = weight                 // 可根据实际升星成功率填写
	shushuDate["use_symbol"] = haveUp                   // 可根据实际使用符文数量填写
	shushuDate["symbol_name"] = per                     // 可根据实际符文名称填写
	// shushuDate["equip_power"] = 0               // 可根据实际装备战力填写
	shushuDate["equip_power_up"] = combat // 可根据实际装备战力提升填写
	shushumain.ReportData(shushuDate)
}

// 洗炼
func (p *Player) equipRefine(itemUid int64, attr pb.Attribute, random bool) error {
	item := p.getItemByUid(itemUid)
	if item == nil || item.Equip == nil {
		return errutil.ErrItemNotFound
	}

	bFind := false
	for _, v := range item.Equip.Attributes {
		if v.Attr == attr {
			bFind = true
		}
	}
	if !bFind {
		return errutil.ErrItemNotFound
	}

	itemCfg, err := global.GetItemCfg(item.CfgId)
	if err != nil {
		return err
	}
	var costId int32
	var costNum int64
	if random {
		for _, v := range global.GlobalCfg.Rifinerandom.Cost {
			if v.ID == itemCfg.Quality {
				if !p.checkTakeSource(v.ItemID, int64(v.ItemNum)) {
					return errutil.ErrItemNotEnough
				}
				costId = v.ItemID
				costNum = v.ItemNum
				break
			}
		}
	} else {
		for _, v := range global.GlobalCfg.Rifinedirect.Cost {
			if v.ID == itemCfg.Quality {
				if !p.checkTakeSource(v.ItemID, int64(v.ItemNum)) {
					return errutil.ErrItemNotEnough
				}
				costId = v.ItemID
				costNum = v.ItemNum
				break
			}
		}
	}
	p.takeSource(costId, costNum, constant.ItemDelRefine)

	if err = p.equipRefineRandom(item, attr, random); err != nil {
		return err
	}
	p.lastEquipRefine.itemId = itemUid
	p.lastEquipRefine.random = random
	p.lastEquipRefine.attr = attr

	return nil
}

func (p *Player) equipRefineRandom(item *pb.ItemInfo, attr pb.Attribute, random bool) error {
	itemCfg, err := global.GetItemCfg(item.CfgId)
	if err != nil {
		return err
	}

	createCfg := global.GetRandomCreatCfg(int32(item.Equip.Quality), int32(itemCfg.Levellimit), int32(itemCfg.Equiptype))
	if createCfg == nil {
		return errutil.ErrConfigNotFound
	}

	var weight, distribUid, attrPart int
	var distribCfg *global.RandomDistrib
	if random {
		tmpKindWeight := []global.KeyValue{}
		for _, v := range createCfg.KindNum {
			// 从随机列表中去除装备已有的高级属性，并且不是要洗炼的属性
			var bFind bool
			for _, v2 := range item.Equip.Attributes {
				if v2.Attr == pb.Attribute(v.Key/1000) && attr != pb.Attribute(v.Key) {
					bFind = true
					break
				}
			}
			if !bFind {
				tmpKindWeight = append(tmpKindWeight, global.KeyValue{Key: v.Key, Value: v.Value})
			}
		}

		sort.Slice(tmpKindWeight, func(i, j int) bool {
			return tmpKindWeight[i].Value > tmpKindWeight[j].Value
		})

		weight = global.RandIn(int(tmpKindWeight[0].Value))
		sort.Slice(tmpKindWeight, func(i, j int) bool {
			return tmpKindWeight[i].Value < tmpKindWeight[j].Value
		})

		for _, wv := range tmpKindWeight {
			if weight < int(wv.Value) {
				distribUid = int(wv.Key)
				break
			}
		}

		// 随机分段
		weight = global.RandIn(10000)
		for _, v := range createCfg.RefiningSet {
			if weight < int(v.Value) {
				attrPart = int(v.Key)
				break
			}
		}

		if distribUid == 0 || attrPart == 0 {
			return errutil.ErrConfigNotFound
		}
		distribCfg = global.GetRandomDistribCfg(int32(distribUid), int32(attrPart))
		if distribCfg == nil {
			return errutil.ErrConfigNotFound
		}
	} else {

		weight = global.RandIn(10000)
		for _, v := range createCfg.RefiningSet {
			if weight < int(v.Value) {
				attrPart = int(v.Key)
				break
			}
		}
		distribCfg = global.GetRandomDistribCfgByKind(int32(attr), int32(attrPart))
		if distribCfg == nil {
			return errutil.ErrConfigNotFound
		}
	}

	disCfg, ok := distribCfg.Interval[itemCfg.Quality]
	if !ok {
		return errutil.ErrConfigNotFound
	}
	var statValue int32
	if disCfg.Max <= disCfg.Min {
		statValue = disCfg.Max
	} else {
		statValue = int32(global.RandIn(int(disCfg.Max)-int(disCfg.Min))+1) + disCfg.Min
	}

	p.lastEquipRefine.newStat = pb.Stats{Attr: pb.Attribute(distribCfg.Kind), Num: statValue}
	return nil
}

func (p *Player) equipRefineReplace(replace bool) (*pb.ItemInfo, error) {
	item := p.getItemByUid(p.lastEquipRefine.itemId)
	if item == nil {
		return nil, errutil.ErrItemNotFound
	}
	if p.lastEquipRefine.itemId == 0 || p.lastEquipRefine.newStat.Attr == 0 ||
		p.lastEquipRefine.newStat.Num == 0 || p.lastEquipRefine.attr == 0 {
		return nil, errutil.ErrItemNotFound
	}

	logger.Debugf("最后于一次洗炼属性 类型: %d 得到的新属性 类型:%d 值: %d", p.lastEquipRefine.attr, p.lastEquipRefine.newStat.Attr, p.lastEquipRefine.newStat.Num)
	if replace {
		newStats := []*pb.Stats{}
		for _, v := range item.Equip.Attributes {
			if v.Attr == p.lastEquipRefine.attr {
				logger.Debugf("去除原有属性与最后洗炼属性相同类型: %d ", v.Attr)
				continue
			}
			newStats = append(newStats, v)
		}
		newStats = append(newStats, &p.lastEquipRefine.newStat)

		item.Equip.Attributes = item.Equip.Attributes[:0]
		for _, attr := range newStats {
			item.Equip.Attributes = append(item.Equip.Attributes, &pb.Stats{Attr: attr.Attr, Num: attr.Num})
			logger.Debugf("添加新属性 类型:%d  值:%d", attr.Attr, attr.Num)
		}
		p.refreshStats()
	}

	p.resetRefineInfo()
	return item, nil
}

func (p *Player) resetRefineInfo() {
	p.lastEquipRefine.attr = 0
	p.lastEquipRefine.itemId = 0
	p.lastEquipRefine.newStat = pb.Stats{}
	p.lastEquipRefine.random = false
}

// 穿脱技能opType 1穿 ，2脱
func (p *Player) skillChange(opType int32, slot pb.SkillsSlot, Uuid int32) error {
	bFind := false
	res := p.getSkill(Uuid)
	if res == nil {
		for _, v := range p.partnerSlots {
			partner := p.getPartnerInfoBySlotId(v.SlotId)
			if partner == nil {
				continue
			}
			for _, v2 := range partner.skills {
				if v2.CfgId == Uuid {
					bFind = true
					break
				}
			}
		}
	} else {
		bFind = true
	}
	if !bFind {
		return errutil.ErrChooseRole
	}

	if slot < pb.SkillsSlot_INVALID_SkillsSlot || slot > pb.SkillsSlot_MaxSkillsSlot {
		return errutil.ErrSkillsSlot
	}
	SkillInfo, ok := global.GetSkillCfg(Uuid)
	if ok != nil {
		return errutil.ErrConfigNotFound
	}
	if SkillInfo.SType != 3 && slot == 0 {
		return errutil.ErrSkillsSlotSpecial
	}
	if opType == 1 {
		p.skillsSlot[slot] = Uuid
	} else if opType == 2 {
		p.skillsSlot[slot] = 0
	}
	return nil
}

func (p *Player) getSlotSkillByType(sType int32) (*PlayerSkill, *MapObject) {
	for _, v := range p.skillsSlot {
		if v/100 == sType {
			if skill := p.getSkill(v); skill != nil {
				return skill, p.MapObject
			}

			for _, v2 := range p.partnerSlots {
				partner := p.getPartnerInfoBySlotId(v2.SlotId)
				if partner == nil {
					continue
				}
				if skill := partner.getSkill(v); skill != nil {
					return skill, partner.MapObject
				}
			}
		}
	}
	return nil, nil
}

// 升级技能
func (p *Player) skillsUpgrade(id int32, paruid int64) (int32, error) {
	// 当前技能等级
	// 查询当前槽位的信息
	var skill_proxy *PlayerSkill
	if paruid == 0 {
		for i, v := range p.skills {
			if v.CfgId == id {
				skill_proxy = p.skills[i]
				break
			}
		}
	} else {
		partner := p.getPartnerInfo(paruid)
		if partner == nil {
			return -1, errutil.ErrChooseRole
		}
		for i, v2 := range partner.skills {
			if v2.CfgId == id {
				skill_proxy = partner.skills[i]
				break
			}
		}
	}

	if skill_proxy == nil {
		return -1, errutil.ErrChooseRole
	}

	maxLevel := int32(0)
	skillCfg := global.SkillCfg
	// 目前升级所需要的信息
	SkillInfo, ok := global.GetSkillCfg(id)
	if ok != nil {
		return id, errutil.ErrConfigNotFound
	}
	for _, v := range skillCfg {
		if SkillInfo.SkillProfession == v.SkillProfession && SkillInfo.SkillID == v.SkillID {
			maxLevel++
		}
	}
	// 最高等级不让操作
	if SkillInfo.Level >= maxLevel {
		return id, errutil.ErrMaxLevel
	}
	// 校验背包情况
	// 目前强化所需要的材料
	for _, v := range SkillInfo.PromoteConsume {
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return id, errutil.ErrItemNotEnough
		}
	}
	// 减少道具
	cnum := int32(0)
	pnum := int32(0)
	for _, v := range SkillInfo.PromoteConsume {
		if v.ID == int32(constant.ItemTypeCoin) {
			cnum += int32(v.Num)
		} else {
			pnum += int32(v.Num)
		}
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelSkillUpgrade)
	}
	// 升级
	skill_proxy.CfgId++

	// 升级之后uuid会变，所以要把槽位里面的信息也同步更新
	for i, v := range p.skillsSlot {
		if v == id {
			p.skillsSlot[i] = skill_proxy.CfgId
			break
		}
	}
	level := SkillInfo.Level
	// 更新战斗力
	oldCombat := p.totalCombat
	p.refreshStats()
	newCombat := p.totalCombat
	p.tasks.dispatch(pb.TargetType_SKILL_UPGRADING_TIMES, 0, 1)

	p.upSkillShushuDate(SkillInfo, cnum, pnum, level, newCombat-oldCombat)
	return skill_proxy.CfgId, nil
}

func (p *Player) upSkillShushuDate(SkillInfo *global.SkillConfig, cnum int32, pnum int32, level int32, combat int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	// 在合并点其他的数据
	shushuDate["name"] = "skill_level_up"
	// 新增字段
	shushuDate["skill_id"] = SkillInfo.UID
	shushuDate["skill_name"] = SkillInfo.Name
	shushuDate["cost_item_num"] = pnum //
	shushuDate["cost_gold"] = cnum     // 没有消耗金币
	// shushuDate["skill_rank"] = genstoneCfg.Quality
	shushuDate["level_before"] = level    // 没有等级变化
	shushuDate["level_after"] = level + 1 // 没有等级变化
	shushuDate["power_up"] = combat
	shushumain.ReportData(shushuDate)
}

// 进阶技能
func (p *Player) skillsAdvanced(id int32, paruid int64) (int32, error) {
	// 技能id
	skillID := id / 100

	var skill_proxy *PlayerSkill
	if paruid == 0 {
		for i, v := range p.skills {
			if v.CfgId == id {
				skill_proxy = p.skills[i]
				break
			}
		}
	} else {
		partner := p.getPartnerInfo(paruid)
		if partner == nil {
			return -1, errutil.ErrChooseRole
		}
		for i, v2 := range partner.skills {
			if v2.CfgId == id {
				skill_proxy = partner.skills[i]
				break
			}
		}
	}

	// 目前的进阶等级
	advancedLevel := skill_proxy.AdvancedLevel

	// 目前进阶最大的等级
	SkillInfo, ok := global.GetSkillCfg(id)
	if ok != nil {
		return advancedLevel, errutil.ErrConfigNotFound
	}
	// 最高等级不让操作
	if advancedLevel >= SkillInfo.CanupperLevel {
		return advancedLevel, errutil.ErrMaxLevel
	}
	// 查询下一进阶需要的材料
	skillCfg := global.SkillCfg
	needConfig := global.SkillConfig{}
	for _, v := range skillCfg {
		if SkillInfo.SkillProfession == v.SkillProfession && advancedLevel+1 == v.CanupperLevel && skillID == v.SkillID && len(v.UpperlevelConsume) > 0 {
			needConfig = global.SkillConfig{
				UpperlevelConsume: v.UpperlevelConsume,
			}
			break
		}
	}

	// 校验背包情况
	// 目前进阶所需要的材料
	for _, v := range needConfig.UpperlevelConsume {
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return advancedLevel, errutil.ErrItemNotEnough
		}
	}
	// 减少道具
	cnum := int32(0)
	pnum := int32(0)
	for _, v := range needConfig.UpperlevelConsume {
		if v.ID == int32(constant.ItemTypeCoin) {
			cnum += int32(v.Num)
		} else {
			pnum += int32(v.Num)
		}
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelSkillUpgrade)
	}
	// 升级
	skill_proxy.AdvancedLevel++
	// 更新战斗力
	oldCombat := p.totalCombat
	p.refreshStats()
	newCombat := p.totalCombat
	p.refreshStats()
	p.upSkillGradeShushuDate(SkillInfo, cnum, pnum, skill_proxy.AdvancedLevel, newCombat-oldCombat)
	p.tasks.dispatch(pb.TargetType_SKILL_ADVANCED_TIMES, 0, 1)
	return skill_proxy.AdvancedLevel, nil
}

func (p *Player) upSkillGradeShushuDate(SkillInfo *global.SkillConfig, cnum int32, pnum int32, level int32, combat int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "skill_upgrade"
	shushuDate["skill_id"] = SkillInfo.UID
	shushuDate["skill_name"] = SkillInfo.Name
	shushuDate["cost_item_num"] = pnum
	shushuDate["cost_gold"] = cnum
	shushuDate["rank_before"] = level - 1
	shushuDate["rank_after"] = level
	shushuDate["power_up"] = combat
	shushuDate["skill_level"] = SkillInfo.Level
	shushumain.ReportData(shushuDate)
}

// 爬塔收益
func (p *Player) soloTowerOfflineReward() error {
	cfg := global.GetEndlesstrialsCfg(p.SoloTowerInfo.Layers)
	if cfg == nil {
		return errutil.ErrConfigNotFound
	}
	offsetMinite := (time.Now().Unix() - p.SoloTowerInfo.RewardTime.Unix())
	reward := []global.IdNum{}

	for _, v := range cfg.OfflineReward {
		multiple := offsetMinite / int64(v.Gap)
		if multiple == 0 {
			continue
		}
		num := int64(v.Num) * multiple
		if num > int64(v.Max) {
			num = int64(v.Max)
		}

		reward = append(reward, global.IdNum{ID: v.ID, Num: int32(num)})

		if !p.checkAddSource(v.ID, num) {
			return errutil.ErrBagFull
		}
	}
	sourceInfo := []global.IdNum{}
	for _, v := range reward {
		sourceInfo = append(sourceInfo, global.IdNum{ID: v.ID, Num: v.Num})
		// p.addSource(v.ID, int64(v.Num), constant.ItemAddSoloTower)
	}
	p.addSourceBatch(sourceInfo, constant.ItemAddSoloTower)
	p.SoloTowerInfo.RewardTime = time.Now()

	return nil
}

// 是否是升星石
func (p *Player) IsStarsLevel(CfgId int32, Num int32) (int32, error) {
	Starlevel := global.GlobalCfg.Starstone.Starlevel
	checkId := int32(-1)
	slot := int32(0)
	for _, v := range Starlevel {
		if v.ID == CfgId {
			checkId = v.Level
			slot = v.Slot
			break
		}
	}
	if checkId != -1 {
		for i := 1; i <= int(Num); i++ {
			if slot != 0 {
				if p.equips[slot].Stars < checkId {
					p.equips[slot].Stars = checkId
				} else {
					return int32(i - 1), errutil.StarsNoUseError
				}
			} else {
				selectId := make([]int, 0)
				for i, v := range p.equips {
					// 需要随机给他一个槽位升满
					if v.Stars < checkId && i != 0 {
						selectId = append(selectId, i)
					}
				}
				if len(selectId) == 0 {
					return int32(i - 1), errutil.StarsNoUseError
				}
				weight := global.RandIn(len(selectId))
				p.equips[selectId[weight]].Stars = checkId
			}
		}
		p.SyncUserEquipsInfo()
		p.refreshStats()

	}
	return Num, nil
}

// 是否是命格石头
func (p *Player) IsFate(CfgId int32, Num int32) ([]int, int32, error) {
	FateInfo := global.GlobalCfg.Fate.Basevalue.Item
	checkId := int32(-1)
	slot := int64(0)
	for _, v := range FateInfo {
		if v.ID == CfgId {
			checkId = v.ID
			slot = int64(v.Num)
			break
		}
	}
	perAll := make([]int, 0)
	if checkId != -1 {
		for i := 1; i <= int(Num); i++ {
			fatecfg, err := global.GetFateCfg(p.fateLevel)
			if err != nil {
				return perAll, int32(i), errutil.ErrConfigNotFound
			}
			per := fateStageRand()
			perAll = append(perAll, per)
			EndExp := p.fateExp + slot*int64(per)
			FateCfg := global.FateCfg
			// 查询这一重需要多少经验
			AllExp := fatecfg.FateExp
			for j := 1; j <= len(FateCfg)+1; j++ {
				fatecfg, err := global.GetFateCfg(p.fateLevel - int32(j))
				if err != nil || len(fatecfg.AdvancedNeed) > 0 {
					break
				}
				AllExp = AllExp + fatecfg.FateExp
			}
			if p.fateExp >= int64(AllExp) {
				if len(fatecfg.AdvancedNeed) > 0 {
					p.refreshStats()
					p.syncUserBase()
					return perAll, int32(i), errutil.ErrMaxLevel
				} else {
					p.fateLevel++
					p.fateExp = EndExp

				}
			} else {
				p.fateExp = EndExp
			}
		}
		p.refreshStats()
		p.syncUserBase()
	}
	return perAll, Num, nil
}

// 是否是命格石头
func (p *Player) IsPartner(CfgId int32, Num int32) ([]int, int32, error) {
	FateInfo := global.GlobalCfg.Fate.Basevalue.Item
	checkId := int32(-1)
	slot := int64(0)
	for _, v := range FateInfo {
		if v.ID == CfgId {
			checkId = v.ID
			slot = int64(v.Num)
			break
		}
	}
	perAll := make([]int, 0)
	if checkId != -1 {
		for i := 1; i <= int(Num); i++ {
			fatecfg, err := global.GetFateCfg(p.fateLevel)
			if err != nil {
				return perAll, int32(i), errutil.ErrConfigNotFound
			}
			per := fateStageRand()
			perAll = append(perAll, per)
			EndExp := p.fateExp + slot*int64(per)
			FateCfg := global.FateCfg
			// 查询这一重需要多少经验
			AllExp := fatecfg.FateExp
			for j := 1; j <= len(FateCfg)+1; j++ {
				fatecfg, err := global.GetFateCfg(p.fateLevel - int32(j))
				if err != nil || len(fatecfg.AdvancedNeed) > 0 {
					break
				}
				AllExp = AllExp + fatecfg.FateExp
			}
			if p.fateExp >= int64(AllExp) {
				if len(fatecfg.AdvancedNeed) > 0 {
					p.refreshStats()
					p.syncUserBase()
					return perAll, int32(i), errutil.ErrMaxLevel
				} else {
					p.fateLevel++
					p.fateExp = EndExp

				}
			} else {
				p.fateExp = EndExp
			}
		}
		p.refreshStats()
		p.syncUserBase()
	}
	return perAll, Num, nil
}

// 命格暴击
func fateStageRand() int {
	fateStage := global.GlobalCfg.Fate.Rate.Stage
	allPriority := 0
	beginPriority := 0
	// 活动内容
	for _, fateStageInfo := range fateStage {
		allPriority = allPriority + int(fateStageInfo.Pro)
	}
	rand := global.RandIn(allPriority) + 1
	for _, fateStageInfo := range fateStage {
		if beginPriority+int(fateStageInfo.Pro) >= rand {
			return fateStageInfo.Rate
		}
		beginPriority = beginPriority + int(fateStageInfo.Pro)
	}
	return -1
}

func (p *Player) SetPlayerInfo(fateLevel int32, fateExp int64) {
	// 升级
	p.fateLevel = fateLevel
	p.fateExp = fateExp
	// 更新战斗力
	p.refreshStats()
}

func (p *Player) SendMail(source string, attachment []protocol.Attachment, args ...string) {
	mail := protocol.MailAddRequest{
		Uid:          -1,
		RecipientUid: p.uid,
		Title:        constant.EmailContentSystem,

		ChannelCode: source,
		Attachment:  attachment,
	}

	if len(args) > 0 {
		mail.Desc = args[0]
	} else {
		mail.Desc = constant.EmailContentBagFull
	}

	if p.mailCounter >= int(global.GlobalCfg.BaseValue.Mailmax.Num) {
		packet := &pb.SystemNotificationPush{Type: pb.Notification_MAIL_FULL}
		p.pushMsg(protocol.RouteSystemNotificationPush, packet)
		return
	}

	p.SetMailCounter(constant.MailCounterOPTypeAdd, 1)

	if err := provide.Mailadd(&mail); err != nil {
		packet := &pb.SystemNotificationPush{Type: pb.Notification_MAIL_FULL}
		p.pushMsg(protocol.RouteSystemNotificationPush, packet)
	}
}

// 合成系统
func (p *Player) compose(composeInfo []*pb.ComposeInfo) ([]*pb.ComposeInfo, error) {
	whileAll := make(map[int32]int32)
	for _, v := range global.ComposeCfg {
		for _, v2 := range v.Raw {
			whileAll[v2] = v.ComposeID
		}
	}
	if len(composeInfo) == 0 {
		return nil, errutil.ErrItemNotFound
	}
	if composeInfo[0] == nil {
		return nil, errutil.ErrItemNotFound
	}
	firstRes := p.getItemByUid(composeInfo[0].Uid)
	if firstRes == nil {
		return nil, errutil.ErrItemNotFound
	}
	// 可堆叠数量大于1的都用商品id
	// 判断是不是能穿戴的装备
	ItemCfg, err := global.GetItemCfg(firstRes.CfgId)
	if err != nil {
		return nil, err
	}
	if ItemCfg.Maxnum > 1 {
		beginComposeId := whileAll[firstRes.CfgId]
		res, ok := global.GetComposeCfg(beginComposeId)
		if ok != nil {
			return nil, ok
		}

		change := true
		sourceDel := make([]global.IdNum, 0)
		sourceLen := 0
		allItemInfo := []*pb.ItemInfo{}
		for _, v := range composeInfo {
			itemInfo := p.getItemByUid(v.Uid)
			if itemInfo == nil {
				return nil, errutil.ErrItemNotFound
			}
			allItemInfo = append(allItemInfo, itemInfo)
			if !p.checkTakeSourceByUid(v.Uid, int64(v.Num)) {
				return nil, errutil.ErrItemNotEnough
			}
			// 查询是不是都是可以合成的 且都是一个种类
			if whileAll[itemInfo.CfgId] == 0 || beginComposeId != whileAll[itemInfo.CfgId] {
				change = false
			}
			sourceLen = sourceLen + int(v.Num)
			sourceDel = append(sourceDel, global.IdNum{
				ID:  int32(itemInfo.CfgId),
				Num: int32(v.Num),
			})
		}
		if sourceLen%int(res.NeedNum) != 0 {
			return nil, errutil.ErrItemNotEnough
		}
		if !change {
			return nil, errutil.ErrItemNotEnough
		}

		needGoldPer := int(res.NeedGold.Num) * sourceLen / int(res.NeedNum)
		// 校验背包情况
		// 目前合成所需要的必要材料
		if !p.checkTakeSource(res.NeedGold.ID, int64(needGoldPer)) {
			return nil, errutil.ErrItemNotEnough
		}
		// 删除合成需要的材料金币
		p.takeSource(res.NeedGold.ID, int64(needGoldPer), constant.ItemDelGemstone)
		// 删除合成需要的物品
		for _, v := range sourceDel {
			if !p.checkTakeSource(v.ID, int64(v.Num)) {
				return nil, errutil.ErrItemNotEnough
			}
		}
		for _, v := range sourceDel {
			p.takeSource(v.ID, int64(v.Num), constant.ItemDelCompose)
		}
		sourceAddReturn := make([]*pb.ComposeInfo, 0)
		// 随机分段
		whileSend := make(map[int32]int32)
		for i := 1; i <= sourceLen/int(res.NeedNum); i++ {
			weight := global.RandIn(int(BaseRate))
			if float64(weight) < float64(res.SuccessRate) {
				randId := composeRand(res)
				whileSend[randId]++
			}
		}
		var sourceAdd []global.IdNum
		if len(whileSend) > 0 {
			for id, num := range whileSend {
				sourceAdd = append(sourceAdd, global.IdNum{
					ID:  id,
					Num: num,
				})
				sourceAddReturn = append(sourceAddReturn, &pb.ComposeInfo{
					Uid: int64(id),
					Num: num,
				})
			}
			if err := p.checkAddSourceBatch(sourceAdd); err == nil {
				p.addSourceBatch(sourceAdd, constant.ItemAddCompose)
			} else {
				attachment := []protocol.Attachment{}
				for _, v := range sourceAdd {
					attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
				}
				p.SendMail(constant.ItemAddElementSweep, attachment)
			}
			for _, info := range sourceAdd {
				itemInfo, err := global.GetItemCfg(info.ID)
				if err != nil {
					return nil, err
				}
				if itemInfo.Quality > 4 && itemInfo.NatureMark > 0 {
					strSend := []string{p.name, strconv.Itoa(int(itemInfo.ID))}
					defaultManager.Lamp(int32(constant.LampGemSynthesis), p.serverID, strSend)
				}

			}
			p.upComposeItemShushuDate(allItemInfo, sourceAdd, sourceLen, needGoldPer)
			return sourceAddReturn, nil
		} else {
			return nil, errutil.ComposeError
		}

	} else {
		beginComposeId := whileAll[firstRes.CfgId]
		res, ok := global.GetComposeCfg(beginComposeId)
		if ok != nil {
			return nil, ok
		}

		change := true
		sourceDel := make([]global.UUIdNum, 0)
		sourceLen := 0
		allItemInfo := []*pb.ItemInfo{}
		for _, v := range composeInfo {
			itemInfo := p.getItemByUid(v.Uid)
			allItemInfo = append(allItemInfo, itemInfo)
			if itemInfo == nil {
				return nil, errutil.ErrItemNotFound
			}
			if !p.checkTakeSourceByUid(v.Uid, int64(v.Num)) {
				return nil, errutil.ErrItemNotEnough
			}
			// 查询是不是都是可以合成的 且都是一个种类
			if whileAll[itemInfo.CfgId] == 0 || beginComposeId != whileAll[itemInfo.CfgId] {
				change = false
			}
			sourceLen = sourceLen + int(v.Num)
			sourceDel = append(sourceDel, global.UUIdNum{
				UID: v.Uid,
				Num: int64(v.Num),
			})
		}
		if sourceLen%int(res.NeedNum) != 0 {
			return nil, errutil.ErrItemNotEnough
		}
		if !change {
			return nil, errutil.ErrItemNotEnough
		}

		// 校验背包情况
		// 目前合成所需要的必要材料
		if !p.checkTakeSource(res.NeedGold.ID, int64(res.NeedGold.Num)) {
			return nil, errutil.ErrItemNotEnough
		}
		// 删除合成需要的材料金币
		p.takeSource(res.NeedGold.ID, int64(res.NeedGold.Num), constant.ItemDelGemstone)
		// 删除合成需要的物品
		p.takeSourceByUidBatch(sourceDel, constant.ItemDelCompose)
		sourceAddReturn := make([]*pb.ComposeInfo, 0)
		// 随机分段
		weight := global.RandIn(int(BaseRate))
		if float64(weight) < float64(res.SuccessRate) {
			whileSend := make(map[int32]int32)
			for i := 1; i <= sourceLen/int(res.NeedNum); i++ {
				randId := composeRand(res)
				whileSend[randId]++

			}
			var sourceAdd []global.IdNum
			for id, num := range whileSend {
				sourceAdd = append(sourceAdd, global.IdNum{
					ID:  id,
					Num: num,
				})
				sourceAddReturn = append(sourceAddReturn, &pb.ComposeInfo{
					Uid: int64(id),
					Num: num,
				})
			}
			if err := p.checkAddSourceBatch(sourceAdd); err == nil {
				p.addSourceBatch(sourceAdd, constant.ItemAddCompose)
			} else {
				attachment := []protocol.Attachment{}
				for _, v := range sourceAdd {
					attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
				}
				p.SendMail(constant.ItemAddElementSweep, attachment)
			}
			p.upComposeItemShushuDate(allItemInfo, sourceAdd, sourceLen, int(res.NeedGold.Num))

			return sourceAddReturn, nil
		} else {
			return nil, errutil.ComposeError
		}
	}

}
func (p *Player) upComposeItemShushuDate(composeInfo []*pb.ItemInfo, sourceAdd []global.IdNum, synthesis_price int, cost_gold int) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "item_synthesis"

	shushuDate["is_three"] = false                  // 是否三合一，按实际赋值
	shushuDate["synthesis_price"] = synthesis_price // 是否三合一，按实际赋值
	cost_item := make([]map[string]interface{}, 0)
	for _, v := range composeInfo {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(int32(v.CfgId)); err == nil {
				GeTNum, err := defaultManager.GetItemCount(p.uid, int32(v.Uid))
				if err != nil {
					continue
				}
				cost_item = append(cost_item, map[string]interface{}{
					"cost_item_id":     ItemCfg.ID,
					"cost_item_name":   ItemCfg.Name,
					"cost_item_num":    v.Num,
					"cost_item_before": GeTNum - int64(v.Num),
					"cost_item_after":  GeTNum,
				})
			}
		}
	}
	shushuDate["cost_item"] = cost_item
	shushuDate["cost_gold"] = cost_gold
	synthesis_item := make([]map[string]interface{}, 0)
	for _, v := range sourceAdd {
		if v.Num > 0 {
			// 查询itemname
			if ItemCfg, err := global.GetItemCfg(v.ID); err == nil {
				GeTNum, err := defaultManager.GetItemCount(p.uid, v.ID)
				if err != nil {
					continue
				}
				synthesis_item = append(synthesis_item, map[string]interface{}{
					"synthesis_item_id":     ItemCfg.ID,
					"synthesis_item_name":   ItemCfg.Name,
					"synthesis_item_num":    v.Num,
					"synthesis_item_before": GeTNum,
					"synthesis_item_after":  GeTNum + int64(v.Num),
				})
			}
		}
	}
	shushuDate["synthesis_item"] = synthesis_item
	shushumain.ReportData(shushuDate)
}

// 合成概率分布
func composeRand(ComposeCfg *global.ComposeConfig) int32 {
	fateStage := ComposeCfg.Produced
	allPriority := 0
	beginPriority := 0
	// 活动内容
	for _, fateStageInfo := range fateStage {
		allPriority = allPriority + int(fateStageInfo.Ratio)
	}
	rand := global.RandIn(allPriority) + 1
	for _, fateStageInfo := range fateStage {
		if beginPriority+int(fateStageInfo.Ratio) >= rand {
			ids := fateStageInfo.ID
			idRand := global.RandIn(len(ids))
			return ids[idRand]
		}
		beginPriority = beginPriority + int(fateStageInfo.Ratio)
	}
	return -1
}

// 穿脱宝石opType 1穿 ，2脱
func (p *Player) genstoneChange(opType int32, slot pb.EquipSlot, stoneSlot pb.GemstoneSlot, itemUid, partnerslotid int32) error {
	if stoneSlot < pb.GemstoneSlot_INVALID_GemstoneSlot || stoneSlot > pb.GemstoneSlot_MaxGemstoneSlot {
		return errutil.ErrGenStoneSlot
	}

	var equipsInfo *Slot
	if partnerslotid == 0 {
		equipsInfo = &p.equips[pb.EquipSlot(slot)]
	} else {
		equipsInfo = &p.getPartnerSlotInfo(partnerslotid).Equips[pb.EquipSlot(slot)]
	}
	if equipsInfo == nil {
		return errutil.ErrEquipSlot
	}

	if opType == 1 {
		canSlot := 0
		for _, StageInfo := range global.GlobalCfg.Diamond.Openlevel.Stage {
			// intLevel, _ := strconv.Atoi(StageInfo.Level)
			// if int(p.level) >= intLevel {
			// 	delInfo, _ := strconv.Atoi(StageInfo.Num)
			// 	if delInfo > canSlot {
			// 		canSlot = delInfo
			// 	}
			// }
			if p.level >= StageInfo.Level {
				if int(StageInfo.Num) > canSlot {
					canSlot = int(StageInfo.Num)
				}
			}
		}
		if stoneSlot >= pb.GemstoneSlot(canSlot) {
			return errutil.ErrGenStoneSlot
		}
		// 校验背包情况
		// 目前合成所需要的必要材料
		if !p.checkTakeSource(itemUid, int64(1)) {
			return errutil.ErrItemNotEnough
		}
		// 判断是不是能穿戴的装备
		genstoneCfg, err := global.GetItemCfg(itemUid)
		if err != nil {
			return errutil.ErrItemNotEnough
		}
		// 魂石大师
		gemLevelOld := int32(0)
		gemLevelNew := int32(0)
		allGems := []int32{}
		for slot, equip := range p.equips {
			if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
				continue
			}
			for GemstoneSlot, Gemstone := range equip.Gemstone {
				if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
					continue
				}
				genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
				if err != nil {
					continue
				}
				// 对应宝石等级的数目加1
				allGems = append(allGems, genstoneCfg.Itemlevel)
			}
		}
		sort.Slice(allGems, func(i, j int) bool {
			return allGems[i] > allGems[j]
		})
		for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
			if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
				equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
				if err != nil {
					continue
				}
				if allGems[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
					gemLevelOld = equipMasterCfg.Limit
					break
				}
			}
		}

		can := false
		for _, stoneSlotInfo := range genstoneCfg.DiamondPosition {
			if stoneSlotInfo == int(pb.EquipSlot(slot)) {
				can = true
				continue
			}
		}
		if !can {
			return errutil.ErrGenStoneSlot
		}
		if equipsInfo.Gemstone[stoneSlot] > 0 {
			p.addSource(int32(equipsInfo.Gemstone[stoneSlot]), 1, constant.ItemAddGemstone)
		}
		p.takeSource(int32(itemUid), 1, constant.ItemDelGemstone)
		equipsInfo.Gemstone[stoneSlot] = pb.GemstoneSlot(itemUid)
		oldcombat := p.totalCombat
		p.refreshStats()
		allGems2 := []int32{}
		for slot, equip := range p.equips {
			if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
				continue
			}
			for GemstoneSlot, Gemstone := range equip.Gemstone {
				if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
					continue
				}
				genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
				if err != nil {
					continue
				}
				// 对应宝石等级的数目加1
				allGems2 = append(allGems2, genstoneCfg.Itemlevel)
			}
		}
		sort.Slice(allGems2, func(i, j int) bool {
			return allGems2[i] > allGems2[j]
		})
		for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
			if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems2) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
				equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
				if err != nil {
					continue
				}
				if allGems2[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
					gemLevelNew = equipMasterCfg.Limit
					break
				}
			}
		}
		if gemLevelNew > gemLevelOld {
			p.upSuitShushuDate(constant.EquipMasterCategoryGemStone, gemLevelNew, partnerslotid)
		}

		newcombat := p.totalCombat
		// 查询装备信息
		var equipCfg *global.ItemConfig
		equip_quality := int32(0)
		if equipsInfo.ItemId != 0 {
			equipCfg2 := p.getItemByUid(equipsInfo.ItemId)
			if equipCfg2 == nil {
				return errutil.ErrItemNotEnough
			}
			equip_quality = int32(equipCfg2.Equip.Quality)
			equipCfg, err = global.GetItemCfg(int32(equipCfg2.CfgId))
			if err != nil {
				return errutil.ErrItemNotEnough
			}
		}
		p.upStoneShushuDate(equipCfg, genstoneCfg, equipsInfo.Enhance, int32(slot), newcombat-oldcombat, equip_quality)
		p.tasks.dispatch(pb.TargetType_GEMSTONE_INLAY_TIMES, 0, 1)
		p.rankChangeTCP(constant.RankTypeSoulStoneLevel, p.calcSoulStoneLevelNum())
		return nil

	} else if opType == 2 {
		if equipsInfo.Gemstone[stoneSlot] == 0 {
			return nil
		}
		if !p.checkAddSource(int32(equipsInfo.Gemstone[stoneSlot]), 1) {
			return errutil.ErrBagFull
		}
		p.addSource(int32(equipsInfo.Gemstone[stoneSlot]), 1, constant.ItemAddGemstone)
		equipsInfo.Gemstone[stoneSlot] = 0
		p.refreshStats()
		return nil
	}
	return nil
}

func (p *Player) upStoneShushuDate(equipsDetail *global.ItemConfig, genstoneCfg *global.ItemConfig, level int32, equip_part int32, combat int32, equip_quality int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "equip_inlay"
	if equipsDetail == nil {
		shushuDate["equip_id"] = 0
		shushuDate["equip_num"] = 0
		shushuDate["equip_quality"] = 0
		shushuDate["equip_name"] = 0
		shushuDate["equip_job"] = 0
	} else {
		shushuDate["equip_id"] = equipsDetail.ID
		shushuDate["equip_num"] = 1
		shushuDate["equip_quality"] = equip_quality
		shushuDate["equip_name"] = equipsDetail.Name
		shushuDate["equip_job"] = equipsDetail.Occupylimit
	}

	shushuDate["equip_part"] = equip_part
	shushuDate["equip_level"] = level

	// soul_stone_info 字段为嵌套map
	re := regexp.MustCompile(`(\d+)级(.)`)
	matches := re.FindStringSubmatch(genstoneCfg.Name)
	shushuDate["soul_stone_info"] = map[string]interface{}{
		"soul_stone_id":        genstoneCfg.ID,         // 宝石ID
		"soul_stone_name":      genstoneCfg.Name,       // 没有名称可填空字符串
		"soul_stone_level":     genstoneCfg.Itemlevel,  // 没有等级可填1
		"soul_stone_type":      matches[2],             // 没有类型可填0
		"soul_stone_attribute": genstoneCfg.NatureMark, //
	}

	// 正则：匹配“9级红魂石(雷)”这种，数字后面是“级”，再跟一个汉字
	// 示例正则: `(\d+)级(.)`
	// 说明: (\d+) 匹配数字，级为固定字，(.) 匹配一个汉字或任意字符
	// 用法示例:

	shushuDate["equip_power_up"] = combat
	shushumain.ReportData(shushuDate)
}

// 一键穿脱宝石opType 1穿 ，2脱
func (p *Player) genstoneChangeAll(slot pb.EquipSlot, partnerslotid int32) ([]int32, error) {
	var equipsInfo *Slot
	if partnerslotid == 0 {
		equipsInfo = &p.equips[pb.EquipSlot(slot)]
	} else {
		equipsInfo = &p.getPartnerSlotInfo(partnerslotid).Equips[pb.EquipSlot(slot)]
	}
	if equipsInfo == nil {
		return nil, errutil.ErrEquipSlot
	}
	// 魂石大师
	gemLevelOld := int32(0)
	gemLevelNew := int32(0)
	allGems := []int32{}
	for slot, equip := range p.equips {
		if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
			continue
		}
		for GemstoneSlot, Gemstone := range equip.Gemstone {
			if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
				continue
			}
			genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
			if err != nil {
				continue
			}
			// 对应宝石等级的数目加1
			allGems = append(allGems, genstoneCfg.Itemlevel)
		}
	}
	sort.Slice(allGems, func(i, j int) bool {
		return allGems[i] > allGems[j]
	})
	for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
		if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
			equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
			if err != nil {
				continue
			}
			if allGems[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
				gemLevelOld = equipMasterCfg.Limit
				break
			}
		}
	}
	canDell := make([]global.IdKeyNum, 0)
	// for k1, _ := range genstoneCfgAll {
	for _, itemInfo := range p.bag.ItemList {
		genstoneCfgAll, err := global.GetItemCfg(itemInfo.CfgId)
		if err != nil {
			return nil, errutil.ErrItemNotEnough
		}
		if len(genstoneCfgAll.DiamondPosition) > 0 {
			whileSolt := false
			for _, diamond := range genstoneCfgAll.DiamondPosition {
				if diamond == int(slot) {
					whileSolt = true
				}
			}
			if !whileSolt {
				continue
			}
			// 1是背包里面的
			for i := 0; i < int(itemInfo.Num); i++ {
				canDell = append(canDell, global.IdKeyNum{
					ID:  itemInfo.CfgId,
					Num: genstoneCfgAll.Itemlevel,
					Key: 1,
				})
			}
		}
	}
	canSlot := 0
	for _, StageInfo := range global.GlobalCfg.Diamond.Openlevel.Stage {
		// intLevel, _ := strconv.Atoi(StageInfo.Level)
		// if int(p.level) >= intLevel {
		// 	delInfo, _ := strconv.Atoi(StageInfo.Num)
		// 	if delInfo > canSlot {
		// 		canSlot = delInfo
		// 	}
		// }
		if p.level >= StageInfo.Level {
			if int(StageInfo.Num) > canSlot {
				canSlot = int(StageInfo.Num)
			}
		}
	}
	// old := make(map[int]IdKeyNum{}, 0)
	old := make([]global.IdKeyNum, 0)
	for _, stoneSlotInfo := range equipsInfo.Gemstone {
		// genstoneCfgInfo, _ := global.GetItemCfg(int32(stoneSlotInfo))
		if stoneSlotInfo > 0 {
			genstoneCfgAll, err := global.GetItemCfg(int32(stoneSlotInfo))
			if err != nil {
				return nil, errutil.ErrItemNotEnough
			}
			// 2是身上的
			old = append(old, global.IdKeyNum{
				ID:  int32(stoneSlotInfo),
				Num: genstoneCfgAll.Itemlevel,
			})
		} else {
			old = append(old, global.IdKeyNum{
				ID:  0,
				Num: 0,
			})
		}
		// old[i] = int32(stoneSlotInfo)
	}
	// newBag := topFourByField(MaxBagStone, canSlot, 1)
	// returnDate := make([]int32, 0, len(newBag))

	// for i, oldInfo := range old {
	// 	for _, newInfo := range newBag {
	// 		// 1是背包，2是身上
	// 		if newInfo.Key == 2 {
	// 			if int32(oldInfo) == newInfo.ID {
	// 				delete(old, i)
	// 				break
	// 			}
	// 		}
	// 	}
	// }
	// 按槽位顺序处理宝石镶嵌
	// for _, newInfo := range newBag {
	// 	// 1是背包，2是身上
	// 	// 跳过相同宝石
	for i := 1; i <= canSlot; i++ {
		newBag := topFourByField(canDell, 1, 1)
		if len(newBag) == 0 {
			continue
		}
		newInfo := newBag[0]
		for j, olditem := range old {
			if newInfo.Num > olditem.Num {
				if int32(olditem.Num) >= newInfo.ID {
					continue
				}
				if j >= canSlot {
					continue
				}
				// 处理旧宝石回包
				if olditem.ID != 0 {
					// 检查背包空间
					if !p.checkAddSource(int32(olditem.ID), 1) {
						return nil, errutil.ErrBagFull
					}
					// 先扣除新宝石
					p.takeSource(newInfo.ID, 1, constant.ItemDelGemstone)
					// 再添加旧宝石
					p.addSource(int32(olditem.ID), 1, constant.ItemAddGemstone)
					for k := 0; k < len(canDell); k++ {
						if canDell[k].ID == newInfo.ID {
							canDell = append(canDell[:k], canDell[k+1:]...)
							break
						}
					}
					canDell = append(canDell, global.IdKeyNum{
						ID:  int32(olditem.ID),
						Num: olditem.Num,
					})

				} else {
					// 直接扣除新宝石
					p.takeSource(newInfo.ID, 1, constant.ItemDelGemstone)
					// canDell = deleteItem(canDell, newInfo.ID)
					for k := 0; k < len(canDell); k++ {
						if canDell[k].ID == newInfo.ID {
							canDell = append(canDell[:k], canDell[k+1:]...)
							break
						}
					}
				}

				// 更新装备槽
				equipsInfo.Gemstone[j] = pb.GemstoneSlot(newInfo.ID)
				// // 删除old这个
				// delete(old, j)
				old[j] = global.IdKeyNum{
					ID:  int32(newInfo.ID),
					Num: newInfo.Num,
				}
				oldcombat := p.totalCombat
				p.refreshStats(WithRefreshStatsPushMessage(false))
				newcombat := p.totalCombat
				genstoneCfg, err := global.GetItemCfg(newInfo.ID)
				if err != nil {
					return nil, errutil.ErrItemNotEnough
				}
				var equipCfg *global.ItemConfig
				equip_quality := int32(0)
				if equipsInfo.ItemId != 0 {
					equipCfg2 := p.getItemByUid(equipsInfo.ItemId)
					if equipCfg2 == nil {
						return nil, errutil.ErrItemNotEnough
					}
					equip_quality = int32(equipCfg2.Equip.Quality)
					equipCfg, err = global.GetItemCfg(int32(equipCfg2.CfgId))
					if err != nil {
						return nil, errutil.ErrItemNotEnough
					}
				}
				p.rankChangeTCP(constant.RankTypeSoulStoneLevel, p.calcSoulStoneLevelNum())
				p.upStoneShushuDate(equipCfg, genstoneCfg, equipsInfo.Enhance, int32(slot), newcombat-oldcombat, equip_quality)
				p.tasks.dispatch(pb.TargetType_GEMSTONE_INLAY_TIMES, 0, 1)
				break
			}
			// 	}
		}
	}
	// for j, olditem := range old {
	// if newInfo.Num > olditem.Num {

	// 	if int32(olditem) == newInfo.ID {
	// 		continue
	// 	}
	// 	if j >= canSlot {
	// 		continue
	// 	}
	// 	// 处理旧宝石回包
	// 	if olditem != 0 {
	// 		// 检查背包空间
	// 		if !p.checkAddSource(int32(olditem), 1) {
	// 			return nil, errutil.ErrBagFull
	// 		}
	// 		// 先扣除新宝石
	// 		p.takeSource(newInfo.ID, 1, constant.ItemDelGemstone)
	// 		// 再添加旧宝石
	// 		p.addSource(int32(olditem), 1, constant.ItemAddGemstone)

	// 	} else {
	// 		// 直接扣除新宝石
	// 		p.takeSource(newInfo.ID, 1, constant.ItemDelGemstone)
	// 	}

	// 	// 更新装备槽
	// 	equipsInfo.Gemstone[j] = pb.GemstoneSlot(newInfo.ID)
	// 	// // 删除old这个
	// 	delete(old, j)
	// 	oldcombat := p.totalCombat
	// 	p.refreshStats(WithRefreshStatsPushMessage(false))
	// 	newcombat := p.totalCombat
	// 	genstoneCfg, err := global.GetItemCfg(newInfo.ID)
	// 	if err != nil {
	// 		return nil, errutil.ErrItemNotEnough
	// 	}
	// 	var equipCfg *global.ItemConfig
	// 	if equipsInfo.ItemId != 0 {
	// 		equipCfg2 := p.getItemByUid(equipsInfo.ItemId)
	// 		if equipCfg2 == nil {
	// 			return nil, errutil.ErrItemNotEnough
	// 		}
	// 		equipCfg, err = global.GetItemCfg(int32(equipCfg2.CfgId))
	// 		if err != nil {
	// 			return nil, errutil.ErrItemNotEnough
	// 		}
	// 	}
	// 	p.upStoneShushuDate(equipCfg, genstoneCfg, equipsInfo.Enhance, int32(slot), newcombat-oldcombat)

	// 	break
	// }
	// 	}
	// }
	// 把equipsInfo.Gemstone[j] 按照num的从大到小排序
	returnDate := make([]int32, 0)
	for _, v := range equipsInfo.Gemstone {
		returnDate = append(returnDate, int32(v))
	}

	p.refreshStats()
	allGems2 := []int32{}
	for slot, equip := range p.equips {
		if slot <= int(pb.EquipSlot_INVALID_EquipSlot) || slot >= int(pb.EquipSlot_MaxEquipSlot) {
			continue
		}
		for GemstoneSlot, Gemstone := range equip.Gemstone {
			if GemstoneSlot < int(pb.GemstoneSlot_INVALID_GemstoneSlot) || GemstoneSlot > int(pb.GemstoneSlot_MaxGemstoneSlot) {
				continue
			}
			genstoneCfg, err := global.GetItemCfg(int32(Gemstone))
			if err != nil {
				continue
			}
			// 对应宝石等级的数目加1
			allGems2 = append(allGems2, genstoneCfg.Itemlevel)
		}
	}
	sort.Slice(allGems2, func(i, j int) bool {
		return allGems2[i] > allGems2[j]
	})
	for i := len(global.GlobalCfg.Diamond.Master.Stage) - 1; i >= 0; i-- {
		if global.GlobalCfg.Diamond.Master.Stage[i].Num > 0 && len(allGems2) >= int(global.GlobalCfg.Diamond.Master.Stage[i].Num) {
			equipMasterCfg, err := global.GetEquipMaster(global.GlobalCfg.Diamond.Master.Stage[i].ID)
			if err != nil {
				continue
			}
			if allGems2[global.GlobalCfg.Diamond.Master.Stage[i].Num-1] >= equipMasterCfg.Limit {
				gemLevelNew = equipMasterCfg.Limit
				break
			}
		}
	}
	if gemLevelNew > gemLevelOld {
		p.upSuitShushuDate(constant.EquipMasterCategoryGemStone, gemLevelNew, partnerslotid)
	}
	return returnDate, nil
}

// 根据Key字段对Item切片进行排序，并返回最大的4个元素
func topFourByField(items []global.IdKeyNum, lenNum int, compare int) []global.IdKeyNum {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Num == items[j].Num {
			// 如果 Num 相同，则比较 Key 字段
			return items[i].Key > items[j].Key // 降序排序 Key
		}
		if compare > 0 {
			return items[i].Num > items[j].Num
		} else {
			return items[i].Num < items[j].Num
		}
	})

	// 然后，如果元素数量超过4，取前4个元素
	if len(items) > lenNum {
		items = items[:lenNum]
	}

	// 返回排序后的最大的4个元素
	return items
}

// 每日必做副本首次挑战
func (p *Player) dailyMapFirstChallenge(roomType int32) bool {
	if _, ok := p.dailyMaps[roomType]; !ok {
		p.dailyMaps[roomType] = 0
		return true
	}
	return false
}

// 每日必做副本挑战次数
func (p *Player) dailyMapTimes(roomType int32) int32 {
	if _, ok := p.dailyMaps[roomType]; !ok {
		p.dailyMaps[roomType] = 0
	}

	return p.dailyMaps[roomType]
}

func (p *Player) subDailyMapTime(roomType int32) {
	if _, ok := p.dailyMaps[roomType]; !ok {
		p.dailyMaps[roomType] = 0
	}
	p.dailyMaps[roomType]--
}

func (p *Player) dailyMapDoubleReward(b bool) (int32, string) {
	var code int32
	var err string
	if !b {
		p.roomType = 0
		p.reward = []global.IdNum{}
	} else {
		if p.roomType > 0 {
			if cfg := global.GetDailyMapDoubleaccess(p.roomType); cfg != nil {
				if p.checkTakeSource(cfg.ID, cfg.Num) {
					p.takeSource(cfg.ID, cfg.Num)
					if err := p.checkAddSourceBatch(p.reward); err == nil {
						p.addSourceBatch(p.reward, constant.ItemAddDailymap)
					} else {
						attachment := []protocol.Attachment{}
						for _, v := range p.reward {
							attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
						}
						p.SendMail(constant.ItemAddDailymap, attachment)
					}
					p.roomType = 0
					p.reward = []global.IdNum{}
				} else {
					code = -1
					err = errutil.ErrCoinNotEnough.Error()
				}
			} else {
				code = -2
				err = errutil.ErrGetAward.Error()
			}
		}
	}
	return code, err
}

func (p *Player) dailyMapCopySweep(roomID int32) (int32, string) {
	var code int32
	var err string

	roomCfg := global.GetRoomCfg(roomID)
	if roomCfg == nil {
		code = -1
		err = errutil.ErrConfigNotFound.Error()
		return code, err
	}

	// VIP等级
	// if p.vip < int32(global.GetDailyMapTime(int32(roomCfg.Roomtype)).VipLimit) {
	// 	code = -2
	// 	err = errutil.VIPNotEnough
	// 	return code, err
	// }

	// 特权等级
	bPrivilege := false
	for _, v := range p.privilege {
		if privilegeCfg := global.GetPrivilegeCfg(v.Pid); privilegeCfg != nil {
			if privilegeCfg.QuicklyDaily > 0 {
				bPrivilege = true
				break
			}
		}
	}
	if !bPrivilege {
		code = -3
		err = errutil.PrivilegeNotEnough
		return code, err
	}

	// 挑战次数
	var vTimes int32
	if vipCfg := global.GetVipCfg(p.vip); vipCfg != nil {
		for _, v := range vipCfg.DailyTime {
			if v.Key == roomCfg.Roomtype {
				vTimes += v.Value
			}
		}
	}
	if p.dailyMapTimes(roomCfg.Roomtype) <= 0 {
		code = -3
		err = errutil.LimitHasBeenReached
		return code, err
	}

	// 减去挑战次数
	p.subDailyMapTime(roomCfg.Roomtype)

	rewards := []global.IdNum{}
	switch roomCfg.Roomtype {
	case constant.ExperienceFBMap:
		baseValue := global.GlobalCfg.DailyMaps.ExpMap.BaseValue
		num := float64(baseValue.Num) * (1 + float64(global.GlobalCfg.DailyMaps.ExpMap.QuickAccess.Score)/BaseRate) * float64(p.level)
		rewards = append(rewards, global.IdNum{ID: baseValue.ID, Num: int32(num)})

	case constant.GoldFBMap:
		baseValue := global.GlobalCfg.DailyMaps.GoldMap.BaseValue
		num := float64(baseValue.Num) * (1 + float64(global.GlobalCfg.DailyMaps.GoldMap.QuickAccess.Score)/BaseRate) * float64(p.level)
		rewards = append(rewards, global.IdNum{ID: baseValue.ID, Num: int32(num)})

	case constant.EnhanceStoneFBMap:
		baseValue := global.GlobalCfg.DailyMaps.StrengthenMap.BaseValue
		num := float64(baseValue.Num) * float64(global.GlobalCfg.DailyMaps.StrengthenMap.QuickAccess.Score) * float64(p.level)
		rewards = append(rewards, global.IdNum{ID: baseValue.ID, Num: int32(num)})

	case constant.BookPageFBMap:
		baseValue := global.GlobalCfg.DailyMaps.SkillPageMap.BaseValue
		num := float64(baseValue.Num) * float64(global.GlobalCfg.DailyMaps.SkillPageMap.QuickAccess.Score) * float64(p.level)
		rewards = append(rewards, global.IdNum{ID: baseValue.ID, Num: int32(num)})

	case constant.StarUpStoneFBMap:
		for _, v := range global.GlobalCfg.DailyMaps.StarStoneMap.QuickAccess.Stages {
			rewards = append(rewards, monsterMgr.MonsterDorp(v.ID, int32(v.Num))...)
		}

	case constant.SoulStoneFBMap:
		for _, v := range global.GlobalCfg.DailyMaps.SoulStoneMap.QuickAccess.Stages {
			rewards = append(rewards, monsterMgr.MonsterDorp(v.ID, int32(v.Num))...)
		}
	}

	if err := p.checkAddSourceBatch(rewards); err == nil {
		p.addSourceBatch(rewards, constant.ItemAddDailymap)
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range rewards {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		p.SendMail(constant.ItemAddDailymap, attachment)
	}

	packet := &pb.GameResultBattlePush{Items: []*pb.ItemInfo{}, RoomId: roomCfg.ID, OverType: pb.GameOverType_GameOverTypeCopySweep}
	for _, v := range rewards {
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	}
	p.pushMsg(protocol.GameResultBattleMapPush, packet)

	p.tasks.dispatch(pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE, 0, 1)

	p.reward = rewards
	p.roomType = roomCfg.Roomtype
	return code, err
}

// 采集
func (p *Player) checkMining() {
	if p.desk == nil {
		return
	}
	if p.miningObj == nil {
		return
	}
	if p.miningObj.dead {
		return
	}
	if p.miningObj.cfg.MType != MTypeGathering {
		return
	}

	if time.Now().Before(p.miningObj.miningStartTime.Add(time.Duration(p.miningObj.cfg.Delay) * time.Millisecond)) {
		return
	}

	p.miningEnd()
}

func (p *Player) miningStart(uid int64) {
	if p.desk == nil {
		return
	}
	p.miningBreak()

	//p.logger.Debugf("采集开始 Name: %s  目标: %d", p.name, uid)

	m, err := p.desk.monsterWithUid(uid)
	if err != nil {
		p.logger.Debugf(" %s 采集失败,怪不存在 err %s", p.name, err.Error())
		return
	}
	monCfg := global.GetMonsterCfg(m.cfgId)
	if monCfg == nil {
		return
	}

	if (m.ownershipType == pb.OwnershipType_OWNERSHIP_PLAYER && p.uid != m.ownershipUid) ||
		m.ownershipType == pb.OwnershipType_OWNERSHIP_GUILD && p.guildId != m.ownershipUid {
		p.logger.Debugf("采集失败,目标有归属 归属类型 %v 归属者ID: %d ", m.ownershipType, m.ownershipUid)
		return
	}

	if m.cfg.MType == 3 && m.cfg.Trigger == constant.MonsterTrigerTypeEquipFuDai {
		if p.bagEquipLeftCount() <= 0 {
			p.logger.Debugf("采集失败,包裹已满 归属类型 %v 归属者ID: %d ", m.ownershipType, m.ownershipUid)
			return
		}
	}

	if m.miningPlayer != nil {
		p.logger.Debugln("采集失败,正在被其他玩家采集 ")
		return
	}

	m.miningStartTime = time.Now()
	m.miningPlayer = p
	p.SetObAction(ObjectActionMining)
	p.miningObj = m

	p.desk.Broadcast(protocol.RouteMiningPush, &pb.MiningPush{PlayerUid: p.uid, TargetUid: m.uid, Action: pb.MiningPush_START})

	if !p.getKillMonControl(monCfg.ID) {
		p.miningLimit()
		packet := &pb.SystemNotificationPush{Type: pb.Notification_INVALID_Notification, Msg: errutil.KillMonControl}
		p.pushMsg(protocol.RouteSystemNotificationPush, packet)
		return
	}

	p.desk.DDailyMap.mingCount++
	p.setAttackTime((monCfg.Delay - 1) * 1000)
}

// 采集结束
func (p *Player) miningEnd() {

	p.logger.Debugf("采集结束 Name: %s  目标 : %d", p.name, p.miningObj.uid)

	p.desk.Broadcast(protocol.RouteMiningPush, &pb.MiningPush{PlayerUid: p.uid, TargetUid: p.miningObj.uid, Action: pb.MiningPush_END})

	p.miningObj.die(p)

	p.miningObj.miningPlayer = nil
	p.miningObj = nil

	switch p.miningObj.cfg.ID {
	case constant.ItemIDLuckyBagLow.Int32(): // 初级福袋（蓝色）
		p.tasks.dispatch(pb.TargetType_BATTLEFIELD_OPEN_BLUE_LUCKY_BAG_TIMES, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OPENED_TOTAL, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OBTAINED_TOTAL, 0, 1)
	case constant.ItemIDLuckyBagMid.Int32(): // 中级福袋（紫色）
		p.tasks.dispatch(pb.TargetType_BATTLEFIELD_OPEN_PURPLE_LUCKY_BAG_TIMES, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OPENED_TOTAL, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OBTAINED_TOTAL, 0, 1)
	case constant.ItemIDLuckyBagHigh.Int32(): // 高级福袋（橙色）
		p.tasks.dispatch(pb.TargetType_BATTLEFIELD_OPEN_ORANGE_LUCKY_BAG_TIMES, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OPENED_TOTAL, 0, 1)
		p.tasks.dispatch(pb.TargetType_LUCKY_BAG_OBTAINED_TOTAL, 0, 1)
	default:
		p.tasks.dispatch(pb.TargetType_BATTLEFIELD_TOTAL_CHESTS_OPENED, 0, 1)
	}

	p.SetObAction(ObjectActionStand)
}

// 采集中断
func (p *Player) miningBreak() {

	if p.GetObAction() != ObjectActionMining || p.miningObj == nil || p.miningObj.miningPlayer.Uid() != p.uid {
		return
	}
	p.logger.Debugf("采集中断 Name: %s 目标 : %d", p.name, p.miningObj.uid)

	p.desk.Broadcast(protocol.RouteMiningPush, &pb.MiningPush{PlayerUid: p.uid, TargetUid: p.miningObj.uid, Action: pb.MiningPush_BREAK})

	p.miningObj.miningPlayer = nil
	p.miningObj = nil

	p.SetObAction(ObjectActionStand)
}

// 采集上限
func (p *Player) miningLimit() {

	if p.GetObAction() != ObjectActionMining || p.miningObj == nil || p.miningObj.miningPlayer.Uid() != p.uid {
		return
	}
	p.logger.Debugf("采集上限 Name: %s 目标 : %d", p.name, p.miningObj.uid)

	p.desk.Broadcast(protocol.RouteMiningPush, &pb.MiningPush{PlayerUid: p.uid, TargetUid: p.miningObj.uid, Action: pb.MiningPush_LIMIT})

	p.miningObj.miningPlayer = nil
	p.miningObj = nil

	p.SetObAction(ObjectActionStand)
}

// 等级奖励
func (p *Player) levelGift() {
	if cfg, ok := global.LevelGiftCfg[p.level]; ok {
		source := protocol.AddSource{Source: []*protocol.SourceInfo{}, Method: constant.ItemAddLevelGift}
		for _, v := range cfg.Gift {
			if (v.Gender != 0 && p.sex != int(v.Gender)) || (v.Job != 0 && p.job != constant.MowJob(v.Job)) {
				continue
			}
			source.Source = append(source.Source, &protocol.SourceInfo{Id: v.ID, Num: int64(v.Num)})
		}
		p.addSourceEvent(&source)
	}
}

func (p *Player) ItemDisassembly(itemInfo *pb.ItemInfo) {
	var getitem []global.KeyValue

	if itemInfo == nil {
		return
	}
	itemCfg, err := global.GetItemCfg(itemInfo.CfgId)
	if err != nil {
		return
	}

	for _, des := range itemCfg.Decompose {
		bFind := false
		for i, gi := range getitem {
			if gi.Key == des.Key {
				getitem[i].Value += des.Value
				bFind = true
				break
			}
		}
		if !bFind {
			getitem = append(getitem, des)
		}
	}
	source := protocol.AddSource{Source: []*protocol.SourceInfo{}, Method: constant.ItemAddDisintegrate}
	for _, v := range getitem {
		source.Source = append(source.Source, &protocol.SourceInfo{Id: v.Key, Num: int64(v.Value)})
	}
	p.addSourceEvent(&source)
}

// 设置宝箱
func (p *Player) homeBoxSetting(nowOpenLevel int32, Job int32) error {
	// 更新数据库
	// p.box.NowOpenLevel = nowOpenLevel
	// p.box.Job = Job
	p.SyncUserBox()
	return nil
}

func (p *Player) openbox() ([]*pb.ItemInfo, error) {

	randomEquip := p.box.Preview
	if len(randomEquip) == 0 {
		return nil, nil
	}
	checkItem := make([]global.IdNum, 0)
	for _, v := range randomEquip {
		checkItem = append(checkItem, global.IdNum{ID: v.ID, Num: v.Num})
	}
	if err := p.checkAddSourceBatch(checkItem); err == nil {
		p.box.NowOpenRes = make([]pb.ItemInfo, len(randomEquip))
		if newItem := p.addEquipBatchWithQuality(randomEquip, constant.ItemAddOpenBox, WithAddSourceUnNotify(AddSourceTypeAutoDis)); newItem != nil {
			for i, v := range newItem {
				p.box.NowOpenRes[i].CfgId = v.CfgId
				p.box.NowOpenRes[i].Uid = v.Uid
				p.box.NowOpenRes[i].Num = v.Num
				p.box.NowOpenRes[i].Equip = v.Equip
			}
			p.box.Preview = nil
			return newItem, nil
		}
		// 触发任务累计
		// p.tasks.dispatch(pb.TargetType_TOTAL_BOX_OPENED, 0, int32(num))
	}
	return nil, errutil.ErrBagFull
}

// 预开宝箱
func (p *Player) mainBattlePrepare(job int32, mainbattleId int32) ([]global.IdNumQuality, error) {
	// 查询num是不是合法的
	// BoxInfo := global.GetBoxLevelCfg(p.box.NowLevel)
	targetBattleID := mainbattleId
	if targetBattleID == 0 {
		targetBattleID = p.mainBattleId + 1
	}
	if targetBattleID > p.mainBattleId+1 {
		return nil, errutil.ErrConfigNotFound
	}
	maxBattleID := int32(len(global.MainBattleCfg))
	if targetBattleID > maxBattleID {
		targetBattleID = maxBattleID
	}
	mainBattleInfo := global.GetMainBattleCfg(targetBattleID)
	if mainBattleInfo == nil {
		return nil, errutil.ErrConfigNotFound
	}
	randNum := len(mainBattleInfo.Monster)
	num := int32(0)
	for i := 1; i <= randNum; i++ {
		for _, v := range global.GlobalCfg.MainBattle.Keys {
			if v.Flag == int32(mainBattleInfo.Flag) {
				minDrop := v.MinDrop
				maxDrop := v.MaxDrop
				num += int32(global.RandIn(int(maxDrop)-int(minDrop)) + int(minDrop))
				break
			}
		}
	}

	if num == 0 {
		return nil, errutil.ErrOpenBoxEmpty
	}
	randomEquip := p.randomEquip(p.level, num, job)
	if len(randomEquip) == 0 {
		return nil, errutil.ErrOpenBoxEmpty
	}
	preview := make([]global.IdNumQuality, 0)
	for _, v := range randomEquip {
		BoxLevelCfg := global.GetBoxLevelCfg(p.box.NowLevel)
		if BoxLevelCfg == nil {
			return nil, errutil.ErrOpenBoxEmpty
		}
		var rate int32
		for _, v2 := range BoxLevelCfg.Qualityrandom {
			rate += v2.Rate
		}
		weight := global.RandIn(int(rate))
		DemoRandom := global.IdRate{}
		beginPriority := 0
		for _, v3 := range BoxLevelCfg.Qualityrandom {
			if beginPriority+int(v3.Rate) >= weight {
				DemoRandom = v3
				break
			}
			beginPriority += int(v3.Rate)
		}
		preview = append(preview, global.IdNumQuality{
			ID:      v.ID,
			Num:     v.Num,
			Quality: DemoRandom.Id,
		})
	}
	p.box.Preview = preview
	return preview, nil
}

func (p *Player) openboxV2(level int32, num int32, job int32) ([]*pb.ItemInfo, error) {
	// level := p.box.NowOpenLevel
	// 如果没设置，就等于角色等级
	// if level == 0 {
	// 	level = p.level
	// }
	// job := p.box.Job
	// needLevel := int32(1)
	// for _, v := range global.GlobalCfg.Openslot.IDLevel {
	// 	if v.ID == 1 {
	// 		needLevel = v.Needlevel
	// 		break
	// 	}
	// }
	// if p.level < needLevel {
	// 	job = int32(p.job)
	// }
	// 查询num是不是合法的
	// BoxInfo := global.GetBoxLevelCfg(p.box.NowLevel)
	// if BoxInfo.Maxtime < num {
	// 	return nil, errutil.ErrEnergyFull
	// }
	level = p.level
	num = global.GlobalCfg.BaseValue.Offlinebox.Equipnum
	randomEquip := p.randomEquip(level, num, job)
	if len(randomEquip) == 0 {
		return nil, errutil.ErrOpenBoxEmpty
	}
	if err := p.checkAddSourceBatch(randomEquip); err == nil {
		itemInfo := p.addSourceBatch(randomEquip, constant.ItemAddOpenBox)
		p.box.NowOpenRes = make([]pb.ItemInfo, len(itemInfo))
		for i, item := range itemInfo {
			p.box.NowOpenRes[i].CfgId = item.CfgId
			p.box.NowOpenRes[i].Uid = item.Uid
			p.box.NowOpenRes[i].Num = item.Num
			p.box.NowOpenRes[i].Equip = item.Equip
		}
		// 触发任务累计
		//p.tasks.dispatch(pb.TargetType_TOTAL_BOX_OPENED, 0, int32(num))

		return itemInfo, nil
	}
	return nil, errutil.ErrBagFull
}

// 升级宝箱
func (p *Player) upgradeBox() error {
	// 判断是否可以升级
	if int(p.box.NowLevel) >= len(global.BoxLevelCfg) {
		return errutil.ErrConfigNotFound
	}
	// 判断等级够不够
	if p.level < global.BoxLevelCfg[p.box.NowLevel].Rolelevel {
		return errutil.ErrLevel
	}
	//主线是不是达标
	if p.mainBattleId < global.BoxLevelCfg[p.box.NowLevel].BattleID {
		return errutil.ErrMainBattle
	}
	// 判断有没有材料
	BoxInfo := global.GetBoxLevelCfg(p.box.NowLevel)
	if BoxInfo == nil {
		return errutil.ErrConfigNotFound
	}
	p.box.CostGold = 0
	for _, v := range BoxInfo.Consume {
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return errutil.ErrItemNotEnough
		}
	}
	for _, v := range BoxInfo.Consume {
		p.takeSource(v.ID, int64(v.Num))
	}
	p.box.NowLevel++
	p.tasks.dispatch(pb.TargetType_CHEST_LEVEL_UPGRADED, 0, p.box.NowLevel)
	p.upBoxShushuDate(p.box)
	p.SyncUserBox()
	return nil
}

// 加快升级宝箱
func (p *Player) UpgradeBoxQuick(num int32) error {
	// 判断是否在升级中
	if p.box.NextBoxRecoveryTime == 0 {
		return errutil.ErrBoxRecovery
	}
	// 判断有没有材料
	BoxInfo := global.GetBoxLevelCfg(p.box.NowLevel)
	if BoxInfo == nil {
		return errutil.ErrConfigNotFound
	}
	// 扣材料
	canNum := p.box.NextBoxRecoveryTime / int64(global.GlobalCfg.BaseValue.Boxlevelup.Accelerate)
	if num > int32(canNum) {
		num = int32(canNum)
		if p.box.NextBoxRecoveryTime%int64(global.GlobalCfg.BaseValue.Boxlevelup.Accelerate) != 0 {
			num++
		}
	}

	// 扣材料
	if !p.checkTakeSource(global.GlobalCfg.BaseValue.Boxlevelup.Itemid, int64(num)) {
		return errutil.ErrItemNotEnough
	}
	p.takeSource(global.GlobalCfg.BaseValue.Boxlevelup.Itemid, int64(num))

	p.box.NextBoxRecoveryTime -= int64(global.GlobalCfg.BaseValue.Boxlevelup.Accelerate * num)
	if p.box.NextBoxRecoveryTime < 0 {
		p.box.NextBoxRecoveryTime = 0
	} else {
		p.SyncUserBox()
	}
	p.box.IsAccelerate = 1
	p.box.AccelerateNum += num
	return nil
}

// 给到随机的奖励
func (p *Player) randomEquip(level int32, num int32, job int32) (rewards []global.IdNum) {

	createCfg := global.OpenBoxCfg
	if len(createCfg) == 0 {
		return []global.IdNum{}
	}
	for _, v := range createCfg {
		// todo校验是否能开这么多箱子
		for _, v1 := range v.Level {
			if level >= v1.MinLevel && level <= v1.MaxLevel {
				var dropInfo []global.IdNum
				//要求是注释掉
				// ranType := int32(0)
				// if len(p.Guidance) > 0 {
				// 	ranType = p.Guidance[0]
				// }
				// switch ranType {
				// case constant.GuidanceTypeBoxOne:
				// 	switch job {
				// 	case 1:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 209, Num: 1},
				// 			global.IdNum{ID: 211, Num: 1},
				// 		)
				// 	case 2:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 2021, Num: 1},
				// 			global.IdNum{ID: 2023, Num: 1},
				// 		)
				// 	case 3:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 3833, Num: 1},
				// 			global.IdNum{ID: 3835, Num: 1},
				// 		)
				// 	default:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 209, Num: 1},
				// 			global.IdNum{ID: 211, Num: 1},
				// 		)
				// 	}
				// 	p.Guidance[0] = constant.GuidanceTypeBoxOne + 1
				// case constant.GuidanceTypeBoxTwo:
				// 	switch job {
				// 	case 1:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 210, Num: 1},
				// 			global.IdNum{ID: 211, Num: 1},
				// 		)
				// 	case 2:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 2022, Num: 1},
				// 			global.IdNum{ID: 2023, Num: 1},
				// 		)
				// 	case 3:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 3834, Num: 1},
				// 			global.IdNum{ID: 3835, Num: 1},
				// 		)

				// 	default:
				// 		dropInfo = append(dropInfo,
				// 			global.IdNum{ID: 210, Num: 1},
				// 			global.IdNum{ID: 211, Num: 1},
				// 		)
				// 	}
				// 	p.Guidance[0] = constant.GuidanceTypeBoxTwo + 1
				// default:
				for i := 0; i < int(num); i++ {
					var rate int32
					for _, v2 := range v.Equiprandom {
						if job == 0 {
							rate += v2.Rate
						} else {
							if v2.Job == job {
								rate += v2.Rate
							}
						}
					}
					weight := global.RandIn(int(rate))
					DemoRandom := global.IdRateJob{}
					beginPriority := 0
					for _, v3 := range v.Equiprandom {
						if job == 0 {
							if beginPriority+int(v3.Rate) >= weight {
								DemoRandom = v3
								break
							}
							beginPriority += int(v3.Rate)
						} else {
							if v3.Job == job {
								if beginPriority+int(v3.Rate) >= weight {
									DemoRandom = v3
									break
								}
								beginPriority += int(v3.Rate)
							}
						}
					}
					dropInfo = append(dropInfo, global.IdNum{ID: DemoRandom.Id, Num: 1})
				}
				// }
				return dropInfo
			} else {
				continue
			}
		}
	}
	return []global.IdNum{}
}

// 随机装备品质
func (p *Player) randomEquipQuality(item *pb.ItemInfo) {
	if item.Type != pb.ItemType_EQUIPMENT {
		return
	}
	BoxLevelCfg := global.GetBoxLevelCfg(p.box.NowLevel)
	if BoxLevelCfg == nil {
		return
	}
	var rate int32
	for _, v2 := range BoxLevelCfg.Qualityrandom {
		rate += v2.Rate
	}
	weight := global.RandIn(int(rate))
	DemoRandom := global.IdRate{}
	beginPriority := 0
	for _, v3 := range BoxLevelCfg.Qualityrandom {
		if beginPriority+int(v3.Rate) >= weight {
			DemoRandom = v3
			break
		}
		beginPriority += int(v3.Rate)
	}
	if item.Equip == nil {
		item.Equip = &pb.EquipInfo{Quality: pb.ItemQuality_INVALID_Quality}
	}
	item.Equip.Quality = pb.ItemQuality(DemoRandom.Id)
}

func (p *Player) boxResIntoBag(Uid int64) error {
	for i := range p.box.NowOpenRes {
		if p.box.NowOpenRes[i].Uid == Uid {
			// 删除这个值
			p.box.NowOpenRes = append(p.box.NowOpenRes[:i], p.box.NowOpenRes[i+1:]...)
			break
		}
	}
	return nil
}

func (p *Player) getCombatPartnerNum() int32 {
	if p.partners == nil {
		return 0
	}
	num := 1
	for _, v := range p.partnerSlots {
		if v.PartnerUid > 0 {
			num++
		}
	}
	return int32(num)
}

// 伙伴信息
func (p *Player) getPartnerInfo(uid int64) *Partner {
	partner, ok := p.partners[uid]
	if !ok {
		return nil
	}
	return partner
}

func (p *Player) getPartnerInfoByCfgId(cfgId int32) *Partner {
	for _, v := range p.partners {
		if v.cfg.Cid == cfgId {
			return v
		}
	}
	return nil
}

func (p *Player) getPartnerInfoBySlotId(slot int32) *Partner {
	if slot < 1 || slot > 4 {
		return nil
	}
	if p.partnerSlots[slot-1].PartnerUid == 0 {
		return nil
	}
	partner, ok := p.partners[p.partnerSlots[slot-1].PartnerUid]
	if !ok {
		return nil
	}
	return partner
}

// 伙伴槽位信息
func (p *Player) getPartnerSlotInfo(slot int32) *PartnerSlot {
	for i := 0; i < len(p.partnerSlots); i++ {
		if p.partnerSlots[i].SlotId == slot {
			return p.partnerSlots[i]
		}
	}
	return nil
}

// 增加伙伴
func (p *Player) addPartner(cfgId int32) error {
	if cfgId == 0 {
		return nil
	}

	// 重复获得伙伴转换成碎片
	for _, v := range p.partners {
		if v.cfg.Cid == cfgId {
			source := protocol.AddSource{Source: []*protocol.SourceInfo{}, Method: constant.ItemAddPartner}
			for _, c := range v.cfg.ChipNum {
				source.Source = append(source.Source, &protocol.SourceInfo{Id: c.Key, Num: int64(c.Value)})
			}
			p.addSourceEvent(&source)
			return nil
		}
	}

	par := CreatePartner(cfgId)
	par.owner = p
	p.partners[par.Uid()] = par

	par.setLevel(p.level)

	// 解锁技能
	par.learnSkill()

	p.syncSavePartnerDb(par.uid, model.Add)
	p.rankChangeTCP(constant.RankTypeRedPartnerStar, p.calcRedPartnerStarNum())
	par.SyncUserPartner()
	p.checkPartnerBonds()
	p.refreshStats()

	return nil
}

// 伙伴上阵
func (p *Player) partnerCombat(slotId int32, partnerUid int64) error {
	partnerSlot := p.getPartnerSlotInfo(slotId)
	if partnerSlot == nil {
		return errutil.ErrPartnerInvalid
	}

	partner := p.getPartnerInfo(partnerUid)
	if partner == nil {
		return errutil.ErrPartnerInvalid
	}

	if partnerSlot.PartnerUid > 0 {
		return errutil.ErrPartnerInvalid
	}

	for _, v := range p.partnerSlots {
		if v.PartnerUid == partnerUid {
			return errutil.ErrPartnerInvalid
		}

		info := p.getPartnerInfo(v.PartnerUid)
		if info != nil {
			if info.cfg.Cid == partner.cfg.Cid {
				return errutil.ErrPartnerInvalid
			}
		}
	}

	// 自动装备伙伴技能
	bskillslotchange := false
	for _, v := range partner.skills {
		if skillCfg, err := global.GetSkillCfg(v.CfgId); err == nil {
			if skillCfg.ActiveorPassive == PassiveSkill {
				continue
			}
			if skillCfg.SType == constant.SkillTypeNormal {
				continue
			}
			bFind := false
			for slot := 1; slot < len(p.skillsSlot); slot++ {
				if p.skillsSlot[slot] == v.CfgId {
					bFind = true
					break
				}
			}
			if !bFind {
				for slot := 1; slot < len(p.skillsSlot); slot++ {
					if p.skillsSlot[slot] == 0 {
						p.skillsSlot[slot] = v.CfgId
						bskillslotchange = true
						break
					}
				}
			}
		}
	}
	if bskillslotchange {
		p.syncUserSkillSlot()
	}

	partnerSlot.PartnerUid = partnerUid
	num := int32(0)
	for _, v := range p.partnerSlots {
		if v.PartnerUid > 0 {
			num++
		}
	}
	if num > 0 {
		p.tasks.dispatch(pb.TargetType_PARTNER_ON_STAGE_TIMES, 0, num)
	}
	p.refreshStats()
	p.SyncUserPartnerSlot()
	return nil
}

// 伙伴休息
func (p *Player) partnerRest(slotId int32) error {
	slotInfo := p.getPartnerSlotInfo(slotId)
	if slotInfo == nil {
		return errutil.ErrPartnerInvalid
	}
	partner := p.getPartnerInfoBySlotId(slotId)
	if partner == nil {
		return errutil.ErrPartnerInvalid
	}

	equipnum := 0
	for _, v := range slotInfo.Equips {
		if v.ItemId > 0 {
			equipnum++
		}
	}

	if equipnum > p.bagEquipLeftCount() {
		return errutil.ErrBagFull
	}

	partnerSlot := p.getPartnerSlotInfo(slotId)
	if partnerSlot == nil {
		return errutil.ErrPartnerInvalid
	}

	// p.partnerSlots[slotId].PartnerUid = 0
	for k, v := range partnerSlot.Equips {
		if v.ItemId <= 0 {
			continue
		}
		p.changeEquip(2, pb.EquipSlot(k), v.ItemId, partnerSlot.SlotId)
	}

	// 自动下阵伙伴技能
	bskillslotchange := false
	for _, v := range partner.skills {
		for slot := 0; slot < len(p.skillsSlot); slot++ {
			if p.skillsSlot[slot] == v.CfgId {
				p.skillsSlot[slot] = 0
				bskillslotchange = true
			}
		}
	}
	if bskillslotchange {
		p.syncUserSkillSlot()
	}

	partnerSlot.PartnerUid = 0

	p.refreshStats()
	p.SyncUserPartnerSlot()
	return nil
}

// 伙伴合成
func (p *Player) partnerCompose(cfgId int32) error {
	for _, v := range p.partners {
		if v.cfg.Cid == cfgId {
			return errutil.ErrPartnerInvalid
		}
	}

	if cfg := global.GetPartnerCfg(cfgId); cfg != nil {
		for _, v := range cfg.ChipNum {
			if !p.checkTakeSource(v.Key, int64(v.Value)) {
				return errutil.ErrItemNotEnough
			}
		}

		for _, v := range cfg.ChipNum {
			p.takeSource(v.Key, int64(v.Value))
		}

		p.addPartner(cfgId)
	}
	return nil
}

// partnerReforge
func (p *Player) partnerReforge(partnerUid int64, isFree bool) error {
	par := p.getPartnerInfo(partnerUid)
	if par == nil {
		return errutil.ErrPartnerInvalid
	}
	if par.level <= 1 {
		return errutil.ErrLevel
	}

	var gotItems []global.IdNum
	var gotExp int64

	if !isFree {
		for _, v := range global.GlobalCfg.Partner.Forge.Charge {
			if v.Quality == par.cfg.Quality {
				if p.getItemCountInBag(int32(v.ID)) < int64(v.Num) {
					return errutil.ErrItemNotEnough
				}
				p.takeSource(v.ID, v.Num)
			}
		}
	}

	for i := par.star; i > 0; i-- {
		for _, v := range par.cfg.StarConsume {
			if v.Id == i {
				for _, cost := range v.Values {
					bFind := false
					for j := 0; j < len(gotItems); j++ {
						if gotItems[j].ID == cost.Key {
							gotItems[j].Num += cost.Value
							bFind = true
						}
					}
					if !bFind {
						gotItems = append(gotItems, global.IdNum{ID: cost.Key, Num: cost.Value})
					}
				}
			}
		}
	}

	for i := par.level - 1; i > 0; i-- {
		cfg, err := global.GetRoleCfg(int32(par.cfg.Job), i)
		if err != nil {
			continue
		}
		gotExp += cfg.NeedExp
	}

	if isFree {
		for i := 0; i < len(gotItems); i++ {
			gotItems[i].Num = int32(math.Ceil(float64(gotItems[i].Num) * float64(global.GlobalCfg.Partner.Forge.Free.Rate/100.0)))
		}
		gotExp = int64(math.Ceil(float64(gotExp) * float64(global.GlobalCfg.Partner.Forge.Free.Rate/100.0)))
	}

	gotItems = append(gotItems, global.Exp2ExpItems(gotExp)...)

	p.addSourceBatch(gotItems, constant.ItemAddDisintegrate)

	// 等级1  星级0 装备全部卸下
	par.star = 0
	par.level = 1
	for j := 0; j < len(p.skillsSlot); j++ {
		for _, v := range par.skills {
			if v.CfgId == p.skillsSlot[j] {
				p.skillsSlot[j] = 0
			}
		}
	}

	for _, slot := range p.partnerSlots {
		if slot.PartnerUid == partnerUid {
			for k, v := range slot.Equips {
				if v.ItemId <= 0 {
					continue
				}
				p.changeEquip(2, pb.EquipSlot(k), v.ItemId, slot.SlotId)
			}
			break
		}
	}
	// 解锁技能
	par.skills = []*PlayerSkill{}
	par.learnSkill()

	par.SyncUserPartner()
	//par.refreshStats(nil, p)
	p.SyncUserPartnerSlot()
	p.syncUserSkillSlot()
	p.checkPartnerBonds()
	p.refreshStats()

	return nil
}

// partnerDisassemble
func (p *Player) partnerDisassemble(partnerUid int64) error {
	par := p.getPartnerInfo(partnerUid)
	if par == nil {
		return errutil.ErrPartnerInvalid
	}

	if par.cfg.Quality < int32(global.GlobalCfg.Partner.Decompose.Quality) {
		return errutil.ErrPartnerInvalid
	}

	if par.star < int32(global.GlobalCfg.Partner.Decompose.Starlevel) {
		return errutil.ErrPartnerInvalid
	}

	var gotItems []global.IdNum
	var gotExp int64
	var tmp []global.IdNum
	for i := par.star; i > 0; i-- {
		for _, v := range par.cfg.StarConsume {
			if v.Id == i {
				for _, cost := range v.Values {
					bFind := false
					for j := 0; j < len(tmp); j++ {
						if tmp[j].ID == cost.Key {
							tmp[j].Num += cost.Value
							bFind = true
							break
						}
					}
					if !bFind {
						tmp = append(tmp, global.IdNum{ID: cost.Key, Num: cost.Value})
					}
				}
			}
		}
	}

	for i := 0; i < len(tmp); i++ {
		if cfg, err := global.GetItemCfg(tmp[i].ID); err == nil {
			if cfg.MainType == constant.PartnerFragment {
				num := int32(math.Ceil(float64(tmp[i].Num) * float64(global.GlobalCfg.Partner.Decompose.Rate/100.0)))
				if num > 0 {
					bFind := false
					for j := 0; j < len(gotItems); j++ {
						if gotItems[j].ID == int32(global.GlobalCfg.Partner.Decompose.ID) {
							gotItems[j].Num += num
							bFind = true
							break
						}
					}
					if !bFind {
						gotItems = append(gotItems, global.IdNum{ID: int32(global.GlobalCfg.Partner.Decompose.ID), Num: num})
					}
				}

			} else {
				gotItems = append(gotItems, tmp[i])
			}
		}
	}

	for i := par.level - 1; i > 0; i-- {
		cfg, err := global.GetRoleCfg(int32(par.cfg.Job), i)
		if err != nil {
			continue
		}
		gotExp += cfg.NeedExp
	}

	gotItems = append(gotItems, global.Exp2ExpItems(gotExp)...)

	p.addSourceBatch(gotItems, constant.ItemAddDisintegrate)

	// 等级1  星级0 装备技能全部卸下
	par.star = 0
	par.level = 1
	for j := 0; j < len(p.skillsSlot); j++ {
		for _, v := range par.skills {
			if v.CfgId == p.skillsSlot[j] {
				p.skillsSlot[j] = 0
			}
		}
	}

	for _, slot := range p.partnerSlots {
		if slot.PartnerUid == partnerUid {
			for k, v := range slot.Equips {
				if v.ItemId <= 0 {
					continue
				}
				p.changeEquip(2, pb.EquipSlot(k), v.ItemId, slot.SlotId)
			}
			break
		}
	}
	// 解锁技能
	par.skills = []*PlayerSkill{}
	par.learnSkill()
	par.SyncUserPartner()
	par.refreshStats(nil, p)
	p.SyncUserPartnerSlot()
	p.syncUserSkillSlot()

	return nil
}

// 是否已装备
func (p *Player) isOnEquip(itemUid int64) bool {
	for _, v := range p.equips {
		if v.ItemId == itemUid {
			return true
		}
	}

	for _, v := range p.partnerSlots {
		for _, vv := range v.Equips {
			if vv.ItemId == itemUid {
				return true
			}
		}
	}
	return false
}

// 一键戴最好装备
func (p *Player) equipBest(objType pb.ObjectType, partnerSlotID int32) {
	if partnerSlotID < 0 || partnerSlotID > 5 {
		return
	}
	levevl_ := p.level
	job_ := p.job
	for slot := 1; slot < int(pb.EquipSlot_MaxEquipSlot); slot++ {

		var equipSlot *Slot
		if objType == pb.ObjectType_PLAYER {
			equipSlot = &p.equips[slot]
		} else {
			if p.getPartnerSlotInfo(partnerSlotID).PartnerUid == 0 {
				return
			}
			partner := p.getPartnerInfoBySlotId(partnerSlotID)
			if partner == nil {
				return
			}
			levevl_ = partner.level
			job_ = constant.MowJob(partner.cfg.Job) % 10
			equipSlot = &p.getPartnerSlotInfo(partnerSlotID).Equips[slot]
		}

		var activeEquip *pb.ItemInfo
		if p.equips[slot].ItemId > 0 {
			activeEquip = p.getItemByUid(equipSlot.ItemId)
		}

		for _, item := range p.bag.ItemList {
			if p.hasHomeOpenBox(item.Uid) {
				continue
			}
			if item.CfgId == 0 {
				continue
			}
			if item.Type != pb.ItemType_EQUIPMENT {
				continue
			}

			itemCfg, err := global.GetItemCfg(item.CfgId)
			if err != nil {
				continue
			}

			if itemCfg.Equiptype != slot {
				continue
			}

			if itemCfg.Levellimit > int(levevl_) {
				continue
			}

			if itemCfg.Occupylimit != 0 && itemCfg.Occupylimit != int(job_) {
				continue
			}

			if p.isOnEquip(item.Uid) {
				continue
			}

			if activeEquip == nil {
				p.equipItem(pb.EquipSlot(slot), item.Uid, int32(partnerSlotID))
				continue
			}

			activeEquip = p.getItemByUid(equipSlot.ItemId)

			if objType == pb.ObjectType_PLAYER {
				p.refreshStats(WithRefreshStatsPushMessage(false))
				currentPower := p.calcCombatRealV2(&Stats{Attrs: p.stats.copyStats()}, true)

				if err := p.unequipItem(pb.EquipSlot(slot), int32(partnerSlotID)); err != nil {
					logger.Errorf("Failed to unequip item: %v", err)
					continue
				}

				if err := p.equipItem(pb.EquipSlot(slot), item.Uid, int32(partnerSlotID)); err != nil {
					logger.Errorf("Failed to equip item: %v", err)
					if activeEquip != nil {
						if err = p.equipItem(pb.EquipSlot(slot), activeEquip.Uid, int32(partnerSlotID)); err != nil {
							logger.Errorf("Failed to rollback equip operation: %v", err)
						}
					}
					continue
				}

				p.refreshStats(WithRefreshStatsPushMessage(false))
				newPower := p.calcCombatRealV2(&Stats{Attrs: p.stats.copyStats()}, true)

				if newPower <= currentPower {
					if err := p.unequipItem(pb.EquipSlot(slot), int32(partnerSlotID)); err != nil {
						logger.Errorf("Failed to unequip item during rollback: %v", err)
						continue
					}
					if activeEquip != nil {
						if err := p.equipItem(pb.EquipSlot(slot), activeEquip.Uid, int32(partnerSlotID)); err != nil {
							logger.Errorf("Failed to rollback equip operation: %v", err)
						}
					}
				}
				logger.Debugf("PLAYER equipBest equipItem slot=%d, item=%d, itemCfgID=%d, newPower=%d, currentPower=%d, activeEquip=%d",
					slot, item.Uid, item.CfgId, newPower, currentPower, func() int64 {
						if activeEquip != nil {
							return activeEquip.Uid
						}
						return 0
					}())
			} else {
				partnerUID := p.partnerSlots[partnerSlotID-1].PartnerUid
				partner := p.getPartnerInfo(partnerUID)

				if partner == nil {
					logger.Errorf("Partner not found for UID: %d", partnerUID)
					continue
				}

				currentPacket := &pb.SyncUserStatsPush{Attributes: []*pb.Stats{}, Partners: []*pb.SyncUserStatsPush_PartnerStats{}}
				partner.refreshStats(currentPacket, p)

				if err := p.unequipItem(pb.EquipSlot(slot), int32(partnerSlotID)); err != nil {
					logger.Errorf("Failed to unequip partner item: %v", err)
					continue
				}

				if err := p.equipItem(pb.EquipSlot(slot), item.Uid, int32(partnerSlotID)); err != nil {
					logger.Errorf("Failed to equip partner item: %v", err)
					if activeEquip != nil {
						if err := p.equipItem(pb.EquipSlot(slot), activeEquip.Uid, int32(partnerSlotID)); err != nil {
							logger.Errorf("Failed to rollback partner equip operation: %v", err)
						}
					}
					continue
				}

				newPacket := &pb.SyncUserStatsPush{Attributes: []*pb.Stats{}, Partners: []*pb.SyncUserStatsPush_PartnerStats{}}
				partner.refreshStats(newPacket, p)
				if len(newPacket.Partners) > 0 && len(currentPacket.Partners) > 0 {
					if newPacket.Partners[0].GetCombat() <= currentPacket.Partners[0].GetCombat() {
						if err := p.unequipItem(pb.EquipSlot(slot), int32(partnerSlotID)); err != nil {
							logger.Errorf("Failed to unequip partner item during rollback: %v", err)
							continue
						}
						if activeEquip != nil {
							if err := p.equipItem(pb.EquipSlot(slot), activeEquip.Uid, int32(partnerSlotID)); err != nil {
								logger.Errorf("Failed to rollback partner equip operation: %v", err)
							}
						}
					}
				}

				logger.Debugf("PARTNER equipBest equipItem slot=%d, item=%d, itemCfgID=%d, newPower=%d, currentPower=%d, activeEquip=%d",
					slot, item.Uid, item.CfgId, newPacket.Combat, currentPacket.Combat, func() int64 {
						if activeEquip != nil {
							return activeEquip.Uid
						}
						return 0
					}())
			}
		}
	}

	if objType == pb.ObjectType_PLAYER {
		p.SyncUserEquipsInfo()
	} else {
		p.SyncUserPartnerSlot()
	}
	p.tasks.dispatch(pb.TargetType_EQUIP_ITEMS_WORN_TOTAL, 0, int32(len(p.equips)))
	p.SyncUserBox()
	p.refreshStats()
}

func (p *Player) getScenePartner() []*pb.ObjectInfo {
	if p.desk == nil {
		return nil
	}

	partnerList := []*pb.ObjectInfo{}
	for _, v := range p.partnerSlots {
		if v == nil {
			continue
		}

		partner := p.getPartnerInfo(v.PartnerUid)
		if partner == nil {
			continue
		}
		if partner.desk == nil {
			continue
		}
		partnerList = append(partnerList, &pb.ObjectInfo{
			Type:      pb.ObjectType(partner.objectType),
			Uid:       partner.uid,
			CfgID:     partner.cfgId,
			Level:     partner.Level(),
			Name:      partner.Name(),
			Location:  partner.location,
			Direction: partner.direction,
			MaxHp:     partner.GetMaxHP(),
			CurHp:     partner.currentHP,
			OwnerUid:  p.uid,
		})
		logger.Debugf(" 同步房间内玩家信息  伙伴ID %d:%s  位置 %v", partner.uid, partner.name, partner.location)
	}
	return partnerList
}

// 伙伴羁绊
func (p *Player) checkPartnerBonds() {
	bond := map[int32]int32{}

	for _, v := range global.PartnerBondCfg {

		unlock := true
		for _, cp := range v.PartnerIDs {
			if par := p.getPartnerInfoByCfgId(cp); par != nil {
				if par.star < v.Star {
					unlock = false
					break
				}
			} else {
				unlock = false
				break
			}
		}

		if unlock {
			if b, ok := bond[v.BondID]; ok {
				if b < v.ID {
					bond[v.BondID] = v.ID
				}
			} else {
				bond[v.BondID] = v.ID
			}
		}
	}
	p.activateBonds = []int32{}
	for _, v := range bond {
		p.activateBonds = append(p.activateBonds, v)
	}
}

// 增加宠物
func (p *Player) addPet(cfgId int32, exp int32, star int32, toSlotId int32, NeedDel int32, Channel string) (newSort int32, err error) {
	//if cfgId == 0 {
	//	return 0, nil
	//}
	//endPetSlot := toSlotId
	//// par := CreatePet(cfgId)
	//// par.owner = p
	//petSlot := p.checkEmptyPetSlot()
	//if petSlot == -1 {
	//	return 0, errutil.ErrPetFull
	//}
	//// 判断这个等级能开多少槽位
	//BagInfo := global.GlobalCfg.Petbag
	//canhave := int32(0)
	//for _, v := range BagInfo {
	//	// 判断当前等级能开多少槽位
	//	if p.level >= int32(v.Level) {
	//		canhave = int32(v.Value)
	//	} else {
	//		continue
	//	}
	//}
	//// 假设他现在能开10个，给到key是10，也是可以放进去的，因为0是上阵的
	//if petSlot > canhave {
	//	return 0, errutil.ErrPetFull
	//}
	//// cfgid对应的宠物信息赋值给petslot
	//petCfg := global.GetPetCfg(cfgId)
	//if petCfg == nil {
	//	return 0, errutil.ErrConfigNotFound
	//}
	//petExpForLevel := int32(1)
	//newExp := int32(0)
	//if exp > 0 {
	//	petExpForLevel, newExp, err = p.petExpForLevel(1, int32(exp), petCfg.Quality)
	//	if err != nil {
	//		return 0, err
	//	}
	//}
	//if toSlotId == 0 {
	//	endPetSlot = petSlot
	//}
	//p.pet[endPetSlot] = PetSlotInfo{
	//	PetCfgId: cfgId,
	//	Quality:  petCfg.Quality,
	//	Level:    petExpForLevel,
	//	Exp:      newExp,
	//	Star:     star,
	//	ItemId:   petCfg.ItemID,
	//	ExpStar:  0,
	//}
	//p.tasks.dispatch(pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL, 0, petExpForLevel)
	//// 判断p.petEverHave有没有，没有的话，就添加新的
	//if p.petEverHave == nil {
	//	p.petEverHave = []int32{cfgId}
	//} else {
	//	nochange := 0
	//	for _, v := range p.petEverHave {
	//		if v == cfgId {
	//			nochange++
	//			continue
	//		}
	//	}
	//	if nochange == 0 {
	//		p.petEverHave = append(p.petEverHave, cfgId)
	//	}
	//}
	//p.syncSavePetDb()
	//p.SyncUserPet(NeedDel)
	//if Channel != "合成" {
	//	p.upPetAddShushuDate(petCfg, Channel)
	//}
	//return petSlot, nil
	return p.PetMgr.addPet(cfgId, exp, star, toSlotId, NeedDel, Channel)
}

//func (p *Player) upPetAddShushuDate(petCfg *global.PetConfig, Channel string) {
//	shushuDate := make(map[string]interface{}, 0)
//	baseInfo := p.ToMapForShushu()
//	for k, v := range baseInfo {
//		shushuDate[k] = v
//	}
//	zun := false
//	if petCfg.Mark == 1 {
//		zun = true
//	}
//	shushuDate["name"] = "pet_catch"
//	shushuDate["item_id"] = petCfg.ID
//	shushuDate["item_name"] = petCfg.Name
//	shushuDate["item_quality"] = petCfg.Quality
//	shushuDate["is_zun"] = zun
//	shushuDate["method"] = Channel
//	shushumain.ReportData(shushuDate)
//}

//// 校验是否有空的宠物槽位 有的话返回具体的槽位
//func (p *Player) checkEmptyPetSlot() int32 {
//	// 判断是否有空的槽位
//	for k, v := range p.pet {
//		if v.PetCfgId == 0 {
//			if k == 0 {
//				continue
//			}
//			return int32(k)
//		}
//	}
//	return -1
//}
//
//// 获取空余槽位数量
//func (p *Player) getEmptyPetSlotNum() int {
//	num := 0
//	for k, v := range p.pet {
//		if k == 0 {
//			continue
//		}
//
//		if v.PetCfgId == 0 {
//			num++
//		}
//	}
//	return num
//}

// 宠物信息
func (p *Player) getPetInfo(SlotId int32) *PetSlotInfo {
	//for i, pet := range p.pet {
	//	if i == int(SlotId) {
	//		return &pet
	//	}
	//}
	//return nil
	return p.PetMgr.getPetInfo(SlotId)
}

// 宠物升星
func (p *Player) petStar(slotId int32, cfgId int32, num int32) (resNum int32, err error) {
	//resper := int32(0)
	//EndExp := int64(p.pet[SlotId].ExpStar)
	//needLevel := int32(0)
	//pet := p.getPetInfo(SlotId)
	//if pet == nil {
	//	return resper, errutil.ErrPetInvalid
	//}
	//
	//petCfg := global.GetPetCfg(pet.PetCfgId)
	//if petCfg == nil {
	//	return resper, errutil.ErrConfigNotFound
	//}
	//
	//if pet.Star >= int32(len(petCfg.StarConsume)) {
	//	return resper, errors.New("已满阶")
	//}
	//// 删除材料
	//if !p.checkTakeSource(CfgId, int64(Num)) {
	//	return resper, errutil.ErrItemNotEnough
	//}
	//
	//oldExpStar := p.pet[SlotId].ExpStar
	//newExpStar := oldExpStar
	//oldStar := p.pet[SlotId].Star
	//// 查询宠物当前的星级所需要的材料
	//for _, v := range petCfg.StarConsume {
	//	// 查询这个道具对应多少经验
	//	// perAll := make([]int, 0)
	//	for i := 0; i <= int(Num); i++ {
	//		if v.Id-1 == pet.Star {
	//			PetAdvanceInfo := global.GlobalCfg.PetAdvance.Basevalue.Item
	//			slot := int64(0)
	//			for _, v2 := range PetAdvanceInfo {
	//				if v2.ID == CfgId {
	//					slot = int64(v2.Num)
	//					break
	//				}
	//			}
	//			per := petStageRand()
	//			resper = int32(per)
	//			if per == -1 {
	//				return 0, errutil.ErrConfigNotFound
	//			}
	//			// perAll = append(perAll, per)
	//			EndExp = int64(p.pet[SlotId].ExpStar) + slot*int64(per)
	//			// 查询这个道具对应多少经验
	//			if EndExp >= int64(v.Values) {
	//				EndExp -= int64(v.Values)
	//				for _, v3 := range petCfg.StarLimit {
	//					if v3.Id == pet.Star+1 {
	//						needLevel = v3.Values
	//						break
	//					}
	//				}
	//				// 如果当前宠物等级不够，不能升星
	//				if pet.Level < needLevel {
	//					// 如果已经满经验了，就不扣除材料
	//					if p.pet[SlotId].ExpStar != v.Values {
	//						EndExp = int64(v.Values)
	//					} else {
	//						return 0, errutil.ErrPetLevelNotEnough // 宠物等级不足
	//					}
	//					break
	//				} else {
	//					p.pet[SlotId].Star++
	//					break
	//				}
	//			} else {
	//				break
	//			}
	//			// pet.ExpStar+
	//		}
	//	}
	//}
	//newStar := p.pet[SlotId].Star
	//p.takeSource(CfgId, int64(Num), constant.ItemDelPetStar)
	//p.pet[SlotId].ExpStar = int32(EndExp)
	//newExpStar = p.pet[SlotId].ExpStar
	//p.syncSavePetDb()
	//p.SyncUserPet(0)
	//oldcombat := p.totalCombat
	//p.refreshStats()
	//newcombat := p.totalCombat
	//// 升星数据上报
	//p.upPetStarShushuDate(petCfg, CfgId, Num, oldExpStar, newExpStar, oldStar, newStar, newcombat-oldcombat)
	//return resper, nil
	return p.PetMgr.petStar(slotId, cfgId, num)
}

func (p *Player) upPetStarShushuDate(petCfg *global.PetConfig, CfgId int32, Num int32, oldExpStar int32, newExpStar int32, oldStar int32, newStar int32, combatDiff int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	zun := false
	if petCfg.Mark == 1 {
		zun = true
	}

	shushuDate["name"] = "pet_overshoot"
	shushuDate["item_id"] = petCfg.ID
	shushuDate["item_name"] = petCfg.Name
	shushuDate["item_quality"] = petCfg.Quality
	shushuDate["is_zun"] = zun
	// 查询config
	shushuDate["cost_item_name"] = CfgId
	shushuDate["cost_item_num"] = Num
	shushuDate["rank_before"] = oldStar
	shushuDate["rank_after"] = newStar
	shushuDate["progress_before"] = oldExpStar // 可根据实际情况填写进度
	shushuDate["progress_after"] = newExpStar  // 可根据实际情况填写进度
	shushuDate["power_up"] = combatDiff
	shushumain.ReportData(shushuDate)
}

// 宠物升星暴击
func petStageRand() int {
	petStage := global.GlobalCfg.PetAdvance.Rate.Stage
	allPriority := 0
	beginPriority := 0
	// 活动内容
	for _, petStageInfo := range petStage {
		allPriority = allPriority + int(petStageInfo.Pro)
	}
	rand := global.RandIn(allPriority) + 1
	for _, petStageInfo := range petStage {
		if beginPriority+int(petStageInfo.Pro) >= rand {
			return petStageInfo.Rate
		}
		beginPriority = beginPriority + int(petStageInfo.Pro)
	}
	return -1
}

func (p *Player) petLevel(slotId int32) error {
	//pet := p.getPetInfo(SlotId)
	//if pet == nil {
	//	return errutil.ErrPetInvalid
	//}
	//
	//petCfg := global.GetPetCfg(pet.PetCfgId)
	//if petCfg == nil {
	//	return errutil.ErrConfigNotFound
	//}
	//roleCfg, _ := global.GetRoleCfg(constant.Pet, pet.Level)
	//if roleCfg == nil {
	//	return errutil.ErrConfigNotFound
	//}
	//if pet.Star >= int32(len(petCfg.StarConsume)) {
	//	return errutil.ErrMaxLevel
	//}
	//// 查询宠物当前的等级所需要的经验
	//per := global.GetpetLevelRate(petCfg.Quality)
	//NeedExp := math.Round(float64(roleCfg.NeedExp)*per) - float64(pet.Exp)
	//if !p.checkTakeSource(constant.ItemTypePetExp, int64(NeedExp)) {
	//	return errutil.ErrItemNotEnough
	//}
	//p.takeSource(constant.ItemTypePetExp, int64(NeedExp), constant.ItemDelPetLevel)
	//p.pet[SlotId].Level++
	//p.pet[SlotId].Exp = 0
	//p.syncSavePetDb()
	//p.SyncUserPet(0)
	//p.refreshStats()
	//return nil
	return p.PetMgr.petLevel(slotId)
}

// 升星查询这个经验对应多少等级
func (p *Player) petExpForStar(level int32, exp int32) (newLevel int32, newExp int32, err error) {
	// for i := 0; i < len(global.RoleCfg); i++ {
	// 	petCfg := global.GetPetCfg(pet.PetCfgId)
	// 	if err != nil {
	// 		return 0, 0, err
	// 	}
	// 	if int64(exp) >= RoleCfg.NeedExp {
	// 		exp -= int32(RoleCfg.NeedExp)
	// 		level++
	// 	} else {
	// 		return level, exp, nil
	// 	}
	// }
	return level, exp, nil
}

//// 升级查询这个经验对应多少等级
//func (p *Player) petExpForLevel(level int32, exp int32, quality int32) (newLevel int32, newExp int32, err error) {
//	for i := 0; i < len(global.RoleCfg); i++ {
//		RoleCfg, err := global.GetRoleCfg(constant.Pet, level)
//		if err != nil {
//			return 0, 0, err
//		}
//		// 每个伙伴品质对应不同的系数
//		per := global.GetpetLevelRate(quality)
//		if float64(exp) >= math.Round(float64(RoleCfg.NeedExp)*per) {
//			exp -= int32(math.Round(float64(RoleCfg.NeedExp) * per))
//			level++
//		} else {
//			return level, exp, nil
//		}
//	}
//	return level, exp, nil
//}
//
//// 查询这个经验对应多少经验
//func (p *Player) petExpForExp(level int32, exp int32, quality int32) (allExp int32, err error) {
//	allExp = int32(0)
//	for i := 1; i < len(global.RoleCfg); i++ {
//		RoleCfg, err := global.GetRoleCfg(constant.Pet, int32(i))
//		if err != nil {
//			return 0, err
//		}
//		// 这级需要的经验，如果没超过这级，不能加到经验里面
//		// 每个伙伴品质对应不同的系数
//		per := global.GetpetLevelRate(quality)
//		if int(level) > i {
//			exp += int32(math.Round(float64(RoleCfg.NeedExp) * per))
//		} else {
//			return allExp + exp, nil
//		}
//	}
//	return allExp + exp, nil
//}

func (p *Player) petAddExp(slotId int32, exp int32) (bool, error) {
	//up := false
	//pet := p.getPetInfo(SlotId)
	//if pet == nil {
	//	return up, errutil.ErrPetInvalid
	//}
	//petCfg := global.GetPetCfg(p.pet[SlotId].PetCfgId)
	//if petCfg == nil {
	//	return up, errutil.ErrConfigNotFound
	//}
	//
	//willexp := int64(exp)
	//oldLevel := p.pet[SlotId].Level
	//petExpForLevel, newExp, err := p.petExpForLevel(p.pet[SlotId].Level, int32(willexp), petCfg.Quality)
	//if err != nil {
	//	return up, err
	//}
	//p.tasks.dispatch(pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL, 0, petExpForLevel)
	//p.pet[SlotId].Level = petExpForLevel
	//p.pet[SlotId].Exp = newExp
	//// pet.Star++
	//if oldLevel != p.pet[SlotId].Level {
	//	up = true
	//}
	//p.syncSavePetDb()
	//p.SyncUserPet(0)
	//return up, nil
	return p.PetMgr.petAddExp(slotId, exp)
}

// 宠物信息
func (p *Player) petMove(fromSlotId, toSlotId int32) (newSlot int32, newPet int32, newPetCfgId int32, err error) {
	//if fromSlotId < 0 || fromSlotId >= int32(len(p.pet)) || toSlotId < 0 || toSlotId >= int32(len(p.pet)) {
	//	return 0, 0, 0, errutil.ErrPetInvalid
	//}
	//// 会有好几种情况，第一种是交换两个空的宠物，第二种是交换两个有宠物的宠物，第三种是交换一个空的宠物和一个有宠物的宠物
	//if p.pet[fromSlotId].PetCfgId == 0 {
	//	return 0, 0, 0, errutil.ErrPetInvalid
	//}
	//newSlot = toSlotId
	//var petCfg *global.PetConfig
	//
	//// 如果拖动的宠物是经验宝宝，那么直接给目的槽位宠物加经验
	//if fromSlotId != 0 {
	//	petCfg = global.GetPetCfg(p.pet[fromSlotId].PetCfgId)
	//}
	//var toPetCfg *global.PetConfig
	//if toSlotId != 0 {
	//	toPetCfg = global.GetPetCfg(p.pet[toSlotId].PetCfgId)
	//}
	//
	//toLevel := p.pet[toSlotId].Level
	//fromLevel := p.pet[fromSlotId].Level
	//fromStar := p.pet[fromSlotId].Star
	//toStar := p.pet[toSlotId].Star
	//// 两个位置都有人，如果他们品级都一样，直接合成
	//if p.pet[fromSlotId].PetCfgId != 0 && p.pet[toSlotId].PetCfgId != 0 {
	//	if fromSlotId != 0 && toSlotId != 0 {
	//		if petCfg.PetKind.Key == 2 {
	//			if p.pet[toSlotId].PetCfgId != 0 {
	//				// 是经验宝宝，有两个经验，一个是经验宝宝累积的经验，还有就是经验宝宝自身配表的经验
	//				fromExp, err := p.petExpForExp(p.pet[fromSlotId].Level, p.pet[fromSlotId].Exp, p.pet[fromSlotId].Quality)
	//				if err != nil {
	//					return newSlot, 0, 0, err
	//				}
	//				allExp := petCfg.PetKind.Value + p.pet[toSlotId].Exp + fromExp
	//				up, err := p.petAddExp(toSlotId, allExp)
	//				if err != nil {
	//					return newSlot, 0, 0, err
	//				}
	//				if up {
	//					p.upPetLevelShushuDate(petCfg, toPetCfg, p.pet[toSlotId].Star, toLevel, fromLevel)
	//				}
	//				p.pet[fromSlotId] = PetSlotInfo{}
	//				p.syncSavePetDb()
	//				p.SyncUserPet(0)
	//				return newSlot, 0, 0, nil
	//			} else {
	//				p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	//				p.syncSavePetDb()
	//				p.SyncUserPet(0)
	//				return newSlot, 0, 0, nil
	//			}
	//		}
	//		if p.pet[fromSlotId].Quality == p.pet[toSlotId].Quality {
	//			// 合成
	//			whileAll := make(map[int32][]int32)
	//			for _, v := range global.ComposeCfg {
	//				for _, v2 := range v.Raw {
	//					whileAll[v2] = append(whileAll[v2], v.ComposeID)
	//				}
	//			}
	//			newId := int32(0)
	//			newPetId := int32(0)
	//			if len(whileAll[p.pet[fromSlotId].ItemId]) == 0 || len(whileAll[p.pet[toSlotId].ItemId]) == 0 {
	//				p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	//				p.syncSavePetDb()
	//				p.SyncUserPet(0)
	//				return newSlot, 0, 0, errutil.ErrComposeError
	//			}
	//
	//			// 两个数组有交集，说明就是可以合成的
	//			if len(whileAll[p.pet[fromSlotId].ItemId]) > 0 && len(whileAll[p.pet[toSlotId].ItemId]) > 0 {
	//				composeIds := make([]int32, 0)
	//				for _, v := range whileAll[p.pet[fromSlotId].ItemId] {
	//					for _, v2 := range whileAll[p.pet[toSlotId].ItemId] {
	//						if v == v2 {
	//							composeIds = append(composeIds, v)
	//							break
	//						}
	//					}
	//				}
	//				if len(composeIds) == 0 {
	//					p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	//					p.syncSavePetDb()
	//					p.SyncUserPet(0)
	//					return newSlot, 0, 0, errutil.ErrComposeError
	//				}
	//				// 查找子集，优先读子集
	//				// 查找最短的Raw数组，把这个数组的第一个元素作为合成id
	//				composeId := composeIds[0]
	//				for _, v := range composeIds {
	//					if len(global.ComposeCfg[v].Raw) < len(global.ComposeCfg[composeId].Raw) {
	//						composeId = v
	//					}
	//				}
	//				logger.Debugf("合成id %d", composeId)
	//				// 合成的id
	//				if global.ComposeCfg[composeId] == nil {
	//					p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	//					p.syncSavePetDb()
	//					p.SyncUserPet(0)
	//					return newSlot, 0, 0, errutil.ErrComposeError
	//				}
	//				newId = composeRand(global.ComposeCfg[composeId])
	//				// 转换成宠物id
	//				for _, v := range global.PetCfg {
	//					if v.ItemID == newId {
	//						newPetId = v.ID
	//						break
	//					}
	//				}
	//
	//				// 两个宠物的经验相加
	//				fromExp, err := p.petExpForExp(p.pet[fromSlotId].Level, p.pet[fromSlotId].Exp, p.pet[fromSlotId].Quality)
	//				if err != nil {
	//					return newSlot, 0, 0, err
	//				}
	//				toExp, err := p.petExpForExp(p.pet[toSlotId].Level, p.pet[toSlotId].Exp, p.pet[toSlotId].Quality)
	//				if err != nil {
	//					return newSlot, 0, 0, err
	//				}
	//				allExp := fromExp + toExp
	//
	//				// 星级取最大星级
	//				star := int32(0)
	//				if p.pet[fromSlotId].Star > p.pet[toSlotId].Star {
	//					star = p.pet[fromSlotId].Star
	//				} else {
	//					star = p.pet[toSlotId].Star
	//				}
	//				fromPetLevel := p.pet[fromSlotId].Level
	//				toPetLevel := p.pet[toSlotId].Level
	//				// 去掉原来的宠物
	//				p.pet[fromSlotId] = PetSlotInfo{}
	//				p.pet[toSlotId] = PetSlotInfo{}
	//				newSlot, _ = p.addPet(newPetId, allExp, star, toSlotId, 0, "合成")
	//				p.syncSavePetDb()
	//				// 查询这个宠物的品质
	//				petInfo := global.GetPetCfg(newPetId)
	//				if petInfo.Quality >= 6 && petInfo.PetKind.Key == 1 {
	//					strSend := []string{p.name, strconv.Itoa(int(petInfo.ItemID))}
	//					defaultManager.Lamp(int32(constant.LampPetComposeSummon), p.serverID, strSend)
	//				}
	//				p.upPetComposeShushuDate(petCfg, toPetCfg, petInfo, fromPetLevel, toPetLevel, p.pet[newSlot].Level)
	//				p.tasks.dispatch(pb.TargetType_PET_COMPOSE_TIMES, 0, 1)
	//				// p.SyncUserPet(0)
	//				if len(global.ComposeCfg[composeId].Return) > 0 {
	//					for _, cv := range global.ComposeCfg[composeId].Return {
	//						if cv.Key == newId {
	//							for _, cv2 := range cv.Values {
	//								// 给玩家添加材料
	//								logger.Debugf("给玩家添加材料 %d %d", cv2.ID, cv2.Num)
	//								p.addSource(cv2.ID, int64(cv2.Num), constant.ItemAddCompose)
	//							}
	//						}
	//					}
	//				}
	//				return newSlot, newPetId, newPetId, nil
	//			} else {
	//			}
	//		}
	//	}
	//}
	//p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	//p.syncSavePetDb()
	//p.SyncUserPet(0)
	//if fromSlotId == 0 || toSlotId == 0 {
	//	oldcombat := p.totalCombat
	//	p.defaultBuff()
	//	newcombat := p.totalCombat
	//	if toSlotId == 0 {
	//		p.upPetOnShushuDate(petCfg, toPetCfg, fromLevel, toLevel, fromStar, toStar, newcombat-oldcombat)
	//	}
	//}
	//return -1, 0, 0, nil
	return p.PetMgr.petMove(fromSlotId, toSlotId)
}

//func (p *Player) upPetLevelShushuDate(petCfg *global.PetConfig, toPetCfg *global.PetConfig, star int32, tolevel int32, fromlevel int32) {
//	shushuDate := make(map[string]interface{}, 0)
//	baseInfo := p.ToMapForShushu()
//	for k, v := range baseInfo {
//		shushuDate[k] = v
//	}
//	zun := false
//	if toPetCfg.Mark == 1 {
//		zun = true
//	}
//	shushuDate["is_zun"] = zun
//	shushuDate["name"] = "pet_gobble"
//	shushuDate["item_id"] = toPetCfg.ID
//	shushuDate["item_name"] = toPetCfg.Name
//	shushuDate["item_quality"] = toPetCfg.Quality
//	shushuDate["item_rank"] = star
//
//	shushuDate["eaten_id"] = petCfg.ID
//	shushuDate["eaten_name"] = petCfg.Name
//	shushuDate["eaten_quality"] = petCfg.Quality
//
//	shushuDate["eaten_level"] = fromlevel
//	shushuDate["level_after"] = tolevel + 1
//	shushuDate["level_before"] = tolevel
//	shushumain.ReportData(shushuDate)
//}

//func (p *Player) upPetComposeShushuDate(fromPet *global.PetConfig, toPet *global.PetConfig, newPet *global.PetConfig, fromLevel int32, toLevel int32, level int32) {
//	shushuDate := make(map[string]interface{}, 0)
//	baseInfo := p.ToMapForShushu()
//	for k, v := range baseInfo {
//		shushuDate[k] = v
//	}
//	zun := false
//	if newPet.Mark == 1 {
//		zun = true
//	}
//
//	shushuDate["name"] = "pet_synthesis"
//	item_from_zun := false
//	if fromPet.Mark == 1 {
//		item_from_zun = true
//	}
//	item_to_zun := false
//	if toPet.Mark == 1 {
//		item_to_zun = true
//	}
//	shushuDate["cost_pet_info"] = make([]map[string]interface{}, 0, 2)
//	shushuDate["cost_pet_info"] = append(shushuDate["cost_pet_info"].([]map[string]interface{}), map[string]interface{}{
//		"item_id":      fromPet.ID,
//		"item_name":    fromPet.Name,
//		"item_quality": fromPet.Quality,
//		"item_level":   fromLevel,
//		"item_zun":     item_from_zun,
//	})
//	shushuDate["cost_pet_info"] = append(shushuDate["cost_pet_info"].([]map[string]interface{}), map[string]interface{}{
//		"item_id":      toPet.ID,
//		"item_name":    toPet.Name,
//		"item_quality": toPet.Quality,
//		"item_level":   toLevel,
//		"item_zun":     item_to_zun,
//	})
//	shushuDate["is_zun"] = zun
//	shushuDate["synthesis_item_id"] = newPet.ItemID
//	shushuDate["synthesis_item_name"] = newPet.Name
//	shushuDate["synthesis_item_quality"] = newPet.Quality
//	shushuDate["synthesis_item_level"] = level
//	shushumain.ReportData(shushuDate)
//}
//
//func (p *Player) upPetOnShushuDate(fromPet *global.PetConfig, toPet *global.PetConfig, fromLevel int32, toLevel int32, fromStar int32, toStar int32, combat int32) {
//	shushuDate := make(map[string]interface{}, 0)
//	baseInfo := p.ToMapForShushu()
//	for k, v := range baseInfo {
//		shushuDate[k] = v
//	}
//	tozun := false
//	fromzun := false
//	if toPet != nil && toPet.Mark == 1 {
//		tozun = true
//	}
//	if fromPet != nil && fromPet.Mark == 1 {
//		fromzun = true
//	}
//	shushuDate["name"] = "pet_equip"
//	if toPet == nil {
//		shushuDate["before_item_id"] = 0
//		shushuDate["before_item_name"] = ""
//		shushuDate["item_level_before"] = 0
//		shushuDate["item_quality_before"] = 0
//		shushuDate["before_is_zun"] = false
//	} else {
//		shushuDate["before_item_id"] = toPet.ID
//		shushuDate["before_item_name"] = toPet.Name
//		shushuDate["item_level_before"] = toLevel
//		shushuDate["item_quality_before"] = toStar
//		shushuDate["before_is_zun"] = tozun
//	}
//
//	shushuDate["is_zun"] = fromzun
//	shushuDate["item_id"] = fromPet.ID
//	shushuDate["item_name"] = fromPet.Name
//	shushuDate["item_level_after"] = fromLevel
//	shushuDate["item_rank"] = toStar
//	shushuDate["item_quality"] = fromPet.Quality
//	shushuDate["power_up"] = combat
//	shushumain.ReportData(shushuDate)
//}

//func (p *Player) syncSavePetDb() error {
//	async.Run(func() {
//		p.lock.Lock()
//		defer p.lock.Unlock()
//
//		data, err := json.Marshal(p.pet)
//		if err != nil {
//			return
//		}
//		data2, err := json.Marshal(p.petEverHave)
//		if err != nil {
//			return
//		}
//		pet := model.UserPet{
//			Uid:         p.uid,
//			PetInfo:     string(data),
//			PetEverHave: string(data2),
//		}
//
//		if err := db.UpdateUserPet(&pet); err != nil {
//			p.logger.Errorf("保存玩家宠物信息失败:%v", err)
//			return
//		}
//	})
//
//	return nil
//}

// 当前应该挑战的关卡ID
func (p *Player) getNextBattleId(CopyType int32) int32 {
	var battleId int32

	var copyType_ int32
	var chapter_ int32
	var layers_ int32

	copyType_ = CopyType
	switch CopyType {
	case constant.BattleMap:
		if p.BattleMap.BattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			copyType_ = p.BattleMap.BattleId / 100000
			chapter_ = p.BattleMap.BattleId % 100000 / 1000
			layers_ = p.BattleMap.BattleId % 1000

			cfg := global.GetBattleMapCfg(p.BattleMap.BattleId)
			if cfg == nil {
				return 0
			}
			if cfg.MapMark == 1 {
				chapter_++
				layers_ = 1

			} else {
				layers_++
			}
			battleId = copyType_*100000 + chapter_*1000 + layers_
		}
	case constant.BossMap:
		if p.BossMap.BattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			battleId = p.BossMap.BattleId + 1
		}
	case constant.TrialTowerMap:
		if p.TrialTowerMap.BattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			battleId = p.TrialTowerMap.BattleId + 1
		}
	case constant.RelicMap:
		if p.RelicMap.BattleId == 0 {
			battleId = copyType_*100000 + 1*1000 + 1
		} else {
			copyType_ = p.RelicMap.BattleId / 100000
			chapter_ = p.RelicMap.BattleId % 100000 / 1000
			layers_ = p.RelicMap.BattleId % 1000

			cfg := global.GetBattleMapCfg(p.RelicMap.BattleId)
			if cfg == nil {
				return 0
			}
			if cfg.MapMark == 1 {
				chapter_++
				layers_ = 1

			} else {
				layers_++
			}

			battleId = copyType_*100000 + chapter_*1000 + layers_
		}
	case constant.FightEndMap:
		if p.FightEndMap.BattleId == 0 {
			battleId = copyType_*100000 + 1
		} else {
			battleId = p.FightEndMap.BattleId + 1
		}

	default:
		return 0
	}
	return battleId
}

// 填充副本信息
func (p *Player) fillInstanceInfo(res *pb.InstanceInfoResponse) {
	res.BattleMap = &pb.InstanceInfoResponse_Battle{
		BattleId:    p.InstanceData.BattleMap.BattleId,
		RewardStage: p.InstanceData.BattleMap.RewardStage,
	}
	res.BossMap = &pb.InstanceInfoResponse_Boss{
		BattleId:      p.InstanceData.BossMap.BattleId,
		ChallengeTime: p.InstanceData.BossMap.ChallengeTime,
	}
	res.TrialTowerMap = &pb.InstanceInfoResponse_TrialTower{
		BattleId: p.InstanceData.TrialTowerMap.BattleId,
	}
	res.RelicMap = &pb.InstanceInfoResponse_Relic{
		BattleId:      p.InstanceData.RelicMap.BattleId,
		ChallengeTime: p.InstanceData.RelicMap.ChallengeTime,
	}

	res.PushingWall = &pb.InstanceInfoResponse_PushingWall{
		UsedCount:           p.InstanceData.PushingWall.UsedCount,
		TicketPurchaseCount: p.InstanceData.PushingWall.TicketPurchaseCount,
	}

	for _, v := range p.InstanceData.RelicMap.Evaluation {
		res.RelicMap.Evaluation = append(res.RelicMap.Evaluation, &pb.KeyValue{
			K: v.Key,
			V: v.Value,
		})
	}
	for _, v := range p.InstanceData.RelicMap.ChapterReward {
		res.RelicMap.ChapterReward = append(res.RelicMap.ChapterReward, &pb.KeyValue{
			K: v.Key,
			V: v.Value,
		})
	}

	res.Afk = &pb.InstanceInfoResponse_AFK{DailyTime: 0}
	res.Afk.DailyTime = p.afkLeftTime()
	if p.afkExpirationAt.Before(time.Now()) {
		res.Afk.ExpirationAt = 0
	} else {
		res.Afk.ExpirationAt = p.afkExpirationAt.Unix()
	}
	res.Afk.AfkMultiplierID = p.AfkMultiplierID
	res.ToEnd = &pb.InstanceInfoResponse_ToEnd{}

	if p.FightEndMap.RefreshAt.Day() != time.Now().Day() {
		res.ToEnd.BattleId = 0
		p.FightEndMap.BattleId = 0
		p.FightEndMap.RefreshAt = time.Now()

		for i := 0; i < len(p.FightEndMap.Evaluation); i++ {
			p.FightEndMap.Evaluation[i].Uid = 0
			p.FightEndMap.Evaluation[i].ChallengeTime = 0
		}
	}
	res.ToEnd.BattleId = p.FightEndMap.BattleId
	for i := 0; i < len(p.FightEndMap.Evaluation); i++ {
		res.ToEnd.Evaluation = append(res.ToEnd.Evaluation, &pb.KeyValue{
			K: p.FightEndMap.Evaluation[i].ChallengeTime,
			V: p.FightEndMap.Evaluation[i].ChallengeTime,
		})
	}
}

// 获取章节奖励
func (p *Player) getBattleReward(map_type pb.MapType, chapter_id int32) error {
	if map_type == pb.MapType_BATTLE_MAP {
		if p.BattleMap.BattleId <= p.BattleMap.RewardStage {
			return errutil.ErrRequestFailed
		}

		if p.bagItemLeftCount() <= 1 {
			return errutil.ErrBagFull
		}

		rewards := []global.IdNum{}
		for i := p.BattleMap.RewardStage; i <= p.BattleMap.BattleId; i++ {
			if cfg := global.GetBattleMapCfg(i); cfg != nil {
				for _, v := range cfg.ChapterAward {
					rewards = append(rewards, global.IdNum{
						ID:  v.ID,
						Num: v.Num,
					})
				}
			}
		}

		p.addSourceBatch(rewards, constant.ItemAddDrop)

		p.BattleMap.RewardStage = p.BattleMap.BattleId

	} else if map_type == pb.MapType_RELIC_MAP {
		rewards := []global.IdNum{}
		chapter_star := int32(0) // 本章节以获得的星
		for _, v := range p.RelicMap.Evaluation {
			if v.Key%100000/1000 != chapter_id {
				continue
			}
			chapter_star += v.Value
		}

		get_star := int32(0) // 已经领取的星
		// get_k := 0

		for _, v := range p.RelicMap.ChapterReward {
			if v.Key == chapter_id {
				if v.Value >= chapter_star {
					return errutil.ErrRequestFailed
				} else {
					get_star = v.Value
					// get_k = k
					break
				}
			}
		}

		if get_star >= chapter_star {
			return errutil.ErrRequestFailed
		}

		var cfg []global.KeyValueEx2
		for _, v := range global.BattleCfg {
			if v.CopyType == int32(map_type) && v.Chapter == chapter_id && len(v.StarAward) > 0 {
				cfg = v.StarAward
				break
			}
		}
		if len(cfg) == 0 {
			return errutil.ErrRequestFailed
		}

		max_get_star := int32(0) // 最新领取的星数
		for _, v := range cfg {
			if v.Key > get_star && v.Key <= chapter_star {
				for _, v2 := range v.Values {
					rewards = append(rewards, global.IdNum{
						ID:  v2.ID,
						Num: v2.Num,
					})
				}
				max_get_star = v.Key
			}
		}

		bFind := false
		for i := 0; i < len(p.RelicMap.ChapterReward); i++ {
			if p.RelicMap.ChapterReward[i].Key == chapter_id {
				p.RelicMap.ChapterReward[i].Value = max_get_star
				bFind = true
				break
			}
		}
		if !bFind {
			p.RelicMap.ChapterReward = append(p.RelicMap.ChapterReward, global.KeyValue{
				Key:   chapter_id,
				Value: max_get_star,
			})
		}

		p.addSourceBatch(rewards, constant.ItemAddDrop)
		p.asyncSaveDb()

		// 跑马灯
		if max_get_star >= 30 {
			context := []string{p.name}
			defaultManager.Lamp(int32(constant.LampRelicChapterReward), p.serverID, context)
		}
	}

	return nil
}

func (p *Player) battleMapCopySweep(roomID, battleID, sweepTimes int32) (int32, string) {
	var code int32
	var err string

	roomCfg := global.GetRoomCfg(roomID)
	if roomCfg == nil {
		code = -1
		err = errutil.ErrConfigNotFound.Error()
		return code, err
	}

	battleCfg := global.GetBattleMapCfg(battleID)
	if battleCfg == nil {
		code = -2
		err = errutil.ErrConfigNotFound.Error()
		return code, err
	}

	switch roomCfg.Roomtype {
	case constant.BossMap:
		if p.BossMap.BattleId < battleID {
			return -3, errutil.ErrRequestFailed.Error()
		}
		if p.BossMap.ChallengeTime < sweepTimes {
			return -4, errutil.ErrRequestFailed.Error()
		}

		// sweepTimes = p.BossMap.ChallengeTime
		p.BossMap.ChallengeTime -= sweepTimes
		p.tasks.dispatch(pb.TargetType_TOTAL_BOSS_CHALLENGE, 0, sweepTimes)
	case constant.RelicMap:
		if p.RelicMap.ChallengeTime < sweepTimes {
			return -4, errutil.ErrRequestFailed.Error()
		}
		for _, v := range p.RelicMap.Evaluation {
			if v.Key == battleID {
				if v.Value < 3 {
					return -5, errutil.ErrRequestFailed.Error()
				}
				break
			}
		}
		p.RelicMap.ChallengeTime -= sweepTimes

		p.tasks.dispatch(pb.TargetType_TOTAL_RELIC_CHALLENGE, 0, sweepTimes)
	}

	rewards := []global.IdNum{}

	for i := 0; i < int(sweepTimes); i++ {
		rewards = append(rewards, battleCfg.Settlement...)
	}

	if err := p.checkAddSourceBatch(rewards); err == nil {
		p.addSourceBatch(rewards, constant.ItemAddDailymap)
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range rewards {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		p.SendMail(constant.ItemAddDailymap, attachment)
	}

	packet := &pb.GameResultBattlePush{BattleId: battleID, OverType: pb.GameOverType_GameOverTypeCopySweep, Items: []*pb.ItemInfo{}}

	for _, v := range rewards {
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	}
	p.pushMsg(protocol.GameResultBattleMapPush, packet)

	p.reward = rewards
	p.roomType = roomCfg.Roomtype
	return code, err
}

// 按权重比例把 N 只怪分配到各 ID
func allocateByWeight(src []global.IdNumP, total int32) map[int32]int32 {
	if total <= 0 || len(src) == 0 {
		return nil
	}

	// 1. 计算总权重
	var sumW int64
	for _, v := range src {
		sumW += int64(v.Num)
	}
	if sumW == 0 { // 权重全 0 就平均分配
		sumW = int64(len(src))
	}

	// 2. 把 total 按权重切成碎片
	remain := int64(total)
	res := make(map[int32]int32, len(src))

	for i, v := range src {
		w := int64(v.Num)
		if sumW != int64(len(src)) && w == 0 {
			continue // 权重为 0 的不参与
		}

		// 最后 1 个直接拿剩余，防止 rounding error
		if i == len(src)-1 {
			res[v.ID] = int32(remain)
			break
		}

		// 向下取整
		n := w * int64(total) / sumW
		// 用随机数做四舍五入，保证期望正确
		if int64(global.RandIn(int(sumW))) < (w*int64(total))%sumW {
			n++
		}
		if n > remain {
			n = remain
		}
		res[v.ID] = int32(n)
		remain -= n
	}
	return res
}

// 扫荡挂机副本
func (p *Player) guardIslandSweep(roomID int32) (int32, string) {
	var code int32
	var err string

	roomCfg := global.GetRoomCfg(roomID)
	if roomCfg == nil {
		code = -1
		err = errutil.ErrConfigNotFound.Error()
		return code, err
	}
	mapCfg := global.GetMapCfg(roomCfg.MapId)
	if mapCfg == nil {
		code = -2
		err = errutil.ErrConfigNotFound.Error()
		return code, err
	}

	// 特权等级
	bPrivilege := false
	for _, v := range p.privilege {
		if privilegeCfg := global.GetPrivilegeCfg(v.Pid); privilegeCfg != nil {
			if privilegeCfg.QuicklyIsiand > 0 {
				bPrivilege = true
				break
			}
		}
	}
	if !bPrivilege {
		code = -3
		err = errutil.PrivilegeNotEnough
		return code, err
	}

	// 剩余时间杀怪数量，四舍五入
	var sweeMonNumItem, sweeMonNumExp float64
	for _, v := range global.GlobalCfg.GuardIsland.Quick {
		if p.level >= v.Minlevel && p.level <= v.Maxlevel {
			sweeMonNumItem = math.Round(float64(p.afkLeftTime()) / float64(v.Battle))
			sweeMonNumExp = math.Round(float64(p.afkLeftTime()) / float64(v.Expbattle))
			break
		}
	}
	monMap := allocateByWeight(mapCfg.MonsterIds, int32(sweeMonNumItem))
	monMapExp := allocateByWeight(mapCfg.MonsterIds, int32(sweeMonNumExp))

	p.logger.Debugf("%s 扫荡挂机副本，时间 %d  道具怪数量 %f   经验怪总数量 %f", p.name, p.afkLeftTime(), sweeMonNumItem, sweeMonNumExp)

	rewards := []global.IdNum{}
	for k, v := range monMap {
		rewards = append(rewards, monsterMgr.MonsterDorp(k, v)...)
	}
	for k, v := range monMapExp {
		if mcfg := global.GetMonsterCfg(k); mcfg != nil {
			rewards = append(rewards, global.IdNum{ID: constant.ItemTypeExp, Num: mcfg.Experience * v})
			logger.Debugf("扫荡挂机副本怪物 %d %s  基础经验 %d  怪物数量 %d  总经验 %d", mcfg.ID, mcfg.Name, mcfg.Experience, v, mcfg.Experience*v)
		}
	}

	rewards = mergeItems(rewards)

	// 守卫海岛挂机倍率
	var multiplier int32 = 0
	if guardIslandHangup, exist := global.GuardIslandCfg[p.AfkMultiplierID]; exist {
		// 当倍率为-1时，代表前端需要弹窗, 这个时候倍率设置为0
		if guardIslandHangup.Multiplier == -1 {
			multiplier = 0
		} else {
			multiplier = guardIslandHangup.Multiplier
		}
	} else {
		logger.Errorf("守卫海岛挂机倍率配置不存在, UserId:%d", p.uid)
	}

	for i := 0; i < len(rewards); i++ {
		rewards[i].Num += int32(float64(multiplier) / BaseRate * float64(rewards[i].Num))
	}

	//包裹空间
	if err := p.checkAddSourceBatch(rewards); err == nil {
		p.addSourceBatch(rewards, constant.ItemAddDailymap)
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range rewards {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		p.SendMail(constant.ItemAddDailymap, attachment)
	}

	packet := &pb.GameResultBattlePush{Items: []*pb.ItemInfo{}, RoomId: roomCfg.ID, OverType: pb.GameOverType_GameOverTypeCopySweep}
	for _, v := range rewards {
		packet.Items = append(packet.Items, &pb.ItemInfo{CfgId: v.ID, Num: int64(v.Num)})
	}
	p.pushMsg(protocol.GameResultBattleMapPush, packet)

	p.delAfkLeftTime(int64(p.afkLeftTime()))
	return code, err
}

// 挂机剩余时间
//	func (p *Player) afkLeftTime() int32 {
//		leftTime := global.GlobalCfg.BaseValue.Expcopy.Dailytime - int32(p.afkDailyTime)
//		if leftTime < 0 {
//			leftTime = 0
//		}
//		return leftTime
//	}

// 挂机剩余时间
func (p *Player) afkLeftTime() int32 {
	// 获取每日使用的挂机卡次数
	afkCardUsageCount, ok := p.DailyItemUsed[constant.ItemTypeAfkCard]
	// 默认剩余时间为配置的每日挂机时间减去已挂机时间
	baseTime := global.GlobalCfg.BaseValue.Expcopy.Dailytime
	usedTime := int32(p.afkDailyTime)
	leftTime := baseTime - usedTime

	// 如果没有使用挂机卡，直接返回剩余时间（不能为负）
	if !ok {
		if leftTime < 0 {
			return 0
		}
		return leftTime
	}

	// 获取挂机卡的配置
	itemUseCfg, err := global.GetItemUseCfg(constant.ItemTypeAfkCard)
	if err != nil || len(itemUseCfg.GiveGoods) == 0 {
		logger.Warnf("GetItemUseCfg err: %v", err)
		if leftTime < 0 {
			return 0
		}
		return leftTime
	}

	// 计算挂机卡总可用时间
	cardTotalTime := int32(afkCardUsageCount) * int32(itemUseCfg.GiveGoods[0].Num) * 60
	leftTime = baseTime + cardTotalTime - usedTime
	if leftTime < 0 {
		return 0
	}
	return leftTime
}

func (p *Player) delAfkLeftTime(del_sec int64) {
	p.afkDailyTime += del_sec
	p.dirty(true)

	p.logger.Debugf("%s 扣除挂机副本时间 %d ", p.name, del_sec)
}

func (p *Player) ToMapForShushu() map[string]interface{} {
	TmpPetID := int64(0)
	TmpPetLevel := int32(0)
	TmpPetQuality := int32(0)
	TmpMonthCard := false
	TmpLifetimeCard := false
	TmpPartnerLevel := int32(0)
	TmpPetName := ""

	bPet := p.PetMgr.GetBattlePet()
	//if len(p.pet) > 0 {
	TmpPetID = int64(bPet.PetCfgId)
	TmpPetLevel = int32(bPet.Level)
	TmpPetQuality = int32(bPet.Quality)
	//}
	zun := false
	// 查询宠物信息
	if petCfg := global.GetPetCfg(int32(TmpPetID)); petCfg != nil {
		TmpPetName = petCfg.Name
		if petCfg.Mark == 1 {
			zun = true
		}
	}

	for _, v := range p.privilege {
		cfg, ok := global.PrivilegeCfg[int32(v.Pid)]
		if !ok {
			continue
		}
		if cfg.Duration < 31 {
			TmpMonthCard = true
		} else {
			TmpLifetimeCard = true
		}
	}
	// for _, v := range p.partnerSlots {
	// 	if v.PartnerUid > 0 {
	// 		if par := p.getPartnerInfo(v.PartnerUid); par != nil {
	// 			TmpPartnerLevel += par.Level()
	// 		}
	// 		break
	// 	}
	// }
	// 查询pvp排名
	rankInfo1 := &protocol.RankOperateRequest{
		Uid:  p.uid,
		Type: 3,
		All:  3,
	}
	res1, err := provide.RankList(rankInfo1)
	tmp_pvp1_rank := int64(0)
	if err == nil {
		if res1.RankOwn.Num > 0 {
			tmp_pvp1_rank = res1.RankOwn.Num
		}
	}

	// 查询pvp排名
	tmp_pvp2_rank := int64(0)
	// rankInfo2 := &protocol.RankOperateRequest{
	// 	Uid:  p.uid,
	// 	Type: 3,
	// 	All:  3,
	// }
	// res2, err := provide.RankList(rankInfo2)
	// if err == nil {
	// 	if res2.RankOwn.Num > 0 {
	// 		tmp_pvp1_rank = res2.RankOwn.Num
	// 	}
	// }

	// 查询任务

	return map[string]interface{}{
		"channel":              "", // 填写渠道信息
		"server_id":            strconv.Itoa(int(p.serverID)),
		"nick_name":            p.name,
		"tmp_power":            int32(p.totalCombat),
		"tmp_level":            int32(p.Level()),
		"tmp_box_level":        int32(p.box.NowLevel),
		"tmp_main_id":          p.InstanceData.BattleMap.BattleId,
		"tmp_task_id":          p.taskProgress,
		"tmp_coin":             int32(p.tokens[constant.ItemTypeCoin]),
		"tmp_coin_cost":        int32(p.tokenDeletes[constant.ItemTypeCoin]),
		"tmp_ingot":            int32(p.tokens[constant.ItemTypeCoupon]),
		"tmp_ingot_cost":       int32(p.tokenDeletes[constant.ItemTypeCoupon]),
		"tmp_battle_soul":      int32(p.tokens[constant.ItemTypeSoul]),
		"tmp_battle_soul_cost": int32(p.tokenDeletes[constant.ItemTypeSoul]),
		"tmp_trial_coin":       int32(p.tokens[constant.ItemTypeTrialVoucher]),
		"tmp_trial_coin_cost":  int32(p.tokenDeletes[constant.ItemTypeTrialVoucher]),
		"tmp_payment_amount":   float64(p.recharge),
		"tmp_days":             int32(p.nowDay),
		"tmp_vip":              int32(p.vip),
		"tmp_guild_id":         int64(p.guildId),
		"tmp_time":             int64(p.totalOnlineTime),
		"tmp_role_job":         constant.JobName[p.job],
		"tmp_account":          p.accountName,
		"tmp_device":           "",
		"tmp_role_id":          p.uid,
		"tmp_pet_id":           TmpPetID,
		"tmp_pet_name":         TmpPetName, // 填写宠物名信息
		"tmp_pet_level":        TmpPetLevel,
		"tmp_pet_rank":         TmpPetQuality,
		"tmp_partner_level":    TmpPartnerLevel,
		"tmp_pvp1_rank":        tmp_pvp1_rank,
		"tmp_pvp2_rank":        tmp_pvp2_rank,
		"tmp_month_card":       TmpMonthCard,
		"tmp_lifetime_card":    TmpLifetimeCard,
		"tmp_pet_zun":          zun,
	}
}

func (p *Player) OfflineAdd() error {
	if err := p.syncBoxFromDB(); err != nil {
		return err
	}
	//if err := p.syncPetFromDB(); err != nil {
	//	return err
	//}
	// 工会信息
	if err := p.syncGuildFromDB(); err != nil {
		return err
	}
	return nil
}

// 添加挑战记录
func (p *Player) addChallengeRecord(roomtype int32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.InstanceData.ChallengeRecord.Records[roomtype]++
}

// 查询是否有挑战记录
func (p *Player) getChallengeRecord(roomtype int32) int32 {
	p.lock.Lock()
	defer p.lock.Unlock()
	if r, ok := p.InstanceData.ChallengeRecord.Records[roomtype]; ok {
		return r
	}
	return 0
}

func (p *Player) SetMailCounter(op constant.MailCounterOPType, num int) error {
	switch op {
	case constant.MailCounterOPTypeAdd:
		p.mailCounter += num
	case constant.MailCounterOPTypeSub:
		p.mailCounter -= num
	case constant.MailCounterOPTypeSet:
		p.mailCounter = num
	}
	if p.mailCounter < 0 {
		p.mailCounter = 0
	}
	return nil
}

// GuardIslandMultiplier 挂机副本
func (p *Player) GuardIslandMultiplier(id int64) error {
	guardIslandRecord, exist := global.GuardIslandCfg[id]
	if !exist {
		return errutil.NewBizError(-1, fmt.Sprintf("config.GuardIsland.id not found: %d, uid: %d", id, p.uid), errutil.ErrConfigNotFound.Error())
	}

	if p.AfkMultiplierID != -1 {
		return errutil.NewBizError(-1, fmt.Sprintf("config.GuardIsland already set, id: %d, uid: %d", p.AfkMultiplierID, p.uid), errutil.ErrRequestFailed.Error())
	}

	logger.Debugf("config.GuardIsland: %+v, uid: %d", guardIslandRecord, p.uid)

	// 处理逻辑
	switch constant.GuardIsLandType(guardIslandRecord.Type) {
	case constant.GuardIsLandTypeFree:
		p.AfkMultiplierID = int64(guardIslandRecord.ID)
	case constant.GuardIsLandTypeAd:
		// TODO: 处理广告逻辑, 现在点了直接过
		p.AfkMultiplierID = int64(guardIslandRecord.ID)
	case constant.GuardIsLandTypeItem:
		// storeConfig, ok := global.StoreCfg[guardIslandRecord.ID]
		// if !ok {
		// 	return errutil.NewBizError(-1, fmt.Errorf("StoreCfg.id not found: %d, uid: %d", guardIslandRecord.ID, p.uid), errutil.ErrConfigNotFound.Error())
		// }
		// if !p.checkTakeSource(int32(storeConfig.Limit.NeedMoney), storeConfig.Limit.Num) {
		// 	return errutil.NewBizError(-1, fmt.Errorf("config.GuardIsland checkAddSourceBatch failed, item_id: %d, item_num: %d, uid: %d", storeConfig.Limit.NeedMoney, storeConfig.Limit.Num, p.uid), errutil.ErrItemNotEnough.Error())
		// }
		// p.takeSource(int32(storeConfig.Limit.NeedMoney), storeConfig.Limit.Num, constant.ItemDelGuardIsland)
		p.AfkMultiplierID = int64(guardIslandRecord.ID)
	case constant.GuardIsLandTypeRecharge:
		// TODO : 处理充值逻辑
		p.AfkMultiplierID = int64(guardIslandRecord.ID)
		// coreEnv := strings.ToLower(viper.GetString("core.env"))
		// if coreEnv == constant.EnvDev || coreEnv == constant.EnvTest {
		// 	p.AfkMultiplierID = int64(guardIslandRecord.ID)
		// } else {
		// 	return errutil.NewBizError(-1, fmt.Errorf("config.GuardIsland recharge not support in env: %s, uid: %d", coreEnv, p.uid), errutil.ErrRequestFailed.Error())
		// }
	default:
		return errutil.NewBizError(-1, fmt.Sprintf("config.GuardIsland.type not found for id: %d, uid: %d", id, p.uid), errutil.ErrRequestFailed.Error())
	}
	return nil
}

// 添加上报记录
func (p *Player) addUpRecord(roomtype string) {
	val, _ := p.firstUp.Get(roomtype)
	p.firstUp.Set(roomtype, val+1)
	p.syncUserBase()
}
func (p *Player) addUpRecordSend(roomtype string) {
	val, _ := p.firstUp.Get(roomtype)
	p.firstUp.Set(roomtype, val+1)
}
func (p *Player) addUpRecordByAmount(roomtype string, amount int64) {
	val, _ := p.firstUp.Get(roomtype)
	p.firstUp.Set(roomtype, val+amount)
	p.syncUserBase()
}
func (p *Player) addUpRecordByAmountTotal(roomtype string, amount int64) {
	p.firstUp.Set(roomtype, amount)
	p.syncUserBase()
}

// 查询是否有挑战记录
func (p *Player) getUpRecord(roomtype string) int64 {
	if r, ok := p.firstUp.Get(roomtype); ok {
		return r
	}
	return 0
}

func (p *Player) elementInfo() *pb.ElementInfoResponse {
	//查询他的主元素
	elementMain := pb.ElementSlot(p.ElementInfo.elementMain)
	elementMainValues := p.ElementInfo.element[pb.ElementSlot_FireEmpower].Level
	empower := &pb.ElementEmpowerList{}
	for k, v := range p.ElementInfo.element {
		//1-5是属强
		if pb.ElementSlot(k) <= pb.ElementSlot_EarthEmpower {
			if v.Level > elementMainValues {
				elementMain = pb.ElementSlot(k)
				elementMainValues = v.Level
				p.ElementInfo.elementMain = int32(k)
			}
		}
		empower.Elements = append(empower.Elements, &pb.ElementEmpower{
			SlotId:     int32(k),
			Level:      v.Level,
			BreakLevel: v.Break,
		})
	}
	uids := &pb.ElementUids{}
	for _, v := range p.ElementInfo.elementInwork {
		uids.Uids = append(uids.Uids, int64(v))
	}
	return &pb.ElementInfoResponse{
		ElementScore: p.ElementMap.BattleId,
		ElementTimes: p.ElementMap.ChallengeTime,
		ElementMain:  int32(elementMain),     // 主元素，需根据实际逻辑赋值
		Empower:      empower,                // 赋能列表，需根据实际逻辑赋值
		NeedSweep:    p.ElementMap.NeedSweep, // 是否需要扫荡，需根据实际逻辑赋值
		Uids:         uids,
	}
}

func (p *Player) elementUpgrade(slotId pb.ElementSlot) (bool, int32, int32, int32, error) {
	elementMain := p.ElementInfo.elementMain

	//判断槽位是否正确
	if slotId < pb.ElementSlot_INVALID_ElementpSlot || slotId > pb.ElementSlot_MaxElementpSlot {
		return false, 0, 0, int32(elementMain), errutil.ErrElementSlot
	}
	//查询槽位对应的uid
	element := p.ElementInfo.element[slotId]
	nextEid := int32(0)
	if element.Eid == 0 {
		nextEid = int32(slotId)*1000 + 1
	} else {
		nextEid = element.Eid + 1
		//查询下一个uid对应的配置信息
		getelement := global.GetElementCfg(element.Eid)
		//判断是否处在突破情况
		if getelement.Breaklevel > 0 && element.Break != getelement.Breaklevel {
			return false, getelement.Level, getelement.Level, int32(elementMain), errutil.ErrElementSlotBrake
		}
	}
	//判断是不是满级
	getelementNext := global.GetElementCfg(nextEid)
	if getelementNext == nil {
		return false, 0, 0, int32(elementMain), errutil.ErrMaxLevel
	}
	// //查询下一个uid对应的配置信息
	// getelement := global.GetElementCfg(nextEid)
	// if getelement == nil {
	// 	return false, getelementNext.Level - 1, getelementNext.Level - 1, int32(elementMain), errutil.ErrConfigNotFound
	// }

	levelconsume := []global.IdNum{}
	for _, v := range getelementNext.Levelconsume {
		//判断材料是否足够
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return false, getelementNext.Level - 1, getelementNext.Level - 1, int32(elementMain), errutil.ErrItemNotEnough
		}
		levelconsume = append(levelconsume, global.IdNum{ID: int32(v.ID), Num: int32(v.Num)})
	}

	//扣除材料
	cost_gold := int32(0)
	cost_item_id := int32(0)
	cost_item_num := int32(0)
	cost_item_before := int64(0)
	cost_item_after := int64(0)
	cost_item_name := ""

	for _, v := range levelconsume {
		if v.ID == constant.ItemTypeCoin {
			cost_gold += int32(v.Num)
		} else {
			cost_item_id = v.ID
			cost_item_num += int32(v.Num)
			GeTNum, err := defaultManager.GetItemCount(p.uid, v.ID)
			if err != nil {
				continue
			}
			if ItemCfg, err := global.GetItemCfg(v.ID); err == nil {
				cost_item_name = ItemCfg.Name
			}
			cost_item_before = int64(GeTNum)
			cost_item_after = int64(GeTNum) - int64(v.Num)
		}
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelElementUpgrade)
	}
	//随机成功概率
	weight := global.RandIn(int(BaseRate))
	// logger.Debug("weight", weight)

	if weight > int(getelementNext.Rate) {
		p.upElementUpShushuDate(int32(slotId), cost_item_id, cost_item_name, cost_item_num, cost_item_before, cost_item_after,
			cost_gold, getelementNext.Level, false, weight)
		return false, getelementNext.Level - 1, getelementNext.Level - 1, int32(elementMain), errutil.ErrUpLevel
	}

	//升级元素
	p.ElementInfo.element[slotId].Eid = nextEid
	p.ElementInfo.element[slotId].Level = getelementNext.Level
	//查询目前最新的主元素
	//查询他的主元素
	elementMainValues := p.ElementInfo.element[elementMain].Level
	for k, v := range p.ElementInfo.element {
		//1-5是属强
		if pb.ElementSlot(k) <= pb.ElementSlot_EarthEmpower {
			if v.Level > elementMainValues {
				elementMain = int32(k)
				p.ElementInfo.elementMain = int32(k)
				elementMainValues = v.Level
			}
		}
	}
	p.tasks.dispatch(pb.TargetType_FIVE_ELEMENT_EMPOWER_TIMES, 0, 1)
	p.rankChangeTCP(constant.RankTypeFiveElementLevel, p.calcFiveElementLevelNum())
	p.upElementUpShushuDate(int32(slotId), cost_item_id, cost_item_name, cost_item_num, cost_item_before, cost_item_after,
		cost_gold, getelementNext.Level, true, weight)
	return true, getelementNext.Level - 1, getelementNext.Level, int32(elementMain), nil
}
func (p *Player) upElementUpShushuDate(wuxing_id int32, CfgId int32, cost_item_name string, Num int32, cost_item_before int64, cost_item_after int64,
	cost_gold, level int32, is_success bool, success_rate int) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "wuxing_get"
	shushuDate["wuxing_id"] = wuxing_id
	shushuDate["wuxing_level"] = level
	shushuDate["cost_item_id"] = CfgId
	shushuDate["cost_item_name"] = cost_item_name
	shushuDate["cost_item_num"] = Num
	shushuDate["cost_item_before"] = cost_item_before
	shushuDate["cost_item_after"] = cost_item_after
	shushuDate["cost_gold"] = cost_gold
	shushuDate["wuxing_after"] = level
	shushuDate["wuxing_before"] = level - 1
	shushuDate["is_success"] = is_success
	shushuDate["success_rate"] = success_rate // 如有实际概率请替换
	shushumain.ReportData(shushuDate)
}

func (p *Player) ElementBreakthrough(slotId pb.ElementSlot) (int32, int32, error) {
	//判断槽位是否正确
	if slotId < pb.ElementSlot_INVALID_ElementpSlot || slotId > pb.ElementSlot_MaxElementpSlot {
		return 0, 0, errutil.ErrElementSlot
	}
	//查询槽位对应的uid
	element := p.ElementInfo.element[slotId]
	//判断是否处在突破情况
	if element.Eid == 0 {
		return 0, 0, errutil.ErrElementSlot
	}
	//查询uid对应的配置信息
	getelement := global.GetElementCfg(element.Eid)
	if getelement == nil {
		return 0, 0, errutil.ErrConfigNotFound
	}
	//判断是否处在突破情况
	if getelement.Breaklevel > 0 && element.Break != getelement.Breaklevel-1 {
		return 0, 0, errutil.ErrElementSlotBrake
	}
	levelconsume := []global.IdNum{}
	for _, v := range getelement.Breakconsume {
		//判断材料是否足够
		if !p.checkTakeSource(v.ID, int64(v.Num)) {
			return 0, 0, errutil.ErrItemNotEnough
		}
		levelconsume = append(levelconsume, global.IdNum{ID: int32(v.ID), Num: int32(v.Num)})
	}
	//扣除材料
	for _, v := range levelconsume {
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelElementUpgrade)
	}
	//升级元素
	p.ElementInfo.element[slotId].Break++
	return p.ElementInfo.element[slotId].Break - 1, p.ElementInfo.element[slotId].Break, nil
}

func (p *Player) elementSweep() error {
	if p.ElementMap.NeedSweep != 1 {
		return errutil.ErrNoNeedSweep
	}
	//把所有的奖励都拿一遍
	res := global.GetElementTrialsSwttlementCfg(p.ElementMap.BattleId)
	if res == nil {
		return errutil.ErrConfigNotFound
	}
	//批量发送道具
	reward := []global.IdNum{}
	for k, v := range res {
		reward = append(reward, global.IdNum{ID: k, Num: v})
	}
	if err := p.checkAddSourceBatch(reward); err == nil {
		p.addSourceBatch(reward, constant.ItemAddElementSweep)
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range reward {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num)})
		}
		p.SendMail(constant.ItemAddElementSweep, attachment)
	}
	p.ElementMap.NeedSweep = 0
	p.addUpRecord(constant.InstanceNum)
	p.rankChangeTCP(constant.RankTypeElement, int64(p.ElementMap.BattleId)*100000000000+int64(p.totalCombat))
	return nil
}
func (p *Player) elementSweepForward(num int32) error {
	if p.ElementMap.NeedSweep == 1 || p.ElementMap.BattleId == 0 {
		return errutil.ErrNoNeedSweep
	}
	//查询这一层的奖励
	res := global.GetElementTrialsCfg(p.ElementMap.BattleId)
	if res == nil {
		return errutil.ErrConfigNotFound
	}
	//批量发送道具
	reward := []global.IdNum{}
	for _, v := range res.Settlement {
		reward = append(reward, global.IdNum{ID: int32(v.ID), Num: int32(v.Num * num)})
	}
	if err := p.checkAddSourceBatch(reward); err == nil {
		p.addSourceBatch(reward, constant.ItemAddElementSweep)
	} else {
		attachment := []protocol.Attachment{}
		for _, v := range reward {
			attachment = append(attachment, protocol.Attachment{UserPropsId: int64(v.ID), Num: int64(v.Num * num)})
		}
		p.SendMail(constant.ItemAddElementSweep, attachment)
	}
	//减去挑战次数
	p.ElementMap.ChallengeTime -= num
	p.addUpRecordByAmount(constant.InstanceNum, int64(num)) // 增加统计
	// p.ElementMap.NeedSweep = 0
	// p.rankChangeTCP(constant.RankTypeElement, int64(p.ElementMap.BattleId)*100000000000+int64(p.totalCombat))
	return nil
}
func (p *Player) resetPvpScore() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pvpScore = DefaultPvpScore
	p.gvgScore = DefaultPvpScore

	for _, v := range p.pvpRecorder {
		v.Status = 2
	}
}

func (p *Player) getZoneId() int32 {
	if !p.bRobot {
		if serverInfo, err := global.ServerListMgrInstance.Get(int(p.serverID)); err == nil {
			return serverInfo.ZoneId
		} else {
			p.logger.Errorf("真实玩家获取战区ID 失败 Uid : %d  Name %s ServerId %d", p.uid, p.name, p.serverID)
		}
	} else {
		if r := robotMgr.getRankRobotInfo(0, p.uid); r != nil {
			return r.zoneId
		} else {
			p.logger.Errorf("机器人获取战区ID 失败 Uid : %d  Name  %s ServerId %d", p.uid, p.name, p.serverID)
		}
	}

	return 0
}

func (p *Player) pveWarMapAddSoul(num int64) {
	if p.warDailySoul >= int64(global.GlobalCfg.War.PK.Daylimit) {
		return
	}
	if p.warDailySoul+num > int64(global.GlobalCfg.War.PK.Daylimit) {
		num = int64(global.GlobalCfg.War.PK.Daylimit) - p.warDailySoul
	}
	p.warDailySoul += num
	p.addSource(constant.ItemTypeSoul, num, constant.ItemAddWar)

	p.logger.Debugf("战场获得战魂 [%s]数量 %d  今日已获得 %d", p.name, num, p.warDailySoul)
}

// titleCompose
func (p *Player) titleCompose(titleId int32) error {

	for _, v := range p.titles.Titles {
		if v.CfgId == titleId {
			return errutil.ErrComposeError
		}
	}
	cfg := global.GetTitleCfg(titleId)
	if cfg == nil {
		return errutil.ErrConfigNotFound
	}

	items := []global.IdNum{}
	for _, v := range cfg.Activity {
		for _, cost := range v.Values {
			if !p.checkTakeSource(cost.Key, int64(cost.Value)) {
				return errutil.ErrItemNotEnough
			}
			items = append(items, global.IdNum{ID: cost.Key, Num: cost.Value})
		}
	}
	for _, v := range items {
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelTitle)
	}

	p.titles.Titles = append(p.titles.Titles, &pb.TitleInfo{CfgId: titleId, ActivationTime: int32(time.Now().Unix())})

	p.checkTitleBonds()
	p.refreshStats()
	return nil
}

//titleStart

func (p *Player) titleStart(titleId int32, luckItem int32) (int32, error) {
	var title *pb.TitleInfo
	for _, v := range p.titles.Titles {
		if v.CfgId == titleId {
			title = v
			break
		}
	}
	if title == nil {
		return 0, errutil.StarsError
	}

	if title.Star >= int32(global.GlobalCfg.War.Titlestar.Maxnum) {
		return title.Star, errutil.ErrMaxLevel
	}

	cfg := global.GetTitleCfg(titleId)
	if cfg == nil {
		return 0, errutil.ErrConfigNotFound
	}

	items := []global.IdNum{}
	for _, cost_items := range cfg.StarConsume {
		for _, v := range cost_items.Values {
			if !p.checkTakeSource(v.Key, int64(v.Value)) {
				return 0, errutil.ErrItemNotEnough
			}
			items = append(items, global.IdNum{ID: v.Key, Num: v.Value})
		}
	}

	if luckItem > 0 {
		if !p.checkTakeSource(luckItem, 1) {
			return 0, errutil.ErrItemNotEnough
		}
	}

	// 扣除材料和幸运符
	for _, v := range items {
		p.takeSource(v.ID, int64(v.Num), constant.ItemDelTitle)
	}

	p.takeSource(luckItem, 1, constant.ItemDelTitle)

	rate := 0
	for _, v := range cfg.StarRate {
		if v.Key == title.Star+1 {
			rate = int(v.Value)
			break
		}
	}

	for _, v := range global.GlobalCfg.War.Luckticket.Item {
		if v.ID == luckItem {
			rate += int(v.Num)
			break
		}
	}

	oldStar := title.Star
	if global.RandIn(int(BaseRate)) < rate {
		title.Star++
	}
	p.checkTitleBonds()

	if title.Star > oldStar {
		p.refreshStats()
	}

	return title.Star, nil
}

// titleEquip
func (p *Player) titleEquip(titleId int32, slot int32) error {
	var title *pb.TitleInfo
	for _, v := range p.titles.Titles {
		if v.CfgId == titleId {
			title = v
			break
		}
	}
	if title == nil && titleId != 0 {
		return errutil.ErrItemNotEnough
	}

	cfg := global.GetTitleCfg(titleId)
	if cfg != nil {
		if cfg.Slot-1 != slot {
			return errutil.ErrEquipSlot
		}
	}

	if slot == 0 {
		p.equipTitleId = titleId
		p.syncUserBase()
	} else {
		if partnerslot := p.getPartnerSlotInfo(slot); partnerslot != nil {
			partnerslot.TitleId = titleId
			if par := p.getPartnerInfoBySlotId(slot); par != nil {
				par.equipTitleId = titleId
			}
			p.SyncUserPartnerSlot()
		}
	}
	p.refreshStats()

	return nil
}

// 伙伴羁绊
func (p *Player) checkTitleBonds() {
	bond := map[int32]int32{}

	for _, v := range global.TitleActivationCfg {

		partIds := []int32{}
		for _, t := range global.TitleCfg {
			if t.Slot == v.Slot {
				partIds = append(partIds, t.UID)
			}
		}
		unlock := true
		for _, cp := range partIds {
			if title := p.getTitle(cp); title != nil {
				if title.Star < v.StarLimit {
					unlock = false
					break
				}
			} else {
				unlock = false
				break
			}
		}

		if unlock {
			if b, ok := bond[v.Slot]; ok {
				if b < v.ID {
					bond[v.Slot] = v.ID
				}
			} else {
				bond[v.Slot] = v.ID
			}
		}
	}

	for k, v := range bond {
		if k == 1 {
			p.activateTitleBonds = []int32{}
			p.activateTitleBonds = append(p.activateTitleBonds, v)
		} else {
			parslot := p.getPartnerSlotInfo(k - 1)
			if parslot == nil {
				continue
			}

			par := p.getPartnerInfoBySlotId(k - 1)
			if par == nil {
				continue
			}
			par.activateTitleBonds = []int32{}
			par.activateTitleBonds = append(par.activateTitleBonds, v)
		}

	}
}

// getTitle
func (p *Player) getTitle(cfgId int32) *pb.TitleInfo {
	for _, v := range p.titles.Titles {
		if v.CfgId == cfgId {
			return v
		}
	}
	return nil
}

// 体力回复
func (p *Player) energyRecovery() {
	p.nextEnergyRecoveryTime = 0
	if p.energy < p.maxEnergy() {
		now := time.Now()
		offlineSec := int64(now.Unix() - p.lastLoginAt.Unix())
		// 离线恢复总体力
		recoverNum := offlineSec / int64(global.GlobalCfg.BaseValue.Energy.Timeinterval*60)
		// 不足恢复一次体力的时间
		leftTime := offlineSec % int64(global.GlobalCfg.BaseValue.Energy.Timeinterval*60)
		p.energy += recoverNum
		if p.energy > p.maxEnergy() {
			p.energy = p.maxEnergy()
			p.energyRecoveryTime = now
		} else {
			p.nextEnergyRecoveryTime = global.GlobalCfg.BaseValue.Energy.Timeinterval*60 - int(leftTime)
			p.energyRecoveryTime = now.Add(-time.Duration(leftTime) * time.Second)
		}
	}
}

func (p *Player) setScore(score int64, pvpType int) {

	p.lock.Lock()
	defer p.lock.Unlock()
	switch pvpType {
	case constant.RankTypeSoul:
		p.pvpScore = int32(score)
	case constant.RankTypeTurnBased:
		p.gvgScore = int32(score)
	}
}

// 最大体力值
func (p *Player) maxEnergy() int64 {
	maxEnergy := int64(global.GlobalCfg.BaseValue.Energy.Fullvalue)
	for _, v := range p.privilege {
		if privilegeCfg := global.GetPrivilegeCfg(v.Pid); privilegeCfg != nil {
			maxEnergy += int64(privilegeCfg.StaminaUp)
		}
	}
	return maxEnergy
}

// 添加记录
func (p *Player) addKillMonControl(monCfgId, num int32) {
	for _, v := range global.KillMonControlCfg {
		for _, id := range v.MonsterIDs {
			if id == monCfgId {
				p.InstanceData.MonKillRecord.MonKillRecords[v.Uid] += num
				break
			}
		}
	}
}

// 查询记录
func (p *Player) getKillMonControl(monCfgId int32) bool {

	for _, v := range global.KillMonControlCfg {
		for _, id := range v.MonsterIDs {
			if id == monCfgId {
				if r, ok := p.InstanceData.MonKillRecord.MonKillRecords[v.Uid]; ok {
					if r >= v.Maxnum {
						return false
					}
				}
			}
		}
	}
	return true
}

func (p *Player) isModuleOpen(moduleID constant.ModuleID) bool {
	modelCfg, ok := global.SwitchCfg[int32(moduleID)]
	if !ok {
		logger.Warnf("isModuleOpen not found moduleID: %d, uid: %d", moduleID, p.uid)
		return false
	}
	if p.Level() < modelCfg.OpenLevel {
		return false
	}
	if p.mainBattleId < modelCfg.MinBattel {
		return false
	}
	if modelCfg.ServerDay > 0 {
		if modelCfg.OpenDay <= 1 {
			return true
		}
		serverInfo, err := global.ServerListMgrInstance.Get(int(p.serverID))
		if err != nil {
			logger.Errorf("isModuleOpen Get serverInfo err: %v, uid: %d", err, p.uid)
			return false
		}
		openTime := time.Unix(serverInfo.OpenBy/1000, 0)
		localOpenZero := time.Date(openTime.Year(), openTime.Month(), openTime.Day(), 0, 0, 0, 0, time.Local)
		diffDays := int32(time.Since(localOpenZero).Hours()/24) + 1
		return diffDays >= modelCfg.OpenDay
	} else {
		return p.createDay >= modelCfg.OpenDay
	}
}

// distributeRewards 统一处理奖励发放逻辑
// 检查包裹空间，空间不足时发送邮件，否则直接发放奖励
func (d *Player) distributeRewards(p *Player, rewards []global.IdNum, addSourceType string, addSourceTypeUnNotify bool) {
	bSendMail := false

	// 检查所有奖励是否能放入包裹
	for _, reward := range rewards {
		if !p.checkAddSource(reward.ID, int64(reward.Num)) {
			bSendMail = true
			break
		}
	}

	if !bSendMail {
		// 包裹空间足够，直接批量发放奖励
		sourceInfo := make([]global.IdNum, 0, len(rewards))
		for _, reward := range rewards {
			sourceInfo = append(sourceInfo, global.IdNum{ID: reward.ID, Num: reward.Num})
		}
		if addSourceTypeUnNotify {
			p.addSourceBatch(sourceInfo, addSourceType, WithAddSourceUnNotify(AddSourceTypeUnNotify))
		} else {
			p.addSourceBatch(sourceInfo, addSourceType)
		}
	} else {
		// 包裹空间不足，通过邮件发送
		attachment := make([]protocol.Attachment, 0, len(rewards))
		for _, reward := range rewards {
			attachment = append(attachment, protocol.Attachment{
				UserPropsId: int64(reward.ID),
				Num:         int64(reward.Num),
			})
		}
		p.SendMail(addSourceType, attachment)
	}
}

// BuildUserReportMap 构建上报给ShuShu的用户数据
func (p *Player) BuildUserReportMap(m tracking.M) map[string]interface{} {
	var (
		petCfgID   int32
		petName    string
		petLevel   int32
		petQuality int32
		petZun     bool
	)
	bPet := p.PetMgr.GetBattlePet()
	if bPet.PetCfgId != 0 {
		petCfgID = bPet.PetCfgId
		petLevel = bPet.Level
		petQuality = bPet.Quality
		if petCfg := global.GetPetCfg(bPet.PetCfgId); petCfg != nil {
			petName = petCfg.Name
			petZun = petCfg.Mark == 1
		}
	}
	//if len(p.pet) > 0 && p.pet[0].PetCfgId != 0 {
	//	petCfgID = p.pet[0].PetCfgId
	//	petLevel = p.pet[0].Level
	//	petQuality = p.pet[0].Quality
	//	if petCfg := global.GetPetCfg(p.pet[0].PetCfgId); petCfg != nil {
	//		petName = petCfg.Name
	//		petZun = petCfg.Mark == 1
	//	}
	//}
	base := map[string]interface{}{
		"channel":       "",                                // 渠道，暂为空
		"server_id":     strconv.Itoa(int(p.serverID)),     // 服务器ID，格式为字符串
		"nick_name":     p.name,                            // 昵称
		"tmp_power":     p.totalCombat,                     // 当前战力值
		"tmp_level":     p.Level(),                         // 当前等级
		"tmp_box_level": p.box.NowLevel,                    // 宝箱level
		"tmp_main_id":   p.InstanceData.BattleMap.BattleId, // 当前战役关卡进度
		"tmp_task_id":   p.taskProgress,                    // 当前任务进度
		"tmp_role_job":  constant.JobName[p.job],           // 主角职业
		"tmp_account":   p.accountName,                     // 账号id
		"tmp_device":    "",                                // 设备id，需补充
		"tmp_role_id":   p.uid,                             // 角色id
		"tmp_pet_id":    petCfgID,                          // 当前上阵宠物id
		"tmp_pet_name":  petName,                           // 当前上阵宠物名
		"tmp_pet_level": petLevel,                          // 当前上阵宠物level
		"tmp_pet_rank":  petQuality,                        // 当前上阵宠物突破等级
		"tmp_pet_zun":   petZun,                            // 当前上阵宠物是否为尊
	}

	for k, v := range m {
		base[k] = v
	}
	return base
}

func (p *Player) AddAdCount() {
	atomic.AddInt32(&p.adTotalToday, 1)
}

func (p *Player) GetAdCount() int32 {
	return atomic.LoadInt32(&p.adTotalToday)
}

// 离线恢复
func (p *Player) offlineRecovery() {

	now := time.Now()
	offlineSec := int64(now.Unix() - p.lastLoginAt.Unix())

	//离线体力恢复
	p.energyRecovery()

	//离线宝箱恢复
	if !p.bCreate && p.getItemCountInBag(constant.ItemTypeTreasureNum) < int64(global.GlobalCfg.BaseValue.Offlinebox.Maxnum) {

		p.box.OfflineTime += offlineSec

		// 离线恢复总宝箱数
		recoverNum := p.box.OfflineTime / int64(global.GlobalCfg.BaseValue.Offlinebox.Interval*60)
		// 不足恢复一个的时间
		leftTime := p.box.OfflineTime % int64(global.GlobalCfg.BaseValue.Offlinebox.Interval*60)
		p.box.OfflineTime = leftTime

		if recoverNum > int64(global.GlobalCfg.BaseValue.Offlinebox.Maxnum)-p.getItemCountInBag(constant.ItemTypeTreasureNum) {
			recoverNum = int64(global.GlobalCfg.BaseValue.Offlinebox.Maxnum) - p.getItemCountInBag(constant.ItemTypeTreasureNum)
		}
		if recoverNum > 0 {
			p.addSource(constant.ItemTypeTreasureNum, recoverNum, constant.ItemAddOfflineBox)
		}
		if p.getItemCountInBag(constant.ItemTypeTreasureNum) >= int64(global.GlobalCfg.BaseValue.Offlinebox.Maxnum) {
			p.box.OfflineTime = 0
		}
		p.logger.Debugf(" %s 上线恢复宝箱数量 %d  当前持有宝箱数量 %d 剩余累计离线恢复宝箱时间 %d", p.name, recoverNum, p.getItemCountInBag(constant.ItemTypeTreasureNum), p.box.OfflineTime)
	}
}

func (p *Player) toDbModuleNew() []byte {
	infos := &pb.ModuleInfos{Infos: []*pb.ModuleInfo{}}
	p.ModuleNew.Range(func(key, value interface{}) bool {
		info, ok := value.(*pb.ModuleInfo)
		if ok {
			infos.Infos = append(infos.Infos, info)
		}
		return true
	})
	if data, err := protobuf.NewSerializer().Marshal(infos); err == nil {
		return data
	}
	return nil
}

func (p *Player) fromDbModuleNew(moduleNew []byte) {
	mis := &pb.ModuleInfos{}
	if err := protobuf.NewSerializer().Unmarshal(moduleNew, mis); err == nil {
		for _, v := range mis.Infos {
			p.ModuleNew.Store(v.GetCfgId(), v)
		}
	} else {
		p.logger.Errorf("新模块激活反序列化失败 err: %s", err.Error())
	}
}

func (p *Player) PushModuleNew() {
	infos := []*pb.ModuleInfo{}
	p.ModuleNew.Range(func(key, value interface{}) bool {
		info, ok := value.(*pb.ModuleInfo)
		if ok {
			infos = append(infos, info)
		}
		return true
	})
	p.pushMsg(protocol.RouteModuleNewPush, &pb.ModuleNewPush{
		Infos: infos,
	})
}

func (p *Player) dispatchTask(taskEventData *event.TaskEventData) {
	p.eventMgr.Publish(event.EventTypeTaskProgress, taskEventData)
}

// GetMinStarByQuality 持有的 X 品质的伙伴全部达到的 {Y} 最大星数
func (p *Player) GetMinStarByQuality(quality constant.Quality) int32 {
	minStar := int32(math.MaxInt32)
	count := 0
	for _, partner := range p.partners {
		if partner.cfg != nil && constant.Quality(partner.cfg.Quality) == quality {
			if partner.star < minStar {
				minStar = partner.star
			}
		}
		count++
	}

	if count == 0 {
		return 0
	}

	return minStar
}

// GetMainRoleAllSlotMinEnhanceLevel 主角全身槽位强化至 X 级
func (p *Player) GetMainRoleAllSlotMinEnhanceLevel() int32 {
	minEnhance := int32(math.MaxInt32)
	count := 0

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		if equip.Enhance < minEnhance {
			minEnhance = equip.Enhance
		}
		count++
	}

	if count != int(pb.EquipSlot_MaxEquipSlot-1) {
		return 0
	}

	return minEnhance
}

// 已方全部槽位强化
func (p *Player) GetAllSlotMinEnhanceLevel() int32 {
	minEnhance := int32(math.MaxInt32)

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		if equip.Enhance < minEnhance {
			minEnhance = equip.Enhance
		}
	}

	for _, slotInfo := range p.partnerSlots {
		if slotInfo == nil {
			continue
		}
		if slotInfo.PartnerUid == 0 {
			continue
		}
		for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
			equip := slotInfo.Equips[i]
			if equip.ItemId == 0 {
				continue
			}
			if equip.Enhance < minEnhance {
				minEnhance = equip.Enhance
			}
		}
	}

	return minEnhance
}

// GetMainRoleAllSlotMinStar 主角全身槽位升星至 X 星
func (p *Player) GetMainRoleAllSlotMinStar() int32 {
	minStar := int32(math.MaxInt32)
	count := 0

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		if equip.Stars < minStar {
			minStar = equip.Stars
		}
		count++
	}

	if count != int(pb.EquipSlot_MaxEquipSlot-1) {
		return 0
	}

	return minStar
}

// GetAllSlotMinStar 已方全部槽位升星
func (p *Player) GetAllSlotMinStar() int32 {
	minStar := int32(math.MaxInt32)

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		if equip.Stars < minStar {
			minStar = equip.Stars
		}
	}

	for _, slotInfo := range p.partnerSlots {
		if slotInfo == nil {
			continue
		}
		if slotInfo.PartnerUid == 0 {
			continue
		}
		for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
			equip := slotInfo.Equips[i]
			if equip.ItemId == 0 {
				continue
			}
			if equip.Stars < minStar {
				minStar = equip.Stars
			}
		}
	}

	return minStar
}

// GetMainRoleAllSlotMinSoulStoneLevel 主角全身镶嵌 X 级以上魂石
func (p *Player) GetMainRoleAllSlotMinSoulStoneLevel() int32 {
	minSoulStoneLevel := int32(math.MaxInt32)
	equipCount := 0
	gemStoreCount := 0

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		equipCount++
		for _, soulStoneId := range equip.Gemstone {
			if soulStoneId == 0 {
				continue
			}
			itemCfg, err := global.GetItemCfg(int32(soulStoneId))
			if err != nil {
				logger.Warnf("GetMainRoleAllSlotMinSoulStoneLevel 获取魂石配置失败 soulStoneId: %d, uid: %d", soulStoneId, p.uid)
				continue
			}
			if itemCfg.Itemlevel < minSoulStoneLevel {
				minSoulStoneLevel = itemCfg.Itemlevel
			}
			gemStoreCount++
		}
	}

	// 根据等级获取开启的魂石槽位数量
	openGemStoreCount := 0
	for _, v := range global.GlobalCfg.Diamond.Openlevel.Stage {
		if p.Level() >= v.Level {
			openGemStoreCount = int(v.Num)
		}
	}

	if equipCount != int(pb.EquipSlot_MaxEquipSlot-1) {
		return 0
	}

	if gemStoreCount != equipCount*openGemStoreCount {
		return 0
	}

	return minSoulStoneLevel
}

// 已方全部槽位镶嵌 X 级以上魂石
func (p *Player) GetAllSlotMinSoulStoneLevel() int32 {
	minSoulStoneLevel := int32(math.MaxInt32)
	partner := 0
	equipCount := 0
	gemStoreCount := 0

	for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
		equip := p.equips[i]
		if equip.ItemId == 0 {
			continue
		}
		equipCount++
		for _, soulStoneId := range equip.Gemstone {
			if soulStoneId == 0 {
				continue
			}
			gemStoreCount++
			itemCfg, err := global.GetItemCfg(int32(soulStoneId))
			if err != nil {
				logger.Warnf("GetAllSlotMinSoulStoneLevel 获取魂石配置失败 soulStoneId: %d, uid: %d", soulStoneId, p.uid)
				continue
			}
			if itemCfg.Itemlevel < minSoulStoneLevel {
				minSoulStoneLevel = itemCfg.Itemlevel
			}
		}
	}

	for _, slotInfo := range p.partnerSlots {
		if slotInfo == nil {
			continue
		}
		if slotInfo.PartnerUid == 0 {
			continue
		}
		partner++
		for i := pb.EquipSlot_Weapon; i < pb.EquipSlot_MaxEquipSlot; i++ {
			equip := slotInfo.Equips[i]
			if equip.ItemId == 0 {
				continue
			}
			equipCount++
			for _, soulStoneId := range equip.Gemstone {
				if soulStoneId == 0 {
					continue
				}
				gemStoreCount++
				itemCfg, err := global.GetItemCfg(int32(soulStoneId))
				if err != nil {
					logger.Warnf("GetAllSlotMinSoulStoneLevel 获取魂石配置失败 soulStoneId: %d, uid: %d", soulStoneId, p.uid)
					continue
				}
				if itemCfg.Itemlevel < minSoulStoneLevel {
					minSoulStoneLevel = itemCfg.Itemlevel
				}
			}
		}
	}

	// 根据等级获取开启的魂石槽位数量
	openGemStoreCount := 0
	for _, v := range global.GlobalCfg.Diamond.Openlevel.Stage {
		if p.Level() >= v.Level {
			openGemStoreCount = int(v.Num)
		}
	}

	if equipCount != (int(pb.EquipSlot_MaxEquipSlot-1) * (partner + 1)) {
		return 0
	}

	if gemStoreCount != equipCount*openGemStoreCount {
		return 0
	}

	return minSoulStoneLevel
}

// 计算粉色词条数量
func (p *Player) calcPinkEntryNum() int64 {
	var count int64 = 0

	// 第一层循环：遍历配置，按属性类型分组收集所有粉色词条的数值范围
	// 使用 map 存储每个属性类型对应的所有数值范围，提高查找效率
	type valueRange struct {
		min int64
		max int64
	}
	pinkRangesMap := make(map[int32][]valueRange)
	for _, cfg := range global.AttributeSegmentCfgV2 {
		// 判断品质颜色是否大于等于6（粉色词条）
		if cfg.QualityColor >= 5 {
			attrType := cfg.AttributeType
			if pinkRangesMap[attrType] == nil {
				pinkRangesMap[attrType] = make([]valueRange, 0)
			}
			pinkRangesMap[attrType] = append(pinkRangesMap[attrType], valueRange{
				min: cfg.ValueMin,
				max: cfg.ValueMax,
			})
		}
	}

	// 第二层循环：遍历保存的属性数据，判断属性类型和数值是否在范围内
	p.refine.attrs.Range(func(key int32, value *pb.Stats) bool {
		attrTypeInt := int32(value.Attr)
		// 查找该属性类型对应的所有数值范围
		if ranges, exists := pinkRangesMap[attrTypeInt]; exists {
			numInt64 := int64(value.Num)
			// 判断数值是否在任何一个配置范围内
			for _, r := range ranges {
				if numInt64 >= r.min && numInt64 <= r.max {
					count++
					break // 找到一个匹配的范围即可
				}
			}
		}
		return true
	})

	return count
}

// 计算红色伙伴总数（因为升星的缘故，1星算是2个伙伴，以此类推）
func (p *Player) calcRedPartnerStarNum() int64 {
	var totalNum int64 = 0
	// 遍历所有拥有的伙伴列表，而不是只计算槽位上的
	for _, partner := range p.partners {
		if partner != nil && partner.cfg != nil {
			// 判断是否为红色品质（品质值为7）
			if partner.cfg.Quality >= 7 {
				// 因为升星的缘故，总数 = 星级 + 1
				// 0星 = 1个，1星 = 2个，2星 = 3个，以此类推
				totalNum += int64(partner.star) + 1
			}
		}
	}
	//加入伙伴排行
	return totalNum
}

// 计算装备总升星数
func (p *Player) calcEquipmentStarNum() int64 {
	var totalStar int64 = 0
	// 遍历主角装备
	for _, equip := range p.equips {
		if equip.ItemId > 0 {
			totalStar += int64(equip.Stars)
		}
	}
	// 遍历伙伴装备
	for _, partnerSlot := range p.partnerSlots {
		if partnerSlot != nil {
			for _, equip := range partnerSlot.Equips {
				if equip.ItemId > 0 {
					totalStar += int64(equip.Stars)
				}
			}
		}
	}
	return totalStar
}

// 计算宠物合成数量
func (p *Player) calcPetSynthesisNum() int64 {
	// TODO: 根据实际需求计算宠物合成数量
	// 可能需要从记录中获取，或者计算已合成的宠物数量
	totalNum := p.getUpRecord(constant.PetSynthesisNum)
	return totalNum
}

// 计算魂石总等级
func (p *Player) calcSoulStoneLevelNum() int64 {
	var totalLevel int64 = 0
	// TODO: 根据实际需求计算魂石总等级
	// Gemstone 字段存储的是宝石槽位索引（int32），需要从背包或其他地方获取宝石信息
	// 遍历主角装备的宝石槽位
	for _, equip := range p.equips {
		if equip.ItemId > 0 {
			for _, gemstoneSlot := range equip.Gemstone {
				if gemstoneSlot > 0 {
					gemstoneCfg, err := global.GetItemCfg(int32(gemstoneSlot))
					if err != nil {
						continue
					}
					totalLevel += int64(gemstoneCfg.Itemlevel)
				}
			}
		}
	}
	// 遍历伙伴装备的宝石槽位
	for _, partnerSlot := range p.partnerSlots {
		if partnerSlot != nil {
			for _, equip := range partnerSlot.Equips {
				if equip.ItemId > 0 {
					for _, gemstoneSlot := range equip.Gemstone {
						if gemstoneSlot > 0 {
							gemstoneCfg, err := global.GetItemCfg(int32(gemstoneSlot))
							if err != nil {
								continue
							}
							totalLevel += int64(gemstoneCfg.Itemlevel)
						}
					}
				}
			}
		}
	}
	return totalLevel
}

// 计算五行赋能总等级
func (p *Player) calcFiveElementLevelNum() int64 {
	var totalLevel int64 = 0
	if p.ElementInfo != nil {
		// 遍历五行元素并累加等级
		for _, element := range p.ElementInfo.element {
			totalLevel += int64(element.Level)
		}
	}
	return totalLevel
}
