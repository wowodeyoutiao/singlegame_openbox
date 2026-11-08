package game

import (
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

type GuildInfo struct {
	Guild        *model.Guild                      // 公会基本信息
	Members      map[int64]*pb.GuildMemberInfo     // 公会成员列表，key是用户ID
	MemberCount  int32                             // 公会成员数量，value是成员数量
	Applications map[int64]*model.GuildApplication // 申请列表
	Skills       map[int32]*GuildSkill             // 公会技能列表[]*GuildSkill          // 公会技能
	BattleId     int32                             // 公会战ID BattleID
}
type GuildSkill struct {
	UUid  int32 `json:"uuid"`
	CfgId int32 `json:"cfgid"`
	Level int32 `json:"level"`
	Exp   int32 `json:"exp"`
}

// GuildManager 管理所有公会
type GuildManager struct {
	component.Base
	guilds map[int64]*GuildInfo // 公会列表，key是公会ID
	mutex  sync.RWMutex         // 读写锁
}

// 响应码
const (
	ResponseCodeSuccess = 0
	ResponseCodeError   = -1
	ResponseCodeInvalid = -2
	ResponseCodeFailed  = -3
)

var guildManager = NewGuildManager()

// 通用操作辅助函数

// validatePlayer 验证会话中的玩家，如果无效则返回错误响应
func (gm *GuildManager) validatePlayer(s *session.Session) (*Player, *pb.CommonResponse) {
	p, err := playerWithSession(s)
	if err != nil {
		res := &pb.CommonResponse{
			Code:  ResponseCodeError,
			Error: err.Error(),
		}
		return nil, res
	}
	return p, nil
}

// validateGuild 验证公会是否存在，如果无效则返回错误响应
func (gm *GuildManager) validateGuild(guildID int64) (*GuildInfo, *pb.CommonResponse) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	if guildID == 0 {
		res := &pb.CommonResponse{
			Code:  ResponseCodeInvalid,
			Error: "还未加入公会",
		}
		return nil, res
	}
	guild, exists := gm.guilds[guildID]
	if !exists {
		res := &pb.CommonResponse{
			Code:  ResponseCodeInvalid,
			Error: "Guild not found",
		}
		return nil, res
	}
	return guild, nil
}

// validateGuildMember 验证玩家是否是公会成员，如果无效则返回错误响应
func (gm *GuildManager) validateGuildMember(guild *GuildInfo, playerID int64) (*pb.GuildMemberInfo, *pb.CommonResponse) {
	// 首先检查内存中是否存在
	member, exists := guild.Members[playerID]
	if !exists {
		// 如果内存中没有，尝试从数据库加载
		members, err := db.QueryGuildMembers(guild.Guild.GuildID)
		if err != nil {
			res := &pb.CommonResponse{
				Code:  ResponseCodeError,
				Error: err.Error(),
			}
			return nil, res
		}

		// 将所有成员缓存到内存中
		for _, m := range members {
			// 获取用户的基本信息，战斗力和等级和最后登录时间
			player := gm.getUserInfo(m.UserID)
			if player == nil {
				res := &pb.CommonResponse{
					Code:  ResponseCodeError,
					Error: "User not found",
				}
				return nil, res
			}
			memberOnline := int32(0)
			if player.curState > PlayerStatusOffLine {
				memberOnline = int32(player.lastLoginAt.Unix())
			}
			guild.Members[m.UserID] = &pb.GuildMemberInfo{
				MemberId:          m.UserID,
				Position:          int32(m.Role),
				MemberName:        player.name,
				MemberLevel:       int32(player.level),
				MemberTotalCombat: int32(player.totalCombat),
				MemberOnline:      memberOnline,
			}
		}

		// 从数据库加载后再次检查
		member, exists = guild.Members[playerID]
		if !exists {
			res := &pb.CommonResponse{
				Code:  ResponseCodeInvalid,
				Error: "该玩家已不在行会",
			}
			return nil, res
		}
	}
	return member, nil
}

// validateGuildPermission 验证成员是否具有所需权限，如果无效则返回错误响应
func (gm *GuildManager) validateGuildPermission(member *pb.GuildMemberInfo, requiredRole int) *pb.CommonResponse {
	if int(member.Position) > requiredRole {
		res := &pb.CommonResponse{
			Code:  ResponseCodeInvalid,
			Error: "Insufficient permissions",
		}
		return res
	}
	return nil
}

// loadGuildMembers 如果尚未加载，则加载公会成员
func (gm *GuildManager) loadGuildMembers(guild *GuildInfo) error {
	if len(guild.Members) == 0 {
		members, err := db.QueryGuildMembers(guild.Guild.GuildID)
		if err != nil {
			return err
		}
		for _, m := range members {
			// 这里可以加载用户的基本信息，比如名字和等级
			// 这里需要从数据库获取
			player := gm.getUserInfo(m.UserID)
			if player == nil {
				continue
			}
			memberOnline := int32(0)
			if player.curState == PlayerStatusOffLine {
				memberOnline = int32(player.lastLoginAt.Unix())
			}
			gm.guilds[guild.Guild.GuildID].Members[m.UserID] = &pb.GuildMemberInfo{
				MemberId:          m.UserID,
				Position:          int32(m.Role),
				MemberName:        player.name,
				MemberLevel:       int32(player.level),
				MemberTotalCombat: int32(player.totalCombat),
				MemberOnline:      memberOnline,
			}
			// 战斗力计算
		}

	} else {
		// 从数据库加载
		for _, m := range guild.Members {
			// 这里可以加载用户的基本信息，比如名字和等级
			// 这里需要从数据库获取
			player := gm.getUserInfo(m.MemberId)
			if player == nil {
				continue
			}
			memberOnline := int32(0)
			if player.curState == PlayerStatusOffLine {
				memberOnline = int32(player.lastLoginAt.Unix())
			}
			gm.guilds[guild.Guild.GuildID].Members[m.MemberId] = &pb.GuildMemberInfo{
				MemberId:          m.MemberId,
				Position:          int32(m.Position),
				MemberName:        player.name,
				MemberLevel:       int32(player.level),
				MemberTotalCombat: int32(player.totalCombat),
				MemberOnline:      memberOnline,
			}
			// 战斗力计算
		}
	}

	return nil
}

// loadGuildApplications 如果尚未加载，则加载公会申请
func (gm *GuildManager) loadGuildApplications(guild *GuildInfo) error {
	if len(guild.Applications) == 0 {
		applications, err := db.QueryGuildApplications(guild.Guild.GuildID)
		if err != nil {
			return err
		}
		for _, a := range applications {
			guild.Applications[a.ApplicantID] = a
		}
	}
	return nil
}

// createSuccessResponse 创建成功响应，可选包含数据
func (gm *GuildManager) createSuccessResponse(data interface{}) *pb.CommonResponse {
	return &pb.CommonResponse{
		Code: ResponseCodeSuccess,
	}
}

// createErrorResponse 创建包含给定错误消息的错误响应
func (gm *GuildManager) createErrorResponse(err error) *pb.CommonResponse {
	return &pb.CommonResponse{
		Code:  ResponseCodeError,
		Error: err.Error(),
	}
}

