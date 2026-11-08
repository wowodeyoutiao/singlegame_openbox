package dispatcher

// RechargeSuccessEvent 充值成功事件
type RechargeSuccessEvent struct {
	PID    int64  // 玩家ID
	Amount string // 充值金额
}

// ActivityPushEvent 活动推送事件
type ActivityPushEvent struct {
	Type ActivityType // 活动类型
	PID  int64        // 玩家ID
	Data interface{}  // 事件数据
}

type UserShowRefreshEvent struct {
	UserID          int64 `json:"user_id"`
	ServerID        int32
	ZoneID          int32
	UserShowRefresh *UserShowRefresh
}

type UserShowRefresh struct {
	Name     string `json:"name"`
	Avatar   int    `json:"avatar"`
	Frame    int    `json:"frame"`
	Sex      int    `json:"sex"`
	Job      int    `json:"job"`
	RobotCID int    `json:"robot_cid"`
	ServerID int    `json:"server_id"`
}

type PVEWarMapBossDieCountEvent struct {
	ZoneID    int32 `json:"zone_id"`
	MonsterID int32 `json:"monster_id"`
}

type ActivatePurchasedSlot struct {
	UID           int64 `json:"uid"`
	PurchasedSlot int32 `json:"purchased_slot"`
}

type CreateCharacter struct {
	AccountId     int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`            // 账号ID
	CharacterId   int64  `protobuf:"varint,2,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`      // 角色ID（游戏内唯一标识）
	CharacterName string `protobuf:"bytes,3,opt,name=character_name,json=characterName,proto3" json:"character_name,omitempty"` // 角色名称
	ServerId      int64  `protobuf:"varint,4,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`               // 服务器ID
	ServerName    string `protobuf:"bytes,5,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`          // 服务器名称
	Level         int32  `protobuf:"varint,6,opt,name=level,proto3" json:"level,omitempty"`                                     // 角色等级
	Avatar        int32  `protobuf:"varint,7,opt,name=avatar,proto3" json:"avatar,omitempty"`                                   // 角色头像ID
	AvatarFrame   int32  `protobuf:"varint,8,opt,name=avatar_frame,json=avatarFrame,proto3" json:"avatar_frame,omitempty"`      // 角色头像框ID
}

type AdClickEvent struct {
	OpenID  string
	ClickID string
	Uid     int64
}
