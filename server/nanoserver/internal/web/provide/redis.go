package provide

import (
	"nanoserver/pkg/errutil"
	"nanoserver/pkg/redis_factory"
)

// 获取有序数组前几名
func ZRevRangeCache(detail ...string) (interface{}, error) {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return nil, errutil.ErrDBOperation
	}
	res, err := tokenCacheRedisFact.ZRevRangeToRedis(detail)
	if err != nil {
		return nil, err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return res, nil
}

// 获取有序数组某个人的名次从大到小
func ZRevRankCache(detail ...string) (interface{}, error) {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return nil, errutil.ErrDBOperation
	}
	res, err := tokenCacheRedisFact.ZRevRankToRedis(detail)
	if err != nil {
		return nil, err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return res, nil
}

// 插入有序数组
func ZAddCache(key string, score int64, members int64) error {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return errutil.ErrDBOperation
	}
	_, err := tokenCacheRedisFact.ZAddToRedis(key, score, members)
	if err != nil {
		return err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return nil

}

func ZincrbyCache(key string, score int64, members int64) error {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return errutil.ErrDBOperation
	}
	_, err := tokenCacheRedisFact.ZincrbyToRedis(key, score, members)
	if err != nil {
		return err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return nil

}

// 获取哈希所有数据
func HGetAllCache(key string) (interface{}, error) {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return nil, errutil.ErrDBOperation
	}
	res, err := tokenCacheRedisFact.HGetAllToRedis(key)
	if err != nil {
		return nil, err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return res, nil
}

// 插入哈希
func HSetCache(detail ...string) error {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return errutil.ErrDBOperation
	}
	_, err := tokenCacheRedisFact.HSetToRedis(detail)
	if err != nil {
		return err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return nil
}

func CmdCache(cmd string, detail ...string) (interface{}, error) {
	tokenCacheRedisFact := CreatetFactory()
	if tokenCacheRedisFact == nil {
		return nil, errutil.ErrDBOperation
	}
	res, err := tokenCacheRedisFact.CmdToRedis(cmd, detail)
	if err != nil {
		return nil, err
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()
	return res, nil
}
func CreatetFactory() *userTokenCacheRedis {
	redCli := redis_factory.GetOneRedisClient()
	if redCli == nil {
		return nil
	}
	return &userTokenCacheRedis{redisClient: redCli, value: "1"}
}

// ReleaseRedisConn 释放redis
func (u *userTokenCacheRedis) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}

type userTokenCacheRedis struct {
	redisClient *redis_factory.RedisClient
	value       string
}

// 获取有序集合
func (u *userTokenCacheRedis) ZRevRangeToRedis(detail []string) (interface{}, error) {
	//key1 := []string{"field1", "value1", "field2", "value2"}
	// 将 []string 转换为 []interface{}
	var args []interface{}
	for _, v := range detail {
		args = append(args, v)
	}
	//conn.Do("hdel", redis.Args{}.Add(keyName).AddFlat(delStr)
	if res, err := u.redisClient.Execute("zrevrange", args...); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

// 获取有序集合某人
func (u *userTokenCacheRedis) ZRevRankToRedis(detail []string) (interface{}, error) {
	//key1 := []string{"field1", "value1", "field2", "value2"}
	// 将 []string 转换为 []interface{}
	var args []interface{}
	for _, v := range detail {
		args = append(args, v)
	}
	//conn.Do("hdel", redis.Args{}.Add(keyName).AddFlat(delStr)
	if res, err := u.redisClient.Execute("zrevrank", "rank:1", "4"); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) ZAddToRedis(key string, score int64, members int64) (bool, error) {
	if _, err := u.redisClient.Execute("zadd", key, score, members); err == nil {
		return true, nil
	} else {
		return false, err
	}
}
func (u *userTokenCacheRedis) ZincrbyToRedis(key string, score int64, members int64) (bool, error) {
	if _, err := u.redisClient.Execute("zincrby", key, score, members); err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) HGetAllToRedis(key string) (interface{}, error) {
	if res, err := u.redisClient.Execute("hgetall", key); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) HSetToRedis(detail []string) (bool, error) {
	//key1 := []string{"field1", "value1", "field2", "value2"}
	// 将 []string 转换为 []interface{}
	var args []interface{}
	for _, v := range detail {
		args = append(args, v)
	}
	//conn.Do("hdel", redis.Args{}.Add(keyName).AddFlat(delStr)
	if _, err := u.redisClient.Execute("hset", args...); err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// SetTokenCache 设置缓存
func (u *userTokenCacheRedis) CmdToRedis(cmd string, detail []string) (interface{}, error) {
	//key1 := []string{"field1", "value1", "field2", "value2"}
	// 将 []string 转换为 []interface{}
	var args []interface{}
	for _, v := range detail {
		args = append(args, v)
	}
	//conn.Do("hdel", redis.Args{}.Add(keyName).AddFlat(delStr)
	if res, err := u.redisClient.Execute(cmd, args...); err == nil {
		return res, nil
	} else {
		return false, err
	}
}
