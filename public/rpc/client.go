package rpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"im/public/protocol"
)

var (
	ConnServerClient protocol.ConnectClient
)

// GetServerClient 获取 grpc 连接
func GetServerClient(addr string) protocol.ConnectClient {
	// todo 安全连接？
	client, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Error("grpc client Dial err, err:", err)
		panic(err)
	}
	ConnServerClient = protocol.NewConnectClient(client)
	return ConnServerClient
}
