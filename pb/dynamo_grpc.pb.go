// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
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
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Forward(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Join(ctx context.Context, in *Node, opts ...grpc.CallOption) (*JoinResponse, error)
	Gossip(ctx context.Context, in *GossipMessage, opts ...grpc.CallOption) (*GossipAck, error)
	// temporarily send the replica to other machines to store
	HintedHandoffWrite(ctx context.Context, in *HintedHandoffWriteRequest, opts ...grpc.CallOption) (*HintedHandoffWriteResponse, error)
	// when the node back alive again, it will send the replica back to the node
	SendReplica(ctx context.Context, in *BulkWriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Delete(ctx context.Context, in *ReplicaDeleteRequest, opts ...grpc.CallOption) (*Empty, error)
	BulkWrite(ctx context.Context, in *BulkWriteRequest, opts ...grpc.CallOption) (*Empty, error)
	Transfer(ctx context.Context, in *ReplicaTransferRequest, opts ...grpc.CallOption) (*Empty, error)
}

type keyValueStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyValueStoreClient(cc grpc.ClientConnInterface) KeyValueStoreClient {
	return &keyValueStoreClient{cc}
}

func (c *keyValueStoreClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (c *keyValueStoreClient) Forward(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Forward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Join(ctx context.Context, in *Node, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Join", in, out, opts...)
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

func (c *keyValueStoreClient) HintedHandoffWrite(ctx context.Context, in *HintedHandoffWriteRequest, opts ...grpc.CallOption) (*HintedHandoffWriteResponse, error) {
	out := new(HintedHandoffWriteResponse)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/HintedHandoffWrite", in, out, opts...)
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

func (c *keyValueStoreClient) Delete(ctx context.Context, in *ReplicaDeleteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) BulkWrite(ctx context.Context, in *BulkWriteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/BulkWrite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueStoreClient) Transfer(ctx context.Context, in *ReplicaTransferRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/dynamo.KeyValueStore/Transfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyValueStoreServer is the server API for KeyValueStore service.
// All implementations must embed UnimplementedKeyValueStoreServer
// for forward compatibility
type KeyValueStoreServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Forward(context.Context, *WriteRequest) (*WriteResponse, error)
	Join(context.Context, *Node) (*JoinResponse, error)
	Gossip(context.Context, *GossipMessage) (*GossipAck, error)
	// temporarily send the replica to other machines to store
	HintedHandoffWrite(context.Context, *HintedHandoffWriteRequest) (*HintedHandoffWriteResponse, error)
	// when the node back alive again, it will send the replica back to the node
	SendReplica(context.Context, *BulkWriteRequest) (*WriteResponse, error)
	Delete(context.Context, *ReplicaDeleteRequest) (*Empty, error)
	BulkWrite(context.Context, *BulkWriteRequest) (*Empty, error)
	Transfer(context.Context, *ReplicaTransferRequest) (*Empty, error)
	mustEmbedUnimplementedKeyValueStoreServer()
}

// UnimplementedKeyValueStoreServer must be embedded to have forward compatible implementations.
type UnimplementedKeyValueStoreServer struct {
}

func (UnimplementedKeyValueStoreServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKeyValueStoreServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedKeyValueStoreServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedKeyValueStoreServer) Forward(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Forward not implemented")
}
func (UnimplementedKeyValueStoreServer) Join(context.Context, *Node) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedKeyValueStoreServer) Gossip(context.Context, *GossipMessage) (*GossipAck, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedKeyValueStoreServer) HintedHandoffWrite(context.Context, *HintedHandoffWriteRequest) (*HintedHandoffWriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HintedHandoffWrite not implemented")
}
func (UnimplementedKeyValueStoreServer) SendReplica(context.Context, *BulkWriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReplica not implemented")
}
func (UnimplementedKeyValueStoreServer) Delete(context.Context, *ReplicaDeleteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKeyValueStoreServer) BulkWrite(context.Context, *BulkWriteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkWrite not implemented")
}
func (UnimplementedKeyValueStoreServer) Transfer(context.Context, *ReplicaTransferRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
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

func _KeyValueStore_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
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

func _KeyValueStore_Forward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Forward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Forward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Forward(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Join(ctx, req.(*Node))
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

func _KeyValueStore_HintedHandoffWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HintedHandoffWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).HintedHandoffWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/HintedHandoffWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).HintedHandoffWrite(ctx, req.(*HintedHandoffWriteRequest))
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

func _KeyValueStore_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicaDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Delete(ctx, req.(*ReplicaDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_BulkWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).BulkWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/BulkWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).BulkWrite(ctx, req.(*BulkWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueStore_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicaTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueStoreServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dynamo.KeyValueStore/Transfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueStoreServer).Transfer(ctx, req.(*ReplicaTransferRequest))
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
			MethodName: "Ping",
			Handler:    _KeyValueStore_Ping_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _KeyValueStore_Write_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _KeyValueStore_Read_Handler,
		},
		{
			MethodName: "Forward",
			Handler:    _KeyValueStore_Forward_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _KeyValueStore_Join_Handler,
		},
		{
			MethodName: "Gossip",
			Handler:    _KeyValueStore_Gossip_Handler,
		},
		{
			MethodName: "HintedHandoffWrite",
			Handler:    _KeyValueStore_HintedHandoffWrite_Handler,
		},
		{
			MethodName: "SendReplica",
			Handler:    _KeyValueStore_SendReplica_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KeyValueStore_Delete_Handler,
		},
		{
			MethodName: "BulkWrite",
			Handler:    _KeyValueStore_BulkWrite_Handler,
		},
		{
			MethodName: "Transfer",
			Handler:    _KeyValueStore_Transfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dynamo.proto",
}
