package tracking

// Tracker 简洁的埋点接口，统一事件与用户属性上报入口。
// accountId: 用户账号ID（如手机号、邮箱、用户ID等），用于标识已登录用户。
// distinctId: 设备ID或匿名用户ID（如设备UUID、cookie），用于标识设备或未登录用户。
// 使用建议：
//   - 未登录用户：仅传 distinctId，accountId 为空。
//   - 已登录用户：同时传 accountId 和 distinctId，建立账号与设备的关联。
//   - 两者至少需传一个，不能同时为空，否则数据无效。
//
// 场景说明：
//   - accountId 用于跨设备、跨平台统一识别用户，适合账号体系场景。
//   - distinctId 用于设备或匿名用户追踪，适合未登录或单设备场景。
//
// 推荐：登录后同时传递 accountId 和 distinctId，保证用户行为的完整追踪与分析。
type Tracker interface {
	// E 事件上报
	// accountID: 用户账号ID（已登录时传入）
	// distinctID: 设备ID或匿名用户ID（必传）
	// event: 事件名称
	// props: 事件属性
	E(accountID, distinctID, event string, props M) error

	// Set 设置用户属性（可覆盖）
	// 适用于：用户基本信息更新（如昵称、等级、性别等）
	Set(accountID, distinctID string, props M) error

	// SetOnce 设置用户属性（只设置一次）
	// 适用于：用户首次行为（如注册时间、首次来源等），已设置则忽略
	SetOnce(accountID, distinctID string, props M) error

	// Add 累加用户属性
	// 适用于：数值型属性累加（如付费总额、登录次数等）
	Add(accountID, distinctID string, props M) error

	// Flush 刷新缓冲区，确保数据及时上报
	Flush() error

	// Close 关闭埋点实例，释放资源
	Close() error
}

// M 属性字典的简写，用于传递事件属性和用户属性
type M map[string]interface{}
