package db

import (
	"context"
	"nanoserver/db/model"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DropControlPrefix = "dropctl:"
)

type DropControlRepo struct {
	rdb redis.UniversalClient
}

func NewDropControlRepo(rdb redis.UniversalClient) *DropControlRepo {
	return &DropControlRepo{rdb: rdb}
}

func (r *DropControlRepo) Load() ([]*model.ItemDrop, error) {
	ctx := context.Background()

	// 使用 Scan 获取所有 dropctl: 前缀的 key，例如 dropctl:1:2592656071600455680:63:10
	var allKeys []string
	iter := r.rdb.Scan(ctx, 0, DropControlPrefix+"*", 1000).Iterator()
	for iter.Next(ctx) {
		allKeys = append(allKeys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	if len(allKeys) == 0 {
		return []*model.ItemDrop{}, nil
	}

	// 批量获取所有 key 的值
	results, err := r.rdb.MGet(ctx, allKeys...).Result()
	if err != nil {
		return nil, err
	}

	var drops []*model.ItemDrop
	for i, result := range results {
		if result == nil {
			continue
		}

		str, ok := result.(string)
		if !ok {
			continue
		}

		// 从 key 中提取 Type, InstanceId, ItemId, MapId
		// key 格式: dropctl:Type:InstanceId:ItemId:mapID
		keyParts := strings.Split(allKeys[i], ":")
		if len(keyParts) != 5 {
			continue
		}

		// 从 value 中提取 BucketStartTime, CurrentCount
		// value 格式: BucketStartTime:CurrentCount
		valueParts := strings.Split(str, ":")
		if len(valueParts) != 2 {
			continue
		}

		drop := &model.ItemDrop{}

		// 解析 key 中的字段
		if val, err := strconv.ParseInt(keyParts[1], 10, 32); err == nil {
			drop.Type = int32(val)
		} else {
			continue
		}

		if val, err := strconv.ParseInt(keyParts[2], 10, 64); err == nil {
			drop.InstanceId = val
		} else {
			continue
		}

		if val, err := strconv.ParseInt(keyParts[3], 10, 64); err == nil {
			drop.ItemId = int32(val)
		} else {
			continue
		}

		if val, err := strconv.ParseInt(keyParts[4], 10, 64); err == nil {
			drop.MapId = int64(val)
		} else {
			continue
		}

		// 解析 value 中的字段
		if val, err := strconv.ParseInt(valueParts[0], 10, 64); err == nil {
			drop.BucketStartTime = val
		} else {
			continue
		}

		if val, err := strconv.ParseInt(valueParts[1], 10, 32); err == nil {
			drop.CurrentCount = int32(val)
		} else {
			continue
		}

		drops = append(drops, drop)
	}

	return drops, nil
}

func (r *DropControlRepo) Save(drops []*model.ItemDrop) error {
	if len(drops) == 0 {
		return nil
	}

	ctx := context.Background()

	// 去重 drops, 只使用最后一条的记录，以前的记录全部丢弃
	uniqueDrops := make(map[string]*model.ItemDrop)
	for _, drop := range drops {
		// 构造 key: dropctl:Type:InstanceId:ItemId:mapID
		key := DropControlPrefix +
			strconv.FormatInt(int64(drop.Type), 10) + ":" +
			strconv.FormatInt(drop.InstanceId, 10) + ":" +
			strconv.FormatInt(int64(drop.ItemId), 10) + ":" +
			strconv.FormatInt(int64(drop.MapId), 10)
		uniqueDrops[key] = drop
	}

	// 创建 pipeline
	pipe := r.rdb.Pipeline()

	// 计算明天0点的过期时间
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	for key, drop := range uniqueDrops {
		// 构造 value: BucketStartTime:CurrentCount
		value := strconv.FormatInt(drop.BucketStartTime, 10) + ":" +
			strconv.FormatInt(int64(drop.CurrentCount), 10)

		// 设置值
		pipe.Set(ctx, key, value, 0)
		// 设置过期时间为明天0点
		pipe.ExpireAt(ctx, key, tomorrow)
	}

	// 执行 pipeline
	_, err := pipe.Exec(ctx)
	return err
}
