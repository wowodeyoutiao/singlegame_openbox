package wechatv2

// import (
// 	//nolint

// 	"bytes"
// 	"context"
// 	"crypto"
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"crypto/sha1"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"sort"
// 	"strings"
// 	"time"

// 	"nanoserver/protocol"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/wechatpay-apiv3/wechatpay-go/core"
// 	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
// 	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
// 	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
// 	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
// 	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
// 	"github.com/wechatpay-apiv3/wechatpay-go/utils"
// )

// type WeChatStore interface {
// 	SetAccessToken(key, accessToken string, expired time.Duration) error
// 	GetAccessToken(key string) (string, error)
// 	SetSessionKey(key, sessionKey string, expired time.Duration) error
// 	GetSessionKey(key string) (string, error)
// }

// type WeChat struct {
// 	appID      string        // 小程序的唯一标识 ID，用于在微信生态中识别一个具体的小程序应用。
// 	offerID    int           // 微信支付时使用的应用ID（旧称），现通常用 appId 表示，用于标识接入微信支付的应用。
// 	appKey     string        // 有时与 appSecret 混淆，实际开发中通常指代小程序或服务接口的密钥（如通信加密用的 key）。
// 	appSecret  string        // 小程序的服务端密钥，用于服务器与微信接口通信的身份验证和获取 token 等敏感操作。
// 	store      WeChatStore   // 存储接口，用于存储和获取 access_token 等信息
// 	httpClient *resty.Client // resty 客户端，用于发送 HTTP 请求

// 	// 商户号相关的配置
// 	mchID string // 商户号 ID，用于标识微信支付的商户账户。  1719427832
// 	// mchCertificateSerialNumber string       // 商户证书序列号，用于微信支付的证书验证。  6E0C644F18A039A5424B11CDF5B5AFB0E124AC99
// 	apiV3Key string // APIv3 密钥，用于微信支付的安全通信。   692Xl4OPJhFyuxy4ApEV6z3pwIafJ10T
// 	// mchPrivateKeyPath          string       // 商户私钥路径，用于签名和验证支付请求。
// 	notifyURL                 string       // 微信支付通知地址，用于接收微信支付结果通知。
// 	client                    *core.Client // 微信支付客户端，用于与微信支付 API 交互。
// 	parseNotifyRequestHandler *notify.Handler
// 	privateKey                *rsa.PrivateKey
// }

// type Config struct {
// 	OfferID                    int
// 	AppID                      string // 小程序的唯一标识 ID，用于在微信生态中识别一个具体的小程序应用。
// 	AppKey                     string // 小程序或服务接口的密钥（如通信加密用的 key），有时与 appSecret 混淆。
// 	AppSecret                  string // 小程序的服务端密钥，用于服务器与微信接口通信的身份验证和获取 token 等敏感操作。
// 	MchID                      string // 商户号 ID，用于标识微信支付的商户账户。
// 	MchCertificateSerialNumber string // 商户证书序列号，用于微信支付的证书验证。
// 	MchPrivateKeyPath          string // 商户私钥路径，用于签名和验证支付请求。
// 	ApiV3Key                   string // APIv3 密钥，用于微信支付的安全通信。
// 	PublicKeyID                string // 公钥ID https://pay.weixin.qq.com/index.php/core/cert/api_cert#/public-key-detail
// 	PublicKey                  string // 公钥内容
// 	NotifyURL                  string // 微信支付通知地址
// }

// // NOTE: 微信支付相关，如何初始化失败则程序不能正常启动
// func New(config *Config, store WeChatStore) *WeChat {
// 	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(config.MchPrivateKeyPath)
// 	if err != nil {
// 		log.Fatalf("failed to load merchant private key: %s", err)
// 	}

