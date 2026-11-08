package game

import (
	"encoding/json"
	"errors"
	"math"
	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/async"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"strconv"
	"time"
)

type PetMgr struct {
	*Player

	// petNormal           int32     //普通宠物次数
	petRecoveryTime     time.Time // 最后一次宠物蛋恢复时间
	nextPetRecoveryTime int       // 下一次恢复宠物蛋时间（秒)

	petGold                 int32     // 普通宠物送的金币
	petRecoveryTimeGold     time.Time // 最后一次宠物送的金币恢复时间
	nextPetRecoveryTimeGold int       // 下一次恢复宠物送的金币时间（秒）

	petExp                 int32     // 普通宠物送的经验
	petRecoveryTimeExp     time.Time // 最后一次宠物送的经验恢复时间
	nextPetRecoveryTimeExp int       // 下一次恢复宠物送的经验时间（秒)

	petEverHave []int32                              // 曾经 已拥有的宠物列表
	pet         [pb.PetsSlot_MaxPetsSlot]PetSlotInfo // 宠物列表
	uPet        *model.UserPet
}

func NewPetMgr(p *Player) *PetMgr {
	a := &PetMgr{
		pet:    [pb.PetsSlot_MaxPetsSlot]PetSlotInfo{},
		Player: p,
	}
	return a
}

// 数据库加载
func (p *PetMgr) syncPetFromDB() error {
	player := p.Player
	uid := player.Uid()
	pet, err := db.QueryUserPet(uid)
	if err != nil {
		player.logger.Errorf("玩家宠物同步错误,userId: %d Error=%v", uid, err)
		return err
	}
	if pet == nil {
		return nil
	}
	if pet.PetInfo != "" && pet.PetInfo != "0" {
		if err = json.Unmarshal([]byte(pet.PetInfo), &p.pet); err != nil {
			player.logger.Errorf("宠物列表反序列化失败 err: %s", err.Error())
			return err
		}
	}
	p.petEverHave = []int32{}
	if pet.PetEverHave != "" && pet.PetEverHave != "0" {
		if err = json.Unmarshal([]byte(pet.PetEverHave), &p.petEverHave); err != nil {
			player.logger.Errorf("宠物曾经拥有反序列化失败 err: %s", err.Error())
			return err
		}
	}
	p.uPet = pet
	return nil
}

// 保存数据库
func (p *PetMgr) syncSaveDb() error {
	uPet := &model.UserPet{}
	player := p.Player
	uPet.Uid = player.Uid()
	if !p.petRecoveryTime.IsZero() {
		uPet.PetRecoveryTime = p.petRecoveryTime.Unix()
	}
	if !p.petRecoveryTimeGold.IsZero() {
		uPet.PetRecoveryTimeGold = p.petRecoveryTimeGold.Unix()
	}
	if !p.petRecoveryTimeExp.IsZero() {
		uPet.PetRecoveryTimeExp = p.petRecoveryTimeExp.Unix()
	}
	uPet.PetExp = p.petExp
	uPet.PetGold = p.petGold
	data, err := json.Marshal(p.pet)
	if err == nil {
		uPet.PetInfo = string(data)
	}
	data, err = json.Marshal(p.petEverHave)
	if err == nil {
		uPet.PetEverHave = string(data)
	}
	if err := db.UpdateUserPet(uPet); err != nil {
		player.logger.Errorf("保存用户宠物信息错误, Error=%v", err)
		return err
	}
	p.uPet = uPet
	return nil
}

// 更新
func (p *PetMgr) syncSavePetDb() {
	async.Run(func() {
		p.Player.lock.Lock()
		defer p.Player.lock.Unlock()

		data, err := json.Marshal(p.pet)
		if err != nil {
			return
		}
		data2, err := json.Marshal(p.petEverHave)
		if err != nil {
			return
		}
		pet := model.UserPet{
			Uid:         p.Player.uid,
			PetInfo:     string(data),
			PetEverHave: string(data2),
		}

		if err := db.UpdateUserPet(&pet); err != nil {
			p.Player.logger.Errorf("保存玩家宠物信息失败:%v", err)
			return
		}
	})
}

