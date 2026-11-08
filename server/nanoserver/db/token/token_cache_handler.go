package token

import (
	"log"
)

func ValidTokenCacheToRedis(userId int64, token string, expires int64) {
	tokenCacheRedisFact := CreateUsersTokenCacheFactory(userId)
	if tokenCacheRedisFact == nil {
		log.Println("redis连接失败，请检查配置")
		return
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()

	tokenCacheRedisFact.SetTokenCache(expires, token)

	// 缓存结束之后删除超过系统设置最大在线数量的token
	tokenCacheRedisFact.DelOverMaxOnlineCache()
}

// DelTokenCacheFromRedis 删除redis所有的token
func DelTokenCacheFromRedis(userId int64) {
	tokenCacheRedisFact := CreateUsersTokenCacheFactory(userId)
	if tokenCacheRedisFact == nil {
		log.Println("redis连接失败，请检查配置")
		return
	}
	tokenCacheRedisFact.ClearUserToken()
	tokenCacheRedisFact.ReleaseRedisConn()
}