// 初始化,从数据库加载所有行会信息
func (gm *GuildManager) AfterInit() {
	// 加载所有行会
	guilds, err := db.QueryGuildList()
	if err != nil {
		logger.Error("加载行会列表失败:", err)
		return
	}
	// 查询每个行会有多少人
	memberCounts, err := db.QueryGuildMemberCount()
	if err != nil {
		logger.Error("加载行会成员数量失败:", err)
		return
	}
	// 加载每个行会的成员数量
	guildMember := make(map[int64]int32)
	for _, memberCount := range memberCounts {
		// 加载成员数量
		guildMember[memberCount.GuildID] = memberCount.MemberCount
	}

	for _, g := range guilds {
		// 加载行会技能，之前是json转的，现在转回来
		var guildSkills map[int32]*GuildSkill
		if g.Skill != "" {
			// 反序列化
			if err := json.Unmarshal([]byte(g.Skill), &guildSkills); err != nil {
				return
			}
		}
		memberCount, ok := guildMember[g.GuildID]
		if !ok {
			memberCount = 0
		}
		guildInfo := &GuildInfo{
			Guild:        g,
			Members:      make(map[int64]*pb.GuildMemberInfo),
			Applications: make(map[int64]*model.GuildApplication),
			Skills:       guildSkills,
			MemberCount:  memberCount,
		}
		gm.guilds[g.GuildID] = guildInfo
	}
	// 加载工会申请列表
	// 如果内存中没有，尝试从数据库加载
	// applications, err := db.QueryGuildApplicationsAll()
	// if err != nil {
	// 	logger.Error("加载行会申请列表失败:", err)
	// }
	// // 将所有申请缓存到内存中
	// for _, a := range applications {
	// 	gm.guilds[a.GuildID].Applications[a.ApplicantID] = 	a
	// }
}

// NewGuildManager 创建新的公会管理器
func NewGuildManager() *GuildManager {
	return &GuildManager{
		guilds: make(map[int64]*GuildInfo),
	}
}

// GuildList 处理获取行会列表请求
func (gm *GuildManager) GuildList(s *session.Session, req *pb.GuildListRequest) error {
	res := &pb.GuildListResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	for _, guildInfo := range gm.guilds {
		if p.serverID != guildInfo.Guild.ServerId {
			continue
		}

		// // 分页
		// if req.Page > 0 && req.PageSize > 0 {
		// 	if int64(i) < int64(req.Page-1)*int64(req.PageSize) || int64(i) >= int64(req.Page)*int64(req.PageSize) {
		// 		continue
		// 	}
		// }
		// 获取行会的LeaderName
		LeaderName := ""
		LeaderInfo := gm.getUserInfo(guildInfo.Guild.LeaderID)
		if LeaderInfo == nil {
			LeaderName = "未知"
		} else {
			LeaderName = LeaderInfo.name
		}
		// 是否申请 0：加入行会 1:  已申请 2: 未申请
		Application := int32(0) // 已加入行会，申请状态为0

		if p.guildId == 0 {
			// 获取申请状态
			app, errRes := gm.validateGuildApplication(guildInfo, p.uid)
			if errRes != nil {
				Application = 2 // 没有申请
			} else {
				Application = int32(app.Status) // 申请状态
			}
		}
		skills := make([]*pb.GuildSkillInfo, 0)
		for _, skill := range guildInfo.Skills {
			skills = append(skills, &pb.GuildSkillInfo{
				CfgId: skill.CfgId,
				Level: skill.Level,
				Exp:   skill.Exp,
			})
		}
		pbGuildInfo := &pb.GuildInfo{
			Id:          guildInfo.Guild.GuildID,
			Name:        guildInfo.Guild.Name,
			LeaderId:    guildInfo.Guild.LeaderID,
			LeaderName:  LeaderName,
			Description: guildInfo.Guild.Description,
			CreatedAt:   guildInfo.Guild.CreatedAt,
			MemberCount: guildInfo.MemberCount,
			MaxMembers:  int32(guildInfo.Guild.MaxMembers),
			Level:       int32(guildInfo.Guild.Level),
			Application: Application,
			GuildExp:    int64(guildInfo.Guild.GuildExp),
			Skills:      skills,
		}
		res.Guilds = append(res.Guilds, pbGuildInfo)
	}
	// logger.Info("GuildList:", res.Guilds)
	return s.Response(res)
}

// GuildCreate 处理创建行会请求
// func (gm *GuildManager) GuildCreate(s *session.Session, req *pb.GuildCreateRequest) error {
// 	res := &pb.GuildCreateResponse{}

// 	// 验证玩家
// 	p, errRes := gm.validatePlayer(s)
// 	if errRes != nil {
// 		res.Code = errRes.Code
// 		res.Error = errRes.Error
// 		return s.Response(res)
// 	}

// 	if strings.ToLower(viper.GetString("core.env")) == constant.EnvProd {
// 		if err := provider.WeChatV2.CheckMessageSecurity(p.accountName, req.Name); err != nil {
// 			res.Code = -1
// 			res.Error = err.Error()
// 			return s.Response(res)
// 		}
// 	}

// 	can := sensitive.ContainsSensitiveWord(req.Name)
// 	if can {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = "行会名称包含敏感词"
// 		return s.Response(res)
// 	}
// 	name := req.Name
// 	runeCount := utf8.RuneCountInString(name)
// 	// 检查名称是否合法
// 	if runeCount < 2 || runeCount > 7 {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = "行会名称长度必须在2到10之间"
// 		return s.Response(res)
// 	}
// 	description := req.Description
// 	leaderId := p.uid

// 	// 在行会中不能创建行会
// 	if p.guildId > 0 {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = errutil.ErrGuildMemberExist.Error()
// 		return s.Response(res)
// 	}

// 	// 检查名称是否已存在
// 	for _, g := range gm.guilds {
// 		if g.Guild.ServerId != p.serverID {
// 			continue
// 		}
// 		if g.Guild.Name == name {
// 			res.Code = ResponseCodeInvalid
// 			res.Error = "行会已存在"
// 			return s.Response(res)
// 		}
// 	}
// 	// 创建公会记录
// 	// 获取1级行会对应的人数上限
// 	MaxMembers := int32(0)
// 	guildCfg := global.GetGuildCfg(1) // 假设1是新公会的等级
// 	if guildCfg != nil {
// 		MaxMembers = guildCfg.MemberNum
// 	}
// 	// 初始化行会技能
// 	skills := gm.GetGuildSkillBegin()
// 	// 转成json
// 	skillsJson, err := json.Marshal(skills)
// 	if err != nil {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = err.Error()
// 		return s.Response(res)
// 	}
// 	guild := &model.Guild{
// 		Name:             name,
// 		LeaderID:         leaderId,
// 		Description:      description,
// 		CreatedBy:        leaderId,
// 		CreatedAt:        time.Now().Unix(),
// 		MaxMembers:       MaxMembers,
// 		Level:            1,
// 		AllowApplication: 1,
// 		RequireApproval:  1,
// 		Status:           1,
// 		Skill:            string(skillsJson),
// 		ServerId:         p.serverID,
// 	}

// 	// 保存到数据库
// 	if err := db.InsertGuild(guild); err != nil {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = err.Error()
// 		return s.Response(res)
// 	}

// 	// 创建公会信息
// 	guildInfo := &GuildInfo{
// 		Guild:        guild,
// 		Members:      make(map[int64]*pb.GuildMemberInfo),
// 		Applications: make(map[int64]*model.GuildApplication),
// 		Skills:       make(map[int32]*GuildSkill, 0),
// 	}
// 	// 初始化行会技能
// 	for _, skill := range skills {
// 		guildInfo.Skills[skill.UUid] = &GuildSkill{
// 			CfgId: skill.CfgId,
// 			Level: 1,
// 			Exp:   0,
// 		}
// 	}
// 	// 添加会长作为第一个成员
// 	member := &model.GuildMember{
// 		GuildID:  guild.GuildID,
// 		UserID:   leaderId,
// 		JoinedAt: time.Now().Unix(),
// 		Role:     1, // 会长角色
// 		Status:   1,
// 	}
// 	if err := db.InsertGuildMember(member); err != nil {
// 		res.Code = ResponseCodeInvalid
// 		res.Error = err.Error()
// 		return s.Response(res)
// 	}
// 	p.guildId = guild.GuildID
// 	guildInfo.Members[leaderId] = &pb.GuildMemberInfo{
// 		MemberId:          member.UserID,
// 		Position:          int32(member.Role),
// 		MemberName:        p.name,
// 		MemberLevel:       int32(p.level),
// 		MemberTotalCombat: int32(p.totalCombat),
// 	}
// 	guildInfo.MemberCount++