// 	mcPublicKey, err := utils.LoadPublicKey(config.PublicKey)
// 	if err != nil {
// 		log.Fatalf("failed to load public key: %s", err)
// 	}
// 	ctx := context.Background()
// 	opts := []core.ClientOption{
// 		// option.WithWechatPayAutoAuthCipher(config.MchID, config.MchCertificateSerialNumber, mchPrivateKey, config.ApiV3Key),
// 		option.WithWechatPayPublicKeyAuthCipher(
// 			config.MchID,                      // 商户号
// 			config.MchCertificateSerialNumber, // 商户证书序列号（需保留）
// 			mchPrivateKey,                     // 商户私钥（需保留）
// 			config.PublicKeyID,                // 公钥ID
// 			mcPublicKey,                       // 公钥内容（从.pem文件读取）
// 		),
// 	}
// 	client, err := core.NewClient(ctx, opts...)
// 	if err != nil {
// 		log.Fatalf("failed to create WeChat Pay client: %s", err)
// 	}

// 	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(
// 		config.PublicKeyID,
// 		*mcPublicKey,
// 	)

// 	parseNotifyRequestHandler, err := notify.NewRSANotifyHandler(
// 		config.ApiV3Key,
// 		verifier,
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to create WeChat Pay notify handler: %s", err)
// 	}

// 	privateKey, err := utils.LoadPrivateKeyWithPath(config.MchPrivateKeyPath)
// 	if err != nil {
// 		log.Fatalf("failed to load WeChat Pay private key: %s", err)
// 	}

// 	return &WeChat{
// 		appID:     config.AppID,
// 		offerID:   config.OfferID,
// 		appKey:    config.AppKey,
// 		appSecret: config.AppSecret,
// 		mchID:     config.MchID,
// 		apiV3Key:  config.ApiV3Key,
// 		store:     store,
// 		httpClient: resty.New().SetTimeout(10 * time.Second).SetRetryCount(3).SetRetryWaitTime(2 * time.Second).AddRetryCondition(
// 			func(r *resty.Response, err error) bool {
// 				return err != nil || (r.StatusCode() >= 500 && r.StatusCode() <= 599)
// 			},
// 		),
// 		client:                    client,
// 		parseNotifyRequestHandler: parseNotifyRequestHandler,
// 		privateKey:                privateKey,
// 		notifyURL:                 config.NotifyURL,
// 	}
// }

// // RequestMidasPaymentGameItem 虚拟支付2.0道具直购 [前端] 支付接口
// // https://developers.weixin.qq.com/minigame/dev/api/midas-payment/wx.requestMidasPaymentGameItem.html
// func (w *WeChat) RequestMidasPaymentGameItem(openID string, params map[string]interface{}) (*RequestMidasPaymentGameItem, error) {
// 	signData, err := MapToJSON(params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert params to JSON: %w", err)
// 	}

// 	method := "requestMidasPaymentGameItem"
// 	needEncodeBody := method + "&" + signData
// 	paySig := GenerateHMACSignature(w.appKey, needEncodeBody)

// 	sessionKey, err := w.store.GetSessionKey(buildKey(PrefixSes, openID))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get session key: %w", err)
// 	}
// 	signature := GenerateHMACSignature(sessionKey, signData)

// 	return &RequestMidasPaymentGameItem{
// 		SignData:  signData,
// 		PaySig:    paySig,
// 		Signature: signature,
// 	}, nil
// }

// // GetAccessToken 获取 access_token，优先从存储中获取，不存在则从微信服务器获取
// // https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
// func (w *WeChat) GetAccessToken() (string, error) {
// 	key := buildKey(PrefixAcc, w.appID)
// 	if w.store == nil {
// 		return "", fmt.Errorf("store not set")
// 	}

// 	// 1. 尝试从存储中获取
// 	var token string
// 	if token, _ = w.store.GetAccessToken(key); token != "" {
// 		return token, nil
// 	}

