package etcd

import (
	"fmt"
	"go.uber.org/zap"
	"im/config"
	"im/server/etcd"
	"time"
)

var (
	DiscoverySer *etcd.Discovery
)

// InitETCD 初始化服务注册发现
// 1. 初始化服务注册，将自己当前启动的 RPC 端口注册到 etcd
// 2. 初始化服务发现，启动 watcher 监听所有 RPC 端口，以便有需要时能直接获取当前注册在 ETCD 的服务
func InitETCD() {
	hostPort := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.RPCPort)
	zap.S().Info("[InitETCD] [hostPort] = ", hostPort)

	// 注册服务并设置 k-v 租约
	err := etcd.RegisterServer(config.Configs.ETCD.ServerList+hostPort, hostPort, 5)
	if err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)

	DiscoverySer, err = etcd.NewDiscovery()
	if err != nil {
		zap.S().Error("[InitETCD] [NewDiscovery] [err] = ", err)
		return
	}

	// 阻塞监听
	DiscoverySer.WatchService(config.Configs.ETCD.ServerList)
}
