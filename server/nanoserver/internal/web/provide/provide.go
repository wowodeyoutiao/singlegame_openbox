package provide

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"nanoserver/db"
	"nanoserver/global"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/maputil"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/gomodule/redigo/redis"
	"github.com/lonng/nano/serialize/protobuf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger = log.WithField("component", "game")
	config protocol.ClientConfig // 远程配置

)

type Rank struct {
	LampLock sync.RWMutex
	chLamp   chan *pb.LampPush
}

func PrivilegeListTCP(uid int64) ([]protocol.PrivilegeDetailTCP, error) {
	privilegeCfg := global.PrivilegeCfg
	// result := make(map[int32][]protocol.SourceInfo)
	result := make([]protocol.PrivilegeDetailTCP, 0)
	// 特权卡
	for _, privilegen := range privilegeCfg {
		endtime := time.Now()
		goodsId := make([]int64, 0)
		request2 := &protocol.StoreListRequest{}
		// request2.Uid = uid
		request2.StoreType = privilegen.Sid / 1000
		request2.GoodsUuid = privilegen.Sid
		// request.StoreType = r.StoreType
		// request.GoodsId = r.GoodsId
		// 查询用户购买情况
		goodsId = append(goodsId, int64(privilegen.Sid))
		resultb, _ := db.StoreList(
			uid,
			request2.StoreType,
			privilegen.Duration,
			goodsId,
			true,
		)
		if len(resultb) > 0 {
			originalTime := time.Unix(resultb[0].CreateTime, 0)
			nextDayMidnight := originalTime.AddDate(0, 0, int(privilegen.Duration)) // 添加N天
			nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
			// 转换回时间戳
			endtime = nextDayMidnight

			result = append(result, protocol.PrivilegeDetailTCP{
				Pid:     privilegen.Pid,
				EndTime: endtime,
			})
		}

	}
	return result, nil
}

// 是不是第一名
func IsFirst(Uid int64, Rtype int, serverId int32) []int {
	returnRankList := make([]int, 0)
	// rTypeName := map[int]string{
	// 	1: "等级", 2: "战力", 3: "乱斗", 4: "试炼", 6: "PVP", 7: "赛季", 20: "世界boss"}
	if Rtype == 0 {
		for k := range constant.RankTypeName {
			if constant.RankTypeServerNoLamp[int64(k)] {
				continue
			}
			// 查看第一名
			resultScores5, _ := zrevrangeTop1Raw(serverId, int32(k))
			resultsValue5 := resultScores5
			// 遍历结果，每两个元素是一个成员和它的分数
			for i := 0; i < len(resultsValue5); i += 2 {
				// 成员
				member, err := redis.Int64(resultsValue5[i], nil)
				if err != nil {
					// logger.Error("Error converting member to string:", err)
					continue
				}
				if strconv.FormatInt(member, 10) == strconv.FormatInt(Uid, 10) {
					// 发送消息给所有在线用户
					returnRankList = append(returnRankList, k)
				}

			}
			// if Firstscore == r.Value {

			// 	// strSend := []string{r.Name, rTypeName[int(r.Type)]}

			// 	// 发送消息给所有在线用户
			// 	// game.Lamp(int32(constant.LampTopRankChanged), strSend)
			// 	// 发送消息给所有在线用户
			// 	// desk :=
			// 	// desk.chLamp <- &pb.LampPush{
			// 	// 	Type: pb.LampType_LampTopRankChanged,
			// 	// 	Data: fmt.Sprintf("%s", r.Name),
			// 	// }
		}
	} else {
		if constant.RankTypeServerNoLamp[int64(Rtype)] {
			return returnRankList
		}
		// 查看第一名
		resultScores5, _ := zrevrangeTop1Raw(serverId, int32(Rtype))
		resultsValue5 := resultScores5
		// 遍历结果，每两个元素是一个成员和它的分数
		for i := 0; i < len(resultsValue5); i += 2 {
			// 成员
			member, err := redis.Int64(resultsValue5[i], nil)
			if err != nil {
				// logger.Error("Error converting member to string:", err)
				continue
			}
			if strconv.FormatInt(member, 10) == strconv.FormatInt(Uid, 10) {
				// 发送消息给所有在线用户
				returnRankList = append(returnRankList, Rtype)
			}
		}
	}
	return returnRankList
}

