package constant

// 工会成员角色
const (
	GuildRoleLeader     = 1 // 会长
	GuildRoleViceLeader = 2 // 副会长
	GuildRoleMember     = 3 // 普通成员
)

// 工会成员状态
const (
	GuildMemberStatusNormal = 1 // 正常
	GuildMemberStatusLeave  = 2 // 退出
	GuildMemberStatusKicked = 3 // 被踢
)

// 工会申请状态
const (
	GuildApplyStatusPending  = 1 // 待处理
	GuildApplyStatusApproved = 2 // 已同意
	GuildApplyStatusRejected = 3 // 已拒绝
)

// 工会操作类型
const (
	GuildActionCreate  = "create"  // 创建工会
	GuildActionDismiss = "dismiss" // 解散工会
	GuildActionJoin    = "join"    // 加入工会
	GuildActionLeave   = "leave"   // 离开工会
	GuildActionKick    = "kick"    // 踢出成员
	GuildActionPromote = "promote" // 升职
	GuildActionDemote  = "demote"  // 降职
	GuildActionSetting = "setting" // 修改设置
)

// 工会默认设置
const (
	GuildDefaultMaxMembers = 100 // 默认最大成员数
	GuildDefaultLevel      = 1   // 默认工会等级
	GuildDefaultPoints     = 0   // 默认工会积分
)

// 工会名称限制
const (
	GuildNameMinLength = 2  // 工会名称最小长度
	GuildNameMaxLength = 12 // 工会名称最大长度
)

// 工会描述限制
const (
	GuildDescMinLength = 0   // 工会描述最小长度
	GuildDescMaxLength = 200 // 工会描述最大长度
)

// 工会相关道具ID
const (
	ItemGuildCreate = 1001 // 创建工会道具
	ItemGuildRename = 1002 // 工会改名道具
)

// 工会相关消耗
const (
	CostGuildCreate = 100000 // 创建工会消耗金币
	CostGuildRename = 50000  // 工会改名消耗金币
)

// 工会日志保存时间(天)
const GuildLogSaveDays = 30

// 工会申请有效期(小时)
const GuildApplyExpireHours = 24

// 工会成员列表分页大小
const GuildMemberPageSize = 20

// 工会申请列表分页大小
const GuildApplyPageSize = 20

// 工会日志列表分页大小
const GuildLogPageSize = 50

// 工会列表分页大小
const GuildListPageSize = 20

// 工会相关事件
const (
	EventGuildCreate  = "guild_create"  // 创建工会
	EventGuildDismiss = "guild_dismiss" // 解散工会
	EventGuildJoin    = "guild_join"    // 加入工会
	EventGuildLeave   = "guild_leave"   // 离开工会
	EventGuildKick    = "guild_kick"    // 踢出成员
	EventGuildPromote = "guild_promote" // 升职
	EventGuildDemote  = "guild_demote"  // 降职
	EventGuildSetting = "guild_setting" // 修改设置
)
