package controller

import (
	"encoding/json"
	"fmt"
	"im/model"
	"im/service"
	"im/util"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

// 添加一个读写锁
var rwLocker sync.RWMutex

// 构建 userId 和 Node 的映射关系
type Node struct {
	Conn *websocket.Conn
	// 并行转串行
	DataQueue chan []byte
	GroupSets set.Interface
}

//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

func Chat(ctx *gin.Context) {
	// todo 检测是否合法
	//checkToken()
	userId, _ := strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	//token := ctx.Query("token")
	// 建立链接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	// 获取 Node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 1),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 将用户的群组信息添加到node 的 group 中
	communities, err := service.GetCommunitiesByUserId(int(userId))
	if err != nil {
		log.Println(err)
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
			fmt.Println(data)
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
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
			log.Println(err.Error())
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
	message := model.Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("收到了数据 %s \n", data)
	//保存数据库
	go service.SaveMessage(message)

	// 对消息分类, 私聊或者是群聊
	switch message.Cmd {
	case model.CMD_PRIVATE:
		// 发送消息给接受者
		go sendMsg(int64(message.ToId), data)
	case model.CMD_GROUP:
		// 群聊, 查找群里面的所有用户发送数据
		for _, v := range clientMap {
			if v.GroupSets.Has(message.ToId) {
				v.DataQueue <- data
			}
		}
	}
}

// 获取好友的聊天消息
func GetFriendMessages(ctx *gin.Context) {
	friendId, ok := ctx.Params.Get("friendId")
	maxId := ctx.Query("maxId")
	if !ok {
		h := util.Response{
			Code:    -1,
			Message: "参数错误",
		}
		h.Fail(ctx.Writer)
		return
	}
	fId, err := strconv.Atoi(friendId)
	if err != nil {
		h := util.Response{
			Code:    -1,
			Message: "参数错误",
		}
		h.Fail(ctx.Writer)
		return
	}
	auth, _ := ctx.Get("auth")
	// 当前最大 id
	mId, err := strconv.Atoi(maxId)
	if err != nil || mId == 0 {
		mId = math.MaxUint32
	}
	user := auth.(model.User)
	messages := service.GetFriendMessages(int(user.ID), fId, mId)
	h := util.Response{
		Code: 0,
		Data: messages,
	}
	h.Success(ctx.Writer)
}