// 	// 2. 存储中没有，则请求微信服务器获取
// 	var tokenResp AccessTokenResponse
// 	resp, err := w.httpClient.R().
// 		SetBody(map[string]string{
// 			"grant_type": "client_credential",
// 			"appid":      w.appID,
// 			"secret":     w.appSecret,
// 		}).
// 		SetResult(&tokenResp).
// 		Post("https://api.weixin.qq.com/cgi-bin/stable_token")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to request wechat api: %w", err)
// 	}

// 	// 检查HTTP响应状态
// 	if !resp.IsSuccess() {
// 		return "", fmt.Errorf("wechat api request failed with status: %s", resp.Status())
// 	}

// 	// 检查是否有错误
// 	if tokenResp.ErrCode != 0 {
// 		return "", fmt.Errorf("wechat api error: code=%d, msg=%s", tokenResp.ErrCode, tokenResp.ErrMsg)
// 	}

// 	// 验证返回的token
// 	if tokenResp.AccessToken == "" {
// 		return "", fmt.Errorf("empty access token received")
// 	}

// 	// 保存到存储中
// 	if err := w.store.SetAccessToken(key, tokenResp.AccessToken, 3600*time.Second); err != nil {
// 		return "", fmt.Errorf("failed to save access token: %w", err)
// 	}

// 	return tokenResp.AccessToken, nil
// }

// // Code2Session 小程序登录凭证校验
// // 通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
// // https://developers.weixin.qq.com/minigame/dev/api-backend/open-api/login/auth.code2Session.html
// func (w *WeChat) Code2Session(jsCode string) (*Code2SessionResponse, error) {
// 	if jsCode == "" {
// 		return nil, fmt.Errorf("js_code cannot be empty")
// 	}

// 	var sessionResp Code2SessionResponse
// 	resp, err := w.httpClient.R().
// 		SetQueryParams(map[string]string{
// 			"appid":      w.appID,
// 			"secret":     w.appSecret,
// 			"js_code":    jsCode,
// 			"grant_type": "authorization_code",
// 		}).
// 		SetResult(&sessionResp).
// 		Get("https://api.weixin.qq.com/sns/jscode2session")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to request wechat code2session api: %w", err)
// 	}

// 	// 检查HTTP响应状态
// 	if !resp.IsSuccess() {
// 		return nil, fmt.Errorf("wechat code2session api request failed with status: %s", resp.Status())
// 	}

// 	// 检查是否有错误
// 	if sessionResp.ErrCode != 0 {
// 		return nil, fmt.Errorf("wechat code2session api error: code=%d, msg=%s", sessionResp.ErrCode, sessionResp.ErrMsg)
// 	}

// 	// 验证必要的返回字段
// 	if sessionResp.OpenID == "" || sessionResp.SessionKey == "" {
// 		return nil, fmt.Errorf("invalid response: missing openid or session_key")
// 	}

// 	// 将 session_key 保存到存储中，以 openid 为 key
// 	if w.store != nil {
// 		if err := w.store.SetSessionKey(buildKey(PrefixSes, sessionResp.OpenID), sessionResp.SessionKey, 24*time.Hour); err != nil {
// 			return nil, fmt.Errorf("failed to save session key: %w", err)
// 		}
// 	}

// 	return &sessionResp, nil
// }

// // CheckSessionKey 校验服务器所保存的登录态 session_key 是否合法
// // 为了保持 session_key 私密性，接口不明文传输 session_key，而是通过校验登录态签名完成
// // https://developers.weixin.qq.com/minigame/dev/api-backend/open-api/login/auth.checkSessionKey.html
// func (w *WeChat) CheckSessionKey(openID string) (*CheckSessionKeyResponse, error) {
// 	if openID == "" {
// 		return nil, fmt.Errorf("openid cannot be empty")
// 	}

// 	// 获取 access_token
// 	accessToken, err := w.GetAccessToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get access token: %w", err)
// 	}

// 	sessionKey, err := w.store.GetSessionKey(buildKey(PrefixSes, openID))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get session key: %w", err)
// 	}

