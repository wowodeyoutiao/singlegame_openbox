package provide

import (
	"context"
	"fmt"
	"strconv"

	"nanoserver/pkg/constant"
	"nanoserver/pkg/goredisx"
	"nanoserver/protocol"

	"github.com/redis/go-redis/v9"
)

// hmsetUserFields writes multiple field/value pairs to the user's redis hash in one call
func hmsetUserFields(serverId int32, uid int64, fieldValues map[string]string) error {
	if len(fieldValues) == 0 {
		return nil
	}
	key := constant.RedisRankUserKeyByServer(serverId, uid)
	ctx := context.TODO()
	// Convert map[string]string to map[string]interface{} for HMSet
	data := make(map[string]interface{}, len(fieldValues))
	for k, v := range fieldValues {
		data[k] = v
	}
	return goredisx.Instance.HMSet(ctx, key, data).Err()
}

// hmsetWorldBossUserName writes only the name field into the special WorldBoss user key space
func hmsetWorldBossUserName(uid int64, name string, serverId int32) error {
	//暂时屏蔽，下次版本更新再换
	// key := constant.RedisRankWorldBossKey(uid)
	key := constant.RedisRankUserKeyByServer(serverId, uid)
	ctx := context.TODO()
	data := map[string]interface{}{"name": name}
	return goredisx.Instance.HMSet(ctx, key, data).Err()
}

// hgetAllUserFields loads all fields from the user's redis hash and converts to string map
func hgetAllUserFields(serverId int32, uid int64) (map[string]string, error) {
	key := constant.RedisRankUserKeyByServer(serverId, uid)
	ctx := context.TODO()
	return goredisx.Instance.HGetAll(ctx, key).Result()
}

// zrevrangeWithScores returns members and scores for a given range, as RankUserInfo with Uid/Value/Num filled
func zrevrangeWithScores(serverId int32, rankType int32, start string, stop string) ([]protocol.RankUserInfo, error) {
	key := constant.RedisRankKey(serverId, rankType)
	// convert start/stop
	s, _ := strconv.ParseInt(start, 10, 64)
	e, _ := strconv.ParseInt(stop, 10, 64)
	ctx := context.TODO()
	zs, err := goredisx.Instance.ZRevRangeWithScores(ctx, key, s, e).Result()
	if err != nil {
		return nil, err
	}
	list := make([]protocol.RankUserInfo, 0, len(zs))
	for i, z := range zs {
		// member to int64
		memberStr := fmt.Sprintf("%v", z.Member)
		member, err := strconv.ParseInt(memberStr, 10, 64)
		if err != nil {
			continue
		}
		list = append(list, protocol.RankUserInfo{
			Uid:   member,
			Num:   int64(i + 1),
			Value: int64(z.Score),
		})
	}
	return list, nil
}

// buildStandardUserFields constructs the common user fields map for HMSET
func buildStandardUserFields(name, avatar, sex, job, robotCid, serverId string, avatarV2 string, frame string) map[string]string {
	return map[string]string{
		"name":      name,
		"avatar":    avatar,
		"sex":       sex,
		"job":       job,
		"robot_cid": robotCid,
		"server_id": serverId,
		"avatar_v2": avatarV2,
		"frame":     frame,
	}
}

// ---- Thin wrappers to avoid scattering raw CmdCache usages ----

// zrevrangeTop1Raw returns raw member/score pairs for top 1
func zrevrangeTop1Raw(serverId int32, rankType int32) ([]interface{}, error) {
	key := constant.RedisRankKey(serverId, rankType)
	ctx := context.TODO()
	zs, err := goredisx.Instance.ZRevRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, 0, len(zs)*2)
	for _, z := range zs {
		memberStr := fmt.Sprintf("%v", z.Member)
		member, _ := strconv.ParseInt(memberStr, 10, 64)
		out = append(out, member, int64(z.Score))
	}
	return out, nil
}

// zrevrangeAtRaw returns raw member/score pairs for [start, stop]
func zrevrangeAtRaw(serverId int32, rankType int32, start string, stop string) ([]interface{}, error) {
	key := constant.RedisRankKey(serverId, rankType)
	s, _ := strconv.ParseInt(start, 10, 64)
	e, _ := strconv.ParseInt(stop, 10, 64)
	ctx := context.TODO()
	zs, err := goredisx.Instance.ZRevRangeWithScores(ctx, key, s, e).Result()
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, 0, len(zs)*2)
	for _, z := range zs {
		memberStr := fmt.Sprintf("%v", z.Member)
		member, _ := strconv.ParseInt(memberStr, 10, 64)
		out = append(out, member, int64(z.Score))
	}
	return out, nil
}

// zrevrankRaw fetches rank (0-based) for member
func zrevrankRaw(serverId int32, rankType int32, member int64) (int64, error) {
	key := constant.RedisRankKey(serverId, rankType)
	ctx := context.TODO()
	return goredisx.Instance.ZRevRank(ctx, key, strconv.FormatInt(member, 10)).Result()
}

// zscoreRaw fetches zset score for member
func zscoreRaw(serverId int32, rankType int32, member int64) (int64, error) {
	key := constant.RedisRankKey(serverId, rankType)
	ctx := context.TODO()
	v, err := goredisx.Instance.ZScore(ctx, key, strconv.FormatInt(member, 10)).Result()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// delRankKey deletes full rank key
func delRankKey(serverId int32, rankType int32) error {
	key := constant.RedisRankKey(serverId, rankType)
	ctx := context.TODO()
	return goredisx.Instance.Del(ctx, key).Err()
}

// hgetallUserRaw returns raw HGETALL values for a user
func hgetallUserRaw(serverId int32, uid int64) ([]interface{}, error) {
	key := constant.RedisRankUserKeyByServer(serverId, uid)
	ctx := context.TODO()
	m, err := goredisx.Instance.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		out = append(out, k, v)
	}
	return out, nil
}

// getStringKey returns GET key as string
func getStringKey(key string) (string, error) {
	ctx := context.TODO()
	return goredisx.Instance.Get(ctx, key).Result()
}

// setStringKey sets key to value
func setStringKey(key string, value string) error {
	ctx := context.TODO()
	return goredisx.Instance.Set(ctx, key, value, 0).Err()
}

// llenKey returns LLEN for list key
func llenKey(key string) (int64, error) {
	ctx := context.TODO()
	return goredisx.Instance.LLen(ctx, key).Result()
}

// zincrbyRank increments the score of member in the sorted set by increment
func zincrbyRank(key string, increment float64, member int64) error {
	ctx := context.TODO()
	_, err := goredisx.Instance.ZIncrBy(ctx, key, increment, strconv.FormatInt(member, 10)).Result()
	return err
}

// ZaddRank adds or updates a member with score in the sorted set
func ZaddRank(key string, score float64, member int64) error {
	ctx := context.TODO()
	_, err := goredisx.Instance.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: strconv.FormatInt(member, 10),
	}).Result()
	return err
}
