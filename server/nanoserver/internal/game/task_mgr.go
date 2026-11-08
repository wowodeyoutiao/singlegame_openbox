package game

import (
	"math"
	"strconv"
	"strings"
	"time"

	"nanoserver/db"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
)

type (
	TargetInfo struct {
		TargetType pb.TargetType // 分类 task.Type
		TargetId   int           // 子分类 task.Typeson
		TargetNum  int           // 达成条件 task.Limit
		Amount     int           // 当前进度
		Completed  bool          // 是否完成
	}
	TaskInfo struct {
		ID        int           // task.UUID
		Class     pb.TaskClass  // task.Class
		Status    pb.TaskStatus // 任务状态
		Target    *TargetInfo
		ResetTime time.Time
	}
	TaskMgr struct {
		*Player
		tasks map[int32]*TaskInfo // 任务
	}
)

func NewTask(p *Player) *TaskMgr {
	return &TaskMgr{tasks: make(map[int32]*TaskInfo), Player: p}
}

// dispatch 任务进度派发器
func (tm *TaskMgr) dispatch(targetType pb.TargetType, i int32, v int32) {
	if tm.isRobot() {
		return
	}

	if targetType <= pb.TargetType_TASK_TARGET_TYPE_INVALID || targetType > pb.TargetType_TASK_TARGET_TYPE_MAX {
		return
	}

	for _, t := range tm.tasks {

		if t.Status != pb.TaskStatus_IN_PROGRESS && t.Class != pb.TaskClass_DAILY_ACTIVITY && t.Class != pb.TaskClass_WEEKLY_ACTIVITY {
			continue
		}
		if t.Target.Completed && t.Class != pb.TaskClass_DAILY_ACTIVITY && t.Class != pb.TaskClass_WEEKLY_ACTIVITY {
			continue
		}
		if t.Target.TargetType != targetType {
			continue
		}

		switch targetType {
		// ====== 直接赋值类（取最大值）传过来可能小于现有的 ======
		case pb.TargetType_ROLE_POWER_REQUIREMENT,
			pb.TargetType_ROLE_LEVEL_REQUIREMENT,
			pb.TargetType_EQUIPMENT_STRENGTHENING_LEVEL_REQUIREMENT,
			pb.TargetType_EQUIPMENT_UPGRADING_STAR_LEVEL_REQUIREMENT,
			pb.TargetType_TOTAL_ATTEMPTS_IN_ENDLESS_TRIALS,
			pb.TargetType_VIP_LEVEL,
			pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL,
			pb.TargetType_PARTNER_ON_STAGE_TIMES,
			pb.TargetType_UNLOCK_PARTNER_SLOT,
			pb.TargetType_CHEST_LEVEL_UPGRADED,
			pb.TargetType_ALL_GREEN_PARTNERS_REACH_STAR,
			pb.TargetType_ALL_BLUE_PARTNERS_REACH_STAR,
			pb.TargetType_ALL_PURPLE_PARTNERS_REACH_STAR,
			pb.TargetType_ALL_ORANGE_PARTNERS_REACH_STAR,
			pb.TargetType_ALL_RED_PARTNERS_REACH_STAR,
			pb.TargetType_ALL_YELLOW_PARTNERS_REACH_STAR,
			pb.TargetType_MAIN_ALL_EQUIP_SLOTS_STRENGTHENED_TO_LEVEL,
			pb.TargetType_MAIN_ALL_EQUIP_SLOTS_STAR_TO_LEVEL,
			pb.TargetType_MAIN_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST,
			pb.TargetType_ALL_EQUIP_SLOTS_STRENGTHENED_TO_LEVEL,
			pb.TargetType_ALL_EQUIP_SLOTS_STAR_TO_LEVEL,
			pb.TargetType_ALL_EQUIP_SLOTS_GEMSTONE_LEVEL_AT_LEAST,
			pb.TargetType_FIVE_ELEMENTS_INSTANCE_CLEARED_FLOOR,
			pb.TargetType_EQUIP_ITEMS_WORN_TOTAL:
			if t.Target.Amount < int(v) {
				t.Target.Amount = int(v)
				tm.update(t)
			}
			// ====== 统计类（增量累加） ======
		case pb.TargetType_KILL_MONSTER_REQUIREMENT,
			pb.TargetType_PVE_INSTANCE_ENTRY_REQUIREMENT,
			pb.TargetType_TOTAL_PHYSICAL_STRENGTH_CONSUMED_REQUIREMENT,
			pb.TargetType_TOTAL_GOLD_COINS_CONSUMED_REQUIREMENT,
			pb.TargetType_DAILY_ACTIVITY_LEVEL,
			pb.TargetType_WEEKLY_ACTIVITY_LEVEL,
			pb.TargetType_TOTAL_BOOST_PURCHASES,
			pb.TargetType_TOTAL_CHALLENGES_IN_ARENA,
			pb.TargetType_TOTAL_KILLS_OF_NORMAL_MONSTERS,
			pb.TargetType_TOTAL_KILLS_OF_BOSS_MONSTERS,
			pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_GOLD,
			pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_COUPONS,
			pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_SOULS,
			pb.TargetType_TOTAL_SHOP_PURCHASES_WITH_TRIAL_TOKENS,
			pb.TargetType_TOTAL_ACQUIRED_SOULS,
			pb.TargetType_TOTAL_ACQUIRED_TRIAL_TOKENS,
			pb.TargetType_TOTAL_BOX_OPENED,
			pb.TargetType_TOTAL_MATERIAL_INSTANCE_CHALLENGE,
			pb.TargetType_TOTAL_MULTI_VS_MULTI_ARENA_CHALLENGE,
			pb.TargetType_TOTAL_BOSS_CHALLENGE,
			pb.TargetType_TOTAL_RELIC_CHALLENGE,
			pb.TargetType_TOTAL_PET_EGGS_OPENED,
			pb.TargetType_TOTAL_EXPERIENCE_FB_CLEARED,
			pb.TargetType_TOTAL_GOLD_FB_CLEARED,
			pb.TargetType_ANY_EQUIP_STAR_UP_TIMES,
			pb.TargetType_SKILL_UPGRADING_TIMES,
			pb.TargetType_SEVENDAY_TASK_GIFT_PURCHASED,
			pb.TargetType_SEVENDAY_TASK_ALL_COMPLETED,
			pb.TargetType_PARTNER_RECEIVE,
			pb.TargetType_GUARD_ISLAND_INSTANCE,
			pb.TargetType_PET_COMPOSE_TIMES,
			pb.TargetType_TURN_BASED_PK_TIMES,
			pb.TargetType_SOUL_STONE_INSTANCE,
			pb.TargetType_SKILL_ADVANCED_TIMES,
			pb.TargetType_JOIN_GUILD,
			pb.TargetType_NORMAL_REFINE_TIMES,
			pb.TargetType_STAR_UP_STONE_INSTANCE,
			pb.TargetType_BOOK_PAGE_INSTANCE,
			pb.TargetType_SPECIAL_REFINE_TIMES,
			pb.TargetType_ENHANCE_STONE_INSTANCE,
			pb.TargetType_FIVE_ELEMENTS_INSTANCE,
			pb.TargetType_GEMSTONE_INLAY_TIMES,
			pb.TargetType_RELIC_TOTAL_STAR,
			pb.TargetType_WORLD_BOSS_PARTICIPATION,
			pb.TargetType_TOTAL_STRENGTHEN_TIMES,
			pb.TargetType_DRAW_GREEN_PET_EGGS,
			pb.TargetType_DRAW_PURPLE_PET_EGGS,
			pb.TargetType_DRAW_ORANGE_PET_EGGS,
			pb.TargetType_COMPOSE_BLUE_COMMON_PETS,
			pb.TargetType_COMPOSE_BLUE_ELITE_PETS,
			pb.TargetType_COMPOSE_PURPLE_COMMON_PETS,
			pb.TargetType_COMPOSE_PURPLE_ELITE_PETS,
			pb.TargetType_COMPOSE_ORANGE_COMMON_PETS,
			pb.TargetType_COMPOSE_ORANGE_ELITE_PETS,
			pb.TargetType_COMPOSE_RED_COMMON_PETS,
			pb.TargetType_COMPOSE_RED_ELITE_PETS,
			pb.TargetType_COMPOSE_YELLOW_ELITE_PETS,
			pb.TargetType_BATTLEFIELD_OPEN_BLUE_LUCKY_BAG_TIMES,
			pb.TargetType_BATTLEFIELD_OPEN_PURPLE_LUCKY_BAG_TIMES,
			pb.TargetType_BATTLEFIELD_OPEN_ORANGE_LUCKY_BAG_TIMES,
			pb.TargetType_BATTLEFIELD_TOTAL_KILLS,
			pb.TargetType_BATTLEFIELD_TOTAL_CHESTS_OPENED,
			pb.TargetType_TOTAL_REFINE_TIMES,
			pb.TargetType_MAIN_ALL_REFINE_SLOTS_OPENED_TOTAL,
			pb.TargetType_ALL_REFINE_ENTRIES_OPENED_TOTAL,
			pb.TargetType_LUCKY_BAG_OBTAINED_TOTAL,
			pb.TargetType_LUCKY_BAG_OPENED_TOTAL,
			pb.TargetType_PUSHING_WALL_PLAY_TIMES,
			pb.TargetType_ARTIFACT_UPGRADE_TIMES,
			pb.TargetType_TICKET_EXCHANGE_TIMES,
			pb.TargetType_FIVE_ELEMENT_EMPOWER_TIMES:
			t.Target.Amount += int(v)
			tm.update(t)
			// ====== 指定子ID 的当前值赋值（例如技能升级直接写等级） ======
		case pb.TargetType_SKILL_UPGRADING_REQUIREMENT:
			if t.Target.TargetId == int(i) {
				t.Target.Amount = int(v)
			}
			// ====== 指定子ID 的增量累加（仅当 i 匹配时才累加） ======
		case pb.TargetType_TOTAL_ITEM_CONSUMPTION,
			pb.TargetType_BATTLE_CHAPTERS_CLEARED,
			pb.TargetType_BOSS_LEVELS_CLEARED,
			pb.TargetType_TRIAL_TOWER_FLOORS_CLEARED,
			pb.TargetType_RELIC_CHAPTERS_CLEARED,
			pb.TargetType_UPGRADE_ANY_PARTNER_LEVEL,
			pb.TargetType_ANY_EQUIP_STRENGTHENED_TO_LEVEL,
			pb.TargetType_MAIN_BATTLE,
			pb.TargetType_EQUIP_STAR_TO_X,
			pb.TargetType_RELIC_CHAPTERS_ALL_3STARS, // ok
			pb.TargetType_ENDLESS_TRIALS_STAGE:
			if t.Target.TargetId == int(i) {
				t.Target.Amount += int(v)
			}
		default:

		}

		tm.checkProgress(t)
	}
}

