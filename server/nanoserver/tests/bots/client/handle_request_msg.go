package client

import (
	"nanoserver/tests/bots/internal/codec"
	"nanoserver/tests/bots/internal/message"
	"nanoserver/tests/bots/internal/packet"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// SendMsg sends a message to the server and uses route to find the default handler
func (c *Client) SendMsg(routeStr string, pbMsg proto.Message) error {
	protoMessage, err := proto.Marshal(pbMsg)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal proto message for route %s", routeStr)
	}

	routes, ok := message.GetDictionary()
	if !ok {
		return errors.Errorf("route dictionary not found, cannot send message for route %s", routeStr)
	}
	routeCode, exists := routes[routeStr]
	if !exists {
		return errors.Errorf("route %s not found in dictionary", routeStr)
	}

	c.lastMessageID++

	requestMessage := &message.Message{
		Type:  message.Request,
		ID:    c.lastMessageID,
		Route: routeStr,
		Data:  protoMessage,
	}

	if handler, exists := c.staticResponseHandlers[uint64(routeCode)]; exists {
		c.dynamicResponseChan <- &DynamicResponse{
			DynamicID: c.lastMessageID,
			StaticID:  uint64(routeCode),
			Handler:   handler,
		}
	}

	encodedMessage, err := requestMessage.Encode()
	if err != nil {
		return errors.Wrapf(err, "failed to encode message for route %s", routeStr)
	}

	encodedPacket, err := codec.Encode(packet.Data, encodedMessage)
	if err != nil {
		return errors.Wrapf(err, "failed to encode packet for route %s", routeStr)
	}

	if c.closed.Load() {
		return errors.Errorf("client %d is closed, cannot send message", c.uid)
	}

	log.Debug().Msgf("Sending message uid=%d, name=%s, route=%s, ID=%d, req=%v", c.uid, c.name, routeStr, c.lastMessageID, pbMsg)
	select {
	case c.sendChan <- encodedPacket:
	case <-c.ctx.Done():
	default:
		log.Warn().Msgf("send channel full, dropping message for route %s", routeStr)
	}
	return nil
}
