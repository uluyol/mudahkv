// Code generated by protoc-gen-go.
// source: mudah.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	mudah.proto

It has these top-level messages:
	Key
	KVChunk
	ListRequest
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Key struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *Key) Reset()                    { *m = Key{} }
func (m *Key) String() string            { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()               {}
func (*Key) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type KVChunk struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *KVChunk) Reset()                    { *m = KVChunk{} }
func (m *KVChunk) String() string            { return proto.CompactTextString(m) }
func (*KVChunk) ProtoMessage()               {}
func (*KVChunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ListRequest struct {
	Prefix string `protobuf:"bytes,1,opt,name=prefix" json:"prefix,omitempty"`
}

func (m *ListRequest) Reset()                    { *m = ListRequest{} }
func (m *ListRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()               {}
func (*ListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Key)(nil), "mudah.Key")
	proto.RegisterType((*KVChunk)(nil), "mudah.KVChunk")
	proto.RegisterType((*ListRequest)(nil), "mudah.ListRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Mudah service

type MudahClient interface {
	Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (Mudah_GetClient, error)
	Set(ctx context.Context, opts ...grpc.CallOption) (Mudah_SetClient, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Mudah_ListClient, error)
}

type mudahClient struct {
	cc *grpc.ClientConn
}

func NewMudahClient(cc *grpc.ClientConn) MudahClient {
	return &mudahClient{cc}
}

func (c *mudahClient) Get(ctx context.Context, in *Key, opts ...grpc.CallOption) (Mudah_GetClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Mudah_serviceDesc.Streams[0], c.cc, "/mudah.Mudah/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &mudahGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Mudah_GetClient interface {
	Recv() (*KVChunk, error)
	grpc.ClientStream
}

type mudahGetClient struct {
	grpc.ClientStream
}

func (x *mudahGetClient) Recv() (*KVChunk, error) {
	m := new(KVChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mudahClient) Set(ctx context.Context, opts ...grpc.CallOption) (Mudah_SetClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Mudah_serviceDesc.Streams[1], c.cc, "/mudah.Mudah/Set", opts...)
	if err != nil {
		return nil, err
	}
	x := &mudahSetClient{stream}
	return x, nil
}

type Mudah_SetClient interface {
	Send(*KVChunk) error
	CloseAndRecv() (*Key, error)
	grpc.ClientStream
}

type mudahSetClient struct {
	grpc.ClientStream
}

func (x *mudahSetClient) Send(m *KVChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mudahSetClient) CloseAndRecv() (*Key, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Key)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mudahClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Mudah_ListClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Mudah_serviceDesc.Streams[2], c.cc, "/mudah.Mudah/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &mudahListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Mudah_ListClient interface {
	Recv() (*KVChunk, error)
	grpc.ClientStream
}

type mudahListClient struct {
	grpc.ClientStream
}

func (x *mudahListClient) Recv() (*KVChunk, error) {
	m := new(KVChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Mudah service

type MudahServer interface {
	Get(*Key, Mudah_GetServer) error
	Set(Mudah_SetServer) error
	List(*ListRequest, Mudah_ListServer) error
}

func RegisterMudahServer(s *grpc.Server, srv MudahServer) {
	s.RegisterService(&_Mudah_serviceDesc, srv)
}

func _Mudah_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Key)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MudahServer).Get(m, &mudahGetServer{stream})
}

type Mudah_GetServer interface {
	Send(*KVChunk) error
	grpc.ServerStream
}

type mudahGetServer struct {
	grpc.ServerStream
}

func (x *mudahGetServer) Send(m *KVChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _Mudah_Set_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MudahServer).Set(&mudahSetServer{stream})
}

type Mudah_SetServer interface {
	SendAndClose(*Key) error
	Recv() (*KVChunk, error)
	grpc.ServerStream
}

type mudahSetServer struct {
	grpc.ServerStream
}

func (x *mudahSetServer) SendAndClose(m *Key) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mudahSetServer) Recv() (*KVChunk, error) {
	m := new(KVChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Mudah_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MudahServer).List(m, &mudahListServer{stream})
}

type Mudah_ListServer interface {
	Send(*KVChunk) error
	grpc.ServerStream
}

type mudahListServer struct {
	grpc.ServerStream
}

func (x *mudahListServer) Send(m *KVChunk) error {
	return x.ServerStream.SendMsg(m)
}

var _Mudah_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mudah.Mudah",
	HandlerType: (*MudahServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _Mudah_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Set",
			Handler:       _Mudah_Set_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "List",
			Handler:       _Mudah_List_Handler,
			ServerStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x2d, 0x4d, 0x49,
	0xcc, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0xc4, 0xb9, 0x98, 0xbd,
	0x53, 0x2b, 0x85, 0x04, 0xb8, 0x98, 0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83,
	0x40, 0x4c, 0x25, 0x43, 0x2e, 0x76, 0xef, 0x30, 0xe7, 0x8c, 0xd2, 0xbc, 0x6c, 0x4c, 0x49, 0x21,
	0x11, 0x2e, 0xd6, 0xb2, 0xc4, 0x9c, 0xd2, 0x54, 0x09, 0x26, 0xa0, 0x18, 0x4f, 0x10, 0x84, 0xa3,
	0xa4, 0xca, 0xc5, 0xed, 0x93, 0x59, 0x5c, 0x12, 0x94, 0x5a, 0x58, 0x9a, 0x5a, 0x5c, 0x22, 0x24,
	0xc6, 0xc5, 0x56, 0x50, 0x94, 0x9a, 0x96, 0x59, 0x01, 0xd5, 0x09, 0xe5, 0x19, 0xb5, 0x32, 0x72,
	0xb1, 0xfa, 0x82, 0x2c, 0x17, 0x52, 0xe5, 0x62, 0x76, 0x4f, 0x2d, 0x11, 0xe2, 0xd2, 0x83, 0x38,
	0x0c, 0xe8, 0x10, 0x29, 0x3e, 0x18, 0x1b, 0x62, 0xb7, 0x12, 0x83, 0x01, 0x23, 0x48, 0x59, 0x30,
	0x50, 0x19, 0x9a, 0x94, 0x14, 0x92, 0x36, 0x25, 0x06, 0x0d, 0x46, 0x21, 0x3d, 0x2e, 0x16, 0x90,
	0xf5, 0x42, 0x42, 0x50, 0x71, 0x24, 0xb7, 0x60, 0x33, 0xd6, 0x89, 0x25, 0x8a, 0xa9, 0x20, 0x29,
	0x89, 0x0d, 0x1c, 0x1c, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe8, 0x77, 0xe7, 0x66, 0x1d,
	0x01, 0x00, 0x00,
}
