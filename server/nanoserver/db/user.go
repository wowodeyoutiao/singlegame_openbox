package db

import (
	"fmt"
	"strconv"
	"time"

	"nanoserver/db/model"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"

	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

// QueryUser get the user by id
func QueryUser(uid int64) (*model.User, error) {
	if uid <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	u := &model.User{
		Uid: uid,
	}
	has, err := database.Get(u)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return u, nil
}

// QueryUserByOpenID get the user by openID
func QueryPlayersByOpenID(openID string) (*model.User, error) {
	if openID == "" {
		return nil, fmt.Errorf("invalid parameter: openID is empty")
	}

	player := &model.User{
		AccountName: openID,
		IsOnline:    UserOnline,
	}
	has, err := database.Get(player)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errutil.ErrUserNotFound
	}

	return player, nil
}

// 查询某个服务器下的所有uid
func QueryUserIDsByServerID(serverID int64) ([]int64, error) {
	if serverID <= 0 {
		return nil, fmt.Errorf("invalid parameter: serverID is empty")
	}

	var userIDs []int64
	err := database.Table("user").Where("server_id = ?", serverID).Cols("uid").Find(&userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}
func QueryUserIDs() ([]int64, error) {

	var userIDs []int64
	err := database.Table("user").Cols("uid").Find(&userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

// UpdateUser update user's info
func UpdateUser(u *model.User) error {
	if u == nil {
		return nil
	}

	var errors []error

	if u.Bag == nil {
		if _, err := database.Where("uid=?", u.Uid).Cols("bag").Update(&model.User{Bag: nil}); err != nil {
			logger.Errorf("UpdateUser: update bag error: %s", err.Error())
			errors = append(errors, fmt.Errorf("update bag failed: %w", err))
		}
	}
	if u.Task == nil {
		if _, err := database.Where("uid=?", u.Uid).Cols("task").Update(&model.User{Task: nil}); err != nil {
			logger.Errorf("UpdateUser: update task error: %s", err.Error())
			errors = append(errors, fmt.Errorf("update task failed: %w", err))
		}
	}
	if u.Buff == nil {
		if _, err := database.Where("uid=?", u.Uid).Cols("buff").Update(&model.User{Buff: nil}); err != nil {
			logger.Errorf("UpdateUser: update buff error: %s", err.Error())
			errors = append(errors, fmt.Errorf("update buff failed: %w", err))
		}
	}

	if _, err := database.Where("uid=?", u.Uid).Cols("world_boss_damage,world_boss_times,match_times_duel,buy_match_num_duel,match_times_melee,buy_match_num_melee,daily_soul,session_soul,afk_multiplier_i_d,afk_daily_time").Update(u); err != nil {
		logger.Errorf("UpdateUser: update user error: %s", err.Error())
		errors = append(errors, fmt.Errorf("update user specific fields failed: %w", err))
	}

	if _, err := database.Where("uid=?", u.Uid).Update(u); err != nil {
		logger.Errorf("UpdateUser: update user error: %s", err.Error())
		errors = append(errors, fmt.Errorf("update user failed: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("UpdateUser failed with %d errors: %v", len(errors), errors)
	}

	return nil
}

// UpdateUserMuteUntil 更新用户的禁言到期时间
func UpdateUserMuteUntil(uid int64, muteUntil int64) error {
	if uid <= 0 {
		return errutil.ErrUserNotFound
	}
	u := &model.User{
		Uid:       uid,
		MuteUntil: muteUntil,
	}
	_, err := database.Where("uid=?", uid).Cols("mute_until").Update(u)
	if err != nil {
		logger.Errorf("UpdateUserMuteUntil: update muteUntil error: %s", err.Error())
		return fmt.Errorf("update muteUntil failed: %w", err)
	}
	return nil
}

// UpdateRegisterBanUntil 更新账号的封号到期时间（封号是针对账号的，不是针对角色的）
func UpdateRegisterBanUntil(uid int64, banUntil int64) error {
	if uid <= 0 {
		return errutil.ErrUserNotFound
	}

	// 1. 通过 UID 查询 User 表获取 account_name（即 IMEI）
	user, err := QueryUser(uid)
	if err != nil {
		return fmt.Errorf("query user failed: %w", err)
	}

	// 2. 通过 account_name（IMEI）更新 Register 表中所有匹配的记录
	reg := &model.Register{
		BanUntil: banUntil,
	}
	_, err = database.Where("imei=?", user.AccountName).Cols("ban_until").Update(reg)
	if err != nil {
		logger.Errorf("UpdateRegisterBanUntil: update register banUntil error: %s", err.Error())
		return fmt.Errorf("update register banUntil failed: %w", err)
	}
	return nil
}

// // UpdateUser update user's info
// func UpdateUser(u *model.User) error {
// 	if u == nil {
// 		return nil
// 	}
// 	if u.Bag == nil {
// 		database.Where("uid=?", u.Uid).Cols("bag").Update(&model.User{Bag: nil})
// 	}
// 	if u.Task == nil {
// 		database.Where("uid=?", u.Uid).Cols("task").Update(&model.User{Task: nil})
// 	}
// 	if u.Buff == nil {
// 		database.Where("uid=?", u.Uid).Cols("buff").Update(&model.User{Buff: nil})
// 	}
// 	database.Where("uid=?", u.Uid).Cols("world_boss_damage,world_boss_times,match_times_duel,buy_match_num_duel,match_times_melee,buy_match_num_melee,daily_soul,session_soul").Update(u)

// 	_, err := database.Where("uid=?", u.Uid).Update(u)
// 	return err
// }

// InsertUser insert a new user
func InsertUser(u *model.User) error {
	if u == nil {
		return nil
	}
	_, err := database.Insert(u)
	return err
}

// DeleteUser delete the user (软删除)
func DeleteUser(uid int64) error {
	u := model.User{
		Status: StatusDeleted,
	}
	_, err := database.Where("uid=?", uid).Update(u)
	return err
}

// DeleteOrRestoreUser 删除/恢复用户（软删除，支持恢复）
// action: 1=删除（软删除），2=恢复
func DeleteOrRestoreUser(uid int64, action int32) error {
	var status int
	if action == 2 {
		status = StatusNormal // 恢复
	} else {
		status = StatusDeleted // 删除
	}
	u := model.User{
		Status: status,
	}
	_, err := database.Where("uid=?", uid).Cols("status").Update(u)
	if err != nil {
		logger.Errorf("DeleteOrRestoreUser: update user status error: %s", err.Error())
		return fmt.Errorf("update user status failed: %w", err)
	}
	return nil
}

func UserAddCoin(uid int64, coin int64) error {
	session := database.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return errutil.ErrDBOperation
	}
	u := &model.User{Id: uid}
	has, err := session.Get(u)
	if err != nil {
		return err
	}
	if !has {
		return errutil.ErrNotFound
	}
	u.Coin += coin
	_, err = session.Where("uid=?", uid).Update(u)
	if err != nil {
		session.Rollback()
		return err
	}
	return session.Commit()
}

func UserLoseCoin(id int64, coin int64) error {
	user := &model.User{
		Id: id,
	}

	has, err := database.Get(user)
	if !has {
		return errutil.ErrNotFound
	}
	if err != nil {
		return err
	}

	return UserLoseCoinByUID(user.Id, coin)
}

func UserLoseCoinByUID(uid int64, coin int64) error {
	session := database.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		logger.Error(err.Error())
		return errutil.ErrDBOperation
	}
	u := &model.User{Id: uid}
	has, err := session.Get(u)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !has {
		logger.Error("user not found")
		return errutil.ErrNotFound
	}
	if u.Coin < coin {
		logger.Error("coin not enough")
		return errutil.ErrCoinNotEnough
	}
	u.Coin -= coin
	// FIXED: 用户剩余1的时候, 扣除不成功
	if _, err := session.Cols("coin").Where("uid=?", uid).Update(u); err != nil {
		session.Rollback()
		logger.Error(err.Error())
		return err
	}

	return session.Commit()
}

func InsertRegister(reg *model.Register) {
	chWrite <- reg
}

func Online(uid int64) error {
	u := &model.User{IsOnline: UserOnline, LastLoginAt: time.Now().Unix()}
	if _, err := database.Where("uid=?", uid).Update(u); err != nil {
		return err
	}
	return nil
}

func Offline(uid int64) error {
	u := &model.User{IsOnline: UserOffline, LastLoginAt: time.Now().Unix()}
	if _, err := database.Where("uid=?", uid).Update(u); err != nil {
		return err
	}
	return nil
}

// 查询是否注册
func QueryRegister(appId string, imei string, serverId int32) (*model.Register, error) {
	bean := &model.Register{
		Imei:     imei,
		AppId:    appId,
		ServerId: serverId,
	}

	ok, err := database.Get(bean)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errutil.ErrUserNotFound
	}

	return bean, nil
}

func QueryRegisterV2(session *xorm.Session, appId string, imei string, serverId int32) (*model.Register, error) {
	bean := &model.Register{
		Imei:     imei,
		AppId:    appId,
		ServerId: serverId,
	}

	ok, err := session.Get(bean)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errutil.ErrUserNotFound
	}

	return bean, nil
}

// QueryRegisterByUid 通过 UID 查询注册信息（用于封号检查）
func QueryRegisterByUid(uid int64) (*model.Register, error) {
	bean := &model.Register{
		Uid: uid,
	}
	has, err := database.Get(bean)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errutil.ErrUserNotFound
	}
	return bean, nil
}

