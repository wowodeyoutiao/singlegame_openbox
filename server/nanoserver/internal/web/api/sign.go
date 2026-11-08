package api

import (
	"context"
	"encoding/json"
	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/internal/game"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"

	//"nanoserver/internal/web"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol"
	"net/http"
	"time"
)

func MakeSignService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/sign/add", nex.Handler(signAdd)).Methods("Put")                                        //添加
	router.Handle("/v1/sign/list", nex.Handler(signList)).Methods("POST")                                     //列表
	router.Handle("/v1/sign/extract", limitHandler(nex.Handler(signExtract), "sign_extract")).Methods("POST") //提取累计

	return router
}

func SignList(r *protocol.SignOperateRequest) (protocol.Awsard, error) {

	//UserRound := int32(1)
	Round := int32(1)
	Day := int32(0)
	UserRound := int32(1)
	//查询用户目前日签到进度
	resultd, _ := db.SignListDay(
		r.Uid,
	)
	if len(resultd) == 0 {
		return protocol.Awsard{
			DayAward: protocol.DayAward{
				Day:   1,
				Round: 1,
				Sign:  false,
			},
			AccumulateAward: make([]protocol.AccumulateAward, 0),
		}, nil
	}
	t1 := time.Unix(resultd[0].CreatedAt, 0)
	t2 := time.Now()
	check := false
	UserRound = resultd[0].UserRound
	Round = resultd[0].Round
	Day = resultd[0].Day
	if t2.Format("2006-01-02") == t1.Format("2006-01-02") {
		check = true

	} else {
		day1 := int32(0)
		roundWhile := make(map[int32]bool)
		result := make(map[int32][]protocol.SourceInfo)
		signCfg := global.SignCfg
		for _, signn := range signCfg["day"] {
			if signn.Round == Round {
				result[signn.Days] = make([]protocol.SourceInfo, len(signn.Content))
				day1++
			}
			if !roundWhile[signn.Round] {
				roundWhile[signn.Round] = true
			}
		}

		if Round >= int32(len(roundWhile)) {
			if Day >= day1 {
				UserRound++
				Round = 1
				Day = 1
			} else {
				Day++
			}
		} else {
			if Day >= day1 {
				Round++
				Day = 1
			} else {
				Day++
			}
		}
		//UserRound = resultd[0].UserRound
		//Round = resultd[0].Round
		//Day = resultd[0].Day
	}
	resulta := []model.Sign{}
	if len(resultd) > 0 {

		// 用户最新  累计签到情况
		resulta, _ = db.SignListAccumulate(
			r.Uid,
			UserRound,
			Round,
		)
	}

	result := protocol.Awsard{
		DayAward: protocol.DayAward{
			Day:   Day,
			Round: Round,
			Sign:  check,
		},
		AccumulateAward: make([]protocol.AccumulateAward, len(resulta)),
	}
	for i, sign := range resulta {
		result.AccumulateAward[i] = protocol.AccumulateAward{
			Day:   sign.Day,
			Round: sign.Round,
		}
	}
	return result, nil
}

func signList(ctx context.Context, r *http.Request) (*protocol.SignListResponse, error) {

	request := &protocol.SignOperateRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	list, err := SignList(request)
	if err != nil {
		return nil, err
	}
	return &protocol.SignListResponse{Data: list}, nil
}

func signAdd(ctx context.Context, r *http.Request) (*protocol.SignOperateResponse, error) {
	request := &protocol.SignOperateRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	AccountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName

	request.Uid = Uid
	request.AccountName = AccountName

	_, err = SignnAdd(request)
	if err != nil {
		if err == errutil.ErrGetAward {
			return &protocol.SignOperateResponse{Data: false, Code: -1}, nil
		}
		return nil, err
	}
	return &protocol.SignOperateResponse{Data: true, Code: 0}, err
}

