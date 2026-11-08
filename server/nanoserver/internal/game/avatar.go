package game

import (
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/utils"
	"nanoserver/protocol/protobuf/pb"
	"time"
)

// 玩家外观数据
type Avatar struct {
	currentAvatarID int32                        // 当前使用头像ID
	currentFrameID  int32                        // 当前使用头像框ID
	ownedAvatars    *utils.SafeMap[int32, int64] // key:头像ID, value:过期时间戳(0表示永久)
	ownedFrames     *utils.SafeMap[int32, int64] // key:头像框ID, value:过期时间戳(0表示永久)
}

func NewAvatar() *Avatar {
	a := &Avatar{
		currentAvatarID: global.GlobalCfg.BaseValue.InitialIcon.Avatar,
		currentFrameID:  global.GlobalCfg.BaseValue.InitialIcon.Avatarframe,
		ownedAvatars:    utils.NewSafeMap[int32, int64](),
		ownedFrames:     utils.NewSafeMap[int32, int64](),
	}
	a.ownedAvatars.Set(global.GlobalCfg.BaseValue.InitialIcon.Avatar, int64(0))
	a.ownedFrames.Set(global.GlobalCfg.BaseValue.InitialIcon.Avatarframe, int64(0))
	return a
}

// AddAvatar 添加头像
func (a *Avatar) AddAvatar(avatarID int32) {
	avatarcfg, ok := global.AvatarCfg[avatarID]
	if !ok {
		logger.Warnf("avatar %d not found in config", avatarID)
		return
	}
	if _, exists := a.ownedAvatars.Get(avatarID); exists {
		logger.Warnf("avatar %d already owned", avatarID)
	}

	// 0 为不过期，时间为获取当天为第一天，然后 expiry 是几天后的0点过期
	if avatarcfg.Period > 0 {
		now := time.Now()
		expiry := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, int(avatarcfg.Period))
		a.ownedAvatars.Set(avatarID, expiry.Unix())
		return
	}
	a.ownedAvatars.Set(avatarID, int64(0))
}

// GetCurrentAvatar 获取当前头像
func (a *Avatar) GetCurrentAvatar() int32 {
	return a.currentAvatarID
}

// GetCurrentFrame 获取当前头像框
func (a *Avatar) GetCurrentFrame() int32 {
	return a.currentFrameID
}

// AddFrame 添加头像框
func (a *Avatar) AddFrame(frameID int32) {
	framecfg, ok := global.AvatarCfg[frameID]
	if !ok {
		logger.Warnf("frame %d not found in config", frameID)
		return
	}
	if _, exists := a.ownedFrames.Get(frameID); exists {
		logger.Warnf("frame %d already owned", frameID)
	}

	// 0 为不过期，时间为获取当天为第一天，然后 expiry 是几天后的0点过期
	if framecfg.Period > 0 {
		now := time.Now()
		expiry := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, int(framecfg.Period))
		a.ownedFrames.Set(frameID, expiry.Unix())
		return
	}
	a.ownedFrames.Set(frameID, int64(0))
}

// SetCurrentAvatar 设置当前头像
func (a *Avatar) SetCurrentAvatar(avatarID int32) error {
	expireTime, exists := a.ownedAvatars.Get(avatarID)
	if !exists {
		return errutil.ErrAvatarNotFound
	}
	if expireTime != 0 && expireTime < time.Now().Unix() {
		return errutil.ErrAvatarExpired
	}
	a.currentAvatarID = avatarID
	return nil
}

// SetCurrentFrame 设置当前头像框
func (a *Avatar) SetCurrentFrame(frameID int32) error {
	expireTime, exists := a.ownedFrames.Get(frameID)
	if !exists {
		return errutil.ErrFrameNotFound
	}
	if expireTime != 0 && expireTime < time.Now().Unix() {
		return errutil.ErrAvatarFrameExpired
	}
	a.currentFrameID = frameID
	return nil
}

