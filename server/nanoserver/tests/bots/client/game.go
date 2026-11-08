package client

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"nanoserver/tests/bots/internal/codec"
	"nanoserver/tests/bots/internal/message"
	"nanoserver/tests/bots/internal/packet"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Client struct {
	// lock          sync.Mutex
	uid           int64
	name          string
	conn          *websocket.Conn
	sendChan      chan []byte
	heartbeatChan chan []byte
	lastMessageID uint64 // 上次发送的消息ID

	pushHandlers           map[string]func([]byte)
	dynamicResponseChan    chan *DynamicResponse   // 动态响应通道
	staticResponseHandlers map[uint64]func([]byte) // 静态响应处理函数
	codec                  *codec.Decoder          // decoder

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool
}

type DynamicResponse struct {
	StaticID  uint64
	DynamicID uint64
	Handler   func([]byte)
}

func New() *Client {
	pctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		ctx:                 pctx,
		cancel:              cancel,
		sendChan:            make(chan []byte, 16),
		heartbeatChan:       make(chan []byte, 1),
		codec:               codec.NewDecoder(),
		dynamicResponseChan: make(chan *DynamicResponse, 1),
	}
	c.initPushHandlers()
	c.initResponseHandlers()
	return c
}

// WebSocket连接
func (c *Client) WebSocketConnect() error {
	wsURL := viper.GetString("server.ws_url")
	if wsURL == "" {
		return errors.Errorf("WebSocket URL is not configured")
	}

	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to WebSocket server at %s", wsURL)
	}

	if err := c.handshake(); err != nil {
		c.conn.Close()
		return err
	}

	c.wg.Add(2)
	go c.readLoop()
	go c.writeLoop()
	return nil
}

func (c *Client) readLoop() {
	defer func() {
		log.Debug().Msg("readLoop stopped")
		c.Close(CloseRead)
		c.wg.Done()
	}()

	for {
		select {
		case <-c.ctx.Done():
			log.Info().Msg("readLoop context done, exiting")
			return
		default:
			if c.closed.Load() {
				return
			}

			_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			messageType, binMsg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Info().Msgf("WebSocket closed by server: %v", err)
					return
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Error().Msgf("read error: %v", err)
				return
			}

			if messageType == websocket.CloseMessage {
				log.Info().Msg("received CloseMessage from server")
				return
			}

			if messageType != websocket.BinaryMessage {
				log.Warn().Msgf("expected BinaryMessage, got %d, binMsg: %s", messageType, binMsg)
				continue
			}

			packets, err := c.codec.Decode(binMsg)
			if err != nil {
				log.Error().Msgf("PackageDecode failed: %v", err)
				continue
			}

			c.handlePackets(packets)
		}
	}
}

func (c *Client) handlePackets(packets []*packet.Packet) {
	for _, packetItem := range packets {
		switch packetItem.Type {
		case packet.Heartbeat:
			c.handleHeartbeat()
		case packet.Handshake:
			log.Debug().Msg("Received handshake")
		case packet.HandshakeAck:
			log.Debug().Msg("Received handshake ack")
		case packet.Data:
			c.handleData(packetItem.Data)
		case packet.Kick:
			log.Warn().Msg("Received kick packet, closing client")
		}
	}
}

func (c *Client) handleData(data []byte) {
	msg, err := message.Decode(data)
	if err != nil {
		log.Error().Msgf("MessageDecode failed: %v, data: %v", err, data)
		return
	}

	switch msg.Type {
	case message.Response:
		c.HandleResponseMsg(msg)
	case message.Push:
		c.HandlePushMsg(msg)
	default:
		log.Warn().Msgf("Unknown message type: %v", msg.Type)
	}
}

func (c *Client) handleHeartbeat() {
	log.Debug().Int64("uid", c.uid).Str("name", c.name).Msgf("Received heartbeat")
	b, err := codec.Encode(packet.Heartbeat, nil)
	if err != nil {
		log.Error().Msgf("encode heartbeat error: %v", err)
		return
	}

	if c.closed.Load() {
		return
	}

	select {
	case c.heartbeatChan <- b:
	case <-c.ctx.Done():
	default:
		log.Warn().Msg("sendChan full, heartbeat dropped")
	}
}

func (c *Client) writeLoop() {
	defer func() {
		log.Debug().Msg("writeLoop stopped")
		c.Close(CloseWrite)
		c.wg.Done()
	}()

	for {
		select {
		case msg, ok := <-c.sendChan:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closed"))
				if err != nil {
					log.Error().Msgf("failed to send close message: %v", err)
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Error().Msgf("write error: %v", err)
				return
			}
		case msg, ok := <-c.heartbeatChan:
			if !ok {
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Error().Msgf("heartbeat write error: %v", err)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

type GameCloseType int

const (
	CloseRead  GameCloseType = iota // 服务端主动关闭
	CloseWrite                      // 客户端主动关闭
)

func (c *Client) Close(closeType GameCloseType) {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}

	log.Info().Msgf("Closing client %s, UID: %d", c.name, c.uid)
	if closeType == CloseRead {
		c.cancel()
		close(c.sendChan)
	} else {
		close(c.sendChan)
		c.cancel()
	}
	c.wg.Wait()

	if c.conn != nil {
		_ = c.conn.Close()
	}
}

// handshake performs the initial handshake with the server
func (c *Client) handshake() error {
	handshakeData := map[string]string{
		"version": "1.1.1",
		"type":    "js-websocket",
	}
	handshakeRaw, err := json.Marshal(handshakeData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal handshake data")
	}
	handshakeEncodedPacket, err := codec.Encode(packet.Handshake, handshakeRaw)
	if err != nil {
		return errors.Wrap(err, "failed to encode handshake packet")
	}
	if err := c.conn.WriteMessage(websocket.TextMessage, handshakeEncodedPacket); err != nil {
		return errors.Wrap(err, "failed to send handshake request")
	}
	// 等待服务端握手响应
	messageType, b, err := c.conn.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "failed to read handshake response")
	}
	if messageType != websocket.BinaryMessage {
		log.Warn().Msgf("expected binary message for handshake ack, got %d", messageType)
	}

	log.Debug().Msgf("received handshake response, message type: %d, data length: %d", messageType, len(b))

	handshakeEncodedPacket, err = codec.Encode(packet.HandshakeAck, nil)
	if err != nil {
		return errors.Wrap(err, "failed to encode handshake ack packet")
	}
	if err := c.conn.WriteMessage(websocket.TextMessage, handshakeEncodedPacket); err != nil {
		return errors.Wrap(err, "failed to send handshake ack")
	}
	return nil
}
