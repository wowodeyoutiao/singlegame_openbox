package api

import (
	"context"
	"encoding/json"
	"math"
	"nanoserver/db"
	"nanoserver/global"
	"nanoserver/internal/game"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/snowflake"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakePrizeDrawService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/prize_draw/detail", nex.Handler(prizeDrawDetail)).Methods("POST")                                                  //抽奖情况
	router.Handle("/v1/prize_draw/partner_receive", limitHandler(nex.Handler(prizeDrawPartnerReceive), "partner_receive")).Methods("PUT") //伙伴抽奖
	router.Handle("/v1/prize_draw/pet_receive", limitHandler(nex.Handler(prizeDrawPetReceive), "pet_receive")).Methods("PUT")             //宠物抽奖
	return router
}

func prizeDrawDetail(ctx context.Context, r *protocol.PrizedrawRequest) (*protocol.PrizedrawListResponse, error) {
	//查询用户登录情况
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	if r.Module == 0 {
		return &protocol.PrizedrawListResponse{Code: -1}, nil
	}
	request := &protocol.PrizedrawRequest{}

	//[{2024-05-09 2} {2024-05-10 1}]
	request.Uid = Uid
	request.Module = r.Module
	CanPlay := make([]protocol.PrizedrawDetail, 0)
	for i := 1; i < int(constant.LotteryMax-1); i++ {
		request.PType = int32(i)
		residue, residueBig, err := PrizeDrawDetail(request)
		if err != nil {
			return &protocol.PrizedrawListResponse{Code: -1, Data: err.Error()}, nil
		}
		CanPlay = append(CanPlay, protocol.PrizedrawDetail{
			Type:    int32(i),
			CanPlay: residue,
			NeedNum: residueBig,
		})
	}
	// if err != nil {
	// 	return &protocol.PrizedrawListResponse{Code: -1, Data: err.Error()}, nil
	// }
	return &protocol.PrizedrawListResponse{CanPlay: CanPlay}, nil
}

// 伙伴抽奖抽奖
func PrizeDrawDetail(r *protocol.PrizedrawRequest) (int32, int32, error) {
	//个人等级
	uerInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return 0, 0, err
	}
	haveNum, err := db.PrizedrawList(r.Uid, r.PType, r.Module, 0)
	haveNum2, err := db.PrizedrawList(r.Uid, r.PType, r.Module, 1)
	if err != nil {
		return 0, 0, err
	}
	PartnerCfg := global.GetSummon(r.Module, r.PType, uerInfo.Level)
	if PartnerCfg == nil {
		return 0, 0, errutil.ErrConfigNotFound
	}
	residue := float64(0)
	if PartnerCfg.MaxTime == 0 {
		residue = -1
	} else {
		residue = float64(PartnerCfg.MaxTime) - haveNum2
	}

	residueBig := float64(0)
	faBegin := float64(haveNum)
	nodeNum := 0
	endNum := 0
	for _, item := range PartnerCfg.Node {
		nodeNum = nodeNum + item
		faBegin = faBegin - float64(item)
		endNum = item
		if faBegin >= 0 {
			continue
		} else {
			residueBig = math.Abs(faBegin)
			break
		}
	}
	if int(haveNum) > nodeNum {
		if endNum == 0 {
			residueBig = 0
		} else {
			residueBig = float64(endNum) - float64((int(haveNum)-nodeNum)%endNum)
		}
	}
	return int32(residue), int32(residueBig), nil
}

