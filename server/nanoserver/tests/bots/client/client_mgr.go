package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nanoserver/protocol"
	"nanoserver/protocol/protobuf/pb"
	"nanoserver/tests/bots/internal/message"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type WebLoginRequest struct{}

type ClientManager struct {
	clients map[int64]*Client
	mu      sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[int64]*Client),
	}
}

func (cm *ClientManager) StartClientsSequential(rangeStart, rangeEnd int) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	message.SetDictionary(protocol.Routemap)

	log.Info().Msgf("正在启动 %d 到 %d 的客户端...", rangeStart, rangeEnd)
	count := rangeEnd - rangeStart + 1
	for i := 0; i < count; i++ {
		currentBotID := rangeStart + i
		clientName := fmt.Sprintf("bot%05d", currentBotID)

		c := New()

		loginResp, err := c.Login(clientName)
		if err != nil {
			return err
		}
		c.uid = loginResp.Uid
		c.name = clientName

		// WebSocket连接
		if err := c.WebSocketConnect(); err != nil {
			return err
		}

		loginToGameServerRequest := &pb.LoginToGameServerRequest{
			Name:  c.name,
			Uid:   loginResp.Uid,
			Token: loginResp.Token,
		}
		if err := c.SendMsg(protocol.RouteManagerLogin, loginToGameServerRequest); err != nil {
			return err
		}

		cm.clients[c.uid] = c

		c.TestLogic()
	}

	log.Info().Msgf("已启动 %d 个客户端", len(cm.clients))

	return nil
}

func (cm *ClientManager) StartClientsConcurrent(rangeStart, rangeEnd int) error {
	log.Info().Msgf("正在启动 %d 到 %d 的客户端...", rangeStart, rangeEnd)

	message.SetDictionary(protocol.Routemap)

	g, ctx := errgroup.WithContext(context.Background())
	count := rangeEnd - rangeStart + 1
	clientChan := make(chan *Client, count)

	// 并发启动客户端
	for BotID := range count {
		BotID := rangeStart + BotID
		g.Go(func() error {
			currentBotID := BotID
			clientName := fmt.Sprintf("bot%05d", currentBotID)

			c := New()

			loginResp, err := c.Login(clientName)
			if err != nil {
				log.Error().Msgf("客户端 %s 登录失败: %v", clientName, err)
				return err
			}

			log.Info().Msgf("客户端 %s 登录成功, UID: %d", clientName, loginResp.Uid)

			if err := c.WebSocketConnect(); err != nil {
				return err
			}

			c.uid = loginResp.Uid
			c.name = clientName

			select {
			case clientChan <- c:
			case <-ctx.Done():
				return ctx.Err()
			}

			loginToGameServerRequest := &pb.LoginToGameServerRequest{
				Name:  c.name,
				Uid:   loginResp.Uid,
				Token: loginResp.Token,
			}
			if err := c.SendMsg(protocol.RouteManagerLogin, loginToGameServerRequest); err != nil {
				return err
			}

			c.TestLogic()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		close(clientChan)
		for c := range clientChan {
			c.Close(CloseWrite)
		}
		return err
	}
	close(clientChan)

	cm.mu.Lock()
	defer cm.mu.Unlock()

	for c := range clientChan {
		cm.clients[c.uid] = c
	}

	log.Info().Msgf("已成功启动 %d 个客户端", len(cm.clients))
	return nil
}

func (cm *ClientManager) StopAllClients() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Info().Msg("正在停止所有客户端...")

	if len(cm.clients) == 0 {
		log.Info().Msg("没有客户端需要停止")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	g, gctx := errgroup.WithContext(ctx)

	for uid, client := range cm.clients {
		client := client
		uid := uid
		g.Go(func() error {
			select {
			case <-gctx.Done():
				log.Warn().Msgf("客户端 UID:%d 关闭超时", uid)
				return gctx.Err()
			default:
				client.Close(CloseWrite)
				log.Debug().Msgf("客户端 UID:%d 已关闭", uid)
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		log.Warn().Msgf("部分客户端关闭时出现错误: %v", err)
	}

	cm.clients = make(map[int64]*Client)
	log.Info().Msg("所有客户端已停止")
}

// StopClient 停止指定的客户端
func (cm *ClientManager) StopClient(uid int64) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[uid]
	if !exists {
		return fmt.Errorf("客户端 UID:%d 不存在", uid)
	}

	client.Close(CloseWrite)
	delete(cm.clients, uid)
	log.Info().Msgf("客户端 UID:%d 已停止", uid)
	return nil
}

func (cm *ClientManager) GetClientCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.clients)
}

// GetClient 获取指定UID的客户端
func (cm *ClientManager) GetClient(uid int64) (*Client, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	client, exists := cm.clients[uid]
	return client, exists
}

// GetAllClients 获取所有客户端的副本（用于只读操作）
func (cm *ClientManager) GetAllClients() map[int64]*Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clients := make(map[int64]*Client, len(cm.clients))
	for uid, client := range cm.clients {
		clients[uid] = client
	}
	return clients
}