// 上线
func (p *PetMgr) online() {
	// var upet *model.UserPet
	// var err error
	player := p.Player
	//uid := player.Uid()
	uPet := p.uPet
	// now := time.Now().Unix()
	// cacheDuration := int64(10) // 缓存10秒

	// // 检查缓存时间
	// if val, ok := userPetCache.Load(uid); ok {
	// 	if ts, ok2 := userPetCacheTimes.Load(uid); ok2 {
	// 		if now-int64(ts.(int64)) < cacheDuration {
	// 			upet = val.(*model.UserPet)
	// 		}
	// 	}
	// }
	// if upet == nil {
	// 	upet, err = db.QueryUserPet(uid)
	// 	if err == nil && upet != nil {
	// 		userPetCache.Store(uid, upet)
	// 		userPetCacheTimes.Store(uid, now)
	// 	}
	// }
	//upet, err := db.QueryUserPet(uid)
	//if err != nil {
	//	p.logger.Errorf("玩家宠物同步错误,userId: %d Error=%v", uid, err)
	//	return err
	//}
	//if upet == nil {
	//	return nil
	//}
	//if upet.PetInfo != "" && upet.PetInfo != "0" {
	//	if err = json.Unmarshal([]byte(upet.PetInfo), &p.pet); err != nil {
	//		p.logger.Errorf("宠物列表反序列化失败 err: %s", err.Error())
	//		return err
	//	}
	//}
	//p.petEverHave = []int32{}
	//if upet.PetEverHave != "" && upet.PetEverHave != "0" {
	//	if err = json.Unmarshal([]byte(upet.PetEverHave), &p.petEverHave); err != nil {
	//		p.logger.Errorf("宠物曾经拥有反序列化失败 err: %s", err.Error())
	//		return err
	//	}
	//}
	// ...existing code...
	// 宠物恢复
	p.nextPetRecoveryTime = 0
	petNormalCount := player.getItemCountInBag(constant.ItemTypePetNormal)
	if petNormalCount >= int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) {
		p.petRecoveryTime = time.Now()
	} else {
		if uPet.PetRecoveryTime == 0 {
			// 当前时间的前一秒
			p.petRecoveryTime = time.Now().Add(-time.Second)
			p.nextPetRecoveryTime = int(global.GlobalCfg.BaseValue.Petfree.Recovertime * 60)
		}
		if uPet.PetRecoveryTime > 0 {
			// 当前时间减去最后更新时间，就是需要恢复的次数
			recoverNum := (time.Now().Unix() - uPet.PetRecoveryTime) / int64(global.GlobalCfg.BaseValue.Petfree.Recovertime*60)
			leftTime := (time.Now().Unix() - uPet.PetRecoveryTime) % int64(global.GlobalCfg.BaseValue.Petfree.Recovertime*60)
			// 次数大于最大的次数，取最大次数
			if recoverNum+petNormalCount > int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) {
				recoverNum = int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) - petNormalCount
			}
			if recoverNum < 0 {
				recoverNum = 0
			}
			// 下次恢复时间
			p.nextPetRecoveryTime = int(global.GlobalCfg.BaseValue.Petfree.Recovertime*60) - int(leftTime)
			// 次数大于0才进行恢复
			if recoverNum > 0 {
				p.addSourceEvent(constant.ItemTypePetNormal, int64(recoverNum), constant.ItemAddPetAFK)
				// 最后一次宠物蛋恢复时间
				p.petRecoveryTime = time.Now().Add(-time.Duration(leftTime) * time.Second)
			} else {
				p.petRecoveryTime = time.Unix(uPet.PetRecoveryTime, 0)
			}
		}
	}
	// ...existing code...
	bPet := p.PetMgr.GetBattlePet()
	if bPet.PetCfgId != 0 {
		petInfo := global.GetPetCfg(bPet.PetCfgId)
		if petInfo != nil && len(petInfo.OfflineAward) > 0 {
			petGoldNum := int32(0)
			petGoldGap := int32(0)
			petGoldMax := int32(0)
			petExpNum := int32(0)
			petExpGap := int32(0)
			petExpMax := int32(0)
			for _, v := range petInfo.OfflineAward {
				if v.Key == bPet.Star {
					for _, v2 := range v.Values {
						if v2.ID == constant.ItemTypeCoin {
							petGoldNum = v2.Num
							petGoldGap = int32(v2.Gap)
							petGoldMax = int32(v2.Max)
						}
						if v2.ID == constant.ItemTypeExp {
							petExpNum = v2.Num
							petExpGap = int32(v2.Gap)
							petExpMax = int32(v2.Max)
						}
					}
				}
			}
			if petGoldNum > 0 {
				if p.petGold < petGoldMax {
					if uPet.PetRecoveryTimeGold > 0 {
						p.nextPetRecoveryTimeGold = 0
						recoverNumGold := (time.Now().Unix() - uPet.PetRecoveryTimeGold) / int64(petGoldGap) * int64(petGoldNum)
						leftTime := (time.Now().Unix() - uPet.PetRecoveryTimeGold) % int64(petGoldGap)
						if p.petGold+int32(recoverNumGold) > petGoldMax {
							p.petGold = petGoldMax
							p.petRecoveryTimeGold = time.Now()
							// upet.PetRecoveryTimeGold = time.Now().Unix()
							if petGoldMax-p.petGold > 0 {
								p.addSourceEvent(constant.ItemTypeCoin, int64(petGoldMax-p.petGold), constant.ItemAddPetGold)
							}
						} else {
							p.petGold += int32(recoverNumGold)
							p.nextPetRecoveryTimeGold = int(petGoldGap) - int(leftTime)
							p.petRecoveryTimeGold = time.Now().Add(-time.Duration(petGoldGap))
							// upet.PetRecoveryTimeGold = time.Now().Unix() - int64(petGoldGap)
							if recoverNumGold > 0 {
								p.addSourceEvent(constant.ItemTypeCoin, int64(recoverNumGold), constant.ItemAddPetGold)
							}
						}
					}
				}
			}
			if petExpNum > 0 {
				if p.petExp < petExpMax {
					if uPet.PetRecoveryTimeExp > 0 {
						recoverNumExp := (time.Now().Unix() - uPet.PetRecoveryTimeExp) / int64(petExpGap) * int64(petExpNum)
						leftTime := (time.Now().Unix() - uPet.PetRecoveryTimeExp) % int64(petExpGap)
						if p.petExp+int32(recoverNumExp) > petExpMax {
							p.petExp = petExpMax
							p.petRecoveryTimeExp = time.Now()
							// upet.PetRecoveryTimeExp = time.Now().Unix()
							if petExpMax-p.petExp > 0 {
								p.addSourceEvent(constant.ItemTypeExp, int64(petExpMax-p.petExp), constant.ItemAddPetExp)
							}
						} else {
							p.petExp += int32(recoverNumExp)
							p.nextPetRecoveryTimeExp = int(petExpGap) - int(leftTime)
							p.petRecoveryTimeExp = time.Now().Add(-time.Duration(petExpGap))
							// upet.PetRecoveryTimeExp = time.Now().Unix() - int64(petExpGap)
							if recoverNumExp > 0 {
								p.addSourceEvent(constant.ItemTypeExp, int64(recoverNumExp), constant.ItemAddPetExp)
							}
						}
					}
				}
			}

		}
	}
}

