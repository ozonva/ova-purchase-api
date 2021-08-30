// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ova_purchase_api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PurchaseServiceClient is the client API for PurchaseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PurchaseServiceClient interface {
	CreatePurchase(ctx context.Context, in *CreatePurchaseRequest, opts ...grpc.CallOption) (*CreatePurchaseResponse, error)
	DescribePurchase(ctx context.Context, in *DescribePurchaseRequest, opts ...grpc.CallOption) (*DescribePurchaseResponse, error)
	ListPurchases(ctx context.Context, in *ListPurchasesRequest, opts ...grpc.CallOption) (*ListPurchasesResponse, error)
	RemovePurchase(ctx context.Context, in *RemovePurchaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type purchaseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPurchaseServiceClient(cc grpc.ClientConnInterface) PurchaseServiceClient {
	return &purchaseServiceClient{cc}
}

func (c *purchaseServiceClient) CreatePurchase(ctx context.Context, in *CreatePurchaseRequest, opts ...grpc.CallOption) (*CreatePurchaseResponse, error) {
	out := new(CreatePurchaseResponse)
	err := c.cc.Invoke(ctx, "/PurchaseService/CreatePurchase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) DescribePurchase(ctx context.Context, in *DescribePurchaseRequest, opts ...grpc.CallOption) (*DescribePurchaseResponse, error) {
	out := new(DescribePurchaseResponse)
	err := c.cc.Invoke(ctx, "/PurchaseService/DescribePurchase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) ListPurchases(ctx context.Context, in *ListPurchasesRequest, opts ...grpc.CallOption) (*ListPurchasesResponse, error) {
	out := new(ListPurchasesResponse)
	err := c.cc.Invoke(ctx, "/PurchaseService/ListPurchases", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) RemovePurchase(ctx context.Context, in *RemovePurchaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/PurchaseService/RemovePurchase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PurchaseServiceServer is the server API for PurchaseService service.
// All implementations must embed UnimplementedPurchaseServiceServer
// for forward compatibility
type PurchaseServiceServer interface {
	CreatePurchase(context.Context, *CreatePurchaseRequest) (*CreatePurchaseResponse, error)
	DescribePurchase(context.Context, *DescribePurchaseRequest) (*DescribePurchaseResponse, error)
	ListPurchases(context.Context, *ListPurchasesRequest) (*ListPurchasesResponse, error)
	RemovePurchase(context.Context, *RemovePurchaseRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedPurchaseServiceServer()
}

// UnimplementedPurchaseServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPurchaseServiceServer struct {
}

func (UnimplementedPurchaseServiceServer) CreatePurchase(context.Context, *CreatePurchaseRequest) (*CreatePurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePurchase not implemented")
}
func (UnimplementedPurchaseServiceServer) DescribePurchase(context.Context, *DescribePurchaseRequest) (*DescribePurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribePurchase not implemented")
}
func (UnimplementedPurchaseServiceServer) ListPurchases(context.Context, *ListPurchasesRequest) (*ListPurchasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPurchases not implemented")
}
func (UnimplementedPurchaseServiceServer) RemovePurchase(context.Context, *RemovePurchaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePurchase not implemented")
}
func (UnimplementedPurchaseServiceServer) mustEmbedUnimplementedPurchaseServiceServer() {}

// UnsafePurchaseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PurchaseServiceServer will
// result in compilation errors.
type UnsafePurchaseServiceServer interface {
	mustEmbedUnimplementedPurchaseServiceServer()
}

func RegisterPurchaseServiceServer(s grpc.ServiceRegistrar, srv PurchaseServiceServer) {
	s.RegisterService(&PurchaseService_ServiceDesc, srv)
}

func _PurchaseService_CreatePurchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).CreatePurchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PurchaseService/CreatePurchase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).CreatePurchase(ctx, req.(*CreatePurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_DescribePurchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribePurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).DescribePurchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PurchaseService/DescribePurchase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).DescribePurchase(ctx, req.(*DescribePurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_ListPurchases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPurchasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).ListPurchases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PurchaseService/ListPurchases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).ListPurchases(ctx, req.(*ListPurchasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_RemovePurchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemovePurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).RemovePurchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PurchaseService/RemovePurchase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).RemovePurchase(ctx, req.(*RemovePurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PurchaseService_ServiceDesc is the grpc.ServiceDesc for PurchaseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PurchaseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PurchaseService",
	HandlerType: (*PurchaseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePurchase",
			Handler:    _PurchaseService_CreatePurchase_Handler,
		},
		{
			MethodName: "DescribePurchase",
			Handler:    _PurchaseService_DescribePurchase_Handler,
		},
		{
			MethodName: "ListPurchases",
			Handler:    _PurchaseService_ListPurchases_Handler,
		},
		{
			MethodName: "RemovePurchase",
			Handler:    _PurchaseService_RemovePurchase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/ova-purchase-api/ova-purchase-api.proto",
}