// 	// 同步玩家基础信息
// 	p.syncUserBase()

// 	// 添加到管理器
// 	gm.guilds[guild.GuildID] = guildInfo
// 	gm.upGuildCreateShushuDate(p, guild.GuildID, guild.Name)

// 	return s.Response(res)
// }

func (gm *GuildManager) GuildCreate(s *session.Session, req *pb.GuildCreateRequest) error {
	res := &pb.GuildCreateResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	if !p.isModuleOpen(constant.ModuleGuild) {
		res.Code = -21
		res.Error = errutil.ErrModuleNotOpen
		return s.Response(res)
	}

	resp, err := gm.guildCreate(p, req)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	return s.Response(resp)
}

func (gm *GuildManager) upGuildCreateShushuDate(p *Player, guild_id int64, guild_name string) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "create_guild"
	// 新增字段
	shushuDate["guild_id"] = guild_id     // 可根据实际装备ID填写
	shushuDate["guild_level"] = 1         // 可根据实际装备数量填写
	shushuDate["guild_name"] = guild_name // 可根据实际装备品质填写
	shushumain.ReportData(shushuDate)
}

// GuildApply 申请加入行会请求
func (gm *GuildManager) GuildApply(s *session.Session, req *pb.GuildApplyRequest) error {
	res := &pb.GuildApplyResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	guildIds := req.GuildIds
	userId := p.uid
	if p.guildId > 0 {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrGuildMemberExist.Error()
		return s.Response(res)
	}
	// 查看用户是不是主动退出行会
	member, err := db.QueryGuildMember(userId, 2)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	if member != nil {
		// 离开行会未超过24小时，不能再次申请
		if time.Now().Unix()-member.LeftAt < 60*60*24 {
			// 告知用户还剩 多少小时多少分钟多少秒
			leftTime := 60*60*24 - (time.Now().Unix() - member.LeftAt)
			// 转换成多少小时多少分钟多少秒
			leftTimeStr := strconv.FormatInt(leftTime/3600, 10) + "小时" + strconv.FormatInt((leftTime%3600)/60, 10) + "分钟" + strconv.FormatInt(leftTime%60, 10) + "秒"
			res.Code = ResponseCodeInvalid
			res.Error = "请" + leftTimeStr + "后再申请"
			return s.Response(res)
		}
	}

	for _, guildId := range guildIds {
		// 申请加入行会
		guildInfo, errRes := gm.validateGuild(guildId)
		if errRes != nil {
			res.Code = errRes.Code
			res.Error = errRes.Error
			continue
		}

		// 检查是否已经申请过
		if _, errRes := gm.validateGuildApplication(guildInfo, userId); errRes == nil {
			res.Code = ResponseCodeInvalid
			res.Error = errutil.ErrApplyRepeat.Error()
			continue
		}

		// 创建申请记录
		application := &model.GuildApplication{
			GuildID:     guildId,
			ApplicantID: userId,
			AppliedAt:   time.Now().Unix(),
			Status:      1,
		}

		if err := db.InsertGuildApplication(application); err != nil {
			logger.Error(err.Error())
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
		gm.guilds[guildId].Applications[userId] = application
	}
	//2025年10月27日11:46:43策划说申请加入公会也算完成任务
	p.tasks.dispatch(pb.TargetType_JOIN_GUILD, 0, 1)
	return s.Response(res)
}

// GuildCancelApply 取消申请
func (gm *GuildManager) GuildCancelApply(s *session.Session, req *pb.GuildCancelApplyRequest) error {
	res := &pb.GuildCancelApplyResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	guildId := req.GuildId
	userId := p.uid
	if p.guildId > 0 {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrGuildMemberExist.Error()
		return s.Response(res)
	}

	guildInfo, errRes := gm.validateGuild(guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查申请是否存在
	app, errRes := gm.validateGuildApplication(guildInfo, userId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = "没有申请过"
		return s.Response(res)
	}

	// 删除申请记录
	app.Status = 4 // 已拒绝状态
	if err := db.UpdateGuildApplication(app); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// 从公会信息中移除
	delete(gm.guilds[guildId].Applications, userId)
	return s.Response(res)
}

// 申请行会列表
func (gm *GuildManager) GuildApplyList(s *session.Session, req *pb.GuildApplyListRequest) error {
	res := &pb.GuildApplyListResponse{}

	// Validate player
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	guildId := p.guildId
	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	memberInfo, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否有权限
	if errRes := gm.validateGuildPermission(memberInfo, 2); errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// Load applications if needed
	if err := gm.loadGuildApplications(guildInfo); err != nil {
		logger.Error("加载行会申请失败:", err)
		return s.Response(res)
	}

	if len(guildInfo.Applications) == 0 {
		return s.Response(res)
	}

	for _, app := range guildInfo.Applications {
		// 获取用户的基本信息，战斗力和等级和最后登录时间
		player := gm.getUserInfo(app.ApplicantID)
		if player == nil {
			res := &pb.CommonResponse{
				Code:  ResponseCodeError,
				Error: errutil.ErrGuildNotFound.Error(),
			}
			return s.Response(res)
		}
		pbApp := &pb.GuildApplyInfo{
			ApplicantId:          app.ApplicantID,
			ApplicantName:        player.name,
			ApplicantLevel:       player.level,
			ApplicantTotalCombat: player.totalCombat,
		}
		res.Applications = append(res.Applications, pbApp)
	}
	return s.Response(res)
}

// GuildHandleApply 处理加入行会申请
func (gm *GuildManager) GuildHandleApply(s *session.Session, req *pb.GuildHandleApplyRequest) error {
	res := &pb.GuildHandleApplyResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	action := req.Accept // 1: 接受, 2: 拒绝

	if len(req.ApplicantIds) == 0 {
		res.Code = ResponseCodeInvalid
		res.Error = "没有申请人"
		return s.Response(res)
	}
	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}
	// 检查是否已经是成员
	member, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查权限
	if errRes := gm.validateGuildPermission(member, 2); errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	for _, userId := range req.ApplicantIds {

		// 检查申请是否存在
		app, errRes := gm.validateGuildApplication(guildInfo, userId)
		if errRes != nil {
			res.Code = errRes.Code
			res.Error = errRes.Error
			return s.Response(res)
		}

		if action == 1 { // 接受
			// 满员不可以加入
			if len(guildInfo.Members) >= int(guildInfo.Guild.MaxMembers) {
				res.Code = ResponseCodeInvalid
				res.Error = "行会已满员"
				return s.Response(res)
			}
			// 加载成员信息
			player := gm.getUserInfo(userId)
			if player == nil {
				res.Code = ResponseCodeError
				res.Error = "player not found"
				return s.Response(res)
			}
			// 检查是否已经是成员
			if player.guildId > 0 {
				res.Code = ResponseCodeInvalid
				res.Error = "对方已加入其他行会"
				app.Status = 5 // 因为已经是成员，取消申请
				if err := db.UpdateGuildApplication(app); err != nil {
					res.Code = ResponseCodeInvalid
					res.Error = err.Error()
					return s.Response(res)
				}
				delete(gm.guilds[p.guildId].Applications, userId)
				return s.Response(res)
			} else {
				// 更新申请表，插入成员表
				app.Status = 2 // 已接受状态
				if err := db.UpdateGuildApplication(app); err != nil {
					res.Code = ResponseCodeInvalid
					res.Error = err.Error()
					return s.Response(res)
				}
			}

			// 插入成员记录
			InsertGuildmember := &model.GuildMember{
				GuildID:  guildInfo.Guild.GuildID,
				UserID:   userId,
				JoinedAt: time.Now().Unix(),
				Role:     3, // 普通成员角色
				Status:   1,
			}
			if err := db.InsertGuildMember(InsertGuildmember); err != nil {
				res.Code = ResponseCodeInvalid
				res.Error = err.Error()
				return s.Response(res)
			}

			memberOnline := int32(0)
			if player.curState == PlayerStatusOffLine {
				memberOnline = int32(player.lastLoginAt.Unix())
			}
			gm.guilds[p.guildId].Members[userId] = &pb.GuildMemberInfo{
				MemberId:          InsertGuildmember.UserID,
				Position:          int32(InsertGuildmember.Role),
				MemberName:        player.name,          //
				MemberLevel:       int32(player.level),  //
				MemberTotalCombat: int32(p.totalCombat), //
				MemberOnline:      memberOnline,         //
			}
			gm.guilds[p.guildId].MemberCount++

			// 更新玩家公会ID
			if owner, ok := defaultManager.player(userId); ok {
				owner.guildId = guildInfo.Guild.GuildID
				owner.guildName = guildInfo.Guild.Name
				// 推送给玩家
				owner.syncUserBase()
				owner.tasks.dispatch(pb.TargetType_JOIN_GUILD, 0, 1)
			}
			gm.upGuildJoinShushuDate(p, guildInfo.Guild.GuildID, guildInfo.Guild.Name, int32(gm.guilds[p.guildId].Guild.Level))
			res.SuccessIds = append(res.SuccessIds, userId)
		} else {
			// 拒绝
			app.Status = 3 // 已拒绝状态
			if err := db.UpdateGuildApplication(app); err != nil {
				res.Code = ResponseCodeInvalid
				res.Error = err.Error()
				return s.Response(res)
			}
		}
		// 从申请列表中移除
		delete(gm.guilds[p.guildId].Applications, userId)
	}

	return s.Response(res)
}

func (gm *GuildManager) upGuildJoinShushuDate(p *Player, guild_id int64, guild_name string, guild_level int32) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "join_guild"
	// 新增字段
	shushuDate["guild_id"] = guild_id       // 可根据实际装备ID填写
	shushuDate["guild_level"] = guild_level // 可根据实际装备数量填写
	shushuDate["guild_name"] = guild_name   // 可根据实际装备品质填写
	shushumain.ReportData(shushuDate)

	shushuUserSetDate := make(map[string]interface{}, 0)
	shushuUserSetDate["current_guild_id"] = guild_id
	shushuUserSetDate["tmp_role_id"] = p.uid
	shushuUserSetDate["tmp_account"] = p.accountName
	shushumain.ReportUserSetData(shushuUserSetDate)
}

