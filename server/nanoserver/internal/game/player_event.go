package game

import (
	"math"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/event"
	"nanoserver/protocol"
)

// registerEventHandlers player注册事件处理器
func (p *Player) registerEventHandlers() {
	p.eventMgr.Subscribe(event.EventTypeTaskProgress, p.onTaskDispatch)
	p.eventMgr.Subscribe(event.EventTypeAddResource, p.onAddSource)
}

// IMapObject interface
func (p *Player) onTaskDispatch(data interface{}) {
	eventData, ok := data.(*event.TaskEventData)
	if !ok {
		return
	}
	p.tasks.dispatch(eventData.Key, eventData.Id, eventData.Num)
}

func (p *Player) onAddSource(data interface{}) {
	eventData, ok := data.(*protocol.AddSource)
	if !ok {
		return
	}
	items := []global.IdNum{}
	for _, source := range eventData.Source {
		if source.Num == 0 {
			continue
		}
		if source.Num > 0 {
			items = append(items, global.IdNum{ID: source.Id, Num: int32(source.Num)})
		} else {
			p.takeSource(source.Id, int64(math.Abs(float64(source.Num))), eventData.Method)
		}
	}
	if len(items) > 0 {
		if err := p.checkAddSourceBatch(items); err == nil {
			p.addSourceBatch(items, eventData.Method)
		} else {
			attachment := []protocol.Attachment{}
			for _, item := range items {
				attachment = append(attachment, protocol.Attachment{UserPropsId: int64(item.ID), Num: int64(item.Num)})
			}
			p.SendMail(constant.ItemAddGame, attachment)
			p.logger.Errorln("发放道具失败,包裹空间不足,邮件发送")
		}
	}
}
