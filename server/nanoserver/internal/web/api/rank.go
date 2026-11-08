package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"nanoserver/db"
	"nanoserver/global"
	"nanoserver/internal/game"
	"nanoserver/internal/web/provide"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/lonng/nano/serialize/protobuf"
	"github.com/lonng/nex"
)

func MakeRankService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/rank/list", nex.Handler(rankList)).Methods("POST")                 // 列表
	router.Handle("/v1/rank/add", nex.Handler(RankChangeTCP)).Methods("PUT")              // 列表
	router.Handle("/v1/rank/addrange", nex.Handler(RankRange)).Methods("PUT")             // 列表
	router.Handle("/v1/rank/send", nex.Handler(SendRankAward)).Methods("PUT")             // 列表
	router.Handle("/v1/rank/del", nex.Handler(DelRank)).Methods("Delete")                 // 列表
	router.Handle("/v1/rank/listmysql", nex.Handler(RankMysql)).Methods("POST")           // 列表
	router.Handle("/v1/rank/sendswason", nex.Handler(SendRankAwardSeason)).Methods("PUT") // 列表
	return router
}

func rankList(ctx context.Context, r *protocol.RankListRequest) (*protocol.RankListResponse, error) {
	// UserRound := int32(1)
	// Round := int32(1)
	// Day := int32(0)
	// UserRound := int32(1)
	// 查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	if r.Type == 0 {
		return &protocol.RankListResponse{Code: -1}, errutil.ErrInvalidParameter
	}
	request := &protocol.RankOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]

	request.Type = r.Type
	request.All = r.All
	if r.Type == 0 {
		return nil, errutil.ErrInvalidParameter
	}
	// 有些榜单需要传入服务器ID
	ress, err := game.UserInfo(Uid)
	if err != nil {
		return nil, err
	}
	if ress != nil {
		request.Name = ress.Name
	}
	if RankTypeServer, ok := constant.RankTypeServer[r.Type]; ok && RankTypeServer {
		request.ServerId = ress.ServerId
	} else {
		request.ServerId = ress.ZoneId
	}
	// 如果是
	if r.Type == constant.RankTypeWorldBossKf {
		request.Uid = int64(ress.ServerId) * 100
	} else {
		request.Uid = Uid
	}
	list, err := provide.RankList(request)
	if err != nil {
		return &protocol.RankListResponse{}, err
	}
	return &protocol.RankListResponse{Data: list}, nil
}

func RankChangeTCP(r *protocol.RankUserInfoInput) (bool, error) {
	rankEnd := "99"

	// for i := 0; i < 200; i++ {
	rank := strconv.Itoa(int(r.ServerId)) + "rank:"
	user := strconv.Itoa(int(r.ServerId)) + "user:"
	name := r.Name
	avatar := r.Avatar
	sex := strconv.Itoa(int(r.Sex))
	job := strconv.Itoa(int(r.Job))
	// 需要比较就查询比大小
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
			r.Type = r.Type*10000 + int32(v.ID)
			break
		}
	}
	key := rank + strconv.Itoa(int(r.Type))
	key2 := user + strconv.Itoa(int(r.Uid))
	_, err2 := provide.CmdCache("hset", key2, "name", name, "avatar", avatar, "sex", sex, "job", job)
	if err2 != nil {
		return false, err2
	}
	if r.Compare {
		// 用户的分数
		resultScores, err5 := provide.CmdCache("zscore", rank+strconv.Itoa(int(r.Type)), strconv.FormatInt(r.Uid, 10))

		// 之前没分数不用比较了
		if resultScores == nil {
			err := provide.ZAddCache(key, r.Value, r.Uid)
			if err != nil {
				return false, err
			}
		} else {
			if err5 != nil {
				return false, err5
			}
			resultScores2, ok2 := resultScores.(interface{})
			if !ok2 {
				logger.Info("resultScores is not a slice of interfaces")
				return false, nil
			}
			OwnScore, err7 := redis.Int64(resultScores2, nil)
			if err7 != nil {
				logger.Error("Error converting member to string:", err7)
			}
			if OwnScore < r.Value {
				err := provide.ZAddCache(key, r.Value, r.Uid)
				if err != nil {
					return false, err
				}
			}
		}

	} else {
		err := provide.ZAddCache(key, r.Value, r.Uid)
		if err != nil {
			return false, err
		}
	}
	if !r.Mysql {
		return true, nil
	}
	// 落地到数据库
	// 如果大于第100个人，那么就更新到数据库
	resultScores, _ := provide.CmdCache("zrevrange", rank+strconv.Itoa(int(r.Type)), rankEnd, rankEnd, "WITHSCORES")
	resultsValue, ok := resultScores.([]interface{})
	if !ok {
		logger.Info("resultsValue is not a slice of interfaces")
	}
	score := int64(0)
	// 遍历结果，每两个元素是一个成员和它的分数
	for i := 0; i < len(resultsValue); i += 2 {
		UserScore, err := redis.Int64(resultsValue[i+1], nil)
		if err != nil {
			logger.Error("Error converting score to float64:", err)
			continue
		}
		score = UserScore
	}
	if r.Value > score {
		rankInfo := &protocol.RankOperateRequest{
			Type: 1,
			All:  2,
		}
		res, _ := provide.RankList(rankInfo)
		zip := pb.RankInfo{
			List: []*pb.Rank{},
		}
		for _, v := range res.RankAll {
			zip.List = append(
				zip.List,
				&pb.Rank{
					Uid:    v.Uid,
					Num:    v.Num,
					Avatar: v.Avatar,
					Name:   v.Name,
					Value:  v.Value,
					Sex:    v.Sex,
				})
		}
		jsonBytes, err := protobuf.NewSerializer().Marshal(&zip)
		if err != nil {
			logger.Info("err:", err)
		}
		_, err1 := db.RankUpdate(r.Type, jsonBytes)
		if err1 != nil {
			logger.Info("err1:", err1)
		}
	}

	return true, nil
}

