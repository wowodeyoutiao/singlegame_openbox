package thinkingdatax

import (
	"context"
	"log"
	"nanoserver/pkg/tracking"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ThinkingDataAnalytics/go-sdk/v2/src/thinkingdata"
)

var TrackerIns tracking.Tracker

// MessageType 消息类型
type MessageType int

const (
	MessageTypeEvent MessageType = iota
	MessageTypeUserSet
	MessageTypeUserSetOnce
	MessageTypeUserAdd
	MessageTypeFlush
)

// Message 异步消息结构
type Message struct {
	Type       MessageType
	AccountID  string
	DistinctID string
	Event      string
	Props      tracking.M
}

// TDTracker ThinkingData 异步实现
type TDTracker struct {
	ta      thinkingdata.TDAnalytics
	msgChan chan *Message
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	once    *sync.Once
	closed  atomic.Bool // 原子布尔标志位，标记是否已关闭

	// 配置参数
	batchSize    int           // 批量处理大小
	flushTimeout time.Duration // 刷新超时时间
	chanSize     int           // 通道缓冲大小
}

// New 创建埋点实例
// consumerType: 消费者类型（ConsumerTypeLog/ConsumerTypeBatch/ConsumerTypeDebug）
// opts: 配置选项
func New(consumerType string, opts ...Option) (tracking.Tracker, error) {
	// 默认配置
	defaultOpts := &options{
		logConfig: &thinkingdata.TDLogConsumerConfig{
			FileNamePrefix: "events",
			Directory:      "./logs",
		},
		// 异步配置默认值
		batchSize:    50,
		flushTimeout: 10 * time.Second,
		chanSize:     1000,
	}

	// 应用选项
	for _, opt := range opts {
		opt(defaultOpts)
	}

	// 创建消费者
	consumer, err := createConsumer(consumerType, defaultOpts)
	if err != nil {
		return nil, err
	}

	// 创建 analytics 实例
	te := thinkingdata.New(consumer)

	// 设置超级属性
	if defaultOpts.superProps != nil {
		te.SetSuperProperties(defaultOpts.superProps)
	}

	// 如果配置了日志记录器，设置到 ThinkingData 实例
	thinkingdata.SetCustomLogger(defaultOpts.logger)

	// 创建异步 tracker
	ctx, cancel := context.WithCancel(context.Background())
	tracker := &TDTracker{
		ta:           te,
		msgChan:      make(chan *Message, defaultOpts.chanSize),
		ctx:          ctx,
		cancel:       cancel,
		wg:           &sync.WaitGroup{},
		once:         &sync.Once{},
		batchSize:    defaultOpts.batchSize,
		flushTimeout: defaultOpts.flushTimeout,
		chanSize:     defaultOpts.chanSize,
	}

	// 启动后台 worker
	tracker.wg.Add(1)
	go tracker.worker()

	return tracker, nil
}

// 创建消费者
func createConsumer(consumerType string, opts *options) (thinkingdata.TDConsumer, error) {
	switch consumerType {
	case ConsumerTypeLog:
		return thinkingdata.NewLogConsumerWithConfig(*opts.logConfig)
	case ConsumerTypeBatch:
		if opts.batchConfig == nil {
			return nil, &ConsumerConfigError{Type: consumerType, Message: "batchConfig is required"}
		}
		return thinkingdata.NewBatchConsumerWithConfig(*opts.batchConfig)
	case ConsumerTypeDebug:
		if opts.debugConfig == nil {
			return nil, &ConsumerConfigError{Type: consumerType, Message: "debugConfig is required"}
		}
		return thinkingdata.NewDebugConsumerWithDeviceId(
			opts.debugConfig.serverUrl,
			opts.debugConfig.appId,
			false,
			opts.debugConfig.deviceId,
		)
	default:
		return thinkingdata.NewLogConsumerWithConfig(*opts.logConfig)
	}
}

// ConsumerConfigError 消费者配置错误
type ConsumerConfigError struct {
	Type    string
	Message string
}

func (e *ConsumerConfigError) Error() string {
	return "consumer config error [type: " + e.Type + "]: " + e.Message
}

// worker 后台处理消息的工作协程
func (t *TDTracker) worker() {
	defer t.wg.Done()

	// 批量处理缓冲区
	messageBatch := make([]*Message, 0, t.batchSize)

	flushTicker := time.NewTicker(t.flushTimeout)
	defer flushTicker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			// 处理剩余消息
			if len(messageBatch) > 0 {
				t.processBatch(messageBatch)
			}
			return

		case msg, ok := <-t.msgChan:
			if !ok {
				// 通道已关闭，退出
				if len(messageBatch) > 0 {
					t.processBatch(messageBatch)
				}
				return
			}

			// 收集消息到批次中
			messageBatch = append(messageBatch, msg)

			// 如果达到批次大小，立即处理
			if len(messageBatch) >= t.batchSize {
				t.processBatch(messageBatch)
				messageBatch = messageBatch[:0] // 重置切片
			}

		case <-flushTicker.C:
			// 定时处理（无论是否有消息都会触发）
			if len(messageBatch) > 0 {
				t.processBatch(messageBatch)
				messageBatch = messageBatch[:0] // 重置切片
			}
		}
	}
}

