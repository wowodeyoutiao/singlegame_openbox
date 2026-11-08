package thinkingdatax

import (
	"nanoserver/global"
	"nanoserver/protocol"
)

type TDItemInfo struct {
	ItemID   int32  `json:"item_id"`
	ItemName string `json:"item_name"`
	ItemNum  int32  `json:"item_num"`
}

func ConvertSourceToItemDetails(source []protocol.SourceInfo) []*TDItemInfo {
	if source == nil {
		return nil
	}

	items := make([]*TDItemInfo, 0, len(source))
	for _, v := range source {
		if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
			items = append(items, &TDItemInfo{
				ItemID:   v.Id,
				ItemName: ItemCfg.Name,
				ItemNum:  int32(v.Num),
			})

		}
	}

	return items
}

type Attr struct {
	Attr  int64 `json:"attr"`
	Value int64 `json:"value"`
}
