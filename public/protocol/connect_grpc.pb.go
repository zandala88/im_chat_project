// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.1
// source: connect.proto

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Connect_DeliverMessage_FullMethodName    = "/protocol.Connect/DeliverMessage"
	Connect_DeliverMessageAll_FullMethodName = "/protocol.Connect/DeliverMessageAll"
)

// ConnectClient is the client API for Connect service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectClient interface {
	// 私聊消息投递
	DeliverMessage(ctx context.Context, in *DeliverMessageReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 群聊消息投递
	DeliverMessageAll(ctx context.Context, in *DeliverMessageAllReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type connectClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectClient(cc grpc.ClientConnInterface) ConnectClient {
	return &connectClient{cc}
}

func (c *connectClient) DeliverMessage(ctx context.Context, in *DeliverMessageReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Connect_DeliverMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectClient) DeliverMessageAll(ctx context.Context, in *DeliverMessageAllReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Connect_DeliverMessageAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectServer is the server API for Connect service.
// All implementations must embed UnimplementedConnectServer
// for forward compatibility.
type ConnectServer interface {
	// 私聊消息投递
	DeliverMessage(context.Context, *DeliverMessageReq) (*emptypb.Empty, error)
	// 群聊消息投递
	DeliverMessageAll(context.Context, *DeliverMessageAllReq) (*emptypb.Empty, error)
	mustEmbedUnimplementedConnectServer()
}

// UnimplementedConnectServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedConnectServer struct{}

func (UnimplementedConnectServer) DeliverMessage(context.Context, *DeliverMessageReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeliverMessage not implemented")
}
func (UnimplementedConnectServer) DeliverMessageAll(context.Context, *DeliverMessageAllReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeliverMessageAll not implemented")
}
func (UnimplementedConnectServer) mustEmbedUnimplementedConnectServer() {}
func (UnimplementedConnectServer) testEmbeddedByValue()                 {}

// UnsafeConnectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectServer will
// result in compilation errors.
type UnsafeConnectServer interface {
	mustEmbedUnimplementedConnectServer()
}

func RegisterConnectServer(s grpc.ServiceRegistrar, srv ConnectServer) {
	// If the following call pancis, it indicates UnimplementedConnectServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Connect_ServiceDesc, srv)
}

func _Connect_DeliverMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeliverMessageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectServer).DeliverMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Connect_DeliverMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectServer).DeliverMessage(ctx, req.(*DeliverMessageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Connect_DeliverMessageAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeliverMessageAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectServer).DeliverMessageAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Connect_DeliverMessageAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectServer).DeliverMessageAll(ctx, req.(*DeliverMessageAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Connect_ServiceDesc is the grpc.ServiceDesc for Connect service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Connect_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Connect",
	HandlerType: (*ConnectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeliverMessage",
			Handler:    _Connect_DeliverMessage_Handler,
		},
		{
			MethodName: "DeliverMessageAll",
			Handler:    _Connect_DeliverMessageAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connect.proto",
}
