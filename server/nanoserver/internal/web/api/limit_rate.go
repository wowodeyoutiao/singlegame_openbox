package api

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

var visitorIPs = make(map[string]*rate.Limiter)
var mu sync.Mutex

func limitHandler(next http.Handler, route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			// 当请求头不存在即不存在代理时直接获取ip
			ip = strings.Split(r.RemoteAddr, ":")[0]
			if ip == "" {
				logger.Infof("limitHandler: %s", route)
				// http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusForbidden)
				return
			}
		}
		ip = ip + route

		mu.Lock()
		limiter, exists := visitorIPs[ip]
		if !exists {
			limiter = rate.NewLimiter(2, 1) //每秒产生一个令牌，桶容量 1. 即每秒允许一个请求通过
			visitorIPs[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			logger.Infof("limitHandler: %s", ip)
			// http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func ClearLimitRate() {
	logger.Debugf("清理令牌桶...")
	if mu.TryLock() {
		defer mu.Unlock()
		for k := range visitorIPs {
			delete(visitorIPs, k)
		}
	}
}