func prizeDrawPartnerReceive(ctx context.Context, r *protocol.PrizedrawRequest) (*protocol.PrizedrawResponse, error) {
	request := &protocol.PrizedrawRequest{}
	//[{2024-05-09 2} {202405-10 1}]
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request.Uid = Uid

	if r.Type != 1 && r.Type != 2 {
		return &protocol.PrizedrawResponse{Code: -1}, nil
	}
	if r.PType != constant.LotteryTypeBronze && r.PType != constant.LotteryTypeGold && r.PType != constant.LotteryTypeSilver {
		return &protocol.PrizedrawResponse{Code: -1}, nil
	}
	request.Type = r.Type
	request.PType = r.PType
	_, award, err := PrizeDrawPartnerReceive(request)
	if err != nil {
		if err == errutil.ErrItemNotEnough {
			return &protocol.PrizedrawResponse{Code: -2, Data: err.Error()}, nil
		}
		return &protocol.PrizedrawResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.PrizedrawResponse{Award: award}, nil
}

// 伙伴抽奖抽奖
func PrizeDrawPartnerReceive(r *protocol.PrizedrawRequest) ([]protocol.SourceInfo, []protocol.PartnerAward, error) {
	//个人等级
	uerInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, nil, err
	}
	//活动内容
	res := global.GetSummon(constant.LotteryModulePartner, r.PType, uerInfo.Level)
	if res == nil {
		return nil, nil, errutil.ErrConfigNotFound
	}

	// //需要领取的日期
	// needContent := ress % 100
	needNum := -1
	switch r.Type {
	case 1:
		needNum = 1
	case 2:
		needNum = 10
	}
	if needNum == -1 {
		return nil, nil, errutil.ErrConfigNotFound
	}
	//减去消耗品
	resultc := make([]protocol.SourceInfo, 1)
	for _, info := range res.Consume {
		if r.Type == info.Type {
			resultc[0].Id = info.Id
			resultc[0].Num = -int64(info.Num)
		}

	}
	//校验背包情况
	err = CheckSourceNum(resultc, r.Uid, false, constant.ItemPrizeDrawPartner)
	if err != nil {
		return nil, nil, err
	}

	haveNum, err := db.PrizedrawList(r.Uid, r.PType, constant.LotteryModulePartner, 0)
	dayhave, err := db.PrizedrawList(r.Uid, r.PType, constant.LotteryModulePartner, 1)
	if err != nil {
		return nil, nil, err
	}

	result2 := make([]protocol.SourceInfo, 0)
	result3 := protocol.AddSource{
		Uid:    r.Uid,
		Source: make([]*protocol.SourceInfo, 0),
	}
	Partner := []protocol.PartnerAward{}
	if len(uerInfo.Guidance) > 0 && uerInfo.Guidance[0] == constant.GuidanceTypePartner && r.PType == 2 {
		if needNum-int(global.GlobalCfg.Firstgacha.Num) > 0 {
			Partner, err = partnerDetail(r.PType, uerInfo, 9, int(haveNum), int(dayhave))
			if err != nil {
				return nil, nil, err
			}
			for i, info := range Partner {
				if info.Kind == constant.LotteryAwardFragment {
					result2 = append(result2, protocol.SourceInfo{
						Id:  Partner[i].Id,
						Num: int64(Partner[i].Num),
					})
				}
				if info.Kind == constant.LotteryAwardWhole {
					result3.Source = append(result3.Source, &protocol.SourceInfo{
						Id:  Partner[i].Id,
						Num: int64(Partner[i].Num),
					})
				}
			}
		}
		Partner = append(Partner, protocol.PartnerAward{
			Kind: constant.LotteryAwardWhole,
			Id:   global.GlobalCfg.Firstgacha.ID,
			Num:  int32(global.GlobalCfg.Firstgacha.Num),
		})
		result3.Source = append(result3.Source, &protocol.SourceInfo{
			Id:  global.GlobalCfg.Firstgacha.ID,
			Num: int64(global.GlobalCfg.Firstgacha.Num),
		})
		//需要更新服务器状态
		game.SetGuidance(r.Uid, constant.GuidanceTypePartner+1)
	} else {
		Partner, err = partnerDetail(r.PType, uerInfo, needNum, int(haveNum), int(dayhave))
		if err != nil {
			return nil, nil, err
		}
		for i, info := range Partner {
			if info.Kind == constant.LotteryAwardFragment {
				result2 = append(result2, protocol.SourceInfo{
					Id:  Partner[i].Id,
					Num: int64(Partner[i].Num),
				})
			}
			if info.Kind == constant.LotteryAwardWhole {
				result3.Source = append(result3.Source, &protocol.SourceInfo{
					Id:  Partner[i].Id,
					Num: int64(Partner[i].Num),
				})
			}
		}
	}
	//给他加碎片
	err = CheckSourceAdd(result2, r.Uid, true, constant.ItemPrizeDrawPartner)
	if err != nil {
		if err == errutil.ErrPlayerNotFound {
			return nil, nil, err
		}
		if err == errutil.ErrBagFull {
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
			if !resultc || err != nil {
				return nil, nil, err
			}
		}
	}
	//给他加整卡
	err = addPartner(result3)
	if err != nil {
		return nil, nil, err
	}
	for _, info := range result3.Source {
		//查询这个伙伴的品质
		partnerInfo := global.GetPartnerCfg(info.Id)
		if partnerInfo.Quality >= 6 {
			strSend := []string{uerInfo.Name, strconv.Itoa(int(partnerInfo.ItemID))}
			game.Lamp(int32(constant.LampHeroSummon), uerInfo.ServerId, strSend)
		}
	}
	jsonDate, err := json.Marshal(Partner)
	if err != nil {
		return nil, nil, err
	}

	jsonDateString := string(jsonDate)
	err = db.PrizedrawAdd(r.Uid, jsonDateString, int32(needNum), r.PType, constant.LotteryModulePartner)
	if err != nil {
		return nil, nil, err
	}
	err = CheckSourceNum(resultc, r.Uid, true, constant.ItemPrizeDrawPartner)
	if err != nil {
		return nil, nil, err
	}
	// return int32(resulte[0].Num), ""
	sourceInfo := make([]protocol.SourceInfo, 0)
	for _, info := range result2 {
		sourceInfo = append(sourceInfo, protocol.SourceInfo{
			Id:  info.Id,
			Num: info.Num,
		})
	}
	partnerInfo := make([]protocol.SourceInfo, 0)
	for _, info := range result3.Source {
		partnerInfo = append(partnerInfo, protocol.SourceInfo{
			Id:  info.Id,
			Num: info.Num,
		})
	}
	game.TaskDispatch(uerInfo.UID, pb.TargetType_PARTNER_RECEIVE, 0, int32(needNum))
	upPartnerShushu(sourceInfo, partnerInfo, resultc, r.Uid, r.PType)
	return result2, Partner, nil
}