func (tm *TaskMgr) resetNewday() {
	logger.Debugf("新的一天刷新日日长任务 uid %d  name %s", tm.uid, tm.name)
	delTask := []int32{}
	for _, v := range tm.tasks {
		// 每日福利
		if v.Class == pb.TaskClass_ROUTINE {
			v.Target.Amount = 0
			v.Target.Completed = false
			v.Status = pb.TaskStatus_IN_PROGRESS

			tm.update(v)

			continue
		}

		// 日常任务
		if v.Class == pb.TaskClass_DAILY || v.Class == pb.TaskClass_DAILY_ACTIVITY {

			if vCfg, err := global.GetTaskCfg(int32(v.ID)); err == nil {
				if len(vCfg.Level) == 2 {
					if tm.level >= vCfg.Level[0] && tm.level <= vCfg.Level[1] {
						v.Target.Amount = 0
						v.Target.Completed = false
						v.Status = pb.TaskStatus_IN_PROGRESS

						tm.update(v)

						logger.Debugf("重置日常任务 uid %d  name %s  taskid %d", tm.uid, tm.name, v.ID)

					} else {
						delTask = append(delTask, int32(v.ID))
						logger.Debugf("删除日常任务，不符合等级段 uid %d  name %s  taskid %d", tm.uid, tm.name, v.ID)

					}
				}
			}
			continue
		}

		if v.Class == pb.TaskClass_SEVENDAY || v.Class == pb.TaskClass_SEVENDAY_ACTIVITY {
			//查询这个任务对应活动表的时间，如果时间已过，直接删掉

			vCfg, err := global.GetTaskCfg(int32(v.ID))
			if err == nil {
				actCfg, err := global.GetActivityCfg(int32(vCfg.ActivID))
				if err != nil {
					continue
				}
				if tm.createDay > actCfg.ResetTime {
					//如果任务完成了，直接给他发奖励
					//但是也要删掉他的任务
					if v.Status == pb.TaskStatus_COMPLETED {
						result3 := make([]protocol.Attachment, len(vCfg.Reward))
						// MailId := strconv.FormatInt(snowflake.GenID(), 10)
						for i, Attachment := range vCfg.Reward {
							result3[i] = protocol.Attachment{
								MailId:      "",
								UserPropsId: int64(Attachment.Key),
								Num:         int64(Attachment.Value),
							}
						}
						taskName := vCfg.TaskDes
						taskName = strings.Replace(taskName, "{0}", strconv.Itoa(int(vCfg.Limit)), -1)
						mailId, _, err := db.MailAdd(
							-1,
							tm.uid,
							"系统邮件：任务奖励未领取",
							taskName,
							result3,
							constant.ItemAddTask,
						)
						if err != nil {
							logger.Error(err)
						}
						upMailShushuDate(mailId, result3)
					}
					delTask = append(delTask, int32(v.ID))
				}
			}
		}
		//解锁伙伴槽位
		if v.Target.TargetType == pb.TargetType_UNLOCK_PARTNER_SLOT {
			if v.Status == pb.TaskStatus_REWARDED || v.Status == pb.TaskStatus_COMPLETED {
				continue
			}
			num := 0
			for _, v := range global.GlobalCfg.Openslot.IDLevel {
				if v.ServerDay > 0 {
					if tm.level >= v.Needlevel && global.ServerListMgrInstance.GetServerOpenDay(tm.serverID) >= v.Openday && tm.mainBattleId >= v.Minbattle {
						num++
					}
				} else {
					if tm.level >= v.Needlevel && tm.createDay >= v.Openday && tm.mainBattleId >= v.Minbattle {
						num++
					}
				}
			}
			v.Target.Amount = num
			if v.Target.Amount >= v.Target.TargetNum {
				v.Target.Completed = true
				v.Status = pb.TaskStatus_COMPLETED
			} else {
				v.Status = pb.TaskStatus_IN_PROGRESS
			}
			tm.update(v)
		}
	}

	for _, v := range delTask {
		tm.del(v)
	}

	tm.check()
}
func upMailShushuDate(mail_id string, sourceInfo []protocol.Attachment) {
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
	shushumain.ReportData(shushuDate)
}
func (tm *TaskMgr) resetNewweek() {

	logger.Debugf("==新的一周刷新周周长任务 uid %d  name %s", tm.uid, tm.name)
	delTask := []int32{}
	for _, v := range tm.tasks {

		// 周常任务
		if v.Class == pb.TaskClass_WEEKLY || v.Class == pb.TaskClass_WEEKLY_ACTIVITY {
			vCfg, err := global.GetTaskCfg(int32(v.ID))
			if err == nil {
				if len(vCfg.Level) == 2 {
					if tm.level >= vCfg.Level[0] && tm.level <= vCfg.Level[1] {
						v.Target.Amount = 0
						v.Target.Completed = false
						v.Status = pb.TaskStatus_IN_PROGRESS

						tm.update(v)

						logger.Debugf("新的一周刷新周长任务，没有跨等级段 uid %d  name %s  任务ID %d", tm.uid, tm.name, v.ID)
					} else {
						delTask = append(delTask, int32(v.ID))
						logger.Debugf("新的一周刷新周长任务，跨等级段，删除任务  uid %d  name %s  任务ID %d", tm.uid, tm.name, v.ID)
					}
				}
			}
		}
	}

	for _, v := range delTask {
		tm.del(v)
	}

	tm.check()
}