func (p *PetMgr) SyncUserPet(needDel int32) {
	packet := &pb.SyncUserPetInfoPush{Pet: []*pb.PetSlotInfo{}, PetRecoveryTime: 0, PetEverHave: []int32{}}
	packet.Pet = []*pb.PetSlotInfo{}
	// 赋值
	for i := range p.pet {
		v := p.pet[i]
		packet.Pet = append(packet.Pet, &pb.PetSlotInfo{
			Uid:     v.PetUid,
			CfgId:   v.PetCfgId,
			Star:    v.Star,
			Level:   v.Level,
			Exp:     v.Exp,
			ExpStar: v.ExpStar,
		})
	}
	checkNum := int64(0)
	// 当前时间+下一次恢复时间
	checkNum = int64(global.GlobalCfg.BaseValue.Petfree.Maxnum) + int64(needDel)
	if p.Player.checkTakeSource(constant.ItemTypePetNormal, checkNum) {
		packet.PetRecoveryTime = 0
	} else {
		// 这里是给前端的时间，就是当前时间加上需要等待的时间
		packet.PetRecoveryTime = int32(time.Now().Unix() + int64(p.nextPetRecoveryTime))
	}
	logger.Infof("SyncUserPetInfoPush PetRecoveryTime:%d, nextPetRecoveryTime:%d", packet.PetRecoveryTime, p.nextPetRecoveryTime)
	packet.PetEverHave = p.petEverHave
	p.Player.pushMsg(protocol.RouteSyncUserPetInfoPush, packet)
}

func (p *PetMgr) newDay() {
	p.petGold = 0
	p.petExp = 0
}