// 排行榜数据录入
func RankChangeTCP(r *protocol.RankUserInfoInput) (bool, error) {
	// 如果是世界boss，需要录入到排行榜
	rankMember := r.Uid
	switch r.Type {
	case constant.RankTypeWorldBossKf:
		// rankEnd := "99"

		rname := r.Name
		// 世界Boss跨服分支使用世界Boss专用用户键写入名称
		if err2 := hmsetWorldBossUserName(rankMember, rname, r.ServerId); err2 != nil {
			return true, err2
		}
		// key := constant.RedisRankWorldBossKey(rankMember)
		key := constant.RedisRankKey(r.ServerId, int32(r.Type))
		err := zincrbyRank(key, float64(r.Value), rankMember)
		if err != nil {
			return true, err
		}
		return true, nil
	case constant.RankTypeWorldBoss:
		if RankTypeServer, ok := constant.RankTypeServer[int64(r.Type)]; ok && RankTypeServer {
		} else {
			if r.ZoneId == 0 {
				r.ServerId = maputil.ABS(int32(r.ServerId))
				//不区分服务器的榜单，是需要区分战区的
				serverDetail, err := global.ServerListMgrInstance.Get(int(r.ServerId))
				if err != nil {
					return true, err
				}
				r.ServerId = -1 * maputil.ABS(int32(serverDetail.ZoneId))
			} else {
				r.ServerId = -1 * maputil.ABS(int32(r.ZoneId))
			}
		}
		// user fields read via helpers
		rname := r.Name
		ravatar := r.Avatar
		rsex := strconv.Itoa(int(r.Sex))
		rjob := strconv.Itoa(int(r.Job))
		rrobot_cid := strconv.Itoa(int(r.RobotCId))
		rserver_id := strconv.Itoa(int(r.UserServerId))
		ravatar2 := strconv.Itoa(int(r.AvatarV2))
		rframe := strconv.Itoa(int(r.Frame))
		key := constant.RedisRankKey(r.ServerId, int32(r.Type))
		// key2 := user + strconv.FormatInt(rankMember, 10)
		if err2 := hmsetUserFields(r.ServerId, rankMember, buildStandardUserFields(rname, ravatar, rsex, rjob, rrobot_cid, rserver_id, ravatar2, rframe)); err2 != nil {
			return true, err2
		}
		err := zincrbyRank(key, float64(r.Value), rankMember)
		if err != nil {
			return true, err
		}
		// 发送跨服消息
		RankChangeTCP(&protocol.RankUserInfoInput{
			Type:     constant.RankTypeWorldBossKf,
			Uid:      int64(r.RankServerId),
			Name:     r.ServerName,
			Value:    r.Value,
			Avatar:   r.Avatar,
			Sex:      r.Sex,
			Job:      r.Job,
			RobotCId: r.RobotCId,
			ServerId: r.ServerId,
		})
		return true, nil
	default:

		if RankTypeServer, ok := constant.RankTypeServer[int64(r.Type)]; ok && RankTypeServer {
		} else {
			if r.ZoneId == 0 {
				r.ServerId = maputil.ABS(int32(r.ServerId))
				//不区分服务器的榜单，是需要区分战区的
				serverDetail, err := global.ServerListMgrInstance.Get(int(r.ServerId))
				if err != nil {
					return true, err
				}
				r.ServerId = -1 * maputil.ABS(int32(serverDetail.ZoneId))
			} else {
				r.ServerId = -1 * maputil.ABS(int32(r.ZoneId))
			}

		}
		// rankEnd := "99"
		// user fields read via helpers
		rname := r.Name
		ravatar := r.Avatar
		rsex := strconv.Itoa(int(r.Sex))
		rjob := strconv.Itoa(int(r.Job))
		rrobot_cid := strconv.Itoa(int(r.RobotCId))
		rserver_id := strconv.Itoa(int(r.UserServerId))
		ravatar2 := strconv.Itoa(int(r.AvatarV2))
		rframe := strconv.Itoa(int(r.Frame))
		key := constant.RedisRankKey(r.ServerId, int32(r.Type))
		// 需要比较就查询比大小
		if r.Type == constant.RankTypeSeason {
			for _, v := range global.GlobalCfg.PVPSeason.Season {
				if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
					if time.Now().Before(start) {
						continue
					}
				} else {
					continue
				}
				if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
					if time.Now().After(end) {
						continue
					}
				}
				r.Type = r.Type*10000 + int32(v.ID)
				break
			}
		}

		if err2 := hmsetUserFields(r.ServerId, rankMember, buildStandardUserFields(rname, ravatar, rsex, rjob, rrobot_cid, rserver_id, ravatar2, rframe)); err2 != nil {
			return true, err2
		}
		if r.Compare {
			// 用户的分数
			resultScores, err5 := zscoreRaw(r.ServerId, r.Type, rankMember)
			if resultScores == 0 || err5 != nil {
				err := ZaddRank(key, float64(r.Value), rankMember)
				if err != nil {
					return true, err
				}
			} else {
				OwnScore := resultScores
				if OwnScore < r.Value {
					err := ZaddRank(key, float64(r.Value), rankMember)
					if err != nil {
						return true, err
					}
				}
			}

		} else {
			err := ZaddRank(key, float64(r.Value), rankMember)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}
}

// // 排行榜
// const (
// RankTypeLevel           = 1  // 等级榜
// RankTypeFight           = 2  // 战力榜
// RankTypeSoul            = 3  // PVP胜场数量
// RankTypeTrial           = 4  // 试炼榜
// RankTypePvp             = 6  // 每日9点的计算(当天战魂)
// RankTypeSeason          = 7  // 赛季排行榜
// RankTypeWorldBoss       = 20 // 世界boss
// RankTypeWorldBossKf     = 21 // 世界boss的跨服榜
// RankTypeGightEnd        = 22 // 一战到底总榜
// RankTypeGightEndlayers  = 23 // 一战到底每一层榜单
// RankTypeTurnBased       = 24 // 回合制排行榜
// RankTypeSoulWeekly      = 25 // pvp胜场周榜
// RankTypeTurnBasedWeekly = 26 // 回合制周榜
// )

