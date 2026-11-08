package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"nanoserver/db"
	"nanoserver/dispatcher"
	"nanoserver/internal/web/api"
	"nanoserver/internal/web/api/provider"
	"nanoserver/pkg/algoutil"
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/jwtv5x"
	"nanoserver/pkg/my_jwt"
	"nanoserver/pkg/redis_factory"
	"nanoserver/pkg/whitelist"
	"nanoserver/protocol"

	_ "net/http/pprof"

	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nex"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Closer func()

var logger = log.WithField("component", "http")

// var userToken string // 广播消息
func dbStartup() func() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.args"))

	return db.MustStartup(
		dsn,
		db.MaxIdleConns(viper.GetInt("database.max_idle_conns")),
		db.MaxOpenConns(viper.GetInt("database.max_open_conns")),
		db.ShowSQL(viper.GetBool("database.show_sql")))
}

func enableWhiteList() {
	whitelist.Setup(viper.GetStringSlice("whitelist.ip"))
}

func version() (*protocol.Version, error) {
	return &protocol.Version{
		Version: viper.GetInt("update.version"),
		Android: viper.GetString("update.android"),
		IOS:     viper.GetString("update.ios"),
	}, nil
}

func pongHandler() (string, error) {
	return "pong", nil
}

// func logRequest(ctx context.Context, r *http.Request) (context.Context, error) {
// 	var (
// 		requestBody interface{}
// 		formData    interface{}
// 	)
// 	// 如果 URI 不是 "/ping"，则打印请求信息
// 	if uri := r.RequestURI; uri != "/ping" {
// 		// 打印请求方法、远程地址、URI 和上下文
// 		// logger.Debugf("Method=%s, RemoteAddr=%s, URL=%s, ctx=%v", r.Method, r.RemoteAddr, uri, ctx)
// 		// 打印查询参数
// 		// logger.Debugf("Query Parameters: %v", queryParams)
// 		// var (
// 		// queryParams = r.URL.Query()
// 		// requestBody interface{}
// 		// )

// 		// 打印表单数据（如果是 POST 请求）
// 		if r.Method == http.MethodPost {
// 			if err := r.ParseForm(); err != nil {
// 				logger.Warnf("Failed to parse form data: %v", err)
// 			} else {
// 				formData = r.Form
// 			}
// 		}

// 		// 打印 JSON 请求体（如果存在）
// 		if r.Body != nil {
// 			bodyBytes, err := ioutil.ReadAll(r.Body)
// 			if err != nil {
// 				logger.Warnf("Failed to read request body: %v", err)
// 			} else {
// 				// 尝试解析 JSON 请求体
// 				var bodyMap map[string]interface{}
// 				if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
// 					// logger.Debugf("JSON Request Body: %v", bodyMap)
// 					// logger.Debugf("logRequest JSON Request Body: %s", string(bodyBytes))
// 					requestBody = bodyMap
// 				} else {
// 					// logger.Debugf("logRequest Raw Request Body: %s", string(bodyBytes))
// 					requestBody = string(bodyBytes)
// 				}
// 				// 恢复请求体，以便后续处理
// 				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
// 			}
// 		}
// 	}

// 	// 打印请求信息
// 	logger.Debugf("HTTP Request: Method=%s, RemoteAddr=%s, URL=%s, QueryParams=%+#v, RequestBody=%+#v, FormData=%+#v", r.Method, r.RemoteAddr, r.RequestURI, r.URL.Query(), requestBody, formData)

// 	return ctx, nil
// }

func logRequest(ctx context.Context, r *http.Request) (context.Context, error) {
	if r.RequestURI == "/ping" {
		return ctx, nil
	}

	var (
		requestBody interface{}
		formData    interface{}
		queryParams = r.URL.Query()
	)

	// Log basic request info
	logger.Debugf("HTTP Request: Method=%s, RemoteAddr=%s, URL=%s, QueryParams=%+#v",
		r.Method, r.RemoteAddr, r.RequestURI, queryParams)

	// Handle form data for POST requests
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			formData = r.Form
		}
	}

	// Read and restore request body
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warnf("Failed to read request body: %v", err)
		} else {
			// Try to parse as JSON
			var bodyMap map[string]interface{}
			if json.Unmarshal(bodyBytes, &bodyMap) == nil {
				requestBody = bodyMap
			} else {
				requestBody = string(bodyBytes)
			}
			// Restore body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	logger.Debugf("HTTP Request Summary: Method=%s, RemoteAddr=%s, URL=%s, QueryParams=%+#v, RequestBody=%+#v, FormData=%+#v",
		r.Method,
		r.RemoteAddr,
		r.RequestURI,
		queryParams,
		requestBody,
		formData,
	)

	return ctx, nil
}

