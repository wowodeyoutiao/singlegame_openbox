package global

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger                = logrus.WithField("component", "global")
	ServerListMgrInstance *ServerListMgr
)

type ServerListResponse struct {
	ErrCode int32       `json:"code"`
	ErrMsg  string      `json:"error"`
	Data    ServerLists `json:"data"`
}

type ServerLists struct {
	Count int32        `json:"count"`
	List  []ServerInfo `json:"list"`
}

type ServerInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IPAddr    string `json:"ipaddr"`
	ServerID  int    `json:"server_id"`
	MergeSID  int32  `json:"merge_sid"`
	OpenBy    int64  `json:"open_by"` // 时间戳（毫秒）
	Status    int    `json:"status"`
	IsShow    int    `json:"is_show"`
	CreatedAt string `json:"createdAt"` // ISO8601 时间格式
	UpdatedAt string `json:"updatedAt"` // ISO8601 时间格式
	CreateBy  int64  `json:"createBy"`
	UpdateBy  int64  `json:"updateBy"`
	ZoneId    int32  `json:"zone_id"` //战区Id
}

type ServerListMgr struct {
	ctx         context.Context
	cancel      context.CancelFunc
	httpClient  *resty.Client
	serverInfos map[int]*ServerInfo
	// mu              sync.RWMutex
	lastUpdate      time.Time
	refreshURL      string
	maxProdServerID int // 最大生产环境的服务器ID，用于分配新服务器ID
	wg              sync.WaitGroup
}

// NewServerList 创建并初始化 ServerListMgr
// 在初始化时必须成功拉取数据，否则返回错误
func NewServerList() error {
	adminURL := viper.GetString("admin.adminUrl")
	if adminURL == "" {
		return fmt.Errorf("admin.adminUrl is not configured")
	}

	ctx, cancel := context.WithCancel(context.Background())

	mgr := &ServerListMgr{
		ctx:         ctx,
		cancel:      cancel,
		httpClient:  resty.New().SetTimeout(10 * time.Second),
		serverInfos: make(map[int]*ServerInfo),
		refreshURL:  fmt.Sprintf("%s/api/v1/servername", adminURL),
	}

	// 首次同步拉取数据
	if err := mgr.fetchAndUpdate(); err != nil {
		cancel()
		return fmt.Errorf("failed to initialize server list: %w", err)
	}

	// 启动后台刷新任务
	mgr.wg.Add(1)
	go mgr.startAutoRefresh()

	ServerListMgrInstance = mgr
	return nil
}

// Close 优雅关闭 ServerListMgr
func (s *ServerListMgr) Close() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

// fetchAndUpdate 拉取数据并更新缓存
func (s *ServerListMgr) fetchAndUpdate() error {
	serverInfos, err := s.fetchServerList()
	if err != nil {
		return err
	}

	// s.mu.Lock()
	s.serverInfos = serverInfos
	s.lastUpdate = time.Now()
	// s.mu.Unlock()

	return nil
}

// startAutoRefresh 启动自动刷新任务
func (s *ServerListMgr) startAutoRefresh() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := s.fetchAndUpdate(); err != nil {
				logger.Errorf("failed to refresh server list: %v", err)
			}
		}
	}
}

// fetchServerList 从API拉取服务器列表
func (s *ServerListMgr) fetchServerList() (map[int]*ServerInfo, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	var resp ServerListResponse

	response, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&resp).
		Post(s.refreshURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("http request failed with status code: %d", response.StatusCode())
	}

	// if resp.ErrCode != 0 {
	// 	return nil, fmt.Errorf("api error: %s (code: %d)", resp.ErrMsg, resp.ErrCode)
	// }

	serverInfos := make(map[int]*ServerInfo, len(resp.Data.List))
	var maxServerID int
	for i := range resp.Data.List {
		info := resp.Data.List[i]
		serverInfos[info.ServerID] = &info
		// IsShow 为 1 表示生产环境的服务器
		if info.ServerID > maxServerID && info.IsShow == 1 {
			maxServerID = info.ServerID
		}
	}

	s.maxProdServerID = maxServerID

	return serverInfos, nil
}

