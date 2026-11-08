package game

import (
	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/async"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/utils"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
)

type Artifact struct {
	*Player
	artifacts *utils.SafeMap[int32, *pb.ArtifactInfo]
}

func NewArtifact(p *Player) *Artifact {
	a := &Artifact{
		artifacts: utils.NewSafeMap[int32, *pb.ArtifactInfo](),
		Player:    p,
	}
	return a
}

func (a *Artifact) saveArtifactDB() error {
	info := &model.Artifact{
		UID:       a.uid,
		Artifacts: make(map[int32]*pb.ArtifactInfo),
	}

	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		info.Artifacts[key] = value
		return true
	})

	if err := db.UpdateArtifact(info); err != nil {
		return nil
	} else {
		logger.Debugf("保存神兵信息成功Uid %d", a.uid)
		return err
	}
}

func (a *Artifact) loadArtifactDB() error {
	if u, err := db.GetArtifactByUid(a.uid); err == nil {
		if u != nil {
			for _, v := range u.Artifacts {
				a.artifacts.Set(v.CfgId, v)
			}
		} else {
			a.initArtifact()
		}
		return nil
	} else {
		logger.Errorf("加载神兵信息失败Uid %d", a.uid)
		return err
	}
}

func (a *Artifact) syncArtifact() {
	packet := &pb.ArtifactListPush{}
	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		packet.Artifacts = append(packet.Artifacts, value)
		return true
	})
	a.pushMsg(protocol.RouteArtifactListPush, packet)
}

func (a *Artifact) initArtifact() {
	if !a.isModuleOpen(constant.ModuleArtifact) {
		return
	}

	if a.artifacts.Len() > 0 {
		a.artifactSkillCheck()
		return
	}

	for _, v := range global.ArtifactAdvanceCfg {
		if a.level < v.OpenLimit {
			continue
		}

		if v.Rank != 0 {
			continue
		}

		if _, ok := a.artifacts.Get(v.UID); !ok {
			a.addArtifact(v.UID)
		}
	}

	info := &model.Artifact{
		UID:       a.uid,
		Artifacts: make(map[int32]*pb.ArtifactInfo),
	}

	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		info.Artifacts[key] = value
		return true
	})

	async.Run(func() {
		if err := db.InsertArtifact(info); err != nil {
			a.logger.Errorf("初始化神兵信息失败Uid %d", a.uid)
		}
	})
}

// 技能解锁
func (a *Artifact) artifactSkillCheck() {
	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		cfg, ok := global.ArtifactAdvanceCfg[key]
		if !ok {
			return true
		}

		if skill := a.getSkill(cfg.SkillID); skill == nil {
			a.addSkill(cfg.SkillID)
			a.logger.Debugf("神兵初始化技能遗漏 %d", cfg.SkillID)
		}
		return true
	})
}

func (a *Artifact) artifactUnlock() {
	if !a.isModuleOpen(constant.ModuleArtifact) {
		return
	}

	for _, v := range global.ArtifactAdvanceCfg {
		if a.level < v.OpenLimit {
			continue
		}

		if v.Rank != 0 {
			continue
		}

		flag := true
		a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
			cfg, ok := global.ArtifactAdvanceCfg[key]
			if !ok {
				return true
			}

			if v.WeaponID == cfg.WeaponID {
				flag = false
				return false
			}
			return true
		})

		if a.getArtifact(v.UID) == nil && flag {
			a.addArtifact(v.UID)
			a.RefreshStats()

			a.logger.Debugf("%s解锁新的神兵 ID :  %d", a.name, v.UID)
		}
	}
}

func (a *Artifact) getArtifact(cfgId int32) *pb.ArtifactInfo {
	if w, ok := a.artifacts.Get(cfgId); ok {
		return w
	}
	return nil
}

