package dispatcher

type Topic string

const (
	// TopicRechargeSuccess Topic = "recharge.success" // 充值成功

	// TopicActivityPush Topic = "activity.push" // 活动推送

	TopicPayRefreshPush        Topic = "pay.refresh.push"           // 支付刷新推送, 不需要迁移到 p.eventMgr
	TopicUserShowRefresh       Topic = "user.show.refresh"          // 用户显示信息刷新, 不需要迁移到 p.eventMgr
	TopicPVEWarMapBossDieCount Topic = "pve.war.map.boss.die.count" // PVE战场地图Boss死亡次数, 不需要迁移到 p.eventMgr
	TopicActivatePurchasedSlot Topic = "activate.purchased.slot"    // 激活购买的洗练槽位

	TopicCreateCharacter Topic = "payment.create.character" // 同步角色到支付系统
	TopicAdClick         Topic = "ad.click"                 // 广告点击进来创角
)

func (t Topic) String() string {
	return string(t)
}

type ActivityType int

const (
	ActivityCumulativeRecharge ActivityType = iota + 1 // 累计充值活动
	ActivityServerCelebration                          // 开服狂欢活动
)