// 一站到底
func CronSendRankAwardOne() {
	rtype := []int64{constant.RankTypeGightEnd}
	err := SendRankAward(rtype, false)
	logger.Debug("22发奖励了", rtype)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardOne", err)
	}
}

// 战力榜和即时赛榜，晚上10点发
func CronSendRankAwardPvp() {
	rtype := []int64{constant.RankTypeFight, constant.RankTypeSoul, constant.RankTypeTurnBased}
	logger.Debug("2,3,24发奖励了", rtype)
	err := SendRankAward(rtype, false)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardPvp", err)
	}
}

// 活动排行榜（31以后），发奖励
func CronSendRankAwardActivity() {
	rtype := []int64{
		constant.RankTypePinkEntry,
		constant.RankTypeRedPartnerStar,
		constant.RankTypeEquipmentStar,
		constant.RankTypePetSynthesis,
		constant.RankTypeSoulStoneLevel,
		constant.RankTypeFiveElementLevel,
		constant.RankTypeSimpleFormation,
	}
	logger.Debug("31-37发奖励了", rtype)
	err := SendRankAward(rtype, false)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardActivity", err)
	}
	if err := DelRank(rtype); err != nil {
		logger.Error("Failed to delete rank:", err)
	}
}

// 世界boss榜，晚上9点发
func CronSendRankAwardBoss() {
	rtype := []int64{constant.RankTypeWorldBossKf}
	logger.Debug("21发奖励了", rtype)
	err := SendRankAwardAll(rtype, false)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardBoss", err)
	}
	rtype2 := []int64{constant.RankTypeWorldBoss}
	logger.Debug("20发奖励了", rtype2)
	err = SendRankAward(rtype2, false)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardPvp", err)
	}
}

// 每周一0点发放周结奖励
func CronSendRankAwardWeekly() {
	rtype := []int64{constant.RankTypeFight, constant.RankTypeSoul, constant.RankTypeTurnBased}
	logger.Debug("2,3,24发奖励了", rtype)
	err := SendRankAward(rtype, true)
	if err != nil {
		// logger.Error("Rank err CronSendRankAwardBoss", err)
	}

	rankType := []int64{
		constant.RankTypeSoul,
		constant.RankTypeTurnBased,
		constant.RankTypeElement,
	}
	if err := DelRank(rankType); err != nil {
		logger.Error("Failed to delete rank:", err)
	}
	//机器人塞进榜
	// game.RankRest()
}

func SendRankAward(rType []int64, weekly bool) error {
	// rankAwardCfg := global.RankingAwardCfg
	// rTypeName := map[int64]string{
	// 	1: "等级", 2: "战力", 3: "乱斗", 4: "试炼", 6: "PVP", 7: "赛季", 20: "世界boss"}
	// result := make(map[int32][]protocol.SourceInfo)
	// 活动内容
	// activity := make(map[int32][]global.KeyValue)
	// for _, activityn := range rankAwardCfg {
	// 	activity[activityn.UUID] = activityn.Award
	// }

	// ServerIds, err := GetAllServerId()
	// if err != nil {
	// 	return err
	// }
	// if len(ServerIds) == 0 {
	// 	return nil
	// }
	ServerIds := make([]int32, 0)
	for i := int(global.ServerIDRangeStart); i <= int(global.ServerIDRangeEnd); i++ {
		ServerIds = append(ServerIds, int32(i))
	}
	//插一个0（已废弃，目前需要遍历所有战区id）

	// 如果你想用 zoneIds 替换 ServerIds，可以这样赋值
	// ServerIds = append(ServerIds, 0)
	ZoneIds := global.ServerListMgrInstance.GetAllZoneId()
	for _, v := range ZoneIds {
		ServerIds = append(ServerIds, -1*v)
	}
	//有些排行榜是限时排行榜，根据开服时间判断的

	for _, serverId := range ServerIds {
		// 查询所有服务器Id
		for _, item := range rType {
			//delete
			needDel := false
			//发奖励需要判断开服几天的榜，需要判断这个服务器的天数是不是对的
			if constant.RankTypeServerOpenDay[item] {
				fa := false
				openDay := global.ServerListMgrInstance.GetServerOpenDay(int32(serverId))
				//天数
				for _, v := range global.GlobalCfg.LimitRanking.Key {
					//实际天数与配置天数一样，需要发奖励
					if v.Value+1 == int64(openDay) && int64(v.ID) == item {
						fa = true
						needDel = true
						break
					}
				}
				if !fa {
					continue
				}
			}
			oldType := item
			NewType := item
			if weekly {
				if constant.RankTypeWeekly[item] != 0 {
					NewType = constant.RankTypeWeekly[item]
					logger.Error("NewType", NewType)
				}
			}
			if item == 7 {
				for _, v := range global.GlobalCfg.PVPSeason.Season {
					if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
						if time.Now().Before(start) {
							continue
						}
					} else {
						continue
					}
					if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
						if time.Now().After(end) {
							continue
						}
					}
					item = oldType*10000 + int64(v.ID)
					NewType = item + int64(v.ID)
					break
				}
			}
			rankBegin := "0"
			rankEnd := "99"
			RankAll := make([]protocol.RankUserInfo, 0)

			list, err2 := zrevrangeWithScores(int32(serverId), int32(item), rankBegin, rankEnd)
			if err2 != nil || list == nil {
				continue
			}
			RankAll = append(RankAll, list...)
			for _, item2 := range RankAll {
				// if activity[int32(NewType)*10000+int32(item2.Num)] == nil {
				// 	continue
				// }
				// result2 := activity[int32(NewType)*10000+int32(item2.Num)]
				// MailId := strconv.FormatInt(snowflake.GenID(), 10)
				result2 := global.GetRankAwardCfg(int32(NewType)*10000 + int32(item2.Num))
				result3 := make([]protocol.Attachment, len(result2.Award))
				// MailId := strconv.FormatInt(snowflake.GenID(), 10)
				for i, Attachment := range result2.Award {
					result3[i] = protocol.Attachment{
						MailId:      "",
						UserPropsId: int64(Attachment.Key),
						Num:         int64(Attachment.Value),
					}
				}
				// resultc, err3 := db.MailAdd(
				// 	-1,
				// 	item2.Uid,
				// 	"系统邮件：排行榜奖励",
				// 	constant.RankTypeName[int(oldType)]+"排行榜奖励",
				// 	result3,
				// 	"RankAward-"+strconv.Itoa(int(item)),
				// )
				shushumain.ReportMailData(protocol.MailAdd{
					Uid:          -1,
					RecipientUid: item2.Uid,
					Title:        "系统邮件：排行榜奖励",
					Desc:         constant.RankTypeName[int(NewType)] + "第" + strconv.Itoa(int(item2.Num)) + "名" + "排行榜奖励",
					Attachments:  result3,
					ChannelCode:  "RankAward-" + strconv.Itoa(int(item)),
				})
				if needDel == true {
					//删除排行榜
					err2 := delRankKey(serverId, int32(item))
					if err2 != nil {
						logger.Error("Failed to delete rank:", err2)
					}
				}
			}
		}
	}

	return nil
}