func (s *ServerListMgr) GetMaxProdServerID() int {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	return s.maxProdServerID
}

// Get 非阻塞获取服务器信息
// 如果缓存中没有该 serverId，则返回错误
func (s *ServerListMgr) Get(serverId int) (ServerInfo, error) {
	// s.mu.RLock()
	info, exists := s.serverInfos[serverId]
	// s.mu.RUnlock()

	if !exists {
		return ServerInfo{}, fmt.Errorf("server with ID %d not found", serverId)
	}

	return *info, nil
}

// GetAll 获取所有服务器信息的副本
func (s *ServerListMgr) GetAll() map[int]ServerInfo {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	result := make(map[int]ServerInfo, len(s.serverInfos))
	for id, info := range s.serverInfos {
		result[id] = *info
	}
	return result
}

// GetAllZoneId 获取所有战区ID
func (s *ServerListMgr) GetAllZoneId() []int32 {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	zoneIdMap := make(map[int32]struct{}, len(s.serverInfos))
	for _, v := range s.serverInfos {
		zoneIdMap[v.ZoneId] = struct{}{}
	}
	// 直接提取 map 的 keys
	ZoneIds := make([]int32, 0, len(zoneIdMap))
	for id := range zoneIdMap {
		ZoneIds = append(ZoneIds, id)
	}
	return ZoneIds
}

// GetLastUpdateTime 获取最后更新时间
func (s *ServerListMgr) GetLastUpdateTime() time.Time {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	return s.lastUpdate
}

// 获取本服的所有战区ID
func (s *ServerListMgr) GetLocalServerZoneIds() []int32 {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	zoneIdMap := make(map[int32]struct{}, len(s.serverInfos))
	for i := ServerIDRangeStart; i <= ServerIDRangeEnd; i++ {
		for _, v := range s.serverInfos {
			if v.ServerID == int(i) {
				zoneIdMap[v.ZoneId] = struct{}{}
			}
		}
	}
	// 直接提取 map 的 keys
	zoneIds := make([]int32, 0, len(zoneIdMap))
	for id := range zoneIdMap {
		zoneIds = append(zoneIds, id)
	}
	return zoneIds
}

// 获取战区内所有服务器ID
func (s *ServerListMgr) GetZoneServerIds(zoneId int32) []int {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	serverIds := make([]int, 0, len(s.serverInfos))
	for _, v := range s.serverInfos {
		if v.ZoneId == zoneId {
			serverIds = append(serverIds, v.ServerID)
		}
	}
	return serverIds
}

func (s *ServerListMgr) GetZoneId(serverID int32) int32 {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	if serverInfo, err := s.Get(int(serverID)); err == nil {
		return serverInfo.ZoneId
	}
	return 0
}

// 获取服务器开启天数
func (s *ServerListMgr) GetServerOpenDay(serverID int32) int32 {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	if info, err := s.Get(int(serverID)); err == nil {
		openTime := time.UnixMilli(info.OpenBy)
		localOpenZero := time.Date(openTime.Year(), openTime.Month(), openTime.Day(),
			0, 0, 0, 0, time.Local)

		diffDays := int32(time.Since(localOpenZero).Hours()/24) + 1

		return diffDays
	}
	return 1
}

// 获取真实服务器serverID
func (s *ServerListMgr) GetRealServerID(serverId int32) int32 {
	if info, exists := s.serverInfos[int(serverId)]; exists {
		if info.MergeSID != 0 {
			return info.MergeSID
		}
	}
	return serverId
}

// 该工会是否合服到当前server
func (s *ServerListMgr) IsGuildMergeToServer(targetSID int32, sourceSID int32) bool {
	if targetSID != sourceSID {
		info := s.serverInfos[int(sourceSID)]
		if info == nil {
			return false
		}
		if info.MergeSID != targetSID {
			return false
		}
	}
	return true
}
