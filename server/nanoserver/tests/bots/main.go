package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	slog "log"
	"nanoserver/bootstrap"
	"nanoserver/tests/bots/client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	viper.SetConfigType("toml")
	viper.SetConfigFile("./configs/robot.toml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Fatal("读取配置文件失败", err)
		return
	}

	// 设置日志
	logLevel := viper.GetString("core.log_level")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":

		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: "2006-01-02 15:04:05",
	}).With().Caller().Logger()

	// 加载业务配置
	_, err := os.Stat("./configs/gameconfig")
	if err != nil {
		log.Error().Msgf("配置文件不存在，请拷贝 gameconfig 至 ./configs/，错误: %v", err)
		return
	}
	bootstrap.Initialize()
	clientManager := client.NewClientManager()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().Msgf("收到信号: %v, 开始优雅关闭", sig)
		clientManager.StopAllClients()
		os.Exit(0)
	}()

	rangeStart := viper.GetInt("bot.range_start")
	rangeEnd := viper.GetInt("bot.range_end")
	loginMode := viper.GetString("bot.login_mode")
	startTime := time.Now()

	switch loginMode {
	case "concurrent":
		if err := clientManager.StartClientsConcurrent(rangeStart, rangeEnd); err != nil {
			log.Error().Msgf("启动并发客户端失败: %+v", err)
			return
		}
	case "sequential":
		if err := clientManager.StartClientsSequential(rangeStart, rangeEnd); err != nil {
			log.Error().Msgf("启动顺序客户端失败: %+v", err)
			return
		}
	default:
		log.Error().Msgf("未知登录模式: %s, 只支持 'concurrent' 或 'sequential'", loginMode)
		return
	}
	duration := time.Since(startTime)
	log.Info().Msgf("客户端启动完成，耗时: %v，当前运行的客户端数量: %d", duration, clientManager.GetClientCount())
	log.Info().Msg("按 Ctrl+C 来优雅关闭所有客户端")

	select {}
}