// 21的榜单是特殊的
// 因为21的榜单是世界boss的跨服，里面是服务器的排行
// 需要查询所有服务器的信息，然后给服务器发奖励
func SendRankAwardAll(rType []int64, weekly bool) error {
	// rankAwardCfg := global.RankingAwardCfg
	// rTypeName := map[int64]string{
	// 	1: "等级", 2: "战力", 3: "乱斗", 4: "试炼", 6: "PVP", 7: "赛季", 20: "世界boss"}
	// result := make(map[int32][]protocol.SourceInfo)
	// 活动内容
	// activity := make(map[int32][]global.KeyValue)
	// for _, activityn := range rankAwardCfg {
	// 	activity[activityn.UUID] = activityn.Award
	// }

	// ServerIds, err := GetAllServerId()
	// if err != nil {
	// 	return err
	// }
	// if len(ServerIds) == 0 {
	// 	return nil
	// }
	// 查询所有服务器Id
	//这个榜是战区榜单，所以需要遍历所有战区
	//查询所有的战区
	// 如果你想用 zoneIds 替换 ServerIds，可以这样赋值
	// ServerIds = append(ServerIds, 0)
	ZoneIds := make([]int32, 0)
	serverDetail := global.ServerListMgrInstance.GetAllZoneId()
	for _, v := range serverDetail {
		ZoneIds = append(ZoneIds, -1*v)
	}
	for _, item := range rType {
		for _, ZoneId := range ZoneIds {
			oldType := item
			NewType := item
			if weekly {
				if constant.RankTypeWeekly[item] != 0 {
					NewType = constant.RankTypeWeekly[item]
				}
			}
			// rank key constructed by helper when needed
			rankBegin := "0"
			rankEnd := "-1"
			RankAll := make([]protocol.RankUserInfo, 0)

			list, err2 := zrevrangeWithScores(ZoneId, int32(item), rankBegin, rankEnd)
			if err2 != nil || list == nil {
				continue
			}
			RankAll = append(RankAll, list...)
			//rankall是服务器的排行榜比如1区，2区，3区
			// 这里的item是一个服务器的
			result3 := make(map[int64][]protocol.Attachment, len(RankAll))
			for _, item2 := range RankAll {
				result2 := global.GetRankAwardCfg(int32(NewType)*10000 + int32(item2.Num))
				// MailId := strconv.FormatInt(snowflake.GenID(), 10)
				for _, Attachment := range result2.Award {
					result3[item2.Uid/100] = append(result3[item2.Uid/100], protocol.Attachment{
						MailId:      "",
						UserPropsId: int64(Attachment.Key),
						Num:         int64(Attachment.Value),
					})
				}
			}
			// 这里的item是一个服务器的

			// if activity[int32(NewType)*10000+int32(item2.Num)] == nil {
			// 	continue
			// }
			// result2 := activity[int32(NewType)*10000+int32(item2.Num)]
			// MailId := strconv.FormatInt(snowflake.GenID(), 10)

			//查询跨服榜单
			// cross-zone rank key constructed by helper when needed
			list2, err2 := zrevrangeWithScores(ZoneId, int32(constant.RankTypeWorldBoss), rankBegin, rankEnd)
			if err2 != nil || list2 == nil {
				continue
			}
			for i := 0; i < len(list2); i++ {
				member := list2[i].Uid
				//查询这个人在哪个服务器
				//然后对应的服务器的奖励
				// 用户的基本信息
				// 遍历结果，每两个元素是一个成员和它的分数
				InfoValues2, err4 := hgetallUserRaw(ZoneId, member)
				if err4 != nil {
					continue
					// logger.Error("InfoValues is not a slice of interfaces")
				}
				serverId := "server_id"
				rserverId := ""
				for j := 0; j < len(InfoValues2); j += 2 {
					member2, err6 := redis.String(InfoValues2[j], nil)
					if err6 != nil {
						// logger.Error("Error converting member to string:", err6)
						continue
					}
					values, err6 := redis.String(InfoValues2[j+1], nil)
					if err6 != nil {
						// logger.Error("Error converting score to String:", err6)
						continue
					}
					if member2 == serverId {
						rserverId = values
					}
				}
				//if err != nil {
				//	continue
				//}
				EndserverId, _ := strconv.Atoi(rserverId)

				//对应服务器的奖励
				if _, ok := result3[int64(EndserverId)]; !ok {
					continue
				}
				shushumain.ReportMailData(protocol.MailAdd{
					Uid:          -1,
					RecipientUid: member,
					Title:        "系统邮件：排行榜奖励",
					Desc:         constant.RankTypeName[int(oldType)] + "第" + strconv.Itoa(i+1) + "名" + "排行榜奖励",
					Attachments:  result3[int64(EndserverId)],
					ChannelCode:  "RankAward-" + strconv.Itoa(int(item)),
				})

			}
		}

		time.Sleep(time.Millisecond)
	}

	return nil
}