// 	signature := GenerateHMACSignature(sessionKey, "")

// 	var checkResp CheckSessionKeyResponse
// 	resp, err := w.httpClient.R().
// 		SetQueryParams(map[string]string{
// 			"access_token": accessToken,
// 			"openid":       openID,
// 			"signature":    signature,
// 			"sig_method":   SignatureMethodHMACSHA256.String(),
// 		}).
// 		SetResult(&checkResp).
// 		Get("https://api.weixin.qq.com/wxa/checksession")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to request wechat checksession api: %w", err)
// 	}

// 	// 检查HTTP响应状态
// 	if !resp.IsSuccess() {
// 		return nil, fmt.Errorf("wechat checksession api request failed with status: %s", resp.Status())
// 	}

// 	// 检查是否有错误
// 	if checkResp.ErrCode != 0 {
// 		return nil, fmt.Errorf("wechat checksession api error: code=%d, msg=%s", checkResp.ErrCode, checkResp.ErrMsg)
// 	}

// 	return &checkResp, nil
// }

// // ResetUserSessionKey 重置指定的登录态 session_key
// // 为了保持 session_key 私密性，接口不明文传入 session_key，而是通过校验登录态签名完成
// // https://developers.weixin.qq.com/minigame/dev/api-backend/open-api/login/auth.resetUserSessionKey.html
// func (w *WeChat) ResetUserSessionKey(openID string) (*ResetUserSessionKeyResponse, error) {
// 	if openID == "" {
// 		return nil, fmt.Errorf("openid cannot be empty")
// 	}

// 	// 获取 access_token
// 	accessToken, err := w.GetAccessToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get access token: %w", err)
// 	}

// 	sessionKey, err := w.store.GetSessionKey(buildKey(PrefixSes, openID))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get session key: %w", err)
// 	}

// 	signature := GenerateHMACSignature(sessionKey, "")

// 	var resetResp ResetUserSessionKeyResponse
// 	resp, err := w.httpClient.R().
// 		SetQueryParams(map[string]string{
// 			"access_token": accessToken,
// 			"openid":       openID,
// 			"signature":    signature,
// 			"sig_method":   SignatureMethodHMACSHA256.String(),
// 		}).
// 		SetResult(&resetResp).
// 		Get("https://api.weixin.qq.com/wxa/resetusersessionkey")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to request wechat resetusersessionkey api: %w", err)
// 	}

// 	// 检查HTTP响应状态
// 	if !resp.IsSuccess() {
// 		return nil, fmt.Errorf("wechat resetusersessionkey api request failed with status: %s", resp.Status())
// 	}

// 	// 检查是否有错误
// 	if resetResp.ErrCode != 0 {
// 		return nil, fmt.Errorf("wechat resetusersessionkey api error: code=%d, msg=%s", resetResp.ErrCode, resetResp.ErrMsg)
// 	}

// 	// 验证必要的返回字段
// 	if resetResp.OpenID == "" || resetResp.SessionKey == "" {
// 		return nil, fmt.Errorf("invalid response: missing openid or session_key")
// 	}

// 	// 更新存储中的 session_key，以 openid 为 key
// 	if w.store != nil {
// 		if err := w.store.SetSessionKey(resetResp.OpenID, resetResp.SessionKey, 30*24*time.Hour); err != nil {
// 			return nil, fmt.Errorf("failed to save reset session key: %w", err)
// 		}
// 	}

// 	return &resetResp, nil
// }

// // HandleWechatMessagePushVerifyURL 处理微信消息推送的URL验证
// // 微信服务器会发送GET请求进行验证，包含signature、timestamp、nonce、echostr参数
// // 验证通过后需要原样返回echostr字符串
// // https://docs.qq.com/doc/DR1hhWlpnQXJXWHRh?_t=1749823652587&nlc=1
// func (w *WeChat) HandleWechatMessagePushVerifyURL(token, timestamp, nonce, signature string) bool {
// 	// 参数验证
// 	if token == "" || timestamp == "" || nonce == "" || signature == "" {
// 		return false
// 	}

