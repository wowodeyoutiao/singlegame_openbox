package constant

// OfficialSuppressionState 表示官职压制状态
type OfficialSuppressionState int32

const (
	SuppressionNone     OfficialSuppressionState = 0 // 无压制
	SuppressionApplied  OfficialSuppressionState = 1 // A 压制 B
	SuppressionReversed OfficialSuppressionState = 3 // A 被 B 压制（即 B 压制 A）
)
