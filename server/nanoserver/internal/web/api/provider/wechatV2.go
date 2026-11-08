package provider

// import (
// 	"os"
// 	"path/filepath"

// 	"nanoserver/internal/web/api/provider/wechatv2"
// 	"nanoserver/pkg/redis_factory"

// 	"github.com/spf13/viper"
// )

// var WeChatV2 *wechatv2.WeChat

// func WeChatV2Init() {
// 	// appid := viper.GetString("wechatv2.appid")
// 	// offerID := viper.GetInt("wechatv2.offerid")
// 	// appKey := viper.GetString("wechatv2.appkey")
// 	// appSecret := viper.GetString("wechatv2.appsecret")
// 	os.Setenv("WECHATPAY_LOG_LEVEL", "debug")

// 	WeChatV2 = wechatv2.New(
// 		&wechatv2.Config{
// 			OfferID:                    viper.GetInt("wechatv2.offer_id"),
// 			AppID:                      viper.GetString("wechatv2.app_id"),
// 			AppKey:                     viper.GetString("wechatv2.app_key"),
// 			AppSecret:                  viper.GetString("wechatv2.app_secret"),
// 			MchID:                      viper.GetString("wechatv2.mc_id"),
// 			MchCertificateSerialNumber: viper.GetString("wechatv2.mch_certificate_serial_number"),
// 			ApiV3Key:                   viper.GetString("wechatv2.api_v3_key"),
// 			MchPrivateKeyPath:          filepath.Join(".", "configs", "cert", "apiclient_key.pem"),
// 			PublicKeyID:                viper.GetString("wechatv2.public_key_id"),
// 			PublicKey:                  viper.GetString("wechatv2.public_key"),
// 			NotifyURL:                  viper.GetString("wechatv2.notify_url"),
// 		},
// 		redis_factory.NewTokenCache(),
// 	)
// }
