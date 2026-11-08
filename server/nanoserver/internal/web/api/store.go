package api

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"nanoserver/db"
	"nanoserver/db/model"
	"nanoserver/global"
	"nanoserver/internal/game"
	"nanoserver/internal/web/provide"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/tracking"
	"nanoserver/pkg/tracking/thinkingdatax"
	"nanoserver/protocol"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	"github.com/spf13/viper"
)

func MakeStoreService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/store/list", nex.Handler(storeList)).Methods("POST")                                // 列表
	router.Handle("/v1/store/buy", limitHandler(nex.Handler(storeBuy), "store_buy")).Methods("PUT")        // 购买
	router.Handle("/v1/store/extractaccount", nex.Handler(storeExtractAccount)).Methods("POST")            // 提取账号
	router.Handle("/v1/store/extractaccountupdate", nex.Handler(storeExtractAccountUpdate)).Methods("PUT") // 提取账号
	router.Handle("/v1/store/bounty", nex.Handler(bounty)).Methods("PUT")                                  // http导入数据
	return router
}

// 传入商城id，查找对应所有商品的限制情况
func StoreList(r *protocol.StoreListRequest) ([]protocol.StoreList, error) {
	storeCfg := global.StoreCfg
	globalStoreType := make(map[int32][]int64)
	globalgoods := make(map[int64]int64)
	globalgoodsAll := make(map[int64]int64)
	// 购买的总次数
	// result := make(map[int32][]protocol.SourceInfo)
	for _, storen := range storeCfg {
		if r.StoreType != 0 {
			if r.StoreType == storen.StoreType {
				LimitType := int32(storen.Limit.LimitType) // 假设 StoreType 是 int 类型的
				goodsUuid := int64(storen.Id)

				// 兼容单个商品的筛选
				if r.GoodsUuid > 0 {
					if r.GoodsUuid == storen.Id {
						// 检查 storeType 是否已经存在于map中，如果不存在则初始化
						if _, exists := globalStoreType[LimitType]; !exists {
							globalStoreType[LimitType] = make([]int64, 0, 10) // 初始化切片，10为切片的预分配容量
						}
						// 追加 goodsId 到对应的 storeType 切片中
						globalStoreType[LimitType] = append(globalStoreType[LimitType], goodsUuid)
						globalgoods[int64(storen.Id)] = 0
					}
				} else {
					// 检查 storeType 是否已经存在于map中，如果不存在则初始化
					if _, exists := globalStoreType[LimitType]; !exists {
						globalStoreType[LimitType] = make([]int64, 0, 10) // 初始化切片，10为切片的预分配容量
					}
					// 追加 goo`dsId 到对应的 storeType 切片中
					globalStoreType[LimitType] = append(globalStoreType[LimitType], goodsUuid)
					globalgoods[int64(storen.Id)] = 0
				}
			}
		} else {
			LimitType := int32(storen.Limit.LimitType) // 假设 StoreType 是 int 类型的
			goodsUuid := int64(storen.Id)

			// 兼容单个商品的筛选
			if r.GoodsUuid > 0 {
				if r.GoodsUuid == storen.Id {
					// 检查 storeType 是否已经存在于map中，如果不存在则初始化
					if _, exists := globalStoreType[LimitType]; !exists {
						globalStoreType[LimitType] = make([]int64, 0, 10) // 初始化切片，10为切片的预分配容量
					}
					// 追加 goodsId 到对应的 storeType 切片中
					globalStoreType[LimitType] = append(globalStoreType[LimitType], goodsUuid)
					globalgoods[int64(storen.Id)] = 0
				}
			} else {
				// 检查 storeType 是否已经存在于map中，如果不存在则初始化
				if _, exists := globalStoreType[LimitType]; !exists {
					globalStoreType[LimitType] = make([]int64, 0, 10) // 初始化切片，10为切片的预分配容量
				}
				// 追加 goodsId 到对应的 storeType 切片中
				globalStoreType[LimitType] = append(globalStoreType[LimitType], goodsUuid)
				globalgoods[int64(storen.Id)] = 0
			}
		}
	}
	// result := make([]int64, 100)
	for i, delstoren := range globalStoreType {
		if i != 999 {
			// result[i] = make([]model.Store, 0)
			day := 1
			switch i {
			case 0:
				day = 365 * 50
			case 1:
				day = 1
			case 2:
				day = 7
			case 3:
				day = 30
			case 4:
				day = 30 * 3
			}
			// 查询用户购买情况
			resultb, _ := db.StoreList(
				r.Uid,
				r.StoreType,
				int32(day),
				delstoren,
				false,
			)
			if len(resultb) > 0 {
				for _, resb := range resultb {
					// 查询初始购买次数
					globalgoods[resb.GoodsUuiD] = globalgoods[resb.GoodsUuiD] + resb.Num
				}
			}
			resultd, _ := db.StoreList(
				r.Uid,
				r.StoreType,
				int32(365*50),
				delstoren,
				false,
			)
			if len(resultd) > 0 {
				for _, resb := range resultd {
					// 查询初始购买次数
					globalgoodsAll[resb.GoodsUuiD] = globalgoodsAll[resb.GoodsUuiD] + resb.Num
				}
			}

			// result[1] = resultb
		} else {
			// 查询用户购买情况
			resultb, _ := db.StoreList(
				r.Uid,
				r.StoreType,
				365*50,
				delstoren,
				false,
			)
			if len(resultb) > 0 {
				for _, resb := range resultb {
					globalgoods[resb.GoodsUuiD] = globalgoods[resb.GoodsUuiD] + resb.Num
					globalgoodsAll[resb.GoodsUuiD] = globalgoodsAll[resb.GoodsUuiD] + resb.Num
				}
			}
		}
	}

	res := []protocol.StoreList{}
	for _, storen := range storeCfg {
		if r.StoreType == storen.StoreType && r.StoreType > 0 {
			if r.GoodsUuid > 0 {
				if r.GoodsUuid == storen.Id {
					Date := int32(0)
					if storen.StoreType == constant.ShopTypeOfferReward {
						if storen.Id == constant.ShopGoodsIdSpecial {
							Date = -1
						} else {
							res, err := ShopTypeOfferReward(storen.Id)
							if err != nil {
								return nil, err
							}
							Date = int32(res)
						}
					}
					buyCount := globalgoods[int64(storen.Id)]
					//buyCount是已经购买了
					//如果已经购买时1，那么已经准备按第二个去算了
					num := int64(0)
					for j := 1; j <= int(r.Num); j++ {
						if len(storen.Consumechange) > 0 {
							if buyCount+int64(j) == 1 {
								num += storen.Limit.Num
							} else {
								if buyCount+int64(j) > int64(len(storen.Consumechange)) {
									//选取最后一个
									num += int64(storen.Consumechange[len(storen.Consumechange)-1].Value)
								} else {
									num += int64(storen.Consumechange[buyCount+int64(j)-2].Value)
								}
							}

						} else {
							num += storen.Limit.Num
						}
					}
					del := protocol.StoreList{
						GoodsUuid:  int64(storen.Id),
						Buy:        buyCount,
						BuyAllNum:  globalgoodsAll[int64(storen.Id)],
						GoodsId:    storen.GoodsId,
						Num:        num,
						LimitNum:   storen.Limit.LimitNum,
						NeedMoney:  int32(storen.Limit.NeedMoney),
						PreLimit:   storen.PreLimit,
						SourceInfo: make([]protocol.SourceInfo, len(storen.Props)),
						Date:       Date,
						Times:      storen.Limit.Times,
						Entrance:   storen.Entrance,
						Method:     storen.Method,
					}
					for i, storen2 := range storen.Props {
						del.SourceInfo[i] = protocol.SourceInfo{
							Id:  storen2.Id,
							Num: int64(storen2.Num),
						}
					}
					res = append(res, del)
				}
			} else {
				Date := int32(0)
				if storen.StoreType == constant.ShopTypeOfferReward {
					if storen.Id == constant.ShopGoodsIdSpecial {
						Date = -1
					} else {
						res, err := ShopTypeOfferReward(storen.Id)
						if err != nil {
							return nil, err
						}
						Date = int32(res)
					}
				}
				buyCount := globalgoods[int64(storen.Id)]
				num := int64(0)
				for j := 1; j <= int(r.Num); j++ {
					if len(storen.Consumechange) > 0 {
						if buyCount+int64(j) == 1 {
							num += storen.Limit.Num
						} else {
							if buyCount+int64(j) > int64(len(storen.Consumechange)) {
								//选取最后一个
								num += int64(storen.Consumechange[len(storen.Consumechange)-1].Value)
							} else {
								num += int64(storen.Consumechange[buyCount+int64(j)-2].Value)
							}
						}

					} else {
						num += storen.Limit.Num
					}
				}
				del := protocol.StoreList{
					GoodsUuid:  int64(storen.Id),
					Buy:        buyCount,
					BuyAllNum:  globalgoodsAll[int64(storen.Id)],
					GoodsId:    storen.GoodsId,
					Num:        num,
					LimitNum:   storen.Limit.LimitNum,
					NeedMoney:  int32(storen.Limit.NeedMoney),
					PreLimit:   storen.PreLimit,
					SourceInfo: make([]protocol.SourceInfo, len(storen.Props)),
					Date:       Date,
					Times:      storen.Limit.Times,
					Entrance:   storen.Entrance,
					Method:     storen.Method,
				}
				for i, storen2 := range storen.Props {
					del.SourceInfo[i] = protocol.SourceInfo{
						Id:  storen2.Id,
						Num: int64(storen2.Num),
					}
				}
				res = append(res, del)
			}
		}
		if r.StoreType == 0 {
			if r.GoodsUuid > 0 {
				if r.GoodsUuid == storen.Id {
					Date := int32(0)
					if storen.StoreType == constant.ShopTypeOfferReward {
						if storen.Id == constant.ShopGoodsIdSpecial {
							Date = -1
						} else {
							res, err := ShopTypeOfferReward(storen.Id)
							if err != nil {
								return nil, err
							}
							Date = int32(res)
						}
					}
					buyCount := globalgoods[int64(storen.Id)]
					num := int64(0)
					for j := 1; j <= int(r.Num); j++ {
						if len(storen.Consumechange) > 0 {
							if buyCount+int64(j) == 1 {
								num += storen.Limit.Num
							} else {
								if buyCount+int64(j) > int64(len(storen.Consumechange)) {
									//选取最后一个
									num += int64(storen.Consumechange[len(storen.Consumechange)-1].Value)
								} else {
									num += int64(storen.Consumechange[buyCount+int64(j)-2].Value)
								}
							}

						} else {
							num += storen.Limit.Num
						}
					}
					del := protocol.StoreList{
						GoodsUuid:  int64(storen.Id),
						Buy:        buyCount,
						BuyAllNum:  globalgoodsAll[int64(storen.Id)],
						GoodsId:    storen.GoodsId,
						Num:        num,
						LimitNum:   storen.Limit.LimitNum,
						NeedMoney:  int32(storen.Limit.NeedMoney),
						PreLimit:   storen.PreLimit,
						SourceInfo: make([]protocol.SourceInfo, len(storen.Props)),
						Date:       Date,
						Times:      storen.Limit.Times,
						Entrance:   storen.Entrance,
						Method:     storen.Method,
					}
					for i, storen2 := range storen.Props {
						del.SourceInfo[i] = protocol.SourceInfo{
							Id:  storen2.Id,
							Num: int64(storen2.Num),
						}
					}
					res = append(res, del)
				}
			} else {
				buyCount := globalgoods[int64(storen.Id)]
				num := int64(0)
				for j := 1; j <= int(r.Num); j++ {
					if len(storen.Consumechange) > 0 {
						if buyCount+int64(j) == 1 {
							num += storen.Limit.Num
						} else {
							if buyCount+int64(j) > int64(len(storen.Consumechange)) {
								//选取最后一个
								num += int64(storen.Consumechange[len(storen.Consumechange)-1].Value)
							} else {
								num += int64(storen.Consumechange[buyCount+int64(j)-2].Value)
							}
						}

					} else {
						num += storen.Limit.Num
					}
				}
				Date := int32(0)
				if storen.StoreType == constant.ShopTypeOfferReward {
					if storen.Id == constant.ShopGoodsIdSpecial {
						Date = -1
					} else {
						res, err := ShopTypeOfferReward(storen.Id)
						if err != nil {
							return nil, err
						}
						Date = int32(res)
					}
				}
				del := protocol.StoreList{
					GoodsUuid:  int64(storen.Id),
					Buy:        buyCount,
					BuyAllNum:  globalgoodsAll[int64(storen.Id)],
					GoodsId:    storen.GoodsId,
					Num:        num,
					LimitNum:   storen.Limit.LimitNum,
					NeedMoney:  int32(storen.Limit.NeedMoney),
					PreLimit:   storen.PreLimit,
					SourceInfo: make([]protocol.SourceInfo, len(storen.Props)),
					Date:       Date,
					Times:      storen.Limit.Times,
					Entrance:   storen.Entrance,
					Method:     storen.Method,
				}
				for i, storen2 := range storen.Props {
					del.SourceInfo[i] = protocol.SourceInfo{
						Id:  storen2.Id,
						Num: int64(storen2.Num),
					}
				}
				res = append(res, del)
			}
		}

	}

	return res, nil
}