// init 任务初始化, 角色上线时调用
func (tm *TaskMgr) init() {
	for _, cfg := range global.TaskCfg {

		if cfg.Job > 0 && tm.job != constant.MowJob(cfg.Job) {
			continue
		}

		// 非首轮任务
		if !cfg.FirstRound {
			continue
		}

		if task := tm.getTaskById(cfg.UUID); task == nil {
			// 等级段筛选
			if len(cfg.Level) == 2 {
				if tm.level < cfg.Level[0] || tm.level > cfg.Level[1] {
					continue
				}
			}

			// vip等级对应的任务
			if cfg.Class == int32(pb.TaskClass_VIP) {
				if tm.vip < cfg.VipLevel-1 {
					continue
				}
			}

			//七日狂欢对应的任务
			if cfg.Class == int32(pb.TaskClass_SEVENDAY) {
				if tm.createDay == 0 {
					tm.createDay = 1
				}
				if int32(tm.createDay) < int32(cfg.SevenDay) {
					continue
				}
			}

			t := &TaskInfo{
				ID:     int(cfg.UUID),
				Class:  pb.TaskClass(cfg.Class),
				Status: pb.TaskStatus_IN_PROGRESS,
				Target: &TargetInfo{
					TargetType: pb.TargetType(cfg.Type),
					TargetId:   int(cfg.Typeson),
					TargetNum:  int(cfg.Limit),
					Amount:     0,
					Completed:  false,
				},
			}

			if t.Class == pb.TaskClass_WEEKLY || t.Class == pb.TaskClass_WEEKLY_ACTIVITY {
				logger.Debugf("周常任务init ID:Name %d:%s 任务ID %d   任务 %v", tm.uid, tm.name, t.ID, t)
			}
			if t.Class == pb.TaskClass_DAILY || t.Class == pb.TaskClass_ROUTINE || t.Class == pb.TaskClass_DAILY_ACTIVITY {
				logger.Debugf("日常任务init ID:Name %d:%s 任务ID %d   任务 %v", tm.uid, tm.name, t.ID, t)
			}

			switch t.Target.TargetType {
			case pb.TargetType_ROLE_POWER_REQUIREMENT:
				t.Target.Amount = int(tm.combat)
			case pb.TargetType_ROLE_LEVEL_REQUIREMENT:
				t.Target.Amount = int(tm.level)
			}

			tm.tasks[int32(t.ID)] = t
		}
	}
}

