// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chittychat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ChittyChatClient is the client API for ChittyChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatClient interface {
	Publish(ctx context.Context, in *PublishMessage, opts ...grpc.CallOption) (*PublishReply, error)
}

type chittyChatClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatClient(cc grpc.ClientConnInterface) ChittyChatClient {
	return &chittyChatClient{cc}
}

func (c *chittyChatClient) Publish(ctx context.Context, in *PublishMessage, opts ...grpc.CallOption) (*PublishReply, error) {
	out := new(PublishReply)
	err := c.cc.Invoke(ctx, "/chittychat.ChittyChat/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServer is the server API for ChittyChat service.
// All implementations must embed UnimplementedChittyChatServer
// for forward compatibility
type ChittyChatServer interface {
	Publish(context.Context, *PublishMessage) (*PublishReply, error)
	mustEmbedUnimplementedChittyChatServer()
}

// UnimplementedChittyChatServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServer struct {
}

func (UnimplementedChittyChatServer) Publish(context.Context, *PublishMessage) (*PublishReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedChittyChatServer) mustEmbedUnimplementedChittyChatServer() {}

// UnsafeChittyChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServer will
// result in compilation errors.
type UnsafeChittyChatServer interface {
	mustEmbedUnimplementedChittyChatServer()
}

func RegisterChittyChatServer(s grpc.ServiceRegistrar, srv ChittyChatServer) {
	s.RegisterService(&ChittyChat_ServiceDesc, srv)
}

func _ChittyChat_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chittychat.ChittyChat/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServer).Publish(ctx, req.(*PublishMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChat_ServiceDesc is the grpc.ServiceDesc for ChittyChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chittychat.ChittyChat",
	HandlerType: (*ChittyChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _ChittyChat_Publish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ChittyChat/chittychat.proto",
}
