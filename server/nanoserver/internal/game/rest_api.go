package game

import (
	"time"

	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
)

// TODO: conc
func Kick(uid int64) error {
	defaultManager.chKick <- uid
	return nil
}

func BroadcastSystemMessage(message string, msgType pb.ChatType) {
	//defaultManager.group.Broadcast(protocol.RouteChatPush, &pb.ChatPush{Type: msgType, Context: message})
	playerMgr.BroadcastToAllGroups(protocol.RouteChatPush, &pb.ChatPush{Type: msgType, Context: message})
}

func Reset(uid int64) {
	defaultManager.chReset <- uid
}

func Recharge(uid, coin int64) {
	// defaultManager.chRecharge <- protocol.RechargeInfo{Uid: uid, Coin: coin}
}

// 发放道具通用接口
func SourceAdd(source protocol.AddSource) {

	defaultManager.AddSourceBatch(source)
}

// 检查道具数量
func CheckSourceNum(source *protocol.AddSource) error {
	return defaultManager.CheckSourceNum(source)
}

// 检测能否发送道具
func CheckSourceAdd(source *protocol.AddSource) error {
	return defaultManager.CheckSourceAdd(source)
}

// 增加特权 privilegeId:配置PID  expirationTime:过期时间
func PrivilegeAdd(shopType int32, uid int64, privilegeId int32, expirationTime time.Time) error {
	return defaultManager.PrivilegeAdd(shopType, uid, privilegeId, expirationTime)
}

// 获取用户基础数据
func UserInfo(uid int64) (*protocol.UserInfo, error) {
	return defaultManager.GetPlayerInfo(uid)
}

// 更新玩家新手引导
func SetGuidance(uid int64, setp int32) error {
	return defaultManager.SetGuidance(uid, setp)
}

// 更新玩家命格
func SetPlayInfo(uid int64, fateLevel int32, fateExp int64) error {
	return defaultManager.SetPlayerInfo(uid, fateLevel, fateExp)
}

// 使用命格石头
func UseFate(uid int64, cfgId int32, num int32) ([]int, int32, error) {
	return defaultManager.userFate(uid, cfgId, num)
}

// 系统通知
func Notification(uid int64, Type pb.Notification) {
	defaultManager.Notification(uid, Type)
}

// 机器人对战测试
func CombatSetting(v *protocol.CombatSetting) error {
	return combatManagerr.Init(v)
}

func CombatResult() *protocol.CombatResult {
	return combatManagerr.CombatResult()
}

// 跑马灯
func Lamp(id int32, serverId int32, context []string) {
	defaultManager.Lamp(id, serverId, context)
}

// 跑马灯
func SystemLamp(stype int32, svalue int64, context string) {
	defaultManager.SystemLamp(stype, svalue, context)
}

// GM指令
func GMCommand(cmd protocol.GmCmd) error {
	if len(defaultManager.chGmCmd) > 32 {
		return errutil.ErrRequestFailed
	}
	defaultManager.chGmCmd <- cmd
	return nil
}

// 增加伙伴
func AddPartner(source *protocol.AddSource) error {
	return defaultManager.addPartnerBatch(source)
}

// 增加宠物
func AddPet(source *protocol.AddSource, NeedDel int32, reason string) error {
	return defaultManager.addPetBatch(source, NeedDel, reason)
}

// 工会登录
func GuildLogin(uid int64) {
	guildManager.GuildLoginPush(uid)
}

// 排行榜第一名登录
//
//	func RankLogin(uid int64) {
//		defaultManager.RankLogin(uid)
//	}
//
// 上报
func AddUpRecord(uid int64, name string, num int64) {
	defaultManager.AddUpRecord(uid, name, num)
}

// 上报
func Upshushu(uid int64, date map[string]interface{}, have bool, first bool) {
	defaultManager.UpShushuDate(uid, date, have, first)
}

func UpshushuSetOnce(uid int64, date map[string]interface{}, have bool) {
	defaultManager.UpShushuDateSetOnce(uid, date, have)
}

func UpshushuSet(uid int64, date map[string]interface{}, have bool) {
	defaultManager.UpShushuDateSet(uid, date, have)
}

// 上报
func Upshushuno(date map[string]interface{}) {
	defaultManager.UpShushuDateNo(date)
}
func Uptencent(date map[string]interface{}) {
	defaultManager.UpShushuTencent(date)
}

// 获取背包内道具数量
func GetBagNum(uid int64, configId int32) int32 {
	count, _ := defaultManager.GetItemCount(uid, configId)
	return int32(count)
}

// // 排行榜登录
// func RankLogin(uid int64) {
// 	provide.GuildLoginPush(uid)
// }

// 修改邮件的计数器，操作为设置，增加，减少 MailCounterSet
func SetMailCounter(uid int64, opType constant.MailCounterOPType, num int) error {
	return defaultManager.MailCounterSet(uid, opType, num)
}

// 同步任务
func TaskDispatch(uid int64, key pb.TargetType, Typeson int32, limit int32) {
	defaultManager.TaskDispatch(uid, key, Typeson, limit)
}

// GuardIslandMultiplier 设置守护岛倍数
func GuardIslandMultiplier(id, guardIslandRecordID int64) error {
	return defaultManager.GuardIslandMultiplierSvc(id, guardIslandRecordID)
}

func ChangeRoleNameREST(uid int64, name string) (string, error) {
	return defaultManager.ChangeRoleNameREST(uid, name)
}

func GuildCreateREST(uid int64, req *pb.GuildCreateRequest) (*pb.GuildCreateResponse, error) {
	res, err := defaultManager.GuildCreateREST(uid, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ChatREST(uid int64, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	res, err := defaultManager.ChatREST(uid, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func AddAdCount(uid int64) (int32, error) {
	return defaultManager.AddAdCount(uid)
}

func SetTempPurchasedSlot(uid int64, characterSlot, equipSlot, refinedSlot int32) error {
	return defaultManager.SetTempPurchasedSlot(uid, characterSlot, equipSlot, refinedSlot)
}

func RefineIsPurchased(uid int64, characterSlot, equipSlot, refinedSlot int32) (bool, error) {
	return defaultManager.RefineIsPurchased(uid, characterSlot, equipSlot, refinedSlot)
}