// 	// 1. 将token、timestamp、nonce三个参数进行字典序排序
// 	tmpArr := []string{token, timestamp, nonce}
// 	sort.Strings(tmpArr)

// 	// 2. 将三个参数字符串拼接成一个字符串进行sha1加密
// 	tmpStr := strings.Join(tmpArr, "")
// 	//nolint: gosec
// 	hash := sha1.Sum([]byte(tmpStr))
// 	hashStr := fmt.Sprintf("%x", hash)

// 	// 3. 将加密后的字符串与signature对比
// 	if hashStr != signature {
// 		return false
// 	}

// 	// 4. 验证通过，返回echostr
// 	return true
// }

// // HandleMessagePushPlainText 处理微信消息推送（明文模式）
// // 验证签名后处理推送的消息内容
// // 返回值：bool表示验证是否通过，string表示回复内容，error表示处理错误
// // https://developers.weixin.qq.com/minigame/dev/guide/base-ability/message-push.html#%E6%B6%88%E6%81%AF%E8%A7%A3%E5%AF%86%E6%96%B9%E5%BC%8F%E4%B8%BA%E6%98%8E%E6%96%87%E6%A8%A1%E5%BC%8F
// // func (w *WeChat) HandleMessagePushPlainText(req *HandleMessagePushPlainTextRequest) (resp *HandleMessagePushPlainTextResponse, err error) {
// // 	// 1. 验证签名
// // 	if !w.HandleWechatMessagePushVerifyURL(req.Token, req.Timestamp, req.Nonce, req.Signature) {
// // 		return nil, fmt.Errorf("signature verification failed")
// // 	}

// // 	// 2. 解析消息体
// // 	// json.Unmarshal([]byte(req.MiniGameRequest.MiniGame), &req.MiniGameRequest)

// // 	fmt.Println(req)

// // 	return &HandleMessagePushPlainTextResponse{
// // 		ErrCode: 0,
// // 		ErrMsg:  "Success",
// // 	}, nil
// // }

// // CreatePrepayOrder 创建预支付订单
// func (w *WeChat) CreatePrepayOrder(ctx context.Context, params *PrepayRequest) (*jsapi.PrepayResponse, error) {
// 	svc := jsapi.JsapiApiService{Client: w.client}

// 	// Example request parameters
// 	// testRequest := jsapi.PrepayRequest{
// 	// 	Appid:         core.String(w.appID),
// 	// 	Mchid:         core.String(w.mchID),
// 	// 	Description:   core.String("Image形象店-深圳腾大-QQ公仔"),
// 	// 	OutTradeNo:    core.String("1217752501201407033233368018"),
// 	// 	TimeExpire:    core.Time(time.Now()),
// 	// 	Attach:        core.String("自定义数据说明"),
// 	// 	NotifyUrl:     core.String("https://www.weixin.qq.com/wxpay/pay.php"),
// 	// 	GoodsTag:      core.String("WXG"),
// 	// 	LimitPay:      []string{"LimitPay_example"},
// 	// 	SupportFapiao: core.Bool(false),
// 	// 	Amount: &jsapi.Amount{
// 	// 		Currency: core.String("CNY"),
// 	// 		Total:    core.Int64(100),
// 	// 	},
// 	// 	Payer: &jsapi.Payer{
// 	// 		Openid: core.String("oUpF8uMuAJO_M2pxb1Q9zNjWeS6o"),
// 	// 	},
// 	// 	Detail: &jsapi.Detail{
// 	// 		CostPrice: core.Int64(608800),
// 	// 		GoodsDetail: []jsapi.GoodsDetail{{
// 	// 			GoodsName:        core.String("iPhoneX 256G"),
// 	// 			MerchantGoodsId:  core.String("ABC"),
// 	// 			Quantity:         core.Int64(1),
// 	// 			UnitPrice:        core.Int64(828800),
// 	// 			WechatpayGoodsId: core.String("1001"),
// 	// 		}},
// 	// 		InvoiceId: core.String("wx123"),
// 	// 	},
// 	// 	SceneInfo: &jsapi.SceneInfo{
// 	// 		DeviceId:      core.String("013467007045764"),
// 	// 		PayerClientIp: core.String("14.23.150.211"),
// 	// 		StoreInfo: &jsapi.StoreInfo{
// 	// 			Address:  core.String("广东省深圳市南山区科技中一道10000号"),
// 	// 			AreaCode: core.String("440305"),
// 	// 			Id:       core.String("0001"),
// 	// 			Name:     core.String("腾讯大厦分店"),
// 	// 		},
// 	// 	},
// 	// 	SettleInfo: &jsapi.SettleInfo{
// 	// 		ProfitSharing: core.Bool(false),
// 	// 	},
// 	// }

