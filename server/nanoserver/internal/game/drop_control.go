package game

import (
	"context"
	"fmt"
	"math"
	"nanoserver/db/model"
	"nanoserver/global"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	DropControlInstance *DropControl
)

type DropControlScope int

const (
	DropControlScopeNone   DropControlScope = iota // 无效类型
	DropControlScopePlayer                         // 用户级别
	DropControlScopeServer                         // 服务器级别
	ItemDropTypeBattleZone                         // 战区级别
	DropControlScopeLimit                          // 最大类型数
)

// DropControlStore 定义了物品掉落的存储接口
type DropControlStore interface {
	Load() ([]*model.ItemDrop, error)
	Save([]*model.ItemDrop) error
}

// BucketDropLimiter 用于限制每个桶内的掉落次数
type BucketDropLimiter struct {
	bucketStartTime time.Time // 当前桶的开始时间
	currentCount    int32     // 当前桶的计数
	// bucketDuration  time.Duration // 桶的时间长度
}

// DropControl 处理物品掉落的核心结构
type DropControl struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	close  atomic.Bool                                                    // 是否已关闭
	ch     chan *model.ItemDrop                                           // 用于异步保存掉落记录的通道
	mu     sync.RWMutex                                                   // 保护 drops 的并发访问
	drops  [DropControlScopeLimit]map[int64]map[string]*BucketDropLimiter // DropControlScope -> 实例ID -> itemID:MapID -> BucketDropLimiter
	repo   DropControlStore                                               // 数据存储接口
}

// NewDropControl 创建一个新的 ItemDrop 实例
func NewDropControl(repo DropControlStore) {
	if repo == nil {
		log.Fatalf("ItemDropStore is nil, cannot create ItemDrop")
	}

	ctx, cancel := context.WithCancel(context.Background())
	id := &DropControl{
		ctx:    ctx,
		cancel: cancel,
		ch:     make(chan *model.ItemDrop, 100),
		repo:   repo,
	}

	for i := int(DropControlScopePlayer); i < int(DropControlScopeLimit); i++ {
		id.drops[i] = make(map[int64]map[string]*BucketDropLimiter)
	}

	id.wg.Add(1)
	go id.worker()

	DropControlInstance = id
}

// Close 关闭 ItemDrop，等待所有 goroutine 完成
func (id *DropControl) Close() {
	if id.close.Swap(true) {
		return
	}
	id.cancel()
	close(id.ch)
	id.wg.Wait()
}

// Load 从数据库加载掉落记录
func (id *DropControl) Load() {
	records, err := id.repo.Load()
	if err != nil {
		log.Errorf("failed to load item drops: %v", err)
		return
	}

	for _, record := range records {
		if record.Type < 0 || record.Type >= int32(DropControlScopeLimit) {
			continue
		}

		if _, exists := id.drops[DropControlScope(record.Type)][record.InstanceId]; !exists {
			id.drops[DropControlScope(record.Type)][record.InstanceId] = make(map[string]*BucketDropLimiter)
		}
		itemMapID := fmt.Sprintf("%d:%d", record.ItemId, record.MapId)
		id.drops[DropControlScope(record.Type)][record.InstanceId][itemMapID] = &BucketDropLimiter{
			bucketStartTime: time.Unix(record.BucketStartTime, 0),
			currentCount:    record.CurrentCount,
		}
	}
}

// worker 处理掉落记录的保存
func (id *DropControl) worker() {
	defer id.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	records := make([]*model.ItemDrop, 0, 100)

	for {
		select {
		case <-id.ctx.Done():
			if len(records) > 0 {
				if err := id.repo.Save(records); err != nil {
					log.Errorf("failed to save item drops: %v", err)
				}
				records = records[:0]
			}
			return
		case record, ok := <-id.ch:
			if !ok {
				return
			}
			records = append(records, record)
			if len(records) >= 100 {
				if err := id.repo.Save(records); err != nil {
					log.Errorf("failed to save item drops: %v", err)
				}
				records = records[:0]
			}
		case <-ticker.C:
			if len(records) > 0 {
				if err := id.repo.Save(records); err != nil {
					log.Errorf("failed to save item drops: %v", err)
				}
				records = records[:0]
			}
		}
	}
}

type DropOption func(*dropOptions)

type dropOptions struct {
	dropMultiplier float64
}

func defaultDropOptions() *dropOptions {
	return &dropOptions{
		dropMultiplier: 1.0,
	}
}

func WithDropMultiplier(multiplier float64) DropOption {
	return func(opts *dropOptions) {
		opts.dropMultiplier = multiplier
	}
}

