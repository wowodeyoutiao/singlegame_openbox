package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"nanoserver/dispatcher"
	"nanoserver/pkg/signature"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request received authMiddleware: %s %s", r.Method, r.URL.Path)

		xSignature := r.Header.Get("X-Signature")
		timestamp := r.Header.Get("X-Timestamp")
		nonce := r.Header.Get("X-Nonce")
		if xSignature == "" || timestamp == "" || nonce == "" {
			log.Errorf("Missing required headers: X-Signature, X-Timestamp, or X-Nonce")
			http.Error(w, "Missing required headers", http.StatusBadRequest)
			return
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		secret := viper.GetString("signer.secret")
		singer := signature.NewSigner(secret)
		ok, err := singer.Verify(
			r.Method,
			r.URL.Path,
			xSignature,
			timestamp,
			nonce,
			r.URL.Query(),
			bodyBytes,
		)
		if err != nil {
			log.Errorf("Signature verification failed: %v", err)
			http.Error(w, "Signature verification failed", http.StatusUnauthorized)
			return
		}
		if !ok {
			log.Errorf("Invalid signature: %s", xSignature)
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	})
}

func MakePayService() http.Handler {
	router := mux.NewRouter()
	payService := &PayService{
		uc: &PayUseCase{},
	}

	router.Handle("/v1/pay/notify", authMiddleware(http.HandlerFunc(payService.PayNotify))).Methods("POST")    // 支付通知的路由
	router.Handle("/v1/pay/userinfo", authMiddleware(http.HandlerFunc(payService.PayUserInfo))).Methods("GET") // 获取用户信息路由

	return router
}

type PayService struct {
	uc *PayUseCase
}

func (p *PayService) PayNotify(w http.ResponseWriter, r *http.Request) {
	var req protocol.PayNotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("PayNotify decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(protocol.PayNotifyResponse{
			Code:    400,
			Message: "Invalid request",
		}); err != nil {
			log.Errorf("PayNotify encode error: %v", err)
		}
		return
	}

	log.Infof("PayNotify request: %+v", req)

	req.Num = 1 // 默认购买数量为1
	resp, err := p.uc.PayNotify(&req)
	if err != nil {
		log.Errorf("PayNotify processing failed, uid=%v, store_id=%v, err: %v", req.Uid, req.StoreId, err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(protocol.PayNotifyResponse{
			Code:    500,
			Message: "Internal server error",
		}); err != nil {
			log.Errorf("PayNotify encode error: %v", err)
		}
		return
	}

	log.Infof("PayNotify response: %+v", resp)

	dispatcher.MustGetInstance().Publish(dispatcher.TopicPayRefreshPush.String(), &pb.PayRefreshPush{RequestId: req.RequestID, UserId: req.Uid, StoreId: int32(req.StoreId)})

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode response: %v", err)
	}
}

func (p *PayService) PayUserInfo(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("uid")
	openIDStr := r.URL.Query().Get("open_id")
	if userIDStr == "" && openIDStr == "" {
		log.Errorf("Missing uid in query params")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(protocol.PayCommonResponse{
			Code:    400,
			Message: "Missing uid or openid in query parameters",
		})
		return
	}

	var (
		userID int64
		err    error
	)
	if userIDStr != "" {
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Errorf("Invalid uid format: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(protocol.PayCommonResponse{
				Code:    400,
				Message: "Invalid uid format",
			})
			return
		}
	}

	req := protocol.PayUserInfoRequest{
		Uid:    userID,
		OpenID: openIDStr,
	}

	log.Infof("PayUserInfo request: %+v", req)
	resp, err := p.uc.PayUserInfo(&req)
	if err != nil {
		log.Errorf("Failed to get user info: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(protocol.PayCommonResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}

	log.Infof("PayUserInfo response: %+v", resp)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode response: %v", err)
	}
}