// 	req := jsapi.PrepayRequest{
// 		Appid:         core.String(w.appID),
// 		Mchid:         core.String(w.mchID),
// 		Description:   core.String(params.Description),
// 		OutTradeNo:    core.String(params.OutTradeNo),
// 		TimeExpire:    core.Time(*params.TimeExpire),
// 		Attach:        core.String(params.Attach),
// 		NotifyUrl:     core.String(w.notifyURL),
// 		GoodsTag:      core.String(params.GoodsTag),
// 		LimitPay:      params.LimitPay,
// 		SupportFapiao: core.Bool(params.SupportFapiao),
// 		Amount:        params.Amount,
// 		Payer:         params.Payer,
// 		Detail:        params.Detail,
// 		SceneInfo:     params.SceneInfo,
// 		SettleInfo:    params.SettleInfo,
// 	}

// 	log.Printf("Creating WeChat prepay order with params: %+v", req)

// 	resp, result, err := svc.Prepay(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("wechat prepay order request error: %w", err)
// 	}

// 	if resp.PrepayId != nil {
// 		return resp, nil
// 	}

// 	var errResp PrepayErrorResponse
// 	bytes, err := io.ReadAll(result.Response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read error response body: %w", err)
// 	}

// 	if err = json.Unmarshal(bytes, &errResp); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
// 	}

// 	if errResp.Code != "" {
// 		return nil, fmt.Errorf("wechat prepay order error: code=%s, message=%s", errResp.Code, errResp.Message)
// 	}

// 	return nil, fmt.Errorf("wechat prepay order failed with unknown error: %v", errResp)
// }

// // ParseNotifyRequest 解析微信支付通知请求
// func (w *WeChat) ParseNotifyRequest(req *http.Request) (*notify.Request, error) {
// 	content := new(payments.Transaction)
// 	notifyReq, err := w.parseNotifyRequestHandler.ParseNotifyRequest(context.Background(), req, content)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse notify request: %w", err)
// 	}
// 	return notifyReq, nil
// }

// // SendCustomMessage 发送客服消息 https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/customer-message/receive.html#%E5%B0%8F%E7%A8%8B%E5%BA%8F%E5%8D%A1%E7%89%87%E6%B6%88%E6%81%AF
// // https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/kf-mgnt/kf-message/sendCustomMessage.html
// func (w *WeChat) SendCustomMessage(ctx context.Context, openID string, msg interface{}) error {
// 	// POST https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN
// 	if openID == "" {
// 		return fmt.Errorf("openID cannot be empty")
// 	}
// 	accessToken, err := w.GetAccessToken()
// 	if err != nil {
// 		return fmt.Errorf("failed to get access token: %w", err)
// 	}
// 	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", accessToken)

// 	body, err := json.Marshal(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal message body: %w", err)
// 	}

