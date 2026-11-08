package game

import (
	"encoding/json"
	"time"
	"unicode/utf8"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol/protobuf/pb"
)

// GuildCreate 处理创建行会请求
func (gm *GuildManager) guildCreate(p *Player, req *pb.GuildCreateRequest) (res *pb.GuildCreateResponse, err error) {
	res = &pb.GuildCreateResponse{}
	// can := sensitive.ContainsSensitiveWord(req.Name)
	// if can {
	// 	res.Code = ResponseCodeInvalid
	// 	res.Error = "行会名称包含敏感词"
	// 	return &pb.GuildCreateResponse{
	// 		Code:  res.Code,
	// 		Error: res.Error,
	// 	}, nil
	// }

	name := req.Name
	runeCount := utf8.RuneCountInString(name)
	// 检查名称是否合法
	if runeCount < 2 || runeCount > 7 {
		res.Code = ResponseCodeInvalid
		res.Error = "行会名称长度必须在2到10之间"
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}
	description := req.Description
	leaderId := p.uid

	// 在行会中不能创建行会
	if p.guildId > 0 {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrGuildMemberExist.Error()
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}

	// 检查名称是否已存在
	for _, g := range gm.guilds {
		if g.Guild.ServerId != p.serverID {
			continue
		}
		if g.Guild.Name == name {
			res.Code = ResponseCodeInvalid
			res.Error = "行会已存在"
			return &pb.GuildCreateResponse{
				Code:  res.Code,
				Error: res.Error,
			}, nil
		}
	}
	//检测是不是有钱
	if !p.checkTakeSource(int32(global.GlobalCfg.Guild.CreatCost.ItemId), int64(global.GlobalCfg.Guild.CreatCost.NeedNum)) {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrItemNotEnough.Error()
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}
	// 创建公会记录
	// 获取1级行会对应的人数上限
	MaxMembers := int32(0)
	guildCfg := global.GetGuildCfg(1) // 假设1是新公会的等级
	if guildCfg != nil {
		MaxMembers = guildCfg.MemberNum
	}
	// 初始化行会技能
	skills := gm.GetGuildSkillBegin()
	// 转成json
	skillsJson, err := json.Marshal(skills)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}
	//生成雪花id
	GuildID := snowflake.GenID()
	guild := &model.Guild{
		GuildID:          GuildID,
		Name:             name,
		LeaderID:         leaderId,
		Description:      description,
		CreatedBy:        leaderId,
		CreatedAt:        time.Now().Unix(),
		MaxMembers:       MaxMembers,
		Level:            1,
		AllowApplication: 1,
		RequireApproval:  1,
		Status:           1,
		Skill:            string(skillsJson),
		ServerId:         p.serverID,
	}

	// 保存到数据库
	if err := db.InsertGuild(guild); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}

	// 创建公会信息
	guildInfo := &GuildInfo{
		Guild:        guild,
		Members:      make(map[int64]*pb.GuildMemberInfo),
		Applications: make(map[int64]*model.GuildApplication),
		Skills:       make(map[int32]*GuildSkill, 0),
	}
	// 初始化行会技能
	for _, skill := range skills {
		guildInfo.Skills[skill.UUid] = &GuildSkill{
			UUid:  skill.UUid,
			CfgId: skill.CfgId,
			Level: 1,
			Exp:   0,
		}
	}
	// 添加会长作为第一个成员
	member := &model.GuildMember{
		GuildID:  guild.GuildID,
		UserID:   leaderId,
		JoinedAt: time.Now().Unix(),
		Role:     1, // 会长角色
		Status:   1,
	}
	if err := db.InsertGuildMember(member); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return &pb.GuildCreateResponse{
			Code:  res.Code,
			Error: res.Error,
		}, nil
	}
	p.guildId = guild.GuildID
	p.guildName = guild.Name
	guildInfo.Members[leaderId] = &pb.GuildMemberInfo{
		MemberId:          member.UserID,
		Position:          int32(member.Role),
		MemberName:        p.name,
		MemberLevel:       int32(p.level),
		MemberTotalCombat: int32(p.totalCombat),
	}
	guildInfo.MemberCount++
	p.takeSource(int32(global.GlobalCfg.Guild.CreatCost.ItemId), int64(global.GlobalCfg.Guild.CreatCost.NeedNum), constant.ItemDelGuildCreate)

	// 同步玩家基础信息
	p.syncUserBase()

	// 添加到管理器
	gm.guilds[guild.GuildID] = guildInfo
	gm.upGuildCreateShushuDate(p, guild.GuildID, guild.Name)
	p.tasks.dispatch(pb.TargetType_JOIN_GUILD, 0, 1)

	return &pb.GuildCreateResponse{
		Code:  ResponseCodeSuccess,
		Error: "",
		Guild: &pb.GuildInfo{
			Id:          guild.GuildID,
			Name:        guild.Name,
			Description: guild.Description,
		},
	}, nil
}