func (a *Artifact) addArtifact(cfgId int32) *pb.ArtifactInfo {
	info := &pb.ArtifactInfo{
		CfgId:   cfgId,
		Soul:    make([]*pb.ArtifactSoul, 0),
		IsEquip: false,
	}

	if cfg, ok := global.ArtifactAdvanceCfg[cfgId]; ok {
		for _, v := range cfg.AdvanceLimit {
			info.Soul = append(info.Soul, &pb.ArtifactSoul{
				CfgId: v.ID,
				Level: 1,
				Exp:   0})
		}
	}

	a.artifacts.Set(cfgId, info)

	if a.artifacts.Len() <= global.GlobalCfg.BaseValue.Epicweapon.MaxNum {
		info.IsEquip = true
	}

	packet := &pb.ArtifactUpdatePush{
		Artifact: info,
		Type:     pb.ArtifactUpdatePush_NEW,
	}

	a.artifactSkill(cfgId, true)

	a.pushMsg(protocol.RouteArtifactUpdatePush, packet)
	return info
}

func (a *Artifact) updateArtifact(cfgId int32) {
	if info := a.getArtifact(cfgId); info != nil {
		packet := &pb.ArtifactUpdatePush{
			Artifact: info,
			Type:     pb.ArtifactUpdatePush_UPDATE,
		}
		a.pushMsg(protocol.RouteArtifactUpdatePush, packet)
	}
}

// 神兵强化
func (a *Artifact) artifactStrength(cfgId, soulCfgId int32) (int64, error) {

	var soul *pb.ArtifactSoul
	artifact := a.getArtifact(cfgId)
	if artifact == nil {
		return 0, errutil.ErrInvalidArtifact
	}
	for _, v := range artifact.Soul {
		if v.CfgId == soulCfgId {
			soul = v
			break
		}
	}
	if soul == nil {
		return 0, errutil.ErrInvalidArtifact
	}

	cfg, ok := global.ArtifactStrengthCfg[soul.CfgId*1000+soul.Level]
	if !ok {
		return 0, errutil.ErrConfigNotFound
	}

	if _, ok := global.ArtifactStrengthCfg[soul.CfgId*1000+soul.Level+1]; !ok {
		return 0, errutil.ErrMaxLevel
	}

	if artifact.CfgId%100 < cfg.LevelLimit {
		return 0, errutil.ErrMaxLevel
	}

	var itemId int32
	var addExp int64

	if len(global.GlobalCfg.EpicWeapon.Basevalue.Item) == 0 {
		return 0, errutil.ErrConfigNotFound
	}
	itemId = global.GlobalCfg.EpicWeapon.Basevalue.Item[0].ID
	addExp = global.GlobalCfg.EpicWeapon.Basevalue.Item[0].Num

	// 扣除消耗
	if a.getItemCountInBag(itemId) == 0 {
		return 0, errutil.ErrItemNotEnough
	}
	a.takeSource(itemId, 1, constant.ItemDelArtifact)

	// 计算经验加成
	var baseRate float64
	for _, v := range global.GlobalCfg.EpicWeapon.Rate.Stage {
		baseRate += float64(v.Pro)
	}
	weight := global.RandIn(int(baseRate))
	var tmpWeight int
	for _, v := range global.GlobalCfg.EpicWeapon.Rate.Stage {
		tmpWeight += v.Pro
		if weight <= tmpWeight {
			addExp = addExp * int64(v.Rate)
			break
		}
	}

	soul.Exp += addExp
	var isUpgrade bool
	for soul.Exp >= cfg.LevelExp {
		soul.Exp -= cfg.LevelExp
		soul.Level++
		isUpgrade = true
	}

	a.updateArtifact(cfgId)

	if isUpgrade {
		a.RefreshStats()
	}
	a.tasks.dispatch(pb.TargetType_ARTIFACT_UPGRADE_TIMES, 0, 1)
	return addExp, nil
}

