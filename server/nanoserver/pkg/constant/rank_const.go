package constant

import (
	"fmt"
)

const (
	//Redis键前缀
	ServerKeyPrefix = "%d"
	//排行key
	ServerRankTypeKey = ServerKeyPrefix + "rank:{s%d}"
	//玩家key
	ServerRankUserKey = "user:%d"

	//世界boss的跨服榜   RankTypeWorldBossKf=21
	ServerRankWorldBossKey = "server:%d"
	//一战到底每一层榜单  RankTypeGightEndlayers=23
	ServerRankFightEndKey = "{" + ServerKeyPrefix + "fightend}:%d"
)

// 排行key
func RedisRankKey(serverID int32, rankType int32) string {
	// rankKey := fmt.Sprintf(ServerRankTypeKey, serverID, rankType)
	rankKey := fmt.Sprintf(ServerKeyPrefix+"rank:%d", serverID, rankType)
	return rankKey
}

// 排行玩家key
func RedisRankUserKey(userID int64) string {
	rankUserKey := fmt.Sprintf(ServerRankUserKey, userID)
	return rankUserKey
}
func RedisRankUserKeyByServer(serverID int32, userID int64) string {
	rankUserKey := fmt.Sprintf(ServerKeyPrefix+"user:%d", serverID, userID)
	return rankUserKey
}

// 世界boss的跨服榜
func RedisRankWorldBossKey(serverID int64) string {
	worldBossKey := fmt.Sprintf(ServerRankWorldBossKey, serverID)
	return worldBossKey
}

// 一战到底每一层榜单
func RedisRankFightEndKey(id int32, storey int32) string {
	fightEndKey := fmt.Sprintf(ServerRankFightEndKey, id, storey)

	return fightEndKey
}
