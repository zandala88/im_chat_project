package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"im/repo"
	"im/util"
	"net/http"
	"sync"
)

// 添加一个读写锁
var rwLocker sync.RWMutex

// Node
// 构建 userId 和 Node 的映射关系
type Node struct {
	Conn *websocket.Conn
	// 并行转串行
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node)

func Chat(c *gin.Context) {
	userId := util.GetUid(c)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("用户%d 升级为 websocket 失败: %v", userId, err)
	}

	// 获取 Node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 1),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 将用户的群组信息添加到node 的 group 中
	communities, err := repo.GetCommunitiesByUserId(int(userId))
	if err != nil {
		zap.S().Errorf("GetCommunitiesByUserId err = %v", err)
	}

	for _, community := range communities {
		node.GroupSets.Add(community.CommunityId)
	}

	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// todo 发送逻辑
	go sendProc(node)

	// todo 完成接受逻辑
	go recvProc(node)

	// todo 发送逻辑
	sendMsg(userId, []byte("Hello, world!"))
}

// 发送逻辑
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			zap.S().Debugf("sendProc data = %v", data)
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Errorf("sendProc.WriteMessage err = %v", err)
				return
			}
		}
	}
}

// 接受逻辑
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Errorf("recvProc.ReadMessage err = %v", err)
			return
		}
		// todo 对 data 进一步处理
		dispatch(data)
	}
}

// 对 node 通道写入数据
func sendMsg(userId int64, data []byte) {
	// 获取 node
	rwLocker.Lock()
	node, ok := clientMap[userId]
	rwLocker.Unlock()

	if ok {
		node.DataQueue <- data
	}
}

// 对接收到的数据进一步处理
func dispatch(data []byte) {
	// 协程保存数据
	message := repo.Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		zap.S().Errorf("json.Unmarshal err=%v", err)
		return
	}
	zap.S().Debugf("收到了数据 %v", message)
	//保存数据库
	go repo.SaveMessage(message)

	// 对消息分类, 私聊或者是群聊
	switch message.Cmd {
	case repo.CMD_PRIVATE:
		// 发送消息给接受者
		go sendMsg(int64(message.ToId), data)
	case repo.CMD_GROUP:
		// 群聊, 查找群里面的所有用户发送数据
		for _, v := range clientMap {
			if v.GroupSets.Has(message.ToId) {
				v.DataQueue <- data
			}
		}
	}
}
