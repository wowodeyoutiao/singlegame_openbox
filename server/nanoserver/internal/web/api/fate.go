package api

import (
	"context"
	"nanoserver/global"
	"nanoserver/internal/game"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"

	//"nanoserver/internal/web"
	"nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol"
	"net/http"
)

func MakeFateService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/fate/detail", nex.Handler(fateList)).Methods("POST")                               //列表
	router.Handle("/v1/fate/up", nex.Handler(fateUp)).Methods("Post")                                     //升级
	router.Handle("/v1/fate/break", limitHandler(nex.Handler(fateBreak), "fate_extract")).Methods("Post") //进阶

	return router
}

func FateList(r *protocol.FateOperateRequest) (int32, int32, error) {
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return 0, 0, err
	}

	return int32(res.FateLevel), int32(res.FateExp), nil
}

func fateList(ctx context.Context, r *http.Request) (*protocol.FateListResponse, error) {

	request := &protocol.FateOperateRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	fateLevel, exp, err := FateList(request)
	if err != nil {
		return &protocol.FateListResponse{FateLevel: 0, Exp: 0}, nil
	}
	return &protocol.FateListResponse{FateLevel: fateLevel, Exp: exp}, nil
}

func fateBreak(ctx context.Context, r *http.Request) (*protocol.FateOperateResponse, error) {
	request := &protocol.FateOperateRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	err = FatenBreak(request)
	if err != nil {
		return &protocol.FateOperateResponse{Data: err.Error(), Code: -1}, nil
	}
	return &protocol.FateOperateResponse{Data: "成功", Code: 0}, err
}
func FatenBreak(r *protocol.FateOperateRequest) error {
	res, err := game.UserInfo(r.Uid)
	if err != nil {
		return err
	}
	fatecfg, err := global.GetFateCfg(res.FateLevel)
	if err != nil {
		return errutil.ErrConfigNotFound
	}
	maxRes, err := global.GetFateMaxCfg()
	if err != nil {
		return errutil.ErrConfigNotFound
	}

	if maxRes.FateLevel == res.FateLevel {
		return errutil.ErrMaxLevel
	}
	//是否需要晋升
	if len(fatecfg.AdvancedNeed) == 0 {
		return errutil.ErrLevel
	}
	if res.FateLevel != fatecfg.FateLevel {
		return errutil.ErrLevel
	}
	FateCfg := global.FateCfg
	//查询这一重需要多少经验
	AllExp := fatecfg.FateExp
	for j := 1; j <= len(FateCfg)+1; j++ {
		fatecfg, err := global.GetFateCfg(res.FateLevel - int32(j))
		if err != nil || len(fatecfg.AdvancedNeed) > 0 {
			break
		}
		AllExp = AllExp + fatecfg.FateExp
	}
	if int32(res.FateExp) < AllExp {
		return errutil.ErrLevel
	}
	resultc := make([]protocol.SourceInfo, len(fatecfg.AdvancedNeed))
	for i, info := range fatecfg.AdvancedNeed {
		resultc[i].Id = info.ID
		resultc[i].Num = -int64(info.Num)
	}
	err = CheckSourceNum(resultc, r.Uid, true, constant.ItemFateBreak)
	if err != nil {
		return err
	}
	err = game.SetPlayInfo(r.Uid, res.FateLevel+1, 0)
	if err != nil {
		return err
	}
	return nil
}
func fateUp(ctx context.Context, r *protocol.FateUpRequest) (*protocol.FateUpResponse, error) {
	request := &protocol.FateUpRequest{}
	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid
	request.ConfigId = r.ConfigId
	if r.Num == 0 {
		r.Num = 1
	}
	request.Num = r.Num
	var whitelist = map[int32]bool{75: true, 72: true, 73: true, 74: true}

	if !whitelist[request.ConfigId] {
		return &protocol.FateUpResponse{}, errutil.ErrInvalidParameter
	}
	perAll, num, err := FatenUp(request)
	if err != nil {
		return &protocol.FateUpResponse{Data: err.Error(), Code: -1}, nil
	}
	return &protocol.FateUpResponse{Data: "成功", Code: 0, Per: perAll, Num: num}, nil
}
func FatenUp(r *protocol.FateUpRequest) ([]int, int32, error) {
	//减去消耗品
	result := make([]protocol.SourceInfo, 1)
	result[0].Id = r.ConfigId
	result[0].Num = -int64(r.Num)
	err := CheckSourceNum(result, r.Uid, true, constant.ItemFateUp)
	if err != nil {
		return nil, r.Num, err
	}
	res, num, err := game.UseFate(r.Uid, r.ConfigId, r.Num)
	if err != nil {
		return nil, num, err
	}
	return res, num, nil
}
