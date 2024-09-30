package server

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"im/config"
	"im/internal/kafka"
	"im/internal/websocket/constant"
	"im/internal/websocket/protocol"
)

type Client struct {
	Conn *websocket.Conn
	Name string
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		MyServer.Ungister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		messageType, message, err := c.Conn.ReadMessage()
		zap.S().Debug("messageType = ", messageType)
		if err != nil {
			zap.S().Error("client read message error", zap.String("client read message error", err.Error()))
			MyServer.Ungister <- c
			c.Conn.Close()
			break
		}

		msg := &protocol.Message{}
		proto.Unmarshal(message, msg)

		// pong
		if msg.Type == constant.HEAT_BEAT {
			pong := &protocol.Message{
				Content: constant.PONG,
				Type:    constant.HEAT_BEAT,
			}
			pongByte, err2 := proto.Marshal(pong)
			if err2 != nil {
				zap.S().Error("client marshal message error", zap.String("client marshal message error", err2.Error()))
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, pongByte)
		} else {
			if config.Configs.ChannelType.ChannelType == constant.KAFKA {
				kafka.Send(message)
			} else {
				zap.S().Debugf("MyServer.Broadcast <- message : %#v", msg)
				MyServer.Broadcast <- message
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		zap.S().Debug("c.Send message = ", message)
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