// 伙伴抽奖抽奖
func partnerDetail(ptype int32, userInfo *protocol.UserInfo, needNum int, haveNum int, dayHave int) ([]protocol.PartnerAward, error) {
	award := make([]protocol.PartnerAward, 0)

	PartnerCfg := global.GetSummon(constant.LotteryModulePartner, ptype, userInfo.Level)
	if needNum+dayHave > int(PartnerCfg.MaxTime) {
		return award, errutil.ErrNotCanBuy
	}
	CommonPriority := 0
	SpecialPriority := 0

	//普通奖池
	for _, Partnercitem := range PartnerCfg.CommonPool {
		CommonPriority = CommonPriority + int(Partnercitem.Priority)
	}
	//必中奖池
	for _, Partnersitem := range PartnerCfg.SpecialPool {
		SpecialPriority = SpecialPriority + int(Partnersitem.Priority)
	}
	fa := make(map[int]bool)
	faEnd := 0
	faAdd := 0
	faBegin := 0
	for _, item := range PartnerCfg.Node {
		faBegin = faBegin + item
		fa[faBegin] = true
		faEnd = faBegin
		faAdd = item
	}
	for i := 1; i <= needNum; i++ {
		//如果不是必中，先看看是不是因为其他的原因
		if !fa[i+int(haveNum)] {
			if (i+int(haveNum)) > faEnd && (i+int(haveNum))%faAdd == 0 {
				fa[i+int(haveNum)] = true
			}
		}
		if fa[i+int(haveNum)] {
			beginPriority := 0
			rand := global.RandIn(SpecialPriority) + 1
			for _, prizeWheelitem := range PartnerCfg.SpecialPool {
				if beginPriority+int(prizeWheelitem.Priority) >= rand {
					// itemInfo, err := global.GetItemCfg(prizeWheelitem.Id)
					// if err != nil {
					// 	return nil, err
					// }
					// if itemInfo.Quality >= 5 {
					// 	strSend := []string{userInfo.Name, strconv.Itoa(int(itemInfo.ID))}
					// 	game.Lamp(int32(constant.LampHeroSummon), strSend)

					// }
					award = append(award, protocol.PartnerAward{
						Id:   prizeWheelitem.Id,
						Kind: prizeWheelitem.Kind,
						Num:  prizeWheelitem.Num,
					})
					break
				}
				beginPriority = beginPriority + int(prizeWheelitem.Priority)
			}
		} else {
			beginPriority := 0
			rand := global.RandIn(CommonPriority) + 1
			for _, prizeWheelitem := range PartnerCfg.CommonPool {
				if beginPriority+int(prizeWheelitem.Priority) >= rand {
					// itemInfo, err := global.GetItemCfg(prizeWheelitem.Id)
					// if err != nil {
					// 	return nil, nil
					// }
					// if itemInfo.Quality >= 5 {
					// 	strSend := []string{userInfo.Name, strconv.Itoa(int(itemInfo.ID))}
					// 	game.Lamp(int32(constant.LampHeroSummon), strSend)
					// }
					award = append(award, protocol.PartnerAward{
						Id:   prizeWheelitem.Id,
						Kind: prizeWheelitem.Kind,
						Num:  prizeWheelitem.Num,
					})
					break
				}
				beginPriority = beginPriority + int(prizeWheelitem.Priority)
			}
		}
	}
	return award, nil
}

