package game

import (
	"fmt"
	"math/rand"
	"time"

	"nanoserver/db"
	"nanoserver/dispatcher"
	"nanoserver/pkg/goredisx"
	"nanoserver/pkg/redis_factory"
	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/serialize/protobuf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version             = "" // 游戏版本
	forceUpdate         = false
	logger              = log.WithField("component", "game")
	mySchedulerInstance scheduler.LocalScheduler
)

// Startup 初始化游戏服务器
func Startup() {
	var (
		err error
	)
	rand.Seed(time.Now().Unix())
	version = viper.GetString("update.version")

	heartbeat := viper.GetInt("core.heartbeat")
	if heartbeat < 5 {
		heartbeat = 5
	}

	// mySchedulerInstance, err = antsx.NewLocalScheduler()
	// if err != nil {
	// 	log.Fatalf("Failed to create local scheduler: %v", err)
	// }

	// if localSched, ok := mySchedulerInstance.(*antsx.LocalScheduler); ok {
	// 	defer localSched.Release()
	// } else {
	// 	log.Fatalf("Failed to assert scheduler type")
	// }

	// 游戏币消耗配置
	forceUpdate = viper.GetBool("update.force")

	logger.Infof("当前游戏服务器版本: %s, 是否强制更新: %t, 当前心跳时间间隔: %d秒", version, forceUpdate, heartbeat)
	logger.Info("game service starup")

	// 初始化世界管理器
	worldMgr.InitWorld()

	// register game handler
	comps := &component.Components{}
	// comps.Register(
	// 	defaultManager,
	// 	component.WithSchedulerName("my-custom-scheduler"),
	// )
	comps.Register(defaultManager)
	comps.Register(roomManager)
	comps.Register(matchManager)
	comps.Register(worldBossMgr)
	comps.Register(guildManager)

	InitTimer()

	redis_factory.Init()

	// 初始化掉落控制
	NewDropControl(db.NewDropControlRepo(goredisx.Instance))
	DropControlInstance.Load()
	defer DropControlInstance.Close()

	err = setupGameEventHandlers()
	if err != nil {
		logger.Panicf("Failed to setup event handlers: %v", err)
	}

	// 加密管道
	// c := newCrypto()
	// pip := pipeline.New()
	// pip.Inbound().PushBack(c.inbound)
	// pip.Outbound().PushBack(c.outbound)

	c := NewGameEventProcess()
	pip := pipeline.New()
	pip.Inbound().PushBack(c.InBound)
	pip.Outbound().PushBack(c.OutBound)

	addr := fmt.Sprintf(":%d", viper.GetInt("game-server.listen_port"))
	// cert := viper.GetString("webserver.certificates.cert")
	// key := viper.GetString("webserver.certificates.key")
	nano.Listen(addr,
		nano.WithPipeline(pip),
		nano.WithHeartbeatInterval(time.Duration(heartbeat)*time.Second),
		nano.WithLogger(log.WithField("component", "nano")),
		nano.WithSerializer(protobuf.NewSerializer()),
		nano.WithComponents(comps),
		// nano.WithDebugMode(),
		nano.WithIsWebsocket(true),
		// nano.WithTSLConfig(cert, key),
		nano.WithDictionary(protocol.Routemap),
		// nano.WithHandshakeValidator(customHandshakeValidator),
	)
	logger.Info("save data befor service stop")

	// 等待所有异步事件处理完成
	event := dispatcher.GetInstance()
	event.WaitAsync()

	// 给大家推送强制下线
	packet := &pb.CompulsoryWithdrawalPush{
		Status: pb.CompulsoryWithdrawalPush_WithdrawalStatus(pb.CompulsoryWithdrawalPush_CloseServer),
	}
	// 数数关闭
	// shushumain.Stop()
	//defaultManager.group.Broadcast(protocol.RouteCompulsoryWithdrawalPush, packet)
	playerMgr.BroadcastToAllGroups(protocol.RouteCompulsoryWithdrawalPush, packet)

	defaultManager.GmCmd(protocol.GmCmd{Cmd: "shutdown"})

	logger.Info("save data finish, now stop server")
	db.MustClose()
}

func setupGameEventHandlers() error {
	event := dispatcher.GetInstance()

	subscriptions := []struct {
		topic   string
		handler interface{}
		async   bool
	}{
		// {dispatcher.TopicRechargeSuccess.String(), ActivityCumulativeRechargeHandler, false},
		// {dispatcher.TopicActivityPush.String(), ActivityPushHandler, true},
		{dispatcher.TopicPayRefreshPush.String(), PayRefreshPushHandler, false},
		{dispatcher.TopicUserShowRefresh.String(), UserShowRefreshHandler, true},
		{dispatcher.TopicPVEWarMapBossDieCount.String(), PVEWarMapBossDieCountHandler, true},
		{dispatcher.TopicActivatePurchasedSlot.String(), ActivatePurchasedSlotHandler, true},
		{dispatcher.TopicAdClick.String(), AdClickHandler, true},
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

// 自定义握手验证器示例
// func customHandshakeValidator(s *session.Session, data []byte) error {
// 	s.Set("my-custom-scheduler", mySchedulerInstance)
// 	return nil
// }
