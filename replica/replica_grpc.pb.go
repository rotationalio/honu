// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: replica/v1/replica.proto

package replica

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

const (
	Replication_Gossip_FullMethodName = "/honu.replica.v1.Replication/Gossip"
)

// ReplicationClient is the client API for Replication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationClient interface {
	// Gossip implements biltateral anti-entropy: during a Gossip session the initiating
	// replica pushes updates to the remote peer and pulls requested changes. Using
	// bidirectional streaming, the initiating peer sends data-less sync messages with
	// the versions of objects it stores locally. The remote replica then responds with
	// data if its local version is later or sends a sync message back requesting the
	// data from the initating replica if its local version is earlier (no exchange)
	// occurs if both replicas have the same version. At the end of a gossip session,
	// both replicas should have synchronized and have identical underlying data stores.
	Gossip(ctx context.Context, opts ...grpc.CallOption) (Replication_GossipClient, error)
}

type replicationClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationClient(cc grpc.ClientConnInterface) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) Gossip(ctx context.Context, opts ...grpc.CallOption) (Replication_GossipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Replication_ServiceDesc.Streams[0], Replication_Gossip_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &replicationGossipClient{stream}
	return x, nil
}

type Replication_GossipClient interface {
	Send(*Sync) error
	Recv() (*Sync, error)
	grpc.ClientStream
}

type replicationGossipClient struct {
	grpc.ClientStream
}

func (x *replicationGossipClient) Send(m *Sync) error {
	return x.ClientStream.SendMsg(m)
}

func (x *replicationGossipClient) Recv() (*Sync, error) {
	m := new(Sync)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReplicationServer is the server API for Replication service.
// All implementations must embed UnimplementedReplicationServer
// for forward compatibility
type ReplicationServer interface {
	// Gossip implements biltateral anti-entropy: during a Gossip session the initiating
	// replica pushes updates to the remote peer and pulls requested changes. Using
	// bidirectional streaming, the initiating peer sends data-less sync messages with
	// the versions of objects it stores locally. The remote replica then responds with
	// data if its local version is later or sends a sync message back requesting the
	// data from the initating replica if its local version is earlier (no exchange)
	// occurs if both replicas have the same version. At the end of a gossip session,
	// both replicas should have synchronized and have identical underlying data stores.
	Gossip(Replication_GossipServer) error
	mustEmbedUnimplementedReplicationServer()
}

// UnimplementedReplicationServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (UnimplementedReplicationServer) Gossip(Replication_GossipServer) error {
	return status.Errorf(codes.Unimplemented, "method Gossip not implemented")
}
func (UnimplementedReplicationServer) mustEmbedUnimplementedReplicationServer() {}

// UnsafeReplicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServer will
// result in compilation errors.
type UnsafeReplicationServer interface {
	mustEmbedUnimplementedReplicationServer()
}

func RegisterReplicationServer(s grpc.ServiceRegistrar, srv ReplicationServer) {
	s.RegisterService(&Replication_ServiceDesc, srv)
}

func _Replication_Gossip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ReplicationServer).Gossip(&replicationGossipServer{stream})
}

type Replication_GossipServer interface {
	Send(*Sync) error
	Recv() (*Sync, error)
	grpc.ServerStream
}

type replicationGossipServer struct {
	grpc.ServerStream
}

func (x *replicationGossipServer) Send(m *Sync) error {
	return x.ServerStream.SendMsg(m)
}

func (x *replicationGossipServer) Recv() (*Sync, error) {
	m := new(Sync)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Replication_ServiceDesc is the grpc.ServiceDesc for Replication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Replication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "honu.replica.v1.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Gossip",
			Handler:       _Replication_Gossip_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "replica/v1/replica.proto",
}