// 宠物抽奖
func prizeDrawPetReceive(ctx context.Context, r *protocol.PrizedrawRequest) (*protocol.PrizedrawResponse, error) {
	request := &protocol.PrizedrawRequest{}
	//[{2024-05-09 2} {202405-10 1}]
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request.Uid = Uid

	if r.Type != 1 && r.Type != 2 {
		return &protocol.PrizedrawResponse{Code: -1}, nil
	}
	if r.PType != constant.LotteryTypeBronze && r.PType != constant.LotteryTypeGold && r.PType != constant.LotteryTypeSilver && r.PType != constant.LotteryTypePlatinum {
		return &protocol.PrizedrawResponse{Code: -1}, nil
	}
	request.Type = r.Type
	request.PType = r.PType
	_, award, err := PrizeDrawPetReceive(request)
	if err != nil {
		return &protocol.PrizedrawResponse{Code: -1, Data: err.Error()}, nil
	}
	return &protocol.PrizedrawResponse{Award: award}, nil
}

// 宠物抽奖抽奖
func PrizeDrawPetReceive(r *protocol.PrizedrawRequest) ([]protocol.SourceInfo, []protocol.PartnerAward, error) {
	//个人等级
	uerInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, nil, err
	}
	//活动内容
	res := global.GetSummon(constant.LotteryModulePet, r.PType, uerInfo.Level)
	if res == nil {
		return nil, nil, errutil.ErrConfigNotFound
	}

	// //需要领取的日期
	// needContent := ress % 100
	needNum := -1
	switch r.Type {
	case 1:
		needNum = 1
	case 2:
		needNum = 10
	}
	if needNum == -1 {
		return nil, nil, errutil.ErrConfigNotFound
	}
	//减去消耗品
	resultc := make([]protocol.SourceInfo, 1)
	for _, info := range res.Consume {
		if r.Type == info.Type {
			resultc[0].Id = info.Id
			resultc[0].Num = -int64(info.Num)
		}

	}
	//校验背包情况
	err = CheckSourceNum(resultc, r.Uid, false, constant.ItemPrizeDrawPet)
	if err != nil {
		//如果没有，再检测有没有点券
		resultc[0].Id = global.GlobalCfg.BaseValue.Petbuy.ItemId
		resultc[0].Num = -int64(global.GlobalCfg.BaseValue.Petbuy.NeedMoney)
		err = CheckSourceNum(resultc, r.Uid, false, constant.ItemPrizeDrawPet)
		if err != nil {
			return nil, nil, err
		}
	}

	haveNum, err := db.PrizedrawList(r.Uid, r.PType, constant.LotteryModulePet, 1)

	if err != nil {
		return nil, nil, err
	}
	petDetail, err := PetDetail(r.PType, uerInfo, needNum, int(haveNum))
	if err != nil {
		return nil, nil, err
	}

	result3 := protocol.AddSource{
		Uid:    r.Uid,
		Source: make([]*protocol.SourceInfo, 0),
	}

	for i, info := range petDetail {
		if info.Kind == constant.LotteryAwardWhole {
			result3.Source = append(result3.Source, &protocol.SourceInfo{
				Id:  petDetail[i].Id,
				Num: int64(petDetail[i].Num),
			})
		}
	}

	//给他加整卡
	//这时候给他卡，他的背包里面还是有道具的
	//所以同步消息的时候，数据可能是错的，比如他一共6个，用掉一个，还剩5个，但是他的背包里面有6个，那他还是要有倒计时的
	channel := "青铜"
	switch r.PType {
	case 1:
		channel = "青铜"
	case 2:
		channel = "白银"
	case 3:
		channel = "黄金"
	case 4:
		channel = "无双宠物卡兑换"
	default:
		channel = "青铜"
	}
	err = addPet(result3, int32(needNum), channel)
	if err != nil {
		return nil, nil, err
	}
	for _, info := range result3.Source {
		//查询这个宠物的品质
		petInfo := global.GetPetCfg(info.Id)
		if petInfo.Quality >= 6 {
			strSend := []string{uerInfo.Name, strconv.Itoa(int(petInfo.ItemID))}
			game.Lamp(int32(constant.LampPetSummon), uerInfo.ServerId, strSend)
		}
	}
	jsonDate, err := json.Marshal(petDetail)
	if err != nil {
		return nil, nil, err
	}

	jsonDateString := string(jsonDate)
	err = db.PrizedrawAdd(r.Uid, jsonDateString, int32(needNum), r.PType, constant.LotteryModulePet)
	if err != nil {
		return nil, nil, err
	}
	//校验背包情况
	err = CheckSourceNum(resultc, r.Uid, true, constant.ItemPrizeDrawPet)
	if err != nil {
		return nil, nil, err
	}
	// return int32(resulte[0].Num), ""
	sourceInfo := make([]protocol.SourceInfo, len(result3.Source))
	for i, v := range result3.Source {
		if v != nil {
			sourceInfo[i] = *v
		}
	}

	switch r.PType {
	case 1:
		game.TaskDispatch(uerInfo.UID, pb.TargetType_DRAW_GREEN_PET_EGGS, 0, int32(needNum))
	case 2:
		game.TaskDispatch(uerInfo.UID, pb.TargetType_DRAW_PURPLE_PET_EGGS, 0, int32(needNum))
	case 3:
		game.TaskDispatch(uerInfo.UID, pb.TargetType_DRAW_ORANGE_PET_EGGS, 0, int32(needNum))
	}

	return nil, petDetail, nil
}
func upPartnerShushu(sourceInfo []protocol.SourceInfo, partnerInfo []protocol.SourceInfo, costInfo []protocol.SourceInfo, uid int64, ptype int32) {
	data := make(map[string]interface{}, 0)
	data["name"] = "summon_parnter"
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
	for _, v := range partnerInfo {
		if v.Num > 0 {
			//查询itemname
			if PartnerCfg := global.GetPartnerCfg(v.Id); PartnerCfg != nil {
				get_item_info = append(get_item_info, map[string]interface{}{
					"item_id":   PartnerCfg.Cid,
					"item_name": PartnerCfg.Name,
					"item_num":  v.Num,
				})
			}
		}
	}
	data["get_item_info"] = get_item_info

	cost_item_info := make([]map[string]interface{}, 0)
	for _, v := range costInfo {
		//查询itemname
		if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
			item_type := "碎片"
			if ItemCfg.MainType == constant.ItemPartner {
				item_type = "整卡"
			}
			cost_item_info = append(cost_item_info, map[string]interface{}{
				"item_id":   ItemCfg.ID,
				"item_name": ItemCfg.Name,
				"item_num":  math.Abs(float64(v.Num)),
				"item_type": item_type,
			})
		}
	}
	data["cost_item_info"] = cost_item_info
	data["summon_type"] = ptype
	game.Upshushu(uid, data, true, true)
}