func (gm *GuildManager) GetGuild(guildId int64) *GuildInfo {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return gm.guilds[guildId]
}

// GuildDismiss 处理解散行会请求
func (gm *GuildManager) GuildDismiss(s *session.Session, req *pb.GuildDismissRequest) error {
	res := &pb.GuildDismissResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	guildId := req.GuildId
	// 检查是否是会长
	guildInfo, errRes := gm.validateGuild(guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	if guildInfo.Guild.LeaderID != p.uid {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrGuildPermission.Error()
		return s.Response(res)
	}

	// 删除公会记录
	if err := db.DeleteGuild(guildId); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// 更新所有成员状态
	for _, member := range guildInfo.Members {
		memberModle := &model.GuildMember{
			GuildID: guildId,
			UserID:  member.MemberId,
			Status:  2, // 退出状态
			LeftAt:  time.Now().Unix(),
		}
		if err := db.UpdateGuildMember(memberModle); err != nil {
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
	}

	// 如果需要，加载申请
	if err := gm.loadGuildApplications(guildInfo); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// 删除所有申请
	for _, app := range guildInfo.Applications {
		app.Status = 3 // 已拒绝状态
		if err := db.UpdateGuildApplication(app); err != nil {
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
	}

	// 从管理器中移除
	delete(gm.guilds, guildId)

	return s.Response(res)
}

// GuildLeave 离开行会请求
func (gm *GuildManager) GuildLeave(s *session.Session, req *pb.GuildLeaveRequest) error {
	res := &pb.GuildLeaveResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	guildId := req.GuildId
	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	member, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 会长不能直接退出
	if member.Position == 1 {
		if len(guildInfo.Members) == 1 {
			delete(gm.guilds, guildId)
			// 同步更新mysql
		} else {
			delete(gm.guilds[guildId].Members, p.uid)
			gm.guilds[p.guildId].MemberCount--

			leaderId := int64(0)
			leaderMaxBattle := int32(0)
			// 给战斗力最高的副会长
			for _, v := range guildInfo.Members {
				if v.MemberId == guildInfo.Guild.LeaderID {
					continue
				}
				if v.Position == 2 {
					player := gm.getUserInfo(v.MemberId)
					if player == nil {
						continue
					}
					// 保存最高的战斗力信息
					if player.totalCombat > leaderMaxBattle {
						leaderId = v.MemberId
						leaderMaxBattle = player.totalCombat
					}
				}
			}
			if leaderId == 0 {
				for _, v := range guildInfo.Members {
					if v.MemberId == guildInfo.Guild.LeaderID {
						continue
					}
					if v.Position == 3 {
						player := gm.getUserInfo(v.MemberId)
						if player == nil {
							continue
						}
						// 保存最高的战斗力信息
						if player.totalCombat > leaderMaxBattle {
							leaderId = v.MemberId
							leaderMaxBattle = player.totalCombat
						}
					}
				}
			}
			gm.guilds[guildId].Guild.LeaderID = leaderId
			// 提升某人的职位
			gm.guilds[guildInfo.Guild.GuildID].Members[leaderId].Position = 1
			guildInfo.Guild.LeaderID = leaderId
			if err := db.UpdateGuild(guildInfo.Guild); err != nil {
				res.Code = ResponseCodeInvalid
				res.Error = err.Error()
				return s.Response(res)
			}
			memberModelTo := &model.GuildMember{
				GuildID: guildId,
				UserID:  p.uid,
				Role:    1, // 会长
			}
			if err := db.UpdateGuildMember(memberModelTo); err != nil {
				res.Code = ResponseCodeFailed
				res.Error = err.Error()
				return s.Response(res)
			}

		}
	} else {
		delete(gm.guilds[guildId].Members, p.uid)
		gm.guilds[p.guildId].MemberCount--
	}
	memberModel := &model.GuildMember{
		GuildID: guildId,
		UserID:  p.uid,
		Status:  2, // 退出状态
		LeftAt:  time.Now().Unix(),
	}
	if err := db.UpdateGuildMember(memberModel); err != nil {
		res.Code = ResponseCodeFailed
		res.Error = err.Error()
		return s.Response(res)
	}

	p.guildId = 0
	p.guildName = ""
	// 如果是最后一个人退出，行会就解散了
	if gm.guilds[guildId] == nil {
		// 同步更新mysql
		if err := db.DeleteGuild(guildId); err != nil {
			return s.Response(res)
		}
	}
	// 同步个人信息
	p.syncUserBase()
	p.refreshStats()
	return s.Response(res)
}

// GuildKick 处理踢出行会成员请求
func (gm *GuildManager) GuildKick(s *session.Session, req *pb.GuildKickRequest) error {
	res := &pb.GuildKickResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是会长
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查成员是否存在
	memberId := req.MemberId
	member, errRes := gm.validateGuildMember(guildInfo, memberId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 踢出成员
	// 不能踢出会长
	if member.Position != 3 {
		res.Code = ResponseCodeInvalid
		res.Error = "不能踢出会长/副会长"
		return s.Response(res)
	}

	// 更新成员状态
	memberModel := &model.GuildMember{
		GuildID: guildInfo.Guild.GuildID,
		UserID:  memberId,
		Status:  3,
		LeftAt:  time.Now().Unix(),
	}
	if err := db.UpdateGuildMember(memberModel); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	delete(gm.guilds[guildInfo.Guild.GuildID].Members, memberId)

	// 更新玩家公会ID
	if kickedPlayer, ok := defaultManager.player(memberId); ok {
		kickedPlayer.guildId = 0
		kickedPlayer.guildName = ""
		kickedPlayer.syncUserBase()
		kickedPlayer.refreshStats()
	}
	gm.guilds[p.guildId].MemberCount--

	return s.Response(res)
}

// GuildPromote 处理提升行会职位请求
// GuildPromote 处理公会成员晋升请求
// 只有会长有权限进行晋升操作
// 参数:
//   - s: 会话信息
//   - req: 晋升请求参数,包含公会ID、成员ID和新的角色
//
// 返回:
//   - error: 操作错误信息
//
// 业务流程:
//  1. 验证操作者是否为会长
//  2. 检查目标成员是否存在
//  3. 更新成员角色信息
//  4. 如果是转让会长,则更新公会信息
func (gm *GuildManager) GuildPromote(s *session.Session, req *pb.GuildPromoteRequest) error {
	res := &pb.GuildPromoteResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是会长
	guildInfo, errRes := gm.validateGuild(req.GuildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	if guildInfo.Guild.LeaderID != p.uid {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrGuildPermission.Error()
		return s.Response(res)
	}

	// 检查成员是否存在
	memberId := req.MemberId
	newRole := req.Role
	_, errRes = gm.validateGuildMember(guildInfo, memberId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 更新角色
	memberModel := &model.GuildMember{
		GuildID: guildInfo.Guild.GuildID,
		UserID:  memberId,
		Role:    int8(newRole),
	}
	if err := db.UpdateGuildMember(memberModel); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	// 更新内存里的信息
	gm.guilds[guildInfo.Guild.GuildID].Members[memberId].Position = newRole
	// 如果是转让会长，更新公会信息
	if newRole == 1 {
		guildInfo.Guild.LeaderID = memberId
		if err := db.UpdateGuild(guildInfo.Guild); err != nil {
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
	}

	return s.Response(res)
}

// validateGuildApplication 验证申请是否存在，如果无效则返回错误响应
func (gm *GuildManager) validateGuildApplication(guild *GuildInfo, applicantID int64) (*model.GuildApplication, *pb.CommonResponse) {
	// 首先检查内存中是否存在
	application, exists := guild.Applications[applicantID]
	if !exists {
		// 如果内存中没有，尝试从数据库加载
		applications, err := db.QueryGuildApplications(guild.Guild.GuildID)
		if err != nil {
			res := &pb.CommonResponse{
				Code:  ResponseCodeError,
				Error: err.Error(),
			}
			return nil, res
		}

		// 将所有申请缓存到内存中
		for _, a := range applications {
			guild.Applications[a.ApplicantID] = a
		}

		// 从数据库加载后再次检查
		application, exists = guild.Applications[applicantID]
		if !exists {
			res := &pb.CommonResponse{
				Code:  ResponseCodeInvalid,
				Error: "暂无行会申请信息",
			}
			return nil, res
		}
	}
	return application, nil
}

// 行会详情
// GuildInfo 获取公会信息
func (gm *GuildManager) GuildInfo(s *session.Session, req *pb.GuildInfoRequest) error {
	res := &pb.GuildInfoResponse{}
	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}
	// 加载成员信息
	// 查询会长信息
	leaderInfo := gm.getUserInfo(guildInfo.Guild.LeaderID)
	if leaderInfo == nil {
		res.Code = ResponseCodeInvalid
		res.Error = "Leader not found"
		return s.Response(res)
	}
	res.Guild = &pb.GuildInfo{
		Id:          guildInfo.Guild.GuildID,
		Name:        guildInfo.Guild.Name,
		LeaderId:    leaderInfo.uid,  // Populate with leader's name if available
		LeaderName:  leaderInfo.name, // Populate with leader's name if available
		Description: guildInfo.Guild.Description,
		CreatedAt:   guildInfo.Guild.CreatedAt,
		MemberCount: int32(len(guildInfo.Members)),
		MaxMembers:  int32(guildInfo.Guild.MaxMembers),
		Level:       int32(guildInfo.Guild.Level),
		GuildExp:    int64(guildInfo.Guild.GuildExp),
		Skills:      make([]*pb.GuildSkillInfo, 0), // Map skills if needed
	}
	for _, skill := range guildInfo.Skills {
		res.Guild.Skills = append(res.Guild.Skills, &pb.GuildSkillInfo{
			CfgId: skill.CfgId,
			Level: skill.Level,
			Exp:   skill.Exp,
		})
	}

	if guildInfo.Members[p.uid] != nil {
		res.Role = int32(guildInfo.Members[p.uid].Position)
	} else {
		res.Role = 0
	}
	return s.Response(res)
}

// 公会成员信息
// GuildMemberInfo 获取公会成员信息
func (gm *GuildManager) GuildMemberInfo(s *session.Session, req *pb.GuildMemberListRequest) error {
	res := &pb.GuildMemberListResponse{}
	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	_, errRes = gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}
	return s.Response(res)
}

// 公会成员列表
// GuildMemberInfo 获取公会成员信息
func (gm *GuildManager) GuildMemberList(s *session.Session, req *pb.GuildMemberListRequest) error {
	res := &pb.GuildMemberListResponse{}
	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	err := gm.loadGuildMembers(guildInfo)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	res.Members = make([]*pb.GuildMemberInfo, 0)
	for _, member := range guildInfo.Members {
		res.Members = append(res.Members, &pb.GuildMemberInfo{
			MemberId:          member.MemberId,
			Position:          member.Position,
			MemberName:        member.MemberName,
			MemberLevel:       member.MemberLevel,
			MemberTotalCombat: member.MemberTotalCombat,
			MemberOnline:      member.MemberOnline,
		})
	}

	return s.Response(res)
}

// GuildDonateInfo 获取公会捐献信息
func (gm *GuildManager) GuildDonateInfo(s *session.Session, req *pb.GuildDonateInfoRequest) error {
	res := &pb.GuildDonateInfoResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	_, errRes = gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 获取用户当天捐献记录
	donates, err := db.QueryGuildDonateInfo(guildInfo.Guild.GuildID, p.uid)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	res.Num = int32(len(donates))

	return s.Response(res)
}

// GuildDonate 公会捐献/签到
func (gm *GuildManager) GuildDonate(s *session.Session, req *pb.GuildDonateRequest) error {
	res := &pb.GuildDonateInfoResponse{}
	res.Num = 0
	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	member, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}
	// 获取用户当天捐献记录
	donates, err := db.QueryGuildDonateInfo(guildInfo.Guild.GuildID, p.uid)
	res.Num = int32(len(donates))
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	if len(donates) >= global.GlobalCfg.Guild.DonateLimit.DailyTime {
		res.Code = ResponseCodeInvalid
		res.Error = "捐献次数已达上限"
		return s.Response(res)
	}

	// 根据类型处理捐献
	oldGuildExp := gm.guilds[p.guildId].Guild.GuildExp
	oldGuildLevel := gm.guilds[p.guildId].Guild.Level
	cost_item_name := ""
	cost_item_num := 0
	get_item_name := ""
	get_item_num := 0
	switch req.Type {
	case 1: // 金币贡献

		if !p.checkTakeSource(int32(global.GlobalCfg.Guild.GoldDonate.NeedItemId), int64(global.GlobalCfg.Guild.GoldDonate.NeedItemNum)) {
			res.Code = ResponseCodeInvalid
			res.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(res)
		}
		p.takeSource(int32(global.GlobalCfg.Guild.GoldDonate.NeedItemId), int64(global.GlobalCfg.Guild.GoldDonate.NeedItemNum), constant.ItemDelGuildDonate)
		// 行会经验增加，自己的公会贡献增加
		err := gm.GuildAddExp(guildInfo.Guild.GuildID, int32(global.GlobalCfg.Guild.GoldDonate.GetItemNum)) // 增加行会经验
		if err != nil {
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
		p.addSource(int32(constant.ItemTypeGuildContribution), int64(global.GlobalCfg.Guild.GoldDonate.GetItemNum), constant.ItemAddGuildDonate)
		gm.guilds[guildInfo.Guild.GuildID].Members[member.MemberId].Points += int32(global.GlobalCfg.Guild.GoldDonate.GetGuildExp)
		if ItemCfg, err := global.GetItemCfg(int32(global.GlobalCfg.Guild.GoldDonate.NeedItemId)); err == nil {
			cost_item_name = ItemCfg.Name
		}
		cost_item_num = global.GlobalCfg.Guild.GoldDonate.NeedItemNum
		member.Points += int32(global.GlobalCfg.Guild.GoldDonate.GetGuildExp)
		if ItemCfg2, err := global.GetItemCfg(int32(constant.ItemTypeGuildContribution)); err == nil {
			get_item_name = ItemCfg2.Name
		}
		get_item_num = global.GlobalCfg.Guild.GoldDonate.GetItemNum
	case 2: // 点券捐献
		// 检查玩家点券是否足够
		if !p.checkTakeSource(int32(global.GlobalCfg.Guild.MoneyDonate.NeedItemId), int64(global.GlobalCfg.Guild.MoneyDonate.NeedItemNum)) {
			res.Code = ResponseCodeInvalid
			res.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(res)
		}
		p.takeSource(int32(global.GlobalCfg.Guild.MoneyDonate.NeedItemId), int64(global.GlobalCfg.Guild.MoneyDonate.NeedItemNum), constant.ItemDelGuildDonate)
		// 行会经验增加，自己的公会贡献增加
		err := gm.GuildAddExp(guildInfo.Guild.GuildID, int32(global.GlobalCfg.Guild.MoneyDonate.GetItemNum)) // 增加行会经验
		if err != nil {
			res.Code = ResponseCodeInvalid
			res.Error = err.Error()
			return s.Response(res)
		}
		p.addSource(int32(constant.ItemTypeGuildContribution), int64(global.GlobalCfg.Guild.MoneyDonate.GetItemNum), constant.ItemAddGuildDonate)
		gm.guilds[guildInfo.Guild.GuildID].Members[member.MemberId].Points += int32(global.GlobalCfg.Guild.MoneyDonate.GetGuildExp)
		if ItemCfg, err := global.GetItemCfg(int32(global.GlobalCfg.Guild.MoneyDonate.NeedItemId)); err == nil {
			cost_item_name = ItemCfg.Name
		}
		cost_item_num = global.GlobalCfg.Guild.MoneyDonate.NeedItemNum
		member.Points += int32(global.GlobalCfg.Guild.MoneyDonate.GetGuildExp)
		if ItemCfg2, err := global.GetItemCfg(int32(constant.ItemTypeGuildContribution)); err == nil {
			get_item_name = ItemCfg2.Name
		}
		get_item_num = global.GlobalCfg.Guild.MoneyDonate.GetItemNum

	default:
		res.Code = ResponseCodeInvalid
		res.Error = "无效的捐献类型"
		return s.Response(res)
	}
	newGuildExp := gm.guilds[p.guildId].Guild.GuildExp
	newGuildLevel := gm.guilds[p.guildId].Guild.Level
	// 更新成员信息// 更新成员捐献信息
	memberModel := &model.GuildMember{
		GuildID: guildInfo.Guild.GuildID,
		UserID:  p.uid,
		Points:  member.Points,
	}
	if err := db.UpdateGuildMember(memberModel); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// // 更新行会捐献信息
	// if err := db.UpdateGuild(guildInfo.Guild); err != nil {
	// 	res.Code = ResponseCodeInvalid
	// 	res.Error = err.Error()
	// 	return s.Response(res)
	// }

	// 添加捐献记录
	donate := &model.GuildDonate{
		GuildID:    guildInfo.Guild.GuildID,
		DonatorID:  p.uid,
		DonateType: req.Type,
		DonateTime: time.Now().Unix(),
	}
	if err := db.InsertGuildDonateInfo(donate); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	res.Num++
	gm.upGuildDonateShushuDate(p, guildInfo.Guild.GuildID, guildInfo.Guild.Name, int32(guildInfo.Guild.Level),
		oldGuildExp, newGuildExp, int32(oldGuildLevel), int32(newGuildLevel), cost_item_name, cost_item_num, get_item_name, get_item_num, global.GlobalCfg.Guild.GoldDonate.NeedItemId)
	return s.Response(res)
}

func (gm *GuildManager) upGuildDonateShushuDate(p *Player, guild_id int64, guild_name string, guildLevel int32,
	oldSkillExp int32, newSkillExp int32, oldSkillLevel int32, newSkillLevel int32, cost_item_name string, cost_item_num int,
	get_item_name string, get_item_num int, cost_item_id int,
) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "guild_donate"
	// 新增字段
	shushuDate["guild_id"] = guild_id      // 可根据实际装备ID填写
	shushuDate["guild_level"] = guildLevel // 可根据实际装备数量填写
	shushuDate["guild_name"] = guild_name  // 可根据实际装备品质填写
	shushuDate["cost_item_id"] = cost_item_id
	shushuDate["cost_item_name"] = cost_item_name
	shushuDate["cost_item_num"] = cost_item_num
	shushuDate["get_item_name"] = get_item_name
	shushuDate["get_item_num"] = get_item_num
	shushuDate["exp_before"] = oldSkillExp
	shushuDate["exp_after"] = newSkillExp
	shushuDate["level_before"] = oldSkillLevel
	shushuDate["level_after"] = newSkillLevel
	shushuDate["pray_times"] = 1
	shushumain.ReportData(shushuDate)
}

// GuildShopBuyInfo 获取公会商城购买情况
func (gm *GuildManager) GuildShopBuyInfo(s *session.Session, req *pb.GuildShopBuyInfoRequest) error {
	res := &pb.GuildShopBuyInfoResponse{}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	member, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	goodsIds := []int64{}
	// 获取公会商城商品
	for _, goods := range global.GuildStoreCfg {
		goodsIds = append(goodsIds, int64(goods.ID))
	}
	// 获取购买记录
	buyInfos, err := db.QueryGuildShopBuyInfo(1, guildInfo.Guild.GuildID, goodsIds, member.MemberId)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// Map buyInfos to the expected type
	var pbBuyInfos []*pb.GuildShopBuyInfo
	for _, info := range buyInfos {
		pbBuyInfos = append(pbBuyInfos, &pb.GuildShopBuyInfo{
			ItemId:   int32(info.GoodsID),
			BuyNum:   int32(info.UserNum),
			TotalNum: int32(info.GuildNum),
		})
	}
	res.BuyInfo = pbBuyInfos
	return s.Response(res)
}

// GuildShopBuy 公会商城购买
func (gm *GuildManager) GuildShopBuy(s *session.Session, req *pb.GuildShopBuyRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_GuildShopBuy)}

	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	member, errRes := gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}
	if req.BuyNum <= 0 {
		res.Code = ResponseCodeInvalid
		res.Error = "购买数量必须大于0"
		return s.Response(res)
	}
	goodsIds := []int64{}
	goodsIds = append(goodsIds, int64(req.Id))
	// 获取购买记录
	buyInfos, err := db.QueryGuildShopBuyInfo(1, guildInfo.Guild.GuildID, goodsIds, member.MemberId)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	// 查看这个商品的信息 global.GuildStoreCfg
	goodsCfg := global.GetGuildStoreCfg(int32(req.Id))
	if goodsCfg == nil {
		res.Code = ResponseCodeInvalid
		res.Error = "商品不存在"
		return s.Response(res)
	}

	// 检查购买次数是否超过限制
	for _, info := range buyInfos {
		if info.GoodsID == req.Id {
			// 个人限制
			if int32(info.UserNum)+req.BuyNum > int32(goodsCfg.PersonalLimit) {
				res.Code = ResponseCodeInvalid
				res.Error = "购买次数超过限制"
				return s.Response(res)
			}
			// 公会限制
			if int32(info.GuildNum)+req.BuyNum > int32(goodsCfg.GuildLimit) {
				res.Code = ResponseCodeInvalid
				res.Error = "公会购买次数超过限制"
				return s.Response(res)
			}
		}
	}

	// 检测有没有足够的道具
	for _, price := range goodsCfg.Price {
		if !p.checkTakeSource(int32(price.ID), int64(price.Num)*int64(req.BuyNum)) {
			res.Code = ResponseCodeInvalid
			res.Error = errutil.ErrItemNotEnough.Error()
			return s.Response(res)
		}
	}

	// 记录购买信息
	buyInfo := model.GuildStoreBuy{
		UserID:  member.MemberId,
		GuildID: guildInfo.Guild.GuildID,
		GoodsID: int64(req.Id),
		BuyNum:  int32(req.BuyNum),
		BuyTime: time.Now().Unix(),
		Status:  1,
	}

	// 插入购买记录
	if err := db.InsertGuildShopBuyInfo(&buyInfo); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}

	// 减去道具
	resultt := []protocol.SourceInfo{}
	cost_id := int32(0)
	for _, price := range goodsCfg.Price {
		cost_id = price.ID
		resultt = append(resultt, protocol.SourceInfo{
			Id:  int32(price.ID),
			Num: int64(price.Num) * int64(req.BuyNum),
		})
		p.takeSource(int32(price.ID), int64(price.Num)*int64(req.BuyNum), constant.ItemDelGuildShopBuy)
	}
	GeTNum, err := defaultManager.GetItemCount(p.uid, cost_id)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	cost_item_before := int64(GeTNum)

	// 给到用户购买的东西

	buy_name := int32(0)
	result2 := []protocol.SourceInfo{}
	for _, goods := range goodsCfg.Goods {
		result2 = append(result2, protocol.SourceInfo{
			Id:  int32(goods.ID),
			Num: int64(goods.Num) * int64(req.BuyNum),
		})
		buy_name = int32(goods.ID)
		p.addSource(int32(goods.ID), int64(goods.Num)*int64(req.BuyNum), constant.ItemAddGuildShopBuy)
	}
	CosNum, err := defaultManager.GetItemCount(p.uid, buy_name)
	if err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	get_item_before := int64(CosNum)
	// 	if result != nil {
	// 		if len(result) > 0 {
	// 			resultt = result[0]
	// 		}
	// 	}
	upStoreOneShushu(buy_name, "商城", result2[0], resultt[0], p.uid, int32(req.Id), cost_item_before, get_item_before)

	return s.Response(res)
}