// TryDrop 尝试掉落物品， dropMultiplier 用于调整掉落数量的倍数, 默认为 1.0
func (id *DropControl) TryDrop(itemID int32, itemNum int64, mapID int32, player *Player, opts ...DropOption) (isDrop bool) {
	options := defaultDropOptions()
	for _, opt := range opts {
		opt(options)
	}

	// 检查掉落控制配置
	key := fmt.Sprintf("%d:%d", itemID, mapID)
	cfg, exists := global.ItemDropControlCfgMap[key]
	if !exists {
		// log.Warnf("dropcontrol not found for ItemGroup: %d, Candropmap: %d", itemID, mapID)
		return true
	}

	// 检查玩家等级是否在控制范围内
	playerLevel := player.Level()
	if playerLevel < cfg.ControlLevel.Min || playerLevel > cfg.ControlLevel.Max {
		log.Warnf("player level %d out of range [%d, %d] for ItemGroup: %d, Candropmap: %d",
			playerLevel, cfg.ControlLevel.Min, cfg.ControlLevel.Max, itemID, mapID)
		return false
	}

	// 计算实际掉落数量
	dropNum := int32(math.Floor(float64(itemNum) * options.dropMultiplier))
	if dropNum <= 0 {
		log.Warnf("calculated dropNum <= 0 for ItemGroup: %d, Candropmap: %d, itemNum: %d, dropMultiplier: %f",
			itemID, mapID, itemNum, options.dropMultiplier)
		return false
	}

	// 线程安全的获取和检查限制器
	id.mu.Lock()
	defer id.mu.Unlock()

	limiter := id.getOrCreateLimiter(cfg.Scope, player, itemID, mapID)

	// 检查桶是否过期，需要重置
	elapsedSeconds := int64(time.Since(limiter.bucketStartTime).Seconds())
	if elapsedSeconds >= int64(cfg.ControlTime) {
		limiter.bucketStartTime = time.Now()
		limiter.currentCount = 0
	}

	// 检查是否超过限制
	if limiter.currentCount+dropNum > cfg.Maxnum {
		log.Warnf("drop limit reached for ItemGroup: %d, Candropmap: %d, currentCount: %d, add: %d, Maxnum: %d, dropMultiplier: %f",
			itemID, mapID, limiter.currentCount, dropNum, cfg.Maxnum, options.dropMultiplier)
		return false
	}

	// 增加计数
	limiter.currentCount += dropNum

	// 异步保存记录
	id.saveDropRecord(cfg.Scope, player, itemID, int64(mapID), limiter)

	return true
}

// getOrCreateLimiter 获取或创建限制器
func (id *DropControl) getOrCreateLimiter(scope int32, player *Player, itemID, mapID int32) *BucketDropLimiter {
	scopeType := DropControlScope(scope)
	var instanceID int64

	switch scopeType {
	case DropControlScopePlayer:
		instanceID = player.Uid()
	case DropControlScopeServer:
		instanceID = int64(player.ServerId())
	case ItemDropTypeBattleZone:
		instanceID = int64(player.getZoneId())
	default:
		log.Errorf("invalid drop control scope: %d", scope)
		return nil
	}

	// 确保 scope 层级存在
	if id.drops[scopeType] == nil {
		id.drops[scopeType] = make(map[int64]map[string]*BucketDropLimiter)
	}

	// 确保实例层级存在
	if id.drops[scopeType][instanceID] == nil {
		id.drops[scopeType][instanceID] = make(map[string]*BucketDropLimiter)
	}

	// 获取或创建限制器
	itemMapID := fmt.Sprintf("%d:%d", itemID, mapID)
	limiter := id.drops[scopeType][instanceID][itemMapID]
	if limiter == nil {
		limiter = &BucketDropLimiter{
			bucketStartTime: time.Now(),
			currentCount:    0,
		}
		id.drops[scopeType][instanceID][itemMapID] = limiter
	}

	return limiter
}

// saveDropRecord 异步保存掉落记录
func (id *DropControl) saveDropRecord(scope int32, player *Player, itemID int32, mapID int64, limiter *BucketDropLimiter) {
	var instanceID int64
	switch DropControlScope(scope) {
	case DropControlScopePlayer:
		instanceID = player.Uid()
	case DropControlScopeServer:
		instanceID = int64(player.ServerId())
	case ItemDropTypeBattleZone:
		instanceID = int64(player.getZoneId())
	default:
		return
	}

	record := &model.ItemDrop{
		Type:            scope,
		InstanceId:      instanceID,
		ItemId:          itemID,
		MapId:           mapID,
		BucketStartTime: limiter.bucketStartTime.Unix(),
		CurrentCount:    limiter.currentCount,
	}

	if id.close.Load() {
		return
	}
	select {
	case <-id.ctx.Done():
		return
	case id.ch <- record:
	default:
		log.Warnf("drop record channel is full, dropping record for item %d", itemID)
	}
}

func (id *DropControl) Reset() {
	id.mu.Lock()
	defer id.mu.Unlock()

	for i := int(DropControlScopePlayer); i < int(DropControlScopeLimit); i++ {
		id.drops[i] = make(map[int64]map[string]*BucketDropLimiter)
	}
}