// 进阶
func (a *Artifact) artifactAdvance(cfgId int32) error {
	artifact := a.getArtifact(cfgId)
	if artifact == nil {
		return errutil.ErrInvalidArtifact
	}
	cfg, ok := global.ArtifactAdvanceCfg[cfgId]
	if !ok {
		return errutil.ErrConfigNotFound
	}
	if _, ok := global.ArtifactAdvanceCfg[cfgId+1]; !ok {
		return errutil.ErrMaxLevel
	}

	// 器魂等级校验
	for _, v := range artifact.Soul {
		for _, v2 := range cfg.AdvanceLimit {
			if v.CfgId == v2.ID && v.Level < v2.Num {
				return errutil.ErrLevel
			}
		}
	}

	// 扣除消耗
	for _, v := range cfg.Consume {
		if !a.checkTakeSource(v.ID, int64(v.Num)) {
			return errutil.ErrItemNotEnough
		}
	}
	for _, v := range cfg.Consume {
		a.takeSource(v.ID, int64(v.Num), constant.ItemDelArtifact)
	}

	a.pushMsg(protocol.RouteArtifactUpdatePush, &pb.ArtifactUpdatePush{
		Artifact: &pb.ArtifactInfo{CfgId: artifact.CfgId},
		Type:     pb.ArtifactUpdatePush_DELETE,
	})

	oldId := artifact.CfgId
	artifact.CfgId++
	a.artifacts.Delete(oldId)
	a.artifacts.Set(artifact.CfgId, artifact)

	a.pushMsg(protocol.RouteArtifactUpdatePush, &pb.ArtifactUpdatePush{
		Artifact: artifact,
		Type:     pb.ArtifactUpdatePush_NEW,
	})

	a.artifactSkill(artifact.CfgId, false)

	a.RefreshStats()

	return nil
}

// 神兵装备
func (a *Artifact) artifactEquip(cfgId int32, isEquip bool) error {

	artifact := a.getArtifact(cfgId)
	if artifact == nil {
		return errutil.ErrInvalidArtifact
	}

	var num int
	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		if value.IsEquip {
			num++
		}
		return true
	})

	if isEquip {
		if num >= global.GlobalCfg.BaseValue.Epicweapon.MaxNum {
			return errutil.ErrArtifactEquipLimit
		}
		artifact.IsEquip = true
	} else {
		artifact.IsEquip = false
	}

	a.updateArtifact(cfgId)
	return nil
}

// 神兵技能解锁
func (a *Artifact) artifactSkill(cfgId int32, isNew bool) {
	cfg, ok := global.ArtifactAdvanceCfg[cfgId]
	if !ok {
		return
	}

	if cfg.SkillID == 0 {
		return
	}

	if isNew {
		a.addSkill(cfg.SkillID)
	} else {
		if oldcfg, ok := global.ArtifactAdvanceCfg[cfgId-1]; ok {
			if skill := a.getSkill(oldcfg.SkillID); skill != nil {
				skill.CfgId = cfg.SkillID
			}
		}
	}
}

func (a *Artifact) isArtifactSkillValid(cfgId int32) bool {
	var using bool
	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		cfg, ok := global.ArtifactAdvanceCfg[key]
		if !ok {
			return true
		}
		if cfg.SkillID == cfgId && value.IsEquip {
			using = true
			return false
		}
		return true
	})

	return using
}

// 神兵属性加成
func (a *Artifact) artifactStatus(obj *MapObject) {
	if obj == nil {
		logger.Debugf(" %d:%s 神兵属性增加失败 ", a.uid, a.name)
		return
	}

	a.artifacts.Range(func(key int32, value *pb.ArtifactInfo) bool {
		if cfg, ok := global.ArtifactAdvanceCfg[key]; ok {
			for _, v := range cfg.AdvanceValue {
				obj.stats.setStats(pb.Attribute(v.ID), v.Num)
			}
		}
		for _, v := range value.Soul {
			if cfg, ok := global.ArtifactStrengthCfg[v.CfgId*1000+v.Level]; ok {
				for _, v2 := range cfg.LevelValue {
					obj.stats.setStats(pb.Attribute(v2.ID), v2.Num)
				}
			}
		}
		return true
	})
}