func upStoreOneShushu(buy_name int32, buy_entrance string, sourceInfo protocol.SourceInfo, costInfo protocol.SourceInfo, uid int64, cfgId int32, cost_item_before int64, get_item_before int64) {
	data := make(map[string]interface{}, 0)
	data["name"] = "shop_buy_item"
	data["buy_name"] = buy_name
	data["buy_entrance"] = buy_entrance
	data["get_item_before"] = get_item_before
	data["get_item_after"] = get_item_before + sourceInfo.Num
	if ItemCfg, err := global.GetItemCfg(sourceInfo.Id); err == nil {
		data["get_item_id"] = ItemCfg.ID
		data["get_item_name"] = ItemCfg.Name
		data["get_item_num"] = sourceInfo.Num
	}
	data["cost_item_before"] = cost_item_before
	data["cost_item_after"] = cost_item_before - costInfo.Num
	if ItemCfg, err := global.GetItemCfg(costInfo.Id); err == nil {
		data["cost_item_id"] = ItemCfg.ID
		data["cost_item_name"] = ItemCfg.Name
		data["cost_item_num"] = math.Abs(float64(costInfo.Num))
	}
	// game.Upshushu(uid, data, true,false)
}

// GuildUpgradeSkill 公会升级技能
func (gm *GuildManager) GuildUpgradeSkill(s *session.Session, req *pb.GuildUpgradeSkillRequest) error {
	res := &pb.CommonResponse{Code: 0, RequestCmdId: uint32(pb.CmdId_C2S_GuildUpgradeSkill)}
	if req.UpgradeCount <= 0 {
		res.Code = ResponseCodeInvalid
		res.Error = "升级次数必须大于0"
		return s.Response(res)
	}
	is_all := false
	if req.UpgradeCount > 1 {
		is_all = true
	}
	// 验证玩家
	p, errRes := gm.validatePlayer(s)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 检查是否是行会成员
	guildInfo, errRes := gm.validateGuild(p.guildId)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	_, errRes = gm.validateGuildMember(guildInfo, p.uid)
	if errRes != nil {
		res.Code = errRes.Code
		res.Error = errRes.Error
		return s.Response(res)
	}

	// 获取技能信息
	skill := global.GetGuildSkillCfg(int32(req.UUId))
	if skill == nil {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrSkillsSlot.Error()
		return s.Response(res)
	}
	if guildInfo.Skills[int32(req.UUId)] == nil {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrSkillsSlot.Error()
		return s.Response(res)
	}
	// 检测有没有祈愿点
	// 检测有没有足够的道具
	if !p.checkTakeSource(int32(constant.ItemTypeGuildWish), int64(req.UpgradeCount)) {
		res.Code = ResponseCodeInvalid
		res.Error = errutil.ErrItemNotEnough.Error()
		return s.Response(res)
	}
	// 提供的技能经验
	upgradeGuild := int64(global.GlobalCfg.Guild.WishValue.GetSKillExp) * int64(req.UpgradeCount)
	// 提供的技能
	upgradeCost := int64(global.GlobalCfg.Guild.WishValue.GetitemNum) * int64(req.UpgradeCount)
	// 给行会技能加经验，如果经验足够，升级技能
	oldSkillExp := guildInfo.Skills[int32(req.UUId)].Exp
	oldSkillLevel := guildInfo.Skills[int32(req.UUId)].Level
	if err := gm.GuildAddSkillExp(guildInfo.Guild.GuildID, int32(req.UUId), upgradeGuild); err != nil {
		res.Code = ResponseCodeInvalid
		res.Error = err.Error()
		return s.Response(res)
	}
	newSkillExp := int32(0)
	newSkillLevel := int32(0)
	if guildInfo.Skills[int32(req.UUId)] != nil {
		newSkillExp = guildInfo.Skills[int32(req.UUId)].Exp
		newSkillLevel = guildInfo.Skills[int32(req.UUId)].Level
	} else {
		// 如果技能不存在，说明可能已经升级到下一级了
		if guildInfo.Skills[int32(req.UUId)+1] != nil {
			newSkillExp = guildInfo.Skills[int32(req.UUId)+1].Exp
			newSkillLevel = guildInfo.Skills[int32(req.UUId)+1].Level
		}
	}
	// //给行会加经验，如果经验足够，升级行会
	// if err := gm.GuildAddExp(guildInfo.Guild.GuildID, int32(upgradeCost)); err != nil {
	// 	res.Code = ResponseCodeInvalid
	// 	res.Error = err.Error()
	// 	return s.Response(res)
	// }
	// 给到玩家行会经验
	p.addSource(constant.ItemTypeGuildContribution, upgradeCost, constant.ItemAddGuildWishExp)

	// 减去道具
	p.takeSource(int32(constant.ItemTypeGuildWish), int64(req.UpgradeCount), constant.ItemDelGuildShopBuy)
	gm.upGuilSkillShushuDate(p, guildInfo.Guild.GuildID, guildInfo.Guild.Name, int32(guildInfo.Guild.Level), skill.SkillName,
		oldSkillExp, newSkillExp, int32(oldSkillLevel), int32(newSkillLevel), req.UpgradeCount, is_all)
	return s.Response(res)
}

