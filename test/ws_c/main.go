package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"im/public/protocol"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	httpAddr       = "http://localhost:9090"
	websocketAddr  = "ws://localhost:9091"
	ResendCountMax = 3 // 超时重传最大次数
)

type Client struct {
	conn                 *websocket.Conn
	token                string
	userId               int64
	clientId             int64
	clientId2Cancel      map[int64]context.CancelFunc // clientId 到 context 的映射
	clientId2CancelMutex sync.Mutex
	seq                  int64 // 本地消息最大同步序列号
}

// websocket 客户端
func main() {
	// http 登录，获取 token
	client := Login()

	// 连接 websocket 服务
	client.Start()
}

func (c *Client) Start() {
	// 连接 websocket
	conn, _, err := websocket.DefaultDialer.Dial(websocketAddr+"/ws", http.Header{})
	if err != nil {
		panic(err)
	}
	c.conn = conn

	fmt.Println("与 websocket 建立连接")
	c.Login()
	go c.Heartbeat()
	time.Sleep(time.Millisecond * 100)
	go c.Sync()
	go c.Receive()
	time.Sleep(time.Millisecond * 100)
	c.ReadLine()
}

// ReadLine 读取用户消息并发送
func (c *Client) ReadLine() {
	var (
		msg         string
		receiverId  int64
		sessionType int64
	)

	readLine := func(hint string, v interface{}) {
		fmt.Println(hint)
		_, err := fmt.Scanln(v)
		if err != nil {
			panic(err)
		}
	}
	for {
		readLine("发送消息", &msg)
		readLine("接收人id(user_id/group_id)：", &receiverId)
		readLine("发送消息类型(1-单聊、2-群聊)：", &sessionType)
		message := &protocol.Message{
			SessionType: protocol.SessionType(sessionType),
			ReceiverId:  receiverId,
			SenderId:    c.userId,
			MessageType: protocol.MessageType_MT_Text,
			Content:     []byte(msg),
			SendTime:    time.Now().UnixMilli(),
		}
		UpMsg := &protocol.UpMsg{
			Msg:      message,
			ClientId: c.GetClientId(),
		}

		c.SendMsg(protocol.CmdType_CT_Message, UpMsg)

		// 启动超时重传
		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context) {
			maxRetry := ResendCountMax // 最大重试次数
			retryCount := 0
			retryInterval := time.Millisecond * 100 // 重试间隔
			ticker := time.NewTicker(retryInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("收到 ACK，不再重试")
					return
				case <-ticker.C:
					if retryCount >= maxRetry {
						fmt.Println("达到最大超时次数，不再重试")
						return
					}
					fmt.Println("消息超时 msg:", msg, "，第", retryCount+1, "次重试")
					c.SendMsg(protocol.CmdType_CT_Message, UpMsg)
					retryCount++
				}
			}
		}(ctx)

		c.clientId2CancelMutex.Lock()
		c.clientId2Cancel[UpMsg.ClientId] = cancel
		c.clientId2CancelMutex.Unlock()

		time.Sleep(time.Second)
	}
}

func (c *Client) Heartbeat() {
	ticker := time.NewTicker(time.Second * 2 * 60)
	for range ticker.C {
		c.SendMsg(protocol.CmdType_CT_Heartbeat, &protocol.HeartbeatMsg{})
	}
}

func (c *Client) Sync() {
	c.SendMsg(protocol.CmdType_CT_Sync, &protocol.SyncInputMsg{Seq: c.seq})
}

func (c *Client) Receive() {
	for {
		_, bytes, err := c.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		c.HandlerMessage(bytes)
	}
}