func (p *PetMgr) update() {
	// 宠物恢复检测
	if p.nextPetRecoveryTime > 0 {
		p.nextPetRecoveryTime--
	}
	player := p.Player
	// 宠物恢复速度加成,换算成秒//体力恢复速度加成,换算成秒
	pettimeinterval := int(global.GlobalCfg.BaseValue.Petfree.Recovertime * 60)
	// 查询背包里面的宠物蛋数目,如果大于设置的数目就不回复了
	if player.checkTakeSource(constant.ItemTypePetNormal, int64(global.GlobalCfg.BaseValue.Petfree.Maxnum)) {
		p.nextPetRecoveryTime = pettimeinterval
		p.petRecoveryTime = time.Now()
	} else {
		if time.Now().After(p.petRecoveryTime.Add(time.Duration(pettimeinterval) * time.Second)) {
			p.addSourceEvent(constant.ItemTypePetNormal, 1, constant.ItemAddPetAuth)
			p.nextPetRecoveryTime = pettimeinterval
			p.petRecoveryTime = time.Now()
			p.SyncUserPet(-1)
		}
	}

	bPet := p.PetMgr.GetBattlePet()
	if bPet != nil {
		if bPet.PetCfgId != 0 {
			petInfo := global.GetPetCfg(bPet.PetCfgId)
			if petInfo != nil && len(petInfo.OfflineAward) > 0 {
				petGoldNum := int32(0)
				petGoldGap := int32(0)
				petGoldMax := int32(0)
				petExpNum := int32(0)
				petExpGap := int32(0)
				petExpMax := int32(0)
				for _, v := range petInfo.OfflineAward {
					if v.Key == bPet.Star {
						for _, v2 := range v.Values {
							if v2.ID == constant.ItemTypeCoin {
								petGoldNum = v2.Num
								petGoldGap = int32(v2.Gap)
								petGoldMax = int32(v2.Max)
							}
							if v2.ID == constant.ItemTypeExp {
								petExpNum = v2.Num
								petExpGap = int32(v2.Gap)
								petExpMax = int32(v2.Max)
							}
						}
					}
				}
				// log.Debug("", p.petGold, "-----------", p.petRecoveryTimeGold, p.petExp, "-----------", p.petRecoveryTimeExp)
				// log.Debug("", p.pet[0].Star, p.petExp, "-----------", petGoldNum, "---", petGoldGap, "---", petGoldMax, "---", petExpNum, "---", petExpGap, "---", petExpMax)
				if petGoldNum > 0 {
					// 宠物金币恢复检测
					if p.nextPetRecoveryTimeGold > 0 {
						p.nextPetRecoveryTimeGold--
					}
					// 宠物金币速度加成,换算成秒
					timeintervalgold := int(float64(petGoldGap))
					if p.petGold < petGoldMax {
						if time.Now().After(p.petRecoveryTimeGold.Add(time.Duration(timeintervalgold) * time.Second)) {
							if p.petGold+petGoldNum > petGoldMax {
								p.petGold = petGoldMax
								p.nextPetRecoveryTimeGold = 0
								if petGoldMax-p.petGold > 0 {
									p.addSourceEvent(constant.ItemTypeCoin, int64(petGoldMax-p.petGold), constant.ItemAddPetGold)
								}
							} else {
								p.petGold += petGoldNum
								p.nextPetRecoveryTimeGold = timeintervalgold
								if petGoldNum > 0 {
									p.addSourceEvent(constant.ItemTypeCoin, int64(petGoldNum), constant.ItemAddPetGold)
								}
							}
							p.petRecoveryTimeGold = time.Now()
							player.syncUserBase()
						}
					}
				}

				if petExpNum > 0 {
					// 宠物经验恢复检测
					if p.nextPetRecoveryTimeExp > 0 {
						p.nextPetRecoveryTimeExp--
					}
					// 宠物经验速度加成,换算成秒
					timeintervalexp := int(float64(petExpGap))
					if p.petExp < petExpMax {
						if time.Now().After(p.petRecoveryTimeExp.Add(time.Duration(timeintervalexp) * time.Second)) {
							if p.petExp+petExpNum > petExpMax {
								p.petExp = petExpMax
								p.nextPetRecoveryTimeExp = 0
								if petExpMax-p.petExp > 0 {
									p.addSourceEvent(constant.ItemTypeExp, int64(petExpMax-p.petExp), constant.ItemAddPetExp)
								}
							} else {
								p.petExp += petExpNum
								p.nextPetRecoveryTimeExp = timeintervalexp
								if petExpNum > 0 {
									p.addSourceEvent(constant.ItemTypeExp, int64(petExpNum), constant.ItemAddPetExp)
								}
							}
							p.petRecoveryTimeExp = time.Now()
							player.syncUserBase()
						}
					}
				}

			}

		}
	}

}

func (p *PetMgr) getTalent(whileBuff *[]int32) {
	pet := p.pet[0]
	if pet.PetCfgId != 0 {
		// 查询宠物buff
		petInfo := global.GetPetCfg(pet.PetCfgId)
		for _, v := range petInfo.Talent {
			*whileBuff = append(*whileBuff, v.ID)
		}
	}
}

// 获取宠物天赋Buff列表
func (p *PetMgr) getPetBuffs() []*PetBuff {
	petBuffs := []*PetBuff{}
	bPet := p.PetMgr.GetBattlePet()
	if bPet.PetCfgId != 0 {
		if petCfg := global.GetPetCfg(bPet.PetCfgId); petCfg != nil {
			for _, v := range petCfg.Talent {
				bufStats := &Stats{Attrs: make(map[pb.Attribute]int32)}
				if len(v.Value) > 0 {
					for _, vv := range v.Value {
						bufStats.setStats(pb.Attribute(vv.ID), vv.Num)
					}
				}
				petBuffs = append(petBuffs, &PetBuff{Id: v.ID, Stats: bufStats, Rate: v.Rate})
			}
		}
	}
	return petBuffs
}

func (p *PetMgr) setStats(petStatus *Stats) {
	// 宠物
	// 宠物因为有些属性加给装备，所以提到前面来了
	bPet := p.GetBattlePet()
	if bPet.PetCfgId != 0 {
		petCfg := global.GetPetCfg(bPet.PetCfgId)
		if petCfg != nil {
			for _, e2 := range petCfg.StarAttribute {
				if e2.Id == bPet.Star {
					for _, e3 := range e2.Values {
						petStatus.setStats(pb.Attribute(e3.Key), e3.Value)
					}
					break
				}
			}

			roleCfg, _ := global.GetRoleCfg(constant.Pet, bPet.Level)
			if roleCfg != nil {
				rate := global.GetpetQualityRate(bPet.Quality)
				if rate == 0 || petCfg.PetKind.Key == 2 {
					rate = 1
				}
				for _, e := range roleCfg.BassProperty {
					petStatus.setStats(pb.Attribute(e.Key), int32(math.Round(float64(e.Value)*rate)))
				}
			}
		}
	}
}

