package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"nanoserver/bootstrap"
	"nanoserver/global"
	"nanoserver/internal/game"
	shushumain "nanoserver/internal/shushu"
	"nanoserver/internal/web"
	"nanoserver/pkg/goredisx"
	"nanoserver/pkg/jwtv5x"
	"nanoserver/pkg/snowflake"
	"nanoserver/pkg/tracking/thinkingdatax"

	thinkingdataOfficial "github.com/ThinkingDataAnalytics/go-sdk/v2/src/thinkingdata"
	_ "github.com/go-sql-driver/mysql"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	// base application info
	app.Name = "mow server"
	app.Author = "MOW"
	app.Version = "0.0.1"
	app.Copyright = "mow team reserved"
	app.Usage = "mow server"

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./configs/config.toml",
			Usage: "load configuration from `FILE`",
		},
		cli.BoolFlag{
			Name:  "cpuprofile",
			Usage: "enable cpu profile",
		},
	}

	app.Action = serve
	app.Run(os.Args)
}

func serve(c *cli.Context) error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(c.String("config"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	logLevel := viper.GetString("core.log_level")
	logOutput := viper.GetString("core.log_output")
	if err := setupLogging(logLevel, logOutput); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	if c.Bool("cpuprofile") {
		filename := fmt.Sprintf("cpuprofile-%d.pprof", time.Now().Unix())
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	global.ServerIDRangeStart = viper.GetInt32("core.server_id_range_start")
	global.ServerIDRangeEnd = viper.GetInt32("core.server_id_range_end")
	snowflake.Init(time.DateOnly, int64(global.ServerIDRangeStart))

	bootstrap.Initialize()

	if err := global.NewServerList(); err != nil {
		log.Fatalf("Failed to initialize server list: %v", err)
	}
	defer global.ServerListMgrInstance.Close()

	redisAddr := fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port"))
	if err := goredisx.NewClient(
		goredisx.WithAddr(redisAddr),
		goredisx.WithAuth(viper.GetString("redis.auth"), viper.GetString("redis.password")),
		goredisx.WithDB(viper.GetInt("redis.index_db")),
	); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	tracker, err := thinkingdatax.New(
		thinkingdatax.ConsumerTypeBatch,
		thinkingdatax.WithBatchConsumer(thinkingdataOfficial.TDBatchConfig{
			ServerUrl: viper.GetString("shushu.shushuUrl"),
			AppId:     viper.GetString("shushu.shushuAppId"),
			AutoFlush: viper.GetBool("shushu.shushuAutoFlush"),
			BatchSize: viper.GetInt("shushu.shushuBatchSize"),
			Interval:  viper.GetInt("shushu.shushuInterval"),
		}),
	)
	if err != nil {
		log.Fatalf("Failed to initialize ThinkingData tracker: %v", err)
	}
	thinkingdatax.TrackerIns = tracker
	defer thinkingdatax.TrackerIns.Close()

	// 创建 JWT 令牌管理器
	jwtv5x.NewTokenManager(jwtv5x.Config{
		Secret:        viper.GetString("token.jwtTokenSignKey"),
		TokenDuration: time.Duration(3600*24*7) * time.Second,
		Repo:          jwtv5x.NewTokenRepo(goredisx.Instance),
	})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() { defer wg.Done(); game.Startup() }()       // 开启游戏服
	go func() { defer wg.Done(); web.Startup() }()        // 开启web服务器
	go func() { defer wg.Done(); shushumain.Startup() }() // 启动数据上报
	wg.Wait()
	return nil
}

// setupLogging 根据配置设置日志输出
func setupLogging(logLevel, logOutput string) error {
	// 设置默认值
	if logLevel == "" {
		logLevel = "info"
	}
	if logOutput == "" {
		logOutput = "file"
	}

	// 根据配置决定输出位置和格式
	switch logOutput {
	case "stdout", "console":
		// 控制台输出使用文本格式
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
			DisableColors:   false,
		})
		log.SetOutput(os.Stdout)
	case "file":
		fallthrough
	default:
		// 文件输出使用JSON格式
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		// 默认输出到文件
		if err := setupFileOutput(); err != nil {
			return err
		}
	}

	// 设置日志级别
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel) // 默认为Info级别
	}

	log.SetReportCaller(true)
	return nil
}

// setupFileOutput 配置文件输出
func setupFileOutput() error {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	writer, err := rotatelogs.New(
		"logs/app_%Y%m%d.log",
		rotatelogs.WithRotationTime(24*time.Hour), // 每天轮转
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留7天
		rotatelogs.WithRotationCount(0),           // 不限制文件数量，按时间清理
	)
	if err != nil {
		return fmt.Errorf("failed to create log rotator: %w", err)
	}

	log.SetOutput(writer)
	return nil
}
