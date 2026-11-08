package thinkingdatax

import (
	"time"

	"github.com/ThinkingDataAnalytics/go-sdk/v2/src/thinkingdata"
)

// Option 函数类型，用于配置选项
type Option func(*options)

// 内部配置选项
type options struct {
	logConfig   *thinkingdata.TDLogConsumerConfig
	batchConfig *thinkingdata.TDBatchConfig
	debugConfig *debugConfig
	superProps  map[string]interface{}
	logger      thinkingdata.TDLogger // ThinkingData 日志记录器

	// 异步 worker 配置
	batchSize    int           // 批量处理大小
	flushTimeout time.Duration // 刷新超时时间
	chanSize     int           // 通道缓冲大小
}

type debugConfig struct {
	serverUrl string
	appId     string
	deviceId  string
}

// Consumer 类型常量
const (
	ConsumerTypeLog   = "log"   // 日志消费者：写入本地文件，适用于开发调试、离线或备份场景
	ConsumerTypeBatch = "batch" // 批量消费者：批量发送到服务器，适用于生产环境高性能数据上报
	ConsumerTypeDebug = "debug" // 调试消费者：实时发送并返回校验结果，适用于测试和数据格式验证
)

// WithLogConsumer 配置日志消费者
// 适用于：开发调试、数据备份、离线环境
func WithLogConsumer(config thinkingdata.TDLogConsumerConfig) Option {
	return func(o *options) {
		o.logConfig = &config
	}
}

// WithBatchConsumer 配置批量消费者
// 适用于：生产环境，高性能批量上报
func WithBatchConsumer(config thinkingdata.TDBatchConfig) Option {
	return func(o *options) {
		o.batchConfig = &config
	}
}

// WithDebugConsumer 配置调试消费者
// 适用于：测试环境，实时验证数据格式
func WithDebugConsumer(serverUrl, appId, deviceId string) Option {
	return func(o *options) {
		o.debugConfig = &debugConfig{
			serverUrl: serverUrl,
			appId:     appId,
			deviceId:  deviceId,
		}
	}
}

// WithSuperProperties 配置超级属性（所有事件都会携带）
func WithSuperProperties(props map[string]interface{}) Option {
	return func(o *options) {
		o.superProps = props
	}
}

// WithAsyncConfig 配置异步 worker 参数
func WithAsyncConfig(batchSize int, flushTimeout time.Duration, chanSize int) Option {
	return func(o *options) {
		if batchSize > 0 {
			o.batchSize = batchSize
		}
		if flushTimeout > 0 {
			o.flushTimeout = flushTimeout
		}
		if chanSize > 0 {
			o.chanSize = chanSize
		}
	}
}

// WithBatchSize 配置批量处理大小
func WithBatchSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.batchSize = size
		}
	}
}

// WithFlushTimeout 配置刷新超时时间
func WithFlushTimeout(timeout time.Duration) Option {
	return func(o *options) {
		if timeout > 0 {
			o.flushTimeout = timeout
		}
	}
}

// WithChannelSize 配置通道缓冲大小
func WithChannelSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.chanSize = size
		}
	}
}

// WithTDLogger 配置 ThinkingData 日志记录器
func WithTDLogger(logger thinkingdata.TDLogger) Option {
	return func(o *options) {
		if logger != nil {
			o.logger = logger
		}
	}
}