// 增加宠物
func (p *PetMgr) addPet(cfgId int32, exp int32, star int32, toSlotId int32, NeedDel int32, Channel string) (newSort int32, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("addPet Panic Error=%v  Stack=%v", err, stack())
		}
	}()

	if cfgId == 0 {
		return 0, nil
	}
	endPetSlot := toSlotId

	player := p.Player
	// par := CreatePet(cfgId)
	// par.owner = p
	petSlot := p.checkEmptyPetSlot()
	if petSlot == -1 {
		return 0, errutil.ErrPetFull
	}
	// 判断这个等级能开多少槽位
	BagInfo := global.GlobalCfg.Petbag
	canhave := int32(0)
	for _, v := range BagInfo {
		// 判断当前等级能开多少槽位
		if player.level >= int32(v.Level) {
			canhave = int32(v.Value)
		} else {
			continue
		}
	}
	// 假设他现在能开10个，给到key是10，也是可以放进去的，因为0是上阵的
	if petSlot > canhave {
		return 0, errutil.ErrPetFull
	}
	// cfgid对应的宠物信息赋值给petslot
	petCfg := global.GetPetCfg(cfgId)
	if petCfg == nil {
		return 0, errutil.ErrConfigNotFound
	}
	petExpForLevel := int32(1)
	newExp := int32(0)
	if exp > 0 {
		petExpForLevel, newExp, err = p.petExpForLevel(1, int32(exp), petCfg.Quality)
		if err != nil {
			return 0, err
		}
	}
	if toSlotId == 0 {
		endPetSlot = petSlot
	}
	p.pet[endPetSlot] = PetSlotInfo{
		PetCfgId: cfgId,
		Quality:  petCfg.Quality,
		Level:    petExpForLevel,
		Exp:      newExp,
		Star:     star,
		ItemId:   petCfg.ItemID,
		ExpStar:  0,
	}
	player.tasks.dispatch(pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL, 0, petExpForLevel)
	// 判断p.petEverHave有没有，没有的话，就添加新的
	if p.petEverHave == nil {
		p.petEverHave = []int32{cfgId}
	} else {
		nochange := 0
		for _, v := range p.petEverHave {
			if v == cfgId {
				nochange++
				continue
			}
		}
		if nochange == 0 {
			p.petEverHave = append(p.petEverHave, cfgId)
		}
	}
	p.syncSavePetDb()
	p.SyncUserPet(NeedDel)
	if Channel != "合成" {
		p.upPetAddShushuDate(petCfg, Channel)
	}
	return petSlot, nil
}

// 校验是否有空的宠物槽位 有的话返回具体的槽位
func (p *PetMgr) checkEmptyPetSlot() int32 {
	// 判断是否有空的槽位
	for k, v := range p.pet {
		if v.PetCfgId == 0 {
			if k == 0 {
				continue
			}
			return int32(k)
		}
	}
	return -1
}

// 获取空余槽位数量
func (p *PetMgr) getEmptyPetSlotNum() int {
	num := 0
	for k, v := range p.pet {
		if k == 0 {
			continue
		}

		if v.PetCfgId == 0 {
			num++
		}
	}
	return num
}

func (p *PetMgr) addSourceEvent(id int32, num int64, reason string) {
	source := protocol.AddSource{Source: []*protocol.SourceInfo{}, Method: reason}
	source.Source = append(source.Source, &protocol.SourceInfo{Id: id, Num: num})
	p.Player.addSourceEvent(&source)
}

// 升级查询这个经验对应多少等级
func (p *PetMgr) petExpForLevel(level int32, exp int32, quality int32) (newLevel int32, newExp int32, err error) {
	for i := 0; i < len(global.RoleCfg); i++ {
		RoleCfg, err := global.GetRoleCfg(constant.Pet, level)
		if err != nil {
			return 0, 0, err
		}
		// 每个伙伴品质对应不同的系数
		per := global.GetpetLevelRate(quality)
		if float64(exp) >= math.Round(float64(RoleCfg.NeedExp)*per) {
			exp -= int32(math.Round(float64(RoleCfg.NeedExp) * per))
			level++
		} else {
			return level, exp, nil
		}
	}
	return level, exp, nil
}

// 查询这个经验对应多少经验
func (p *PetMgr) petExpForExp(level int32, exp int32, quality int32) (allExp int32, err error) {
	allExp = int32(0)
	for i := 1; i < len(global.RoleCfg); i++ {
		RoleCfg, err := global.GetRoleCfg(constant.Pet, int32(i))
		if err != nil {
			return 0, err
		}
		// 这级需要的经验，如果没超过这级，不能加到经验里面
		// 每个伙伴品质对应不同的系数
		per := global.GetpetLevelRate(quality)
		if int(level) > i {
			exp += int32(math.Round(float64(RoleCfg.NeedExp) * per))
		} else {
			return allExp + exp, nil
		}
	}
	return allExp + exp, nil
}

// 宠物信息
func (p *PetMgr) getPetInfo(SlotId int32) *PetSlotInfo {
	for i, pet := range p.pet {
		if i == int(SlotId) {
			return &pet
		}
	}
	return nil
}

