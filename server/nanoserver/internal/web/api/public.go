package api

import (
	"math"
	"nanoserver/internal/game"
	"nanoserver/protocol"
)

// 发送道具
func CheckSourceAdd(award []protocol.SourceInfo, Uid int64, Send bool, methon string) error {
	addpacket := protocol.AddSource{
		Uid:    Uid,
		Source: make([]*protocol.SourceInfo, len(award)),
		Method: methon,
	}
	// 复制 Attachments 字段
	for j, detail := range award {
		// 假设 Attachment 结构体中有可以复制的字段
		addpacket.Source[j] = &protocol.SourceInfo{
			Id:  detail.Id,
			Num: detail.Num,
		}
		addpacket.ShopType = detail.ShopType
	}
	result := game.CheckSourceAdd(&addpacket)
	if result != nil {
		return result
	}
	if Send {
		game.SourceAdd(addpacket)
	}
	return nil
}

// 扣减道具
func CheckSourceNum(award []protocol.SourceInfo, Uid int64, Send bool, methon string) error {
	checkpacket := protocol.AddSource{
		Uid:    Uid,
		Source: make([]*protocol.SourceInfo, len(award)),
		Method: methon,
	}
	checkpacketDel := protocol.AddSource{
		Uid:    Uid,
		Source: make([]*protocol.SourceInfo, len(award)),
		Method: methon,
	}
	// 复制 Attachments 字段
	for j, detail := range award {
		// 假设 Attachment 结构体中有可以复制的字段
		checkpacket.Source[j] = &protocol.SourceInfo{
			Id:       detail.Id,
			Num:      detail.Num,
			ShopType: detail.ShopType,
			Method:   methon,
		}
		checkpacketDel.Source[j] = &protocol.SourceInfo{
			Id:       detail.Id,
			Num:      int64(math.Abs(float64(detail.Num))),
			ShopType: detail.ShopType,
			Method:   methon,
		}
	}
	result := game.CheckSourceNum(&checkpacketDel)
	if result != nil {
		return result
	}
	if Send {
		game.SourceAdd(checkpacket)
	}
	return nil
}
func Notification(uid int64) error {
	game.Notification(uid, 1)
	return nil
}
func addPartner(award protocol.AddSource) error {
	game.AddPartner(&award)
	return nil
}
func addPet(award protocol.AddSource, NeedDel int32, reason string) error {
	err := game.AddPet(&award, NeedDel, reason)
	if err != nil {
		return err
	}
	return nil
}