func SendRankAward() (bool, error) {
	rType := []int64{1, 2, 3, 4}
	rTypeName := []string{"", "等级", "战力", "乱斗", "试炼"}
	// rankAwardCfg := global.RankingAwardCfg
	// result := make(map[int32][]protocol.SourceInfo)
	// 活动内容
	// activity := make(map[int32][]global.KeyValue)
	// for _, activityn := range rankAwardCfg {
	// 	activity[activityn.UUID] = activityn.Award
	// }
	for _, item := range rType {

		rank := strconv.Itoa(int(global.ServerIDRangeStart)) + "rank:"
		rankBegin := "0"
		rankEnd := "99"
		RankAll := make([]protocol.RankUserInfo, 0)

		results, err2 := provide.CmdCache("zrevrange", rank+strconv.Itoa(int(item)), rankBegin, rankEnd, "WITHSCORES")
		if err2 != nil {
			logger.Error("自动脚本报错", err2)
			return false, nil
		}
		if results == nil {
			logger.Info("没有数据", err2)
			return false, nil
		}
		// 初始化一个映射来存储每个人及其对应的分数
		// 检查结果是否为 []interface{}
		resultsValue, ok := results.([]interface{})
		if !ok {
			logger.Info("resultsValue is not a slice of interfaces")
		}
		// 遍历结果，每两个元素是一个成员和它的分数
		for i := 0; i < len(resultsValue); i += 2 {
			member, err := redis.Int64(resultsValue[i], nil)
			if err != nil {
				logger.Errorf("Error converting member to string: %v", err)
				continue
			}
			score, err := redis.Int64(resultsValue[i+1], nil)
			if err != nil {
				logger.Errorf("Error converting member to string: %v", err)
				continue
			}

			//if err != nil {
			//	continue
			//}
			RankAll = append(RankAll, protocol.RankUserInfo{
				Uid:   member,
				Num:   int64(i/2 + 1),
				Type:  item,
				Value: score,
			})
		}
		for _, item2 := range RankAll {
			result2 := global.GetRankAwardCfg(int32(item)*10000 + int32(item2.Num))
			MailId := strconv.FormatInt(snowflake.GenID(), 10)
			result3 := make([]protocol.Attachment, len(result2.Award))
			for i, Attachment := range result2.Award {
				result3[i] = protocol.Attachment{
					MailId:      MailId,
					UserPropsId: int64(Attachment.Key),
					Num:         int64(Attachment.Value),
				}
			}
			request3 := &protocol.MailAddRequest{
				Uid:          item2.Uid,
				RecipientUid: -1,
				Title:        "系统邮件：排行榜奖励",
				Desc:         rTypeName[item] + "排行榜奖励",
				Attachment:   result3,
				ChannelCode:  "0",
			}
			resultc, err := Mailadd(request3)
			if !resultc {
				return true, err
			}
		}
		time.Sleep(time.Millisecond)
	}
	return true, nil
}

func DelRank() (bool, error) {
	rType := []int64{4}
	for _, item := range rType {
		rank := strconv.Itoa(int(global.ServerIDRangeStart)) + "rank:"
		_, err := provide.CmdCache("del", rank+strconv.Itoa(int(item)))
		if err != nil {
			logger.Errorf("DelRank error: %v", err)
		}

	}
	return true, nil
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
					Sex:     item2.Sex,
					Job:     item2.Job,
				}
				_, err1 := RankChangeTCP(&input)
				if err1 != nil {
					return nil, err
				}
			}
			// logger.Errorf("恢复数据成功 : %d", err)
			RankUserInfo = zip.List
		} else {
			logger.Errorf("数据反序列化失败 : %d", err)
		}
	}

	return RankUserInfo, nil
}

func RankRange() (bool, error) {
	for i := 1; i < 8; i++ {
		for j := 1; j < 1; j++ {
			input := protocol.RankUserInfoInput{
				Uid:     int64(j),
				Type:    int32(i),
				Compare: false,
				Name:    strconv.Itoa(j) + strconv.Itoa(j),
				Value:   int64(j)*10000 + int64(j)*1000 + int64(j)*100 + int64(j),
				Mysql:   true,
			}
			_, err := RankChangeTCP(&input)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

// 1234每日定时发奖励
func SendRankAwardSeason() (bool, error) {
	provide.SendRankAward([]int64{1, 2, 3, 4}, false)
	return true, nil
}

// func RankLamp(name string, rankName string) {
// 	strSend := []string{name, rankName}
// 	game.Lamp(int32(constant.LampPetSummon), strSend)
// }
