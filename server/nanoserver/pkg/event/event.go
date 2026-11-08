package event

import "nanoserver/protocol/protobuf/pb"

type EventType int32

const (
	EventTypeAddBuff        EventType = 1 //添加buff事件
	EventTypeFateCoinEffect EventType = 2 //命运金币事件-效果
	EventTypeRefreshStats   EventType = 3 //刷新属性
	//Player
	EventTypeTaskProgress EventType = 10001 //任务进度更新
	EventTypeAddResource  EventType = 10002 //添加/删除道具
	//Rank
	EventTypeRankPlayer EventType = 11001 //排行 玩家相关排行数据
)

type EventData struct {
	UId   int64
	CfgId int32
	Fn    func()
}

type TaskEventData struct {
	Key pb.TargetType //任务类型
	Id  int32
	Num int32
}