// processBatch 批量处理消息
func (t *TDTracker) processBatch(messages []*Message) {
	for _, msg := range messages {
		var err error

		switch msg.Type {
		case MessageTypeEvent:
			err = t.ta.Track(msg.AccountID, msg.DistinctID, msg.Event, msg.Props)
		case MessageTypeUserSet:
			err = t.ta.UserSet(msg.AccountID, msg.DistinctID, msg.Props)
		case MessageTypeUserSetOnce:
			err = t.ta.UserSetOnce(msg.AccountID, msg.DistinctID, msg.Props)
		case MessageTypeUserAdd:
			err = t.ta.UserAdd(msg.AccountID, msg.DistinctID, msg.Props)
		case MessageTypeFlush:
			// 只有显式的 Flush 消息才调用 SDK 的 Flush
			err = t.ta.Flush()
		}

		// 如果有错误，直接使用 helper 打印日志
		if err != nil {
			// t.helper.Errorf("ThinkingData operation failed: type=%d, accountID=%s, distinctID=%s, event=%s, error=%v",
			// 	msg.Type, msg.AccountID, msg.DistinctID, msg.Event, err)
			log.Printf("ThinkingData operation failed: type=%d, accountID=%s, distinctID=%s, event=%s, error=%v",
				msg.Type, msg.AccountID, msg.DistinctID, msg.Event, err)
		}
	}
}

// sendMessage 发送消息到工作协程（异步）
func (t *TDTracker) sendMessage(msg *Message) error {
	// 检查是否已关闭
	if t.closed.Load() {
		return &AsyncError{Message: "tracker is closed"}
	}

	select {
	case t.msgChan <- msg:
		return nil
	case <-t.ctx.Done():
		return context.Canceled
	default:
		// 通道满了，返回错误
		return &AsyncError{Message: "message channel is full"}
	}
}

// sendMessageSync 发送消息到工作协程（同步）
func (t *TDTracker) sendMessageSync(msg *Message) error {
	// 对于 Flush 操作，直接在当前协程中执行
	if msg.Type == MessageTypeFlush {
		return t.ta.Flush()
	}

	// 其他操作转为异步处理
	return t.sendMessage(msg)
}

// AsyncError 异步错误
type AsyncError struct {
	Message string
}

func (e *AsyncError) Error() string {
	return "async error: " + e.Message
}

// E 事件上报（异步）
// accountID: 用户账号ID（已登录时传入，如手机号、用户ID）
// distinctID: 设备ID（必传，如设备UUID、匿名ID）
// event: 事件名称（如 "login"、"payment"）
// props: 事件属性（如 M{"amount": 99.99, "channel": "wechat"}）
func (t *TDTracker) E(accountID, distinctID, event string, props tracking.M) error {
	if props == nil {
		props = make(tracking.M)
	}
	// 自动添加时间戳（如果未设置）
	if _, exists := props["#time"]; !exists {
		props["#time"] = time.Now()
	}

	msg := &Message{
		Type:       MessageTypeEvent,
		AccountID:  accountID,
		DistinctID: distinctID,
		Event:      event,
		Props:      props,
	}

	return t.sendMessage(msg)
}

// Set 设置用户属性（可覆盖）（异步）
// 适用于：用户基本信息更新（昵称、等级等）
func (t *TDTracker) Set(accountID, distinctID string, props tracking.M) error {
	msg := &Message{
		Type:       MessageTypeUserSet,
		AccountID:  accountID,
		DistinctID: distinctID,
		Props:      props,
	}

	return t.sendMessage(msg)
}

// SetOnce 设置用户属性（只设置一次）（异步）
// 适用于：用户首次行为（注册时间、首次来源等）
func (t *TDTracker) SetOnce(accountID, distinctID string, props tracking.M) error {
	msg := &Message{
		Type:       MessageTypeUserSetOnce,
		AccountID:  accountID,
		DistinctID: distinctID,
		Props:      props,
	}

	return t.sendMessage(msg)
}

// Add 累加用户属性（异步）
// 适用于：数值型属性累加（付费总额、登录次数等）
func (t *TDTracker) Add(accountID, distinctID string, props tracking.M) error {
	msg := &Message{
		Type:       MessageTypeUserAdd,
		AccountID:  accountID,
		DistinctID: distinctID,
		Props:      props,
	}

	return t.sendMessage(msg)
}

// Flush 刷新缓冲区（同步）
func (t *TDTracker) Flush() error {
	msg := &Message{
		Type: MessageTypeFlush,
	}

	return t.sendMessageSync(msg)
}

// Close 关闭实例，确保完整的生命周期管理
func (t *TDTracker) Close() error {
	t.once.Do(func() {
		// 设置关闭标志，防止新消息发送
		t.closed.Store(true)

		// 取消上下文，通知 worker 停止
		t.cancel()

		// 等待 worker 处理完成
		t.wg.Wait()

		// 关闭通道
		close(t.msgChan)

		// 关闭底层实例
		t.ta.Close()
	})

	return nil
}
