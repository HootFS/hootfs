// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

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

// DiscoverServiceClient is the client API for DiscoverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscoverServiceClient interface {
	JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error)
	GetActive(ctx context.Context, in *GetActiveRequest, opts ...grpc.CallOption) (*GetActiveResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type discoverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscoverServiceClient(cc grpc.ClientConnInterface) DiscoverServiceClient {
	return &discoverServiceClient{cc}
}

func (c *discoverServiceClient) JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error) {
	out := new(JoinClusterResponse)
	err := c.cc.Invoke(ctx, "/hootfs.discovery.DiscoverService/JoinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoverServiceClient) GetActive(ctx context.Context, in *GetActiveRequest, opts ...grpc.CallOption) (*GetActiveResponse, error) {
	out := new(GetActiveResponse)
	err := c.cc.Invoke(ctx, "/hootfs.discovery.DiscoverService/GetActive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoverServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/hootfs.discovery.DiscoverService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscoverServiceServer is the server API for DiscoverService service.
// All implementations must embed UnimplementedDiscoverServiceServer
// for forward compatibility
type DiscoverServiceServer interface {
	JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error)
	GetActive(context.Context, *GetActiveRequest) (*GetActiveResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedDiscoverServiceServer()
}

// UnimplementedDiscoverServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDiscoverServiceServer struct {
}

func (UnimplementedDiscoverServiceServer) JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (UnimplementedDiscoverServiceServer) GetActive(context.Context, *GetActiveRequest) (*GetActiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActive not implemented")
}
func (UnimplementedDiscoverServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDiscoverServiceServer) mustEmbedUnimplementedDiscoverServiceServer() {}

// UnsafeDiscoverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscoverServiceServer will
// result in compilation errors.
type UnsafeDiscoverServiceServer interface {
	mustEmbedUnimplementedDiscoverServiceServer()
}

func RegisterDiscoverServiceServer(s grpc.ServiceRegistrar, srv DiscoverServiceServer) {
	s.RegisterService(&DiscoverService_ServiceDesc, srv)
}

func _DiscoverService_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoverServiceServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hootfs.discovery.DiscoverService/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoverServiceServer).JoinCluster(ctx, req.(*JoinClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscoverService_GetActive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoverServiceServer).GetActive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hootfs.discovery.DiscoverService/GetActive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoverServiceServer).GetActive(ctx, req.(*GetActiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscoverService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoverServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hootfs.discovery.DiscoverService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoverServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DiscoverService_ServiceDesc is the grpc.ServiceDesc for DiscoverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DiscoverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hootfs.discovery.DiscoverService",
	HandlerType: (*DiscoverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinCluster",
			Handler:    _DiscoverService_JoinCluster_Handler,
		},
		{
			MethodName: "GetActive",
			Handler:    _DiscoverService_GetActive_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _DiscoverService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "discover.proto",
}