func (gm *GuildManager) upGuilSkillShushuDate(p *Player, guild_id int64, guild_name string, guildLevel int32, skillName string,
	oldSkillExp int32, newSkillExp int32, oldSkillLevel int32, newSkillLevel int32, prayTimes int32, is_all bool,
) {
	// 创建一interface{}类型的map
	shushuDate := make(map[string]interface{}, 0)
	// upBaseInfo的数据合并进去
	baseInfo := p.ToMapForShushu()
	for k, v := range baseInfo {
		shushuDate[k] = v
	}
	shushuDate["name"] = "guild_skill"
	// 新增字段
	shushuDate["guild_id"] = guild_id      // 可根据实际装备ID填写
	shushuDate["guild_level"] = guildLevel // 可根据实际装备数量填写
	shushuDate["guild_name"] = guild_name  // 可根据实际装备品质填写
	shushuDate["skill_name"] = skillName
	shushuDate["exp_before"] = oldSkillExp
	shushuDate["exp_after"] = newSkillExp
	shushuDate["is_all"] = is_all
	shushuDate["level_before"] = oldSkillLevel
	shushuDate["level_after"] = newSkillLevel
	shushuDate["pray_times"] = prayTimes
	shushumain.ReportData(shushuDate)
}

// 给行会增加经验，如果升级经验充足，行会就升级
func (gm *GuildManager) GuildAddExp(guildId int64, exp int32) error {
	guildInfo := gm.GetGuild(guildId) // 获取行会信息
	if guildInfo == nil {
		return errutil.ErrGuildNotFound
	}
	gm.guilds[guildId].Guild.GuildExp += exp // 增加经验
	// 检查是否升级
	// 获取当前等级升级需要多少经验
	guildCfg := global.GetGuildCfg(int32(guildInfo.Guild.Level))
	if guildCfg == nil {
		return errutil.ErrGuildNotFound
	}
	// 获取当前等级升级需要多少经验
	if guildInfo.Guild.GuildExp >= guildCfg.NeedEXP {
		// 查询下一级的配置
		guildCfgNext := global.GetGuildCfg(int32(guildInfo.Guild.Level + 1))
		if guildCfgNext == nil {
			return nil
		}
		gm.guilds[guildId].Guild.Level++
		gm.guilds[guildId].Guild.GuildExp -= guildCfg.NeedEXP // 升级经验清零
		// 行会的属性同步更新
		gm.guilds[guildId].Guild.MaxMembers = guildCfgNext.MemberNum // 行会最大成员数
	}
	// 更新行会信息
	if err := db.UpdateGuild(gm.guilds[guildId].Guild); err != nil {
		return err
	}
	return nil
}

