package game

import (
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/randx"
	"nanoserver/pkg/utils"
	"nanoserver/protocol/protobuf/pb"
	"sort"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
)

type RefineV2 struct {
	purchasedSlots     *utils.SafeMap[int32, struct{}]  // 购买的槽位
	attrs              *utils.SafeMap[int32, *pb.Stats] // 使用三元组作为key: [CharacterSlot, EquipmentSlot, RefinedSlot] -> Stats
	tempAttrs          *pb.Stats                        // 临时用：洗练的属性
	tempPurchasedSlots int32                            // 临时用：购买的槽位
	lastRefineType     pb.RefineType                    // 最后一次洗练类型
}

func NewRefineV2() *RefineV2 {
	return &RefineV2{
		purchasedSlots: utils.NewSafeMap[int32, struct{}](),
		attrs:          utils.NewSafeMap[int32, *pb.Stats](),
	}
}

// validateIndex 验证索引范围
func (r *RefineV2) validateIndex(characterSlot pb.CharacterSlot, equipSlot pb.EquipSlot, refinedSlot pb.RefinedSlot) error {
	if characterSlot <= pb.CharacterSlot_CHARACTER_SLOT_UNSPECIFIED || characterSlot >= pb.CharacterSlot_CHARACTER_SLOT_MAX {
		return errutil.ErrInvalidCharacterSlot
	}
	if equipSlot <= pb.EquipSlot_INVALID_EquipSlot || equipSlot >= pb.EquipSlot_MaxEquipSlot {
		return errutil.ErrInvalidEquipmentSlot
	}
	if refinedSlot <= pb.RefinedSlot_REFINED_SLOT_UNSPECIFIED || refinedSlot >= pb.RefinedSlot_REFINED_SLOT_MAX {
		return errutil.ErrInvalidRefinedSlot
	}
	return nil
}

// calculateKey 计算组合key: characterSlot * 10000 + equipSlot * 100 + refinedSlot
func (r *RefineV2) calculateKey(characterSlot pb.CharacterSlot, equipSlot pb.EquipSlot, refinedSlot pb.RefinedSlot) int32 {
	return int32(characterSlot)*10000 + int32(equipSlot)*100 + int32(refinedSlot)
}

// parseKey 从组合key中解析出原始索引
func (r *RefineV2) parseKey(key int32) (pb.CharacterSlot, pb.EquipSlot, pb.RefinedSlot) {
	characterSlot := pb.CharacterSlot(key / 10000)
	remaining := key % 10000
	equipSlot := pb.EquipSlot(remaining / 100)
	refinedSlot := pb.RefinedSlot(remaining % 100)
	return characterSlot, equipSlot, refinedSlot
}

// get 获取洗练属性
func (r *RefineV2) get(characterSlot pb.CharacterSlot, equipSlot pb.EquipSlot, refinedSlot pb.RefinedSlot) (*pb.Stats, bool) {
	key := r.calculateKey(characterSlot, equipSlot, refinedSlot)
	return r.attrs.Get(key)
}

