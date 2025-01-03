package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
	"im/public/protocol"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ResendCountMax = 3 // 超时重传最大次数
)

type Client struct {
	conn                 *websocket.Conn
	token                string
	userId               int64
	clientId             int64
	clientId2Cancel      map[int64]context.CancelFunc
	clientId2CancelMutex sync.Mutex
	seq                  int64
	sendCh               chan []byte
}

func NewClient(userId, token, host string) *Client {
	c := &Client{
		clientId2Cancel: make(map[int64]context.CancelFunc),
		token:           token,
		userId:          cast.ToInt64(userId),
		sendCh:          make(chan []byte, 1024),
	}

	// 连接 websocket
	conn, _, err := websocket.DefaultDialer.Dial(host+"/ws", http.Header{})
	if err != nil {
		panic(err)
	}
	c.conn = conn
	// 向 websocket 发送登录请求
	c.login()
	go c.heartbeat()
	go c.write()
	go c.read()
	return c
}

func (c *Client) read() {
	for {
		_, bytes, err := c.conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		outputBatchMsg := new(protocol.OutputBatch)
		err = proto.Unmarshal(bytes, outputBatchMsg)
		if err != nil {
			panic(err)
		}
		for _, output := range outputBatchMsg.Outputs {
			msg := new(protocol.Output)
			err = proto.Unmarshal(output, msg)
			if err != nil {
				panic(err)
			}

			// 只收两种，Message 收取下行消息和 ACK，上行消息ACK回复
			switch msg.Type {
			case protocol.CmdType_CT_Message:
				// 计算接收消息数量
				atomic.AddInt64(&receivedMessageCount, 1)
				msgTimer.updateEndTime()

				pushMsg := new(protocol.PushMsg)
				err = proto.Unmarshal(msg.Data, pushMsg)
				if err != nil {
					panic(err)
				}
				// 更新 seq
				seq := pushMsg.Msg.Seq
				if c.seq < seq {
					c.seq = seq
				}
			case protocol.CmdType_CT_ACK: // 收到 ACK
				ackMsg := new(protocol.ACKMsg)
				err = proto.Unmarshal(msg.Data, ackMsg)
				if err != nil {
					panic(err)
				}

				switch ackMsg.Type {
				case protocol.ACKType_AT_Up: // 收到上行消息的 ACK
					// 计算接收消息数量
					atomic.AddInt64(&receivedMessageCount, 1)
					msgTimer.updateEndTime()

					// 取消超时重传
					clientId := ackMsg.ClientId
					c.clientId2CancelMutex.Lock()
					if cancel, ok := c.clientId2Cancel[clientId]; ok {
						// 取消超时重传
						cancel()
						delete(c.clientId2Cancel, clientId)
					}
					c.clientId2CancelMutex.Unlock()
					// 更新客户端本地维护的 seq
					seq := ackMsg.Seq
					if c.seq < seq {
						c.seq = seq
					}
				}
			default:
				fmt.Println("未知消息类型")
			}
		}
	}

}

func (c *Client) write() {
	for {
		select {
		case bytes, ok := <-c.sendCh:
			if !ok {
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
				return
			}
		}
	}
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(time.Minute * 2)
	for range ticker.C {
		c.sendMsg(protocol.CmdType_CT_Heartbeat, &protocol.HeartbeatMsg{})
	}
}

func (c *Client) login() {
	c.sendMsg(protocol.CmdType_CT_Login, &protocol.LoginMsg{
		Token: []byte(c.token),
	})
}

// send 发送消息，启动超时重试
func (c *Client) send(chatId int64) {
	message := &protocol.Message{
		SessionType: protocol.SessionType_ST_Group,              // 群聊
		ReceiverId:  chatId,                                     // 发送到该群
		SenderId:    c.userId,                                   // 发送者
		MessageType: protocol.MessageType_MT_Text,               // 文本
		Content:     []byte("文本聊天消息" + cast.ToString(c.userId)), // 消息
		SendTime:    time.Now().UnixMilli(),                     // 发送时间
	}
	UpMsg := &protocol.UpMsg{
		Msg:      message,
		ClientId: c.getClientId(),
	}
	// 发送消息
	c.sendMsg(protocol.CmdType_CT_Message, UpMsg)

	// 启动超时重传
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		maxRetry := ResendCountMax // 最大重试次数
		retryCount := 0
		retryInterval := time.Millisecond * 500 // 重试间隔
		ticker := time.NewTicker(retryInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if retryCount >= maxRetry {
					return
				}
				c.sendMsg(protocol.CmdType_CT_Message, UpMsg)
				retryCount++
			}
		}
	}(ctx)

	c.clientId2CancelMutex.Lock()
	c.clientId2Cancel[UpMsg.ClientId] = cancel
	c.clientId2CancelMutex.Unlock()

}

func (c *Client) getClientId() int64 {
	c.clientId++
	return c.clientId
}

// 客户端向服务端发送上行消息
func (c *Client) sendMsg(cmdType protocol.CmdType, msg proto.Message) {
	// 组装顶层数据
	cmdMsg := &protocol.Input{
		Type: cmdType,
		Data: nil,
	}
	if msg != nil {
		data, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		cmdMsg.Data = data
	}

	bytes, err := proto.Marshal(cmdMsg)
	if err != nil {
		panic(err)
	}

	// 发送
	c.sendCh <- bytes
}
