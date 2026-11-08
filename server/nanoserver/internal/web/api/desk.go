package api

import (
	"net/http"
	"strconv"
	"time"

	"nanoserver/db"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/whitelist"
	"nanoserver/protocol"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakeDeskService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/desk/player/{id}", nex.Handler(deskList)).Methods("POST") //获取desk列表(lite)
	router.Handle("/v1/desk/{id}", nex.Handler(deskByID)).Methods("POST")        //获取desk记录
	return router
}

func DeskByID(id int64) (*protocol.Desk, error) {
	p, err := db.QueryDesk(id)
	if err != nil {
		return nil, err
	}
	return &protocol.Desk{
		Id:      p.Id,
		Creator: p.Creator,
		Mode:    p.Mode,
		Round:   p.Round,
		DeskNo:  p.DeskNo,

		CreatedAt: p.CreatedAt,
		DismissAt: p.Creator,
	}, nil

}

func DeskList(playerId int64) ([]protocol.Desk, int64, error) {
	//默认全部
	ps, total, err := db.DeskList(playerId)
	if err != nil {
		return nil, 0, err
	}
	list := make([]protocol.Desk, total)

	const (
		format = "2006-01-02 15:04:05"
	)

	for i, p := range ps {

		createdAtStr := time.Unix(p.CreatedAt, 0).Format(format)

		list[i] = protocol.Desk{
			Id:           p.Id,
			Creator:      p.Creator,
			Round:        p.Round,
			Mode:         p.Mode,
			DeskNo:       p.DeskNo,
			CreatedAt:    p.CreatedAt,
			CreatedAtStr: createdAtStr,
			DismissAt:    p.Creator,
		}
	}
	return list, int64(len(list)), nil
}

func deskList(r *http.Request) (*protocol.DeskListResponse, error) {
	if !whitelist.VerifyIP(r.RemoteAddr) {
		return nil, errutil.ErrPermissionDenied
	}
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		return nil, errutil.ErrInvalidParameter
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return nil, errutil.ErrInvalidParameter
	}

	list, t, err := DeskList(id)
	if err != nil {
		return nil, err
	}
	return &protocol.DeskListResponse{Data: list, Total: t}, nil
}

func deskByID(r *http.Request) (*protocol.DeskByIDResponse, error) {
	if !whitelist.VerifyIP(r.RemoteAddr) {
		return nil, errutil.ErrPermissionDenied
	}
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		return nil, errutil.ErrInvalidParameter
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return nil, errutil.ErrInvalidParameter
	}

	h, err := DeskByID(id)
	if err != nil {
		return nil, err
	}
	return &protocol.DeskByIDResponse{Data: h}, nil
}