func ShopTypeOfferReward(goodsUuid int32) (int64, error) {
	return 0, nil

	bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
	key := bounty + strconv.Itoa(int(goodsUuid))
	cmd := "LLEN"
	res2, err := provide.CmdCache(cmd, key)
	if err != nil {
		return 0, err
	}
	OwnScore, err := redis.Int64(res2, nil)
	if err != nil {
		return 0, err
	}
	return OwnScore, nil
}

func storeList(ctx context.Context, r *protocol.StoreListRequest) (*protocol.StoreListResponse, error) {
	request := &protocol.StoreListRequest{}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	//StoreType, err1 := strconv.ParseInt(r.StoreType, 10, 32)
	//if err1 != nil {
	//	return nil, err1
	//}

	request.Uid = Uid
	request.StoreType = r.StoreType

	list, err := StoreList(request)
	res := make([]protocol.StoreListLittle, len(list))
	for i, storen := range list {
		res[i].GoodsUuid = storen.GoodsUuid
		res[i].Buy = storen.Buy
		res[i].BuyAllNum = storen.BuyAllNum
		res[i].Date = storen.Date
	}
	if err != nil {
		return nil, err
	}
	// fmt.Println("b", list)
	return &protocol.StoreListResponse{Data: res}, nil
}

func storeBuy(ctx context.Context, data *protocol.StoreListRequest, httprequest *http.Request) (*protocol.StoreOperateResponse, error) {
	logger.Debugf("====================storeBuy: %+#v", data)
	request := &protocol.StoreListRequest{
		Num:        data.Num,
		Notify:     data.Notify,
		RequestID:  data.RequestID,
		Extra:      data.Extra,
		Platform:   data.Platform,
		RefineSlot: data.RefineSlot,
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	goods_uuid := data.GoodsUuid
	// request.Num = data.Num
	// request.Notify = data.Notify
	// request.RequestID = data.RequestID
	if goods_uuid == 0 || data.Num <= 0 {
		return nil, errutil.ErrInvalidParameter
	}

	// request.Extra = data.Extra
	// request.Platform = data.Platform

	request.Uid = Uid
	//// 使用整数除法得到 5
	//request.StoreType = int32(goods_id / 1000)
	//// 使用取余操作得到 10
	//request.GoodsId = goods_id % 1000
	request.GoodsUuid = goods_uuid
	request.Ip = ip(httprequest.RemoteAddr)
	resp, err := StoreBuy(request)
	if err != nil {
		return &protocol.StoreOperateResponse{Code: -1, Data: false, Error: err.Error()}, nil
	}

	if resp != nil {
		return resp, nil
	}

	return &protocol.StoreOperateResponse{Code: 0, Data: true, Error: ""}, nil
}

// 判断是否能够购买
func StoreBuyCan(request *protocol.StoreListRequest, userInfo *protocol.UserInfo) ([]protocol.StoreList, int32, error) {
	// 查询目前的情况
	// request.StoreType = r.StoreType
	// request.GoodsId = r.GoodsId
	list, err := StoreList(request)
	buyNum := int32(1)
	if err != nil || list == nil {
		return nil, buyNum, err
	}
	// 校验有没有购买次数
	if len(list) == 0 {
		return nil, buyNum, errutil.ErrLevel
	}
	if list[0].LimitNum != 0 {
		if request.Num != 1 {
			buyNum = int32(request.Num)
		}
		if int32(list[0].Buy)+buyNum > int32(list[0].LimitNum) {
			return nil, buyNum, errutil.ErrNotCanBuy
		}
	}
	// 如果是悬赏，需要对比vip
	if request.StoreType == constant.ShopTypeOfferReward && request.GoodsUuid != constant.ShopGoodsIdSpecial {
		if userInfo.Vip < list[0].PreLimit {
			return nil, buyNum, errutil.ErrLevel
		}
		bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
		key := bounty + strconv.Itoa(int(request.GoodsUuid))
		cmd := "LLEN"
		res2, err := provide.CmdCache(cmd, key)
		if err != nil {
			return nil, buyNum, err
		}
		OwnScore, err := redis.Int64(res2, nil)
		if err != nil {
			return nil, buyNum, err
		}
		if OwnScore < 1 {
			return nil, buyNum, errutil.ErrNotCanBuy
		}
	}
	// 首充
	if request.StoreType == constant.ShopTypeFirstCharge {

		////校验当天是不是登录
		storeCfg := global.ActivityCfg
		// 活动内容
		actId := int32(0)
		typeId := int32(0)
		conId := int32(0)
		GoodsId := make([]int64, 0)
		for _, storen := range storeCfg {
			if storen.PreLimit == request.GoodsUuid {
				actId = storen.Acttag
				typeId = storen.Type
				conId = storen.Content
			}
		}
		if conId != 1 {
			for _, storen := range storeCfg {
				if storen.Type == typeId && storen.Acttag == actId && storen.Content == 1 {
					GoodsId = append(GoodsId, int64(storen.PreLimit))
				}
			}

			// request.StoreType = r.StoreType
			// request.GoodsId = r.GoodsId
			if len(GoodsId) == 0 {
				return nil, buyNum, errutil.ErrChooseRound
			}
			StoType := int32(GoodsId[0] / 1000)

			resultb, _ := db.StoreList(
				request.Uid,
				StoType,
				int32(365*50),
				GoodsId,
				false,
			)
			if len(resultb) == 0 {
				return nil, buyNum, errutil.ErrChooseRound
			}
			// 查询有没有购买过1，且距离1的天数对不对

			endTime := time.Now()
			// resultb[0].CreateTime当天凌晨的时间
			createdTime := time.Unix(int64(resultb[0].CreateTime), 0)
			createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
			// 计算时间差
			duration := endTime.Sub(createdDayMidnight)
			// 计算相差的天数，即目前第几天
			days := int64(duration.Hours()/24) + 1
			if days != int64(conId) {
				return nil, buyNum, errutil.ErrChooseRound
			}
		}

		// if err != nil {
		// 	return nil, buyNum, err
		// }
		// strSend := []string{userInfo.Name}
		// game.Lamp(int32(constant.LampFirstRechargeReward), strSend)
	}

	// logger.Debugf("购买商品", uid)
	return list, buyNum, nil
}

func StoreBuy(r *protocol.StoreListRequest) (resp *protocol.StoreOperateResponse, err error) {
	NeedMoney := int32(constant.NeedMoney)
	NeedAd := int32(constant.NeedAd)
	// 查询目前的情况
	request := &protocol.StoreListRequest{}
	request.Uid = r.Uid
	request.Num = r.Num
	request.GoodsUuid = r.GoodsUuid
	request.StoreType = r.GoodsUuid / 1000 // 根据商品ID计算商店类型

	// createTime := time.Now().Unix()
	// 个人信息
	userInfo, err := game.UserInfo(r.Uid)
	if err != nil {
		return nil, err
	}
	if userInfo.Status > int(game.PlayerStatusPlaying) {
		return nil, errutil.ErrReLogin
	}

	list, err := StoreList(request)

	buyNum := int32(1)
	// 校验有没有购买次数
	if len(list) == 0 {
		return nil, errutil.ErrLevel
	}
	if list[0].LimitNum != 0 {
		if request.Num != 1 {
			buyNum = int32(request.Num)
		}
		if int32(list[0].Buy)+buyNum > int32(list[0].LimitNum) {
			return nil, errutil.ErrNotCanBuy
		}
	}
	buy_name := list[0].Entrance
	buy_entrance := list[0].Method
	if list[0].PreLimit != 0 {
		//查询前置商品.
		request2 := &protocol.StoreListRequest{}
		request2.Uid = r.Uid
		request2.Num = r.Num
		request2.GoodsUuid = list[0].PreLimit
		request2.StoreType = list[0].PreLimit / 1000
		_, err := StoreList(request2)
		if err != nil {
			return nil, err
		}

	}
	// 如果是试炼商城，需要对比当前的层数
	if request.StoreType == constant.ShopTypeTrialReward {
		if userInfo.TrialTowerLevel < list[0].PreLimit {
			return nil, errutil.ErrLevel
		}
	}

	// 如果是悬赏，需要对比vip
	if request.StoreType == constant.ShopTypeOfferReward && request.GoodsUuid != constant.ShopGoodsIdSpecial {
		if userInfo.Vip < list[0].PreLimit {
			return nil, errutil.ErrLevel
		}
		bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
		key := bounty + strconv.Itoa(int(request.GoodsUuid))
		cmd := "LLEN"
		res2, err := provide.CmdCache(cmd, key)
		if err != nil {
			return nil, err
		}
		OwnScore, err := redis.Int64(res2, nil)
		if err != nil {
			return nil, err
		}
		if OwnScore < 1 {
			return nil, errutil.ErrNotCanBuy
		}
	}

	// 首充
	if request.StoreType == constant.ShopTypeFirstCharge {
		////校验当天是不是登录
		storeCfg := global.ActivityCfg
		// 活动内容
		actId := int32(0)
		typeId := int32(0)
		conId := int32(0)
		GoodsId := make([]int64, 0)
		for _, storen := range storeCfg {
			if storen.PreLimit == request.GoodsUuid {
				actId = storen.Acttag
				typeId = storen.Type
				conId = storen.Content
			}
		}
		if conId != 1 {
			for _, storen := range storeCfg {
				if storen.Type == typeId && storen.Acttag == actId && storen.Content == 1 {
					GoodsId = append(GoodsId, int64(storen.PreLimit))
				}
			}

			// request.StoreType = r.StoreType
			// request.GoodsId = r.GoodsId
			if len(GoodsId) == 0 {
				return nil, errutil.ErrChooseRound
			}
			StoType := int32(GoodsId[0] / 1000)

			resultb, _ := db.StoreList(
				request.Uid,
				StoType,
				int32(365*50),
				GoodsId,
				false,
			)
			if len(resultb) == 0 {
				return nil, errutil.ErrChooseRound
			}
			endTime := time.Now()
			// resultb[0].CreateTime当天凌晨的时间
			createdTime := time.Unix(int64(resultb[0].CreateTime), 0)
			createdDayMidnight := time.Date(createdTime.Year(), createdTime.Month(), createdTime.Day(), 0, 0, 0, 0, createdTime.Location())
			// 计算时间差
			duration := endTime.Sub(createdDayMidnight)
			// 计算相差的天数，即目前第几天
			days := int64(duration.Hours()/24) + 1
			if days != int64(conId) {
				return nil, errutil.ErrChooseRound
			}
		}

		if err != nil {
			return nil, err
		}
	}

	if request.StoreType == constant.ShopTypeRefine {
		if r.RefineSlot != nil {
			// 判断是否买过，如果买过则直接返回
			issPurchased, err := game.RefineIsPurchased(r.Uid, int32(r.RefineSlot.CharacterSlot), int32(r.RefineSlot.EquipSlot), int32(r.RefineSlot.RefinedSlot))
			if err != nil {
				return nil, err
			}

			if issPurchased {
				logger.Warnf("refine solt is purchased uid %d, character %d, equip %d, refine %d", r.Uid, r.RefineSlot.CharacterSlot, r.RefineSlot.EquipSlot, r.RefineSlot.RefinedSlot)
				return nil, errutil.ErrNotCanBuy
			}

			if err := game.SetTempPurchasedSlot(r.Uid, int32(r.RefineSlot.CharacterSlot), int32(r.RefineSlot.EquipSlot), int32(r.RefineSlot.RefinedSlot)); err != nil {
				return nil, err
			}
		}
	}

	if list[0].NeedMoney == NeedMoney {
		if list[0].Times != 0 {
			if int32(list[0].BuyAllNum) < list[0].Times {
				if request.Num > 1 {
					return nil, errutil.ErrNotCanBuy
				}
				if list[0].Times >= int32(list[0].BuyAllNum)+int32(request.Num) {
					appEnv := strings.ToLower(viper.GetString("core.env"))
					if appEnv == constant.EnvProd || appEnv == constant.EnvReview || appEnv == constant.EnvStag {
						createOrderRequest := &protocol.CreateOrderRequest{
							ProductId:      int(request.GoodsUuid),
							ProductionName: strconv.FormatInt(int64(request.GoodsUuid), 10),
							Device: protocol.Device{
								Remote: ip(r.Ip),
							},
							Uid:       r.Uid,
							Extra:     r.Extra,
							Platform:  r.Platform,
							Amount:    int(math.Ceil(float64(list[0].NeedMoney) * float64(request.Num))),
							OpenID:    userInfo.AccountName,
							ServerID:  uint64(userInfo.ServerId),
							RequestID: r.RequestID,
						}
						data, err := CreateOrderV2(createOrderRequest)
						if err != nil {
							return nil, err
						}
						upStoreInitOrderShushu(buy_name, buy_entrance, data.PrePayInternalID, int(math.Ceil(float64(list[0].Num)*float64(request.Num))), r.Uid)
						logger.Debugf("storeBuy CreateOrderV2 data: %v", data)
						return data, nil
					}
				}
			}
		} else {
			appEnv := strings.ToLower(viper.GetString("core.env"))
			if appEnv == constant.EnvProd || appEnv == constant.EnvReview || appEnv == constant.EnvStag {
				createOrderRequest := &protocol.CreateOrderRequest{
					ProductId:      int(request.GoodsUuid),
					ProductionName: strconv.FormatInt(int64(request.GoodsUuid), 10),
					Device: protocol.Device{
						Remote: ip(r.Ip),
					},
					Uid:       r.Uid,
					Extra:     r.Extra,
					Platform:  r.Platform,
					Amount:    int(math.Ceil(float64(list[0].NeedMoney) * float64(request.Num))),
					OpenID:    userInfo.AccountName,
					ServerID:  uint64(userInfo.ServerId),
					RequestID: r.RequestID,
				}
				data, err := CreateOrderV2(createOrderRequest)
				if err != nil {
					return nil, err
				}
				upStoreInitOrderShushu(buy_name, buy_entrance, data.PrePayInternalID, int(math.Ceil(float64(list[0].Num)*float64(request.Num))), r.Uid)
				logger.Debugf("storeBuy CreateOrderV2 data: %v", data)
				return data, nil
			}
		}

		// host, _, err := net.SplitHostPort(cfg.Addr)
		// userIp := ""
		// if err == nil {
		// 	userIp = host
		// }
		// 正式环境和提审环境需要打开

	}

	// 是否是人民币购买
	cost_item_before := int64(0)
	get_item_before := int64(0)
	////发放签到的奖励
	result := make([]protocol.SourceInfo, 0)

	// 减去消耗品
	// 这里人民币支付需要拉起支付，看广告的，需要看广告，
	// 还有个特殊情况就是，加了新功能，购买X次之后，就免费了，所以需要判断一下
	if list[0].NeedMoney == NeedMoney {
		// 累加vip等级
		// game.Recharge(request.Uid, list[0].Num)
	} else if list[0].NeedMoney == NeedAd {
		// 后续会接入广告系统
	} else {
		// 减去消耗品
		if list[0].Times != 0 {
			// 如果有购买n次的条件，且已经满足购买n次的条件就免费
			if int32(list[0].BuyAllNum) < list[0].Times {
				if list[0].Times >= int32(list[0].BuyAllNum)+int32(request.Num) {
					for _, info := range list {
						result = append(result, protocol.SourceInfo{
							Id:  info.NeedMoney,
							Num: -info.Num,
						})
					}
				} else {
					for j := 1; j <= int(list[0].Times)-int(list[0].BuyAllNum); j++ {
						for _, info := range list {
							result = append(result, protocol.SourceInfo{
								Id:  info.NeedMoney,
								Num: -info.Num,
							})
						}
					}
				}

				if result[0].Num != 0 && result[0].Id != 0 {
					cost_item_before = int64(game.GetBagNum(r.Uid, result[0].Id))
					err2 := CheckSourceNum(result, request.Uid, true, list[0].Method+"-"+list[0].Entrance)
					if err2 != nil {
						return nil, err2
					}
				}
			}
		} else {
			for _, info := range list {
				result = append(result, protocol.SourceInfo{
					Id: info.NeedMoney,

					Num: -info.Num,
				})
			}
			if result[0].Num != 0 && result[0].Id != 0 {
				cost_item_before = int64(game.GetBagNum(r.Uid, result[0].Id))
				err2 := CheckSourceNum(result, request.Uid, true, list[0].Method+"-"+list[0].Entrance)
				if err2 != nil {
					return nil, err2
				}
			}
		}
	}

	// 校验背包情况
	result2 := make([]protocol.SourceInfo, len(list[0].SourceInfo))
	bei := int64(request.Num)
	//首冲双倍 点券商城
	if request.StoreType == constant.ShopTypeRecharge && list[0].Buy == 0 {
		bei = int64(request.Num) + 1
	}
	for i, info := range list[0].SourceInfo {
		// get_item_before = int64(game.GetBagNum(r.Uid, info.Id))
		result2[i].Id = info.Id
		result2[i].Num = info.Num * bei
		result2[i].ShopType = request.StoreType
	}

	// 之前默认是背包满了才发邮件
	// 现在多一个判断，有些是必须发邮件
	// if request.StoreType == constant.ShopTypeSubscription {
	// 	MailId := strconv.FormatInt(snowflake.GenID(), 10)
	// 	result3 := make([]protocol.Attachment, len(result2))
	// 	for i, Attachment := range result2 {
	// 		result3[i] = protocol.Attachment{
	// 			MailId:      MailId,
	// 			UserPropsId: int64(Attachment.Id),
	// 			Num:         Attachment.Num,
	// 		}
	// 	}
	// 	request3 := &protocol.MailAddRequest{
	// 		Uid:          -1,
	// 		RecipientUid: r.Uid,
	// 		Title:        "系统邮件：特权卡奖励",
	// 		Desc:         "特权卡购买立返奖励",
	// 		Attachment:   result3,
	// 		ChannelCode:  "0",
	// 	}
	// 	resultc, err := Mailadd(request3)
	// 	if !resultc {
	// 		return nil, err
	// 	}
	// } else {
	// 	err3 := CheckSourceAdd(result2, request.Uid, true, list[0].Method+"-"+list[0].Entrance)
	// 	if err3 != nil {
	// 		if err3 == errutil.ErrPlayerNotFound {
	// 			return nil, err3
	// 		}
	// 		if err3 == errutil.ErrBagFull {

	// 			MailId := strconv.FormatInt(snowflake.GenID(), 10)
	// 			result3 := make([]protocol.Attachment, len(result2))
	// 			for i, Attachment := range result2 {
	// 				result3[i] = protocol.Attachment{
	// 					MailId:      MailId,
	// 					UserPropsId: int64(Attachment.Id),
	// 					Num:         Attachment.Num,
	// 				}
	// 			}

	// 			request3 := &protocol.MailAddRequest{
	// 				Uid:          request.Uid,
	// 				RecipientUid: -1,
	// 				Title:        "系统邮件",
	// 				Desc:         "包裹已满，请及时清理",
	// 				Attachment:   result3,
	// 				ChannelCode:  "0",
	// 			}
	// 			resultc, err3 := Mailadd(request3)
	// 			if !resultc || err3 != nil {
	// 				return nil, err3
	// 			}
	// 		}
	// 	}
	// }

	// if request.StoreType == constant.ShopTypeSubscription {
	// 	privilegeCfg := global.PrivilegeCfg
	// 	Name := ""
	// 	for _, privilegen := range privilegeCfg {
	// 		if privilegen.Sid == request.GoodsUuid {
	// 			Name = privilegen.Pname
	// 			originalTime := time.Now()
	// 			nextDayMidnight := originalTime.AddDate(0, 0, int(privilegen.Duration)) // 添加N天
	// 			nextDayMidnight = time.Date(nextDayMidnight.Year(), nextDayMidnight.Month(), nextDayMidnight.Day(), 0, 0, 0, 0, nextDayMidnight.Location())
	// 			err := game.PrivilegeAdd(request.StoreType, request.Uid, privilegen.Pid, nextDayMidnight)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 		}
	// 	}

	// 	strSend := []string{userInfo.Name, Name}
	// 	game.Lamp(int32(constant.LampPrivilegeCardPurchased), userInfo.ServerId, strSend)
	// }

	// if request.StoreType == constant.ShopTypeAFK {
	// 	err := game.PrivilegeAdd(request.StoreType, request.Uid, request.GoodsUuid, time.Time{})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// // 购买首冲走跑马灯
	// if request.StoreType == constant.ShopTypeFirstCharge {
	// 	strSend := []string{userInfo.Name}
	// 	game.Lamp(int32(constant.LampFirstRecharge), userInfo.ServerId, strSend)
	// }

	// // 购买试验基金走跑马灯
	// if request.StoreType == constant.ShopTypeTrialReward {
	// 	strSend := []string{userInfo.Name}
	// 	game.Lamp(int32(constant.LampTrialFundPurchased), userInfo.ServerId, strSend)
	// }

	// if request.StoreType == constant.ShopTypeGuardIsland {
	// 	// 这里需要将 Store.ID 转换为 config.xml 中的 GuardIsland.id
	// 	var (
	// 		guardIslandRecrodID int64
	// 		found               bool
	// 	)
	// 	for _, v := range global.GuardIslandCfg {
	// 		if v.ID == request.GoodsUuid {
	// 			guardIslandRecrodID = int64(v.ID)
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		logger.Errorf("StoreBuy GuardIslandMultiplier not found, uid: %d, store.uuid: %d", request.Uid, request.GoodsUuid)
	// 		return nil, errutil.ErrNotFound
	// 	}
	// 	if err := game.GuardIslandMultiplier(request.Uid, guardIslandRecrodID); err != nil {
	// 		logger.Errorf("StoreBuy GuardIslandMultiplier error: %v, uid: %d, store.uuid: %d", err, request.Uid, request.GoodsUuid)
	// 	}
	// }

	// jsonDate, err := json.Marshal(result2)
	// if err != nil {
	// 	return nil, err
	// }

	// jsonDateString := string(jsonDate)
	// // 插入购买记录表
	// resultb, err := db.StoreAdd(
	// 	request.Uid,
	// 	request.GoodsUuid,
	// 	jsonDateString,
	// 	buyNum,
	// 	createTime,
	// )
	// if !resultb {
	// 	return nil, err
	// }

	// if request.StoreType == constant.ShopTypeOfferReward {
	// 	if request.GoodsUuid != constant.ShopGoodsIdSpecial {
	// 		bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
	// 		key := bounty + strconv.Itoa(int(request.GoodsUuid))
	// 		cmd := "RPOP"
	// 		res2, err := provide.CmdCache(cmd, key)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		OwnScore, err := redis.Int64(res2, nil)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		logger.Debug("error:", OwnScore)
	// 	}
	// 	// todo 需要给这个用户发钱
	// 	storeExtractAccount, err := StoreExtractAccount(request.Uid)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	money := float64(0)
	// 	sendId := int32(0)
	// 	for _, info := range list[0].SourceInfo {
	// 		sendId = info.Id
	// 	}
	// 	for i := len(global.GlobalCfg.Redbag) - 1; i >= 0; i-- {
	// 		if sendId == int32(global.GlobalCfg.Redbag[i].Id) {
	// 			money = float64(global.GlobalCfg.Redbag[i].Value) / 100
	// 		}
	// 	}
	// 	order := protocol.FundTransParams{
	// 		AliName:     storeExtractAccount[0].Name,
	// 		Identity:    storeExtractAccount[0].Detail,
	// 		OrderTitle:  strconv.Itoa(int(sendId)),
	// 		OutBizNo:    strconv.Itoa(int(createTime)) + strconv.Itoa(int(request.Uid)),
	// 		TransAmount: strconv.FormatFloat(money, 'f', 2, 64),
	// 	}
	// 	err = provider.Byrtypay.FundTransUniTransfer(order)
	// 	if err != nil {
	// 		if request.GoodsUuid != constant.ShopGoodsIdSpecial {
	// 			bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
	// 			key := bounty + strconv.Itoa(int(request.GoodsUuid))
	// 			cmd := "RPUSH"
	// 			_, err := provide.CmdCache(cmd, key)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 		}
	// 		return nil, err
	// 	}
	// }

	// if request.StoreType == constant.ShopTypeSevenDaysFun {
	// 	game.TaskDispatch(request.Uid, pb.TargetType_SEVENDAY_TASK_GIFT_PURCHASED, 0, 1)
	// }

	_, err = NewPayUseCase().PayNotify(&protocol.PayNotifyRequest{
		Uid:     uint64(request.Uid),
		StoreId: int64(request.GoodsUuid),
		Num:     request.Num,
	})
	if err != nil {
		logger.Errorf("StoreBuy PayNotify error: %v, uid: %d, store.uuid: %d", err, request.Uid, request.GoodsUuid)
		return nil, errutil.ErrServerInternal
	}

	buy_number := list[0].BuyAllNum
	// 	上报
	// 没有库存就发邮件
	if list[0].NeedMoney == NeedAd {
		upAdClickShushu(buy_name, buy_entrance, r.Uid)
		// 计数增加
		count, err := game.AddAdCount(r.Uid)
		if err != nil {
			logger.Errorf("StoreBuy AddAdCount error: %v, uid: %d", err, r.Uid)
		}
		// 设置属性
		if err == nil {
			thinkingdatax.TrackerIns.Set(strconv.FormatInt(r.Uid, 10), userInfo.AccountName, tracking.M{
				"ad_total_today": count,
			})
		}
	}
	if list[0].NeedMoney != NeedMoney {
		is_ad := list[0].NeedMoney == NeedAd
		resultt := protocol.SourceInfo{}
		if result != nil {
			if len(result) > 0 {
				resultt = result[0]
			}
		}
		if len(result2) > 0 {
			if len(result2) > 1 {
				// 有多个物品
				upStorePkgShushu(buy_name, buy_entrance, result2, resultt, r.Uid, cost_item_before, is_ad, buy_number)
			} else {
				upStoreOneShushu(buy_name, buy_entrance, result2[0], resultt, r.Uid, int32(r.GoodsId), cost_item_before, get_item_before, is_ad)
			}
		}
	}

	return nil, nil
}

//	func upStoreInitOrderShushu(buy_name int32, buy_entrance int32, order_id string, pay_amount int, uid int64) {
//		data := make(map[string]interface{}, 0)
//		data["name"] = "init_order"
//		data["order_id"] = order_id
//		data["pay_amount"] = pay_amount
//		data["payment_entrance"] = buy_name
//		data["payment_name"] = buy_entrance
//		game.Upshushu(uid, data, true, false)
//	}
func upStoreInitOrderShushu(buy_name string, buy_entrance string, order_id string, pay_amount int, uid int64) {
	data := make(map[string]interface{}, 0)
	data["name"] = "init_order"
	data["order_id"] = order_id
	data["pay_amount"] = pay_amount
	data["payment_entrance"] = buy_name
	data["payment_name"] = buy_entrance
	game.Upshushu(uid, data, true, false)
}

func upStorePkgShushu(buy_name string, buy_entrance string, sourceInfo []protocol.SourceInfo, costInfo protocol.SourceInfo, uid int64, cost_item_before int64, is_ad bool, buy_number int64) {
	data := make(map[string]interface{}, 0)
	data["name"] = "shop_buy_package"
	data["buy_name"] = buy_name
	data["buy_entrance"] = buy_entrance
	get_item_info := make([]map[string]interface{}, 0)
	for _, v := range sourceInfo {
		// 查询itemname
		if ItemCfg, err := global.GetItemCfg(v.Id); err == nil {
			get_item_info = append(get_item_info, map[string]interface{}{
				"item_id":   ItemCfg.ID,
				"item_name": ItemCfg.Name,
				"item_num":  math.Abs(float64(v.Num)),
			})
		}
	}
	data["is_ad"] = is_ad
	data["buy_number"] = buy_number
	data["get_item_info"] = get_item_info
	data["cost_item_before"] = cost_item_before
	data["cost_item_after"] = cost_item_before + costInfo.Num
	if ItemCfg, err := global.GetItemCfg(costInfo.Id); err == nil {
		data["cost_item_id"] = ItemCfg.ID
		data["cost_item_name"] = ItemCfg.Name
		data["cost_item_num"] = math.Abs(float64(costInfo.Num))
	}

	game.Upshushu(uid, data, true, false)
}

func upStoreOneShushu(buy_name string, buy_entrance string, sourceInfo protocol.SourceInfo, costInfo protocol.SourceInfo, uid int64, cfgId int32, cost_item_before int64, get_item_before int64, is_ad bool) {
	data := make(map[string]interface{}, 0)
	data["name"] = "shop_buy_item"
	data["buy_name"] = buy_name
	data["buy_entrance"] = buy_entrance
	data["is_ad"] = is_ad
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
	game.Upshushu(uid, data, true, false)
}

func upAdClickShushu(buy_name string, buy_entrance string, uid int64) {
	data := make(map[string]interface{}, 0)
	data["name"] = "ad_click"
	data["ad_scene"] = buy_entrance
	data["ad_position"] = buy_name
	game.Upshushu(uid, data, true, true)
}

func storeExtractAccount(ctx context.Context) (*protocol.StoreExtractAccountResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	list, err := StoreExtractAccount(Uid)
	if err != nil {
		return &protocol.StoreExtractAccountResponse{Code: -1, Data: err.Error()}, nil
	}
	if len(list) > 0 {
		return &protocol.StoreExtractAccountResponse{Type: list[0].Type, Detail: list[0].Detail, Name: list[0].Name}, nil
	}
	return &protocol.StoreExtractAccountResponse{Type: 0, Detail: ""}, nil
}

func StoreExtractAccount(uid int64) ([]model.UserExtractAccount, error) {
	// 插入购买记录表
	resultb, err := db.ExtractAccountInfo(
		uid,
	)
	if err != nil {
		return nil, err
	}
	return resultb, nil
}

func storeExtractAccountUpdate(ctx context.Context, r *protocol.StoreExtractAccountAddRequest) (*protocol.StoreOperateResponse, error) {
	request := &protocol.StoreExtractAccountAddRequest{}
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	request.Uid = Uid
	if (r.Type != 1 && r.Type != 2) || r.Detail == "" {
		return nil, errutil.ErrInvalidParameter
	}
	request.Type = r.Type
	request.Detail = r.Detail
	request.Uid = Uid
	request.Name = r.Name
	err := StoreExtractAccountUpdate(request)
	if err != nil {
		return &protocol.StoreOperateResponse{Code: -1, Error: err.Error(), Data: false}, nil
	}
	return &protocol.StoreOperateResponse{Data: true}, nil
}

func StoreExtractAccountUpdate(r *protocol.StoreExtractAccountAddRequest) error {
	list, err := StoreExtractAccount(r.Uid)
	if err != nil {
		return nil
	}
	if len(list) == 0 {
		// 插入购买记录表
		err = db.ExtractAccountAdd(
			r.Uid,
			r.Detail,
			r.Type,
			r.Name,
		)
		if err != nil {
			return err
		}
	} else {
		// 插入购买记录表
		err = db.ExtractAccountUpdate(
			r.Uid,
			r.Detail,
			r.Type,
			r.Name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// 入redis
// bounty
func bounty(ctx context.Context) (*protocol.StoreOperateResponse, error) {
	bounty := strconv.Itoa(int(global.ServerIDRangeStart)) + "bounty:"
	// 获取昨日总金额
	// allMoney, err := db.Storetotal(0)
	// if err != nil {
	// 	return nil, err
	// }
	BountyCfg := global.GlobalCfg.Bounty.Grade
	// result := make(map[int32][]protocol.SourceInfo)
	result := &protocol.StoreOperateResponse{}
	// 特权卡
	allWeight := 0
	for _, privilegen := range BountyCfg {
		allWeight = allWeight + int(privilegen.Weight)
	}
	endMoney := 10000
	for _, privilegen := range BountyCfg {
		store, err := global.GetStoreCfg(privilegen.ID)
		if err != nil {
			return nil, err
		}
		itemInfo, err := global.GetItemUseCfg(store.Props[0].Id)
		if err != nil {
			return nil, err
		}
		price := itemInfo.GiveGoods[0].Num / 100
		key := bounty + strconv.Itoa(int(privilegen.ID))
		num := int(float32(endMoney) * float32(privilegen.Weight) / float32(allWeight) / float32(price))
		for i := 0; i < num; i++ {
			cmd := "LPUSH"
			_, err := provide.CmdCache(cmd, key, strconv.Itoa(int(price)))
			if err != nil {
				return nil, err
			}
		}

	}
	return result, nil
}
