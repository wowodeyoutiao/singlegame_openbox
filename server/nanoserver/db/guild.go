package db

import (
	"fmt"
	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
	"time"
)

// QueryGuild 根据工会ID查询工会信息
func QueryGuild(guildId int64) (*model.Guild, error) {
	guild := &model.Guild{GuildID: guildId}
	has, err := database.Get(guild)
	if err != nil {
		return nil, fmt.Errorf("查询工会信息失败: %v", err)
	}
	if !has {
		return nil, fmt.Errorf("工会不存在")
	}
	return guild, nil
}

// QueryGuildList 查询工会列表，顺便查一下有多少人在工会里
func QueryGuildList() ([]*model.Guild, error) {
	guilds := make([]*model.Guild, 0)
	err := database.Where("status = ?", StatusNormal).Find(&guilds)
	if err != nil {
		return nil, fmt.Errorf("查询工会列表失败: %v", err)
	}
	return guilds, nil
}

type GuildMemberCount struct {
	GuildID     int64 `xorm:"guild_id"`
	MemberCount int32 `xorm:"member_count"`
}

// 用sql group by分类查询每个公会的人员数量
func QueryGuildMemberCount() ([]*GuildMemberCount, error) {
	guilds := []*GuildMemberCount{}
	err := database.SQL("SELECT guild_id, COUNT(*) AS member_count FROM guild_member WHERE status = 1 GROUP BY guild_id").Find(&guilds) // 这里的leader_id是工会的创建者
	if err != nil {
		return nil, fmt.Errorf("查询工会成员数量失败: %v", err)
	}
	return guilds, nil
}

// InsertGuild 创建新工会
func InsertGuild(guild *model.Guild) error {
	guild.CreatedAt = time.Now().Unix()
	_, err := database.Insert(guild)
	if err != nil {
		return fmt.Errorf("创建工会失败: %v", err)
	}
	return nil
}

// UpdateGuild 更新工会信息
func UpdateGuild(guild *model.Guild) error {
	_, err := database.Where("guild_id = ?", guild.GuildID).Update(guild) // database.ID(guild.GuildID).Update(guild)
	if err != nil {
		return fmt.Errorf("更新工会信息失败: %v", err)
	}
	return nil
}

// DeleteGuild 删除工会
func DeleteGuild(guildId int64) error {
	_, err := database.Where("guild_id = ?", guildId).Update(&model.Guild{Status: StatusDeleted})
	if err != nil {
		return fmt.Errorf("删除工会失败: %v", err)
	}
	return nil
}

// QueryGuildMember 查询某人所在的工会
func QueryGuildMember(userId int64, status int8) (*model.GuildMember, error) {
	member := &model.GuildMember{UserID: userId, Status: status}
	has, err := database.Get(member)
	if err != nil {
		return nil, fmt.Errorf("查询工会成员失败: %v", err)
	}
	if !has {
		return nil, nil
	}
	return member, nil
}

// QueryGuildMembers 查询工会所有成员
func QueryGuildMembers(guildId int64) ([]*model.GuildMember, error) {
	members := make([]*model.GuildMember, 0)
	err := database.Where("guild_id = ? AND status = 1", guildId).Find(&members)
	if err != nil {
		return nil, fmt.Errorf("查询工会成员列表失败: %v", err)
	}
	return members, nil
}

// InsertGuildMember 添加工会成员
func InsertGuildMember(member *model.GuildMember) error {
	member.JoinedAt = time.Now().Unix()
	member.Status = 1
	_, err := database.Insert(member)
	if err != nil {
		return fmt.Errorf("添加工会成员失败: %v", err)
	}
	return nil
}

// UpdateGuildMember 更新工会成员信息
func UpdateGuildMember(member *model.GuildMember) error {
	_, err := database.Where("guild_id = ? AND user_id = ? and status = 1", member.GuildID, member.UserID).Update(member)
	if err != nil {
		return fmt.Errorf("更新工会成员失败: %v", err)
	}
	return nil
}

// DeleteGuildMember 删除工会成员
func DeleteGuildMember(memberId int64) error {
	_, err := database.ID(memberId).Update(&model.GuildMember{Status: 3})
	if err != nil {
		return fmt.Errorf("删除工会成员失败: %v", err)
	}
	return nil
}

// QueryGuildApplications 查询工会申请列表
func QueryGuildApplications(guildId int64) ([]*model.GuildApplication, error) {
	applications := make([]*model.GuildApplication, 0)
	err := database.Where("guild_id = ? AND status = 1", guildId).Find(&applications)
	if err != nil {
		return nil, fmt.Errorf("查询工会申请列表失败: %v", err)
	}
	return applications, nil
}

// QueryGuildApplications 查询工会所有申请列表
func QueryGuildApplicationsAll() ([]*model.GuildApplication, error) {
	applications := make([]*model.GuildApplication, 0)
	err := database.Where("status = 1").Find(&applications)
	if err != nil {
		return nil, fmt.Errorf("查询工会申请列表失败: %v", err)
	}
	return applications, nil
}

