package game

import (
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
)

func (p *PetMgr) upPetLevelShushuDate(petCfg *global.PetConfig, toPetCfg *global.PetConfig, star int32, tolevel int32, fromlevel int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.Player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	zun := false
	if toPetCfg.Mark == 1 {
		zun = true
	}
	shushuDate["is_zun"] = zun
	shushuDate["name"] = "pet_gobble"
	shushuDate["item_id"] = toPetCfg.ID
	shushuDate["item_name"] = toPetCfg.Name
	shushuDate["item_quality"] = toPetCfg.Quality
	shushuDate["item_rank"] = star

	shushuDate["eaten_id"] = petCfg.ID
	shushuDate["eaten_name"] = petCfg.Name
	shushuDate["eaten_quality"] = petCfg.Quality

	shushuDate["eaten_level"] = fromlevel
	shushuDate["level_after"] = tolevel + 1
	shushuDate["level_before"] = tolevel
	shushumain.ReportData(shushuDate)
}

func (p *PetMgr) upPetComposeShushuDate(fromPet *global.PetConfig, toPet *global.PetConfig, newPet *global.PetConfig, fromLevel int32, toLevel int32, level int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.Player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	zun := false
	if newPet.Mark == 1 {
		zun = true
	}

	shushuDate["name"] = "pet_synthesis"
	item_from_zun := false
	if fromPet.Mark == 1 {
		item_from_zun = true
	}
	item_to_zun := false
	if toPet.Mark == 1 {
		item_to_zun = true
	}
	shushuDate["cost_pet_info"] = make([]map[string]interface{}, 0, 2)
	shushuDate["cost_pet_info"] = append(shushuDate["cost_pet_info"].([]map[string]interface{}), map[string]interface{}{
		"item_id":      fromPet.ID,
		"item_name":    fromPet.Name,
		"item_quality": fromPet.Quality,
		"item_level":   fromLevel,
		"item_zun":     item_from_zun,
	})
	shushuDate["cost_pet_info"] = append(shushuDate["cost_pet_info"].([]map[string]interface{}), map[string]interface{}{
		"item_id":      toPet.ID,
		"item_name":    toPet.Name,
		"item_quality": toPet.Quality,
		"item_level":   toLevel,
		"item_zun":     item_to_zun,
	})
	shushuDate["is_zun"] = zun
	shushuDate["synthesis_item_id"] = newPet.ItemID
	shushuDate["synthesis_item_name"] = newPet.Name
	shushuDate["synthesis_item_quality"] = newPet.Quality
	shushuDate["synthesis_item_level"] = level
	shushumain.ReportData(shushuDate)
}

func (p *PetMgr) upPetOnShushuDate(fromPet *global.PetConfig, toPet *global.PetConfig, fromLevel int32, toLevel int32, fromStar int32, toStar int32, combat int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.Player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	tozun := false
	fromzun := false
	if toPet != nil && toPet.Mark == 1 {
		tozun = true
	}
	if fromPet != nil && fromPet.Mark == 1 {
		fromzun = true
	}
	shushuDate["name"] = "pet_equip"
	if toPet == nil {
		shushuDate["before_item_id"] = 0
		shushuDate["before_item_name"] = ""
		shushuDate["item_level_before"] = 0
		shushuDate["item_quality_before"] = 0
		shushuDate["before_is_zun"] = false
	} else {
		shushuDate["before_item_id"] = toPet.ID
		shushuDate["before_item_name"] = toPet.Name
		shushuDate["item_level_before"] = toLevel
		shushuDate["item_quality_before"] = toStar
		shushuDate["before_is_zun"] = tozun
	}

	shushuDate["is_zun"] = fromzun
	shushuDate["item_id"] = fromPet.ID
	shushuDate["item_name"] = fromPet.Name
	shushuDate["item_level_after"] = fromLevel
	shushuDate["item_rank"] = toStar
	shushuDate["item_quality"] = fromPet.Quality
	shushuDate["power_up"] = combat
	shushumain.ReportData(shushuDate)
}

func (p *PetMgr) upPetStarShushuDate(petCfg *global.PetConfig, CfgId int32, Num int32, oldExpStar int32, newExpStar int32, oldStar int32, newStar int32, combatDiff int32) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.Player.ToMapForShushu()
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

func (p *PetMgr) upPetAddShushuDate(petCfg *global.PetConfig, Channel string) {
	shushuDate := make(map[string]interface{}, 0)
	baseInfo := p.Player.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	zun := false
	if petCfg.Mark == 1 {
		zun = true
	}
	shushuDate["name"] = "pet_catch"
	shushuDate["item_id"] = petCfg.ID
	shushuDate["item_name"] = petCfg.Name
	shushuDate["item_quality"] = petCfg.Quality
	shushuDate["is_zun"] = zun
	shushuDate["method"] = Channel
	shushumain.ReportData(shushuDate)
}
