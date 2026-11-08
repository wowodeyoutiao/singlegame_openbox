package goredisx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Instance redis.UniversalClient
)

// Config 封装 Redis 客户端的配置参数
type Config struct {
	// 基础连接配置
	addr     string // Redis 服务器地址
	username string // Redis 用户名（Redis 6.0+ ACL）
	password string // Redis 密码
	db       int    // Redis 数据库编号

	// 连接池配置
	poolSize     int // 连接池大小
	minIdleConns int // 最小空闲连接数

	// 超时配置
	dialTimeout  time.Duration // 连接超时
	readTimeout  time.Duration // 读取超时
	writeTimeout time.Duration // 写入超时

	// 安全配置
	tlsConfig *tls.Config // TLS 配置
}

// Option 定义配置选项函数类型
type Option func(*Config)

// WithAddr 设置 Redis 服务器地址
func WithAddr(addr string) Option {
	return func(c *Config) {
		c.addr = addr
	}
}

// WithAuth 设置用户名和密码
func WithAuth(username, password string) Option {
	return func(c *Config) {
		c.username = username
		c.password = password
	}
}

// WithDB 设置数据库编号
func WithDB(db int) Option {
	return func(c *Config) {
		c.db = db
	}
}

// WithPool 设置连接池
func WithPool(poolSize, minIdleConns int) Option {
	return func(c *Config) {
		c.poolSize = poolSize
		c.minIdleConns = minIdleConns
	}
}

// WithTimeouts 设置超时时间
func WithTimeouts(dial, read, write time.Duration) Option {
	return func(c *Config) {
		c.dialTimeout = dial
		c.readTimeout = read
		c.writeTimeout = write
	}
}

// WithTLS 设置 TLS 配置
func WithTLS(config *tls.Config) Option {
	return func(c *Config) {
		c.tlsConfig = config
	}
}

func NewConfig() *Config {
	poolSize := runtime.NumCPU() * 20

	return &Config{
		db:           0,
		poolSize:     poolSize,
		minIdleConns: poolSize / 5,
		dialTimeout:  2 * time.Second,
		readTimeout:  1 * time.Second,
		writeTimeout: 1 * time.Second,
	}
}

// NewClient 创建新的 Redis 客户端
func NewClient(opts ...Option) error {
	config := NewConfig()

	for _, opt := range opts {
		opt(config)
	}

	if config.addr == "" {
		return errors.New("redis address is required")
	}

	options := &redis.Options{
		Addr:         config.addr,
		Username:     config.username,
		Password:     config.password,
		DB:           config.db,
		PoolSize:     config.poolSize,
		MinIdleConns: config.minIdleConns,
		DialTimeout:  config.dialTimeout,
		ReadTimeout:  config.readTimeout,
		WriteTimeout: config.writeTimeout,
		TLSConfig:    config.tlsConfig,
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), config.dialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	Instance = client

	return nil
}
