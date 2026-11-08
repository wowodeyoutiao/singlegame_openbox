package web

import (
	"context"
	"nanoserver/pkg/my_jwt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimitError 频率限制错误
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// NewRateLimitError 创建频率限制错误
func NewRateLimitError(message string) *RateLimitError {
	return &RateLimitError{
		Message: message,
	}
}

// 用户级别锁相关
var (
	userLocks              sync.Map           // 用户ID -> *userLockInfo
	defaultUserLockTimeout = 10 * time.Second // 默认用户锁超时时间
)

func ShouldLockUrl(url string) bool {
	_, exists := ValidRouterEndpoints[url]
	return exists // 在白名单才需要加锁
}

// RouterWhitelist 存储允许的 路由
var ValidRouterEndpoints = map[string]struct{}{
	"/v1/prize_draw/pet_receive":     {},
	"/v1/store/buy":                  {},
	"/v1/prize_draw/partner_receive": {},
}

// userLockInfo 用户锁信息
type userLockInfo struct {
	mu       sync.Mutex
	locked   bool
	timer    *time.Timer
	lockedAt time.Time
}

// userLockMiddleware 用户级别锁中间件
func userLockMiddleware(ctx context.Context, r *http.Request) (context.Context, error) {
	if !ShouldLockUrl(r.RequestURI) {
		return ctx, nil // 白名单路由不加锁
	}

	userInfo, ok := ctx.Value("userDate").(my_jwt.CustomClaims)
	if !ok {
		return ctx, nil // 白名单路由直接通过
	}

	userID := userInfo.UserId
	if userID <= 0 {
		return ctx, nil // 无效用户ID直接通过
	}

	key := r.Method + ":" + r.RequestURI + ":" + strconv.FormatInt(userID, 10)

	// 获取或创建锁信息
	lockInfo, _ := userLocks.LoadOrStore(key, &userLockInfo{})
	info := lockInfo.(*userLockInfo)

	info.mu.Lock()
	defer info.mu.Unlock()

	if info.locked {
		logger.Warnf("[UserLock] User %d is locked, rejecting request: %s %s", userID, r.Method, r.RequestURI)
		return ctx, NewRateLimitError("请求过于频繁，请稍后再试")
	}

	// 设置锁定状态
	info.locked = true
	info.lockedAt = time.Now()

	// 将key存储到context中，供解锁时使用
	ctx = context.WithValue(ctx, "lockKey", key)

	// 停止之前的定时器并设置新的超时解锁
	if info.timer != nil {
		info.timer.Stop()
	}

	info.timer = time.AfterFunc(defaultUserLockTimeout, func() {
		unlockUserByKey(key)
	})

	logger.Debugf("[UserLock] User %d locked for request: %s %s", userID, r.Method, r.RequestURI)
	return ctx, nil
}

// userUnLockMiddleware 用户级别解锁中间件
func userUnLockMiddleware(ctx context.Context, w http.ResponseWriter) (context.Context, error) {
	key, ok := ctx.Value("lockKey").(string)
	if !ok {
		return ctx, nil
	}

	unlockUserByKey(key)
	return ctx, nil
}

// unlockUserByKey 解锁指定用户（提取公共逻辑）
func unlockUserByKey(key string) {
	if lockInfo, ok := userLocks.Load(key); ok {
		info := lockInfo.(*userLockInfo)
		info.mu.Lock()
		defer info.mu.Unlock()

		if info.locked {
			info.locked = false
			if info.timer != nil {
				info.timer.Stop()
				info.timer = nil
			}
			logger.Debugf("[UserLock] Key %s unlocked", key)
		}
	}
}