// // TokenMiddleware  是一个中间件，用于验证请求中的 token
// func TokenMiddleware(ctx context.Context, r *http.Request) (context.Context, error) {
// 	routerToCheck := strings.Split(r.RequestURI, "?")[0]
// 	// log.Printf("ctx", ctx)
// 	// log.Printf("routerToCheck", routerToCheck)
// 	// log.Printf("VerifyUrl", VerifyUrl(routerToCheck))
// 	if VerifyUrl(routerToCheck) {
// 		return ctx, nil
// 	}
// 	// 从请求头中获取 token
// 	utoken := r.Header.Get("Authorization")
// 	// 检查字符串是否以"Bearer "开头
// 	if strings.HasPrefix(utoken, "Bearer ") {
// 		// 去除"Bearer "前缀
// 		utoken = utoken[len("Bearer "):]
// 	}
// 	// 使用 token 服务来验证 token
// 	tokenFactory := token.CreateUserFactory()
// 	check := tokenFactory.IsEffective(utoken)
// 	// log.Printf("check", check)
// 	if !check {
// 		return ctx, errutil.ErrPermissionDenied
// 	}
// 	// log.Printf("Error generating token: %v", check)

// 	// 解析token
// 	userInfo, err := tokenFactory.ParseToken(utoken)
// 	if err != nil {
// 		return ctx, err
// 	}
// 	// 创建一个MyData结构体实例并填充数据
// 	customClaims := my_jwt.CustomClaims{
// 		UserId:      userInfo.UserId,
// 		AccountName: userInfo.AccountName,
// 		ServerId:    userInfo.ServerId,
// 		Platform:    userInfo.Platform,
// 		Name:        userInfo.Name,
// 	}
// 	ctx = context.WithValue(ctx, "userDate", customClaims)
// 	r = r.WithContext(ctx)

// 	log.Debugf("TokenMiddleware: UserId=%d, AccountName=%s, ServerId=%d, Platform=%s, Name=%s", customClaims.UserId, customClaims.AccountName, customClaims.ServerId, customClaims.Platform, customClaims.Name)

//		return ctx, nil
//	}
//

// TokenMiddleware 是一个中间件，用于验证请求中的 token
// func TokenMiddleware(ctx context.Context, r *http.Request) (context.Context, error) {
// 	// 提取请求路径（去除查询参数）
// 	routerToCheck := strings.Split(r.RequestURI, "?")[0]

// 	// 如果该路径不需要验证 token，直接放行
// 	if VerifyUrl(routerToCheck) {
// 		return ctx, nil
// 	}

// 	// 从请求头中获取 token
// 	utoken := r.Header.Get("Authorization")

// 	// 去除 "Bearer " 前缀（如果存在）
// 	if strings.HasPrefix(utoken, "Bearer ") {
// 		utoken = utoken[len("Bearer "):]
// 	}

// 	// 使用 token 工厂验证 token 是否有效
// 	tokenFactory := token.CreateUserFactory()
// 	if !tokenFactory.IsEffective(utoken) {
// 		return ctx, errutil.ErrPermissionDenied
// 	}

// 	// 解析 token 获取用户信息
// 	userInfo, err := tokenFactory.ParseToken(utoken)
// 	if err != nil {
// 		return ctx, err
// 	}

// 	// 构建用户信息并注入上下文
// 	customClaims := my_jwt.CustomClaims{
// 		UserId:      userInfo.UserId,
// 		AccountName: userInfo.AccountName,
// 		ServerId:    userInfo.ServerId,
// 		Platform:    userInfo.Platform,
// 		Name:        userInfo.Name,
// 	}

// 	ctx = context.WithValue(ctx, "userDate", customClaims)
// 	r = r.WithContext(ctx)

// 	// 记录中间件处理结果
// 	log.Debugf("TokenMiddleware: UserId=%d, AccountName=%s, ServerId=%d, Platform=%s, Name=%s",
// 		customClaims.UserId, customClaims.AccountName, customClaims.ServerId, customClaims.Platform, customClaims.Name)

// 	return ctx, nil
// }

