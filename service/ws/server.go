package ws

import (
	"go.uber.org/zap"
	"im/config"
	"sync"
)

var (
	ConnManager *Server
	once        sync.Once
)

// Server 连接管理
// 1. 连接管理
// 2. 工作队列
type Server struct {
	connMap   sync.Map    // 登录的用户连接 k-用户userid v-连接
	taskQueue []chan *Req // 工作池
}

func GetServer() *Server {
	once.Do(func() {
		ConnManager = &Server{
			taskQueue: make([]chan *Req, config.Configs.App.WorkerPoolSize), // 初始worker队列中，worker个数
		}
	})
	return ConnManager
}

// Stop 关闭服务
func (cm *Server) Stop() {
	zap.S().Debug("server stop ...")
	ch := make(chan struct{}, 1000) // 控制并发数
	var wg sync.WaitGroup
	connAll := cm.GetConnAll()
	for _, conn := range connAll {
		ch <- struct{}{}
		wg.Add(1)
		c := conn
		go func() {
			defer func() {
				wg.Done()
				<-ch
			}()
			c.Stop()
		}()
	}
	close(ch)
	wg.Wait()
}

// AddConn 添加连接
func (cm *Server) AddConn(userId int64, conn *Conn) {
	cm.connMap.Store(userId, conn)
	zap.S().Debugf("connection UserId=%d add to Server\n", userId)
}

// RemoveConn 删除连接
func (cm *Server) RemoveConn(userId int64) {
	cm.connMap.Delete(userId)
	zap.S().Debugf("connection UserId=%d remove from Server\n", userId)
}

// GetConn 根据userid获取相应的连接
func (cm *Server) GetConn(userId int64) *Conn {
	value, ok := cm.connMap.Load(userId)
	if ok {
		return value.(*Conn)
	}
	return nil
}

// GetConnAll 获取全部连接
func (cm *Server) GetConnAll() []*Conn {
	conns := make([]*Conn, 0)
	cm.connMap.Range(func(key, value interface{}) bool {
		conn := value.(*Conn)
		conns = append(conns, conn)
		return true
	})
	return conns
}

// SendMessageAll 进行本地推送
func (cm *Server) SendMessageAll(userId2Msg map[int64][]byte) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 5) // 限制并发数
	for userId, data := range userId2Msg {
		ch <- struct{}{}
		wg.Add(1)
		go func(userId int64, data []byte) {
			defer func() {
				<-ch
				wg.Done()
			}()
			conn := ConnManager.GetConn(userId)
			if conn != nil {
				conn.SendMsg(userId, data)
			}
		}(userId, data)
	}
	close(ch)
	wg.Wait()
}

// StartWorkerPool 启动 worker 工作池
func (cm *Server) StartWorkerPool() {
	// 初始化并启动 worker 工作池
	for i := 0; i < len(cm.taskQueue); i++ {
		// 初始化
		cm.taskQueue[i] = make(chan *Req, config.Configs.App.MaxWorkerTask) // 初始化worker队列中，每个worker的队列长度
		// 启动worker
		go cm.StartOneWorker(i, cm.taskQueue[i])
	}
}

// StartOneWorker 启动 worker 的工作流程
func (cm *Server) StartOneWorker(workerID int, taskQueue chan *Req) {
	zap.S().Debug("Worker ID = ", workerID, " is started.")
	for {
		select {
		case req := <-taskQueue:
			req.f()
		}
	}
}

// SendMsgToTaskQueue 将消息交给 taskQueue，由 worker 调度处理
func (cm *Server) SendMsgToTaskQueue(req *Req) {
	if len(cm.taskQueue) > 0 {
		// 根据ConnID来分配当前的连接应该由哪个worker负责处理，保证同一个连接的消息处理串行
		// 轮询的平均分配法则

		//得到需要处理此条连接的workerID
		workerID := req.conn.ConnId % int64(len(cm.taskQueue))

		// 将消息发给对应的 taskQueue
		cm.taskQueue[workerID] <- req
	} else {
		// 可能导致消息乱序
		go req.f()
	}
}