//	func DelRedis(rType []int64) error {
//		for _, item := range rType {
//			rank := "rank:"
//			rankBegin := "0"
//			rankEnd := "99"
//			results, err2 := CmdCache("zrevrange", rank+strconv.Itoa(int(item)), rankBegin, rankEnd, "WITHSCORES")
//			if err2 != nil {
//			}
//		}
//	}

// 查询所有服务器Id
func GetAllServerId() ([]int32, error) {
	respJson, err := RequestWeb(http.MethodPost, viper.GetString("admin.adminUrl")+"/api/v1/servername", BaseRequest{})
	if err != nil {
		return nil, err
	}
	var resp respParam
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		return nil, err
	}
	ServerId := make([]int32, 0)
	for _, server := range resp.Date.List {
		ServerId = append(ServerId, server.ServerId)
	}
	return ServerId, nil
}

// 根据具体返回的值填写
type respParam struct {
	ErrCode int32      `json:"code"`
	ErrMsg  string     `json:"error"`
	Date    ServerList `json:"data"`
}
type ServerList struct {
	Count int32        `json:"count"`
	List  []ServerInfo `json:"list"`
}
type ServerInfo struct {
	ServerId   int32  `json:"server_id"`
	ServerName string `json:"name"`
}

func DelRank(rType []int64) error {
	// 查询所有服务器Id
	// ServerIds, err := GetAllServerId()
	// if err != nil {
	// 	return err
	// }
	ServerIds := make([]int32, 0)
	for i := int(global.ServerIDRangeStart); i <= int(global.ServerIDRangeEnd); i++ {
		ServerIds = append(ServerIds, int32(i))
		// serverDetail, err := global.ServerListMgrInstance.Get(i)
		// if err != nil {
		// 	return err
		// }
		// ServerIds = append(ServerIds, serverDetail.ZoneId)

	}
	if len(ServerIds) == 0 {
		return nil
	}
	// ZoneIds := make([]int32, 0)
	serverDetail := global.ServerListMgrInstance.GetAllZoneId()
	for _, v := range serverDetail {
		ServerIds = append(ServerIds, -1*v)
	}
	// 遍历所有服务器Id
	// ServerIds = append(ServerIds, 0)
	for _, serverId := range ServerIds {
		for _, item := range rType {
			err2 := delRankKey(serverId, int32(item))
			if err2 != nil {
			}
		}
	}

	return nil
}

