// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: api/llm.proto

package api

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

// LLMClient is the client API for LLM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LLMClient interface {
	Completion(ctx context.Context, in *CompletionRequest, opts ...grpc.CallOption) (*CompletionResponse, error)
	CompletionStream(ctx context.Context, in *CompletionStreamRequest, opts ...grpc.CallOption) (LLM_CompletionStreamClient, error)
	ChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, in *ChatCompletionStreamRequest, opts ...grpc.CallOption) (LLM_ChatCompletionStreamClient, error)
}

type lLMClient struct {
	cc grpc.ClientConnInterface
}

func NewLLMClient(cc grpc.ClientConnInterface) LLMClient {
	return &lLMClient{cc}
}

func (c *lLMClient) Completion(ctx context.Context, in *CompletionRequest, opts ...grpc.CallOption) (*CompletionResponse, error) {
	out := new(CompletionResponse)
	err := c.cc.Invoke(ctx, "/llm.LLM/Completion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLMClient) CompletionStream(ctx context.Context, in *CompletionStreamRequest, opts ...grpc.CallOption) (LLM_CompletionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &LLM_ServiceDesc.Streams[0], "/llm.LLM/CompletionStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &lLMCompletionStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LLM_CompletionStreamClient interface {
	Recv() (*CompletionStreamResponse, error)
	grpc.ClientStream
}

type lLMCompletionStreamClient struct {
	grpc.ClientStream
}

func (x *lLMCompletionStreamClient) Recv() (*CompletionStreamResponse, error) {
	m := new(CompletionStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *lLMClient) ChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error) {
	out := new(ChatCompletionResponse)
	err := c.cc.Invoke(ctx, "/llm.LLM/ChatCompletion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLMClient) ChatCompletionStream(ctx context.Context, in *ChatCompletionStreamRequest, opts ...grpc.CallOption) (LLM_ChatCompletionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &LLM_ServiceDesc.Streams[1], "/llm.LLM/ChatCompletionStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &lLMChatCompletionStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LLM_ChatCompletionStreamClient interface {
	Recv() (*ChatCompletionStreamResponse, error)
	grpc.ClientStream
}

type lLMChatCompletionStreamClient struct {
	grpc.ClientStream
}

func (x *lLMChatCompletionStreamClient) Recv() (*ChatCompletionStreamResponse, error) {
	m := new(ChatCompletionStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LLMServer is the server API for LLM service.
// All implementations must embed UnimplementedLLMServer
// for forward compatibility
type LLMServer interface {
	Completion(context.Context, *CompletionRequest) (*CompletionResponse, error)
	CompletionStream(*CompletionStreamRequest, LLM_CompletionStreamServer) error
	ChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(*ChatCompletionStreamRequest, LLM_ChatCompletionStreamServer) error
	mustEmbedUnimplementedLLMServer()
}

// UnimplementedLLMServer must be embedded to have forward compatible implementations.
type UnimplementedLLMServer struct {
}

func (UnimplementedLLMServer) Completion(context.Context, *CompletionRequest) (*CompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Completion not implemented")
}
func (UnimplementedLLMServer) CompletionStream(*CompletionStreamRequest, LLM_CompletionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CompletionStream not implemented")
}
func (UnimplementedLLMServer) ChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChatCompletion not implemented")
}
func (UnimplementedLLMServer) ChatCompletionStream(*ChatCompletionStreamRequest, LLM_ChatCompletionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ChatCompletionStream not implemented")
}
func (UnimplementedLLMServer) mustEmbedUnimplementedLLMServer() {}

// UnsafeLLMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LLMServer will
// result in compilation errors.
type UnsafeLLMServer interface {
	mustEmbedUnimplementedLLMServer()
}

func RegisterLLMServer(s grpc.ServiceRegistrar, srv LLMServer) {
	s.RegisterService(&LLM_ServiceDesc, srv)
}

func _LLM_Completion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompletionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLMServer).Completion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/llm.LLM/Completion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLMServer).Completion(ctx, req.(*CompletionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLM_CompletionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CompletionStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LLMServer).CompletionStream(m, &lLMCompletionStreamServer{stream})
}

type LLM_CompletionStreamServer interface {
	Send(*CompletionStreamResponse) error
	grpc.ServerStream
}

type lLMCompletionStreamServer struct {
	grpc.ServerStream
}

func (x *lLMCompletionStreamServer) Send(m *CompletionStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _LLM_ChatCompletion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatCompletionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLMServer).ChatCompletion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/llm.LLM/ChatCompletion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLMServer).ChatCompletion(ctx, req.(*ChatCompletionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLM_ChatCompletionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChatCompletionStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LLMServer).ChatCompletionStream(m, &lLMChatCompletionStreamServer{stream})
}

type LLM_ChatCompletionStreamServer interface {
	Send(*ChatCompletionStreamResponse) error
	grpc.ServerStream
}

type lLMChatCompletionStreamServer struct {
	grpc.ServerStream
}

func (x *lLMChatCompletionStreamServer) Send(m *ChatCompletionStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// LLM_ServiceDesc is the grpc.ServiceDesc for LLM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LLM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "llm.LLM",
	HandlerType: (*LLMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Completion",
			Handler:    _LLM_Completion_Handler,
		},
		{
			MethodName: "ChatCompletion",
			Handler:    _LLM_ChatCompletion_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CompletionStream",
			Handler:       _LLM_CompletionStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ChatCompletionStream",
			Handler:       _LLM_ChatCompletionStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/llm.proto",
}