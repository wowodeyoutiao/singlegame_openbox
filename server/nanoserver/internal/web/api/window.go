package api

import (
	"context"
	"fmt"
	"nanoserver/db"
	"nanoserver/internal/game"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/pkg/algoutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/timex"
	"nanoserver/protocol"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	"github.com/pkg/errors"
)

// API密钥，用于外部接口的加密验证
const apiSecretKey = "zdyjb202507"

func MakeWindowService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/window/register_remain", nex.Handler(registerRemain)).Methods("GET") //列表
	router.Handle("/v1/window/payment_remain", nex.Handler(paymentRemain)).Methods("GET")   //列表
	router.Handle("/v1/window/pay_per", nex.Handler(payPer)).Methods("GET")                 //列表
	router.Handle("/v1/window/ltv", nex.Handler(ltv)).Methods("GET")                        //列表
	router.Handle("/v1/window/online_log", nex.Handler(online)).Methods("GET")              //列表
	router.Handle("/v1/window/ingot", nex.Handler(ingot)).Methods("GET")                    //列表
	router.Handle("/v1/window/have_money", nex.Handler(haveMoney)).Methods("GET")           //列表
	router.Handle("/v1/window/reward", nex.Handler(reward)).Methods("GET")                  //列表
	router.Handle("/v1/window/gold", nex.Handler(gold)).Methods("GET")                      //列表
	router.Handle("/v1/window/soul", nex.Handler(soul)).Methods("GET")                      //列表
	router.Handle("/v1/window/shushu", nex.Handler(shushu)).Methods("POST")
	router.Handle("/v1/window/shushuregister", nex.Handler(shushuRegister)).Methods("POST")
	router.Handle("/v1/window/send_lamp", nex.Handler(sendLamp)).Methods("POST")
	router.Handle("/v1/window/send_mail", nex.Handler(sendMail)).Methods("POST")
	router.Handle("/v1/window/send_mute", nex.Handler(sendMute)).Methods("POST")
	router.Handle("/v1/window/ban_user", nex.Handler(banUser)).Methods("POST")
	router.Handle("/v1/window/delete_user", nex.Handler(deleteUser)).Methods("POST")
	router.Handle("/v1/window/query_user_info", nex.Handler(queryUserInfo)).Methods("POST")
	//列表
	router.Handle("/v1/window/shushuno", nex.Handler(shushuno)).Methods("POST")      //列表
	router.Handle("/v1/window/tenxun", nex.Handler(tenxun)).Methods("GET")           //列表
	router.Handle("/v1/window/ta_mp_hide", nex.Handler(mpHide)).Methods("POST")      //列表
	router.Handle("/v1/window/completed", nex.Handler(completed)).Methods("POST")    //列表
	router.Handle("/v1/window/openBox", nex.Handler(openbox)).Methods("POST")        //列表
	router.Handle("/v1/window/ad_show", nex.Handler(adshow)).Methods("POST")         //列表
	router.Handle("/v1/window/ad_close", nex.Handler(adclose)).Methods("POST")       //列表
	router.Handle("/v1/window/click_event", nex.Handler(clickEvent)).Methods("POST") //列表

	return router
}
func registerRemain(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  1,
			Two:    2,
			Three:  3,
			Seven:  4,
			Thirty: 5,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  6,
			Two:    7,
			Three:  8,
			Seven:  9,
			Thirty: 10,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    12,
			Three:  13,
			Seven:  14,
			Thirty: 15,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 20,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func paymentRemain(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  13,
			Two:    24,
			Three:  35,
			Seven:  46,
			Thirty: 57,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  68,
			Two:    79,
			Three:  810,
			Seven:  911,
			Thirty: 1012,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  113,
			Two:    214,
			Three:  315,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func payPer(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111,
			Two:    222,
			Three:  333,
			Seven:  444,
			Thirty: 555,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  666,
			Two:    777,
			Three:  888,
			Seven:  999,
			Thirty: 101010,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  1111,
			Two:    2222,
			Three:  3333,
			Seven:  4444,
			Thirty: 5555,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66666,
			Two:    77777,
			Three:  88888,
			Seven:  99999,
			Thirty: 10101010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func ltv(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func online(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func ingot(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}

func haveMoney(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func reward(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}

func gold(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}
func soul(r *protocol.WindowLtvListRequest) (*protocol.WindowLtvResponse, error) {
	count := 2
	RankListDetail := protocol.WindowLtvList{
		Count:  count,
		Detail: make([]protocol.WindowLtvDetail, count),
	}
	if !r.Export {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  111111,
			Two:    1111112,
			Three:  1111113,
			Seven:  1111114,
			Thirty: 1111115,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  16,
			Two:    17,
			Three:  18,
			Seven:  19,
			Thirty: 110,
		}
	} else {
		RankListDetail.Detail[0] = protocol.WindowLtvDetail{
			Date:   "2024-09-01",
			Today:  11,
			Two:    22,
			Three:  33,
			Seven:  44,
			Thirty: 55,
		}
		RankListDetail.Detail[1] = protocol.WindowLtvDetail{
			Date:   "2024-09-02",
			Today:  66,
			Two:    77,
			Three:  88,
			Seven:  99,
			Thirty: 1010,
		}
	}

	return &protocol.WindowLtvResponse{Data: RankListDetail}, nil
}

func shushu(ctx context.Context, r *protocol.WindowShushuRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = r.Name
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func shushuRegister(ctx context.Context, r *protocol.WindowShushuRegisterRequest) (*protocol.WindowShushuResponse, error) {
	Uid := int64(0)
	data := make(map[string]interface{}, 0)
	data["name"] = "register"
	data["create_time"] = time.Now().Unix()
	data["tmp_role_id"] = int64(1)
	data["tmp_account"] = r.Name
	data["register_type"] = "wechat"
	game.Upshushu(Uid, data, false, false)

	userSet := make(map[string]interface{}, 0)
	userSet["ip"] = r.Ip
	userSet["tmp_role_id"] = int64(1)
	userSet["tmp_account"] = r.Name
	userSet["register_time"] = r.Register
	userSet["register_time_v2"] = time.Unix(r.Register, 0).Format(timex.FormatDateTimeMillis)
	game.UpshushuSet(Uid, userSet, false)

	userSetOnce := make(map[string]interface{}, 0)
	userSetOnce["tmp_role_id"] = r.Name
	userSetOnce["tmp_account"] = r.Name
	userSetOnce["campaign_id"] = r.Channel
	game.UpshushuSetOnce(Uid, userSetOnce, false)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func shushuno(ctx context.Context, r *protocol.WindowShushuRequest) (*protocol.WindowShushuResponse, error) {
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = r.Name
	upPkgShushuNo(PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func tenxun(ctx context.Context, r *http.Request) (*protocol.WindowShushuResponse, error) {
	//获取url上携带的参数
	query := r.URL.Query()
	adgroup_id := query.Get("adgroup_id")
	click_id := query.Get("click_id")
	count_id := query.Get("count_id")
	call_back := query.Get("call_back")
	data := make(map[string]interface{}, 0)
	data["name"] = "register"
	data["tmp_role_id"] = int64(2576266495554232320)
	data["tmp_account"] = "2576266495554232320"
	data["count_id"] = count_id
	data["call_back"] = call_back
	data["adgroup_id"] = adgroup_id
	data["click_id"] = click_id
	game.Upshushu(int64(2576266495554232320), data, false, false)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return nil, nil
}
func mpHide(ctx context.Context, r *protocol.WindowShushuHumpHideRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = r.Name
	PostDate["#duration"] = r.Duration
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func completed(ctx context.Context, r *protocol.WindowShushuCompletedRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = r.Name
	PostDate["guide_step"] = r.Guide_step
	PostDate["guide_name"] = r.Guide_name
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func openbox(ctx context.Context, r *protocol.WindowShusOpenBoxRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = r.Name
	PostDate["box_item_info"] = r.BoxItemInfo
	// 处理 box_item_info 为数组
	boxItemInfos := make([]map[string]interface{}, 0, len(r.BoxItemInfo))
	for _, item := range r.BoxItemInfo {
		boxItemInfos = append(boxItemInfos, map[string]interface{}{
			"item_id":       item.ItemID,
			"item_level":    item.ItemLevel,
			"item_quality":  item.ItemQuality,
			"item_name":     item.ItemName,
			"count":         item.Count,
			"item_property": item.ItemProperty,
		})
	}
	PostDate["box_item_info"] = boxItemInfos
	PostDate["is_auto"] = r.IsAuto
	PostDate["power_up"] = r.PowerUp
	PostDate["open_num"] = r.OpenNum
	PostDate["role_power_stop"] = r.RolePowerStop
	PostDate["partner_power_stop"] = r.PartnerPowerStop
	PostDate["box_num_before"] = r.BoxNumBefore
	PostDate["box_num_after"] = r.BoxNumAfter
	PostDate["box_job_limited"] = r.BoxJobLimited
	PostDate["box_level_limited"] = r.BoxLevelLimited
	PostDate["retain_property"] = r.RetainProperty
	PostDate["retain_property_num"] = r.RetainPropertyNum

	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func adclose(ctx context.Context, r *protocol.WindowShushuAdCloseRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["ad_id"] = r.AdID
	PostDate["ad_type"] = r.AdType
	PostDate["ad_scene"] = r.AdScene
	PostDate["ad_time"] = r.AdTime
	PostDate["ad_duration"] = r.AdDuration
	PostDate["view_complete"] = r.ViewComplete
	PostDate["reward_type"] = r.RewardType
	PostDate["ad_position"] = r.AdPosition
	PostDate["name"] = "ad_close"
	boxItemInfos := make([]map[string]interface{}, 0, len(r.GetItemInfo))
	for _, item := range r.GetItemInfo {
		boxItemInfos = append(boxItemInfos, map[string]interface{}{
			"item_id":   item.ItemID,
			"item_name": item.ItemName,
			"item_num":  item.ItemNum,
		})
	}
	PostDate["get_item_info"] = boxItemInfos
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func adshow(ctx context.Context, r *protocol.WindowShushuAdCloseRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = "ad_show"
	PostDate["ad_id"] = r.AdID
	PostDate["ad_type"] = r.AdType
	PostDate["ad_scene"] = r.AdScene
	PostDate["ad_position"] = r.AdPosition
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func clickEvent(ctx context.Context, r *protocol.WindowShushuClickEventRequest) (*protocol.WindowShushuResponse, error) {
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId
	PostDate := make(map[string]interface{}, 0)
	PostDate["name"] = "click_event"
	PostDate["page"] = r.Page
	PostDate["component"] = r.Component
	PostDate["position"] = r.Position
	PostDate["event_click"] = r.EventClick
	upPkgShushu(Uid, PostDate)
	// 这里可以根据需要处理 params
	// 例如：fmt.Println(params)
	return &protocol.WindowShushuResponse{Code: 0}, nil
}
func upPkgShushu(uid int64, postDate map[string]interface{}) {
	game.Upshushu(uid, postDate, true, false)
}
func upPkgShushuNo(postDate map[string]interface{}) {
	game.Upshushuno(postDate)
}

// verifyKey 验证加密key
// secret: 加密密钥，与 sendMail 保持一致
func verifyKey(value int64, key string, secret string) error {
	delKey := fmt.Sprintf("%d%s", value, secret)
	expectedKey := algoutil.MD5String(delKey)
	if expectedKey != key {
		logger.Error("verifyKey error", "expectedKey:", expectedKey, ",get key:", key)
		return errors.New("key error")
	}
	return nil
}
func sendLamp(ctx context.Context, r *protocol.SendLampRequest) (*protocol.WindowShushuResponse, error) {
	game.SystemLamp(int32(r.Stype), r.Svalue, r.Title)
	return nil, nil
}
func sendMute(ctx context.Context, r *protocol.SendMuteRequest) (*protocol.WindowShushuResponse, error) {
	// 验证加密key
	if err := verifyKey(r.Uid, r.Key, apiSecretKey); err != nil {
		return nil, err
	}

	// 通过 GM 指令执行禁言操作
	cmd := protocol.GmCmd{
		Cmd:   "mute",
		Param: []string{fmt.Sprintf("%d", r.Uid), fmt.Sprintf("%d", r.Duration)},
	}
	if err := game.GMCommand(cmd); err != nil {
		return nil, err
	}

	return &protocol.WindowShushuResponse{Code: 0}, nil
}

func banUser(ctx context.Context, r *protocol.BanUserRequest) (*protocol.WindowShushuResponse, error) {
	// 验证加密key
	if err := verifyKey(r.Uid, r.Key, apiSecretKey); err != nil {
		return nil, err
	}

	// 通过 GM 指令执行封号操作
	cmd := protocol.GmCmd{
		Cmd:   "ban",
		Param: []string{fmt.Sprintf("%d", r.Uid), fmt.Sprintf("%d", r.Duration)},
	}
	if err := game.GMCommand(cmd); err != nil {
		return nil, err
	}

	return &protocol.WindowShushuResponse{Code: 0}, nil
}

func deleteUser(ctx context.Context, r *protocol.DeleteUserRequest) (*protocol.WindowShushuResponse, error) {
	// 验证加密key
	if err := verifyKey(r.Uid, r.Key, apiSecretKey); err != nil {
		return nil, err
	}

	// 验证 action 参数
	if r.Action != 1 && r.Action != 2 {
		return nil, fmt.Errorf("action参数错误，只能是1（删除）或2（恢复）")
	}

	// 通过 GM 指令执行删除/恢复操作（会自动先踢下线）
	cmd := protocol.GmCmd{
		Cmd:   "delete",
		Param: []string{fmt.Sprintf("%d", r.Uid), fmt.Sprintf("%d", r.Action)},
	}
	if err := game.GMCommand(cmd); err != nil {
		return nil, err
	}

	return &protocol.WindowShushuResponse{Code: 0}, nil
}

func queryUserInfo(ctx context.Context, r *protocol.QueryUserInfoRequest) (*protocol.QueryUserInfoResponse, error) {
	// 验证加密key
	if err := verifyKey(r.Uid, r.Key, apiSecretKey); err != nil {
		return nil, err
	}

	// 先尝试从内存中获取玩家信息
	memoryInfo, memoryErr := game.UserInfo(r.Uid)

	// 再从数据库查询（用于补充或兜底）
	user, dbErr := db.QueryUser(r.Uid)
	if dbErr != nil && memoryInfo == nil {
		return nil, fmt.Errorf("query user info failed: %w", dbErr)
	}

	var userInfo *protocol.UserInfo
	if memoryErr == nil && memoryInfo != nil {
		// 内存数据优先
		userInfo = memoryInfo
		userInfo.IsOnline = 1
	}

	if userInfo == nil && user != nil {
		// 仅有数据库数据的情况，保持原有行为
		userInfo = &protocol.UserInfo{
			UID:         user.Uid,
			Name:        user.Name,
			Level:       user.Level,
			ServerId:    user.ServerId,
			AccountName: user.AccountName,
			RegisterAt:  user.RegisterAt,
			LastLoginAt: user.LastLoginAt,
			Status:      user.Status,
			IsOnline:    user.IsOnline,
			Role:        user.Role,
			MuteUntil:   user.MuteUntil,
		}
	} else if userInfo != nil && user != nil {
		// 内存数据为主，缺失字段再补齐
		if userInfo.Name == "" && user.Name != "" {
			userInfo.Name = user.Name
		}
		if userInfo.Level == 0 && user.Level != 0 {
			userInfo.Level = user.Level
		}
		if userInfo.ServerId == 0 && user.ServerId != 0 {
			userInfo.ServerId = user.ServerId
		}
		if userInfo.AccountName == "" && user.AccountName != "" {
			userInfo.AccountName = user.AccountName
		}
		if userInfo.RegisterAt == 0 && user.RegisterAt != 0 {
			userInfo.RegisterAt = user.RegisterAt
		}
		if userInfo.LastLoginAt == 0 && user.LastLoginAt != 0 {
			userInfo.LastLoginAt = user.LastLoginAt
		}
		if userInfo.Status == 0 && user.Status != 0 {
			userInfo.Status = user.Status
		}
		if userInfo.Role == 0 && user.Role != 0 {
			userInfo.Role = user.Role
		}
		if userInfo.MuteUntil == 0 && user.MuteUntil != 0 {
			userInfo.MuteUntil = user.MuteUntil
		}
		userInfo.IsOnline = user.IsOnline
	}

	// 查询注册信息（获取封号时间）
	// 封号是针对账号的，需要通过 AccountName（IMEI）查询 Register 表
	banUntil := int64(0)
	accountName := ""
	serverID := int32(0)
	if userInfo != nil {
		accountName = userInfo.AccountName
		serverID = userInfo.ServerId
	}
	if accountName == "" && user != nil {
		accountName = user.AccountName
		serverID = user.ServerId
	}
	if accountName != "" {
		if register, err := db.QueryRegister("", accountName, serverID); err == nil && register != nil {
			banUntil = register.BanUntil
		}
	}
	if userInfo != nil {
		userInfo.BanUntil = banUntil
	}

	return &protocol.QueryUserInfoResponse{
		Code: 0,
		Data: userInfo,
	}, nil
}

func sendMail(ctx context.Context, r *protocol.SendMailRequest) (*protocol.WindowShushuResponse, error) {
	if err := verifyKey(r.SValue, r.Key, apiSecretKey); err != nil {
		return nil, err
	}
	result := make([]protocol.Attachment, len(r.Attachment))
	// MailId := strconv.FormatInt(snowflake.GenID(), 10)
	for i, Attachment := range r.Attachment {
		result[i] = protocol.Attachment{
			MailId:      "",
			UserPropsId: int64(Attachment.UserPropsId),
			Num:         int64(Attachment.Num),
		}
	}
	Peoples := []int64{}
	switch r.Stype {
	case 1: // 全部的人
		//查询这个服务器的下的所有uid
		// shushumain.ReportMailData(protocol.MailAdd{})
	case 2: // 区服
	case 3: // 服务器邮件
		//查询这个服务器的下的所有uid
		// shushumain.ReportMailData(protocol.MailAdd{})
		res, err := db.QueryUserIDsByServerID(r.SValue)
		if err != nil {
			logger.Error("sendMail db.QueryUserIDsByServerID error", "err", err)
			return nil, err
		}
		Peoples = res
	case 4: // 指定用户邮件
		Peoples = []int64{r.SValue}
	}
	for _, people := range Peoples {
		shushumain.ReportMailData(protocol.MailAdd{
			Uid:          -1,
			RecipientUid: people,
			Title:        r.Title,
			Desc:         r.Desc,
			Attachments:  result,
			ChannelCode:  "system-mail",
		})
		Notification(people)
	}
	return nil, nil
}
