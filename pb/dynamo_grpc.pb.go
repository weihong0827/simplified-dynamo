// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: dynamo.proto

package dynamo

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

// KeyValueStoreClient is the client API for KeyValueStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyValueStoreClient interface {
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Gossip(ctx context.Context, in *GossipMessage, opts ...grpc.CallOption) (*GossipAck, error)
	// temporarily send the replica to other machines to store
	HintedHandoff(ctx context.Context, in *HintedHandoffWriteRequest, opts ...grpc.CallOption) (*Empty, error)
	// when the node back alive again, it will send the replica back to the node
	SendReplica(ctx context.Context, in *BulkWriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
}

type keyValueStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyValueStoreClient(cc grpc.ClientConnInterface) KeyValueStoreClient {
	return &keyValueStoreClient{cc}
}

func (c *keyValueStoreClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Gossip(ctx context.Context, in *GossipMessage, opts ...grpc.CallOption) (*GossipAck, error) {
	out := new(GossipAck)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Gossip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) HintedHandoff(ctx context.Context, in *HintedHandoffWriteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/HintedHandoff", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) SendReplica(ctx context.Context, in *BulkWriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/SendReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyValueStoreServer is the server API for KeyValueStore service.
// All implementations must embed UnimplementedKeyValueStoreServer
// for forward compatibility
type KeyValueStoreServer interface {
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Gossip(context.Context, *GossipMessage) (*GossipAck, error)
	// temporarily send the replica to other machines to store
	HintedHandoff(context.Context, *HintedHandoffWriteRequest) (*Empty, error)
	// when the node back alive again, it will send the replica back to the node
	SendReplica(context.Context, *BulkWriteRequest) (*WriteResponse, error)
	mustEmbedUnimplementedKeyValueStoreServer()
}

// UnimplementedKeyValueStoreServer must be embedded to have forward compatible implementations.
type UnimplementedKeyValueStoreServer struct {
}

func (UnimplementedKeyValueStoreServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedKeyValueStoreServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedKeyValueStoreServer) Gossip(context.Context, *GossipMessage) (*GossipAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedKeyValueStoreServer) HintedHandoff(context.Context, *HintedHandoffWriteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HintedHandoff not implemented")
}
func (UnimplementedKeyValueStoreServer) SendReplica(context.Context, *BulkWriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReplica not implemented")
}
func (UnimplementedKeyValueStoreServer) mustEmbedUnimplementedKeyValueStoreServer() {}

// UnsafeKeyValueStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyValueStoreServer will
// result in compilation errors.
type UnsafeKeyValueStoreServer interface {
	mustEmbedUnimplementedKeyValueStoreServer()
}

func RegisterKeyValueStoreServer(s grpc.ServiceRegistrar, srv KeyValueStoreServer) {
	s.RegisterService(&KeyValueStore_ServiceDesc, srv)
}

func _KeyValueStore_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Gossip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Gossip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Gossip",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Gossip(ctx, req.(*GossipMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_HintedHandoff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HintedHandoffWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).HintedHandoff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/HintedHandoff",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).HintedHandoff(ctx, req.(*HintedHandoffWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_SendReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).SendReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/SendReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).SendReplica(ctx, req.(*BulkWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyValueStore_ServiceDesc is the grpc.ServiceDesc for KeyValueStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyValueStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dynamo.KeyValueStore",
	HandlerType: (*KeyValueStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _KeyValueStore_Write_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _KeyValueStore_Read_Handler,
		},
		{
			MethodName: "Gossip",
			Handler:    _KeyValueStore_Gossip_Handler,
		},
		{
			MethodName: "HintedHandoff",
			Handler:    _KeyValueStore_HintedHandoff_Handler,
		},
		{
			MethodName: "SendReplica",
			Handler:    _KeyValueStore_SendReplica_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dynamo.proto",
}

// PTPClient is the client API for PTP service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PTPClient interface {
	AddNewNode(ctx context.Context, in *Node, opts ...grpc.CallOption) (*MembershipList, error)
	Gossip(ctx context.Context, in *GossipMessage, opts ...grpc.CallOption) (*GossipAck, error)
}

type pTPClient struct {
	cc grpc.ClientConnInterface
}

func NewPTPClient(cc grpc.ClientConnInterface) PTPClient {
	return &pTPClient{cc}
}

func (c *pTPClient) AddNewNode(ctx context.Context, in *Node, opts ...grpc.CallOption) (*MembershipList, error) {
	out := new(MembershipList)
	err := c.cc.Invoke(ctx, "/dynamo.PTP/AddNewNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pTPClient) Gossip(ctx context.Context, in *GossipMessage, opts ...grpc.CallOption) (*GossipAck, error) {
	out := new(GossipAck)
	err := c.cc.Invoke(ctx, "/dynamo.PTP/Gossip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PTPServer is the server API for PTP service.
// All implementations must embed UnimplementedPTPServer
// for forward compatibility
type PTPServer interface {
	AddNewNode(context.Context, *Node) (*MembershipList, error)
	Gossip(context.Context, *GossipMessage) (*GossipAck, error)
	mustEmbedUnimplementedPTPServer()
}

// UnimplementedPTPServer must be embedded to have forward compatible implementations.
type UnimplementedPTPServer struct {
}

func (UnimplementedPTPServer) AddNewNode(context.Context, *Node) (*MembershipList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNewNode not implemented")
}
func (UnimplementedPTPServer) Gossip(context.Context, *GossipMessage) (*GossipAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedPTPServer) mustEmbedUnimplementedPTPServer() {}

// UnsafePTPServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PTPServer will
// result in compilation errors.
type UnsafePTPServer interface {
	mustEmbedUnimplementedPTPServer()
}

func RegisterPTPServer(s grpc.ServiceRegistrar, srv PTPServer) {
	s.RegisterService(&PTP_ServiceDesc, srv)
}

func _PTP_AddNewNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PTPServer).AddNewNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.PTP/AddNewNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PTPServer).AddNewNode(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _PTP_Gossip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PTPServer).Gossip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.PTP/Gossip",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PTPServer).Gossip(ctx, req.(*GossipMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// PTP_ServiceDesc is the grpc.ServiceDesc for PTP service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PTP_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dynamo.PTP",
	HandlerType: (*PTPServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddNewNode",
			Handler:    _PTP_AddNewNode_Handler,
		},
		{
			MethodName: "Gossip",
			Handler:    _PTP_Gossip_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dynamo.proto",
}