// LoadFromDB 从数据库加载数据
func (a *Avatar) LoadFromDB(p *Player, dbAvatar *model.Avatar) error {
	if a.ownedAvatars == nil {
		a.ownedAvatars = utils.NewSafeMap[int32, int64]()
	}
	for id, expiry := range dbAvatar.OwnedAvatars {
		a.ownedAvatars.Set(id, expiry)
	}

	if a.ownedFrames == nil {
		a.ownedFrames = utils.NewSafeMap[int32, int64]()
	}
	for id, expiry := range dbAvatar.OwnedFrames {
		a.ownedFrames.Set(id, expiry)
	}
	a.currentAvatarID = dbAvatar.CurrentAvatar
	a.currentFrameID = dbAvatar.CurrentFrame
	return nil
}

// ToDB 保存数据到数据库
func (a *Avatar) ToDB(p *Player) (dbAvatar *model.Avatar) {
	ownedAvatars := make(map[int32]int64, a.ownedAvatars.Len())
	a.ownedAvatars.Range(func(key int32, value int64) bool {
		ownedAvatars[key] = value
		return true
	})
	ownedFrames := make(map[int32]int64, a.ownedFrames.Len())
	a.ownedFrames.Range(func(key int32, value int64) bool {
		ownedFrames[key] = value
		return true
	})
	dbAvatar = &model.Avatar{
		UID:           p.Uid(),
		OwnedAvatars:  ownedAvatars,
		OwnedFrames:   ownedFrames,
		CurrentAvatar: a.currentAvatarID,
		CurrentFrame:  a.currentFrameID,
	}
	return dbAvatar
}

// ToProto 转换为protobuf格式
func (a *Avatar) ToProto() *pb.SyncUserAvatarPush {
	avatars := make([]*pb.PlayerAvatarInfo, 0)
	a.ownedAvatars.Range(func(id int32, expiry int64) bool {
		avatars = append(avatars, &pb.PlayerAvatarInfo{Id: id, ExpireTime: expiry})
		return true
	})
	frames := make([]*pb.PlayerAvatarInfo, 0)
	a.ownedFrames.Range(func(id int32, expiry int64) bool {
		frames = append(frames, &pb.PlayerAvatarInfo{Id: id, ExpireTime: expiry})
		return true
	})
	return &pb.SyncUserAvatarPush{
		Avatars:         avatars,
		Frames:          frames,
		CurrentAvatarId: a.currentAvatarID,
		CurrentFrameId:  a.currentFrameID,
	}
}

// Stats 计算属性
func (a *Avatar) Stats(p *Player) {
	now := time.Now().Unix()
	a.ownedAvatars.Range(func(id int32, expiry int64) bool {
		if avatarCfg, ok := global.AvatarCfg[id]; ok {
			if expiry != 0 && expiry < now {
				return true
			}
			for _, avatarAttr := range avatarCfg.Value {
				if avatarAttr.Job == int32(p.job) {
					p.stats.setStats(pb.Attribute(avatarAttr.Key), avatarAttr.Value)
				}
			}
		}
		return true
	})
	a.ownedFrames.Range(func(id int32, expiry int64) bool {
		if frameCfg, ok := global.AvatarCfg[id]; ok {
			if expiry != 0 && expiry < now {
				return true
			}
			for _, avatarAttr := range frameCfg.Value {
				if avatarAttr.Job == int32(p.job) {
					p.stats.setStats(pb.Attribute(avatarAttr.Key), avatarAttr.Value)
				}
			}
		}
		return true
	})
}

// AddInitialAvatarAndFrame 兼容代码,添加等级1就给的头像和头像框
func (a *Avatar) AddInitialAvatarAndFrame() {
	cfg, ok := global.LevelGiftCfg[1]
	if !ok {
		return
	}

	for _, gift := range cfg.Gift {
		item, ok := global.ItemCfg[gift.ID]
		if !ok {
			continue
		}

		switch item.MainType {
		case constant.ItemMainTypeAvatar:
			if !a.ownedAvatars.Has(gift.ID) {
				a.AddAvatar(gift.ID)
			}
		case constant.ItemMainTypeAvatarFrame:
			if !a.ownedFrames.Has(gift.ID) {
				a.AddFrame(gift.ID)
			}
		}
	}
}
