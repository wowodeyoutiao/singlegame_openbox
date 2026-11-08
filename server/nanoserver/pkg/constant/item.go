package constant

type ItemID int32

// 定义道具种类的枚举 item.Type
const (
	// 代表不进入包裹的数值类道具
	NotInBag int32 = iota
	// 代表进入包裹且可被销毁的道具
	InBagDisposable
	// 代表进入包裹且可被销毁和使用的道具
	InBagDisposableAndUsable
	// 代表进入包裹、可被销毁、带有时效的道具
	InBagTimeLimited
	// 代表进入包裹、可被销毁和使用、带有时效的道具
	InBagDisposableUsableTimeLimited
	// 代表装备类，该类道具可装备、可分解
	Equipment
	// 代表宝石类，该类道具可镶嵌、可合成
	Gem
	// 代表伙伴碎片，该道具可合成，成为1个整卡伙伴
	PartnerFragment
	// 代表伙伴整卡
	ItemPartner
	// 代表战宠
	ItemPet
	//称号碎片
	ItemMainTypeTitleFragment
	// 可采集宝箱
	ItemMainTypeMiningBox
	// 头像
	ItemMainTypeAvatar
	// 头像框
	ItemMainTypeAvatarFrame
)

// 定义道具子类型常量 item.SonType
const (
	ItemTypeNone       = 0 // 代表数值类道具
	ItemTypeMaterial   = 1 // 代表材料类道具
	ItemTypeEquipment  = 2 // 代表装备类道具
	ItemTypeGem        = 3 // 代表宝石类道具
	ItemTypeConsumable = 4 // 代表消耗类道具
)

const (
	ItemIDLuckyBagLow  ItemID = 25064 // 初级福袋 蓝色
	ItemIDLuckyBagMid  ItemID = 25065 // 中级福袋 紫色
	ItemIDLuckyBagHigh ItemID = 25066 // 高级福袋 橙色
	ItemIDOfficialExp  ItemID = 52002 // 官职经验
)

func (id ItemID) Int32() int32 {
	return int32(id)
}
