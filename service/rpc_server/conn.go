package rpc_server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"im/config"
	"im/public/protocol"
	"im/service/ws"
	"log"
	"net"
)

type ConnectServer struct {
	protocol.UnsafeConnectServer // 禁止向前兼容
}

func (*ConnectServer) DeliverMessage(ctx context.Context, req *protocol.DeliverMessageReq) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}

	// 获取本地连接
	conn := ws.GetServer().GetConn(req.ReceiverId)
	if conn == nil || conn.GetUserId() != req.ReceiverId {
		zap.S().Debug("[DeliverMessage] 连接不存在 user_id:", req.ReceiverId)
		return resp, nil
	}

	// 消息发送
	conn.SendMsg(req.ReceiverId, req.Data)

	return resp, nil
}

func (*ConnectServer) DeliverMessageAll(ctx context.Context, req *protocol.DeliverMessageAllReq) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}

	// 进行本地推送
	ws.GetServer().SendMessageAll(req.GetReceiverId_2Data())

	return resp, nil
}

func InitRPCServer() {
	rpcPort := config.Configs.App.RPCPort

	server := grpc.NewServer()
	protocol.RegisterConnectServer(server, &ConnectServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	zap.S().Debug("rpc server 启动 ", rpcPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to rpc serve: %v", err)
	}
}
