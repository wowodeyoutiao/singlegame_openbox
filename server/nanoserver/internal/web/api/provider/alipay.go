package provider

import (
	"context"
	"nanoserver/global"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
)

type AlipayMgr struct {
	privateKey    string
	appId         string
	unifyOrderURL string
	callbackURL   string
	client        *alipay.Client
}

var logger = log.WithFields(log.Fields{"component": "http", "service": "pay"})
var Alipay = &AlipayMgr{}

func (ali *AlipayMgr) Init() {
	ali.appId = viper.GetString("alipay.appid")
	ali.privateKey = viper.GetString("alipay.appsecret")
	ali.unifyOrderURL = viper.GetString("alipay.unify_order_url")
	ali.callbackURL = viper.GetString("alipay.callback_url")

	ali.client, _ = alipay.New(ali.appId, ali.privateKey, true)
	// 加载应用公钥证书
	if err := ali.client.LoadAppCertPublicKeyFromFile(global.BasePath + "/configs/cert/appCertPublicKey_2021004185615400.crt"); err != nil {
		logger.Errorf("加载应用公钥证书失败 %s ", err.Error())
	}
	// 加载支付宝公钥证书
	if err := ali.client.LoadAlipayCertPublicKeyFromFile(global.BasePath + "/configs/cert/alipayCertPublicKey_RSA2.crt"); err != nil {
		logger.Errorf("加载支付宝公钥证书 %s", err.Error())
	}
	// 加载支付宝根证书
	if err := ali.client.LoadAliPayRootCertFromFile(global.BasePath + "/configs/cert/alipayRootCert.crt"); err != nil {
		logger.Errorf("加载支付宝根证书 %s", err.Error())
	}
	logger.Debugln("支付宝证书初始化完成")
}

// 创建订单
// func (ali *AlipayMgr) CreateOrderResponse(order *model.Order) (string, error) {
// 	pay := alipay.TradeAppPay{}
// 	pay.ReturnURL = ali.callbackURL   // 支付后调转页面
// 	pay.NotifyURL = ali.unifyOrderURL // 支付宝回调
// 	pay.Subject = strconv.Itoa(order.ProductId)
// 	pay.OutTradeNo = strings.Replace(uuid.New(), "-", "", -1) // 传递一个唯一单号
// 	pay.TotalAmount = strconv.Itoa(order.Money)               // 金额
// 	pay.ProductCode = "QUICK_MSECURITY_PAY"                   // APP支付

// 	url, err := ali.client.TradeAppPay(pay)
// 	if err != nil {
// 		return "", err
// 	}
// 	payURL := url
// 	log.Println("支付页面URL:", payURL)

// 	return payURL, nil
// }

// 单笔转账
func (ali *AlipayMgr) FundTransUniTransfer(ctx context.Context, p alipay.FundTransUniTransfer) error {
	// var p = alipay.FundTransUniTransfer{
	// 	OutBizNo:    "商户端的唯一订单号",
	// 	TransAmount: "0.1",
	// 	ProductCode: "TRANS_ACCOUNT_NO_PWD",
	// 	BizScene:    "DIRECT_TRANSFER",
	// 	PayeeInfo: &alipay.PayeeInfo{
	// 		Identity:     "支付宝账号",
	// 		IdentityType: "ALIPAY_LOGON_ID",
	// 		Name:         "真实姓名",
	// 	},
	// 	OrderTitle: "转账业务的标题",
	// 	Remark:     "转账备注",
	// }

	p.ProductCode = "TRANS_ACCOUNT_NO_PWD"
	p.BizScene = "DIRECT_TRANSFER"
	result, err := ali.client.FundTransUniTransfer(ctx, p)
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Printf("转账结果: %+v\n", result)

	return nil
}

func (ali *AlipayMgr) VerifySign(r *http.Request) error {

	return ali.client.VerifySign(r.Form)
}