// 宠物升星
func (p *PetMgr) petStar(slotId int32, cfgId int32, Num int32) (resNum int32, err error) {
	resper := int32(0)
	EndExp := int64(p.pet[slotId].ExpStar)
	needLevel := int32(0)
	pet := p.getPetInfo(slotId)
	if pet == nil {
		return resper, errutil.ErrPetInvalid
	}

	petCfg := global.GetPetCfg(pet.PetCfgId)
	if petCfg == nil {
		return resper, errutil.ErrConfigNotFound
	}

	if pet.Star >= int32(len(petCfg.StarConsume)) {
		return resper, errors.New("已满阶")
	}
	// 删除材料
	if !p.Player.checkTakeSource(cfgId, int64(Num)) {
		return resper, errutil.ErrItemNotEnough
	}

	oldExpStar := p.pet[slotId].ExpStar
	newExpStar := oldExpStar
	oldStar := p.pet[slotId].Star
	// 查询宠物当前的星级所需要的材料
	for _, v := range petCfg.StarConsume {
		// 查询这个道具对应多少经验
		// perAll := make([]int, 0)
		for i := 0; i <= int(Num); i++ {
			if v.Id-1 == pet.Star {
				PetAdvanceInfo := global.GlobalCfg.PetAdvance.Basevalue.Item
				slot := int64(0)
				for _, v2 := range PetAdvanceInfo {
					if v2.ID == cfgId {
						slot = int64(v2.Num)
						break
					}
				}
				per := petStageRand()
				resper = int32(per)
				if per == -1 {
					return 0, errutil.ErrConfigNotFound
				}
				// perAll = append(perAll, per)
				EndExp = int64(p.pet[slotId].ExpStar) + slot*int64(per)
				// 查询这个道具对应多少经验
				if EndExp >= int64(v.Values) {
					EndExp -= int64(v.Values)
					for _, v3 := range petCfg.StarLimit {
						if v3.Id == pet.Star+1 {
							needLevel = v3.Values
							break
						}
					}
					// 如果当前宠物等级不够，不能升星
					if pet.Level < needLevel {
						// 如果已经满经验了，就不扣除材料
						if p.pet[slotId].ExpStar != v.Values {
							EndExp = int64(v.Values)
						} else {
							return 0, errutil.ErrPetLevelNotEnough // 宠物等级不足
						}
						break
					} else {
						p.pet[slotId].Star++
						break
					}
				} else {
					break
				}
				// pet.ExpStar+
			}
		}
	}
	newStar := p.pet[slotId].Star
	p.Player.takeSource(cfgId, int64(Num), constant.ItemDelPetStar)
	p.pet[slotId].ExpStar = int32(EndExp)
	newExpStar = p.pet[slotId].ExpStar
	p.syncSavePetDb()
	p.SyncUserPet(0)
	oldcombat := p.Player.totalCombat
	p.Player.refreshStats()
	newcombat := p.Player.totalCombat
	// 升星数据上报
	p.upPetStarShushuDate(petCfg, cfgId, Num, oldExpStar, newExpStar, oldStar, newStar, newcombat-oldcombat)
	return resper, nil
}

func (p *PetMgr) petLevel(SlotId int32) error {
	pet := p.getPetInfo(SlotId)
	if pet == nil {
		return errutil.ErrPetInvalid
	}

	petCfg := global.GetPetCfg(pet.PetCfgId)
	if petCfg == nil {
		return errutil.ErrConfigNotFound
	}
	roleCfg, _ := global.GetRoleCfg(constant.Pet, pet.Level)
	if roleCfg == nil {
		return errutil.ErrConfigNotFound
	}
	if pet.Star >= int32(len(petCfg.StarConsume)) {
		return errutil.ErrMaxLevel
	}
	// 查询宠物当前的等级所需要的经验
	per := global.GetpetLevelRate(petCfg.Quality)
	NeedExp := math.Round(float64(roleCfg.NeedExp)*per) - float64(pet.Exp)
	if !p.Player.checkTakeSource(constant.ItemTypePetExp, int64(NeedExp)) {
		return errutil.ErrItemNotEnough
	}
	p.Player.takeSource(constant.ItemTypePetExp, int64(NeedExp), constant.ItemDelPetLevel)
	p.pet[SlotId].Level++
	p.pet[SlotId].Exp = 0
	p.syncSavePetDb()
	p.SyncUserPet(0)
	p.Player.refreshStats()
	return nil
}

func (p *PetMgr) petAddExp(slotId int32, exp int32) (bool, error) {
	up := false
	pet := p.getPetInfo(slotId)
	if pet == nil {
		return up, errutil.ErrPetInvalid
	}
	petCfg := global.GetPetCfg(p.pet[slotId].PetCfgId)
	if petCfg == nil {
		return up, errutil.ErrConfigNotFound
	}

	willexp := int64(exp)
	oldLevel := p.pet[slotId].Level
	petExpForLevel, newExp, err := p.petExpForLevel(p.pet[slotId].Level, int32(willexp), petCfg.Quality)
	if err != nil {
		return up, err
	}
	p.Player.tasks.dispatch(pb.TargetType_ANY_PET_UPGRADED_TO_LEVEL, 0, petExpForLevel)
	p.pet[slotId].Level = petExpForLevel
	p.pet[slotId].Exp = newExp
	// pet.Star++
	if oldLevel != p.pet[slotId].Level {
		up = true
	}
	p.syncSavePetDb()
	p.SyncUserPet(0)
	return up, nil
}