// 与开服时间相关的活动的排行榜
func DelayDelRank() error {
	for _, v := range global.GlobalCfg.LimitRanking.Key {
		rtype := []int64{int64(v.ID)}
		logger.Infof("自动脚本开始%d", v)
		err := SendRankAward(rtype, false)
		if err != nil {
			logger.Error("自动脚本报错", err)
		}
	}
	return nil
}
func RankList(r *protocol.RankOperateRequest) (protocol.RankDetail, error) {

	if RankTypeServer, ok := constant.RankTypeServer[int64(r.Type)]; ok && RankTypeServer {
	} else {
		//如果没有传入战区ID，那么需要传入服务器ID
		if r.ZoneId != 0 {
			r.ServerId = -1 * maputil.ABS(r.ZoneId)
		} else {
			r.ServerId = -1 * maputil.ABS(r.ServerId)
		}
	}

	RankDetail := protocol.RankDetail{}
	RankDetail.RankAll = make([]protocol.RankUserInfo, 0)

	if r.Type == constant.RankTypeWorldBossKf {
		rankBegin := "0"
		rankEnd := "99"
		if r.Count == 0 {
			rankEnd = strconv.Itoa(r.Count - 1)
		}
		Uid := int64(r.Uid)
		beginType := r.Type

		RankDetail.RankOwn = protocol.RankUserInfo{
			Uid:  Uid,
			Name: r.Name,
		}
		if r.Type == 7 {
			for _, v := range global.GlobalCfg.PVPSeason.Season {
				if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
					if time.Now().Before(start) {
						continue
					}
				} else {
					continue
				}
				if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
					if time.Now().After(end) {
						continue
					}
				}
				r.Type = r.Type*10000 + int64(v.ID)
				break
			}
		}
		if r.All == 1 || r.All == 2 {
			list, err2 := zrevrangeWithScores(int32(r.ServerId), int32(r.Type), rankBegin, rankEnd)
			if err2 != nil {
				return RankDetail, err2
			}
			for _, e := range list {
				fields, _ := hgetAllUserFields(int32(r.ServerId), e.Uid)
				RankDetail.RankAll = append(RankDetail.RankAll, protocol.RankUserInfo{
					Uid:    e.Uid,
					Num:    e.Num,
					Avatar: fields["avatar"],
					Name:   fields["name"],
					Value:  e.Value,
					Type:   beginType,
				})
			}
		}
		if r.All == 1 || r.All == 3 {
			// 个人的排名
			resultsValue2, err := zrevrankRaw(r.ServerId, int32(r.Type), Uid)
			if err != nil {
				return RankDetail, nil
			}
			// 用户的分数
			resultScores3, err := zscoreRaw(r.ServerId, int32(r.Type), Uid)
			if err != nil {
				return RankDetail, nil
			}
			OwnScore := resultScores3
			// 用户的基本信息
			fields2, _ := hgetAllUserFields(int32(r.ServerId), Uid)
			rname2 := fields2["name"]
			ravatar := fields2["avatar"]
			rsex2 := fields2["sex"]
			rjob2 := fields2["job"]
			//if err != nil {
			//	continue
			//}
			EndSex2, _ := strconv.Atoi(rsex2)
			EnJob2, _ := strconv.Atoi(rjob2)
			EnFrame2, _ := strconv.Atoi(fields2["frame"])
			EnAvatar2, _ := strconv.Atoi(fields2["avatar_v2"])
			EnServerId, _ := strconv.Atoi(fields2["server_id"])
			serverInfo, err := global.ServerListMgrInstance.Get(EnServerId)
			if err != nil {
				serverInfo.Name = ""
			}
			RankDetail.RankOwn = protocol.RankUserInfo{
				Uid:        Uid,
				Num:        resultsValue2 + 1,
				Avatar:     ravatar,
				Name:       rname2,
				Value:      OwnScore,
				Type:       beginType,
				Sex:        int32(EndSex2),
				Job:        int32(EnJob2),
				ServerName: serverInfo.Name,
				AvatarV2:   int32(EnAvatar2),
				Frame:      int32(EnFrame2),
			}
		}
	} else {
		rankBegin := "0"
		rankEnd := "99"
		if r.Count == 0 {
			rankEnd = strconv.Itoa(r.Count - 1)
		}
		Uid := r.Uid
		beginType := r.Type
		value := int64(0)
		if r.Type == constant.RankTypeSoul || r.Type == constant.RankTypeTurnBased {
			value = int64(global.GlobalCfg.BaseValue.PVPinitscore.Num)
		}
		RankDetail.RankOwn = protocol.RankUserInfo{
			Value: value,
			Uid:   Uid,
			Name:  r.Name,
		}
		if r.Type == 7 {
			for _, v := range global.GlobalCfg.PVPSeason.Season {
				if start, err := time.Parse(time.DateOnly, v.StartTime); err == nil {
					if time.Now().Before(start) {
						continue
					}
				} else {
					continue
				}
				if end, err := time.Parse(time.DateOnly, v.EndTime); err == nil {
					if time.Now().After(end) {
						continue
					}
				}
				r.Type = r.Type*10000 + int64(v.ID)
				break
			}
		}
		if r.All == 1 || r.All == 2 {
			list, err2 := zrevrangeWithScores(int32(r.ServerId), int32(r.Type), rankBegin, rankEnd)
			if err2 != nil {
				return RankDetail, nil
			}
			for _, e := range list {
				fields, _ := hgetAllUserFields(int32(r.ServerId), e.Uid)
				EndSex, _ := strconv.Atoi(fields["sex"])
				EnJob, _ := strconv.Atoi(fields["job"])
				EnRobotCId, _ := strconv.Atoi(fields["robot_cid"])
				EnFrame2, _ := strconv.Atoi(fields["frame"])
				EnAvatar2, _ := strconv.Atoi(fields["avatar_v2"])
				EnServerId, _ := strconv.Atoi(fields["server_id"])
				serverInfo, err := global.ServerListMgrInstance.Get(EnServerId)
				if err != nil {
					serverInfo.Name = ""
				}
				RankDetail.RankAll = append(RankDetail.RankAll, protocol.RankUserInfo{
					Uid:        e.Uid,
					Num:        e.Num,
					Avatar:     fields["avatar"],
					Name:       fields["name"],
					Sex:        int32(EndSex),
					Job:        int32(EnJob),
					RobotCId:   int32(EnRobotCId),
					Value:      e.Value,
					Type:       beginType,
					ServerName: serverInfo.Name,
					AvatarV2:   int32(EnAvatar2),
					Frame:      int32(EnFrame2),
				})
			}
		}
		if r.All == 1 || r.All == 3 {
			// 个人的排名
			resultsValue2, err := zrevrankRaw(r.ServerId, int32(r.Type), Uid)
			if err != nil {
				return RankDetail, nil
			}
			// resultsValue2 already int64 from zrevrankRaw
			// 用户的分数
			resultScores3, err := zscoreRaw(r.ServerId, int32(r.Type), Uid)
			if err != nil {
				return RankDetail, nil
			}
			OwnScore := resultScores3
			fields2, _ := hgetAllUserFields(int32(r.ServerId), Uid)
			rname2 := fields2["name"]
			ravatar2 := fields2["avatar"]
			rsex2 := fields2["sex"]
			rjob2 := fields2["job"]
			//if err != nil {
			//	continue
			//}
			EndSex2, _ := strconv.Atoi(rsex2)
			EnJob2, _ := strconv.Atoi(rjob2)
			EnFrame2, _ := strconv.Atoi(fields2["frame"])
			EnAvatar2, _ := strconv.Atoi(fields2["avatar_v2"])
			EnServerId, _ := strconv.Atoi(fields2["server_id"])
			serverInfo, err := global.ServerListMgrInstance.Get(EnServerId)
			if err != nil {
				serverInfo.Name = ""
			}
			RankDetail.RankOwn = protocol.RankUserInfo{
				Uid:        Uid,
				Num:        resultsValue2 + 1,
				Avatar:     ravatar2,
				Name:       rname2,
				Value:      OwnScore,
				Type:       beginType,
				Sex:        int32(EndSex2),
				Job:        int32(EnJob2),
				ServerName: serverInfo.Name,
				AvatarV2:   int32(EnAvatar2),
				Frame:      int32(EnFrame2),
			}
		}
	}

	return RankDetail, nil
}

