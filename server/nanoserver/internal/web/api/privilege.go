package api

import (
	"context"
	"nanoserver/db"
	"nanoserver/global"
	"nanoserver/internal/game"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakePrivilegeService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/privilege/listcontinuity", nex.Handler(privilegeListContinuity)).Methods("POST")                    //列表
	router.Handle("/v1/privilege/receive", limitHandler(nex.Handler(privilegeReceive), "privilegeReceive")).Methods("PUT") //领取奖励

	return router
}

func privilegeListContinuity(ctx context.Context) (*protocol.PrivilegeListResponse, error) {

	//UserRound := int32(1)
	//Round := int32(1)
	//Day := int32(0)
	//UserRound := int32(1)
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.PrivilegeOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	list, err := PrivilegeList(request)
	if err != nil {
		return &protocol.PrivilegeListResponse{}, err
	}
	return &protocol.PrivilegeListResponse{Data: list}, nil
}
func PrivilegeList(r *protocol.PrivilegeOperateRequest) ([]protocol.PrivilegeDetailDay, error) {

	privilegeCfg := global.PrivilegeCfg
	//result := make(map[int32][]protocol.SourceInfo)
	result := make([]protocol.PrivilegeDetailDay, 0)
	//查询用户今日领取情况
	resultd, _ := db.PrivilegeListDay(
		r.Uid,
	)
	PidWhile := make(map[int32]bool)
	for _, item := range resultd {
		PidWhile[int32(item.Type)] = true
	}

	//特权卡
	for _, privilegen := range privilegeCfg {
		buys := 1
		gets := 1
		endtime := int64(-1)
		goodsId := make([]int64, 0)
		request2 := &protocol.StoreListRequest{}
		request2.Uid = r.Uid
		request2.StoreType = privilegen.Sid / 1000
		request2.GoodsUuid = privilegen.Sid
		//request.StoreType = r.StoreType
		//request.GoodsId = r.GoodsId
		//查询用户购买情况
		day := 30
		switch privilegen.Pid {
		case 1:
			day = 30
		case 2:
			day = 3650
		}
		goodsId = append(goodsId, int64(privilegen.Sid))
		resultb, _ := db.StoreList(
			r.Uid,
			request2.StoreType,
			int32(day),
			goodsId,
			true,
		)
		if len(resultb) > 0 {
			originalTime := time.Unix(resultb[0].CreateTime, 0)
			nextDayMidnight := originalTime.AddDate(0, 0, int(privilegen.Duration)) // 添加N天
			nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
			// 转换回时间戳
			nextDayMidnightTimestamp := nextDayMidnight.Unix()
			buys = 2
			endtime = nextDayMidnightTimestamp
		}

		if PidWhile[privilegen.Pid] {
			gets = 2
		}
		result = append(result, protocol.PrivilegeDetailDay{
			Pid:     privilegen.Pid,
			Buy:     buys,
			Get:     gets,
			EndTime: endtime,
		})
	}
	return result, nil
}

func privilegeReceive(ctx context.Context, r *protocol.PrivilegeOperateRequest) (*protocol.PrivilegeReceiveDayResponse, error) {

	//UserRound := int32(1)
	//Round := int32(1)
	//Day := int32(0)
	//UserRound := int32(1)
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request := &protocol.PrivilegeOperateRequest{}
	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.PID = r.PID
	if request.PID == 0 {
		return &protocol.PrivilegeReceiveDayResponse{}, errutil.ErrInvalidParameter
	}

	Cumulative, err := PrivilegeCumulative(request)
	if !Cumulative {
		return &protocol.PrivilegeReceiveDayResponse{Code: -1, Data: err}, nil
	}
	return &protocol.PrivilegeReceiveDayResponse{Data: "成功"}, nil
}

// 签到，给到日期
func PrivilegeCumulative(r *protocol.PrivilegeOperateRequest) (bool, string) {
	privilegeCfg := global.PrivilegeCfg
	//result := make(map[int32][]protocol.SourceInfo)
	//查询用户今日领取情况
	resultd, _ := db.PrivilegeListDay(
		r.Uid,
	)

	PidWhile := make(map[int32]bool)
	for _, item := range resultd {
		PidWhile[int32(item.Type)] = true
	}
	//特权卡
	lista := make([]global.Award, 0)
	for _, privilegen := range privilegeCfg {
		if privilegen.Pid == r.PID {
			request2 := &protocol.StoreListRequest{}
			request2.Uid = r.Uid
			request2.StoreType = privilegen.Sid / 1000
			request2.GoodsUuid = privilegen.Sid
			//request.StoreType = r.StoreType
			//request.GoodsId = r.GoodsId
			lists, err2 := StoreList(request2)
			if err2 != nil || lists == nil {
				return false, "获取出错"
			}
			if lists[0].Buy == 0 {
				return false, "没有购买"
			}
			if PidWhile[privilegen.Pid] {
				return false, "已经领取过"
			}
			for _, awards := range privilegen.Dailyaward {
				lista = append(lista, awards)
			}
		}
	}

	//插入
	_, err := db.PrivilegeAddDay(
		r.Uid,
		r.PID,
	)
	if err != nil {
		return false, err.Error()
	}
	//校验背包情况
	result2 := make([]protocol.SourceInfo, len(lista))
	for i, info := range lista {
		result2[i].Id = info.Id
		result2[i].Num = int64(info.Num)
	}
	err3 := CheckSourceAdd(result2, r.Uid, true, constant.ItemPrivilege)
	if err3 != nil {
		if err3 == errutil.ErrPlayerNotFound {
			return false, err3.Error()
		}
		if err3 == errutil.ErrBagFull {

			MailId := strconv.FormatInt(snowflake.GenID(), 10)
			result3 := make([]protocol.Attachment, len(result2))
			for i, Attachment := range result2 {
				result3[i] = protocol.Attachment{
					MailId:      MailId,
					UserPropsId: int64(Attachment.Id),
					Num:         Attachment.Num,
				}
			}

			request3 := &protocol.MailAddRequest{
				Uid:          r.Uid,
				RecipientUid: -1,
				Title:        "系统邮件",
				Desc:         "包裹已满，请及时清理",
				Attachment:   result3,
				ChannelCode:  "0",
			}
			resultc, err := Mailadd(request3)
			if !resultc {
				logger.Errorf("MailAdd error: %v", err)
				return false, err.Error()
			}
		}

	}
	vip_card_type := ""
	switch r.PID {
	case 1:
		vip_card_type = "日卡"
	case 2:
		vip_card_type = "终身卡"
	}

	upPrivilegeShushuDate(r.Uid, vip_card_type)

	return true, ""
}
func upPrivilegeShushuDate(Uid int64, vip_card_type string) {
	shushuDate := make(map[string]interface{}, 0)
	//在合并点其他的数据
	shushuDate["vip_card_type"] = vip_card_type
	shushuDate["is_success"] = true
	shushuDate["name"] = "vip_card"
	game.Upshushu(Uid, shushuDate, true, false)
}
