package history

import (
	"encoding/json"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/protocol"
)

type SnapShot struct {
	BasicInfo *protocol.DeskBasicInfo `json:"basicInfo"`
}

type History struct {
	mode    int
	beginAt int64
	endAt   int64
	deskID  int64

	SnapShot
}

func New(deskID int64, mode int, name0, name1, name2, name3 string, basic *protocol.DeskBasicInfo) *History {
	return &History{
		beginAt: time.Now().Unix(),
		deskID:  deskID,
		mode:    mode,

		SnapShot: SnapShot{
			BasicInfo: basic,
		},
	}
}

func (h *History) Save() error {
	data, err := json.Marshal(&h.SnapShot)
	if err != nil {
		return err
	}

	t := &model.History{
		DeskId:  h.deskID,
		BeginAt: h.beginAt,
		Mode:    h.mode,
		EndAt:   time.Now().Unix(),

		Snapshot: string(data),
	}

	return db.InsertHistory(t)
}