// 宠物信息
func (p *PetMgr) petMove(fromSlotId, toSlotId int32) (newSlot int32, newPet int32, newPetCfgId int32, err error) {
	if fromSlotId < 0 || fromSlotId >= int32(len(p.pet)) || toSlotId < 0 || toSlotId >= int32(len(p.pet)) {
		return 0, 0, 0, errutil.ErrPetInvalid
	}
	// 会有好几种情况，第一种是交换两个空的宠物，第二种是交换两个有宠物的宠物，第三种是交换一个空的宠物和一个有宠物的宠物
	if p.pet[fromSlotId].PetCfgId == 0 {
		return 0, 0, 0, errutil.ErrPetInvalid
	}
	newSlot = toSlotId
	var petCfg *global.PetConfig

	// 如果拖动的宠物是经验宝宝，那么直接给目的槽位宠物加经验
	if fromSlotId != 0 {
		petCfg = global.GetPetCfg(p.pet[fromSlotId].PetCfgId)
	}
	var toPetCfg *global.PetConfig
	if toSlotId != 0 {
		toPetCfg = global.GetPetCfg(p.pet[toSlotId].PetCfgId)
	}

	toLevel := p.pet[toSlotId].Level
	fromLevel := p.pet[fromSlotId].Level
	fromStar := p.pet[fromSlotId].Star
	toStar := p.pet[toSlotId].Star
	// 两个位置都有人，如果他们品级都一样，直接合成
	if p.pet[fromSlotId].PetCfgId != 0 && p.pet[toSlotId].PetCfgId != 0 {
		if fromSlotId != 0 && toSlotId != 0 {
			if petCfg.PetKind.Key == 2 {
				if p.pet[toSlotId].PetCfgId != 0 {
					// 是经验宝宝，有两个经验，一个是经验宝宝累积的经验，还有就是经验宝宝自身配表的经验
					fromExp, err := p.petExpForExp(p.pet[fromSlotId].Level, p.pet[fromSlotId].Exp, p.pet[fromSlotId].Quality)
					if err != nil {
						return newSlot, 0, 0, err
					}
					allExp := petCfg.PetKind.Value + p.pet[toSlotId].Exp + fromExp
					up, err := p.petAddExp(toSlotId, allExp)
					if err != nil {
						return newSlot, 0, 0, err
					}
					if up {
						p.upPetLevelShushuDate(petCfg, toPetCfg, p.pet[toSlotId].Star, toLevel, fromLevel)
					}
					p.pet[fromSlotId] = PetSlotInfo{}
					p.syncSavePetDb()
					p.SyncUserPet(0)
					return newSlot, 0, 0, nil
				} else {
					p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
					p.syncSavePetDb()
					p.SyncUserPet(0)
					return newSlot, 0, 0, nil
				}
			}
			if p.pet[fromSlotId].Quality == p.pet[toSlotId].Quality {
				// 合成
				whileAll := make(map[int32][]int32)
				for _, v := range global.ComposeCfg {
					for _, v2 := range v.Raw {
						whileAll[v2] = append(whileAll[v2], v.ComposeID)
					}
				}
				newId := int32(0)
				newPetId := int32(0)
				if len(whileAll[p.pet[fromSlotId].ItemId]) == 0 || len(whileAll[p.pet[toSlotId].ItemId]) == 0 {
					p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
					p.syncSavePetDb()
					p.SyncUserPet(0)
					return newSlot, 0, 0, errutil.ErrComposeError
				}

				// 两个数组有交集，说明就是可以合成的
				if len(whileAll[p.pet[fromSlotId].ItemId]) > 0 && len(whileAll[p.pet[toSlotId].ItemId]) > 0 {
					composeIds := make([]int32, 0)
					for _, v := range whileAll[p.pet[fromSlotId].ItemId] {
						for _, v2 := range whileAll[p.pet[toSlotId].ItemId] {
							if v == v2 {
								composeIds = append(composeIds, v)
								break
							}
						}
					}
					if len(composeIds) == 0 {
						p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
						p.syncSavePetDb()
						p.SyncUserPet(0)
						return newSlot, 0, 0, errutil.ErrComposeError
					}
					// 查找子集，优先读子集
					// 查找最短的Raw数组，把这个数组的第一个元素作为合成id
					composeId := composeIds[0]
					for _, v := range composeIds {
						if len(global.ComposeCfg[v].Raw) < len(global.ComposeCfg[composeId].Raw) {
							composeId = v
						}
					}
					logger.Debugf("合成id %d", composeId)
					// 合成的id
					if global.ComposeCfg[composeId] == nil {
						p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
						p.syncSavePetDb()
						p.SyncUserPet(0)
						return newSlot, 0, 0, errutil.ErrComposeError
					}
					newId = composeRand(global.ComposeCfg[composeId])
					// 转换成宠物id
					for _, v := range global.PetCfg {
						if v.ItemID == newId {
							newPetId = v.ID
							break
						}
					}

					// 两个宠物的经验相加
					fromExp, err := p.petExpForExp(p.pet[fromSlotId].Level, p.pet[fromSlotId].Exp, p.pet[fromSlotId].Quality)
					if err != nil {
						return newSlot, 0, 0, err
					}
					toExp, err := p.petExpForExp(p.pet[toSlotId].Level, p.pet[toSlotId].Exp, p.pet[toSlotId].Quality)
					if err != nil {
						return newSlot, 0, 0, err
					}
					allExp := fromExp + toExp

					// 星级取最大星级
					star := int32(0)
					if p.pet[fromSlotId].Star > p.pet[toSlotId].Star {
						star = p.pet[fromSlotId].Star
					} else {
						star = p.pet[toSlotId].Star
					}
					fromPetLevel := p.pet[fromSlotId].Level
					toPetLevel := p.pet[toSlotId].Level
					// 去掉原来的宠物
					p.pet[fromSlotId] = PetSlotInfo{}
					p.pet[toSlotId] = PetSlotInfo{}
					newSlot, _ = p.addPet(newPetId, allExp, star, toSlotId, 0, "合成")
					p.syncSavePetDb()
					// 查询这个宠物的品质
					petInfo := global.GetPetCfg(newPetId)
					if petInfo.Quality >= 6 && petInfo.PetKind.Key == 1 {
						strSend := []string{p.Player.name, strconv.Itoa(int(petInfo.ItemID))}
						defaultManager.Lamp(int32(constant.LampPetComposeSummon), p.Player.serverID, strSend)
					}
					p.addUpRecord(constant.PetSynthesisNum)
					p.upPetComposeShushuDate(petCfg, toPetCfg, petInfo, fromPetLevel, toPetLevel, p.pet[newSlot].Level)
					p.rankChangeTCP(constant.RankTypePetSynthesis, p.getUpRecord(constant.PetSynthesisNum))
					p.Player.tasks.dispatch(pb.TargetType_PET_COMPOSE_TIMES, 0, 1)

					switch constant.Quality(petInfo.Quality) {
					case constant.QualityBlue:
						switch petInfo.Mark {
						case 1: // 代表“尊”
							p.tasks.dispatch(pb.TargetType_COMPOSE_BLUE_ELITE_PETS, 0, 1)
						case 2: // 代表“凡”
							p.tasks.dispatch(pb.TargetType_COMPOSE_BLUE_COMMON_PETS, 0, 1)
						}
					case constant.QualityPurple:
						switch petInfo.Mark {
						case 1:
							p.tasks.dispatch(pb.TargetType_COMPOSE_PURPLE_ELITE_PETS, 0, 1)
						case 2:
							p.tasks.dispatch(pb.TargetType_COMPOSE_PURPLE_COMMON_PETS, 0, 1)
						}
					case constant.QualityOrange:
						switch petInfo.Mark {
						case 1:
							p.tasks.dispatch(pb.TargetType_COMPOSE_ORANGE_ELITE_PETS, 0, 1)
						case 2:
							p.tasks.dispatch(pb.TargetType_COMPOSE_ORANGE_COMMON_PETS, 0, 1)
						}
					case constant.QualityRed:
						switch petInfo.Mark {
						case 1:
							p.tasks.dispatch(pb.TargetType_COMPOSE_RED_ELITE_PETS, 0, 1)
						case 2:
							p.tasks.dispatch(pb.TargetType_COMPOSE_RED_COMMON_PETS, 0, 1)
						}
					case constant.QualityYellow:
						switch petInfo.Mark {
						case 1:
							p.tasks.dispatch(pb.TargetType_COMPOSE_YELLOW_ELITE_PETS, 0, 1)
							// case 2:
							// p.tasks.dispatch(pb.TargetType_COMPOSE_YELLOW_COMMON_PETS, 0, 1)
						}
					}
					// p.SyncUserPet(0)
					if len(global.ComposeCfg[composeId].Return) > 0 {
						for _, cv := range global.ComposeCfg[composeId].Return {
							if cv.Key == newId {
								for _, cv2 := range cv.Values {
									// 给玩家添加材料
									logger.Debugf("给玩家添加材料 %d %d", cv2.ID, cv2.Num)
									p.addSourceEvent(cv2.ID, int64(cv2.Num), constant.ItemAddCompose)
								}
							}
						}
					}
					return newSlot, newPetId, newPetId, nil
				} else {
				}
			}
		}
	}
	p.pet[toSlotId], p.pet[fromSlotId] = p.pet[fromSlotId], p.pet[toSlotId]
	p.syncSavePetDb()
	p.SyncUserPet(0)
	if fromSlotId == 0 || toSlotId == 0 {
		oldCombat := p.Player.totalCombat
		p.Player.defaultBuff()
		newCombat := p.Player.totalCombat
		if toSlotId == 0 {
			p.upPetOnShushuDate(petCfg, toPetCfg, fromLevel, toLevel, fromStar, toStar, newCombat-oldCombat)
		}
	}
	return -1, 0, 0, nil
}

// 获取出战宠物信息
func (p *PetMgr) GetBattlePetInfo() *pb.PetSlotInfo {
	battlePet := p.pet[0]
	pet := &pb.PetSlotInfo{
		CfgId:   battlePet.PetCfgId,
		Level:   battlePet.Level,
		Star:    battlePet.Star,
		Enhance: battlePet.Enhance,
		Quality: battlePet.Quality,
		Exp:     battlePet.Exp,
		ExpStar: battlePet.ExpStar,
	}
	return pet
}

// 获取出战宠物
func (p *PetMgr) GetBattlePet() *PetSlotInfo {
	return &p.pet[0]
}