// 给行会技能增加经验，如果升级经验充足，行会技能就升级
func (gm *GuildManager) GuildAddSkillExp(guildId int64, SkillCfgId int32, exp int64) error {
	guildInfo := gm.GetGuild(guildId) // 获取行会信息
	if guildInfo == nil {
		return errutil.ErrGuildNotFound
	}
	// 检查是否升级

	// 获取当前等级升级需要多少经验
	guildCfg := global.GetGuildSkillCfg(int32(SkillCfgId))
	if guildCfg == nil {
		return errutil.ErrGuildNotFound
	}

	// 检查技能是否存在
	currentSkill, exists := gm.guilds[guildId].Skills[int32(SkillCfgId)]
	if !exists || currentSkill == nil {
		return errutil.ErrSkillsSlot
	}

	if currentSkill.Exp+int32(exp) >= guildCfg.EXP {
		// 升级 - 计算剩余经验
		remainingExp := currentSkill.Exp + int32(exp) - guildCfg.EXP

		// 创建下一级技能
		guildInfo.Skills[int32(SkillCfgId)+1] = &GuildSkill{
			UUid:  int32(SkillCfgId + 1),
			CfgId: int32(guildCfg.SkillID),
			Exp:   remainingExp,
			Level: int32(guildCfg.SkillLevel + 1),
		}

		// 删除当前等级技能
		delete(guildInfo.Skills, int32(SkillCfgId))
	} else {
		// 未达到升级条件，只增加经验
		guildInfo.Skills[int32(SkillCfgId)].Exp += int32(exp)
	}
	// 更新行会信息
	// 转成json
	skillsJson, err := json.Marshal(guildInfo.Skills)
	if err != nil {
		return err
	}
	guildInfo.Guild.Skill = string(skillsJson)
	if err := db.UpdateGuild(guildInfo.Guild); err != nil {
		return err
	}
	return nil
}