func RankMysql() ([]*pb.Rank, error) {
	RankUserInfo := []*pb.Rank{}
	res, _ := db.RankList()
	for _, item := range res {
		zip := pb.RankInfo{
			List: []*pb.Rank{},
		}
		if err := protobuf.NewSerializer().Unmarshal(item.List, &zip); err == nil {
			for _, item2 := range zip.List {
				input := protocol.RankUserInfoInput{
					Uid:     item2.Uid,
					Type:    item.Type,
					Compare: false,
					Name:    item2.Name,
					Avatar:  item2.Avatar,
					Value:   item2.Value,
				}
				_, err1 := RankChangeTCP(&input)
				if err1 != nil {
					return nil, err
				}
			}

			RankUserInfo = zip.List
		} else {
			// logger.Error("反序列化失败:", err)
		}
	}

	return RankUserInfo, nil
}

func CronSendRankAwardSeason() {
	// //查询所有服务器Id
	// ServerIds, err := GetAllServerId()
	// if err != nil {
	// 	return
	// }
	ServerIds := make([]int32, 0)
	for i := int(global.ServerIDRangeStart); i <= int(global.ServerIDRangeEnd); i++ {
		ServerIds = append(ServerIds, int32(i))
	}

	if len(ServerIds) == 0 {
		return
	}
	// 遍历所有服务器Id
	for _, serverId := range ServerIds {
		rank := strconv.Itoa(int(serverId)) + "rankSend:"

		// 将字符串转换为time.Time类型，并定义时间段数组
		var timePeriods []struct {
			start time.Time
			end   time.Time
		}
		timePeriodsStr := global.GlobalCfg.PVPSeason.Season

		for _, periodStr := range timePeriodsStr {
			start, err := time.Parse(time.DateOnly, periodStr.StartTime)
			if err != nil {
				// logger.Error("Error parsing start time:", err)
			}
			end, err := time.Parse(time.DateOnly, periodStr.EndTime)
			if err != nil {
				// logger.Error("Error parsing end time:", err)
			}
			timePeriods = append(timePeriods, struct {
				start time.Time
				end   time.Time
			}{start, end})
		}
		// 获取当前时间
		currentTime := time.Now()
		awardStatus := 0
		// 检查当前时间是否大于第一个时间段的开始时间
		// 遍历时间段，检查当前时间是否在时间段的间隙中
		var lastPeriodIndex int
		triggered := false
		for i, period := range timePeriods {
			if currentTime.Before(period.start) {
				if i == 0 || currentTime.After(timePeriods[i-1].end) {
					// 当前时间在时间段间隙中，且是刚过去的时间段
					lastPeriodIndex = i - 1
					triggered = true
					break
				}
			} else if currentTime.After(period.end) {
				// 更新刚过去的时间段索引
				lastPeriodIndex = i
			}
		}
		if triggered {
			if lastPeriodIndex >= 0 {
				end := lastPeriodIndex + 1
				logger.Infof("当前时间在时间段间隙中，刚过完第 %d 个时间段。", end)
				// 查询是否已经发过
				resStr, err2 := getStringKey(rank + strconv.Itoa(int(end)))
				if err2 != nil {
					// logger.Error("获取redis失败", err2)
					return
				}
				if resStr == "" {
					// 查数据库
					res1, _ := db.RankList()
					if len(res1) == 0 {
						// logger.Error("数据出错。")
						return
					}
					for _, item := range res1 {
						if int(item.Type) == 70000+end {
							awardStatus = item.SendStatus
						}
					}
					if awardStatus == 1 {
						// logger.Error("已经更新过了。")
						return
					}
					err3 := setStringKey(rank+strconv.Itoa(int(end)), "1")
					if err3 != nil {
						// logger.Error("获取redis失败", err3)
						return
					}
					_, err4 := db.RankUpdateSend(int32(70000 + end))
					if err4 != nil {
						// logger.Error("更新到数据库失败")
						return
					}
					// logger.Error("成功。")
					rtype := []int64{7}
					err := SendRankAward(rtype, false)
					if err != nil {
					}
				}
			} else {
				// logger.Debug("当前时间在第一个时间段之前，没有刚过完的时间段。")
			}
		} else {
			// logger.Debug("当前时间位于时间段内，不触发间隙触发条件")
		}
	}

	return
	//results, err2 := CmdCache("set", rank+strconv.Itoa(int(item)), rankBegin, rankEnd, "WITHSCORES")
	//if results == nil {
	//	continue
	//}
	//rtype := []int64{7}
	//err := SendRankAward(rtype)
	//if err != nil {
	//}
}

