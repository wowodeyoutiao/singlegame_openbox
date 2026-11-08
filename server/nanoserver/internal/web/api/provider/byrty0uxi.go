package provider

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"nanoserver/pkg/errutil"
	"nanoserver/protocol"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

// Byrty0uxi 签名，MD5(key=value&key1=value1&.....&key=密钥) 所有请求参数按key的ascii码排序
type Byrty0uxi struct {
	Name           string `json:"name"`           // 项目名称，name,和密钥从外部获取
	Subject        string `json:"subject"`        // 订单标题
	OutTradeNo     string `json:"outTradeNo"`     // 商户订单号。由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	TotalAmount    string `json:"totalAmount"`    // 订单总金额。单位为元，精确到小数点后两位，取值范围：[0.01,100000000]
	Body           string `json:"body"`           // 订单附加信息
	ClientIp       string `json:"clientIp"`       // 客户端下单IP
	NotifyUrl      string `json:"notifyUrl"`      // 通知异步回调地址
	ProductCode    string `json:"productCode"`    // APP支付：QUICK_MSECURITY_PAY ；网站支付:QUICK_WAP_WAY，默认APP支付
	TimeoutExpress string `json:"timeoutExpress"` // 订单关闭时间 取值范围：5m～15d。m-分钟，h-小时，d-天，1c-当天，默认30m
	Uid            string `json:"uid"`            // 用户ID,对于投诉用户风控，关联支付宝ID
}

type Byrty0uxiMgr struct {
	*Byrty0uxi
	payUrl      string //支付接口
	transUrl    string //转账接口
	realNameUrl string //实名认证接口

	appName   string
	appKey    string
	notifyUrl string
}

var Byrtypay = &Byrty0uxiMgr{
	payUrl:      "http://pay.byrty0uxi.com/api/order/createOrder",
	transUrl:    "http://pay.byrty0uxi.com/api/order/createTransferOrder",
	appName:     "bZsVg",
	appKey:      "3652168abb4f7fd46967d965bc76e22d",
	notifyUrl:   "http://127.0.0.1:12307/v1/order/notify/byrty0uxi",
	realNameUrl: "http://47.102.156.146/realname/realname/realnameReq",
}

func (mgr *Byrty0uxiMgr) Init() {
	mgr.appName = viper.GetString("byrty0uxi.appid")
	mgr.appKey = viper.GetString("byrty0uxi.appsecret")
	mgr.notifyUrl = viper.GetString("byrty0uxi.unify_order_url")
	mgr.payUrl = viper.GetString("byrty0uxi.pay_url")
	mgr.transUrl = viper.GetString("byrty0uxi.trans_url")
}

type OrderInfo struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
	Data    string `json:"data"`
}

// 创建订单
// func (mgr *Byrty0uxiMgr) CreateOrderResponse(order *model.Order) (string, error) {
// 	amount := float64(order.RealMoney)
// 	if viper.GetBool("core.debug") {
// 		amount = 0.01
// 	}
// 	b := &Byrty0uxi{
// 		Name:           mgr.appName,
// 		Subject:        order.ProductName,
// 		OutTradeNo:     order.OrderId,
// 		TotalAmount:    fmt.Sprintf("%.2f", amount),
// 		Body:           fmt.Sprintf("%d", order.ProductId),
// 		ClientIp:       order.Remote,
// 		NotifyUrl:      mgr.notifyUrl,
// 		ProductCode:    "QUICK_MSECURITY_PAY",
// 		TimeoutExpress: "5m",
// 		Uid:            strconv.Itoa(int(order.Uid)),
// 	}

// 	// 签名map
// 	m := make(map[string]interface{}, 0)
// 	m["name"] = b.Name
// 	m["subject"] = b.Subject
// 	m["outTradeNo"] = b.OutTradeNo
// 	m["totalAmount"] = b.TotalAmount
// 	m["body"] = b.Body
// 	m["clientIp"] = b.ClientIp
// 	m["notifyUrl"] = b.NotifyUrl
// 	m["productCode"] = b.ProductCode
// 	m["timeoutExpress"] = b.TimeoutExpress
// 	m["uid"] = b.Uid

// 	sign, err := mgr.signCalculator(m)
// 	if err != nil {
// 		return "sign failed ", err
// 	}

// 	formData := url.Values{}
// 	formData.Set("name", b.Name)
// 	formData.Set("subject", b.Subject)
// 	formData.Set("outTradeNo", b.OutTradeNo)
// 	formData.Set("totalAmount", b.TotalAmount)
// 	formData.Set("body", b.Body)
// 	formData.Set("clientIp", b.ClientIp)
// 	formData.Set("notifyUrl", b.NotifyUrl)
// 	formData.Set("productCode", b.ProductCode)
// 	formData.Set("timeoutExpress", b.TimeoutExpress)
// 	formData.Set("uid", b.Uid)
// 	formData.Set("sign", sign)

// 	request, err := http.NewRequest(http.MethodPost, mgr.payUrl, strings.NewReader(formData.Encode()))
// 	if err != nil {
// 		logger.Errorf("create unify order failed: %s", err.Error())
// 		return "", err
// 	}
// 	parma := request.URL.Query()
// 	request.URL.RawQuery = parma.Encode()

// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	c := http.Client{}
// 	response, err := c.Do(request)
// 	if err != nil {
// 		logger.Errorf("request unify order failed: %s", err.Error())
// 		return "", err
// 	}
// 	defer response.Body.Close()

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	resp := &OrderInfo{}
// 	if err := json.Unmarshal(body, &resp); err != nil {
// 		logger.Errorf("unify order response unmarsharl failed: %s", err.Error())
// 		return "", err
// 	}

// 	return resp.Message, nil
// }

type Response struct {
	Status  int    `json:"Status"`  // 1成功，其他失败
	Message string `json:"Message"` // 描述
	Data    Data   `json:"data"`    // 包含订单信息的数据结构
}

type Data struct {
	OrderId        string      `json:"OrderId"`        // 支付宝订单号
	PayFundOrderId string      `json:"PayFundOrderId"` // 支付资金订单号
	SubCode        interface{} `json:"SubCode"`        // 子代码
	SubMsg         interface{} `json:"SubMsg"`         // 子消息
	Body           string      `json:"body"`           // 包含具体响应内容的JSON字符串
	AlipayCertSn   string      `json:"alipay_cert_sn"` // 支付宝证书序列号
	Sign           string      `json:"sign"`           // 签名
}
type BodyContent struct {
	AlipayFundTransUniTransferResponse struct {
		Code           string `json:"code"`              // 响应代码
		Msg            string `json:"msg"`               // 响应消息
		OrderId        string `json:"order_id"`          // 订单ID
		OutBizNo       string `json:"out_biz_no"`        // 外部业务编号
		PayFundOrderId string `json:"pay_fund_order_id"` // 支付资金订单号
		Status         string `json:"status"`            // 状态
		TransDate      string `json:"trans_date"`        // 交易日期
	} `json:"alipay_fund_trans_uni_transfer_response"`
	AlipayCertSn string `json:"alipay_cert_sn"` // 支付宝证书序列号
	Sign         string `json:"sign"`           // 签名
}

// 单笔转账
func (mgr *Byrty0uxiMgr) FundTransUniTransfer(p protocol.FundTransParams) error {
	// 签名map
	m := make(map[string]interface{}, 0)
	m["name"] = mgr.appName
	m["aliName"] = p.AliName
	m["identity"] = p.Identity
	m["identityType"] = "ALIPAY_LOGON_ID"
	m["orderTitle"] = p.OrderTitle
	m["outBizNo"] = p.OutBizNo
	m["remark"] = p.Remark
	m["transAmount"] = p.TransAmount

	sign, err := mgr.signCalculator(m)
	if err != nil {
		return err
	}
	formData := url.Values{}
	formData.Set("name", mgr.appName)
	formData.Set("aliName", p.AliName)
	formData.Set("identity", p.Identity)
	formData.Set("identityType", "ALIPAY_LOGON_ID")
	formData.Set("orderTitle", p.OrderTitle)
	formData.Set("outBizNo", p.OutBizNo)
	formData.Set("remark", p.Remark)
	formData.Set("transAmount", p.TransAmount)
	formData.Set("sign", sign)

	request, err := http.NewRequest(http.MethodPost, mgr.transUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		logger.Errorf("create unify order failed: %s", err.Error())
		return err
	}
	parma := request.URL.Query()
	request.URL.RawQuery = parma.Encode()

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := http.Client{}
	response, err := c.Do(request)
	if err != nil {
		logger.Errorf("request unify order failed: %s", err.Error())
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	resp := &Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		logger.Errorf("unify order response unmarsharl failed: %s", err.Error())
		return err
	}
	if resp.Status != 1 {
		logger.Errorf("转账失败: Status: %d", resp.Status)
		return errutil.ErrTrans
	}
	respbody := &BodyContent{}
	err = json.Unmarshal([]byte(resp.Data.Body), respbody)
	if err != nil {
		return err
	}
	logger.Debugf("转账成功: AliName: %s Identity %s  TransAmount: %s", p.AliName, p.Identity, p.TransAmount)
	return nil
}

// 微信支付计算签名的函数
func (mgr *Byrty0uxiMgr) signCalculator(params map[string]interface{}) (string, error) {
	key := mgr.appKey
	if params == nil || key == "" {
		return "", errutil.ErrIllegalParameter
	}

	sortedKeys := make([]string, 0)
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	buf := &bytes.Buffer{}
	for _, k := range sortedKeys {
		if params[k] == nil {
			continue
		}

		switch params[k].(type) {
		case string:
			if len(params[k].(string)) == 0 {
				continue
			}

		case int:
			if params[k].(int) == 0 {
				continue
			}
		}

		fmt.Fprintf(buf, "%s=%v&", k, params[k])
	}
	fmt.Fprintf(buf, "key=%s", key)
	md5Ctx := md5.New()
	signStr := buf.Bytes()
	md5Ctx.Write(signStr)
	cipherStr := md5Ctx.Sum(nil)
	sign := strings.ToLower(hex.EncodeToString(cipherStr))
	return sign, nil
}

func (mgr *Byrty0uxiMgr) CheckSign(params map[string]interface{}, sign string) bool {
	mySign, err := mgr.signCalculator(params)
	if err != nil {
		return false
	}
	if mySign != sign {
		return false

	}
	return true
}
