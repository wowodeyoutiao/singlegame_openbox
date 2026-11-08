package api

import (
	"context"
	"nanoserver/internal/game"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"

	//"nanoserver/internal/web"

	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol"
	"net/http"
)

func MakeRemainService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/remain/detail", nex.Handler(remainList)).Methods("POST") //列表

	return router
}

func RemainList(r *protocol.RemainOperateRequest) (int32, int32, error) {
	_, err := game.UserInfo(r.Uid)
	if err != nil {
		return 0, 0, err
	}

	return 0, 0, nil
}

func remainList(ctx context.Context, r *http.Request) (*protocol.RemainListResponse, error) {

	request := &protocol.RemainOperateRequest{}
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	Uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	request.Uid = Uid

	remainLevel, exp, err := RemainList(request)
	if err != nil {
		return &protocol.RemainListResponse{RemainLevel: 0, Exp: 0}, nil
	}
	return &protocol.RemainListResponse{RemainLevel: remainLevel, Exp: exp}, nil
}