// 宠物抽奖抽奖
func PetDetail(ptype int32, userInfo *protocol.UserInfo, needNum int, haveNum int) ([]protocol.PartnerAward, error) {
	award := make([]protocol.PartnerAward, 0)

	PartnerCfg := global.GetSummon(constant.LotteryModulePet, ptype, userInfo.Level)
	if needNum+haveNum > int(PartnerCfg.MaxTime) && PartnerCfg.MaxTime != 0 {
		return award, errutil.ErrNotCanBuy
	}
	CommonPriority := 0

	//普通奖池
	for _, Partnercitem := range PartnerCfg.CommonPool {
		CommonPriority = CommonPriority + int(Partnercitem.Priority)
	}
	for i := 1; i <= needNum; i++ {
		beginPriority := 0
		rand := global.RandIn(CommonPriority) + 1
		for _, prizeWheelitem := range PartnerCfg.CommonPool {
			if beginPriority+int(prizeWheelitem.Priority) >= rand {
				award = append(award, protocol.PartnerAward{
					Id:   prizeWheelitem.Id,
					Kind: prizeWheelitem.Kind,
					Num:  prizeWheelitem.Num,
				})
				break
			}
			beginPriority = beginPriority + int(prizeWheelitem.Priority)
		}
	}
	return award, nil
}
