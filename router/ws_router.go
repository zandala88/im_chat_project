package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"im/config"
	"im/service/ws"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSRouter websocket 路由
func WSRouter() {
	server := ws.GetServer()

	// 开启worker工作池
	server.StartWorkerPool()

	// 开启心跳超时检测
	checker := ws.NewHeartbeatChecker(time.Second*time.Duration(config.Configs.App.HeartbeatInterval), server)
	go checker.Start()

	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	pprofRouter := gin.New()
	pprof.Register(pprofRouter)
	r.Any("/debug/pprof/*any", gin.WrapH(pprofRouter))
	var connID int64

	r.GET("/ws", func(c *gin.Context) {
		// 升级协议  http -> websocket
		WsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.S().Error("[WSRouter] [upgrader.Upgrade] [err] = ", err)
			return
		}

		// 初始化连接
		conn := ws.NewConnection(server, WsConn, connID)
		connID++

		// 开启读写线程
		go conn.Start()
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.WebsocketPort),
		Handler: r,
	}

	go func() {
		zap.S().Info("[WSRouter] websocket 启动：", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalf("websocket listen: %s\n", err)
			zap.S().Error("[WSRouter] [ListenAndServe] [err] = ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.S().Info("Shutdown Websocket Server ...")

	// 关闭服务
	server.Stop()
	checker.Stop()

	// 5s 超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Error("[WSRouter] [Shutdown] [err] = ", err)
	}

	zap.S().Info("[WSRouter] Server exiting")
}