// 签到 并领奖
func SignnAdd(r *protocol.SignOperateRequest) (bool, error) {
	//查询目前签到情况
	resultd, err := db.SignListDay(
		r.Uid,
	)
	if err != nil {
		return false, err
	}
	UserRound := int32(1)
	Round := int32(1)
	//默认为1
	Day := int32(1)
	day := int32(0)

	//查询vip等级
	// 个人信息
	userInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return false, err
	}
	vipLevel := userInfo.Vip

	roundWhile := make(map[int32]bool)
	result := make(map[int32][]protocol.SourceInfo)
	signCfg := global.SignCfg
	for _, signn := range signCfg["day"] {
		if signn.Round == Round {
			result[signn.Days] = make([]protocol.SourceInfo, len(signn.Content))
			if len(signn.Content) > 0 {
				bei := int32(1)
				viplimit, _ := strconv.ParseInt(signn.VIPLimit, 10, 32)
				if vipLevel >= int32(viplimit) {
					bei = 2
				}
				for j, prop := range signn.Content {
					result[signn.Days][j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  int32(prop.Id),
						Num: int64(prop.Num * bei),
					}
				}
			}
			day++
		}
		if !roundWhile[signn.Round] {
			roundWhile[signn.Round] = true
		}
	}
	if len(resultd) > 0 {
		t1 := time.Unix(resultd[0].CreatedAt, 0)
		t2 := time.Now()
		if t2.Format("2006-01-02") == t1.Format("2006-01-02") {
			return false, errutil.ErrGetAward
		}
		//如果这轮结束了，那么就是新的一轮，如果整个轮数都结束了，那么重新循环
		UserRound = resultd[0].UserRound
		Round = resultd[0].Round
		Day = resultd[0].Day

		if Round >= int32(len(roundWhile)) {
			if Day >= day {
				UserRound++
				Round = 1
				Day = 1
			} else {
				Day++
			}
		} else {
			if Day >= day {
				Round++
				Day = 1
			} else {
				Day++
			}
		}

	}

	//发放签到的奖励
	err2 := CheckSourceAdd(result[Day], r.Uid, true, constant.ItemSign)
	if err2 != nil {
		return false, err2
	}
	jsonDate, err := json.Marshal(result[Day])
	if err != nil {
		return false, err
	}

	jsonDateString := string(jsonDate)
	//成功了直接记录
	resulta, err := db.SignAddDay(
		r.Uid,
		vipLevel,
		UserRound,
		Round,
		Day,
		jsonDateString,
	)
	upSignShushuData(Day, false, result[Day], r.Uid, 1)
	if err != nil || !resulta {
		return false, err
	}
	return true, nil
}
func upSignShushuData(check_in_day int32, is_supplement bool, sourceInfo []protocol.SourceInfo, uid int64, ptype int32) {
	data := make(map[string]interface{}, 0)
	data["name"] = "check_in"
	data["check_in_type"] = check_in_day
	data["check_in_day"] = check_in_day
	data["is_supplement"] = is_supplement
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		if v.Num > 0 {
			//查询itemname
			if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   ItemCfg.ID,
					"item_name": ItemCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	data["get_item_info"] = get_item_info
	game.Upshushu(uid, data, true, false)
}
func signExtract(ctx context.Context, r *http.Request, AccumulateAward *protocol.AccumulateAwardRequest) (*protocol.SignOperateResponse, error) {
	request := &protocol.AccumulateAwardRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	AccountName := ctx.Value("userDate").(my_jwt.CustomClaims).AccountName

	if !(AccumulateAward.Day > int32(0)) || !(AccumulateAward.Round > int32(0)) {
		return nil, errutil.ErrIllegalParameter
	}
	request.Uid = Uid
	request.AccountName = AccountName
	request.Day = AccumulateAward.Day
	request.Round = AccumulateAward.Round

	_, err = SignnExtract(request)
	if err != nil {
		return nil, err
	}
	return &protocol.SignOperateResponse{Data: true, Code: 0}, err
}
func SignnExtract(r *protocol.AccumulateAwardRequest) (bool, error) {
	//查询目前签到情况
	resultd, err := db.SignListDay(
		r.Uid,
	)
	if err != nil {
		return false, err
	}
	UserRound := int32(1)
	Round := int32(1)
	//默认为1
	Day := int32(1)
	resulta := []model.Sign{}
	if len(resultd) == 0 {
		return false, errutil.ErrChooseRound
	}
	//如果这轮结束了，那么就是新的一轮，如果整个轮数都结束了，那么重新循环
	UserRound = resultd[0].UserRound
	Round = resultd[0].Round
	Day = resultd[0].Day
	//查询目前累计签到情况
	resulta, err = db.SignListAccumulate(
		r.Uid,
		UserRound,
		Round,
	)
	//只能领取这个轮次的
	if r.Round != Round {
		return false, errutil.ErrChooseRound
	}
	if r.Day > Day {
		return false, errutil.ErrChooseRound
	}
	//映射
	roundWhile := make(map[int32]bool)
	if len(resulta) > 0 {
		for _, signn := range resulta {
			if !roundWhile[signn.Day] {
				roundWhile[signn.Day] = true
			}
		}
	}
	//查询vip等级
	check, err := db.UserCheck(r.Uid, r.AccountName)
	if err != nil {
		return false, err
	}
	if len(check) == 0 {
		return false, errutil.ErrChooseRole
	}
	vipLevel := int32(check[0].Vip)

	signCfg := global.SignCfg
	globalWhile := make(map[int32]bool)
	result := make(map[int32][]protocol.SourceInfo)
	for _, signn := range signCfg["accumulate"] {
		if signn.Round == Round && signn.Days == r.Day {
			if len(signn.Content) > 0 {
				bei := int32(1)
				viplimit, _ := strconv.ParseInt(signn.VIPLimit, 10, 32)
				if vipLevel >= int32(viplimit) {
					bei = 2
				}
				result[signn.Days] = make([]protocol.SourceInfo, len(signn.Content))
				for j, prop := range signn.Content {
					result[signn.Days][j] = protocol.SourceInfo{
						// ... 根据字段一一复制值
						Id:  int32(prop.Id),
						Num: int64(prop.Num * bei),
					}
				}
			}
			globalWhile[signn.Round] = true
		}
	}
	if roundWhile[r.Day] {
		return false, errutil.ErrGetAward
	}
	//超过的天数，不能领取

	//todo发放签到的奖励
	if len(result[r.Day]) == 0 {
		return false, errutil.ErrNotFindAward
	}
	////发放签到的奖励
	err2 := CheckSourceAdd(result[r.Day], r.Uid, true, constant.ItemSign)
	if err2 != nil {
		return false, err2
	}
	jsonDate, err := json.Marshal(result[r.Day])
	if err != nil {
		return false, err
	}

	jsonDateString := string(jsonDate)
	//插入
	resultEnd, err := db.SignExtract(
		r.Uid,
		vipLevel,
		UserRound,
		Round,
		r.Day,
		jsonDateString,
	)

	if err != nil || !resultEnd {
		return false, err
	}
	return true, nil
}