func (tm *TaskMgr) check() {
	for _, v := range global.TaskCfg {

		if v.Job > 0 && tm.job != constant.MowJob(v.Job) {
			continue
		}
		// 非首轮任务
		if !v.FirstRound {
			continue
		}
		if task := tm.getTaskById(v.UUID); task == nil {
			tm.add(v.UUID, 0)
		}
	}
	// 对已存在的任务重算一次进度：DB 里存的 amount 可能是旧的（例如"穿戴X件装备"接取时被算成 0 且之后不再触发）
	for _, t := range tm.tasks {
		tm.checkProgress(t)
	}
}

func (tm *TaskMgr) getTaskById(id int32) *TaskInfo {
	t, ok := tm.tasks[id]
	if ok {
		return t
	}
	return nil
}

func (tm *TaskMgr) toPB() *pb.DBTask {
	dbtask := &pb.DBTask{}
	for _, v := range tm.tasks {
		dbtask.Tasks = append(dbtask.Tasks, &pb.TaskInfo{
			CfgId: int32(v.ID), Class: v.Class, Status: v.Status, RefreshTime: v.ResetTime.Unix(),
			Target: &pb.TargetInfo{
				Type: v.Target.TargetType, TargetId: int32(v.Target.TargetId),
				TargetNum: int32(v.Target.TargetNum), Amount: int32(v.Target.Amount), Completed: v.Target.Completed,
			},
		})
	}
	return dbtask
}