// Refine 洗练
func (r *RefineV2) Refine(req *pb.RefineRequest, p *Player) (stats *pb.Stats, err error) {
	var (
		itemID int32
		num    int64
	)

	if err := r.validateIndex(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot()); err != nil {
		return nil, err
	}

	// 判断是否是购买的槽位
	for _, v := range global.GlobalCfg.GearRefining.Keys {
		// 1代表等级限制
		if v.Type == 1 && v.Slot == int32(req.GetRefinedSlot()) {
			if p.level < v.Limit {
				logger.Errorf("Refined level not enough: playerLevel %v, requiredLevel %v, characterSlot %v, equipSlot %v, refinedSlot %v", p.level, v.Limit, req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
				return nil, errutil.ErrLevel
			}
		}

		// 2是需要花钱购买的槽位, 3是花元宝购买的
		if (v.Type == 2 || v.Type == 3) && v.Slot == int32(req.GetRefinedSlot()) {
			key := r.calculateKey(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
			if !r.purchasedSlots.Has(key) {
				logger.Errorf("Refined slot not purchased: characterSlot %v, equipSlot %v, refinedSlot %v", req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
				return nil, errutil.ErrBadRequest
			}
		}
	}
	// 资源检查并执行洗练
	switch req.GetType() {
	case pb.RefineType_REFINE_TYPE_RANDOM: // 随机洗练
		itemID = global.GlobalCfg.BaseValue.ValueRefine.RandomID
		num = global.GlobalCfg.BaseValue.ValueRefine.RandomNum
		if !p.checkTakeSource(itemID, int64(num)) {
			return nil, errutil.ErrItemNotEnough
		}
		p.takeSource(itemID, int64(num), constant.ModelRefined.String()) // 扣除资源
		stats, err = r.randomRefine(req, p.level)
		p.tasks.dispatch(pb.TargetType_NORMAL_REFINE_TIMES, 0, 1)
		p.refine.lastRefineType = req.GetType()

	case pb.RefineType_REFINE_TYPE_DIRECTIONAL: // 定向洗练
		itemID = global.GlobalCfg.BaseValue.ValueRefine.AssignID
		num = global.GlobalCfg.BaseValue.ValueRefine.AssignNum
		if !p.checkTakeSource(itemID, int64(num)) {
			return nil, errutil.ErrItemNotEnough
		}
		p.takeSource(itemID, int64(num), constant.ModelRefined.String()) // 扣除资源
		stats, err = r.directionalRefine(req, p.level)
		p.tasks.dispatch(pb.TargetType_SPECIAL_REFINE_TIMES, 0, 1)
		p.refine.lastRefineType = req.GetType()

	default:
		logger.Warnf("Refine invalid refine type: %v", req.GetType())
		return nil, errutil.ErrBadRequest
	}
	p.addUpRecord(constant.RefineNum) // 增加统计

	p.tasks.dispatch(pb.TargetType_TOTAL_REFINE_TIMES, 0, 1)

	return
}

// RefineSelect 替换洗练属性
func (s *RefineV2) RefineSelect(req *pb.RefineSelectRequest) (resp *pb.RefineSelectResponse, err error) {
	resp = &pb.RefineSelectResponse{
		CharacterSlot: req.GetCharacterSlot(),
		EquipSlot:     req.GetEquipSlot(),
		RefinedSlot:   req.GetRefinedSlot(),
	}

	if err := s.validateIndex(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot()); err != nil {
		return resp, err
	}

	if s.tempAttrs == nil {
		return resp, errutil.ErrRefineNoTempAttrs // 没有临时属性，不能替换
	}

	key := s.calculateKey(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
	if proto.Equal(s.tempAttrs, req.GetSelectedAttributes()) {
		s.attrs.Set(key, req.GetSelectedAttributes())
		s.tempAttrs = nil
	}
	if stored, ok := s.attrs.Get(key); ok {
		resp.SelectedAttributes = stored
	}

	return resp, nil
}

// ToPb 转换为协议格式
func (r *RefineV2) ToPb() ([]*pb.CharacterRefineData, []*pb.RefineSlot) {
	var (
		characterRefineDataList []*pb.CharacterRefineData
		resultPurchasedSlots    []*pb.RefineSlot
	)
	characterMap := make(map[pb.CharacterSlot]map[pb.EquipSlot]map[pb.RefinedSlot]*pb.Stats)

	r.attrs.Range(func(key int32, value *pb.Stats) bool {
		characterSlot, equipSlot, refinedSlot := r.parseKey(key)
		if _, exists := characterMap[characterSlot]; !exists {
			characterMap[characterSlot] = make(map[pb.EquipSlot]map[pb.RefinedSlot]*pb.Stats)
		}
		if _, exists := characterMap[characterSlot][equipSlot]; !exists {
			characterMap[characterSlot][equipSlot] = make(map[pb.RefinedSlot]*pb.Stats)
		}
		characterMap[characterSlot][equipSlot][refinedSlot] = value
		return true
	})

	for charSlot, equipMap := range characterMap {
		charData := &pb.CharacterRefineData{
			CharacterSlot: pb.CharacterSlot(charSlot),
		}
		for equipSlot, refinedMap := range equipMap {
			equipData := &pb.EquipRefineData{
				EquipSlot: equipSlot,
			}
			for refinedSlot, stats := range refinedMap {
				refinedData := &pb.RefineData{
					RefinedSlot: refinedSlot,
					Stats: &pb.Stats{
						Attr: stats.Attr,
						Num:  stats.Num,
					},
				}
				equipData.RefinedStats = append(equipData.RefinedStats, refinedData)
			}
			charData.EquipRefines = append(charData.EquipRefines, equipData)
		}
		characterRefineDataList = append(characterRefineDataList, charData)
	}

	r.purchasedSlots.Range(func(key int32, _ struct{}) bool {
		characterSlot, equipSlot, refinedSlot := r.parseKey(key)
		resultPurchasedSlots = append(resultPurchasedSlots, &pb.RefineSlot{
			CharacterSlot: characterSlot,
			EquipSlot:     equipSlot,
			RefinedSlot:   refinedSlot,
		})
		return true
	})

	return characterRefineDataList, resultPurchasedSlots
}

// ToDB 转换为数据库存储结构
func (r *RefineV2) ToDB(p *Player) *model.Refine {
	purchasedSlots := make([]int32, 0, 10)
	p.refine.purchasedSlots.Range(func(key int32, _ struct{}) bool {
		purchasedSlots = append(purchasedSlots, key)
		return true
	})

	sort.Slice(purchasedSlots, func(i, j int) bool {
		return purchasedSlots[i] < purchasedSlots[j]
	})

	attrs := make(map[int32]*pb.Stats)
	r.attrs.Range(func(key int32, value *pb.Stats) bool {
		attrs[key] = value
		return true
	})

	return &model.Refine{
		UID:            p.uid,
		PurchasedSlots: purchasedSlots,
		Attrs:          attrs,
	}
}

// FromDB 从数据库加载数据
func (r *RefineV2) FromDB(refineModel *model.Refine) error {
	if refineModel == nil {
		return nil
	}
	if r.purchasedSlots == nil {
		r.purchasedSlots = utils.NewSafeMap[int32, struct{}]()
	}
	for _, v := range refineModel.PurchasedSlots {
		r.purchasedSlots.Set(v, struct{}{})
	}

	r.attrs = utils.NewSafeMap[int32, *pb.Stats]()
	for key, value := range refineModel.Attrs {
		r.attrs.Set(key, value)
	}
	return nil
}

// randomRefine 随机洗练
func (r *RefineV2) randomRefine(req *pb.RefineRequest, level int32) (stats *pb.Stats, err error) {
	// 获取对应槽位的等级
	randomCreatCfgV2, err := global.GetRandomCreatCfgV2(int32(req.GetEquipSlot()), level)
	if err != nil {
		return nil, err
	}

	attributeTypePicker := randx.New(randomCreatCfgV2.AttributeTypeSet, func(item global.KeyValue) int { return int(item.Value) })
	attributeType, err := attributeTypePicker.Pick()
	if err != nil {
		logger.Warnf("randomRefine pick kind error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	refineSegmentPicker := randx.New(randomCreatCfgV2.RefiningSegmentRandomSet, func(item global.KeyValue) int { return int(item.Value) })
	refineSegment, err := refineSegmentPicker.Pick()
	if err != nil {
		logger.Warnf("randomRefine pick refineSegment error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	attributeSegmentCfg, err := global.GetAttributeSegmentCfgV2(attributeType.Key, level, refineSegment.Key)
	if err != nil {
		logger.Warnf("randomRefine get attributeSegmentCfg error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	value, err := randx.RandIntRange(attributeSegmentCfg.ValueMin, attributeSegmentCfg.ValueMax)
	if err != nil {
		logger.Warnf("randomRefine RandIntRange error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	r.tempAttrs = &pb.Stats{
		Attr: pb.Attribute(attributeSegmentCfg.AttributeType),
		Num:  int32(value),
	}

	return r.tempAttrs, nil
}

// directionalRefine 定向洗练
func (r *RefineV2) directionalRefine(req *pb.RefineRequest, level int32) (stats *pb.Stats, err error) {
	// 获取对应槽位的等级
	randomCreatCfgV2, err := global.GetRandomCreatCfgV2(int32(req.GetEquipSlot()), level)
	if err != nil {
		return nil, err
	}

	refineSegmentPicker := randx.New(randomCreatCfgV2.RefiningSegmentDirectionalSet, func(item global.KeyValue) int { return int(item.Value) })
	refineSegment, err := refineSegmentPicker.Pick()
	if err != nil {
		logger.Warnf("randomRefine pick refineSegment error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	attr, exist := r.get(req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
	if !exist {
		logger.Warnf("directionalRefine no exist attr for characterSlot %v, equipSlot %v, refinedSlot %v", req.GetCharacterSlot(), req.GetEquipSlot(), req.GetRefinedSlot())
		return nil, errutil.ErrRefineNoExist
	}

	attributeSegmentCfg, err := global.GetAttributeSegmentCfgV2(int32(attr.Attr), level, refineSegment.Key)
	if err != nil {
		logger.Warnf("randomRefine get attributeSegmentCfg error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	value, err := randx.RandIntRange(attributeSegmentCfg.ValueMin, attributeSegmentCfg.ValueMax)
	if err != nil {
		logger.Warnf("randomRefine RandIntRange error: %v", err)
		return nil, errutil.ErrServerInternal
	}

	r.tempAttrs = &pb.Stats{
		Attr: pb.Attribute(attributeSegmentCfg.AttributeType),
		Num:  int32(value),
	}

	return r.tempAttrs, nil
}

// ActivatePurchasedSlot 激活已购买的槽位
func (r *RefineV2) ActivatePurchasedSlot(p *Player) error {
	characterSlot, equipSlot, refinedSlot := r.parseKey(r.GetTempPurchasedSlot())
	if err := r.validateIndex(characterSlot, equipSlot, refinedSlot); err != nil {
		logger.Warnf("ActivatePurchasedSlot validateIndex error: %v", err)
		return err
	}
	r.purchasedSlots.Set(r.GetTempPurchasedSlot(), struct{}{})
	r.SetTempPurchasedSlot(0)

	p.tasks.dispatch(pb.TargetType_MAIN_ALL_REFINE_SLOTS_OPENED_TOTAL, 0, r.GetMainRoleAllSlotRefineOpenedCount())
	p.tasks.dispatch(pb.TargetType_ALL_REFINE_ENTRIES_OPENED_TOTAL, 0, r.GetAllSlotMinRefineSlot())
	return nil
}

// SetTempPurchasedSlot 设置临时购买槽位
func (r *RefineV2) SetTempPurchasedSlot(slot int32) {
	atomic.StoreInt32(&r.tempPurchasedSlots, slot)
}

// GetTempPurchasedSlot 获取临时购买槽位
func (r *RefineV2) GetTempPurchasedSlot() int32 {
	return atomic.LoadInt32(&r.tempPurchasedSlots)
}

// IsPurchased 是否购买过
func (r *RefineV2) IsPurchased(characterSlot pb.CharacterSlot, equipSlot pb.EquipSlot, refinedSlot pb.RefinedSlot) bool {
	key := r.calculateKey(characterSlot, equipSlot, refinedSlot)
	return r.purchasedSlots.Has(key)
}

// Stats 获取洗练属性 如果是主角则传入player, 如果是伙伴则传入partner
// player 是 slot 1, partner 是 slot 2，3，4，5
func (r *RefineV2) Stats(p *Player, partner *Partner, slot pb.CharacterSlot) {
	if p != nil {
		if slot != pb.CharacterSlot_CHARACTER_SLOT_1 {
			logger.Errorf("Refine Stats invalid slot for player: %v", slot)
			return
		}
		r.attrs.Range(func(key int32, value *pb.Stats) bool {
			characterSlot, _, _ := r.parseKey(key)
			if characterSlot != slot {
				return true
			}
			attr := value
			p.stats.setStats(attr.GetAttr(), attr.GetNum())
			return true
		})
	}

	if partner != nil {
		if slot == pb.CharacterSlot_CHARACTER_SLOT_1 {
			logger.Errorf("Refine Stats invalid slot for partner: %v", slot)
			return
		}
		r.attrs.Range(func(key int32, value *pb.Stats) bool {
			characterSlot, _, _ := r.parseKey(key)
			if characterSlot != slot {
				return true
			}
			attr := value
			partner.stats.setStats(attr.GetAttr(), attr.GetNum())
			return true
		})
	}
}

// GetPurchasedSlotCount 分别获取槽位3和槽位4开启的槽位数量
func (r *RefineV2) GetPurchasedSlotCount() (int, int) {
	count3 := 0
	count4 := 0
	r.purchasedSlots.Range(func(key int32, _ struct{}) bool {
		_, _, refinedSlot := r.parseKey(key)
		if refinedSlot == pb.RefinedSlot_REFINED_SLOT_3 {
			count3++
		} else if refinedSlot == pb.RefinedSlot_REFINED_SLOT_4 {
			count4++
		}
		return true
	})
	return count3, count4
}

// GetMainRoleAllSlotRefineOpenedCount 主角全身洗练槽位开启总量
func (r *RefineV2) GetMainRoleAllSlotRefineOpenedCount() int32 {
	count := int32(0)

	r.purchasedSlots.Range(func(key int32, _ struct{}) bool {
		character, _, refinedSlot := r.parseKey(key)
		if character != pb.CharacterSlot_CHARACTER_SLOT_1 {
			return true
		}

		switch refinedSlot {
		case pb.RefinedSlot_REFINED_SLOT_3:
			count++
		case pb.RefinedSlot_REFINED_SLOT_4:
			count++
		}

		return true
	})

	return count
}

// GetAllSlotMinRefineSlot 获取所有槽位的最低洗练槽位数量
func (r *RefineV2) GetAllSlotMinRefineSlot() int32 {
	countMap := make(map[pb.CharacterSlot]int32)

	r.purchasedSlots.Range(func(key int32, _ struct{}) bool {
		character, _, refinedSlot := r.parseKey(key)
		if _, exists := countMap[character]; !exists {
			countMap[character] = 0
		}
		switch refinedSlot {
		case pb.RefinedSlot_REFINED_SLOT_3:
			countMap[character]++
		case pb.RefinedSlot_REFINED_SLOT_4:
			countMap[character]++
		}
		return true
	})

	minCount := int32(0)
	for _, v := range countMap {
		if v < minCount {
			minCount = v
		}
	}

	return minCount
}
