package main

import (
	"im/config"
	_ "im/public"
	"im/public/etcd"
	"im/router"
	"im/server/mq"
	"im/service/rpc_server"
)

func main() {
	mq.InitMessageMQ(config.Configs.RabbitMQ.URL)
	// 初始化服务注册发现
	go etcd.InitETCD()

	// 启动 http 服务
	go router.HTTPRouter()

	// 启动 rpc 服务
	go rpc_server.InitRPCServer()

	// 启动 websocket 服务
	router.WSRouter()
}