func QueryUserListByAccountName(account string, serverId int32) ([]model.User, error) {
	user := &model.User{
		Status:      StatusNormal,
		AccountName: account,
		ServerId:    serverId,
	}
	_, err := database.Count(user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to count users with account %s and serverId %d", account, serverId)
	}
	result := make([]model.User, 0)
	err = database.Where("status = 1 and account_name=? and server_id=?", account, serverId).Limit(10, 0).Desc("id").Find(&result)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query users with account %s and serverId %d", account, serverId)
	}

	return result, nil
}

func QueryUserListByAccountNameV2(session *xorm.Session, account string, serverId int32) ([]model.User, error) {
	user := &model.User{
		Status:      StatusNormal,
		AccountName: account,
		ServerId:    serverId,
	}

	_, err := session.Count(user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to count users with account %s and serverId %d", account, serverId)
	}

	result := make([]model.User, 0)
	err = session.Where("status = 1 and account_name=? and server_id=?", account, serverId).Limit(10, 0).Desc("id").Find(&result)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query users with account %s and serverId %d", account, serverId)
	}

	return result, nil
}

// 用uid和角色名字查询是否存在
func UserCheck(uid int64, account string) ([]model.User, error) {
	User := &model.User{}
	result := make([]model.User, 0)
	err := database.Where("status = ? and uid=?  and account_name = ?", 1, uid, account).
		Desc("id").Find(&result, User)
	if err != nil {
		logger.Error(err)
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

func DelUserByAccountName(account string, uid int64) (bool, error) {
	u := &model.User{
		Status: StatusDeleted,
	}
	_, err := database.Where(" uid =? and account_name=? ", uid, account).Update(u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func QueryGuestUser(appId string, imei string) (*model.User, error) {
	bean := &model.Register{
		Imei:  imei,
		AppId: appId,
	}

	ok, err := database.Get(bean)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errutil.ErrUserNotFound
	}

	user := &model.User{
		Id: bean.Uid,
	}
	ok, err = database.Get(user)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errutil.ErrUserNotFound
	}

	return user, nil
}

func RegisterUserLog(uid int64, d protocol.Device, appid string, channelID string, regType int, serverId int32) {
	// Insert user register record
	t := time.Now().Unix()
	reg := &model.Register{
		Uid:          uid,
		Remote:       d.Remote,
		Ip:           d.IP,
		Imei:         d.IMEI,
		Os:           d.OS,
		Model:        d.Model,
		AppId:        appid,
		ChannelId:    channelID,
		RegisterAt:   t,
		RegisterType: regType,
		ServerId:     serverId,
	}
	InsertRegister(reg)
}

func RegisterUserLogV2(session *xorm.Session, uid int64, d protocol.Device, appid string, channelID string, regType int, serverId int32, adSource string) error {
	reg := &model.Register{
		Uid:          uid,
		Remote:       d.Remote,
		Ip:           d.IP,
		Imei:         d.IMEI,
		Os:           d.OS,
		Model:        d.Model,
		AppId:        appid,
		ChannelId:    channelID,
		RegisterAt:   time.Now().Unix(),
		RegisterType: regType,
		ServerId:     serverId,
		AdSource:     adSource,
	}
	if _, err := session.Insert(reg); err != nil {
		return fmt.Errorf("failed to insert register log: %w", err)
	}
	return nil
}

func InsertLoginLog(uid int64, d protocol.Device, appid string, channelID string) {
	// Insert user operation record
	log := &model.Login{
		Uid:       uid,
		Remote:    d.Remote,
		Ip:        d.IP,
		Imei:      d.IMEI,
		Os:        d.OS,
		Model:     d.Model,
		AppId:     appid,
		ChannelId: channelID,
		LoginAt:   time.Now().Unix(),
	}
	// userOnline(uid)
	chWrite <- log
}

func InsertRoleLoginLog(uid int64, d protocol.Device) {
	// Insert user operation record
	log := &model.LoginRole{
		Uid:     uid,
		Remote:  d.Remote,
		Ip:      d.IP,
		Imei:    d.IMEI,
		Os:      d.OS,
		Model:   d.Model,
		LoginAt: time.Now().Unix(),
	}
	chWrite <- log
}

// 查询某个用户一段时间内的所有登录情况
func RoleLoginNum(uid int64, begin_at int64, end_at int64) ([]protocol.LoginRoleNum, error) {
	sql := "SELECT  DATE(FROM_UNIXTIME(`login_at`)) AS date, COUNT(*) AS num " +
		"FROM login_role WHERE `uid` = ? and login_at>= ? and login_at<= ?  GROUP BY date ORDER BY date"
	mQuery, err := database.Query(
		sql,
		uid,
		begin_at,
		end_at,
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	LoginRoleNum := make([]protocol.LoginRoleNum, len(mQuery))

	for i, iteam := range mQuery {
		LoginRoleNum[i].Date = string(iteam["date"])
		LoginRoleNum[i].Num = string(iteam["num"])
	}

	return LoginRoleNum, nil
}

// 查询某个用户创角至今的天数
func RoleFirstLoginNum(uid int64) (int32, error) {
	sql := "SELECT MIN(login_at) AS first_login FROM login_role WHERE uid = ?"
	mQuery, err := database.Query(sql, uid)
	if err != nil || len(mQuery) == 0 {
		return 0, errutil.ErrDBOperation
	}
	firstLoginStr := string(mQuery[0]["first_login"])
	if firstLoginStr == "" {
		return 0, nil
	}
	firstLogin, err := strconv.ParseInt(firstLoginStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// 取首次登录和当前时间的日期（零点）
	firstDate := time.Unix(firstLogin, 0)
	now := time.Now()
	firstZero := time.Date(firstDate.Year(), firstDate.Month(), firstDate.Day(), 0, 0, 0, 0, firstDate.Location())
	nowZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	days := int32(nowZero.Sub(firstZero).Hours()/24) + 1
	return days, nil
}

// QueryUserInfo get the user by id
func QueryUserInfo(id int64) (*protocol.UserStatsInfo, error) {
	if id <= 0 {
		return nil, errutil.ErrUserNotFound
	}
	u := &model.User{
		Id: id,
	}
	has, err := database.Get(u)
	if !has {
		err = errutil.ErrUserNotFound
	}

	if err != nil {
		logger.Errorf("查询用户出错, Error%s", err.Error())
		return nil, err
	}

	// 注册记录
	r := &model.Register{
		Uid: u.Id,
	}
	database.Get(r)

	// 登录记录
	l := &model.Login{
		Uid: u.Id,
	}
	database.Desc("login_at").Get(l)

	ta := &model.ThirdAccount{
		Uid: u.Id,
	}
	database.Get(ta)

	// 总对局数
	match, _ := database.Where("player0 =? OR player1 =? OR player2 = ?",
		id, id, id).Count(model.Desk{})

	usi := &protocol.UserStatsInfo{
		ID:             u.Id,
		Uid:            u.Id,
		Name:           ta.ThirdName,
		RegisterAt:     r.RegisterAt,
		RegisterIP:     r.Remote,
		LastestLoginAt: l.LoginAt,
		LastestLoginIP: l.Remote,
		RemainCard:     u.Coin,
		TotalMatch:     match,
		Stats:          make(map[int64]*protocol.DailyStats),
		StatsAt:        []int64{},
	}

	// 最多查最近一月的数据
	now := time.Now()
	temp := time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, time.Local)

	begin := temp.Unix()
	if begin < r.RegisterAt {
		t1 := time.Unix(r.RegisterAt, 0)
		begin = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local).Unix()
	}

	temp2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	f := func(id, begin, end int64) *protocol.DailyStats {
		ret := &protocol.DailyStats{}
		// 桌数
		asCreator, _ := database.Where("creator = ? AND created_at BETWEEN ? AND ?",
			id,
			begin,
			end).Count(&model.Desk{})

		ret.AsCreator = asCreator

		// 参与过的房号
		desks := []model.Desk{}
		database.Where("(player0 =? OR player1 =? OR player2 = ? ) AND created_at BETWEEN ? AND ?",
			id, id, id,
			begin, end,
		).Find(&desks)

		// 胜局数
		wins := 0

		// 战绩
		score := 0
		for _, d := range desks {

			if d.Player0 == id {
				if d.ScoreChange0 > 0 {
					wins++
				}
				score += d.ScoreChange0

			}
			if d.Player1 == id {
				if d.ScoreChange0 > 0 {
					wins++
				}
				score += d.ScoreChange1
			}
			if d.Player2 == id {
				if d.ScoreChange0 > 0 {
					wins++
				}
				score += d.ScoreChange2
			}

			ret.DeskNos = append(ret.DeskNos, d.DeskNo)

		}
		ret.Score = score
		ret.Win = wins

		return ret
	}

	for i := begin; i <= temp2.Unix(); i += 86400 {
		usi.Stats[i] = f(id, i, i+86400)
		usi.StatsAt = append(usi.StatsAt, i)
	}

	return usi, nil
}

// IsUserExists is the user with uid is existed
func IsUserExists(uid int64) bool {
	u := &model.User{
		Status: StatusNormal,
		Id:     uid,
	}
	has, _ := database.Get(u)
	return has
}

func IsUserExistsByName(name string) bool {
	u := &model.User{
		Status: StatusNormal,
		Name:   name,
	}
	has, _ := database.Get(u)
	return has
}

func IsUserExistsByNameV2(session *xorm.Session, name string) (bool, error) {
	u := &model.User{
		Status: StatusNormal,
		Name:   name,
	}
	return session.Get(u)
}

func QueryUserList(offset, count int) ([]model.User, int64, error) {
	user := &model.User{
		Status: StatusNormal,
	}
	total, err := database.Count(user)
	if err != nil {
		logger.Error(err)
		return nil, 0, errutil.ErrDBOperation
	}
	result := make([]model.User, 0)
	err = database.Where("status=?", StatusNormal).Limit(count, offset).
		Desc("id").Find(&result)
	if err != nil {
		return nil, 0, errutil.ErrDBOperation
	}
	return result, total, nil
}

// 注册用户数
func QueryRegisterUsers(begin, end int64) (int, error) {
	if begin > end {
		return 0, errutil.ErrIllegalParameter
	}

	user := &model.User{
		Status: StatusNormal,
	}

	total, err := database.Where("`register_at` BETWEEN ? AND ?", begin, end).Count(user)
	if err != nil {
		logger.Error(err)
		return 0, errutil.ErrDBOperation
	}

	return int(total), nil
}

// 活跃人数
func QueryActivationUser(from, to int64) ([]*protocol.ActivationUser, error) {
	fn := func(from, to int64) *protocol.ActivationUser {
		mQuery, err := database.Query("SELECT COUNT(DISTINCT(uid)) AS users FROM `login` WHERE login_at BETWEEN ? AND ? ",
			from,
			to)

		cc := &protocol.ActivationUser{
			Date: from,
		}

		if len(mQuery) < 1 || err != nil {
			return cc
		}

		temp := string(mQuery[0]["users"])
		if temp != "" {
			cc.Value, _ = strconv.ParseInt(temp, 10, 0)
		}

		return cc
	}
	begin := time.Unix(from, 0)

	var ret []*protocol.ActivationUser

	t := time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)

	for i := t.Unix(); i < to; i += dayInSecond {
		cc := fn(i, i+dayInSecond-1)

		ret = append(ret, cc)
	}
	return ret, nil
}

// 实时在线人/桌数
func OnlineStatsLite() (*model.Online, error) {
	ol := &model.Online{}

	has, err := database.Desc("time").Get(ol)
	if err != nil || !has {
		return nil, err
	}

	return ol, err
}

func retentionHelper(current int) (*retentionStats, error) {
	f := func(step int) int64 {
		sql := fmt.Sprintf("SELECT COUNT( DISTINCT(login.uid)) AS retention FROM login JOIN register ON login.uid = register.uid	" +
			" WHERE register.register_at BETWEEN ? AND ? AND login.login_at BETWEEN ? AND ? ")

		fmt.Print(sql)
		fmt.Println(current, current+step)
		m, err := database.Query(
			sql,
			current,
			current+dayInSecond,
			current+step,
			current+step+dayInSecond)
		if len(m) < 1 || err != nil {
			return 0
		}

		retentionStr := string(m[0]["retention"])
		if retentionStr != "" {
			retention, err := strconv.ParseInt(retentionStr, 10, 0)
			if err != nil {
				return 0
			}

			return retention
		}

		return 0
	}

	r := &retentionStats{}

	r.loginDay1 = f(day1)
	r.loginDay2 = f(day2)
	r.loginDay3 = f(day3)
	r.loginDay7 = f(day7)
	r.loginDay14 = f(day14)
	r.loginDay30 = f(day30)

	sql := "SELECT COUNT(DISTINCT(uid)) AS register FROM `register` WHERE register_at BETWEEN ? AND ? "

	mQuery, err := database.Query(
		sql,
		current,
		current+dayInSecond,
	)

	if len(mQuery) < 1 || err != nil {
		return nil, err
	}

	registerStr := string(mQuery[0]["register"])
	if registerStr != "" {
		r.register, _ = strconv.ParseInt(registerStr, 10, 0)
	}

	return r, nil
}

// 某注册天的 n日留存
func RetentionList(current int) (*protocol.Retention, error) {
	st, err := retentionHelper(current)
	if err != nil {
		return nil, err
	}

	fill := func(rl *protocol.RetentionLite, register, login int64) {
		rl.Login = login

		if register != 0 {
			rl.Rate = fmt.Sprintf("%.2f", float32(login*100)/float32(register))
		} else {
			rl.Rate = "0.00"
		}
	}

	ret := &protocol.Retention{
		Date:     current,
		Register: st.register,
	}

	fill(&ret.Retention_1, st.register, st.loginDay1)
	fill(&ret.Retention_2, st.register, st.loginDay2)
	fill(&ret.Retention_3, st.register, st.loginDay3)
	fill(&ret.Retention_7, st.register, st.loginDay7)
	fill(&ret.Retention_14, st.register, st.loginDay14)
	fill(&ret.Retention_30, st.register, st.loginDay30)

	return ret, nil
}

// 查询所有机器人
func QueryRobotUserList() ([]model.User, error) {
	result := make([]model.User, 0)
	err := database.Where("is_robot=1").Find(&result)
	if err != nil {
		return nil, errutil.ErrDBOperation
	}
	return result, nil
}

// 0点清理
func DeleteUserItem() error {
	// 每天0点 定时任务开始执行...
	u := &model.User{
		DayOnlineTime: 0,
		VitalityTime:  0,
	}
	_, err := database.Update(u)
	if err != nil {
		return errutil.ErrDBOperation
	}
	return nil
}

// 查询单个服务器创建角色的数量
func QueryUserCountByServerId() ([]protocol.RegisterUserNum, error) {
	sql := "SELECT server_id, COUNT(*) AS num FROM `user` GROUP BY server_id"
	mQuery, err := database.Query(
		sql,
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	RegisterUserNum := make([]protocol.RegisterUserNum, len(mQuery))

	for i, iteam := range mQuery {
		RegisterUserNum[i].ServerId, err = strconv.Atoi(string(iteam["server_id"]))
		RegisterUserNum[i].Num, err = strconv.ParseInt(string(iteam["num"]), 10, 64)
	}

	return RegisterUserNum, nil
}
