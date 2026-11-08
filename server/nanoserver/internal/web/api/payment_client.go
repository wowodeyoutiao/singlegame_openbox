package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nanoserver/dispatcher"
	constant "nanoserver/pkg/constant"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/signature"
	"nanoserver/protocol"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type PaymentClient struct{}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{}
}

// CheckContentSecurity 执行内容安全检查
func (p *PaymentClient) CheckContentSecurity(openID string, content string, scene int) error {
	if openID == "" || content == "" {
		logger.Warnf("CheckContentSecurity: openID or content is empty")
		return nil
	}

	// 如果是生产环境或提审环境，才进行安全检查
	appEnv := strings.ToLower(viper.GetString("core.env"))
	if appEnv != constant.EnvProd && appEnv != constant.EnvReview {
		return nil
	}

	msgSecCheckAddr := fmt.Sprintf("%s/v1/wx/msg-sec-check", constant.PaymentURL)
	httpClient := resty.New()

	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)

	msg := &protocol.MsgSecCheckRequest{
		Openid:  openID,
		Scene:   int32(scene),
		Version: 2, // 版本号，2表示最新版本
		Content: content,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("安全检查请求体序列化失败: %v", err)
	}

	signResult := signer.Sign("POST", "/v1/wx/msg-sec-check", nil, []byte(jsonData))
	msgSecCheckResp, err := httpClient.R().
		SetContext(context.TODO()).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		Post(msgSecCheckAddr)
	if err != nil {
		return fmt.Errorf("安全检查请求失败: %v", err)
	}

	if msgSecCheckResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("安全检查失败, 状态码: %d", msgSecCheckResp.StatusCode())
	}

	msgSecCheck := &protocol.MsgSecCheckResponse{}
	if err := json.Unmarshal(msgSecCheckResp.Body(), msgSecCheck); err != nil {
		return fmt.Errorf("安全检查响应解析失败: %v", err)
	}

	if msgSecCheck.Result.Suggest != "pass" {
		return fmt.Errorf("内容含有敏感信息")
	}

	return nil
}

// MidasPaymentGameItem 虚拟支付
func (p *PaymentClient) MidasPaymentGameItem(ctx context.Context, req *protocol.MidasPaymentGameItemRequest) (*protocol.MidasPaymentGameItemResponse, error) {
	midasPaymentsGameItemAddr := fmt.Sprintf("%s/v1/wx/midas-payments/game-item", constant.PaymentURL)
	httpClient := resty.New()

	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "虚拟支付请求体序列化失败: %v", err)
	}
	midasPaymentsGameItemResp := &protocol.MidasPaymentGameItemResponse{}

	signResult := signer.Sign("POST", "/v1/wx/midas-payments/game-item", nil, []byte(jsonData))
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		SetResult(midasPaymentsGameItemResp).
		Post(midasPaymentsGameItemAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "虚拟支付请求失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("虚拟支付请求失败, 状态码: %d, 错误信息: %s", resp.StatusCode(), string(resp.Body()))
	}

	return midasPaymentsGameItemResp, nil
}

// JSAPIPrepay 微信JSAPI预支付订单
func (p *PaymentClient) JSAPIPrepay(ctx context.Context, req *protocol.PrepayRequest) (*protocol.PrepayResponse, error) {
	JSAPIPrepaAddr := fmt.Sprintf("%s/v1/wx/merchant/prepay", constant.PaymentURL)
	httpClient := resty.New()

	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "JSAPI预支付请求体序列化失败: %v", err)
	}
	prepayResp := &protocol.PrepayResponse{}
	signResult := signer.Sign("POST", "/v1/wx/merchant/prepay", nil, []byte(jsonData))
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		SetResult(prepayResp).
		Post(JSAPIPrepaAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "JSAPI预支付请求失败: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("JSAPI预支付请求失败, 状态码: %d, 错误信息: %s", resp.StatusCode(), string(resp.Body()))
	}

	return prepayResp, nil
}

func (p *PaymentClient) VerifyActiveCode(ctx context.Context, req *protocol.VerifyActiveCodeRequest) (*protocol.VerifyActiveCodeResponse, error) {
	verifyActiveCodeAddr := fmt.Sprintf("%s/v1/active-codes/verify", constant.PaymentURL)
	httpClient := resty.New()
	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)
	jsonData, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("VerifyActiveCode: 激活码验证请求体序列化失败: %v", err)
		return nil, errutil.ErrServerInternal
	}
	verifyActiveCodeResp := &protocol.VerifyActiveCodeResponse{}
	signResult := signer.Sign("POST", "/v1/active-codes/verify", nil, []byte(jsonData))
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		SetResult(verifyActiveCodeResp).
		Post(verifyActiveCodeAddr)
	if err != nil {
		logger.Errorf("VerifyActiveCode: 激活码验证请求失败: %v", err)
		return nil, errutil.ErrServerInternal
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Errorf("VerifyActiveCode: 激活码验证请求失败, 状态码: %d, 错误信息: %s", resp.StatusCode(), string(resp.Body()))
		return nil, errutil.ErrServerInternal
	}

	return verifyActiveCodeResp, nil
}

// CreatePayCharacterRequest 创建游戏角色
func (p *PaymentClient) CreatePayCharacterRequest(ctx context.Context, req *dispatcher.CreateCharacter) (*protocol.PayCommonResponse, error) {
	createCharacterAddr := fmt.Sprintf("%s/v1/characters", constant.PaymentURL)
	httpClient := resty.New()
	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "创建游戏角色请求体序列化失败: %v", err)
	}
	createCharacterResp := &protocol.PayCommonResponse{}
	signResult := signer.Sign("POST", "/v1/characters", nil, []byte(jsonData))
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		SetResult(createCharacterResp).
		Post(createCharacterAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "创建游戏角色请求失败: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("创建游戏角色请求失败, 状态码: %d, 错误信息: %s", resp.StatusCode(), string(resp.Body()))
	}

	return createCharacterResp, nil
}

// UpdatePayCharacterRequest 更新游戏角色
func (p *PaymentClient) UpdatePayCharacterRequest(ctx context.Context, req *dispatcher.CreateCharacter) (*protocol.PayCommonResponse, error) {
	updateCharacterAddr := fmt.Sprintf("%s/v1/characters/%d", constant.PaymentURL, req.CharacterId)
	httpClient := resty.New()
	signerSecret := viper.GetString("signer.secret")
	signer := signature.NewSigner(signerSecret)
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "更新游戏角色请求体序列化失败: %v", err)
	}
	updateCharacterResp := &protocol.PayCommonResponse{}
	signResult := signer.Sign("PUT", fmt.Sprintf("/v1/characters/%d", req.CharacterId), nil, []byte(jsonData))
	resp, err := httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetHeader("X-Signature", signResult.Signature).
		SetHeader("X-Timestamp", signResult.Timestamp).
		SetHeader("X-Nonce", signResult.Nonce).
		SetBody(jsonData).
		SetResult(updateCharacterResp).
		Put(updateCharacterAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "更新游戏角色请求失败: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("更新游戏角色请求失败, 状态码: %d, 错误信息: %s", resp.StatusCode(), string(resp.Body()))
	}

	return updateCharacterResp, nil
}
