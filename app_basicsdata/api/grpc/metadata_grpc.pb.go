// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: metadata.proto

package metadata

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

// MetadataServerClient is the client API for MetadataServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetadataServerClient interface {
	GetMetadata(ctx context.Context, in *GetMetadataRequest, opts ...grpc.CallOption) (*MetadataResponse, error)
	PutMetadata(ctx context.Context, in *MetadataRequest, opts ...grpc.CallOption) (*MetadataResponse, error)
}

type metadataServerClient struct {
	cc grpc.ClientConnInterface
}

func NewMetadataServerClient(cc grpc.ClientConnInterface) MetadataServerClient {
	return &metadataServerClient{cc}
}

func (c *metadataServerClient) GetMetadata(ctx context.Context, in *GetMetadataRequest, opts ...grpc.CallOption) (*MetadataResponse, error) {
	out := new(MetadataResponse)
	err := c.cc.Invoke(ctx, "/pb_metadata.MetadataServer/GetMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServerClient) PutMetadata(ctx context.Context, in *MetadataRequest, opts ...grpc.CallOption) (*MetadataResponse, error) {
	out := new(MetadataResponse)
	err := c.cc.Invoke(ctx, "/pb_metadata.MetadataServer/PutMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetadataServerServer is the server API for MetadataServer service.
// All implementations must embed UnimplementedMetadataServerServer
// for forward compatibility
type MetadataServerServer interface {
	GetMetadata(context.Context, *GetMetadataRequest) (*MetadataResponse, error)
	PutMetadata(context.Context, *MetadataRequest) (*MetadataResponse, error)
	mustEmbedUnimplementedMetadataServerServer()
}

// UnimplementedMetadataServerServer must be embedded to have forward compatible implementations.
type UnimplementedMetadataServerServer struct {
}

func (UnimplementedMetadataServerServer) GetMetadata(context.Context, *GetMetadataRequest) (*MetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetadata not implemented")
}
func (UnimplementedMetadataServerServer) PutMetadata(context.Context, *MetadataRequest) (*MetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutMetadata not implemented")
}
func (UnimplementedMetadataServerServer) mustEmbedUnimplementedMetadataServerServer() {}

// UnsafeMetadataServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetadataServerServer will
// result in compilation errors.
type UnsafeMetadataServerServer interface {
	mustEmbedUnimplementedMetadataServerServer()
}

func RegisterMetadataServerServer(s grpc.ServiceRegistrar, srv MetadataServerServer) {
	s.RegisterService(&MetadataServer_ServiceDesc, srv)
}

func _MetadataServer_GetMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServerServer).GetMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb_metadata.MetadataServer/GetMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServerServer).GetMetadata(ctx, req.(*GetMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataServer_PutMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServerServer).PutMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb_metadata.MetadataServer/PutMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServerServer).PutMetadata(ctx, req.(*MetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetadataServer_ServiceDesc is the grpc.ServiceDesc for MetadataServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetadataServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb_metadata.MetadataServer",
	HandlerType: (*MetadataServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMetadata",
			Handler:    _MetadataServer_GetMetadata_Handler,
		},
		{
			MethodName: "PutMetadata",
			Handler:    _MetadataServer_PutMetadata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metadata.proto",
}
