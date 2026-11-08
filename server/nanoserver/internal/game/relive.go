package game

import "nanoserver/protocol/protobuf/pb"

// ReliveProcess 复活流程
func (d *Desk) reliveProcess(data *ReliveData) {
	switch data.reliveType {
	case pb.ReliveType_RT_RELIVE_NORMAL:
		d.reliveNormal(data)
	}
}

// ReliveNormal 正常复活
func (d *Desk) reliveNormal(data *ReliveData) {
	uid := data.uid
	objectType := data.objectType
	obj := d.getObject(uid, objectType)
	if obj != nil {
		obj.Revive(nil, 0, data.reviveHp)
	}
}

func (d *Desk) getObject(uid int64, objType ObjectType) IMapObject {
	switch objType {
	case ObjectTypeMonster:
		if obj, err := d.monsterWithUid(uid); err == nil {
			return obj
		}

	case ObjectTypePlayer:
		if obj, err := d.playerWithId(uid); err == nil {
			return obj
		}

	case ObjectTypeRobot:
		if obj, err := d.robotWithId(uid); err == nil {
			return obj
		}

	case ObjectTypePartner:
		if obj, err := d.objectWithId(uid); err == nil && obj.Owner() != nil {
			if owner, err := d.playerWithId(obj.Owner().Uid()); err == nil {
				if par := owner.getPartnerInfo(uid); par != nil {
					return obj
				}
			} else if owner, err := d.robotWithId(obj.Owner().Uid()); err == nil {
				if par := owner.getPartnerInfo(uid); par != nil {
					return obj
				}
			}
		}

	case ObjectTypeDecoy:
		if obj, err := d.objectWithId(uid); err == nil {
			return obj
		}
	}
	return nil
}
