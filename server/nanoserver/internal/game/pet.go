package game

import (
	"nanoserver/global"
)

const PetWarriorNormalSkillID = 2000000
const PetWizardNormalSkillID = 2000100
const PetRitualistNormalSkillID = 2000200

type Pet struct {
	*MapObject

	level  int32          //等级
	star   int32          //星级
	skills []*PlayerSkill //技能
	cfg    *global.PetConfig
}

//func CreatePet(cid int32) *Pet {
//	pet := &Pet{
//		skills: make([]*PlayerSkill, 0),
//		level:  1,
//		star:   0,
//	}
//	pet.cfg = global.GetPetCfg(cid)
//	pet.cfgId = cid
//
//	return pet
//}
