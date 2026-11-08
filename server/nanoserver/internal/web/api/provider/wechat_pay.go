package provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const (
	AppKey  = "qvwgpt6PD5lL2vDQeTdJksvUv5vfYdIz" // 替换为实际的AppKey
	OfferId = 1450393955
)

type WeChatPay struct {
	appKey     string // 米大师密钥
	sessionKey string // 小程序session_key
}

func NewWeChatPay(appKey, sessionKey string) *WeChatPay {
	return &WeChatPay{
		appKey:     appKey,
		sessionKey: sessionKey,
	}
}

// PaySig 虚拟支付2.0道具直购支付签名
// https://developers.weixin.qq.com/minigame/dev/api/midas-payment/wx.requestMidasPaymentGameItem.html
func (w *WeChatPay) PaySig(signData string) (signature string, err error) {
	method := "requestMidasPaymentGameItem"
	// signData = string(rawData) // 直接使用传入的JSON字符串作为签名数据
	needEncodeBody := method + "&" + signData
	h := hmac.New(sha256.New, []byte(w.appKey))
	h.Write([]byte(needEncodeBody))
	signature = hex.EncodeToString(h.Sum(nil))

	return
}

// Signature 虚拟支付2.0道具直购 用户态签名
// https://developers.weixin.qq.com/minigame/dev/guide/open-ability/signature.html#%E7%94%A8%E6%88%B7%E7%99%BB%E5%BD%95%E6%80%81%E7%AD%BE%E5%90%8D
func (w *WeChatPay) Signature(signData string) (string, error) {
	// 将参数转换为JSON字符串作为rawData
	// rawData, err := json.Marshal(params)
	// if err != nil {
	// 	return "", err
	// }

	// 使用HMAC-SHA256生成签名: signature = hmac_sha256(session_key, rawData)
	h := hmac.New(sha256.New, []byte(w.sessionKey))
	h.Write([]byte(signData))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature, nil
}

// func (w *WeChatPay) CreateOrderResponse(order *model.Order) (*RequestMidasPaymentGameItem, error) {
// 	// 传入必要的参数
// 	params, err := utils.ConvertJSONToMap(order.Extra)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert order extra to map: %w", err)
// 	}

// 	// params["mode"] = "goods" // 支付模式
// 	// signData := fmt.Sprintf(`{"mode":"%v","offerId":"%v","buyQuantity":%v,"env":%v,"currencyType":"%v","platform":"%v","zoneId":"%v","productId":"%v","goodsPrice":%v,"outTradeNo":"%v","attach":"%v"}`,
// 	signData := fmt.Sprintf(`{"mode":"%v","offerId":"%v","buyQuantity":%v,"env":%v,"currencyType":"%v","platform":"%v","productId":"%v","goodsPrice":%v,"outTradeNo":"%v"}`,
// 		params["mode"],
// 		OfferId, // offerId
// 		params["buyQuantity"],
// 		params["env"],
// 		params["currencyType"],
// 		params["platform"],
// 		// params["zoneId"],
// 		params["productId"],
// 		params["goodsPrice"],
// 		order.OrderId,
// 	// params["attach"]
// 	)
// 	paySig, err := w.PaySig(signData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate pay signature: %w", err)
// 	}

// 	signature, err := w.Signature(signData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate user signature: %w", err)
// 	}

// 	requestMidasPaymentGameItem := &RequestMidasPaymentGameItem{
// 		SignData:  signData,
// 		PaySig:    paySig,
// 		Signature: signature,
// 	}

// 	logger.Debugf("Player %v: PaySig: %s", order.Uid, paySig)
// 	logger.Debugf("Player %v: Signature: %s", order.Uid, signature)
// 	logger.Debugf("Player %v: SignData: %s", order.Uid, signData)

// 	return requestMidasPaymentGameItem, nil
// }

// RequestMidasPaymentGameItem 微信小游戏虚拟支付请求参数
type RequestMidasPaymentGameItem struct {
	SignData  string `json:"signData"`  // 签名数据JSON字符串
	PaySig    string `json:"paySig"`    // 支付签名
	Signature string `json:"signature"` // 用户态签名
}

// PaymentSignData 支付签名数据结构
type PaymentSignData struct {
	Mode         string `json:"mode"`         // 支付模式，如"goods"
	OfferID      string `json:"offerId"`      // 商品ID
	BuyQuantity  int    `json:"buyQuantity"`  // 购买数量
	Env          int    `json:"env"`          // 环境，0-正式环境，1-沙箱环境
	CurrencyType string `json:"currencyType"` // 货币类型，如"CNY"
	Platform     string `json:"platform"`     // 平台，如"android"
	ZoneID       string `json:"zoneId"`       // 区域ID
	ProductID    string `json:"productId"`    // 产品ID
	GoodsPrice   int    `json:"goodsPrice"`   // 商品价格（分）
	OutTradeNo   string `json:"outTradeNo"`   // 商户订单号
	Attach       string `json:"attach"`       // 附加数据
}
