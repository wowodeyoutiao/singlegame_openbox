package api

import (
	"context"
	"net/http"

	"nanoserver/internal/game"
	"nanoserver/pkg/my_jwt"
	"nanoserver/protocol/protobuf/pb"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakeSensitiveService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/sensitive/guilds", nex.Handler(GuildCreate)).Methods("POST")      // 帮派创建
	router.Handle("/v1/sensitive/messages", nex.Handler(SendMessage)).Methods("POST")    // 聊天
	router.Handle("/v1/sensitive/role/name", nex.Handler(UpdateRoleName)).Methods("PUT") // 更改角色名字
	return router
}

func GuildCreate(ctx context.Context, req *pb.GuildCreateRequest) (resp *pb.GuildCreateResponse, err error) {
	uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	userInfo, err := game.UserInfo(uid)
	if err != nil {
		logger.Errorf("GuildCreate 获取用户信息失败: %v, uid: %d", err, uid)
		return &pb.GuildCreateResponse{
			Code:  -1,
			Error: "获取用户信息失败",
		}, nil
	}

	if err := NewPaymentClient().CheckContentSecurity(userInfo.AccountName, req.GetName(), 1); err != nil {
		logger.Errorf("GuildCreate 安全检查失败: %v, uid: %d", err, uid)
		return &pb.GuildCreateResponse{
			Code:  -1,
			Error: "内容含有敏感信息",
		}, nil
	}

	resp, err = game.GuildCreateREST(uid, req)
	if err != nil {
		logger.Errorf("GuildCreate 创建公会失败: %v, uid: %d", err, uid)
		return &pb.GuildCreateResponse{
			Code:  -1,
			Error: err.Error(),
		}, nil
	}

	return resp, nil
}

func SendMessage(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	userInfo, err := game.UserInfo(uid)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, uid: %d", err, uid)
		return &pb.ChatResponse{
			Code:  -1,
			Error: "获取用户信息失败",
		}, nil
	}

	if err := NewPaymentClient().CheckContentSecurity(userInfo.AccountName, req.GetContext(), 1); err != nil {
		upMessageShushuDate(uid, int32(req.Type), 1, req.Context, false)
		logger.Errorf("SendMessage 安全检查失败: %v, uid: %d", err, uid)
		return &pb.ChatResponse{
			Code:  -1,
			Error: "内容含有敏感信息",
		}, nil
	}

	resp, err := game.ChatREST(uid, req)
	if err != nil {
		upMessageShushuDate(uid, int32(req.Type), 2, req.Context, false)
		logger.Errorf("发送聊天消息失败: %v, uid: %d", err, uid)
		return &pb.ChatResponse{
			Code:  -1,
			Error: err.Error(),
		}, nil
	}
	upMessageShushuDate(uid, int32(req.Type), 3, req.Context, true)
	return resp, nil
}

func upMessageShushuDate(Uid int64, channel_type int32, message_status int32, message_content string, is_success bool) {
	shushuDate := make(map[string]interface{}, 0)
	//在合并点其他的数据
	shushuDate["channel_type"] = channel_type
	shushuDate["message_content"] = message_content
	shushuDate["message_status"] = message_status
	shushuDate["is_success"] = is_success
	shushuDate["name"] = "message"
	game.Upshushu(Uid, shushuDate, true, false)
}

func UpdateRoleName(ctx context.Context, req *pb.ChangeRoleNameRequest) (*pb.ChangeRoleNameResponse, error) {
	uid := ctx.Value("userDate").(my_jwt.CustomClaims).UserId

	userInfo, err := game.UserInfo(uid)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, uid: %d", err, uid)
		return &pb.ChangeRoleNameResponse{
			Code:     -1,
			Error:    "获取用户信息失败",
			RoleName: req.GetRoleName(),
		}, nil
	}

	if err := NewPaymentClient().CheckContentSecurity(userInfo.AccountName, req.GetRoleName(), 1); err != nil {
		logger.Errorf("UpdateRoleName 安全检查失败: %v, uid: %d", err, uid)
		return &pb.ChangeRoleNameResponse{
			Code:     -1,
			Error:    "内容含有敏感信息",
			RoleName: req.GetRoleName(),
		}, nil
	}

	name, err := game.ChangeRoleNameREST(uid, req.RoleName)
	if err != nil {
		logger.Errorf("修改角色名失败: %v, uid: %d", err, uid)
		return &pb.ChangeRoleNameResponse{
			Code:     -1,
			Error:    err.Error(),
			RoleName: req.GetRoleName(),
		}, nil
	}

	return &pb.ChangeRoleNameResponse{
		Code:     0,
		Error:    "",
		RoleName: name,
	}, nil
}