func (tm *TaskMgr) toPBTaskInfo(t *TaskInfo) *pb.TaskInfo {
	return &pb.TaskInfo{
		CfgId: int32(t.ID), Class: t.Class, Status: t.Status,
		Target: &pb.TargetInfo{
			Type: t.Target.TargetType, TargetId: int32(t.Target.TargetId),
			TargetNum: int32(t.Target.TargetNum),
			Amount:    int32(t.Target.Amount), Completed: t.Target.Completed,
		},
	}
}

func (tm *TaskMgr) fromPB(dbt *pb.DBTask) {
	for _, v := range dbt.Tasks {
		rt := time.Unix(v.RefreshTime, 0)
		tm.tasks[v.CfgId] = &TaskInfo{
			ID:        int(v.CfgId),
			Class:     v.Class,
			Status:    v.Status,
			ResetTime: rt,
			Target: &TargetInfo{
				TargetType: v.Target.Type,
				TargetId:   int(v.Target.TargetId),
				TargetNum:  int(v.Target.TargetNum),
				Amount:     int(v.Target.Amount),
				Completed:  v.Target.Completed,
			},
		}
	}
}

func (tm *TaskMgr) add(id int32, preTaskId int32) *TaskInfo {
	cfg, err := global.GetTaskCfg(id)
	if err != nil {
		return nil
	}

	for _, v := range tm.tasks {
		if v.ID == int(id) {
			return nil
		}
		vCfg, err := global.GetTaskCfg(int32(v.ID))
		if err != nil {
			return nil
		}
		//日活跃 周活跃任务不同等级互斥 Class相同&&Typeson相同
		if vCfg.Class == int32(pb.TaskClass_DAILY_ACTIVITY) || vCfg.Class == int32(pb.TaskClass_WEEKLY_ACTIVITY) ||
			vCfg.Class == int32(pb.TaskClass_DAILY) || vCfg.Class == int32(pb.TaskClass_WEEKLY) {
			if cfg.Typeson == vCfg.Typeson && cfg.Class == vCfg.Class && cfg.Type == vCfg.Type {
				return nil
			}
		}

	}

	if len(cfg.Level) == 2 {
		if tm.level < cfg.Level[0] || tm.level > cfg.Level[1] {
			return nil
		}
	}
	// vip等级对应的任务
	if cfg.Class == int32(pb.TaskClass_VIP) {
		if tm.vip < cfg.VipLevel-1 {
			return nil
		}
	}
	//七日狂欢对应的任务
	if cfg.Class == int32(pb.TaskClass_SEVENDAY) {
		if tm.createDay == 0 {
			tm.createDay = 1
		}
		//如果活动已经结束了，就不要给他派了
		actCfg, err := global.GetActivityCfg(int32(cfg.ActivID))
		if err != nil {
			return nil
		}
		// logger.Debug("actCfg.ResetTime", actCfg.ResetTime, "tm.createDay", tm.createDay, "cfg.SevenDay", cfg.SevenDay)
		if tm.createDay > actCfg.ResetTime {
			return nil
		}
		if int32(tm.createDay) < int32(cfg.SevenDay) {
			return nil
		}
	}
	t := &TaskInfo{
		ID: int(cfg.UUID), Class: pb.TaskClass(cfg.Class), Status: pb.TaskStatus_IN_PROGRESS,
		Target: &TargetInfo{TargetType: pb.TargetType(cfg.Type), TargetId: int(cfg.Typeson), TargetNum: int(cfg.Limit), Amount: 0, Completed: false},
	}

	if t.Class == pb.TaskClass_WEEKLY || t.Class == pb.TaskClass_WEEKLY_ACTIVITY {
		logger.Debugf("接取新周常任务 ID:Name %d:%s 任务ID %d  任务 %v", tm.uid, tm.name, t.ID, t)
	}

	if t.Class == pb.TaskClass_DAILY || t.Class == pb.TaskClass_ROUTINE || t.Class == pb.TaskClass_DAILY_ACTIVITY {
		logger.Debugf("接取新日常任务 ID:Name %d:%s 任务ID %d  任务 %v", tm.uid, tm.name, t.ID, t)
	}

	// 当前任务的继承任务是否需要重置计数条件
	var preTask *TaskInfo
	if preTaskId > 0 && !cfg.ResetMark {
		if preTask = tm.getTaskById(preTaskId); preTask != nil {
			t.Target.Amount = preTask.Target.Amount
		}
	}

	switch t.Target.TargetType {
	case pb.TargetType_ROLE_POWER_REQUIREMENT:
		t.Target.Amount = int(tm.totalCombat)
	case pb.TargetType_ROLE_LEVEL_REQUIREMENT:
		t.Target.Amount = int(tm.level)
	}
	tm.tasks[int32(t.ID)] = t

	packet := &pb.TaskAcceptPush{Task: tm.toPBTaskInfo(t)}
	tm.pushMsg(protocol.RouteTaskAcceptPush, packet)
	logger.Infof("add task %d", t.ID)
	tm.checkProgress(t)
	return t
}