// TokenMiddleware 是一个中间件，用于验证请求中的 token
func TokenMiddleware(ctx context.Context, r *http.Request) (context.Context, error) {
	// 提取请求路径（去除查询参数）
	routerToCheck := strings.Split(r.RequestURI, "?")[0]

	// 如果该路径不需要验证 token，直接放行
	if VerifyUrl(routerToCheck) {
		return ctx, nil
	}

	// 从请求头中获取 token
	utoken := r.Header.Get("Authorization")

	// 去除 "Bearer " 前缀（如果存在）
	utoken = strings.TrimPrefix(utoken, "Bearer ")

	gameClaims, err := jwtv5x.Instance.VerifyToken(context.TODO(), utoken)
	if err != nil {
		return ctx, errutil.ErrPermissionDenied
	}

	// 将 UserId 从字符串转换为 int64
	userID, err := strconv.ParseInt(gameClaims.UserId, 10, 64)
	if err != nil {
		return ctx, errutil.ErrPermissionDenied
	}

	// 构建用户信息并注入上下文
	customClaims := my_jwt.CustomClaims{
		UserId:      userID,
		AccountName: gameClaims.AccountName,
		ServerId:    gameClaims.ServerId,
		Platform:    gameClaims.Platform,
		Name:        gameClaims.Name,
		Ad:          gameClaims.Ad,
	}

	ctx = context.WithValue(ctx, "userDate", customClaims)
	// r = r.WithContext(ctx)

	// 记录中间件处理结果
	log.Debugf("TokenMiddleware: UserId=%d, AccountName=%s, ServerId=%d, Platform=%s, Name=%s",
		customClaims.UserId, customClaims.AccountName, customClaims.ServerId, customClaims.Platform, customClaims.Name)

	return ctx, nil
}

// VerifyUrl 检查给定的 url 是否在白名单中
func VerifyUrl(url string) bool {
	_, exists := RouterWhitelist[url]
	return exists
}

// RouterWhitelist 存储允许的 路由
var RouterWhitelist = map[string]struct{}{
	"/v1/user/login/guest":                        {},
	"/v1/gm/combat/setting":                       {},
	"/v1/gm/combat/result":                        {},
	"/v1/window/register_remain":                  {},
	"/v1/window/payment_remain":                   {},
	"/v1/window/pay_per":                          {},
	"/v1/window/ltv":                              {},
	"/v1/window/ingot":                            {},
	"/v1/window/online_log":                       {},
	"/v1/window/have_money":                       {},
	"/v1/window/reward":                           {},
	"/v1/window/gold":                             {},
	"/v1/window/soul":                             {},
	"/v1/window/shushuno":                         {},
	"/v1/window/shushuregister":                   {},
	"/v1/user/login/3rd":                          {},
	"/v1/order/notify/byrty0uxi":                  {},
	"/v1/mail/add":                                {},
	"/v1/order/wx/ios/callbacks/customer/message": {},
	"/v1/order/wx/android/callbacks/paymentns":    {},
	"/v1/order/wx/ios/callbacks/payments":         {},
	"/v1/window/tenxun":                           {},
	"/v1/window/send_lamp":                        {},
	"/v1/window/send_mail":                        {},
	"/v1/window/send_mute":                        {},
	"/v1/window/ban_user":                         {},
	"/v1/window/delete_user":                      {},
	"/v1/window/query_user_info":                  {},
}

