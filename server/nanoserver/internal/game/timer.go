package game

import (
	"log"
	"nanoserver/internal/web/provide"
	"nanoserver/pkg/constant"

	"github.com/robfig/cron"
)

// 每隔5秒执行一次：*/5 * * * * ?
// 每隔1分钟执行一次：0 */1 * * * ?
// 每天23点执行一次：0 0 23 * * ?
// 每天凌晨1点执行一次：0 0 1 * * ?
// 每月1号凌晨1点执行一次：0 0 1 1 * ?
// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?

func InitTimer() {
	// 创建一个定时任务调度器 ,表达式格式为：秒 分 时 日 月 星期几 年份
	c := cron.New()

	var err error
	// 每天零点
	if err = c.AddFunc("0 0 0 * * *", TimerFor0Oclock); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 每月1日
	if err = c.AddFunc("0 0 1 1 *", TimerFor1stEveryMonth); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 每周一
	if err = c.AddFunc("0 0 0 * * 1", TimerForMonday); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}

	// pvp开启 10:00:00
	if err = c.AddFunc("0 0 10 * * *", matchManager.matchingOn); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// pvp关闭 20:30:00
	if err = c.AddFunc("0 30 20 * * *", matchManager.matchingOff); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// // 21:00准时发排行榜奖励
	// if err = c.AddFunc("0 0 21 * * *", provide.CronSendRankAwardOne); err != nil {
	// 	log.Panicln("添加定时任务失败:", err)
	// 	return
	// }
	// 战力榜和即时赛榜和回合制，晚上10点发
	if err = c.AddFunc("0 0 22 * * *", provide.CronSendRankAwardPvp); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 世界boss榜，晚上9点发
	if err = c.AddFunc("0 0 21 * * *", provide.CronSendRankAwardBoss); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 一战到底榜，晚上21点发
	if err = c.AddFunc("0 0 21 * * *", provide.CronSendRankAwardOne); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 当日达，晚上0点发
	if err = c.AddFunc("0 0 0 * * *", provide.CronSendRankAwardActivity); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}

	// // 每分钟都执行
	// if err = c.AddFunc("0 */1 * * * ", provide.CronSendRankAwardSeason); err != nil {
	// 	log.Panicln("添加定时任务失败:", err)
	// 	return
	// }

	// 每秒都执行
	if err = c.AddFunc("*/1 * * * * ?", TimerForSecond); err != nil {
		log.Panicln("添加定时任务失败:", err)
		return
	}
	// 启动定时任务调度器
	c.Start()

	logger.Info("定时任务调度器已启动...")
}

// 每天0点 定时任务开始执行
func TimerFor0Oclock() {
	log.Println("每天0点 定时任务开始执行...")
	defaultManager.TimerFor0Oclock()
	worldBossMgr.reset()
	rankType := []int64{
		constant.RankTypeWorldBoss,
		constant.RankTypeWorldBossKf,
		constant.RankTypeGightEnd,
		constant.RankTypeGightEndlayers,
	}
	provide.DelRank(rankType)
	provide.DelayDelRank()
	DropControlInstance.Reset()

	playerMgr.IteratePlayers(func(val *Player) bool {
		if val.curState == PlayerStatusIdle || val.curState == PlayerStatusPlaying {
			val.SyncUserAvatar(true)
		}
		return true
	})
}

// 每月1日 定时任务开始执行
func TimerFor1stEveryMonth() {
	log.Println("每月1日1点 定时任务开始执行...")
}

// 每周一 定时任务开始执行
func TimerForMonday() {
	// 试练塔每周一重置
	log.Println("每周一 定时任务开始执行...")
	provide.CronSendRankAwardWeekly()
	robotMgr.updateAllServerRank()
	defaultManager.resetPvpScore()
}

// 每秒 定时任务开始执行
func TimerForSecond() {
	// itemRecoderMgr.update()
	worldBossMgr.checkWorldBossTime()
}