func (tm *TaskMgr) del(id int32) {
	if t := tm.getTaskById(id); t != nil {

		packet := &pb.TaskRemovePush{Task: tm.toPBTaskInfo(t)}

		tm.pushMsg(protocol.RouteTaskRemovePush, packet)
		//删除任务
		delete(tm.tasks, id)
	}
}

func (tm *TaskMgr) checkProgress(t *TaskInfo) {
	if t.Status == pb.TaskStatus_COMPLETED || t.Status == pb.TaskStatus_REWARDED {
		return
	}

	if t.Target.Completed {
		t.Status = pb.TaskStatus_COMPLETED
	}

	if t.Class == pb.TaskClass_ROUTINE {
		if tm.checkTime(t) {
			t.Target.Completed = true
			t.Status = pb.TaskStatus_COMPLETED
			tm.update(t)
		}
		return
	}
	if t.Target.TargetType == pb.TargetType_SKILL_UPGRADING_REQUIREMENT {
		if skill := tm.getSkill(int32(t.Target.TargetId)); skill != nil {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
			t.Status = pb.TaskStatus_COMPLETED
		} else {
			for _, v := range tm.skills {
				if t.Target.TargetId/100 == int(v.CfgId)/100 {
					if int(v.CfgId)%100 >= t.Target.TargetId%100 {
						t.Target.Amount = t.Target.TargetNum
						t.Target.Completed = true
						t.Status = pb.TaskStatus_COMPLETED
					}
				}
			}
		}
	}
	if t.Target.TargetType == pb.TargetType_VIP_LEVEL {
		t.Target.Amount = int(tm.vip)
		tm.update(t)
	}
	if t.Target.TargetType == pb.TargetType_BATTLE_CHAPTERS_CLEARED {
		// 查询当前章节
		if tm.InstanceData != nil {
			if int(tm.InstanceData.BattleMap.BattleId) >= t.Target.TargetId {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
			}
		}
	}
	if t.Target.TargetType == pb.TargetType_TRIAL_TOWER_FLOORS_CLEARED {
		// 查询当前章节
		if tm.InstanceData != nil {
			if int(tm.InstanceData.TrialTowerMap.BattleId) >= t.Target.TargetId {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
			}
		}
	}

	if t.Target.TargetType == pb.TargetType_RELIC_CHAPTERS_CLEARED {
		// 查询当前章节
		if tm.InstanceData != nil {
			if int(tm.InstanceData.RelicMap.BattleId) >= t.Target.TargetId {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
			}
		}
	}
	if t.Target.TargetType == pb.TargetType_BOSS_LEVELS_CLEARED {
		// 查询当前章节
		if tm.InstanceData != nil {
			if int(tm.InstanceData.BossMap.BattleId) >= t.Target.TargetId {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
			}
		}
	}
	//任务：伙伴等级
	if t.Target.TargetType == pb.TargetType_UPGRADE_ANY_PARTNER_LEVEL {
		if tm.InstanceData != nil {
			for _, v := range tm.partnerSlots {
				if v.PartnerUid > 0 {
					if par := tm.getPartnerInfo(v.PartnerUid); par != nil {
						if par.Level() >= int32(t.Target.TargetId) {
							t.Target.Amount = t.Target.TargetNum
							t.Target.Completed = true
							break
						}
					}
				}
			}

		}
	}
	//解锁伙伴槽位
	if t.Target.TargetType == pb.TargetType_UNLOCK_PARTNER_SLOT {
		num := 0
		for _, v := range global.GlobalCfg.Openslot.IDLevel {
			if v.ServerDay > 0 {
				if tm.level >= v.Needlevel && global.ServerListMgrInstance.GetServerOpenDay(tm.serverID) >= v.Openday && tm.mainBattleId >= v.Minbattle {
					num++
				}
			} else {
				if tm.level >= v.Needlevel && tm.createDay >= v.Openday && tm.mainBattleId >= v.Minbattle {
					num++
				}
			}
		}
		t.Target.Amount = num
		tm.update(t)
	}
	//装备强化等级
	if t.Target.TargetType == pb.TargetType_ANY_EQUIP_STRENGTHENED_TO_LEVEL {
		for _, v := range tm.equips {
			if v.Enhance >= int32(t.Target.TargetId) {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
				break
			}
		}
	}
	//所有装备强化到多少级
	if t.Target.TargetType == pb.TargetType_EQUIPMENT_STRENGTHENING_LEVEL_REQUIREMENT {
		// 任务上报：全身装备强化等级需求
		minEnhance := int32(math.MaxInt32)
		for i, equip := range tm.equips {
			if i == 0 {
				continue
			}
			if equip.Enhance < minEnhance {
				minEnhance = equip.Enhance
			}
		}
		if minEnhance >= int32(t.Target.TargetNum) {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
			tm.update(t)
		}
	}
	//任意装备升星到多少级
	if t.Target.TargetType == pb.TargetType_EQUIP_STAR_TO_X {
		for _, v := range tm.equips {
			if v.Stars >= int32(t.Target.TargetId) {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
				break
			}
		}
	}
	//所有装备升星到多少级
	if t.Target.TargetType == pb.TargetType_EQUIPMENT_UPGRADING_STAR_LEVEL_REQUIREMENT {
		// 任务上报：全身装备升星等级需求
		minStars := int32(math.MaxInt32)
		for i, equip := range tm.equips {
			if i == 0 {
				continue
			}
			if equip.Stars < minStars {
				minStars = equip.Stars
			}
		}
		if minStars >= int32(t.Target.TargetNum) {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
			tm.update(t)
		}
	}
	if t.Target.TargetType == pb.TargetType_RELIC_TOTAL_STAR {
		allStar := 0
		for i := 0; i < len(tm.RelicMap.Evaluation); i++ {
			allStar += int(tm.RelicMap.Evaluation[i].Value)
		}
		t.Target.Amount = allStar
		if allStar >= t.Target.TargetNum {
			t.Target.Completed = true
			tm.update(t)
		}
	}
	//宠物等级
	if t.Target.TargetType == pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL {
		for _, v := range tm.PetMgr.pet {
			if v.Level >= int32(t.Target.TargetNum) {
				t.Target.Amount = t.Target.TargetNum
				t.Target.Completed = true
				break
			}
		}
		tm.update(t)
	}
	//宝箱等级
	if t.Target.TargetType == pb.TargetType_CHEST_LEVEL_UPGRADED {
		t.Target.Amount = int(tm.box.NowLevel)
		tm.update(t)
	}
	//查询购买次数
	if pb.TargetType(t.Target.TargetType) == pb.TargetType_SEVENDAY_TASK_GIFT_PURCHASED {
		//查询购买次数
		ress, err := db.StoreTotal(tm.uid, constant.ShopTypeSevenDaysFun)
		if err != nil {
			return
		}
		t.Target.Amount = len(ress)
		tm.update(t)
	}
	//主线任务
	if t.Target.TargetType == pb.TargetType_MAIN_BATTLE {
		if int(tm.mainBattleId) >= t.Target.TargetId {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
		}
	}

	//任务-进入工会
	if t.Target.TargetType == pb.TargetType_JOIN_GUILD {
		if tm.guildId > 0 {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
		}
		tm.update(t)
	}
	//任务-上阵伙伴的数量
	if t.Target.TargetType == pb.TargetType_PARTNER_ON_STAGE_TIMES {
		num := 0
		for _, v := range tm.partnerSlots {
			if v.PartnerUid > 0 {
				num++
			}
		}
		if num >= t.Target.TargetNum {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
			tm.update(t)
		}
	}
	//任务-穿戴装备
	if t.Target.TargetType == pb.TargetType_EQUIP_ITEMS_WORN_TOTAL {
		num := 0
		for _, v := range tm.equips {
			if v.ItemId > 0 {
				num++
			}
		}
		//无论是否达标都回写当前进度：原实现只有达标才赋值，导致接取任务时进度永远是 0
		t.Target.Amount = num
		if num >= t.Target.TargetNum {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
		}
		logger.Infof("[TASK][wear] uid=%d task=%d worn=%d targetNum=%d amount=%d", tm.uid, t.ID, num, t.Target.TargetNum, t.Target.Amount)
		tm.update(t)
	}
	if t.Target.TargetNum <= t.Target.Amount || t.Status == pb.TaskStatus_COMPLETED {
		t.Target.Completed = true
		t.Status = pb.TaskStatus_COMPLETED
		tm.update(t)
	}

	// 遗迹第X章所有关卡均三星
	if t.Target.TargetType == pb.TargetType_RELIC_CHAPTERS_ALL_3STARS {
		if t.Target.Amount >= t.Target.TargetNum {
			t.Target.Amount = t.Target.TargetNum
			t.Target.Completed = true
			t.Status = pb.TaskStatus_COMPLETED
			tm.update(t)
		}
	}

	// logger.Debugf("任务进度 ID=%d  Amount=%d  TargetNum=%d", t.ID, t.Target.Amount, t.Target.TargetNum)
}