func startupService() http.Handler {
	var (
		mux    = http.NewServeMux()
		webDir = viper.GetString("webserver.static_dir")
	)

	nex.Before(logRequest, TokenMiddleware, userLockMiddleware)
	nex.After(userUnLockMiddleware)
	mux.Handle("/v1/user/", api.MakeLoginService())
	mux.Handle("/v1/order/", api.MakeOrderService())
	mux.Handle("/v1/history/", api.MakeHistoryService())
	mux.Handle("/v1/desk/", api.MakeDeskService())
	mux.Handle("/v1/version", nex.Handler(version))
	mux.Handle("/v1/mail/", api.MakeMailService())
	mux.Handle("/v1/sign/", api.MakeSignService())
	mux.Handle("/v1/store/", api.MakeStoreService())
	mux.Handle("/v1/activation/", api.MakeActivationService())
	mux.Handle("/v1/activity/", api.MakeActivityService())
	mux.Handle("/v1/fate/", api.MakeFateService())
	mux.Handle("/v1/window/", api.MakeWindowService())
	mux.Handle("/v1/sensitive/", api.MakeSensitiveService())

	mux.Handle("/v1/privilege/", api.MakePrivilegeService())
	mux.Handle("/v1/rank/", api.MakeRankService())
	mux.Handle("/v1/prize_draw/", api.MakePrizeDrawService())
	mux.Handle("/v1/pay/", api.MakePayService()) // 没使用 nex 直接使用了标准库

	// GM系统命令
	mux.Handle("/v1/gm/reset", nex.Handler(resetPlayerHandler).Before(authFilter))   // 重置玩家未完成房间状态
	mux.Handle("/v1/gm/consume", nex.Handler(cardConsumeHandler).Before(authFilter)) // 设置消耗
	mux.Handle("/v1/gm/broadcast", nex.Handler(broadcast).Before(authFilter))        // 消息广播
	mux.Handle("/v1/gm/kick", nex.Handler(kickHandler).Before(authFilter))           // 踢人
	mux.Handle("/v1/gm/online", nex.Handler(onlineHandler).Before(authFilter))       // 在线信息
	mux.Handle("/v1/gm/recharge", nex.Handler(rechargeHandler).Before(authFilter))   // 玩家充值
	mux.Handle("/v1/gm/query/user/", nex.Handler(userInfoHandler))                   // 玩家信息查询
	mux.Handle("/v1/gm/combat/setting", nex.Handler(combatSetting))                  // 对战信息设置
	mux.Handle("/v1/gm/combat/result", nex.Handler(combatResult))                    // 对战结果查询

	// 统计后台
	mux.Handle("/v1/stats/user/register", nex.Handler(registerUsersHandler).Before(authFilter))     // 注册人数
	mux.Handle("/v1/stats/user/activation", nex.Handler(activationUsersHandler).Before(authFilter)) // 活跃人数
	mux.Handle("/v1/stats/online", nex.Handler(onlineLiteHandler).Before(authFilter))               // 同时在线人、桌数
	mux.Handle("/v1/stats/retention", nex.Handler(retentionHandler).Before(authFilter))             // 留存
	mux.Handle("/v1/stats/consume", nex.Handler(cardConsumeStatsHandler).Before(authFilter))        // 游戏币消耗

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(webDir))))
	mux.Handle("/ping", nex.Handler(pongHandler))

	// 监控指标
	mux.Handle("/6fab1940/metrics", promhttp.Handler())

	if logLevel := viper.GetString("core.log_level"); logLevel == "debug" {
		mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}

	return algoutil.AccessControl(algoutil.OptionControl(mux))
}

func Startup() {
	// setup database
	closer := dbStartup()
	defer closer()

	// enable white list
	enableWhiteList()

	redis_factory.Init()
	// provider.Alipay.Init()
	provider.Byrtypay.Init()
	// provider.WeChatV2Init()

	var (
		addr      = viper.GetString("webserver.addr")
		cert      = viper.GetString("webserver.certificates.cert")
		key       = viper.GetString("webserver.certificates.key")
		enableSSL = viper.GetBool("webserver.enable_ssl")
	)

	err := setupWebEventHandlers()
	if err != nil {
		logger.Panicf("Failed to setup event handlers: %v", err)
	}

	logger.Infof("Web service addr: %s(enable ssl: %v)", addr, enableSSL)
	go func() {
		// http service
		// todo mysql 同步redis
		// _, err := api.RankMysql()
		// if err != nil {
		// 	log.Errorln("mysql redis ", err)
		// }

		// 24小时清一次令牌桶
		scheduler.NewTimer(time.Hour*24, func() {
			api.ClearLimitRate()
		})

		mux := startupService()
		if enableSSL {
			log.Fatal(http.ListenAndServeTLS(addr, cert, key, mux))
		} else {
			log.Fatal(http.ListenAndServe(addr, mux))
		}
	}()

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	// stop server
	select {
	case s := <-sg:
		log.Infof("got signal: %s", s.String())
	}
}

func setupWebEventHandlers() error {
	event := dispatcher.GetInstance()

	subscriptions := []struct {
		topic   string
		handler interface{}
		async   bool
	}{
		{dispatcher.TopicCreateCharacter.String(), CreateCharacterHandler, true},
	}

	var errors []error
	for _, sub := range subscriptions {
		var err error
		if sub.async {
			err = event.SubscribeAsync(sub.topic, sub.handler, false)
		} else {
			err = event.Subscribe(sub.topic, sub.handler)
		}

		if err != nil {
			errors = append(errors, fmt.Errorf("failed to subscribe to %s: %w", sub.topic, err))
		}
	}

	if len(errors) > 0 {
		var errMsg string
		for i, err := range errors {
			if i > 0 {
				errMsg += "; "
			}
			errMsg += err.Error()
		}
		return fmt.Errorf("event subscription errors: %s", errMsg)
	}

	return nil
}