// HandlerMessage 客户端消息处理
func (c *Client) HandlerMessage(bytes []byte) {
	outputBatchMsg := new(protocol.OutputBatch)
	err := proto.Unmarshal(bytes, outputBatchMsg)
	if err != nil {
		panic(err)
	}
	for _, output := range outputBatchMsg.Outputs {
		msg := new(protocol.Output)
		err := proto.Unmarshal(output, msg)
		if err != nil {
			panic(err)
		}

		fmt.Println("收到顶层 OutPut 消息：", msg)

		switch msg.Type {
		case protocol.CmdType_CT_Sync:
			syncMsg := new(protocol.SyncOutputMsg)
			err = proto.Unmarshal(msg.Data, syncMsg)
			if err != nil {
				panic(err)
			}

			seq := c.seq
			for _, message := range syncMsg.Messages {
				fmt.Println("收到离线消息：", message)
				if seq < message.Seq {
					seq = message.Seq
				}
			}
			c.seq = seq
			// 如果还有，继续拉取
			if syncMsg.HasMore {
				c.Sync()
			}
		case protocol.CmdType_CT_Message:
			pushMsg := new(protocol.PushMsg)
			err = proto.Unmarshal(msg.Data, pushMsg)
			if err != nil {
				panic(err)
			}
			fmt.Printf("收到消息：%s, 发送人id：%d, 会话类型：%s, 接收时间:%s\n", pushMsg.Msg.GetContent(), pushMsg.Msg.GetSenderId(), pushMsg.Msg.SessionType, time.Now().Format("2006-01-02 15:04:05"))
			// 更新 seq
			seq := pushMsg.Msg.Seq
			if c.seq < seq {
				c.seq = seq
			}
			fmt.Println("更新 seq:", c.seq)
		case protocol.CmdType_CT_ACK: // 收到 ACK
			ackMsg := new(protocol.ACKMsg)
			err = proto.Unmarshal(msg.Data, ackMsg)
			if err != nil {
				panic(err)
			}

			switch ackMsg.Type {
			case protocol.ACKType_AT_Login:
				fmt.Println("登录成功")
			case protocol.ACKType_AT_Up: // 收到上行消息的 ACK
				// 取消超时重传
				clientId := ackMsg.ClientId
				c.clientId2CancelMutex.Lock()
				if cancel, ok := c.clientId2Cancel[clientId]; ok {
					// 取消超时重传
					cancel()
					delete(c.clientId2Cancel, clientId)
					fmt.Println("取消超时重传，clientId:", clientId)
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

// Login websocket 登录
func (c *Client) Login() {
	fmt.Println("websocket login...")
	// 组装底层数据
	loginMsg := &protocol.LoginMsg{
		Token: []byte(c.token),
	}
	c.SendMsg(protocol.CmdType_CT_Login, loginMsg)
}

// SendMsg 客户端向服务端发送上行消息
func (c *Client) SendMsg(cmdType protocol.CmdType, msg proto.Message) {
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
	err = c.conn.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		panic(err)
	}
}

func (c *Client) GetClientId() int64 {
	c.clientId++
	return c.clientId
}

// Login 用户http登录获取 token
func Login() *Client {
	// 读取 phone_number 和 password 参数
	var phoneNumber, password string
	fmt.Print("Enter phone_number: ")
	fmt.Scanln(&phoneNumber)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	// 创建一个 url.Values 对象，并将 phone_number 和 password 参数添加到其中
	data := url.Values{}
	data.Set("phone_number", phoneNumber)
	data.Set("password", password)

	// 向服务器发送 POST 请求
	resp, err := http.PostForm(httpAddr+"/login", data)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		panic(err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected HTTP status code: %d\n", resp.StatusCode)
		panic(err)
	}

	// 读取返回数据
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 获取 token
	var respData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token  string `json:"token"`
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	err = json.Unmarshal(bodyData, &respData)
	if err != nil {
		panic(err)
	}

	if respData.Code != 200 {
		panic(fmt.Sprintf("登录失败, %s", respData))
	}
	// 获取客户端 seq 序列号
	var seq int64
	fmt.Print("Enter seq: ")
	fmt.Scanln(&seq)

	client := &Client{
		clientId2Cancel: make(map[int64]context.CancelFunc),
		seq:             seq,
	}

	client.token = respData.Data.Token
	clientStr := respData.Data.UserId
	client.userId, err = strconv.ParseInt(clientStr, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("client code:", respData.Code, " ", respData.Msg)
	fmt.Println("token:", client.token, "userId", client.userId)
	return client
}