// 邮件
func Mailadd(r *protocol.MailAddRequest) error {
	if r == nil {
		return errutil.ErrIllegalParameter
	}
	_, result, err := db.MailAdd(
		r.Uid,
		r.RecipientUid,
		r.Title,
		r.Desc,
		r.Attachment,
		r.ChannelCode,
	)
	if err != nil || !result {
		return err
	}

	return nil
}

// 昨日充值入库
func Bounty(uid int64) ([]protocol.PrivilegeDetailTCP, error) {
	rank := strconv.Itoa(int(global.ServerIDRangeStart)) + "Bounty:"
	// 获取昨日总金额
	allMoney, err := db.Storetotal(0)
	if err != nil {
		return nil, err
	}
	BountyCfg := global.GlobalCfg.Bounty.Grade
	// result := make(map[int32][]protocol.SourceInfo)
	result := make([]protocol.PrivilegeDetailTCP, 0)
	// 特权卡
	allWeight := 0
	for _, privilegen := range BountyCfg {
		allWeight = allWeight + int(privilegen.Weight)
	}
	endMoney := allMoney[0].Num * int64(global.GlobalCfg.BaseValue.Bonus.Rate)
	for _, privilegen := range BountyCfg {
		key := rank + strconv.Itoa(int(privilegen.ID))
		num := int(int(endMoney) * (int(privilegen.Weight) / allWeight))
		err := ZaddRank(key, float64(num), int64(privilegen.ID))
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// 昨日充值入库
func ltv(uid int64) ([]RankListDetail, error) {
	result := make([]RankListDetail, 0)
	return result, nil
}

type RankListDetail struct {
	Today  int  `json:"today"`  // 当天
	Two    int  `json:"two"`    // 次日
	Three  int  `json:"three"`  // 3日
	Seven  bool `json:"seven"`  // 7日
	Thirty bool `json:"thirty"` // 7日
}

func RequestWeb(methodPost string, url string, reqBody interface{}) ([]byte, error) {
	reqJson, _ := json.Marshal(&reqBody)
	urlReq, _ := http.NewRequest(methodPost, url, bytes.NewReader(reqJson))
	parma := urlReq.URL.Query()
	urlReq.URL.RawQuery = parma.Encode()
	res, err := HttpDo(urlReq)
	if err != nil {
		// logger.Error("请求失败", err)
		return nil, err
	}
	return res, nil
}

type BaseRequest struct {
	Code   string `form:"code"`
	Ipaddr string `form:"ipaddr"`
	Show   int    `form:"show"`
}

func HttpDo(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 保存用户信息
func SaveUserInfo(r *protocol.RankUserInfoInput) error {
	rankMember := r.Uid
	rname := r.Name
	rsex := strconv.Itoa(int(r.Sex))
	rjob := strconv.Itoa(int(r.Job))
	rrobot_cid := strconv.Itoa(int(r.RobotCId))
	rserver_id := strconv.Itoa(int(r.ServerId))
	// key2 := constant.RedisRankUserKey(r.Uid)
	ravatar := r.Avatar
	ravatar2 := strconv.Itoa(int(r.AvatarV2))
	rframe := strconv.Itoa(int(r.Frame))
	if err2 := hmsetUserFields(r.ServerId, rankMember, buildStandardUserFields(rname, ravatar, rsex, rjob, rrobot_cid, rserver_id, ravatar2, rframe)); err2 != nil {
		return err2
	}
	return nil
}