// InsertGuildApplication 添加工会申请
func InsertGuildApplication(application *model.GuildApplication) error {
	application.AppliedAt = time.Now().Unix()
	application.Status = 1
	_, err := database.Insert(application)
	if err != nil {
		return fmt.Errorf("添加工会申请失败: %v", err)
	}
	return nil
}

// UpdateGuildApplication 更新工会申请状态
func UpdateGuildApplication(application *model.GuildApplication) error {
	application.HandledAt = time.Now().Unix()
	_, err := database.ID(application.ApplicationID).Update(application)
	if err != nil {
		return fmt.Errorf("更新工会申请失败: %v", err)
	}
	return nil
}

// InsertGuildLog 添加工会操作日志
func InsertGuildLog(log *model.GuildLog) error {
	log.ActionTime = time.Now().Unix()
	_, err := database.Insert(log)
	if err != nil {
		return fmt.Errorf("添加工会日志失败: %v", err)
	}
	return nil
}

// QueryGuildLogs 查询工会操作日志
func QueryGuildLogs(guildId int64) ([]*model.GuildLog, error) {
	logs := make([]*model.GuildLog, 0)
	err := database.Where("guild_id = ?", guildId).Desc("action_time").Find(&logs)
	if err != nil {
		return nil, fmt.Errorf("查询工会日志失败: %v", err)
	}
	return logs, nil
}

// QueryGuildShopBuyInfo 查询公会商城全部或者某个商品，全员或者某个个人今日购买次数
// QueryGuildShopBuyInfoToday 查询公会商城今日购买次数
type BuyStore struct {
	GoodsID  int64 `xorm:"goods_id"`
	UserNum  int64 `xorm:"user_num"`
	GuildNum int64 `xorm:"guild_num"`
}

func QueryGuildShopBuyInfo(dayType int32, guildId int64, goodsIds []int64, uid int64) ([]BuyStore, error) {
	storeAll := []*model.GuildStoreBuy{}
	storeUser := []*model.GuildStoreBuy{}
	storeUserDel := make(map[int64]int32)
	storeOther := []BuyStore{}
	now := time.Now()                                                                                         // 获取当前时间
	tomorrow := now.AddDate(0, 0, 1)                                                                          // 增加一天
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()) // 设置为0点0分0秒
	timea := midnight.Unix() - int64((dayType)*3600*24)

	err := database.Select("goods_id,sum(buy_num) as buy_num ,min(buy_time) buy_time").Where("status = ? and buy_time>? and guild_id=? ",
		StatusNormal, timea, guildId).In("goods_id", goodsIds).GroupBy("goods_id").
		OrderBy("buy_time").Find(&storeAll)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	err = database.Select("goods_id,sum(buy_num) as buy_num,min(buy_time) buy_time").Where("status = ? and buy_time>? and guild_id=? ",
		StatusNormal, timea, guildId).In("user_id", uid).In("goods_id", goodsIds).GroupBy("goods_id").
		OrderBy("buy_time").Find(&storeUser)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	for _, v := range storeUser {
		storeUserDel[v.GoodsID] = v.BuyNum
	}
	for _, v := range storeAll {
		userNum := int32(0)
		if storeUserDel[v.GoodsID] > 0 {
			userNum = storeUserDel[v.GoodsID]
		}
		storeOther = append(storeOther, BuyStore{
			GoodsID:  v.GoodsID,
			UserNum:  int64(userNum),
			GuildNum: int64(v.BuyNum),
		})
	}
	return storeOther, nil
}

// InsertGuildShopBuyInfo 插入公会商城购买记录
func InsertGuildShopBuyInfo(buyInfo *model.GuildStoreBuy) error {
	_, err := database.Insert(buyInfo)
	if err != nil {
		return fmt.Errorf("插入公会商城购买记录失败: %v", err)
	}
	return nil
}

// QueryGuildDonateInfo 查询用户当天的捐献记录
func QueryGuildDonateInfo(guildId int64, uid int64) ([]*model.GuildDonate, error) {
	donates := make([]*model.GuildDonate, 0)
	//获取当天凌晨的时间戳
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	timea := midnight.Unix()
	err := database.Where("guild_id = ? AND donate_time >= ? and donator_id = ?", guildId, timea, uid).Find(&donates)
	if err != nil {
		return nil, fmt.Errorf("查询公会捐献记录失败: %v", err)
	}
	return donates, nil
}

// InsertGuildDonateInfo 插入公会捐献记录
func InsertGuildDonateInfo(donate *model.GuildDonate) error {
	_, err := database.Insert(donate)
	if err != nil {
		return fmt.Errorf("插入公会捐献记录失败: %v", err)
	}
	return nil
}

func UpdateUserGuild(p *model.GuildMember) error {
	if p == nil {
		return nil
	}
	_, err := database.Where("user_id=? and guild_id=? and status=1", p.UserID, p.GuildID).Update(p)
	return err
}