func (tm *TaskMgr) update(t *TaskInfo) {
	packet := &pb.TaskUpdatePush{Task: tm.toPBTaskInfo(t)}
	tm.pushMsg(protocol.RouteTaskUpdatePush, packet)
}

// 完成任务
func (tm *TaskMgr) complate(id int32) error {
	t := tm.getTaskById(id)
	if t == nil {
		return errutil.ErrTaskNotFound
	}

	if t.Status != pb.TaskStatus_COMPLETED {
		return errutil.ErrTaskStatus
	}
	cfg, err := global.GetTaskCfg(int32(t.ID))
	if err != nil {
		return errutil.ErrTaskNotFound
	}

	if cfg.Class == int32(pb.TaskClass_ACHIEVEMENT) {
		if !tm.checkTime(t) {
			return errutil.ErrTaskStatus
		}
	}
	t.Status = pb.TaskStatus_REWARDED
	tm.update(t)

	for _, v := range cfg.Reward {
		tm.addSource(v.Key, int64(v.Value), constant.ItemAddTask)
	}

	if cfg.NextAchievement != 0 {
		tm.add(cfg.NextAchievement, int32(t.ID))
	}
	// 活跃度任务
	if cfg.ActivityReward.Value > 0 {
		if t.Class == pb.TaskClass_DAILY {
			tm.dispatch(pb.TargetType_DAILY_ACTIVITY_LEVEL, 0, cfg.ActivityReward.Value)
		}
		if t.Class == pb.TaskClass_WEEKLY {
			tm.dispatch(pb.TargetType_WEEKLY_ACTIVITY_LEVEL, 0, cfg.ActivityReward.Value)
		}
		if t.Class == pb.TaskClass_SEVENDAY {
			tm.dispatch(pb.TargetType_SEVENDAY_TASK_ALL_COMPLETED, 0, cfg.ActivityReward.Value)
		}
	}
	if t.Class == pb.TaskClass_MAINLINE {
		tm.taskProgress = t.ID
	}
	//上报
	tm.upTaskShushuDate(cfg)
	return nil
}

