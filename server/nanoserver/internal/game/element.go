package game

import (
	"nanoserver/protocol/protobuf/pb"
)

type ElementInfo struct {
	element       [pb.ElementSlot_MaxElementpSlot]ElementSlot // 元素槽位
	elementInwork [ElementPartnerMax]int64                    //上阵槽位
	elementMain   int32
}

func (p *Player) ElementInfoElementToPb() *pb.ElementEmpowerList {
	empower := &pb.ElementEmpowerList{}
	for _, v := range p.ElementInfo.element {
		empower.Elements = append(empower.Elements, &pb.ElementEmpower{
			SlotId:     v.SlotId,
			Level:      v.Level,
			BreakLevel: v.Break,
		})
	}
	return empower
}
func (p *Player) ElementInfoElementInworkToPb() *pb.ElementUids {
	uid := &pb.ElementUids{}
	for _, v := range p.ElementInfo.elementInwork {
		uid.Uids = append(uid.Uids, v)
	}
	return uid
}