// 	resp, err := w.httpClient.R().
// 		SetContext(ctx).
// 		SetHeader("Content-Type", "application/json; charset=utf-8").
// 		SetBody(bytes.NewBuffer(body)).
// 		Post(url)
// 	if err != nil {
// 		return fmt.Errorf("request failed: %w", err)
// 	}

// 	result := &protocol.WxPayAndroidNotifyResponse{}
// 	if err := json.Unmarshal(resp.Body(), result); err != nil {
// 		return fmt.Errorf("failed to unmarshal response: %w", err)
// 	}

// 	if result.ErrCode != 0 {
// 		return fmt.Errorf("wechat api error: %d - %s", result.ErrCode, result.ErrMsg)
// 	}

// 	return nil
// }

// // JSAPI 签名
// func (w *WeChat) JSAPISign(timestamp, nonceStr, packageStr string) (string, error) {
// 	signData := fmt.Sprintf("%s\n%s\n%s\n%s\n", w.appID, timestamp, nonceStr, packageStr)

// 	// 计算SHA256哈希
// 	hash := sha256.Sum256([]byte(signData))

// 	// 使用RSA私钥签名
// 	signature, err := rsa.SignPKCS1v15(rand.Reader, w.privateKey, crypto.SHA256, hash[:])
// 	if err != nil {
// 		return "", fmt.Errorf("RSA签名失败: %w", err)
// 	}

// 	// Base64编码
// 	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

// 	return signatureBase64, err
// }

// // 微信屏蔽字库

// type SecurityCheckRequest struct {
// 	OpenID    string `json:"openid"`
// 	Scene     int    `json:"scene"`
// 	Version   int    `json:"version"`
// 	Content   string `json:"content"`
// 	Title     string `json:"title,omitempty"`
// 	Nickname  string `json:"nickname,omitempty"`
// 	Signature string `json:"signature,omitempty"`
// }

// type SecurityCheckResponse struct {
// 	Errcode int    `json:"errcode"`
// 	Errmsg  string `json:"errmsg"`
// 	Result  struct {
// 		Suggest string `json:"suggest"`
// 		Label   int    `json:"label"`
// 	} `json:"result"`
// 	Detail []struct {
// 		Strategy string `json:"strategy"`
// 		Errcode  int    `json:"errcode"`
// 		Suggest  string `json:"suggest"`
// 		Label    int    `json:"label"`
// 		Prob     int    `json:"prob,omitempty"`
// 		Keyword  string `json:"keyword,omitempty"`
// 	} `json:"detail"`
// 	TraceID string `json:"trace_id"`
// }

// func checkMessageSecurity(accessToken string, request SecurityCheckRequest) (*SecurityCheckResponse, error) {
// 	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", accessToken)
// 	jsonData, err := json.Marshal(request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response SecurityCheckResponse
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &response, nil
// }

// func (w *WeChat) CheckMessageSecurity(openID string, content string) error {
// 	if openID == "" {
// 		return fmt.Errorf("openID cannot be empty")
// 	}
// 	accessToken, err := w.GetAccessToken()
// 	if err != nil {
// 		return fmt.Errorf("failed to get access token: %w", err)
// 	}

// 	request := SecurityCheckRequest{
// 		OpenID:    openID,           // 替换为用户的 OpenID
// 		Scene:     1,                // 场景枚举值
// 		Version:   2,                // 接口版本号
// 		Content:   content,          // 需检测的文本内容
// 		Title:     "Test Title",     // 可选：文本标题
// 		Nickname:  "Test Nickname",  // 可选：用户昵称
// 		Signature: "Test Signature", // 可选：个性签名
// 	}

// 	response, err := checkMessageSecurity(accessToken, request)
// 	if err != nil {
// 		fmt.Printf("Failed to check message security: %v\n", err)
// 		return fmt.Errorf("Failed to check message security")
// 	}

// 	fmt.Printf("Response: %+v\n", response)

// 	if response.Result.Suggest != "pass" {
// 		return fmt.Errorf("所输入内容含有非法字符！")
// 	}
// 	return nil
// }