func (p *Player) upTaskShushuDate(cfg *global.TaskConfig) {
	if p.bCreate {
		return
	}
	//创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	//upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	//在合并点其他的数据

	// 替换cfg.TaskDes中的{0}为cfg.Limit的值
	taskName := cfg.TaskDes
	if cfg.Limit > 0 {
		taskName = strings.Replace(taskName, "{0}", strconv.Itoa(int(cfg.Limit)), -1)
	}
	shushuDate["task_name"] = taskName
	shushuDate["name"] = "task_complete"
	var taskType string
	classIdx := int32(cfg.Class)
	if classIdx >= 0 && int(classIdx) <= len(constant.TaskName) {
		taskType = constant.TaskName[classIdx]
	} else {
		taskType = ""
	}
	shushuDate["task_type"] = taskType
	shushuDate["task_id"] = cfg.UUID
	itemInfo := make([]map[string]interface{}, 0)
	for _, v := range cfg.Reward {
		if v.Value > 0 {
			//查询itemname
			if itemCfg, err := global.GetItemCfg(int32(v.Key)); err == nil {
				item := map[string]interface{}{
					"item_id":   itemCfg.ID,
					"item_name": itemCfg.Name,
					"item_num":  v.Value,
				}
				itemInfo = append(itemInfo, item)
			}
		}
	}
	shushuDate["get_item_info"] = itemInfo

	shushumain.ReportData(shushuDate)

}

// 时间段比较
func (tm *TaskMgr) checkTime(t *TaskInfo) bool {
	cfg, err := global.GetTaskCfg(int32(t.ID))
	if err != nil {
		return false
	}

	if cfg.Class == int32(pb.TaskClass_ROUTINE) {
		if start, err := time.Parse("15:04", cfg.LimitTime.Start); err == nil {
			if start.Hour() > time.Now().Hour() {
				return false
			} else if start.Hour() == time.Now().Hour() && start.Minute() > time.Now().Minute() {
				return false
			}
		}
		if end, err := time.Parse("15:04", cfg.LimitTime.End); err == nil {
			if end.Hour() < time.Now().Hour() {
				return false
			} else if end.Hour() == time.Now().Hour() && end.Minute() < time.Now().Minute() {
				return false
			}
		}
	}
	return true
}

// 查询某一个任务的完成情况
func (tm *TaskMgr) queryTask(id int32) *TaskInfo {
	if v, ok := tm.tasks[id]; ok {
		return v
	}
	return nil
}