// 获取等级为1的行会技能（初始化行会技能列表）
func (gm *GuildManager) GetGuildSkillBegin() map[int32]*GuildSkill {
	skills := make(map[int32]*GuildSkill)
	for _, skill := range global.GuildSkillCfg {
		if skill.SkillLevel == 1 {
			skills[skill.UUID] = &GuildSkill{
				UUid:  skill.UUID,
				CfgId: skill.SkillID,
				Level: 1,
			}
		}
	}
	return skills
}

// 获取某个人的信息
func (gm *GuildManager) getUserInfo(uid int64) *Player {
	// 获取用户信息
	// 检查玩家是否在线 defaultManager.players(uid)
	if p, ok := defaultManager.player(uid); ok {
		p.refreshStats()
		return p
	}
	p := newPlayer(nil, uid, "", "", "")
	if !p.syncFromDB() {
		return nil
	}
	p.curState = PlayerStatusOffLine
	p.refreshStats()
	defaultManager.setPlayer(uid, p)
	// 获取用户信息
	return p
}

// ///////////////// GuildInfo /////////////////////////
func (g *GuildInfo) SetGuidBattleId(battleId int32) {
	g.BattleId = battleId
}

func (g *GuildInfo) GetGuidBattleId() int32 {
	return g.BattleId
}

// 会长登录，给全公会人发走马灯
func (gm *GuildManager) GuildLoginPush(uid int64) {
	// 这里需要从数据库获取
	Info := gm.getUserInfo(uid)
	if Info == nil {
		return
	}
	if Info.guildId > 0 {
		// 查询行会信息
		if gm.guilds[Info.guildId] != nil {
			if gm.guilds[Info.guildId].Guild.LeaderID == uid {
				strSend := []string{Info.name}
				for member := range gm.guilds[Info.guildId].Members {
					Info := gm.getUserInfo(member)
					if Info.curState != PlayerStatusOffLine {
						packet := &pb.LampPush{
							Uid:     int32(constant.LampGuildMasterLogin),
							Context: strSend,
						}
						Info.pushMsg(protocol.RouteLampPush, packet)
					}
				}
			}
		}
	}
}
